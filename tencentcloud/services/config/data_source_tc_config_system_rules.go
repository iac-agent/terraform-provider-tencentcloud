package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudConfigSystemRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudConfigSystemRulesRead,
		Schema: map[string]*schema.Schema{
			"keyword": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search keyword. Supports identifier/名称/标签/描述 search。",
			},

			"risk_level": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "风险等级 对于 filtering. 有效值：1 (high risk)，2 (medium risk)，3 (low risk)。",
			},

			"rule_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "System preset 规则 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule 唯一 identifier。",
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
							Description: "风险等级 有效值：1 (high risk)，2 (medium risk)，3 (low risk)。",
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
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last 更新时间。",
						},
						"trigger_type": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Trigger 类型 列表。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"resource_type": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Supported 资源类型 列表。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"label": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Rule 标签 列表。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"reference_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 times 此 规则 是 referenced。",
						},
						"identifier_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule 类型",
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

func dataSourceTencentCloudConfigSystemRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_config_system_rules.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ConfigService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := buildSystemRulesParamMap(d)

	respData, reqErr := service.DescribeConfigSystemRulesByFilter(ctx, paramMap)
	if reqErr != nil {
		return reqErr
	}

	ruleList := flattenSystemConfigRuleList(respData)
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

func buildSystemRulesParamMap(d *schema.ResourceData) map[string]interface{} {
	paramMap := make(map[string]interface{})

	if v, ok := d.GetOk("keyword"); ok {
		paramMap["Keyword"] = v.(string)
	}

	if v, ok := d.GetOk("risk_level"); ok {
		paramMap["RiskLevel"] = v.(int)
	}

	return paramMap
}

func flattenSystemConfigRuleList(items []*configv20220802.SystemConfigRule) []map[string]interface{} {
	ruleList := make([]map[string]interface{}, 0, len(items))
	for _, rule := range items {
		ruleMap := map[string]interface{}{}

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

		if rule.UpdateTime != nil {
			ruleMap["update_time"] = rule.UpdateTime
		}

		if rule.TriggerType != nil {
			triggerTypes := make([]string, 0, len(rule.TriggerType))
			for _, t := range rule.TriggerType {
				if t != nil {
					triggerTypes = append(triggerTypes, *t)
				}
			}

			ruleMap["trigger_type"] = triggerTypes
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

		if rule.Label != nil {
			labels := make([]string, 0, len(rule.Label))
			for _, l := range rule.Label {
				if l != nil {
					labels = append(labels, *l)
				}
			}

			ruleMap["label"] = labels
		}

		if rule.ReferenceCount != nil {
			ruleMap["reference_count"] = int(*rule.ReferenceCount)
		}

		if rule.IdentifierType != nil {
			ruleMap["identifier_type"] = rule.IdentifierType
		}

		ruleList = append(ruleList, ruleMap)
	}

	return ruleList
}
