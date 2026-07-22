package cynosdb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTencentCloudCynosdbZoneConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCynosdbZoneConfigRead,
		Schema: map[string]*schema.Schema{
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// Computed values
			"list": {
				Type:        schema.TypeList,
				Description: "区域列表。每个元素包含以下属性：",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例CPU，单位：核。",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例内存，单位：GB。",
						},
						"max_storage_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例的最大可用存储空间，单位GB。",
						},
						"min_storage_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例的最小可用存储空间，单位：GB。",
						},
						"machine_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "机器类型。",
						},
						"max_io_bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最大 io 带宽。",
						},
						"zone_stock_infos": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "区域库存信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "可用区。",
									},
									"has_stock": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "有库存。",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCynosdbZoneConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cynosdb_zone_config.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region

	instanceSpecSet, err := service.DescribeRedisZoneConfig(ctx)
	if err != nil {
		return fmt.Errorf("api[DescribeRedisZoneConfig]fail, return %s", err.Error())
	}

	result := make([]map[string]interface{}, 0)

	for _, instanceSpec := range instanceSpecSet {
		resultItem := make(map[string]interface{})
		resultItem["cpu"] = *instanceSpec.Cpu
		resultItem["memory"] = *instanceSpec.Memory
		resultItem["max_storage_size"] = *instanceSpec.MaxStorageSize
		resultItem["min_storage_size"] = *instanceSpec.MinStorageSize
		resultItem["machine_type"] = *instanceSpec.MachineType
		resultItem["max_io_bandwidth"] = *instanceSpec.MaxIoBandWidth
		zoneStockInfos := make([]map[string]interface{}, 0)
		for _, zoneStockInfoItem := range instanceSpec.ZoneStockInfos {
			zoneStockInfo := make(map[string]interface{})
			zoneStockInfo["zone"] = *zoneStockInfoItem.Zone
			zoneStockInfo["has_stock"] = *zoneStockInfoItem.HasStock

			zoneStockInfos = append(zoneStockInfos, zoneStockInfo)
		}
		resultItem["zone_stock_infos"] = zoneStockInfos
		result = append(result, resultItem)
	}

	id := "cynosdb_zoneconfig_" + region
	d.SetId(id)
	_ = d.Set("list", result)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {

		if err := tccommon.WriteToFile(output.(string), result); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}

	}
	return nil
}
