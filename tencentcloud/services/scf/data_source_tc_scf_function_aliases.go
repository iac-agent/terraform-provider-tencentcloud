package scf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudScfFunctionAliases() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfFunctionAliasesRead,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Function 名称",
			},

			"namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Function namespace。",
			},

			"function_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "If this parameter is provided，only aliases associated with this function 版本 will be returned。",
			},

			"aliases": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Alias list。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"function_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Master 版本 pointed to by the alias。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alias 名称",
						},
						"routing_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Routing information of aliasNote: this field may return null，indicating that no valid values can be obtained。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"additional_version_weights": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Additional 版本 with random 权重-based routing。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"version": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Function 版本 名称",
												},
												"weight": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "版本 权重",
												},
											},
										},
									},
									"addition_version_matchs": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Additional 版本 with rule-based routing。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"version": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Function 版本 名称",
												},
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Matching rule 键 When the API is called，pass in the `键` to route the request to the specified 版本 based on the matching ruleHeader method:Enter invoke.headers.用户 for `键` and pass in `RoutingKey:{用户:值}` when invoking a function through `invoke` for invocation based on rule matching。",
												},
												"method": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Match method. Valid values:range: range matchexact: exact string match。",
												},
												"expression": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Rule requirements for range match:It should be described in an open or closed range，i.e.，`(a,b)` or `[a,b]`，where both a and b are integersRule requirements for exact match:Exact string match。",
												},
											},
										},
									},
								},
							},
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DescriptionNote: this field may return null，indicating that no valid values can be obtained。",
						},
						"add_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation timeNote: this field may return null，indicating that no valid values can be obtained。",
						},
						"mod_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update timeNote: this field may return null，indicating that no valid values can be obtained。",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudScfFunctionAliasesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_scf_function_aliases.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("function_name"); ok {
		paramMap["FunctionName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		paramMap["Namespace"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("function_version"); ok {
		paramMap["FunctionVersion"] = helper.String(v.(string))
	}

	service := ScfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var aliases []*scf.Alias

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeScfFunctionAliasesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		aliases = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(aliases))
	tmpList := make([]map[string]interface{}, 0, len(aliases))

	if aliases != nil {
		for _, alias := range aliases {
			aliasMap := map[string]interface{}{}

			if alias.FunctionVersion != nil {
				aliasMap["function_version"] = alias.FunctionVersion
			}

			if alias.Name != nil {
				aliasMap["name"] = alias.Name
			}

			if alias.RoutingConfig != nil {
				routingConfigMap := map[string]interface{}{}

				if alias.RoutingConfig.AdditionalVersionWeights != nil {
					additionalVersionWeightsList := []interface{}{}
					for _, additionalVersionWeights := range alias.RoutingConfig.AdditionalVersionWeights {
						additionalVersionWeightsMap := map[string]interface{}{}

						if additionalVersionWeights.Version != nil {
							additionalVersionWeightsMap["version"] = additionalVersionWeights.Version
						}

						if additionalVersionWeights.Weight != nil {
							additionalVersionWeightsMap["weight"] = additionalVersionWeights.Weight
						}

						additionalVersionWeightsList = append(additionalVersionWeightsList, additionalVersionWeightsMap)
					}

					routingConfigMap["additional_version_weights"] = additionalVersionWeightsList
				}

				if alias.RoutingConfig.AddtionVersionMatchs != nil {
					addtionVersionMatchsList := []interface{}{}
					for _, addtionVersionMatchs := range alias.RoutingConfig.AddtionVersionMatchs {
						addtionVersionMatchsMap := map[string]interface{}{}

						if addtionVersionMatchs.Version != nil {
							addtionVersionMatchsMap["version"] = addtionVersionMatchs.Version
						}

						if addtionVersionMatchs.Key != nil {
							addtionVersionMatchsMap["key"] = addtionVersionMatchs.Key
						}

						if addtionVersionMatchs.Method != nil {
							addtionVersionMatchsMap["method"] = addtionVersionMatchs.Method
						}

						if addtionVersionMatchs.Expression != nil {
							addtionVersionMatchsMap["expression"] = addtionVersionMatchs.Expression
						}

						addtionVersionMatchsList = append(addtionVersionMatchsList, addtionVersionMatchsMap)
					}

					routingConfigMap["addition_version_matchs"] = addtionVersionMatchsList
				}

				aliasMap["routing_config"] = []interface{}{routingConfigMap}
			}

			if alias.Description != nil {
				aliasMap["description"] = alias.Description
			}

			if alias.AddTime != nil {
				aliasMap["add_time"] = alias.AddTime
			}

			if alias.ModTime != nil {
				aliasMap["mod_time"] = alias.ModTime
			}

			ids = append(ids, *alias.Name)
			tmpList = append(tmpList, aliasMap)
		}

		_ = d.Set("aliases", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
