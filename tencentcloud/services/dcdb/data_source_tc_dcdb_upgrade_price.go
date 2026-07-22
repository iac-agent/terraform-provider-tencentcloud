package dcdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDcdbUpgradePrice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbUpgradePriceRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID",
			},

			"upgrade_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Upgrade 类型，ADD: add new shard，EXPAND: upgrade the existing shard，SPLIT: split existing shard。",
			},

			"add_shard_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "配置 for adding new shard。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"shard_count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The 数量 new shards。",
						},
						"shard_memory": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Shard memory size （GB）。",
						},
						"shard_storage": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Shard storage capacity （GB）。",
						},
					},
				},
			},

			"expand_shard_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "配置 for expanding existing shard。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"shard_instance_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "列表 shard ID。",
						},
						"shard_memory": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Shard memory size （GB）。",
						},
						"shard_storage": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Shard storage capacity （GB）。",
						},
						"shard_node_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Shard node count。",
						},
					},
				},
			},

			"split_shard_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "配置 for splitting existing shard。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"shard_instance_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "列表 shard ID。",
						},
						"split_rate": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Data split ratio，fixed at 50%。",
						},
						"shard_memory": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Shard memory size （GB）。",
						},
						"shard_storage": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Shard storage capacity （GB）。",
						},
					},
				},
			},

			"amount_unit": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Price unit. 有效值：`pent` (cent)，`microPent` (microcent)。",
			},

			"original_price": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Original price. 单位：Cent (default). If the request parameter 包含`AmountUnit`，see `AmountUnit` 描述 Currency: CNY (Chinese site)，USD (international site)。",
			},

			"price": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The actual price may be different from the original price due to discounts. 单位：Cent (default). If the request parameter 包含`AmountUnit`，see `AmountUnit` 描述 Currency: CNY (Chinese site)，USD (international site)。",
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

func dataSourceTencentCloudDcdbUpgradePriceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dcdb_upgrade_price.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var (
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("upgrade_type"); ok {
		paramMap["UpgradeType"] = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "add_shard_config"); ok {
		addShardConfig := dcdb.AddShardConfig{}
		if v, ok := dMap["shard_count"]; ok {
			addShardConfig.ShardCount = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["shard_memory"]; ok {
			addShardConfig.ShardMemory = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["shard_storage"]; ok {
			addShardConfig.ShardStorage = helper.IntInt64(v.(int))
		}
		paramMap["AddShardConfig"] = &addShardConfig
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "expand_shard_config"); ok {
		expandShardConfig := dcdb.ExpandShardConfig{}
		if v, ok := dMap["shard_instance_ids"]; ok {
			shardInstanceIdsSet := v.(*schema.Set).List()
			expandShardConfig.ShardInstanceIds = helper.InterfacesStringsPoint(shardInstanceIdsSet)
		}
		if v, ok := dMap["shard_memory"]; ok {
			expandShardConfig.ShardMemory = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["shard_storage"]; ok {
			expandShardConfig.ShardStorage = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["shard_node_count"]; ok {
			expandShardConfig.ShardNodeCount = helper.IntInt64(v.(int))
		}
		paramMap["ExpandShardConfig"] = &expandShardConfig
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "split_shard_config"); ok {
		splitShardConfig := dcdb.SplitShardConfig{}
		if v, ok := dMap["shard_instance_ids"]; ok {
			shardInstanceIdsSet := v.(*schema.Set).List()
			splitShardConfig.ShardInstanceIds = helper.InterfacesStringsPoint(shardInstanceIdsSet)
		}
		if v, ok := dMap["split_rate"]; ok {
			splitShardConfig.SplitRate = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["shard_memory"]; ok {
			splitShardConfig.ShardMemory = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["shard_storage"]; ok {
			splitShardConfig.ShardStorage = helper.IntInt64(v.(int))
		}
		paramMap["SplitShardConfig"] = &splitShardConfig
	}

	if v, ok := d.GetOk("amount_unit"); ok {
		paramMap["AmountUnit"] = helper.String(v.(string))
	}

	service := DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var result *dcdb.DescribeDCDBUpgradePriceResponseParams
	var e error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e = service.DescribeDcdbUpgradePriceByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if result != nil {
		if result.OriginalPrice != nil {
			_ = d.Set("original_price", result.OriginalPrice)
		}

		if result.Price != nil {
			_ = d.Set("price", result.Price)
		}

		if result.Formula != nil {
			_ = d.Set("formula", result.Formula)
		}
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
