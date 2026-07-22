package eb

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudEbSearch() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEbSearchRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "开始时间。",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "结束时间。",
			},

			"event_bus_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "事件 bus ID。",
			},

			"group_field": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "aggregate 字段，当 querying 日志 索引 dimension 值，您 必须 enter。",
			},

			"filter": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 criteria。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "过滤字段名称",
						},
						"operator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "操作者，congruent eq，不 equal neq，similar like，exclude similar 不 like，less 比 lt，less 比 和 equal 到 lte，greater 比 gt，greater 比 和 equal 到 gte，在 范围 范围，不 在 范围 norange。",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "过滤值，范围 operation needs 到 enter two 值 在 same 时间，separated 通过 commas。",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "logical relationship 的 级别 filters， 值 AND 或 OR。",
						},
						"filters": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "LogFilters 数组。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "过滤字段名称",
									},
									"operator": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "操作者，congruent eq，不 equal neq，similar like，exclude similar 不 like，less 比 lt，less 比 和 equal 到 lte，greater 比 gt，greater 比 和 equal 到 gte，within 范围 范围，不 within 范围 norange。",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "过滤器 值，范围 operations need 到 enter two 值 在 same 时间，separated 通过 commas。",
									},
								},
							},
						},
					},
				},
			},

			"order_fields": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "sort 数组，take effect 当 日志 是 retrieved。",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "排序方式，asc 从 old 到 new，desc 从 new 到 old，take effect 当 日志 是 retrieved。",
			},

			"dimension_values": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "索引 retrieves dimension 值。",
			},

			"results": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Log search results，note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "报告 时间 的 单个 日志，note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log 内容 details，note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event 来源，note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event 类型，note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"rule_ids": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event matching 规则，note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"subject": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID，note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域，注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event 状态，note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
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

func dataSourceTencentCloudEbSearchRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_eb_search.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		startTime  string
		endTime    string
		eventBusId string
		groupField string
	)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("start_time"); v != nil {
		startTime = strconv.Itoa(v.(int))
		paramMap["StartTime"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("end_time"); v != nil {
		endTime = strconv.Itoa(v.(int))
		paramMap["EndTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("event_bus_id"); ok {
		eventBusId = v.(string)
		paramMap["EventBusId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_field"); ok {
		groupField = v.(string)
		paramMap["GroupField"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter"); ok {
		filterSet := v.([]interface{})
		tmpSet := make([]*eb.LogFilter, 0, len(filterSet))

		for _, item := range filterSet {
			logFilter := eb.LogFilter{}
			logFilterMap := item.(map[string]interface{})

			if v, ok := logFilterMap["key"]; ok {
				logFilter.Key = helper.String(v.(string))
			}
			if v, ok := logFilterMap["operator"]; ok {
				logFilter.Operator = helper.String(v.(string))
			}
			if v, ok := logFilterMap["value"]; ok {
				logFilter.Value = helper.String(v.(string))
			}
			if v, ok := logFilterMap["type"]; ok {
				logFilter.Type = helper.String(v.(string))
			}
			if v, ok := logFilterMap["filters"]; ok {
				for _, item := range v.([]interface{}) {
					filtersMap := item.(map[string]interface{})
					logFilters := eb.LogFilters{}
					if v, ok := filtersMap["key"]; ok {
						logFilters.Key = helper.String(v.(string))
					}
					if v, ok := filtersMap["operator"]; ok {
						logFilters.Operator = helper.String(v.(string))
					}
					if v, ok := filtersMap["value"]; ok {
						logFilters.Value = helper.String(v.(string))
					}
					logFilter.Filters = append(logFilter.Filters, &logFilters)
				}
			}
			tmpSet = append(tmpSet, &logFilter)
		}
		paramMap["Filter"] = tmpSet
	}

	if v, ok := d.GetOk("order_fields"); ok {
		orderFieldsSet := v.(*schema.Set).List()
		paramMap["OrderFields"] = helper.InterfacesStringsPoint(orderFieldsSet)
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	service := EbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	if groupField != "" {
		var searchResults []*string
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			response, e := service.DescribeEbSearchByFilter(ctx, paramMap)
			if e != nil {
				return tccommon.RetryError(e)
			}
			searchResults = response
			return nil
		})
		if err != nil {
			return err
		}

		if searchResults != nil {
			_ = d.Set("dimension_values", searchResults)
		}
	}

	var results []*eb.SearchLogResult
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeEbSearchLogByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		results = response
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(results))
	if results != nil {
		for _, searchLogResult := range results {
			searchLogResultMap := map[string]interface{}{}

			if searchLogResult.Timestamp != nil {
				searchLogResultMap["timestamp"] = searchLogResult.Timestamp
			}

			if searchLogResult.Message != nil {
				searchLogResultMap["message"] = searchLogResult.Message
			}

			if searchLogResult.Source != nil {
				searchLogResultMap["source"] = searchLogResult.Source
			}

			if searchLogResult.Type != nil {
				searchLogResultMap["type"] = searchLogResult.Type
			}

			if searchLogResult.RuleIds != nil {
				searchLogResultMap["rule_ids"] = searchLogResult.RuleIds
			}

			if searchLogResult.Subject != nil {
				searchLogResultMap["subject"] = searchLogResult.Subject
			}

			if searchLogResult.Region != nil {
				searchLogResultMap["region"] = searchLogResult.Region
			}

			if searchLogResult.Status != nil {
				searchLogResultMap["status"] = searchLogResult.Status
			}

			tmpList = append(tmpList, searchLogResultMap)
		}

		_ = d.Set("results", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash([]string{startTime, endTime, eventBusId}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
