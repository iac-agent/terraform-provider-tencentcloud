package monitor

import (
	"context"
	"crypto/md5"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMonitorAlarmNotices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentMonitorAlarmNoticesRead,
		Schema: map[string]*schema.Schema{
			"order": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ASC",
				Description: "排序方式 更新时间 ASC=forward 顺序 DESC=reverse 顺序",
			},
			"owner_uid": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "primary 账号 uid 是 用于create preset 通知。",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Alarm 通知 模板名称 用于fuzzy search。",
			},
			"receiver_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "To 过滤器 告警 通知 templates according 到 recipients，您 need 到 select 通知 用户 类型 USER=用户 GROUP=用户 组 Leave blank = 不 过滤器 通过 recipient。",
			},
			"user_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "列表 recipients。",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"group_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Receive 组 列表。",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"notice_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Receive 组 列表。",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于存储结果。",
			},

			"alarm_notice": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Alarm 通知 template 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm 通知 模板 ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm 通知 模板名称",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "最后修改时间。",
						},
						"updated_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last Modified By。",
						},
						"notice_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm 通知 类型 ALARM=Notification 不 restored OK=Notification restored ALL。",
						},
						"user_notices": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Alarm 通知 template 列表.(At most five)。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"receiver_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Recipient 类型 USER=用户 GROUP=用户 Group。",
									},
									"start_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 秒 since 通知 开始时间 00:00:00 (值 范围 0-86399)。",
									},
									"end_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 秒 since 通知 结束时间 00:00:00 (值 范围 0-86399)。",
									},
									"notice_way": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Notification Channel List EMAIL=Mail SMS=SMS CALL=Telephone WECHAT=WeChat RTX=Enterprise WeChat。",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"user_ids": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "用户 UID List。",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"group_ids": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "用户 组 ID 列表。",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"phone_order": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Telephone polling 列表。",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"phone_circle_times": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 telephone polls (取值范围：1-5)。",
									},
									"phone_inner_interval": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 秒 between calls 在 polling 会话 (取值范围：60-900)。",
									},
									"phone_circle_interval": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 秒 between polls (取值范围：60-900)。",
									},
									"need_phone_arrive_notice": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Contact 通知 必填 0= No 1= Yes。",
									},
									"phone_call_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Call 类型 SYNC= Simultaneous call CIRCLE= Round call 如果 此 参数 是 不 指定， 默认值为 round call。",
									},
									"weekday": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Notification 周期 1-7 表示Monday 到 Sunday。",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
								},
							},
						},
						"url_notices": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "最大callback notifications 是 3。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Callback URL (limited 到 256 字符)。",
									},
									"start_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Notification Start Time 数量 秒 在 start 的 day。",
									},
									"end_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Notification End Time Seconds 在 start 的 day。",
									},
									"weekday": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Notification 周期 1-7 表示Monday 到 Sunday。",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
								},
							},
						},
						"cls_notices": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A 最大 的 一个 告警 通知 可以 是 pushed 到 CLS 服务。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Regional。",
									},
									"log_set_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Log collection ID。",
									},
									"topic_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Theme ID。",
									},
									"enable": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Start-stop 状态，可以 不 是 transmitted，默认值 已启用 0= 已禁用，1= 已启用",
									},
								},
							},
						},
						"is_preset": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否为the 系统 默认值 通知 template 0=No 1=Yes。",
						},
						"notice_language": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Notification 语言 zh-CN=Chinese en-US=English。",
						},
						"policy_ids": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "列表 告警 策略 IDs bound 到 告警 通知 template。",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"amp_consumer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AMP 消费者 ID。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentMonitorAlarmNoticesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_alarm_notices.read")()

	var (
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		err            error
		alarmNotices   []interface{}
		alarmNotice    []*monitor.AlarmNotice
	)

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	alarmNoticeMap := make(map[string]interface{})
	alarmNoticeMap["order"] = helper.String(d.Get("order").(string))

	if v, ok := d.GetOk("owner_uid"); ok {
		alarmNoticeMap["ownerUid"] = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("name"); ok {
		alarmNoticeMap["name"] = helper.String(v.(string))
	}
	if v, ok := d.GetOk("receiver_type"); ok {
		alarmNoticeMap["receiverType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_ids"); ok {
		userIds := v.(*schema.Set).List()
		userIdsArr := make([]*int64, 0, len(userIds))
		for _, userId := range userIds {
			userIdsArr = append(userIdsArr, helper.Int64(userId.(int64)))
		}
		alarmNoticeMap["userIdArr"] = userIdsArr
	}

	if v, ok := d.GetOk("group_ids"); ok {
		groupIds := v.(*schema.Set).List()
		groupIdsArr := make([]*int64, 0, len(groupIds))
		for _, groupId := range groupIds {
			groupIdsArr = append(groupIdsArr, helper.Int64(groupId.(int64)))
		}
		alarmNoticeMap["groupArr"] = groupIdsArr
	}

	if v, ok := d.GetOk("notice_ids"); ok {
		noticeIds := v.(*schema.Set).List()
		noticeIdsArr := make([]*string, 0, len(noticeIds))
		for _, noticeId := range noticeIds {
			noticeIdsArr = append(noticeIdsArr, helper.String(noticeId.(string)))
		}
		alarmNoticeMap["noticeArr"] = noticeIdsArr
	}

	alarmNotice, err = monitorService.DescribeAlarmNoticeById(ctx, alarmNoticeMap)
	if err != nil {
		return err
	}
	for _, noticesItem := range alarmNotice {
		noticesItemMap := map[string]interface{}{
			"id":              noticesItem.Id,
			"name":            noticesItem.Name,
			"updated_at":      noticesItem.UpdatedAt,
			"updated_by":      noticesItem.UpdatedBy,
			"notice_type":     noticesItem.NoticeType,
			"is_preset":       noticesItem.IsPreset,
			"notice_language": noticesItem.NoticeLanguage,
			"policy_ids":      noticesItem.PolicyIds,
			"amp_consumer_id": noticesItem.AMPConsumerId,
		}

		userNoticesItems := make([]interface{}, 0, 100)
		for _, userNotices := range noticesItem.UserNotices {
			userNoticesItems = append(userNoticesItems, map[string]interface{}{
				"receiver_type":            userNotices.ReceiverType,
				"start_time":               userNotices.StartTime,
				"end_time":                 userNotices.EndTime,
				"notice_way":               userNotices.NoticeWay,
				"user_ids":                 userNotices.UserIds,
				"group_ids":                userNotices.GroupIds,
				"phone_order":              userNotices.PhoneOrder,
				"phone_circle_times":       userNotices.PhoneCircleTimes,
				"phone_inner_interval":     userNotices.PhoneInnerInterval,
				"phone_circle_interval":    userNotices.PhoneCircleInterval,
				"need_phone_arrive_notice": userNotices.NeedPhoneArriveNotice,
				"phone_call_type":          userNotices.PhoneCallType,
				"weekday":                  userNotices.Weekday,
			})
		}

		urlNoticesItems := make([]interface{}, 0, 100)
		for _, urlNotice := range noticesItem.URLNotices {
			urlNoticesItems = append(urlNoticesItems, map[string]interface{}{
				"url":        urlNotice.URL,
				"start_time": urlNotice.StartTime,
				"end_time":   urlNotice.EndTime,
				"weekday":    urlNotice.Weekday,
			})
		}

		clsNoticesItems := make([]interface{}, 0, 100)
		for _, clsNotice := range noticesItem.CLSNotices {
			clsNoticesItems = append(clsNoticesItems, map[string]interface{}{
				"region":     clsNotice.Region,
				"log_set_id": clsNotice.LogSetId,
				"topic_id":   clsNotice.TopicId,
				"enable":     clsNotice.Enable,
			})
		}
		noticesItemMap["user_notices"] = userNoticesItems
		noticesItemMap["url_notices"] = urlNoticesItems
		noticesItemMap["cls_notices"] = clsNoticesItems
		alarmNotices = append(alarmNotices, noticesItemMap)
	}

	md := md5.New()
	id := fmt.Sprintf("%x", md.Sum(nil))
	d.SetId(id)

	if err = d.Set("alarm_notice", alarmNotices); err != nil {
		return err
	}
	if output, ok := d.GetOk("result_output_file"); ok {
		return tccommon.WriteToFile(output.(string), alarmNotices)
	}
	return nil
}
