package fwm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// _ resource.Resource is intentionally omitted: the factory function
// NewFwmSecurityGroupRuleResource returns resource.Resource, which already
// performs the implicit compile-time check at the return statement. Only the
// extended-interface assertions are kept here.
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
//
// One Terraform resource manages exactly one enterprise security group rule.
// The remote identifier is composed of the owning rule group ID and the rule
// ID as `<group_id>#<rule_id>`, because both values are required by every
// rule level cloud API.
type FwmSecurityGroupRuleResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// FwmSecurityGroupRuleRuleModel maps one element of the `rules` attribute.
type FwmSecurityGroupRuleRuleModel struct {
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

// FwmSecurityGroupRuleDetailModel maps one element of the computed `rule`
// attribute, which mirrors the cloud side rule detail.
type FwmSecurityGroupRuleDetailModel struct {
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

// FwmSecurityGroupRuleResourceModel maps the schema attributes.
type FwmSecurityGroupRuleResourceModel struct {
	Id      types.String                      `tfsdk:"id"`
	GroupId types.String                      `tfsdk:"group_id"`
	Rules   []FwmSecurityGroupRuleRuleModel   `tfsdk:"rules"`
	Rule    []FwmSecurityGroupRuleDetailModel `tfsdk:"rule"`
}

func (r *FwmSecurityGroupRuleResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_fwm_security_group_rule"
}

func (r *FwmSecurityGroupRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage a FWM (Firewall Manager) enterprise security group rule. " +
			"The remote identifier of this resource is `<group_id>#<rule_id>`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier of the security group rule, in the form `<group_id>#<rule_id>`.",
				PlanModifiers: []planmodifier.String{
					stringUseStateForUnknownRule{},
				},
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Description: "ID of the rule group the security group rule belongs to.",
				PlanModifiers: []planmodifier.String{
					stringRequiresReplaceRule{},
				},
			},
			"rules": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Rule list, exactly one element is supported and it describes the rule to create.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip_version": schema.StringAttribute{
							Optional:    true,
							Description: "IP version of the rule.",
						},
						"source_content": schema.StringAttribute{
							Optional:    true,
							Description: "Source address of the rule.",
						},
						"source_type": schema.StringAttribute{
							Optional:    true,
							Description: "Source type of the rule.",
						},
						"dest_content": schema.StringAttribute{
							Optional:    true,
							Description: "Destination address of the rule.",
						},
						"dest_type": schema.StringAttribute{
							Optional:    true,
							Description: "Destination type of the rule.",
						},
						"protocol": schema.StringAttribute{
							Optional:    true,
							Description: "Protocol of the rule.",
						},
						"port": schema.StringAttribute{
							Optional:    true,
							Description: "Port of the rule.",
						},
						"service_template_id": schema.StringAttribute{
							Optional:    true,
							Description: "ID of the protocol and port template used by the rule.",
						},
						"rule_action": schema.StringAttribute{
							Optional:    true,
							Description: "Action of the rule, `1` for block and `2` for allow.",
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Description: "Description of the rule.",
						},
						"order_index": schema.Int64Attribute{
							Optional:    true,
							Description: "Priority of the rule.",
						},
						"rule_id": schema.StringAttribute{
							Optional:    true,
							Description: "ID of the rule. Leave it empty on creation; the cloud side generates the ID.",
						},
						"scope": schema.StringAttribute{
							Optional:    true,
							Description: "Effective scope of the rule.",
						},
						"protocol_port_type": schema.Int64Attribute{
							Optional:    true,
							Description: "Protocol and port type of the rule.",
						},
						"belong_member_id": schema.StringAttribute{
							Optional:    true,
							Description: "ID of the member account the rule belongs to. It is required when the group is a `cfwg-xxx` group or when the source type or destination type is `instance` or `tag`.",
						},
					},
				},
			},
			"rule": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Detail of the security group rule returned by the cloud side.",
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
							Description: "Protocol and port type of the rule.",
						},
						"service_template_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the protocol and port template used by the rule.",
						},
						"dst_port": schema.StringAttribute{
							Computed:    true,
							Description: "Destination port of the rule, a single port, multiple ports or a port range.",
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
							Description: "Source of the rule, `0` means it was created from the user console.",
						},
						"scope": schema.StringAttribute{
							Computed:    true,
							Description: "Effective scope of the rule, a bit mask where `1` means enterprise security group and `2` means lighthouse instance.",
						},
					},
				},
			},
		},
	}
}

// Create handles `terraform apply` for a new security group rule.
func (r *FwmSecurityGroupRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FwmSecurityGroupRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(plan.Rules) != 1 {
		resp.Diagnostics.AddError(
			"Invalid rules",
			fmt.Sprintf("tencentcloud_fwm_security_group_rule manages exactly one rule, but %d element(s) were given in `rules`.", len(plan.Rules)),
		)
		return
	}

	groupId, ok := resolveGroupIdOfRule(plan)
	if !ok {
		resp.Diagnostics.AddError(
			"Missing group_id",
			"`group_id` must be set to the ID of the rule group the security group rule belongs to.",
		)
		return
	}

	plan.GroupId = types.StringValue(groupId)

	request := &fwmv20250611.CreateSecurityGroupRuleRequest{
		GroupId: &groupId,
		Rules:   []*fwmv20250611.SecurityGroupRule{buildSecurityGroupRule(plan.Rules[0])},
	}

	ruleIds, err := NewFwmService(r.Client()).CreateFwmSecurityGroupRule(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating FWM security group rule",
			fmt.Sprintf("CreateSecurityGroupRule failed for group %s: %v", groupId, err),
		)
		return
	}

	if len(ruleIds) == 0 || ruleIds[0] == nil || *ruleIds[0] == "" {
		resp.Diagnostics.AddError(
			"Error creating FWM security group rule",
			fmt.Sprintf("CreateSecurityGroupRule returned an empty rule ID for group %s.", groupId),
		)
		return
	}

	ruleId := *ruleIds[0]
	plan.Id = types.StringValue(fmt.Sprintf("%s%s%s", groupId, tccommon.FILED_SP, ruleId))

	// Read back immediately after creation so that state matches the cloud side.
	detail, err := NewFwmService(r.Client()).DescribeFwmSecurityGroupRuleById(ctx, groupId, ruleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading FWM security group rule after create",
			fmt.Sprintf("DescribeSecurityGroupRule failed for rule %s of group %s: %v", ruleId, groupId, err),
		)
		return
	}
	setSecurityGroupRuleDetail(&plan, detail)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh`, `terraform plan` refresh and imports.
func (r *FwmSecurityGroupRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FwmSecurityGroupRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId, ruleId, ok := splitSecurityGroupRuleId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource identifier",
			fmt.Sprintf("Expected an identifier in the form `<group_id>#<rule_id>`, got %q.", state.Id.ValueString()),
		)
		return
	}

	if !state.GroupId.IsNull() && !state.GroupId.IsUnknown() && state.GroupId.ValueString() != "" {
		groupId = state.GroupId.ValueString()
	}

	detail, err := NewFwmService(r.Client()).DescribeFwmSecurityGroupRuleById(ctx, groupId, ruleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading FWM security group rule",
			fmt.Sprintf("DescribeSecurityGroupRule failed for rule %s of group %s: %v", ruleId, groupId, err),
		)
		return
	}

	if detail == nil {
		log.Printf("[WARN]%s resource `tencentcloud_fwm_security_group_rule` [%s] not found, please check if it has been deleted.", tccommon.GetLogId(ctx), state.Id.ValueString())
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("tencentcloud_fwm_security_group_rule %s not found, removing from state.", state.Id.ValueString()),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	state.GroupId = types.StringValue(groupId)
	state.Rule = []FwmSecurityGroupRuleDetailModel{toSecurityGroupRuleDetailModel(detail)}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates of the rule.
func (r *FwmSecurityGroupRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FwmSecurityGroupRuleResourceModel
	var state FwmSecurityGroupRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId, ruleId, ok := splitSecurityGroupRuleId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource identifier",
			fmt.Sprintf("Expected an identifier in the form `<group_id>#<rule_id>`, got %q.", state.Id.ValueString()),
		)
		return
	}

	if !plan.GroupId.IsNull() && !plan.GroupId.IsUnknown() && plan.GroupId.ValueString() != "" {
		groupId = plan.GroupId.ValueString()
	}

	plannedRule := FwmSecurityGroupRuleRuleModel{}
	if len(plan.Rules) > 0 {
		plannedRule = plan.Rules[0]
	}

	// The cloud API identifies the rule to modify by the ID inside the rule body.
	if plannedRule.RuleId.IsNull() || plannedRule.RuleId.IsUnknown() || plannedRule.RuleId.ValueString() == "" {
		plannedRule.RuleId = types.StringValue(ruleId)
	}

	request := &fwmv20250611.ModifySecurityGroupRuleRequest{
		GroupId: &groupId,
		Rule:    buildSecurityGroupRule(plannedRule),
	}

	if err := NewFwmService(r.Client()).ModifyFwmSecurityGroupRule(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating FWM security group rule",
			fmt.Sprintf("ModifySecurityGroupRule failed for rule %s of group %s: %v", ruleId, groupId, err),
		)
		return
	}

	detail, err := NewFwmService(r.Client()).DescribeFwmSecurityGroupRuleById(ctx, groupId, ruleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading FWM security group rule after update",
			fmt.Sprintf("DescribeSecurityGroupRule failed for rule %s of group %s: %v", ruleId, groupId, err),
		)
		return
	}
	setSecurityGroupRuleDetail(&plan, detail)

	if len(plan.Rules) > 0 {
		plan.Rules[0].RuleId = types.StringValue(ruleId)
	}
	plan.GroupId = types.StringValue(groupId)
	plan.Id = types.StringValue(fmt.Sprintf("%s%s%s", groupId, tccommon.FILED_SP, ruleId))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles rule deletion.
func (r *FwmSecurityGroupRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FwmSecurityGroupRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId, ruleId, ok := splitSecurityGroupRuleId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource identifier",
			fmt.Sprintf("Expected an identifier in the form `<group_id>#<rule_id>`, got %q.", state.Id.ValueString()),
		)
		return
	}

	if !state.GroupId.IsNull() && !state.GroupId.IsUnknown() && state.GroupId.ValueString() != "" {
		groupId = state.GroupId.ValueString()
	}

	if err := NewFwmService(r.Client()).DeleteFwmSecurityGroupRuleById(ctx, groupId, ruleId); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting FWM security group rule",
			fmt.Sprintf("DeleteSecurityGroupRule failed for rule %s of group %s: %v", ruleId, groupId, err),
		)
		return
	}

	// The framework removes the resource from state after a successful delete.
}

// ImportState is embedded via fw.WithImportByID for `terraform import <id>` passthrough,
// where `<id>` must be `<group_id>#<rule_id>`.

// resolveGroupIdOfRule returns the rule group ID configured on the resource.
func resolveGroupIdOfRule(plan FwmSecurityGroupRuleResourceModel) (string, bool) {
	if plan.GroupId.IsNull() || plan.GroupId.IsUnknown() {
		return "", false
	}

	groupId := plan.GroupId.ValueString()
	if groupId == "" {
		return "", false
	}

	return groupId, true
}

// splitSecurityGroupRuleId parses `<group_id>#<rule_id>`.
func splitSecurityGroupRuleId(id string) (string, string, bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}

	return parts[0], parts[1], true
}

// buildSecurityGroupRule converts the model of one rule into the SDK struct.
func buildSecurityGroupRule(item FwmSecurityGroupRuleRuleModel) *fwmv20250611.SecurityGroupRule {
	rule := &fwmv20250611.SecurityGroupRule{}

	if !item.IpVersion.IsNull() && !item.IpVersion.IsUnknown() {
		rule.IpVersion = item.IpVersion.ValueStringPointer()
	}
	if !item.SourceContent.IsNull() && !item.SourceContent.IsUnknown() {
		rule.SourceContent = item.SourceContent.ValueStringPointer()
	}
	if !item.SourceType.IsNull() && !item.SourceType.IsUnknown() {
		rule.SourceType = item.SourceType.ValueStringPointer()
	}
	if !item.DestContent.IsNull() && !item.DestContent.IsUnknown() {
		rule.DestContent = item.DestContent.ValueStringPointer()
	}
	if !item.DestType.IsNull() && !item.DestType.IsUnknown() {
		rule.DestType = item.DestType.ValueStringPointer()
	}
	if !item.Protocol.IsNull() && !item.Protocol.IsUnknown() {
		rule.Protocol = item.Protocol.ValueStringPointer()
	}
	if !item.Port.IsNull() && !item.Port.IsUnknown() {
		rule.Port = item.Port.ValueStringPointer()
	}
	if !item.ServiceTemplateId.IsNull() && !item.ServiceTemplateId.IsUnknown() {
		rule.ServiceTemplateId = item.ServiceTemplateId.ValueStringPointer()
	}
	if !item.RuleAction.IsNull() && !item.RuleAction.IsUnknown() {
		rule.RuleAction = item.RuleAction.ValueStringPointer()
	}
	if !item.Description.IsNull() && !item.Description.IsUnknown() {
		rule.Description = item.Description.ValueStringPointer()
	}
	if !item.OrderIndex.IsNull() && !item.OrderIndex.IsUnknown() {
		rule.OrderIndex = item.OrderIndex.ValueInt64Pointer()
	}
	if !item.RuleId.IsNull() && !item.RuleId.IsUnknown() {
		rule.RuleId = item.RuleId.ValueStringPointer()
	}
	if !item.Scope.IsNull() && !item.Scope.IsUnknown() {
		rule.Scope = item.Scope.ValueStringPointer()
	}
	if !item.ProtocolPortType.IsNull() && !item.ProtocolPortType.IsUnknown() {
		rule.ProtocolPortType = item.ProtocolPortType.ValueInt64Pointer()
	}
	if !item.BelongMemberId.IsNull() && !item.BelongMemberId.IsUnknown() {
		rule.BelongMemberId = item.BelongMemberId.ValueStringPointer()
	}

	return rule
}

// toSecurityGroupRuleDetailModel converts the cloud side rule detail into the model.
func toSecurityGroupRuleDetailModel(item *fwmv20250611.SgRuleResp) FwmSecurityGroupRuleDetailModel {
	if item == nil {
		return FwmSecurityGroupRuleDetailModel{}
	}

	return FwmSecurityGroupRuleDetailModel{
		RuleId:            types.StringPointerValue(item.RuleId),
		Sequence:          types.Int64PointerValue(item.Sequence),
		Region:            types.StringPointerValue(item.Region),
		IpVersion:         types.StringPointerValue(item.IpVersion),
		SrcContent:        types.StringPointerValue(item.SrcContent),
		SrcType:           types.StringPointerValue(item.SrcType),
		DstContent:        types.StringPointerValue(item.DstContent),
		DstType:           types.StringPointerValue(item.DstType),
		Protocol:          types.StringPointerValue(item.Protocol),
		ProtocolPortType:  types.Int64PointerValue(item.ProtocolPortType),
		ServiceTemplateId: types.StringPointerValue(item.ServiceTemplateId),
		DstPort:           types.StringPointerValue(item.DstPort),
		RuleAction:        types.StringPointerValue(item.RuleAction),
		Detail:            types.StringPointerValue(item.Detail),
		RuleSource:        types.StringPointerValue(item.RuleSource),
		Scope:             types.StringPointerValue(item.Scope),
	}
}

// setSecurityGroupRuleDetail writes the computed rule detail into the model.
func setSecurityGroupRuleDetail(model *FwmSecurityGroupRuleResourceModel, item *fwmv20250611.SgRuleResp) {
	model.Rule = []FwmSecurityGroupRuleDetailModel{toSecurityGroupRuleDetailModel(item)}
}

// stringRequiresReplaceRule forces resource replacement when a string attribute changes.
type stringRequiresReplaceRule struct{}

func (stringRequiresReplaceRule) Description(_ context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}

func (stringRequiresReplaceRule) MarkdownDescription(ctx context.Context) string {
	return stringRequiresReplaceRule{}.Description(ctx)
}

func (stringRequiresReplaceRule) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() {
		// Resource is being created; nothing to replace.
		return
	}
	if req.PlanValue.IsNull() {
		// Resource is being destroyed; nothing to replace.
		return
	}
	if !req.PlanValue.Equal(req.StateValue) {
		resp.RequiresReplace = true
	}
}

// stringUseStateForUnknownRule copies the prior state value into the plan when
// the new plan value is unknown, used for the computed `id` attribute so that
// it stays stable across refreshes.
type stringUseStateForUnknownRule struct{}

func (stringUseStateForUnknownRule) Description(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}

func (stringUseStateForUnknownRule) MarkdownDescription(ctx context.Context) string {
	return stringUseStateForUnknownRule{}.Description(ctx)
}

func (stringUseStateForUnknownRule) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() {
		return
	}
	if !req.PlanValue.IsUnknown() {
		return
	}
	resp.PlanValue = req.StateValue
}
