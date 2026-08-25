package vpc

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcTrafficMirrorFilterRulesV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcTrafficMirrorFilterRulesV2Create,
		Read:   resourceTencentCloudVpcTrafficMirrorFilterRulesV2Read,
		Update: resourceTencentCloudVpcTrafficMirrorFilterRulesV2Update,
		Delete: resourceTencentCloudVpcTrafficMirrorFilterRulesV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"traffic_mirror_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Traffic mirror instance ID.",
			},

			"ingress_filter_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Ingress filter rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"src_net": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Source network segment of filter rule.",
						},
						"dst_net": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Destination network segment of filter rule.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Protocol of filter rule.",
						},
						"src_port": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Source port of filter rule, default 1-65535.",
						},
						"dst_port": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Destination port of filter rule, default 1-65535.",
						},
						"traffic_mirror_filter_rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Traffic mirror filter rule ID.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Traffic mirror filter rule priority.",
						},
						"action": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Traffic mirror filter rule policy, support types: `ACCEPT`, `DROP`.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Traffic mirror filter rule description.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
					},
				},
			},

			"egress_filter_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Egress filter rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"src_net": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Source network segment of filter rule.",
						},
						"dst_net": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Destination network segment of filter rule.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Protocol of filter rule.",
						},
						"src_port": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Source port of filter rule, default 1-65535.",
						},
						"dst_port": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Destination port of filter rule, default 1-65535.",
						},
						"traffic_mirror_filter_rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Traffic mirror filter rule ID.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Traffic mirror filter rule priority.",
						},
						"action": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Traffic mirror filter rule policy, support types: `ACCEPT`, `DROP`.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Traffic mirror filter rule description.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudVpcTrafficMirrorFilterRulesV2Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_traffic_mirror_filter_rules_v2.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request         = vpc.NewCreateTrafficMirrorFilterRulesRequest()
		trafficMirrorId string
	)

	if v, ok := d.GetOk("traffic_mirror_id"); ok {
		trafficMirrorId = v.(string)
		request.TrafficMirrorId = helper.String(trafficMirrorId)
	}

	if v, ok := d.GetOk("ingress_filter_rules"); ok {
		ingressRules := expandTrafficMirrorFilterRules(v.([]interface{}))
		request.IngressFilterRules = ingressRules
	}

	if v, ok := d.GetOk("egress_filter_rules"); ok {
		egressRules := expandTrafficMirrorFilterRules(v.([]interface{}))
		request.EgressFilterRules = egressRules
	}

	hasIngress := d.Get("ingress_filter_rules").([]interface{})
	hasEgress := d.Get("egress_filter_rules").([]interface{})
	if len(hasIngress) == 0 && len(hasEgress) == 0 {
		return fmt.Errorf("At least one of `ingress_filter_rules` or `egress_filter_rules` must be provided.")
	}

	var response *vpc.CreateTrafficMirrorFilterRulesResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateTrafficMirrorFilterRulesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create vpc traffic_mirror_filter_rules_v2 failed, reason:%+v", logId, err)
		return err
	}

	if response == nil || response.Response == nil {
		log.Printf("[CRITAL]%s create vpc traffic_mirror_filter_rules_v2 failed, Response is nil.", logId)
		return fmt.Errorf("CreateTrafficMirrorFilterRules response is nil.")
	}

	log.Printf("[DEBUG]%s api[%s] success, traffic_mirror_id=%s", logId, request.GetAction(), trafficMirrorId)
	d.SetId(trafficMirrorId)

	return resourceTencentCloudVpcTrafficMirrorFilterRulesV2Read(d, meta)
}

func resourceTencentCloudVpcTrafficMirrorFilterRulesV2Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_traffic_mirror_filter_rules_v2.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		trafficMirrorId = d.Id()
		request         = vpc.NewDescribeTrafficMirrorFilterRulesRequest()
	)

	request.TrafficMirrorId = helper.String(trafficMirrorId)
	request.Limit = helper.IntUint64(100)

	var response *vpc.DescribeTrafficMirrorFilterRulesResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeTrafficMirrorFilterRulesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read vpc traffic_mirror_filter_rules_v2 failed, reason:%+v", logId, err)
		return err
	}

	if response == nil || response.Response == nil {
		log.Printf("[CRUD]%s resource `tencentcloud_vpc_traffic_mirror_filter_rules_v2` response is nil, id=%s", logId, d.Id())
		d.SetId("")
		return nil
	}

	ingressRules := response.Response.IngressFilterRules
	egressRules := response.Response.EgressFilterRules

	if (ingressRules == nil || len(ingressRules) == 0) && (egressRules == nil || len(egressRules) == 0) {
		log.Printf("[CRUD]%s resource `tencentcloud_vpc_traffic_mirror_filter_rules_v2` rules are empty, id=%s", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("traffic_mirror_id", trafficMirrorId)

	if ingressRules != nil {
		_ = d.Set("ingress_filter_rules", flattenTrafficMirrorFilterRules(ingressRules))
	}

	if egressRules != nil {
		_ = d.Set("egress_filter_rules", flattenTrafficMirrorFilterRules(egressRules))
	}

	return nil
}

func resourceTencentCloudVpcTrafficMirrorFilterRulesV2Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_traffic_mirror_filter_rules_v2.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request         = vpc.NewModifyTrafficMirrorFilterRulesRequest()
		trafficMirrorId = d.Id()
	)

	request.TrafficMirrorId = helper.String(trafficMirrorId)

	needChange := false
	if d.HasChange("ingress_filter_rules") {
		needChange = true
		if v, ok := d.GetOk("ingress_filter_rules"); ok {
			request.IngressFilterRules = expandTrafficMirrorFilterRules(v.([]interface{}))
		} else {
			request.IngressFilterRules = []*vpc.TrafficMirrorFilter{}
		}
	}

	if d.HasChange("egress_filter_rules") {
		needChange = true
		if v, ok := d.GetOk("egress_filter_rules"); ok {
			request.EgressFilterRules = expandTrafficMirrorFilterRules(v.([]interface{}))
		} else {
			request.EgressFilterRules = []*vpc.TrafficMirrorFilter{}
		}
	}

	if needChange {
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ModifyTrafficMirrorFilterRulesWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			_ = result
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update vpc traffic_mirror_filter_rules_v2 failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudVpcTrafficMirrorFilterRulesV2Read(d, meta)
}

func resourceTencentCloudVpcTrafficMirrorFilterRulesV2Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_traffic_mirror_filter_rules_v2.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request         = vpc.NewDeleteTrafficMirrorFilterRulesRequest()
		trafficMirrorId = d.Id()
	)

	request.TrafficMirrorId = helper.String(trafficMirrorId)

	ingressRuleIds := make([]*string, 0)
	if v, ok := d.GetOk("ingress_filter_rules"); ok {
		for _, item := range v.([]interface{}) {
			ruleMap := item.(map[string]interface{})
			if ruleId, ok := ruleMap["traffic_mirror_filter_rule_id"].(string); ok && ruleId != "" {
				ingressRuleIds = append(ingressRuleIds, helper.String(ruleId))
			}
		}
	}
	request.IngressFilterRuleIds = ingressRuleIds

	egressRuleIds := make([]*string, 0)
	if v, ok := d.GetOk("egress_filter_rules"); ok {
		for _, item := range v.([]interface{}) {
			ruleMap := item.(map[string]interface{})
			if ruleId, ok := ruleMap["traffic_mirror_filter_rule_id"].(string); ok && ruleId != "" {
				egressRuleIds = append(egressRuleIds, helper.String(ruleId))
			}
		}
	}
	request.EgressFilterRuleIds = egressRuleIds

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DeleteTrafficMirrorFilterRulesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		_ = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete vpc traffic_mirror_filter_rules_v2 failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

func expandTrafficMirrorFilterRules(rules []interface{}) []*vpc.TrafficMirrorFilter {
	result := make([]*vpc.TrafficMirrorFilter, 0, len(rules))
	for _, item := range rules {
		ruleMap := item.(map[string]interface{})
		rule := &vpc.TrafficMirrorFilter{}

		if v, ok := ruleMap["src_net"].(string); ok && v != "" {
			rule.SrcNet = helper.String(v)
		}

		if v, ok := ruleMap["dst_net"].(string); ok && v != "" {
			rule.DstNet = helper.String(v)
		}

		if v, ok := ruleMap["protocol"].(string); ok && v != "" {
			rule.Protocol = helper.String(v)
		}

		if v, ok := ruleMap["src_port"].(string); ok && v != "" {
			rule.SrcPort = helper.String(v)
		}

		if v, ok := ruleMap["dst_port"].(string); ok && v != "" {
			rule.DstPort = helper.String(v)
		}

		if v, ok := ruleMap["traffic_mirror_filter_rule_id"].(string); ok && v != "" {
			rule.TrafficMirrorFilterRuleId = helper.String(v)
		}

		if v, ok := ruleMap["priority"].(int); ok {
			rule.Priority = helper.IntUint64(v)
		}

		if v, ok := ruleMap["action"].(string); ok && v != "" {
			rule.Action = helper.String(v)
		}

		if v, ok := ruleMap["description"].(string); ok && v != "" {
			rule.Description = helper.String(v)
		}

		result = append(result, rule)
	}
	return result
}

func flattenTrafficMirrorFilterRules(rules []*vpc.TrafficMirrorFilter) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		ruleMap := map[string]interface{}{}

		if rule.SrcNet != nil {
			ruleMap["src_net"] = rule.SrcNet
		}

		if rule.DstNet != nil {
			ruleMap["dst_net"] = rule.DstNet
		}

		if rule.Protocol != nil {
			ruleMap["protocol"] = rule.Protocol
		}

		if rule.SrcPort != nil {
			ruleMap["src_port"] = rule.SrcPort
		}

		if rule.DstPort != nil {
			ruleMap["dst_port"] = rule.DstPort
		}

		if rule.TrafficMirrorFilterRuleId != nil {
			ruleMap["traffic_mirror_filter_rule_id"] = rule.TrafficMirrorFilterRuleId
		}

		if rule.Priority != nil {
			ruleMap["priority"] = rule.Priority
		}

		if rule.Action != nil {
			ruleMap["action"] = rule.Action
		}

		if rule.Description != nil {
			ruleMap["description"] = rule.Description
		}

		if rule.CreatedTime != nil {
			ruleMap["created_time"] = rule.CreatedTime
		}

		result = append(result, ruleMap)
	}
	return result
}
