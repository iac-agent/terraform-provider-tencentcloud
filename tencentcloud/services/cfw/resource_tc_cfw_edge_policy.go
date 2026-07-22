package cfw

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCfwEdgePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwEdgePolicyCreate,
		Read:   resourceTencentCloudCfwEdgePolicyRead,
		Update: resourceTencentCloudCfwEdgePolicyUpdate,
		Delete: resourceTencentCloudCfwEdgePolicyDelete,
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
				Description: "协议 如果 Direction=1 && 范围=serial，可选 值: TCP UDP ICMP ANY HTTP HTTPS HTTP/HTTPS SMTP SMTPS SMTP/SMTPS FTP DNS; 如果 Direction=1 && 范围!=serial，可选 值: TCP; 如果 Direction=0 && 范围=serial，可选 值: TCP UDP ICMP ANY HTTP HTTPS HTTP/HTTPS SMTP SMTPS SMTP/SMTPS FTP DNS; 如果 Direction=0 && 范围!=serial，可选 值: TCP HTTP/HTTPS TLS/SSL。",
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
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(POLICY_SCOPE),
				Description:  "Effective 范围. serial: serial; side: bypass; all: 全局，默认为 all。",
			},
			"param_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Parameter template ID。",
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

func resourceTencentCloudCfwEdgePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_edge_policy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		request        = cfw.NewAddAclRuleRequest()
		response       = cfw.NewAddAclRuleResponse()
		createRuleItem = cfw.CreateRuleItem{}
		uuid           string
	)

	if v, ok := d.GetOk("source_content"); ok {
		createRuleItem.SourceContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_type"); ok {
		createRuleItem.SourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_content"); ok {
		createRuleItem.TargetContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_type"); ok {
		createRuleItem.TargetType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol"); ok {
		createRuleItem.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_action"); ok {
		createRuleItem.RuleAction = helper.String(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		createRuleItem.Port = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("direction"); ok {
		createRuleItem.Direction = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("enable"); ok {
		createRuleItem.Enable = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		createRuleItem.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("scope"); ok {
		createRuleItem.Scope = helper.String(v.(string))
	}

	if v, ok := d.GetOk("param_template_id"); ok {
		createRuleItem.ParamTemplateId = helper.String(v.(string))
	}

	request.Rules = append(request.Rules, &createRuleItem)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwClient().AddAclRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return tccommon.RetryError(fmt.Errorf("Create cfw edgePolicy failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cfw edgePolicy failed, reason:%+v", logId, err)
		return err
	}

	if len(response.Response.RuleUuid) == 0 {
		return fmt.Errorf("RuleUuid is nil.")
	}

	ruleUuid := *response.Response.RuleUuid[0]
	uuid = strconv.FormatInt(ruleUuid, 10)
	d.SetId(uuid)

	return resourceTencentCloudCfwEdgePolicyRead(d, meta)
}

func resourceTencentCloudCfwEdgePolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_edge_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		ruleUuid   = d.Id()
		sourceType string
		targetType string
	)

	edgePolicy, err := service.DescribeCfwEdgePolicyById(ctx, ruleUuid)
	if err != nil {
		return err
	}

	if edgePolicy == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cfw_edge_policy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if edgePolicy.SourceType != nil {
		_ = d.Set("source_type", edgePolicy.SourceType)
		sourceType = *edgePolicy.SourceType
	}

	if edgePolicy.SourceContent != nil {
		if sourceType == "tag" {
			params := strings.Split(*edgePolicy.SourceContent, "|")
			key := params[0]
			value := params[1]
			var obj SourceContentJson
			obj.Key = key
			obj.Value = value
			tmpStr, _ := json.Marshal(obj)
			_ = d.Set("source_content", string(tmpStr))
		} else {
			_ = d.Set("source_content", edgePolicy.SourceContent)
		}
	}

	if edgePolicy.TargetType != nil {
		_ = d.Set("target_type", edgePolicy.TargetType)
		targetType = *edgePolicy.TargetType
	}

	if edgePolicy.TargetContent != nil {
		if targetType == "tag" {
			params := strings.Split(*edgePolicy.TargetContent, "|")
			key := params[0]
			value := params[1]
			var obj TargetContentJson
			obj.Key = key
			obj.Value = value
			tmpStr, _ := json.Marshal(obj)
			_ = d.Set("target_content", string(tmpStr))
		} else {
			_ = d.Set("target_content", edgePolicy.TargetContent)
		}
	}

	if edgePolicy.Protocol != nil {
		_ = d.Set("protocol", edgePolicy.Protocol)
	}

	if edgePolicy.RuleAction != nil {
		_ = d.Set("rule_action", edgePolicy.RuleAction)
	}

	if edgePolicy.Port != nil {
		_ = d.Set("port", edgePolicy.Port)
	}

	if edgePolicy.Direction != nil {
		_ = d.Set("direction", edgePolicy.Direction)
	}

	if edgePolicy.Uuid != nil {
		_ = d.Set("uuid", edgePolicy.Uuid)
	}

	if edgePolicy.Enable != nil {
		_ = d.Set("enable", edgePolicy.Enable)
	}

	if edgePolicy.Description != nil {
		_ = d.Set("description", edgePolicy.Description)
	}

	if edgePolicy.Scope != nil {
		_ = d.Set("scope", edgePolicy.Scope)
	}

	if edgePolicy.ParamTemplateId != nil {
		_ = d.Set("param_template_id", edgePolicy.ParamTemplateId)
	}

	if edgePolicy.OrderIndex != nil {
		_ = d.Set("order_index", edgePolicy.OrderIndex)
	}

	return nil
}

func resourceTencentCloudCfwEdgePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_edge_policy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		request        = cfw.NewModifyAclRuleRequest()
		modifyRuleItem = cfw.CreateRuleItem{}
		uuid           = d.Id()
	)

	immutableArgs := []string{"uuid", "direction"}
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

	if v, ok := d.GetOk("scope"); ok {
		modifyRuleItem.Scope = helper.String(v.(string))
	}

	if v, ok := d.GetOk("param_template_id"); ok {
		modifyRuleItem.ParamTemplateId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("order_index"); ok {
		modifyRuleItem.OrderIndex = helper.IntInt64(v.(int))
	}

	request.Rules = append(request.Rules, &modifyRuleItem)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwClient().ModifyAclRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cfw edgePolicy failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCfwEdgePolicyRead(d, meta)
}

func resourceTencentCloudCfwEdgePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_edge_policy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		uuid    = d.Id()
	)

	if err := service.DeleteCfwEdgePolicyById(ctx, uuid); err != nil {
		return err
	}

	return nil
}
