package fwm

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudFwmEdgeAclRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudFwmEdgeAclRulesRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Edge ACL rule group ID.",
			},

			"direction": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Rule direction: 0-egress, 1-ingress.",
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter list. Each filter contains name, values and operator_type (1-equal, 7-in, 9-fuzzy match). Supported filter names: RuleId, Direction, Protocol, RuleAction, etc.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter field name.",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Filter value list.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"operator_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Operator type: 1-equal, 7-in, 9-fuzzy match.",
						},
					},
				},
			},

			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Page size, default 100, max 1000.",
			},

			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Page offset, default 0.",
			},

			"order": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort order: asc-ascending, desc-descending.",
			},

			"by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort field, supports Sequence, RuleId, etc.",
			},

			// computed
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of rules matching the filter.",
			},

			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of edge ACL rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule ID.",
						},
						"sequence": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Rule execution sequence.",
						},
						"direction": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Rule direction: 0-egress, 1-ingress.",
						},
						"source_content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source address content.",
						},
						"source_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source address type.",
						},
						"source_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source address name (returned when the type is template/instance/tag).",
						},
						"target_content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target address content.",
						},
						"target_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target address type.",
						},
						"target_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target address name (returned when the type is template/instance/tag).",
						},
						"port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Port.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol type.",
						},
						"rule_action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule action: accept-pass, drop-block, log-observe.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule description.",
						},
						"scope": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Effective scope of the rule.",
						},
						"param_template_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter template ID.",
						},
						"invalid": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the rule is invalid: 0-valid, 1-invalid.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule creation time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule latest update time.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudFwmEdgeAclRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_fwm_edge_acl_rules.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = NewFwmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("group_id"); ok {
		paramMap["GroupId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("direction"); ok {
		paramMap["Direction"] = helper.IntInt64(v.(int))
	}

	if rawFilters, ok := d.GetOk("filters"); ok {
		filters := rawFilters.([]interface{})
		filterList := make([]*fwmv20250611.CommonFilter, 0, len(filters))
		for _, v := range filters {
			filterMap, filterMapOk := v.(map[string]interface{})
			if !filterMapOk {
				continue
			}

			filter := fwmv20250611.CommonFilter{}
			if name, nameOk := filterMap["name"]; nameOk && name.(string) != "" {
				filter.Name = helper.String(name.(string))
			}

			if rawValues, valuesOk := filterMap["values"]; valuesOk {
				values := rawValues.([]interface{})
				filter.Values = make([]*string, 0, len(values))
				for _, value := range values {
					filter.Values = append(filter.Values, helper.String(value.(string)))
				}
			}

			if operatorType, operatorTypeOk := filterMap["operator_type"]; operatorTypeOk {
				filter.OperatorType = helper.IntInt64(operatorType.(int))
			}

			filterList = append(filterList, &filter)
		}

		if len(filterList) > 0 {
			paramMap["Filters"] = filterList
		}
	}

	if v, ok := d.GetOk("limit"); ok {
		paramMap["Limit"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("offset"); ok {
		paramMap["Offset"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("by"); ok {
		paramMap["By"] = helper.String(v.(string))
	}

	var rules []*fwmv20250611.EdgeAclRuleResp
	var totalCount int64
	rErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		var e error
		rules, totalCount, e = service.DescribeFwmEdgeAclRulesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if rErr != nil {
		log.Printf("[CRITAL]%s read fwm edge acl rules failed, reason:%+v", logId, rErr)
		return rErr
	}

	ids := make([]string, 0, len(rules))
	rulesList := make([]map[string]interface{}, 0, len(rules))
	if rules != nil {
		for _, rule := range rules {
			ruleMap := map[string]interface{}{}

			if rule.RuleId != nil {
				ruleMap["rule_id"] = rule.RuleId
				ids = append(ids, *rule.RuleId)
			}

			if rule.Sequence != nil {
				ruleMap["sequence"] = rule.Sequence
			}

			if rule.Direction != nil {
				ruleMap["direction"] = rule.Direction
			}

			if rule.SourceContent != nil {
				ruleMap["source_content"] = rule.SourceContent
			}

			if rule.SourceType != nil {
				ruleMap["source_type"] = rule.SourceType
			}

			if rule.SourceName != nil {
				ruleMap["source_name"] = rule.SourceName
			}

			if rule.TargetContent != nil {
				ruleMap["target_content"] = rule.TargetContent
			}

			if rule.TargetType != nil {
				ruleMap["target_type"] = rule.TargetType
			}

			if rule.TargetName != nil {
				ruleMap["target_name"] = rule.TargetName
			}

			if rule.Port != nil {
				ruleMap["port"] = rule.Port
			}

			if rule.Protocol != nil {
				ruleMap["protocol"] = rule.Protocol
			}

			if rule.RuleAction != nil {
				ruleMap["rule_action"] = rule.RuleAction
			}

			if rule.Description != nil {
				ruleMap["description"] = rule.Description
			}

			if rule.Scope != nil {
				ruleMap["scope"] = rule.Scope
			}

			if rule.ParamTemplateId != nil {
				ruleMap["param_template_id"] = rule.ParamTemplateId
			}

			if rule.Invalid != nil {
				ruleMap["invalid"] = rule.Invalid
			}

			if rule.CreateTime != nil {
				ruleMap["create_time"] = rule.CreateTime
			}

			if rule.UpdateTime != nil {
				ruleMap["update_time"] = rule.UpdateTime
			}

			rulesList = append(rulesList, ruleMap)
		}
	}

	_ = d.Set("total_count", totalCount)
	_ = d.Set("rules", rulesList)

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), rulesList); e != nil {
			return e
		}
	}

	return nil
}
