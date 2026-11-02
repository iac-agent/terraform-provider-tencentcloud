package fwm

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// _ resource.Resource is intentionally omitted: the factory function
// NewFwmVpcAclRule2Resource returns resource.Resource, which already performs
// the implicit compile-time check at the return statement. Only the
// extended-interface assertions are kept here.
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
// tencentcloud_fwm_vpc_acl_rule2.
//
// One Terraform resource maps onto ONE VPC border firewall ACL rule. The cloud
// API expresses rules as a list (CreateVpcAclRule.Rules / DeleteVpcAclRule.RuleIds),
// so this resource always sends a single-element list on Create and a single
// id on Delete, and validates that exactly one rule comes back on Read.
//
// Deleting this resource only removes the referenced rule from the rule group
// (DeleteVpcAclRule with RuleIds=[rule_id]); the rule group itself and every
// other rule in the group are left untouched.
type FwmVpcAclRule2Resource struct {
	fw.ResourceWithConfigure
}

// FwmVpcAclRule2ResourceModel maps the schema attributes.
//
// group_id and rule_id together form the composite identity
// `<group_id>#<rule_id>`; rule_id is known after Create and is therefore
// Computed.
type FwmVpcAclRule2ResourceModel struct {
	Id              types.String `tfsdk:"id"`
	GroupId         types.String `tfsdk:"group_id"`
	RuleId          types.String `tfsdk:"rule_id"`
	SourceContent   types.String `tfsdk:"source_content"`
	SourceType      types.String `tfsdk:"source_type"`
	DestContent     types.String `tfsdk:"dest_content"`
	DestType        types.String `tfsdk:"dest_type"`
	Protocol        types.String `tfsdk:"protocol"`
	Port            types.String `tfsdk:"port"`
	RuleAction      types.String `tfsdk:"rule_action"`
	OrderIndex      types.Int64  `tfsdk:"order_index"`
	EdgeId          types.String `tfsdk:"edge_id"`
	FwGroupId       types.String `tfsdk:"fw_group_id"`
	Description     types.String `tfsdk:"description"`
	ParamTemplateId types.String `tfsdk:"param_template_id"`
	BelongMemberId  types.String `tfsdk:"belong_member_id"`

	Sequence          types.Int64  `tfsdk:"sequence"`
	IpVersion         types.String `tfsdk:"ip_version"`
	SourceName        types.String `tfsdk:"source_name"`
	DestName          types.String `tfsdk:"dest_name"`
	ParamTemplateName types.String `tfsdk:"param_template_name"`
	Invalid           types.Int64  `tfsdk:"invalid"`
	AppId             types.String `tfsdk:"app_id"`
	Uin               types.String `tfsdk:"uin"`
	Nickname          types.String `tfsdk:"nickname"`
	NodeName          types.String `tfsdk:"node_name"`
	CreateTime        types.String `tfsdk:"create_time"`
	UpdateTime        types.String `tfsdk:"update_time"`
	DnsParseCnt       types.Int64  `tfsdk:"dns_parse_cnt"`
	FwGroupName       types.String `tfsdk:"fw_group_name"`
	HitCount          types.Int64  `tfsdk:"hit_count"`
	LastHitTime       types.String `tfsdk:"last_hit_time"`
}

func (r *FwmVpcAclRule2Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_fwm_vpc_acl_rule2"
}

func (r *FwmVpcAclRule2Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a VPC border firewall ACL rule of Cloud Firewall (FWM). " +
			"One resource manages exactly one ACL rule inside the rule group given by group_id.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier in the form `<group_id>#<rule_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.StringAttribute{
				Required: true,
				Description: "ID of the ACL rule group that owns the rule, such as `fwmrg-9z03vakjj8`. " +
					"The rule group is the parent container and must already exist before creating the rule. " +
					"Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rule_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the ACL rule, returned by the API on create and used to locate the rule afterwards.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_content": schema.StringAttribute{
				Required:    true,
				Description: "Source address content. Its format depends on source_type: CIDR for ip, domain name for url, template ID for template, instance ID list for instance, tag for tag.",
			},
			"source_type": schema.StringAttribute{
				Required:    true,
				Description: "Source address type: ip, url, template, instance, tag.",
			},
			"dest_content": schema.StringAttribute{
				Required:    true,
				Description: "Destination address content. Its format depends on dest_type and follows the same rules as source_content.",
			},
			"dest_type": schema.StringAttribute{
				Required:    true,
				Description: "Destination address type: ip, url, template, instance, tag.",
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "Protocol, such as TCP, UDP, ICMP, ANY, HTTP, HTTPS, DNS, FTP.",
			},
			"port": schema.StringAttribute{
				Optional:    true,
				Description: "Port or port range. Empty for the ICMP protocol.",
			},
			"rule_action": schema.StringAttribute{
				Required:    true,
				Description: "Action of the rule: accept, drop, log.",
			},
			"order_index": schema.Int64Attribute{
				Required: true,
				Description: "Priority of the rule, starting from 1. The rule is placed at this position inside the rule group and " +
					"the other rules of the group may be re-sequenced by the cloud side as a result. " +
					"The effective sequence is exposed by the computed `sequence` attribute.",
			},
			"edge_id": schema.StringAttribute{
				Required:    true,
				Description: "Border firewall ID: ALL means global, CFWS-xxx means a specific border.",
			},
			"fw_group_id": schema.StringAttribute{
				Required:    true,
				Description: "Firewall instance ID that the rule applies to: ALL means globally effective, ccn-xxx a Cloud Connect Network instance, cfwg-xxx a firewall group instance.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the rule.",
			},
			"param_template_id": schema.StringAttribute{
				Optional:    true,
				Description: "ID of the port template used by the rule.",
			},
			"belong_member_id": schema.StringAttribute{
				Optional: true,
				Description: "ID of the member account that the rule belongs to. Required when fw_group_id is a cfwg-xxx firewall group " +
					"or when source_type/dest_type is instance or tag.",
			},

			// Computed, read-only values returned by DescribeVpcAclRules.
			"sequence": schema.Int64Attribute{
				Computed:    true,
				Description: "Effective priority of the rule inside the rule group, as returned by the cloud side.",
			},
			"ip_version": schema.StringAttribute{
				Computed:    true,
				Description: "IP version of the rule: ipv4 or ipv6.",
			},
			"source_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the source asset, when the source is an instance or a tag.",
			},
			"dest_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the destination asset, when the destination is an instance or a tag.",
			},
			"param_template_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the port template used by the rule.",
			},
			"invalid": schema.Int64Attribute{
				Computed:    true,
				Description: "Whether the rule is invalid: 0 means valid, 1 means invalid.",
			},
			"app_id": schema.StringAttribute{
				Computed:    true,
				Description: "AppId of the member account that the rule belongs to.",
			},
			"uin": schema.StringAttribute{
				Computed:    true,
				Description: "Uin of the member account that the rule belongs to.",
			},
			"nickname": schema.StringAttribute{
				Computed:    true,
				Description: "Nickname of the member account that the rule belongs to.",
			},
			"node_name": schema.StringAttribute{
				Computed:    true,
				Description: "Department that the member account belongs to.",
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
				Description: "Number of domains resolved by the rule.",
			},
			"fw_group_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the firewall group that the rule applies to.",
			},
			"hit_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Accumulated hit count of the rule over all member accounts it was delivered to, as returned by the cloud side. Zero when the rule has never been hit.",
			},
			"last_hit_time": schema.StringAttribute{
				Computed:    true,
				Description: "Latest hit time of the rule over all member accounts it was delivered to. \"-\" when the rule has never been hit.",
			},
		},
	}
}

// Create handles `terraform apply` for a new ACL rule.
func (r *FwmVpcAclRule2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FwmVpcAclRule2ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId := plan.GroupId.ValueString()
	rule := buildVpcAclRule2(plan)

	service := NewFwmService(r.Client())
	ruleId, err := service.CreateFwmVpcAclRule(ctx, groupId, rule)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating FWM VPC ACL rule",
			fmt.Sprintf("CreateVpcAclRule failed for group %q: %v", groupId, err),
		)
		return
	}

	plan.Id = types.StringValue(strings.Join([]string{groupId, ruleId}, tccommon.FILED_SP))
	plan.RuleId = types.StringValue(ruleId)

	// Persist the identity BEFORE the read-back: if the read-back fails, the
	// resource is still in state and `terraform destroy` can clean it up.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := service.DescribeFwmVpcAclRuleById(ctx, groupId, ruleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading FWM VPC ACL rule after create",
			fmt.Sprintf("DescribeVpcAclRules failed for group %q rule %q: %v", groupId, ruleId, err),
		)
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"FWM VPC ACL rule not readable after create",
			fmt.Sprintf("Rule %q in group %q was created but cannot be read back.", ruleId, groupId),
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
			fmt.Sprintf("FWM VPC ACL rule id is broken, id is %q (expected `<group_id>%s<rule_id>`)", state.Id.ValueString(), tccommon.FILED_SP),
		)
		return
	}

	service := NewFwmService(r.Client())
	item, err := service.DescribeFwmVpcAclRuleById(ctx, groupId, ruleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading FWM VPC ACL rule",
			fmt.Sprintf("DescribeVpcAclRules failed for group %q rule %q: %v", groupId, ruleId, err),
		)
		return
	}

	// Only a successful query with no matching rule means "deleted on the cloud
	// side"; a nil response is an API failure and is reported above.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("FWM VPC ACL rule %q in group %q not found, removing from state.", ruleId, groupId),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	// group_id and rule_id are not echoed back by DescribeVpcAclRules, so they
	// are re-derived from the composite id instead of being read from the item.
	state.GroupId = types.StringValue(groupId)
	state.RuleId = types.StringValue(ruleId)
	refreshVpcAclRule2Model(&state, item)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications of an existing ACL rule.
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
			fmt.Sprintf("FWM VPC ACL rule id is broken, id is %q (expected `<group_id>%s<rule_id>`)", state.Id.ValueString(), tccommon.FILED_SP),
		)
		return
	}

	// ModifyVpcAclRule replaces the whole rule record, so the request must carry
	// the complete rule. Optional fields the user has removed from the
	// configuration are sent as the prior state value: the clear-vs-keep
	// semantics of an omitted optional field is not documented for this API, and
	// guessing "clear" would silently wipe the cloud-side value.
	plan.RuleId = types.StringValue(ruleId)
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

	rule := buildVpcAclRule2(plan)

	service := NewFwmService(r.Client())
	if err := service.ModifyFwmVpcAclRule(ctx, groupId, rule); err != nil {
		resp.Diagnostics.AddError(
			"Error updating FWM VPC ACL rule",
			fmt.Sprintf("ModifyVpcAclRule failed for group %q rule %q: %v", groupId, ruleId, err),
		)
		return
	}

	item, err := service.DescribeFwmVpcAclRuleById(ctx, groupId, ruleId)
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
			fmt.Sprintf("Rule %q in group %q was updated but cannot be read back.", ruleId, groupId),
		)
		return
	}

	plan.Id = state.Id
	plan.GroupId = types.StringValue(groupId)
	plan.RuleId = types.StringValue(ruleId)
	refreshVpcAclRule2Model(&plan, item)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles `terraform destroy` by removing only the referenced rule.
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
			fmt.Sprintf("FWM VPC ACL rule id is broken, id is %q (expected `<group_id>%s<rule_id>`)", state.Id.ValueString(), tccommon.FILED_SP),
		)
		return
	}

	service := NewFwmService(r.Client())
	if err := service.DeleteFwmVpcAclRule(ctx, groupId, ruleId); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting FWM VPC ACL rule",
			fmt.Sprintf("DeleteVpcAclRule failed for group %q rule %q: %v", groupId, ruleId, err),
		)
		return
	}
	// The Framework removes the entry from state after a successful Delete.
}

// ImportState supports `terraform import` with the composite id
// `<group_id>#<rule_id>` and fills both identity segments in state.
func (r *FwmVpcAclRule2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	groupId, ruleId, ok := splitFwmVpcAclRule2Id(req.ID)
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

// buildVpcAclRule2 converts the model into the SDK rule struct used by both
// CreateVpcAclRule and ModifyVpcAclRule.
func buildVpcAclRule2(m FwmVpcAclRule2ResourceModel) *fwmv20250611.VpcAclRule {
	rule := &fwmv20250611.VpcAclRule{
		SourceContent:   stringPtrOrNil(m.SourceContent),
		SourceType:      stringPtrOrNil(m.SourceType),
		DestContent:     stringPtrOrNil(m.DestContent),
		DestType:        stringPtrOrNil(m.DestType),
		Protocol:        stringPtrOrNil(m.Protocol),
		RuleAction:      stringPtrOrNil(m.RuleAction),
		EdgeId:          stringPtrOrNil(m.EdgeId),
		FwGroupId:       stringPtrOrNil(m.FwGroupId),
		RuleId:          stringPtrOrNil(m.RuleId),
		Port:            stringPtrOrNil(m.Port),
		Description:     stringPtrOrNil(m.Description),
		ParamTemplateId: stringPtrOrNil(m.ParamTemplateId),
		BelongMemberId:  stringPtrOrNil(m.BelongMemberId),
	}

	if !m.OrderIndex.IsNull() && !m.OrderIndex.IsUnknown() {
		rule.OrderIndex = helper.IntInt64(int(m.OrderIndex.ValueInt64()))
	}

	return rule
}

// refreshVpcAclRule2Model copies the cloud-owned values of a described rule
// into the model.
//
// group_id has no counterpart in the response and rule_id is part of the
// identity, so only the remaining business fields are refreshed here.
func refreshVpcAclRule2Model(m *FwmVpcAclRule2ResourceModel, item *fwmv20250611.VpcAclRuleResp) {
	if item == nil {
		return
	}

	m.SourceContent = types.StringPointerValue(item.SourceContent)
	m.SourceType = types.StringPointerValue(item.SourceType)
	m.DestContent = types.StringPointerValue(item.DestContent)
	m.DestType = types.StringPointerValue(item.DestType)
	m.Protocol = types.StringPointerValue(item.Protocol)
	m.RuleAction = types.StringPointerValue(item.RuleAction)
	m.EdgeId = types.StringPointerValue(item.EdgeId)
	m.FwGroupId = types.StringPointerValue(item.FwGroupId)

	// Optional inputs: the API may omit them in its response, in which case the
	// previously known value is kept instead of being nulled out (that would
	// produce a permanent diff against the configuration).
	m.Port = stringPointerValueOrKeep(item.Port, m.Port)
	m.Description = stringPointerValueOrKeep(item.Description, m.Description)
	m.ParamTemplateId = stringPointerValueOrKeep(item.ParamTemplateId, m.ParamTemplateId)
	if item.BelongMember != nil {
		m.BelongMemberId = stringPointerValueOrKeep(item.BelongMember.MemberId, m.BelongMemberId)
	} else {
		m.BelongMemberId = keepStringValue(m.BelongMemberId)
	}

	// order_index is the writable priority. The API only documents the
	// response-side equivalent as `Sequence`, but every response in practice
	// mirrors the requested order index, so the sequence is used as the
	// read-back value (documented as such in the resource documentation).
	if item.Sequence != nil {
		m.OrderIndex = types.Int64PointerValue(item.Sequence)
	} else {
		m.OrderIndex = keepInt64Value(m.OrderIndex)
	}

	m.Sequence = types.Int64PointerValue(item.Sequence)
	m.IpVersion = types.StringPointerValue(item.IpVersion)
	m.SourceName = types.StringPointerValue(item.SourceName)
	m.DestName = types.StringPointerValue(item.DestName)
	m.ParamTemplateName = types.StringPointerValue(item.ParamTemplateName)
	m.Invalid = types.Int64PointerValue(item.Invalid)
	m.CreateTime = types.StringPointerValue(item.CreateTime)
	m.UpdateTime = types.StringPointerValue(item.UpdateTime)
	m.DnsParseCnt = types.Int64PointerValue(item.DnsParseCnt)
	m.FwGroupName = types.StringPointerValue(item.FwGroupName)
	m.HitCount = types.Int64PointerValue(item.HitCount)
	m.LastHitTime = types.StringPointerValue(item.LastHitTime)

	if item.BelongMember != nil {
		m.AppId = types.StringPointerValue(item.BelongMember.AppId)
		m.Uin = types.StringPointerValue(item.BelongMember.Uin)
		m.Nickname = types.StringPointerValue(item.BelongMember.Nickname)
		m.NodeName = types.StringPointerValue(item.BelongMember.NodeName)
	} else {
		m.AppId = types.StringNull()
		m.Uin = types.StringNull()
		m.Nickname = types.StringNull()
		m.NodeName = types.StringNull()
	}
}

// splitFwmVpcAclRule2Id parses the composite id `<group_id>#<rule_id>`.
func splitFwmVpcAclRule2Id(id string) (groupId, ruleId string, ok bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// stringPointerValueOrKeep maps a response pointer onto types.String, keeping
// the previous value when the cloud side returned nothing.
func stringPointerValueOrKeep(v *string, prev types.String) types.String {
	if v == nil {
		return keepStringValue(prev)
	}
	return types.StringValue(*v)
}

// keepStringValue returns the previous value unless it was null or unknown,
// in which case it falls back to null.
func keepStringValue(prev types.String) types.String {
	if prev.IsNull() || prev.IsUnknown() {
		return types.StringNull()
	}
	return prev
}

// keepInt64Value returns the previous value unless it was null or unknown.
func keepInt64Value(prev types.Int64) types.Int64 {
	if prev.IsNull() || prev.IsUnknown() {
		return types.Int64Null()
	}
	return prev
}

// stringPtrOrNil converts an optional attribute into an SDK string pointer.
func stringPtrOrNil(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	return helper.String(v.ValueString())
}
