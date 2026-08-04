package teo

import (
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
				Required:    true,
				Description: "The list of L7 acceleration rules. Each rule contains rule_name, status, description, and branches.",
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
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions for the Describe API. The upper limit of Filters.Values is 20.",
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
							Description: "Filter field values.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"rule_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The rule IDs returned by the Create API.",
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
		zoneId  string
		ruleIds []string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("rules"); ok {
		rules := v.([]interface{})
		if len(rules) == 0 {
			return fmt.Errorf("rules list cannot be empty")
		}

		request := teov20220901.NewCreateL7AccRulesRequest()
		request.ZoneId = helper.String(zoneId)

		ruleList := make([]*teov20220901.RuleEngineItem, 0, len(rules))
		for _, item := range rules {
			rulesMap := item.(map[string]interface{})
			ruleEngineItem := teov20220901.RuleEngineItem{}
			if v, ok := rulesMap["status"].(string); ok && v != "" {
				ruleEngineItem.Status = helper.String(v)
			}
			if v, ok := rulesMap["rule_name"].(string); ok && v != "" {
				ruleEngineItem.RuleName = helper.String(v)
			}
			if v, ok := rulesMap["description"]; ok {
				descriptionSet := v.([]interface{})
				for i := range descriptionSet {
					description := descriptionSet[i].(string)
					ruleEngineItem.Description = append(ruleEngineItem.Description, helper.String(description))
				}
			}
			if v, ok := rulesMap["branches"]; ok {
				ruleEngineItem.Branches = resourceTencentCloudTeoL7AccRuleGetBranchs(map[string]interface{}{"branches": v})
			}
			ruleList = append(ruleList, &ruleEngineItem)
		}
		request.Rules = ruleList

		var response *teov20220901.CreateL7AccRulesResponse
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateL7AccRules(request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			return err
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil {
			log.Printf("[CRITAL]%s create teo l7_acc_rules failed, zone_id=%s, response is nil", logId, zoneId)
			return fmt.Errorf("create teo l7_acc_rules failed: response is nil")
		}

		if response.Response.RuleIds != nil {
			for _, ruleId := range response.Response.RuleIds {
				if ruleId != nil {
					ruleIds = append(ruleIds, *ruleId)
				}
			}
		}

		if len(ruleIds) == 0 {
			log.Printf("[CRITAL]%s create teo l7_acc_rules failed, zone_id=%s, rule_ids is empty", logId, zoneId)
			return fmt.Errorf("create teo l7_acc_rules failed: rule_ids is empty")
		}
	}

	log.Printf("[DEBUG]%s setting id to zone_id=%s", logId, zoneId)
	d.SetId(zoneId)

	return resourceTencentCloudTeoL7AccRulesRead(d, meta)
}

func resourceTencentCloudTeoL7AccRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rules.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	zoneId := d.Id()

	_ = d.Set("zone_id", zoneId)

	request := teov20220901.NewDescribeL7AccRulesRequest()
	request.ZoneId = helper.String(zoneId)
	request.Limit = helper.IntInt64(1000)
	request.Offset = helper.IntInt64(0)

	if v, ok := d.GetOk("filters"); ok {
		filters := v.([]interface{})
		filterList := make([]*teov20220901.Filter, 0, len(filters))
		for _, item := range filters {
			filterMap := item.(map[string]interface{})
			filter := teov20220901.Filter{}
			if v, ok := filterMap["name"].(string); ok && v != "" {
				filter.Name = helper.String(v)
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.([]interface{})
				for i := range valuesSet {
					values := valuesSet[i].(string)
					filter.Values = append(filter.Values, helper.String(values))
				}
			}
			filterList = append(filterList, &filter)
		}
		request.Filters = filterList
	}

	var response *teov20220901.DescribeL7AccRulesResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeL7AccRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		log.Printf("[CRITAL]%s read teo l7_acc_rules failed, zone_id=%s, response is nil", logId, zoneId)
		return fmt.Errorf("read teo l7_acc_rules failed: response is nil")
	}

	if response.Response.Rules == nil || len(response.Response.Rules) == 0 {
		log.Printf("[CRUD] teo l7_acc_rules id=%s", d.Id())
		log.Printf("[WARN]%s read teo l7_acc_rules empty, zone_id=%s, skip SetId", logId, zoneId)
		return fmt.Errorf("read teo l7_acc_rules failed: rules is empty")
	}

	rulesList := make([]map[string]interface{}, 0, len(response.Response.Rules))
	ruleIds := make([]string, 0, len(response.Response.Rules))
	for _, rule := range response.Response.Rules {
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
			descs := make([]string, 0, len(rule.Description))
			for _, desc := range rule.Description {
				if desc != nil {
					descs = append(descs, *desc)
				}
			}
			rulesMap["description"] = descs
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

	return nil
}

func resourceTencentCloudTeoL7AccRulesUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rules.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	zoneId := d.Id()

	if !d.HasChange("rules") {
		return resourceTencentCloudTeoL7AccRulesRead(d, meta)
	}

	oldRaw, newRaw := d.GetChange("rules")
	oldRules := oldRaw.([]interface{})
	newRules := newRaw.([]interface{})

	// Build maps keyed by rule_name for comparison
	oldRuleMap := make(map[string]map[string]interface{})
	for _, item := range oldRules {
		rulesMap := item.(map[string]interface{})
		if ruleName, ok := rulesMap["rule_name"].(string); ok && ruleName != "" {
			oldRuleMap[ruleName] = rulesMap
		}
	}

	newRuleMap := make(map[string]map[string]interface{})
	for _, item := range newRules {
		rulesMap := item.(map[string]interface{})
		if ruleName, ok := rulesMap["rule_name"].(string); ok && ruleName != "" {
			newRuleMap[ruleName] = rulesMap
		}
	}

	// Find rules to delete (in old but not in new)
	for ruleName, oldRule := range oldRuleMap {
		if _, exists := newRuleMap[ruleName]; !exists {
			// Delete this rule
			if ruleIdVal, ok := oldRule["rule_id"].(string); ok && ruleIdVal != "" {
				deleteRequest := teov20220901.NewDeleteL7AccRulesRequest()
				deleteRequest.ZoneId = helper.String(zoneId)
				deleteRequest.RuleIds = helper.Strings([]string{ruleIdVal})

				err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteL7AccRules(deleteRequest)
					if e != nil {
						return tccommon.RetryError(e)
					}
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, deleteRequest.GetAction(), deleteRequest.ToJsonString(), result.ToJsonString())
					return nil
				})
				if err != nil {
					return err
				}
			}
		}
	}

	// Find rules to create (in new but not in old)
	newRulesToCreate := make([]*teov20220901.RuleEngineItem, 0)
	for ruleName, newRule := range newRuleMap {
		if _, exists := oldRuleMap[ruleName]; !exists {
			ruleEngineItem := buildRuleEngineItem(newRule)
			newRulesToCreate = append(newRulesToCreate, &ruleEngineItem)
		}
	}
	if len(newRulesToCreate) > 0 {
		createRequest := teov20220901.NewCreateL7AccRulesRequest()
		createRequest.ZoneId = helper.String(zoneId)
		createRequest.Rules = newRulesToCreate

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateL7AccRules(createRequest)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, createRequest.GetAction(), createRequest.ToJsonString(), result.ToJsonString())
			return nil
		})
		if err != nil {
			return err
		}
	}

	// Find rules to modify (in both old and new, but differ)
	for ruleName, newRule := range newRuleMap {
		if oldRule, exists := oldRuleMap[ruleName]; exists {
			// Check if rule has changed
			if rulesChanged(oldRule, newRule) {
				ruleIdVal := ""
				if v, ok := oldRule["rule_id"].(string); ok {
					ruleIdVal = v
				}
				if ruleIdVal == "" {
					continue
				}

				ruleEngineItem := buildRuleEngineItem(newRule)
				ruleEngineItem.RuleId = helper.String(ruleIdVal)

				modifyRequest := teov20220901.NewModifyL7AccRuleRequest()
				modifyRequest.ZoneId = helper.String(zoneId)
				modifyRequest.Rule = &ruleEngineItem

				err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyL7AccRule(modifyRequest)
					if e != nil {
						return tccommon.RetryError(e)
					}
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyRequest.GetAction(), modifyRequest.ToJsonString(), result.ToJsonString())
					return nil
				})
				if err != nil {
					return err
				}
			}
		}
	}

	return resourceTencentCloudTeoL7AccRulesRead(d, meta)
}

func resourceTencentCloudTeoL7AccRulesDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rules.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	zoneId := d.Id()

	// Get all rule IDs from state
	if v, ok := d.GetOk("rule_ids"); ok {
		ruleIds := v.([]interface{})
		if len(ruleIds) > 0 {
			ruleIdList := make([]string, 0, len(ruleIds))
			for _, id := range ruleIds {
				ruleIdList = append(ruleIdList, id.(string))
			}

			request := teov20220901.NewDeleteL7AccRulesRequest()
			request.ZoneId = helper.String(zoneId)
			request.RuleIds = helper.Strings(ruleIdList)

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteL7AccRules(request)
				if e != nil {
					return tccommon.RetryError(e)
				}
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				return nil
			})
			if err != nil {
				return err
			}
		}
	}

	d.SetId("")

	return nil
}

func buildRuleEngineItem(rulesMap map[string]interface{}) teov20220901.RuleEngineItem {
	ruleEngineItem := teov20220901.RuleEngineItem{}
	if v, ok := rulesMap["status"].(string); ok && v != "" {
		ruleEngineItem.Status = helper.String(v)
	}
	if v, ok := rulesMap["rule_name"].(string); ok && v != "" {
		ruleEngineItem.RuleName = helper.String(v)
	}
	if v, ok := rulesMap["description"]; ok {
		descriptionSet := v.([]interface{})
		for i := range descriptionSet {
			description := descriptionSet[i].(string)
			ruleEngineItem.Description = append(ruleEngineItem.Description, helper.String(description))
		}
	}
	if v, ok := rulesMap["branches"]; ok {
		ruleEngineItem.Branches = resourceTencentCloudTeoL7AccRuleGetBranchs(map[string]interface{}{"branches": v})
	}
	return ruleEngineItem
}

func rulesChanged(oldRule, newRule map[string]interface{}) bool {
	if oldRule["status"] != newRule["status"] {
		return true
	}
	if oldRule["rule_name"] != newRule["rule_name"] {
		return true
	}
	if fmt.Sprintf("%v", oldRule["description"]) != fmt.Sprintf("%v", newRule["description"]) {
		return true
	}
	if fmt.Sprintf("%v", oldRule["branches"]) != fmt.Sprintf("%v", newRule["branches"]) {
		return true
	}
	return false
}
