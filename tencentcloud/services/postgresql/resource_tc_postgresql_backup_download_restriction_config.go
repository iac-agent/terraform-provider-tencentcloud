package postgresql

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresqlBackupDownloadRestrictionConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlBackupDownloadRestrictionConfigCreate,
		Read:   resourceTencentCloudPostgresqlBackupDownloadRestrictionConfigRead,
		Update: resourceTencentCloudPostgresqlBackupDownloadRestrictionConfigUpdate,
		Delete: resourceTencentCloudPostgresqlBackupDownloadRestrictionConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"restriction_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Backup 文件 download restriction 类型: NONE:Unlimited，both 内部 和 外部 networks 可以 是 downloaded. INTRANET:Only intranet downloads 是 allowed. CUSTOMIZE:Customize vpc 或 ip 该 limits downloads。",
			},

			"vpc_restriction_effect": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "vpc 限制 Strategy: ALLOW，DENY。",
			},

			"vpc_id_set": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "列表 vpcIds 该 allow 或 deny downloading 的 备份 files。",
			},

			"ip_restriction_effect": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ip 限制 Strategy: ALLOW，DENY。",
			},

			"ip_set": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "列表 ips 该 是 allowed 或 denied 到 download 备份 files。",
			},
		},
	}
}

func resourceTencentCloudPostgresqlBackupDownloadRestrictionConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_backup_download_restriction_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var resType string
	if v, ok := d.GetOk("restriction_type"); ok {
		resType = v.(string)
	}

	d.SetId(resType)

	return resourceTencentCloudPostgresqlBackupDownloadRestrictionConfigUpdate(d, meta)
}

func resourceTencentCloudPostgresqlBackupDownloadRestrictionConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_backup_download_restriction_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	resType := d.Id()

	BackupDownloadRestrictionConfig, err := service.DescribePostgresqlBackupDownloadRestrictionConfigById(ctx, resType)
	if err != nil {
		return err
	}

	if BackupDownloadRestrictionConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PostgresqlBackupDownloadRestrictionConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if BackupDownloadRestrictionConfig.RestrictionType != nil {
		_ = d.Set("restriction_type", BackupDownloadRestrictionConfig.RestrictionType)
	}

	if BackupDownloadRestrictionConfig.VpcRestrictionEffect != nil {
		_ = d.Set("vpc_restriction_effect", BackupDownloadRestrictionConfig.VpcRestrictionEffect)
	}

	if BackupDownloadRestrictionConfig.VpcIdSet != nil {
		_ = d.Set("vpc_id_set", BackupDownloadRestrictionConfig.VpcIdSet)
	}

	if BackupDownloadRestrictionConfig.IpRestrictionEffect != nil {
		_ = d.Set("ip_restriction_effect", BackupDownloadRestrictionConfig.IpRestrictionEffect)
	}

	if BackupDownloadRestrictionConfig.IpSet != nil {
		_ = d.Set("ip_set", BackupDownloadRestrictionConfig.IpSet)
	}

	return nil
}

func resourceTencentCloudPostgresqlBackupDownloadRestrictionConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_backup_download_restriction_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := postgresql.NewModifyBackupDownloadRestrictionRequest()

	resType := d.Id()

	if d.HasChange("restriction_type") {
		if v, ok := d.GetOk("restriction_type"); ok {
			resType = v.(string)
		}
	}
	request.RestrictionType = &resType

	if d.HasChange("vpc_restriction_effect") {
		if v, ok := d.GetOk("vpc_restriction_effect"); ok {
			request.VpcRestrictionEffect = helper.String(v.(string))
		}
	}

	if d.HasChange("vpc_id_set") {
		if v, ok := d.GetOk("vpc_id_set"); ok {
			vpcIdSetSet := v.(*schema.Set).List()
			for i := range vpcIdSetSet {
				if vpcIdSetSet[i] != nil {
					vpcIdSet := vpcIdSetSet[i].(string)
					request.VpcIdSet = append(request.VpcIdSet, &vpcIdSet)
				}
			}
		}
	}

	if d.HasChange("ip_restriction_effect") {
		if v, ok := d.GetOk("ip_restriction_effect"); ok {
			request.IpRestrictionEffect = helper.String(v.(string))
		}
	}

	if d.HasChange("ip_set") {
		if v, ok := d.GetOk("ip_set"); ok {
			ipSetSet := v.(*schema.Set).List()
			for i := range ipSetSet {
				if ipSetSet[i] != nil {
					ipSet := ipSetSet[i].(string)
					request.IpSet = append(request.IpSet, &ipSet)
				}
			}
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().ModifyBackupDownloadRestriction(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update postgresql BackupDownloadRestrictionConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPostgresqlBackupDownloadRestrictionConfigRead(d, meta)
}

func resourceTencentCloudPostgresqlBackupDownloadRestrictionConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_backup_download_restriction_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
