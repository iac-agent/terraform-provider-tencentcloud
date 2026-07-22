package ckafka

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCkafkaZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaZoneRead,
		Schema: map[string]*schema.Schema{
			"cdc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "cdc professional 集群 business 参数。",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "查询 结果 complex 对象 entity。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "可用区 列表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "可用区 ID",
									},
									"is_internal_app": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "内部 APP。",
									},
									"app_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "app ID。",
									},
									"flag": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "flag。",
									},
									"zone_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "可用区 名称",
									},
									"zone_status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "可用区 状态",
									},
									"exflag": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "extra flag。",
									},
									"sold_out": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "json 对象，键 是 model，值 true 是 sold out，false 是 不 sold out。",
									},
									"sales_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Standard Edition Sold Out Information。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"flag": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Manually 集合 flags。",
												},
												"version": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ckakfa 版本(1.1.1/2.4.2/0.10.2)。",
												},
												"platform": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Professional Edition，Standard Edition flag。",
												},
												"sold_out": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "sold out flag: true sold out。",
												},
											},
										},
									},
								},
							},
						},
						"max_buy_instance_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最大purchased 实例。",
						},
						"max_bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum purchased 带宽 在 Mbs。",
						},
						"unit_price": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Postpaid 单位 价格。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"real_total_cost": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "discount 价格。",
									},
									"total_cost": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "original 价格。",
									},
								},
							},
						},
						"message_price": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Postpaid 消息 单位 价格。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"real_total_cost": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "discount 价格。",
									},
									"total_cost": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "original 价格。",
									},
								},
							},
						},
						"cluster_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "用户 exclusive 集群 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "集群 ID",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ClusterName。",
									},
									"max_disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "largest 磁盘 在 集群，（GB）。",
									},
									"max_band_width": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum 集群 带宽 在 MBs。",
									},
									"available_disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "当前 可用 磁盘 的 集群，（GB）。",
									},
									"available_band_width": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "当前 可用 带宽 的 集群 在 MBs。",
									},
									"zone_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Availability 可用区 到 其中 集群 belongs，indicating availability 可用区 到 其中 集群 belongs。",
									},
									"zone_ids": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
										Computed:    true,
										Description: "availability 可用区 其中 集群 节点 是 located. 如果 集群 是 cross-availability 可用区 集群，它 includes 多个 availability zones 其中 集群 节点 是 located。",
									},
								},
							},
						},
						"standard": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Purchase Standard Edition Configuration。",
						},
						"standard_s2": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Standard Edition S2 配置。",
						},
						"profession": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Professional Edition 配置。",
						},
						"physical": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Physical Exclusive Edition Configuration。",
						},
						"public_network": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public 网络 带宽。",
						},
						"public_network_limit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public 网络 带宽 配置。",
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

func dataSourceTencentCloudCkafkaZoneRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ckafka_zone.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cdc_id"); ok {
		paramMap["CdcId"] = helper.String(v.(string))
	}

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var result *ckafka.ZoneResponse

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeCkafkaCkafkaZoneByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	zoneResponseMapList := make([]interface{}, 0)
	if result != nil {
		zoneResponseMap := map[string]interface{}{}

		if result.ZoneList != nil {
			zoneListList := []interface{}{}
			for _, zoneList := range result.ZoneList {
				zoneListMap := map[string]interface{}{}

				if zoneList.ZoneId != nil {
					ids = append(ids, *zoneList.ZoneId)
					zoneListMap["zone_id"] = zoneList.ZoneId
				}

				if zoneList.IsInternalApp != nil {
					zoneListMap["is_internal_app"] = zoneList.IsInternalApp
				}

				if zoneList.AppId != nil {
					zoneListMap["app_id"] = zoneList.AppId
				}

				if zoneList.Flag != nil {
					zoneListMap["flag"] = zoneList.Flag
				}

				if zoneList.ZoneName != nil {
					zoneListMap["zone_name"] = zoneList.ZoneName
				}

				if zoneList.ZoneStatus != nil {
					zoneListMap["zone_status"] = zoneList.ZoneStatus
				}

				if zoneList.Exflag != nil {
					zoneListMap["exflag"] = zoneList.Exflag
				}

				if zoneList.SoldOut != nil {
					zoneListMap["sold_out"] = zoneList.SoldOut
				}

				if zoneList.SalesInfo != nil {
					salesInfoList := []interface{}{}
					for _, salesInfo := range zoneList.SalesInfo {
						salesInfoMap := map[string]interface{}{}

						if salesInfo.Flag != nil {
							salesInfoMap["flag"] = salesInfo.Flag
						}

						if salesInfo.Version != nil {
							salesInfoMap["version"] = salesInfo.Version
						}

						if salesInfo.Platform != nil {
							salesInfoMap["platform"] = salesInfo.Platform
						}

						if salesInfo.SoldOut != nil {
							salesInfoMap["sold_out"] = salesInfo.SoldOut
						}

						salesInfoList = append(salesInfoList, salesInfoMap)
					}

					zoneListMap["sales_info"] = salesInfoList
				}

				zoneListList = append(zoneListList, zoneListMap)
			}

			zoneResponseMap["zone_list"] = zoneListList
		}

		if result.MaxBuyInstanceNum != nil {
			zoneResponseMap["max_buy_instance_num"] = result.MaxBuyInstanceNum
		}

		if result.MaxBandwidth != nil {
			zoneResponseMap["max_bandwidth"] = result.MaxBandwidth
		}

		if result.UnitPrice != nil {
			unitPriceMap := map[string]interface{}{}

			if result.UnitPrice.RealTotalCost != nil {
				unitPriceMap["real_total_cost"] = result.UnitPrice.RealTotalCost
			}

			if result.UnitPrice.TotalCost != nil {
				unitPriceMap["total_cost"] = result.UnitPrice.TotalCost
			}

			zoneResponseMap["unit_price"] = []interface{}{unitPriceMap}
		}

		if result.MessagePrice != nil {
			messagePriceMap := map[string]interface{}{}

			if result.MessagePrice.RealTotalCost != nil {
				messagePriceMap["real_total_cost"] = result.MessagePrice.RealTotalCost
			}

			if result.MessagePrice.TotalCost != nil {
				messagePriceMap["total_cost"] = result.MessagePrice.TotalCost
			}

			zoneResponseMap["message_price"] = []interface{}{messagePriceMap}
		}

		if result.ClusterInfo != nil {
			clusterInfoList := []interface{}{}
			for _, clusterInfo := range result.ClusterInfo {
				clusterInfoMap := map[string]interface{}{}

				if clusterInfo.ClusterId != nil {
					clusterInfoMap["cluster_id"] = clusterInfo.ClusterId
				}

				if clusterInfo.ClusterName != nil {
					clusterInfoMap["cluster_name"] = clusterInfo.ClusterName
				}

				if clusterInfo.MaxDiskSize != nil {
					clusterInfoMap["max_disk_size"] = clusterInfo.MaxDiskSize
				}

				if clusterInfo.MaxBandWidth != nil {
					clusterInfoMap["max_band_width"] = clusterInfo.MaxBandWidth
				}

				if clusterInfo.AvailableDiskSize != nil {
					clusterInfoMap["available_disk_size"] = clusterInfo.AvailableDiskSize
				}

				if clusterInfo.AvailableBandWidth != nil {
					clusterInfoMap["available_band_width"] = clusterInfo.AvailableBandWidth
				}

				if clusterInfo.ZoneId != nil {
					clusterInfoMap["zone_id"] = clusterInfo.ZoneId
				}

				if clusterInfo.ZoneIds != nil {
					clusterInfoMap["zone_ids"] = clusterInfo.ZoneIds
				}

				clusterInfoList = append(clusterInfoList, clusterInfoMap)
			}

			zoneResponseMap["cluster_info"] = clusterInfoList
		}

		if result.Standard != nil {
			zoneResponseMap["standard"] = result.Standard
		}

		if result.StandardS2 != nil {
			zoneResponseMap["standard_s2"] = result.StandardS2
		}

		if result.Profession != nil {
			zoneResponseMap["profession"] = result.Profession
		}

		if result.Physical != nil {
			zoneResponseMap["physical"] = result.Physical
		}

		if result.PublicNetwork != nil {
			zoneResponseMap["public_network"] = result.PublicNetwork
		}

		if result.PublicNetworkLimit != nil {
			zoneResponseMap["public_network_limit"] = result.PublicNetworkLimit
		}
		zoneResponseMapList = append(zoneResponseMapList, zoneResponseMap)
		_ = d.Set("result", zoneResponseMapList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), zoneResponseMapList); e != nil {
			return e
		}
	}
	return nil
}
