package tse

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTseCngwStrategy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTseCngwStrategyCreate,
		Read:   resourceTencentCloudTseCngwStrategyRead,
		Update: resourceTencentCloudTseCngwStrategyUpdate,
		Delete: resourceTencentCloudTseCngwStrategyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "网关 ID。",
			},

			"strategy_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "strategy 名称，up 到 20 字符。",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "描述 信息，up 到 120 字符。",
			},

			"config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "配置 的 metric scaling。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_replicas": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "max 数量 副本 对于 metric scaling。",
						},
						"metrics": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "metric 列表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "metric 类型 Deafault 值\n- Resource。",
									},
									"resource_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "指标名称 Reference 值:\n- cpu\n- 内存\n注意：此字段可能返回 null，表示有效值不可用。",
									},
									"target_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "目标 类型 metric，currently 仅 支持 `Utilization`\n注意：此字段可能返回 null，表示有效值不可用。",
									},
									"target_value": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "目标 值 的 metric\n注意：此字段可能返回 null，表示有效值不可用。",
									},
								},
							},
						},
						"create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "创建时间\n注意：此字段可能返回 null，表示有效值不可用。",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "修改时间\n注意：此字段可能返回 null，表示有效值不可用。",
						},
						"strategy_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "strategy ID\n注意：此字段可能返回 null，表示有效值不可用。",
						},
						"behavior": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "behavior 配置 的 metric\n注意：此字段可能返回 null，表示有效值不可用。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"scale_up": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "配置 的 up scale\n注意：此字段可能返回 null，表示有效值不可用。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"stabilization_window_seconds": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "stability window 时间，单位:second，默认值 0 当 scale up\n注意：此字段可能返回 null，表示有效值不可用。",
												},
												"select_policy": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "类型 策略，默认值：max\n注意：此字段可能返回 null，表示有效值不可用。",
												},
												"policies": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "policies 的 scale up\n注意：此字段可能返回 null，表示有效值不可用。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "类型，默认值：Pods\n注意：此字段可能返回 null，表示有效值不可用。",
															},
															"value": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "值\n注意：此字段可能返回 null，表示有效值不可用。",
															},
															"period_seconds": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "周期 的 scale up\n注意：此字段可能返回 null，表示有效值不可用。",
															},
														},
													},
												},
											},
										},
									},
									"scale_down": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "配置 的 down scale\n注意：此字段可能返回 null，表示有效值不可用。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"stabilization_window_seconds": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "stability window 时间，单位:second，默认值 300 当 scale down\n注意：此字段可能返回 null，表示有效值不可用。",
												},
												"select_policy": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "类型 策略，默认值：max\n注意：此字段可能返回 null，表示有效值不可用。",
												},
												"policies": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "policies 的 scale down\n注意：此字段可能返回 null，表示有效值不可用。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "类型，默认值：Pods\n注意：此字段可能返回 null，表示有效值不可用。",
															},
															"value": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "值\n注意：此字段可能返回 null，表示有效值不可用。",
															},
															"period_seconds": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "周期 的 scale down\n注意：此字段可能返回 null，表示有效值不可用。",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"cron_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "配置 的 timed scaling。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"params": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "参数 列表 timed scaling\n注意：此字段可能返回 null，表示有效值不可用。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"period": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "周期 的 timed scaling\n注意：此字段可能返回 null，表示有效值不可用。",
									},
									"start_at": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "开始时间 的 timed scaling\n注意：此字段可能返回 null，表示有效值不可用。",
									},
									"target_replicas": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "数量 目标 nodes 对于 timed scaling. Do 不 exceed max 数量 副本 对于 metric scaling\n注意：此字段可能返回 null，表示有效值不可用。",
									},
									"crontab": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "cron expression 的 timed scaling，无 input 必填\n注意：此字段可能返回 null，表示有效值不可用。",
									},
								},
							},
						},
						"strategy_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "strategy ID\n注意：此字段可能返回 null，表示有效值不可用。",
						},
					},
				},
			},

			"strategy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "strategy ID\n注意：此字段可能返回 null，表示有效值不可用。",
			},
		},
	}
}

func resourceTencentCloudTseCngwStrategyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_strategy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = tse.NewCreateAutoScalerResourceStrategyRequest()
		response   = tse.NewCreateAutoScalerResourceStrategyResponse()
		gatewayId  string
		strategyId string
	)
	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
		request.GatewayId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("strategy_name"); ok {
		request.StrategyName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "config"); ok {
		cloudNativeAPIGatewayStrategyAutoScalerConfig := tse.CloudNativeAPIGatewayStrategyAutoScalerConfig{}
		if v, ok := dMap["max_replicas"]; ok {
			cloudNativeAPIGatewayStrategyAutoScalerConfig.MaxReplicas = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["metrics"]; ok {
			for _, item := range v.([]interface{}) {
				metricsMap := item.(map[string]interface{})
				cloudNativeAPIGatewayStrategyAutoScalerConfigMetric := tse.CloudNativeAPIGatewayStrategyAutoScalerConfigMetric{}
				if v, ok := metricsMap["type"]; ok {
					cloudNativeAPIGatewayStrategyAutoScalerConfigMetric.Type = helper.String(v.(string))
				}
				if v, ok := metricsMap["resource_name"]; ok {
					cloudNativeAPIGatewayStrategyAutoScalerConfigMetric.ResourceName = helper.String(v.(string))
				}
				if v, ok := metricsMap["target_type"]; ok {
					cloudNativeAPIGatewayStrategyAutoScalerConfigMetric.TargetType = helper.String(v.(string))
				}
				if v, ok := metricsMap["target_value"]; ok {
					cloudNativeAPIGatewayStrategyAutoScalerConfigMetric.TargetValue = helper.IntInt64(v.(int))
				}
				cloudNativeAPIGatewayStrategyAutoScalerConfig.Metrics = append(cloudNativeAPIGatewayStrategyAutoScalerConfig.Metrics, &cloudNativeAPIGatewayStrategyAutoScalerConfigMetric)
			}
		}
		if behaviorMap, ok := helper.InterfaceToMap(dMap, "behavior"); ok {
			autoScalerBehavior := tse.AutoScalerBehavior{}
			if scaleUpMap, ok := helper.InterfaceToMap(behaviorMap, "scale_up"); ok {
				autoScalerRules := tse.AutoScalerRules{}
				if v, ok := scaleUpMap["stabilization_window_seconds"]; ok {
					autoScalerRules.StabilizationWindowSeconds = helper.IntInt64(v.(int))
				}
				if v, ok := scaleUpMap["select_policy"]; ok {
					autoScalerRules.SelectPolicy = helper.String(v.(string))
				}
				if v, ok := scaleUpMap["policies"]; ok {
					for _, item := range v.([]interface{}) {
						policiesMap := item.(map[string]interface{})
						autoScalerPolicy := tse.AutoScalerPolicy{}
						if v, ok := policiesMap["type"]; ok {
							autoScalerPolicy.Type = helper.String(v.(string))
						}
						if v, ok := policiesMap["value"]; ok {
							autoScalerPolicy.Value = helper.IntInt64(v.(int))
						}
						if v, ok := policiesMap["period_seconds"]; ok {
							autoScalerPolicy.PeriodSeconds = helper.IntInt64(v.(int))
						}
						autoScalerRules.Policies = append(autoScalerRules.Policies, &autoScalerPolicy)
					}
				}
				autoScalerBehavior.ScaleUp = &autoScalerRules
			}
			if scaleDownMap, ok := helper.InterfaceToMap(behaviorMap, "scale_down"); ok {
				autoScalerRules := tse.AutoScalerRules{}
				if v, ok := scaleDownMap["stabilization_window_seconds"]; ok {
					autoScalerRules.StabilizationWindowSeconds = helper.IntInt64(v.(int))
				}
				if v, ok := scaleDownMap["select_policy"]; ok {
					autoScalerRules.SelectPolicy = helper.String(v.(string))
				}
				if v, ok := scaleDownMap["policies"]; ok {
					for _, item := range v.([]interface{}) {
						policiesMap := item.(map[string]interface{})
						autoScalerPolicy := tse.AutoScalerPolicy{}
						if v, ok := policiesMap["type"]; ok {
							autoScalerPolicy.Type = helper.String(v.(string))
						}
						if v, ok := policiesMap["value"]; ok {
							autoScalerPolicy.Value = helper.IntInt64(v.(int))
						}
						if v, ok := policiesMap["period_seconds"]; ok {
							autoScalerPolicy.PeriodSeconds = helper.IntInt64(v.(int))
						}
						autoScalerRules.Policies = append(autoScalerRules.Policies, &autoScalerPolicy)
					}
				}
				autoScalerBehavior.ScaleDown = &autoScalerRules
			}
			cloudNativeAPIGatewayStrategyAutoScalerConfig.Behavior = &autoScalerBehavior
		}
		request.Config = &cloudNativeAPIGatewayStrategyAutoScalerConfig
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "cron_config"); ok {
		cloudNativeAPIGatewayStrategyCronScalerConfig := tse.CloudNativeAPIGatewayStrategyCronScalerConfig{}
		if v, ok := dMap["params"]; ok {
			for _, item := range v.([]interface{}) {
				paramsMap := item.(map[string]interface{})
				cloudNativeAPIGatewayStrategyCronScalerConfigParam := tse.CloudNativeAPIGatewayStrategyCronScalerConfigParam{}
				if v, ok := paramsMap["period"]; ok {
					cloudNativeAPIGatewayStrategyCronScalerConfigParam.Period = helper.String(v.(string))
				}
				if v, ok := paramsMap["start_at"]; ok {
					cloudNativeAPIGatewayStrategyCronScalerConfigParam.StartAt = helper.String(v.(string))
				}
				if v, ok := paramsMap["target_replicas"]; ok {
					cloudNativeAPIGatewayStrategyCronScalerConfigParam.TargetReplicas = helper.IntInt64(v.(int))
				}
				if v, ok := paramsMap["crontab"]; ok {
					cloudNativeAPIGatewayStrategyCronScalerConfigParam.Crontab = helper.String(v.(string))
				}
				cloudNativeAPIGatewayStrategyCronScalerConfig.Params = append(cloudNativeAPIGatewayStrategyCronScalerConfig.Params, &cloudNativeAPIGatewayStrategyCronScalerConfigParam)
			}
		}
		request.CronConfig = &cloudNativeAPIGatewayStrategyCronScalerConfig
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTseClient().CreateAutoScalerResourceStrategy(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tse cngwStrategy failed, reason:%+v", logId, err)
		return err
	}

	strategyId = *response.Response.StrategyId
	d.SetId(gatewayId + tccommon.FILED_SP + strategyId)

	return resourceTencentCloudTseCngwStrategyRead(d, meta)
}

func resourceTencentCloudTseCngwStrategyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_strategy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	strategyId := idSplit[1]

	cngwStrategy, err := service.DescribeTseCngwStrategyById(ctx, gatewayId, strategyId)
	if err != nil {
		return err
	}

	if cngwStrategy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TseCngwStrategy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("gateway_id", gatewayId)
	_ = d.Set("strategy_id", strategyId)

	if cngwStrategy.StrategyName != nil {
		_ = d.Set("strategy_name", cngwStrategy.StrategyName)
	}

	if cngwStrategy.Description != nil {
		_ = d.Set("description", cngwStrategy.Description)
	}

	if cngwStrategy.Config != nil {
		configMap := map[string]interface{}{}

		if cngwStrategy.Config.MaxReplicas != nil {
			configMap["max_replicas"] = cngwStrategy.Config.MaxReplicas
		}

		if cngwStrategy.Config.Metrics != nil {
			metricsList := []interface{}{}
			for _, metrics := range cngwStrategy.Config.Metrics {
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

		if cngwStrategy.Config.Behavior != nil {
			behaviorMap := map[string]interface{}{}

			if cngwStrategy.Config.Behavior.ScaleUp != nil {
				scaleUpMap := map[string]interface{}{}

				if cngwStrategy.Config.Behavior.ScaleUp.StabilizationWindowSeconds != nil {
					scaleUpMap["stabilization_window_seconds"] = cngwStrategy.Config.Behavior.ScaleUp.StabilizationWindowSeconds
				}

				if cngwStrategy.Config.Behavior.ScaleUp.SelectPolicy != nil {
					scaleUpMap["select_policy"] = cngwStrategy.Config.Behavior.ScaleUp.SelectPolicy
				}

				if cngwStrategy.Config.Behavior.ScaleUp.Policies != nil {
					policiesList := []interface{}{}
					for _, policies := range cngwStrategy.Config.Behavior.ScaleUp.Policies {
						policiesMap := map[string]interface{}{}

						if policies.Type != nil {
							policiesMap["type"] = policies.Type
						}

						if policies.Value != nil {
							policiesMap["value"] = policies.Value
						}

						if policies.PeriodSeconds != nil {
							policiesMap["period_seconds"] = policies.PeriodSeconds
						}

						policiesList = append(policiesList, policiesMap)
					}

					scaleUpMap["policies"] = policiesList
				}

				behaviorMap["scale_up"] = []interface{}{scaleUpMap}
			}

			if cngwStrategy.Config.Behavior.ScaleDown != nil {
				scaleDownMap := map[string]interface{}{}

				if cngwStrategy.Config.Behavior.ScaleDown.StabilizationWindowSeconds != nil {
					scaleDownMap["stabilization_window_seconds"] = cngwStrategy.Config.Behavior.ScaleDown.StabilizationWindowSeconds
				}

				if cngwStrategy.Config.Behavior.ScaleDown.SelectPolicy != nil {
					scaleDownMap["select_policy"] = cngwStrategy.Config.Behavior.ScaleDown.SelectPolicy
				}

				if cngwStrategy.Config.Behavior.ScaleDown.Policies != nil {
					policiesList := []interface{}{}
					for _, policies := range cngwStrategy.Config.Behavior.ScaleDown.Policies {
						policiesMap := map[string]interface{}{}

						if policies.Type != nil {
							policiesMap["type"] = policies.Type
						}

						if policies.Value != nil {
							policiesMap["value"] = policies.Value
						}

						if policies.PeriodSeconds != nil {
							policiesMap["period_seconds"] = policies.PeriodSeconds
						}

						policiesList = append(policiesList, policiesMap)
					}

					scaleDownMap["policies"] = policiesList
				}

				behaviorMap["scale_down"] = []interface{}{scaleDownMap}
			}

			configMap["behavior"] = []interface{}{behaviorMap}
		}

		_ = d.Set("config", []interface{}{configMap})
	}

	if cngwStrategy.CronConfig != nil {
		cronConfigMap := map[string]interface{}{}

		if cngwStrategy.CronConfig.Params != nil {
			paramsList := []interface{}{}
			for _, params := range cngwStrategy.CronConfig.Params {
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

		_ = d.Set("cron_config", []interface{}{cronConfigMap})
	}

	return nil
}

func resourceTencentCloudTseCngwStrategyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_strategy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := tse.NewModifyAutoScalerResourceStrategyRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	strategyId := idSplit[1]

	request.GatewayId = &gatewayId
	request.StrategyId = &strategyId

	immutableArgs := []string{"gateway_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("strategy_name") {
		if v, ok := d.GetOk("strategy_name"); ok {
			request.StrategyName = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "config"); ok {
		cloudNativeAPIGatewayStrategyAutoScalerConfig := tse.CloudNativeAPIGatewayStrategyAutoScalerConfig{}
		if v, ok := dMap["max_replicas"]; ok {
			cloudNativeAPIGatewayStrategyAutoScalerConfig.MaxReplicas = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["metrics"]; ok {
			for _, item := range v.([]interface{}) {
				metricsMap := item.(map[string]interface{})
				cloudNativeAPIGatewayStrategyAutoScalerConfigMetric := tse.CloudNativeAPIGatewayStrategyAutoScalerConfigMetric{}
				if v, ok := metricsMap["type"]; ok {
					cloudNativeAPIGatewayStrategyAutoScalerConfigMetric.Type = helper.String(v.(string))
				}
				if v, ok := metricsMap["resource_name"]; ok {
					cloudNativeAPIGatewayStrategyAutoScalerConfigMetric.ResourceName = helper.String(v.(string))
				}
				if v, ok := metricsMap["target_type"]; ok {
					cloudNativeAPIGatewayStrategyAutoScalerConfigMetric.TargetType = helper.String(v.(string))
				}
				if v, ok := metricsMap["target_value"]; ok {
					cloudNativeAPIGatewayStrategyAutoScalerConfigMetric.TargetValue = helper.IntInt64(v.(int))
				}
				cloudNativeAPIGatewayStrategyAutoScalerConfig.Metrics = append(cloudNativeAPIGatewayStrategyAutoScalerConfig.Metrics, &cloudNativeAPIGatewayStrategyAutoScalerConfigMetric)
			}
		}
		if behaviorMap, ok := helper.InterfaceToMap(dMap, "behavior"); ok {
			autoScalerBehavior := tse.AutoScalerBehavior{}
			if scaleUpMap, ok := helper.InterfaceToMap(behaviorMap, "scale_up"); ok {
				autoScalerRules := tse.AutoScalerRules{}
				if v, ok := scaleUpMap["stabilization_window_seconds"]; ok {
					autoScalerRules.StabilizationWindowSeconds = helper.IntInt64(v.(int))
				}
				if v, ok := scaleUpMap["select_policy"]; ok {
					autoScalerRules.SelectPolicy = helper.String(v.(string))
				}
				if v, ok := scaleUpMap["policies"]; ok {
					for _, item := range v.([]interface{}) {
						policiesMap := item.(map[string]interface{})
						autoScalerPolicy := tse.AutoScalerPolicy{}
						if v, ok := policiesMap["type"]; ok {
							autoScalerPolicy.Type = helper.String(v.(string))
						}
						if v, ok := policiesMap["value"]; ok {
							autoScalerPolicy.Value = helper.IntInt64(v.(int))
						}
						if v, ok := policiesMap["period_seconds"]; ok {
							autoScalerPolicy.PeriodSeconds = helper.IntInt64(v.(int))
						}
						autoScalerRules.Policies = append(autoScalerRules.Policies, &autoScalerPolicy)
					}
				}
				autoScalerBehavior.ScaleUp = &autoScalerRules
			}
			if scaleDownMap, ok := helper.InterfaceToMap(behaviorMap, "scale_down"); ok {
				autoScalerRules := tse.AutoScalerRules{}
				if v, ok := scaleDownMap["stabilization_window_seconds"]; ok {
					autoScalerRules.StabilizationWindowSeconds = helper.IntInt64(v.(int))
				}
				if v, ok := scaleDownMap["select_policy"]; ok {
					autoScalerRules.SelectPolicy = helper.String(v.(string))
				}
				if v, ok := scaleDownMap["policies"]; ok {
					for _, item := range v.([]interface{}) {
						policiesMap := item.(map[string]interface{})
						autoScalerPolicy := tse.AutoScalerPolicy{}
						if v, ok := policiesMap["type"]; ok {
							autoScalerPolicy.Type = helper.String(v.(string))
						}
						if v, ok := policiesMap["value"]; ok {
							autoScalerPolicy.Value = helper.IntInt64(v.(int))
						}
						if v, ok := policiesMap["period_seconds"]; ok {
							autoScalerPolicy.PeriodSeconds = helper.IntInt64(v.(int))
						}
						autoScalerRules.Policies = append(autoScalerRules.Policies, &autoScalerPolicy)
					}
				}
				autoScalerBehavior.ScaleDown = &autoScalerRules
			}
			cloudNativeAPIGatewayStrategyAutoScalerConfig.Behavior = &autoScalerBehavior
		}
		request.Config = &cloudNativeAPIGatewayStrategyAutoScalerConfig
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "cron_config"); ok {
		cloudNativeAPIGatewayStrategyCronScalerConfig := tse.CloudNativeAPIGatewayStrategyCronScalerConfig{}
		if v, ok := dMap["params"]; ok {
			for _, item := range v.([]interface{}) {
				paramsMap := item.(map[string]interface{})
				cloudNativeAPIGatewayStrategyCronScalerConfigParam := tse.CloudNativeAPIGatewayStrategyCronScalerConfigParam{}
				if v, ok := paramsMap["period"]; ok {
					cloudNativeAPIGatewayStrategyCronScalerConfigParam.Period = helper.String(v.(string))
				}
				if v, ok := paramsMap["start_at"]; ok {
					cloudNativeAPIGatewayStrategyCronScalerConfigParam.StartAt = helper.String(v.(string))
				}
				if v, ok := paramsMap["target_replicas"]; ok {
					cloudNativeAPIGatewayStrategyCronScalerConfigParam.TargetReplicas = helper.IntInt64(v.(int))
				}
				if v, ok := paramsMap["crontab"]; ok {
					cloudNativeAPIGatewayStrategyCronScalerConfigParam.Crontab = helper.String(v.(string))
				}
				cloudNativeAPIGatewayStrategyCronScalerConfig.Params = append(cloudNativeAPIGatewayStrategyCronScalerConfig.Params, &cloudNativeAPIGatewayStrategyCronScalerConfigParam)
			}
		}
		request.CronConfig = &cloudNativeAPIGatewayStrategyCronScalerConfig
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTseClient().ModifyAutoScalerResourceStrategy(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tse cngwStrategy failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTseCngwStrategyRead(d, meta)
}

func resourceTencentCloudTseCngwStrategyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tse_cngw_strategy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	strategyId := idSplit[1]

	if err := service.DeleteTseCngwStrategyById(ctx, gatewayId, strategyId); err != nil {
		return err
	}

	return nil
}
