package mongodb

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudMongodbInstanceBackupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMongodbInstanceBackupRuleCreate,
		Read:   resourceTencentCloudMongodbInstanceBackupRuleRead,
		Update: resourceTencentCloudMongodbInstanceBackupRuleUpdate,
		Delete: resourceTencentCloudMongodbInstanceBackupRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "实例 ID",
			},

			"backup_method": {
				Required: true,
				Type:     schema.TypeInt,
				Description: "Set automatic 备份 方法. 有效 值:\n" +
					"- 0: Logical backup;\n" +
					"- 1: Physical backup;\n" +
					"- 3: Snapshot backup (supported only in cloud disk version).",
			},

			"backup_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Set 开始时间 对于 automatic 备份. 值 范围 是: [0,23]. For 示例，setting 此 参数 到 2 表示 该 备份 starts 在 02:00。",
			},

			"backup_frequency": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "指定daily automatic 备份 频率. 12: Back up twice day，approximately 12 hours apart; 24: Back up once day (默认值)，approximately 24 hours apart。",
			},

			"notify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set 是否send failure alerts 当 automatic 备份 errors occur.\n- true: Send.\n- false: Do 不 send。",
			},

			"backup_retention_period": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "指定retention 周期 对于 备份 数据. 单位：days，默认为 7 days. 取值范围：[7，365]。",
			},

			"active_weekdays": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "指定specific dates 对于 automatic backups 到 是 performed each week. 格式: Enter 数量 between 0 和 6 到 represent Sunday through Saturday (e.g.，1 表示 Monday). Separate 多个 dates 使用 commas (,). Example: Entering 1,3,5 表示 系统 将 perform backups 在 Mondays，Wednesdays，和 Fridays every week. 默认值：如果 不 集合， 默认为 full cycle (0,1,2,3,4,5,6)，meaning backups 将 是 performed daily。",
			},

			"long_term_unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Long-term retention 周期 Supports selecting 特定 dates 对于 backups 在 weekly 或 monthly basis (e.g.，备份 数据 对于 1st 和 15th 的 each month) 到 retain 对于 longer 周期 已禁用 (默认值): Long-term retention 是 已禁用 Weekly retention: 指定`weekly`. Monthly retention: 指定`monthly`。",
			},

			"long_term_active_days": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "指定specific 备份 dates 到 是 retained long-term. 此 setting 仅 takes effect 当 LongTermUnit 是 集合 到 weekly 或 monthly. Weekly Retention: Enter 数量 between 0 和 6 到 represent Sunday through Saturday. Separate 多个 dates 使用 commas. Monthly Retention: Enter 数量 between 1 和 31 到 represent 特定 dates within month. Separate 多个 dates 使用 commas。",
			},

			"long_term_expired_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Long-term 备份 retention 周期 值 范围 [30，1075]。",
			},

			"oplog_expired_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Incremental 备份 retention 周期 单位：days. 默认值：7 days. 取值范围：[7,365]。",
			},

			"backup_version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Backup 版本 Old 版本 备份 是 0，advanced 备份 是 1. Set 此 值 到 1 当 enabling advanced 备份。",
			},

			"alarm_water_level": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Sets 告警 阈值 对于 备份 dataset 存储 space usage. 单位：%. 默认值：100. 取值范围：[50，300]。",
			},
		},
	}
}

func resourceTencentCloudMongodbInstanceBackupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_backup_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		instanceId string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)
	return resourceTencentCloudMongodbInstanceBackupRuleUpdate(d, meta)
}

func resourceTencentCloudMongodbInstanceBackupRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_backup_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = mongodb.NewDescribeBackupRulesRequest()
		response   = mongodb.NewDescribeBackupRulesResponse()
		instanceId = d.Id()
	)

	request.InstanceId = &instanceId
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().DescribeBackupRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe mongodb backup rules failed, Response is nil"))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s describe mongodb backup rules failed, reason:%+v", logId, err)
		return err
	}

	_ = d.Set("instance_id", instanceId)

	if response.Response.BackupMethod != nil {
		_ = d.Set("backup_method", response.Response.BackupMethod)
	}

	if response.Response.BackupTime != nil {
		_ = d.Set("backup_time", response.Response.BackupTime)
	}

	if response.Response.BackupSaveTime != nil {
		_ = d.Set("backup_retention_period", response.Response.BackupSaveTime)
	}

	return nil
}

func resourceTencentCloudMongodbInstanceBackupRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_backup_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = mongodb.NewSetBackupRulesRequest()
		instanceId = d.Id()
	)

	if v, ok := d.GetOkExists("backup_method"); ok {
		request.BackupMethod = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("backup_time"); ok {
		request.BackupTime = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("backup_frequency"); ok {
		request.BackupFrequency = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("notify"); ok {
		request.Notify = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("backup_retention_period"); ok {
		request.BackupRetentionPeriod = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("active_weekdays"); ok {
		request.ActiveWeekdays = helper.String(v.(string))
	}

	if v, ok := d.GetOk("long_term_unit"); ok {
		request.LongTermUnit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("long_term_active_days"); ok {
		request.LongTermActiveDays = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("long_term_expired_days"); ok {
		request.LongTermExpiredDays = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("oplog_expired_days"); ok {
		request.OplogExpiredDays = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("backup_version"); ok {
		request.BackupVersion = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("alarm_water_level"); ok {
		request.AlarmWaterLevel = helper.IntInt64(v.(int))
	}

	request.InstanceId = &instanceId
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().SetBackupRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate mongodb backupRule failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMongodbInstanceBackupRuleRead(d, meta)
}

func resourceTencentCloudMongodbInstanceBackupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_backup_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
