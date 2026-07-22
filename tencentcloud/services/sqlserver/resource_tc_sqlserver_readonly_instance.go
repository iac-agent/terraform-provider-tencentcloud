package sqlserver

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcpostgresql "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/postgresql"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"context"
	"fmt"
	"log"

	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSqlserverReadonlyInstance() *schema.Resource {
	readonlyInstanceInfo := map[string]*schema.Schema{
		"master_instance_id": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
			Description: "Indicates master 实例 ID 的 recovery 实例.",
		},
		"readonly_group_type": {
			Type:         schema.TypeInt,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: tccommon.ValidateAllowedIntValue([]int{1, 3}),
			Description:  "Type 的 readonly 组. 有效 值: `1`, `3`. `1` 对于 一个 auto-assigned readonly 实例 per 一个 readonly 组, `2` 对于 creating new readonly 组, `3` 对于 all exist readonly 实例 stay 在 exist readonly 组. For now, 仅 `1` 和 `3` 是 支持.",
		},
		"force_upgrade": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Default:     false,
			Description: "Indicate 该 master 实例 upgrade 或 不. `true` 对于 upgrading master SQL Server 实例 到 集群 类型 通过 force. Default 是 false. 注意: 此 是 不 支持 使用 `DUAL`(ha_type), `2017`(engine_version) master SQL Server 实例, 对于 它 将 cause ha_type 的 master SQL Server 实例 change.",
		},
		"readonly_group_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "ID 的 readonly 组 该 此 实例 belongs 到. 当 `readonly_group_type` 集合 值 `3`, 它 必须 是 集合 使用 有效 值.",
		},
		"readonly_group_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Required 当 `readonly_group_type`=2, 名称 的 newly 创建 read-仅 组.",
		},
		"readonly_groups_is_offline_delay": {
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Description: "Required 当 `readonly_group_type`=2, whether newly 创建 read-仅 组 has delay elimination 已启用, 1-已启用, 0-已禁用. 当 delay between read-仅 copy 和 primary 实例 exceeds 阈值, 它 是 automatically removed.",
		},
		"readonly_groups_max_delay_time": {
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Description: "Required 当 `readonly_group_type`=2 和 `readonly_groups_is_offline_delay`=1, 阈值 对于 delayed elimination 的 newly 创建 read-仅 groups.",
		},
		"readonly_groups_min_in_group": {
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Description: "当 `readonly_group_type`=2 和 `readonly_groups_is_offline_delay`=1, 它 是 必填. After newly 创建 read-仅 组 是 delayed 和 removed, 在 least 数量 的 read-仅 copies should 是 retained.",
		},
		"engine_version": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Version 的 SQL Server 数据库 引擎.",
		},
		"ha_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "实例 类型.",
		},
		"project_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Project ID.",
		},
	}

	basic := TencentSqlServerBasicInfo(true)
	for k, v := range basic {
		readonlyInstanceInfo[k] = v
	}

	return &schema.Resource{
		Create: resourceTencentCloudSqlserverReadonlyInstanceCreate,
		Read:   resourceTencentCloudSqlserverReadonlyInstanceRead,
		Update: resourceTencentCloudSqlserverReadonlyInstanceUpdate,
		Delete: resourceTencentCloudSqlserverReadonlyInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: readonlyInstanceInfo,
	}
}

func resourceTencentCloudSqlserverReadonlyInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_readonly_instance.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	sqlserverService := SqlserverService{client: client}
	tagService := svctag.NewTagService(client)
	region := client.Region

	var (
		name                        = d.Get("name").(string)
		masterInstanceId            = d.Get("master_instance_id").(string)
		payType                     = d.Get("charge_type").(string)
		readonlyGroupType           = d.Get("readonly_group_type").(int)
		subnetId                    = d.Get("subnet_id").(string)
		vpcId                       = d.Get("vpc_id").(string)
		zone                        = d.Get("availability_zone").(string)
		storage                     = d.Get("storage").(int)
		memory                      = d.Get("memory").(int)
		readOnlyGroupIsOfflineDelay = d.Get("readonly_groups_is_offline_delay").(int)
		forceUpgrade                = d.Get("force_upgrade").(bool)
		readonlyGroupId             = ""
		securityGroups              = make([]string, 0)
	)

	if v, ok := d.GetOk("readonly_group_id"); ok && readonlyGroupType == 3 {
		readonlyGroupId = v.(string)
	}

	if temp, ok := d.GetOk("security_groups"); ok {
		sgGroup := temp.(*schema.Set).List()
		for _, sg := range sgGroup {
			securityGroups = append(securityGroups, sg.(string))
		}
	}

	request := sqlserver.NewCreateReadOnlyDBInstancesRequest()

	if v, ok := d.GetOk("readonly_group_name"); ok && readonlyGroupType == 2 {
		request.ReadOnlyGroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("readonly_groups_is_offline_delay"); ok && readonlyGroupType == 2 {
		request.ReadOnlyGroupIsOfflineDelay = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("readonly_groups_max_delay_time"); ok && readonlyGroupType == 2 && readOnlyGroupIsOfflineDelay == 1 {
		request.ReadOnlyGroupMaxDelayTime = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("readonly_groups_min_in_group"); ok && readonlyGroupType == 2 && readOnlyGroupIsOfflineDelay == 1 {
		request.ReadOnlyGroupMinInGroup = helper.IntInt64(v.(int))
	}

	request.InstanceId = &masterInstanceId
	request.InstanceChargeType = &payType
	request.Memory = helper.IntInt64(memory)
	request.Storage = helper.IntInt64(storage)
	request.SubnetId = &subnetId
	request.VpcId = &vpcId
	request.GoodsNum = helper.IntInt64(1)

	request.ReadOnlyGroupType = helper.IntInt64(readonlyGroupType)
	if readonlyGroupId != "" {
		request.ReadOnlyGroupId = &readonlyGroupId
	}

	if forceUpgrade {
		request.ReadOnlyGroupForcedUpgrade = helper.BoolToInt64Ptr(forceUpgrade)
	}
	request.Zone = &zone
	request.SecurityGroupList = make([]*string, 0, len(securityGroups))
	for _, v := range securityGroups {
		request.SecurityGroupList = append(request.SecurityGroupList, &v)
	}

	if payType == svcpostgresql.COMMON_PAYTYPE_POSTPAID {
		request.InstanceChargeType = helper.String("POSTPAID")
	}
	if payType == svcpostgresql.COMMON_PAYTYPE_PREPAID {
		request.InstanceChargeType = helper.String("PREPAID")

		if v, ok := d.Get("period").(int); ok {
			request.Period = helper.IntInt64(v)
		}
	}

	if v, ok := d.Get("auto_voucher").(int); ok {
		request.AutoVoucher = helper.IntInt64(v)
	}

	if v, ok := d.Get("voucher_ids").([]interface{}); ok {
		request.VoucherIds = helper.InterfacesStringsPoint(v)
	}

	if v, ok := d.GetOk("time_zone"); ok {
		request.TimeZone = helper.String(v.(string))
	}

	var instanceId string
	var outErr, inErr error

	outErr = resource.Retry(12*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		instanceId, inErr = sqlserverService.CreateSqlserverReadonlyInstance(ctx, request)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	d.SetId(instanceId)

	//set name
	outErr = resource.Retry(3*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		inErr := sqlserverService.ModifySqlserverInstanceName(ctx, instanceId, name)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		resourceName := tccommon.BuildTagResourceName("sqlserver", "instance", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	return resourceTencentCloudSqlserverReadonlyInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverReadonlyInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_readonly_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()
	sqlserverService := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instance, has, err := tencentSqlServerBasicInfoRead(ctx, d, meta)
	if err != nil {
		return err
	}

	if !has {
		d.SetId("")
		return nil
	}
	_ = d.Set("project_id", instance.ProjectId)
	_ = d.Set("engine_version", instance.Version)
	_ = d.Set("ha_type", SQLSERVER_HA_TYPE_FLAGS[*instance.HAFlag])

	//readonly group ID
	readOnlyInstance, err := sqlserverService.DescribeReadonlyGroupListByReadonlyInstanceId(ctx, instanceId)

	if err != nil {
		return err
	}

	if readOnlyInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `sqlserver_readonly_instance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("master_instance_id", readOnlyInstance.MasterInstanceId)
	_ = d.Set("readonly_group_id", readOnlyInstance.ReadOnlyGroupId)
	_ = d.Set("readonly_group_name", readOnlyInstance.ReadOnlyGroupName)
	_ = d.Set("readonly_groups_is_offline_delay", readOnlyInstance.IsOfflineDelay)
	_ = d.Set("readonly_groups_max_delay_time", readOnlyInstance.ReadOnlyMaxDelayTime)
	_ = d.Set("readonly_groups_min_in_group", readOnlyInstance.MinReadOnlyInGroup)

	readOnlyGroup, err := sqlserverService.DescribeReadOnlyGroupListById(ctx, *readOnlyInstance.MasterInstanceId, *readOnlyInstance.ReadOnlyGroupId)
	if readOnlyGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `sqlserver_readonly_instance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if readOnlyGroup.ReadOnlyGroupType != nil {
		_ = d.Set("readonly_group_type", readOnlyGroup.ReadOnlyGroupType)
	}

	if readOnlyGroup.ReadOnlyGroupForcedUpgrade != nil {
		if *readOnlyGroup.ReadOnlyGroupForcedUpgrade == 1 {
			_ = d.Set("force_upgrade", true)
		} else {
			_ = d.Set("force_upgrade", false)
		}
	}

	if instance.TimeZone != nil {
		_ = d.Set("time_zone", instance.TimeZone)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "sqlserver", "instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudSqlserverReadonlyInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_readonly_instance.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	immutableArgs := []string{"time_zone"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	//basic info update
	if err := sqlServerAllInstanceRoleUpdate(ctx, d, meta); err != nil {
		return err
	}

	return resourceTencentCloudSqlserverReadonlyInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverReadonlyInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_readonly_instance.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()
	sqlserverService := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var outErr, inErr error
	var has bool

	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, has, inErr = sqlserverService.DescribeSqlserverInstanceById(ctx, d.Id())
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}

	if !has {
		return nil
	}

	//terminate sql instance
	outErr = sqlserverService.TerminateSqlserverInstance(ctx, instanceId)

	if outErr != nil {
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = sqlserverService.TerminateSqlserverInstance(ctx, instanceId)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
	}

	if outErr != nil {
		return outErr
	}

	outErr = sqlserverService.DeleteSqlserverInstance(ctx, instanceId)

	if outErr != nil {
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = sqlserverService.DeleteSqlserverInstance(ctx, instanceId)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
	}

	if outErr != nil {
		return outErr
	}

	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, has, inErr := sqlserverService.DescribeSqlserverInstanceById(ctx, d.Id())
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		if has {
			inErr = fmt.Errorf("delete SQL Server readonly instance %s fail, instance still exists from SDK DescribeSqlserverInstanceById", instanceId)
			return resource.RetryableError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}
	return nil
}
