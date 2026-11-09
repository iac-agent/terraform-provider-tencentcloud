package teo

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var _ datasource.DataSourceWithConfigure = &TeoZoneCustomVariablesV5DataSource{}

// NewTeoZoneCustomVariablesV5DataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewTeoZoneCustomVariablesV5DataSource() datasource.DataSource {
	return &TeoZoneCustomVariablesV5DataSource{}
}

// TeoZoneCustomVariablesV5DataSource implements datasource.DataSource for
// the tencentcloud_teo_zone_custom_variables_v5 data source.
type TeoZoneCustomVariablesV5DataSource struct {
	fw.DataSourceWithConfigure
}

// TeoZoneCustomVariablesV5DataSourceModel maps the schema attributes.
type TeoZoneCustomVariablesV5DataSourceModel struct {
	ZoneId                   types.String `tfsdk:"zone_id"`
	CustomVariables          types.List   `tfsdk:"custom_variables"`
	CustomVariableOperations types.List   `tfsdk:"custom_variable_operations"`
}

// zoneCustomVariableModel is the element type of custom_variables.
type zoneCustomVariableModel struct {
	Name         types.String `tfsdk:"name"`
	InitialValue types.String `tfsdk:"initial_value"`
	Description  types.String `tfsdk:"description"`
}

// customVariableOperationModel is the element type of custom_variable_operations.
type customVariableOperationModel struct {
	Branches    types.List `tfsdk:"branches"`
	Description types.List `tfsdk:"description"`
}

// customVariableOperationBranchModel is the element type of the branches
// lists exposed at the first two levels (custom_variable_operations.branches
// and custom_variable_operations.branches.sub_rules.branches). Both levels
// carry the same field names; the element type of the nested lists differs.
type customVariableOperationBranchModel struct {
	Condition types.String `tfsdk:"condition"`
	Actions   types.List   `tfsdk:"actions"`
	SubRules  types.List   `tfsdk:"sub_rules"`
}

// customVariableOperationActionModel is the element type of the actions lists.
type customVariableOperationActionModel struct {
	Name          types.String `tfsdk:"name"`
	SetParameters types.Object `tfsdk:"set_parameters"`
}

// customVariableOperationSubRuleModel is the element type of the sub_rules
// lists exposed at the first two levels. Both levels carry the same field
// names; the element type of Branches differs.
type customVariableOperationSubRuleModel struct {
	Branches    types.List `tfsdk:"branches"`
	Description types.List `tfsdk:"description"`
}

// customVariableOperationBranchLeafModel is the element type of the deepest
// branches list exposed under sub_rules.branches.sub_rules.branches. The
// recursive branch structure of the cloud API is truncated at this level
// because a Terraform schema cannot declare unbounded recursion, and the
// cloud API states that only one level of SubRules is currently supported.
type customVariableOperationBranchLeafModel struct {
	Condition types.String `tfsdk:"condition"`
	Actions   types.List   `tfsdk:"actions"`
}

// Exported aliases used by unit tests in the external test package to
// unmarshal state elements via ElementsAs reflection.
type (
	// ZoneCustomVariableModelForTest mirrors zoneCustomVariableModel.
	ZoneCustomVariableModelForTest = zoneCustomVariableModel
	// CustomVariableOperationModelForTest mirrors customVariableOperationModel.
	CustomVariableOperationModelForTest = customVariableOperationModel
	// CustomVariableOperationBranchModelForTest mirrors customVariableOperationBranchModel.
	CustomVariableOperationBranchModelForTest = customVariableOperationBranchModel
	// CustomVariableOperationActionModelForTest mirrors customVariableOperationActionModel.
	CustomVariableOperationActionModelForTest = customVariableOperationActionModel
	// CustomVariableOperationSubRuleModelForTest mirrors customVariableOperationSubRuleModel.
	CustomVariableOperationSubRuleModelForTest = customVariableOperationSubRuleModel
	// CustomVariableOperationBranchLeafModelForTest mirrors customVariableOperationBranchLeafModel.
	CustomVariableOperationBranchLeafModelForTest = customVariableOperationBranchLeafModel
)

// Attribute type builders ----------------------------------------------------

// setParametersAttrTypes returns the attribute types of a set_parameters object.
func setParametersAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	}
}

// customVariableAttrTypes returns the attribute types of a custom_variables element.
func customVariableAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":          types.StringType,
		"initial_value": types.StringType,
		"description":   types.StringType,
	}
}

// customVariableOperationAttrTypes returns the attribute types of a
// custom_variable_operations element.
func customVariableOperationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"branches": types.ListType{
			ElemType: types.ObjectType{AttrTypes: customVariableOperationBranchAttrTypes()},
		},
		"description": types.ListType{ElemType: types.StringType},
	}
}

// customVariableOperationActionAttrTypes returns the attribute types of an
// actions element.
func customVariableOperationActionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
		"set_parameters": types.ObjectType{
			AttrTypes: setParametersAttrTypes(),
		},
	}
}

// customVariableOperationBranchAttrTypes returns the attribute types of a
// level-one branches element. Its sub_rules.branches elements are level-two
// branches, which still carry sub_rules.
func customVariableOperationBranchAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"condition": types.StringType,
		"actions": types.ListType{
			ElemType: types.ObjectType{AttrTypes: customVariableOperationActionAttrTypes()},
		},
		"sub_rules": types.ListType{
			ElemType: types.ObjectType{AttrTypes: customVariableOperationSubRuleAttrTypes()},
		},
	}
}

// customVariableOperationSubRuleAttrTypes returns the attribute types of a
// level-one sub_rules element. Its branches elements are level-two branches.
func customVariableOperationSubRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"branches": types.ListType{
			ElemType: types.ObjectType{AttrTypes: customVariableOperationBranchL2AttrTypes()},
		},
		"description": types.ListType{ElemType: types.StringType},
	}
}

// customVariableOperationBranchL2AttrTypes returns the attribute types of a
// level-two branches element. Its sub_rules.branches elements are truncated
// leaf branches without further sub_rules.
func customVariableOperationBranchL2AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"condition": types.StringType,
		"actions": types.ListType{
			ElemType: types.ObjectType{AttrTypes: customVariableOperationActionAttrTypes()},
		},
		"sub_rules": types.ListType{
			ElemType: types.ObjectType{AttrTypes: customVariableOperationSubRuleL2AttrTypes()},
		},
	}
}

// customVariableOperationSubRuleL2AttrTypes returns the attribute types of a
// level-two sub_rules element. Its branches elements are truncated leaf
// branches.
func customVariableOperationSubRuleL2AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"branches": types.ListType{
			ElemType: types.ObjectType{AttrTypes: customVariableOperationBranchLeafAttrTypes()},
		},
		"description": types.ListType{ElemType: types.StringType},
	}
}

// customVariableOperationBranchLeafAttrTypes returns the attribute types of
// the deepest branches element exposed under sub_rules.branches.sub_rules.
func customVariableOperationBranchLeafAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"condition": types.StringType,
		"actions": types.ListType{
			ElemType: types.ObjectType{AttrTypes: customVariableOperationActionAttrTypes()},
		},
	}
}

func (d *TeoZoneCustomVariablesV5DataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_zone_custom_variables_v5"
}

func (d *TeoZoneCustomVariablesV5DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the zone-level custom variables of TencentCloud EdgeOne (TEO), " +
			"including variable definitions and variable operation rules.",
		Attributes: map[string]schema.Attribute{
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Zone ID.",
			},
			"custom_variables": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of zone-level custom variables. Empty when no custom variables are configured.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
							Description: "Variable name. The full prefix is required: user.zone.* indicates a zone-level " +
								"custom variable, and user.rule.* indicates a rule-level custom variable. The custom part " +
								"after the prefix only supports upper and lower case letters, numbers and underscores. " +
								"Variable names are case-sensitive and can be up to 50 characters long. The name cannot " +
								"be modified after creation.",
						},
						"initial_value": schema.StringAttribute{
							Computed: true,
							Description: "Initial value of the variable. Constant strings, variables and formulas are " +
								"supported. The length cannot exceed 255 characters.",
						},
						"description": schema.StringAttribute{
							Computed: true,
							Description: "Variable description. The length cannot exceed 60 characters.",
						},
					},
				},
			},
			"custom_variable_operations": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Zone-level custom variable operation rules. Empty when no operation rules are configured.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"branches": schema.ListNestedAttribute{
							Computed: true,
							Description: "Sub-rule branches. Currently only one rule is supported in this list; " +
								"additional rules are invalid. Note: this field may return null, which indicates a " +
								"failure to obtain a valid value.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"condition": schema.StringAttribute{
										Computed:    true,
										Description: "Match condition.",
									},
									"actions": schema.ListNestedAttribute{
										Computed: true,
										Description: "Operations. Note: Actions and SubRules cannot both be empty. " +
											"Note: this field may return null, which indicates a failure to obtain a valid value.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"name": schema.StringAttribute{
													Computed: true,
													Description: "Operation name. The name needs to correspond to the parameter struct, " +
														"for example, when Name=Set, SetParameters is required. Currently only `Set` is supported.",
												},
												"set_parameters": schema.SingleNestedAttribute{
													Computed: true,
													Description: "Custom variable setting parameters. When multiple operations exist in this " +
														"parameter, they are executed in the order of the array. Required when Name is `Set`. " +
														"Note: this field may return null, which indicates a failure to obtain a valid value.",
													Attributes: map[string]schema.Attribute{
														"name": schema.StringAttribute{
															Computed: true,
															Description: "Custom variable name. The custom variable must be defined " +
																"before it can be used in operations.",
														},
														"value": schema.StringAttribute{
															Computed: true,
															Description: "Custom variable value. Constant strings, variables and formulas are " +
																"supported, while Chinese is not supported. The length cannot exceed " +
																"1000 characters.",
														},
													},
												},
											},
										},
									},
									"sub_rules": schema.ListNestedAttribute{
										Computed: true,
										Description: "Sub-rule list. When multiple rules exist in this list, they are executed " +
											"from top to bottom. Note: SubRules and Actions cannot both be empty, and only one " +
											"level of SubRules is currently supported. Note: this field may return null, " +
											"which indicates a failure to obtain a valid value.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"branches": schema.ListNestedAttribute{
													Computed: true,
													Description: "Sub-rule branches. Note: this field may return null, which indicates " +
														"a failure to obtain a valid value.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"condition": schema.StringAttribute{
																Computed:    true,
																Description: "Match condition.",
															},
															"actions": schema.ListNestedAttribute{
																Computed: true,
																Description: "Operations. Note: Actions and SubRules cannot both be empty. " +
																	"Note: this field may return null, which indicates a failure to obtain a valid value.",
																NestedObject: schema.NestedAttributeObject{
																	Attributes: map[string]schema.Attribute{
																		"name": schema.StringAttribute{
																			Computed: true,
																			Description: "Operation name. The name needs to correspond to the parameter struct, " +
																				"for example, when Name=Set, SetParameters is required. Currently only `Set` is supported.",
																		},
																		"set_parameters": schema.SingleNestedAttribute{
																			Computed: true,
																			Description: "Custom variable setting parameters. When multiple operations exist in this " +
																				"parameter, they are executed in the order of the array. Required when Name is `Set`. " +
																				"Note: this field may return null, which indicates a failure to obtain a valid value.",
																			Attributes: map[string]schema.Attribute{
																				"name": schema.StringAttribute{
																					Computed: true,
																					Description: "Custom variable name. The custom variable must be defined " +
																						"before it can be used in operations.",
																				},
																				"value": schema.StringAttribute{
																					Computed: true,
																					Description: "Custom variable value. Constant strings, variables and formulas are " +
																						"supported, while Chinese is not supported. The length cannot exceed " +
																						"1000 characters.",
																				},
																			},
																		},
																	},
																},
															},
															"sub_rules": schema.ListNestedAttribute{
																Computed: true,
																Description: "Sub-rule list. When multiple rules exist in this list, they are executed " +
																	"from top to bottom. Note: SubRules and Actions cannot both be empty, and only one " +
																	"level of SubRules is currently supported. Note: this field may return null, " +
																	"which indicates a failure to obtain a valid value.",
																NestedObject: schema.NestedAttributeObject{
																	Attributes: map[string]schema.Attribute{
																		"branches": schema.ListNestedAttribute{
																			Computed: true,
																			Description: "Sub-rule branches. The recursive branch structure of the cloud API is truncated " +
																				"at this level because only one level of SubRules is currently supported. " +
																				"Note: this field may return null, which indicates a failure to obtain a valid value.",
																			NestedObject: schema.NestedAttributeObject{
																				Attributes: map[string]schema.Attribute{
																					"condition": schema.StringAttribute{
																						Computed:    true,
																						Description: "Match condition.",
																					},
																					"actions": schema.ListNestedAttribute{
																						Computed: true,
																						Description: "Operations. Note: Actions and SubRules cannot both be empty. " +
																							"Note: this field may return null, which indicates a failure to obtain a valid value.",
																						NestedObject: schema.NestedAttributeObject{
																							Attributes: map[string]schema.Attribute{
																								"name": schema.StringAttribute{
																									Computed: true,
																									Description: "Operation name. The name needs to correspond to the parameter struct, " +
																										"for example, when Name=Set, SetParameters is required. Currently only `Set` is supported.",
																								},
																								"set_parameters": schema.SingleNestedAttribute{
																									Computed: true,
																									Description: "Custom variable setting parameters. When multiple operations exist in this " +
																										"parameter, they are executed in the order of the array. Required when Name is `Set`. " +
																										"Note: this field may return null, which indicates a failure to obtain a valid value.",
																									Attributes: map[string]schema.Attribute{
																										"name": schema.StringAttribute{
																											Computed: true,
																											Description: "Custom variable name. The custom variable must be defined " +
																												"before it can be used in operations.",
																										},
																										"value": schema.StringAttribute{
																											Computed: true,
																											Description: "Custom variable value. Constant strings, variables and formulas are " +
																												"supported, while Chinese is not supported. The length cannot exceed " +
																												"1000 characters.",
																										},
																									},
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
																			Description: "Rule remarks. Note: this field may return null, which indicates a failure " +
																				"to obtain a valid value.",
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
													Description: "Rule remarks. Note: this field may return null, which indicates a failure " +
														"to obtain a valid value.",
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
							Description: "Rule remarks. Multiple remarks can be filled in. Note: this field may return null, " +
								"which indicates a failure to obtain a valid value.",
						},
					},
				},
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *TeoZoneCustomVariablesV5DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TeoZoneCustomVariablesV5DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ZoneId.IsNull() || config.ZoneId.IsUnknown() || config.ZoneId.ValueString() == "" {
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

	zoneId := config.ZoneId.ValueString()
	service := NewTeoService(d.Client())
	result, err := service.DescribeTeoZoneCustomVariables(ctx, zoneId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error querying TEO zone custom variables for zone (%s)", zoneId),
			err.Error(),
		)
		return
	}
	if result == nil || result.Response == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error querying TEO zone custom variables for zone (%s)", zoneId),
			"DescribeZoneCustomVariables returned an empty response.",
		)
		return
	}

	// custom_variables: an empty list is a normal result when the zone has no
	// custom variables configured, so an empty slice is written as an empty
	// list instead of an error.
	customVariables := make([]zoneCustomVariableModel, 0, len(result.Response.CustomVariables))
	for _, cv := range result.Response.CustomVariables {
		if cv == nil {
			continue
		}
		customVariables = append(customVariables, zoneCustomVariableModel{
			Name:         types.StringPointerValue(cv.Name),
			InitialValue: types.StringPointerValue(cv.InitialValue),
			Description:  types.StringPointerValue(cv.Description),
		})
	}
	customVariablesList, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: customVariableAttrTypes(),
	}, customVariables)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.CustomVariables = customVariablesList

	// custom_variable_operations: an empty list is a normal result when the
	// zone has no custom variable operations configured.
	operations := make([]customVariableOperationModel, 0, len(result.Response.CustomVariableOperations))
	for _, op := range result.Response.CustomVariableOperations {
		if op == nil {
			continue
		}
		operations = append(operations, customVariableOperationModel{
			Branches:    flattenCustomVariableOperationBranches(ctx, op.Branches),
			Description: stringPointerListValue(op.Description),
		})
	}
	operationsList, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: customVariableOperationAttrTypes(),
	}, operations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.CustomVariableOperations = operationsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// Flatten helpers --------------------------------------------------------

// stringPointerListValue converts a []*string into a types.List of strings.
// A nil slice produces a null list, which faithfully represents the
// "may return null" semantics of the cloud API.
func stringPointerListValue(items []*string) types.List {
	if items == nil {
		return types.ListNull(types.StringType)
	}
	values := make([]attr.Value, 0, len(items))
	for _, item := range items {
		values = append(values, types.StringPointerValue(item))
	}
	list, diags := types.ListValue(types.StringType, values)
	if diags.HasError() {
		return types.ListNull(types.StringType)
	}
	return list
}

// setParametersObject converts an SDK SetParameters pointer into the single
// nested object value used by actions.set_parameters. The cloud API returns
// a single SetParameters object per action.
func setParametersObject(sp *teo.SetParameters) types.Object {
	if sp == nil {
		return types.ObjectNull(setParametersAttrTypes())
	}
	obj, diags := types.ObjectValue(setParametersAttrTypes(), map[string]attr.Value{
		"name":  types.StringPointerValue(sp.Name),
		"value": types.StringPointerValue(sp.Value),
	})
	if diags.HasError() {
		return types.ObjectNull(setParametersAttrTypes())
	}
	return obj
}

// flattenCustomVariableOperationActions converts a list of SDK rule actions
// into the actions model list. A nil slice produces a null list.
func flattenCustomVariableOperationActions(ctx context.Context, actions []*teo.CustomVariableOperationRuleAction) types.List {
	if actions == nil {
		return types.ListNull(types.ObjectType{AttrTypes: customVariableOperationActionAttrTypes()})
	}
	items := make([]customVariableOperationActionModel, 0, len(actions))
	for _, action := range actions {
		if action == nil {
			continue
		}
		items = append(items, customVariableOperationActionModel{
			Name:          types.StringPointerValue(action.Name),
			SetParameters: setParametersObject(action.SetParameters),
		})
	}
	list, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: customVariableOperationActionAttrTypes(),
	}, items)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: customVariableOperationActionAttrTypes()})
	}
	return list
}

// flattenCustomVariableOperationBranches converts a list of SDK level-one
// rule branches into the branches model list. A nil slice produces a null
// list.
func flattenCustomVariableOperationBranches(ctx context.Context, branches []*teo.CustomVariableOperationRuleBranch) types.List {
	if branches == nil {
		return types.ListNull(types.ObjectType{AttrTypes: customVariableOperationBranchAttrTypes()})
	}
	items := make([]customVariableOperationBranchModel, 0, len(branches))
	for _, branch := range branches {
		if branch == nil {
			continue
		}
		items = append(items, customVariableOperationBranchModel{
			Condition: types.StringPointerValue(branch.Condition),
			Actions:   flattenCustomVariableOperationActions(ctx, branch.Actions),
			SubRules:  flattenCustomVariableOperationSubRules(ctx, branch.SubRules),
		})
	}
	list, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: customVariableOperationBranchAttrTypes(),
	}, items)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: customVariableOperationBranchAttrTypes()})
	}
	return list
}

// flattenCustomVariableOperationSubRules converts a list of SDK level-one
// sub rules into the sub_rules model list. A nil slice produces a null list.
func flattenCustomVariableOperationSubRules(ctx context.Context, subRules []*teo.CustomVariableOperationSubRule) types.List {
	if subRules == nil {
		return types.ListNull(types.ObjectType{AttrTypes: customVariableOperationSubRuleAttrTypes()})
	}
	items := make([]customVariableOperationSubRuleModel, 0, len(subRules))
	for _, subRule := range subRules {
		if subRule == nil {
			continue
		}
		items = append(items, customVariableOperationSubRuleModel{
			Branches:    flattenCustomVariableOperationBranchesL2(ctx, subRule.Branches),
			Description: stringPointerListValue(subRule.Description),
		})
	}
	list, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: customVariableOperationSubRuleAttrTypes(),
	}, items)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: customVariableOperationSubRuleAttrTypes()})
	}
	return list
}

// flattenCustomVariableOperationBranchesL2 converts a list of SDK level-two
// rule branches into the branches model list. A nil slice produces a null
// list.
func flattenCustomVariableOperationBranchesL2(ctx context.Context, branches []*teo.CustomVariableOperationRuleBranch) types.List {
	if branches == nil {
		return types.ListNull(types.ObjectType{AttrTypes: customVariableOperationBranchL2AttrTypes()})
	}
	items := make([]customVariableOperationBranchModel, 0, len(branches))
	for _, branch := range branches {
		if branch == nil {
			continue
		}
		items = append(items, customVariableOperationBranchModel{
			Condition: types.StringPointerValue(branch.Condition),
			Actions:   flattenCustomVariableOperationActions(ctx, branch.Actions),
			SubRules:  flattenCustomVariableOperationSubRulesL2(ctx, branch.SubRules),
		})
	}
	list, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: customVariableOperationBranchL2AttrTypes(),
	}, items)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: customVariableOperationBranchL2AttrTypes()})
	}
	return list
}

// flattenCustomVariableOperationSubRulesL2 converts a list of SDK level-two
// sub rules into the sub_rules model list. The branches under this level are
// truncated leaf branches. A nil slice produces a null list.
func flattenCustomVariableOperationSubRulesL2(ctx context.Context, subRules []*teo.CustomVariableOperationSubRule) types.List {
	if subRules == nil {
		return types.ListNull(types.ObjectType{AttrTypes: customVariableOperationSubRuleL2AttrTypes()})
	}
	items := make([]customVariableOperationSubRuleModel, 0, len(subRules))
	for _, subRule := range subRules {
		if subRule == nil {
			continue
		}
		items = append(items, customVariableOperationSubRuleModel{
			Branches:    flattenCustomVariableOperationBranchLeaves(ctx, subRule.Branches),
			Description: stringPointerListValue(subRule.Description),
		})
	}
	list, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: customVariableOperationSubRuleL2AttrTypes(),
	}, items)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: customVariableOperationSubRuleL2AttrTypes()})
	}
	return list
}

// flattenCustomVariableOperationBranchLeaves converts a list of SDK rule
// branches into the truncated leaf branch model list. A nil slice produces a
// null list.
func flattenCustomVariableOperationBranchLeaves(ctx context.Context, branches []*teo.CustomVariableOperationRuleBranch) types.List {
	if branches == nil {
		return types.ListNull(types.ObjectType{AttrTypes: customVariableOperationBranchLeafAttrTypes()})
	}
	items := make([]customVariableOperationBranchLeafModel, 0, len(branches))
	for _, branch := range branches {
		if branch == nil {
			continue
		}
		items = append(items, customVariableOperationBranchLeafModel{
			Condition: types.StringPointerValue(branch.Condition),
			Actions:   flattenCustomVariableOperationActions(ctx, branch.Actions),
		})
	}
	list, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: customVariableOperationBranchLeafAttrTypes(),
	}, items)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: customVariableOperationBranchLeafAttrTypes()})
	}
	return list
}
