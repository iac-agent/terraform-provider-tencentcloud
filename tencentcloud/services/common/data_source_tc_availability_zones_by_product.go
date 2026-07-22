package common

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/api/v20201106"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudAvailabilityZonesByProduct() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAvailabilityZonesByProductRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "当 指定，仅 可用区 使用 exactly 名称 match 将 是 返回。",
			},
			"product": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A 字符串 variable 表示that 查询 将 使用 product 信息。",
			},
			"include_unavailable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "A bool variable 表示that 查询 将 include `UNAVAILABLE` zones。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// Computed values.
			"zones": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 zones 将 是 exported 和 its every element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An 内部 ID 对于 可用区，like `200003`，usually 不 so useful。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 可用区，like `ap-guangzhou-3`。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 可用区，like `Guangzhou 可用区 3`。",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "state 的 可用区，indicate availability 使用 `AVAILABLE` 和 `UNAVAILABLE` 值。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAvailabilityZonesByProductRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_availability_zones.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	apiService := APIService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	var name string
	var product string
	var includeUnavailable = false
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	if v, ok := d.GetOk("product"); ok {
		product = v.(string)
	}
	if v, ok := d.GetOkExists("include_unavailable"); ok {
		includeUnavailable = v.(bool)
	}

	var zones []*api.ZoneInfo
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		zones, errRet = apiService.DescribeZonesWithProduct(ctx, product)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	zoneList := make([]map[string]interface{}, 0, len(zones))
	ids := make([]string, 0, len(zones))
	for _, zone := range zones {
		if name != "" && name != *zone.Zone {
			continue
		}
		if !includeUnavailable && *zone.ZoneState == svccvm.ZONE_STATE_UNAVAILABLE {
			continue
		}
		mapping := map[string]interface{}{
			"id":          zone.ZoneId,
			"name":        zone.Zone,
			"description": zone.ZoneName,
			"state":       zone.ZoneState,
		}
		zoneList = append(zoneList, mapping)
		ids = append(ids, *zone.ZoneId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("zones", zoneList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set zones list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), zoneList); err != nil {
			return err
		}
	}

	return nil
}
