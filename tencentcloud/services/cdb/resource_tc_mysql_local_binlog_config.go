package cdb

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlLocalBinlogConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlLocalBinlogConfigCreate,
		Read:   resourceTencentCloudMysqlLocalBinlogConfigRead,
		Update: resourceTencentCloudMysqlLocalBinlogConfigUpdate,
		Delete: resourceTencentCloudMysqlLocalBinlogConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "实例ID，格式为cdb-c1nl9rpv。与腾讯数据库控制台显示的实例ID相同。",
			},

			"save_hours": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "本地binlog的保留期限。有效范围：72-168小时。当有灾难恢复实例时，有效范围为120-168小时。",
			},

			"max_usage": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "本地binlog的空间利用率。值范围：[30,50]。",
			},
		},
	}
}

func resourceTencentCloudMysqlLocalBinlogConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_local_binlog_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	d.SetId(d.Get("instance_id").(string))

	return resourceTencentCloudMysqlLocalBinlogConfigUpdate(d, meta)
}

func resourceTencentCloudMysqlLocalBinlogConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_local_binlog_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	localBinlogConfig, err := service.DescribeMysqlLocalBinlogConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if localBinlogConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `tencentcloud_mysql_local_binlog_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil

	}

	_ = d.Set("instance_id", instanceId)

	if localBinlogConfig.SaveHours != nil {
		_ = d.Set("save_hours", localBinlogConfig.SaveHours)
	}

	if localBinlogConfig.MaxUsage != nil {
		_ = d.Set("max_usage", localBinlogConfig.MaxUsage)
	}

	return nil
}

func resourceTencentCloudMysqlLocalBinlogConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_local_binlog_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := mysql.NewModifyLocalBinlogConfigRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if v, _ := d.GetOk("save_hours"); v != nil {
		request.SaveHours = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("max_usage"); v != nil {
		request.MaxUsage = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().ModifyLocalBinlogConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mysql localBinlogConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMysqlLocalBinlogConfigRead(d, meta)
}

func resourceTencentCloudMysqlLocalBinlogConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_local_binlog_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
