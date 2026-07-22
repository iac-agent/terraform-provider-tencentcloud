package tse

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTseGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTseGroupsRead,
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "网关 ID。",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 conditions，有效 值:名称,GroupId。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤名称",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "过滤器 值。",
						},
					},
				},
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "groups 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "总数",
						},
						"gateway_group_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "组 列表 网关。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "组 ID。",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "组名称",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "组 描述",
									},
									"node_config": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "组 节点 configration。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"specification": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "组 规格，1c2g|2c4g|4c8g|8c16g。",
												},
												"number": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "组 节点 数量，2-50。",
												},
											},
										},
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "组 状态",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "组 创建时间。",
									},
									"is_first_group": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "whether 它 是 默认值 组- 0: false.- 1: yes。",
									},
									"binding_strategy": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "associated strategy information注意：此字段可能返回 null，表示有效值不可用。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"strategy_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "strategy ID。",
												},
												"strategy_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "strategy name注意：此字段可能返回 null，表示有效值不可用。",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "create time注意：此字段可能返回 null，表示有效值不可用。",
												},
												"modify_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "modify time注意：此字段可能返回 null，表示有效值不可用。",
												},
												"description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "描述 strategy注意：此字段可能返回 null，表示有效值不可用。",
												},
												"config": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "auto scaling configuration注意：此字段可能返回 null，表示有效值不可用。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"max_replicas": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "最大replicas注意：此字段可能返回 null，表示有效值不可用。",
															},
															"metrics": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "metric list注意：此字段可能返回 null，表示有效值不可用。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "metric type注意：此字段可能返回 null，表示有效值不可用。",
																		},
																		"resource_name": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "metric 资源 name注意：此字段可能返回 null，表示有效值不可用。",
																		},
																		"target_type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "metric 目标 type注意：此字段可能返回 null，表示有效值不可用。",
																		},
																		"target_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "metric 目标 value注意：此字段可能返回 null，表示有效值不可用。",
																		},
																	},
																},
															},
															"enabled": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否enable metric auto scaling注意：此字段可能返回 null，表示有效值不可用。",
															},
															"create_time": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "create time注意：此字段可能返回 null，表示有效值不可用。",
															},
															"modify_time": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "modify time注意：此字段可能返回 null，表示有效值不可用。",
															},
															"strategy_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "strategy ID注意：此字段可能返回 null，表示有效值不可用。",
															},
															"auto_scaler_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "auto scaler ID注意：此字段可能返回 null，表示有效值不可用。",
															},
														},
													},
												},
												"gateway_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "网关 ID注意：此字段可能返回 null，表示有效值不可用。",
												},
												"cron_config": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "timing scaling configuration注意：此字段可能返回 null，表示有效值不可用。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"enabled": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否enable timing auto scaling。",
															},
															"params": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "params 的 timing auto scaling。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"period": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "周期 的 timing auto scaling。",
																		},
																		"start_at": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "开始时间。",
																		},
																		"target_replicas": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "目标 replicas。",
																		},
																		"crontab": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "cron expression。",
																		},
																	},
																},
															},
															"create_time": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "创建时间。",
															},
															"modify_time": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "修改时间。",
															},
															"strategy_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "strategy ID。",
															},
														},
													},
												},
												"max_replicas": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "最大replicas。",
												},
											},
										},
									},
									"gateway_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "网关 ID。",
									},
									"internet_max_bandwidth_out": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "公有 网络 outbound 流量 带宽。",
									},
									"modify_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "修改时间。",
									},
									"subnet_ids": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "子网 IDs。",
									},
								},
							},
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

func dataSourceTencentCloudTseGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tse_groups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("gateway_id"); ok {
		paramMap["GatewayId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*tse.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := tse.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	service := TseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var result *tse.NativeGatewayServerGroups
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeTseGroupsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result.GatewayGroupList))
	nativeGatewayServerGroupsMap := map[string]interface{}{}
	if result != nil {
		if result.TotalCount != nil {
			nativeGatewayServerGroupsMap["total_count"] = result.TotalCount
		}

		if result.GatewayGroupList != nil {
			gatewayGroupListList := []interface{}{}
			for _, gatewayGroupList := range result.GatewayGroupList {
				gatewayGroupListMap := map[string]interface{}{}

				if gatewayGroupList.GroupId != nil {
					gatewayGroupListMap["group_id"] = gatewayGroupList.GroupId
				}

				if gatewayGroupList.Name != nil {
					gatewayGroupListMap["name"] = gatewayGroupList.Name
				}

				if gatewayGroupList.Description != nil {
					gatewayGroupListMap["description"] = gatewayGroupList.Description
				}

				if gatewayGroupList.NodeConfig != nil {
					nodeConfigMap := map[string]interface{}{}

					if gatewayGroupList.NodeConfig.Specification != nil {
						nodeConfigMap["specification"] = gatewayGroupList.NodeConfig.Specification
					}

					if gatewayGroupList.NodeConfig.Number != nil {
						nodeConfigMap["number"] = gatewayGroupList.NodeConfig.Number
					}

					gatewayGroupListMap["node_config"] = []interface{}{nodeConfigMap}
				}

				if gatewayGroupList.Status != nil {
					gatewayGroupListMap["status"] = gatewayGroupList.Status
				}

				if gatewayGroupList.CreateTime != nil {
					gatewayGroupListMap["create_time"] = gatewayGroupList.CreateTime
				}

				if gatewayGroupList.IsFirstGroup != nil {
					gatewayGroupListMap["is_first_group"] = gatewayGroupList.IsFirstGroup
				}

				if gatewayGroupList.BindingStrategy != nil {
					bindingStrategyMap := map[string]interface{}{}

					if gatewayGroupList.BindingStrategy.StrategyId != nil {
						bindingStrategyMap["strategy_id"] = gatewayGroupList.BindingStrategy.StrategyId
					}

					if gatewayGroupList.BindingStrategy.StrategyName != nil {
						bindingStrategyMap["strategy_name"] = gatewayGroupList.BindingStrategy.StrategyName
					}

					if gatewayGroupList.BindingStrategy.CreateTime != nil {
						bindingStrategyMap["create_time"] = gatewayGroupList.BindingStrategy.CreateTime
					}

					if gatewayGroupList.BindingStrategy.ModifyTime != nil {
						bindingStrategyMap["modify_time"] = gatewayGroupList.BindingStrategy.ModifyTime
					}

					if gatewayGroupList.BindingStrategy.Description != nil {
						bindingStrategyMap["description"] = gatewayGroupList.BindingStrategy.Description
					}

					if gatewayGroupList.BindingStrategy.Config != nil {
						configMap := map[string]interface{}{}

						if gatewayGroupList.BindingStrategy.Config.MaxReplicas != nil {
							configMap["max_replicas"] = gatewayGroupList.BindingStrategy.Config.MaxReplicas
						}

						if gatewayGroupList.BindingStrategy.Config.Metrics != nil {
							metricsList := []interface{}{}
							for _, metrics := range gatewayGroupList.BindingStrategy.Config.Metrics {
								metricsMap := map[string]interface{}{}

								if metrics.Type != nil {
									metricsMap["type"] = metrics.Type
								}

								if metrics.ResourceName != nil {
									metricsMap["resource_name"] = metrics.ResourceName
								}

								if metrics.TargetType != nil {
									metricsMap["target_type"] = metrics.TargetType
								}

								if metrics.TargetValue != nil {
									metricsMap["target_value"] = metrics.TargetValue
								}

								metricsList = append(metricsList, metricsMap)
							}

							configMap["metrics"] = metricsList
						}

						if gatewayGroupList.BindingStrategy.Config.Enabled != nil {
							configMap["enabled"] = gatewayGroupList.BindingStrategy.Config.Enabled
						}

						if gatewayGroupList.BindingStrategy.Config.CreateTime != nil {
							configMap["create_time"] = gatewayGroupList.BindingStrategy.Config.CreateTime
						}

						if gatewayGroupList.BindingStrategy.Config.ModifyTime != nil {
							configMap["modify_time"] = gatewayGroupList.BindingStrategy.Config.ModifyTime
						}

						if gatewayGroupList.BindingStrategy.Config.StrategyId != nil {
							configMap["strategy_id"] = gatewayGroupList.BindingStrategy.Config.StrategyId
						}

						if gatewayGroupList.BindingStrategy.Config.AutoScalerId != nil {
							configMap["auto_scaler_id"] = gatewayGroupList.BindingStrategy.Config.AutoScalerId
						}

						bindingStrategyMap["config"] = []interface{}{configMap}
					}

					if gatewayGroupList.BindingStrategy.GatewayId != nil {
						bindingStrategyMap["gateway_id"] = gatewayGroupList.BindingStrategy.GatewayId
					}

					if gatewayGroupList.BindingStrategy.CronConfig != nil {
						cronConfigMap := map[string]interface{}{}

						if gatewayGroupList.BindingStrategy.CronConfig.Enabled != nil {
							cronConfigMap["enabled"] = gatewayGroupList.BindingStrategy.CronConfig.Enabled
						}

						if gatewayGroupList.BindingStrategy.CronConfig.Params != nil {
							paramsList := []interface{}{}
							for _, params := range gatewayGroupList.BindingStrategy.CronConfig.Params {
								paramsMap := map[string]interface{}{}

								if params.Period != nil {
									paramsMap["period"] = params.Period
								}

								if params.StartAt != nil {
									paramsMap["start_at"] = params.StartAt
								}

								if params.TargetReplicas != nil {
									paramsMap["target_replicas"] = params.TargetReplicas
								}

								if params.Crontab != nil {
									paramsMap["crontab"] = params.Crontab
								}

								paramsList = append(paramsList, paramsMap)
							}

							cronConfigMap["params"] = paramsList
						}

						if gatewayGroupList.BindingStrategy.CronConfig.CreateTime != nil {
							cronConfigMap["create_time"] = gatewayGroupList.BindingStrategy.CronConfig.CreateTime
						}

						if gatewayGroupList.BindingStrategy.CronConfig.ModifyTime != nil {
							cronConfigMap["modify_time"] = gatewayGroupList.BindingStrategy.CronConfig.ModifyTime
						}

						if gatewayGroupList.BindingStrategy.CronConfig.StrategyId != nil {
							cronConfigMap["strategy_id"] = gatewayGroupList.BindingStrategy.CronConfig.StrategyId
						}

						bindingStrategyMap["cron_config"] = []interface{}{cronConfigMap}
					}

					if gatewayGroupList.BindingStrategy.MaxReplicas != nil {
						bindingStrategyMap["max_replicas"] = gatewayGroupList.BindingStrategy.MaxReplicas
					}

					gatewayGroupListMap["binding_strategy"] = []interface{}{bindingStrategyMap}
				}

				if gatewayGroupList.GatewayId != nil {
					gatewayGroupListMap["gateway_id"] = gatewayGroupList.GatewayId
				}

				if gatewayGroupList.InternetMaxBandwidthOut != nil {
					gatewayGroupListMap["internet_max_bandwidth_out"] = gatewayGroupList.InternetMaxBandwidthOut
				}

				if gatewayGroupList.ModifyTime != nil {
					gatewayGroupListMap["modify_time"] = gatewayGroupList.ModifyTime
				}

				if gatewayGroupList.SubnetIds != nil {
					gatewayGroupListMap["subnet_ids"] = gatewayGroupList.SubnetIds
				}

				gatewayGroupListList = append(gatewayGroupListList, gatewayGroupListMap)
				ids = append(ids, *gatewayGroupList.GroupId)
			}

			nativeGatewayServerGroupsMap["gateway_group_list"] = gatewayGroupListList
		}

		_ = d.Set("result", []interface{}{nativeGatewayServerGroupsMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), nativeGatewayServerGroupsMap); e != nil {
			return e
		}
	}
	return nil
}
