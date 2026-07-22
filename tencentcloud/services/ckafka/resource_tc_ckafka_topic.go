package ckafka

import (
	"context"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCkafkaTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCkafkaTopicCreate,
		Read:   resourceTencentCloudCkafkaTopicRead,
		Update: resourceTencentCloudCkafkaTopicUpdate,
		Delete: resourceTencentCLoudCkafkaTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Ckafka 实例 ID。",
			},
			"topic_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "名称 CKafka 主题. It 必须 start 使用 letter， rest 可以 contain letters，numbers 和 dashes(-)。",
			},
			"partition_num": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "数量 分区。",
			},
			"replica_num": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "数量 副本。",
			},
			"enable_white_list": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "是否open ip whitelist，`true`: open，`false`: close。",
			},
			"ip_white_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Ip whitelist，配额 限制，必填 当 enableWhileList=true。",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"note": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "subject note. It 必须 start 使用 letter，和 remaining part 可以 contain letters，numbers 和 dashes (-)。",
			},
			"retention": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      60000,
				ValidateFunc: tccommon.ValidateIntegerMin(60000),
				Description:  "消息 可以 是 selected. Retention 时间，单位 是 ms， 当前 最小 值 是 60000ms。",
			},
			"sync_replica_min_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Min 数量 sync replicas，默认为 `1`。",
			},
			"clean_up_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "delete",
				Description: "Clear 日志 策略，日志 clear 模式，默认为 `delete`. `delete`: logs 是 删除 according 到 存储 时间. `compact`: logs 是 compressed according 到 键 `compact，delete`: logs 是 compressed according 到 键 和 将 是 删除 according 到 存储 时间。",
			},
			"unclean_leader_election_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "是否allow unsynchronized replicas 到 是 selected 作为 leader，默认为 `false`，`true: `allowed，`false`: 不 allowed。",
			},
			"max_message_bytes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Max 消息 bytes. min: 1024 Byte(1KB)，max: 8388608 Byte(8MB)。",
			},
			"segment": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateIntegerMin(3600000),
				Description:  "Segment scrolling 时间，在 ms， 当前 最小 是 3600000ms。",
			},
			"message_storage_location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "消息 存储 location。",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "创建时间 的 CKafka 主题。",
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
				Description: "Data 备份 cos 状态 有效值：`0`，`1`. `1`: do 不 open 数据 备份，`0`: open 数据 备份。",
			},
			"segment_bytes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "数量 bytes rolled 通过 分片。",
			},
		},
	}
}

func resourceTencentCloudCkafkaTopicCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_topic.create")()
	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		ckafkcService   = CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request         = ckafka.NewCreateTopicRequest()
		instanceId      string
		topicName       string
		whiteListSwitch bool
		// ipWhiteLists    = d.Get("ip_white_list").([]interface{})
		// ipWhiteList     = make([]*string, 0, len(ipWhiteLists))
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("topic_name"); ok {
		request.TopicName = helper.String(v.(string))
		topicName = v.(string)
	}

	if v, ok := d.GetOkExists("partition_num"); ok {
		request.PartitionNum = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("replica_num"); ok {
		request.ReplicaNum = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("enable_white_list"); ok {
		request.EnableWhiteList = helper.BoolToInt64Ptr(v.(bool))
		whiteListSwitch = v.(bool)
	}

	if v, ok := d.GetOk("ip_white_list"); ok {
		for _, item := range v.([]interface{}) {
			request.IpWhiteList = append(request.IpWhiteList, helper.String(item.(string)))
		}
	}

	if v, ok := d.GetOk("note"); ok {
		request.Note = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("retention"); ok {
		request.RetentionMs = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("sync_replica_min_num"); ok {
		request.MinInsyncReplicas = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("clean_up_policy"); ok {
		request.CleanUpPolicy = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("unclean_leader_election_enable"); ok {
		request.UncleanLeaderElectionEnable = helper.BoolToInt64Ptr(v.(bool))
	}

	if v, ok := d.GetOkExists("max_message_bytes"); ok {
		request.MaxMessageBytes = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("segment"); ok {
		request.SegmentMs = helper.IntInt64(v.(int))
	}

	if whiteListSwitch {
		if len(request.IpWhiteList) == 0 {
			return fmt.Errorf("this Topic %s Create Failed, reason: ip whitelist switch is on, ip whitelist cannot be empty", topicName)
		}
	}

	//Before create topic,Check if kafka exists
	_, has, error := ckafkcService.DescribeInstanceById(ctx, instanceId)
	if error != nil {
		return error
	}

	if !has {
		return fmt.Errorf("ckafka %s does not exist", instanceId)
	}

	errCreate := ckafkcService.CreateCkafkaTopic(ctx, request)
	if errCreate != nil {
		return errCreate
	}

	_, hasExist, err := ckafkcService.DescribeCkafkaTopicByName(ctx, instanceId, topicName)
	if err != nil {
		return err
	}

	if !hasExist {
		return fmt.Errorf("this Topic %s Create Failed", topicName)
	}

	if len(request.IpWhiteList) > 0 && whiteListSwitch {
		err = ckafkcService.AddCkafkaTopicIpWhiteList(ctx, instanceId, topicName, request.IpWhiteList)
		if err != nil {
			return err
		}
	}

	resourceId := instanceId + tccommon.FILED_SP + topicName
	d.SetId(resourceId)
	return resourceTencentCloudCkafkaTopicRead(d, meta)
}

func resourceTencentCloudCkafkaTopicRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_topic.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		ckafkcService = CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := items[0]
	topicName := items[1]

	topicListInfo, hasExist, e := ckafkcService.DescribeCkafkaTopicByName(ctx, instanceId, topicName)
	if e != nil {
		return e
	}

	if !hasExist {
		d.SetId("")
		return nil
	}

	var topicinfo *ckafka.TopicAttributesResponse
	errInfo := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		topicDetail, e := ckafkcService.DescribeCkafkaTopicAttributes(ctx, instanceId, topicName)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if topicDetail == nil {
			d.SetId("")
			return nil
		}

		topicinfo = topicDetail
		return nil
	})

	if errInfo != nil {
		return fmt.Errorf("[API]Describe kafka topic fail,reason:%s", errInfo.Error())
	}

	_ = d.Set("instance_id", instanceId)

	if topicListInfo.TopicName != nil {
		_ = d.Set("topic_name", topicListInfo.TopicName)
	}

	if topicinfo.PartitionNum != nil {
		_ = d.Set("partition_num", topicinfo.PartitionNum)
	}

	if topicListInfo.ReplicaNum != nil {
		_ = d.Set("replica_num", topicListInfo.ReplicaNum)
	}

	if topicinfo.EnableWhiteList != nil {
		if *topicinfo.EnableWhiteList == 1 {
			_ = d.Set("enable_white_list", true)
		} else {
			_ = d.Set("enable_white_list", false)
		}
	}

	if topicinfo.IpWhiteList != nil {
		_ = d.Set("ip_white_list", topicinfo.IpWhiteList)
	}

	if topicinfo.Note != nil {
		_ = d.Set("note", topicinfo.Note)
	}

	if topicinfo.Config != nil {
		if topicinfo.Config.Retention != nil {
			_ = d.Set("retention", topicinfo.Config.Retention)
		}

		if topicinfo.Config.MinInsyncReplicas != nil {
			_ = d.Set("sync_replica_min_num", topicinfo.Config.MinInsyncReplicas)
		}

		if topicinfo.Config.CleanUpPolicy != nil {
			_ = d.Set("clean_up_policy", topicinfo.Config.CleanUpPolicy)
		}

		if topicinfo.Config.UncleanLeaderElectionEnable != nil {
			if *topicinfo.Config.UncleanLeaderElectionEnable == 1 {
				_ = d.Set("unclean_leader_election_enable", true)
			} else {
				_ = d.Set("unclean_leader_election_enable", false)
			}
		}

		if topicinfo.Config.MaxMessageBytes != nil {
			_ = d.Set("max_message_bytes", topicinfo.Config.MaxMessageBytes)
		}

		if topicinfo.Config.SegmentMs != nil {
			_ = d.Set("segment", topicinfo.Config.SegmentMs)
		}

		if topicinfo.Config.SegmentBytes != nil {
			_ = d.Set("segment_bytes", topicinfo.Config.SegmentBytes)
		}
	}

	if topicinfo.CreateTime != nil {
		_ = d.Set("create_time", helper.FormatUnixTime(uint64(*topicinfo.CreateTime)))
	}

	if topicListInfo.ForwardInterval != nil {
		_ = d.Set("forward_interval", topicListInfo.ForwardInterval)
	}

	if topicListInfo.ForwardCosBucket != nil {
		_ = d.Set("forward_cos_bucket", topicListInfo.ForwardCosBucket)
	}

	if topicListInfo.ForwardStatus != nil {
		_ = d.Set("forward_status", topicListInfo.ForwardStatus)
	}

	return nil
}

func resourceTencentCloudCkafkaTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_topic.update")()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		ckafkcService   = CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		whiteListSwitch bool
	)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := items[0]
	topicName := items[1]

	request := ckafka.NewModifyTopicAttributesRequest()
	request.InstanceId = &instanceId
	request.TopicName = &topicName
	if v, ok := d.GetOkExists("replica_num"); ok {
		request.ReplicaNum = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("enable_white_list"); ok {
		request.EnableWhiteList = helper.BoolToInt64Ptr(v.(bool))
		whiteListSwitch = v.(bool)
	}

	if v, ok := d.GetOk("clean_up_policy"); ok {
		request.CleanUpPolicy = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("retention"); ok {
		request.RetentionMs = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("note"); ok {
		request.Note = helper.String(v.(string))
	}

	if v, ok := d.GetOk("segment"); ok {
		request.SegmentMs = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("sync_replica_min_num"); ok {
		request.MinInsyncReplicas = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("unclean_leader_election_enable"); ok {
		request.UncleanLeaderElectionEnable = helper.BoolToInt64Ptr(v.(bool))
	}

	if d.Get("max_message_bytes").(int) != 0 {
		request.MaxMessageBytes = helper.IntInt64(d.Get("max_message_bytes").(int))
	}

	//Update ip white List
	if whiteListSwitch {
		_, newInterface := d.GetChange("ip_white_list")
		newIpWhiteListInterface := newInterface.([]interface{})
		var newIpWhiteList []*string
		for _, value := range newIpWhiteListInterface {
			newIpWhiteList = append(newIpWhiteList, helper.String(value.(string)))
		}

		if len(newIpWhiteList) == 0 {
			return fmt.Errorf("this Topic %s Create Failed, reason: ip whitelist switch is on, ip whitelist cannot be empty", topicName)
		}

		request.IpWhiteList = newIpWhiteList
	} else {
		//IP whiteList Switch not turned on, and the ip whitelist cannot be modified
		if d.HasChange("ip_white_list") {
			return fmt.Errorf("this Topic %s IP whitelist Modification failed, reason: The Ip Whitelist Switch is not turned on", topicName)
		}
	}

	//Update partition num
	oldPartitionNum, newPartitionNum := d.GetChange("partition_num")
	if newPartitionNum.(int) < oldPartitionNum.(int) {
		return fmt.Errorf("this Topic %s partition Modification failed, reason: The partitonNum must more than current partitionNum", topicName)
	} else {
		if newPartitionNum.(int) > oldPartitionNum.(int) {
			err := ckafkcService.AddCkafkaTopicPartition(ctx, instanceId, topicName, int64(newPartitionNum.(int)))
			if err != nil {
				return err
			}
		}
	}

	err := ckafkcService.ModifyCkafkaTopicAttribute(ctx, request)
	if err != nil {
		return err
	}

	_, has, errDes := ckafkcService.DescribeCkafkaTopicByName(ctx, instanceId, topicName)
	if errDes != nil {
		return errDes
	}

	if !has {
		errDes = fmt.Errorf("this Topic %s Update Failed", topicName)
		return errDes
	}

	return resourceTencentCloudCkafkaTopicRead(d, meta)
}

func resourceTencentCLoudCkafkaTopicDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_topic.delete")()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		ckafkcService = CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := items[0]
	topicName := items[1]

	err := ckafkcService.DeleteCkafkaTopic(ctx, instanceId, topicName)
	if err != nil {
		return err
	}

	return nil
}
