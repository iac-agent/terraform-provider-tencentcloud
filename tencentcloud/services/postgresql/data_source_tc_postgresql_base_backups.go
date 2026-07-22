package postgresql

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPostgresqlBaseBackups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresqlBaseBackupsRead,
		Schema: map[string]*schema.Schema{
			"min_finish_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Minimum 结束时间 的 备份 在 格式 的 `2018-01-01 00:00:00`. It 是 7 days ago 通过 默认值。",
			},

			"max_finish_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Maximum 结束时间 的 备份 在 格式 的 `2018-01-01 00:00:00`. It 是 当前 时间 通过 默认值。",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 实例 使用 一个 或 more criteria. 有效 过滤器 names: `db-实例-ID`: 过滤器 通过 实例 ID (在 字符串 格式). `db-实例-名称`: 过滤器 通过 实例名称 (在 字符串 格式). `db-实例-ip`: 过滤器 通过 实例 VPC IP (在 字符串 格式). `base-备份-ID`: 过滤器 通过 base 备份 ID (在 字符串 格式)。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "过滤名称",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "一个或多个过滤值",
						},
					},
				},
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting 字段. 有效值：`StartTime`，`FinishTime`，`Size`。",
			},

			"order_by_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting 顺序 有效值：`asc` (ascending)，`desc` (descending)。",
			},

			"base_backup_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 full 备份 details。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique ID 备份 文件。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup 文件 名称",
						},
						"backup_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup 方法，包括 physical 和 logical。",
						},
						"backup_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup 模式，包括 automatic 和 manual。",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup 任务 状态",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Backup 集合 大小 在 bytes。",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup 开始时间。",
						},
						"finish_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup 结束时间。",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup 过期时间。",
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

func dataSourceTencentCloudPostgresqlBaseBackupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_postgresql_base_backups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("min_finish_time"); ok {
		paramMap["MinFinishTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("max_finish_time"); ok {
		paramMap["MaxFinishTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*postgresql.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := postgresql.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["Filters"] = tmpSet
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by_type"); ok {
		paramMap["OrderByType"] = helper.String(v.(string))
	}

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var baseBackupSet []*postgresql.BaseBackup

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePostgresqlBaseBackupsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		baseBackupSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(baseBackupSet))
	tmpList := make([]map[string]interface{}, 0, len(baseBackupSet))

	if baseBackupSet != nil {
		for _, baseBackup := range baseBackupSet {
			baseBackupMap := map[string]interface{}{}

			if baseBackup.DBInstanceId != nil {
				baseBackupMap["db_instance_id"] = baseBackup.DBInstanceId
			}

			if baseBackup.Id != nil {
				baseBackupMap["id"] = baseBackup.Id
			}

			if baseBackup.Name != nil {
				baseBackupMap["name"] = baseBackup.Name
			}

			if baseBackup.BackupMethod != nil {
				baseBackupMap["backup_method"] = baseBackup.BackupMethod
			}

			if baseBackup.BackupMode != nil {
				baseBackupMap["backup_mode"] = baseBackup.BackupMode
			}

			if baseBackup.State != nil {
				baseBackupMap["state"] = baseBackup.State
			}

			if baseBackup.Size != nil {
				baseBackupMap["size"] = baseBackup.Size
			}

			if baseBackup.StartTime != nil {
				baseBackupMap["start_time"] = baseBackup.StartTime
			}

			if baseBackup.FinishTime != nil {
				baseBackupMap["finish_time"] = baseBackup.FinishTime
			}

			if baseBackup.ExpireTime != nil {
				baseBackupMap["expire_time"] = baseBackup.ExpireTime
			}

			ids = append(ids, *baseBackup.DBInstanceId)
			tmpList = append(tmpList, baseBackupMap)
		}

		_ = d.Set("base_backup_set", tmpList)
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
