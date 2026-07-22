package ssm

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func ResourceTencentCloudSsmSecretVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSsmSecretVersionCreate,
		Read:   resourceTencentCloudSsmSecretVersionRead,
		Update: resourceTencentCloudSsmSecretVersionUpdate,
		Delete: resourceTencentCloudSsmSecretVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"secret_name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "名称 secret 其中 不能 是 repeated 在 same 地域 最大 长度 是 128 bytes. 名称 可以 仅 contain English letters，numbers，underscore 和 hyphen '-'. first character 必须 是 letter 或 数量。",
			},
			"version_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "版本 的 secret. 最大 长度 是 64 bytes. version_id 可以 仅 contain English letters，numbers，underscore 和 hyphen '-'. first character 必须 是 letter 或 数量。",
			},
			"secret_binary": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"secret_string"},
				Description:  "base64-encoded binary secret. secret_binary 和 secret_string 必须 是 集合 仅 一个，和 最大 support 是 4096 bytes. 当 secret 状态 是 `已禁用`，此 字段 将 不 update anymore。",
			},
			"secret_string": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"secret_binary"},
				Description:  "字符串 text 的 secret. secret_binary 和 secret_string 必须 是 集合 仅 一个，和 最大 support 是 4096 bytes. 当 secret 状态 是 `已禁用`，此 字段 将 不 update anymore。",
			},
		},
	}
}

func resourceTencentCloudSsmSecretVersionCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_secret_version.create")()

	var (
		logId                 = tccommon.GetLogId(tccommon.ContextNil)
		ctx                   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		ssmService            = SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		outErr, inErr         error
		secretName, versionId string
	)

	param := make(map[string]interface{})
	param["secret_name"] = d.Get("secret_name").(string)
	param["version_id"] = d.Get("version_id").(string)
	if v, ok := d.GetOk("secret_binary"); ok {
		param["secret_binary"] = v.(string)
	}

	if v, ok := d.GetOk("secret_string"); ok {
		param["secret_string"] = v.(string)
	}

	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		secretName, versionId, inErr = ssmService.PutSecretValue(ctx, param)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}

		return nil
	})

	if outErr != nil {
		return outErr
	}

	d.SetId(strings.Join([]string{secretName, versionId}, tccommon.FILED_SP))
	return resourceTencentCloudSsmSecretVersionRead(d, meta)
}

func resourceTencentCloudSsmSecretVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_secret_version.read")()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		ssmService    = SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		outErr, inErr error
		secretInfo    *SecretInfo
	)

	ids := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("SSM secret version id can't read, id is borken, id is %s", d.Id())
	}
	secretName := ids[0]
	versionId := ids[1]

	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		secretInfo, inErr = ssmService.DescribeSecretByName(ctx, secretName)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}

		return nil
	})

	if outErr != nil {
		return outErr
	}

	var versionIds []string
	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		versionIds, inErr = ssmService.DescribeSecretVersionIdsByName(ctx, secretName)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}

		return nil
	})

	if outErr != nil {
		return outErr
	}

	var hasVersionId bool
	for _, id := range versionIds {
		if id == versionId {
			hasVersionId = true
			break
		}
	}

	if !hasVersionId {
		d.SetId("")
		return nil
	}

	if secretInfo.status == SSM_STATUS_ENABLED {
		var secretVersionInfo *SecretVersionInfo
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			secretVersionInfo, inErr = ssmService.DescribeSecretVersion(ctx, secretName, versionId)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}

			return nil
		})

		if outErr != nil {
			return outErr
		}

		_ = d.Set("secret_name", secretVersionInfo.secretName)
		_ = d.Set("version_id", secretVersionInfo.versionId)
		_ = d.Set("secret_binary", secretVersionInfo.secretBinary)
		_ = d.Set("secret_string", secretVersionInfo.secretString)
	}

	return nil
}

func resourceTencentCloudSsmSecretVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_secret_version.update")()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		ssmService    = SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		outErr, inErr error
		secretInfo    *SecretInfo
	)

	ids := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("SSM secret version id can't read, id is borken, id is %s", d.Id())
	}
	secretName := ids[0]
	versionId := ids[1]

	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		secretInfo, inErr = ssmService.DescribeSecretByName(ctx, secretName)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}

		return nil
	})

	if outErr != nil {
		return outErr
	}

	if secretInfo.status == SSM_STATUS_ENABLED {
		d.Partial(true)

		param := make(map[string]interface{})
		param["secret_name"] = secretName
		param["version_id"] = versionId
		if v, ok := d.GetOk("secret_binary"); ok {
			param["secret_binary"] = v.(string)
		}

		if v, ok := d.GetOk("secret_string"); ok {
			param["secret_string"] = v.(string)
		}

		if d.HasChange("secret_binary") || d.HasChange("secret_string") {
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				e := ssmService.UpdateSecret(ctx, param)
				if e != nil {
					return tccommon.RetryError(e)
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s modify SSM secret content failed, reason:%+v", logId, err)
				return err
			}
		}

		d.Partial(false)
	}

	return resourceTencentCloudSsmSecretVersionRead(d, meta)
}

func resourceTencentCloudSsmSecretVersionDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_secret_version.delete")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		ssmService = SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	ids := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("SSM secret version id can't read, id is borken, id is %s", d.Id())
	}
	secretName := ids[0]
	versionId := ids[1]

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := ssmService.DeleteSecretVersion(ctx, secretName, versionId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete SSM secret version failed, reason:%+v", logId, err)
		return err
	}

	return resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, e := ssmService.DescribeSecretVersion(ctx, secretName, versionId)
		if e != nil {
			if sdkError, ok := e.(*sdkErrors.TencentCloudSDKError); ok && sdkError.Code == "ResourceNotFound" {
				return nil
			}

			return tccommon.RetryError(err)
		}

		return resource.RetryableError(fmt.Errorf("delete fail"))
	})
}
