package monitor

import (
	"fmt"
	"log"
	"strconv"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/go-multierror"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func AlarmPolicyRule() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metric_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "指标名称 或 事件名称",
		},
		"period": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Statistical 周期 （秒）。",
		},
		"operator": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "操作者",
		},
		"value": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Threshold。",
		},
		"continue_period": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "数量 periods。",
		},
		"notice_frequency": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Alarm 间隔 （秒）。",
		},
		"is_power_notice": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "是否alarm 频率 increases exponentially。",
		},
		"filter": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "过滤器 condition 对于 一个 单个 触发器 规则. Must 集合 它 当 create tke-xxx 规则。",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "过滤器 条件类型 有效值：DIMENSION (uses dimensions 对于 filtering)。",
					},
					"dimensions": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "JSON 字符串 generated 通过 serializing AlarmPolicyDimension two-dimensional 数组。",
					},
				},
			},
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Metric display 名称，其中 是 使用 在 output 参数。",
		},
		"unit": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Unit，其中 是 使用 在 output 参数。",
		},
		"rule_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Trigger 条件类型",
		},
	}
}

func ResourceTencentCloudMonitorAlarmPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentMonitorAlarmPolicyCreate,
		Read:   resourceTencentMonitorAlarmPolicyRead,
		Update: resourceTencentMonitorAlarmPolicyUpdate,
		Delete: resourceTencentMonitorAlarmPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "名称 策略。",
			},
			"monitor_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "类型 监控。",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "类型 告警。",
			},
			"remark": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(0, 100),
				Description:  "备注 的 策略 组。",
			},
			"enable": {
				Type:        schema.TypeInt,
				Default:     1,
				Optional:    true,
				Description: "是否enable，默认为 `1`。",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     -1,
				Description: "项目 ID For products 使用 different projects， 值 other 比 -1 必须 是 passed 在。",
			},
			//nolint:misspell
			"conditon_template_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "ID 触发器 condition template。",
			},
			"conditions": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Computed:    true,
				Description: "A 列表 metric 触发器 condition。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_union_rule": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
							Description:  "和 或 relation 的 indicator 告警 规则。",
						},
						"rules": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: AlarmPolicyRule(),
							},
							Description: "A 列表 metric 触发器 condition。",
						},
					},
				},
			},
			"event_conditions": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "A 列表 事件 触发器 condition。",
				Elem: &schema.Resource{
					Schema: AlarmPolicyRule(),
				},
			},
			"notice_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "列表 通知 规则 IDs。",
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "ID 通知 规则 到 是 queried。",
				},
			},
			"trigger_tasks": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Triggered 任务 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Triggered 任务 类型",
						},
						"task_config": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration 信息 在 JSON 格式",
						},
					},
				},
			},
			"policy_tag": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Policy 标签 到 bind 对象。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "标签键",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "标签值",
						},
					},
				},
			},

			"group_by": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Aggregate dimension 列表，指定which dimension keys 到 使用 对于 组 通过。",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"filter": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Global filters。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤器 条件类型 有效值：DIMENSION (uses dimensions 对于 filtering)。",
						},
						"dimensions": {
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "JSON 字符串 generated 通过 serializing AlarmPolicyDimension two-dimensional 数组。",
							DiffSuppressFunc: helper.DiffSupressJSON,
						},
					},
				},
			},

			// compute
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "告警 策略 创建时间。",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "告警 策略 更新时间。",
			},
		},
	}
}

func resourceTencentMonitorAlarmPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_alarm_policy.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	var (
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request        = monitor.NewCreateAlarmPolicyRequest()
	)
	request.Module = helper.String("monitor")
	request.PolicyName = helper.String(d.Get("policy_name").(string))
	request.MonitorType = helper.String(d.Get("monitor_type").(string))
	request.Namespace = helper.String(d.Get("namespace").(string))

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	enable := d.Get("enable").(int)
	request.Enable = helper.IntInt64(enable)

	//if v, ok := d.GetOk("enable"); ok {
	//	request.Enable = helper.IntInt64(v.(int))
	//}

	projectId := d.Get("project_id").(int)
	if projectId != -1 {
		request.ProjectId = helper.IntInt64(projectId)
	}

	//nolint:misspell
	if v, ok := d.GetOk("conditon_template_id"); ok {
		request.ConditionTemplateId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("conditions"); ok {
		conditions := v.([]interface{})

		if len(conditions) != 1 {
			return fmt.Errorf("need only one conditions.")
		}

		condition := conditions[0].(map[string]interface{})
		var policy = monitor.AlarmPolicyCondition{}
		policy.IsUnionRule = helper.IntInt64(condition["is_union_rule"].(int))
		policy.Rules = make([]*monitor.AlarmPolicyRule, 0, 10)

		rules := condition["rules"]

		for _, item := range rules.([]interface{}) {
			m := item.(map[string]interface{})
			alarmPolicyRule := monitor.AlarmPolicyRule{}
			if m["metric_name"] != nil {
				alarmPolicyRule.MetricName = helper.String(m["metric_name"].(string))
			}
			if m["period"] != nil {
				alarmPolicyRule.Period = helper.IntInt64(m["period"].(int))
			}
			if m["operator"] != nil {
				alarmPolicyRule.Operator = helper.String(m["operator"].(string))
			}
			if m["value"] != nil {
				alarmPolicyRule.Value = helper.String(m["value"].(string))
			}
			if m["continue_period"] != nil {
				alarmPolicyRule.ContinuePeriod = helper.IntInt64(m["continue_period"].(int))
			}
			if m["notice_frequency"] != nil {
				alarmPolicyRule.NoticeFrequency = helper.IntInt64(m["notice_frequency"].(int))
			}
			if m["is_power_notice"] != nil {
				alarmPolicyRule.IsPowerNotice = helper.IntInt64(m["is_power_notice"].(int))
			}
			if v, ok := m["filter"]; ok {
				filters := v.([]interface{})
				if len(filters) > 0 {
					filter := filters[0].(map[string]interface{})
					alarmPolicyFilter := monitor.AlarmPolicyFilter{
						Type:       helper.String(filter["type"].(string)),
						Dimensions: helper.String(filter["dimensions"].(string)),
					}
					alarmPolicyRule.Filter = &alarmPolicyFilter
				}
			}

			if m["description"] != nil {
				alarmPolicyRule.Description = helper.String(m["description"].(string))
			}
			if m["unit"] != nil {
				alarmPolicyRule.Unit = helper.String(m["unit"].(string))
			}
			if m["rule_type"] != nil {
				alarmPolicyRule.RuleType = helper.String(m["rule_type"].(string))
			}
			policy.Rules = append(policy.Rules, &alarmPolicyRule)
		}
		request.Condition = &policy
	}

	if v, ok := d.GetOk("event_conditions"); ok {
		eventCondition := monitor.AlarmPolicyEventCondition{}
		rules := make([]*monitor.AlarmPolicyRule, 0, 10)
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			alarmPolicyRule := monitor.AlarmPolicyRule{}
			if m["metric_name"] != nil {
				alarmPolicyRule.MetricName = helper.String(m["metric_name"].(string))
			}
			if m["period"] != nil {
				alarmPolicyRule.Period = helper.IntInt64(m["period"].(int))
			}
			if m["operator"] != nil {
				alarmPolicyRule.Operator = helper.String(m["operator"].(string))
			}
			if m["value"] != nil {
				alarmPolicyRule.Value = helper.String(m["value"].(string))
			}
			if m["continue_period"] != nil {
				alarmPolicyRule.ContinuePeriod = helper.IntInt64(m["continue_period"].(int))
			}
			if m["notice_frequency"] != nil {
				alarmPolicyRule.NoticeFrequency = helper.IntInt64(m["notice_frequency"].(int))
			}
			if m["is_power_notice"] != nil {
				alarmPolicyRule.IsPowerNotice = helper.IntInt64(m["is_power_notice"].(int))
			}
			if m["filter"] != nil {
				filters := m["filter"].([]interface{})
				if len(filters) > 0 {
					filter := filters[0].(map[string]interface{})
					alarmPolicyFilter := monitor.AlarmPolicyFilter{
						Type:       helper.String(filter["type"].(string)),
						Dimensions: helper.String(filter["dimensions"].(string)),
					}
					alarmPolicyRule.Filter = &alarmPolicyFilter
				}
			}
			if m["description"] != nil {
				alarmPolicyRule.Description = helper.String(m["description"].(string))
			}
			if m["unit"] != nil {
				alarmPolicyRule.Unit = helper.String(m["unit"].(string))
			}
			if m["rule_type"] != nil {
				alarmPolicyRule.RuleType = helper.String(m["rule_type"].(string))
			}
			rules = append(rules, &alarmPolicyRule)
		}
		eventCondition.Rules = rules
		request.EventCondition = &eventCondition
	}

	if v, ok := d.GetOk("notice_ids"); ok {
		notice := make([]*string, 0, 10)
		for _, item := range v.([]interface{}) {
			notice = append(notice, helper.String(item.(string)))
		}
		request.NoticeIds = notice
	}

	if v, ok := d.GetOk("trigger_tasks"); ok {
		tasks := make([]*monitor.AlarmPolicyTriggerTask, 0, 10)
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			triggerTask := monitor.AlarmPolicyTriggerTask{}
			triggerTask.Type = helper.String(m["type"].(string))
			triggerTask.TaskConfig = helper.String(m["task_config"].(string))
			tasks = append(tasks, &triggerTask)
		}
		request.TriggerTasks = tasks
	}

	if v, ok := d.GetOk("group_by"); ok {
		groupBySet := v.(*schema.Set).List()
		for i := range groupBySet {
			groupBy := groupBySet[i].(string)
			request.GroupBy = append(request.GroupBy, &groupBy)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "filter"); ok {
		alarmPolicyFilter := monitor.AlarmPolicyFilter{}
		if v, ok := dMap["type"]; ok {
			alarmPolicyFilter.Type = helper.String(v.(string))
		}
		if v, ok := dMap["dimensions"]; ok {
			alarmPolicyFilter.Dimensions = helper.String(v.(string))
		}
		request.Filter = &alarmPolicyFilter
	}

	var groupId *string
	var policyId *string
	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := monitorService.client.UseMonitorClient().CreateAlarmPolicy(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		policyId = response.Response.PolicyId
		groupId = response.Response.OriginId
		return nil
	}); err != nil {
		return err
	}
	d.SetId(*policyId)

	// binding tag
	if v, ok := d.GetOk("policy_tag"); ok {
		request := monitor.NewBindingPolicyTagRequest()

		request.Module = helper.String("monitor")
		request.PolicyId = helper.String(*policyId)
		request.ServiceType = helper.String(d.Get("namespace").(string))
		request.GroupId = helper.String(*groupId)
		tagSet := make([]*monitor.PolicyTag, 0, 10)
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			tagInfo := monitor.PolicyTag{
				Key:   helper.String(m["key"].(string)),
				Value: helper.String(m["value"].(string)),
			}
			tagSet = append(tagSet, &tagInfo)
		}
		request.Tag = tagSet[0]

		if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err := monitorService.client.UseMonitorClient().BindingPolicyTag(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), err.Error())
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			return nil
		}); err != nil {
			return err
		}
	}

	return resourceTencentMonitorAlarmPolicyRead(d, meta)
}

func resourceTencentMonitorAlarmPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_alarm_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	//logId := tccommon.GetLogId(tccommon.ContextNil)
	//ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request        = monitor.NewDescribeAlarmPolicyRequest()
		policy         *monitor.AlarmPolicy
	)

	policyId := d.Id()
	request.PolicyId = &policyId
	request.Module = helper.String("monitor")

	if err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := monitorService.client.UseMonitorClient().DescribeAlarmPolicy(request)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		policy = response.Response.Policy
		return nil
	}); err != nil {
		return err
	}

	if policy == nil {
		d.SetId("")
		return nil
	}

	var errs []error
	errs = append(errs,
		d.Set("policy_name", policy.PolicyName),
		d.Set("monitor_type", policy.MonitorType),
		d.Set("namespace", policy.Namespace),
		d.Set("remark", policy.Remark),
		d.Set("enable", policy.Enable),
		d.Set("project_id", policy.ProjectId),
	)

	if policy.GroupBy != nil {
		groupBy := []string{}
		for _, v := range policy.GroupBy {
			groupBy = append(groupBy, *v.Id)
		}
		_ = d.Set("group_by", groupBy)
	}

	if policy.Filter != nil {
		filterMap := map[string]interface{}{}

		if policy.Filter.Type != nil {
			filterMap["type"] = policy.Filter.Type
		}

		if policy.Filter.Dimensions != nil {
			filterMap["dimensions"] = policy.Filter.Dimensions
		}

		_ = d.Set("filter", []interface{}{filterMap})
	}

	//nolint:misspell
	if policy.ConditionTemplateId != nil && *policy.ConditionTemplateId != "" {
		id, err := strconv.ParseInt(*policy.ConditionTemplateId, 10, 64)
		if id != 0 && err == nil {
			errs = append(errs, d.Set("conditon_template_id", id))
		}
	}

	if policy.InsertTime != nil {
		t := time.Unix(*policy.InsertTime, 0)
		tFmt := t.Format("2006-01-02 15:04:05")
		errs = append(errs, d.Set("create_time", tFmt))
	}
	if policy.UpdateTime != nil {
		t := time.Unix(*policy.UpdateTime, 0)
		tFmt := t.Format("2006-01-02 15:04:05")
		errs = append(errs, d.Set("update_time", tFmt))
	}

	var rules = make([]interface{}, 0, 100)
	for _, rule := range policy.Condition.Rules {

		m := map[string]interface{}{
			"metric_name":      rule.MetricName,
			"period":           rule.Period,
			"operator":         rule.Operator,
			"value":            rule.Value,
			"continue_period":  rule.ContinuePeriod,
			"notice_frequency": rule.NoticeFrequency,
			"description":      rule.Description,
			"unit":             rule.Unit,
			"rule_type":        rule.RuleType,
		}
		if rule.Filter != nil {
			if *rule.Filter.Type != "" || *rule.Filter.Dimensions != "" {
				var filter = make([]interface{}, 0, 10)
				alarmPolicyFilter := map[string]interface{}{
					"type":       rule.Filter.Type,
					"dimensions": rule.Filter.Dimensions,
				}
				filter = append(filter, alarmPolicyFilter)
				if len(filter) > 0 {
					m["filter"] = filter
				}
			}
		}

		rules = append(rules, m)
	}

	conditions := map[string]interface{}{
		"is_union_rule": policy.Condition.IsUnionRule,
		"rules":         rules,
	}
	_ = d.Set("conditions", []interface{}{conditions})

	eventConditions := make([]map[string]interface{}, 0, len(policy.EventCondition.Rules))
	for _, eventRule := range policy.EventCondition.Rules {

		m := make(map[string]interface{}, 5)
		m["metric_name"] = eventRule.MetricName
		m["period"] = eventRule.Period
		m["operator"] = eventRule.Operator
		m["value"] = eventRule.Value
		m["continue_period"] = eventRule.ContinuePeriod
		m["notice_frequency"] = eventRule.NoticeFrequency
		m["is_power_notice"] = eventRule.IsPowerNotice
		m["notice_frequency"] = eventRule.NoticeFrequency
		m["description"] = eventRule.Description
		m["unit"] = eventRule.Unit
		m["rule_type"] = eventRule.RuleType
		if eventRule.Filter != nil {
			if *eventRule.Filter.Type != "" || *eventRule.Filter.Dimensions != "" {
				var filter = make([]interface{}, 0, 10)
				alarmPolicyFilter := map[string]interface{}{
					"type":       eventRule.Filter.Type,
					"dimensions": eventRule.Filter.Dimensions,
				}
				filter = append(filter, alarmPolicyFilter)
				if len(filter) > 0 {
					m["filter"] = filter
				}
			}
		}
		eventConditions = append(eventConditions, m)
	}
	_ = d.Set("event_conditions", eventConditions)
	var noticeIds = make([]interface{}, 0, 100)
	for _, notice := range policy.NoticeIds {
		noticeIds = append(noticeIds, notice)
	}
	errs = append(errs, d.Set("notice_ids", noticeIds))

	var triggerTasks = make([]interface{}, 0, 100)
	for _, task := range policy.TriggerTasks {
		m := map[string]interface{}{}
		m["type"] = task.Type
		m["task_config"] = task.TaskConfig
		triggerTasks = append(triggerTasks, m)
	}
	errs = append(errs, d.Set("trigger_tasks", triggerTasks))

	tagSets := make([]map[string]interface{}, 0, len(policy.TagInstances))
	for _, item := range policy.TagInstances {
		tagSets = append(tagSets, map[string]interface{}{
			"key":   item.Key,
			"value": item.Value,
		})
	}
	_ = d.Set("policy_tag", tagSets)

	var errResults *multierror.Error
	for i := range errs {
		err := errs[i]
		if err != nil {
			errResults = multierror.Append(errResults, err)
		}
	}
	return errResults.ErrorOrNil()
}

func resourceTencentMonitorAlarmPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_alarm_policy.update")()
	//logId := tccommon.GetLogId(tccommon.ContextNil)
	//ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	if d.HasChange("policy_name") {
		request := monitor.NewModifyAlarmPolicyInfoRequest()
		request.Module = helper.String("monitor")
		request.PolicyId = helper.String(d.Id())
		request.Key = helper.String("NAME")
		value := d.Get("policy_name").(string)
		request.Value = helper.String(value)

		if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			if _, err := monitorService.client.UseMonitorClient().ModifyAlarmPolicyInfo(request); err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	if d.HasChange("remark") {
		request := monitor.NewModifyAlarmPolicyInfoRequest()
		request.Module = helper.String("monitor")
		request.PolicyId = helper.String(d.Id())
		request.Key = helper.String("REMARK")
		value := d.Get("remark").(string)
		request.Value = helper.String(value)

		if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			if _, err := monitorService.client.UseMonitorClient().ModifyAlarmPolicyInfo(request); err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	if d.HasChange("enable") {
		request := monitor.NewModifyAlarmPolicyStatusRequest()
		request.Module = helper.String("monitor")
		request.PolicyId = helper.String(d.Id())

		enable := d.Get("enable").(int)
		request.Enable = helper.IntInt64(enable)

		if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			if _, err := monitorService.client.UseMonitorClient().ModifyAlarmPolicyStatus(request); err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	if d.HasChange("conditions") || d.HasChange("event_conditions") || d.HasChange("group_by") || d.HasChange("filter") {
		request := monitor.NewModifyAlarmPolicyConditionRequest()
		request.Module = helper.String("monitor")
		request.PolicyId = helper.String(d.Id())

		if v, ok := d.GetOk("conditions"); ok {
			conditions := v.([]interface{})

			if len(conditions) != 1 {
				return fmt.Errorf("need only one conditions.")
			}

			condition := conditions[0].(map[string]interface{})
			var policy = monitor.AlarmPolicyCondition{}
			policy.IsUnionRule = helper.IntInt64(condition["is_union_rule"].(int))
			policy.Rules = make([]*monitor.AlarmPolicyRule, 0, 10)

			rules := condition["rules"]

			for _, item := range rules.([]interface{}) {
				m := item.(map[string]interface{})
				alarmPolicyRule := monitor.AlarmPolicyRule{}
				if m["metric_name"] != nil {
					alarmPolicyRule.MetricName = helper.String(m["metric_name"].(string))
				}
				if m["period"] != nil {
					alarmPolicyRule.Period = helper.IntInt64(m["period"].(int))
				}
				if m["value"] != nil {
					alarmPolicyRule.Value = helper.String(m["value"].(string))
				}
				if m["operator"] != nil {
					alarmPolicyRule.Operator = helper.String(m["operator"].(string))
				}
				if m["continue_period"] != nil {
					alarmPolicyRule.ContinuePeriod = helper.IntInt64(m["continue_period"].(int))
				}
				if m["notice_frequency"] != nil {
					alarmPolicyRule.NoticeFrequency = helper.IntInt64(m["notice_frequency"].(int))
				}
				if m["is_power_notice"] != nil {
					alarmPolicyRule.IsPowerNotice = helper.IntInt64(m["is_power_notice"].(int))
				}
				if m["filter"] != nil {
					filters := m["filter"].([]interface{})
					// Max Items is 1
					if len(filters) > 0 {
						filter := filters[0].(map[string]interface{})
						alarmPolicyFilter := monitor.AlarmPolicyFilter{
							Type:       helper.String(filter["type"].(string)),
							Dimensions: helper.String(filter["dimensions"].(string)),
						}
						alarmPolicyRule.Filter = &alarmPolicyFilter
					}
				}
				if m["description"] != nil {
					alarmPolicyRule.Description = helper.String(m["description"].(string))
				}
				if m["unit"] != nil {
					alarmPolicyRule.Unit = helper.String(m["unit"].(string))
				}
				if m["rule_type"] != nil {
					alarmPolicyRule.RuleType = helper.String(m["rule_type"].(string))
				}
				policy.Rules = append(policy.Rules, &alarmPolicyRule)
			}
			request.Condition = &policy
		}
		if v, ok := d.GetOk("event_conditions"); ok {
			eventCondition := monitor.AlarmPolicyEventCondition{}
			rules := make([]*monitor.AlarmPolicyRule, 0, 10)
			for _, item := range v.([]interface{}) {
				m := item.(map[string]interface{})
				alarmPolicyRule := monitor.AlarmPolicyRule{}
				if m["metric_name"] != nil {
					alarmPolicyRule.MetricName = helper.String(m["metric_name"].(string))
				}
				if m["period"] != nil {
					alarmPolicyRule.Period = helper.IntInt64(m["period"].(int))
				}
				if m["operator"] != nil {
					alarmPolicyRule.Operator = helper.String(m["operator"].(string))
				}
				if m["value"] != nil {
					alarmPolicyRule.Value = helper.String(m["value"].(string))
				}
				if m["continue_period"] != nil {
					alarmPolicyRule.ContinuePeriod = helper.IntInt64(m["continue_period"].(int))
				}
				if m["notice_frequency"] != nil {
					alarmPolicyRule.NoticeFrequency = helper.IntInt64(m["notice_frequency"].(int))
				}
				if m["is_power_notice"] != nil {
					alarmPolicyRule.IsPowerNotice = helper.IntInt64(m["is_power_notice"].(int))
				}
				if m["filter"] != nil {
					filters := m["filter"].([]interface{})
					// Max Items is 1
					if len(filters) > 0 {
						filter := filters[0].(map[string]interface{})
						alarmPolicyFilter := monitor.AlarmPolicyFilter{
							Type:       helper.String(filter["type"].(string)),
							Dimensions: helper.String(filter["dimensions"].(string)),
						}
						alarmPolicyRule.Filter = &alarmPolicyFilter
					}
				}
				if m["description"] != nil {
					alarmPolicyRule.Description = helper.String(m["description"].(string))
				}
				if m["unit"] != nil {
					alarmPolicyRule.Unit = helper.String(m["unit"].(string))
				}
				if m["rule_type"] != nil {
					alarmPolicyRule.RuleType = helper.String(m["rule_type"].(string))
				}
				rules = append(rules, &alarmPolicyRule)
			}
			eventCondition.Rules = rules
			request.EventCondition = &eventCondition
		}

		if v, ok := d.GetOk("group_by"); ok {
			groupBySet := v.(*schema.Set).List()
			for i := range groupBySet {
				groupBy := groupBySet[i].(string)
				request.GroupBy = append(request.GroupBy, &groupBy)
			}
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "filter"); ok {
			alarmPolicyFilter := monitor.AlarmPolicyFilter{}
			if v, ok := dMap["type"]; ok {
				alarmPolicyFilter.Type = helper.String(v.(string))
			}
			if v, ok := dMap["dimensions"]; ok {
				alarmPolicyFilter.Dimensions = helper.String(v.(string))
			}
			request.Filter = &alarmPolicyFilter
		}

		if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			if _, err := monitorService.client.UseMonitorClient().ModifyAlarmPolicyCondition(request); err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	if d.HasChange("notice_ids") {
		request := monitor.NewModifyAlarmPolicyNoticeRequest()
		request.Module = helper.String("monitor")
		request.PolicyId = helper.String(d.Id())

		if v, ok := d.GetOk("notice_ids"); ok {
			notice := make([]*string, 0, 10)
			for _, item := range v.([]interface{}) {
				notice = append(notice, helper.String(item.(string)))
			}
			request.NoticeIds = notice
		}

		if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			if _, err := monitorService.client.UseMonitorClient().ModifyAlarmPolicyNotice(request); err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	if d.HasChange("trigger_tasks") {
		request := monitor.NewModifyAlarmPolicyTasksRequest()
		request.Module = helper.String("monitor")
		request.PolicyId = helper.String(d.Id())
		if v, ok := d.GetOk("trigger_tasks"); ok {
			tasks := make([]*monitor.AlarmPolicyTriggerTask, 0, 10)
			for _, item := range v.([]interface{}) {
				m := item.(map[string]interface{})
				triggerTask := monitor.AlarmPolicyTriggerTask{}
				triggerTask.Type = helper.String(m["type"].(string))
				triggerTask.TaskConfig = helper.String(m["task_config"].(string))
				tasks = append(tasks, &triggerTask)
			}
			request.TriggerTasks = tasks
		}
		if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			if _, err := monitorService.client.UseMonitorClient().ModifyAlarmPolicyTasks(request); err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return resourceTencentMonitorAlarmPolicyRead(d, meta)
}

func resourceTencentMonitorAlarmPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_alarm_policy.delete")()

	var (
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request        = monitor.NewDeleteAlarmPolicyRequest()
	)
	request.Module = helper.String("monitor")
	policyIds := []*string{helper.String(d.Id())}
	request.PolicyIds = policyIds

	if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if _, err := monitorService.client.UseMonitorClient().DeleteAlarmPolicy(request); err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
