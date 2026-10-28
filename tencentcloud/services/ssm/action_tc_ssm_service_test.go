package ssm_test

import (
	"context"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	actiontimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/action"
	action_schema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	ssmSDK "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ssm"
)

func ptrStringSsmService(s string) *string {
	return &s
}

func ptrBoolSsmService(b bool) *bool {
	return &b
}

func ptrInt64SsmService(i int64) *int64 {
	return &i
}

// ssmServiceSchema returns the framework action schema used to build the raw
// config value for InvokeRequest.
func ssmServiceSchema(t *testing.T) action_schema.Schema {
	a := ssm.NewSsmServiceAction()
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build action schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// ssmServiceTimeoutsObjType builds the tftypes.Object type for the timeouts
// block.
func ssmServiceTimeoutsObjType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"invoke": tftypes.String,
		},
	}
}

// nullSsmServiceTimeoutsValue is the null value for an unconfigured timeouts
// block.
func nullSsmServiceTimeoutsValue() tftypes.Value {
	return tftypes.NewValue(ssmServiceTimeoutsObjType(), nil)
}

// newSsmServiceInvokeRequest builds an InvokeRequest carrying the given
// pay_mode / client_token raw values and the given timeouts block value.
// A nil payModeVal means pay_mode is null; a nil clientTokenVal means
// client_token is null.
func newSsmServiceInvokeRequest(t *testing.T, payModeVal, clientTokenVal, timeoutsVal tftypes.Value) action.InvokeRequest {
	s := ssmServiceSchema(t)

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"pay_mode":     tftypes.String,
			"client_token": tftypes.String,
			"timeouts":     ssmServiceTimeoutsObjType(),
		},
	}, map[string]tftypes.Value{
		"pay_mode":     payModeVal,
		"client_token": clientTokenVal,
		"timeouts":     timeoutsVal,
	})

	req := action.InvokeRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

// newSsmServiceInvokeRequestWithToken builds an InvokeRequest with a concrete
// pay_mode and client_token.
func newSsmServiceInvokeRequestWithToken(t *testing.T, payMode, clientToken string) action.InvokeRequest {
	return newSsmServiceInvokeRequest(t,
		tftypes.NewValue(tftypes.String, payMode),
		tftypes.NewValue(tftypes.String, clientToken),
		nullSsmServiceTimeoutsValue(),
	)
}

// newSsmServiceInvokeRequestWithTimeouts builds an InvokeRequest with a
// concrete pay_mode and the given timeouts.invoke duration string.
func newSsmServiceInvokeRequestWithTimeouts(t *testing.T, payMode, invokeTimeout string) action.InvokeRequest {
	return newSsmServiceInvokeRequest(t,
		tftypes.NewValue(tftypes.String, payMode),
		tftypes.NewValue(tftypes.String, nil),
		tftypes.NewValue(ssmServiceTimeoutsObjType(), map[string]tftypes.Value{
			"invoke": tftypes.NewValue(tftypes.String, invokeTimeout),
		}),
	)
}

// setSsmServiceClient injects a mock connectivity client into the action
// instance via the Configure method.
func setSsmServiceClient(a *ssm.SsmServiceAction, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := action.ConfigureRequest{ProviderData: meta}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), configureReq, configureResp)
}

// go test ./tencentcloud/services/ssm/ -run "TestSsmServiceAction" -v -count=1 -gcflags="all=-l"

// TestSsmServiceAction_Metadata verifies the registered Terraform type name.
func TestSsmServiceAction_Metadata(t *testing.T) {
	a := ssm.NewSsmServiceAction()
	resp := &action.MetadataResponse{}
	a.Metadata(context.Background(), action.MetadataRequest{}, resp)
	assert.Equal(t, "tencentcloud_ssm_service", resp.TypeName)
}

// TestSsmServiceAction_Schema verifies the attribute set and that the timeouts
// block exposes the framework-timeouts custom type with an `invoke` attribute.
func TestSsmServiceAction_Schema(t *testing.T) {
	s := ssmServiceSchema(t)

	attrs := s.GetAttributes()
	assert.Contains(t, attrs, "pay_mode")
	assert.Contains(t, attrs, "client_token")

	payMode, ok := attrs["pay_mode"].(action_schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, payMode.Required)

	clientToken, ok := attrs["client_token"].(action_schema.StringAttribute)
	assert.True(t, ok)
	assert.True(t, clientToken.Optional)
	assert.False(t, clientToken.Required)

	blocks := s.GetBlocks()
	assert.Contains(t, blocks, "timeouts")
	timeoutsBlock, ok := blocks["timeouts"].(action_schema.SingleNestedBlock)
	assert.True(t, ok)
	assert.IsType(t, actiontimeouts.Type{}, timeoutsBlock.CustomType)
	_, ok = timeoutsBlock.Attributes["invoke"].(action_schema.StringAttribute)
	assert.True(t, ok)
}

// TestSsmServiceAction_Invoke_Success covers the happy path: the activation
// order is placed and the first status probe reports the service as enabled.
func TestSsmServiceAction_Invoke_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	ssmClient := &ssmSDK.Client{}
	patches.ApplyMethodReturn(client, "UseSsmClient", ssmClient)

	var gotCreateRequest *ssmSDK.CreateServiceRequest
	patches.ApplyMethodFunc(ssmClient, "CreateService", func(request *ssmSDK.CreateServiceRequest) (*ssmSDK.CreateServiceResponse, error) {
		gotCreateRequest = request
		resp := ssmSDK.NewCreateServiceResponse()
		resp.Response = &ssmSDK.CreateServiceResponseParams{
			ResourceId: ptrStringSsmService("ssm_1250000000"),
			OrderId:    ptrStringSsmService("2026010112000001"),
			RequestId:  ptrStringSsmService("fake-request-id"),
		}
		return resp, nil
	})
	patches.ApplyMethodFunc(ssmClient, "GetServiceStatus", func(_ *ssmSDK.GetServiceStatusRequest) (*ssmSDK.GetServiceStatusResponse, error) {
		resp := ssmSDK.NewGetServiceStatusResponse()
		resp.Response = &ssmSDK.GetServiceStatusResponseParams{
			ServiceEnabled: ptrBoolSsmService(true),
			InvalidType:    ptrInt64SsmService(1),
			ResourceId:     ptrStringSsmService("ssm_1250000000"),
			PayModel:       ptrStringSsmService("POSTPAY"),
			RequestId:      ptrStringSsmService("fake-request-id"),
		}
		return resp, nil
	})

	a := ssm.NewSsmServiceAction().(*ssm.SsmServiceAction)
	setSsmServiceClient(a, client)

	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), newSsmServiceInvokeRequestWithToken(t, "PostPay", "ssm-fixed-token"), resp)

	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, 0, resp.Diagnostics.WarningsCount())

	if assert.NotNil(t, gotCreateRequest) {
		assert.Equal(t, "PostPay", *gotCreateRequest.PayMode)
		assert.Equal(t, "ssm-fixed-token", *gotCreateRequest.ClientToken)
	}
}

// TestSsmServiceAction_Invoke_GeneratesClientToken verifies that an omitted
// client_token is generated with the `ssm-` prefix and stays within the 64
// character API limit.
func TestSsmServiceAction_Invoke_GeneratesClientToken(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	ssmClient := &ssmSDK.Client{}
	patches.ApplyMethodReturn(client, "UseSsmClient", ssmClient)

	var gotToken string
	patches.ApplyMethodFunc(ssmClient, "CreateService", func(request *ssmSDK.CreateServiceRequest) (*ssmSDK.CreateServiceResponse, error) {
		if request.ClientToken != nil {
			gotToken = *request.ClientToken
		}
		resp := ssmSDK.NewCreateServiceResponse()
		resp.Response = &ssmSDK.CreateServiceResponseParams{
			ResourceId: ptrStringSsmService("ssm_1250000000"),
		}
		return resp, nil
	})
	patches.ApplyMethodFunc(ssmClient, "GetServiceStatus", func(_ *ssmSDK.GetServiceStatusRequest) (*ssmSDK.GetServiceStatusResponse, error) {
		resp := ssmSDK.NewGetServiceStatusResponse()
		resp.Response = &ssmSDK.GetServiceStatusResponseParams{
			ServiceEnabled: ptrBoolSsmService(true),
			InvalidType:    ptrInt64SsmService(1),
			ResourceId:     ptrStringSsmService("ssm_1250000000"),
		}
		return resp, nil
	})

	a := ssm.NewSsmServiceAction().(*ssm.SsmServiceAction)
	setSsmServiceClient(a, client)

	resp := &action.InvokeResponse{}
	req := newSsmServiceInvokeRequest(t,
		tftypes.NewValue(tftypes.String, "PostPay"),
		tftypes.NewValue(tftypes.String, nil),
		nullSsmServiceTimeoutsValue(),
	)
	a.Invoke(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())
	assert.True(t, strings.HasPrefix(gotToken, "ssm-"))
	assert.LessOrEqual(t, len(gotToken), 64)
	assert.GreaterOrEqual(t, len(gotToken), 1)
}

// TestSsmServiceAction_Invoke_AlreadyEnabled covers the idempotent re-invoke
// path: CreateService reports ResourceInUse.ServiceAlreadyEnabled, which must
// not fail the action.
func TestSsmServiceAction_Invoke_AlreadyEnabled(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	ssmClient := &ssmSDK.Client{}
	patches.ApplyMethodReturn(client, "UseSsmClient", ssmClient)

	patches.ApplyMethodFunc(ssmClient, "CreateService", func(_ *ssmSDK.CreateServiceRequest) (*ssmSDK.CreateServiceResponse, error) {
		return nil, &sdkErrors.TencentCloudSDKError{
			Code:    ssmSDK.RESOURCEINUSE_SERVICEALREADYENABLED,
			Message: "service already enabled",
		}
	})
	patches.ApplyMethodFunc(ssmClient, "GetServiceStatus", func(_ *ssmSDK.GetServiceStatusRequest) (*ssmSDK.GetServiceStatusResponse, error) {
		resp := ssmSDK.NewGetServiceStatusResponse()
		resp.Response = &ssmSDK.GetServiceStatusResponseParams{
			ServiceEnabled: ptrBoolSsmService(true),
			InvalidType:    ptrInt64SsmService(1),
			ResourceId:     ptrStringSsmService("ssm_1250000000"),
		}
		return resp, nil
	})

	a := ssm.NewSsmServiceAction().(*ssm.SsmServiceAction)
	setSsmServiceClient(a, client)

	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), newSsmServiceInvokeRequestWithToken(t, "PostPay", "ssm-fixed-token"), resp)

	assert.False(t, resp.Diagnostics.HasError())
}

// TestSsmServiceAction_Invoke_CreateError verifies that a non-idempotent API
// failure surfaces as an error diagnostic.
func TestSsmServiceAction_Invoke_CreateError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	ssmClient := &ssmSDK.Client{}
	patches.ApplyMethodReturn(client, "UseSsmClient", ssmClient)

	patches.ApplyMethodFunc(ssmClient, "CreateService", func(_ *ssmSDK.CreateServiceRequest) (*ssmSDK.CreateServiceResponse, error) {
		return nil, &sdkErrors.TencentCloudSDKError{
			Code:    "FailedOperation",
			Message: "internal error",
		}
	})

	a := ssm.NewSsmServiceAction().(*ssm.SsmServiceAction)
	setSsmServiceClient(a, client)

	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), newSsmServiceInvokeRequestWithToken(t, "PostPay", "ssm-fixed-token"), resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "activating SSM service")
}

// TestSsmServiceAction_Invoke_StatusError verifies that a failing
// GetServiceStatus call surfaces as an error diagnostic.
func TestSsmServiceAction_Invoke_StatusError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	ssmClient := &ssmSDK.Client{}
	patches.ApplyMethodReturn(client, "UseSsmClient", ssmClient)

	patches.ApplyMethodFunc(ssmClient, "CreateService", func(_ *ssmSDK.CreateServiceRequest) (*ssmSDK.CreateServiceResponse, error) {
		resp := ssmSDK.NewCreateServiceResponse()
		resp.Response = &ssmSDK.CreateServiceResponseParams{
			ResourceId: ptrStringSsmService("ssm_1250000000"),
		}
		return resp, nil
	})
	patches.ApplyMethodFunc(ssmClient, "GetServiceStatus", func(_ *ssmSDK.GetServiceStatusRequest) (*ssmSDK.GetServiceStatusResponse, error) {
		return nil, &sdkErrors.TencentCloudSDKError{
			Code:    "InternalError",
			Message: "status query failed",
		}
	})

	a := ssm.NewSsmServiceAction().(*ssm.SsmServiceAction)
	setSsmServiceClient(a, client)

	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), newSsmServiceInvokeRequestWithToken(t, "PostPay", "ssm-fixed-token"), resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "querying SSM service status")
}

// TestSsmServiceAction_Invoke_InvalidArgs verifies that every input validation
// failure is reported before any cloud API call is made.
func TestSsmServiceAction_Invoke_InvalidArgs(t *testing.T) {
	testCases := []struct {
		name        string
		payMode     tftypes.Value
		clientToken tftypes.Value
		wantSummary string
	}{
		{
			name:        "missing pay_mode",
			payMode:     tftypes.NewValue(tftypes.String, nil),
			clientToken: tftypes.NewValue(tftypes.String, nil),
			wantSummary: "Missing pay_mode",
		},
		{
			name:        "empty pay_mode",
			payMode:     tftypes.NewValue(tftypes.String, ""),
			clientToken: tftypes.NewValue(tftypes.String, nil),
			wantSummary: "Missing pay_mode",
		},
		{
			name:        "unsupported pay_mode",
			payMode:     tftypes.NewValue(tftypes.String, "PrePay"),
			clientToken: tftypes.NewValue(tftypes.String, nil),
			wantSummary: "Invalid pay_mode",
		},
		{
			name:        "client_token with space",
			payMode:     tftypes.NewValue(tftypes.String, "PostPay"),
			clientToken: tftypes.NewValue(tftypes.String, "ssm token"),
			wantSummary: "Invalid client_token",
		},
		{
			name:        "client_token with non-ascii",
			payMode:     tftypes.NewValue(tftypes.String, "PostPay"),
			clientToken: tftypes.NewValue(tftypes.String, "ssm-密钥"),
			wantSummary: "Invalid client_token",
		},
		{
			name:        "client_token too long",
			payMode:     tftypes.NewValue(tftypes.String, "PostPay"),
			clientToken: tftypes.NewValue(tftypes.String, strings.Repeat("a", 65)),
			wantSummary: "Invalid client_token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			client := &connectivity.TencentCloudClient{}
			ssmClient := &ssmSDK.Client{}
			patches.ApplyMethodReturn(client, "UseSsmClient", ssmClient)

			apiCalled := false
			patches.ApplyMethodFunc(ssmClient, "CreateService", func(_ *ssmSDK.CreateServiceRequest) (*ssmSDK.CreateServiceResponse, error) {
				apiCalled = true
				return ssmSDK.NewCreateServiceResponse(), nil
			})

			a := ssm.NewSsmServiceAction().(*ssm.SsmServiceAction)
			setSsmServiceClient(a, client)

			resp := &action.InvokeResponse{}
			a.Invoke(context.Background(), newSsmServiceInvokeRequest(t, tc.payMode, tc.clientToken, nullSsmServiceTimeoutsValue()), resp)

			assert.True(t, resp.Diagnostics.HasError())
			assert.Equal(t, tc.wantSummary, resp.Diagnostics.Errors()[0].Summary())
			assert.False(t, apiCalled)
		})
	}
}

// TestSsmServiceAction_Invoke_NilClient verifies the guard for an
// unconfigured provider client.
func TestSsmServiceAction_Invoke_NilClient(t *testing.T) {
	a := ssm.NewSsmServiceAction().(*ssm.SsmServiceAction)

	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), newSsmServiceInvokeRequestWithToken(t, "PostPay", "ssm-fixed-token"), resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Equal(t, "Provider not configured", resp.Diagnostics.Errors()[0].Summary())
}

// TestSsmServiceAction_Invoke_FailFastOnSuspended verifies that a suspended
// (invalid_type=2) service fails immediately instead of burning the budget.
func TestSsmServiceAction_Invoke_FailFastOnSuspended(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	ssmClient := &ssmSDK.Client{}
	patches.ApplyMethodReturn(client, "UseSsmClient", ssmClient)

	patches.ApplyMethodFunc(ssmClient, "CreateService", func(_ *ssmSDK.CreateServiceRequest) (*ssmSDK.CreateServiceResponse, error) {
		resp := ssmSDK.NewCreateServiceResponse()
		resp.Response = &ssmSDK.CreateServiceResponseParams{
			ResourceId: ptrStringSsmService("ssm_1250000000"),
		}
		return resp, nil
	})

	statusCalls := 0
	patches.ApplyMethodFunc(ssmClient, "GetServiceStatus", func(_ *ssmSDK.GetServiceStatusRequest) (*ssmSDK.GetServiceStatusResponse, error) {
		statusCalls++
		resp := ssmSDK.NewGetServiceStatusResponse()
		resp.Response = &ssmSDK.GetServiceStatusResponseParams{
			ServiceEnabled: ptrBoolSsmService(false),
			InvalidType:    ptrInt64SsmService(2),
			ResourceId:     ptrStringSsmService("ssm_1250000000"),
		}
		return resp, nil
	})

	a := ssm.NewSsmServiceAction().(*ssm.SsmServiceAction)
	setSsmServiceClient(a, client)

	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), newSsmServiceInvokeRequestWithToken(t, "PostPay", "ssm-fixed-token"), resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "suspended due to arrears")
	assert.Equal(t, 1, statusCalls)
}

// TestSsmServiceAction_Invoke_ConvergenceTimeout verifies that a service that
// never becomes enabled ends with a timeout diagnostic carrying the last
// observed snapshot. A 1ms `timeouts.invoke` is used so the case finishes fast.
func TestSsmServiceAction_Invoke_ConvergenceTimeout(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	ssmClient := &ssmSDK.Client{}
	patches.ApplyMethodReturn(client, "UseSsmClient", ssmClient)

	patches.ApplyMethodFunc(ssmClient, "CreateService", func(_ *ssmSDK.CreateServiceRequest) (*ssmSDK.CreateServiceResponse, error) {
		resp := ssmSDK.NewCreateServiceResponse()
		resp.Response = &ssmSDK.CreateServiceResponseParams{
			ResourceId: ptrStringSsmService("ssm_1250000000"),
			OrderId:    ptrStringSsmService("2026010112000001"),
		}
		return resp, nil
	})
	patches.ApplyMethodFunc(ssmClient, "GetServiceStatus", func(_ *ssmSDK.GetServiceStatusRequest) (*ssmSDK.GetServiceStatusResponse, error) {
		resp := ssmSDK.NewGetServiceStatusResponse()
		resp.Response = &ssmSDK.GetServiceStatusResponseParams{
			ServiceEnabled: ptrBoolSsmService(false),
			InvalidType:    ptrInt64SsmService(1),
			ResourceId:     ptrStringSsmService("ssm_1250000000"),
		}
		return resp, nil
	})

	a := ssm.NewSsmServiceAction().(*ssm.SsmServiceAction)
	setSsmServiceClient(a, client)

	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), newSsmServiceInvokeRequestWithTimeouts(t, "PostPay", "1ms"), resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Equal(t, "SSM service did not converge", resp.Diagnostics.Errors()[0].Summary())
	detail := resp.Diagnostics.Errors()[0].Detail()
	assert.Contains(t, detail, "service_enabled=false")
	assert.Contains(t, detail, "invalid_type=1")
	assert.Contains(t, detail, "resource_id=ssm_1250000000")
}

// TestSsmServiceAction_Invoke_InvalidTimeout verifies that a malformed
// timeouts.invoke value is rejected before any API call.
func TestSsmServiceAction_Invoke_InvalidTimeout(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	ssmClient := &ssmSDK.Client{}
	patches.ApplyMethodReturn(client, "UseSsmClient", ssmClient)

	apiCalled := false
	patches.ApplyMethodFunc(ssmClient, "CreateService", func(_ *ssmSDK.CreateServiceRequest) (*ssmSDK.CreateServiceResponse, error) {
		apiCalled = true
		return ssmSDK.NewCreateServiceResponse(), nil
	})

	a := ssm.NewSsmServiceAction().(*ssm.SsmServiceAction)
	setSsmServiceClient(a, client)

	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), newSsmServiceInvokeRequestWithTimeouts(t, "PostPay", "not-a-duration"), resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Timeout Cannot Be Parsed")
	assert.False(t, apiCalled)
}

// TestSsmServiceAction_Invoke_ResourceIdMismatchWarns verifies that a
// divergence between the CreateService and GetServiceStatus resource_id is a
// warning, not a failure.
func TestSsmServiceAction_Invoke_ResourceIdMismatchWarns(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	ssmClient := &ssmSDK.Client{}
	patches.ApplyMethodReturn(client, "UseSsmClient", ssmClient)

	patches.ApplyMethodFunc(ssmClient, "CreateService", func(_ *ssmSDK.CreateServiceRequest) (*ssmSDK.CreateServiceResponse, error) {
		resp := ssmSDK.NewCreateServiceResponse()
		resp.Response = &ssmSDK.CreateServiceResponseParams{
			ResourceId: ptrStringSsmService("ssm_1250000000"),
		}
		return resp, nil
	})
	patches.ApplyMethodFunc(ssmClient, "GetServiceStatus", func(_ *ssmSDK.GetServiceStatusRequest) (*ssmSDK.GetServiceStatusResponse, error) {
		resp := ssmSDK.NewGetServiceStatusResponse()
		resp.Response = &ssmSDK.GetServiceStatusResponseParams{
			ServiceEnabled: ptrBoolSsmService(true),
			InvalidType:    ptrInt64SsmService(1),
			ResourceId:     ptrStringSsmService("ssm_1250000001"),
		}
		return resp, nil
	})

	a := ssm.NewSsmServiceAction().(*ssm.SsmServiceAction)
	setSsmServiceClient(a, client)

	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), newSsmServiceInvokeRequestWithToken(t, "PostPay", "ssm-fixed-token"), resp)

	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, 1, resp.Diagnostics.WarningsCount())
	assert.Equal(t, "SSM resource_id mismatch", resp.Diagnostics.Warnings()[0].Summary())
}
