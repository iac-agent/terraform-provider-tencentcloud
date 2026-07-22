package lighthouse

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudLighthouseDiskConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseDiskConfigRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 列表.zoneFilter 通过 availability 可用区类型: StringRequired: 无。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "待过滤字段",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "过滤值 的 字段。",
						},
					},
				},
			},

			"disk_config_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 云 磁盘 configurations。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability 可用区",
						},
						"disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud 磁盘 类型",
						},
						"disk_sales_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud 磁盘 sale 状态",
						},
						"max_disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum 云 磁盘 大小。",
						},
						"min_disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum 云 磁盘 大小。",
						},
						"disk_step_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cloud 磁盘 increment。",
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

func dataSourceTencentCloudLighthouseDiskConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_lighthouse_disk_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*lighthouse.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := lighthouse.Filter{}
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
		paramMap["filters"] = tmpSet
	}

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var diskConfigSet []*lighthouse.DiskConfig

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseDiskConfigByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		diskConfigSet = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(diskConfigSet))

	if diskConfigSet != nil {
		for _, diskConfig := range diskConfigSet {
			diskConfigMap := map[string]interface{}{}

			if diskConfig.Zone != nil {
				diskConfigMap["zone"] = diskConfig.Zone
			}

			if diskConfig.DiskType != nil {
				diskConfigMap["disk_type"] = diskConfig.DiskType
			}

			if diskConfig.DiskSalesState != nil {
				diskConfigMap["disk_sales_state"] = diskConfig.DiskSalesState
			}

			if diskConfig.MaxDiskSize != nil {
				diskConfigMap["max_disk_size"] = diskConfig.MaxDiskSize
			}

			if diskConfig.MinDiskSize != nil {
				diskConfigMap["min_disk_size"] = diskConfig.MinDiskSize
			}

			if diskConfig.DiskStepSize != nil {
				diskConfigMap["disk_step_size"] = diskConfig.DiskStepSize
			}

			tmpList = append(tmpList, diskConfigMap)
		}

		_ = d.Set("disk_config_set", tmpList)
	}
	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
