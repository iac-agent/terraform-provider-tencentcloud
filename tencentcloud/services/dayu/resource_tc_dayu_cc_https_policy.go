package dayu

import (
	"context"
	"errors"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDayuCCHttpsPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuCCHttpsPolicyCreate,
		Read:   resourceTencentCloudDayuCCHttpsPolicyRead,
		Update: resourceTencentCloudDayuCCHttpsPolicyUpdate,
		Delete: resourceTencentCloudDayuCCHttpsPolicyDelete,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID 资源 该 CC self-define https 策略 works 对于。",
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_RESOURCE_TYPE_HTTPS),
				ForceNew:     true,
				Description:  "类型 资源 该 CC self-define https 策略 works 对于，有效 值 是 `bgpip`。",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 20),
				Description:  "名称 CC self-define https 策略. Length should between 1 和 20。",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "域名 该 CC self-define https 策略 works 对于，仅 有效 当 `协议` 是 `https`。",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Rule ID 域名 该 CC self-define https 策略 works 对于，仅 有效 当 `协议` 是 `https`。",
			},
			"switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicate CC self-define https 策略 takes effect 或 不。",
			},
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_CC_POLICY_ACTION),
				Description:  "操作 模式 有效 值 是 `alg` 和 `drop`。",
			},
			"rule_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"skey": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_CC_POLICY_HTTPS_CHECK_TYPE),
							Description:  "键 的 规则. 有效 值 是 `cgi`，`ua` 和 `referer`。",
						},
						"operator": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_CC_POLICY_CHECK_OP_HTTPS),
							Description:  "操作者 的 规则. 有效 值 是 `include` 和 `equal`。",
						},
						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateStringLengthInRange(0, 31),
							Description:  "Rule 值，then 长度 should 是 less 比 31 bytes。",
						},
					},
				},
				Description: "Rule 列表 CC self-define https 策略。",
			},
			//computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "创建时间 的 CC self-define https 策略。",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID 的 CC self-define https 策略。",
			},
			"ip_list": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Ip 的 CC self-define https 策略。",
			},
		},
	}
}

func resourceTencentCloudDayuCCHttpsPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_cc_https_policy.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	resourceId := d.Get("resource_id").(string)
	resourceType := d.Get("resource_type").(string)
	//set CCPolicy
	var ccPolicy dayu.CCPolicy
	ccPolicy.Name = helper.String(d.Get("name").(string))
	ccPolicy.Smode = helper.String(DAYU_CC_POLICY_SMODE_MATCH)
	ccPolicy.Protocol = helper.String(DAYU_L7_RULE_PROTOCOL_HTTPS)
	ccPolicy.Domain = helper.String(d.Get("domain").(string))
	ccPolicy.RuleId = helper.String(d.Get("rule_id").(string))
	ccPolicy.ExeMode = helper.String(d.Get("action").(string))

	ccPolicy.IpList = []*string{}

	if v, ok := d.GetOk("rule_id"); ok {
		ccPolicy.RuleId = helper.String(v.(string))
	}

	switchFlag := d.Get("switch").(bool)
	if switchFlag {
		ccPolicy.Switch = helper.IntUint64(1)
	} else {
		ccPolicy.Switch = helper.IntUint64(0)
	}

	ruleList := d.Get("rule_list").(*schema.Set).List()
	ccPolicy.RuleList = make([]*dayu.CCRule, 0, len(ruleList))
	for _, rule := range ruleList {
		var ccRule dayu.CCRule
		ruleMap := rule.(map[string]interface{})
		ccRule.Skey = helper.String(ruleMap["skey"].(string))
		ccRule.Operator = helper.String(ruleMap["operator"].(string))
		ccRule.Value = helper.String(ruleMap["value"].(string))
		ccPolicy.RuleList = append(ccPolicy.RuleList, &ccRule)
	}

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	policyId := ""

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := dayuService.CreateCCSelfdefinePolicy(ctx, resourceType, resourceId, ccPolicy)
		if e != nil {
			return tccommon.RetryError(e)
		}
		policyId = result
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(resourceType + tccommon.FILED_SP + resourceId + tccommon.FILED_SP + policyId)

	return resourceTencentCloudDayuCCHttpPolicyRead(d, meta)
}

func resourceTencentCloudDayuCCHttpsPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_cc_https_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of dayu CC policy")
	}
	resourceType := items[0]
	resourceId := items[1]
	policyId := items[2]

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	policy, has, err := dayuService.DescribeCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			policy, has, err = dayuService.DescribeCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", policy.Name)
	_ = d.Set("create_time", policy.CreateTime)
	_ = d.Set("policy_id", policy.SetId)
	_ = d.Set("action", policy.ExeMode)
	_ = d.Set("rule_id", policy.RuleId)
	_ = d.Set("domain", policy.Domain)
	_ = d.Set("switch", *policy.Switch > 0)
	_ = d.Set("rule_list", flattenCCRuleList(policy.RuleList))
	_ = d.Set("ip_list", helper.StringsInterfaces(policy.IpList))

	return nil
}

func resourceTencentCloudDayuCCHttpsPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_cc_https_policy.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of CC policy")
	}
	resourceType := items[0]
	resourceId := items[1]
	policyId := items[2]

	//set CCPolicy
	var ccPolicy dayu.CCPolicy
	ccPolicy.Name = helper.String(d.Get("name").(string))
	ccPolicy.Smode = helper.String(DAYU_CC_POLICY_SMODE_MATCH)
	ccPolicy.Protocol = helper.String(DAYU_L7_RULE_PROTOCOL_HTTP)
	ccPolicy.SetId = helper.String(policyId)

	if v, ok := d.GetOk("rule_id"); ok {
		ccPolicy.RuleId = helper.String(v.(string))
	}

	switchFlag := d.Get("switch").(bool)
	if switchFlag {
		ccPolicy.Switch = helper.IntUint64(1)
	} else {
		ccPolicy.Switch = helper.IntUint64(0)
	}

	ccPolicy.ExeMode = helper.String(d.Get("action").(string))
	ruleList := d.Get("rule_list").(*schema.Set).List()
	ccPolicy.RuleList = make([]*dayu.CCRule, 0, len(ruleList))
	for _, rule := range ruleList {
		var ccRule dayu.CCRule
		ruleMap := rule.(map[string]interface{})
		ccRule.Skey = helper.String(ruleMap["skey"].(string))
		ccRule.Operator = helper.String(ruleMap["operator"].(string))
		ccRule.Value = helper.String(ruleMap["value"].(string))
		ccPolicy.RuleList = append(ccPolicy.RuleList, &ccRule)
	}

	//the sdk really designed error, it need this para
	ipList := d.Get("ip_list").(*schema.Set).List()
	for _, ip := range ipList {
		ccPolicy.IpList = append(ccPolicy.IpList, helper.String(ip.(string)))
	}
	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := dayuService.ModifyCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId, ccPolicy)
		if e != nil {
			return tccommon.RetryError(e, tccommon.InternalError)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudDayuCCHttpPolicyRead(d, meta)
}

func resourceTencentCloudDayuCCHttpsPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_cc_https_policy.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of CC policy")
	}
	resourceType := items[0]
	resourceId := items[1]
	policyId := items[2]

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := dayuService.DeleteCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	_, has, err := dayuService.DescribeCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
	if err != nil || has {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, has, err = dayuService.DescribeCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId)
			if err != nil {
				return tccommon.RetryError(err)
			}

			if has {
				err = fmt.Errorf("delete DDoS policy fail, CC policy still exist from sdk DescribeCCSelfDefinePolicy")
				return resource.RetryableError(err)
			}

			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		return nil
	} else {
		return errors.New("delete CC policy fail, CC policy still exist from sdk DescribeCCSelfDefinePolicy")
	}
}
