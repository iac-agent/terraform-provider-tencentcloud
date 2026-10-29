package fwm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// _ resource.Resource is intentionally omitted: the factory function
// NewFwmStrategyV2Resource returns resource.Resource, which already performs
// the implicit compile-time check at the return statement. Only the
// extended-interface assertions are kept here.
var (
	_ resource.ResourceWithConfigure   = &FwmStrategyV2Resource{}
	_ resource.ResourceWithImportState = &FwmStrategyV2Resource{}
)

// NewFwmStrategyV2Resource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewFwmStrategyV2Resource() resource.Resource {
	return &FwmStrategyV2Resource{}
}

// FwmStrategyV2Resource implements resource.Resource for
// tencentcloud_fwm_strategy_v2.
type FwmStrategyV2Resource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// FwmStrategyV2ResourceModel maps the schema attributes.
type FwmStrategyV2ResourceModel struct {
	Id             types.String `tfsdk:"id"`
	Product        types.String `tfsdk:"product"`
	ReceiveAccount types.Set    `tfsdk:"receive_account"`
	PreStrategy    types.List   `tfsdk:"pre_strategy"`
	PostStrategy   types.List   `tfsdk:"post_strategy"`
	ReceiveGroup   types.Set    `tfsdk:"receive_group"`
	Strategy       types.List   `tfsdk:"strategy"`
}

// strategyGroupModel maps one entry of pre_strategy or post_strategy.
type strategyGroupModel struct {
	GroupId  types.String `tfsdk:"group_id"`
	Sequence types.Int64  `tfsdk:"sequence"`
}

// strategyReceiveAccountModel maps one entry of the describe result
// response.Response.Strategy[].ReceiveAccount.
type strategyReceiveAccountModel struct {
	Uin          types.String          `tfsdk:"uin"`
	Nickname     types.String          `tfsdk:"nickname"`
	ReceiverType types.Int64           `tfsdk:"receiver_type"`
	Members      []strategyMemberModel `tfsdk:"members"`
}

// strategyMemberModel maps one entry of
// response.Response.Strategy[].ReceiveAccount[].Members.
type strategyMemberModel struct {
	AppId    types.String `tfsdk:"app_id"`
	Uin      types.String `tfsdk:"uin"`
	Nickname types.String `tfsdk:"nickname"`
	MemberId types.String `tfsdk:"member_id"`
	NodeName types.String `tfsdk:"node_name"`
}

// strategyDetailModel maps one entry of the describe result
// response.Response.Strategy.
type strategyDetailModel struct {
	StrategyId     types.String                  `tfsdk:"strategy_id"`
	GroupId        types.String                  `tfsdk:"group_id"`
	GroupName      types.String                  `tfsdk:"group_name"`
	RuleCount      types.Int64                   `tfsdk:"rule_count"`
	RuleStatus     types.Int64                   `tfsdk:"rule_status"`
	ReceiveAccount []strategyReceiveAccountModel `tfsdk:"receive_account"`
	Sequence       types.Int64                   `tfsdk:"sequence"`
	ErrMsg         types.String                  `tfsdk:"err_msg"`
	ErrorType      types.String                  `tfsdk:"error_type"`
	UpdateTime     types.String                  `tfsdk:"update_time"`
	CreateBy       types.String                  `tfsdk:"create_by"`
	UpdateBy       types.String                  `tfsdk:"update_by"`
	ExecArea       types.String                  `tfsdk:"exec_area"`
	CreateName     types.String                  `tfsdk:"create_name"`
	UpdateName     types.String                  `tfsdk:"update_name"`
	CreateTime     types.String                  `tfsdk:"create_time"`
}

func (r *FwmStrategyV2Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_fwm_strategy_v2"
}

func (r *FwmStrategyV2Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a FWM (Firewall Manager) strategy. " +
			"A strategy defines the rule groups executed before and after the firewall and " +
			"the accounts or account groups the rules are delivered to.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the FWM strategy.",
				PlanModifiers: []planmodifier.String{
					stringUseStateForUnknown{},
				},
			},
			"product": schema.StringAttribute{
				Required:    true,
				Description: "Type of the product the strategy belongs to.",
			},
			"receive_account": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Uins of the accounts that receive the delivered rules.",
			},
			"pre_strategy": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Rule groups executed before the firewall.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"group_id": schema.StringAttribute{
							Optional:    true,
							Description: "Rule group ID.",
						},
						"sequence": schema.Int64Attribute{
							Optional:    true,
							Description: "Priority of the rule group.",
						},
					},
				},
			},
			"post_strategy": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Rule groups executed after the firewall.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"group_id": schema.StringAttribute{
							Optional:    true,
							Description: "Rule group ID.",
						},
						"sequence": schema.Int64Attribute{
							Optional:    true,
							Description: "Priority of the rule group.",
						},
					},
				},
			},
			"receive_group": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "IDs of the account groups that receive the delivered rules.",
			},
			"strategy": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Details of the strategies generated by the creation request.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"strategy_id": schema.StringAttribute{
							Computed:    true,
							Description: "Strategy ID.",
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
						"rule_status": schema.Int64Attribute{
							Computed:    true,
							Description: "Status of the strategy.",
						},
						"receive_account": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Accounts that receive the delivered rules.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"uin": schema.StringAttribute{
										Computed:    true,
										Description: "Uin of the tenant.",
									},
									"nickname": schema.StringAttribute{
										Computed:    true,
										Description: "Name of the tenant.",
									},
									"receiver_type": schema.Int64Attribute{
										Computed:    true,
										Description: "Receiver type. `0`: account uin, `1`: account group.",
									},
									"members": schema.ListNestedAttribute{
										Computed:    true,
										Description: "Members of the account group. Only returned when receiver_type is `1`.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"app_id": schema.StringAttribute{
													Computed:    true,
													Description: "AppId of the member.",
												},
												"uin": schema.StringAttribute{
													Computed:    true,
													Description: "Uin of the member.",
												},
												"nickname": schema.StringAttribute{
													Computed:    true,
													Description: "Nickname of the member.",
												},
												"member_id": schema.StringAttribute{
													Computed:    true,
													Description: "Member ID. Empty for personal accounts.",
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
						"sequence": schema.Int64Attribute{
							Computed:    true,
							Description: "Priority of the strategy.",
						},
						"err_msg": schema.StringAttribute{
							Computed:    true,
							Description: "Reason why the delivery failed.",
						},
						"error_type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of the delivery failure reason.",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "Update time of the strategy.",
						},
						"create_by": schema.StringAttribute{
							Computed:    true,
							Description: "Creator of the strategy.",
						},
						"update_by": schema.StringAttribute{
							Computed:    true,
							Description: "Updater of the strategy.",
						},
						"exec_area": schema.StringAttribute{
							Computed:    true,
							Description: "Execution area of the strategy.",
						},
						"create_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the creator.",
						},
						"update_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the updater.",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "Create time of the strategy.",
						},
					},
				},
			},
		},
	}
}

// Create handles `terraform apply` for a new FWM strategy.
func (r *FwmStrategyV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FwmStrategyV2ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	request := fwmv20250611.NewCreateStrategyRequest()
	request.Product = plan.Product.ValueStringPointer()

	if !plan.ReceiveAccount.IsNull() && !plan.ReceiveAccount.IsUnknown() {
		var accounts []types.String
		resp.Diagnostics.Append(plan.ReceiveAccount.ElementsAs(ctx, &accounts, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, account := range accounts {
			request.ReceiveAccount = append(request.ReceiveAccount, account.ValueStringPointer())
		}
	}

	if !plan.ReceiveGroup.IsNull() && !plan.ReceiveGroup.IsUnknown() {
		var groups []types.String
		resp.Diagnostics.Append(plan.ReceiveGroup.ElementsAs(ctx, &groups, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, group := range groups {
			request.ReceiveGroup = append(request.ReceiveGroup, group.ValueStringPointer())
		}
	}

	buildStrategyGroups := func(list types.List) ([]*fwmv20250611.StrategyReq, bool) {
		if list.IsNull() || list.IsUnknown() {
			return nil, true
		}
		var items []strategyGroupModel
		diags := list.ElementsAs(ctx, &items, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return nil, false
		}
		groups := make([]*fwmv20250611.StrategyReq, 0, len(items))
		for _, item := range items {
			strategyReq := &fwmv20250611.StrategyReq{}
			if !item.GroupId.IsNull() && !item.GroupId.IsUnknown() {
				strategyReq.GroupId = item.GroupId.ValueStringPointer()
			}
			if !item.Sequence.IsNull() && !item.Sequence.IsUnknown() {
				strategyReq.Sequence = item.Sequence.ValueInt64Pointer()
			}
			groups = append(groups, strategyReq)
		}
		return groups, true
	}

	preStrategies, ok := buildStrategyGroups(plan.PreStrategy)
	if !ok {
		return
	}
	request.PreStrategy = preStrategies

	postStrategies, ok := buildStrategyGroups(plan.PostStrategy)
	if !ok {
		return
	}
	request.PostStrategy = postStrategies

	logId := tccommon.GetLogId(ctx)
	log.Printf("[DEBUG]%s fwm_strategy_v2 create request, product=%s", logId, plan.Product.ValueString())

	response, err := r.Client().UseFwmV20250611Client().CreateStrategyWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating fwm_strategy_v2",
			fmt.Sprintf("CreateStrategy failed, product=%s: %v", plan.Product.ValueString(), err),
		)
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		resp.Diagnostics.AddError(
			"Error creating fwm_strategy_v2",
			fmt.Sprintf("CreateStrategy returned an empty response, product=%s", plan.Product.ValueString()),
		)
		return
	}

	strategyIds := response.Response.StrategyIds
	if len(strategyIds) == 0 {
		resp.Diagnostics.AddError(
			"Error creating fwm_strategy_v2",
			fmt.Sprintf("CreateStrategy returned an empty StrategyIds, product=%s", plan.Product.ValueString()),
		)
		return
	}

	id := ""
	for _, strategyId := range strategyIds {
		if strategyId == nil || *strategyId == "" {
			continue
		}
		if id != "" {
			id += tccommon.FILED_SP
		}
		id += *strategyId
	}
	if id == "" {
		resp.Diagnostics.AddError(
			"Error creating fwm_strategy_v2",
			fmt.Sprintf("CreateStrategy returned empty StrategyIds values, product=%s", plan.Product.ValueString()),
		)
		return
	}

	// The strategy IDs are written to state before reading back so that a
	// partially created strategy can still be destroyed by Terraform.
	plan.Id = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	strategies, err := r.read(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading fwm_strategy_v2 after create",
			fmt.Sprintf("id=%s: %v", id, err),
		)
		return
	}
	plan.Strategy = strategies
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *FwmStrategyV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FwmStrategyV2ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	id := state.Id.ValueString()
	strategies, err := r.read(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading fwm_strategy_v2",
			fmt.Sprintf("id=%s: %v", id, err),
		)
		return
	}

	if strategies.IsNull() {
		log.Printf("[WARN] fwm_strategy_v2 id=%s not found on the cloud side, removing it from state.", id)
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("fwm_strategy_v2 %s was not found, removing it from state.", id),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	state.Strategy = strategies
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates of a FWM strategy.
func (r *FwmStrategyV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state FwmStrategyV2ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	ids := splitFwmStrategyIds(state.Id.ValueString())
	if len(ids) != 1 {
		resp.Diagnostics.AddError(
			"Unsupported fwm_strategy_v2 update",
			fmt.Sprintf("ModifyStrategy can only update a strategy composed of exactly one strategy ID, current id=%s", state.Id.ValueString()),
		)
		return
	}

	// receive_account and receive_group are shared by every strategy created by
	// the same create request, so changing them on a multi strategy resource
	// would leave the remaining strategies inconsistent.
	if len(splitFwmStrategyIds(plan.Id.ValueString())) > 1 {
		if !plan.ReceiveAccount.Equal(state.ReceiveAccount) || !plan.ReceiveGroup.Equal(state.ReceiveGroup) {
			resp.Diagnostics.AddError(
				"Unsupported fwm_strategy_v2 update",
				fmt.Sprintf("receive_account and receive_group cannot be modified on the resource id=%s because it contains more than one strategy; destroy and recreate the resource instead.", state.Id.ValueString()),
			)
			return
		}
	}

	// Only strategy_id, receive_account, sequence, group_id and receive_group
	// are accepted by ModifyStrategy, so every provided optional argument is
	// sent and the rule group scope (pre_strategy / post_strategy) comes from
	// the rule group type of group_id.
	request := fwmv20250611.NewModifyStrategyRequest()
	request.StrategyId = helper.String(ids[0])

	if !plan.ReceiveAccount.IsNull() && !plan.ReceiveAccount.IsUnknown() {
		var accounts []types.String
		resp.Diagnostics.Append(plan.ReceiveAccount.ElementsAs(ctx, &accounts, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, account := range accounts {
			request.ReceiveAccount = append(request.ReceiveAccount, account.ValueStringPointer())
		}
	}

	if !plan.ReceiveGroup.IsNull() && !plan.ReceiveGroup.IsUnknown() {
		var groups []types.String
		resp.Diagnostics.Append(plan.ReceiveGroup.ElementsAs(ctx, &groups, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, group := range groups {
			request.ReceiveGroup = append(request.ReceiveGroup, group.ValueStringPointer())
		}
	}

	group, groupDiags := resolveUpdatedStrategyGroup(ctx, plan, state)
	resp.Diagnostics.Append(groupDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.GroupId = group.GroupId
	request.Sequence = group.Sequence

	logId := tccommon.GetLogId(ctx)
	response, err := r.Client().UseFwmV20250611Client().ModifyStrategyWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating fwm_strategy_v2",
			fmt.Sprintf("ModifyStrategy failed, strategy_id=%s: %v", ids[0], err),
		)
		return
	}

	if response == nil || response.Response == nil {
		resp.Diagnostics.AddError(
			"Error updating fwm_strategy_v2",
			fmt.Sprintf("ModifyStrategy returned an empty response, strategy_id=%s", ids[0]),
		)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	id := state.Id.ValueString()
	strategies, err := r.read(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading fwm_strategy_v2 after update",
			fmt.Sprintf("id=%s: %v", id, err),
		)
		return
	}
	if strategies.IsNull() {
		log.Printf("[WARN] fwm_strategy_v2 id=%s not found after update, removing it from state.", id)
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("fwm_strategy_v2 %s was not found after update, removing it from state.", id),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	plan.Id = types.StringValue(id)
	plan.Strategy = strategies
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles the deletion of a FWM strategy.
func (r *FwmStrategyV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FwmStrategyV2ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	ids := splitFwmStrategyIds(state.Id.ValueString())
	if len(ids) == 0 {
		resp.Diagnostics.AddError(
			"Error deleting fwm_strategy_v2",
			fmt.Sprintf("Empty strategy id in state: %s", state.Id.ValueString()),
		)
		return
	}

	request := fwmv20250611.NewDeleteStrategyRequest()
	for _, id := range ids {
		request.StrategyIds = append(request.StrategyIds, helper.String(id))
	}

	logId := tccommon.GetLogId(ctx)
	response, err := r.Client().UseFwmV20250611Client().DeleteStrategyWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting fwm_strategy_v2",
			fmt.Sprintf("DeleteStrategy failed, id=%s: %v", state.Id.ValueString(), err),
		)
		return
	}

	if response == nil || response.Response == nil {
		resp.Diagnostics.AddError(
			"Error deleting fwm_strategy_v2",
			fmt.Sprintf("DeleteStrategy returned an empty response, id=%s", state.Id.ValueString()),
		)
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
}

// splitFwmStrategyIds splits the resource ID into the strategy IDs returned by
// CreateStrategy.
func splitFwmStrategyIds(id string) []string {
	ids := make([]string, 0)
	start := 0
	for i := 0; i <= len(id); i++ {
		if i == len(id) || id[i:i+len(tccommon.FILED_SP)] == tccommon.FILED_SP {
			part := id[start:i]
			if part != "" {
				ids = append(ids, part)
			}
			start = i + len(tccommon.FILED_SP)
		}
	}
	return ids
}

// resolveUpdatedStrategyGroup picks the rule group entry that carries
// strategy_id, so that ModifyStrategy only touches the arguments that changed.
func resolveUpdatedStrategyGroup(ctx context.Context, plan, state FwmStrategyV2ResourceModel) (*fwmv20250611.StrategyReq, diag.Diagnostics) {
	var diags diag.Diagnostics

	preItems, ok := strategyGroupsFromList(ctx, plan.PreStrategy, &diags)
	if !ok {
		return nil, diags
	}
	postItems, ok := strategyGroupsFromList(ctx, plan.PostStrategy, &diags)
	if !ok {
		return nil, diags
	}

	var item strategyGroupModel
	switch {
	case !plan.PreStrategy.Equal(state.PreStrategy) && len(preItems) > 0:
		item = preItems[0]
	case !plan.PostStrategy.Equal(state.PostStrategy) && len(postItems) > 0:
		item = postItems[0]
	case len(preItems) == 1:
		item = preItems[0]
	case len(postItems) == 1:
		item = postItems[0]
	case len(postItems) > 0:
		item = postItems[0]
	case len(preItems) > 0:
		item = preItems[0]
	default:
		return nil, diags
	}

	group := &fwmv20250611.StrategyReq{}
	if !item.GroupId.IsNull() && !item.GroupId.IsUnknown() {
		group.GroupId = item.GroupId.ValueStringPointer()
	}
	if !item.Sequence.IsNull() && !item.Sequence.IsUnknown() {
		group.Sequence = item.Sequence.ValueInt64Pointer()
	}
	return group, diags
}

// strategyGroupsFromList converts a pre_strategy / post_strategy list attribute
// into its Go representation.
func strategyGroupsFromList(ctx context.Context, list types.List, diags *diag.Diagnostics) ([]strategyGroupModel, bool) {
	if list.IsNull() || list.IsUnknown() {
		return nil, true
	}
	var items []strategyGroupModel
	d := list.ElementsAs(ctx, &items, false)
	diags.Append(d...)
	if diags.HasError() {
		return nil, false
	}
	return items, true
}

// read fetches every strategy referenced by the resource ID. It returns a null
// list when none of them exists on the cloud side, so the caller can decide to
// remove the resource from state.
func (r *FwmStrategyV2Resource) read(ctx context.Context, id string) (types.List, error) {
	ids := splitFwmStrategyIds(id)
	strategyType := strategyDetailAttrType()

	items := make([]strategyDetailModel, 0, len(ids))
	for _, strategyId := range ids {
		request := fwmv20250611.NewDescribeStrategyRequest()
		request.StrategyId = helper.String(strategyId)

		response, err := r.Client().UseFwmV20250611Client().DescribeStrategyWithContext(ctx, request)
		if err != nil {
			return types.ListNull(strategyType), fmt.Errorf("DescribeStrategy failed, strategy_id=%s: %v", strategyId, err)
		}
		if response == nil || response.Response == nil {
			return types.ListNull(strategyType), fmt.Errorf("DescribeStrategy returned an empty response, strategy_id=%s", strategyId)
		}
		if response.Response.Strategy == nil {
			log.Printf("[WARN] fwm_strategy_v2 strategy_id=%s not found on the cloud side.", strategyId)
			continue
		}
		items = append(items, flattenStrategyDetail(response.Response.Strategy))
	}

	if len(items) == 0 {
		return types.ListNull(strategyType), nil
	}

	list, diags := types.ListValueFrom(ctx, strategyType, items)
	if diags.HasError() {
		return types.ListNull(strategyType), fmt.Errorf("failed to flatten fwm_strategy_v2 details: %v", diags)
	}
	return list, nil
}

// flattenStrategyDetail converts one StrategyResp into its model.
func flattenStrategyDetail(item *fwmv20250611.StrategyResp) strategyDetailModel {
	result := strategyDetailModel{
		StrategyId: types.StringPointerValue(item.StrategyId),
		GroupId:    types.StringPointerValue(item.GroupId),
		GroupName:  types.StringPointerValue(item.GroupName),
		RuleCount:  types.Int64PointerValue(item.RuleCount),
		RuleStatus: types.Int64PointerValue(item.RuleStatus),
		Sequence:   types.Int64PointerValue(item.Sequence),
		ErrMsg:     types.StringPointerValue(item.ErrMsg),
		ErrorType:  types.StringPointerValue(item.ErrorType),
		UpdateTime: types.StringPointerValue(item.UpdateTime),
		CreateBy:   types.StringPointerValue(item.CreateBy),
		UpdateBy:   types.StringPointerValue(item.UpdateBy),
		ExecArea:   types.StringPointerValue(item.ExecArea),
		CreateName: types.StringPointerValue(item.CreateName),
		UpdateName: types.StringPointerValue(item.UpdateName),
		CreateTime: types.StringPointerValue(item.CreateTime),
	}

	accounts := make([]strategyReceiveAccountModel, 0, len(item.ReceiveAccount))
	for _, account := range item.ReceiveAccount {
		if account == nil {
			continue
		}
		members := make([]strategyMemberModel, 0, len(account.Members))
		for _, member := range account.Members {
			if member == nil {
				continue
			}
			members = append(members, strategyMemberModel{
				AppId:    types.StringPointerValue(member.AppId),
				Uin:      types.StringPointerValue(member.Uin),
				Nickname: types.StringPointerValue(member.Nickname),
				MemberId: types.StringPointerValue(member.MemberId),
				NodeName: types.StringPointerValue(member.NodeName),
			})
		}
		accounts = append(accounts, strategyReceiveAccountModel{
			Uin:          types.StringPointerValue(account.Uin),
			Nickname:     types.StringPointerValue(account.Nickname),
			ReceiverType: types.Int64PointerValue(account.ReceiverType),
			Members:      members,
		})
	}
	result.ReceiveAccount = accounts

	return result
}

// strategyDetailAttrType returns the object type of one strategy entry.
func strategyDetailAttrType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"strategy_id": types.StringType,
			"group_id":    types.StringType,
			"group_name":  types.StringType,
			"rule_count":  types.Int64Type,
			"rule_status": types.Int64Type,
			"receive_account": types.ListType{
				ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"uin":           types.StringType,
						"nickname":      types.StringType,
						"receiver_type": types.Int64Type,
						"members": types.ListType{
							ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"app_id":    types.StringType,
									"uin":       types.StringType,
									"nickname":  types.StringType,
									"member_id": types.StringType,
									"node_name": types.StringType,
								},
							},
						},
					},
				},
			},
			"sequence":    types.Int64Type,
			"err_msg":     types.StringType,
			"error_type":  types.StringType,
			"update_time": types.StringType,
			"create_by":   types.StringType,
			"update_by":   types.StringType,
			"exec_area":   types.StringType,
			"create_name": types.StringType,
			"update_name": types.StringType,
			"create_time": types.StringType,
		},
	}
}

// stringUseStateForUnknown keeps a computed string attribute stable across
// applies when the plan value is not yet known.
type stringUseStateForUnknown struct{}

func (stringUseStateForUnknown) Description(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}

func (m stringUseStateForUnknown) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
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
