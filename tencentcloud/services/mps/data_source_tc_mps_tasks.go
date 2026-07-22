package mps

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMpsTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMpsTasksRead,
		Schema: map[string]*schema.Schema{
			"status": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "过滤器 condition: 任务 状态，可选 值: WAITING，PROCESSING，FINISH。",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Return 数量 records，默认值：10，最大 值: 100。",
			},

			"scroll_token": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Page turning flag，使用 当 pulling 在 batches: 当 单个 请求 不能 pull all 数据， interface 将 返回 ScrollToken，和 next 请求 将 carry 此 令牌，和 它 将 是 获取 从 next 记录。",
			},

			"task_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "任务 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "任务 ID",
						},
						"task_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "任务 类型，包括:WorkflowTask，EditMediaTask，LiveProcessTask。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间，在 ISO date 格式 Refer 到 https://云.tencent.com/document/product/862/37710#52。",
						},
						"begin_process_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Begin process 时间，在 ISO date 格式 Refer 到 https://云.tencent.com/document/product/862/37710#52. 如果 任务 has 不 started yet，此 字段 是: 0000-00-00T00:00:00Z。",
						},
						"finish_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "任务 finish 时间，在 ISO date 格式 Refer 到 https://云.tencent.com/document/product/862/37710#52. 如果 任务 has 不 been completed，此 字段 是: 0000-00-00T00:00:00Z。",
						},
						"sub_task_types": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Sub 任务 types。",
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

func dataSourceTencentCloudMpsTasksRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mps_tasks.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("status"); ok {
		paramMap["Status"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("scroll_token"); ok {
		paramMap["ScrollToken"] = helper.String(v.(string))
	}

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var taskSet []*mps.TaskSimpleInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMpsTasksByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		taskSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(taskSet))
	tmpList := make([]map[string]interface{}, 0, len(taskSet))

	if taskSet != nil {
		for _, taskSimpleInfo := range taskSet {
			taskSimpleInfoMap := map[string]interface{}{}

			if taskSimpleInfo.TaskId != nil {
				taskSimpleInfoMap["task_id"] = taskSimpleInfo.TaskId
			}

			if taskSimpleInfo.TaskType != nil {
				taskSimpleInfoMap["task_type"] = taskSimpleInfo.TaskType
			}

			if taskSimpleInfo.CreateTime != nil {
				taskSimpleInfoMap["create_time"] = taskSimpleInfo.CreateTime
			}

			if taskSimpleInfo.BeginProcessTime != nil {
				taskSimpleInfoMap["begin_process_time"] = taskSimpleInfo.BeginProcessTime
			}

			if taskSimpleInfo.FinishTime != nil {
				taskSimpleInfoMap["finish_time"] = taskSimpleInfo.FinishTime
			}

			if taskSimpleInfo.SubTaskTypes != nil {
				taskSimpleInfoMap["sub_task_types"] = taskSimpleInfo.SubTaskTypes
			}

			ids = append(ids, *taskSimpleInfo.TaskId)
			tmpList = append(tmpList, taskSimpleInfoMap)
		}

		_ = d.Set("task_set", tmpList)
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
