package monitor

import (
	"context"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMonitorAlarmBasicAlarms() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMonitorAlarmBasicAlarmsRead,
		Schema: map[string]*schema.Schema{
			"module": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Interface 模块 名称，当前值 监控。",
			},

			"start_time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "开始时间，默认为 一个 day 是 时间戳。",
			},

			"end_time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "结束时间，默认为 当前 时间戳。",
			},

			"occur_time_order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "排序方式 occurrence 时间，taking ASC 或 DESC 值。",
			},

			"project_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "过滤器 based 在 项目 ID。",
			},

			"view_names": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "过滤器 based 在 策略 类型",
			},

			"alarm_status": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "过滤器 based 在 告警状态",
			},

			"obj_like": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "过滤器 based 在 告警 objects。",
			},

			"instance_group_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "过滤器 based 在 实例 组 ID",
			},

			"metric_names": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "过滤器 通过 indicator 名称",
			},

			"alarms": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Alarm List。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID 此 告警。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "项目 ID",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Entry 名称",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "告警状态 ID，0 表示not recovered; 1 表示that 它 has been restored; 2,3,5 表示insufficient 数据; 4 表示it has expired。",
						},
						"alarm_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "告警状态，ALARM 表示not recovered; OK 表示that 它 has been restored; NO_ DATA 表示insufficient 数据; NO_ CONF 表示that 它 has expired。",
						},
						"group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Policy 组 ID",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy Group 名称",
						},
						"first_occur_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time 的 occurrence。",
						},
						"duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Duration （秒）。",
						},
						"last_occur_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "结束时间。",
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm 内容",
						},
						"obj_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm Object。",
						},
						"obj_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm 对象 ID。",
						},
						"view_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy 类型",
						},
						"vpc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC，仅 CVM has。",
						},
						"metric_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicator ID。",
						},
						"metric_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicator 名称",
						},
						"alarm_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Alarm 类型，0 表示 indicator 告警，2 表示 product 事件 告警，和 3 表示 平台 事件 告警。",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域",
						},
						"dimensions": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm 对象 dimension 信息。",
						},
						"notify_way": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Notification 方法。",
						},
						"instance_group": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "实例 Group Information。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_group_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "实例 组 ID",
									},
									"instance_group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例 Group 名称",
									},
								},
							},
						},
					},
				},
			},

			"warning": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "备注",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudMonitorAlarmBasicAlarmsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_alarm_basic_alarms.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("module"); ok {
		paramMap["Module"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("start_time"); ok {
		paramMap["StartTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("end_time"); ok {
		paramMap["EndTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("occur_time_order"); ok {
		paramMap["OccurTimeOrder"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_ids"); ok {
		projectIdsSet := v.(*schema.Set).List()
		paramMap["ProjectIds"] = helper.InterfacesIntInt64Point(projectIdsSet)
	}

	if v, ok := d.GetOk("view_names"); ok {
		viewNamesSet := v.(*schema.Set).List()
		paramMap["ViewNames"] = helper.InterfacesStringsPoint(viewNamesSet)
	}

	if v, ok := d.GetOk("alarm_status"); ok {
		alarmStatusSet := v.(*schema.Set).List()
		paramMap["AlarmStatus"] = helper.InterfacesIntInt64Point(alarmStatusSet)
	}

	if v, ok := d.GetOk("obj_like"); ok {
		paramMap["ObjLike"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_group_ids"); ok {
		instanceGroupIdsSet := v.(*schema.Set).List()
		paramMap["InstanceGroupIds"] = helper.InterfacesIntInt64Point(instanceGroupIdsSet)
	}

	if v, ok := d.GetOk("metric_names"); ok {
		metricNamesSet := v.(*schema.Set).List()
		paramMap["MetricNames"] = helper.InterfacesStringsPoint(metricNamesSet)
	}

	service := MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var alarms []*monitor.DescribeBasicAlarmListAlarms
	var warning *string
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, w, e := service.DescribeMonitorAlarmBasicAlarmsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		alarms = result
		warning = w
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(alarms))
	tmpList := make([]map[string]interface{}, 0, len(alarms))

	if alarms != nil {
		for _, describeBasicAlarmListAlarms := range alarms {
			describeBasicAlarmListAlarmsMap := map[string]interface{}{}

			if describeBasicAlarmListAlarms.Id != nil {
				describeBasicAlarmListAlarmsMap["id"] = describeBasicAlarmListAlarms.Id
			}

			if describeBasicAlarmListAlarms.ProjectId != nil {
				describeBasicAlarmListAlarmsMap["project_id"] = describeBasicAlarmListAlarms.ProjectId
			}

			if describeBasicAlarmListAlarms.ProjectName != nil {
				describeBasicAlarmListAlarmsMap["project_name"] = describeBasicAlarmListAlarms.ProjectName
			}

			if describeBasicAlarmListAlarms.Status != nil {
				describeBasicAlarmListAlarmsMap["status"] = describeBasicAlarmListAlarms.Status
			}

			if describeBasicAlarmListAlarms.AlarmStatus != nil {
				describeBasicAlarmListAlarmsMap["alarm_status"] = describeBasicAlarmListAlarms.AlarmStatus
			}

			if describeBasicAlarmListAlarms.GroupId != nil {
				describeBasicAlarmListAlarmsMap["group_id"] = describeBasicAlarmListAlarms.GroupId
			}

			if describeBasicAlarmListAlarms.GroupName != nil {
				describeBasicAlarmListAlarmsMap["group_name"] = describeBasicAlarmListAlarms.GroupName
			}

			if describeBasicAlarmListAlarms.FirstOccurTime != nil {
				describeBasicAlarmListAlarmsMap["first_occur_time"] = describeBasicAlarmListAlarms.FirstOccurTime
			}

			if describeBasicAlarmListAlarms.Duration != nil {
				describeBasicAlarmListAlarmsMap["duration"] = describeBasicAlarmListAlarms.Duration
			}

			if describeBasicAlarmListAlarms.LastOccurTime != nil {
				describeBasicAlarmListAlarmsMap["last_occur_time"] = describeBasicAlarmListAlarms.LastOccurTime
			}

			if describeBasicAlarmListAlarms.Content != nil {
				describeBasicAlarmListAlarmsMap["content"] = describeBasicAlarmListAlarms.Content
			}

			if describeBasicAlarmListAlarms.ObjName != nil {
				describeBasicAlarmListAlarmsMap["obj_name"] = describeBasicAlarmListAlarms.ObjName
			}

			if describeBasicAlarmListAlarms.ObjId != nil {
				describeBasicAlarmListAlarmsMap["obj_id"] = describeBasicAlarmListAlarms.ObjId
			}

			if describeBasicAlarmListAlarms.ViewName != nil {
				describeBasicAlarmListAlarmsMap["view_name"] = describeBasicAlarmListAlarms.ViewName
			}

			if describeBasicAlarmListAlarms.Vpc != nil {
				describeBasicAlarmListAlarmsMap["vpc"] = describeBasicAlarmListAlarms.Vpc
			}

			if describeBasicAlarmListAlarms.MetricId != nil {
				describeBasicAlarmListAlarmsMap["metric_id"] = describeBasicAlarmListAlarms.MetricId
			}

			if describeBasicAlarmListAlarms.MetricName != nil {
				describeBasicAlarmListAlarmsMap["metric_name"] = describeBasicAlarmListAlarms.MetricName
			}

			if describeBasicAlarmListAlarms.AlarmType != nil {
				describeBasicAlarmListAlarmsMap["alarm_type"] = describeBasicAlarmListAlarms.AlarmType
			}

			if describeBasicAlarmListAlarms.Region != nil {
				describeBasicAlarmListAlarmsMap["region"] = describeBasicAlarmListAlarms.Region
			}

			if describeBasicAlarmListAlarms.Dimensions != nil {
				describeBasicAlarmListAlarmsMap["dimensions"] = describeBasicAlarmListAlarms.Dimensions
			}

			if describeBasicAlarmListAlarms.NotifyWay != nil {
				describeBasicAlarmListAlarmsMap["notify_way"] = describeBasicAlarmListAlarms.NotifyWay
			}

			if describeBasicAlarmListAlarms.InstanceGroup != nil {
				instanceGroupList := []interface{}{}
				for _, instanceGroup := range describeBasicAlarmListAlarms.InstanceGroup {
					instanceGroupMap := map[string]interface{}{}

					if instanceGroup.InstanceGroupId != nil {
						instanceGroupMap["instance_group_id"] = instanceGroup.InstanceGroupId
					}

					if instanceGroup.InstanceGroupName != nil {
						instanceGroupMap["instance_group_name"] = instanceGroup.InstanceGroupName
					}

					instanceGroupList = append(instanceGroupList, instanceGroupMap)
				}

				describeBasicAlarmListAlarmsMap["instance_group"] = instanceGroupList
			}

			ids = append(ids, strconv.Itoa(int(*describeBasicAlarmListAlarms.Id)))
			tmpList = append(tmpList, describeBasicAlarmListAlarmsMap)
		}

		_ = d.Set("alarms", tmpList)
	}

	if warning != nil {
		_ = d.Set("warning", warning)
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
