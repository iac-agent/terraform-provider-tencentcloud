package cfs

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCfsAvailableZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCfsAvailableZoneRead,
		Schema: map[string]*schema.Schema{
			"region_zones": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Information such as resource availability in each AZ and the supported storage classes and protocols。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域名称，such as `ap-beijing`。",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域名称，such as `bj`。",
						},
						"region_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 availability. If a 地域 has at least one AZ where resources are purchasable，this 值 will be AVAILABLE; otherwise，it will be UNAVAILABLE。",
						},
						"zones": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "数组 AZs。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "AZ 名称",
									},
									"zone_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "AZ ID。",
									},
									"zone_cn_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Chinese 名称 an AZ。",
									},
									"types": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "数组 classes。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"protocols": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "协议 and sale details。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"sale_status": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "	Sale 状态 有效值：sale_out (sold out)，saling (purchasable)，no_saling (non-purchasable)。",
															},
															"protocol": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "协议 类型 有效值：NFS，CIFS。",
															},
														},
													},
												},
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Storage class. 有效值：SD (standard storage) and HP (high-performance storage)。",
												},
												"prepayment": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "表示是否prepaid is supported. true: yes; false: no。",
												},
											},
										},
									},
									"zone_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Chinese and English names of an AZ。",
									},
								},
							},
						},
						"region_cn_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 chinese 名称，such as `Guangzhou`。",
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

func dataSourceTencentCloudCfsAvailableZoneRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cfs_available_zone.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CfsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var regionZones []*cfs.AvailableRegion

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCfsAvailableZoneByFilter(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}
		regionZones = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(regionZones))
	tmpList := make([]map[string]interface{}, 0, len(regionZones))

	for _, availableRegion := range regionZones {
		availableRegionMap := map[string]interface{}{}

		if availableRegion.Region != nil {
			availableRegionMap["region"] = availableRegion.Region
		}

		if availableRegion.RegionName != nil {
			availableRegionMap["region_name"] = availableRegion.RegionName
		}

		if availableRegion.RegionStatus != nil {
			availableRegionMap["region_status"] = availableRegion.RegionStatus
		}

		if availableRegion.Zones != nil {
			zonesList := []interface{}{}
			for _, zones := range availableRegion.Zones {
				zonesMap := map[string]interface{}{}

				if zones.Zone != nil {
					zonesMap["zone"] = zones.Zone
				}

				if zones.ZoneId != nil {
					zonesMap["zone_id"] = zones.ZoneId
				}

				if zones.ZoneCnName != nil {
					zonesMap["zone_cn_name"] = zones.ZoneCnName
				}

				if zones.Types != nil {
					typesList := []interface{}{}
					for _, types := range zones.Types {
						typesMap := map[string]interface{}{}

						if types.Protocols != nil {
							protocolsList := []interface{}{}
							for _, protocols := range types.Protocols {
								protocolsMap := map[string]interface{}{}

								if protocols.SaleStatus != nil {
									protocolsMap["sale_status"] = protocols.SaleStatus
								}

								if protocols.Protocol != nil {
									protocolsMap["protocol"] = protocols.Protocol
								}

								protocolsList = append(protocolsList, protocolsMap)
							}

							typesMap["protocols"] = protocolsList
						}

						if types.Type != nil {
							typesMap["type"] = types.Type
						}

						if types.Prepayment != nil {
							typesMap["prepayment"] = types.Prepayment
						}

						typesList = append(typesList, typesMap)
					}

					zonesMap["types"] = typesList
				}

				if zones.ZoneName != nil {
					zonesMap["zone_name"] = zones.ZoneName
				}

				zonesList = append(zonesList, zonesMap)
			}

			availableRegionMap["zones"] = zonesList
		}

		if availableRegion.RegionCnName != nil {
			availableRegionMap["region_cn_name"] = availableRegion.RegionCnName
		}
		ids = append(ids, *availableRegion.Region)
		tmpList = append(tmpList, availableRegionMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("region_zones", tmpList)
	if err != nil {
		return err
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
