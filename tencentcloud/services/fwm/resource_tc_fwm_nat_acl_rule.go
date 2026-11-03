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

func ResourceTencentCloudFwmNatAclRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudFwmNatAclRuleCreate,
		Read:   resourceTencentCloudFwmNatAclRuleRead,
		Update: resourceTencentCloudFwmNatAclRuleUpdate,
		Delete: resourceTencentCloudFwmNatAclRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "NAT ACL rule group ID. The group is the parent container and must already exist. Changing it will recreate the rule.",
			},

			"direction": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Rule direction: 0-outbound, 1-inbound. It is also part of the resource id, so changing it will recreate the rule.",
			},

			"source_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Source address content.",
			},

			"source_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Source type: ip/url/template/instance/tag.",
			},

			"target_content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target address content.",
			},

			"target_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target type: ip/url/template/instance/tag.",
			},

			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Protocol: TCP/UDP/ICMP/ANY/HTTP/HTTPS/DNS/FTP...",
			},

			"rule_action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Action: accept/drop/log.",
			},

			"order_index": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Rule priority, starting from 1.",
			},

			"scope": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Effective range: ALL-global, ap-xxx-region, cfwnat-xxx-NAT firewall instance.",
			},

			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Rule ID. It is filled by the server after creation and is part of the resource id `group_id#rule_id#direction`.",
			},

			"port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Port, e.g. 80, 80,443, 80-90; leave empty when protocol is ICMP.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule description.",
			},

			"param_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Port template ID.",
			},

			"belong_member_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Member account ID the rule belongs to. Required when scope is cfwnat-xxx or source_type/target_type is instance/tag.",
			},

			// computed
			"sequence": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Rule priority returned by the server, the readable counterpart of `order_index`.",
			},

			"source_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source asset name.",
			},

			"target_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target asset name.",
			},

			"scope_desc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Effective range description.",
			},

			"fw_ins_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Firewall instance ID.",
			},

			"country_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Country name.",
			},

			"city_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "City name.",
			},

			"country_code": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Country code.",
			},

			"city_code": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "City code.",
			},

			"country_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Country key.",
			},

			"city_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "City key.",
			},

			"is_region": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether the rule is a region rule: 0-no, 1-yes.",
			},

			"cloud_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud vendor code.",
			},

			"is_cloud": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether the rule is a cloud vendor rule: 0-no, 1-yes.",
			},

			"param_template_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port template name.",
			},

			"invalid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether the rule is invalid: 0-valid, 1-invalid.",
			},

			"belong_member": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Member account the rule belongs to.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member app ID.",
						},
						"uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member UIN.",
						},
						"nickname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member nickname.",
						},
						"member_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member ID (empty in personal scenario).",
						},
						"node_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Department the member belongs to.",
						},
					},
				},
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Update time.",
			},

			"dns_parse_cnt": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of domains.",
			},

			"hit_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Accumulated hit count of the rule under all delivered member accounts. It is 0 when never hit.",
			},

			"last_hit_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Latest hit time of the rule under all delivered member accounts. It is \"-\" when never hit.",
			},
		},
	}
}

func resourceTencentCloudFwmNatAclRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_nat_acl_rule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = fwmv20250611.NewCreateNatAclRuleRequest()
		response  = fwmv20250611.NewCreateNatAclRuleResponse()
		groupId   string
		direction int64
	)

	groupId = d.Get("group_id").(string)
	request.GroupId = helper.String(groupId)
	direction = int64(d.Get("direction").(int))

	rule := fwmv20250611.NatAclRule{}
	rule.SourceContent = helper.String(d.Get("source_content").(string))
	rule.SourceType = helper.String(d.Get("source_type").(string))
	rule.TargetContent = helper.String(d.Get("target_content").(string))
	rule.TargetType = helper.String(d.Get("target_type").(string))
	rule.Protocol = helper.String(d.Get("protocol").(string))
	rule.RuleAction = helper.String(d.Get("rule_action").(string))
	rule.OrderIndex = helper.Int64(int64(d.Get("order_index").(int)))
	rule.Scope = helper.String(d.Get("scope").(string))
	rule.Direction = helper.Int64(direction)
	if v, ok := d.GetOk("rule_id"); ok {
		rule.RuleId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		rule.Port = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		rule.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("param_template_id"); ok {
		rule.ParamTemplateId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("belong_member_id"); ok {
		rule.BelongMemberId = helper.String(v.(string))
	}

	request.Rules = []*fwmv20250611.NatAclRule{&rule}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().CreateNatAclRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create fwm_nat_acl_rule failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create fwm_nat_acl_rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create fwm_nat_acl_rule success, logId=%s, group_id=%s, rule_count=%d", logId, logId, groupId, len(response.Response.RuleIds))
	if len(response.Response.RuleIds) == 0 {
		return fmt.Errorf("Create fwm_nat_acl_rule failed, RuleIds is empty, group_id=%s, logId=%s", groupId, logId)
	}

	// A single NatAclRule is sent per Create call, so exactly one rule id must come back.
	if len(response.Response.RuleIds) > 1 {
		return fmt.Errorf("Create fwm_nat_acl_rule failed, expected 1 rule id but got %d, group_id=%s, logId=%s", len(response.Response.RuleIds), groupId, logId)
	}

	if response.Response.RuleIds[0] == nil || *response.Response.RuleIds[0] == "" {
		return fmt.Errorf("Create fwm_nat_acl_rule failed, RuleIds[0] is empty, group_id=%s, logId=%s", groupId, logId)
	}

	ruleId := *response.Response.RuleIds[0]
	_ = d.Set("rule_id", ruleId)
	d.SetId(strings.Join([]string{groupId, ruleId, strconv.FormatInt(direction, 10)}, tccommon.FILED_SP))
	return resourceTencentCloudFwmNatAclRuleRead(d, meta)
}

func resourceTencentCloudFwmNatAclRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_nat_acl_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = NewFwmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	groupId, ruleId, direction, err := parseFwmNatAclRuleId(d.Id())
	if err != nil {
		return err
	}

	respData, err := service.DescribeFwmNatAclRuleById(ctx, groupId, ruleId, direction)
	if err != nil {
		if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			if sdkerr.Code == "ResourceNotFound" {
				log.Printf("[CRUD] fwm_nat_acl_rule id=%s, ResourceNotFound, remove it from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] fwm_nat_acl_rule id=%s, describe returned empty, remove it from state", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("group_id", groupId)
	_ = d.Set("rule_id", ruleId)
	_ = d.Set("direction", int(direction))

	if respData.RuleId != nil {
		_ = d.Set("rule_id", respData.RuleId)
	}

	if respData.Direction != nil {
		_ = d.Set("direction", respData.Direction)
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

	// `port` and `description` are clearable on the cloud side: an empty value must reset the
	// state to "" as well, otherwise the old value stays in state forever (permanent diff).
	if respData.Port != nil {
		_ = d.Set("port", respData.Port)
	} else {
		_ = d.Set("port", "")
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	} else {
		_ = d.Set("description", "")
	}

	if respData.Protocol != nil {
		_ = d.Set("protocol", respData.Protocol)
	}

	if respData.RuleAction != nil {
		_ = d.Set("rule_action", respData.RuleAction)
	}

	if respData.Scope != nil {
		_ = d.Set("scope", respData.Scope)
	}

	if respData.ScopeDesc != nil {
		_ = d.Set("scope_desc", respData.ScopeDesc)
	}

	if respData.FwInsId != nil {
		_ = d.Set("fw_ins_id", respData.FwInsId)
	}

	if respData.CountryName != nil {
		_ = d.Set("country_name", respData.CountryName)
	}

	if respData.CityName != nil {
		_ = d.Set("city_name", respData.CityName)
	}

	if respData.CountryCode != nil {
		_ = d.Set("country_code", respData.CountryCode)
	}

	if respData.CityCode != nil {
		_ = d.Set("city_code", respData.CityCode)
	}

	if respData.CountryKey != nil {
		_ = d.Set("country_key", respData.CountryKey)
	}

	if respData.CityKey != nil {
		_ = d.Set("city_key", respData.CityKey)
	}

	if respData.IsRegion != nil {
		_ = d.Set("is_region", respData.IsRegion)
	}

	if respData.CloudCode != nil {
		_ = d.Set("cloud_code", respData.CloudCode)
	}

	if respData.IsCloud != nil {
		_ = d.Set("is_cloud", respData.IsCloud)
	}

	// `param_template_id` / `belong_member_id` are clearable on the cloud side as well.
	if respData.ParamTemplateId != nil {
		_ = d.Set("param_template_id", respData.ParamTemplateId)
	} else {
		_ = d.Set("param_template_id", "")
	}

	if respData.ParamTemplateName != nil {
		_ = d.Set("param_template_name", respData.ParamTemplateName)
	}

	if respData.Invalid != nil {
		_ = d.Set("invalid", respData.Invalid)
	}

	var memberId string
	if respData.BelongMember != nil {
		if respData.BelongMember.MemberId != nil {
			memberId = *respData.BelongMember.MemberId
		}

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
	} else {
		_ = d.Set("belong_member", []interface{}{})
	}

	if memberId != "" {
		_ = d.Set("belong_member_id", memberId)
	} else {
		_ = d.Set("belong_member_id", "")
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}

	if respData.DnsParseCnt != nil {
		_ = d.Set("dns_parse_cnt", respData.DnsParseCnt)
	}

	if respData.HitCount != nil {
		_ = d.Set("hit_count", respData.HitCount)
	}

	if respData.LastHitTime != nil {
		_ = d.Set("last_hit_time", respData.LastHitTime)
	}

	// Sequence is the readable counterpart of order_index, keep the configured value in state
	// when the server has nothing to report and expose the server value as `sequence`.
	if respData.Sequence != nil {
		_ = d.Set("sequence", respData.Sequence)
		_ = d.Set("order_index", int(*respData.Sequence))
	}

	return nil
}

func resourceTencentCloudFwmNatAclRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_nat_acl_rule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = fwmv20250611.NewModifyNatAclRuleRequest()
	)

	groupId, ruleId, direction, err := parseFwmNatAclRuleId(d.Id())
	if err != nil {
		return err
	}

	request.GroupId = helper.String(groupId)

	// ModifyNatAclRule overwrites the whole rule, so every writable field is always sent,
	// including the ones the user cleared (empty string means clear on the cloud side).
	rule := fwmv20250611.NatAclRule{}
	rule.RuleId = helper.String(ruleId)
	rule.Direction = helper.Int64(direction)
	rule.SourceContent = helper.String(d.Get("source_content").(string))
	rule.SourceType = helper.String(d.Get("source_type").(string))
	rule.TargetContent = helper.String(d.Get("target_content").(string))
	rule.TargetType = helper.String(d.Get("target_type").(string))
	rule.Protocol = helper.String(d.Get("protocol").(string))
	rule.RuleAction = helper.String(d.Get("rule_action").(string))
	rule.OrderIndex = helper.Int64(int64(d.Get("order_index").(int)))
	rule.Scope = helper.String(d.Get("scope").(string))
	rule.Port = helper.String(d.Get("port").(string))
	rule.Description = helper.String(d.Get("description").(string))
	rule.ParamTemplateId = helper.String(d.Get("param_template_id").(string))
	rule.BelongMemberId = helper.String(d.Get("belong_member_id").(string))
	request.Rule = &rule

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().ModifyNatAclRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify fwm_nat_acl_rule failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s modify fwm_nat_acl_rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudFwmNatAclRuleRead(d, meta)
}

func resourceTencentCloudFwmNatAclRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_nat_acl_rule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = fwmv20250611.NewDeleteNatAclRuleRequest()
	)

	groupId, ruleId, _, err := parseFwmNatAclRuleId(d.Id())
	if err != nil {
		return err
	}

	request.GroupId = helper.String(groupId)
	request.RuleIds = []*string{helper.String(ruleId)}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmV20250611Client().DeleteNatAclRuleWithContext(ctx, request)
		if e != nil {
			if sdkerr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "ResourceNotFound" {
					return nil
				}
			}
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete fwm_nat_acl_rule failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete fwm_nat_acl_rule failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

// parseFwmNatAclRuleId splits the resource id <groupId>#<ruleId>#<direction>.
func parseFwmNatAclRuleId(id string) (groupId, ruleId string, direction int64, err error) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 3 {
		return "", "", 0, fmt.Errorf("id is broken, %s", id)
	}

	if parts[0] == "" || parts[1] == "" {
		return "", "", 0, fmt.Errorf("id is broken, group_id or rule_id is empty in %s", id)
	}

	direction, err = strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return "", "", 0, fmt.Errorf("id is broken, invalid direction in %s: %s", id, err.Error())
	}

	return parts[0], parts[1], direction, nil
}
