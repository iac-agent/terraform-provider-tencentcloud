package clb

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcas "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/as"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var clbActionMu = &sync.Mutex{}

func ResourceTencentCloudClbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbInstanceCreate,
		Read:   resourceTencentCloudClbInstanceRead,
		Update: resourceTencentCloudClbInstanceUpdate,
		Delete: resourceTencentCloudClbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"network_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CLB_NETWORK_TYPE),
				Description: "CLB实例类型。有效值：“OPEN”和“INTERNAL”。",
			},
			"clb_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description: "CLB 的名称。名称只能包含汉字、英文字母、数字、下划线和连字符“-”。",
			},
			"clb_vips": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "CLB的虚拟服务地址表。",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "CLB 实例中项目的 ID，“0”- 默认项目。",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "CLB的VPC ID。",
			},
			"subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(2, 60),
				Description: "如果购买“INTERNAL”clb 实例，则必须指定子网 ID。 “INTERNAL” clb 实例的 VIP 将从该子网生成。",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "集群 ID。",
			},
			"address_ip_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "仅适用于公网CLB实例。 IP版本。值：“IPV4”、“IPV6”和“IPv6FullChain”（不区分大小写）。默认值：“IPV4”。注：IPV6 表示 IPv6 NAT64，而 IPv6FullChain 表示 IPv6。",
			},
			"internet_charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "网费类型，仅适用于开通CLB。有效值为“TRAFFIC_POSTPAID_BY_HOUR”、“BANDWIDTH_POSTPAID_BY_HOUR”和“BANDWIDTH_PACKAGE”。",
			},
			"delete_protect": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "是否启用删除保护。",
			},
			"bandwidth_package_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "带宽包 ID。如果设置，“internet_charge_type”必须是“BANDWIDTH_PACKAGE”。",
			},
			"internet_bandwidth_max_out": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "最大带宽输出，仅适用于开放的CLB。有效值范围为 [1，2048]。单位为Mbps。",
			},
			"security_groups": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "CLB实例的安全组。支持“OPEN”和“INTERNAL”CLB。",
			},
			"target_region_info_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "后端服务的地域信息附加在CLB实例上。仅支持“OPEN”CLB。",
			},
			"target_region_info_vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "后端服务的VPC信息附着在CLB实例上。仅支持“OPEN”CLB。",
			},
			"snat_pro": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "其他VPC的绑定IP是否切换。",
			},
			"snat_ips": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Snat Ip 列表，需要 `snat_pro=true`。注意：这里无法读取和修改该参数，因为动态ip无法追踪，请导入资源“tencentcloud_clb_snat_ip”来处理固定ip。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Snat IP 地址，如果设置为空将自动分配。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Snat 子网 ID。",
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "此 CLB 中的可用标签。",
			},
			"sla_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "创建 LCU 支持的实例需要此参数。价值观：" +
					"`SLA`: Super Large 4. When you have activated Super Large models, `SLA` refers to Super Large 4; " +
					"`clb.c2.medium`: Standard; " +
					"`clb.c3.small`: Advanced 1; " +
					"`clb.c3.medium`: Advanced 1; " +
					"`clb.c4.small`: Super Large 1; " +
					"`clb.c4.medium`: Super Large 2; " +
					"`clb.c4.large`: Super Large 3; " +
					"`clb.c4.xlarge`: Super Large 4. " +
					"For more details, see [Instance Specifications](https://intl.cloud.tencent.com/document/product/214/84689?from_cn_redirect=1).",
			},
			"vip_isp": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "网络运营商，仅适用于开通CLB。有效值为“CMCC”（中国移动）、“CTCC”（电信）、“CUCC”（中国联通）和“BGP”。如果指定了该ISP，则网络计费方式只能使用带宽套餐计费（BANDWIDTH_PACKAGE）。",
			},
			"vip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "指定CLB实例应用的VIP。该参数是可选的。如果不指定该参数，系统会自动为该参数分配一个值。 IPv4和IPv6 CLB实例支持此参数，但IPv6 NAT64 CLB实例不支持。",
			},
			"load_balancer_pass_to_target": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "目标允许流量是否来自clb。如果值为true，则仅检查clb的安全组，或者同时检查clb和后端实例安全组。",
			},
			"master_zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "设置跨可用区容灾的主域id，仅适用于开放的CLB。",
			},
			"zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "可用的zone id，仅适用于开放的CLB。",
			},
			"slave_zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "设置跨可用区容灾的从属区id，仅适用于开放的CLB。当主设备宕机时，该区域将承担流量。",
			},
			"log_set_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "日志集的id。",
			},
			"log_topic_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "日志主题的id。",
			},
			"dynamic_vip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "是否创建动态vip CLB实例，“true”或“false”。",
			},
			"eip_address_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "EIP唯一ID，例如eip-1v2rmbwk，仅适用于内网负载均衡绑定EIP。 EIP变更期间，可能会出现短暂的网络中断。",
			},
			"associate_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "关联的终端节点ID；传递空字符串表示取消关联节点。",
			},
			"exclusive_cluster": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "专用CLB实例的信息。在私网创建专用CLB实例时必须指定该参数。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"l4_clusters": {
							Type:        schema.TypeSet,
							Optional:    true,
							ForceNew:    true,
							Description: "四层专用簇列表\n注意：该字段可能返回null，表示取不到有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_id": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "唯一的集群 ID。",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "集群名称。",
									},
									"zone": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "集群AZ，如ap-guangzhou-1\n注意：该字段可能返回null，表示取不到有效值。",
									},
								},
							},
						},
						"l7_clusters": {
							Type:        schema.TypeSet,
							Optional:    true,
							ForceNew:    true,
							Description: "七层专用簇列表\n注意：该字段可能返回null，表示取不到有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_id": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "唯一的集群 ID。",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "集群名称。",
									},
									"zone": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "集群AZ，如ap-guangzhou-1\n注意：该字段可能返回null，表示取不到有效值。",
									},
								},
							},
						},
						"classical_cluster": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "vpcgw cluster\n注意：该字段可能返回null，表示取不到有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_id": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "唯一的集群 ID。",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "集群名称。",
									},
									"zone": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "集群AZ，如ap-guangzhou-1\n注意：该字段可能返回null，表示取不到有效值。",
									},
								},
							},
						},
					},
				},
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CLB实例的域名。",
			},
			"ipv6_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "当IP地址版本为ipv6，`IPv6Nat64`时，该字段有意义 | “IPv6全链”。",
			},
			"address_ipv6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "负载均衡实例的IPv6地址。",
			},
		},
	}
}

func resourceTencentCloudClbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_instance.create")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	networkType := d.Get("network_type").(string)
	clbName := d.Get("clb_name").(string)
	flag, e := checkSameName(clbName, meta)
	if e != nil {
		return e
	}

	if flag {
		return fmt.Errorf("[CHECK][CLB instance][Create] check: Same CLB name %s exists!", clbName)
	}

	targetRegionInfoRegion := ""
	targetRegionInfoVpcId := ""
	if v, ok := d.GetOk("target_region_info_region"); ok {
		targetRegionInfoRegion = v.(string)
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support this operation with target_region_info")
		}
	}

	if v, ok := d.GetOk("target_region_info_vpc_id"); ok {
		targetRegionInfoVpcId = v.(string)
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support this operation with target_region_info")
		}
	}

	if (targetRegionInfoRegion != "" && targetRegionInfoVpcId == "") || (targetRegionInfoRegion == "" && targetRegionInfoVpcId != "") {
		return fmt.Errorf("[CHECK][CLB instance][Create] check: region and vpc_id must be set at same time")
	}

	request := clb.NewCreateLoadBalancerRequest()
	request.LoadBalancerType = helper.String(networkType)
	request.LoadBalancerName = helper.String(clbName)
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		projectId := int64(v.(int))
		request.ProjectId = &projectId
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterIds = []*string{helper.String(v.(string))}
	}

	//vip_isp
	if v, ok := d.GetOk("vip_isp"); ok {
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support vip ISP setting")
		}

		request.VipIsp = helper.String(v.(string))
	}

	//vip
	if v, ok := d.GetOk("vip"); ok {
		request.Vip = helper.String(v.(string))
	}

	//SlaType
	if v, ok := d.GetOk("sla_type"); ok {
		request.SlaType = helper.String(v.(string))
	}

	//ip version
	if v, ok := d.GetOk("address_ip_version"); ok {
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support IP version setting")
		}

		request.AddressIPVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("snat_pro"); ok {
		request.SnatPro = helper.Bool(v.(bool))
	}

	if v, ok := d.Get("snat_ips").([]interface{}); ok && len(v) > 0 {
		for i := range v {
			item := v[i].(map[string]interface{})
			snatIp := &clb.SnatIp{}
			if v, ok := item["subnet_id"].(string); ok && v != "" {
				snatIp.SubnetId = &v
			}

			if v, ok := item["ip"].(string); ok && v != "" {
				snatIp.Ip = &v
			}

			request.SnatIps = append(request.SnatIps, snatIp)
		}
	}

	v, ok := d.GetOk("internet_charge_type")
	bv, bok := d.GetOk("internet_bandwidth_max_out")
	pv, pok := d.GetOk("bandwidth_package_id")

	chargeType := v.(string)

	//internet charge type
	if ok {
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support internet charge type setting")
		}
	}

	if ok || bok || pok {
		request.InternetAccessible = &clb.InternetAccessible{}
	}

	if ok {
		request.InternetAccessible.InternetChargeType = helper.String(chargeType)
	}

	if pok {
		if chargeType != svcas.INTERNET_CHARGE_TYPE_BANDWIDTH_PACKAGE {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: internet_charge_type must `BANDWIDTH_PACKAGE` when bandwidth_package_id was set")
		}

		request.BandwidthPackageId = helper.String(pv.(string))
	}

	// open or internal
	if bok {
		request.InternetAccessible.InternetMaxBandwidthOut = helper.IntInt64(bv.(int))
	}

	if v, ok := d.GetOk("master_zone_id"); ok {
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support master zone id setting")
		}

		request.MasterZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("zone_id"); ok {
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support zone id setting")
		}

		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("slave_zone_id"); ok {
		if networkType == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance][Create] check: INTERNAL network_type do not support slave zone id setting")
		}

		request.SlaveZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("load_balancer_pass_to_target"); ok {
		request.LoadBalancerPassToTarget = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("dynamic_vip"); ok {
		request.DynamicVip = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("eip_address_id"); ok {
		request.EipAddressId = helper.String(v.(string))
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		for k, v := range tags {
			tmpKey := k
			tmpValue := v
			request.Tags = append(request.Tags, &clb.TagInfo{
				TagKey:   &tmpKey,
				TagValue: &tmpValue,
			})
		}
	}

	if exclusiveClusterMap, ok := helper.InterfacesHeadMap(d, "exclusive_cluster"); ok {
		exclusiveCluster := clb.ExclusiveCluster{}
		if v, ok := exclusiveClusterMap["l4_clusters"]; ok {
			for _, item := range v.(*schema.Set).List() {
				l4ClustersMap := item.(map[string]interface{})
				clusterItem := clb.ClusterItem{}
				if v, ok := l4ClustersMap["cluster_id"].(string); ok && v != "" {
					clusterItem.ClusterId = helper.String(v)
				}

				if v, ok := l4ClustersMap["cluster_name"].(string); ok && v != "" {
					clusterItem.ClusterName = helper.String(v)
				}

				if v, ok := l4ClustersMap["zone"].(string); ok && v != "" {
					clusterItem.Zone = helper.String(v)
				}

				exclusiveCluster.L4Clusters = append(exclusiveCluster.L4Clusters, &clusterItem)
			}
		}

		if v, ok := exclusiveClusterMap["l7_clusters"]; ok {
			for _, item := range v.(*schema.Set).List() {
				l7ClustersMap := item.(map[string]interface{})
				clusterItem := clb.ClusterItem{}
				if v, ok := l7ClustersMap["cluster_id"].(string); ok && v != "" {
					clusterItem.ClusterId = helper.String(v)
				}

				if v, ok := l7ClustersMap["cluster_name"].(string); ok && v != "" {
					clusterItem.ClusterName = helper.String(v)
				}

				if v, ok := l7ClustersMap["zone"].(string); ok && v != "" {
					clusterItem.Zone = helper.String(v)
				}

				exclusiveCluster.L7Clusters = append(exclusiveCluster.L7Clusters, &clusterItem)
			}
		}

		if classicalClusterMap, ok := helper.ConvertInterfacesHeadToMap(exclusiveClusterMap["classical_cluster"]); ok {
			clusterItem := clb.ClusterItem{}
			if v, ok := classicalClusterMap["cluster_id"].(string); ok && v != "" {
				clusterItem.ClusterId = helper.String(v)
			}

			if v, ok := classicalClusterMap["cluster_name"].(string); ok && v != "" {
				clusterItem.ClusterName = helper.String(v)
			}

			if v, ok := classicalClusterMap["zone"].(string); ok && v != "" {
				clusterItem.Zone = helper.String(v)
			}

			exclusiveCluster.ClassicalCluster = &clusterItem
		}

		request.ExclusiveCluster = &exclusiveCluster
	}

	var response *clb.CreateLoadBalancerResponse
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().CreateLoadBalancer(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.RequestId == nil {
			return resource.NonRetryableError(fmt.Errorf("Create CLB instance failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create CLB instance failed, reason:%+v", logId, err)
		return err
	}

	// wait
	requestId := *response.Response.RequestId
	clbId, err := waitForTaskFinishGetIDWithTimeout(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	if clbId == "" {
		return fmt.Errorf("[CHECK][CLB instance][Create] check: response error, load balancer id is nil")
	}

	d.SetId(clbId)

	if v, ok := d.GetOk("security_groups"); ok {
		sgRequest := clb.NewSetLoadBalancerSecurityGroupsRequest()
		sgRequest.LoadBalancerId = helper.String(clbId)
		securityGroups := v.([]interface{})
		sgRequest.SecurityGroups = make([]*string, 0, len(securityGroups))
		for i := range securityGroups {
			if securityGroups[i] != nil {
				if securityGroup, ok := securityGroups[i].(string); ok && securityGroup != "" {
					sgRequest.SecurityGroups = append(sgRequest.SecurityGroups, &securityGroup)
				}
			}
		}

		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			sgResponse, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().SetLoadBalancerSecurityGroups(sgRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, sgRequest.GetAction(), sgRequest.ToJsonString(), sgResponse.ToJsonString())
				requestId := *sgResponse.Response.RequestId
				retryErr := waitForTaskFinishWithTimeout(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient(), d.Timeout(schema.TimeoutCreate))
				if retryErr != nil {
					return tccommon.RetryError(errors.WithStack(retryErr))
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s create CLB instance security_groups failed, reason:%+v", logId, err)
			return err
		}
	}

	if v, ok := d.GetOk("log_set_id"); ok {
		if u, ok := d.GetOk("log_topic_id"); ok {
			logRequest := clb.NewSetLoadBalancerClsLogRequest()
			logRequest.LoadBalancerId = helper.String(clbId)
			logRequest.LogSetId = helper.String(v.(string))
			logRequest.LogTopicId = helper.String(u.(string))
			err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				logResponse, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().SetLoadBalancerClsLog(logRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
						logId, logRequest.GetAction(), logRequest.ToJsonString(), logResponse.ToJsonString())
					requestId := *logResponse.Response.RequestId
					retryErr := waitForTaskFinishWithTimeout(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient(), d.Timeout(schema.TimeoutCreate))
					if retryErr != nil {
						return tccommon.RetryError(errors.WithStack(retryErr))
					}
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s set CLB instance log failed, reason:%+v", logId, err)
				return err
			}

		} else {
			return fmt.Errorf("log_topic_id and log_set_id must be set together.")
		}
	}

	if targetRegionInfoRegion != "" {
		isLoadBalancePassToTgt := d.Get("load_balancer_pass_to_target").(bool)
		targetRegionInfo := clb.TargetRegionInfo{
			Region: &targetRegionInfoRegion,
			VpcId:  &targetRegionInfoVpcId,
		}

		mRequest := clb.NewModifyLoadBalancerAttributesRequest()
		mRequest.LoadBalancerId = helper.String(clbId)
		mRequest.TargetRegionInfo = &targetRegionInfo
		mRequest.LoadBalancerPassToTarget = &isLoadBalancePassToTgt
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			mResponse, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyLoadBalancerAttributes(mRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, mRequest.GetAction(), mRequest.ToJsonString(), mResponse.ToJsonString())
				requestId := *mResponse.Response.RequestId
				retryErr := waitForTaskFinishWithTimeout(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient(), d.Timeout(schema.TimeoutCreate))
				if retryErr != nil {
					return tccommon.RetryError(errors.WithStack(retryErr))
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s create CLB instance failed, reason:%+v", logId, err)
			return err
		}
	}

	if v, ok := d.GetOkExists("delete_protect"); ok {
		isDeleteProect := v.(bool)
		if isDeleteProect {
			mRequest := clb.NewModifyLoadBalancerAttributesRequest()
			mRequest.LoadBalancerId = helper.String(clbId)
			mRequest.DeleteProtect = &isDeleteProect
			err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				mResponse, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyLoadBalancerAttributes(mRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
						logId, mRequest.GetAction(), mRequest.ToJsonString(), mResponse.ToJsonString())
					requestId := *mResponse.Response.RequestId
					retryErr := waitForTaskFinishWithTimeout(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient(), d.Timeout(schema.TimeoutCreate))
					if retryErr != nil {
						return tccommon.RetryError(errors.WithStack(retryErr))
					}
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s create CLB instance failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	if v, ok := d.GetOkExists("associate_endpoint"); ok {
		endpointId := v.(string)
		if endpointId != "" {
			mRequest := clb.NewModifyLoadBalancerAttributesRequest()
			mRequest.LoadBalancerId = helper.String(clbId)
			mRequest.AssociateEndpoint = &endpointId
			err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				mResponse, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyLoadBalancerAttributes(mRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, mRequest.GetAction(), mRequest.ToJsonString(), mResponse.ToJsonString())
					if mResponse == nil || mResponse.Response == nil || mResponse.Response.RequestId == nil {
						return resource.NonRetryableError(fmt.Errorf("Modify load balancer attributes failed, Response is nil."))
					}

					requestId := *mResponse.Response.RequestId
					retryErr := waitForTaskFinishWithTimeout(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient(), d.Timeout(schema.TimeoutCreate))
					if retryErr != nil {
						return tccommon.RetryError(errors.WithStack(retryErr))
					}
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s create CLB instance failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudClbInstanceRead(d, meta)
}

func resourceTencentCloudClbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instance   *clb.LoadBalancer
		clbId      = d.Id()
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := clbService.DescribeLoadBalancerById(ctx, clbId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		instance = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read CLB instance failed, reason:%+v", logId, err)
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("network_type", instance.LoadBalancerType)
	_ = d.Set("clb_name", instance.LoadBalancerName)
	_ = d.Set("clb_vips", helper.StringsInterfaces(instance.LoadBalancerVips))
	_ = d.Set("subnet_id", instance.SubnetId)
	_ = d.Set("vpc_id", instance.VpcId)
	_ = d.Set("target_region_info_region", instance.TargetRegionInfo.Region)
	_ = d.Set("target_region_info_vpc_id", instance.TargetRegionInfo.VpcId)
	_ = d.Set("project_id", instance.ProjectId)
	_ = d.Set("security_groups", helper.StringsInterfaces(instance.SecureGroups))
	_ = d.Set("domain", instance.LoadBalancerDomain)
	_ = d.Set("ipv6_mode", instance.IPv6Mode)
	_ = d.Set("address_ipv6", instance.AddressIPv6)

	if instance.ClusterIds != nil && len(instance.ClusterIds) > 0 {
		_ = d.Set("cluster_id", instance.ClusterIds[0])
	}

	if instance.SlaType != nil {
		_ = d.Set("sla_type", instance.SlaType)
	}

	if instance.VipIsp != nil {
		_ = d.Set("vip_isp", instance.VipIsp)
	}

	if instance.LoadBalancerVips != nil && len(instance.LoadBalancerVips) > 0 {
		_ = d.Set("vip", instance.LoadBalancerVips[0])
	}

	if instance.AddressIPVersion != nil {
		if *instance.AddressIPVersion == "ipv6" && instance.IPv6Mode != nil && *instance.IPv6Mode == "IPv6FullChain" {
			_ = d.Set("address_ip_version", instance.IPv6Mode)
		} else {
			_ = d.Set("address_ip_version", instance.AddressIPVersion)
		}
	}

	if instance.NetworkAttributes != nil {
		_ = d.Set("internet_bandwidth_max_out", instance.NetworkAttributes.InternetMaxBandwidthOut)
		_ = d.Set("internet_charge_type", instance.NetworkAttributes.InternetChargeType)
	}

	if instance.MasterZone != nil {
		_ = d.Set("master_zone_id", instance.MasterZone.Zone)
		_ = d.Set("zone_id", instance.MasterZone.Zone)
	}

	if instance.BackupZoneSet != nil && len(instance.BackupZoneSet) > 0 {
		_ = d.Set("slave_zone_id", instance.BackupZoneSet[0].Zone)
	}

	_ = d.Set("load_balancer_pass_to_target", instance.LoadBalancerPassToTarget)
	_ = d.Set("log_set_id", instance.LogSetId)
	_ = d.Set("log_topic_id", instance.LogTopicId)

	if _, ok := d.GetOk("snat_pro"); ok {
		_ = d.Set("snat_pro", instance.SnatPro)
	}

	if *instance.LoadBalancerType == "INTERNAL" {
		request := vpc.NewDescribeAddressesRequest()
		request.Filters = []*vpc.Filter{
			{
				Name:   helper.String("instance-id"),
				Values: helper.Strings([]string{clbId}),
			},
		}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeAddresses(request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if result == nil || result.Response == nil || result.Response.AddressSet == nil {
				e = fmt.Errorf("Describe CLB instance EIP failed")
				return resource.NonRetryableError(e)
			}

			if len(result.Response.AddressSet) == 1 {
				if result.Response.AddressSet[0].AddressId != nil {
					_ = d.Set("eip_address_id", result.Response.AddressSet[0].AddressId)
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s Describe CLB instance EIP failed, reason:%+v", logId, err)
			return err
		}
	}

	if instance.AssociateEndpoint != nil {
		_ = d.Set("associate_endpoint", instance.AssociateEndpoint)
	}

	if instance.ExclusiveCluster != nil {
		exclusiveClusterMap := map[string]interface{}{}
		if instance.ExclusiveCluster.L4Clusters != nil {
			l4ClustersList := make([]map[string]interface{}, 0, len(instance.ExclusiveCluster.L4Clusters))
			for _, l4Clusters := range instance.ExclusiveCluster.L4Clusters {
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

		if instance.ExclusiveCluster.L7Clusters != nil {
			l7ClustersList := make([]map[string]interface{}, 0, len(instance.ExclusiveCluster.L7Clusters))
			for _, l7Clusters := range instance.ExclusiveCluster.L7Clusters {
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

		if instance.ExclusiveCluster.ClassicalCluster != nil {
			classicalClusterMap := map[string]interface{}{}
			if instance.ExclusiveCluster.ClassicalCluster.ClusterId != nil {
				classicalClusterMap["cluster_id"] = instance.ExclusiveCluster.ClassicalCluster.ClusterId
			}

			if instance.ExclusiveCluster.ClassicalCluster.ClusterName != nil {
				classicalClusterMap["cluster_name"] = instance.ExclusiveCluster.ClassicalCluster.ClusterName
			}

			if instance.ExclusiveCluster.ClassicalCluster.Zone != nil {
				classicalClusterMap["zone"] = instance.ExclusiveCluster.ClassicalCluster.Zone
			}

			exclusiveClusterMap["classical_cluster"] = []interface{}{classicalClusterMap}
		}

		_ = d.Set("exclusive_cluster", []interface{}{exclusiveClusterMap})
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "clb", "clb", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)
	return nil
}

func resourceTencentCloudClbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_instance.update")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbId = d.Id()
	)

	immutableArgs := []string{"snat_ips", "dynamic_vip", "master_zone_id", "slave_zone_id", "vpc_id", "subnet_id", "address_ip_version", "bandwidth_package_id", "zone_id"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	d.Partial(true)

	request := clb.NewModifyLoadBalancerAttributesRequest()
	request.LoadBalancerId = helper.String(clbId)
	clbName := ""
	targetRegionInfo := clb.TargetRegionInfo{}
	internet := clb.InternetAccessible{}
	changed := false
	isLoadBalancerPassToTgt := false
	isDeleteProtect := false
	snatPro := d.Get("snat_pro").(bool)

	if d.HasChange("clb_name") {
		changed = true
		clbName = d.Get("clb_name").(string)
		flag, err := checkSameName(clbName, meta)
		if err != nil {
			return err
		}

		if flag {
			return fmt.Errorf("[CHECK][CLB instance][Update] check: Same CLB name %s exists!", clbName)
		}

		request.LoadBalancerName = helper.String(clbName)
	}

	if d.HasChange("target_region_info_region") || d.HasChange("target_region_info_vpc_id") {
		if d.Get("network_type") == CLB_NETWORK_TYPE_INTERNAL {
			return fmt.Errorf("[CHECK][CLB instance %s][Update] check: INTERNAL network_type do not support this operation with target_region_info", clbId)
		}

		changed = true
		region := d.Get("target_region_info_region").(string)
		vpcId := d.Get("target_region_info_vpc_id").(string)
		targetRegionInfo = clb.TargetRegionInfo{
			Region: &region,
			VpcId:  &vpcId,
		}
		request.TargetRegionInfo = &targetRegionInfo
	}

	if d.HasChange("sla_type") {
		slaRequest := clb.NewModifyLoadBalancerSlaRequest()
		param := clb.SlaUpdateParam{}
		param.LoadBalancerId = &clbId
		param.SlaType = helper.String(d.Get("sla_type").(string))
		slaRequest.LoadBalancerSla = []*clb.SlaUpdateParam{&param}
		var taskId string
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyLoadBalancerSla(slaRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			taskId = *result.Response.RequestId
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update clb instanceSlaConfig failed, reason:%+v", logId, err)
			return err
		}

		retryErr := waitForTaskFinishWithTimeout(taskId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient(), d.Timeout(schema.TimeoutUpdate))
		if retryErr != nil {
			return retryErr
		}
	}

	if d.HasChange("internet_charge_type") || d.HasChange("internet_bandwidth_max_out") {
		changed = true
		chargeType := d.Get("internet_charge_type").(string)
		bandwidth := d.Get("internet_bandwidth_max_out").(int)
		if chargeType != "" {
			internet.InternetChargeType = &chargeType
		}

		if bandwidth > 0 {
			internet.InternetMaxBandwidthOut = helper.IntInt64(bandwidth)
		}

		request.InternetChargeInfo = &internet
	}

	if d.HasChange("load_balancer_pass_to_target") {
		changed = true
		isLoadBalancerPassToTgt = d.Get("load_balancer_pass_to_target").(bool)
		request.LoadBalancerPassToTarget = &isLoadBalancerPassToTgt
	}

	if d.HasChange("snat_pro") {
		changed = true
		request.SnatPro = &snatPro
	}

	if d.HasChange("delete_protect") {
		changed = true
		isDeleteProtect = d.Get("delete_protect").(bool)
		request.DeleteProtect = &isDeleteProtect
	}

	if d.HasChange("associate_endpoint") {
		changed = true
		associateEndpoint := d.Get("associate_endpoint").(string)
		request.AssociateEndpoint = &associateEndpoint
	}

	if changed {
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyLoadBalancerAttributes(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
				requestId := *response.Response.RequestId
				retryErr := waitForTaskFinishWithTimeout(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient(), d.Timeout(schema.TimeoutUpdate))
				if retryErr != nil {
					return tccommon.RetryError(retryErr)
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update CLB instance failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("security_groups") {
		sgRequest := clb.NewSetLoadBalancerSecurityGroupsRequest()
		sgRequest.LoadBalancerId = helper.String(clbId)
		securityGroups := d.Get("security_groups").([]interface{})
		sgRequest.SecurityGroups = make([]*string, 0, len(securityGroups))
		for i := range securityGroups {
			if securityGroup, ok := securityGroups[i].(string); ok && securityGroup != "" {
				sgRequest.SecurityGroups = append(sgRequest.SecurityGroups, &securityGroup)
			}
		}

		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			sgResponse, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().SetLoadBalancerSecurityGroups(sgRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, sgRequest.GetAction(), sgRequest.ToJsonString(), sgResponse.ToJsonString())
				requestId := *sgResponse.Response.RequestId
				retryErr := waitForTaskFinishWithTimeout(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient(), d.Timeout(schema.TimeoutUpdate))
				if retryErr != nil {
					return tccommon.RetryError(errors.WithStack(retryErr))
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update CLB instance security_group failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("log_set_id") || d.HasChange("log_topic_id") {
		logSetId := d.Get("log_set_id")
		logTopicId := d.Get("log_topic_id")
		logRequest := clb.NewSetLoadBalancerClsLogRequest()
		logRequest.LoadBalancerId = helper.String(clbId)
		logRequest.LogSetId = helper.String(logSetId.(string))
		logRequest.LogTopicId = helper.String(logTopicId.(string))
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			logResponse, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().SetLoadBalancerClsLog(logRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, logRequest.GetAction(), logRequest.ToJsonString(), logResponse.ToJsonString())
				requestId := *logResponse.Response.RequestId
				retryErr := waitForTaskFinishWithTimeout(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient(), d.Timeout(schema.TimeoutUpdate))
				if retryErr != nil {
					return tccommon.RetryError(errors.WithStack(retryErr))
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s set CLB instance log failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("project_id") {
		var projectId int
		if v, ok := d.GetOkExists("project_id"); ok {
			projectId = v.(int)
		}

		pRequest := clb.NewModifyLoadBalancersProjectRequest()
		pRequest.LoadBalancerIds = []*string{&clbId}
		pRequest.ProjectId = helper.IntUint64(projectId)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			pResponse, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyLoadBalancersProject(pRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, pRequest.GetAction(), pRequest.ToJsonString(), pResponse.ToJsonString())
				requestId := *pResponse.Response.RequestId
				retryErr := waitForTaskFinishWithTimeout(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient(), d.Timeout(schema.TimeoutUpdate))
				if retryErr != nil {
					return tccommon.RetryError(errors.WithStack(retryErr))
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update CLB instance project_id failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("eip_address_id") {
		oldEip, newEip := d.GetChange("eip_address_id")
		oldEipStr := oldEip.(string)
		newEipStr := newEip.(string)
		// delete old first
		if oldEipStr != "" {
			request := vpc.NewDisassociateAddressRequest()
			request.AddressId = helper.String(oldEipStr)
			err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DisassociateAddress(request)
				if e != nil {
					return tccommon.RetryError(e)
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s Disassociate EIP failed, reason:%+v", logId, err)
				return err
			}

			// wait
			eipRequest := vpc.NewDescribeAddressesRequest()
			eipRequest.AddressIds = helper.Strings([]string{oldEipStr})
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeAddresses(eipRequest)
				if e != nil {
					return tccommon.RetryError(e)
				}

				if result == nil || result.Response == nil || result.Response.AddressSet == nil || len(result.Response.AddressSet) != 1 {
					e = fmt.Errorf("Describe CLB instance EIP failed")
					return resource.NonRetryableError(e)
				}

				if *result.Response.AddressSet[0].AddressStatus != "UNBIND" {
					return resource.RetryableError(fmt.Errorf("EIP status is still %s", *result.Response.AddressSet[0].AddressStatus))
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s Describe CLB instance EIP failed, reason:%+v", logId, err)
				return err
			}
		}

		// attach new
		if newEipStr != "" {
			request := vpc.NewAssociateAddressRequest()
			request.AddressId = helper.String(newEipStr)
			request.InstanceId = helper.String(clbId)
			err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AssociateAddress(request)
				if e != nil {
					return tccommon.RetryError(e)
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s Associate EIP failed, reason:%+v", logId, err)
				return err
			}

			// wait
			eipRequest := vpc.NewDescribeAddressesRequest()
			eipRequest.AddressIds = helper.Strings([]string{newEipStr})
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeAddresses(eipRequest)
				if e != nil {
					return tccommon.RetryError(e)
				}

				if result == nil || result.Response == nil || result.Response.AddressSet == nil || len(result.Response.AddressSet) != 1 {
					e = fmt.Errorf("Describe CLB instance EIP failed")
					return resource.NonRetryableError(e)
				}

				if *result.Response.AddressSet[0].AddressStatus != "BIND" {
					return resource.RetryableError(fmt.Errorf("EIP status is still %s", *result.Response.AddressSet[0].AddressStatus))
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s Describe CLB instance EIP failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("clb", "clb", tcClient.Region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}

	d.Partial(false)
	return nil
}

func resourceTencentCloudClbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_instance.delete")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		clbId      = d.Id()
	)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := clbService.DeleteLoadBalancerById(ctx, clbId)
		if e != nil {
			if ve, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if ve.GetCode() == "FailedOperation.ResourceInOperating" {
					return tccommon.RetryError(e, "FailedOperation.ResourceInOperating")
				}
			}

			return tccommon.RetryError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete CLB instance failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

func checkSameName(name string, meta interface{}) (flag bool, errRet error) {
	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	flag = false
	params := make(map[string]interface{})
	params["clb_name"] = name
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		clbs, e := clbService.DescribeLoadBalancerByFilter(ctx, params)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if len(clbs) > 0 {
			//this function is a fuzzy query
			// so take a further check
			for _, clbInfo := range clbs {
				if *clbInfo.LoadBalancerName == name {
					flag = true
					return nil
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read CLB instance failed, reason:%+v", logId, err)
	}

	errRet = err
	return
}
