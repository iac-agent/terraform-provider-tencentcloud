package dbbrain

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbbrainSlowLogTopSqls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainSlowLogTopSqlsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start 时间, such 作为 `2019-09-10 12:13:14`.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "deadline, such 作为 `2019-09-11 10:13:14`, 间隔 between deadline 和 start 时间 是 less 比 7 days.",
			},

			"sort_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort 键, currently 支持 sort keys such 作为 QueryTime, ExecTimes, RowsSent, LockTime 和 RowsExamined, 默认值 是 QueryTime.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "sorting 方法 支持 ASC (ascending) 和 DESC (descending). 默认值 是 DESC.",
			},

			"schema_list": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Array 的 数据库 names.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schema": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "DB 名称.",
						},
					},
				},
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product 类型, 支持 值 include: `mysql` - 云 数据库 MySQL, `cynosdb` - 云 数据库 CynosDB 对于 MySQL, 默认值 是 `mysql`.",
			},

			"rows": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Slow 日志 top sql 列表.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lock_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "SQL 总数 lock waiting 时间, 在 秒.",
						},
						"lock_time_max": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Maximum lock waiting 时间, 在 秒.",
						},
						"lock_time_min": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Minimum lock waiting 时间, 在 秒.",
						},
						"rows_examined": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "总数 scan lines.",
						},
						"rows_examined_max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum 数量 的 scan lines.",
						},
						"rows_examined_min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum 数量 的 scan lines.",
						},
						"query_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Total 时间, 在 秒.",
						},
						"query_time_max": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "最大 execution 时间, 在 秒.",
						},
						"query_time_min": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "最小 execution 时间, 在 秒.",
						},
						"rows_sent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "总数 数量 的 rows 返回.",
						},
						"rows_sent_max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum 数量 的 rows 返回.",
						},
						"rows_sent_min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum 数量 的 rows 返回.",
						},
						"exec_times": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Execution times.",
						},
						"sql_template": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "sql template.",
						},
						"sql_text": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SQL 使用 参数 (random).",
						},
						"schema": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database 名称.",
						},
						"query_time_ratio": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Total 时间-consuming ratio, 单位 %.",
						},
						"lock_time_ratio": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "ratio 的 总数 lock waiting 时间 的 SQL, 在 %.",
						},
						"rows_examined_ratio": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "proportion 的 总数 数量 的 scanned lines, 单位 %.",
						},
						"rows_sent_ratio": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "proportion 的 总数 数量 的 rows 返回, 在 %.",
						},
						"query_time_avg": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Average execution 时间, 在 秒.",
						},
						"rows_sent_avg": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "average 数量 的 rows 返回.",
						},
						"lock_time_avg": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Average lock waiting 时间, 在 秒.",
						},
						"rows_examined_avg": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "average 数量 的 lines scanned.",
						},
						"md5": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "MD5 值 的 SOL template.",
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

func dataSourceTencentCloudDbbrainSlowLogTopSqlsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_slow_log_top_sqls.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var id string
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
		id = v.(string)
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["start_time"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["end_time"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_by"); ok {
		paramMap["sort_by"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["order_by"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("schema_list"); ok {
		schemaListSet := v.([]interface{})
		tmpSet := make([]*dbbrain.SchemaItem, 0, len(schemaListSet))

		for _, item := range schemaListSet {
			schemaItem := dbbrain.SchemaItem{}
			schemaItemMap := item.(map[string]interface{})

			if v, ok := schemaItemMap["schema"]; ok {
				schemaItem.Schema = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &schemaItem)
		}
		paramMap["schema_list"] = tmpSet
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var rows []*dbbrain.SlowLogTopSqlItem

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainSlowLogTopSqlsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		rows = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(rows))

	if rows != nil {
		for _, slowLogTopSqlItem := range rows {
			slowLogTopSqlItemMap := map[string]interface{}{}

			if slowLogTopSqlItem.LockTime != nil {
				slowLogTopSqlItemMap["lock_time"] = slowLogTopSqlItem.LockTime
			}

			if slowLogTopSqlItem.LockTimeMax != nil {
				slowLogTopSqlItemMap["lock_time_max"] = slowLogTopSqlItem.LockTimeMax
			}

			if slowLogTopSqlItem.LockTimeMin != nil {
				slowLogTopSqlItemMap["lock_time_min"] = slowLogTopSqlItem.LockTimeMin
			}

			if slowLogTopSqlItem.RowsExamined != nil {
				slowLogTopSqlItemMap["rows_examined"] = slowLogTopSqlItem.RowsExamined
			}

			if slowLogTopSqlItem.RowsExaminedMax != nil {
				slowLogTopSqlItemMap["rows_examined_max"] = slowLogTopSqlItem.RowsExaminedMax
			}

			if slowLogTopSqlItem.RowsExaminedMin != nil {
				slowLogTopSqlItemMap["rows_examined_min"] = slowLogTopSqlItem.RowsExaminedMin
			}

			if slowLogTopSqlItem.QueryTime != nil {
				slowLogTopSqlItemMap["query_time"] = slowLogTopSqlItem.QueryTime
			}

			if slowLogTopSqlItem.QueryTimeMax != nil {
				slowLogTopSqlItemMap["query_time_max"] = slowLogTopSqlItem.QueryTimeMax
			}

			if slowLogTopSqlItem.QueryTimeMin != nil {
				slowLogTopSqlItemMap["query_time_min"] = slowLogTopSqlItem.QueryTimeMin
			}

			if slowLogTopSqlItem.RowsSent != nil {
				slowLogTopSqlItemMap["rows_sent"] = slowLogTopSqlItem.RowsSent
			}

			if slowLogTopSqlItem.RowsSentMax != nil {
				slowLogTopSqlItemMap["rows_sent_max"] = slowLogTopSqlItem.RowsSentMax
			}

			if slowLogTopSqlItem.RowsSentMin != nil {
				slowLogTopSqlItemMap["rows_sent_min"] = slowLogTopSqlItem.RowsSentMin
			}

			if slowLogTopSqlItem.ExecTimes != nil {
				slowLogTopSqlItemMap["exec_times"] = slowLogTopSqlItem.ExecTimes
			}

			if slowLogTopSqlItem.SqlTemplate != nil {
				slowLogTopSqlItemMap["sql_template"] = slowLogTopSqlItem.SqlTemplate
			}

			if slowLogTopSqlItem.SqlText != nil {
				slowLogTopSqlItemMap["sql_text"] = slowLogTopSqlItem.SqlText
			}

			if slowLogTopSqlItem.Schema != nil {
				slowLogTopSqlItemMap["schema"] = slowLogTopSqlItem.Schema
			}

			if slowLogTopSqlItem.QueryTimeRatio != nil {
				slowLogTopSqlItemMap["query_time_ratio"] = slowLogTopSqlItem.QueryTimeRatio
			}

			if slowLogTopSqlItem.LockTimeRatio != nil {
				slowLogTopSqlItemMap["lock_time_ratio"] = slowLogTopSqlItem.LockTimeRatio
			}

			if slowLogTopSqlItem.RowsExaminedRatio != nil {
				slowLogTopSqlItemMap["rows_examined_ratio"] = slowLogTopSqlItem.RowsExaminedRatio
			}

			if slowLogTopSqlItem.RowsSentRatio != nil {
				slowLogTopSqlItemMap["rows_sent_ratio"] = slowLogTopSqlItem.RowsSentRatio
			}

			if slowLogTopSqlItem.QueryTimeAvg != nil {
				slowLogTopSqlItemMap["query_time_avg"] = slowLogTopSqlItem.QueryTimeAvg
			}

			if slowLogTopSqlItem.RowsSentAvg != nil {
				slowLogTopSqlItemMap["rows_sent_avg"] = slowLogTopSqlItem.RowsSentAvg
			}

			if slowLogTopSqlItem.LockTimeAvg != nil {
				slowLogTopSqlItemMap["lock_time_avg"] = slowLogTopSqlItem.LockTimeAvg
			}

			if slowLogTopSqlItem.RowsExaminedAvg != nil {
				slowLogTopSqlItemMap["rows_examined_avg"] = slowLogTopSqlItem.RowsExaminedAvg
			}

			if slowLogTopSqlItem.Md5 != nil {
				slowLogTopSqlItemMap["md5"] = slowLogTopSqlItem.Md5
			}

			tmpList = append(tmpList, slowLogTopSqlItemMap)
		}

		_ = d.Set("rows", tmpList)
	}

	d.SetId(id)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
