package dts

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDtsCompareTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDtsCompareTasksRead,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "作业 ID",
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "compare 任务 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "作业 ID",
						},
						"compare_task_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "compare 任务 ID",
						},
						"task_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "compare 任务 名称",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "compare 任务 状态，可选 值 是 创建/readyRun/running/success/stopping/failed/canceled。",
						},
						"config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "配置",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"object_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "对象 模式",
									},
									"object_items": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "对象 items。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"db_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "数据库 名称",
												},
												"db_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "数据库 模式",
												},
												"schema_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "schema 名称",
												},
												"table_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "表 模式",
												},
												"tables": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "表 列表。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"table_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "表 名称",
															},
														},
													},
												},
												"view_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "view 模式",
												},
												"views": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "view 列表。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"view_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "view 名称",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"check_process": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "compare check info。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "状态",
									},
									"percent": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "progress info。",
									},
									"step_all": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "all step counts。",
									},
									"step_now": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "当前 step 数量。",
									},
									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "消息",
									},
									"step": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "step info。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"step_no": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "step 数量。",
												},
												"step_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "step 名称",
												},
												"step_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "step ID。",
												},
												"status": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "状态",
												},
												"start_time": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "开始时间。",
												},
												"step_message": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "step 消息",
												},
												"percent": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "step progress。",
												},
												"errors": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "errors info。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"message": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "消息",
															},
															"solution": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "解决方案",
															},
															"help_doc": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "帮助文档",
															},
														},
													},
												},
												"warnings": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "warnings info。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"message": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "消息",
															},
															"solution": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "解决方案",
															},
															"help_doc": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "帮助文档",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"compare_process": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "compare processing info。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "状态",
									},
									"percent": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "progress info。",
									},
									"step_all": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "all step counts。",
									},
									"step_now": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "当前 step 数量。",
									},
									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "消息",
									},
									"step": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "step info。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"step_no": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "step 数量。",
												},
												"step_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "step 名称",
												},
												"step_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "step ID。",
												},
												"status": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "状态",
												},
												"start_time": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "开始时间。",
												},
												"step_message": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "step 消息",
												},
												"percent": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "step progress。",
												},
												"errors": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "errors info。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"message": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "消息",
															},
															"solution": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "解决方案",
															},
															"help_doc": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "帮助文档",
															},
														},
													},
												},
												"warnings": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "warnings info。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"message": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "消息",
															},
															"solution": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "解决方案",
															},
															"help_doc": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "帮助文档",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"conclusion": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "结论",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"started_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "开始时间。",
						},
						"finished_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "finished 时间。",
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

func dataSourceTencentCloudDtsCompareTasksRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dts_compare_tasks.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("job_id"); ok {
		paramMap["job_id"] = helper.String(v.(string))
	}

	dtsService := DtsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var compareTaskItems []*dts.CompareTaskItem
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := dtsService.DescribeDtsCompareTasksByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		compareTaskItems = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Dts items failed, reason:%+v", logId, err)
		return err
	}

	ids := make([]string, 0, len(compareTaskItems))
	itemList := make([]map[string]interface{}, 0, len(compareTaskItems))

	if compareTaskItems != nil {
		for _, item := range compareTaskItems {
			itemMap := map[string]interface{}{}
			if item.JobId != nil {
				itemMap["job_id"] = item.JobId
			}
			if item.CompareTaskId != nil {
				itemMap["compare_task_id"] = item.CompareTaskId
			}
			if item.TaskName != nil {
				itemMap["task_name"] = item.TaskName
			}
			if item.Status != nil {
				itemMap["status"] = item.Status
			}
			if item.Config != nil {
				configMap := map[string]interface{}{}
				if item.Config.ObjectMode != nil {
					configMap["object_mode"] = item.Config.ObjectMode
				}
				if item.Config.ObjectItems != nil {
					objectItemsList := []interface{}{}
					for _, objectItems := range item.Config.ObjectItems {
						objectItemsMap := map[string]interface{}{}
						if objectItems.DbName != nil {
							objectItemsMap["db_name"] = objectItems.DbName
						}
						if objectItems.DbMode != nil {
							objectItemsMap["db_mode"] = objectItems.DbMode
						}
						if objectItems.SchemaName != nil {
							objectItemsMap["schema_name"] = objectItems.SchemaName
						}
						if objectItems.TableMode != nil {
							objectItemsMap["table_mode"] = objectItems.TableMode
						}
						if objectItems.Tables != nil {
							tablesList := []interface{}{}
							for _, tables := range objectItems.Tables {
								tablesMap := map[string]interface{}{}
								if tables.TableName != nil {
									tablesMap["table_name"] = tables.TableName
								}

								tablesList = append(tablesList, tablesMap)
							}
							objectItemsMap["tables"] = tablesList
						}
						if objectItems.ViewMode != nil {
							objectItemsMap["view_mode"] = objectItems.ViewMode
						}
						if objectItems.Views != nil {
							viewsList := []interface{}{}
							for _, views := range objectItems.Views {
								viewsMap := map[string]interface{}{}
								if views.ViewName != nil {
									viewsMap["view_name"] = views.ViewName
								}

								viewsList = append(viewsList, viewsMap)
							}
							objectItemsMap["views"] = viewsList
						}

						objectItemsList = append(objectItemsList, objectItemsMap)
					}
					configMap["object_items"] = objectItemsList
				}

				itemMap["config"] = []interface{}{configMap}
			}
			if item.CheckProcess != nil {
				checkProcessMap := map[string]interface{}{}
				if item.CheckProcess.Status != nil {
					checkProcessMap["status"] = item.CheckProcess.Status
				}
				if item.CheckProcess.Percent != nil {
					checkProcessMap["percent"] = item.CheckProcess.Percent
				}
				if item.CheckProcess.StepAll != nil {
					checkProcessMap["step_all"] = item.CheckProcess.StepAll
				}
				if item.CheckProcess.StepNow != nil {
					checkProcessMap["step_now"] = item.CheckProcess.StepNow
				}
				if item.CheckProcess.Message != nil {
					checkProcessMap["message"] = item.CheckProcess.Message
				}
				if item.CheckProcess.Steps != nil {
					stepList := []interface{}{}
					for _, step := range item.CheckProcess.Steps {
						stepMap := map[string]interface{}{}
						if step.StepNo != nil {
							stepMap["step_no"] = step.StepNo
						}
						if step.StepName != nil {
							stepMap["step_name"] = step.StepName
						}
						if step.StepId != nil {
							stepMap["step_id"] = step.StepId
						}
						if step.Status != nil {
							stepMap["status"] = step.Status
						}
						if step.StartTime != nil {
							stepMap["start_time"] = step.StartTime
						}
						if step.StepMessage != nil {
							stepMap["step_message"] = step.StepMessage
						}
						if step.Percent != nil {
							stepMap["percent"] = step.Percent
						}
						if step.Errors != nil {
							errorsList := []interface{}{}
							for _, errors := range step.Errors {
								errorsMap := map[string]interface{}{}
								if errors.Message != nil {
									errorsMap["message"] = errors.Message
								}
								if errors.Solution != nil {
									errorsMap["solution"] = errors.Solution
								}
								if errors.HelpDoc != nil {
									errorsMap["help_doc"] = errors.HelpDoc
								}

								errorsList = append(errorsList, errorsMap)
							}
							stepMap["errors"] = errorsList
						}
						if step.Warnings != nil {
							warningsList := []interface{}{}
							for _, warnings := range step.Warnings {
								warningsMap := map[string]interface{}{}
								if warnings.Message != nil {
									warningsMap["message"] = warnings.Message
								}
								if warnings.Solution != nil {
									warningsMap["solution"] = warnings.Solution
								}
								if warnings.HelpDoc != nil {
									warningsMap["help_doc"] = warnings.HelpDoc
								}

								warningsList = append(warningsList, warningsMap)
							}
							stepMap["warnings"] = warningsList
						}

						stepList = append(stepList, stepMap)
					}
					checkProcessMap["step"] = stepList
				}

				itemMap["check_process"] = []interface{}{checkProcessMap}
			}
			if item.CompareProcess != nil {
				compareProcessMap := map[string]interface{}{}
				if item.CompareProcess.Status != nil {
					compareProcessMap["status"] = item.CompareProcess.Status
				}
				if item.CompareProcess.Percent != nil {
					compareProcessMap["percent"] = item.CompareProcess.Percent
				}
				if item.CompareProcess.StepAll != nil {
					compareProcessMap["step_all"] = item.CompareProcess.StepAll
				}
				if item.CompareProcess.StepNow != nil {
					compareProcessMap["step_now"] = item.CompareProcess.StepNow
				}
				if item.CompareProcess.Message != nil {
					compareProcessMap["message"] = item.CompareProcess.Message
				}
				if item.CompareProcess.Steps != nil {
					stepList := []interface{}{}
					for _, step := range item.CompareProcess.Steps {
						stepMap := map[string]interface{}{}
						if step.StepNo != nil {
							stepMap["step_no"] = step.StepNo
						}
						if step.StepName != nil {
							stepMap["step_name"] = step.StepName
						}
						if step.StepId != nil {
							stepMap["step_id"] = step.StepId
						}
						if step.Status != nil {
							stepMap["status"] = step.Status
						}
						if step.StartTime != nil {
							stepMap["start_time"] = step.StartTime
						}
						if step.StepMessage != nil {
							stepMap["step_message"] = step.StepMessage
						}
						if step.Percent != nil {
							stepMap["percent"] = step.Percent
						}
						if step.Errors != nil {
							errorsList := []interface{}{}
							for _, errors := range step.Errors {
								errorsMap := map[string]interface{}{}
								if errors.Message != nil {
									errorsMap["message"] = errors.Message
								}
								if errors.Solution != nil {
									errorsMap["solution"] = errors.Solution
								}
								if errors.HelpDoc != nil {
									errorsMap["help_doc"] = errors.HelpDoc
								}

								errorsList = append(errorsList, errorsMap)
							}
							stepMap["errors"] = errorsList
						}
						if step.Warnings != nil {
							warningsList := []interface{}{}
							for _, warnings := range step.Warnings {
								warningsMap := map[string]interface{}{}
								if warnings.Message != nil {
									warningsMap["message"] = warnings.Message
								}
								if warnings.Solution != nil {
									warningsMap["solution"] = warnings.Solution
								}
								if warnings.HelpDoc != nil {
									warningsMap["help_doc"] = warnings.HelpDoc
								}

								warningsList = append(warningsList, warningsMap)
							}
							stepMap["warnings"] = warningsList
						}

						stepList = append(stepList, stepMap)
					}
					compareProcessMap["step"] = stepList
				}

				itemMap["compare_process"] = []interface{}{compareProcessMap}
			}
			if item.Conclusion != nil {
				itemMap["conclusion"] = item.Conclusion
			}
			if item.CreatedAt != nil {
				itemMap["created_at"] = item.CreatedAt
			}
			if item.StartedAt != nil {
				itemMap["started_at"] = item.StartedAt
			}
			if item.FinishedAt != nil {
				itemMap["finished_at"] = item.FinishedAt
			}
			ids = append(ids, *item.JobId)
			itemList = append(itemList, itemMap)
		}
		d.SetId(helper.DataResourceIdsHash(ids))
		_ = d.Set("list", itemList)
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), itemList); e != nil {
			return e
		}
	}

	return nil
}
