package dc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDcAccessPoints() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcAccessPointsRead,
		Schema: map[string]*schema.Schema{
			"region_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Access point 地域，which can be queried through `DescribeRegions`.You can call `DescribeRegions` to get the 地域 ID。",
			},

			"access_point_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Access point information。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_point_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access point 名称",
						},
						"access_point_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique access point ID。",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access point 状态 有效值：available，unavailable。",
						},
						"location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access point location。",
						},
						"line_operator": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "列表 ISPs supported by access point。",
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 地域 that manages the access point。",
						},
						"available_port_type": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Available 端口 类型 at the access point. 有效值：1000BASE-T: gigabit electrical 端口; 1000BASE-LX: 10 km gigabit single-模式 optical 端口; 1000BASE-ZX: 80 km gigabit single-模式 optical 端口; 10GBASE-LR: 10 km 10-gigabit single-模式 optical 端口; 10GBASE-ZR: 80 km 10-gigabit single-模式 optical 端口; 10GBASE-LH: 40 km 10-gigabit single-模式 optical 端口; 100GBASE-LR4: 10 km 100-gigabit single-模式 optical portfiber optic 端口Note: this field may return `null`，indicating that no valid 值 is obtained。",
						},
						"coordinate": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Latitude and longitude of the access pointNote: this field may return `null`，indicating that no valid values can be obtained。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lat": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Latitude。",
									},
									"lng": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Longitude。",
									},
								},
							},
						},
						"city": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "City where the access point is locatedNote: this field may return `null`，indicating that no valid values can be obtained。",
						},
						"area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access point regionNote: this field may return `null`，indicating that no valid values can be obtained。",
						},
						"access_point_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access point 类型 有效值：`VXLAN`，`QCPL`，and `QCAR`.Note: this field may return `null`，indicating that no valid values can be obtained。",
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

func dataSourceTencentCloudDcAccessPointsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dc_access_points.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("region_id"); ok {
		paramMap["RegionId"] = helper.String(v.(string))
	}

	service := DcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var accessPointSet []*dc.AccessPoint

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDcAccessPointsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		accessPointSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(accessPointSet))
	tmpList := make([]map[string]interface{}, 0, len(accessPointSet))

	if accessPointSet != nil {
		for _, accessPoint := range accessPointSet {
			accessPointMap := map[string]interface{}{}

			if accessPoint.AccessPointName != nil {
				accessPointMap["access_point_name"] = accessPoint.AccessPointName
			}

			if accessPoint.AccessPointId != nil {
				accessPointMap["access_point_id"] = accessPoint.AccessPointId
			}

			if accessPoint.State != nil {
				accessPointMap["state"] = accessPoint.State
			}

			if accessPoint.Location != nil {
				accessPointMap["location"] = accessPoint.Location
			}

			if accessPoint.LineOperator != nil {
				accessPointMap["line_operator"] = accessPoint.LineOperator
			}

			if accessPoint.RegionId != nil {
				accessPointMap["region_id"] = accessPoint.RegionId
			}

			if accessPoint.AvailablePortType != nil {
				accessPointMap["available_port_type"] = accessPoint.AvailablePortType
			}

			if accessPoint.Coordinate != nil {
				coordinateMap := map[string]interface{}{}

				if accessPoint.Coordinate.Lat != nil {
					coordinateMap["lat"] = accessPoint.Coordinate.Lat
				}

				if accessPoint.Coordinate.Lng != nil {
					coordinateMap["lng"] = accessPoint.Coordinate.Lng
				}

				accessPointMap["coordinate"] = []interface{}{coordinateMap}
			}

			if accessPoint.City != nil {
				accessPointMap["city"] = accessPoint.City
			}

			if accessPoint.Area != nil {
				accessPointMap["area"] = accessPoint.Area
			}

			if accessPoint.AccessPointType != nil {
				accessPointMap["access_point_type"] = accessPoint.AccessPointType
			}

			ids = append(ids, *accessPoint.AccessPointId)
			tmpList = append(tmpList, accessPointMap)
		}

		_ = d.Set("access_point_set", tmpList)
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
