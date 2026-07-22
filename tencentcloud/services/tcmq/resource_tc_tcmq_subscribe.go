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
				Description: "Topic 名称，其中 必须 是 唯一 在 same 主题 under same 账号 在 same 地域 It 可以 contain up 到 64 letters，digits，和 hyphens 和 必须 begin 使用 letter。",
			},

			"subscription_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Subscription 名称，其中 必须 是 唯一 在 same 主题 under same 账号 在 same 地域 It 可以 contain up 到 64 letters，digits，和 hyphens 和 必须 begin 使用 letter。",
			},

			"protocol": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ubscription 协议 Currently，two protocols 是 支持: `http` 和 `queue`. To 使用 `http` 协议，您 need 到 build your own web 服务器 到 receive messages. With `queue` 协议，messages 是 automatically pushed 到 CMQ queue 和 您 可以 pull them concurrently。",
			},

			"endpoint": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "`Endpoint` 对于 通知 receipt，其中 是 distinguished 通过 `协议`. For `http`，`Endpoint` 必须 begin 使用 `http://` 和 `主机` 可以 是 域名 名称 或 IP. For `Queue`，enter `QueueName`. 注意 该 currently push 服务 不能 push messages 到 VPC; therefore，如果 VPC 域名 名称 或 地址 是 entered 对于 `Endpoint`，pushed messages 将 不 是 received. Currently，messages 可以 是 pushed 仅 到 公有 网络 和 classic 网络。",
			},

			"notify_strategy": {
				Default:     "EXPONENTIAL_DECAY_RETRY",
				Optional:    true,
				Type:        schema.TypeString,
				Description: "CMQ push 服务器 retry 策略 在 case 错误 occurs while pushing 消息 到 `Endpoint`. 有效值：1. `BACKOFF_RETRY`: backoff retry，其中 是 到 retry 在 fixed 间隔，discard 消息 after certain 数量 retries，和 continue 到 push next 消息; 2. `EXPONENTIAL_DECAY_RETRY`: exponential decay retry，其中 是 到 retry 在 exponentially increasing 间隔，such 作为 1s，2s，4s，8s，和 so 在. As 消息 可以 是 retained 在 主题 对于 一个 day，failed messages 将 是 discarded 在 most after 一个 day 的 retry. 默认值：`EXPONENTIAL_DECAY_RETRY`。",
			},

			"filter_tags": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "消息 正文 标签 (用于message filtering). 数量 标签 不能 exceed 5，和 each 标签 可以 contain up 到 16 字符. It 是 使用 在 conjunction 使用 `MsgTag` 参数 的 `(Batch)PublishMessage`. Rules: 1. 如果 `FilterTag` 是 不 已配置，无 matter whether `MsgTag` 是 已配置， subscription 将 receive all messages published 到 主题; 2. 如果 数组 `FilterTag` 值 has 值，仅 当 在 least 一个 的 值 在 数组 also exists 在 数组 `MsgTag` 值 (i.e.，`FilterTag` 和 `MsgTag` have intersection) 可以 subscription receive messages published 到 主题; 3. 如果 数组 `FilterTag` 值 has 值，但 `MsgTag` 是 不 已配置，then 无 消息 published 到 主题 将 是 received，其中 可以 是 considered 作为 special case 的 规则 2 作为 `FilterTag` 和 `MsgTag` do 不 intersect 在 此 case. overall design idea 的 规则 是 based 在 intention 的 subscriber。",
			},

			"binding_key": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "数量 `BindingKey` 不能 exceed 5，和 长度 的 each `BindingKey` 不能 exceed 64 bytes. 此 字段 表示filtering 策略 对于 subscribing 到 和 receiving messages. Each `BindingKey` includes up 到 15 dots (namely up 到 16 segments)。",
			},

			"notify_content_format": {
				Default:     "JSON",
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Push 内容 格式 有效值：1. JSON; 2. SIMPLIFIED，i.e.， raw 格式 如果 `协议` 是 `queue`，此 值 必须 是 `SIMPLIFIED`. 如果 `协议` 是 `http`，both options 是 acceptable，和 默认值为 `JSON`。",
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
