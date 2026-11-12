package fwm

import (
	"context"

	"strconv"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudFwmEdgeAclRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudFwmEdgeAclRulesRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Edge ACL rule group ID, the rule group to query rules from.",
			},
			"direction": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Rule direction: 0-egress, 1-ingress.",
			},
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    20,
				Description: "Filter conditions, support filtering by RuleId, Protocol, RuleAction etc.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter field name, e.g. RuleId, Protocol, RuleAction.",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Filter values.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"operator_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Operator type, 1 means equal.",
						},
					},
				},
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "Page size, default 100, max 1000.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Pagination offset, default 0.",
			},
			"order": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort order: asc or desc.",
			},
			"by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort field, support Sequence, RuleId etc.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of rules matched.",
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Edge ACL rule list.",
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
							Description: "Rule execution order.",
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
							Description: "Source address name.",
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
							Description: "Target address name.",
						},
						"port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target port.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol type.",
						},
						"rule_action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule action: accept, drop, log.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule description.",
						},
						"scope": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule effective scope.",
						},
						"hit_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Rule hit count.",
						},
						"last_hit_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule last hit time.",
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
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = NewFwmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	groupId := d.Get("group_id").(string)
	direction := d.Get("direction").(int)
	limit := d.Get("limit").(int)
	offset := d.Get("offset").(int)
	order := d.Get("order").(string)
	by := d.Get("by").(string)

	filters := make([]*fwmv20250611.CommonFilter, 0)
	if rawFilters, ok := d.GetOk("filters"); ok {
		for _, item := range rawFilters.([]interface{}) {
			filterMap := item.(map[string]interface{})
			name := filterMap["name"].(string)
			filter := &fwmv20250611.CommonFilter{
				Name:   &name,
				Values: nil,
			}
			for _, value := range filterMap["values"].([]interface{}) {
				value := value.(string)
				filter.Values = append(filter.Values, &value)
			}
			if v, ok := filterMap["operator_type"]; ok {
				filter.OperatorType = helper.Int64(int64(v.(int)))
			}
			filters = append(filters, filter)
		}
	}

	var (
		rules      []*fwmv20250611.EdgeAclRuleResp
		totalCount int64
		err        error
	)

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		rules, totalCount, err = service.DescribeFwmEdgeAclRulesByFilter(ctx, groupId, int64(direction), filters, int64(limit), int64(offset), order, by)
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	_ = d.Set("total_count", totalCount)

	tmpList := make([]interface{}, 0, len(rules))
	for _, rule := range rules {
		ruleMap := map[string]interface{}{}
		if rule.RuleId != nil {
			ruleMap["rule_id"] = rule.RuleId
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
		if rule.HitCount != nil {
			ruleMap["hit_count"] = rule.HitCount
		}
		if rule.LastHitTime != nil {
			ruleMap["last_hit_time"] = rule.LastHitTime
		}
		if rule.CreateTime != nil {
			ruleMap["create_time"] = rule.CreateTime
		}
		if rule.UpdateTime != nil {
			ruleMap["update_time"] = rule.UpdateTime
		}
		tmpList = append(tmpList, ruleMap)
	}

	_ = d.Set("data", tmpList)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
