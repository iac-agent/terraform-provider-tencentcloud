package scf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudScfRequestStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfRequestStatusRead,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Function 名称",
			},

			"function_request_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID 请求 到 是 queried。",
			},

			"namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Function 命名空间。",
			},

			"start_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "开始时间 的 查询，对于 示例 `2017-05-16 20:00:00`. 如果 它's left 空，它 默认为 15 minutes before 当前 时间。",
			},

			"end_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "结束时间 的 查询. such 作为 `2017-05-16 20:59:59`. 如果 `StartTime` 是 不 指定，`EndTime` 默认为 当前 时间. 如果 `StartTime` 是 指定，`EndTime` 为必填项，和 它 need 到 是 later 比 `StartTime`。",
			},

			"data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Details 的 函数 running statusNote: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"function_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Function 名称",
						},
						"ret_msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "返回值 after 函数 是 executed。",
						},
						"request_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "请求 ID",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Request 开始时间。",
						},
						"ret_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "结果 的 请求. `0`: succeeded，`1`: running，`-1`: exception。",
						},
						"duration": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Time consumed 对于 请求 在 ms。",
						},
						"mem_usage": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Time consumed 通过 请求 （MB）。",
						},
						"retry_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Retry Attempts。",
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

func dataSourceTencentCloudScfRequestStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_scf_request_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("function_name"); ok {
		paramMap["FunctionName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("function_request_id"); ok {
		paramMap["FunctionRequestId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		paramMap["Namespace"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	service := ScfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var data []*scf.RequestStatus

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeScfRequestStatusByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		data = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(data))
	tmpList := make([]map[string]interface{}, 0, len(data))

	if data != nil {
		for _, requestStatus := range data {
			requestStatusMap := map[string]interface{}{}

			if requestStatus.FunctionName != nil {
				requestStatusMap["function_name"] = requestStatus.FunctionName
			}

			if requestStatus.RetMsg != nil {
				requestStatusMap["ret_msg"] = requestStatus.RetMsg
			}

			if requestStatus.RequestId != nil {
				requestStatusMap["request_id"] = requestStatus.RequestId
			}

			if requestStatus.StartTime != nil {
				requestStatusMap["start_time"] = requestStatus.StartTime
			}

			if requestStatus.RetCode != nil {
				requestStatusMap["ret_code"] = requestStatus.RetCode
			}

			if requestStatus.Duration != nil {
				requestStatusMap["duration"] = requestStatus.Duration
			}

			if requestStatus.MemUsage != nil {
				requestStatusMap["mem_usage"] = requestStatus.MemUsage
			}

			if requestStatus.RetryNum != nil {
				requestStatusMap["retry_num"] = requestStatus.RetryNum
			}

			ids = append(ids, *requestStatus.FunctionName)
			tmpList = append(tmpList, requestStatusMap)
		}

		_ = d.Set("data", tmpList)
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
