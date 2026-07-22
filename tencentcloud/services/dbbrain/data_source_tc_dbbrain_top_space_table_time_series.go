package dbbrain

import (
	"context"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbbrainTopSpaceTableTimeSeries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainTopSpaceTableTimeSeriesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Default:     20,
				Description: "数量 的 Top tables 返回, 最大 值 是 100, 和 默认值 是 20.",
			},

			"sort_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "sorting 字段 使用 到 过滤器 Top 表. 可选 字段 include DataLength, IndexLength, TotalLength, DataFree, FragRatio, TableRows, 和 PhysicalFileSize. 默认值 是 PhysicalFileSize.",
			},

			"start_date": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "start date, such 作为 2021-01-01, earliest 是 29th day before 当前 day, 和 默认值 是 6th day before deadline.",
			},

			"end_date": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "deadline, such 作为 2021-01-01, earliest 是 29th day before 当前 day, 和 默认值 是 当前 day.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product 类型, 支持 值 include: mysql - 云 数据库 MySQL, cynosdb - 云 数据库 CynosDB 对于 MySQL, 默认值 是 mysql.",
			},

			"top_space_table_time_series": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "时间-series 数据 列表 的 返回 Top tablespace 统计.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"table_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "表 名称.",
						},
						"table_schema": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "databases 名称.",
						},
						"engine": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Storage 引擎 对于 数据库 tables.",
						},
						"series_data": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Spatial index 数据 在 单位 时间 间隔.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"series": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Monitor metrics.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"metric": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Indicator 名称.",
												},
												"unit": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Indicator 单位.",
												},
												"values": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeFloat,
													},
													Description: "Index 值. 注意: 此 字段 可能 返回 null, indicating 该 无 有效 值 可以 是 获取.",
												},
											},
										},
									},
									"timestamp": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
										Computed:    true,
										Description: "timestamp corresponding 到 监控 indicator.",
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
				Description: "Used 到 save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainTopSpaceTableTimeSeriesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_top_space_table_time_series.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var instanceId string

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("sort_by"); ok {
		paramMap["SortBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_date"); ok {
		paramMap["StartDate"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_date"); ok {
		paramMap["EndDate"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var topSpaceTableTimeSeries []*dbbrain.TableSpaceTimeSeries

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainTopSpaceTableTimeSeriesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		topSpaceTableTimeSeries = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(topSpaceTableTimeSeries))
	tmpList := make([]map[string]interface{}, 0, len(topSpaceTableTimeSeries))

	if topSpaceTableTimeSeries != nil {
		for _, tableSpaceTimeSeries := range topSpaceTableTimeSeries {
			tableSpaceTimeSeriesMap := map[string]interface{}{}

			if tableSpaceTimeSeries.TableName != nil {
				tableSpaceTimeSeriesMap["table_name"] = tableSpaceTimeSeries.TableName
			}

			if tableSpaceTimeSeries.TableSchema != nil {
				tableSpaceTimeSeriesMap["table_schema"] = tableSpaceTimeSeries.TableSchema
			}

			if tableSpaceTimeSeries.Engine != nil {
				tableSpaceTimeSeriesMap["engine"] = tableSpaceTimeSeries.Engine
			}

			if tableSpaceTimeSeries.SeriesData != nil {
				seriesDataMap := map[string]interface{}{}

				if tableSpaceTimeSeries.SeriesData.Series != nil {
					seriesList := []interface{}{}
					for _, series := range tableSpaceTimeSeries.SeriesData.Series {
						seriesMap := map[string]interface{}{}

						if series.Metric != nil {
							seriesMap["metric"] = series.Metric
						}

						if series.Unit != nil {
							seriesMap["unit"] = series.Unit
						}

						if series.Values != nil {
							seriesMap["values"] = series.Values
						}

						seriesList = append(seriesList, seriesMap)
					}

					seriesDataMap["series"] = seriesList
				}

				if tableSpaceTimeSeries.SeriesData.Timestamp != nil {
					seriesDataMap["timestamp"] = tableSpaceTimeSeries.SeriesData.Timestamp
				}

				tableSpaceTimeSeriesMap["series_data"] = []interface{}{seriesDataMap}
			}

			ids = append(ids, strings.Join([]string{instanceId, *tableSpaceTimeSeries.TableSchema, *tableSpaceTimeSeries.TableName}, tccommon.FILED_SP))
			tmpList = append(tmpList, tableSpaceTimeSeriesMap)
		}

		_ = d.Set("top_space_table_time_series", tmpList)
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
