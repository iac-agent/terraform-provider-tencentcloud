package tcmq

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTcmqTopic() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcmqTopicRead,
		Schema: map[string]*schema.Schema{
			"offset": {
				Default:     0,
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Starting position 的 列表 topics 到 是 返回 在 当前 页面 在 case 的 paginated 返回. 如果 值 是 entered，限制 为必填项. 如果此参数为空，0 将 是 使用 通过 默认值。",
			},

			"limit": {
				Default:     20,
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "数量 topics 到 是 返回 per 页面 在 case 的 paginated 返回. 如果 此 参数 是 不 passed 在，20 将 是 使用 通过 默认值. Maximum 值: 50。",
			},

			"topic_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Fuzzy search 通过 TopicName。",
			},

			"topic_name_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "过滤器 通过 CMQ 主题 名称",
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

			"topic_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Topic 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_id": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Topic ID。",
						},
						"topic_name": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Topic 名称",
						},
						"msg_retention_seconds": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Maximum lifecycle 的 消息 在 主题. After 周期 指定 通过 此 参数 has elapsed since 消息 是 sent 到 主题， 消息 将 是 删除 无 matter whether 它 has been successfully pushed 到 用户 此 参数 是 measured （秒） 和 defaulted 到 一个 day (86,400 秒)，其中 不能 是 modified。",
						},
						"max_msg_size": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Maximum 消息 大小，其中 ranges 从 1,024 到 1,048,576 bytes (i.e.，1-1,024 KB). 默认值为 65,536。",
						},
						"qps": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "数量 messages published per second。",
						},
						"filter_type": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Filtering 策略 selected 当 subscription 是 创建: 如果 `filterType` 是 1，`FilterTag` 将 是 用于filtering. 如果 `filterType` 是 2，`BindingKey` 将 是 用于filtering。",
						},
						"create_time": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Topic 创建时间. A Unix 时间戳 accurate down 到 millisecond 将 是 返回。",
						},
						"last_modify_time": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Time 当 主题 attribute 是 last modified. A Unix 时间戳 accurate down 到 millisecond 将 是 返回。",
						},
						"msg_count": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "数量 当前 messages 在 主题 (数量 retained messages)。",
						},
						"create_uin": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "创建者 `Uin`. `资源` 字段 对于 CAM authentication 是 composed 的 此 字段。",
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
						"broker_type": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "有效值：`0` (Pulsar)，`1` (RocketMQ)。",
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

func dataSourceTencentCloudTcmqTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tcmq_topic.read")()
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

	if v, ok := d.GetOk("topic_name"); ok {
		paramMap["topic_name"] = v.(string)
	}

	if v, ok := d.GetOk("topic_name_list"); ok {
		topicNameListSet := v.(*schema.Set).List()
		topicNameList := make([]string, 0)
		for i := range topicNameListSet {
			topicName := topicNameListSet[i].(string)
			topicNameList = append(topicNameList, topicName)
		}
		paramMap["topic_name_list"] = topicNameList

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
		paramMap["filters"] = filters

	}

	service := TcmqService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var topicList []*tdmq.CmqTopic
	topicNames := make([]string, 0)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTcmqTopicByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		topicList = result
		return nil
	})
	if err != nil {
		return err
	}

	result := make([]map[string]interface{}, 0)
	for _, topic := range topicList {
		topicNames = append(topicNames, *topic.TopicName)
		topicItem := make(map[string]interface{})
		if topic.TopicId != nil {
			topicItem["topic_id"] = *topic.TopicId
		}
		if topic.TopicName != nil {
			topicItem["topic_name"] = *topic.TopicName
		}
		if topic.MsgRetentionSeconds != nil {
			topicItem["msg_retention_seconds"] = *topic.MsgRetentionSeconds
		}
		if topic.MaxMsgSize != nil {
			topicItem["max_msg_size"] = *topic.MaxMsgSize
		}
		if topic.Qps != nil {
			topicItem["qps"] = *topic.Qps
		}
		if topic.FilterType != nil {
			topicItem["filter_type"] = *topic.FilterType
		}
		if topic.CreateTime != nil {
			topicItem["create_time"] = *topic.CreateTime
		}
		if topic.LastModifyTime != nil {
			topicItem["last_modify_time"] = *topic.LastModifyTime
		}
		if topic.MsgCount != nil {
			topicItem["msg_count"] = *topic.MsgCount
		}
		if topic.CreateUin != nil {
			topicItem["create_uin"] = *topic.CreateUin
		}
		if topic.Trace != nil {
			topicItem["trace"] = *topic.Trace
		}
		if topic.TenantId != nil {
			topicItem["tenant_id"] = *topic.TenantId
		}
		if topic.NamespaceName != nil {
			topicItem["namespace_name"] = *topic.NamespaceName
		}
		if topic.Status != nil {
			topicItem["status"] = *topic.Status
		}
		if topic.BrokerType != nil {
			topicItem["broker_type"] = *topic.BrokerType
		}

		if topic.Tags != nil {
			tags := make([]map[string]interface{}, 0)
			for _, item := range topic.Tags {
				tags = append(tags, map[string]interface{}{
					"tag_key":   *item.TagKey,
					"tag_value": *item.TagValue,
				})
			}
			topicItem["tags"] = tags
		}

		result = append(result, topicItem)
	}
	d.SetId(helper.DataResourceIdsHash(topicNames))
	_ = d.Set("topic_list", result)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
