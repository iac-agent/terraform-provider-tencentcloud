package tse

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTseGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTseGatewaysRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 conditions，有效 值:类型,名称,GatewayId,标签,TradeType,InternetPaymode,地域",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤名称",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "过滤值",
						},
					},
				},
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "gateways 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "总数",
						},
						"gateway_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "网关 列表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"gateway_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "网关 ID。",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "状态 网关. May 返回 值: `Creating`，`CreateFailed`，`Running`，`Modifying`，`UpdatingSpec`，`UpdateFailed`，`Deleting`，`DeleteFailed`，`Isolating`。",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "网关 名称",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "网关 类型",
									},
									"gateway_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "网关 版本 Reference 值: `2.4.1`，`2.5.1`。",
									},
									"node_config": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "original 节点 配置",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"specification": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "规格，1c2g|2c4g|4c8g|8c16g。",
												},
												"number": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "节点 数量，2-50。",
												},
											},
										},
									},
									"vpc_config": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "vpc 信息。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "子网 ID. Assign IP 地址 到 引擎 在 VPC 子网。",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "子网 ID. Assign IP 地址 到 引擎 在 VPC 子网。",
												},
											},
										},
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "描述 网关。",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创建时间。",
									},
									"tags": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "标签 信息 的 gateway注意：此字段可能返回 null，表示有效值不可用。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"tag_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "标签键",
												},
												"tag_value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "标签值",
												},
											},
										},
									},
									"enable_cls": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否enable CLS 日志。",
									},
									"trade_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "trade 类型 `0`: postpaid，`1`: Prepaid。",
									},
									"feature_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "product 版本 `TRIAL`，`STANDARD`(默认值)，`PROFESSIONAL`。",
									},
									"internet_max_bandwidth_out": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "公有 网络 outbound 流量 带宽。",
									},
									"auto_renew_flag": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "auto 续费标识，`0`: 默认值 状态，`1`: auto renew，`2`: auto 不 renew。",
									},
									"cur_deadline": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "expire date，对于 prepaid 类型注意：此字段可能返回 null，表示有效值不可用。",
									},
									"isolate_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "isolation 时间，使用 当 网关 是 isolated。",
									},
									"enable_internet": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否open 公有 网络 的 客户端.注意：此字段可能返回 null，表示有效值不可用。",
									},
									"engine_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "引擎 地域 的 网关。",
									},
									"ingress_class_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ingress class 名称",
									},
									"internet_pay_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "trade 类型 internet. `BANDWIDTH`，`TRAFFIC`。",
									},
									"gateway_minor_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "minor 版本 的 网关。",
									},
									"instance_port": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "端口 信息 该 实例 monitors。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"http_port": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "http 端口",
												},
												"https_port": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "https 端口",
												},
											},
										},
									},
									"load_balancer_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "load balance 类型 公有 internet。",
									},
									"public_ip_addresses": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "addresses 的 公有 internet。",
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

func dataSourceTencentCloudTseGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tse_gateways.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*tse.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := tse.Filter{}
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
		paramMap["filters"] = tmpSet
	}

	service := TseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var result *tse.ListCloudNativeAPIGatewayResult
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeTseGatewaysByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result.GatewayList))
	listCloudNativeAPIGatewayResultMap := map[string]interface{}{}
	if result != nil {
		if result.TotalCount != nil {
			listCloudNativeAPIGatewayResultMap["total_count"] = result.TotalCount
		}

		if result.GatewayList != nil {
			gatewayListList := []interface{}{}
			for _, gatewayList := range result.GatewayList {
				gatewayListMap := map[string]interface{}{}

				if gatewayList.GatewayId != nil {
					gatewayListMap["gateway_id"] = gatewayList.GatewayId
				}

				if gatewayList.Status != nil {
					gatewayListMap["status"] = gatewayList.Status
				}

				if gatewayList.Name != nil {
					gatewayListMap["name"] = gatewayList.Name
				}

				if gatewayList.Type != nil {
					gatewayListMap["type"] = gatewayList.Type
				}

				if gatewayList.GatewayVersion != nil {
					gatewayListMap["gateway_version"] = gatewayList.GatewayVersion
				}

				if gatewayList.NodeConfig != nil {
					nodeConfigMap := map[string]interface{}{}

					if gatewayList.NodeConfig.Specification != nil {
						nodeConfigMap["specification"] = gatewayList.NodeConfig.Specification
					}

					if gatewayList.NodeConfig.Number != nil {
						nodeConfigMap["number"] = gatewayList.NodeConfig.Number
					}

					gatewayListMap["node_config"] = []interface{}{nodeConfigMap}
				}

				if gatewayList.VpcConfig != nil {
					vpcConfigMap := map[string]interface{}{}

					if gatewayList.VpcConfig.VpcId != nil {
						vpcConfigMap["vpc_id"] = gatewayList.VpcConfig.VpcId
					}

					if gatewayList.VpcConfig.SubnetId != nil {
						vpcConfigMap["subnet_id"] = gatewayList.VpcConfig.SubnetId
					}

					gatewayListMap["vpc_config"] = []interface{}{vpcConfigMap}
				}

				if gatewayList.Description != nil {
					gatewayListMap["description"] = gatewayList.Description
				}

				if gatewayList.CreateTime != nil {
					gatewayListMap["create_time"] = gatewayList.CreateTime
				}

				if gatewayList.Tags != nil {
					tagsList := []interface{}{}
					for _, tags := range gatewayList.Tags {
						tagsMap := map[string]interface{}{}

						if tags.TagKey != nil {
							tagsMap["tag_key"] = tags.TagKey
						}

						if tags.TagValue != nil {
							tagsMap["tag_value"] = tags.TagValue
						}

						tagsList = append(tagsList, tagsMap)
					}

					gatewayListMap["tags"] = tagsList
				}

				if gatewayList.EnableCls != nil {
					gatewayListMap["enable_cls"] = gatewayList.EnableCls
				}

				if gatewayList.TradeType != nil {
					gatewayListMap["trade_type"] = gatewayList.TradeType
				}

				if gatewayList.FeatureVersion != nil {
					gatewayListMap["feature_version"] = gatewayList.FeatureVersion
				}

				if gatewayList.InternetMaxBandwidthOut != nil {
					gatewayListMap["internet_max_bandwidth_out"] = gatewayList.InternetMaxBandwidthOut
				}

				if gatewayList.AutoRenewFlag != nil {
					gatewayListMap["auto_renew_flag"] = gatewayList.AutoRenewFlag
				}

				if gatewayList.CurDeadline != nil {
					gatewayListMap["cur_deadline"] = gatewayList.CurDeadline
				}

				if gatewayList.IsolateTime != nil {
					gatewayListMap["isolate_time"] = gatewayList.IsolateTime
				}

				if gatewayList.EnableInternet != nil {
					gatewayListMap["enable_internet"] = gatewayList.EnableInternet
				}

				if gatewayList.EngineRegion != nil {
					gatewayListMap["engine_region"] = gatewayList.EngineRegion
				}

				if gatewayList.IngressClassName != nil {
					gatewayListMap["ingress_class_name"] = gatewayList.IngressClassName
				}

				if gatewayList.InternetPayMode != nil {
					gatewayListMap["internet_pay_mode"] = gatewayList.InternetPayMode
				}

				if gatewayList.GatewayMinorVersion != nil {
					gatewayListMap["gateway_minor_version"] = gatewayList.GatewayMinorVersion
				}

				if gatewayList.InstancePort != nil {
					instancePortMap := map[string]interface{}{}

					if gatewayList.InstancePort.HttpPort != nil {
						instancePortMap["http_port"] = gatewayList.InstancePort.HttpPort
					}

					if gatewayList.InstancePort.HttpsPort != nil {
						instancePortMap["https_port"] = gatewayList.InstancePort.HttpsPort
					}

					gatewayListMap["instance_port"] = []interface{}{instancePortMap}
				}

				if gatewayList.LoadBalancerType != nil {
					gatewayListMap["load_balancer_type"] = gatewayList.LoadBalancerType
				}

				if gatewayList.PublicIpAddresses != nil {
					gatewayListMap["public_ip_addresses"] = gatewayList.PublicIpAddresses
				}

				gatewayListList = append(gatewayListList, gatewayListMap)
				ids = append(ids, *gatewayList.GatewayId)
			}

			listCloudNativeAPIGatewayResultMap["gateway_list"] = gatewayListList
		}

		_ = d.Set("result", []interface{}{listCloudNativeAPIGatewayResultMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), listCloudNativeAPIGatewayResultMap); e != nil {
			return e
		}
	}
	return nil
}
