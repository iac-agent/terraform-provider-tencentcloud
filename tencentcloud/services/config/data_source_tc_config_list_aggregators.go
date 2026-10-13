package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudConfigListAggregators() *schema.Resource {
	return &schema.Resource{
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			return diag.FromErr(dataSourceTencentCloudConfigListAggregatorsRead(d, meta))
		},
		Schema: map[string]*schema.Schema{
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total count of aggregators.",
			},
			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Aggregator (account group) list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account group name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account group description.",
						},
						"owner_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Owner user ID.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"account_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Count of members in the account group.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account group type, RD: global account group.",
						},
						"account_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account group ID.",
						},
						"aggregator_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Account group status.",
						},
						"member_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account group name.",
						},
					},
				},
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudConfigListAggregatorsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_config_list_aggregators.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ConfigService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	var (
		total uint64
		items []*configv20220802.Aggregator
	)

	var e error
	total, items, e = service.DescribeConfigListAggregators(ctx)
	if e != nil {
		return e
	}

	_ = d.Set("total", total)

	itemList := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		itemMap := map[string]interface{}{}
		if item.Name != nil {
			itemMap["name"] = item.Name
		}

		if item.Description != nil {
			itemMap["description"] = item.Description
		}

		if item.OwnerUin != nil {
			itemMap["owner_uin"] = int(*item.OwnerUin)
		}

		if item.CreateTime != nil {
			itemMap["create_time"] = item.CreateTime
		}

		if item.AccountCount != nil {
			itemMap["account_count"] = int(*item.AccountCount)
		}

		if item.Type != nil {
			itemMap["type"] = item.Type
		}

		if item.AccountGroupId != nil {
			itemMap["account_group_id"] = item.AccountGroupId
		}

		if item.AggregatorStatus != nil {
			itemMap["aggregator_status"] = int(*item.AggregatorStatus)
		}

		if item.MemberName != nil {
			itemMap["member_name"] = item.MemberName
		}

		itemList = append(itemList, itemMap)
	}

	_ = d.Set("items", itemList)

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
