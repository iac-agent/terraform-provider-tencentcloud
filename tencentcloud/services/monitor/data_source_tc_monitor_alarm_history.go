package monitor

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMonitorAlarmHistory() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMonitorAlarmHistoryRead,
		Schema: map[string]*schema.Schema{
			"module": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "值 fixed 在 监控。",
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "排序方式 first occurrence 时间 在 descending 排序依据 默认值. 有效值：ASC (ascending)，DESC (descending)。",
			},

			"start_time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "开始时间，其中 是 时间戳 一个 day ago 通过 默认值 和 时间 当 告警 FirstOccurTime first occurs. An 告警 记录 可以 是 searched 仅 如果 its FirstOccurTime 是 later 比 StartTime。",
			},

			"end_time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "结束时间，其中 是 当前 时间戳 和 时间 当 告警 FirstOccurTime first occurs. An 告警 记录 可以 是 searched 仅 如果 its FirstOccurTime 是 earlier 比 EndTime。",
			},

			"monitor_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "过滤器 通过 监控类型 有效值：MT_QCE (Tencent Cloud 服务 监控)，MT_TAW (应用 performance 监控)，MT_RUM (frontend performance 监控)，MT_PROBE (云 automated testing). 如果此参数为空，all types 将 是 queried 通过 默认值。",
			},

			"alarm_object": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "过滤器 通过 告警 对象. Fuzzy search 使用 字符串 是 支持。",
			},

			"alarm_status": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "过滤器 通过 告警状态 有效值：ALARM (不 resolved)，OK (resolved)，NO_CONF (expired)，NO_DATA (insufficient 数据). 如果此参数为空，all 将 是 queried 通过 默认值。",
			},

			"project_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "过滤器 通过 项目 ID. 有效值：-1 (无 项目)，0 (默认值 项目)。",
			},

			"instance_group_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "过滤器 通过 实例 组 ID",
			},

			"namespaces": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 通过 策略 类型 Monitoring 类型 和 策略 类型 是 first-级别 和 second-级别 filters respectively 和 both need 到 是 passed 在. For 示例，[{MonitorType: MT_QCE，Namespace: cvm_device}]。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"monitor_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "监控类型",
						},
						"namespace": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Policy 类型",
						},
					},
				},
			},

			"metric_names": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "过滤器 通过 指标名称",
			},

			"policy_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Fuzzy search 通过 策略 名称",
			},

			"content": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Fuzzy search 通过 告警 内容",
			},

			"receiver_uids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Search 通过 recipient。",
			},

			"receiver_groups": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Search 通过 recipient 组。",
			},

			"policy_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Search 通过 告警策略 ID 列表。",
			},

			"alarm_levels": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Alarm levels。",
			},

			"histories": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Alarm 记录 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alarm_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm 记录 ID。",
						},
						"monitor_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "监控类型",
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy 类型",
						},
						"alarm_object": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm 对象。",
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm 内容",
						},
						"first_occur_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "时间戳 的 first occurrence。",
						},
						"last_occur_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "时间戳 的 last occurrence。",
						},
						"alarm_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "告警状态 有效值：ALARM (不 resolved)，OK (resolved)，NO_CONF (expired)，NO_DATA (insufficient 数据)。",
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "告警策略 ID",
						},
						"policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy 名称",
						},
						"vpc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC 的 告警 对象 对于 basic product 告警。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "项目 ID",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "项目名称",
						},
						"instance_group": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "实例 组 的 告警 对象。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "实例 组 ID",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例 组名称",
									},
								},
							},
						},
						"receiver_uids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "Recipient 列表。",
						},
						"receiver_groups": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "Recipient 组 列表。",
						},
						"notice_ways": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Alarm channel 列表. 有效值：SMS (SMS)，EMAIL (email)，CALL (phone)，WECHAT (WeChat)。",
						},
						"origin_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "告警策略 ID，其中 可以 是 使用 当 您 call APIs (BindingPolicyObject，UnBindingAllPolicyObject，UnBindingPolicyObject) 到 bind/unbind 实例 或 实例 groups 到/从 告警 策略。",
						},
						"alarm_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm 类型",
						},
						"event_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "事件 ID",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域",
						},
						"policy_exists": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否policy exists. 有效值：0 (无)，1 (yes)。",
						},
						"metrics_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Metric informationNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"qce_namespace": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace 用于query 数据 通过 Tencent Cloud 服务 监控 类型",
									},
									"metric_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "指标名称",
									},
									"period": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Statistical 周期",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "值 triggering 告警。",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Metric display 名称",
									},
								},
							},
						},
						"dimensions": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dimension 信息 的 实例 该 triggered alarms.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"alarm_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm 级别Note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
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

func dataSourceTencentCloudMonitorAlarmHistoryRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_alarm_history.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("module"); ok {
		paramMap["Module"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("start_time"); ok {
		paramMap["StartTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("end_time"); ok {
		paramMap["EndTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("monitor_types"); ok {
		monitorTypesSet := v.(*schema.Set).List()
		paramMap["MonitorTypes"] = helper.InterfacesStringsPoint(monitorTypesSet)
	}

	if v, ok := d.GetOk("alarm_object"); ok {
		paramMap["AlarmObject"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("alarm_status"); ok {
		alarmStatusSet := v.(*schema.Set).List()
		paramMap["AlarmStatus"] = helper.InterfacesStringsPoint(alarmStatusSet)
	}

	if v, ok := d.GetOk("project_ids"); ok {
		projectIdsSet := v.(*schema.Set).List()
		paramMap["ProjectIds"] = helper.InterfacesIntInt64Point(projectIdsSet)
	}

	if v, ok := d.GetOk("instance_group_ids"); ok {
		instanceGroupIdsSet := v.(*schema.Set).List()
		paramMap["InstanceGroupIds"] = helper.InterfacesIntInt64Point(instanceGroupIdsSet)
	}

	if v, ok := d.GetOk("namespaces"); ok {
		namespacesSet := v.([]interface{})
		tmpSet := make([]*monitor.MonitorTypeNamespace, 0, len(namespacesSet))

		for _, item := range namespacesSet {
			monitorTypeNamespace := monitor.MonitorTypeNamespace{}
			monitorTypeNamespaceMap := item.(map[string]interface{})

			if v, ok := monitorTypeNamespaceMap["monitor_type"]; ok {
				monitorTypeNamespace.MonitorType = helper.String(v.(string))
			}
			if v, ok := monitorTypeNamespaceMap["namespace"]; ok {
				monitorTypeNamespace.Namespace = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &monitorTypeNamespace)
		}
		paramMap["namespaces"] = tmpSet
	}

	if v, ok := d.GetOk("metric_names"); ok {
		metricNamesSet := v.(*schema.Set).List()
		paramMap["MetricNames"] = helper.InterfacesStringsPoint(metricNamesSet)
	}

	if v, ok := d.GetOk("policy_name"); ok {
		paramMap["PolicyName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("content"); ok {
		paramMap["Content"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("receiver_uids"); ok {
		receiverUidsSet := v.(*schema.Set).List()
		paramMap["ReceiverUids"] = helper.InterfacesIntInt64Point(receiverUidsSet)
	}

	if v, ok := d.GetOk("receiver_groups"); ok {
		receiverGroupsSet := v.(*schema.Set).List()
		paramMap["ReceiverGroups"] = helper.InterfacesIntInt64Point(receiverGroupsSet)
	}

	if v, ok := d.GetOk("policy_ids"); ok {
		policyIdsSet := v.(*schema.Set).List()
		paramMap["PolicyIds"] = helper.InterfacesStringsPoint(policyIdsSet)
	}

	if v, ok := d.GetOk("alarm_levels"); ok {
		alarmLevelsSet := v.(*schema.Set).List()
		paramMap["AlarmLevels"] = helper.InterfacesStringsPoint(alarmLevelsSet)
	}

	service := MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var histories []*monitor.AlarmHistory

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMonitorAlarmHistoryByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		histories = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(histories))
	tmpList := make([]map[string]interface{}, 0, len(histories))

	if histories != nil {
		for _, alarmHistory := range histories {
			alarmHistoryMap := map[string]interface{}{}

			if alarmHistory.AlarmId != nil {
				alarmHistoryMap["alarm_id"] = alarmHistory.AlarmId
			}

			if alarmHistory.MonitorType != nil {
				alarmHistoryMap["monitor_type"] = alarmHistory.MonitorType
			}

			if alarmHistory.Namespace != nil {
				alarmHistoryMap["namespace"] = alarmHistory.Namespace
			}

			if alarmHistory.AlarmObject != nil {
				alarmHistoryMap["alarm_object"] = alarmHistory.AlarmObject
			}

			if alarmHistory.Content != nil {
				alarmHistoryMap["content"] = alarmHistory.Content
			}

			if alarmHistory.FirstOccurTime != nil {
				alarmHistoryMap["first_occur_time"] = alarmHistory.FirstOccurTime
			}

			if alarmHistory.LastOccurTime != nil {
				alarmHistoryMap["last_occur_time"] = alarmHistory.LastOccurTime
			}

			if alarmHistory.AlarmStatus != nil {
				alarmHistoryMap["alarm_status"] = alarmHistory.AlarmStatus
			}

			if alarmHistory.PolicyId != nil {
				alarmHistoryMap["policy_id"] = alarmHistory.PolicyId
			}

			if alarmHistory.PolicyName != nil {
				alarmHistoryMap["policy_name"] = alarmHistory.PolicyName
			}

			if alarmHistory.VPC != nil {
				alarmHistoryMap["vpc"] = alarmHistory.VPC
			}

			if alarmHistory.ProjectId != nil {
				alarmHistoryMap["project_id"] = alarmHistory.ProjectId
			}

			if alarmHistory.ProjectName != nil {
				alarmHistoryMap["project_name"] = alarmHistory.ProjectName
			}

			if alarmHistory.InstanceGroup != nil {
				instanceGroupList := []interface{}{}
				for _, instanceGroup := range alarmHistory.InstanceGroup {
					instanceGroupMap := map[string]interface{}{}

					if instanceGroup.Id != nil {
						instanceGroupMap["id"] = instanceGroup.Id
					}

					if instanceGroup.Name != nil {
						instanceGroupMap["name"] = instanceGroup.Name
					}

					instanceGroupList = append(instanceGroupList, instanceGroupMap)
				}

				alarmHistoryMap["instance_group"] = instanceGroupList
			}

			if alarmHistory.ReceiverUids != nil {
				alarmHistoryMap["receiver_uids"] = alarmHistory.ReceiverUids
			}

			if alarmHistory.ReceiverGroups != nil {
				alarmHistoryMap["receiver_groups"] = alarmHistory.ReceiverGroups
			}

			if alarmHistory.NoticeWays != nil {
				alarmHistoryMap["notice_ways"] = alarmHistory.NoticeWays
			}

			if alarmHistory.OriginId != nil {
				alarmHistoryMap["origin_id"] = alarmHistory.OriginId
			}

			if alarmHistory.AlarmType != nil {
				alarmHistoryMap["alarm_type"] = alarmHistory.AlarmType
			}

			if alarmHistory.EventId != nil {
				alarmHistoryMap["event_id"] = alarmHistory.EventId
			}

			if alarmHistory.Region != nil {
				alarmHistoryMap["region"] = alarmHistory.Region
			}

			if alarmHistory.PolicyExists != nil {
				alarmHistoryMap["policy_exists"] = alarmHistory.PolicyExists
			}

			if alarmHistory.MetricsInfo != nil {
				metricsInfoList := []interface{}{}
				for _, metricsInfo := range alarmHistory.MetricsInfo {
					metricsInfoMap := map[string]interface{}{}

					if metricsInfo.QceNamespace != nil {
						metricsInfoMap["qce_namespace"] = metricsInfo.QceNamespace
					}

					if metricsInfo.MetricName != nil {
						metricsInfoMap["metric_name"] = metricsInfo.MetricName
					}

					if metricsInfo.Period != nil {
						metricsInfoMap["period"] = metricsInfo.Period
					}

					if metricsInfo.Value != nil {
						metricsInfoMap["value"] = metricsInfo.Value
					}

					if metricsInfo.Description != nil {
						metricsInfoMap["description"] = metricsInfo.Description
					}

					metricsInfoList = append(metricsInfoList, metricsInfoMap)
				}

				alarmHistoryMap["metrics_info"] = metricsInfoList
			}

			if alarmHistory.Dimensions != nil {
				alarmHistoryMap["dimensions"] = alarmHistory.Dimensions
			}

			if alarmHistory.AlarmLevel != nil {
				alarmHistoryMap["alarm_level"] = alarmHistory.AlarmLevel
			}

			ids = append(ids, *alarmHistory.AlarmId)
			tmpList = append(tmpList, alarmHistoryMap)
		}

		_ = d.Set("histories", tmpList)
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
