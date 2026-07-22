package dcdb

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDcdbEncryptAttributesConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbEncryptAttributesConfigCreate,
		Read:   resourceTencentCloudDcdbEncryptAttributesConfigRead,
		Update: resourceTencentCloudDcdbEncryptAttributesConfigUpdate,
		Delete: resourceTencentCloudDcdbEncryptAttributesConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID",
			},

			"encrypt_enabled": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "是否enable 数据 加密. Notice: 它 是 不 支持 到 turn 它 关闭 after 它 是 turned 在. 可选 值: 0-disable，1-启用。",
			},
		},
	}
}

func resourceTencentCloudDcdbEncryptAttributesConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_encrypt_attributes_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}
	d.SetId(instanceId)

	return resourceTencentCloudDcdbEncryptAttributesConfigUpdate(d, meta)
}

func resourceTencentCloudDcdbEncryptAttributesConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_encrypt_attributes_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	encryptAttributesConfig, err := service.DescribeDcdbEncryptAttributesConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if encryptAttributesConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DcdbEncryptAttributesConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if encryptAttributesConfig.EncryptStatus != nil {
		_ = d.Set("encrypt_enabled", encryptAttributesConfig.EncryptStatus)
	}

	return nil
}

func resourceTencentCloudDcdbEncryptAttributesConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_encrypt_attributes_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := dcdb.NewModifyDBEncryptAttributesRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if d.HasChange("encrypt_enabled") {
		if v, ok := d.GetOkExists("encrypt_enabled"); ok {
			request.EncryptEnabled = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDcdbClient().ModifyDBEncryptAttributes(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dcdb encryptAttributesConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDcdbEncryptAttributesConfigRead(d, meta)
}

func resourceTencentCloudDcdbEncryptAttributesConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dcdb_encrypt_attributes_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
