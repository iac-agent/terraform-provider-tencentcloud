package clb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClbResources() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbResourcesRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤查询可用区资源列表，具体如下：可用区 - 字符串 - 可选 - 按可用区过滤，如 ap-guangzhou-1。 isp -- 字符串 - 可选 - 按 ISP 过滤。值：BGP、CMCC、CUCC 和 CTCC。",
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

			"zone_resource_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "可用区支持的资源列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"master_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "主AZ，如ap-guangzhou-1。",
						},
						"resource_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "资源列表。注意：该字段可能返回null，表示取不到有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "具体ISP资源信息，取值：CMCC、CUCC、CTCC、BGP、INTERNAL。",
									},
									"isp": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ISP信息，如CMCC、CUCC、CTCC、BGP、INTERNAL等。",
									},
									"availability_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "可用资源。注意：该字段可能返回null，表示取不到有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "具体ISP资源信息。值：CMCC、CUCC、CTCC、BGP。",
												},
												"availability": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "资源是否可用。值：可用、不可用。",
												},
											},
										},
									},
								},
							},
						},
						"slave_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "辅助AZ，例如ap-guangzhou-2。注意：该字段可能返回null，表示取不到有效值。",
						},
						"ip_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP版本。值：IPv4、IPv6 和 IPv6_Nat。",
						},
						"zone_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AZ区域，如ap-广州。",
						},
						"local_zone": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "该AZ是否为LocalZone。值：真、假。",
						},
						"zone_resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "区域内的资源类型。价值观：共享、独占。",
						},
						"edge_zone": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "该AZ是否为边缘区域。值：真、假。",
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

func dataSourceTencentCloudClbResourcesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clb_resources.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
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

	var zoneResourceSet []*clb.ZoneResource

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClbResourcesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		zoneResourceSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(zoneResourceSet))
	tmpList := make([]map[string]interface{}, 0, len(zoneResourceSet))

	if zoneResourceSet != nil {
		for _, zoneResource := range zoneResourceSet {
			zoneResourceMap := map[string]interface{}{}

			if zoneResource.MasterZone != nil {
				zoneResourceMap["master_zone"] = zoneResource.MasterZone
			}

			if zoneResource.ResourceSet != nil {
				resourceSetList := []interface{}{}
				for _, resourceSet := range zoneResource.ResourceSet {
					resourceSetMap := map[string]interface{}{}

					if resourceSet.Type != nil {
						resourceSetMap["type"] = resourceSet.Type
					}

					if resourceSet.Isp != nil {
						resourceSetMap["isp"] = resourceSet.Isp
					}

					if resourceSet.AvailabilitySet != nil {
						availabilitySetList := []interface{}{}
						for _, availabilitySet := range resourceSet.AvailabilitySet {
							availabilitySetMap := map[string]interface{}{}

							if availabilitySet.Type != nil {
								availabilitySetMap["type"] = availabilitySet.Type
							}

							if availabilitySet.Availability != nil {
								availabilitySetMap["availability"] = availabilitySet.Availability
							}

							availabilitySetList = append(availabilitySetList, availabilitySetMap)
						}

						resourceSetMap["availability_set"] = availabilitySetList
					}

					resourceSetList = append(resourceSetList, resourceSetMap)
				}

				zoneResourceMap["resource_set"] = resourceSetList
			}

			if zoneResource.SlaveZone != nil {
				zoneResourceMap["slave_zone"] = zoneResource.SlaveZone
			}

			if zoneResource.IPVersion != nil {
				zoneResourceMap["ip_version"] = zoneResource.IPVersion
			}

			if zoneResource.ZoneRegion != nil {
				zoneResourceMap["zone_region"] = zoneResource.ZoneRegion
			}

			if zoneResource.LocalZone != nil {
				zoneResourceMap["local_zone"] = zoneResource.LocalZone
			}

			if zoneResource.ZoneResourceType != nil {
				zoneResourceMap["zone_resource_type"] = zoneResource.ZoneResourceType
			}

			if zoneResource.EdgeZone != nil {
				zoneResourceMap["edge_zone"] = zoneResource.EdgeZone
			}

			ids = append(ids, *zoneResource.MasterZone)
			tmpList = append(tmpList, zoneResourceMap)
		}

		_ = d.Set("zone_resource_set", tmpList)
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
