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
	albRulesTestId       = "rule-12345"
	albRulesListenerId   = "lst-12345"
	albRulesLbId         = "alb-12345"
	albRulesFalseRequest = "fake-request-id"
)

func albRulesStr(v string) *string { return &v }
func albRulesInt(v int64) *int64   { return &v }

// albRulesTestRule builds a minimal RuleInput used by Create stubs.
func albRulesTestRuleInput() *albv20251030.RuleInput {
	return &albv20251030.RuleInput{
		Actions: []*albv20251030.RuleAction{{
			Order: albRulesInt(1),
			Type:  albRulesStr("TargetGroup"),
			TargetGroupConfig: &albv20251030.TargetGroupConfig{
				TargetGroups: []*albv20251030.TargetGroupTuple{{
					TargetGroupId: albRulesStr("lbtg-12345"),
					Weight:        albRulesInt(10),
				}},
			},
		}},
		Conditions: []*albv20251030.RuleCondition{{
			Type:       albRulesStr("Host"),
			HostConfig: []*string{albRulesStr("example.com")},
		}},
		Priority: albRulesInt(1),
		RuleName: albRulesStr("tf-test-rule"),
	}
}

// albRulesTestRuleObject builds the types.Object value for one rule using
// the same AttrTypes map declared by the resource implementation.
func albRulesTestRuleObject() types.Object {
	actions := types.ListValueMust(types.ObjectType{AttrTypes: svalb.RulesActionsAttrTypes}, []attr.Value{
		types.ObjectValueMust(svalb.RulesActionsAttrTypes, map[string]attr.Value{
			"order":                 types.Int64Value(1),
			"type":                  types.StringValue("TargetGroup"),
			"fixed_response_config": types.ObjectNull(svalb.RulesActionsFixedResponseConfigAttrTypes),
			"insert_header_config":  types.ObjectNull(svalb.RulesActionsInsertHeaderConfigAttrTypes),
			"redirect_config":       types.ObjectNull(svalb.RulesActionsRedirectConfigAttrTypes),
			"remove_header_config":  types.ObjectNull(svalb.RulesActionsRemoveHeaderConfigAttrTypes),
			"rewrite_config":        types.ObjectNull(svalb.RulesActionsRewriteConfigAttrTypes),
			"target_group_config": types.ObjectValueMust(svalb.RulesActionsTargetGroupConfigAttrTypes, map[string]attr.Value{
				"target_groups": types.ListValueMust(types.ObjectType{AttrTypes: svalb.RulesActionsTargetGroupConfigTargetGroupsAttrTypes}, []attr.Value{
					types.ObjectValueMust(svalb.RulesActionsTargetGroupConfigTargetGroupsAttrTypes, map[string]attr.Value{
						"target_group_id": types.StringValue("lbtg-12345"),
						"weight":          types.Int64Value(10),
					}),
				}),
				"target_group_sticky_session": types.ObjectNull(svalb.RulesActionsTargetGroupConfigTargetGroupStickySessionAttrTypes),
			}),
		}),
	})
	conditions := types.ListValueMust(types.ObjectType{AttrTypes: svalb.RulesConditionsAttrTypes}, []attr.Value{
		types.ObjectValueMust(svalb.RulesConditionsAttrTypes, map[string]attr.Value{
			"type":                types.StringValue("Host"),
			"cookie_config":       types.ListNull(types.ObjectType{AttrTypes: svalb.RulesConditionsCookieConfigAttrTypes}),
			"header_config":       types.ObjectNull(svalb.RulesConditionsHeaderConfigAttrTypes),
			"host_config":         types.ListValueMust(types.StringType, []attr.Value{types.StringValue("example.com")}),
			"method_config":       types.ListNull(types.StringType),
			"path_config":         types.ListNull(types.StringType),
			"query_string_config": types.ListNull(types.ObjectType{AttrTypes: svalb.RulesConditionsQueryStringConfigAttrTypes}),
			"source_ip_config":    types.ListNull(types.StringType),
		}),
	})
	return types.ObjectValueMust(svalb.RulesAttrTypes, map[string]attr.Value{
		"actions":    actions,
		"conditions": conditions,
		"priority":   types.Int64Value(1),
		"direction":  types.StringNull(),
		"rule_name":  types.StringValue("tf-test-rule"),
		"tags":       types.ListNull(types.ObjectType{AttrTypes: svalb.RulesTagsAttrTypes}),
	})
}

func albRulesTestModel() svalb.AlbRulesResourceModel {
	return svalb.AlbRulesResourceModel{
		Id:             types.StringValue(albRulesTestId),
		ListenerId:     types.StringValue(albRulesListenerId),
		LoadBalancerId: types.StringValue(albRulesLbId),
		Rules: types.ListValueMust(types.ObjectType{AttrTypes: svalb.RulesAttrTypes}, []attr.Value{
			albRulesTestRuleObject(),
		}),
		ClientToken: types.StringNull(),
		DryRun:      types.BoolNull(),
	}
}

func albRulesTestOutput() *albv20251030.RuleOutput {
	in := albRulesTestRuleInput()
	return &albv20251030.RuleOutput{
		Actions:    in.Actions,
		Conditions: in.Conditions,
		Priority:   in.Priority,
		Direction:  in.Direction,
		RuleId:     albRulesStr(albRulesTestId),
		RuleName:   in.RuleName,
		Tags:       in.Tags,
	}
}

type albRulesTestContext struct {
	resource resource.Resource
	schema   schema.Schema
	client   *connectivity.TencentCloudClient
}

func newAlbRulesTestContext(t *testing.T) *albRulesTestContext {
	t.Helper()
	connectivityClient := &connectivity.TencentCloudClient{}
	res := svalb.NewAlbRulesResource()
	configured, ok := res.(resource.ResourceWithConfigure)
	if !ok {
		t.Fatalf("resource does not implement ResourceWithConfigure")
	}
	var configureResp resource.ConfigureResponse
	configured.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: connectivityClient}}, &configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("configure failed: %v", configureResp.Diagnostics.Errors())
	}
	var schemaResp resource.SchemaResponse
	res.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("schema failed: %v", schemaResp.Diagnostics.Errors())
	}
	return &albRulesTestContext{resource: res, schema: schemaResp.Schema, client: connectivityClient}
}

func (c *albRulesTestContext) plan(t *testing.T, m svalb.AlbRulesResourceModel) tfsdk.Plan {
	t.Helper()
	p := tfsdk.Plan{Schema: c.schema}
	if d := p.Set(context.Background(), m); d.HasError() {
		t.Fatalf("plan: %v", d.Errors())
	}
	return p
}

func (c *albRulesTestContext) state(t *testing.T, m svalb.AlbRulesResourceModel) tfsdk.State {
	t.Helper()
	s := tfsdk.State{Schema: c.schema}
	if d := s.Set(context.Background(), m); d.HasError() {
		t.Fatalf("state: %v", d.Errors())
	}
	return s
}

func patchAlbRulesClient(patches *gomonkey.Patches, client *connectivity.TencentCloudClient) *albv20251030.Client {
	sdk := &albv20251030.Client{}
	patches.ApplyMethodReturn(client, "UseAlbV20251030Client", sdk)
	return sdk
}

func TestAlbRulesCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()
	ctx := newAlbRulesTestContext(t)
	sdk := patchAlbRulesClient(patches, ctx.client)
	var captured *albv20251030.CreateRulesRequest
	patches.ApplyMethodFunc(sdk, "CreateRules", func(req *albv20251030.CreateRulesRequest) (*albv20251030.CreateRulesResponse, error) {
		captured = req
		return &albv20251030.CreateRulesResponse{Response: &albv20251030.CreateRulesResponseParams{
			RuleIds: []*string{albRulesStr(albRulesTestId)}, RequestId: albRulesStr(albRulesFalseRequest),
		}}, nil
	})
	patches.ApplyMethodFunc(sdk, "DescribeRules", func(req *albv20251030.DescribeRulesRequest) (*albv20251030.DescribeRulesResponse, error) {
		return &albv20251030.DescribeRulesResponse{Response: &albv20251030.DescribeRulesResponseParams{
			Rules: []*albv20251030.RuleOutput{albRulesTestOutput()}, RequestId: albRulesStr(albRulesFalseRequest),
		}}, nil
	})
	model := albRulesTestModel()
	model.Id = types.StringNull()
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}
	if captured == nil || captured.ListenerId == nil || *captured.ListenerId != albRulesListenerId {
		t.Errorf("create request listener id not passed: %v", captured)
	}
	if captured != nil && (len(captured.Rules) != 1 || captured.Rules[0].Priority == nil || *captured.Rules[0].Priority != 1) {
		t.Errorf("create request rules not passed: %v", captured.Rules)
	}
	var got svalb.AlbRulesResourceModel
	if d := createResp.State.Get(context.Background(), &got); d.HasError() {
		t.Fatalf("read state: %v", d.Errors())
	}
	if got.Id.ValueString() != albRulesTestId {
		t.Errorf("id = %q, want %q", got.Id.ValueString(), albRulesTestId)
	}
}

func TestAlbRulesCreateEmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()
	ctx := newAlbRulesTestContext(t)
	sdk := patchAlbRulesClient(patches, ctx.client)
	patches.ApplyMethodFunc(sdk, "CreateRules", func(req *albv20251030.CreateRulesRequest) (*albv20251030.CreateRulesResponse, error) {
		return &albv20251030.CreateRulesResponse{Response: &albv20251030.CreateRulesResponseParams{RequestId: albRulesStr(albRulesFalseRequest)}}, nil
	})
	model := albRulesTestModel()
	model.Id = types.StringNull()
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail with an empty rule id")
	}
}

func TestAlbRulesCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()
	ctx := newAlbRulesTestContext(t)
	sdk := patchAlbRulesClient(patches, ctx.client)
	patches.ApplyMethodFunc(sdk, "CreateRules", func(req *albv20251030.CreateRulesRequest) (*albv20251030.CreateRulesResponse, error) {
		return nil, fmt.Errorf("create failed")
	})
	model := albRulesTestModel()
	model.Id = types.StringNull()
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail on an API error")
	}
}

func TestAlbRulesReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()
	ctx := newAlbRulesTestContext(t)
	sdk := patchAlbRulesClient(patches, ctx.client)
	var gotId string
	patches.ApplyMethodFunc(sdk, "DescribeRules", func(req *albv20251030.DescribeRulesRequest) (*albv20251030.DescribeRulesResponse, error) {
		if len(req.RuleIds) == 1 && req.RuleIds[0] != nil {
			gotId = *req.RuleIds[0]
		}
		return &albv20251030.DescribeRulesResponse{Response: &albv20251030.DescribeRulesResponseParams{
			Rules: []*albv20251030.RuleOutput{albRulesTestOutput()}, RequestId: albRulesStr(albRulesFalseRequest),
		}}, nil
	})
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, albRulesTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, albRulesTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}
	if gotId != albRulesTestId {
		t.Errorf("describe called with %q, want %q", gotId, albRulesTestId)
	}
	var got svalb.AlbRulesResourceModel
	if d := readResp.State.Get(context.Background(), &got); d.HasError() {
		t.Fatalf("read state: %v", d.Errors())
	}
	if got.Id.ValueString() != albRulesTestId {
		t.Errorf("id = %q", got.Id.ValueString())
	}
}

func TestAlbRulesReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()
	ctx := newAlbRulesTestContext(t)
	sdk := patchAlbRulesClient(patches, ctx.client)
	patches.ApplyMethodFunc(sdk, "DescribeRules", func(req *albv20251030.DescribeRulesRequest) (*albv20251030.DescribeRulesResponse, error) {
		return &albv20251030.DescribeRulesResponse{Response: &albv20251030.DescribeRulesResponseParams{RequestId: albRulesStr(albRulesFalseRequest)}}, nil
	})
	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, albRulesTestModel())
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, albRulesTestModel())}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on not-found: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state should be removed on not-found")
	}
}

func TestAlbRulesUpdateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()
	ctx := newAlbRulesTestContext(t)
	sdk := patchAlbRulesClient(patches, ctx.client)
	var captured *albv20251030.ModifyRulesAttributesRequest
	patches.ApplyMethodFunc(sdk, "ModifyRulesAttributes", func(req *albv20251030.ModifyRulesAttributesRequest) (*albv20251030.ModifyRulesAttributesResponse, error) {
		captured = req
		return &albv20251030.ModifyRulesAttributesResponse{Response: &albv20251030.ModifyRulesAttributesResponseParams{RequestId: albRulesStr(albRulesFalseRequest)}}, nil
	})
	patches.ApplyMethodFunc(sdk, "DescribeRules", func(req *albv20251030.DescribeRulesRequest) (*albv20251030.DescribeRulesResponse, error) {
		return &albv20251030.DescribeRulesResponse{Response: &albv20251030.DescribeRulesResponseParams{
			Rules: []*albv20251030.RuleOutput{albRulesTestOutput()}, RequestId: albRulesStr(albRulesFalseRequest),
		}}, nil
	})
	planned := albRulesTestModel()
	// Change a covered attribute so the write API is invoked.
	plannedRules := albRulesTestRuleObject().Attributes()
	plannedRules["rule_name"] = types.StringValue("tf-test-rule-renamed")
	planned.Rules = types.ListValueMust(types.ObjectType{AttrTypes: svalb.RulesAttrTypes}, []attr.Value{
		types.ObjectValueMust(svalb.RulesAttrTypes, plannedRules),
	})
	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, albRulesTestModel())
	ctx.resource.Update(context.Background(), resource.UpdateRequest{Plan: ctx.plan(t, planned), State: ctx.state(t, albRulesTestModel())}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}
	if captured == nil || len(captured.Rules) != 1 || captured.Rules[0].RuleId == nil || *captured.Rules[0].RuleId != albRulesTestId {
		t.Errorf("modify request did not carry the rule id: %v", captured)
	}
	var got svalb.AlbRulesResourceModel
	if d := updateResp.State.Get(context.Background(), &got); d.HasError() {
		t.Fatalf("read state: %v", d.Errors())
	}
	if got.Id.ValueString() != albRulesTestId {
		t.Errorf("id changed after update: %q", got.Id.ValueString())
	}
}

func TestAlbRulesDeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()
	ctx := newAlbRulesTestContext(t)
	sdk := patchAlbRulesClient(patches, ctx.client)
	var captured *albv20251030.DeleteRulesRequest
	patches.ApplyMethodFunc(sdk, "DeleteRules", func(req *albv20251030.DeleteRulesRequest) (*albv20251030.DeleteRulesResponse, error) {
		captured = req
		return &albv20251030.DeleteRulesResponse{Response: &albv20251030.DeleteRulesResponseParams{RequestId: albRulesStr(albRulesFalseRequest)}}, nil
	})
	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, albRulesTestModel())}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}
	if captured == nil || len(captured.RuleIds) != 1 || captured.RuleIds[0] == nil || *captured.RuleIds[0] != albRulesTestId {
		t.Errorf("delete request did not carry the rule id: %v", captured)
	}
}
