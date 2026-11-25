package alb_test

import (
	"context"
	"fmt"
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
	svalb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/alb"
)

const (
	albListenerTestId         = "lst-12345678"
	albListenerTestLoadBal    = "alb-12345678"
	albListenerTestName       = "tf-test-listener"
	albListenerTestRenamed    = "tf-test-listener-renamed"
	albListenerTestCloudName  = "tf-test-listener-normalized"
	albListenerTestPort       = int64(80)
	albListenerTestTargetGrp  = "lbtg-abcdefgh"
	albListenerTestTagKey     = "tag-key"
	albListenerTestTagValue   = "tag-value"
	albListenerTestCert       = "lsc-12345678"
	albListenerTestCaCert     = "lsc-abcdefgh"
	albListenerTestPolicyId   = "tls-12345678"
	albListenerTestToken      = "tf-test-token"
	albListenerTestIdleSecond = int64(15)
	albListenerTestReqSecond  = int64(60)
)

// Attribute type maps mirroring the nested schema definitions of the resource.
var (
	albListenerTargetGroupAttrTypes = map[string]attr.Type{
		"target_group_id": types.StringType,
		"weight":          types.Int64Type,
	}
	albListenerStickySessionAttrTypes = map[string]attr.Type{
		"enabled": types.BoolType,
		"timeout": types.Int64Type,
	}
	albListenerTargetGroupConfigAttrTypes = map[string]attr.Type{
		"target_groups":               types.ListType{ElemType: types.ObjectType{AttrTypes: albListenerTargetGroupAttrTypes}},
		"target_group_sticky_session": types.ObjectType{AttrTypes: albListenerStickySessionAttrTypes},
	}
	albListenerDefaultActionAttrTypes = map[string]attr.Type{
		"target_group_config": types.ObjectType{AttrTypes: albListenerTargetGroupConfigAttrTypes},
		"type":                types.StringType,
	}
	albListenerTagAttrTypes = map[string]attr.Type{
		"tag_key":   types.StringType,
		"tag_value": types.StringType,
	}
	albListenerXForwardedForAttrTypes = map[string]attr.Type{
		"x_forwarded_for_alb_id_enabled":          types.BoolType,
		"x_forwarded_for_client_src_port_enabled": types.BoolType,
		"x_forwarded_for_host_enabled":            types.BoolType,
		"x_forwarded_for_mode":                    types.StringType,
		"x_forwarded_for_port_enabled":            types.BoolType,
		"x_forwarded_for_proto_enabled":           types.BoolType,
		"x_tencent_client_idn_enabled":            types.BoolType,
		"x_tencent_client_sdn_enabled":            types.BoolType,
		"x_tencent_client_serial_enabled":         types.BoolType,
		"x_tencent_client_verify_enabled":         types.BoolType,
	}

	albListenerDefaultActionObjectType = types.ObjectType{AttrTypes: albListenerDefaultActionAttrTypes}
	albListenerTagObjectType           = types.ObjectType{AttrTypes: albListenerTagAttrTypes}
)

func albListenerStr(v string) *string { return &v }

func albListenerDefaultActions() types.List {
	return types.ListValueMust(albListenerDefaultActionObjectType, []attr.Value{
		types.ObjectValueMust(albListenerDefaultActionAttrTypes, map[string]attr.Value{
			"type": types.StringValue("ForwardGroup"),
			"target_group_config": types.ObjectValueMust(albListenerTargetGroupConfigAttrTypes, map[string]attr.Value{
				"target_groups": types.ListValueMust(types.ObjectType{AttrTypes: albListenerTargetGroupAttrTypes}, []attr.Value{
					types.ObjectValueMust(albListenerTargetGroupAttrTypes, map[string]attr.Value{
						"target_group_id": types.StringValue(albListenerTestTargetGrp),
						"weight":          types.Int64Value(10),
					}),
				}),
				"target_group_sticky_session": types.ObjectNull(albListenerStickySessionAttrTypes),
			}),
		}),
	})
}

// albListenerTestModel builds the terraform model used for plan/state. Every
// unset attribute carries an explicit typed null value.
func albListenerTestModel() svalb.AlbListenerResourceModel {
	return svalb.AlbListenerResourceModel{
		Id:               types.StringValue(albListenerTestLoadBal + "#" + albListenerTestId),
		ListenerId:       types.StringValue(albListenerTestId),
		LoadBalancerId:   types.StringValue(albListenerTestLoadBal),
		ListenerPort:     types.Int64Value(albListenerTestPort),
		ListenerProtocol: types.StringValue("HTTP"),
		ListenerName:     types.StringValue(albListenerTestName),
		DefaultActions:   albListenerDefaultActions(),
		CaCertificateIds: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue(albListenerTestCaCert),
		}),
		CaEnabled: types.BoolValue(false),
		CertificateIds: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue(albListenerTestCert),
		}),
		ClientToken:      types.StringValue(albListenerTestToken),
		GzipEnabled:      types.BoolValue(true),
		Http2Enabled:     types.BoolValue(false),
		IdleTimeout:      types.Int64Value(albListenerTestIdleSecond),
		RequestTimeout:   types.Int64Value(albListenerTestReqSecond),
		SecurityPolicyId: types.StringValue(albListenerTestPolicyId),
		Tags: types.ListValueMust(albListenerTagObjectType, []attr.Value{
			types.ObjectValueMust(albListenerTagAttrTypes, map[string]attr.Value{
				"tag_key":   types.StringValue(albListenerTestTagKey),
				"tag_value": types.StringValue(albListenerTestTagValue),
			}),
		}),
		XForwardedForConfig: types.ObjectValueMust(albListenerXForwardedForAttrTypes, map[string]attr.Value{
			"x_forwarded_for_alb_id_enabled":          types.BoolValue(true),
			"x_forwarded_for_client_src_port_enabled": types.BoolNull(),
			"x_forwarded_for_host_enabled":            types.BoolNull(),
			"x_forwarded_for_mode":                    types.StringValue("append"),
			"x_forwarded_for_port_enabled":            types.BoolNull(),
			"x_forwarded_for_proto_enabled":           types.BoolNull(),
			"x_tencent_client_idn_enabled":            types.BoolNull(),
			"x_tencent_client_sdn_enabled":            types.BoolNull(),
			"x_tencent_client_serial_enabled":         types.BoolNull(),
			"x_tencent_client_verify_enabled":         types.BoolNull(),
		}),
	}
}

// albListenerTestItem builds the DescribeListeners response item.
func albListenerTestItem() *albv20251030.ListenerOutput {
	return &albv20251030.ListenerOutput{
		ListenerId:          albListenerStr(albListenerTestId),
		ListenerName:        albListenerStr(albListenerTestName),
		ListenerPort:        func() *uint64 { v := uint64(albListenerTestPort); return &v }(),
		ListenerProtocol:    albListenerStr("HTTP"),
		CaEnable:            func() *bool { v := false; return &v }(),
		GzipEnabled:         func() *bool { v := true; return &v }(),
		Http2Enable:         func() *bool { v := false; return &v }(),
		IdleTimeout:         func() *uint64 { v := uint64(albListenerTestIdleSecond); return &v }(),
		RequestTimeout:      func() *uint64 { v := uint64(albListenerTestReqSecond); return &v }(),
		TlsSecurityPolicyId: albListenerStr(albListenerTestPolicyId),
		ListenerStatus:      albListenerStr("Active"),
		Tags: []*albv20251030.TagInfo{
			{
				TagKey:   albListenerStr(albListenerTestTagKey),
				TagValue: albListenerStr(albListenerTestTagValue),
			},
		},
		XForwardedForConfig: &albv20251030.XForwardedForConfig{
			XForwardedForAlbIdEnabled: func() *bool { v := true; return &v }(),
			XForwardedForMode:         albListenerStr("append"),
		},
	}
}

// albListenerTestContext wires the resource to a configured client; the SDK
// client methods are stubbed per test case.
type albListenerTestContext struct {
	resource resource.Resource
	schema   schema.Schema
	client   *connectivity.TencentCloudClient
}

func newAlbListenerTestContext(t *testing.T) *albListenerTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}

	res := svalb.NewAlbListenerResource()
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

	return &albListenerTestContext{resource: res, schema: schemaResp.Schema, client: connectivityClient}
}

func (c *albListenerTestContext) plan(t *testing.T, m svalb.AlbListenerResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}

	return plan
}

func (c *albListenerTestContext) state(t *testing.T, m svalb.AlbListenerResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}

	return state
}

// patchAlbListenerSdkClient stubs the versioned client accessor and returns the
// mock SDK client whose methods the test patches next.
func patchAlbListenerSdkClient(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) *albv20251030.Client {
	albClient := &albv20251030.Client{}
	patches.ApplyMethodReturn(client, "UseAlbV20251030Client", albClient)

	return albClient
}

// TestAlbListenerSchema verifies the registered type name.
func TestAlbListenerSchema(t *testing.T) {
	ctx := newAlbListenerTestContext(t)

	var metadataResp resource.MetadataResponse
	ctx.resource.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "tencentcloud"}, &metadataResp)
	if metadataResp.TypeName != "tencentcloud_alb_listener" {
		t.Errorf("type name = %q, want %q", metadataResp.TypeName, "tencentcloud_alb_listener")
	}
}

// TestAlbListenerCreateSuccess verifies that Create sends every configured
// field, records the returned listener id and refreshes the state.
func TestAlbListenerCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)

	var captured *albv20251030.CreateListenerRequest
	patches.ApplyMethodFunc(albClient, "CreateListener", func(req *albv20251030.CreateListenerRequest) (*albv20251030.CreateListenerResponse, error) {
		captured = req

		return &albv20251030.CreateListenerResponse{Response: &albv20251030.CreateListenerResponseParams{
			ListenerId: albListenerStr(albListenerTestId),
			RequestId:  albListenerStr("fake-request-id"),
		}}, nil
	})
	patches.ApplyMethodFunc(albClient, "DescribeListeners", func(req *albv20251030.DescribeListenersRequest) (*albv20251030.DescribeListenersResponse, error) {
		return &albv20251030.DescribeListenersResponse{Response: &albv20251030.DescribeListenersResponseParams{
			Listeners: []*albv20251030.ListenerOutput{albListenerTestItem()},
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})

	model := albListenerTestModel()

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("CreateListener request is nil")
	}
	if captured.LoadBalancerId == nil || *captured.LoadBalancerId != albListenerTestLoadBal {
		t.Errorf("create request load_balancer_id = %v, want %q", captured.LoadBalancerId, albListenerTestLoadBal)
	}
	if captured.ListenerPort == nil || *captured.ListenerPort != uint64(albListenerTestPort) {
		t.Errorf("create request listener_port = %v, want %d", captured.ListenerPort, albListenerTestPort)
	}
	if captured.ListenerProtocol == nil || *captured.ListenerProtocol != "HTTP" {
		t.Errorf("create request listener_protocol = %v, want %q", captured.ListenerProtocol, "HTTP")
	}
	if captured.ListenerName == nil || *captured.ListenerName != albListenerTestName {
		t.Errorf("create request listener_name = %v, want %q", captured.ListenerName, albListenerTestName)
	}
	if captured.IdleTimeout == nil || *captured.IdleTimeout != uint64(albListenerTestIdleSecond) {
		t.Errorf("create request idle_timeout = %v, want %d", captured.IdleTimeout, albListenerTestIdleSecond)
	}
	if captured.RequestTimeout == nil || *captured.RequestTimeout != uint64(albListenerTestReqSecond) {
		t.Errorf("create request request_timeout = %v, want %d", captured.RequestTimeout, albListenerTestReqSecond)
	}
	if captured.SecurityPolicyId == nil || *captured.SecurityPolicyId != albListenerTestPolicyId {
		t.Errorf("create request security_policy_id = %v, want %q", captured.SecurityPolicyId, albListenerTestPolicyId)
	}
	if len(captured.CertificateIds) != 1 || captured.CertificateIds[0] == nil || *captured.CertificateIds[0] != albListenerTestCert {
		t.Errorf("create request certificate_ids = %v, want [%q]", captured.CertificateIds, albListenerTestCert)
	}
	if len(captured.CaCertificateIds) != 1 || captured.CaCertificateIds[0] == nil || *captured.CaCertificateIds[0] != albListenerTestCaCert {
		t.Errorf("create request ca_certificate_ids = %v, want [%q]", captured.CaCertificateIds, albListenerTestCaCert)
	}
	if len(captured.Tags) != 1 || captured.Tags[0] == nil || captured.Tags[0].TagKey == nil || *captured.Tags[0].TagKey != albListenerTestTagKey {
		t.Errorf("create request tags = %v, want one tag with key %q", captured.Tags, albListenerTestTagKey)
	}
	if len(captured.DefaultActions) != 1 || captured.DefaultActions[0] == nil || captured.DefaultActions[0].Type == nil || *captured.DefaultActions[0].Type != "ForwardGroup" {
		t.Fatalf("create request default_actions = %v, want one ForwardGroup action", captured.DefaultActions)
	}
	if captured.DefaultActions[0].TargetGroupConfig == nil || len(captured.DefaultActions[0].TargetGroupConfig.TargetGroups) != 1 {
		t.Fatalf("create request target_group_config = %v, want one target group", captured.DefaultActions[0].TargetGroupConfig)
	}
	if id := captured.DefaultActions[0].TargetGroupConfig.TargetGroups[0].TargetGroupId; id == nil || *id != albListenerTestTargetGrp {
		t.Errorf("create request target_group_id = %v, want %q", id, albListenerTestTargetGrp)
	}
	if captured.XForwardedForConfig == nil || captured.XForwardedForConfig.XForwardedForAlbIdEnabled == nil || !*captured.XForwardedForConfig.XForwardedForAlbIdEnabled {
		t.Errorf("create request x_forwarded_for_config = %v, want alb id header enabled", captured.XForwardedForConfig)
	}

	var got svalb.AlbListenerResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != albListenerTestLoadBal+"#"+albListenerTestId {
		t.Errorf("id = %q, want %q", got.Id.ValueString(), albListenerTestLoadBal+"#"+albListenerTestId)
	}
	if got.ListenerId.ValueString() != albListenerTestId {
		t.Errorf("listener_id = %q, want %q", got.ListenerId.ValueString(), albListenerTestId)
	}
}

// TestAlbListenerCreateEmptyId verifies that Create fails loudly when the API
// returns no listener id.
func TestAlbListenerCreateEmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(albClient, "CreateListener", func(req *albv20251030.CreateListenerRequest) (*albv20251030.CreateListenerResponse, error) {
		return &albv20251030.CreateListenerResponse{Response: &albv20251030.CreateListenerResponseParams{
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, albListenerTestModel())
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, albListenerTestModel())}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail when the API returns an empty listener id")
	}
}

// TestAlbListenerCreateApiError verifies API failures surface as errors.
func TestAlbListenerCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(albClient, "CreateListener", func(req *albv20251030.CreateListenerRequest) (*albv20251030.CreateListenerResponse, error) {
		return nil, fmt.Errorf("create failed")
	})

	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, albListenerTestModel())
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, albListenerTestModel())}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail on an API error")
	}
}

// TestAlbListenerReadSuccess verifies Read refreshes the cloud-owned values.
// The describe stub reports a listener name DIFFERENT from the state, so a Read
// that skipped the refresh could not pass.
func TestAlbListenerReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)

	var captured *albv20251030.DescribeListenersRequest
	patches.ApplyMethodFunc(albClient, "DescribeListeners", func(req *albv20251030.DescribeListenersRequest) (*albv20251030.DescribeListenersResponse, error) {
		captured = req

		drifted := albListenerTestItem()
		drifted.ListenerName = albListenerStr(albListenerTestCloudName)

		return &albv20251030.DescribeListenersResponse{Response: &albv20251030.DescribeListenersResponseParams{
			Listeners: []*albv20251030.ListenerOutput{drifted},
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})

	model := albListenerTestModel()

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, model)
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, model)}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("DescribeListeners request is nil")
	}
	if captured.LoadBalancerId == nil || *captured.LoadBalancerId != albListenerTestLoadBal {
		t.Errorf("describe request load_balancer_id = %v, want %q", captured.LoadBalancerId, albListenerTestLoadBal)
	}
	if len(captured.ListenerIds) != 1 || captured.ListenerIds[0] == nil || *captured.ListenerIds[0] != albListenerTestId {
		t.Errorf("describe request listener_ids = %v, want [%q]", captured.ListenerIds, albListenerTestId)
	}

	var got svalb.AlbListenerResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.ListenerName.ValueString() != albListenerTestCloudName {
		t.Errorf("listener_name = %q, want the cloud value %q", got.ListenerName.ValueString(), albListenerTestCloudName)
	}
	if got.SecurityPolicyId.ValueString() != albListenerTestPolicyId {
		t.Errorf("security_policy_id = %q, want %q", got.SecurityPolicyId.ValueString(), albListenerTestPolicyId)
	}
}

// TestAlbListenerReadNotFound verifies that a describe list without the tracked
// listener drops the resource from state (confirmed not-found).
func TestAlbListenerReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(albClient, "DescribeListeners", func(req *albv20251030.DescribeListenersRequest) (*albv20251030.DescribeListenersResponse, error) {
		return &albv20251030.DescribeListenersResponse{Response: &albv20251030.DescribeListenersResponseParams{
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})

	model := albListenerTestModel()

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, model)
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, model)}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on confirmed not-found: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state should be removed after the listener is confirmed gone")
	}
}

// TestAlbListenerUpdateSuccess verifies that Update sends the planned values and
// records the post-update cloud state from the read-back, not the plan.
func TestAlbListenerUpdateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)

	var captured *albv20251030.ModifyListenerAttributesRequest
	patches.ApplyMethodFunc(albClient, "ModifyListenerAttributes", func(req *albv20251030.ModifyListenerAttributesRequest) (*albv20251030.ModifyListenerAttributesResponse, error) {
		captured = req

		return &albv20251030.ModifyListenerAttributesResponse{Response: &albv20251030.ModifyListenerAttributesResponseParams{
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})
	patches.ApplyMethodFunc(albClient, "DescribeListeners", func(req *albv20251030.DescribeListenersRequest) (*albv20251030.DescribeListenersResponse, error) {
		normalized := albListenerTestItem()
		normalized.ListenerName = albListenerStr(albListenerTestCloudName)

		return &albv20251030.DescribeListenersResponse{Response: &albv20251030.DescribeListenersResponseParams{
			Listeners: []*albv20251030.ListenerOutput{normalized},
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})

	plannedModel := albListenerTestModel()
	plannedModel.ListenerName = types.StringValue(albListenerTestRenamed)

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, albListenerTestModel())
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, plannedModel),
			State: ctx.state(t, albListenerTestModel()),
		}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ModifyListenerAttributes request is nil")
	}
	if captured.ListenerId == nil || *captured.ListenerId != albListenerTestId {
		t.Errorf("update request listener_id = %v, want %q", captured.ListenerId, albListenerTestId)
	}
	if captured.LoadBalancerId == nil || *captured.LoadBalancerId != albListenerTestLoadBal {
		t.Errorf("update request load_balancer_id = %v, want %q", captured.LoadBalancerId, albListenerTestLoadBal)
	}
	if captured.ListenerName == nil || *captured.ListenerName != albListenerTestRenamed {
		t.Errorf("update request listener_name = %v, want the planned %q", captured.ListenerName, albListenerTestRenamed)
	}
	if len(captured.DefaultActions) != 1 {
		t.Errorf("update request default_actions = %v, want one action", captured.DefaultActions)
	}

	var got svalb.AlbListenerResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.ListenerName.ValueString() != albListenerTestCloudName {
		t.Errorf("listener_name = %q, want the read-back cloud value %q", got.ListenerName.ValueString(), albListenerTestCloudName)
	}
	if got.Id.ValueString() != albListenerTestLoadBal+"#"+albListenerTestId {
		t.Errorf("id = %q, want the unchanged %q", got.Id.ValueString(), albListenerTestLoadBal+"#"+albListenerTestId)
	}
}

// TestAlbListenerUpdateApiError verifies that a modify failure surfaces as an
// error.
func TestAlbListenerUpdateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(albClient, "ModifyListenerAttributes", func(req *albv20251030.ModifyListenerAttributesRequest) (*albv20251030.ModifyListenerAttributesResponse, error) {
		return nil, fmt.Errorf("modify failed")
	})

	plannedModel := albListenerTestModel()
	plannedModel.ListenerName = types.StringValue(albListenerTestRenamed)

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, albListenerTestModel())
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{
			Plan:  ctx.plan(t, plannedModel),
			State: ctx.state(t, albListenerTestModel()),
		}, &updateResp)
	if !updateResp.Diagnostics.HasError() {
		t.Fatal("Update should fail on an API error")
	}
}

// TestAlbListenerDeleteSuccess verifies Delete sends the load balancer id and
// the tracked listener id.
func TestAlbListenerDeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)

	var captured *albv20251030.DeleteListenerRequest
	patches.ApplyMethodFunc(albClient, "DeleteListener", func(req *albv20251030.DeleteListenerRequest) (*albv20251030.DeleteListenerResponse, error) {
		captured = req

		return &albv20251030.DeleteListenerResponse{Response: &albv20251030.DeleteListenerResponseParams{
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})

	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, albListenerTestModel())}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("DeleteListener request is nil")
	}
	if captured.LoadBalancerId == nil || *captured.LoadBalancerId != albListenerTestLoadBal {
		t.Errorf("delete request load_balancer_id = %v, want %q", captured.LoadBalancerId, albListenerTestLoadBal)
	}
	if len(captured.ListenerIds) != 1 || captured.ListenerIds[0] == nil || *captured.ListenerIds[0] != albListenerTestId {
		t.Errorf("delete request listener_ids = %v, want [%q]", captured.ListenerIds, albListenerTestId)
	}
}

// TestAlbListenerDeleteApiError verifies that a delete failure surfaces as an
// error.
func TestAlbListenerDeleteApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(albClient, "DeleteListener", func(req *albv20251030.DeleteListenerRequest) (*albv20251030.DeleteListenerResponse, error) {
		return nil, fmt.Errorf("delete failed")
	})

	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, albListenerTestModel())}, &deleteResp)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatal("Delete should fail on an API error")
	}
}

// TestAlbListenerImportState verifies the composite id is split into the two
// locating attributes.
func TestAlbListenerImportState(t *testing.T) {
	ctx := newAlbListenerTestContext(t)

	importer, ok := ctx.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatal("resource does not implement resource.ResourceWithImportState")
	}

	importResp := resource.ImportStateResponse{State: ctx.state(t, svalb.AlbListenerResourceModel{
		Id:                  types.StringNull(),
		ListenerId:          types.StringNull(),
		LoadBalancerId:      types.StringNull(),
		ListenerPort:        types.Int64Null(),
		ListenerProtocol:    types.StringNull(),
		ListenerName:        types.StringNull(),
		DefaultActions:      types.ListNull(albListenerDefaultActionObjectType),
		CaCertificateIds:    types.ListNull(types.StringType),
		CaEnabled:           types.BoolNull(),
		CertificateIds:      types.ListNull(types.StringType),
		ClientToken:         types.StringNull(),
		GzipEnabled:         types.BoolNull(),
		Http2Enabled:        types.BoolNull(),
		IdleTimeout:         types.Int64Null(),
		RequestTimeout:      types.Int64Null(),
		SecurityPolicyId:    types.StringNull(),
		Tags:                types.ListNull(albListenerTagObjectType),
		XForwardedForConfig: types.ObjectNull(albListenerXForwardedForAttrTypes),
	})}

	importer.ImportState(context.Background(), resource.ImportStateRequest{
		ID: albListenerTestLoadBal + "#" + albListenerTestId,
	}, &importResp)
	if importResp.Diagnostics.HasError() {
		t.Fatalf("ImportState failed: %v", importResp.Diagnostics.Errors())
	}

	var got svalb.AlbListenerResourceModel
	if diags := importResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != albListenerTestLoadBal+"#"+albListenerTestId {
		t.Errorf("id = %q, want %q", got.Id.ValueString(), albListenerTestLoadBal+"#"+albListenerTestId)
	}
	if got.LoadBalancerId.ValueString() != albListenerTestLoadBal {
		t.Errorf("load_balancer_id = %q, want %q", got.LoadBalancerId.ValueString(), albListenerTestLoadBal)
	}
	if got.ListenerId.ValueString() != albListenerTestId {
		t.Errorf("listener_id = %q, want %q", got.ListenerId.ValueString(), albListenerTestId)
	}
}
