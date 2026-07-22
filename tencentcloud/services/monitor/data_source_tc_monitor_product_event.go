package monitor

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func DataSourceTencentCloudMonitorProductEvent() *schema.Resource {

	defaultStartTime := time.Now().Unix() - 3600

	return &schema.Resource{
		Read: dataSourceTencentMonitorProductEventRead,
		Schema: map[string]*schema.Schema{
			"product_name": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "产品类型 filtering，such as `cvm` for cloud server。",
			},
			"event_name": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "事件名称 filtering，such as `guest_reboot` 表示that the machine restart。",
			},
			"instance_id": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Affect objects，such as `ins-19708ino`。",
			},
			"dimensions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Instance dimension 名称，eg: `deviceWanIp` for internet ip。",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Instance dimension 值，eg: `119.119.119.119` for internet ip。",
						},
					},
				},
				Description: "Dimensional composition of instance objects。",
			},
			"region_list": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "地域 filter，such as `gz`。",
			},
			"type": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Event type filtering, with value range " + helper.SliceFieldSerialize(monitorEventTypes) + ", indicating state change and abnormal events.",
			},
			"status": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Event status filter, value range " + helper.SliceFieldSerialize(monitorEventStatus) + ", indicating recovered, unrecovered and stateless.",
			},
			"project_id": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "项目 ID filter。",
			},
			"is_alarm_config": {
				Type:        schema.TypeInt,
				Default:     0,
				Optional:    true,
				Description: "告警状态 configuration filter，1means configured，0(default) means not configured。",
			},
			"start_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     defaultStartTime,
				Description: "Start 时间戳 for this query，eg:`1588230000`. Default 开始时间 is `now-3600`。",
			},
			"end_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     defaultStartTime + 600,
				Description: "End 时间戳 for this query，eg:`1588232111`. Default 开始时间 is `now-3000`。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于存储结果。",
			},
			// Computed values
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list events. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "事件 ID",
						},
						"event_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event short 名称",
						},
						"event_ename": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event english 名称",
						},
						"event_cname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event chinese 名称",
						},
						"product_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Product short 名称",
						},
						"product_ename": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Product english 名称",
						},
						"product_cname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Product chinese 名称",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance ID this event。",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 this instance。",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The 地域 of this instance。",
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "项目 ID this instance。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态 this event。",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 this event。",
						},
						"start_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The start 时间戳 of this event。",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The update 时间戳 of this event。",
						},
						"support_alarm": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否support alarm。",
						},
						"is_alarm_config": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否configure alarm。",
						},
						"dimensions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A 列表 event dimensions. Each element 包含following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The 键 of this dimension。",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "名称 this dimension。",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The 值 of this dimension。",
									},
								},
							},
						},
						"addition_msg": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A 列表 addition 消息 Each element 包含following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The 键 of this addition 消息",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "名称 this addition 消息",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The 值 of this addition 消息",
									},
								},
							},
						},
						"group_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A 列表 group info. Each element 包含following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Policy 组 ID",
									},
									"group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Policy 组名称",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentMonitorProductEventRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_product_event.read")()

	var (
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request        = monitor.NewDescribeProductEventListRequest()
		response       *monitor.DescribeProductEventListResponse
		events         []*monitor.DescribeProductEventListEvents
		offset         int64 = 0
		limit          int64 = 20
		err            error
	)

	request.IsAlarmConfig = helper.IntInt64(d.Get("is_alarm_config").(int))
	request.StartTime = helper.IntInt64(d.Get("start_time").(int))
	request.EndTime = helper.IntInt64(d.Get("end_time").(int))

	if *request.EndTime < *request.StartTime {
		return fmt.Errorf("field `start_time` should >= `end_time`")
	}

	for k, v := range map[string]*[]*string{
		"product_name": &request.ProductName,
		"event_name":   &request.EventName,
		"instance_id":  &request.InstanceId,
		"region_list":  &request.RegionList,
		"type":         &request.Type,
		"status":       &request.Status,
		"project_id":   &request.Project,
	} {
		if iface, ok := d.GetOk(k); ok {
			*v = helper.Strings(helper.InterfacesStrings(iface.([]interface{})))
		}
	}

	if iface, ok := d.GetOk("dimensions"); ok {

		request.Dimensions = make([]*monitor.DescribeProductEventListDimensions, 0, 100)
		for _, mapIface := range iface.([]interface{}) {

			m := mapIface.(map[string]interface{})

			var mtrDimension monitor.DescribeProductEventListDimensions

			if m["name"] == nil || m["value"] == nil {
				return fmt.Errorf("miss `name` or `value` from field `dimensions`")
			}
			mtrDimension.Name = helper.String(m["name"].(string))
			mtrDimension.Value = helper.String(m["value"].(string))
			request.Dimensions = append(request.Dimensions, &mtrDimension)
		}
	}

	if iface, ok := d.GetOk("region_list"); ok {
		request.RegionList = helper.InterfacesStringsPoint(iface.([]interface{}))
	}

	if iface, ok := d.GetOk("type"); ok {
		request.Type = make([]*string, 0, 3)
		for _, strIface := range iface.([]interface{}) {
			str := strIface.(string)
			if !helper.StringsContain(monitorEventTypes, str) {
				return fmt.Errorf("type `%s` no support now", str)
			} else {
				request.Type = append(request.Type, &str)
			}
		}
	}

	if iface, ok := d.GetOk("status"); ok {
		request.Status = make([]*string, 0, 3)
		for _, strIface := range iface.([]interface{}) {
			str := strIface.(string)
			if !helper.StringsContain(monitorEventStatus, str) {
				return fmt.Errorf("status `%s` no support now", str)
			} else {
				request.Status = append(request.Status, &str)
			}
		}
	}

	if iface, ok := d.GetOk("project_id"); ok {
		request.Project = make([]*string, 0, 3)
		for _, strIface := range iface.([]interface{}) {
			str := strIface.(string)
			if _, err := strconv.ParseInt(str, 10, 64); err != nil {
				return fmt.Errorf("each element of `project_id` must be an integer ,`%s` not fit the rule ", str)
			} else {
				request.Project = append(request.Project, &str)
			}
		}
	}

	request.IsAlarmConfig = helper.IntInt64(d.Get("is_alarm_config").(int))

	startTime := d.Get("start_time").(int)
	endTime := d.Get("end_time").(int)

	if endTime < startTime {
		return fmt.Errorf("`start_time` should <= `end_time`, config not fit the rule")
	}

	request.StartTime = helper.IntInt64(startTime)
	request.EndTime = helper.IntInt64(endTime)
	request.Offset = &offset
	request.Limit = &limit
	request.Module = helper.String("monitor")

	var finish = false
	for {

		if finish {
			break
		}

		if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			if response, err = monitorService.client.UseMonitorClient().DescribeProductEventList(request); err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			events = append(events, response.Response.Events...)
			if len(response.Response.Events) < int(limit) {
				finish = true
			}
			return nil
		}); err != nil {
			return err
		}

		offset = offset + limit
	}

	var list = make([]interface{}, 0, len(events))

	for _, event := range events {

		var listItem = map[string]interface{}{}

		var dimensions []interface{}
		var additionMsg []interface{}
		var groupInfo []interface{}

		for _, item := range event.Dimensions {
			dimensions = append(dimensions,
				map[string]interface{}{
					"key":   item.Key,
					"name":  item.Name,
					"value": item.Value,
				})
		}

		for _, item := range event.AdditionMsg {
			additionMsg = append(additionMsg,
				map[string]interface{}{
					"key":   item.Key,
					"name":  item.Name,
					"value": item.Value,
				})
		}
		for _, item := range event.GroupInfo {
			groupInfo = append(groupInfo,
				map[string]interface{}{
					"group_id":   item.GroupId,
					"group_name": item.GroupName,
				})
		}

		listItem["dimensions"] = dimensions
		listItem["addition_msg"] = additionMsg
		listItem["group_info"] = groupInfo
		listItem["event_id"] = event.EventId
		listItem["event_name"] = event.EventName
		listItem["event_ename"] = event.EventEName
		listItem["event_cname"] = event.EventCName
		listItem["product_name"] = event.ProductName
		listItem["product_ename"] = event.ProductEName
		listItem["product_cname"] = event.ProductCName
		listItem["instance_id"] = event.InstanceId
		listItem["instance_name"] = event.InstanceName
		listItem["region"] = event.Region
		listItem["project_id"] = event.ProjectId
		listItem["status"] = event.Status
		listItem["type"] = event.Type
		listItem["start_time"] = event.StartTime
		listItem["update_time"] = event.UpdateTime
		listItem["support_alarm"] = event.SupportAlarm
		listItem["is_alarm_config"] = event.IsAlarmConfig

		list = append(list, listItem)
	}

	if err = d.Set("list", list); err != nil {
		return err
	}

	md := md5.New()
	_, _ = md.Write([]byte(request.ToJsonString()))
	id := fmt.Sprintf("%x", md.Sum(nil))
	d.SetId(id)

	if output, ok := d.GetOk("result_output_file"); ok {
		return tccommon.WriteToFile(output.(string), list)
	}
	return nil
}
