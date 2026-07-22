package postgresql

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudPostgresqlInstanceHAConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlInstanceHAConfigCreate,
		Read:   resourceTencentCloudPostgresqlInstanceHAConfigRead,
		Update: resourceTencentCloudPostgresqlInstanceHAConfigUpdate,
		Delete: resourceTencentCloudPostgresqlInstanceHAConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID",
			},
			"sync_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(SYNC_MODE),
				Description:  "Master slave synchronization 方法，Semi-sync: Semi synchronous; Async: Asynchronous. Main 实例 默认值：Semi-sync，Read-仅 实例 默认值：Async。",
			},
			"max_standby_latency": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1073741824, 322122547200),
				Description:  "Maximum 延迟 数据 卷 对于 highly 可用 备份 machines. 当 延迟 数据 amount 的 备份 节点 是 less 比 或 equal 到 此 值，和 延迟 时间 的 备份 节点 是 less 比 或 equal 到 MaxStandbyLag，它 可以 switch 到 main 节点. 单位：byte; Parameter 范围: [1073741824，322122547200]。",
			},
			"max_standby_lag": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(5, 10),
				Description:  "Maximum 延迟 的 highly 可用 备份 machines. 当 延迟 时间 的 备份 节点 是 less 比 或 equal 到 此 值，和 amount 的 延迟 数据 的 备份 节点 是 less 比 或 equal 到 MaxStandbyLatency， primary 节点 可以 是 switched. 单位：s; Parameter 范围: [5，10]。",
			},
			"max_sync_standby_latency": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum 延迟 数据 对于 synchronous 备份. 当 amount 的 数据 delayed 通过 备份 machine 是 less 比 或 equal 到 此 值，和 延迟 时间 的 备份 machine 是 less 比 或 equal 到 MaxSyncStandbyLag，then 备份 machine adopts synchronous 复制; Otherwise，adopt asynchronous 复制. 此 参数 值 是 有效 对于 实例 其中 SyncMode 是 集合 到 Semi sync. 当 semi synchronous 实例 prohibits degradation 到 asynchronous 复制，MaxSyncStandbyLatency 和 MaxSyncStandbyLag 是 不 集合. 当 semi synchronous 实例 allow degenerate asynchronous 复制，PostgreSQL 版本 9 实例 必须 have MaxSyncStandbyLatency 集合 和 MaxSyncStandbyLag 不 集合，while PostgreSQL 版本 10 和 above 实例 必须 have MaxSyncStandbyLatency 和 MaxSyncStandbyLag 集合。",
			},
			"max_sync_standby_lag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum 延迟 时间 对于 synchronous 备份. 当 延迟 时间 的 standby machine 是 less 比 或 equal 到 此 值，和 amount 的 延迟 数据 的 standby machine 是 less 比 或 equal 到 MaxSyncStandbyLatency，then standby machine adopts synchronous 复制; Otherwise，adopt asynchronous 复制. 此 参数 值 是 有效 对于 实例 其中 SyncMode 是 集合 到 Semi sync. 当 semi synchronous 实例 prohibits degradation 到 asynchronous 复制，MaxSyncStandbyLatency 和 MaxSyncStandbyLag 是 不 集合. 当 semi synchronous 实例 allow degenerate asynchronous 复制，PostgreSQL 版本 9 实例 必须 have MaxSyncStandbyLatency 集合 和 MaxSyncStandbyLag 不 集合，while PostgreSQL 版本 10 和 above 实例 必须 have MaxSyncStandbyLatency 和 MaxSyncStandbyLag 集合。",
			},
		},
	}
}

func resourceTencentCloudPostgresqlInstanceHAConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_instance_ha_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudPostgresqlInstanceHAConfigUpdate(d, meta)
}

func resourceTencentCloudPostgresqlInstanceHAConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_instance_ha_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	haConfig, err := service.DescribePostgresqlInstanceHAConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if haConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PostgresqlInstanceHAConfig` [%s] not found.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if haConfig.SyncMode != nil {
		_ = d.Set("sync_mode", haConfig.SyncMode)
	}

	if haConfig.MaxStandbyLatency != nil {
		_ = d.Set("max_standby_latency", haConfig.MaxStandbyLatency)
	}

	if haConfig.MaxStandbyLag != nil {
		_ = d.Set("max_standby_lag", haConfig.MaxStandbyLag)
	}

	if haConfig.MaxSyncStandbyLatency != nil {
		_ = d.Set("max_sync_standby_latency", haConfig.MaxSyncStandbyLatency)
	}

	if haConfig.MaxSyncStandbyLag != nil {
		_ = d.Set("max_sync_standby_lag", haConfig.MaxSyncStandbyLag)
	}

	return nil
}

func resourceTencentCloudPostgresqlInstanceHAConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_instance_ha_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		request           = postgresql.NewModifyDBInstanceHAConfigRequest()
		instanceId        = d.Id()
		syncMode          = d.Get("sync_mode").(string)
		maxStandbyLatency = d.Get("max_standby_latency").(int)
		maxStandbyLag     = d.Get("max_standby_lag").(int)
	)

	request.DBInstanceId = &instanceId
	request.SyncMode = &syncMode
	request.MaxStandbyLatency = helper.IntUint64(maxStandbyLatency)
	request.MaxStandbyLag = helper.IntUint64(maxStandbyLag)
	if v, ok := d.GetOkExists("max_sync_standby_latency"); ok {
		request.MaxSyncStandbyLatency = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("max_sync_standby_lag"); ok {
		request.MaxSyncStandbyLag = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().ModifyDBInstanceHAConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate postgresql ModifyDBInstanceHAConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPostgresqlInstanceHAConfigRead(d, meta)
}

func resourceTencentCloudPostgresqlInstanceHAConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_instance_ha_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
