package fwm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var _ datasource.DataSourceWithConfigure = &FwmOrganSummaryDataSource{}

// NewFwmOrganSummaryDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewFwmOrganSummaryDataSource() datasource.DataSource {
	return &FwmOrganSummaryDataSource{}
}

// FwmOrganSummaryDataSource implements datasource.DataSource for
// tencentcloud_fwm_organ_summary.
//
// DescribeOrganSummary accepts no business parameters: it always returns the
// organisation overview of the account authenticated by the provider. This
// data source therefore exposes no query inputs, only computed outputs.
type FwmOrganSummaryDataSource struct {
	fw.DataSourceWithConfigure
}

// FwmOrganSummaryDataSourceModel maps the data source schema attributes.
type FwmOrganSummaryDataSourceModel struct {
	Summary types.Object `tfsdk:"summary"`
}

// fwmOrganSummaryAdminInfoAttrTypes describes the admin_info nested object.
var fwmOrganSummaryAdminInfoAttrTypes = map[string]attr.Type{
	"app_id":    types.StringType,
	"uin":       types.StringType,
	"nickname":  types.StringType,
	"member_id": types.StringType,
	"node_name": types.StringType,
}

// fwmOrganSummaryAttrTypes describes the summary nested object.
var fwmOrganSummaryAttrTypes = map[string]attr.Type{
	"group_name":            types.StringType,
	"admin_info":            types.ObjectType{AttrTypes: fwmOrganSummaryAdminInfoAttrTypes},
	"admin_count":           types.Int64Type,
	"joined_member_count":   types.Int64Type,
	"member_limit":          types.Int64Type,
	"member_limit_display":  types.StringType,
	"cfw_sharer_count":      types.Int64Type,
	"cfw_user_count":        types.Int64Type,
	"departments":           types.SetType{ElemType: types.StringType},
	"managed_member_count":  types.Int64Type,
	"managed_product_count": types.Int64Type,
	"cfw_manage_count":      types.Int64Type,
	"fwm_share_user_quota":  types.Int64Type,
}

// fwmOrganSummaryScalarAttrPaths maps the scalar summary attributes to their
// schema paths.
var fwmOrganSummaryScalarAttrPaths = map[string]path.Path{
	"group_name":            path.Root("summary").AtName("group_name"),
	"admin_count":           path.Root("summary").AtName("admin_count"),
	"joined_member_count":   path.Root("summary").AtName("joined_member_count"),
	"member_limit":          path.Root("summary").AtName("member_limit"),
	"member_limit_display":  path.Root("summary").AtName("member_limit_display"),
	"cfw_sharer_count":      path.Root("summary").AtName("cfw_sharer_count"),
	"cfw_user_count":        path.Root("summary").AtName("cfw_user_count"),
	"managed_member_count":  path.Root("summary").AtName("managed_member_count"),
	"managed_product_count": path.Root("summary").AtName("managed_product_count"),
	"cfw_manage_count":      path.Root("summary").AtName("cfw_manage_count"),
	"fwm_share_user_quota":  path.Root("summary").AtName("fwm_share_user_quota"),
}

func (d *FwmOrganSummaryDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_fwm_organ_summary"
}

func (d *FwmOrganSummaryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the organisation overview of FWM (Firewall Manager). " +
			"The result is always scoped to the account authenticated by the provider.",
		Attributes: map[string]schema.Attribute{
			"summary": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Organisation overview returned by the FWM API.",
				Attributes: map[string]schema.Attribute{
					"group_name": schema.StringAttribute{
						Computed:    true,
						Description: "Name of the organisation group.",
					},
					"admin_info": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Account information of the administrator.",
						Attributes: map[string]schema.Attribute{
							"app_id": schema.StringAttribute{
								Computed:    true,
								Description: "AppId of the administrator account.",
							},
							"uin": schema.StringAttribute{
								Computed:    true,
								Description: "Uin of the administrator account.",
							},
							"nickname": schema.StringAttribute{
								Computed:    true,
								Description: "Nickname of the administrator account.",
							},
							"member_id": schema.StringAttribute{
								Computed:    true,
								Description: "Member ID of the administrator account. Empty in personal (non-organisation) scenarios.",
							},
							"node_name": schema.StringAttribute{
								Computed:    true,
								Description: "Department the administrator account belongs to.",
							},
						},
					},
					"admin_count": schema.Int64Attribute{
						Computed:    true,
						Description: "Number of administrators and delegated administrators.",
					},
					"joined_member_count": schema.Int64Attribute{
						Computed:    true,
						Description: "Number of members that have joined the organisation.",
					},
					"member_limit": schema.Int64Attribute{
						Computed:    true,
						Description: "Upper limit of members that can join the organisation. The value -1 means no limit.",
					},
					"member_limit_display": schema.StringAttribute{
						Computed:    true,
						Description: "Display text of the member limit.",
					},
					"cfw_sharer_count": schema.Int64Attribute{
						Computed:    true,
						Description: "Number of specification sharers.",
					},
					"cfw_user_count": schema.Int64Attribute{
						Computed:    true,
						Description: "Number of specification users.",
					},
					"departments": schema.SetAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "Names of the departments in the organisation.",
					},
					"managed_member_count": schema.Int64Attribute{
						Computed:    true,
						Description: "Number of managed accounts.",
					},
					"managed_product_count": schema.Int64Attribute{
						Computed:    true,
						Description: "Number of managed products.",
					},
					"cfw_manage_count": schema.Int64Attribute{
						Computed:    true,
						Description: "Number of accounts managed through specification sharing.",
					},
					"fwm_share_user_quota": schema.Int64Attribute{
						Computed:    true,
						Description: "Total quota allocated to specification sharing users.",
					},
				},
			},
		},
	}
}

// Read queries the organisation overview and writes the result into
// resp.State. DescribeOrganSummary takes no business parameters, so no query
// input is forwarded.
func (d *FwmOrganSummaryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config FwmOrganSummaryDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if d.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The FWM data source was invoked before the provider client was configured.",
		)
		return
	}

	summary, err := NewFwmService(d.Client()).DescribeFwmOrganSummary(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying FWM organisation summary",
			fmt.Sprintf("DescribeOrganSummary failed: %v", err),
		)
		return
	}

	// A nil summary is an abnormal API response rather than an empty result:
	// the organisation overview always exists for an authenticated account.
	if summary == nil {
		resp.Diagnostics.AddError(
			"Error querying FWM organisation summary",
			"DescribeOrganSummary returned an empty response.",
		)
		return
	}

	setFwmOrganSummaryState(ctx, summary, &resp.State, &resp.Diagnostics)
}

// setFwmOrganSummaryState flattens the API result into state attribute by
// attribute. Account identities are written straight into state and are never
// held in the configuration model struct.
func setFwmOrganSummaryState(ctx context.Context, summary *fwmv20250611.OrganSummary, state *tfsdk.State, diags *diag.Diagnostics) {
	writeScalar := func(name string, value attr.Value) {
		p, ok := fwmOrganSummaryScalarAttrPaths[name]
		if !ok {
			diags.AddError(
				"Internal error",
				fmt.Sprintf("unknown summary attribute %q", name),
			)
			return
		}
		diags.Append(state.SetAttribute(ctx, p, value)...)
	}

	writeScalar("group_name", types.StringPointerValue(summary.GroupName))
	writeScalar("admin_count", types.Int64PointerValue(summary.AdminCount))
	writeScalar("joined_member_count", types.Int64PointerValue(summary.JoinedMemberCount))
	writeScalar("member_limit", types.Int64PointerValue(summary.MemberLimit))
	writeScalar("member_limit_display", types.StringPointerValue(summary.MemberLimitDisplay))
	writeScalar("cfw_sharer_count", types.Int64PointerValue(summary.CfwSharerCount))
	writeScalar("cfw_user_count", types.Int64PointerValue(summary.CfwUserCount))
	writeScalar("managed_member_count", types.Int64PointerValue(summary.ManagedMemberCount))
	writeScalar("managed_product_count", types.Int64PointerValue(summary.ManagedProductCount))
	writeScalar("cfw_manage_count", types.Int64PointerValue(summary.CfwManageCount))
	writeScalar("fwm_share_user_quota", types.Int64PointerValue(summary.FwmShareUserQuota))

	// departments is an unordered string collection; an empty result is a
	// valid value, so it is written as an empty set instead of null.
	departments := make([]types.String, 0, len(summary.Departments))
	for _, department := range summary.Departments {
		if department == nil {
			continue
		}
		departments = append(departments, types.StringValue(*department))
	}
	departmentSet, ds := types.SetValueFrom(ctx, types.StringType, departments)
	diags.Append(ds...)
	if diags.HasError() {
		return
	}
	writeScalar("departments", departmentSet)

	// admin_info is a fixed single object; nil attributes become null values.
	adminInfo := types.ObjectNull(fwmOrganSummaryAdminInfoAttrTypes)
	if summary.AdminInfo != nil {
		v, ds := types.ObjectValue(
			fwmOrganSummaryAdminInfoAttrTypes,
			map[string]attr.Value{
				"app_id":    types.StringPointerValue(summary.AdminInfo.AppId),
				"uin":       types.StringPointerValue(summary.AdminInfo.Uin),
				"nickname":  types.StringPointerValue(summary.AdminInfo.Nickname),
				"member_id": types.StringPointerValue(summary.AdminInfo.MemberId),
				"node_name": types.StringPointerValue(summary.AdminInfo.NodeName),
			},
		)
		diags.Append(ds...)
		if diags.HasError() {
			return
		}
		adminInfo = v
	}

	summaryValue, summaryDiags := types.ObjectValue(
		fwmOrganSummaryAttrTypes,
		map[string]attr.Value{
			"group_name":            types.StringPointerValue(summary.GroupName),
			"admin_info":            adminInfo,
			"admin_count":           types.Int64PointerValue(summary.AdminCount),
			"joined_member_count":   types.Int64PointerValue(summary.JoinedMemberCount),
			"member_limit":          types.Int64PointerValue(summary.MemberLimit),
			"member_limit_display":  types.StringPointerValue(summary.MemberLimitDisplay),
			"cfw_sharer_count":      types.Int64PointerValue(summary.CfwSharerCount),
			"cfw_user_count":        types.Int64PointerValue(summary.CfwUserCount),
			"departments":           departmentSet,
			"managed_member_count":  types.Int64PointerValue(summary.ManagedMemberCount),
			"managed_product_count": types.Int64PointerValue(summary.ManagedProductCount),
			"cfw_manage_count":      types.Int64PointerValue(summary.CfwManageCount),
			"fwm_share_user_quota":  types.Int64PointerValue(summary.FwmShareUserQuota),
		},
	)
	diags.Append(summaryDiags...)
	if diags.HasError() {
		return
	}

	diags.Append(state.Set(ctx, &FwmOrganSummaryDataSourceModel{Summary: summaryValue})...)
}
