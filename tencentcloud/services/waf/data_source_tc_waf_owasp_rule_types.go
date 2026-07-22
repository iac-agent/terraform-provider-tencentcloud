package waf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wafv20180125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWafOwaspRuleTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafOwaspRuleTypesRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "域名 names to be queried。",
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions. supports RuleId，CveID，and Desc。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field 名称，用于filtering\nFilter the sub-顺序 number (值) by DealName。",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Values after filtering。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"exact_match": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Exact search or not。",
						},
					},
				},
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Rule 类型 list and information。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "类型 ID。",
						},
						"type_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 名称",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 描述",
						},
						"classification": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data 类型 category。",
						},
						"action": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Protection 模式 of the rule 类型 有效值：0 (observation)，1 (intercept)。",
						},
						"level": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Protection 级别 of the rule 类型 有效值：100 (loose)，200 (normal)，300 (strict)，400 (ultra-strict)。",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The switch 状态 rule 类型 有效值：0 (已禁用)，1 (已启用)。",
						},
						"total_rule": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "指定all rules under the rule 类型 always。",
						},
						"active_rule": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "表示total 数量 rules 已启用 under the rule 类型",
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

func dataSourceTencentCloudWafOwaspRuleTypesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_waf_owasp_rule_types.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		domain  string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("domain"); ok {
		paramMap["Domain"] = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*wafv20180125.FiltersItemNew, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filtersItemNew := wafv20180125.FiltersItemNew{}
			if v, ok := filtersMap["name"].(string); ok && v != "" {
				filtersItemNew.Name = helper.String(v)
			}

			if v, ok := filtersMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				for i := range valuesSet {
					values := valuesSet[i].(string)
					filtersItemNew.Values = append(filtersItemNew.Values, helper.String(values))
				}
			}

			if v, ok := filtersMap["exact_match"].(bool); ok {
				filtersItemNew.ExactMatch = helper.Bool(v)
			}

			tmpSet = append(tmpSet, &filtersItemNew)
		}

		paramMap["Filters"] = tmpSet
	}

	var respData []*wafv20180125.OwaspRuleType
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafOwaspRuleTypesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	listList := make([]map[string]interface{}, 0, len(respData))
	for _, list := range respData {
		listMap := map[string]interface{}{}
		if list.TypeId != nil {
			listMap["type_id"] = list.TypeId
		}

		if list.TypeName != nil {
			listMap["type_name"] = list.TypeName
		}

		if list.Description != nil {
			listMap["description"] = list.Description
		}

		if list.Classification != nil {
			listMap["classification"] = list.Classification
		}

		if list.Action != nil {
			listMap["action"] = list.Action
		}

		if list.Level != nil {
			listMap["level"] = list.Level
		}

		if list.Status != nil {
			listMap["status"] = list.Status
		}

		if list.TotalRule != nil {
			listMap["total_rule"] = list.TotalRule
		}

		if list.ActiveRule != nil {
			listMap["active_rule"] = list.ActiveRule
		}

		listList = append(listList, listMap)
	}

	_ = d.Set("list", listList)
	d.SetId(domain)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
