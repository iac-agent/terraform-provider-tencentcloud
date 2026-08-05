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

func ResourceTencentCloudTeoL7AccRules() *schema.Resource {
	return &schema.Resource{
		Create: ResourceTencentCloudTeoL7AccRulesCreate,
		Read:   ResourceTencentCloudTeoL7AccRulesRead,
		Update: ResourceTencentCloudTeoL7AccRulesUpdate,
		Delete: ResourceTencentCloudTeoL7AccRulesDelete,
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
				Description: "List of rule configurations. Supports a single rule per resource instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule name. The name length limit is 255 characters.",
						},
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule status. The possible values are: `enable`: enabled; `disable`: disabled.",
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
					},
				},
			},
		},
	}
}

func ResourceTencentCloudTeoL7AccRulesCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rules.create")()
	defer tccommon.InconsistentCheck(d, meta)()
	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		zoneId string
		ruleId string
	)
	zoneId = d.Get("zone_id").(string)
	request := teov20220901.NewCreateL7AccRulesRequest()
	request.ZoneId = helper.String(zoneId)

	rulesList := d.Get("rules").([]interface{})
	rule := &teov20220901.RuleEngineItem{}

	if len(rulesList) > 0 {
		ruleMap := rulesList[0].(map[string]interface{})

		if v, ok := ruleMap["rule_name"]; ok && v.(string) != "" {
			rule.RuleName = helper.String(v.(string))
		}
		if v, ok := ruleMap["status"]; ok && v.(string) != "" {
			rule.Status = helper.String(v.(string))
		}
		if v, ok := ruleMap["description"]; ok {
			descriptionSet := v.([]interface{})
			for i := range descriptionSet {
				description := descriptionSet[i].(string)
				rule.Description = append(rule.Description, helper.String(description))
			}
		}
		if v, ok := ruleMap["branches"]; ok {
			rule.Branches = resourceTencentCloudTeoL7AccRuleGetBranchs(map[string]interface{}{"branches": v})
		}
	}

	request.Rules = []*teov20220901.RuleEngineItem{rule}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateL7AccRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if result == nil || result.Response == nil || len(result.Response.RuleIds) == 0 || result.Response.RuleIds[0] == nil {
			return resource.NonRetryableError(fmt.Errorf("l7_acc_rules %s create response is empty", logId))
		}
		ruleId = *result.Response.RuleIds[0]
		return nil
	})
	if err != nil {
		return err
	}

	if ruleId == "" {
		log.Printf("[CRITAL]%s l7_acc_rules create returned empty rule id, zoneId: %s\n", logId, zoneId)
		return fmt.Errorf("l7_acc_rules create returned empty rule id")
	}

	d.SetId(zoneId + tccommon.FILED_SP + ruleId)

	return ResourceTencentCloudTeoL7AccRulesRead(d, meta)
}

func ResourceTencentCloudTeoL7AccRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rules.read")()
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

	respData, err := service.DescribeTeoL7AccRuleById(ctx, zoneId, ruleId)
	if err != nil {
		return err
	}

	if respData == nil || len(respData.Rules) == 0 {
		log.Printf("[WARN]%s resource `teo_l7_acc_rules` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	rule := respData.Rules[0]
	rulesList := []interface{}{}
	ruleMap := make(map[string]interface{})

	if rule.RuleName != nil {
		ruleMap["rule_name"] = *rule.RuleName
	}
	if rule.Status != nil {
		ruleMap["status"] = *rule.Status
	}
	ruleMap["description"] = rule.Description
	if rule.RuleId != nil {
		ruleMap["rule_id"] = *rule.RuleId
	}
	if rule.RulePriority != nil {
		ruleMap["rule_priority"] = *rule.RulePriority
	}
	ruleMap["branches"] = resourceTencentCloudTeoL7AccRuleSetBranchs(rule.Branches)

	rulesList = append(rulesList, ruleMap)
	_ = d.Set("rules", rulesList)

	return nil
}

func ResourceTencentCloudTeoL7AccRulesUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rules.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	ruleId := idSplit[1]

	if d.HasChange("rules") {
		request := teov20220901.NewModifyL7AccRuleRequest()
		request.ZoneId = helper.String(zoneId)
		rule := &teov20220901.RuleEngineItem{}
		rule.RuleId = &ruleId

		rulesList := d.Get("rules").([]interface{})
		if len(rulesList) > 0 {
			ruleMap := rulesList[0].(map[string]interface{})

			if v, ok := ruleMap["rule_name"]; ok && v.(string) != "" {
				rule.RuleName = helper.String(v.(string))
			}
			if v, ok := ruleMap["status"]; ok && v.(string) != "" {
				rule.Status = helper.String(v.(string))
			}
			if v, ok := ruleMap["description"]; ok {
				descriptionSet := v.([]interface{})
				for i := range descriptionSet {
					description := descriptionSet[i].(string)
					rule.Description = append(rule.Description, helper.String(description))
				}
			}
			if v, ok := ruleMap["branches"]; ok {
				rule.Branches = resourceTencentCloudTeoL7AccRuleGetBranchs(map[string]interface{}{"branches": v})
			}
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

	return ResourceTencentCloudTeoL7AccRulesRead(d, meta)
}

func ResourceTencentCloudTeoL7AccRulesDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_l7_acc_rules.delete")()
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

	return nil
}
