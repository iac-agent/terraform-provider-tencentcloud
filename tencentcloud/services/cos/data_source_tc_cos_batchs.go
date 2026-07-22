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
				Description: "The task 状态 information you need to query. If you do not 指定a task 状态，COS 返回status of all tasks that have been executed，including those that are in progress. If you 指定a task 状态，COS 返回task in the specified state. 可选 task states include: 活跃，Cancelled，Cancelling，Complete，Completing，Failed，Failing，New，Paused，Pausing，Preparing，Ready，Suspended。",
			},
			"jobs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Multiple batch processing task information。",
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
							Description: "Mission 描述 The length is limited to 0-256 bytes。",
						},
						"job_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "作业 ID The length is limited to 1-64 bytes。",
						},
						"operation": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Actions performed on objects in a batch processing job. For example，COSPutObjectCopy。",
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Mission 优先级 Tasks with higher values will be given 优先级 The 优先级 size is limited to 0-2147483647。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task execution 状态 Legal parameter values include 活跃，Cancelled，Cancelling，Complete，Completing，Failed，Failing，New，Paused，Pausing，Preparing，Ready，Suspended。",
						},
						"termination_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The 结束时间 of the batch processing job。",
						},
						"progress_summary": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Summary of 状态 task implementation. Describe the total 数量 operations performed in this task，the 数量 successful operations，and the 数量 failed operations。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"number_of_tasks_failed": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The current failed Operand。",
									},
									"number_of_tasks_succeeded": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The current successful Operand。",
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
