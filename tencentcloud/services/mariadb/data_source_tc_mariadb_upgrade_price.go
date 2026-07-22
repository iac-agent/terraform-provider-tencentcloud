package mariadb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMariadbUpgradePrice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbUpgradePriceRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID",
			},
			"memory": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Memory 大小 （GB）， 其中 可以 是 获取 通过 querying 实例 规格 through `DescribeDBInstanceSpecs` API。",
			},
			"storage": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Storage 容量 （GB）。 最大 和 最小 存储 space 可以 是 获取 通过 querying 实例 规格 through `DescribeDBInstanceSpecs` API。",
			},
			"node_count": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "New 实例 nodes，zero 表示 不 change。",
			},
			"amount_unit": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Price 单位. 有效值：`* pent` (cent)，`* microPent` (microcent)。",
			},
			"original_price": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Original 价格 * 单位：Cent (默认值). 如果 请求 参数 包含`AmountUnit`，see `AmountUnit` 描述 * Currency: CNY (Chinese site)，USD (international site)。",
			},
			"price": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "actual 价格 可能 是 different 从 original 价格 due 到 discounts. * 单位：Cent (默认值). 如果 请求 参数 包含`AmountUnit`，see `AmountUnit` 描述 * Currency: CNY (Chinese site)，USD (international site)。",
			},
			"formula": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Price calculation formula。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudMariadbUpgradePriceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mariadb_upgrade_price.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		price      *mariadb.DescribeUpgradePriceResponseParams
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, _ := d.GetOk("memory"); v != nil {
		paramMap["Memory"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("storage"); v != nil {
		paramMap["Storage"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("node_count"); v != nil {
		paramMap["NodeCount"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("amount_unit"); ok {
		paramMap["AmountUnit"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbUpgradePriceByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		price = result
		return nil
	})

	if err != nil {
		return err
	}

	if price.OriginalPrice != nil {
		_ = d.Set("original_price", price.OriginalPrice)
	}

	if price.Price != nil {
		_ = d.Set("price", price.Price)
	}

	if price.Formula != nil {
		_ = d.Set("formula", price.Formula)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
