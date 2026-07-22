package monitor

import (
	"context"
	"fmt"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudMonitorBindingAlarmReceiver() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentMonitorBindingAlarmReceiverCreate,
		Read:   resourceTencentMonitorBindingAlarmReceiverRead,
		Update: resourceTencentMonitorBindingAlarmReceiverUpdate,
		Delete: resourceTencentMonitorBindingAlarmReceiverDelete,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Policy 组 ID 对于 binding receivers。",
			},
			"receivers": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "A 列表 receivers(将 overwrite 配置 的 服务器 或 other resources). Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 86399),
							Description:  "Alarm 周期 开始时间. 有效 值 ranges: (0~86399). 其中 removes date after 它 是 converted 到 Beijing 时间 作为 Unix 时间戳，对于 示例 7200 表示 '10:0:0'。",
						},
						"end_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      86400,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 86399),
							Description:  "End 的 告警 周期 Meaning 使用 `start_time`。",
						},
						"notify_way": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							MinItems:    1,
							Description: `Method of warning notification.Optional ` + helper.SliceFieldSerialize(monitorNotifyWays) + `.`,
						},
						"receiver_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(monitorReceiverTypes),
							Description:  "Receive 类型. Optional" + helper.SliceFieldSerialize(monitorReceiverTypes) + ".",
						},
						"receiver_group_list": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
							Description: "Alarm receive 组 ID 列表。",
						},
						"receiver_user_list": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
							Description: "Alarm receiver ID 列表。",
						},
						"receive_language": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(monitorReceiveLanguages),
							Default:      monitorReceiveLanguageCN,
							Description:  "Alert sending language. Optional" + helper.SliceFieldSerialize(monitorReceiveLanguages) + ".",
						},
					},
				},
			},
		},
	}
}

func resourceTencentMonitorBindingAlarmReceiverCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_binding_receiver.create")()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request        = monitor.NewModifyAlarmReceiversRequest()
		groupId        = int64(d.Get("group_id").(int))
	)

	info, err := monitorService.DescribePolicyGroup(ctx, groupId)
	if err != nil {
		return err
	}

	if info == nil {
		return fmt.Errorf("policy group %d not exist", groupId)
	}

	request.GroupId = &groupId
	request.Module = helper.String("monitor")
	request.ReceiverInfos = make([]*monitor.ReceiverInfo, 0, 10)

	for _, iface := range d.Get("receivers").([]interface{}) {
		var receiverInfo monitor.ReceiverInfo
		ifaceMap := iface.(map[string]interface{})

		receiverInfo.StartTime = helper.IntInt64(ifaceMap["start_time"].(int))
		receiverInfo.EndTime = helper.IntInt64(ifaceMap["end_time"].(int))
		receiverInfo.NotifyWay = helper.InterfacesStringsPoint(ifaceMap["notify_way"].([]interface{}))
		receiverInfo.ReceiverType = helper.String(ifaceMap["receiver_type"].(string))

		if ifaceMap["receiver_group_list"] != nil {
			receiverInfo.ReceiverGroupList = helper.InterfacesIntInt64Point(ifaceMap["receiver_group_list"].([]interface{}))
		}
		if ifaceMap["receiver_user_list"] != nil {
			receiverInfo.ReceiverUserList = helper.InterfacesIntInt64Point(ifaceMap["receiver_user_list"].([]interface{}))
		}

		if *receiverInfo.ReceiverType == monitorReceiverTypeGroup {
			if len(receiverInfo.ReceiverGroupList) < 1 {
				return fmt.Errorf("miss field receiver_group_list, this array at least  has one element when you choose `group` receiver_type")
			}
		}
		if *receiverInfo.ReceiverType == monitorReceiverTypeUser {
			if len(receiverInfo.ReceiverUserList) < 1 {
				return fmt.Errorf("miss field receiver_user_list, this array at least  has one element when you choose `user` receiver_type")
			}
		}

		receiverInfo.ReceiveLanguage = helper.String(ifaceMap["receive_language"].(string))
		request.ReceiverInfos = append(request.ReceiverInfos, &receiverInfo)
	}

	if err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if _, err = monitorService.client.UseMonitorClient().ModifyAlarmReceivers(request); err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%d", groupId))
	time.Sleep(3 * time.Second)

	return resourceTencentMonitorBindingAlarmReceiverRead(d, meta)
}

func resourceTencentMonitorBindingAlarmReceiverRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_binding_receiver.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		groupId        = int64(d.Get("group_id").(int))
	)

	info, err := monitorService.DescribePolicyGroup(ctx, groupId)
	if err != nil {
		return err
	}

	if info == nil {
		d.SetId("")
		return nil
	}

	list := make([]interface{}, 0, len(info.ReceiverInfos))

	for _, receiver := range info.ReceiverInfos {
		var receiverMap = make(map[string]interface{})
		receiverMap["start_time"] = receiver.StartTime
		receiverMap["end_time"] = receiver.EndTime
		receiverMap["notify_way"] = receiver.NotifyWay
		receiverMap["receiver_type"] = receiver.ReceiverType
		receiverMap["receiver_group_list"] = receiver.ReceiverGroupList
		receiverMap["receiver_user_list"] = receiver.ReceiverUserList
		receiverMap["receive_language"] = receiver.ReceiveLanguage
		list = append(list, receiverMap)

	}
	return d.Set("receivers", list)

}
func resourceTencentMonitorBindingAlarmReceiverUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_binding_receiver.update")()
	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request        = monitor.NewModifyAlarmReceiversRequest()
		groupId        = int64(d.Get("group_id").(int))
	)

	info, err := monitorService.DescribePolicyGroup(ctx, groupId)
	if err != nil {
		return err
	}

	if info == nil {
		d.SetId("")
		return nil
	}

	request.GroupId = &groupId
	request.Module = helper.String("monitor")
	request.ReceiverInfos = make([]*monitor.ReceiverInfo, 0, 10)

	for _, iface := range d.Get("receivers").([]interface{}) {
		var receiverInfo monitor.ReceiverInfo
		ifaceMap := iface.(map[string]interface{})

		receiverInfo.StartTime = helper.IntInt64(ifaceMap["start_time"].(int))
		receiverInfo.EndTime = helper.IntInt64(ifaceMap["end_time"].(int))
		receiverInfo.NotifyWay = helper.InterfacesStringsPoint(ifaceMap["notify_way"].([]interface{}))
		receiverInfo.ReceiverType = helper.String(ifaceMap["receiver_type"].(string))

		if ifaceMap["receiver_group_list"] != nil {
			receiverInfo.ReceiverGroupList = helper.InterfacesIntInt64Point(ifaceMap["receiver_group_list"].([]interface{}))
		}
		if ifaceMap["receiver_user_list"] != nil {
			receiverInfo.ReceiverUserList = helper.InterfacesIntInt64Point(ifaceMap["receiver_user_list"].([]interface{}))
		}
		if *receiverInfo.ReceiverType == monitorReceiverTypeGroup {
			if len(receiverInfo.ReceiverGroupList) < 1 {
				return fmt.Errorf("miss field receiver_group_list, this array at least  has one element when you choose `group` receiver_type")
			}
		}
		if *receiverInfo.ReceiverType == monitorReceiverTypeUser {
			if len(receiverInfo.ReceiverUserList) < 1 {
				return fmt.Errorf("miss field receiver_user_list, this array at least  has one element when you choose `user` receiver_type")
			}
		}
		receiverInfo.ReceiveLanguage = helper.String(ifaceMap["receive_language"].(string))
		request.ReceiverInfos = append(request.ReceiverInfos, &receiverInfo)
	}

	if err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if _, err = monitorService.client.UseMonitorClient().ModifyAlarmReceivers(request); err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	time.Sleep(3 * time.Second)
	return resourceTencentMonitorBindingAlarmReceiverRead(d, meta)
}

func resourceTencentMonitorBindingAlarmReceiverDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_binding_receiver.delete")()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request        = monitor.NewModifyAlarmReceiversRequest()
		groupId        = int64(d.Get("group_id").(int))
	)

	info, err := monitorService.DescribePolicyGroup(ctx, groupId)
	if err != nil {
		return err
	}

	if info == nil {
		d.SetId("")
		return nil
	}

	request.GroupId = &groupId
	request.Module = helper.String("monitor")
	if err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if _, err = monitorService.client.UseMonitorClient().ModifyAlarmReceivers(request); err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
