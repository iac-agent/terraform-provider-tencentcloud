package tat

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTatInvocationTask() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTatInvocationTaskRead,
		Schema: map[string]*schema.Schema{
			"invocation_task_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "列表 execution 任务 IDs. Up 到 100 IDs 是 allowed 对于 each 请求. InvocationTaskIds 和 Filters 不能 是 指定 在 same 时间。",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 conditions.invocation-ID - String - 必填: No - (过滤器 condition) 过滤器 通过 execution activity ID.invocation-任务-ID - String - 必填: No - (过滤器 condition) 过滤器 通过 execution 任务 ID.实例-ID - String - 必填: No - (过滤器 condition) 过滤器 通过 实例 ID.command-ID - String - 必填: No - (过滤器 condition) 过滤器 通过 命令 IDUp 到 10 Filters 是 allowed 对于 each 请求. Each 过滤器 可以 have up 到 five 过滤器.Values. InvocationTaskIds 和 Filters 不能 是 指定 在 same 时间。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "待过滤字段",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "过滤器 值 的 字段。",
						},
					},
				},
			},

			"hide_output": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "是否hide output. 有效 值:True (默认值): Hide outputFalse: Show output。",
			},

			"invocation_task_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 execution tasks。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"invocation_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Execution activity ID。",
						},
						"invocation_task_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Execution 任务 ID。",
						},
						"command_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "命令 ID",
						},
						"task_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Execution 任务 状态 有效 值:PENDING: PendingDELIVERING: DeliveringDELIVER_DELAYED: Delivery delayedDELIVER_FAILED: Delivery failedSTART_FAILED: Failed 到 start commandRUNNING: RunningSUCCESS: SuccessFAILED: Failed 到 execute command. exit 代码 是 不 0 after execution.TIMEOUT: Command timed outTASK_TIMEOUT: 任务 timed outCANCELLING: CancelingCANCELLED: Canceled (canceled before execution)TERMINATED: Terminated (canceled during execution)。",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID",
						},
						"task_result": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Execution 结果",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"exit_code": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ExitCode 的 execution。",
									},
									"output": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Base64-encoded command output. 最大 长度 是 24 KB。",
									},
									"exec_start_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Time 当 execution 是 started。",
									},
									"exec_end_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Time 当 execution 是 ended。",
									},
									"dropped": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Dropped bytes 的 command output。",
									},
									"output_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "COS URL 的 logs。",
									},
									"output_upload_cos_error_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "错误信息 对于 uploading logs 到 COS。",
									},
								},
							},
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "开始时间 的 execution 任务。",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "结束时间 的 execution 任务。",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "更新时间。",
						},
						"command_document": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Command details 的 execution 任务。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Base64-encoded command。",
									},
									"command_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "命令类型",
									},
									"timeout": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Timeout 周期",
									},
									"working_directory": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Execution 路径",
									},
									"username": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "用户 who executes command。",
									},
									"output_cos_bucket_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "URL 的 COS 存储桶 到 store output。",
									},
									"output_cos_key_prefix": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Prefix 的 输出文件 名称",
									},
								},
							},
						},
						"error_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "错误信息 displayed 当 execution 任务 fails。",
						},
						"invocation_source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Invocation 来源",
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

func dataSourceTencentCloudTatInvocationTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tat_invocation_task.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("invocation_task_ids"); ok {
		invocationTaskIdsSet := v.(*schema.Set).List()
		paramMap["InvocationTaskIds"] = helper.InterfacesStringsPoint(invocationTaskIdsSet)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*tat.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := tat.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	if v, _ := d.GetOk("hide_output"); v != nil {
		paramMap["HideOutput"] = helper.Bool(v.(bool))
	}

	service := TatService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var invocationTaskSet []*tat.InvocationTask

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTatInvocationTaskByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		invocationTaskSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(invocationTaskSet))
	tmpList := make([]map[string]interface{}, 0, len(invocationTaskSet))

	if invocationTaskSet != nil {
		for _, invocationTask := range invocationTaskSet {
			invocationTaskMap := map[string]interface{}{}

			if invocationTask.InvocationId != nil {
				invocationTaskMap["invocation_id"] = invocationTask.InvocationId
			}

			if invocationTask.InvocationTaskId != nil {
				invocationTaskMap["invocation_task_id"] = invocationTask.InvocationTaskId
			}

			if invocationTask.CommandId != nil {
				invocationTaskMap["command_id"] = invocationTask.CommandId
			}

			if invocationTask.TaskStatus != nil {
				invocationTaskMap["task_status"] = invocationTask.TaskStatus
			}

			if invocationTask.InstanceId != nil {
				invocationTaskMap["instance_id"] = invocationTask.InstanceId
			}

			if invocationTask.TaskResult != nil {
				taskResultMap := map[string]interface{}{}

				if invocationTask.TaskResult.ExitCode != nil {
					taskResultMap["exit_code"] = invocationTask.TaskResult.ExitCode
				}

				if invocationTask.TaskResult.Output != nil {
					taskResultMap["output"] = invocationTask.TaskResult.Output
				}

				if invocationTask.TaskResult.ExecStartTime != nil {
					taskResultMap["exec_start_time"] = invocationTask.TaskResult.ExecStartTime
				}

				if invocationTask.TaskResult.ExecEndTime != nil {
					taskResultMap["exec_end_time"] = invocationTask.TaskResult.ExecEndTime
				}

				if invocationTask.TaskResult.Dropped != nil {
					taskResultMap["dropped"] = invocationTask.TaskResult.Dropped
				}

				if invocationTask.TaskResult.OutputUrl != nil {
					taskResultMap["output_url"] = invocationTask.TaskResult.OutputUrl
				}

				if invocationTask.TaskResult.OutputUploadCOSErrorInfo != nil {
					taskResultMap["output_upload_cos_error_info"] = invocationTask.TaskResult.OutputUploadCOSErrorInfo
				}

				invocationTaskMap["task_result"] = []interface{}{taskResultMap}
			}

			if invocationTask.StartTime != nil {
				invocationTaskMap["start_time"] = invocationTask.StartTime
			}

			if invocationTask.EndTime != nil {
				invocationTaskMap["end_time"] = invocationTask.EndTime
			}

			if invocationTask.CreatedTime != nil {
				invocationTaskMap["created_time"] = invocationTask.CreatedTime
			}

			if invocationTask.UpdatedTime != nil {
				invocationTaskMap["updated_time"] = invocationTask.UpdatedTime
			}

			if invocationTask.CommandDocument != nil {
				commandDocumentMap := map[string]interface{}{}

				if invocationTask.CommandDocument.Content != nil {
					commandDocumentMap["content"] = invocationTask.CommandDocument.Content
				}

				if invocationTask.CommandDocument.CommandType != nil {
					commandDocumentMap["command_type"] = invocationTask.CommandDocument.CommandType
				}

				if invocationTask.CommandDocument.Timeout != nil {
					commandDocumentMap["timeout"] = invocationTask.CommandDocument.Timeout
				}

				if invocationTask.CommandDocument.WorkingDirectory != nil {
					commandDocumentMap["working_directory"] = invocationTask.CommandDocument.WorkingDirectory
				}

				if invocationTask.CommandDocument.Username != nil {
					commandDocumentMap["username"] = invocationTask.CommandDocument.Username
				}

				if invocationTask.CommandDocument.OutputCOSBucketUrl != nil {
					commandDocumentMap["output_cos_bucket_url"] = invocationTask.CommandDocument.OutputCOSBucketUrl
				}

				if invocationTask.CommandDocument.OutputCOSKeyPrefix != nil {
					commandDocumentMap["output_cos_key_prefix"] = invocationTask.CommandDocument.OutputCOSKeyPrefix
				}

				invocationTaskMap["command_document"] = []interface{}{commandDocumentMap}
			}

			if invocationTask.ErrorInfo != nil {
				invocationTaskMap["error_info"] = invocationTask.ErrorInfo
			}

			if invocationTask.InvocationSource != nil {
				invocationTaskMap["invocation_source"] = invocationTask.InvocationSource
			}

			ids = append(ids, *invocationTask.InvocationTaskId)
			tmpList = append(tmpList, invocationTaskMap)
		}

		_ = d.Set("invocation_task_set", tmpList)
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
