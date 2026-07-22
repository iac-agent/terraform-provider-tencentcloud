package sqlserver

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcpostgresql "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/postgresql"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func TencentSqlServerBasicInfo(isROInstance bool) map[string]*schema.Schema {
	basicSchema := map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
			Description:  "Name 的 SQL Server 实例.",
		},
		"charge_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      svcpostgresql.COMMON_PAYTYPE_POSTPAID,
			ForceNew:     true,
			ValidateFunc: tccommon.ValidateAllowedStringValue([]string{svcpostgresql.COMMON_PAYTYPE_PREPAID, svcpostgresql.COMMON_PAYTYPE_POSTPAID}),
			Description:  "Pay 类型 的 SQL Server 实例. Available 值 `PREPAID`, `POSTPAID_BY_HOUR`.",
		},
		"period": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: tccommon.ValidateIntegerInRange(1, 48),
			Description:  "Purchase 实例 周期 在 month. 值 does 不 exceed 48.",
		},
		"auto_voucher": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "Whether 到 使用 voucher automatically; 1 对于 yes, 0 对于 无, 默认值 是 0.",
		},
		"voucher_ids": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "An 数组 的 voucher IDs, currently 仅 一个 可以 是 使用 对于 单个 order.",
		},
		"vpc_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "ID 的 VPC.",
		},
		"subnet_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "ID 的 子网.",
		},
		"storage": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Disk 大小 (在 GB). Allowed 值 必须 是 多个 的 10. 存储 必须 是 集合 使用 限制 的 `storage_min` 和 `storage_max` 其中 数据 source `tencentcloud_sqlserver_specinfos` provides.",
		},
		"memory": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Memory 大小 (在 GB). Allowed 值 必须 是 larger 比 `内存` 该 数据 source `tencentcloud_sqlserver_specinfos` provides.",
		},
		"availability_zone": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Computed:    true,
			Description: "Availability zone.",
		},
		"security_groups": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "Security 组 bound 到 实例.",
		},
		"time_zone": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "System 时间 zone, 默认值: `China Standard Time`.",
		},
		//Computed values
		"ro_flag": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Readonly flag. `RO` (read-仅 实例), `MASTER` (primary 实例 使用 read-仅 实例). 如果 它 是 left 空, 它 refers 到 实例 其中 是 不 read-仅 和 has 无 RO 组.",
		},
		"vip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "IP 对于 私有 访问.",
		},
		"vport": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Port 对于 私有 访问.",
		},
		"create_time": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Create 时间 的 SQL Server 实例.",
		},
		"status": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Status 的 SQL Server 实例. 1 对于 applying, 2 对于 running, 3 对于 running 使用 限制, 4 对于 isolated, 5 对于 recycling, 6 对于 recycled, 7 对于 running 使用 任务, 8 对于 关闭-line, 9 对于 expanding, 10 对于 migrating, 11 对于 readonly, 12 对于 rebooting.",
		},
		"tags": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "tags 的 SQL Server.",
		},
		"wait_switch": {
			Type:        schema.TypeInt,
			Optional:    true,
			Deprecated:  "It has been deprecated from version 1.81.2.",
			Description: "way 到 execute allocation. Supported 值 include: 0 - execute immediately, 1 - execute 在 maintenance window.",
		},
		"dns_pod_domain": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Internet 地址 域名 名称.",
		},
		"tgw_wan_vport": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "External 端口 数量.",
		},
	}

	if !isROInstance {
		basicSchema["auto_renew"] = &schema.Schema{
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Automatic renewal sign. 0 对于 normal renewal, 1 对于 automatic renewal (Default). Only 有效 当 purchasing prepaid 实例.",
		}
	}

	return basicSchema
}

func ResourceTencentCloudSqlserverInstance() *schema.Resource {
	specialInfo := map[string]*schema.Schema{
		"multi_zones": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Computed:    true,
			Description: "Indicate whether 到 deploy across availability zones.",
		},
		//RO computed values
		"engine_version": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Default:     "2008R2",
			Description: "Version 的 SQL Server 数据库 引擎. Allowed 值 是 `2008R2`(SQL Server 2008 Enterprise), `2012SP3`(SQL Server 2012 Enterprise), `2016SP1` (SQL Server 2016 Enterprise), `201602`(SQL Server 2016 Standard) 和 `2017`(SQL Server 2017 Enterprise). Default 是 `2008R2`.",
		},
		"ha_type": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Default:     "DUAL",
			Description: "实例 类型. `DUAL` (dual-服务器 high availability), `CLUSTER` (集群). Default 是 `DUAL`.",
			Deprecated:  "It has been deprecated from version 1.81.136.",
		},
		"maintenance_week_set": {
			Type:        schema.TypeSet,
			Optional:    true,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeInt},
			Description: "A 列表 的 整数 indicates weekly maintenance. For 示例, [2,7] presents do weekly maintenance 在 every Tuesday 和 Sunday.",
		},
		"maintenance_start_time": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Start 时间 的 maintenance 在 一个 day, 格式 like `HH:mm`.",
		},
		"maintenance_time_span": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "timespan 的 maintenance 在 一个 day, 单位 是 hour.",
		},
		"project_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Project ID, 默认值 值 是 0.",
		},
	}
	basic := TencentSqlServerBasicInfo(false)
	for k, v := range basic {
		specialInfo[k] = v
	}
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverInstanceCreate,
		Read:   resourceTencentCloudSqlserverInstanceRead,
		Update: resourceTencentCloudSqlserverInstanceUpdate,
		Delete: resourceTencentCLoudSqlserverInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: helper.ImportWithDefaultValue(map[string]interface{}{
				"auto_voucher": 0,
			}),
		},
		Schema: specialInfo,
	}
}

func resourceTencentCloudSqlserverInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_instance.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	sqlserverService := SqlserverService{client: client}
	tagService := svctag.NewTagService(client)
	region := client.Region
	var (
		name           = d.Get("name").(string)
		dbVersion      = d.Get("engine_version").(string)
		payType        = d.Get("charge_type").(string)
		projectId      = d.Get("project_id").(int)
		subnetId       = d.Get("subnet_id").(string)
		vpcId          = d.Get("vpc_id").(string)
		zone           = d.Get("availability_zone").(string)
		storage        = d.Get("storage").(int)
		memory         = d.Get("memory").(int)
		weekSet        = make([]int, 0)
		startTime      = d.Get("maintenance_start_time").(string)
		timeSpan       = d.Get("maintenance_time_span").(int)
		securityGroups = make([]string, 0)
	)

	var multiZones bool
	if v, ok := d.GetOkExists("multi_zones"); ok {
		multiZones = v.(bool)
	}

	if v, ok := d.GetOk("maintenance_week_set"); ok {
		mWeekSet := v.(*schema.Set).List()
		for _, vv := range mWeekSet {
			weekSet = append(weekSet, vv.(int))
		}
	}

	var instanceId string
	var outErr, inErr error

	if temp, ok := d.GetOkExists("security_groups"); ok {
		sgGroup := temp.(*schema.Set).List()
		for _, sg := range sgGroup {
			securityGroups = append(securityGroups, sg.(string))
		}
	}

	request := sqlserver.NewCreateDBInstancesRequest()
	request.DBVersion = &dbVersion
	request.Memory = helper.IntInt64(memory)
	request.Storage = helper.IntInt64(storage)
	request.SubnetId = &subnetId
	request.VpcId = &vpcId
	request.MultiZones = &multiZones

	if payType == svcpostgresql.COMMON_PAYTYPE_POSTPAID {
		request.InstanceChargeType = helper.String("POSTPAID")
	}
	if payType == svcpostgresql.COMMON_PAYTYPE_PREPAID {
		request.InstanceChargeType = helper.String("PREPAID")
		if v, ok := d.Get("auto_renew").(int); ok {
			request.AutoRenewFlag = helper.IntInt64(v)
		}

		if v, ok := d.Get("period").(int); ok {
			request.Period = helper.IntInt64(v)
		}
	}

	if v, ok := d.Get("auto_voucher").(int); ok {
		request.AutoVoucher = helper.IntInt64(v)
	}

	if v, ok := d.GetOk("voucher_ids"); ok {
		voucherIds := v.(*schema.Set).List()
		request.VoucherIds = helper.InterfacesStringsPoint(voucherIds)
	}

	if projectId != 0 {
		request.ProjectId = helper.IntInt64(projectId)
	}

	if len(weekSet) > 0 {
		request.Weekly = make([]*int64, 0)
		for _, i := range weekSet {
			request.Weekly = append(request.Weekly, helper.IntInt64(i))
		}
	}
	if startTime != "" {
		request.StartTime = &startTime
	}
	if timeSpan != 0 {
		request.Span = helper.IntInt64(timeSpan)
	}

	request.SecurityGroupList = make([]*string, 0, len(securityGroups))
	for _, v := range securityGroups {
		request.SecurityGroupList = append(request.SecurityGroupList, &v)
	}

	if v, ok := d.GetOk("time_zone"); ok {
		request.TimeZone = helper.String(v.(string))
	}

	request.GoodsNum = helper.IntInt64(1)
	request.Zone = &zone

	outErr = resource.Retry(6*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		instanceId, inErr = sqlserverService.CreateSqlserverInstance(ctx, request)
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
	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
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
	return resourceTencentCloudSqlserverInstanceRead(d, meta)
}

func sqlServerAllInstanceRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	sqlserverService := SqlserverService{client: client}
	tagService := svctag.NewTagService(client)
	region := client.Region
	instanceId := d.Id()

	var outErr, inErr error

	//upgrade storage and memory size
	if d.HasChange("memory") || d.HasChange("storage") || d.HasChange("auto_voucher") || d.HasChange("voucher_ids") {
		memory := d.Get("memory").(int)
		storage := d.Get("storage").(int)
		autoVoucher := d.Get("auto_voucher").(int)
		voucherIds := d.Get("voucher_ids").(*schema.Set).List()
		outErr = sqlserverService.UpgradeSqlserverInstance(ctx, instanceId, memory, storage, autoVoucher, helper.InterfacesStringsPoint(voucherIds))

		if outErr != nil {
			return outErr
		}

	}

	//update name
	if d.HasChange("name") {
		name := d.Get("name").(string)
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = sqlserverService.ModifySqlserverInstanceName(ctx, instanceId, name)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}
	}

	if d.HasChange("security_groups") {
		o, n := d.GetChange("security_groups")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		oldSet := os.List()
		newSet := ns.List()

		for _, v := range oldSet {
			sgId := v.(string)
			outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				inErr := sqlserverService.RemoveSecurityGroup(ctx, instanceId, sgId)
				if inErr != nil {
					return tccommon.RetryError(inErr)
				}
				return nil
			})
			if outErr != nil {
				return outErr
			}
		}
		for _, v := range newSet {
			sgId := v.(string)
			outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				inErr := sqlserverService.AddSecurityGroup(ctx, instanceId, sgId)
				if inErr != nil {
					return tccommon.RetryError(inErr)
				}
				return nil
			})
			if outErr != nil {
				return outErr
			}
		}
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		resourceName := tccommon.BuildTagResourceName("sqlserver", "instance", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return nil
}

func sqlServerAllInstanceNetUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		request     = sqlserver.NewModifyDBInstanceNetworkRequest()
		flowRequest = sqlserver.NewDescribeFlowStatusRequest()
		flowId      int64
		instanceId  = d.Id()
	)

	if d.HasChange("vpc_id") || d.HasChange("subnet_id") {
		vpcId := d.Get("vpc_id").(string)
		subnetId := d.Get("subnet_id").(string)
		request.InstanceId = &instanceId
		request.NewVpcId = &vpcId
		request.NewSubnetId = &subnetId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().ModifyDBInstanceNetwork(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil {
				e = fmt.Errorf("sqlserver configInstanceNetwork not exists")
				return resource.NonRetryableError(e)
			}

			flowId = *result.Response.FlowId
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update sqlserver configInstanceNetwork failed, reason:%+v", logId, err)
			return err
		}

		flowRequest.FlowId = &flowId
		err = resource.Retry(10*tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().DescribeFlowStatus(flowRequest)
			if e != nil {
				return tccommon.RetryError(e)
			}

			if *result.Response.Status == SQLSERVER_TASK_SUCCESS {
				return nil
			} else if *result.Response.Status == SQLSERVER_TASK_RUNNING {
				return resource.RetryableError(fmt.Errorf("sqlserver configInstanceNetwork status is running"))
			} else if *result.Response.Status == int64(SQLSERVER_TASK_FAIL) {
				return resource.NonRetryableError(fmt.Errorf("sqlserver configInstanceNetwork status is fail"))
			} else {
				e = fmt.Errorf("sqlserver configInstanceNetwork status illegal")
				return resource.NonRetryableError(e)
			}
		})

		if err != nil {
			log.Printf("[CRITAL]%s create sqlserver configInstanceNetwork failed, reason:%+v", logId, err)
			return err
		}
	}
	return nil
}

func resourceTencentCloudSqlserverInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_instance.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	d.Partial(true)

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

	//update network
	if err := sqlServerAllInstanceNetUpdate(d, meta); err != nil {
		return err
	}

	var outErr, inErr error
	instanceId := d.Id()

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	sqlserverService := SqlserverService{client: client}
	tagService := svctag.NewTagService(client)
	region := client.Region
	//update project id
	if d.HasChange("project_id") {
		projectId := d.Get("project_id").(int)
		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = sqlserverService.ModifySqlserverInstanceProjectId(ctx, instanceId, projectId)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

	}

	if d.HasChange("maintenance_week_set") || d.HasChange("maintenance_start_time") || d.HasChange("maintenance_time_span") {
		weekSet := make([]int, 0)
		if v, ok := d.GetOk("maintenance_week_set"); ok {
			mWeekSet := v.(*schema.Set).List()
			for _, vv := range mWeekSet {
				weekSet = append(weekSet, vv.(int))
			}
		}
		startTime := d.Get("maintenance_start_time").(string)
		timeSpan := d.Get("maintenance_time_span").(int)

		outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = sqlserverService.ModifySqlserverInstanceMaintenanceSpan(ctx, instanceId, weekSet, startTime, timeSpan)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			return outErr
		}

	}
	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		resourceName := tccommon.BuildTagResourceName("sqlserver", "instance", region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

	}

	d.Partial(false)

	return resourceTencentCloudSqlserverInstanceRead(d, meta)
}
func tencentSqlServerBasicInfoRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (instance *sqlserver.DBInstance,
	has bool, errRet error) {

	if d.Id() == "" {
		return
	}
	instanceId := d.Id()
	sqlserverService := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var outErr, inErr error
	instance, has, outErr = sqlserverService.DescribeSqlserverInstanceById(ctx, d.Id())
	if outErr != nil {
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			instance, has, inErr = sqlserverService.DescribeSqlserverInstanceById(ctx, instanceId)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
	}
	if outErr != nil {
		errRet = outErr
	}

	if !has {
		return
	}
	//basic info
	_ = d.Set("project_id", instance.ProjectId)
	_ = d.Set("availability_zone", instance.Zone)
	_ = d.Set("vpc_id", instance.UniqVpcId)
	_ = d.Set("subnet_id", instance.UniqSubnetId)
	_ = d.Set("name", instance.Name)
	_ = d.Set("charge_type", helper.Int64ToStr(*instance.PayMode))

	if int(*instance.PayMode) == 1 {
		_ = d.Set("charge_type", svcpostgresql.COMMON_PAYTYPE_PREPAID)
		if _, ok := d.GetOk("auto_renew"); ok {
			_ = d.Set("auto_renew", instance.RenewFlag)
		}
	} else {
		_ = d.Set("charge_type", svcpostgresql.COMMON_PAYTYPE_POSTPAID)
	}

	if instance.DnsPodDomain != nil {
		_ = d.Set("dns_pod_domain", instance.DnsPodDomain)
	}

	if instance.TgwWanVPort != nil {
		_ = d.Set("tgw_wan_vport", instance.TgwWanVPort)
	}

	var securityGroup []string
	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		securityGroup, inErr = sqlserverService.DescribeInstanceSecurityGroups(ctx, instanceId)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})

	if outErr != nil {
		errRet = outErr
	}

	_ = d.Set("ro_flag", instance.ROFlag)
	_ = d.Set("create_time", instance.CreateTime)
	_ = d.Set("status", instance.Status)
	_ = d.Set("memory", instance.Memory)
	_ = d.Set("storage", instance.Storage)
	_ = d.Set("vip", instance.Vip)
	_ = d.Set("vport", instance.Vport)
	_ = d.Set("security_groups", securityGroup)

	if instance.TimeZone != nil {
		_ = d.Set("time_zone", instance.TimeZone)
	}
	return
}

func resourceTencentCloudSqlserverInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_instance.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var outErr, inErr error
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
	_ = d.Set("multi_zones", instance.IsDrZone)

	//maintanence
	weekSet, startTime, timeSpan, outErr := sqlserverService.DescribeMaintenanceSpan(ctx, instanceId)
	if outErr != nil {
		outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			weekSet, startTime, timeSpan, inErr = sqlserverService.DescribeMaintenanceSpan(ctx, instanceId)
			if inErr != nil {
				return tccommon.RetryError(inErr)
			}
			return nil
		})
	}
	if outErr != nil {
		return outErr
	}
	_ = d.Set("maintenance_week_set", weekSet)
	_ = d.Set("maintenance_start_time", startTime)
	_ = d.Set("maintenance_time_span", timeSpan)

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "sqlserver", "instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCLoudSqlserverInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_instance.delete")()

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

	// Wait for instance to be isolated (status = 4)
	outErr = resource.Retry(tccommon.ReadRetryTimeout*10, func() *resource.RetryError {
		// Query instance status using DescribeDBInstances API
		instance, inErr := sqlserverService.DescribeSqlserverRestartDBInstanceById(ctx, instanceId)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}

		// Check if instance exists
		if instance == nil {
			return resource.NonRetryableError(fmt.Errorf("instance %s not found", instanceId))
		}

		// Check instance status
		if instance.Status != nil {
			status := *instance.Status
			log.Printf("[DEBUG]%s instance %s current status: %d", logId, instanceId, status)

			if status == 4 {
				// Instance is isolated, ready to delete
				log.Printf("[INFO]%s instance %s is isolated (status=4), ready to delete", logId, instanceId)
				return nil
			}

			// Continue waiting for other statuses
			return resource.RetryableError(fmt.Errorf("waiting for instance %s to be isolated, current status: %d", instanceId, status))
		}

		return resource.RetryableError(fmt.Errorf("instance %s status is nil", instanceId))
	})

	if outErr != nil {
		log.Printf("[CRITAL]%s wait for instance %s isolation failed, reason: %+v", logId, instanceId, outErr)
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

	return nil
}
