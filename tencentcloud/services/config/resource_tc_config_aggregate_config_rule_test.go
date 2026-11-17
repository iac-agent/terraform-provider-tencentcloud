package config

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// --- Schema / Metadata / Factory ---

// TestConfigAggregateConfigRuleResource_Metadata validates the type name.
func TestConfigAggregateConfigRuleResource_Metadata(t *testing.T) {
	r := &ConfigAggregateConfigRuleResource{}
	req := resource.MetadataRequest{}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	if resp.TypeName != "tencentcloud_config_aggregate_config_rule" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_config_aggregate_config_rule")
	}
}

// TestConfigAggregateConfigRuleResource_Schema validates required/optional/computed
// flags on every top-level attribute.
func TestConfigAggregateConfigRuleResource_Schema(t *testing.T) {
	r := &ConfigAggregateConfigRuleResource{}
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
		"id":                         {computed: true},
		"identifier":                 {required: true},
		"identifier_type":            {required: true},
		"rule_name":                  {required: true},
		"resource_type":              {required: true},
		"trigger_type":               {required: true},
		"risk_level":                 {required: true},
		"account_group_id":           {required: true},
		"input_parameter":            {optional: true},
		"description":                {optional: true},
		"region_scope":               {optional: true},
		"tags_scope":                 {optional: true},
		"exclude_resource_ids_scope": {optional: true},
		"config_rule_id":             {computed: true},
		"status":                     {computed: true},
		"compliance_result":          {computed: true},
		"create_time":                {computed: true},
		"service_function":           {computed: true},
		"compliance_pack_id":         {computed: true},
		"compliance_pack_name":       {computed: true},
		"account_group_name":         {computed: true},
		"rule_owner_id":              {computed: true},
		"labels":                     {computed: true},
		"regions_scope":              {computed: true},
		"manage_trigger_type":        {computed: true},
		"source_condition":           {computed: true},
		"manage_input_parameter":     {computed: true},
		"config_rule_invoked_time":   {computed: true},
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
		case schema.Int64Attribute:
			if a.Required != want.required {
				t.Errorf("attribute %q Required=%v, want %v", name, a.Required, want.required)
			}
			if a.Optional != want.optional {
				t.Errorf("attribute %q Optional=%v, want %v", name, a.Optional, want.optional)
			}
			if a.Computed != want.computed {
				t.Errorf("attribute %q Computed=%v, want %v", name, a.Computed, want.computed)
			}
		case schema.BoolAttribute:
			if a.Computed != want.computed {
				t.Errorf("attribute %q Computed=%v, want %v", name, a.Computed, want.computed)
			}
		case schema.ListAttribute:
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
			t.Errorf("attribute %q unexpected type %T", name, attr)
		}
	}

	// identifier, identifier_type and account_group_id must force replacement.
	for _, name := range []string{"identifier", "identifier_type", "account_group_id"} {
		cr, ok := resp.Schema.Attributes[name].(schema.StringAttribute)
		if !ok {
			t.Fatalf("%s not a StringAttribute", name)
		}
		if len(cr.PlanModifiers) == 0 {
			t.Errorf("%s should have a RequiresReplace plan modifier", name)
		}
	}
}

// TestConfigAggregateConfigRuleResource_ImplementsInterfaces verifies the
// resource satisfies the framework interfaces advertised in the _ asserts.
func TestConfigAggregateConfigRuleResource_ImplementsInterfaces(t *testing.T) {
	var _ resource.ResourceWithConfigure = &ConfigAggregateConfigRuleResource{}
	var _ resource.ResourceWithImportState = &ConfigAggregateConfigRuleResource{}
}

// TestConfigAggregateConfigRuleResource_Factory validates the registry factory.
func TestConfigAggregateConfigRuleResource_Factory(t *testing.T) {
	got := NewConfigAggregateConfigRuleResource()
	if got == nil {
		t.Fatal("NewConfigAggregateConfigRuleResource returned nil")
	}
	if _, ok := got.(*ConfigAggregateConfigRuleResource); !ok {
		t.Fatalf("factory returned %T, want *ConfigAggregateConfigRuleResource", got)
	}
}

// --- ID split ---

// TestSplitAggregateConfigRuleId validates parsing of the composite id.
func TestSplitAggregateConfigRuleId(t *testing.T) {
	good := "ca-xxxxxxxx" + tccommon.FILED_SP + "cr-yyyyyyyy"
	ag, cr, ok := splitAggregateConfigRuleId(good)
	if !ok || ag != "ca-xxxxxxxx" || cr != "cr-yyyyyyyy" {
		t.Errorf("splitAggregateConfigRuleId(%q) = (%q, %q, %v); want (\"ca-xxxxxxxx\", \"cr-yyyyyyyy\", true)", good, ag, cr, ok)
	}

	bad := []string{
		"",
		"only-group",
		tccommon.FILED_SP + "cr",
		"ca" + tccommon.FILED_SP,
		"a" + tccommon.FILED_SP + "b" + tccommon.FILED_SP + "c",
	}
	for _, in := range bad {
		if _, _, ok := splitAggregateConfigRuleId(in); ok {
			t.Errorf("splitAggregateConfigRuleId(%q) ok=true, want false", in)
		}
	}
}

// --- expand / flatten ---

// buildAggregateConfigRuleTriggerTypesList constructs a types.List of
// trigger_type for tests.
func buildAggregateConfigRuleTriggerTypesList(ctx context.Context, t *testing.T) types.List {
	t.Helper()
	entry := aggregateConfigRuleTriggerTypeModel{
		MessageType:               types.StringValue("ScheduledNotification"),
		MaximumExecutionFrequency: types.StringValue("TwentyFour_Hours"),
	}
	list, diags := types.ListValueFrom(ctx, triggerTypeObjectType(), []aggregateConfigRuleTriggerTypeModel{entry})
	if diags.HasError() {
		t.Fatalf("build trigger_type list failed: %v", diags)
	}
	return list
}

// TestExpandAggregateConfigRuleTriggerTypes verifies the plan list -> SDK
// struct conversion covers all leaf fields.
func TestExpandAggregateConfigRuleTriggerTypes(t *testing.T) {
	ctx := context.Background()
	list := buildAggregateConfigRuleTriggerTypesList(ctx, t)

	entries, diags := expandAggregateConfigRuleTriggerTypes(ctx, list)
	if diags.HasError() {
		t.Fatalf("expandAggregateConfigRuleTriggerTypes returned diagnostics: %v", diags)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 trigger type, got %d", len(entries))
	}
	if entries[0].MessageType == nil || *entries[0].MessageType != "ScheduledNotification" {
		t.Errorf("MessageType = %v, want ScheduledNotification", entries[0].MessageType)
	}
	if entries[0].MaximumExecutionFrequency == nil || *entries[0].MaximumExecutionFrequency != "TwentyFour_Hours" {
		t.Errorf("MaximumExecutionFrequency = %v, want TwentyFour_Hours", entries[0].MaximumExecutionFrequency)
	}
}

// TestExpandAggregateConfigRuleTriggerTypes_Null verifies a null list yields
// no entries.
func TestExpandAggregateConfigRuleTriggerTypes_Null(t *testing.T) {
	entries, diags := expandAggregateConfigRuleTriggerTypes(context.Background(), types.ListNull(triggerTypeObjectType()))
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if entries != nil {
		t.Errorf("expected nil entries for null list, got %v", entries)
	}
}

// TestExpandAggregateConfigRuleInputParameters verifies the plan list -> SDK
// struct conversion covers all leaf fields.
func TestExpandAggregateConfigRuleInputParameters(t *testing.T) {
	ctx := context.Background()
	entry := aggregateConfigRuleInputParameterModel{
		ParameterKey: types.StringValue("MinimumLength"),
		Type:         types.StringValue("Optional"),
		Value:        types.StringValue("8"),
	}
	list, diags := types.ListValueFrom(ctx, aggregateInputParameterObjectType(), []aggregateConfigRuleInputParameterModel{entry})
	if diags.HasError() {
		t.Fatalf("build input_parameter list failed: %v", diags)
	}

	params, expandDiags := expandAggregateConfigRuleInputParameters(ctx, list)
	if expandDiags.HasError() {
		t.Fatalf("expandAggregateConfigRuleInputParameters returned diagnostics: %v", expandDiags)
	}
	if len(params) != 1 {
		t.Fatalf("expected 1 parameter, got %d", len(params))
	}
	if params[0].ParameterKey == nil || *params[0].ParameterKey != "MinimumLength" {
		t.Errorf("ParameterKey = %v, want MinimumLength", params[0].ParameterKey)
	}
	if params[0].Type == nil || *params[0].Type != "Optional" {
		t.Errorf("Type = %v, want Optional", params[0].Type)
	}
	if params[0].Value == nil || *params[0].Value != "8" {
		t.Errorf("Value = %v, want 8", params[0].Value)
	}
}

// TestExpandAggregateConfigRuleTagsScope verifies the plan list -> SDK struct
// conversion covers all leaf fields.
func TestExpandAggregateConfigRuleTagsScope(t *testing.T) {
	ctx := context.Background()
	entry := aggregateConfigRuleTagsScopeModel{
		TagKey:   types.StringValue("env"),
		TagValue: types.StringValue("prod"),
	}
	list, diags := types.ListValueFrom(ctx, tagsScopeObjectType(), []aggregateConfigRuleTagsScopeModel{entry})
	if diags.HasError() {
		t.Fatalf("build tags_scope list failed: %v", diags)
	}

	tags, expandDiags := expandAggregateConfigRuleTagsScope(ctx, list)
	if expandDiags.HasError() {
		t.Fatalf("expandAggregateConfigRuleTagsScope returned diagnostics: %v", expandDiags)
	}
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(tags))
	}
	if tags[0].TagKey == nil || *tags[0].TagKey != "env" {
		t.Errorf("TagKey = %v, want env", tags[0].TagKey)
	}
	if tags[0].TagValue == nil || *tags[0].TagValue != "prod" {
		t.Errorf("TagValue = %v, want prod", tags[0].TagValue)
	}
}

// TestFlattenAggregateConfigRuleStrings verifies the SDK slice -> list
// conversion and the skip of nil elements.
func TestFlattenAggregateConfigRuleStrings(t *testing.T) {
	ctx := context.Background()
	list, diags := flattenAggregateConfigRuleStrings(ctx, []*string{helper.String("QCS::CAM::User"), nil}, types.ListNull(types.StringType))
	if diags.HasError() {
		t.Fatalf("flattenAggregateConfigRuleStrings returned diagnostics: %v", diags)
	}
	if list.IsNull() || len(list.Elements()) != 1 {
		t.Errorf("expected 1 element, got %v", list)
	}

	// An empty cloud list must not clear a configured list.
	prev, _ := types.ListValueFrom(ctx, types.StringType, []string{"keep-me"})
	kept, diags := flattenAggregateConfigRuleStrings(ctx, nil, prev)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if len(kept.Elements()) != 1 {
		t.Errorf("expected the previous value to be kept, got %v", kept)
	}
}

// TestFlattenAggregateConfigRuleTriggerTypes verifies the SDK slice -> list
// conversion reads back leaf fields.
func TestFlattenAggregateConfigRuleTriggerTypes(t *testing.T) {
	ctx := context.Background()
	prevList := buildAggregateConfigRuleTriggerTypesList(ctx, t)

	entries := []*configv20220802.TriggerType{
		{
			MessageType:               helper.String("ConfigurationItemChangeNotification"),
			MaximumExecutionFrequency: nil,
		},
	}

	list, diags := flattenAggregateConfigRuleTriggerTypes(ctx, entries, prevList)
	if diags.HasError() {
		t.Fatalf("flattenAggregateConfigRuleTriggerTypes returned diagnostics: %v", diags)
	}
	if list.IsNull() {
		t.Fatal("flattened list is null")
	}

	var out []aggregateConfigRuleTriggerTypeModel
	if diags := list.ElementsAs(ctx, &out, false); diags.HasError() {
		t.Fatalf("elements as failed: %v", diags)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].MessageType.ValueString() != "ConfigurationItemChangeNotification" {
		t.Errorf("MessageType = %q, want ConfigurationItemChangeNotification", out[0].MessageType.ValueString())
	}
	// maximum_execution_frequency is not echoed by the cloud here; it must be
	// carried over from the previous state.
	if out[0].MaximumExecutionFrequency.ValueString() != "TwentyFour_Hours" {
		t.Errorf("MaximumExecutionFrequency = %q, want preserved value", out[0].MaximumExecutionFrequency.ValueString())
	}
}

// TestFlattenAggregateConfigRuleSourceConditions verifies the read-only
// source_condition list conversion.
func TestFlattenAggregateConfigRuleSourceConditions(t *testing.T) {
	ctx := context.Background()
	conditions := []*configv20220802.SourceConditionForManage{
		{
			EmptyAs:      helper.String("NOT_APPLICABLE"),
			SelectPath:   helper.String("passwordPolicy.MinimumLength"),
			Operator:     helper.String("gte"),
			Required:     helper.Bool(true),
			DesiredValue: helper.String("8"),
		},
		nil,
	}

	list, diags := flattenAggregateConfigRuleSourceConditions(ctx, conditions)
	if diags.HasError() {
		t.Fatalf("flattenAggregateConfigRuleSourceConditions returned diagnostics: %v", diags)
	}
	if list.IsNull() || len(list.Elements()) != 1 {
		t.Fatalf("expected 1 condition, got %v", list)
	}

	var out []aggregateConfigRuleSourceConditionModel
	if diags := list.ElementsAs(ctx, &out, false); diags.HasError() {
		t.Fatalf("elements as failed: %v", diags)
	}
	if out[0].EmptyAs.ValueString() != "NOT_APPLICABLE" {
		t.Errorf("EmptyAs = %q, want NOT_APPLICABLE", out[0].EmptyAs.ValueString())
	}
	if out[0].SelectPath.ValueString() != "passwordPolicy.MinimumLength" {
		t.Errorf("SelectPath = %q", out[0].SelectPath.ValueString())
	}
	if out[0].Operator.ValueString() != "gte" {
		t.Errorf("Operator = %q", out[0].Operator.ValueString())
	}
	if !out[0].Required.ValueBool() {
		t.Errorf("Required = %v, want true", out[0].Required.ValueBool())
	}
	if out[0].DesiredValue.ValueString() != "8" {
		t.Errorf("DesiredValue = %q", out[0].DesiredValue.ValueString())
	}
}

// TestFlattenAggregateConfigRuleManageInputParameters verifies the read-only
// manage_input_parameter list conversion.
func TestFlattenAggregateConfigRuleManageInputParameters(t *testing.T) {
	ctx := context.Background()
	params := []*configv20220802.InputParameterForManage{
		{
			ValueType:    helper.String("Integer"),
			ParameterKey: helper.String("MinimumLength"),
			Type:         helper.String("Optional"),
			DefaultValue: helper.String("8"),
			Description:  helper.String("minimum length"),
		},
	}

	list, diags := flattenAggregateConfigRuleManageInputParameters(ctx, params)
	if diags.HasError() {
		t.Fatalf("flattenAggregateConfigRuleManageInputParameters returned diagnostics: %v", diags)
	}
	if list.IsNull() || len(list.Elements()) != 1 {
		t.Fatalf("expected 1 parameter, got %v", list)
	}

	var out []aggregateConfigRuleManageInputParameterModel
	if diags := list.ElementsAs(ctx, &out, false); diags.HasError() {
		t.Fatalf("elements as failed: %v", diags)
	}
	if out[0].ValueType.ValueString() != "Integer" {
		t.Errorf("ValueType = %q", out[0].ValueType.ValueString())
	}
	if out[0].ParameterKey.ValueString() != "MinimumLength" {
		t.Errorf("ParameterKey = %q", out[0].ParameterKey.ValueString())
	}
	if out[0].DefaultValue.ValueString() != "8" {
		t.Errorf("DefaultValue = %q", out[0].DefaultValue.ValueString())
	}
}

// TestRefreshAggregateConfigRuleModel verifies cloud-owned top-level fields are
// written into the model and identity parts are reconstructed.
func TestRefreshAggregateConfigRuleModel(t *testing.T) {
	ctx := context.Background()
	m := &ConfigAggregateConfigRuleResourceModel{
		RuleName:     types.StringValue("tf-example-rule"),
		ResourceType: types.ListNull(types.StringType),
		TriggerType:  types.ListNull(triggerTypeObjectType()),
		RiskLevel:    types.Int64Value(1),
	}

	riskLevel := int64(2)
	ruleName := "tf-example-rule"
	item := &configv20220802.ConfigRule{
		Identifier:     helper.String("cam-password-policy-check"),
		IdentifierType: helper.String("SYSTEM"),
		RuleName:       &ruleName,
		RiskLevel:      &riskLevel,
		Status:         helper.String("ACTIVE"),
		CreateTime:     helper.String("2024-01-01 00:00:00"),
		ConfigRuleId:   helper.String("cr-zzz"),
		AccountGroupId: helper.String("ca-aaa"),
	}

	if diags := refreshAggregateConfigRuleModel(ctx, m, "ca-aaa", "cr-zzz", item); diags.HasError() {
		t.Fatalf("refreshAggregateConfigRuleModel returned diagnostics: %v", diags)
	}

	if m.ConfigRuleId.ValueString() != "cr-zzz" {
		t.Errorf("ConfigRuleId = %q, want cr-zzz", m.ConfigRuleId.ValueString())
	}
	wantId := "ca-aaa" + tccommon.FILED_SP + "cr-zzz"
	if m.Id.ValueString() != wantId {
		t.Errorf("Id = %q, want %q", m.Id.ValueString(), wantId)
	}
	if m.Identifier.ValueString() != "cam-password-policy-check" {
		t.Errorf("Identifier = %q", m.Identifier.ValueString())
	}
	if m.IdentifierType.ValueString() != "SYSTEM" {
		t.Errorf("IdentifierType = %q", m.IdentifierType.ValueString())
	}
	if m.Status.ValueString() != "ACTIVE" {
		t.Errorf("Status = %q", m.Status.ValueString())
	}
	if m.CreateTime.ValueString() != "2024-01-01 00:00:00" {
		t.Errorf("CreateTime = %q", m.CreateTime.ValueString())
	}
	if m.RiskLevel.ValueInt64() != 2 {
		t.Errorf("RiskLevel = %d, want 2", m.RiskLevel.ValueInt64())
	}
	if m.ResourceType.IsNull() || len(m.ResourceType.Elements()) != 0 {
		t.Errorf("ResourceType should be a known empty list, got %v", m.ResourceType)
	}
}

// TestRefreshAggregateConfigRuleModel_Malformed verifies a response missing
// required attributes is reported instead of producing an invalid state.
func TestRefreshAggregateConfigRuleModel_Malformed(t *testing.T) {
	ctx := context.Background()
	m := &ConfigAggregateConfigRuleResourceModel{
		ResourceType: types.ListNull(types.StringType),
		TriggerType:  types.ListNull(triggerTypeObjectType()),
	}
	item := &configv20220802.ConfigRule{
		ConfigRuleId: helper.String("cr-zzz"),
		// no rule_name / risk_level
	}

	if diags := refreshAggregateConfigRuleModel(ctx, m, "ca-aaa", "cr-zzz", item); !diags.HasError() {
		t.Fatal("expected an error diagnostic for a response without rule_name or risk_level")
	}
}
