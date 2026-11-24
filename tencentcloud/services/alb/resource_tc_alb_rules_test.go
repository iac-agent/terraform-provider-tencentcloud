package alb

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// --- Schema / Metadata / Factory ---

// TestAlbRulesResource_Metadata validates the type name.
func TestAlbRulesResource_Metadata(t *testing.T) {
	r := &AlbRulesResource{}
	req := resource.MetadataRequest{}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	if resp.TypeName != "tencentcloud_alb_rules" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_alb_rules")
	}
}

// TestAlbRulesResource_Schema validates required/optional/computed flags on
// the top-level attributes.
func TestAlbRulesResource_Schema(t *testing.T) {
	r := &AlbRulesResource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}

	type expect struct {
		required, optional, computed bool
	}
	cases := map[string]expect{
		"id":               {computed: true},
		"listener_id":      {required: true},
		"load_balancer_id": {required: true},
		"rules":            {required: true},
	}

	for name, want := range cases {
		attr, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Errorf("attribute %q missing from schema", name)
			continue
		}
		switch a := attr.(type) {
		case schema.StringAttribute:
			if a.Required != want.required || a.Optional != want.optional || a.Computed != want.computed {
				t.Errorf("attribute %q Required=%v Optional=%v Computed=%v, want %v", name, a.Required, a.Optional, a.Computed, want)
			}
		case schema.ListNestedAttribute:
			if a.Required != want.required || a.Optional != want.optional || a.Computed != want.computed {
				t.Errorf("attribute %q Required=%v Optional=%v Computed=%v, want %v", name, a.Required, a.Optional, a.Computed, want)
			}
		default:
			t.Errorf("attribute %q unexpected type %T", name, attr)
		}
	}

	// The two addressing attributes must force replacement.
	for _, name := range []string{"listener_id", "load_balancer_id"} {
		attr, ok := resp.Schema.Attributes[name].(schema.StringAttribute)
		if !ok {
			t.Fatalf("%s not a StringAttribute", name)
		}
		if len(attr.PlanModifiers) == 0 {
			t.Errorf("%s should have a RequiresReplace plan modifier", name)
		}
	}
}

// TestAlbRulesResource_ImplementsInterfaces verifies the resource satisfies
// the framework interfaces advertised in the _ asserts.
func TestAlbRulesResource_ImplementsInterfaces(t *testing.T) {
	var _ resource.ResourceWithConfigure = &AlbRulesResource{}
	var _ resource.ResourceWithImportState = &AlbRulesResource{}
}

// TestAlbRulesResource_Factory validates the registry factory.
func TestAlbRulesResource_Factory(t *testing.T) {
	got := NewAlbRulesResource()
	if got == nil {
		t.Fatal("NewAlbRulesResource returned nil")
	}
	if _, ok := got.(*AlbRulesResource); !ok {
		t.Fatalf("factory returned %T, want *AlbRulesResource", got)
	}
}

// TestAlbRulesResource_ConfigureNilProviderData ensures the pre-Configure
// invocation pattern (nil ProviderData) is silently no-op.
func TestAlbRulesResource_ConfigureNilProviderData(t *testing.T) {
	r := &AlbRulesResource{}
	req := resource.ConfigureRequest{}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Configure with nil ProviderData returned diagnostics: %v", resp.Diagnostics)
	}
}

// --- ID split ---

// TestSplitAlbRulesId validates parsing of the composite id.
func TestSplitAlbRulesId(t *testing.T) {
	good := "alb-xxxxxxxx" + tccommon.FILED_SP + "lst-yyyyyyyy" + tccommon.FILED_SP + "rule-zzzzzzzz"
	lb, lst, rule, ok := splitAlbRulesId(good)
	if !ok || lb != "alb-xxxxxxxx" || lst != "lst-yyyyyyyy" || rule != "rule-zzzzzzzz" {
		t.Errorf("splitAlbRulesId(%q) = (%q, %q, %q, %v)", good, lb, lst, rule, ok)
	}

	bad := []string{
		"",
		"alb-xxxxxxxx",
		"alb-xxxxxxxx" + tccommon.FILED_SP + "lst-yyyyyyyy",
		tccommon.FILED_SP + "lst" + tccommon.FILED_SP + "rule",
		"alb" + tccommon.FILED_SP + tccommon.FILED_SP + "rule-zzzzzzzz",
	}
	for _, in := range bad {
		if _, _, _, ok := splitAlbRulesId(in); ok {
			t.Errorf("splitAlbRulesId(%q) should fail", in)
		}
	}
}

// TestAlbRulesIdRoundTrip validates id building and parsing are inverse.
func TestAlbRulesIdRoundTrip(t *testing.T) {
	id := albRulesId("alb-aaaaaaaa", "lst-bbbbbbbb", "rule-cccccccc")
	lb, lst, rule, ok := splitAlbRulesId(id)
	if !ok || lb != "alb-aaaaaaaa" || lst != "lst-bbbbbbbb" || rule != "rule-cccccccc" {
		t.Errorf("round trip of %q failed: (%q, %q, %q, %v)", id, lb, lst, rule, ok)
	}
}

// --- Flattening ---

// TestFlattenAlbRules validates the rule list built from a DescribeRules item.
func TestFlattenAlbRules(t *testing.T) {
	item := &albv20251030.RuleOutput{
		RuleId:    helper.String("rule-cccccccc"),
		Priority:  helper.Int64(1),
		Direction: helper.String("Request"),
		RuleName:  helper.String("tf-rule"),
		Actions: []*albv20251030.RuleAction{
			{
				Order: helper.Int64(20),
				Type:  helper.String("TargetGroup"),
				TargetGroupConfig: &albv20251030.TargetGroupConfig{
					TargetGroups: []*albv20251030.TargetGroupTuple{
						{
							TargetGroupId: helper.String("lbtg-dddddddd"),
							Weight:        helper.Int64(10),
						},
					},
				},
			},
		},
		Conditions: []*albv20251030.RuleCondition{
			{
				Type:       helper.String("Path"),
				PathConfig: []*string{helper.String("/api")},
			},
		},
	}

	list, diags := flattenAlbRules(context.Background(), []*albv20251030.RuleOutput{item}, types.ListNull(albTagObjectType()))
	if diags.HasError() {
		t.Fatalf("flattenAlbRules returned diagnostics: %v", diags)
	}
	if list.IsNull() || len(list.Elements()) != 1 {
		t.Fatalf("flattenAlbRules produced %d rules, want 1", len(list.Elements()))
	}

	var rules []albRuleModel
	if diags := list.ElementsAs(context.Background(), &rules, false); diags.HasError() {
		t.Fatalf("cannot read flattened rules: %v", diags)
	}
	if len(rules) != 1 {
		t.Fatalf("got %d rules, want 1", len(rules))
	}
	if rules[0].Priority.ValueInt64() != 1 {
		t.Errorf("priority = %d, want 1", rules[0].Priority.ValueInt64())
	}
	if rules[0].Direction.ValueString() != "Request" {
		t.Errorf("direction = %q, want Request", rules[0].Direction.ValueString())
	}

	var actions []albRuleActionModel
	if diags := rules[0].Actions.ElementsAs(context.Background(), &actions, false); diags.HasError() {
		t.Fatalf("cannot read flattened actions: %v", diags)
	}
	if len(actions) != 1 || actions[0].Type.ValueString() != "TargetGroup" {
		t.Fatalf("unexpected actions: %+v", actions)
	}

	var conditions []albRuleConditionModel
	if diags := rules[0].Conditions.ElementsAs(context.Background(), &conditions, false); diags.HasError() {
		t.Fatalf("cannot read flattened conditions: %v", diags)
	}
	if len(conditions) != 1 || conditions[0].Type.ValueString() != "Path" {
		t.Fatalf("unexpected conditions: %+v", conditions)
	}
}

// TestFlattenAlbRulesAbsentConfigBlocks validates config blocks the cloud does
// not report stay null so they do not show up as spurious diffs.
func TestFlattenAlbRulesAbsentConfigBlocks(t *testing.T) {
	item := &albv20251030.RuleOutput{
		RuleId:   helper.String("rule-cccccccc"),
		Priority: helper.Int64(2),
		Actions: []*albv20251030.RuleAction{
			{Order: helper.Int64(1), Type: helper.String("Redirect")},
		},
		Conditions: []*albv20251030.RuleCondition{
			{Type: helper.String("Host")},
		},
	}

	list, diags := flattenAlbRules(context.Background(), []*albv20251030.RuleOutput{item}, types.ListNull(albTagObjectType()))
	if diags.HasError() {
		t.Fatalf("flattenAlbRules returned diagnostics: %v", diags)
	}

	var rules []albRuleModel
	if diags := list.ElementsAs(context.Background(), &rules, false); diags.HasError() {
		t.Fatalf("cannot read flattened rules: %v", diags)
	}

	var actions []albRuleActionModel
	if diags := rules[0].Actions.ElementsAs(context.Background(), &actions, false); diags.HasError() {
		t.Fatalf("cannot read flattened actions: %v", diags)
	}
	if !actions[0].FixedResponseConfig.IsNull() {
		t.Error("fixed_response_config should stay null when the cloud omits it")
	}
	if !actions[0].TargetGroupConfig.IsNull() {
		t.Error("target_group_config should stay null when the cloud omits it")
	}

	var conditions []albRuleConditionModel
	if diags := rules[0].Conditions.ElementsAs(context.Background(), &conditions, false); diags.HasError() {
		t.Fatalf("cannot read flattened conditions: %v", diags)
	}
	if !conditions[0].HeaderConfig.IsNull() {
		t.Error("header_config should stay null when the cloud omits it")
	}
	if !conditions[0].HostConfig.IsNull() {
		t.Error("host_config should stay null when the cloud omits it")
	}
}

// --- Tags ---

// TestAlbTagsChanged validates the tag change detection.
func TestAlbTagsChanged(t *testing.T) {
	objType := albTagObjectType()
	tagObj := types.ObjectValueMust(objType.(types.ObjectType).AttributeTypes(), map[string]attr.Value{
		"tag_key":   types.StringValue("env"),
		"tag_value": types.StringValue("test"),
	})
	planned := types.ListValueMust(objType, []attr.Value{tagObj})

	if !albTagsChanged(planned, types.ListNull(objType)) {
		t.Error("planned tags versus empty state should report a change")
	}
	if albTagsChanged(planned, planned) {
		t.Error("equal tags should not report a change")
	}
	if albTagsChanged(types.ListNull(objType), planned) {
		t.Error("a removed tags block is handled by read-back, not the change check")
	}
}
