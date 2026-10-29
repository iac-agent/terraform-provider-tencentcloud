package fwm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// _ resource.Resource is intentionally omitted: the factory function
// NewFwmRuleGroupResource returns resource.Resource, which already performs
// the implicit compile-time check at the return statement. Only the
// extended-interface assertions are kept here.
var (
	_ resource.ResourceWithConfigure   = &FwmRuleGroupResource{}
	_ resource.ResourceWithImportState = &FwmRuleGroupResource{}
)

// NewFwmRuleGroupResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewFwmRuleGroupResource() resource.Resource {
	return &FwmRuleGroupResource{}
}

// FwmRuleGroupResource implements resource.Resource for
// tencentcloud_fwm_rule_group.
type FwmRuleGroupResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// FwmRuleGroupMemberModel maps one member of a rule group receive account.
type FwmRuleGroupMemberModel struct {
	AppId    types.String `tfsdk:"app_id"`
	Uin      types.String `tfsdk:"uin"`
	Nickname types.String `tfsdk:"nickname"`
	MemberId types.String `tfsdk:"member_id"`
	NodeName types.String `tfsdk:"node_name"`
}

// FwmRuleGroupReceiveAccountModel maps one receive account entry.
type FwmRuleGroupReceiveAccountModel struct {
	Uin          types.String              `tfsdk:"uin"`
	Nickname     types.String              `tfsdk:"nickname"`
	ReceiverType types.Int64               `tfsdk:"receiver_type"`
	Members      []FwmRuleGroupMemberModel `tfsdk:"members"`
}

// FwmRuleGroupFilterModel maps one filter of the describe request.
type FwmRuleGroupFilterModel struct {
	Name         types.String   `tfsdk:"name"`
	Values       []types.String `tfsdk:"values"`
	OperatorType types.Int64    `tfsdk:"operator_type"`
}

// FwmRuleGroupRuleGroupModel maps one element of the rule_groups list.
type FwmRuleGroupRuleGroupModel struct {
	AppId          types.String                      `tfsdk:"app_id"`
	GroupId        types.String                      `tfsdk:"group_id"`
	GroupName      types.String                      `tfsdk:"group_name"`
	RuleCount      types.Int64                       `tfsdk:"rule_count"`
	Product        types.String                      `tfsdk:"product"`
	CreateTime     types.String                      `tfsdk:"create_time"`
	UpdateTime     types.String                      `tfsdk:"update_time"`
	CreateBy       types.String                      `tfsdk:"create_by"`
	CreateName     types.String                      `tfsdk:"create_name"`
	UpdateBy       types.String                      `tfsdk:"update_by"`
	UpdateName     types.String                      `tfsdk:"update_name"`
	ReceiveAccount []FwmRuleGroupReceiveAccountModel `tfsdk:"receive_account"`
	ExecAreas      []types.String                    `tfsdk:"exec_areas"`
	Direction      types.Int64                       `tfsdk:"direction"`
}

// FwmRuleGroupResourceModel maps the schema attributes.
type FwmRuleGroupResourceModel struct {
	Id         types.String                 `tfsdk:"id"`
	GroupName  types.String                 `tfsdk:"group_name"`
	Product    types.String                 `tfsdk:"product"`
	Filters    []FwmRuleGroupFilterModel    `tfsdk:"filters"`
	Order      types.String                 `tfsdk:"order"`
	By         types.String                 `tfsdk:"by"`
	RuleGroups []FwmRuleGroupRuleGroupModel `tfsdk:"rule_groups"`
	GroupId    types.String                 `tfsdk:"group_id"`
	GroupIds   []types.String               `tfsdk:"group_ids"`
}

func (r *FwmRuleGroupResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_fwm_rule_group"
}

func (r *FwmRuleGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a FWM (Firewall Manager) rule group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the FWM rule group resource.",
				PlanModifiers: []planmodifier.String{
					stringUseStateForUnknown{},
				},
			},
			"group_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the rule group. It cannot be empty and can be up to 50 characters long (Chinese characters count as one character).",
			},
			"product": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Product the rule group belongs to. Valid values: `enterprise_sg` (enterprise security group), " +
					"`cfw_vpc_acl` (VPC border firewall), `cfw_nat_acl` (NAT border firewall), `cfw_edge_acl` (internet border firewall). " +
					"It is required and does not accept `all`.",
			},
			"filters": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Filters used to query the rule group list.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Name of the filter field.",
						},
						"values": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Computed:    true,
							Description: "Values of the filter field.",
						},
						"operator_type": schema.Int64Attribute{
							Optional:    true,
							Computed:    true,
							Description: "Operator type. Valid values: `1` (exact match), `7` (in), `9` (fuzzy match).",
						},
					},
				},
			},
			"order": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Sorting order of the query result, e.g. `desc` or `asc`.",
			},
			"by": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Field used to sort the query result.",
			},
			"rule_groups": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of rule groups.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"app_id": schema.StringAttribute{
							Computed:    true,
							Description: "Tenant AppId.",
						},
						"group_id": schema.StringAttribute{
							Computed:    true,
							Description: "Rule group ID.",
						},
						"group_name": schema.StringAttribute{
							Computed:    true,
							Description: "Rule group name.",
						},
						"rule_count": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of rules in the rule group.",
						},
						"product": schema.StringAttribute{
							Computed:    true,
							Description: "Product type.",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "Creation time.",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "Update time.",
						},
						"create_by": schema.StringAttribute{
							Computed:    true,
							Description: "SubUin of the creator.",
						},
						"create_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the creator.",
						},
						"update_by": schema.StringAttribute{
							Computed:    true,
							Description: "SubUin of the last updater.",
						},
						"update_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the last updater.",
						},
						"receive_account": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Associated accounts the rule group is delivered to.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"uin": schema.StringAttribute{
										Computed:    true,
										Description: "Tenant Uin.",
									},
									"nickname": schema.StringAttribute{
										Computed:    true,
										Description: "Tenant name.",
									},
									"receiver_type": schema.Int64Attribute{
										Computed:    true,
										Description: "Receiver type. `0` means account Uin, `1` means account group.",
									},
									"members": schema.ListNestedAttribute{
										Computed:    true,
										Description: "Member list, returned only when `receiver_type` is `1`.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"app_id": schema.StringAttribute{
													Computed:    true,
													Description: "Member AppId.",
												},
												"uin": schema.StringAttribute{
													Computed:    true,
													Description: "Member Uin.",
												},
												"nickname": schema.StringAttribute{
													Computed:    true,
													Description: "Member nickname.",
												},
												"member_id": schema.StringAttribute{
													Computed:    true,
													Description: "Member ID, empty in personal scenarios.",
												},
												"node_name": schema.StringAttribute{
													Computed:    true,
													Description: "Department the member belongs to.",
												},
											},
										},
									},
								},
							},
						},
						"exec_areas": schema.SetAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "Execution stages the rule group belongs to.",
						},
						"direction": schema.Int64Attribute{
							Computed:    true,
							Description: "Rule direction. `1` means inbound, `0` means outbound, `-1` means no direction.",
						},
					},
				},
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "ID of the rule group used by the modify operation. It must be the ID of the rule group managed by this resource.",
				PlanModifiers: []planmodifier.String{
					stringUseStateForUnknown{},
				},
			},
			"group_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "IDs of the rule groups used by the delete operation. It must contain the ID of the rule group managed by this resource.",
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *FwmRuleGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FwmRuleGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.GroupName.IsNull() || plan.GroupName.IsUnknown() || plan.GroupName.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("group_name"),
			"Invalid group_name",
			"group_name is required by CreateRuleGroup and cannot be empty.",
		)
		return
	}
	if plan.Product.IsNull() || plan.Product.IsUnknown() || plan.Product.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("product"),
			"Invalid product",
			"product is required by CreateRuleGroup and cannot be empty.",
		)
		return
	}

	service := NewFwmService(r.Client())
	groupId, err := service.CreateFwmRuleGroup(ctx, plan.GroupName.ValueString(), plan.Product.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating FWM rule group",
			fmt.Sprintf("CreateRuleGroup failed for rule group %q: %v", plan.GroupName.ValueString(), err),
		)
		return
	}

	plan.Id = types.StringValue(groupId)
	plan.GroupId = types.StringValue(groupId)
	plan.GroupIds = []types.String{types.StringValue(groupId)}

	ruleGroup, err := service.DescribeFwmRuleGroupById(ctx, groupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading FWM rule group after create",
			fmt.Sprintf("DescribeRuleGroups failed for rule group %q: %v", groupId, err),
		)
		return
	}
	if ruleGroup != nil {
		r.refreshModel(&plan, ruleGroup)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *FwmRuleGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FwmRuleGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId := state.Id.ValueString()

	service := NewFwmService(r.Client())
	ruleGroup, err := service.DescribeFwmRuleGroupById(ctx, groupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading FWM rule group",
			fmt.Sprintf("DescribeRuleGroups failed for rule group %q: %v", groupId, err),
		)
		return
	}

	if ruleGroup == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("FWM rule group %s not found, removing from state.", groupId),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	r.refreshModel(&state, ruleGroup)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates.
func (r *FwmRuleGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FwmRuleGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state FwmRuleGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId := state.Id.ValueString()

	if !plan.GroupName.IsNull() && !plan.GroupName.IsUnknown() && !plan.GroupName.Equal(state.GroupName) {
		service := NewFwmService(r.Client())
		if err := service.ModifyFwmRuleGroupName(ctx, groupId, plan.GroupName.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				"Error updating FWM rule group",
				fmt.Sprintf("ModifyRuleGroup failed for rule group %q: %v", groupId, err),
			)
			return
		}
	}

	plan.Id = state.Id
	plan.GroupId = types.StringValue(groupId)

	service := NewFwmService(r.Client())
	ruleGroup, err := service.DescribeFwmRuleGroupById(ctx, groupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading FWM rule group after update",
			fmt.Sprintf("DescribeRuleGroups failed for rule group %q: %v", groupId, err),
		)
		return
	}
	if ruleGroup != nil {
		r.refreshModel(&plan, ruleGroup)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles `terraform destroy`.
func (r *FwmRuleGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FwmRuleGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewFwmService(r.Client())
	if err := service.DeleteFwmRuleGroupById(ctx, state.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting FWM rule group",
			fmt.Sprintf("DeleteRuleGroup failed for rule group %q: %v", state.Id.ValueString(), err),
		)
		return
	}
	// After a successful delete the Framework removes the state automatically.
}

// ImportState supports `terraform import <address> <group_id>`.
func (r *FwmRuleGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// refreshModel copies the API item fields into the model. Only cloud-owned
// attributes are refreshed; the query-only inputs (filters/order/by) and the
// write-only operation inputs (group_id/group_ids) keep their plan values.
func (r *FwmRuleGroupResource) refreshModel(m *FwmRuleGroupResourceModel, item *fwmv20250611.RuleGroupResp) {
	if item == nil {
		return
	}

	ruleGroup := buildRuleGroupItem(item)
	m.GroupName = ruleGroup.GroupName
	m.Product = ruleGroup.Product
	m.RuleGroups = []FwmRuleGroupRuleGroupModel{ruleGroup}
}

// buildRuleGroupItem converts one cloud rule group into the rule_groups element model.
func buildRuleGroupItem(item *fwmv20250611.RuleGroupResp) FwmRuleGroupRuleGroupModel {
	return FwmRuleGroupRuleGroupModel{
		AppId:          types.StringPointerValue(item.AppId),
		GroupId:        types.StringPointerValue(item.GroupId),
		GroupName:      types.StringPointerValue(item.GroupName),
		RuleCount:      types.Int64PointerValue(item.RuleCount),
		Product:        types.StringPointerValue(item.Product),
		CreateTime:     types.StringPointerValue(item.CreateTime),
		UpdateTime:     types.StringPointerValue(item.UpdateTime),
		CreateBy:       types.StringPointerValue(item.CreateBy),
		CreateName:     types.StringPointerValue(item.CreateName),
		UpdateBy:       types.StringPointerValue(item.UpdateBy),
		UpdateName:     types.StringPointerValue(item.UpdateName),
		ReceiveAccount: buildReceiveAccounts(item.ReceiveAccount),
		ExecAreas:      stringSliceToTypeSlice(item.ExecAreas),
		Direction:      types.Int64PointerValue(item.Direction),
	}
}

// BuildReceiveAccountsForTest exposes buildReceiveAccounts to unit tests.
func BuildReceiveAccountsForTest(in []*fwmv20250611.ReceiveAccount) []FwmRuleGroupReceiveAccountModel {
	return buildReceiveAccounts(in)
}

// StringSliceToTypeSliceForTest exposes stringSliceToTypeSlice to unit tests.
func StringSliceToTypeSliceForTest(in []*string) []types.String {
	return stringSliceToTypeSlice(in)
}

// BuildRuleGroupModelForTest converts a cloud rule group into the rule_groups
// element model, so unit tests can assert the mapping without a live client.
func BuildRuleGroupModelForTest(item *fwmv20250611.RuleGroupResp) FwmRuleGroupRuleGroupModel {
	return buildRuleGroupItem(item)
}

// stringSliceToTypeSlice converts a string pointer slice into a framework
// types.String slice. An empty cloud value produces an empty list so the
// optional list attribute never fails to be persisted.
func stringSliceToTypeSlice(in []*string) []types.String {
	out := make([]types.String, 0, len(in))
	for _, v := range in {
		if v == nil {
			out = append(out, types.StringNull())
			continue
		}
		out = append(out, types.StringValue(*v))
	}
	return out
}

// buildReceiveAccounts converts the SDK receive account entries into the
// framework nested list model.
func buildReceiveAccounts(in []*fwmv20250611.ReceiveAccount) []FwmRuleGroupReceiveAccountModel {
	if len(in) == 0 {
		return nil
	}

	out := make([]FwmRuleGroupReceiveAccountModel, 0, len(in))
	for _, v := range in {
		if v == nil {
			continue
		}

		account := FwmRuleGroupReceiveAccountModel{
			Uin:          types.StringPointerValue(v.Uin),
			Nickname:     types.StringPointerValue(v.Nickname),
			ReceiverType: types.Int64PointerValue(v.ReceiverType),
		}

		if len(v.Members) > 0 {
			members := make([]FwmRuleGroupMemberModel, 0, len(v.Members))
			for _, member := range v.Members {
				if member == nil {
					continue
				}
				members = append(members, FwmRuleGroupMemberModel{
					AppId:    types.StringPointerValue(member.AppId),
					Uin:      types.StringPointerValue(member.Uin),
					Nickname: types.StringPointerValue(member.Nickname),
					MemberId: types.StringPointerValue(member.MemberId),
					NodeName: types.StringPointerValue(member.NodeName),
				})
			}
			account.Members = members
		}

		out = append(out, account)
	}
	return out
}

// stringUseStateForUnknown copies the prior state value into the plan when
// the new plan value is unknown, used for computed attributes so they stay
// stable across refreshes.
type stringUseStateForUnknown struct{}

func (stringUseStateForUnknown) Description(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}

func (stringUseStateForUnknown) MarkdownDescription(ctx context.Context) string {
	return stringUseStateForUnknown{}.Description(ctx)
}

func (stringUseStateForUnknown) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() {
		return
	}
	if !req.PlanValue.IsUnknown() {
		return
	}
	resp.PlanValue = req.StateValue
}
