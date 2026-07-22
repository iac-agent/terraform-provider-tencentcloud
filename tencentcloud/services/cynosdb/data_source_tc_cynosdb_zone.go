package cynosdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCynosdbZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbZoneRead,
		Schema: map[string]*schema.Schema{
			"include_virtual_zones": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "是否包含虚拟区域。",
			},

			"show_permission": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "是否显示该区域下的所有可用区域，并显示用户每个可用区域的权限。",
			},

			"region_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "地区信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地区英文。",
						},
						"region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "区域 ID。",
						},
						"region_zh": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地区中文名称。",
						},
						"zone_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "可供出售的区域列表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "区域名称为英文。",
									},
									"zone_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "区域 ID。",
									},
									"zone_zh": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "中文区域名称。",
									},
									"is_support_serverless": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "是否支持无服务器集群，0：不支持，1：支持。",
									},
									"is_support_normal": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "是否支持普通集群，0：不支持，1：支持。",
									},
									"physical_zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "物理区域。",
									},
									"has_permission": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "用户是否具有区域权限注：该字段可能返回null，表示取不到有效值。",
									},
									"is_whole_rdma_zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "是 Rdma 区域。",
									},
								},
							},
						},
						"db_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "数据库类型。",
						},
						"modules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "区域模块支持。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_disable": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "区域是否出售，可选值：是、否。",
									},
									"module_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "模块名称。",
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

func dataSourceTencentCloudCynosdbZoneRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cynosdb_zone.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("include_virtual_zones"); v != nil {
		paramMap["IncludeVirtualZones"] = helper.Bool(v.(bool))
	}

	if v, _ := d.GetOk("show_permission"); v != nil {
		paramMap["ShowPermission"] = helper.Bool(v.(bool))
	}

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var regionSet []*cynosdb.SaleRegion
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbZoneByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		regionSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(regionSet))
	tmpList := make([]map[string]interface{}, 0, len(regionSet))
	if regionSet != nil {
		for _, saleRegion := range regionSet {
			saleRegionMap := map[string]interface{}{}

			if saleRegion.Region != nil {
				saleRegionMap["region"] = saleRegion.Region
			}

			if saleRegion.RegionId != nil {
				saleRegionMap["region_id"] = saleRegion.RegionId
			}

			if saleRegion.RegionZh != nil {
				saleRegionMap["region_zh"] = saleRegion.RegionZh
			}

			if saleRegion.ZoneSet != nil {
				zoneSetList := []interface{}{}
				for _, zoneSet := range saleRegion.ZoneSet {
					zoneSetMap := map[string]interface{}{}

					if zoneSet.Zone != nil {
						zoneSetMap["zone"] = zoneSet.Zone
					}

					if zoneSet.ZoneId != nil {
						zoneSetMap["zone_id"] = zoneSet.ZoneId
					}

					if zoneSet.ZoneZh != nil {
						zoneSetMap["zone_zh"] = zoneSet.ZoneZh
					}

					if zoneSet.IsSupportServerless != nil {
						zoneSetMap["is_support_serverless"] = zoneSet.IsSupportServerless
					}

					if zoneSet.IsSupportNormal != nil {
						zoneSetMap["is_support_normal"] = zoneSet.IsSupportNormal
					}

					if zoneSet.PhysicalZone != nil {
						zoneSetMap["physical_zone"] = zoneSet.PhysicalZone
					}

					if zoneSet.HasPermission != nil {
						zoneSetMap["has_permission"] = zoneSet.HasPermission
					}

					if zoneSet.IsWholeRdmaZone != nil {
						zoneSetMap["is_whole_rdma_zone"] = zoneSet.IsWholeRdmaZone
					}

					zoneSetList = append(zoneSetList, zoneSetMap)
				}

				saleRegionMap["zone_set"] = zoneSetList
			}

			if saleRegion.DbType != nil {
				saleRegionMap["db_type"] = saleRegion.DbType
			}

			if saleRegion.Modules != nil {
				modulesList := []interface{}{}
				for _, modules := range saleRegion.Modules {
					modulesMap := map[string]interface{}{}

					if modules.IsDisable != nil {
						modulesMap["is_disable"] = modules.IsDisable
					}

					if modules.ModuleName != nil {
						modulesMap["module_name"] = modules.ModuleName
					}

					modulesList = append(modulesList, modulesMap)
				}

				saleRegionMap["modules"] = modulesList
			}

			ids = append(ids, *saleRegion.Region)
			tmpList = append(tmpList, saleRegionMap)
		}

		_ = d.Set("region_set", tmpList)
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
