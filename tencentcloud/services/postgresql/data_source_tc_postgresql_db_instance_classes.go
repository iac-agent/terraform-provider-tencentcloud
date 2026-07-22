package postgresql

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPostgresqlDbInstanceClasses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresqlDbInstanceClassesRead,
		Schema: map[string]*schema.Schema{
			"zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "AZ ID，其中 可以 是 获取 through `DescribeZones` API。",
			},

			"db_engine": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Database engines. 有效值：1. `postgresql` (TencentDB 对于 PostgreSQL) 2. `mssql_compatible` (MSSQL compatible-TencentDB 对于 PostgreSQL)。",
			},

			"db_major_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Major 版本 的 数据库，such 作为 12 或 13，其中 可以 是 获取 through `DescribeDBVersions` API。",
			},

			"storage_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Storage 类型 过滤器. 有效值：`PHYSICAL_LOCAL_SSD` (本地 SSD)，`CLOUD_PREMIUM` (premium 云 磁盘)，`CLOUD_SSD` (云 SSD)，`CLOUD_HSSD` (enhanced 云 SSD)。",
			},

			"class_info_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 数据库 specifications。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spec_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specification ID。",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CPU 核数",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory 大小 （MB）。",
						},
						"max_storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum 存储 容量 在 GB 支持 通过 此 规格。",
						},
						"min_storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum 存储 容量 在 GB 支持 通过 此 规格。",
						},
						"qps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Estimated QPS 对于 此 规格。",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudPostgresqlDbInstanceClassesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_postgresql_db_instance_classes.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone"); ok {
		paramMap["Zone"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_engine"); ok {
		paramMap["DBEngine"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_major_version"); ok {
		paramMap["DBMajorVersion"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("storage_type"); ok {
		paramMap["StorageType"] = helper.String(v.(string))
	}

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var classInfoSet []*postgresql.ClassInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePostgresqlDbInstanceClassesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		classInfoSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(classInfoSet))
	tmpList := make([]map[string]interface{}, 0, len(classInfoSet))

	if classInfoSet != nil {
		for _, classInfo := range classInfoSet {
			classInfoMap := map[string]interface{}{}

			if classInfo.SpecCode != nil {
				classInfoMap["spec_code"] = classInfo.SpecCode
			}

			if classInfo.CPU != nil {
				classInfoMap["cpu"] = classInfo.CPU
			}

			if classInfo.Memory != nil {
				classInfoMap["memory"] = classInfo.Memory
			}

			if classInfo.MaxStorage != nil {
				classInfoMap["max_storage"] = classInfo.MaxStorage
			}

			if classInfo.MinStorage != nil {
				classInfoMap["min_storage"] = classInfo.MinStorage
			}

			if classInfo.QPS != nil {
				classInfoMap["qps"] = classInfo.QPS
			}

			ids = append(ids, *classInfo.SpecCode)
			tmpList = append(tmpList, classInfoMap)
		}

		_ = d.Set("class_info_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
