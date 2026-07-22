package emr

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudEmrAutoScaleStrategy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEmrAutoScaleStrategyCreate,
		Read:   resourceTencentCloudEmrAutoScaleStrategyRead,
		Update: resourceTencentCloudEmrAutoScaleStrategyUpdate,
		Delete: resourceTencentCloudEmrAutoScaleStrategyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "实例 ID.",
			},

			"strategy_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "1 表示 expansion 和 contraction according 到 load 规则, 2 表示 expansion 和 contraction according 到 时间 规则. Must 是 filled 在 和 match following 规则 策略.",
			},

			"load_auto_scale_strategy": {
				Type:          schema.TypeList,
				ConflictsWith: []string{"time_auto_scale_strategy"},
				Optional:      true,
				Computed:      true,
				Description:   "Expansion 规则 based 在 load.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"strategy_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Rule ID.",
						},
						"strategy_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule 名称.",
						},
						"calm_down_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Cooling 时间 对于 规则 到 take effect.",
						},
						"scale_action": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Expansion 和 contraction actions, 1 表示 expansion, 2 表示 shrinkage.",
						},
						"scale_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "amount 的 expansion 和 contraction each 时间 规则 takes effect.",
						},
						"process_method": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Indicator processing 方法, 1 表示 MAX, 2 表示 MIN, 和 3 表示 AVG.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Rule 优先级, 无效 当 added, defaults 到 auto-increment.",
						},
						"strategy_status": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Rule 状态, 1 表示 已启用, 3 表示 已禁用.",
						},
						"yarn_node_label": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule expansion specifies yarn 节点 label.",
						},
						"period_valid": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Effective 时间 对于 规则 到 take effect.",
						},
						"grace_down_flag": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Elegant shrink switch.",
						},
						"grace_down_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Graceful downsizing waiting 时间.",
						},
						"tags": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Binding 标签 列表.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "标签 键.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "标签 值.",
									},
								},
							},
						},
						"config_group_assigned": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Default 配置 组.",
						},
						"measure_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Expansion 资源 calculation methods, \"DEFAULT\", \"INSTANCE\", \"CPU\", \"MEMORYGB\".\r\n\"DEFAULT\" 表示 默认值 模式, 其中 has same meaning 作为 \"INSTANCE\".\r\n\"INSTANCE\" 表示 calculation based 在 nodes, 默认值 方法.\r\n\"CPU\" 表示 calculated based 在 数量 的 cores 的 machine.\r\n\"MEMORYGB\" 表示 calculated based 在 数量 的 machine 内存.",
						},
						"load_metrics_conditions": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Multiple indicator 触发器 conditions.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"load_metrics": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Trigger 规则 conditions.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"statistic_period": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "regular statistical 周期 provides 1min, 3min, 和 5min.",
												},
												"trigger_threshold": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "数量 的 triggers. 当 数量 的 consecutive triggers exceeds TriggerThreshold, expansion 和 contraction 将 begin.",
												},
												"load_metrics": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Expansion 和 contraction load indicators.",
												},
												"metric_id": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Rule metadata 记录 ID.",
												},
												"conditions": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Trigger condition.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"compare_method": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Conditional comparison 方法, 1 表示 greater 比, 2 表示 less 比, 3 表示 greater 比 或 equal 到, 4 表示 less 比 或 equal 到.",
															},
															"threshold": {
																Type:        schema.TypeFloat,
																Optional:    true,
																Description: "Conditional 阈值.",
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

			"time_auto_scale_strategy": {
				Type:          schema.TypeList,
				ConflictsWith: []string{"load_auto_scale_strategy"},
				Optional:      true,
				Description:   "Rules 对于 scaling up 和 down over 时间.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"strategy_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Policy 名称, 唯一 within 集群.",
						},
						"interval_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "cooling 时间 after 策略 是 triggered. During 此 周期, elastic expansion 和 contraction 将 不 是 triggered.",
						},
						"scale_action": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Expansion 和 contraction actions, 1 表示 expansion, 2 表示 shrinkage.",
						},
						"scale_num": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "数量 的 expansions 和 contractions.",
						},
						"strategy_status": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Rule 状态, 1 表示 有效, 2 表示 无效, 和 3 表示 suspended. Required.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Rule 优先级, smaller 它 是, higher 它 是.",
						},
						"retry_valid_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "当 多个 规则 是 triggered 在 same 时间 和 some 的 them 是 不 actually executed, retries 将 是 made within 此 时间 范围.",
						},
						"repeat_strategy": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "Time expansion 和 contraction repetition strategy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"repeat_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "值 范围 是 \"DAY\", \"DOW\", \"DOM\", 和 \"NONE\", 其中 respectively represent daily repetition, weekly repetition, monthly repetition 和 一个-时间 execution. Required.",
									},
									"day_repeat": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Repeat 规则 通过 day, 有效 当 RepeatType 是 \"DAY\".",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"execute_at_time_of_day": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Repeat 特定 时间 当 任务 是 executed, such 作为 \"01:02:00\".",
												},
												"step": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Executed every Step day.",
												},
											},
										},
									},
									"week_repeat": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Repeat 规则 通过 week, 有效 当 RepeatType 是 \"DOW\".",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"execute_at_time_of_day": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Repeat 特定 时间 当 任务 是 executed, such 作为 \"01:02:00\".",
												},
												"days_of_week": {
													Type:        schema.TypeSet,
													Required:    true,
													Description: "numerical 描述 的 days 的 week, 对于 示例, [1,3,4] 表示 Monday, Wednesday, 和 Thursday every week.",
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
											},
										},
									},
									"month_repeat": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Repeat 规则 通过 month, 有效 当 RepeatType 是 \"DOM\".",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"execute_at_time_of_day": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Repeat 特定 时间 当 任务 是 executed, such 作为 \"01:02:00\".",
												},
												"days_of_month_range": {
													Type:        schema.TypeSet,
													Required:    true,
													Description: "描述 的 day 周期 在 each month, 长度 可以 仅 是 2, 对于 示例, [2,10] 表示 2-10th 的 each month.",
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
											},
										},
									},
									"not_repeat": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Execute 规则 once, effective 当 RepeatType 是 \"NONE\".",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"execute_at": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "特定 和 完整 时间 的 任务 execution, 格式 是 \"2020-07-13 00:00:00\".",
												},
											},
										},
									},
									"expire": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Rule expiration 时间. After 此 时间, 规则 将 automatically 是 placed 在 suspended state, 在 form 的 \"2020-07-23 00:00:00\". Required.",
									},
								},
							},
						},
						"strategy_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Policy 唯一 ID.",
						},
						"grace_down_flag": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Elegant shrink switch.",
						},
						"grace_down_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Graceful downsizing waiting 时间.",
						},
						"tags": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Binding 标签 列表.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "标签 键.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "标签 值.",
									},
								},
							},
						},
						"config_group_assigned": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Default 配置 组.",
						},
						"measure_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Expansion 资源 calculation methods, \"DEFAULT\", \"INSTANCE\", \"CPU\", \"MEMORYGB\".\r\n\"DEFAULT\" 表示 默认值 模式, 其中 has same meaning 作为 \"INSTANCE\".\r\n\"INSTANCE\" 表示 calculation based 在 nodes, 默认值 方法.\r\n\"CPU\" 表示 calculated based 在 数量 的 cores 的 machine.\r\n\"MEMORYGB\" 表示 calculated based 在 数量 的 machine 内存.",
						},
						"terminate_policy": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Destruction strategy, \"DEFAULT\", 默认值 destruction strategy, shrinkage 是 triggered 通过 shrinkage 规则, \"TIMING\" 表示 scheduled destruction.",
						},
						"max_use": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Maximum usage 时间, 秒, 最小 1 hour, 最大 24 hours.",
						},
						"soft_deploy_info": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Node 部署 服务 列表. Only fill 在 HDFS 和 YARN 对于 部署 services. [Mapping relationship 表 corresponding 到 组件 names](https://云.tencent.com/document/product/589/98760).",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"service_node_info": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Start process 列表.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"compensate_flag": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Compensation expansion, 0 表示 不 已启用, 1 表示 已启用.",
						},
						"group_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "scaling 组 ID.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudEmrAutoScaleStrategyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_auto_scale_strategy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		instanceId   string
		strategyType int64
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOkExists("strategy_type"); ok {
		strategyType = int64(v.(int))
	}

	if v, ok := d.GetOk("load_auto_scale_strategy"); ok && strategyType == 1 {
		for _, loadAutoScaleStrategy := range v.([]interface{}) {
			request := emr.NewAddMetricScaleStrategyRequest()
			request.InstanceId = helper.String(instanceId)
			request.StrategyType = helper.Int64(strategyType)
			loadAutoScaleStrategyMap := loadAutoScaleStrategy.(map[string]interface{})
			loadAutoScaleStrategy := emr.LoadAutoScaleStrategy{}
			if v, ok := loadAutoScaleStrategyMap["strategy_id"]; ok {
				loadAutoScaleStrategy.StrategyId = helper.IntInt64(v.(int))
			}
			if v, ok := loadAutoScaleStrategyMap["strategy_name"]; ok {
				loadAutoScaleStrategy.StrategyName = helper.String(v.(string))
			}
			if v, ok := loadAutoScaleStrategyMap["calm_down_time"]; ok {
				loadAutoScaleStrategy.CalmDownTime = helper.IntInt64(v.(int))
			}
			if v, ok := loadAutoScaleStrategyMap["scale_action"]; ok {
				loadAutoScaleStrategy.ScaleAction = helper.IntInt64(v.(int))
			}
			if v, ok := loadAutoScaleStrategyMap["scale_num"]; ok {
				loadAutoScaleStrategy.ScaleNum = helper.IntInt64(v.(int))
			}
			if v, ok := loadAutoScaleStrategyMap["process_method"]; ok {
				loadAutoScaleStrategy.ProcessMethod = helper.IntInt64(v.(int))
			}
			if v, ok := loadAutoScaleStrategyMap["priority"]; ok {
				loadAutoScaleStrategy.Priority = helper.IntInt64(v.(int))
			}
			if v, ok := loadAutoScaleStrategyMap["strategy_status"]; ok {
				loadAutoScaleStrategy.StrategyStatus = helper.IntInt64(v.(int))
			}
			if v, ok := loadAutoScaleStrategyMap["yarn_node_label"]; ok {
				loadAutoScaleStrategy.YarnNodeLabel = helper.String(v.(string))
			}
			if v, ok := loadAutoScaleStrategyMap["period_valid"]; ok {
				loadAutoScaleStrategy.PeriodValid = helper.String(v.(string))
			}
			if v, ok := loadAutoScaleStrategyMap["grace_down_flag"]; ok {
				loadAutoScaleStrategy.GraceDownFlag = helper.Bool(v.(bool))
			}
			if v, ok := loadAutoScaleStrategyMap["grace_down_time"]; ok {
				loadAutoScaleStrategy.GraceDownTime = helper.IntInt64(v.(int))
			}
			if v, ok := loadAutoScaleStrategyMap["tags"]; ok {
				for _, item := range v.([]interface{}) {
					tagsMap := item.(map[string]interface{})
					tag := emr.Tag{}
					if v, ok := tagsMap["tag_key"]; ok {
						tag.TagKey = helper.String(v.(string))
					}
					if v, ok := tagsMap["tag_value"]; ok {
						tag.TagValue = helper.String(v.(string))
					}
					loadAutoScaleStrategy.Tags = append(loadAutoScaleStrategy.Tags, &tag)
				}
			}
			if v, ok := loadAutoScaleStrategyMap["config_group_assigned"]; ok {
				loadAutoScaleStrategy.ConfigGroupAssigned = helper.String(v.(string))
			}
			if v, ok := loadAutoScaleStrategyMap["measure_method"]; ok {
				loadAutoScaleStrategy.MeasureMethod = helper.String(v.(string))
			}
			if loadMetricsConditionsMap, ok := helper.ConvertInterfacesHeadToMap(loadAutoScaleStrategyMap["load_metrics_conditions"]); ok {
				loadMetricsConditions := emr.LoadMetricsConditions{}
				if v, ok := loadMetricsConditionsMap["load_metrics"]; ok {
					for _, item := range v.([]interface{}) {
						loadMetricsMap := item.(map[string]interface{})
						loadMetricsCondition := emr.LoadMetricsCondition{}
						if v, ok := loadMetricsMap["statistic_period"]; ok {
							loadMetricsCondition.StatisticPeriod = helper.IntInt64(v.(int))
						}
						if v, ok := loadMetricsMap["trigger_threshold"]; ok {
							loadMetricsCondition.TriggerThreshold = helper.IntInt64(v.(int))
						}
						if v, ok := loadMetricsMap["load_metrics"]; ok {
							loadMetricsCondition.LoadMetrics = helper.String(v.(string))
						}
						if v, ok := loadMetricsMap["metric_id"]; ok {
							loadMetricsCondition.MetricId = helper.IntInt64(v.(int))
						}
						if v, ok := loadMetricsMap["conditions"]; ok {
							for _, item := range v.([]interface{}) {
								conditionsMap := item.(map[string]interface{})
								triggerCondition := emr.TriggerCondition{}
								if v, ok := conditionsMap["compare_method"]; ok {
									triggerCondition.CompareMethod = helper.IntInt64(v.(int))
								}
								if v, ok := conditionsMap["threshold"]; ok {
									triggerCondition.Threshold = helper.Float64(v.(float64))
								}
								loadMetricsCondition.Conditions = append(loadMetricsCondition.Conditions, &triggerCondition)
							}
						}
						loadMetricsConditions.LoadMetrics = append(loadMetricsConditions.LoadMetrics, &loadMetricsCondition)
					}
				}
				loadAutoScaleStrategy.LoadMetricsConditions = &loadMetricsConditions
			}
			request.LoadAutoScaleStrategy = &loadAutoScaleStrategy
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().AddMetricScaleStrategyWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s create emr auto scale strategy failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	if v, ok := d.GetOk("time_auto_scale_strategy"); ok && strategyType == 2 {
		for _, timeAutoScaleStrategy := range v.([]interface{}) {
			request := emr.NewAddMetricScaleStrategyRequest()
			request.InstanceId = helper.String(instanceId)
			request.StrategyType = helper.Int64(strategyType)
			timeAutoScaleStrategyMap := timeAutoScaleStrategy.(map[string]interface{})
			timeAutoScaleStrategy := emr.TimeAutoScaleStrategy{}
			if v, ok := timeAutoScaleStrategyMap["strategy_name"]; ok {
				timeAutoScaleStrategy.StrategyName = helper.String(v.(string))
			}
			if v, ok := timeAutoScaleStrategyMap["interval_time"]; ok {
				timeAutoScaleStrategy.IntervalTime = helper.IntUint64(v.(int))
			}
			if v, ok := timeAutoScaleStrategyMap["scale_action"]; ok {
				timeAutoScaleStrategy.ScaleAction = helper.IntUint64(v.(int))
			}
			if v, ok := timeAutoScaleStrategyMap["scale_num"]; ok {
				timeAutoScaleStrategy.ScaleNum = helper.IntUint64(v.(int))
			}
			if v, ok := timeAutoScaleStrategyMap["strategy_status"]; ok {
				timeAutoScaleStrategy.StrategyStatus = helper.IntUint64(v.(int))
			}
			if v, ok := timeAutoScaleStrategyMap["priority"]; ok {
				timeAutoScaleStrategy.Priority = helper.IntUint64(v.(int))
			}
			if v, ok := timeAutoScaleStrategyMap["retry_valid_time"]; ok {
				timeAutoScaleStrategy.RetryValidTime = helper.IntUint64(v.(int))
			}
			if repeatStrategyMap, ok := helper.ConvertInterfacesHeadToMap(timeAutoScaleStrategyMap["repeat_strategy"]); ok {
				repeatStrategy := emr.RepeatStrategy{}
				if v, ok := repeatStrategyMap["repeat_type"]; ok {
					repeatStrategy.RepeatType = helper.String(v.(string))
				}
				if dayRepeatMap, ok := helper.ConvertInterfacesHeadToMap(repeatStrategyMap["day_repeat"]); ok {
					dayRepeatStrategy := emr.DayRepeatStrategy{}
					if v, ok := dayRepeatMap["execute_at_time_of_day"]; ok {
						dayRepeatStrategy.ExecuteAtTimeOfDay = helper.String(v.(string))
					}
					if v, ok := dayRepeatMap["step"]; ok {
						dayRepeatStrategy.Step = helper.IntUint64(v.(int))
					}
					repeatStrategy.DayRepeat = &dayRepeatStrategy
				}
				if weekRepeatMap, ok := helper.ConvertInterfacesHeadToMap(repeatStrategyMap["week_repeat"]); ok {
					weekRepeatStrategy := emr.WeekRepeatStrategy{}
					if v, ok := weekRepeatMap["execute_at_time_of_day"]; ok {
						weekRepeatStrategy.ExecuteAtTimeOfDay = helper.String(v.(string))
					}
					if v, ok := weekRepeatMap["days_of_week"]; ok {
						daysOfWeekSet := v.(*schema.Set).List()
						for i := range daysOfWeekSet {
							daysOfWeek := daysOfWeekSet[i].(int)
							weekRepeatStrategy.DaysOfWeek = append(weekRepeatStrategy.DaysOfWeek, helper.IntUint64(daysOfWeek))
						}
					}
					repeatStrategy.WeekRepeat = &weekRepeatStrategy
				}
				if monthRepeatMap, ok := helper.ConvertInterfacesHeadToMap(repeatStrategyMap["month_repeat"]); ok {
					monthRepeatStrategy := emr.MonthRepeatStrategy{}
					if v, ok := monthRepeatMap["execute_at_time_of_day"]; ok {
						monthRepeatStrategy.ExecuteAtTimeOfDay = helper.String(v.(string))
					}
					if v, ok := monthRepeatMap["days_of_month_range"]; ok {
						daysOfMonthRangeSet := v.(*schema.Set).List()
						for i := range daysOfMonthRangeSet {
							daysOfMonthRange := daysOfMonthRangeSet[i].(int)
							monthRepeatStrategy.DaysOfMonthRange = append(monthRepeatStrategy.DaysOfMonthRange, helper.IntUint64(daysOfMonthRange))
						}
					}
					repeatStrategy.MonthRepeat = &monthRepeatStrategy
				}
				if notRepeatMap, ok := helper.ConvertInterfacesHeadToMap(repeatStrategyMap["not_repeat"]); ok {
					notRepeatStrategy := emr.NotRepeatStrategy{}
					if v, ok := notRepeatMap["execute_at"]; ok {
						notRepeatStrategy.ExecuteAt = helper.String(v.(string))
					}
					repeatStrategy.NotRepeat = &notRepeatStrategy
				}
				if v, ok := repeatStrategyMap["expire"]; ok {
					repeatStrategy.Expire = helper.String(v.(string))
				}
				timeAutoScaleStrategy.RepeatStrategy = &repeatStrategy
			}
			if v, ok := timeAutoScaleStrategyMap["strategy_id"]; ok {
				timeAutoScaleStrategy.StrategyId = helper.IntUint64(v.(int))
			}
			if v, ok := timeAutoScaleStrategyMap["grace_down_flag"]; ok {
				timeAutoScaleStrategy.GraceDownFlag = helper.Bool(v.(bool))
			}
			if v, ok := timeAutoScaleStrategyMap["grace_down_time"]; ok {
				timeAutoScaleStrategy.GraceDownTime = helper.IntInt64(v.(int))
			}
			if v, ok := timeAutoScaleStrategyMap["tags"]; ok {
				for _, item := range v.([]interface{}) {
					tagsMap := item.(map[string]interface{})
					tag := emr.Tag{}
					if v, ok := tagsMap["tag_key"]; ok {
						tag.TagKey = helper.String(v.(string))
					}
					if v, ok := tagsMap["tag_value"]; ok {
						tag.TagValue = helper.String(v.(string))
					}
					timeAutoScaleStrategy.Tags = append(timeAutoScaleStrategy.Tags, &tag)
				}
			}
			if v, ok := timeAutoScaleStrategyMap["config_group_assigned"]; ok {
				timeAutoScaleStrategy.ConfigGroupAssigned = helper.String(v.(string))
			}
			if v, ok := timeAutoScaleStrategyMap["measure_method"]; ok {
				timeAutoScaleStrategy.MeasureMethod = helper.String(v.(string))
			}
			if v, ok := timeAutoScaleStrategyMap["terminate_policy"]; ok {
				timeAutoScaleStrategy.TerminatePolicy = helper.String(v.(string))
			}
			if v, ok := timeAutoScaleStrategyMap["max_use"]; ok {
				timeAutoScaleStrategy.MaxUse = helper.IntInt64(v.(int))
			}
			if v, ok := timeAutoScaleStrategyMap["soft_deploy_info"]; ok {
				softDeployInfoSet := v.(*schema.Set).List()
				for i := range softDeployInfoSet {
					softDeployInfo := softDeployInfoSet[i].(int)
					timeAutoScaleStrategy.SoftDeployInfo = append(timeAutoScaleStrategy.SoftDeployInfo, helper.IntInt64(softDeployInfo))
				}
			}
			if v, ok := timeAutoScaleStrategyMap["service_node_info"]; ok {
				serviceNodeInfoSet := v.(*schema.Set).List()
				for i := range serviceNodeInfoSet {
					serviceNodeInfo := serviceNodeInfoSet[i].(int)
					timeAutoScaleStrategy.ServiceNodeInfo = append(timeAutoScaleStrategy.ServiceNodeInfo, helper.IntInt64(serviceNodeInfo))
				}
			}
			if v, ok := timeAutoScaleStrategyMap["compensate_flag"]; ok {
				timeAutoScaleStrategy.CompensateFlag = helper.IntInt64(v.(int))
			}
			if v, ok := timeAutoScaleStrategyMap["group_id"]; ok {
				timeAutoScaleStrategy.GroupId = helper.IntInt64(v.(int))
			}
			request.TimeAutoScaleStrategy = &timeAutoScaleStrategy
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().AddMetricScaleStrategyWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s create emr auto scale strategy failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	d.SetId(strings.Join([]string{instanceId, helper.Int64ToStr(strategyType)}, tccommon.FILED_SP))

	return resourceTencentCloudEmrAutoScaleStrategyRead(d, meta)
}

func resourceTencentCloudEmrAutoScaleStrategyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_auto_scale_strategy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := EMRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	strategyType := idSplit[1]
	_ = d.Set("instance_id", instanceId)
	_ = d.Set("strategy_type", helper.StrToInt(strategyType))
	respData, err := service.DescribeEmrAutoScaleStrategyById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `emr_auto_scale_strategy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	loadAutoScaleStrategiesList := make([]map[string]interface{}, 0, len(respData.LoadAutoScaleStrategies))
	if respData.LoadAutoScaleStrategies != nil && strategyType == "1" {
		for _, loadAutoScaleStrategies := range respData.LoadAutoScaleStrategies {
			loadAutoScaleStrategiesMap := map[string]interface{}{}

			if loadAutoScaleStrategies.StrategyId != nil {
				loadAutoScaleStrategiesMap["strategy_id"] = loadAutoScaleStrategies.StrategyId
			}

			if loadAutoScaleStrategies.StrategyName != nil {
				loadAutoScaleStrategiesMap["strategy_name"] = loadAutoScaleStrategies.StrategyName
			}

			if loadAutoScaleStrategies.CalmDownTime != nil {
				loadAutoScaleStrategiesMap["calm_down_time"] = loadAutoScaleStrategies.CalmDownTime
			}

			if loadAutoScaleStrategies.ScaleAction != nil {
				loadAutoScaleStrategiesMap["scale_action"] = loadAutoScaleStrategies.ScaleAction
			}

			if loadAutoScaleStrategies.ScaleNum != nil {
				loadAutoScaleStrategiesMap["scale_num"] = loadAutoScaleStrategies.ScaleNum
			}

			if loadAutoScaleStrategies.ProcessMethod != nil {
				loadAutoScaleStrategiesMap["process_method"] = loadAutoScaleStrategies.ProcessMethod
			}

			if loadAutoScaleStrategies.Priority != nil {
				loadAutoScaleStrategiesMap["priority"] = loadAutoScaleStrategies.Priority
			}

			if loadAutoScaleStrategies.StrategyStatus != nil {
				loadAutoScaleStrategiesMap["strategy_status"] = loadAutoScaleStrategies.StrategyStatus
			}

			if loadAutoScaleStrategies.YarnNodeLabel != nil {
				loadAutoScaleStrategiesMap["yarn_node_label"] = loadAutoScaleStrategies.YarnNodeLabel
			}

			if loadAutoScaleStrategies.PeriodValid != nil {
				loadAutoScaleStrategiesMap["period_valid"] = loadAutoScaleStrategies.PeriodValid
			}

			if loadAutoScaleStrategies.GraceDownFlag != nil {
				loadAutoScaleStrategiesMap["grace_down_flag"] = loadAutoScaleStrategies.GraceDownFlag
			}

			if loadAutoScaleStrategies.GraceDownTime != nil {
				loadAutoScaleStrategiesMap["grace_down_time"] = loadAutoScaleStrategies.GraceDownTime
			}

			if loadAutoScaleStrategies.ConfigGroupAssigned != nil {
				loadAutoScaleStrategiesMap["config_group_assigned"] = loadAutoScaleStrategies.ConfigGroupAssigned
			}

			if loadAutoScaleStrategies.MeasureMethod != nil {
				loadAutoScaleStrategiesMap["measure_method"] = loadAutoScaleStrategies.MeasureMethod
			}

			loadMetricsConditionsMap := map[string]interface{}{}

			if loadAutoScaleStrategies.LoadMetricsConditions != nil {
				loadMetricsList := make([]map[string]interface{}, 0, len(loadAutoScaleStrategies.LoadMetricsConditions.LoadMetrics))
				if loadAutoScaleStrategies.LoadMetricsConditions.LoadMetrics != nil {
					for _, loadMetrics := range loadAutoScaleStrategies.LoadMetricsConditions.LoadMetrics {
						loadMetricsMap := map[string]interface{}{}

						if loadMetrics.StatisticPeriod != nil {
							loadMetricsMap["statistic_period"] = loadMetrics.StatisticPeriod
						}

						if loadMetrics.TriggerThreshold != nil {
							loadMetricsMap["trigger_threshold"] = loadMetrics.TriggerThreshold
						}

						if loadMetrics.LoadMetrics != nil {
							loadMetricsMap["load_metrics"] = loadMetrics.LoadMetrics
						}

						if loadMetrics.MetricId != nil {
							loadMetricsMap["metric_id"] = loadMetrics.MetricId
						}

						conditionsList := make([]map[string]interface{}, 0, len(loadMetrics.Conditions))
						if loadMetrics.Conditions != nil {
							for _, conditions := range loadMetrics.Conditions {
								conditionsMap := map[string]interface{}{}

								if conditions.CompareMethod != nil {
									conditionsMap["compare_method"] = conditions.CompareMethod
								}

								if conditions.Threshold != nil {
									conditionsMap["threshold"] = conditions.Threshold
								}

								conditionsList = append(conditionsList, conditionsMap)
							}

							loadMetricsMap["conditions"] = conditionsList
						}
						loadMetricsList = append(loadMetricsList, loadMetricsMap)
					}

					loadMetricsConditionsMap["load_metrics"] = loadMetricsList
				}
				loadAutoScaleStrategiesMap["load_metrics_conditions"] = []interface{}{loadMetricsConditionsMap}
			}

			loadAutoScaleStrategiesList = append(loadAutoScaleStrategiesList, loadAutoScaleStrategiesMap)
		}

		_ = d.Set("load_auto_scale_strategy", loadAutoScaleStrategiesList)
	}

	timeBasedAutoScaleStrategiesList := make([]map[string]interface{}, 0, len(respData.TimeBasedAutoScaleStrategies))
	if respData.TimeBasedAutoScaleStrategies != nil && strategyType == "2" {
		for _, timeBasedAutoScaleStrategies := range respData.TimeBasedAutoScaleStrategies {
			timeBasedAutoScaleStrategiesMap := map[string]interface{}{}

			if timeBasedAutoScaleStrategies.StrategyName != nil {
				timeBasedAutoScaleStrategiesMap["strategy_name"] = timeBasedAutoScaleStrategies.StrategyName
			}

			if timeBasedAutoScaleStrategies.IntervalTime != nil {
				timeBasedAutoScaleStrategiesMap["interval_time"] = timeBasedAutoScaleStrategies.IntervalTime
			}

			if timeBasedAutoScaleStrategies.ScaleAction != nil {
				timeBasedAutoScaleStrategiesMap["scale_action"] = timeBasedAutoScaleStrategies.ScaleAction
			}

			if timeBasedAutoScaleStrategies.ScaleNum != nil {
				timeBasedAutoScaleStrategiesMap["scale_num"] = timeBasedAutoScaleStrategies.ScaleNum
			}

			if timeBasedAutoScaleStrategies.StrategyStatus != nil {
				timeBasedAutoScaleStrategiesMap["strategy_status"] = timeBasedAutoScaleStrategies.StrategyStatus
			}

			if timeBasedAutoScaleStrategies.Priority != nil {
				timeBasedAutoScaleStrategiesMap["priority"] = timeBasedAutoScaleStrategies.Priority
			}

			if timeBasedAutoScaleStrategies.RetryValidTime != nil {
				timeBasedAutoScaleStrategiesMap["retry_valid_time"] = timeBasedAutoScaleStrategies.RetryValidTime
			}

			repeatStrategyMap := map[string]interface{}{}

			if timeBasedAutoScaleStrategies.RepeatStrategy != nil {
				if timeBasedAutoScaleStrategies.RepeatStrategy.RepeatType != nil {
					repeatStrategyMap["repeat_type"] = timeBasedAutoScaleStrategies.RepeatStrategy.RepeatType
				}

				if timeBasedAutoScaleStrategies.RepeatStrategy.DayRepeat != nil {
					isEmpty := true
					dayRepeatMap := map[string]interface{}{}
					if timeBasedAutoScaleStrategies.RepeatStrategy.DayRepeat.ExecuteAtTimeOfDay != nil && *timeBasedAutoScaleStrategies.RepeatStrategy.DayRepeat.ExecuteAtTimeOfDay != "" {
						isEmpty = false
						dayRepeatMap["execute_at_time_of_day"] = timeBasedAutoScaleStrategies.RepeatStrategy.DayRepeat.ExecuteAtTimeOfDay
					}

					if timeBasedAutoScaleStrategies.RepeatStrategy.DayRepeat.Step != nil {
						dayRepeatMap["step"] = timeBasedAutoScaleStrategies.RepeatStrategy.DayRepeat.Step
					}
					if !isEmpty {
						repeatStrategyMap["day_repeat"] = []interface{}{dayRepeatMap}
					}
				}

				if timeBasedAutoScaleStrategies.RepeatStrategy.WeekRepeat != nil {
					isEmpty := true
					weekRepeatMap := map[string]interface{}{}
					if timeBasedAutoScaleStrategies.RepeatStrategy.WeekRepeat.ExecuteAtTimeOfDay != nil && *timeBasedAutoScaleStrategies.RepeatStrategy.WeekRepeat.ExecuteAtTimeOfDay != "" {
						weekRepeatMap["execute_at_time_of_day"] = timeBasedAutoScaleStrategies.RepeatStrategy.WeekRepeat.ExecuteAtTimeOfDay
						isEmpty = false
					}

					if timeBasedAutoScaleStrategies.RepeatStrategy.WeekRepeat.DaysOfWeek != nil {
						weekRepeatMap["days_of_week"] = timeBasedAutoScaleStrategies.RepeatStrategy.WeekRepeat.DaysOfWeek
					}

					if !isEmpty {
						repeatStrategyMap["week_repeat"] = []interface{}{weekRepeatMap}
					}
				}

				if timeBasedAutoScaleStrategies.RepeatStrategy.MonthRepeat != nil {
					monthRepeatMap := map[string]interface{}{}
					isEmpty := true
					if timeBasedAutoScaleStrategies.RepeatStrategy.MonthRepeat.ExecuteAtTimeOfDay != nil && *timeBasedAutoScaleStrategies.RepeatStrategy.MonthRepeat.ExecuteAtTimeOfDay != "" {
						monthRepeatMap["execute_at_time_of_day"] = timeBasedAutoScaleStrategies.RepeatStrategy.MonthRepeat.ExecuteAtTimeOfDay
						isEmpty = false
					}

					if timeBasedAutoScaleStrategies.RepeatStrategy.MonthRepeat.DaysOfMonthRange != nil {
						monthRepeatMap["days_of_month_range"] = timeBasedAutoScaleStrategies.RepeatStrategy.MonthRepeat.DaysOfMonthRange
					}
					if !isEmpty {
						repeatStrategyMap["month_repeat"] = []interface{}{monthRepeatMap}
					}
				}

				if timeBasedAutoScaleStrategies.RepeatStrategy.NotRepeat != nil {
					isEmpty := true
					notRepeatMap := map[string]interface{}{}
					if timeBasedAutoScaleStrategies.RepeatStrategy.NotRepeat.ExecuteAt != nil && *timeBasedAutoScaleStrategies.RepeatStrategy.NotRepeat.ExecuteAt != "" {
						notRepeatMap["execute_at"] = timeBasedAutoScaleStrategies.RepeatStrategy.NotRepeat.ExecuteAt
						isEmpty = false
					}

					if !isEmpty {
						repeatStrategyMap["not_repeat"] = []interface{}{notRepeatMap}
					}
				}

				if timeBasedAutoScaleStrategies.RepeatStrategy.Expire != nil {
					repeatStrategyMap["expire"] = timeBasedAutoScaleStrategies.RepeatStrategy.Expire
				}

				timeBasedAutoScaleStrategiesMap["repeat_strategy"] = []interface{}{repeatStrategyMap}
			}

			if timeBasedAutoScaleStrategies.StrategyId != nil {
				timeBasedAutoScaleStrategiesMap["strategy_id"] = timeBasedAutoScaleStrategies.StrategyId
			}

			if timeBasedAutoScaleStrategies.GraceDownFlag != nil {
				timeBasedAutoScaleStrategiesMap["grace_down_flag"] = timeBasedAutoScaleStrategies.GraceDownFlag
			}

			if timeBasedAutoScaleStrategies.GraceDownTime != nil {
				timeBasedAutoScaleStrategiesMap["grace_down_time"] = timeBasedAutoScaleStrategies.GraceDownTime
			}

			tagsList := make([]map[string]interface{}, 0, len(timeBasedAutoScaleStrategies.Tags))
			if timeBasedAutoScaleStrategies.Tags != nil {
				for _, tags := range timeBasedAutoScaleStrategies.Tags {
					tagsMap := map[string]interface{}{}

					if tags.TagKey != nil {
						tagsMap["tag_key"] = tags.TagKey
					}

					if tags.TagValue != nil {
						tagsMap["tag_value"] = tags.TagValue
					}

					tagsList = append(tagsList, tagsMap)
				}

				timeBasedAutoScaleStrategiesMap["tags"] = tagsList
			}
			if timeBasedAutoScaleStrategies.ConfigGroupAssigned != nil {
				timeBasedAutoScaleStrategiesMap["config_group_assigned"] = timeBasedAutoScaleStrategies.ConfigGroupAssigned
			}

			if timeBasedAutoScaleStrategies.MeasureMethod != nil {
				timeBasedAutoScaleStrategiesMap["measure_method"] = timeBasedAutoScaleStrategies.MeasureMethod
			}

			if timeBasedAutoScaleStrategies.TerminatePolicy != nil {
				timeBasedAutoScaleStrategiesMap["terminate_policy"] = timeBasedAutoScaleStrategies.TerminatePolicy
			}

			if timeBasedAutoScaleStrategies.MaxUse != nil {
				timeBasedAutoScaleStrategiesMap["max_use"] = timeBasedAutoScaleStrategies.MaxUse
			}

			if timeBasedAutoScaleStrategies.SoftDeployInfo != nil {
				timeBasedAutoScaleStrategiesMap["soft_deploy_info"] = timeBasedAutoScaleStrategies.SoftDeployInfo
			}

			if timeBasedAutoScaleStrategies.ServiceNodeInfo != nil {
				timeBasedAutoScaleStrategiesMap["service_node_info"] = timeBasedAutoScaleStrategies.ServiceNodeInfo
			}

			if timeBasedAutoScaleStrategies.CompensateFlag != nil {
				timeBasedAutoScaleStrategiesMap["compensate_flag"] = timeBasedAutoScaleStrategies.CompensateFlag
			}

			if timeBasedAutoScaleStrategies.GroupId != nil {
				timeBasedAutoScaleStrategiesMap["group_id"] = timeBasedAutoScaleStrategies.GroupId
			}

			timeBasedAutoScaleStrategiesList = append(timeBasedAutoScaleStrategiesList, timeBasedAutoScaleStrategiesMap)
		}

		_ = d.Set("time_auto_scale_strategy", timeBasedAutoScaleStrategiesList)
	}

	return nil
}

func resourceTencentCloudEmrAutoScaleStrategyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_auto_scale_strategy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := EMRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	strategyType := idSplit[1]

	needChange := false
	mutableArgs := []string{"load_auto_scale_strategy", "time_auto_scale_strategy"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}
	if needChange {
		respData, err := service.DescribeEmrAutoScaleStrategyById(ctx, instanceId)
		if err != nil {
			return err
		}

		if respData == nil {
			return fmt.Errorf("can not find strategy")
		}
		loadAutoScaleStrategyNameIdMap := make(map[string]int64)
		for _, loadAutoScaleStrategy := range respData.LoadAutoScaleStrategies {
			strategyName := loadAutoScaleStrategy.StrategyName
			StrategyId := loadAutoScaleStrategy.StrategyId
			if strategyName != nil && StrategyId != nil {
				loadAutoScaleStrategyNameIdMap[*strategyName] = *StrategyId
			}

		}
		timeBasedScaleStrategyNameIdMap := make(map[string]uint64)
		for _, timeBasedScaleStrategy := range respData.TimeBasedAutoScaleStrategies {
			strategyName := timeBasedScaleStrategy.StrategyName
			StrategyId := timeBasedScaleStrategy.StrategyId
			if strategyName != nil && StrategyId != nil {
				timeBasedScaleStrategyNameIdMap[*strategyName] = *StrategyId
			}

		}
		request := emr.NewModifyAutoScaleStrategyRequest()

		if v, ok := d.GetOk("instance_id"); ok {
			request.InstanceId = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("strategy_type"); ok {
			request.StrategyType = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("load_auto_scale_strategy"); ok && strategyType == "1" {
			for idx, item := range v.([]interface{}) {
				loadAutoScaleStrategiesMap := item.(map[string]interface{})
				loadAutoScaleStrategy := emr.LoadAutoScaleStrategy{}
				if v, ok := loadAutoScaleStrategiesMap["strategy_id"]; ok {
					loadAutoScaleStrategy.StrategyId = helper.IntInt64(v.(int))
				}
				if d.HasChange(fmt.Sprintf("load_auto_scale_strategy.%d.strategy_name", idx)) {
					return fmt.Errorf("can not change strategy name")
				}
				if v, ok := loadAutoScaleStrategiesMap["strategy_name"]; ok {
					loadAutoScaleStrategy.StrategyId = helper.Int64(loadAutoScaleStrategyNameIdMap[v.(string)])
				}

				if v, ok := loadAutoScaleStrategiesMap["calm_down_time"]; ok {
					loadAutoScaleStrategy.CalmDownTime = helper.IntInt64(v.(int))
				}
				if v, ok := loadAutoScaleStrategiesMap["scale_action"]; ok {
					loadAutoScaleStrategy.ScaleAction = helper.IntInt64(v.(int))
				}
				if v, ok := loadAutoScaleStrategiesMap["scale_num"]; ok {
					loadAutoScaleStrategy.ScaleNum = helper.IntInt64(v.(int))
				}
				if v, ok := loadAutoScaleStrategiesMap["process_method"]; ok {
					loadAutoScaleStrategy.ProcessMethod = helper.IntInt64(v.(int))
				}
				if v, ok := loadAutoScaleStrategiesMap["priority"]; ok {
					loadAutoScaleStrategy.Priority = helper.IntInt64(v.(int))
				}
				if v, ok := loadAutoScaleStrategiesMap["strategy_status"]; ok {
					loadAutoScaleStrategy.StrategyStatus = helper.IntInt64(v.(int))
				}
				if v, ok := loadAutoScaleStrategiesMap["yarn_node_label"]; ok {
					loadAutoScaleStrategy.YarnNodeLabel = helper.String(v.(string))
				}
				if v, ok := loadAutoScaleStrategiesMap["period_valid"]; ok {
					loadAutoScaleStrategy.PeriodValid = helper.String(v.(string))
				}
				if v, ok := loadAutoScaleStrategiesMap["grace_down_flag"]; ok {
					loadAutoScaleStrategy.GraceDownFlag = helper.Bool(v.(bool))
				}
				if v, ok := loadAutoScaleStrategiesMap["grace_down_time"]; ok {
					loadAutoScaleStrategy.GraceDownTime = helper.IntInt64(v.(int))
				}
				if v, ok := loadAutoScaleStrategiesMap["tags"]; ok {
					for _, item := range v.([]interface{}) {
						tagsMap := item.(map[string]interface{})
						tag := emr.Tag{}
						if v, ok := tagsMap["tag_key"]; ok {
							tag.TagKey = helper.String(v.(string))
						}
						if v, ok := tagsMap["tag_value"]; ok {
							tag.TagValue = helper.String(v.(string))
						}
						loadAutoScaleStrategy.Tags = append(loadAutoScaleStrategy.Tags, &tag)
					}
				}
				if v, ok := loadAutoScaleStrategiesMap["config_group_assigned"]; ok {
					loadAutoScaleStrategy.ConfigGroupAssigned = helper.String(v.(string))
				}
				if v, ok := loadAutoScaleStrategiesMap["measure_method"]; ok {
					loadAutoScaleStrategy.MeasureMethod = helper.String(v.(string))
				}
				if loadMetricsConditionsMap, ok := helper.ConvertInterfacesHeadToMap(loadAutoScaleStrategiesMap["load_metrics_conditions"]); ok {
					loadMetricsConditions := emr.LoadMetricsConditions{}
					if v, ok := loadMetricsConditionsMap["load_metrics"]; ok {
						for _, item := range v.([]interface{}) {
							loadMetricsMap := item.(map[string]interface{})
							loadMetricsCondition := emr.LoadMetricsCondition{}
							if v, ok := loadMetricsMap["statistic_period"]; ok {
								loadMetricsCondition.StatisticPeriod = helper.IntInt64(v.(int))
							}
							if v, ok := loadMetricsMap["trigger_threshold"]; ok {
								loadMetricsCondition.TriggerThreshold = helper.IntInt64(v.(int))
							}
							if v, ok := loadMetricsMap["load_metrics"]; ok {
								loadMetricsCondition.LoadMetrics = helper.String(v.(string))
							}
							if v, ok := loadMetricsMap["metric_id"]; ok {
								loadMetricsCondition.MetricId = helper.IntInt64(v.(int))
							}
							if v, ok := loadMetricsMap["conditions"]; ok {
								for _, item := range v.([]interface{}) {
									conditionsMap := item.(map[string]interface{})
									triggerCondition := emr.TriggerCondition{}
									if v, ok := conditionsMap["compare_method"]; ok {
										triggerCondition.CompareMethod = helper.IntInt64(v.(int))
									}
									if v, ok := conditionsMap["threshold"]; ok {
										triggerCondition.Threshold = helper.Float64(v.(float64))
									}
									loadMetricsCondition.Conditions = append(loadMetricsCondition.Conditions, &triggerCondition)
								}
							}
							loadMetricsConditions.LoadMetrics = append(loadMetricsConditions.LoadMetrics, &loadMetricsCondition)
						}
					}
					loadAutoScaleStrategy.LoadMetricsConditions = &loadMetricsConditions
				}
				request.LoadAutoScaleStrategies = append(request.LoadAutoScaleStrategies, &loadAutoScaleStrategy)
			}
		}

		if v, ok := d.GetOk("time_auto_scale_strategy"); ok && strategyType == "2" {
			for idx, item := range v.([]interface{}) {
				timeAutoScaleStrategiesMap := item.(map[string]interface{})
				timeAutoScaleStrategy := emr.TimeAutoScaleStrategy{}
				if d.HasChange(fmt.Sprintf("load_auto_scale_strategy.%d.strategy_name", idx)) {
					return fmt.Errorf("can not change strategy name")
				}
				if v, ok := timeAutoScaleStrategiesMap["strategy_name"]; ok {

					timeAutoScaleStrategy.StrategyId = helper.Uint64(timeBasedScaleStrategyNameIdMap[v.(string)])
				}
				if v, ok := timeAutoScaleStrategiesMap["interval_time"]; ok {
					timeAutoScaleStrategy.IntervalTime = helper.IntUint64(v.(int))
				}
				if v, ok := timeAutoScaleStrategiesMap["scale_action"]; ok {
					timeAutoScaleStrategy.ScaleAction = helper.IntUint64(v.(int))
				}
				if v, ok := timeAutoScaleStrategiesMap["scale_num"]; ok {
					timeAutoScaleStrategy.ScaleNum = helper.IntUint64(v.(int))
				}
				if v, ok := timeAutoScaleStrategiesMap["strategy_status"]; ok {
					timeAutoScaleStrategy.StrategyStatus = helper.IntUint64(v.(int))
				}
				if v, ok := timeAutoScaleStrategiesMap["priority"]; ok {
					timeAutoScaleStrategy.Priority = helper.IntUint64(v.(int))
				}
				if v, ok := timeAutoScaleStrategiesMap["retry_valid_time"]; ok {
					timeAutoScaleStrategy.RetryValidTime = helper.IntUint64(v.(int))
				}
				if repeatStrategyMap, ok := helper.ConvertInterfacesHeadToMap(timeAutoScaleStrategiesMap["repeat_strategy"]); ok {
					repeatStrategy := emr.RepeatStrategy{}
					if v, ok := repeatStrategyMap["repeat_type"]; ok {
						repeatStrategy.RepeatType = helper.String(v.(string))
					}
					if dayRepeatMap, ok := helper.ConvertInterfacesHeadToMap(repeatStrategyMap["day_repeat"]); ok {
						dayRepeatStrategy := emr.DayRepeatStrategy{}
						if v, ok := dayRepeatMap["execute_at_time_of_day"]; ok {
							dayRepeatStrategy.ExecuteAtTimeOfDay = helper.String(v.(string))
						}
						if v, ok := dayRepeatMap["step"]; ok {
							dayRepeatStrategy.Step = helper.IntUint64(v.(int))
						}
						repeatStrategy.DayRepeat = &dayRepeatStrategy
					}
					if weekRepeatMap, ok := helper.ConvertInterfacesHeadToMap(repeatStrategyMap["week_repeat"]); ok {
						weekRepeatStrategy := emr.WeekRepeatStrategy{}
						if v, ok := weekRepeatMap["execute_at_time_of_day"]; ok {
							weekRepeatStrategy.ExecuteAtTimeOfDay = helper.String(v.(string))
						}
						if v, ok := weekRepeatMap["days_of_week"]; ok {
							daysOfWeekSet := v.(*schema.Set).List()
							for i := range daysOfWeekSet {
								daysOfWeek := daysOfWeekSet[i].(int)
								weekRepeatStrategy.DaysOfWeek = append(weekRepeatStrategy.DaysOfWeek, helper.IntUint64(daysOfWeek))
							}
						}
						repeatStrategy.WeekRepeat = &weekRepeatStrategy
					}
					if monthRepeatMap, ok := helper.ConvertInterfacesHeadToMap(repeatStrategyMap["month_repeat"]); ok {
						monthRepeatStrategy := emr.MonthRepeatStrategy{}
						if v, ok := monthRepeatMap["execute_at_time_of_day"]; ok {
							monthRepeatStrategy.ExecuteAtTimeOfDay = helper.String(v.(string))
						}
						if v, ok := monthRepeatMap["days_of_month_range"]; ok {
							daysOfMonthRangeSet := v.(*schema.Set).List()
							for i := range daysOfMonthRangeSet {
								daysOfMonthRange := daysOfMonthRangeSet[i].(int)
								monthRepeatStrategy.DaysOfMonthRange = append(monthRepeatStrategy.DaysOfMonthRange, helper.IntUint64(daysOfMonthRange))
							}
						}
						repeatStrategy.MonthRepeat = &monthRepeatStrategy
					}
					if notRepeatMap, ok := helper.ConvertInterfacesHeadToMap(repeatStrategyMap["not_repeat"]); ok {
						notRepeatStrategy := emr.NotRepeatStrategy{}
						if v, ok := notRepeatMap["execute_at"]; ok {
							notRepeatStrategy.ExecuteAt = helper.String(v.(string))
						}
						repeatStrategy.NotRepeat = &notRepeatStrategy
					}
					if v, ok := repeatStrategyMap["expire"]; ok {
						repeatStrategy.Expire = helper.String(v.(string))
					}
					timeAutoScaleStrategy.RepeatStrategy = &repeatStrategy
				}
				if v, ok := timeAutoScaleStrategiesMap["grace_down_flag"]; ok {
					timeAutoScaleStrategy.GraceDownFlag = helper.Bool(v.(bool))
				}
				if v, ok := timeAutoScaleStrategiesMap["grace_down_time"]; ok {
					timeAutoScaleStrategy.GraceDownTime = helper.IntInt64(v.(int))
				}
				if v, ok := timeAutoScaleStrategiesMap["tags"]; ok {
					for _, item := range v.([]interface{}) {
						tagsMap := item.(map[string]interface{})
						tag := emr.Tag{}
						if v, ok := tagsMap["tag_key"]; ok {
							tag.TagKey = helper.String(v.(string))
						}
						if v, ok := tagsMap["tag_value"]; ok {
							tag.TagValue = helper.String(v.(string))
						}
						timeAutoScaleStrategy.Tags = append(timeAutoScaleStrategy.Tags, &tag)
					}
				}
				if v, ok := timeAutoScaleStrategiesMap["config_group_assigned"]; ok {
					timeAutoScaleStrategy.ConfigGroupAssigned = helper.String(v.(string))
				}
				if v, ok := timeAutoScaleStrategiesMap["measure_method"]; ok {
					timeAutoScaleStrategy.MeasureMethod = helper.String(v.(string))
				}
				if v, ok := timeAutoScaleStrategiesMap["terminate_policy"]; ok {
					timeAutoScaleStrategy.TerminatePolicy = helper.String(v.(string))
				}
				if v, ok := timeAutoScaleStrategiesMap["max_use"]; ok {
					timeAutoScaleStrategy.MaxUse = helper.IntInt64(v.(int))
				}
				if v, ok := timeAutoScaleStrategiesMap["soft_deploy_info"]; ok {
					softDeployInfoSet := v.(*schema.Set).List()
					for i := range softDeployInfoSet {
						softDeployInfo := softDeployInfoSet[i].(int)
						timeAutoScaleStrategy.SoftDeployInfo = append(timeAutoScaleStrategy.SoftDeployInfo, helper.IntInt64(softDeployInfo))
					}
				}
				if v, ok := timeAutoScaleStrategiesMap["service_node_info"]; ok {
					serviceNodeInfoSet := v.(*schema.Set).List()
					for i := range serviceNodeInfoSet {
						serviceNodeInfo := serviceNodeInfoSet[i].(int)
						timeAutoScaleStrategy.ServiceNodeInfo = append(timeAutoScaleStrategy.ServiceNodeInfo, helper.IntInt64(serviceNodeInfo))
					}
				}
				if v, ok := timeAutoScaleStrategiesMap["compensate_flag"]; ok {
					timeAutoScaleStrategy.CompensateFlag = helper.IntInt64(v.(int))
				}
				if v, ok := timeAutoScaleStrategiesMap["group_id"]; ok {
					timeAutoScaleStrategy.GroupId = helper.IntInt64(v.(int))
				}
				request.TimeAutoScaleStrategies = append(request.TimeAutoScaleStrategies, &timeAutoScaleStrategy)
			}
		}

		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().ModifyAutoScaleStrategyWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update emr auto scale strategy failed, reason:%+v", logId, err)
			return err
		}
	}

	_ = instanceId
	return resourceTencentCloudEmrAutoScaleStrategyRead(d, meta)
}

func resourceTencentCloudEmrAutoScaleStrategyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_auto_scale_strategy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	strategyType := idSplit[1]

	service := EMRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	respData, err := service.DescribeEmrAutoScaleStrategyById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		return nil
	}
	if strategyType == "1" {
		for _, loadAutoScaleStrategy := range respData.LoadAutoScaleStrategies {
			if err := service.DeleteAutoScaleStrategy(ctx, instanceId, strategyType, *loadAutoScaleStrategy.StrategyId); err != nil {
				return err
			}
		}
	}

	if strategyType == "2" {
		for _, timeBasedAutoScaleStrategy := range respData.TimeBasedAutoScaleStrategies {
			if err := service.DeleteAutoScaleStrategy(ctx, instanceId, strategyType, *helper.UInt64Int64(*timeBasedAutoScaleStrategy.StrategyId)); err != nil {
				return err
			}
		}
	}

	return nil
}
