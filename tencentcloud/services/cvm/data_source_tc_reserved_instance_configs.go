package cvm

import (
	"context"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvmintl "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudReservedInstanceConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudReservedInstanceConfigsRead,

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "可用 可用区 该 reserved 实例 locates 在。",
			},
			"duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{31536000, 94608000}),
				Description:  "Validity 周期 的 reserved 实例. 有效 值 是 `31536000`(1 year) 和 `94608000`(3 years)。",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "类型 reserved 实例。",
			},
			"offering_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "过滤器 通过 Payment 类型 Such 作为 All Upfront。",
			},
			"product_description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "过滤器 通过 Platform 描述 (该 是，operating 系统) 对于 Reserved 实例 billing. Shaped like: linux。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// computed
			"config_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An 信息 列表 reserved 实例 配置. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Configuration ID purchasable reserved 实例。",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability 可用区 的 purchasable reserved 实例。",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 类型 reserved 实例。",
						},
						"duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Validity 周期 的 reserved 实例。",
						},
						"price": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Purchase 价格 的 reserved 实例。",
						},
						"currency_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Settlement currency 的 reserved 实例，其中 是 standard currency 代码 作为 listed 在 ISO 4217。",
						},
						"platform": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Platform 的 reserved 实例。",
						},
						"offering_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OfferingType 的 reserved 实例。",
						},
						"usage_price": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "UsagePrice 的 reserved 实例。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudReservedInstanceConfigsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_reserved_instance_configs.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cvmService := CvmService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	filter := make(map[string]string)
	if v, ok := d.GetOk("availability_zone"); ok {
		filter["zone"] = v.(string)
	}
	if v, ok := d.GetOk("duration"); ok {
		filter["duration"] = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("instance_type"); ok {
		filter["instance-type"] = v.(string)
	}
	if v, ok := d.GetOk("offering_type"); ok {
		filter["offering-type"] = v.(string)
	}
	if v, ok := d.GetOk("product_description"); ok {
		filter["product-description"] = v.(string)
	}

	var configs []*cvmintl.ReservedInstancesOffering
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		configs, errRet = cvmService.DescribeReservedInstanceConfigs(ctx, filter)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	configList := make([]map[string]interface{}, 0, len(configs))
	ids := make([]string, 0, len(configs))
	for _, config := range configs {
		mapping := map[string]interface{}{
			"config_id":         config.ReservedInstancesOfferingId,
			"availability_zone": config.Zone,
			"instance_type":     config.InstanceType,
			"duration":          config.Duration,
			"price":             config.FixedPrice,
			"currency_code":     config.CurrencyCode,
			"platform":          config.ProductDescription,
			"offering_type":     config.OfferingType,
			"usage_price":       config.UsagePrice,
		}
		configList = append(configList, mapping)
		ids = append(ids, *config.ReservedInstancesOfferingId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("config_list", configList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set config list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), configList); err != nil {
			return err
		}
	}
	return nil
}
