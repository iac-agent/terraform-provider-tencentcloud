package cynosdb

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svcpostgresql "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/postgresql"
)

const (
	CYNOSDB_CHARGE_TYPE_POSTPAID = svcpostgresql.COMMON_PAYTYPE_POSTPAID
	CYNOSDB_CHARGE_TYPE_PREPAID  = svcpostgresql.COMMON_PAYTYPE_PREPAID
	CYNOSDB_SERVERLESS           = "SERVERLESS"
	CYNOSDB_NORMAL               = "NORMAL"

	CYNOSDB_STATUS_RUNNING  = "running"
	CYNOSDB_STATUS_OFFLINE  = "offlined"
	CYNOSDB_STATUS_ISOLATED = "isolated"
	CYNOSDB_STATUS_DELETED  = "deleted"

	CYNOSDB_UPGRADE_IMMEDIATE = "upgradeImmediate"

	CYNOSDB_INSTANCE_RW_TYPE = "rw"
	CYNOSDB_INSTANCE_RO_TYPE = "ro"

	CYNOSDB_DEFAULT_OFFSET = 0
	CYNOSDB_MAX_LIMIT      = 100

	CYNOSDB_INSGRP_HA       = "ha"
	CYNOSDB_INSGRP_RO       = "ro"
	CYNOSDB_INSGRP_SINGLERO = "singleRo"

	// 0-成功，1-失败，2-处理中
	CYNOSDB_FLOW_STATUS_SUCCESSFUL = "0"
	CYNOSDB_FLOW_STATUS_FAILED     = "1"
	CYNOSDB_FLOW_STATUS_PROCESSING = "2"
)

const (
	STATUS_YES = "yes"
	STATUS_NO  = "no"

	RW_TYPE = "READWRITE"
	RO_TYPE = "READONLY"
)

var (
	CYNOSDB_PREPAID_PERIOD = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36}
	CYNOSDB_CHARGE_TYPE    = map[int64]string{
		0: svcpostgresql.COMMON_PAYTYPE_POSTPAID,
		1: svcpostgresql.COMMON_PAYTYPE_PREPAID,
	}
)

func TencentCynosdbInstanceBaseInfo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_cpu_core": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "CynosDB集群中读写类型实例的CPU核数。创建普通集群时需要。注意：该字段的修改将立即生效，如果要在维护时段升级，请从控制台升级。",
		},
		"instance_memory_size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "读写类型实例的内存容量，单位为GB。创建普通集群时需要。注意：该字段的修改将立即生效，如果要在维护时段升级，请从控制台升级。",
		},
		"instance_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "实例ID。",
		},
		"instance_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "实例名称。",
		},
		"instance_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "实例的状态。",
		},
		"instance_storage_size": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "实例存储大小，单位GB。",
		},
		"instance_maintain_weekdays": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			// DefaultFunc doesn't work but wil remain it
			DefaultFunc: func() (interface{}, error) {
				weekdays := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
				return weekdays, nil
			},
			Elem: &schema.Schema{Type: schema.TypeString},
			Set: func(v interface{}) int {
				return helper.HashString(v.(string))
			},
			Description: "平日进行维护。默认情况下`[\"星期一\"、\"星期二\"、\"星期三\"、\"星期四\"、\"星期五\"、\"星期六\"、\"星期日\"]`。",
		},
		"instance_maintain_start_time": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     10800,
			Description: "从 00:00 开始的偏移时间，单位为秒。例如，凌晨 03:00 应为“10800”。默认为“10800”。",
		},
		"instance_maintain_duration": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     3600,
			Description: "维护持续时间，单位为秒。默认为“3600”。",
		},
	}
}

func TencentCynosdbClusterBaseInfo() map[string]*schema.Schema {
	cluster := map[string]*schema.Schema{
		"project_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     0,
			Description: "项目 ID。默认为“0”。",
		},
		"available_zone": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "CynosDB集群的可用区。",
		},
		"vpc_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "专有网络ID。",
		},
		"subnet_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "该VPC内子网的ID。",
		},
		"old_ip_reserve_hours": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "旧地址回收时间，修改vpc时必须填写旧地址回收时间，修改vpc时必须填写。",
		},
		"port": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     5432,
			Description: "CynosDB集群的端口。",
		},
		"db_type": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "CynosDB 类型，可用值包括“MYSQL”。",
		},
		"db_version": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "CynosDB 的版本，与 db_type 相关。对于“MYSQL”，可用值为“5.7”、“8.0”。",
		},
		"storage_limit": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "CynosDB集群实例的存储限制，单位为GB。非无服务器实例的最大存储空间（以 GB 为单位）。注意：如果db_type为“MYSQL”，charge_type为“PREPAID”，则该值不能超过CPU和内存规格对应的最大存储，交易模式为“下单付款”。当 charge_type 为“POSTPAID_BY_HOUR”时，该参数是不必要的。",
		},
		"storage_pay_mode": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "集群存储计费方式，按量付费：`0`-包年/包月：`1`-默认按量付费。当DbType为MYSQL时，集群计算计费方式为后付费时（其中DbMode为SERVERLESS），存储计费方式只能按量计费；回滚和克隆不支持按年订阅按月存储。",
		},
		"cluster_name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "CynosDB 集群的名称。",
		},
		"password": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "“root”帐户的密码。",
		},
		"instance_count": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Computed:    true,
			Description: "实例数量，范围为（0,16]，默认值为2（即1个RW实例+1个Ro实例），传入的n表示1个RW实例+n-1个Ro实例（规格相同），如果需要更准确的集群组成，请使用InstanceInitInfos。",
		},
		// payment
		"charge_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      CYNOSDB_CHARGE_TYPE_POSTPAID,
			ValidateFunc: tccommon.ValidateAllowedStringValue([]string{CYNOSDB_CHARGE_TYPE_POSTPAID, CYNOSDB_CHARGE_TYPE_PREPAID}),
			Description: "实例的收费类型。有效值为“PREPAID”和“POSTPAID_BY_HOUR”。默认值为“POSTPAID_BY_HOUR”。",
		},
		"prepaid_period": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: tccommon.ValidateAllowedIntValue(CYNOSDB_PREPAID_PERIOD),
			Description: "预付费实例的租期（时间单位为月）。有效值为“1”、“2”、“3”、“4”、“5”、“6”、“7”、“8”、“9”、“10”、“11”、“12”、“24”、“36”。注意：仅当 charge_type 设置为“PREPAID”时才有效。",
		},
		"auto_renew_flag": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "自动更新标志。有效值为“0”(MANUAL_RENEW)、“1”(AUTO_RENEW)。默认值为“0”。仅适用于 PREPAID 集群。",
		},
		"force_delete": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "是否直接删除集群实例。默认为 false。如果设置为 true，集群及其“所有相关实例”将被删除，而不是保留在回收站中。注意：适用于“PREPAID”和“POSTPAID_BY_HOUR”集群。",
		},
		"tags": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "CynosDB集群的标签。",
		},
		// Computed
		"charset": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "CynosDB 集群使用的字符集。",
		},
		"cluster_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Cynosdb 集群的状态。",
		},
		"create_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "CynosDB 集群的创建时间。",
		},
		"storage_used": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "CynosDB集群已使用的存储，单位为MB。",
		},
		// rw instance group infos
		"rw_group_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "读写实例组ID。",
		},
		"rw_group_instances": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "读写实例组中的实例列表。",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"instance_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "实例ID。",
					},
					"instance_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "实例名称。",
					},
				},
			},
		},
		"rw_group_sg": {
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "`rw_group` 的安全组 ID。",
		},
		"rw_group_addr": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "读写地址。每个元素包含以下属性：",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ip": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "用于读写连接的 IP 地址。",
					},
					"port": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "读写连接的端口号。",
					},
				},
			},
		},
		// ro instance group infos
		"ro_group_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "只读实例组ID。",
		},
		"ro_group_instances": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "只读实例组中的实例列表。",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"instance_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "实例ID。",
					},
					"instance_name": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "实例名称。",
					},
				},
			},
		},
		"ro_group_sg": {
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "`ro_group` 的安全组 ID。",
		},
		"ro_group_addr": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "只读地址。每个元素包含以下属性：",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ip": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "只读连接的 IP 地址。",
					},
					"port": {
						Type:        schema.TypeInt,
						Computed:    true,
						Description: "只读连接的端口号。",
					},
				},
			},
		},
		"param_items": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "指定数据库的参数列表。创建集群时设置 param_template_id 时有效。使用 数据.tencentcloud_mysql_default_params 查询可用参数详细信息。",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "参数名称，例如`字符集服务器`。",
					},
					"old_value": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Param old 值，表示已经设置的值，修改current_value时需要该值。",
					},
					"current_value": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "要设置的参数预期值。",
					},
				},
			},
		},
		"prarm_template_id": {
			Type:          schema.TypeInt,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"param_template_id"},
			Deprecated:    "It will be deprecated. Use `param_template_id` instead.",
			Description: "参数模板的ID。",
		},
		"param_template_id": {
			Type:          schema.TypeInt,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"prarm_template_id"},
			Description: "参数模板的ID。",
		},
		"instance_init_infos": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			Description: "实例初始化配置信息，主要用于购买集群时选择不同规格的实例。",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"cpu": {
						Type:        schema.TypeInt,
						Required:    true,
						ForceNew:    true,
						Description: "实例的CPU。",
					},
					"memory": {
						Type:        schema.TypeInt,
						Required:    true,
						ForceNew:    true,
						Description: "实例记忆。",
					},
					"instance_type": {
						Type:        schema.TypeString,
						Required:    true,
						ForceNew:    true,
						Description: "实例类型。值：`rw`、`ro`。",
					},
					"instance_count": {
						Type:        schema.TypeInt,
						Required:    true,
						ForceNew:    true,
						Description: "实例计数。范围：[1，15]。",
					},
					"min_ro_count": {
						Type:        schema.TypeInt,
						Optional:    true,
						ForceNew:    true,
						Description: "无服务器实例的最小数量。范围[1,15]。",
					},
					"max_ro_count": {
						Type:        schema.TypeInt,
						Optional:    true,
						ForceNew:    true,
						Description: "无服务器实例的最大数量。范围[1,15]。",
					},
					"min_ro_cpu": {
						Type:        schema.TypeFloat,
						Optional:    true,
						ForceNew:    true,
						Description: "最低无服务器实例规格。",
					},
					"max_ro_cpu": {
						Type:        schema.TypeFloat,
						Optional:    true,
						ForceNew:    true,
						Description: "最大无服务器实例规格。",
					},
					"device_type": {
						Type:        schema.TypeString,
						Optional:    true,
						ForceNew:    true,
						Description: "实例机器类型。值：“通用”、“专有”。",
					},
				},
			},
		},
		"db_mode": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "指定DB模式，仅当`db_type`为`MYSQL`时可用。值：“正常”（默认）、“无服务器”。",
		},
		"min_cpu": {
			Optional:    true,
			Type:        schema.TypeFloat,
			Description: "当“db_mode”为“SERVERLESS”时需要最小 CPU 核心数，请请求DescribeServerlessInstanceSpecs 以获取更多参考。",
		},
		"max_cpu": {
			Optional:    true,
			Type:        schema.TypeFloat,
			Description: "最大 CPU 核心数，当 `db_mode` 为 `SERVERLESS` 时需要，请请求DescribeServerlessInstanceSpecs 以获取更多参考。",
		},
		"auto_pause": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "指定当“db_mode”为“SERVERLESS”时集群是否可以自动暂停。值：“是”（默认）、“否”。",
		},
		"auto_pause_delay": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "当“db_mode”为“SERVERLESS”时，指定自动暂停延迟（以秒为单位）。值范围：“[600，691200]”。默认值：`600`。",
		},
		"slave_zone": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "CynosDB 集群的多区域地址。",
		},
		"serverless_status_flag": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"resume", "pause"}),
			Description: "指定是否暂停或恢复无服务器集群。值：“继续”、“暂停”。",
		},
		"serverless_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "无服务器集群状态。注意：这是一个只读属性，要修改，请设置 `serverless_status_flag`。",
		},
		"cynos_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "内核次要版本，例如“3.1.16.002”。",
		},
	}

	for k, v := range TencentCynosdbInstanceBaseInfo() {
		cluster[k] = v
	}

	return cluster
}
