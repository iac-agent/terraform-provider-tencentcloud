package alb_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcalb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/alb"
)

const (
	albTestLoadBalancerId = "alb-12345678"
	albTestListenerId     = "lst-87654321"
	albTestTargetGroupId  = "lbtg-11111111"
)

func albString(value string) *string { return &value }

func albBool(value bool) *bool { return &value }

func albUint64(value uint64) *uint64 { return &value }

type albListenerTestContext struct {
	resource resource.Resource
	client   *albv20251030.Client
	schema   schema.Schema
}

func newAlbListenerTestContext(t *testing.T, patches *gomonkey.Patches) *albListenerTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}
	albClient := &albv20251030.Client{}
	patches.ApplyMethodReturn(connectivityClient, "UseAlbV20251030Client", albClient)

	res := svcalb.NewAlbListenerResource()
	configured, ok := res.(resource.ResourceWithConfigure)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithConfigure")
	}

	var configureResp resource.ConfigureResponse
	configured.Configure(
		context.Background(),
		resource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: connectivityClient}},
		&configureResp,
	)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("failed to configure resource: %v", configureResp.Diagnostics.Errors())
	}

	var schemaResp resource.SchemaResponse
	res.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build schema: %v", schemaResp.Diagnostics.Errors())
	}

	return &albListenerTestContext{resource: res, client: albClient, schema: schemaResp.Schema}
}

func (c *albListenerTestContext) plan(t *testing.T, model svcalb.AlbListenerResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}
	return plan
}

func (c *albListenerTestContext) state(t *testing.T, model svcalb.AlbListenerResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}
	return state
}

func albSampleListener() *albv20251030.ListenerOutput {
	return &albv20251030.ListenerOutput{
		ListenerId:       albString(albTestListenerId),
		ListenerName:     albString("tf-example"),
		ListenerPort:     albUint64(80),
		ListenerProtocol: albString("HTTP"),
		ListenerStatus:   albString("Active"),
		GzipEnabled:      albBool(true),
		IdleTimeout:      albUint64(30),
		RequestTimeout:   albUint64(60),
		Tags: []*albv20251030.TagInfo{
			{TagKey: albString("createdBy"), TagValue: albString("terraform")},
		},
		XForwardedForConfig: &albv20251030.XForwardedForConfig{
			XForwardedForMode: albString("append"),
		},
	}
}

func albDescribeResponse(listeners ...*albv20251030.ListenerOutput) *albv20251030.DescribeListenersResponse {
	return &albv20251030.DescribeListenersResponse{
		Response: &albv20251030.DescribeListenersResponseParams{
			Listeners: listeners,
		},
	}
}

func albXffAttrTypes(t *testing.T, testContext *albListenerTestContext) map[string]attr.Type {
	t.Helper()

	nested, ok := testContext.schema.Attributes["x_forwarded_for_config"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("attribute x_forwarded_for_config is not a single nested attribute")
	}
	objectType, ok := nested.GetType().(types.ObjectType)
	if !ok {
		t.Fatalf("attribute x_forwarded_for_config does not have an object type")
	}
	return objectType.AttrTypes
}

func albDefaultActionType(t *testing.T, testContext *albListenerTestContext) attr.Type {
	t.Helper()

	nested, ok := testContext.schema.Attributes["default_actions"].(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("attribute default_actions is not a list nested attribute")
	}
	listType, ok := nested.GetType().(types.ListType)
	if !ok {
		t.Fatalf("attribute default_actions does not have a list type")
	}
	return listType.ElemType
}

func albAssertNoErrorDiagnostics(t *testing.T, diags interface{ HasError() bool }) {
	t.Helper()

	if diags.HasError() {
		t.Fatalf("unexpected error diagnostics")
	}
}

func TestAlbListenerSchema(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbListenerTestContext(t, patches)

	for _, name := range []string{
		"id",
		"load_balancer_id",
		"listener_port",
		"listener_protocol",
		"listener_name",
		"default_actions",
		"certificate_ids",
		"ca_enabled",
		"ca_certificate_ids",
		"security_policy_id",
		"gzip_enabled",
		"http2enabled",
		"idle_timeout",
		"request_timeout",
		"x_forwarded_for_config",
		"tags",
		"listener_status",
		"create_time",
		"modify_time",
	} {
		if _, ok := testContext.schema.Attributes[name]; !ok {
			t.Fatalf("schema attribute %q is missing", name)
		}
	}
}

func TestAlbListenerCreate(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbListenerTestContext(t, patches)

	var capturedCreate *albv20251030.CreateListenerRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"CreateListenerWithContext",
		func(_ context.Context, request *albv20251030.CreateListenerRequest) (*albv20251030.CreateListenerResponse, error) {
			capturedCreate = request
			return &albv20251030.CreateListenerResponse{
				Response: &albv20251030.CreateListenerResponseParams{ListenerId: albString(albTestListenerId)},
			}, nil
		},
	)
	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeListenersWithContext",
		func(_ context.Context, request *albv20251030.DescribeListenersRequest) (*albv20251030.DescribeListenersResponse, error) {
			if request.LoadBalancerId == nil || *request.LoadBalancerId != albTestLoadBalancerId {
				t.Fatalf("DescribeListeners LoadBalancerId is not %q", albTestLoadBalancerId)
			}
			if len(request.ListenerIds) != 1 || request.ListenerIds[0] == nil || *request.ListenerIds[0] != albTestListenerId {
				t.Fatalf("DescribeListeners ListenerIds is not [%s]", albTestListenerId)
			}
			return albDescribeResponse(albSampleListener()), nil
		},
	)

	model := svcalb.AlbListenerResourceModel{
		LoadBalancerId:      types.StringValue(albTestLoadBalancerId),
		ListenerPort:        types.Int64Value(80),
		ListenerProtocol:    types.StringValue("HTTP"),
		ListenerName:        types.StringValue("tf-example"),
		DefaultActions:      types.ListNull(albDefaultActionType(t, testContext)),
		GzipEnabled:         types.BoolValue(true),
		IdleTimeout:         types.Int64Value(30),
		RequestTimeout:      types.Int64Value(60),
		Tags:                types.MapNull(types.StringType),
		CertificateIds:      types.ListNull(types.StringType),
		CaCertificateIds:    types.ListNull(types.StringType),
		XForwardedForConfig: types.ObjectNull(albXffAttrTypes(t, testContext)),
	}
	createResp := &resource.CreateResponse{State: tfsdk.State{Schema: testContext.schema}}
	testContext.resource.Create(context.Background(), resource.CreateRequest{Plan: testContext.plan(t, model)}, createResp)
	albAssertNoErrorDiagnostics(t, createResp.Diagnostics)

	if capturedCreate == nil {
		t.Fatalf("CreateListener was not called")
	}
	if capturedCreate.LoadBalancerId == nil || *capturedCreate.LoadBalancerId != albTestLoadBalancerId {
		t.Fatalf("CreateListener LoadBalancerId is not %q", albTestLoadBalancerId)
	}
	if capturedCreate.ListenerPort == nil || *capturedCreate.ListenerPort != 80 {
		t.Fatalf("CreateListener ListenerPort is not 80")
	}
	if capturedCreate.ListenerProtocol == nil || *capturedCreate.ListenerProtocol != "HTTP" {
		t.Fatalf("CreateListener ListenerProtocol is not HTTP")
	}
	if capturedCreate.ListenerName == nil || *capturedCreate.ListenerName != "tf-example" {
		t.Fatalf("CreateListener ListenerName is not tf-example")
	}
	if capturedCreate.GzipEnabled == nil || !*capturedCreate.GzipEnabled {
		t.Fatalf("CreateListener GzipEnabled is not true")
	}
	if capturedCreate.IdleTimeout == nil || *capturedCreate.IdleTimeout != 30 {
		t.Fatalf("CreateListener IdleTimeout is not 30")
	}
	if capturedCreate.ClientToken == nil || *capturedCreate.ClientToken == "" {
		t.Fatalf("CreateListener ClientToken is empty")
	}

	var state svcalb.AlbListenerResourceModel
	if diags := createResp.State.Get(context.Background(), &state); diags.HasError() {
		t.Fatalf("failed to read created state: %v", diags.Errors())
	}
	if state.Id.ValueString() != albTestLoadBalancerId+"#"+albTestListenerId {
		t.Fatalf("state id is %q, want %q", state.Id.ValueString(), albTestLoadBalancerId+"#"+albTestListenerId)
	}
	if state.ListenerStatus.ValueString() != "Active" {
		t.Fatalf("state listener_status is %q, want Active", state.ListenerStatus.ValueString())
	}
}

func TestAlbListenerReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbListenerTestContext(t, patches)

	patches.ApplyMethodFunc(
		testContext.client,
		"DescribeListenersWithContext",
		func(_ context.Context, _ *albv20251030.DescribeListenersRequest) (*albv20251030.DescribeListenersResponse, error) {
			return albDescribeResponse(), nil
		},
	)

	model := svcalb.AlbListenerResourceModel{
		Id:                  types.StringValue(albTestLoadBalancerId + "#" + albTestListenerId),
		LoadBalancerId:      types.StringValue(albTestLoadBalancerId),
		ListenerPort:        types.Int64Value(80),
		DefaultActions:      types.ListNull(albDefaultActionType(t, testContext)),
		CertificateIds:      types.ListNull(types.StringType),
		CaCertificateIds:    types.ListNull(types.StringType),
		Tags:                types.MapNull(types.StringType),
		XForwardedForConfig: types.ObjectNull(albXffAttrTypes(t, testContext)),
	}
	readResp := &resource.ReadResponse{State: testContext.state(t, model)}
	testContext.resource.Read(context.Background(), resource.ReadRequest{State: testContext.state(t, model)}, readResp)
	albAssertNoErrorDiagnostics(t, readResp.Diagnostics)

	if !readResp.State.Raw.IsNull() {
		t.Fatalf("state must be removed when the ALB listener does not exist")
	}
}

func TestAlbListenerDelete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	testContext := newAlbListenerTestContext(t, patches)

	var capturedDelete *albv20251030.DeleteListenerRequest
	patches.ApplyMethodFunc(
		testContext.client,
		"DeleteListenerWithContext",
		func(_ context.Context, request *albv20251030.DeleteListenerRequest) (*albv20251030.DeleteListenerResponse, error) {
			capturedDelete = request
			return &albv20251030.DeleteListenerResponse{
				Response: &albv20251030.DeleteListenerResponseParams{},
			}, nil
		},
	)

	model := svcalb.AlbListenerResourceModel{
		Id:                  types.StringValue(albTestLoadBalancerId + "#" + albTestListenerId),
		LoadBalancerId:      types.StringValue(albTestLoadBalancerId),
		ListenerPort:        types.Int64Value(80),
		DefaultActions:      types.ListNull(albDefaultActionType(t, testContext)),
		CertificateIds:      types.ListNull(types.StringType),
		CaCertificateIds:    types.ListNull(types.StringType),
		Tags:                types.MapNull(types.StringType),
		XForwardedForConfig: types.ObjectNull(albXffAttrTypes(t, testContext)),
	}
	deleteResp := &resource.DeleteResponse{State: testContext.state(t, model)}
	testContext.resource.Delete(context.Background(), resource.DeleteRequest{State: testContext.state(t, model)}, deleteResp)
	albAssertNoErrorDiagnostics(t, deleteResp.Diagnostics)

	if capturedDelete == nil {
		t.Fatalf("DeleteListener was not called")
	}
	if capturedDelete.LoadBalancerId == nil || *capturedDelete.LoadBalancerId != albTestLoadBalancerId {
		t.Fatalf("DeleteListener LoadBalancerId is not %q", albTestLoadBalancerId)
	}
	if len(capturedDelete.ListenerIds) != 1 || capturedDelete.ListenerIds[0] == nil || *capturedDelete.ListenerIds[0] != albTestListenerId {
		t.Fatalf("DeleteListener ListenerIds is not [%s]", albTestListenerId)
	}
}
