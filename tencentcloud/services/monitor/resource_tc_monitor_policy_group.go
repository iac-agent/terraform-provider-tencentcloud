package monitor

import (
	"context"
	"fmt"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudMonitorPolicyGroup() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource has been deprecated in Terraform TencentCloud provider version 1.59.18. Please use 'tencentcloud_monitor_alarm_policy' instead.",
		Create:             resourceTencentMonitorPolicyGroupCreate,
		Read:               resourceTencentMonitorPolicyGroupRead,
		Update:             resourceTencentMonitorPolicyGroupUpdate,
		Delete:             resourceTencentMonitorPolicyGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Policy 组名称，长度 should between 1 和 20。",
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 20),
			},
			"policy_view_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Policy view 名称，eg:`cvm_device`,`BANDWIDTHPACKAGE`，refer 到 `数据.tencentcloud_monitor_policy_conditions(policy_view_name)`。",
			},
			"remark": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(0, 100),
				Description:  "Policy 组's 备注 信息。",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     0,
				Description: "项目 ID 到 其中 策略 组 belongs，默认为 `0`。",
			},
			"is_union_rule": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description:  "和 或 relation 的 indicator 告警 规则. 有效值：`0`，`1`. `0` 表示 或 规则 (如果 any 规则 是 met， 告警 将 是 raised)，`1` 表示 和 规则 (如果 all 规则 是 met， 告警 将 是 raised). 默认为 0。",
			},
			"conditions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A 列表 阈值 规则. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "ID 的 metric，refer 到 `数据.tencentcloud_monitor_policy_conditions(metric_id)`。",
						},
						"alarm_notify_type": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
							Description:  "Alarm sending convergence 类型 `0` continuous 告警，`1` 索引 告警。",
						},
						"alarm_notify_period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Alarm sending cycle per second. <0 does 不 fire，`0` 仅 fires once，和 >0 fires every triggerTime second。",
						},
						"calc_type": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(1, 12),
							Description:  "Compare 类型 有效 值 ranges: [1~12]. `1` 表示 more 比，`2` 表示 greater 比 或 equal，`3` 表示 less 比，`4` 表示 less 比 或 equal 到，`5` 表示 equal，`6` 表示 不 equal，`7` 表示 days rose，`8` 表示 days fell，`9` 表示 weeks rose，`10` 表示 weeks fell，`11` 表示 周期 rise，`12` 表示 周期 fell，refer 到 `数据.tencentcloud_monitor_policy_conditions(calc_type_keys)`。",
						},
						"calc_value": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Computed:    true,
							Description: "Threshold 值，refer 到 `数据.tencentcloud_monitor_policy_conditions(calc_value_*)`。",
						},
						"calc_period": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Data aggregation cycle (单位 的 second)，如果 metric has 默认值 可以 不 是 filled，refer 到 `数据.tencentcloud_monitor_policy_conditions(period_keys)`。",
						},
						"continue_period": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "规则 triggers alert 该 lasts 对于 several detection cycles，refer 到 `数据.tencentcloud_monitor_policy_conditions(period_num_keys)`。",
						},
					},
				},
			},
			"event_conditions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A 列表 事件 规则. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "ID 此 事件 metric，refer 到 `数据.tencentcloud_monitor_policy_conditions(event_id)。",
						},
						"alarm_notify_type": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
							Description:  "Alarm sending convergence 类型 `0` continuous 告警，`1` 索引 告警。",
						},
						"alarm_notify_period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Alarm sending cycle per second. <0 does 不 fire，`0` 仅 fires once，和 >0 fires every triggerTime second。",
						},
					},
				},
			},
			// computed value
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
							Description: "Alarm receive 组 ID 列表。",
						},
						"receiver_user_list": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Computed:    true,
							Description: "Alarm receiver ID 列表。",
						},
						"uid_list": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Computed:    true,
							Description: "phone alerts receiver uid。",
						},
						"start_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Alarm 周期 开始时间.Range [0,86400]，其中 removes date after 它 是 converted 到 Beijing 时间 作为 Unix 时间戳，对于 示例 7200 表示 '10:0:0'。",
						},
						"end_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "End 的 告警 周期 Meaning 使用 `start_time`。",
						},
						"notify_way": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: `Method of warning notification. Valid values: "SMS", "SITE", "EMAIL", "CALL", "WECHAT".`,
						},
						"receiver_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Receive 类型 有效值：组，用户 '组' (receiving 组) 或 '用户' (receiver)。",
						},
						"round_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Telephone 告警 数量。",
						},
						"round_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Telephone 告警 间隔 per round (秒)。",
						},
						"person_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Telephone 警告 到 individual 间隔 (秒)。",
						},
						"need_send_notice": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Do need telephone 告警 contact prompt. You don't need `0`，您 need `1`。",
						},
						"send_for": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: `Telephone warning time. Valid values: "OCCUR","RECOVER".`,
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
			"binding_objects": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 binding objects(列表 仅 those 在 `provider.地域`). Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unique_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Object 唯一 ID。",
						},
						"dimensions_json": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Represents collection 的 dimensions 的 对象 实例，json 格式",
						},
						"is_shielded": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否object 是 shielded 或 不，0 表示 unshielded 和 1 表示 shielded。",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 其中 对象 是 located。",
						},
					},
				},
			},
			"dimension_group": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "A 列表 dimensions 对于 此 策略 组。",
			},
			"support_regions": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "Support regions 此 策略 组。",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "策略 组 更新时间。",
			},
			"last_edit_uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Recently edited 用户 uin。",
			},
		},
	}
}
func resourceTencentMonitorPolicyGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_policy_group.create")()

	var (
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request        = monitor.NewCreatePolicyGroupRequest()
	)

	request.GroupName = helper.String(d.Get("group_name").(string))
	request.ViewName = helper.String(d.Get("policy_view_name").(string))
	if iface, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(iface.(string))
	}
	request.IsUnionRule = helper.IntInt64(d.Get("is_union_rule").(int))
	request.ProjectId = helper.IntInt64(d.Get("project_id").(int))
	request.Module = helper.String("monitor")

	if iface, ok := d.GetOk("conditions"); ok {
		request.Conditions = make([]*monitor.CreatePolicyGroupCondition, 0, 10)
		for _, item := range iface.([]interface{}) {
			m := item.(map[string]interface{})
			createPolicyGroupCondition := monitor.CreatePolicyGroupCondition{}
			createPolicyGroupCondition.MetricId = helper.IntInt64(m["metric_id"].(int))
			createPolicyGroupCondition.AlarmNotifyType = helper.IntInt64(m["alarm_notify_type"].(int))
			createPolicyGroupCondition.AlarmNotifyPeriod = helper.IntInt64(m["alarm_notify_period"].(int))
			if m["calc_type"] != nil {
				createPolicyGroupCondition.CalcType = helper.IntInt64(m["calc_type"].(int))
			}
			if m["calc_value"] != nil {
				createPolicyGroupCondition.CalcValue = helper.Float64(m["calc_value"].(float64))
			}
			if m["calc_period"] != nil {
				createPolicyGroupCondition.CalcPeriod = helper.IntInt64(m["calc_period"].(int))
			}
			if m["continue_period"] != nil {
				createPolicyGroupCondition.ContinuePeriod = helper.IntInt64(m["continue_period"].(int))
			}
			request.Conditions = append(request.Conditions, &createPolicyGroupCondition)
		}
	}

	if iface, ok := d.GetOk("event_conditions"); ok {
		request.EventConditions = make([]*monitor.CreatePolicyGroupEventCondition, 0, 10)
		for _, item := range iface.([]interface{}) {
			m := item.(map[string]interface{})
			createPolicyGroupCondition := monitor.CreatePolicyGroupEventCondition{}
			createPolicyGroupCondition.EventId = helper.IntInt64(m["event_id"].(int))
			createPolicyGroupCondition.AlarmNotifyType = helper.IntInt64(m["alarm_notify_type"].(int))
			createPolicyGroupCondition.AlarmNotifyPeriod = helper.IntInt64(m["alarm_notify_period"].(int))
			request.EventConditions = append(request.EventConditions, &createPolicyGroupCondition)
		}
	}

	var groupId *int64
	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := monitorService.client.UseMonitorClient().CreatePolicyGroup(request)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		groupId = response.Response.GroupId

		return nil
	}); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%d", *groupId))
	return resourceTencentMonitorPolicyGroupRead(d, meta)
}

func resourceTencentMonitorPolicyGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_policy_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request        = monitor.NewDescribePolicyGroupInfoRequest()
		response       *monitor.DescribePolicyGroupInfoResponse
	)

	groupId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("id [%s] is broken", d.Id())
	}

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

	if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if response, err = monitorService.client.UseMonitorClient().DescribePolicyGroupInfo(request); err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	if response == nil {
		d.SetId("")
		return nil
	}
	var errs []error
	errs = append(errs,
		d.Set("group_name", response.Response.GroupName),
		d.Set("policy_view_name", response.Response.ViewName),
		d.Set("remark", response.Response.Remark),
		d.Set("is_union_rule", response.Response.IsUnionRule),
		d.Set("project_id", response.Response.ProjectId),
	)

	var conditions = make([]interface{}, 0, 100)
	for _, condition := range response.Response.ConditionsConfig {
		m := map[string]interface{}{}
		m["metric_id"] = condition.MetricId
		m["alarm_notify_type"] = condition.AlarmNotifyType
		m["alarm_notify_period"] = condition.AlarmNotifyPeriod
		m["calc_type"] = condition.CalcType
		m["calc_value"] = condition.CalcValue
		m["calc_period"] = condition.Period
		m["continue_period"] = condition.ContinueTime
		conditions = append(conditions, m)
	}
	errs = append(errs, d.Set("conditions", conditions))

	var eventConditions = make([]interface{}, 0, 100)
	for _, condition := range response.Response.EventConfig {
		m := map[string]interface{}{}
		m["event_id"] = condition.EventId
		m["alarm_notify_type"] = condition.AlarmNotifyType
		m["alarm_notify_period"] = condition.AlarmNotifyPeriod
		eventConditions = append(eventConditions, m)
	}
	errs = append(errs, d.Set("event_conditions", eventConditions))

	receivers := make([]interface{}, 0, 100)
	for _, item := range response.Response.ReceiverInfos {

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
	errs = append(errs, d.Set("receivers", receivers))

	errs = append(errs,
		d.Set("group_name", response.Response.GroupName),
		d.Set("policy_view_name", response.Response.ViewName),
		d.Set("remark", response.Response.Remark),
		d.Set("is_union_rule", response.Response.IsUnionRule),
		d.Set("support_regions", response.Response.Region),
		d.Set("dimension_group", response.Response.DimensionGroup),
		d.Set("update_time", response.Response.UpdateTime),
		d.Set("last_edit_uin", response.Response.LastEditUin),
	)

	objects, err := monitorService.DescribeBindingPolicyObjectList(ctx, groupId)
	if err != nil {
		return err
	}
	bindingObjects := make([]interface{}, 0, len(objects))

	for _, event := range objects {
		var listItem = map[string]interface{}{}
		listItem["region"] = event.Region
		listItem["unique_id"] = event.UniqueId
		listItem["dimensions_json"] = event.Dimensions
		listItem["is_shielded"] = event.IsShielded
		listItem["region"] = event.Region
		bindingObjects = append(bindingObjects, listItem)
	}
	errs = append(errs, d.Set("binding_objects", bindingObjects))

	if len(errs) > 0 {
		return errs[0]
	} else {
		return nil
	}
}

func resourceTencentMonitorPolicyGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_policy_group.update")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request        = monitor.NewModifyPolicyGroupRequest()
	)
	groupId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("id [%s] is broken", d.Id())
	}

	info, err := monitorService.DescribePolicyGroup(ctx, groupId)
	if err != nil {
		return err
	}
	if info == nil {
		d.SetId("")
		return nil
	}
	request.GroupId = &groupId
	request.GroupName = helper.String(d.Get("group_name").(string))
	request.ViewName = helper.String(d.Get("policy_view_name").(string))
	request.IsUnionRule = helper.IntInt64(d.Get("is_union_rule").(int))
	request.Module = helper.String("monitor")

	if iface, ok := d.GetOk("conditions"); ok {
		request.Conditions = make([]*monitor.ModifyPolicyGroupCondition, 0, 10)
		for _, item := range iface.([]interface{}) {
			m := item.(map[string]interface{})
			condition := monitor.ModifyPolicyGroupCondition{}
			condition.MetricId = helper.IntInt64(m["metric_id"].(int))
			condition.AlarmNotifyType = helper.IntInt64(m["alarm_notify_type"].(int))
			condition.AlarmNotifyPeriod = helper.IntInt64(m["alarm_notify_period"].(int))
			if m["calc_type"] != nil {
				condition.CalcType = helper.IntInt64(m["calc_type"].(int))
			}
			if m["calc_value"] != nil {
				condition.CalcValue = helper.String(fmt.Sprintf("%f", m["calc_value"].(float64)))
			}
			if m["calc_period"] != nil {
				condition.CalcPeriod = helper.IntInt64(m["calc_period"].(int))
			}
			if m["continue_period"] != nil {
				condition.ContinuePeriod = helper.IntInt64(m["continue_period"].(int))
			}
			request.Conditions = append(request.Conditions, &condition)
		}
	}

	if iface, ok := d.GetOk("event_conditions"); ok {
		request.EventConditions = make([]*monitor.ModifyPolicyGroupEventCondition, 0, 10)
		for _, item := range iface.([]interface{}) {
			m := item.(map[string]interface{})
			condition := monitor.ModifyPolicyGroupEventCondition{}
			condition.EventId = helper.IntInt64(m["event_id"].(int))
			condition.AlarmNotifyType = helper.IntInt64(m["alarm_notify_type"].(int))
			condition.AlarmNotifyPeriod = helper.IntInt64(m["alarm_notify_period"].(int))
			request.EventConditions = append(request.EventConditions, &condition)
		}
	}

	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := monitorService.client.UseMonitorClient().ModifyPolicyGroup(request)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	return resourceTencentMonitorPolicyGroupRead(d, meta)
}

func resourceTencentMonitorPolicyGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_policy_group.delete")()

	var (
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request        = monitor.NewDeletePolicyGroupRequest()
	)

	groupId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("id [%s] is broken", d.Id())
	}
	request.GroupId = []*int64{&groupId}
	request.Module = helper.String("monitor")

	if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if _, err = monitorService.client.UseMonitorClient().DeletePolicyGroup(request); err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
