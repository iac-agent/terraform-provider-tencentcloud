package config

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwhelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// Compile-time interface assertions (repo convention: only assert the
// extension interfaces; the base interface is checked implicitly by the
// factory return value).
var (
	_ resource.ResourceWithConfigure   = &ConfigAggregateConfigRuleResource{}
	_ resource.ResourceWithImportState = &ConfigAggregateConfigRuleResource{}

	// errAggregateConfigRuleNotFound marks an empty Describe response. The
	// Config API answers with an empty body once the rule is gone, so the case
	// must be told apart from a transport or API failure.
	errAggregateConfigRuleNotFound = errors.New("aggregate config rule not found")
)

// NewConfigAggregateConfigRuleResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewConfigAggregateConfigRuleResource() resource.Resource {
	return &ConfigAggregateConfigRuleResource{}
}

// ConfigAggregateConfigRuleResource implements resource.Resource for
// tencentcloud_config_aggregate_config_rule.
//
// The cloud side addresses an account group rule with an account group id plus
// a rule id, so the Terraform id is both parts joined with `#`
// (tccommon.FILED_SP) and split again for every CRUD call.
type ConfigAggregateConfigRuleResource struct {
	fw.ResourceWithConfigure
}

// ConfigAggregateConfigRuleResourceModel maps the top-level schema
// attributes. Every schema attribute appears here exactly once.
type ConfigAggregateConfigRuleResourceModel struct {
	Id                       types.String `tfsdk:"id"`
	Identifier               types.String `tfsdk:"identifier"`
	IdentifierType           types.String `tfsdk:"identifier_type"`
	RuleName                 types.String `tfsdk:"rule_name"`
	ResourceType             types.List   `tfsdk:"resource_type"`
	TriggerType              types.List   `tfsdk:"trigger_type"`
	RiskLevel                types.Int64  `tfsdk:"risk_level"`
	AccountGroupId           types.String `tfsdk:"account_group_id"`
	InputParameter           types.List   `tfsdk:"input_parameter"`
	Description              types.String `tfsdk:"description"`
	RegionScope              types.List   `tfsdk:"region_scope"`
	TagsScope                types.List   `tfsdk:"tags_scope"`
	ExcludeResourceIdsScope  types.List   `tfsdk:"exclude_resource_ids_scope"`
	ConfigRuleId             types.String `tfsdk:"config_rule_id"`
	Status                   types.String `tfsdk:"status"`
	ComplianceResult         types.String `tfsdk:"compliance_result"`
	CreateTime               types.String `tfsdk:"create_time"`
	ServiceFunction          types.String `tfsdk:"service_function"`
	CompliancePackId         types.String `tfsdk:"compliance_pack_id"`
	CompliancePackName       types.String `tfsdk:"compliance_pack_name"`
	AccountGroupName         types.String `tfsdk:"account_group_name"`
	RuleOwnerId              types.Int64  `tfsdk:"rule_owner_id"`
	Labels                   types.List   `tfsdk:"labels"`
	RegionsScope             types.List   `tfsdk:"regions_scope"`
	ManageTriggerType        types.List   `tfsdk:"manage_trigger_type"`
	SourceCondition          types.List   `tfsdk:"source_condition"`
	ManageInputParameter     types.List   `tfsdk:"manage_input_parameter"`
	ConfigRuleInvokedTime    types.String `tfsdk:"config_rule_invoked_time"`
}

// aggregateConfigRuleTriggerTypeModel maps one element of trigger_type.
type aggregateConfigRuleTriggerTypeModel struct {
	MessageType               types.String `tfsdk:"message_type"`
	MaximumExecutionFrequency types.String `tfsdk:"maximum_execution_frequency"`
}

// aggregateConfigRuleInputParameterModel maps one element of input_parameter.
type aggregateConfigRuleInputParameterModel struct {
	ParameterKey types.String `tfsdk:"parameter_key"`
	Type         types.String `tfsdk:"type"`
	Value        types.String `tfsdk:"value"`
}

// aggregateConfigRuleTagsScopeModel maps one element of tags_scope.
type aggregateConfigRuleTagsScopeModel struct {
	TagKey   types.String `tfsdk:"tag_key"`
	TagValue types.String `tfsdk:"tag_value"`
}

// aggregateConfigRuleSourceConditionModel maps one element of the read-only
// source_condition list.
type aggregateConfigRuleSourceConditionModel struct {
	EmptyAs      types.String `tfsdk:"empty_as"`
	SelectPath   types.String `tfsdk:"select_path"`
	Operator     types.String `tfsdk:"operator"`
	Required     types.Bool   `tfsdk:"required"`
	DesiredValue types.String `tfsdk:"desired_value"`
}

// aggregateConfigRuleManageInputParameterModel maps one element of the
// read-only manage_input_parameter list.
type aggregateConfigRuleManageInputParameterModel struct {
	ValueType    types.String `tfsdk:"value_type"`
	ParameterKey types.String `tfsdk:"parameter_key"`
	Type         types.String `tfsdk:"type"`
	DefaultValue types.String `tfsdk:"default_value"`
	Description  types.String `tfsdk:"description"`
}

func (r *ConfigAggregateConfigRuleResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_config_aggregate_config_rule"
}

func (r *ConfigAggregateConfigRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create an aggregate (account group) config rule of Config.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier in the form `<account_group_id>#<rule_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"identifier": schema.StringAttribute{
				Required: true,
				Description: "Rule template identifier. For system preset rules use the identifier name; " +
					"for custom rules use the cloud function ARN (region:functionName). " +
					"Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"identifier_type": schema.StringAttribute{
				Required: true,
				Description: "Rule template type. Valid values: `SYSTEM` (system preset), `CUSTOMIZE` (custom). " +
					"Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rule_name": schema.StringAttribute{
				Required:    true,
				Description: "Rule name.",
			},
			"resource_type": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "Supported resource type list (e.g. QCS::CAM::User).",
			},
			"trigger_type": schema.ListNestedAttribute{
				Required:    true,
				Description: "Trigger type list, up to 2 entries.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"message_type": schema.StringAttribute{
							Required:    true,
							Description: "Trigger message type. Valid values: `ScheduledNotification`, `ConfigurationItemChangeNotification`.",
						},
						"maximum_execution_frequency": schema.StringAttribute{
							Optional:    true,
							Description: "Maximum execution frequency (only for `ScheduledNotification`). e.g. `TwentyFour_Hours`.",
						},
					},
				},
			},
			"risk_level": schema.Int64Attribute{
				Required:    true,
				Description: "Risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).",
			},
			"account_group_id": schema.StringAttribute{
				Required:    true,
				Description: "Account group id. Immutable: changing it forces resource replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"input_parameter": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Rule input parameter list.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"parameter_key": schema.StringAttribute{
							Required:    true,
							Description: "Parameter key.",
						},
						"type": schema.StringAttribute{
							Optional:    true,
							Description: "Parameter type. Required: `Require`, optional: `Optional`.",
						},
						"value": schema.StringAttribute{
							Optional:    true,
							Description: "Parameter value.",
						},
					},
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Rule description.",
			},
			"region_scope": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Region scope list; rule only applies to resources in the specified regions.",
			},
			"tags_scope": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Tag scope list; rule only applies to resources with the specified tags.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tag_key": schema.StringAttribute{
							Optional:    true,
							Description: "Tag key.",
						},
						"tag_value": schema.StringAttribute{
							Optional:    true,
							Description: "Tag value.",
						},
					},
				},
			},
			"exclude_resource_ids_scope": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Resource ID list excluded from rule evaluation.",
			},

			// Computed
			"config_rule_id": schema.StringAttribute{
				Computed:    true,
				Description: "Config rule id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Rule status. Valid values: `ACTIVE` (enabled), `NO_ACTIVE` (disabled).",
			},
			"compliance_result": schema.StringAttribute{
				Computed:    true,
				Description: "Compliance result. Valid values: `COMPLIANT`, `NON_COMPLIANT`, `NOT_APPLICABLE`.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time. Format: YYYY-MM-DD h:i:s.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_function": schema.StringAttribute{
				Computed:    true,
				Description: "Cloud function corresponding to the rule.",
			},
			"compliance_pack_id": schema.StringAttribute{
				Computed:    true,
				Description: "Compliance pack id the rule belongs to (empty when the rule is standalone).",
			},
			"compliance_pack_name": schema.StringAttribute{
				Computed:    true,
				Description: "Compliance pack name the rule belongs to (empty when the rule is standalone).",
			},
			"account_group_name": schema.StringAttribute{
				Computed:    true,
				Description: "Account group name.",
			},
			"rule_owner_id": schema.Int64Attribute{
				Computed:    true,
				Description: "User id the rule belongs to.",
			},
			"labels": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Labels the rule belongs to.",
			},
			"regions_scope": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Region scope reported by the cloud.",
			},
			"manage_trigger_type": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Trigger types supported by the preset rule. Valid values: `ScheduledNotification`, `ConfigurationItemChangeNotification`.",
			},
			"source_condition": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Trigger conditions of the preset rule.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"empty_as": schema.StringAttribute{
							Computed:    true,
							Description: "Behavior when the condition is empty: `COMPLIANT`, `NON_COMPLIANT`, `NOT_APPLICABLE`.",
						},
						"select_path": schema.StringAttribute{
							Computed:    true,
							Description: "Configuration path.",
						},
						"operator": schema.StringAttribute{
							Computed:    true,
							Description: "Comparison operator.",
						},
						"required": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the condition is required.",
						},
						"desired_value": schema.StringAttribute{
							Computed:    true,
							Description: "Desired value.",
						},
					},
				},
			},
			"manage_input_parameter": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Managed (preset) rule parameter details.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"value_type": schema.StringAttribute{
							Computed:    true,
							Description: "Value type. Integer for numeric values, String for strings.",
						},
						"parameter_key": schema.StringAttribute{
							Computed:    true,
							Description: "Parameter key.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Parameter type. Required: `Require`, optional: `Optional`.",
						},
						"default_value": schema.StringAttribute{
							Computed:    true,
							Description: "Default value.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Parameter description.",
						},
					},
				},
			},
			"config_rule_invoked_time": schema.StringAttribute{
				Computed:    true,
				Description: "Last rule evaluation time. Format: YYYY-MM-DD h:i:s.",
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *ConfigAggregateConfigRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ConfigAggregateConfigRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	accountGroupId := plan.AccountGroupId.ValueString()

	request := configv20220802.NewAddAggregateConfigRuleRequest()
	request.AccountGroupId = helper.String(accountGroupId)
	request.Identifier = fwhelper.StringPointerFromValue(plan.Identifier)
	request.IdentifierType = fwhelper.StringPointerFromValue(plan.IdentifierType)
	request.RuleName = fwhelper.StringPointerFromValue(plan.RuleName)
	request.RiskLevel = fwhelper.Uint64PointerFromValue(plan.RiskLevel)
	request.Description = fwhelper.StringPointerFromValue(plan.Description)

	resourceTypes, resourceTypeDiags := expandAggregateConfigRuleStrings(ctx, plan.ResourceType)
	resp.Diagnostics.Append(resourceTypeDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.ResourceType = resourceTypes

	triggerTypes, triggerTypeDiags := expandAggregateConfigRuleTriggerTypes(ctx, plan.TriggerType)
	resp.Diagnostics.Append(triggerTypeDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.TriggerType = triggerTypes

	inputParameters, inputParameterDiags := expandAggregateConfigRuleInputParameters(ctx, plan.InputParameter)
	resp.Diagnostics.Append(inputParameterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.InputParameter = inputParameters

	regionScope, regionScopeDiags := expandAggregateConfigRuleStrings(ctx, plan.RegionScope)
	resp.Diagnostics.Append(regionScopeDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.RegionScope = regionScope

	tagsScope, tagsScopeDiags := expandAggregateConfigRuleTagsScope(ctx, plan.TagsScope)
	resp.Diagnostics.Append(tagsScopeDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.TagsScope = tagsScope

	excludeResourceIdsScope, excludeScopeDiags := expandAggregateConfigRuleStrings(ctx, plan.ExcludeResourceIdsScope)
	resp.Diagnostics.Append(excludeScopeDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.ExcludeResourceIdsScope = excludeResourceIdsScope

	ruleId, err := service.AddAggregateConfigRule(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Config aggregate config rule",
			fmt.Sprintf("AddAggregateConfigRule failed (%s): %s", aggregateConfigRuleScope(accountGroupId, ""), err),
		)
		return
	}

	plan.ConfigRuleId = types.StringValue(ruleId)
	plan.Id = types.StringValue(aggregateConfigRuleId(accountGroupId, ruleId))

	// Read the rule back so the very first state already holds the cloud-owned
	// attributes (status, create_time, computed details).
	item, err := describeAggregateConfigRule(ctx, service, accountGroupId, ruleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Config aggregate config rule after create",
			fmt.Sprintf(
				"AddAggregateConfigRule succeeded but the rule cannot be read back (%s). "+
					"Import it with `terraform import tencentcloud_config_aggregate_config_rule.<name> %s` to keep it managed: %s",
				aggregateConfigRuleScope(accountGroupId, ruleId), plan.Id.ValueString(), err,
			),
		)
		return
	}

	resp.Diagnostics.Append(refreshAggregateConfigRuleModel(ctx, &plan, accountGroupId, ruleId, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *ConfigAggregateConfigRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ConfigAggregateConfigRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountGroupId, ruleId, ok := parseAggregateConfigRuleId(&resp.Diagnostics, state.Id.ValueString())
	if !ok {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	item, err := describeAggregateConfigRule(ctx, service, accountGroupId, ruleId)
	if err != nil {
		// A confirmed "not found" means the rule was removed outside Terraform,
		// so it is dropped from state instead of failing the refresh.
		if errors.Is(err, errAggregateConfigRuleNotFound) {
			resp.Diagnostics.AddWarning(
				"Resource not found",
				fmt.Sprintf("Config aggregate config rule %s no longer exists, removing it from state.", state.Id.ValueString()),
			)
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading Config aggregate config rule",
			fmt.Sprintf("DescribeAggregateConfigRule failed (%s): %s", aggregateConfigRuleScope(accountGroupId, ruleId), err),
		)
		return
	}

	resp.Diagnostics.Append(refreshAggregateConfigRuleModel(ctx, &state, accountGroupId, ruleId, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates.
func (r *ConfigAggregateConfigRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ConfigAggregateConfigRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ConfigAggregateConfigRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountGroupId, ruleId, ok := parseAggregateConfigRuleId(&resp.Diagnostics, state.Id.ValueString())
	if !ok {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	// UpdateAggregateConfigRule requires the identity and the mandatory
	// trigger/risk fields; the remaining fields are sent when configured.
	request := configv20220802.NewUpdateAggregateConfigRuleRequest()
	request.AccountGroupId = helper.String(accountGroupId)
	request.RuleId = helper.String(ruleId)
	request.RiskLevel = fwhelper.Uint64PointerFromValue(plan.RiskLevel)
	request.RuleName = fwhelper.StringPointerFromValue(plan.RuleName)
	request.Description = fwhelper.StringPointerFromValue(plan.Description)

	triggerTypes, triggerTypeDiags := expandAggregateConfigRuleTriggerTypes(ctx, plan.TriggerType)
	resp.Diagnostics.Append(triggerTypeDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.TriggerType = triggerTypes

	inputParameters, inputParameterDiags := expandAggregateConfigRuleInputParameters(ctx, plan.InputParameter)
	resp.Diagnostics.Append(inputParameterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.InputParameter = inputParameters

	regionScope, regionScopeDiags := expandAggregateConfigRuleStrings(ctx, plan.RegionScope)
	resp.Diagnostics.Append(regionScopeDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.RegionScope = regionScope

	tagsScope, tagsScopeDiags := expandAggregateConfigRuleTagsScope(ctx, plan.TagsScope)
	resp.Diagnostics.Append(tagsScopeDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.TagsScope = tagsScope

	excludeResourceIdsScope, excludeScopeDiags := expandAggregateConfigRuleStrings(ctx, plan.ExcludeResourceIdsScope)
	resp.Diagnostics.Append(excludeScopeDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request.ExcludeResourceIdsScope = excludeResourceIdsScope

	if err := service.UpdateAggregateConfigRule(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating Config aggregate config rule",
			fmt.Sprintf("UpdateAggregateConfigRule failed (%s): %s", aggregateConfigRuleScope(accountGroupId, ruleId), err),
		)
		return
	}

	item, err := describeAggregateConfigRule(ctx, service, accountGroupId, ruleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Config aggregate config rule after update",
			fmt.Sprintf("UpdateAggregateConfigRule succeeded but the rule cannot be read back (%s): %s", aggregateConfigRuleScope(accountGroupId, ruleId), err),
		)
		return
	}

	resp.Diagnostics.Append(refreshAggregateConfigRuleModel(ctx, &plan, accountGroupId, ruleId, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *ConfigAggregateConfigRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ConfigAggregateConfigRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountGroupId, ruleId, ok := parseAggregateConfigRuleId(&resp.Diagnostics, state.Id.ValueString())
	if !ok {
		return
	}

	service, ok := r.service(&resp.Diagnostics)
	if !ok {
		return
	}

	request := configv20220802.NewDeleteAggregateConfigRuleRequest()
	request.AccountGroupId = helper.String(accountGroupId)
	request.RuleId = helper.String(ruleId)

	if err := service.DeleteAggregateConfigRule(ctx, request); err != nil {
		// An already deleted rule is a successful deletion.
		if tccommon.IsExpectError(err, []string{"ResourceNotFound", "ResourceInexistent"}) {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting Config aggregate config rule",
			fmt.Sprintf("DeleteAggregateConfigRule failed (%s): %s", aggregateConfigRuleScope(accountGroupId, ruleId), err),
		)
	}
}

// ImportState supports `terraform import` using the
// `<account_group_id>#<rule_id>` form. The remaining attributes are filled in
// by the Read that follows.
func (r *ConfigAggregateConfigRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _, _, ok := splitAggregateConfigRuleId(req.ID); !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Import id %q is invalid, expected `<account_group_id>%s<rule_id>`.", req.ID, tccommon.FILED_SP),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// service returns the Config service bound to the shared client. When the
// provider has not been configured it records a diagnostic and reports false,
// so the CRUD methods never dereference a nil SDK client.
func (r *ConfigAggregateConfigRuleResource) service(diags *diag.Diagnostics) (ConfigService, bool) {
	if r.Client() == nil {
		diags.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return ConfigService{}, false
	}
	return NewConfigService(r.Client()), true
}

// describeAggregateConfigRule reads the rule from the cloud. An empty response
// means the rule no longer exists and is reported as
// errAggregateConfigRuleNotFound; API and transport failures are returned
// verbatim.
func describeAggregateConfigRule(ctx context.Context, service ConfigService, accountGroupId, ruleId string) (*configv20220802.ConfigRule, error) {
	item, err := service.DescribeAggregateConfigRuleById(ctx, accountGroupId, ruleId)
	if err != nil {
		return nil, err
	}
	if item == nil || item.ConfigRuleId == nil {
		return nil, errAggregateConfigRuleNotFound
	}
	return item, nil
}

// aggregateConfigRuleId builds the Terraform id from its two cloud parts.
func aggregateConfigRuleId(accountGroupId, ruleId string) string {
	return strings.Join([]string{accountGroupId, ruleId}, tccommon.FILED_SP)
}

// splitAggregateConfigRuleId parses `<account_group_id>#<rule_id>`.
func splitAggregateConfigRuleId(id string) (string, string, bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// parseAggregateConfigRuleId resolves the Terraform id into the pair of ids
// the Config API expects, recording a diagnostic when the id is malformed.
func parseAggregateConfigRuleId(diags *diag.Diagnostics, id string) (accountGroupId, ruleId string, ok bool) {
	accountGroupId, ruleId, ok = splitAggregateConfigRuleId(id)
	if !ok {
		diags.AddError(
			"Invalid resource ID",
			fmt.Sprintf("Resource id %q is invalid, expected `<account_group_id>%s<rule_id>`.", id, tccommon.FILED_SP),
		)
	}
	return
}

// aggregateConfigRuleScope renders the cloud-side identity of a rule for
// diagnostics, so every failure can be traced back to a concrete account group
// and rule id.
func aggregateConfigRuleScope(accountGroupId, ruleId string) string {
	if ruleId == "" {
		return fmt.Sprintf("account_group_id=%q", accountGroupId)
	}
	return fmt.Sprintf("account_group_id=%q, rule_id=%q", accountGroupId, ruleId)
}

// refreshAggregateConfigRuleModel copies the cloud state into the model.
// Attributes the Describe API may omit keep their plan/state value, so a value
// configured by the practitioner is never nulled out after apply.
func refreshAggregateConfigRuleModel(ctx context.Context, m *ConfigAggregateConfigRuleResourceModel, accountGroupId, ruleId string, item *configv20220802.ConfigRule) diag.Diagnostics {
	var diags diag.Diagnostics

	if item == nil {
		diags.AddError(
			"Malformed Config aggregate config rule response",
			fmt.Sprintf("DescribeAggregateConfigRule returned an empty response (%s).", aggregateConfigRuleScope(accountGroupId, ruleId)),
		)
		return diags
	}

	m.AccountGroupId = types.StringValue(accountGroupId)
	m.ConfigRuleId = types.StringValue(ruleId)
	m.Id = types.StringValue(aggregateConfigRuleId(accountGroupId, ruleId))

	m.Identifier = fwhelper.StringValueOrPrev(item.Identifier, m.Identifier)
	m.IdentifierType = fwhelper.StringValueOrPrev(item.IdentifierType, m.IdentifierType)
	m.RuleName = fwhelper.StringValueOrPrev(item.RuleName, m.RuleName)
	m.Description = fwhelper.StringValueOrPrev(item.Description, m.Description)
	m.Status = types.StringPointerValue(item.Status)
	m.ComplianceResult = types.StringPointerValue(item.ComplianceResult)
	m.CreateTime = types.StringPointerValue(item.CreateTime)
	m.ServiceFunction = types.StringPointerValue(item.ServiceFunction)
	m.CompliancePackId = types.StringPointerValue(item.CompliancePackId)
	m.CompliancePackName = types.StringPointerValue(item.CompliancePackName)
	m.AccountGroupName = types.StringPointerValue(item.AccountGroupName)
	m.ConfigRuleInvokedTime = types.StringPointerValue(item.ConfigRuleInvokedTime)

	// ConfigRule.RiskLevel is *int64.
	if item.RiskLevel != nil {
		m.RiskLevel = types.Int64Value(*item.RiskLevel)
	} else if m.RiskLevel.IsUnknown() {
		m.RiskLevel = types.Int64Null()
	}

	if item.RuleOwnerId != nil {
		m.RuleOwnerId = types.Int64Value(int64(*item.RuleOwnerId))
	} else {
		m.RuleOwnerId = types.Int64Null()
	}

	resourceTypes, resourceTypeDiags := flattenAggregateConfigRuleStrings(ctx, item.ResourceType, m.ResourceType)
	diags.Append(resourceTypeDiags...)
	m.ResourceType = resourceTypes

	triggerTypes, triggerTypeDiags := flattenAggregateConfigRuleTriggerTypes(ctx, item.TriggerType, m.TriggerType)
	diags.Append(triggerTypeDiags...)
	m.TriggerType = triggerTypes

	inputParameters, inputParameterDiags := flattenAggregateConfigRuleInputParameters(ctx, item.InputParameter, m.InputParameter)
	diags.Append(inputParameterDiags...)
	m.InputParameter = inputParameters

	regionScope, regionScopeDiags := flattenAggregateConfigRuleStrings(ctx, item.RegionsScope, m.RegionScope)
	diags.Append(regionScopeDiags...)
	m.RegionScope = regionScope

	tagsScope, tagsScopeDiags := flattenAggregateConfigRuleTagsScope(ctx, item.TagsScope, m.TagsScope)
	diags.Append(tagsScopeDiags...)
	m.TagsScope = tagsScope

	excludeResourceIdsScope, excludeScopeDiags := flattenAggregateConfigRuleStrings(ctx, item.ExcludeResourceIdsScope, m.ExcludeResourceIdsScope)
	diags.Append(excludeScopeDiags...)
	m.ExcludeResourceIdsScope = excludeResourceIdsScope

	labels, labelsDiags := flattenAggregateConfigRuleStrings(ctx, item.Labels, types.ListNull(types.StringType))
	diags.Append(labelsDiags...)
	m.Labels = labels

	manageTriggerTypes, manageTriggerTypeDiags := flattenAggregateConfigRuleStrings(ctx, item.ManageTriggerType, types.ListNull(types.StringType))
	diags.Append(manageTriggerTypeDiags...)
	m.ManageTriggerType = manageTriggerTypes

	sourceConditions, sourceConditionDiags := flattenAggregateConfigRuleSourceConditions(ctx, item.SourceCondition)
	diags.Append(sourceConditionDiags...)
	m.SourceCondition = sourceConditions

	manageInputParameters, manageInputParameterDiags := flattenAggregateConfigRuleManageInputParameters(ctx, item.ManageInputParameter)
	diags.Append(manageInputParameterDiags...)
	m.ManageInputParameter = manageInputParameters

	// Required attributes must never be written as null, otherwise Terraform
	// rejects the refreshed object.
	if !diags.HasError() && (m.RuleName.IsNull() || m.ResourceType.IsNull() || m.TriggerType.IsNull() || m.RiskLevel.IsNull()) {
		diags.AddError(
			"Malformed Config aggregate config rule response",
			fmt.Sprintf("DescribeAggregateConfigRule returned no rule_name, resource_type, trigger_type or risk_level (%s).", aggregateConfigRuleScope(accountGroupId, ruleId)),
		)
	}
	return diags
}

// expandAggregateConfigRuleStrings converts a planned string list into the SDK
// string slice.
func expandAggregateConfigRuleStrings(ctx context.Context, list types.List) ([]*string, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var values []string
	diags.Append(list.ElementsAs(ctx, &values, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*string, 0, len(values))
	for _, value := range values {
		out = append(out, helper.String(value))
	}
	return out, diags
}

// flattenAggregateConfigRuleStrings builds a string list from the SDK string
// slice. The previous plan/state list is used to keep configured values when
// the API omits them.
func flattenAggregateConfigRuleStrings(ctx context.Context, values []*string, prevList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	out := make([]string, 0, len(values))
	for _, value := range values {
		if value == nil {
			continue // a malformed response must never panic the provider
		}
		out = append(out, *value)
	}

	// An empty response must not clear a list that was planned with entries.
	if len(out) == 0 && !prevList.IsNull() && !prevList.IsUnknown() {
		return prevList, diags
	}

	list, listDiags := types.ListValueFrom(ctx, types.StringType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(types.StringType), diags
	}
	return list, diags
}

// expandAggregateConfigRuleTriggerTypes converts the planned trigger_type list
// into the SDK TriggerType slice.
func expandAggregateConfigRuleTriggerTypes(ctx context.Context, list types.List) ([]*configv20220802.TriggerType, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var entries []aggregateConfigRuleTriggerTypeModel
	diags.Append(list.ElementsAs(ctx, &entries, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*configv20220802.TriggerType, 0, len(entries))
	for _, entry := range entries {
		out = append(out, &configv20220802.TriggerType{
			MessageType:               fwhelper.StringPointerFromValue(entry.MessageType),
			MaximumExecutionFrequency: fwhelper.StringPointerFromValue(entry.MaximumExecutionFrequency),
		})
	}
	return out, diags
}

// flattenAggregateConfigRuleTriggerTypes builds the trigger_type list from the
// Describe response. The previous list is used, position by position, to keep
// values the API omits.
func flattenAggregateConfigRuleTriggerTypes(ctx context.Context, entries []*configv20220802.TriggerType, prevList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := triggerTypeObjectType()

	prev := make([]aggregateConfigRuleTriggerTypeModel, 0)
	if !prevList.IsNull() && !prevList.IsUnknown() {
		diags.Append(prevList.ElementsAs(ctx, &prev, false)...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}
	}

	out := make([]aggregateConfigRuleTriggerTypeModel, 0, len(entries))
	for _, entry := range entries {
		if entry == nil {
			continue // a malformed response must never panic the provider
		}

		// The API keeps trigger order stable, so the entry at the same
		// position holds the previously configured values.
		var prior aggregateConfigRuleTriggerTypeModel
		if len(out) < len(prev) {
			prior = prev[len(out)]
		}

		out = append(out, aggregateConfigRuleTriggerTypeModel{
			MessageType:               fwhelper.StringValueOrPrev(entry.MessageType, prior.MessageType),
			MaximumExecutionFrequency: fwhelper.StringValueOrPrev(entry.MaximumExecutionFrequency, prior.MaximumExecutionFrequency),
		})
	}

	// A Required list must stay consistent with what was planned, so an empty
	// response keeps the previous value instead of clearing configured
	// triggers.
	if len(out) == 0 && !prevList.IsNull() && !prevList.IsUnknown() {
		return prevList, diags
	}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

// expandAggregateConfigRuleInputParameters converts the planned
// input_parameter list into the SDK InputParameter slice.
func expandAggregateConfigRuleInputParameters(ctx context.Context, list types.List) ([]*configv20220802.InputParameter, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var params []aggregateConfigRuleInputParameterModel
	diags.Append(list.ElementsAs(ctx, &params, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*configv20220802.InputParameter, 0, len(params))
	for _, param := range params {
		out = append(out, &configv20220802.InputParameter{
			ParameterKey: fwhelper.StringPointerFromValue(param.ParameterKey),
			Type:         fwhelper.StringPointerFromValue(param.Type),
			Value:        fwhelper.StringPointerFromValue(param.Value),
		})
	}
	return out, diags
}

// flattenAggregateConfigRuleInputParameters builds the input_parameter list
// from the Describe response. The previous list is used, position by position,
// to keep values the API omits.
func flattenAggregateConfigRuleInputParameters(ctx context.Context, params []*configv20220802.InputParameter, prevList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := aggregateInputParameterObjectType()

	prev := make([]aggregateConfigRuleInputParameterModel, 0)
	if !prevList.IsNull() && !prevList.IsUnknown() {
		diags.Append(prevList.ElementsAs(ctx, &prev, false)...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}
	}

	out := make([]aggregateConfigRuleInputParameterModel, 0, len(params))
	for _, param := range params {
		if param == nil {
			continue // a malformed response must never panic the provider
		}

		// The API keeps parameter order stable, so the entry at the same
		// position holds the previously configured values.
		var prior aggregateConfigRuleInputParameterModel
		if len(out) < len(prev) {
			prior = prev[len(out)]
		}

		out = append(out, aggregateConfigRuleInputParameterModel{
			ParameterKey: fwhelper.StringValueOrPrev(param.ParameterKey, prior.ParameterKey),
			Type:         fwhelper.StringValueOrPrev(param.Type, prior.Type),
			Value:        fwhelper.StringValueOrPrev(param.Value, prior.Value),
		})
	}

	// An empty response must not clear a list that was planned with entries.
	if len(out) == 0 && !prevList.IsNull() && !prevList.IsUnknown() {
		return prevList, diags
	}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

// expandAggregateConfigRuleTagsScope converts the planned tags_scope list into
// the SDK Tag slice.
func expandAggregateConfigRuleTagsScope(ctx context.Context, list types.List) ([]*configv20220802.Tag, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var tags []aggregateConfigRuleTagsScopeModel
	diags.Append(list.ElementsAs(ctx, &tags, false)...)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]*configv20220802.Tag, 0, len(tags))
	for _, tag := range tags {
		out = append(out, &configv20220802.Tag{
			TagKey:   fwhelper.StringPointerFromValue(tag.TagKey),
			TagValue: fwhelper.StringPointerFromValue(tag.TagValue),
		})
	}
	return out, diags
}

// flattenAggregateConfigRuleTagsScope builds the tags_scope list from the
// Describe response. The previous list is used, position by position, to keep
// values the API omits.
func flattenAggregateConfigRuleTagsScope(ctx context.Context, tags []*configv20220802.Tag, prevList types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	objType := tagsScopeObjectType()

	prev := make([]aggregateConfigRuleTagsScopeModel, 0)
	if !prevList.IsNull() && !prevList.IsUnknown() {
		diags.Append(prevList.ElementsAs(ctx, &prev, false)...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}
	}

	out := make([]aggregateConfigRuleTagsScopeModel, 0, len(tags))
	for _, tag := range tags {
		if tag == nil {
			continue // a malformed response must never panic the provider
		}

		var prior aggregateConfigRuleTagsScopeModel
		if len(out) < len(prev) {
			prior = prev[len(out)]
		}

		out = append(out, aggregateConfigRuleTagsScopeModel{
			TagKey:   fwhelper.StringValueOrPrev(tag.TagKey, prior.TagKey),
			TagValue: fwhelper.StringValueOrPrev(tag.TagValue, prior.TagValue),
		})
	}

	// An empty response must not clear a list that was planned with entries.
	if len(out) == 0 && !prevList.IsNull() && !prevList.IsUnknown() {
		return prevList, diags
	}

	list, listDiags := types.ListValueFrom(ctx, objType, out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}
	return list, diags
}

// flattenAggregateConfigRuleSourceConditions builds the read-only
// source_condition list from the Describe response.
func flattenAggregateConfigRuleSourceConditions(ctx context.Context, conditions []*configv20220802.SourceConditionForManage) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	out := make([]aggregateConfigRuleSourceConditionModel, 0, len(conditions))
	for _, condition := range conditions {
		if condition == nil {
			continue // a malformed response must never panic the provider
		}

		out = append(out, aggregateConfigRuleSourceConditionModel{
			EmptyAs:      types.StringPointerValue(condition.EmptyAs),
			SelectPath:   types.StringPointerValue(condition.SelectPath),
			Operator:     types.StringPointerValue(condition.Operator),
			Required:     fwhelper.BoolValueOrNull(condition.Required),
			DesiredValue: types.StringPointerValue(condition.DesiredValue),
		})
	}

	list, listDiags := types.ListValueFrom(ctx, sourceConditionObjectType(), out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(sourceConditionObjectType()), diags
	}
	return list, diags
}

// flattenAggregateConfigRuleManageInputParameters builds the read-only
// manage_input_parameter list from the Describe response.
func flattenAggregateConfigRuleManageInputParameters(ctx context.Context, params []*configv20220802.InputParameterForManage) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	out := make([]aggregateConfigRuleManageInputParameterModel, 0, len(params))
	for _, param := range params {
		if param == nil {
			continue // a malformed response must never panic the provider
		}

		out = append(out, aggregateConfigRuleManageInputParameterModel{
			ValueType:    types.StringPointerValue(param.ValueType),
			ParameterKey: types.StringPointerValue(param.ParameterKey),
			Type:         types.StringPointerValue(param.Type),
			DefaultValue: types.StringPointerValue(param.DefaultValue),
			Description:  types.StringPointerValue(param.Description),
		})
	}

	list, listDiags := types.ListValueFrom(ctx, manageInputParameterObjectType(), out)
	diags.Append(listDiags...)
	if diags.HasError() {
		return types.ListNull(manageInputParameterObjectType()), diags
	}
	return list, diags
}

// triggerTypeObjectType is the object type of the trigger_type elements.
func triggerTypeObjectType() attr.Type {
	return types.ObjectType{AttrTypes: triggerTypeAttrTypes()}
}

// triggerTypeAttrTypes is the object type of the trigger_type nested object,
// used to build types.List values.
func triggerTypeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"message_type":                types.StringType,
		"maximum_execution_frequency": types.StringType,
	}
}

// aggregateInputParameterObjectType is the object type of the input_parameter
// elements.
func aggregateInputParameterObjectType() attr.Type {
	return types.ObjectType{AttrTypes: aggregateInputParameterAttrTypes()}
}

// aggregateInputParameterAttrTypes is the object type of the input_parameter
// nested object, used to build types.List values.
func aggregateInputParameterAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"parameter_key": types.StringType,
		"type":          types.StringType,
		"value":         types.StringType,
	}
}

// tagsScopeObjectType is the object type of the tags_scope elements.
func tagsScopeObjectType() attr.Type {
	return types.ObjectType{AttrTypes: tagsScopeAttrTypes()}
}

// tagsScopeAttrTypes is the object type of the tags_scope nested object, used
// to build types.List values.
func tagsScopeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"tag_key":   types.StringType,
		"tag_value": types.StringType,
	}
}

// sourceConditionObjectType is the object type of the source_condition
// elements.
func sourceConditionObjectType() attr.Type {
	return types.ObjectType{AttrTypes: sourceConditionAttrTypes()}
}

// sourceConditionAttrTypes is the object type of the source_condition nested
// object, used to build types.List values.
func sourceConditionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"empty_as":      types.StringType,
		"select_path":   types.StringType,
		"operator":      types.StringType,
		"required":      types.BoolType,
		"desired_value": types.StringType,
	}
}

// manageInputParameterObjectType is the object type of the
// manage_input_parameter elements.
func manageInputParameterObjectType() attr.Type {
	return types.ObjectType{AttrTypes: manageInputParameterAttrTypes()}
}

// manageInputParameterAttrTypes is the object type of the
// manage_input_parameter nested object, used to build types.List values.
func manageInputParameterAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value_type":    types.StringType,
		"parameter_key": types.StringType,
		"type":          types.StringType,
		"default_value": types.StringType,
		"description":   types.StringType,
	}
}
