package gaap

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudGaapDestRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapDestRegionsRead,
		Schema: map[string]*schema.Schema{
			"dest_region_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "来源 Site Area Details List。",
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
							Description: "地域名称",
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
							Description: "Property bitmap，其中 each bit 表示 属性，其中:0，表示that 功能 是 不 支持;1，表示support 对于 此 功能. meaning 的 功能 bitmap 是 作为 follows (从 right 到 left): first bit 支持 4-layer acceleration; second bit 支持 7-layer acceleration; third bit 支持 Http3 访问; fourth bit 支持 IPv6; fifth bit 支持 high-quality BGP 访问; 6th bit 支持 three 网络 访问; 7th bit 支持 QoS acceleration 在 访问 segment.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"support_feature": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Ability 到 访问 regional support注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_type": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "A 列表 网络 types 支持 通过 访问 area，使用 normal indicating support 对于 regular BGP，cn2 indicating premium BGP，triple indicating three networks，和 secure_EIP 表示 自定义 secure EIP。",
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

func dataSourceTencentCloudGaapDestRegionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_gaap_dest_regions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := GaapService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var destRegionSet []*gaap.RegionDetail

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeGaapDestRegions(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}
		destRegionSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(destRegionSet))
	tmpList := make([]map[string]interface{}, 0, len(destRegionSet))

	if destRegionSet != nil {
		for _, regionDetail := range destRegionSet {
			regionDetailMap := map[string]interface{}{}

			if regionDetail.RegionId != nil {
				regionDetailMap["region_id"] = regionDetail.RegionId
			}

			if regionDetail.RegionName != nil {
				regionDetailMap["region_name"] = regionDetail.RegionName
			}

			if regionDetail.RegionArea != nil {
				regionDetailMap["region_area"] = regionDetail.RegionArea
			}

			if regionDetail.RegionAreaName != nil {
				regionDetailMap["region_area_name"] = regionDetail.RegionAreaName
			}

			if regionDetail.IDCType != nil {
				regionDetailMap["idc_type"] = regionDetail.IDCType
			}

			if regionDetail.FeatureBitmap != nil {
				regionDetailMap["feature_bitmap"] = regionDetail.FeatureBitmap
			}

			if regionDetail.SupportFeature != nil {
				supportFeatureMap := map[string]interface{}{}

				if regionDetail.SupportFeature.NetworkType != nil {
					supportFeatureMap["network_type"] = regionDetail.SupportFeature.NetworkType
				}

				regionDetailMap["support_feature"] = []interface{}{supportFeatureMap}
			}

			ids = append(ids, *regionDetail.RegionId)
			tmpList = append(tmpList, regionDetailMap)
		}

		_ = d.Set("dest_region_set", tmpList)
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
