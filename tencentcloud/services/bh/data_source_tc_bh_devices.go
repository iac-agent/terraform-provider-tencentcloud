package bh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bhv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudBhDevices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudBhDevicesRead,
		Schema: map[string]*schema.Schema{
			"id_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Asset ID collection。",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Asset 名称 或 asset IP，fuzzy search。",
			},

			"ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Not currently 使用。",
			},

			"ap_code_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "地域 代码 collection。",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"kind": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Operating 系统 类型，1 - Linux，2 - Windows，3 - MySQL，4 - SQLServer。",
			},

			"authorized_user_id_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "用户 ID collection 使用 访问 到 此 asset。",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"resource_id_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "过滤器 condition，asset-bound bastion 主机 服务 ID collection。",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"kind_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Can 过滤器 通过 多个 types，1 - Linux，2 - Windows，3 - MySQL，4 - SQLServer。",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"managed_account": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "是否asset 包含managed accounts. 1，contains; 0，does 不 contain。",
			},

			"department_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "过滤器 condition，可以 过滤器 通过 department ID。",
			},

			"account_id_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Cloud 账号 ID 到 其中 asset belongs。",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"provider_type_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Cloud provider 类型，1 - Tencent Cloud，2 - Alibaba Cloud。",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"cloud_device_status_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Synchronized 云 asset 状态，marking 状态 synchronized assets，0 - 删除，1 - normal，2 - isolated，3 - expired。",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"tag_filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "过滤器 condition, 可以 过滤器 通过 标签 键 和 标签 值. 如果 both 标签 键 和 标签 值 过滤器 conditions 是 指定, they have \"AND\" relationship.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "标签键",
						},
						"tag_value": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "标签值",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "过滤器 数组。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field 到 过滤器. Support: BindingStatus，实例 ID，DeviceAccount，VpcId，DomainId，ResourceId，名称，Ip，ManageDimension。",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "过滤器 值 对于 字段. \nIf 多个 Filters exist， relationship between Filters 是 logical AND. \nIf 多个 Values exist 对于 same 过滤器， relationship between Values under same 过滤器 是 logical OR。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"device_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Asset 信息 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Asset ID。",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID，corresponding 到 CVM，CDB 和 other 实例 IDs。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Asset 名称",
						},
						"public_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public IP。",
						},
						"private_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Private IP。",
						},
						"ap_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 代码",
						},
						"ap_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域名称",
						},
						"os_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operating 系统 名称",
						},
						"kind": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Asset 类型 1 - Linux，2 - Windows，3 - MySQL，4 - SQLServer。",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Management 端口",
						},
						"group_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Asset 组 列表 到 其中 它 belongs。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "组 ID",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "组名称",
									},
									"department": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Department 信息 到 其中 它 belongs。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Department ID。",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Department 名称，1 - 256 字符。",
												},
												"managers": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "Department administrator 账号 ID。",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"manager_users": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Administrator users。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"manager_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Administrator ID。",
															},
															"manager_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Administrator 名称",
															},
														},
													},
												},
											},
										},
									},
									"count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Count。",
									},
								},
							},
						},
						"account_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 accounts bound 到 asset。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "私有网络 ID",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "子网 ID",
						},
						"resource": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Bastion 主机 服务 信息，note 该 它 是 null 当 无 服务 是 bound。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service 实例 ID，such 作为 bh-saas-s3ed4r5e。",
									},
									"ap_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "地域 代码",
									},
									"sv_args": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service 实例 规格 信息。",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "私有网络 ID",
									},
									"nodes": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 assets corresponding 到 服务 规格。",
									},
									"renew_flag": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Auto-renewal flag，0 - 默认值 state，1 - auto-renewal，2 - explicitly 不 auto-renewal。",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "过期时间。",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Resource 状态，0 - 不 initialized，1 - normal，2 - isolated，3 - destroyed，4 - initialization failed，5 - initializing。",
									},
									"resource_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service 实例名称，such 作为 T-Sec-Bastion 主机 (SaaS 类型)。",
									},
									"pid": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Pricing model ID。",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resource 创建时间。",
									},
									"product_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Product 代码，p_cds_dasb。",
									},
									"sub_product_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Sub-product 代码，sp_cds_dasb_bh_saas。",
									},
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Availability 可用区",
									},
									"expired": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether expired，true - expired，false - 不 expired。",
									},
									"deployed": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether deployed，true - deployed，false - 不 deployed。",
									},
									"vpc_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC 名称 其中 服务 是 deployed。",
									},
									"vpc_cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CIDR block 的 VPC 其中 服务 是 deployed。",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "子网 ID 其中 服务 是 deployed。",
									},
									"subnet_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet 名称 其中 服务 是 deployed。",
									},
									"cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CIDR block 的 子网 其中 服务 是 deployed。",
									},
									"public_ip_set": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "External IP。",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"private_ip_set": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Internal IP。",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"module_set": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Advanced 功能 列表 已启用 对于 服务，such 作为: [DB]。",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"used_nodes": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 使用 authorization points。",
									},
									"extend_points": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Extension points。",
									},
									"package_bandwidth": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 带宽 extension packages (4M)。",
									},
									"package_node": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 authorization point extension packages (50 points)。",
									},
									"log_delivery_args": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Log delivery 规格 信息。",
									},
									"clb_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Bastion 主机 资源 load balancer。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"clb_ip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Load balancer IP。",
												},
											},
										},
									},
									"domain_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 网络 domains。",
									},
									"used_domain_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 网络 domains already 使用。",
									},
									"trial": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "0 non-trial 版本，1 trial 版本",
									},
									"log_delivery": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Log delivery 规格 信息。",
									},
									"cdc_cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CDC 集群 ID。",
									},
									"deploy_model": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Deployment 模式，默认值 0，0-cvm 1-tke。",
									},
									"intranet_access": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "0 默认值，non-intranet 访问，1 intranet 访问，2 intranet 访问 opening，3 intranet 访问 closing。",
									},
									"intranet_private_ip_set": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "IP addresses 对于 intranet 访问。",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"intranet_vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC 对于 enabling intranet 访问。",
									},
									"intranet_subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "子网 ID 对于 enabling intranet 访问。",
									},
									"intranet_vpc_cidr": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CIDR block 的 VPC 对于 enabling intranet 访问。",
									},
									"domain_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom 域名 名称 对于 bastion 主机 intranet IP。",
									},
									"share_clb": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否share CLB，true-shared CLB，false-dedicated CLB。",
									},
									"open_clb_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Shared CLB ID",
									},
									"lb_vip_isp": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ISP 信息。",
									},
									"tui_cmd_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Linux asset command line operation 端口",
									},
									"tui_direct_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Linux asset direct 连接 端口",
									},
									"web_access": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "1 默认值，web 访问 已启用，0 web 访问 已禁用，2 web 访问 opening，3 web 访问 closing。",
									},
									"client_access": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "1 默认值，客户端 访问 已启用，0 客户端 访问 已禁用，2 客户端 访问 opening，3 客户端 访问 closing。",
									},
									"external_access": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "1 默认值，外部 访问 已启用，0 外部 访问 已禁用，2 外部 访问 opening，3 外部 访问 closing。",
									},
									"ioa_resource": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "0 默认值，0-free 版本 (trial 版本) IOA，1-paid 版本 IOA。",
									},
									"package_ioa_user_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 zero trust bastion 主机 用户 extension packages，1 extension 包 corresponds 到 20 users。",
									},
									"package_ioa_bandwidth": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 zero trust bastion 主机 带宽 extension packages，一个 extension 包 表示 4M 带宽。",
									},
									"ioa_resource_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Zero trust 实例 ID corresponding 到 bastion 主机 实例。",
									},
								},
							},
						},
						"department": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Department 到 其中 asset belongs。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Department ID。",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Department 名称，1 - 256 字符。",
									},
									"managers": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Department administrator 账号 ID。",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"manager_users": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Administrator users。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"manager_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Administrator ID。",
												},
												"manager_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Administrator 名称",
												},
											},
										},
									},
								},
							},
						},
						"ip_port_set": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Multi-节点 信息 对于 数据库 assets。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network 域名 ID。",
						},
						"domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network 域名 名称",
						},
						"enable_ssl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether SSL 是 已启用，仅 支持 Redis assets，0: 已禁用 1: 已启用",
						},
						"ssl_cert_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 uploaded SSL 证书。",
						},
						"ioa_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "资源 ID 在 IOA side。",
						},
						"manage_dimension": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "K8S 集群 management dimension，1-集群，2-命名空间，3-workload。",
						},
						"manage_account_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "K8S 集群 management 账号 ID。",
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "K8S 集群 命名空间。",
						},
						"workload": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "K8S 集群 workload。",
						},
						"sync_pod_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 synchronized pods 在 K8S 集群。",
						},
						"total_pod_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total 数量 pods 在 K8S 集群。",
						},
						"cloud_account_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cloud 账号 ID。",
						},
						"cloud_account_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud 账号 名称",
						},
						"provider_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cloud provider 类型，1-Tencent Cloud，2-Alibaba Cloud。",
						},
						"provider_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud provider 名称",
						},
						"sync_cloud_device_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Synchronized 云 asset 状态，marking 状态 synchronized assets，0-删除，1-normal，2-isolated，3-expired。",
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

func dataSourceTencentCloudBhDevicesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_bh_devices.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = BhService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("id_set"); ok {
		idSetList := []*uint64{}
		idSetSet := v.(*schema.Set).List()
		for i := range idSetSet {
			idSet := idSetSet[i].(int)
			idSetList = append(idSetList, helper.IntUint64(idSet))
		}

		paramMap["IdSet"] = idSetList
	}

	if v, ok := d.GetOk("name"); ok {
		paramMap["Name"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ip"); ok {
		paramMap["Ip"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ap_code_set"); ok {
		apCodeSetList := []*string{}
		apCodeSetSet := v.(*schema.Set).List()
		for i := range apCodeSetSet {
			apCodeSet := apCodeSetSet[i].(string)
			apCodeSetList = append(apCodeSetList, helper.String(apCodeSet))
		}

		paramMap["ApCodeSet"] = apCodeSetList
	}

	if v, ok := d.GetOkExists("kind"); ok {
		paramMap["Kind"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("authorized_user_id_set"); ok {
		authorizedUserIdSetList := []*uint64{}
		authorizedUserIdSetSet := v.(*schema.Set).List()
		for i := range authorizedUserIdSetSet {
			authorizedUserIdSet := authorizedUserIdSetSet[i].(int)
			authorizedUserIdSetList = append(authorizedUserIdSetList, helper.IntUint64(authorizedUserIdSet))
		}

		paramMap["AuthorizedUserIdSet"] = authorizedUserIdSetList
	}

	if v, ok := d.GetOk("resource_id_set"); ok {
		resourceIdSetList := []*string{}
		resourceIdSetSet := v.(*schema.Set).List()
		for i := range resourceIdSetSet {
			resourceIdSet := resourceIdSetSet[i].(string)
			resourceIdSetList = append(resourceIdSetList, helper.String(resourceIdSet))
		}

		paramMap["ResourceIdSet"] = resourceIdSetList
	}

	if v, ok := d.GetOk("kind_set"); ok {
		kindSetList := []*uint64{}
		kindSetSet := v.(*schema.Set).List()
		for i := range kindSetSet {
			kindSet := kindSetSet[i].(int)
			kindSetList = append(kindSetList, helper.IntUint64(kindSet))
		}

		paramMap["KindSet"] = kindSetList
	}

	if v, ok := d.GetOk("managed_account"); ok {
		paramMap["ManagedAccount"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("department_id"); ok {
		paramMap["DepartmentId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("account_id_set"); ok {
		accountIdSetList := []*uint64{}
		accountIdSetSet := v.(*schema.Set).List()
		for i := range accountIdSetSet {
			accountIdSet := accountIdSetSet[i].(int)
			accountIdSetList = append(accountIdSetList, helper.IntUint64(accountIdSet))
		}

		paramMap["AccountIdSet"] = accountIdSetList
	}

	if v, ok := d.GetOk("provider_type_set"); ok {
		providerTypeSetList := []*uint64{}
		providerTypeSetSet := v.(*schema.Set).List()
		for i := range providerTypeSetSet {
			providerTypeSet := providerTypeSetSet[i].(int)
			providerTypeSetList = append(providerTypeSetList, helper.IntUint64(providerTypeSet))
		}

		paramMap["ProviderTypeSet"] = providerTypeSetList
	}

	if v, ok := d.GetOk("cloud_device_status_set"); ok {
		cloudDeviceStatusSetList := []*uint64{}
		cloudDeviceStatusSetSet := v.(*schema.Set).List()
		for i := range cloudDeviceStatusSetSet {
			cloudDeviceStatusSet := cloudDeviceStatusSetSet[i].(int)
			cloudDeviceStatusSetList = append(cloudDeviceStatusSetList, helper.IntUint64(cloudDeviceStatusSet))
		}

		paramMap["CloudDeviceStatusSet"] = cloudDeviceStatusSetList
	}

	if v, ok := d.GetOk("tag_filters"); ok {
		tagFiltersSet := v.([]interface{})
		tmpSet := make([]*bhv20230418.TagFilter, 0, len(tagFiltersSet))
		for _, item := range tagFiltersSet {
			tagFiltersMap := item.(map[string]interface{})
			tagFilter := bhv20230418.TagFilter{}
			if v, ok := tagFiltersMap["tag_key"].(string); ok && v != "" {
				tagFilter.TagKey = helper.String(v)
			}

			if v, ok := tagFiltersMap["tag_value"]; ok {
				tagValueSet := v.(*schema.Set).List()
				for i := range tagValueSet {
					tagValue := tagValueSet[i].(string)
					tagFilter.TagValue = append(tagFilter.TagValue, helper.String(tagValue))
				}
			}

			tmpSet = append(tmpSet, &tagFilter)
		}

		paramMap["TagFilters"] = tmpSet
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*bhv20230418.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := bhv20230418.Filter{}
			if v, ok := filtersMap["name"].(string); ok && v != "" {
				filter.Name = helper.String(v)
			}

			if v, ok := filtersMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				for i := range valuesSet {
					values := valuesSet[i].(string)
					filter.Values = append(filter.Values, helper.String(values))
				}
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	var respData []*bhv20230418.Device
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeBhDevicesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	deviceSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, deviceSet := range respData {
			deviceSetMap := map[string]interface{}{}
			if deviceSet.Id != nil {
				deviceSetMap["id"] = deviceSet.Id
			}

			if deviceSet.InstanceId != nil {
				deviceSetMap["instance_id"] = deviceSet.InstanceId
			}

			if deviceSet.Name != nil {
				deviceSetMap["name"] = deviceSet.Name
			}

			if deviceSet.PublicIp != nil {
				deviceSetMap["public_ip"] = deviceSet.PublicIp
			}

			if deviceSet.PrivateIp != nil {
				deviceSetMap["private_ip"] = deviceSet.PrivateIp
			}

			if deviceSet.ApCode != nil {
				deviceSetMap["ap_code"] = deviceSet.ApCode
			}

			if deviceSet.ApName != nil {
				deviceSetMap["ap_name"] = deviceSet.ApName
			}

			if deviceSet.OsName != nil {
				deviceSetMap["os_name"] = deviceSet.OsName
			}

			if deviceSet.Kind != nil {
				deviceSetMap["kind"] = deviceSet.Kind
			}

			if deviceSet.Port != nil {
				deviceSetMap["port"] = deviceSet.Port
			}

			groupSetList := make([]map[string]interface{}, 0, len(deviceSet.GroupSet))
			if deviceSet.GroupSet != nil {
				for _, groupSet := range deviceSet.GroupSet {
					groupSetMap := map[string]interface{}{}
					if groupSet.Id != nil {
						groupSetMap["id"] = groupSet.Id
					}

					if groupSet.Name != nil {
						groupSetMap["name"] = groupSet.Name
					}

					departmentMap := map[string]interface{}{}
					if groupSet.Department != nil {
						if groupSet.Department.Id != nil {
							departmentMap["id"] = groupSet.Department.Id
						}

						if groupSet.Department.Name != nil {
							departmentMap["name"] = groupSet.Department.Name
						}

						if groupSet.Department.Managers != nil {
							departmentMap["managers"] = groupSet.Department.Managers
						}

						managerUsersList := make([]map[string]interface{}, 0, len(groupSet.Department.ManagerUsers))
						if groupSet.Department.ManagerUsers != nil {
							for _, managerUsers := range groupSet.Department.ManagerUsers {
								managerUsersMap := map[string]interface{}{}
								if managerUsers.ManagerId != nil {
									managerUsersMap["manager_id"] = managerUsers.ManagerId
								}

								if managerUsers.ManagerName != nil {
									managerUsersMap["manager_name"] = managerUsers.ManagerName
								}

								managerUsersList = append(managerUsersList, managerUsersMap)
							}

							departmentMap["manager_users"] = managerUsersList
						}

						groupSetMap["department"] = []interface{}{departmentMap}
					}

					if groupSet.Count != nil {
						groupSetMap["count"] = groupSet.Count
					}

					groupSetList = append(groupSetList, groupSetMap)
				}

				deviceSetMap["group_set"] = groupSetList
			}

			if deviceSet.AccountCount != nil {
				deviceSetMap["account_count"] = deviceSet.AccountCount
			}

			if deviceSet.VpcId != nil {
				deviceSetMap["vpc_id"] = deviceSet.VpcId
			}

			if deviceSet.SubnetId != nil {
				deviceSetMap["subnet_id"] = deviceSet.SubnetId
			}

			resourceMap := map[string]interface{}{}
			if deviceSet.Resource != nil {
				if deviceSet.Resource.ResourceId != nil {
					resourceMap["resource_id"] = deviceSet.Resource.ResourceId
				}

				if deviceSet.Resource.ApCode != nil {
					resourceMap["ap_code"] = deviceSet.Resource.ApCode
				}

				if deviceSet.Resource.SvArgs != nil {
					resourceMap["sv_args"] = deviceSet.Resource.SvArgs
				}

				if deviceSet.Resource.VpcId != nil {
					resourceMap["vpc_id"] = deviceSet.Resource.VpcId
				}

				if deviceSet.Resource.Nodes != nil {
					resourceMap["nodes"] = deviceSet.Resource.Nodes
				}

				if deviceSet.Resource.RenewFlag != nil {
					resourceMap["renew_flag"] = deviceSet.Resource.RenewFlag
				}

				if deviceSet.Resource.ExpireTime != nil {
					resourceMap["expire_time"] = deviceSet.Resource.ExpireTime
				}

				if deviceSet.Resource.Status != nil {
					resourceMap["status"] = deviceSet.Resource.Status
				}

				if deviceSet.Resource.ResourceName != nil {
					resourceMap["resource_name"] = deviceSet.Resource.ResourceName
				}

				if deviceSet.Resource.Pid != nil {
					resourceMap["pid"] = deviceSet.Resource.Pid
				}

				if deviceSet.Resource.CreateTime != nil {
					resourceMap["create_time"] = deviceSet.Resource.CreateTime
				}

				if deviceSet.Resource.ProductCode != nil {
					resourceMap["product_code"] = deviceSet.Resource.ProductCode
				}

				if deviceSet.Resource.SubProductCode != nil {
					resourceMap["sub_product_code"] = deviceSet.Resource.SubProductCode
				}

				if deviceSet.Resource.Zone != nil {
					resourceMap["zone"] = deviceSet.Resource.Zone
				}

				if deviceSet.Resource.Expired != nil {
					resourceMap["expired"] = deviceSet.Resource.Expired
				}

				if deviceSet.Resource.Deployed != nil {
					resourceMap["deployed"] = deviceSet.Resource.Deployed
				}

				if deviceSet.Resource.VpcName != nil {
					resourceMap["vpc_name"] = deviceSet.Resource.VpcName
				}

				if deviceSet.Resource.VpcCidrBlock != nil {
					resourceMap["vpc_cidr_block"] = deviceSet.Resource.VpcCidrBlock
				}

				if deviceSet.Resource.SubnetId != nil {
					resourceMap["subnet_id"] = deviceSet.Resource.SubnetId
				}

				if deviceSet.Resource.SubnetName != nil {
					resourceMap["subnet_name"] = deviceSet.Resource.SubnetName
				}

				if deviceSet.Resource.CidrBlock != nil {
					resourceMap["cidr_block"] = deviceSet.Resource.CidrBlock
				}

				if deviceSet.Resource.PublicIpSet != nil {
					resourceMap["public_ip_set"] = deviceSet.Resource.PublicIpSet
				}

				if deviceSet.Resource.PrivateIpSet != nil {
					resourceMap["private_ip_set"] = deviceSet.Resource.PrivateIpSet
				}

				if deviceSet.Resource.ModuleSet != nil {
					resourceMap["module_set"] = deviceSet.Resource.ModuleSet
				}

				if deviceSet.Resource.UsedNodes != nil {
					resourceMap["used_nodes"] = deviceSet.Resource.UsedNodes
				}

				if deviceSet.Resource.ExtendPoints != nil {
					resourceMap["extend_points"] = deviceSet.Resource.ExtendPoints
				}

				if deviceSet.Resource.PackageBandwidth != nil {
					resourceMap["package_bandwidth"] = deviceSet.Resource.PackageBandwidth
				}

				if deviceSet.Resource.PackageNode != nil {
					resourceMap["package_node"] = deviceSet.Resource.PackageNode
				}

				if deviceSet.Resource.LogDeliveryArgs != nil {
					resourceMap["log_delivery_args"] = deviceSet.Resource.LogDeliveryArgs
				}

				clbSetList := make([]map[string]interface{}, 0, len(deviceSet.Resource.ClbSet))
				if deviceSet.Resource.ClbSet != nil {
					for _, clbSet := range deviceSet.Resource.ClbSet {
						clbSetMap := map[string]interface{}{}
						if clbSet.ClbIp != nil {
							clbSetMap["clb_ip"] = clbSet.ClbIp
						}

						clbSetList = append(clbSetList, clbSetMap)
					}

					resourceMap["clb_set"] = clbSetList
				}

				if deviceSet.Resource.DomainCount != nil {
					resourceMap["domain_count"] = deviceSet.Resource.DomainCount
				}

				if deviceSet.Resource.UsedDomainCount != nil {
					resourceMap["used_domain_count"] = deviceSet.Resource.UsedDomainCount
				}

				if deviceSet.Resource.Trial != nil {
					resourceMap["trial"] = deviceSet.Resource.Trial
				}

				if deviceSet.Resource.LogDelivery != nil {
					resourceMap["log_delivery"] = deviceSet.Resource.LogDelivery
				}

				if deviceSet.Resource.CdcClusterId != nil {
					resourceMap["cdc_cluster_id"] = deviceSet.Resource.CdcClusterId
				}

				if deviceSet.Resource.DeployModel != nil {
					resourceMap["deploy_model"] = deviceSet.Resource.DeployModel
				}

				if deviceSet.Resource.IntranetAccess != nil {
					resourceMap["intranet_access"] = deviceSet.Resource.IntranetAccess
				}

				if deviceSet.Resource.IntranetPrivateIpSet != nil {
					resourceMap["intranet_private_ip_set"] = deviceSet.Resource.IntranetPrivateIpSet
				}

				if deviceSet.Resource.IntranetVpcId != nil {
					resourceMap["intranet_vpc_id"] = deviceSet.Resource.IntranetVpcId
				}

				if deviceSet.Resource.IntranetSubnetId != nil {
					resourceMap["intranet_subnet_id"] = deviceSet.Resource.IntranetSubnetId
				}

				if deviceSet.Resource.IntranetVpcCidr != nil {
					resourceMap["intranet_vpc_cidr"] = deviceSet.Resource.IntranetVpcCidr
				}

				if deviceSet.Resource.DomainName != nil {
					resourceMap["domain_name"] = deviceSet.Resource.DomainName
				}

				if deviceSet.Resource.ShareClb != nil {
					resourceMap["share_clb"] = deviceSet.Resource.ShareClb
				}

				if deviceSet.Resource.OpenClbId != nil {
					resourceMap["open_clb_id"] = deviceSet.Resource.OpenClbId
				}

				if deviceSet.Resource.LbVipIsp != nil {
					resourceMap["lb_vip_isp"] = deviceSet.Resource.LbVipIsp
				}

				if deviceSet.Resource.TUICmdPort != nil {
					resourceMap["tui_cmd_port"] = deviceSet.Resource.TUICmdPort
				}

				if deviceSet.Resource.TUIDirectPort != nil {
					resourceMap["tui_direct_port"] = deviceSet.Resource.TUIDirectPort
				}

				if deviceSet.Resource.WebAccess != nil {
					resourceMap["web_access"] = deviceSet.Resource.WebAccess
				}

				if deviceSet.Resource.ClientAccess != nil {
					resourceMap["client_access"] = deviceSet.Resource.ClientAccess
				}

				if deviceSet.Resource.ExternalAccess != nil {
					resourceMap["external_access"] = deviceSet.Resource.ExternalAccess
				}

				if deviceSet.Resource.IOAResource != nil {
					resourceMap["ioa_resource"] = deviceSet.Resource.IOAResource
				}

				if deviceSet.Resource.PackageIOAUserCount != nil {
					resourceMap["package_ioa_user_count"] = deviceSet.Resource.PackageIOAUserCount
				}

				if deviceSet.Resource.PackageIOABandwidth != nil {
					resourceMap["package_ioa_bandwidth"] = deviceSet.Resource.PackageIOABandwidth
				}

				if deviceSet.Resource.IOAResourceId != nil {
					resourceMap["ioa_resource_id"] = deviceSet.Resource.IOAResourceId
				}

				deviceSetMap["resource"] = []interface{}{resourceMap}
			}

			departmentMap := map[string]interface{}{}
			if deviceSet.Department != nil {
				if deviceSet.Department.Id != nil {
					departmentMap["id"] = deviceSet.Department.Id
				}

				if deviceSet.Department.Name != nil {
					departmentMap["name"] = deviceSet.Department.Name
				}

				if deviceSet.Department.Managers != nil {
					departmentMap["managers"] = deviceSet.Department.Managers
				}

				managerUsersList := make([]map[string]interface{}, 0, len(deviceSet.Department.ManagerUsers))
				if deviceSet.Department.ManagerUsers != nil {
					for _, managerUsers := range deviceSet.Department.ManagerUsers {
						managerUsersMap := map[string]interface{}{}
						if managerUsers.ManagerId != nil {
							managerUsersMap["manager_id"] = managerUsers.ManagerId
						}

						if managerUsers.ManagerName != nil {
							managerUsersMap["manager_name"] = managerUsers.ManagerName
						}

						managerUsersList = append(managerUsersList, managerUsersMap)
					}

					departmentMap["manager_users"] = managerUsersList
				}

				deviceSetMap["department"] = []interface{}{departmentMap}
			}

			if deviceSet.IpPortSet != nil {
				deviceSetMap["ip_port_set"] = deviceSet.IpPortSet
			}

			if deviceSet.DomainId != nil {
				deviceSetMap["domain_id"] = deviceSet.DomainId
			}

			if deviceSet.DomainName != nil {
				deviceSetMap["domain_name"] = deviceSet.DomainName
			}

			if deviceSet.EnableSSL != nil {
				deviceSetMap["enable_ssl"] = deviceSet.EnableSSL
			}

			if deviceSet.SSLCertName != nil {
				deviceSetMap["ssl_cert_name"] = deviceSet.SSLCertName
			}

			if deviceSet.IOAId != nil {
				deviceSetMap["ioa_id"] = deviceSet.IOAId
			}

			if deviceSet.ManageDimension != nil {
				deviceSetMap["manage_dimension"] = deviceSet.ManageDimension
			}

			if deviceSet.ManageAccountId != nil {
				deviceSetMap["manage_account_id"] = deviceSet.ManageAccountId
			}

			if deviceSet.Namespace != nil {
				deviceSetMap["namespace"] = deviceSet.Namespace
			}

			if deviceSet.Workload != nil {
				deviceSetMap["workload"] = deviceSet.Workload
			}

			if deviceSet.SyncPodCount != nil {
				deviceSetMap["sync_pod_count"] = deviceSet.SyncPodCount
			}

			if deviceSet.TotalPodCount != nil {
				deviceSetMap["total_pod_count"] = deviceSet.TotalPodCount
			}

			if deviceSet.CloudAccountId != nil {
				deviceSetMap["cloud_account_id"] = deviceSet.CloudAccountId
			}

			if deviceSet.CloudAccountName != nil {
				deviceSetMap["cloud_account_name"] = deviceSet.CloudAccountName
			}

			if deviceSet.ProviderType != nil {
				deviceSetMap["provider_type"] = deviceSet.ProviderType
			}

			if deviceSet.ProviderName != nil {
				deviceSetMap["provider_name"] = deviceSet.ProviderName
			}

			if deviceSet.SyncCloudDeviceStatus != nil {
				deviceSetMap["sync_cloud_device_status"] = deviceSet.SyncCloudDeviceStatus
			}

			deviceSetList = append(deviceSetList, deviceSetMap)
		}

		_ = d.Set("device_set", deviceSetList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), deviceSetList); e != nil {
			return e
		}
	}

	return nil
}
