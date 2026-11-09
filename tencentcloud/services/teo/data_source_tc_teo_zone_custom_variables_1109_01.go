package teo

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	teosdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var _ datasource.DataSourceWithConfigure = &TeoZoneCustomVariablesDataSource{}

// NewTeoZoneCustomVariablesDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewTeoZoneCustomVariablesDataSource() datasource.DataSource {
	return &TeoZoneCustomVariablesDataSource{}
}

// TeoZoneCustomVariablesDataSource implements datasource.DataSource for
// tencentcloud_teo_zone_custom_variables_1109_01.
type TeoZoneCustomVariablesDataSource struct {
	fw.DataSourceWithConfigure
}

// TeoZoneCustomVariablesDataSourceModel maps the schema attributes.
type TeoZoneCustomVariablesDataSourceModel struct {
	ZoneId                  types.String `tfsdk:"zone_id"`
	CustomVariables         types.List   `tfsdk:"custom_variables"`
	CustomVariableOperation types.List   `tfsdk:"custom_variable_operations"`
}

// customVariableItemModel is the element type of custom_variables.
type customVariableItemModel struct {
	Name         types.String `tfsdk:"name"`
	InitialValue types.String `tfsdk:"initial_value"`
	Description  types.String `tfsdk:"description"`
}

// customVariableOperationItemModel is the element type of custom_variable_operations.
type customVariableOperationItemModel struct {
	Branches    types.List `tfsdk:"branches"`
	Description types.List `tfsdk:"description"`
}

// customVariableOperationBranchModel is the element type of branches.
type customVariableOperationBranchModel struct {
	Condition types.String `tfsdk:"condition"`
	Actions   types.List  `tfsdk:"actions"`
	SubRules  types.List  `tfsdk:"sub_rules"`
}

// customVariableOperationActionModel is the element type of actions.
type customVariableOperationActionModel struct {
	Name          types.String                  `tfsdk:"name"`
	SetParameters *customVariableSetParamModel `tfsdk:"set_parameters"`
}

// customVariableSetParamModel is the nested object type of set_parameters.
type customVariableSetParamModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

// customVariableOperationSubRuleModel is the element type of sub_rules.
type customVariableOperationSubRuleModel struct {
	Branches    types.List `tfsdk:"branches"`
	Description types.List `tfsdk:"description"`
}

// deepSubRuleDescriptionModel is the element type of the deepest sub_rules
// list, which only exposes the rule descriptions.
type deepSubRuleDescriptionModel struct {
	Description types.List `tfsdk:"description"`
}

func (d *TeoZoneCustomVariablesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_zone_custom_variables_1109_01"
}

func (d *TeoZoneCustomVariablesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// The schema of the deepest sub_rules list only exposes descriptions,
	// matching the documented one-layer SubRules limitation of the API.
	deepSubRulesAttr := schema.ListNestedAttribute{
		Computed:    true,
		Description: "Deeper sub-rule list. Rules in this list are executed from top to bottom. Only one layer of SubRules is supported; deeper nesting only exposes rule comments.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"description": schema.ListAttribute{
					Computed:    true,
					ElementType: types.StringType,
					Description: "Rule comments.",
				},
			},
		},
	}
	deepBranchActionsAttr := schema.ListNestedAttribute{
		Computed:    true,
		Description: "Actions of the branch. Actions and SubRules cannot both be empty. Valid value of name: `Set`.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Computed:    true,
					Description: "Action name, which needs to correspond to the parameter struct, e.g. when Name=Set, SetParameters is required. Currently only `Set` is supported.",
				},
				"set_parameters": schema.SingleNestedAttribute{
					Computed:    true,
					Description: "Custom variable setting parameters. When there are multiple operations, they are executed in array order. Required when Name is `Set`.",
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Custom variable name. The custom variable must be defined before it can be used in operations.",
						},
						"value": schema.StringAttribute{
							Computed:    true,
							Description: "Custom variable assignment. Supports constant strings, variables and formulas; Chinese characters are not supported. Up to 1000 characters.",
						},
					},
				},
			},
		},
	}

	resp.Schema = schema.Schema{
		Description: "Use this data source to query zone-level custom variables of a TEO zone, " +
			"including variable definitions and variable operation rules.",
		Attributes: map[string]schema.Attribute{
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Site ID.",
			},
			"custom_variables": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of zone-level custom variable definitions.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Variable name. Use the full prefix: user.zone.* for zone-level custom variables, user.rule.* for rule-level custom variables. The custom part after the prefix only supports letters, digits and underscores. Variable names are case-sensitive and up to 50 characters long. The name cannot be modified after creation.",
						},
						"initial_value": schema.StringAttribute{
							Computed:    true,
							Description: "Initial value of the variable. Supports constant strings, variables and formulas, up to 255 characters.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Variable description, up to 60 characters.",
						},
					},
				},
			},
			"custom_variable_operations": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Zone-level custom variable operation rules.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"branches": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Sub-rule branches. This list currently only supports one rule entry; multiple entries are invalid.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"condition": schema.StringAttribute{
										Computed:    true,
										Description: "Match condition expression.",
									},
									"actions": schema.ListNestedAttribute{
										Computed:    true,
										Description: "Actions of the branch. Actions and SubRules cannot both be empty. Valid value of name: `Set`.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"name": schema.StringAttribute{
													Computed:    true,
													Description: "Action name, which needs to correspond to the parameter struct, e.g. when Name=Set, SetParameters is required. Currently only `Set` is supported.",
												},
												"set_parameters": schema.SingleNestedAttribute{
													Computed:    true,
													Description: "Custom variable setting parameters. When there are multiple operations, they are executed in array order. Required when Name is `Set`.",
													Attributes: map[string]schema.Attribute{
														"name": schema.StringAttribute{
															Computed:    true,
															Description: "Custom variable name. The custom variable must be defined before it can be used in operations.",
														},
														"value": schema.StringAttribute{
															Computed:    true,
															Description: "Custom variable assignment. Supports constant strings, variables and formulas; Chinese characters are not supported. Up to 1000 characters.",
														},
													},
												},
											},
										},
									},
									"sub_rules": schema.ListNestedAttribute{
										Computed:    true,
										Description: "Sub-rule list. Rules in this list are executed from top to bottom. SubRules and Actions cannot both be empty, and only one layer of SubRules is supported.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"branches": schema.ListNestedAttribute{
													Computed:    true,
													Description: "Sub-rule branches.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"condition": schema.StringAttribute{
																Computed:    true,
																Description: "Match condition expression.",
															},
															"actions": deepBranchActionsAttr,
															"sub_rules": deepSubRulesAttr,
														},
													},
												},
												"description": schema.ListAttribute{
													Computed:    true,
													ElementType: types.StringType,
													Description: "Rule comments.",
												},
											},
										},
									},
								},
							},
						},
						"description": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "Rule comments. Multiple comments can be filled in.",
						},
					},
				},
			},
		},
	}
}

// attribute type maps for building typed list values. Each map describes the
// attributes of a nested object; nested list attributes are typed as
// types.ListType whose element type is the nested object type.
var (
	customVariableAttrTypes = map[string]attr.Type{
		"name":          types.StringType,
		"initial_value": types.StringType,
		"description":   types.StringType,
	}
	customVariableOperationAttrTypes = map[string]attr.Type{
		"branches":    types.ListType{ElemType: types.ObjectType{AttrTypes: branchAttrTypes}},
		"description": types.ListType{ElemType: types.StringType},
	}
	actionAttrTypes = map[string]attr.Type{
		"name":           types.StringType,
		"set_parameters": setParametersAttrType,
	}
	setParametersAttrType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":  types.StringType,
			"value": types.StringType,
		},
	}
	branchAttrTypes = map[string]attr.Type{
		"condition": types.StringType,
		"actions":   types.ListType{ElemType: types.ObjectType{AttrTypes: actionAttrTypes}},
		"sub_rules": types.ListType{ElemType: types.ObjectType{AttrTypes: subRuleAttrTypes}},
	}
	subRuleAttrTypes = map[string]attr.Type{
		"branches":    types.ListType{ElemType: types.ObjectType{AttrTypes: deepBranchAttrTypes}},
		"description": types.ListType{ElemType: types.StringType},
	}
	deepBranchAttrTypes = map[string]attr.Type{
		"condition": types.StringType,
		"actions":   types.ListType{ElemType: types.ObjectType{AttrTypes: actionAttrTypes}},
		"sub_rules": types.ListType{ElemType: types.ObjectType{AttrTypes: deepSubRuleAttrTypes}},
	}
	deepSubRuleAttrTypes = map[string]attr.Type{
		"description": types.ListType{ElemType: types.StringType},
	}
)

// Read queries the cloud API and writes the result into resp.State.
func (d *TeoZoneCustomVariablesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TeoZoneCustomVariablesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId := config.ZoneId.ValueString()
	if zoneId == "" {
		resp.Diagnostics.AddError(
			"Missing zone_id",
			"zone_id is required and cannot be empty.",
		)
		return
	}

	if d.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	service := NewTeoService(d.Client())
	variables, operations, err := service.DescribeTeoZoneCustomVariables(ctx, zoneId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("querying TEO zone custom variables for zone (%s)", zoneId),
			err.Error(),
		)
		return
	}

	// custom_variables: list of variable definitions; empty list is a valid result.
	variableItems := make([]customVariableItemModel, 0, len(variables))
	for _, v := range variables {
		variableItems = append(variableItems, customVariableItemModel{
			Name:         types.StringPointerValue(v.Name),
			InitialValue: types.StringPointerValue(v.InitialValue),
			Description:  types.StringPointerValue(v.Description),
		})
	}
	variablesList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: customVariableAttrTypes}, variableItems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.CustomVariables = variablesList

	// custom_variable_operations: list of operation rules; empty list is valid.
	operationItems := make([]customVariableOperationItemModel, 0, len(operations))
	for _, op := range operations {
		item, diags := flattenCustomVariableOperation(ctx, op)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		operationItems = append(operationItems, item)
	}
	operationsList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: customVariableOperationAttrTypes}, operationItems)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.CustomVariableOperation = operationsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// flattenCustomVariableOperation converts a CustomVariableOperation into the
// data source model.
func flattenCustomVariableOperation(ctx context.Context, op *teosdk.CustomVariableOperation) (customVariableOperationItemModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	item := customVariableOperationItemModel{}

	descriptionList, d := types.ListValueFrom(ctx, types.StringType, stringPointersToStrings(op.Description))
	diags.Append(d...)
	item.Description = descriptionList

	branches := make([]customVariableOperationBranchModel, 0, len(op.Branches))
	for _, branch := range op.Branches {
		branchModel, d := flattenCustomVariableBranch(ctx, branch)
		diags.Append(d...)
		branches = append(branches, branchModel)
	}
	branchesList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: branchAttrTypes}, branches)
	diags.Append(d...)
	item.Branches = branchesList

	return item, diags
}

// flattenCustomVariableBranch converts a top-level rule branch into the data
// source model.
func flattenCustomVariableBranch(ctx context.Context, branch *teosdk.CustomVariableOperationRuleBranch) (customVariableOperationBranchModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	model := customVariableOperationBranchModel{
		Condition: types.StringPointerValue(branch.Condition),
	}

	actionsList, d := flattenCustomVariableActions(ctx, branch.Actions)
	diags.Append(d...)
	model.Actions = actionsList

	subRules := make([]customVariableOperationSubRuleModel, 0, len(branch.SubRules))
	for _, subRule := range branch.SubRules {
		subRuleModel := customVariableOperationSubRuleModel{}

		descriptionList, d := types.ListValueFrom(ctx, types.StringType, stringPointersToStrings(subRule.Description))
		diags.Append(d...)
		subRuleModel.Description = descriptionList

		branches := make([]customVariableOperationBranchModel, 0, len(subRule.Branches))
		for _, subBranch := range subRule.Branches {
			subBranchModel := customVariableOperationBranchModel{
				Condition: types.StringPointerValue(subBranch.Condition),
			}

			actionsList, d := flattenCustomVariableActions(ctx, subBranch.Actions)
			diags.Append(d...)
			subBranchModel.Actions = actionsList

			deepSubRules := make([]deepSubRuleDescriptionModel, 0, len(subBranch.SubRules))
			for _, deepSubRule := range subBranch.SubRules {
				descriptionList, d := types.ListValueFrom(ctx, types.StringType, stringPointersToStrings(deepSubRule.Description))
				diags.Append(d...)
				deepSubRules = append(deepSubRules, deepSubRuleDescriptionModel{
					Description: descriptionList,
				})
			}
			deepSubRulesList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: deepSubRuleAttrTypes}, deepSubRules)
			diags.Append(d...)
			subBranchModel.SubRules = deepSubRulesList

			branches = append(branches, subBranchModel)
		}
		branchesList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: deepBranchAttrTypes}, branches)
		diags.Append(d...)
		subRuleModel.Branches = branchesList

		subRules = append(subRules, subRuleModel)
	}
	subRulesList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: subRuleAttrTypes}, subRules)
	diags.Append(d...)
	model.SubRules = subRulesList

	return model, diags
}

// flattenCustomVariableActions converts an action list into a typed list.
func flattenCustomVariableActions(ctx context.Context, actions []*teosdk.CustomVariableOperationRuleAction) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	actionItems := make([]customVariableOperationActionModel, 0, len(actions))
	for _, action := range actions {
		actionModel := customVariableOperationActionModel{
			Name: types.StringPointerValue(action.Name),
		}
		if action.SetParameters != nil {
			actionModel.SetParameters = &customVariableSetParamModel{
				Name:  types.StringPointerValue(action.SetParameters.Name),
				Value: types.StringPointerValue(action.SetParameters.Value),
			}
		}
		actionItems = append(actionItems, actionModel)
	}
	actionsList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: actionAttrTypes}, actionItems)
	diags.Append(d...)
	return actionsList, diags
}

// stringPointersToStrings converts []*string into []string, treating nil as
// an empty slice. Nil elements are skipped.
func stringPointersToStrings(input []*string) []string {
	out := make([]string, 0, len(input))
	for _, s := range input {
		if s != nil {
			out = append(out, *s)
		}
	}
	return out
}
