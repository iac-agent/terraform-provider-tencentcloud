package cdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMysqlBinlogBackupOverview() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlBinlogBackupOverviewRead,
		Schema: map[string]*schema.Schema{
			"product": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "查询的云数据库产品类型，目前仅支持`mysql`。",
			},

			"binlog_backup_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "总日志备份容量，包括异地日志备份（单位为字节）。",
			},

			"binlog_backup_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "日志备份总数，包括远程日志备份。",
			},

			"remote_binlog_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "远程日志备份容量（以字节为单位）。",
			},

			"remote_binlog_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "远程日志备份的数量。",
			},

			"binlog_archive_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "归档日志备份容量（以字节为单位）。",
			},

			"binlog_archive_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "归档日志备份的数量。",
			},

			"binlog_standby_volume": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "标准存储日志备份容量（以字节为单位）。",
			},

			"binlog_standby_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "标准存储日志备份的数量。",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudMysqlBinlogBackupOverviewRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_binlog_backup_overview.read")()
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
	var binlogBackupOverview *cdb.DescribeBinlogBackupOverviewResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlBinlogBackupOverviewByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		binlogBackupOverview = result
		return nil
	})
	if err != nil {
		return err
	}

	if binlogBackupOverview.BinlogBackupVolume != nil {
		_ = d.Set("binlog_backup_volume", binlogBackupOverview.BinlogBackupVolume)
	}

	if binlogBackupOverview.BinlogBackupCount != nil {
		_ = d.Set("binlog_backup_count", binlogBackupOverview.BinlogBackupCount)
	}

	if binlogBackupOverview.RemoteBinlogVolume != nil {
		_ = d.Set("remote_binlog_volume", binlogBackupOverview.RemoteBinlogVolume)
	}

	if binlogBackupOverview.RemoteBinlogCount != nil {
		_ = d.Set("remote_binlog_count", binlogBackupOverview.RemoteBinlogCount)
	}

	if binlogBackupOverview.BinlogArchiveVolume != nil {
		_ = d.Set("binlog_archive_volume", binlogBackupOverview.BinlogArchiveVolume)
	}

	if binlogBackupOverview.BinlogArchiveCount != nil {
		_ = d.Set("binlog_archive_count", binlogBackupOverview.BinlogArchiveCount)
	}

	if binlogBackupOverview.BinlogStandbyVolume != nil {
		_ = d.Set("binlog_standby_volume", binlogBackupOverview.BinlogStandbyVolume)
	}

	if binlogBackupOverview.BinlogStandbyCount != nil {
		_ = d.Set("binlog_standby_count", binlogBackupOverview.BinlogStandbyCount)
	}

	d.SetId(helper.DataResourceIdsHash([]string{product}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), binlogBackupOverview); e != nil {
			return e
		}
	}
	return nil
}
