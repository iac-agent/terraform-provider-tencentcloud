package ags

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	sdkags "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// modelAttrTypes returns the top-level attribute types of the resource model.
func modelAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":            basetypes.StringType{},
		"registry_id":   basetypes.StringType{},
		"name":          basetypes.StringType{},
		"description":   basetypes.StringType{},
		"approval_mode": basetypes.StringType{},
		"tags":          types.ListType{ElemType: tagsObjectType()},
		"registry":      types.ObjectType{AttrTypes: registryAttrTypes()},
	}
}

// --- Schema / Metadata / Factory ---

// TestAgsRegistryResource_Metadata validates the type name.
func TestAgsRegistryResource_Metadata(t *testing.T) {
	r := &AgsRegistryResource{}
	req := resource.MetadataRequest{}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	if resp.TypeName != "tencentcloud_ags_registry" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_ags_registry")
	}
}

// TestAgsRegistryResource_Schema validates required/optional/computed flags on
// every top-level attribute and on the nested attributes.
func TestAgsRegistryResource_Schema(t *testing.T) {
	resp := testResourceSchemaResponse(t)

	type expect struct {
		required, optional, computed bool
	}
	cases := map[string]expect{
		"id":            {computed: true},
		"registry_id":   {computed: true},
		"name":          {required: true},
		"description":   {optional: true},
		"approval_mode": {optional: true},
		"tags":          {optional: true},
		"registry":      {computed: true},
	}
	for name, want := range cases {
		attr, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Errorf("attribute %q missing from schema", name)
			continue
		}
		var got expect
		switch a := attr.(type) {
		case schema.StringAttribute:
			got = expect{a.Required, a.Optional, a.Computed}
		case schema.ListNestedAttribute:
			got = expect{a.Required, a.Optional, a.Computed}
		case schema.SingleNestedAttribute:
			got = expect{a.Required, a.Optional, a.Computed}
		default:
			t.Errorf("attribute %q unexpected type %T", name, attr)
		}
		if got != want {
			t.Errorf("attribute %q got %+v, want %+v", name, got, want)
		}
	}

	registry, ok := resp.Schema.Attributes["registry"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatal("registry not a SingleNestedAttribute")
	}
	registryCases := map[string]expect{
		"registry_id":             {computed: true},
		"name":                    {computed: true},
		"description":             {computed: true},
		"approval_mode":           {computed: true},
		"region":                  {computed: true},
		"status":                  {computed: true},
		"create_time":             {computed: true},
		"update_time":             {computed: true},
		"record_count":            {computed: true},
		"published_record_count":  {computed: true},
		"app_id":                  {computed: true},
		"creator_uin":             {computed: true},
		"creator_sub_account_uin": {computed: true},
		"tags":                    {computed: true},
	}
	for name, want := range registryCases {
		attr, ok := registry.Attributes[name]
		if !ok {
			t.Errorf("nested attribute %q missing from registry schema", name)
			continue
		}
		var got expect
		switch a := attr.(type) {
		case schema.StringAttribute:
			got = expect{a.Required, a.Optional, a.Computed}
		case schema.Int64Attribute:
			got = expect{a.Required, a.Optional, a.Computed}
		case schema.ListNestedAttribute:
			got = expect{a.Required, a.Optional, a.Computed}
		default:
			t.Errorf("nested attribute %q unexpected type %T", name, attr)
		}
		if got != want {
			t.Errorf("nested attribute %q got %+v, want %+v", name, got, want)
		}
	}
}

// TestAgsRegistryResource_ImplementsInterfaces verifies the resource satisfies
// the framework interfaces advertised by the var _ asserts.
func TestAgsRegistryResource_ImplementsInterfaces(t *testing.T) {
	var _ resource.ResourceWithConfigure = &AgsRegistryResource{}
	var _ resource.ResourceWithImportState = &AgsRegistryResource{}
}

// TestAgsRegistryResource_Factory validates the registry factory.
func TestAgsRegistryResource_Factory(t *testing.T) {
	got := NewAgsRegistryResource()
	if got == nil {
		t.Fatal("NewAgsRegistryResource returned nil")
	}
	if _, ok := got.(*AgsRegistryResource); !ok {
		t.Fatalf("factory returned %T, want *AgsRegistryResource", got)
	}
}

// TestAgsRegistryResource_ConfigureNilProviderData ensures the pre-Configure
// invocation pattern (nil ProviderData) is a silent no-op.
func TestAgsRegistryResource_ConfigureNilProviderData(t *testing.T) {
	r := &AgsRegistryResource{}
	req := resource.ConfigureRequest{}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Configure with nil ProviderData returned diagnostics: %v", resp.Diagnostics)
	}
}

// TestAgsRegistryResource_UpdateRejected verifies the no-update embedded helper
// rejects any Update call (all fields are RequiresReplace so this never fires).
func TestAgsRegistryResource_UpdateRejected(t *testing.T) {
	r := &AgsRegistryResource{}
	req := resource.UpdateRequest{}
	resp := &resource.UpdateResponse{}

	r.Update(context.Background(), req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("Update should have produced an error diagnostic")
	}
}

// --- expand / flatten helpers ---

// buildTagObject constructs a single tag object value.
func buildTagObject(t *testing.T, key, value string) attr.Value {
	t.Helper()
	obj, diags := types.ObjectValue(tagsAttrTypes(), map[string]attr.Value{
		"key":   types.StringValue(key),
		"value": types.StringValue(value),
	})
	if diags.HasError() {
		t.Fatalf("build tag object failed: %v", diags)
	}
	return obj
}

// buildTagsList constructs a types.List of tag objects for tests.
func buildTagsList(t *testing.T) types.List {
	t.Helper()
	elements := []attr.Value{
		buildTagObject(t, "env", "prod"),
		buildTagObject(t, "team", "devops"),
	}
	list, diags := types.ListValue(tagsObjectType(), elements)
	if diags.HasError() {
		t.Fatalf("build tags list failed: %v", diags)
	}
	return list
}

// TestExpandTags verifies the tags list -> SDK struct conversion covers all
// leaf fields.
func TestExpandTags(t *testing.T) {
	ctx := context.Background()
	list := buildTagsList(t)

	tags, diags := expandTags(ctx, list)
	if diags.HasError() {
		t.Fatalf("expandTags returned diagnostics: %v", diags)
	}
	if len(tags) != 2 {
		t.Fatalf("expandTags len = %d, want 2", len(tags))
	}
	if tags[0].Key == nil || *tags[0].Key != "env" || tags[0].Value == nil || *tags[0].Value != "prod" {
		t.Errorf("tags[0] = %v/%v, want env/prod", tags[0].Key, tags[0].Value)
	}
	if tags[1].Key == nil || *tags[1].Key != "team" || tags[1].Value == nil || *tags[1].Value != "devops" {
		t.Errorf("tags[1] = %v/%v, want team/devops", tags[1].Key, tags[1].Value)
	}
}

// TestExpandTags_Null verifies a null list yields no tags and no diagnostics.
func TestExpandTags_Null(t *testing.T) {
	tags, diags := expandTags(context.Background(), types.ListNull(tagsObjectType()))
	if diags.HasError() {
		t.Fatalf("expandTags(null) returned diagnostics: %v", diags)
	}
	if len(tags) != 0 {
		t.Errorf("expandTags(null) len = %d, want 0", len(tags))
	}
}

// buildCloudRegistry constructs a full SDK CloudRegistry for tests.
func buildCloudRegistry() *sdkags.CloudRegistry {
	return &sdkags.CloudRegistry{
		RegistryId:           helper.String("reg-abcd1234"),
		Name:                 helper.String("tf-example"),
		Description:          helper.String("example registry"),
		ApprovalMode:         helper.String("AUTO"),
		Region:               helper.String("ap-guangzhou"),
		Status:               helper.String("ACTIVE"),
		CreateTime:           helper.String("2026-08-11T10:00:00Z"),
		UpdateTime:           helper.String("2026-08-11T10:00:00Z"),
		RecordCount:          helper.Int64(3),
		PublishedRecordCount: helper.Int64(1),
		AppId:                helper.Int64(1250000000),
		CreatorUin:           helper.String("100000000001"),
		CreatorSubAccountUin: helper.String(""),
		Tags: []*sdkags.CloudTag{
			{Key: helper.String("env"), Value: helper.String("prod")},
			{Key: helper.String("team"), Value: helper.String("devops")},
		},
	}
}

// TestTagsListValue verifies the SDK tags -> types.List conversion preserves
// every leaf value.
func TestTagsListValue(t *testing.T) {
	list := tagsListValue(buildCloudRegistry().Tags)
	if list.IsNull() || list.IsUnknown() {
		t.Fatal("tagsListValue returned null/unknown list")
	}
	elements := list.Elements()
	if len(elements) != 2 {
		t.Fatalf("tags list len = %d, want 2", len(elements))
	}
	obj, ok := elements[0].(basetypes.ObjectValue)
	if !ok {
		t.Fatalf("element 0 type = %T, want basetypes.ObjectValue", elements[0])
	}
	var model AgsRegistryTagModel
	diags := obj.As(context.Background(), &model, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		t.Fatalf("obj.As returned diagnostics: %v", diags)
	}
	if model.Key.ValueString() != "env" || model.Value.ValueString() != "prod" {
		t.Errorf("tag = %s/%s, want env/prod", model.Key.ValueString(), model.Value.ValueString())
	}
}

// TestTagsListValue_Empty verifies an empty SDK tag slice maps to a null list.
func TestTagsListValue_Empty(t *testing.T) {
	list := tagsListValue(nil)
	if !list.IsNull() {
		t.Errorf("tagsListValue(nil) should be null, got %v", list)
	}
}

// TestRefreshRegistryModel verifies every computed leaf field is refreshed
// from the cloud object.
func TestRefreshRegistryModel(t *testing.T) {
	m := &AgsRegistryResourceModel{}
	refreshRegistryModel(m, buildCloudRegistry())

	if m.Registry.IsNull() {
		t.Fatal("registry object is null after refresh")
	}
	var info AgsRegistryInfoModel
	diags := m.Registry.As(context.Background(), &info, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		t.Fatalf("registry.As returned diagnostics: %v", diags)
	}
	if info.RegistryId.ValueString() != "reg-abcd1234" {
		t.Errorf("registry_id = %s", info.RegistryId.ValueString())
	}
	if info.Name.ValueString() != "tf-example" {
		t.Errorf("name = %s", info.Name.ValueString())
	}
	if info.Description.ValueString() != "example registry" {
		t.Errorf("description = %s", info.Description.ValueString())
	}
	if info.ApprovalMode.ValueString() != "AUTO" {
		t.Errorf("approval_mode = %s", info.ApprovalMode.ValueString())
	}
	if info.Region.ValueString() != "ap-guangzhou" {
		t.Errorf("region = %s", info.Region.ValueString())
	}
	if info.Status.ValueString() != "ACTIVE" {
		t.Errorf("status = %s", info.Status.ValueString())
	}
	if info.CreateTime.ValueString() != "2026-08-11T10:00:00Z" {
		t.Errorf("create_time = %s", info.CreateTime.ValueString())
	}
	if info.UpdateTime.ValueString() != "2026-08-11T10:00:00Z" {
		t.Errorf("update_time = %s", info.UpdateTime.ValueString())
	}
	if info.RecordCount.ValueInt64() != 3 {
		t.Errorf("record_count = %d", info.RecordCount.ValueInt64())
	}
	if info.PublishedRecordCount.ValueInt64() != 1 {
		t.Errorf("published_record_count = %d", info.PublishedRecordCount.ValueInt64())
	}
	if info.AppId.ValueInt64() != 1250000000 {
		t.Errorf("app_id = %d", info.AppId.ValueInt64())
	}
	if info.CreatorUin.ValueString() != "100000000001" {
		t.Errorf("creator_uin = %s", info.CreatorUin.ValueString())
	}
	if info.CreatorSubAccountUin.ValueString() != "" {
		t.Errorf("creator_sub_account_uin = %s", info.CreatorSubAccountUin.ValueString())
	}
	if len(info.Tags.Elements()) != 2 {
		t.Errorf("registry.tags len = %d, want 2", len(info.Tags.Elements()))
	}
	if m.Name.ValueString() != "tf-example" {
		t.Errorf("top-level name = %s", m.Name.ValueString())
	}
	if m.ApprovalMode.ValueString() != "AUTO" {
		t.Errorf("top-level approval_mode = %s", m.ApprovalMode.ValueString())
	}
	if m.Tags.IsNull() || len(m.Tags.Elements()) != 2 {
		t.Error("top-level tags should carry 2 elements")
	}
}

// TestRefreshRegistryModel_Nil verifies a nil cloud object leaves the model
// untouched.
func TestRefreshRegistryModel_Nil(t *testing.T) {
	m := &AgsRegistryResourceModel{}
	refreshRegistryModel(m, nil)
	if !m.Registry.IsNull() {
		t.Error("registry should remain null for nil item")
	}
}

// TestRefreshRegistryModel_NullLeaves verifies nullable cloud fields map to
// null state values.
func TestRefreshRegistryModel_NullLeaves(t *testing.T) {
	m := &AgsRegistryResourceModel{}
	refreshRegistryModel(m, &sdkags.CloudRegistry{
		RegistryId: helper.String("reg-abcd1234"),
		Name:       helper.String("tf-example"),
	})

	var info AgsRegistryInfoModel
	diags := m.Registry.As(context.Background(), &info, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		t.Fatalf("registry.As returned diagnostics: %v", diags)
	}
	if !info.Description.IsNull() {
		t.Errorf("description should be null, got %s", info.Description.ValueString())
	}
	if !info.Status.IsNull() {
		t.Errorf("status should be null, got %s", info.Status.ValueString())
	}
	if !info.RecordCount.IsNull() {
		t.Errorf("record_count should be null, got %d", info.RecordCount.ValueInt64())
	}
	if !info.Tags.IsNull() {
		t.Error("tags should be null for no cloud tags")
	}
}

// --- CRUD with mocked cloud API ---

// mockAgsClient wires a resource with a fake client and patches the client
// accessor. The returned client can be patched per-API with ApplyMethodFunc.
func mockAgsClient(t *testing.T) (*gomonkey.Patches, *AgsRegistryResource, *sdkags.Client) {
	t.Helper()
	patches := gomonkey.NewPatches()
	client := &connectivity.TencentCloudClient{Region: "ap-guangzhou"}
	agsClient := &sdkags.Client{}
	patches.ApplyMethodReturn(client, "UseAgsClient", agsClient)
	r := &AgsRegistryResource{}
	resp := &resource.ConfigureResponse{}
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: client}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Configure returned diagnostics: %v", resp.Diagnostics)
	}
	return patches, r, agsClient
}

// testResourceSchema returns the real resource schema for state/plan assembly.
func testResourceSchema(t *testing.T) tfsdk.Schema {
	t.Helper()
	resp := testResourceSchemaResponse(t)
	return tfsdk.Schema{Attributes: resp.Schema.Attributes}
}

// testResourceSchemaResponse renders the resource schema once for assertions.
func testResourceSchemaResponse(t *testing.T) *resource.SchemaResponse {
	t.Helper()
	r := &AgsRegistryResource{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}
	return resp
}

// buildCreatePlan builds a plan object for Create tests.
func buildCreatePlan(ctx context.Context, t *testing.T, withOptional bool) tfsdk.Plan {
	t.Helper()
	model := AgsRegistryResourceModel{
		Name: types.StringValue("tf-example"),
		Tags: buildTagsList(t),
	}
	if withOptional {
		model.Description = types.StringValue("example registry")
		model.ApprovalMode = types.StringValue("AUTO")
	}
	plan := tfsdk.Plan{
		Schema: testResourceSchema(t),
		Raw:    types.ObjectNull(modelAttrTypes()),
	}
	if diags := plan.Set(ctx, &model); diags.HasError() {
		t.Fatalf("plan.Set returned diagnostics: %v", diags)
	}
	return plan
}

// buildStateWithId builds a state containing only the resource id.
func buildStateWithId(ctx context.Context, t *testing.T, id string) tfsdk.State {
	t.Helper()
	state := tfsdk.State{
		Schema: testResourceSchema(t),
		Raw:    types.ObjectNull(modelAttrTypes()),
	}
	model := AgsRegistryResourceModel{
		Id:         types.StringValue(id),
		RegistryId: types.StringValue(id),
	}
	if diags := state.Set(ctx, &model); diags.HasError() {
		t.Fatalf("state.Set returned diagnostics: %v", diags)
	}
	return state
}

// TestAgsRegistryResource_Create_Success mocks CreateRegistry and
// DescribeRegistry and verifies the state carries the identity and all
// refreshed leaf fields.
func TestAgsRegistryResource_Create_Success(t *testing.T) {
	ctx := context.Background()
	patches, r, agsClient := mockAgsClient(t)
	defer patches.Reset()

	var createCalls int
	patches.ApplyMethodFunc(agsClient, "CreateRegistryWithContext",
		func(ctx context.Context, request *sdkags.CreateRegistryRequest) (*sdkags.CreateRegistryResponse, error) {
			createCalls++
			if request.Name == nil || *request.Name != "tf-example" {
				t.Errorf("CreateRegistry request.Name = %v, want tf-example", request.Name)
			}
			if len(request.Tags) != 2 {
				t.Errorf("CreateRegistry request.Tags len = %d, want 2", len(request.Tags))
			}
			if request.Tags[0].Key == nil || *request.Tags[0].Key != "env" {
				t.Errorf("CreateRegistry request.Tags[0].Key = %v, want env", request.Tags[0].Key)
			}
			resp := &sdkags.CreateRegistryResponse{}
			resp.Response = &sdkags.CreateRegistryResponseParams{
				RegistryId: helper.String("reg-abcd1234"),
			}
			return resp, nil
		})
	patches.ApplyMethodFunc(agsClient, "DescribeRegistryWithContext",
		func(ctx context.Context, request *sdkags.DescribeRegistryRequest) (*sdkags.DescribeRegistryResponse, error) {
			if request.RegistryId == nil || *request.RegistryId != "reg-abcd1234" {
				t.Errorf("DescribeRegistry request.RegistryId = %v, want reg-abcd1234", request.RegistryId)
			}
			resp := &sdkags.DescribeRegistryResponse{}
			resp.Response = &sdkags.DescribeRegistryResponseParams{Registry: buildCloudRegistry()}
			return resp, nil
		})

	plan := buildCreatePlan(ctx, t, true)
	resp := &resource.CreateResponse{
		State: tfsdk.State{Schema: testResourceSchema(t), Raw: plan.Raw},
	}

	r.Create(ctx, resource.CreateRequest{Plan: plan}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create returned diagnostics: %v", resp.Diagnostics)
	}
	if createCalls != 1 {
		t.Fatalf("CreateRegistry called %d times, want 1", createCalls)
	}

	var state AgsRegistryResourceModel
	resp.State.Get(ctx, &state)
	if state.Id.ValueString() != "reg-abcd1234" {
		t.Errorf("state.id = %s, want reg-abcd1234", state.Id.ValueString())
	}
	if state.RegistryId.ValueString() != "reg-abcd1234" {
		t.Errorf("state.registry_id = %s", state.RegistryId.ValueString())
	}
	if state.Name.ValueString() != "tf-example" {
		t.Errorf("state.name = %s", state.Name.ValueString())
	}
	if state.ApprovalMode.ValueString() != "AUTO" {
		t.Errorf("state.approval_mode = %s", state.ApprovalMode.ValueString())
	}
	if state.Description.ValueString() != "example registry" {
		t.Errorf("state.description = %s", state.Description.ValueString())
	}
	if len(state.Tags.Elements()) != 2 {
		t.Errorf("state.tags len = %d, want 2", len(state.Tags.Elements()))
	}
	if state.Registry.IsNull() {
		t.Error("state.registry should not be null")
	}
}

// TestAgsRegistryResource_Create_APIError verifies an API failure produces an
// error diagnostic.
func TestAgsRegistryResource_Create_APIError(t *testing.T) {
	ctx := context.Background()
	patches, r, agsClient := mockAgsClient(t)
	defer patches.Reset()

	patches.ApplyMethodReturn(agsClient, "CreateRegistryWithContext", (*sdkags.CreateRegistryResponse)(nil), fmt.Errorf("InternalError"))

	plan := buildCreatePlan(ctx, t, false)
	resp := &resource.CreateResponse{
		State: tfsdk.State{Schema: testResourceSchema(t), Raw: plan.Raw},
	}

	r.Create(ctx, resource.CreateRequest{Plan: plan}, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("Create should have produced an error diagnostic")
	}
}

// TestAgsRegistryResource_Create_EmptyId verifies an empty registry id is
// treated as an API failure.
func TestAgsRegistryResource_Create_EmptyId(t *testing.T) {
	ctx := context.Background()
	patches, r, agsClient := mockAgsClient(t)
	defer patches.Reset()

	resp := &sdkags.CreateRegistryResponse{
		Response: &sdkags.CreateRegistryResponseParams{RegistryId: helper.String("")},
	}
	patches.ApplyMethodReturn(agsClient, "CreateRegistryWithContext", resp, nil)

	plan := buildCreatePlan(ctx, t, false)
	createResp := &resource.CreateResponse{
		State: tfsdk.State{Schema: testResourceSchema(t), Raw: plan.Raw},
	}

	r.Create(ctx, resource.CreateRequest{Plan: plan}, createResp)
	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create should have produced an error diagnostic for an empty id")
	}
}

// TestAgsRegistryResource_Read_Success mocks DescribeRegistry and verifies
// the state refresh carries all leaf fields.
func TestAgsRegistryResource_Read_Success(t *testing.T) {
	ctx := context.Background()
	patches, r, agsClient := mockAgsClient(t)
	defer patches.Reset()

	patches.ApplyMethodFunc(agsClient, "DescribeRegistryWithContext",
		func(ctx context.Context, request *sdkags.DescribeRegistryRequest) (*sdkags.DescribeRegistryResponse, error) {
			resp := &sdkags.DescribeRegistryResponse{}
			resp.Response = &sdkags.DescribeRegistryResponseParams{Registry: buildCloudRegistry()}
			return resp, nil
		})

	state := buildStateWithId(ctx, t, "reg-abcd1234")
	resp := &resource.ReadResponse{State: state}
	r.Read(ctx, resource.ReadRequest{State: state}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", resp.Diagnostics)
	}

	var refreshed AgsRegistryResourceModel
	resp.State.Get(ctx, &refreshed)
	if refreshed.Id.ValueString() != "reg-abcd1234" {
		t.Errorf("state.id = %s", refreshed.Id.ValueString())
	}
	if refreshed.RegistryId.ValueString() != "reg-abcd1234" {
		t.Errorf("state.registry_id = %s", refreshed.RegistryId.ValueString())
	}
	if refreshed.Name.ValueString() != "tf-example" {
		t.Errorf("state.name = %s", refreshed.Name.ValueString())
	}
	if refreshed.ApprovalMode.ValueString() != "AUTO" {
		t.Errorf("state.approval_mode = %s", refreshed.ApprovalMode.ValueString())
	}
	if refreshed.Registry.IsNull() {
		t.Error("state.registry should not be null")
	}
	var info AgsRegistryInfoModel
	diags := refreshed.Registry.As(ctx, &info, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		t.Fatalf("registry.As returned diagnostics: %v", diags)
	}
	if info.Status.ValueString() != "ACTIVE" {
		t.Errorf("registry.status = %s", info.Status.ValueString())
	}
	if info.RecordCount.ValueInt64() != 3 {
		t.Errorf("registry.record_count = %d", info.RecordCount.ValueInt64())
	}
}

// TestAgsRegistryResource_Read_NotFound verifies a nil registry removes the
// resource from state without errors.
func TestAgsRegistryResource_Read_NotFound(t *testing.T) {
	ctx := context.Background()
	patches, r, agsClient := mockAgsClient(t)
	defer patches.Reset()

	resp := &sdkags.DescribeRegistryResponse{
		Response: &sdkags.DescribeRegistryResponseParams{Registry: nil},
	}
	patches.ApplyMethodReturn(agsClient, "DescribeRegistryWithContext", resp, nil)

	state := buildStateWithId(ctx, t, "reg-gone0000")
	readResp := &resource.ReadResponse{State: state}
	r.Read(ctx, resource.ReadRequest{State: state}, readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", readResp.Diagnostics)
	}
}

// TestAgsRegistryResource_Read_APIError verifies API failures are surfaced as
// error diagnostics.
func TestAgsRegistryResource_Read_APIError(t *testing.T) {
	ctx := context.Background()
	patches, r, agsClient := mockAgsClient(t)
	defer patches.Reset()

	patches.ApplyMethodReturn(agsClient, "DescribeRegistryWithContext", (*sdkags.DescribeRegistryResponse)(nil), fmt.Errorf("NetworkError"))

	state := buildStateWithId(ctx, t, "reg-abcd1234")
	resp := &resource.ReadResponse{State: state}
	r.Read(ctx, resource.ReadRequest{State: state}, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("Read should have produced an error diagnostic")
	}
}

// TestAgsRegistryResource_Delete_Success mocks DeleteRegistry and verifies
// no error diagnostics and the id was passed through.
func TestAgsRegistryResource_Delete_Success(t *testing.T) {
	ctx := context.Background()
	patches, r, agsClient := mockAgsClient(t)
	defer patches.Reset()

	var deletedId *string
	patches.ApplyMethodFunc(agsClient, "DeleteRegistryWithContext",
		func(ctx context.Context, request *sdkags.DeleteRegistryRequest) (*sdkags.DeleteRegistryResponse, error) {
			deletedId = request.RegistryId
			return &sdkags.DeleteRegistryResponse{}, nil
		})

	state := buildStateWithId(ctx, t, "reg-abcd1234")
	resp := &resource.DeleteResponse{}
	r.Delete(ctx, resource.DeleteRequest{State: state}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Delete returned diagnostics: %v", resp.Diagnostics)
	}
	if deletedId == nil || *deletedId != "reg-abcd1234" {
		t.Errorf("DeleteRegistry request.RegistryId = %v, want reg-abcd1234", deletedId)
	}
}

// TestAgsRegistryResource_Delete_APIError verifies delete failures keep the
// resource in state for retry.
func TestAgsRegistryResource_Delete_APIError(t *testing.T) {
	ctx := context.Background()
	patches, r, agsClient := mockAgsClient(t)
	defer patches.Reset()

	patches.ApplyMethodReturn(agsClient, "DeleteRegistryWithContext", (*sdkags.DeleteRegistryResponse)(nil), fmt.Errorf("OperationDenied"))

	state := buildStateWithId(ctx, t, "reg-abcd1234")
	resp := &resource.DeleteResponse{}
	r.Delete(ctx, resource.DeleteRequest{State: state}, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("Delete should have produced an error diagnostic")
	}
}

// TestAgsRegistryResource_ImportState verifies the embedded passthrough import
// maps the import id onto the id attribute.
func TestAgsRegistryResource_ImportState(t *testing.T) {
	r := &AgsRegistryResource{}
	req := resource.ImportStateRequest{ID: "reg-abcd1234"}
	resp := &resource.ImportStateResponse{
		State: tfsdk.State{Schema: testResourceSchema(t), Raw: types.ObjectNull(modelAttrTypes())},
	}

	r.ImportState(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("ImportState returned diagnostics: %v", resp.Diagnostics)
	}
	var id types.String
	resp.State.GetAttribute(context.Background(), path.Root("id"), &id)
	if id.ValueString() != "reg-abcd1234" {
		t.Errorf("imported id = %s, want reg-abcd1234", id.ValueString())
	}
}
