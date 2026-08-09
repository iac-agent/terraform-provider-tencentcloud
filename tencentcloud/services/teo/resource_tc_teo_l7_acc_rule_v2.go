package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoL7AccRuleV2() *schema.Resource {
	return &schema.Resource{
		Create: ResourceTencentCloudTeoL7AccRuleV2Create,
		Read:   ResourceTencentCloudTeoL7AccRuleV2Read,
		Update: ResourceTencentCloudTeoL7AccRuleV2Update,
		Delete: ResourceTencentCloudTeoL7AccRuleV2Delete,
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
			"branches": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Sub-Rule branch. this list currently supports filling in only one rule; multiple entries are invalid.",
				Elem: &schema.Resource{
					Schema: TencentTeoL7RuleBranchBasicInfo(1),
				},
			},
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule ID. Unique identifier of the rule.",
			},
			"rule_priority": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Rule priority. only used as an output parameter.",
			},
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions for querying L7 acceleration rules. Each filter contains a name and a list of values.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter name, e.g. rule-id.",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Filter values.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of L7 acceleration rules matched by the filters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule status.",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule ID.",
						},
						"rule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule name.",
						},
						"description": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Rule description.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"branches": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Sub-Rule branch.",
							Elem: &schema.Resource{
								Schema: TencentTeoL7RuleBranchBasicInfo(1),
							},
						},
						"rule_priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Rule priority.",
						},
					},
				},
			},
		},
	}
}

func ResourceTencentCloudTeoL7AccRuleV2Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()
	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		zoneId string
		ruleId string
	)
	zoneId = d.Get("zone_id").(string)
	request := teov20220901.NewCreateL7AccRulesRequest()
	request.ZoneId = helper.String(zoneId)
	rule := &teov20220901.RuleEngineItem{}

	if v, ok := d.GetOk("status"); ok {
		rule.Status = helper.String(v.(string))
	}
	if v, ok := d.GetOk("rule_name"); ok {
		rule.RuleName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("description"); ok {
		descriptionSet := v.([]interface{})
		for i := range descriptionSet {
			description := descriptionSet[i].(string)
			rule.Description = append(rule.Description, helper.String(description))
		}
	}

	if v, ok := d.GetOk("branches"); ok {
		rule.Branches = resourceTencentCloudTeoL7AccRuleGetBranchs(map[string]interface{}{"branches": v})
	}

	request.Rules = []*teov20220901.RuleEngineItem{rule}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateL7AccRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if result.Response != nil && len(result.Response.RuleIds) > 0 && result.Response.RuleIds[0] != nil {
			ruleId = *result.Response.RuleIds[0]
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(zoneId + tccommon.FILED_SP + ruleId)

	return ResourceTencentCloudTeoL7AccRuleV2Read(d, meta)
}

func ResourceTencentCloudTeoL7AccRuleV2Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	ruleId := idSplit[1]

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("rule_id", ruleId)

	respData, err := service.DescribeTeoL7AccRuleById(ctx, zoneId, ruleId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `teo_l7_acc_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if len(respData.Rules) > 0 {
		rule := respData.Rules[0]
		_ = d.Set("status", rule.Status)
		_ = d.Set("rule_name", rule.RuleName)
		_ = d.Set("description", rule.Description)
		_ = d.Set("rule_priority", rule.RulePriority)
		_ = d.Set("branches", resourceTencentCloudTeoL7AccRuleSetBranchs(rule.Branches))
	}

	// Handle filters and rules computed output
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		filters := make([]*teov20220901.Filter, 0, len(filtersSet))
		for i := range filtersSet {
			filterMap := filtersSet[i].(map[string]interface{})
			filter := &teov20220901.Filter{}
			if name, ok := filterMap["name"]; ok {
				filter.Name = helper.String(name.(string))
			}
			if values, ok := filterMap["values"]; ok {
				valuesList := values.([]interface{})
				filter.Values = make([]*string, 0, len(valuesList))
				for j := range valuesList {
					filter.Values = append(filter.Values, helper.String(valuesList[j].(string)))
				}
			}
			filters = append(filters, filter)
		}

		rulesRespData, err := service.DescribeTeoL7AccRuleByFilters(ctx, zoneId, filters)
		if err != nil {
			return err
		}

		if rulesRespData != nil && len(rulesRespData.Rules) > 0 {
			_ = d.Set("rules", flattenTeoL7AccRules(rulesRespData.Rules))
		} else {
			_ = d.Set("rules", []interface{}{})
		}
	}

	return nil
}

func ResourceTencentCloudTeoL7AccRuleV2Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	ruleId := idSplit[1]

	request := teov20220901.NewModifyL7AccRuleRequest()
	request.ZoneId = helper.String(zoneId)
	rule := &teov20220901.RuleEngineItem{}
	rule.RuleId = &ruleId

	if d.HasChange("status") || d.HasChange("rule_name") || d.HasChange("description") || d.HasChange("branches") {
		if v, ok := d.GetOk("status"); ok {
			rule.Status = helper.String(v.(string))
		}
		if v, ok := d.GetOk("rule_name"); ok {
			rule.RuleName = helper.String(v.(string))
		}
		if v, ok := d.GetOk("description"); ok {
			descriptionSet := v.([]interface{})
			for i := range descriptionSet {
				description := descriptionSet[i].(string)
				rule.Description = append(rule.Description, helper.String(description))
			}
		}
		if v, ok := d.GetOk("branches"); ok {
			rule.Branches = resourceTencentCloudTeoL7AccRuleGetBranchs(map[string]interface{}{"branches": v})
		}

		request.Rule = rule

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyL7AccRule(request)
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

	return ResourceTencentCloudTeoL7AccRuleV2Read(d, meta)
}

func ResourceTencentCloudTeoL7AccRuleV2Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	ruleId := idSplit[1]

	request := teov20220901.NewDeleteL7AccRulesRequest()
	request.ZoneId = helper.String(zoneId)
	request.RuleIds = helper.Strings([]string{ruleId})

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

	return ResourceTencentCloudTeoL7AccRuleV2Read(d, meta)
}

func flattenTeoL7AccRules(rules []*teov20220901.RuleEngineItem) []interface{} {
	result := make([]interface{}, 0, len(rules))
	for _, rule := range rules {
		if rule == nil {
			continue
		}
		ruleMap := map[string]interface{}{}

		if rule.Status != nil {
			ruleMap["status"] = rule.Status
		}

		if rule.RuleId != nil {
			ruleMap["rule_id"] = rule.RuleId
		}

		if rule.RuleName != nil {
			ruleMap["rule_name"] = rule.RuleName
		}

		if rule.Description != nil {
			descList := make([]string, 0, len(rule.Description))
			for _, desc := range rule.Description {
				if desc != nil {
					descList = append(descList, *desc)
				}
			}
			ruleMap["description"] = descList
		}

		if rule.Branches != nil {
			ruleMap["branches"] = resourceTencentCloudTeoL7AccRuleSetBranchs(rule.Branches)
		}

		if rule.RulePriority != nil {
			ruleMap["rule_priority"] = int(*rule.RulePriority)
		}

		result = append(result, ruleMap)
	}
	return result
}
