package ssm

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	actiontimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// _ action.Action is intentionally omitted: the factory function
// NewSsmServiceAction returns action.Action, which already performs the
// implicit compile-time check at the return statement. Only the
// extended-interface assertion is kept here to verify the Configure method
// signature (promoted from fw.ActionWithConfigure) at compile time.
var _ action.ActionWithConfigure = &SsmServiceAction{}

// NewSsmServiceAction is the factory referenced by
// tencentcloud/framework/registry.go to register this action.
//
// The "Action" suffix is required because the ssm package already declares
// NewSsmService(client *connectivity.TencentCloudClient) SsmService in
// service_tencentcloud_ssm.go; Go has no overloading, so a same-named factory
// would be a duplicate declaration.
func NewSsmServiceAction() action.Action {
	return &SsmServiceAction{}
}

// defaultSsmServiceInvokeTimeout is the built-in maximum duration to wait for
// the SSM service to become enabled after the activation order is placed.
// Practitioners can override it through the `timeouts` block, e.g.
// `timeouts { invoke = "30m" }`.
const defaultSsmServiceInvokeTimeout = 10 * time.Minute

// ssmServicePollInterval is the delay between two GetServiceStatus probes of
// the convergence loop. It is a package-level variable so unit tests can
// shorten it.
var ssmServicePollInterval = 5 * time.Second

// ssmServicePayModePostPay is the only PayMode accepted by the CreateService
// API.
const ssmServicePayModePostPay = "PostPay"

// ssmServiceClientTokenPattern enforces the ClientToken contract of the
// CreateService API: 1 to 64 visible ASCII characters ('!'..'~'), with no
// space, non-ASCII or control character.
var ssmServiceClientTokenPattern = regexp.MustCompile(`^[\x21-\x7E]{1,64}$`)

// SsmServiceAction implements action.Action for tencentcloud_ssm_service.
type SsmServiceAction struct {
	fw.ActionWithConfigure
}

// SsmServiceActionModel maps the schema attributes.
type SsmServiceActionModel struct {
	PayMode     types.String `tfsdk:"pay_mode"`
	ClientToken types.String `tfsdk:"client_token"`

	// Timeouts is the value type defined by the official
	// terraform-plugin-framework-timeouts module. Only the `invoke` stage
	// exists because an action has a single lifecycle stage.
	Timeouts actiontimeouts.Value `tfsdk:"timeouts"`
}

func (a *SsmServiceAction) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "tencentcloud_ssm_service"
}

func (a *SsmServiceAction) Schema(ctx context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an action to activate TencentCloud Secrets Manager (SSM) for the account and region " +
			"of the current provider configuration via the CreateService API, and then wait until " +
			"GetServiceStatus reports that the service is enabled. " +
			"SSM is an account-level, per-region service, so every other SSM resource fails until it is activated. " +
			"This is a one-time trigger: no Terraform state is persisted, no output attribute is exposed, " +
			"the action cannot be imported, and it must be run explicitly with `terraform apply -invoke=...`. " +
			"Each invoke places a pay-as-you-go activation order, so repeated invokes re-trigger the cloud call " +
			"(pin `client_token` to make retries replayable). Destroying the configuration does not deactivate SSM.",
		Attributes: map[string]schema.Attribute{
			"pay_mode": schema.StringAttribute{
				Required: true,
				Description: "Billing mode of the SSM activation order. Only `PostPay` (pay-as-you-go) is supported by " +
					"the CreateService API and the value is case sensitive.",
			},
			"client_token": schema.StringAttribute{
				Optional: true,
				Description: "Client-generated identifier of the activation request, used by the billing side to " +
					"deduplicate repeated creation requests. Must be 1 to 64 visible ASCII characters " +
					"('!'..'~'), with no space, non-ASCII or control character. " +
					"When omitted, the action generates `ssm-<random token>` on every invoke; provide an explicit " +
					"value to replay the same logical activation after a partial failure.",
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": ssmServiceTimeoutsBlock(ctx),
		},
	}
}

// ssmServiceTimeoutsBlock builds the `timeouts` block from the official
// terraform-plugin-framework-timeouts module, whose single optional `invoke`
// attribute is validated as a time.Duration at plan time. That module leaves
// the block itself undescribed (only the attribute carries a description), so
// the block description is restored here to keep the schema - and the
// documentation generated from it by gendoc - self-explanatory.
func ssmServiceTimeoutsBlock(ctx context.Context) schema.Block {
	block := actiontimeouts.BlockWithOpts(ctx, actiontimeouts.Opts{
		InvokeDescription: "A string that can be parsed as a duration (https://pkg.go.dev/time#ParseDuration), e.g. `10m`. " +
			"It limits how long the action waits for the SSM service to become enabled after the activation order " +
			"is placed. Default is `10m`.",
	})

	if nested, ok := block.(schema.SingleNestedBlock); ok {
		nested.Description = "The timeouts block allows you to specify the timeout for the invoke operation."
		return nested
	}

	return block
}

// Invoke is called to run the logic of the action.
func (a *SsmServiceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data SsmServiceActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.PayMode.IsNull() || data.PayMode.IsUnknown() || data.PayMode.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing pay_mode",
			"pay_mode is required and cannot be empty.",
		)
		return
	}

	payMode := data.PayMode.ValueString()
	if payMode != ssmServicePayModePostPay {
		resp.Diagnostics.AddError(
			"Invalid pay_mode",
			fmt.Sprintf("pay_mode (%s) is not supported: only `PostPay` is supported by the CreateService API.", payMode),
		)
		return
	}

	if !data.ClientToken.IsNull() && !data.ClientToken.IsUnknown() {
		clientToken := data.ClientToken.ValueString()
		if !ssmServiceClientTokenPattern.MatchString(clientToken) {
			resp.Diagnostics.AddError(
				"Invalid client_token",
				fmt.Sprintf("client_token (%s) is invalid: it must be 1 to 64 visible ASCII characters ('!'..'~'), "+
					"with no space, non-ASCII or control character.", clientToken),
			)
			return
		}
	}

	if a.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	// An unset/null timeouts block falls back to the built-in default, and a
	// malformed duration is reported as an error diagnostic. The schema
	// validator normally rejects malformed values at plan time; this is the
	// runtime safety net.
	timeout, timeoutDiags := data.Timeouts.Invoke(ctx, defaultSsmServiceInvokeTimeout)
	resp.Diagnostics.Append(timeoutDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The ClientToken is resolved exactly once, before the retry block, so
	// that every retry inside this invoke reuses the same value: the API
	// contract requires the same logical creation to reuse the same token.
	clientToken := ""
	if !data.ClientToken.IsNull() && !data.ClientToken.IsUnknown() && data.ClientToken.ValueString() != "" {
		clientToken = data.ClientToken.ValueString()
	} else {
		clientToken = fmt.Sprintf("ssm-%s", helper.BuildToken())
	}

	progress := func(msg string) {
		if resp.SendProgress != nil {
			resp.SendProgress(action.InvokeProgressEvent{Message: msg})
		}
	}

	progress(fmt.Sprintf("activating TencentCloud SSM service (pay_mode=%s, client_token=%s)", payMode, clientToken))

	service := NewSsmService(a.Client())
	var (
		created         *ssm.CreateServiceResponseParams
		alreadyEnabled  bool
		createStartTime = time.Now()
	)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		created = nil
		result, e := service.CreateSsmService(ctx, &payMode, &clientToken)
		if e != nil {
			// ResourceInUse.ServiceAlreadyEnabled means the service is already
			// at the target state: it is never worth retrying, and it must not
			// be treated as a failure either.
			if isSsmServiceAlreadyEnabled(e) {
				log.Printf("[WARN] SSM CreateService reported already enabled: %s", e.Error())
				alreadyEnabled = true
				return nil
			}
			return tccommon.RetryError(e)
		}
		created = result
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("activating SSM service (pay_mode=%s)", payMode),
			err.Error(),
		)
		return
	}

	if alreadyEnabled {
		// A previous invoke, a manual console activation or a retry after a
		// lost response already enabled the service: this is not a failure,
		// so continue to the convergence check.
		log.Printf("[WARN] SSM service already enabled after %s, skip the activation order", time.Since(createStartTime))
		progress("TencentCloud SSM service is already enabled, skip the activation order")
	}

	createdResourceId := ""
	if created != nil {
		progress(fmt.Sprintf("activation order submitted (resource_id=%s, order_id=%s)",
			ssmServiceStringValue(created.ResourceId), ssmServiceStringValue(created.OrderId)))
		log.Printf("[DEBUG] SSM CreateService succeeded, resource_id=%s, order_id=%s, request_id=%s",
			ssmServiceStringValue(created.ResourceId), ssmServiceStringValue(created.OrderId), ssmServiceStringValue(created.RequestId))
		createdResourceId = ssmServiceStringValue(created.ResourceId)
	}

	deadline := time.Now().Add(timeout)
	var (
		lastServiceEnabled   bool
		lastInvalidType      int64
		lastResourceId       string
		observedOnce         bool
		enabledAlreadyWarned bool
	)

	for {
		status, e := service.DescribeSsmServiceStatusByFilter(ctx)
		if e != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("querying SSM service status while waiting for activation (pay_mode=%s)", payMode),
				e.Error(),
			)
			return
		}

		if status != nil {
			observedOnce = true
			lastServiceEnabled = status.ServiceEnabled != nil && *status.ServiceEnabled
			if status.InvalidType != nil {
				lastInvalidType = *status.InvalidType
			} else {
				lastInvalidType = 0
			}
			lastResourceId = ssmServiceStringValue(status.ResourceId)

			log.Printf("[DEBUG] SSM GetServiceStatus: service_enabled=%v, invalid_type=%d, resource_id=%s, "+
				"access_key_escrow_enabled=%v, pay_model=%s, expire_time=%s, qps_limit=%d, secret_limit=%d, "+
				"renew_flag=%d, resource_region=%d, total_count=%d",
				lastServiceEnabled, lastInvalidType, lastResourceId,
				ssmServiceBoolValue(status.AccessKeyEscrowEnabled), ssmServiceStringValue(status.PayModel),
				ssmServiceStringValue(status.ExpireTime), ssmServiceInt64Value(status.QPSLimit), ssmServiceInt64Value(status.SecretLimit),
				ssmServiceInt64Value(status.RenewFlag), ssmServiceInt64Value(status.ResourceRegion), ssmServiceInt64Value(status.TotalCount))

			if lastServiceEnabled {
				if createdResourceId != "" && lastResourceId != "" && createdResourceId != lastResourceId &&
					!enabledAlreadyWarned {
					enabledAlreadyWarned = true
					resp.Diagnostics.AddWarning(
						"SSM resource_id mismatch",
						fmt.Sprintf("the resource_id returned by CreateService (%s) differs from the one reported by "+
							"GetServiceStatus (%s). The mapping is owned by the cloud side, so this is informational only.",
							createdResourceId, lastResourceId),
					)
				}

				progress(fmt.Sprintf("TencentCloud SSM service is enabled (resource_id=%s)", lastResourceId))
				return
			}

			// InvalidType 2 means suspended for arrears and 3 means the
			// resource was released: neither can be fixed by waiting, so fail
			// fast instead of burning the whole budget.
			if lastInvalidType == 2 || lastInvalidType == 3 {
				reason := "suspended due to arrears"
				if lastInvalidType == 3 {
					reason = "resource released"
				}
				resp.Diagnostics.AddError(
					fmt.Sprintf("activating SSM service (pay_mode=%s)", payMode),
					fmt.Sprintf("the SSM service cannot become enabled because it is %s (invalid_type=%d, resource_id=%s). "+
						"Please resolve the account state in the console and invoke the action again.",
						reason, lastInvalidType, lastResourceId),
				)
				return
			}
		}

		if time.Now().Add(ssmServicePollInterval).After(deadline) {
			if !observedOnce {
				resp.Diagnostics.AddError(
					"SSM service did not converge",
					fmt.Sprintf("the GetServiceStatus API did not return any status within %s "+
						"after the activation order was placed (pay_mode=%s, client_token=%s). "+
						"Please check the activation order state in the console.",
						timeout, payMode, clientToken),
				)
				return
			}

			resp.Diagnostics.AddError(
				"SSM service did not converge",
				fmt.Sprintf("the SSM service did not become enabled within %s. "+
					"last observed: service_enabled=%v, invalid_type=%d, resource_id=%s "+
					"(pay_mode=%s, client_token=%s).",
					timeout, lastServiceEnabled, lastInvalidType, lastResourceId, payMode, clientToken),
			)
			return
		}

		progress(fmt.Sprintf("waiting for TencentCloud SSM service to become enabled "+
			"(service_enabled=%v, invalid_type=%d)", lastServiceEnabled, lastInvalidType))

		time.Sleep(ssmServicePollInterval)
	}
}

// isSsmServiceAlreadyEnabled reports whether err signals that the SSM service
// is already activated, which makes a repeated invoke idempotent instead of a
// failure.
func isSsmServiceAlreadyEnabled(err error) bool {
	if err == nil {
		return false
	}

	if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
		if sdkErr.Code == ssm.RESOURCEINUSE_SERVICEALREADYENABLED {
			return true
		}
	}

	return regexp.MustCompile(regexp.QuoteMeta(ssm.RESOURCEINUSE_SERVICEALREADYENABLED)).MatchString(err.Error())
}

func ssmServiceStringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func ssmServiceBoolValue(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}

func ssmServiceInt64Value(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}
