package fwm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
)

func ptrStringFwmOrganSummary(s string) *string {
	return &s
}

func ptrInt64FwmOrganSummary(i int64) *int64 {
	return &i
}

// fwmOrganSummarySchema returns the data source schema used to build the raw
// config/state values for ReadRequest and ReadResponse.
func fwmOrganSummarySchema(t *testing.T) schema.Schema {
	t.Helper()
	d := fwm.NewFwmOrganSummaryDataSource()
	schemaResp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// fwmOrganSummaryNullRawValue builds a raw tftypes value for the whole schema
// with every attribute null, which is the shape of an empty config and of an
// empty state before Read runs.
func fwmOrganSummaryNullRawValue(t *testing.T, s schema.Schema) tftypes.Value {
	t.Helper()
	return tftypes.NewValue(s.Type().TerraformType(context.Background()), nil)
}

// newFwmOrganSummaryReadRequest builds a ReadRequest whose config carries no
// input, matching a schema that has only computed attributes.
func newFwmOrganSummaryReadRequest(t *testing.T, s schema.Schema) datasource.ReadRequest {
	t.Helper()
	req := datasource.ReadRequest{}
	req.Config.Raw = fwmOrganSummaryNullRawValue(t, s)
	req.Config.Schema = s
	return req
}

// newFwmOrganSummaryReadResponse builds a ReadResponse with an empty state
// matching the schema.
func newFwmOrganSummaryReadResponse(t *testing.T, s schema.Schema) datasource.ReadResponse {
	t.Helper()
	resp := datasource.ReadResponse{}
	resp.State.Raw = fwmOrganSummaryNullRawValue(t, s)
	resp.State.Schema = s
	return resp
}

// setFwmOrganSummaryClient injects a mock connectivity client into the data
// source via the Configure method.
func setFwmOrganSummaryClient(d datasource.DataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

// readFwmOrganSummaryState extracts the summary object from the state that
// Read wrote, failing the test when Read reported an error.
func readFwmOrganSummaryState(t *testing.T, resp *datasource.ReadResponse) types.Object {
	t.Helper()
	assert.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)

	var summary types.Object
	diags := resp.State.GetAttribute(context.Background(), path.Root("summary"), &summary)
	assert.False(t, diags.HasError(), "failed to read summary from state: %v", diags)
	return summary
}

// TestFwmOrganSummary_Read_Success tests a successful read that maps every
// API field into state.
func TestFwmOrganSummary_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	apiClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmV20250611Client", apiClient)

	var gotRequest *fwmv20250611.DescribeOrganSummaryRequest
	patches.ApplyMethodFunc(apiClient, "DescribeOrganSummaryWithContext", func(_ context.Context, request *fwmv20250611.DescribeOrganSummaryRequest) (*fwmv20250611.DescribeOrganSummaryResponse, error) {
		gotRequest = request
		resp := fwmv20250611.NewDescribeOrganSummaryResponse()
		resp.Response = &fwmv20250611.DescribeOrganSummaryResponseParams{
			RequestId: ptrStringFwmOrganSummary("fake-request-id"),
			Summary: &fwmv20250611.OrganSummary{
				GroupName: ptrStringFwmOrganSummary("tf-example-group"),
				AdminInfo: &fwmv20250611.MemberInfo{
					AppId:    ptrStringFwmOrganSummary("1300000001"),
					Uin:      ptrStringFwmOrganSummary("100000000001"),
					Nickname: ptrStringFwmOrganSummary("admin"),
					MemberId: ptrStringFwmOrganSummary("mem-xxxx"),
					NodeName: ptrStringFwmOrganSummary("root"),
				},
				AdminCount:         ptrInt64FwmOrganSummary(2),
				JoinedMemberCount:  ptrInt64FwmOrganSummary(3),
				MemberLimit:        ptrInt64FwmOrganSummary(-1),
				MemberLimitDisplay: ptrStringFwmOrganSummary("Unlimited"),
				CfwSharerCount:     ptrInt64FwmOrganSummary(1),
				CfwUserCount:       ptrInt64FwmOrganSummary(4),
				Departments: []*string{
					ptrStringFwmOrganSummary("dept-a"),
					ptrStringFwmOrganSummary("dept-b"),
				},
				ManagedMemberCount:  ptrInt64FwmOrganSummary(5),
				ManagedProductCount: ptrInt64FwmOrganSummary(6),
				CfwManageCount:      ptrInt64FwmOrganSummary(7),
				FwmShareUserQuota:   ptrInt64FwmOrganSummary(8),
			},
		}
		return resp, nil
	})

	ds := fwm.NewFwmOrganSummaryDataSource()
	setFwmOrganSummaryClient(ds, client)

	s := fwmOrganSummarySchema(t)
	readReq := newFwmOrganSummaryReadRequest(t, s)
	readResp := newFwmOrganSummaryReadResponse(t, s)
	ds.Read(context.Background(), readReq, readResp)

	assert.NotNil(t, gotRequest, "the API must be called once")

	summary := readFwmOrganSummaryState(t, &readResp)
	attrs := summary.Attributes()

	assert.Equal(t, "tf-example-group", attrs["group_name"].(types.String).ValueString())
	assert.Equal(t, int64(2), attrs["admin_count"].(types.Int64).ValueInt64())
	assert.Equal(t, int64(3), attrs["joined_member_count"].(types.Int64).ValueInt64())
	assert.Equal(t, int64(-1), attrs["member_limit"].(types.Int64).ValueInt64())
	assert.Equal(t, "Unlimited", attrs["member_limit_display"].(types.String).ValueString())
	assert.Equal(t, int64(1), attrs["cfw_sharer_count"].(types.Int64).ValueInt64())
	assert.Equal(t, int64(4), attrs["cfw_user_count"].(types.Int64).ValueInt64())
	assert.Equal(t, int64(5), attrs["managed_member_count"].(types.Int64).ValueInt64())
	assert.Equal(t, int64(6), attrs["managed_product_count"].(types.Int64).ValueInt64())
	assert.Equal(t, int64(7), attrs["cfw_manage_count"].(types.Int64).ValueInt64())
	assert.Equal(t, int64(8), attrs["fwm_share_user_quota"].(types.Int64).ValueInt64())

	departments, ok := attrs["departments"].(types.Set)
	assert.True(t, ok, "departments must be a set")
	assert.Equal(t, 2, len(departments.Elements()))

	adminInfo, ok := attrs["admin_info"].(types.Object)
	assert.True(t, ok, "admin_info must be an object")
	adminAttrs := adminInfo.Attributes()
	assert.Equal(t, "1300000001", adminAttrs["app_id"].(types.String).ValueString())
	assert.Equal(t, "100000000001", adminAttrs["uin"].(types.String).ValueString())
	assert.Equal(t, "admin", adminAttrs["nickname"].(types.String).ValueString())
	assert.Equal(t, "mem-xxxx", adminAttrs["member_id"].(types.String).ValueString())
	assert.Equal(t, "root", adminAttrs["node_name"].(types.String).ValueString())
}

// TestFwmOrganSummary_Read_EmptyDepartments tests that an empty department
// list is a valid empty set rather than an error.
func TestFwmOrganSummary_Read_EmptyDepartments(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	apiClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmV20250611Client", apiClient)

	patches.ApplyMethodFunc(apiClient, "DescribeOrganSummaryWithContext", func(_ context.Context, request *fwmv20250611.DescribeOrganSummaryRequest) (*fwmv20250611.DescribeOrganSummaryResponse, error) {
		resp := fwmv20250611.NewDescribeOrganSummaryResponse()
		resp.Response = &fwmv20250611.DescribeOrganSummaryResponseParams{
			RequestId: ptrStringFwmOrganSummary("fake-request-id"),
			Summary: &fwmv20250611.OrganSummary{
				GroupName:   ptrStringFwmOrganSummary("personal"),
				Departments: []*string{},
			},
		}
		return resp, nil
	})

	ds := fwm.NewFwmOrganSummaryDataSource()
	setFwmOrganSummaryClient(ds, client)

	s := fwmOrganSummarySchema(t)
	readReq := newFwmOrganSummaryReadRequest(t, s)
	readResp := newFwmOrganSummaryReadResponse(t, s)
	ds.Read(context.Background(), readReq, readResp)

	summary := readFwmOrganSummaryState(t, &readResp)
	attrs := summary.Attributes()

	departments := attrs["departments"].(types.Set)
	assert.False(t, departments.IsNull())
	assert.Equal(t, 0, len(departments.Elements()))

	// admin_info absent from the API response must be null, not an error.
	assert.True(t, attrs["admin_info"].(types.Object).IsNull())

	// Unset numeric fields must stay null so -1 keeps its sentinel meaning.
	assert.True(t, attrs["member_limit"].(types.Int64).IsNull())
	assert.Equal(t, "personal", attrs["group_name"].(types.String).ValueString())
}

// TestFwmOrganSummary_Read_APIError tests API error propagation.
func TestFwmOrganSummary_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	apiClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmV20250611Client", apiClient)

	patches.ApplyMethodFunc(apiClient, "DescribeOrganSummaryWithContext", func(_ context.Context, request *fwmv20250611.DescribeOrganSummaryRequest) (*fwmv20250611.DescribeOrganSummaryResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=UnauthorizedOperation, Message=not authorized")
	})

	ds := fwm.NewFwmOrganSummaryDataSource()
	setFwmOrganSummaryClient(ds, client)

	s := fwmOrganSummarySchema(t)
	readReq := newFwmOrganSummaryReadRequest(t, s)
	readResp := newFwmOrganSummaryReadResponse(t, s)
	ds.Read(context.Background(), readReq, readResp)

	assert.True(t, readResp.Diagnostics.HasError())
	assert.Contains(t, readResp.Diagnostics.Errors()[0].Summary(), "Error querying FWM organisation summary")
	assert.Contains(t, readResp.Diagnostics.Errors()[0].Detail(), "UnauthorizedOperation")
}

// TestFwmOrganSummary_Read_NilResponse tests that a nil API response is
// reported as an error instead of being treated as an empty result.
func TestFwmOrganSummary_Read_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	apiClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmV20250611Client", apiClient)

	patches.ApplyMethodFunc(apiClient, "DescribeOrganSummaryWithContext", func(_ context.Context, request *fwmv20250611.DescribeOrganSummaryRequest) (*fwmv20250611.DescribeOrganSummaryResponse, error) {
		return nil, nil
	})

	ds := fwm.NewFwmOrganSummaryDataSource()
	setFwmOrganSummaryClient(ds, client)

	s := fwmOrganSummarySchema(t)
	readReq := newFwmOrganSummaryReadRequest(t, s)
	readResp := newFwmOrganSummaryReadResponse(t, s)
	ds.Read(context.Background(), readReq, readResp)

	assert.True(t, readResp.Diagnostics.HasError())
	assert.Contains(t, readResp.Diagnostics.Errors()[0].Summary(), "Error querying FWM organisation summary")
}

// TestFwmOrganSummary_Read_NilSummary tests that a response without the
// summary object is reported as an error.
func TestFwmOrganSummary_Read_NilSummary(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	apiClient := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(client, "UseFwmV20250611Client", apiClient)

	patches.ApplyMethodFunc(apiClient, "DescribeOrganSummaryWithContext", func(_ context.Context, request *fwmv20250611.DescribeOrganSummaryRequest) (*fwmv20250611.DescribeOrganSummaryResponse, error) {
		resp := fwmv20250611.NewDescribeOrganSummaryResponse()
		resp.Response = &fwmv20250611.DescribeOrganSummaryResponseParams{
			RequestId: ptrStringFwmOrganSummary("fake-request-id"),
		}
		return resp, nil
	})

	ds := fwm.NewFwmOrganSummaryDataSource()
	setFwmOrganSummaryClient(ds, client)

	s := fwmOrganSummarySchema(t)
	readReq := newFwmOrganSummaryReadRequest(t, s)
	readResp := newFwmOrganSummaryReadResponse(t, s)
	ds.Read(context.Background(), readReq, readResp)

	assert.True(t, readResp.Diagnostics.HasError())
	assert.Contains(t, readResp.Diagnostics.Errors()[0].Detail(), "empty response")
}

// TestFwmOrganSummary_Read_ClientNotConfigured tests the nil-client guard.
func TestFwmOrganSummary_Read_ClientNotConfigured(t *testing.T) {
	ds := fwm.NewFwmOrganSummaryDataSource()

	s := fwmOrganSummarySchema(t)
	readReq := newFwmOrganSummaryReadRequest(t, s)
	readResp := newFwmOrganSummaryReadResponse(t, s)
	ds.Read(context.Background(), readReq, readResp)

	assert.True(t, readResp.Diagnostics.HasError())
	assert.Contains(t, readResp.Diagnostics.Errors()[0].Summary(), "Provider not configured")
}

// TestFwmOrganSummary_MetadataAndSchema validates the type name and that the
// data source exposes computed outputs only.
func TestFwmOrganSummary_MetadataAndSchema(t *testing.T) {
	ds := fwm.NewFwmOrganSummaryDataSource()

	metaResp := &datasource.MetadataResponse{}
	ds.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_fwm_organ_summary", metaResp.TypeName)

	schemaResp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())

	attrs := schemaResp.Schema.Attributes
	assert.Equal(t, 1, len(attrs))

	summaryAttr, ok := attrs["summary"].(schema.SingleNestedAttribute)
	assert.True(t, ok, "summary must be a single nested attribute")
	assert.True(t, summaryAttr.Computed)
	assert.False(t, summaryAttr.Optional)
	assert.False(t, summaryAttr.Required)

	for name, nested := range summaryAttr.Attributes {
		assert.True(t, nested.IsComputed(), "summary attribute %q must be computed", name)
		assert.False(t, nested.IsOptional(), "summary attribute %q must not be optional", name)
		assert.False(t, nested.IsRequired(), "summary attribute %q must not be required", name)
	}

	departmentsAttr, ok := summaryAttr.Attributes["departments"].(schema.SetAttribute)
	assert.True(t, ok, "departments must be a set attribute")
	assert.Equal(t, types.StringType, departmentsAttr.ElementType)

	adminInfoAttr, ok := summaryAttr.Attributes["admin_info"].(schema.SingleNestedAttribute)
	assert.True(t, ok, "admin_info must be a single nested attribute")
	for _, name := range []string{"app_id", "uin", "nickname", "member_id", "node_name"} {
		_, exists := adminInfoAttr.Attributes[name]
		assert.True(t, exists, "admin_info must expose %q", name)
	}

	objType, ok := summaryAttr.GetType().(attr.TypeWithAttributeTypes)
	assert.True(t, ok, "summary type must expose its attribute types")
	for name := range summaryAttr.Attributes {
		_, exists := objType.AttributeTypes()[name]
		assert.True(t, exists, "summary object type must contain %q", name)
	}
}
