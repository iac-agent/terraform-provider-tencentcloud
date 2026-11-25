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
	svcalb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/alb"
)

const (
	albRulesTestLoadBalancerId  = "alb-12345678"
	albRulesTestListenerId      = "lst-12345678"
	albRulesTestRuleId          = "rule-12345678"
	albRulesTestRuleName        = "tf-example-rule"
	albRulesTestRenamedRuleName = "tf-example-rule-renamed"
	albRulesTestCloudRuleName   = "tf-example-rule-normalized"
	albRulesTestTargetGroupId   = "lbtg-12345678"
	albRulesTestHost            = "www.example.com"
	albRulesTestPlannedPriority = int64(10)
	albRulesTestTagKey          = "env"
	albRulesTestTagValue        = "test"
	albRulesTestCompositeId     = "alb-12345678#lst-12345678#rule-12345678"
	albRulesTestActionType      = "TargetGroup"
	albRulesTestConditionType   = "Host"
	albRulesTestStickyTimeout   = int64(1000)
)

func albRulesTestStr(v string) *string { return &v }

func albRulesTestInt64(v int64) *int64 { return &v }

func albRulesTestBool(v bool) *bool { return &v }

// albRulesTestContext wires the resource to a configured client; the SDK
// client methods are stubbed per test case.
type albRulesTestContext struct {
	resource resource.Resource
	schema   schema.Schema
	client   *connectivity.TencentCloudClient
	sdk      *albv20251030.Client
}

func newAlbRulesTestContext(t *testing.T) *albRulesTestContext {
	t.Helper()

	patches := gomonkey.NewPatches()
	sdkClient := &albv20251030.Client{}
	connectivityClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(connectivityClient, "UseAlbV20251030Client", sdkClient)
	t.Cleanup(patches.Reset)

	res := svcalb.NewAlbRulesResource()
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

	return &albRulesTestContext{resource: res, schema: schemaResp.Schema, client: connectivityClient, sdk: sdkClient}
}

func (c *albRulesTestContext) plan(t *testing.T, model svcalb.AlbRulesResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: c.schema}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags.Errors())
	}

	return plan
}

func (c *albRulesTestContext) state(t *testing.T, model svcalb.AlbRulesResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: c.schema}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags.Errors())
	}

	return state
}

// objectAttrTypes returns the attribute types of a top-level single nested
// attribute so the test values mirror the schema exactly.
func (c *albRulesTestContext) objectAttrTypes(t *testing.T, name string) map[string]attr.Type {
	t.Helper()

	nested, ok := c.schema.Attributes[name].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("attribute %q is not a single nested attribute", name)
	}
	objectType, ok := nested.GetType().(types.ObjectType)
	if !ok {
		t.Fatalf("attribute %q does not have an object type", name)
	}

	return objectType.AttrTypes
}

// listAttrTypes returns the attribute types of a top-level list nested
// attribute.
func (c *albRulesTestContext) listAttrTypes(t *testing.T, name string) map[string]attr.Type {
	t.Helper()

	nested, ok := c.schema.Attributes[name].(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("attribute %q is not a list nested attribute", name)
	}
	listType, ok := nested.GetType().(types.ListType)
	if !ok {
		t.Fatalf("attribute %q does not have a list type", name)
	}
	objectType, ok := listType.ElemType.(types.ObjectType)
	if !ok {
		t.Fatalf("attribute %q does not have an object element type", name)
	}

	return objectType.AttrTypes
}

// listElementAttrTypes returns the attribute types of a nested list below the
// given parent attribute types.
func albRulesTestListElementAttrTypes(t *testing.T, parent map[string]attr.Type, name string) map[string]attr.Type {
	t.Helper()

	listType, ok := parent[name].(types.ListType)
	if !ok {
		t.Fatalf("attribute %q does not have a list type", name)
	}
	objectType, ok := listType.ElemType.(types.ObjectType)
	if !ok {
		t.Fatalf("attribute %q does not have an object element type", name)
	}

	return objectType.AttrTypes
}

// objectAttrTypesOf returns the attribute types of a nested single object
// below the given parent attribute types.
func albRulesTestObjectAttrTypes(t *testing.T, parent map[string]attr.Type, name string) map[string]attr.Type {
	t.Helper()

	objectType, ok := parent[name].(types.ObjectType)
	if !ok {
		t.Fatalf("attribute %q does not have an object type", name)
	}

	return objectType.AttrTypes
}

func albRulesTestObject(t *testing.T, attributeTypes map[string]attr.Type, values map[string]attr.Value) types.Object {
	t.Helper()

	object, diags := types.ObjectValue(attributeTypes, values)
	if diags.HasError() {
		t.Fatalf("failed to build object value: %v", diags.Errors())
	}

	return object
}

func albRulesTestList(t *testing.T, elementTypes map[string]attr.Type, elements ...map[string]attr.Value) types.List {
	t.Helper()

	objectType := types.ObjectType{AttrTypes: elementTypes}
	values := make([]attr.Value, 0, len(elements))
	for _, element := range elements {
		object, diags := types.ObjectValue(elementTypes, element)
		if diags.HasError() {
			t.Fatalf("failed to build list element: %v", diags.Errors())
		}
		values = append(values, object)
	}

	list, diags := types.ListValue(objectType, values)
	if diags.HasError() {
		t.Fatalf("failed to build list value: %v", diags.Errors())
	}

	return list
}

func albRulesTestStringList(t *testing.T, values ...string) types.List {
	t.Helper()

	elements := make([]attr.Value, 0, len(values))
	for _, value := range values {
		elements = append(elements, types.StringValue(value))
	}

	list, diags := types.ListValue(types.StringType, elements)
	if diags.HasError() {
		t.Fatalf("failed to build string list: %v", diags.Errors())
	}

	return list
}

// sampleModel builds the model used for plan/state: every schema attribute gets
// an explicit value or an explicit typed null.
func (c *albRulesTestContext) sampleModel(t *testing.T) svcalb.AlbRulesResourceModel {
	t.Helper()

	actionTypes := c.listAttrTypes(t, "actions")
	conditionTypes := c.listAttrTypes(t, "conditions")
	tagTypes := c.listAttrTypes(t, "tags")
	targetGroupConfigTypes := albRulesTestObjectAttrTypes(t, actionTypes, "target_group_config")
	targetGroupTupleTypes := albRulesTestListElementAttrTypes(t, targetGroupConfigTypes, "target_groups")
	stickySessionTypes := albRulesTestObjectAttrTypes(t, targetGroupConfigTypes, "target_group_sticky_session")

	return svcalb.AlbRulesResourceModel{
		Id:             types.StringNull(),
		LoadBalancerId: types.StringValue(albRulesTestLoadBalancerId),
		ListenerId:     types.StringValue(albRulesTestListenerId),
		Priority:       types.Int64Value(albRulesTestPlannedPriority),
		Direction:      types.StringValue("Request"),
		RuleName:       types.StringValue(albRulesTestRuleName),
		Actions: albRulesTestList(t, actionTypes, map[string]attr.Value{
			"order":                 types.Int64Value(1),
			"type":                  types.StringValue(albRulesTestActionType),
			"fixed_response_config": types.ObjectNull(albRulesTestObjectAttrTypes(t, actionTypes, "fixed_response_config")),
			"insert_header_config":  types.ObjectNull(albRulesTestObjectAttrTypes(t, actionTypes, "insert_header_config")),
			"redirect_config":       types.ObjectNull(albRulesTestObjectAttrTypes(t, actionTypes, "redirect_config")),
			"remove_header_config":  types.ObjectNull(albRulesTestObjectAttrTypes(t, actionTypes, "remove_header_config")),
			"rewrite_config":        types.ObjectNull(albRulesTestObjectAttrTypes(t, actionTypes, "rewrite_config")),
			"target_group_config": albRulesTestObject(t, targetGroupConfigTypes, map[string]attr.Value{
				"target_groups": albRulesTestList(t, targetGroupTupleTypes, map[string]attr.Value{
					"target_group_id": types.StringValue(albRulesTestTargetGroupId),
					"weight":          types.Int64Value(10),
				}),
				"target_group_sticky_session": albRulesTestObject(t, stickySessionTypes, map[string]attr.Value{
					"enabled": types.BoolValue(true),
					"timeout": types.Int64Value(albRulesTestStickyTimeout),
				}),
			}),
		}),
		Conditions: albRulesTestList(t, conditionTypes, map[string]attr.Value{
			"type":                types.StringValue(albRulesTestConditionType),
			"cookie_config":       types.ListNull(types.ObjectType{AttrTypes: albRulesTestListElementAttrTypes(t, conditionTypes, "cookie_config")}),
			"header_config":       types.ObjectNull(albRulesTestObjectAttrTypes(t, conditionTypes, "header_config")),
			"host_config":         albRulesTestStringList(t, albRulesTestHost),
			"method_config":       types.ListNull(types.StringType),
			"path_config":         types.ListNull(types.StringType),
			"query_string_config": types.ListNull(types.ObjectType{AttrTypes: albRulesTestListElementAttrTypes(t, conditionTypes, "query_string_config")}),
			"source_ip_config":    types.ListNull(types.StringType),
		}),
		Tags: albRulesTestList(t, tagTypes, map[string]attr.Value{
			"tag_key":   types.StringValue(albRulesTestTagKey),
			"tag_value": types.StringValue(albRulesTestTagValue),
		}),
	}
}

// albRulesTestItem builds the DescribeRules response item. It deliberately
// reports values that differ from the plan where a case asserts the read-back.
func albRulesTestItem(t *testing.T, model svcalb.AlbRulesResourceModel) *albv20251030.RuleOutput {
	t.Helper()

	item := &albv20251030.RuleOutput{
		RuleId:     albRulesTestStr(albRulesTestRuleId),
		Priority:   albRulesTestInt64(albRulesTestPlannedPriority),
		Direction:  albRulesTestStr("Request"),
		RuleName:   albRulesTestStr(albRulesTestRuleName),
		CreateTime: albRulesTestStr("2026-01-01 00:00:00"),
		ModifyTime: albRulesTestStr("2026-01-01 00:00:00"),
		Status:     albRulesTestStr("Active"),
		Actions: []*albv20251030.RuleAction{
			{
				Order: albRulesTestInt64(1),
				Type:  albRulesTestStr(albRulesTestActionType),
				TargetGroupConfig: &albv20251030.TargetGroupConfig{
					TargetGroups: []*albv20251030.TargetGroupTuple{
						{
							TargetGroupId: albRulesTestStr(albRulesTestTargetGroupId),
							Weight:        albRulesTestInt64(10),
						},
					},
					TargetGroupStickySession: &albv20251030.TargetGroupStickySession{
						Enabled: albRulesTestBool(true),
						Timeout: albRulesTestInt64(albRulesTestStickyTimeout),
					},
				},
			},
		},
		Conditions: []*albv20251030.RuleCondition{
			{
				Type:       albRulesTestStr(albRulesTestConditionType),
				HostConfig: []*string{albRulesTestStr(albRulesTestHost)},
			},
		},
		Tags: []*albv20251030.TagInfo{
			{
				TagKey:   albRulesTestStr(albRulesTestTagKey),
				TagValue: albRulesTestStr(albRulesTestTagValue),
			},
		},
	}

	return item
}

// TestAlbRulesCreateSuccess verifies that Create sends the location fields and
// the full rule payload, records the composite id and refreshes the state from
// the read-back.
func TestAlbRulesCreateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbRulesTestContext(t)

	var captured *albv20251030.CreateRulesRequest
	patches.ApplyMethodFunc(ctx.sdk, "CreateRules", func(req *albv20251030.CreateRulesRequest) (*albv20251030.CreateRulesResponse, error) {
		captured = req

		return &albv20251030.CreateRulesResponse{Response: &albv20251030.CreateRulesResponseParams{
			RuleIds:   []*string{albRulesTestStr(albRulesTestRuleId)},
			RequestId: albRulesTestStr("fake-request-id"),
		}}, nil
	})
	patches.ApplyMethodFunc(ctx.sdk, "DescribeRules", func(req *albv20251030.DescribeRulesRequest) (*albv20251030.DescribeRulesResponse, error) {
		return &albv20251030.DescribeRulesResponse{Response: &albv20251030.DescribeRulesResponseParams{
			Rules:     []*albv20251030.RuleOutput{albRulesTestItem(t, ctx.sampleModel(t))},
			RequestId: albRulesTestStr("fake-request-id"),
		}}, nil
	})

	model := ctx.sampleModel(t)
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create failed: %v", createResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("CreateRules request is nil")
	}
	if captured.LoadBalancerId == nil || *captured.LoadBalancerId != albRulesTestLoadBalancerId {
		t.Errorf("create request load_balancer_id = %v, want %q", captured.LoadBalancerId, albRulesTestLoadBalancerId)
	}
	if captured.ListenerId == nil || *captured.ListenerId != albRulesTestListenerId {
		t.Errorf("create request listener_id = %v, want %q", captured.ListenerId, albRulesTestListenerId)
	}
	if len(captured.Rules) != 1 || captured.Rules[0] == nil {
		t.Fatalf("create request rules = %v, want exactly one rule", captured.Rules)
	}

	rule := captured.Rules[0]
	if rule.Priority == nil || *rule.Priority != albRulesTestPlannedPriority {
		t.Errorf("create request priority = %v, want %d", rule.Priority, albRulesTestPlannedPriority)
	}
	if rule.RuleName == nil || *rule.RuleName != albRulesTestRuleName {
		t.Errorf("create request rule_name = %v, want %q", rule.RuleName, albRulesTestRuleName)
	}
	if rule.Direction == nil || *rule.Direction != "Request" {
		t.Errorf("create request direction = %v, want Request", rule.Direction)
	}
	if len(rule.Actions) != 1 || rule.Actions[0].TargetGroupConfig == nil ||
		len(rule.Actions[0].TargetGroupConfig.TargetGroups) != 1 ||
		rule.Actions[0].TargetGroupConfig.TargetGroups[0].TargetGroupId == nil ||
		*rule.Actions[0].TargetGroupConfig.TargetGroups[0].TargetGroupId != albRulesTestTargetGroupId {
		t.Errorf("create request actions = %v, want the target group %q", rule.Actions, albRulesTestTargetGroupId)
	}
	if len(rule.Conditions) != 1 || rule.Conditions[0].HostConfig == nil || len(rule.Conditions[0].HostConfig) != 1 ||
		rule.Conditions[0].HostConfig[0] == nil || *rule.Conditions[0].HostConfig[0] != albRulesTestHost {
		t.Errorf("create request conditions = %v, want the host %q", rule.Conditions, albRulesTestHost)
	}
	if len(rule.Tags) != 1 || rule.Tags[0].TagKey == nil || *rule.Tags[0].TagKey != albRulesTestTagKey {
		t.Errorf("create request tags = %v, want tag key %q", rule.Tags, albRulesTestTagKey)
	}

	var got svcalb.AlbRulesResourceModel
	if diags := createResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != albRulesTestCompositeId {
		t.Errorf("id = %q, want %q", got.Id.ValueString(), albRulesTestCompositeId)
	}
	if got.Priority.ValueInt64() != albRulesTestPlannedPriority {
		t.Errorf("priority = %d, want %d", got.Priority.ValueInt64(), albRulesTestPlannedPriority)
	}
	if got.RuleName.ValueString() != albRulesTestRuleName {
		t.Errorf("rule_name = %q, want %q", got.RuleName.ValueString(), albRulesTestRuleName)
	}
}

// TestAlbRulesCreateEmptyId verifies that Create fails loudly when the API
// returns no rule id.
func TestAlbRulesCreateEmptyId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbRulesTestContext(t)
	patches.ApplyMethodFunc(ctx.sdk, "CreateRules", func(req *albv20251030.CreateRulesRequest) (*albv20251030.CreateRulesResponse, error) {
		return &albv20251030.CreateRulesResponse{Response: &albv20251030.CreateRulesResponseParams{
			RequestId: albRulesTestStr("fake-request-id"),
		}}, nil
	})

	model := ctx.sampleModel(t)
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail when the API returns an empty rule id")
	}
}

// TestAlbRulesCreateApiError verifies that API failures surface as errors.
func TestAlbRulesCreateApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbRulesTestContext(t)
	patches.ApplyMethodFunc(ctx.sdk, "CreateRules", func(req *albv20251030.CreateRulesRequest) (*albv20251030.CreateRulesResponse, error) {
		return nil, fmt.Errorf("create failed")
	})

	model := ctx.sampleModel(t)
	var createResp resource.CreateResponse
	createResp.State = ctx.state(t, model)
	ctx.resource.Create(context.Background(), resource.CreateRequest{Plan: ctx.plan(t, model)}, &createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should fail on an API error")
	}
}

// TestAlbRulesReadSuccess verifies that Read addresses the rule by the three id
// segments and refreshes the model from the cloud.
func TestAlbRulesReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbRulesTestContext(t)

	var captured *albv20251030.DescribeRulesRequest
	patches.ApplyMethodFunc(ctx.sdk, "DescribeRules", func(req *albv20251030.DescribeRulesRequest) (*albv20251030.DescribeRulesResponse, error) {
		captured = req

		return &albv20251030.DescribeRulesResponse{Response: &albv20251030.DescribeRulesResponseParams{
			Rules:     []*albv20251030.RuleOutput{albRulesTestItem(t, ctx.sampleModel(t))},
			RequestId: albRulesTestStr("fake-request-id"),
		}}, nil
	})

	model := ctx.sampleModel(t)
	model.Id = types.StringValue(albRulesTestCompositeId)

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, model)
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, model)}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read failed: %v", readResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("DescribeRules request is nil")
	}
	if captured.LoadBalancerId == nil || *captured.LoadBalancerId != albRulesTestLoadBalancerId {
		t.Errorf("describe request load_balancer_id = %v, want %q", captured.LoadBalancerId, albRulesTestLoadBalancerId)
	}
	if captured.ListenerId == nil || *captured.ListenerId != albRulesTestListenerId {
		t.Errorf("describe request listener_id = %v, want %q", captured.ListenerId, albRulesTestListenerId)
	}
	if len(captured.RuleIds) != 1 || captured.RuleIds[0] == nil || *captured.RuleIds[0] != albRulesTestRuleId {
		t.Errorf("describe request rule_ids = %v, want [%q]", captured.RuleIds, albRulesTestRuleId)
	}

	var got svcalb.AlbRulesResourceModel
	if diags := readResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != albRulesTestCompositeId {
		t.Errorf("id = %q, want %q", got.Id.ValueString(), albRulesTestCompositeId)
	}
	if got.Priority.ValueInt64() != albRulesTestPlannedPriority {
		t.Errorf("priority = %d, want %d", got.Priority.ValueInt64(), albRulesTestPlannedPriority)
	}
}

// TestAlbRulesReadNotFound verifies that a describe answer without the tracked
// rule drops the resource from state.
func TestAlbRulesReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbRulesTestContext(t)
	patches.ApplyMethodFunc(ctx.sdk, "DescribeRules", func(req *albv20251030.DescribeRulesRequest) (*albv20251030.DescribeRulesResponse, error) {
		return &albv20251030.DescribeRulesResponse{Response: &albv20251030.DescribeRulesResponseParams{
			RequestId: albRulesTestStr("fake-request-id"),
		}}, nil
	})

	model := ctx.sampleModel(t)
	model.Id = types.StringValue(albRulesTestCompositeId)

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, model)
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, model)}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read should not error on a confirmed not-found: %v", readResp.Diagnostics.Errors())
	}
	if !readResp.State.Raw.IsNull() {
		t.Error("state should be removed after the rule is confirmed gone")
	}
}

// TestAlbRulesReadApiError verifies that a failed describe surfaces as an error
// and keeps the state.
func TestAlbRulesReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbRulesTestContext(t)
	patches.ApplyMethodFunc(ctx.sdk, "DescribeRules", func(req *albv20251030.DescribeRulesRequest) (*albv20251030.DescribeRulesResponse, error) {
		return nil, fmt.Errorf("describe failed")
	})

	model := ctx.sampleModel(t)
	model.Id = types.StringValue(albRulesTestCompositeId)

	var readResp resource.ReadResponse
	readResp.State = ctx.state(t, model)
	ctx.resource.Read(context.Background(), resource.ReadRequest{State: ctx.state(t, model)}, &readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("Read should fail on an API error")
	}
}

// TestAlbRulesUpdateSuccess verifies that Update sends the planned values to
// ModifyRulesAttributes and records the POST-update cloud state.
func TestAlbRulesUpdateSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbRulesTestContext(t)

	var captured *albv20251030.ModifyRulesAttributesRequest
	patches.ApplyMethodFunc(ctx.sdk, "ModifyRulesAttributes", func(req *albv20251030.ModifyRulesAttributesRequest) (*albv20251030.ModifyRulesAttributesResponse, error) {
		captured = req

		return &albv20251030.ModifyRulesAttributesResponse{Response: &albv20251030.ModifyRulesAttributesResponseParams{
			RequestId: albRulesTestStr("fake-request-id"),
		}}, nil
	})
	patches.ApplyMethodFunc(ctx.sdk, "DescribeRules", func(req *albv20251030.DescribeRulesRequest) (*albv20251030.DescribeRulesResponse, error) {
		normalized := albRulesTestItem(t, ctx.sampleModel(t))
		normalized.RuleName = albRulesTestStr(albRulesTestCloudRuleName)

		return &albv20251030.DescribeRulesResponse{Response: &albv20251030.DescribeRulesResponseParams{
			Rules:     []*albv20251030.RuleOutput{normalized},
			RequestId: albRulesTestStr("fake-request-id"),
		}}, nil
	})

	stateModel := ctx.sampleModel(t)
	stateModel.Id = types.StringValue(albRulesTestCompositeId)

	plannedModel := ctx.sampleModel(t)
	plannedModel.Id = types.StringValue(albRulesTestCompositeId)
	plannedModel.RuleName = types.StringValue(albRulesTestRenamedRuleName)

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, stateModel)
	ctx.resource.Update(context.Background(), resource.UpdateRequest{
		Plan:  ctx.plan(t, plannedModel),
		State: ctx.state(t, stateModel),
	}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("ModifyRulesAttributes request is nil")
	}
	if captured.LoadBalancerId == nil || *captured.LoadBalancerId != albRulesTestLoadBalancerId {
		t.Errorf("update request load_balancer_id = %v, want %q", captured.LoadBalancerId, albRulesTestLoadBalancerId)
	}
	if captured.ListenerId == nil || *captured.ListenerId != albRulesTestListenerId {
		t.Errorf("update request listener_id = %v, want %q", captured.ListenerId, albRulesTestListenerId)
	}
	if len(captured.Rules) != 1 || captured.Rules[0] == nil {
		t.Fatalf("update request rules = %v, want exactly one rule", captured.Rules)
	}
	if captured.Rules[0].RuleId == nil || *captured.Rules[0].RuleId != albRulesTestRuleId {
		t.Errorf("update request rule_id = %v, want %q", captured.Rules[0].RuleId, albRulesTestRuleId)
	}
	if captured.Rules[0].RuleName == nil || *captured.Rules[0].RuleName != albRulesTestRenamedRuleName {
		t.Errorf("update request rule_name = %v, want the planned %q", captured.Rules[0].RuleName, albRulesTestRenamedRuleName)
	}

	var got svcalb.AlbRulesResourceModel
	if diags := updateResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != albRulesTestCompositeId {
		t.Errorf("id = %q, want the unchanged %q", got.Id.ValueString(), albRulesTestCompositeId)
	}
	if got.RuleName.ValueString() != albRulesTestCloudRuleName {
		t.Errorf("rule_name = %q, want the read-back cloud value %q", got.RuleName.ValueString(), albRulesTestCloudRuleName)
	}
}

// TestAlbRulesUpdateWithoutChanges verifies that no write API is called when
// none of the attributes covered by ModifyRulesAttributes changed.
func TestAlbRulesUpdateWithoutChanges(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbRulesTestContext(t)

	modifyCalls := 0
	patches.ApplyMethodFunc(ctx.sdk, "ModifyRulesAttributes", func(req *albv20251030.ModifyRulesAttributesRequest) (*albv20251030.ModifyRulesAttributesResponse, error) {
		modifyCalls++

		return nil, fmt.Errorf("the modify API must not be called")
	})
	patches.ApplyMethodFunc(ctx.sdk, "DescribeRules", func(req *albv20251030.DescribeRulesRequest) (*albv20251030.DescribeRulesResponse, error) {
		return &albv20251030.DescribeRulesResponse{Response: &albv20251030.DescribeRulesResponseParams{
			Rules:     []*albv20251030.RuleOutput{albRulesTestItem(t, ctx.sampleModel(t))},
			RequestId: albRulesTestStr("fake-request-id"),
		}}, nil
	})

	model := ctx.sampleModel(t)
	model.Id = types.StringValue(albRulesTestCompositeId)

	var updateResp resource.UpdateResponse
	updateResp.State = ctx.state(t, model)
	ctx.resource.Update(context.Background(), resource.UpdateRequest{
		Plan:  ctx.plan(t, model),
		State: ctx.state(t, model),
	}, &updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("Update failed: %v", updateResp.Diagnostics.Errors())
	}
	if modifyCalls != 0 {
		t.Errorf("ModifyRulesAttributes was called %d time(s) for an unchanged plan", modifyCalls)
	}
}

// TestAlbRulesDeleteSuccess verifies that Delete addresses the rule by the
// three id segments.
func TestAlbRulesDeleteSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ctx := newAlbRulesTestContext(t)

	var captured *albv20251030.DeleteRulesRequest
	patches.ApplyMethodFunc(ctx.sdk, "DeleteRules", func(req *albv20251030.DeleteRulesRequest) (*albv20251030.DeleteRulesResponse, error) {
		captured = req

		return &albv20251030.DeleteRulesResponse{Response: &albv20251030.DeleteRulesResponseParams{
			RequestId: albRulesTestStr("fake-request-id"),
		}}, nil
	})

	model := ctx.sampleModel(t)
	model.Id = types.StringValue(albRulesTestCompositeId)

	var deleteResp resource.DeleteResponse
	ctx.resource.Delete(context.Background(), resource.DeleteRequest{State: ctx.state(t, model)}, &deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed: %v", deleteResp.Diagnostics.Errors())
	}

	if captured == nil {
		t.Fatal("DeleteRules request is nil")
	}
	if captured.LoadBalancerId == nil || *captured.LoadBalancerId != albRulesTestLoadBalancerId {
		t.Errorf("delete request load_balancer_id = %v, want %q", captured.LoadBalancerId, albRulesTestLoadBalancerId)
	}
	if captured.ListenerId == nil || *captured.ListenerId != albRulesTestListenerId {
		t.Errorf("delete request listener_id = %v, want %q", captured.ListenerId, albRulesTestListenerId)
	}
	if len(captured.RuleIds) != 1 || captured.RuleIds[0] == nil || *captured.RuleIds[0] != albRulesTestRuleId {
		t.Errorf("delete request rule_ids = %v, want [%q]", captured.RuleIds, albRulesTestRuleId)
	}
}

// TestAlbRulesImportState verifies that the composite import id is split into
// the location attributes.
func TestAlbRulesImportState(t *testing.T) {
	ctx := newAlbRulesTestContext(t)

	importable, ok := ctx.resource.(resource.ResourceWithImportState)
	if !ok {
		t.Fatal("resource does not implement resource.ResourceWithImportState")
	}

	importResp := resource.ImportStateResponse{State: ctx.state(t, ctx.sampleModel(t))}
	importable.ImportState(context.Background(), resource.ImportStateRequest{ID: albRulesTestCompositeId}, &importResp)
	if importResp.Diagnostics.HasError() {
		t.Fatalf("ImportState failed: %v", importResp.Diagnostics.Errors())
	}

	var got svcalb.AlbRulesResourceModel
	if diags := importResp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags.Errors())
	}
	if got.Id.ValueString() != albRulesTestCompositeId {
		t.Errorf("id = %q, want %q", got.Id.ValueString(), albRulesTestCompositeId)
	}
	if got.LoadBalancerId.ValueString() != albRulesTestLoadBalancerId {
		t.Errorf("load_balancer_id = %q, want %q", got.LoadBalancerId.ValueString(), albRulesTestLoadBalancerId)
	}
	if got.ListenerId.ValueString() != albRulesTestListenerId {
		t.Errorf("listener_id = %q, want %q", got.ListenerId.ValueString(), albRulesTestListenerId)
	}
}
