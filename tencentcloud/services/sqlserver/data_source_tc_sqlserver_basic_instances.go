package sqlserver

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcpostgresql "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/postgresql"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSqlserverBasicInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverBasicInstanceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID 的 SQL Server basic 实例 到 是 查询.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name 的 SQL Server basic 实例 到 是 查询.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Project ID 的 SQL Server basic 实例 到 是 查询.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Vpc ID 的 SQL Server basic 实例 到 是 查询.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subnet ID 的 SQL Server basic 实例 到 是 查询.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 save results.",
			},
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 的 SQL Server basic 实例. Each element contains following attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 的 SQL Server basic 实例.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name 的 SQL Server basic 实例.",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pay 类型 的 SQL Server basic 实例. For now, 仅 `POSTPAID_BY_HOUR` 是 有效.",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version 的 SQL Server basic 数据库 引擎. Allowed 值 是 `2008R2`(SQL Server 2008 Enterprise), `2012SP3`(SQL Server 2012 Enterprise), `2016SP1` (SQL Server 2016 Enterprise), `201602`(SQL Server 2016 Standard) 和 `2017`(SQL Server 2017 Enterprise). Default 是 `2008R2`.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 的 VPC.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 的 子网.",
						},
						"storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk 大小 (在 GB). Allowed 值 必须 是 多个 的 10. 存储 必须 是 集合 使用 限制 的 `storage_min` 和 `storage_max` 其中 数据 source `tencentcloud_sqlserver_specinfos` provides.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory 大小 (在 GB). Allowed 值 必须 是 larger 比 `内存` 该 数据 source `tencentcloud_sqlserver_specinfos` provides.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CPU 数量 的 SQL Server basic 实例.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project ID, 默认值 值 是 `0`.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone.",
						},
						"used_storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Used 存储.",
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
							Description: "Status 的 SQL Server basic 实例. `1` 对于 applying, `2` 对于 running, `3` 对于 running 使用 限制, `4` 对于 isolated, `5` 对于 recycling, `6` 对于 recycled, `7` 对于 running 使用 任务, `8` 对于 关闭-line, `9` 对于 expanding, `10` 对于 migrating, `11` 对于 readonly, `12` 对于 rebooting.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Tags 的 SQL Server basic 实例.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudSqlserverBasicInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_sqlserver_basic_instances.read")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		tcClient   = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService = svctag.NewTagService(tcClient)
		service    = SqlserverService{client: tcClient}
		id         = d.Get("id").(string)
		name       = d.Get("name").(string)
		projectId  = -1
		vpcId      string
		subnetId   string
	)
	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(int)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		vpcId = v.(string)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		subnetId = v.(string)
	}
	instanceList, err := service.DescribeSqlserverInstances(ctx, id, name, projectId, vpcId, subnetId, 1)
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceList))
	list := make([]map[string]interface{}, 0, len(instanceList))
	for _, v := range instanceList {
		listItem := make(map[string]interface{})
		listItem["id"] = v.InstanceId
		listItem["name"] = v.Name
		listItem["project_id"] = v.ProjectId
		listItem["storage"] = v.Storage
		listItem["memory"] = v.Memory
		listItem["availability_zone"] = v.Zone
		listItem["create_time"] = v.CreateTime
		listItem["vpc_id"] = v.UniqVpcId
		listItem["subnet_id"] = v.UniqSubnetId
		listItem["engine_version"] = v.Version
		listItem["vip"] = v.Vip
		listItem["vport"] = v.Vport
		listItem["used_storage"] = v.UsedStorage
		listItem["status"] = v.Status
		listItem["cpu"] = v.Cpu

		if *v.PayMode == 1 {
			listItem["charge_type"] = svcpostgresql.COMMON_PAYTYPE_PREPAID
		} else {
			listItem["charge_type"] = svcpostgresql.COMMON_PAYTYPE_POSTPAID
		}
		tagList, err := tagService.DescribeResourceTags(ctx, "sqlserver", "instance", tcClient.Region, *v.InstanceId)
		if err != nil {
			return err
		}

		listItem["tags"] = tagList
		list = append(list, listItem)
		ids = append(ids, *v.InstanceId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("instance_list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}

	return nil

}
