package cvm

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCvmLaunchTemplateVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmLaunchTemplateVersionCreate,
		Read:   resourceTencentCloudCvmLaunchTemplateVersionRead,
		Delete: resourceTencentCloudCvmLaunchTemplateVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"placement": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Location 的 实例. You 可以 使用 此 参数 到 指定attributes 的 实例，such 作为 its availability 可用区，项目，和 CDH (对于 dedicated CVMs)。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "ID availability 可用区 其中 实例 resides. You 可以 call DescribeZones API 和 obtain ID 在 返回 可用区 字段。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "ID 项目 到 其中 实例 belongs. 此 参数 可以 是 获取 从 projectId 返回 通过 DescribeProject. 如果 此 是 left 空， 默认值 项目 是 使用。",
						},
						"host_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "ID 列表 CDHs 从 其中 实例 可以 是 创建. 如果 您 have purchased CDHs 和 指定this 参数， 实例 您 purchase 将 是 randomly deployed 在 CDHs。",
						},
						"host_ips": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "IPs 的 hosts 到 create CVMs。",
						},
					},
				},
			},

			"launch_template_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "实例 launch 模板 ID 此 参数 是 使用 作为 basis 对于 creating new template versions。",
			},

			"launch_template_version": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "此 参数，当 指定，是 用于create 实例 launch templates. 如果 此 参数 是 不 指定， 默认值 版本 将 是 使用。",
			},

			"launch_template_version_description": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "描述 实例 launch template versions. 此 参数 可以 contain 2-256 字符。",
			},

			"instance_type": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "类型 实例. 如果 此 参数 是 不 指定， 系统 将 dynamically 指定default model according 到 资源 sales 在 当前 地域",
			},

			"image_id": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Image ID。",
			},

			"system_disk": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "System 磁盘 配置 信息 的 实例. 如果 此 参数 是 不 指定，它 是 assigned according 到 系统 默认值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "类型 系统 磁盘. 默认值：类型 hard 磁盘 currently 在 stock。",
						},
						"disk_id": {
							Type:        schema.TypeString,
							Computed:    true,
							ForceNew:    true,
							Description: "System 磁盘 ID. System disks whose 类型 是 LOCAL_BASIC 或 LOCAL_SSD do 不 have ID 和 do 不 support 此 参数. It 是 仅 使用 作为 response 参数 对于 APIs such 作为 DescribeInstances，和 不能 是 使用 作为 请求 参数 对于 APIs such 作为 RunInstances。",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "System 磁盘 大小; 单位: GB; 默认值：50 GB。",
						},
						"cdc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "ID dedicated 集群 到 其中 实例 belongs。",
						},
					},
				},
			},

			"data_disks": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "配置 信息 的 实例 数据 disks. 如果 此 参数 是 不 指定，无 数据 磁盘 将 是 purchased 通过 默认值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_size": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    true,
							Description: "Data 磁盘 大小 (在 GB). 最小 adjustment increment 是 10 GB. 值 范围 varies 通过 数据 磁盘 类型",
						},
						"disk_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "类型 数据 磁盘。",
						},
						"disk_id": {
							Type:        schema.TypeString,
							Computed:    true,
							ForceNew:    true,
							Description: "System 磁盘 ID. System disks whose 类型 是 LOCAL_BASIC 或 LOCAL_SSD do 不 have ID 和 do 不 support 此 参数. It 是 仅 使用 作为 response 参数 对于 APIs such 作为 DescribeInstances，和 不能 是 使用 作为 请求 参数 对于 APIs such 作为 RunInstances。",
						},
						"delete_with_instance": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "是否terminate 数据 磁盘 当 its CVM 是 terminated. 默认值：`true`。",
						},
						"snapshot_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Data 磁盘 快照 ID. 大小 的 selected 数据 磁盘 快照 必须 是 smaller 比 该 的 数据 磁盘. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 found。",
						},
						"encrypt": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "指定是否data 磁盘 是 encrypted。",
						},
						"kms_key_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "ID 自定义 CMK 在 格式 的 UUID 或 `kms-abcd1234`。",
						},
						"throughput_performance": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Cloud 磁盘 performance 在 MB/s。",
						},
						"cdc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "ID dedicated 集群 到 其中 实例 belongs。",
						},
					},
				},
			},

			"virtual_private_cloud": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Describes 信息 在 VPC，包括 subnets，IP addresses，etc。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "私有网络 ID 在 格式 的 vpc-xxx，如果 您 指定DEFAULT 对于 both VpcId 和 SubnetId 当 creating 实例， 默认值 VPC 将 是 使用。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "VPC 子网 ID 在 格式 子网-xxx，如果 您 指定DEFAULT 对于 both VpcId 和 SubnetId 当 creating 实例， 默认值 VPC 将 是 使用。",
						},
						"as_vpc_gateway": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "是否use CVM 实例 作为 公有 网关. 公有 网关 是 仅 可用 当 实例 has 公有 IP 和 resides 在 VPC。",
						},
						"private_ip_addresses": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "数组 VPC 子网 IPs. You 可以 使用 此 参数 当 creating 实例 或 modifying VPC attributes 的 实例. Currently 您 可以 指定multiple IPs 在 一个 子网 仅 当 creating 多个 实例 在 same 时间。",
						},
						"ipv6_address_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "数量 IPv6 addresses randomly generated 对于 ENI。",
						},
					},
				},
			},

			"internet_accessible": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Describes accessibility 的 实例 在 公有 网络，包括 its 网络 billing 方法，最大 带宽，etc。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"internet_charge_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Network 连接 billing plan。",
						},
						"internet_max_bandwidth_out": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "最大 outbound 带宽 的 公有 网络，在 Mbps. 默认值为 0 Mbps。",
						},
						"public_ip_assigned": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "是否assign 公有 IP。",
						},
						"bandwidth_package_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Bandwidth 包 ID。",
						},
					},
				},
			},

			"instance_count": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "数量 实例 到 是 purchased。",
			},

			"instance_name": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "实例名称 到 是 displayed。",
			},

			"login_settings": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Describes login settings 的 实例。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Login 密码 的 实例。",
						},
						"key_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "列表 键 IDs. After 实例 是 associated 使用 键，您 可以 访问 实例 使用 私有 键 在 键 pair。",
						},
						"keep_image_login": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "是否keep original settings 的 镜像。",
						},
					},
				},
			},

			"security_group_ids": {
				Optional: true,
				Computed: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security groups 到 其中 实例 belongs. 如果 此 参数 是 不 指定， 实例 将 是 associated 使用 默认值 安全 groups。",
			},

			"enhanced_service": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Enhanced 服务. You 可以 使用 此 参数 到 指定是否enable services such 作为 Anti-DDoS 和 Cloud Monitor. 如果 此 参数 是 不 指定，Cloud Monitor 和 Anti-DDoS 是 已启用 对于 公有 images 通过 默认值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_service": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Enables 云 安全 服务. 如果 此 参数 是 不 指定， 云 安全 服务 将 是 已启用 通过 默认值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "是否enable Cloud Security。",
									},
								},
							},
						},
						"monitor_service": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Enables 云 监控 服务. 如果 此 参数 是 不 指定， 云 监控 服务 将 是 已启用 通过 默认值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "是否enable Cloud Monitor。",
									},
								},
							},
						},
						"automation_service": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "是否enable TAT 服务. 如果 此 参数 是 不 指定， TAT 服务 是 已启用 对于 公有 images 和 已禁用 对于 other images 通过 默认值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "是否enable TAT 服务。",
									},
								},
							},
						},
					},
				},
			},

			"client_token": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "A 唯一 字符串 supplied 通过 客户端 到 ensure 该 请求 是 idempotent. Its 最大 长度 是 64 ASCII 字符. 如果 此 参数 是 不 指定， idem-potency 的 请求 不能 是 guaranteed。",
			},

			"host_name": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Hostname 的 CVM。",
			},

			"action_timer": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Scheduled tasks. You 可以 使用 此 参数 到 指定scheduled tasks 对于 实例. Only scheduled termination 是 支持。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timer_action": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Timer 名称 Currently TerminateInstances 是 仅 支持 值",
						},
						"action_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "执行时间，displayed according 到 ISO8601 standard，和 UTC 时间 是 使用. 格式 是 YYYY-MM-DDThh:mm:ssZ. For 示例，2018-05-29T11:26:40Z， execution 必须 是 在 least 5 minutes later 比 当前 时间。",
						},
						"externals": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Additional 数据。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"release_address": {
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "Release 地址",
									},
									"unsupport_networks": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "Not 支持 网络。",
									},
									"storage_block_attr": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "Information 在 本地 HDD 存储。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Local HDD 存储 类型 值: LOCAL_PRO。",
												},
												"min_size": {
													Type:        schema.TypeInt,
													Required:    true,
													ForceNew:    true,
													Description: "Minimum 容量 的 本地 HDD 存储。",
												},
												"max_size": {
													Type:        schema.TypeInt,
													Required:    true,
													ForceNew:    true,
													Description: "Maximum 容量 的 本地 HDD 存储。",
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

			"disaster_recover_group_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Placement 组 ID You 可以 仅 指定one。",
			},

			"tag_specification": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "描述 标签 associated 使用 资源 实例 during 实例 creation。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "类型 资源 该 标签 是 bound 到。",
						},
						"tags": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "列表 标签",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "标签键",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "标签值",
									},
								},
							},
						},
					},
				},
			},

			"instance_market_options": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Options related 到 bidding requests。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spot_options": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							ForceNew:    true,
							Description: "Options related 到 bidding。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_price": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Bidding 价格。",
									},
									"spot_instance_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "Bidding 请求 类型 Currently 仅 一个-时间 是 支持。",
									},
								},
							},
						},
						"market_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Market 选项 类型 Currently spot 是 仅 支持 值",
						},
					},
				},
			},

			"user_data": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "用户 数据 提供 到 实例. 此 参数 needs 到 是 encoded 在 base64 格式 使用 最大 大小 的 16 KB。",
			},

			"dry_run": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "是否request 是 dry run 仅。",
			},

			"cam_role_name": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "角色 名称 CAM。",
			},

			"hpc_cluster_id": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "HPC 集群 ID. HPC 集群 必须 和 可以 仅 是 指定 对于 high-performance computing 实例。",
			},

			"instance_charge_type": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "charge 类型 实例。",
			},

			"instance_charge_prepaid": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Describes billing 方法 的 实例。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    true,
							Description: "Subscription 周期; 单位: month; 有效值：1，2，3，4，5，6，7，8，9，10，11，12，24，36，48，60。",
						},
						"renew_flag": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "自动续费标识 有效值：NOTIFY_AND_AUTO_RENEW: notify upon expiration 和 renew automatically NOTIFY_AND_MANUAL_RENEW: notify upon expiration 但 do 不 renew automatically DISABLE_NOTIFY_AND_MANUAL_RENEW: neither notify upon expiration nor renew automatically &lt;br&gt;&lt;br&gt;默认值：NOTIFY_AND_MANUAL_RENEW. 如果 此 参数 是 指定 作为 NOTIFY_AND_AUTO_RENEW， 实例 将 是 automatically renewed 在 monthly basis 如果 账号 balance 是 sufficient。",
						},
					},
				},
			},

			"disable_api_termination": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "是否termination protection 是 已启用 `TRUE`: Enable 实例 protection，其中 表示 该 此 实例 可以 不 是 删除 通过 API 操作`FALSE`: Do 不 启用 实例 protection. 默认值：`FALSE`。",
			},
		},
	}
}

func resourceTencentCloudCvmLaunchTemplateVersionCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_launch_template_version.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request          = cvm.NewCreateLaunchTemplateVersionRequest()
		response         = cvm.NewCreateLaunchTemplateVersionResponse()
		launchTemplateId string
	)
	if dMap, ok := helper.InterfacesHeadMap(d, "placement"); ok {
		placement := cvm.Placement{}
		if v, ok := dMap["zone"]; ok {
			placement.Zone = helper.String(v.(string))
		}
		if v, ok := dMap["project_id"]; ok {
			placement.ProjectId = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["host_ids"]; ok {
			hostIdsSet := v.(*schema.Set).List()
			for i := range hostIdsSet {
				hostIds := hostIdsSet[i].(string)
				placement.HostIds = append(placement.HostIds, &hostIds)
			}
		}
		//if v, ok := dMap["host_ips"]; ok {
		//	hostIpsSet := v.(*schema.Set).List()
		//	for i := range hostIpsSet {
		//		hostIps := hostIpsSet[i].(string)
		//		placement.HostIps = append(placement.HostIps, &hostIps)
		//	}
		//}
		request.Placement = &placement
	}

	if v, ok := d.GetOk("launch_template_id"); ok {
		launchTemplateId = v.(string)
		request.LaunchTemplateId = helper.String(launchTemplateId)
	}

	if v, ok := d.GetOkExists("launch_template_version"); ok {
		request.LaunchTemplateVersion = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("launch_template_version_description"); ok {
		request.LaunchTemplateVersionDescription = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("image_id"); ok {
		request.ImageId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "system_disk"); ok {
		systemDisk := cvm.SystemDisk{}
		if v, ok := dMap["disk_type"]; ok {
			systemDisk.DiskType = helper.String(v.(string))
		}
		if v, ok := dMap["disk_size"]; ok {
			systemDisk.DiskSize = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["cdc_id"]; ok {
			systemDisk.CdcId = helper.String(v.(string))
		}
		request.SystemDisk = &systemDisk
	}

	if v, ok := d.GetOk("data_disks"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			dataDisk := cvm.DataDisk{}
			if v, ok := dMap["disk_size"]; ok {
				dataDisk.DiskSize = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["disk_type"]; ok {
				dataDisk.DiskType = helper.String(v.(string))
			}
			if v, ok := dMap["delete_with_instance"]; ok {
				dataDisk.DeleteWithInstance = helper.Bool(v.(bool))
			}
			if v, ok := dMap["snapshot_id"]; ok {
				dataDisk.SnapshotId = helper.String(v.(string))
			}
			if v, ok := dMap["encrypt"]; ok {
				dataDisk.Encrypt = helper.Bool(v.(bool))
			}
			if v, ok := dMap["kms_key_id"]; ok {
				dataDisk.KmsKeyId = helper.String(v.(string))
			}
			if v, ok := dMap["throughput_performance"]; ok {
				dataDisk.ThroughputPerformance = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["cdc_id"]; ok {
				dataDisk.CdcId = helper.String(v.(string))
			}
			request.DataDisks = append(request.DataDisks, &dataDisk)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "virtual_private_cloud"); ok {
		virtualPrivateCloud := cvm.VirtualPrivateCloud{}
		if v, ok := dMap["vpc_id"]; ok {
			virtualPrivateCloud.VpcId = helper.String(v.(string))
		}
		if v, ok := dMap["subnet_id"]; ok {
			virtualPrivateCloud.SubnetId = helper.String(v.(string))
		}
		if v, ok := dMap["as_vpc_gateway"]; ok {
			virtualPrivateCloud.AsVpcGateway = helper.Bool(v.(bool))
		}
		if v, ok := dMap["private_ip_addresses"]; ok {
			privateIpAddressesSet := v.(*schema.Set).List()
			for i := range privateIpAddressesSet {
				privateIpAddresses := privateIpAddressesSet[i].(string)
				virtualPrivateCloud.PrivateIpAddresses = append(virtualPrivateCloud.PrivateIpAddresses, &privateIpAddresses)
			}
		}
		if v, ok := dMap["ipv6_address_count"]; ok {
			virtualPrivateCloud.Ipv6AddressCount = helper.IntUint64(v.(int))
		}
		request.VirtualPrivateCloud = &virtualPrivateCloud
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "internet_accessible"); ok {
		internetAccessible := cvm.InternetAccessible{}
		if v, ok := dMap["internet_charge_type"].(string); ok && v != "" {
			internetAccessible.InternetChargeType = helper.String(v)
		}
		if v, ok := dMap["internet_max_bandwidth_out"]; ok {
			internetAccessible.InternetMaxBandwidthOut = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["public_ip_assigned"]; ok {
			internetAccessible.PublicIpAssigned = helper.Bool(v.(bool))
		}
		if v, ok := dMap["bandwidth_package_id"].(string); ok && v != "" {
			internetAccessible.BandwidthPackageId = helper.String(v)
		}
		request.InternetAccessible = &internetAccessible
	}

	if v, ok := d.GetOkExists("instance_count"); ok {
		request.InstanceCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "login_settings"); ok {
		loginSettings := cvm.LoginSettings{}
		if v, ok := dMap["password"].(string); ok && v != "" {
			loginSettings.Password = helper.String(v)
		}
		if v, ok := dMap["key_ids"]; ok {
			keyIdsSet := v.(*schema.Set).List()
			for i := range keyIdsSet {
				keyIds := keyIdsSet[i].(string)
				loginSettings.KeyIds = append(loginSettings.KeyIds, &keyIds)
			}
		}
		if v, ok := dMap["keep_image_login"].(string); ok && v != "" {
			loginSettings.KeepImageLogin = helper.String(v)
		}
		request.LoginSettings = &loginSettings
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "enhanced_service"); ok {
		enhancedService := cvm.EnhancedService{}
		if securityServiceMap, ok := helper.InterfaceToMap(dMap, "security_service"); ok {
			runSecurityServiceEnabled := cvm.RunSecurityServiceEnabled{}
			if v, ok := securityServiceMap["enabled"]; ok {
				runSecurityServiceEnabled.Enabled = helper.Bool(v.(bool))
			}
			enhancedService.SecurityService = &runSecurityServiceEnabled
		}
		if monitorServiceMap, ok := helper.InterfaceToMap(dMap, "monitor_service"); ok {
			runMonitorServiceEnabled := cvm.RunMonitorServiceEnabled{}
			if v, ok := monitorServiceMap["enabled"]; ok {
				runMonitorServiceEnabled.Enabled = helper.Bool(v.(bool))
			}
			enhancedService.MonitorService = &runMonitorServiceEnabled
		}
		if automationServiceMap, ok := helper.InterfaceToMap(dMap, "automation_service"); ok {
			runAutomationServiceEnabled := cvm.RunAutomationServiceEnabled{}
			if v, ok := automationServiceMap["enabled"]; ok {
				runAutomationServiceEnabled.Enabled = helper.Bool(v.(bool))
			}
			enhancedService.AutomationService = &runAutomationServiceEnabled
		}
		request.EnhancedService = &enhancedService
	}

	if v, ok := d.GetOk("client_token"); ok {
		request.ClientToken = helper.String(v.(string))
	}

	if v, ok := d.GetOk("host_name"); ok {
		request.HostName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "action_timer"); ok {
		actionTimer := cvm.ActionTimer{}
		if v, ok := dMap["timer_action"]; ok {
			actionTimer.TimerAction = helper.String(v.(string))
		}
		if v, ok := dMap["action_time"]; ok {
			actionTimer.ActionTime = helper.String(v.(string))
		}
		if externalsMap, ok := helper.InterfaceToMap(dMap, "externals"); ok {
			externals := cvm.Externals{}
			if v, ok := externalsMap["release_address"]; ok {
				externals.ReleaseAddress = helper.Bool(v.(bool))
			}
			if v, ok := externalsMap["unsupport_networks"]; ok {
				unsupportNetworksSet := v.(*schema.Set).List()
				for i := range unsupportNetworksSet {
					unsupportNetworks := unsupportNetworksSet[i].(string)
					externals.UnsupportNetworks = append(externals.UnsupportNetworks, &unsupportNetworks)
				}
			}
			if storageBlockAttrMap, ok := helper.InterfaceToMap(externalsMap, "storage_block_attr"); ok {
				storageBlock := cvm.StorageBlock{}
				if v, ok := storageBlockAttrMap["type"]; ok {
					storageBlock.Type = helper.String(v.(string))
				}
				if v, ok := storageBlockAttrMap["min_size"]; ok {
					storageBlock.MinSize = helper.IntInt64(v.(int))
				}
				if v, ok := storageBlockAttrMap["max_size"]; ok {
					storageBlock.MaxSize = helper.IntInt64(v.(int))
				}
				externals.StorageBlockAttr = &storageBlock
			}
			actionTimer.Externals = &externals
		}
		request.ActionTimer = &actionTimer
	}

	if v, ok := d.GetOk("disaster_recover_group_ids"); ok {
		disasterRecoverGroupIdsSet := v.(*schema.Set).List()
		for i := range disasterRecoverGroupIdsSet {
			disasterRecoverGroupIds := disasterRecoverGroupIdsSet[i].(string)
			request.DisasterRecoverGroupIds = append(request.DisasterRecoverGroupIds, &disasterRecoverGroupIds)
		}
	}

	if v, ok := d.GetOk("tag_specification"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			tagSpecification := cvm.TagSpecification{}
			if v, ok := dMap["resource_type"]; ok {
				tagSpecification.ResourceType = helper.String(v.(string))
			}
			if v, ok := dMap["tags"]; ok {
				for _, item := range v.([]interface{}) {
					tagsMap := item.(map[string]interface{})
					tag := cvm.Tag{}
					if v, ok := tagsMap["key"]; ok {
						tag.Key = helper.String(v.(string))
					}
					if v, ok := tagsMap["value"]; ok {
						tag.Value = helper.String(v.(string))
					}
					tagSpecification.Tags = append(tagSpecification.Tags, &tag)
				}
			}
			request.TagSpecification = append(request.TagSpecification, &tagSpecification)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "instance_market_options"); ok {
		instanceMarketOptionsRequest := cvm.InstanceMarketOptionsRequest{}
		if spotOptionsMap, ok := helper.InterfaceToMap(dMap, "spot_options"); ok {
			spotMarketOptions := cvm.SpotMarketOptions{}
			if v, ok := spotOptionsMap["max_price"]; ok {
				spotMarketOptions.MaxPrice = helper.String(v.(string))
			}
			if v, ok := spotOptionsMap["spot_instance_type"]; ok {
				spotMarketOptions.SpotInstanceType = helper.String(v.(string))
			}
			instanceMarketOptionsRequest.SpotOptions = &spotMarketOptions
		}
		if v, ok := dMap["market_type"]; ok {
			instanceMarketOptionsRequest.MarketType = helper.String(v.(string))
		}
		request.InstanceMarketOptions = &instanceMarketOptionsRequest
	}

	if v, ok := d.GetOk("user_data"); ok {
		request.UserData = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request.DryRun = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("cam_role_name"); ok {
		request.CamRoleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("hpc_cluster_id"); ok {
		request.HpcClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "instance_charge_prepaid"); ok {
		instanceChargePrepaid := cvm.InstanceChargePrepaid{}
		if v, ok := dMap["period"]; ok {
			instanceChargePrepaid.Period = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["renew_flag"]; ok {
			instanceChargePrepaid.RenewFlag = helper.String(v.(string))
		}
		request.InstanceChargePrepaid = &instanceChargePrepaid
	}

	if v, ok := d.GetOkExists("disable_api_termination"); ok {
		request.DisableApiTermination = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().CreateLaunchTemplateVersion(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cvm launchTemplateVersion failed, reason:%+v", logId, err)
		return err
	}

	launchTemplateVersionNumber := *response.Response.LaunchTemplateVersionNumber
	launchTemplateVersionNumberString := strconv.FormatInt(launchTemplateVersionNumber, 10)
	d.SetId(launchTemplateId + tccommon.FILED_SP + launchTemplateVersionNumberString)

	return resourceTencentCloudCvmLaunchTemplateVersionRead(d, meta)
}

func resourceTencentCloudCvmLaunchTemplateVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_launch_template_version.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	launchTemplateId := idSplit[0]
	launchTemplateVersionNumber := idSplit[1]

	launchTemplateVersion, err := service.DescribeCvmLaunchTemplateVersionById(ctx, launchTemplateId, launchTemplateVersionNumber)
	if err != nil {
		return err
	}

	if launchTemplateVersion == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CvmLaunchTemplateVersion` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if launchTemplateVersion.LaunchTemplateId != nil {
		_ = d.Set("launch_template_id", launchTemplateVersion.LaunchTemplateId)
	}

	if launchTemplateVersion.LaunchTemplateVersion != nil {
		_ = d.Set("launch_template_version", launchTemplateVersion.LaunchTemplateVersion)
	}

	if launchTemplateVersion.LaunchTemplateVersionDescription != nil {
		_ = d.Set("launch_template_version_description", launchTemplateVersion.LaunchTemplateVersionDescription)
	}

	if launchTemplateVersion.LaunchTemplateVersionData != nil {
		if launchTemplateVersion.LaunchTemplateVersionData.InstanceType != nil {
			_ = d.Set("instance_type", launchTemplateVersion.LaunchTemplateVersionData.InstanceType)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.ImageId != nil {
			_ = d.Set("image_id", launchTemplateVersion.LaunchTemplateVersionData.ImageId)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.Placement != nil {
			placementMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.Placement.Zone != nil {
				placementMap["zone"] = launchTemplateVersion.LaunchTemplateVersionData.Placement.Zone
			}

			if launchTemplateVersion.LaunchTemplateVersionData.Placement.ProjectId != nil {
				placementMap["project_id"] = launchTemplateVersion.LaunchTemplateVersionData.Placement.ProjectId
			}

			if launchTemplateVersion.LaunchTemplateVersionData.Placement.HostIds != nil {
				placementMap["host_ids"] = launchTemplateVersion.LaunchTemplateVersionData.Placement.HostIds
			}

			//if launchTemplateVersion.LaunchTemplateVersionData.Placement.HostIps != nil {
			//	placementMap["host_ips"] = launchTemplateVersion.LaunchTemplateVersionData.Placement.HostIps
			//}

			_ = d.Set("placement", []interface{}{placementMap})
		}

		if launchTemplateVersion.LaunchTemplateVersionData.SystemDisk != nil {
			systemDiskMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.DiskType != nil {
				systemDiskMap["disk_type"] = launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.DiskType
			}

			if launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.DiskId != nil {
				systemDiskMap["disk_id"] = launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.DiskId
			}

			if launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.DiskSize != nil {
				systemDiskMap["disk_size"] = launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.DiskSize
			}

			if launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.CdcId != nil {
				systemDiskMap["cdc_id"] = launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.CdcId
			}

			_ = d.Set("system_disk", []interface{}{systemDiskMap})
		}

		if launchTemplateVersion.LaunchTemplateVersionData.DataDisks != nil {
			dataDisksList := []interface{}{}
			for _, dataDisk := range launchTemplateVersion.LaunchTemplateVersionData.DataDisks {
				dataDisksMap := map[string]interface{}{}

				if dataDisk.DiskSize != nil {
					dataDisksMap["disk_size"] = dataDisk.DiskSize
				}

				if dataDisk.DiskType != nil {
					dataDisksMap["disk_type"] = dataDisk.DiskType
				}

				if dataDisk.DiskId != nil {
					dataDisksMap["disk_id"] = dataDisk.DiskId
				}

				if dataDisk.DeleteWithInstance != nil {
					dataDisksMap["delete_with_instance"] = dataDisk.DeleteWithInstance
				}

				if dataDisk.SnapshotId != nil {
					dataDisksMap["snapshot_id"] = dataDisk.SnapshotId
				}

				if dataDisk.Encrypt != nil {
					dataDisksMap["encrypt"] = dataDisk.Encrypt
				}

				if dataDisk.KmsKeyId != nil {
					dataDisksMap["kms_key_id"] = dataDisk.KmsKeyId
				}

				if dataDisk.ThroughputPerformance != nil {
					dataDisksMap["throughput_performance"] = dataDisk.ThroughputPerformance
				}

				if dataDisk.CdcId != nil {
					dataDisksMap["cdc_id"] = dataDisk.CdcId
				}

				dataDisksList = append(dataDisksList, dataDisksMap)
			}

			_ = d.Set("data_disks", dataDisksList)

		}

		if launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud != nil {
			virtualPrivateCloudMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.VpcId != nil {
				virtualPrivateCloudMap["vpc_id"] = launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.VpcId
			}

			if launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.SubnetId != nil {
				virtualPrivateCloudMap["subnet_id"] = launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.SubnetId
			}

			if launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.AsVpcGateway != nil {
				virtualPrivateCloudMap["as_vpc_gateway"] = launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.AsVpcGateway
			}

			if launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.PrivateIpAddresses != nil {
				virtualPrivateCloudMap["private_ip_addresses"] = launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.PrivateIpAddresses
			}

			if launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.Ipv6AddressCount != nil {
				virtualPrivateCloudMap["ipv6_address_count"] = launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.Ipv6AddressCount
			}

			_ = d.Set("virtual_private_cloud", []interface{}{virtualPrivateCloudMap})
		}

		if launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible != nil {
			internetAccessibleMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.InternetChargeType != nil {
				internetAccessibleMap["internet_charge_type"] = launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.InternetChargeType
			}

			if launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.InternetMaxBandwidthOut != nil {
				internetAccessibleMap["internet_max_bandwidth_out"] = launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.InternetMaxBandwidthOut
			}

			if launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.PublicIpAssigned != nil {
				internetAccessibleMap["public_ip_assigned"] = launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.PublicIpAssigned
			}

			if launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.BandwidthPackageId != nil {
				internetAccessibleMap["bandwidth_package_id"] = launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.BandwidthPackageId
			}

			_ = d.Set("internet_accessible", []interface{}{internetAccessibleMap})
		}

		if launchTemplateVersion.LaunchTemplateVersionData.InstanceCount != nil {
			_ = d.Set("instance_count", launchTemplateVersion.LaunchTemplateVersionData.InstanceCount)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.InstanceName != nil {
			_ = d.Set("instance_name", launchTemplateVersion.LaunchTemplateVersionData.InstanceName)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.LoginSettings != nil {
			loginSettingsMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.LoginSettings.Password != nil {
				loginSettingsMap["password"] = launchTemplateVersion.LaunchTemplateVersionData.LoginSettings.Password
			}

			if launchTemplateVersion.LaunchTemplateVersionData.LoginSettings.KeyIds != nil {
				loginSettingsMap["key_ids"] = launchTemplateVersion.LaunchTemplateVersionData.LoginSettings.KeyIds
			}

			if launchTemplateVersion.LaunchTemplateVersionData.LoginSettings.KeepImageLogin != nil {
				loginSettingsMap["keep_image_login"] = launchTemplateVersion.LaunchTemplateVersionData.LoginSettings.KeepImageLogin
			}

			_ = d.Set("login_settings", []interface{}{loginSettingsMap})
		}

		if launchTemplateVersion.LaunchTemplateVersionData.SecurityGroupIds != nil {
			_ = d.Set("security_group_ids", launchTemplateVersion.LaunchTemplateVersionData.SecurityGroupIds)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.EnhancedService != nil {
			enhancedServiceMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.SecurityService != nil {
				securityServiceMap := map[string]interface{}{}

				if launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.SecurityService.Enabled != nil {
					securityServiceMap["enabled"] = launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.SecurityService.Enabled
				}

				enhancedServiceMap["security_service"] = []interface{}{securityServiceMap}
			}

			if launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.MonitorService != nil {
				monitorServiceMap := map[string]interface{}{}

				if launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.MonitorService.Enabled != nil {
					monitorServiceMap["enabled"] = launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.MonitorService.Enabled
				}

				enhancedServiceMap["monitor_service"] = []interface{}{monitorServiceMap}
			}

			if launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.AutomationService != nil {
				automationServiceMap := map[string]interface{}{}

				if launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.AutomationService.Enabled != nil {
					automationServiceMap["enabled"] = launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.AutomationService.Enabled
				}

				enhancedServiceMap["automation_service"] = []interface{}{automationServiceMap}
			}

			_ = d.Set("enhanced_service", []interface{}{enhancedServiceMap})
		}

		if launchTemplateVersion.LaunchTemplateVersionData.ClientToken != nil {
			_ = d.Set("client_token", launchTemplateVersion.LaunchTemplateVersionData.ClientToken)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.HostName != nil {
			_ = d.Set("host_name", launchTemplateVersion.LaunchTemplateVersionData.HostName)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer != nil {
			actionTimerMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.TimerAction != nil {
				actionTimerMap["timer_action"] = launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.TimerAction
			}

			if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.ActionTime != nil {
				actionTimerMap["action_time"] = launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.ActionTime
			}

			if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals != nil {
				externalsMap := map[string]interface{}{}

				if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.ReleaseAddress != nil {
					externalsMap["release_address"] = launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.ReleaseAddress
				}

				if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.UnsupportNetworks != nil {
					externalsMap["unsupport_networks"] = launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.UnsupportNetworks
				}

				if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.StorageBlockAttr != nil {
					storageBlockAttrMap := map[string]interface{}{}

					if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.StorageBlockAttr.Type != nil {
						storageBlockAttrMap["type"] = launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.StorageBlockAttr.Type
					}

					if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.StorageBlockAttr.MinSize != nil {
						storageBlockAttrMap["min_size"] = launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.StorageBlockAttr.MinSize
					}

					if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.StorageBlockAttr.MaxSize != nil {
						storageBlockAttrMap["max_size"] = launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.StorageBlockAttr.MaxSize
					}

					externalsMap["storage_block_attr"] = []interface{}{storageBlockAttrMap}
				}

				actionTimerMap["externals"] = []interface{}{externalsMap}
			}

			_ = d.Set("action_timer", []interface{}{actionTimerMap})
		}

		if launchTemplateVersion.LaunchTemplateVersionData.DisasterRecoverGroupIds != nil {
			_ = d.Set("disaster_recover_group_ids", launchTemplateVersion.LaunchTemplateVersionData.DisasterRecoverGroupIds)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.TagSpecification != nil {
			tagSpecificationList := []interface{}{}
			for _, tagSpecification := range launchTemplateVersion.LaunchTemplateVersionData.TagSpecification {
				tagSpecificationMap := map[string]interface{}{}

				if tagSpecification.ResourceType != nil {
					tagSpecificationMap["resource_type"] = tagSpecification.ResourceType
				}

				if tagSpecification.Tags != nil {
					tagsList := []interface{}{}
					for _, tag := range tagSpecification.Tags {
						tagsMap := map[string]interface{}{}

						if tag.Key != nil {
							tagsMap["key"] = tag.Key
						}

						if tag.Value != nil {
							tagsMap["value"] = tag.Value
						}

						tagsList = append(tagsList, tagsMap)
					}

					tagSpecificationMap["tags"] = []interface{}{tagsList}
				}

				tagSpecificationList = append(tagSpecificationList, tagSpecificationMap)
			}

			_ = d.Set("tag_specification", tagSpecificationList)

		}

		if launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions != nil {
			instanceMarketOptionsMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions.SpotOptions != nil {
				spotOptionsMap := map[string]interface{}{}

				if launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions.SpotOptions.MaxPrice != nil {
					spotOptionsMap["max_price"] = launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions.SpotOptions.MaxPrice
				}

				if launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions.SpotOptions.SpotInstanceType != nil {
					spotOptionsMap["spot_instance_type"] = launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions.SpotOptions.SpotInstanceType
				}

				instanceMarketOptionsMap["spot_options"] = []interface{}{spotOptionsMap}
			}

			if launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions.MarketType != nil {
				instanceMarketOptionsMap["market_type"] = launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions.MarketType
			}

			_ = d.Set("instance_market_options", []interface{}{instanceMarketOptionsMap})
		}

		if launchTemplateVersion.LaunchTemplateVersionData.UserData != nil {
			_ = d.Set("user_data", launchTemplateVersion.LaunchTemplateVersionData.UserData)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.CamRoleName != nil {
			_ = d.Set("cam_role_name", launchTemplateVersion.LaunchTemplateVersionData.CamRoleName)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.HpcClusterId != nil {
			_ = d.Set("hpc_cluster_id", launchTemplateVersion.LaunchTemplateVersionData.HpcClusterId)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.InstanceChargeType != nil {
			_ = d.Set("instance_charge_type", launchTemplateVersion.LaunchTemplateVersionData.InstanceChargeType)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.InstanceChargePrepaid != nil {
			instanceChargePrepaidMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.InstanceChargePrepaid.Period != nil {
				instanceChargePrepaidMap["period"] = launchTemplateVersion.LaunchTemplateVersionData.InstanceChargePrepaid.Period
			}

			if launchTemplateVersion.LaunchTemplateVersionData.InstanceChargePrepaid.RenewFlag != nil {
				instanceChargePrepaidMap["renew_flag"] = launchTemplateVersion.LaunchTemplateVersionData.InstanceChargePrepaid.RenewFlag
			}

			_ = d.Set("instance_charge_prepaid", []interface{}{instanceChargePrepaidMap})
		}

		if launchTemplateVersion.LaunchTemplateVersionData.DisableApiTermination != nil {
			_ = d.Set("disable_api_termination", launchTemplateVersion.LaunchTemplateVersionData.DisableApiTermination)
		}
	}

	return nil
}

func resourceTencentCloudCvmLaunchTemplateVersionDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_launch_template_version.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	launchTemplateId := idSplit[0]
	launchTemplateVersionNumber := idSplit[1]

	if err := service.DeleteCvmLaunchTemplateVersionById(ctx, launchTemplateId, launchTemplateVersionNumber); err != nil {
		return err
	}

	return nil
}
