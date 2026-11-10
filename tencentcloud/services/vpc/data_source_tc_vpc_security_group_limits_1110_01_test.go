package vpc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"
)

// go test ./tencentcloud/services/vpc/ -run "TestVpcSecurityGroupLimits111001" -v -count=1 -gcflags="all=-l"

const vpcSecurityGroupLimits111001TypeName = "tencentcloud_vpc_security_group_limits_1110_01"

// vpcSecurityGroupLimits111001LeafAttrs lists the leaf attributes of
// security_group_limit_set, in schema order.
var vpcSecurityGroupLimits111001LeafAttrs = []string{
	"security_group_limit",
	"security_group_policy_limit",
	"refered_security_group_limit",
	"security_group_instance_limit",
	"instance_security_group_limit",
	"security_group_extended_policy_limit",
	"security_group_refered_cvm_and_eni_limit",
	"security_group_refered_svc_limit",
}

// newVpcSecurityGroupLimits111001DataSource builds the data source under test
// and returns it together with its schema.
func newVpcSecurityGroupLimits111001DataSource(t *testing.T) (datasource.DataSource, schema.Schema) {
	t.Helper()

	ds := svcvpc.NewVpcSecurityGroupLimits111001DataSource()
	if ds == nil {
		t.Fatal("NewVpcSecurityGroupLimits111001DataSource returned nil")
	}

	schemaResp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", schemaResp.Diagnostics)
	}

	listAttr, ok := schemaResp.Schema.Attributes["security_group_limit_set"].(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("security_group_limit_set has type %T, want schema.ListNestedAttribute",
			schemaResp.Schema.Attributes["security_group_limit_set"])
	}
	if listAttr.NestedObject.Attributes == nil {
		t.Fatal("security_group_limit_set is missing its nested attributes")
	}

	return ds, schemaResp.Schema
}

// configureVpcSecurityGroupLimits111001DataSource injects a mock connectivity
// client through the ProviderData hand-off used by the framework provider.
func configureVpcSecurityGroupLimits111001DataSource(t *testing.T, ds datasource.DataSource, client *connectivity.TencentCloudClient) {
	t.Helper()

	configurable, ok := ds.(datasource.DataSourceWithConfigure)
	if !ok {
		t.Fatalf("data source %T does not implement datasource.DataSourceWithConfigure", ds)
	}
	configResp := &datasource.ConfigureResponse{}
	configurable.Configure(
		context.Background(),
		datasource.ConfigureRequest{ProviderData: &sharedmeta.ProviderMeta{Client: client}},
		configResp,
	)
	if configResp.Diagnostics.HasError() {
		t.Fatalf("Configure returned diagnostics: %v", configResp.Diagnostics)
	}
}

// newVpcSecurityGroupLimits111001ReadResponse builds a ReadResponse whose State
// is initialised with an empty object of the data source type, which is what
// the framework hands over before Read populates it.
func newVpcSecurityGroupLimits111001ReadResponse(t *testing.T, s schema.Schema) *datasource.ReadResponse {
	t.Helper()

	objType, ok := s.Type().(types.ObjectType)
	if !ok {
		t.Fatalf("data source schema type is %T, want types.ObjectType", s.Type())
	}

	attrs := make(map[string]attr.Value, len(objType.AttrTypes))
	for name, typ := range objType.AttrTypes {
		attrs[name] = nullValueOf(t, typ)
	}
	obj, diags := types.ObjectValue(objType.AttrTypes, attrs)
	if diags.HasError() {
		t.Fatalf("building empty state object failed: %v", diags)
	}

	resp := &datasource.ReadResponse{}
	resp.Diagnostics.Append(resp.State.Set(context.Background(), obj)...)
	if resp.Diagnostics.HasError() {
		t.Fatalf("initialising state failed: %v", resp.Diagnostics)
	}
	return resp
}

// nullValueOf returns the null value of the given attribute type.
func nullValueOf(t *testing.T, typ attr.Type) attr.Value {
	t.Helper()

	switch typ := typ.(type) {
	case types.ListType:
		return types.ListNull(typ.ElemType)
	default:
		t.Fatalf("unsupported attribute type %T in test helper", typ)
		return nil
	}
}

// readVpcSecurityGroupLimits111001Model decodes the whole state into the data
// source model.
func readVpcSecurityGroupLimits111001Model(t *testing.T, resp *datasource.ReadResponse) svcvpc.VpcSecurityGroupLimits111001DataSourceModel {
	t.Helper()

	var model svcvpc.VpcSecurityGroupLimits111001DataSourceModel
	diags := resp.State.Get(context.Background(), &model)
	if diags.HasError() {
		t.Fatalf("reading state failed: %v", diags)
	}
	if model.SecurityGroupLimitSet.IsNull() {
		t.Fatal("security_group_limit_set is null in state, want a one element list")
	}
	if model.SecurityGroupLimitSet.IsUnknown() {
		t.Fatal("security_group_limit_set is unknown in state, want a known value")
	}
	return model
}

// readVpcSecurityGroupLimits111001Item decodes security_group_limit_set[0] into
// a map keyed by attribute name.
func readVpcSecurityGroupLimits111001Item(t *testing.T, model svcvpc.VpcSecurityGroupLimits111001DataSourceModel) map[string]types.Int64 {
	t.Helper()

	elements := model.SecurityGroupLimitSet.Elements()
	if len(elements) != 1 {
		t.Fatalf("security_group_limit_set has %d elements in state, want 1", len(elements))
	}
	object, ok := elements[0].(types.Object)
	if !ok {
		t.Fatalf("security_group_limit_set element has type %T, want types.Object", elements[0])
	}

	attrs := object.Attributes()
	out := make(map[string]types.Int64, len(attrs))
	for name, value := range attrs {
		intValue, ok := value.(types.Int64)
		if !ok {
			t.Fatalf("attribute %s has type %T, want types.Int64", name, value)
		}
		out[name] = intValue
	}
	return out
}

// newVpcSecurityGroupLimits111001LimitSet builds a fully populated API payload.
func newVpcSecurityGroupLimits111001LimitSet() *vpc.SecurityGroupLimitSet {
	return &vpc.SecurityGroupLimitSet{
		SecurityGroupLimit:                 common.Uint64Ptr(20),
		SecurityGroupPolicyLimit:           common.Uint64Ptr(100),
		ReferedSecurityGroupLimit:          common.Uint64Ptr(10),
		SecurityGroupInstanceLimit:         common.Uint64Ptr(2000),
		InstanceSecurityGroupLimit:         common.Uint64Ptr(5),
		SecurityGroupExtendedPolicyLimit:   common.Uint64Ptr(180),
		SecurityGroupReferedCvmAndEniLimit: common.Uint64Ptr(1),
		SecurityGroupReferedSvcLimit:       common.Uint64Ptr(2),
	}
}

// TestVpcSecurityGroupLimits111001_DataSource_Metadata validates the registered
// Terraform type name.
func TestVpcSecurityGroupLimits111001_DataSource_Metadata(t *testing.T) {
	ds := svcvpc.NewVpcSecurityGroupLimits111001DataSource()

	resp := &datasource.MetadataResponse{}
	ds.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "tencentcloud"}, resp)

	if resp.TypeName != vpcSecurityGroupLimits111001TypeName {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, vpcSecurityGroupLimits111001TypeName)
	}
}

// TestVpcSecurityGroupLimits111001_DataSource_Interfaces validates that the
// data source exposes the data source lifecycle and accepts a configured
// client.
func TestVpcSecurityGroupLimits111001_DataSource_Interfaces(t *testing.T) {
	ds := svcvpc.NewVpcSecurityGroupLimits111001DataSource()

	if _, ok := ds.(datasource.DataSource); !ok {
		t.Fatalf("data source %T does not implement datasource.DataSource", ds)
	}
	if _, ok := ds.(datasource.DataSourceWithConfigure); !ok {
		t.Fatalf("data source %T does not implement datasource.DataSourceWithConfigure", ds)
	}
}

// TestVpcSecurityGroupLimits111001_DataSource_Schema validates that the schema
// exposes exactly one computed output list with all documented leaf attributes.
func TestVpcSecurityGroupLimits111001_DataSource_Schema(t *testing.T) {
	_, s := newVpcSecurityGroupLimits111001DataSource(t)

	if len(s.Attributes) != 1 {
		t.Fatalf("schema has %d attributes, want 1", len(s.Attributes))
	}

	limitSetAttr, ok := s.Attributes["security_group_limit_set"]
	if !ok {
		t.Fatal("schema is missing attribute security_group_limit_set")
	}
	listAttr, ok := limitSetAttr.(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("security_group_limit_set has type %T, want schema.ListNestedAttribute", limitSetAttr)
	}
	if !listAttr.Computed {
		t.Error("security_group_limit_set should be Computed")
	}
	if listAttr.Optional || listAttr.Required {
		t.Error("security_group_limit_set is an output only attribute, it should be neither Optional nor Required")
	}
	if listAttr.Description == "" {
		t.Error("security_group_limit_set is missing its Description")
	}

	if len(listAttr.NestedObject.Attributes) != len(vpcSecurityGroupLimits111001LeafAttrs) {
		t.Fatalf("security_group_limit_set has %d nested attributes, want %d",
			len(listAttr.NestedObject.Attributes), len(vpcSecurityGroupLimits111001LeafAttrs))
	}
	for _, name := range vpcSecurityGroupLimits111001LeafAttrs {
		nested, ok := listAttr.NestedObject.Attributes[name]
		if !ok {
			t.Errorf("security_group_limit_set is missing nested attribute %s", name)
			continue
		}
		intAttr, ok := nested.(schema.Int64Attribute)
		if !ok {
			t.Errorf("nested attribute %s has type %T, want schema.Int64Attribute", name, nested)
			continue
		}
		if !intAttr.Computed {
			t.Errorf("nested attribute %s should be Computed", name)
		}
		if intAttr.Optional || intAttr.Required {
			t.Errorf("nested attribute %s is an output only attribute, it should be neither Optional nor Required", name)
		}
		if intAttr.Description == "" {
			t.Errorf("nested attribute %s is missing its Description", name)
		}
	}
}

// TestVpcSecurityGroupLimits111001_DataSource_ReadSuccess validates that every
// leaf field returned by the cloud API is written back into state.
func TestVpcSecurityGroupLimits111001_DataSource_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(client, "UseVpcClient", vpcClient)

	calls := 0
	patches.ApplyMethodFunc(vpcClient, "DescribeSecurityGroupLimitsWithContext",
		func(_ context.Context, request *vpc.DescribeSecurityGroupLimitsRequest) (*vpc.DescribeSecurityGroupLimitsResponse, error) {
			calls++
			if request == nil {
				t.Fatal("DescribeSecurityGroupLimits was called with a nil request")
			}
			if request.GetAction() != "DescribeSecurityGroupLimits" {
				t.Errorf("request action = %q, want %q", request.GetAction(), "DescribeSecurityGroupLimits")
			}
			return &vpc.DescribeSecurityGroupLimitsResponse{
				Response: &vpc.DescribeSecurityGroupLimitsResponseParams{
					SecurityGroupLimitSet: newVpcSecurityGroupLimits111001LimitSet(),
					RequestId:             common.StringPtr("fake-request-id"),
				},
			}, nil
		})

	ds, s := newVpcSecurityGroupLimits111001DataSource(t)
	configureVpcSecurityGroupLimits111001DataSource(t, ds, client)

	resp := newVpcSecurityGroupLimits111001ReadResponse(t, s)
	ds.Read(context.Background(), datasource.ReadRequest{}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", resp.Diagnostics)
	}
	if calls != 1 {
		t.Errorf("DescribeSecurityGroupLimits called %d times, want 1", calls)
	}

	item := readVpcSecurityGroupLimits111001Item(t, readVpcSecurityGroupLimits111001Model(t, resp))
	expectations := map[string]int64{
		"security_group_limit":                     20,
		"security_group_policy_limit":              100,
		"refered_security_group_limit":             10,
		"security_group_instance_limit":            2000,
		"instance_security_group_limit":            5,
		"security_group_extended_policy_limit":     180,
		"security_group_refered_cvm_and_eni_limit": 1,
		"security_group_refered_svc_limit":         2,
	}
	for name, want := range expectations {
		value, ok := item[name]
		if !ok {
			t.Errorf("attribute %s is missing from state", name)
			continue
		}
		if value.IsNull() || value.IsUnknown() {
			t.Errorf("attribute %s is null or unknown in state, want %d", name, want)
			continue
		}
		if value.ValueInt64() != want {
			t.Errorf("attribute %s = %d, want %d", name, value.ValueInt64(), want)
		}
	}
}

// TestVpcSecurityGroupLimits111001_DataSource_ReadNilFields validates that a
// partially populated API response keeps the missing leaves null instead of
// turning them into zero values.
func TestVpcSecurityGroupLimits111001_DataSource_ReadNilFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(client, "UseVpcClient", vpcClient)

	patches.ApplyMethodFunc(vpcClient, "DescribeSecurityGroupLimitsWithContext",
		func(_ context.Context, _ *vpc.DescribeSecurityGroupLimitsRequest) (*vpc.DescribeSecurityGroupLimitsResponse, error) {
			return &vpc.DescribeSecurityGroupLimitsResponse{
				Response: &vpc.DescribeSecurityGroupLimitsResponseParams{
					SecurityGroupLimitSet: &vpc.SecurityGroupLimitSet{
						SecurityGroupLimit:       common.Uint64Ptr(20),
						SecurityGroupPolicyLimit: common.Uint64Ptr(100),
					},
					RequestId: common.StringPtr("fake-request-id"),
				},
			}, nil
		})

	ds, s := newVpcSecurityGroupLimits111001DataSource(t)
	configureVpcSecurityGroupLimits111001DataSource(t, ds, client)

	resp := newVpcSecurityGroupLimits111001ReadResponse(t, s)
	ds.Read(context.Background(), datasource.ReadRequest{}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", resp.Diagnostics)
	}

	item := readVpcSecurityGroupLimits111001Item(t, readVpcSecurityGroupLimits111001Model(t, resp))

	populated := map[string]int64{
		"security_group_limit":        20,
		"security_group_policy_limit": 100,
	}
	for name, want := range populated {
		value, ok := item[name]
		if !ok {
			t.Errorf("attribute %s is missing from state", name)
			continue
		}
		if value.IsNull() || value.IsUnknown() {
			t.Errorf("attribute %s is null or unknown in state, want %d", name, want)
			continue
		}
		if value.ValueInt64() != want {
			t.Errorf("attribute %s = %d, want %d", name, value.ValueInt64(), want)
		}
	}

	empty := []string{
		"refered_security_group_limit",
		"security_group_instance_limit",
		"instance_security_group_limit",
		"security_group_extended_policy_limit",
		"security_group_refered_cvm_and_eni_limit",
		"security_group_refered_svc_limit",
	}
	for _, name := range empty {
		value, ok := item[name]
		if !ok {
			t.Errorf("attribute %s is missing from state", name)
			continue
		}
		if !value.IsNull() {
			t.Errorf("attribute %s = %d, want null because the API did not return it", name, value.ValueInt64())
		}
	}
}

// TestVpcSecurityGroupLimits111001_DataSource_ReadNilLimitSet validates that a
// response without the quota payload is reported as an error instead of being
// written to state as an empty list.
func TestVpcSecurityGroupLimits111001_DataSource_ReadNilLimitSet(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(client, "UseVpcClient", vpcClient)

	patches.ApplyMethodFunc(vpcClient, "DescribeSecurityGroupLimitsWithContext",
		func(_ context.Context, _ *vpc.DescribeSecurityGroupLimitsRequest) (*vpc.DescribeSecurityGroupLimitsResponse, error) {
			return &vpc.DescribeSecurityGroupLimitsResponse{
				Response: &vpc.DescribeSecurityGroupLimitsResponseParams{
					RequestId: common.StringPtr("fake-request-id"),
				},
			}, nil
		})

	ds, s := newVpcSecurityGroupLimits111001DataSource(t)
	configureVpcSecurityGroupLimits111001DataSource(t, ds, client)

	resp := newVpcSecurityGroupLimits111001ReadResponse(t, s)
	ds.Read(context.Background(), datasource.ReadRequest{}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("Read should report an error when SecurityGroupLimitSet is missing")
	}
}

// TestVpcSecurityGroupLimits111001_DataSource_ReadNilResponse validates that a
// nil response envelope is reported as an error.
func TestVpcSecurityGroupLimits111001_DataSource_ReadNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(client, "UseVpcClient", vpcClient)

	patches.ApplyMethodFunc(vpcClient, "DescribeSecurityGroupLimitsWithContext",
		func(_ context.Context, _ *vpc.DescribeSecurityGroupLimitsRequest) (*vpc.DescribeSecurityGroupLimitsResponse, error) {
			return &vpc.DescribeSecurityGroupLimitsResponse{}, nil
		})

	ds, s := newVpcSecurityGroupLimits111001DataSource(t)
	configureVpcSecurityGroupLimits111001DataSource(t, ds, client)

	resp := newVpcSecurityGroupLimits111001ReadResponse(t, s)
	ds.Read(context.Background(), datasource.ReadRequest{}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("Read should report an error when the response envelope is nil")
	}
}

// TestVpcSecurityGroupLimits111001_DataSource_ReadAPIError validates that a
// cloud API failure is surfaced as an error diagnostic.
func TestVpcSecurityGroupLimits111001_DataSource_ReadAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(client, "UseVpcClient", vpcClient)

	patches.ApplyMethodFunc(vpcClient, "DescribeSecurityGroupLimitsWithContext",
		func(_ context.Context, _ *vpc.DescribeSecurityGroupLimitsRequest) (*vpc.DescribeSecurityGroupLimitsResponse, error) {
			return nil, errors.New("AuthFailure.UnauthorizedOperation")
		})

	ds, s := newVpcSecurityGroupLimits111001DataSource(t)
	configureVpcSecurityGroupLimits111001DataSource(t, ds, client)

	resp := newVpcSecurityGroupLimits111001ReadResponse(t, s)
	ds.Read(context.Background(), datasource.ReadRequest{}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("Read should report an error when the cloud API fails")
	}
}
