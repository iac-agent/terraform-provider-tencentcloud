package fwm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// _ resource.Resource is intentionally omitted: the factory function
// NewFwmEdgeAclRuleGroupResource returns resource.Resource, which already
// performs the implicit compile-time check at the return statement. Only the
// extended-interface assertions are kept here.
var (
	_ resource.ResourceWithConfigure   = &FwmEdgeAclRuleGroupResource{}
	_ resource.ResourceWithImportState = &FwmEdgeAclRuleGroupResource{}
)

// NewFwmEdgeAclRuleGroupResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewFwmEdgeAclRuleGroupResource() resource.Resource {
	return &FwmEdgeAclRuleGroupResource{}
}

// FwmEdgeAclRuleGroupResource implements resource.Resource for
// tencentcloud_fwm_edge_acl_rule_group: one internet border rule group
// (product cfw_edge_acl) together with the rules it contains.
//
// Contract notes:
//   - identity form is `single`: the cloud assigned GroupId is the resource id;
//   - group_name is updatable through ModifyRuleGroup;
//   - every other attribute (including the whole rules list) is immutable,
//     because ModifyRuleGroup accepts nothing but GroupId and GroupName. A
//     change of an immutable attribute therefore forces a replacement built
//     from the plan, which avoids relying on an undocumented "absent means
//     clear" semantic;
//   - the rule details are read back from DescribeEdgeAclRules, since
//     DescribeRuleGroups only exposes RuleCount.
type FwmEdgeAclRuleGroupResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// FwmEdgeAclRuleGroupRuleModel maps one element of the rules attribute.
type FwmEdgeAclRuleGroupRuleModel struct {
	RuleId          types.String `tfsdk:"rule_id"`
	OrderIndex      types.Int64  `tfsdk:"order_index"`
	Direction       types.Int64  `tfsdk:"direction"`
	SourceContent   types.String `tfsdk:"source_content"`
	SourceType      types.String `tfsdk:"source_type"`
	TargetContent   types.String `tfsdk:"target_content"`
	TargetType      types.String `tfsdk:"target_type"`
	Port            types.String `tfsdk:"port"`
	Protocol        types.String `tfsdk:"protocol"`
	RuleAction      types.String `tfsdk:"rule_action"`
	Description     types.String `tfsdk:"description"`
	Scope           types.String `tfsdk:"scope"`
	BelongMemberId  types.String `tfsdk:"belong_member_id"`
	ParamTemplateId types.String `tfsdk:"param_template_id"`
}

// FwmEdgeAclRuleGroupResourceModel maps the schema attributes.
type FwmEdgeAclRuleGroupResourceModel struct {
	Id        types.String `tfsdk:"id"`
	GroupId   types.String `tfsdk:"group_id"`
	GroupName types.String `tfsdk:"group_name"`
	Product   types.String `tfsdk:"product"`
	RuleCount types.Int64  `tfsdk:"rule_count"`
	Rules     types.List   `tfsdk:"rules"`
}

// fwmEdgeAclRuleGroupRuleAttrTypes is the attribute type map of one rules
// element, shared by the schema description and the state conversion.
var fwmEdgeAclRuleGroupRuleAttrTypes = map[string]attr.Type{
	"rule_id":           types.StringType,
	"order_index":       types.Int64Type,
	"direction":         types.Int64Type,
	"source_content":    types.StringType,
	"source_type":       types.StringType,
	"target_content":    types.StringType,
	"target_type":       types.StringType,
	"port":              types.StringType,
	"protocol":          types.StringType,
	"rule_action":       types.StringType,
	"description":       types.StringType,
	"scope":             types.StringType,
	"belong_member_id":  types.StringType,
	"param_template_id": types.StringType,
}

func (r *FwmEdgeAclRuleGroupResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_fwm_edge_acl_rule_group"
}

func (r *FwmEdgeAclRuleGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage a FWM internet border rule group (product cfw_edge_acl). " +
			"A rule group bundles the rules that the internet border firewall applies; the provider creates the " +
			"group and its rules in one call and removes both on destroy. " +
			"Only group_name can be updated in place, any change of the rules list replaces the rule group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the rule group, equal to group_id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the rule group assigned by the cloud.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the rule group, 1 to 50 characters. It can be updated in place.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"product": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Product the rule group belongs to. The internet border rule group is fixed to `cfw_edge_acl`; this attribute mirrors the value returned by the cloud.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rule_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of rules in the rule group.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"rules": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Rules of the rule group. Their order is significant: order_index starts at 1 and increases strictly, the smaller the number the higher the priority. Any change replaces the rule group, because the update API cannot modify rules.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rule_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the rule assigned by the cloud.",
						},
						"order_index": schema.Int64Attribute{
							Optional:    true,
							Computed:    true,
							Description: "Execution order of the rule, starting at 1 and strictly increasing. The smaller the number the higher the priority.",
						},
						"direction": schema.Int64Attribute{
							Optional:    true,
							Computed:    true,
							Description: "Direction of the rule: 0 egress, 1 ingress.",
						},
						"source_content": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Source content. Format depends on source_type: IP/CIDR for ip, a domain for domain, a template id for template, a comma separated instance id list for instance, and key:value for tag.",
						},
						"source_type": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Source address type: ip, domain, template, instance, tag.",
						},
						"target_content": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Target content. Format is the same as source_content.",
						},
						"target_type": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Target address type: ip, domain, template, instance, tag.",
						},
						"port": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Port of the rule. Supports a single port, a port range and comma separated ports, such as 80, 1-65535 or 80,443,8080.",
						},
						"protocol": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Protocol of the rule: TCP, UDP, ICMP, ANY.",
						},
						"rule_action": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Action of the rule: accept to allow, drop to block, log to observe.",
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Description of the rule, 0 to 256 characters.",
						},
						"scope": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Effective scope of the rule: serial, side or all.",
						},
						"belong_member_id": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "ID of the member account the rule belongs to, used in multi account scenarios. Required when source_type or target_type is instance or tag.",
						},
						"param_template_id": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "ID of the parameter template used by the rule.",
						},
					},
				},
			},
		},
	}
}

// Create handles `terraform apply` for a new rule group.
func (r *FwmEdgeAclRuleGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FwmEdgeAclRuleGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rules, diags := expandFwmEdgeAclRuleGroupRules(ctx, plan.Rules)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewFwmService(r.Client())
	groupId, err := service.CreateFwmEdgeAclRuleGroup(ctx, plan.GroupName.ValueString(), rules)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating FWM edge ACL rule group",
			fmt.Sprintf("CreateEdgeAclRuleGroup failed: %v", err),
		)
		return
	}

	// Persist the identity before reading back: a failing read keeps the id in
	// state so `terraform destroy` can still clean the rule group up.
	plan.Id = types.StringValue(groupId)
	plan.GroupId = types.StringValue(groupId)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.read(ctx, groupId)
	if err != nil {
		resp.Diagnostics.AddError("Error reading FWM edge ACL rule group after create", err.Error())
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"FWM edge ACL rule group not readable after create",
			fmt.Sprintf("rule group %s was created but cannot be read back", groupId),
		)
		return
	}

	r.refreshModel(ctx, &plan, item, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and the post-apply state refresh.
func (r *FwmEdgeAclRuleGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FwmEdgeAclRuleGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId := state.Id.ValueString()
	item, err := r.read(ctx, groupId)
	if err != nil {
		resp.Diagnostics.AddError("Error reading FWM edge ACL rule group", err.Error())
		return
	}

	// A confirmed "not found" on the cloud side drops the resource from state so
	// Terraform plans a recreation. A nil API response is an error above and
	// never lands here.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("FWM edge ACL rule group %s not found, removing from state.", groupId),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	state.Id = types.StringValue(groupId)
	r.refreshModel(ctx, &state, item, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles the only in-place modification the cloud supports: renaming
// the rule group. Every other attribute carries a RequiresReplace plan
// modifier, so Update only ever receives group_name changes.
func (r *FwmEdgeAclRuleGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FwmEdgeAclRuleGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state FwmEdgeAclRuleGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId := state.Id.ValueString()

	if !plan.GroupName.Equal(state.GroupName) {
		if plan.GroupName.IsUnknown() {
			resp.Diagnostics.AddError(
				"Invalid planned value for group_name",
				"group_name is unknown during apply, so the rule group cannot be renamed.",
			)
			return
		}

		name := plan.GroupName.ValueString()
		if name == "" {
			// ModifyRuleGroup documents no "empty means clear" behaviour, so an
			// explicit clear would either be silently ignored or wipe the name.
			resp.Diagnostics.AddError(
				"Unsupported group_name value",
				"group_name is optional but cannot be cleared once the rule group exists: the API documents no "+
					"clear semantics for an empty name. Set a non-empty name or remove the whole rule group.",
			)
			return
		}

		service := NewFwmService(r.Client())
		if err := service.ModifyFwmEdgeAclRuleGroup(ctx, groupId, name); err != nil {
			resp.Diagnostics.AddError(
				"Error updating FWM edge ACL rule group",
				fmt.Sprintf("ModifyRuleGroup failed for %s: %v", groupId, err),
			)
			return
		}
	}

	item, err := r.read(ctx, groupId)
	if err != nil {
		resp.Diagnostics.AddError("Error reading FWM edge ACL rule group after update", err.Error())
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"FWM edge ACL rule group not readable after update",
			fmt.Sprintf("rule group %s was updated but cannot be read back", groupId),
		)
		return
	}

	plan.Id = state.Id
	r.refreshModel(ctx, &plan, item, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete removes the managed rule group. Only the single managed id is passed
// to the batch delete API.
func (r *FwmEdgeAclRuleGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FwmEdgeAclRuleGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewFwmService(r.Client())
	if err := service.DeleteFwmEdgeAclRuleGroup(ctx, state.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting FWM edge ACL rule group",
			fmt.Sprintf("DeleteRuleGroup failed for %s: %v", state.Id.ValueString(), err),
		)
		return
	}
}

// ImportState is embedded via fw.WithImportByID for `terraform import <id>` passthrough.

// fwmEdgeAclRuleGroupDetail is the aggregated read result of one rule group.
// A nil *fwmEdgeAclRuleGroupDetail returned by read means "does not exist".
type fwmEdgeAclRuleGroupDetail struct {
	Group *fwmEdgeAclRuleGroupGroupDetail
	Rules []*fwmEdgeAclRuleGroupRuleDetail
}

// fwmEdgeAclRuleGroupGroupDetail carries the group level fields.
type fwmEdgeAclRuleGroupGroupDetail struct {
	GroupName string
	Product   string
	RuleCount int64
}

// fwmEdgeAclRuleGroupRuleDetail carries the fields of one rule.
type fwmEdgeAclRuleGroupRuleDetail struct {
	RuleId          string
	OrderIndex      int64
	Direction       int64
	SourceContent   string
	SourceType      string
	TargetContent   string
	TargetType      string
	Port            string
	Protocol        string
	RuleAction      string
	Description     string
	Scope           string
	BelongMemberId  string
	ParamTemplateId string
}

// read fetches one rule group together with its rules.
//
// Two different "empty" outcomes stay distinguishable:
//   - API failure (nil response / nil Response) -> error
//   - queried successfully but nothing matched -> (nil, nil) => "does not exist"
func (r *FwmEdgeAclRuleGroupResource) read(ctx context.Context, groupId string) (*fwmEdgeAclRuleGroupDetail, error) {
	if groupId == "" {
		return nil, fmt.Errorf("rule group id is empty")
	}

	service := NewFwmService(r.Client())

	group, err := service.DescribeFwmEdgeAclRuleGroupById(ctx, groupId)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, nil
	}

	detail := &fwmEdgeAclRuleGroupDetail{
		Group: &fwmEdgeAclRuleGroupGroupDetail{},
	}
	if group.GroupName != nil {
		detail.Group.GroupName = *group.GroupName
	}
	if group.Product != nil {
		detail.Group.Product = *group.Product
	}
	if group.RuleCount != nil {
		detail.Group.RuleCount = *group.RuleCount
	}

	rules, err := service.DescribeFwmEdgeAclRulesByGroupId(ctx, groupId)
	if err != nil {
		return nil, err
	}

	for _, rule := range rules {
		if rule == nil {
			continue
		}
		item := &fwmEdgeAclRuleGroupRuleDetail{}
		if rule.RuleId != nil {
			item.RuleId = *rule.RuleId
		}
		if rule.Sequence != nil {
			item.OrderIndex = *rule.Sequence
		}
		if rule.Direction != nil {
			item.Direction = *rule.Direction
		}
		if rule.SourceContent != nil {
			item.SourceContent = *rule.SourceContent
		}
		if rule.SourceType != nil {
			item.SourceType = *rule.SourceType
		}
		if rule.TargetContent != nil {
			item.TargetContent = *rule.TargetContent
		}
		if rule.TargetType != nil {
			item.TargetType = *rule.TargetType
		}
		if rule.Port != nil {
			item.Port = *rule.Port
		}
		if rule.Protocol != nil {
			item.Protocol = *rule.Protocol
		}
		if rule.RuleAction != nil {
			item.RuleAction = *rule.RuleAction
		}
		if rule.Description != nil {
			item.Description = *rule.Description
		}
		if rule.Scope != nil {
			item.Scope = *rule.Scope
		}
		if rule.BelongMember != nil && rule.BelongMember.MemberId != nil {
			item.BelongMemberId = *rule.BelongMember.MemberId
		}
		if rule.ParamTemplateId != nil {
			item.ParamTemplateId = *rule.ParamTemplateId
		}
		detail.Rules = append(detail.Rules, item)
	}

	return detail, nil
}

// refreshModel writes the cloud state into the model in one place, so Read,
// Create and Update all converge on the same code path.
func (r *FwmEdgeAclRuleGroupResource) refreshModel(
	ctx context.Context,
	m *FwmEdgeAclRuleGroupResourceModel,
	detail *fwmEdgeAclRuleGroupDetail,
	diags *diag.Diagnostics,
) {
	if detail == nil || detail.Group == nil {
		return
	}

	m.GroupName = fwmEdgeAclRuleGroupStringValue(detail.Group.GroupName)
	m.Product = fwmEdgeAclRuleGroupStringValue(detail.Group.Product)
	m.RuleCount = types.Int64Value(detail.Group.RuleCount)

	items := make([]FwmEdgeAclRuleGroupRuleModel, 0, len(detail.Rules))
	for _, rule := range detail.Rules {
		items = append(items, FwmEdgeAclRuleGroupRuleModel{
			RuleId:          fwmEdgeAclRuleGroupStringValue(rule.RuleId),
			OrderIndex:      types.Int64Value(rule.OrderIndex),
			Direction:       types.Int64Value(rule.Direction),
			SourceContent:   fwmEdgeAclRuleGroupStringValue(rule.SourceContent),
			SourceType:      fwmEdgeAclRuleGroupStringValue(rule.SourceType),
			TargetContent:   fwmEdgeAclRuleGroupStringValue(rule.TargetContent),
			TargetType:      fwmEdgeAclRuleGroupStringValue(rule.TargetType),
			Port:            fwmEdgeAclRuleGroupStringValue(rule.Port),
			Protocol:        fwmEdgeAclRuleGroupStringValue(rule.Protocol),
			RuleAction:      fwmEdgeAclRuleGroupStringValue(rule.RuleAction),
			Description:     fwmEdgeAclRuleGroupStringValue(rule.Description),
			Scope:           fwmEdgeAclRuleGroupStringValue(rule.Scope),
			BelongMemberId:  fwmEdgeAclRuleGroupStringValue(rule.BelongMemberId),
			ParamTemplateId: fwmEdgeAclRuleGroupStringValue(rule.ParamTemplateId),
		})
	}

	list, listDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: fwmEdgeAclRuleGroupRuleAttrTypes}, items)
	diags.Append(listDiags...)
	if listDiags.HasError() {
		return
	}
	m.Rules = list
}

// fwmEdgeAclRuleGroupStringValue maps an optional string to its framework
// value. An absent cloud field becomes null instead of an empty string so the
// three value states stay distinguishable.
func fwmEdgeAclRuleGroupStringValue(v string) types.String {
	if v == "" {
		return types.StringNull()
	}
	return types.StringValue(v)
}

// expandFwmEdgeAclRuleGroupRules converts the planned rules list into API edge
// rule infos. rule_id is not sent: a brand new rule group never carries rule
// ids, they are assigned by the cloud.
func expandFwmEdgeAclRuleGroupRules(ctx context.Context, list types.List) ([]*fwmv20250611.EdgeAclRuleInfo, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var items []FwmEdgeAclRuleGroupRuleModel
	diags := list.ElementsAs(ctx, &items, false)
	if diags.HasError() {
		return nil, diags
	}

	rules := make([]*fwmv20250611.EdgeAclRuleInfo, 0, len(items))
	for i, item := range items {
		rule := &fwmv20250611.EdgeAclRuleInfo{}
		rule.OrderIndex = fwmEdgeAclRuleGroupInt64Pointer(item.OrderIndex)
		rule.Direction = fwmEdgeAclRuleGroupInt64Pointer(item.Direction)
		rule.SourceContent = fwmEdgeAclRuleGroupStringPointer(item.SourceContent)
		rule.SourceType = fwmEdgeAclRuleGroupStringPointer(item.SourceType)
		rule.TargetContent = fwmEdgeAclRuleGroupStringPointer(item.TargetContent)
		rule.TargetType = fwmEdgeAclRuleGroupStringPointer(item.TargetType)
		rule.Port = fwmEdgeAclRuleGroupStringPointer(item.Port)
		rule.Protocol = fwmEdgeAclRuleGroupStringPointer(item.Protocol)
		rule.RuleAction = fwmEdgeAclRuleGroupStringPointer(item.RuleAction)
		rule.Description = fwmEdgeAclRuleGroupStringPointer(item.Description)
		rule.Scope = fwmEdgeAclRuleGroupStringPointer(item.Scope)
		rule.BelongMemberId = fwmEdgeAclRuleGroupStringPointer(item.BelongMemberId)
		rule.ParamTemplateId = fwmEdgeAclRuleGroupStringPointer(item.ParamTemplateId)

		if rule.OrderIndex == nil {
			// The API requires order_index to start at 1 and increase strictly;
			// falling back to the position keeps a rule list without explicit
			// order_index usable and deterministic.
			rule.OrderIndex = helper.IntInt64(i + 1)
		}

		rules = append(rules, rule)
	}

	return rules, diags
}

// fwmEdgeAclRuleGroupStringPointer returns nil for null and unknown values so
// an unset attribute is never sent as an empty string.
func fwmEdgeAclRuleGroupStringPointer(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	return helper.String(v.ValueString())
}

func fwmEdgeAclRuleGroupInt64Pointer(v types.Int64) *int64 {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	return helper.Int64(v.ValueInt64())
}
