package region

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	regionv20220627 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/region/v20220627"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudZonesRead,
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

			"zone_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "可用区 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用区 名称，e.g. `ap-guangzhou-3`。",
						},
						"zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用区 描述，e.g. `Guangzhou 可用区 3`。",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用区 ID",
						},
						"zone_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用区 状态，`AVAILABLE` 或 `UNAVAILABLE`。",
						},
						"parent_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parent 可用区 identifier。",
						},
						"parent_zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parent 可用区 ID。",
						},
						"parent_zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parent 可用区 描述",
						},
						"zone_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用区 类型",
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

func dataSourceTencentCloudZonesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_zones.read")()
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

	var respData []*regionv20220627.ZoneInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeZonesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	zoneListData := make([]map[string]interface{}, 0, len(respData))
	for _, zoneInfo := range respData {
		zoneMap := map[string]interface{}{}
		if zoneInfo.Zone != nil {
			zoneMap["zone"] = zoneInfo.Zone
		}
		if zoneInfo.ZoneName != nil {
			zoneMap["zone_name"] = zoneInfo.ZoneName
		}
		if zoneInfo.ZoneId != nil {
			zoneMap["zone_id"] = zoneInfo.ZoneId
		}
		if zoneInfo.ZoneState != nil {
			zoneMap["zone_state"] = zoneInfo.ZoneState
		}
		if zoneInfo.ParentZone != nil {
			zoneMap["parent_zone"] = zoneInfo.ParentZone
		}
		if zoneInfo.ParentZoneId != nil {
			zoneMap["parent_zone_id"] = zoneInfo.ParentZoneId
		}
		if zoneInfo.ParentZoneName != nil {
			zoneMap["parent_zone_name"] = zoneInfo.ParentZoneName
		}
		if zoneInfo.ZoneType != nil {
			zoneMap["zone_type"] = zoneInfo.ZoneType
		}
		zoneListData = append(zoneListData, zoneMap)
	}

	_ = d.Set("zone_list", zoneListData)

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
