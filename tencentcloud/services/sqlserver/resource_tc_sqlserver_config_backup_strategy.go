package sqlserver

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSqlserverConfigBackupStrategy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverConfigBackupStrategyCreate,
		Read:   resourceTencentCloudSqlserverConfigBackupStrategyRead,
		Update: resourceTencentCloudSqlserverConfigBackupStrategyUpdate,
		Delete: resourceTencentCloudSqlserverConfigBackupStrategyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID.",
			},

			"backup_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Backup 类型. 有效 值: weekly (当 长度(BackupDay) <=7 && 长度(BackupDay) >=2), daily (当 长度(BackupDay)=1). Default 值: daily.",
			},

			"backup_time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Backup 时间. Value 范围: 整数 从 0 到 23.",
			},

			"backup_day": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Backup 间隔 在 days 当 BackupType 是 daily. 当前 值 可以 仅 是 1.",
			},

			"backup_model": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Backup 模式. 有效 值: master_pkg (archive 备份 files 的 primary 节点), master_no_pkg (do 不 archive 备份 files 的 primary 节点), slave_pkg (archive 备份 files 的 副本 节点), slave_no_pkg (do 不 archive 备份 files 的 副本 节点). Backup files 的 副本 节点 是 支持 仅 当 Always On disaster recovery 是 已启用.",
			},

			"backup_cycle": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "days 的 week 在 其中 备份 将 是 performed 当 `BackupType` 是 weekly. 如果 数据 备份 retention 周期 是 less 比 7 days, 值 将 是 1-7, indicating 该 备份 将 是 performed everyday 通过 默认值; 如果 数据 备份 retention 周期 是 greater 比 或 equal 到 7 days, 值 将 是 在 least any two days, indicating 该 备份 将 是 performed 在 least twice 在 week 通过 默认值.",
			},

			"backup_save_days": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Data (日志) 备份 retention 周期. Value 范围: 3-1830 days, 默认值 值: 7 days.",
			},

			"regular_backup_enable": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Archive 备份 状态. 有效 值: 启用 (已启用); disable (已禁用). Default 值: disable.",
			},

			"regular_backup_save_days": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Archive 备份 retention days. Value 范围: 90-3650 days. Default 值: 365 days.",
			},

			"regular_backup_strategy": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Archive 备份 策略. 有效 值: years (yearly); quarters (quarterly); months(monthly); Default 值: `months`.",
			},

			"regular_backup_counts": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "数量 的 retained archive backups. Default 值: 1.",
			},

			"regular_backup_start_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Archive 备份 start date 在 YYYY-MM-DD 格式, 其中 是 当前 时间 通过 默认值.",
			},
		},
	}
}

func resourceTencentCloudSqlserverConfigBackupStrategyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_config_backup_strategy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverConfigBackupStrategyUpdate(d, meta)
}

func resourceTencentCloudSqlserverConfigBackupStrategyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_config_backup_strategy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instanceId := d.Id()

	configBackupStrategy, err := service.DescribeSqlserverConfigBackupStrategyById(ctx, instanceId)
	if err != nil {
		return err
	}

	if configBackupStrategy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverConfigBackupStrategy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if configBackupStrategy.InstanceId != nil {
		_ = d.Set("instance_id", configBackupStrategy.InstanceId)
	}

	if configBackupStrategy.BackupCycleType != nil {
		_ = d.Set("backup_type", configBackupStrategy.BackupCycleType)
		if configBackupStrategy.BackupCycleType == helper.String(SQLSERVER_BACKUP_CYCLETYPE_DAILY) {
			//Backup interval in days. When the BackupType is daily, valid value is 1.
			_ = d.Set("backup_day", 1)
		}
	}

	if configBackupStrategy.BackupTime != nil {
		_ = d.Set("backup_time", helper.StrToInt(*configBackupStrategy.BackupTime))
	}

	if configBackupStrategy.BackupModel != nil {
		_ = d.Set("backup_model", configBackupStrategy.BackupModel)
	}

	if configBackupStrategy.BackupCycle != nil {
		_ = d.Set("backup_cycle", configBackupStrategy.BackupCycle)
	}

	if configBackupStrategy.BackupSaveDays != nil {
		_ = d.Set("backup_save_days", configBackupStrategy.BackupSaveDays)
	}

	// if configBackupStrategy.RegularBackupEnable != nil {
	// 	_ = d.Set("regular_backup_enable", configBackupStrategy.RegularBackupEnable)
	// }

	// if configBackupStrategy.RegularBackupSaveDays != nil {
	// 	_ = d.Set("regular_backup_save_days", configBackupStrategy.RegularBackupSaveDays)
	// }

	// if configBackupStrategy.RegularBackupStrategy != nil {
	// 	_ = d.Set("regular_backup_strategy", configBackupStrategy.RegularBackupStrategy)
	// }

	// if configBackupStrategy.RegularBackupCounts != nil {
	// 	_ = d.Set("regular_backup_counts", configBackupStrategy.RegularBackupCounts)
	// }

	// if configBackupStrategy.RegularBackupStartTime != nil {
	// 	_ = d.Set("regular_backup_start_time", configBackupStrategy.RegularBackupStartTime)
	// }

	return nil
}

func resourceTencentCloudSqlserverConfigBackupStrategyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_config_backup_strategy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := sqlserver.NewModifyBackupStrategyRequest()

	needChange := false

	request.InstanceId = helper.String(d.Id())

	mutableArgs := []string{"backup_type", "backup_time", "backup_day", "backup_model", "backup_cycle", "backup_save_days", "regular_backup_enable", "regular_backup_save_days", "regular_backup_strategy", "regular_backup_counts", "regular_backup_start_time"}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {

		if v, ok := d.GetOk("backup_type"); ok {
			request.BackupType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("backup_model"); ok {
			request.BackupModel = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("backup_time"); ok {
			request.BackupTime = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("backup_day"); ok {
			request.BackupDay = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("backup_cycle"); ok {
			backupCycleSet := v.(*schema.Set).List()
			for i := range backupCycleSet {
				backupCycle := backupCycleSet[i].(int)
				request.BackupCycle = append(request.BackupCycle, helper.IntUint64(backupCycle))
			}
		}

		if v, ok := d.GetOkExists("backup_save_days"); ok {
			request.BackupSaveDays = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("regular_backup_enable"); ok {
			request.RegularBackupEnable = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("regular_backup_save_days"); ok {
			request.RegularBackupSaveDays = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("regular_backup_strategy"); ok {
			request.RegularBackupStrategy = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("regular_backup_counts"); ok {
			request.RegularBackupCounts = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("regular_backup_start_time"); ok {
			request.RegularBackupStartTime = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().ModifyBackupStrategy(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update sqlserver configBackupStrategy failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudSqlserverConfigBackupStrategyRead(d, meta)
}

func resourceTencentCloudSqlserverConfigBackupStrategyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_config_backup_strategy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
