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
				Description: "任务 状态 info。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current 使用 来源 URL",
						},
						"looped_times": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 times VOD 来源 任务 是 played 在 loop。",
						},
						"offset_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "playback 偏移量 的 VOD 来源，（秒）。",
						},
						"report_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "latest heartbeat 报告 时间 在 UTC 格式，对于 示例: 2022-02-11T10:00:00Z.注意: UTC 时间 是 8 hours ahead 的 Beijing 时间。",
						},
						"run_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Real run 状态:活跃,inactive。",
						},
						"file_duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "时长 的 VOD 来源 文件，（秒）。",
						},
						"next_file_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL 的 next progress VOD 文件。",
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
