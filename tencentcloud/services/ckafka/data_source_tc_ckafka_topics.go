package ckafka

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCkafkaTopics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaTopicsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Ckafka 实例 ID。",
			},
			"topic_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 64),
				Description:  "名称 CKafka 主题. It 必须 start 使用 letter， rest 可以 contain letters，numbers 和 dashes(-). 长度 范围 是 从 1 到 64。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于存储结果。",
			},
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 实例. Each element 包含following attributes。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID CKafka 主题。",
						},
						"topic_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 CKafka 主题。",
						},
						"partition_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 分区。",
						},
						"replica_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 副本。",
						},
						"note": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CKafka 主题 note 描述",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 CKafka 主题。",
						},
						"enable_white_list": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否open IP Whitelist. `true`: open，`false`: close。",
						},
						"ip_white_list_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "IP Whitelist count。",
						},
						"forward_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Periodic 频率 的 数据 备份 到 cos。",
						},
						"forward_cos_bucket": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data 备份 COS 存储桶: 存储桶 地址 该 是 dumped 到 cos。",
						},
						"forward_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Data 备份 cos 状态 `1`: do 不 open 数据 备份，`0`: open 数据 备份。",
						},
						"retention": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "消息 可以 是 selected. Retention 时间(单位 ms)。",
						},
						"sync_replica_min_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Min 数量 sync replicas。",
						},
						"clean_up_policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Clear 日志 策略，日志 clear 模式 `delete`: logs 是 删除 according 到 存储 时间，`compact`: logs 是 compressed according 到 键，`compact，delete`: logs 是 compressed according 到 键 和 将 是 删除 according 到 存储 时间。",
						},
						"unclean_leader_election_enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否allow unsynchronized replicas 到 是 selected 作为 leader，默认为 `false`，`true: `allowed，`false`: 不 allowed。",
						},
						"max_message_bytes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Max 消息 bytes。",
						},
						"segment": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Segment scrolling 时间，在 ms。",
						},
						"segment_bytes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 bytes rolled 通过 分片。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCkafkaTopicsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ckafka_topics.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	instanceId := d.Get("instance_id").(string)
	topicName := d.Get("topic_name").(string)
	ckafkcService := CkafkaService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	topicDetails, err := ckafkcService.DescribeCkafkaTopics(ctx, instanceId, topicName)
	if err != nil {
		return err
	}

	instanceList := make([]map[string]interface{}, 0, len(topicDetails))
	ids := make([]string, 0, len(topicDetails))

	for _, topic := range topicDetails {
		var uncleanLeaderElectionEnable bool
		if topic.Config.UncleanLeaderElectionEnable != nil {
			uncleanLeaderElectionEnable = *topic.Config.UncleanLeaderElectionEnable != 0
		}
		instance := map[string]interface{}{
			"topic_name":                     topic.TopicName,
			"topic_id":                       topic.TopicId,
			"partition_num":                  topic.PartitionNum,
			"replica_num":                    topic.ReplicaNum,
			"note":                           topic.Note,
			"create_time":                    helper.FormatUnixTime(uint64(*topic.CreateTime)),
			"enable_white_list":              topic.EnableWhiteList,
			"ip_white_list_count":            topic.IpWhiteListCount,
			"forward_interval":               topic.ForwardInterval,
			"forward_cos_bucket":             topic.ForwardCosBucket,
			"forward_status":                 topic.ForwardStatus,
			"retention":                      topic.Config.Retention,
			"sync_replica_min_num":           topic.Config.MinInsyncReplicas,
			"clean_up_policy":                topic.Config.CleanUpPolicy,
			"unclean_leader_election_enable": uncleanLeaderElectionEnable,
			"max_message_bytes":              topic.Config.MaxMessageBytes,
			"segment":                        topic.Config.SegmentMs,
			"segment_bytes":                  topic.Config.SegmentBytes,
		}
		resourceId := instanceId + tccommon.FILED_SP + *topic.TopicName
		instanceList = append(instanceList, instance)
		ids = append(ids, resourceId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if err = d.Set("instance_list", instanceList); err != nil {
		log.Printf("[CRITAL]%s provider set ckafka topic list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), instanceList); err != nil {
			return err
		}
	}

	return nil
}
