package teo

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

const teoZone1110_01TestZoneId = "zone-2noz78a8ev6k"

func teoZone1110_01Ptr(s string) *string { return &s }

// teoZone1110_01NewTestResource returns a configured resource. The shared client
// is a bare TencentCloudClient; the cloud calls themselves are patched per test.
func teoZone1110_01NewTestResource(t *testing.T) resource.Resource {
	t.Helper()

	res := NewTeoZone1110_01Resource()
	configured, ok := res.(resource.ResourceWithConfigure)
	if !ok {
		t.Fatal("resource does not implement resource.ResourceWithConfigure")
	}
	configured.Configure(
		context.Background(),
		resource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: &connectivity.TencentCloudClient{}}},
		&resource.ConfigureResponse{},
	)
	return res
}

// teoZone1110_01Schema returns the resource schema.
func teoZone1110_01Schema(t *testing.T) schema.Schema {
	t.Helper()

	res := NewTeoZone1110_01Resource()
	resp := &resource.SchemaResponse{}
	res.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics from Schema: %v", resp.Diagnostics)
	}
	return resp.Schema
}

// TestTeoZone1110_01Metadata verifies the registered type name.
func TestTeoZone1110_01Metadata(t *testing.T) {
	res := NewTeoZone1110_01Resource()

	resp := &resource.MetadataResponse{}
	res.Metadata(context.Background(), resource.MetadataRequest{}, resp)
	if resp.TypeName != "tencentcloud_teo_zone_1110_01" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_teo_zone_1110_01")
	}
}

// TestTeoZone1110_01SchemaShape verifies the schema attribute flags.
func TestTeoZone1110_01SchemaShape(t *testing.T) {
	zoneSchema := teoZone1110_01Schema(t)

	expectedStrings := map[string]struct{ required, optional, computed bool }{
		"id":              {computed: true},
		"zone_id":         {computed: true},
		"type":            {optional: true, computed: true},
		"zone_name":       {optional: true},
		"area":            {optional: true, computed: true},
		"plan_id":         {optional: true},
		"alias_zone_name": {optional: true, computed: true},
	}
	for name, want := range expectedStrings {
		attribute, ok := zoneSchema.Attributes[name].(schema.StringAttribute)
		if !ok {
			t.Errorf("attribute %q is missing or is not a string attribute", name)
			continue
		}
		if attribute.Required != want.required {
			t.Errorf("attribute %q Required = %v, want %v", name, attribute.Required, want.required)
		}
		if attribute.Optional != want.optional {
			t.Errorf("attribute %q Optional = %v, want %v", name, attribute.Optional, want.optional)
		}
		if attribute.Computed != want.computed {
			t.Errorf("attribute %q Computed = %v, want %v", name, attribute.Computed, want.computed)
		}
	}

	for _, name := range []string{"allow_duplicates", "jump_start"} {
		attribute, ok := zoneSchema.Attributes[name].(schema.BoolAttribute)
		if !ok {
			t.Errorf("attribute %q is missing or is not a bool attribute", name)
			continue
		}
		if !attribute.Optional || attribute.Required || attribute.Computed {
			t.Errorf("attribute %q must be optional only", name)
		}
		if len(attribute.PlanModifiers) != 1 {
			t.Errorf("attribute %q must carry exactly one plan modifier", name)
		}
	}

	tagsAttribute, ok := zoneSchema.Attributes["tags"].(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("attribute tags is missing or is not a list nested attribute")
	}
	if !tagsAttribute.Optional || !tagsAttribute.Computed {
		t.Errorf("attribute tags must be optional and computed")
	}
	for _, name := range []string{"tag_key", "tag_value"} {
		nested, ok := tagsAttribute.NestedObject.Attributes[name].(schema.StringAttribute)
		if !ok {
			t.Errorf("nested attribute tags.%s is missing", name)
			continue
		}
		if !nested.Required {
			t.Errorf("nested attribute tags.%s must be required", name)
		}
	}

	ownership, ok := zoneSchema.Attributes["ownership_verification"].(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("attribute ownership_verification is missing or is not a list nested attribute")
	}
	if !ownership.Computed || ownership.Optional || ownership.Required {
		t.Errorf("attribute ownership_verification must be computed only")
	}
	for _, name := range []string{"dns_verification", "file_verification", "ns_verification"} {
		if _, ok := ownership.NestedObject.Attributes[name]; !ok {
			t.Errorf("nested attribute ownership_verification.%s is missing", name)
		}
	}

	dns, ok := ownership.NestedObject.Attributes["dns_verification"].(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("ownership_verification.dns_verification is not a list nested attribute")
	}
	for _, name := range []string{"subdomain", "record_type", "record_value"} {
		if _, ok := dns.NestedObject.Attributes[name]; !ok {
			t.Errorf("nested attribute ownership_verification.dns_verification.%s is missing", name)
		}
	}

	file, ok := ownership.NestedObject.Attributes["file_verification"].(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("ownership_verification.file_verification is not a list nested attribute")
	}
	for _, name := range []string{"path", "content"} {
		if _, ok := file.NestedObject.Attributes[name]; !ok {
			t.Errorf("nested attribute ownership_verification.file_verification.%s is missing", name)
		}
	}

	ns, ok := ownership.NestedObject.Attributes["ns_verification"].(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("ownership_verification.ns_verification is not a list nested attribute")
	}
	nameServers, ok := ns.NestedObject.Attributes["name_servers"].(schema.SetAttribute)
	if !ok {
		t.Fatalf("ownership_verification.ns_verification.name_servers is not a set attribute")
	}
	if !nameServers.Computed {
		t.Errorf("ownership_verification.ns_verification.name_servers must be computed")
	}
}

// TestTeoZone1110_01BuildCreateZoneRequest verifies every configured field lands
// in the CreateZone request.
func TestTeoZone1110_01BuildCreateZoneRequest(t *testing.T) {
	tags, diags := types.ListValue(teoZone1110_01TagObjectType, []attr.Value{
		teoZone1110_01Object(t, "env", "test"),
		teoZone1110_01Object(t, "owner", "iac"),
	})
	if diags.HasError() {
		t.Fatalf("failed to build tags: %v", diags)
	}

	plan := &TeoZone111001ResourceModel{
		Type:            types.StringValue("partial"),
		ZoneName:        types.StringValue("example.com"),
		Area:            types.StringValue("mainland"),
		PlanId:          types.StringValue("edgeone-2wdo315m2y4c"),
		AliasZoneName:   types.StringValue("alias-example"),
		Tags:            tags,
		AllowDuplicates: types.BoolValue(true),
		JumpStart:       types.BoolValue(false),
	}

	request, diags := buildTeoZone1110_01CreateZoneRequest(context.Background(), plan)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if request.Type == nil || *request.Type != "partial" {
		t.Errorf("unexpected Type: %+v", request.Type)
	}
	if request.ZoneName == nil || *request.ZoneName != "example.com" {
		t.Errorf("unexpected ZoneName: %+v", request.ZoneName)
	}
	if request.Area == nil || *request.Area != "mainland" {
		t.Errorf("unexpected Area: %+v", request.Area)
	}
	if request.PlanId == nil || *request.PlanId != "edgeone-2wdo315m2y4c" {
		t.Errorf("unexpected PlanId: %+v", request.PlanId)
	}
	if request.AliasZoneName == nil || *request.AliasZoneName != "alias-example" {
		t.Errorf("unexpected AliasZoneName: %+v", request.AliasZoneName)
	}
	if request.AllowDuplicates == nil || !*request.AllowDuplicates {
		t.Errorf("unexpected AllowDuplicates: %+v", request.AllowDuplicates)
	}
	if request.JumpStart == nil || *request.JumpStart {
		t.Errorf("unexpected JumpStart: %+v", request.JumpStart)
	}
	if len(request.Tags) != 2 {
		t.Fatalf("unexpected number of tags: %d", len(request.Tags))
	}
	if request.Tags[0].TagKey == nil || *request.Tags[0].TagKey != "env" ||
		request.Tags[0].TagValue == nil || *request.Tags[0].TagValue != "test" {
		t.Errorf("unexpected first tag: %+v", request.Tags[0])
	}
	if request.Tags[1].TagKey == nil || *request.Tags[1].TagKey != "owner" ||
		request.Tags[1].TagValue == nil || *request.Tags[1].TagValue != "iac" {
		t.Errorf("unexpected second tag: %+v", request.Tags[1])
	}
}

// TestTeoZone1110_01BuildCreateZoneRequestSkipsUnset verifies null and unknown
// values are not sent.
func TestTeoZone1110_01BuildCreateZoneRequestSkipsUnset(t *testing.T) {
	plan := &TeoZone111001ResourceModel{
		Type:            types.StringValue("noDomainAccess"),
		ZoneName:        types.StringUnknown(),
		Area:            types.StringNull(),
		PlanId:          types.StringNull(),
		AliasZoneName:   types.StringUnknown(),
		Tags:            types.ListNull(teoZone1110_01TagObjectType),
		AllowDuplicates: types.BoolNull(),
		JumpStart:       types.BoolUnknown(),
	}

	request, diags := buildTeoZone1110_01CreateZoneRequest(context.Background(), plan)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if request.Type == nil || *request.Type != "noDomainAccess" {
		t.Errorf("unexpected Type: %+v", request.Type)
	}
	if request.ZoneName != nil {
		t.Errorf("ZoneName must not be sent when unknown: %+v", request.ZoneName)
	}
	if request.Area != nil {
		t.Errorf("Area must not be sent when null: %+v", request.Area)
	}
	if request.PlanId != nil {
		t.Errorf("PlanId must not be sent when null: %+v", request.PlanId)
	}
	if request.AliasZoneName != nil {
		t.Errorf("AliasZoneName must not be sent when unknown: %+v", request.AliasZoneName)
	}
	if request.AllowDuplicates != nil {
		t.Errorf("AllowDuplicates must not be sent when null: %+v", request.AllowDuplicates)
	}
	if request.JumpStart != nil {
		t.Errorf("JumpStart must not be sent when unknown: %+v", request.JumpStart)
	}
	if len(request.Tags) != 0 {
		t.Errorf("Tags must not be sent when null: %+v", request.Tags)
	}
}

// TestTeoZone1110_01BuildModifyZoneRequest verifies only changed mutable values
// are sent.
func TestTeoZone1110_01BuildModifyZoneRequest(t *testing.T) {
	state := &TeoZone111001ResourceModel{
		ZoneName:      types.StringValue("example.com"),
		Type:          types.StringValue("partial"),
		Area:          types.StringValue("overseas"),
		AliasZoneName: types.StringValue("alias-example"),
	}
	plan := &TeoZone111001ResourceModel{
		ZoneName:      types.StringValue("example.com"),
		Type:          types.StringValue("partial"),
		Area:          types.StringValue("mainland"),
		AliasZoneName: types.StringValue("alias-new"),
	}

	request, diags := buildTeoZone1110_01ModifyZoneRequest(plan, state, teoZone1110_01TestZoneId)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if request.ZoneId == nil || *request.ZoneId != teoZone1110_01TestZoneId {
		t.Errorf("unexpected ZoneId: %+v", request.ZoneId)
	}
	if request.Type == nil || *request.Type != "partial" {
		t.Errorf("unexpected Type: %+v", request.Type)
	}
	if request.Area == nil || *request.Area != "mainland" {
		t.Errorf("unexpected Area: %+v", request.Area)
	}
	if request.AliasZoneName == nil || *request.AliasZoneName != "alias-new" {
		t.Errorf("unexpected AliasZoneName: %+v", request.AliasZoneName)
	}
	if request.ZoneName != nil {
		t.Errorf("ZoneName must not be sent when unchanged: %+v", request.ZoneName)
	}
}

// TestTeoZone1110_01BuildModifyZoneRequestZoneNameChange verifies a changed
// zone_name is sent when the zone switches from no domain access to CNAME access.
func TestTeoZone1110_01BuildModifyZoneRequestZoneNameChange(t *testing.T) {
	state := &TeoZone111001ResourceModel{
		ZoneName: types.StringNull(),
		Type:     types.StringValue("noDomainAccess"),
	}
	plan := &TeoZone111001ResourceModel{
		ZoneName: types.StringValue("example.com"),
		Type:     types.StringValue("partial"),
	}

	request, diags := buildTeoZone1110_01ModifyZoneRequest(plan, state, teoZone1110_01TestZoneId)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if request.ZoneName == nil || *request.ZoneName != "example.com" {
		t.Errorf("unexpected ZoneName: %+v", request.ZoneName)
	}
}

// TestTeoZone1110_01BuildModifyZoneRequestRemovedAlias verifies removing
// alias_zone_name is refused.
func TestTeoZone1110_01BuildModifyZoneRequestRemovedAlias(t *testing.T) {
	state := &TeoZone111001ResourceModel{
		ZoneName:      types.StringValue("example.com"),
		Type:          types.StringValue("partial"),
		Area:          types.StringValue("overseas"),
		AliasZoneName: types.StringValue("alias-example"),
	}
	plan := &TeoZone111001ResourceModel{
		ZoneName:      types.StringValue("example.com"),
		Type:          types.StringValue("partial"),
		Area:          types.StringValue("overseas"),
		AliasZoneName: types.StringNull(),
	}

	_, diags := buildTeoZone1110_01ModifyZoneRequest(plan, state, teoZone1110_01TestZoneId)
	if !diags.HasError() {
		t.Fatal("expected an error when alias_zone_name is removed from the configuration")
	}
}

// TestTeoZone1110_01RefreshModel verifies the cloud owned fields are copied while
// the create-only switches keep their configured values.
func TestTeoZone1110_01RefreshModel(t *testing.T) {
	model := &TeoZone111001ResourceModel{
		AllowDuplicates: types.BoolValue(true),
		JumpStart:       types.BoolValue(true),
	}
	zone := &teov20220901.Zone{
		ZoneId:        teoZone1110_01Ptr(teoZone1110_01TestZoneId),
		ZoneName:      teoZone1110_01Ptr("example.com"),
		AliasZoneName: teoZone1110_01Ptr("alias-example"),
		Area:          teoZone1110_01Ptr("overseas"),
		Type:          teoZone1110_01Ptr("partial"),
		Tags: []*teov20220901.Tag{
			{TagKey: teoZone1110_01Ptr("env"), TagValue: teoZone1110_01Ptr("test")},
		},
	}

	diags := refreshTeoZone1110_01Model(context.Background(), model, zone)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if model.ZoneId.ValueString() != teoZone1110_01TestZoneId {
		t.Errorf("unexpected zone_id %q", model.ZoneId.ValueString())
	}
	if model.ZoneName.ValueString() != "example.com" {
		t.Errorf("unexpected zone_name %q", model.ZoneName.ValueString())
	}
	if model.AliasZoneName.ValueString() != "alias-example" {
		t.Errorf("unexpected alias_zone_name %q", model.AliasZoneName.ValueString())
	}
	if model.Area.ValueString() != "overseas" {
		t.Errorf("unexpected area %q", model.Area.ValueString())
	}
	if model.Type.ValueString() != "partial" {
		t.Errorf("unexpected type %q", model.Type.ValueString())
	}
	if !model.AllowDuplicates.ValueBool() {
		t.Errorf("allow_duplicates must keep the configured value")
	}
	if !model.JumpStart.ValueBool() {
		t.Errorf("jump_start must keep the configured value")
	}

	var tags []TeoZone1110_01TagModel
	if d := model.Tags.ElementsAs(context.Background(), &tags, false); d.HasError() {
		t.Fatalf("failed to decode tags: %v", d)
	}
	if len(tags) != 1 {
		t.Fatalf("unexpected number of tags: %d", len(tags))
	}
	if tags[0].TagKey.ValueString() != "env" {
		t.Errorf("unexpected tag_key %q", tags[0].TagKey.ValueString())
	}
	if tags[0].TagValue.ValueString() != "test" {
		t.Errorf("unexpected tag_value %q", tags[0].TagValue.ValueString())
	}
}

// TestTeoZone1110_01RefreshModelNilZone verifies a nil zone leaves the model untouched.
func TestTeoZone1110_01RefreshModelNilZone(t *testing.T) {
	model := &TeoZone111001ResourceModel{
		ZoneName: types.StringValue("example.com"),
	}

	diags := refreshTeoZone1110_01Model(context.Background(), model, nil)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if model.ZoneName.ValueString() != "example.com" {
		t.Errorf("model must not be modified for a nil zone")
	}
}

// TestTeoZone1110_01BuildOwnershipVerification verifies all three verification
// blocks are flattened into state.
func TestTeoZone1110_01BuildOwnershipVerification(t *testing.T) {
	verification := &teov20220901.OwnershipVerification{
		DnsVerification: &teov20220901.DnsVerification{
			Subdomain:   teoZone1110_01Ptr("_teo-verification"),
			RecordType:  teoZone1110_01Ptr("TXT"),
			RecordValue: teoZone1110_01Ptr("record-value"),
		},
		FileVerification: &teov20220901.FileVerification{
			Path:    teoZone1110_01Ptr("/.well-known/teo-verification/abc.txt"),
			Content: teoZone1110_01Ptr("file-content"),
		},
		NsVerification: &teov20220901.NsVerification{
			NameServers: []*string{teoZone1110_01Ptr("ns1.teo.com"), teoZone1110_01Ptr("ns2.teo.com")},
		},
	}

	list := buildOwnershipVerificationForTest(t, verification)
	if list.IsNull() {
		t.Fatal("expected a non null ownership_verification list")
	}
	if len(list.Elements()) != 1 {
		t.Fatalf("expected exactly one ownership_verification block, got %d", len(list.Elements()))
	}

	var blocks []TeoZone1110_01OwnershipVerificationModel
	if d := list.ElementsAs(context.Background(), &blocks, false); d.HasError() {
		t.Fatalf("failed to decode ownership_verification: %v", d)
	}

	var dnsModels []TeoZone1110_01DnsVerificationModel
	if d := blocks[0].DnsVerification.ElementsAs(context.Background(), &dnsModels, false); d.HasError() {
		t.Fatalf("failed to decode dns_verification: %v", d)
	}
	if len(dnsModels) != 1 {
		t.Fatalf("unexpected number of dns_verification blocks: %d", len(dnsModels))
	}
	if dnsModels[0].Subdomain.ValueString() != "_teo-verification" {
		t.Errorf("unexpected subdomain %q", dnsModels[0].Subdomain.ValueString())
	}
	if dnsModels[0].RecordType.ValueString() != "TXT" {
		t.Errorf("unexpected record_type %q", dnsModels[0].RecordType.ValueString())
	}
	if dnsModels[0].RecordValue.ValueString() != "record-value" {
		t.Errorf("unexpected record_value %q", dnsModels[0].RecordValue.ValueString())
	}

	var fileModels []TeoZone1110_01FileVerificationModel
	if d := blocks[0].FileVerification.ElementsAs(context.Background(), &fileModels, false); d.HasError() {
		t.Fatalf("failed to decode file_verification: %v", d)
	}
	if len(fileModels) != 1 {
		t.Fatalf("unexpected number of file_verification blocks: %d", len(fileModels))
	}
	if fileModels[0].Path.ValueString() != "/.well-known/teo-verification/abc.txt" {
		t.Errorf("unexpected path %q", fileModels[0].Path.ValueString())
	}
	if fileModels[0].Content.ValueString() != "file-content" {
		t.Errorf("unexpected content %q", fileModels[0].Content.ValueString())
	}

	var nsModels []TeoZone1110_01NsVerificationModel
	if d := blocks[0].NsVerification.ElementsAs(context.Background(), &nsModels, false); d.HasError() {
		t.Fatalf("failed to decode ns_verification: %v", d)
	}
	if len(nsModels) != 1 {
		t.Fatalf("unexpected number of ns_verification blocks: %d", len(nsModels))
	}
	var nameServers []string
	if d := nsModels[0].NameServers.ElementsAs(context.Background(), &nameServers, false); d.HasError() {
		t.Fatalf("failed to decode name_servers: %v", d)
	}
	if len(nameServers) != 2 {
		t.Fatalf("unexpected number of name servers: %d", len(nameServers))
	}
	found := map[string]bool{}
	for _, ns := range nameServers {
		found[ns] = true
	}
	if !found["ns1.teo.com"] || !found["ns2.teo.com"] {
		t.Errorf("unexpected name servers: %+v", nameServers)
	}
}

// TestTeoZone1110_01BuildOwnershipVerificationNil verifies a nil verification
// input produces a null list rather than an empty one.
func TestTeoZone1110_01BuildOwnershipVerificationNil(t *testing.T) {
	list := buildOwnershipVerificationForTest(t, nil)
	if !list.IsNull() {
		t.Errorf("expected a null list for a nil verification, got %+v", list)
	}
}

// TestTeoZone1110_01BuildOwnershipVerificationEmpty verifies a non nil but empty
// verification produces empty nested collections.
func TestTeoZone1110_01BuildOwnershipVerificationEmpty(t *testing.T) {
	list := buildOwnershipVerificationForTest(t, &teov20220901.OwnershipVerification{})
	if list.IsNull() {
		t.Fatal("expected a non null list for a non nil verification")
	}

	var blocks []TeoZone1110_01OwnershipVerificationModel
	if d := list.ElementsAs(context.Background(), &blocks, false); d.HasError() {
		t.Fatalf("failed to decode ownership_verification: %v", d)
	}
	if len(blocks) != 1 {
		t.Fatalf("expected exactly one ownership_verification block, got %d", len(blocks))
	}
	if len(blocks[0].DnsVerification.Elements()) != 0 {
		t.Errorf("expected an empty dns_verification, got %d elements", len(blocks[0].DnsVerification.Elements()))
	}
	if len(blocks[0].FileVerification.Elements()) != 0 {
		t.Errorf("expected an empty file_verification, got %d elements", len(blocks[0].FileVerification.Elements()))
	}
	if len(blocks[0].NsVerification.Elements()) != 0 {
		t.Errorf("expected an empty ns_verification, got %d elements", len(blocks[0].NsVerification.Elements()))
	}
}

// buildOwnershipVerificationForTest calls the ownership verification builder and
// fails the test when it reports diagnostics.
func buildOwnershipVerificationForTest(t *testing.T, verification *teov20220901.OwnershipVerification) types.List {
	t.Helper()

	var diags diag.Diagnostics
	list := buildTeoZone1110_01OwnershipVerification(context.Background(), verification, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	return list
}

// teoZone1110_01Object builds one tags element.
func teoZone1110_01Object(t *testing.T, tagKey, tagValue string) attr.Value {
	t.Helper()

	obj, diags := types.ObjectValue(teoZone1110_01TagAttrTypes, map[string]attr.Value{
		"tag_key":   types.StringValue(tagKey),
		"tag_value": types.StringValue(tagValue),
	})
	if diags.HasError() {
		t.Fatalf("failed to build tag object: %v", diags)
	}
	return obj
}

// TestTeoZone1110_01FlattenStrings verifies nil entries are skipped.
func TestTeoZone1110_01FlattenStrings(t *testing.T) {
	got := flattenTeoZone1110_01Strings([]*string{teoZone1110_01Ptr("a"), nil, teoZone1110_01Ptr("b")})
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Errorf("unexpected flattened strings: %+v", got)
	}
	if flattened := flattenTeoZone1110_01Strings(nil); len(flattened) != 0 {
		t.Errorf("unexpected flattened strings for nil input: %+v", flattened)
	}
}

// TestTeoZone1110_01ServiceDeleteNotFoundIsSwallowed verifies the not found
// detection used by Delete only matches documented codes.
func TestTeoZone1110_01ServiceDeleteNotFoundIsSwallowed(t *testing.T) {
	client := &connectivity.TencentCloudClient{}
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "DeleteZoneWithContext",
		func(_ context.Context, request *teov20220901.DeleteZoneRequest) (*teov20220901.DeleteZoneResponse, error) {
			return nil, fmt.Errorf("ResourceNotFound: zone not found")
		})

	service := NewTeoService(client)
	err := service.DeleteTeoZoneBy1110_01(context.Background(), teoZone1110_01TestZoneId)
	if err == nil {
		t.Fatal("expected the service to return the API error")
	}
	if !tccommon.IsExpectError(err, []string{"ResourceNotFound", "InvalidParameter.ZoneNotExists"}) {
		t.Errorf("the not found error must be recognised as an expected error: %v", err)
	}
}

// teoZone1110_01EmptyState returns a state whose every attribute is null.
func teoZone1110_01EmptyState(t *testing.T, zoneSchema schema.Schema) tfsdk.State {
	t.Helper()

	raw := tftypes.NewValue(zoneSchema.Type().TerraformType(context.Background()), nil)

	return tfsdk.State{Schema: zoneSchema, Raw: raw}
}

// teoZone1110_01Plan builds a plan from the given model.
func teoZone1110_01Plan(t *testing.T, zoneSchema schema.Schema, model *TeoZone111001ResourceModel) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{Schema: zoneSchema}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build plan: %v", diags)
	}
	return plan
}

// teoZone1110_01State builds a state from the given model.
func teoZone1110_01State(t *testing.T, zoneSchema schema.Schema, model *TeoZone111001ResourceModel) tfsdk.State {
	t.Helper()

	state := tfsdk.State{Schema: zoneSchema}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("failed to build state: %v", diags)
	}
	return state
}

// teoZone1110_01FullZone returns a zone with every readable field populated.
func teoZone1110_01FullZone() *teov20220901.Zone {
	return &teov20220901.Zone{
		ZoneId:        teoZone1110_01Ptr(teoZone1110_01TestZoneId),
		ZoneName:      teoZone1110_01Ptr("example.com"),
		AliasZoneName: teoZone1110_01Ptr("alias-example"),
		Area:          teoZone1110_01Ptr("overseas"),
		Type:          teoZone1110_01Ptr("partial"),
		Tags: []*teov20220901.Tag{
			{TagKey: teoZone1110_01Ptr("env"), TagValue: teoZone1110_01Ptr("test")},
		},
	}
}

// teoZone1110_01DescribeResponse builds a DescribeZones response.
func teoZone1110_01DescribeResponse(zones []*teov20220901.Zone) *teov20220901.DescribeZonesResponse {
	resp := teov20220901.NewDescribeZonesResponse()
	total := int64(len(zones))
	resp.Response = &teov20220901.DescribeZonesResponseParams{
		TotalCount: &total,
		Zones:      zones,
		RequestId:  teoZone1110_01Ptr("request-id"),
	}
	return resp
}

// TestTeoZone1110_01CreateLifecycle drives Create through the framework and
// verifies the request, the state and the ownership verification output.
func TestTeoZone1110_01CreateLifecycle(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res := teoZone1110_01NewTestResource(t)
	client := res.(*TeoZone111001Resource).Client()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(client, "UseTeoClient", teoClient)

	var capturedCreate *teov20220901.CreateZoneRequest
	patches.ApplyMethodFunc(teoClient, "CreateZoneWithContext",
		func(_ context.Context, request *teov20220901.CreateZoneRequest) (*teov20220901.CreateZoneResponse, error) {
			capturedCreate = request
			resp := teov20220901.NewCreateZoneResponse()
			resp.Response = &teov20220901.CreateZoneResponseParams{
				ZoneId:    teoZone1110_01Ptr(teoZone1110_01TestZoneId),
				RequestId: teoZone1110_01Ptr("request-id"),
				OwnershipVerification: &teov20220901.OwnershipVerification{
					DnsVerification: &teov20220901.DnsVerification{
						Subdomain:   teoZone1110_01Ptr("_teo-verification"),
						RecordType:  teoZone1110_01Ptr("TXT"),
						RecordValue: teoZone1110_01Ptr("record-value"),
					},
					FileVerification: &teov20220901.FileVerification{
						Path:    teoZone1110_01Ptr("/.well-known/teo-verification/abc.txt"),
						Content: teoZone1110_01Ptr("file-content"),
					},
					NsVerification: &teov20220901.NsVerification{
						NameServers: []*string{teoZone1110_01Ptr("ns1.teo.com"), teoZone1110_01Ptr("ns2.teo.com")},
					},
				},
			}
			return resp, nil
		})
	patches.ApplyMethodFunc(teoClient, "DescribeZonesWithContext",
		func(_ context.Context, request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
			if len(request.Filters) != 1 || request.Filters[0].Name == nil ||
				*request.Filters[0].Name != "zone-id" {
				t.Errorf("DescribeZones must be filtered by zone-id, got %+v", request.Filters)
			}
			if len(request.Filters[0].Values) != 1 || request.Filters[0].Values[0] == nil ||
				*request.Filters[0].Values[0] != teoZone1110_01TestZoneId {
				t.Errorf("unexpected filter values: %+v", request.Filters[0].Values)
			}
			return teoZone1110_01DescribeResponse([]*teov20220901.Zone{teoZone1110_01FullZone()}), nil
		})

	zoneSchema := teoZone1110_01Schema(t)
	tags, diags := types.ListValue(teoZone1110_01TagObjectType, []attr.Value{teoZone1110_01Object(t, "env", "test")})
	if diags.HasError() {
		t.Fatalf("failed to build tags: %v", diags)
	}

	planModel := &TeoZone111001ResourceModel{
		Id:                    types.StringUnknown(),
		ZoneId:                types.StringUnknown(),
		Type:                  types.StringValue("partial"),
		ZoneName:              types.StringValue("example.com"),
		Area:                  types.StringValue("overseas"),
		PlanId:                types.StringValue("edgeone-2wdo315m2y4c"),
		AliasZoneName:         types.StringValue("alias-example"),
		Tags:                  tags,
		AllowDuplicates:       types.BoolValue(false),
		JumpStart:             types.BoolValue(true),
		OwnershipVerification: types.ListNull(teoZone1110_01OwnershipVerificationObjectType),
	}

	createResp := &resource.CreateResponse{State: teoZone1110_01EmptyState(t, zoneSchema)}
	res.Create(context.Background(),
		resource.CreateRequest{Plan: teoZone1110_01Plan(t, zoneSchema, planModel)}, createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("unexpected create diagnostics: %v", createResp.Diagnostics)
	}

	if capturedCreate == nil {
		t.Fatal("CreateZone was not called")
	}
	if capturedCreate.Type == nil || *capturedCreate.Type != "partial" {
		t.Errorf("unexpected Type: %+v", capturedCreate.Type)
	}
	if capturedCreate.ZoneName == nil || *capturedCreate.ZoneName != "example.com" {
		t.Errorf("unexpected ZoneName: %+v", capturedCreate.ZoneName)
	}
	if capturedCreate.Area == nil || *capturedCreate.Area != "overseas" {
		t.Errorf("unexpected Area: %+v", capturedCreate.Area)
	}
	if capturedCreate.PlanId == nil || *capturedCreate.PlanId != "edgeone-2wdo315m2y4c" {
		t.Errorf("unexpected PlanId: %+v", capturedCreate.PlanId)
	}
	if capturedCreate.AliasZoneName == nil || *capturedCreate.AliasZoneName != "alias-example" {
		t.Errorf("unexpected AliasZoneName: %+v", capturedCreate.AliasZoneName)
	}
	if capturedCreate.AllowDuplicates == nil || *capturedCreate.AllowDuplicates {
		t.Errorf("unexpected AllowDuplicates: %+v", capturedCreate.AllowDuplicates)
	}
	if capturedCreate.JumpStart == nil || !*capturedCreate.JumpStart {
		t.Errorf("unexpected JumpStart: %+v", capturedCreate.JumpStart)
	}
	if len(capturedCreate.Tags) != 1 || capturedCreate.Tags[0].TagKey == nil ||
		*capturedCreate.Tags[0].TagKey != "env" || capturedCreate.Tags[0].TagValue == nil ||
		*capturedCreate.Tags[0].TagValue != "test" {
		t.Errorf("unexpected Tags: %+v", capturedCreate.Tags)
	}

	var stateModel TeoZone111001ResourceModel
	if diags := createResp.State.Get(context.Background(), &stateModel); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags)
	}

	if stateModel.Id.ValueString() != teoZone1110_01TestZoneId {
		t.Errorf("unexpected id %q", stateModel.Id.ValueString())
	}
	if stateModel.ZoneId.ValueString() != teoZone1110_01TestZoneId {
		t.Errorf("unexpected zone_id %q", stateModel.ZoneId.ValueString())
	}
	if stateModel.ZoneName.ValueString() != "example.com" {
		t.Errorf("unexpected zone_name %q", stateModel.ZoneName.ValueString())
	}
	if stateModel.Type.ValueString() != "partial" {
		t.Errorf("unexpected type %q", stateModel.Type.ValueString())
	}
	if stateModel.Area.ValueString() != "overseas" {
		t.Errorf("unexpected area %q", stateModel.Area.ValueString())
	}
	if stateModel.PlanId.ValueString() != "edgeone-2wdo315m2y4c" {
		t.Errorf("unexpected plan_id %q", stateModel.PlanId.ValueString())
	}
	if stateModel.AliasZoneName.ValueString() != "alias-example" {
		t.Errorf("unexpected alias_zone_name %q", stateModel.AliasZoneName.ValueString())
	}
	if stateModel.AllowDuplicates.ValueBool() {
		t.Errorf("allow_duplicates must keep the configured value")
	}
	if !stateModel.JumpStart.ValueBool() {
		t.Errorf("jump_start must keep the configured value")
	}

	var stateTags []TeoZone1110_01TagModel
	if diags := stateModel.Tags.ElementsAs(context.Background(), &stateTags, false); diags.HasError() {
		t.Fatalf("failed to decode tags: %v", diags)
	}
	if len(stateTags) != 1 || stateTags[0].TagKey.ValueString() != "env" ||
		stateTags[0].TagValue.ValueString() != "test" {
		t.Errorf("unexpected tags in state: %+v", stateTags)
	}

	var verifications []TeoZone1110_01OwnershipVerificationModel
	if diags := stateModel.OwnershipVerification.ElementsAs(context.Background(), &verifications, false); diags.HasError() {
		t.Fatalf("failed to decode ownership_verification: %v", diags)
	}
	if len(verifications) != 1 {
		t.Fatalf("expected one ownership_verification block, got %d", len(verifications))
	}
	var dnsModels []TeoZone1110_01DnsVerificationModel
	if diags := verifications[0].DnsVerification.ElementsAs(context.Background(), &dnsModels, false); diags.HasError() {
		t.Fatalf("failed to decode dns_verification: %v", diags)
	}
	if len(dnsModels) != 1 || dnsModels[0].Subdomain.ValueString() != "_teo-verification" ||
		dnsModels[0].RecordType.ValueString() != "TXT" ||
		dnsModels[0].RecordValue.ValueString() != "record-value" {
		t.Errorf("unexpected dns_verification: %+v", dnsModels)
	}
	var fileModels []TeoZone1110_01FileVerificationModel
	if diags := verifications[0].FileVerification.ElementsAs(context.Background(), &fileModels, false); diags.HasError() {
		t.Fatalf("failed to decode file_verification: %v", diags)
	}
	if len(fileModels) != 1 || fileModels[0].Path.ValueString() != "/.well-known/teo-verification/abc.txt" ||
		fileModels[0].Content.ValueString() != "file-content" {
		t.Errorf("unexpected file_verification: %+v", fileModels)
	}
	var nsModels []TeoZone1110_01NsVerificationModel
	if diags := verifications[0].NsVerification.ElementsAs(context.Background(), &nsModels, false); diags.HasError() {
		t.Fatalf("failed to decode ns_verification: %v", diags)
	}
	if len(nsModels) != 1 {
		t.Fatalf("expected one ns_verification block, got %d", len(nsModels))
	}
	var nameServers []string
	if diags := nsModels[0].NameServers.ElementsAs(context.Background(), &nameServers, false); diags.HasError() {
		t.Fatalf("failed to decode name_servers: %v", diags)
	}
	if len(nameServers) != 2 {
		t.Errorf("unexpected name servers: %+v", nameServers)
	}
}

// TestTeoZone1110_01CreateEmptyZoneId verifies an empty zone id is an error and
// no identity is written into state.
func TestTeoZone1110_01CreateEmptyZoneId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res := teoZone1110_01NewTestResource(t)
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(res.(*TeoZone111001Resource).Client(), "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "CreateZoneWithContext",
		func(_ context.Context, request *teov20220901.CreateZoneRequest) (*teov20220901.CreateZoneResponse, error) {
			return teov20220901.NewCreateZoneResponse(), nil
		})

	zoneSchema := teoZone1110_01Schema(t)
	planModel := &TeoZone111001ResourceModel{
		Id:                    types.StringUnknown(),
		ZoneId:                types.StringUnknown(),
		Type:                  types.StringValue("noDomainAccess"),
		ZoneName:              types.StringNull(),
		Area:                  types.StringNull(),
		PlanId:                types.StringNull(),
		AliasZoneName:         types.StringNull(),
		Tags:                  types.ListNull(teoZone1110_01TagObjectType),
		AllowDuplicates:       types.BoolNull(),
		JumpStart:             types.BoolNull(),
		OwnershipVerification: types.ListNull(teoZone1110_01OwnershipVerificationObjectType),
	}

	createResp := &resource.CreateResponse{State: teoZone1110_01EmptyState(t, zoneSchema)}
	res.Create(context.Background(),
		resource.CreateRequest{Plan: teoZone1110_01Plan(t, zoneSchema, planModel)}, createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("expected an error for an empty zone id")
	}
}

// TestTeoZone1110_01CreateAPIError verifies CreateZone errors are reported.
func TestTeoZone1110_01CreateAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res := teoZone1110_01NewTestResource(t)
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(res.(*TeoZone111001Resource).Client(), "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "CreateZoneWithContext",
		func(_ context.Context, request *teov20220901.CreateZoneRequest) (*teov20220901.CreateZoneResponse, error) {
			return nil, fmt.Errorf("CreateZone failed")
		})

	zoneSchema := teoZone1110_01Schema(t)
	planModel := &TeoZone111001ResourceModel{
		Id:                    types.StringUnknown(),
		ZoneId:                types.StringUnknown(),
		Type:                  types.StringValue("partial"),
		ZoneName:              types.StringValue("example.com"),
		Area:                  types.StringValue("overseas"),
		PlanId:                types.StringNull(),
		AliasZoneName:         types.StringNull(),
		Tags:                  types.ListNull(teoZone1110_01TagObjectType),
		AllowDuplicates:       types.BoolNull(),
		JumpStart:             types.BoolNull(),
		OwnershipVerification: types.ListNull(teoZone1110_01OwnershipVerificationObjectType),
	}

	createResp := &resource.CreateResponse{State: teoZone1110_01EmptyState(t, zoneSchema)}
	res.Create(context.Background(),
		resource.CreateRequest{Plan: teoZone1110_01Plan(t, zoneSchema, planModel)}, createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("expected an error when CreateZone fails")
	}
}

// TestTeoZone1110_01ReadLifecycle verifies Read refreshes the cloud owned fields
// and keeps the create-only switches.
func TestTeoZone1110_01ReadLifecycle(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res := teoZone1110_01NewTestResource(t)
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(res.(*TeoZone111001Resource).Client(), "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "DescribeZonesWithContext",
		func(_ context.Context, request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
			return teoZone1110_01DescribeResponse([]*teov20220901.Zone{teoZone1110_01FullZone()}), nil
		})

	zoneSchema := teoZone1110_01Schema(t)
	stateModel := &TeoZone111001ResourceModel{
		Id:                    types.StringValue(teoZone1110_01TestZoneId),
		ZoneId:                types.StringValue(teoZone1110_01TestZoneId),
		Type:                  types.StringValue("partial"),
		ZoneName:              types.StringValue("example.com"),
		Area:                  types.StringValue("overseas"),
		PlanId:                types.StringValue("edgeone-2wdo315m2y4c"),
		AliasZoneName:         types.StringValue("alias-example"),
		Tags:                  types.ListNull(teoZone1110_01TagObjectType),
		AllowDuplicates:       types.BoolValue(true),
		JumpStart:             types.BoolValue(true),
		OwnershipVerification: types.ListNull(teoZone1110_01OwnershipVerificationObjectType),
	}

	readResp := &resource.ReadResponse{State: teoZone1110_01State(t, zoneSchema, stateModel)}
	res.Read(context.Background(),
		resource.ReadRequest{State: teoZone1110_01State(t, zoneSchema, stateModel)}, readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("unexpected read diagnostics: %v", readResp.Diagnostics)
	}

	var refreshed TeoZone111001ResourceModel
	if diags := readResp.State.Get(context.Background(), &refreshed); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags)
	}
	if refreshed.Id.ValueString() != teoZone1110_01TestZoneId {
		t.Errorf("unexpected id %q", refreshed.Id.ValueString())
	}
	if refreshed.ZoneId.ValueString() != teoZone1110_01TestZoneId {
		t.Errorf("unexpected zone_id %q", refreshed.ZoneId.ValueString())
	}
	if refreshed.Type.ValueString() != "partial" {
		t.Errorf("unexpected type %q", refreshed.Type.ValueString())
	}
	if refreshed.Area.ValueString() != "overseas" {
		t.Errorf("unexpected area %q", refreshed.Area.ValueString())
	}
	if refreshed.AliasZoneName.ValueString() != "alias-example" {
		t.Errorf("unexpected alias_zone_name %q", refreshed.AliasZoneName.ValueString())
	}
	if !refreshed.AllowDuplicates.ValueBool() {
		t.Errorf("allow_duplicates must be kept from state")
	}
	if !refreshed.JumpStart.ValueBool() {
		t.Errorf("jump_start must be kept from state")
	}
	var refreshedTags []TeoZone1110_01TagModel
	if diags := refreshed.Tags.ElementsAs(context.Background(), &refreshedTags, false); diags.HasError() {
		t.Fatalf("failed to decode tags: %v", diags)
	}
	if len(refreshedTags) != 1 || refreshedTags[0].TagKey.ValueString() != "env" {
		t.Errorf("unexpected tags in state: %+v", refreshedTags)
	}
}

// TestTeoZone1110_01ReadNotFound verifies a zone missing on the cloud side is
// removed from state.
func TestTeoZone1110_01ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res := teoZone1110_01NewTestResource(t)
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(res.(*TeoZone111001Resource).Client(), "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "DescribeZonesWithContext",
		func(_ context.Context, request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
			return teoZone1110_01DescribeResponse(nil), nil
		})

	zoneSchema := teoZone1110_01Schema(t)
	stateModel := &TeoZone111001ResourceModel{
		Id:                    types.StringValue(teoZone1110_01TestZoneId),
		ZoneId:                types.StringValue(teoZone1110_01TestZoneId),
		Type:                  types.StringNull(),
		ZoneName:              types.StringValue("example.com"),
		Area:                  types.StringNull(),
		PlanId:                types.StringNull(),
		AliasZoneName:         types.StringNull(),
		Tags:                  types.ListNull(teoZone1110_01TagObjectType),
		AllowDuplicates:       types.BoolNull(),
		JumpStart:             types.BoolNull(),
		OwnershipVerification: types.ListNull(teoZone1110_01OwnershipVerificationObjectType),
	}

	readResp := &resource.ReadResponse{State: teoZone1110_01State(t, zoneSchema, stateModel)}
	res.Read(context.Background(),
		resource.ReadRequest{State: teoZone1110_01State(t, zoneSchema, stateModel)}, readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("unexpected read diagnostics: %v", readResp.Diagnostics)
	}
	if !readResp.State.Raw.IsNull() {
		t.Errorf("state must be removed when the zone no longer exists")
	}
}

// TestTeoZone1110_01ReadNilResponse verifies a nil response is an error and does
// not clear state.
func TestTeoZone1110_01ReadNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res := teoZone1110_01NewTestResource(t)
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(res.(*TeoZone111001Resource).Client(), "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "DescribeZonesWithContext",
		func(_ context.Context, request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
			return nil, nil
		})

	zoneSchema := teoZone1110_01Schema(t)
	stateModel := &TeoZone111001ResourceModel{
		Id:                    types.StringValue(teoZone1110_01TestZoneId),
		ZoneId:                types.StringValue(teoZone1110_01TestZoneId),
		Type:                  types.StringNull(),
		ZoneName:              types.StringValue("example.com"),
		Area:                  types.StringNull(),
		PlanId:                types.StringNull(),
		AliasZoneName:         types.StringNull(),
		Tags:                  types.ListNull(teoZone1110_01TagObjectType),
		AllowDuplicates:       types.BoolNull(),
		JumpStart:             types.BoolNull(),
		OwnershipVerification: types.ListNull(teoZone1110_01OwnershipVerificationObjectType),
	}

	readResp := &resource.ReadResponse{State: teoZone1110_01State(t, zoneSchema, stateModel)}
	res.Read(context.Background(),
		resource.ReadRequest{State: teoZone1110_01State(t, zoneSchema, stateModel)}, readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("expected an error for a nil DescribeZones response")
	}
}

// TestTeoZone1110_01ReadAPIError verifies DescribeZones errors are reported.
func TestTeoZone1110_01ReadAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res := teoZone1110_01NewTestResource(t)
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(res.(*TeoZone111001Resource).Client(), "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "DescribeZonesWithContext",
		func(_ context.Context, request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
			return nil, fmt.Errorf("DescribeZones failed")
		})

	zoneSchema := teoZone1110_01Schema(t)
	stateModel := &TeoZone111001ResourceModel{
		Id:                    types.StringValue(teoZone1110_01TestZoneId),
		ZoneId:                types.StringValue(teoZone1110_01TestZoneId),
		Type:                  types.StringNull(),
		ZoneName:              types.StringValue("example.com"),
		Area:                  types.StringNull(),
		PlanId:                types.StringNull(),
		AliasZoneName:         types.StringNull(),
		Tags:                  types.ListNull(teoZone1110_01TagObjectType),
		AllowDuplicates:       types.BoolNull(),
		JumpStart:             types.BoolNull(),
		OwnershipVerification: types.ListNull(teoZone1110_01OwnershipVerificationObjectType),
	}

	readResp := &resource.ReadResponse{State: teoZone1110_01State(t, zoneSchema, stateModel)}
	res.Read(context.Background(),
		resource.ReadRequest{State: teoZone1110_01State(t, zoneSchema, stateModel)}, readResp)
	if !readResp.Diagnostics.HasError() {
		t.Fatal("expected an error when DescribeZones fails")
	}
}

// TestTeoZone1110_01UpdateLifecycle verifies Update sends the changed mutable
// fields and refreshes state from the cloud side.
func TestTeoZone1110_01UpdateLifecycle(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res := teoZone1110_01NewTestResource(t)
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(res.(*TeoZone111001Resource).Client(), "UseTeoClient", teoClient)

	var capturedModify *teov20220901.ModifyZoneRequest
	patches.ApplyMethodFunc(teoClient, "ModifyZoneWithContext",
		func(_ context.Context, request *teov20220901.ModifyZoneRequest) (*teov20220901.ModifyZoneResponse, error) {
			capturedModify = request
			return teov20220901.NewModifyZoneResponse(), nil
		})
	patches.ApplyMethodFunc(teoClient, "DescribeZonesWithContext",
		func(_ context.Context, request *teov20220901.DescribeZonesRequest) (*teov20220901.DescribeZonesResponse, error) {
			zone := teoZone1110_01FullZone()
			zone.Area = teoZone1110_01Ptr("mainland")
			zone.AliasZoneName = teoZone1110_01Ptr("alias-new")
			return teoZone1110_01DescribeResponse([]*teov20220901.Zone{zone}), nil
		})

	zoneSchema := teoZone1110_01Schema(t)
	stateModel := &TeoZone111001ResourceModel{
		Id:                    types.StringValue(teoZone1110_01TestZoneId),
		ZoneId:                types.StringValue(teoZone1110_01TestZoneId),
		Type:                  types.StringValue("partial"),
		ZoneName:              types.StringValue("example.com"),
		Area:                  types.StringValue("overseas"),
		PlanId:                types.StringValue("edgeone-2wdo315m2y4c"),
		AliasZoneName:         types.StringValue("alias-example"),
		Tags:                  types.ListNull(teoZone1110_01TagObjectType),
		AllowDuplicates:       types.BoolValue(false),
		JumpStart:             types.BoolValue(false),
		OwnershipVerification: types.ListNull(teoZone1110_01OwnershipVerificationObjectType),
	}
	planModel := &TeoZone111001ResourceModel{
		Id:                    types.StringValue(teoZone1110_01TestZoneId),
		ZoneId:                types.StringValue(teoZone1110_01TestZoneId),
		Type:                  types.StringValue("partial"),
		ZoneName:              types.StringValue("example.com"),
		Area:                  types.StringValue("mainland"),
		PlanId:                types.StringValue("edgeone-2wdo315m2y4c"),
		AliasZoneName:         types.StringValue("alias-new"),
		Tags:                  types.ListNull(teoZone1110_01TagObjectType),
		AllowDuplicates:       types.BoolValue(false),
		JumpStart:             types.BoolValue(false),
		OwnershipVerification: types.ListNull(teoZone1110_01OwnershipVerificationObjectType),
	}

	updateResp := &resource.UpdateResponse{State: teoZone1110_01State(t, zoneSchema, stateModel)}
	res.Update(context.Background(), resource.UpdateRequest{
		Plan:  teoZone1110_01Plan(t, zoneSchema, planModel),
		State: teoZone1110_01State(t, zoneSchema, stateModel),
	}, updateResp)
	if updateResp.Diagnostics.HasError() {
		t.Fatalf("unexpected update diagnostics: %v", updateResp.Diagnostics)
	}

	if capturedModify == nil {
		t.Fatal("ModifyZone was not called")
	}
	if capturedModify.ZoneId == nil || *capturedModify.ZoneId != teoZone1110_01TestZoneId {
		t.Errorf("unexpected ZoneId: %+v", capturedModify.ZoneId)
	}
	if capturedModify.Area == nil || *capturedModify.Area != "mainland" {
		t.Errorf("unexpected Area: %+v", capturedModify.Area)
	}
	if capturedModify.AliasZoneName == nil || *capturedModify.AliasZoneName != "alias-new" {
		t.Errorf("unexpected AliasZoneName: %+v", capturedModify.AliasZoneName)
	}
	if capturedModify.ZoneName != nil {
		t.Errorf("zone_name must not be sent when unchanged: %+v", capturedModify.ZoneName)
	}

	var updated TeoZone111001ResourceModel
	if diags := updateResp.State.Get(context.Background(), &updated); diags.HasError() {
		t.Fatalf("failed to read state: %v", diags)
	}
	if updated.Area.ValueString() != "mainland" {
		t.Errorf("unexpected area %q", updated.Area.ValueString())
	}
	if updated.AliasZoneName.ValueString() != "alias-new" {
		t.Errorf("unexpected alias_zone_name %q", updated.AliasZoneName.ValueString())
	}
}

// TestTeoZone1110_01UpdateRemovedAliasZoneName verifies removing
// alias_zone_name is refused before the API is called.
func TestTeoZone1110_01UpdateRemovedAliasZoneName(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res := teoZone1110_01NewTestResource(t)
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(res.(*TeoZone111001Resource).Client(), "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "ModifyZoneWithContext",
		func(_ context.Context, request *teov20220901.ModifyZoneRequest) (*teov20220901.ModifyZoneResponse, error) {
			t.Error("ModifyZone must not be called when alias_zone_name is removed")
			return teov20220901.NewModifyZoneResponse(), nil
		})

	zoneSchema := teoZone1110_01Schema(t)
	stateModel := &TeoZone111001ResourceModel{
		Id:                    types.StringValue(teoZone1110_01TestZoneId),
		ZoneId:                types.StringValue(teoZone1110_01TestZoneId),
		Type:                  types.StringValue("partial"),
		ZoneName:              types.StringValue("example.com"),
		Area:                  types.StringValue("overseas"),
		PlanId:                types.StringValue("edgeone-2wdo315m2y4c"),
		AliasZoneName:         types.StringValue("alias-example"),
		Tags:                  types.ListNull(teoZone1110_01TagObjectType),
		AllowDuplicates:       types.BoolValue(false),
		JumpStart:             types.BoolValue(false),
		OwnershipVerification: types.ListNull(teoZone1110_01OwnershipVerificationObjectType),
	}
	planModel := &TeoZone111001ResourceModel{
		Id:                    types.StringValue(teoZone1110_01TestZoneId),
		ZoneId:                types.StringValue(teoZone1110_01TestZoneId),
		Type:                  types.StringValue("partial"),
		ZoneName:              types.StringValue("example.com"),
		Area:                  types.StringValue("overseas"),
		PlanId:                types.StringValue("edgeone-2wdo315m2y4c"),
		AliasZoneName:         types.StringNull(),
		Tags:                  types.ListNull(teoZone1110_01TagObjectType),
		AllowDuplicates:       types.BoolValue(false),
		JumpStart:             types.BoolValue(false),
		OwnershipVerification: types.ListNull(teoZone1110_01OwnershipVerificationObjectType),
	}

	updateResp := &resource.UpdateResponse{State: teoZone1110_01State(t, zoneSchema, stateModel)}
	res.Update(context.Background(), resource.UpdateRequest{
		Plan:  teoZone1110_01Plan(t, zoneSchema, planModel),
		State: teoZone1110_01State(t, zoneSchema, stateModel),
	}, updateResp)
	if !updateResp.Diagnostics.HasError() {
		t.Fatal("expected an error when alias_zone_name is removed")
	}
}

// TestTeoZone1110_01DeleteLifecycle verifies Delete sends the zone id.
func TestTeoZone1110_01DeleteLifecycle(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res := teoZone1110_01NewTestResource(t)
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(res.(*TeoZone111001Resource).Client(), "UseTeoClient", teoClient)

	var capturedDelete *teov20220901.DeleteZoneRequest
	patches.ApplyMethodFunc(teoClient, "DeleteZoneWithContext",
		func(_ context.Context, request *teov20220901.DeleteZoneRequest) (*teov20220901.DeleteZoneResponse, error) {
			capturedDelete = request
			return teov20220901.NewDeleteZoneResponse(), nil
		})

	zoneSchema := teoZone1110_01Schema(t)
	stateModel := &TeoZone111001ResourceModel{
		Id:                    types.StringValue(teoZone1110_01TestZoneId),
		ZoneId:                types.StringValue(teoZone1110_01TestZoneId),
		Type:                  types.StringValue("partial"),
		ZoneName:              types.StringValue("example.com"),
		Area:                  types.StringValue("overseas"),
		PlanId:                types.StringValue("edgeone-2wdo315m2y4c"),
		AliasZoneName:         types.StringNull(),
		Tags:                  types.ListNull(teoZone1110_01TagObjectType),
		AllowDuplicates:       types.BoolValue(false),
		JumpStart:             types.BoolValue(false),
		OwnershipVerification: types.ListNull(teoZone1110_01OwnershipVerificationObjectType),
	}

	deleteResp := &resource.DeleteResponse{State: teoZone1110_01State(t, zoneSchema, stateModel)}
	res.Delete(context.Background(),
		resource.DeleteRequest{State: teoZone1110_01State(t, zoneSchema, stateModel)}, deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("unexpected delete diagnostics: %v", deleteResp.Diagnostics)
	}
	if capturedDelete == nil || capturedDelete.ZoneId == nil ||
		*capturedDelete.ZoneId != teoZone1110_01TestZoneId {
		t.Errorf("unexpected DeleteZone request: %+v", capturedDelete)
	}
}

// TestTeoZone1110_01DeleteAPIError verifies an unrelated delete error keeps the
// zone in state.
func TestTeoZone1110_01DeleteAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	res := teoZone1110_01NewTestResource(t)
	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(res.(*TeoZone111001Resource).Client(), "UseTeoClient", teoClient)
	patches.ApplyMethodFunc(teoClient, "DeleteZoneWithContext",
		func(_ context.Context, request *teov20220901.DeleteZoneRequest) (*teov20220901.DeleteZoneResponse, error) {
			return nil, fmt.Errorf("InternalError: something went wrong")
		})

	zoneSchema := teoZone1110_01Schema(t)
	stateModel := &TeoZone111001ResourceModel{
		Id:                    types.StringValue(teoZone1110_01TestZoneId),
		ZoneId:                types.StringValue(teoZone1110_01TestZoneId),
		Type:                  types.StringValue("partial"),
		ZoneName:              types.StringValue("example.com"),
		Area:                  types.StringValue("overseas"),
		PlanId:                types.StringValue("edgeone-2wdo315m2y4c"),
		AliasZoneName:         types.StringNull(),
		Tags:                  types.ListNull(teoZone1110_01TagObjectType),
		AllowDuplicates:       types.BoolValue(false),
		JumpStart:             types.BoolValue(false),
		OwnershipVerification: types.ListNull(teoZone1110_01OwnershipVerificationObjectType),
	}

	deleteResp := &resource.DeleteResponse{State: teoZone1110_01State(t, zoneSchema, stateModel)}
	res.Delete(context.Background(),
		resource.DeleteRequest{State: teoZone1110_01State(t, zoneSchema, stateModel)}, deleteResp)
	if !deleteResp.Diagnostics.HasError() {
		t.Fatal("expected an error when DeleteZone fails")
	}
}

// TestTeoZone1110_01ImportState verifies both id and zone_id are set so the
// first plan after import does not try to replace the zone.
func TestTeoZone1110_01ImportState(t *testing.T) {
	res := teoZone1110_01NewTestResource(t)
	zoneSchema := teoZone1110_01Schema(t)

	importResp := &resource.ImportStateResponse{State: teoZone1110_01EmptyState(t, zoneSchema)}
	res.(resource.ResourceWithImportState).ImportState(context.Background(),
		resource.ImportStateRequest{ID: teoZone1110_01TestZoneId}, importResp)
	if importResp.Diagnostics.HasError() {
		t.Fatalf("unexpected import diagnostics: %v", importResp.Diagnostics)
	}

	var importedId types.String
	if diags := importResp.State.GetAttribute(context.Background(), path.Root("id"), &importedId); diags.HasError() {
		t.Fatalf("failed to read id: %v", diags)
	}
	if importedId.ValueString() != teoZone1110_01TestZoneId {
		t.Errorf("unexpected id %q", importedId.ValueString())
	}

	var importedZoneId types.String
	if diags := importResp.State.GetAttribute(context.Background(), path.Root("zone_id"), &importedZoneId); diags.HasError() {
		t.Fatalf("failed to read zone_id: %v", diags)
	}
	if importedZoneId.ValueString() != teoZone1110_01TestZoneId {
		t.Errorf("unexpected zone_id %q", importedZoneId.ValueString())
	}
}
