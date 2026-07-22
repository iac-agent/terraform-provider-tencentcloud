package dayuv2

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcantiddos "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/antiddos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAntiddosDefaultAlarmThreshold() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAntiddosDefaultAlarmThresholdCreate,
		Read:   resourceTencentCloudAntiddosDefaultAlarmThresholdRead,
		Update: resourceTencentCloudAntiddosDefaultAlarmThresholdUpdate,
		Delete: resourceTencentCloudAntiddosDefaultAlarmThresholdDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"default_alarm_config": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Alarm 阈值 配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alarm_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Alarm 阈值 类型，值 [1 (incoming 流量 告警 阈值) 2 (attack cleaning 流量 告警 阈值)]。",
						},
						"alarm_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Alarm 阈值，在 Mbps，使用 值 的&gt;=0; 当 使用 作为 input 参数，setting 0 将 delete 告警 阈值 配置;。",
						},
					},
				},
			},

			"instance_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "产品类型，值 [bgp (表示 advanced defense 包 product) bgpip (表示 advanced defense IP product)]。",
			},
		},
	}
}

func resourceTencentCloudAntiddosDefaultAlarmThresholdCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_default_alarm_threshold.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	d.SetId(d.Get("instance_type").(string))

	return resourceTencentCloudAntiddosDefaultAlarmThresholdUpdate(d, meta)
}

func resourceTencentCloudAntiddosDefaultAlarmThresholdRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_default_alarm_threshold.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcantiddos.NewAntiddosService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	instanceType := d.Id()

	defaultAlarmThreshold, err := service.DescribeAntiddosDefaultAlarmThresholdById(ctx, instanceType, 1)
	if err != nil {
		return err
	}

	if defaultAlarmThreshold == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `AntiddosDefaultAlarmThreshold` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if defaultAlarmThreshold != nil {
		defaultAlarmConfigMap := map[string]interface{}{}

		if defaultAlarmThreshold.AlarmType != nil {
			defaultAlarmConfigMap["alarm_type"] = defaultAlarmThreshold.AlarmType
		}

		if defaultAlarmThreshold.AlarmThreshold != nil {
			defaultAlarmConfigMap["alarm_threshold"] = defaultAlarmThreshold.AlarmThreshold
		}

		_ = d.Set("default_alarm_config", []interface{}{defaultAlarmConfigMap})
	}

	_ = d.Set("instance_type", instanceType)

	return nil
}

func resourceTencentCloudAntiddosDefaultAlarmThresholdUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_default_alarm_threshold.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := antiddos.NewCreateDefaultAlarmThresholdRequest()

	instanceType := d.Id()

	request.InstanceType = &instanceType

	if d.HasChange("default_alarm_config") {
		if dMap, ok := helper.InterfacesHeadMap(d, "default_alarm_config"); ok {
			defaultAlarmThreshold := antiddos.DefaultAlarmThreshold{}
			if v, ok := dMap["alarm_type"]; ok {
				defaultAlarmThreshold.AlarmType = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["alarm_threshold"]; ok {
				defaultAlarmThreshold.AlarmThreshold = helper.IntUint64(v.(int))
			}
			request.DefaultAlarmConfig = &defaultAlarmThreshold
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAntiddosClient().CreateDefaultAlarmThreshold(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update antiddos defaultAlarmThreshold failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudAntiddosDefaultAlarmThresholdRead(d, meta)
}

func resourceTencentCloudAntiddosDefaultAlarmThresholdDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_antiddos_default_alarm_threshold.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
