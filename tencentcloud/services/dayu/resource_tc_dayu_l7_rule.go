package dayu

import (
	"context"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDayuL7Rule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuL7RuleCreate,
		Read:   resourceTencentCloudDayuL7RuleRead,
		Update: resourceTencentCloudDayuL7RuleUpdate,
		Delete: resourceTencentCloudDayuL7RuleDelete,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID 资源 该 layer 7 规则 works 对于。",
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_RESOURCE_TYPE),
				ForceNew:     true,
				Description:  "类型 资源 该 layer 7 规则 works 对于，有效 值 是 `bgpip`。",
			},
			"domain": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(0, 80),
				Description:  "域名 该 layer 7 规则 works 对于. 有效 字符串 长度 ranges 从 0 到 80。",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_L7_RULE_PROTOCOL),
				Description:  "协议 的 规则. 有效值：`http`，`https`。",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "名称 规则。",
			},
			"switch": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicate 规则 将 take effect 或 不。",
			},
			"source_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue(DAYU_L7_RULE_SOURCE_TYPE),
				Description:  "来源 类型，`1` 对于 来源 的 主机，`2` 对于 来源 的 IP。",
			},
			"source_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "来源 IP 或 域名，有效 格式 的 ip 是 like `1.1.1.1:60` 或 `1.1.1.1` 和 有效 格式 的 主机 来源 是 like `abc.com` 或 `abc.com:80`。",
				},
				MinItems:    1,
				MaxItems:    16,
				Description: "来源 列表 规则，它 可以 是 集合 的 ip sources 或 集合 的 域名 sources. 数量 items ranges 从 1 到 16。",
			},
			"ssl_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SSL ID，当 `协议` 是 `https`， 字段 should 是 集合 使用 有效 SSL ID。",
			},
			"health_check_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "表示是否health check 是 已启用 默认为 `false`。",
			},
			"health_check_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(10, 60),
				Description:  "Interval 时间 的 health check. 有效 值 ranges: [10~60]sec. 默认为 15 sec。",
			},
			"health_check_health_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 10),
				Description:  "Health 阈值 的 health check，和 默认为 `3`. 如果 success 结果 是 返回 对于 health check 3 consecutive times，表示that forwarding 是 normal. 值 范围 是 [2-10]。",
			},
			"health_check_unhealth_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 10),
				Description:  "Unhealthy 阈值 的 health check，和 默认为 `3`. 如果 unhealthy 结果 是 返回 3 consecutive times，表示that forwarding 是 abnormal. 值 范围 是 [2-10]。",
			},
			"health_check_code": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 31),
				Description:  "HTTP 状态 代码 默认为 `26`. 有效 值 ranges: [1~31]. `1` 表示 返回值 '1xx' 是 health. `2` 表示 返回值 '2xx' 是 health. `4` 表示 返回值 '3xx' 是 health. `8` 表示 返回值 '4xx' 是 health. `16` 表示 返回值 '5xx' 是 health. 如果 您 want 多个 返回 codes 到 indicate health，need 到 add corresponding 值。",
			},
			"health_check_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "路径 的 health check. 默认为 `/`。",
			},
			"health_check_method": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_RULE_METHOD),
				Description:  "Methods 的 health check. 默认为 'HEAD'， 可用 值 是 'HEAD' 和 'GET'。",
			},
			//computed
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID layer 7 规则。",
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "状态 规则. `0` 对于 create/modify success，`2` 对于 create/modify fail，`3` 对于 delete success，`5` 对于 delete failed，`6` 对于 waiting 到 是 创建/modified，`7` 对于 waiting 到 是 删除 和 8 对于 waiting 到 get SSL ID。",
			},
		},
	}
}

func resourceTencentCloudDayuL7RuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_l7_rule.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	resourceId := d.Get("resource_id").(string)
	resourceType := d.Get("resource_type").(string)

	domain := d.Get("domain").(string)

	//set L4RuleEntry
	var rule dayu.L7RuleEntry
	rule.LbType = helper.IntUint64(1)
	//test that the keep time para will make effect
	protocol := d.Get("protocol").(string)
	sslId := ""
	if protocol == "https" {
		if v, ok := d.GetOk("ssl_id"); ok {
			sslId = v.(string)
		}
		if sslId == "" {
			return fmt.Errorf("`ssl_id` should be set when `protocol` is `https`.")
		}
		rule.SSLId = &sslId
		rule.CertType = helper.IntUint64(2)
	} else {
		rule.CertType = helper.IntUint64(0)
	}
	rule.Protocol = &protocol
	rule.RuleName = helper.String(d.Get("name").(string))
	sourceType := d.Get("source_type").(int)
	//check that there is no check with the source list and sdk returns
	rule.SourceType = helper.IntUint64(sourceType)
	rule.Domain = &domain

	sourceList := d.Get("source_list").(*schema.Set).List()
	//check
	healthCheckSwitch := d.Get("health_check_switch").(bool)
	if healthCheckSwitch {
		if len(sourceList) <= 1 {
			return fmt.Errorf("The `health_check_switch` cannot be set `true` when `source_list` has only one item.")
		}
	}

	for _, source := range sourceList {
		var l4RuleSource dayu.L4RuleSource
		l4RuleSource.Source = helper.String(source.(string))
		l4RuleSource.Weight = helper.IntUint64(0)
		rule.SourceList = append(rule.SourceList, &l4RuleSource)
	}

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	ruleId := ""
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := dayuService.CreateL7Rule(ctx, resourceType, resourceId, rule)
		if e != nil {
			return tccommon.RetryError(e)
		}
		ruleId = result
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(resourceType + tccommon.FILED_SP + resourceId + tccommon.FILED_SP + ruleId)

	readyFlag, rErr := checkL7RuleStatus(meta, resourceType, resourceId, ruleId, "create")
	if rErr != nil {
		return rErr
	}
	if !readyFlag {
		return fmt.Errorf("Creating operation is timeout %s", ruleId)
	}

	//set health check
	var healthCheck dayu.L7HealthConfig
	healthCheck.Protocol = helper.String(d.Get("protocol").(string))
	healthCheck.Domain = &domain
	healthCheck.Enable = helper.BoolToInt64Pointer(d.Get("health_check_switch").(bool))
	healthCheck.Interval = helper.IntUint64(d.Get("health_check_interval").(int))
	healthCheck.Method = helper.String(d.Get("health_check_method").(string))
	healthCheck.Url = helper.String(d.Get("health_check_path").(string))
	healthCheck.KickNum = helper.IntUint64(d.Get("health_check_unhealth_num").(int))
	healthCheck.AliveNum = helper.IntUint64(d.Get("health_check_health_num").(int))
	healthCheck.StatusCode = helper.IntUint64(d.Get("health_check_code").(int))

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := dayuService.SetL7Health(ctx, resourceType, resourceId, healthCheck)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	readyFlag, rErr = checkL7RuleStatus(meta, resourceType, resourceId, ruleId, "check_health")
	if rErr != nil {
		return rErr
	}
	if !readyFlag {
		return fmt.Errorf("Set health is timeout %s", ruleId)
	}

	//set switch
	switchFlag := d.Get("switch").(bool)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := dayuService.SetRuleSwitch(ctx, resourceType, resourceId, ruleId, switchFlag, protocol)
		if e != nil {
			return tccommon.RetryError(e, tccommon.InternalError)
		}
		return nil
	})

	if err != nil {
		return err
	}

	//check switch status
	readyFlag, rErr = checkL7RuleStatus(meta, resourceType, resourceId, ruleId, fmt.Sprintf("check_switch_%t", switchFlag))
	if rErr != nil {
		return rErr
	}
	if !readyFlag {
		return fmt.Errorf("Set switch is timeout %s", ruleId)
	}

	return resourceTencentCloudDayuL7RuleRead(d, meta)
}

func resourceTencentCloudDayuL7RuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_l7_rule.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of dayu L7 rule")
	}
	resourceType := items[0]
	resourceId := items[1]
	ruleId := items[2]

	domain := d.Get("domain").(string)
	sourceList := d.Get("source_list").(*schema.Set).List()

	d.Partial(true)
	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	ruleFlag := false
	ruleKey := []string{"protocol", "source_type", "source_list", "ssl_id"}

	for _, key := range ruleKey {
		if d.HasChange(key) {
			ruleFlag = true
		}
	}
	if ruleFlag {
		//set L4RuleEntry
		var rule dayu.L7RuleEntry
		rule.LbType = helper.IntUint64(1)
		rule.RuleId = helper.String(ruleId)
		//test that the keep time para will make effect
		protocol := d.Get("protocol").(string)
		sslId := ""
		if protocol == DAYU_L7_RULE_PROTOCOL_HTTPS {
			if v, ok := d.GetOk("ssl_id"); ok {
				sslId = v.(string)
			}
			if sslId == "" {
				return fmt.Errorf("`ssl_id` should be set when `protocol` is `https`.")
			}
			rule.SSLId = &sslId
			rule.CertType = helper.IntUint64(2)
		} else {
			rule.CertType = helper.IntUint64(0)
		}
		rule.RuleName = helper.String(d.Get("name").(string))
		sourceType := d.Get("source_type").(int)
		//check that there is no check with the source list and sdk returns
		rule.SourceType = helper.IntUint64(sourceType)
		rule.Domain = &domain
		rule.Protocol = &protocol

		for _, source := range sourceList {
			var l4RuleSource dayu.L4RuleSource
			l4RuleSource.Source = helper.String(source.(string))
			l4RuleSource.Weight = helper.IntUint64(0)
			rule.SourceList = append(rule.SourceList, &l4RuleSource)
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := dayuService.ModifyL7Rule(ctx, resourceType, resourceId, rule)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})

		if err != nil {
			return err
		}

		readyFlag, rErr := checkL7RuleStatus(meta, resourceType, resourceId, ruleId, "modify")
		if rErr != nil {
			return rErr
		}
		if !readyFlag {
			return fmt.Errorf("Modify operation is timeout %s", ruleId)
		}

	}

	healthFlag := false
	healthKey := []string{"health_check_switch", "health_check_interval", "health_check_path", "health_check_method", "health_check_unhealth_num", "health_check_health_num", "health_check_code"}

	for _, key := range healthKey {
		if d.HasChange(key) {
			healthFlag = true
		}
	}

	if healthFlag {
		//check
		sourceList := d.Get("source_list").(*schema.Set).List()
		if len(sourceList) <= 1 {
			return fmt.Errorf("The `health_check_switch` cannot be set when `source_list` has only one item.")
		}

		//set health check
		var healthCheck dayu.L7HealthConfig
		healthCheck.Protocol = helper.String(d.Get("protocol").(string))
		healthCheck.Domain = &domain
		healthCheck.Enable = helper.BoolToInt64Pointer(d.Get("health_check_switch").(bool))
		healthCheck.Interval = helper.IntUint64(d.Get("health_check_interval").(int))
		healthCheck.Method = helper.String(d.Get("health_check_method").(string))
		healthCheck.Url = helper.String(d.Get("health_check_path").(string))
		healthCheck.KickNum = helper.IntUint64(d.Get("health_check_unhealth_num").(int))
		healthCheck.AliveNum = helper.IntUint64(d.Get("health_check_health_num").(int))
		healthCheck.StatusCode = helper.IntUint64(d.Get("health_check_code").(int))

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := dayuService.SetL7Health(ctx, resourceType, resourceId, healthCheck)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})

		if err != nil {
			return err
		}

		readyFlag, rErr := checkL7RuleStatus(meta, resourceType, resourceId, ruleId, "check_health")
		if rErr != nil {
			return rErr
		}
		if !readyFlag {
			return fmt.Errorf("Set health is timeout %s", ruleId)
		}
	}

	if d.HasChange("switch") {
		//set switch
		switchFlag := d.Get("switch").(bool)
		protocol := d.Get("protocol").(string)
		if d.HasChange("protocol") {
			//set old protocol para close first
			oldInterface, newInterface := d.GetChange("protocol")
			oldProtocol := oldInterface.(string)
			newProtocol := newInterface.(string)
			protocol = oldProtocol
			//open new only
			if switchFlag {
				protocol = newProtocol
			} else {
				protocol = ""
			}
		}
		if protocol != "" {

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				e := dayuService.SetRuleSwitch(ctx, resourceType, resourceId, ruleId, switchFlag, protocol)
				if e != nil {
					return tccommon.RetryError(e, tccommon.InternalError)
				}
				return nil
			})

			if err != nil {
				return err
			}

			//check switch status
			readyFlag, rErr := checkL7RuleStatus(meta, resourceType, resourceId, ruleId, fmt.Sprintf("check_switch_%t", switchFlag))
			if rErr != nil {
				return rErr
			}
			if !readyFlag {
				return fmt.Errorf("Set switch is timeout %s", ruleId)
			}
		}

	}

	d.Partial(false)

	return resourceTencentCloudDayuL7RuleRead(d, meta)
}

func resourceTencentCloudDayuL7RuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_l7_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of dayu L7 rule")
	}
	resourceType := items[0]
	resourceId := items[1]
	ruleId := items[2]

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	//set rule
	rule, health, has, err := dayuService.DescribeL7Rule(ctx, resourceType, resourceId, ruleId)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			rule, health, has, err = dayuService.DescribeL7Rule(ctx, resourceType, resourceId, ruleId)
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
	_ = d.Set("protocol", rule.Protocol)
	_ = d.Set("domain", rule.Domain)
	_ = d.Set("rule_id", rule.RuleId)
	_ = d.Set("ssl_id", rule.SSLId)
	_ = d.Set("name", rule.RuleName)
	_ = d.Set("source_type", int(*rule.SourceType))
	_ = d.Set("status", int(*rule.Status))
	sourceList := make([]*string, 0, len(rule.SourceList))
	for _, v := range rule.SourceList {
		sourceList = append(sourceList, v.Source)
	}
	_ = d.Set("source_list", helper.StringsInterfaces(sourceList))

	if *rule.Protocol == DAYU_L7_RULE_PROTOCOL_HTTPS {
		_ = d.Set("switch", *rule.CCEnable > 0)
	} else {
		_ = d.Set("switch", *rule.CCStatus > 0)
	}
	//set health check
	if health == nil {
		_ = d.Set("health_check_switch", false)
		return nil
	}

	_ = d.Set("health_check_switch", *health.Enable > 0)
	_ = d.Set("health_check_path", health.Url)
	_ = d.Set("health_check_method", health.Method)
	_ = d.Set("health_check_health_num", int(*health.AliveNum))
	_ = d.Set("health_check_unhealth_num", int(*health.KickNum))
	_ = d.Set("health_check_interval", int(*health.Interval))
	_ = d.Set("health_check_code", int(*health.StatusCode))

	return nil
}

func resourceTencentCloudDayuL7RuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_l7_rule.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of L7 rule")
	}
	resourceType := items[0]
	resourceId := items[1]
	ruleId := items[2]

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := dayuService.DeleteL7Rule(ctx, resourceType, resourceId, ruleId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	readyFlag, rErr := checkL7RuleStatus(meta, resourceType, resourceId, ruleId, "delete")
	if rErr != nil {
		return rErr
	}
	if !readyFlag {
		return fmt.Errorf("Delete is timeout %s", ruleId)
	}

	return nil
}

func checkL7RuleStatus(meta interface{}, resourceType string, resourceId string, ruleId string, checkType string) (status bool, errRrt error) {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_l7_rule.check_status")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		sRule, health, has, err := dayuService.DescribeL7Rule(ctx, resourceType, resourceId, ruleId)
		if err != nil {
			return tccommon.RetryError(err)
		}

		if has {
			//created failed
			if *sRule.Status == DAYU_L7_STATUS_SET_FAIL && (checkType == "create" || checkType == "modify") {
				err = fmt.Errorf("%s rule %s failed...", checkType, ruleId)
				status = false
				return resource.NonRetryableError(err)
			} else if *sRule.Status == DAYU_L7_STATUS_SET_DONE && (checkType == "create" || checkType == "modify") {
				//action completed
				status = true
				return nil
			} else if *sRule.Status == DAYU_L7_STATUS_DEL_FAIL && checkType == "delete" {
				//delete failed
				err = fmt.Errorf("%s rule %s failed...", checkType, ruleId)
				status = false
				return resource.NonRetryableError(err)
			} else if health != nil && *health.Status == DAYU_L7_HEALTH_STATUS_SET_DONE && checkType == "check_health" {
				//check health setting completed
				status = true
				return nil
			} else if health != nil && *health.Status == DAYU_L7_HEALTH_STATUS_SET_FAIL && checkType == "check_health" {
				//check health setting failed
				status = false
				err = fmt.Errorf("%s rule %s failed...status %d", checkType, ruleId, *sRule.Status)
				return resource.NonRetryableError(err)
			} else if checkType == "check_switch_true" {
				//check switch set on completed, the para of http is different from https
				if (*sRule.Protocol == DAYU_L7_RULE_PROTOCOL_HTTPS && *sRule.CCEnable == 1) || (*sRule.Protocol == DAYU_L7_RULE_PROTOCOL_HTTP && *sRule.CCStatus == 1) {
					status = true
					return nil
				} else {
					//check switch set on failed
					status = false
					err = fmt.Errorf("%s rule %s ...", checkType, ruleId)
					return resource.RetryableError(err)
				}
			} else if checkType == "check_switch_false" {
				//check switch set off completed, same as above
				if (*sRule.Protocol == DAYU_L7_RULE_PROTOCOL_HTTPS && *sRule.CCEnable == 0) || (*sRule.Protocol == DAYU_L7_RULE_PROTOCOL_HTTP && *sRule.CCStatus == 0) {
					status = true
					return nil
				} else {
					//check switch set off failed
					status = false
					err = fmt.Errorf("%s rule %s ...", checkType, ruleId)
					return resource.RetryableError(err)
				}
			} else {
				if *sRule.Status == DAYU_L7_STATUS_SSL_WAIT {
					//wait to check ssl
					err = fmt.Errorf("SSL is not ready")
				} else {
					//other cases lead to retry(delete done, set waiting, delete waiting, health setting)
					err = fmt.Errorf("%s rule %s wait to be done, Status %d... describe retry", checkType, ruleId, *sRule.Status)
				}
				return resource.RetryableError(err)
			}
		} else {
			if checkType == "delete" {
				//check delete and do not exist, consider success
				status = true
				return nil
			} else {
				//other cases with no exist, report error
				err = fmt.Errorf("cannot find %s rule", ruleId)
				return resource.NonRetryableError(err)
			}
		}
	})

	if err != nil {
		status = false
	}
	return status, err
}
