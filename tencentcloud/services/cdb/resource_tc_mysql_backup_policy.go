package cdb

import (
	"bytes"
	"context"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudMysqlBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlBackupPolicyCreate,
		Read:   resourceTencentCloudMysqlBackupPolicyRead,
		Update: resourceTencentCloudMysqlBackupPolicyUpdate,
		Delete: resourceTencentCloudMysqlBackupPolicyDelete,

		Schema: map[string]*schema.Schema{
			"mysql_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "将应用策略的实例 ID。",
			},
			"retention_period": {
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateIntegerInRange(7, 1830),
				Optional:     true,
				Default:      7,
				Description: "备份文件的保留时间，以天为单位。最小值为 7 天，最大值为 1830 天。默认值为“7”。",
			},
			"backup_model": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      MYSQL_ALLOW_BACKUP_MODEL[0],
				ValidateFunc: tccommon.ValidateAllowedStringValue(MYSQL_ALLOW_BACKUP_MODEL),
				Description: "备份方法。支持的值包括： `physical` - 物理备份； `快照` - 快照备份。多节点仅支持“物理”，单节点仅支持“快照”。",
			},
			"backup_time": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      MYSQL_ALLOW_BACKUP_TIME[0],
				ValidateFunc: tccommon.ValidateAllowedStringValue(MYSQL_ALLOW_BACKUP_TIME),
				Description: "实例备份时间，格式为“HH:mm-HH:mm”。时间设置间隔为四小时。默认为“02:00-06:00”。可以支持以下值：`02:00-06:00`、`06:00-10:00`、`10:00-14:00`、`14:00-18:00`、`18:00-22:00`和`22:00-02:00`。",
			},

			"binlog_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Binlog 保留时间，以天为单位。最小值为 7 天，最大值为 1830 天。该值不能设置大于备份文件保留时间。",
			},

			"enable_binlog_standby": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "off",
				Description: "是否启用日志备份标准存储策略，“off”-关闭，“on”-打开，默认关闭。",
			},

			"binlog_standby_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "日志备份存储的标准起始天数。当日志备份达到标准起始存储天数时将进行转换。最短为 30 天，且不得大于日志备份保留天数。",
			},
		},
	}
}

func resourceTencentCloudMysqlBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_backup_policy.create")()

	d.SetId(d.Get("mysql_id").(string))

	return resourceTencentCloudMysqlBackupPolicyUpdate(d, meta)
}

func resourceTencentCloudMysqlBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_backup_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	mysqlService := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		desResponse, e := mysqlService.DescribeBackupConfigByMysqlId(ctx, d.Id())
		if e != nil {
			if mysqlService.NotFoundMysqlInstance(e) {
				d.SetId("")
				return nil
			}
			return tccommon.RetryError(e)
		}
		_ = d.Set("mysql_id", d.Id())
		_ = d.Set("retention_period", int(*desResponse.Response.BackupExpireDays))
		_ = d.Set("backup_model", *desResponse.Response.BackupMethod)
		var buf bytes.Buffer

		if *desResponse.Response.StartTimeMin < 10 {
			buf.WriteString("0")
		}
		buf.WriteString(fmt.Sprintf("%d:00-", *desResponse.Response.StartTimeMin))

		if *desResponse.Response.StartTimeMax < 10 {
			buf.WriteString("0")
		}
		buf.WriteString(fmt.Sprintf("%d:00", *desResponse.Response.StartTimeMax))
		_ = d.Set("backup_time", buf.String())
		_ = d.Set("binlog_period", int(*desResponse.Response.BinlogExpireDays))

		if desResponse.Response.EnableBinlogStandby != nil {
			_ = d.Set("enable_binlog_standby", desResponse.Response.EnableBinlogStandby)
		}

		if desResponse.Response.BinlogStandbyDays != nil {
			_ = d.Set("binlog_standby_days", desResponse.Response.BinlogStandbyDays)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("[API]Describe mysql backup policy fail,reason:%s", err.Error())
	}
	return nil
}

func resourceTencentCloudMysqlBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_backup_policy.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	mysqlService := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		mysqlId             = d.Get("mysql_id").(string)
		retentionPeriod     = int64(d.Get("retention_period").(int))
		backupModel         = d.Get("backup_model").(string)
		backupTime          = d.Get("backup_time").(string)
		binlogExpireDays    int64
		enableBinlogStandby string
		binlogStandbyDays   int64
	)

	if v, ok := d.GetOkExists("binlog_period"); ok {
		binlogExpireDays = int64(v.(int))
	}

	if v, ok := d.GetOk("enable_binlog_standby"); ok {
		enableBinlogStandby = v.(string)
	}

	if v, ok := d.GetOkExists("binlog_standby_days"); ok {
		binlogStandbyDays = int64(v.(int))
	}

	if d.HasChange("retention_period") || d.HasChange("backup_model") || d.HasChange("backup_time") ||
		d.HasChange("binlog_period") || d.HasChange("enable_binlog_standby") || d.HasChange("binlog_standby_days") {
		err := mysqlService.ModifyBackupConfigByMysqlId(ctx, mysqlId, retentionPeriod, backupModel, backupTime, binlogExpireDays, enableBinlogStandby, binlogStandbyDays)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudMysqlBackupPolicyRead(d, meta)
}

// set all config to default
func resourceTencentCloudMysqlBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_backup_policy.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	mysqlService := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		retentionPeriod     int64 = 7
		backupModel               = MYSQL_ALLOW_BACKUP_MODEL[1]
		backupTime                = MYSQL_ALLOW_BACKUP_TIME[0]
		binlogExpireDays    int64 = 7
		enableBinlogStandby       = "off"
		binlogStandbyDays   int64 = 180
	)
	err := mysqlService.ModifyBackupConfigByMysqlId(ctx, d.Id(), retentionPeriod, backupModel, backupTime, binlogExpireDays, enableBinlogStandby, binlogStandbyDays)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
