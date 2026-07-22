package trocket

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmqRocketmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTdmqRocketmqGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqRocketmqGroupRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "集群 ID",
			},

			"namespace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Namespace。",
			},

			"filter_topic": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Topic 名称，which can be 用于query all subscription groups under the topic。",
			},

			"filter_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Consumer group query by consumer 组名称 Fuzzy query is supported。",
			},

			"filter_one_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subscription 组名称 After it is specified，the information of only this subscription group will be returned。",
			},

			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "列表 subscription groups。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Consumer 组名称",
						},
						"consumer_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The 数量 online consumers。",
						},
						"tps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Consumption TPS。",
						},
						"total_accumulative": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total 数量 heaped messages。",
						},
						"consumption_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "`0`: Cluster consumption 模式; `1`: Broadcast consumption 模式; `-1`: Unknown。",
						},
						"read_enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否enable consumption。",
						},
						"retry_partition_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The 数量 partitions in a retry topic。",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "创建时间 （毫秒）。",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "修改时间 （毫秒）。",
						},
						"client_protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Client 协议",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "备注 (up to 128 characters)。",
						},
						"consumer_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Consumer 类型 Enumerated values: ACTIVELY or PASSIVELY。",
						},
						"broadcast_enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否enable broadcast consumption。",
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

func dataSourceTencentCloudTdmqRocketmqGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tdmqRocketmq_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["cluster_id"] = v.(string)
	}

	if v, ok := d.GetOk("namespace_id"); ok {
		paramMap["namespace_id"] = v.(string)
	}

	if v, ok := d.GetOk("filter_topic"); ok {
		paramMap["filter_topic"] = v.(string)
	}

	if v, ok := d.GetOk("filter_group"); ok {
		paramMap["filter_group"] = v.(string)
	}

	if v, ok := d.GetOk("filter_one_group"); ok {
		paramMap["filter_one_group"] = v.(string)
	}

	tdmqRocketmqService := TdmqRocketmqService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var groups []*tdmqRocketmq.RocketMQGroup
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := tdmqRocketmqService.DescribeTdmqRocketmqGroupByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		groups = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read TdmqRocketmq groups failed, reason:%+v", logId, err)
		return err
	}

	groupList := []interface{}{}
	ids := make([]string, 0)
	for _, group := range groups {
		groupMap := map[string]interface{}{}
		if group.Name != nil {
			groupMap["name"] = group.Name
		}
		if group.ConsumerNum != nil {
			groupMap["consumer_num"] = group.ConsumerNum
		}
		if group.TPS != nil {
			groupMap["tps"] = group.TPS
		}
		if group.TotalAccumulative != nil {
			groupMap["total_accumulative"] = group.TotalAccumulative
		}
		if group.ConsumptionMode != nil {
			groupMap["consumption_mode"] = group.ConsumptionMode
		}
		if group.ReadEnabled != nil {
			groupMap["read_enable"] = group.ReadEnabled
		}
		if group.RetryPartitionNum != nil {
			groupMap["retry_partition_num"] = group.RetryPartitionNum
		}
		if group.CreateTime != nil {
			groupMap["create_time"] = group.CreateTime
		}
		if group.UpdateTime != nil {
			groupMap["update_time"] = group.UpdateTime
		}
		if group.ClientProtocol != nil {
			groupMap["client_protocol"] = group.ClientProtocol
		}
		if group.Remark != nil {
			groupMap["remark"] = group.Remark
		}
		if group.ConsumerType != nil {
			groupMap["consumer_type"] = group.ConsumerType
		}
		if group.BroadcastEnabled != nil {
			groupMap["broadcast_enable"] = group.BroadcastEnabled
		}

		groupList = append(groupList, groupMap)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("groups", groupList)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), groupList); e != nil {
			return e
		}
	}

	return nil
}
