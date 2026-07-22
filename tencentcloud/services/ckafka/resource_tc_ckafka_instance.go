package ckafka

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCkafkaInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCkafkaInstanceCreate,
		Read:   resourceTencentCloudCkafkaInstanceRead,
		Update: resourceTencentCloudCkafkaInstanceUpdate,
		Delete: resourceTencentCLoudCkafkaInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "实例名称",
			},
			"zone_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Available 可用区 ID",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					multiZone := d.Get("multi_zone_flag").(bool)
					zoneId := d.Get("zone_id").(int)
					v, ok := d.GetOk("zone_ids")

					if !multiZone || !ok || old == "" {
						return old == new
					}

					zoneIds := v.(*schema.Set)
					return zoneIds.Contains(zoneId)
				},
			},
			"specifications_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "profession",
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"standard", "profession", "premium"}),
				Description:  "Specifications 类型 实例. Allowed 值 是 `profession`，`premium`. 默认为 `profession`。",
			},
			"charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      CKAFKA_CHARGE_TYPE_PREPAID,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{CKAFKA_CHARGE_TYPE_POSTPAID, CKAFKA_CHARGE_TYPE_PREPAID}),
				Description:  "charge 类型 实例. 有效 值 是 `PREPAID` 和 `POSTPAID_BY_HOUR`. 默认值为 `PREPAID`。",
			},
			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Prepaid purchase 时间，such 作为 1，是 一个 month。",
			},
			"instance_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
				Description:  "描述 实例类型 `profession`: 1，`standard`: 1(general)，2(standard)，3(advanced)，4(容量)，5(specialized-1)，6(specialized-2)，7(specialized-3)，8(specialized-4)，9(exclusive)。",
			},
			"upgrade_strategy": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				Description: "POSTPAID_BY_HOUR scale-down 模式\n" +
					"- 1: stable transformation;\n" +
					"- 2: High-speed transformer.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "私有网络 ID，它 将 是 basic 网络 如果 不 集合。",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "子网 ID，它 将 是 basic 网络 如果 不 集合。",
			},
			"msg_retention_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "最大 retention 时间 的 实例 logs, 在 minutes." +
					" the default is 10080 (7 days), the maximum is 30 days, and the default 0 is not filled," +
					" which means that the log retention time recovery policy is not enabled.",
			},
			"renew_flag": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "Prepaid automatic renewal mark, 0 表示 默认值 state, initial state," +
					" 1 means automatic renewal, 2 means clear no automatic renewal (user setting).",
			},
			"kafka_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Kafka 版本 (0.10.2/1.1.1/2.4.1)。",
			},
			"band_width": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "实例 带宽 在 MBps。",
			},
			"disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "Disk Size. Its 间隔 varies 使用 带宽, 和 input 必须 是 within 间隔, 其中 可以 是 viewed through control." +
					"If it is not within the interval, the plan will cause a change when first created.",
			},
			"partition": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "Partition Size. Its 间隔 varies 使用 带宽, 和 input 必须 是 within 间隔, 其中 可以 是 viewed through control." +
					"If it is not within the interval, the plan will cause a change when first created.",
			},
			"multi_zone_flag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "表示是否instance 是 multi zones. NOTE: 如果 集合 到 `true`，`zone_ids` 必须 集合 together。",
			},
			"zone_ids": {
				Type:         schema.TypeSet,
				Optional:     true,
				Description:  "列表 可用 可用区 ID NOTE: 此 argument 必须 集合 together 使用 `multi_zone_flag`。",
				RequiredWith: []string{"multi_zone_flag"},
				Elem:         &schema.Schema{Type: schema.TypeInt},
			},
			"tags": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				Deprecated:    "It has been deprecated from version 1.78.5, because it do not support change. Use `tag_set` instead.",
				ConflictsWith: []string{"tag_set"},
				Description:   "标签 的 实例. Partition 大小， professional 版本 does 不 need 标签",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "标签键",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "标签值",
						},
					},
				},
			},
			"tag_set": {
				Type:          schema.TypeMap,
				Optional:      true,
				Computed:      true,
				Description:   "标签 集合 的 实例。",
				ConflictsWith: []string{"tags"},
			},
			"disk_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "类型 磁盘。",
			},
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_create_topic_enable": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Automatic creation. true: 已启用，false: 不 已启用",
						},
						"default_num_partitions": {
							Type:     schema.TypeInt,
							Required: true,
							Description: "如果 auto.create.主题.启用 是 集合 到 true 和 此 值 是 不 集合," +
								"3 will be used by default.",
						},
						"default_replication_factor": {
							Type:     schema.TypeInt,
							Required: true,
							Description: "如果 auto.create.主题.启用 是 集合 到 true 但 此 值 是 不 集合," +
								"2 will be used by default.",
						},
					},
				},
				Description: "实例 配置。",
			},
			"dynamic_retention_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							Description: "Whether 动态 消息 retention 时间 配置 是" +
								"enabled. 0: disabled; 1: enabled.",
						},
						"disk_quota_percentage": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							Description: "Disk 配额 阈值 (在 percentage) 对于 triggering" +
								"the message retention time change event.",
						},
						"step_forward_percentage": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							Description: "Percentage 通过 其中 消息 retention" +
								"time is shortened each time.",
						},
						"bottom_retention": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Minimum retention 时间，在 minutes。",
						},
					},
				},
				Description: "Dynamic 消息 retention 策略 配置。",
			},
			"rebalance_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Deprecated:  "It has been deprecated from version 1.82.37.",
				Description: "Modification 的 rebalancing 时间 after upgrade。",
			},
			"public_network": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerMin(3),
				Description:  "Bandwidth 的 公有 网络。",
			},
			"max_message_byte": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1024, 12*1024*1024),
				Description:  "大小 的 单个 消息 在 bytes 在 实例 级别 取值范围：`1024 - 12*1024*1024 bytes (i.e.，1KB-12MB)。",
			},
			"elastic_bandwidth_switch": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Elastic 带宽 switch 0 不 turned 在 1 turned 在 (0 默认值). 此 takes effect 仅 当 实例 是 创建。",
			},
			"custom_ssl_cert_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom 证书 ID，仅 effective 当 `specifications_type` 是 集合 到 `profession`，支持 自定义 证书 capabilities。",
			},
			"vip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Vip 的 实例。",
			},
			"vport": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "类型 实例。",
			},
		},
	}
}

func ckafkaRequestSetParams(request interface{}, d *schema.ResourceData) {
	values := reflect.ValueOf(request).Elem()

	instanceName := d.Get("instance_name").(string)
	zoneId := d.Get("zone_id").(int)
	values.FieldByName("InstanceName").Set(reflect.ValueOf(helper.String(instanceName)))
	values.FieldByName("ZoneId").Set(reflect.ValueOf(helper.IntInt64(zoneId)))

	requestType := reflect.TypeOf(request)
	if strings.Contains(requestType.String(), "CreateInstancePreRequest") {
		if v, ok := d.GetOk("period"); ok {
			period := int64(v.(int))
			values.FieldByName("Period").Set(reflect.ValueOf(helper.String(fmt.Sprintf("%dm", period))))
		}
		if v, ok := d.GetOk("renew_flag"); ok {
			values.FieldByName("RenewFlag").Set(reflect.ValueOf(helper.Int64(int64(v.(int)))))
		}
	}

	instanceType := helper.IntInt64(1)
	if v, ok := d.GetOkExists("instance_type"); ok {
		instanceType = helper.IntInt64(v.(int))
	}
	values.FieldByName("InstanceType").Set(reflect.ValueOf(instanceType))

	if v, ok := d.GetOk("specifications_type"); ok {
		values.FieldByName("SpecificationsType").Set(reflect.ValueOf(helper.String(v.(string))))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		values.FieldByName("VpcId").Set(reflect.ValueOf(helper.String(v.(string))))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		values.FieldByName("SubnetId").Set(reflect.ValueOf(helper.String(v.(string))))
	}

	if v, ok := d.GetOk("kafka_version"); ok {
		values.FieldByName("KafkaVersion").Set(reflect.ValueOf(helper.String(v.(string))))
	}

	if v, ok := d.GetOk("disk_size"); ok {
		values.FieldByName("DiskSize").Set(reflect.ValueOf(helper.Int64(int64(v.(int)))))

	}

	if v, ok := d.GetOk("band_width"); ok {
		values.FieldByName("BandWidth").Set(reflect.ValueOf(helper.Int64(int64(v.(int)))))
	}

	if v, ok := d.GetOk("partition"); ok {
		values.FieldByName("Partition").Set(reflect.ValueOf(helper.Int64(int64(v.(int)))))

	}

	if v, ok := d.GetOk("tags"); ok {
		tagSet := make([]*ckafka.Tag, 0, 10)
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			tagInfo := ckafka.Tag{
				TagKey:   helper.String(m["key"].(string)),
				TagValue: helper.String(m["value"].(string)),
			}
			tagSet = append(tagSet, &tagInfo)
		}
		values.FieldByName("Tags").Set(reflect.ValueOf(tagSet))
	}

	if v, ok := d.GetOk("disk_type"); ok {
		values.FieldByName("DiskType").Set(reflect.ValueOf(helper.String(v.(string))))
	}

	if flag := d.Get("multi_zone_flag").(bool); flag {
		values.FieldByName("MultiZoneFlag").Set(reflect.ValueOf(helper.Bool(flag)))

		ids := d.Get("zone_ids").(*schema.Set).List()
		zoneIds := make([]*int64, 0)
		for _, v := range ids {
			zoneIds = append(zoneIds, helper.IntInt64(v.(int)))
		}
		values.FieldByName("ZoneIds").Set(reflect.ValueOf(zoneIds))
	}

	if v, ok := d.GetOk("elastic_bandwidth_switch"); ok {
		values.FieldByName("ElasticBandwidthSwitch").Set(reflect.ValueOf(helper.Int64(int64(v.(int)))))
	}

	if v, ok := d.GetOk("custom_ssl_cert_id"); ok {
		values.FieldByName("CustomSSLCertId").Set(reflect.ValueOf(helper.String(v.(string))))
	}
}

func createCkafkaInstancePostPaid(ctx context.Context, d *schema.ResourceData, meta interface{}) (instanceId *string, err error) {
	logId := tccommon.GetLogId(ctx)
	request := ckafka.NewCreatePostPaidInstanceRequest()
	ckafkaRequestSetParams(request, d)
	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCkafkaClient().CreatePostPaidInstance(request)
	if err != nil {
		log.Printf("[CRITAL]%s create ckafka instance failed, reason:%s\n", logId, err.Error())
		return
	}
	if response.Response == nil || response.Response.Result.Data == nil {
		err = fmt.Errorf("CreatePostPaidInstance response is nil")
		return
	}
	instanceId = response.Response.Result.Data.InstanceId
	return
}
func createCkafkaInstancePrePaid(ctx context.Context, d *schema.ResourceData, meta interface{}) (instanceId *string, err error) {
	logId := tccommon.GetLogId(ctx)
	request := ckafka.NewCreateInstancePreRequest()
	ckafkaRequestSetParams(request, d)
	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCkafkaClient().CreateInstancePre(request)
	if err != nil {
		log.Printf("[CRITAL]%s create ckafka instance failed, reason:%s\n", logId, err.Error())
		return
	}
	if response.Response == nil || response.Response.Result.Data == nil {
		err = fmt.Errorf("CreateInstancePre response is nil")
		return
	}
	instanceId = response.Response.Result.Data.InstanceId
	return
}

func resourceTencentCloudCkafkaInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_instance.create")()
	var (
		instanceId *string
		createErr  error
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		service    = CkafkaService{
			client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
		}
		ctx = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	)

	chargeType := d.Get("charge_type").(string)
	if chargeType == CKAFKA_CHARGE_TYPE_POSTPAID {
		instanceId, createErr = createCkafkaInstancePostPaid(ctx, d, meta)
	} else if chargeType == CKAFKA_CHARGE_TYPE_PREPAID {
		instanceId, createErr = createCkafkaInstancePrePaid(ctx, d, meta)
	} else {
		return fmt.Errorf("invalid `charge_type` value")
	}
	if createErr != nil {
		return createErr
	}
	if instanceId == nil {
		return fmt.Errorf("instanceId is nil")
	}
	// wait sync instance
	time.Sleep(5 * time.Second)

	err := resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		has, ready, err := service.CheckCkafkaInstanceReady(ctx, *instanceId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if !has {
			return resource.NonRetryableError(fmt.Errorf("ckafka instance not exists."))
		}
		if ready {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("create ckafka instance task is processing"))
	})
	if err != nil {
		return err
	}
	d.SetId(*instanceId)

	// modify instance attributes
	var (
		needModify    = false
		modifyRequest = ckafka.NewModifyInstanceAttributesRequest()
	)
	modifyRequest.InstanceId = instanceId

	if v, ok := d.GetOk("msg_retention_time"); ok {
		needModify = true
		retentionTime := int64(v.(int))
		modifyRequest.MsgRetentionTime = helper.Int64(retentionTime)
	}

	if v, ok := d.GetOk("config"); ok {
		needModify = true
		config := make([]*ckafka.ModifyInstanceAttributesConfig, 0, 10)
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			configInfo := ckafka.ModifyInstanceAttributesConfig{}
			if autoCreateTopicEnable, ok := dMap["auto_create_topic_enable"]; ok {
				configInfo.AutoCreateTopicEnable = helper.Bool(autoCreateTopicEnable.(bool))
			}
			if defaultNumPartitions, ok := dMap["default_num_partitions"]; ok {
				configInfo.DefaultNumPartitions = helper.Int64(int64(defaultNumPartitions.(int)))
			}
			if defaultReplicationFactor, ok := dMap["default_replication_factor"]; ok {
				configInfo.DefaultReplicationFactor = helper.Int64(int64(defaultReplicationFactor.(int)))
			}
			config = append(config, &configInfo)
		}
		modifyRequest.Config = config[0]
	}

	if v, ok := d.GetOk("dynamic_retention_config"); ok {
		needModify = true
		dynamic := make([]*ckafka.DynamicRetentionTime, 0, 10)
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			dynamicInfo := ckafka.DynamicRetentionTime{}
			if enable, ok := dMap["enable"]; ok {
				dynamicInfo.Enable = helper.Int64(int64(enable.(int)))
			}
			if diskQuotaPercentage, ok := dMap["disk_quota_percentage"]; ok {
				dynamicInfo.DiskQuotaPercentage = helper.Int64(int64(diskQuotaPercentage.(int)))
			}
			if stepForwardPercentage, ok := dMap["step_forward_percentage"]; ok {
				dynamicInfo.StepForwardPercentage = helper.Int64(int64(stepForwardPercentage.(int)))
			}
			if bottomRetention, ok := dMap["bottom_retention"]; ok {
				dynamicInfo.BottomRetention = helper.Int64(int64(bottomRetention.(int)))
			}
			dynamic = append(dynamic, &dynamicInfo)
		}
		modifyRequest.DynamicRetentionConfig = dynamic[0]
	}

	if v, ok := d.GetOk("rebalance_time"); ok {
		needModify = true
		modifyRequest.RebalanceTime = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("public_network"); ok {
		needModify = true
		modifyRequest.PublicNetwork = helper.Int64(int64(v.(int)))
	}

	//if v, ok := d.GetOk("dynamic_disk_config"); ok {
	//	needModify = true
	//	dynamic := make([]*ckafka.DynamicDiskConfig, 0, 10)
	//	for _, item := range v.([]interface{}) {
	//		dMap := item.(map[string]interface{})
	//		dynamicInfo := ckafka.DynamicDiskConfig{}
	//		if enable, ok := dMap["enable"]; ok {
	//			dynamicInfo.Enable = helper.Int64(int64(enable.(int)))
	//		}
	//		if stepForwardPercentage, ok := dMap["step_forward_percentage"]; ok {
	//			dynamicInfo.StepForwardPercentage = helper.Int64(int64(stepForwardPercentage.(int)))
	//		}
	//		if diskQuotaPercentage, ok := dMap["disk_quota_percentage"]; ok {
	//			dynamicInfo.DiskQuotaPercentage = helper.Int64(int64(diskQuotaPercentage.(int)))
	//		}
	//		if maxDiskSpace, ok := dMap["max_disk_space"]; ok {
	//			dynamicInfo.MaxDiskSpace = helper.Int64(int64(maxDiskSpace.(int)))
	//		}
	//		dynamic = append(dynamic, &dynamicInfo)
	//	}
	//	modifyRequest.DynamicDiskConfig = dynamic[0]
	//}

	if v, ok := d.GetOkExists("max_message_byte"); ok {
		needModify = true
		modifyRequest.MaxMessageByte = helper.Uint64(uint64(v.(int)))
	}

	if needModify {
		err := service.ModifyCkafkaInstanceAttributes(ctx, modifyRequest)
		if err != nil {
			return fmt.Errorf("[API]Set kafka instance attributes fail, reason:%s", err.Error())
		}
	}

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(client)
	region := client.Region

	if tags := helper.GetTags(d, "tag_set"); len(tags) > 0 {
		resourceName := tccommon.BuildTagResourceName("ckafka", "ckafkaId", region, *instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudCkafkaInstanceRead(d, meta)
}

func resourceTencentCloudCkafkaInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var service = CkafkaService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	instanceId := d.Id()

	var info *ckafka.InstanceDetail
	var isExist = true

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		res, has, e := service.DescribeCkafkaInstanceById(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if !has {
			d.SetId("")
			isExist = false
			return nil
		}
		info = res
		return nil
	})
	if err != nil {
		return fmt.Errorf("[API]Describe kafka instance fail, reason:%s", err.Error())
	}

	if !isExist {
		return nil
	}

	if info.ExpireTime != nil {
		if *info.ExpireTime > 0 {
			_ = d.Set("charge_type", CKAFKA_CHARGE_TYPE_PREPAID)
		} else if *info.ExpireTime == 0 {
			_ = d.Set("charge_type", CKAFKA_CHARGE_TYPE_POSTPAID)
		} else {
			return fmt.Errorf("expireTime less than 0")
		}
	}
	_ = d.Set("instance_name", info.InstanceName)
	_ = d.Set("zone_id", info.ZoneId)
	_ = d.Set("vpc_id", info.VpcId)
	_ = d.Set("subnet_id", info.SubnetId)
	_ = d.Set("renew_flag", info.RenewFlag)
	_ = d.Set("kafka_version", info.Version)
	_ = d.Set("disk_size", info.DiskSize)
	bandWidth := info.Bandwidth
	_ = d.Set("vip", info.Vip)
	_ = d.Set("vport", info.Vport)
	_ = d.Set("band_width", *bandWidth/8)
	_ = d.Set("partition", info.MaxPartitionNumber)
	if *info.InstanceType == "profession" {
		_ = d.Set("specifications_type", "profession")
		_ = d.Set("instance_type", 1)
	} else {
		_ = d.Set("specifications_type", "standard")
		_ = d.Set("instance_type", CKAFKA_INSTANCE_TYPE[*info.InstanceType])
	}

	if len(info.ZoneIds) > 1 {
		_ = d.Set("multi_zone_flag", true)
		ids := helper.Int64sInterfaces(info.ZoneIds)
		idSet := schema.NewSet(func(i interface{}) int {
			return i.(int)
		}, ids)
		_ = d.Set("zone_ids", idSet)
	}

	tagSets := make([]map[string]interface{}, 0, len(info.Tags))
	for _, item := range info.Tags {
		tagSets = append(tagSets, map[string]interface{}{
			"key":   item.TagKey,
			"value": item.TagValue,
		})
	}
	_ = d.Set("tags", tagSets)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(client)
	region := client.Region

	tags, err := tagService.DescribeResourceTags(ctx, "ckafka", "ckafkaId", region, instanceId)
	if err != nil {
		return err
	}
	_ = d.Set("tag_set", tags)

	_ = d.Set("disk_type", info.DiskType)

	// query msg_retention_time
	var (
		request  = ckafka.NewDescribeInstanceAttributesRequest()
		response = ckafka.NewDescribeInstanceAttributesResponse()
	)
	request.InstanceId = &instanceId
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.client.UseCkafkaClient().DescribeInstanceAttributes(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		attr := response.Response.Result
		_ = d.Set("msg_retention_time", attr.MsgRetentionTime)

		if attr.Config != nil {
			config := make([]map[string]interface{}, 0)
			config = append(config, map[string]interface{}{
				"auto_create_topic_enable":   attr.Config.AutoCreateTopicsEnable,
				"default_num_partitions":     attr.Config.DefaultNumPartitions,
				"default_replication_factor": attr.Config.DefaultReplicationFactor,
			})
			_ = d.Set("config", config)
		}

		dynamicConfig := make([]map[string]interface{}, 0)
		dynamicConfig = append(dynamicConfig, map[string]interface{}{
			"enable":                  attr.RetentionTimeConfig.Enable,
			"disk_quota_percentage":   attr.RetentionTimeConfig.DiskQuotaPercentage,
			"step_forward_percentage": attr.RetentionTimeConfig.StepForwardPercentage,
			"bottom_retention":        attr.RetentionTimeConfig.BottomRetention,
		})
		_ = d.Set("dynamic_retention_config", dynamicConfig)
		_ = d.Set("public_network", attr.PublicNetwork)

		//dynamicDiskConfig := make([]map[string]interface{}, 0)
		//dynamicDiskConfig = append(dynamicDiskConfig, map[string]interface{}{
		//	"enable":                  attr.DynamicDiskConfig.Enable,
		//	"disk_quota_percentage":   attr.DynamicDiskConfig.DiskQuotaPercentage,
		//	"step_forward_percentage": attr.DynamicDiskConfig.StepForwardPercentage,
		//	"max_disk_space":          attr.DynamicDiskConfig.MaxDiskSpace,
		//})
		//_ = d.Set("dynamic_disk_config", dynamicDiskConfig)

		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Ckafka Instance Attributes failed, reason:%s\n", logId, err.Error())
		return err
	}
	return nil
}

func resourceTencentCloudCkafkaInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_instance.update")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := CkafkaService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	immutableArgs := []string{
		"zone_id", "period", "vpc_id",
		"subnet_id", "renew_flag", "kafka_version",
		"multi_zone_flag", "zone_ids", "disk_type",
		"specifications_type", "instance_type",
		"elastic_bandwidth_switch", "custom_ssl_cert_id",
	}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	instanceId := d.Id()
	modifyInstanceAttributesFlag := false
	request := ckafka.NewModifyInstanceAttributesRequest()
	request.InstanceId = &instanceId
	if d.HasChange("instance_name") {
		if v, ok := d.GetOk("instance_name"); ok {
			request.InstanceName = helper.String(v.(string))
			modifyInstanceAttributesFlag = true
		}
	}

	if d.HasChange("msg_retention_time") {
		if v, ok := d.GetOk("msg_retention_time"); ok {
			request.MsgRetentionTime = helper.Int64(int64(v.(int)))
			modifyInstanceAttributesFlag = true
		}
	}

	if d.HasChange("config") {
		if v, ok := d.GetOk("config"); ok {
			items := v.([]interface{})
			dMap := items[0].(map[string]interface{})
			configInfo := ckafka.ModifyInstanceAttributesConfig{}
			if autoCreateTopicEnable, ok := dMap["auto_create_topic_enable"]; ok {
				configInfo.AutoCreateTopicEnable = helper.Bool(autoCreateTopicEnable.(bool))
			}
			if defaultNumPartitions, ok := dMap["default_num_partitions"]; ok {
				configInfo.DefaultNumPartitions = helper.Int64(int64(defaultNumPartitions.(int)))
			}
			if defaultReplicationFactor, ok := dMap["default_replication_factor"]; ok {
				configInfo.DefaultReplicationFactor = helper.Int64(int64(defaultReplicationFactor.(int)))
			}
			request.Config = &configInfo
			modifyInstanceAttributesFlag = true
		}
	}

	if d.HasChange("dynamic_retention_config") {
		if v, ok := d.GetOk("dynamic_retention_config"); ok {
			items := v.([]interface{})
			dMap := items[0].(map[string]interface{})
			dynamicInfo := ckafka.DynamicRetentionTime{}
			if enable, ok := dMap["enable"]; ok {
				dynamicInfo.Enable = helper.Int64(int64(enable.(int)))
			}
			if diskQuotaPercentage, ok := dMap["disk_quota_percentage"]; ok {
				dynamicInfo.DiskQuotaPercentage = helper.Int64(int64(diskQuotaPercentage.(int)))
			}
			if stepForwardPercentage, ok := dMap["step_forward_percentage"]; ok {
				dynamicInfo.StepForwardPercentage = helper.Int64(int64(stepForwardPercentage.(int)))
			}
			if bottomRetention, ok := dMap["bottom_retention"]; ok {
				dynamicInfo.BottomRetention = helper.Int64(int64(bottomRetention.(int)))
			}
			request.DynamicRetentionConfig = &dynamicInfo
			modifyInstanceAttributesFlag = true
		}
	}

	if d.HasChange("rebalance_time") {
		if v, ok := d.GetOk("rebalance_time"); ok {
			request.RebalanceTime = helper.Int64(int64(v.(int)))
			modifyInstanceAttributesFlag = true
		}
	}

	if d.HasChange("public_network") {
		if v, ok := d.GetOk("public_network"); ok {
			request.PublicNetwork = helper.Int64(int64(v.(int)))
			modifyInstanceAttributesFlag = true
		}
	}

	if d.HasChange("max_message_byte") {
		if v, ok := d.GetOkExists("max_message_byte"); ok {
			request.MaxMessageByte = helper.Uint64(uint64(v.(int)))
			modifyInstanceAttributesFlag = true
		}
	}

	if modifyInstanceAttributesFlag {
		err := service.ModifyCkafkaInstanceAttributes(ctx, request)
		if err != nil {
			return fmt.Errorf("[API]Set kafka instance attributes fail, reason:%s", err.Error())
		}
	}

	if d.HasChange("band_width") || d.HasChange("disk_size") || d.HasChange("partition") {
		chargeType := d.Get("charge_type").(string)
		if chargeType == CKAFKA_CHARGE_TYPE_POSTPAID {
			request := ckafka.NewInstanceScalingDownRequest()
			request.InstanceId = helper.String(instanceId)
			upgradeStrategy := d.Get("upgrade_strategy").(int)
			request.UpgradeStrategy = helper.IntInt64(upgradeStrategy)
			if v, ok := d.GetOk("band_width"); ok && d.HasChange("band_width") {
				request.BandWidth = helper.Int64(int64(v.(int)))
			}
			if v, ok := d.GetOk("disk_size"); ok && d.HasChange("disk_size") {
				request.DiskSize = helper.Int64(int64(v.(int)))
			}
			if v, ok := d.GetOk("partition"); ok && d.HasChange("partition") {
				request.Partition = helper.Int64(int64(v.(int)))
			}

			_, err := service.client.UseCkafkaClient().InstanceScalingDown(request)
			if err != nil {
				return fmt.Errorf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]", logId,
					request.GetAction(), request.ToJsonString(), err.Error())
			}
		} else if chargeType == CKAFKA_CHARGE_TYPE_PREPAID {
			request := ckafka.NewModifyInstancePreRequest()
			request.InstanceId = helper.String(instanceId)

			if v, ok := d.GetOk("band_width"); ok {
				request.BandWidth = helper.Int64(int64(v.(int)))
			}
			if v, ok := d.GetOk("disk_size"); ok {
				request.DiskSize = helper.Int64(int64(v.(int)))
			}
			if v, ok := d.GetOk("partition"); ok {
				request.Partition = helper.Int64(int64(v.(int)))
			}

			_, err := service.client.UseCkafkaClient().ModifyInstancePre(request)
			if err != nil {
				return fmt.Errorf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]", logId,
					request.GetAction(), request.ToJsonString(), err.Error())
			}
		} else {
			return fmt.Errorf("invalid `charge_type` value")
		}
		// InstanceScalingDown statue delay
		time.Sleep(5 * time.Second)

		err := resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, ready, err := service.CheckCkafkaInstanceReady(ctx, instanceId)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if ready {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("upgrade ckafka instance task is processing"))
		})

		if err != nil {
			return err
		}
	}

	if d.HasChange("tag_set") {

		client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(client)
		region := client.Region

		oldTags, newTags := d.GetChange("tag_set")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		resourceName := tccommon.BuildTagResourceName("ckafka", "ckafkaId", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudCkafkaInstanceRead(d, meta)
}

func resourceTencentCLoudCkafkaInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_instance.delete")()
	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = CkafkaService{
			client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
		}
	)
	instanceId := d.Id()
	chargeType := d.Get("charge_type").(string)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		if chargeType == CKAFKA_CHARGE_TYPE_POSTPAID {
			request := ckafka.NewDeleteInstancePostRequest()
			request.InstanceId = &instanceId
			_, err := service.client.UseCkafkaClient().DeleteInstancePost(request)
			if err != nil {
				return tccommon.RetryError(err, "UnsupportedOperation")
			}

		} else if chargeType == CKAFKA_CHARGE_TYPE_PREPAID {
			request := ckafka.NewDeleteInstancePreRequest()
			request.InstanceId = &instanceId
			_, err := service.client.UseCkafkaClient().DeleteInstancePre(request)
			if err != nil {
				return tccommon.RetryError(err, "UnsupportedOperation")
			}
		} else {
			return resource.NonRetryableError(fmt.Errorf("invalid `charge_type` value"))
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		has, _, err := service.CheckCkafkaInstanceReady(ctx, instanceId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if !has {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("delete ckafka instance task is processing"))
	})
	if err != nil {
		return err
	}
	return nil
}
