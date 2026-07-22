package mongodb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMongodbInstanceBackups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMongodbInstanceBackupsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID， 格式 是: cmgo-9d0p6umb.Same 作为 实例 ID displayed 在 云 数据库 console 页面。",
			},

			"backup_method": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Backup 模式，currently 支持: 0-logic 备份，1-physical 备份，2-all backups. 默认为 logical 备份。",
			},

			"backup_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "备份 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID",
						},
						"backup_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Backup 模式 类型",
						},
						"backup_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup 模式 名称",
						},
						"backup_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "备注 的 备份。",
						},
						"backup_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size 的 备份(KN)。",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "开始时间 的 备份。",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "结束时间 的 备份。",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Backup 状态",
						},
						"backup_method": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Backup 方法。",
						},
						"back_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Backup 记录 ID。",
						},
						"delete_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheduled deletion 时间 对于 备份。",
						},
						"backup_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 其中 备份 是 stored (对于 cross-地域 backups)。",
						},
						"restore_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time point 支持 对于 备份 恢复。",
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

func dataSourceTencentCloudMongodbInstanceBackupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mongodb_instance_backups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("backup_method"); ok {
		paramMap["backup_method"] = helper.IntInt64(v.(int))
	}

	service := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var backupList []*mongodb.BackupInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMongodbInstanceBackupsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		backupList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(backupList))
	tmpList := make([]map[string]interface{}, 0, len(backupList))

	if backupList != nil {
		for _, backupInfo := range backupList {
			backupInfoMap := map[string]interface{}{}

			if backupInfo.InstanceId != nil {
				backupInfoMap["instance_id"] = backupInfo.InstanceId
			}

			if backupInfo.BackupType != nil {
				backupInfoMap["backup_type"] = backupInfo.BackupType
			}

			if backupInfo.BackupName != nil {
				backupInfoMap["backup_name"] = backupInfo.BackupName
			}

			if backupInfo.BackupDesc != nil {
				backupInfoMap["backup_desc"] = backupInfo.BackupDesc
			}

			if backupInfo.BackupSize != nil {
				backupInfoMap["backup_size"] = backupInfo.BackupSize
			}

			if backupInfo.StartTime != nil {
				backupInfoMap["start_time"] = backupInfo.StartTime
			}

			if backupInfo.EndTime != nil {
				backupInfoMap["end_time"] = backupInfo.EndTime
			}

			if backupInfo.Status != nil {
				backupInfoMap["status"] = backupInfo.Status
			}

			if backupInfo.BackupMethod != nil {
				backupInfoMap["backup_method"] = backupInfo.BackupMethod
			}

			if backupInfo.BackId != nil {
				backupInfoMap["back_id"] = backupInfo.BackId
			}

			if backupInfo.DeleteTime != nil {
				backupInfoMap["delete_time"] = backupInfo.DeleteTime
			}

			if backupInfo.BackupRegion != nil {
				backupInfoMap["backup_region"] = backupInfo.BackupRegion
			}

			if backupInfo.RestoreTime != nil {
				backupInfoMap["restore_time"] = backupInfo.RestoreTime
			}

			ids = append(ids, *backupInfo.InstanceId)
			tmpList = append(tmpList, backupInfoMap)
		}

		_ = d.Set("backup_list", tmpList)
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
