package tcmq

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTcmqSubscribe() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcmqSubscribeCreate,
		Read:   resourceTencentCloudTcmqSubscribeRead,
		Update: resourceTencentCloudTcmqSubscribeUpdate,
		Delete: resourceTencentCloudTcmqSubscribeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"topic_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Topic 名称，which must be unique in the same topic under the same 账号 in the same 地域 It can contain up to 64 letters，digits，and hyphens and must begin with a letter。",
			},

			"subscription_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Subscription 名称，which must be unique in the same topic under the same 账号 in the same 地域 It can contain up to 64 letters，digits，and hyphens and must begin with a letter。",
			},

			"protocol": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ubscription 协议 Currently，two protocols are supported: `http` and `queue`. To use the `http` 协议，you need to build your own web server to receive messages. With the `queue` 协议，messages are automatically pushed to a CMQ queue and you can pull them concurrently。",
			},

			"endpoint": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "`Endpoint` for notification receipt，which is distinguished by `协议`. For `http`，`Endpoint` must begin with `http://` and `主机` can be a 域名 名称 or IP. For `Queue`，enter `QueueName`. Note that currently the push service cannot push messages to a VPC; therefore，if a VPC 域名 名称 or 地址 is entered for `Endpoint`，pushed messages will not be received. Currently，messages can be pushed only to the public network and classic network。",
			},

			"notify_strategy": {
				Default:     "EXPONENTIAL_DECAY_RETRY",
				Optional:    true,
				Type:        schema.TypeString,
				Description: "CMQ push server retry policy in case an 错误 occurs while pushing a 消息 to `Endpoint`. 有效值：1. `BACKOFF_RETRY`: backoff retry，which is to retry at a fixed interval，discard the 消息 after a certain 数量 retries，and continue to push the next 消息; 2. `EXPONENTIAL_DECAY_RETRY`: exponential decay retry，which is to retry at an exponentially increasing interval，such as 1s，2s，4s，8s，and so on. As a 消息 can be retained in a topic for one day，failed messages will be discarded at most after one day of retry. 默认值：`EXPONENTIAL_DECAY_RETRY`。",
			},

			"filter_tags": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "消息 body 标签 (用于message filtering). The 数量 标签 cannot exceed 5，and each 标签 can contain up to 16 characters. It is used in conjunction with the `MsgTag` parameter of `(Batch)PublishMessage`. Rules: 1. If `FilterTag` is not configured，no matter whether `MsgTag` is configured，the subscription will receive all messages published to the topic; 2. If the 数组 `FilterTag` values has a 值，only when at least one of the values in the array also exists in the 数组 `MsgTag` values (i.e.，`FilterTag` and `MsgTag` have an intersection) can the subscription receive messages published to the topic; 3. If the 数组 `FilterTag` values has a 值，but `MsgTag` is not configured，then no 消息 published to the topic will be received，which can be considered as a special case of rule 2 as `FilterTag` and `MsgTag` do not intersect in this case. The overall design idea of rules is based on the intention of the subscriber。",
			},

			"binding_key": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The 数量 `BindingKey` cannot exceed 5，and the length of each `BindingKey` cannot exceed 64 bytes. This field 表示filtering policy for subscribing to and receiving messages. Each `BindingKey` includes up to 15 dots (namely up to 16 segments)。",
			},

			"notify_content_format": {
				Default:     "JSON",
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Push 内容 格式 有效值：1. JSON; 2. SIMPLIFIED，i.e.，the raw 格式 If `协议` is `queue`，this 值 must be `SIMPLIFIED`. If `协议` is `http`，both options are acceptable，and the 默认值为 `JSON`。",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签描述列表",
			},
		},
	}
}

func resourceTencentCloudTcmqSubscribeCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcmq_subscribe.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request          = tcmq.NewCreateCmqSubscribeRequest()
		subscriptionName string
		topicName        string
	)
	if v, ok := d.GetOk("topic_name"); ok {
		topicName = v.(string)
		request.TopicName = helper.String(topicName)
	}

	if v, ok := d.GetOk("subscription_name"); ok {
		subscriptionName = v.(string)
		request.SubscriptionName = helper.String(subscriptionName)
	}

	if v, ok := d.GetOk("protocol"); ok {
		request.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("endpoint"); ok {
		request.Endpoint = helper.String(v.(string))
	}

	if v, ok := d.GetOk("notify_strategy"); ok {
		request.NotifyStrategy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter_tags"); ok {
		filterTagSet := v.(*schema.Set).List()
		for i := range filterTagSet {
			filterTag := filterTagSet[i].(string)
			request.FilterTag = append(request.FilterTag, &filterTag)
		}
	}

	if v, ok := d.GetOk("binding_key"); ok {
		bindingKeySet := v.(*schema.Set).List()
		for i := range bindingKeySet {
			bindingKey := bindingKeySet[i].(string)
			request.BindingKey = append(request.BindingKey, &bindingKey)
		}
	}

	if v, ok := d.GetOk("notify_content_format"); ok {
		request.NotifyContentFormat = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().CreateCmqSubscribe(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tcmq subscribe failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(topicName + tccommon.FILED_SP + subscriptionName)

	return resourceTencentCloudTcmqSubscribeRead(d, meta)
}

func resourceTencentCloudTcmqSubscribeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcmq_subscribe.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TcmqService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	topicName := idSplit[0]
	subscriptionName := idSplit[1]

	subscribe, err := service.DescribeTcmqSubscribeById(ctx, topicName, subscriptionName)
	if err != nil {
		return err
	}

	if subscribe == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	_ = d.Set("topic_name", topicName)

	if subscribe.SubscriptionName != nil {
		_ = d.Set("subscription_name", subscribe.SubscriptionName)
	}

	if subscribe.Protocol != nil {
		_ = d.Set("protocol", subscribe.Protocol)
	}

	if subscribe.Endpoint != nil {
		_ = d.Set("endpoint", subscribe.Endpoint)
	}

	if subscribe.NotifyStrategy != nil {
		_ = d.Set("notify_strategy", subscribe.NotifyStrategy)
	}

	if subscribe.FilterTags != nil {
		_ = d.Set("filter_tags", subscribe.FilterTags)
	}

	if subscribe.BindingKey != nil {
		_ = d.Set("binding_key", subscribe.BindingKey)
	}

	if subscribe.NotifyContentFormat != nil {
		_ = d.Set("notify_content_format", subscribe.NotifyContentFormat)
	}

	return nil
}

func resourceTencentCloudTcmqSubscribeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcmq_subscribe.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := tcmq.NewModifyCmqSubscriptionAttributeRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	topicName := idSplit[0]
	subscriptionName := idSplit[1]

	request.TopicName = &topicName
	request.SubscriptionName = &subscriptionName
	if d.HasChange("topic_name") {
		if v, ok := d.GetOk("topic_name"); ok {
			request.TopicName = helper.String(v.(string))
		}
	}

	if d.HasChange("subscription_name") {
		if v, ok := d.GetOk("subscription_name"); ok {
			request.SubscriptionName = helper.String(v.(string))
		}
	}

	if d.HasChange("notify_strategy") {
		if v, ok := d.GetOk("notify_strategy"); ok {
			request.NotifyStrategy = helper.String(v.(string))
		}
	}

	if d.HasChange("binding_key") {
		if v, ok := d.GetOk("binding_key"); ok {
			bindingKeySet := v.(*schema.Set).List()
			for i := range bindingKeySet {
				bindingKey := bindingKeySet[i].(string)
				request.BindingKey = append(request.BindingKey, &bindingKey)
			}
		}
	}

	if d.HasChange("notify_content_format") {
		if v, ok := d.GetOk("notify_content_format"); ok {
			request.NotifyContentFormat = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().ModifyCmqSubscriptionAttribute(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tcmq subscribe failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTcmqSubscribeRead(d, meta)
}

func resourceTencentCloudTcmqSubscribeDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcmq_subscribe.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TcmqService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	topicName := idSplit[0]
	subscriptionName := idSplit[1]

	if err := service.DeleteTcmqSubscribeById(ctx, topicName, subscriptionName); err != nil {
		return err
	}

	return nil
}
