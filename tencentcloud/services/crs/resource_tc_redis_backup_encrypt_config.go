package crs

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudRedisBackupEncryptConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudRedisBackupEncryptConfigCreate,
		Read:   resourceTencentCloudRedisBackupEncryptConfigRead,
		Update: resourceTencentCloudRedisBackupEncryptConfigUpdate,
		Delete: resourceTencentCloudRedisBackupEncryptConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the redis instance.",
			},
			"backup_encrypt": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates whether to enable backup encryption. `true`: enable, `false`: disable.",
			},
			"encrypt_key_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Encryption key type. Valid values: `0`: TencentDB managed key (default), `1`: customer-defined KMS key.",
			},
			"kms_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "KMS key ID, format: `xxxxxxx-xxxx-xxxx-xxxx-xxxxxxx`.",
			},
			"kms_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "KMS region, format: `ap-xxxx`.",
			},
		},
	}
}

func resourceTencentCloudRedisBackupEncryptConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_encrypt.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	instanceId := d.Get("instance_id").(string)
	d.SetId(instanceId)

	return resourceTencentCloudRedisBackupEncryptConfigUpdate(d, meta)
}

func resourceTencentCloudRedisBackupEncryptConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_encrypt.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request = redis.NewDescribeBackupEncryptConfigRequest()
	)

	request.InstanceId = helper.String(d.Id())

	response := redis.NewDescribeBackupEncryptConfigResponse()
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().DescribeBackupEncryptConfigWithContext(ctx, request)
		if e != nil {
			if sdkErr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == RedisInstanceNotFound {
					log.Printf("[CRUD] redis backup_encrypt id=%s", d.Id())
					d.SetId("")
					return nil
				}
			}
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Read redis backup_encrypt config failed, Response is nil."))
		}

		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read redis backup_encrypt config failed, reason:%+v", logId, err)
		return err
	}

	if d.Id() == "" {
		return nil
	}

	if response.Response.BackupEncrypt != nil {
		_ = d.Set("backup_encrypt", response.Response.BackupEncrypt)
	}

	if response.Response.EncryptKeyType != nil {
		_ = d.Set("encrypt_key_type", response.Response.EncryptKeyType)
	}

	if response.Response.KmsKeyId != nil {
		_ = d.Set("kms_key_id", response.Response.KmsKeyId)
	}

	if response.Response.KmsRegion != nil {
		_ = d.Set("kms_region", response.Response.KmsRegion)
	}

	return nil
}

func resourceTencentCloudRedisBackupEncryptConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_encrypt.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	)

	needChange := false
	mutableArgs := []string{"backup_encrypt", "encrypt_key_type", "kms_key_id", "kms_region"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := redis.NewModifyBackupEncryptRequest()
		request.InstanceId = helper.String(d.Id())
		request.BackupEncrypt = helper.Bool(d.Get("backup_encrypt").(bool))

		if v, ok := d.GetOkExists("encrypt_key_type"); ok {
			request.EncryptKeyType = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("kms_key_id"); ok {
			request.KmsKeyId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("kms_region"); ok {
			request.KmsRegion = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().ModifyBackupEncryptWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Update redis backup_encrypt config failed, Response is nil."))
			}

			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update redis backup_encrypt config failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudRedisBackupEncryptConfigRead(d, meta)
}

func resourceTencentCloudRedisBackupEncryptConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_backup_encrypt.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
