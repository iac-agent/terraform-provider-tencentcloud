package mariadb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
)

func DataSourceTencentCloudMariadbOrders() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbOrdersRead,
		Schema: map[string]*schema.Schema{
			"deal_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "列表 long 顺序 numbers to be queried，which are returned for the APIs for creating，renewing，or scaling instances。",
			},
			"deals": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "顺序 information list。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deal_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "顺序 number。",
						},
						"owner_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "账号",
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 items。",
						},
						"flow_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID associated process，which can be 用于query the process execution 状态",
						},
						"instance_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "The ID created instance，which 为必填项 only for the 顺序 that creates an instance.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"pay_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Payment 模式 有效值：0 (postpaid)，1 (prepaid)。",
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

func dataSourceTencentCloudMariadbOrdersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mariadb_orders.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service  = MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		deals    []*mariadb.Deal
		dealName string
	)

	if v, ok := d.GetOk("deal_name"); ok {
		dealName = v.(string)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbOrdersByFilter(ctx, dealName)
		if e != nil {
			return tccommon.RetryError(e)
		}

		deals = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(deals))

	if deals != nil {
		for _, deal := range deals {
			dealMap := map[string]interface{}{}

			if deal.DealName != nil {
				dealMap["deal_name"] = deal.DealName
			}

			if deal.OwnerUin != nil {
				dealMap["owner_uin"] = deal.OwnerUin
			}

			if deal.Count != nil {
				dealMap["count"] = deal.Count
			}

			if deal.FlowId != nil {
				dealMap["flow_id"] = deal.FlowId
			}

			if deal.InstanceIds != nil {
				dealMap["instance_ids"] = deal.InstanceIds
			}

			if deal.PayMode != nil {
				dealMap["pay_mode"] = deal.PayMode
			}

			tmpList = append(tmpList, dealMap)
		}

		_ = d.Set("deals", tmpList)
	}

	d.SetId(dealName)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
