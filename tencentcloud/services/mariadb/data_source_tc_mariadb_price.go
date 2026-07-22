package mariadb

import (
	"context"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMariadbPrice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbPriceRead,
		Schema: map[string]*schema.Schema{
			"zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "AZ ID purchased 实例。",
			},
			"node_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "数量 实例 nodes，其中 可以 是 获取 通过 querying 实例 规格 through `DescribeDBInstanceSpecs` API。",
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
			"buy_count": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateIntegerMin(1),
				Description:  "quantity 您 want 到 purchase 是 queried 通过 默认值 对于 价格 的 purchasing 1 实例。",
			},
			"period": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Purchase 周期 在 months。",
			},
			"paymode": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Billing 类型 有效值：`postpaid` (pay-作为-您-go)，`prepaid` (monthly subscription)。",
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
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudMariadbPriceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mariadb_price.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		price   *mariadb.DescribePriceResponseParams
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone"); ok {
		paramMap["Zone"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("node_count"); v != nil {
		paramMap["NodeCount"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("memory"); v != nil {
		paramMap["Memory"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("storage"); v != nil {
		paramMap["Storage"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("buy_count"); v != nil {
		paramMap["Count"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("period"); v != nil {
		paramMap["Period"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("paymode"); ok {
		paramMap["Paymode"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("amount_unit"); ok {
		paramMap["AmountUnit"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbPriceByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		price = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0)
	if price.OriginalPrice != nil {
		_ = d.Set("original_price", price.OriginalPrice)
		ids = append(ids, strconv.Itoa(int(*price.OriginalPrice)))
	}

	if price.Price != nil {
		_ = d.Set("price", price.Price)
		ids = append(ids, strconv.Itoa(int(*price.Price)))
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
