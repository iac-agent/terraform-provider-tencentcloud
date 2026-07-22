package region

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	regionv20220627 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/region/v20220627"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRegionsRead,
		Schema: map[string]*schema.Schema{
			"product": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Product 名称 到 查询，e.g. `cvm`. Use `tencentcloud_products` 到 get 可用 product names。",
			},

			"scene": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Scene control 参数. `0` 或 不 集合 表示 do 不 查询 可选 business whitelist; `1` 表示 查询 可选 business whitelist。",
			},

			"region_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "地域 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 identifier，e.g. `ap-guangzhou`。",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域名称，e.g. `South China (Guangzhou)`。",
						},
						"region_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 availability 状态，e.g. `AVAILABLE`。",
						},
						"region_type_m_c": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Console 类型，null 当 called via API。",
						},
						"location_m_c": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 描述 在 different languages。",
						},
						"region_name_m_c": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 描述 displayed 在 console。",
						},
						"region_id_m_c": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 ID 对于 console。",
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

func dataSourceTencentCloudRegionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_regions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = RegionService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("scene"); ok {
		paramMap["Scene"] = helper.IntInt64(v.(int))
	}

	var respData []*regionv20220627.RegionInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRegionsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	regionListData := make([]map[string]interface{}, 0, len(respData))
	for _, regionInfo := range respData {
		regionMap := map[string]interface{}{}
		if regionInfo.Region != nil {
			regionMap["region"] = regionInfo.Region
		}
		if regionInfo.RegionName != nil {
			regionMap["region_name"] = regionInfo.RegionName
		}
		if regionInfo.RegionState != nil {
			regionMap["region_state"] = regionInfo.RegionState
		}
		if regionInfo.RegionTypeMC != nil {
			regionMap["region_type_m_c"] = regionInfo.RegionTypeMC
		}
		if regionInfo.LocationMC != nil {
			regionMap["location_m_c"] = regionInfo.LocationMC
		}
		if regionInfo.RegionNameMC != nil {
			regionMap["region_name_m_c"] = regionInfo.RegionNameMC
		}
		if regionInfo.RegionIdMC != nil {
			regionMap["region_id_m_c"] = regionInfo.RegionIdMC
		}
		regionListData = append(regionListData, regionMap)
	}

	_ = d.Set("region_list", regionListData)

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
