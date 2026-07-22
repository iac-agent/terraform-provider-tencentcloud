package cos

import (
	"context"
	"encoding/json"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCosBatchs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCosBatchsRead,

		Schema: map[string]*schema.Schema{
			"uin": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Uin。",
			},
			"appid": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Appid。",
			},
			"job_statuses": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "任务 状态 信息 您 need 到 查询. 如果 您 do 不 指定a 任务 状态，COS 返回status 的 all tasks 该 have been executed，包括 those 该 是 在 progress. 如果 您 指定a 任务 状态，COS 返回task 在 指定 state. 可选 任务 states include: 活跃，Cancelled，Cancelling，Complete，Completing，Failed，Failing，New，Paused，Pausing，Preparing，Ready，Suspended。",
			},
			"jobs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Multiple batch processing 任务 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"creation_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Job 创建时间。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mission 描述 长度 是 limited 到 0-256 bytes。",
						},
						"job_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "作业 ID 长度 是 limited 到 1-64 bytes。",
						},
						"operation": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Actions performed 在 objects 在 batch processing 作业. For 示例，COSPutObjectCopy。",
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Mission 优先级 Tasks 使用 higher 值 将 是 given 优先级 优先级 大小 是 limited 到 0-2147483647。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "任务 execution 状态 Legal 参数 值 include 活跃，Cancelled，Cancelling，Complete，Completing，Failed，Failing，New，Paused，Pausing，Preparing，Ready，Suspended。",
						},
						"termination_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "结束时间 的 batch processing 作业。",
						},
						"progress_summary": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Summary 的 状态 任务 implementation. Describe 总数 数量 operations performed 在 此 任务， 数量 successful operations，和 数量 failed operations。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"number_of_tasks_failed": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "当前 failed Operand。",
									},
									"number_of_tasks_succeeded": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "当前 successful Operand。",
									},
									"total_number_of_tasks": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total Operand。",
									},
								},
							},
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

func dataSourceTencentCloudCosBatchsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cos_batchs.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	uin := d.Get("uin").(string)
	appid := d.Get("appid").(int)
	jobs := make([]map[string]interface{}, 0)

	opt := &cos.BatchListJobsOptions{}
	if v, ok := d.GetOk("job_statuses"); ok {
		opt.JobStatuses = v.(string)
	}
	headers := &cos.BatchRequestHeaders{
		XCosAppid: appid,
	}
	ids := make([]string, 0)
	for {
		req, _ := json.Marshal(opt)
		result, response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosBatchClient(uin).Batch.ListJobs(ctx, nil, headers)
		responseBody, _ := json.Marshal(response.Body)
		log.Printf("[DEBUG]%s api[ListJobs] success, request body [%s], response body [%v]\n", logId, req, responseBody)
		if err != nil {
			return err
		}
		for _, item := range result.Jobs.Members {
			jobItem := make(map[string]interface{})
			jobItem["creation_time"] = item.CreationTime
			jobItem["description"] = item.Description
			jobItem["job_id"] = item.JobId
			jobItem["operation"] = item.Operation
			jobItem["priority"] = item.Priority
			jobItem["status"] = item.Status
			jobItem["termination_date"] = item.TerminationDate
			progressSummary := map[string]interface{}{
				"number_of_tasks_failed":    item.ProgressSummary.NumberOfTasksFailed,
				"number_of_tasks_succeeded": item.ProgressSummary.NumberOfTasksSucceeded,
				"total_number_of_tasks":     item.ProgressSummary.TotalNumberOfTasks,
			}
			jobItem["progress_summary"] = []interface{}{progressSummary}
			ids = append(ids, item.JobId)
			jobs = append(jobs, jobItem)
		}
		if result.NextToken != "" {
			opt.NextToken = result.NextToken
		} else {
			break
		}
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("jobs", jobs)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), jobs); err != nil {
			return err
		}
	}

	return nil
}
