package as

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAsSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsScheduleCreate,
		Read:   resourceTencentCloudAsScheduleRead,
		Update: resourceTencentCloudAsScheduleUpdate,
		Delete: resourceTencentCloudAsScheduleDelete,

		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID scaling 组。",
			},
			"schedule_action_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "名称 此 scaling 操作",
			},
			"max_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "最大 大小 对于 Auto Scaling 组。",
			},
			"min_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "最小 大小 对于 Auto Scaling 组。",
			},
			"desired_capacity": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "desired 数量 CVM 实例 该 should 是 running 在 组。",
			},
			"start_time": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAsScheduleTimestamp,
				Description:  "时间 对于 此 action 到 start, 在 \"YYYY-MM-DDThh:mm:ss+08:00\" 格式 (UTC+8).",
			},
			"end_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAsScheduleTimestamp,
				Description:  "时间 对于 此 action 到 end, 在 \"YYYY-MM-DDThh:mm:ss+08:00\" 格式 (UTC+8).",
			},
			"recurrence": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "时间 当 recurring future actions 将 start. 开始时间 是 指定 通过 用户 following Unix cron syntax 格式 And 此 argument should 是 集合 使用 end_time together。",
			},
			"disable_update_desired_capacity": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "此 flag disables normal update 的 DesiredCapacityproperty 该 would otherwise occur 当 scheduled scaling 任务 是 triggered.\n指定是否scheduled 任务 triggers proactive modification 的 DesiredCapacity 当 值 是 True. DesiredCapacity 可能 是 modified 通过 minSize 和 maxSize mechanism.\nThe following cases assume 该 DisableUpdateDesiredCapacity 是 True:\n- 当 scheduled 任务 triggered， original DesiredCapacity 是 5. scheduled 任务 changes minSize 到 10， maxSize 到 20，和 DesiredCapacity 到 15. Since DesiredCapacity update 是 已禁用，15 does 不 take effect. However， original DesiredCapacity 5 是 less 比 minSize 10，so final new DesiredCapacity 是 10.\n- 当 scheduled 任务 triggered， original DesiredCapacity 是 25. scheduled 任务 changes minSize 到 10 和 maxSize 到 20，和 DesiredCapacity 到 15. Since DesiredCapacity update 是 已禁用，15 does 不 take effect. However， original DesiredCapacity 25 是 greater 比 maxSize 20，so final new DesiredCapacity 是 20.\n- 当 scheduled 任务 triggered， original DesiredCapacity 是 13. scheduled 任务 changes minSize 到 10 和 maxSize 到 20，和 DesiredCapacity 到 15. Since DesiredCapacity update 是 已禁用，15 does 不 take effect，和 DesiredCapacity 是 still 13。",
			},
		},
	}
}

func resourceTencentCloudAsScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_schedule.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := as.NewCreateScheduledActionRequest()
	request.AutoScalingGroupId = helper.String(d.Get("scaling_group_id").(string))
	request.ScheduledActionName = helper.String(d.Get("schedule_action_name").(string))
	request.MaxSize = helper.IntUint64(d.Get("max_size").(int))
	request.MinSize = helper.IntUint64(d.Get("min_size").(int))
	request.DesiredCapacity = helper.IntUint64(d.Get("desired_capacity").(int))
	request.StartTime = helper.String(d.Get("start_time").(string))

	// end_time and recurrence must be specified at the same time
	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = helper.String(v.(string))
		if vv, ok := d.GetOk("recurrence"); ok {
			request.Recurrence = helper.String(vv.(string))
		} else {
			return fmt.Errorf("end_time and recurrence must be specified at the same time.")
		}
	} else {
		if _, ok := d.GetOk("recurrence"); ok {
			return fmt.Errorf("end_time and recurrence must be specified at the same time.")
		}
	}

	if v, ok := d.GetOkExists("disable_update_desired_capacity"); ok {
		request.DisableUpdateDesiredCapacity = helper.Bool(v.(bool))
	}

	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().CreateScheduledAction(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.ScheduledActionId == nil {
		return fmt.Errorf("schedule action id is nil")
	}
	d.SetId(*response.Response.ScheduledActionId)

	return resourceTencentCloudAsScheduleRead(d, meta)
}

func resourceTencentCloudAsScheduleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_schedule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	scheduledActionId := d.Id()
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		scheduledAction, has, e := asService.DescribeScheduledActionById(ctx, scheduledActionId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if has == 0 {
			d.SetId("")
			return nil
		}

		_ = d.Set("scaling_group_id", *scheduledAction.AutoScalingGroupId)
		_ = d.Set("schedule_action_name", *scheduledAction.ScheduledActionName)
		_ = d.Set("max_size", *scheduledAction.MaxSize)
		_ = d.Set("min_size", *scheduledAction.MinSize)
		_ = d.Set("desired_capacity", *scheduledAction.DesiredCapacity)
		_ = d.Set("start_time", *scheduledAction.StartTime)

		if scheduledAction.EndTime != nil {
			_ = d.Set("end_time", *scheduledAction.EndTime)
		}
		if scheduledAction.Recurrence != nil {
			_ = d.Set("recurrence", *scheduledAction.Recurrence)
		}
		if scheduledAction.DisableUpdateDesiredCapacity != nil {
			_ = d.Set("disable_update_desired_capacity", *scheduledAction.DisableUpdateDesiredCapacity)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudAsScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_schedule.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := as.NewModifyScheduledActionRequest()
	scheduledActionId := d.Id()
	request.ScheduledActionId = &scheduledActionId
	if d.HasChange("schedule_action_name") {
		request.ScheduledActionName = helper.String(d.Get("schedule_action_name").(string))
	}
	if d.HasChange("max_size") {
		request.MaxSize = helper.IntUint64(d.Get("max_size").(int))
	}
	if d.HasChange("min_size") {
		request.MinSize = helper.IntUint64(d.Get("min_size").(int))
	}
	if d.HasChange("desired_capacity") {
		request.DesiredCapacity = helper.IntUint64(d.Get("desired_capacity").(int))
	}
	if d.HasChange("start_time") {
		request.StartTime = helper.String(d.Get("start_time").(string))
	}
	if d.HasChange("end_time") {
		request.EndTime = helper.String(d.Get("end_time").(string))
		request.Recurrence = helper.String(d.Get("recurrence").(string))
	}
	if d.HasChange("recurrence") {
		request.Recurrence = helper.String(d.Get("recurrence").(string))
		request.EndTime = helper.String(d.Get("end_time").(string))
	}
	if d.HasChange("disable_update_desired_capacity") {
		request.DisableUpdateDesiredCapacity = helper.Bool(d.Get("disable_update_desired_capacity").(bool))
	}
	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().ModifyScheduledAction(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func resourceTencentCloudAsScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_schedule.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	scheduledActionId := d.Id()
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := asService.DeleteScheduledAction(ctx, scheduledActionId)
	if err != nil {
		return err
	}

	return nil
}
