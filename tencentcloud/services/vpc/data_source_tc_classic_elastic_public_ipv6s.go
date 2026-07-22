package vpc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClassicElasticPublicIpv6s() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClassicElasticPublicIpv6sRead,
		Schema: map[string]*schema.Schema{
			"ip6_address_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "列表 唯一 IDs 该 identify IPV6. IPV6 唯一 ID 是 shaped like `eip-11112222`. Parameters do 不 support specifying both `Ip6AddressIds` 和 `Filters`。",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "upper 限制 对于 `Filters per 请求 是 10，和 upper 限制 对于`过滤器.Values` 是 100. Parameters do 不 support specifying both AddressIds 和 Filters. detailed filtering conditions 是 作为 follows:\n - 地址-ip: 过滤器 according 到 IPV6 IP 地址\n - 网络-interface-ID: 过滤器 according 到 唯一 ID Elastic Network Interface。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "属性名称 如果 there 是 多个 Filters， relationship between Filters 是 logical AND relationship。",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Attribute 值 如果 there 是 多个 Values 在 same 过滤器， relationship between Values under same 过滤器 是 logical OR relationship. 当 值 类型 是 Boolean 类型， 值 可以 是 directly taken 到 字符串 TRUE 或 FALSE。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"address_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "列表 IPV6 details。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID EIP 是 唯一 identifier 的 EIP。",
						},
						"address_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "`EIP` 名称",
						},
						"address_status": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "`EIP` 状态，包括 `CREATING`(Creating),`BINDING`(Binding),`BIND`(Unbinding),`UNBIND`(Unbinding),`OFFLINING`(Releasing),`BIND_ENI`(Binding Suspend Elastic Network Interface)。",
						},
						"address_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "External 网络 IP 地址",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "bound 资源 实例 `ID`. It 可能 是 `CVM`,`NAT`。",
						},
						"created_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "创建时间. It 是 expressed 在 accordance 使用 ISO8601 standard 和 uses UTC 时间. 格式 是: `Y-MM-DDThh:mm:ssZ`。",
						},
						"network_interface_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "bound Elastic Network Interface ID。",
						},
						"private_address_ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Binding resources intranet IP。",
						},
						"is_arrears": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Resource isolation 状态 true 表示 eip 是 在 isolation，false 表示 资源 是 在 non-isolation state。",
						},
						"is_blocked": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Resource blocking 状态 true 表示 eip 是 blocked，false 表示 eip 是 不 blocked。",
						},
						"is_eip_direct_connection": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether eip 支持 pass-through 模式 true 表示 eip 支持 pass-through 模式，false 表示 resources do 不 support pass-through 模式",
						},
						"address_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "EIP 资源 types，包括 CalcIP，WanIP，EIP，AnycastEIP，和 high-defense EIP. Among them: `CalcIP` 表示 device IP,`WanIP` 表示 ordinary 公有 网络 IP,`EIP` 表示 elastic 公有 网络 IP,`AnycastEIP` 表示 accelerated EIP，和 `AntiDDoSEIP` 表示 highly resistant EIP。",
						},
						"cascade_release": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether eip 是 automatically released after unbinding. true 表示 该 eip 将 是 automatically released after unbinding，false 表示 该 eip 将 不 是 automatically released after unbinding。",
						},
						"eip_alg_type": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "类型 协议 opened 通过 EIP ALG。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ftp": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "是否Ftp 协议 Alg 函数 是 已启用",
									},
									"sip": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "是否Sip 协议 Alg 函数 是 已启用",
									},
								},
							},
						},
						"internet_service_provider": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "操作者 信息 的 elastic 公有 网络 IP. Current possible 返回 值 include `CMCC`,`CTCC`,`CUCC`,`BGP`。",
						},
						"local_bgp": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether 本地 带宽 EIP。",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "带宽 值 的 elastic 公有 网络 IP. 注意 该 elastic 公有 IP 的 traditional 账号 types has 无 带宽 attribute 和 值 是 null。",
						},
						"internet_charge_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Network charging model 对于 elastic 公有 网络 IP. 注意 该 elastic 公有 IP 的 traditional 账号 types does 不 have 网络 charging 模式 attribute 和 值 是 blank. 注意：此字段可能返回 null，表示无法获取有效值。 Includes:\n - BANDWIDTH_PREPAID_BY_MONTH: 表示 prepaid monthly 带宽;\n - TRAFFIC_POSTPAID_BY_HOUR: 表示 postpayment per hour 的 流量;\n - BANDWIDTH_POSTPAID_BY_HOUR: 表示 pay-per-hour 带宽;\n - BANDWIDTH_PACKAGE: 表示 shared 带宽 包。",
						},
						"tag_set": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "列表 标签 associated 使用 elastic 公有 IP。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "标签键",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "标签值",
									},
								},
							},
						},
						"deadline_date": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过期时间。",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "实例 类型 EIP binding。",
						},
						"egress": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Static 单个-wire IP 网络 exit。",
						},
						"anti_ddos_package_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "High-defense 包 ID. 当 EIP 类型 是 high-defense EIP，它 返回high-defense 包 ID 到 其中 EIP 是 bound。",
						},
						"renew_flag": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "是否current EIP 是 automatically renewed，此 字段 将 是 displayed 仅 对于 EIP prepaid 通过 monthly 带宽. Examples 的 特定 值 是 作为 follows:\n - NOTIFY_AND_MANUAL_RENEW: Normal renewal\n - NOTIFY_AND_AUTO_RENEW: Automatic renewal\n - DISABLE_NOTIFY_AND_MANUAL_RENEW: No renewal after expiration。",
						},
						"bandwidth_package_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "带宽 包 ID associated 使用 当前 公有 IP. 如果 公有 IP does 不 使用 带宽 packages 对于 charging， 返回 将 是 blank。",
						},
						"un_vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique ID vpc 到 其中 traditional Elastic IPv6 belongs。",
						},
						"dedicated_cluster_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "CDC 唯一 ID。",
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

func dataSourceTencentCloudClassicElasticPublicIpv6sRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_classic_elastic_public_ipv6s.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("ip6_address_ids"); ok {
		ip6AddressIdsList := []*string{}
		ip6AddressIdsSet := v.(*schema.Set).List()
		for i := range ip6AddressIdsSet {
			ip6AddressIds := ip6AddressIdsSet[i].(string)
			ip6AddressIdsList = append(ip6AddressIdsList, helper.String(ip6AddressIds))
		}
		paramMap["Ip6AddressIds"] = ip6AddressIdsList
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*vpc.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := vpc.Filter{}
			if v, ok := filtersMap["name"]; ok {
				filter.Name = helper.String(v.(string))
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

	var respData []*vpc.Address
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClassicElasticPublicIpv6sByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	var ids []string
	addressSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, addressSet := range respData {
			addressSetMap := map[string]interface{}{}

			var addressId string
			if addressSet.AddressId != nil {
				addressSetMap["address_id"] = addressSet.AddressId
				addressId = *addressSet.AddressId
			}

			if addressSet.AddressName != nil {
				addressSetMap["address_name"] = addressSet.AddressName
			}

			if addressSet.AddressStatus != nil {
				addressSetMap["address_status"] = addressSet.AddressStatus
			}

			if addressSet.AddressIp != nil {
				addressSetMap["address_ip"] = addressSet.AddressIp
			}

			if addressSet.InstanceId != nil {
				addressSetMap["instance_id"] = addressSet.InstanceId
			}

			if addressSet.CreatedTime != nil {
				addressSetMap["created_time"] = addressSet.CreatedTime
			}

			if addressSet.NetworkInterfaceId != nil {
				addressSetMap["network_interface_id"] = addressSet.NetworkInterfaceId
			}

			if addressSet.PrivateAddressIp != nil {
				addressSetMap["private_address_ip"] = addressSet.PrivateAddressIp
			}

			if addressSet.IsArrears != nil {
				addressSetMap["is_arrears"] = addressSet.IsArrears
			}

			if addressSet.IsBlocked != nil {
				addressSetMap["is_blocked"] = addressSet.IsBlocked
			}

			if addressSet.IsEipDirectConnection != nil {
				addressSetMap["is_eip_direct_connection"] = addressSet.IsEipDirectConnection
			}

			if addressSet.AddressType != nil {
				addressSetMap["address_type"] = addressSet.AddressType
			}

			if addressSet.CascadeRelease != nil {
				addressSetMap["cascade_release"] = addressSet.CascadeRelease
			}

			eipAlgTypeMap := map[string]interface{}{}

			if addressSet.EipAlgType != nil {
				if addressSet.EipAlgType.Ftp != nil {
					eipAlgTypeMap["ftp"] = addressSet.EipAlgType.Ftp
				}

				if addressSet.EipAlgType.Sip != nil {
					eipAlgTypeMap["sip"] = addressSet.EipAlgType.Sip
				}

				addressSetMap["eip_alg_type"] = []interface{}{eipAlgTypeMap}
			}

			if addressSet.InternetServiceProvider != nil {
				addressSetMap["internet_service_provider"] = addressSet.InternetServiceProvider
			}

			if addressSet.LocalBgp != nil {
				addressSetMap["local_bgp"] = addressSet.LocalBgp
			}

			if addressSet.Bandwidth != nil {
				addressSetMap["bandwidth"] = addressSet.Bandwidth
			}

			if addressSet.InternetChargeType != nil {
				addressSetMap["internet_charge_type"] = addressSet.InternetChargeType
			}

			tagSetList := make([]map[string]interface{}, 0, len(addressSet.TagSet))
			if addressSet.TagSet != nil {
				for _, tagSet := range addressSet.TagSet {
					tagSetMap := map[string]interface{}{}

					if tagSet.Key != nil {
						tagSetMap["key"] = tagSet.Key
					}

					if tagSet.Value != nil {
						tagSetMap["value"] = tagSet.Value
					}

					tagSetList = append(tagSetList, tagSetMap)
				}

				addressSetMap["tag_set"] = tagSetList
			}
			if addressSet.DeadlineDate != nil {
				addressSetMap["deadline_date"] = addressSet.DeadlineDate
			}

			if addressSet.InstanceType != nil {
				addressSetMap["instance_type"] = addressSet.InstanceType
			}

			if addressSet.Egress != nil {
				addressSetMap["egress"] = addressSet.Egress
			}

			if addressSet.AntiDDoSPackageId != nil {
				addressSetMap["anti_ddos_package_id"] = addressSet.AntiDDoSPackageId
			}

			if addressSet.RenewFlag != nil {
				addressSetMap["renew_flag"] = addressSet.RenewFlag
			}

			if addressSet.BandwidthPackageId != nil {
				addressSetMap["bandwidth_package_id"] = addressSet.BandwidthPackageId
			}

			if addressSet.UnVpcId != nil {
				addressSetMap["un_vpc_id"] = addressSet.UnVpcId
			}

			if addressSet.DedicatedClusterId != nil {
				addressSetMap["dedicated_cluster_id"] = addressSet.DedicatedClusterId
			}

			ids = append(ids, addressId)
			addressSetList = append(addressSetList, addressSetMap)
		}

		_ = d.Set("address_set", addressSetList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), addressSetList); e != nil {
			return e
		}
	}

	return nil
}
