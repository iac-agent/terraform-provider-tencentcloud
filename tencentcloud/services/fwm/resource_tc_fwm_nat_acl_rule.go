package fwm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// _ resource.Resource is intentionally omitted: the factory function
// NewFwmNatAclRuleResource returns resource.Resource, which already performs
// the implicit compile-time check at the return statement. Only the
// extended-interface assertions are kept here.
var (
	_ resource.ResourceWithConfigure   = &FwmNatAclRuleResource{}
	_ resource.ResourceWithImportState = &FwmNatAclRuleResource{}
)

// NewFwmNatAclRuleResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewFwmNatAclRuleResource() resource.Resource {
	return &FwmNatAclRuleResource{}
}

// FwmNatAclRuleResource implements resource.Resource for
// tencentcloud_fwm_nat_acl_rule.
type FwmNatAclRuleResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// FwmNatAclRuleResourceModel maps the schema attributes.
//
// The cloud API works on a single NAT boundary ACL rule per call, so the
// resource exposes the flattened rule body (Rule / Rules in the cloud API)
// directly at the top level. RuleIds is the create response and only carries
// a single element; `rule_id` (part of the composite resource id) is the
// identifier actually used by Read / Update / Delete.
type FwmNatAclRuleResourceModel struct {
	Id              types.String `tfsdk:"id"`
	GroupId         types.String `tfsdk:"group_id"`
	RuleIds         types.List   `tfsdk:"rule_ids"`
	RuleId          types.String `tfsdk:"rule_id"`
	SourceContent   types.String `tfsdk:"source_content"`
	SourceType      types.String `tfsdk:"source_type"`
	TargetContent   types.String `tfsdk:"target_content"`
	TargetType      types.String `tfsdk:"target_type"`
	Protocol        types.String `tfsdk:"protocol"`
	RuleAction      types.String `tfsdk:"rule_action"`
	OrderIndex      types.Int64  `tfsdk:"order_index"`
	Scope           types.String `tfsdk:"scope"`
	Direction       types.Int64  `tfsdk:"direction"`
	Port            types.String `tfsdk:"port"`
	Description     types.String `tfsdk:"description"`
	ParamTemplateId types.String `tfsdk:"param_template_id"`
	BelongMemberId  types.String `tfsdk:"belong_member_id"`

	// Read-only attributes returned by DescribeNatAclRules.
	Sequence          types.Int64  `tfsdk:"sequence"`
	SourceName        types.String `tfsdk:"source_name"`
	TargetName        types.String `tfsdk:"target_name"`
	ScopeDesc         types.String `tfsdk:"scope_desc"`
	FwInsId           types.String `tfsdk:"fw_ins_id"`
	CountryName       types.String `tfsdk:"country_name"`
	CityName          types.String `tfsdk:"city_name"`
	CountryCode       types.Int64  `tfsdk:"country_code"`
	CityCode          types.Int64  `tfsdk:"city_code"`
	CountryKey        types.String `tfsdk:"country_key"`
	CityKey           types.String `tfsdk:"city_key"`
	IsRegion          types.Int64  `tfsdk:"is_region"`
	CloudCode         types.String `tfsdk:"cloud_code"`
	IsCloud           types.Int64  `tfsdk:"is_cloud"`
	ParamTemplateName types.String `tfsdk:"param_template_name"`
	Invalid           types.Int64  `tfsdk:"invalid"`
	BelongMember      types.Object `tfsdk:"belong_member"`
	CreateTime        types.String `tfsdk:"create_time"`
	UpdateTime        types.String `tfsdk:"update_time"`
	DnsParseCnt       types.Int64  `tfsdk:"dns_parse_cnt"`
	HitCount          types.Int64  `tfsdk:"hit_count"`
	LastHitTime       types.String `tfsdk:"last_hit_time"`
	AllTotalCount     types.Int64  `tfsdk:"all_total_count"`
}

// fwmNatAclRuleBelongMemberModel maps the read-only BelongMember block.
type fwmNatAclRuleBelongMemberModel struct {
	AppId    types.String `tfsdk:"app_id"`
	Uin      types.String `tfsdk:"uin"`
	Nickname types.String `tfsdk:"nickname"`
	MemberId types.String `tfsdk:"member_id"`
	NodeName types.String `tfsdk:"node_name"`
}

func (r *FwmNatAclRuleResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_fwm_nat_acl_rule"
}

func (r *FwmNatAclRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage a FWM (Firewall Manager) NAT boundary ACL rule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier in the form `<group_id>#<rule_id>`.",
				PlanModifiers: []planmodifier.String{
					stringUseStateForUnknown{},
				},
			},
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the NAT boundary ACL rule group the rule belongs to.",
				PlanModifiers: []planmodifier.String{
					stringRequiresReplace{},
				},
			},
			"rule_ids": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Rule ID list returned by the create API. The cloud API creates one rule per call, so it holds a single element.",
			},
			"rule_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "ID of the rule. It is returned by the create API and is required when the rule is modified; it is the identifier used by read, update and delete.",
				PlanModifiers: []planmodifier.String{
					stringUseStateForUnknown{},
				},
			},
			"source_content": schema.StringAttribute{
				Required:    true,
				Description: "Source address content.",
			},
			"source_type": schema.StringAttribute{
				Required:    true,
				Description: "Source type, valid values: `ip`, `url`, `template`, `instance`, `tag`.",
			},
			"target_content": schema.StringAttribute{
				Required:    true,
				Description: "Target address content.",
			},
			"target_type": schema.StringAttribute{
				Required:    true,
				Description: "Target type, valid values: `ip`, `url`, `template`, `instance`, `tag`.",
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "Protocol, such as `TCP`, `UDP`, `ICMP`, `ANY`, `HTTP`, `HTTPS`, `DNS`, `FTP`.",
			},
			"rule_action": schema.StringAttribute{
				Required:    true,
				Description: "Action of the rule, valid values: `accept`, `drop`, `log`.",
			},
			"order_index": schema.Int64Attribute{
				Required:    true,
				Description: "Priority of the rule, starting from 1.",
			},
			"scope": schema.StringAttribute{
				Required:    true,
				Description: "Effective scope of the rule: `ALL` takes effect globally, `ap-xxx` takes effect in a region, `cfwnat-xxx` takes effect on a NAT firewall instance.",
			},
			"direction": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Direction of the rule: `0` outbound, `1` inbound.",
			},
			"port": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Port of the rule. It is empty when the protocol is `ICMP`.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Description of the rule.",
			},
			"param_template_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "ID of the port template.",
			},
			"belong_member_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "ID of the member account the rule belongs to. Required when `scope` is `cfwnat-xxx` or `source_type`/`target_type` is `instance`/`tag`.",
			},
			"sequence": schema.Int64Attribute{
				Computed:    true,
				Description: "Priority of the rule as returned by the cloud API.",
			},
			"source_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the source asset.",
			},
			"target_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the target asset.",
			},
			"scope_desc": schema.StringAttribute{
				Computed:    true,
				Description: "Description of the effective scope.",
			},
			"fw_ins_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the firewall instance.",
			},
			"country_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the country.",
			},
			"city_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the city.",
			},
			"country_code": schema.Int64Attribute{
				Computed:    true,
				Description: "Code of the country.",
			},
			"city_code": schema.Int64Attribute{
				Computed:    true,
				Description: "Code of the city.",
			},
			"country_key": schema.StringAttribute{
				Computed:    true,
				Description: "Key of the country.",
			},
			"city_key": schema.StringAttribute{
				Computed:    true,
				Description: "Key of the city.",
			},
			"is_region": schema.Int64Attribute{
				Computed:    true,
				Description: "Whether the rule is a region rule: `0` no, `1` yes.",
			},
			"cloud_code": schema.StringAttribute{
				Computed:    true,
				Description: "Code of the cloud provider.",
			},
			"is_cloud": schema.Int64Attribute{
				Computed:    true,
				Description: "Whether the rule is a cloud provider rule: `0` no, `1` yes.",
			},
			"param_template_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the port template.",
			},
			"invalid": schema.Int64Attribute{
				Computed:    true,
				Description: "Whether the rule is invalid: `0` valid, `1` invalid.",
			},
			"belong_member": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Member account the rule belongs to.",
				Attributes: map[string]schema.Attribute{
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
						Description: "Nickname of the member account.",
					},
					"member_id": schema.StringAttribute{
						Computed:    true,
						Description: "ID of the member account. It is empty in the personal account scenario.",
					},
					"node_name": schema.StringAttribute{
						Computed:    true,
						Description: "Organization node the member account belongs to.",
					},
				},
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the rule.",
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "Last update time of the rule.",
			},
			"dns_parse_cnt": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of parsed domains.",
			},
			"hit_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Cumulative hit count of the rule across all delivered member accounts.",
			},
			"last_hit_time": schema.StringAttribute{
				Computed:    true,
				Description: "Latest hit time of the rule across all delivered member accounts. It is `-` when the rule has never been hit.",
			},
			"all_total_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Total number of rules in the rule group without filters.",
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *FwmNatAclRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FwmNatAclRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId := plan.GroupId.ValueString()

	request := fwmv20250611.NewCreateNatAclRuleRequest()
	request.GroupId = plan.GroupId.ValueStringPointer()
	request.Rules = []*fwmv20250611.NatAclRule{plan.toRuleParam()}

	log.Printf("[DEBUG] tencentcloud_fwm_nat_acl_rule create, group_id=%s", groupId)
	service := NewFwmService(r.Client())
	result, err := service.CreateFwmNatAclRule(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating FWM NAT ACL rule",
			fmt.Sprintf("CreateNatAclRule failed for group %q: %v", groupId, err),
		)
		return
	}

	if result == nil || result.Response == nil {
		resp.Diagnostics.AddError(
			"Error creating FWM NAT ACL rule",
			fmt.Sprintf("CreateNatAclRule returned an empty response for group %q", groupId),
		)
		return
	}

	ruleId := ""
	for _, id := range result.Response.RuleIds {
		if id != nil && *id != "" {
			ruleId = *id
			break
		}
	}
	if ruleId == "" {
		resp.Diagnostics.AddError(
			"Error creating FWM NAT ACL rule",
			fmt.Sprintf("CreateNatAclRule returned an empty rule id for group %q, request id %s", groupId, stringValue(result.Response.RequestId)),
		)
		return
	}

	plan.RuleId = types.StringValue(ruleId)
	plan.RuleIds = stringListValue(result.Response.RuleIds)
	plan.Id = types.StringValue(buildFwmNatAclRuleId(groupId, ruleId))
	if plan.BelongMember.IsNull() {
		plan.BelongMember = types.ObjectNull(fwmNatAclRuleBelongMemberAttrTypes())
	}

	// Persist the id before reading back so that a failed refresh still leaves a
	// destroyable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.read(ctx, groupId, ruleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading FWM NAT ACL rule after create",
			fmt.Sprintf("DescribeNatAclRules failed for group %q rule %q: %v", groupId, ruleId, err),
		)
		return
	}
	if rule == nil {
		resp.Diagnostics.AddError(
			"Error reading FWM NAT ACL rule after create",
			fmt.Sprintf("FWM NAT ACL rule %q of group %q was not found right after creation", ruleId, groupId),
		)
		return
	}

	r.refreshModel(ctx, &plan, rule)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *FwmNatAclRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FwmNatAclRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId, ruleId, ok := splitFwmNatAclRuleId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("FWM NAT ACL rule id is broken, id is %q (expected `<group_id>%s<rule_id>`)", state.Id.ValueString(), tccommon.FILED_SP),
		)
		return
	}

	rule, err := r.read(ctx, groupId, ruleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading FWM NAT ACL rule",
			fmt.Sprintf("DescribeNatAclRules failed for group %q rule %q: %v", groupId, ruleId, err),
		)
		return
	}

	if rule == nil {
		log.Printf("[WARN] tencentcloud_fwm_nat_acl_rule id=%s group_id=%s rule_id=%s not found, removing from state", state.Id.ValueString(), groupId, ruleId)
		resp.Diagnostics.AddWarning(
			"FWM NAT ACL rule not found",
			fmt.Sprintf("FWM NAT ACL rule %q of group %q does not exist any more, removing it from state.", ruleId, groupId),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	state.GroupId = types.StringValue(groupId)
	state.RuleId = types.StringValue(ruleId)
	if state.RuleIds.IsNull() || state.RuleIds.IsUnknown() {
		state.RuleIds = stringListValue([]*string{&ruleId})
	}

	r.refreshModel(ctx, &state, rule)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates.
func (r *FwmNatAclRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FwmNatAclRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state FwmNatAclRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId, ruleId, ok := splitFwmNatAclRuleId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("FWM NAT ACL rule id is broken, id is %q (expected `<group_id>%s<rule_id>`)", state.Id.ValueString(), tccommon.FILED_SP),
		)
		return
	}

	ruleParam := plan.toRuleParam()
	ruleParam.RuleId = &ruleId

	request := fwmv20250611.NewModifyNatAclRuleRequest()
	request.GroupId = &groupId
	request.Rule = ruleParam

	log.Printf("[DEBUG] tencentcloud_fwm_nat_acl_rule update, group_id=%s rule_id=%s", groupId, ruleId)
	service := NewFwmService(r.Client())
	result, err := service.ModifyFwmNatAclRule(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating FWM NAT ACL rule",
			fmt.Sprintf("ModifyNatAclRule failed for group %q rule %q: %v", groupId, ruleId, err),
		)
		return
	}
	if result == nil || result.Response == nil {
		resp.Diagnostics.AddError(
			"Error updating FWM NAT ACL rule",
			fmt.Sprintf("ModifyNatAclRule returned an empty response for group %q rule %q", groupId, ruleId),
		)
		return
	}

	plan.Id = state.Id
	plan.GroupId = types.StringValue(groupId)
	plan.RuleId = types.StringValue(ruleId)
	plan.RuleIds = state.RuleIds

	rule, err := r.read(ctx, groupId, ruleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading FWM NAT ACL rule after update",
			fmt.Sprintf("DescribeNatAclRules failed for group %q rule %q: %v", groupId, ruleId, err),
		)
		return
	}
	if rule == nil {
		log.Printf("[WARN] tencentcloud_fwm_nat_acl_rule id=%s group_id=%s rule_id=%s not found after update", plan.Id.ValueString(), groupId, ruleId)
		resp.Diagnostics.AddWarning(
			"FWM NAT ACL rule not found",
			fmt.Sprintf("FWM NAT ACL rule %q of group %q does not exist any more, removing it from state.", ruleId, groupId),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	r.refreshModel(ctx, &plan, rule)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles `terraform destroy`.
func (r *FwmNatAclRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FwmNatAclRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupId, ruleId, ok := splitFwmNatAclRuleId(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("FWM NAT ACL rule id is broken, id is %q (expected `<group_id>%s<rule_id>`)", state.Id.ValueString(), tccommon.FILED_SP),
		)
		return
	}

	log.Printf("[DEBUG] tencentcloud_fwm_nat_acl_rule delete, group_id=%s rule_id=%s", groupId, ruleId)
	service := NewFwmService(r.Client())
	if err := service.DeleteFwmNatAclRuleByRuleId(ctx, groupId, ruleId); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting FWM NAT ACL rule",
			fmt.Sprintf("DeleteNatAclRule failed for group %q rule %q: %v", groupId, ruleId, err),
		)
		return
	}
}

// ImportState is embedded via fw.WithImportByID for `<group_id>#<rule_id>` id passthrough.

// read fetches one NAT ACL rule of the given group; it returns nil when the
// cloud side no longer has the rule.
func (r *FwmNatAclRuleResource) read(ctx context.Context, groupId, ruleId string) (*fwmv20250611.NatAclRuleResp, error) {
	service := NewFwmService(r.Client())
	rule, err := service.DescribeFwmNatAclRuleByRuleId(ctx, groupId, ruleId)
	if err != nil {
		if tccommon.IsExpectError(err, []string{"ResourceNotFound", "ResourceNotFound.RuleNotFound"}) {
			return nil, nil
		}
		return nil, err
	}
	return rule, nil
}

// refreshModel copies cloud-owned fields of the rule into the model. User
// supplied fields are left untouched when the cloud side returns nothing, so
// that the plan value is preserved and no spurious diff is produced.
func (r *FwmNatAclRuleResource) refreshModel(ctx context.Context, m *FwmNatAclRuleResourceModel, rule *fwmv20250611.NatAclRuleResp) {
	if rule == nil {
		return
	}

	setStringIfNotEmpty(&m.RuleId, rule.RuleId)
	setStringIfNotEmpty(&m.SourceContent, rule.SourceContent)
	setStringIfNotEmpty(&m.SourceType, rule.SourceType)
	setStringIfNotEmpty(&m.TargetContent, rule.TargetContent)
	setStringIfNotEmpty(&m.TargetType, rule.TargetType)
	setStringIfNotEmpty(&m.Protocol, rule.Protocol)
	setStringIfNotEmpty(&m.RuleAction, rule.RuleAction)
	setStringIfNotEmpty(&m.Scope, rule.Scope)
	setStringIfNotEmpty(&m.Port, rule.Port)
	setStringIfNotEmpty(&m.Description, rule.Description)
	setStringIfNotEmpty(&m.ParamTemplateId, rule.ParamTemplateId)
	setInt64IfNotNil(&m.OrderIndex, rule.Sequence)
	setInt64IfNotNil(&m.Direction, rule.Direction)

	setInt64Pointer(&m.Sequence, rule.Sequence)
	setStringPointer(&m.SourceName, rule.SourceName)
	setStringPointer(&m.TargetName, rule.TargetName)
	setStringPointer(&m.ScopeDesc, rule.ScopeDesc)
	setStringPointer(&m.FwInsId, rule.FwInsId)
	setStringPointer(&m.CountryName, rule.CountryName)
	setStringPointer(&m.CityName, rule.CityName)
	setInt64Pointer(&m.CountryCode, rule.CountryCode)
	setInt64Pointer(&m.CityCode, rule.CityCode)
	setStringPointer(&m.CountryKey, rule.CountryKey)
	setStringPointer(&m.CityKey, rule.CityKey)
	setInt64Pointer(&m.IsRegion, rule.IsRegion)
	setStringPointer(&m.CloudCode, rule.CloudCode)
	setInt64Pointer(&m.IsCloud, rule.IsCloud)
	setStringPointer(&m.ParamTemplateName, rule.ParamTemplateName)
	setInt64Pointer(&m.Invalid, rule.Invalid)
	setStringPointer(&m.CreateTime, rule.CreateTime)
	setStringPointer(&m.UpdateTime, rule.UpdateTime)
	setInt64Pointer(&m.DnsParseCnt, rule.DnsParseCnt)
	setInt64Pointer(&m.HitCount, rule.HitCount)
	setStringPointer(&m.LastHitTime, rule.LastHitTime)

	if rule.BelongMember != nil {
		member := fwmNatAclRuleBelongMemberModel{
			AppId:    types.StringPointerValue(rule.BelongMember.AppId),
			Uin:      types.StringPointerValue(rule.BelongMember.Uin),
			Nickname: types.StringPointerValue(rule.BelongMember.Nickname),
			MemberId: types.StringPointerValue(rule.BelongMember.MemberId),
			NodeName: types.StringPointerValue(rule.BelongMember.NodeName),
		}
		obj, diags := types.ObjectValueFrom(ctx, fwmNatAclRuleBelongMemberAttrTypes(), member)
		if !diags.HasError() {
			m.BelongMember = obj
		}
	} else if m.BelongMember.IsNull() {
		// Keep a typed null object so the value always matches the schema even
		// when the cloud side returns no member information.
		m.BelongMember = types.ObjectNull(fwmNatAclRuleBelongMemberAttrTypes())
	}
}

// toRuleParam converts the plan/state into the cloud API rule body.
func (m *FwmNatAclRuleResourceModel) toRuleParam() *fwmv20250611.NatAclRule {
	rule := &fwmv20250611.NatAclRule{
		SourceContent: m.SourceContent.ValueStringPointer(),
		SourceType:    m.SourceType.ValueStringPointer(),
		TargetContent: m.TargetContent.ValueStringPointer(),
		TargetType:    m.TargetType.ValueStringPointer(),
		Protocol:      m.Protocol.ValueStringPointer(),
		RuleAction:    m.RuleAction.ValueStringPointer(),
		Scope:         m.Scope.ValueStringPointer(),
	}

	if !m.OrderIndex.IsNull() && !m.OrderIndex.IsUnknown() {
		rule.OrderIndex = m.OrderIndex.ValueInt64Pointer()
	}
	if !m.Direction.IsNull() && !m.Direction.IsUnknown() {
		rule.Direction = m.Direction.ValueInt64Pointer()
	}
	if !m.RuleId.IsNull() && !m.RuleId.IsUnknown() {
		rule.RuleId = m.RuleId.ValueStringPointer()
	}
	if !m.Port.IsNull() && !m.Port.IsUnknown() {
		rule.Port = m.Port.ValueStringPointer()
	}
	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		rule.Description = m.Description.ValueStringPointer()
	}
	if !m.ParamTemplateId.IsNull() && !m.ParamTemplateId.IsUnknown() {
		rule.ParamTemplateId = m.ParamTemplateId.ValueStringPointer()
	}
	if !m.BelongMemberId.IsNull() && !m.BelongMemberId.IsUnknown() {
		rule.BelongMemberId = m.BelongMemberId.ValueStringPointer()
	}

	return rule
}

func fwmNatAclRuleBelongMemberAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"app_id":    types.StringType,
		"uin":       types.StringType,
		"nickname":  types.StringType,
		"member_id": types.StringType,
		"node_name": types.StringType,
	}
}

// buildFwmNatAclRuleId joins the composite resource id with tccommon.FILED_SP.
func buildFwmNatAclRuleId(groupId, ruleId string) string {
	return strings.Join([]string{groupId, ruleId}, tccommon.FILED_SP)
}

// splitFwmNatAclRuleId parses `<group_id>#<rule_id>`.
func splitFwmNatAclRuleId(id string) (string, string, bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func stringListValue(values []*string) types.List {
	out := make([]attr.Value, 0, len(values))
	for _, v := range values {
		out = append(out, types.StringPointerValue(v))
	}
	return types.ListValueMust(types.StringType, out)
}

func stringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func setStringPointer(target *types.String, source *string) {
	*target = types.StringPointerValue(source)
}

func setInt64Pointer(target *types.Int64, source *int64) {
	*target = types.Int64PointerValue(source)
}

// setStringIfNotEmpty only overwrites the model when the cloud side returns a
// non-empty value, so user inputs are kept when the API omits the field.
func setStringIfNotEmpty(target *types.String, source *string) {
	if source == nil || *source == "" {
		return
	}
	*target = types.StringValue(*source)
}

// setInt64IfNotNil only overwrites the model when the cloud side returns a value.
func setInt64IfNotNil(target *types.Int64, source *int64) {
	if source == nil {
		return
	}
	*target = types.Int64Value(*source)
}

// stringRequiresReplace forces resource replacement when a string attribute changes.
type stringRequiresReplace struct{}

func (stringRequiresReplace) Description(_ context.Context) string {
	return "If the value of this attribute changes, Terraform will destroy and recreate the resource."
}

func (stringRequiresReplace) MarkdownDescription(ctx context.Context) string {
	return stringRequiresReplace{}.Description(ctx)
}

func (stringRequiresReplace) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
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

// stringUseStateForUnknown copies the prior state value into the plan when the
// new plan value is unknown, so computed attributes stay stable across refreshes.
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
