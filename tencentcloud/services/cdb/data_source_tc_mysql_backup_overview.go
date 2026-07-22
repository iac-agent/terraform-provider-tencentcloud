package cdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMysqlBackupOverview() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlBackupOverviewRead,
		Schema: map[string]*schema.Schema{
			"product": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "查询的云数据库产品类型，目前仅支持`mysql`。",
			},

			"backup_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "当前区域的用户备份总数（包括数据备份和日志备份）。",
			},

			"backup_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "当前区域用户的总备份容量。",
			},

			"billing_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "当前区域用户备份的可计费容量，即超出赠送容量的部分。",
			},

			"free_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "用户在当前区域获得的免费备份容量。",
			},

			"remote_backup_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "当前区域用户异地备份总容量。注意：该字段可能返回null，表示取不到有效值。",
			},

			"backup_archive_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "归档备份能力，包括数据备份和日志备份。注意：该字段可能返回null，表示取不到有效值。",
			},

			"backup_standby_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "标准存储备份能力，包括数据备份和日志备份。注意：该字段可能返回null，表示取不到有效值。",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudMysqlBackupOverviewRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_backup_overview.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	product := ""
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("product"); ok {
		product = v.(string)
		paramMap["Product"] = helper.String(v.(string))
	}

	var backupCount *cdb.DescribeBackupOverviewResponseParams
	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlBackupOverviewByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		backupCount = result
		return nil
	})
	if err != nil {
		return err
	}

	if backupCount.BackupCount != nil {
		_ = d.Set("backup_count", backupCount.BackupCount)
	}

	if backupCount.BackupVolume != nil {
		_ = d.Set("backup_volume", backupCount.BackupVolume)
	}

	if backupCount.BillingVolume != nil {
		_ = d.Set("billing_volume", backupCount.BillingVolume)
	}

	if backupCount.FreeVolume != nil {
		_ = d.Set("free_volume", backupCount.FreeVolume)
	}

	if backupCount.RemoteBackupVolume != nil {
		_ = d.Set("remote_backup_volume", backupCount.RemoteBackupVolume)
	}

	if backupCount.BackupArchiveVolume != nil {
		_ = d.Set("backup_archive_volume", backupCount.BackupArchiveVolume)
	}

	if backupCount.BackupStandbyVolume != nil {
		_ = d.Set("backup_standby_volume", backupCount.BackupStandbyVolume)
	}

	d.SetId(helper.DataResourceIdsHash([]string{product}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), backupCount); e != nil {
			return e
		}
	}
	return nil
}
