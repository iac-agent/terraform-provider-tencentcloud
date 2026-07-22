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
				Description: "Topic 名称，which must be unique in the same topic under the same 账号 in the same 地域 It can contain up to 64 letters，digits，and hyphens and must begin with a letter。",
			},

			"offset": {
				Default:     0,
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Starting position of the 列表 topics to be returned on the current page in case of paginated return. If a 值 is entered，限制 为必填项. 如果此参数为空，0 will be used by default。",
			},

			"limit": {
				Default:     20,
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "数量 topics to be returned per page in case of paginated return. If this parameter is not passed in，20 will be used by default. Maximum 值: 50。",
			},

			"subscription_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Fuzzy search by SubscriptionName。",
			},

			"subscription_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Set of subscription attributes。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subscription_name": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Subscription 名称，which must be unique in the same topic under the same 账号 in the same 地域 It can contain up to 64 letters，digits，and hyphens and must begin with a letter。",
						},
						"subscription_id": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Subscription ID，which will be used during monitoring data pull。",
						},
						"topic_owner": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Subscription 所有者 APPID。",
						},
						"msg_count": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "数量 messages to be delivered in the subscription。",
						},
						"last_modify_time": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Time when the subscription attribute is last modified. A Unix 时间戳 accurate down to the millisecond will be returned。",
						},
						"create_time": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "Subscription 创建时间. A Unix 时间戳 accurate down to the millisecond will be returned。",
						},
						"binding_key": {
							Computed: true,
							Type:     schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Filtering policy for subscribing to and receiving messages。",
						},
						"endpoint": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Endpoint that receives notifications，which varies by `协议`: for HTTP，the endpoint must start with `http://`，and the `主机` can be a 域名 or IP; for `queue`，`queueName` should be entered。",
						},
						"filter_tags": {
							Computed: true,
							Type:     schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Filtering policy selected when a subscription is created:If `filterType` is 1，`filterTag` will be 用于filtering. If `filterType` is 2，`bindingKey` will be 用于filtering。",
						},
						"protocol": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Subscription 协议 Currently，two protocols are supported: HTTP and queue. To use the HTTP 协议，you need to build your own web server to receive messages. With the queue 协议，messages are automatically pushed to a CMQ queue and you can pull them concurrently。",
						},
						"notify_strategy": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "CMQ push server retry policy in case an 错误 occurs while pushing a 消息 to `Endpoint`. 有效值：1. `BACKOFF_RETRY`: backoff retry，which is to retry at a fixed interval，discard the 消息 after a certain 数量 retries，and continue to push the next 消息; 2. `EXPONENTIAL_DECAY_RETRY`: exponential decay retry，which is to retry at an exponentially increasing interval，such as 1s，2s，4s，8s，and so on. As a 消息 can be retained in a topic for one day，failed messages will be discarded at most after one day of retry. 默认值：`EXPONENTIAL_DECAY_RETRY`。",
						},
						"notify_content_format": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Push 内容 格式 有效值：1. `JSON`; 2. `SIMPLIFIED`，i.e.，the raw 格式 If `协议` is `queue`，this 值 must be `SIMPLIFIED`. If `协议` is `http`，both options are acceptable，and the 默认值为 `JSON`。",
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
