package fwm_test

import (
	"context"
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/fwm"
)

func newOrganMembersDataSource(t *testing.T) datasource.DataSourceWithConfigure {
	t.Helper()

	ds, ok := fwm.NewFwmOrganMembersDataSource().(datasource.DataSourceWithConfigure)
	if !ok {
		t.Fatal("expect tencentcloud_fwm_organ_members to implement datasource.DataSourceWithConfigure")
	}
	return ds
}

// buildOrganMembersConfig builds a Config value for the given optional inputs.
func buildOrganMembersConfig(t *testing.T, filters []fwmOrganMemberFilterInput, by, order *string) tfsdk.Config {
	t.Helper()

	ds := newOrganMembersDataSource(t)
	schemaResp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)

	filterAttrTypes := organMemberFilterAttrTypes()

	filtersValue := types.ListNull(types.ObjectType{AttrTypes: filterAttrTypes})
	if filters != nil {
		elements := make([]attr.Value, 0, len(filters))
		for _, f := range filters {
			values := types.SetNull(types.StringType)
			if f.values != nil {
				elems := make([]attr.Value, 0, len(f.values))
				for _, v := range f.values {
					elems = append(elems, types.StringValue(v))
				}
				set, diags := types.SetValue(types.StringType, elems)
				if diags.HasError() {
					t.Fatalf("build filters values: %v", diags)
				}
				values = set
			}

			o, diags := types.ObjectValue(filterAttrTypes, map[string]attr.Value{
				"name":          types.StringPointerValue(f.name),
				"values":        values,
				"operator_type": types.Int64PointerValue(f.operatorType),
			})
			if diags.HasError() {
				t.Fatalf("build filters element: %v", diags)
			}
			elements = append(elements, o)
		}

		list, diags := types.ListValue(types.ObjectType{AttrTypes: filterAttrTypes}, elements)
		if diags.HasError() {
			t.Fatalf("build filters: %v", diags)
		}
		filtersValue = list
	}

	value, diags := types.ObjectValue(organMembersDataSourceAttrTypes(), map[string]attr.Value{
		"filters":     filtersValue,
		"by":          types.StringPointerValue(by),
		"order":       types.StringPointerValue(order),
		"members":     types.ListNull(types.ObjectType{AttrTypes: organMemberAttrTypes()}),
		"total_count": types.Int64Null(),
	})
	if diags.HasError() {
		t.Fatalf("build config: %v", diags)
	}

	raw, err := value.ToTerraformValue(context.Background())
	if err != nil {
		t.Fatalf("convert config to terraform value: %v", err)
	}

	return tfsdk.Config{
		Schema: schemaResp.Schema,
		Raw:    raw,
	}
}

// organMembersDataSourceAttrTypes returns the top-level attribute types of the
// data source, used to build a complete Config object in tests.
func organMembersDataSourceAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"filters":     types.ListType{ElemType: types.ObjectType{AttrTypes: organMemberFilterAttrTypes()}},
		"by":          types.StringType,
		"order":       types.StringType,
		"members":     types.ListType{ElemType: types.ObjectType{AttrTypes: organMemberAttrTypes()}},
		"total_count": types.Int64Type,
	}
}

func organMemberFilterAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":          types.StringType,
		"values":        types.SetType{ElemType: types.StringType},
		"operator_type": types.Int64Type,
	}
}

type fwmOrganMemberFilterInput struct {
	name         *string
	values       []string
	operatorType *int64
}

func organMemberAccountGroupAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"group_id":   types.StringType,
		"group_name": types.StringType,
	}
}

func organMemberAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"member_id":               types.StringType,
		"app_id":                  types.StringType,
		"uin":                     types.StringType,
		"nickname":                types.StringType,
		"sub_account_count":       types.Int64Type,
		"node_name":               types.StringType,
		"role":                    types.StringType,
		"role_display":            types.StringType,
		"account_group":           types.ListType{ElemType: types.ObjectType{AttrTypes: organMemberAccountGroupAttrTypes()}},
		"cfw_managed":             types.Int64Type,
		"cfw_share_role":          types.StringType,
		"cfw_share_role_display":  types.StringType,
		"cfw_sharer_app_id":       types.StringType,
		"cfw_instance_id":         types.StringType,
		"policy_analysis_enabled": types.Int64Type,
		"member_create_time":      types.StringType,
		"join_type":               types.StringType,
		"cfw_pay_status":          types.Int64Type,
		"cfw_capable":             types.Int64Type,
		"sg_managed":              types.Int64Type,
		"is_cfw_post_pay":         types.Int64Type,
	}
}

func buildFullOrganMemberItem() *fwmv20250611.OrganMemberItem {
	return &fwmv20250611.OrganMemberItem{
		MemberId:              ptrStringFwm("mem-xxxx"),
		AppId:                 ptrStringFwm("1300000001"),
		Uin:                   ptrStringFwm("100000000001"),
		Nickname:              ptrStringFwm("tf-example"),
		SubAccountCount:       ptrInt64OrganMembers(3),
		NodeName:              ptrStringFwm("root"),
		Role:                  ptrStringFwm("admin"),
		RoleDisplay:           ptrStringFwm("Administrator"),
		AccountGroup:          &fwmv20250611.AccountGroupInfo{GroupId: ptrStringFwm("acg-xxxx"), GroupName: ptrStringFwm("group-name")},
		CfwManaged:            ptrInt64OrganMembers(1),
		CfwShareRole:          ptrStringFwm("sharer"),
		CfwShareRoleDisplay:   ptrStringFwm("Sharer"),
		CfwSharerAppId:        ptrStringFwm("1300000002"),
		CfwInstanceId:         ptrStringFwm("cfw-xxxx"),
		PolicyAnalysisEnabled: ptrInt64OrganMembers(1),
		MemberCreateTime:      ptrStringFwm("2025-06-11 10:00:00"),
		JoinType:              ptrStringFwm("invite"),
		CfwPayStatus:          ptrInt64OrganMembers(2),
		CfwCapable:            ptrInt64OrganMembers(1),
		SgManaged:             ptrInt64OrganMembers(0),
		IsCfwPostPay:          ptrInt64OrganMembers(-1),
	}
}

func ptrInt64OrganMembers(v int64) *int64 {
	return &v
}

func mockOrganMembersAPI(patches *gomonkey.Patches, client *fwmv20250611.Client, members []*fwmv20250611.OrganMemberItem, total int64) *fwmv20250611.DescribeOrganMembersRequest {
	var captured *fwmv20250611.DescribeOrganMembersRequest
	patches.ApplyMethodFunc(client, "DescribeOrganMembersWithContext", func(_ context.Context, request *fwmv20250611.DescribeOrganMembersRequest) (*fwmv20250611.DescribeOrganMembersResponse, error) {
		captured = request
		resp := fwmv20250611.NewDescribeOrganMembersResponse()
		resp.Response = &fwmv20250611.DescribeOrganMembersResponseParams{
			RequestId:  ptrStringFwm("fake-request-id"),
			TotalCount: ptrInt64OrganMembers(total),
			Members:    members,
		}
		return resp, nil
	})
	return captured
}

func readOrganMembers(t *testing.T, ds datasource.DataSourceWithConfigure, config tfsdk.Config) (*datasource.ReadResponse, *fwmOrganMembersState) {
	t.Helper()

	schemaResp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)

	resp := &datasource.ReadResponse{
		State: tfsdk.State{Schema: schemaResp.Schema, Raw: tftypes.NewValue(schemaResp.Schema.Type().TerraformType(context.Background()), nil)},
	}
	ds.Read(context.Background(), datasource.ReadRequest{Config: config}, resp)
	if resp.Diagnostics.HasError() {
		return resp, nil
	}

	var state fwmOrganMembersState
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &state)...)
	return resp, &state
}

type fwmOrganMembersState struct {
	Members    []fwmOrganMemberState       `tfsdk:"members"`
	Filters    []fwmOrganMemberFilterState `tfsdk:"filters"`
	By         types.String                `tfsdk:"by"`
	Order      types.String                `tfsdk:"order"`
	TotalCount types.Int64                 `tfsdk:"total_count"`
}

type fwmOrganMemberFilterState struct {
	Name         types.String `tfsdk:"name"`
	Values       types.Set    `tfsdk:"values"`
	OperatorType types.Int64  `tfsdk:"operator_type"`
}

type fwmOrganMemberAccountGroupState struct {
	GroupId   types.String `tfsdk:"group_id"`
	GroupName types.String `tfsdk:"group_name"`
}

type fwmOrganMemberState struct {
	MemberId              types.String                      `tfsdk:"member_id"`
	AppId                 types.String                      `tfsdk:"app_id"`
	Uin                   types.String                      `tfsdk:"uin"`
	Nickname              types.String                      `tfsdk:"nickname"`
	SubAccountCount       types.Int64                       `tfsdk:"sub_account_count"`
	NodeName              types.String                      `tfsdk:"node_name"`
	Role                  types.String                      `tfsdk:"role"`
	RoleDisplay           types.String                      `tfsdk:"role_display"`
	AccountGroup          []fwmOrganMemberAccountGroupState `tfsdk:"account_group"`
	CfwManaged            types.Int64                       `tfsdk:"cfw_managed"`
	CfwShareRole          types.String                      `tfsdk:"cfw_share_role"`
	CfwShareRoleDisplay   types.String                      `tfsdk:"cfw_share_role_display"`
	CfwSharerAppId        types.String                      `tfsdk:"cfw_sharer_app_id"`
	CfwInstanceId         types.String                      `tfsdk:"cfw_instance_id"`
	PolicyAnalysisEnabled types.Int64                       `tfsdk:"policy_analysis_enabled"`
	MemberCreateTime      types.String                      `tfsdk:"member_create_time"`
	JoinType              types.String                      `tfsdk:"join_type"`
	CfwPayStatus          types.Int64                       `tfsdk:"cfw_pay_status"`
	CfwCapable            types.Int64                       `tfsdk:"cfw_capable"`
	SgManaged             types.Int64                       `tfsdk:"sg_managed"`
	IsCfwPostPay          types.Int64                       `tfsdk:"is_cfw_post_pay"`
}

func setupOrganMembersMock(t *testing.T) (*gomonkey.Patches, *fwmv20250611.Client, datasource.DataSourceWithConfigure) {
	t.Helper()

	patches := gomonkey.NewPatches()
	client := &fwmv20250611.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseFwmV20250611Client", client)

	ds := newOrganMembersDataSource(t)
	ds.Configure(context.Background(), datasource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: &connectivity.TencentCloudClient{}},
	}, &datasource.ConfigureResponse{})

	return patches, client, ds
}

func TestFwmOrganMembersDataSource_Read_AllFields(t *testing.T) {
	patches, client, ds := setupOrganMembersMock(t)
	defer patches.Reset()

	mockOrganMembersAPI(patches, client, []*fwmv20250611.OrganMemberItem{buildFullOrganMemberItem()}, 1)

	config := buildOrganMembersConfig(t, nil, nil, nil)
	resp, state := readOrganMembers(t, ds, config)
	if state == nil {
		t.Fatalf("read failed: %v", resp.Diagnostics)
	}

	if !state.TotalCount.Equal(types.Int64Value(1)) {
		t.Fatalf("expect total_count 1, got %s", state.TotalCount.String())
	}
	if len(state.Members) != 1 {
		t.Fatalf("expect 1 member, got %d", len(state.Members))
	}

	m := state.Members[0]
	assertOrganMemberString(t, "member_id", m.MemberId, "mem-xxxx")
	assertOrganMemberString(t, "app_id", m.AppId, "1300000001")
	assertOrganMemberString(t, "uin", m.Uin, "100000000001")
	assertOrganMemberString(t, "nickname", m.Nickname, "tf-example")
	assertOrganMemberInt(t, "sub_account_count", m.SubAccountCount, 3)
	assertOrganMemberString(t, "node_name", m.NodeName, "root")
	assertOrganMemberString(t, "role", m.Role, "admin")
	assertOrganMemberString(t, "role_display", m.RoleDisplay, "Administrator")
	assertOrganMemberInt(t, "cfw_managed", m.CfwManaged, 1)
	assertOrganMemberString(t, "cfw_share_role", m.CfwShareRole, "sharer")
	assertOrganMemberString(t, "cfw_share_role_display", m.CfwShareRoleDisplay, "Sharer")
	assertOrganMemberString(t, "cfw_sharer_app_id", m.CfwSharerAppId, "1300000002")
	assertOrganMemberString(t, "cfw_instance_id", m.CfwInstanceId, "cfw-xxxx")
	assertOrganMemberInt(t, "policy_analysis_enabled", m.PolicyAnalysisEnabled, 1)
	assertOrganMemberString(t, "member_create_time", m.MemberCreateTime, "2025-06-11 10:00:00")
	assertOrganMemberString(t, "join_type", m.JoinType, "invite")
	assertOrganMemberInt(t, "cfw_pay_status", m.CfwPayStatus, 2)
	assertOrganMemberInt(t, "cfw_capable", m.CfwCapable, 1)
	assertOrganMemberInt(t, "sg_managed", m.SgManaged, 0)
	assertOrganMemberInt(t, "is_cfw_post_pay", m.IsCfwPostPay, -1)

	if len(m.AccountGroup) != 1 {
		t.Fatalf("expect 1 account group, got %d", len(m.AccountGroup))
	}
	assertOrganMemberString(t, "account_group.group_id", m.AccountGroup[0].GroupId, "acg-xxxx")
	assertOrganMemberString(t, "account_group.group_name", m.AccountGroup[0].GroupName, "group-name")
}

func TestFwmOrganMembersDataSource_Read_EmptyAndNullAccountGroup(t *testing.T) {
	patches, client, ds := setupOrganMembersMock(t)
	defer patches.Reset()

	member := &fwmv20250611.OrganMemberItem{
		MemberId: ptrStringFwm("mem-yyyy"),
		AppId:    ptrStringFwm("1300000003"),
	}
	mockOrganMembersAPI(patches, client, []*fwmv20250611.OrganMemberItem{member}, 1)

	config := buildOrganMembersConfig(t, nil, nil, nil)
	resp, state := readOrganMembers(t, ds, config)
	if state == nil {
		t.Fatalf("read failed: %v", resp.Diagnostics)
	}

	if len(state.Members) != 1 {
		t.Fatalf("expect 1 member, got %d", len(state.Members))
	}
	m := state.Members[0]
	if len(m.AccountGroup) != 0 {
		t.Fatalf("expect empty account_group, got %d elements", len(m.AccountGroup))
	}
	// Nil SDK pointers must become null values, not zero values.
	if !m.SubAccountCount.IsNull() {
		t.Fatalf("expect sub_account_count null, got %s", m.SubAccountCount.String())
	}
	if !m.Uin.IsNull() {
		t.Fatalf("expect uin null, got %s", m.Uin.String())
	}
	if !m.CfwManaged.IsNull() {
		t.Fatalf("expect cfw_managed null, got %s", m.CfwManaged.String())
	}
}

func TestFwmOrganMembersDataSource_Read_NoMatchReturnsEmptyList(t *testing.T) {
	patches, client, ds := setupOrganMembersMock(t)
	defer patches.Reset()

	mockOrganMembersAPI(patches, client, nil, 0)

	config := buildOrganMembersConfig(t, nil, nil, nil)
	resp, state := readOrganMembers(t, ds, config)
	if state == nil {
		t.Fatalf("empty result must not fail the read: %v", resp.Diagnostics)
	}

	if state.Members == nil {
		t.Fatal("expect members to be a list value, got nil slice")
	}
	if len(state.Members) != 0 {
		t.Fatalf("expect empty members, got %d", len(state.Members))
	}
	if !state.TotalCount.Equal(types.Int64Value(0)) {
		t.Fatalf("expect total_count 0, got %s", state.TotalCount.String())
	}
}

func TestFwmOrganMembersDataSource_Read_ForwardsQueryInputs(t *testing.T) {
	patches, client, ds := setupOrganMembersMock(t)
	defer patches.Reset()

	var captured *fwmv20250611.DescribeOrganMembersRequest
	patches.ApplyMethodFunc(client, "DescribeOrganMembersWithContext", func(_ context.Context, request *fwmv20250611.DescribeOrganMembersRequest) (*fwmv20250611.DescribeOrganMembersResponse, error) {
		captured = request
		resp := fwmv20250611.NewDescribeOrganMembersResponse()
		resp.Response = &fwmv20250611.DescribeOrganMembersResponseParams{
			RequestId:  ptrStringFwm("fake-request-id"),
			TotalCount: ptrInt64OrganMembers(1),
			Members:    []*fwmv20250611.OrganMemberItem{buildFullOrganMemberItem()},
		}
		return resp, nil
	})

	name := "MemberId"
	operatorType := int64(1)
	by := "MemberCreateTime"
	order := "desc"
	config := buildOrganMembersConfig(t, []fwmOrganMemberFilterInput{
		{name: &name, values: []string{"mem-xxxx"}, operatorType: &operatorType},
	}, &by, &order)

	resp, state := readOrganMembers(t, ds, config)
	if state == nil {
		t.Fatalf("read failed: %v", resp.Diagnostics)
	}

	if captured == nil {
		t.Fatal("expect DescribeOrganMembers to be called")
	}
	if len(captured.Filters) != 1 {
		t.Fatalf("expect 1 filter, got %d", len(captured.Filters))
	}
	if captured.Filters[0].Name == nil || *captured.Filters[0].Name != "MemberId" {
		t.Fatalf("expect filter name MemberId, got %v", captured.Filters[0].Name)
	}
	if len(captured.Filters[0].Values) != 1 || captured.Filters[0].Values[0] == nil || *captured.Filters[0].Values[0] != "mem-xxxx" {
		t.Fatalf("expect filter values [mem-xxxx], got %v", captured.Filters[0].Values)
	}
	if captured.Filters[0].OperatorType == nil || *captured.Filters[0].OperatorType != 1 {
		t.Fatalf("expect operator type 1, got %v", captured.Filters[0].OperatorType)
	}
	if captured.By == nil || *captured.By != "MemberCreateTime" {
		t.Fatalf("expect by MemberCreateTime, got %v", captured.By)
	}
	if captured.Order == nil || *captured.Order != "desc" {
		t.Fatalf("expect order desc, got %v", captured.Order)
	}
	if captured.Limit == nil || *captured.Limit != 100 {
		t.Fatalf("expect limit 100, got %v", captured.Limit)
	}
	if captured.Offset == nil || *captured.Offset != 0 {
		t.Fatalf("expect offset 0, got %v", captured.Offset)
	}
}

func TestFwmOrganMembersDataSource_Read_Paginates(t *testing.T) {
	patches, client, ds := setupOrganMembersMock(t)
	defer patches.Reset()

	var offsets []int64
	patches.ApplyMethodFunc(client, "DescribeOrganMembersWithContext", func(_ context.Context, request *fwmv20250611.DescribeOrganMembersRequest) (*fwmv20250611.DescribeOrganMembersResponse, error) {
		var offset int64
		if request.Offset != nil {
			offset = *request.Offset
		}
		offsets = append(offsets, offset)

		var members []*fwmv20250611.OrganMemberItem
		if offset == 0 {
			members = make([]*fwmv20250611.OrganMemberItem, 0, 100)
			for i := 0; i < 100; i++ {
				members = append(members, &fwmv20250611.OrganMemberItem{
					MemberId: ptrStringFwm("mem-first-page"),
				})
			}
		} else {
			members = []*fwmv20250611.OrganMemberItem{{MemberId: ptrStringFwm("mem-second-page")}}
		}

		resp := fwmv20250611.NewDescribeOrganMembersResponse()
		resp.Response = &fwmv20250611.DescribeOrganMembersResponseParams{
			RequestId:  ptrStringFwm("fake-request-id"),
			TotalCount: ptrInt64OrganMembers(101),
			Members:    members,
		}
		return resp, nil
	})

	config := buildOrganMembersConfig(t, nil, nil, nil)
	resp, state := readOrganMembers(t, ds, config)
	if state == nil {
		t.Fatalf("read failed: %v", resp.Diagnostics)
	}

	if len(offsets) != 2 {
		t.Fatalf("expect 2 paged requests, got %d (%v)", len(offsets), offsets)
	}
	if offsets[0] != 0 || offsets[1] != 100 {
		t.Fatalf("expect offsets [0 100], got %v", offsets)
	}
	if len(state.Members) != 101 {
		t.Fatalf("expect 101 members, got %d", len(state.Members))
	}
	if !state.TotalCount.Equal(types.Int64Value(101)) {
		t.Fatalf("expect total_count 101, got %s", state.TotalCount.String())
	}
}

func TestFwmOrganMembersDataSource_Read_ApiError(t *testing.T) {
	patches, client, ds := setupOrganMembersMock(t)
	defer patches.Reset()

	patches.ApplyMethodFunc(client, "DescribeOrganMembersWithContext", func(_ context.Context, _ *fwmv20250611.DescribeOrganMembersRequest) (*fwmv20250611.DescribeOrganMembersResponse, error) {
		return nil, errors.New("fake api error")
	})

	config := buildOrganMembersConfig(t, nil, nil, nil)
	resp, state := readOrganMembers(t, ds, config)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expect an error diagnostic when the API call fails")
	}
	if state != nil {
		t.Fatal("expect no state to be written when the API call fails")
	}
}

func TestFwmOrganMembersDataSource_Read_NilResponse(t *testing.T) {
	patches, client, ds := setupOrganMembersMock(t)
	defer patches.Reset()

	patches.ApplyMethodFunc(client, "DescribeOrganMembersWithContext", func(_ context.Context, _ *fwmv20250611.DescribeOrganMembersRequest) (*fwmv20250611.DescribeOrganMembersResponse, error) {
		return nil, nil
	})

	config := buildOrganMembersConfig(t, nil, nil, nil)
	resp, _ := readOrganMembers(t, ds, config)
	if !resp.Diagnostics.HasError() {
		t.Fatal("a nil response must be reported as an error instead of an empty result")
	}
}

func TestFwmOrganMembersDataSource_Schema(t *testing.T) {
	ds := newOrganMembersDataSource(t)

	resp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema diagnostics: %v", resp.Diagnostics)
	}

	for _, name := range []string{"filters", "by", "order"} {
		attrSchema, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Fatalf("schema attribute %q is missing", name)
		}
		if !attrSchema.IsOptional() || attrSchema.IsRequired() || attrSchema.IsComputed() {
			t.Fatalf("schema attribute %q must be optional-only", name)
		}
	}

	for _, name := range []string{"members", "total_count"} {
		attrSchema, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Fatalf("schema attribute %q is missing", name)
		}
		if !attrSchema.IsComputed() {
			t.Fatalf("schema attribute %q must be computed", name)
		}
	}
}

// assertNoErrorDiagnostics fails the test when the read reported any error
// diagnostic.
func assertNoErrorDiagnostics(t *testing.T, diags diag.Diagnostics) {
	t.Helper()

	if diags.HasError() {
		t.Fatalf("expect no error diagnostics, got %v", diags)
	}
}

func TestFwmOrganMembersDataSource_Read_NoDiagnostics(t *testing.T) {
	patches, client, ds := setupOrganMembersMock(t)
	defer patches.Reset()

	mockOrganMembersAPI(patches, client, []*fwmv20250611.OrganMemberItem{buildFullOrganMemberItem()}, 1)

	resp, state := readOrganMembers(t, ds, buildOrganMembersConfig(t, nil, nil, nil))
	assertNoErrorDiagnostics(t, resp.Diagnostics)
	if state == nil {
		t.Fatal("expect state to be written")
	}
}

func assertOrganMemberString(t *testing.T, name string, got types.String, want string) {
	t.Helper()

	if got.IsNull() || got.IsUnknown() {
		t.Fatalf("%s: expect %q, got a null/unknown value", name, want)
	}
	if got.ValueString() != want {
		t.Fatalf("%s: expect %q, got %q", name, want, got.ValueString())
	}
}

func assertOrganMemberInt(t *testing.T, name string, got types.Int64, want int64) {
	t.Helper()

	if got.IsNull() || got.IsUnknown() {
		t.Fatalf("%s: expect %d, got a null/unknown value", name, want)
	}
	if got.ValueInt64() != want {
		t.Fatalf("%s: expect %d, got %d", name, want, got.ValueInt64())
	}
}
