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
	albListenerTestId       = "rule-12345"
	albListenerTestListener = "lst-12345"
	albListenerTestLb       = "alb-12345"
	albListenerTestRuleName = "tf-test-rule"
)

func albListenerStr(v string) *string { return &v }
func albListenerInt(v int64) *int64   { return &v }

// albListenerTestItem builds a DescribeRules response item.
func albListenerTestItem() *albv20251030.RuleOutput {
	return &albv20251030.RuleOutput{
		RuleId:     albListenerStr(albListenerTestId),
		RuleName:   albListenerStr(albListenerTestRuleName),
		Priority:   albListenerInt(10),
		Direction:  albListenerStr("Request"),
		Status:     albListenerStr("Active"),
		CreateTime: albListenerStr("2024-01-01T00:00:00Z"),
		ModifyTime: albListenerStr("2024-01-02T00:00:00Z"),
		Actions: []*albv20251030.RuleAction{{
			Order: albListenerInt(1),
			Type:  albListenerStr("TargetGroup"),
		}},
		Conditions: []*albv20251030.RuleCondition{{
			Type:       albListenerStr("Host"),
			HostConfig: []*string{albListenerStr("www.example.com")},
		}},
	}
}

// albListenerTestModel builds a model with one rule; every nested field is set
// with an explicit typed value so tfsdk reflection succeeds.
func albListenerTestModel() svalb.AlbListenerResourceModel {
	action := types.ObjectValueMust(map[string]attr.Type{
		"order":                 types.Int64Type,
		"type":                  types.StringType,
		"fixed_response_config": types.ObjectType{AttrTypes: albFixedResponseTypes()},
		"insert_header_config":  types.ObjectType{AttrTypes: albInsertHeaderTypes()},
		"redirect_config":       types.ObjectType{AttrTypes: albRedirectTypes()},
		"remove_header_config":  types.ObjectType{AttrTypes: albRemoveHeaderTypes()},
		"rewrite_config":        types.ObjectType{AttrTypes: albRewriteTypes()},
		"target_group_config":   types.ObjectType{AttrTypes: albTargetGroupConfigTypes()},
	}, map[string]attr.Value{
		"order":                 types.Int64Value(1),
		"type":                  types.StringValue("TargetGroup"),
		"fixed_response_config": types.ObjectNull(albFixedResponseTypes()),
		"insert_header_config":  types.ObjectNull(albInsertHeaderTypes()),
		"redirect_config":       types.ObjectNull(albRedirectTypes()),
		"remove_header_config":  types.ObjectNull(albRemoveHeaderTypes()),
		"rewrite_config":        types.ObjectNull(albRewriteTypes()),
		"target_group_config":   types.ObjectNull(albTargetGroupConfigTypes()),
	})
	condition := types.ObjectValueMust(map[string]attr.Type{
		"type":                types.StringType,
		"cookie_config":       types.ListType{ElemType: types.ObjectType{AttrTypes: albCookieTypes()}},
		"header_config":       types.ObjectType{AttrTypes: albHeaderTypes()},
		"host_config":         types.ListType{ElemType: types.StringType},
		"method_config":       types.ListType{ElemType: types.StringType},
		"path_config":         types.ListType{ElemType: types.StringType},
		"query_string_config": types.ListType{ElemType: types.ObjectType{AttrTypes: albQueryStringTypes()}},
		"source_ip_config":    types.ListType{ElemType: types.StringType},
	}, map[string]attr.Value{
		"type":                types.StringValue("Host"),
		"cookie_config":       types.ListNull(types.ObjectType{AttrTypes: albCookieTypes()}),
		"header_config":       types.ObjectNull(albHeaderTypes()),
		"host_config":         types.ListValueMust(types.StringType, []attr.Value{types.StringValue("www.example.com")}),
		"method_config":       types.ListNull(types.StringType),
		"path_config":         types.ListNull(types.StringType),
		"query_string_config": types.ListNull(types.ObjectType{AttrTypes: albQueryStringTypes()}),
		"source_ip_config":    types.ListNull(types.StringType),
	})
	rule := types.ObjectValueMust(albRuleTypes(), map[string]attr.Value{
		"actions":     types.ListValueMust(types.ObjectType{AttrTypes: albActionTypes()}, []attr.Value{action}),
		"conditions":  types.ListValueMust(types.ObjectType{AttrTypes: albConditionTypes()}, []attr.Value{condition}),
		"priority":    types.Int64Value(10),
		"direction":   types.StringValue("Request"),
		"rule_name":   types.StringValue(albListenerTestRuleName),
		"tags":        types.ListNull(types.ObjectType{AttrTypes: albTagTypes()}),
		"rule_id":     types.StringNull(),
		"create_time": types.StringNull(),
		"modify_time": types.StringNull(),
		"status":      types.StringNull(),
	})

	return svalb.AlbListenerResourceModel{
		Id:             types.StringValue(albListenerTestId),
		ListenerId:     types.StringValue(albListenerTestListener),
		LoadBalancerId: types.StringValue(albListenerTestLb),
		Rules:          types.ListValueMust(types.ObjectType{AttrTypes: albRuleTypes()}, []attr.Value{rule}),
	}
}

func albRuleTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"actions":     types.ListType{ElemType: types.ObjectType{AttrTypes: albActionTypes()}},
		"conditions":  types.ListType{ElemType: types.ObjectType{AttrTypes: albConditionTypes()}},
		"priority":    types.Int64Type,
		"direction":   types.StringType,
		"rule_name":   types.StringType,
		"tags":        types.ListType{ElemType: types.ObjectType{AttrTypes: albTagTypes()}},
		"rule_id":     types.StringType,
		"create_time": types.StringType,
		"modify_time": types.StringType,
		"status":      types.StringType,
	}
}

func albActionTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"order":                 types.Int64Type,
		"type":                  types.StringType,
		"fixed_response_config": types.ObjectType{AttrTypes: albFixedResponseTypes()},
		"insert_header_config":  types.ObjectType{AttrTypes: albInsertHeaderTypes()},
		"redirect_config":       types.ObjectType{AttrTypes: albRedirectTypes()},
		"remove_header_config":  types.ObjectType{AttrTypes: albRemoveHeaderTypes()},
		"rewrite_config":        types.ObjectType{AttrTypes: albRewriteTypes()},
		"target_group_config":   types.ObjectType{AttrTypes: albTargetGroupConfigTypes()},
	}
}

func albConditionTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":                types.StringType,
		"cookie_config":       types.ListType{ElemType: types.ObjectType{AttrTypes: albCookieTypes()}},
		"header_config":       types.ObjectType{AttrTypes: albHeaderTypes()},
		"host_config":         types.ListType{ElemType: types.StringType},
		"method_config":       types.ListType{ElemType: types.StringType},
		"path_config":         types.ListType{ElemType: types.StringType},
		"query_string_config": types.ListType{ElemType: types.ObjectType{AttrTypes: albQueryStringTypes()}},
		"source_ip_config":    types.ListType{ElemType: types.StringType},
	}
}

func albFixedResponseTypes() map[string]attr.Type {
	return map[string]attr.Type{"http_code": types.Int64Type, "content": types.StringType, "content_type": types.StringType}
}

func albInsertHeaderTypes() map[string]attr.Type {
	return map[string]attr.Type{"key": types.StringType, "value": types.StringType, "value_type": types.StringType}
}

func albRedirectTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"http_code": types.Int64Type, "host": types.StringType, "path": types.StringType,
		"port": types.StringType, "protocol": types.StringType, "query": types.StringType,
	}
}

func albRemoveHeaderTypes() map[string]attr.Type {
	return map[string]attr.Type{"key": types.StringType}
}

func albRewriteTypes() map[string]attr.Type {
	return map[string]attr.Type{"host": types.StringType, "path": types.StringType, "query": types.StringType}
}

func albTargetGroupConfigTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"target_groups":               types.ListType{ElemType: types.ObjectType{AttrTypes: albTargetGroupTupleTypes()}},
		"target_group_sticky_session": types.ObjectType{AttrTypes: albStickySessionTypes()},
	}
}

func albTargetGroupTupleTypes() map[string]attr.Type {
	return map[string]attr.Type{"target_group_id": types.StringType, "weight": types.Int64Type}
}

func albStickySessionTypes() map[string]attr.Type {
	return map[string]attr.Type{"enabled": types.BoolType, "timeout": types.Int64Type}
}

func albCookieTypes() map[string]attr.Type {
	return map[string]attr.Type{"key": types.StringType, "value": types.StringType}
}

func albHeaderTypes() map[string]attr.Type {
	return map[string]attr.Type{"key": types.StringType, "values": types.ListType{ElemType: types.StringType}}
}

func albQueryStringTypes() map[string]attr.Type {
	return map[string]attr.Type{"key": types.StringType, "value": types.StringType}
}

func albTagTypes() map[string]attr.Type {
	return map[string]attr.Type{"tag_key": types.StringType, "tag_value": types.StringType}
}

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
	configured.Configure(context.Background(), resource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: connectivityClient},
	}, &configureResp)
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

func patchAlbListenerSdkClient(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) *albv20251030.Client {
	albClient := &albv20251030.Client{}
	patches.ApplyMethodReturn(client, "UseAlbV20251030Client", albClient)
	return albClient
}

// TestAlbListenerMetadata verifies the type name registered in the schema.
func TestAlbListenerMetadata(t *testing.T) {
	ctx := newAlbListenerTestContext(t)
	var resp resource.MetadataResponse
	ctx.resource.Metadata(context.Background(), resource.MetadataRequest{}, &resp)
	if resp.TypeName != "tencentcloud_alb_listener" {
		t.Errorf("type name = %q, want tencentcloud_alb_listener", resp.TypeName)
	}
}

// TestAlbListenerCreateSuccess verifies request mapping and read-back.
func TestAlbListenerCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)

	var captured *albv20251030.CreateRulesRequest
	patches.ApplyMethodFunc(albClient, "CreateRules", func(req *albv20251030.CreateRulesRequest) (*albv20251030.CreateRulesResponse, error) {
		captured = req
		return &albv20251030.CreateRulesResponse{Response: &albv20251030.CreateRulesResponseParams{
			RuleIds: []*string{albListenerStr(albListenerTestId)},
		}}, nil
	})
	patches.ApplyMethodFunc(albClient, "DescribeRules", func(req *albv20251030.DescribeRulesRequest) (*albv20251030.DescribeRulesResponse, error) {
		return &albv20251030.DescribeRulesResponse{Response: &albv20251030.DescribeRulesResponseParams{
			Rules:      []*albv20251030.RuleOutput{albListenerTestItem()},
			TotalCount: albListenerInt(1),
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
		t.Fatal("CreateRules request is nil")
	}
	if captured.ListenerId == nil || *captured.ListenerId != albListenerTestListener {
		t.Errorf("listener_id = %v, want %q", captured.ListenerId, albListenerTestListener)
	}
	if captured.LoadBalancerId == nil || *captured.LoadBalancerId != albListenerTestLb {
		t.Errorf("load_balancer_id = %v, want %q", captured.LoadBalancerId, albListenerTestLb)
	}
	if len(captured.Rules) != 1 || captured.Rules[0].Priority == nil || *captured.Rules[0].Priority != 10 {
		t.Fatalf("create rules = %v, want one rule with priority 10", captured.Rules)
	}

	var got svalb.AlbListenerResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != albListenerTestId {
		t.Errorf("id = %q, want %q", got.Id.ValueString(), albListenerTestId)
	}
}

// TestAlbListenerCreateEmptyId verifies Create fails on an empty id.
func TestAlbListenerCreateEmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(albClient, "CreateRules", func(req *albv20251030.CreateRulesRequest) (*albv20251030.CreateRulesResponse, error) {
		return &albv20251030.CreateRulesResponse{Response: &albv20251030.CreateRulesResponseParams{}}, nil
	})

	model := albListenerTestModel()
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail when CreateRules returns an empty rule id")
	}
}

// TestAlbListenerReadNotFound verifies that a missing rule removes state.
func TestAlbListenerReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(albClient, "DescribeRules", func(req *albv20251030.DescribeRulesRequest) (*albv20251030.DescribeRulesResponse, error) {
		return &albv20251030.DescribeRulesResponse{Response: &albv20251030.DescribeRulesResponseParams{}}, nil
	})

	model := albListenerTestModel()
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, model)
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, model)}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Fatal("Read should remove a rule that no longer exists")
	}
}

// TestAlbListenerUpdateSuccess verifies the modify request and read-back.
func TestAlbListenerUpdateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)

	var captured *albv20251030.ModifyRulesAttributesRequest
	patches.ApplyMethodFunc(albClient, "ModifyRulesAttributes", func(req *albv20251030.ModifyRulesAttributesRequest) (*albv20251030.ModifyRulesAttributesResponse, error) {
		captured = req
		return &albv20251030.ModifyRulesAttributesResponse{Response: &albv20251030.ModifyRulesAttributesResponseParams{}}, nil
	})
	patches.ApplyMethodFunc(albClient, "DescribeRules", func(req *albv20251030.DescribeRulesRequest) (*albv20251030.DescribeRulesResponse, error) {
		item := albListenerTestItem()
		item.RuleName = albListenerStr("tf-test-rule-renamed")
		return &albv20251030.DescribeRulesResponse{Response: &albv20251030.DescribeRulesResponseParams{
			Rules: []*albv20251030.RuleOutput{item},
		}}, nil
	})

	stateModel := albListenerTestModel()
	planModel := albListenerTestModel()
	var planRules []struct {
		Actions    types.List   `tfsdk:"actions"`
		Conditions types.List   `tfsdk:"conditions"`
		Priority   types.Int64  `tfsdk:"priority"`
		Direction  types.String `tfsdk:"direction"`
		RuleName   types.String `tfsdk:"rule_name"`
		Tags       types.List   `tfsdk:"tags"`
		RuleId     types.String `tfsdk:"rule_id"`
		CreateTime types.String `tfsdk:"create_time"`
		ModifyTime types.String `tfsdk:"modify_time"`
		Status     types.String `tfsdk:"status"`
	}
	// Change one user-configurable field so the update guard fires.
	if diags := planModel.Rules.ElementsAs(context.Background(), &planRules, false); diags.HasError() {
		t.Fatalf("failed to decode plan rules: %v", diags.Errors())
	}
	planRules[0].RuleName = types.StringValue("tf-test-rule-renamed")
	planModel.Rules = types.ListValueMust(types.ObjectType{AttrTypes: albRuleTypes()}, []attr.Value{types.ObjectValueMust(albRuleTypes(), map[string]attr.Value{
		"actions":     planRules[0].Actions,
		"conditions":  planRules[0].Conditions,
		"priority":    planRules[0].Priority,
		"direction":   planRules[0].Direction,
		"rule_name":   planRules[0].RuleName,
		"tags":        planRules[0].Tags,
		"rule_id":     planRules[0].RuleId,
		"create_time": planRules[0].CreateTime,
		"modify_time": planRules[0].ModifyTime,
		"status":      planRules[0].Status,
	})})

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, stateModel)
	ctx.resource.Update(context.Background(), resource.UpdateRequest{
		Plan:  ctx.plan(t, planModel),
		State: ctx.state(t, stateModel),
	}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ModifyRulesAttributes request is nil")
	}
	if len(captured.Rules) != 1 || captured.Rules[0].RuleId == nil || *captured.Rules[0].RuleId != albListenerTestId {
		t.Fatalf("modify rules = %v, want one rule with rule id %q", captured.Rules, albListenerTestId)
	}

	var got svalb.AlbListenerResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	var gotRules []types.Object
	if diags := got.Rules.ElementsAs(context.Background(), &gotRules, false); diags.HasError() {
		t.Fatalf("failed to decode state rules: %v", diags.Errors())
	}
	if len(gotRules) != 1 {
		t.Fatalf("state rules = %d, want 1", len(gotRules))
	}
	attrs := gotRules[0].Attributes()
	name, ok := attrs["rule_name"].(types.String)
	if !ok || name.ValueString() != "tf-test-rule-renamed" {
		t.Errorf("state rule_name = %v, want cloud value tf-test-rule-renamed", attrs["rule_name"])
	}
}

// TestAlbListenerDeleteSuccess verifies the delete request carries the rule id.
func TestAlbListenerDeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)

	var captured *albv20251030.DeleteRulesRequest
	patches.ApplyMethodFunc(albClient, "DeleteRules", func(req *albv20251030.DeleteRulesRequest) (*albv20251030.DeleteRulesResponse, error) {
		captured = req
		return &albv20251030.DeleteRulesResponse{Response: &albv20251030.DeleteRulesResponseParams{}}, nil
	})

	model := albListenerTestModel()
	var deleteResp resource.DeleteResponse
	deleteResp.State = ctx.state(t, model)
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, model)}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}

	if captured == nil || len(captured.RuleIds) != 1 || captured.RuleIds[0] == nil || *captured.RuleIds[0] != albListenerTestId {
		t.Fatalf("delete request rule ids = %v, want [%q]", captured, albListenerTestId)
	}
}

// TestAlbListenerReadApiError verifies SDK errors surface as diagnostics.
func TestAlbListenerReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbListenerTestContext(t)
	albClient := patchAlbListenerSdkClient(patches, ctx.client)
	patches.ApplyMethodFunc(albClient, "DescribeRules", func(req *albv20251030.DescribeRulesRequest) (*albv20251030.DescribeRulesResponse, error) {
		return nil, fmt.Errorf("network down")
	})

	model := albListenerTestModel()
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, model)
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, model)}, &readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("Read should surface API errors")
	}
}
