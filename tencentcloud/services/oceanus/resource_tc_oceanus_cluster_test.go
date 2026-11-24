package oceanus

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"
)

// TestOceanusClusterResource_Metadata validates the resource type name.
func TestOceanusClusterResource_Metadata(t *testing.T) {
	r := &OceanusClusterResource{}
	req := resource.MetadataRequest{}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	if resp.TypeName != "tencentcloud_oceanus_cluster" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_oceanus_cluster")
	}
}

// TestOceanusClusterResource_Schema validates required / optional / computed
// flags on the top-level attributes.
func TestOceanusClusterResource_Schema(t *testing.T) {
	r := &OceanusClusterResource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}

	type expect struct {
		required, optional, computed, sensitive bool
	}
	cases := map[string]expect{
		"id":                               {computed: true},
		"cluster_name":                     {required: true},
		"region_id":                        {required: true},
		"zone_id":                          {required: true},
		"login_password":                   {required: true, sensitive: true},
		"vpc_descriptions":                 {required: true},
		"default_cos_bucket":               {required: true},
		"cu":                               {optional: true, computed: true},
		"remark":                           {optional: true, computed: true},
		"period":                           {optional: true},
		"instance_charge_type":             {optional: true},
		"cluster_type":                     {optional: true},
		"renew_flag":                       {optional: true},
		"flink_web_ui_network_access_type": {optional: true},
		"slave_vpc_descriptions":           {optional: true},
		"cu_memory":                        {optional: true, computed: true},
		"status":                           {computed: true},
		"status_desc":                      {computed: true},
		"version":                          {computed: true},
		"tags":                             {computed: true},
		"cc_ns":                            {computed: true},
		"slave_zones":                      {computed: true},
		"security_group_ids":               {computed: true},
		"create_time":                      {computed: true},
		"update_time":                      {computed: true},
	}

	for name, want := range cases {
		attr, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Errorf("attribute %q missing from schema", name)
			continue
		}

		required := attr.IsRequired()
		optional := attr.IsOptional()
		computed := attr.IsComputed()

		if required != want.required {
			t.Errorf("attribute %q Required=%v, want %v", name, required, want.required)
		}
		if optional != want.optional {
			t.Errorf("attribute %q Optional=%v, want %v", name, optional, want.optional)
		}
		if computed != want.computed {
			t.Errorf("attribute %q Computed=%v, want %v", name, computed, want.computed)
		}

		if stringAttr, ok := attr.(schema.StringAttribute); ok && stringAttr.Sensitive != want.sensitive {
			t.Errorf("attribute %q Sensitive=%v, want %v", name, stringAttr.Sensitive, want.sensitive)
		}
	}
}

// TestOceanusClusterResource_Factory validates the registry factory.
func TestOceanusClusterResource_Factory(t *testing.T) {
	got := NewOceanusClusterResource()
	if got == nil {
		t.Fatal("NewOceanusClusterResource returned nil")
	}
	if _, ok := got.(*OceanusClusterResource); !ok {
		t.Fatalf("factory returned %T, want *OceanusClusterResource", got)
	}
}

// TestOceanusClusterResource_ConfigureNilProviderData ensures the
// pre-Configure invocation pattern (nil ProviderData) is a silent no-op.
func TestOceanusClusterResource_ConfigureNilProviderData(t *testing.T) {
	r := &OceanusClusterResource{}
	req := resource.ConfigureRequest{}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Configure with nil ProviderData returned diagnostics: %v", resp.Diagnostics)
	}

	if r.Client() != nil {
		t.Error("Client() should stay nil when ProviderData is nil")
	}
}

// TestExpandOceanusClusterVpcDescriptions validates the vpc_descriptions
// expansion including the optional app_id / owner_uin fields.
func TestExpandOceanusClusterVpcDescriptions(t *testing.T) {
	list, diags := types.ListValueFrom(context.Background(), oceanusClusterVpcDescriptionObjectType(t), []oceanusClusterVpcDescriptionModel{
		{
			VpcId:    types.StringValue("vpc-123"),
			SubnetId: types.StringValue("subnet-123"),
			AppId:    types.Int64Null(),
			OwnerUin: types.StringNull(),
		},
	})
	if diags.HasError() {
		t.Fatalf("ListValueFrom returned diagnostics: %v", diags)
	}

	out, expandDiags := expandOceanusClusterVpcDescriptions(context.Background(), list)
	if expandDiags.HasError() {
		t.Fatalf("expandOceanusClusterVpcDescriptions returned diagnostics: %v", expandDiags)
	}

	if len(out) != 1 {
		t.Fatalf("len(out) = %d, want 1", len(out))
	}
	if out[0].VpcId == nil || *out[0].VpcId != "vpc-123" {
		t.Errorf("VpcId = %v, want vpc-123", out[0].VpcId)
	}
	if out[0].SubnetId == nil || *out[0].SubnetId != "subnet-123" {
		t.Errorf("SubnetId = %v, want subnet-123", out[0].SubnetId)
	}
	if out[0].AppId != nil {
		t.Errorf("AppId = %v, want nil when unset", out[0].AppId)
	}
	if out[0].OwnerUin != nil {
		t.Errorf("OwnerUin = %v, want nil when unset", out[0].OwnerUin)
	}
}

// TestExpandOceanusClusterVpcDescriptionsNull validates that an absent list
// expands to nothing.
func TestExpandOceanusClusterVpcDescriptionsNull(t *testing.T) {
	out, diags := expandOceanusClusterVpcDescriptions(context.Background(), types.ListNull(oceanusClusterVpcDescriptionObjectType(t)))
	if diags.HasError() {
		t.Fatalf("expandOceanusClusterVpcDescriptions returned diagnostics: %v", diags)
	}
	if out != nil {
		t.Errorf("out = %v, want nil", out)
	}
}

// TestFloat64ValueOrNull validates the float conversion helper.
func TestFloat64ValueOrNull(t *testing.T) {
	if v := float64ValueOrNull(nil); !v.IsNull() {
		t.Errorf("float64ValueOrNull(nil) = %v, want Null", v)
	}

	value := 12.5
	if v := float64ValueOrNull(&value); v.IsNull() || v.ValueFloat64() != 12.5 {
		t.Errorf("float64ValueOrNull(12.5) = %v, want 12.5", v)
	}
}

// TestStringSetOrNull validates the string set conversion helper.
func TestStringSetOrNull(t *testing.T) {
	if v := stringSetOrNull(nil); !v.IsNull() {
		t.Errorf("stringSetOrNull(nil) = %v, want Null", v)
	}

	first, second := "a", "b"
	set := stringSetOrNull([]*string{&first, nil, &second})
	if set.IsNull() {
		t.Fatal("stringSetOrNull returned Null, want 2 elements")
	}
	if len(set.Elements()) != 2 {
		t.Errorf("len(set.Elements()) = %d, want 2", len(set.Elements()))
	}
}

// TestRefreshOceanusClusterModelKeepsUnreadableFields checks that the fields
// the Describe API does not echo back keep their previous values.
func TestRefreshOceanusClusterModelKeepsUnreadableFields(t *testing.T) {
	model := OceanusClusterResourceModel{
		ClusterName:                 types.StringValue("tf-example"),
		RegionId:                    types.Int64Value(1),
		ZoneId:                      types.Int64Value(100003),
		LoginPassword:               types.StringValue("Password@123"),
		VpcDescriptions:             types.ListNull(oceanusClusterVpcDescriptionObjectType(t)),
		DefaultCosBucket:            types.StringValue("bucket-1250000000"),
		Cu:                          types.Int64Value(12),
		Remark:                      types.StringValue("remark"),
		FlinkWebUiNetworkAccessType: types.StringValue("NetworkAccess_EXTERNAL"),
		SlaveVpcDescriptions:        types.ListNull(oceanusClusterVpcDescriptionObjectType(t)),
		CuMemory:                    types.Int64Value(4),
	}

	name := "tf-example"
	diags := refreshOceanusClusterModel(context.Background(), &model, &oceanus.Cluster{
		ClusterId: nil,
		Name:      &name,
	})
	if diags.HasError() {
		t.Fatalf("refreshOceanusClusterModel returned diagnostics: %v", diags)
	}

	if model.ClusterName.ValueString() != "tf-example" {
		t.Errorf("ClusterName = %q, want tf-example", model.ClusterName.ValueString())
	}
	if model.RegionId.ValueInt64() != 1 || model.ZoneId.ValueInt64() != 100003 {
		t.Errorf("RegionId/ZoneId changed to %v/%v", model.RegionId, model.ZoneId)
	}
	if model.LoginPassword.ValueString() != "Password@123" {
		t.Errorf("LoginPassword changed to %q", model.LoginPassword.ValueString())
	}
	if model.DefaultCosBucket.ValueString() != "bucket-1250000000" {
		t.Errorf("DefaultCosBucket changed to %q", model.DefaultCosBucket.ValueString())
	}
	if model.CuMemory.ValueInt64() != 4 {
		t.Errorf("CuMemory changed to %v", model.CuMemory)
	}
	if model.FlinkWebUiNetworkAccessType.ValueString() != "NetworkAccess_EXTERNAL" {
		t.Errorf("FlinkWebUiNetworkAccessType changed to %q", model.FlinkWebUiNetworkAccessType.ValueString())
	}
	if !model.Version.IsNull() {
		t.Errorf("Version = %v, want Null when the cloud omits it", model.Version)
	}
}

// oceanusClusterVpcDescriptionObjectType mirrors the vpc_descriptions nested
// object type for the tests above.
func oceanusClusterVpcDescriptionObjectType(t *testing.T) types.ObjectType {
	t.Helper()

	r := &OceanusClusterResource{}
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}

	attr, ok := resp.Schema.Attributes["vpc_descriptions"].(schema.ListNestedAttribute)
	if !ok {
		t.Fatal("vpc_descriptions is not a ListNestedAttribute")
	}

	attrTypes := make(map[string]attr.Type)
	for name, nestedAttr := range attr.NestedObject.GetAttributes() {
		attrTypes[name] = nestedAttr.GetType()
	}

	return types.ObjectType{AttrTypes: attrTypes}
}
