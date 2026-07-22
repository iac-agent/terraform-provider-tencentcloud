package clb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClbInstanceDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbInstanceDetailRead,
		Schema: map[string]*schema.Schema{
			"fields": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "字段列表。仅返回指定的字段。如果留空，则返回“null”。默认添加字段“LoadBalancerId”和“LoadBalancerName”。有关字段的详细信息。",
			},

			"target_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "目标类型。有效值：NODE 和 GROUP。如果字段列表中包含“TargetId”、“TargetAddress”、“TargetPort”、“TargetWeight”等字段，则必须导出目标组或非目标组的“Target”。",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "描述 CLB 实例详细信息的查询列表过滤条件：loadbalancer-ID - String - 必填：否 - （过滤条件）CLB 实例 ID，如 lb-12345678； 项目-ID - String - 必填：否 - （过滤条件）项目ID，如0、123； 网络 - String - 必填：否 - （过滤条件）CLB实例的网络类型，例如Public、Private等。&amp;lt;/li&gt;&amp;lt;li&gt; VIP - String - 必填：否 - （过滤条件）CLB实例VIP，如1.1.1.1、2204::22:3； 目标-ip - String - 必填：否 - （过滤条件）目标真实服务器的私有IP，例如1.1.1.1和2203::214:4； vpcid - String - 必填：否 - （过滤条件）CLB实例所属VPC实例标识符，如vpc-12345678； 可用区 - String - 必填：否 - （过滤条件）CLB实例所在可用区，如ap-guangzhou-1； 标签-键 - String - 必填：否 - （过滤条件）CLB实例的Tag 键，例如名称； 标签:* - 字符串 - 必填：否 - （过滤条件）CLB 实例标记，冒号后跟标记键。例如，使用 {名称: 标签:名称,Values: [zhangsan，lisi]} 过滤标签键“名称”，标签值“zhangsan”和“lisi”； fuzzy-search - 字符串 - 必填：否 - （过滤条件）对 CLB 实例 VIP 和 CLB 实例名称进行模糊搜索，如 1。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤器名称。",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "过滤值数组。",
						},
					},
				},
			},

			"load_balancer_detail_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "CLB实例详细信息列表。注意：该字段可能返回null，表示取不到有效值。",
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
							Description: "CLB实例网络类型：Public：公网； 私有：私有网络。注意：该字段可能返回null，表示取不到有效值。",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CLB实例状态，包括：0：正在创建； 1：正在运行。注意：该字段可能返回null，表示取不到有效值。",
						},
						"address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB实例VIP。注意：该字段可能返回null，表示取不到有效值。",
						},
						"address_ipv6": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB实例的IPv6 VIP地址。注意：该字段可能返回null，表示取不到有效值。",
						},
						"address_ip_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB实例的IP版本。有效值：IPv4、IPv6。注意：该字段可能返回null，表示取不到有效值。",
						},
						"ipv6_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB实例的IPv6地址类型。有效值：IPv6Nat64、IPv6FullChain。注意：该字段可能返回null，表示取不到有效值。",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB实例所在的可用区。注意：该字段可能返回null，表示取不到有效值。",
						},
						"address_isp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB IP地址所属的ISP。注意：该字段可能返回null，表示取不到有效值。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB实例所属VPC实例ID。注意：该字段可能返回null，表示取不到有效值。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CLB实例所属项目ID。 0：默认项目。注意：该字段可能返回null，表示取不到有效值。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB实例创建时间。注意：该字段可能返回null，表示取不到有效值。",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB实例计费方式。注意：该字段可能返回null，表示取不到有效值。",
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
										Description: "TRAFFIC_POSTPAID_BY_HOUR：按小时按流量按量付费； BANDWIDTH_POSTPAID_BY_HOUR：按小时按带宽按量付费；BANDWIDTH_PACKAGE：按带宽套餐计费（目前只有指定ISP才支持此方式）。",
									},
									"internet_max_bandwidth_out": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "最大出方向带宽（Mbps），仅适用于公网CLB。值范围：0-65,535。默认值：10。",
									},
									"bandwidth_pkg_sub_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "带宽套餐类型，如SINGLEISP 注：该字段可能返回null，表示取不到有效值。",
									},
								},
							},
						},
						"prepaid_attributes": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "CLB实例按量付费属性。注意：该字段可能返回null，表示取不到有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"renew_flag": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "续订类型。 AUTO_RENEW：自动续订； MANUAL_RENEW：手动续订注意：该字段可能返回null，表示取不到有效值。",
									},
									"period": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "周期，表示月数（保留字段）注意：该字段可能返回null，表示取不到有效值。",
									},
								},
							},
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
										Description: "是否开启VIP直连注：该字段可能返回null，表示取不到有效值。",
									},
									"tgw_group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "TgwGroup 名称 注：该字段可能返回null，表示取不到有效值。",
									},
								},
							},
						},
						"config_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB实例的自定义配置ID。多个ID之间必须用逗号（,）分隔。注意：该字段可能返回null，表示取不到有效值。",
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
						"listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB监听ID。注意：该字段可能返回null，表示取不到有效值。",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "监听协议。注意：该字段可能返回null，表示取不到有效值。",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "监听端口。注意：该字段可能返回null，表示取不到有效值。",
						},
						"location_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "转发规则ID。注意：该字段可能返回null，表示取不到有效值。",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "转发规则的域名。注意：该字段可能返回null，表示取不到有效值。",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "转发规则路径。注意：该字段可能返回null，表示取不到有效值。",
						},
						"target_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "目标真实服务器ID。注意：该字段可能返回null，表示取不到有效值。",
						},
						"target_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "目标真实服务器地址。注意：该字段可能返回null，表示取不到有效值。",
						},
						"target_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "目标真实服务器的监听端口。注意：该字段可能返回null，表示取不到有效值。",
						},
						"target_weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "目标真实服务器的转发权重。注意：该字段可能返回null，表示取不到有效值。",
						},
						"isolation": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "0：不隔离； 1：isolated。注意：该字段可能返回null，表示取不到有效值。",
						},
						"security_group": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "CLB实例绑定的安全组列表。注意：该字段可能返回null，表示取不到有效值。",
						},
						"load_balancer_pass_to_target": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CLB实例是否按IP计费。注意：该字段可能返回null，表示取不到有效值。",
						},
						"target_health": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "目标真实服务器的健康状态。注意：该字段可能返回null，表示取不到有效值。",
						},
						"domains": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "转发规则关联的域名列表 注：该字段可能返回null，表示取不到有效值。",
						},
						"slave_zone": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "多AZ CLB实例次可用区注意：该字段可能返回null，表示取不到有效值。",
						},
						"zones": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "私有CLB实例的AZ。仅测试版用户可用。注意：该字段可能返回“null”，表示无法获取有效值。",
						},
						"sni_switch": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否启用SNI。该参数仅对HTTPS监听有意义。注意：该字段可能返回null，表示取不到有效值。",
						},
						"load_balancer_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB实例的域名。注意：该字段可能返回null，表示取不到有效值。",
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

func dataSourceTencentCloudClbInstanceDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clb_instance_detail.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("fields"); ok {
		fieldsSet := v.(*schema.Set).List()
		paramMap["Fields"] = helper.InterfacesStringsPoint(fieldsSet)
	}

	if v, ok := d.GetOk("target_type"); ok {
		paramMap["TargetType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*clb.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := clb.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["Filters"] = tmpSet
	}

	service := ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var loadBalancerDetailSet []*clb.LoadBalancerDetail

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClbInstanceDetailByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		loadBalancerDetailSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(loadBalancerDetailSet))
	tmpList := make([]map[string]interface{}, 0, len(loadBalancerDetailSet))

	if loadBalancerDetailSet != nil {
		for _, loadBalancerDetail := range loadBalancerDetailSet {
			loadBalancerDetailMap := map[string]interface{}{}

			if loadBalancerDetail.LoadBalancerId != nil {
				loadBalancerDetailMap["load_balancer_id"] = loadBalancerDetail.LoadBalancerId
			}

			if loadBalancerDetail.LoadBalancerName != nil {
				loadBalancerDetailMap["load_balancer_name"] = loadBalancerDetail.LoadBalancerName
			}

			if loadBalancerDetail.LoadBalancerType != nil {
				loadBalancerDetailMap["load_balancer_type"] = loadBalancerDetail.LoadBalancerType
			}

			if loadBalancerDetail.Status != nil {
				loadBalancerDetailMap["status"] = loadBalancerDetail.Status
			}

			if loadBalancerDetail.Address != nil {
				loadBalancerDetailMap["address"] = loadBalancerDetail.Address
			}

			if loadBalancerDetail.AddressIPv6 != nil {
				loadBalancerDetailMap["address_ipv6"] = loadBalancerDetail.AddressIPv6
			}

			if loadBalancerDetail.AddressIPVersion != nil {
				loadBalancerDetailMap["address_ip_version"] = loadBalancerDetail.AddressIPVersion
			}

			if loadBalancerDetail.IPv6Mode != nil {
				loadBalancerDetailMap["ipv6_mode"] = loadBalancerDetail.IPv6Mode
			}

			if loadBalancerDetail.Zone != nil {
				loadBalancerDetailMap["zone"] = loadBalancerDetail.Zone
			}

			if loadBalancerDetail.AddressIsp != nil {
				loadBalancerDetailMap["address_isp"] = loadBalancerDetail.AddressIsp
			}

			if loadBalancerDetail.VpcId != nil {
				loadBalancerDetailMap["vpc_id"] = loadBalancerDetail.VpcId
			}

			if loadBalancerDetail.ProjectId != nil {
				loadBalancerDetailMap["project_id"] = loadBalancerDetail.ProjectId
			}

			if loadBalancerDetail.CreateTime != nil {
				loadBalancerDetailMap["create_time"] = loadBalancerDetail.CreateTime
			}

			if loadBalancerDetail.ChargeType != nil {
				loadBalancerDetailMap["charge_type"] = loadBalancerDetail.ChargeType
			}

			if loadBalancerDetail.NetworkAttributes != nil {
				networkAttributesMap := map[string]interface{}{}

				if loadBalancerDetail.NetworkAttributes.InternetChargeType != nil {
					networkAttributesMap["internet_charge_type"] = loadBalancerDetail.NetworkAttributes.InternetChargeType
				}

				if loadBalancerDetail.NetworkAttributes.InternetMaxBandwidthOut != nil {
					networkAttributesMap["internet_max_bandwidth_out"] = loadBalancerDetail.NetworkAttributes.InternetMaxBandwidthOut
				}

				if loadBalancerDetail.NetworkAttributes.BandwidthpkgSubType != nil {
					networkAttributesMap["bandwidth_pkg_sub_type"] = loadBalancerDetail.NetworkAttributes.BandwidthpkgSubType
				}

				loadBalancerDetailMap["network_attributes"] = []interface{}{networkAttributesMap}
			}

			if loadBalancerDetail.PrepaidAttributes != nil {
				prepaidAttributesMap := map[string]interface{}{}

				if loadBalancerDetail.PrepaidAttributes.RenewFlag != nil {
					prepaidAttributesMap["renew_flag"] = loadBalancerDetail.PrepaidAttributes.RenewFlag
				}

				if loadBalancerDetail.PrepaidAttributes.Period != nil {
					prepaidAttributesMap["period"] = loadBalancerDetail.PrepaidAttributes.Period
				}

				loadBalancerDetailMap["prepaid_attributes"] = []interface{}{prepaidAttributesMap}
			}

			if loadBalancerDetail.ExtraInfo != nil {
				extraInfoMap := map[string]interface{}{}

				if loadBalancerDetail.ExtraInfo.ZhiTong != nil {
					extraInfoMap["zhi_tong"] = loadBalancerDetail.ExtraInfo.ZhiTong
				}

				if loadBalancerDetail.ExtraInfo.TgwGroupName != nil {
					extraInfoMap["tgw_group_name"] = loadBalancerDetail.ExtraInfo.TgwGroupName
				}

				loadBalancerDetailMap["extra_info"] = []interface{}{extraInfoMap}
			}

			if loadBalancerDetail.ConfigId != nil {
				loadBalancerDetailMap["config_id"] = loadBalancerDetail.ConfigId
			}

			if loadBalancerDetail.Tags != nil {
				tagsList := []interface{}{}
				for _, tags := range loadBalancerDetail.Tags {
					tagsMap := map[string]interface{}{}

					if tags.TagKey != nil {
						tagsMap["tag_key"] = tags.TagKey
					}

					if tags.TagValue != nil {
						tagsMap["tag_value"] = tags.TagValue
					}

					tagsList = append(tagsList, tagsMap)
				}

				loadBalancerDetailMap["tags"] = tagsList
			}

			if loadBalancerDetail.ListenerId != nil {
				loadBalancerDetailMap["listener_id"] = loadBalancerDetail.ListenerId
			}

			if loadBalancerDetail.Protocol != nil {
				loadBalancerDetailMap["protocol"] = loadBalancerDetail.Protocol
			}

			if loadBalancerDetail.Port != nil {
				loadBalancerDetailMap["port"] = loadBalancerDetail.Port
			}

			if loadBalancerDetail.LocationId != nil {
				loadBalancerDetailMap["location_id"] = loadBalancerDetail.LocationId
			}

			if loadBalancerDetail.Domain != nil {
				loadBalancerDetailMap["domain"] = loadBalancerDetail.Domain
			}

			if loadBalancerDetail.Url != nil {
				loadBalancerDetailMap["url"] = loadBalancerDetail.Url
			}

			if loadBalancerDetail.TargetId != nil {
				loadBalancerDetailMap["target_id"] = loadBalancerDetail.TargetId
			}

			if loadBalancerDetail.TargetAddress != nil {
				loadBalancerDetailMap["target_address"] = loadBalancerDetail.TargetAddress
			}

			if loadBalancerDetail.TargetPort != nil {
				loadBalancerDetailMap["target_port"] = loadBalancerDetail.TargetPort
			}

			if loadBalancerDetail.TargetWeight != nil {
				loadBalancerDetailMap["target_weight"] = loadBalancerDetail.TargetWeight
			}

			if loadBalancerDetail.Isolation != nil {
				loadBalancerDetailMap["isolation"] = loadBalancerDetail.Isolation
			}

			if loadBalancerDetail.SecurityGroup != nil {
				loadBalancerDetailMap["security_group"] = loadBalancerDetail.SecurityGroup
			}

			if loadBalancerDetail.LoadBalancerPassToTarget != nil {
				loadBalancerDetailMap["load_balancer_pass_to_target"] = loadBalancerDetail.LoadBalancerPassToTarget
			}

			if loadBalancerDetail.TargetHealth != nil {
				loadBalancerDetailMap["target_health"] = loadBalancerDetail.TargetHealth
			}

			if loadBalancerDetail.Domains != nil {
				loadBalancerDetailMap["domains"] = loadBalancerDetail.Domains
			}

			if loadBalancerDetail.SlaveZone != nil {
				loadBalancerDetailMap["slave_zone"] = loadBalancerDetail.SlaveZone
			}

			if loadBalancerDetail.Zones != nil {
				loadBalancerDetailMap["zones"] = loadBalancerDetail.Zones
			}

			if loadBalancerDetail.SniSwitch != nil {
				loadBalancerDetailMap["sni_switch"] = loadBalancerDetail.SniSwitch
			}

			if loadBalancerDetail.LoadBalancerDomain != nil {
				loadBalancerDetailMap["load_balancer_domain"] = loadBalancerDetail.LoadBalancerDomain
			}

			ids = append(ids, *loadBalancerDetail.LoadBalancerId)
			tmpList = append(tmpList, loadBalancerDetailMap)
		}

		_ = d.Set("load_balancer_detail_set", tmpList)
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
