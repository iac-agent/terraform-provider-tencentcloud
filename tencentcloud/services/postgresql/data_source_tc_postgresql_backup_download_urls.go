package postgresql

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPostgresqlBackupDownloadUrls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresqlBackupDownloadUrlsRead,
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID",
			},

			"backup_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Backup 类型 有效值：`LogBackup`，`BaseBackup`。",
			},

			"backup_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Unique 备份 ID。",
			},

			"url_expire_time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Validity 周期 的 URL，其中 是 12 hours 通过 默认值。",
			},

			"backup_download_restriction": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Backup download restriction。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"restriction_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "类型 网络 restrictions 对于 downloading 备份 files. 有效值：`NONE` (backups 可以 是 downloaded over both 私有 和 公有 networks)，`INTRANET` (backups 可以 仅 是 downloaded over 私有 网络)，`CUSTOMIZE` (backups 可以 是 downloaded over 指定 VPCs 或 在 指定 IPs)。",
						},
						"vpc_restriction_effect": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether VPC 是 allowed. 有效值：`ALLOW` (allow)，`DENY` (deny)。",
						},
						"vpc_id_set": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "是否为allowed 到 download 私有网络 ID 列表 备份 files。",
						},
						"ip_restriction_effect": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether IP 是 allowed. 有效值：`ALLOW` (allow)，`DENY` (deny)。",
						},
						"ip_set": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "是否为allowed 到 download IP 列表 备份 files。",
						},
					},
				},
			},

			"backup_download_url": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Backup download URL",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudPostgresqlBackupDownloadUrlsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_postgresql_backup_download_urls.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	var id string

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("db_instance_id"); ok {
		paramMap["DBInstanceId"] = helper.String(v.(string))
		id = v.(string)
	}

	if v, ok := d.GetOk("backup_type"); ok {
		paramMap["BackupType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("backup_id"); ok {
		paramMap["BackupId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("url_expire_time"); v != nil {
		paramMap["URLExpireTime"] = helper.IntUint64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "backup_download_restriction"); ok {
		backupDownloadRestriction := postgresql.BackupDownloadRestriction{}
		if v, ok := dMap["restriction_type"]; ok {
			backupDownloadRestriction.RestrictionType = helper.String(v.(string))
		}
		if v, ok := dMap["vpc_restriction_effect"]; ok {
			backupDownloadRestriction.VpcRestrictionEffect = helper.String(v.(string))
		}
		if v, ok := dMap["vpc_id_set"]; ok {
			vpcIdSetSet := v.(*schema.Set).List()
			backupDownloadRestriction.VpcIdSet = helper.InterfacesStringsPoint(vpcIdSetSet)
		}
		if v, ok := dMap["ip_restriction_effect"]; ok {
			backupDownloadRestriction.IpRestrictionEffect = helper.String(v.(string))
		}
		if v, ok := dMap["ip_set"]; ok {
			ipSetSet := v.(*schema.Set).List()
			backupDownloadRestriction.IpSet = helper.InterfacesStringsPoint(ipSetSet)
		}
		paramMap["BackupDownloadRestriction"] = &backupDownloadRestriction
	}

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var backupDownloadURL *string
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePostgresqlBackupDownloadUrlsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		backupDownloadURL = result
		return nil
	})
	if err != nil {
		return err
	}

	if backupDownloadURL != nil {
		_ = d.Set("backup_download_url", backupDownloadURL)
	}

	d.SetId(helper.DataResourceIdHash(id))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), backupDownloadURL); e != nil {
			return e
		}
	}
	return nil
}
