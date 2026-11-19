package config

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudConfigAggregateRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudConfigAggregateRuleCreate,
		Read:   resourceTencentCloudConfigAggregateRuleRead,
		Update: resourceTencentCloudConfigAggregateRuleUpdate,
		Delete: resourceTencentCloudConfigAggregateRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Rule template identifier. For system preset rules use the identifier name; for custom rules use the cloud function ARN (region:functionName).",
			},

			"identifier_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Rule template type. Valid values: SYSTEM, CUSTOMIZE.",
			},

			"rule_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule name.",
			},

			"resource_type": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Supported resource type list of the rule.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"account_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Account group id of the aggregate rule.",
			},

			"trigger_type": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Trigger type list, up to 2 entries.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"message_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Trigger message type. Valid values: ScheduledNotification, ConfigurationItemChangeNotification.",
						},
						"maximum_execution_frequency": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Maximum execution frequency (only for ScheduledNotification).",
						},
					},
				},
			},

			"risk_level": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).",
			},

			"input_parameter": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Rule input parameter list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter key.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameter type: Require or Optional.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameter value.",
						},
					},
				},
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule description.",
			},

			"region_scope": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Region scope list; rule only applies to resources in the specified regions.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags_scope": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tag scope list; rule only applies to resources with the specified tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"exclude_resource_ids_scope": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Resource ID list excluded from rule evaluation.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// Computed
			"config_rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Config rule ID.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule status. Valid values: ACTIVE (enabled), NO_ACTIVE (stopped).",
			},

			"compliance_result": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Compliance result. Valid values: COMPLIANT, NON_COMPLIANT, NOT_APPLICABLE.",
			},

			"config_rule_invoked_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last rule evaluation time.",
			},
		},
	}
}

func resourceTencentCloudConfigAggregateRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_aggregate_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = configv20220802.NewAddAggregateConfigRuleRequest()
	)

	if v, ok := d.GetOk("identifier"); ok {
		request.Identifier = helper.String(v.(string))
	}

	if v, ok := d.GetOk("identifier_type"); ok {
		request.IdentifierType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_name"); ok {
		request.RuleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_type"); ok {
		rawList := v.([]interface{})
		for _, item := range rawList {
			request.ResourceType = append(request.ResourceType, helper.String(item.(string)))
		}
	}

	if v, ok := d.GetOk("account_group_id"); ok {
		request.AccountGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("trigger_type"); ok {
		request.TriggerType = buildConfigRuleTriggerTypes(v.([]interface{}))
	}

	if v, ok := d.GetOk("risk_level"); ok {
		request.RiskLevel = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("input_parameter"); ok {
		request.InputParameter = buildInputParameters(v.([]interface{}))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("region_scope"); ok {
		rawList := v.([]interface{})
		for _, item := range rawList {
			request.RegionScope = append(request.RegionScope, helper.String(item.(string)))
		}
	}

	if v, ok := d.GetOk("tags_scope"); ok {
		request.TagsScope = buildConfigRuleTagsScope(v.([]interface{}))
	}

	if v, ok := d.GetOk("exclude_resource_ids_scope"); ok {
		rawList := v.([]interface{})
		for _, item := range rawList {
			request.ExcludeResourceIdsScope = append(request.ExcludeResourceIdsScope, helper.String(item.(string)))
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().AddAggregateConfigRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create aggregate config rule failed, Response is nil"))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create aggregate config rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	accountGroupId := d.Get("account_group_id").(string)

	service := ConfigService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	ruleId, err := service.WaitAggregateConfigRuleCreated(ctx, accountGroupId, request)
	if err != nil {
		return err
	}

	d.SetId(helper.IdFormat(accountGroupId, ruleId))

	return resourceTencentCloudConfigAggregateRuleRead(d, meta)
}

func resourceTencentCloudConfigAggregateRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_aggregate_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		idParts = strings.Split(d.Id(), tccommon.FILED_SP)
	)

	if len(idParts) != 2 {
		return fmt.Errorf("id is malformed, id:%s", d.Id())
	}

	accountGroupId := idParts[0]
	ruleId := idParts[1]

	service := ConfigService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	respData, err := service.DescribeAggregateConfigRuleById(ctx, accountGroupId, ruleId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_config_aggregate_rule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("account_group_id", accountGroupId)

	if respData.Identifier != nil {
		_ = d.Set("identifier", respData.Identifier)
	}

	if respData.IdentifierType != nil {
		_ = d.Set("identifier_type", respData.IdentifierType)
	}

	if respData.RuleName != nil {
		_ = d.Set("rule_name", respData.RuleName)
	}

	if respData.ResourceType != nil {
		resourceTypes := make([]string, 0, len(respData.ResourceType))
		for _, v := range respData.ResourceType {
			if v != nil {
				resourceTypes = append(resourceTypes, *v)
			}
		}

		_ = d.Set("resource_type", resourceTypes)
	}

	if respData.TriggerType != nil {
		triggerTypes := make([]map[string]interface{}, 0, len(respData.TriggerType))
		for _, tt := range respData.TriggerType {
			ttMap := map[string]interface{}{}
			if tt.MessageType != nil {
				ttMap["message_type"] = tt.MessageType
			}

			if tt.MaximumExecutionFrequency != nil {
				ttMap["maximum_execution_frequency"] = tt.MaximumExecutionFrequency
			}

			triggerTypes = append(triggerTypes, ttMap)
		}

		_ = d.Set("trigger_type", triggerTypes)
	}

	if respData.RiskLevel != nil {
		_ = d.Set("risk_level", int(*respData.RiskLevel))
	}

	if respData.InputParameter != nil {
		inputParams := make([]map[string]interface{}, 0, len(respData.InputParameter))
		for _, param := range respData.InputParameter {
			paramMap := map[string]interface{}{}
			if param.ParameterKey != nil {
				paramMap["parameter_key"] = param.ParameterKey
			}

			if param.Type != nil {
				paramMap["type"] = param.Type
			}

			if param.Value != nil {
				paramMap["value"] = param.Value
			}

			inputParams = append(inputParams, paramMap)
		}

		_ = d.Set("input_parameter", inputParams)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.RegionsScope != nil {
		regions := make([]string, 0, len(respData.RegionsScope))
		for _, v := range respData.RegionsScope {
			if v != nil {
				regions = append(regions, *v)
			}
		}

		_ = d.Set("region_scope", regions)
	}

	if respData.TagsScope != nil {
		tagsScope := make([]map[string]interface{}, 0, len(respData.TagsScope))
		for _, tag := range respData.TagsScope {
			tagMap := map[string]interface{}{}
			if tag.TagKey != nil {
				tagMap["tag_key"] = tag.TagKey
			}

			if tag.TagValue != nil {
				tagMap["tag_value"] = tag.TagValue
			}

			tagsScope = append(tagsScope, tagMap)
		}

		_ = d.Set("tags_scope", tagsScope)
	}

	if respData.ExcludeResourceIdsScope != nil {
		excludeIds := make([]string, 0, len(respData.ExcludeResourceIdsScope))
		for _, v := range respData.ExcludeResourceIdsScope {
			if v != nil {
				excludeIds = append(excludeIds, *v)
			}
		}

		_ = d.Set("exclude_resource_ids_scope", excludeIds)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.ConfigRuleId != nil {
		_ = d.Set("config_rule_id", respData.ConfigRuleId)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.ComplianceResult != nil {
		_ = d.Set("compliance_result", respData.ComplianceResult)
	}

	if respData.ConfigRuleInvokedTime != nil {
		_ = d.Set("config_rule_invoked_time", respData.ConfigRuleInvokedTime)
	}

	return nil
}

func resourceTencentCloudConfigAggregateRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_aggregate_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		idParts = strings.Split(d.Id(), tccommon.FILED_SP)
	)

	if len(idParts) != 2 {
		return fmt.Errorf("id is malformed, id:%s", d.Id())
	}

	accountGroupId := idParts[0]
	ruleId := idParts[1]

	request := configv20220802.NewUpdateAggregateConfigRuleRequest()
	request.RuleId = &ruleId
	request.AccountGroupId = &accountGroupId

	if v, ok := d.GetOk("rule_name"); ok {
		request.RuleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("trigger_type"); ok {
		request.TriggerType = buildConfigRuleTriggerTypes(v.([]interface{}))
	}

	if v, ok := d.GetOk("risk_level"); ok {
		request.RiskLevel = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("input_parameter"); ok {
		request.InputParameter = buildInputParameters(v.([]interface{}))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("region_scope"); ok {
		rawList := v.([]interface{})
		for _, item := range rawList {
			request.RegionScope = append(request.RegionScope, helper.String(item.(string)))
		}
	}

	if v, ok := d.GetOk("tags_scope"); ok {
		request.TagsScope = buildConfigRuleTagsScope(v.([]interface{}))
	}

	if v, ok := d.GetOk("exclude_resource_ids_scope"); ok {
		rawList := v.([]interface{})
		for _, item := range rawList {
			request.ExcludeResourceIdsScope = append(request.ExcludeResourceIdsScope, helper.String(item.(string)))
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().UpdateAggregateConfigRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update aggregate config rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudConfigAggregateRuleRead(d, meta)
}

func resourceTencentCloudConfigAggregateRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_config_aggregate_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		idParts = strings.Split(d.Id(), tccommon.FILED_SP)
	)

	if len(idParts) != 2 {
		return fmt.Errorf("id is malformed, id:%s", d.Id())
	}

	accountGroupId := idParts[0]
	ruleId := idParts[1]

	request := configv20220802.NewDeleteAggregateConfigRuleRequest()
	request.RuleId = &ruleId
	request.AccountGroupId = &accountGroupId

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseConfigV20220802Client().DeleteAggregateConfigRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete aggregate config rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
