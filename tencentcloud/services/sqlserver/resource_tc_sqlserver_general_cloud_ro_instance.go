package sqlserver

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSqlserverGeneralCloudRoInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverGeneralCloudRoInstanceCreate,
		Read:   resourceTencentCloudSqlserverGeneralCloudRoInstanceRead,
		Update: resourceTencentCloudSqlserverGeneralCloudRoInstanceUpdate,
		Delete: resourceTencentCloudSqlserverGeneralCloudRoInstanceDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(CreateDefaultTimeout * time.Second),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout * time.Second),
			Update: schema.DefaultTimeout(UpdateDefaultTimeout * time.Second),
			Delete: schema.DefaultTimeout(DeleteDefaultTimeout * time.Second),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Primary 实例 ID, 在 格式: mssql-3l3fgqn7.",
			},
			"zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 Availability Zone, similar 到 ap-guangzhou-1 (Guangzhou District 1); 实例 sales area 可以 是 获取 through interface DescribeZones.",
			},
			"read_only_group_type": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 3),
				Description:  "Read-仅 组 类型 选项, 1- Ship according 到 一个 实例 和 一个 read-仅 组, 2 - Ship after creating read-仅 组, all 实例 是 under 此 read-仅 组, 3 - All 实例 shipped 是 在 existing Some read-仅 groups below.",
			},
			"memory": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "实例 内存 大小, 在 GB.",
			},
			"storage": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "实例 磁盘 大小, 在 GB.",
			},
			"cpu": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Number 的 实例 cores.",
			},
			"machine_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "主机 磁盘 类型 的 purchased 实例, CLOUD_HSSD-enhanced SSD 云 磁盘 对于 virtual machines, CLOUD_TSSD-extremely fast SSD 云 磁盘 对于 virtual machines, CLOUD_BSSD-universal SSD 云 磁盘 对于 virtual machines.",
			},
			"read_only_group_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Required 当 ReadOnlyGroupType=3, existing read-仅 组 ID.",
			},
			"read_only_group_name": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Required 当 ReadOnlyGroupType=2, 名称 的 newly 创建 read-仅 组.",
			},
			"read_only_group_is_offline_delay": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Required 当 ReadOnlyGroupType=2, whether 到 启用 delayed elimination 函数 对于 newly 创建 read-仅 组, 1-在, 0-关闭. 当 delay between read-仅 副本 和 primary 实例 是 greater 比 阈值, 它 将 是 automatically removed.",
			},
			"read_only_group_max_delay_time": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Mandatory 当 ReadOnlyGroupType=2 和 ReadOnlyGroupIsOfflineDelay=1, 阈值 对于 delay culling 的 newly 创建 read-仅 groups.",
			},
			"read_only_group_min_in_group": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Required 当 ReadOnlyGroupType=2 和 ReadOnlyGroupIsOfflineDelay=1, newly 创建 read-仅 组 retains 在 least 数量 的 read-仅 replicas after delay elimination.",
			},
			"instance_charge_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Payment 模式, 值 支持 PREPAID (prepaid), POSTPAID (postpaid).",
			},
			"subnet_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "VPC 子网 ID, 在 form 的 子网-bdoe83fa; SubnetId 和 VpcId need 到 是 集合 在 same 时间 或 不 集合 在 same 时间.",
			},
			"vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "VPC 网络 ID, 在 form 的 vpc-dsp338hz; SubnetId 和 VpcId need 到 是 集合 在 same 时间 或 不 集合 在 same 时间.",
			},
			"period": {
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 48),
				Description:  "Purchase 实例 周期, 默认值 值 是 1, 其中 表示 一个 month. 值 不能 exceed 48.",
			},
			"security_group_list": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Security 组 列表, fill 在 安全 组 ID 在 form 的 sg-xxx.",
			},
			"resource_tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag 描述 列表.",
			},
			"collation": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "System character 集合 collation, 默认值: Chinese_PRC_CI_AS.",
			},
			"time_zone": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "System 时间 zone, 默认值: China Standard Time.",
			},
			"ro_instance_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Primary read 仅 实例 ID, 在 格式: mssqlro-lbljc5qd.",
			},
		},
	}
}

func resourceTencentCloudSqlserverGeneralCloudRoInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_general_cloud_ro_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service      = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request      = sqlserver.NewCreateCloudReadOnlyDBInstancesRequest()
		response     = sqlserver.NewCreateCloudReadOnlyDBInstancesResponse()
		timeout      = d.Timeout(schema.TimeoutCreate)
		instanceId   string
		roInstanceId string
		dealNames    string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("read_only_group_type"); ok {
		roGroupType := v.(int)
		if roGroupType == 1 {
			request.ReadOnlyGroupForcedUpgrade = helper.IntInt64(1)
		}
		if roGroupType == 2 {
			if v, ok := d.GetOk("read_only_group_name"); ok {
				request.ReadOnlyGroupName = helper.String(v.(string))
			}

			if v, ok := d.GetOkExists("read_only_group_is_offline_delay"); ok {
				readOnlyGroupIsOfflineDelay := v.(int)
				if readOnlyGroupIsOfflineDelay == 1 {
					if v, ok := d.GetOk("read_only_group_max_delay_time"); ok {
						request.ReadOnlyGroupMaxDelayTime = helper.IntInt64(v.(int))
					}

					if v, ok := d.GetOk("read_only_group_min_in_group"); ok {
						request.ReadOnlyGroupMinInGroup = helper.IntInt64(v.(int))
					}
				}
				request.ReadOnlyGroupIsOfflineDelay = helper.IntInt64(readOnlyGroupIsOfflineDelay)
			}

		} else if roGroupType == 3 {
			if v, ok := d.GetOk("read_only_group_id"); ok {
				request.ReadOnlyGroupId = helper.String(v.(string))
			}
			request.ReadOnlyGroupForcedUpgrade = helper.IntInt64(1)
		}
		request.ReadOnlyGroupType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("memory"); ok {
		request.Memory = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("storage"); ok {
		request.Storage = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("cpu"); ok {
		request.Cpu = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("machine_type"); ok {
		request.MachineType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("period"); ok {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("security_group_list"); ok {
		securityGroupListSet := v.(*schema.Set).List()
		for i := range securityGroupListSet {
			securityGroupList := securityGroupListSet[i].(string)
			request.SecurityGroupList = append(request.SecurityGroupList, &securityGroupList)
		}
	}

	if v, ok := d.GetOk("collation"); ok {
		request.Collation = helper.String(v.(string))
	}

	if v, ok := d.GetOk("time_zone"); ok {
		request.TimeZone = helper.String(v.(string))
	}

	request.GoodsNum = helper.IntInt64(1)

	err := resource.Retry(timeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().CreateCloudReadOnlyDBInstances(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver generalCloudRoInstance failed, reason:%+v", logId, err)
		return err
	}

	dealNames = *response.Response.DealNames[0]
	roInstanceId, err = service.GetInfoFromDeal(ctx, dealNames, timeout)
	if err != nil {
		return err
	}

	if tags := helper.GetTags(d, "resource_tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::sqlserver:%s:uin/:instance/%s", region, roInstanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	d.SetId(strings.Join([]string{instanceId, roInstanceId}, tccommon.FILED_SP))

	return resourceTencentCloudSqlserverGeneralCloudRoInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralCloudRoInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_general_cloud_ro_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	instanceId := idSplit[0]
	roInstanceId := idSplit[1]

	generalCloudRoInstance, err := service.DescribeSqlserverGeneralCloudRoInstanceById(ctx, roInstanceId)
	if err != nil {
		return err
	}

	if generalCloudRoInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverGeneralCloudRoInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if generalCloudRoInstance.InstanceId != nil {
		_ = d.Set("ro_instance_id", generalCloudRoInstance.InstanceId)
	}

	if generalCloudRoInstance.Zone != nil {
		_ = d.Set("zone", generalCloudRoInstance.Zone)
	}

	if generalCloudRoInstance.Memory != nil {
		_ = d.Set("memory", generalCloudRoInstance.Memory)
	}

	if generalCloudRoInstance.Storage != nil {
		_ = d.Set("storage", generalCloudRoInstance.Storage)
	}

	if generalCloudRoInstance.Cpu != nil {
		_ = d.Set("cpu", generalCloudRoInstance.Cpu)
	}

	if generalCloudRoInstance.Type != nil {
		_ = d.Set("machine_type", generalCloudRoInstance.Type)
	}

	if generalCloudRoInstance.PayMode != nil {
		if *generalCloudRoInstance.PayMode == 0 {
			_ = d.Set("instance_charge_type", SQLSERVER_TYPE_POSTPAID)
		} else {
			_ = d.Set("instance_charge_type", SQLSERVER_TYPE_PREPAID)
		}
	}

	if generalCloudRoInstance.UniqSubnetId != nil {
		_ = d.Set("subnet_id", generalCloudRoInstance.UniqSubnetId)
	}

	if generalCloudRoInstance.UniqVpcId != nil {
		_ = d.Set("vpc_id", generalCloudRoInstance.UniqVpcId)
	}

	if generalCloudRoInstance.Collation != nil {
		_ = d.Set("collation", generalCloudRoInstance.Collation)
	}

	if generalCloudRoInstance.TimeZone != nil {
		_ = d.Set("time_zone", generalCloudRoInstance.TimeZone)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "sqlserver", "instance", tcClient.Region, roInstanceId)
	if err != nil {
		return err
	}

	_ = d.Set("resource_tags", tags)

	securityGroupList, err := service.DescribeInstanceSecurityGroups(ctx, roInstanceId)
	if err != nil {
		return err
	}

	if securityGroupList == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlservereSecurityGroups` [%s] not found, please check if it has been deleted.", logId, d.Id())
		return nil
	}

	if securityGroupList != nil {
		_ = d.Set("security_group_list", securityGroupList)
	}

	roGroupList, err := service.DescribeReadonlyGroupList(ctx, instanceId)
	if err != nil {
		return err
	}

	if roGroupList == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlservereReadonlyGroup` [%s] not found, please check if it has been deleted.", logId, d.Id())
		return nil
	}

	for _, v := range roGroupList {
		readOnlyInstanceSet := v.ReadOnlyInstanceSet
		for _, readOnlyInstance := range readOnlyInstanceSet {
			if roInstanceId == *readOnlyInstance.InstanceId {
				if v.ReadOnlyGroupId != nil {
					_ = d.Set("read_only_group_id", v.ReadOnlyGroupId)
				}

				if v.ReadOnlyGroupName != nil {
					_ = d.Set("read_only_group_name", v.ReadOnlyGroupName)
				}

				if v.IsOfflineDelay != nil {
					_ = d.Set("read_only_group_is_offline_delay", v.IsOfflineDelay)
				}

				if v.ReadOnlyMaxDelayTime != nil {
					_ = d.Set("read_only_group_max_delay_time", v.ReadOnlyMaxDelayTime)
				}

				if v.MinReadOnlyInGroup != nil {
					_ = d.Set("read_only_group_min_in_group", v.MinReadOnlyInGroup)
				}
			}
		}
	}

	return nil
}

func resourceTencentCloudSqlserverGeneralCloudRoInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_general_cloud_ro_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		ctx              = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		client           = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		sqlserverService = SqlserverService{client: client}
		request          = sqlserver.NewUpgradeDBInstanceRequest()
		timeout          = d.Timeout(schema.TimeoutUpdate)
		dealId           string
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	roInstanceId := idSplit[1]

	immutableArgs := []string{"instance_id", "zone", "read_only_group_type", "machine_type", "read_only_group_forced_upgrade", "read_only_group_id", "read_only_group_name", "read_only_group_is_offline_delay", "read_only_group_max_delay_time", "read_only_group_min_in_group", "instance_charge_type", "subnet_id", "vpc_id", "period", "security_group_list", "auto_voucher", "voucher_ids", "resource_tags", "collation", "time_zone"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.InstanceId = &roInstanceId
	request.WaitSwitch = helper.IntInt64(0)

	if d.HasChange("memory") || d.HasChange("storage") || d.HasChange("cpu") {
		if v, ok := d.GetOkExists("memory"); ok {
			request.Memory = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("storage"); ok {
			request.Storage = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("cpu"); ok {
			request.Cpu = helper.IntInt64(v.(int))
		}

		err := resource.Retry(timeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().UpgradeDBInstance(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			dealId = *result.Response.DealName
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update sqlserver generalCloudRoInstance failed, reason:%+v", logId, err)
			return err
		}

		_, err = sqlserverService.GetInfoFromDeal(ctx, dealId, timeout)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudSqlserverGeneralCloudRoInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralCloudRoInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_general_cloud_ro_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	roInstanceId := idSplit[1]

	if err := service.TerminateSqlserverInstance(ctx, roInstanceId); err != nil {
		return err
	}

	if err := service.DeleteSqlserverGeneralCloudRoInstanceById(ctx, roInstanceId); err != nil {
		return err
	}

	return nil
}
