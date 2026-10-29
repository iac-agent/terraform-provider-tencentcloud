package fwm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var (
	_ resource.ResourceWithConfigure   = &FwmSecurityGroupRuleResource{}
	_ resource.ResourceWithImportState = &FwmSecurityGroupRuleResource{}
)

// NewFwmSecurityGroupRuleResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewFwmSecurityGroupRuleResource() resource.Resource {
	return &FwmSecurityGroupRuleResource{}
}

// FwmSecurityGroupRuleResource implements resource.Resource for
// tencentcloud_fwm_security_group_rule.
type FwmSecurityGroupRuleResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// FwmSecurityGroupRuleResourceModel maps the schema attributes of
// tencentcloud_fwm_security_group_rule.
type FwmSecurityGroupRuleResourceModel struct {
	Id      types.String                     `tfsdk:"id"`
	GroupId types.String                     `tfsdk:"group_id"`
	Rules   []FwmSecurityGroupRuleRulesModel `tfsdk:"rules"`
	Rule    []FwmSecurityGroupRuleRuleModel  `tfsdk:"rule"`
}

// FwmSecurityGroupRuleRulesModel maps one element of the `rules` argument used
// by CreateSecurityGroupRule / ModifySecurityGroupRule.
type FwmSecurityGroupRuleRulesModel struct {
	IpVersion         types.String `tfsdk:"ip_version"`
	SourceContent     types.String `tfsdk:"source_content"`
	SourceType        types.String `tfsdk:"source_type"`
	DestContent       types.String `tfsdk:"dest_content"`
	DestType          types.String `tfsdk:"dest_type"`
	Protocol          types.String `tfsdk:"protocol"`
	Port              types.String `tfsdk:"port"`
	ServiceTemplateId types.String `tfsdk:"service_template_id"`
	RuleAction        types.String `tfsdk:"rule_action"`
	Description       types.String `tfsdk:"description"`
	OrderIndex        types.Int64  `tfsdk:"order_index"`
	RuleId            types.String `tfsdk:"rule_id"`
	Scope             types.String `tfsdk:"scope"`
	ProtocolPortType  types.Int64  `tfsdk:"protocol_port_type"`
	BelongMemberId    types.String `tfsdk:"belong_member_id"`
}

// FwmSecurityGroupRuleRuleModel maps one element of the computed `rule`
// attribute returned by DescribeSecurityGroupRule.
type FwmSecurityGroupRuleRuleModel struct {
	RuleId            types.String `tfsdk:"rule_id"`
	Sequence          types.Int64  `tfsdk:"sequence"`
	Region            types.String `tfsdk:"region"`
	IpVersion         types.String `tfsdk:"ip_version"`
	SrcContent        types.String `tfsdk:"src_content"`
	SrcType           types.String `tfsdk:"src_type"`
	DstContent        types.String `tfsdk:"dst_content"`
	DstType           types.String `tfsdk:"dst_type"`
	Protocol          types.String `tfsdk:"protocol"`
	ProtocolPortType  types.Int64  `tfsdk:"protocol_port_type"`
	ServiceTemplateId types.String `tfsdk:"service_template_id"`
	DstPort           types.String `tfsdk:"dst_port"`
	RuleAction        types.String `tfsdk:"rule_action"`
	Detail            types.String `tfsdk:"detail"`
	RuleSource        types.String `tfsdk:"rule_source"`
	Scope             types.String `tfsdk:"scope"`
}

func (r *FwmSecurityGroupRuleResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_fwm_security_group_rule"
}

func (r *FwmSecurityGroupRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage FWM enterprise security group rules. " +
			"Each apply creates one rule per element of `rules`, and destroy deletes every rule tracked by this resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier in the form `<group_id>#<rule_id>#<rule_id>...`.",
				PlanModifiers: []planmodifier.String{
					fwmStringUseStateForUnknown{},
				},
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "ID of the enterprise security group rule group.",
				PlanModifiers: []planmodifier.String{
					fwmStringUseStateForUnknown{},
				},
			},
			"rules": schema.ListNestedAttribute{
				Optional:    true,
				Description: "List of security group rules to create.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip_version": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "IP version of the rule.",
						},
						"source_content": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Source address of the rule.",
						},
						"source_type": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Source type of the rule.",
						},
						"dest_content": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Destination address of the rule.",
						},
						"dest_type": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Destination type of the rule.",
						},
						"protocol": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Protocol of the rule.",
						},
						"port": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Port range of the rule.",
						},
						"service_template_id": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Protocol port template ID of the rule.",
						},
						"rule_action": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Action of the rule, for example `1` for block and `2` for allow.",
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Description of the rule.",
						},
						"order_index": schema.Int64Attribute{
							Optional:    true,
							Computed:    true,
							Description: "Priority of the rule.",
						},
						"rule_id": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "ID of the rule.",
						},
						"scope": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Effective scope of the rule.",
						},
						"protocol_port_type": schema.Int64Attribute{
							Optional:    true,
							Computed:    true,
							Description: "Protocol port type of the rule.",
						},
						"belong_member_id": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "ID of the member account the rule belongs to. It is required when group_id is a cfwg-xxx group, or when source_type/dest_type is `instance` or `tag`.",
						},
					},
				},
			},
			"rule": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Rule details returned by the cloud API.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rule_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the rule.",
						},
						"sequence": schema.Int64Attribute{
							Computed:    true,
							Description: "Priority of the rule.",
						},
						"region": schema.StringAttribute{
							Computed:    true,
							Description: "Region of the rule.",
						},
						"ip_version": schema.StringAttribute{
							Computed:    true,
							Description: "IP version of the rule.",
						},
						"src_content": schema.StringAttribute{
							Computed:    true,
							Description: "Source content of the rule.",
						},
						"src_type": schema.StringAttribute{
							Computed:    true,
							Description: "Source type of the rule.",
						},
						"dst_content": schema.StringAttribute{
							Computed:    true,
							Description: "Destination content of the rule.",
						},
						"dst_type": schema.StringAttribute{
							Computed:    true,
							Description: "Destination type of the rule.",
						},
						"protocol": schema.StringAttribute{
							Computed:    true,
							Description: "Protocol of the rule.",
						},
						"protocol_port_type": schema.Int64Attribute{
							Computed:    true,
							Description: "Protocol port template parameter of the rule.",
						},
						"service_template_id": schema.StringAttribute{
							Computed:    true,
							Description: "Protocol port template ID of the rule.",
						},
						"dst_port": schema.StringAttribute{
							Computed:    true,
							Description: "Destination port of the rule. Single port, multiple ports and port ranges are supported.",
						},
						"rule_action": schema.StringAttribute{
							Computed:    true,
							Description: "Action of the rule, `1` for block and `2` for allow.",
						},
						"detail": schema.StringAttribute{
							Computed:    true,
							Description: "Description of the rule.",
						},
						"rule_source": schema.StringAttribute{
							Computed:    true,
							Description: "Source of the rule, `0` means it was added from the user console.",
						},
						"scope": schema.StringAttribute{
							Computed:    true,
							Description: "Effective scope of the rule, bitwise, `1` for SG (enterprise security group) and `2` for LH (lightweight server).",
						},
					},
				},
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *FwmSecurityGroupRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FwmSecurityGroupRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.GroupId.IsNull() || plan.GroupId.IsUnknown() || plan.GroupId.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing required argument",
			"group_id must be specified when creating tencentcloud_fwm_security_group_rule.",
		)
		return
	}

	if len(plan.Rules) == 0 {
		resp.Diagnostics.AddError(
			"Missing required argument",
			"rules must contain at least one element when creating tencentcloud_fwm_security_group_rule.",
		)
		return
	}

	groupId := plan.GroupId.ValueString()

	request := fwmv20250611.NewCreateSecurityGroupRuleRequest()
	request.GroupId = &groupId
	request.Rules = expandFwmSecurityGroupRuleRules(plan.Rules)

	service := NewFwmService(r.Client())
	response, err := service.CreateFwmSecurityGroupRule(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating FWM security group rule",
			fmt.Sprintf("CreateSecurityGroupRule failed for group %s: %v", groupId, err),
		)
		return
	}

	ruleIds := response.Response.RuleIds
	if len(ruleIds) == 0 {
		resp.Diagnostics.AddError(
			"Error creating FWM security group rule",
			fmt.Sprintf("CreateSecurityGroupRule returned empty RuleIds for group %s.", groupId),
		)
		return
	}

	idParts := []string{groupId}
	for _, ruleId := range ruleIds {
		if ruleId == nil || *ruleId == "" {
			resp.Diagnostics.AddError(
				"Error creating FWM security group rule",
				fmt.Sprintf("CreateSecurityGroupRule returned an empty rule ID for group %s.", groupId),
			)
			return
		}
		idParts = append(idParts, *ruleId)
	}

	// Write the ID before any follow-up API call so that the created rules can
	// always be destroyed even if a later step fails.
	plan.Id = types.StringValue(strings.Join(idParts, tccommon.FILED_SP))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Back-fill the rule IDs of the `rules` elements so that config and state agree.
	for i := range plan.Rules {
		if i+1 < len(idParts) {
			plan.Rules[i].RuleId = types.StringValue(idParts[i+1])
		}
	}

	r.refreshPlan(ctx, &plan, idParts[1:])

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *FwmSecurityGroupRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FwmSecurityGroupRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId, ruleIds, ok := splitFwmSecurityGroupRuleId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("tencentcloud_fwm_security_group_rule id is broken, id is %q (expected `<group_id>%s<rule_id>...`)", state.Id.ValueString(), tccommon.FILED_SP),
		)
		return
	}

	r.refreshPlan(ctx, &state, ruleIds)

	if len(state.Rule) == 0 {
		log.Printf("[WARN] tencentcloud_fwm_security_group_rule id=%s not found, removing from state.", state.Id.ValueString())
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("tencentcloud_fwm_security_group_rule %s not found, removing from state.", state.Id.ValueString()),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	state.GroupId = types.StringValue(groupId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications.
func (r *FwmSecurityGroupRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FwmSecurityGroupRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state FwmSecurityGroupRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId, ruleIds, ok := splitFwmSecurityGroupRuleId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("tencentcloud_fwm_security_group_rule id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	if !plan.GroupId.IsNull() && !plan.GroupId.IsUnknown() && plan.GroupId.ValueString() != groupId {
		resp.Diagnostics.AddError(
			"Unsupported change of group_id",
			fmt.Sprintf("group_id cannot be changed from %s to %s, please recreate tencentcloud_fwm_security_group_rule instead.", groupId, plan.GroupId.ValueString()),
		)
		return
	}

	// Match the planned rules with the existing rule IDs by position so that
	// every rule keeps its identity across updates.
	existingIds := make(map[string]string, len(state.Rules))
	for i, rule := range state.Rules {
		if !rule.RuleId.IsNull() && !rule.RuleId.IsUnknown() && i < len(ruleIds) {
			existingIds[rule.RuleId.ValueString()] = ruleIds[i]
		}
	}

	service := NewFwmService(r.Client())
	newRuleIds := make([]string, 0, len(plan.Rules))

	for i := range plan.Rules {
		var ruleId string
		if i < len(ruleIds) {
			ruleId = ruleIds[i]
		}
		if !plan.Rules[i].RuleId.IsNull() && !plan.Rules[i].RuleId.IsUnknown() {
			if id, exists := existingIds[plan.Rules[i].RuleId.ValueString()]; exists {
				ruleId = id
			}
		}
		plan.Rules[i].RuleId = types.StringValue(ruleId)

		request := fwmv20250611.NewModifySecurityGroupRuleRequest()
		request.GroupId = &groupId
		request.Rule = expandFwmSecurityGroupRuleRule(&plan.Rules[i])

		if err := service.ModifyFwmSecurityGroupRule(ctx, request); err != nil {
			resp.Diagnostics.AddError(
				"Error updating FWM security group rule",
				fmt.Sprintf("ModifySecurityGroupRule failed for group %s rule %s: %v", groupId, ruleId, err),
			)
			return
		}

		newRuleIds = append(newRuleIds, ruleId)
	}

	plan.Id = types.StringValue(strings.Join(append([]string{groupId}, newRuleIds...), tccommon.FILED_SP))
	plan.GroupId = types.StringValue(groupId)

	r.refreshPlan(ctx, &plan, newRuleIds)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles `terraform destroy`.
func (r *FwmSecurityGroupRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FwmSecurityGroupRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId, ruleIds, ok := splitFwmSecurityGroupRuleId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("tencentcloud_fwm_security_group_rule id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	ruleIdPointers := make([]*string, 0, len(ruleIds))
	for _, ruleId := range ruleIds {
		id := ruleId
		ruleIdPointers = append(ruleIdPointers, &id)
	}

	service := NewFwmService(r.Client())
	if err := service.DeleteFwmSecurityGroupRule(ctx, groupId, ruleIdPointers); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting FWM security group rule",
			fmt.Sprintf("DeleteSecurityGroupRule failed for group %s rules %v: %v", groupId, ruleIds, err),
		)
		return
	}
}

// ImportState supports `terraform import` using the `<group_id>#<rule_id>...` form.
func (r *FwmSecurityGroupRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _, _, ok := splitFwmSecurityGroupRuleId(req.ID); !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in the form `<group_id>%s<rule_id>...`, got %q", tccommon.FILED_SP, req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// refreshPlan queries every rule tracked by this resource and fills the
// computed `rule` attribute. It also back-fills `rules` elements when the
// configured element carries the same rule ID.
func (r *FwmSecurityGroupRuleResource) refreshPlan(ctx context.Context, m *FwmSecurityGroupRuleResourceModel, ruleIds []string) {
	groupId := m.GroupId.ValueString()
	if groupId == "" {
		parsedGroupId, _, ok := splitFwmSecurityGroupRuleId(m.Id.ValueString())
		if ok {
			groupId = parsedGroupId
		}
	}

	service := NewFwmService(r.Client())
	details := make([]FwmSecurityGroupRuleRuleModel, 0, len(ruleIds))
	for _, ruleId := range ruleIds {
		rule, err := service.DescribeFwmSecurityGroupRuleById(ctx, groupId, ruleId)
		if err != nil {
			log.Printf("[WARN] describe tencentcloud_fwm_security_group_rule failed, group_id=%s rule_id=%s, error: %v", groupId, ruleId, err)
			continue
		}
		if rule == nil {
			log.Printf("[WARN] tencentcloud_fwm_security_group_rule not found, group_id=%s rule_id=%s", groupId, ruleId)
			continue
		}
		details = append(details, flattenFwmSecurityGroupRuleRule(rule))
	}

	if len(details) > 0 {
		m.Rule = details
	} else if len(m.Rule) == 0 {
		m.Rule = nil
	}
}

// expandFwmSecurityGroupRuleRules converts the `rules` list of the model into
// the SDK type consumed by CreateSecurityGroupRule.
func expandFwmSecurityGroupRuleRules(rules []FwmSecurityGroupRuleRulesModel) []*fwmv20250611.SecurityGroupRule {
	result := make([]*fwmv20250611.SecurityGroupRule, 0, len(rules))
	for i := range rules {
		result = append(result, expandFwmSecurityGroupRuleRule(&rules[i]))
	}
	return result
}

// expandFwmSecurityGroupRuleRule converts one `rules` element into the SDK type.
func expandFwmSecurityGroupRuleRule(rule *FwmSecurityGroupRuleRulesModel) *fwmv20250611.SecurityGroupRule {
	item := &fwmv20250611.SecurityGroupRule{}
	if !rule.IpVersion.IsNull() && !rule.IpVersion.IsUnknown() {
		item.IpVersion = rule.IpVersion.ValueStringPointer()
	}
	if !rule.SourceContent.IsNull() && !rule.SourceContent.IsUnknown() {
		item.SourceContent = rule.SourceContent.ValueStringPointer()
	}
	if !rule.SourceType.IsNull() && !rule.SourceType.IsUnknown() {
		item.SourceType = rule.SourceType.ValueStringPointer()
	}
	if !rule.DestContent.IsNull() && !rule.DestContent.IsUnknown() {
		item.DestContent = rule.DestContent.ValueStringPointer()
	}
	if !rule.DestType.IsNull() && !rule.DestType.IsUnknown() {
		item.DestType = rule.DestType.ValueStringPointer()
	}
	if !rule.Protocol.IsNull() && !rule.Protocol.IsUnknown() {
		item.Protocol = rule.Protocol.ValueStringPointer()
	}
	if !rule.Port.IsNull() && !rule.Port.IsUnknown() {
		item.Port = rule.Port.ValueStringPointer()
	}
	if !rule.ServiceTemplateId.IsNull() && !rule.ServiceTemplateId.IsUnknown() {
		item.ServiceTemplateId = rule.ServiceTemplateId.ValueStringPointer()
	}
	if !rule.RuleAction.IsNull() && !rule.RuleAction.IsUnknown() {
		item.RuleAction = rule.RuleAction.ValueStringPointer()
	}
	if !rule.Description.IsNull() && !rule.Description.IsUnknown() {
		item.Description = rule.Description.ValueStringPointer()
	}
	if !rule.OrderIndex.IsNull() && !rule.OrderIndex.IsUnknown() {
		item.OrderIndex = rule.OrderIndex.ValueInt64Pointer()
	}
	if !rule.RuleId.IsNull() && !rule.RuleId.IsUnknown() {
		item.RuleId = rule.RuleId.ValueStringPointer()
	}
	if !rule.Scope.IsNull() && !rule.Scope.IsUnknown() {
		item.Scope = rule.Scope.ValueStringPointer()
	}
	if !rule.ProtocolPortType.IsNull() && !rule.ProtocolPortType.IsUnknown() {
		item.ProtocolPortType = rule.ProtocolPortType.ValueInt64Pointer()
	}
	if !rule.BelongMemberId.IsNull() && !rule.BelongMemberId.IsUnknown() {
		item.BelongMemberId = rule.BelongMemberId.ValueStringPointer()
	}
	return item
}

// flattenFwmSecurityGroupRuleRule converts the SDK rule detail into the
// computed `rule` element of the model.
func flattenFwmSecurityGroupRuleRule(rule *fwmv20250611.SgRuleResp) FwmSecurityGroupRuleRuleModel {
	return FwmSecurityGroupRuleRuleModel{
		RuleId:            types.StringPointerValue(rule.RuleId),
		Sequence:          types.Int64PointerValue(rule.Sequence),
		Region:            types.StringPointerValue(rule.Region),
		IpVersion:         types.StringPointerValue(rule.IpVersion),
		SrcContent:        types.StringPointerValue(rule.SrcContent),
		SrcType:           types.StringPointerValue(rule.SrcType),
		DstContent:        types.StringPointerValue(rule.DstContent),
		DstType:           types.StringPointerValue(rule.DstType),
		Protocol:          types.StringPointerValue(rule.Protocol),
		ProtocolPortType:  types.Int64PointerValue(rule.ProtocolPortType),
		ServiceTemplateId: types.StringPointerValue(rule.ServiceTemplateId),
		DstPort:           types.StringPointerValue(rule.DstPort),
		RuleAction:        types.StringPointerValue(rule.RuleAction),
		Detail:            types.StringPointerValue(rule.Detail),
		RuleSource:        types.StringPointerValue(rule.RuleSource),
		Scope:             types.StringPointerValue(rule.Scope),
	}
}

// splitFwmSecurityGroupRuleId parses `<group_id>#<rule_id>...#<rule_id>`.
func splitFwmSecurityGroupRuleId(id string) (string, []string, bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) < 2 || parts[0] == "" {
		return "", nil, false
	}

	ruleIds := make([]string, 0, len(parts)-1)
	for _, part := range parts[1:] {
		if part == "" {
			return "", nil, false
		}
		ruleIds = append(ruleIds, part)
	}
	return parts[0], ruleIds, true
}

// fwmStringUseStateForUnknown copies the prior state value into the plan when
// the new plan value is unknown, keeping computed identifiers stable.
type fwmStringUseStateForUnknown struct{}

func (fwmStringUseStateForUnknown) Description(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}

func (fwmStringUseStateForUnknown) MarkdownDescription(ctx context.Context) string {
	return fwmStringUseStateForUnknown{}.Description(ctx)
}

func (fwmStringUseStateForUnknown) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() {
		return
	}
	if !req.PlanValue.IsUnknown() {
		return
	}
	resp.PlanValue = req.StateValue
}
