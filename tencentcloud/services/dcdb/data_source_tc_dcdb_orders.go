package dcdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDcdbOrders() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbOrdersRead,
		Schema: map[string]*schema.Schema{
			"deal_names": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "列表 long 顺序 numbers 到 是 queried，其中 是 返回 对于 APIs 对于 creating，renewing，或 scaling 实例。",
			},

			"deals": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "顺序 信息 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deal_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "顺序 数量。",
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
							Description: "ID associated process，其中 可以 是 用于query process execution 状态",
						},
						"instance_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "ID 创建 实例，其中 为必填项 仅 对于 顺序 该 creates 实例.注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudDcdbOrdersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dcdb_orders.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("deal_names"); ok {
		dealNamesSet := v.(*schema.Set).List()
		paramMap["DealNames"] = helper.InterfacesStringsPoint(dealNamesSet)
	}

	service := DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var deals []*dcdb.Deal

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDcdbOrdersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		deals = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(deals))
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

			ids = append(ids, *deal.DealName)
			tmpList = append(tmpList, dealMap)
		}

		_ = d.Set("deals", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
