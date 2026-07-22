package monitor

import (
	"crypto/md5"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func DataSourceTencentCloudMonitorPolicyGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentMonitorPolicyGroupsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Policy 组名称 for query。",
			},
			"policy_view_names": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "The policy view for query。",
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
				Description: "A list policy groups. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The policy group id。",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The policy 组名称",
						},
						"is_open": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether open or not。",
						},
						"policy_view_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The policy group view 名称",
						},
						"last_edit_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Recently edited 用户 uin。",
						},
						"use_sum": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 instances of policy group bindings。",
						},
						"no_shielded_sum": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 unmasked instances of policy group bindings。",
						},
						"is_default": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "If is default policy group or not，`0` represents the non-default policy，and `1` represents the default policy。",
						},
						"conditions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A 列表 threshold rules. Each element 包含following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_show_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "名称 this metric。",
									},
									"period": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Data aggregation cycle (unit second)。",
									},
									"metric_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ID this metric。",
									},
									"rule_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Threshold rule ID。",
									},
									"metric_unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unit of this metric。",
									},
									"alarm_notify_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Alarm sending convergence 类型 `0` continuous alarm，`1` 索引 alarm。",
									},
									"alarm_notify_period": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Alarm sending cycle per second. `<0` does not fire，`0` only fires once，and `>0` fires every triggerTime second。",
									},
									"calc_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Compare 类型，`1` means more than，`2`  means greater than or equal，`3` means less than，`4` means less than or equal to，`5` means equal，`6` means not equal，`7` means days rose，`8` means days fell，`9` means weeks rose，`10` means weeks fell，`11` means 周期 rise，`12` means 周期 fell。",
									},
									"calc_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Threshold 值",
									},
									"continue_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "How long does the triggering rule last (per second)。",
									},
								},
							},
						},
						"event_conditions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A 列表 event rules. Each element 包含following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"event_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ID this event metric。",
									},
									"event_show_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "名称 this event metric。",
									},
									"rule_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Threshold rule ID。",
									},
									"alarm_notify_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Alarm sending convergence 类型 `0` continuous alarm，`1` 索引 alarm。",
									},
									"alarm_notify_period": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Alarm sending cycle per second. `<0` does not fire，`0` only fires once，and `>0` fires every triggerTime second。",
									},
								},
							},
						},
						"receivers": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A 列表 receivers. Each element 包含following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"receiver_group_list": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Computed:    true,
										Description: "Alarm receive 组 ID list。",
									},
									"receiver_user_list": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Computed:    true,
										Description: "Alarm receiver ID list。",
									},
									"uid_list": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Computed:    true,
										Description: "The phone alerts the receiver uid。",
									},
									"start_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Alarm 周期 开始时间.Range [0,86399]，which removes the date after it is converted to Beijing time as a Unix 时间戳，for example 7200 means '10:0:0'。",
									},
									"end_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "End of alarm 周期 Meaning with `start_time`。",
									},
									"notify_way": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Computed:    true,
										Description: `Method of warning notification.Optional ` + helper.SliceFieldSerialize(monitorNotifyWays) + `.`,
									},
									"receiver_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Receive 类型 可选 'group' or '用户'。",
									},
									"round_number": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Telephone alarm number。",
									},
									"round_interval": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Telephone alarm interval per round (seconds)。",
									},
									"person_interval": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Telephone 警告 to individual interval (seconds)。",
									},
									"need_send_notice": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Do need a telephone alarm contact prompt.You don't need 0，you need 1。",
									},
									"send_for": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Computed:    true,
										Description: `Telephone warning time.Option "OCCUR", "RECOVER".`,
									},
									"recover_notify": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Computed:    true,
										Description: `Restore notification mode. Optional "SMS".`,
									},
									"receive_language": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Alert sending 语言",
									},
								},
							},
						},
						"can_set_default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether it can be set as the default policy。",
						},
						"parent_group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Parent policy 组 ID",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy group 备注",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The project ID to which the policy group belongs。",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The policy group update 时间戳。",
						},
						"insert_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The policy group create 时间戳。",
						},
					},
				},
			},
		},
	}
}
func dataSourceTencentMonitorPolicyGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_policy_groups.read")()

	var (
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request        = monitor.NewDescribePolicyGroupListRequest()
		response       *monitor.DescribePolicyGroupListResponse
		err            error

		name            = d.Get("name").(string)
		policyViewNames = helper.InterfacesStrings(d.Get("policy_view_names").([]interface{}))

		list            = make([]interface{}, 0, 100)
		offset    int64 = 0
		limit     int64 = 20
		groupList       = make([]*monitor.DescribePolicyGroupListGroup, 0, 10)
		finish    bool
	)

	request.Module = helper.String("monitor")
	request.Offset = &offset
	request.Limit = &limit

	if name != "" {
		request.Like = &name
	}

	if len(policyViewNames) != 0 {
		request.ViewNames = helper.Strings(policyViewNames)
	}

	for {
		if finish {
			break
		}
		if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			if response, err = monitorService.client.UseMonitorClient().DescribePolicyGroupList(request); err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			groupList = append(groupList, response.Response.GroupList...)
			if len(response.Response.GroupList) < int(limit) {
				finish = true
			}
			return nil
		}); err != nil {
			return err
		}
		offset = offset + limit
	}

	for _, group := range groupList {
		var listItem = map[string]interface{}{}
		listItem["group_id"] = group.GroupId
		listItem["group_name"] = group.GroupName
		listItem["is_open"] = group.IsOpen
		listItem["policy_view_name"] = group.ViewName
		listItem["last_edit_uin"] = group.LastEditUin
		listItem["use_sum"] = group.UseSum
		listItem["no_shielded_sum"] = group.NoShieldedSum
		listItem["is_default"] = group.IsDefault
		listItem["can_set_default"] = group.CanSetDefault
		listItem["parent_group_id"] = group.ParentGroupId
		listItem["remark"] = group.Remark
		listItem["project_id"] = group.ProjectId
		listItem["update_time"] = group.UpdateTime
		listItem["insert_time"] = group.InsertTime

		conditions := make([]interface{}, 0, 100)
		for _, item := range group.Conditions {
			conditions = append(conditions, map[string]interface{}{
				"metric_show_name":    item.MetricShowName,
				"period":              item.Period,
				"metric_id":           item.MetricId,
				"rule_id":             item.RuleId,
				"metric_unit":         item.Unit,
				"alarm_notify_type":   item.AlarmNotifyType,
				"alarm_notify_period": item.AlarmNotifyPeriod,
				"calc_type":           item.CalcType,
				"calc_value":          item.CalcValue,
				"continue_time":       item.ContinueTime,
			})
		}
		listItem["conditions"] = conditions

		eventConditions := make([]interface{}, 0, 100)
		for _, item := range group.EventConditions {
			eventConditions = append(eventConditions, map[string]interface{}{
				"event_id":            item.EventId,
				"event_show_name":     item.EventShowName,
				"rule_id":             item.RuleId,
				"alarm_notify_type":   item.AlarmNotifyType,
				"alarm_notify_period": item.AlarmNotifyPeriod,
			})
		}
		listItem["event_conditions"] = eventConditions

		receivers := make([]interface{}, 0, 100)
		for _, item := range group.ReceiverInfos {

			receiver := map[string]interface{}{
				"start_time":       item.StartTime,
				"end_time":         item.EndTime,
				"receiver_type":    item.ReceiverType,
				"round_number":     item.RoundNumber,
				"round_interval":   item.RoundInterval,
				"person_interval":  item.PersonInterval,
				"need_send_notice": item.NeedSendNotice,
				"receive_language": item.ReceiveLanguage,
			}
			{
				slice := make([]int64, 0, len(item.ReceiverGroupList))
				for _, value := range item.ReceiverGroupList {
					slice = append(slice, *value)
				}
				receiver["receiver_group_list"] = slice
			}

			{
				slice := make([]int64, 0, len(item.ReceiverUserList))
				for _, value := range item.ReceiverUserList {
					slice = append(slice, *value)
				}
				receiver["receiver_user_list"] = slice
			}

			{
				slice := make([]int64, 0, len(item.UidList))
				for _, value := range item.UidList {
					slice = append(slice, *value)
				}
				receiver["uid_list"] = slice
			}

			{
				slice := make([]string, 0, len(item.NotifyWay))
				for _, value := range item.NotifyWay {
					slice = append(slice, *value)
				}
				receiver["notify_way"] = slice
			}

			{
				slice := make([]string, 0, len(item.SendFor))
				for _, value := range item.SendFor {
					slice = append(slice, *value)
				}
				receiver["send_for"] = slice
			}

			{
				slice := make([]string, 0, len(item.RecoverNotify))
				for _, value := range item.RecoverNotify {
					slice = append(slice, *value)
				}
				receiver["recover_notify"] = slice
			}
			receivers = append(receivers, receiver)

		}
		listItem["receivers"] = receivers

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
