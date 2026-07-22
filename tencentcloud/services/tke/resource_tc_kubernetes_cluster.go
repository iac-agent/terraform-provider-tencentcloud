package tke

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svcas "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/as"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
)

func ResourceTencentCloudKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKubernetesClusterCreate,
		Read:   resourceTencentCloudKubernetesClusterRead,
		Update: resourceTencentCloudKubernetesClusterUpdate,
		Delete: resourceTencentCloudKubernetesClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: customResourceImporter,
		},
		CustomizeDiff: customdiff.All(
			customizeDiffForContainerRuntimeDefault,
		),
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 集群。",
			},

			"cluster_desc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "描述 集群。",
			},

			"cluster_os": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "tlinux2.4x86_64",
				Description: "Cluster operating 系统，支持 setting 公有 images ( 字段 passes corresponding 镜像 名称) 和 自定义 images ( 字段 passes corresponding 镜像 ID). For details，please refer 到: https://云.tencent.com/document/product/457/68289。",
			},

			"cluster_subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Control Plane Subnet Information. 此 字段 为必填项 仅 在 following scenarios: 当 容器 网络 插件 是 CiliumOverlay，TKE 将 obtain 2 IPs 从 此 子网 到 create 内部 load balancer; 当 creating managed 集群 该 支持 CDC 使用 VPC-CNI 网络 插件，在 least 12 IPs 必须 是 reserved。",
			},

			"cluster_os_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "GENERAL",
				Description:  "Image 类型 集群 os， 可用 值 include: 'GENERAL'. 默认为 'GENERAL'。",
				ValidateFunc: tccommon.ValidateAllowedStringValue(TKE_CLUSTER_OS_TYPES),
			},

			"container_runtime": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "Runtime 类型 集群， 可用 值 include: 'docker' 和 'containerd'. Kubernetes v1.24 has removed dockershim，so please 使用 containerd 在 v1.24 或 higher. 默认值为 `docker` 对于 versions below v1.24 和 `containerd` 对于 versions above v1.24。",
				ValidateFunc: tccommon.ValidateAllowedStringValue(TKE_RUNTIMES),
			},

			"cluster_deploy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "MANAGED_CLUSTER",
				Description:  "Deployment 类型 集群， 可用 值 include: 'MANAGED_CLUSTER' 和 'INDEPENDENT_CLUSTER'. 默认为 'MANAGED_CLUSTER'。",
				ValidateFunc: tccommon.ValidateAllowedStringValue(TKE_DEPLOY_TYPES),
			},

			"cluster_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "版本 的 集群. Use `tencentcloud_kubernetes_available_cluster_versions` 到 get upgradable 集群 版本",
			},

			"upgrade_instances_follow_cluster": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "表示是否upgrade all 集群 实例. 默认为 false。",
			},

			"cluster_ipvs": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     true,
				Description: "表示是否`ipvs` 是 已启用 默认为 true. False 表示 `iptables` 是 已启用",
			},

			"cluster_as_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "表示是否enable 集群 节点 auto scaling. 默认为 false。",
				Deprecated:  "This argument is deprecated because the TKE auto-scaling group was no longer available.",
			},

			"cluster_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "指定cluster 级别，有效 对于 managed 集群，使用 数据 来源 `tencentcloud_kubernetes_cluster_levels` 到 查询 可用 levels. Available 值 examples `L5`，`L20`，`L50`，`L100`，etc。",
			},

			"auto_upgrade_cluster_level": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "是否cluster 级别 auto upgraded，有效 对于 managed 集群。",
			},

			"acquire_cluster_admin_role": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "如果 集合 到 true，它 将 acquire ClusterRole tke:admin. NOTE: 此 arguments 不能 revoke 到 `false` after acquired。",
			},

			"node_pool_global_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Global 配置 effective 对于 all 节点 pools。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_scale_in_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "表示是否enable scale-在。",
						},
						"expander": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "表示which scale-out 方法 将 是 使用 当 there 是 多个 scaling groups. 有效值：`random` - select random scaling 组，`most-pods` - select scaling 组 该 可以 调度 most pods，`least-waste` - select scaling 组 该 可以 ensure fewest remaining resources after Pod scheduling。",
						},
						"max_concurrent_scale_in": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Max concurrent scale-在 卷。",
						},
						"scale_in_delay": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "数量 minutes after 集群 scale-out 当 系统 starts judging 是否perform scale-在。",
						},
						"scale_in_unneeded_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "数量 consecutive minutes 的 idleness after 其中 节点 是 subject 到 scale-在。",
						},
						"scale_in_utilization_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Percentage 的 节点 资源 usage below 其中 节点 是 considered 到 是 idle。",
						},
						"ignore_daemon_sets_utilization": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "是否ignore DaemonSet pods 通过 默认值 当 calculating 资源 usage。",
						},
						"skip_nodes_with_local_storage": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "During scale-在，ignore nodes 使用 本地 存储 pods。",
						},
						"skip_nodes_with_system_pods": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "During scale-在，ignore nodes 使用 pods 在 kube-系统 命名空间 该 是 不 managed 通过 DaemonSet。",
						},
					},
				},
			},

			"cluster_extra_args": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Customized 参数 对于 master 组件,such 作为 kube-apiserver，kube-controller-manager，kube-scheduler。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kube_apiserver": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "customized 参数 对于 kube-apiserver。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"kube_controller_manager": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "customized 参数 对于 kube-controller-manager。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"kube_scheduler": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "customized 参数 对于 kube-scheduler。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"is_dual_stack": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "In VPC-CNI 模式 的 集群， dual stack 集群 状态 默认为 false，indicating non dual stack 集群。",
			},

			"node_name_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "lan-ip",
				Description:  "节点名称 类型 Cluster， 可用 值 include: 'lan-ip' 和 'hostname'，默认为 'lan-ip'。",
				ValidateFunc: tccommon.ValidateAllowedStringValue(TKE_CLUSTER_NODE_NAME_TYPE),
			},

			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "GR",
				Description:  "Cluster 网络 类型， 可用 值 include: 'GR' 和 'VPC-CNI' 和 'CiliumOverlay'. 默认为 GR。",
				ValidateFunc: tccommon.ValidateAllowedStringValue(TKE_CLUSTER_NETWORK_TYPE),
			},

			"enable_customized_pod_cidr": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "是否enable 自定义 模式 的 节点 podCIDR 大小. 默认为 false。",
			},

			"base_pod_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "数量 basic pods. 有效 当 enable_customized_pod_cidr=true。",
			},

			"is_non_static_ip_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "表示是否non-静态 ip 模式 是 已启用 默认为 false。",
			},

			"data_plane_v2": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "是否enable DataPlaneV2 (replace kube-proxy 使用 cilium). `data_plane_v2` 和 `cluster_ipvs` should 不 是 集合 在 same 时间。",
			},

			"deletion_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "表示是否cluster 删除保护 是 已启用 默认为 false。",
			},

			"resource_delete_options": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "资源 deletion 策略 当 集群 是 删除. Currently，CBS 是 支持 (CBS 是 retained 通过 默认值). Only 有效 当 deleting 集群。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "资源类型，有效 值 是 `CBS`，`CLB`，和 `CVM`。",
						},
						"delete_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "deletion 模式 的 CBS resources 当 集群 是 删除，`terminate` (destroy)，`retain` (retain). Other resources 是 删除 通过 默认值。",
						},
						"skip_deletion_protection": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "是否skip resources 使用 删除保护 已启用， 默认为 false。",
						},
					},
				},
			},

			"kube_proxy_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Cluster kube-proxy 模式， 可用 值 include: 'kube-proxy-bpf'. 默认为 不 集合.当 集合 到 kube-proxy-bpf，集群 版本 greater 比 1.14 和 使用 Tencent Linux 2.4 为必填项。",
			},

			"vpc_cni_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Distinguish between shared 网络 card multi-IP 模式 和 independent 网络 card 模式 Fill 在 `tke-路由-eni` 对于 shared 网络 card multi-IP 模式 和 `tke-direct-eni` 对于 independent 网络 card 模式 默认为 shared 网络 card 模式 当 它 是 necessary 到 turn 关闭 vpc-cni 容器 网络 capability，both `eni_subnet_ids` 和 `vpc_cni_type` 必须 是 集合 到 空。",
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"tke-route-eni", "tke-direct-eni"}),
			},

			"vpc_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "私有网络 ID 集群。",
				ValidateFunc: tccommon.ValidateStringLengthInRange(4, 100),
			},

			"cluster_internet": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Open internet 访问 或 不. 如果 此 字段 是 集合 'true'， 字段 below `worker_config` 必须 是 集合. Because 仅 集群 使用 节点 是 allowed 启用 访问 端点. You 可能 open 它 through `tencentcloud_kubernetes_cluster_endpoint`。",
			},

			"cluster_internet_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "域名 名称 对于 集群 Kube-apiserver internet 访问. Be careful 如果 您 modify 值 的 此 参数， cluster_external_endpoint 值 可能 是 changed automatically too。",
			},

			"cluster_intranet": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Open intranet 访问 或 不. 如果 此 字段 是 集合 'true'， 字段 below `worker_config` 必须 是 集合. Because 仅 集群 使用 节点 是 allowed 启用 访问 端点. You 可能 open 它 through `tencentcloud_kubernetes_cluster_endpoint`。",
			},

			"cluster_intranet_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "域名 名称 对于 集群 Kube-apiserver intranet 访问. Be careful 如果 您 modify 值 的 此 参数， pgw_endpoint 值 可能 是 changed automatically too。",
			},

			"cluster_internet_security_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "指定security 组，NOTE: 此 argument 必须 不 是 空 如果 集群 internet 已启用",
			},

			"managed_cluster_internet_security_policies": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Security policies 对于 managed 集群 internet，like:'192.168.1.0/24' 或 '113.116.51.27'，'0.0.0.0/0' 表示 all. 此 字段 可以 仅 集合 当 字段 `cluster_deploy_type` 是 'MANAGED_CLUSTER' 和 `cluster_internet` 是 true. `managed_cluster_internet_security_policies` 可以 不 delete 或 空 once 是 集合。",
				Deprecated:  "this argument was deprecated, use `cluster_internet_security_group` instead.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"cluster_intranet_subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "子网 ID who 可以 访问 此 independent 集群，此 字段 必须 和 可以 仅 集合 当 `cluster_intranet` 是 true. `cluster_intranet_subnet_id` 可以 不 modify once 是 集合。",
			},

			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "项目 ID，默认值为 0。",
			},

			"cluster_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A 网络 地址 block 的 集群. Different 从 vpc cidr 和 cidr 的 other clusters within 此 vpc. Must 是 在 10./192.168/172.[16-31] segments。",
				// ValidateFunc: clusterCidrValidateFunc,
			},

			"ignore_cluster_cidr_conflict": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "表示是否ignore 集群 cidr conflict 错误 默认为 false。",
			},

			"ignore_service_cidr_conflict": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "表示是否ignore 服务 cidr conflict 错误 Only 有效 在 `VPC-CNI` 模式",
			},

			"cluster_max_pod_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     256,
				Description: "最大Pods per 节点 在 集群. 默认为 256. 最小 值 是 4. 当 its power unequal 到 2，它 将 round upward 到 closest power 的 2。",
			},

			"cluster_max_service_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     256,
				Description: "最大services 在 集群. 默认为 256. 范围 是 从 32 到 32768. 当 its power unequal 到 2，它 将 round upward 到 closest power 的 2。",
			},

			"service_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A 网络 地址 block 的 服务. Different 从 vpc cidr 和 cidr 的 other clusters within 此 vpc. Must 是 在 10./192.168/172.[16-31] segments。",
				// ValidateFunc: serviceCidrValidateFunc,
			},

			"eni_subnet_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Subnet Ids 对于 集群 使用 VPC-CNI 网络 模式 此 字段 可以 仅 集合 当 字段 `network_type` 是 'VPC-CNI'. `eni_subnet_ids` 可以 不 空 once 是 集合。",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"claim_expired_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "Claim expired 秒 到 recycle ENI. 此 字段 可以 仅 集合 当 字段 `network_type` 是 'VPC-CNI'. `claim_expired_seconds` 必须 greater 或 equal 比 300 和 less 比 15768000。",
				ValidateFunc: claimExpiredSecondsValidateFunc,
			},

			"master_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Deploy machine 配置 信息 的 'MASTER_ETCD' 服务，和 create <=7 units 对于 common users。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Default:     1,
							Description: "数量 cvm。",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "表示which availability 可用区 将 是 使用。",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Default:     "sub machine of tke",
							Description: "名称 CVMs。",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Specified types 的 CVM 实例。",
						},
						"instance_charge_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "POSTPAID_BY_HOUR",
							Description:  "charge 类型 实例. 有效 值 是 `PREPAID` 和 `POSTPAID_BY_HOUR`. 默认为 `POSTPAID_BY_HOUR`. 注意: TencentCloud International 仅 支持 `POSTPAID_BY_HOUR`，`PREPAID` 实例 将 不 terminated after 集群 删除，和 可能 不 allow 到 delete before expired。",
							ValidateFunc: tccommon.ValidateAllowedStringValue(TKE_INSTANCE_CHARGE_TYPE),
						},
						"instance_charge_type_prepaid_period": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							Default:      1,
							Description:  "tenancy (时间 单位 是 month) 的 prepaid 实例. NOTE: 它 仅 works 当 instance_charge_type 是 集合 到 `PREPAID`. 有效 值 是 `1`，`2`，`3`，`4`，`5`，`6`，`7`，`8`，`9`，`10`，`11`，`12`，`24`，`36`。",
							ValidateFunc: tccommon.ValidateAllowedIntValue(svccvm.CVM_PREPAID_PERIOD),
						},
						"instance_charge_type_prepaid_renew_flag": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							Description:  "自动续费标识 有效值：`NOTIFY_AND_AUTO_RENEW`: notify upon expiration 和 renew automatically，`NOTIFY_AND_MANUAL_RENEW`: notify upon expiration 但 do 不 renew automatically，`DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. 默认值：`NOTIFY_AND_MANUAL_RENEW`. 如果 此 参数 是 指定 作为 `NOTIFY_AND_AUTO_RENEW`， 实例 将 是 automatically renewed 在 monthly basis 如果 账号 balance 是 sufficient. NOTE: 它 仅 works 当 instance_charge_type 是 集合 到 `PREPAID`。",
							ValidateFunc: tccommon.ValidateAllowedStringValue(svccvm.CVM_PREPAID_RENEW_FLAG),
						},
						"subnet_id": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							Description:  "Private 网络 ID。",
							ValidateFunc: tccommon.ValidateStringLengthInRange(4, 100),
						},
						"system_disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "CLOUD_PREMIUM",
							Description:  "System 磁盘 类型 For more 信息 在 limits 的 系统 磁盘 types，see [Storage Overview](https://intl.云.tencent.com/document/product/213/4952). 有效值：`LOCAL_BASIC`: 本地 磁盘，`LOCAL_SSD`: 本地 SSD 磁盘，`CLOUD_SSD`: SSD，`CLOUD_PREMIUM`: Premium Cloud Storage. NOTE: `CLOUD_BASIC`，`LOCAL_BASIC` 和 `LOCAL_SSD` 是 已弃用",
							ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.SYSTEM_DISK_ALLOW_TYPE),
						},
						"system_disk_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							Default:      50,
							Description:  "Volume 的 系统 磁盘 （GB）。 默认为 `50`。",
							ValidateFunc: tccommon.ValidateIntegerInRange(20, 1024),
						},
						"data_disk": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    11,
							Description: "数据盘配置",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ForceNew:     true,
										Default:      "CLOUD_PREMIUM",
										Description:  "Types 的 磁盘，可用值：`CLOUD_PREMIUM` 和 `CLOUD_SSD` 和 `CLOUD_HSSD` 和 `CLOUD_TSSD`。",
										ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.SYSTEM_DISK_ALLOW_TYPE),
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Optional:    true,
										ForceNew:    true,
										Default:     0,
										Description: "Volume 的 磁盘 （GB）。 默认为 `0`。",
									},
									"snapshot_id": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Data 磁盘 快照 ID。",
									},
									"encrypt": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "表示是否encrypt 数据 磁盘，默认值 `false`。",
									},
									"kms_key_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ID 自定义 CMK 在 格式 的 UUID 或 `kms-abcd1234`. 此 参数 是 用于encrypt 云 disks。",
									},
									"file_system": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "File 系统，e.g. `ext3/ext4/xfs`。",
									},
									"auto_format_and_mount": {
										Type:        schema.TypeBool,
										Optional:    true,
										ForceNew:    true,
										Default:     false,
										Description: "Indicate 是否auto 格式 和 mount 或 不. 默认为 `false`。",
									},
									"mount_target": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "挂载目标",
									},
									"disk_partition": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "名称 device 或 分区 到 mount。",
									},
								},
							},
						},
						"internet_charge_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "TRAFFIC_POSTPAID_BY_HOUR",
							Description:  "Charge types 对于 网络 流量. Available 值 include `TRAFFIC_POSTPAID_BY_HOUR`。",
							ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.INTERNET_CHARGE_ALLOW_TYPE),
						},
						"internet_max_bandwidth_out": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Max 带宽 的 Internet 访问 在 Mbps. 默认为 0。",
						},
						"bandwidth_package_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "带宽 包 ID. 如果 用户 是 standard 用户，then bandwidth_package_id 是 needed，或 默认值 has bandwidth_package_id。",
						},
						"public_ip_assigned": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: "指定是否assign 公网 IP 地址",
						},
						"password": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Sensitive:    true,
							Description:  "密码 到 访问，should 是 集合 如果 `key_ids` 不 集合。",
							ValidateFunc: tccommon.ValidateAsConfigPassword,
						},
						"key_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "ID 列表 keys，should 是 集合 如果 `密码` 不 集合。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"security_group_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "Security groups 到 其中 CVM 实例 belongs。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"enhanced_security_service": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Default:     true,
							Description: "To 指定是否enable 云 安全 服务. 默认为 TRUE。",
						},
						"enhanced_monitor_service": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Default:     true,
							Description: "To 指定是否enable 云 监控 服务. 默认为 TRUE。",
						},
						"user_data": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "ase64-encoded 用户 Data text， 长度 限制 是 16KB。",
						},
						"cam_role_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "被授权访问的 CAM 角色名称",
						},
						"hostname": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "主机 名称 attached 实例. Dot (.) 和 dash (-) 不能 是 使用 作为 first 和 last 字符 的 HostName 和 不能 是 使用 consecutively. Windows 示例: 长度 的 名称 character 是 [2，15]，letters (capitalization 是 不 restricted)，numbers 和 dashes (-) 是 allowed，dots (.) 是 不 支持，和 不 all numbers 是 allowed. Examples 的 other types (Linux，etc.): character 长度 是 [2，60]，和 多个 dots 是 allowed. There 是 segment between dots. Each segment allows letters (使用 无 limitation 在 capitalization)，numbers 和 dashes (-)。",
						},
						"disaster_recover_group_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "Disaster recover groups 到 其中 CVM 实例 belongs. Only support 最大 1。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"img_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "有效 镜像 ID，格式 的 img-xxx. 注意: `img_id` 将 是 replaced 使用 镜像 corresponding 到 TKE `cluster_os`。",
							ValidateFunc: tccommon.ValidateImageID,
						},
						"desired_pod_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Default:     0,
							Description: "Indicate 到 集合 desired pod 数量 在 节点. 有效 当 enable_customized_pod_cidr=true，和 它 override `[globe_]desired_pod_num` 对于 当前 节点. Either all 字段 `desired_pod_num` 或 none。",
						},
						"hpc_cluster_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID 的 cvm hpc 集群。",
						},
					},
				},
			},

			"worker_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Deploy machine 配置 信息 的 'WORKER' 服务，和 create <=20 units 对于 common users. other 'WORK' 服务 是 added 通过 'tencentcloud_kubernetes_scale_worker'。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Default:     1,
							Description: "数量 cvm。",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "表示which availability 可用区 将 是 使用。",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Default:     "sub machine of tke",
							Description: "名称 CVMs。",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Specified types 的 CVM 实例。",
						},
						"instance_charge_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "POSTPAID_BY_HOUR",
							Description:  "charge 类型 实例. 有效 值 是 `PREPAID` 和 `POSTPAID_BY_HOUR`. 默认为 `POSTPAID_BY_HOUR`. 注意: TencentCloud International 仅 支持 `POSTPAID_BY_HOUR`，`PREPAID` 实例 将 不 terminated after 集群 删除，和 可能 不 allow 到 delete before expired。",
							ValidateFunc: tccommon.ValidateAllowedStringValue(TKE_INSTANCE_CHARGE_TYPE),
						},
						"instance_charge_type_prepaid_period": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							Default:      1,
							Description:  "tenancy (时间 单位 是 month) 的 prepaid 实例. NOTE: 它 仅 works 当 instance_charge_type 是 集合 到 `PREPAID`. 有效 值 是 `1`，`2`，`3`，`4`，`5`，`6`，`7`，`8`，`9`，`10`，`11`，`12`，`24`，`36`。",
							ValidateFunc: tccommon.ValidateAllowedIntValue(svccvm.CVM_PREPAID_PERIOD),
						},
						"instance_charge_type_prepaid_renew_flag": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							Description:  "自动续费标识 有效值：`NOTIFY_AND_AUTO_RENEW`: notify upon expiration 和 renew automatically，`NOTIFY_AND_MANUAL_RENEW`: notify upon expiration 但 do 不 renew automatically，`DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. 默认值：`NOTIFY_AND_MANUAL_RENEW`. 如果 此 参数 是 指定 作为 `NOTIFY_AND_AUTO_RENEW`， 实例 将 是 automatically renewed 在 monthly basis 如果 账号 balance 是 sufficient. NOTE: 它 仅 works 当 instance_charge_type 是 集合 到 `PREPAID`。",
							ValidateFunc: tccommon.ValidateAllowedStringValue(svccvm.CVM_PREPAID_RENEW_FLAG),
						},
						"subnet_id": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							Description:  "Private 网络 ID。",
							ValidateFunc: tccommon.ValidateStringLengthInRange(4, 100),
						},
						"system_disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "CLOUD_PREMIUM",
							Description:  "System 磁盘 类型 For more 信息 在 limits 的 系统 磁盘 types，see [Storage Overview](https://intl.云.tencent.com/document/product/213/4952). 有效值：`LOCAL_BASIC`: 本地 磁盘，`LOCAL_SSD`: 本地 SSD 磁盘，`CLOUD_SSD`: SSD，`CLOUD_PREMIUM`: Premium Cloud Storage. NOTE: `CLOUD_BASIC`，`LOCAL_BASIC` 和 `LOCAL_SSD` 是 已弃用",
							ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.SYSTEM_DISK_ALLOW_TYPE),
						},
						"system_disk_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							Default:      50,
							Description:  "Volume 的 系统 磁盘 （GB）。 默认为 `50`。",
							ValidateFunc: tccommon.ValidateIntegerInRange(20, 1024),
						},
						"data_disk": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    11,
							Description: "数据盘配置",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ForceNew:     true,
										Default:      "CLOUD_PREMIUM",
										Description:  "Types 的 磁盘，可用值：`CLOUD_PREMIUM` 和 `CLOUD_SSD` 和 `CLOUD_HSSD` 和 `CLOUD_TSSD`。",
										ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.SYSTEM_DISK_ALLOW_TYPE),
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Optional:    true,
										ForceNew:    true,
										Default:     0,
										Description: "Volume 的 磁盘 （GB）。 默认为 `0`。",
									},
									"snapshot_id": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Data 磁盘 快照 ID。",
									},
									"encrypt": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "表示是否encrypt 数据 磁盘，默认值 `false`。",
									},
									"kms_key_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ID 自定义 CMK 在 格式 的 UUID 或 `kms-abcd1234`. 此 参数 是 用于encrypt 云 disks。",
									},
									"file_system": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "File 系统，e.g. `ext3/ext4/xfs`。",
									},
									"auto_format_and_mount": {
										Type:        schema.TypeBool,
										Optional:    true,
										ForceNew:    true,
										Default:     false,
										Description: "Indicate 是否auto 格式 和 mount 或 不. 默认为 `false`。",
									},
									"mount_target": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "挂载目标",
									},
									"disk_partition": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "名称 device 或 分区 到 mount。",
									},
								},
							},
						},
						"internet_charge_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "TRAFFIC_POSTPAID_BY_HOUR",
							Description:  "Charge types 对于 网络 流量. Available 值 include `TRAFFIC_POSTPAID_BY_HOUR`。",
							ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.INTERNET_CHARGE_ALLOW_TYPE),
						},
						"internet_max_bandwidth_out": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Max 带宽 的 Internet 访问 在 Mbps. 默认为 0。",
						},
						"bandwidth_package_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "带宽 包 ID. 如果 用户 是 standard 用户，then bandwidth_package_id 是 needed，或 默认值 has bandwidth_package_id。",
						},
						"public_ip_assigned": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: "指定是否assign 公网 IP 地址",
						},
						"password": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Sensitive:    true,
							Description:  "密码 到 访问，should 是 集合 如果 `key_ids` 不 集合。",
							ValidateFunc: tccommon.ValidateAsConfigPassword,
						},
						"key_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "ID 列表 keys，should 是 集合 如果 `密码` 不 集合。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"security_group_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "Security groups 到 其中 CVM 实例 belongs。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"enhanced_security_service": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Default:     true,
							Description: "To 指定是否enable 云 安全 服务. 默认为 TRUE。",
						},
						"enhanced_monitor_service": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Default:     true,
							Description: "To 指定是否enable 云 监控 服务. 默认为 TRUE。",
						},
						"user_data": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "ase64-encoded 用户 Data text， 长度 限制 是 16KB。",
						},
						"cam_role_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "被授权访问的 CAM 角色名称",
						},
						"hostname": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "主机 名称 attached 实例. Dot (.) 和 dash (-) 不能 是 使用 作为 first 和 last 字符 的 HostName 和 不能 是 使用 consecutively. Windows 示例: 长度 的 名称 character 是 [2，15]，letters (capitalization 是 不 restricted)，numbers 和 dashes (-) 是 allowed，dots (.) 是 不 支持，和 不 all numbers 是 allowed. Examples 的 other types (Linux，etc.): character 长度 是 [2，60]，和 多个 dots 是 allowed. There 是 segment between dots. Each segment allows letters (使用 无 limitation 在 capitalization)，numbers 和 dashes (-)。",
						},
						"disaster_recover_group_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "Disaster recover groups 到 其中 CVM 实例 belongs. Only support 最大 1。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"img_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "有效 镜像 ID，格式 的 img-xxx. 注意: `img_id` 将 是 replaced 使用 镜像 corresponding 到 TKE `cluster_os`。",
							ValidateFunc: tccommon.ValidateImageID,
						},
						"desired_pod_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Default:     0,
							Description: "Indicate 到 集合 desired pod 数量 在 节点. 有效 当 enable_customized_pod_cidr=true，和 它 override `[globe_]desired_pod_num` 对于 当前 节点. Either all 字段 `desired_pod_num` 或 none。",
						},
						"hpc_cluster_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID 的 cvm hpc 集群。",
						},
					},
				},
			},

			"exist_instance": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Create tke 集群 通过 existed 实例。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_role": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "角色 的 existed 节点. 值: MASTER_ETCD 或 WORKER。",
						},
						"instances_para": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Reinstallation 参数 的 existing 实例。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_ids": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Cluster IDs。",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"security_group_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Security groups 到 其中 CVM 实例 belongs。",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"password": {
										Type:         schema.TypeString,
										Optional:     true,
										Sensitive:    true,
										Description:  "密码 到 访问，should 是 集合 如果 `key_ids` 不 集合。",
										ValidateFunc: tccommon.ValidateAsConfigPassword,
									},
									"key_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "ID 列表 keys，should 是 集合 如果 `密码` 不 集合。",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"enhanced_security_service": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     true,
										Description: "To 指定是否enable 云 安全 服务. 默认为 TRUE。",
									},
									"enhanced_monitor_service": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     true,
										Description: "To 指定是否enable 云 监控 服务. 默认为 TRUE。",
									},
									"master_config": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Advanced Node Settings. commonly 用于attach existing 实例。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mount_target": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "挂载目标 默认为 不 mounting。",
												},
												"docker_graph_path": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Docker graph 路径 默认为 `/var/lib/docker`。",
												},
												"user_script": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "User 脚本 encoded 在 base64, 其中 将 是 executed after k8s 组件 runs. 用户 needs 到 ensure 脚本's reentrant 和 retry logic. 脚本 和 its generated 日志 files 可以 是 viewed 在 节点 路径 /数据/ccs_userscript/. 如果 节点 needs 到 是 initialized before joining 调度, 它 可以 是 使用 在 conjunction 使用 `unschedulable` 参数. After final initialization 的 userScript 是 completed, add command \"kubectl uncordon nodename --kubeconfig=/root/.kube/config\" 到 add 节点 到 调度.",
												},
												"unschedulable": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Set 是否joined nodes participate 在 scheduling，使用 默认值 的 0，indicating participation 在 scheduling; Non 0 表示 不 participating 在 scheduling。",
												},
												"labels": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Node 标签 列表。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "名称 map。",
															},
															"value": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "值 的 map。",
															},
														},
													},
												},
												"data_disk": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "数据盘配置",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"disk_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Types 的 磁盘. 有效 值: `LOCAL_BASIC`，`LOCAL_SSD`，`CLOUD_BASIC`，`CLOUD_PREMIUM`，`CLOUD_SSD`，`CLOUD_HSSD`，`CLOUD_TSSD` 和 `CLOUD_BSSD`。",
															},
															"file_system": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "File 系统，e.g. `ext3/ext4/xfs`。",
															},
															"disk_size": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Volume 的 磁盘 （GB）。 默认为 `0`。",
															},
															"auto_format_and_mount": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Indicate 是否auto 格式 和 mount 或 不. 默认为 `false`。",
															},
															"mount_target": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "挂载目标",
															},
															"disk_partition": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "名称 device 或 分区 到 mount. NOTE: 此 argument doesn't support setting 在 节点 池，或 将 leads 到 mount 错误",
															},
														},
													},
												},
												"extra_args": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Custom 参数 信息 related 到 节点. 此 是 white-列表 参数。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"kubelet": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Kubelet 自定义 参数. 参数 格式 是 [\"k1=v1\", \"k1=v2\"].",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
												"desired_pod_number": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Indicate 到 集合 desired pod 数量 在 节点. 有效 当 集群 是 podCIDR。",
												},
												"gpu_args": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "GPU 驱动 参数。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mig_enable": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "是否enable MIG。",
															},
															"driver": {
																Type:         schema.TypeMap,
																Optional:     true,
																Description:  "GPU 驱动 版本 格式 like: `{ 版本: String，名称: String }`. `版本`: 版本 的 GPU 驱动 或 CUDA; `名称`: 名称 GPU 驱动 或 CUDA。",
																ValidateFunc: tccommon.ValidateTkeGpuDriverVersion,
															},
															"cuda": {
																Type:         schema.TypeMap,
																Optional:     true,
																Description:  "CUDA 版本 格式 like: `{ 版本: String，名称: String }`. `版本`: 版本 的 GPU 驱动 或 CUDA; `名称`: 名称 GPU 驱动 或 CUDA。",
																ValidateFunc: tccommon.ValidateTkeGpuDriverVersion,
															},
															"cudnn": {
																Type:         schema.TypeMap,
																Optional:     true,
																Description:  "cuDNN 版本 格式 like: `{ 版本: String，名称: String，doc_name: String，dev_name: String }`. `版本`: cuDNN 版本; `名称`: cuDNN 名称; `doc_name`: Doc 名称 cuDNN; `dev_name`: Dev 名称 cuDNN。",
																ValidateFunc: tccommon.ValidateTkeGpuDriverVersion,
															},
															"custom_driver": {
																Type:        schema.TypeMap,
																Optional:    true,
																Description: "Custom GPU 驱动. 格式 like: `{地址: String}`. `地址`: URL 的 自定义 GPU 驱动 地址",
															},
														},
													},
												},
												"taints": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Node taint。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "键 的 taint。",
															},
															"value": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "值 的 taint。",
															},
															"effect": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Effect 的 taint。",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"desired_pod_numbers": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Custom 模式 集群，您 可以 指定number 的 pods 对于 each 节点. corresponding 到 existed_instances_para.instance_ids 参数。",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},

			"auth_options": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "指定cluster authentication 配置. Only 可用 对于 managed 集群 和 `cluster_version` >= 1.20。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"use_tke_default": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "如果 集合 到 `true`， issuer 和 jwks_uri 将 是 generated automatically 通过 tke，please do 不 集合 issuer 和 jwks_uri，和 they 将 是 ignored。",
						},
						"jwks_uri": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "指定service-账号-jwks-uri. 如果 use_tke_默认为 集合 到 `true`，please do 不 集合 此 字段，它 将 是 ignored anyway。",
						},
						"issuer": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "指定service-账号-issuer. 如果 use_tke_默认为 集合 到 `true`，please do 不 集合 此 字段，它 将 是 ignored anyway。",
						},
						"auto_create_discovery_anonymous_auth": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "如果 集合 到 `true`， rbac 规则 将 是 创建 automatically 其中 allow anonymous 用户 到 访问 '/.well-known/openid-配置' 和 '/openid/v1/jwks'。",
						},
					},
				},
			},

			"extension_addon": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Information 的 add-在 到 是 installed. It 是 recommended 到 使用 资源 `tencentcloud_kubernetes_addon` management 集群 addon。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Add-在 名称",
						},
						"param": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Parameter 的 add-在 资源 对象 在 JSON 字符串 格式，please check 示例 在 top 的 页面 对于 reference。",
							DiffSuppressFunc: helper.DiffSupressJSON,
						},
					},
				},
			},

			"log_agent": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "指定cluster 日志 agent 配置",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "是否log agent 已启用",
						},
						"kubelet_root_dir": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Kubelet root directory 作为 literal。",
						},
					},
				},
			},

			"event_persistence": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "指定cluster Event Persistence 配置 NOTE: Please make sure your TKE CamRole have 权限 到 访问 CLS 服务。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "指定weather Event Persistence 已启用",
						},
						"log_set_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "指定id 的 existing CLS 日志 集合，或 auto create new 集合 通过 leave 它 空。",
						},
						"topic_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "指定id 的 existing CLS 日志 主题，或 auto create new 主题 通过 leave 它 空。",
						},
						"delete_event_log_and_topic": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "当 您 want 到 close 集群 事件 persistence 或 delete 集群，您 可以 使用 此 参数 到 determine 是否event persistence 日志 集合 和 主题 创建 通过 默认值 将 是 删除。",
						},
					},
				},
			},

			"cluster_audit": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "指定Cluster Audit 配置 NOTE: Please make sure your TKE CamRole have 权限 到 访问 CLS 服务。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "指定weather Cluster Audit 已启用 NOTE: Enable Cluster Audit 将 also auto install Log Agent。",
						},
						"log_set_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "指定id 的 existing CLS 日志 集合，或 auto create new 集合 通过 leave 它 空。",
						},
						"topic_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "指定id 的 existing CLS 日志 主题，或 auto create new 主题 通过 leave 它 空。",
						},
						"delete_audit_log_and_topic": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "当 您 want 到 close 集群 审计日志 或 delete 集群，您 可以 使用 此 参数 到 determine 是否audit 日志 集合 和 主题 创建 通过 默认值 将 是 删除。",
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 的 集群。",
			},

			"cluster_node_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "节点数量 在 集群。",
			},

			"worker_instances_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An 信息 列表 cvm within 'WORKER' clusters. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID cvm。",
						},
						"instance_role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "角色 的 cvm。",
						},
						"instance_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State 的 cvm。",
						},
						"failed_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Information 的 cvm 当 它 是 failed。",
						},
						"lan_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "LAN IP 的 cvm。",
						},
					},
				},
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Labels 的 tke 集群 nodes。",
			},

			"unschedulable": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				Default:          0,
				Description:      "Sets 是否joining 节点 participates 在 调度. 默认为 '0'. Participate 在 scheduling。",
				DiffSuppressFunc: unschedulableDiffSuppressFunc,
			},

			"mount_target": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "挂载目标 默认为 不 mounting。",
			},

			"globe_desired_pod_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Indicate 到 集合 desired pod 数量 在 节点. 有效 当 enable_customized_pod_cidr=true，和 它 takes effect 对于 all nodes。",
			},

			"docker_graph_path": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          "/var/lib/docker",
				Description:      "Docker graph 路径 默认为 `/var/lib/docker`。",
				DiffSuppressFunc: dockerGraphPathDiffSuppressFunc,
			},

			"pre_start_user_script": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Base64-encoded 用户 脚本，executed before initializing 节点，currently 仅 effective 对于 adding existing nodes。",
			},

			"extra_args": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Custom 参数 信息 related 到 节点。",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"runtime_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Container Runtime 版本",
			},

			"kube_config": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Kubernetes 配置",
			},

			"kube_config_intranet": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Kubernetes 配置 的 私有 网络。",
			},

			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "用户 名称 账号",
			},

			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "密码 的 账号",
			},

			"certification_authority": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "证书 用于access。",
			},

			"cluster_external_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "External 网络 地址 到 访问。",
			},

			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "域名 名称 对于 访问。",
			},

			"pgw_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Intranet 地址 用于access。",
			},

			"security_policy": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Access 策略。",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"cdc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CDC ID。",
			},

			"instance_delete_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "strategy 对于 deleting 集群 实例: terminate (destroy 实例，仅 support pay 作为 您 go 云 主机 实例) retain (remove 仅，keep 实例)，默认为 terminate。",
			},

			"disable_addons": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "To prevent installation 的 特定 Addon 组件，enter corresponding AddonName。",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTencentCloudKubernetesClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		clusterId string
	)
	var (
		request  = tkev20180525.NewCreateClusterRequest()
		response = tkev20180525.NewCreateClusterResponse()
	)

	if v, ok := d.GetOk("cdc_id"); ok {
		request.CdcId = helper.String(v.(string))
	}

	clusterCIDRSettings := tkev20180525.ClusterCIDRSettings{}
	if v, ok := d.GetOk("cluster_cidr"); ok {
		clusterCIDRSettings.ClusterCIDR = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("ignore_cluster_cidr_conflict"); ok {
		clusterCIDRSettings.IgnoreClusterCIDRConflict = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("ignore_service_cidr_conflict"); ok {
		clusterCIDRSettings.IgnoreServiceCIDRConflict = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("cluster_max_service_num"); ok {
		clusterCIDRSettings.MaxClusterServiceNum = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOkExists("cluster_max_pod_num"); ok {
		clusterCIDRSettings.MaxNodePodNum = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("service_cidr"); ok {
		clusterCIDRSettings.ServiceCIDR = helper.String(v.(string))
	}
	request.ClusterCIDRSettings = &clusterCIDRSettings

	clusterBasicSettings := tkev20180525.ClusterBasicSettings{}
	if v, ok := d.GetOk("cluster_version"); ok {
		clusterBasicSettings.ClusterVersion = helper.String(v.(string))
	}
	if v, ok := d.GetOk("cluster_name"); ok {
		clusterBasicSettings.ClusterName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("cluster_desc"); ok {
		clusterBasicSettings.ClusterDescription = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("project_id"); ok {
		clusterBasicSettings.ProjectId = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("cluster_os_type"); ok {
		clusterBasicSettings.OsCustomizeType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("cluster_subnet_id"); ok {
		clusterBasicSettings.SubnetId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("cluster_level"); ok {
		clusterBasicSettings.ClusterLevel = helper.String(v.(string))
	}
	autoUpgradeClusterLevel := tkev20180525.AutoUpgradeClusterLevel{}
	if v, ok := d.GetOkExists("auto_upgrade_cluster_level"); ok {
		autoUpgradeClusterLevel.IsAutoUpgrade = helper.Bool(v.(bool))
	}
	clusterBasicSettings.AutoUpgradeClusterLevel = &autoUpgradeClusterLevel
	request.ClusterBasicSettings = &clusterBasicSettings

	clusterAdvancedSettings := tkev20180525.ClusterAdvancedSettings{}
	if v, ok := d.GetOkExists("cluster_ipvs"); ok {
		clusterAdvancedSettings.IPVS = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("cluster_as_enabled"); ok {
		clusterAdvancedSettings.AsEnabled = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOk("container_runtime"); ok {
		clusterAdvancedSettings.ContainerRuntime = helper.String(v.(string))
	}
	if v, ok := d.GetOk("node_name_type"); ok {
		clusterAdvancedSettings.NodeNameType = helper.String(v.(string))
	}
	if extraArgsMap, ok := helper.InterfacesHeadMap(d, "cluster_extra_args"); ok {
		clusterExtraArgs := tkev20180525.ClusterExtraArgs{}
		if v, ok := extraArgsMap["kube_apiserver"]; ok {
			kubeAPIServerSet := v.([]interface{})
			for i := range kubeAPIServerSet {
				if kubeAPIServer, ok := kubeAPIServerSet[i].(string); ok && kubeAPIServer != "" {
					clusterExtraArgs.KubeAPIServer = append(clusterExtraArgs.KubeAPIServer, helper.String(kubeAPIServer))
				}
			}
		}
		if v, ok := extraArgsMap["kube_controller_manager"]; ok {
			kubeControllerManagerSet := v.([]interface{})
			for i := range kubeControllerManagerSet {
				if kubeControllerManager, ok := kubeControllerManagerSet[i].(string); ok && kubeControllerManager != "" {
					clusterExtraArgs.KubeControllerManager = append(clusterExtraArgs.KubeControllerManager, helper.String(kubeControllerManager))
				}
			}
		}
		if v, ok := extraArgsMap["kube_scheduler"]; ok {
			kubeSchedulerSet := v.([]interface{})
			for i := range kubeSchedulerSet {
				if kubeScheduler, ok := kubeSchedulerSet[i].(string); ok && kubeScheduler != "" {
					clusterExtraArgs.KubeScheduler = append(clusterExtraArgs.KubeScheduler, helper.String(kubeScheduler))
				}
			}
		}
		clusterAdvancedSettings.ExtraArgs = &clusterExtraArgs
	}
	if v, ok := d.GetOkExists("is_dual_stack"); ok {
		clusterAdvancedSettings.IsDualStack = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOk("network_type"); ok {
		clusterAdvancedSettings.NetworkType = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("is_non_static_ip_mode"); ok {
		clusterAdvancedSettings.IsNonStaticIpMode = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("data_plane_v2"); ok {
		clusterAdvancedSettings.DataPlaneV2 = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("deletion_protection"); ok {
		clusterAdvancedSettings.DeletionProtection = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOk("kube_proxy_mode"); ok {
		clusterAdvancedSettings.KubeProxyMode = helper.String(v.(string))
	}
	if v, ok := d.GetOk("runtime_version"); ok {
		clusterAdvancedSettings.RuntimeVersion = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("enable_customized_pod_cidr"); ok {
		clusterAdvancedSettings.EnableCustomizedPodCIDR = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("base_pod_num"); ok {
		clusterAdvancedSettings.BasePodNumber = helper.IntInt64(v.(int))
	}
	request.ClusterAdvancedSettings = &clusterAdvancedSettings

	instanceAdvancedSettings := tkev20180525.InstanceAdvancedSettings{}
	if v, ok := d.GetOkExists("globe_desired_pod_num"); ok {
		instanceAdvancedSettings.DesiredPodNumber = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("mount_target"); ok {
		instanceAdvancedSettings.MountTarget = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("unschedulable"); ok {
		instanceAdvancedSettings.Unschedulable = helper.IntInt64(v.(int))
	}
	request.InstanceAdvancedSettings = &instanceAdvancedSettings

	if v, ok := d.GetOk("extension_addon"); ok {
		for _, item := range v.([]interface{}) {
			extensionAddonsMap := item.(map[string]interface{})
			extensionAddon := tkev20180525.ExtensionAddon{}
			if v, ok := extensionAddonsMap["name"]; ok {
				extensionAddon.AddonName = helper.String(v.(string))
			}
			if v, ok := extensionAddonsMap["param"]; ok {
				extensionAddon.AddonParam = helper.String(v.(string))
			}
			request.ExtensionAddons = append(request.ExtensionAddons, &extensionAddon)
		}
	}

	if v, ok := d.GetOk("disable_addons"); ok {
		for _, item := range v.([]interface{}) {
			if disableAddon, ok := item.(string); ok {
				request.DisableAddons = append(request.DisableAddons, &disableAddon)
			}
		}
	}

	if err := resourceTencentCloudKubernetesClusterCreatePostFillRequest0(ctx, request); err != nil {
		return err
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().CreateClusterWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create kubernetes cluster failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId

	if err := resourceTencentCloudKubernetesClusterCreatePostHandleResponse0(ctx, response); err != nil {
		return err
	}

	d.SetId(clusterId)

	return resourceTencentCloudKubernetesClusterRead(d, meta)
}

func resourceTencentCloudKubernetesClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	clusterId := d.Id()

	respData, err := service.DescribeKubernetesClusterById(ctx, clusterId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `kubernetes_cluster` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData.CdcId != nil {
		_ = d.Set("cdc_id", respData.CdcId)
	}

	if respData.ClusterName != nil {
		_ = d.Set("cluster_name", respData.ClusterName)
	}

	if respData.ClusterDescription != nil {
		_ = d.Set("cluster_desc", respData.ClusterDescription)
	}

	if respData.ClusterVersion != nil {
		_ = d.Set("cluster_version", respData.ClusterVersion)
	}

	if respData.ClusterType != nil {
		_ = d.Set("cluster_deploy_type", respData.ClusterType)
	}

	if respData.ClusterNetworkSettings != nil {
		if respData.ClusterNetworkSettings.ClusterCIDR != nil {
			_ = d.Set("cluster_cidr", respData.ClusterNetworkSettings.ClusterCIDR)
		}

		if respData.ClusterNetworkSettings.IgnoreClusterCIDRConflict != nil {
			_ = d.Set("ignore_cluster_cidr_conflict", respData.ClusterNetworkSettings.IgnoreClusterCIDRConflict)
		}

		if respData.ClusterNetworkSettings.IgnoreServiceCIDRConflict != nil {
			_ = d.Set("ignore_service_cidr_conflict", respData.ClusterNetworkSettings.IgnoreServiceCIDRConflict)
		}

		if respData.ClusterNetworkSettings.MaxNodePodNum != nil {
			_ = d.Set("cluster_max_pod_num", respData.ClusterNetworkSettings.MaxNodePodNum)
		}

		if respData.ClusterNetworkSettings.MaxClusterServiceNum != nil {
			_ = d.Set("cluster_max_service_num", respData.ClusterNetworkSettings.MaxClusterServiceNum)
		}

		if respData.ClusterNetworkSettings.Ipvs != nil {
			_ = d.Set("cluster_ipvs", respData.ClusterNetworkSettings.Ipvs)
		}

		if respData.ClusterNetworkSettings.VpcId != nil {
			_ = d.Set("vpc_id", respData.ClusterNetworkSettings.VpcId)
		}

		if respData.ClusterNetworkSettings.Subnets != nil {
			_ = d.Set("eni_subnet_ids", respData.ClusterNetworkSettings.Subnets)
		}

		if respData.ClusterNetworkSettings.DataPlaneV2 != nil {
			_ = d.Set("data_plane_v2", respData.ClusterNetworkSettings.DataPlaneV2)
		}
	}

	if respData.ClusterNodeNum != nil {
		_ = d.Set("cluster_node_num", respData.ClusterNodeNum)
	}

	if respData.ProjectId != nil {
		_ = d.Set("project_id", respData.ProjectId)
	}

	if respData.DeletionProtection != nil {
		_ = d.Set("deletion_protection", respData.DeletionProtection)
	}

	if respData.ClusterLevel != nil {
		_ = d.Set("cluster_level", respData.ClusterLevel)
	}

	if err := resourceTencentCloudKubernetesClusterReadPostHandleResponse0(ctx, respData); err != nil {
		return err
	}

	var respData1 *tkev20180525.DescribeClusterInstancesResponseParams
	reqErr1 := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKubernetesClusterById1(ctx, clusterId)
		if e != nil {
			if err := resourceTencentCloudKubernetesClusterReadRequestOnError1(ctx, result, e); err != nil {
				return err
			}
			return tccommon.RetryError(e)
		}
		respData1 = result
		return nil
	})
	if reqErr1 != nil {
		log.Printf("[CRITAL]%s read kubernetes cluster failed, reason:%+v", logId, reqErr1)
		return reqErr1
	}

	if respData1 == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `kubernetes_cluster` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	instanceSetList := make([]map[string]interface{}, 0, len(respData1.InstanceSet))
	if respData1.InstanceSet != nil {
		for _, instanceSet := range respData1.InstanceSet {
			instanceSetMap := map[string]interface{}{}

			if instanceSet.InstanceId != nil {
				instanceSetMap["instance_id"] = instanceSet.InstanceId
			}

			if instanceSet.InstanceRole != nil {
				instanceSetMap["instance_role"] = instanceSet.InstanceRole
			}

			if instanceSet.InstanceState != nil {
				instanceSetMap["instance_state"] = instanceSet.InstanceState
			}

			if instanceSet.FailedReason != nil {
				instanceSetMap["failed_reason"] = instanceSet.FailedReason
			}

			if instanceSet.LanIP != nil {
				instanceSetMap["lan_ip"] = instanceSet.LanIP
			}

			instanceSetList = append(instanceSetList, instanceSetMap)
		}

		_ = d.Set("worker_instances_list", instanceSetList)
	}

	var respData2 *tkev20180525.DescribeClusterSecurityResponseParams
	reqErr2 := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKubernetesClusterById2(ctx, clusterId)
		if e != nil {
			if err := resourceTencentCloudKubernetesClusterReadRequestOnError2(ctx, result, e); err != nil {
				return err
			}
			return tccommon.RetryError(e)
		}
		respData2 = result
		return nil
	})
	if reqErr2 != nil {
		log.Printf("[CRITAL]%s read kubernetes cluster failed, reason:%+v", logId, reqErr2)
		return reqErr2
	}

	if respData2 == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `kubernetes_cluster` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData2.UserName != nil {
		_ = d.Set("user_name", respData2.UserName)
	}

	if respData2.Password != nil {
		_ = d.Set("password", respData2.Password)
	}

	if respData2.CertificationAuthority != nil {
		_ = d.Set("certification_authority", respData2.CertificationAuthority)
	}

	if respData2.ClusterExternalEndpoint != nil {
		_ = d.Set("cluster_external_endpoint", respData2.ClusterExternalEndpoint)
	}

	if respData2.Domain != nil {
		_ = d.Set("domain", respData2.Domain)
	}

	if respData2.PgwEndpoint != nil {
		_ = d.Set("pgw_endpoint", respData2.PgwEndpoint)
	}

	if err := resourceTencentCloudKubernetesClusterReadPostHandleResponse2(ctx, respData2); err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudKubernetesClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{"cdc_id", "extension_addon", "disable_addons"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	clusterId := d.Id()

	if err := resourceTencentCloudKubernetesClusterUpdateOnStart(ctx); err != nil {
		return err
	}

	needChange := false
	mutableArgs := []string{"project_id", "cluster_name", "cluster_desc", "cluster_level", "auto_upgrade_cluster_level"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := tkev20180525.NewModifyClusterAttributeRequest()

		request.ClusterId = helper.String(clusterId)

		if v, ok := d.GetOkExists("project_id"); ok {
			request.ProjectId = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("cluster_name"); ok {
			request.ClusterName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("cluster_desc"); ok {
			request.ClusterDesc = helper.String(v.(string))
		}

		if err := resourceTencentCloudKubernetesClusterUpdatePostFillRequest0(ctx, request); err != nil {
			return err
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().ModifyClusterAttributeWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update kubernetes cluster failed, reason:%+v", logId, err)
			return err
		}
	}

	needChange1 := false
	mutableArgs1 := []string{"cluster_version"}
	for _, v := range mutableArgs1 {
		if d.HasChange(v) {
			needChange1 = true
			break
		}
	}

	if needChange1 {
		request1 := tkev20180525.NewUpdateClusterVersionRequest()

		response1 := tkev20180525.NewUpdateClusterVersionResponse()

		request1.ClusterId = helper.String(clusterId)

		if v, ok := d.GetOk("cluster_version"); ok {
			request1.DstVersion = helper.String(v.(string))
		}

		if err := resourceTencentCloudKubernetesClusterUpdatePostFillRequest1(ctx, request1); err != nil {
			return err
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().UpdateClusterVersionWithContext(ctx, request1)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request1.GetAction(), request1.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update kubernetes cluster failed, reason:%+v", logId, err)
			return err
		}
		if err := resourceTencentCloudKubernetesClusterUpdatePostHandleResponse1(ctx, response1); err != nil {
			return err
		}

	}

	// upgrade node version(instances)
	if v, ok := d.GetOkExists("upgrade_instances_follow_cluster"); ok {
		if v.(bool) {
			tkeService := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
			err := upgradeClusterInstances(tkeService, ctx, clusterId)
			if err != nil {
				return err
			}
		}
	}

	needChange2 := false
	mutableArgs2 := []string{"node_pool_global_config"}
	for _, v := range mutableArgs2 {
		if d.HasChange(v) {
			needChange2 = true
			break
		}
	}

	if needChange2 {
		request2 := tkev20180525.NewModifyClusterAsGroupOptionAttributeRequest()

		request2.ClusterId = helper.String(clusterId)

		if clusterAsGroupOptionMap, ok := helper.InterfacesHeadMap(d, "node_pool_global_config"); ok {
			clusterAsGroupOption := tkev20180525.ClusterAsGroupOption{}
			if v, ok := clusterAsGroupOptionMap["is_scale_in_enabled"]; ok {
				clusterAsGroupOption.IsScaleDownEnabled = helper.Bool(v.(bool))
			}
			if v, ok := clusterAsGroupOptionMap["expander"]; ok {
				clusterAsGroupOption.Expander = helper.String(v.(string))
			}
			if v, ok := clusterAsGroupOptionMap["max_concurrent_scale_in"]; ok {
				clusterAsGroupOption.MaxEmptyBulkDelete = helper.IntInt64(v.(int))
			}
			if v, ok := clusterAsGroupOptionMap["scale_in_delay"]; ok {
				clusterAsGroupOption.ScaleDownDelay = helper.IntInt64(v.(int))
			}
			if v, ok := clusterAsGroupOptionMap["scale_in_unneeded_time"]; ok {
				clusterAsGroupOption.ScaleDownUnneededTime = helper.IntInt64(v.(int))
			}
			if v, ok := clusterAsGroupOptionMap["scale_in_utilization_threshold"]; ok {
				clusterAsGroupOption.ScaleDownUtilizationThreshold = helper.IntInt64(v.(int))
			}
			if v, ok := clusterAsGroupOptionMap["ignore_daemon_sets_utilization"]; ok {
				clusterAsGroupOption.IgnoreDaemonSetsUtilization = helper.Bool(v.(bool))
			}
			if v, ok := clusterAsGroupOptionMap["skip_nodes_with_local_storage"]; ok {
				clusterAsGroupOption.SkipNodesWithLocalStorage = helper.Bool(v.(bool))
			}
			if v, ok := clusterAsGroupOptionMap["skip_nodes_with_system_pods"]; ok {
				clusterAsGroupOption.SkipNodesWithSystemPods = helper.Bool(v.(bool))
			}
			request2.ClusterAsGroupOption = &clusterAsGroupOption
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().ModifyClusterAsGroupOptionAttributeWithContext(ctx, request2)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request2.GetAction(), request2.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update kubernetes cluster failed, reason:%+v", logId, err)
			return err
		}
	}

	needChange3 := false
	mutableArgs3 := []string{"cluster_os"}
	for _, v := range mutableArgs3 {
		if d.HasChange(v) {
			needChange3 = true
			break
		}
	}

	if needChange3 {
		request3 := tkev20180525.NewModifyClusterImageRequest()

		request3.ClusterId = helper.String(clusterId)

		if v, ok := d.GetOk("cluster_os"); ok {
			request3.ImageId = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().ModifyClusterImageWithContext(ctx, request3)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request3.GetAction(), request3.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update kubernetes cluster failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("exist_instance") {
		tkeService := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		cvmService := svccvm.NewCvmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

		oldInterface, newInterface := d.GetChange("exist_instance")
		oldInstances := oldInterface.(*schema.Set)
		newInstances := newInterface.(*schema.Set)

		remove := oldInstances.Difference(newInstances).List()
		add := newInstances.Difference(oldInstances).List()

		// scale out first
		if len(add) > 0 {
			tmpNew := make([]*tkev20180525.ExistedInstancesForNode, 0, len(add))
			instanceIds := make([]*string, 0)
			instanceInfo := make([]map[string]interface{}, 0)
			for index := range add {
				if add[index] != nil {
					instance := add[index].(map[string]interface{})
					existedInstance, _ := tkeGetCvmExistInstancesPara(instance)
					tmpNew = append(tmpNew, &existedInstance)

					// get all new cvm IDs
					if len(existedInstance.ExistedInstancesPara.InstanceIds) > 0 {
						dMap := make(map[string]interface{}, 0)
						instanceIds = append(instanceIds, existedInstance.ExistedInstancesPara.InstanceIds...)
						dMap["instance_ids"] = instanceIds
						dMap["node_role"] = existedInstance.NodeRole
						instanceInfo = append(instanceInfo, dMap)
					}
				}
			}

			if len(tmpNew) > 0 {
				request := tkev20180525.NewScaleOutClusterMasterRequest()
				request.ClusterId = &clusterId
				request.ExistedInstancesForNode = tmpNew
				err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().ScaleOutClusterMasterWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s scale out cluster failed, reason:%+v", logId, err)
					return err
				}

				// wait for cvm status
				err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
					result, e := cvmService.DescribeInstanceByFilter(ctx, instanceIds, nil)
					if e != nil {
						return tccommon.RetryError(e, tccommon.InternalError)
					}

					initFlag := true
					if result != nil {
						for _, item := range result {
							if item.InstanceState != nil && *item.InstanceState != "RUNNING" {
								initFlag = false
								break
							}
						}

						if initFlag {
							return nil
						}
					}

					return resource.RetryableError(fmt.Errorf("cvm instance status is not RUNNING, retry..."))
				})

				if err != nil {
					return err
				}

				// wait for tke node init
				for _, item := range instanceInfo {
					tmpInsIds := item["instance_ids"].([]*string)
					nodeRole := item["node_role"].(*string)
					err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
						result, e := tkeService.DescribeKubernetesClusterMasterAttachmentByIds(ctx, clusterId, tmpInsIds, nodeRole)
						if e != nil {
							return tccommon.RetryError(e, tccommon.InternalError)
						}

						initFlag := true
						if result != nil && result.InstanceSet != nil {
							for _, item := range result.InstanceSet {
								if item.InstanceState != nil && *item.InstanceState != "running" {
									initFlag = false
									break
								}
							}

							if initFlag {
								return nil
							}
						}

						return resource.RetryableError(fmt.Errorf("tke master node cvm instance status is not running, retry..."))
					})

					if err != nil {
						return err
					}
				}

				// wait for tke cluster status
				err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
					result, e := tkeService.DescribeKubernetesClusterById(ctx, clusterId)
					if e != nil {
						return tccommon.RetryError(e, tccommon.InternalError)
					}

					if result == nil {

					}

					if result.ClusterStatus != nil && *result.ClusterStatus == "Running" {
						return nil
					}

					return resource.RetryableError(fmt.Errorf("tke status is not RUNNING, retry..."))
				})

				if err != nil {
					return err
				}

			}
		}

		// scale in
		if len(remove) > 0 {
			tmpOld := make([]map[string]interface{}, 0)
			for index := range remove {
				if remove[index] != nil {
					instance := remove[index].(map[string]interface{})
					existedInstance, _ := tkeGetCvmExistInstancesPara(instance)

					insMap := make(map[string]interface{})
					if existedInstance.NodeRole != nil {
						insMap["node_role"] = *existedInstance.NodeRole
					}

					if len(existedInstance.ExistedInstancesPara.InstanceIds) > 0 {
						for _, item := range existedInstance.ExistedInstancesPara.InstanceIds {
							if item != nil {
								insMap["instance_id"] = *item
							}

							tmpOld = append(tmpOld, insMap)
						}
					}
				}
			}

			if len(tmpOld) > 0 {
				request := tkev20180525.NewScaleInClusterMasterRequest()
				request.ClusterId = &clusterId
				for _, item := range tmpOld {
					tmp := tkev20180525.ScaleInMaster{}
					if v, ok := item["node_role"].(string); ok && v != "" {
						tmp.NodeRole = &v
					}

					if v, ok := item["instance_id"].(string); ok && v != "" {
						tmp.InstanceId = &v
					}

					tmp.InstanceDeleteMode = helper.String("retain")
					request.ScaleInMasters = append(request.ScaleInMasters, &tmp)
				}

				err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().ScaleInClusterMasterWithContext(ctx, request)
					if e != nil {
						if sdkErr, ok := e.(*errors.TencentCloudSDKError); ok {
							if sdkErr.GetCode() == "ResourceNotFound" {
								return nil
							}

							if sdkErr.GetCode() == "InvalidParameter" && strings.Contains(sdkErr.GetMessage(), `is not exist`) {
								return nil
							}
						}

						return tccommon.RetryError(e, tccommon.InternalError)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s scale in cluster failed, reason:%+v", logId, err)
					return err
				}

				// wait for tke cluster status
				err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
					result, e := tkeService.DescribeKubernetesClusterById(ctx, clusterId)
					if e != nil {
						return tccommon.RetryError(e, tccommon.InternalError)
					}

					if result == nil {

					}

					if result.ClusterStatus != nil && *result.ClusterStatus == "Running" {
						return nil
					}

					return resource.RetryableError(fmt.Errorf("tke status is not RUNNING, retry..."))
				})

				if err != nil {
					return err
				}
			}
		}
	}

	if err := resourceTencentCloudKubernetesClusterUpdateOnExit(ctx); err != nil {
		return err
	}

	return resourceTencentCloudKubernetesClusterRead(d, meta)
}

func resourceTencentCloudKubernetesClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	clusterId := d.Id()

	var (
		request  = tkev20180525.NewDeleteClusterRequest()
		response = tkev20180525.NewDeleteClusterResponse()
	)

	request.ClusterId = helper.String(clusterId)

	instanceDeleteMode := "terminate"
	if v, ok := d.GetOk("instance_delete_mode"); ok {
		instanceDeleteMode = v.(string)
	}

	request.InstanceDeleteMode = &instanceDeleteMode

	if v, ok := d.GetOk("resource_delete_options"); ok {
		for _, item := range v.(*schema.Set).List() {
			resourceDeleteOptionsMap := item.(map[string]interface{})
			resourceDeleteOption := tkev20180525.ResourceDeleteOption{}
			if v, ok := resourceDeleteOptionsMap["resource_type"]; ok {
				resourceDeleteOption.ResourceType = helper.String(v.(string))
			}
			if v, ok := resourceDeleteOptionsMap["delete_mode"]; ok {
				resourceDeleteOption.DeleteMode = helper.String(v.(string))
			}
			if v, ok := resourceDeleteOptionsMap["skip_deletion_protection"]; ok {
				resourceDeleteOption.SkipDeletionProtection = helper.Bool(v.(bool))
			}
			request.ResourceDeleteOptions = append(request.ResourceDeleteOptions, &resourceDeleteOption)
		}
	}

	if err := resourceTencentCloudKubernetesClusterDeletePostFillRequest0(ctx, request); err != nil {
		return err
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().DeleteClusterWithContext(ctx, request)
		if e != nil {
			if err := resourceTencentCloudKubernetesClusterDeleteRequestOnError0(ctx, e); err != nil {
				return err
			}
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete kubernetes cluster failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	if err := resourceTencentCloudKubernetesClusterDeletePostHandleResponse0(ctx, response); err != nil {
		return err
	}

	return nil
}

func customizeDiffForContainerRuntimeDefault(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// example 1.22.5(maybe 1.22.5-tke.21)
	if clusterVersion, ok := d.GetOk("cluster_version"); ok {
		version := clusterVersion.(string)
		parts := strings.Split(version, ".")
		fmt.Println(parts)
		if len(parts) < 2 {
			log.Printf("[WARN] Invalid cluster version format: %s", version)
			return nil
		}

		mainVersionStr := strings.Split(parts[1], "-")[0]
		mainVersion, err := strconv.Atoi(mainVersionStr)
		fmt.Println(mainVersion)
		if err != nil {
			log.Printf("[WARN] Failed to parse cluster version: %v", err)
			return nil
		}

		runtimeValue := "docker"
		if mainVersion >= 24 {
			runtimeValue = "containerd"
		}

		if _, ok := d.GetOk("container_runtime"); !ok {
			if err := d.SetNew("container_runtime", runtimeValue); err != nil {
				return err
			}
		}
	} else {
		if _, ok := d.GetOk("container_runtime"); !ok {
			if err := d.SetNew("container_runtime", "containerd"); err != nil {
				return err
			}
		}
	}

	return nil
}
