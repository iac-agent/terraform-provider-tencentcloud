package mariadb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMariadbRenewalPrice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbRenewalPriceRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID",
			},
			"period": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Renewal duration，默认值：1 month。",
			},
			"amount_unit": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Price unit. 有效值：`* pent` (cent)，`* microPent` (microcent)。",
			},
			"original_price": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Original price * 单位：Cent (default). If the request parameter 包含`AmountUnit`，see `AmountUnit` 描述 * Currency: CNY (Chinese site)，USD (international site)。",
			},
			"price": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The actual price may be different from the original price due to discounts. * 单位：Cent (default). If the request parameter 包含`AmountUnit`，see `AmountUnit` 描述 * Currency: CNY (Chinese site)，USD (international site)。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudMariadbRenewalPriceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mariadb_renewal_price.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		price      *mariadb.DescribeRenewalPriceResponseParams
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, _ := d.GetOk("period"); v != nil {
		paramMap["Period"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("amount_unit"); ok {
		paramMap["AmountUnit"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbRenewalPriceByFilter(ctx, paramMap)
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

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
