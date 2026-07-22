package lighthouse

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudLighthouseRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseRegionRead,
		Schema: map[string]*schema.Schema{
			"region_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "地域 信息 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域名称",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 描述",
						},
						"region_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 availability 状态",
						},
						"is_china_mainland": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否region 是 在 Chinese mainland。",
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

func dataSourceTencentCloudLighthouseRegionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_lighthouse_region.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var regionSet []*lighthouse.RegionInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseRegionByFilter(ctx, paramMap)
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
		for _, regionInfo := range regionSet {
			regionInfoMap := map[string]interface{}{}

			if regionInfo.Region != nil {
				regionInfoMap["region"] = regionInfo.Region
			}

			if regionInfo.RegionName != nil {
				regionInfoMap["region_name"] = regionInfo.RegionName
			}

			if regionInfo.RegionState != nil {
				regionInfoMap["region_state"] = regionInfo.RegionState
			}

			if regionInfo.IsChinaMainland != nil {
				regionInfoMap["is_china_mainland"] = regionInfo.IsChinaMainland
			}

			ids = append(ids, *regionInfo.Region)
			tmpList = append(tmpList, regionInfoMap)
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
