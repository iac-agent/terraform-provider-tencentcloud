package fwm

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// Compile-time interface assertions. The base resource.Resource interface is
// checked implicitly by the factory return value.
var (
	_ resource.ResourceWithConfigure   = &FwmVpcAclRule2Resource{}
	_ resource.ResourceWithImportState = &FwmVpcAclRule2Resource{}
)

// NewFwmVpcAclRule2Resource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewFwmVpcAclRule2Resource() resource.Resource {
	return &FwmVpcAclRule2Resource{}
}

// FwmVpcAclRule2Resource implements resource.Resource for
// tencentcloud_fwm_vpc_acl_rule2: a single VPC border firewall ACL rule inside
// a rule group.
//
// The VPC ACL cloud APIs are collection oriented (CreateVpcAclRule takes a
// Rules array, DeleteVpcAclRule takes RuleIds). This resource always sends and
// manages exactly one rule, and reads the rule back by rule id, so the
// Terraform resource stays one-rule-per-instance instead of turning into a
// multi-object resource.
//
// The identity is `<group_id>#<rule_id>`. GroupId is the mandatory query scope
// of every VPC ACL API and has no read-back leaf of its own, which is why the
// resource id carries it as an immutable segment.
type FwmVpcAclRule2Resource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// FwmVpcAclRule2ResourceModel maps the schema attributes. Every schema
// attribute appears here exactly once.
//
// OrderIndex is kept even though DescribeVpcAclRules returns the same value
// under Rules.Sequence: Sequence is the read-back candidate of OrderIndex (name
// difference only, same meaning of priority starting from 1). Both are declared
// so a priority change on the cloud side can be detected. If Sequence ever
// comes back empty for a rule, the refresh keeps the configured order_index
// instead of clearing it.
//
// BelongMemberId (write side) and BelongMember.MemberId (read side) are the
// same value in different containers and are mapped onto each other.
type FwmVpcAclRule2ResourceModel struct {
	Id                types.String `tfsdk:"id"`
	GroupId           types.String `tfsdk:"group_id"`
	RuleId            types.String `tfsdk:"rule_id"`
	SourceContent     types.String `tfsdk:"source_content"`
	SourceType        types.String `tfsdk:"source_type"`
	DestContent       types.String `tfsdk:"dest_content"`
	DestType          types.String `tfsdk:"dest_type"`
	Protocol          types.String `tfsdk:"protocol"`
	Port              types.String `tfsdk:"port"`
	RuleAction        types.String `tfsdk:"rule_action"`
	OrderIndex        types.Int64  `tfsdk:"order_index"`
	EdgeId            types.String `tfsdk:"edge_id"`
	FwGroupId         types.String `tfsdk:"fw_group_id"`
	Description       types.String `tfsdk:"description"`
	ParamTemplateId   types.String `tfsdk:"param_template_id"`
	BelongMemberId    types.String `tfsdk:"belong_member_id"`
	IpVersion         types.String `tfsdk:"ip_version"`
	Invalid           types.Int64  `tfsdk:"invalid"`
	ParamTemplateName types.String `tfsdk:"param_template_name"`
	FwGroupName       types.String `tfsdk:"fw_group_name"`
	SourceName        types.String `tfsdk:"source_name"`
	DestName          types.String `tfsdk:"dest_name"`
	CreateTime        types.String `tfsdk:"create_time"`
	UpdateTime        types.String `tfsdk:"update_time"`
	DnsParseCnt       types.Int64  `tfsdk:"dns_parse_cnt"`
	HitCount          types.Int64  `tfsdk:"hit_count"`
	LastHitTime       types.String `tfsdk:"last_hit_time"`
}

func (r *FwmVpcAclRule2Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_fwm_vpc_acl_rule2"
}

func (r *FwmVpcAclRule2Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a FWM VPC border firewall ACL rule. " +
			"A rule group may hold many rules and the cloud API creates and deletes rules in batches, " +
			"but this resource always manages exactly one rule: the rule id is stored in the resource id " +
			"(`<group_id>#<rule_id>`) and the rule is read back by id.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier in the form `<group_id>#<rule_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the rule group the rule belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rule_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the rule, assigned by the cloud side on creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_content": schema.StringAttribute{
				Required:    true,
				Description: "Source address content.",
			},
			"source_type": schema.StringAttribute{
				Required:    true,
				Description: "Source type: ip/url/template/instance/tag. Changing it recreates the rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dest_content": schema.StringAttribute{
				Required:    true,
				Description: "Destination address content.",
			},
			"dest_type": schema.StringAttribute{
				Required:    true,
				Description: "Destination type: ip/url/template/instance/tag. Changing it recreates the rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "Protocol: TCP/UDP/ICMP/ANY/HTTP/HTTPS/DNS/FTP and so on.",
			},
			"port": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Port, for example `80`, `80,443` or `80-90`. It is empty for the ICMP protocol.",
			},
			"rule_action": schema.StringAttribute{
				Required:    true,
				Description: "Action of the rule: accept/drop/log.",
			},
			"order_index": schema.Int64Attribute{
				Required:    true,
				Description: "Priority of the rule, starting from 1.",
			},
			"edge_id": schema.StringAttribute{
				Required:    true,
				Description: "Border firewall ID: `ALL` for global, `CFWS-xxx` for a specific border. Changing it recreates the rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"fw_group_id": schema.StringAttribute{
				Required:    true,
				Description: "Firewall instance ID (effective scope of the rule): `ALL` for global, `ccn-xxx` for a Cloud Connect Network instance, `cfwg-xxx` for a firewall group instance. Changing it recreates the rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"param_template_id": schema.StringAttribute{
				Optional:    true,
				Description: "Port template ID.",
			},
			"belong_member_id": schema.StringAttribute{
				Optional:    true,
				Description: "ID of the member account owning the rule. Required when `fw_group_id` is `cfwg-xxx` or when `source_type`/`dest_type` is `instance`/`tag`. Changing it recreates the rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the rule.",
			},
			// Computed attributes owned by the cloud side.
			"ip_version": schema.StringAttribute{
				Computed:    true,
				Description: "IP version: ipv4 or ipv6.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"invalid": schema.Int64Attribute{
				Computed:    true,
				Description: "Whether the rule is invalid: 0 - valid, 1 - invalid.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"param_template_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the port template.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"fw_group_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the firewall group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the source asset.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dest_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the destination asset.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "Last modification time of the rule.",
			},
			"dns_parse_cnt": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of parsed domains of the rule.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"hit_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Cumulative hit count of the rule across all delivered member accounts, as returned by the cloud API. It is 0 when the rule has never been hit.",
			},
			"last_hit_time": schema.StringAttribute{
				Computed:    true,
				Description: "Latest hit time of the rule across all delivered member accounts. It is `-` when the rule has never been hit.",
			},
		},
	}
}

// Create handles `terraform apply` for a new rule.
func (r *FwmVpcAclRule2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FwmVpcAclRule2ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule := &fwmv20250611.VpcAclRule{
		SourceContent: plan.SourceContent.ValueStringPointer(),
		SourceType:    plan.SourceType.ValueStringPointer(),
		DestContent:   plan.DestContent.ValueStringPointer(),
		DestType:      plan.DestType.ValueStringPointer(),
		Protocol:      plan.Protocol.ValueStringPointer(),
		RuleAction:    plan.RuleAction.ValueStringPointer(),
		OrderIndex:    plan.OrderIndex.ValueInt64Pointer(),
		EdgeId:        plan.EdgeId.ValueStringPointer(),
		FwGroupId:     plan.FwGroupId.ValueStringPointer(),
	}
	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		rule.Port = plan.Port.ValueStringPointer()
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		rule.Description = plan.Description.ValueStringPointer()
	}
	if !plan.ParamTemplateId.IsNull() && !plan.ParamTemplateId.IsUnknown() {
		rule.ParamTemplateId = plan.ParamTemplateId.ValueStringPointer()
	}
	if !plan.BelongMemberId.IsNull() && !plan.BelongMemberId.IsUnknown() {
		rule.BelongMemberId = plan.BelongMemberId.ValueStringPointer()
	}

	// The API accepts a rule array, this resource creates exactly one rule.
	request := fwmv20250611.NewCreateVpcAclRuleRequest()
	request.GroupId = plan.GroupId.ValueStringPointer()
	request.Rules = []*fwmv20250611.VpcAclRule{rule}

	ruleIds, err := NewFwmService(r.Client()).CreateFwmVpcAclRule(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating FWM VPC ACL rule",
			fmt.Sprintf("CreateVpcAclRule failed for group %q: %v", plan.GroupId.ValueString(), err),
		)
		return
	}

	// A single-element request must come back with exactly one rule id. An empty
	// or ambiguous response is an API failure, never a partially created rule:
	// writing an empty id into state would leave an unmanageable resource.
	if len(ruleIds) != 1 || ruleIds[0] == nil || *ruleIds[0] == "" {
		resp.Diagnostics.AddError(
			"Error creating FWM VPC ACL rule",
			fmt.Sprintf("CreateVpcAclRule returned %d rule ids for a single-rule request; exactly one non-empty rule id is required.", len(ruleIds)),
		)
		return
	}

	plan.RuleId = types.StringValue(*ruleIds[0])
	plan.Id = types.StringValue(strings.Join([]string{plan.GroupId.ValueString(), plan.RuleId.ValueString()}, tccommon.FILED_SP))

	// Persist the identity before reading back, so a failed read-back still
	// leaves a destroyable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := NewFwmService(r.Client()).DescribeFwmVpcAclRuleById(ctx, plan.GroupId.ValueString(), plan.RuleId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading FWM VPC ACL rule after create",
			fmt.Sprintf("DescribeVpcAclRules failed for group %q rule %q: %v", plan.GroupId.ValueString(), plan.RuleId.ValueString(), err),
		)
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"FWM VPC ACL rule not readable after create",
			fmt.Sprintf("rule %q in group %q was created but cannot be read back.", plan.RuleId.ValueString(), plan.GroupId.ValueString()),
		)
		return
	}

	refreshVpcAclRule2Model(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and the post-apply state refresh.
func (r *FwmVpcAclRule2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FwmVpcAclRule2ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId, ruleId, ok := splitFwmVpcAclRule2Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("fwm vpc acl rule id is broken, id is %q (expected `<group_id>%s<rule_id>`)", state.Id.ValueString(), tccommon.FILED_SP),
		)
		return
	}

	item, err := NewFwmService(r.Client()).DescribeFwmVpcAclRuleById(ctx, groupId, ruleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading FWM VPC ACL rule",
			fmt.Sprintf("DescribeVpcAclRules failed for group %q rule %q: %v", groupId, ruleId, err),
		)
		return
	}

	// A rule that is no longer present in the group has been deleted outside of
	// Terraform: drop it from state so it gets recreated. API errors never land
	// here, so a permission or timeout error cannot silently clear the state.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("fwm vpc acl rule %q in group %q not found, removing from state.", ruleId, groupId),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	state.GroupId = types.StringValue(groupId)
	refreshVpcAclRule2Model(&state, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications of a single rule.
func (r *FwmVpcAclRule2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FwmVpcAclRule2ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state FwmVpcAclRule2ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId, ruleId, ok := splitFwmVpcAclRule2Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("fwm vpc acl rule id is broken, id is %q (expected `<group_id>%s<rule_id>`)", state.Id.ValueString(), tccommon.FILED_SP),
		)
		return
	}

	mergeVpcAclRule2Plan(state, &plan)

	rule := &fwmv20250611.VpcAclRule{
		RuleId:        helper.String(ruleId),
		SourceContent: plan.SourceContent.ValueStringPointer(),
		SourceType:    plan.SourceType.ValueStringPointer(),
		DestContent:   plan.DestContent.ValueStringPointer(),
		DestType:      plan.DestType.ValueStringPointer(),
		Protocol:      plan.Protocol.ValueStringPointer(),
		RuleAction:    plan.RuleAction.ValueStringPointer(),
		OrderIndex:    plan.OrderIndex.ValueInt64Pointer(),
		EdgeId:        plan.EdgeId.ValueStringPointer(),
		FwGroupId:     plan.FwGroupId.ValueStringPointer(),
	}
	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		rule.Port = plan.Port.ValueStringPointer()
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		rule.Description = plan.Description.ValueStringPointer()
	}
	if !plan.ParamTemplateId.IsNull() && !plan.ParamTemplateId.IsUnknown() {
		rule.ParamTemplateId = plan.ParamTemplateId.ValueStringPointer()
	}
	if !plan.BelongMemberId.IsNull() && !plan.BelongMemberId.IsUnknown() {
		rule.BelongMemberId = plan.BelongMemberId.ValueStringPointer()
	}

	request := fwmv20250611.NewModifyVpcAclRuleRequest()
	request.GroupId = helper.String(groupId)
	request.Rule = rule

	if err := NewFwmService(r.Client()).ModifyFwmVpcAclRule(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating FWM VPC ACL rule",
			fmt.Sprintf("ModifyVpcAclRule failed for group %q rule %q: %v", groupId, ruleId, err),
		)
		return
	}

	item, err := NewFwmService(r.Client()).DescribeFwmVpcAclRuleById(ctx, groupId, ruleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading FWM VPC ACL rule after update",
			fmt.Sprintf("DescribeVpcAclRules failed for group %q rule %q: %v", groupId, ruleId, err),
		)
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"FWM VPC ACL rule not readable after update",
			fmt.Sprintf("rule %q in group %q was updated but cannot be read back.", ruleId, groupId),
		)
		return
	}

	plan.Id = state.Id
	plan.RuleId = state.RuleId
	plan.GroupId = types.StringValue(groupId)
	refreshVpcAclRule2Model(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles `terraform destroy`. DeleteVpcAclRule takes a rule id list,
// only this resource's rule id is sent.
func (r *FwmVpcAclRule2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FwmVpcAclRule2ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId, ruleId, ok := splitFwmVpcAclRule2Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("fwm vpc acl rule id is broken, id is %q (expected `<group_id>%s<rule_id>`)", state.Id.ValueString(), tccommon.FILED_SP),
		)
		return
	}

	// A rule that is already gone is treated as deleted, so a repeated destroy
	// converges instead of failing.
	if err := NewFwmService(r.Client()).DeleteFwmVpcAclRule(ctx, groupId, ruleId); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting FWM VPC ACL rule",
			fmt.Sprintf("DeleteVpcAclRule failed for group %q rule %q: %v", groupId, ruleId, err),
		)
		return
	}
}

// ImportState supports `terraform import` using the `<group_id>#<rule_id>` form.
func (r *FwmVpcAclRule2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	groupId, ruleId, ok := splitFwmVpcAclRule2Id(req.ID)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("expected import ID in the form `<group_id>%s<rule_id>`, got %q", tccommon.FILED_SP, req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_id"), groupId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("rule_id"), ruleId)...)
}

// splitFwmVpcAclRule2Id parses `<group_id>#<rule_id>`.
func splitFwmVpcAclRule2Id(id string) (groupId, ruleId string, ok bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// mergeVpcAclRule2Plan fills the unset (null) optional attributes of the plan
// from the prior state.
//
// Optional attributes are declared as Optional+Computed for `port` and plain
// Optional for the rest, so a user who removes a field from the configuration
// gets Null in the plan. ModifyVpcAclRule does not document whether an omitted
// field keeps or clears the cloud value; sending nothing for such a field (the
// historical provider behaviour) keeps the old cloud value, therefore the plan
// value must be reverted to the previous state value as well, otherwise the
// refreshed state would always diverge from the configuration.
func mergeVpcAclRule2Plan(state FwmVpcAclRule2ResourceModel, plan *FwmVpcAclRule2ResourceModel) {
	if plan.Port.IsNull() || plan.Port.IsUnknown() {
		plan.Port = state.Port
	}
	if plan.Description.IsNull() || plan.Description.IsUnknown() {
		plan.Description = state.Description
	}
	if plan.ParamTemplateId.IsNull() || plan.ParamTemplateId.IsUnknown() {
		plan.ParamTemplateId = state.ParamTemplateId
	}
	if plan.BelongMemberId.IsNull() || plan.BelongMemberId.IsUnknown() {
		plan.BelongMemberId = state.BelongMemberId
	}
	if plan.IpVersion.IsNull() || plan.IpVersion.IsUnknown() {
		plan.IpVersion = state.IpVersion
	}
	if plan.Invalid.IsNull() || plan.Invalid.IsUnknown() {
		plan.Invalid = state.Invalid
	}
	if plan.ParamTemplateName.IsNull() || plan.ParamTemplateName.IsUnknown() {
		plan.ParamTemplateName = state.ParamTemplateName
	}
	if plan.FwGroupName.IsNull() || plan.FwGroupName.IsUnknown() {
		plan.FwGroupName = state.FwGroupName
	}
	if plan.SourceName.IsNull() || plan.SourceName.IsUnknown() {
		plan.SourceName = state.SourceName
	}
	if plan.DestName.IsNull() || plan.DestName.IsUnknown() {
		plan.DestName = state.DestName
	}
	if plan.CreateTime.IsNull() || plan.CreateTime.IsUnknown() {
		plan.CreateTime = state.CreateTime
	}
	if plan.UpdateTime.IsNull() || plan.UpdateTime.IsUnknown() {
		plan.UpdateTime = state.UpdateTime
	}
	if plan.DnsParseCnt.IsNull() || plan.DnsParseCnt.IsUnknown() {
		plan.DnsParseCnt = state.DnsParseCnt
	}
	if plan.HitCount.IsNull() || plan.HitCount.IsUnknown() {
		plan.HitCount = state.HitCount
	}
	if plan.LastHitTime.IsNull() || plan.LastHitTime.IsUnknown() {
		plan.LastHitTime = state.LastHitTime
	}
}

// refreshVpcAclRule2Model copies the cloud returned rule into the model.
//
// OrderIndex is refreshed from Rules.Sequence of the same rule. Only a
// non-empty Sequence overwrites the configured order_index, so an API that
// stops returning the priority cannot silently blank out the configuration.
func refreshVpcAclRule2Model(m *FwmVpcAclRule2ResourceModel, item *fwmv20250611.VpcAclRuleResp) {
	if item == nil {
		return
	}

	if item.RuleId != nil {
		m.RuleId = types.StringValue(*item.RuleId)
	}
	if item.SourceContent != nil {
		m.SourceContent = types.StringValue(*item.SourceContent)
	}
	if item.SourceType != nil {
		m.SourceType = types.StringValue(*item.SourceType)
	}
	if item.DestContent != nil {
		m.DestContent = types.StringValue(*item.DestContent)
	}
	if item.DestType != nil {
		m.DestType = types.StringValue(*item.DestType)
	}
	if item.Protocol != nil {
		m.Protocol = types.StringValue(*item.Protocol)
	}
	m.Port = stringPtrValueOrKeep(item.Port, m.Port)
	m.ParamTemplateId = stringPtrValueOrKeep(item.ParamTemplateId, m.ParamTemplateId)
	m.Description = stringPtrValueOrKeep(item.Description, m.Description)
	if item.RuleAction != nil {
		m.RuleAction = types.StringValue(*item.RuleAction)
	}
	if item.Sequence != nil && *item.Sequence >= 0 {
		m.OrderIndex = types.Int64Value(*item.Sequence)
	}
	if item.EdgeId != nil {
		m.EdgeId = types.StringValue(*item.EdgeId)
	}
	if item.FwGroupId != nil {
		m.FwGroupId = types.StringValue(*item.FwGroupId)
	}
	m.FwGroupName = stringPtrValueOrKeep(item.FwGroupName, m.FwGroupName)
	m.SourceName = stringPtrValueOrKeep(item.SourceName, m.SourceName)
	m.DestName = stringPtrValueOrKeep(item.DestName, m.DestName)
	if item.BelongMember != nil && item.BelongMember.MemberId != nil && *item.BelongMember.MemberId != "" {
		m.BelongMemberId = types.StringValue(*item.BelongMember.MemberId)
	}
	if item.IpVersion != nil {
		m.IpVersion = types.StringValue(*item.IpVersion)
	}
	if item.Invalid != nil {
		m.Invalid = types.Int64Value(*item.Invalid)
	}
	m.ParamTemplateName = stringPtrValueOrKeep(item.ParamTemplateName, m.ParamTemplateName)
	m.CreateTime = stringPtrValueOrKeep(item.CreateTime, m.CreateTime)
	m.UpdateTime = stringPtrValueOrKeep(item.UpdateTime, m.UpdateTime)
	m.LastHitTime = stringPtrValueOrKeep(item.LastHitTime, m.LastHitTime)
	if item.DnsParseCnt != nil {
		m.DnsParseCnt = types.Int64Value(*item.DnsParseCnt)
	}
	if item.HitCount != nil {
		m.HitCount = types.Int64Value(*item.HitCount)
	}
}

// stringPtrValueOrKeep returns the cloud value when it is present, otherwise it
// keeps the current model value (Null stays Null, a configured value is kept).
func stringPtrValueOrKeep(v *string, current types.String) types.String {
	if v == nil {
		return current
	}
	return types.StringValue(*v)
}