package vpc

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudSecurityGroupRuleSet() *schema.Resource {
	ruleElem := map[string]*schema.Schema{
		"action": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: tccommon.ValidateAllowedStringValueIgnoreCase([]string{"ACCEPT", "DROP"}),
			Description:  "Rule 策略 的 安全 组. 有效值：`ACCEPT` 和 `DROP`。",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "描述 安全 组 规则。",
		},
		"cidr_block": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "An IP 地址 网络 或 CIDR segment. NOTE: `cidr_block`，`ipv6_cidr_block`，`source_security_id` 和 `address_template_*` 是 exclusive 和 不能 是 集合 在 same 时间; One 的 `cidr_block`，`ipv6_cidr_block`，`source_security_id` 和 `address_template_*` 必须 是 集合。",
		},
		"ipv6_cidr_block": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "An IPV6 地址 网络 或 CIDR segment，和 conflict 使用 `source_security_id` 和 `address_template_*`. NOTE: One 的 `cidr_block`，`ipv6_cidr_block`，`source_security_id` 和 `address_template_*` 必须 是 集合。",
		},
		"source_security_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "ID nested 安全 组，和 conflicts 使用 `cidr_block` 和 `address_template_*`. NOTE: One 的 `cidr_block`，`ipv6_cidr_block`，`source_security_id` 和 `address_template_*` 必须 是 集合。",
		},
		"address_template_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "指定Address 模板 ID like `ipm-xxxxxxxx`，conflict 使用 `source_security_id` 和 `cidr_block`. NOTE: One 的 `cidr_block`，`ipv6_cidr_block`，`source_security_id` 和 `address_template_*` 必须 是 集合。",
		},
		"address_template_group": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "指定Group ID 地址 template like `ipmg-xxxxxxxx`，conflict 使用 `source_security_id` 和 `cidr_block`. NOTE: One 的 `cidr_block`，`ipv6_cidr_block`，`source_security_id` 和 `address_template_*` 必须 是 集合。",
		},
		"service_template_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "指定Protocol 模板 ID like `ppm-xxxxxxxx`，conflict 使用 `协议` 和 `端口`。",
		},
		"service_template_group": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "指定Group ID 协议 模板 ID like `ppmg-xxxxxxxx`，conflict 使用 `协议` 和 `端口`。",
		},
		"protocol": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "类型 IP 协议 有效值：`TCP`，`UDP`，`ICMP`，`ICMPv6` 和 `ALL`. 默认为 all types 协议，和 conflicts 使用 `service_template_*`。",
		},
		"port": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Range 的 端口 可用 值 可以 是 `all`， 单个 端口，或 端口 范围. E.g. `80`，`80,90`，`80-90` 或 `all`. 注意: 如果 `协议` 值 是 集合 到 `ALL`， `端口` 值 also needs 到 是 集合 到 `all`. 默认为 all ports，和 conflicts 使用 `service_template_*`。",
		},
		"policy_index": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "安全 组 规则 索引 数量，whose 值 dynamically changes 使用 changes 在 安全 组 规则。",
		},
	}
	return &schema.Resource{
		Create: resourceTencentCloudSecurityGroupRuleSetCreate,
		Read:   resourceTencentCloudSecurityGroupRuleSetRead,
		Update: resourceTencentCloudSecurityGroupRuleSetUpdate,
		Delete: resourceTencentCloudSecurityGroupRuleSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID 安全 组 到 是 queried。",
			},
			"ingress": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "列表 ingress 规则. NOTE: 此 block 是 ordered， first 规则 has highest 优先级",
				Elem:        &schema.Resource{Schema: ruleElem},
			},
			"egress": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "列表 egress 规则. NOTE: 此 block 是 ordered， first 规则 has highest 优先级",
				Elem:        &schema.Resource{Schema: ruleElem},
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Security policies 版本，auto increment 对于 every update。",
			},
		},
	}
}

func resourceTencentCloudSecurityGroupRuleSetCreate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_security_group_rule_set.create")()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service         = VpcService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request         = vpc.NewModifySecurityGroupPoliciesRequest()
		securityGroupId string
		err             error
	)

	if v, ok := d.GetOk("security_group_id"); ok {
		request.SecurityGroupId = helper.String(v.(string))
		securityGroupId = v.(string)
	}

	request.SecurityGroupPolicySet = &vpc.SecurityGroupPolicySet{}

	if v, ok := d.GetOk("ingress"); ok {
		ingressRules := v.([]interface{})
		request.SecurityGroupPolicySet.Ingress, err = unmarshalSecurityPolicy(ingressRules)
		if err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("egress"); ok {
		egressRules := v.([]interface{})
		request.SecurityGroupPolicySet.Egress, err = unmarshalSecurityPolicy(egressRules)
		if err != nil {
			return err
		}
	}

	err = service.ModifySecurityGroupPolicies(ctx, request)
	if err != nil {
		return err
	}

	d.SetId(securityGroupId)
	return resourceTencentCloudSecurityGroupRuleSetRead(d, m)
}

func resourceTencentCloudSecurityGroupRuleSetRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_security_group_rule_set.read")()
	defer tccommon.InconsistentCheck(d, m)()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service         = VpcService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}
		securityGroupId = d.Id()
	)

	request := vpc.NewDescribeSecurityGroupPoliciesRequest()
	request.SecurityGroupId = &securityGroupId

	result, err := service.DescribeSecurityGroupPolicies(ctx, securityGroupId)
	if err != nil {
		return err
	}

	if result == nil {
		d.SetId("")
		return fmt.Errorf("resource `tencentcloud_security_group_rule_set` %s does not exist", d.Id())
	}

	_ = d.Set("security_group_id", securityGroupId)

	if result.Version != nil {
		_ = d.Set("version", result.Version)
	}

	if result.Ingress != nil {
		_ = d.Set("ingress", marshalSecurityPolicy(result.Ingress))
	}

	if result.Egress != nil {
		_ = d.Set("egress", marshalSecurityPolicy(result.Egress))
	}

	return nil
}

func resourceTencentCloudSecurityGroupRuleSetUpdate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("tencentcloud_security_group_rule_set.update")()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service         = VpcService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request         = vpc.NewModifySecurityGroupPoliciesRequest()
		securityGroupId = d.Id()
		nextVer         string
		needChange      bool
		err             error
	)

	mutableArgs := []string{"ingress", "egress"}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		version := d.Get("version").(string)
		ver, _ := strconv.ParseInt(version, 10, 64)
		ver += 1
		request.SecurityGroupId = helper.String(securityGroupId)
		request.SecurityGroupPolicySet = &vpc.SecurityGroupPolicySet{}
		request.SortPolicys = helper.Bool(true)

		oldIngress, newIngress := d.GetChange("ingress")
		oldEgress, newEgress := d.GetChange("egress")

		oldIngressList := oldIngress.([]interface{})
		newIngressList := newIngress.([]interface{})
		oldEgressList := oldEgress.([]interface{})
		newEgressList := newEgress.([]interface{})

		if len(newIngressList) == 0 && len(newEgressList) == 0 {
			ver = 0

		} else if len(newIngressList) != 0 && len(newEgressList) == 0 {
			request.SecurityGroupPolicySet.Ingress, err = unmarshalSecurityPolicy(newIngressList)
			if err != nil {
				return err
			}

			if len(oldEgressList) > 0 {
				tmpList := make([]*int64, 0)
				for _, v := range oldEgressList {
					item := v.(map[string]interface{})
					policyIndex := int64(item["policy_index"].(int))
					tmpList = append(tmpList, &policyIndex)
				}

				e := service.DeleteSecurityGroupPolicyByPolicyIndexList(ctx, securityGroupId, tmpList, "egress")
				if e != nil {
					return e
				}

				ver += 1
			}

		} else if len(newIngressList) == 0 && len(newEgressList) != 0 {
			request.SecurityGroupPolicySet.Egress, err = unmarshalSecurityPolicy(newEgressList)
			if err != nil {
				return err
			}

			if len(oldIngressList) > 0 {
				tmpList := make([]*int64, 0)
				for _, v := range oldIngressList {
					item := v.(map[string]interface{})
					policyIndex := int64(item["policy_index"].(int))
					tmpList = append(tmpList, &policyIndex)
				}

				e := service.DeleteSecurityGroupPolicyByPolicyIndexList(ctx, securityGroupId, tmpList, "ingress")
				if e != nil {
					return e
				}

				ver += 1
			}

		} else {
			request.SecurityGroupPolicySet.Ingress, err = unmarshalSecurityPolicy(newIngressList)
			if err != nil {
				return err
			}

			request.SecurityGroupPolicySet.Egress, err = unmarshalSecurityPolicy(newEgressList)
			if err != nil {
				return err
			}
		}

		nextVer = fmt.Sprintf("%d", ver)
		request.SecurityGroupPolicySet.Version = helper.String(nextVer)
		err = service.ModifySecurityGroupPolicies(ctx, request)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudSecurityGroupRuleSetRead(d, m)
}

func resourceTencentCloudSecurityGroupRuleSetDelete(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_security_group_rule_set.delete")()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service         = VpcService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request         = vpc.NewModifySecurityGroupPoliciesRequest()
		securityGroupId = d.Id()
	)

	request.SecurityGroupId = &securityGroupId
	request.SecurityGroupPolicySet = &vpc.SecurityGroupPolicySet{
		Version: helper.String("0"),
		Ingress: []*vpc.SecurityGroupPolicy{},
		Egress:  []*vpc.SecurityGroupPolicy{},
	}

	err := service.ModifySecurityGroupPolicies(ctx, request)
	if err != nil {
		log.Printf("[CRITAL]%s security group rule delete failed: %s\n ", logId, err.Error())
		return err
	}

	return nil
}

func unmarshalSecurityPolicy(policies []interface{}) (output []*vpc.SecurityGroupPolicy, err error) {
	for i := range policies {
		policy := policies[i].(map[string]interface{})
		result := &vpc.SecurityGroupPolicy{
			Action: helper.String(policy["action"].(string)),
		}
		// CidrBlock, Ipv6CidrBlock, SecurityGroupId, AddressTemplate are exclusive, and Protocol + Port, ServiceTemplate are also exclusive
		var (
			cidrBlock            = policy["cidr_block"].(string)
			ipv6CidrBlock        = policy["ipv6_cidr_block"].(string)
			sgId                 = policy["source_security_id"].(string)
			addressTemplateId    = policy["address_template_id"].(string)
			addressTemplateGroup = policy["address_template_group"].(string)
			protocol             = policy["protocol"].(string)
			port                 = policy["port"].(string)
			serviceTemplate      = policy["service_template_id"].(string)
			serviceTemplateGroup = policy["service_template_group"].(string)
			desc                 = policy["description"].(string)
		)

		// check if exclusive arguments both set
		checkExcludeValues := func(item map[string]string) (result []string) {
			for k, v := range item {
				if v != "" {
					result = append(result, k)
				}
			}
			return result
		}

		if excludes := checkExcludeValues(map[string]string{
			"cidr_block":             cidrBlock,
			"ipv6_cidr_block":        ipv6CidrBlock,
			"source_security_id":     sgId,
			"address_template_id":    addressTemplateId,
			"address_template_group": addressTemplateGroup,
		}); len(excludes) > 1 {
			err = fmt.Errorf("conflict at rule.%d, cannot set %s in time", i, strings.Join(excludes, ","))
			return
		}

		if excludes := checkExcludeValues(map[string]string{
			"protocol + port":        protocol + port,
			"service_template_id":    serviceTemplate,
			"service_template_group": serviceTemplateGroup,
		}); len(excludes) > 1 {
			err = fmt.Errorf("conflict at rule.%d, cannot set %s in time", i, strings.Join(excludes, ","))
			return
		}

		if cidrBlock != "" {
			result.CidrBlock = &cidrBlock
		}
		if ipv6CidrBlock != "" {
			result.Ipv6CidrBlock = &ipv6CidrBlock
		}
		if sgId != "" {
			result.SecurityGroupId = &sgId
		}
		if addressTemplateId != "" || addressTemplateGroup != "" {
			result.AddressTemplate = &vpc.AddressTemplateSpecification{}
		}
		if addressTemplateId != "" {
			result.AddressTemplate.AddressId = &addressTemplateId
		}
		if addressTemplateGroup != "" {
			result.AddressTemplate.AddressGroupId = &addressTemplateGroup
		}
		if protocol != "" {
			result.Protocol = &protocol
		}
		if port != "" {
			result.Port = &port
		}
		if serviceTemplate != "" || serviceTemplateGroup != "" {
			result.ServiceTemplate = &vpc.ServiceTemplateSpecification{}
		}
		if serviceTemplate != "" {
			result.ServiceTemplate.ServiceId = &serviceTemplate
		}
		if serviceTemplateGroup != "" {
			result.ServiceTemplate.ServiceGroupId = &serviceTemplateGroup
		}
		if desc != "" {
			result.PolicyDescription = &desc
		}

		output = append(output, result)
	}
	return
}

func marshalSecurityPolicy(policies []*vpc.SecurityGroupPolicy) []interface{} {
	result := make([]interface{}, 0, len(policies))
	for i := range policies {
		policy := policies[i]
		dMap := map[string]interface{}{
			"action": policy.Action,
		}
		if policy.PolicyIndex != nil {
			dMap["policy_index"] = policy.PolicyIndex
		}
		if policy.CidrBlock != nil {
			dMap["cidr_block"] = policy.CidrBlock
		}
		if policy.Ipv6CidrBlock != nil {
			dMap["ipv6_cidr_block"] = policy.Ipv6CidrBlock
		}
		if policy.Ipv6CidrBlock != nil {
			dMap["source_security_id"] = policy.SecurityGroupId
		}
		if policy.AddressTemplate != nil && policy.AddressTemplate.AddressId != nil {
			dMap["address_template_id"] = policy.AddressTemplate.AddressId
		}
		if policy.AddressTemplate != nil && policy.AddressTemplate.AddressGroupId != nil {
			dMap["address_template_group"] = policy.AddressTemplate.AddressGroupId
		}
		if policy.Protocol != nil /*!checkPolicyPortIgnore(policy.Protocol, policy)*/ {
			dMap["protocol"] = strings.ToUpper(*policy.Protocol)
		}
		if policy.Port != nil /*!checkPolicyPortIgnore(policy.Port, policy)*/ {
			dMap["port"] = policy.Port
		}
		if policy.ServiceTemplate != nil && policy.ServiceTemplate.ServiceId != nil {
			dMap["service_template_id"] = policy.ServiceTemplate.ServiceId
		}
		if policy.ServiceTemplate != nil && policy.ServiceTemplate.ServiceGroupId != nil {
			dMap["service_template_group"] = policy.ServiceTemplate.ServiceGroupId
		}
		if policy.PolicyDescription != nil {
			dMap["description"] = policy.PolicyDescription
		}
		result = append(result, dMap)
	}
	return result
}
