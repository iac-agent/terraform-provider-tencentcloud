package tcmq

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTcmqSubscribe() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcmqSubscribeRead,
		Schema: map[string]*schema.Schema{
			"topic_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Topic 名称，其中 必须 是 唯一 在 same 主题 under same 账号 在 same 地域 It 可以 contain up 到 64 letters，digits，和 hyphens 和 必须 begin 使用 letter。",
			},

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

			"subscription_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Fuzzy search 通过 SubscriptionName。",
			},

			"subscription_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Set 的 subscription attributes。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subscription_name": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Subscription 名称，其中 必须 是 唯一 在 same 主题 under same 账号 在 same 地域 It 可以 contain up 到 64 letters，digits，和 hyphens 和 必须 begin 使用 letter。",
						},
						"subscription_id": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Subscription ID，其中 将 是 使用 during 监控 数据 pull。",
						},
						"topic_owner": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Subscription 所有者 APPID。",
						},
						"msg_count": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "数量 messages 到 是 delivered 在 subscription。",
						},
						"last_modify_time": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Time 当 subscription attribute 是 last modified. A Unix 时间戳 accurate down 到 millisecond 将 是 返回。",
						},
						"create_time": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Subscription 创建时间. A Unix 时间戳 accurate down 到 millisecond 将 是 返回。",
						},
						"binding_key": {
							Computed: true,
							Type:     schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Filtering 策略 对于 subscribing 到 和 receiving messages。",
						},
						"endpoint": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Endpoint 该 receives notifications，其中 varies 通过 `协议`: 对于 HTTP， 端点 必须 start 使用 `http://`，和 `主机` 可以 是 域名 或 IP; 对于 `queue`，`queueName` should 是 entered。",
						},
						"filter_tags": {
							Computed: true,
							Type:     schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Filtering 策略 selected 当 subscription 是 创建:如果 `filterType` 是 1，`filterTag` 将 是 用于filtering. 如果 `filterType` 是 2，`bindingKey` 将 是 用于filtering。",
						},
						"protocol": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Subscription 协议 Currently，two protocols 是 支持: HTTP 和 queue. To 使用 HTTP 协议，您 need 到 build your own web 服务器 到 receive messages. With queue 协议，messages 是 automatically pushed 到 CMQ queue 和 您 可以 pull them concurrently。",
						},
						"notify_strategy": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "CMQ push 服务器 retry 策略 在 case 错误 occurs while pushing 消息 到 `Endpoint`. 有效值：1. `BACKOFF_RETRY`: backoff retry，其中 是 到 retry 在 fixed 间隔，discard 消息 after certain 数量 retries，和 continue 到 push next 消息; 2. `EXPONENTIAL_DECAY_RETRY`: exponential decay retry，其中 是 到 retry 在 exponentially increasing 间隔，such 作为 1s，2s，4s，8s，和 so 在. As 消息 可以 是 retained 在 主题 对于 一个 day，failed messages 将 是 discarded 在 most after 一个 day 的 retry. 默认值：`EXPONENTIAL_DECAY_RETRY`。",
						},
						"notify_content_format": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Push 内容 格式 有效值：1. `JSON`; 2. `SIMPLIFIED`，i.e.， raw 格式 如果 `协议` 是 `queue`，此 值 必须 是 `SIMPLIFIED`. 如果 `协议` 是 `http`，both options 是 acceptable，和 默认值为 `JSON`。",
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

func dataSourceTencentCloudTcmqSubscribeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tcmq_subscribe.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("topic_name"); ok {
		paramMap["topic_name"] = v.(string)
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["offset"] = v.(int)
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["limit"] = v.(int)
	}

	if v, ok := d.GetOk("subscription_name"); ok {
		paramMap["subscription_name"] = v.(string)
	}

	service := TcmqService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var subscriptionList []*tdmq.CmqSubscription
	subscriptionNames := make([]string, 0)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTcmqSubscribeByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		subscriptionList = result
		return nil
	})
	if err != nil {
		return err
	}
	result := make([]map[string]interface{}, 0)
	for _, subscription := range subscriptionList {
		resultItem := make(map[string]interface{})
		if subscription.SubscriptionName != nil {
			subscriptionNames = append(subscriptionNames, *subscription.SubscriptionName)
			resultItem["subscription_name"] = *subscription.SubscriptionName
		}
		if subscription.SubscriptionId != nil {
			resultItem["subscription_id"] = *subscription.SubscriptionId
		}
		if subscription.TopicOwner != nil {
			resultItem["topic_owner"] = *subscription.TopicOwner
		}
		if subscription.MsgCount != nil {
			resultItem["msg_count"] = *subscription.MsgCount
		}
		if subscription.LastModifyTime != nil {
			resultItem["last_modify_time"] = *subscription.LastModifyTime
		}
		if subscription.CreateTime != nil {
			resultItem["create_time"] = *subscription.CreateTime
		}
		if subscription.Endpoint != nil {
			resultItem["endpoint"] = *subscription.Endpoint
		}
		if subscription.Protocol != nil {
			resultItem["protocol"] = *subscription.Protocol
		}
		if subscription.NotifyStrategy != nil {
			resultItem["notify_strategy"] = *subscription.NotifyStrategy
		}
		if subscription.NotifyContentFormat != nil {
			resultItem["notify_content_format"] = *subscription.NotifyContentFormat
		}
		if subscription.BindingKey != nil {
			bindingKeys := make([]string, 0)
			for _, item := range subscription.BindingKey {
				bindingKeys = append(bindingKeys, *item)
			}
			resultItem["binding_key"] = bindingKeys
		}
		if subscription.FilterTags != nil {
			filterTags := make([]string, 0)
			for _, item := range subscription.FilterTags {
				filterTags = append(filterTags, *item)
			}
			resultItem["filter_tags"] = filterTags
		}
		result = append(result, resultItem)
	}

	d.SetId(helper.DataResourceIdsHash(subscriptionNames))
	_ = d.Set("subscription_list", result)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
