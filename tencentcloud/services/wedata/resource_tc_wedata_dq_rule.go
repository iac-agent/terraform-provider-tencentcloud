package wedata

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataDqRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataDqRuleCreate,
		Read:   resourceTencentCloudWedataDqRuleRead,
		Update: resourceTencentCloudWedataDqRuleUpdate,
		Delete: resourceTencentCloudWedataDqRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "项目 ID",
			},
			"rule_group_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Rule 组 ID。",
			},
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Rule 名称",
			},
			"table_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Table ID。",
			},
			"rule_template_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Rule template ID。",
			},
			"type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Rule 类型 1. System 模板，2. Custom 模板，3. Custom SQL。",
			},
			"quality_dim": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Rules belong 到 quality dimensions (1. accuracy，2. uniqueness，3. completeness，4. consistency，5. timeliness，6. effectiveness)。",
			},
			"source_object_data_type_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "来源 字段 类型 int，字符串。",
			},
			"source_object_value": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "来源 字段 名称",
			},
			"condition_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Detection 范围 1. Full Table 2. Conditional scan。",
			},
			"condition_expression": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Condition scans WHERE condition expressions。",
			},
			"custom_sql": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Custom sql。",
			},
			"compare_rule": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Alarm 触发器 condition。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"items": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Comparison condition list注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"compare_type": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "比较类型 1. Fixed 值 2. Fluctuating 值 3. Comparison 的 值 范围 4. Enumeration 范围 comparison 5. Do 不 compare注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"operator": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Comparison 操作类型 &amp;lt; &amp;lt;= == =&amp;gt; &amp;gt;注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"value_compute_type": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Quality 统计 类型 1. Absolute 值 2. Increase 3. Decrease 4. C 包含5. N C does 不 contain注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"value_list": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Compare 阈值 list注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value_type": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Threshold 类型 1. Low 阈值 2. High 阈值 3. Common 阈值 4. Enumerated value注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Threshold value注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
								},
							},
						},
						"cycle_step": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Periodic 表示default 周期 的 template，在 seconds注意：此字段可能返回 null，表示无法获取有效值。",
						},
					},
				},
			},
			"alarm_level": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Alarm 触发器 levels 1. Low，2. Medium，3. High。",
			},
			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Rule 描述",
			},
			//"datasource_id": {
			//	Required:    true,
			//	Type:        schema.TypeString,
			//	Description: "Datasource ID.",
			//},
			//"database_id": {
			//	Required:    true,
			//	Type:        schema.TypeString,
			//	Description: "Database ID.",
			//},
			"target_database_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Target 数据库 ID。",
			},
			"target_table_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Target 表 ID。",
			},
			"target_condition_expr": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Target 过滤器 condition expression。",
			},
			"rel_condition_expr": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "来源 字段 和 目标 字段 是 associated 使用 conditional 在 expression。",
			},
			"field_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Custom template sql expression 字段 replacement 参数。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"where_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Where variable注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Field key注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"field_value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Field value注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"field_data_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Field type注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"table_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Library 表 variable注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Database id注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"database_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Database name注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"table_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Table id注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"table_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Table name注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"table_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Table key注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"field_config": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Field variable注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field_key": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Field key注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"field_value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Field value注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"field_data_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Field type注意：此字段可能返回 null，表示无法获取有效值。",
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
			"target_object_value": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Target 字段 名称 CITY。",
			},
			"source_engine_types": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "列表 execution engines 支持 通过 此 规则。",
			},
			"rule_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Rule ID。",
			},
		},
	}
}

func resourceTencentCloudWedataDqRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_dq_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		request   = wedata.NewCreateRuleRequest()
		response  = wedata.NewCreateRuleResponse()
		projectId string
		ruleId    string
	)

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOkExists("rule_group_id"); ok {
		request.RuleGroupId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("table_id"); ok {
		request.TableId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("rule_template_id"); ok {
		request.RuleTemplateId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("type"); ok {
		request.Type = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("quality_dim"); ok {
		request.QualityDim = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("source_object_data_type_name"); ok {
		request.SourceObjectDataTypeName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_object_value"); ok {
		request.SourceObjectValue = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("condition_type"); ok {
		request.ConditionType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("condition_expression"); ok {
		request.ConditionExpression = helper.String(v.(string))
	}

	if v, ok := d.GetOk("custom_sql"); ok {
		request.CustomSql = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "compare_rule"); ok {
		compareRule := wedata.CompareRule{}
		if v, ok := dMap["items"]; ok {
			for _, item := range v.([]interface{}) {
				itemsMap := item.(map[string]interface{})
				compareRuleItem := wedata.CompareRuleItem{}
				if v, ok := itemsMap["compare_type"]; ok {
					compareRuleItem.CompareType = helper.IntUint64(v.(int))
				}

				if v, ok := itemsMap["operator"]; ok {
					compareRuleItem.Operator = helper.String(v.(string))
				}

				if v, ok := itemsMap["value_compute_type"]; ok {
					compareRuleItem.ValueComputeType = helper.IntUint64(v.(int))
				}

				if v, ok := itemsMap["value_list"]; ok {
					for _, item := range v.([]interface{}) {
						valueListMap := item.(map[string]interface{})
						thresholdValue := wedata.ThresholdValue{}
						if v, ok := valueListMap["value_type"]; ok {
							thresholdValue.ValueType = helper.IntUint64(v.(int))
						}

						if v, ok := valueListMap["value"]; ok {
							thresholdValue.Value = helper.String(v.(string))
						}

						compareRuleItem.ValueList = append(compareRuleItem.ValueList, &thresholdValue)
					}
				}

				compareRule.Items = append(compareRule.Items, &compareRuleItem)
			}
		}

		if v, ok := dMap["cycle_step"]; ok {
			compareRule.CycleStep = helper.IntUint64(v.(int))
		}

		request.CompareRule = &compareRule
	}

	if v, ok := d.GetOkExists("alarm_level"); ok {
		request.AlarmLevel = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	//if v, ok := d.GetOk("datasource_id"); ok {
	//	request.DatasourceId = helper.String(v.(string))
	//}
	//
	//if v, ok := d.GetOk("database_id"); ok {
	//	request.DatabaseId = helper.String(v.(string))
	//}

	if v, ok := d.GetOk("target_database_id"); ok {
		request.TargetDatabaseId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_table_id"); ok {
		request.TargetTableId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_condition_expr"); ok {
		request.TargetConditionExpr = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rel_condition_expr"); ok {
		request.RelConditionExpr = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "field_config"); ok {
		ruleFieldConfig := wedata.RuleFieldConfig{}
		if v, ok := dMap["where_config"]; ok {
			for _, item := range v.([]interface{}) {
				whereConfigMap := item.(map[string]interface{})
				fieldConfig := wedata.FieldConfig{}
				if v, ok := whereConfigMap["field_key"]; ok {
					fieldConfig.FieldKey = helper.String(v.(string))
				}

				if v, ok := whereConfigMap["field_value"]; ok {
					fieldConfig.FieldValue = helper.String(v.(string))
				}

				if v, ok := whereConfigMap["field_data_type"]; ok {
					fieldConfig.FieldDataType = helper.String(v.(string))
				}

				ruleFieldConfig.WhereConfig = append(ruleFieldConfig.WhereConfig, &fieldConfig)
			}
		}

		if v, ok := dMap["table_config"]; ok {
			for _, item := range v.([]interface{}) {
				tableConfigMap := item.(map[string]interface{})
				tableConfig := wedata.TableConfig{}
				if v, ok := tableConfigMap["database_id"]; ok {
					tableConfig.DatabaseId = helper.String(v.(string))
				}

				if v, ok := tableConfigMap["database_name"]; ok {
					tableConfig.DatabaseName = helper.String(v.(string))
				}

				if v, ok := tableConfigMap["table_id"]; ok {
					tableConfig.TableId = helper.String(v.(string))
				}

				if v, ok := tableConfigMap["table_name"]; ok {
					tableConfig.TableName = helper.String(v.(string))
				}

				if v, ok := tableConfigMap["table_key"]; ok {
					tableConfig.TableKey = helper.String(v.(string))
				}

				if v, ok := tableConfigMap["field_config"]; ok {
					for _, item := range v.([]interface{}) {
						fieldConfigMap := item.(map[string]interface{})
						fieldConfig := wedata.FieldConfig{}
						if v, ok := fieldConfigMap["field_key"]; ok {
							fieldConfig.FieldKey = helper.String(v.(string))
						}

						if v, ok := fieldConfigMap["field_value"]; ok {
							fieldConfig.FieldValue = helper.String(v.(string))
						}

						if v, ok := fieldConfigMap["field_data_type"]; ok {
							fieldConfig.FieldDataType = helper.String(v.(string))
						}

						tableConfig.FieldConfig = append(tableConfig.FieldConfig, &fieldConfig)
					}
				}

				ruleFieldConfig.TableConfig = append(ruleFieldConfig.TableConfig, &tableConfig)
			}
		}

		request.FieldConfig = &ruleFieldConfig
	}

	if v, ok := d.GetOk("target_object_value"); ok {
		request.TargetObjectValue = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_engine_types"); ok {
		sourceEngineTypesSet := v.(*schema.Set).List()
		for i := range sourceEngineTypesSet {
			sourceEngineTypes := sourceEngineTypesSet[i].(int)
			request.SourceEngineTypes = append(request.SourceEngineTypes, helper.IntUint64(sourceEngineTypes))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataClient().CreateRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("wedata dqRule not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create wedata dqRule failed, reason:%+v", logId, err)
		return err
	}

	ruleIdInt := *response.Response.Data.RuleId
	ruleId = helper.UInt64ToStr(ruleIdInt)
	d.SetId(strings.Join([]string{projectId, ruleId}, tccommon.FILED_SP))

	return resourceTencentCloudWedataDqRuleRead(d, meta)
}

func resourceTencentCloudWedataDqRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_dq_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId := idSplit[0]
	ruleId := idSplit[1]

	dqRule, err := service.DescribeWedataDqRuleById(ctx, projectId, ruleId)
	if err != nil {
		return err
	}

	if dqRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WedataDqRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("project_id", projectId)
	_ = d.Set("rule_id", ruleId)

	if dqRule.RuleGroupId != nil {
		_ = d.Set("rule_group_id", dqRule.RuleGroupId)
	}

	if dqRule.Name != nil {
		_ = d.Set("name", dqRule.Name)
	}

	if dqRule.TableId != nil {
		_ = d.Set("table_id", dqRule.TableId)
	}

	if dqRule.RuleTemplateId != nil {
		_ = d.Set("rule_template_id", dqRule.RuleTemplateId)
	}

	if dqRule.Type != nil {
		_ = d.Set("type", dqRule.Type)
	}

	if dqRule.QualityDim != nil {
		_ = d.Set("quality_dim", dqRule.QualityDim)
	}

	if dqRule.SourceObjectDataTypeName != nil {
		_ = d.Set("source_object_data_type_name", dqRule.SourceObjectDataTypeName)
	}

	if dqRule.SourceObjectValue != nil {
		_ = d.Set("source_object_value", dqRule.SourceObjectValue)
	}

	if dqRule.ConditionType != nil {
		_ = d.Set("condition_type", dqRule.ConditionType)
	}

	if dqRule.ConditionExpression != nil {
		_ = d.Set("condition_expression", dqRule.ConditionExpression)
	}

	if dqRule.CustomSql != nil {
		_ = d.Set("custom_sql", dqRule.CustomSql)
	}

	if dqRule.CompareRule != nil {
		compareRuleMap := map[string]interface{}{}
		if dqRule.CompareRule.Items != nil {
			itemsList := []interface{}{}
			for _, items := range dqRule.CompareRule.Items {
				itemsMap := map[string]interface{}{}

				if items.CompareType != nil {
					itemsMap["compare_type"] = items.CompareType
				}

				if items.Operator != nil {
					itemsMap["operator"] = items.Operator
				}

				if items.ValueComputeType != nil {
					itemsMap["value_compute_type"] = items.ValueComputeType
				}

				if items.ValueList != nil {
					valueListList := []interface{}{}
					for _, valueList := range items.ValueList {
						valueListMap := map[string]interface{}{}

						if valueList.ValueType != nil {
							valueListMap["value_type"] = valueList.ValueType
						}

						if valueList.Value != nil {
							valueListMap["value"] = valueList.Value
						}

						valueListList = append(valueListList, valueListMap)
					}

					itemsMap["value_list"] = valueListList
				}

				itemsList = append(itemsList, itemsMap)
			}

			compareRuleMap["items"] = itemsList
		}

		if dqRule.CompareRule.CycleStep != nil {
			compareRuleMap["cycle_step"] = dqRule.CompareRule.CycleStep
		}

		_ = d.Set("compare_rule", []interface{}{compareRuleMap})
	}

	if dqRule.AlarmLevel != nil {
		_ = d.Set("alarm_level", dqRule.AlarmLevel)
	}

	if dqRule.Description != nil {
		_ = d.Set("description", dqRule.Description)
	}

	if dqRule.TargetDatabaseId != nil {
		_ = d.Set("target_database_id", dqRule.TargetDatabaseId)
	}

	if dqRule.TargetTableId != nil {
		_ = d.Set("target_table_id", dqRule.TargetTableId)
	}

	if dqRule.TargetConditionExpr != nil {
		_ = d.Set("target_condition_expr", dqRule.TargetConditionExpr)
	}

	if dqRule.RelConditionExpr != nil {
		_ = d.Set("rel_condition_expr", dqRule.RelConditionExpr)
	}

	if dqRule.FieldConfig != nil {
		fieldConfigMap := map[string]interface{}{}

		if dqRule.FieldConfig.WhereConfig != nil {
			whereConfigList := []interface{}{}
			for _, whereConfig := range dqRule.FieldConfig.WhereConfig {
				whereConfigMap := map[string]interface{}{}

				if whereConfig.FieldKey != nil {
					whereConfigMap["field_key"] = whereConfig.FieldKey
				}

				if whereConfig.FieldValue != nil {
					whereConfigMap["field_value"] = whereConfig.FieldValue
				}

				if whereConfig.FieldDataType != nil {
					whereConfigMap["field_data_type"] = whereConfig.FieldDataType
				}

				whereConfigList = append(whereConfigList, whereConfigMap)
			}

			fieldConfigMap["where_config"] = whereConfigList
		}

		if dqRule.FieldConfig.TableConfig != nil {
			tableConfigList := []interface{}{}
			for _, tableConfig := range dqRule.FieldConfig.TableConfig {
				tableConfigMap := map[string]interface{}{}

				if tableConfig.DatabaseId != nil {
					tableConfigMap["database_id"] = tableConfig.DatabaseId
				}

				if tableConfig.DatabaseName != nil {
					tableConfigMap["database_name"] = tableConfig.DatabaseName
				}

				if tableConfig.TableId != nil {
					tableConfigMap["table_id"] = tableConfig.TableId
				}

				if tableConfig.TableName != nil {
					tableConfigMap["table_name"] = tableConfig.TableName
				}

				if tableConfig.TableKey != nil {
					tableConfigMap["table_key"] = tableConfig.TableKey
				}

				if tableConfig.FieldConfig != nil {
					fieldConfigList := []interface{}{}
					for _, fieldConfig := range tableConfig.FieldConfig {
						configMap := map[string]interface{}{}

						if fieldConfig.FieldKey != nil {
							configMap["field_key"] = fieldConfig.FieldKey
						}

						if fieldConfig.FieldValue != nil {
							configMap["field_value"] = fieldConfig.FieldValue
						}

						if fieldConfig.FieldDataType != nil {
							configMap["field_data_type"] = fieldConfig.FieldDataType
						}

						fieldConfigList = append(fieldConfigList, configMap)
					}

					tableConfigMap["field_config"] = fieldConfigList
				}

				tableConfigList = append(tableConfigList, tableConfigMap)
			}

			fieldConfigMap["table_config"] = tableConfigList
		}

		_ = d.Set("field_config", []interface{}{fieldConfigMap})
	}

	if dqRule.TargetObjectValue != nil {
		_ = d.Set("target_object_value", dqRule.TargetObjectValue)
	}

	if dqRule.SourceEngineTypes != nil {
		_ = d.Set("source_engine_types", dqRule.SourceEngineTypes)
	}

	return nil
}

func resourceTencentCloudWedataDqRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_dq_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = wedata.NewModifyRuleRequest()
	)

	immutableArgs := []string{"project_id", "datasource_id", "database_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId := idSplit[0]

	request.ProjectId = &projectId

	if v, ok := d.GetOkExists("rule_group_id"); ok {
		request.RuleGroupId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("table_id"); ok {
		request.TableId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("rule_template_id"); ok {
		request.RuleTemplateId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("type"); ok {
		request.Type = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("quality_dim"); ok {
		request.QualityDim = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("source_object_data_type_name"); ok {
		request.SourceObjectDataTypeName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_object_value"); ok {
		request.SourceObjectValue = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("condition_type"); ok {
		request.ConditionType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("condition_expression"); ok {
		request.ConditionExpression = helper.String(v.(string))
	}

	if v, ok := d.GetOk("custom_sql"); ok {
		request.CustomSql = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "compare_rule"); ok {
		compareRule := wedata.CompareRule{}
		if v, ok := dMap["items"]; ok {
			for _, item := range v.([]interface{}) {
				itemsMap := item.(map[string]interface{})
				compareRuleItem := wedata.CompareRuleItem{}
				if v, ok := itemsMap["compare_type"]; ok {
					compareRuleItem.CompareType = helper.IntUint64(v.(int))
				}

				if v, ok := itemsMap["operator"]; ok {
					compareRuleItem.Operator = helper.String(v.(string))
				}

				if v, ok := itemsMap["value_compute_type"]; ok {
					compareRuleItem.ValueComputeType = helper.IntUint64(v.(int))
				}

				if v, ok := itemsMap["value_list"]; ok {
					for _, item := range v.([]interface{}) {
						valueListMap := item.(map[string]interface{})
						thresholdValue := wedata.ThresholdValue{}
						if v, ok := valueListMap["value_type"]; ok {
							thresholdValue.ValueType = helper.IntUint64(v.(int))
						}

						if v, ok := valueListMap["value"]; ok {
							thresholdValue.Value = helper.String(v.(string))
						}

						compareRuleItem.ValueList = append(compareRuleItem.ValueList, &thresholdValue)
					}
				}

				compareRule.Items = append(compareRule.Items, &compareRuleItem)
			}
		}

		if v, ok := dMap["cycle_step"]; ok {
			compareRule.CycleStep = helper.IntUint64(v.(int))
		}

		request.CompareRule = &compareRule
	}

	if v, ok := d.GetOkExists("alarm_level"); ok {
		request.AlarmLevel = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_database_id"); ok {
		request.TargetDatabaseId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_table_id"); ok {
		request.TargetTableId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_condition_expr"); ok {
		request.TargetConditionExpr = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rel_condition_expr"); ok {
		request.RelConditionExpr = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "field_config"); ok {
		ruleFieldConfig := wedata.RuleFieldConfig{}
		if v, ok := dMap["where_config"]; ok {
			for _, item := range v.([]interface{}) {
				whereConfigMap := item.(map[string]interface{})
				fieldConfig := wedata.FieldConfig{}
				if v, ok := whereConfigMap["field_key"]; ok {
					fieldConfig.FieldKey = helper.String(v.(string))
				}

				if v, ok := whereConfigMap["field_value"]; ok {
					fieldConfig.FieldValue = helper.String(v.(string))
				}

				if v, ok := whereConfigMap["field_data_type"]; ok {
					fieldConfig.FieldDataType = helper.String(v.(string))
				}

				ruleFieldConfig.WhereConfig = append(ruleFieldConfig.WhereConfig, &fieldConfig)
			}
		}

		if v, ok := dMap["table_config"]; ok {
			for _, item := range v.([]interface{}) {
				tableConfigMap := item.(map[string]interface{})
				tableConfig := wedata.TableConfig{}
				if v, ok := tableConfigMap["database_id"]; ok {
					tableConfig.DatabaseId = helper.String(v.(string))
				}

				if v, ok := tableConfigMap["database_name"]; ok {
					tableConfig.DatabaseName = helper.String(v.(string))
				}

				if v, ok := tableConfigMap["table_id"]; ok {
					tableConfig.TableId = helper.String(v.(string))
				}

				if v, ok := tableConfigMap["table_name"]; ok {
					tableConfig.TableName = helper.String(v.(string))
				}

				if v, ok := tableConfigMap["table_key"]; ok {
					tableConfig.TableKey = helper.String(v.(string))
				}

				if v, ok := tableConfigMap["field_config"]; ok {
					for _, item := range v.([]interface{}) {
						fieldConfigMap := item.(map[string]interface{})
						fieldConfig := wedata.FieldConfig{}
						if v, ok := fieldConfigMap["field_key"]; ok {
							fieldConfig.FieldKey = helper.String(v.(string))
						}

						if v, ok := fieldConfigMap["field_value"]; ok {
							fieldConfig.FieldValue = helper.String(v.(string))
						}

						if v, ok := fieldConfigMap["field_data_type"]; ok {
							fieldConfig.FieldDataType = helper.String(v.(string))
						}

						tableConfig.FieldConfig = append(tableConfig.FieldConfig, &fieldConfig)
					}
				}

				ruleFieldConfig.TableConfig = append(ruleFieldConfig.TableConfig, &tableConfig)
			}
		}

		request.FieldConfig = &ruleFieldConfig
	}

	if v, ok := d.GetOk("target_object_value"); ok {
		request.TargetObjectValue = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_engine_types"); ok {
		sourceEngineTypesSet := v.(*schema.Set).List()
		for i := range sourceEngineTypesSet {
			sourceEngineTypes := sourceEngineTypesSet[i].(int)
			request.SourceEngineTypes = append(request.SourceEngineTypes, helper.IntUint64(sourceEngineTypes))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataClient().ModifyRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update wedata dqRule failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataDqRuleRead(d, meta)
}

func resourceTencentCloudWedataDqRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_dq_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId := idSplit[0]
	ruleId := idSplit[1]

	if err := service.DeleteWedataDqRuleById(ctx, projectId, ruleId); err != nil {
		return err
	}

	return nil
}
