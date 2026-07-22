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

func ResourceTencentCloudSqlserverBasicInstance() *schema.Resource {

	return &schema.Resource{
		Create: resourceTencentCloudSqlserverBasicInstanceCreate,
		Read:   resourceTencentCloudSqlserverBasicInstanceRead,
		Update: resourceTencentCloudSqlserverBasicInstanceUpdate,
		Delete: resourceTencentCLoudSqlserverBasicInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "Name 的 SQL Server basic 实例.",
			},
			"cpu": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "CPU 数量 的 SQL Server basic 实例.",
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
			"machine_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{CLOUD_PREMIUM, CLOUD_SSD, CLOUD_HSSD, CLOUD_BSSD}),
				Description:  "主机 类型 的 purchased 实例, `CLOUD_PREMIUM` 对于 virtual machine high-performance 云 磁盘, `CLOUD_SSD` 对于 virtual machine SSD 云 磁盘, `CLOUD_HSSD` 对于 virtual machine enhanced 云 磁盘, `CLOUD_BSSD` 对于 virtual machine general purpose SSD 云 磁盘.",
			},
			"charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      svcpostgresql.COMMON_PAYTYPE_POSTPAID,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{svcpostgresql.COMMON_PAYTYPE_PREPAID, svcpostgresql.COMMON_PAYTYPE_POSTPAID}),
				Description:  "Pay 类型 的 SQL Server basic 实例. For now, 仅 `POSTPAID_BY_HOUR` 是 有效.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "ID 的 VPC.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "ID 的 子网.",
			},
			"engine_version": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Default:     "2008R2",
				Description: "Version 的 SQL Server basic 数据库 引擎. Allowed 值 是 `2008R2`(SQL Server 2008 Enterprise), `2012SP3`(SQL Server 2012 Enterprise), `2016SP1` (SQL Server 2016 Enterprise), `201602`(SQL Server 2016 Standard) 和 `2017`(SQL Server 2017 Enterprise). Default 是 `2008R2`.",
			},
			"time_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "System timezone 对于 SQL Server 实例. Default 是 `China Standard Time`. 此 setting 不能 是 changed after creation.",
			},
			"disk_encrypt_flag": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(0, 1),
				Description:  "Disk 加密 flag. `0` - Disabled (默认值), `1` - Enabled. Disk 加密 不能 是 changed after 实例 creation.",
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 48),
				Description:  "Purchase 实例 周期, 默认值 值 是 1, 其中 表示 一个 month. 值 does 不 exceed 48.",
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security 组 bound 到 实例.",
			},
			"auto_renew": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Automatic renewal sign. 0 对于 normal renewal, 1 对于 automatic renewal, 默认值 是 1 automatic renewal. Only 有效 当 purchasing prepaid 实例.",
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
			"maintenance_week_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "A 列表 的 整数 indicates weekly maintenance. For 示例, [1,7] presents do weekly maintenance 在 every Monday 和 Sunday.",
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
			"availability_zone": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "Availability zone.",
			},
			"collation": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Chinese_PRC_CI_AS",
				Description: "System character 集合 sorting 规则, 默认值: Chinese_PRC_CI_AS.",
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
				Description: "Create 时间 的 SQL Server basic 实例.",
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Status 的 SQL Server basic 实例. 1 对于 applying, 2 对于 running, 3 对于 running 使用 限制, 4 对于 isolated, 5 对于 recycling, 6 对于 recycled, 7 对于 running 使用 任务, 8 对于 关闭-line, 9 对于 expanding, 10 对于 migrating, 11 对于 readonly, 12 对于 rebooting.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "tags 的 SQL Server basic 实例.",
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
		},
	}
}

func resourceTencentCloudSqlserverBasicInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_basic_instance.create")()

	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		ctx              = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		client           = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		sqlserverService = SqlserverService{client: client}
		tagService       = svctag.NewTagService(client)
		region           = client.Region
		paramMap         = make(map[string]interface{})
		name             = d.Get("name").(string)
		payType          = d.Get("charge_type").(string)
		securityGroups   = make([]string, 0)
		voucherIds       = make([]string, 0)
		weekSet          = make([]int, 0)
	)
	if payType == svcpostgresql.COMMON_PAYTYPE_POSTPAID {
		payType = "POSTPAID"
		paramMap["autoRenew"] = 0
	} else {
		if v, ok := d.GetOk("auto_renew"); ok {
			paramMap["autoRenew"] = v.(int)
		} else {
			paramMap["autoRenew"] = 1
		}
	}
	paramMap["cpu"] = d.Get("cpu").(int)
	paramMap["memory"] = d.Get("memory").(int)
	paramMap["storage"] = d.Get("storage").(int)
	paramMap["subnetId"] = d.Get("subnet_id").(string)
	paramMap["vpcId"] = d.Get("vpc_id").(string)
	paramMap["machineType"] = d.Get("machine_type").(string)
	paramMap["payType"] = payType
	paramMap["engineVersion"] = d.Get("engine_version").(string)
	paramMap["period"] = d.Get("period").(int)
	paramMap["autoVoucher"] = d.Get("auto_voucher").(int)
	paramMap["availabilityZone"] = d.Get("availability_zone").(string)
	paramMap["collation"] = d.Get("collation").(string)

	// time_zone
	if v, ok := d.GetOk("time_zone"); ok {
		paramMap["time_zone"] = v.(string)
	}

	// disk_encrypt_flag
	if v, ok := d.GetOkExists("disk_encrypt_flag"); ok {
		paramMap["disk_encrypt_flag"] = v.(int)
	}

	if v, ok := d.GetOk("project_id"); ok {
		paramMap["projectId"] = v.(int)
	}
	if v, ok := d.GetOk("maintenance_start_time"); ok {
		paramMap["startTime"] = v.(string)
	}
	if v, ok := d.GetOk("maintenance_time_span"); ok {
		paramMap["timeSpan"] = v.(int)
	}
	// weekSet
	if v, ok := d.GetOk("maintenance_week_set"); ok {
		mWeekSet := v.(*schema.Set).List()
		for _, vv := range mWeekSet {
			weekSet = append(weekSet, vv.(int))
		}
		paramMap["weekSet"] = weekSet
	}
	// securityGroups
	if temp, ok := d.GetOk("security_groups"); ok {
		sgGroup := temp.(*schema.Set).List()
		for _, sg := range sgGroup {
			securityGroups = append(securityGroups, sg.(string))
		}
		paramMap["securityGroups"] = securityGroups
	}
	// voucherIds
	if temp, ok := d.GetOk("voucher_ids"); ok {
		voucherId := temp.(*schema.Set).List()
		for _, id := range voucherId {
			voucherIds = append(voucherIds, id.(string))
		}
		paramMap["voucherIds"] = voucherIds
	}

	var instanceId string
	var outErr, inErr error
	outErr = resource.Retry(12*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		instanceId, inErr = sqlserverService.CreateSqlserverBasicInstance(ctx, paramMap, weekSet, voucherIds, securityGroups)
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
	return resourceTencentCloudSqlserverBasicInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverBasicInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_basic_instance.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var outErr, inErr error
	instanceId := d.Id()
	sqlserverService := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	instance, has, outErr := sqlserverService.DescribeSqlserverInstanceById(ctx, d.Id())
	if outErr != nil {
		return outErr
	}
	if !has {
		d.SetId("")
		return nil
	}
	chargeType := instance.PayMode
	_ = d.Set("cpu", instance.Cpu)
	_ = d.Set("subnet_id", instance.UniqSubnetId)
	_ = d.Set("vpc_id", instance.UniqVpcId)
	_ = d.Set("machine_type", instance.Type)
	if int(*chargeType) == 1 {
		_ = d.Set("charge_type", svcpostgresql.COMMON_PAYTYPE_PREPAID)
		_ = d.Set("auto_renew", instance.RenewFlag)
	} else {
		_ = d.Set("charge_type", svcpostgresql.COMMON_PAYTYPE_POSTPAID)
		_ = d.Set("auto_renew", 0)
	}
	_ = d.Set("name", instance.Name)
	_ = d.Set("engine_version", instance.Version)

	_ = d.Set("availability_zone", instance.Zone)
	_ = d.Set("project_id", instance.ProjectId)
	_ = d.Set("create_time", instance.CreateTime)
	_ = d.Set("status", instance.Status)
	_ = d.Set("cpu", instance.Cpu)
	_ = d.Set("memory", instance.Memory)
	_ = d.Set("storage", instance.Storage)
	_ = d.Set("vip", instance.Vip)
	_ = d.Set("vport", instance.Vport)
	if instance.DnsPodDomain != nil {
		_ = d.Set("dns_pod_domain", instance.DnsPodDomain)
	}

	if instance.TgwWanVPort != nil {
		_ = d.Set("tgw_wan_vport", instance.TgwWanVPort)
	}

	// time_zone
	if instance.TimeZone != nil {
		_ = d.Set("time_zone", instance.TimeZone)
	}

	// Get disk encryption flag from attributes API
	var attribute *sqlserver.DescribeDBInstancesAttributeResponseParams
	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		attribute, inErr = sqlserverService.DescribeSqlserverInstanceAttributeById(ctx, instanceId)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		log.Printf("[WARN]%s describe sqlserver instance attribute failed, reason: %v", logId, outErr)
		// Don't fail the entire read, just log the warning
	}

	// disk_encrypt_flag
	if attribute != nil && attribute.IsDiskEncryptFlag != nil {
		_ = d.Set("disk_encrypt_flag", int(*attribute.IsDiskEncryptFlag))
	}

	//maintanence
	var weekSet []int
	var startTime string
	var timeSpan int
	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		weekSet, startTime, timeSpan, inErr = sqlserverService.DescribeMaintenanceSpan(ctx, instanceId)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}
	_ = d.Set("maintenance_week_set", weekSet)
	_ = d.Set("maintenance_start_time", startTime)
	_ = d.Set("maintenance_time_span", timeSpan)

	var securityGroup []string
	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		securityGroup, inErr = sqlserverService.DescribeInstanceSecurityGroups(ctx, instanceId)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}
	_ = d.Set("security_groups", securityGroup)

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "sqlserver", "instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)
	return nil
}

func resourceTencentCloudSqlserverBasicInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_basic_instance.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	sqlserverService := SqlserverService{client: client}
	tagService := svctag.NewTagService(client)
	region := client.Region
	payType := d.Get("charge_type").(string)

	immutableArgs := []string{"collation", "time_zone", "disk_encrypt_flag"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	var outErr, inErr error
	instanceId := d.Id()
	d.Partial(true)
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
	//upgrade storage and memory size
	if d.HasChange("memory") || d.HasChange("storage") ||
		d.HasChange("cpu") || d.HasChange("auto_voucher") {
		voucherIds := make([]string, 0)
		memory := d.Get("memory").(int)
		storage := d.Get("storage").(int)
		cpu := d.Get("cpu").(int)
		autoVoucher := d.Get("auto_voucher").(int)
		if temp, ok := d.GetOk("voucher_ids"); ok {
			voucherId := temp.(*schema.Set).List()
			for _, id := range voucherId {
				voucherIds = append(voucherIds, id.(string))
			}
		}
		outErr = resource.Retry(12*tccommon.WriteRetryTimeout, func() *resource.RetryError {
			inErr = sqlserverService.UpgradeSqlserverBasicInstance(ctx, instanceId, memory, storage, cpu, autoVoucher, voucherIds)
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

	if payType == svcpostgresql.COMMON_PAYTYPE_PREPAID {
		if d.HasChange("auto_renew") {
			var renewFlag int
			_, newValue := d.GetChange("auto_renew")
			renewFlag = newValue.(int)
			outErr = resource.Retry(2*tccommon.WriteRetryTimeout, func() *resource.RetryError {
				inErr = sqlserverService.NewModifyDBInstanceRenewFlag(ctx, instanceId, renewFlag)
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
	d.Partial(false)
	return resourceTencentCloudSqlserverBasicInstanceRead(d, meta)
}

func resourceTencentCLoudSqlserverBasicInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_basic_instance.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Id()
	sqlserverService := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var outErr, inErr error
	var has bool
	var instance *sqlserver.DBInstance

	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, has, inErr = sqlserverService.DescribeSqlserverInstanceById(ctx, d.Id())
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
	// PREPAID
	if *instance.PayMode == 1 {
		return fmt.Errorf("PREPAID instances are not allowed to be deleted now, please terminate them on console")
	}
	//terminate sql instance
	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		inErr = sqlserverService.TerminateSqlserverInstance(ctx, instanceId)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}

	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		inErr = sqlserverService.DeleteSqlserverInstance(ctx, instanceId)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}

	outErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		_, has, inErr := sqlserverService.DescribeSqlserverInstanceById(ctx, d.Id())
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		if has {
			inErr = fmt.Errorf("delete SQL Server basic instance %s fail, instance still exists from SDK DescribeSqlserverInstanceById", instanceId)
			return resource.RetryableError(inErr)
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}
	return nil
}
