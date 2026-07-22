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
				Description: "Starting position 的 queue 列表 到 是 返回 在 当前 页面 在 case 的 paginated 返回. 如果 值 是 entered，限制 必须 是 指定. 如果此参数为空，0 将 是 使用 通过 默认值。",
			},

			"limit": {
				Default:     20,
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "数量 queues 到 是 返回 per 页面 在 case 的 paginated 返回. 如果 此 参数 是 不 passed 在，20 将 是 使用 通过 默认值. Maximum 值: 50。",
			},

			"queue_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "过滤器 通过 QueueName。",
			},

			"queue_name_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "过滤器 通过 CMQ queue 名称",
			},

			"is_tag_filter": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "For filtering 通过 标签，此 参数 必须 是 集合 到 `true`。",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器. Currently，您 可以 过滤器 通过 标签 标签 名称 必须 是 prefixed 使用 `标签:`，such 作为 `标签: 所有者`，`标签: 环境`，或 `标签: business`。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "过滤器 参数 名称",
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
				Description: "Queue 列表。",
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
							Description: "限制 的 数量 messages produced per second. 值 对于 consumed messages 是 1.1 times 此 值",
						},
						"bps": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Bandwidth 限制",
						},
						"max_delay_seconds": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Maximum retention 周期 对于 inflight messages。",
						},
						"max_msg_heap_num": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "最大heaped messages. 值 范围 是 1,000,000-10,000,000 during beta 测试 和 可以 是 1,000,000-1,000,000,000 after product 是 officially released. 默认值为 10,000,000 during beta 测试 和 将 是 100,000,000 after product 是 officially released。",
						},
						"polling_wait_seconds": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Long polling wait 时间 对于 消息 reception. 取值范围：0-30 秒. 默认值：0。",
						},
						"msg_retention_seconds": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "max 周期 during 其中 消息 是 retained before 它 是 automatically acknowledged. 取值范围：30-43,200 秒 (30 秒 到 12 hours). 默认值：3600 秒 (1 hour)。",
						},
						"visibility_timeout": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "消息 visibility 超时 周期 取值范围：1-43200 秒 (i.e.，12 hours). 默认值：30。",
						},
						"max_msg_size": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Maximum 消息 长度. 取值范围：1024-65536 bytes (i.e.，1-64 KB). 默认值：65536。",
						},
						"rewind_seconds": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Rewindable 时间 的 messages 在 queue. 取值范围：0-1,296,000s (如果 消息 rewind 是 已启用). 值 `0` 表示that 消息 rewind 是 不 已启用",
						},
						"create_time": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Queue 创建时间. A Unix 时间戳 accurate down 到 millisecond 将 是 返回。",
						},
						"last_modify_time": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Time 当 queue attribute 是 last modified. A Unix 时间戳 accurate down 到 millisecond 将 是 返回。",
						},
						"active_msg_num": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Total 数量 messages 在 `活跃` 状态 (i.e.，unconsumed) 在 queue，其中 是 approximate 值",
						},
						"inactive_msg_num": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Total 数量 messages 在 `Inactive` 状态 (i.e.，being consumed) 在 queue，其中 是 approximate 值",
						},
						"delay_msg_num": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "数量 delayed messages。",
						},
						"rewind_msg_num": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "数量 retained messages 其中 have been 删除 通过 `DelMsg` API 但 是 still within their rewind 时间 范围。",
						},
						"min_msg_time": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Minimum unconsumed 时间 的 消息 （秒）。",
						},
						"transaction": {
							Computed:    true,
							Type:        schema.TypeBool,
							Description: "1: 事务 queue; 0: general queue。",
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
							Description: "Dead letter queue 策略。",
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
										Description: "Dead letter queue 策略。",
									},
									"max_time_to_live": {
										Computed:    true,
										Type:        schema.TypeInt,
										Description: "Maximum 周期 （秒） before unconsumed 消息 expires，其中 为必填项 如果 `Policy` 是 1. 取值范围：300-43200. 此 值 should 是 smaller 比 `MsgRetentionSeconds` (最大 消息 retention 周期)。",
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
							Description: "Transaction 消息 策略。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"first_query_interval": {
										Computed:    true,
										Type:        schema.TypeInt,
										Description: "First lookback 时间。",
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
										Description: "值 的 标签键",
									},
									"tag_value": {
										Computed:    true,
										Type:        schema.TypeString,
										Description: "值 的 标签值",
									},
								},
							},
						},
						"trace": {
							Computed:    true,
							Type:        schema.TypeBool,
							Description: "消息 trace. true: 已启用; false: 不 已启用",
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
							Description: "Cluster 状态 `0`: creating; `1`: normal; `2`: terminating; `3`: 删除; `4`: isolated; `5`: creation failed; `6`: deletion failed。",
						},
						"max_unacked_msg_num": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "最大unacknowledged messages。",
						},
						"max_msg_backlog_size": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Maximum 大小 的 heaped messages 在 bytes。",
						},
						"retention_size_in_mb": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Queue 存储 space 已配置 对于 消息 rewind. 取值范围：1,024-10,240 MB (如果 消息 rewind 是 已启用). 值 `0` 表示that 消息 rewind 是 不 已启用",
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
