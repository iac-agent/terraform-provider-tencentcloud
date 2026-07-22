package clb

import (
	"context"
	"encoding/json"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbInstancesRead,

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "需要查询的CLB ID。",
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CLB_NETWORK_TYPE),
				Description: "CLB实例的类型，可用值包括“OPEN”和“INTERNAL”。",
			},
			"clb_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "需要查询的CLB名称。",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "CLB 的项目 ID。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"master_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "主可用区域 ID。",
			},
			"clb_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "云负载均衡器列表。每个元素包含以下属性：",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"clb_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB 的 ID。",
						},
						"clb_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB 名称。",
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB 的类型。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "项目 ID。",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "集群ID。",
						},
						"clb_vips": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "CLB的虚拟服务地址表。",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CLB 的状态。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB的创建时间。",
						},
						"status_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB的最新状态转换时间。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "专有网络ID。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "子网的 ID。",
						},
						"security_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "安全组的ID集。",
						},
						"target_region_info_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "后端服务的区域信息附在CLB中。",
						},
						"target_region_info_vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB中附有后端服务的VpcId信息。",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "此 CLB 中的可用标签。",
						},
						"address_ip_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP版本，仅适用于开放的CLB。有效值为“IPV4”、“IPV6”和“IPv6FullChain”。",
						},
						"vip_isp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "网络运营商，仅适用于开通CLB。有效值为“CMCC”（中国移动）、“CTCC”（电信）、“CUCC”（中国联通）和“BGP”。如果指定了该ISP，则网络计费方式只能使用带宽套餐计费（BANDWIDTH_PACKAGE）。",
						},
						"internet_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "网费类型，仅适用于开通CLB。有效值为“TRAFFIC_POSTPAID_BY_HOUR”、“BANDWIDTH_POSTPAID_BY_HOUR”和“BANDWIDTH_PACKAGE”。",
						},
						"internet_bandwidth_max_out": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最大带宽输出，仅适用于开放的CLB。有效值范围为 [1，2048]。单位是MB。",
						},
						"zone_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "可用区域唯一id（数字表示），该字段可能为空，表示无法获取有效值。",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用区域唯一id（字符串表示），该字段可能为空，表示无法获取有效值。",
						},
						"zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用区域名称，该字段可能为空，表示无法获取有效值。",
						},
						"zone_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "该可用区所属的区域，该字段可能为空，表示无法获取有效值。",
						},
						"local_zone": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "该可用区域是否为本地区域，该字段可能为空，表示无法获取有效值。",
						},
						"numerical_vpc_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数字形式的 私有网络 ID。注意：该字段可能返回null，表示取不到有效值。",
						},
						"zones": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "就近访问模式的VPC内部负载均衡器部署规则的区域。注意：该字段可能返回null，表示取不到有效值。",
						},
						// Basic info fields
						"forward": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CLB类型标识符，1：CLB，0：经典CLB。",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB域名（仅适用于公网经典CLB），逐渐弃用。",
						},
						"load_balancer_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB实例的域。",
						},
						// Network config fields
						"address_ipv6": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB实例的IPv6地址。",
						},
						"ipv6_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "当 IP 版本为 ipv6、IPv6Nat64 或 IPv6FullChain 时的 IPv6 模式。",
						},
						"mix_ip_target": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "IPv6FullChain CLB 第 7 层侦听器支持 IPv4/IPv6 目标的混合绑定。",
						},
						"anycast_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "任播 CLB 发布区域，对于非任播 CLB 返回空字符串。",
						},
						"egress": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "网络出口。",
						},
						"local_bgp": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "IP类型是否为本地BGP。",
						},
						// Billing and lifecycle fields
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "计费类型，PREPAID：预付费，POSTPAID_BY_HOUR：按量付费。",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB实例的过期时间，仅适用于预付费CLB，格式：YYYY-MM-DD HH:mm:ss。",
						},
						"prepaid_period": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "预付费购买期限，单位：月。",
						},
						"prepaid_renew_flag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "预付费续订标志，NOTIFY_AND_AUTO_RENEW：通知并自动续订，NOTIFY_AND_MANUAL_RENEW：通知但不自动续订，DISABLE_NOTIFY_AND_MANUAL_RENEW：不通知且不自动续订。",
						},
						// Log config fields
						"log_set_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "日志服务 (CLS) 日志集 ID。",
						},
						"log_topic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "日志服务（CLS）日志主题ID。",
						},
						"health_log_set_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "日志服务（CLS）健康检查日志集ID。",
						},
						"health_log_topic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "日志服务（CLS）健康检查日志主题ID。",
						},
						// Security and isolation fields
						"open_bgp": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "DDoS高防LB标识，1：DDoS高防，0：非DDoS高防。",
						},
						"snat": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否启用SNAT。",
						},
						"snat_pro": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否启用 SnatPro。",
						},
						"snat_ips": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "启用S​​natPro后的SnatIp列表（JSON格式）。",
						},
						"isolation": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否隔离，0：不隔离，1：隔离。",
						},
						"isolated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB实例被隔离的时间，格式：YYYY-MM-DD HH:mm:ss。",
						},
						"is_block": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "VIP是否被屏蔽。",
						},
						"is_block_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "阻塞或解阻塞时间，格式：YYYY-MM-DD HH:mm:ss。",
						},
						"is_ddos": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否可以绑定DDoS高防。",
						},
						// Performance and capacity fields
						"sla_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "性能容量类型规范（clb.c1.small/clb.c2.medium/clb.c3.small/clb.c3.medium/clb.c4.small/clb.c4.medium/clb.c4.large/clb.c4.xlarge 或空字符串）。",
						},
						"exclusive": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例类型是否独占，1：独占，0：不独占。",
						},
						"target_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "绑定的后端服务数量。",
						},
						// Cluster and deployment fields
						"cluster_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "集群ID列表。",
						},
						"cluster_tag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "七层专属标签。",
						},
						"nfv_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB是否是NFV，空：否，l7nfv：七层是NFV。",
						},
						"backup_zone_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "备份区域列表，每个元素包含zone_id/可用区/zone_name/zone_region/local_zone。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "可用区域唯一ID（数字表示）。",
									},
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "可用区域唯一ID（字符串表示）。",
									},
									"zone_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "可用区域名称。",
									},
									"zone_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "该可用区所属的区域。",
									},
									"local_zone": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "该可用区域是否为本地区域。",
									},
								},
							},
						},
						"available_zone_affinity_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用区域转发关联信息（JSON 格式）。",
						},
						// Advanced config fields
						"config_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB维度个性化配置ID。",
						},
						"load_balancer_pass_to_target": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "后端服务是否允许来自 CLB 的流量。",
						},
						"attribute_flags": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "CLB 属性标志数组。",
						},
						"exclusive_cluster": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "内部独占集群信息（JSON格式）。",
						},
						"extra_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "保留字段，一般不需要关注（JSON格式）。",
						},
						"associate_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "与 CLB 实例关联的端点 ID。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudClbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clb_instances.read")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		clbService = ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		clbs       []*clb.LoadBalancer
	)

	params := make(map[string]interface{})
	if v, ok := d.GetOk("clb_id"); ok {
		params["clb_id"] = v.(string)
	}

	if v, ok := d.GetOk("clb_name"); ok {
		params["clb_name"] = v.(string)
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		params["project_id"] = v.(int)
	}

	if v, ok := d.GetOk("network_type"); ok {
		params["network_type"] = v.(string)
	}

	if v, ok := d.GetOk("master_zone"); ok {
		params["master_zone"] = v.(string)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := clbService.DescribeLoadBalancerByFilter(ctx, params)
		if e != nil {
			return tccommon.RetryError(e)
		}

		clbs = results
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read CLB instances failed, reason:%+v", logId, err)
		return err
	}

	clbList := make([]map[string]interface{}, 0, len(clbs))
	ids := make([]string, 0, len(clbs))
	for _, clbInstance := range clbs {
		mapping := map[string]interface{}{
			"clb_id":                    clbInstance.LoadBalancerId,
			"clb_name":                  clbInstance.LoadBalancerName,
			"network_type":              clbInstance.LoadBalancerType,
			"status":                    clbInstance.Status,
			"create_time":               clbInstance.CreateTime,
			"status_time":               clbInstance.StatusTime,
			"project_id":                clbInstance.ProjectId,
			"vpc_id":                    clbInstance.VpcId,
			"subnet_id":                 clbInstance.SubnetId,
			"clb_vips":                  helper.StringsInterfaces(clbInstance.LoadBalancerVips),
			"target_region_info_region": clbInstance.TargetRegionInfo.Region,
			"target_region_info_vpc_id": clbInstance.TargetRegionInfo.VpcId,
			"address_ip_version":        clbInstance.AddressIPVersion,
			"vip_isp":                   clbInstance.VipIsp,
			"security_groups":           helper.StringsInterfaces(clbInstance.SecureGroups),
		}

		if clbInstance.ClusterIds != nil && len(clbInstance.ClusterIds) > 0 {
			mapping["cluster_id"] = *clbInstance.ClusterIds[0]
		}

		if clbInstance.NetworkAttributes != nil {
			mapping["internet_charge_type"] = *clbInstance.NetworkAttributes.InternetChargeType
			mapping["internet_bandwidth_max_out"] = *clbInstance.NetworkAttributes.InternetMaxBandwidthOut
		}

		if clbInstance.MasterZone != nil {
			mapping["zone_id"] = *clbInstance.MasterZone.ZoneId
			mapping["zone"] = *clbInstance.MasterZone.Zone
			mapping["zone_name"] = *clbInstance.MasterZone.ZoneName
			mapping["zone_region"] = *clbInstance.MasterZone.ZoneRegion
			mapping["local_zone"] = *clbInstance.MasterZone.LocalZone
		}

		if clbInstance.Tags != nil {
			tags := make(map[string]interface{}, len(clbInstance.Tags))
			for _, t := range clbInstance.Tags {
				tags[*t.TagKey] = *t.TagValue
			}

			mapping["tags"] = tags
		}

		if clbInstance.NumericalVpcId != nil {
			mapping["numerical_vpc_id"] = clbInstance.NumericalVpcId
		}

		if clbInstance.Zones != nil {
			mapping["zones"] = helper.StringsInterfaces(clbInstance.Zones)
		}

		// Basic info fields
		if clbInstance.Forward != nil {
			mapping["forward"] = clbInstance.Forward
		}
		if clbInstance.Domain != nil {
			mapping["domain"] = clbInstance.Domain
		}
		if clbInstance.LoadBalancerDomain != nil {
			mapping["load_balancer_domain"] = clbInstance.LoadBalancerDomain
		}

		// Network config fields
		if clbInstance.AddressIPv6 != nil {
			mapping["address_ipv6"] = clbInstance.AddressIPv6
		}
		if clbInstance.IPv6Mode != nil {
			mapping["ipv6_mode"] = clbInstance.IPv6Mode
		}
		if clbInstance.MixIpTarget != nil {
			mapping["mix_ip_target"] = clbInstance.MixIpTarget
		}
		if clbInstance.AnycastZone != nil {
			mapping["anycast_zone"] = clbInstance.AnycastZone
		}
		if clbInstance.Egress != nil {
			mapping["egress"] = clbInstance.Egress
		}
		if clbInstance.LocalBgp != nil {
			mapping["local_bgp"] = clbInstance.LocalBgp
		}

		// Billing and lifecycle fields
		if clbInstance.ChargeType != nil {
			mapping["charge_type"] = clbInstance.ChargeType
		}
		if clbInstance.ExpireTime != nil {
			mapping["expire_time"] = clbInstance.ExpireTime
		}
		if clbInstance.PrepaidAttributes != nil {
			if clbInstance.PrepaidAttributes.Period != nil {
				mapping["prepaid_period"] = clbInstance.PrepaidAttributes.Period
			}
			if clbInstance.PrepaidAttributes.RenewFlag != nil {
				mapping["prepaid_renew_flag"] = clbInstance.PrepaidAttributes.RenewFlag
			}
		}

		// Log config fields
		if clbInstance.LogSetId != nil {
			mapping["log_set_id"] = clbInstance.LogSetId
		}
		if clbInstance.LogTopicId != nil {
			mapping["log_topic_id"] = clbInstance.LogTopicId
		}
		if clbInstance.HealthLogSetId != nil {
			mapping["health_log_set_id"] = clbInstance.HealthLogSetId
		}
		if clbInstance.HealthLogTopicId != nil {
			mapping["health_log_topic_id"] = clbInstance.HealthLogTopicId
		}

		// Security and isolation fields
		if clbInstance.OpenBgp != nil {
			mapping["open_bgp"] = clbInstance.OpenBgp
		}
		if clbInstance.Snat != nil {
			mapping["snat"] = clbInstance.Snat
		}
		if clbInstance.SnatPro != nil {
			mapping["snat_pro"] = clbInstance.SnatPro
		}
		if clbInstance.SnatIps != nil {
			snatIpsJSON, _ := json.Marshal(clbInstance.SnatIps)
			mapping["snat_ips"] = string(snatIpsJSON)
		}
		if clbInstance.Isolation != nil {
			mapping["isolation"] = clbInstance.Isolation
		}
		if clbInstance.IsolatedTime != nil {
			mapping["isolated_time"] = clbInstance.IsolatedTime
		}
		if clbInstance.IsBlock != nil {
			mapping["is_block"] = clbInstance.IsBlock
		}
		if clbInstance.IsBlockTime != nil {
			mapping["is_block_time"] = clbInstance.IsBlockTime
		}
		if clbInstance.IsDDos != nil {
			mapping["is_ddos"] = clbInstance.IsDDos
		}

		// Performance and capacity fields
		if clbInstance.SlaType != nil {
			mapping["sla_type"] = clbInstance.SlaType
		}
		if clbInstance.Exclusive != nil {
			mapping["exclusive"] = clbInstance.Exclusive
		}
		if clbInstance.TargetCount != nil {
			mapping["target_count"] = clbInstance.TargetCount
		}

		// Cluster and deployment fields
		if clbInstance.ClusterIds != nil {
			mapping["cluster_ids"] = helper.StringsInterfaces(clbInstance.ClusterIds)
		}
		if clbInstance.ClusterTag != nil {
			mapping["cluster_tag"] = clbInstance.ClusterTag
		}
		if clbInstance.NfvInfo != nil {
			mapping["nfv_info"] = clbInstance.NfvInfo
		}
		if clbInstance.BackupZoneSet != nil {
			backupZones := make([]map[string]interface{}, 0, len(clbInstance.BackupZoneSet))
			for _, zone := range clbInstance.BackupZoneSet {
				backupZone := make(map[string]interface{})
				if zone.ZoneId != nil {
					backupZone["zone_id"] = *zone.ZoneId
				}
				if zone.Zone != nil {
					backupZone["zone"] = *zone.Zone
				}
				if zone.ZoneName != nil {
					backupZone["zone_name"] = *zone.ZoneName
				}
				if zone.ZoneRegion != nil {
					backupZone["zone_region"] = *zone.ZoneRegion
				}
				if zone.LocalZone != nil {
					backupZone["local_zone"] = *zone.LocalZone
				}
				backupZones = append(backupZones, backupZone)
			}
			mapping["backup_zone_set"] = backupZones
		}
		if clbInstance.AvailableZoneAffinityInfo != nil {
			availableZoneAffinityJSON, _ := json.Marshal(clbInstance.AvailableZoneAffinityInfo)
			mapping["available_zone_affinity_info"] = string(availableZoneAffinityJSON)
		}

		// Advanced config fields
		if clbInstance.ConfigId != nil {
			mapping["config_id"] = clbInstance.ConfigId
		}
		if clbInstance.LoadBalancerPassToTarget != nil {
			mapping["load_balancer_pass_to_target"] = clbInstance.LoadBalancerPassToTarget
		}
		if clbInstance.AttributeFlags != nil {
			mapping["attribute_flags"] = helper.StringsInterfaces(clbInstance.AttributeFlags)
		}
		if clbInstance.ExclusiveCluster != nil {
			exclusiveClusterJSON, _ := json.Marshal(clbInstance.ExclusiveCluster)
			mapping["exclusive_cluster"] = string(exclusiveClusterJSON)
		}
		if clbInstance.ExtraInfo != nil {
			extraInfoJSON, _ := json.Marshal(clbInstance.ExtraInfo)
			mapping["extra_info"] = string(extraInfoJSON)
		}
		if clbInstance.AssociateEndpoint != nil {
			mapping["associate_endpoint"] = clbInstance.AssociateEndpoint
		}

		clbList = append(clbList, mapping)
		ids = append(ids, *clbInstance.LoadBalancerId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("clb_list", clbList); e != nil {
		log.Printf("[CRITAL]%s provider set CLB list fail, reason:%+v", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), clbList); e != nil {
			return e
		}
	}

	return nil
}
