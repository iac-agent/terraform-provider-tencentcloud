package teo

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoL7AccRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoL7AccRulesCreate,
		Read:   resourceTencentCloudTeoL7AccRulesRead,
		Update: resourceTencentCloudTeoL7AccRulesUpdate,
		Delete: resourceTencentCloudTeoL7AccRulesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone id.",
			},

			"rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Rules content.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule ID. Unique identifier of the rule.",
						},
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule status. The possible values are: `enable`: enabled; `disable`: disabled.",
						},
						"rule_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule name. The name length limit is 255 characters.",
						},
						"description": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Rule annotation. multiple annotations can be added.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"rule_priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Rule priority. only used as an output parameter.",
						},
						"branches": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Sub-Rule branch. this list currently supports filling in only one rule; multiple entries are invalid.",
							Elem: &schema.Resource{
								Schema: TencentTeoL7RuleBranchBasicInfo(1),
							},
						},
					},
				},
			},

			"rule_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of rule IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTencentCloudTeoL7AccRulesCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rules.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		zoneId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("rules"); ok {
		rules := buildRuleEngineItemsFromRules(v.([]interface{}))

		request := teov20220901.NewCreateL7AccRulesRequest()
		request.ZoneId = helper.String(zoneId)
		request.Rules = rules

		var ruleIds []*string
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateL7AccRules(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("create l7_acc_rules returned empty response"))
			}
			ruleIds = result.Response.RuleIds
			return nil
		})
		if err != nil {
			return err
		}

		ruleIdsStr := make([]string, 0, len(ruleIds))
		for _, id := range ruleIds {
			if id != nil {
				ruleIdsStr = append(ruleIdsStr, *id)
			}
		}
		_ = d.Set("rule_ids", ruleIdsStr)
	}

	d.SetId(zoneId)

	return resourceTencentCloudTeoL7AccRulesRead(d, meta)
}

func resourceTencentCloudTeoL7AccRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rules.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	zoneId := d.Id()

	_ = d.Set("zone_id", zoneId)

	respData, err := service.DescribeTeoL7AccRuleById(ctx, zoneId, "")
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `teo_l7_acc_rules` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	rulesList := make([]map[string]interface{}, 0, len(respData.Rules))
	ruleIds := make([]string, 0, len(respData.Rules))
	if respData.Rules != nil {
		for _, rule := range respData.Rules {
			rulesMap := map[string]interface{}{}

			if rule.RuleId != nil {
				rulesMap["rule_id"] = *rule.RuleId
				ruleIds = append(ruleIds, *rule.RuleId)
			}

			if rule.Status != nil {
				rulesMap["status"] = *rule.Status
			}

			if rule.RuleName != nil {
				rulesMap["rule_name"] = *rule.RuleName
			}

			if rule.Description != nil {
				desc := make([]string, 0, len(rule.Description))
				for _, d := range rule.Description {
					if d != nil {
						desc = append(desc, *d)
					}
				}
				rulesMap["description"] = desc
			}

			if rule.RulePriority != nil {
				rulesMap["rule_priority"] = *rule.RulePriority
			}

			if rule.Branches != nil {
				rulesMap["branches"] = resourceTencentCloudTeoL7AccRuleSetBranchs(rule.Branches)
			}
			rulesList = append(rulesList, rulesMap)
		}

		_ = d.Set("rules", rulesList)
		_ = d.Set("rule_ids", ruleIds)
	}

	_ = zoneId
	return nil
}

func resourceTencentCloudTeoL7AccRulesUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rules.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	zoneId := d.Id()

	if d.HasChange("rules") {
		oldRaw, newRaw := d.GetChange("rules")
		oldRules := oldRaw.([]interface{})
		newRules := newRaw.([]interface{})

		// Build maps for comparison
		oldRuleMap := make(map[string]map[string]interface{})
		for _, item := range oldRules {
			ruleMap := item.(map[string]interface{})
			if ruleId, ok := ruleMap["rule_id"].(string); ok && ruleId != "" {
				oldRuleMap[ruleId] = ruleMap
			}
		}

		newRuleMap := make(map[string]map[string]interface{})
		newRulesWithoutId := make([]interface{}, 0)
		for _, item := range newRules {
			ruleMap := item.(map[string]interface{})
			if ruleId, ok := ruleMap["rule_id"].(string); ok && ruleId != "" {
				newRuleMap[ruleId] = ruleMap
			} else {
				newRulesWithoutId = append(newRulesWithoutId, item)
			}
		}

		// Delete removed rules
		deleteRuleIds := make([]string, 0)
		for oldRuleId := range oldRuleMap {
			if _, exists := newRuleMap[oldRuleId]; !exists {
				deleteRuleIds = append(deleteRuleIds, oldRuleId)
			}
		}
		if len(deleteRuleIds) > 0 {
			deleteRequest := teov20220901.NewDeleteL7AccRulesRequest()
			deleteRequest.ZoneId = helper.String(zoneId)
			deleteRequest.RuleIds = helper.Strings(deleteRuleIds)
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteL7AccRules(deleteRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, deleteRequest.GetAction(), deleteRequest.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				return err
			}
		}

		// Modify existing rules
		for ruleId, newRule := range newRuleMap {
			if _, exists := oldRuleMap[ruleId]; exists {
				modifyRequest := teov20220901.NewModifyL7AccRuleRequest()
				modifyRequest.ZoneId = helper.String(zoneId)
				rule := buildRuleEngineItemFromRuleMap(newRule)
				rule.RuleId = helper.String(ruleId)
				modifyRequest.Rule = rule

				err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyL7AccRule(modifyRequest)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyRequest.GetAction(), modifyRequest.ToJsonString(), result.ToJsonString())
					}
					return nil
				})
				if err != nil {
					return err
				}
			}
		}

		// Create new rules
		if len(newRulesWithoutId) > 0 {
			newRules := buildRuleEngineItemsFromRules(newRulesWithoutId)
			createRequest := teov20220901.NewCreateL7AccRulesRequest()
			createRequest.ZoneId = helper.String(zoneId)
			createRequest.Rules = newRules

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateL7AccRules(createRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, createRequest.GetAction(), createRequest.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
	}

	_ = zoneId
	return resourceTencentCloudTeoL7AccRulesRead(d, meta)
}

func resourceTencentCloudTeoL7AccRulesDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rules.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	zoneId := d.Id()

	ruleIds := make([]string, 0)
	if v, ok := d.GetOk("rule_ids"); ok {
		for _, item := range v.([]interface{}) {
			ruleIds = append(ruleIds, item.(string))
		}
	}

	if len(ruleIds) > 0 {
		request := teov20220901.NewDeleteL7AccRulesRequest()
		request.ZoneId = helper.String(zoneId)
		request.RuleIds = helper.Strings(ruleIds)

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteL7AccRules(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudTeoL7AccRulesRead(d, meta)
}

func buildRuleEngineItemsFromRules(rules []interface{}) []*teov20220901.RuleEngineItem {
	result := make([]*teov20220901.RuleEngineItem, 0, len(rules))
	for _, item := range rules {
		ruleMap := item.(map[string]interface{})
		rule := buildRuleEngineItemFromRuleMap(ruleMap)
		result = append(result, rule)
	}
	return result
}

func buildRuleEngineItemFromRuleMap(ruleMap map[string]interface{}) *teov20220901.RuleEngineItem {
	rule := &teov20220901.RuleEngineItem{}
	if v, ok := ruleMap["status"].(string); ok && v != "" {
		rule.Status = helper.String(v)
	}
	if v, ok := ruleMap["rule_name"].(string); ok && v != "" {
		rule.RuleName = helper.String(v)
	}
	if v, ok := ruleMap["description"]; ok {
		descriptionSet := v.([]interface{})
		for i := range descriptionSet {
			description := descriptionSet[i].(string)
			rule.Description = append(rule.Description, helper.String(description))
		}
	}
	if _, ok := ruleMap["branches"]; ok {
		rule.Branches = resourceTencentCloudTeoL7AccRuleGetBranchs(ruleMap)
	}
	return rule
}
