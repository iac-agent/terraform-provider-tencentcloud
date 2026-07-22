package dayu

import (
	"context"
	"errors"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	//sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"
)

func ResourceTencentCloudDayuCCHttpPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuCCHttpPolicyCreate,
		Read:   resourceTencentCloudDayuCCHttpPolicyRead,
		Update: resourceTencentCloudDayuCCHttpPolicyUpdate,
		Delete: resourceTencentCloudDayuCCHttpPolicyDelete,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID 资源 该 CC self-define http 策略 works 对于。",
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_RESOURCE_TYPE),
				ForceNew:     true,
				Description:  "类型 资源 该 CC self-define http 策略 works 对于，有效 值 是 `bgpip`，`bgp`，`bgp-multip` 和 `net`。",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 20),
				Description:  "名称 CC self-define http 策略. Length should between 1 和 20。",
			},
			"smode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_CC_POLICY_SMODE),
				Default:      DAYU_CC_POLICY_SMODE_MATCH,
				Description:  "Match 模式，和 有效 值 是 `matching`，`speedlimit`. 注意: speed 限制 类型 CC self-define 策略 可以 仅 集合 一个。",
			},
			"frequency": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 10000),
				Description:  "Max 频率 per minute，仅 有效 当 `smode` 是 `speedlimit`， 有效 值 ranges 从 1 到 10000。",
			},
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_CC_POLICY_ACTION),
				Description:  "操作 模式，仅 有效 当 `smode` 是 `matching`. 有效 值 是 `alg` 和 `drop`。",
			},
			"switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicate CC self-define http 策略 takes effect 或 不。",
			},
			"rule_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"skey": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "host",
							ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_CC_POLICY_HTTP_CHECK_TYPE),
							Description:  "键 的 规则. 有效值：`主机`，`cgi`，`ua`，`referer`。",
						},
						"operator": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "include",
							ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_CC_POLICY_CHECK_OP),
							Description:  "操作者 的 规则. 有效值：`include`，`not_include`，`equal`。",
						},
						"value": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "",
							ValidateFunc: tccommon.ValidateStringLengthInRange(0, 31),
							Description:  "Rule 值，then 长度 should 是 less 比 31 bytes。",
						},
					},
				},
				Description: "Rule 列表 CC self-define http 策略， 仅 有效 当 `smode` 是 `matching`。",
			},
			"ip": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: tccommon.ValidateIp,
				Description:  "Ip 的 CC self-define http 策略，仅 有效 当 `resource_type` 是 `bgp-multip`. num 的 列表 items 可以 仅 是 集合 一个。",
			},
			//computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "创建时间 的 CC self-define http 策略。",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID 的 CC self-define http 策略。",
			},
		},
	}
}

func resourceTencentCloudDayuCCHttpPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_cc_http_policy.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	resourceId := d.Get("resource_id").(string)
	resourceType := d.Get("resource_type").(string)
	//set CCPolicy
	var ccPolicy dayu.CCPolicy
	ccPolicy.Name = helper.String(d.Get("name").(string))
	smode := d.Get("smode").(string)
	ccPolicy.Smode = &smode
	frequency := 0
	if v, ok := d.GetOk("frequency"); ok {
		frequency = v.(int)
	}

	if smode == DAYU_CC_POLICY_SMODE_SPEED_LIMIT {
		if frequency == 0 {
			return fmt.Errorf("`frequencys` should be set when `smode` is `speedlimit`.")
		}
		ccPolicy.Frequency = helper.IntUint64(frequency)
	} else {
		ccPolicy.ExeMode = helper.String(d.Get("action").(string))
	}
	ccPolicy.Protocol = helper.String(DAYU_L7_RULE_PROTOCOL_HTTP)
	switchFlag := d.Get("switch").(bool)
	if switchFlag {
		ccPolicy.Switch = helper.IntUint64(1)
	} else {
		ccPolicy.Switch = helper.IntUint64(0)
	}

	ip := ""
	if v, ok := d.GetOk("ip"); ok {
		ip = v.(string)
	}
	if ip != "" {
		ccPolicy.IpList = []*string{&ip}
	} else {
		ccPolicy.IpList = []*string{}
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

func resourceTencentCloudDayuCCHttpPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_cc_http_policy.read")()
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
	_ = d.Set("smode", policy.Smode)
	_ = d.Set("policy_id", policy.SetId)
	_ = d.Set("action", policy.ExeMode)
	ipList := helper.StringsInterfaces(policy.IpList)
	if len(ipList) == 1 {
		_ = d.Set("ip", ipList[0])
	}
	_ = d.Set("switch", *policy.Switch > 0)

	if policy.Frequency != nil && *policy.Smode == "frequency" {
		_ = d.Set("frequency", policy.Frequency)
	}
	if policy.RuleList != nil && *policy.Smode == "matching" {
		_ = d.Set("rule_list", flattenCCRuleList(policy.RuleList))
	}

	return nil
}

func resourceTencentCloudDayuCCHttpPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_cc_http_policy.update")()

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
	smode := d.Get("smode").(string)
	ccPolicy.Smode = &smode
	frequency := 0
	if v, ok := d.GetOk("frequency"); ok {
		frequency = v.(int)
	}
	if smode == DAYU_CC_POLICY_SMODE_SPEED_LIMIT {
		if frequency == 0 {
			return fmt.Errorf("`speedlimit` should be set when `smode` is `speedlimit`.")
		}
		ccPolicy.Frequency = helper.IntUint64(frequency)
	} else {
		ccPolicy.ExeMode = helper.String(d.Get("action").(string))
	}
	ccPolicy.Protocol = helper.String(DAYU_L7_RULE_PROTOCOL_HTTP)

	switchFlag := d.Get("switch").(bool)
	if switchFlag {
		ccPolicy.Switch = helper.IntUint64(1)
	} else {
		ccPolicy.Switch = helper.IntUint64(0)
	}

	ruleList := d.Get("rule_list").([]interface{})
	ccPolicy.RuleList = make([]*dayu.CCRule, 0, len(ruleList))
	for _, rule := range ruleList {
		var ccRule dayu.CCRule
		ruleMap := rule.(map[string]interface{})
		ccRule.Skey = helper.String(ruleMap["skey"].(string))
		ccRule.Operator = helper.String(ruleMap["operator"].(string))
		ccRule.Value = helper.String(ruleMap["value"].(string))
		ccPolicy.RuleList = append(ccPolicy.RuleList, &ccRule)
	}
	ip := ""
	if v, ok := d.GetOk("ip"); ok {
		ip = v.(string)
	}
	if ip != "" {
		ccPolicy.IpList = []*string{&ip}
	} else {
		ccPolicy.IpList = []*string{}
	}
	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := dayuService.ModifyCCSelfdefinePolicy(ctx, resourceType, resourceId, policyId, ccPolicy)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudDayuCCHttpPolicyRead(d, meta)
}

func resourceTencentCloudDayuCCHttpPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_cc_http_policy.delete")()

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
