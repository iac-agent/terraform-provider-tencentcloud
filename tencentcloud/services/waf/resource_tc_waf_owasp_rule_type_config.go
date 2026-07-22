package waf

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wafv20180125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWafOwaspRuleTypeConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafOwaspRuleTypeConfigCreate,
		Read:   resourceTencentCloudWafOwaspRuleTypeConfigRead,
		Update: resourceTencentCloudWafOwaspRuleTypeConfigUpdate,
		Delete: resourceTencentCloudWafOwaspRuleTypeConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "域名 名称",
			},

			"type_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Rule 类型 ID。",
			},

			"rule_type_status": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "switch 状态 规则 类型 有效值：0 (已禁用)，1 (已启用)。",
			},

			"rule_type_action": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Protection 模式 的 规则 类型 有效值：0 (observation)，1 (intercept)。",
			},

			"rule_type_level": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Protection 级别 的 规则. 有效值：100 (loose)，200 (normal)，300 (strict)，400 (ultra-strict)。",
			},

			// computed
			"rule_type_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule 类型 名称",
			},

			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule 类型 描述",
			},

			"classification": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Data 类型 category。",
			},

			"total_rule": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "指定all 规则 under 规则 类型 always。",
			},

			"active_rule": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "表示total 数量 规则 已启用 under 规则 类型",
			},
		},
	}
}

func resourceTencentCloudWafOwaspRuleTypeConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_owasp_rule_type_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		domain string
		typeId string
	)

	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}

	if v, ok := d.GetOk("type_id"); ok {
		typeId = v.(string)
	}

	d.SetId(strings.Join([]string{domain, typeId}, tccommon.FILED_SP))
	return resourceTencentCloudWafOwaspRuleTypeConfigUpdate(d, meta)
}

func resourceTencentCloudWafOwaspRuleTypeConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_owasp_rule_type_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	domain := idSplit[0]
	typeId := idSplit[1]

	respData, err := service.DescribeWafOwaspRuleTypeConfigById(ctx, domain, typeId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_waf_owasp_rule_type_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("domain", domain)
	_ = d.Set("type_id", typeId)

	if respData.Status != nil {
		_ = d.Set("rule_type_status", respData.Status)
	}

	if respData.Action != nil {
		_ = d.Set("rule_type_action", respData.Action)
	}

	if respData.Level != nil {
		_ = d.Set("rule_type_level", respData.Level)
	}

	if respData.TypeName != nil {
		_ = d.Set("rule_type_name", respData.TypeName)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.Classification != nil {
		_ = d.Set("classification", respData.Classification)
	}

	if respData.TotalRule != nil {
		_ = d.Set("total_rule", respData.TotalRule)
	}

	if respData.ActiveRule != nil {
		_ = d.Set("active_rule", respData.ActiveRule)
	}

	return nil
}

func resourceTencentCloudWafOwaspRuleTypeConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_owasp_rule_type_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	domain := idSplit[0]
	typeId := idSplit[1]

	if v, ok := d.GetOkExists("rule_type_status"); ok {
		request := wafv20180125.NewModifyOwaspRuleTypeStatusRequest()
		request.RuleTypeStatus = helper.IntInt64(v.(int))
		request.Domain = &domain
		request.TypeIDs = helper.Strings([]string{typeId})
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyOwaspRuleTypeStatusWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update waf owasp rule type status failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if v, ok := d.GetOkExists("rule_type_action"); ok {
		request := wafv20180125.NewModifyOwaspRuleTypeActionRequest()
		request.RuleTypeAction = helper.IntInt64(v.(int))
		request.Domain = &domain
		request.TypeIDs = helper.Strings([]string{typeId})
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyOwaspRuleTypeActionWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update waf owasp rule type action failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if v, ok := d.GetOkExists("rule_type_level"); ok {
		request := wafv20180125.NewModifyOwaspRuleTypeLevelRequest()
		request.RuleTypeLevel = helper.IntInt64(v.(int))
		request.Domain = &domain
		request.TypeIDs = helper.Strings([]string{typeId})
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafV20180125Client().ModifyOwaspRuleTypeLevelWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update waf owasp rule type level failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudWafOwaspRuleTypeConfigRead(d, meta)
}

func resourceTencentCloudWafOwaspRuleTypeConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_owasp_rule_type_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
