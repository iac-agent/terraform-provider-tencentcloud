package crs

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudRedisBackup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRedisBackupRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ID 实例。",
			},

			"begin_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "开始时间，such 作为 2017-02-08 19:09:26.Query 列表 backups 该 实例 started backing up during [beginTime，endTime] 时间 周期",
			},

			"end_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "结束时间，such 作为 2017-02-08 19:09:26.Query 列表 backups 该 实例 started backing up during [beginTime，endTime] 时间 周期",
			},

			"status": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "状态 备份 任务:1: Backup 是 在 process.2: 备份 是 normal.3: Backup 到 RDB 文件 processing.4: RDB conversion completed.-1: 备份 has expired.-2: Backup 删除。",
			},

			"instance_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "实例名称，其中 支持 fuzzy search based 在 实例名称",
			},

			"backup_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "An 数组 backups 对于 实例。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup 开始时间。",
						},
						"backup_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup ID。",
						},
						"backup_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup 类型1: 用户-initiated manual 备份.0: System-initiated 备份 在 early morning。",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Backup 状态1: 备份 是 locked 通过 another process.2: 备份 是 normal 和 不 locked 通过 any process.-1: 备份 has expired.3: 备份 是 being exported.4: 备份 export 是 successful。",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Notes 信息 对于 备份。",
						},
						"locked": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否backup 是 locked.0: Not locked.1: Has been locked。",
						},
						"backup_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Internal 字段，其中 可以 是 ignored 通过 用户",
						},
						"full_backup": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Internal 字段，其中 可以 是 ignored 通过 用户",
						},
						"instance_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Internal 字段，其中 可以 是 ignored 通过 用户",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 实例。",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 实例。",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 其中 备份 是 located。",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup 结束时间。",
						},
						"file_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Back up 文件 types。",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup 文件 过期时间。",
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

func dataSourceTencentCloudRedisBackupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_redis_backup.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("begin_time"); ok {
		paramMap["begin_time"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["end_time"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		statusSet := v.(*schema.Set).List()
		statusList := []*int64{}
		for i := range statusSet {
			status := statusSet[i].(int)
			statusList = append(statusList, helper.IntInt64(status))
		}
		paramMap["status"] = statusList
	}

	if v, ok := d.GetOk("instance_name"); ok {
		paramMap["InstanceName"] = helper.String(v.(string))
	}

	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var backupSet []*redis.RedisBackupSet

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRedisBackupByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		backupSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(backupSet))
	tmpList := make([]map[string]interface{}, 0, len(backupSet))

	if backupSet != nil {
		for _, redisBackupSet := range backupSet {
			redisBackupSetMap := map[string]interface{}{}

			if redisBackupSet.StartTime != nil {
				redisBackupSetMap["start_time"] = redisBackupSet.StartTime
			}

			if redisBackupSet.BackupId != nil {
				redisBackupSetMap["backup_id"] = redisBackupSet.BackupId
			}

			if redisBackupSet.BackupType != nil {
				redisBackupSetMap["backup_type"] = redisBackupSet.BackupType
			}

			if redisBackupSet.Status != nil {
				redisBackupSetMap["status"] = redisBackupSet.Status
			}

			if redisBackupSet.Remark != nil {
				redisBackupSetMap["remark"] = redisBackupSet.Remark
			}

			if redisBackupSet.Locked != nil {
				redisBackupSetMap["locked"] = redisBackupSet.Locked
			}

			if redisBackupSet.BackupSize != nil {
				redisBackupSetMap["backup_size"] = redisBackupSet.BackupSize
			}

			if redisBackupSet.FullBackup != nil {
				redisBackupSetMap["full_backup"] = redisBackupSet.FullBackup
			}

			if redisBackupSet.InstanceType != nil {
				redisBackupSetMap["instance_type"] = redisBackupSet.InstanceType
			}

			if redisBackupSet.InstanceId != nil {
				redisBackupSetMap["instance_id"] = redisBackupSet.InstanceId
			}

			if redisBackupSet.InstanceName != nil {
				redisBackupSetMap["instance_name"] = redisBackupSet.InstanceName
			}

			if redisBackupSet.Region != nil {
				redisBackupSetMap["region"] = redisBackupSet.Region
			}

			if redisBackupSet.EndTime != nil {
				redisBackupSetMap["end_time"] = redisBackupSet.EndTime
			}

			if redisBackupSet.FileType != nil {
				redisBackupSetMap["file_type"] = redisBackupSet.FileType
			}

			if redisBackupSet.ExpireTime != nil {
				redisBackupSetMap["expire_time"] = redisBackupSet.ExpireTime
			}

			ids = append(ids, *redisBackupSet.InstanceId)
			tmpList = append(tmpList, redisBackupSetMap)
		}

		_ = d.Set("backup_set", tmpList)
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
