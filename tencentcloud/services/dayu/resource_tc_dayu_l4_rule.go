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

func ResourceTencentCloudDayuL4Rule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuL4RuleCreate,
		Read:   resourceTencentCloudDayuL4RuleRead,
		Update: resourceTencentCloudDayuL4RuleUpdate,
		Delete: resourceTencentCloudDayuL4RuleDelete,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID 资源 该 layer 4 规则 works 对于。",
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_RESOURCE_TYPE_RULE),
				ForceNew:     true,
				Description:  "类型 资源 该 layer 4 规则 works 对于. 有效值：`bgpip` 和 `net`。",
			},
			"source_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue(DAYU_L7_RULE_SOURCE_TYPE),
				Description:  "来源 类型，`1` 对于 来源 的 主机，`2` 对于 来源 的 IP。",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "名称 规则. 当 `resource_type` 是 `net`，此 字段 should 是 集合 使用 有效 域名",
			},
			"s_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidatePort,
				Description:  "来源 端口 的 L4 规则。",
			},
			"d_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidatePort,
				Description:  "destination 端口 的 L4 规则。",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_L4_RULE_PROTOCOL),
				Description:  "协议 的 规则. 有效值：`http`，`https`. 当 `source_type` 是 1(主机 来源)， 值 的 此 字段 可以 仅 集合 使用 `tcp`。",
			},
			"source_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "来源 IP 或 域名，有效 格式 的 ip 是 like `1.1.1.1` 和 有效 格式 的 主机 来源 是 like `abc.com`。",
						},
						"weight": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 100),
							Description:  "权重 的 来源， 有效 值 ranges 从 0 到 100。",
						},
					},
				},
				MinItems:    1,
				MaxItems:    20,
				Description: "来源 列表 规则，它 可以 是 集合 的 ip sources 或 集合 的 域名 sources. 数量 items ranges 从 1 到 20。",
			},
			"health_check_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "表示是否health check 是 已启用 默认为 `false`. Only 有效 当 来源 列表 has more 比 一个 来源 item。",
			},
			"health_check_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(10, 60),
				Description:  "Interval 时间 的 health check. 值 范围 是 10-60 sec，和 默认为 15 sec。",
			},
			"health_check_health_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 10),
				Description:  "Health 阈值 的 health check，和 默认为 3. 如果 success 结果 是 返回 对于 health check 3 consecutive times，表示that forwarding 是 normal. 值 范围 是 2-10。",
			},
			"health_check_unhealth_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 10),
				Description:  "Unhealthy 阈值 的 health check，和 默认为 3. 如果 unhealthy 结果 是 返回 3 consecutive times，表示that forwarding 是 abnormal. 值 范围 是 2-10。",
			},
			"health_check_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 60),
				Description:  "HTTP 状态 代码 默认为 26 和 值 范围 是 2-60。",
			},
			"session_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate 该 会话 将 keep 或 不，和 默认值为 `false`。",
			},
			"session_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 300),
				Description:  "Session keep 时间，仅 有效 当 `session_switch` 是 true， 可用 值 ranges 从 1 到 300 和 单位 是 second。",
			},
			//computed
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID layer 4 规则。",
			},
			"lb_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "LB 类型 规则. 有效值：`1`，`2`. `1` 对于 权重 cycling 和 `2` 对于 IP hash。",
			},
		},
	}
}

func resourceTencentCloudDayuL4RuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_l4_rule.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	resourceId := d.Get("resource_id").(string)
	resourceType := d.Get("resource_type").(string)

	destPort := d.Get("d_port").(int)

	//check
	protocol := d.Get("protocol").(string)
	source_type := d.Get("source_type").(int)
	sourceList := d.Get("source_list").(*schema.Set).List()

	if source_type == DAYU_L7_RULE_SOURCE_TYPE_HOST && protocol != DAYU_L4_RULE_PROTOCOL_TCP {
		return fmt.Errorf("`protocol` can only be set with `TCP` when `source_type` is 1(host source).")
	}

	healthCheckSwitch := d.Get("health_check_switch").(bool)
	if healthCheckSwitch {
		if len(sourceList) <= 1 {
			return fmt.Errorf("The `health_check_switch` cannot be set `true` when `source_list` has only one item.")
		}
	}

	//check
	timeout := 0
	interval := 0
	if v, ok := d.GetOk("health_check_timeout"); ok {
		timeout = v.(int)
	}
	if v, ok := d.GetOk("health_check_interval"); ok {
		interval = v.(int)
	}

	if timeout > 0 && interval > 0 {
		if timeout > interval {
			return fmt.Errorf("The `health_check_interval` should be greater than `health_check_timeout`.")
		}
	}

	//set L4RuleEntry
	var rule dayu.L4RuleEntry
	rule.LbType = helper.IntUint64(1)
	rule.SourcePort = helper.IntUint64(d.Get("s_port").(int))
	rule.VirtualPort = helper.IntUint64(destPort)
	rule.Protocol = helper.String(d.Get("protocol").(string))
	rule.RuleName = helper.String(d.Get("name").(string))
	sourceType := d.Get("source_type").(int)
	//check that there is no check with the source list and sdk returns
	rule.SourceType = helper.IntUint64(sourceType)

	rule.SourceList = make([]*dayu.L4RuleSource, 0, len(sourceList))
	for _, source := range sourceList {
		sourceMap := source.(map[string]interface{})
		var l4RuleSource dayu.L4RuleSource
		l4RuleSource.Source = helper.String(sourceMap["source"].(string))
		l4RuleSource.Weight = helper.IntUint64(sourceMap["weight"].(int))
		rule.SourceList = append(rule.SourceList, &l4RuleSource)
	}

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	ruleId := ""
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := dayuService.CreateL4Rule(ctx, resourceType, resourceId, rule)
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

	//set health check
	var healthCheck dayu.L4HealthConfig
	healthCheck.Protocol = helper.String(d.Get("protocol").(string))
	healthCheck.Enable = helper.BoolToInt64Pointer(d.Get("health_check_switch").(bool))
	healthCheck.Interval = helper.IntUint64(d.Get("health_check_interval").(int))
	healthCheck.KickNum = helper.IntUint64(d.Get("health_check_unhealth_num").(int))
	healthCheck.AliveNum = helper.IntUint64(d.Get("health_check_health_num").(int))
	healthCheck.TimeOut = helper.IntUint64(d.Get("health_check_timeout").(int))
	healthCheck.VirtualPort = helper.IntUint64(destPort)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := dayuService.SetL4Health(ctx, resourceType, resourceId, healthCheck)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	//set session
	sessionFlag := d.Get("session_switch").(bool)
	sessionTime := 0
	if v, ok := d.GetOk("session_time"); ok {
		sessionTime = v.(int)
	}
	if sessionTime == 0 && sessionFlag {
		return fmt.Errorf("`session_time` should be set when `session_switch` is true.")
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := dayuService.SetSession(ctx, resourceType, resourceId, ruleId, sessionFlag, sessionTime)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudDayuL4RuleRead(d, meta)
}

func resourceTencentCloudDayuL4RuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_l4_rule.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of dayu L4 rule")
	}
	resourceType := items[0]
	resourceId := items[1]
	ruleId := items[2]

	d.Partial(true)

	sourceType := d.Get("source_type").(int)
	protocol := d.Get("protocol").(string)
	//check
	if sourceType == 1 && protocol != DAYU_L4_RULE_PROTOCOL_TCP {
		return fmt.Errorf("`protocol` can only be set with `TCP` when `source_type` is 1(host source).")
	}
	sourceList := d.Get("source_list").(*schema.Set).List()

	//check
	timeout := 0
	interval := 0
	if v, ok := d.GetOk("health_check_timeout"); ok {
		timeout = v.(int)
	}
	if v, ok := d.GetOk("health_check_interval"); ok {
		interval = v.(int)
	}

	if timeout > 0 && interval > 0 {
		if timeout > interval {
			return fmt.Errorf("The `health_check_interval` should be greater than `health_check_timeout`.")
		}
	}

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	ruleFlag := false
	ruleKey := []string{"s_port", "d_port", "protocol", "source_list"}

	for _, key := range ruleKey {
		if d.HasChange(key) {
			ruleFlag = true
		}
	}

	if ruleFlag {
		//set L4RuleEntry
		var rule dayu.L4RuleEntry
		rule.LbType = helper.IntUint64(1)
		rule.SourcePort = helper.IntUint64(d.Get("s_port").(int))
		rule.VirtualPort = helper.IntUint64(d.Get("d_port").(int))
		rule.Protocol = helper.String(d.Get("protocol").(string))
		rule.RuleName = helper.String(d.Get("name").(string))
		rule.RuleId = &ruleId

		rule.SourceType = helper.IntUint64(sourceType)
		rule.Protocol = &protocol

		rule.SourceList = make([]*dayu.L4RuleSource, 0, len(sourceList))
		for _, source := range sourceList {
			sourceMap := source.(map[string]interface{})
			var l4RuleSource dayu.L4RuleSource
			l4RuleSource.Source = helper.String(sourceMap["source"].(string))
			l4RuleSource.Weight = helper.IntUint64(sourceMap["weight"].(int))
			rule.SourceList = append(rule.SourceList, &l4RuleSource)
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := dayuService.ModifyL4Rule(ctx, resourceType, resourceId, rule)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	healthFlag := false
	healthKey := []string{"health_check_switch", "health_check_interval", "health_check_timeout", "health_check_unhealth_num", "health_check_health_num", "d_port"}

	for _, key := range healthKey {
		if d.HasChange(key) {
			healthFlag = true
		}
	}

	if healthFlag {
		//set health check
		if len(sourceList) <= 1 {
			return fmt.Errorf("The `health_check_switch` cannot be set when `source_list` has only one item.")
		}

		var healthCheck dayu.L4HealthConfig
		healthCheck.Protocol = helper.String(d.Get("protocol").(string))
		healthCheck.Enable = helper.BoolToInt64Pointer(d.Get("health_check_switch").(bool))
		healthCheck.Interval = helper.IntUint64(d.Get("health_check_interval").(int))
		healthCheck.KickNum = helper.IntUint64(d.Get("health_check_unhealth_num").(int))
		healthCheck.AliveNum = helper.IntUint64(d.Get("health_check_health_num").(int))
		healthCheck.TimeOut = helper.IntUint64(d.Get("health_check_timeout").(int))
		healthCheck.VirtualPort = helper.IntUint64(d.Get("d_port").(int))

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := dayuService.SetL4Health(ctx, resourceType, resourceId, healthCheck)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})

		if err != nil {
			return err
		}

	}

	if d.HasChange("session_switch") || d.HasChange("session_time") {
		sessionFlag := d.Get("session_switch").(bool)
		sessionTime := 0
		if v, ok := d.GetOk("session_time"); ok {
			sessionTime = v.(int)
		}
		if sessionTime == 0 && sessionFlag {
			return fmt.Errorf("`session_time` should be set when `session_switch` is true.")
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := dayuService.SetSession(ctx, resourceType, resourceId, ruleId, sessionFlag, sessionTime)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})

		if err != nil {
			return err
		}

	}

	d.Partial(false)

	return resourceTencentCloudDayuL4RuleRead(d, meta)
}

func resourceTencentCloudDayuL4RuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_l4_rule.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of dayu L4 rule")
	}
	resourceType := items[0]
	resourceId := items[1]
	ruleId := items[2]

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	//set rule
	rule, health, has, err := dayuService.DescribeL4Rule(ctx, resourceType, resourceId, ruleId)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			rule, health, has, err = dayuService.DescribeL4Rule(ctx, resourceType, resourceId, ruleId)
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
	_ = d.Set("s_port", int(*rule.SourcePort))
	_ = d.Set("d_port", int(*rule.VirtualPort))
	_ = d.Set("rule_id", rule.RuleId)
	_ = d.Set("lb_type", int(*rule.LbType))
	_ = d.Set("name", rule.RuleName)
	_ = d.Set("source_type", int(*rule.SourceType))
	_ = d.Set("session_time", int(*rule.KeepTime))
	_ = d.Set("session_switch", *rule.KeepEnable > 0)
	_ = d.Set("source_list", flattenSourceList(rule.SourceList))

	//set health check
	if health == nil {
		_ = d.Set("health_check_switch", false)
		return nil
	}
	_ = d.Set("health_check_switch", *health.Enable > 0)
	_ = d.Set("health_check_timeout", int(*health.TimeOut))
	_ = d.Set("health_check_health_num", int(*health.AliveNum))
	_ = d.Set("health_check_unhealth_num", int(*health.KickNum))
	_ = d.Set("health_check_interval", int(*health.Interval))

	return nil
}

func resourceTencentCloudDayuL4RuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_l4_rule.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of L4 rule")
	}
	resourceType := items[0]
	resourceId := items[1]
	ruleId := items[2]

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := dayuService.DeleteL4Rule(ctx, resourceType, resourceId, ruleId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	_, _, has, err := dayuService.DescribeL4Rule(ctx, resourceType, resourceId, ruleId)
	if err != nil || has {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, _, has, err = dayuService.DescribeL4Rule(ctx, resourceType, resourceId, ruleId)
			if err != nil {
				return tccommon.RetryError(err)
			}

			if has {
				err = fmt.Errorf("delete L4 rule fail, L4 rule %s still exist from sdk", ruleId)
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
		return errors.New("delete CC policy fail, CC policy still exist from sdk")
	}
}
