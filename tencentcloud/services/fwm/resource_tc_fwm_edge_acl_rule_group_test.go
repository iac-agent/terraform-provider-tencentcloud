package fwm_test

import (
	"context"
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
)

const (
	edgeAclRuleGroupId       = "fwmrg-ed9v0example"
	edgeAclRuleGroupName     = "tf-example-edge-acl-group"
	edgeAclRuleGroupRuleId   = "rule-0001"
	edgeAclRuleGroupRuleIdV2 = "rule-0002"
)

// edgeAclRuleGroupTestModel mirrors the resource model so assertions can read
// every leaf field back from state (state -> struct -> typed values).
type edgeAclRuleGroupTestModel struct {
	Id        types.String `tfsdk:"id"`
	GroupId   types.String `tfsdk:"group_id"`
	GroupName types.String `tfsdk:"group_name"`
	Product   types.String `tfsdk:"product"`
	RuleCount types.Int64  `tfsdk:"rule_count"`
	Rules     types.List   `tfsdk:"rules"`
}

type edgeAclRuleGroupTestRuleModel struct {
	RuleId          types.String `tfsdk:"rule_id"`
	OrderIndex      types.Int64  `tfsdk:"order_index"`
	Direction       types.Int64  `tfsdk:"direction"`
	SourceContent   types.String `tfsdk:"source_content"`
	SourceType      types.String `tfsdk:"source_type"`
	TargetContent   types.String `tfsdk:"target_content"`
	TargetType      types.String `tfsdk:"target_type"`
	Port            types.String `tfsdk:"port"`
	Protocol        types.String `tfsdk:"protocol"`
	RuleAction      types.String `tfsdk:"rule_action"`
	Description     types.String `tfsdk:"description"`
	Scope           types.String `tfsdk:"scope"`
	BelongMemberId  types.String `tfsdk:"belong_member_id"`
	ParamTemplateId types.String `tfsdk:"param_template_id"`
}

// newEdgeAclRuleGroupResource builds the resource and injects the shared meta
// through the real Configure implementation.
func newEdgeAclRuleGroupResource(t *testing.T) (resource.Resource, *sharedmeta.ProviderMeta) {
	t.Helper()

	res := fwm.NewFwmEdgeAclRuleGroupResource()
	configurable, ok := res.(resource.ResourceWithConfigure)
	if !ok {
		t.Fatalf("resource does not implement resource.ResourceWithConfigure")
	}

	if _, ok := res.(resource.ResourceWithImportState); !ok {
		t.Fatalf("resource does not implement resource.ResourceWithImportState")
	}

	var configureResp resource.ConfigureResponse
	meta := &sharedmeta.ProviderMeta{Client: &connectivity.TencentCloudClient{}}
	configurable.Configure(context.Background(), resource.ConfigureRequest{ProviderData: meta}, &configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("configure returned diagnostics: %v", configureResp.Diagnostics)
	}

	return res, meta
}

func mustEdgeAclRuleGroupSchema(t *testing.T, res resource.Resource) *schema.Schema {
	t.Helper()

	var schemaResp resource.SchemaResponse
	res.Schema(context.Background(), resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("schema returned diagnostics: %v", schemaResp.Diagnostics)
	}
	return &schemaResp.Schema
}

// buildEdgeAclRuleGroupData builds a fully known tftypes value for the
// resource schema, with a null rules list.
func buildEdgeAclRuleGroupData(
	t *testing.T,
	sch *schema.Schema,
	id *string,
	groupName *string,
	ruleCount *int64,
) tftypes.Value {
	t.Helper()

	objType := sch.Type().TerraformType(context.Background()).(tftypes.Object)

	stringValue := func(v *string) tftypes.Value {
		if v == nil {
			return tftypes.NewValue(tftypes.String, nil)
		}
		return tftypes.NewValue(tftypes.String, *v)
	}
	numberValue := func(v *int64) tftypes.Value {
		if v == nil {
			return tftypes.NewValue(tftypes.Number, nil)
		}
		return tftypes.NewValue(tftypes.Number, *v)
	}

	product := fwm.EdgeAclRuleGroupProduct
	return tftypes.NewValue(objType, map[string]tftypes.Value{
		"id":         stringValue(id),
		"group_id":   stringValue(id),
		"group_name": stringValue(groupName),
		"product":    stringValue(&product),
		"rule_count": numberValue(ruleCount),
		"rules":      tftypes.NewValue(tftypes.List{ElementType: ruleObjectType(t, sch)}, nil),
	})
}

func mockEdgeAclRuleGroupClient(patches *gomonkey.Patches, meta *sharedmeta.ProviderMeta) *fwmv20250611.Client {
	client := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(meta.Client, "UseFwmV20250611Client", client)
	return client
}

func buildEdgeAclRuleGroupResp() *fwmv20250611.RuleGroupResp {
	return &fwmv20250611.RuleGroupResp{
		AppId:      ptrString("1250000000"),
		GroupId:    ptrString(edgeAclRuleGroupId),
		GroupName:  ptrString(edgeAclRuleGroupName),
		RuleCount:  ptrInt64(1),
		Product:    ptrString(fwm.EdgeAclRuleGroupProduct),
		CreateTime: ptrString("2025-01-01 00:00:00"),
		UpdateTime: ptrString("2025-01-02 00:00:00"),
		Direction:  ptrInt64(1),
	}
}

func buildEdgeAclRuleGroupRuleResp() *fwmv20250611.EdgeAclRuleResp {
	return &fwmv20250611.EdgeAclRuleResp{
		RuleId:          ptrString(edgeAclRuleGroupRuleId),
		Sequence:        ptrInt64(1),
		Direction:       ptrInt64(1),
		SourceContent:   ptrString("0.0.0.0/0"),
		SourceType:      ptrString("ip"),
		TargetContent:   ptrString("1.1.1.1"),
		TargetType:      ptrString("ip"),
		Port:            ptrString("80"),
		Protocol:        ptrString("TCP"),
		RuleAction:      ptrString("accept"),
		Description:     ptrString("allow web traffic"),
		Scope:           ptrString("all"),
		ParamTemplateId: ptrString("tmpl-0001"),
		BelongMember: &fwmv20250611.MemberInfo{
			MemberId: ptrString("mem-0001"),
		},
	}
}

func mockDescribeRuleGroupsByGroup(patches *gomonkey.Patches, client *fwmv20250611.Client, groups ...*fwmv20250611.RuleGroupResp) {
	patches.ApplyMethodFunc(client, "DescribeRuleGroupsWithContext", func(_ context.Context, _ *fwmv20250611.DescribeRuleGroupsRequest) (*fwmv20250611.DescribeRuleGroupsResponse, error) {
		resp := fwmv20250611.NewDescribeRuleGroupsResponse()
		resp.Response = &fwmv20250611.DescribeRuleGroupsResponseParams{
			TotalCount: ptrInt64(int64(len(groups))),
			RuleGroups: groups,
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})
}

func mockDescribeEdgeAclRulesByDirection(patches *gomonkey.Patches, client *fwmv20250611.Client, rulesByDirection map[int64][]*fwmv20250611.EdgeAclRuleResp) {
	patches.ApplyMethodFunc(client, "DescribeEdgeAclRulesWithContext", func(_ context.Context, request *fwmv20250611.DescribeEdgeAclRulesRequest) (*fwmv20250611.DescribeEdgeAclRulesResponse, error) {
		var direction int64
		if request != nil && request.Direction != nil {
			direction = *request.Direction
		}
		rules := rulesByDirection[direction]
		resp := fwmv20250611.NewDescribeEdgeAclRulesResponse()
		resp.Response = &fwmv20250611.DescribeEdgeAclRulesResponseParams{
			TotalCount: ptrInt64(int64(len(rules))),
			Rules:      rules,
			RequestId:  ptrString("fake-request-id"),
		}
		return resp, nil
	})
}

// TestFwmEdgeAclRuleGroup_Schema pins the schema contract: group_name is the
// only updatable attribute, everything else is either computed or replaced.
func TestFwmEdgeAclRuleGroup_Schema(t *testing.T) {
	res := fwm.NewFwmEdgeAclRuleGroupResource()

	var metadataResp resource.MetadataResponse
	res.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "tencentcloud"}, &metadataResp)
	if metadataResp.TypeName != "tencentcloud_fwm_edge_acl_rule_group" {
		t.Fatalf("unexpected type name %q", metadataResp.TypeName)
	}

	sch := mustEdgeAclRuleGroupSchema(t, res)

	for _, name := range []string{"id", "group_id", "rule_count", "rules"} {
		attr, ok := sch.Attributes[name]
		if !ok {
			t.Fatalf("attribute %q is missing from the schema", name)
		}
		if !attr.IsComputed() {
			t.Fatalf("attribute %q must be computed", name)
		}
	}

	groupName, ok := sch.Attributes["group_name"]
	if !ok {
		t.Fatal("attribute group_name is missing from the schema")
	}
	if groupName.IsRequired() || !groupName.IsOptional() || !groupName.IsComputed() {
		t.Fatal("group_name must be optional and computed")
	}

	product, ok := sch.Attributes["product"]
	if !ok {
		t.Fatal("attribute product is missing from the schema")
	}
	if !product.IsOptional() || !product.IsComputed() {
		t.Fatal("product must be optional and computed")
	}

	rulesAttr, ok := sch.Attributes["rules"]
	if !ok {
		t.Fatal("attribute rules is missing from the schema")
	}
	if !rulesAttr.IsOptional() || !rulesAttr.IsComputed() {
		t.Fatal("rules must be optional and computed")
	}

	rules, ok := rulesAttr.(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("rules must be a ListNestedAttribute, got %T", rulesAttr)
	}

	expected := []string{
		"rule_id", "order_index", "direction", "source_content", "source_type",
		"target_content", "target_type", "port", "protocol", "rule_action",
		"description", "scope", "belong_member_id", "param_template_id",
	}
	if len(rules.NestedObject.Attributes) != len(expected) {
		t.Fatalf("unexpected number of rule attributes: %d", len(rules.NestedObject.Attributes))
	}
	for _, name := range expected {
		if _, ok := rules.NestedObject.Attributes[name]; !ok {
			t.Fatalf("rule attribute %q is missing from the schema", name)
		}
	}

	ruleId, ok := rules.NestedObject.Attributes["rule_id"]
	if !ok || !ruleId.IsComputed() || ruleId.IsOptional() {
		t.Fatal("rule_id must be computed only")
	}
}

// ruleRawValue builds the tftypes value of one rules element.
func ruleRawValue(t *testing.T, ruleType tftypes.Object, values map[string]tftypes.Value) tftypes.Value {
	t.Helper()
	return tftypes.NewValue(ruleType, values)
}

// ruleObjectType returns the tftypes.Object of one rules element.
func ruleObjectType(t *testing.T, sch *schema.Schema) tftypes.Object {
	t.Helper()

	rules, ok := sch.Attributes["rules"].(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("rules must be a ListNestedAttribute")
	}

	attrTypes := map[string]tftypes.Type{}
	for name := range rules.NestedObject.Attributes {
		if name == "order_index" || name == "direction" {
			attrTypes[name] = tftypes.Number
			continue
		}
		attrTypes[name] = tftypes.String
	}
	return tftypes.Object{AttributeTypes: attrTypes}
}

// buildRulesPlanValue builds the tftypes list for one rule with the fields
// used by the create test.
func buildRulesPlanValue(t *testing.T, sch *schema.Schema) tftypes.Value {
	t.Helper()

	ruleType := ruleObjectType(t, sch)
	rule := ruleRawValue(t, ruleType, map[string]tftypes.Value{
		"rule_id":           tftypes.NewValue(tftypes.String, nil),
		"order_index":       tftypes.NewValue(tftypes.Number, 1),
		"direction":         tftypes.NewValue(tftypes.Number, 1),
		"source_content":    tftypes.NewValue(tftypes.String, "0.0.0.0/0"),
		"source_type":       tftypes.NewValue(tftypes.String, "ip"),
		"target_content":    tftypes.NewValue(tftypes.String, "1.1.1.1"),
		"target_type":       tftypes.NewValue(tftypes.String, "ip"),
		"port":              tftypes.NewValue(tftypes.String, "80"),
		"protocol":          tftypes.NewValue(tftypes.String, "TCP"),
		"rule_action":       tftypes.NewValue(tftypes.String, "accept"),
		"description":       tftypes.NewValue(tftypes.String, "allow web traffic"),
		"scope":             tftypes.NewValue(tftypes.String, "all"),
		"belong_member_id":  tftypes.NewValue(tftypes.String, nil),
		"param_template_id": tftypes.NewValue(tftypes.String, nil),
	})

	return tftypes.NewValue(tftypes.List{ElementType: ruleType}, []tftypes.Value{rule})
}

// buildCreatePlan builds a create-time plan (id, group_id, product and
// rule_count are unknown).
func buildCreatePlan(t *testing.T, sch *schema.Schema, groupName string, withRules bool) tftypes.Value {
	t.Helper()

	rulesValue := tftypes.NewValue(tftypes.List{ElementType: ruleObjectType(t, sch)}, nil)
	if withRules {
		rulesValue = buildRulesPlanValue(t, sch)
	}

	return tftypes.NewValue(sch.Type().TerraformType(context.Background()), map[string]tftypes.Value{
		"id":         tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"group_id":   tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"group_name": tftypes.NewValue(tftypes.String, groupName),
		"product":    tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"rule_count": tftypes.NewValue(tftypes.Number, tftypes.UnknownValue),
		"rules":      rulesValue,
	})
}

// TestFwmEdgeAclRuleGroup_Create covers the happy path: the group is created
// with its rules and every leaf field is read back into state.
func TestFwmEdgeAclRuleGroup_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res, meta := newEdgeAclRuleGroupResource(t)
	client := mockEdgeAclRuleGroupClient(patches, meta)
	sch := mustEdgeAclRuleGroupSchema(t, res)
	ctx := context.Background()

	var captured *fwmv20250611.CreateEdgeAclRuleGroupRequest
	patches.ApplyMethodFunc(client, "CreateEdgeAclRuleGroupWithContext", func(_ context.Context, request *fwmv20250611.CreateEdgeAclRuleGroupRequest) (*fwmv20250611.CreateEdgeAclRuleGroupResponse, error) {
		captured = request
		resp := fwmv20250611.NewCreateEdgeAclRuleGroupResponse()
		resp.Response = &fwmv20250611.CreateEdgeAclRuleGroupResponseParams{
			GroupId:   ptrString(edgeAclRuleGroupId),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	mockDescribeRuleGroupsByGroup(patches, client, buildEdgeAclRuleGroupResp())
	mockDescribeEdgeAclRulesByDirection(patches, client, map[int64][]*fwmv20250611.EdgeAclRuleResp{
		1: {buildEdgeAclRuleGroupRuleResp()},
	})

	planRaw := buildCreatePlan(t, sch, edgeAclRuleGroupName, true)
	createResp := &resource.CreateResponse{
		State: tfsdk.State{Schema: *sch, Raw: tftypes.NewValue(sch.Type().TerraformType(ctx), nil)},
	}
	res.Create(ctx, resource.CreateRequest{Plan: tfsdk.Plan{Schema: *sch, Raw: planRaw}}, createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("create returned diagnostics: %v", createResp.Diagnostics)
	}

	if captured == nil {
		t.Fatal("CreateEdgeAclRuleGroup was not called")
	}
	if captured.Product == nil || *captured.Product != fwm.EdgeAclRuleGroupProduct {
		t.Fatalf("product must be hard-coded to %q, got %v", fwm.EdgeAclRuleGroupProduct, captured.Product)
	}
	if captured.GroupName == nil || *captured.GroupName != edgeAclRuleGroupName {
		t.Fatalf("unexpected group name sent to the API: %v", captured.GroupName)
	}
	if len(captured.Rules) != 1 {
		t.Fatalf("expected 1 rule in the create request, got %d", len(captured.Rules))
	}

	sentRule := captured.Rules[0]
	if sentRule.RuleId != nil {
		t.Fatal("rule_id must not be sent when creating a rule group")
	}
	if sentRule.OrderIndex == nil || *sentRule.OrderIndex != 1 {
		t.Fatalf("unexpected order_index sent: %v", sentRule.OrderIndex)
	}
	if sentRule.Direction == nil || *sentRule.Direction != 1 {
		t.Fatalf("unexpected direction sent: %v", sentRule.Direction)
	}
	if sentRule.SourceContent == nil || *sentRule.SourceContent != "0.0.0.0/0" {
		t.Fatalf("unexpected source_content sent: %v", sentRule.SourceContent)
	}
	if sentRule.TargetContent == nil || *sentRule.TargetContent != "1.1.1.1" {
		t.Fatalf("unexpected target_content sent: %v", sentRule.TargetContent)
	}
	if sentRule.Port == nil || *sentRule.Port != "80" {
		t.Fatalf("unexpected port sent: %v", sentRule.Port)
	}
	if sentRule.Protocol == nil || *sentRule.Protocol != "TCP" {
		t.Fatalf("unexpected protocol sent: %v", sentRule.Protocol)
	}
	if sentRule.RuleAction == nil || *sentRule.RuleAction != "accept" {
		t.Fatalf("unexpected rule_action sent: %v", sentRule.RuleAction)
	}
	if sentRule.Description == nil || *sentRule.Description != "allow web traffic" {
		t.Fatalf("unexpected description sent: %v", sentRule.Description)
	}
	if sentRule.Scope == nil || *sentRule.Scope != "all" {
		t.Fatalf("unexpected scope sent: %v", sentRule.Scope)
	}
	if sentRule.BelongMemberId != nil {
		t.Fatalf("belong_member_id must not be sent when it is null, got %v", *sentRule.BelongMemberId)
	}

	var got edgeAclRuleGroupTestModel
	if diags := createResp.State.Get(ctx, &got); diags.HasError() {
		t.Fatalf("reading state returned diagnostics: %v", diags)
	}

	if got.Id.ValueString() != edgeAclRuleGroupId {
		t.Fatalf("unexpected id %q", got.Id.ValueString())
	}
	if got.GroupId.ValueString() != edgeAclRuleGroupId {
		t.Fatalf("unexpected group_id %q", got.GroupId.ValueString())
	}
	if got.GroupName.ValueString() != edgeAclRuleGroupName {
		t.Fatalf("unexpected group_name %q", got.GroupName.ValueString())
	}
	if got.Product.ValueString() != fwm.EdgeAclRuleGroupProduct {
		t.Fatalf("unexpected product %q", got.Product.ValueString())
	}
	if got.RuleCount.ValueInt64() != 1 {
		t.Fatalf("unexpected rule_count %d", got.RuleCount.ValueInt64())
	}

	var rules []edgeAclRuleGroupTestRuleModel
	if diags := got.Rules.ElementsAs(ctx, &rules, false); diags.HasError() {
		t.Fatalf("reading rules returned diagnostics: %v", diags)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule in state, got %d", len(rules))
	}

	rule := rules[0]
	checks := map[string]struct {
		got  string
		want string
	}{
		"rule_id":           {rule.RuleId.ValueString(), edgeAclRuleGroupRuleId},
		"source_content":    {rule.SourceContent.ValueString(), "0.0.0.0/0"},
		"source_type":       {rule.SourceType.ValueString(), "ip"},
		"target_content":    {rule.TargetContent.ValueString(), "1.1.1.1"},
		"target_type":       {rule.TargetType.ValueString(), "ip"},
		"port":              {rule.Port.ValueString(), "80"},
		"protocol":          {rule.Protocol.ValueString(), "TCP"},
		"rule_action":       {rule.RuleAction.ValueString(), "accept"},
		"description":       {rule.Description.ValueString(), "allow web traffic"},
		"scope":             {rule.Scope.ValueString(), "all"},
		"belong_member_id":  {rule.BelongMemberId.ValueString(), "mem-0001"},
		"param_template_id": {rule.ParamTemplateId.ValueString(), "tmpl-0001"},
	}
	for name, c := range checks {
		if c.got != c.want {
			t.Fatalf("rule field %s: got %q, want %q", name, c.got, c.want)
		}
	}
	if rule.OrderIndex.ValueInt64() != 1 {
		t.Fatalf("rule field order_index: got %d, want 1", rule.OrderIndex.ValueInt64())
	}
	if rule.Direction.ValueInt64() != 1 {
		t.Fatalf("rule field direction: got %d, want 1", rule.Direction.ValueInt64())
	}
}

// TestFwmEdgeAclRuleGroup_Create_EmptyGroupId verifies an empty cloud id fails
// Create instead of writing an unusable id into state.
func TestFwmEdgeAclRuleGroup_Create_EmptyGroupId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res, meta := newEdgeAclRuleGroupResource(t)
	client := mockEdgeAclRuleGroupClient(patches, meta)
	sch := mustEdgeAclRuleGroupSchema(t, res)
	ctx := context.Background()

	patches.ApplyMethodFunc(client, "CreateEdgeAclRuleGroupWithContext", func(_ context.Context, _ *fwmv20250611.CreateEdgeAclRuleGroupRequest) (*fwmv20250611.CreateEdgeAclRuleGroupResponse, error) {
		resp := fwmv20250611.NewCreateEdgeAclRuleGroupResponse()
		resp.Response = &fwmv20250611.CreateEdgeAclRuleGroupResponseParams{RequestId: ptrString("fake-request-id")}
		return resp, nil
	})

	planRaw := buildCreatePlan(t, sch, edgeAclRuleGroupName, false)
	createResp := &resource.CreateResponse{
		State: tfsdk.State{Schema: *sch, Raw: tftypes.NewValue(sch.Type().TerraformType(ctx), nil)},
	}
	res.Create(ctx, resource.CreateRequest{Plan: tfsdk.Plan{Schema: *sch, Raw: planRaw}}, createResp)

	if !createResp.Diagnostics.HasError() {
		t.Fatal("expected an error when the API returns an empty group id")
	}
}

// TestFwmEdgeAclRuleGroup_Create_ApiError verifies API failures are surfaced.
func TestFwmEdgeAclRuleGroup_Create_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res, meta := newEdgeAclRuleGroupResource(t)
	client := mockEdgeAclRuleGroupClient(patches, meta)
	sch := mustEdgeAclRuleGroupSchema(t, res)
	ctx := context.Background()

	patches.ApplyMethodFunc(client, "CreateEdgeAclRuleGroupWithContext", func(_ context.Context, _ *fwmv20250611.CreateEdgeAclRuleGroupRequest) (*fwmv20250611.CreateEdgeAclRuleGroupResponse, error) {
		return nil, errors.New("fake create failure")
	})

	planRaw := buildCreatePlan(t, sch, edgeAclRuleGroupName, false)
	createResp := &resource.CreateResponse{
		State: tfsdk.State{Schema: *sch, Raw: tftypes.NewValue(sch.Type().TerraformType(ctx), nil)},
	}
	res.Create(ctx, resource.CreateRequest{Plan: tfsdk.Plan{Schema: *sch, Raw: planRaw}}, createResp)

	if !createResp.Diagnostics.HasError() {
		t.Fatal("expected an error when the create API fails")
	}
}

// buildExistingState builds a state value for an existing rule group.
func buildExistingState(t *testing.T, sch *schema.Schema, groupName string) tftypes.Value {
	t.Helper()

	return buildEdgeAclRuleGroupData(
		t,
		sch,
		ptrString(edgeAclRuleGroupId),
		ptrString(groupName),
		ptrInt64(1),
	)
}

// TestFwmEdgeAclRuleGroup_Read covers the refresh path, including the rules
// that are only visible through DescribeEdgeAclRules.
func TestFwmEdgeAclRuleGroup_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res, meta := newEdgeAclRuleGroupResource(t)
	client := mockEdgeAclRuleGroupClient(patches, meta)
	sch := mustEdgeAclRuleGroupSchema(t, res)
	ctx := context.Background()

	mockDescribeRuleGroupsByGroup(patches, client, buildEdgeAclRuleGroupResp())
	mockDescribeEdgeAclRulesByDirection(patches, client, map[int64][]*fwmv20250611.EdgeAclRuleResp{
		1: {buildEdgeAclRuleGroupRuleResp()},
	})

	stateRaw := buildExistingState(t, sch, "stale-name")
	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: *sch, Raw: stateRaw}}
	res.Read(ctx, resource.ReadRequest{State: tfsdk.State{Schema: *sch, Raw: stateRaw}}, readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("read returned diagnostics: %v", readResp.Diagnostics)
	}

	var got edgeAclRuleGroupTestModel
	if diags := readResp.State.Get(ctx, &got); diags.HasError() {
		t.Fatalf("reading state returned diagnostics: %v", diags)
	}

	if got.Id.ValueString() != edgeAclRuleGroupId {
		t.Fatalf("unexpected id %q", got.Id.ValueString())
	}
	if got.GroupName.ValueString() != edgeAclRuleGroupName {
		t.Fatalf("unexpected group_name %q", got.GroupName.ValueString())
	}
	if got.RuleCount.ValueInt64() != 1 {
		t.Fatalf("unexpected rule_count %d", got.RuleCount.ValueInt64())
	}

	var rules []edgeAclRuleGroupTestRuleModel
	if diags := got.Rules.ElementsAs(ctx, &rules, false); diags.HasError() {
		t.Fatalf("reading rules returned diagnostics: %v", diags)
	}
	if len(rules) != 1 || rules[0].RuleId.ValueString() != edgeAclRuleGroupRuleId {
		t.Fatalf("unexpected rules in state: %v", rules)
	}
	if rules[0].Protocol.ValueString() != "TCP" {
		t.Fatalf("unexpected protocol in state: %q", rules[0].Protocol.ValueString())
	}
}

// TestFwmEdgeAclRuleGroup_Read_NotFound verifies a successful query that
// matches nothing removes the resource from state.
func TestFwmEdgeAclRuleGroup_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res, meta := newEdgeAclRuleGroupResource(t)
	client := mockEdgeAclRuleGroupClient(patches, meta)
	sch := mustEdgeAclRuleGroupSchema(t, res)
	ctx := context.Background()

	mockDescribeRuleGroupsByGroup(patches, client)

	stateRaw := buildExistingState(t, sch, edgeAclRuleGroupName)
	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: *sch, Raw: stateRaw}}
	res.Read(ctx, resource.ReadRequest{State: tfsdk.State{Schema: *sch, Raw: stateRaw}}, readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("read returned diagnostics: %v", readResp.Diagnostics)
	}

	if !readResp.State.Raw.IsNull() {
		t.Fatal("expected the resource to be removed from state when it does not exist")
	}
}

// TestFwmEdgeAclRuleGroup_Read_NilResponse verifies a nil API response is an
// error and never treated as "does not exist".
func TestFwmEdgeAclRuleGroup_Read_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res, meta := newEdgeAclRuleGroupResource(t)
	client := mockEdgeAclRuleGroupClient(patches, meta)
	sch := mustEdgeAclRuleGroupSchema(t, res)
	ctx := context.Background()

	patches.ApplyMethodFunc(client, "DescribeRuleGroupsWithContext", func(_ context.Context, _ *fwmv20250611.DescribeRuleGroupsRequest) (*fwmv20250611.DescribeRuleGroupsResponse, error) {
		return nil, nil
	})

	stateRaw := buildExistingState(t, sch, edgeAclRuleGroupName)
	readResp := &resource.ReadResponse{State: tfsdk.State{Schema: *sch, Raw: stateRaw}}
	res.Read(ctx, resource.ReadRequest{State: tfsdk.State{Schema: *sch, Raw: stateRaw}}, readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("expected an error when the describe API returns a nil response")
	}
}

// TestFwmEdgeAclRuleGroup_Update_Rename verifies the only supported in-place
// update goes through ModifyRuleGroup with the managed group id.
func TestFwmEdgeAclRuleGroup_Update_Rename(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res, meta := newEdgeAclRuleGroupResource(t)
	client := mockEdgeAclRuleGroupClient(patches, meta)
	sch := mustEdgeAclRuleGroupSchema(t, res)
	ctx := context.Background()

	var captured *fwmv20250611.ModifyRuleGroupRequest
	patches.ApplyMethodFunc(client, "ModifyRuleGroupWithContext", func(_ context.Context, request *fwmv20250611.ModifyRuleGroupRequest) (*fwmv20250611.ModifyRuleGroupResponse, error) {
		captured = request
		resp := fwmv20250611.NewModifyRuleGroupResponse()
		resp.Response = &fwmv20250611.ModifyRuleGroupResponseParams{RequestId: ptrString("fake-request-id")}
		return resp, nil
	})

	renamed := buildEdgeAclRuleGroupResp()
	renamed.GroupName = ptrString("renamed-group")
	mockDescribeRuleGroupsByGroup(patches, client, renamed)
	mockDescribeEdgeAclRulesByDirection(patches, client, map[int64][]*fwmv20250611.EdgeAclRuleResp{
		1: {buildEdgeAclRuleGroupRuleResp()},
	})

	stateRaw := buildExistingState(t, sch, edgeAclRuleGroupName)
	planRaw := buildExistingState(t, sch, "renamed-group")

	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: *sch, Raw: stateRaw}}
	res.Update(ctx,
		resource.UpdateRequest{
			Plan:  tfsdk.Plan{Schema: *sch, Raw: planRaw},
			State: tfsdk.State{Schema: *sch, Raw: stateRaw},
		},
		updateResp,
	)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("update returned diagnostics: %v", updateResp.Diagnostics)
	}

	if captured == nil {
		t.Fatal("ModifyRuleGroup was not called")
	}
	if captured.GroupId == nil || *captured.GroupId != edgeAclRuleGroupId {
		t.Fatalf("unexpected group id sent to the API: %v", captured.GroupId)
	}
	if captured.GroupName == nil || *captured.GroupName != "renamed-group" {
		t.Fatalf("unexpected group name sent to the API: %v", captured.GroupName)
	}

	var got edgeAclRuleGroupTestModel
	if diags := updateResp.State.Get(ctx, &got); diags.HasError() {
		t.Fatalf("reading state returned diagnostics: %v", diags)
	}
	if got.GroupName.ValueString() != "renamed-group" {
		t.Fatalf("unexpected group_name %q", got.GroupName.ValueString())
	}
	if got.Id.ValueString() != edgeAclRuleGroupId {
		t.Fatalf("unexpected id %q", got.Id.ValueString())
	}
}

// TestFwmEdgeAclRuleGroup_Update_ClearNameRejected verifies clearing
// group_name fails loudly instead of sending an undocumented empty value.
func TestFwmEdgeAclRuleGroup_Update_ClearNameRejected(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res, meta := newEdgeAclRuleGroupResource(t)
	client := mockEdgeAclRuleGroupClient(patches, meta)
	sch := mustEdgeAclRuleGroupSchema(t, res)
	ctx := context.Background()

	modifyCalled := false
	patches.ApplyMethodFunc(client, "ModifyRuleGroupWithContext", func(_ context.Context, _ *fwmv20250611.ModifyRuleGroupRequest) (*fwmv20250611.ModifyRuleGroupResponse, error) {
		modifyCalled = true
		return fwmv20250611.NewModifyRuleGroupResponse(), nil
	})

	stateRaw := buildExistingState(t, sch, edgeAclRuleGroupName)
	planRaw := buildEdgeAclRuleGroupData(t, sch, ptrString(edgeAclRuleGroupId), nil, ptrInt64(1))

	updateResp := &resource.UpdateResponse{State: tfsdk.State{Schema: *sch, Raw: stateRaw}}
	res.Update(ctx,
		resource.UpdateRequest{
			Plan:  tfsdk.Plan{Schema: *sch, Raw: planRaw},
			State: tfsdk.State{Schema: *sch, Raw: stateRaw},
		},
		updateResp,
	)

	if !updateResp.Diagnostics.HasError() {
		t.Fatal("expected an error when group_name is cleared")
	}
	if modifyCalled {
		t.Fatal("ModifyRuleGroup must not be called with an empty group name")
	}
}

// TestFwmEdgeAclRuleGroup_Delete verifies exactly one group id is sent to the
// batch delete API.
func TestFwmEdgeAclRuleGroup_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res, meta := newEdgeAclRuleGroupResource(t)
	client := mockEdgeAclRuleGroupClient(patches, meta)
	sch := mustEdgeAclRuleGroupSchema(t, res)
	ctx := context.Background()

	var captured *fwmv20250611.DeleteRuleGroupRequest
	patches.ApplyMethodFunc(client, "DeleteRuleGroupWithContext", func(_ context.Context, request *fwmv20250611.DeleteRuleGroupRequest) (*fwmv20250611.DeleteRuleGroupResponse, error) {
		captured = request
		resp := fwmv20250611.NewDeleteRuleGroupResponse()
		resp.Response = &fwmv20250611.DeleteRuleGroupResponseParams{RequestId: ptrString("fake-request-id")}
		return resp, nil
	})

	stateRaw := buildExistingState(t, sch, edgeAclRuleGroupName)
	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: *sch, Raw: stateRaw}}
	res.Delete(ctx, resource.DeleteRequest{State: tfsdk.State{Schema: *sch, Raw: stateRaw}}, deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("delete returned diagnostics: %v", deleteResp.Diagnostics)
	}

	if captured == nil {
		t.Fatal("DeleteRuleGroup was not called")
	}
	if len(captured.GroupIds) != 1 {
		t.Fatalf("expected exactly one group id, got %d", len(captured.GroupIds))
	}
	if captured.GroupIds[0] == nil || *captured.GroupIds[0] != edgeAclRuleGroupId {
		t.Fatalf("unexpected group id: %v", captured.GroupIds[0])
	}
}

// TestFwmEdgeAclRuleGroup_Delete_ApiError verifies delete failures keep the
// resource in state.
func TestFwmEdgeAclRuleGroup_Delete_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res, meta := newEdgeAclRuleGroupResource(t)
	client := mockEdgeAclRuleGroupClient(patches, meta)
	sch := mustEdgeAclRuleGroupSchema(t, res)
	ctx := context.Background()

	patches.ApplyMethodFunc(client, "DeleteRuleGroupWithContext", func(_ context.Context, _ *fwmv20250611.DeleteRuleGroupRequest) (*fwmv20250611.DeleteRuleGroupResponse, error) {
		return nil, errors.New("fake delete failure")
	})

	stateRaw := buildExistingState(t, sch, edgeAclRuleGroupName)
	deleteResp := &resource.DeleteResponse{State: tfsdk.State{Schema: *sch, Raw: stateRaw}}
	res.Delete(ctx, resource.DeleteRequest{State: tfsdk.State{Schema: *sch, Raw: stateRaw}}, deleteResp)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatal("expected an error when the delete API fails")
	}
}

// TestFwmEdgeAclRuleGroup_Import verifies the group id is passed through.
func TestFwmEdgeAclRuleGroup_Import(t *testing.T) {
	res := fwm.NewFwmEdgeAclRuleGroupResource()
	sch := mustEdgeAclRuleGroupSchema(t, res)
	ctx := context.Background()

	importable, ok := res.(resource.ResourceWithImportState)
	if !ok {
		t.Fatal("resource does not implement resource.ResourceWithImportState")
	}

	state := tfsdk.State{Schema: *sch, Raw: tftypes.NewValue(sch.Type().TerraformType(ctx), nil)}
	importResp := &resource.ImportStateResponse{State: state}
	importable.ImportState(ctx, resource.ImportStateRequest{ID: edgeAclRuleGroupId}, importResp)
	if importResp.Diagnostics.HasError() {
		t.Fatalf("import returned diagnostics: %v", importResp.Diagnostics)
	}

	var got edgeAclRuleGroupTestModel
	if diags := importResp.State.Get(ctx, &got); diags.HasError() {
		t.Fatalf("reading state returned diagnostics: %v", diags)
	}
	if got.Id.ValueString() != edgeAclRuleGroupId {
		t.Fatalf("unexpected imported id %q", got.Id.ValueString())
	}
}
