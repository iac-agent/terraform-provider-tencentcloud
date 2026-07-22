package monitor

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMonitorAlarmPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMonitorAlarmPolicyRead,
		Schema: map[string]*schema.Schema{
			"module": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "值 fixed 在 监控。",
			},

			"policy_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Fuzzy search 通过 策略 名称",
			},

			"monitor_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "过滤器 通过 监控类型 有效值：MT_QCE (Tencent Cloud 服务 监控). 如果此参数为空，all 将 是 queried 通过 默认值。",
			},

			"namespaces": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "过滤器 通过 命名空间. For 值 的 different 策略 types，please see:[Poicy 类型 List](https://www.tencentcloud.com/document/product/248/39565?has_map=1)。",
			},

			"dimensions": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "告警 对象 列表，其中 是 JSON 字符串. outer 数组 corresponds 到 多个 实例，和 inner 数组 是 dimension 的 对象.For 示例，'CVM - Basic Monitor' 可以 是 written 作为: [ {Dimensions: {unInstanceId: ins-qr8d555g}}，{Dimensions: {unInstanceId: ins-qr8d555h}} ]You 可以 also refer 到 'Example 2' below.For more 信息 在 参数 samples 的 different Tencent Cloud services，see [Product Policy 类型 和 Dimension Information](https://www.tencentcloud.com/document/product/248/39565?has_map=1).注意: 如果 1 是 passed 在 对于 NeedCorrespondence， relationship between 策略 和 实例 needs 到 是 返回. You 可以 pass 在 up 到 20 告警 对象 dimensions 到 avoid 请求 超时。",
			},

			"receiver_uids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Search 通过 recipient. You 可以 get 用户 列表 使用 API [ListUsers](https://www.tencentcloud.com/document/product/598/34587?from_cn_redirect=1) 在 'Cloud Access Management' 或 查询 sub-用户 信息 使用 API [GetUser](https://www.tencentcloud.com/document/product/598/34590?from_cn_redirect=1). Uid 字段 在 返回 结果 should 是 entered here。",
			},

			"receiver_groups": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Search 通过 recipient 组. You 可以 get 用户 组 列表 使用 API [ListGroups](https://www.tencentcloud.com/document/product/598/34589?from_cn_redirect=1) 在 'Cloud Access Management' 或 查询 用户 组 列表 其中 sub-用户 是 在 使用 API [ListGroupsForUser](https://www.tencentcloud.com/document/product/598/34588?from_cn_redirect=1). GroupId 字段 在 返回 结果 should 是 entered here。",
			},

			"policy_type": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "过滤器 通过 默认值 策略. 有效值：DEFAULT (display 默认值 策略)，NOT_DEFAULT (display non-默认值 policies). 如果此参数为空，all policies 将 是 displayed。",
			},

			"field": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "排序方式 字段. For 示例，到 排序方式 last 修改时间，使用 Field: UpdateTime。",
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "排序顺序 有效值：ASC (ascending)，DESC (descending)。",
			},

			"project_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "ID 数组 策略 项目，其中 可以 是 viewed 在 following 页面: [Project Management](https://console.tencentcloud.com/项目)。",
			},

			"notice_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "列表 通知 template IDs，其中 可以 是 获取 通过 querying 通知 template 列表.It 可以 是 queried 使用 API [DescribeAlarmNotices](https://www.tencentcloud.com/document/product/248/39300)。",
			},

			"rule_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "过滤器 通过 触发器 condition. 有效值：STATIC (display policies 使用 静态 阈值)，DYNAMIC (display policies 使用 动态 阈值). 如果此参数为空，all policies 将 是 displayed。",
			},

			"enable": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "过滤器 通过 告警状态 有效值：[1]: 已启用; [0]: 已禁用; [0，1]: all。",
			},

			"not_binding_notice_rule": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "如果 1 是 passed 在，告警 policies 使用 无 通知 规则 已配置 是 queried. 如果 它 是 left 空 或 other 值 是 passed 在，all 告警 policies 是 queried。",
			},

			"instance_group_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "实例 组 ID",
			},

			"need_correspondence": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "是否relationship between 策略 和 input 参数 过滤器 dimension 为必填项. 1: Yes. 0: No. 默认值：0。",
			},

			"trigger_tasks": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 告警 策略 通过 triggered 任务 (such 作为 auto scaling 任务). Up 到 10 tasks 可以 是 指定。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Triggered 任务 类型 有效 值: AS (auto scaling)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"task_config": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration 信息 在 JSON 格式，such 作为 {Key1:Value1,Key2:Value2}注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
					},
				},
			},

			"one_click_policy_type": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "过滤器 通过 quick 告警 策略. 如果此参数为空，all policies 是 displayed. ONECLICK: Display quick 告警 policies; NOT_ONECLICK: Display non-quick 告警 policies。",
			},

			"not_bind_all": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "是否returned 结果 needs 到 过滤器 policies associated 使用 all objects. 有效值：1 (Yes)，0 (No)。",
			},

			"not_instance_group": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "是否returned 结果 needs 到 过滤器 policies associated 使用 实例 groups. 有效值：1 (Yes)，0 (No)。",
			},

			"prom_ins_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ID TencentCloud Managed Service 对于 Prometheus 实例，其中 是 用于customizing metric 策略。",
			},

			"receiver_on_call_form_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Search 通过 调度。",
			},

			"policies": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Policy 数组。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm 策略 IDNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm 策略 nameNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "RemarksNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"monitor_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "监控类型 有效值：MT_QCE (Tencent Cloud 服务 监控)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"enable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "状态 有效值：0 (已禁用)，1 (已启用)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"use_sum": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 实例 bound 到 策略 groupNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "项目 ID 有效值：-1 (无 项目)，0 (默认值 项目)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project nameNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm 策略 typeNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"condition_template_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Trigger condition template IDNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"condition": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Metric 触发器 conditionNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_union_rule": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Judgment condition 的 告警 触发器 condition (0: Any; 1: All; 2: Composite). 当 值 是 集合 到 2 (i.e.，composite 触发器 conditions)，此 参数 should 是 使用 together 使用 ComplexExpression.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"rules": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Alarm 触发器 condition listNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"metric_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "指标名称 或 事件名称 支持 metrics 可以 是 queried via [DescribeAlarmMetrics](https://www.tencentcloud.com/document/product/248/39322) 和 支持 events via [DescribeAlarmEvents](https://www.tencentcloud.com/document/product/248/39324).注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
												},
												"period": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Statistical 周期 （秒）。 有效 值 可以 是 queried via DescribeAlarmMetrics.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
												},
												"operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Operatorintelligent = intelligent detection without thresholdeq = equal toge = greater 比 或 equal togt = greater thanle = less 比 或 equal tolt = less thanne = 不 equal today_increase = day-在-day increaseday_decrease = day-在-day decreaseday_wave = day-在-day fluctuationweek_increase = week-在-week increaseweek_decrease = week-在-week decreaseweek_wave = week-在-week fluctuationcycle_increase = cyclical increasecycle_decrease = cyclical decreasecycle_wave = cyclical fluctuationre = regex matchThe 有效 值 可以 是 queried via [DescribeAlarmMetrics](https://www.tencentcloud.com/document/product/248/39322)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Threshold. 有效 值 范围 可以 是 queried via [DescribeAlarmMetrics](https://www.tencentcloud.com/document/product/248/39322)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
												},
												"continue_period": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "数量 periods. 1: continue 对于 一个 周期; 2: continue 对于 two periods; 和 so 在. 有效 值 可以 是 queried via [DescribeAlarmMetrics](https://www.tencentcloud.com/document/product/248/39322)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
												},
												"notice_frequency": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Alarm 间隔 （秒）。 有效值：0 (do 不 repeat)，300 (告警 once every 5 minutes)，600 (告警 once every 10 minutes)，900 (告警 once every 15 minutes)，1800 (告警 once every 30 minutes)，3600 (告警 once every hour)，7200 (告警 once every 2 hours)，10800 (告警 once every 3 hours)，21600 (告警 once every 6 hours)，43200 (告警 once every 12 hours)，86400 (告警 once every dayNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"is_power_notice": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "是否alarm 频率 increases exponentially. 有效值：0 (无)，1 (yesNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"filter": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "过滤器 condition 对于 一个 单个 触发器 rulNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "过滤器 条件类型 有效值：DIMENSION (uses dimensions 对于 filteringNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
															},
															"dimensions": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "JSON 字符串 generated 通过 serializing AlarmPolicyDimension two-dimensional 数组. 一个-dimensional arrays 是 在 OR relationship，和 elements 在 一个-dimensional 数组 是 在 AND relationshiNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
															},
														},
													},
												},
												"description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Metric display 名称，其中 是 使用 在 output parameteNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"unit": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Unit，其中 是 使用 在 output parameteNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"rule_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Trigger 条件类型 STATIC: 静态 阈值; 动态: 动态 阈值. 如果 您 do 不 指定this 参数 当 creating 或 editing 策略，STATIC 是 使用 通过 defaultNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
												},
												"is_advanced": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "是否为an advanced metric. 0: No; 1: Yes注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"is_open": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "是否advanced metric 功能 是 已启用 0: No; 1: Yes注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"product_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Integration center product ID注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"value_max": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Maximum valu注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"value_min": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Minimum valu注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"hierarchical_value": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "配置 的 告警 级别 threshol注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"remind": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Threshold 对于 Remind leve注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"warn": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Threshold 对于 Warn leve注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"serious": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Threshold 对于 Serious leve注意：此字段可能返回 null，表示无法获取有效值。",
															},
														},
													},
												},
											},
										},
									},
									"complex_expression": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "judgment expression 的 composite 告警 触发器 conditions，其中 是 有效 当 值 的 IsUnionRule 是 2. 此 参数 是 用于determine 该 告警 condition 是 met 仅 当 expression 值 是 True 对于 多个 触发器 conditions注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"event_condition": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Event 触发器 conditioNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rules": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Alarm 触发器 condition lisNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"metric_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "指标名称 或 事件名称 支持 metrics 可以 是 queried via DescribeAlarmMetrics 和 支持 events via DescribeAlarmEventsNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
												},
												"period": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Statistical 周期 （秒）。 有效 值 可以 是 queried via DescribeAlarmMetricsNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
												},
												"operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Statistical 周期 （秒）。 有效 值 可以 是 queried via DescribeAlarmMetrics.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取.操作者	String	No	Operatorintelligent = intelligent detection without thresholdeq = equal toge = greater 比 或 equal togt = greater thanle = less 比 或 equal tolt = less thanne = 不 equal today_increase = day-在-day increaseday_decrease = day-在-day decreaseday_wave = day-在-day fluctuationweek_increase = week-在-week increaseweek_decrease = week-在-week decreaseweek_wave = week-在-week fluctuationcycle_increase = cyclical increasecycle_decrease = cyclical decreasecycle_wave = cyclical fluctuationre = regex matchThe 有效 值 可以 是 queried via DescribeAlarmMetrics.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Threshold. 有效 值 范围 可以 是 queried via DescribeAlarmMetrics.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
												},
												"continue_period": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "数量 periods. 1: continue 对于 一个 周期; 2: continue 对于 two periods; 和 so 在. 有效 值 可以 是 queried via DescribeAlarmMetrics.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
												},
												"notice_frequency": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Alarm 间隔 （秒）。 有效值：0 (do 不 repeat)，300 (告警 once every 5 minutes)，600 (告警 once every 10 minutes)，900 (告警 once every 15 minutes)，1800 (告警 once every 30 minutes)，3600 (告警 once every hour)，7200 (告警 once every 2 hours)，10800 (告警 once every 3 hours)，21600 (告警 once every 6 hours)，43200 (告警 once every 12 hours)，86400 (告警 once every day)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"is_power_notice": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "是否alarm 频率 increases exponentially. 有效值：0 (无)，1 (yes)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"filter": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "过滤器 condition 对于 一个 单个 触发器 ruleNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "过滤器 条件类型 有效值：DIMENSION (uses dimensions 对于 filtering)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
															},
															"dimensions": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "JSON 字符串 generated 通过 serializing AlarmPolicyDimension two-dimensional 数组. 一个-dimensional arrays 是 在 OR relationship，和 elements 在 一个-dimensional 数组 是 在 AND relationshipNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
															},
														},
													},
												},
												"description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Metric display 名称，其中 是 使用 在 output parameterNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"unit": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Unit，其中 是 使用 在 output parameterNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"rule_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Trigger 条件类型 STATIC: 静态 阈值; 动态: 动态 阈值. 如果 您 do 不 指定this 参数 当 creating 或 editing 策略，STATIC 是 使用 通过 默认值.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
												},
												"is_advanced": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "是否为an advanced metric. 0: No; 1: Yes.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"is_open": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "是否advanced metric 功能 是 已启用 0: No; 1: Yes.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"product_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Integration center product ID.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"value_max": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Maximum value注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"value_min": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Minimum value注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"hierarchical_value": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "配置 的 告警 级别 threshold注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"remind": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Threshold 对于 Remind level注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"warn": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Threshold 对于 Warn level注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"serious": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Threshold 对于 Serious level注意：此字段可能返回 null，表示无法获取有效值。",
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
						"notice_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Notification 规则 ID listNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"notices": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Notification 规则 listNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Alarm 通知 template IDNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Alarm 通知 template nameNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
									},
									"updated_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Last modified timeNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
									},
									"updated_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Last modified byNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
									},
									"notice_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Alarm 通知 类型 有效值：ALARM (对于 unresolved alarms)，OK (对于 resolved alarms)，ALL (对于 all alarms)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
									},
									"user_notices": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "用户 通知 listNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"receiver_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Recipient 类型 有效值：USER (用户)，GROUP (用户 组)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"start_time": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Notification 开始时间，其中 是 expressed 通过 数量 秒 since 00:00:00. 取值范围：0-86399Note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"end_time": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Notification 结束时间，其中 是 expressed 通过 数量 秒 since 00:00:00. 取值范围：0-86399Note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"notice_way": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "Notification channel 列表. 有效值：EMAIL (email)，SMS (SMS)，CALL (phone)，WECHAT (WeChat)，RTX (WeCom)注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"user_ids": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Computed:    true,
													Description: "用户 uid listNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"group_ids": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Computed:    true,
													Description: "用户 组 ID listNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"phone_order": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Computed:    true,
													Description: "Phone polling listNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"phone_circle_times": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "数量 phone pollings. 取值范围：1-5Note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"phone_inner_interval": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Call 间隔 （秒） within 一个 polling. 取值范围：60-900Note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"phone_circle_interval": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Polling 间隔 （秒）。 取值范围：60-900Note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"need_phone_arrive_notice": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Whether receipt 通知 为必填项. 有效值：0 (无)，1 (yes)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"phone_call_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Dial 类型 SYNC (simultaneous dial)，CIRCLE (polled dial). 默认值：CIRCLE.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"weekday": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Computed:    true,
													Description: "Notification cycle. 值 1-7 indicate Monday 到 Sunday.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"on_call_form_ids": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "列表 调度 IDsNote: u200dThis 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
											},
										},
									},
									"url_notices": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Callback 通知 listNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Callback URL，其中 可以 contain up 到 256 charactersNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"is_valid": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Whether verification 是 passed. 有效值：0 (无)，1 (yes)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"validation_code": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Verification codeNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"start_time": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "开始时间 的 通知 （秒）， 其中 是 calculated 从 00:00:00.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"end_time": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "结束时间 的 通知 （秒）， 其中 是 calculated 从 00:00:00.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"weekday": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Computed:    true,
													Description: "Notification cycle. 值 1-7 indicate Monday 到 Sunday.注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"is_preset": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "是否为the 系统 默认值 通知 template. 有效值：0 (无)，1 (yes)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
									},
									"notice_language": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Notification 语言 有效值：zh-CN (Chinese)，en-US (English)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
									},
									"policy_ids": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "列表 IDs 的 告警 policies bound 到 告警 通知 templateNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
									},
									"amp_consumer_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Backend AMP 消费者 ID.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"cls_notices": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Channel 到 push 告警 notifications 到 CLS.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "地域",
												},
												"log_set_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Logset ID。",
												},
												"topic_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Topic ID。",
												},
												"enable": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "状态 有效值：0 (已禁用)，1 (已启用). 默认值：1 (已启用). 此 参数 可以 是 left 空。",
												},
											},
										},
									},
									"tags": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "标签 bound 到 通知 template注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "标签键",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "标签值",
												},
											},
										},
									},
								},
							},
						},
						"trigger_tasks": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Triggered 任务 listNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Triggered 任务 类型 有效 值: AS (auto scaling)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
									},
									"task_config": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration 信息 在 JSON 格式，such 作为 {Key1:Value1,Key2:Value2}注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
									},
								},
							},
						},
						"conditions_temp": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "模板 策略 groupNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"template_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "模板 nameNote: u200dThis 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
									},
									"condition": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Metric 触发器 conditionNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_union_rule": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Judgment condition 的 告警 触发器 condition (0: Any; 1: All; 2: Composite). 当 值 是 集合 到 2 (i.e.，composite 触发器 conditions)，此 参数 should 是 使用 together 使用 ComplexExpression.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"rules": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Alarm 触发器 condition listNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"metric_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "指标名称 或 事件名称 支持 metrics 可以 是 queried via DescribeAlarmMetrics 和 支持 events via DescribeAlarmEvents.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
															},
															"period": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Statistical 周期 （秒）。 有效 值 可以 是 queried via DescribeAlarmMetrics.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
															},
															"operator": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Operatorintelligent = intelligent detection without thresholdeq = equal toge = greater 比 或 equal togt = greater thanle = less 比 或 equal tolt = less thanne = 不 equal today_increase = day-在-day increaseday_decrease = day-在-day decreaseday_wave = day-在-day fluctuationweek_increase = week-在-week increaseweek_decrease = week-在-week decreaseweek_wave = week-在-week fluctuationcycle_increase = cyclical increasecycle_decrease = cyclical decreasecycle_wave = cyclical fluctuationre = regex matchThe 有效 值 可以 是 queried via DescribeAlarmMetrics.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
															},
															"value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Threshold. 有效 值 范围 可以 是 queried via DescribeAlarmMetrics.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
															},
															"continue_period": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "数量 periods. 1: continue 对于 一个 周期; 2: continue 对于 two periods; 和 so 在. 有效 值 可以 是 queried via DescribeAlarmMetrics.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
															},
															"notice_frequency": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Alarm 间隔 （秒）。 有效值：0 (do 不 repeat)，300 (告警 once every 5 minutes)，600 (告警 once every 10 minutes)，900 (告警 once every 15 minutes)，1800 (告警 once every 30 minutes)，3600 (告警 once every hour)，7200 (告警 once every 2 hours)，10800 (告警 once every 3 hours)，21600 (告警 once every 6 hours)，43200 (告警 once every 12 hours)，86400 (告警 once every day)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
															},
															"is_power_notice": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "是否alarm 频率 increases exponentially. 有效值：0 (无)，1 (yes)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
															},
															"filter": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "过滤器 condition 对于 一个 单个 触发器 ruleNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "过滤器 条件类型 有效值：DIMENSION (uses dimensions 对于 filtering)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
																		},
																		"dimensions": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "JSON 字符串 generated 通过 serializing AlarmPolicyDimension two-dimensional 数组. 一个-dimensional arrays 是 在 OR relationship，和 elements 在 一个-dimensional 数组 是 在 AND relationshipNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
																		},
																	},
																},
															},
															"description": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Metric display 名称，其中 是 使用 在 output parameterNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
															},
															"unit": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Unit，其中 是 使用 在 output parameterNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
															},
															"rule_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Trigger 条件类型 STATIC: 静态 阈值; 动态: 动态 阈值. 如果 您 do 不 指定this 参数 当 creating 或 editing 策略，STATIC 是 使用 通过 默认值.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
															},
															"is_advanced": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "是否为an advanced metric. 0: No; 1: Yes.注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"is_open": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "是否advanced metric 功能 是 已启用 0: No; 1: Yes.注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"product_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Integration center product ID.注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"value_max": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Maximum value注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"value_min": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Minimum value注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"hierarchical_value": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "配置 的 告警 级别 threshold注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"remind": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Threshold 对于 Remind level注意：此字段可能返回 null，表示无法获取有效值。",
																		},
																		"warn": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Threshold 对于 Warn level注意：此字段可能返回 null，表示无法获取有效值。",
																		},
																		"serious": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Threshold 对于 Serious level注意：此字段可能返回 null，表示无法获取有效值。",
																		},
																	},
																},
															},
														},
													},
												},
												"complex_expression": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "judgment expression 的 composite 告警 触发器 conditions，其中 是 有效 当 值 的 IsUnionRule 是 2. 此 参数 是 用于determine 该 告警 condition 是 met 仅 当 expression 值 是 True 对于 多个 触发器 conditions.注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"event_condition": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Event 触发器 conditionNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rules": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Alarm 触发器 condition listNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"metric_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "指标名称 或 事件名称 支持 metrics 可以 是 queried via DescribeAlarmMetrics 和 支持 events via DescribeAlarmEvents.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
															},
															"period": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Statistical 周期 （秒）。 有效 值 可以 是 queried via DescribeAlarmMetrics.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
															},
															"operator": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Operatorintelligent = intelligent detection without thresholdeq = equal toge = greater 比 或 equal togt = greater thanle = less 比 或 equal tolt = less thanne = 不 equal today_increase = day-在-day increaseday_decrease = day-在-day decreaseday_wave = day-在-day fluctuationweek_increase = week-在-week increaseweek_decrease = week-在-week decreaseweek_wave = week-在-week fluctuationcycle_increase = cyclical increasecycle_decrease = cyclical decreasecycle_wave = cyclical fluctuationre = regex matchThe 有效 值 可以 是 queried via DescribeAlarmMetrics.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
															},
															"value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Threshold. 有效 值 范围 可以 是 queried via DescribeAlarmMetrics.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
															},
															"continue_period": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "数量 periods. 1: continue 对于 一个 周期; 2: continue 对于 two periods; 和 so 在. 有效 值 可以 是 queried via DescribeAlarmMetrics.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
															},
															"notice_frequency": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Alarm 间隔 （秒）。 有效值：0 (do 不 repeat)，300 (告警 once every 5 minutes)，600 (告警 once every 10 minutes)，900 (告警 once every 15 minutes)，1800 (告警 once every 30 minutes)，3600 (告警 once every hour)，7200 (告警 once every 2 hours)，10800 (告警 once every 3 hours)，21600 (告警 once every 6 hours)，43200 (告警 once every 12 hours)，86400 (告警 once every day)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
															},
															"is_power_notice": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "是否alarm 频率 increases exponentially. 有效值：0 (无)，1 (yes)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
															},
															"filter": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "过滤器 condition 对于 一个 单个 触发器 ruleNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "过滤器 条件类型 有效值：DIMENSION (uses dimensions 对于 filtering)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
																		},
																		"dimensions": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "JSON 字符串 generated 通过 serializing AlarmPolicyDimension two-dimensional 数组. 一个-dimensional arrays 是 在 OR relationship，和 elements 在 一个-dimensional 数组 是 在 AND relationshipNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
																		},
																	},
																},
															},
															"description": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Metric display 名称，其中 是 使用 在 output parameterNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
															},
															"unit": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Unit，其中 是 使用 在 output parameterNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
															},
															"rule_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Trigger 条件类型 STATIC: 静态 阈值; 动态: 动态 阈值. 如果 您 do 不 指定this 参数 当 creating 或 editing 策略，STATIC 是 使用 通过 默认值.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
															},
															"is_advanced": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "是否为an advanced metric. 0: No; 1: Yes.注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"is_open": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "是否advanced metric 功能 是 已启用 0: No; 1: Yes.注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"product_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Integration center product ID.注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"value_max": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Maximum value注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"value_min": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Minimum value注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"hierarchical_value": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "配置 的 告警 级别 threshold注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"remind": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Threshold 对于 Remind level注意：此字段可能返回 null，表示无法获取有效值。",
																		},
																		"warn": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Threshold 对于 Warn level注意：此字段可能返回 null，表示无法获取有效值。",
																		},
																		"serious": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Threshold 对于 Serious level注意：此字段可能返回 null，表示无法获取有效值。",
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
						"last_edit_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Uin 的 last modifying userNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Update timeNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"insert_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Creation timeNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"region": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "RegionNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"namespace_show_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace display nameNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"is_default": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否为the 默认值 策略. 有效值：1 (yes)，0 (无)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"can_set_default": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否default 策略 可以 是 集合. 有效值：1 (yes)，0 (无)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"instance_group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例 组 IDNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"instance_sum": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total 数量 实例 在 实例 groupNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"instance_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 组 nameNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"rule_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Trigger 条件类型 有效值：STATIC (静态 阈值)，DYNAMIC (动态)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"origin_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy ID 对于 实例/实例 组 binding 和 unbinding APIs (BindingPolicyObject，UnBindingAllPolicyObject，UnBindingPolicyObject)注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"tag_instances": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Tag注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签 key注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签 value注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"instance_sum": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 instances注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"service_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service 类型，对于 示例，CVM注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"region_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "地域 ID注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"binding_status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Binding 状态 2: bound; 1: binding注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"tag_status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "标签 状态 2: existent; 1: nonexistent注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"filter_dimensions_param": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Information 在 过滤器 dimension associated 使用 策略.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"is_one_click": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否为a quick 告警 策略.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"one_click_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否quick 告警 策略 是 已启用注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"advanced_metric_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 advanced metrics.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"is_bind_all": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否policy 是 associated 使用 all objects注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Policy tag注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签键",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签值",
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

func dataSourceTencentCloudMonitorAlarmPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_alarm_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("module"); ok {
		paramMap["Module"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("policy_name"); ok {
		paramMap["PolicyName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("monitor_types"); ok {
		monitorTypesSet := v.(*schema.Set).List()
		paramMap["MonitorTypes"] = helper.InterfacesStringsPoint(monitorTypesSet)
	}

	if v, ok := d.GetOk("namespaces"); ok {
		namespacesSet := v.(*schema.Set).List()
		paramMap["Namespaces"] = helper.InterfacesStringsPoint(namespacesSet)
	}

	if v, ok := d.GetOk("dimensions"); ok {
		paramMap["Dimensions"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("receiver_uids"); ok {
		receiverUidsList := []*int64{}
		receiverUidsSet := v.(*schema.Set).List()
		for i := range receiverUidsSet {
			receiverUids := receiverUidsSet[i].(int)
			receiverUidsList = append(receiverUidsList, helper.IntInt64(receiverUids))
		}
		paramMap["ReceiverUids"] = receiverUidsList
	}

	if v, ok := d.GetOk("receiver_groups"); ok {
		receiverGroupsList := []*int64{}
		receiverGroupsSet := v.(*schema.Set).List()
		for i := range receiverGroupsSet {
			receiverGroups := receiverGroupsSet[i].(int)
			receiverGroupsList = append(receiverGroupsList, helper.IntInt64(receiverGroups))
		}
		paramMap["ReceiverGroups"] = receiverGroupsList
	}

	if v, ok := d.GetOk("policy_type"); ok {
		policyTypeSet := v.(*schema.Set).List()
		paramMap["PolicyType"] = helper.InterfacesStringsPoint(policyTypeSet)
	}

	if v, ok := d.GetOk("field"); ok {
		paramMap["Field"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_ids"); ok {
		projectIdsList := []*int64{}
		projectIdsSet := v.(*schema.Set).List()
		for i := range projectIdsSet {
			projectIds := projectIdsSet[i].(int)
			projectIdsList = append(projectIdsList, helper.IntInt64(projectIds))
		}
		paramMap["ProjectIds"] = projectIdsList
	}

	if v, ok := d.GetOk("notice_ids"); ok {
		noticeIdsSet := v.(*schema.Set).List()
		paramMap["NoticeIds"] = helper.InterfacesStringsPoint(noticeIdsSet)
	}

	if v, ok := d.GetOk("rule_types"); ok {
		ruleTypesSet := v.(*schema.Set).List()
		paramMap["RuleTypes"] = helper.InterfacesStringsPoint(ruleTypesSet)
	}

	if v, ok := d.GetOk("enable"); ok {
		enableList := []*int64{}
		enableSet := v.(*schema.Set).List()
		for i := range enableSet {
			enable := enableSet[i].(int)
			enableList = append(enableList, helper.IntInt64(enable))
		}
		paramMap["Enable"] = enableList
	}

	if v, ok := d.GetOkExists("not_binding_notice_rule"); ok {
		paramMap["NotBindingNoticeRule"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("instance_group_id"); ok {
		paramMap["InstanceGroupId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("need_correspondence"); ok {
		paramMap["NeedCorrespondence"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("trigger_tasks"); ok {
		triggerTasksSet := v.([]interface{})
		tmpSet := make([]*monitor.AlarmPolicyTriggerTask, 0, len(triggerTasksSet))

		for _, item := range triggerTasksSet {
			alarmPolicyTriggerTask := monitor.AlarmPolicyTriggerTask{}
			alarmPolicyTriggerTaskMap := item.(map[string]interface{})

			if v, ok := alarmPolicyTriggerTaskMap["type"]; ok {
				alarmPolicyTriggerTask.Type = helper.String(v.(string))
			}
			if v, ok := alarmPolicyTriggerTaskMap["task_config"]; ok {
				alarmPolicyTriggerTask.TaskConfig = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &alarmPolicyTriggerTask)
		}
		paramMap["trigger_tasks"] = tmpSet
	}

	if v, ok := d.GetOk("one_click_policy_type"); ok {
		oneClickPolicyTypeSet := v.(*schema.Set).List()
		paramMap["OneClickPolicyType"] = helper.InterfacesStringsPoint(oneClickPolicyTypeSet)
	}

	if v, ok := d.GetOkExists("not_bind_all"); ok {
		paramMap["NotBindAll"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("not_instance_group"); ok {
		paramMap["NotInstanceGroup"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("prom_ins_id"); ok {
		paramMap["PromInsId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("receiver_on_call_form_ids"); ok {
		receiverOnCallFormIDsSet := v.(*schema.Set).List()
		paramMap["ReceiverOnCallFormIDs"] = helper.InterfacesStringsPoint(receiverOnCallFormIDsSet)
	}

	service := MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var policies []*monitor.AlarmPolicy
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMonitorAlarmPolicyByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		policies = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(policies))
	tmpList := make([]map[string]interface{}, 0, len(policies))

	if policies != nil {
		for _, alarmPolicy := range policies {
			alarmPolicyMap := map[string]interface{}{}

			if alarmPolicy.PolicyId != nil {
				alarmPolicyMap["policy_id"] = alarmPolicy.PolicyId
			}

			if alarmPolicy.PolicyName != nil {
				alarmPolicyMap["policy_name"] = alarmPolicy.PolicyName
			}

			if alarmPolicy.Remark != nil {
				alarmPolicyMap["remark"] = alarmPolicy.Remark
			}

			if alarmPolicy.MonitorType != nil {
				alarmPolicyMap["monitor_type"] = alarmPolicy.MonitorType
			}

			if alarmPolicy.Enable != nil {
				alarmPolicyMap["enable"] = alarmPolicy.Enable
			}

			if alarmPolicy.UseSum != nil {
				alarmPolicyMap["use_sum"] = alarmPolicy.UseSum
			}

			if alarmPolicy.ProjectId != nil {
				alarmPolicyMap["project_id"] = alarmPolicy.ProjectId
			}

			if alarmPolicy.ProjectName != nil {
				alarmPolicyMap["project_name"] = alarmPolicy.ProjectName
			}

			if alarmPolicy.Namespace != nil {
				alarmPolicyMap["namespace"] = alarmPolicy.Namespace
			}

			if alarmPolicy.ConditionTemplateId != nil {
				alarmPolicyMap["condition_template_id"] = alarmPolicy.ConditionTemplateId
			}

			if alarmPolicy.Condition != nil {
				conditionMap := map[string]interface{}{}

				if alarmPolicy.Condition.IsUnionRule != nil {
					conditionMap["is_union_rule"] = alarmPolicy.Condition.IsUnionRule
				}

				if alarmPolicy.Condition.Rules != nil {
					rulesList := []interface{}{}
					for _, rules := range alarmPolicy.Condition.Rules {
						rulesMap := map[string]interface{}{}

						if rules.MetricName != nil {
							rulesMap["metric_name"] = rules.MetricName
						}

						if rules.Period != nil {
							rulesMap["period"] = rules.Period
						}

						if rules.Operator != nil {
							rulesMap["operator"] = rules.Operator
						}

						if rules.Value != nil {
							rulesMap["value"] = rules.Value
						}

						if rules.ContinuePeriod != nil {
							rulesMap["continue_period"] = rules.ContinuePeriod
						}

						if rules.NoticeFrequency != nil {
							rulesMap["notice_frequency"] = rules.NoticeFrequency
						}

						if rules.IsPowerNotice != nil {
							rulesMap["is_power_notice"] = rules.IsPowerNotice
						}

						if rules.Filter != nil {
							filterMap := map[string]interface{}{}

							if rules.Filter.Type != nil {
								filterMap["type"] = rules.Filter.Type
							}

							if rules.Filter.Dimensions != nil {
								filterMap["dimensions"] = rules.Filter.Dimensions
							}

							rulesMap["filter"] = []interface{}{filterMap}
						}

						if rules.Description != nil {
							rulesMap["description"] = rules.Description
						}

						if rules.Unit != nil {
							rulesMap["unit"] = rules.Unit
						}

						if rules.RuleType != nil {
							rulesMap["rule_type"] = rules.RuleType
						}

						if rules.IsAdvanced != nil {
							rulesMap["is_advanced"] = rules.IsAdvanced
						}

						if rules.IsOpen != nil {
							rulesMap["is_open"] = rules.IsOpen
						}

						if rules.ProductId != nil {
							rulesMap["product_id"] = rules.ProductId
						}

						if rules.ValueMax != nil {
							rulesMap["value_max"] = rules.ValueMax
						}

						if rules.ValueMin != nil {
							rulesMap["value_min"] = rules.ValueMin
						}

						if rules.HierarchicalValue != nil {
							hierarchicalValueMap := map[string]interface{}{}

							if rules.HierarchicalValue.Remind != nil {
								hierarchicalValueMap["remind"] = rules.HierarchicalValue.Remind
							}

							if rules.HierarchicalValue.Warn != nil {
								hierarchicalValueMap["warn"] = rules.HierarchicalValue.Warn
							}

							if rules.HierarchicalValue.Serious != nil {
								hierarchicalValueMap["serious"] = rules.HierarchicalValue.Serious
							}

							rulesMap["hierarchical_value"] = []interface{}{hierarchicalValueMap}
						}

						rulesList = append(rulesList, rulesMap)
					}

					conditionMap["rules"] = rulesList
				}

				if alarmPolicy.Condition.ComplexExpression != nil {
					conditionMap["complex_expression"] = alarmPolicy.Condition.ComplexExpression
				}

				alarmPolicyMap["condition"] = []interface{}{conditionMap}
			}

			if alarmPolicy.EventCondition != nil {
				eventConditionMap := map[string]interface{}{}

				if alarmPolicy.EventCondition.Rules != nil {
					rulesList := []interface{}{}
					for _, rules := range alarmPolicy.EventCondition.Rules {
						rulesMap := map[string]interface{}{}

						if rules.MetricName != nil {
							rulesMap["metric_name"] = rules.MetricName
						}

						if rules.Period != nil {
							rulesMap["period"] = rules.Period
						}

						if rules.Operator != nil {
							rulesMap["operator"] = rules.Operator
						}

						if rules.Value != nil {
							rulesMap["value"] = rules.Value
						}

						if rules.ContinuePeriod != nil {
							rulesMap["continue_period"] = rules.ContinuePeriod
						}

						if rules.NoticeFrequency != nil {
							rulesMap["notice_frequency"] = rules.NoticeFrequency
						}

						if rules.IsPowerNotice != nil {
							rulesMap["is_power_notice"] = rules.IsPowerNotice
						}

						if rules.Filter != nil {
							filterMap := map[string]interface{}{}

							if rules.Filter.Type != nil {
								filterMap["type"] = rules.Filter.Type
							}

							if rules.Filter.Dimensions != nil {
								filterMap["dimensions"] = rules.Filter.Dimensions
							}

							rulesMap["filter"] = []interface{}{filterMap}
						}

						if rules.Description != nil {
							rulesMap["description"] = rules.Description
						}

						if rules.Unit != nil {
							rulesMap["unit"] = rules.Unit
						}

						if rules.RuleType != nil {
							rulesMap["rule_type"] = rules.RuleType
						}

						if rules.IsAdvanced != nil {
							rulesMap["is_advanced"] = rules.IsAdvanced
						}

						if rules.IsOpen != nil {
							rulesMap["is_open"] = rules.IsOpen
						}

						if rules.ProductId != nil {
							rulesMap["product_id"] = rules.ProductId
						}

						if rules.ValueMax != nil {
							rulesMap["value_max"] = rules.ValueMax
						}

						if rules.ValueMin != nil {
							rulesMap["value_min"] = rules.ValueMin
						}

						if rules.HierarchicalValue != nil {
							hierarchicalValueMap := map[string]interface{}{}

							if rules.HierarchicalValue.Remind != nil {
								hierarchicalValueMap["remind"] = rules.HierarchicalValue.Remind
							}

							if rules.HierarchicalValue.Warn != nil {
								hierarchicalValueMap["warn"] = rules.HierarchicalValue.Warn
							}

							if rules.HierarchicalValue.Serious != nil {
								hierarchicalValueMap["serious"] = rules.HierarchicalValue.Serious
							}

							rulesMap["hierarchical_value"] = []interface{}{hierarchicalValueMap}
						}

						rulesList = append(rulesList, rulesMap)
					}

					eventConditionMap["rules"] = rulesList
				}

				alarmPolicyMap["event_condition"] = []interface{}{eventConditionMap}
			}

			if alarmPolicy.NoticeIds != nil {
				alarmPolicyMap["notice_ids"] = alarmPolicy.NoticeIds
			}

			if alarmPolicy.Notices != nil {
				noticesList := []interface{}{}
				for _, notices := range alarmPolicy.Notices {
					noticesMap := map[string]interface{}{}

					if notices.Id != nil {
						noticesMap["id"] = notices.Id
					}

					if notices.Name != nil {
						noticesMap["name"] = notices.Name
					}

					if notices.UpdatedAt != nil {
						noticesMap["updated_at"] = notices.UpdatedAt
					}

					if notices.UpdatedBy != nil {
						noticesMap["updated_by"] = notices.UpdatedBy
					}

					if notices.NoticeType != nil {
						noticesMap["notice_type"] = notices.NoticeType
					}

					if notices.UserNotices != nil {
						userNoticesList := []interface{}{}
						for _, userNotices := range notices.UserNotices {
							userNoticesMap := map[string]interface{}{}

							if userNotices.ReceiverType != nil {
								userNoticesMap["receiver_type"] = userNotices.ReceiverType
							}

							if userNotices.StartTime != nil {
								userNoticesMap["start_time"] = userNotices.StartTime
							}

							if userNotices.EndTime != nil {
								userNoticesMap["end_time"] = userNotices.EndTime
							}

							if userNotices.NoticeWay != nil {
								userNoticesMap["notice_way"] = userNotices.NoticeWay
							}

							if userNotices.UserIds != nil {
								userNoticesMap["user_ids"] = userNotices.UserIds
							}

							if userNotices.GroupIds != nil {
								userNoticesMap["group_ids"] = userNotices.GroupIds
							}

							if userNotices.PhoneOrder != nil {
								userNoticesMap["phone_order"] = userNotices.PhoneOrder
							}

							if userNotices.PhoneCircleTimes != nil {
								userNoticesMap["phone_circle_times"] = userNotices.PhoneCircleTimes
							}

							if userNotices.PhoneInnerInterval != nil {
								userNoticesMap["phone_inner_interval"] = userNotices.PhoneInnerInterval
							}

							if userNotices.PhoneCircleInterval != nil {
								userNoticesMap["phone_circle_interval"] = userNotices.PhoneCircleInterval
							}

							if userNotices.NeedPhoneArriveNotice != nil {
								userNoticesMap["need_phone_arrive_notice"] = userNotices.NeedPhoneArriveNotice
							}

							if userNotices.PhoneCallType != nil {
								userNoticesMap["phone_call_type"] = userNotices.PhoneCallType
							}

							if userNotices.Weekday != nil {
								userNoticesMap["weekday"] = userNotices.Weekday
							}

							if userNotices.OnCallFormIDs != nil {
								userNoticesMap["on_call_form_ids"] = userNotices.OnCallFormIDs
							}

							userNoticesList = append(userNoticesList, userNoticesMap)
						}

						noticesMap["user_notices"] = userNoticesList
					}

					if notices.URLNotices != nil {
						uRLNoticesList := []interface{}{}
						for _, uRLNotices := range notices.URLNotices {
							uRLNoticesMap := map[string]interface{}{}

							if uRLNotices.URL != nil {
								uRLNoticesMap["url"] = uRLNotices.URL
							}

							if uRLNotices.IsValid != nil {
								uRLNoticesMap["is_valid"] = uRLNotices.IsValid
							}

							if uRLNotices.ValidationCode != nil {
								uRLNoticesMap["validation_code"] = uRLNotices.ValidationCode
							}

							if uRLNotices.StartTime != nil {
								uRLNoticesMap["start_time"] = uRLNotices.StartTime
							}

							if uRLNotices.EndTime != nil {
								uRLNoticesMap["end_time"] = uRLNotices.EndTime
							}

							if uRLNotices.Weekday != nil {
								uRLNoticesMap["weekday"] = uRLNotices.Weekday
							}

							uRLNoticesList = append(uRLNoticesList, uRLNoticesMap)
						}

						noticesMap["url_notices"] = uRLNoticesList
					}

					if notices.IsPreset != nil {
						noticesMap["is_preset"] = notices.IsPreset
					}

					if notices.NoticeLanguage != nil {
						noticesMap["notice_language"] = notices.NoticeLanguage
					}

					if notices.PolicyIds != nil {
						noticesMap["policy_ids"] = notices.PolicyIds
					}

					if notices.AMPConsumerId != nil {
						noticesMap["amp_consumer_id"] = notices.AMPConsumerId
					}

					if notices.CLSNotices != nil {
						cLSNoticesList := []interface{}{}
						for _, cLSNotices := range notices.CLSNotices {
							cLSNoticesMap := map[string]interface{}{}

							if cLSNotices.Region != nil {
								cLSNoticesMap["region"] = cLSNotices.Region
							}

							if cLSNotices.LogSetId != nil {
								cLSNoticesMap["log_set_id"] = cLSNotices.LogSetId
							}

							if cLSNotices.TopicId != nil {
								cLSNoticesMap["topic_id"] = cLSNotices.TopicId
							}

							if cLSNotices.Enable != nil {
								cLSNoticesMap["enable"] = cLSNotices.Enable
							}

							cLSNoticesList = append(cLSNoticesList, cLSNoticesMap)
						}

						noticesMap["cls_notices"] = cLSNoticesList
					}

					if notices.Tags != nil {
						tagsList := []interface{}{}
						for _, tags := range notices.Tags {
							tagsMap := map[string]interface{}{}

							if tags.Key != nil {
								tagsMap["key"] = tags.Key
							}

							if tags.Value != nil {
								tagsMap["value"] = tags.Value
							}

							tagsList = append(tagsList, tagsMap)
						}

						noticesMap["tags"] = tagsList
					}

					noticesList = append(noticesList, noticesMap)
				}

				alarmPolicyMap["notices"] = noticesList
			}

			if alarmPolicy.TriggerTasks != nil {
				triggerTasksList := []interface{}{}
				for _, triggerTasks := range alarmPolicy.TriggerTasks {
					triggerTasksMap := map[string]interface{}{}

					if triggerTasks.Type != nil {
						triggerTasksMap["type"] = triggerTasks.Type
					}

					if triggerTasks.TaskConfig != nil {
						triggerTasksMap["task_config"] = triggerTasks.TaskConfig
					}

					triggerTasksList = append(triggerTasksList, triggerTasksMap)
				}

				alarmPolicyMap["trigger_tasks"] = triggerTasksList
			}

			if alarmPolicy.ConditionsTemp != nil {
				conditionsTempMap := map[string]interface{}{}

				if alarmPolicy.ConditionsTemp.TemplateName != nil {
					conditionsTempMap["template_name"] = alarmPolicy.ConditionsTemp.TemplateName
				}

				if alarmPolicy.ConditionsTemp.Condition != nil {
					conditionMap := map[string]interface{}{}

					if alarmPolicy.ConditionsTemp.Condition.IsUnionRule != nil {
						conditionMap["is_union_rule"] = alarmPolicy.ConditionsTemp.Condition.IsUnionRule
					}

					if alarmPolicy.ConditionsTemp.Condition.Rules != nil {
						rulesList := []interface{}{}
						for _, rules := range alarmPolicy.ConditionsTemp.Condition.Rules {
							rulesMap := map[string]interface{}{}

							if rules.MetricName != nil {
								rulesMap["metric_name"] = rules.MetricName
							}

							if rules.Period != nil {
								rulesMap["period"] = rules.Period
							}

							if rules.Operator != nil {
								rulesMap["operator"] = rules.Operator
							}

							if rules.Value != nil {
								rulesMap["value"] = rules.Value
							}

							if rules.ContinuePeriod != nil {
								rulesMap["continue_period"] = rules.ContinuePeriod
							}

							if rules.NoticeFrequency != nil {
								rulesMap["notice_frequency"] = rules.NoticeFrequency
							}

							if rules.IsPowerNotice != nil {
								rulesMap["is_power_notice"] = rules.IsPowerNotice
							}

							if rules.Filter != nil {
								filterMap := map[string]interface{}{}

								if rules.Filter.Type != nil {
									filterMap["type"] = rules.Filter.Type
								}

								if rules.Filter.Dimensions != nil {
									filterMap["dimensions"] = rules.Filter.Dimensions
								}

								rulesMap["filter"] = []interface{}{filterMap}
							}

							if rules.Description != nil {
								rulesMap["description"] = rules.Description
							}

							if rules.Unit != nil {
								rulesMap["unit"] = rules.Unit
							}

							if rules.RuleType != nil {
								rulesMap["rule_type"] = rules.RuleType
							}

							if rules.IsAdvanced != nil {
								rulesMap["is_advanced"] = rules.IsAdvanced
							}

							if rules.IsOpen != nil {
								rulesMap["is_open"] = rules.IsOpen
							}

							if rules.ProductId != nil {
								rulesMap["product_id"] = rules.ProductId
							}

							if rules.ValueMax != nil {
								rulesMap["value_max"] = rules.ValueMax
							}

							if rules.ValueMin != nil {
								rulesMap["value_min"] = rules.ValueMin
							}

							if rules.HierarchicalValue != nil {
								hierarchicalValueMap := map[string]interface{}{}

								if rules.HierarchicalValue.Remind != nil {
									hierarchicalValueMap["remind"] = rules.HierarchicalValue.Remind
								}

								if rules.HierarchicalValue.Warn != nil {
									hierarchicalValueMap["warn"] = rules.HierarchicalValue.Warn
								}

								if rules.HierarchicalValue.Serious != nil {
									hierarchicalValueMap["serious"] = rules.HierarchicalValue.Serious
								}

								rulesMap["hierarchical_value"] = []interface{}{hierarchicalValueMap}
							}

							rulesList = append(rulesList, rulesMap)
						}

						conditionMap["rules"] = rulesList
					}

					if alarmPolicy.ConditionsTemp.Condition.ComplexExpression != nil {
						conditionMap["complex_expression"] = alarmPolicy.ConditionsTemp.Condition.ComplexExpression
					}

					conditionsTempMap["condition"] = []interface{}{conditionMap}
				}

				if alarmPolicy.ConditionsTemp.EventCondition != nil {
					eventConditionMap := map[string]interface{}{}

					if alarmPolicy.ConditionsTemp.EventCondition.Rules != nil {
						rulesList := []interface{}{}
						for _, rules := range alarmPolicy.ConditionsTemp.EventCondition.Rules {
							rulesMap := map[string]interface{}{}

							if rules.MetricName != nil {
								rulesMap["metric_name"] = rules.MetricName
							}

							if rules.Period != nil {
								rulesMap["period"] = rules.Period
							}

							if rules.Operator != nil {
								rulesMap["operator"] = rules.Operator
							}

							if rules.Value != nil {
								rulesMap["value"] = rules.Value
							}

							if rules.ContinuePeriod != nil {
								rulesMap["continue_period"] = rules.ContinuePeriod
							}

							if rules.NoticeFrequency != nil {
								rulesMap["notice_frequency"] = rules.NoticeFrequency
							}

							if rules.IsPowerNotice != nil {
								rulesMap["is_power_notice"] = rules.IsPowerNotice
							}

							if rules.Filter != nil {
								filterMap := map[string]interface{}{}

								if rules.Filter.Type != nil {
									filterMap["type"] = rules.Filter.Type
								}

								if rules.Filter.Dimensions != nil {
									filterMap["dimensions"] = rules.Filter.Dimensions
								}

								rulesMap["filter"] = []interface{}{filterMap}
							}

							if rules.Description != nil {
								rulesMap["description"] = rules.Description
							}

							if rules.Unit != nil {
								rulesMap["unit"] = rules.Unit
							}

							if rules.RuleType != nil {
								rulesMap["rule_type"] = rules.RuleType
							}

							if rules.IsAdvanced != nil {
								rulesMap["is_advanced"] = rules.IsAdvanced
							}

							if rules.IsOpen != nil {
								rulesMap["is_open"] = rules.IsOpen
							}

							if rules.ProductId != nil {
								rulesMap["product_id"] = rules.ProductId
							}

							if rules.ValueMax != nil {
								rulesMap["value_max"] = rules.ValueMax
							}

							if rules.ValueMin != nil {
								rulesMap["value_min"] = rules.ValueMin
							}

							if rules.HierarchicalValue != nil {
								hierarchicalValueMap := map[string]interface{}{}

								if rules.HierarchicalValue.Remind != nil {
									hierarchicalValueMap["remind"] = rules.HierarchicalValue.Remind
								}

								if rules.HierarchicalValue.Warn != nil {
									hierarchicalValueMap["warn"] = rules.HierarchicalValue.Warn
								}

								if rules.HierarchicalValue.Serious != nil {
									hierarchicalValueMap["serious"] = rules.HierarchicalValue.Serious
								}

								rulesMap["hierarchical_value"] = []interface{}{hierarchicalValueMap}
							}

							rulesList = append(rulesList, rulesMap)
						}

						eventConditionMap["rules"] = rulesList
					}

					conditionsTempMap["event_condition"] = []interface{}{eventConditionMap}
				}

				alarmPolicyMap["conditions_temp"] = []interface{}{conditionsTempMap}
			}

			if alarmPolicy.LastEditUin != nil {
				alarmPolicyMap["last_edit_uin"] = alarmPolicy.LastEditUin
			}

			if alarmPolicy.UpdateTime != nil {
				alarmPolicyMap["update_time"] = alarmPolicy.UpdateTime
			}

			if alarmPolicy.InsertTime != nil {
				alarmPolicyMap["insert_time"] = alarmPolicy.InsertTime
			}

			if alarmPolicy.Region != nil {
				alarmPolicyMap["region"] = alarmPolicy.Region
			}

			if alarmPolicy.NamespaceShowName != nil {
				alarmPolicyMap["namespace_show_name"] = alarmPolicy.NamespaceShowName
			}

			if alarmPolicy.IsDefault != nil {
				alarmPolicyMap["is_default"] = alarmPolicy.IsDefault
			}

			if alarmPolicy.CanSetDefault != nil {
				alarmPolicyMap["can_set_default"] = alarmPolicy.CanSetDefault
			}

			if alarmPolicy.InstanceGroupId != nil {
				alarmPolicyMap["instance_group_id"] = alarmPolicy.InstanceGroupId
			}

			if alarmPolicy.InstanceSum != nil {
				alarmPolicyMap["instance_sum"] = alarmPolicy.InstanceSum
			}

			if alarmPolicy.InstanceGroupName != nil {
				alarmPolicyMap["instance_group_name"] = alarmPolicy.InstanceGroupName
			}

			if alarmPolicy.RuleType != nil {
				alarmPolicyMap["rule_type"] = alarmPolicy.RuleType
			}

			if alarmPolicy.OriginId != nil {
				alarmPolicyMap["origin_id"] = alarmPolicy.OriginId
			}

			if alarmPolicy.TagInstances != nil {
				tagInstancesList := []interface{}{}
				for _, tagInstances := range alarmPolicy.TagInstances {
					tagInstancesMap := map[string]interface{}{}

					if tagInstances.Key != nil {
						tagInstancesMap["key"] = tagInstances.Key
					}

					if tagInstances.Value != nil {
						tagInstancesMap["value"] = tagInstances.Value
					}

					if tagInstances.InstanceSum != nil {
						tagInstancesMap["instance_sum"] = tagInstances.InstanceSum
					}

					if tagInstances.ServiceType != nil {
						tagInstancesMap["service_type"] = tagInstances.ServiceType
					}

					if tagInstances.RegionId != nil {
						tagInstancesMap["region_id"] = tagInstances.RegionId
					}

					if tagInstances.BindingStatus != nil {
						tagInstancesMap["binding_status"] = tagInstances.BindingStatus
					}

					if tagInstances.TagStatus != nil {
						tagInstancesMap["tag_status"] = tagInstances.TagStatus
					}

					tagInstancesList = append(tagInstancesList, tagInstancesMap)
				}

				alarmPolicyMap["tag_instances"] = tagInstancesList
			}

			if alarmPolicy.FilterDimensionsParam != nil {
				alarmPolicyMap["filter_dimensions_param"] = alarmPolicy.FilterDimensionsParam
			}

			if alarmPolicy.IsOneClick != nil {
				alarmPolicyMap["is_one_click"] = alarmPolicy.IsOneClick
			}

			if alarmPolicy.OneClickStatus != nil {
				alarmPolicyMap["one_click_status"] = alarmPolicy.OneClickStatus
			}

			if alarmPolicy.AdvancedMetricNumber != nil {
				alarmPolicyMap["advanced_metric_number"] = alarmPolicy.AdvancedMetricNumber
			}

			if alarmPolicy.IsBindAll != nil {
				alarmPolicyMap["is_bind_all"] = alarmPolicy.IsBindAll
			}

			if alarmPolicy.Tags != nil {
				tagsList := []interface{}{}
				for _, tags := range alarmPolicy.Tags {
					tagsMap := map[string]interface{}{}

					if tags.Key != nil {
						tagsMap["key"] = tags.Key
					}

					if tags.Value != nil {
						tagsMap["value"] = tags.Value
					}

					tagsList = append(tagsList, tagsMap)
				}

				alarmPolicyMap["tags"] = tagsList
			}

			ids = append(ids, *alarmPolicy.PolicyId)
			tmpList = append(tmpList, alarmPolicyMap)
		}

		_ = d.Set("policies", tmpList)
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
