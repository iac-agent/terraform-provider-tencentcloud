package wedata

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataRuleTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataRuleTemplatesRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "模板 类型 `1` 表示 System template，`2` 表示 Custom template。",
			},

			"source_object_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "来源 数据 对象 类型 `1`: Constant，`2`: Offline 表 级别，`3`: Offline 字段 级别",
			},

			"project_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "项目 ID",
			},

			"source_engine_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Applicable 类型 来源 数据。",
			},

			"data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "规则 template 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_template_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID 规则 template。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 规则 template。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 规则 template。",
						},
						"type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "模板 类型 `1` 表示 System template，`2` 表示 Custom template。",
						},
						"source_object_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "来源 对象 类型 `1`: Constant，`2`: Offline 表 级别，`3`: Offline 字段 级别",
						},
						"source_object_data_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "来源 数据 对象 类型 `1`: 值，`2`: 字符串。",
						},
						"source_content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "内容 的 规则 template。",
						},
						"source_engine_types": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "Applicable 类型 来源 数据。",
						},
						"quality_dim": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Quality inspection dimensions. `1`: Accuracy，`2`: Uniqueness，`3`: Completeness，`4`: Consistency，`5`: Timeliness，`6`: Effectiveness。",
						},
						"compare_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "类型 comparison 方法 支持 通过 规则 (1: fixed 值 comparison，greater 比，less 比，greater 比 或 equal 到，etc. 2: fluctuating 值 comparison，absolute 值，rise，fall)。",
						},
						"citation_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Citations。",
						},
						"user_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "用户 ID。",
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "用户 名称",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "更新时间，like: yyyy-MM-dd HH:mm:ss。",
						},
						"where_flag": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "如果 add 其中。",
						},
						"multi_source_flag": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否associate other 库 tables。",
						},
						"sql_expression": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sql Expression。",
						},
						"sub_quality_dim": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sub Quality inspection dimensions. `1`: Accuracy，`2`: Uniqueness，`3`: Completeness，`4`: Consistency，`5`: Timeliness，`6`: Effectiveness。",
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

func dataSourceTencentCloudWedataRuleTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_rule_templates.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOkExists("type"); ok {
		paramMap["Type"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("source_object_type"); ok {
		paramMap["SourceObjectType"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("project_id"); ok {
		paramMap["ProjectId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_engine_types"); ok {
		sourceEngineTypesSet := v.(*schema.Set).List()
		var sourceEngineTypes []*uint64
		for i := range sourceEngineTypesSet {
			sourceEngineType := sourceEngineTypesSet[i].(int)
			sourceEngineTypes = append(sourceEngineTypes, helper.IntUint64(sourceEngineType))
		}
		paramMap["SourceEngineTypes"] = sourceEngineTypes
	}

	service := WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var data []*wedata.RuleTemplate

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataRuleTemplatesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		data = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(data))
	tmpList := make([]map[string]interface{}, 0, len(data))

	if data != nil {
		for _, ruleTemplate := range data {
			ruleTemplateMap := map[string]interface{}{}

			if ruleTemplate.RuleTemplateId != nil {
				ruleTemplateMap["rule_template_id"] = ruleTemplate.RuleTemplateId
			}

			if ruleTemplate.Name != nil {
				ruleTemplateMap["name"] = ruleTemplate.Name
			}

			if ruleTemplate.Description != nil {
				ruleTemplateMap["description"] = ruleTemplate.Description
			}

			if ruleTemplate.Type != nil {
				ruleTemplateMap["type"] = ruleTemplate.Type
			}

			if ruleTemplate.SourceObjectType != nil {
				ruleTemplateMap["source_object_type"] = ruleTemplate.SourceObjectType
			}

			if ruleTemplate.SourceObjectDataType != nil {
				ruleTemplateMap["source_object_data_type"] = ruleTemplate.SourceObjectDataType
			}

			if ruleTemplate.SourceContent != nil {
				ruleTemplateMap["source_content"] = ruleTemplate.SourceContent
			}

			if ruleTemplate.SourceEngineTypes != nil {
				ruleTemplateMap["source_engine_types"] = ruleTemplate.SourceEngineTypes
			}

			if ruleTemplate.QualityDim != nil {
				ruleTemplateMap["quality_dim"] = ruleTemplate.QualityDim
			}

			if ruleTemplate.CompareType != nil {
				ruleTemplateMap["compare_type"] = ruleTemplate.CompareType
			}

			if ruleTemplate.CitationCount != nil {
				ruleTemplateMap["citation_count"] = ruleTemplate.CitationCount
			}

			if ruleTemplate.UserId != nil {
				ruleTemplateMap["user_id"] = ruleTemplate.UserId
			}

			if ruleTemplate.UserName != nil {
				ruleTemplateMap["user_name"] = ruleTemplate.UserName
			}

			if ruleTemplate.UpdateTime != nil {
				ruleTemplateMap["update_time"] = ruleTemplate.UpdateTime
			}

			if ruleTemplate.WhereFlag != nil {
				ruleTemplateMap["where_flag"] = ruleTemplate.WhereFlag
			}

			if ruleTemplate.MultiSourceFlag != nil {
				ruleTemplateMap["multi_source_flag"] = ruleTemplate.MultiSourceFlag
			}

			if ruleTemplate.SqlExpression != nil {
				ruleTemplateMap["sql_expression"] = ruleTemplate.SqlExpression
			}

			if ruleTemplate.SubQualityDim != nil {
				ruleTemplateMap["sub_quality_dim"] = ruleTemplate.SubQualityDim
			}

			ids = append(ids, helper.UInt64ToStr(*ruleTemplate.RuleTemplateId))
			tmpList = append(tmpList, ruleTemplateMap)
		}

		_ = d.Set("data", tmpList)
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
