package tcmq

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTcmqQueue() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcmqQueueRead,
		Schema: map[string]*schema.Schema{
			"offset": {
				Default:     0,
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Starting position of a queue list to be returned on the current page in case of paginated return. If a 值 is entered，限制 must be specified. 如果此参数为空，0 will be used by default。",
			},

			"limit": {
				Default:     20,
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The 数量 queues to be returned per page in case of paginated return. If this parameter is not passed in，20 will be used by default. Maximum 值: 50。",
			},

			"queue_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter by QueueName。",
			},

			"queue_name_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Filter by CMQ queue 名称",
			},

			"is_tag_filter": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "For filtering by 标签，this parameter must be set to `true`。",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter. Currently，you can filter by 标签 The 标签 名称 must be prefixed with `标签:`，such as `标签: 所有者`，`标签: environment`，or `标签: business`。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filter parameter 名称",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "值",
						},
					},
				},
			},

			"queue_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Queue list。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"queue_id": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "消息 queue ID。",
						},
						"queue_name": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "消息 queue 名称",
						},
						"qps": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "限制 of the 数量 messages produced per second. The 值 for consumed messages is 1.1 times this 值",
						},
						"bps": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Bandwidth 限制",
						},
						"max_delay_seconds": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Maximum retention 周期 for inflight messages。",
						},
						"max_msg_heap_num": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "最大heaped messages. The 值 range is 1,000,000-10,000,000 during the beta test and can be 1,000,000-1,000,000,000 after the product is officially released. The 默认值为 10,000,000 during the beta test and will be 100,000,000 after the product is officially released。",
						},
						"polling_wait_seconds": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Long polling wait time for 消息 reception. 取值范围：0-30 seconds. 默认值：0。",
						},
						"msg_retention_seconds": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "The max 周期 during which a 消息 is retained before it is automatically acknowledged. 取值范围：30-43,200 seconds (30 seconds to 12 hours). 默认值：3600 seconds (1 hour)。",
						},
						"visibility_timeout": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "消息 visibility timeout 周期 取值范围：1-43200 seconds (i.e.，12 hours). 默认值：30。",
						},
						"max_msg_size": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Maximum 消息 length. 取值范围：1024-65536 bytes (i.e.，1-64 KB). 默认值：65536。",
						},
						"rewind_seconds": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Rewindable time of messages in the queue. 取值范围：0-1,296,000s (if 消息 rewind is 已启用). The 值 `0` 表示that 消息 rewind is not 已启用",
						},
						"create_time": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Queue 创建时间. A Unix 时间戳 accurate down to the millisecond will be returned。",
						},
						"last_modify_time": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Time when the queue attribute is last modified. A Unix 时间戳 accurate down to the millisecond will be returned。",
						},
						"active_msg_num": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Total 数量 messages in `活跃` 状态 (i.e.，unconsumed) in the queue，which is an approximate 值",
						},
						"inactive_msg_num": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Total 数量 messages in `Inactive` 状态 (i.e.，being consumed) in the queue，which is an approximate 值",
						},
						"delay_msg_num": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "数量 delayed messages。",
						},
						"rewind_msg_num": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "数量 retained messages which have been deleted by the `DelMsg` API but are still within their rewind time range。",
						},
						"min_msg_time": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Minimum unconsumed time of 消息 （秒）。",
						},
						"transaction": {
							Computed:    true,
							Type:        schema.TypeBool,
							Description: "1: transaction queue; 0: general queue。",
						},
						"dead_letter_source": {
							Computed:    true,
							Type:        schema.TypeList,
							Description: "Dead letter queue。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"queue_id": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "消息 queue ID。",
									},
									"queue_name": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "消息 queue 名称",
									},
								},
							},
						},
						"dead_letter_policy": {
							Computed:    true,
							Type:        schema.TypeList,
							Description: "Dead letter queue policy。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dead_letter_queue": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "Dead letter queue。",
									},
									"policy": {
										Computed:    true,
										Type:        schema.TypeInt,
										Description: "Dead letter queue policy。",
									},
									"max_time_to_live": {
										Computed:    true,
										Type:        schema.TypeInt,
										Description: "Maximum 周期 （秒） before an unconsumed 消息 expires，which 为必填项 if `Policy` is 1. 取值范围：300-43200. This 值 should be smaller than `MsgRetentionSeconds` (maximum 消息 retention 周期)。",
									},
									"max_receive_count": {
										Computed:    true,
										Type:        schema.TypeInt,
										Description: "最大receipts。",
									},
								},
							},
						},
						"transaction_policy": {
							Computed:    true,
							Type:        schema.TypeList,
							Description: "Transaction 消息 policy。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"first_query_interval": {
										Computed:    true,
										Type:        schema.TypeInt,
										Description: "First lookback time。",
									},
									"max_query_count": {
										Computed:    true,
										Type:        schema.TypeInt,
										Description: "最大queries。",
									},
								},
							},
						},
						"create_uin": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "创建者 `Uin`。",
						},
						"tags": {
							Computed:    true,
							Type:        schema.TypeList,
							Description: "Associated 标签",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "值 of the 标签键",
									},
									"tag_value": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "值 of the 标签值",
									},
								},
							},
						},
						"trace": {
							Computed:    true,
							Type:        schema.TypeBool,
							Description: "消息 trace. true: 已启用; false: not 已启用",
						},
						"tenant_id": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Tenant ID。",
						},
						"namespace_name": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Namespace 名称",
						},
						"status": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Cluster 状态 `0`: creating; `1`: normal; `2`: terminating; `3`: deleted; `4`: isolated; `5`: creation failed; `6`: deletion failed。",
						},
						"max_unacked_msg_num": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "The 最大unacknowledged messages。",
						},
						"max_msg_backlog_size": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Maximum size of heaped messages in bytes。",
						},
						"retention_size_in_mb": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Queue storage space configured for 消息 rewind. 取值范围：1,024-10,240 MB (if 消息 rewind is 已启用). The 值 `0` 表示that 消息 rewind is not 已启用",
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

func dataSourceTencentCloudTcmqQueueRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tcmq_queue.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["offset"] = v.(int)
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["limit"] = v.(int)
	}

	if v, ok := d.GetOk("queue_name"); ok {
		paramMap["queue_name"] = v.(string)
	}

	if v, ok := d.GetOk("queue_name_list"); ok {
		queueNameListSet := v.(*schema.Set).List()
		queueNameList := make([]string, 0)
		for i := range queueNameListSet {
			queueName := queueNameListSet[i].(string)
			queueNameList = append(queueNameList, queueName)
		}
		paramMap["queue_name_list"] = queueNameList
	}

	if v, _ := d.GetOk("is_tag_filter"); v != nil {
		paramMap["is_tag_filter"] = v.(bool)
	}

	if v, ok := d.GetOk("filters"); ok {
		filters := make([]map[string]interface{}, 0)
		for _, item := range v.(*schema.Set).List() {
			filter := item.(map[string]interface{})
			name := filter["name"].(string)
			values := make([]string, 0)
			values = append(values, filter["values"].([]string)...)
			filters = append(filters, map[string]interface{}{
				"name":   name,
				"values": values,
			})
		}
		paramMap["fileters"] = filters

	}

	service := TcmqService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var queueList []*tcmq.CmqQueue
	queueNames := make([]string, 0)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTcmqQueueByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		queueList = result
		return nil
	})
	if err != nil {
		return err
	}
	result := make([]map[string]interface{}, 0)
	for _, queue := range queueList {
		queueNames = append(queueNames, *queue.QueueName)
		queueItem := make(map[string]interface{})
		if queue.QueueId != nil {
			queueItem["queue_id"] = *queue.QueueId
		}
		if queue.QueueName != nil {
			queueItem["queue_name"] = *queue.QueueName
		}
		if queue.Qps != nil {
			queueItem["qps"] = *queue.Qps
		}
		if queue.Bps != nil {
			queueItem["bps"] = *queue.Bps
		}
		if queue.MaxDelaySeconds != nil {
			queueItem["max_delay_seconds"] = *queue.MaxDelaySeconds
		}
		if queue.MaxMsgHeapNum != nil {
			queueItem["max_msg_heap_num"] = *queue.MaxMsgHeapNum
		}
		if queue.PollingWaitSeconds != nil {
			queueItem["polling_wait_seconds"] = *queue.PollingWaitSeconds
		}
		if queue.MsgRetentionSeconds != nil {
			queueItem["msg_retention_seconds"] = *queue.MsgRetentionSeconds
		}
		if queue.VisibilityTimeout != nil {
			queueItem["visibility_timeout"] = *queue.VisibilityTimeout
		}
		if queue.MaxMsgSize != nil {
			queueItem["max_msg_size"] = *queue.MaxMsgSize
		}
		if queue.RewindSeconds != nil {
			queueItem["rewind_seconds"] = *queue.RewindSeconds
		}
		if queue.CreateTime != nil {
			queueItem["create_time"] = *queue.CreateTime
		}
		if queue.LastModifyTime != nil {
			queueItem["last_modify_time"] = *queue.LastModifyTime
		}
		if queue.ActiveMsgNum != nil {
			queueItem["active_msg_num"] = *queue.ActiveMsgNum
		}
		if queue.InactiveMsgNum != nil {
			queueItem["inactive_msg_num"] = *queue.InactiveMsgNum
		}
		if queue.DelayMsgNum != nil {
			queueItem["delay_msg_num"] = *queue.DelayMsgNum
		}
		if queue.RewindMsgNum != nil {
			queueItem["rewind_msg_num"] = *queue.RewindMsgNum
		}
		if queue.MinMsgTime != nil {
			queueItem["min_msg_time"] = *queue.MinMsgTime
		}
		if queue.Transaction != nil {
			queueItem["transaction"] = *queue.Transaction
		}
		if queue.CreateUin != nil {
			queueItem["create_uin"] = *queue.CreateUin
		}
		if queue.Trace != nil {
			queueItem["trace"] = *queue.Trace
		}
		if queue.TenantId != nil {
			queueItem["tenant_id"] = *queue.TenantId
		}
		if queue.NamespaceName != nil {
			queueItem["namespace_name"] = *queue.NamespaceName
		}
		if queue.Status != nil {
			queueItem["status"] = *queue.Status
		}
		if queue.MaxUnackedMsgNum != nil {
			queueItem["max_unacked_msg_num"] = *queue.MaxUnackedMsgNum
		}
		if queue.MaxMsgBacklogSize != nil {
			queueItem["max_msg_backlog_size"] = *queue.MaxMsgBacklogSize
		}
		if queue.RetentionSizeInMB != nil {
			queueItem["retention_size_in_mb"] = *queue.RetentionSizeInMB
		}
		if queue.DeadLetterSource != nil {
			deadLetterSource := make([]map[string]interface{}, 0)
			for _, item := range queue.DeadLetterSource {
				deadLetterSource = append(deadLetterSource, map[string]interface{}{
					"queue_id":   item.QueueId,
					"queue_name": item.QueueName,
				})
			}
			queueItem["dead_letter_source"] = deadLetterSource
		}
		if queue.Tags != nil {
			tags := make([]map[string]interface{}, 0)
			for _, item := range queue.Tags {
				tags = append(tags, map[string]interface{}{
					"tag_key":   item.TagKey,
					"tag_value": item.TagValue,
				})
			}
			queueItem["tags"] = tags
		}
		if queue.DeadLetterPolicy != nil {
			deadLetterPolicy := make(map[string]interface{})
			deadLetterPolicy["dead_letter_queue"] = queue.DeadLetterPolicy.DeadLetterQueue
			deadLetterPolicy["policy"] = queue.DeadLetterPolicy.Policy
			deadLetterPolicy["max_time_to_live"] = queue.DeadLetterPolicy.MaxTimeToLive
			deadLetterPolicy["max_receive_count"] = queue.DeadLetterPolicy.MaxReceiveCount
			queueItem["dead_letter_policy"] = deadLetterPolicy
		}
		if queue.TransactionPolicy != nil {
			transactionPolicy := make(map[string]interface{})
			transactionPolicy["first_query_interval"] = queue.TransactionPolicy.FirstQueryInterval
			transactionPolicy["max_query_count"] = queue.TransactionPolicy.MaxQueryCount
			queueItem["transaction_policy"] = transactionPolicy
		}

		result = append(result, queueItem)
	}
	d.SetId(helper.DataResourceIdsHash(queueNames))
	_ = d.Set("queue_list", result)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
