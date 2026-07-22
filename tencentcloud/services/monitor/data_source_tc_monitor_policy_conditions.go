package monitor

import (
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func DataSourceTencentCloudMonitorPolicyConditions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentMonitorPolicyConditionRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 策略 名称，support partial matching，eg:`Cloud Virtual Machine`,`Virtual`,`Cloud Load Banlancer-Private CLB Listener`。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于存储结果。",
			},
			// Computed values
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 策略 condition. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_view_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy view 名称，eg:`cvm_device`,`BANDWIDTHPACKAGE`，refer 到 `数据.tencentcloud_monitor_policy_conditions(policy_view_name)`。",
						},
						"is_support_multi_region": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否support multi 地域",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 此 策略 名称",
						},
						"support_regions": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Support regions 的 此 策略 view。",
						},
						"event_metrics": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A 列表 事件 condition metrics. Each element 包含following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"event_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ID 此 事件 metric。",
									},
									"event_show_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "名称 此 事件 metric。",
									},
									"need_recovered": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否recover。",
									},
								},
							},
						},
						"metrics": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A 列表 事件 condition metrics. Each element 包含following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ID 此 metric。",
									},
									"metric_show_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "名称 此 metric。",
									},
									"metric_unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "单位 的 此 metric。",
									},
									"calc_type_keys": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Computed:    true,
										Description: "Calculate 类型 此 metric。",
									},
									"calc_type_need": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether `calc_type` 必填 在 配置。",
									},
									"calc_value_default": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "默认值 calculate 值 的 此 metric。",
									},
									"calc_value_fixed": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "fixed calculate 值 的 此 metric。",
									},
									"calc_value_min": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "min calculate 值 的 此 metric。",
									},
									"calc_value_max": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "max calculate 值 的 此 metric。",
									},
									"calc_value_need": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether `calc_value` 必填 在 配置。",
									},
									"continue_time_default": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "默认值 continue 时间(秒) 配置 对于 此 metric。",
									},
									"continue_time_keys": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Computed:    true,
										Description: "continue 时间(秒) keys 对于 此 metric。",
									},
									"continue_time_need": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether `continue_time` 必填 在 配置。",
									},
									"period_default": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "默认值 数据 时间(秒) 配置 对于 此 metric。",
									},
									"period_keys": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Computed:    true,
										Description: "数据 时间(秒) keys 对于 此 metric。",
									},
									"period_need": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether `周期` 必填 在 配置。",
									},
									"period_num_default": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "默认值 周期 数量 配置 对于 此 metric。",
									},
									"period_num_keys": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Computed:    true,
										Description: "周期 数量 keys 对于 此 metric。",
									},
									"period_num_need": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether `period_num` 必填 在 配置。",
									},
									"stat_type_p5": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data aggregation 模式，cycle 的 5 秒。",
									},
									"stat_type_p10": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data aggregation 模式，cycle 的 10 秒。",
									},
									"stat_type_p60": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data aggregation 模式，cycle 的 60 秒。",
									},
									"stat_type_p300": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data aggregation 模式，cycle 的 300 秒。",
									},
									"stat_type_p600": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data aggregation 模式，cycle 的 600 秒。",
									},
									"stat_type_p1800": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data aggregation 模式，cycle 的 1800 秒。",
									},
									"stat_type_p3600": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data aggregation 模式，cycle 的 3600 秒。",
									},
									"stat_type_p86400": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data aggregation 模式，cycle 的 86400 秒。",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentMonitorPolicyConditionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_policy_conditions.read")()

	var (
		monitorService = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request        = monitor.NewDescribePolicyConditionListRequest()
		response       *monitor.DescribePolicyConditionListResponse
		err            error
		name           = d.Get("name").(string)
		policyViewList = make([]interface{}, 0, 100)
	)

	request.Module = helper.String("monitor")

	if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if response, err = monitorService.client.UseMonitorClient().DescribePolicyConditionList(request); err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	for _, policyViewItem := range response.Response.Conditions {
		if name != "" && !strings.Contains(*policyViewItem.Name, name) {
			continue
		}
		policyViewMap := map[string]interface{}{
			"policy_view_name":        policyViewItem.PolicyViewName,
			"name":                    policyViewItem.Name,
			"is_support_multi_region": policyViewItem.IsSupportMultiRegion,
			"support_regions":         policyViewItem.SupportRegions,
		}
		eventMetrics := make([]interface{}, 0, 100)
		for _, eventItem := range policyViewItem.EventMetrics {
			eventMetrics = append(eventMetrics, map[string]interface{}{
				"event_id":        eventItem.EventId,
				"event_show_name": eventItem.EventShowName,
				"need_recovered":  eventItem.NeedRecovered,
			})
		}
		policyViewMap["event_metrics"] = eventMetrics

		metrics := make([]interface{}, 0, 100)
		for _, metric := range policyViewItem.Metrics {
			metricMap := map[string]interface{}{
				"metric_id":             metric.MetricId,
				"metric_show_name":      metric.MetricShowName,
				"metric_unit":           metric.MetricUnit,
				"calc_type_keys":        metric.ConfigManual.CalcType.Keys,
				"calc_type_need":        metric.ConfigManual.CalcType.Need,
				"calc_value_default":    metric.ConfigManual.CalcValue.Default,
				"calc_value_fixed":      metric.ConfigManual.CalcValue.Fixed,
				"calc_value_min":        metric.ConfigManual.CalcValue.Min,
				"calc_value_max":        metric.ConfigManual.CalcValue.Max,
				"calc_value_need":       metric.ConfigManual.CalcValue.Need,
				"continue_time_default": 0,
				"continue_time_keys":    nil,
				"continue_time_need":    false,
				"period_default":        0,
				"period_keys":           nil,
				"period_need":           false,
				"period_num_default":    0,
				"period_num_keys":       nil,
				"period_num_need":       false,
				"stat_type_p5":          nil,
				"stat_type_p10":         nil,
				"stat_type_p60":         nil,
				"stat_type_p300":        nil,
				"stat_type_p600":        nil,
				"stat_type_p1800":       nil,
				"stat_type_p3600":       nil,
				"stat_type_p86400":      nil,
			}

			if metric.ConfigManual.ContinueTime != nil {
				metricMap["continue_time_default"] = metric.ConfigManual.ContinueTime.Default
				metricMap["continue_time_keys"] = metric.ConfigManual.ContinueTime.Keys
				metricMap["continue_time_need"] = metric.ConfigManual.ContinueTime.Need
			}

			if metric.ConfigManual.Period != nil {
				metricMap["period_default"] = metric.ConfigManual.Period.Default
				metricMap["period_keys"] = metric.ConfigManual.Period.Keys
				metricMap["period_need"] = metric.ConfigManual.Period.Need
			}

			if metric.ConfigManual.PeriodNum != nil {
				metricMap["period_num_default"] = metric.ConfigManual.PeriodNum.Default
				metricMap["period_num_keys"] = metric.ConfigManual.PeriodNum.Keys
				metricMap["period_num_need"] = metric.ConfigManual.PeriodNum.Need
			}

			if metric.ConfigManual.StatType != nil {
				metricMap["stat_type_p5"] = metric.ConfigManual.StatType.P5
				metricMap["stat_type_p10"] = metric.ConfigManual.StatType.P10
				metricMap["stat_type_p60"] = metric.ConfigManual.StatType.P60
				metricMap["stat_type_p300"] = metric.ConfigManual.StatType.P300
				metricMap["stat_type_p600"] = metric.ConfigManual.StatType.P600
				metricMap["stat_type_p1800"] = metric.ConfigManual.StatType.P1800
				metricMap["stat_type_p3600"] = metric.ConfigManual.StatType.P3600
				metricMap["stat_type_p86400"] = metric.ConfigManual.StatType.P86400
			}

			metrics = append(metrics, metricMap)
		}
		policyViewMap["metrics"] = metrics
		policyViewList = append(policyViewList, policyViewMap)
	}
	d.SetId("policy_conditions_" + name)
	if err = d.Set("list", policyViewList); err != nil {
		return err
	}
	if output, ok := d.GetOk("result_output_file"); ok {
		return tccommon.WriteToFile(output.(string), policyViewList)
	}
	return nil
}
