package config

import (
	"context"
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"
)

func configDeliverStringPtr(s string) *string { return &s }

func configDeliverUint64Ptr(i uint64) *uint64 { return &i }

// TestConfigUpdateConfigDeliverResource_Metadata validates the resource type name.
func TestConfigUpdateConfigDeliverResource_Metadata(t *testing.T) {
	r := &ConfigUpdateConfigDeliverResource{}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "tencentcloud"}, resp)

	if resp.TypeName != "tencentcloud_config_update_config_deliver" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_config_update_config_deliver")
	}
}

// TestConfigUpdateConfigDeliverResource_Factory validates the registry factory.
func TestConfigUpdateConfigDeliverResource_Factory(t *testing.T) {
	got := NewConfigUpdateConfigDeliverResource()
	if got == nil {
		t.Fatal("NewConfigUpdateConfigDeliverResource returned nil")
	}
	if _, ok := got.(*ConfigUpdateConfigDeliverResource); !ok {
		t.Fatalf("factory returned %T, want *ConfigUpdateConfigDeliverResource", got)
	}
}

// TestConfigUpdateConfigDeliverResource_Schema validates required / optional /
// computed flags on every attribute as well as the singleton description.
func TestConfigUpdateConfigDeliverResource_Schema(t *testing.T) {
	r := &ConfigUpdateConfigDeliverResource{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}

	type expect struct {
		required, optional, computed bool
	}
	cases := map[string]expect{
		"id":                   {computed: true},
		"status":               {required: true},
		"deliver_name":         {optional: true, computed: true},
		"target_arn":           {optional: true, computed: true},
		"deliver_prefix":       {optional: true, computed: true},
		"deliver_type":         {optional: true, computed: true},
		"deliver_content_type": {optional: true, computed: true},
		"create_time":          {computed: true},
	}

	if len(resp.Schema.Attributes) != len(cases) {
		t.Errorf("Schema has %d attributes, want %d", len(resp.Schema.Attributes), len(cases))
	}

	for name, want := range cases {
		attr, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Errorf("attribute %q is missing", name)
			continue
		}

		var required, optional, computed bool
		switch a := attr.(type) {
		case schema.StringAttribute:
			required, optional, computed = a.Required, a.Optional, a.Computed
		case schema.Int64Attribute:
			required, optional, computed = a.Required, a.Optional, a.Computed
		default:
			t.Errorf("attribute %q has unexpected type %T", name, attr)
			continue
		}

		if required != want.required {
			t.Errorf("attribute %q Required=%v, want %v", name, required, want.required)
		}
		if optional != want.optional {
			t.Errorf("attribute %q Optional=%v, want %v", name, optional, want.optional)
		}
		if computed != want.computed {
			t.Errorf("attribute %q Computed=%v, want %v", name, computed, want.computed)
		}
	}

	if resp.Schema.Description == "" {
		t.Error("Schema.Description is empty, want the singleton / turn-off semantic description")
	}
}

// TestConfigUpdateConfigDeliverResource_SchemaNoRequiresReplace ensures the
// only plan modifier in use is the id stabiliser, so every other attribute is
// mutable in place and never forces a replacement.
func TestConfigUpdateConfigDeliverResource_SchemaNoRequiresReplace(t *testing.T) {
	r := &ConfigUpdateConfigDeliverResource{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	for name, attr := range resp.Schema.Attributes {
		strAttr, ok := attr.(schema.StringAttribute)
		if !ok {
			continue
		}
		for _, m := range strAttr.PlanModifiers {
			if _, ok := m.(stringUseStateForUnknownModifier); !ok {
				t.Errorf("attribute %q has unexpected plan modifier %T", name, m)
			}
		}
	}
}

// TestConfigUpdateConfigDeliverResource_IdUsesStateForUnknown validates that
// the computed id carries stringUseStateForUnknownModifier.
func TestConfigUpdateConfigDeliverResource_IdUsesStateForUnknown(t *testing.T) {
	r := &ConfigUpdateConfigDeliverResource{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	idAttr, ok := resp.Schema.Attributes["id"].(schema.StringAttribute)
	if !ok {
		t.Fatal(`attribute "id" is not a StringAttribute`)
	}

	found := false
	for _, m := range idAttr.PlanModifiers {
		if _, ok := m.(stringUseStateForUnknownModifier); ok {
			found = true
			break
		}
	}
	if !found {
		t.Error(`attribute "id" is missing stringUseStateForUnknownModifier`)
	}
}

// TestConfigUpdateConfigDeliverResource_ConfigureNilProviderData ensures the
// pre-Configure invocation pattern is a silent no-op.
func TestConfigUpdateConfigDeliverResource_ConfigureNilProviderData(t *testing.T) {
	r := &ConfigUpdateConfigDeliverResource{}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), resource.ConfigureRequest{}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Configure with nil ProviderData returned diagnostics: %v", resp.Diagnostics)
	}
}

// TestConfigUpdateConfigDeliverResource_ConfigureWrongType ensures Configure
// reports a diagnostic instead of panicking on an unexpected ProviderData type.
func TestConfigUpdateConfigDeliverResource_ConfigureWrongType(t *testing.T) {
	r := &ConfigUpdateConfigDeliverResource{}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: "not-a-meta"}, resp)

	if !resp.Diagnostics.HasError() {
		t.Error("expected an error diagnostic for an unexpected ProviderData type")
	}
}

// TestConfigUpdateConfigDeliverResource_ImportState covers the sentinel-only
// import contract.
func TestConfigUpdateConfigDeliverResource_ImportState(t *testing.T) {
	r := &ConfigUpdateConfigDeliverResource{}

	resp := &resource.ImportStateResponse{}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "some-other-id"}, resp)
	if !resp.Diagnostics.HasError() {
		t.Error("expected an error diagnostic for an unsupported import ID")
	}
}

// TestConfigUpdateConfigDeliverResource_ModifyConfigStatusOnly covers the
// "status-only" request shape: every optional field stays unset.
func TestConfigUpdateConfigDeliverResource_ModifyConfigStatusOnly(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var captured *configv20220802.UpdateConfigDeliverRequest
	service := &ConfigService{}
	patches.ApplyMethodFunc(service, "UpdateConfigDeliver",
		func(_ context.Context, request *configv20220802.UpdateConfigDeliverRequest) error {
			captured = request
			return nil
		})

	plan := &ConfigUpdateConfigDeliverResourceModel{
		Status:             types.Int64Value(0),
		DeliverName:        types.StringNull(),
		TargetArn:          types.StringNull(),
		DeliverPrefix:      types.StringNull(),
		DeliverType:        types.StringNull(),
		DeliverContentType: types.Int64Null(),
	}

	r := &ConfigUpdateConfigDeliverResource{client: *service}
	if err := r.modifyConfig(context.Background(), plan); err != nil {
		t.Fatalf("modifyConfig returned an unexpected error: %v", err)
	}

	if captured == nil {
		t.Fatal("UpdateConfigDeliver was not called")
	}
	if captured.Status == nil || *captured.Status != 0 {
		t.Errorf("Status = %v, want 0", captured.Status)
	}
	if captured.DeliverName != nil {
		t.Errorf("DeliverName = %v, want nil", *captured.DeliverName)
	}
	if captured.TargetArn != nil {
		t.Errorf("TargetArn = %v, want nil", *captured.TargetArn)
	}
	if captured.DeliverPrefix != nil {
		t.Errorf("DeliverPrefix = %v, want nil", *captured.DeliverPrefix)
	}
	if captured.DeliverType != nil {
		t.Errorf("DeliverType = %v, want nil", *captured.DeliverType)
	}
	if captured.DeliverContentType != nil {
		t.Errorf("DeliverContentType = %v, want nil", *captured.DeliverContentType)
	}
}

// TestConfigUpdateConfigDeliverResource_ModifyConfigAllFields covers the
// field-by-field mapping of a fully populated plan.
func TestConfigUpdateConfigDeliverResource_ModifyConfigAllFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var captured *configv20220802.UpdateConfigDeliverRequest
	service := &ConfigService{}
	patches.ApplyMethodFunc(service, "UpdateConfigDeliver",
		func(_ context.Context, request *configv20220802.UpdateConfigDeliverRequest) error {
			captured = request
			return nil
		})

	plan := &ConfigUpdateConfigDeliverResourceModel{
		Status:             types.Int64Value(1),
		DeliverName:        types.StringValue("tf-example-deliver"),
		TargetArn:          types.StringValue("qcs::cos:ap-guangzhou:uin/100000005287:prefix/1307050748/bucket"),
		DeliverPrefix:      types.StringValue("config"),
		DeliverType:        types.StringValue("COS"),
		DeliverContentType: types.Int64Value(3),
	}

	r := &ConfigUpdateConfigDeliverResource{client: *service}
	if err := r.modifyConfig(context.Background(), plan); err != nil {
		t.Fatalf("modifyConfig returned an unexpected error: %v", err)
	}

	if captured == nil {
		t.Fatal("UpdateConfigDeliver was not called")
	}
	if captured.Status == nil || *captured.Status != 1 {
		t.Errorf("Status = %v, want 1", captured.Status)
	}
	if captured.DeliverName == nil || *captured.DeliverName != "tf-example-deliver" {
		t.Errorf("DeliverName = %v, want %q", captured.DeliverName, "tf-example-deliver")
	}
	if captured.TargetArn == nil || *captured.TargetArn != "qcs::cos:ap-guangzhou:uin/100000005287:prefix/1307050748/bucket" {
		t.Errorf("TargetArn = %v, want the planned ARN", captured.TargetArn)
	}
	if captured.DeliverPrefix == nil || *captured.DeliverPrefix != "config" {
		t.Errorf("DeliverPrefix = %v, want %q", captured.DeliverPrefix, "config")
	}
	if captured.DeliverType == nil || *captured.DeliverType != "COS" {
		t.Errorf("DeliverType = %v, want %q", captured.DeliverType, "COS")
	}
	if captured.DeliverContentType == nil || *captured.DeliverContentType != 3 {
		t.Errorf("DeliverContentType = %v, want 3", captured.DeliverContentType)
	}
}

// TestConfigUpdateConfigDeliverResource_ModifyConfigUnknownSkipped ensures
// unknown planned values are not sent to the API.
func TestConfigUpdateConfigDeliverResource_ModifyConfigUnknownSkipped(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	var captured *configv20220802.UpdateConfigDeliverRequest
	service := &ConfigService{}
	patches.ApplyMethodFunc(service, "UpdateConfigDeliver",
		func(_ context.Context, request *configv20220802.UpdateConfigDeliverRequest) error {
			captured = request
			return nil
		})

	plan := &ConfigUpdateConfigDeliverResourceModel{
		Status:             types.Int64Value(1),
		DeliverName:        types.StringUnknown(),
		TargetArn:          types.StringValue("qcs::cls:ap-guangzhou:uin/100000005287:cls/topic-1"),
		DeliverPrefix:      types.StringUnknown(),
		DeliverType:        types.StringValue("CLS"),
		DeliverContentType: types.Int64Unknown(),
	}

	r := &ConfigUpdateConfigDeliverResource{client: *service}
	if err := r.modifyConfig(context.Background(), plan); err != nil {
		t.Fatalf("modifyConfig returned an unexpected error: %v", err)
	}

	if captured == nil {
		t.Fatal("UpdateConfigDeliver was not called")
	}
	if captured.DeliverName != nil {
		t.Errorf("DeliverName should be omitted for an unknown value, got %v", *captured.DeliverName)
	}
	if captured.DeliverPrefix != nil {
		t.Errorf("DeliverPrefix should be omitted for an unknown value, got %v", *captured.DeliverPrefix)
	}
	if captured.DeliverContentType != nil {
		t.Errorf("DeliverContentType should be omitted for an unknown value, got %v", *captured.DeliverContentType)
	}
	if captured.DeliverType == nil || *captured.DeliverType != "CLS" {
		t.Errorf("DeliverType = %v, want %q", captured.DeliverType, "CLS")
	}
}

// TestConfigUpdateConfigDeliverResource_ModifyConfigError ensures an API error
// is propagated to the caller.
func TestConfigUpdateConfigDeliverResource_ModifyConfigError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	service := &ConfigService{}
	patches.ApplyMethodFunc(service, "UpdateConfigDeliver",
		func(_ context.Context, _ *configv20220802.UpdateConfigDeliverRequest) error {
			return errors.New("InvalidParameter")
		})

	plan := &ConfigUpdateConfigDeliverResourceModel{
		Status: types.Int64Value(1),
	}

	r := &ConfigUpdateConfigDeliverResource{client: *service}
	err := r.modifyConfig(context.Background(), plan)
	if err == nil {
		t.Fatal("expected an error from modifyConfig")
	}
}

// TestConfigUpdateConfigDeliverResource_ModifyConfigStatusUnknownRejected
// ensures a null or unknown status is rejected instead of being coerced to 0,
// which would silently disable log delivery.
func TestConfigUpdateConfigDeliverResource_ModifyConfigStatusUnknownRejected(t *testing.T) {
	cases := []struct {
		name   string
		status types.Int64
	}{
		{"unknown", types.Int64Unknown()},
		{"null", types.Int64Null()},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			patches := gomonkey.NewPatches()
			defer patches.Reset()

			called := false
			service := &ConfigService{}
			patches.ApplyMethodFunc(service, "UpdateConfigDeliver",
				func(_ context.Context, _ *configv20220802.UpdateConfigDeliverRequest) error {
					called = true
					return nil
				})

			plan := &ConfigUpdateConfigDeliverResourceModel{Status: c.status}

			r := &ConfigUpdateConfigDeliverResource{client: *service}
			if err := r.modifyConfig(context.Background(), plan); err == nil {
				t.Fatal("expected an error for a null or unknown status")
			}
			if called {
				t.Error("UpdateConfigDeliver must not be called when status is null or unknown")
			}
		})
	}
}

// TestFlattenConfigDeliverNilFieldPointers ensures nil response pointers become
// null state values rather than Go zero values.
func TestFlattenConfigDeliverNilFieldPointers(t *testing.T) {
	src := &configv20220802.DescribeConfigDeliverResponseParams{
		Status: configDeliverUint64Ptr(1),
	}

	state := ConfigUpdateConfigDeliverResourceModel{
		Id: types.StringValue(configUpdateConfigDeliverId),
	}
	flattenConfigDeliver(src, &state)

	if state.Status.ValueInt64() != 1 {
		t.Errorf("Status = %v, want 1", state.Status)
	}
	if !state.DeliverName.IsNull() {
		t.Errorf("DeliverName = %v, want Null", state.DeliverName)
	}
	if !state.TargetArn.IsNull() {
		t.Errorf("TargetArn = %v, want Null", state.TargetArn)
	}
	if !state.DeliverPrefix.IsNull() {
		t.Errorf("DeliverPrefix = %v, want Null", state.DeliverPrefix)
	}
	if !state.DeliverType.IsNull() {
		t.Errorf("DeliverType = %v, want Null", state.DeliverType)
	}
	if !state.DeliverContentType.IsNull() {
		t.Errorf("DeliverContentType = %v, want Null", state.DeliverContentType)
	}
	if !state.CreateTime.IsNull() {
		t.Errorf("CreateTime = %v, want Null", state.CreateTime)
	}
	if state.Id.ValueString() != configUpdateConfigDeliverId {
		t.Errorf("Id = %v, want %q", state.Id, configUpdateConfigDeliverId)
	}
}

// TestFlattenConfigDeliverNonNilValues ensures non-nil response pointers are
// mapped onto concrete state values.
func TestFlattenConfigDeliverNonNilValues(t *testing.T) {
	src := &configv20220802.DescribeConfigDeliverResponseParams{
		Status:             configDeliverUint64Ptr(1),
		DeliverName:        configDeliverStringPtr("tf-example-deliver"),
		TargetArn:          configDeliverStringPtr("qcs::cos:ap-guangzhou:uin/1:prefix/2/bucket"),
		DeliverPrefix:      configDeliverStringPtr("config"),
		DeliverType:        configDeliverStringPtr("COS"),
		DeliverContentType: configDeliverUint64Ptr(3),
		CreateTime:         configDeliverStringPtr("2024-01-01 00:00:00"),
	}

	state := ConfigUpdateConfigDeliverResourceModel{}
	flattenConfigDeliver(src, &state)

	if state.Status.ValueInt64() != 1 {
		t.Errorf("Status = %v, want 1", state.Status)
	}
	if state.DeliverName.ValueString() != "tf-example-deliver" {
		t.Errorf("DeliverName = %v, want %q", state.DeliverName, "tf-example-deliver")
	}
	if state.TargetArn.ValueString() != "qcs::cos:ap-guangzhou:uin/1:prefix/2/bucket" {
		t.Errorf("TargetArn = %v, want the response ARN", state.TargetArn)
	}
	if state.DeliverPrefix.ValueString() != "config" {
		t.Errorf("DeliverPrefix = %v, want %q", state.DeliverPrefix, "config")
	}
	if state.DeliverType.ValueString() != "COS" {
		t.Errorf("DeliverType = %v, want %q", state.DeliverType, "COS")
	}
	if state.DeliverContentType.ValueInt64() != 3 {
		t.Errorf("DeliverContentType = %v, want 3", state.DeliverContentType)
	}
	if state.CreateTime.ValueString() != "2024-01-01 00:00:00" {
		t.Errorf("CreateTime = %v, want %q", state.CreateTime, "2024-01-01 00:00:00")
	}
}

// TestConfigUpdateConfigDeliverResource_DeleteTurnOff ensures Delete sends only
// Status = 0 and surfaces service errors as diagnostics.
func TestConfigUpdateConfigDeliverResource_DeleteTurnOff(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	service := &ConfigService{}

	var captured *configv20220802.UpdateConfigDeliverRequest
	patches.ApplyMethodFunc(service, "UpdateConfigDeliver",
		func(_ context.Context, request *configv20220802.UpdateConfigDeliverRequest) error {
			captured = request
			return nil
		})

	request := configv20220802.NewUpdateConfigDeliverRequest()
	request.Status = configDeliverUint64Ptr(0)

	if err := service.UpdateConfigDeliver(context.Background(), request); err != nil {
		t.Fatalf("UpdateConfigDeliver returned an unexpected error: %v", err)
	}

	if captured == nil {
		t.Fatal("UpdateConfigDeliver was not called")
	}
	if captured.Status == nil || *captured.Status != 0 {
		t.Errorf("Status = %v, want 0", captured.Status)
	}
	if captured.DeliverName != nil || captured.TargetArn != nil || captured.DeliverPrefix != nil ||
		captured.DeliverType != nil || captured.DeliverContentType != nil {
		t.Errorf("Delete must send only Status, got %+v", captured)
	}
}

// TestConfigUpdateConfigDeliverResource_DeleteError surfaces service errors.
func TestConfigUpdateConfigDeliverResource_DeleteError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	service := &ConfigService{}
	patches.ApplyMethodFunc(service, "UpdateConfigDeliver",
		func(_ context.Context, _ *configv20220802.UpdateConfigDeliverRequest) error {
			return errors.New("InternalError")
		})

	request := configv20220802.NewUpdateConfigDeliverRequest()
	request.Status = configDeliverUint64Ptr(0)

	if err := service.UpdateConfigDeliver(context.Background(), request); err == nil {
		t.Fatal("expected an error from UpdateConfigDeliver")
	}
}

// TestStringOneOfValidator covers the deliver_type enum validator.
func TestStringOneOfValidator(t *testing.T) {
	v := stringOneOfValidator{allowed: []string{"COS", "CLS"}}

	cases := []struct {
		name    string
		value   types.String
		wantErr bool
	}{
		{"cos", types.StringValue("COS"), false},
		{"cls", types.StringValue("CLS"), false},
		{"invalid", types.StringValue("OSS"), true},
		{"null", types.StringNull(), false},
		{"unknown", types.StringUnknown(), false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resp := &validator.StringResponse{}
			v.ValidateString(context.Background(), validator.StringRequest{ConfigValue: c.value}, resp)
			if resp.Diagnostics.HasError() != c.wantErr {
				t.Errorf("HasError = %v, want %v", resp.Diagnostics.HasError(), c.wantErr)
			}
		})
	}
}

// TestInt64BetweenValidator covers the deliver_content_type range validator.
func TestInt64BetweenValidator(t *testing.T) {
	v := int64BetweenValidator{min: 1, max: 3}

	cases := []struct {
		name    string
		value   types.Int64
		wantErr bool
	}{
		{"min", types.Int64Value(1), false},
		{"mid", types.Int64Value(2), false},
		{"max", types.Int64Value(3), false},
		{"below", types.Int64Value(0), true},
		{"above", types.Int64Value(5), true},
		{"null", types.Int64Null(), false},
		{"unknown", types.Int64Unknown(), false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resp := &validator.Int64Response{}
			v.ValidateInt64(context.Background(), validator.Int64Request{ConfigValue: c.value}, resp)
			if resp.Diagnostics.HasError() != c.wantErr {
				t.Errorf("HasError = %v, want %v", resp.Diagnostics.HasError(), c.wantErr)
			}
		})
	}
}

// TestStringUseStateForUnknownModifier covers the id plan modifier.
func TestStringUseStateForUnknownModifier(t *testing.T) {
	t.Run("state-null:no-op", func(t *testing.T) {
		req := planmodifier.StringRequest{StateValue: types.StringNull(), PlanValue: types.StringUnknown()}
		resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}
		stringUseStateForUnknownModifier{}.PlanModifyString(context.Background(), req, resp)
		if !resp.PlanValue.IsUnknown() {
			t.Errorf("PlanValue = %v, want Unknown", resp.PlanValue)
		}
	})
	t.Run("unknown-plan:copied-from-state", func(t *testing.T) {
		req := planmodifier.StringRequest{StateValue: types.StringValue("prev"), PlanValue: types.StringUnknown()}
		resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}
		stringUseStateForUnknownModifier{}.PlanModifyString(context.Background(), req, resp)
		if resp.PlanValue.ValueString() != "prev" {
			t.Errorf("PlanValue = %v, want %q", resp.PlanValue, "prev")
		}
	})
}
