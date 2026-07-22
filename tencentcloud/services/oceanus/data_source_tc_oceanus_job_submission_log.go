package oceanus

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudOceanusJobSubmissionLog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOceanusJobSubmissionLogRead,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "作业 ID",
			},
			"start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "开始时间，unix 时间戳，（毫秒）。",
			},
			"end_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "结束时间，unix 时间戳，（毫秒）。",
			},
			"running_order_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Default:     0,
				Description: "Job 实例 ID。",
			},
			"keyword": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Keyword，默认值 空。",
			},
			"cursor": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Cursor，默认值 空，first 请求 does 不 need 到 pass 在。",
			},
			"order_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Default:     "asc",
				Description: "Sorting 方法，默认值 asc，asc: ascending，desc: descending。",
			},
			"list_over": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "是否list 是 over。",
			},
			"job_request_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "请求 ID starting 作业。",
			},
			"log_list": {
				Computed:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Log 列表，已弃用",
			},
			"job_instance_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Job 实例 列表 during 指定 时间 周期",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"running_order_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID 实例，starting 从 1 在 顺序 的 startup 时间。",
						},
						"job_instance_start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "startup 时间 的 实例。",
						},
						"starting_millis": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "startup 时间 的 实例 （毫秒）。",
						},
					},
				},
			},
			"log_content_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 日志 contents。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "内容 的 日志。",
						},
						"time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "时间戳 （毫秒）。",
						},
						"pkg_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 日志 组。",
						},
						"pkg_log_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID 日志，其中 是 唯一 within 日志 组。",
						},
						"container_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 容器 到 其中 日志 belongs。",
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

func dataSourceTencentCloudOceanusJobSubmissionLogRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_oceanus_job_submission_log.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		ctx              = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service          = OceanusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		jobSubmissionLog *oceanus.DescribeJobSubmissionLogResponseParams
		jobId            string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("job_id"); ok {
		paramMap["JobId"] = helper.String(v.(string))
		jobId = v.(string)
	}

	if v, ok := d.GetOkExists("start_time"); ok {
		paramMap["StartTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("end_time"); ok {
		paramMap["EndTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("running_order_id"); ok {
		paramMap["RunningOrderId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("keyword"); ok {
		paramMap["Keyword"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_type"); ok {
		paramMap["OrderType"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("cursor"); v != nil {
		paramMap["Cursor"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOceanusJobSubmissionLogByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		jobSubmissionLog = result
		return nil
	})

	if err != nil {
		return err
	}

	if jobSubmissionLog.Cursor != nil {
		_ = d.Set("cursor", jobSubmissionLog.Cursor)
	}

	if jobSubmissionLog.ListOver != nil {
		_ = d.Set("list_over", jobSubmissionLog.ListOver)
	}

	if jobSubmissionLog.JobRequestId != nil {
		_ = d.Set("job_request_id", jobSubmissionLog.JobRequestId)
	}

	if jobSubmissionLog.LogList != nil {
		tmpList := make([]string, 0, len(jobSubmissionLog.LogList))
		for _, log := range jobSubmissionLog.LogList {
			tmpList = append(tmpList, *log)
		}

		_ = d.Set("log_list", tmpList)
	}

	if jobSubmissionLog.JobInstanceList != nil {
		tmpList := make([]map[string]interface{}, 0, len(jobSubmissionLog.JobInstanceList))
		for _, jobInstance := range jobSubmissionLog.JobInstanceList {
			jobInstanceMap := map[string]interface{}{}

			if jobInstance.RunningOrderId != nil {
				jobInstanceMap["running_order_id"] = jobInstance.RunningOrderId
			}

			if jobInstance.JobInstanceStartTime != nil {
				jobInstanceMap["job_instance_start_time"] = jobInstance.JobInstanceStartTime
			}

			if jobInstance.StartingMillis != nil {
				jobInstanceMap["starting_millis"] = jobInstance.StartingMillis
			}

			tmpList = append(tmpList, jobInstanceMap)
		}

		_ = d.Set("job_instance_list", tmpList)
	}

	if jobSubmissionLog.LogContentList != nil {
		tmpList := make([]map[string]interface{}, 0, len(jobSubmissionLog.LogContentList))
		for _, logContent := range jobSubmissionLog.LogContentList {
			logContentMap := map[string]interface{}{}

			if logContent.Log != nil {
				logContentMap["log"] = logContent.Log
			}

			if logContent.Time != nil {
				logContentMap["time"] = logContent.Time
			}

			if logContent.PkgId != nil {
				logContentMap["pkg_id"] = logContent.PkgId
			}

			if logContent.PkgLogId != nil {
				logContentMap["pkg_log_id"] = logContent.PkgLogId
			}

			if logContent.ContainerName != nil {
				logContentMap["container_name"] = logContent.ContainerName
			}

			tmpList = append(tmpList, logContentMap)
		}

		_ = d.Set("log_content_list", tmpList)
	}

	d.SetId(jobId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
