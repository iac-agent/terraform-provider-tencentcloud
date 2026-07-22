package cdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMysqlDataBackupOverview() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlDataBackupOverviewRead,
		Schema: map[string]*schema.Schema{
			"product": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "查询的云数据库产品类型，目前仅支持`mysql`。",
			},

			"data_backup_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "当前区域的数据备份总容量（包括自动备份和手动备份，单位为字节）。",
			},

			"data_backup_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "当前区域的数据备份总数。",
			},

			"auto_backup_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "当前区域的自动备份总容量。",
			},

			"auto_backup_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "当前区域自动备份的总数。",
			},

			"manual_backup_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "当前区域的手动备份总容量。",
			},

			"manual_backup_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "当前区域的手动备份总数。",
			},

			"remote_backup_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "远程备份总容量。",
			},

			"remote_backup_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "远程备份总数。",
			},

			"data_backup_archive_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "当前区域归档备份的总容量。",
			},

			"data_backup_archive_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "当前区域中的存档备份总数。",
			},

			"data_backup_standby_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "当前区域标准存储的备份总容量。",
			},

			"data_backup_standby_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "当前区域中标准存储备份的总数。",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudMysqlDataBackupOverviewRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_data_backup_overview.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	product := ""
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("product"); ok {
		product = v.(string)
		paramMap["Product"] = helper.String(v.(string))
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var dataBackupOverview *cdb.DescribeDataBackupOverviewResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlDataBackupOverviewByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		dataBackupOverview = result
		return nil
	})
	if err != nil {
		return err
	}

	if dataBackupOverview.DataBackupVolume != nil {
		_ = d.Set("data_backup_volume", dataBackupOverview.DataBackupVolume)
	}

	if dataBackupOverview.DataBackupCount != nil {
		_ = d.Set("data_backup_count", dataBackupOverview.DataBackupCount)
	}

	if dataBackupOverview.AutoBackupVolume != nil {
		_ = d.Set("auto_backup_volume", dataBackupOverview.AutoBackupVolume)
	}

	if dataBackupOverview.AutoBackupCount != nil {
		_ = d.Set("auto_backup_count", dataBackupOverview.AutoBackupCount)
	}

	if dataBackupOverview.ManualBackupVolume != nil {
		_ = d.Set("manual_backup_volume", dataBackupOverview.ManualBackupVolume)
	}

	if dataBackupOverview.ManualBackupCount != nil {
		_ = d.Set("manual_backup_count", dataBackupOverview.ManualBackupCount)
	}

	if dataBackupOverview.RemoteBackupVolume != nil {
		_ = d.Set("remote_backup_volume", dataBackupOverview.RemoteBackupVolume)
	}

	if dataBackupOverview.RemoteBackupCount != nil {
		_ = d.Set("remote_backup_count", dataBackupOverview.RemoteBackupCount)
	}

	if dataBackupOverview.DataBackupArchiveVolume != nil {
		_ = d.Set("data_backup_archive_volume", dataBackupOverview.DataBackupArchiveVolume)
	}

	if dataBackupOverview.DataBackupArchiveCount != nil {
		_ = d.Set("data_backup_archive_count", dataBackupOverview.DataBackupArchiveCount)
	}

	if dataBackupOverview.DataBackupStandbyVolume != nil {
		_ = d.Set("data_backup_standby_volume", dataBackupOverview.DataBackupStandbyVolume)
	}

	if dataBackupOverview.DataBackupStandbyCount != nil {
		_ = d.Set("data_backup_standby_count", dataBackupOverview.DataBackupStandbyCount)
	}

	d.SetId(helper.DataResourceIdsHash([]string{product}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataBackupOverview); e != nil {
			return e
		}
	}
	return nil
}
