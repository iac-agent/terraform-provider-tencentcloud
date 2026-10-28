package fwm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// Compile-time interface assertions: the base resource.Resource interface is
// already checked implicitly by the factory return value, so only the
// extension interfaces are asserted here.
var (
	_ resource.ResourceWithConfigure   = &FwmVpcAclRuleResource{}
	_ resource.ResourceWithImportState = &FwmVpcAclRuleResource{}
	_ resource.ResourceWithModifyPlan  = &FwmVpcAclRuleResource{}
)

// NewFwmVpcAclRuleResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewFwmVpcAclRuleResource() resource.Resource {
	return &FwmVpcAclRuleResource{}
}

// FwmVpcAclRuleResource implements resource.Resource for
// tencentcloud_fwm_vpc_acl_rule.
type FwmVpcAclRuleResource struct {
	fw.ResourceWithConfigure
}

// FwmVpcAclRuleRuleModel maps one VPC border rule under the `rules` attribute.
type FwmVpcAclRuleRuleModel struct {
	SourceContent     types.String `tfsdk:"source_content"`
	SourceType        types.String `tfsdk:"source_type"`
	SourceName        types.String `tfsdk:"source_name"`
	DestContent       types.String `tfsdk:"dest_content"`
	DestType          types.String `tfsdk:"dest_type"`
	DestName          types.String `tfsdk:"dest_name"`
	Protocol          types.String `tfsdk:"protocol"`
	Port              types.String `tfsdk:"port"`
	ParamTemplateId   types.String `tfsdk:"param_template_id"`
	ParamTemplateName types.String `tfsdk:"param_template_name"`
	RuleAction        types.String `tfsdk:"rule_action"`
	OrderIndex        types.Int64  `tfsdk:"order_index"`
	EdgeId            types.String `tfsdk:"edge_id"`
	FwGroupId         types.String `tfsdk:"fw_group_id"`
	FwGroupName       types.String `tfsdk:"fw_group_name"`
	RuleId            types.String `tfsdk:"rule_id"`
	Description       types.String `tfsdk:"description"`
	BelongMemberId    types.String `tfsdk:"belong_member_id"`
	BelongMember      types.Object `tfsdk:"belong_member"`
	Sequence          types.Int64  `tfsdk:"sequence"`
	IpVersion         types.String `tfsdk:"ip_version"`
	Invalid           types.Int64  `tfsdk:"invalid"`
	CreateTime        types.String `tfsdk:"create_time"`
	UpdateTime        types.String `tfsdk:"update_time"`
	DnsParseCnt       types.Int64  `tfsdk:"dns_parse_cnt"`
	HitCount          types.Int64  `tfsdk:"hit_count"`
	LastHitTime       types.String `tfsdk:"last_hit_time"`
}

// FwmVpcAclRuleBelongMemberModel maps the nested `belong_member` object.
type FwmVpcAclRuleBelongMemberModel struct {
	AppId    types.String `tfsdk:"app_id"`
	Uin      types.String `tfsdk:"uin"`
	Nickname types.String `tfsdk:"nickname"`
	MemberId types.String `tfsdk:"member_id"`
	NodeName types.String `tfsdk:"node_name"`
}

// FwmVpcAclRuleResourceModel maps the schema attributes.
type FwmVpcAclRuleResourceModel struct {
	Id            types.String `tfsdk:"id"`
	GroupId       types.String `tfsdk:"group_id"`
	FwGroupId     types.String `tfsdk:"fw_group_id"`
	RuleAction    types.String `tfsdk:"rule_action"`
	Description   types.String `tfsdk:"description"`
	OrderIndex    types.Int64  `tfsdk:"order_index"`
	RuleId        types.String `tfsdk:"rule_id"`
	RuleIds       types.List   `tfsdk:"rule_ids"`
	AllTotalCount types.Int64  `tfsdk:"all_total_count"`
	Rules         types.List   `tfsdk:"rules"`
}

func (r *FwmVpcAclRuleResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_fwm_vpc_acl_rule"
}

func (r *FwmVpcAclRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a FWM VPC border ACL rule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier in the form `<group_id>#<rule_id>`.",
				PlanModifiers: []planmodifier.String{
					stringUseStateForUnknownFwmVpcAclRule{},
				},
			},
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the VPC border ACL rule group the rule belongs to. Changing it forces a new rule.",
				PlanModifiers: []planmodifier.String{
					stringRequiresReplaceFwmVpcAclRule{},
				},
			},
			"fw_group_id": schema.StringAttribute{
				Required:    true,
				Description: "Firewall instance ID that the rule applies to. `ALL` means global, `ccn-xxx` is a Cloud Connect Network instance and `cfwg-xxx` is a firewall group instance. Changing it forces a new rule.",
				PlanModifiers: []planmodifier.String{
					stringRequiresReplaceFwmVpcAclRule{},
				},
			},
			"rule_action": schema.StringAttribute{
				Required:    true,
				Description: "Action of the rule. Valid values: `accept`, `drop`, `log`. Changing it forces a new rule.",
				PlanModifiers: []planmodifier.String{
					stringRequiresReplaceFwmVpcAclRule{},
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Description of the rule.",
			},
			"order_index": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Priority of the rule, starting from `1`.",
				PlanModifiers: []planmodifier.Int64{
					int64UseStateForUnknownFwmVpcAclRule{},
				},
			},
			"rule_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the VPC border ACL rule managed by this resource.",
				PlanModifiers: []planmodifier.String{
					stringUseStateForUnknownFwmVpcAclRule{},
				},
			},
			"rule_ids": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "IDs of the VPC border ACL rules created by this resource.",
			},
			"all_total_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Total number of rules of the rule group without any filter.",
				PlanModifiers: []planmodifier.Int64{
					int64UseStateForUnknownFwmVpcAclRule{},
				},
			},
			"rules": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Full list of rules of the rule group as returned by the cloud API, including the rule managed by this resource.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"source_content": schema.StringAttribute{
							Computed:    true,
							Description: "Source address content of the rule.",
						},
						"source_type": schema.StringAttribute{
							Computed:    true,
							Description: "Source type of the rule. Valid values: `ip`, `url`, `template`, `instance`, `tag`.",
						},
						"source_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the source asset.",
						},
						"dest_content": schema.StringAttribute{
							Computed:    true,
							Description: "Destination address content of the rule.",
						},
						"dest_type": schema.StringAttribute{
							Computed:    true,
							Description: "Destination type of the rule. Valid values: `ip`, `url`, `template`, `instance`, `tag`.",
						},
						"dest_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the destination asset.",
						},
						"protocol": schema.StringAttribute{
							Computed:    true,
							Description: "Protocol of the rule, e.g. `TCP`, `UDP`, `ICMP`, `ANY`, `HTTP`, `HTTPS`, `DNS`, `FTP`.",
						},
						"port": schema.StringAttribute{
							Computed:    true,
							Description: "Port of the rule. It is empty when the protocol is `ICMP`.",
						},
						"param_template_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the port template.",
						},
						"param_template_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the port template.",
						},
						"rule_action": schema.StringAttribute{
							Computed:    true,
							Description: "Action of the rule. Valid values: `accept`, `drop`, `log`.",
						},
						"order_index": schema.Int64Attribute{
							Computed:    true,
							Description: "Priority of the rule, starting from `1`.",
						},
						"edge_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the border firewall. `ALL` means global and `CFWS-xxx` means the specified border.",
						},
						"fw_group_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the firewall instance.",
						},
						"fw_group_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the firewall group.",
						},
						"rule_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the rule.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Description of the rule.",
						},
						"belong_member_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the member account the rule belongs to.",
						},
						"belong_member": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Member account the rule belongs to.",
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
									Description: "ID of the member. It is empty for personal accounts.",
								},
								"node_name": schema.StringAttribute{
									Computed:    true,
									Description: "Department the member belongs to.",
								},
							},
						},
						"sequence": schema.Int64Attribute{
							Computed:    true,
							Description: "Priority of the rule.",
						},
						"ip_version": schema.StringAttribute{
							Computed:    true,
							Description: "IP version of the rule. Valid values: `ipv4`, `ipv6`.",
						},
						"invalid": schema.Int64Attribute{
							Computed:    true,
							Description: "Whether the rule is invalid. `0` means valid, `1` means invalid.",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "Creation time of the rule.",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "Last modification time of the rule.",
						},
						"dns_parse_cnt": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of parsed domain names of the rule.",
						},
						"hit_count": schema.Int64Attribute{
							Computed:    true,
							Description: "Accumulated hit count of the rule across all delivered member accounts.",
						},
						"last_hit_time": schema.StringAttribute{
							Computed:    true,
							Description: "Last hit time of the rule. It is `-` when the rule has never been hit.",
						},
					},
				},
			},
		},
	}
}

// ModifyPlan rejects in-place changes of the optional top-level attributes
// when the cloud API has no way to update them on an existing rule.
func (r *FwmVpcAclRuleResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var plan FwmVpcAclRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state FwmVpcAclRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.checkImmutablePlanChanges(plan, state, resp)
}

// checkImmutablePlanChanges reports an error when an attribute that the cloud
// API cannot update on an existing rule is modified in place.
func (r *FwmVpcAclRuleResource) checkImmutablePlanChanges(plan, state FwmVpcAclRuleResourceModel, resp *resource.ModifyPlanResponse) {
	if !plan.Description.Equal(state.Description) {
		resp.Diagnostics.AddError(
			"Attribute description cannot be updated in place",
			"The FWM VPC border ACL rule does not support updating description. "+
				"Remove the resource and let Terraform recreate it instead.",
		)
		return
	}

	if !plan.OrderIndex.Equal(state.OrderIndex) {
		resp.Diagnostics.AddError(
			"Attribute order_index cannot be updated in place",
			"The FWM VPC border ACL rule does not support updating order_index. "+
				"Remove the resource and let Terraform recreate it instead.",
		)
	}
}

// Create handles `terraform apply` for a new VPC border ACL rule.
func (r *FwmVpcAclRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FwmVpcAclRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule := &fwmv20250611.VpcAclRule{
		FwGroupId:     plan.FwGroupId.ValueStringPointer(),
		RuleAction:    plan.RuleAction.ValueStringPointer(),
		SourceType:    helper.String("ip"),
		SourceContent: helper.String("0.0.0.0/0"),
		DestType:      helper.String("ip"),
		DestContent:   helper.String("0.0.0.0/0"),
		Protocol:      helper.String("ANY"),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		rule.Description = plan.Description.ValueStringPointer()
	}
	if !plan.OrderIndex.IsNull() && !plan.OrderIndex.IsUnknown() {
		rule.OrderIndex = helper.Int64(plan.OrderIndex.ValueInt64())
	}

	groupId := plan.GroupId.ValueString()
	service := NewFwmService(r.Client())
	ruleIds, err := service.CreateFwmVpcAclRule(ctx, groupId, []*fwmv20250611.VpcAclRule{rule})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating FWM VPC border ACL rule",
			fmt.Sprintf("CreateVpcAclRule failed for fwm_vpc_acl_rule group_id %q: %v", groupId, err),
		)
		return
	}

	if len(ruleIds) == 0 || ruleIds[0] == nil || *ruleIds[0] == "" {
		resp.Diagnostics.AddError(
			"Error creating FWM VPC border ACL rule",
			fmt.Sprintf("CreateVpcAclRule returned an empty rule id for fwm_vpc_acl_rule group_id %q", groupId),
		)
		return
	}

	ruleId := *ruleIds[0]
	plan.Id = types.StringValue(strings.Join([]string{groupId, ruleId}, tccommon.FILED_SP))
	plan.RuleId = types.StringValue(ruleId)

	ruleIdsList, diags := types.ListValueFrom(ctx, types.StringType, []string{ruleId})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.RuleIds = ruleIdsList

	// Persist the id before polling so a failed poll still leaves the rule in
	// state and `terraform destroy` is able to clean it up.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.read(ctx, groupId, ruleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading FWM VPC border ACL rule after create",
			fmt.Sprintf("DescribeVpcAclRules failed for fwm_vpc_acl_rule %q: %v", plan.Id.ValueString(), err),
		)
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"FWM VPC border ACL rule not found after create",
			fmt.Sprintf("fwm_vpc_acl_rule %q cannot be read back after create", plan.Id.ValueString()),
		)
		return
	}

	plan.RuleId = types.StringPointerValue(item.RuleId)
	r.refreshModel(ctx, &plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *FwmVpcAclRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FwmVpcAclRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId, ruleId, ok := splitFwmVpcAclRuleId(state.Id.ValueString())
	if !ok {
		groupId = state.GroupId.ValueString()
		ruleId = state.RuleId.ValueString()
	}
	if groupId == "" || ruleId == "" {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("fwm_vpc_acl_rule id is broken, id is %q (expected `<group_id>#<rule_id>`)", state.Id.ValueString()),
		)
		return
	}

	item, err := r.read(ctx, groupId, ruleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading FWM VPC border ACL rule",
			fmt.Sprintf("DescribeVpcAclRules failed for fwm_vpc_acl_rule %q: %v", state.Id.ValueString(), err),
		)
		return
	}

	if item == nil {
		logMissingFwmVpcAclRule(state.Id.ValueString())
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("fwm_vpc_acl_rule %s not found, removing from state.", state.Id.ValueString()),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	state.GroupId = types.StringValue(groupId)
	state.RuleId = types.StringPointerValue(item.RuleId)
	r.refreshModel(ctx, &state, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update is required by the resource.Resource interface but no attribute of
// this resource is mutable in place, so it only re-reads the remote state.
func (r *FwmVpcAclRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FwmVpcAclRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state FwmVpcAclRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId := state.GroupId.ValueString()
	ruleId := state.RuleId.ValueString()

	item, err := r.read(ctx, groupId, ruleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading FWM VPC border ACL rule after update",
			fmt.Sprintf("DescribeVpcAclRules failed for fwm_vpc_acl_rule %q: %v", state.Id.ValueString(), err),
		)
		return
	}
	if item == nil {
		logMissingFwmVpcAclRule(state.Id.ValueString())
		resp.State.RemoveResource(ctx)
		return
	}

	plan.Id = state.Id
	plan.GroupId = types.StringValue(groupId)
	plan.RuleId = types.StringPointerValue(item.RuleId)
	r.refreshModel(ctx, &plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles `terraform destroy`.
func (r *FwmVpcAclRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FwmVpcAclRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId := state.GroupId.ValueString()
	ruleId := state.RuleId.ValueString()
	if groupId == "" || ruleId == "" {
		var ok bool
		groupId, ruleId, ok = splitFwmVpcAclRuleId(state.Id.ValueString())
		if !ok {
			resp.Diagnostics.AddError(
				"Invalid resource ID",
				fmt.Sprintf("fwm_vpc_acl_rule id is broken, id is %q (expected `<group_id>#<rule_id>`)", state.Id.ValueString()),
			)
			return
		}
	}

	service := NewFwmService(r.Client())
	if err := service.DeleteFwmVpcAclRule(ctx, groupId, []string{ruleId}); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting FWM VPC border ACL rule",
			fmt.Sprintf("DeleteVpcAclRule failed for fwm_vpc_acl_rule %q: %v", state.Id.ValueString(), err),
		)
		return
	}
}

// ImportState supports `terraform import` using `<group_id>#<rule_id>`.
func (r *FwmVpcAclRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	groupId, ruleId, ok := splitFwmVpcAclRuleId(req.ID)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in the form `<group_id>%s<rule_id>`, got %q", tccommon.FILED_SP, req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_id"), groupId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("rule_id"), ruleId)...)
}

// read fetches one rule of the group; it returns nil when the cloud side no
// longer has it.
func (r *FwmVpcAclRuleResource) read(ctx context.Context, groupId, ruleId string) (*fwmv20250611.VpcAclRuleResp, error) {
	service := NewFwmService(r.Client())
	return service.DescribeFwmVpcAclRuleById(ctx, groupId, ruleId)
}

// refreshModel copies the cloud-owned fields of one rule into the model.
func (r *FwmVpcAclRuleResource) refreshModel(ctx context.Context, m *FwmVpcAclRuleResourceModel, item *fwmv20250611.VpcAclRuleResp) {
	if item == nil {
		return
	}

	m.Description = types.StringPointerValue(item.Description)
	m.OrderIndex = types.Int64PointerValue(item.Sequence)
	if item.FwGroupId != nil {
		m.FwGroupId = types.StringPointerValue(item.FwGroupId)
	}
	if item.RuleAction != nil {
		m.RuleAction = types.StringPointerValue(item.RuleAction)
	}

	belongMember := types.ObjectNull(fwmVpcAclRuleBelongMemberAttrTypes())
	belongMemberId := types.StringNull()
	if item.BelongMember != nil {
		belongMemberId = types.StringPointerValue(item.BelongMember.MemberId)
		obj, diags := types.ObjectValueFrom(ctx, fwmVpcAclRuleBelongMemberAttrTypes(), FwmVpcAclRuleBelongMemberModel{
			AppId:    types.StringPointerValue(item.BelongMember.AppId),
			Uin:      types.StringPointerValue(item.BelongMember.Uin),
			Nickname: types.StringPointerValue(item.BelongMember.Nickname),
			MemberId: types.StringPointerValue(item.BelongMember.MemberId),
			NodeName: types.StringPointerValue(item.BelongMember.NodeName),
		})
		if diags.HasError() {
			belongMember = types.ObjectNull(fwmVpcAclRuleBelongMemberAttrTypes())
		} else {
			belongMember = obj
		}
	}

	rule := FwmVpcAclRuleRuleModel{
		SourceContent:     types.StringPointerValue(item.SourceContent),
		SourceType:        types.StringPointerValue(item.SourceType),
		SourceName:        types.StringPointerValue(item.SourceName),
		DestContent:       types.StringPointerValue(item.DestContent),
		DestType:          types.StringPointerValue(item.DestType),
		DestName:          types.StringPointerValue(item.DestName),
		Protocol:          types.StringPointerValue(item.Protocol),
		Port:              types.StringPointerValue(item.Port),
		ParamTemplateId:   types.StringPointerValue(item.ParamTemplateId),
		ParamTemplateName: types.StringPointerValue(item.ParamTemplateName),
		RuleAction:        types.StringPointerValue(item.RuleAction),
		OrderIndex:        types.Int64PointerValue(item.Sequence),
		EdgeId:            types.StringPointerValue(item.EdgeId),
		FwGroupId:         types.StringPointerValue(item.FwGroupId),
		FwGroupName:       types.StringPointerValue(item.FwGroupName),
		RuleId:            types.StringPointerValue(item.RuleId),
		Description:       types.StringPointerValue(item.Description),
		BelongMemberId:    belongMemberId,
		BelongMember:      belongMember,
		Sequence:          types.Int64PointerValue(item.Sequence),
		IpVersion:         types.StringPointerValue(item.IpVersion),
		Invalid:           types.Int64PointerValue(item.Invalid),
		CreateTime:        types.StringPointerValue(item.CreateTime),
		UpdateTime:        types.StringPointerValue(item.UpdateTime),
		DnsParseCnt:       types.Int64PointerValue(item.DnsParseCnt),
		HitCount:          types.Int64PointerValue(item.HitCount),
		LastHitTime:       types.StringPointerValue(item.LastHitTime),
	}

	ruleList, diags := types.ListValueFrom(ctx, fwmVpcAclRuleRuleObjectType(), []FwmVpcAclRuleRuleModel{rule})
	if diags.HasError() {
		return
	}
	m.Rules = ruleList
}

// logMissingFwmVpcAclRule keeps the resource id in the log before the state
// entry is dropped, so the failed read can still be located.
func logMissingFwmVpcAclRule(id string) {
	log.Printf("[WARN] fwm_vpc_acl_rule %s not found, removing it from state.", id)
}

// splitFwmVpcAclRuleId parses `<group_id>#<rule_id>`.
func splitFwmVpcAclRuleId(id string) (string, string, bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// fwmVpcAclRuleBelongMemberAttrTypes returns the attribute types of the
// nested `belong_member` object.
func fwmVpcAclRuleBelongMemberAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"app_id":    types.StringType,
		"uin":       types.StringType,
		"nickname":  types.StringType,
		"member_id": types.StringType,
		"node_name": types.StringType,
	}
}

// fwmVpcAclRuleRuleObjectType returns the object type of one `rules` element.
func fwmVpcAclRuleRuleObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"source_content":      types.StringType,
			"source_type":         types.StringType,
			"source_name":         types.StringType,
			"dest_content":        types.StringType,
			"dest_type":           types.StringType,
			"dest_name":           types.StringType,
			"protocol":            types.StringType,
			"port":                types.StringType,
			"param_template_id":   types.StringType,
			"param_template_name": types.StringType,
			"rule_action":         types.StringType,
			"order_index":         types.Int64Type,
			"edge_id":             types.StringType,
			"fw_group_id":         types.StringType,
			"fw_group_name":       types.StringType,
			"rule_id":             types.StringType,
			"description":         types.StringType,
			"belong_member_id":    types.StringType,
			"belong_member": types.ObjectType{
				AttrTypes: fwmVpcAclRuleBelongMemberAttrTypes(),
			},
			"sequence":      types.Int64Type,
			"ip_version":    types.StringType,
			"invalid":       types.Int64Type,
			"create_time":   types.StringType,
			"update_time":   types.StringType,
			"dns_parse_cnt": types.Int64Type,
			"hit_count":     types.Int64Type,
			"last_hit_time": types.StringType,
		},
	}
}

// stringRequiresReplaceFwmVpcAclRule forces replacement of tencentcloud_fwm_vpc_acl_rule
// when a string attribute changes.
type stringRequiresReplaceFwmVpcAclRule struct{}

func (stringRequiresReplaceFwmVpcAclRule) Description(_ context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}

func (stringRequiresReplaceFwmVpcAclRule) MarkdownDescription(ctx context.Context) string {
	return stringRequiresReplaceFwmVpcAclRule{}.Description(ctx)
}

func (stringRequiresReplaceFwmVpcAclRule) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() {
		return
	}
	if req.PlanValue.IsNull() {
		return
	}
	if !req.PlanValue.Equal(req.StateValue) {
		resp.RequiresReplace = true
	}
}

// stringUseStateForUnknownFwmVpcAclRule keeps the prior state value when the
// planned value is unknown, so computed strings stay stable across refreshes.
type stringUseStateForUnknownFwmVpcAclRule struct{}

func (stringUseStateForUnknownFwmVpcAclRule) Description(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}

func (stringUseStateForUnknownFwmVpcAclRule) MarkdownDescription(ctx context.Context) string {
	return stringUseStateForUnknownFwmVpcAclRule{}.Description(ctx)
}

func (stringUseStateForUnknownFwmVpcAclRule) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() {
		return
	}
	if !req.PlanValue.IsUnknown() {
		return
	}
	resp.PlanValue = req.StateValue
}

// int64UseStateForUnknownFwmVpcAclRule keeps the prior state value when the
// planned value is unknown, so computed integers stay stable across refreshes.
type int64UseStateForUnknownFwmVpcAclRule struct{}

func (int64UseStateForUnknownFwmVpcAclRule) Description(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}

func (int64UseStateForUnknownFwmVpcAclRule) MarkdownDescription(ctx context.Context) string {
	return int64UseStateForUnknownFwmVpcAclRule{}.Description(ctx)
}

func (int64UseStateForUnknownFwmVpcAclRule) PlanModifyInt64(_ context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if req.StateValue.IsNull() {
		return
	}
	if !req.PlanValue.IsUnknown() {
		return
	}
	resp.PlanValue = req.StateValue
}
