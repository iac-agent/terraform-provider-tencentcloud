package crs

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

//internal version: replace import begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
//internal version: replace import end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

func ResourceTencentCloudRedisInstance() *schema.Resource {
	types := []string{}
	for _, v := range REDIS_NAMES {
		types = append(types, "`"+v+"`")
	}
	sort.Strings(types)
	typeStr := strings.Trim(strings.Join(types, ","), ",")

	return &schema.Resource{
		Create: resourceTencentCloudRedisInstanceCreate,
		Read:   resourceTencentCloudRedisInstanceRead,
		Update: resourceTencentCloudRedisInstanceUpdate,
		Delete: resourceTencentCloudRedisInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "可用 可用区 的 实例 到 是 创建，like `ap-beijing-7`，please refer 到 `tencentcloud_redis_zone_config.列表`。",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "实例名称",
			},
			"type_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateIntegerMin(2),
				Description: "实例 类型. Available 值 reference 数据 source `tencentcloud_redis_zone_config` 或 [document](https://intl.云.tencent.com/document/product/239/32069), toggle immediately 当 modified." +
					"<ul><li>2: Redis 2.8 Memory Edition (standard architecture);</li> " +
					"<li>3: CKV 3.2 Memory Edition (standard architecture);</li> " +
					"<li>4: CKV 3.2 Memory Edition (cluster architecture);</li> " +
					"<li>6: Redis 4.0 Memory Edition (standard architecture);</li> " +
					"<li>7: Redis 4.0 Memory Edition (cluster architecture);</li> " +
					"<li>8: Redis 5.0 Memory Edition (standard architecture);</li> " +
					"<li>9: Redis 5.0 Memory Edition (cluster architecture);</li> " +
					"<li>15: Redis 6.2 Memory Edition (standard architecture);</li> " +
					"<li>16: Redis 6.2 Memory Edition (cluster architecture);</li> " +
					"<li>17: Redis 7.0 Memory Edition (standard architecture);</li> " +
					"<li>18: Redis 7.0 Memory Edition (cluster architecture). </li> " +
					"<li>200: Memcached 1.6 Memory Edition (cluster architecture). </li>Note: The CKV version is currently used by existing users and is temporarily retained.</ul>.",
			},
			"redis_shard_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{1, 3, 5, 8, 12, 16, 24, 32, 40, 48, 64, 80, 96, 128}),
				Description:  "数量 实例 shards; 此 参数 does 不 need 到 是 已配置 对于 standard 版本 实例; 对于 集群 版本 实例， 数量 shards ranges 从: [`1`，`3`，`5`，`8`，`12`，`16`，`24 `，`32`，`40`，`48`，`64`，`80`，`96`，`128`]。",
			},
			"redis_replicas_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{1, 2, 3, 4, 5}),
				Description:  "数量 实例 copies. 此 不是必填项 对于 standalone 和 master slave versions 和 必须 equal 到 count 的 `replica_zone_ids`，Non-multi-AZ does 不 require `replica_zone_ids`; Redis 内存 版本 4.0，5.0，6.2 standard architecture 和 集群 architecture support 数量 copies 在 范围 [1，2，3，4，5]; Redis 2.8 standard 版本 和 CKV standard 版本 仅 support 1 copy。",
			},
			"replica_zone_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "ID 副本 nodes 可用 可用区 此 不是必填项 对于 standalone 和 master slave versions. NOTE: Removing some 的 same 可用区 的 replicas (e.g. removing 100001 的 [100001，100001，100002]) 将 pick first hit 到 remove。",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					for _, name := range REDIS_NAMES {
						if name == value {
							return
						}
					}
					errors = append(errors, fmt.Errorf("this redis type %s not support now", value))
					return
				},
				Deprecated:  "It has been deprecated from version 1.33.1. Please use 'type_id' instead.",
				Description: "实例 类型. Available 值:" + typeStr + ", specific region support specific types, need to refer data `tencentcloud_redis_zone_config`.",
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: tccommon.ValidateMysqlPassword,
				Description:  "密码 对于 Redis 用户，其中 should 是 8 到 16 字符. NOTE: Only `no_auth=true` 指定 可以 make 密码 空。",
			},
			"no_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "表示是否redis 实例 support 无-auth 访问. NOTE: Only 可用 在 私有 云 环境。",
			},
			"replicas_read_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether copy read-仅 是 支持，Redis 2.8 Standard Edition 和 CKV Standard Edition do 不 support 副本 read-仅，turn 在 副本 read-仅， 实例 将 automatically read 和 write separate，write requests 是 routed 到 primary 节点，read requests 是 routed 到 副本 节点，如果 您 need 到 open 副本 read-仅， recommended 数量 replicas >=2。",
			},
			"mem_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "内存 卷 的 可用 实例(在 MB)，please refer 到 `tencentcloud_redis_zone_config.列表[可用区].shard_memories`. 当 redis 是 standard 类型，它 表示 总数 内存 大小 的 实例; 当 Redis 是 集群类型，它 表示 内存 大小 的 per sharding. `512MB` 是 支持 仅 在 master-slave 实例。",
			},
			"vpc_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 100),
				Description:  "ID vpc 使用 其中 实例 是 到 是 associated. 当 `operation_network` 是 `changeVpc` 或 `changeBaseToVpc`，此 参数 needs 到 是 已配置。",
			},
			"subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 100),
				Description:  "指定which 子网 实例 should belong 到. 当 `operation_network` 是 `changeVpc` 或 `changeBaseToVpc`，此 参数 needs 到 是 已配置。",
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set: func(v interface{}) int {
					return helper.HashString(v.(string))
				},
				Description: "ID 安全 组. 如果 both vpc_id 和 subnet_id 是 不 集合，此 argument should 不 是 集合 either。",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "指定which 项目 实例 should belong 到。",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     6379,
				Description: "端口 用于access redis 实例. 默认值为 6379. 当 `operation_network` 是 `changeVPort` 或 `changeVip`，此 参数 needs 到 是 已配置。",
			},
			"params_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "指定params template ID. 如果 不 集合，将 使用 默认值 template。",
			},
			"operation_network": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(REDIS_MODIFY_NETWORK_CONFIG),
				Description:  "Refers 到 category 的 pre-modified 网络，包括: `changeVip`: refers 到 switching 私有 网络，包括 its intranet IPv4 地址 和 端口; `changeVpc`: refers 到 switching 子网 到 其中 私有 网络 belongs; `changeBaseToVpc`: refers 到 switching basic 网络 到 私有 网络; `changeVPort`: refers 到 仅 modifying 实例 网络 端口",
			},
			"recycle": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue(REDIS_RECYCLE_TIME),
				Description:  "Original intranet IPv4 地址 retention 时间: 单位: day，取值范围：`0`，`1`，`2`，`3`，`7`，`15`。",
			},
			"ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "IP 地址 的 实例. 当 `operation_network` 是 `changeVip`，此 参数 needs 到 是 已配置。",
			},
			"wait_switch": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Switching 模式: `1`-maintenance 时间 window switching，`2`-immediate switching，默认值 `2`。",
			},
			"product_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "指定product 版本 的 实例. `本地`: Local 磁盘 版本，`云`: Cloud 磁盘 版本，`cdc`: Exclusive 集群 版本 默认为 `本地`。",
			},
			"redis_cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Exclusive 集群 ID. 当 `product_version` 是 集合 到 `cdc`，此 参数 必须 是 集合。",
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				//internal version: replace tagComputed begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
				//internal version: replace tagComputed end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
				Description: "实例 标签",
			},
			// Computed values
			"dedicated_cluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Dedicated 集群 ID",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current 状态 实例，maybe: init，processing，online，isolate 和 todelete。",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "时间 当 实例 是 创建。",
			},
			// payment
			"charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      REDIS_CHARGE_TYPE_POSTPAID,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{REDIS_CHARGE_TYPE_POSTPAID, REDIS_CHARGE_TYPE_PREPAID}),
				Description:  "charge 类型 实例. 有效值：`PREPAID` 和 `POSTPAID`. 默认值为 `POSTPAID`。",
			},
			"prepaid_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue(REDIS_PREPAID_PERIOD),
				Description:  "tenancy (时间 单位 是 month) 的 prepaid 实例，NOTE: 它 仅 works 当 charge_type 是 集合 到 `PREPAID`. 有效 值 是 `1`，`2`，`3`，`4`，`5`，`6`，`7`，`8`，`9`，`10`，`11`，`12`，`24`，`36`。",
			},
			"force_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicate 是否delete Redis 实例 directly 或 不. 默认为 false. 如果 集合 true， 实例 将 是 删除 instead 的 staying recycle bin。",
			},
			"auto_renew_flag": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1, 2}),
				Default:      0,
				Description:  "Auto-续费标识 0 - 默认值 state (manual renewal); 1 - automatic renewal; 2 - explicit 无 automatic renewal。",
			},
			"node_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Readonly Primary/Replica nodes。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"master": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "表示是否node 是 master。",
						},
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID master 或 副本 节点。",
						},
						"zone_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID availability 可用区 的 master 或 副本 节点。",
						},
					},
				},
			},
			"wan_address_switch": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Wan 地址 switch，默认值 `close`，值: `open`，`close`。",
			},
			"wan_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Allocate Wan 地址",
			},
		},
	}
}

func resourceTencentCloudRedisInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_instance.create")()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		client       = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		redisService = RedisService{client: client}
		tagService   = svctag.NewTagService(client)
		region       = client.Region
	)

	//internal version: replace clientCreate begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	//internal version: replace clientCreate end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

	availabilityZone := d.Get("availability_zone").(string)
	redisName := d.Get("name").(string)
	redisType := d.Get("type").(string)
	typeId := int64(d.Get("type_id").(int))
	redisShardNum := 1
	if v, ok := d.GetOk("redis_shard_num"); ok {
		redisShardNum = v.(int)
	}

	redisReplicasNum := d.Get("redis_replicas_num").(int)
	password := d.Get("password").(string)
	noAuth := d.Get("no_auth").(bool)
	memSize := d.Get("mem_size").(int)
	vpcId := d.Get("vpc_id").(string)
	subnetId := d.Get("subnet_id").(string)
	securityGroups := d.Get("security_groups").(*schema.Set).List()
	projectId := d.Get("project_id").(int)
	port := d.Get("port").(int)
	chargeType := d.Get("charge_type").(string)
	autoRenewFlag := d.Get("auto_renew_flag").(int)
	paramsTemplateId := d.Get("params_template_id").(string)
	operation := d.Get("operation_network").(string)
	productVersion := d.Get("product_version").(string)
	redisClusterId := d.Get("redis_cluster_id").(string)
	if productVersion == "cdc" && redisClusterId == "" {
		return fmt.Errorf("If `product_version` is set to `cdc`, params `redis_cluster_id` must be set\n.")
	}

	chargeTypeID := REDIS_CHARGE_TYPE_ID[chargeType]
	var replicasReadonly bool
	if v, ok := d.GetOk("replicas_read_only"); ok {
		replicasReadonly = v.(bool)
	}

	var chargePeriod uint64 = 1
	if chargeType == REDIS_CHARGE_TYPE_PREPAID {
		if period, ok := d.GetOk("prepaid_period"); ok {
			chargePeriod = uint64(period.(int))
		} else {
			return fmt.Errorf("instance charge type prepaid period can not be empty when charge type is %s", chargeType)
		}
	}

	if (typeId == 0 && redisType == "") || (typeId != 0 && redisType != "") {
		return fmt.Errorf("`type_id` and `type` set one item and only one item")
	}

	if password == "" && !noAuth {
		return fmt.Errorf("`password` must not be empty unless `no_auth` is `true`")
	}

	if noAuth && (vpcId == "" || subnetId == "") {
		return fmt.Errorf("cannot set `no_auth=true` if `vpc_id` and `subnet_id` is empty")
	}

	for id, name := range REDIS_NAMES {
		if redisType == name {
			typeId = id
			break
		}
	}

	sellConfigures, err := redisService.DescribeRedisZoneConfig(ctx)
	if err != nil {
		return fmt.Errorf("api[DescribeRedisZoneConfig]fail, return %s", err.Error())
	}

	var regionItem *redis.RegionConf
	var zoneItem *redis.ZoneCapacityConf
	var redisItem *redis.ProductConf
	for _, regionItem = range sellConfigures {
		if *regionItem.RegionId == region {
			break
		}
	}

	if regionItem == nil {
		return fmt.Errorf("all redis in this region `%s` be sold out", region)
	}

	for _, zones := range regionItem.ZoneSet {
		if *zones.IsSaleout {
			continue
		}

		if *zones.ZoneName == availabilityZone {
			zoneItem = zones
			break
		}
	}

	if zoneItem == nil {
		return fmt.Errorf("all redis in this zone `%s` be sold out", availabilityZone)
	}

	for _, reds := range zoneItem.ProductSet {
		if *reds.Type == typeId {
			redisItem = reds
			break
		}
	}

	if redisItem == nil {
		return fmt.Errorf("redis type_id `%d` be sold out or this type_id is not supports", typeId)
	}

	var redisShardNums []string
	var redisReplicasNums []string
	var numErrors []string
	for _, v := range redisItem.ShardNum {
		redisShardNums = append(redisShardNums, *v)
	}

	for _, v := range redisItem.ReplicaNum {
		redisReplicasNums = append(redisReplicasNums, *v)
	}

	if !tccommon.IsContains(redisShardNums, fmt.Sprintf("%d", redisShardNum)) {
		numErrors = append(numErrors, fmt.Sprintf("redis_shard_num : %s", strings.Join(redisShardNums, ",")))
	}

	if !tccommon.IsContains(redisReplicasNums, fmt.Sprintf("%d", redisReplicasNum)) {
		numErrors = append(numErrors, fmt.Sprintf(" redis_replicas_num : %s", strings.Join(redisReplicasNums, ",")))
	}

	if len(numErrors) > 0 {
		return fmt.Errorf("redis type_id `%d` only supports %s", typeId, strings.Join(numErrors, ","))
	}

	if operation != "" {
		return fmt.Errorf("This parameter `operation_network` is not required when redis is created")
	}

	requestSecurityGroup := make([]string, 0, len(securityGroups))

	for _, v := range securityGroups {
		requestSecurityGroup = append(requestSecurityGroup, v.(string))
	}

	//internal version: replace null begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	//internal version: replace null end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

	nodeInfo := make([]*redis.RedisNodeInfo, 0)
	if raw, ok := d.GetOk("replica_zone_ids"); ok {
		zoneIds := raw.([]interface{})
		//internal version: replace redisServer begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
		masterZoneId, err := service.getZoneId(availabilityZone)
		//internal version: replace redisServer end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

		if err != nil {
			return err
		}

		// insert master node
		nodeInfo = append(nodeInfo, &redis.RedisNodeInfo{
			NodeType: helper.Int64(0),
			ZoneId:   helper.Int64Uint64(masterZoneId),
		})

		for _, v := range zoneIds {
			id := v.(int)
			nodeInfo = append(nodeInfo, &redis.RedisNodeInfo{
				NodeType: helper.Int64(1),
				ZoneId:   helper.IntUint64(id),
			})
		}
	}

	instanceIds, err := redisService.CreateInstances(ctx,
		availabilityZone,
		typeId,
		password,
		vpcId,
		subnetId,
		redisName,
		int64(memSize),
		int64(projectId),
		int64(port),
		requestSecurityGroup,
		redisShardNum,
		redisReplicasNum,
		chargeTypeID,
		chargePeriod,
		nodeInfo,
		noAuth,
		autoRenewFlag,
		replicasReadonly,
		paramsTemplateId,
		productVersion,
		redisClusterId,
	)

	//internal version: replace varId begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	//internal version: replace varId end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	if err != nil {
		//internal version: replace bpass begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
		return err
		//internal version: replace bpass end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	}

	if len(instanceIds) == 0 {
		return fmt.Errorf("redis api CreateInstances return empty redis id")
	}

	//internal version: replace getId begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	var redisId = *instanceIds[0]
	//internal version: replace getId end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

	//internal version: replace setTag begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	//internal version: replace setTag end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

	//internal version: replace queryAndSetId begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	_, _, _, err = redisService.CheckRedisOnlineOk(ctx, redisId, 20*tccommon.ReadRetryTimeout)

	if err != nil {
		log.Printf("[CRITAL]%s create redis task fail, reason:%s\n", logId, err.Error())
		return err
	}

	d.SetId(redisId)
	//internal version: replace queryAndSetId end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

	//internal version: replace null begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		resourceName := tccommon.BuildTagResourceName("redis", "instance", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	//internal version: replace null end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

	if v, ok := d.GetOk("wan_address_switch"); ok {
		err := resourceRedisWanAddressModify(ctx, &redisService, meta, d.Id(), v.(string))
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudRedisInstanceRead(d, meta)
}

func resourceTencentCloudRedisInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	//internal version: replace clientRead begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	//internal version: replace clientRead end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

	var onlineHas = true
	var (
		has  bool
		info *redis.InstanceSet
		e    error
	)

	if v, ok := d.GetOkExists("wait_switch"); ok && v.(int) == 1 {
		info, e = service.DescribeRedisInstanceById(ctx, d.Id())
		if e != nil {
			return e
		}
	} else {
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			has, _, info, e = service.CheckRedisOnlineOk(ctx, d.Id(), tccommon.ReadRetryTimeout*20)
			if info != nil {
				if *info.Status == REDIS_STATUS_ISOLATE || *info.Status == REDIS_STATUS_TODELETE {
					d.SetId("")
					onlineHas = false
					return nil
				}
			}
			if e != nil {
				return resource.NonRetryableError(e)
			}
			if !has {
				d.SetId("")
				onlineHas = false
				return nil
			}
			return nil
		})
		if err != nil {
			//internal version: replace redisFail begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
			return fmt.Errorf("Fail to get info from redis, reaseon %s", err.Error())
			//internal version: replace redisFail end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
		}
		if !onlineHas {
			return nil
		}
	}

	statusName := REDIS_STATUS[*info.Status]
	if statusName == "" {
		err := fmt.Errorf("redis read unkwnow status %d", *info.Status)
		log.Printf("[CRITAL]%s redis read status name error, reason:%s\n", logId, err.Error())
		return err
	}
	_ = d.Set("status", statusName)

	_ = d.Set("name", *info.InstanceName)

	zoneName, err := service.getZoneName(*info.ZoneId)
	if err != nil {
		return err
	}
	// not set field type_id
	// process import case
	if d.Get("type_id").(int) == 0 && d.Get("type").(string) != "" {
		typeName := REDIS_NAMES[*info.Type]
		if typeName == "" {
			err = fmt.Errorf("redis read unkwnow type %d", *info.Type)
			log.Printf("[CRITAL]%s redis read type name error, reason:%s\n", logId, err.Error())
			return err
		}
		_ = d.Set("type", typeName)
	} else {
		_ = d.Set("type_id", info.Type)
	}

	_ = d.Set("redis_shard_num", info.RedisShardNum)
	_ = d.Set("redis_replicas_num", info.RedisReplicasNum)
	_ = d.Set("availability_zone", zoneName)
	_ = d.Set("mem_size", info.RedisShardSize)
	_ = d.Set("vpc_id", info.UniqVpcId)
	_ = d.Set("subnet_id", info.UniqSubnetId)
	_ = d.Set("project_id", info.ProjectId)
	_ = d.Set("port", info.Port)
	_ = d.Set("ip", info.WanIp)
	_ = d.Set("create_time", info.Createtime)
	_ = d.Set("auto_renew_flag", info.AutoRenewFlag)
	_ = d.Set("product_version", info.ProductVersion)
	_ = d.Set("redis_cluster_id", info.RedisClusterId)
	_ = d.Set("dedicated_cluster_id", info.DedicatedClusterId)
	slaveReadWeight := *info.SlaveReadWeight
	if slaveReadWeight == 0 {
		_ = d.Set("replicas_read_only", false)
	} else if slaveReadWeight == 100 {
		_ = d.Set("replicas_read_only", true)
	}

	// only true or user explicit declared will set for import case.
	if _, ok := d.GetOk("no_auth"); ok || *info.NoAuth {
		_ = d.Set("no_auth", info.NoAuth)
	}

	if d.Get("vpc_id").(string) != "" {
		securityGroups, err := service.DescribeInstanceSecurityGroup(ctx, d.Id())
		if err != nil {
			return err
		}
		if len(securityGroups) > 0 {
			_ = d.Set("security_groups", securityGroups)
		}
	}

	if info.NodeSet != nil {
		var zoneIds []int
		var nodeInfos []interface{}
		for i := range info.NodeSet {
			nodeInfo := info.NodeSet[i]
			nodeInfos = append(nodeInfos, map[string]interface{}{
				"master":  *nodeInfo.NodeType == 0,
				"zone_id": *nodeInfo.ZoneId,
				"id":      *nodeInfo.NodeId,
			})
			if *nodeInfo.NodeType == 0 {
				continue
			}
			zoneIds = append(zoneIds, int(*nodeInfo.ZoneId))
		}

		_ = d.Set("node_info", nodeInfos)

		var zoneIdsEqual = false

		replicaZones, replicaZonesOk := d.GetOk("replica_zone_ids")
		if replicaZonesOk {
			oldIds := helper.InterfacesIntegers(replicaZones.([]interface{}))
			zoneIdsEqual = checkIdsEqual(oldIds, zoneIds)
		}

		if !zoneIdsEqual {
			_ = d.Set("replica_zone_ids", zoneIds)
		}
	}
	//internal version: replace resourceTag begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "redis", "instance", tcClient.Region, d.Id())
	//internal version: replace resourceTag end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	_ = d.Set("charge_type", REDIS_CHARGE_TYPE_NAME[*info.BillingMode])

	if info.WanAddress != nil && *info.WanAddress != "" {
		_ = d.Set("wan_address", info.WanAddress)
		_ = d.Set("wan_address_switch", "open")
	} else {
		_ = d.Set("wan_address_switch", "close")
	}

	return nil
}

func resourceTencentCloudRedisInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_redis_instance.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	redisService := RedisService{client: client}
	tagService := svctag.NewTagService(client)
	region := client.Region

	d.Partial(true)

	unsupportedUpdateFields := []string{
		"product_version",
		"redis_cluster_id",
	}
	for _, field := range unsupportedUpdateFields {
		if d.HasChange(field) {
			return fmt.Errorf("tencentcloud_redis_instance update on %s is not support yet", field)
		}
	}

	// charge_type
	if d.HasChange("charge_type") {
		newChargeType := d.Get("charge_type").(string)
		period := 0
		if newChargeType == REDIS_CHARGE_TYPE_PREPAID {
			if v, ok := d.GetOk("prepaid_period"); ok {
				period = v.(int)
			} else {
				return fmt.Errorf("prepaid_period must be set when charge_type is changed to PREPAID")
			}
		}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := redisService.ModifyInstanceChargeType(ctx, id, newChargeType, period)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s redis ModifyInstanceChargeType error, reason:%s\n", logId, err.Error())
			return err
		}
	}

	// name\mem_size\password\project_id

	if d.HasChange("name") {
		name := d.Get("name").(string)
		if name == "" {
			name = id
		}
		err := redisService.ModifyInstanceName(ctx, id, name)
		if err != nil {
			return err
		}
	}

	// MemSize, ShardNum and ReplicaNum can only change one for each upgrade invoke
	if d.HasChange("mem_size") {

		_, newInter := d.GetChange("mem_size")
		newMemSize := newInter.(int)
		oShard, _ := d.GetChange("redis_shard_num")
		redisShardNum := oShard.(int)
		oReplica, _ := d.GetChange("redis_replicas_num")
		redisReplicasNum := oReplica.(int)

		if newMemSize < 1 {
			return fmt.Errorf("redis mem_size value cannot be set to less than 1")
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			_, err := redisService.UpgradeInstance(ctx, id, newMemSize, redisShardNum, redisReplicasNum, nil)
			if err != nil {
				// Upgrade memory will cause instance lock and cannot acknowledge by polling status, wait until lock release
				return tccommon.RetryError(err, redis.FAILEDOPERATION_UNKNOWN, redis.FAILEDOPERATION_SYSTEMERROR)
			}
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s redis upgrade instance error, reason:%s\n", logId, err.Error())
			return err
		}

		err = redisService.CheckRedisUpdateOk(ctx, id)

		if err != nil {
			log.Printf("[CRITAL]%s redis update mem size fail , reason:%s\n", logId, err.Error())
			return err
		}

		// temp solution for wait
		if d.HasChange("redis_shard_num") || d.HasChange("redis_replicas_num") || d.HasChange("replica_zone_ids") || d.HasChange("type_id") {
			time.Sleep(time.Minute * 2)
		}
	}

	// MemSize, ShardNum and ReplicaNum can only change one for each upgrade invoke
	if d.HasChange("redis_shard_num") {
		redisShardNum := d.Get("redis_shard_num").(int)
		oReplica, _ := d.GetChange("redis_replicas_num")
		redisReplicasNum := oReplica.(int)
		memSize := d.Get("mem_size").(int)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			_, err := redisService.UpgradeInstance(ctx, id, memSize, redisShardNum, redisReplicasNum, nil)
			if err != nil {
				// Upgrade memory will cause instance lock and cannot acknowledge by polling status, wait until lock release
				return tccommon.RetryError(err, redis.FAILEDOPERATION_UNKNOWN, redis.FAILEDOPERATION_SYSTEMERROR)
			}
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s redis upgrade instance error, reason:%s\n", logId, err.Error())
			return err
		}

		err = redisService.CheckRedisUpdateOk(ctx, id)

		if err != nil {
			log.Printf("[CRITAL]%s redis update shard num fail , reason:%s\n", logId, err.Error())
			return err
		}

		// temp solution for wait
		if d.HasChange("redis_replicas_num") || d.HasChange("replica_zone_ids") || d.HasChange("type_id") {
			time.Sleep(time.Minute * 2)
		}
	}

	if d.HasChange("redis_replicas_num") || d.HasChange("replica_zone_ids") {
		err := resourceRedisNodeSetModify(ctx, &redisService, d)
		if err != nil {
			return err
		}

		// temp solution for wait
		if d.HasChange("type_id") {
			time.Sleep(time.Minute * 2)
		}
	}

	if d.HasChange("password") || d.HasChange("no_auth") {
		var (
			taskId   int64
			password = d.Get("password").(string)
			noAuth   = d.Get("no_auth").(bool)
			err      error
		)

		// After redis spec modified, reset password may not successfully response immediately.
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			taskId, err = redisService.ResetPassword(ctx, id, password, noAuth)
			if err != nil {
				log.Printf("[CRITAL]%s redis change password error, reason:%s\n", logId, err.Error())
				return tccommon.RetryError(err, redis.FAILEDOPERATION_SYSTEMERROR)
			}
			return nil
		})

		if err != nil {
			return err
		}

		err = resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ok, err := redisService.DescribeTaskInfo(ctx, id, taskId)
			if err != nil {
				if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
					return resource.RetryableError(err)
				} else {
					return resource.NonRetryableError(err)
				}
			}
			if ok {
				return nil
			} else {
				return resource.RetryableError(fmt.Errorf("change password is processing"))
			}
		})

		if err != nil {
			log.Printf("[CRITAL]%s redis change password fail, reason:%s\n", logId, err.Error())
			return err
		}
	}

	if d.HasChange("params_template_id") {
		request := redis.NewApplyParamsTemplateRequest()
		request.InstanceIds = []*string{&id}
		request.TemplateId = helper.String(d.Get("params_template_id").(string))
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			_, err := redisService.ApplyParamsTemplate(ctx, request)
			if err != nil {
				return tccommon.RetryError(err, redis.FAILEDOPERATION_SYSTEMERROR, redis.RESOURCEUNAVAILABLE_INSTANCELOCKEDERROR)
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	if d.HasChange("project_id") {
		projectId := d.Get("project_id").(int)
		err := redisService.ModifyInstanceProjectId(ctx, id, int64(projectId))
		if err != nil {
			return err
		}
	}

	if d.HasChanges("security_groups") {
		sgs := d.Get("security_groups").(*schema.Set).List()
		var sgIds []*string
		for _, sgId := range sgs {
			sgIds = append(sgIds, helper.String(sgId.(string)))
		}
		err := redisService.ModifyDBInstanceSecurityGroups(ctx, "redis", d.Id(), sgIds)
		if err != nil {
			return err
		}
	}

	if d.HasChanges("type_id") {
		request := redis.NewUpgradeInstanceVersionRequest()
		typeId := d.Get("type_id").(int)
		request.InstanceId = &id
		request.TargetInstanceType = helper.String(strconv.Itoa(typeId))
		waitSwitch := 2
		if v, ok := d.GetOkExists("wait_switch"); ok {
			waitSwitch = v.(int)
		}

		request.SwitchOption = helper.IntInt64(waitSwitch)
		startTime := time.Now().Format("2006-01-02 15:04:05")
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().UpgradeInstanceVersion(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate redis upgradeVersionOperation failed, reason:%+v", logId, err)
			return err
		}

		service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		if waitSwitch == 2 {
			_, _, _, err = service.CheckRedisOnlineOk(ctx, id, 20*tccommon.ReadRetryTimeout)
			if err != nil {
				log.Printf("[CRITAL]%s redis upgradeVersionOperation fail, reason:%s\n", logId, err.Error())
				return err
			}
		} else {
			time.Sleep(10 * time.Second)
			paramMap := make(map[string]interface{})
			paramMap["InstanceId"] = &id
			paramMap["BeginTime"] = &startTime
			paramMap["TaskTypes"] = helper.StringsStringsPoint([]string{"043"})
			err := resource.Retry(5*tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := service.DescribeRedisInstanceTaskListByFilter(ctx, paramMap)
				if e != nil {
					return tccommon.RetryError(e)
				}
				if result == nil || len(result) < 1 {
					return resource.RetryableError(fmt.Errorf("redis upgradeVersion fail, result is nil"))
				}
				for _, v := range result {
					if *v.Result == 0 || *v.Result == 1 {
						return resource.RetryableError(fmt.Errorf("redis upgradeVersion state is %v, retry...", *v.Result))
					}
					if *v.Result == 4 {
						return resource.NonRetryableError(fmt.Errorf("redis upgradeVersion fail, task status is %v", *v.Result))
					}
					if *v.Result == 2 {
						return nil
					}
				}
				return resource.RetryableError(fmt.Errorf("redis upgradeVersion fail, retry..."))
			})
			if err != nil {
				log.Printf("[CRITAL]%s redis upgradeVersion failed, reason:%+v", logId, err)
				return err
			}
		}

	}

	if d.HasChange("vpc_id") || d.HasChange("subnet_id") || d.HasChange("port") || d.HasChange("recycle") || d.HasChange("ip") {
		if _, ok := d.GetOk("operation_network"); !ok {
			return fmt.Errorf("When modifying `vpc_id`, `subnet_id`, `port`, `recycle`, `ip`, the `operation_network` parameter is required")
		}

		request := redis.NewModifyNetworkConfigRequest()
		request.InstanceId = &id

		operation := d.Get("operation_network").(string)
		request.Operation = &operation

		switch operation {
		case REDIS_MODIFY_NETWORK_CONFIG[0]:
			if v, ok := d.GetOk("ip"); ok {
				request.Vip = helper.String(v.(string))
			} else {
				return fmt.Errorf("When `operation_network` is %v, this parameter must be filled in", operation)
			}

			if v, ok := d.GetOk("port"); ok {
				request.VPort = helper.IntInt64(v.(int))
			} else {
				return fmt.Errorf("When `operation_network` is %v, this parameter must be filled in", operation)
			}
		case REDIS_MODIFY_NETWORK_CONFIG[1], REDIS_MODIFY_NETWORK_CONFIG[2]:
			if v, ok := d.GetOk("vpc_id"); ok {
				request.VpcId = helper.String(v.(string))
			} else {
				return fmt.Errorf("When `operation_network` is %v, this parameter must be filled in", operation)
			}

			if v, ok := d.GetOk("subnet_id"); ok {
				request.SubnetId = helper.String(v.(string))
			} else {
				return fmt.Errorf("When `operation_network` is %v, this parameter must be filled in", operation)
			}
		case REDIS_MODIFY_NETWORK_CONFIG[3]:
			if v, ok := d.GetOk("port"); ok {
				request.VPort = helper.IntInt64(v.(int))
			} else {
				return fmt.Errorf("When `operation_network` is %v, this parameter must be filled in", operation)
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().ModifyNetworkConfig(request)
			if e != nil {
				if _, ok := e.(*sdkErrors.TencentCloudSDKError); !ok {
					return resource.RetryableError(e)
				} else {
					return resource.NonRetryableError(e)
				}
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate redis networkConfig failed, reason:%+v", logId, err)
			return err
		}

		service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		_, _, _, err = service.CheckRedisOnlineOk(ctx, id, 20*tccommon.ReadRetryTimeout)
		if err != nil {
			log.Printf("[CRITAL]%s redis networkConfig fail, reason:%s\n", logId, err.Error())
			return err
		}

		_ = d.Set("operation_network", operation)
	}

	if d.HasChange("wan_address_switch") {
		err := resourceRedisWanAddressModify(ctx, &redisService, meta, d.Id(), d.Get("wan_address_switch").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("replicas_read_only") {
		if v, ok := d.GetOkExists("replicas_read_only"); ok {
			var taskId int64
			if v.(bool) {
				// enable
				request := redis.NewEnableReplicaReadonlyRequest()
				request.InstanceId = &id
				err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().EnableReplicaReadonly(request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil || result.Response.TaskId == nil {
						return resource.RetryableError(fmt.Errorf("Enable replica read only fail, Response is nil."))
					}

					taskId = *result.Response.TaskId
					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s enable replica read only failed, reason:%+v", logId, err)
					return err
				}
			} else {
				// disable
				request := redis.NewDisableReplicaReadonlyRequest()
				request.InstanceId = &id
				err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().DisableReplicaReadonly(request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil || result.Response.TaskId == nil {
						return resource.RetryableError(fmt.Errorf("Disable replica read only fail, Response is nil."))
					}

					taskId = *result.Response.TaskId
					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s disable replica read only failed, reason:%+v", logId, err)
					return err
				}
			}

			// wait
			request := redis.NewDescribeTaskInfoRequest()
			request.TaskId = helper.Int64Uint64(taskId)
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().DescribeTaskInfo(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil || result.Response.Status == nil {
					return resource.RetryableError(fmt.Errorf("Describe task info fail, Response is nil."))
				}

				if *result.Response.Status == "succeed" {
					return nil
				} else if *result.Response.Status == "failed" || *result.Response.Status == "error" {
					return resource.NonRetryableError(fmt.Errorf("Replica read only fail, task status is %v", *result.Response.Status))
				} else {
					return resource.RetryableError(fmt.Errorf("Replica read only is still in progress...Status is %v", *result.Response.Status))
				}
			})

			if err != nil {
				log.Printf("[CRITAL]%s replica read only failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		//internal version: replace setTagUpdate begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
		resourceName := tccommon.BuildTagResourceName("redis", "instance", region, id)
		//internal version: replace setTagUpdate end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

		//internal version: replace waitTag begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
		//internal version: replace waitTag end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	}

	d.Partial(false)

	return resourceTencentCloudRedisInstanceRead(d, meta)
}

func resourceTencentCloudRedisInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	defer tccommon.LogElapsed("resource.tencentcloud_redis_instance.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	// Collect infos before deleting action
	var chargeType string
	has, _, info, err := service.CheckRedisOnlineOk(ctx, d.Id(), 20*tccommon.ReadRetryTimeout)

	if err != nil {
		log.Printf("[CRITAL]%s redis querying before deleting task fail, reason:%s\n", logId, err.Error())
		return err
	}

	if !has {
		return nil
	}

	chargeType = REDIS_CHARGE_TYPE_NAME[*info.BillingMode]

	var wait = func(action string, taskInfo interface{}) (errRet error) {

		errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			var ok bool
			var err error
			switch v := taskInfo.(type) {
			case int64:
				ok, err = service.DescribeTaskInfo(ctx, d.Id(), v)
			case string:
				ok, _, err = service.DescribeInstanceDealDetail(ctx, v)
			}
			if err != nil {
				if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
					return resource.RetryableError(err)
				} else {
					return resource.NonRetryableError(err)
				}
			}
			if ok {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("%s timeout.", action))
		})

		if errRet != nil {
			log.Printf("[CRITAL]%s redis %s fail, reason:%s\n", logId, action, errRet.Error())
		}
		return errRet
	}

	forceDelete := d.Get("force_delete").(bool)
	if chargeType == REDIS_CHARGE_TYPE_POSTPAID {
		taskId, err := service.DestroyPostpaidInstance(ctx, d.Id())
		if err != nil {
			log.Printf("[CRITAL]%s redis %s fail, reason:%s\n", logId, "DestroyPostpaidInstance", err.Error())
			return err
		}
		if err = wait("DestroyPostpaidInstance", taskId); err != nil {
			return err
		}

	} else {
		if _, err := service.DestroyPrepaidInstance(ctx, d.Id()); err != nil {
			log.Printf("[CRITAL]%s redis %s fail, reason:%s\n", logId, "DestroyPrepaidInstance", err.Error())
			return err
		}

		// Deal info only support create and renew and resize, need to check destroy status by describing api.
		if errDestroyChecking := resource.Retry(20*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			has, isolated, err := service.CheckRedisDestroyOk(ctx, d.Id())
			if err != nil {
				log.Printf("[CRITAL]%s CheckRedisDestroyOk fail, reason:%s\n", logId, err.Error())
				return resource.NonRetryableError(err)
			}
			if !has || isolated {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("instance is not ready to be destroyed"))
		}); errDestroyChecking != nil {
			log.Printf("[CRITAL]%s redis querying before deleting task fail, reason:%s\n", logId, errDestroyChecking.Error())
			return errDestroyChecking
		}
	}

	if forceDelete {
		taskId, err := service.CleanUpInstance(ctx, d.Id())
		if err != nil {
			log.Printf("[CRITAL]%s redis %s fail, reason:%s\n", logId, "CleanUpInstance", err.Error())
			return err
		}

		return wait("CleanUpInstance", taskId)
	}
	return nil
}

func checkIdsEqual(o []int, n []int) bool {
	if len(o) != len(n) {
		return false
	}

	sort.Ints(o)
	sort.Ints(n)

	for i, v := range o {
		if v != n[i] {
			return false
		}
	}
	return true
}

func resourceRedisNodeSetModify(ctx context.Context, service *RedisService, d *schema.ResourceData) error {
	id := d.Id()
	memSize := d.Get("mem_size").(int)
	shardNum := d.Get("redis_shard_num").(int)
	o, n := d.GetChange("replica_zone_ids")
	oz := helper.InterfacesIntegers(o.([]interface{}))
	nz := helper.InterfacesIntegers(n.([]interface{}))
	log.Printf("o = %v, n = %v", oz, nz)
	adds, lacks := tccommon.GetListDiffs(oz, nz)

	var redisNodeInfos []*redis.RedisNodeInfo

	if len(adds) > 0 {
		_, _, info, err := service.CheckRedisOnlineOk(ctx, id, tccommon.ReadRetryTimeout)
		if err != nil {
			return err
		}
		redisNodeInfos = info.NodeSet
		redisReplicaCount := len(redisNodeInfos) - 1

		log.Printf("%v will be add", adds)
		var addNodes []*redis.RedisNodeInfo
		for _, zoneId := range adds {
			addNodes = append(addNodes, &redis.RedisNodeInfo{
				NodeType: helper.IntInt64(1),
				ZoneId:   helper.IntUint64(zoneId),
			})
		}
		if redisReplicaCount+len(adds) == 0 && len(adds) == 1 {
			// Processing the change from a single-AZ instance to a multi-AZ instance
			request := redis.NewModifyInstanceAvailabilityZonesRequest()
			if v, ok := d.GetOkExists("wait_switch"); ok {
				request.SwitchOption = helper.IntInt64(v.(int))
			} else {
				request.SwitchOption = helper.IntInt64(2)
			}
			request.InstanceId = &id
			request.NodeSet = append(request.NodeSet,
				&redis.RedisNodeInfo{
					NodeType: helper.IntInt64(1),
					ZoneId:   helper.IntUint64(adds[0]),
				},
				&redis.RedisNodeInfo{
					NodeType: helper.IntInt64(0),
					ZoneId:   helper.IntUint64(int(*info.ZoneId)),
				})
			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				_, err := service.client.UseRedisClient().ModifyInstanceAvailabilityZones(request)
				if err != nil {
					return tccommon.RetryError(err, redis.INTERNALERROR)
				}
				return nil
			})
		} else {
			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				_, err := service.UpgradeInstance(ctx, d.Id(), memSize, shardNum, redisReplicaCount+len(adds), addNodes)
				if err != nil {
					return tccommon.RetryError(err, redis.FAILEDOPERATION_UNKNOWN)
				}
				return nil
			})
		}
		if err != nil {
			return err
		}
		err = service.CheckRedisUpdateOk(ctx, id)
		if err != nil {
			return err
		}
	}

	if len(lacks) > 0 {
		_, _, info, err := service.CheckRedisOnlineOk(ctx, id, tccommon.ReadRetryTimeout)
		if err != nil {
			return err
		}
		redisNodeInfos = info.NodeSet
		redisReplicaCount := len(redisNodeInfos) - 1
		removeNodes := TencentCloudRedisGetRemoveNodesByIds(lacks[:], redisNodeInfos)
		replicasParam := redisReplicaCount - len(lacks)
		if replicasParam <= 0 {
			return fmt.Errorf("cannot delete replica %d which is your only replica on instance %s", removeNodes[0].NodeId, id)
		}
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			_, err := service.UpgradeInstance(ctx, id, memSize, shardNum, replicasParam, removeNodes)
			if err != nil {
				return tccommon.RetryError(err, redis.FAILEDOPERATION_UNKNOWN)
			}
			return nil
		})
		if err != nil {
			return err
		}
		err = service.CheckRedisUpdateOk(ctx, id)
		if err != nil {
			return err
		}
	}

	// Non-Multi-AZ modification redis_replicas_num
	if d.HasChange("redis_replicas_num") && len(oz) == 0 && len(nz) == 0 {
		_, replica := d.GetChange("redis_replicas_num")
		redisReplicasNum := replica.(int)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			_, err := service.UpgradeInstance(ctx, id, memSize, shardNum, redisReplicasNum, nil)
			if err != nil {
				// Upgrade memory will cause instance lock and cannot acknowledge by polling status, wait until lock release
				return tccommon.RetryError(err, redis.FAILEDOPERATION_UNKNOWN, redis.FAILEDOPERATION_SYSTEMERROR)
			}
			return nil
		})
		if err != nil {
			return err
		}

		err = service.CheckRedisUpdateOk(ctx, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func TencentCloudRedisGetRemoveNodesByIds(ids []int, nodes []*redis.RedisNodeInfo) (result []*redis.RedisNodeInfo) {
	for i := range nodes {
		node := nodes[i]
		if *node.NodeType == 0 {
			continue
		}
		index := tccommon.FindIntListIndex(ids, int(*node.ZoneId))
		if index == -1 {
			continue
		}
		result = append(result, node)
		ids = append(ids[:index], ids[index+1:]...)
	}
	return
}

func resourceRedisWanAddressModify(ctx context.Context, service *RedisService, meta interface{}, instanceId, addressSwitch string) error {
	instance, err := service.DescribeRedisInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if addressSwitch == "close" {
		if instance.WanAddress != nil && *instance.WanAddress != "" {
			request := redis.NewReleaseWanAddressRequest()
			request.InstanceId = helper.String(instanceId)

			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().ReleaseWanAddressWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG] api[%s] success, request body [%s], response body [%s]\n", request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if reqErr != nil {
				log.Printf("[CRITAL] delete redis wan address failed, reason:%+v", reqErr)
				return reqErr
			}

			_, _, _, err := service.CheckRedisOnlineOk(ctx, instanceId, 20*tccommon.ReadRetryTimeout)
			if err != nil {
				log.Printf("[CRITAL] redis networkConfig fail, reason:%s\n", err.Error())
				return err
			}
		}
	} else if addressSwitch == "open" {
		if instance.WanAddress == nil || *instance.WanAddress == "" {
			request := redis.NewAllocateWanAddressRequest()
			request.InstanceId = helper.String(instanceId)

			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient().AllocateWanAddressWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG] api[%s] success, request body [%s], response body [%s]\n", request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if reqErr != nil {
				log.Printf("[CRITAL] create redis wan address failed, reason:%+v", reqErr)
				return reqErr
			}

			service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
			_, _, _, err := service.CheckRedisOnlineOk(ctx, instanceId, 20*tccommon.ReadRetryTimeout)
			if err != nil {
				log.Printf("[CRITAL] redis networkConfig fail, reason:%s\n", err.Error())
				return err
			}
		}
	} else {
		return fmt.Errorf("invalid address_switch %s", addressSwitch)
	}

	return nil
}
