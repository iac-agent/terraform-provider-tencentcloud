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

func ResourceTencentCloudAsScalingPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsScalingPolicyCreate,
		Read:   resourceTencentCloudAsScalingPolicyRead,
		Update: resourceTencentCloudAsScalingPolicyUpdate,
		Delete: resourceTencentCloudAsScalingPolicyDelete,

		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID scaling 组。",
			},
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "名称 策略 用于define reaction 当 告警 是 triggered。",
			},
			"adjustment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(SCALING_GROUP_ADJUSTMENT_TYPE),
				Description:  "指定是否adjustment 是 absolute 数量 或 percentage 的 当前 容量. 有效值：`CHANGE_IN_CAPACITY`，`EXACT_CAPACITY` 和 `PERCENT_CHANGE_IN_CAPACITY`。",
			},
			"adjustment_value": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Define 数量 实例 通过 其中 到 scale.For `CHANGE_IN_CAPACITY` 类型 或 PERCENT_CHANGE_IN_CAPACITY， positive increment adds 到 当前 容量 和 negative 值 removes 从 当前 容量. For `EXACT_CAPACITY` 类型，它 defines absolute 数量 existing Auto Scaling 组 大小。",
			},
			"comparison_operator": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(SCALING_GROUP_COMPARISON_OPERATOR),
				Description:  "比较运算符 有效值：`GREATER_THAN`，`GREATER_THAN_OR_EQUAL_TO`，`LESS_THAN`，`LESS_THAN_OR_EQUAL_TO`，`EQUAL_TO` 和 `NOT_EQUAL_TO`。",
			},
			"metric_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(SCALING_GROUP_METRIC_NAME),
				Description:  "名称 indicator. 有效值：`CPU_UTILIZATION`，`MEM_UTILIZATION`，`LAN_TRAFFIC_OUT`，`LAN_TRAFFIC_IN`，`WAN_TRAFFIC_OUT` 和 `WAN_TRAFFIC_IN`。",
			},
			"threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Alarm 阈值。",
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{60, 300}),
				Description:  "Time 周期 在 second. 有效值：`60` 和 `300`。",
			},
			"continuous_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 10),
				Description:  "重试次数 有效 值 ranges: (1~10)。",
			},
			"statistic": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(SCALING_GROUP_STATISTIC),
				Description:  "Statistic types. 有效值：`AVERAGE`，`MAXIMUM` 和 `MINIMUM`. 默认为 `AVERAGE`。",
			},
			"cooldown": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Cooldwon 时间 在 second. 默认为 `300`。",
			},
			"notification_user_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "An ID 组 的 users 到 是 notified 当 告警 是 triggered。",
			},
			"policy_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Alarm triggering 策略 类型， 默认值 类型 是 SIMPLE. 取值范围：SIMPLE: Simple 策略; TARGET_TRACKING: Target tracking 策略。",
			},
			"predefined_metric_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Predefined 监控 items，applicable 仅 到 目标 tracking policies，和 必填 在 目标 tracking 策略 scenarios. 取值范围：ASG_AVG_CPU_UTILIZATION: Average CPU utilization; ASG_AVG_LAN_TRAFFIC_OUT: Average intranet outbound 带宽; ASG_AVG_LAN_TRAFFIC_IN: Average intranet inbound 带宽; ASG_AVG_WAN_TRAFFIC_OUT: Average internet outbound 带宽; ASG_AVG_WAN_TRAFFIC_IN: Average internet inbound 带宽。",
			},
			"target_value": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Target 值，applicable 仅 到 目标 tracking strategies，和 必填 在 目标 tracking strategy scenarios. ASG_AVG_CPU_UTILIZATION: [1，100)，单位：%; ASG_AVG_LAN_TRAFFIC_OUT: >0，单位：Mbps; ASG_AVG_LAN_TRAFFIC_IN: >0，单位：Mbps; ASG_AVG_WAN_TRAFFIC_OUT: >0，单位：Mbps; ASG_AVG_WAN_TRAFFIC_IN: >0，单位：Mbps。",
			},
			"estimated_instance_warmup": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "实例 warm-up 时间，（秒）， applicable 仅 到 目标 tracking strategies. 取值范围为 0-3600，使用 默认值 warm-up 时间 的 300 秒。",
			},
			"disable_scale_in": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "是否disable scaling down applies 仅 到 目标 tracking strategy; 默认值为 false. 取值范围：true: 目标 tracking strategy 仅 triggers scaling up; false: 目标 tracking strategy triggers both scaling up 和 scaling down。",
			},
		},
	}
}

func resourceTencentCloudAsScalingPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_policy.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := as.NewCreateScalingPolicyRequest()
	request.AutoScalingGroupId = helper.String(d.Get("scaling_group_id").(string))
	request.ScalingPolicyName = helper.String(d.Get("policy_name").(string))
	if v, ok := d.GetOk("adjustment_type"); ok {
		request.AdjustmentType = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("adjustment_value"); ok {
		request.AdjustmentValue = helper.IntInt64(v.(int))
	}
	metricAlarm := &as.MetricAlarm{}
	var hasMetricAlarm bool
	if v, ok := d.GetOk("comparison_operator"); ok {
		metricAlarm.ComparisonOperator = helper.String(v.(string))
		hasMetricAlarm = true
	}
	if v, ok := d.GetOk("metric_name"); ok {
		metricAlarm.MetricName = helper.String(v.(string))
		hasMetricAlarm = true
	}
	if v, ok := d.GetOkExists("threshold"); ok {
		metricAlarm.Threshold = helper.IntUint64(v.(int))
		hasMetricAlarm = true
	}
	if v, ok := d.GetOkExists("period"); ok {
		metricAlarm.Period = helper.IntUint64(v.(int))
		hasMetricAlarm = true
	}
	if v, ok := d.GetOkExists("continuous_time"); ok {
		metricAlarm.ContinuousTime = helper.IntUint64(v.(int))
		hasMetricAlarm = true
	}
	if v, ok := d.GetOk("statistic"); ok {
		metricAlarm.Statistic = helper.String(v.(string))
		hasMetricAlarm = true
	}
	if hasMetricAlarm {
		request.MetricAlarm = metricAlarm
	}
	if v, ok := d.GetOk("cooldown"); ok {
		request.Cooldown = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("notification_user_group_ids"); ok {
		notificationUserGroupIds := v.([]interface{})
		request.NotificationUserGroupIds = make([]*string, 0, len(notificationUserGroupIds))
		for _, value := range notificationUserGroupIds {
			request.NotificationUserGroupIds = append(request.NotificationUserGroupIds, helper.String(value.(string)))
		}
	}

	if v, ok := d.GetOk("policy_type"); ok {
		request.ScalingPolicyType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("predefined_metric_type"); ok {
		request.PredefinedMetricType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("target_value"); ok {
		request.TargetValue = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("estimated_instance_warmup"); ok {
		request.EstimatedInstanceWarmup = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("disable_scale_in"); ok {
		request.DisableScaleIn = helper.Bool(v.(bool))
	}

	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().CreateScalingPolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.AutoScalingPolicyId == nil {
		return fmt.Errorf("scaling policy id is nil")
	}
	d.SetId(*response.Response.AutoScalingPolicyId)

	return resourceTencentCloudAsScalingPolicyRead(d, meta)
}

func resourceTencentCloudAsScalingPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	scalingPolicyId := d.Id()
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		scalingPolicy, has, e := asService.DescribeScalingPolicyById(ctx, scalingPolicyId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if has == 0 {
			d.SetId("")
			return nil
		}

		if scalingPolicy.AutoScalingGroupId != nil {
			_ = d.Set("scaling_group_id", *scalingPolicy.AutoScalingGroupId)
		}
		if scalingPolicy.ScalingPolicyName != nil {
			_ = d.Set("policy_name", *scalingPolicy.ScalingPolicyName)
		}
		if scalingPolicy.AdjustmentType != nil {
			_ = d.Set("adjustment_type", *scalingPolicy.AdjustmentType)
		}
		if scalingPolicy.AdjustmentValue != nil {
			_ = d.Set("adjustment_value", *scalingPolicy.AdjustmentValue)
		}
		if scalingPolicy.MetricAlarm != nil {
			if scalingPolicy.MetricAlarm.ComparisonOperator != nil {
				_ = d.Set("comparison_operator", *scalingPolicy.MetricAlarm.ComparisonOperator)
			}
			if scalingPolicy.MetricAlarm.MetricName != nil {
				_ = d.Set("metric_name", *scalingPolicy.MetricAlarm.MetricName)
			}
			if scalingPolicy.MetricAlarm.Threshold != nil {
				_ = d.Set("threshold", *scalingPolicy.MetricAlarm.Threshold)
			}
			if scalingPolicy.MetricAlarm.Period != nil {
				_ = d.Set("period", *scalingPolicy.MetricAlarm.Period)
			}
			if scalingPolicy.MetricAlarm.ContinuousTime != nil {
				_ = d.Set("continuous_time", *scalingPolicy.MetricAlarm.ContinuousTime)
			}
			if scalingPolicy.MetricAlarm.Statistic != nil {
				_ = d.Set("statistic", *scalingPolicy.MetricAlarm.Statistic)
			}
		}
		if scalingPolicy.Cooldown != nil {
			_ = d.Set("cooldown", *scalingPolicy.Cooldown)
		}
		if scalingPolicy.NotificationUserGroupIds != nil {
			_ = d.Set("notification_user_group_ids", helper.StringsInterfaces(scalingPolicy.NotificationUserGroupIds))
		}
		if scalingPolicy.ScalingPolicyType != nil {
			_ = d.Set("policy_type", *scalingPolicy.ScalingPolicyType)
		}
		if scalingPolicy.PredefinedMetricType != nil {
			_ = d.Set("predefined_metric_type", *scalingPolicy.PredefinedMetricType)
		}
		if scalingPolicy.TargetValue != nil {
			_ = d.Set("target_value", *scalingPolicy.TargetValue)
		}
		if scalingPolicy.EstimatedInstanceWarmup != nil {
			_ = d.Set("estimated_instance_warmup", *scalingPolicy.EstimatedInstanceWarmup)
		}
		if scalingPolicy.DisableScaleIn != nil {
			_ = d.Set("disable_scale_in", *scalingPolicy.DisableScaleIn)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func resourceTencentCloudAsScalingPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_policy.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := as.NewModifyScalingPolicyRequest()
	scalingPolicyId := d.Id()
	request.AutoScalingPolicyId = &scalingPolicyId
	if d.HasChange("policy_name") {
		request.ScalingPolicyName = helper.String(d.Get("policy_name").(string))
	}
	if d.HasChange("adjustment_type") {
		if v, ok := d.GetOk("adjustment_type"); ok {
			request.AdjustmentType = helper.String(v.(string))
		}
	}
	if d.HasChange("adjustment_value") {
		if v, ok := d.GetOkExists("adjustment_value"); ok {
			request.AdjustmentValue = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("comparison_operator") || d.HasChange("threshold") || d.HasChange("metric_name") || d.HasChange("period") || d.HasChange("continuous_time") || d.HasChange("statistic") {
		request.MetricAlarm = &as.MetricAlarm{}

		if v, ok := d.GetOk("comparison_operator"); ok {
			request.MetricAlarm.ComparisonOperator = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("threshold"); ok {
			request.MetricAlarm.Threshold = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("metric_name"); ok {
			request.MetricAlarm.MetricName = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("period"); ok {
			request.MetricAlarm.Period = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("continuous_time"); ok {
			request.MetricAlarm.ContinuousTime = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("statistic"); ok {
			request.MetricAlarm.Statistic = helper.String(v.(string))
		}
	}

	if d.HasChange("cooldown") {
		request.Cooldown = helper.IntUint64(d.Get("cooldown").(int))
	}
	if d.HasChange("notification_user_group_ids") {
		notificationUserGroupIds := d.Get("notification_user_group_ids").([]interface{})
		request.NotificationUserGroupIds = make([]*string, 0, len(notificationUserGroupIds))
		for _, value := range notificationUserGroupIds {
			request.NotificationUserGroupIds = append(request.NotificationUserGroupIds, helper.String(value.(string)))
		}
	}

	if v, ok := d.GetOk("predefined_metric_type"); ok {
		request.PredefinedMetricType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("target_value"); ok {
		request.TargetValue = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("estimated_instance_warmup"); ok {
		request.EstimatedInstanceWarmup = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("disable_scale_in"); ok {
		request.DisableScaleIn = helper.Bool(v.(bool))
	}

	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().ModifyScalingPolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func resourceTencentCloudAsScalingPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_policy.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	scalingPolicyId := d.Id()
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := asService.DeleteScalingPolicy(ctx, scalingPolicyId)
	if err != nil {
		return err
	}

	return nil
}
