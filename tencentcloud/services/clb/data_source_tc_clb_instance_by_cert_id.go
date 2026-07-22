package clb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClbInstanceByCertId() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbInstanceByCertIdRead,
		Schema: map[string]*schema.Schema{
			"cert_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "服务器或客户端证书 ID。",
			},

			"cert_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "证书 ID 以及与其关联的 CLB 实例列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "证书编号。",
						},
						"load_balancers": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "与证书关联的 CLB 实例的列表。注意：该字段可能返回null，表示取不到有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"load_balancer_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB实例ID。",
									},
									"load_balancer_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB 实例名称。",
									},
									"load_balancer_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB实例网络类型：OPEN：公网；内部：专用网络。",
									},
									"forward": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "CLB 类型标识符。取值范围：1（CLB）； 0（经典 CLB）。",
									},
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB实例的域名。它仅适用于公共经典 CLB。该参数很快就会停止。请改用 LoadBalancerDomain。注意：该字段可能返回null，表示取不到有效值。",
									},
									"load_balancer_vips": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "CLB实例的VIP列表。注意：该字段可能返回null，表示取不到有效值。",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "CLB实例状态，包括：0：正在创建； 1：跑步。注意：该字段可能返回null，表示取不到有效值。",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB实例创建时间。注意：该字段可能返回null，表示取不到有效值。",
									},
									"status_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB实例最后一次状态变化时间。注意：该字段可能返回null，表示取不到有效值。",
									},
									"project_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "CLB实例所属项目ID。 0：默认项目。",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "私有网络 ID 注意：该字段可能返回null，表示取不到有效值。",
									},
									"open_bgp": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "保护性 CLB 标识符。取值范围：1（保护）、0（不保护）。注意：该字段可能返回null，表示取不到有效值。",
									},
									"snat": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "2016年12月之前创建的所有私网经典CLB均已启用SNAT。 注意：该字段可能返回null，表示取不到有效值。",
									},
									"isolation": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "0：不隔离； 1：孤立。注意：该字段可能返回null，表示取不到有效值。",
									},
									"log": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "日志信息。只有具有HTTP或HTTPS监听的公网CLB才能生成日志。注意：该字段可能返回null，表示取不到有效值。",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB实例所在子网（仅对私网VPC CLB有意义）。注意：该字段可能返回null，表示取不到有效值。",
									},
									"tags": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "CLB实例标签信息。注意：该字段可能返回null，表示取不到有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"tag_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "标签键。",
												},
												"tag_value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "标签值。",
												},
											},
										},
									},
									"secure_groups": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "CLB实例的安全组。注意：该字段可能返回null，表示取不到有效值。",
									},
									"target_region_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "CLB实例绑定的后端服务器基本信息。注意：该字段可能返回null，表示取不到有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "目标区域，例如ap-guangzhou。",
												},
												"vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "目标网络，VPC格式为vpc-abcd1234，基础网络格式为0。",
												},
											},
										},
									},
									"anycast_zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Anycast CLB 发布区域。对于非选播 CLB，此字段返回空字符串。注意：该字段可能返回null，表示取不到有效值。",
									},
									"address_i_p_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IP版本。有效值：ipv4、ipv6。注意：该字段可能返回null，表示取不到有效值。",
									},
									"numerical_vpc_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数字形式的 私有网络 ID。注意：该字段可能返回null，表示取不到有效值。",
									},
									"vip_isp": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB IP地址所属的ISP。注意：该字段可能返回null，表示取不到有效值。",
									},
									"master_zone": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "主要AZ。注意：该字段可能返回null，表示取不到有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"zone_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "。",
												},
												"zone": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "数字形式的唯一AZ ID，如100001。 注意：该字段可能返回null，表示取不到有效值。",
												},
												"zone_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "可用区名称，如广州一区。 注意：该字段可能返回null，表示取不到有效值。",
												},
												"zone_region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "AZ区域，例如ap-广州。注意：该字段可能返回null，表示取不到有效值。",
												},
												"local_zone": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "AZ是否是LocalZone，例如，false。注意：该字段可能返回null，表示取不到有效值。",
												},
												"edge_zone": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "该AZ是否为边缘区域。值：真、假。注意：该字段可能返回null，表示取不到有效值。",
												},
											},
										},
									},
									"backup_zone_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "备份区。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"zone_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "数字形式的唯一AZ ID，如100001。 注意：该字段可能返回null，表示取不到有效值。",
												},
												"zone": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "数字形式的唯一AZ ID，如100001。 注意：该字段可能返回null，表示取不到有效值。",
												},
												"zone_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "可用区名称，如广州一区。 注意：该字段可能返回null，表示取不到有效值。",
												},
												"zone_region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "AZ区域，例如ap-广州。注意：该字段可能返回null，表示取不到有效值。",
												},
												"local_zone": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "AZ是否是LocalZone，例如，false。注意：该字段可能返回null，表示取不到有效值。",
												},
												"edge_zone": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "该AZ是否为边缘区域。值：真、假。注意：该字段可能返回null，表示取不到有效值。",
												},
											},
										},
									},
									"isolated_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB实例隔离时间。注意：该字段可能返回null，表示取不到有效值。",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB实例过期时间，仅对预付费实例有效。注意：该字段可能返回null，表示取不到有效值。",
									},
									"charge_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB实例的计费方式。有效值：PREPAID（每月订阅）、POSTPAID_BY_HOUR（按需付费）。注意：该字段可能返回null，表示取不到有效值。",
									},
									"network_attributes": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "CLB实例网络属性。注意：该字段可能返回null，表示取不到有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"internet_charge_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "TRAFFIC_POSTPAID_BY_HOUR：按小时按流量按量付费； BANDWIDTH_POSTPAID_BY_HOUR：按带宽按小时付费； BANDWIDTH_PACKAGE：按带宽套餐计费（目前只有指定ISP才支持此方式）。",
												},
												"internet_max_bandwidth_out": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "最大出方向带宽（Mbps），仅适用于公网CLB。值范围：0-65,535。默认值：10。",
												},
												"bandwidthpkg_sub_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "带宽封装类型，如SINGLEISP。注意：该字段可能返回null，表示取不到有效值。",
												},
											},
										},
									},
									"prepaid_attributes": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "CLB实例的预付费计费属性。注意：该字段可能返回null，表示取不到有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"renew_flag": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "续订类型。 AUTO_RENEW：自动续订； MANUAL_RENEW：手动续订。注意：该字段可能返回null，表示取不到有效值。",
												},
												"period": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Cycle，表示月数（保留字段）。注意：该字段可能返回null，表示取不到有效值。",
												},
											},
										},
									},
									"log_set_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB日志服务（CLS）的日志集ID。注意：该字段可能返回null，表示取不到有效值。",
									},
									"log_topic_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB日志服务（CLS）的日志主题ID。注意：该字段可能返回null，表示取不到有效值。",
									},
									"address_i_pv6": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB实例的IPv6地址。注意：该字段可能返回null，表示取不到有效值。",
									},
									"extra_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "保留字段，一般可以忽略。注意：该字段可能返回null，表示取不到有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"zhi_tong": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否启用VIP直连。注意：该字段可能返回null，表示取不到有效值。",
												},
												"tgw_group_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Tgw 组名称。注意：该字段可能返回null，表示取不到有效值。",
												},
											},
										},
									},
									"is_ddos": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否可以绑定DDoS高防实例。注意：该字段可能返回null，表示取不到有效值。",
									},
									"config_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB 实例级别的自定义配置 ID。注意：该字段可能返回null，表示取不到有效值。",
									},
									"load_balancer_pass_to_target": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "真实服务器是否开放CLB实例到互联网的流量。注意：该字段可能返回null，表示取不到有效值。",
									},
									"exclusive_cluster": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "专网专用集群。注意：该字段可能返回null，表示取不到有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"l4_clusters": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "第 4 层专用集群列表。注意：该字段可能返回null，表示取不到有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cluster_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "唯一的集群 ID。",
															},
															"cluster_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "集群名称。注意：该字段可能返回null，表示取不到有效值。",
															},
															"zone": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "集群AZ，例如ap-guangzhou-1。注意：该字段可能返回null，表示取不到有效值。",
															},
														},
													},
												},
												"l7_clusters": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "第 7 层专用集群列表。注意：该字段可能返回null，表示取不到有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cluster_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "唯一的集群 ID。",
															},
															"cluster_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "集群名称。注意：该字段可能返回null，表示取不到有效值。",
															},
															"zone": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "集群AZ，例如ap-guangzhou-1。注意：该字段可能返回null，表示取不到有效值。",
															},
														},
													},
												},
												"classical_cluster": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "vpcgw 集群。注意：该字段可能返回null，表示取不到有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cluster_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "唯一的集群 ID。",
															},
															"cluster_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "集群名称。注意：该字段可能返回null，表示取不到有效值。",
															},
															"zone": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "集群AZ，例如ap-guangzhou-1。注意：该字段可能返回null，表示取不到有效值。",
															},
														},
													},
												},
											},
										},
									},
									"ipv6_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "仅当IP地址版本为ipv6时，该字段才有意义。有效值：IPv6Nat64、IPv6FullChain。注意：该字段可能返回null，表示取不到有效值。",
									},
									"snat_pro": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否启用SnatPro。注意：该字段可能返回null，表示取不到有效值。",
									},
									"snat_ips": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "启用S​​natPro负载均衡后的SnatIp列表。注意：该字段可能返回null，表示取不到有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"subnet_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "唯一的 VPC 子网 ID，例如subnet-12345678。",
												},
												"ip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "IP 地址，例如 192.168.0.1。",
												},
											},
										},
									},
									"sla_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "LCU支持的实例规格。注意：该字段可能返回null，表示取不到有效值。",
									},
									"is_block": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "VIP是否被屏蔽。注意：该字段可能返回null，表示取不到有效值。",
									},
									"is_block_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "时间被封锁或被解除。注意：该字段可能返回null，表示取不到有效值。",
									},
									"local_bgp": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "IP类型是否为本地BGP。",
									},
									"cluster_tag": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "专用第 7 层标签。注意：该字段可能返回null，表示取不到有效值。",
									},
									"mix_ip_target": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "如果启用了 IPv6FullChain CLB 实例的七层监听，则该 CLB 实例可以同时绑定 IPv4 和 IPv6 CVM 实例。注意：该字段可能返回null，表示取不到有效值。",
									},
									"zones": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "基于VPC的私网CLB实例的可用区。注意：该字段可能返回null，表示取不到有效值。",
									},
									"nfv_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "是否是NFV CLB实例。无返回信息：无； l7nfv：是的。注意：该字段可能返回null，表示取不到有效值。",
									},
									"health_log_set_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB CLS的健康检查日志集ID。注意：该字段可能返回null，表示取不到有效值。",
									},
									"health_log_topic_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB CLS的健康检查日志主题ID。注意：该字段可能返回null，表示取不到有效值。",
									},
									"cluster_ids": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "集群 ID。注意：该字段可能返回null，表示取不到有效值。",
									},
									"attribute_flags": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "集群 ID。注意：该字段可能返回null，表示取不到有效值。",
									},
									"load_balancer_domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB实例的域名。注意：该字段可能返回null，表示取不到有效值。",
									},
								},
							},
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

func dataSourceTencentCloudClbInstanceByCertIdRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clb_instance_by_cert_id.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cert_ids"); ok {
		certIdsSet := v.(*schema.Set).List()
		paramMap["CertIds"] = helper.InterfacesStringsPoint(certIdsSet)
	}

	service := ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var certSet []*clb.CertIdRelatedWithLoadBalancers

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClbInstanceByCertId(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		certSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(certSet))
	tmpList := make([]map[string]interface{}, 0, len(certSet))

	if certSet != nil {
		for _, certIdRelatedWithLoadBalancers := range certSet {
			certIdRelatedWithLoadBalancersMap := map[string]interface{}{}

			if certIdRelatedWithLoadBalancers.CertId != nil {
				certIdRelatedWithLoadBalancersMap["cert_id"] = certIdRelatedWithLoadBalancers.CertId
			}

			if certIdRelatedWithLoadBalancers.LoadBalancers != nil {
				loadBalancersList := []interface{}{}
				for _, loadBalancers := range certIdRelatedWithLoadBalancers.LoadBalancers {
					loadBalancersMap := map[string]interface{}{}

					if loadBalancers.LoadBalancerId != nil {
						loadBalancersMap["load_balancer_id"] = loadBalancers.LoadBalancerId
					}

					if loadBalancers.LoadBalancerName != nil {
						loadBalancersMap["load_balancer_name"] = loadBalancers.LoadBalancerName
					}

					if loadBalancers.LoadBalancerType != nil {
						loadBalancersMap["load_balancer_type"] = loadBalancers.LoadBalancerType
					}

					if loadBalancers.Forward != nil {
						loadBalancersMap["forward"] = loadBalancers.Forward
					}

					if loadBalancers.Domain != nil {
						loadBalancersMap["domain"] = loadBalancers.Domain
					}

					if loadBalancers.LoadBalancerVips != nil {
						loadBalancersMap["load_balancer_vips"] = loadBalancers.LoadBalancerVips
					}

					if loadBalancers.Status != nil {
						loadBalancersMap["status"] = loadBalancers.Status
					}

					if loadBalancers.CreateTime != nil {
						loadBalancersMap["create_time"] = loadBalancers.CreateTime
					}

					if loadBalancers.StatusTime != nil {
						loadBalancersMap["status_time"] = loadBalancers.StatusTime
					}

					if loadBalancers.ProjectId != nil {
						loadBalancersMap["project_id"] = loadBalancers.ProjectId
					}

					if loadBalancers.VpcId != nil {
						loadBalancersMap["vpc_id"] = loadBalancers.VpcId
					}

					if loadBalancers.OpenBgp != nil {
						loadBalancersMap["open_bgp"] = loadBalancers.OpenBgp
					}

					if loadBalancers.Snat != nil {
						loadBalancersMap["snat"] = loadBalancers.Snat
					}

					if loadBalancers.Isolation != nil {
						loadBalancersMap["isolation"] = loadBalancers.Isolation
					}

					if loadBalancers.Log != nil {
						loadBalancersMap["log"] = loadBalancers.Log
					}

					if loadBalancers.SubnetId != nil {
						loadBalancersMap["subnet_id"] = loadBalancers.SubnetId
					}

					if loadBalancers.Tags != nil {
						tagsList := []interface{}{}
						for _, tags := range loadBalancers.Tags {
							tagsMap := map[string]interface{}{}

							if tags.TagKey != nil {
								tagsMap["tag_key"] = tags.TagKey
							}

							if tags.TagValue != nil {
								tagsMap["tag_value"] = tags.TagValue
							}

							tagsList = append(tagsList, tagsMap)
						}

						loadBalancersMap["tags"] = tagsList
					}

					if loadBalancers.SecureGroups != nil {
						loadBalancersMap["secure_groups"] = loadBalancers.SecureGroups
					}

					if loadBalancers.TargetRegionInfo != nil {
						targetRegionInfoMap := map[string]interface{}{}

						if loadBalancers.TargetRegionInfo.Region != nil {
							targetRegionInfoMap["region"] = loadBalancers.TargetRegionInfo.Region
						}

						if loadBalancers.TargetRegionInfo.VpcId != nil {
							targetRegionInfoMap["vpc_id"] = loadBalancers.TargetRegionInfo.VpcId
						}

						loadBalancersMap["target_region_info"] = []interface{}{targetRegionInfoMap}
					}

					if loadBalancers.AnycastZone != nil {
						loadBalancersMap["anycast_zone"] = loadBalancers.AnycastZone
					}

					if loadBalancers.AddressIPVersion != nil {
						loadBalancersMap["address_i_p_version"] = loadBalancers.AddressIPVersion
					}

					if loadBalancers.NumericalVpcId != nil {
						loadBalancersMap["numerical_vpc_id"] = loadBalancers.NumericalVpcId
					}

					if loadBalancers.VipIsp != nil {
						loadBalancersMap["vip_isp"] = loadBalancers.VipIsp
					}

					if loadBalancers.MasterZone != nil {
						masterZoneMap := map[string]interface{}{}

						if loadBalancers.MasterZone.ZoneId != nil {
							masterZoneMap["zone_id"] = loadBalancers.MasterZone.ZoneId
						}

						if loadBalancers.MasterZone.Zone != nil {
							masterZoneMap["zone"] = loadBalancers.MasterZone.Zone
						}

						if loadBalancers.MasterZone.ZoneName != nil {
							masterZoneMap["zone_name"] = loadBalancers.MasterZone.ZoneName
						}

						if loadBalancers.MasterZone.ZoneRegion != nil {
							masterZoneMap["zone_region"] = loadBalancers.MasterZone.ZoneRegion
						}

						if loadBalancers.MasterZone.LocalZone != nil {
							masterZoneMap["local_zone"] = loadBalancers.MasterZone.LocalZone
						}

						if loadBalancers.MasterZone.EdgeZone != nil {
							masterZoneMap["edge_zone"] = loadBalancers.MasterZone.EdgeZone
						}

						loadBalancersMap["master_zone"] = []interface{}{masterZoneMap}
					}

					if loadBalancers.BackupZoneSet != nil {
						backupZoneSetList := []interface{}{}
						for _, backupZoneSet := range loadBalancers.BackupZoneSet {
							backupZoneSetMap := map[string]interface{}{}

							if backupZoneSet.ZoneId != nil {
								backupZoneSetMap["zone_id"] = backupZoneSet.ZoneId
							}

							if backupZoneSet.Zone != nil {
								backupZoneSetMap["zone"] = backupZoneSet.Zone
							}

							if backupZoneSet.ZoneName != nil {
								backupZoneSetMap["zone_name"] = backupZoneSet.ZoneName
							}

							if backupZoneSet.ZoneRegion != nil {
								backupZoneSetMap["zone_region"] = backupZoneSet.ZoneRegion
							}

							if backupZoneSet.LocalZone != nil {
								backupZoneSetMap["local_zone"] = backupZoneSet.LocalZone
							}

							if backupZoneSet.EdgeZone != nil {
								backupZoneSetMap["edge_zone"] = backupZoneSet.EdgeZone
							}

							backupZoneSetList = append(backupZoneSetList, backupZoneSetMap)
						}

						loadBalancersMap["backup_zone_set"] = backupZoneSetList
					}

					if loadBalancers.IsolatedTime != nil {
						loadBalancersMap["isolated_time"] = loadBalancers.IsolatedTime
					}

					if loadBalancers.ExpireTime != nil {
						loadBalancersMap["expire_time"] = loadBalancers.ExpireTime
					}

					if loadBalancers.ChargeType != nil {
						loadBalancersMap["charge_type"] = loadBalancers.ChargeType
					}

					if loadBalancers.NetworkAttributes != nil {
						networkAttributesMap := map[string]interface{}{}

						if loadBalancers.NetworkAttributes.InternetChargeType != nil {
							networkAttributesMap["internet_charge_type"] = loadBalancers.NetworkAttributes.InternetChargeType
						}

						if loadBalancers.NetworkAttributes.InternetMaxBandwidthOut != nil {
							networkAttributesMap["internet_max_bandwidth_out"] = loadBalancers.NetworkAttributes.InternetMaxBandwidthOut
						}

						if loadBalancers.NetworkAttributes.BandwidthpkgSubType != nil {
							networkAttributesMap["bandwidthpkg_sub_type"] = loadBalancers.NetworkAttributes.BandwidthpkgSubType
						}

						loadBalancersMap["network_attributes"] = []interface{}{networkAttributesMap}
					}

					if loadBalancers.PrepaidAttributes != nil {
						prepaidAttributesMap := map[string]interface{}{}

						if loadBalancers.PrepaidAttributes.RenewFlag != nil {
							prepaidAttributesMap["renew_flag"] = loadBalancers.PrepaidAttributes.RenewFlag
						}

						if loadBalancers.PrepaidAttributes.Period != nil {
							prepaidAttributesMap["period"] = loadBalancers.PrepaidAttributes.Period
						}

						loadBalancersMap["prepaid_attributes"] = []interface{}{prepaidAttributesMap}
					}

					if loadBalancers.LogSetId != nil {
						loadBalancersMap["log_set_id"] = loadBalancers.LogSetId
					}

					if loadBalancers.LogTopicId != nil {
						loadBalancersMap["log_topic_id"] = loadBalancers.LogTopicId
					}

					if loadBalancers.AddressIPv6 != nil {
						loadBalancersMap["address_i_pv6"] = loadBalancers.AddressIPv6
					}

					if loadBalancers.ExtraInfo != nil {
						extraInfoMap := map[string]interface{}{}

						if loadBalancers.ExtraInfo.ZhiTong != nil {
							extraInfoMap["zhi_tong"] = loadBalancers.ExtraInfo.ZhiTong
						}

						if loadBalancers.ExtraInfo.TgwGroupName != nil {
							extraInfoMap["tgw_group_name"] = loadBalancers.ExtraInfo.TgwGroupName
						}

						loadBalancersMap["extra_info"] = []interface{}{extraInfoMap}
					}

					if loadBalancers.IsDDos != nil {
						loadBalancersMap["is_ddos"] = loadBalancers.IsDDos
					}

					if loadBalancers.ConfigId != nil {
						loadBalancersMap["config_id"] = loadBalancers.ConfigId
					}

					if loadBalancers.LoadBalancerPassToTarget != nil {
						loadBalancersMap["load_balancer_pass_to_target"] = loadBalancers.LoadBalancerPassToTarget
					}

					if loadBalancers.ExclusiveCluster != nil {
						exclusiveClusterMap := map[string]interface{}{}

						if loadBalancers.ExclusiveCluster.L4Clusters != nil {
							l4ClustersList := []interface{}{}
							for _, l4Clusters := range loadBalancers.ExclusiveCluster.L4Clusters {
								l4ClustersMap := map[string]interface{}{}

								if l4Clusters.ClusterId != nil {
									l4ClustersMap["cluster_id"] = l4Clusters.ClusterId
								}

								if l4Clusters.ClusterName != nil {
									l4ClustersMap["cluster_name"] = l4Clusters.ClusterName
								}

								if l4Clusters.Zone != nil {
									l4ClustersMap["zone"] = l4Clusters.Zone
								}

								l4ClustersList = append(l4ClustersList, l4ClustersMap)
							}

							exclusiveClusterMap["l4_clusters"] = l4ClustersList
						}

						if loadBalancers.ExclusiveCluster.L7Clusters != nil {
							l7ClustersList := []interface{}{}
							for _, l7Clusters := range loadBalancers.ExclusiveCluster.L7Clusters {
								l7ClustersMap := map[string]interface{}{}

								if l7Clusters.ClusterId != nil {
									l7ClustersMap["cluster_id"] = l7Clusters.ClusterId
								}

								if l7Clusters.ClusterName != nil {
									l7ClustersMap["cluster_name"] = l7Clusters.ClusterName
								}

								if l7Clusters.Zone != nil {
									l7ClustersMap["zone"] = l7Clusters.Zone
								}

								l7ClustersList = append(l7ClustersList, l7ClustersMap)
							}

							exclusiveClusterMap["l7_clusters"] = l7ClustersList
						}

						if loadBalancers.ExclusiveCluster.ClassicalCluster != nil {
							classicalClusterMap := map[string]interface{}{}

							if loadBalancers.ExclusiveCluster.ClassicalCluster.ClusterId != nil {
								classicalClusterMap["cluster_id"] = loadBalancers.ExclusiveCluster.ClassicalCluster.ClusterId
							}

							if loadBalancers.ExclusiveCluster.ClassicalCluster.ClusterName != nil {
								classicalClusterMap["cluster_name"] = loadBalancers.ExclusiveCluster.ClassicalCluster.ClusterName
							}

							if loadBalancers.ExclusiveCluster.ClassicalCluster.Zone != nil {
								classicalClusterMap["zone"] = loadBalancers.ExclusiveCluster.ClassicalCluster.Zone
							}

							exclusiveClusterMap["classical_cluster"] = []interface{}{classicalClusterMap}
						}

						loadBalancersMap["exclusive_cluster"] = []interface{}{exclusiveClusterMap}
					}

					if loadBalancers.IPv6Mode != nil {
						loadBalancersMap["ipv6_mode"] = loadBalancers.IPv6Mode
					}

					if loadBalancers.SnatPro != nil {
						loadBalancersMap["snat_pro"] = loadBalancers.SnatPro
					}

					if loadBalancers.SnatIps != nil {
						snatIpsList := []interface{}{}
						for _, snatIps := range loadBalancers.SnatIps {
							snatIpsMap := map[string]interface{}{}

							if snatIps.SubnetId != nil {
								snatIpsMap["subnet_id"] = snatIps.SubnetId
							}

							if snatIps.Ip != nil {
								snatIpsMap["ip"] = snatIps.Ip
							}

							snatIpsList = append(snatIpsList, snatIpsMap)
						}

						loadBalancersMap["snat_ips"] = snatIpsList
					}

					if loadBalancers.SlaType != nil {
						loadBalancersMap["sla_type"] = loadBalancers.SlaType
					}

					if loadBalancers.IsBlock != nil {
						loadBalancersMap["is_block"] = loadBalancers.IsBlock
					}

					if loadBalancers.IsBlockTime != nil {
						loadBalancersMap["is_block_time"] = loadBalancers.IsBlockTime
					}

					if loadBalancers.LocalBgp != nil {
						loadBalancersMap["local_bgp"] = loadBalancers.LocalBgp
					}

					if loadBalancers.ClusterTag != nil {
						loadBalancersMap["cluster_tag"] = loadBalancers.ClusterTag
					}

					if loadBalancers.MixIpTarget != nil {
						loadBalancersMap["mix_ip_target"] = loadBalancers.MixIpTarget
					}

					if loadBalancers.Zones != nil {
						loadBalancersMap["zones"] = loadBalancers.Zones
					}

					if loadBalancers.NfvInfo != nil {
						loadBalancersMap["nfv_info"] = loadBalancers.NfvInfo
					}

					if loadBalancers.HealthLogSetId != nil {
						loadBalancersMap["health_log_set_id"] = loadBalancers.HealthLogSetId
					}

					if loadBalancers.HealthLogTopicId != nil {
						loadBalancersMap["health_log_topic_id"] = loadBalancers.HealthLogTopicId
					}

					if loadBalancers.ClusterIds != nil {
						loadBalancersMap["cluster_ids"] = loadBalancers.ClusterIds
					}

					if loadBalancers.AttributeFlags != nil {
						loadBalancersMap["attribute_flags"] = loadBalancers.AttributeFlags
					}

					if loadBalancers.LoadBalancerDomain != nil {
						loadBalancersMap["load_balancer_domain"] = loadBalancers.LoadBalancerDomain
					}

					loadBalancersList = append(loadBalancersList, loadBalancersMap)
				}

				certIdRelatedWithLoadBalancersMap["load_balancers"] = loadBalancersList
			}

			ids = append(ids, *certIdRelatedWithLoadBalancers.CertId)
			tmpList = append(tmpList, certIdRelatedWithLoadBalancersMap)
		}

		_ = d.Set("cert_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
