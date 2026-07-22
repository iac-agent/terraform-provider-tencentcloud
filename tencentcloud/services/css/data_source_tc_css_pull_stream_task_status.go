package css

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCssPullStreamTaskStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCssPullStreamTaskStatusRead,
		Schema: map[string]*schema.Schema{
			"task_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "任务 ID",
			},

			"task_status_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Task 状态 info。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current use 来源 URL",
						},
						"looped_times": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The 数量 times a VOD 来源 task is played in a loop。",
						},
						"offset_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The playback 偏移量 of the VOD 来源，（秒）。",
						},
						"report_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The latest heartbeat reporting time in UTC 格式，for example: 2022-02-11T10:00:00Z.Note: UTC time is 8 hours ahead of Beijing time。",
						},
						"run_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Real run 状态:活跃,inactive。",
						},
						"file_duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The duration of the VOD 来源 file，（秒）。",
						},
						"next_file_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of the next progress VOD file。",
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

func dataSourceTencentCloudCssPullStreamTaskStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_css_pull_stream_task_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var taskId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("task_id"); ok {
		taskId = v.(string)
		paramMap["TaskId"] = helper.String(v.(string))
	}

	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var taskStatusInfo *css.TaskStatusInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCssPullStreamTaskStatusByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		taskStatusInfo = result
		return nil
	})
	if err != nil {
		return err
	}

	taskStatusInfoMap := map[string]interface{}{}
	if taskStatusInfo != nil {
		if taskStatusInfo.FileUrl != nil {
			taskStatusInfoMap["file_url"] = taskStatusInfo.FileUrl
		}

		if taskStatusInfo.LoopedTimes != nil {
			taskStatusInfoMap["looped_times"] = taskStatusInfo.LoopedTimes
		}

		if taskStatusInfo.OffsetTime != nil {
			taskStatusInfoMap["offset_time"] = taskStatusInfo.OffsetTime
		}

		if taskStatusInfo.ReportTime != nil {
			taskStatusInfoMap["report_time"] = taskStatusInfo.ReportTime
		}

		if taskStatusInfo.RunStatus != nil {
			taskStatusInfoMap["run_status"] = taskStatusInfo.RunStatus
		}

		if taskStatusInfo.FileDuration != nil {
			taskStatusInfoMap["file_duration"] = taskStatusInfo.FileDuration
		}

		if taskStatusInfo.NextFileUrl != nil {
			taskStatusInfoMap["next_file_url"] = taskStatusInfo.NextFileUrl
		}

		_ = d.Set("task_status_info", []interface{}{taskStatusInfoMap})
	}

	d.SetId(helper.DataResourceIdsHash([]string{taskId}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), taskStatusInfoMap); e != nil {
			return e
		}
	}
	return nil
}
