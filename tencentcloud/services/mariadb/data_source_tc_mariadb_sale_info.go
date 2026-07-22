package mariadb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMariadbSaleInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbSaleInfoRead,
		Schema: map[string]*schema.Schema{
			"region_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 sale 地域 info。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域名称(en)。",
						},
						"region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "地域 ID",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域名称(zh)。",
						},
						"zone_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "列表 az 可用区",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "可用区 名称(en)。",
									},
									"zone_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "可用区 ID",
									},
									"zone_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "可用区 名称(zh)。",
									},
									"on_sale": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是 可用区 在 sale。",
									},
								},
							},
						},
						"available_choice": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "可用 可用区 choice。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"master_zone": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "master 可用区",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"zone": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "可用区 名称(en)。",
												},
												"zone_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "可用区 ID",
												},
												"zone_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "可用区 名称(zh)。",
												},
												"on_sale": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是 可用区 在 sale。",
												},
											},
										},
									},
									"slave_zones": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "slave zones。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"zone": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "可用区 名称(en)。",
												},
												"zone_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "可用区 ID",
												},
												"zone_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "可用区 名称(zh)。",
												},
												"on_sale": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是 可用区 在 sale。",
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
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudMariadbSaleInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mariadb_sale_info.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		regionList []*mariadb.RegionInfo
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbSaleInfoByFilter(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}

		regionList = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(regionList))
	tmpList := make([]map[string]interface{}, 0, len(regionList))

	if regionList != nil {
		for _, regionInfo := range regionList {
			regionInfoMap := map[string]interface{}{}

			if regionInfo.Region != nil {
				regionInfoMap["region"] = regionInfo.Region
			}

			if regionInfo.RegionId != nil {
				regionInfoMap["region_id"] = regionInfo.RegionId
			}

			if regionInfo.RegionName != nil {
				regionInfoMap["region_name"] = regionInfo.RegionName
			}

			if regionInfo.ZoneList != nil {
				zoneListList := []interface{}{}
				for _, zoneList := range regionInfo.ZoneList {
					zoneListMap := map[string]interface{}{}

					if zoneList.Zone != nil {
						zoneListMap["zone"] = zoneList.Zone
					}

					if zoneList.ZoneId != nil {
						zoneListMap["zone_id"] = zoneList.ZoneId
					}

					if zoneList.ZoneName != nil {
						zoneListMap["zone_name"] = zoneList.ZoneName
					}

					if zoneList.OnSale != nil {
						zoneListMap["on_sale"] = zoneList.OnSale
					}

					zoneListList = append(zoneListList, zoneListMap)
				}

				regionInfoMap["zone_list"] = zoneListList
			}

			if regionInfo.AvailableChoice != nil {
				availableChoiceList := []interface{}{}
				for _, availableChoice := range regionInfo.AvailableChoice {
					availableChoiceMap := map[string]interface{}{}

					if availableChoice.MasterZone != nil {
						masterZoneList := []interface{}{}
						masterZoneMap := map[string]interface{}{}

						if availableChoice.MasterZone.Zone != nil {
							masterZoneMap["zone"] = availableChoice.MasterZone.Zone
						}

						if availableChoice.MasterZone.ZoneId != nil {
							masterZoneMap["zone_id"] = availableChoice.MasterZone.ZoneId
						}

						if availableChoice.MasterZone.ZoneName != nil {
							masterZoneMap["zone_name"] = availableChoice.MasterZone.ZoneName
						}

						if availableChoice.MasterZone.OnSale != nil {
							masterZoneMap["on_sale"] = availableChoice.MasterZone.OnSale
						}

						masterZoneList = append(masterZoneList, masterZoneMap)
						availableChoiceMap["master_zone"] = masterZoneList
					}

					if availableChoice.SlaveZones != nil {
						slaveZonesList := []interface{}{}
						for _, slaveZones := range availableChoice.SlaveZones {
							slaveZonesMap := map[string]interface{}{}

							if slaveZones.Zone != nil {
								slaveZonesMap["zone"] = slaveZones.Zone
							}

							if slaveZones.ZoneId != nil {
								slaveZonesMap["zone_id"] = slaveZones.ZoneId
							}

							if slaveZones.ZoneName != nil {
								slaveZonesMap["zone_name"] = slaveZones.ZoneName
							}

							if slaveZones.OnSale != nil {
								slaveZonesMap["on_sale"] = slaveZones.OnSale
							}

							slaveZonesList = append(slaveZonesList, slaveZonesMap)
						}

						availableChoiceMap["slave_zones"] = slaveZonesList
					}

					availableChoiceList = append(availableChoiceList, availableChoiceMap)
				}

				regionInfoMap["available_choice"] = availableChoiceList
			}

			ids = append(ids, *regionInfo.Region)
			tmpList = append(tmpList, regionInfoMap)
		}

		_ = d.Set("region_list", tmpList)
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
