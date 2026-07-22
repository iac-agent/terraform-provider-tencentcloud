package cvm

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccbs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cbs"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudInstanceCreate,
		Read:   resourceTencentCloudInstanceRead,
		Update: resourceTencentCloudInstanceUpdate,
		Delete: resourceTencentCloudInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"image_id", "launch_template_id"},
				Description:  "镜像 到 使用 对于 实例. Modifications 可能 lead 到 reinstallation 的 实例's operating 系统。",
			},
			"availability_zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"availability_zone", "launch_template_id"},
				Description:  "可用 可用区 对于 CVM 实例。",
			},
			"dedicated_cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Exclusive 集群 ID",
			},
			"dedicated_resource_pack_tenancy": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Dedicated 资源 pack tenancy strategy. 有效值：`ResourcePool` (使用 实例 资源 池 对于 资源 pre-deduction)。",
			},
			"dedicated_resource_pack_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				RequiredWith: []string{"dedicated_resource_pack_tenancy"},
				Description:  "列表 dedicated 资源 pack IDs (e.g.，rpp-xxxxxxxx). 当 creating 实例 使用 pre-purchased 资源 池 packs，此 参数 必须 是 指定 together 使用 `dedicated_resource_pack_tenancy` 到 match corresponding tenancy strategy. Related 资源: `tencentcloud_cvm_resource_pool_packs`。",
			},
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(2, 128),
				Description:  "名称 实例. max 长度 的 instance_name 是 128，和 默认值为 `Terraform-CVM-实例`。",
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateInstanceType,
				Description:  "类型 实例。",
			},
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "hostname 的 实例. Windows 实例: 名称 should 是 combination 的 2 到 15 字符 comprised 的 letters (case insensitive)，numbers，和 hyphens (-). 周期 (.) 是 不 支持，和 名称 不能 是 字符串 的 pure numbers. Other types (such 作为 Linux) 的 实例: 名称 should 是 combination 的 2 到 60 字符，supporting 多个 periods (.). piece between two periods 是 composed 的 letters (case insensitive)，numbers，和 hyphens (-). Changing `hostname` 将 cause 实例 系统 到 restart。",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "项目 实例 belongs 到，默认为 0。",
			},
			"running_flag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Set 实例 到 running 或 stop. 默认值为 true， 实例 将 shutdown 当 此 flag 是 false。",
			},
			"stop_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "实例 shutdown 模式 有效值：SOFT_FIRST: perform soft shutdown first，和 force shut down 实例 如果 soft shutdown fails; HARD: force shut down 实例 directly; SOFT: soft shutdown 仅. 默认值：SOFT。",
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{
					CVM_STOP_TYPE_SOFT_FIRST,
					CVM_STOP_TYPE_HARD,
					CVM_STOP_TYPE_SOFT,
				}),
			},
			"stopped_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Billing 方法 的 pay-作为-您-go 实例 after shutdown. 可用值：`KEEP_CHARGING`,`STOP_CHARGING`. Default `KEEP_CHARGING`。",
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{
					CVM_STOP_MODE_KEEP_CHARGING,
					CVM_STOP_MODE_STOP_CHARGING,
				}),
			},
			"disaster_recover_group_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"placement_group_id"},
				Description:   "Placement 组 ID",
			},
			"placement_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID placement 组。",
			},
			"force_replace_placement_group_id": {
				Type:         schema.TypeBool,
				Optional:     true,
				RequiredWith: []string{"placement_group_id"},
				Description:  "是否force 实例 主机 到 是 replaced. 取值范围：true: Allows 实例 到 change 主机 和 restart 实例. Local 磁盘 machines do 不 support specifying 此 参数; false: Does 不 allow 实例 到 change 主机 和 仅 join placement 组 在 当前 主机 此 可能 cause placement 组 到 fail 到 change. Only useful 对于 change `placement_group_id`，默认为 false。",
			},
			// payment
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CVM_CHARGE_TYPE),
				Description:  "charge 类型 实例. 有效 值 是 `PREPAID`，`POSTPAID_BY_HOUR`，`SPOTPAID`，`CDHPAID` 和 `CDCPAID`. 默认为 `POSTPAID_BY_HOUR`. 注意: TencentCloud International 仅 支持 `POSTPAID_BY_HOUR` 和 `CDHPAID`. `PREPAID` 实例 可能 不 allow 到 delete before expired. `SPOTPAID` 实例 必须 集合 `spot_instance_type` 和 `spot_max_price` 在 same 时间. `CDHPAID` 实例 必须 集合 `cdh_instance_type` 和 `cdh_host_id`。",
			},
			"instance_charge_type_prepaid_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue(CVM_PREPAID_PERIOD),
				Description:  "tenancy (时间 单位 是 month) 的 prepaid 实例，NOTE: 它 仅 works 当 instance_charge_type 是 集合 到 `PREPAID`. 有效 值 是 `1`，`2`，`3`，`4`，`5`，`6`，`7`，`8`，`9`，`10`，`11`，`12`，`24`，`36`，`48`，`60`。",
			},
			"instance_charge_type_prepaid_renew_flag": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CVM_PREPAID_RENEW_FLAG),
				Description:  "自动续费标识 有效值：`NOTIFY_AND_AUTO_RENEW`: notify upon expiration 和 renew automatically，`NOTIFY_AND_MANUAL_RENEW`: notify upon expiration 但 do 不 renew automatically，`DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. 默认值：`NOTIFY_AND_MANUAL_RENEW`. 如果 此 参数 是 指定 作为 `NOTIFY_AND_AUTO_RENEW`， 实例 将 是 automatically renewed 在 monthly basis 如果 账号 balance 是 sufficient. NOTE: 它 仅 works 当 instance_charge_type 是 集合 到 `PREPAID`。",
			},
			"spot_instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CVM_SPOT_INSTANCE_TYPE),
				Description:  "类型 spot 实例，仅 support `ONE-TIME` now. 注意: 它 仅 works 当 instance_charge_type 是 集合 到 `SPOTPAID`。",
			},
			"spot_max_price": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateStringNumber,
				Description: "竞价型实例的最高价格，是十进制字符串的格式，例如“0.50”。注意：仅当instance_charge_type设置为“SPOTPAID”时才有效。",
			},
			"cdh_instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateStringPrefix("CDH_"),
				Description:  "类型 实例 创建 在 cdh， 值 的 此 参数 是 在 格式 的 CDH_XCXG based 在 数量 CPU 核数 和 内存 容量. 注意: 它 仅 works 当 instance_charge_type 是 集合 到 `CDHPAID`。",
			},
			"cdh_host_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "ID 的 cdh 实例. 注意: 它 仅 works 当 instance_charge_type 是 集合 到 `CDHPAID`。",
			},
			// network
			"internet_charge_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					stopMode := d.Get("stopped_mode").(string)
					if stopMode != CVM_STOP_MODE_STOP_CHARGING || !d.HasChange("running_flag") {
						return old == new
					}
					return old == "" || new == ""
				},
				ValidateFunc: tccommon.ValidateAllowedStringValue(CVM_INTERNET_CHARGE_TYPE),
				Description:  "Internet charge 类型 实例，有效 值 是 `BANDWIDTH_PREPAID`，`TRAFFIC_POSTPAID_BY_HOUR`，`BANDWIDTH_POSTPAID_BY_HOUR` 和 `BANDWIDTH_PACKAGE`. 如果 不 集合，公网计费类型 是 consistent 使用 cvm 计费类型 通过 默认值. 此 值 takes NO Effect 当 changing 和 does 不 need 到 是 集合 当 `allocate_public_ip` 是 false。",
			},
			"bandwidth_package_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "带宽 包 ID. 如果 用户 是 standard 用户，then bandwidth_package_id 是 needed，或 默认值 has bandwidth_package_id。",
			},
			"internet_max_bandwidth_out": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Maximum outgoing 带宽 到 公有 网络，measured 在 Mbps (Mega bits per second). 此 值 does 不 need 到 是 集合 当 `allocate_public_ip` 是 false。",
			},
			"allocate_public_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Associate 公网 IP 地址 使用 实例 在 VPC 或 Classic. Boolean 值，默认为 false。",
			},
			"ipv4_address_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"WanIP", "HighQualityEIP", "AntiDDoSEIP"}),
				Description:  "AddressType. 默认值：WanIP. For beta users 的 dedicated IP. 值 可以 是: HighQualityEIP: Dedicated IP. 注意 该 dedicated IPs 是 仅 可用 在 partial regions. For beta users 的 Anti-DDoS IP， 值 可以 是: AntiDDoSEIP: Anti-DDoS EIP. 注意 该 Anti-DDoS IPs 是 仅 可用 在 partial regions。",
			},
			"ipv6_address_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"EIPv6", "HighQualityEIPv6"}),
				Description:  "IPv6 AddressType. 默认值：WanIP. EIPv6: Elastic IPv6; HighQualityEIPv6: Premium IPv6，仅 China Hong Kong 支持 premium IPv6. To allocate IPv6 addresses 到 resources，please 指定Elastic IPv6 类型",
			},
			"ipv6_address_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "指定number 的 randomly generated IPv6 addresses 对于 Elastic Network Interface。",
			},
			"anti_ddos_package_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Anti-DDoS 服务 包 ID. 此 为必填项 当 您 want 到 请求 AntiDDoS IP。",
			},
			// vpc
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID VPC 网络. 如果 您 want 到 create 实例 在 VPC 网络，此 参数 必须 是 集合。",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID VPC 子网. 如果 您 want 到 create 实例 在 VPC 网络，此 参数 必须 是 集合。",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "私有 IP 到 是 assigned 到 此 实例，必须 是 在 提供 子网 和 可用。",
			},
			// security group
			"security_groups": {
				Type:          schema.TypeSet,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"orderly_security_groups"},
				Description:   "A 列表 安全 组 IDs 到 associate 使用。",
				Deprecated:    "It will be deprecated. Use `orderly_security_groups` instead.",
			},

			"orderly_security_groups": {
				Type:          schema.TypeList,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"security_groups"},
				Description:   "A 列表 orderly 安全 组 IDs 到 associate 使用。",
			},
			// storage
			"system_disk_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CVM_DISK_TYPE),
				Description:  "System 磁盘 类型 For more 信息 在 limits 的 系统 磁盘 types，see [Storage Overview](https://intl.云.tencent.com/document/product/213/4952). 有效值：`LOCAL_BASIC`: 本地 磁盘，`LOCAL_SSD`: 本地 SSD 磁盘，`CLOUD_BASIC`: 云 磁盘，`CLOUD_SSD`: 云 SSD 磁盘，`CLOUD_PREMIUM`: Premium Cloud Storage，`CLOUD_BSSD`: Basic SSD，`CLOUD_HSSD`: Enhanced SSD，`CLOUD_TSSD`: Tremendous SSD. NOTE: 如果 modified， 实例 可能 强制停止",
			},
			"system_disk_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Size 的 系统 磁盘. 单位 是 GB，默认为 50GB. 如果 modified， 实例 可能 强制停止",
			},
			"system_disk_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "System 磁盘 快照 ID 用于initialize 系统 磁盘. 当 系统 磁盘 类型 是 `LOCAL_BASIC` 和 `LOCAL_SSD`，磁盘 ID 是 不 支持。",
			},
			"system_disk_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "名称 系统 磁盘。",
			},
			"system_disk_resize_online": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Resize online。",
			},
			"data_disks": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Settings 对于 数据 disks。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_disk_type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Data 磁盘 类型 For more 信息 about limits 在 different 数据 磁盘 types，see [Storage Overview](https://intl.云.tencent.com/document/product/213/4952). 有效值：LOCAL_BASIC: 本地 磁盘，LOCAL_SSD: 本地 SSD 磁盘，LOCAL_NVME: 本地 NVME 磁盘，指定 在 InstanceType，LOCAL_PRO: 本地 HDD 磁盘，指定 在 InstanceType，CLOUD_BASIC: HDD 云 磁盘，CLOUD_PREMIUM: Premium Cloud Storage，CLOUD_SSD: SSD，CLOUD_HSSD: Enhanced SSD，CLOUD_TSSD: Tremendous SSD，CLOUD_BSSD: Balanced SSD。",
						},
						"data_disk_size": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Size 的 数据 磁盘，和 单位 是 GB。",
						},
						"data_disk_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "名称 数据 磁盘。",
						},
						"data_disk_snapshot_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Snapshot ID 数据 磁盘. selected 数据 磁盘 快照 大小 必须 是 smaller 比 数据 磁盘 大小。",
						},
						"data_disk_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Data 磁盘 ID 用于initialize 数据 磁盘. 当 数据 磁盘 类型 是 `LOCAL_BASIC` 和 `LOCAL_SSD`，磁盘 ID 是 不 支持。",
						},
						"delete_with_instance": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							ForceNew:    true,
							Description: "Decides 是否disk 是 删除 使用 实例(仅 applied 到 `CLOUD_BASIC`，`CLOUD_SSD` 和 `CLOUD_PREMIUM` 磁盘 使用 `POSTPAID_BY_HOUR` 实例)，默认为 true。",
						},
						"delete_with_instance_prepaid": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							ForceNew:    true,
							Description: "Decides 是否disk 是 删除 使用 实例(仅 applied 到 `CLOUD_BASIC`，`CLOUD_SSD` 和 `CLOUD_PREMIUM` 磁盘 使用 `PREPAID` 实例)，默认为 false。",
						},
						"kms_key_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "可选 参数. 当 purchasing 加密 磁盘，customize 键 当 此 参数 是 passed 在， `encrypt` 参数 need 是 集合。",
						},
						"encrypt": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							ForceNew:    true,
							Description: "Decides 是否disk 是 encrypted. 默认为 `false`。",
						},
						"throughput_performance": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							ForceNew:    true,
							Description: "Add extra performance 到 数据 磁盘. Only works 当 磁盘 类型 是 `CLOUD_TSSD` 或 `CLOUD_HSSD`。",
						},
					},
				},
			},
			// enhance services
			"disable_security_service": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disable enhance 服务 对于 安全，它 是 已启用 通过 默认值. 当 此 options 是 集合，安全 agent won't 是 installed. Modifications 可能 lead 到 reinstallation 的 实例's operating 系统。",
			},
			"disable_monitor_service": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disable enhance 服务 对于 监控，它 是 已启用 通过 默认值. 当 此 options 是 集合，监控 agent won't 是 installed. Modifications 可能 lead 到 reinstallation 的 实例's operating 系统。",
			},
			"disable_automation_service": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disable enhance 服务 对于 automation，它 是 已启用 通过 默认值. 当 此 options 是 集合，监控 agent won't 是 installed. Modifications 可能 lead 到 reinstallation 的 实例's operating 系统。",
			},
			// login
			"key_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Please use `key_ids` instead.",
				ConflictsWith: []string{"key_ids"},
				Description:   "键 pair 到 使用 对于 实例，它 looks like `skey-16jig7tx`. Modifications 可能 lead 到 reinstallation 的 实例's operating 系统。",
			},
			"key_ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"key_name", "password"},
				Description:   "键 pair 到 使用 对于 实例，它 looks like `skey-16jig7tx`. Modifications 可能 lead 到 reinstallation 的 实例's operating 系统。",
				Set:           schema.HashString,
				Elem:          &schema.Schema{Type: schema.TypeString},
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "密码 对于 实例. In 顺序 对于 new 密码 到 take effect， 实例 将 是 restarted after 密码 change. Modifications 可能 lead 到 reinstallation 的 实例's operating 系统。",
			},
			"keep_image_login": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "false" && old == "" || old == "false" && new == "" {
						return true
					} else {
						return old == new
					}
				},
				ConflictsWith: []string{"key_name", "key_ids", "password"},
				Description:   "是否keep 镜像 login 或 不，默认为 `false`. 当 镜像 类型 是 私有 或 shared 或 imported，此 参数 可以 是 集合 `true`. Modifications 可能 lead 到 reinstallation 的 实例's operating 系统.。",
			},
			"user_data": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"user_data_raw"},
				Description:   "用户 数据 到 是 injected into 此 实例. Must 是 base64 encoded 和 up 到 16 KB. 如果 `user_data_replace_on_change` 是 集合 到 `true`，updates 到 此 字段 将 触发器 destruction 和 recreation 的 CVM 实例。",
			},
			"user_data_raw": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"user_data"},
				Description:   "用户 数据 到 是 injected into 此 实例，在 plain text. Conflicts 使用 `user_data`. Up 到 16 KB after base64 encoded. 如果 `user_data_replace_on_change` 是 集合 到 `true`，updates 到 此 字段 将 触发器 destruction 和 recreation 的 CVM 实例。",
			},
			"user_data_replace_on_change": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "当 使用 在 combination 使用 `user_data` 或 `user_data_raw` 将 触发器 destroy 和 recreate 的 CVM 实例 当 集合 到 `true`. 默认为 `false`。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "A mapping 的 标签 到 assign 到 资源. For 标签 limits，please refer 到 [Use Limits](https://intl.云.tencent.com/document/product/651/13354)。",
			},
			"force_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate 是否force delete 实例. 默认为 `false`. 如果 集合 true， 实例 将 是 permanently 删除 instead 的 being moved into recycle bin. 注意: 仅 works 对于 `PREPAID` 实例。",
			},
			"disable_api_termination": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "是否termination protection 是 已启用 默认为 `false`. 如果 集合 true，其中 表示 该 此 实例 可以 不 是 删除 通过 API 操作",
			},
			// role
			"cam_role_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "被授权访问的 CAM 角色名称",
			},
			"hpc_cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "High-performance computing 集群 ID. 如果 实例 创建 是 high-performance computing 实例，您 need 到 指定cluster 在 其中 实例 是 placed，otherwise 它 不能 是 指定。",
			},
			// template
			"launch_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "实例 launch 模板 ID 此 参数 allows 您 到 create 实例 使用 preset 参数 在 实例 template。",
			},
			"launch_template_version": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "实例 launch template 版本 数量. 如果 given， new 实例 launch template 将 是 创建 based 在 given 版本 数量。",
			},
			"release_address": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Release elastic IP. Under EIP 2.0，仅 first EIP under primary 网络 card 是 提供，和 EIP types 是 limited 到 HighQualityEIP，AntiDDoSEIP，EIPv6，和 HighQualityEIPv6. Default behavior 是 不 released。",
			},
			// Computed values.
			"instance_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current 状态 实例。",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP 的 实例。",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Globally 唯一 ID 实例。",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "创建时间 的 实例。",
			},
			"expired_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "过期时间 的 实例。",
			},
			"cpu": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "数量 CPU 核数 的 实例。",
			},
			"memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "实例 内存 容量，单位 （GB）。",
			},
			"os_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "实例 os 名称",
			},
			"rack_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "rack ID 实例 资源 池 到 其中 实例 belongs。",
			},
			"ipv6_addresses": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "IPv6 地址 的 实例。",
			},
			"public_ipv6_addresses": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "公有 IPv6 地址 到 其中 实例 是 bound。",
			},
		},

		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIf("user_data", func(_ context.Context, diff *schema.ResourceDiff, meta interface{}) bool {
				return diff.Get("user_data_replace_on_change").(bool)
			}),

			customdiff.ForceNewIf("user_data_raw", func(_ context.Context, diff *schema.ResourceDiff, meta interface{}) bool {
				return diff.Get("user_data_replace_on_change").(bool)
			}),
		),
	}
}

func resourceTencentCloudInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_instance.create")()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		ctx                = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		cvmService         = CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceChargeType = CVM_CHARGE_TYPE_POSTPAID
	)

	request := cvm.NewRunInstancesRequest()
	if v, ok := d.GetOk("image_id"); ok {
		request.ImageId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("availability_zone"); ok {
		request.Placement = &cvm.Placement{
			Zone: helper.String(v.(string)),
		}
	}

	if v, ok := d.GetOk("dedicated_cluster_id"); ok {
		request.DedicatedClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		projectId := int64(v.(int))
		request.Placement.ProjectId = &projectId
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("hostname"); ok {
		request.HostName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cam_role_name"); ok {
		request.CamRoleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("hpc_cluster_id"); ok {
		request.HpcClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		instanceChargeType = v.(string)
		request.InstanceChargeType = &instanceChargeType
		if instanceChargeType == CVM_CHARGE_TYPE_PREPAID || instanceChargeType == CVM_CHARGE_TYPE_UNDERWRITE {
			request.InstanceChargePrepaid = &cvm.InstanceChargePrepaid{}
			if period, ok := d.GetOk("instance_charge_type_prepaid_period"); ok {
				periodInt64 := int64(period.(int))
				request.InstanceChargePrepaid.Period = &periodInt64
			}

			if renewFlag, ok := d.GetOk("instance_charge_type_prepaid_renew_flag"); ok {
				request.InstanceChargePrepaid.RenewFlag = helper.String(renewFlag.(string))
			}
		}

		if instanceChargeType == CVM_CHARGE_TYPE_SPOTPAID {
			spotInstanceType, sitOk := d.GetOk("spot_instance_type")
			spotMaxPrice, smpOk := d.GetOk("spot_max_price")
			if sitOk || smpOk {
				request.InstanceMarketOptions = &cvm.InstanceMarketOptionsRequest{}
				request.InstanceMarketOptions.MarketType = helper.String(CVM_MARKET_TYPE_SPOT)
				request.InstanceMarketOptions.SpotOptions = &cvm.SpotMarketOptions{}
			}

			if sitOk {
				request.InstanceMarketOptions.SpotOptions.SpotInstanceType = helper.String(strings.ToLower(spotInstanceType.(string)))
			}

			if smpOk {
				request.InstanceMarketOptions.SpotOptions.MaxPrice = helper.String(spotMaxPrice.(string))
			}
		}

		if instanceChargeType == CVM_CHARGE_TYPE_CDHPAID {
			if v, ok := d.GetOk("cdh_instance_type"); ok {
				request.InstanceType = helper.String(v.(string))
			} else {
				return fmt.Errorf("cdh_instance_type can not be empty when instance_charge_type is %s", instanceChargeType)
			}

			if v, ok := d.GetOk("cdh_host_id"); ok {
				request.Placement.HostIds = append(request.Placement.HostIds, helper.String(v.(string)))
			} else {
				return fmt.Errorf("cdh_host_id can not be empty when instance_charge_type is %s", instanceChargeType)
			}
		}
	}

	// Dedicated resource pack placement parameters
	if v, ok := d.GetOk("dedicated_resource_pack_tenancy"); ok {
		request.Placement.DedicatedResourcePackTenancy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dedicated_resource_pack_ids"); ok {
		packIds := v.(*schema.Set).List()
		for _, packId := range packIds {
			request.Placement.DedicatedResourcePackIds = append(
				request.Placement.DedicatedResourcePackIds,
				helper.String(packId.(string)),
			)
		}
	}

	// Check for disaster_recover_group_ids first (new field)
	if v, ok := d.GetOk("disaster_recover_group_ids"); ok {
		disasterRecoverGroupIdsSet := v.(*schema.Set).List()
		for i := range disasterRecoverGroupIdsSet {
			disasterRecoverGroupId := disasterRecoverGroupIdsSet[i].(string)
			request.DisasterRecoverGroupIds = append(request.DisasterRecoverGroupIds, &disasterRecoverGroupId)
		}
	}

	var rpgFlag bool
	if v, ok := d.GetOkExists("force_replace_placement_group_id"); ok {
		rpgFlag = v.(bool)
	}

	if !rpgFlag {
		if v, ok := d.GetOk("placement_group_id"); ok {
			request.DisasterRecoverGroupIds = []*string{helper.String(v.(string))}
		}
	}

	// network
	var (
		internetAccessible cvm.InternetAccessible
		netWorkFlag        bool
	)

	if v, ok := d.GetOk("internet_charge_type"); ok {
		internetAccessible.InternetChargeType = helper.String(v.(string))
		netWorkFlag = true
	}

	if v, ok := d.GetOkExists("internet_max_bandwidth_out"); ok {
		maxBandwidthOut := int64(v.(int))
		internetAccessible.InternetMaxBandwidthOut = &maxBandwidthOut
		netWorkFlag = true
	}

	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		internetAccessible.BandwidthPackageId = helper.String(v.(string))
		netWorkFlag = true
	}

	if v, ok := d.GetOkExists("allocate_public_ip"); ok {
		allocatePublicIp := v.(bool)
		internetAccessible.PublicIpAssigned = &allocatePublicIp
		netWorkFlag = true
	}

	if v, ok := d.GetOk("ipv4_address_type"); ok {
		internetAccessible.IPv4AddressType = helper.String(v.(string))
		netWorkFlag = true
	}

	if v, ok := d.GetOk("ipv6_address_type"); ok {
		internetAccessible.IPv6AddressType = helper.String(v.(string))
		netWorkFlag = true
	}

	if v, ok := d.GetOk("anti_ddos_package_id"); ok {
		internetAccessible.AntiDDoSPackageId = helper.String(v.(string))
		netWorkFlag = true
	}

	if netWorkFlag {
		request.InternetAccessible = &internetAccessible
	}

	// vpc
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{}
		request.VirtualPrivateCloud.VpcId = helper.String(v.(string))

		if v, ok = d.GetOk("subnet_id"); ok {
			request.VirtualPrivateCloud.SubnetId = helper.String(v.(string))
		}

		if v, ok = d.GetOk("private_ip"); ok {
			request.VirtualPrivateCloud.PrivateIpAddresses = []*string{helper.String(v.(string))}
		}
		if v, ok = d.GetOkExists("ipv6_address_count"); ok {
			request.VirtualPrivateCloud.Ipv6AddressCount = helper.IntUint64(v.(int))
		}
	}

	if v, ok := d.GetOk("security_groups"); ok {
		securityGroups := v.(*schema.Set).List()
		request.SecurityGroupIds = make([]*string, 0, len(securityGroups))
		for _, securityGroup := range securityGroups {
			request.SecurityGroupIds = append(request.SecurityGroupIds, helper.String(securityGroup.(string)))
		}
	}

	if v, ok := d.GetOk("orderly_security_groups"); ok {
		securityGroups := v.([]interface{})
		request.SecurityGroupIds = make([]*string, 0, len(securityGroups))
		for _, securityGroup := range securityGroups {
			request.SecurityGroupIds = append(request.SecurityGroupIds, helper.String(securityGroup.(string)))
		}
	}

	// storage
	var (
		systemDisk     cvm.SystemDisk
		systemDiskFlag bool
	)

	if v, ok := d.GetOk("system_disk_type"); ok {
		systemDisk.DiskType = helper.String(v.(string))
		systemDiskFlag = true
	}

	if v, ok := d.GetOkExists("system_disk_size"); ok {
		diskSize := int64(v.(int))
		systemDisk.DiskSize = &diskSize
		systemDiskFlag = true
	}

	if v, ok := d.GetOk("system_disk_id"); ok {
		systemDisk.DiskId = helper.String(v.(string))
		systemDiskFlag = true
	}

	if v, ok := d.GetOk("system_disk_name"); ok {
		systemDisk.DiskName = helper.String(v.(string))
		systemDiskFlag = true
	}

	if systemDiskFlag {
		request.SystemDisk = &systemDisk
	}

	if v, ok := d.GetOk("data_disks"); ok {
		dataDisks := v.([]interface{})
		for _, d := range dataDisks {
			value := d.(map[string]interface{})
			diskType := value["data_disk_type"].(string)
			diskSize := int64(value["data_disk_size"].(int))
			throughputPerformance := int64(value["throughput_performance"].(int))
			dataDisk := cvm.DataDisk{
				DiskType:              &diskType,
				DiskSize:              &diskSize,
				ThroughputPerformance: &throughputPerformance,
			}

			if v, ok := value["data_disk_name"]; ok && v != nil {
				diskName := v.(string)
				if diskName != "" {
					dataDisk.DiskName = helper.String(diskName)
				}
			}

			if v, ok := value["data_disk_snapshot_id"]; ok && v != nil {
				snapshotId := v.(string)
				if snapshotId != "" {
					dataDisk.SnapshotId = helper.String(snapshotId)
				}
			}

			if value["data_disk_id"] != "" {
				dataDisk.DiskId = helper.String(value["data_disk_id"].(string))
			}

			if deleteWithInstance, ok := value["delete_with_instance"]; ok {
				deleteWithInstanceBool := deleteWithInstance.(bool)
				if (instanceChargeType != CVM_CHARGE_TYPE_POSTPAID) && deleteWithInstanceBool {
					return fmt.Errorf("param `delete_with_instance` only can be true when `instance_charge_type` is %s", CVM_CHARGE_TYPE_POSTPAID)
				}

				dataDisk.DeleteWithInstance = &deleteWithInstanceBool
			}

			if v, ok := value["kms_key_id"]; ok && v != "" {
				dataDisk.KmsKeyId = helper.String(v.(string))
			}

			if encrypt, ok := value["encrypt"]; ok {
				encryptBool := encrypt.(bool)
				dataDisk.Encrypt = &encryptBool
			}

			request.DataDisks = append(request.DataDisks, &dataDisk)
		}
	}

	// enhanced service
	var (
		enhancedService     cvm.EnhancedService
		enhancedServiceFlag bool
	)

	if v, ok := d.GetOkExists("disable_security_service"); ok {
		securityService := !(v.(bool))
		enhancedService.SecurityService = &cvm.RunSecurityServiceEnabled{
			Enabled: &securityService,
		}
		enhancedServiceFlag = true
	}

	if v, ok := d.GetOkExists("disable_monitor_service"); ok {
		monitorService := !(v.(bool))
		enhancedService.MonitorService = &cvm.RunMonitorServiceEnabled{
			Enabled: &monitorService,
		}
		enhancedServiceFlag = true
	}

	if v, ok := d.GetOkExists("disable_automation_service"); ok {
		automationService := !(v.(bool))
		enhancedService.AutomationService = &cvm.RunAutomationServiceEnabled{
			Enabled: &automationService,
		}
		enhancedServiceFlag = true
	}

	if enhancedServiceFlag {
		request.EnhancedService = &enhancedService
	}

	// login
	var (
		loginSettings     cvm.LoginSettings
		loginSettingsFlag bool
	)

	if v, ok := d.GetOk("key_name"); ok {
		loginSettings.KeyIds = []*string{helper.String(v.(string))}
		loginSettingsFlag = true
	}

	if v, ok := d.GetOk("key_ids"); ok {
		keyIds := v.(*schema.Set).List()
		if len(keyIds) > 0 {
			loginSettings.KeyIds = helper.InterfacesStringsPoint(keyIds)
			loginSettingsFlag = true
		}
	}

	if v, ok := d.GetOk("password"); ok {
		loginSettings.Password = helper.String(v.(string))
		loginSettingsFlag = true
	}

	if v, ok := d.GetOkExists("keep_image_login"); ok {
		if v.(bool) {
			loginSettings.KeepImageLogin = helper.String(CVM_IMAGE_LOGIN)
		} else {
			loginSettings.KeepImageLogin = helper.String(CVM_IMAGE_LOGIN_NOT)
		}

		loginSettingsFlag = true
	}

	if loginSettingsFlag {
		request.LoginSettings = &loginSettings
	}

	if v, ok := d.GetOk("user_data"); ok {
		request.UserData = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_data_raw"); ok {
		userData := base64.StdEncoding.EncodeToString([]byte(v.(string)))
		request.UserData = &userData
	}

	if v, ok := d.GetOkExists("disable_api_termination"); ok {
		request.DisableApiTermination = helper.Bool(v.(bool))
	}

	var launchTemplate cvm.LaunchTemplate
	if v, ok := d.GetOk("launch_template_id"); ok {
		launchTemplate.LaunchTemplateId = helper.String(v.(string))
		request.LaunchTemplate = &launchTemplate
	}

	if v, ok := d.GetOkExists("launch_template_version"); ok {
		launchTemplate.LaunchTemplateVersion = helper.IntUint64(v.(int))
		request.LaunchTemplate = &launchTemplate
	}

	if v := helper.GetTags(d, "tags"); len(v) > 0 {
		tags := make([]*cvm.Tag, 0)
		for tagKey, tagValue := range v {
			tag := cvm.Tag{
				Key:   helper.String(tagKey),
				Value: helper.String(tagValue),
			}

			tags = append(tags, &tag)
		}

		tagSpecification := cvm.TagSpecification{
			ResourceType: helper.String("instance"),
			Tags:         tags,
		}

		request.TagSpecification = append(request.TagSpecification, &tagSpecification)
	}

	clientToken := helper.BuildToken()
	request.ClientToken = &clientToken

	instanceId := ""
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check("create")
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().RunInstances(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			e, ok := err.(*sdkErrors.TencentCloudSDKError)
			if ok && tccommon.IsContains(CVM_RETRYABLE_ERROR, e.Code) {
				return resource.RetryableError(fmt.Errorf("cvm create error: %s, retrying", e.Error()))
			}

			return tccommon.RetryError(err)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		if len(response.Response.InstanceIdSet) < 1 {
			err = fmt.Errorf("instance id is nil")
			return resource.NonRetryableError(err)
		}

		instanceId = *response.Response.InstanceIdSet[0]
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(instanceId)

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}

		if instance != nil && *instance.InstanceState == CVM_STATUS_LAUNCH_FAILED {
			//LatestOperationCodeMode
			if instance.LatestOperationErrorMsg != nil {
				return resource.NonRetryableError(fmt.Errorf("cvm instance %s launch failed. Error msg: %s.\n", *instance.InstanceId, *instance.LatestOperationErrorMsg))
			}

			return resource.NonRetryableError(fmt.Errorf("cvm instance %s launch failed, this resource will not be stored to tfstate and will auto removed\n.", *instance.InstanceId))
		}

		if instance != nil && *instance.InstanceState == CVM_STATUS_RUNNING {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
	})

	if err != nil {
		return err
	}

	// set placement group id
	if rpgFlag {
		if v, ok := d.GetOk("placement_group_id"); ok && v != "" {
			request := cvm.NewModifyInstancesDisasterRecoverGroupRequest()
			request.InstanceIds = helper.Strings([]string{instanceId})
			request.DisasterRecoverGroupId = helper.String(v.(string))
			request.Force = helper.Bool(rpgFlag)
			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ModifyInstancesDisasterRecoverGroup(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				return err
			}

			// wait
			err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
				if errRet != nil {
					return tccommon.RetryError(errRet, tccommon.InternalError)
				}

				if instance != nil && *instance.InstanceState == CVM_STATUS_LAUNCH_FAILED {
					//LatestOperationCodeMode
					if instance.LatestOperationErrorMsg != nil {
						return resource.NonRetryableError(fmt.Errorf("cvm instance %s launch failed. Error msg: %s.\n", *instance.InstanceId, *instance.LatestOperationErrorMsg))
					}

					return resource.NonRetryableError(fmt.Errorf("cvm instance %s launch failed, this resource will not be stored to tfstate and will auto removed\n.", *instance.InstanceId))
				}

				if instance != nil && *instance.InstanceState == CVM_STATUS_RUNNING {
					return nil
				}

				return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
			})

			if err != nil {
				return err
			}
		}
	}

	// Wait for the tags attached to the vm since tags attachment it's async while vm creation.
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			actualTags, e := tagService.DescribeResourceTags(ctx, "cvm", "instance", tcClient.Region, instanceId)
			if e != nil {
				return resource.RetryableError(e)
			}

			for tagKey, tagValue := range tags {
				if v, ok := actualTags[tagKey]; !ok || v != tagValue {
					return resource.RetryableError(fmt.Errorf("tag(%s, %s) modification is not completed", tagKey, tagValue))
				}
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	if v, ok := d.GetOkExists("running_flag"); ok {
		if !v.(bool) {
			stopType := d.Get("stop_type").(string)
			stoppedMode := d.Get("stopped_mode").(string)
			err = cvmService.StopInstance(ctx, instanceId, stopType, stoppedMode)
			if err != nil {
				return err
			}

			err = resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
				instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
				if errRet != nil {
					return tccommon.RetryError(errRet, tccommon.InternalError)
				}

				if instance != nil && *instance.InstanceState == CVM_STATUS_STOPPED {
					return nil
				}

				return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
			})

			if err != nil {
				return err
			}
		}
	}

	return resourceTencentCloudInstanceRead(d, meta)
}

func resourceTencentCloudInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		client     = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		cvmService = CvmService{client: client}
		cbsService = svccbs.NewCbsService(client)
		instanceId = d.Id()
	)

	forceDelete := false
	if v, ok := d.GetOkExists("force_delete"); ok {
		forceDelete = v.(bool)
		_ = d.Set("force_delete", forceDelete)
	}

	var instance *cvm.Instance
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet = cvmService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}

		if instance != nil && instance.LatestOperationState != nil && *instance.LatestOperationState == "OPERATING" {
			return resource.RetryableError(fmt.Errorf("waiting for instance %s operation", *instance.InstanceId))
		}

		return nil
	})

	if err != nil {
		return err
	}

	if instance == nil || *instance.InstanceState == CVM_STATUS_LAUNCH_FAILED {
		d.SetId("")
		log.Printf("[CRITAL]instance %s not exist or launch failed", instanceId)
		return nil
	}

	var cvmImages []string
	var response *cvm.DescribeImagesResponse
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		request := cvm.NewDescribeImagesRequest()
		response, errRet = client.UseCvmClient().DescribeImages(request)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}

		if *response.Response.TotalCount > 0 {
			for i := range response.Response.ImageSet {
				image := response.Response.ImageSet[i]
				cvmImages = append(cvmImages, *image.ImageId)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	if d.Get("image_id").(string) == "" || instance.ImageId == nil || !tccommon.IsContains(cvmImages, *instance.ImageId) {
		_ = d.Set("image_id", instance.ImageId)
	}

	_ = d.Set("availability_zone", instance.Placement.Zone)
	_ = d.Set("dedicated_cluster_id", instance.DedicatedClusterId)
	_ = d.Set("instance_name", instance.InstanceName)
	_ = d.Set("instance_type", instance.InstanceType)
	_ = d.Set("project_id", instance.Placement.ProjectId)

	// Set dedicated resource pack placement parameters
	if instance.Placement != nil {
		if instance.Placement.DedicatedResourcePackTenancy != nil {
			_ = d.Set("dedicated_resource_pack_tenancy", instance.Placement.DedicatedResourcePackTenancy)
		}
		// if len(instance.Placement.DedicatedResourcePackIds) > 0 {
		// 	_ = d.Set("dedicated_resource_pack_ids", helper.StringsInterfaces(instance.Placement.DedicatedResourcePackIds))
		// }
		if instance.Placement.RackId != nil {
			_ = d.Set("rack_id", instance.Placement.RackId)
		}
	}

	_ = d.Set("instance_charge_type", instance.InstanceChargeType)
	_ = d.Set("instance_charge_type_prepaid_renew_flag", instance.RenewFlag)
	_ = d.Set("internet_charge_type", instance.InternetAccessible.InternetChargeType)
	_ = d.Set("internet_max_bandwidth_out", instance.InternetAccessible.InternetMaxBandwidthOut)
	_ = d.Set("vpc_id", instance.VirtualPrivateCloud.VpcId)
	_ = d.Set("subnet_id", instance.VirtualPrivateCloud.SubnetId)
	_ = d.Set("security_groups", instance.SecurityGroupIds)
	_ = d.Set("orderly_security_groups", instance.SecurityGroupIds)
	_ = d.Set("system_disk_type", instance.SystemDisk.DiskType)
	_ = d.Set("system_disk_size", instance.SystemDisk.DiskSize)
	_ = d.Set("system_disk_id", instance.SystemDisk.DiskId)
	_ = d.Set("instance_status", instance.InstanceState)
	_ = d.Set("create_time", instance.CreatedTime)
	_ = d.Set("expired_time", instance.ExpiredTime)
	_ = d.Set("cam_role_name", instance.CamRoleName)
	_ = d.Set("disable_api_termination", instance.DisableApiTermination)
	_ = d.Set("cpu", instance.CPU)
	_ = d.Set("memory", instance.Memory)
	_ = d.Set("os_name", instance.OsName)
	_ = d.Set("hpc_cluster_id", instance.HpcClusterId)
	_ = d.Set("ipv6_addresses", instance.IPv6Addresses)
	_ = d.Set("public_ipv6_addresses", instance.PublicIPv6Addresses)

	if instance.Uuid != nil {
		_ = d.Set("uuid", instance.Uuid)
	}

	if instance.DisasterRecoverGroupId != nil {
		_ = d.Set("placement_group_id", instance.DisasterRecoverGroupId)
	}

	if *instance.InstanceChargeType == CVM_CHARGE_TYPE_CDHPAID {
		_ = d.Set("cdh_instance_type", instance.InstanceType)
	}

	if _, ok := d.GetOkExists("allocate_public_ip"); !ok {
		_ = d.Set("allocate_public_ip", len(instance.PublicIpAddresses) > 0)
	}

	if instance.InternetAccessible != nil {
		if instance.InternetAccessible.IPv4AddressType != nil {
			_ = d.Set("ipv4_address_type", instance.InternetAccessible.IPv4AddressType)
		}

		if instance.InternetAccessible.IPv6AddressType != nil {
			_ = d.Set("ipv6_address_type", instance.InternetAccessible.IPv6AddressType)
		}

		if instance.InternetAccessible.AntiDDoSPackageId != nil {
			_ = d.Set("anti_ddos_package_id", instance.InternetAccessible.AntiDDoSPackageId)
		}
	}

	tagService := svctag.NewTagService(client)

	tags, err := tagService.DescribeResourceTags(ctx, "cvm", "instance", client.Region, d.Id())
	if err != nil {
		return err
	}

	// as attachment add tencentcloud:autoscaling:auto-scaling-group-id tag automatically
	// we should remove this tag, otherwise it will cause terraform state change
	delete(tags, "tencentcloud:autoscaling:auto-scaling-group-id")
	_ = d.Set("tags", tags)

	// set system_disk_name
	if instance.SystemDisk.DiskId != nil && strings.HasPrefix(*instance.SystemDisk.DiskId, "disk-") {
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			disks, err := cbsService.DescribeDiskList(ctx, []*string{instance.SystemDisk.DiskId})
			if err != nil {
				return tccommon.RetryError(err)
			}

			for i := range disks {
				disk := disks[i]
				if *disk.DiskState == "EXPANDING" {
					return resource.RetryableError(fmt.Errorf("data_disk[%d] is expending", i))
				}

				if *disk.DiskId == *instance.SystemDisk.DiskId {
					_ = d.Set("system_disk_name", disk.DiskName)
				}
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	// set data_disks
	var hasDataDisks, isCombineDataDisks, hasDataDisksId bool
	dataDiskList := make([]map[string]interface{}, 0, len(instance.DataDisks))
	diskSizeMap := map[string]*uint64{}
	diskOrderMap := make(map[string]int)
	dataDiskIds := make([]*string, 0, len(instance.DataDisks))
	refreshDataDisks := make([]interface{}, 0, len(instance.DataDisks))

	if v, ok := d.GetOk("data_disks"); ok {
		hasDataDisks = true
		// check has data disk id and name
		dataDisks := v.([]interface{})
		for _, item := range dataDisks {
			value := item.(map[string]interface{})
			if v, ok := value["data_disk_id"]; ok && v != nil {
				diskId := v.(string)
				if diskId != "" && strings.HasPrefix(diskId, "disk-") {
					dataDiskIds = append(dataDiskIds, &diskId)
					hasDataDisksId = true
				}
			}
		}
	}

	// refresh data disk name and size
	if hasDataDisksId && len(dataDiskIds) > 0 {
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			disks, err := cbsService.DescribeDiskList(ctx, dataDiskIds)
			if err != nil {
				return tccommon.RetryError(err)
			}

			if v, ok := d.GetOk("data_disks"); ok {
				dataDisks := v.([]interface{})
				for _, item := range dataDisks {
					value := item.(map[string]interface{})
					for _, item := range disks {
						if value["data_disk_id"].(string) == *item.DiskId {
							value["data_disk_name"] = *item.DiskName
							value["data_disk_size"] = int(*item.DiskSize)
							value["data_disk_type"] = *item.DiskType
							if item.KmsKeyId != nil {
								value["kms_key_id"] = *item.KmsKeyId
							}

							value["encrypt"] = *item.Encrypt
							value["throughput_performance"] = *item.ThroughputPerformance
							value["delete_with_instance"] = *item.DeleteWithInstance
							break
						}
					}
				}

				refreshDataDisks = dataDisks
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	// scene with has disks name
	if len(instance.DataDisks) > 0 && !hasDataDisksId {
		var diskIds []*string
		for i := range instance.DataDisks {
			id := instance.DataDisks[i].DiskId
			size := instance.DataDisks[i].DiskSize
			if id == nil {
				continue
			}

			if strings.HasPrefix(*id, "disk-") {
				diskIds = append(diskIds, id)
			} else {
				diskSizeMap[*id] = helper.Int64Uint64(*size)
			}
		}

		if len(diskIds) > 0 {
			err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				disks, err := cbsService.DescribeDiskList(ctx, diskIds)
				if err != nil {
					return tccommon.RetryError(err)
				}

				for i := range disks {
					disk := disks[i]
					if *disk.DiskState == "EXPANDING" {
						return resource.RetryableError(fmt.Errorf("data_disk[%d] is expending", i))
					}

					diskSizeMap[*disk.DiskId] = disk.DiskSize
					if hasDataDisks {
						items := strings.Split(*disk.DiskName, "_")
						diskOrder := items[len(items)-1]
						diskOrderInt, err := strconv.Atoi(diskOrder)
						if err != nil {
							isCombineDataDisks = true
							continue
						}

						diskOrderMap[*disk.DiskId] = diskOrderInt
					}
				}

				return nil
			})

			if err != nil {
				return err
			}

			tmpDataDisks := make([]interface{}, 0, len(instance.DataDisks))
			if v, ok := d.GetOk("data_disks"); ok {
				tmpDataDisks = v.([]interface{})
			}

			for _, disk := range instance.DataDisks {
				dataDisk := make(map[string]interface{})
				if !strings.HasPrefix(*disk.DiskId, "disk-") {
					continue
				}

				dataDisk["data_disk_id"] = disk.DiskId
				if disk.DiskId == nil {
					dataDisk["data_disk_size"] = disk.DiskSize
				} else if size, ok := diskSizeMap[*disk.DiskId]; ok {
					dataDisk["data_disk_size"] = size
				}

				dataDisk["data_disk_type"] = disk.DiskType
				dataDisk["data_disk_snapshot_id"] = disk.SnapshotId
				dataDisk["delete_with_instance"] = disk.DeleteWithInstance
				dataDisk["kms_key_id"] = disk.KmsKeyId
				dataDisk["encrypt"] = disk.Encrypt
				dataDisk["throughput_performance"] = disk.ThroughputPerformance
				dataDiskList = append(dataDiskList, dataDisk)
			}

			if hasDataDisks && !isCombineDataDisks {
				sort.SliceStable(dataDiskList, func(idx1, idx2 int) bool {
					dataDiskIdIdx1 := *dataDiskList[idx1]["data_disk_id"].(*string)
					dataDiskIdIdx2 := *dataDiskList[idx2]["data_disk_id"].(*string)
					return diskOrderMap[dataDiskIdIdx1] < diskOrderMap[dataDiskIdIdx2]
				})
			}

			// set data disk name
			finalDiskIds := make([]*string, 0, len(dataDiskList))
			for _, item := range dataDiskList {
				diskId := item["data_disk_id"].(*string)
				finalDiskIds = append(finalDiskIds, diskId)
			}

			if len(finalDiskIds) != 0 {
				err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
					disks, err := cbsService.DescribeDiskList(ctx, finalDiskIds)
					if err != nil {
						return tccommon.RetryError(err)
					}

					for _, disk := range disks {
						diskId := disk.DiskId
						for _, v := range dataDiskList {
							tmpDiskId := v["data_disk_id"].(*string)
							if *diskId == *tmpDiskId {
								v["data_disk_name"] = disk.DiskName
								break
							}
						}
					}

					return nil
				})

				if err != nil {
					return err
				}
			}

			sortedDataDiskList, err := sortDataDisks(tmpDataDisks, dataDiskList)
			if err != nil {
				return err
			}

			// set data disk delete_with_instance_prepaid
			for i := range sortedDataDiskList {
				sortedDataDiskList[i]["delete_with_instance_prepaid"] = false
				if hasDataDisks {
					tmpDataDisk := tmpDataDisks[i].(map[string]interface{})
					if deleteWithInstancePrepaidBool, ok := tmpDataDisk["delete_with_instance_prepaid"].(bool); ok {
						sortedDataDiskList[i]["delete_with_instance_prepaid"] = deleteWithInstancePrepaidBool
					}
				}
			}

			_ = d.Set("data_disks", sortedDataDiskList)
		}
	} else if len(instance.DataDisks) > 0 && hasDataDisksId {
		// scene with no disks name
		dDiskHash := make([]map[string]interface{}, 0)
		// get source disk hash
		if v, ok := d.GetOk("data_disks"); ok {
			dataDisks := v.([]interface{})
			if hasDataDisksId {
				dataDisks = refreshDataDisks
			}

			for index, item := range dataDisks {
				value := item.(map[string]interface{})
				tmpMap := make(map[string]interface{})
				diskName := strconv.Itoa(index)
				diskType := value["data_disk_type"].(string)
				diskSize := int64(value["data_disk_size"].(int))
				deleteWithInstance := value["delete_with_instance"].(bool)
				kmsKeyId := value["kms_key_id"].(string)
				encrypt := value["encrypt"].(bool)
				if tmpV, ok := value["data_disk_name"].(string); ok && tmpV != "" {
					diskName = tmpV
				}

				diskObj := diskHash{
					diskType:           diskType,
					diskSize:           diskSize,
					deleteWithInstance: deleteWithInstance,
					kmsKeyId:           kmsKeyId,
					encrypt:            encrypt,
				}

				// set hash
				tmpMap[diskName] = getDataDiskHash(diskObj)
				tmpMap["index"] = index
				tmpMap["flag"] = 0
				dDiskHash = append(dDiskHash, tmpMap)
			}
		}

		tmpDataDiskMap := make(map[int]interface{}, 0)
		var diskIds []*string
		var cbsDisks []*cbs.Disk
		for i := range instance.DataDisks {
			id := instance.DataDisks[i].DiskId
			if id == nil {
				continue
			}

			if strings.HasPrefix(*id, "disk-") {
				diskIds = append(diskIds, id)
			}
		}

		if len(diskIds) > 0 {
			err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				cbsDisks, err = cbsService.DescribeDiskList(ctx, diskIds)
				if err != nil {
					return tccommon.RetryError(err)
				}

				for i := range cbsDisks {
					disk := cbsDisks[i]
					if *disk.DiskState == "EXPANDING" {
						return resource.RetryableError(fmt.Errorf("data_disk[%d] is expending", i))
					}
				}

				return nil
			})

			if err != nil {
				return err
			}

			// make data disks data
			sourceDataDisks := make([]*map[string]interface{}, 0)
			for _, cvmDisk := range instance.DataDisks {
				for _, cbsDisk := range cbsDisks {
					if *cvmDisk.DiskId == *cbsDisk.DiskId {
						dataDisk := make(map[string]interface{}, 10)
						dataDisk["data_disk_id"] = cvmDisk.DiskId
						dataDisk["data_disk_size"] = cvmDisk.DiskSize
						dataDisk["data_disk_name"] = cbsDisk.DiskName
						dataDisk["data_disk_type"] = cvmDisk.DiskType
						dataDisk["data_disk_snapshot_id"] = cvmDisk.SnapshotId
						dataDisk["delete_with_instance"] = cvmDisk.DeleteWithInstance
						dataDisk["kms_key_id"] = cvmDisk.KmsKeyId
						dataDisk["encrypt"] = cvmDisk.Encrypt
						dataDisk["throughput_performance"] = cvmDisk.ThroughputPerformance
						dataDisk["flag"] = 0
						sourceDataDisks = append(sourceDataDisks, &dataDisk)
						break
					}
				}
			}

			// has set disk name first
			for v := range sourceDataDisks {
				for i := range dDiskHash {
					var kmsKeyId *string
					disk := *sourceDataDisks[v]
					diskFlag := disk["flag"].(int)
					diskName := disk["data_disk_name"].(*string)
					diskType := disk["data_disk_type"].(*string)
					diskSize := disk["data_disk_size"].(*int64)
					deleteWithInstance := disk["delete_with_instance"].(*bool)
					if v, ok := disk["kms_key_id"].(*string); ok && v != nil {
						kmsKeyId = v
					} else {
						kmsKeyId = helper.String("")
					}
					encrypt := disk["encrypt"].(*bool)
					tmpHash := getDataDiskHash(diskHash{
						diskType:           *diskType,
						diskSize:           *diskSize,
						deleteWithInstance: *deleteWithInstance,
						kmsKeyId:           *kmsKeyId,
						encrypt:            *encrypt,
					})

					// get disk name
					hashItem := dDiskHash[i]
					if _, ok := hashItem[*diskName]; ok {
						// check hash and flag
						if hashItem["flag"] == 0 && diskFlag == 0 && tmpHash == hashItem[*diskName] {
							dataDisk := make(map[string]interface{}, 8)
							dataDisk["data_disk_id"] = disk["data_disk_id"]
							dataDisk["data_disk_size"] = disk["data_disk_size"]
							dataDisk["data_disk_name"] = disk["data_disk_name"]
							dataDisk["data_disk_type"] = disk["data_disk_type"]
							dataDisk["data_disk_snapshot_id"] = disk["data_disk_snapshot_id"]
							dataDisk["delete_with_instance"] = disk["delete_with_instance"]
							dataDisk["kms_key_id"] = disk["kms_key_id"]
							dataDisk["encrypt"] = disk["encrypt"]
							dataDisk["throughput_performance"] = disk["throughput_performance"]
							tmpDataDiskMap[hashItem["index"].(int)] = dataDisk
							hashItem["flag"] = 1
							disk["flag"] = 1
							break
						}
					}
				}
			}

			// no set disk name last
			for v := range sourceDataDisks {
				for i := range dDiskHash {
					var kmsKeyId *string
					disk := *sourceDataDisks[v]
					diskFlag := disk["flag"].(int)
					diskType := disk["data_disk_type"].(*string)
					diskSize := disk["data_disk_size"].(*int64)
					deleteWithInstance := disk["delete_with_instance"].(*bool)
					if v, ok := disk["kms_key_id"].(*string); ok && v != nil {
						kmsKeyId = v
					} else {
						kmsKeyId = helper.String("")
					}
					encrypt := disk["encrypt"].(*bool)
					tmpHash := getDataDiskHash(diskHash{
						diskType:           *diskType,
						diskSize:           *diskSize,
						deleteWithInstance: *deleteWithInstance,
						kmsKeyId:           *kmsKeyId,
						encrypt:            *encrypt,
					})

					// check hash and flag
					hashItem := dDiskHash[i]
					if hashItem["flag"] == 0 && diskFlag == 0 && tmpHash == hashItem[strconv.Itoa(i)] {
						dataDisk := make(map[string]interface{}, 8)
						dataDisk["data_disk_id"] = disk["data_disk_id"]
						dataDisk["data_disk_size"] = disk["data_disk_size"]
						dataDisk["data_disk_name"] = disk["data_disk_name"]
						dataDisk["data_disk_type"] = disk["data_disk_type"]
						dataDisk["data_disk_snapshot_id"] = disk["data_disk_snapshot_id"]
						dataDisk["delete_with_instance"] = disk["delete_with_instance"]
						dataDisk["kms_key_id"] = disk["kms_key_id"]
						dataDisk["encrypt"] = disk["encrypt"]
						dataDisk["throughput_performance"] = disk["throughput_performance"]
						tmpDataDiskMap[hashItem["index"].(int)] = dataDisk
						hashItem["flag"] = 1
						disk["flag"] = 1
						break
					}
				}
			}

			keys := make([]int, 0, len(tmpDataDiskMap))
			for k := range tmpDataDiskMap {
				keys = append(keys, k)
			}

			sort.Ints(keys)
			for _, v := range keys {
				tmpDataDisk := tmpDataDiskMap[v].(map[string]interface{})
				dataDiskList = append(dataDiskList, tmpDataDisk)
			}

			// set data disk delete_with_instance_prepaid
			if v, ok := d.GetOk("data_disks"); ok {
				tmpDataDisks := v.([]interface{})
				for i := range dataDiskList {
					dataDiskList[i]["delete_with_instance_prepaid"] = false
					if hasDataDisks {
						tmpDataDisk := tmpDataDisks[i].(map[string]interface{})
						if deleteWithInstancePrepaidBool, ok := tmpDataDisk["delete_with_instance_prepaid"].(bool); ok {
							dataDiskList[i]["delete_with_instance_prepaid"] = deleteWithInstancePrepaidBool
						}
					}
				}
			}

			_ = d.Set("data_disks", dataDiskList)
		}
	} else {
		_ = d.Set("data_disks", dataDiskList)
	}

	if len(instance.PrivateIpAddresses) > 0 {
		_ = d.Set("private_ip", instance.PrivateIpAddresses[0])
	}

	if len(instance.PublicIpAddresses) > 0 {
		_ = d.Set("public_ip", instance.PublicIpAddresses[0])
	}

	if len(instance.LoginSettings.KeyIds) > 0 {
		_ = d.Set("key_name", instance.LoginSettings.KeyIds[0])
		_ = d.Set("key_ids", instance.LoginSettings.KeyIds)
	} else {
		_ = d.Set("key_name", "")
		_ = d.Set("key_ids", []*string{})
	}

	if instance.LoginSettings.KeepImageLogin != nil {
		_ = d.Set("keep_image_login", *instance.LoginSettings.KeepImageLogin == CVM_IMAGE_LOGIN)
	}

	if *instance.InstanceState == CVM_STATUS_STOPPED {
		_ = d.Set("running_flag", false)
	} else {
		_ = d.Set("running_flag", true)
	}

	var instanceAttribute *cvm.InstanceAttribute
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		request := cvm.NewDescribeInstancesAttributesRequest()
		request.InstanceIds = helper.Strings([]string{instanceId})
		request.Attributes = helper.Strings([]string{"UserData"})
		response, errRet := client.UseCvmClient().DescribeInstancesAttributes(request)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}

		if len(response.Response.InstanceSet) > 0 {
			instanceAttribute = response.Response.InstanceSet[0]
		}
		return nil
	})

	if err != nil {
		return err
	}
	if instanceAttribute != nil && instanceAttribute.Attributes != nil && instanceAttribute.Attributes.UserData != nil {
		_ = d.Set("user_data", instanceAttribute.Attributes.UserData)
		userDataRaw, e := base64.StdEncoding.DecodeString(*(instanceAttribute.Attributes.UserData))
		if e != nil {
			return e
		}
		_ = d.Set("user_data_raw", string(userDataRaw))
	}

	if instance.VirtualPrivateCloud != nil && instance.VirtualPrivateCloud.Ipv6AddressCount != nil {
		_ = d.Set("ipv6_address_count", instance.VirtualPrivateCloud.Ipv6AddressCount)
	}

	return nil
}

func resourceTencentCloudInstanceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	defer tccommon.LogElapsed("resource.tencentcloud_instance.update")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		instanceId = d.Id()
		cvmService = CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	d.Partial(true)

	// Get the latest instance info from actual resource.
	instanceInfo, err := cvmService.DescribeInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	var (
		periodSet         = false
		renewFlagSet      = false
		expectChargeType  = CVM_CHARGE_TYPE_POSTPAID
		currentChargeType = *instanceInfo.InstanceChargeType
	)

	chargeType, chargeOk := d.GetOk("instance_charge_type")
	if chargeOk {
		expectChargeType = chargeType.(string)
	}

	if d.HasChange("instance_charge_type") && expectChargeType != currentChargeType {
		var (
			period    = -1
			renewFlag string
		)

		if v, ok := d.GetOk("instance_charge_type_prepaid_period"); ok {
			period = v.(int)
		}

		if v, ok := d.GetOk("instance_charge_type_prepaid_renew_flag"); ok {
			renewFlag = v.(string)
		}

		// change charge type
		err := cvmService.ModifyInstanceChargeType(ctx, instanceId, expectChargeType, period, renewFlag)
		if err != nil {
			return err
		}

		// query cvm status
		err = waitForOperationFinished(d, meta, 5*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
		if err != nil {
			return err
		}

		periodSet = true
		renewFlagSet = true
	}

	// When instance is prepaid but period was empty and set to 1, skip this case.
	op, np := d.GetChange("instance_charge_type_prepaid_period")
	if _, ok := op.(int); !ok && np.(int) == 1 {
		periodSet = true
	}

	if d.HasChange("instance_charge_type_prepaid_period") && !periodSet {
		chargeType := d.Get("instance_charge_type").(string)
		period := d.Get("instance_charge_type_prepaid_period").(int)
		renewFlag := ""

		if v, ok := d.GetOk("instance_charge_type_prepaid_renew_flag"); ok {
			renewFlag = v.(string)
		}

		err := cvmService.ModifyInstanceChargeType(ctx, instanceId, chargeType, period, renewFlag)
		if err != nil {
			return err
		}

		// query cvm status
		err = waitForOperationFinished(d, meta, 5*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
		if err != nil {
			return err
		}

		renewFlagSet = true
	}

	if d.HasChange("instance_charge_type_prepaid_renew_flag") && !renewFlagSet {
		//renew api
		err := cvmService.ModifyRenewParam(ctx, instanceId, d.Get("instance_charge_type_prepaid_renew_flag").(string))
		if err != nil {
			return err
		}

		//check success
		err = waitForOperationFinished(d, meta, 2*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
		if err != nil {
			return err
		}

		time.Sleep(tccommon.ReadRetryTimeout)
	}

	if d.HasChange("instance_name") {
		err := cvmService.ModifyInstanceName(ctx, instanceId, d.Get("instance_name").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("hostname") {
		err := cvmService.ModifyHostName(ctx, instanceId, d.Get("hostname").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("disable_api_termination") {
		err := cvmService.ModifyDisableApiTermination(ctx, instanceId, d.Get("disable_api_termination").(bool))
		if err != nil {
			return err
		}
	}

	if d.HasChange("cam_role_name") {
		err := cvmService.ModifyCamRoleName(ctx, instanceId, d.Get("cam_role_name").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("security_groups") {
		securityGroups := d.Get("security_groups").(*schema.Set).List()
		securityGroupIds := make([]*string, 0, len(securityGroups))
		for _, securityGroup := range securityGroups {
			securityGroupIds = append(securityGroupIds, helper.String(securityGroup.(string)))
		}

		err := cvmService.ModifySecurityGroups(ctx, instanceId, securityGroupIds)
		if err != nil {
			return err
		}
	}

	if d.HasChange("orderly_security_groups") {
		orderlySecurityGroups := d.Get("orderly_security_groups").([]interface{})
		orderlySecurityGroupIds := make([]*string, 0, len(orderlySecurityGroups))
		for _, securityGroup := range orderlySecurityGroups {
			orderlySecurityGroupIds = append(orderlySecurityGroupIds, helper.String(securityGroup.(string)))
		}

		err := cvmService.ModifySecurityGroups(ctx, instanceId, orderlySecurityGroupIds)
		if err != nil {
			return err
		}
	}

	if d.HasChange("project_id") {
		projectId := d.Get("project_id").(int)
		err := cvmService.ModifyProjectId(ctx, instanceId, int64(projectId))
		if err != nil {
			return err
		}
	}

	// Reset Instance
	// Keep Login Info
	if d.HasChange("image_id") ||
		d.HasChange("disable_security_service") ||
		d.HasChange("disable_monitor_service") ||
		d.HasChange("disable_automation_service") ||
		d.HasChange("keep_image_login") {

		request := cvm.NewResetInstanceRequest()
		request.InstanceId = helper.String(d.Id())

		if v, ok := d.GetOk("image_id"); ok {
			request.ImageId = helper.String(v.(string))
		}

		// enhanced service
		var (
			enhancedService     cvm.EnhancedService
			enhancedServiceFlag bool
		)

		if v, ok := d.GetOkExists("disable_security_service"); ok {
			securityService := !(v.(bool))
			enhancedService.SecurityService = &cvm.RunSecurityServiceEnabled{
				Enabled: &securityService,
			}
			enhancedServiceFlag = true
		}

		if v, ok := d.GetOkExists("disable_monitor_service"); ok {
			monitorService := !(v.(bool))
			enhancedService.MonitorService = &cvm.RunMonitorServiceEnabled{
				Enabled: &monitorService,
			}
			enhancedServiceFlag = true
		}

		if v, ok := d.GetOkExists("disable_automation_service"); ok {
			automationService := !(v.(bool))
			enhancedService.AutomationService = &cvm.RunAutomationServiceEnabled{
				Enabled: &automationService,
			}
			enhancedServiceFlag = true
		}

		if enhancedServiceFlag {
			request.EnhancedService = &enhancedService
		}

		// login
		var (
			loginSettings     cvm.LoginSettings
			loginSettingsFlag bool
		)

		if v, ok := d.GetOk("key_name"); ok {
			loginSettings.KeyIds = []*string{helper.String(v.(string))}
			loginSettingsFlag = true
		}

		if v, ok := d.GetOk("key_ids"); ok {
			keyIds := v.(*schema.Set).List()
			if len(keyIds) > 0 {
				loginSettings.KeyIds = helper.InterfacesStringsPoint(keyIds)
				loginSettingsFlag = true
			}
		}

		if v, ok := d.GetOk("password"); ok {
			loginSettings.Password = helper.String(v.(string))
			loginSettingsFlag = true
		}

		if v, ok := d.GetOkExists("keep_image_login"); ok {
			if v.(bool) {
				loginSettings.KeepImageLogin = helper.String(CVM_IMAGE_LOGIN)
			} else {
				loginSettings.KeepImageLogin = helper.String(CVM_IMAGE_LOGIN_NOT)
			}

			loginSettingsFlag = true
		}

		if loginSettingsFlag {
			request.LoginSettings = &loginSettings
		}

		if err := cvmService.ResetInstance(ctx, request); err != nil {
			return err
		}
		// Modify Login Info Directly
	} else {
		if d.HasChange("password") {
			err := cvmService.ModifyPassword(ctx, instanceId, d.Get("password").(string))
			if err != nil {
				return err
			}

			err = waitForOperationFinished(d, meta, 2*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
			if err != nil {
				return err
			}
		}

		if d.HasChange("key_name") {
			o, n := d.GetChange("key_name")
			oldKeyId := o.(string)
			keyId := n.(string)

			if oldKeyId != "" {
				err := cvmService.UnbindKeyPair(ctx, []*string{&oldKeyId}, []*string{&instanceId})
				if err != nil {
					return err
				}

				err = waitForOperationFinished(d, meta, 2*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
				if err != nil {
					return err
				}
			}

			if keyId != "" {
				err = cvmService.BindKeyPair(ctx, []*string{&keyId}, instanceId)
				if err != nil {
					return err
				}

				err = waitForOperationFinished(d, meta, 2*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
				if err != nil {
					return err
				}
			}
		}

		// support remove old `key_name` to `key_ids`, so do not follow "else"
		if d.HasChange("key_ids") {
			o, n := d.GetChange("key_ids")
			ov := o.(*schema.Set)
			nv := n.(*schema.Set)
			adds := nv.Difference(ov)
			removes := ov.Difference(nv)
			adds.Remove("")
			removes.Remove("")

			if removes.Len() > 0 {
				err := cvmService.UnbindKeyPair(ctx, helper.InterfacesStringsPoint(removes.List()), []*string{&instanceId})
				if err != nil {
					return err
				}

				err = waitForOperationFinished(d, meta, 2*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
				if err != nil {
					return err
				}
			}

			if adds.Len() > 0 {
				err = cvmService.BindKeyPair(ctx, helper.InterfacesStringsPoint(adds.List()), instanceId)
				if err != nil {
					return err
				}

				err = waitForOperationFinished(d, meta, 2*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
				if err != nil {
					return err
				}
			}
		}
	}

	if d.HasChange("data_disks") {
		o, n := d.GetChange("data_disks")
		ov := o.([]interface{})
		nv := n.([]interface{})

		if len(ov) != len(nv) {
			return fmt.Errorf("error: data disk count has changed (%d -> %d) but doesn't support add or remove for now", len(ov), len(nv))
		}

		cbsService := svccbs.NewCbsService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

		for i := range nv {
			sizeKey := fmt.Sprintf("data_disks.%d.data_disk_size", i)
			idKey := fmt.Sprintf("data_disks.%d.data_disk_id", i)
			nameKey := fmt.Sprintf("data_disks.%d.data_disk_name", i)
			if d.HasChange(sizeKey) {
				size := d.Get(sizeKey).(int)
				diskId := d.Get(idKey).(string)
				err := cbsService.ResizeDisk(ctx, diskId, size)
				if err != nil {
					return fmt.Errorf("an error occurred when modifying data disk size: %s, reason: %s", sizeKey, err.Error())
				}
			}
			if d.HasChange(nameKey) {
				name := d.Get(nameKey).(string)
				diskId := d.Get(idKey).(string)
				err := cbsService.ModifyDiskAttributes(ctx, "", diskId, name, -1, "")
				if err != nil {
					return fmt.Errorf("an error occurred when modifying data disk name: %s, reason: %s", name, err.Error())
				}
			}
		}
	}

	var flag bool
	if d.HasChange("running_flag") {
		flag = d.Get("running_flag").(bool)
		if err := switchInstance(&cvmService, ctx, d, flag); err != nil {
			return err
		}
	}

	if d.HasChange("system_disk_size") || d.HasChange("system_disk_type") {
		size := d.Get("system_disk_size").(int)
		diskType := d.Get("system_disk_type").(string)
		//diskId := d.Get("system_disk_id").(string)
		req := cvm.NewResizeInstanceDisksRequest()
		req.InstanceId = &instanceId
		req.ForceStop = helper.Bool(true)
		req.SystemDisk = &cvm.SystemDisk{
			DiskSize: helper.IntInt64(size),
			DiskType: &diskType,
		}
		if v, ok := d.GetOkExists("system_disk_resize_online"); ok {
			req.ResizeOnline = helper.Bool(v.(bool))
		}

		err := cvmService.ResizeInstanceDisks(ctx, req)
		if err != nil {
			return fmt.Errorf("an error occurred when modifying system_disk, reason: %s", err.Error())
		}

		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			instance, err := cvmService.DescribeInstanceById(ctx, instanceId)
			if err != nil {
				return tccommon.RetryError(err)
			}

			if instance != nil && instance.LatestOperationState != nil {
				if *instance.InstanceState == "FAILED" {
					return resource.NonRetryableError(fmt.Errorf("instance operation failed"))
				}

				if *instance.InstanceState == "OPERATING" {
					return resource.RetryableError(fmt.Errorf("instance operating"))
				}
			}

			if instance != nil && instance.SystemDisk != nil {
				//wait until disk result as expected
				if *instance.SystemDisk.DiskType != diskType || int(*instance.SystemDisk.DiskSize) != size {
					return resource.RetryableError(fmt.Errorf("waiting for expanding success"))
				}
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	if d.HasChange("system_disk_name") {
		systemDiskName := d.Get("system_disk_name").(string)
		if v, ok := d.GetOk("system_disk_id"); ok {
			systemDiskId := v.(string)
			cbsService := svccbs.NewCbsService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
			err := cbsService.ModifyDiskAttributes(ctx, "", systemDiskId, systemDiskName, -1, "")
			if err != nil {
				return fmt.Errorf("an error occurred when modifying system disk name %s, reason: %s", systemDiskName, err.Error())
			}
		} else {
			return fmt.Errorf("system disk name do not support change because of no system disk ID.")
		}
	}

	if d.HasChange("instance_type") {
		err := cvmService.ModifyInstanceType(ctx, instanceId, d.Get("instance_type").(string))
		if err != nil {
			return err
		}

		err = waitForOperationFinished(d, meta, 2*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
		if err != nil {
			return err
		}
	}

	if d.HasChange("cdh_instance_type") {
		err := cvmService.ModifyInstanceType(ctx, instanceId, d.Get("cdh_instance_type").(string))
		if err != nil {
			return err
		}

		err = waitForOperationFinished(d, meta, 2*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
		if err != nil {
			return err
		}
	}

	if d.HasChange("vpc_id") || d.HasChange("subnet_id") || d.HasChange("private_ip") {
		vpcId := d.Get("vpc_id").(string)
		subnetId := d.Get("subnet_id").(string)
		privateIp := d.Get("private_ip").(string)
		err := cvmService.ModifyVpc(ctx, instanceId, vpcId, subnetId, privateIp)
		if err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		oldInterface, newInterface := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldInterface.(map[string]interface{}), newInterface.(map[string]interface{}))
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := tccommon.BuildTagResourceName("cvm", "instance", region, instanceId)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}

		//except instance ,system disk and data disk will be tagged
		//keep logical consistence with the console
		//tag system disk
		if systemDiskId, ok := d.GetOk("system_disk_id"); ok {
			if systemDiskId.(string) != "" {
				resourceName = tccommon.BuildTagResourceName("cvm", "volume", region, systemDiskId.(string))
				if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
					return err
				}
			}
		}

		//tag disk ids
		if dataDisks, ok := d.GetOk("data_disks"); ok {
			dataDiskList := dataDisks.([]interface{})
			for _, dataDisk := range dataDiskList {
				disk := dataDisk.(map[string]interface{})
				dataDiskId := disk["data_disk_id"].(string)
				resourceName = tccommon.BuildTagResourceName("cvm", "volume", region, dataDiskId)
				if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
					return err
				}
			}
		}
	}

	if d.HasChange("internet_max_bandwidth_out") {
		chargeType := d.Get("internet_charge_type").(string)
		bandWidthOut := int64(d.Get("internet_max_bandwidth_out").(int))
		if chargeType != "TRAFFIC_POSTPAID_BY_HOUR" && chargeType != "BANDWIDTH_POSTPAID_BY_HOUR" && chargeType != "BANDWIDTH_PACKAGE" {
			return fmt.Errorf("charge type should be one of `TRAFFIC_POSTPAID_BY_HOUR BANDWIDTH_POSTPAID_BY_HOUR BANDWIDTH_PACKAGE` when adjusting internet_max_bandwidth_out")
		}

		err := cvmService.ModifyInternetMaxBandwidthOut(ctx, instanceId, chargeType, bandWidthOut)
		if err != nil {
			return err
		}

		err = waitForOperationFinished(d, meta, 2*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
		if err != nil {
			return err
		}
	}

	if d.HasChange("user_data") {
		err := cvmService.ModifyUserData(ctx, instanceId, d.Get("user_data").(string))
		if err != nil {
			return err
		}

		err = waitForOperationFinished(d, meta, 2*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
		if err != nil {
			return err
		}
	}

	if d.HasChange("user_data_raw") {
		userDataRaw := d.Get("user_data_raw").(string)
		userData := base64.StdEncoding.EncodeToString([]byte(userDataRaw))
		err := cvmService.ModifyUserData(ctx, instanceId, userData)
		if err != nil {
			return err
		}

		err = waitForOperationFinished(d, meta, 2*tccommon.ReadRetryTimeout, CVM_LATEST_OPERATION_STATE_OPERATING, false)
		if err != nil {
			return err
		}
	}

	if d.HasChange("placement_group_id") || d.HasChange("force_replace_placement_group_id") {
		oldPGI, newPGI := d.GetChange("placement_group_id")
		oldPGIStr := oldPGI.(string)
		newPGIStr := newPGI.(string)
		if newPGIStr == "" {
			// wait cvm support delete DisasterRecoverGroupId
			return fmt.Errorf("Deleting `placement_group_id` is not currently supported.")
		} else {
			if oldPGIStr == newPGIStr {
				return fmt.Errorf("It is not possible to change only `force_replace_placement_group_id`, it needs to be modified together with `placement_group_id`.")
			}

			request := cvm.NewModifyInstancesDisasterRecoverGroupRequest()
			if v, ok := d.GetOkExists("force_replace_placement_group_id"); ok {
				request.Force = helper.Bool(v.(bool))
			}

			request.InstanceIds = helper.Strings([]string{instanceId})
			request.DisasterRecoverGroupId = helper.String(newPGIStr)
			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ModifyInstancesDisasterRecoverGroup(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				return err
			}

			// wait
			err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
				if errRet != nil {
					return tccommon.RetryError(errRet, tccommon.InternalError)
				}

				if instance != nil && *instance.InstanceState == CVM_STATUS_LAUNCH_FAILED {
					//LatestOperationCodeMode
					if instance.LatestOperationErrorMsg != nil {
						return resource.NonRetryableError(fmt.Errorf("cvm instance %s launch failed. Error msg: %s.\n", *instance.InstanceId, *instance.LatestOperationErrorMsg))
					}

					return resource.NonRetryableError(fmt.Errorf("cvm instance %s launch failed, this resource will not be stored to tfstate and will auto removed\n.", *instance.InstanceId))
				}

				if instance != nil && *instance.InstanceState == CVM_STATUS_RUNNING {
					return nil
				}

				return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
			})

			if err != nil {
				return err
			}
		}
	}

	d.Partial(false)

	return resourceTencentCloudInstanceRead(d, meta)
}

func resourceTencentCloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_instance.delete")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		cvmService = CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	instanceId := d.Id()
	//check is force delete or not
	forceDelete := d.Get("force_delete").(bool)
	instanceChargeType := d.Get("instance_charge_type").(string)

	instance, err := cvmService.DescribeInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	var releaseAddress bool
	if v, ok := d.GetOkExists("release_address"); ok {
		releaseAddress = v.(bool)
	}
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		errRet := cvmService.DeleteInstance(ctx, instanceId, releaseAddress)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}

		return nil
	})

	if err != nil {
		return err
	}

	//check recycling
	notExist := false

	//check exist
	err = resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}

		if instance == nil {
			notExist = true
			return nil
		}

		if *instance.InstanceState == CVM_STATUS_SHUTDOWN && *instance.LatestOperationState != CVM_LATEST_OPERATION_STATE_OPERATING {
			//in recycling
			return nil
		}

		return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
	})

	if err != nil {
		return err
	}

	vpcService := vpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	if notExist {
		err := waitIpRelease(ctx, vpcService, instance)
		if err != nil {
			return err
		}

		return nil
	}

	if instanceChargeType == CVM_CHARGE_TYPE_PREPAID {
		if v, ok := d.GetOk("data_disks"); ok {
			dataDisks := v.([]interface{})
			for _, d := range dataDisks {
				value := d.(map[string]interface{})
				deleteWithInstancePrepaid := value["delete_with_instance_prepaid"].(bool)
				if deleteWithInstancePrepaid {
					diskId := value["data_disk_id"].(string)
					cbsService := svccbs.NewCbsService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
					err = resource.Retry(tccommon.ReadRetryTimeout*2, func() *resource.RetryError {
						diskInfo, e := cbsService.DescribeDiskById(ctx, diskId)
						if e != nil {
							return tccommon.RetryError(e, tccommon.InternalError)
						}

						if *diskInfo.DiskState != svccbs.CBS_STORAGE_STATUS_ATTACHED {
							return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *diskInfo.DiskState))
						}

						return nil
					})

					if err != nil {
						log.Printf("[CRITAL]%s delete cbs failed, reason:%s\n ", logId, err.Error())
						return err
					}

					err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
						e := cbsService.DetachDisk(ctx, diskId, instanceId)
						if e != nil {
							return tccommon.RetryError(e, tccommon.InternalError)
						}

						return nil
					})

					if err != nil {
						log.Printf("[CRITAL]%s detach cbs failed, reason:%s\n ", logId, err.Error())
						return err
					}

					err = resource.Retry(tccommon.ReadRetryTimeout*2, func() *resource.RetryError {
						diskInfo, e := cbsService.DescribeDiskById(ctx, diskId)
						if e != nil {
							return tccommon.RetryError(e, tccommon.InternalError)
						}

						if *diskInfo.DiskState != svccbs.CBS_STORAGE_STATUS_UNATTACHED {
							return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *diskInfo.DiskState))
						}

						return nil
					})

					if err != nil {
						log.Printf("[CRITAL]%s read cbs status failed, reason:%s\n ", logId, err.Error())
						return err
					}

					err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
						e := cbsService.DeleteDiskById(ctx, diskId)
						if e != nil {
							return tccommon.RetryError(e, tccommon.InternalError)
						}

						return nil
					})

					if err != nil {
						log.Printf("[CRITAL]%s delete cbs failed, reason:%s\n ", logId, err.Error())
						return err
					}

					err = resource.Retry(tccommon.ReadRetryTimeout*2, func() *resource.RetryError {
						diskInfo, e := cbsService.DescribeDiskById(ctx, diskId)
						if e != nil {
							return tccommon.RetryError(e, tccommon.InternalError)
						}

						if *diskInfo.DiskState != svccbs.CBS_STORAGE_STATUS_TORECYCLE {
							return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *diskInfo.DiskState))
						}

						return nil
					})

					if err != nil {
						log.Printf("[CRITAL]%s read cbs status failed, reason:%s\n ", logId, err.Error())
						return err
					}
				}
			}
		}
	}

	if !forceDelete {
		return nil
	}

	// exist in recycle, delete again
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		errRet := cvmService.DeleteInstance(ctx, instanceId, releaseAddress)
		//when state is terminating, do not delete but check exist
		if errRet != nil {
			//check InvalidInstanceState.Terminating
			ee, ok := errRet.(*sdkErrors.TencentCloudSDKError)
			if !ok {
				return tccommon.RetryError(errRet)
			}

			if ee.Code == "InvalidInstanceState.Terminating" {
				return nil
			}

			return tccommon.RetryError(errRet, "OperationDenied.InstanceOperationInProgress")
		}

		return nil
	})

	if err != nil {
		return err
	}

	//describe and check not exist
	err = resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}

		if instance == nil {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
	})

	if err != nil {
		return err
	}

	if v, ok := d.GetOk("data_disks"); ok {
		dataDisks := v.([]interface{})
		for _, d := range dataDisks {
			value := d.(map[string]interface{})
			diskId := value["data_disk_id"].(string)
			deleteWithInstance := value["delete_with_instance"].(bool)
			if deleteWithInstance && instanceChargeType == CVM_CHARGE_TYPE_POSTPAID {
				cbsService := svccbs.NewCbsService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
				err = resource.Retry(tccommon.ReadRetryTimeout*2, func() *resource.RetryError {
					diskInfo, e := cbsService.DescribeDiskById(ctx, diskId)
					if e != nil {
						return tccommon.RetryError(e, tccommon.InternalError)
					}

					if *diskInfo.DiskState != svccbs.CBS_STORAGE_STATUS_UNATTACHED {
						return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *diskInfo.DiskState))
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s delete cbs failed, reason:%s\n ", logId, err.Error())
					return err
				}

				err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					e := cbsService.DeleteDiskById(ctx, diskId)
					if e != nil {
						return tccommon.RetryError(e, tccommon.InternalError)
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s delete cbs failed, reason:%s\n ", logId, err.Error())
					return err
				}

				err = resource.Retry(tccommon.ReadRetryTimeout*2, func() *resource.RetryError {
					diskInfo, e := cbsService.DescribeDiskById(ctx, diskId)
					if e != nil {
						return tccommon.RetryError(e, tccommon.InternalError)
					}

					if *diskInfo.DiskState == svccbs.CBS_STORAGE_STATUS_TORECYCLE {
						return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *diskInfo.DiskState))
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s read cbs status failed, reason:%s\n ", logId, err.Error())
					return err
				}

				err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					e := cbsService.DeleteDiskById(ctx, diskId)
					if e != nil {
						return tccommon.RetryError(e, tccommon.InternalError)
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s delete cbs failed, reason:%s\n ", logId, err.Error())
					return err
				}
				err = resource.Retry(tccommon.ReadRetryTimeout*2, func() *resource.RetryError {
					diskInfo, e := cbsService.DescribeDiskById(ctx, diskId)
					if e != nil {
						return tccommon.RetryError(e, tccommon.InternalError)
					}

					if diskInfo != nil {
						return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *diskInfo.DiskState))
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s read cbs status failed, reason:%s\n ", logId, err.Error())
					return err
				}
			}

			deleteWithInstancePrepaid := value["delete_with_instance_prepaid"].(bool)
			if deleteWithInstancePrepaid && instanceChargeType == CVM_CHARGE_TYPE_PREPAID {
				cbsService := svccbs.NewCbsService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
				err = resource.Retry(tccommon.ReadRetryTimeout*2, func() *resource.RetryError {
					diskInfo, e := cbsService.DescribeDiskById(ctx, diskId)
					if e != nil {
						return tccommon.RetryError(e, tccommon.InternalError)
					}

					if *diskInfo.DiskState != svccbs.CBS_STORAGE_STATUS_TORECYCLE {
						return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *diskInfo.DiskState))
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s delete cbs failed, reason:%s\n ", logId, err.Error())
					return err
				}

				err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					e := cbsService.DeleteDiskById(ctx, diskId)
					if e != nil {
						return tccommon.RetryError(e, tccommon.InternalError)
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s delete cbs failed, reason:%s\n ", logId, err.Error())
					return err
				}

				err = resource.Retry(tccommon.ReadRetryTimeout*2, func() *resource.RetryError {
					diskInfo, e := cbsService.DescribeDiskById(ctx, diskId)
					if e != nil {
						return tccommon.RetryError(e, tccommon.InternalError)
					}

					if diskInfo != nil {
						return resource.RetryableError(fmt.Errorf("cbs storage status is %s", *diskInfo.DiskState))
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s read cbs status failed, reason:%s\n ", logId, err.Error())
					return err
				}
			}
		}
	}

	err = waitIpRelease(ctx, vpcService, instance)
	if err != nil {
		return err
	}

	return nil
}

func switchInstance(cvmService *CvmService, ctx context.Context, d *schema.ResourceData, flag bool) (err error) {
	instanceId := d.Id()
	if flag {
		err = cvmService.StartInstance(ctx, instanceId)
		if err != nil {
			return err
		}

		err = resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
			if errRet != nil {
				return tccommon.RetryError(errRet, tccommon.InternalError)
			}

			if instance != nil && *instance.InstanceState == CVM_STATUS_RUNNING {
				return nil
			}

			return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
		})

		if err != nil {
			return err
		}
	} else {
		stopType := d.Get("stop_type").(string)
		stoppedMode := d.Get("stopped_mode").(string)
		skipStopApi := false
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			// when retry polling instance status, stop instance should skipped
			if !skipStopApi {
				err := cvmService.StopInstance(ctx, instanceId, stopType, stoppedMode)
				if err != nil {
					return resource.NonRetryableError(err)
				}
			}

			instance, err := cvmService.DescribeInstanceById(ctx, instanceId)
			if err != nil {
				return resource.NonRetryableError(err)
			}

			if instance == nil {
				return resource.NonRetryableError(fmt.Errorf("instance %s not found", instanceId))
			}

			if instance.LatestOperationState != nil {
				operationState := *instance.LatestOperationState
				if operationState == "OPERATING" {
					skipStopApi = true
					return resource.RetryableError(fmt.Errorf("instance %s stop operating, retrying", instanceId))
				}

				if operationState == "FAILED" {
					skipStopApi = false
					return resource.RetryableError(fmt.Errorf("instance %s stop failed, retrying", instanceId))
				}
			}

			return nil
		})

		if err != nil {
			return err
		}

		err = resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
			if errRet != nil {
				return tccommon.RetryError(errRet, tccommon.InternalError)
			}

			if instance != nil && *instance.InstanceState == CVM_STATUS_STOPPED {
				return nil
			}

			return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func waitForOperationFinished(d *schema.ResourceData, meta interface{}, timeout time.Duration, state string, immediately bool) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	cvmService := CvmService{client}
	instanceId := d.Id()
	// We cannot catch LatestOperationState change immediately after modification returns, we must wait for LatestOperationState update to expected.
	if !immediately {
		time.Sleep(time.Second * 10)
	}

	err := resource.Retry(timeout, func() *resource.RetryError {
		instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}

		if instance == nil {
			return resource.NonRetryableError(fmt.Errorf("%s not exists", instanceId))
		}

		if instance.LatestOperationState == nil {
			return resource.RetryableError(fmt.Errorf("wait for operation update"))
		}

		if *instance.LatestOperationState == state {
			return resource.RetryableError(fmt.Errorf("waiting for instance %s operation", instanceId))
		}

		if *instance.LatestOperationState == CVM_LATEST_OPERATION_STATE_FAILED {
			return resource.NonRetryableError(fmt.Errorf("failed operation"))
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func waitIpRelease(ctx context.Context, vpcService vpc.VpcService, instance *cvm.Instance) error {
	// wait ip release
	if len(instance.PrivateIpAddresses) > 0 {
		params := make(map[string]interface{})
		params["VpcId"] = instance.VirtualPrivateCloud.VpcId
		params["SubnetId"] = instance.VirtualPrivateCloud.SubnetId
		params["IpAddresses"] = instance.PrivateIpAddresses
		err := resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			usedIpAddress, errRet := vpcService.DescribeVpcUsedIpAddressByFilter(ctx, params)
			if errRet != nil {
				return tccommon.RetryError(errRet, tccommon.InternalError)
			}

			if len(usedIpAddress) > 0 {
				return resource.RetryableError(fmt.Errorf("wait cvm private ip release..."))
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

type diskHash struct {
	diskType           string
	diskSize           int64
	deleteWithInstance bool
	kmsKeyId           string
	encrypt            bool
}

func getDataDiskHash(obj diskHash) string {
	h := sha256.New()
	h.Write([]byte(obj.diskType))
	h.Write([]byte(fmt.Sprintf("%d", obj.diskSize)))
	h.Write([]byte(fmt.Sprintf("%t", obj.deleteWithInstance)))
	h.Write([]byte(obj.kmsKeyId))
	h.Write([]byte(fmt.Sprintf("%t", obj.encrypt)))
	return hex.EncodeToString(h.Sum(nil))
}

func sortDataDisks(tmpDataDisks []interface{}, dataDiskList []map[string]interface{}) (sortedList []map[string]interface{}, err error) {
	// import
	if len(tmpDataDisks) == 0 {
		return dataDiskList, nil
	}

	if len(tmpDataDisks) != len(dataDiskList) {
		err = fmt.Errorf("Inconsistent number of data disks.")
		return
	}

	remainingDisks := make([]map[string]interface{}, len(dataDiskList))
	copy(remainingDisks, dataDiskList)

	for _, tmpDisk := range tmpDataDisks {
		dMap := tmpDisk.(map[string]interface{})
		tmpName, _ := dMap["data_disk_name"].(string)
		tmpSizeRaw := dMap["data_disk_size"]
		tmpSize, e := extractInt(tmpSizeRaw)
		if e != nil {
			return nil, e
		}

		tmpType, _ := dMap["data_disk_type"].(string)
		tmpKmsKeyId, _ := dMap["kms_key_id"].(string)
		tmpEncrypt, _ := dMap["encrypt"].(bool)
		tmpDelWithIns, _ := dMap["delete_with_instance"].(bool)
		tmpTpRaw := dMap["throughput_performance"]
		tmpTp, e := extractInt(tmpTpRaw)
		if e != nil {
			return nil, e
		}

		var matchedDisk map[string]interface{}
		matchedIndex := -1

		for i, dataDisk := range remainingDisks {
			dataName, _ := dataDisk["data_disk_name"].(*string)
			dataSizeRaw := dataDisk["data_disk_size"]
			dataSize, e := extractInt(dataSizeRaw)
			if e != nil {
				return nil, e
			}

			dataType, _ := dataDisk["data_disk_type"].(*string)
			dataKmsKeyId, _ := dataDisk["kms_key_id"].(*string)
			dataEncrypt, _ := dataDisk["encrypt"].(*bool)
			dataDelWithIns, _ := dataDisk["delete_with_instance"].(*bool)
			dataTpRaw := dataDisk["throughput_performance"]
			dataTp, e := extractInt(dataTpRaw)
			if e != nil {
				return nil, e
			}

			match := true
			if tmpName != "" && *dataName != tmpName {
				match = false
			}

			if tmpKmsKeyId != "" && dataKmsKeyId != nil {
				if tmpKmsKeyId != *dataKmsKeyId {
					match = false
				}
			}

			if dataSize != tmpSize || *dataType != tmpType || *dataEncrypt != tmpEncrypt || *dataDelWithIns != tmpDelWithIns || dataTp != tmpTp {
				match = false
			}

			if match {
				matchedDisk = dataDisk
				matchedIndex = i
				break
			}
		}

		if matchedIndex == -1 {
			err = fmt.Errorf("Unable to find match: tmpDisk = %v", tmpDisk)
			return
		}

		sortedList = append(sortedList, matchedDisk)
		remainingDisks = append(remainingDisks[:matchedIndex], remainingDisks[matchedIndex+1:]...)
	}

	return
}

func extractInt(value interface{}) (int, error) {
	if value == nil {
		return 0, fmt.Errorf("value is nil.")
	}

	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		ptrValue := reflect.ValueOf(value).Elem().Interface()
		return extractInt(ptrValue)
	}

	switch v := value.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case uint64:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("Unrecognized numerical type: %T.", value)
	}
}
