package monitor

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMonitorAlarmBasicMetric() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMonitorAlarmBasicMetricRead,
		Schema: map[string]*schema.Schema{
			"namespace": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "business 命名空间 是 different 对于 each 云 product. To obtain business 命名空间，please go 到 product 监控 indicator documents，such 作为 命名空间 的 云 服务器，其中 可以 是 found 在 [Cloud Server Monitoring Indicators](https://云.tencent.com/document/product/248/6843 )。",
			},

			"metric_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Indicator names 是 different 对于 each 云 product. To obtain indicator names，please go 到 监控 indicator documents 的 each product，such 作为 indicator names 的 云 servers，其中 可以 是 found 在 [Cloud Server Monitoring Indicators]( https://云.tencent.com/document/product/248/6843)。",
			},

			"dimensions": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "可选 参数，filtered 通过 dimension。",
			},

			"metric_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 indicator descriptions 获取 从 查询。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespaces，each 云 product 将 have 命名空间。",
						},
						"metric_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicator 名称",
						},
						"unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Units 用于indicators。",
						},
						"unit_cname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Units 用于indicators。",
						},
						"period": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "statistical 周期 支持 通过 indicator，（秒）， such 作为 60，300。",
						},
						"periods": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Indicator 方法 within statistical cycle。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"period": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cycle。",
									},
									"stat_type": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Statistical methods。",
									},
								},
							},
						},
						"meaning": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Explanation 的 meaning 的 statistical indicators。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"en": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Explanation 的 indicators 在 English。",
									},
									"zh": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Chinese interpretation 的 indicators。",
									},
								},
							},
						},
						"dimensions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Dimension 描述 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dimensions": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Dimension 名称 数组。",
									},
								},
							},
						},
						"metric_c_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicator Chinese 名称",
						},
						"metric_e_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicator English 名称",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudMonitorAlarmBasicMetricRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_alarm_metric.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("namespace"); ok {
		paramMap["Namespace"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("metric_name"); ok {
		paramMap["MetricName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dimensions"); ok {
		dimensionsSet := v.(*schema.Set).List()
		paramMap["Dimensions"] = helper.InterfacesStringsPoint(dimensionsSet)
	}

	service := MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var metricSet []*monitor.MetricSet

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMonitorAlarmBasicMetricByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		metricSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(metricSet))
	tmpList := make([]map[string]interface{}, 0, len(metricSet))

	if metricSet != nil {
		for _, metricSet := range metricSet {
			metricSetMap := map[string]interface{}{}

			if metricSet.Namespace != nil {
				metricSetMap["namespace"] = metricSet.Namespace
			}

			if metricSet.MetricName != nil {
				metricSetMap["metric_name"] = metricSet.MetricName
			}

			if metricSet.Unit != nil {
				metricSetMap["unit"] = metricSet.Unit
			}

			if metricSet.UnitCname != nil {
				metricSetMap["unit_cname"] = metricSet.UnitCname
			}

			if metricSet.Period != nil {
				metricSetMap["period"] = metricSet.Period
			}

			if metricSet.Periods != nil {
				periodsList := []interface{}{}
				for _, periods := range metricSet.Periods {
					periodsMap := map[string]interface{}{}

					if periods.Period != nil {
						periodsMap["period"] = periods.Period
					}

					if periods.StatType != nil {
						periodsMap["stat_type"] = periods.StatType
					}

					periodsList = append(periodsList, periodsMap)
				}

				metricSetMap["periods"] = periodsList
			}

			if metricSet.Meaning != nil {
				meaningMap := map[string]interface{}{}

				if metricSet.Meaning.En != nil {
					meaningMap["en"] = metricSet.Meaning.En
				}

				if metricSet.Meaning.Zh != nil {
					meaningMap["zh"] = metricSet.Meaning.Zh
				}

				metricSetMap["meaning"] = []interface{}{meaningMap}
			}

			if metricSet.Dimensions != nil {
				dimensionsList := []interface{}{}
				for _, dimensions := range metricSet.Dimensions {
					dimensionsMap := map[string]interface{}{}

					if dimensions.Dimensions != nil {
						dimensionsMap["dimensions"] = dimensions.Dimensions
					}

					dimensionsList = append(dimensionsList, dimensionsMap)
				}

				metricSetMap["dimensions"] = dimensionsList
			}

			if metricSet.MetricCName != nil {
				metricSetMap["metric_c_name"] = metricSet.MetricCName
			}

			if metricSet.MetricEName != nil {
				metricSetMap["metric_e_name"] = metricSet.MetricEName
			}

			ids = append(ids, *metricSet.MetricName)
			tmpList = append(tmpList, metricSetMap)
		}

		_ = d.Set("metric_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
