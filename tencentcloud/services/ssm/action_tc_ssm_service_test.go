package ssm

import (
	"context"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

const (
	testSsmServicePayMode     = "PostPay"
	testSsmServiceClientToken = "ssm-4f4d0cef-0000-1111-2222-333344445555"
	testSsmServiceResourceId  = "ssm_1250000000"
	testSsmServiceOrderId     = "20250101120000001234567890"
)

func ssmServiceTestStringPtr(v string) *string { return &v }
func ssmServiceTestBoolPtr(v bool) *bool       { return &v }
func ssmServiceTestInt64Ptr(v int64) *int64    { return &v }

// buildSsmServiceConfig builds a tfsdk.Config for the action schema from the two raw input values.
// A nil pointer means the attribute is null.
func buildSsmServiceConfig(t *testing.T, s schema.Schema, payMode, clientToken *string) tfsdk.Config {
	t.Helper()

	payModeValue := tftypes.NewValue(tftypes.String, nil)
	if payMode != nil {
		payModeValue = tftypes.NewValue(tftypes.String, *payMode)
	}

	clientTokenValue := tftypes.NewValue(tftypes.String, nil)
	if clientToken != nil {
		clientTokenValue = tftypes.NewValue(tftypes.String, *clientToken)
	}

	return tfsdk.Config{
		Raw: tftypes.NewValue(
			tftypes.Object{AttributeTypes: map[string]tftypes.Type{
				"pay_mode":     tftypes.String,
				"client_token": tftypes.String,
			}},
			map[string]tftypes.Value{
				"pay_mode":     payModeValue,
				"client_token": clientTokenValue,
			},
		),
		Schema: s,
	}
}

func newSsmServiceTestMethodResponse() *ssm.GetServiceStatusResponse {
	resp := ssm.NewGetServiceStatusResponse()
	resp.Response = &ssm.GetServiceStatusResponseParams{
		ServiceEnabled:         ssmServiceTestBoolPtr(true),
		InvalidType:            ssmServiceTestInt64Ptr(1),
		AccessKeyEscrowEnabled: ssmServiceTestBoolPtr(true),
		ExpireTime:             ssmServiceTestStringPtr("2030-01-01 00:00:00"),
		QPSLimit:               ssmServiceTestInt64Ptr(100),
		SecretLimit:            ssmServiceTestInt64Ptr(1000),
		PayModel:               ssmServiceTestStringPtr(testSsmServicePayMode),
		RenewFlag:              ssmServiceTestInt64Ptr(0),
		ResourceId:             ssmServiceTestStringPtr(testSsmServiceResourceId),
		ResourceRegion:         ssmServiceTestInt64Ptr(1),
		RequestId:              ssmServiceTestStringPtr("fake-request-id"),
	}
	return resp
}

// newSsmServiceTestAction returns an action whose service holds a non-nil (but empty) client, so
// the "SSM client is not configured" guard does not trigger. The cloud API calls themselves are
// mocked with gomonkey, therefore the empty client is never dereferenced.
func newSsmServiceTestAction() *SsmServiceAction {
	return &SsmServiceAction{client: SsmService{client: &connectivity.TencentCloudClient{}}}
}

// collectSsmServiceProgress returns a SendProgress function that records every message.
func collectSsmServiceProgress(messages *[]string) func(action.InvokeProgressEvent) {
	return func(event action.InvokeProgressEvent) {
		*messages = append(*messages, event.Message)
	}
}

func TestSsmServiceAction_Metadata(t *testing.T) {
	a := &SsmServiceAction{}
	resp := &action.MetadataResponse{}

	a.Metadata(context.Background(), action.MetadataRequest{}, resp)

	assert.Equal(t, "tencentcloud_ssm_service", resp.TypeName)
}

func TestSsmServiceAction_Schema(t *testing.T) {
	a := &SsmServiceAction{}
	resp := &action.SchemaResponse{}

	a.Schema(context.Background(), action.SchemaRequest{}, resp)

	require.False(t, resp.Diagnostics.HasError(), "Schema returned diagnostics: %v", resp.Diagnostics)
	require.Len(t, resp.Schema.Attributes, 2, "action must declare exactly the two input attributes")

	payMode, ok := resp.Schema.Attributes["pay_mode"].(schema.StringAttribute)
	require.True(t, ok, "pay_mode must be a StringAttribute")
	assert.True(t, payMode.Required)
	assert.False(t, payMode.Optional)
	assert.NotEmpty(t, payMode.Validators)
	assert.NotEmpty(t, payMode.Description)

	clientToken, ok := resp.Schema.Attributes["client_token"].(schema.StringAttribute)
	require.True(t, ok, "client_token must be a StringAttribute")
	assert.True(t, clientToken.Required)
	assert.False(t, clientToken.Optional)
	assert.NotEmpty(t, clientToken.Validators)
	assert.NotEmpty(t, clientToken.Description)
}

func TestSsmServiceAction_SchemaValidators(t *testing.T) {
	a := &SsmServiceAction{}
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)

	validate := func(attribute string, value string) bool {
		attr, ok := schemaResp.Schema.Attributes[attribute].(schema.StringAttribute)
		if !ok {
			t.Fatalf("attribute %q is not a StringAttribute", attribute)
		}

		for _, v := range attr.Validators {
			resp := &validator.StringResponse{}
			v.ValidateString(context.Background(), validator.StringRequest{
				Path:        path.Root(attribute),
				ConfigValue: types.StringValue(value),
			}, resp)
			if resp.Diagnostics.HasError() {
				return false
			}
		}
		return true
	}

	assert.True(t, validate("pay_mode", testSsmServicePayMode))
	assert.False(t, validate("pay_mode", "PrePay"))
	assert.False(t, validate("pay_mode", "postpay"))

	assert.True(t, validate("client_token", testSsmServiceClientToken))
	assert.False(t, validate("client_token", ""))
	assert.False(t, validate("client_token", "ssm token with spaces"))
	assert.False(t, validate("client_token", "ssm-凭据"))
}

func TestSsmServiceAction_Configure(t *testing.T) {
	a := &SsmServiceAction{}

	nilResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), action.ConfigureRequest{ProviderData: nil}, nilResp)
	assert.False(t, nilResp.Diagnostics.HasError(), "nil ProviderData must return early without diagnostics")

	badResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), action.ConfigureRequest{ProviderData: "not-a-meta"}, badResp)
	require.True(t, badResp.Diagnostics.HasError())
	assert.Contains(t, badResp.Diagnostics[0].Summary(), "Unexpected ProviderData type")
}

func TestSsmServiceAction_Invoke_InvalidConfig(t *testing.T) {
	a := newSsmServiceTestAction()

	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	require.False(t, schemaResp.Diagnostics.HasError())

	cases := []struct {
		name        string
		payMode     *string
		clientToken *string
		summary     string
	}{
		{
			name:        "null pay_mode",
			payMode:     nil,
			clientToken: ssmServiceTestStringPtr(testSsmServiceClientToken),
			summary:     "Missing pay_mode",
		},
		{
			name:        "empty pay_mode",
			payMode:     ssmServiceTestStringPtr(""),
			clientToken: ssmServiceTestStringPtr(testSsmServiceClientToken),
			summary:     "Missing pay_mode",
		},
		{
			name:        "unsupported pay_mode",
			payMode:     ssmServiceTestStringPtr("PrePay"),
			clientToken: ssmServiceTestStringPtr(testSsmServiceClientToken),
			summary:     "Invalid pay_mode",
		},
		{
			name:        "null client_token",
			payMode:     ssmServiceTestStringPtr(testSsmServicePayMode),
			clientToken: nil,
			summary:     "Missing client_token",
		},
		{
			name:        "empty client_token",
			payMode:     ssmServiceTestStringPtr(testSsmServicePayMode),
			clientToken: ssmServiceTestStringPtr(""),
			summary:     "Missing client_token",
		},
		{
			name:        "too long client_token",
			payMode:     ssmServiceTestStringPtr(testSsmServicePayMode),
			clientToken: ssmServiceTestStringPtr("ssm-01234567890123456789012345678901234567890123456789012345678901234"),
			summary:     "Invalid client_token",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var messages []string
			resp := &action.InvokeResponse{SendProgress: collectSsmServiceProgress(&messages)}

			a.Invoke(context.Background(), action.InvokeRequest{
				Config: buildSsmServiceConfig(t, schemaResp.Schema, c.payMode, c.clientToken),
			}, resp)

			require.True(t, resp.Diagnostics.HasError())
			assert.Equal(t, c.summary, resp.Diagnostics[0].Summary())
			assert.Empty(t, messages, "no progress event must be emitted when validation fails")
		})
	}
}

func TestSsmServiceAction_Invoke_UnconfiguredClient(t *testing.T) {
	a := &SsmServiceAction{}

	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)

	var messages []string
	resp := &action.InvokeResponse{SendProgress: collectSsmServiceProgress(&messages)}

	a.Invoke(context.Background(), action.InvokeRequest{
		Config: buildSsmServiceConfig(t, schemaResp.Schema,
			ssmServiceTestStringPtr(testSsmServicePayMode),
			ssmServiceTestStringPtr(testSsmServiceClientToken)),
	}, resp)

	require.True(t, resp.Diagnostics.HasError())
	assert.Equal(t, "SSM client is not configured", resp.Diagnostics[0].Summary())
}

func TestSsmServiceAction_Invoke_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var capturedPayMode, capturedClientToken string
	createCalls := 0

	patches.ApplyMethodFunc(&SsmService{}, "CreateSsmService",
		func(_ context.Context, payMode, clientToken string) (*ssm.CreateServiceResponseParams, error) {
			createCalls++
			capturedPayMode = payMode
			capturedClientToken = clientToken
			return &ssm.CreateServiceResponseParams{
				ResourceId: ssmServiceTestStringPtr(testSsmServiceResourceId),
				OrderId:    ssmServiceTestStringPtr(testSsmServiceOrderId),
			}, nil
		})

	patches.ApplyMethodFunc(&SsmService{}, "DescribeSsmServiceStatusByFilter",
		func(_ context.Context) (*ssm.GetServiceStatusResponseParams, error) {
			return newSsmServiceTestMethodResponse().Response, nil
		})

	a := newSsmServiceTestAction()
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)

	var messages []string
	resp := &action.InvokeResponse{SendProgress: collectSsmServiceProgress(&messages)}

	a.Invoke(context.Background(), action.InvokeRequest{
		Config: buildSsmServiceConfig(t, schemaResp.Schema,
			ssmServiceTestStringPtr(testSsmServicePayMode),
			ssmServiceTestStringPtr(testSsmServiceClientToken)),
	}, resp)

	require.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)
	assert.Equal(t, 1, createCalls)
	assert.Equal(t, testSsmServicePayMode, capturedPayMode)
	assert.Equal(t, testSsmServiceClientToken, capturedClientToken)

	joined := joinSsmServiceMessages(messages)
	assert.Contains(t, joined, testSsmServiceResourceId)
	assert.Contains(t, joined, testSsmServiceOrderId)
	assert.Contains(t, joined, "service_enabled=true")
	assert.Contains(t, joined, "invalid_type=1")
	assert.Contains(t, joined, "secret_limit=1000")
	assert.Contains(t, joined, "resource_region=1")
}

func TestSsmServiceAction_Invoke_AlreadyEnabled(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	alreadyErr := sdkErrors.NewTencentCloudSDKError(ssm.RESOURCEINUSE_SERVICEALREADYENABLED,
		"service already enabled", "fake-request-id")

	patches.ApplyMethodReturn(&SsmService{}, "CreateSsmService",
		(*ssm.CreateServiceResponseParams)(nil), alreadyErr)

	statusCalls := 0
	patches.ApplyMethodFunc(&SsmService{}, "DescribeSsmServiceStatusByFilter",
		func(_ context.Context) (*ssm.GetServiceStatusResponseParams, error) {
			statusCalls++
			return newSsmServiceTestMethodResponse().Response, nil
		})

	a := newSsmServiceTestAction()
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)

	var messages []string
	resp := &action.InvokeResponse{SendProgress: collectSsmServiceProgress(&messages)}

	a.Invoke(context.Background(), action.InvokeRequest{
		Config: buildSsmServiceConfig(t, schemaResp.Schema,
			ssmServiceTestStringPtr(testSsmServicePayMode),
			ssmServiceTestStringPtr(testSsmServiceClientToken)),
	}, resp)

	require.False(t, resp.Diagnostics.HasError(),
		"ResourceInUse.ServiceAlreadyEnabled must not fail the action: %v", resp.Diagnostics)
	assert.Equal(t, 1, statusCalls, "the polling branch must still run")
	assert.Contains(t, joinSsmServiceMessages(messages), "service_enabled=true")
}

func TestSsmServiceAction_Invoke_CreateServiceError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	unauthorizedErr := sdkErrors.NewTencentCloudSDKError("UnauthorizedOperation",
		"no permission to subscribe SSM", "fake-request-id")

	patches.ApplyMethodReturn(&SsmService{}, "CreateSsmService",
		(*ssm.CreateServiceResponseParams)(nil), unauthorizedErr)

	statusCalled := false
	patches.ApplyMethodFunc(&SsmService{}, "DescribeSsmServiceStatusByFilter",
		func(_ context.Context) (*ssm.GetServiceStatusResponseParams, error) {
			statusCalled = true
			return newSsmServiceTestMethodResponse().Response, nil
		})

	a := newSsmServiceTestAction()
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)

	var messages []string
	resp := &action.InvokeResponse{SendProgress: collectSsmServiceProgress(&messages)}

	a.Invoke(context.Background(), action.InvokeRequest{
		Config: buildSsmServiceConfig(t, schemaResp.Schema,
			ssmServiceTestStringPtr(testSsmServicePayMode),
			ssmServiceTestStringPtr(testSsmServiceClientToken)),
	}, resp)

	require.True(t, resp.Diagnostics.HasError())
	assert.Equal(t, "Error subscribing to SSM", resp.Diagnostics[0].Summary())
	assert.Contains(t, resp.Diagnostics[0].Detail(), "CreateService")
	assert.False(t, statusCalled, "the polling branch must not run when CreateService fails")
}

func TestSsmServiceAction_Invoke_PollingTimeout(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodReturn(&SsmService{}, "CreateSsmService",
		&ssm.CreateServiceResponseParams{
			ResourceId: ssmServiceTestStringPtr(testSsmServiceResourceId),
			OrderId:    ssmServiceTestStringPtr(testSsmServiceOrderId),
		}, nil)

	notEnabled := &ssm.GetServiceStatusResponseParams{
		ServiceEnabled: ssmServiceTestBoolPtr(false),
		InvalidType:    ssmServiceTestInt64Ptr(0),
		ResourceId:     ssmServiceTestStringPtr(testSsmServiceResourceId),
	}
	patches.ApplyMethodReturn(&SsmService{}, "DescribeSsmServiceStatusByFilter", notEnabled, nil)

	originalInterval := ssmServiceStatusPollInterval
	originalTimeout := ssmServiceStatusPollTimeout
	ssmServiceStatusPollInterval = time.Millisecond
	ssmServiceStatusPollTimeout = 20 * time.Millisecond
	defer func() {
		ssmServiceStatusPollInterval = originalInterval
		ssmServiceStatusPollTimeout = originalTimeout
	}()

	a := newSsmServiceTestAction()
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)

	var messages []string
	resp := &action.InvokeResponse{SendProgress: collectSsmServiceProgress(&messages)}

	a.Invoke(context.Background(), action.InvokeRequest{
		Config: buildSsmServiceConfig(t, schemaResp.Schema,
			ssmServiceTestStringPtr(testSsmServicePayMode),
			ssmServiceTestStringPtr(testSsmServiceClientToken)),
	}, resp)

	require.True(t, resp.Diagnostics.HasError())
	assert.Equal(t, "SSM service subscription did not become enabled", resp.Diagnostics[0].Summary())
	assert.Contains(t, resp.Diagnostics[0].Detail(), "ServiceEnabled was still false")
	assert.Contains(t, resp.Diagnostics[0].Detail(), "invalid_type=0")
	assert.Contains(t, resp.Diagnostics[0].Detail(), testSsmServiceResourceId)
}

func TestSsmServiceAction_Invoke_EmptyStatusResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodReturn(&SsmService{}, "CreateSsmService",
		&ssm.CreateServiceResponseParams{
			ResourceId: ssmServiceTestStringPtr(testSsmServiceResourceId),
			OrderId:    ssmServiceTestStringPtr(testSsmServiceOrderId),
		}, nil)

	patches.ApplyMethodReturn(&SsmService{}, "DescribeSsmServiceStatusByFilter",
		(*ssm.GetServiceStatusResponseParams)(nil), nil)

	originalInterval := ssmServiceStatusPollInterval
	originalTimeout := ssmServiceStatusPollTimeout
	ssmServiceStatusPollInterval = time.Millisecond
	ssmServiceStatusPollTimeout = 20 * time.Millisecond
	defer func() {
		ssmServiceStatusPollInterval = originalInterval
		ssmServiceStatusPollTimeout = originalTimeout
	}()

	a := newSsmServiceTestAction()
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)

	var messages []string
	resp := &action.InvokeResponse{SendProgress: collectSsmServiceProgress(&messages)}

	a.Invoke(context.Background(), action.InvokeRequest{
		Config: buildSsmServiceConfig(t, schemaResp.Schema,
			ssmServiceTestStringPtr(testSsmServicePayMode),
			ssmServiceTestStringPtr(testSsmServiceClientToken)),
	}, resp)

	require.True(t, resp.Diagnostics.HasError())
	assert.Equal(t, "SSM service subscription did not become enabled", resp.Diagnostics[0].Summary())
	assert.Contains(t, resp.Diagnostics[0].Detail(), "empty response")
}

func TestFormatSsmServiceStatus_NilSafe(t *testing.T) {
	assert.Equal(t, "service_enabled=null", formatSsmServiceStatus(nil))

	empty := formatSsmServiceStatus(&ssm.GetServiceStatusResponseParams{})
	assert.Contains(t, empty, "service_enabled=null")
	assert.Contains(t, empty, "invalid_type=null")
	assert.Contains(t, empty, "access_key_escrow_enabled=null")
	assert.Contains(t, empty, "expire_time=null")
	assert.Contains(t, empty, "qps_limit=null")
	assert.Contains(t, empty, "secret_limit=null")
	assert.Contains(t, empty, "pay_model=null")
	assert.Contains(t, empty, "renew_flag=null")
	assert.Contains(t, empty, "resource_id=null")
	assert.Contains(t, empty, "resource_region=null")

	full := formatSsmServiceStatus(newSsmServiceTestMethodResponse().Response)
	assert.Contains(t, full, "service_enabled=true")
	assert.Contains(t, full, "qps_limit=100")
	assert.Contains(t, full, "pay_model=PostPay")
	assert.Contains(t, full, "expire_time=2030-01-01 00:00:00")
}

func TestIsSsmServiceAlreadyEnabled(t *testing.T) {
	assert.False(t, isSsmServiceAlreadyEnabled(nil))
	assert.True(t, isSsmServiceAlreadyEnabled(
		sdkErrors.NewTencentCloudSDKError(ssm.RESOURCEINUSE_SERVICEALREADYENABLED, "enabled", "id")))
	assert.False(t, isSsmServiceAlreadyEnabled(
		sdkErrors.NewTencentCloudSDKError("UnauthorizedOperation", "denied", "id")))
}

func joinSsmServiceMessages(messages []string) string {
	var joined string
	for _, message := range messages {
		joined += message + "\n"
	}
	return joined
}
