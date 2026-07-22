package tcmq

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tcmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTcmqQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcmqQueueCreate,
		Read:   resourceTencentCloudTcmqQueueRead,
		Update: resourceTencentCloudTcmqQueueUpdate,
		Delete: resourceTencentCloudTcmqQueueDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"queue_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Queue 名称，其中 必须 是 唯一 under same 账号 在 same 地域 It 可以 contain up 到 64 letters，digits，和 hyphens 和 必须 begin 使用 letter。",
			},

			"max_msg_heap_num": {
				Optional:    true,
				Default:     10000000,
				Type:        schema.TypeInt,
				Description: "最大heaped messages. 值 范围 是 1,000,000-10,000,000 during beta 测试 和 可以 是 1,000,000-1,000,000,000 after product 是 officially released. 默认值为 10,000,000 during beta 测试 和 将 是 100,000,000 after product 是 officially released。",
			},

			"polling_wait_seconds": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Long polling wait 时间 对于 消息 reception. 取值范围：0-30 秒. 默认值：0。",
			},

			"visibility_timeout": {
				Optional:    true,
				Default:     30,
				Type:        schema.TypeInt,
				Description: "消息 visibility 超时 周期 取值范围：1-43200 秒 (i.e.，12 hours). 默认值：30。",
			},

			"max_msg_size": {
				Optional:    true,
				Default:     65536,
				Type:        schema.TypeInt,
				Description: "Maximum 消息 长度. 取值范围：1024-65536 bytes (i.e.，1-64 KB). 默认值：65536。",
			},

			"msg_retention_seconds": {
				Optional:    true,
				Default:     3600,
				Type:        schema.TypeInt,
				Description: "max 周期 during 其中 消息 是 retained before 它 是 automatically acknowledged. 取值范围：30-43,200 秒 (30 秒 到 12 hours). 默认值：3600 秒 (1 hour)。",
			},

			"rewind_seconds": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Rewindable 时间 的 messages 在 queue. 取值范围：0-1,296,000s (如果 消息 rewind 是 已启用). 值 `0` 表示that 消息 rewind 是 不 已启用",
			},

			"transaction": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "1: 事务 queue; 0: general queue。",
			},

			"first_query_interval": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "First lookback 间隔。",
			},

			"max_query_count": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "最大lookbacks。",
			},

			"dead_letter_queue_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Dead letter queue 名称",
			},

			"policy": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Dead letter 策略. 0: 消息 has been consumed 多个 times 但 不 删除; 1: `Time-To-Live` has elapsed。",
			},

			"max_receive_count": {
				Optional:    true,
				Default:     50,
				Type:        schema.TypeInt,
				Description: "Maximum receipt times. 取值范围：1-1000。",
			},

			"max_time_to_live": {
				Optional:    true,
				Default:     300,
				Type:        schema.TypeInt,
				Description: "Maximum 周期 （秒） before unconsumed 消息 expires，其中 为必填项 如果 `策略` 是 1. 取值范围：300-43200. 此 值 should 是 smaller 比 `msgRetentionSeconds` (最大 消息 retention 周期)。",
			},

			"trace": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "是否enable 消息 trace. true: yes; false: 无. 如果 此 字段 是 不 已配置， 功能 将 不 是 已启用",
			},

			"retention_size_in_mb": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Queue 存储 space 已配置 对于 消息 rewind. 取值范围：10,240-512,000 MB (如果 消息 rewind 是 已启用). 值 `0` 表示that 消息 rewind 是 不 已启用",
			},
		},
	}
}

func resourceTencentCloudTcmqQueueCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcmq_queue.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request   = tcmq.NewCreateCmqQueueRequest()
		queueName string
	)
	if v, ok := d.GetOk("queue_name"); ok {
		queueName = v.(string)
		request.QueueName = helper.String(queueName)
	}

	if v, _ := d.GetOk("max_msg_heap_num"); v != nil {
		request.MaxMsgHeapNum = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("polling_wait_seconds"); v != nil {
		request.PollingWaitSeconds = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("visibility_timeout"); v != nil {
		request.VisibilityTimeout = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("max_msg_size"); v != nil {
		request.MaxMsgSize = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("msg_retention_seconds"); v != nil {
		request.MsgRetentionSeconds = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("rewind_seconds"); v != nil {
		request.RewindSeconds = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("transaction"); v != nil {
		request.Transaction = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("first_query_interval"); v != nil {
		request.FirstQueryInterval = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("max_query_count"); v != nil {
		request.MaxQueryCount = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("dead_letter_queue_name"); ok {
		request.DeadLetterQueueName = helper.String(v.(string))
	}

	if v, _ := d.GetOk("policy"); v != nil {
		request.Policy = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("max_receive_count"); v != nil {
		request.MaxReceiveCount = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("max_time_to_live"); v != nil {
		request.MaxTimeToLive = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("trace"); v != nil {
		request.Trace = helper.Bool(v.(bool))
	}

	if v, _ := d.GetOk("retention_size_in_mb"); v != nil {
		request.RetentionSizeInMB = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().CreateCmqQueue(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tcmq queue failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(queueName)

	return resourceTencentCloudTcmqQueueRead(d, meta)
}

func resourceTencentCloudTcmqQueueRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcmq_queue.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TcmqService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	queueName := d.Id()

	queue, err := service.DescribeTcmqQueueById(ctx, queueName)
	if err != nil {
		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "ResourceNotFound" {
				return nil
			}
		}
		return err
	}

	if queue == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	if queue.QueueName != nil {
		_ = d.Set("queue_name", queue.QueueName)
	}

	if queue.MaxMsgHeapNum != nil {
		_ = d.Set("max_msg_heap_num", queue.MaxMsgHeapNum)
	}

	if queue.PollingWaitSeconds != nil {
		_ = d.Set("polling_wait_seconds", queue.PollingWaitSeconds)
	}

	if queue.VisibilityTimeout != nil {
		_ = d.Set("visibility_timeout", queue.VisibilityTimeout)
	}

	if queue.MaxMsgSize != nil {
		_ = d.Set("max_msg_size", queue.MaxMsgSize)
	}

	if queue.MsgRetentionSeconds != nil {
		_ = d.Set("msg_retention_seconds", queue.MsgRetentionSeconds)
	}

	if queue.RewindSeconds != nil {
		_ = d.Set("rewind_seconds", queue.RewindSeconds)
	}

	if queue.Transaction != nil {
		_ = d.Set("transaction", queue.Transaction)
	}

	if queue.TransactionPolicy != nil {
		if queue.TransactionPolicy.FirstQueryInterval != nil {
			_ = d.Set("first_query_interval", queue.TransactionPolicy.FirstQueryInterval)
		}

		if queue.TransactionPolicy.MaxQueryCount != nil {
			_ = d.Set("max_query_count", queue.TransactionPolicy.MaxQueryCount)
		}
	}

	if len(queue.DeadLetterSource) > 0 && queue.DeadLetterSource[0].QueueName != nil {
		_ = d.Set("dead_letter_queue_name", queue.DeadLetterSource[0].QueueName)
	}

	if queue.DeadLetterPolicy != nil {
		if queue.DeadLetterPolicy.Policy != nil {
			_ = d.Set("policy", queue.DeadLetterPolicy.Policy)
		}

		if queue.DeadLetterPolicy.MaxReceiveCount != nil {
			_ = d.Set("max_receive_count", queue.DeadLetterPolicy.MaxReceiveCount)
		}

		if queue.DeadLetterPolicy.MaxTimeToLive != nil {
			_ = d.Set("max_time_to_live", queue.DeadLetterPolicy.MaxTimeToLive)
		}
	}

	if queue.Trace != nil {
		_ = d.Set("trace", queue.Trace)
	}

	if queue.RetentionSizeInMB != nil {
		_ = d.Set("retention_size_in_mb", queue.RetentionSizeInMB)
	}

	return nil
}

func resourceTencentCloudTcmqQueueUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcmq_queue.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := tcmq.NewModifyCmqQueueAttributeRequest()

	queueName := d.Id()

	request.QueueName = &queueName
	if d.HasChange("queue_name") {
		if v, ok := d.GetOk("queue_name"); ok {
			request.QueueName = helper.String(v.(string))
		}
	}

	if d.HasChange("max_msg_heap_num") {
		if v, _ := d.GetOk("max_msg_heap_num"); v != nil {
			request.MaxMsgHeapNum = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("polling_wait_seconds") {
		if v, _ := d.GetOk("polling_wait_seconds"); v != nil {
			request.PollingWaitSeconds = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("visibility_timeout") {
		if v, _ := d.GetOk("visibility_timeout"); v != nil {
			request.VisibilityTimeout = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("max_msg_size") {
		if v, _ := d.GetOk("max_msg_size"); v != nil {
			request.MaxMsgSize = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("msg_retention_seconds") {
		if v, _ := d.GetOk("msg_retention_seconds"); v != nil {
			request.MsgRetentionSeconds = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("rewind_seconds") {
		if v, _ := d.GetOk("rewind_seconds"); v != nil {
			request.RewindSeconds = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("transaction") {
		if v, _ := d.GetOk("transaction"); v != nil {
			request.Transaction = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("first_query_interval") {
		if v, _ := d.GetOk("first_query_interval"); v != nil {
			request.FirstQueryInterval = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("max_query_count") {
		if v, _ := d.GetOk("max_query_count"); v != nil {
			request.MaxQueryCount = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("dead_letter_queue_name") {
		if v, ok := d.GetOk("dead_letter_queue_name"); ok {
			request.DeadLetterQueueName = helper.String(v.(string))
		}
	}

	if d.HasChange("policy") {
		if v, _ := d.GetOk("policy"); v != nil {
			request.Policy = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("max_receive_count") {
		if v, _ := d.GetOk("max_receive_count"); v != nil {
			request.MaxReceiveCount = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("max_time_to_live") {
		if v, _ := d.GetOk("max_time_to_live"); v != nil {
			request.MaxTimeToLive = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("trace") {
		if v, _ := d.GetOk("trace"); v != nil {
			request.Trace = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("retention_size_in_mb") {
		if v, _ := d.GetOk("retention_size_in_mb"); v != nil {
			request.RetentionSizeInMB = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().ModifyCmqQueueAttribute(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tcmq queue failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTcmqQueueRead(d, meta)
}

func resourceTencentCloudTcmqQueueDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcmq_queue.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TcmqService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	queueName := d.Id()

	if err := service.DeleteTcmqQueueById(ctx, queueName); err != nil {
		return err
	}

	return nil
}
