package cdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMysqlInstanceParamRecord() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlInstanceParamRecordRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例ID，格式为：cdb-c1nl9rpv，与云数据库控制台页面显示的实例ID相同，可以通过【查询实例列表】（https://云.tencent.com/document/api/236/15872）接口获取输出参数中InstanceId字段的值。",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "参数修改记录。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID。",
						},
						"param_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "参数名称。",
						},
						"old_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "修改前的参数值。",
						},
						"new_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "参数的修改值。",
						},
						"is_success": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "参数是否修改成功。",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "改变时间。",
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

func dataSourceTencentCloudMysqlInstanceParamRecordRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_instance_param_record.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var instanceParamRecord []*cdb.ParamRecord
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlInstanceParamRecordByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instanceParamRecord = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceParamRecord))
	tmpList := make([]map[string]interface{}, 0, len(instanceParamRecord))
	if instanceParamRecord != nil {
		for _, paramRecord := range instanceParamRecord {
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

			if paramRecord.IsSucess != nil {
				paramRecordMap["is_success"] = paramRecord.IsSucess
			}

			if paramRecord.ModifyTime != nil {
				paramRecordMap["modify_time"] = paramRecord.ModifyTime
			}

			ids = append(ids, *paramRecord.InstanceId)
			tmpList = append(tmpList, paramRecordMap)
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
