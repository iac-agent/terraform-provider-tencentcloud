package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataSqlScriptRuns() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataSqlScriptRunsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "项目 ID",
			},

			"script_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Script ID。",
			},

			"job_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "作业 ID",
			},

			"search_word": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search keyword。",
			},

			"execute_user_uin": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Execute 用户 UIN。",
			},

			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "开始时间。",
			},

			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "结束时间。",
			},

			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Data exploration tasks。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data exploration 任务 ID。",
						},
						"job_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data exploration 任务 名称",
						},
						"job_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Job 类型",
						},
						"script_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Script ID。",
						},
						"job_execution_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Subtask 列表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"job_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data exploration 任务 ID。",
									},
									"job_execution_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subquery 任务 ID。",
									},
									"job_execution_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subquery 名称",
									},
									"script_content": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subquery SQL 内容",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subquery 状态",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创建时间。",
									},
									"execute_stage_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Execution phase。",
									},
									"log_file_path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Log 文件 路径",
									},
									"result_file_path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "结果 文件 路径",
									},
									"result_preview_file_path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Preview 结果 文件 路径",
									},
									"result_total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total 数量 rows 在 任务 execution 结果",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "更新时间。",
									},
									"end_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "结束时间。",
									},
									"time_cost": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Time consumed。",
									},
									"context_script_content": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Context SQL 内容",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"result_preview_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 rows 对于 previewing 任务 execution results。",
									},
									"result_effect_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 rows affected 通过 任务 execution 结果",
									},
									"collecting_total_result": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether collecting full results: 默认值 false，true 表示 collecting full results，用于frontend polling。",
									},
									"script_content_truncate": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否script 内容 是 truncated。",
									},
								},
							},
						},
						"script_content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Script 内容",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "任务 状态",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "任务 创建时间。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "更新时间。",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "结束时间。",
						},
						"owner_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud 所有者 账号 UIN。",
						},
						"user_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "账号 UIN。",
						},
						"time_cost": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Time consumed。",
						},
						"script_content_truncate": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否script 内容 是 truncated。",
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

func dataSourceTencentCloudWedataSqlScriptRunsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_sql_script_runs.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("project_id"); ok {
		paramMap["ProjectId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("script_id"); ok {
		paramMap["ScriptId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("job_id"); ok {
		paramMap["JobId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("execute_user_uin"); ok {
		paramMap["ExecuteUserUin"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	var respData []*wedatav20250806.JobDto
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataSqlScriptRunsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	dataList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, data := range respData {
			dataMap := map[string]interface{}{}
			if data.JobId != nil {
				dataMap["job_id"] = data.JobId
			}

			if data.JobName != nil {
				dataMap["job_name"] = data.JobName
			}

			if data.JobType != nil {
				dataMap["job_type"] = data.JobType
			}

			if data.ScriptId != nil {
				dataMap["script_id"] = data.ScriptId
			}

			jobExecutionListList := make([]map[string]interface{}, 0, len(data.JobExecutionList))
			if data.JobExecutionList != nil {
				for _, jobExecutionList := range data.JobExecutionList {
					jobExecutionListMap := map[string]interface{}{}

					if jobExecutionList.JobId != nil {
						jobExecutionListMap["job_id"] = jobExecutionList.JobId
					}

					if jobExecutionList.JobExecutionId != nil {
						jobExecutionListMap["job_execution_id"] = jobExecutionList.JobExecutionId
					}

					if jobExecutionList.JobExecutionName != nil {
						jobExecutionListMap["job_execution_name"] = jobExecutionList.JobExecutionName
					}

					if jobExecutionList.ScriptContent != nil {
						jobExecutionListMap["script_content"] = jobExecutionList.ScriptContent
					}

					if jobExecutionList.Status != nil {
						jobExecutionListMap["status"] = jobExecutionList.Status
					}

					if jobExecutionList.CreateTime != nil {
						jobExecutionListMap["create_time"] = jobExecutionList.CreateTime
					}

					if jobExecutionList.ExecuteStageInfo != nil {
						jobExecutionListMap["execute_stage_info"] = jobExecutionList.ExecuteStageInfo
					}

					if jobExecutionList.LogFilePath != nil {
						jobExecutionListMap["log_file_path"] = jobExecutionList.LogFilePath
					}

					if jobExecutionList.ResultFilePath != nil {
						jobExecutionListMap["result_file_path"] = jobExecutionList.ResultFilePath
					}

					if jobExecutionList.ResultPreviewFilePath != nil {
						jobExecutionListMap["result_preview_file_path"] = jobExecutionList.ResultPreviewFilePath
					}

					if jobExecutionList.ResultTotalCount != nil {
						jobExecutionListMap["result_total_count"] = jobExecutionList.ResultTotalCount
					}

					if jobExecutionList.UpdateTime != nil {
						jobExecutionListMap["update_time"] = jobExecutionList.UpdateTime
					}

					if jobExecutionList.EndTime != nil {
						jobExecutionListMap["end_time"] = jobExecutionList.EndTime
					}

					if jobExecutionList.TimeCost != nil {
						jobExecutionListMap["time_cost"] = jobExecutionList.TimeCost
					}

					if jobExecutionList.ContextScriptContent != nil {
						jobExecutionListMap["context_script_content"] = jobExecutionList.ContextScriptContent
					}

					if jobExecutionList.ResultPreviewCount != nil {
						jobExecutionListMap["result_preview_count"] = jobExecutionList.ResultPreviewCount
					}

					if jobExecutionList.ResultEffectCount != nil {
						jobExecutionListMap["result_effect_count"] = jobExecutionList.ResultEffectCount
					}

					if jobExecutionList.CollectingTotalResult != nil {
						jobExecutionListMap["collecting_total_result"] = jobExecutionList.CollectingTotalResult
					}

					if jobExecutionList.ScriptContentTruncate != nil {
						jobExecutionListMap["script_content_truncate"] = jobExecutionList.ScriptContentTruncate
					}

					jobExecutionListList = append(jobExecutionListList, jobExecutionListMap)
				}

				dataMap["job_execution_list"] = jobExecutionListList
			}
			if data.ScriptContent != nil {
				dataMap["script_content"] = data.ScriptContent
			}

			if data.Status != nil {
				dataMap["status"] = data.Status
			}

			if data.CreateTime != nil {
				dataMap["create_time"] = data.CreateTime
			}

			if data.UpdateTime != nil {
				dataMap["update_time"] = data.UpdateTime
			}

			if data.EndTime != nil {
				dataMap["end_time"] = data.EndTime
			}

			if data.OwnerUin != nil {
				dataMap["owner_uin"] = data.OwnerUin
			}

			if data.UserUin != nil {
				dataMap["user_uin"] = data.UserUin
			}

			if data.TimeCost != nil {
				dataMap["time_cost"] = data.TimeCost
			}

			if data.ScriptContentTruncate != nil {
				dataMap["script_content_truncate"] = data.ScriptContentTruncate
			}

			dataList = append(dataList, dataMap)
		}

		_ = d.Set("data", dataList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataList); e != nil {
			return e
		}
	}

	return nil
}
