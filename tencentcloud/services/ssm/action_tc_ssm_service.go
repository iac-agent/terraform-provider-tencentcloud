package ssm

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

const (
	ssmServicePayModePostPay = "PostPay"
	ssmServiceClientTokenMin = 1
	ssmServiceClientTokenMax = 64
)

var (
	_ action.Action              = &SsmServiceAction{}
	_ action.ActionWithConfigure = &SsmServiceAction{}
)

// ssmServiceStatusPollInterval is the interval between two GetServiceStatus probes while the
// subscription is being provisioned. It is a variable so that unit tests can shorten it.
var ssmServiceStatusPollInterval = 5 * time.Second

// ssmServiceStatusPollTimeout is the total budget of the provisioning wait. It defaults to the
// provider-wide write retry timeout and is overridable by unit tests.
var ssmServiceStatusPollTimeout = tccommon.WriteRetryTimeout

// NewSsmService is the factory referenced by tencentcloud/framework/registry.go to register
// tencentcloud_ssm_service. The factory name deliberately carries no "Action" suffix.
func NewSsmService() action.Action {
	return &SsmServiceAction{}
}

// SsmServiceAction implements action.Action for tencentcloud_ssm_service. It subscribes the
// current account to Secrets Manager by calling CreateService and keeps no Terraform state.
type SsmServiceAction struct {
	client SsmService
}

// SsmServiceActionModel maps the action input attributes.
type SsmServiceActionModel struct {
	PayMode     types.String `tfsdk:"pay_mode"`
	ClientToken types.String `tfsdk:"client_token"`
}

// ssmServicePayModeValidator enforces the only billing mode supported by CreateService.
type ssmServicePayModeValidator struct{}

func (v ssmServicePayModeValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be %q", ssmServicePayModePostPay)
}

func (v ssmServicePayModeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v ssmServicePayModeValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if value := req.ConfigValue.ValueString(); value != ssmServicePayModePostPay {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid pay_mode",
			fmt.Sprintf("pay_mode must be %q (case-sensitive), got %q.", ssmServicePayModePostPay, value),
		)
	}
}

// ssmServiceClientTokenValidator enforces the 1 to 64 visible ASCII characters contract of
// CreateServiceRequest.ClientToken.
type ssmServiceClientTokenValidator struct{}

func (v ssmServiceClientTokenValidator) Description(_ context.Context) string {
	return fmt.Sprintf("value must be %d to %d visible ASCII characters", ssmServiceClientTokenMin, ssmServiceClientTokenMax)
}

func (v ssmServiceClientTokenValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v ssmServiceClientTokenValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	if len(value) < ssmServiceClientTokenMin || len(value) > ssmServiceClientTokenMax {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid client_token length",
			fmt.Sprintf("client_token must be %d to %d characters, got %d.",
				ssmServiceClientTokenMin, ssmServiceClientTokenMax, len(value)),
		)
		return
	}

	for _, c := range []byte(value) {
		if c < '!' || c > '~' {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid client_token character",
				"client_token must consist of visible ASCII characters (! to ~) only, "+
					"spaces, Chinese characters and control characters are not allowed.",
			)
			return
		}
	}
}

func (a *SsmServiceAction) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_ssm_service"
}

func (a *SsmServiceAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Subscribes the current account to Secrets Manager (SSM) by calling CreateService. " +
			"This is a one-time action: it keeps no Terraform state, every invoke performs a real " +
			"account-level API call, and no API can unsubscribe, so `terraform destroy` does not " +
			"revert the subscription.",
		Attributes: map[string]schema.Attribute{
			"pay_mode": schema.StringAttribute{
				Required: true,
				Description: "Billing mode of the SSM subscription. Only `PostPay` is supported " +
					"(case-sensitive).",
				Validators: []validator.String{ssmServicePayModeValidator{}},
			},
			"client_token": schema.StringAttribute{
				Required: true,
				Description: "Client-generated idempotency token of the subscription request, " +
					"1 to 64 visible ASCII characters. Reuse the same value when retrying the same " +
					"logical creation, otherwise a duplicated billing order may be created.",
				Validators: []validator.String{ssmServiceClientTokenValidator{}},
			},
		},
	}
}

// Configure receives the *sharedmeta.ProviderMeta populated by framework.Provider.Configure and
// constructs the SDKv2-style SsmService from the shared *connectivity.TencentCloudClient.
func (a *SsmServiceAction) Configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	meta, ok := req.ProviderData.(*sharedmeta.ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData type",
			fmt.Sprintf("Expected *sharedmeta.ProviderMeta, got: %T. "+
				"This is a bug in the provider; please report it.", req.ProviderData),
		)
		return
	}

	a.client = NewSsmServiceClient(meta.Client)
}

// Invoke subscribes the account to SSM and waits until the service reports itself enabled.
// It never writes state: the subscription snapshot is reported through progress messages.
func (a *SsmServiceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data SsmServiceActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.PayMode.IsNull() || data.PayMode.IsUnknown() || data.PayMode.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing pay_mode",
			"pay_mode is required and must not be empty. Only \"PostPay\" is supported by CreateService.",
		)
		return
	}

	if data.ClientToken.IsNull() || data.ClientToken.IsUnknown() || data.ClientToken.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing client_token",
			"client_token is required and must not be empty. Generate a value such as \"ssm-<uuid>\" and "+
				"reuse it when retrying the same logical creation.",
		)
		return
	}

	payMode := data.PayMode.ValueString()
	clientToken := data.ClientToken.ValueString()
	if payMode != ssmServicePayModePostPay {
		resp.Diagnostics.AddError(
			"Invalid pay_mode",
			fmt.Sprintf("pay_mode must be %q (case-sensitive), got %q.", ssmServicePayModePostPay, payMode),
		)
		return
	}

	if len(clientToken) < ssmServiceClientTokenMin || len(clientToken) > ssmServiceClientTokenMax {
		resp.Diagnostics.AddError(
			"Invalid client_token",
			fmt.Sprintf("client_token must be %d to %d characters, got %d.",
				ssmServiceClientTokenMin, ssmServiceClientTokenMax, len(clientToken)),
		)
		return
	}

	for _, c := range []byte(clientToken) {
		if c < '!' || c > '~' {
			resp.Diagnostics.AddError(
				"Invalid client_token",
				"client_token must consist of visible ASCII characters (! to ~) only, "+
					"spaces, Chinese characters and control characters are not allowed.",
			)
			return
		}
	}

	if a.client.client == nil {
		resp.Diagnostics.AddError(
			"SSM client is not configured",
			"The tencentcloud_ssm_service action was invoked before the provider configured its client. "+
				"Please report this as a provider bug.",
		)
		return
	}

	progress := func(message string) {
		if resp.SendProgress != nil {
			resp.SendProgress(action.InvokeProgressEvent{Message: message})
		}
	}

	log.Printf("[INFO] tencentcloud_ssm_service start to subscribe the account to SSM, pay_mode=%s", payMode)
	progress(fmt.Sprintf("Subscribing the current account to Secrets Manager (pay_mode=%s)", payMode))

	var (
		createService *ssm.CreateServiceResponseParams
		createErr     error
	)
	// client_token is captured once and reused unchanged on every retry, so that a network level
	// retry can never create a second billing order for the same logical creation.
	retryErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		createService, createErr = a.client.CreateSsmService(ctx, payMode, clientToken)
		if createErr != nil {
			if isSsmServiceAlreadyEnabled(createErr) {
				log.Printf("[INFO] tencentcloud_ssm_service is already enabled for the current account, skip CreateService error")
				return nil
			}
			return tccommon.RetryError(createErr)
		}
		return nil
	})
	if retryErr != nil && !isSsmServiceAlreadyEnabled(retryErr) {
		resp.Diagnostics.AddError(
			"Error subscribing to SSM",
			fmt.Sprintf("CreateService failed: %v", retryErr),
		)
		return
	}

	if createService != nil {
		resourceId := ssmServiceStringValue(createService.ResourceId)
		orderId := ssmServiceStringValue(createService.OrderId)
		log.Printf("[INFO] tencentcloud_ssm_service CreateService succeeded, resource_id=%s, order_id=%s", resourceId, orderId)
		progress(fmt.Sprintf("CreateService returned resource_id=%s, order_id=%s", resourceId, orderId))
	}

	status, err := a.waitForSsmServiceEnabled(ctx, progress)
	if err != nil {
		resp.Diagnostics.AddError("SSM service subscription did not become enabled", err.Error())
		return
	}

	snapshot := formatSsmServiceStatus(status)
	log.Printf("[INFO] tencentcloud_ssm_service is enabled, snapshot: %s", snapshot)
	progress(fmt.Sprintf("SSM service is enabled, status snapshot: %s", snapshot))
}

// waitForSsmServiceEnabled polls GetServiceStatus until ServiceEnabled is true or the polling
// budget is exhausted. CreateService is asynchronous, so its response does not imply that the
// service is usable yet.
func (a *SsmServiceAction) waitForSsmServiceEnabled(ctx context.Context, progress func(string)) (*ssm.GetServiceStatusResponseParams, error) {
	deadline := time.Now().Add(ssmServiceStatusPollTimeout)
	var lastStatus *ssm.GetServiceStatusResponseParams
	var lastErr error

	for {
		status, err := a.client.DescribeSsmServiceStatusByFilter(ctx)
		switch {
		case err != nil:
			lastErr = err
			log.Printf("[WARN] tencentcloud_ssm_service GetServiceStatus failed: %v", err)
		case status == nil:
			lastErr = fmt.Errorf("GetServiceStatus returned an empty response")
			log.Printf("[WARN] tencentcloud_ssm_service GetServiceStatus returned an empty response")
		default:
			lastErr = nil
			lastStatus = status
			if status.ServiceEnabled != nil && *status.ServiceEnabled {
				return status, nil
			}
		}

		if time.Now().After(deadline) {
			if lastErr != nil {
				return nil, fmt.Errorf("GetServiceStatus failed: %v", lastErr)
			}
			if lastStatus == nil {
				return nil, fmt.Errorf("GetServiceStatus never returned a usable response within %s", ssmServiceStatusPollTimeout)
			}
			return nil, fmt.Errorf("ServiceEnabled was still false after %s, last observed snapshot: %s",
				ssmServiceStatusPollTimeout, formatSsmServiceStatus(lastStatus))
		}

		progress(fmt.Sprintf("Waiting for the SSM subscription to become effective, last snapshot: %s",
			formatSsmServiceStatus(lastStatus)))

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled while waiting for the SSM subscription: %v", ctx.Err())
		case <-time.After(ssmServiceStatusPollInterval):
		}
	}
}

// isSsmServiceAlreadyEnabled reports whether err is the ResourceInUse.ServiceAlreadyEnabled code,
// which means the desired end state has already been reached.
func isSsmServiceAlreadyEnabled(err error) bool {
	if err == nil {
		return false
	}

	sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError)
	if !ok {
		// resource.Retry wraps the original error, therefore fall back to a plain match.
		return strings.Contains(err.Error(), ssm.RESOURCEINUSE_SERVICEALREADYENABLED)
	}

	code := sdkErr.GetCode()
	if code == ssm.RESOURCEINUSE_SERVICEALREADYENABLED {
		return true
	}

	// The API may return a longer code such as ResourceInUse.ServiceAlreadyEnabled.xxx.
	return strings.HasPrefix(code, ssm.RESOURCEINUSE_SERVICEALREADYENABLED)
}

func ssmServiceStringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

// formatSsmServiceStatus renders the read-only GetServiceStatus fields as a single line. Every
// pointer is nil-checked and rendered as null when absent.
func formatSsmServiceStatus(status *ssm.GetServiceStatusResponseParams) string {
	if status == nil {
		return "service_enabled=null"
	}

	toBool := func(value *bool) string {
		if value == nil {
			return "null"
		}
		return fmt.Sprintf("%t", *value)
	}
	toInt64 := func(value *int64) string {
		if value == nil {
			return "null"
		}
		return fmt.Sprintf("%d", *value)
	}
	toString := func(value *string) string {
		if value == nil {
			return "null"
		}
		return *value
	}

	return fmt.Sprintf(
		"service_enabled=%s invalid_type=%s access_key_escrow_enabled=%s expire_time=%s "+
			"qps_limit=%s secret_limit=%s pay_model=%s renew_flag=%s resource_id=%s resource_region=%s",
		toBool(status.ServiceEnabled),
		toInt64(status.InvalidType),
		toBool(status.AccessKeyEscrowEnabled),
		toString(status.ExpireTime),
		toInt64(status.QPSLimit),
		toInt64(status.SecretLimit),
		toString(status.PayModel),
		toInt64(status.RenewFlag),
		toString(status.ResourceId),
		toInt64(status.ResourceRegion),
	)
}
