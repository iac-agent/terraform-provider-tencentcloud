package es

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	elasticsearch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudElasticsearchDiagnose() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudElasticsearchDiagnoseRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID",
			},

			"date": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Report date，格式 20210301。",
			},

			"limit": {
				Optional:    true,
				Default:     10,
				Type:        schema.TypeInt,
				Description: "数量 copies 返回 在 报告. 默认值 1。",
			},

			"diagnose_results": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 diagnostic reports。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID",
						},
						"request_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Request ID。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"completed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否diagnosis 是 完整 或 不。",
						},
						"score": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total diagnostic score。",
						},
						"job_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Diagnosis 类型，2 timing diagnosis，3 customer manual 触发器 diagnosis。",
						},
						"job_param": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Diagnostic 参数 such 作为 diagnostic 时间，diagnostic 索引，etc。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"jobs": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Diagnostic item 列表。",
									},
									"indices": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Diagnostic indices。",
									},
									"interval": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Historical diagnosis 时间。",
									},
								},
							},
						},
						"job_results": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Diagnostic item 结果 列表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"job_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Diagnostic item 名称",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Diagnostic item 状态:-2 failed,-1 到 是 retried，0 running，1 successful。",
									},
									"score": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Diagnostic item score。",
									},
									"summary": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Diagnostic summary。",
									},
									"advise": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Diagnostic advice。",
									},
									"detail": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Diagnosis details。",
									},
									"metric_details": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Details 的 diagnostic metrics。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Metric detail 名称",
												},
												"metrics": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Metric detail 值",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"dimensions": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "索引 dimension family。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"key": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Intelligent operation 和 maintenance 索引 dimension 键",
																		},
																		"value": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Dimension 值 的 intelligent operation 和 maintenance 索引",
																		},
																	},
																},
															},
															"value": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "值",
															},
														},
													},
												},
											},
										},
									},
									"log_details": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Diagnostic 日志 details。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Log exception 名称",
												},
												"advise": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Log exception handling recommendation。",
												},
												"count": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "数量 occurrences 的 日志 exception names。",
												},
											},
										},
									},
									"setting_details": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Diagnostic 配置 detail。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "键",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "值",
												},
												"advise": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Configuration processing recommendations。",
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

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudElasticsearchDiagnoseRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_elasticsearch_diagnose.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = v.(string)
	}

	if v, ok := d.GetOk("date"); ok {
		paramMap["Date"] = v.(string)
	}

	if v, ok := d.GetOk("limit"); ok {
		paramMap["Limit"] = v.(int)
	}

	service := ElasticsearchService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var diagnoseResults []*elasticsearch.DiagnoseResult

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeElasticsearchDiagnoseByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		diagnoseResults = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(diagnoseResults))
	tmpList := make([]map[string]interface{}, 0, len(diagnoseResults))

	if diagnoseResults != nil {
		for _, diagnoseResult := range diagnoseResults {
			diagnoseResultMap := map[string]interface{}{}

			if diagnoseResult.InstanceId != nil {
				diagnoseResultMap["instance_id"] = diagnoseResult.InstanceId
			}

			if diagnoseResult.RequestId != nil {
				diagnoseResultMap["request_id"] = diagnoseResult.RequestId
			}

			if diagnoseResult.CreateTime != nil {
				diagnoseResultMap["create_time"] = diagnoseResult.CreateTime
			}

			if diagnoseResult.Completed != nil {
				diagnoseResultMap["completed"] = diagnoseResult.Completed
			}

			if diagnoseResult.Score != nil {
				diagnoseResultMap["score"] = diagnoseResult.Score
			}

			if diagnoseResult.JobType != nil {
				diagnoseResultMap["job_type"] = diagnoseResult.JobType
			}

			if diagnoseResult.JobParam != nil {
				jobParamMap := map[string]interface{}{}

				if diagnoseResult.JobParam.Jobs != nil {
					jobParamMap["jobs"] = diagnoseResult.JobParam.Jobs
				}

				if diagnoseResult.JobParam.Indices != nil {
					jobParamMap["indices"] = diagnoseResult.JobParam.Indices
				}

				if diagnoseResult.JobParam.Interval != nil {
					jobParamMap["interval"] = diagnoseResult.JobParam.Interval
				}

				diagnoseResultMap["job_param"] = []interface{}{jobParamMap}
			}

			if diagnoseResult.JobResults != nil {
				jobResultsList := []interface{}{}
				for _, jobResults := range diagnoseResult.JobResults {
					jobResultsMap := map[string]interface{}{}

					if jobResults.JobName != nil {
						jobResultsMap["job_name"] = jobResults.JobName
					}

					if jobResults.Status != nil {
						jobResultsMap["status"] = jobResults.Status
					}

					if jobResults.Score != nil {
						jobResultsMap["score"] = jobResults.Score
					}

					if jobResults.Summary != nil {
						jobResultsMap["summary"] = jobResults.Summary
					}

					if jobResults.Advise != nil {
						jobResultsMap["advise"] = jobResults.Advise
					}

					if jobResults.Detail != nil {
						jobResultsMap["detail"] = jobResults.Detail
					}

					if jobResults.MetricDetails != nil {
						metricDetailsList := []interface{}{}
						for _, metricDetails := range jobResults.MetricDetails {
							metricDetailsMap := map[string]interface{}{}

							if metricDetails.Key != nil {
								metricDetailsMap["key"] = metricDetails.Key
							}

							if metricDetails.Metrics != nil {
								metricsList := []interface{}{}
								for _, metrics := range metricDetails.Metrics {
									metricsMap := map[string]interface{}{}

									if metrics.Dimensions != nil {
										dimensionsList := []interface{}{}
										for _, dimensions := range metrics.Dimensions {
											dimensionsMap := map[string]interface{}{}

											if dimensions.Key != nil {
												dimensionsMap["key"] = dimensions.Key
											}

											if dimensions.Value != nil {
												dimensionsMap["value"] = dimensions.Value
											}

											dimensionsList = append(dimensionsList, dimensionsMap)
										}

										metricsMap["dimensions"] = dimensionsList
									}

									if metrics.Value != nil {
										metricsMap["value"] = metrics.Value
									}

									metricsList = append(metricsList, metricsMap)
								}

								metricDetailsMap["metrics"] = metricsList
							}

							metricDetailsList = append(metricDetailsList, metricDetailsMap)
						}

						jobResultsMap["metric_details"] = metricDetailsList
					}

					if jobResults.LogDetails != nil {
						logDetailsList := []interface{}{}
						for _, logDetails := range jobResults.LogDetails {
							logDetailsMap := map[string]interface{}{}

							if logDetails.Key != nil {
								logDetailsMap["key"] = logDetails.Key
							}

							if logDetails.Advise != nil {
								logDetailsMap["advise"] = logDetails.Advise
							}

							if logDetails.Count != nil {
								logDetailsMap["count"] = logDetails.Count
							}

							logDetailsList = append(logDetailsList, logDetailsMap)
						}

						jobResultsMap["log_details"] = logDetailsList
					}

					if jobResults.SettingDetails != nil {
						settingDetailsList := []interface{}{}
						for _, settingDetails := range jobResults.SettingDetails {
							settingDetailsMap := map[string]interface{}{}

							if settingDetails.Key != nil {
								settingDetailsMap["key"] = settingDetails.Key
							}

							if settingDetails.Value != nil {
								settingDetailsMap["value"] = settingDetails.Value
							}

							if settingDetails.Advise != nil {
								settingDetailsMap["advise"] = settingDetails.Advise
							}

							settingDetailsList = append(settingDetailsList, settingDetailsMap)
						}

						jobResultsMap["setting_details"] = settingDetailsList
					}

					jobResultsList = append(jobResultsList, jobResultsMap)
				}

				diagnoseResultMap["job_results"] = jobResultsList
			}

			ids = append(ids, *diagnoseResult.InstanceId)
			tmpList = append(tmpList, diagnoseResultMap)
		}

		_ = d.Set("diagnose_results", tmpList)
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
