package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataResourceGroupMetrics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataResourceGroupMetricsRead,
		Schema: map[string]*schema.Schema{
			"resource_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Execution 资源 组 ID",
			},

			"start_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Usage trend 开始时间 (milliseconds)，默认为 last hour。",
			},

			"end_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Usage trend 结束时间 (milliseconds)，默认为 当前 时间。",
			},

			"metric_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Metric dimension.\n\n- all --- All\n- 任务 --- 任务 metrics\n- 系统 --- System metrics。",
			},

			"granularity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Metric collection granularity，单位 在 minutes，默认值 1 minute。",
			},

			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Execution 组 metric 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Resource 组 规格 related: CPU count。",
						},
						"disk_volume": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Resource 组 规格 related: 磁盘 规格。",
						},
						"mem_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Resource 组 规格 related: 内存 大小，单位: G。",
						},
						"life_cycle": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Resource 组 lifecycle，单位: days。",
						},
						"maximum_concurrency": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Resource 组 规格 related: 最大 并发",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Resource 组 状态\n\n- 0 --- Initializing\n- 1 --- Running\n- 2 --- Running abnormally\n- 3 --- Releasing\n- 4 --- Released\n- 5 --- Creating\n- 6 --- Creation failed\n- 7 --- Updating\n- 8 --- Update failed\n- 9 --- Expired\n- 10 --- Release failed\n- 11 --- In 使用\n- 12 --- Not 在 使用。",
						},
						"metric_snapshots": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Metric details。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "指标名称\n\n- ConcurrencyUsage --- 并发 usage 速率\n- CpuCoreUsage --- CPU usage 速率\n- CpuLoad --- CPU load\n- DevelopQueueTask --- 数量 development tasks 在 queue\n- DevelopRunningTask --- 数量 running development tasks\n- DevelopSchedulingTask --- 数量 scheduling development tasks\n- DiskUsage --- Disk usage\n- DiskUsed --- Disk 使用 amount\n- MaximumConcurrency --- Maximum 并发\n- MemoryLoad --- Memory load\n- MemoryUsage --- Memory usage。",
									},
									"snapshot_value": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "当前值",
									},
									"trend_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Metric trend。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"timestamp": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "时间戳。",
												},
												"value": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "指标值",
												},
											},
										},
									},
								},
							},
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

func dataSourceTencentCloudWedataResourceGroupMetricsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_resource_group_metrics.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId           = tccommon.GetLogId(nil)
		ctx             = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service         = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		resourceGroupId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("resource_group_id"); ok {
		paramMap["ResourceGroupId"] = helper.String(v.(string))
		resourceGroupId = v.(string)
	}

	if v, ok := d.GetOkExists("start_time"); ok {
		paramMap["StartTime"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("end_time"); ok {
		paramMap["EndTime"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("metric_type"); ok {
		paramMap["MetricType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("granularity"); ok {
		paramMap["Granularity"] = helper.IntUint64(v.(int))
	}

	var respData *wedatav20250806.ResourceGroupMetrics
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataResourceGroupMetricsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	dataMap := map[string]interface{}{}
	if respData != nil {
		if respData.CpuNum != nil {
			dataMap["cpu_num"] = respData.CpuNum
		}

		if respData.DiskVolume != nil {
			dataMap["disk_volume"] = respData.DiskVolume
		}

		if respData.MemSize != nil {
			dataMap["mem_size"] = respData.MemSize
		}

		if respData.LifeCycle != nil {
			dataMap["life_cycle"] = respData.LifeCycle
		}

		if respData.MaximumConcurrency != nil {
			dataMap["maximum_concurrency"] = respData.MaximumConcurrency
		}

		if respData.Status != nil {
			dataMap["status"] = respData.Status
		}

		metricSnapshotsList := make([]map[string]interface{}, 0, len(respData.MetricSnapshots))
		if respData.MetricSnapshots != nil {
			for _, metricSnapshots := range respData.MetricSnapshots {
				metricSnapshotsMap := map[string]interface{}{}
				if metricSnapshots.MetricName != nil {
					metricSnapshotsMap["metric_name"] = metricSnapshots.MetricName
				}

				if metricSnapshots.SnapshotValue != nil {
					metricSnapshotsMap["snapshot_value"] = metricSnapshots.SnapshotValue
				}

				trendListList := make([]map[string]interface{}, 0, len(metricSnapshots.TrendList))
				if metricSnapshots.TrendList != nil {
					for _, trendList := range metricSnapshots.TrendList {
						trendListMap := map[string]interface{}{}

						if trendList.Timestamp != nil {
							trendListMap["timestamp"] = trendList.Timestamp
						}

						if trendList.Value != nil {
							trendListMap["value"] = trendList.Value
						}

						trendListList = append(trendListList, trendListMap)
					}

					metricSnapshotsMap["trend_list"] = trendListList
				}

				metricSnapshotsList = append(metricSnapshotsList, metricSnapshotsMap)
			}

			dataMap["metric_snapshots"] = metricSnapshotsList
		}

		_ = d.Set("data", []interface{}{dataMap})
	}

	d.SetId(resourceGroupId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataMap); e != nil {
			return e
		}
	}

	return nil
}
