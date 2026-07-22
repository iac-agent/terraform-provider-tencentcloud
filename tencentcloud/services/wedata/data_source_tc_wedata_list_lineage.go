package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataListLineage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataListLineageRead,
		Schema: map[string]*schema.Schema{
			"resource_unique_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Entity 唯一 ID。",
			},

			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Entity 类型: TABLE|METRIC|MODEL|SERVICE|COLUMN。",
			},

			"direction": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Lineage direction: INPUT|OUTPUT。",
			},

			"platform": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "来源: WEDATA|THIRD，默认为 WEDATA。",
			},

			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Lineage 记录 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Current 资源。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_unique_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Entity original 唯一 ID。",
									},
									"resource_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Business 名称: 数据库.表|指标名称|model 名称|字段 名称",
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Entity 类型: TABLE|METRIC|MODEL|SERVICE|COLUMN。",
									},
									"lineage_node_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Lineage 节点 唯一 identifier。",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "描述: 表 类型|metric 描述|model 描述|字段 描述",
									},
									"platform": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "来源: WEDATA|THIRD，默认为 WEDATA。",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创建时间。",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "更新时间。",
									},
									"resource_properties": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Resource additional extension 参数。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "属性名称",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "属性值",
												},
											},
										},
									},
								},
							},
						},
						"relation": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Relation。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"relation_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Association ID。",
									},
									"source_unique_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "来源 唯一 lineage ID。",
									},
									"target_unique_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Target 唯一 lineage ID。",
									},
									"processes": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Lineage processing process。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"process_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Original 唯一 ID。",
												},
												"process_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "任务 类型: SCHEDULE_TASK，INTEGRATION_TASK，THIRD_REPORT，TABLE_MODEL，MODEL_METRIC，METRIC_METRIC，DATA_SERVICE。",
												},
												"platform": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "WEDATA，THIRD。",
												},
												"process_sub_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "任务 subtype: SQL_TASK，INTEGRATED_STREAM，INTEGRATED_OFFLINE。",
												},
												"process_properties": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Additional extension 参数。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "属性名称",
															},
															"value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "属性值",
															},
														},
													},
												},
												"lineage_node_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Lineage 任务 唯一 节点 ID",
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

func dataSourceTencentCloudWedataListLineageRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_list_lineage.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("resource_unique_id"); ok {
		paramMap["ResourceUniqueId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_type"); ok {
		paramMap["ResourceType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("direction"); ok {
		paramMap["Direction"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("platform"); ok {
		paramMap["Platform"] = helper.String(v.(string))
	}

	var respData []*wedatav20250806.LineageNodeInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataListLineageByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	itemsList := make([]map[string]interface{}, 0, len(respData))
	for _, items := range respData {
		itemsMap := map[string]interface{}{}
		resourceMap := map[string]interface{}{}
		if items.Resource != nil {
			if items.Resource.ResourceUniqueId != nil {
				resourceMap["resource_unique_id"] = items.Resource.ResourceUniqueId
			}

			if items.Resource.ResourceName != nil {
				resourceMap["resource_name"] = items.Resource.ResourceName
			}

			if items.Resource.ResourceType != nil {
				resourceMap["resource_type"] = items.Resource.ResourceType
			}

			if items.Resource.LineageNodeId != nil {
				resourceMap["lineage_node_id"] = items.Resource.LineageNodeId
			}

			if items.Resource.Description != nil {
				resourceMap["description"] = items.Resource.Description
			}

			if items.Resource.Platform != nil {
				resourceMap["platform"] = items.Resource.Platform
			}

			if items.Resource.CreateTime != nil {
				resourceMap["create_time"] = items.Resource.CreateTime
			}

			if items.Resource.UpdateTime != nil {
				resourceMap["update_time"] = items.Resource.UpdateTime
			}

			resourcePropertiesList := make([]map[string]interface{}, 0, len(items.Resource.ResourceProperties))
			if items.Resource.ResourceProperties != nil {
				for _, resourceProperties := range items.Resource.ResourceProperties {
					resourcePropertiesMap := map[string]interface{}{}

					if resourceProperties.Name != nil {
						resourcePropertiesMap["name"] = resourceProperties.Name
					}

					if resourceProperties.Value != nil {
						resourcePropertiesMap["value"] = resourceProperties.Value
					}

					resourcePropertiesList = append(resourcePropertiesList, resourcePropertiesMap)
				}

				resourceMap["resource_properties"] = resourcePropertiesList
			}

			itemsMap["resource"] = []interface{}{resourceMap}
		}

		relationMap := map[string]interface{}{}
		if items.Relation != nil {
			if items.Relation.RelationId != nil {
				relationMap["relation_id"] = items.Relation.RelationId
			}

			if items.Relation.SourceUniqueId != nil {
				relationMap["source_unique_id"] = items.Relation.SourceUniqueId
			}

			if items.Relation.TargetUniqueId != nil {
				relationMap["target_unique_id"] = items.Relation.TargetUniqueId
			}

			processesList := make([]map[string]interface{}, 0, len(items.Relation.Processes))
			if items.Relation.Processes != nil {
				for _, processes := range items.Relation.Processes {
					processesMap := map[string]interface{}{}
					if processes.ProcessId != nil {
						processesMap["process_id"] = processes.ProcessId
					}

					if processes.ProcessType != nil {
						processesMap["process_type"] = processes.ProcessType
					}

					if processes.Platform != nil {
						processesMap["platform"] = processes.Platform
					}

					if processes.ProcessSubType != nil {
						processesMap["process_sub_type"] = processes.ProcessSubType
					}

					processPropertiesList := make([]map[string]interface{}, 0, len(processes.ProcessProperties))
					if processes.ProcessProperties != nil {
						for _, processProperties := range processes.ProcessProperties {
							processPropertiesMap := map[string]interface{}{}

							if processProperties.Name != nil {
								processPropertiesMap["name"] = processProperties.Name
							}

							if processProperties.Value != nil {
								processPropertiesMap["value"] = processProperties.Value
							}

							processPropertiesList = append(processPropertiesList, processPropertiesMap)
						}

						processesMap["process_properties"] = processPropertiesList
					}

					if processes.LineageNodeId != nil {
						processesMap["lineage_node_id"] = processes.LineageNodeId
					}

					processesList = append(processesList, processesMap)
				}

				relationMap["processes"] = processesList
			}

			itemsMap["relation"] = []interface{}{relationMap}
		}

		itemsList = append(itemsList, itemsMap)
	}

	_ = d.Set("items", itemsList)

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), itemsList); e != nil {
			return e
		}
	}

	return nil
}
