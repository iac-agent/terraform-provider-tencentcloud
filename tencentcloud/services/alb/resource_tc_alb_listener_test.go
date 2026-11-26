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
	svsalb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/alb"
)

const (
	albListenerTestRuleId     = "rule-12345"
	albListenerTestListenerId = "lst-12345"
	albListenerTestLbId       = "alb-12345"
	albListenerTestGroupName  = "lbtg-12345"
	albListenerTestRuleName   = "tf-test-rule"
	albListenerTestNormalized = "tf-test-rule-normalized"
)

// Local mirrors of the schema's nested attribute types (derived from the
// generated schema), used to build typed tfsdk values for the tests.
func albListenerTestTagAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{"tag_key": types.StringType, "tag_value": types.StringType}
}

func albListenerTestConditionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":                types.StringType,
		"cookie_config":       types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"key": types.StringType, "value": types.StringType}}},
		"header_config":       types.ObjectType{AttrTypes: map[string]attr.Type{"key": types.StringType, "values": types.ListType{ElemType: types.StringType}}},
		"host_config":         types.ListType{ElemType: types.StringType},
		"method_config":       types.ListType{ElemType: types.StringType},
		"path_config":         types.ListType{ElemType: types.StringType},
		"query_string_config": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"key": types.StringType, "value": types.StringType}}},
		"source_ip_config":    types.ListType{ElemType: types.StringType},
	}
}

func albListenerTestStickySessionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{"enabled": types.BoolType, "timeout": types.Int64Type}
}

func albListenerTestTargetGroupTupleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{"target_group_id": types.StringType, "weight": types.Int64Type}
}

func albListenerTestTargetGroupConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"target_groups":               types.ListType{ElemType: types.ObjectType{AttrTypes: albListenerTestTargetGroupTupleAttrTypes()}},
		"target_group_sticky_session": types.ObjectType{AttrTypes: albListenerTestStickySessionAttrTypes()},
	}
}

func albListenerTestActionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"order":                 types.Int64Type,
		"type":                  types.StringType,
		"fixed_response_config": types.ObjectType{AttrTypes: map[string]attr.Type{"http_code": types.Int64Type, "content": types.StringType, "content_type": types.StringType}},
		"insert_header_config":  types.ObjectType{AttrTypes: map[string]attr.Type{"key": types.StringType, "value": types.StringType, "value_type": types.StringType}},
		"redirect_config":       types.ObjectType{AttrTypes: map[string]attr.Type{"http_code": types.Int64Type, "host": types.StringType, "path": types.StringType, "port": types.StringType, "protocol": types.StringType, "query": types.StringType}},
		"remove_header_config":  types.ObjectType{AttrTypes: map[string]attr.Type{"key": types.StringType}},
		"rewrite_config":        types.ObjectType{AttrTypes: map[string]attr.Type{"host": types.StringType, "path": types.StringType, "query": types.StringType}},
		"target_group_config":   types.ObjectType{AttrTypes: albListenerTestTargetGroupConfigAttrTypes()},
	}
}

func albListenerTestRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"actions":     types.ListType{ElemType: types.ObjectType{AttrTypes: albListenerTestActionAttrTypes()}},
		"conditions":  types.ListType{ElemType: types.ObjectType{AttrTypes: albListenerTestConditionAttrTypes()}},
		"priority":    types.Int64Type,
		"direction":   types.StringType,
		"rule_name":   types.StringType,
		"tags":        types.ListType{ElemType: types.ObjectType{AttrTypes: albListenerTestTagAttrTypes()}},
		"rule_id":     types.StringType,
		"status":      types.StringType,
		"create_time": types.StringType,
		"modify_time": types.StringType,
	}
}

// albListenerTestRule mirrors the rule model for typed ElementsAs reads.
type albListenerTestRule struct {
	Actions    types.List   `tfsdk:"actions"`
	Conditions types.List   `tfsdk:"conditions"`
	Priority   types.Int64  `tfsdk:"priority"`
	Direction  types.String `tfsdk:"direction"`
	RuleName   types.String `tfsdk:"rule_name"`
	Tags       types.List   `tfsdk:"tags"`
	RuleId     types.String `tfsdk:"rule_id"`
	Status     types.String `tfsdk:"status"`
	CreateTime types.String `tfsdk:"create_time"`
	ModifyTime types.String `tfsdk:"modify_time"`
}

func albListenerStr(v string) *string { return &v }

func albListenerInt(v int64) *int64 { return &v }

func albListenerTestAction(t *testing.T) attr.Value {
	t.Helper()
	targetGroup := types.ObjectValueMust(
		albListenerTestTargetGroupConfigAttrTypes(),
		map[string]attr.Value{
			"target_groups": types.ListValueMust(
				types.ObjectType{AttrTypes: albListenerTestTargetGroupTupleAttrTypes()},
				[]attr.Value{types.ObjectValueMust(
					albListenerTestTargetGroupTupleAttrTypes(),
					map[string]attr.Value{
						"target_group_id": types.StringValue(albListenerTestGroupName),
						"weight":          types.Int64Value(10),
					},
				)},
			),
			"target_group_sticky_session": types.ObjectNull(albListenerTestStickySessionAttrTypes()),
		},
	)
	return types.ObjectValueMust(
		albListenerTestActionAttrTypes(),
		map[string]attr.Value{
			"order":                 types.Int64Value(1),
			"type":                  types.StringValue("TargetGroup"),
			"fixed_response_config": types.ObjectNull(map[string]attr.Type{"http_code": types.Int64Type, "content": types.StringType, "content_type": types.StringType}),
			"insert_header_config":  types.ObjectNull(map[string]attr.Type{"key": types.StringType, "value": types.StringType, "value_type": types.StringType}),
			"redirect_config":       types.ObjectNull(map[string]attr.Type{"http_code": types.Int64Type, "host": types.StringType, "path": types.StringType, "port": types.StringType, "protocol": types.StringType, "query": types.StringType}),
			"remove_header_config":  types.ObjectNull(map[string]attr.Type{"key": types.StringType}),
			"rewrite_config":        types.ObjectNull(map[string]attr.Type{"host": types.StringType, "path": types.StringType, "query": types.StringType}),
			"target_group_config":   targetGroup,
		},
	)
}

func albListenerTestRules(t *testing.T, ruleName string, withComputed bool) types.List {
	t.Helper()
	rule := types.ObjectValueMust(
		albListenerTestRuleAttrTypes(),
		map[string]attr.Value{
			"actions":     types.ListValueMust(types.ObjectType{AttrTypes: albListenerTestActionAttrTypes()}, []attr.Value{albListenerTestAction(t)}),
			"conditions":  types.ListNull(types.ObjectType{AttrTypes: albListenerTestConditionAttrTypes()}),
			"priority":    types.Int64Value(10),
			"direction":   types.StringValue("Request"),
			"rule_name":   types.StringValue(ruleName),
			"tags":        types.ListNull(types.ObjectType{AttrTypes: albListenerTestTagAttrTypes()}),
			"rule_id":     albListenerTestString(withComputed, albListenerTestRuleId),
			"status":      albListenerTestString(withComputed, "Active"),
			"create_time": albListenerTestString(withComputed, "2025-01-01T00:00:00Z"),
			"modify_time": albListenerTestString(withComputed, "2025-01-02T00:00:00Z"),
		},
	)
	return types.ListValueMust(types.ObjectType{AttrTypes: albListenerTestRuleAttrTypes()}, []attr.Value{rule})
}

func albListenerTestString(set bool, v string) types.String {
	if set {
		return types.StringValue(v)
	}
	return types.StringNull()
}

func albListenerTestModel(t *testing.T, ruleName string, withComputed bool) svsalb.AlbListenerResourceModel {
	t.Helper()
	return svsalb.AlbListenerResourceModel{
		Id:             types.StringValue(albListenerTestRuleId),
		ListenerId:     types.StringValue(albListenerTestListenerId),
		LoadBalancerId: types.StringValue(albListenerTestLbId),
		ClientToken:    types.StringNull(),
		DryRun:         types.BoolNull(),
		Rules:          albListenerTestRules(t, ruleName, withComputed),
	}
}

func albListenerTestItem(ruleName string) *albv20251030.RuleOutput {
	return &albv20251030.RuleOutput{
		RuleId:     albListenerStr(albListenerTestRuleId),
		RuleName:   albListenerStr(ruleName),
		Priority:   albListenerInt(10),
		Direction:  albListenerStr("Request"),
		Status:     albListenerStr("Active"),
		CreateTime: albListenerStr("2025-01-01T00:00:00Z"),
		ModifyTime: albListenerStr("2025-01-02T00:00:00Z"),
		Actions: []*albv20251030.RuleAction{{
			Order: albListenerInt(1),
			Type:  albListenerStr("TargetGroup"),
			TargetGroupConfig: &albv20251030.TargetGroupConfig{
				TargetGroups: []*albv20251030.TargetGroupTuple{{
					TargetGroupId: albListenerStr(albListenerTestGroupName),
					Weight:        albListenerInt(10),
				}},
			},
		}},
	}
}

func patchAlbListenerSdkClient(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) *albv20251030.Client {
	albClient := &albv20251030.Client{}
	patches.ApplyMethodReturn(client, "UseAlbV20251030Client", albClient)
	return albClient
}

type albListenerTestContext struct {
	resource resource.Resource
	schema   schema.Schema
	client   *connectivity.TencentCloudClient
}

func newAlbListenerTestContext(t *testing.T) *albListenerTestContext {
	t.Helper()

	connectivityClient := &connectivity.TencentCloudClient{}
	res := svsalb.NewAlbListenerResource()
	configured, ok := res.(resource.ResourceWithConfigure)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithConfigure")
	}
	var configureResp resource.ConfigureResponse
	configured.Configure(context.Background(),
		resource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: connectivityClient}},
		&configureResp)
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

func (c *albListenerTestContext) plan(t *testing.T, m svsalb.AlbListenerResourceModel) tfsdk.Plan {
	t.Helper()
	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}
	return plan
}

func (c *albListenerTestContext) state(t *testing.T, m svsalb.AlbListenerResourceModel) tfsdk.State {
	t.Helper()
	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), m); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}
	return state
}

func assertAlbListenerRules(t *testing.T, list types.List, wantRuleName, wantRuleId, wantStatus string) {
	t.Helper()
	var rules []albListenerTestRule
	if diags := list.ElementsAs(context.Background(), &rules, false); diags.HasError() {
		t.Fatalf("rules: %v", diags.Errors())
	}
	if len(rules) != 1 {
		t.Fatalf("rules length = %d, want 1", len(rules))
	}
	if rules[0].RuleName.ValueString() != wantRuleName {
		t.Errorf("rule_name = %q, want %q", rules[0].RuleName.ValueString(), wantRuleName)
	}
	if rules[0].RuleId.ValueString() != wantRuleId {
		t.Errorf("rule_id = %q, want %q", rules[0].RuleId.ValueString(), wantRuleId)
	}
	if rules[0].Status.ValueString() != wantStatus {
		t.Errorf("status = %q, want %q", rules[0].Status.ValueString(), wantStatus)
	}
}

func TestAlbListenerMetadata(t *testing.T) {
	ctx := newAlbListenerTestContext(t)
	var resp resource.MetadataResponse
	ctx.resource.Metadata(context.Background(), resource.MetadataRequest{}, &resp)
	if resp.TypeName != "tencentcloud_alb_listener" {
		t.Errorf("type name = %q, want tencentcloud_alb_listener", resp.TypeName)
	}
}

func TestAlbListenerCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)

	var captured *albv20251030.CreateRulesRequest
	patches.ApplyMethodFunc(albClient, "CreateRules", func(req *albv20251030.CreateRulesRequest) (*albv20251030.CreateRulesResponse, error) {
		captured = req
		return &albv20251030.CreateRulesResponse{Response: &albv20251030.CreateRulesResponseParams{
			RuleIds:   []*string{albListenerStr(albListenerTestRuleId)},
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})
	patches.ApplyMethodFunc(albClient, "DescribeRules", func(req *albv20251030.DescribeRulesRequest) (*albv20251030.DescribeRulesResponse, error) {
		return &albv20251030.DescribeRulesResponse{Response: &albv20251030.DescribeRulesResponseParams{
			Rules:     []*albv20251030.RuleOutput{albListenerTestItem(albListenerTestRuleName)},
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})

	model := albListenerTestModel(t, albListenerTestRuleName, false)
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("CreateRules request is nil")
	}
	if captured.ListenerId == nil || *captured.ListenerId != albListenerTestListenerId {
		t.Errorf("create request listener_id = %v, want %q", captured.ListenerId, albListenerTestListenerId)
	}
	if captured.LoadBalancerId == nil || *captured.LoadBalancerId != albListenerTestLbId {
		t.Errorf("create request load_balancer_id = %v, want %q", captured.LoadBalancerId, albListenerTestLbId)
	}
	if len(captured.Rules) != 1 || len(captured.Rules[0].Actions) != 1 || captured.Rules[0].Actions[0].TargetGroupConfig == nil {
		t.Fatal("create request rules were not expanded")
	}

	var got svsalb.AlbListenerResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != albListenerTestRuleId {
		t.Errorf("id = %q, want %q", got.Id.ValueString(), albListenerTestRuleId)
	}
	assertAlbListenerRules(t, got.Rules, albListenerTestRuleName, albListenerTestRuleId, "Active")
}

func TestAlbListenerCreateEmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(albClient, "CreateRules", func(req *albv20251030.CreateRulesRequest) (*albv20251030.CreateRulesResponse, error) {
		return &albv20251030.CreateRulesResponse{Response: &albv20251030.CreateRulesResponseParams{
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})

	model := albListenerTestModel(t, albListenerTestRuleName, false)
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail when CreateRules returns no rule id")
	}
}

func TestAlbListenerCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(albClient, "CreateRules", func(req *albv20251030.CreateRulesRequest) (*albv20251030.CreateRulesResponse, error) {
		return nil, fmt.Errorf("create failed")
	})

	model := albListenerTestModel(t, albListenerTestRuleName, false)
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail on an API error")
	}
}

func TestAlbListenerReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)
	var gotId string
	patches.ApplyMethodFunc(albClient, "DescribeRules", func(req *albv20251030.DescribeRulesRequest) (*albv20251030.DescribeRulesResponse, error) {
		if len(req.RuleIds) == 1 && req.RuleIds[0] != nil {
			gotId = *req.RuleIds[0]
		}
		return &albv20251030.DescribeRulesResponse{Response: &albv20251030.DescribeRulesResponseParams{
			Rules:     []*albv20251030.RuleOutput{albListenerTestItem(albListenerTestRuleName)},
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})

	model := albListenerTestModel(t, albListenerTestRuleName, true)
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, model)
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, model)}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}
	if gotId != albListenerTestRuleId {
		t.Errorf("DescribeRules called with %q, want %q", gotId, albListenerTestRuleId)
	}
	var got svsalb.AlbListenerResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertAlbListenerRules(t, got.Rules, albListenerTestRuleName, albListenerTestRuleId, "Active")
}

func TestAlbListenerReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(albClient, "DescribeRules", func(req *albv20251030.DescribeRulesRequest) (*albv20251030.DescribeRulesResponse, error) {
		return &albv20251030.DescribeRulesResponse{Response: &albv20251030.DescribeRulesResponseParams{
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})

	model := albListenerTestModel(t, albListenerTestRuleName, true)
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, model)
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, model)}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on confirmed not-found: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state should be removed after the rule is confirmed gone")
	}
}

func TestAlbListenerUpdateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)

	var captured *albv20251030.ModifyRulesAttributesRequest
	patches.ApplyMethodFunc(albClient, "ModifyRulesAttributes", func(req *albv20251030.ModifyRulesAttributesRequest) (*albv20251030.ModifyRulesAttributesResponse, error) {
		captured = req
		return &albv20251030.ModifyRulesAttributesResponse{Response: &albv20251030.ModifyRulesAttributesResponseParams{
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})
	patches.ApplyMethodFunc(albClient, "DescribeRules", func(req *albv20251030.DescribeRulesRequest) (*albv20251030.DescribeRulesResponse, error) {
		return &albv20251030.DescribeRulesResponse{Response: &albv20251030.DescribeRulesResponseParams{
			Rules:     []*albv20251030.RuleOutput{albListenerTestItem(albListenerTestNormalized)},
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})

	planned := albListenerTestModel(t, "tf-test-rule-renamed", false)
	state := albListenerTestModel(t, albListenerTestRuleName, true)
	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, state)
	ctx.resource.Update(context.Background(),
		resource.UpdateRequest{Plan: ctx.plan(t, planned), State: ctx.state(t, state)}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}
	if captured == nil {
		t.Fatal("ModifyRulesAttributes request is nil")
	}
	if len(captured.Rules) != 1 || captured.Rules[0].RuleId == nil || *captured.Rules[0].RuleId != albListenerTestRuleId {
		t.Errorf("update request did not carry the tracked rule id: %+v", captured.Rules)
	}

	var got svsalb.AlbListenerResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	assertAlbListenerRules(t, got.Rules, albListenerTestNormalized, albListenerTestRuleId, "Active")
}

func TestAlbListenerDeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)

	var captured *albv20251030.DeleteRulesRequest
	patches.ApplyMethodFunc(albClient, "DeleteRules", func(req *albv20251030.DeleteRulesRequest) (*albv20251030.DeleteRulesResponse, error) {
		captured = req
		return &albv20251030.DeleteRulesResponse{Response: &albv20251030.DeleteRulesResponseParams{
			RequestId: albListenerStr("fake-request-id"),
		}}, nil
	})

	model := albListenerTestModel(t, albListenerTestRuleName, true)
	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, model)}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}
	if captured == nil {
		t.Fatal("DeleteRules request is nil")
	}
	if len(captured.RuleIds) != 1 || captured.RuleIds[0] == nil || *captured.RuleIds[0] != albListenerTestRuleId {
		t.Errorf("delete request rule_ids = %v, want [%q]", captured.RuleIds, albListenerTestRuleId)
	}
}
