package tse

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTseGatewayCanaryRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTseGatewayCanaryRulesRead,
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "网关 ID。",
			},

			"service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "服务 ID",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "canary 规则 配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"canary_rule_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "canary 规则 列表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"priority": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "优先级 值 ranges 从 0 到 100; larger 值， higher 优先级; 优先级 不能 是 repeated between different 规则。",
									},
									"enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "状态 canary 规则。",
									},
									"condition_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "参数 matching condition 列表。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "类型Reference 值:- 路径- 方法- 查询- 头部- cookie- 正文- 系统。",
												},
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "参数 名称",
												},
												"operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "操作者Reference 值:`le`，`eq`，`lt`，`ne`，`ge`，`gt`，`regex`，`exists`，`在`，`不 在`， `prefix`，`exact`，`regex`。",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "参数 值",
												},
												"delimiter": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "delimiter. 有效 当 操作者 是 在 或 不 在，reference 值:`,`，`;`,`\\n`。",
												},
												"global_config_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "全局 配置 ID。",
												},
												"global_config_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "全局 配置 名称",
												},
											},
										},
									},
									"balanced_service_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "服务 权重 配置。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"service_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "服务 ID",
												},
												"service_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "服务名称",
												},
												"upstream_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "upstream 名称",
												},
												"percent": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "percent，10 是 10%，有效值：0 到 100。",
												},
											},
										},
									},
									"service_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "服务 ID",
									},
									"service_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "服务名称",
									},
								},
							},
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "总数",
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

func dataSourceTencentCloudTseGatewayCanaryRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tse_gateway_canary_rules.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var gatewayId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
		paramMap["GatewayId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_id"); ok {
		paramMap["ServiceId"] = helper.String(v.(string))
	}

	service := TseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var result *tse.CloudAPIGatewayCanaryRuleList

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeTseGatewayCanaryRulesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	cloudAPIGatewayCanaryRuleListMap := map[string]interface{}{}
	if result != nil {
		if result.CanaryRuleList != nil {
			canaryRuleListList := []interface{}{}
			for _, canaryRuleList := range result.CanaryRuleList {
				canaryRuleListMap := map[string]interface{}{}

				if canaryRuleList.Priority != nil {
					canaryRuleListMap["priority"] = canaryRuleList.Priority
				}

				if canaryRuleList.Enabled != nil {
					canaryRuleListMap["enabled"] = canaryRuleList.Enabled
				}

				if canaryRuleList.ConditionList != nil {
					conditionListList := []interface{}{}
					for _, conditionList := range canaryRuleList.ConditionList {
						conditionListMap := map[string]interface{}{}

						if conditionList.Type != nil {
							conditionListMap["type"] = conditionList.Type
						}

						if conditionList.Key != nil {
							conditionListMap["key"] = conditionList.Key
						}

						if conditionList.Operator != nil {
							conditionListMap["operator"] = conditionList.Operator
						}

						if conditionList.Value != nil {
							conditionListMap["value"] = conditionList.Value
						}

						if conditionList.Delimiter != nil {
							conditionListMap["delimiter"] = conditionList.Delimiter
						}

						if conditionList.GlobalConfigId != nil {
							conditionListMap["global_config_id"] = conditionList.GlobalConfigId
						}

						if conditionList.GlobalConfigName != nil {
							conditionListMap["global_config_name"] = conditionList.GlobalConfigName
						}

						conditionListList = append(conditionListList, conditionListMap)
					}

					canaryRuleListMap["condition_list"] = conditionListList
				}

				if canaryRuleList.BalancedServiceList != nil {
					balancedServiceListList := []interface{}{}
					for _, balancedServiceList := range canaryRuleList.BalancedServiceList {
						balancedServiceListMap := map[string]interface{}{}

						if balancedServiceList.ServiceID != nil {
							balancedServiceListMap["service_id"] = balancedServiceList.ServiceID
						}

						if balancedServiceList.ServiceName != nil {
							balancedServiceListMap["service_name"] = balancedServiceList.ServiceName
						}

						if balancedServiceList.UpstreamName != nil {
							balancedServiceListMap["upstream_name"] = balancedServiceList.UpstreamName
						}

						if balancedServiceList.Percent != nil {
							balancedServiceListMap["percent"] = balancedServiceList.Percent
						}

						balancedServiceListList = append(balancedServiceListList, balancedServiceListMap)
					}

					canaryRuleListMap["balanced_service_list"] = balancedServiceListList
				}

				if canaryRuleList.ServiceId != nil {
					canaryRuleListMap["service_id"] = canaryRuleList.ServiceId
				}

				if canaryRuleList.ServiceName != nil {
					canaryRuleListMap["service_name"] = canaryRuleList.ServiceName
				}

				canaryRuleListList = append(canaryRuleListList, canaryRuleListMap)
			}

			cloudAPIGatewayCanaryRuleListMap["canary_rule_list"] = canaryRuleListList
		}

		if result.TotalCount != nil {
			cloudAPIGatewayCanaryRuleListMap["total_count"] = result.TotalCount
		}

		_ = d.Set("result", []interface{}{cloudAPIGatewayCanaryRuleListMap})
	}

	d.SetId(gatewayId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), cloudAPIGatewayCanaryRuleListMap); e != nil {
			return e
		}
	}
	return nil
}
