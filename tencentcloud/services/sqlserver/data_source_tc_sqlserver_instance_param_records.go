package sqlserver

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSqlserverInstanceParamRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverInstanceParamRecordsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID 在 格式 的 mssql-dj5i29c5n. It 是 same 作为 实例 ID displayed 在 TencentDB console 和 response 参数 InstanceId 的 DescribeDBInstances API.",
			},
			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Parameter modification records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID.",
						},
						"param_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter 名称.",
						},
						"old_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter 值 before modification.",
						},
						"new_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter 值 after modification.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Parameter modification 状态. 有效 值: 1 (initializing 和 waiting 对于 modification), 2 (modification succeed), 3 (modification failed), 4 (modifying).",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modification 时间.",
						},
					},
				},
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 save results.",
			},
		},
	}
}

func dataSourceTencentCloudSqlserverInstanceParamRecordsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_sqlserver_instance_param_records.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	var items []*sqlserver.ParamRecord

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverInstanceParamRecordsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		items = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(items))

	if items != nil {
		for _, paramRecord := range items {
			paramRecordMap := map[string]interface{}{}

			if paramRecord.InstanceId != nil {
				paramRecordMap["instance_id"] = paramRecord.InstanceId
			}

			if paramRecord.ParamName != nil {
				paramRecordMap["param_name"] = paramRecord.ParamName
			}

			if paramRecord.OldValue != nil {
				paramRecordMap["old_value"] = paramRecord.OldValue
			}

			if paramRecord.NewValue != nil {
				paramRecordMap["new_value"] = paramRecord.NewValue
			}

			if paramRecord.Status != nil {
				paramRecordMap["status"] = paramRecord.Status
			}

			if paramRecord.ModifyTime != nil {
				paramRecordMap["modify_time"] = paramRecord.ModifyTime
			}

			tmpList = append(tmpList, paramRecordMap)
		}

		_ = d.Set("items", tmpList)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
