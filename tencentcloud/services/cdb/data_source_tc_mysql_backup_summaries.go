package cdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMysqlBackupSummaries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlBackupSummariesRead,
		Schema: map[string]*schema.Schema{
			"product": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "查询的云数据库产品类型，目前仅支持`mysql`。",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "指定按某项排序，可选值包括：BackupVolume：备份卷、DataBackupVolume：数据备份卷、BinlogBackupVolume：日志备份卷、AutoBackupVolume：自动备份卷、ManualBackupVolume：手动备份卷。默认情况下，它们按备份卷排序。",
			},

			"order_direction": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "指定排序方向，可选值包括：ASC：正序，DESC：逆序。默认为 ASC。",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "实例备份统计条目。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例ID。",
						},
						"auto_backup_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "该实例的自动数据备份数量。",
						},
						"auto_backup_volume": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "该实例的自动数据备份能力。",
						},
						"manual_backup_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "该实例的手动数据备份数量。",
						},
						"manual_backup_volume": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "本实例手动数据备份的容量。",
						},
						"data_backup_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例的数据备份总数（包括自动备份和手动备份）。",
						},
						"data_backup_volume": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "该实例的总数据备份容量。",
						},
						"binlog_backup_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "此实例的日志备份数。",
						},
						"binlog_backup_volume": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例日志备份的容量。",
						},
						"backup_volume": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例的总备份（包括数据备份和日志备份）占用容量。",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudMysqlBackupSummariesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_backup_summaries.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_direction"); ok {
		paramMap["OrderDirection"] = helper.String(v.(string))
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var backupSummaries []*cdb.BackupSummaryItem
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlBackupSummariesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		backupSummaries = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(backupSummaries))
	tmpList := make([]map[string]interface{}, 0, len(backupSummaries))

	if backupSummaries != nil {
		for _, backupSummaryItem := range backupSummaries {
			backupSummaryItemMap := map[string]interface{}{}

			if backupSummaryItem.InstanceId != nil {
				backupSummaryItemMap["instance_id"] = backupSummaryItem.InstanceId
			}

			if backupSummaryItem.AutoBackupCount != nil {
				backupSummaryItemMap["auto_backup_count"] = backupSummaryItem.AutoBackupCount
			}

			if backupSummaryItem.AutoBackupVolume != nil {
				backupSummaryItemMap["auto_backup_volume"] = backupSummaryItem.AutoBackupVolume
			}

			if backupSummaryItem.ManualBackupCount != nil {
				backupSummaryItemMap["manual_backup_count"] = backupSummaryItem.ManualBackupCount
			}

			if backupSummaryItem.ManualBackupVolume != nil {
				backupSummaryItemMap["manual_backup_volume"] = backupSummaryItem.ManualBackupVolume
			}

			if backupSummaryItem.DataBackupCount != nil {
				backupSummaryItemMap["data_backup_count"] = backupSummaryItem.DataBackupCount
			}

			if backupSummaryItem.DataBackupVolume != nil {
				backupSummaryItemMap["data_backup_volume"] = backupSummaryItem.DataBackupVolume
			}

			if backupSummaryItem.BinlogBackupCount != nil {
				backupSummaryItemMap["binlog_backup_count"] = backupSummaryItem.BinlogBackupCount
			}

			if backupSummaryItem.BinlogBackupVolume != nil {
				backupSummaryItemMap["binlog_backup_volume"] = backupSummaryItem.BinlogBackupVolume
			}

			if backupSummaryItem.BackupVolume != nil {
				backupSummaryItemMap["backup_volume"] = backupSummaryItem.BackupVolume
			}

			ids = append(ids, *backupSummaryItem.InstanceId)
			tmpList = append(tmpList, backupSummaryItemMap)
		}

		_ = d.Set("items", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
