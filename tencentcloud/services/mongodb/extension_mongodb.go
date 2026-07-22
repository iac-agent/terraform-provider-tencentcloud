package mongodb

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svcpostgresql "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/postgresql"
)

const (
	MONGODB_INSTANCE_STATUS_INITIAL    = 0
	MONGODB_INSTANCE_STATUS_PROCESSING = 1
	MONGODB_INSTANCE_STATUS_RUNNING    = 2
	MONGODB_INSTANCE_STATUS_EXPIRED    = -2

	MONGODB_ENGINE_VERSION_3_WT    = "MONGO_3_WT"
	MONGODB_ENGINE_VERSION_36_WT   = "MONGO_36_WT"
	MONGODB_ENGINE_VERSION_3_ROCKS = "MONGO_3_ROCKS"
	MONGODB_ENGINE_VERSION_4_WT    = "MONGO_40_WT"

	MONGODB_MACHINE_TYPE_GIO    = "GIO"
	MONGODB_MACHINE_TYPE_TGIO   = "TGIO"
	MONGODB_MACHINE_TYPE_HIO    = "HIO"
	MONGODB_MACHINE_TYPE_HIO10G = "HIO10G"

	MONGODB_CLUSTER_TYPE_REPLSET = "REPLSET"
	MONGODB_CLUSTER_TYPE_SHARD   = "SHARD"

	MONGO_INSTANCE_TYPE_FORMAL   = 1
	MONGO_INSTANCE_TYPE_READONLY = 3
	MONGO_INSTANCE_TYPE_STANDBY  = 4
)

var MONGODB_CLUSTER_TYPE = []string{
	MONGODB_CLUSTER_TYPE_REPLSET,
	MONGODB_CLUSTER_TYPE_SHARD,
}

const (
	MONGODB_DEFAULT_LIMIT  = 20
	MONGODB_MAX_LIMIT      = 100
	MONGODB_DEFAULT_OFFSET = 0
)

const (
	MONGODB_CHARGE_TYPE_POSTPAID = svcpostgresql.COMMON_PAYTYPE_POSTPAID
	MONGODB_CHARGE_TYPE_PREPAID  = svcpostgresql.COMMON_PAYTYPE_PREPAID
)

var MONGODB_CHARGE_TYPE = map[uint64]string{
	0: MONGODB_CHARGE_TYPE_POSTPAID,
	1: MONGODB_CHARGE_TYPE_PREPAID,
}

var MONGODB_AUTO_RENEW_FLAG = map[int]string{
	0: "NOTIFY_AND_MANUAL_RENEW",
	1: "NOTIFY_AND_AUTO_RENEW",
	2: "DISABLE_NOTIFY_AND_MANUAL_RENEW",
}

var MONGODB_PREPAID_PERIOD = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36}

const (
	MONGODB_TASK_FAILED  = "failed"
	MONGODB_TASK_PAUSED  = "paused"
	MONGODB_TASK_RUNNING = "running"
	MONGODB_TASK_SUCCESS = "success"
)

const (
	MONGODB_STATUS_DELIVERY_SUCCESS = 4
	MONGODB_STATUS_RETURN_SUCCESS   = 6
)

func TencentMongodbBasicInfo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "名称 Mongodb 实例。",
		},
		"memory": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: tccommon.ValidateIntegerMin(2),
			Description:  "Memory 大小. 最小 值 是 2，和 单位 是 GB. Memory 和 卷 必须 是 upgraded 或 degraded simultaneously。",
		},
		"volume": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: tccommon.ValidateIntegerMin(25),
			Description:  "Disk 大小. 最小 值 是 25，和 单位 是 GB. Memory 和 卷 必须 是 upgraded 或 degraded simultaneously。",
		},
		"engine_version": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Refers 到 版本 信息. DescribeSpecInfo API 可以 是 called 到 obtain detailed 信息 about 支持 versions.\n- MONGO_40_WT: 版本 的 MongoDB 4.0 WiredTiger 存储 引擎.\n- MONGO_42_WT: 版本 的 MongoDB 4.2 WiredTiger 存储 引擎.\n- MONGO_44_WT: 版本 的 MongoDB 4.4 WiredTiger 存储 引擎.\n- MONGO_50_WT: 版本 的 MongoDB 5.0 WiredTiger 存储 引擎.\n- MONGO_60_WT: 版本 的 MongoDB 6.0 WiredTiger 存储 引擎.\n- MONGO_70_WT: 版本 的 MongoDB 7.0 WiredTiger 存储 引擎.\n- MONGO_80_WT: 版本 的 MongoDB 8.0 WiredTiger 存储 引擎。",
		},
		"machine_type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			DiffSuppressFunc: func(k, olds, news string, d *schema.ResourceData) bool {
				if (olds == MONGODB_MACHINE_TYPE_GIO && news == MONGODB_MACHINE_TYPE_HIO) ||
					(olds == MONGODB_MACHINE_TYPE_HIO && news == MONGODB_MACHINE_TYPE_GIO) {
					return true
				} else if (olds == MONGODB_MACHINE_TYPE_TGIO && news == MONGODB_MACHINE_TYPE_HIO10G) ||
					(olds == MONGODB_MACHINE_TYPE_HIO10G && news == MONGODB_MACHINE_TYPE_TGIO) {
					return true
				}
				return olds == news
			},
			Description: "类型 Mongodb 实例，和 可用 值 include `HIO`(或 `GIO` 其中 将 是 已弃用，表示 high IO) 和 `HIO10G`(或 `TGIO` 其中 将 是 已弃用，表示 10-gigabit high IO)。",
		},
		"available_zone": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "可用 可用区 的 Mongodb。",
		},
		"vpc_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     "",
			Description: "ID VPC。",
		},
		"subnet_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "ID 子网 within 此 VPC. 值 为必填项 如果 `vpc_id` 是 集合。",
		},
		"project_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "ID 项目 其中 实例 belongs。",
		},
		"security_groups": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Set: func(v interface{}) int {
				return helper.HashString(v.(string))
			},
			Description: "ID 安全 组。",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "密码 的 此 Mongodb 账号",
		},
		"tags": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "标签 的 Mongodb. 键 名称 `项目` 是 系统 reserved 和 可以't 是 使用。",
		},
		"mongos_cpu": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "数量 mongos cpu。",
		},
		"mongos_memory": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Mongos 内存 大小 （GB）。",
		},
		"mongos_node_num": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "数量 mongos。",
		},
		// payment
		"charge_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      MONGODB_CHARGE_TYPE_POSTPAID,
			ValidateFunc: tccommon.ValidateAllowedStringValue([]string{MONGODB_CHARGE_TYPE_POSTPAID, MONGODB_CHARGE_TYPE_PREPAID}),
			Description:  "charge 类型 实例. 有效 值 是 `PREPAID` 和 `POSTPAID_BY_HOUR`. 默认值为 `POSTPAID_BY_HOUR`. 注意: TencentCloud International 仅 支持 `POSTPAID_BY_HOUR`. Caution 该 update operation 在 此 字段 将 delete old 实例 和 create new 一个 使用 new 计费类型",
		},
		"prepaid_period": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: tccommon.ValidateAllowedIntValue(MONGODB_PREPAID_PERIOD),
			Description:  "tenancy (时间 单位 是 month) 的 prepaid 实例. 有效 值 是 1，2，3，4，5，6，7，8，9，10，11，12，24，36. NOTE: 它 仅 works 当 charge_type 是 集合 到 `PREPAID`。",
		},
		"auto_renew_flag": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "Auto 续费标识 有效 值 是 `0`(NOTIFY_AND_MANUAL_RENEW)，`1`(NOTIFY_AND_AUTO_RENEW) 和 `2`(DISABLE_NOTIFY_AND_MANUAL_RENEW). 默认值为 `0`. 注意: 仅 works 对于 PREPAID 实例. Only 支持`0` 和 `1` 对于 creation。",
		},
		"in_maintenance": {
			Type:     schema.TypeInt,
			Optional: true,
			Description: "Switch 时间 对于 实例 配置 changes.\n" +
				"	- 0: When the adjustment is completed, perform the configuration task immediately. Default is 0.\n" +
				"	- 1: Perform reconfiguration tasks within the maintenance time window.\n" +
				"Note: Adjusting the number of nodes and slices does not support changes within the maintenance window.",
		},
		// Computed
		"status": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "状态 Mongodb 实例，和 可用 值 include pending initialization(expressed 使用 0)， processing(expressed 使用 1)，running(expressed 使用 2) 和 expired(expressed 使用 -2)。",
		},
		"vip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IP 的 Mongodb 实例。",
		},
		"vport": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "IP 端口 的 Mongodb 实例。",
		},
		"create_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "创建时间 的 Mongodb 实例。",
		},
	}
}
