package cdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMysqlInstanceRebootTime() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlInstanceRebootTimeRead,
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "实例ID与云数据库控制台页面显示的实例ID一致，格式为：cdb-c1nl9rpv。",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "返回的参数信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例ID，格式为：cdb-c1nl9rpv，与云数据库控制台页面显示的实例ID相同。",
						},
						"time_in_seconds": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "预计重启时间。",
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

func dataSourceTencentCloudMysqlInstanceRebootTimeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_instance_reboot_time.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		paramMap["InstanceIds"] = helper.InterfacesStringsPoint(instanceIdsSet)
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var instanceRebootTime []*cdb.InstanceRebootTime
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlInstanceRebootTimeByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instanceRebootTime = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceRebootTime))
	tmpList := make([]map[string]interface{}, 0, len(instanceRebootTime))
	if instanceRebootTime != nil {
		for _, instanceRebootTime := range instanceRebootTime {
			instanceRebootTimeMap := map[string]interface{}{}

			if instanceRebootTime.InstanceId != nil {
				instanceRebootTimeMap["instance_id"] = instanceRebootTime.InstanceId
			}

			if instanceRebootTime.TimeInSeconds != nil {
				instanceRebootTimeMap["time_in_seconds"] = instanceRebootTime.TimeInSeconds
			}

			ids = append(ids, *instanceRebootTime.InstanceId)
			tmpList = append(tmpList, instanceRebootTimeMap)
		}

		_ = d.Set("items", tmpList)
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
