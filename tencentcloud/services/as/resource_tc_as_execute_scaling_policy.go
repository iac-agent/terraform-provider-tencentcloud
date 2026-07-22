package as

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAsExecuteScalingPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsExecuteScalingPolicyCreate,
		Read:   resourceTencentCloudAsExecuteScalingPolicyRead,
		Delete: resourceTencentCloudAsExecuteScalingPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"auto_scaling_policy_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Auto-scaling 策略 ID. 此 参数 是 不 可用 到 目标 tracking 策略。",
			},

			"honor_cooldown": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "是否check 如果 auto scaling 组 是 在 cooldown 周期 默认值：false。",
			},

			"trigger_source": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "来源 该 triggers scaling 策略. 有效值：API 和 CLOUD_MONITOR. 默认值：API. 值 CLOUD_MONITOR 是 特定 到 Cloud Monitor 服务。",
			},
		},
	}
}

func resourceTencentCloudAsExecuteScalingPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_execute_scaling_policy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = as.NewExecuteScalingPolicyRequest()
		response   = as.NewExecuteScalingPolicyResponse()
		activityId string
	)
	if v, ok := d.GetOk("auto_scaling_policy_id"); ok {
		request.AutoScalingPolicyId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("honor_cooldown"); ok {
		request.HonorCooldown = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("trigger_source"); ok {
		request.TriggerSource = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().ExecuteScalingPolicy(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate as executeScalingPolicy failed, reason:%+v", logId, err)
		return err
	}

	activityId = *response.Response.ActivityId
	d.SetId(activityId)

	return resourceTencentCloudAsExecuteScalingPolicyRead(d, meta)
}

func resourceTencentCloudAsExecuteScalingPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_execute_scaling_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudAsExecuteScalingPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_execute_scaling_policy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
