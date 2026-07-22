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
				Description: "值 fixed at monitor。",
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "排序方式 the first occurrence time in descending 排序依据 default. 有效值：ASC (ascending)，DESC (descending)。",
			},

			"start_time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "开始时间，which is the 时间戳 one day ago by default and the time when the alarm FirstOccurTime first occurs. An alarm record can be searched only if its FirstOccurTime is later than the StartTime。",
			},

			"end_time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "结束时间，which is the current 时间戳 and the time when the alarm FirstOccurTime first occurs. An alarm record can be searched only if its FirstOccurTime is earlier than the EndTime。",
			},

			"monitor_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Filter by 监控类型 有效值：MT_QCE (Tencent Cloud service monitoring)，MT_TAW (application performance monitoring)，MT_RUM (frontend performance monitoring)，MT_PROBE (cloud automated testing). 如果此参数为空，all types will be queried by default。",
			},

			"alarm_object": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter by alarm object. Fuzzy search with string is supported。",
			},

			"alarm_status": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Filter by 告警状态 有效值：ALARM (not resolved)，OK (resolved)，NO_CONF (expired)，NO_DATA (insufficient data). 如果此参数为空，all will be queried by default。",
			},

			"project_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Filter by project ID. 有效值：-1 (no project)，0 (default project)。",
			},

			"instance_group_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Filter by instance 组 ID",
			},

			"namespaces": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter by policy 类型 Monitoring 类型 and policy 类型 are first-级别 and second-级别 filters respectively and both need to be passed in. For example，[{MonitorType: MT_QCE，Namespace: cvm_device}]。",
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
				Description: "Filter by 指标名称",
			},

			"policy_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Fuzzy search by policy 名称",
			},

			"content": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Fuzzy search by alarm 内容",
			},

			"receiver_uids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Search by recipient。",
			},

			"receiver_groups": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Search by recipient group。",
			},

			"policy_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Search by 告警策略 ID list。",
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
				Description: "Alarm record list。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alarm_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm record ID。",
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
							Description: "Alarm object。",
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm 内容",
						},
						"first_occur_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "时间戳 of the first occurrence。",
						},
						"last_occur_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "时间戳 of the last occurrence。",
						},
						"alarm_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "告警状态 有效值：ALARM (not resolved)，OK (resolved)，NO_CONF (expired)，NO_DATA (insufficient data)。",
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
							Description: "VPC of alarm object for basic product alarm。",
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
							Description: "Instance group of alarm object。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Instance 组 ID",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance 组名称",
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
							Description: "Recipient list。",
						},
						"receiver_groups": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "Recipient group list。",
						},
						"notice_ways": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Alarm channel list. 有效值：SMS (SMS)，EMAIL (email)，CALL (phone)，WECHAT (WeChat)。",
						},
						"origin_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "告警策略 ID，which can be used when you call APIs (BindingPolicyObject，UnBindingAllPolicyObject，UnBindingPolicyObject) to bind/unbind instances or instance groups to/from an alarm policy。",
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
							Description: "是否policy exists. 有效值：0 (no)，1 (yes)。",
						},
						"metrics_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Metric informationNote: this field may return null，indicating that no valid values can be obtained。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"qce_namespace": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace 用于query data by Tencent Cloud service monitoring 类型",
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
										Description: "值 triggering alarm。",
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
							Description: "Dimension information of an instance that triggered alarms.Note: this field may return null，indicating that no valid values can be obtained。",
						},
						"alarm_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm 级别Note: this field may return null，indicating that no valid values can be obtained。",
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
