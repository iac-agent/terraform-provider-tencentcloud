package sqlserver

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSqlserverGeneralCloudInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverGeneralCloudInstanceCreate,
		Read:   resourceTencentCloudSqlserverGeneralCloudInstanceRead,
		Update: resourceTencentCloudSqlserverGeneralCloudInstanceUpdate,
		Delete: resourceTencentCloudSqlserverGeneralCloudInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "Name 的 SQL Server 实例.",
			},
			"zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 AZ, such 作为 ap-guangzhou-1 (Guangzhou Zone 1). Purchasable AZs 对于 实例 可以 是 获取 through DescribeZones API.",
			},
			"memory": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Memory, 单位: GB.",
			},
			"storage": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "实例 磁盘 存储, 单位: GB.",
			},
			"cpu": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Cpu, 单位: CORE.",
			},
			"machine_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "主机 磁盘 类型 的 purchased 实例, CLOUD_HSSD-enhanced SSD 云 磁盘 对于 virtual machines, CLOUD_TSSD-extremely fast SSD 云 磁盘 对于 virtual machines, CLOUD_BSSD-universal SSD 云 磁盘 对于 virtual machines.",
			},
			"instance_charge_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Payment 模式, 值 支持 PREPAID (prepaid), POSTPAID (postpaid).",
			},
			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "项目 ID.",
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
				Description:  "Purchase 实例 周期, 默认值 值 是 1, 其中 表示 一个 month. 值 不能 exceed 48. 有效 仅 当 'instance_charge_type' 参数 值 是 'PREPAID'.",
			},
			"db_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "sqlserver 版本, currently all 支持 versions 是: 2008R2 (SQL Server 2008 R2 Enterprise), 2012SP3 (SQL Server 2012 Enterprise), 201202 (SQL Server 2012 Standard), 2014SP2 (SQL Server 2014 Enterprise), 201402 (SQL Server 2014 Standard), 2016SP1 (SQL Server 2016 Enterprise), 201602 (SQL Server 2016 Standard), 2017 (SQL Server 2017 Enterprise), 201702 (SQL Server 2017 Standard), 2019 (SQL Server 2019 Enterprise), 201902 (SQL Server 2019 Standard). Each 地域 支持 different versions 对于 sale, 和 版本 信息 该 可以 是 sold 在 each 地域 可以 是 pulled through DescribeProductConfig interface. 如果 left blank, 默认值 版本 是 2008R2.",
			},
			"auto_renew_flag": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Automatic renewal flag: 0-normal renewal 1-automatic renewal, 默认值 是 1 automatic renewal. 有效 仅 当 purchasing prepaid 实例. 有效 仅 当 'instance_charge_type' 参数 值 是 'PREPAID'.",
			},
			"security_group_list": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Security 组 列表, fill 在 安全 组 ID 在 form 的 sg-xxx.",
			},
			"weekly": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Maintainable 时间 window 配置, 在 weeks, indicates days 的 week 该 allow maintenance, 1-7 represent Monday 到 weekend respectively.",
			},
			"start_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Maintainable 时间 window 配置, daily maintainable start 时间.",
			},
			"span": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Maintainable 时间 window 配置, 时长, 单位: hour.",
			},
			"resource_tags": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "A collection 的 tags bound 到 new 实例.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "标签 键.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "标签 值.",
						},
					},
				},
			},
			"collation": {
				Optional:    true,
				Type:        schema.TypeString,
				Default:     "Chinese_PRC_CI_AS",
				Description: "System character 集合 collation, 默认值: Chinese_PRC_CI_AS.",
			},
			"time_zone": {
				Optional:    true,
				Type:        schema.TypeString,
				Default:     "China Standard Time",
				Description: "System 时间 zone, 默认值: China Standard Time.",
			},
			"ha_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Deprecated:  "It has been deprecated from version 1.81.2.",
				Description: "Upgrade high-availability architecture 的 sqlserver, upgrade 从 mirror disaster recovery 到 always 在 集群 disaster recovery, 仅 support 2017 和 above 和 support always 在 high-availability 实例, do 不 support downgrading 到 mirror disaster recovery, CLUSTER-upgrade 到 always 在 容量 Disaster, 如果 不 filled, high-availability architecture 将 不 是 modified.",
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
			"multi_zones": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether 到 deploy across availability zones, 默认值 值 是 false.",
			},
			"multi_nodes": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether 它 是 multi-节点 architecture 实例, 默认值 值 是 false. 当 MultiNodes = true, 参数 MultiZones 必须 是 true.",
			},
			"dr_zones": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "standby 节点 availability area 是 空 通过 默认值. 当 MultiNodes = true, primary 节点 和 standby 节点 availability areas 不能 all 是 same. 最小 数量 的 standby availability areas 集合 是 2, 和 最大 数量 是 无 more 比 5.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"disk_encrypt_flag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Disk 加密 identification, 0-不 encrypted, 1-encrypted.",
			},
		},
	}
}

func resourceTencentCloudSqlserverGeneralCloudInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_general_cloud_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service      = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request      = sqlserver.NewCreateCloudDBInstancesRequest()
		instanceId   string
		instanceName string
		dealId       string
	)

	if v, ok := d.GetOk("name"); ok {
		instanceName = v.(string)
	}

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
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
		if v.(string) == SQLSERVER_TYPE_PREPAID {
			if v, ok := d.GetOk("period"); ok {
				request.Period = helper.IntInt64(v.(int))
			}

			if v, ok := d.GetOk("auto_renew_flag"); ok {
				request.AutoRenewFlag = helper.IntInt64(v.(int))
			}
		}
	}

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_version"); ok {
		request.DBVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_list"); ok {
		securityGroupListSet := v.(*schema.Set).List()
		for i := range securityGroupListSet {
			securityGroupList := securityGroupListSet[i].(string)
			request.SecurityGroupList = append(request.SecurityGroupList, &securityGroupList)
		}
	}

	if v, ok := d.GetOk("weekly"); ok {
		weeklySet := v.(*schema.Set).List()
		for i := range weeklySet {
			weekly := weeklySet[i].(int)
			request.Weekly = append(request.Weekly, helper.IntInt64(weekly))
		}
	}

	if v, ok := d.GetOk("start_time"); ok {
		request.StartTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("span"); ok {
		request.Span = helper.IntInt64(v.(int))
	}

	request.MultiZones = helper.Bool(true)

	if v, ok := d.GetOk("resource_tags"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			resourceTag := sqlserver.ResourceTag{}
			if v, ok := dMap["tag_key"]; ok {
				resourceTag.TagKey = helper.String(v.(string))
			}
			if v, ok := dMap["tag_value"]; ok {
				resourceTag.TagValue = helper.String(v.(string))
			}
			request.ResourceTags = append(request.ResourceTags, &resourceTag)
		}
	}

	if v, ok := d.GetOk("collation"); ok {
		request.Collation = helper.String(v.(string))
	}

	if v, ok := d.GetOk("time_zone"); ok {
		request.TimeZone = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("multi_zones"); ok {
		request.MultiZones = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("multi_nodes"); ok {
		request.MultiNodes = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOk("dr_zones"); ok {
		drZones := v.(*schema.Set).List()
		for i := range drZones {
			drZone := drZones[i].(string)
			request.DrZones = append(request.DrZones, &drZone)
		}
	}
	if v, ok := d.GetOkExists("disk_encrypt_flag"); ok {
		request.DiskEncryptFlag = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().CreateCloudDBInstances(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		dealId = *result.Response.DealName
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver generalCloudInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId, err = service.GetInfoFromDeal(ctx, dealId)
	if err != nil {
		return err
	}

	// set name
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		inErr := service.ModifySqlserverInstanceName(ctx, instanceId, instanceName)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}

		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverGeneralCloudInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralCloudInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_general_cloud_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	generalCloudInstance, err := service.DescribeSqlserverGeneralCloudInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if generalCloudInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverGeneralCloudInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if generalCloudInstance.Name != nil {
		_ = d.Set("name", generalCloudInstance.Name)
	}

	if generalCloudInstance.Zone != nil {
		_ = d.Set("zone", generalCloudInstance.Zone)
	}

	if generalCloudInstance.Memory != nil {
		_ = d.Set("memory", generalCloudInstance.Memory)
	}

	if generalCloudInstance.Storage != nil {
		_ = d.Set("storage", generalCloudInstance.Storage)
	}

	if generalCloudInstance.Cpu != nil {
		_ = d.Set("cpu", generalCloudInstance.Cpu)
	}

	if generalCloudInstance.Type != nil {
		_ = d.Set("machine_type", generalCloudInstance.Type)
	}

	if generalCloudInstance.PayMode != nil {
		if *generalCloudInstance.PayMode == 0 {
			_ = d.Set("instance_charge_type", SQLSERVER_TYPE_POSTPAID)
		} else {
			_ = d.Set("instance_charge_type", SQLSERVER_TYPE_PREPAID)
		}
	}

	if generalCloudInstance.ProjectId != nil {
		_ = d.Set("project_id", generalCloudInstance.ProjectId)
	}

	if generalCloudInstance.UniqSubnetId != nil {
		_ = d.Set("subnet_id", generalCloudInstance.UniqSubnetId)
	}

	if generalCloudInstance.UniqVpcId != nil {
		_ = d.Set("vpc_id", generalCloudInstance.UniqVpcId)
	}

	if generalCloudInstance.Version != nil {
		_ = d.Set("db_version", generalCloudInstance.Version)
	}

	if generalCloudInstance.RenewFlag != nil {
		_ = d.Set("auto_renew_flag", generalCloudInstance.RenewFlag)
	}

	if generalCloudInstance.ResourceTags != nil {
		resourceTagsList := []interface{}{}
		for _, resourceTags := range generalCloudInstance.ResourceTags {
			resourceTagsMap := map[string]interface{}{}

			if resourceTags.TagKey != nil {
				resourceTagsMap["tag_key"] = resourceTags.TagKey
			}

			if resourceTags.TagValue != nil {
				resourceTagsMap["tag_value"] = resourceTags.TagValue
			}

			resourceTagsList = append(resourceTagsList, resourceTagsMap)
		}

		_ = d.Set("resource_tags", resourceTagsList)

	}

	if generalCloudInstance.Collation != nil {
		_ = d.Set("collation", generalCloudInstance.Collation)
	}

	if generalCloudInstance.TimeZone != nil {
		_ = d.Set("time_zone", generalCloudInstance.TimeZone)
	}

	if generalCloudInstance.DnsPodDomain != nil {
		_ = d.Set("dns_pod_domain", generalCloudInstance.DnsPodDomain)
	}

	if generalCloudInstance.TgwWanVPort != nil {
		_ = d.Set("tgw_wan_vport", generalCloudInstance.TgwWanVPort)
	}

	if generalCloudInstance.IsDrZone != nil {
		_ = d.Set("multi_zones", generalCloudInstance.IsDrZone)
	}

	if len(generalCloudInstance.MultiSlaveZones) > 0 {
		_ = d.Set("multi_nodes", true)
		drZones := make([]string, 0)
		for _, multiSlaveZone := range generalCloudInstance.MultiSlaveZones {
			drZones = append(drZones, *multiSlaveZone.SlaveZone)
		}
		if len(drZones) > 0 {
			_ = d.Set("dr_zones", drZones)
		}
	} else {
		_ = d.Set("multi_nodes", false)
	}

	var insAttribute *sqlserver.DescribeDBInstancesAttributeResponseParams
	paramMap := map[string]interface{}{
		"InstanceId": helper.String(instanceId),
	}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverInsAttributeByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		insAttribute = result
		return nil
	})

	if err != nil {
		return err
	}

	if insAttribute.IsDiskEncryptFlag != nil {
		_ = d.Set("disk_encrypt_flag", insAttribute.IsDiskEncryptFlag)
	}
	maintenanceSpan, err := service.DescribeMaintenanceSpanById(ctx, instanceId)
	if err != nil {
		return err
	}

	if maintenanceSpan == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlservereMaintenanceSpan` [%s] not found, please check if it has been deleted.", logId, d.Id())
		return nil
	}

	if maintenanceSpan.Span != nil {
		_ = d.Set("span", maintenanceSpan.Span)
	}

	if maintenanceSpan.StartTime != nil {
		_ = d.Set("start_time", maintenanceSpan.StartTime)
	}

	if maintenanceSpan.Weekly != nil {
		_ = d.Set("weekly", maintenanceSpan.Weekly)
	}

	securityGroupList, err := service.DescribeInstanceSecurityGroups(ctx, instanceId)
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

	return nil
}

func resourceTencentCloudSqlserverGeneralCloudInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_general_cloud_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		ctx              = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		client           = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		sqlserverService = SqlserverService{client: client}
		request          = sqlserver.NewUpgradeDBInstanceRequest()
		instanceId       = d.Id()
		waitSwitch       int64
		dealId           string
		instanceName     string
	)

	request.InstanceId = &instanceId
	immutableArgs := []string{"zone", "machine_type", "instance_charge_type", "project_id", "subnet_id", "vpc_id", "period", "security_group_list", "weekly", "start_time", "span", "resource_tags", "collation", "time_zone", "multi_zones", "multi_nodes", "dr_zones", "disk_encrypt_flag"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			instanceName = v.(string)

			// set name
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				inErr := sqlserverService.ModifySqlserverInstanceName(ctx, instanceId, instanceName)
				if inErr != nil {
					return tccommon.RetryError(inErr)
				}

				return nil
			})

			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("memory") {
		if v, ok := d.GetOk("memory"); ok {
			request.Memory = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("storage") {
		if v, ok := d.GetOk("storage"); ok {
			request.Storage = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("cpu") {
		if v, ok := d.GetOk("cpu"); ok {
			request.Cpu = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("db_version") {
		if v, ok := d.GetOk("db_version"); ok {
			request.DBVersion = helper.String(v.(string))
		}
	}

	waitSwitch = 0
	request.WaitSwitch = &waitSwitch

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
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
		log.Printf("[CRITAL]%s update sqlserver generalCloudInstance failed, reason:%+v", logId, err)
		return err
	}

	_, err = sqlserverService.GetInfoFromDeal(ctx, dealId)
	if err != nil {
		return err
	}

	return resourceTencentCloudSqlserverGeneralCloudInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverGeneralCloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_general_cloud_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	if err := service.TerminateSqlserverInstanceById(ctx, instanceId); err != nil {
		return err
	}

	// Wait for instance to be isolated (status = 4)
	err := resource.Retry(tccommon.ReadRetryTimeout*10, func() *resource.RetryError {
		// Query instance status using DescribeDBInstances API
		instance, err := service.DescribeSqlserverRestartDBInstanceById(ctx, instanceId)
		if err != nil {
			return tccommon.RetryError(err)
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

	if err != nil {
		log.Printf("[CRITAL]%s wait for instance %s isolation failed, reason: %+v", logId, instanceId, err)
		return err
	}

	if err := service.DeleteSqlserverInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
