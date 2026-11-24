package alb

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"
)

// --- Schema / Metadata / Factory ---

// TestAlbRulesResourceMetadata validates the type name.
func TestAlbRulesResourceMetadata(t *testing.T) {
	r := &AlbRulesResource{}
	req := resource.MetadataRequest{}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	if resp.TypeName != "tencentcloud_alb_rules" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_alb_rules")
	}
}

// TestAlbRulesResourceSchema validates required/optional/computed flags on
// the top-level attributes and the nested rule attributes.
func TestAlbRulesResourceSchema(t *testing.T) {
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
		"load_balancer_id": {required: true},
		"listener_id":      {required: true},
		"rules":            {required: true},
		"next_token":       {computed: true},
	}

	for name, want := range cases {
		attr, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Errorf("attribute %q missing from schema", name)
			continue
		}
		switch a := attr.(type) {
		case schema.StringAttribute:
			if a.Required != want.required {
				t.Errorf("attribute %q Required=%v, want %v", name, a.Required, want.required)
			}
			if a.Optional != want.optional {
				t.Errorf("attribute %q Optional=%v, want %v", name, a.Optional, want.optional)
			}
			if a.Computed != want.computed {
				t.Errorf("attribute %q Computed=%v, want %v", name, a.Computed, want.computed)
			}
		case schema.ListNestedAttribute:
			if a.Required != want.required {
				t.Errorf("attribute %q Required=%v, want %v", name, a.Required, want.required)
			}
			if a.Optional != want.optional {
				t.Errorf("attribute %q Optional=%v, want %v", name, a.Optional, want.optional)
			}
			if a.Computed != want.computed {
				t.Errorf("attribute %q Computed=%v, want %v", name, a.Computed, want.computed)
			}
		default:
			t.Errorf("attribute %q has unexpected type %T", name, attr)
		}
	}

	// Nested rule attributes.
	rulesAttr, ok := resp.Schema.Attributes["rules"].(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("rules attribute is not a ListNestedAttribute")
	}
	ruleCases := map[string]expect{
		"actions":    {required: true},
		"conditions": {required: true},
		"priority":   {required: true},
		"direction":  {optional: true, computed: true},
		"rule_name":  {optional: true, computed: true},
		"tags":       {optional: true},
		"rule_id":    {computed: true},
		"status":     {computed: true},
	}
	for name, want := range ruleCases {
		attr, ok := rulesAttr.NestedObject.Attributes[name]
		if !ok {
			t.Errorf("nested attribute %q missing from rules schema", name)
			continue
		}
		switch a := attr.(type) {
		case schema.Int64Attribute:
			if a.Required != want.required || a.Optional != want.optional || a.Computed != want.computed {
				t.Errorf("nested attribute %q Required=%v Optional=%v Computed=%v, want %v", name, a.Required, a.Optional, a.Computed, want)
			}
		case schema.StringAttribute:
			if a.Required != want.required || a.Optional != want.optional || a.Computed != want.computed {
				t.Errorf("nested attribute %q Required=%v Optional=%v Computed=%v, want %v", name, a.Required, a.Optional, a.Computed, want)
			}
		case schema.ListNestedAttribute:
			if a.Required != want.required || a.Optional != want.optional || a.Computed != want.computed {
				t.Errorf("nested attribute %q Required=%v Optional=%v Computed=%v, want %v", name, a.Required, a.Optional, a.Computed, want)
			}
		default:
			t.Errorf("nested attribute %q has unexpected type %T", name, attr)
		}
	}
}

// TestAlbRulesResourceFactory validates the factory return value.
func TestAlbRulesResourceFactory(t *testing.T) {
	res := NewAlbRulesResource()
	if res == nil {
		t.Fatal("NewAlbRulesResource returned nil")
	}
}

// --- Expand / flatten round trips ---

// TestExpandAlbRules maps a minimal plan into the SDK request structs.
func TestExpandAlbRules(t *testing.T) {
	priority := int64(1)
	rule := AlbRulesRuleModel{
		Priority: types.Int64Value(priority),
		Actions:  buildActionList(t, AlbRulesActionModel{Order: types.Int64Value(1), Type: types.StringValue("FixedResponse")}),
	}

	ruleInputs, expandDiags := expandAlbRules(context.Background(), buildRuleList(t, rule))
	if expandDiags.HasError() {
		t.Fatalf("expandAlbRules returned diagnostics: %v", expandDiags)
	}

	if len(ruleInputs) != 1 {
		t.Fatalf("len(ruleInputs) = %d, want 1", len(ruleInputs))
	}
	if ruleInputs[0].Priority == nil || *ruleInputs[0].Priority != priority {
		t.Errorf("rule priority = %v, want %v", ruleInputs[0].Priority, priority)
	}
	if len(ruleInputs[0].Actions) != 1 {
		t.Fatalf("len(actions) = %d, want 1", len(ruleInputs[0].Actions))
	}
	if ruleInputs[0].Actions[0].Type == nil || *ruleInputs[0].Actions[0].Type != "FixedResponse" {
		t.Errorf("action type = %v, want FixedResponse", ruleInputs[0].Actions[0].Type)
	}
}

// TestRefreshAlbRulesResourceModel flattens cloud rules into the model.
func TestRefreshAlbRulesResourceModel(t *testing.T) {
	ruleId := "rule-12345678"
	priority := int64(1)
	cloudRules := []*albv20251030.RuleOutput{
		{
			RuleId:   &ruleId,
			Priority: &priority,
		},
	}

	model := &AlbRulesResourceModel{}
	diags := &diag.Diagnostics{}
	refreshAlbRulesResourceModel(context.Background(), model, cloudRules, diags)
	if diags.HasError() {
		t.Fatalf("refreshAlbRulesResourceModel returned diagnostics: %v", diags.diags)
	}

	var ruleModels []AlbRulesRuleModel
	model.Rules.ElementsAs(context.Background(), &ruleModels, false)
	if len(ruleModels) != 1 {
		t.Fatalf("len(rules) = %d, want 1", len(ruleModels))
	}
	if ruleModels[0].RuleId.ValueString() != ruleId {
		t.Errorf("rule_id = %q, want %q", ruleModels[0].RuleId.ValueString(), ruleId)
	}
	if ruleModels[0].Priority.ValueInt64() != priority {
		t.Errorf("priority = %d, want %d", ruleModels[0].Priority.ValueInt64(), priority)
	}
}

// --- test helpers ---

// buildRuleList builds a rules types.List from one model.
func buildRuleList(t *testing.T, rule AlbRulesRuleModel) types.List {
	t.Helper()

	list, listDiags := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: albRulesAttrTypes()}, []AlbRulesRuleModel{rule})
	if listDiags.HasError() {
		t.Fatalf("buildRuleList returned diagnostics: %v", listDiags)
	}
	return list
}

// buildActionList builds an actions types.List from one model.
func buildActionList(t *testing.T, action AlbRulesActionModel) types.List {
	t.Helper()

	list, listDiags := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: albRulesActionAttrTypes()}, []AlbRulesActionModel{action})
	if listDiags.HasError() {
		t.Fatalf("buildActionList returned diagnostics: %v", listDiags)
	}
	return list
}
