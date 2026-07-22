package dlc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlcv20210125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDlcTaskResult() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcTaskResultRead,
		Schema: map[string]*schema.Schema{
			"task_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique 任务 ID。",
			},

			"next_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "pagination 信息 返回 通过 last response. 此 参数 可以 是 omitted 对于 first response，其中 数据 将 是 返回 从 beginning. 数据 使用 卷 集合 通过 `MaxResults` 字段 是 返回 each 时间。",
			},

			"max_results": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "最大returned rows. 取值范围：0-1,000. 默认值：1,000。",
			},

			"is_transform_data_type": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "是否convert 数据 类型",
			},

			"task_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "queried 任务 信息. 如果 返回 值 是 空， 任务 使用 entered 任务 ID does 不 exist. 任务 结果 将 是 返回 仅 如果 任务 状态 是 `2` (succeeded).\n注意：此字段可能返回 null，表示无法获取有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique 任务 ID。",
						},
						"datasource_connection_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "名称 默认值 selected 数据 来源 当 当前 作业 是 executed\n注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"database_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "名称 默认值 selected 数据库 当 当前 作业 是 executed\n注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"sql": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "currently executed SQL statement. Each 任务 包含one SQL statement。",
						},
						"sql_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "类型 executed 任务. 有效值：`DDL`，`DML`，`DQL`。",
						},
						"state": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "u200cThe 当前 任务 状态 有效值：`0` (initializing)，`1` (executing)，`2` (executed)，`3` (writing 数据)，`4` (queuing)，u200c`-1` (failed)，和 `-3` (canceled). Only 当 任务 是 successfully executed， 任务 execution 结果 将 是 返回。",
						},
						"data_amount": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Amount 的 数据 scanned 在 bytes。",
						},
						"used_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "compute 时间 在 ms。",
						},
						"output_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "地址 的 COS 存储桶 对于 storing 任务 结果",
						},
						"create_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "任务 creation 时间戳。",
						},
						"output_message": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "任务 execution 信息. `success` 将 是 返回 如果 任务 succeeds; otherwise， failure cause 将 是 返回。",
						},
						"row_affect_info": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "数量 affected rows。",
						},
						"result_schema": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Schema 信息 的 结果\n注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Column 名称，其中 是 case-insensitive 和 可以 contain up 到 25 字符。",
									},
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Column 类型 有效 值:\nstring|tinyint|smallint|int|bigint|布尔值|float|double|decimal|时间戳|date|binary|数组<data_type>|map<primitive_type，data_type>|struct<col_name : data_type [COMMENT col_comment]，...>|uniontype<data_type，data_type，...>。",
									},
									"comment": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Class 注释\n注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"precision": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Length 的 entire numeric 值\n注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"scale": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Length 的 decimal part\n注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"nullable": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "是否column 是 null.\n注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"position": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Field position\n注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"create_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Field 创建时间\n注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"modified_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Field 修改时间\n注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"is_partition": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "是否column 是 分区 字段.\n注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"result_set": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "结果 信息. After 它 是 unescaped，each element 的 outer 数组 是 数据 row.\n注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"next_token": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Pagination 信息. 如果 there 是 无 more 结果 数据，`nextToken` 将 是 空。",
						},
						"percentage": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "任务 progress (%)。",
						},
						"progress_detail": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "任务 progress details。",
						},
						"display_format": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Console display 格式 有效值：`表`，`text`。",
						},
						"total_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "任务 时间 在 ms。",
						},
						"query_result_time": {
							Type:        schema.TypeFloat,
							Required:    true,
							Description: "Time consumed 到 get results\n注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudDlcTaskResultRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_task_result.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		taskId  string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("task_id"); ok {
		paramMap["TaskId"] = helper.String(v.(string))
		taskId = v.(string)
	}

	if v, ok := d.GetOk("next_token"); ok {
		paramMap["NextToken"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("max_results"); ok {
		paramMap["MaxResults"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("is_transform_data_type"); ok {
		paramMap["IsTransformDataType"] = helper.Bool(v.(bool))
	}

	var respData *dlcv20210125.DescribeTaskResultResponseParams
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcTaskResultByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	taskInfoMap := map[string]interface{}{}
	if respData.TaskInfo != nil {
		if respData.TaskInfo.TaskId != nil {
			taskInfoMap["task_id"] = respData.TaskInfo.TaskId
		}

		if respData.TaskInfo.DatasourceConnectionName != nil {
			taskInfoMap["datasource_connection_name"] = respData.TaskInfo.DatasourceConnectionName
		}

		if respData.TaskInfo.DatabaseName != nil {
			taskInfoMap["database_name"] = respData.TaskInfo.DatabaseName
		}

		if respData.TaskInfo.SQL != nil {
			taskInfoMap["sql"] = respData.TaskInfo.SQL
		}

		if respData.TaskInfo.SQLType != nil {
			taskInfoMap["sql_type"] = respData.TaskInfo.SQLType
		}

		if respData.TaskInfo.State != nil {
			taskInfoMap["state"] = respData.TaskInfo.State
		}

		if respData.TaskInfo.DataAmount != nil {
			taskInfoMap["data_amount"] = respData.TaskInfo.DataAmount
		}

		if respData.TaskInfo.UsedTime != nil {
			taskInfoMap["used_time"] = respData.TaskInfo.UsedTime
		}

		if respData.TaskInfo.OutputPath != nil {
			taskInfoMap["output_path"] = respData.TaskInfo.OutputPath
		}

		if respData.TaskInfo.CreateTime != nil {
			taskInfoMap["create_time"] = respData.TaskInfo.CreateTime
		}

		if respData.TaskInfo.OutputMessage != nil {
			taskInfoMap["output_message"] = respData.TaskInfo.OutputMessage
		}

		if respData.TaskInfo.RowAffectInfo != nil {
			taskInfoMap["row_affect_info"] = respData.TaskInfo.RowAffectInfo
		}

		resultSchemaList := make([]map[string]interface{}, 0, len(respData.TaskInfo.ResultSchema))
		if respData.TaskInfo.ResultSchema != nil {
			for _, resultSchema := range respData.TaskInfo.ResultSchema {
				resultSchemaMap := map[string]interface{}{}

				if resultSchema.Name != nil {
					resultSchemaMap["name"] = resultSchema.Name
				}

				if resultSchema.Type != nil {
					resultSchemaMap["type"] = resultSchema.Type
				}

				if resultSchema.Comment != nil {
					resultSchemaMap["comment"] = resultSchema.Comment
				}

				if resultSchema.Precision != nil {
					resultSchemaMap["precision"] = resultSchema.Precision
				}

				if resultSchema.Scale != nil {
					resultSchemaMap["scale"] = resultSchema.Scale
				}

				if resultSchema.Nullable != nil {
					resultSchemaMap["nullable"] = resultSchema.Nullable
				}

				if resultSchema.Position != nil {
					resultSchemaMap["position"] = resultSchema.Position
				}

				if resultSchema.CreateTime != nil {
					resultSchemaMap["create_time"] = resultSchema.CreateTime
				}

				if resultSchema.ModifiedTime != nil {
					resultSchemaMap["modified_time"] = resultSchema.ModifiedTime
				}

				if resultSchema.IsPartition != nil {
					resultSchemaMap["is_partition"] = resultSchema.IsPartition
				}

				resultSchemaList = append(resultSchemaList, resultSchemaMap)
			}

			taskInfoMap["result_schema"] = resultSchemaList
		}
		if respData.TaskInfo.ResultSet != nil {
			taskInfoMap["result_set"] = respData.TaskInfo.ResultSet
		}

		if respData.TaskInfo.NextToken != nil {
			taskInfoMap["next_token"] = respData.TaskInfo.NextToken
		}

		if respData.TaskInfo.Percentage != nil {
			taskInfoMap["percentage"] = respData.TaskInfo.Percentage
		}

		if respData.TaskInfo.ProgressDetail != nil {
			taskInfoMap["progress_detail"] = respData.TaskInfo.ProgressDetail
		}

		if respData.TaskInfo.DisplayFormat != nil {
			taskInfoMap["display_format"] = respData.TaskInfo.DisplayFormat
		}

		if respData.TaskInfo.TotalTime != nil {
			taskInfoMap["total_time"] = respData.TaskInfo.TotalTime
		}

		if respData.TaskInfo.QueryResultTime != nil {
			taskInfoMap["query_result_time"] = respData.TaskInfo.QueryResultTime
		}

		_ = d.Set("task_info", []interface{}{taskInfoMap})
	}

	d.SetId(taskId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), taskInfoMap); e != nil {
			return e
		}
	}

	return nil
}
