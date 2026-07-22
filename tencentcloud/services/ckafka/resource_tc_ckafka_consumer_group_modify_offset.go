package ckafka

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCkafkaConsumerGroupModifyOffset() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCkafkaConsumerGroupModifyOffsetCreate,
		Read:   resourceTencentCloudCkafkaConsumerGroupModifyOffsetRead,
		Delete: resourceTencentCloudCkafkaConsumerGroupModifyOffsetDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Kafka 实例 ID",
			},

			"group": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "kafka 组。",
			},

			"strategy": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeInt,
				Description: "Reset 策略 的 偏移量.\n" +
					"`0`: Move the offset forward or backward shift bar;\n" +
					"`1`: Alignment reference (by-duration,to-datetime,to-earliest,to-latest), which means moving the offset to the location of the specified timestamp;\n" +
					"`2`: Alignment reference (to-offset), which means to move the offset to the specified offset location.",
			},

			"topics": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "表示topics 该 needs 到 是 reset. Leave 它 空 表示 all。",
			},

			"shift": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "此 字段 必须 是 included 当 strategy 是 0. 如果 它 是 greater 比 zero， 偏移量 将 是 moved backward 通过 shift bars，和 如果 它 是 less 比 zero， 偏移量 将 是 traced back 到 数量 shift entries. After correct reset， new 偏移量 should 是 (old_offset + shift). It should 是 noted 该 如果 new 偏移量 是 less 比 分区's earliest，它 将 是 集合 到 earliest，和 如果 latest greater 比 分区 将 是 集合 到 latest。",
			},

			"shift_timestamp": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Unit ms. 当 strategy 是 1，您 必须 include 此 字段，其中-2 表示 到 reset 偏移量 到 beginning,-1 表示 到 reset 到 latest position (equivalent 到 emptying)，和 other 值 represent 指定 时间. You 将 get 偏移量 的 指定 时间 在 主题 和 then reset 它. 如果 there 是 无 消息 在 指定 时间，get last 偏移量",
			},

			"offset": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "偏移量 location 该 needs 到 是 reset. 当 strategy 是 2，此 字段 必须 是 included。",
			},

			"partitions": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "列表 分区 该 needs 到 是 reset 如果 无 Topics 参数 是 指定. Resets 分区 在 corresponding Partition 列表 all topics. 当 Topics 是 指定， 分区 的 corresponding 主题 列表 指定 Partitions 列表 是 reset。",
			},
		},
	}
}

func resourceTencentCloudCkafkaConsumerGroupModifyOffsetCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_consumer_group_modify_offset.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = ckafka.NewModifyGroupOffsetsRequest()
		instanceId string
		group      string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOk("group"); ok {
		group = v.(string)
		request.Group = helper.String(group)
	}

	if v, _ := d.GetOk("strategy"); v != nil {
		request.Strategy = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("topics"); ok {
		topicsSet := v.(*schema.Set).List()
		for i := range topicsSet {
			topics := topicsSet[i].(string)
			request.Topics = append(request.Topics, &topics)
		}
	}

	if v, _ := d.GetOk("shift"); v != nil {
		request.Shift = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("shift_timestamp"); v != nil {
		request.ShiftTimestamp = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("offset"); v != nil {
		request.Offset = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("partitions"); ok {
		partitionsSet := v.(*schema.Set).List()
		for i := range partitionsSet {
			partitions := partitionsSet[i].(int)
			request.Partitions = append(request.Partitions, helper.IntInt64(partitions))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCkafkaClient().ModifyGroupOffsets(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ckafka consumerGroupModifyOffset failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + tccommon.FILED_SP + group)

	return resourceTencentCloudCkafkaConsumerGroupModifyOffsetRead(d, meta)
}

func resourceTencentCloudCkafkaConsumerGroupModifyOffsetRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_consumer_group_modify_offset.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCkafkaConsumerGroupModifyOffsetDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_consumer_group_modify_offset.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
