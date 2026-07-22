package common

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudAvailabilityRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAvailabilityRegionsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "当 指定，仅 地域 使用 exactly 名称 match 将 是 返回. `默认值` 值 表示 它 consistent 使用 provider 地域",
			},
			"include_unavailable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "A bool variable 表示that 查询 将 include `UNAVAILABLE` regions。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// Computed values.
			"regions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 regions 将 是 exported 和 its every element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 地域，like `ap-guangzhou`。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 地域，like `Guangzhou 地域`。",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "state 的 地域，indicate availability 使用 `AVAILABLE` 和 `UNAVAILABLE` 值。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAvailabilityRegionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_availability_regions.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cvmService := svccvm.NewCvmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	var name string
	var includeUnavailable = false
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	if name == "default" {
		name = meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
	}
	if v, ok := d.GetOkExists("include_unavailable"); ok {
		includeUnavailable = v.(bool)
	}

	var regions []*cvm.RegionInfo
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		regions, errRet = cvmService.DescribeRegions(ctx)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	regionList := make([]map[string]interface{}, 0, len(regions))
	ids := make([]string, 0, len(regions))
	for _, region := range regions {
		if name != "" && name != *region.Region {
			continue
		}
		if !includeUnavailable && *region.RegionState == svccvm.ZONE_STATE_UNAVAILABLE {
			continue
		}
		mapping := map[string]interface{}{
			"name":        region.Region,
			"description": region.RegionName,
			"state":       region.RegionState,
		}
		regionList = append(regionList, mapping)
		ids = append(ids, *region.Region)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("regions", regionList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set regions list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), regionList); err != nil {
			return err
		}
	}

	return nil
}
