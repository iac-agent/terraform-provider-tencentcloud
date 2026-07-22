package gaap

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudGaapAccessRegionsByDestRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapAccessRegionsByDestRegionRead,
		Schema: map[string]*schema.Schema{
			"dest_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Origin 地域",
			},

			"ip_address_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "IP 版本，可以 是 taken 作为 IPv4 或 IPv6，使用 默认值 的 IPv4。",
			},

			"package_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Channel 包 类型，其中 Thunder 表示 standard proxy 组，Accelerator 表示 game accelerator proxy，和 CrossBorder 表示 cross-border proxy。",
			},

			"access_region_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 可用 acceleration 可用区 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 ID。",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Chinese 或 English 名称 地域",
						},
						"concurrent_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "可选 并发 值 数组。",
						},
						"bandwidth_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "可选 带宽 值 数组。",
						},
						"region_area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 的 computer room。",
						},
						"region_area_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域名称 的 computer room。",
						},
						"idc_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 computer room，其中 dc 表示 DataCenter 数据 center 和 ec 表示 EdgeComputing edge 节点。",
						},
						"feature_bitmap": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "类型 computer room，其中 dc 表示 DataCenter 数据 center，ec 表示 功能 bitmap，和 each bit 表示 功能，其中:0，表示that 功能 是 不 支持;1，表示support 对于 此 功能. meaning 的 功能 bitmap 是 作为 follows (从 right 到 left): first bit 支持 4-layer acceleration; second bit 支持 7-layer acceleration; third bit 支持 Http3 访问; fourth bit 支持 IPv6; fifth bit 支持 high-quality BGP 访问; 6th bit 支持 three 网络 访问; 7th bit 支持 QoS acceleration 在 访问 segment.注意：此字段可能返回 null，表示无法获取有效值。 Edge nodes。",
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

func dataSourceTencentCloudGaapAccessRegionsByDestRegionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_gaap_access_regions_by_dest_region.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("dest_region"); ok {
		paramMap["dest_region"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ip_address_version"); ok {
		paramMap["ip_address_version"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("package_type"); ok {
		paramMap["package_type"] = helper.String(v.(string))
	}

	service := GaapService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var accessRegionSet []*gaap.AccessRegionDetial

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeGaapAccessRegionsByDestRegionByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		accessRegionSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(accessRegionSet))
	tmpList := make([]map[string]interface{}, 0, len(accessRegionSet))

	if accessRegionSet != nil {
		for _, accessRegionDetial := range accessRegionSet {
			accessRegionDetialMap := map[string]interface{}{}

			if accessRegionDetial.RegionId != nil {
				accessRegionDetialMap["region_id"] = accessRegionDetial.RegionId
			}

			if accessRegionDetial.RegionName != nil {
				accessRegionDetialMap["region_name"] = accessRegionDetial.RegionName
			}

			if accessRegionDetial.ConcurrentList != nil {
				accessRegionDetialMap["concurrent_list"] = accessRegionDetial.ConcurrentList
			}

			if accessRegionDetial.BandwidthList != nil {
				accessRegionDetialMap["bandwidth_list"] = accessRegionDetial.BandwidthList
			}

			if accessRegionDetial.RegionArea != nil {
				accessRegionDetialMap["region_area"] = accessRegionDetial.RegionArea
			}

			if accessRegionDetial.RegionAreaName != nil {
				accessRegionDetialMap["region_area_name"] = accessRegionDetial.RegionAreaName
			}

			if accessRegionDetial.IDCType != nil {
				accessRegionDetialMap["idc_type"] = accessRegionDetial.IDCType
			}

			if accessRegionDetial.FeatureBitmap != nil {
				accessRegionDetialMap["feature_bitmap"] = accessRegionDetial.FeatureBitmap
			}

			ids = append(ids, *accessRegionDetial.RegionId)
			tmpList = append(tmpList, accessRegionDetialMap)
		}

		_ = d.Set("access_region_set", tmpList)
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
