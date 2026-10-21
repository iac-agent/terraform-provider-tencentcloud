package fwm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

var (
	_ resource.Resource                = &FwmStrategyResource{}
	_ resource.ResourceWithConfigure   = &FwmStrategyResource{}
	_ resource.ResourceWithImportState = &FwmStrategyResource{}
)

// NewFwmStrategyResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewFwmStrategyResource() resource.Resource {
	return &FwmStrategyResource{}
}

// FwmStrategyResource implements resource.Resource for
// tencentcloud_fwm_strategy.
type FwmStrategyResource struct {
	client *sharedmeta.ProviderMeta
}

// FwmStrategyGroupModel maps a pre/post strategy rule group.
type FwmStrategyGroupModel struct {
	GroupId  types.String `tfsdk:"group_id"`
	Sequence types.Int64  `tfsdk:"sequence"`
}

// FwmStrategyMemberModel maps a member of a receive account group.
type FwmStrategyMemberModel struct {
	AppId    types.String `tfsdk:"app_id"`
	Uin      types.String `tfsdk:"uin"`
	Nickname types.String `tfsdk:"nickname"`
	MemberId types.String `tfsdk:"member_id"`
	NodeName types.String `tfsdk:"node_name"`
}

// FwmStrategyReceiveAccountModel maps a receive account detail returned by
// DescribeStrategy.
type FwmStrategyReceiveAccountModel struct {
	Uin          types.String `tfsdk:"uin"`
	Nickname     types.String `tfsdk:"nickname"`
	ReceiverType types.Int64  `tfsdk:"receiver_type"`
	Members      types.List   `tfsdk:"members"`
}

// FwmStrategyResourceModel maps the schema attributes.
type FwmStrategyResourceModel struct {
	Id                   types.String           `tfsdk:"id"`
	Product              types.String           `tfsdk:"product"`
	ReceiveAccount       types.List             `tfsdk:"receive_account"`
	PreStrategy          *FwmStrategyGroupModel `tfsdk:"pre_strategy"`
	PostStrategy         *FwmStrategyGroupModel `tfsdk:"post_strategy"`
	ReceiveGroup         types.List             `tfsdk:"receive_group"`
	GroupId              types.String           `tfsdk:"group_id"`
	Sequence             types.Int64            `tfsdk:"sequence"`
	StrategyId           types.String           `tfsdk:"strategy_id"`
	GroupName            types.String           `tfsdk:"group_name"`
	RuleCount            types.Int64            `tfsdk:"rule_count"`
	RuleStatus           types.Int64            `tfsdk:"rule_status"`
	ReceiveAccountDetail types.List             `tfsdk:"receive_account_detail"`
	ErrMsg               types.String           `tfsdk:"err_msg"`
	ErrorType            types.String           `tfsdk:"error_type"`
	UpdateTime           types.String           `tfsdk:"update_time"`
	CreateBy             types.String           `tfsdk:"create_by"`
	UpdateBy             types.String           `tfsdk:"update_by"`
	ExecArea             types.String           `tfsdk:"exec_area"`
	CreateName           types.String           `tfsdk:"create_name"`
	UpdateName           types.String           `tfsdk:"update_name"`
	CreateTime           types.String           `tfsdk:"create_time"`
}

func (r *FwmStrategyResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_fwm_strategy"
}

func (r *FwmStrategyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage a FWM (Firewall Management) strategy.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "ID of the strategy, same as `strategy_id`.",
			},
			"product": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Product type. This field is only used when creating and cannot be changed afterwards.",
			},
			"receive_account": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of accounts (uins) receiving the dispatch rules.",
			},
			"pre_strategy": schema.SingleNestedAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Description: "Pre-zone rule group. This field is only used when creating and cannot be changed afterwards.",
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
			"post_strategy": schema.SingleNestedAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Description: "Post-zone rule group. This field is only used when creating and cannot be changed afterwards.",
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
			"receive_group": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of account groups receiving the dispatch rules.",
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Rule group ID used when updating the strategy.",
			},
			"sequence": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Priority of the rule group used when updating the strategy.",
			},
			"strategy_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "ID of the strategy.",
			},
			"group_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the rule group.",
			},
			"rule_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of rules in the strategy.",
			},
			"rule_status": schema.Int64Attribute{
				Computed:    true,
				Description: "Status of the strategy.",
			},
			"receive_account_detail": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Details of the accounts receiving the dispatch rules.",
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
							Description: "Type of the receiver. `0` means account uin and `1` means account group.",
						},
						"members": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Members of the account group. Only returned when `receiver_type` is `1`.",
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
										Description: "ID of the member. Empty in personal scenarios.",
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
			"err_msg": schema.StringAttribute{
				Computed:    true,
				Description: "Reason why the dispatch failed.",
			},
			"error_type": schema.StringAttribute{
				Computed:    true,
				Description: "Type of the dispatch failure reason.",
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "Last updated time of the strategy.",
			},
			"create_by": schema.StringAttribute{
				Computed:    true,
				Description: "Creator of the strategy.",
			},
			"update_by": schema.StringAttribute{
				Computed:    true,
				Description: "Last updater of the strategy.",
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
				Description: "Name of the last updater.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the strategy.",
			},
		},
	}
}

// Configure receives the *sharedmeta.ProviderMeta populated by
// framework.Provider.Configure and stores the shared client.
func (r *FwmStrategyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	meta, ok := req.ProviderData.(*sharedmeta.ProviderMeta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData type",
			fmt.Sprintf("Expected *sharedmeta.ProviderMeta, got: %T. "+
				"This is a bug in the provider; please report it.", req.ProviderData),
		)
		return
	}

	r.client = meta
}

// Create handles `terraform apply` for new resources.
func (r *FwmStrategyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FwmStrategyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := sdk.NewCreateStrategyRequest()
	request.Product = plan.Product.ValueStringPointer()
	if !plan.ReceiveAccount.IsNull() && !plan.ReceiveAccount.IsUnknown() {
		var accounts []string
		resp.Diagnostics.Append(plan.ReceiveAccount.ElementsAs(ctx, &accounts, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, v := range accounts {
			request.ReceiveAccount = append(request.ReceiveAccount, helper.String(v))
		}
	}
	if !plan.ReceiveGroup.IsNull() && !plan.ReceiveGroup.IsUnknown() {
		var groups []string
		resp.Diagnostics.Append(plan.ReceiveGroup.ElementsAs(ctx, &groups, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, v := range groups {
			request.ReceiveGroup = append(request.ReceiveGroup, helper.String(v))
		}
	}
	if plan.PreStrategy != nil {
		pre := &sdk.StrategyReq{}
		if !plan.PreStrategy.GroupId.IsNull() && !plan.PreStrategy.GroupId.IsUnknown() {
			pre.GroupId = plan.PreStrategy.GroupId.ValueStringPointer()
		}
		if !plan.PreStrategy.Sequence.IsNull() && !plan.PreStrategy.Sequence.IsUnknown() {
			pre.Sequence = helper.Int64(plan.PreStrategy.Sequence.ValueInt64())
		}
		request.PreStrategy = []*sdk.StrategyReq{pre}
	}
	if plan.PostStrategy != nil {
		post := &sdk.StrategyReq{}
		if !plan.PostStrategy.GroupId.IsNull() && !plan.PostStrategy.GroupId.IsUnknown() {
			post.GroupId = plan.PostStrategy.GroupId.ValueStringPointer()
		}
		if !plan.PostStrategy.Sequence.IsNull() && !plan.PostStrategy.Sequence.IsUnknown() {
			post.Sequence = helper.Int64(plan.PostStrategy.Sequence.ValueInt64())
		}
		request.PostStrategy = []*sdk.StrategyReq{post}
	}

	logId := "fwm_strategy_create"
	response, err := r.client.Client.UseFwmClient().CreateStrategy(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating fwm_strategy resource",
			fmt.Sprintf("CreateStrategy failed: %v", err),
		)
		return
	}
	if response == nil || response.Response == nil || len(response.Response.StrategyIds) == 0 || response.Response.StrategyIds[0] == nil || *response.Response.StrategyIds[0] == "" {
		log.Printf("[CRUD] %s: CreateStrategy returned empty strategy id", logId)
		resp.Diagnostics.AddError(
			"Error creating fwm_strategy resource",
			fmt.Sprintf("CreateStrategy returned empty strategy id, logId=%s", logId),
		)
		return
	}

	plan.Id = types.StringValue(*response.Response.StrategyIds[0])

	// Read back immediately after creation to keep state consistent with the cloud side.
	item, err := r.read(ctx, plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading fwm_strategy resource after create",
			fmt.Sprintf("DescribeStrategy failed for strategy %s: %v", plan.Id.ValueString(), err),
		)
		return
	}
	r.refreshModel(&plan, item)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *FwmStrategyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FwmStrategyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.read(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading fwm_strategy resource",
			fmt.Sprintf("DescribeStrategy failed for strategy %s: %v", state.Id.ValueString(), err),
		)
		return
	}

	// No longer exists on the cloud side => remove from state so it can be recreated.
	if item == nil {
		log.Printf("[CRUD] fwm_strategy id=%s not found on cloud side, removing from state", state.Id.ValueString())
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("fwm_strategy resource %s not found, removing from state.", state.Id.ValueString()),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	r.refreshModel(&state, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates.
func (r *FwmStrategyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FwmStrategyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := sdk.NewModifyStrategyRequest()
	request.StrategyId = plan.Id.ValueStringPointer()
	// Take update values from the Plan (not the old State values)
	if !plan.ReceiveAccount.IsNull() && !plan.ReceiveAccount.IsUnknown() {
		var accounts []string
		resp.Diagnostics.Append(plan.ReceiveAccount.ElementsAs(ctx, &accounts, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, v := range accounts {
			request.ReceiveAccount = append(request.ReceiveAccount, helper.String(v))
		}
	}
	if !plan.ReceiveGroup.IsNull() && !plan.ReceiveGroup.IsUnknown() {
		var groups []string
		resp.Diagnostics.Append(plan.ReceiveGroup.ElementsAs(ctx, &groups, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, v := range groups {
			request.ReceiveGroup = append(request.ReceiveGroup, helper.String(v))
		}
	}
	if !plan.Sequence.IsNull() && !plan.Sequence.IsUnknown() {
		request.Sequence = helper.Int64(plan.Sequence.ValueInt64())
	}
	if !plan.GroupId.IsNull() && !plan.GroupId.IsUnknown() {
		request.GroupId = plan.GroupId.ValueStringPointer()
	}

	_, err := r.client.Client.UseFwmClient().ModifyStrategy(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating fwm_strategy resource",
			fmt.Sprintf("ModifyStrategy failed for strategy %s: %v", plan.Id.ValueString(), err),
		)
		return
	}

	item, err := r.read(ctx, plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading fwm_strategy resource after update",
			fmt.Sprintf("DescribeStrategy failed for strategy %s: %v", plan.Id.ValueString(), err),
		)
		return
	}
	r.refreshModel(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *FwmStrategyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FwmStrategyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := sdk.NewDeleteStrategyRequest()
	request.StrategyIds = []*string{state.Id.ValueStringPointer()}

	_, err := r.client.Client.UseFwmClient().DeleteStrategy(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting fwm_strategy resource",
			fmt.Sprintf("DeleteStrategy failed for strategy %s: %v", state.Id.ValueString(), err),
		)
		return
	}
}

// ImportState supports `terraform import` using the strategy id passthrough.
func (r *FwmStrategyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// read fetches the strategy by ID; returns nil when cloud-side no longer has it.
func (r *FwmStrategyResource) read(ctx context.Context, id string) (*sdk.StrategyResp, error) {
	request := sdk.NewDescribeStrategyRequest()
	request.StrategyId = helper.String(id)

	response, err := r.client.Client.UseFwmClient().DescribeStrategy(request)
	if err != nil {
		return nil, err
	}
	if response == nil || response.Response == nil || response.Response.Strategy == nil {
		log.Printf("[CRUD] fwm_strategy id=%s: DescribeStrategy returned empty strategy", id)
		return nil, nil
	}
	return response.Response.Strategy, nil
}

var fwmStrategyMemberAttrTypes = map[string]attr.Type{
	"app_id":    types.StringType,
	"uin":       types.StringType,
	"nickname":  types.StringType,
	"member_id": types.StringType,
	"node_name": types.StringType,
}

var fwmStrategyReceiveAccountAttrTypes = map[string]attr.Type{
	"uin":           types.StringType,
	"nickname":      types.StringType,
	"receiver_type": types.Int64Type,
	"members":       types.ListType{ElemType: types.ObjectType{AttrTypes: fwmStrategyMemberAttrTypes}},
}

// refreshModel copies API item fields into the model (only cloud-owned fields).
func (r *FwmStrategyResource) refreshModel(m *FwmStrategyResourceModel, item *sdk.StrategyResp) {
	if item == nil {
		return
	}
	m.Id = types.StringPointerValue(item.StrategyId)
	m.StrategyId = types.StringPointerValue(item.StrategyId)
	m.GroupId = types.StringPointerValue(item.GroupId)
	m.Sequence = types.Int64PointerValue(item.Sequence)
	m.GroupName = types.StringPointerValue(item.GroupName)
	m.RuleCount = types.Int64PointerValue(item.RuleCount)
	m.RuleStatus = types.Int64PointerValue(item.RuleStatus)
	m.ErrMsg = types.StringPointerValue(item.ErrMsg)
	m.ErrorType = types.StringPointerValue(item.ErrorType)
	m.UpdateTime = types.StringPointerValue(item.UpdateTime)
	m.CreateBy = types.StringPointerValue(item.CreateBy)
	m.UpdateBy = types.StringPointerValue(item.UpdateBy)
	m.ExecArea = types.StringPointerValue(item.ExecArea)
	m.CreateName = types.StringPointerValue(item.CreateName)
	m.UpdateName = types.StringPointerValue(item.UpdateName)
	m.CreateTime = types.StringPointerValue(item.CreateTime)

	accounts := make([]FwmStrategyReceiveAccountModel, 0, len(item.ReceiveAccount))
	for _, a := range item.ReceiveAccount {
		if a == nil {
			continue
		}
		account := FwmStrategyReceiveAccountModel{
			Uin:          types.StringPointerValue(a.Uin),
			Nickname:     types.StringPointerValue(a.Nickname),
			ReceiverType: types.Int64PointerValue(a.ReceiverType),
		}
		members := make([]FwmStrategyMemberModel, 0, len(a.Members))
		for _, mem := range a.Members {
			if mem == nil {
				continue
			}
			members = append(members, FwmStrategyMemberModel{
				AppId:    types.StringPointerValue(mem.AppId),
				Uin:      types.StringPointerValue(mem.Uin),
				Nickname: types.StringPointerValue(mem.Nickname),
				MemberId: types.StringPointerValue(mem.MemberId),
				NodeName: types.StringPointerValue(mem.NodeName),
			})
		}
		memberList, diags := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: fwmStrategyMemberAttrTypes}, members)
		if diags.HasError() {
			account.Members = types.ListNull(types.ObjectType{AttrTypes: fwmStrategyMemberAttrTypes})
		} else {
			account.Members = memberList
		}
		accounts = append(accounts, account)
	}
	accountList, diags := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: fwmStrategyReceiveAccountAttrTypes}, accounts)
	if !diags.HasError() {
		m.ReceiveAccountDetail = accountList
	}
}
