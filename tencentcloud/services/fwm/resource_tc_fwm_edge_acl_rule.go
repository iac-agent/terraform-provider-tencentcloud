package fwm

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	fwmv20250611 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var fwmEdgeAclRuleMutableArgs = []string{
	"source_content",
	"target_content",
	"port",
	"protocol",
	"rule_action",
	"description",
	"scope",
	"belong_member_id",
	"param_template_id",
}

// fwmEdgeAclRuleReplaceArgs returns the changed arguments that can only be
// applied by recreating the rule.
//
// group_id / direction are parts of the resource id (`group_id#rule_id#direction`)
// used by every Read/Update/Delete call, so they always force a replacement.
// source_type / target_type are mutable through ModifyEdgeAclRule, but when the
// type changes alone while belong_member_id is configured, the configured member
// account may no longer match the new address type, so the rule is recreated.
func fwmEdgeAclRuleReplaceArgs(d *schema.ResourceDiff) []string {
	args := make([]string, 0, 2)

	for _, arg := range []string{"group_id", "direction"} {
		if d.HasChange(arg) {
			args = append(args, arg)
		}
	}

	if _, ok := d.GetOk("belong_member_id"); !ok {
		return args
	}

	typePairs := [][2]string{
		{"source_type", "source_content"},
		{"target_type", "target_content"},
	}
	for _, pair := range typePairs {
		if d.HasChange(pair[0]) && !d.HasChange(pair[1]) {
			args = append(args, pair[0])
		}
	}

	return args
}

func ResourceTencentCloudFwmEdgeAclRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudFwmEdgeAclRuleCreate,
		Read:   resourceTencentCloudFwmEdgeAclRuleRead,
		Update: resourceTencentCloudFwmEdgeAclRuleUpdate,
		Delete: resourceTencentCloudFwmEdgeAclRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
			// ForceNew cannot be declared on a schema attribute for a conditional
			// replacement, so the replacement decision is made here.
			for _, arg := range fwmEdgeAclRuleReplaceArgs(d) {
				if err := d.ForceNew(arg); err != nil {
					return err
				}
			}

			return nil
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Edge ACL rule group ID. The rule group is the parent container and must already exist before creating the rule. " +
					"The rule is always queried within this group, so changing it recreates the rule.",
			},

			"direction": {
				Type:     schema.TypeInt,
				Required: true,
				Description: "Rule direction: 0-egress, 1-ingress. The direction is part of the resource id (`group_id#rule_id#direction`), " +
					"so changing it recreates the rule.",
			},

			"source_content": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Source address content. Format depends on source_type: IP/CIDR for ip, domain for domain, template ID for template, " +
					"instance ID list (comma separated) for instance, tag key:value for tag.",
			},

			"source_type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Source address type: ip, domain, template, instance, tag. It can be modified in place; it only recreates the rule when it is changed " +
					"alone while belong_member_id is configured, because the member account may no longer match the new address type.",
			},

			"target_content": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Target address content. Format is the same as source_content.",
			},

			"target_type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Target address type: ip, domain, template, instance, tag. It can be modified in place; it only recreates the rule when it is changed " +
					"alone while belong_member_id is configured, because the member account may no longer match the new address type.",
			},

			"port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Port. Supports single port, port range and comma separated multi ports, such as 80, 1-65535, 80,443,8080.",
			},

			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Protocol type: TCP, UDP, ICMP, ANY.",
			},

			"rule_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule action: accept-pass, drop-block, log-observe.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule description, length 0-256 characters.",
			},

			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Effective scope: serial, side, all.",
			},

			"belong_member_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Member account ID that the rule belongs to (used in multi-account scenarios). Required when source_type or target_type is instance or tag.",
			},

			"param_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Parameter template ID.",
			},

			// computed
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule ID, it is part of the resource id `group_id#rule_id#direction`.",
			},

			"sequence": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Rule execution sequence returned by the server (the read side alias of the create `OrderIndex`).",
			},

			"source_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source address name returned by the server, only set when source_type is template, instance or tag.",
			},

			"target_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target address name returned by the server, only set when target_type is template, instance or tag.",
			},

			"param_template_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Parameter template name returned by the server, only set when the rule uses a parameter template.",
			},

			"country_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region name level 1, returned for region rules (used by regular matching).",
			},

			"city_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region name level 2, returned for region rules (used by regular matching).",
			},

			"country_code": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Country ID, returned for region rules.",
			},

			"city_code": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "City ID, returned for region rules.",
			},

			"country_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Country short name, returned for region rules.",
			},

			"city_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Province/city short name, returned for region rules.",
			},

			"is_region": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "0 means a normal rule, 1 means a region rule.",
			},

			"cloud_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud vendor code, returned for cloud vendor rules.",
			},

			"is_cloud": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "0 means a normal rule, 1 means a cloud vendor rule.",
			},

			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance name returned by the server.",
			},

			"invalid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether the rule is invalid: 0-valid, 1-invalid.",
			},

			"dns_parse_cnt": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of parsed domains.",
			},

			"hit_count": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "Accumulated hit count of the rule over all delivered member accounts. It is the Cloud Firewall value without any time window conversion, " +
					"and 0 when the rule never hit.",
			},

			"last_hit_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Latest hit time over all delivered member accounts, \"-\" when the rule never hit.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule creation time.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule last update time.",
			},

			"belong_member": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Member account that the rule belongs to.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member AppId.",
						},
						"uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member Uin.",
						},
						"nickname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member nickname.",
						},
						"member_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member ID, empty in the personal account scenario.",
						},
						"node_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Department the member belongs to.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudFwmEdgeAclRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_edge_acl_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = fwmv20250611.NewCreateEdgeAclRuleRequest()
		response  = fwmv20250611.NewCreateEdgeAclRuleResponse()
		groupId   string
		ruleId    string
		direction int64
	)

	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = helper.String(v.(string))
		groupId = v.(string)
	}

	rule := fwmv20250611.EdgeAclRuleInfo{}
	if v, ok := d.GetOkExists("direction"); ok {
		direction = int64(v.(int))
		rule.Direction = helper.Int64(direction)
	}

	if v, ok := d.GetOk("source_content"); ok {
		rule.SourceContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_type"); ok {
		rule.SourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_content"); ok {
		rule.TargetContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_type"); ok {
		rule.TargetType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		rule.Port = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol"); ok {
		rule.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_action"); ok {
		rule.RuleAction = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		rule.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("scope"); ok {
		rule.Scope = helper.String(v.(string))
	}

	if v, ok := d.GetOk("belong_member_id"); ok {
		rule.BelongMemberId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("param_template_id"); ok {
		rule.ParamTemplateId = helper.String(v.(string))
	}

	// CreateEdgeAclRule requires OrderIndex to start from 1 and increase strictly,
	// a single managed rule therefore always creates with OrderIndex 1.
	rule.OrderIndex = helper.IntInt64(1)
	request.Rules = []*fwmv20250611.EdgeAclRuleInfo{&rule}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().CreateEdgeAclRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create tencentcloud_fwm_edge_acl_rule failed, Response is nil."))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create tencentcloud_fwm_edge_acl_rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create tencentcloud_fwm_edge_acl_rule success, logId=%s, rule_ids=%v", logId, logId, response.Response.RuleIds)

	if len(response.Response.RuleIds) != 1 {
		return fmt.Errorf("Create tencentcloud_fwm_edge_acl_rule failed, expect exactly 1 rule id but got %d, group_id=%s", len(response.Response.RuleIds), groupId)
	}

	if response.Response.RuleIds[0] == nil || *response.Response.RuleIds[0] == "" {
		return fmt.Errorf("Create tencentcloud_fwm_edge_acl_rule failed, rule id is empty, group_id=%s", groupId)
	}

	ruleId = *response.Response.RuleIds[0]
	d.SetId(strings.Join([]string{groupId, ruleId, strconv.FormatInt(direction, 10)}, tccommon.FILED_SP))
	log.Printf("[DEBUG]%s tencentcloud_fwm_edge_acl_rule created, logId=%s, id=%s", logId, logId, d.Id())

	return resourceTencentCloudFwmEdgeAclRuleRead(d, meta)
}

func resourceTencentCloudFwmEdgeAclRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_edge_acl_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = NewFwmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	groupId, ruleId, direction, err := parseFwmEdgeAclRuleId(d.Id())
	if err != nil {
		return err
	}

	respData, err := service.DescribeFwmEdgeAclRuleById(ctx, groupId, ruleId, direction)
	if err != nil {
		if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok && sdkErr.Code == "ResourceNotFound" {
			log.Printf("[CRUD] tencentcloud_fwm_edge_acl_rule not found, clear id, id=%s", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] tencentcloud_fwm_edge_acl_rule id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if respData.RuleId == nil || *respData.RuleId != ruleId {
		log.Printf("[CRUD] tencentcloud_fwm_edge_acl_rule id=%s, describe returns unexpected rule, skip state update", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("group_id", groupId)

	if respData.RuleId != nil {
		_ = d.Set("rule_id", respData.RuleId)
	}

	if respData.Direction != nil {
		_ = d.Set("direction", int(*respData.Direction))
	}

	if respData.Sequence != nil {
		_ = d.Set("sequence", int(*respData.Sequence))
	}

	if respData.SourceContent != nil {
		_ = d.Set("source_content", respData.SourceContent)
	}

	if respData.SourceType != nil {
		_ = d.Set("source_type", respData.SourceType)
	}

	if respData.SourceName != nil {
		_ = d.Set("source_name", respData.SourceName)
	}

	if respData.TargetContent != nil {
		_ = d.Set("target_content", respData.TargetContent)
	}

	if respData.TargetType != nil {
		_ = d.Set("target_type", respData.TargetType)
	}

	if respData.TargetName != nil {
		_ = d.Set("target_name", respData.TargetName)
	}

	if respData.Port != nil {
		_ = d.Set("port", respData.Port)
	}

	if respData.Protocol != nil {
		_ = d.Set("protocol", respData.Protocol)
	}

	if respData.RuleAction != nil {
		_ = d.Set("rule_action", respData.RuleAction)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.Scope != nil {
		_ = d.Set("scope", respData.Scope)
	}

	if respData.ParamTemplateId != nil {
		_ = d.Set("param_template_id", respData.ParamTemplateId)
	}

	if respData.ParamTemplateName != nil {
		_ = d.Set("param_template_name", respData.ParamTemplateName)
	}

	if respData.CountryName != nil {
		_ = d.Set("country_name", respData.CountryName)
	}

	if respData.CityName != nil {
		_ = d.Set("city_name", respData.CityName)
	}

	if respData.CountryCode != nil {
		_ = d.Set("country_code", int(*respData.CountryCode))
	}

	if respData.CityCode != nil {
		_ = d.Set("city_code", int(*respData.CityCode))
	}

	if respData.CountryKey != nil {
		_ = d.Set("country_key", respData.CountryKey)
	}

	if respData.CityKey != nil {
		_ = d.Set("city_key", respData.CityKey)
	}

	if respData.IsRegion != nil {
		_ = d.Set("is_region", int(*respData.IsRegion))
	}

	if respData.CloudCode != nil {
		_ = d.Set("cloud_code", respData.CloudCode)
	}

	if respData.IsCloud != nil {
		_ = d.Set("is_cloud", int(*respData.IsCloud))
	}

	if respData.InstanceName != nil {
		_ = d.Set("instance_name", respData.InstanceName)
	}

	if respData.Invalid != nil {
		_ = d.Set("invalid", int(*respData.Invalid))
	}

	if respData.DnsParseCnt != nil {
		_ = d.Set("dns_parse_cnt", int(*respData.DnsParseCnt))
	}

	if respData.HitCount != nil {
		_ = d.Set("hit_count", int(*respData.HitCount))
	}

	if respData.LastHitTime != nil {
		_ = d.Set("last_hit_time", respData.LastHitTime)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}

	if respData.BelongMember != nil {
		belongMember := map[string]interface{}{
			"app_id":    "",
			"uin":       "",
			"nickname":  "",
			"member_id": "",
			"node_name": "",
		}
		if respData.BelongMember.AppId != nil {
			belongMember["app_id"] = *respData.BelongMember.AppId
		}
		if respData.BelongMember.Uin != nil {
			belongMember["uin"] = *respData.BelongMember.Uin
		}
		if respData.BelongMember.Nickname != nil {
			belongMember["nickname"] = *respData.BelongMember.Nickname
		}
		if respData.BelongMember.MemberId != nil {
			belongMember["member_id"] = *respData.BelongMember.MemberId
		}
		if respData.BelongMember.NodeName != nil {
			belongMember["node_name"] = *respData.BelongMember.NodeName
		}

		_ = d.Set("belong_member", []interface{}{belongMember})

		if respData.BelongMember.MemberId != nil {
			_ = d.Set("belong_member_id", respData.BelongMember.MemberId)
		}
	}

	return nil
}

func resourceTencentCloudFwmEdgeAclRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_edge_acl_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	groupId, ruleId, direction, err := parseFwmEdgeAclRuleId(d.Id())
	if err != nil {
		return err
	}

	// Every field of the rule is expressible by ModifyEdgeAclRule, only the parts
	// of the composite id (group_id / rule_id / direction) are immutable.
	immutableArgs := []string{"group_id", "rule_id"}
	for _, arg := range immutableArgs {
		if d.HasChange(arg) {
			return fmt.Errorf("argument `%s` of tencentcloud_fwm_edge_acl_rule is immutable and cannot be updated, id=%s", arg, d.Id())
		}
	}

	if d.HasChange("direction") {
		return fmt.Errorf("argument `direction` of tencentcloud_fwm_edge_acl_rule is a part of the resource id and cannot be updated in place, id=%s", d.Id())
	}

	needChange := false
	for _, arg := range fwmEdgeAclRuleMutableArgs {
		if d.HasChange(arg) {
			needChange = true
			break
		}
	}

	if d.HasChange("source_type") || d.HasChange("target_type") {
		needChange = true
	}

	if !needChange {
		return resourceTencentCloudFwmEdgeAclRuleRead(d, meta)
	}

	request := fwmv20250611.NewModifyEdgeAclRuleRequest()
	rule := fwmv20250611.EdgeAclRuleInfo{}

	if v, ok := d.GetOk("source_content"); ok {
		rule.SourceContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_type"); ok {
		rule.SourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_content"); ok {
		rule.TargetContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_type"); ok {
		rule.TargetType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		rule.Port = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol"); ok {
		rule.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_action"); ok {
		rule.RuleAction = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		rule.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("scope"); ok {
		rule.Scope = helper.String(v.(string))
	}

	if v, ok := d.GetOk("belong_member_id"); ok {
		rule.BelongMemberId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("param_template_id"); ok {
		rule.ParamTemplateId = helper.String(v.(string))
	}

	rule.RuleId = helper.String(ruleId)
	rule.Direction = helper.Int64(direction)
	// OrderIndex follows the current server side Sequence, otherwise Modify would
	// reset the rule priority to a fixed value.
	if v, ok := d.GetOkExists("sequence"); ok {
		rule.OrderIndex = helper.Int64(int64(v.(int)))
	}

	request.GroupId = helper.String(groupId)
	request.Rule = &rule

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().ModifyEdgeAclRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify tencentcloud_fwm_edge_acl_rule failed, Response is nil."))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s modify tencentcloud_fwm_edge_acl_rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudFwmEdgeAclRuleRead(d, meta)
}

func resourceTencentCloudFwmEdgeAclRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_edge_acl_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = fwmv20250611.NewDeleteEdgeAclRuleRequest()
	)

	groupId, ruleId, _, err := parseFwmEdgeAclRuleId(d.Id())
	if err != nil {
		return err
	}

	// DeleteEdgeAclRule deletes the given RuleIds only, a single managed rule is
	// always deleted with an one-element list. A rule that already disappeared is
	// treated as deleted so that terraform destroy stays idempotent.
	request.GroupId = helper.String(groupId)
	request.RuleIds = []*string{helper.String(ruleId)}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().DeleteEdgeAclRuleWithContext(ctx, request)
		if e != nil {
			if sdkErr, ok := e.(*sdkErrors.TencentCloudSDKError); ok && sdkErr.Code == "ResourceNotFound" {
				log.Printf("[WARN]%s delete tencentcloud_fwm_edge_acl_rule, rule not found, treat as success, id=%s", logId, d.Id())
				return nil
			}
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete tencentcloud_fwm_edge_acl_rule failed, Response is nil."))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete tencentcloud_fwm_edge_acl_rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

// parseFwmEdgeAclRuleId splits the resource id <groupId>#<ruleId>#<direction>.
func parseFwmEdgeAclRuleId(id string) (groupId, ruleId string, direction int64, err error) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 3 {
		return "", "", 0, fmt.Errorf("id is broken, %s", id)
	}

	direction, err = strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return "", "", 0, fmt.Errorf("id is broken, invalid direction in %s: %s", id, err.Error())
	}

	return parts[0], parts[1], direction, nil
}
