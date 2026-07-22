package postgresql

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTencentCloudPostgresqlSpecinfos() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresqlSpecinfosRead,
		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "可用区 的 postgresql 实例 到 查询。",
			},
			"storage_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Storage 类型 过滤器. 有效值：`PHYSICAL_LOCAL_SSD` (本地 SSD)，`CLOUD_PREMIUM` (premium 云 磁盘)，`CLOUD_SSD` (云 SSD)，`CLOUD_HSSD` (enhanced 云 SSD)。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 zones 将 是 exported 和 its every element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID postgresql 实例 speccode。",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory 大小(在 GB)。",
						},
						"storage_min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最小 卷 大小(在 GB)。",
						},
						"storage_max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最大 卷 大小(在 GB)。",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CPU 数量 postgresql 实例。",
						},
						"qps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "QPS 的 postgresql 实例。",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "版本 的 postgresql 数据库 引擎。",
						},
						"engine_version_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "版本 名称 postgresql 数据库 引擎。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudPostgresqlSpecinfosRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_postgresql_specinfos.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := PostgresqlService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	zone := d.Get("availability_zone").(string)
	storageType := ""
	if v, ok := d.GetOk("storage_type"); ok {
		storageType = v.(string)
	}
	speccodes, err := service.DescribeSpecinfos(ctx, zone, storageType)
	if err != nil {
		speccodes, err = service.DescribeSpecinfos(ctx, zone, storageType)
	}
	if err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(speccodes))
	for _, v := range speccodes {
		listItem := make(map[string]interface{})
		listItem["id"] = v.SpecCode
		listItem["memory"] = *v.Memory / 1024
		listItem["storage_min"] = v.MinStorage
		listItem["storage_max"] = v.MaxStorage
		listItem["cpu"] = v.Cpu
		listItem["qps"] = v.Qps
		listItem["engine_version"] = v.Version
		listItem["engine_version_name"] = v.VersionName
		list = append(list, listItem)
	}

	d.SetId("speccode." + zone)
	if e := d.Set("list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}

	return nil
}
