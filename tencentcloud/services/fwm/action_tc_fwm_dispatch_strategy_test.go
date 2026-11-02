package fwm_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/action"
	action_schema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	fwm20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
)

func ptrStringFwmDispatchStrategy(s string) *string { return &s }

func ptrInt64FwmDispatchStrategy(v int64) *int64 { return &v }

func ptrFloat64FwmDispatchStrategy(v float64) *float64 { return &v }

// fwmDispatchStrategySchema returns the action schema used to build the raw
// config value for InvokeRequest.
func fwmDispatchStrategySchema(t *testing.T) action_schema.Schema {
	a := fwm.NewFwmDispatchStrategy()
	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build action schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// fwmDispatchStrategyObjectType returns the tftypes object describing the action
// config, including the optional timeouts block.
func fwmDispatchStrategyObjectType() tftypes.Type {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"product": tftypes.String,
			"status":  tftypes.Number,
			"timeouts": tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"invoke": tftypes.String,
				},
			},
		},
	}
}

// newFwmDispatchStrategyInvokeRequest builds an action.InvokeRequest with the
// given product and status. A nil product means the attribute is null, a nil
// status means the attribute is null.
func newFwmDispatchStrategyInvokeRequest(t *testing.T, product *string, status *int64) action.InvokeRequest {
	s := fwmDispatchStrategySchema(t)

	productVal := tftypes.NewValue(tftypes.String, nil)
	if product != nil {
		productVal = tftypes.NewValue(tftypes.String, *product)
	}

	statusVal := tftypes.NewValue(tftypes.Number, nil)
	if status != nil {
		statusVal = tftypes.NewValue(tftypes.Number, *status)
	}

	timeoutsVal := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"invoke": tftypes.String,
		},
	}, nil)

	raw := tftypes.NewValue(fwmDispatchStrategyObjectType(), map[string]tftypes.Value{
		"product":  productVal,
		"status":   statusVal,
		"timeouts": timeoutsVal,
	})

	req := action.InvokeRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

// setFwmDispatchStrategyClient injects a mock connectivity client into the action
// instance through the Configure method.
func setFwmDispatchStrategyClient(a *fwm.FwmDispatchStrategy, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := action.ConfigureRequest{ProviderData: meta}
	configureResp := &action.ConfigureResponse{}
	a.Configure(context.Background(), configureReq, configureResp)
}

// shortFwmDispatchStrategyPollInterval shortens the package level poll interval
// so that waiting tests do not sleep for the production value.
func shortFwmDispatchStrategyPollInterval(t *testing.T) {
	t.Helper()
	previous := fwm.SetFwmDispatchStrategyPollIntervalForTest(10 * time.Millisecond)
	t.Cleanup(func() {
		fwm.SetFwmDispatchStrategyPollIntervalForTest(previous)
	})
}

// go test ./tencentcloud/services/fwm/ -run "TestFwmDispatchStrategy" -v -count=1 -gcflags="all=-l"

// TestFwmDispatchStrategy_MetadataAndSchema validates metadata and schema.
func TestFwmDispatchStrategy_MetadataAndSchema(t *testing.T) {
	a := fwm.NewFwmDispatchStrategy()

	metaResp := &action.MetadataResponse{}
	a.Metadata(context.Background(), action.MetadataRequest{}, metaResp)
	if metaResp.TypeName != "tencentcloud_fwm_dispatch_strategy" {
		t.Fatalf("unexpected TypeName: %s", metaResp.TypeName)
	}

	schemaResp := &action.SchemaResponse{}
	a.Schema(context.Background(), action.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("unexpected schema diagnostics: %v", schemaResp.Diagnostics)
	}

	attrs := schemaResp.Schema.Attributes
	productAttr, ok := attrs["product"].(action_schema.StringAttribute)
	if !ok {
		t.Fatalf("product is not a string attribute: %T", attrs["product"])
	}
	if productAttr.Required || !productAttr.Optional {
		t.Fatalf("product must be optional, got required=%v optional=%v", productAttr.Required, productAttr.Optional)
	}
	if len(productAttr.Validators) != 1 {
		t.Fatalf("product must carry the non-empty validator, got %d validators", len(productAttr.Validators))
	}

	statusAttr, ok := attrs["status"].(action_schema.Int64Attribute)
	if !ok {
		t.Fatalf("status is not an int64 attribute: %T", attrs["status"])
	}
	if statusAttr.Required || !statusAttr.Optional {
		t.Fatalf("status must be optional, got required=%v optional=%v", statusAttr.Required, statusAttr.Optional)
	}
	if len(statusAttr.Validators) != 1 {
		t.Fatalf("status must carry the enum validator, got %d validators", len(statusAttr.Validators))
	}

	if _, ok := schemaResp.Schema.Blocks["timeouts"]; !ok {
		t.Fatalf("timeouts block is missing from the schema")
	}
}

// TestFwmDispatchStrategy_Invoke_DispatchSuccess covers the successful dispatch
// path: the request carries the configured scope and status, and the action waits
// until the dispatch status reaches success.
func TestFwmDispatchStrategy_Invoke_DispatchSuccess(t *testing.T) {
	shortFwmDispatchStrategyPollInterval(t)

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwm20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmClient", fwmClient)

	var gotDispatchRequest *fwm20250611.DispatchStrategyRequest
	patches.ApplyMethodFunc(fwmClient, "DispatchStrategyWithContext", func(_ context.Context, request *fwm20250611.DispatchStrategyRequest) (*fwm20250611.DispatchStrategyResponse, error) {
		gotDispatchRequest = request
		resp := &fwm20250611.DispatchStrategyResponse{
			Response: &fwm20250611.DispatchStrategyResponseParams{
				Status:    ptrInt64FwmDispatchStrategy(1),
				RequestId: ptrStringFwmDispatchStrategy("fake-request-id"),
			},
		}
		return resp, nil
	})

	statusCalls := 0
	var gotStatusProduct *string
	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyDispatchStatusWithContext", func(_ context.Context, request *fwm20250611.DescribeStrategyDispatchStatusRequest) (*fwm20250611.DescribeStrategyDispatchStatusResponse, error) {
		statusCalls++
		gotStatusProduct = request.Product

		status := int64(1)
		if statusCalls > 1 {
			status = 2
		}
		return &fwm20250611.DescribeStrategyDispatchStatusResponse{
			Response: &fwm20250611.DescribeStrategyDispatchStatusResponseParams{
				Progress:             ptrFloat64FwmDispatchStrategy(100),
				StartTime:            ptrStringFwmDispatchStrategy("2025-01-01 00:00:00"),
				EndTime:              ptrStringFwmDispatchStrategy("2025-01-01 00:01:00"),
				Status:               ptrInt64FwmDispatchStrategy(status),
				RuleGroupNum:         ptrInt64FwmDispatchStrategy(3),
				DispatchStrategyList: []*string{ptrStringFwmDispatchStrategy("strategy-1"), ptrStringFwmDispatchStrategy("strategy-2")},
				RequestId:            ptrStringFwmDispatchStrategy("fake-request-id"),
			},
		}, nil
	})

	a := fwm.NewFwmDispatchStrategy().(*fwm.FwmDispatchStrategy)
	setFwmDispatchStrategyClient(a, client)

	var progressMessages []string
	req := newFwmDispatchStrategyInvokeRequest(t, ptrStringFwmDispatchStrategy("cfw_edge_acl"), ptrInt64FwmDispatchStrategy(1))
	resp := &action.InvokeResponse{
		SendProgress: func(event action.InvokeProgressEvent) {
			progressMessages = append(progressMessages, event.Message)
		},
	}
	a.Invoke(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error diagnostics, got: %v", resp.Diagnostics)
	}
	if gotDispatchRequest == nil {
		t.Fatal("DispatchStrategy was not called")
	}
	if gotDispatchRequest.Product == nil || *gotDispatchRequest.Product != "cfw_edge_acl" {
		t.Fatalf("unexpected dispatch product: %v", gotDispatchRequest.Product)
	}
	if gotDispatchRequest.Status == nil || *gotDispatchRequest.Status != 1 {
		t.Fatalf("unexpected dispatch status: %v", gotDispatchRequest.Status)
	}
	if gotStatusProduct == nil || *gotStatusProduct != "cfw_edge_acl" {
		t.Fatalf("unexpected status query product: %v", gotStatusProduct)
	}
	if statusCalls < 2 {
		t.Fatalf("expected the action to poll until success, got %d status calls", statusCalls)
	}
	if len(progressMessages) == 0 {
		t.Fatal("expected progress events to be reported")
	}
	joined := strings.Join(progressMessages, " | ")
	if !strings.Contains(joined, "dispatch finished") {
		t.Fatalf("expected a final snapshot progress message, got: %s", joined)
	}
	if !strings.Contains(joined, "rule_group_num=3") {
		t.Fatalf("expected the final snapshot to carry rule_group_num, got: %s", joined)
	}
}

// TestFwmDispatchStrategy_Invoke_AllProducts covers the configuration where every
// optional attribute is omitted: no scope is sent and no client-side default
// status is printed as "0".
func TestFwmDispatchStrategy_Invoke_AllProducts(t *testing.T) {
	shortFwmDispatchStrategyPollInterval(t)

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwm20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmClient", fwmClient)

	var gotDispatchRequest *fwm20250611.DispatchStrategyRequest
	patches.ApplyMethodFunc(fwmClient, "DispatchStrategyWithContext", func(_ context.Context, request *fwm20250611.DispatchStrategyRequest) (*fwm20250611.DispatchStrategyResponse, error) {
		gotDispatchRequest = request
		return &fwm20250611.DispatchStrategyResponse{
			Response: &fwm20250611.DispatchStrategyResponseParams{Status: ptrInt64FwmDispatchStrategy(1)},
		}, nil
	})

	var gotStatusProduct *string
	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyDispatchStatusWithContext", func(_ context.Context, request *fwm20250611.DescribeStrategyDispatchStatusRequest) (*fwm20250611.DescribeStrategyDispatchStatusResponse, error) {
		gotStatusProduct = request.Product
		return &fwm20250611.DescribeStrategyDispatchStatusResponse{
			Response: &fwm20250611.DescribeStrategyDispatchStatusResponseParams{
				Status: ptrInt64FwmDispatchStrategy(2),
			},
		}, nil
	})

	a := fwm.NewFwmDispatchStrategy().(*fwm.FwmDispatchStrategy)
	setFwmDispatchStrategyClient(a, client)

	req := newFwmDispatchStrategyInvokeRequest(t, nil, nil)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error diagnostics, got: %v", resp.Diagnostics)
	}
	if gotDispatchRequest == nil {
		t.Fatal("DispatchStrategy was not called")
	}
	if gotDispatchRequest.Product != nil {
		t.Fatalf("product must be omitted, got: %v", *gotDispatchRequest.Product)
	}
	if gotDispatchRequest.Status == nil || *gotDispatchRequest.Status != 1 {
		t.Fatalf("status must default to 1, got: %v", gotDispatchRequest.Status)
	}
	if gotStatusProduct != nil {
		t.Fatalf("status query product must be omitted, got: %v", *gotStatusProduct)
	}
}

// TestFwmDispatchStrategy_Invoke_Abort covers status = 2: the abort request is
// sent and the action does not wait for a successful dispatch.
func TestFwmDispatchStrategy_Invoke_Abort(t *testing.T) {
	shortFwmDispatchStrategyPollInterval(t)

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwm20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmClient", fwmClient)

	var gotDispatchRequest *fwm20250611.DispatchStrategyRequest
	patches.ApplyMethodFunc(fwmClient, "DispatchStrategyWithContext", func(_ context.Context, request *fwm20250611.DispatchStrategyRequest) (*fwm20250611.DispatchStrategyResponse, error) {
		gotDispatchRequest = request
		return &fwm20250611.DispatchStrategyResponse{
			Response: &fwm20250611.DispatchStrategyResponseParams{Status: ptrInt64FwmDispatchStrategy(1)},
		}, nil
	})

	statusCalls := 0
	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyDispatchStatusWithContext", func(_ context.Context, request *fwm20250611.DescribeStrategyDispatchStatusRequest) (*fwm20250611.DescribeStrategyDispatchStatusResponse, error) {
		statusCalls++
		return &fwm20250611.DescribeStrategyDispatchStatusResponse{
			Response: &fwm20250611.DescribeStrategyDispatchStatusResponseParams{
				Status: ptrInt64FwmDispatchStrategy(1),
			},
		}, nil
	})

	a := fwm.NewFwmDispatchStrategy().(*fwm.FwmDispatchStrategy)
	setFwmDispatchStrategyClient(a, client)

	req := newFwmDispatchStrategyInvokeRequest(t, ptrStringFwmDispatchStrategy("cfw_edge_acl"), ptrInt64FwmDispatchStrategy(2))
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error diagnostics, got: %v", resp.Diagnostics)
	}
	if gotDispatchRequest == nil {
		t.Fatal("DispatchStrategy was not called")
	}
	if gotDispatchRequest.Status == nil || *gotDispatchRequest.Status != 2 {
		t.Fatalf("status must be 2, got: %v", gotDispatchRequest.Status)
	}
	if statusCalls != 0 {
		t.Fatalf("an abort must not wait for a successful dispatch, got %d status calls", statusCalls)
	}
}

// TestFwmDispatchStrategy_Invoke_DispatchFailed covers a terminal failure status
// reported by the cloud: the action fails with the cloud-side reason.
func TestFwmDispatchStrategy_Invoke_DispatchFailed(t *testing.T) {
	shortFwmDispatchStrategyPollInterval(t)

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwm20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmClient", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DispatchStrategyWithContext", func(_ context.Context, request *fwm20250611.DispatchStrategyRequest) (*fwm20250611.DispatchStrategyResponse, error) {
		return &fwm20250611.DispatchStrategyResponse{
			Response: &fwm20250611.DispatchStrategyResponseParams{Status: ptrInt64FwmDispatchStrategy(1)},
		}, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyDispatchStatusWithContext", func(_ context.Context, request *fwm20250611.DescribeStrategyDispatchStatusRequest) (*fwm20250611.DescribeStrategyDispatchStatusResponse, error) {
		return &fwm20250611.DescribeStrategyDispatchStatusResponse{
			Response: &fwm20250611.DescribeStrategyDispatchStatusResponseParams{
				Status:   ptrInt64FwmDispatchStrategy(3),
				ErrorMsg: ptrStringFwmDispatchStrategy("rule group is invalid"),
			},
		}, nil
	})

	a := fwm.NewFwmDispatchStrategy().(*fwm.FwmDispatchStrategy)
	setFwmDispatchStrategyClient(a, client)

	req := newFwmDispatchStrategyInvokeRequest(t, ptrStringFwmDispatchStrategy("cfw_edge_acl"), ptrInt64FwmDispatchStrategy(1))
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic for a failed dispatch")
	}
	if !strings.Contains(resp.Diagnostics.Errors()[0].Summary(), "waiting for FWM strategy dispatch") {
		t.Fatalf("unexpected error summary: %s", resp.Diagnostics.Errors()[0].Summary())
	}
	if !strings.Contains(resp.Diagnostics.Errors()[0].Detail(), "rule group is invalid") {
		t.Fatalf("expected the cloud-side reason in the detail, got: %s", resp.Diagnostics.Errors()[0].Detail())
	}
}

// TestFwmDispatchStrategy_Invoke_DispatchAPIError covers a DispatchStrategy API
// error: it is surfaced as a diagnostic and no polling happens.
func TestFwmDispatchStrategy_Invoke_DispatchAPIError(t *testing.T) {
	shortFwmDispatchStrategyPollInterval(t)

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwm20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmClient", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DispatchStrategyWithContext", func(_ context.Context, request *fwm20250611.DispatchStrategyRequest) (*fwm20250611.DispatchStrategyResponse, error) {
		return nil, errors.New("AuthFailure.UnauthorizedOperation")
	})

	statusCalls := 0
	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyDispatchStatusWithContext", func(_ context.Context, request *fwm20250611.DescribeStrategyDispatchStatusRequest) (*fwm20250611.DescribeStrategyDispatchStatusResponse, error) {
		statusCalls++
		return nil, errors.New("should not be called")
	})

	a := fwm.NewFwmDispatchStrategy().(*fwm.FwmDispatchStrategy)
	setFwmDispatchStrategyClient(a, client)

	req := newFwmDispatchStrategyInvokeRequest(t, ptrStringFwmDispatchStrategy("cfw_edge_acl"), ptrInt64FwmDispatchStrategy(1))
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic for a failed DispatchStrategy call")
	}
	if !strings.Contains(resp.Diagnostics.Errors()[0].Detail(), "AuthFailure.UnauthorizedOperation") {
		t.Fatalf("expected the API error in the detail, got: %s", resp.Diagnostics.Errors()[0].Detail())
	}
	if statusCalls != 0 {
		t.Fatalf("no status query may happen after a dispatch failure, got %d", statusCalls)
	}
}

// TestFwmDispatchStrategy_Invoke_StatusQueryError covers a nil Response payload:
// the status query is treated as an API failure, not as an empty result.
func TestFwmDispatchStrategy_Invoke_StatusQueryError(t *testing.T) {
	shortFwmDispatchStrategyPollInterval(t)

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwm20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmClient", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DispatchStrategyWithContext", func(_ context.Context, request *fwm20250611.DispatchStrategyRequest) (*fwm20250611.DispatchStrategyResponse, error) {
		return &fwm20250611.DispatchStrategyResponse{
			Response: &fwm20250611.DispatchStrategyResponseParams{Status: ptrInt64FwmDispatchStrategy(1)},
		}, nil
	})

	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyDispatchStatusWithContext", func(_ context.Context, request *fwm20250611.DescribeStrategyDispatchStatusRequest) (*fwm20250611.DescribeStrategyDispatchStatusResponse, error) {
		return &fwm20250611.DescribeStrategyDispatchStatusResponse{}, nil
	})

	a := fwm.NewFwmDispatchStrategy().(*fwm.FwmDispatchStrategy)
	setFwmDispatchStrategyClient(a, client)

	req := newFwmDispatchStrategyInvokeRequest(t, ptrStringFwmDispatchStrategy("cfw_edge_acl"), ptrInt64FwmDispatchStrategy(1))
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic when the status response is nil")
	}
	if !strings.Contains(resp.Diagnostics.Errors()[0].Detail(), "Response is nil") {
		t.Fatalf("expected a nil response error, got: %s", resp.Diagnostics.Errors()[0].Detail())
	}
}

// TestFwmDispatchStrategy_Invoke_MissingProduct covers an empty product, which
// must be rejected before any API call.
func TestFwmDispatchStrategy_Invoke_MissingProduct(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwm20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmClient", fwmClient)

	apiCalled := false
	patches.ApplyMethodFunc(fwmClient, "DispatchStrategyWithContext", func(_ context.Context, request *fwm20250611.DispatchStrategyRequest) (*fwm20250611.DispatchStrategyResponse, error) {
		apiCalled = true
		return &fwm20250611.DispatchStrategyResponse{
			Response: &fwm20250611.DispatchStrategyResponseParams{Status: ptrInt64FwmDispatchStrategy(1)},
		}, nil
	})

	a := fwm.NewFwmDispatchStrategy().(*fwm.FwmDispatchStrategy)
	setFwmDispatchStrategyClient(a, client)

	req := newFwmDispatchStrategyInvokeRequest(t, ptrStringFwmDispatchStrategy(""), nil)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic for an empty product")
	}
	if !strings.Contains(resp.Diagnostics.Errors()[0].Summary(), "Invalid product") {
		t.Fatalf("unexpected error summary: %s", resp.Diagnostics.Errors()[0].Summary())
	}
	if apiCalled {
		t.Fatal("no API call may happen with an empty product")
	}
}

// TestFwmDispatchStrategy_Invoke_InvalidStatus covers a status value outside the
// documented enum.
func TestFwmDispatchStrategy_Invoke_InvalidStatus(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwm20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmClient", fwmClient)

	apiCalled := false
	patches.ApplyMethodFunc(fwmClient, "DispatchStrategyWithContext", func(_ context.Context, request *fwm20250611.DispatchStrategyRequest) (*fwm20250611.DispatchStrategyResponse, error) {
		apiCalled = true
		return &fwm20250611.DispatchStrategyResponse{
			Response: &fwm20250611.DispatchStrategyResponseParams{Status: ptrInt64FwmDispatchStrategy(1)},
		}, nil
	})

	a := fwm.NewFwmDispatchStrategy().(*fwm.FwmDispatchStrategy)
	setFwmDispatchStrategyClient(a, client)

	req := newFwmDispatchStrategyInvokeRequest(t, ptrStringFwmDispatchStrategy("cfw_edge_acl"), ptrInt64FwmDispatchStrategy(3))
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic for an unsupported status")
	}
	if !strings.Contains(resp.Diagnostics.Errors()[0].Summary(), "Invalid status") {
		t.Fatalf("unexpected error summary: %s", resp.Diagnostics.Errors()[0].Summary())
	}
	if apiCalled {
		t.Fatal("no API call may happen with an unsupported status")
	}
}

// TestFwmDispatchStrategy_Invoke_AbortWithoutProduct covers the cross-field rule:
// aborting needs an explicit product scope.
func TestFwmDispatchStrategy_Invoke_AbortWithoutProduct(t *testing.T) {
	a := fwm.NewFwmDispatchStrategy().(*fwm.FwmDispatchStrategy)

	req := newFwmDispatchStrategyInvokeRequest(t, nil, ptrInt64FwmDispatchStrategy(2))
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic when aborting without a product")
	}
	if !strings.Contains(resp.Diagnostics.Errors()[0].Summary(), "Missing product") {
		t.Fatalf("unexpected error summary: %s", resp.Diagnostics.Errors()[0].Summary())
	}
}

// TestFwmDispatchStrategy_Invoke_ClientNotConfigured covers the nil-client guard.
func TestFwmDispatchStrategy_Invoke_ClientNotConfigured(t *testing.T) {
	a := fwm.NewFwmDispatchStrategy().(*fwm.FwmDispatchStrategy)

	req := newFwmDispatchStrategyInvokeRequest(t, ptrStringFwmDispatchStrategy("cfw_edge_acl"), nil)
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic when the client is not configured")
	}
	if !strings.Contains(resp.Diagnostics.Errors()[0].Summary(), "Provider not configured") {
		t.Fatalf("unexpected error summary: %s", resp.Diagnostics.Errors()[0].Summary())
	}
}

// TestFwmDispatchStrategy_ValidateConfig covers the offline validation rules.
func TestFwmDispatchStrategy_ValidateConfig(t *testing.T) {
	a := fwm.NewFwmDispatchStrategy().(*fwm.FwmDispatchStrategy)

	req := action.ValidateConfigRequest{}
	req.Config.Raw = tftypes.NewValue(fwmDispatchStrategyObjectType(), map[string]tftypes.Value{
		"product": tftypes.NewValue(tftypes.String, ""),
		"status":  tftypes.NewValue(tftypes.Number, int64(1)),
		"timeouts": tftypes.NewValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{"invoke": tftypes.String},
		}, nil),
	})
	req.Config.Schema = fwmDispatchStrategySchema(t)

	resp := &action.ValidateConfigResponse{}
	a.ValidateConfig(context.Background(), req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error diagnostic for an empty product")
	}
	if !strings.Contains(resp.Diagnostics.Errors()[0].Summary(), "Invalid product") {
		t.Fatalf("unexpected error summary: %s", resp.Diagnostics.Errors()[0].Summary())
	}
}

// TestFwmDispatchStrategy_StatusValidator covers the status enum validator through
// the schema attribute declaration.
func TestFwmDispatchStrategy_StatusValidator(t *testing.T) {
	s := fwmDispatchStrategySchema(t)

	attr, ok := s.Attributes["status"].(action_schema.Int64Attribute)
	if !ok {
		t.Fatalf("status is not an int64 attribute: %T", s.Attributes["status"])
	}
	if len(attr.Validators) != 1 {
		t.Fatalf("expected exactly one validator on status, got %d", len(attr.Validators))
	}

	for _, value := range []int64{1, 2} {
		resp := &validator.Int64Response{}
		attr.Validators[0].ValidateInt64(context.Background(), validator.Int64Request{
			Path:        path.Root("status"),
			ConfigValue: types.Int64Value(value),
		}, resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("status %d must be accepted, got: %v", value, resp.Diagnostics)
		}
	}

	resp := &validator.Int64Response{}
	attr.Validators[0].ValidateInt64(context.Background(), validator.Int64Request{
		Path:        path.Root("status"),
		ConfigValue: types.Int64Value(9),
	}, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected a validation error for status = 9")
	}
	if !strings.Contains(resp.Diagnostics.Errors()[0].Detail(), "status must be 1 (dispatch) or 2 (abort)") {
		t.Fatalf("unexpected validation error: %s", resp.Diagnostics.Errors()[0].Detail())
	}
}

// TestFwmDispatchStrategy_ProductValidator covers the non-empty product validator.
func TestFwmDispatchStrategy_ProductValidator(t *testing.T) {
	s := fwmDispatchStrategySchema(t)

	attr, ok := s.Attributes["product"].(action_schema.StringAttribute)
	if !ok {
		t.Fatalf("product is not a string attribute: %T", s.Attributes["product"])
	}
	if len(attr.Validators) != 1 {
		t.Fatalf("expected exactly one validator on product, got %d", len(attr.Validators))
	}

	resp := &validator.StringResponse{}
	attr.Validators[0].ValidateString(context.Background(), validator.StringRequest{
		Path:        path.Root("product"),
		ConfigValue: types.StringValue("cfw_edge_acl"),
	}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("a non-empty product must be accepted, got: %v", resp.Diagnostics)
	}

	resp = &validator.StringResponse{}
	attr.Validators[0].ValidateString(context.Background(), validator.StringRequest{
		Path:        path.Root("product"),
		ConfigValue: types.StringValue(""),
	}, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected a validation error for an empty product")
	}
}

// TestFwmDispatchStrategy_Wait_Timeout covers the timeout branch: the wait reports
// the last observed snapshot instead of a bare timeout.
func TestFwmDispatchStrategy_Wait_Timeout(t *testing.T) {
	shortFwmDispatchStrategyPollInterval(t)

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwm20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmClient", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyDispatchStatusWithContext", func(_ context.Context, request *fwm20250611.DescribeStrategyDispatchStatusRequest) (*fwm20250611.DescribeStrategyDispatchStatusResponse, error) {
		return &fwm20250611.DescribeStrategyDispatchStatusResponse{
			Response: &fwm20250611.DescribeStrategyDispatchStatusResponseParams{
				Progress:     ptrFloat64FwmDispatchStrategy(42.5),
				StartTime:    ptrStringFwmDispatchStrategy("2025-01-01 00:00:00"),
				Status:       ptrInt64FwmDispatchStrategy(4),
				RuleGroupNum: ptrInt64FwmDispatchStrategy(5),
			},
		}, nil
	})

	service := fwm.NewFwmService(client)
	err := service.WaitFwmStrategyDispatch(context.Background(), "cfw_edge_acl", 60*time.Millisecond, 10*time.Millisecond)
	if err == nil {
		t.Fatal("expected a timeout error")
	}
	for _, want := range []string{"did not converge", "status=4", "progress=42.5", "rule_group_num=5"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("expected %q in the timeout error, got: %s", want, err.Error())
		}
	}
}

// TestFwmDispatchStrategy_Wait_NilStatus covers an unset status field, which keeps
// the wait running until the timeout.
func TestFwmDispatchStrategy_Wait_NilStatus(t *testing.T) {
	shortFwmDispatchStrategyPollInterval(t)

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwm20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmClient", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyDispatchStatusWithContext", func(_ context.Context, request *fwm20250611.DescribeStrategyDispatchStatusRequest) (*fwm20250611.DescribeStrategyDispatchStatusResponse, error) {
		return &fwm20250611.DescribeStrategyDispatchStatusResponse{
			Response: &fwm20250611.DescribeStrategyDispatchStatusResponseParams{},
		}, nil
	})

	service := fwm.NewFwmService(client)
	err := service.WaitFwmStrategyDispatch(context.Background(), "", 40*time.Millisecond, 10*time.Millisecond)
	if err == nil {
		t.Fatal("expected a timeout error for a nil status")
	}
	if !strings.Contains(err.Error(), "status=<nil>") {
		t.Fatalf("expected the nil status snapshot in the error, got: %s", err.Error())
	}
}

// TestFwmDispatchStrategy_Wait_FailedWithoutMessage covers the failure branch when
// the cloud reports no error message.
func TestFwmDispatchStrategy_Wait_FailedWithoutMessage(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwm20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmClient", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyDispatchStatusWithContext", func(_ context.Context, request *fwm20250611.DescribeStrategyDispatchStatusRequest) (*fwm20250611.DescribeStrategyDispatchStatusResponse, error) {
		return &fwm20250611.DescribeStrategyDispatchStatusResponse{
			Response: &fwm20250611.DescribeStrategyDispatchStatusResponseParams{
				Status: ptrInt64FwmDispatchStrategy(3),
			},
		}, nil
	})

	service := fwm.NewFwmService(client)
	err := service.WaitFwmStrategyDispatch(context.Background(), "cfw_edge_acl", time.Second, 10*time.Millisecond)
	if err == nil {
		t.Fatal("expected a failure error")
	}
	if !strings.Contains(err.Error(), "the cloud did not report an error message") {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

// TestFwmDispatchStrategy_Invoke_ProgressNilSafe covers that Invoke does not panic
// when SendProgress is nil.
func TestFwmDispatchStrategy_Invoke_ProgressNilSafe(t *testing.T) {
	shortFwmDispatchStrategyPollInterval(t)

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	fwmClient := &fwm20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmClient", fwmClient)

	patches.ApplyMethodFunc(fwmClient, "DispatchStrategyWithContext", func(_ context.Context, request *fwm20250611.DispatchStrategyRequest) (*fwm20250611.DispatchStrategyResponse, error) {
		return &fwm20250611.DispatchStrategyResponse{
			Response: &fwm20250611.DispatchStrategyResponseParams{Status: ptrInt64FwmDispatchStrategy(1)},
		}, nil
	})
	patches.ApplyMethodFunc(fwmClient, "DescribeStrategyDispatchStatusWithContext", func(_ context.Context, request *fwm20250611.DescribeStrategyDispatchStatusRequest) (*fwm20250611.DescribeStrategyDispatchStatusResponse, error) {
		return &fwm20250611.DescribeStrategyDispatchStatusResponse{
			Response: &fwm20250611.DescribeStrategyDispatchStatusResponseParams{Status: ptrInt64FwmDispatchStrategy(2)},
		}, nil
	})

	a := fwm.NewFwmDispatchStrategy().(*fwm.FwmDispatchStrategy)
	setFwmDispatchStrategyClient(a, client)

	req := newFwmDispatchStrategyInvokeRequest(t, ptrStringFwmDispatchStrategy("cfw_edge_acl"), ptrInt64FwmDispatchStrategy(1))
	resp := &action.InvokeResponse{}
	a.Invoke(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error diagnostics, got: %v", resp.Diagnostics)
	}
}
