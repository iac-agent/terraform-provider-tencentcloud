package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudConfigRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudConfigRulesRead,
		Schema: map[string]*schema.Schema{
			"rule_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule 名称 对于 filtering。",
			},

			"risk_level": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "风险等级 列表 对于 filtering. 有效值：1 (high risk)，2 (medium risk)，3 (low risk)。",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule state 对于 filtering. 有效值：ACTIVE，UN_ACTIVE。",
			},

			"compliance_result": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Compliance 结果 列表 对于 filtering. 有效值：COMPLIANT，NON_COMPLIANT。",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"order_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort 类型 通过 规则 名称 有效值：desc (descending)，asc (ascending)。",
			},

			"rule_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "配置 规则 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "配置 规则 ID。",
						},
						"identifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule identifier。",
						},
						"rule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule 名称",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule 描述",
						},
						"risk_level": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "风险等级 有效值：1 (low risk)，2 (medium risk)，3 (high risk)。",
						},
						"service_function": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Corresponding 服务 函数。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule 状态 有效值：ACTIVE，NO_ACTIVE。",
						},
						"compliance_result": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Compliance 结果 有效值：COMPLIANT，NON_COMPLIANT，NOT_APPLICABLE。",
						},
						"identifier_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule 类型 有效值：CUSTOMIZE (自定义 规则)，SYSTEM (managed 规则)。",
						},
						"compliance_pack_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Compliance pack ID。",
						},
						"compliance_pack_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Compliance pack 名称",
						},
						"resource_type": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Supported 资源类型 列表。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"labels": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Rule 标签 列表。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"config_rule_invoked_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule evaluation 时间。",
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

func dataSourceTencentCloudConfigRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_config_rules.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ConfigService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := buildConfigRulesParamMap(d)

	respData, reqErr := service.DescribeConfigRulesByFilter(ctx, paramMap)
	if reqErr != nil {
		return reqErr
	}

	ruleList := flattenConfigRuleList(respData)
	_ = d.Set("rule_list", ruleList)

	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}

func buildConfigRulesParamMap(d *schema.ResourceData) map[string]interface{} {
	paramMap := make(map[string]interface{})

	if v, ok := d.GetOk("rule_name"); ok {
		paramMap["RuleName"] = v.(string)
	}

	if v, ok := d.GetOk("risk_level"); ok {
		rawList := v.([]interface{})
		riskLevels := make([]*uint64, 0, len(rawList))
		for _, item := range rawList {
			val := uint64(item.(int))
			riskLevels = append(riskLevels, &val)
		}

		paramMap["RiskLevel"] = riskLevels
	}

	if v, ok := d.GetOk("state"); ok {
		paramMap["State"] = v.(string)
	}

	if v, ok := d.GetOk("compliance_result"); ok {
		rawList := v.([]interface{})
		results := make([]*string, 0, len(rawList))
		for _, item := range rawList {
			val := item.(string)
			results = append(results, &val)
		}

		paramMap["ComplianceResult"] = results
	}

	if v, ok := d.GetOk("order_type"); ok {
		paramMap["OrderType"] = v.(string)
	}

	return paramMap
}

func flattenConfigRuleList(items []*configv20220802.ConfigRule) []map[string]interface{} {
	ruleList := make([]map[string]interface{}, 0, len(items))
	for _, rule := range items {
		ruleMap := map[string]interface{}{}

		if rule.ConfigRuleId != nil {
			ruleMap["config_rule_id"] = rule.ConfigRuleId
		}

		if rule.Identifier != nil {
			ruleMap["identifier"] = rule.Identifier
		}

		if rule.RuleName != nil {
			ruleMap["rule_name"] = rule.RuleName
		}

		if rule.Description != nil {
			ruleMap["description"] = rule.Description
		}

		if rule.RiskLevel != nil {
			ruleMap["risk_level"] = int(*rule.RiskLevel)
		}

		if rule.ServiceFunction != nil {
			ruleMap["service_function"] = rule.ServiceFunction
		}

		if rule.CreateTime != nil {
			ruleMap["create_time"] = rule.CreateTime
		}

		if rule.Status != nil {
			ruleMap["status"] = rule.Status
		}

		if rule.ComplianceResult != nil {
			ruleMap["compliance_result"] = rule.ComplianceResult
		}

		if rule.IdentifierType != nil {
			ruleMap["identifier_type"] = rule.IdentifierType
		}

		if rule.CompliancePackId != nil {
			ruleMap["compliance_pack_id"] = rule.CompliancePackId
		}

		if rule.CompliancePackName != nil {
			ruleMap["compliance_pack_name"] = rule.CompliancePackName
		}

		if rule.ResourceType != nil {
			resourceTypes := make([]string, 0, len(rule.ResourceType))
			for _, rt := range rule.ResourceType {
				if rt != nil {
					resourceTypes = append(resourceTypes, *rt)
				}
			}

			ruleMap["resource_type"] = resourceTypes
		}

		if rule.Labels != nil {
			labels := make([]string, 0, len(rule.Labels))
			for _, l := range rule.Labels {
				if l != nil {
					labels = append(labels, *l)
				}
			}

			ruleMap["labels"] = labels
		}

		if rule.ConfigRuleInvokedTime != nil {
			ruleMap["config_rule_invoked_time"] = rule.ConfigRuleInvokedTime
		}

		ruleList = append(ruleList, ruleMap)
	}

	return ruleList
}
