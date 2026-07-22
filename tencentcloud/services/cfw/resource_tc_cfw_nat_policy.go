package cfw

import (
	"context"
	"fmt"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCfwNatPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwNatPolicyCreate,
		Read:   resourceTencentCloudCfwNatPolicyRead,
		Update: resourceTencentCloudCfwNatPolicyUpdate,
		Delete: resourceTencentCloudCfwNatPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"source_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access 来源 示例: net:IP/CIDR(192.168.0.2)。",
			},
			"source_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access 来源 类型: 对于 inbound 规则， 类型 可以 是 net，location，vendor，template; 对于 outbound 规则，它 可以 是 net，实例，标签，template，组。",
			},
			"target_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Example 的 访问 purpose: net: IP/CIDR(192.168.0.2) 域名: 域名 名称 规则，such 作为 *.qq.com。",
			},
			"target_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access purpose 类型: For inbound 规则， 类型 可以 是 net，实例，标签，template，组; 对于 outbound 规则，它 可以 是 net，location，vendor，template。",
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "协议 如果 Direction=1，可选 值: TCP，UDP，ANY; 如果 Direction=0，可选 值: TCP，UDP，ICMP，ANY，HTTP，HTTPS，HTTP/HTTPS，SMTP，SMTPS，SMTP/SMTPS，FTP，和 DNS。",
			},
			"rule_action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(POLICY_RULE_ACTION),
				Description:  "How 流量 集合 在 访问 control 策略 passes through 云 firewall. Values: accept: allow; drop: reject; 日志: observe。",
			},
			"port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "端口 对于 访问 control 策略. 值: -1/-1: All ports 80: 端口 80。",
			},
			"direction": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Rule direction: 1，inbound; 0，outbound。",
			},
			"enable": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      POLICY_ENABLE_TRUE,
				ValidateFunc: tccommon.ValidateAllowedStringValue(POLICY_ENABLE),
				Description:  "Rule 状态，true 表示 已启用，false 表示 已禁用 默认为 true。",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "描述",
			},
			"param_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Parameter template ID. 注意：此字段可能返回 null，表示无法获取有效值。",
			},
			"internal_uuid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Internal ID。",
			},
			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "范围 的 effective 规则. ALL: Global effectiveness; ap-guangzhou: Effective territory; cfwnat-xxx: Effectiveness based 在 实例 dimension。",
			},
			// computed
			"uuid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "唯一 ID corresponding 到 规则，无 need 到 fill 在 当 creating 规则。",
			},
			"order_index": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Execution 顺序",
			},
		},
	}
}

func resourceTencentCloudCfwNatPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_nat_policy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		request           = cfw.NewAddNatAcRuleRequest()
		response          = cfw.NewAddNatAcRuleResponse()
		createNatRuleItem = cfw.CreateNatRuleItem{}
		uuid              string
	)

	if v, ok := d.GetOk("source_content"); ok {
		createNatRuleItem.SourceContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_type"); ok {
		createNatRuleItem.SourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_content"); ok {
		createNatRuleItem.TargetContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_type"); ok {
		createNatRuleItem.TargetType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol"); ok {
		createNatRuleItem.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_action"); ok {
		createNatRuleItem.RuleAction = helper.String(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		createNatRuleItem.Port = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("direction"); ok {
		createNatRuleItem.Direction = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("enable"); ok {
		createNatRuleItem.Enable = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		createNatRuleItem.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("param_template_id"); ok {
		createNatRuleItem.ParamTemplateId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("scope"); ok {
		createNatRuleItem.Scope = helper.String(v.(string))
	}

	request.Rules = append(request.Rules, &createNatRuleItem)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwClient().AddNatAcRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.RuleUuid == nil || len(result.Response.RuleUuid) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Create cfw natPolicy failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cfw natPolicy failed, reason:%+v", logId, err)
		return err
	}

	ruleUuid := *response.Response.RuleUuid[0]
	uuid = strconv.FormatInt(ruleUuid, 10)
	d.SetId(uuid)

	return resourceTencentCloudCfwNatPolicyRead(d, meta)
}

func resourceTencentCloudCfwNatPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_nat_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service  = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		ruleUuid = d.Id()
	)

	natPolicy, err := service.DescribeCfwNatPolicyById(ctx, ruleUuid)
	if err != nil {
		return err
	}

	if natPolicy == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cfw_nat_policy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if natPolicy.SourceContent != nil {
		_ = d.Set("source_content", natPolicy.SourceContent)
	}

	if natPolicy.SourceType != nil {
		_ = d.Set("source_type", natPolicy.SourceType)
	}

	if natPolicy.TargetContent != nil {
		_ = d.Set("target_content", natPolicy.TargetContent)
	}

	if natPolicy.TargetType != nil {
		_ = d.Set("target_type", natPolicy.TargetType)
	}

	if natPolicy.Protocol != nil {
		_ = d.Set("protocol", natPolicy.Protocol)
	}

	if natPolicy.RuleAction != nil {
		_ = d.Set("rule_action", natPolicy.RuleAction)
	}

	if natPolicy.Port != nil {
		_ = d.Set("port", natPolicy.Port)
	}

	if natPolicy.Direction != nil {
		_ = d.Set("direction", natPolicy.Direction)
	}

	if natPolicy.Enable != nil {
		_ = d.Set("enable", natPolicy.Enable)
	}

	if natPolicy.Description != nil {
		_ = d.Set("description", natPolicy.Description)
	}

	if natPolicy.Scope != nil {
		_ = d.Set("scope", natPolicy.Scope)
	}

	if natPolicy.ParamTemplateId != nil {
		_ = d.Set("param_template_id", natPolicy.ParamTemplateId)
	}

	if natPolicy.InternalUuid != nil {
		_ = d.Set("internal_uuid", natPolicy.InternalUuid)
	}

	if natPolicy.Uuid != nil {
		_ = d.Set("uuid", natPolicy.Uuid)
	}

	if natPolicy.OrderIndex != nil {
		_ = d.Set("order_index", natPolicy.OrderIndex)
	}

	return nil
}

func resourceTencentCloudCfwNatPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_nat_policy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		request        = cfw.NewModifyNatAcRuleRequest()
		modifyRuleItem = cfw.CreateNatRuleItem{}
		uuid           = d.Id()
	)

	immutableArgs := []string{"direction"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	uuidInt, _ := strconv.ParseInt(uuid, 10, 64)
	modifyRuleItem.Uuid = &uuidInt

	if v, ok := d.GetOk("source_content"); ok {
		modifyRuleItem.SourceContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_type"); ok {
		modifyRuleItem.SourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_content"); ok {
		modifyRuleItem.TargetContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_type"); ok {
		modifyRuleItem.TargetType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol"); ok {
		modifyRuleItem.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_action"); ok {
		modifyRuleItem.RuleAction = helper.String(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		modifyRuleItem.Port = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("direction"); ok {
		modifyRuleItem.Direction = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("enable"); ok {
		modifyRuleItem.Enable = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		modifyRuleItem.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("param_template_id"); ok {
		modifyRuleItem.ParamTemplateId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("scope"); ok {
		modifyRuleItem.Scope = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("order_index"); ok {
		modifyRuleItem.OrderIndex = helper.IntInt64(v.(int))
	}

	request.Rules = append(request.Rules, &modifyRuleItem)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwClient().ModifyNatAcRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cfw natPolicy failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCfwNatPolicyRead(d, meta)
}

func resourceTencentCloudCfwNatPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_nat_policy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		uuid    = d.Id()
	)

	if err := service.DeleteCfwNatPolicyById(ctx, uuid); err != nil {
		return err
	}

	return nil
}
