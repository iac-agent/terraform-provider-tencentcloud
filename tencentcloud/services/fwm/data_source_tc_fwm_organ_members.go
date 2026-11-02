package fwm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// defaultOrganMemberPageSize is the page size used while walking the paging
// API. DescribeOrganMembers documents a default of 20 but no upper bound, so a
// value accepted by the sibling read paths of this service is used.
const defaultOrganMemberPageSize int64 = 100

// maxOrganMemberPages caps the paging loop so an unexpected TotalCount/Offset
// combination cannot turn into an endless request storm.
const maxOrganMemberPages int64 = 1000

var _ datasource.DataSourceWithConfigure = &FwmOrganMembersDataSource{}

// NewFwmOrganMembersDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewFwmOrganMembersDataSource() datasource.DataSource {
	return &FwmOrganMembersDataSource{}
}

// FwmOrganMembersDataSource implements datasource.DataSource for
// tencentcloud_fwm_organ_members.
type FwmOrganMembersDataSource struct {
	fw.DataSourceWithConfigure
}

// FwmOrganMembersDataSourceModel maps the schema attributes.
type FwmOrganMembersDataSourceModel struct {
	Members    types.List   `tfsdk:"members"`
	Filters    types.List   `tfsdk:"filters"`
	By         types.String `tfsdk:"by"`
	Order      types.String `tfsdk:"order"`
	TotalCount types.Int64  `tfsdk:"total_count"`
}

// fwmOrganMemberFilterModel is the element type of filters.
type fwmOrganMemberFilterModel struct {
	Name         types.String `tfsdk:"name"`
	Values       types.Set    `tfsdk:"values"`
	OperatorType types.Int64  `tfsdk:"operator_type"`
}

// fwmOrganMemberAccountGroupModel is the element type of members.account_group.
type fwmOrganMemberAccountGroupModel struct {
	GroupId   types.String `tfsdk:"group_id"`
	GroupName types.String `tfsdk:"group_name"`
}

// fwmOrganMemberModel is the element type of members.
type fwmOrganMemberModel struct {
	MemberId              types.String `tfsdk:"member_id"`
	AppId                 types.String `tfsdk:"app_id"`
	Uin                   types.String `tfsdk:"uin"`
	Nickname              types.String `tfsdk:"nickname"`
	SubAccountCount       types.Int64  `tfsdk:"sub_account_count"`
	NodeName              types.String `tfsdk:"node_name"`
	Role                  types.String `tfsdk:"role"`
	RoleDisplay           types.String `tfsdk:"role_display"`
	AccountGroup          types.List   `tfsdk:"account_group"`
	CfwManaged            types.Int64  `tfsdk:"cfw_managed"`
	CfwShareRole          types.String `tfsdk:"cfw_share_role"`
	CfwShareRoleDisplay   types.String `tfsdk:"cfw_share_role_display"`
	CfwSharerAppId        types.String `tfsdk:"cfw_sharer_app_id"`
	CfwInstanceId         types.String `tfsdk:"cfw_instance_id"`
	PolicyAnalysisEnabled types.Int64  `tfsdk:"policy_analysis_enabled"`
	MemberCreateTime      types.String `tfsdk:"member_create_time"`
	JoinType              types.String `tfsdk:"join_type"`
	CfwPayStatus          types.Int64  `tfsdk:"cfw_pay_status"`
	CfwCapable            types.Int64  `tfsdk:"cfw_capable"`
	SgManaged             types.Int64  `tfsdk:"sg_managed"`
	IsCfwPostPay          types.Int64  `tfsdk:"is_cfw_post_pay"`
}

// fwmOrganMemberAccountGroupAttrTypes is the attribute type map of one
// account_group element, shared by the schema conversion and the tests.
func fwmOrganMemberAccountGroupAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"group_id":   types.StringType,
		"group_name": types.StringType,
	}
}

func fwmOrganMemberFilterAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":          types.StringType,
		"values":        types.SetType{ElemType: types.StringType},
		"operator_type": types.Int64Type,
	}
}

func fwmOrganMemberAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"member_id":               types.StringType,
		"app_id":                  types.StringType,
		"uin":                     types.StringType,
		"nickname":                types.StringType,
		"sub_account_count":       types.Int64Type,
		"node_name":               types.StringType,
		"role":                    types.StringType,
		"role_display":            types.StringType,
		"account_group":           types.ListType{ElemType: types.ObjectType{AttrTypes: fwmOrganMemberAccountGroupAttrTypes()}},
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

func (d *FwmOrganMembersDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_fwm_organ_members"
}

func (d *FwmOrganMembersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the member list of a TencentCloud FWM (Cloud Firewall Manager) organization.",
		Attributes: map[string]schema.Attribute{
			// ---- Query inputs: Required / Optional only ----
			"filters": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Search conditions used to filter the returned members, such as member ID, account name, role or management status. All conditions are evaluated together.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Optional:    true,
							Description: "Name of the field to filter on, for example `MemberId`, `AppId` or `Nickname`.",
						},
						"values": schema.SetAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Description: "Values of the field to filter on.",
						},
						"operator_type": schema.Int64Attribute{
							Optional:    true,
							Description: "Operator of this condition: 1 - exact match, 7 - in, 9 - fuzzy match.",
						},
					},
				},
			},
			"by": schema.StringAttribute{
				Optional:    true,
				Description: "Field used to sort the returned members, for example `MemberCreateTime`.",
			},
			"order": schema.StringAttribute{
				Optional:    true,
				Description: "Sort direction: `asc` for ascending, `desc` for descending.",
			},
			// ---- Query outputs: Computed ----
			"total_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of members that match the query conditions.",
			},
			"members": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of FWM organization members.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"member_id": schema.StringAttribute{
							Computed:    true,
							Description: "Member ID.",
						},
						"app_id": schema.StringAttribute{
							Computed:    true,
							Description: "AppId of the member account.",
						},
						"uin": schema.StringAttribute{
							Computed:    true,
							Description: "Uin of the member account.",
						},
						"nickname": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the member account.",
						},
						"sub_account_count": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of sub-accounts of the member.",
						},
						"node_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the organization node the member belongs to.",
						},
						"role": schema.StringAttribute{
							Computed:    true,
							Description: "Role of the member: `admin` - administrator, `delegatedAdmin` - delegated administrator, `member` - ordinary member.",
						},
						"role_display": schema.StringAttribute{
							Computed:    true,
							Description: "Display name of the member role.",
						},
						"account_group": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Account group the member belongs to.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"group_id": schema.StringAttribute{
										Computed:    true,
										Description: "Account group ID.",
									},
									"group_name": schema.StringAttribute{
										Computed:    true,
										Description: "Account group name.",
									},
								},
							},
						},
						"cfw_managed": schema.Int64Attribute{
							Computed:    true,
							Description: "Cloud firewall management status: 0 - not managed, 1 - managed.",
						},
						"cfw_share_role": schema.StringAttribute{
							Computed:    true,
							Description: "Cloud firewall sharing role: `sharer` - sharer, `user` - user, `none` - not set.",
						},
						"cfw_share_role_display": schema.StringAttribute{
							Computed:    true,
							Description: "Display name of the cloud firewall sharing role.",
						},
						"cfw_sharer_app_id": schema.StringAttribute{
							Computed:    true,
							Description: "AppId of the cloud firewall sharer, set when the member acts as a user.",
						},
						"cfw_instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "Cloud firewall billing instance ID. A non-empty value means cloud firewall has been purchased.",
						},
						"policy_analysis_enabled": schema.Int64Attribute{
							Computed:    true,
							Description: "Policy analysis permission: 0 - disabled, 1 - enabled.",
						},
						"member_create_time": schema.StringAttribute{
							Computed:    true,
							Description: "Time when the member joined the organization.",
						},
						"join_type": schema.StringAttribute{
							Computed:    true,
							Description: "Way the account joined the organization.",
						},
						"cfw_pay_status": schema.Int64Attribute{
							Computed:    true,
							Description: "Cloud firewall billing status: 0 - not purchased, 2 - purchased, 3 - trial, 4 - expired.",
						},
						"cfw_capable": schema.Int64Attribute{
							Computed:    true,
							Description: "Whether the member is capable of using cloud firewall.",
						},
						"sg_managed": schema.Int64Attribute{
							Computed:    true,
							Description: "Private security group management switch: 0 - not managed, 1 - managed.",
						},
						"is_cfw_post_pay": schema.Int64Attribute{
							Computed:    true,
							Description: "Whether the member uses the pay-as-you-go cloud firewall edition: 1 - pay-as-you-go, 0 - not pay-as-you-go, -1 - unknown.",
						},
					},
				},
			},
		},
	}
}

// Read queries every page of DescribeOrganMembers and writes the aggregated
// result into resp.State.
func (d *FwmOrganMembersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config FwmOrganMembersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := FwmService{client: d.Client()}

	filters := make([]*fwmv20250611.CommonFilter, 0)
	if !config.Filters.IsNull() && !config.Filters.IsUnknown() {
		var filterModels []fwmOrganMemberFilterModel
		resp.Diagnostics.Append(config.Filters.ElementsAs(ctx, &filterModels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, item := range filterModels {
			filter := &fwmv20250611.CommonFilter{}
			if !item.Name.IsNull() && !item.Name.IsUnknown() {
				filter.Name = item.Name.ValueStringPointer()
			}
			if !item.Values.IsNull() && !item.Values.IsUnknown() {
				var values []types.String
				resp.Diagnostics.Append(item.Values.ElementsAs(ctx, &values, false)...)
				if resp.Diagnostics.HasError() {
					return
				}

				stringValues := make([]*string, 0, len(values))
				for _, v := range values {
					if v.IsNull() || v.IsUnknown() {
						continue
					}
					stringValues = append(stringValues, v.ValueStringPointer())
				}
				filter.Values = stringValues
			}
			if !item.OperatorType.IsNull() && !item.OperatorType.IsUnknown() {
				filter.OperatorType = item.OperatorType.ValueInt64Pointer()
			}
			filters = append(filters, filter)
		}
	}

	var by, order string
	if !config.By.IsNull() && !config.By.IsUnknown() {
		by = config.By.ValueString()
	}
	if !config.Order.IsNull() && !config.Order.IsUnknown() {
		order = config.Order.ValueString()
	}

	items := make([]fwmOrganMemberModel, 0)
	var totalCount int64
	for page := int64(0); page < maxOrganMemberPages; page++ {
		var (
			offset  = page * defaultOrganMemberPageSize
			members []*fwmv20250611.OrganMemberItem
			err     error
		)
		members, totalCount, err = service.DescribeFwmOrganMembersByFilter(ctx, filters, by, order, defaultOrganMemberPageSize, offset)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error querying FWM organization members",
				fmt.Sprintf("DescribeOrganMembers failed: %v", err),
			)
			return
		}

		for _, member := range members {
			if member == nil {
				continue
			}
			item, diags := flattenFwmOrganMember(ctx, member)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			items = append(items, item)
		}

		if len(members) < int(defaultOrganMemberPageSize) {
			break
		}
		if totalCount > 0 && int64(len(items)) >= totalCount {
			break
		}
	}

	membersList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: fwmOrganMemberAttrTypes()}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// An empty result is a valid business state (no member matched), so it is
	// written back as an empty list instead of being reported as an error.
	config.Members = membersList
	config.TotalCount = types.Int64Value(totalCount)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// flattenFwmOrganMember converts one SDK member item into its state model.
func flattenFwmOrganMember(ctx context.Context, member *fwmv20250611.OrganMemberItem) (fwmOrganMemberModel, diag.Diagnostics) {
	accountGroups := make([]fwmOrganMemberAccountGroupModel, 0, 1)
	if member.AccountGroup != nil {
		accountGroups = append(accountGroups, fwmOrganMemberAccountGroupModel{
			GroupId:   types.StringPointerValue(member.AccountGroup.GroupId),
			GroupName: types.StringPointerValue(member.AccountGroup.GroupName),
		})
	}

	// A member without account group yields an empty list rather than null, so
	// that the same value round-trips between plan and state.
	accountGroupList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: fwmOrganMemberAccountGroupAttrTypes()}, accountGroups)

	return fwmOrganMemberModel{
		MemberId:              types.StringPointerValue(member.MemberId),
		AppId:                 types.StringPointerValue(member.AppId),
		Uin:                   types.StringPointerValue(member.Uin),
		Nickname:              types.StringPointerValue(member.Nickname),
		SubAccountCount:       types.Int64PointerValue(member.SubAccountCount),
		NodeName:              types.StringPointerValue(member.NodeName),
		Role:                  types.StringPointerValue(member.Role),
		RoleDisplay:           types.StringPointerValue(member.RoleDisplay),
		AccountGroup:          accountGroupList,
		CfwManaged:            types.Int64PointerValue(member.CfwManaged),
		CfwShareRole:          types.StringPointerValue(member.CfwShareRole),
		CfwShareRoleDisplay:   types.StringPointerValue(member.CfwShareRoleDisplay),
		CfwSharerAppId:        types.StringPointerValue(member.CfwSharerAppId),
		CfwInstanceId:         types.StringPointerValue(member.CfwInstanceId),
		PolicyAnalysisEnabled: types.Int64PointerValue(member.PolicyAnalysisEnabled),
		MemberCreateTime:      types.StringPointerValue(member.MemberCreateTime),
		JoinType:              types.StringPointerValue(member.JoinType),
		CfwPayStatus:          types.Int64PointerValue(member.CfwPayStatus),
		CfwCapable:            types.Int64PointerValue(member.CfwCapable),
		SgManaged:             types.Int64PointerValue(member.SgManaged),
		IsCfwPostPay:          types.Int64PointerValue(member.IsCfwPostPay),
	}, diags
}
