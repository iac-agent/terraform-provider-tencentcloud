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

func ResourceTencentCloudSqlserverBusinessIntelligenceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverBusinessIntelligenceInstanceCreate,
		Read:   resourceTencentCloudSqlserverBusinessIntelligenceInstanceRead,
		Update: resourceTencentCloudSqlserverBusinessIntelligenceInstanceUpdate,
		Delete: resourceTencentCloudSqlserverBusinessIntelligenceInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 AZ, such 作为 ap-guangzhou-1 (Guangzhou Zone 1). Purchasable AZs 对于 实例 可以 是 获取 through theDescribeZones API.",
			},
			"project_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},
			"memory": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "实例 内存 大小 在 GB.",
			},
			"storage": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "实例 磁盘 大小 在 GB.",
			},
			"subnet_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "VPC 子网 ID 在 格式 的 子网-bdoe83fa. Both SubnetId 和 VpcId need 到 是 集合 或 unset 在 same 时间.",
			},
			"vpc_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "VPC ID 在 格式 的 vpc-dsp338hz. Both SubnetId 和 VpcId need 到 是 集合 或 unset 在 same 时间.",
			},
			"db_version": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Supported versions 的 business intelligence 服务器. 有效 值: 201603 (SQL Server 2016 Integration Services), 201703 (SQL Server 2017 Integration Services), 201903 (SQL Server 2019 Integration Services). Default 值: 201903. As purchasable versions 是 地域-特定, 您 可以 使用 DescribeProductConfig API 到 查询 信息 的 purchasable versions 在 each 地域.",
			},
			"security_group_list": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Security 组 列表, 其中 contains 安全 组 IDs 在 格式 的 sg-xxx.",
			},
			"weekly": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Configuration 的 maintenance window, 其中 specifies day 的 week 当 maintenance 可以 是 performed. 有效 值: 1 (Monday), 2 (Tuesday), 3 (Wednesday), 4 (Thursday), 5 (Friday), 6 (Saturday), 7 (Sunday).",
			},
			"start_time": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Configuration 的 maintenance window, 其中 specifies start 时间 的 daily maintenance.",
			},
			"span": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Configuration 的 maintenance window, 其中 specifies maintenance 时长 在 hours.",
			},
			"resource_tags": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Tags associated 使用 实例 到 是 创建.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag 键.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Tag 值.",
						},
					},
				},
			},
			"cpu": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "数量 的 CPU cores 的 实例 您 want 到 purchase.",
			},
			"machine_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "主机 类型 的 purchased 实例. 有效 值: CLOUD_PREMIUM (virtual machine 使用 premium 云 磁盘), CLOUD_SSD (virtual machine 使用 SSD).",
			},
			"instance_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 Name.",
			},
		},
	}
}

func resourceTencentCloudSqlserverBusinessIntelligenceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_business_intelligence_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		CreateDBIRequest = sqlserver.NewCreateBusinessDBInstancesRequest()
		DescRequest      = sqlserver.NewDescribeDBInstancesRequest()
		ModifyRequest    = sqlserver.NewModifyDBInstanceNameRequest()
		DescResponse     = sqlserver.DBInstance{}
		instanceId       string
		instanceName     string
	)

	if v, ok := d.GetOk("zone"); ok {
		CreateDBIRequest.Zone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		CreateDBIRequest.ProjectId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("memory"); ok {
		CreateDBIRequest.Memory = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("storage"); ok {
		CreateDBIRequest.Storage = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		CreateDBIRequest.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		CreateDBIRequest.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_version"); ok {
		CreateDBIRequest.DBVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_list"); ok {
		CreateDBIRequest.SecurityGroupList = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOk("weekly"); ok {
		CreateDBIRequest.Weekly = helper.InterfacesIntInt64Point(v.([]interface{}))

	}

	if v, ok := d.GetOk("start_time"); ok {
		CreateDBIRequest.StartTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("span"); ok {
		CreateDBIRequest.Span = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("resource_tags"); ok {
		for _, item := range v.([]interface{}) {
			if item != nil {
				dMap := item.(map[string]interface{})
				resourceTag := sqlserver.ResourceTag{}
				if t, h := dMap["tag_key"]; h {
					resourceTag.TagKey = helper.String(t.(string))
				}
				if t, h := dMap["tag_value"]; h {
					resourceTag.TagValue = helper.String(t.(string))
				}
				CreateDBIRequest.ResourceTags = append(CreateDBIRequest.ResourceTags, &resourceTag)
			}
		}
	}

	if v, ok := d.GetOk("cpu"); ok {
		CreateDBIRequest.Cpu = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("machine_type"); ok {
		CreateDBIRequest.MachineType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		instanceName = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().CreateBusinessDBInstances(CreateDBIRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]  ", logId, CreateDBIRequest.GetAction(), CreateDBIRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("sqlserver CreateBusinessDBInstances not exists")
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sqlserver businessIntelligenceInstance failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(5*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().DescribeDBInstances(DescRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]  ", logId, DescRequest.GetAction(), DescRequest.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("sqlserver business DescribeDBInstances not exists")
			return resource.NonRetryableError(e)
		}

		if *result.Response.TotalCount == 0 {
			e = fmt.Errorf("sqlserver business DescribeDBInstances not exists")
			return resource.NonRetryableError(e)
		}

		dbInstance := *result.Response.DBInstances[0]
		DescResponse = dbInstance
		if *dbInstance.Status == SQLSERVER_BSDBINSTANCE_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("sqlserver business DescribeDBInstances status is running"))
		} else if *dbInstance.Status == SQLSERVER_BSDBINSTANCE_STATUS_SUCCESS {
			return nil
		} else {
			e = fmt.Errorf("sqlserver business DescribeDBInstances status illegal")
			return resource.NonRetryableError(e)
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s describe sqlserver businessIntelligenceInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *DescResponse.InstanceId

	ModifyRequest.InstanceId = &instanceId
	ModifyRequest.InstanceName = &instanceName
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().ModifyDBInstanceName(ModifyRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]  ", logId, ModifyRequest.GetAction(), ModifyRequest.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver businessIntelligenceInstance failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverBusinessIntelligenceInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverBusinessIntelligenceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_business_intelligence_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	businessIntelligenceInstance, err := service.DescribeSqlserverBusinessIntelligenceInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if businessIntelligenceInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverBusinessIntelligenceInstance` [%s] not found, please check if it has been deleted.  ", logId, d.Id())
		return nil
	}

	if businessIntelligenceInstance.Zone != nil {
		_ = d.Set("zone", businessIntelligenceInstance.Zone)
	}

	if businessIntelligenceInstance.Memory != nil {
		_ = d.Set("memory", businessIntelligenceInstance.Memory)
	}

	if businessIntelligenceInstance.Storage != nil {
		_ = d.Set("storage", businessIntelligenceInstance.Storage)
	}

	if businessIntelligenceInstance.Cpu != nil {
		_ = d.Set("cpu", businessIntelligenceInstance.Cpu)
	}

	if businessIntelligenceInstance.Type != nil {
		_ = d.Set("machine_type", businessIntelligenceInstance.Type)
	}

	if businessIntelligenceInstance.ProjectId != nil {
		_ = d.Set("project_id", businessIntelligenceInstance.ProjectId)
	}

	if businessIntelligenceInstance.UniqSubnetId != nil {
		_ = d.Set("subnet_id", businessIntelligenceInstance.UniqSubnetId)
	}

	if businessIntelligenceInstance.UniqVpcId != nil {
		_ = d.Set("vpc_id", businessIntelligenceInstance.UniqVpcId)
	}

	if businessIntelligenceInstance.VersionName != nil {
		var dbVersion string
		if *businessIntelligenceInstance.VersionName == SQLSERVER_DB_VERSION_NAME_2016 {
			dbVersion = SQLSERVER_DB_VERSION_2016
		} else if *businessIntelligenceInstance.VersionName == SQLSERVER_DB_VERSION_NAME_2017 {
			dbVersion = SQLSERVER_DB_VERSION_2017
		} else if *businessIntelligenceInstance.VersionName == SQLSERVER_DB_VERSION_NAME_2019 {
			dbVersion = SQLSERVER_DB_VERSION_2019
		} else {
			dbVersion = SQLSERVER_DB_VERSION_2019
		}
		_ = d.Set("db_version", dbVersion)
	}

	if businessIntelligenceInstance.Name != nil {
		_ = d.Set("instance_name", businessIntelligenceInstance.Name)
	}

	if businessIntelligenceInstance.ResourceTags != nil {
		resourceTagsList := []interface{}{}
		for _, resourceTags := range businessIntelligenceInstance.ResourceTags {
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

	maintenanceSpan, err := service.DescribeMaintenanceSpanById(ctx, instanceId)
	if err != nil {
		return err
	}

	if maintenanceSpan == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlservereMaintenanceSpan` [%s] not found, please check if it has been deleted.  ", logId, d.Id())
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

	return nil
}

func resourceTencentCloudSqlserverBusinessIntelligenceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_business_intelligence_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{"zone", "memory", "storage", "cpu", "machine_type", "project_id", "goods_num", "subnet_id", "vpc_id", "d_b_version", "security_group_list", "weekly", "start_time", "span", "resource_tags"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = sqlserver.NewModifyDBInstanceNameRequest()
		instanceId = d.Id()
	)

	request.InstanceId = &instanceId

	if d.HasChange("instance_name") {
		if v, ok := d.GetOk("instance_name"); ok {
			request.InstanceName = helper.String(v.(string))
		}
	} else {
		return nil
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().ModifyDBInstanceName(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]  ", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver businessIntelligenceInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverBusinessIntelligenceInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverBusinessIntelligenceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_business_intelligence_instance.delete")()
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

	if err := service.DeleteSqlserverInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
