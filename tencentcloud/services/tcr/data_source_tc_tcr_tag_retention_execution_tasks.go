package tcr

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTcrTagRetentionExecutionTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcrTagRetentionExecutionTasksRead,
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID",
			},

			"retention_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "retention ID。",
			},

			"execution_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "execution ID。",
			},

			"retention_task_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 版本 retention tasks。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "任务 ID",
						},
						"execution_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "规则 execution ID。",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "任务 开始时间。",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "任务 结束时间。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "execution 状态 任务: Failed，Succeed，Stopped，InProgress。",
						},
						"total": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total 数量 标签",
						},
						"retained": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total 数量 retained 标签",
						},
						"repository": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "repository 名称",
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

func dataSourceTencentCloudTcrTagRetentionExecutionTasksRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tcr_tag_retention_execution_tasks.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		registryId  string
		retentionId string
		executionId string
	)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("registry_id"); ok {
		paramMap["registry_id"] = helper.String(v.(string))
		registryId = v.(string)
	}

	if v, _ := d.GetOk("retention_id"); v != nil {
		paramMap["retention_id"] = helper.IntInt64(v.(int))
		retentionId = helper.IntToStr(v.(int))
	}

	if v, _ := d.GetOk("execution_id"); v != nil {
		paramMap["execution_id"] = helper.IntInt64(v.(int))
		executionId = helper.IntToStr(v.(int))
	}

	service := TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var retentionTaskList []*tcr.RetentionTask

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTcrTagRetentionExecutionTasksByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		retentionTaskList = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(retentionTaskList))

	if retentionTaskList != nil {
		for _, retentionTask := range retentionTaskList {
			retentionTaskMap := map[string]interface{}{}

			if retentionTask.TaskId != nil {
				retentionTaskMap["task_id"] = retentionTask.TaskId
			}

			if retentionTask.ExecutionId != nil {
				retentionTaskMap["execution_id"] = retentionTask.ExecutionId
			}

			if retentionTask.StartTime != nil {
				retentionTaskMap["start_time"] = retentionTask.StartTime
			}

			if retentionTask.EndTime != nil {
				retentionTaskMap["end_time"] = retentionTask.EndTime
			}

			if retentionTask.Status != nil {
				retentionTaskMap["status"] = retentionTask.Status
			}

			if retentionTask.Total != nil {
				retentionTaskMap["total"] = retentionTask.Total
			}

			if retentionTask.Retained != nil {
				retentionTaskMap["retained"] = retentionTask.Retained
			}

			if retentionTask.Repository != nil {
				retentionTaskMap["repository"] = retentionTask.Repository
			}

			tmpList = append(tmpList, retentionTaskMap)
		}

		_ = d.Set("retention_task_list", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash([]string{registryId, retentionId, executionId}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
