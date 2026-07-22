package cdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMysqlSlowLogData() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlSlowLogDataRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID。",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "开始时间戳。例如 1585142640。",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "结束时间戳。例如 1585142640。",
			},

			"user_hosts": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "客户端主机列表。",
			},

			"user_names": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "客户端用户名列表。",
			},

			"data_bases": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "访问的数据库列表。",
			},

			"sort_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "排序字段。目前支持：时间戳、QueryTime、LockTime、RowsExamined、RowsSent。",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "按升序或降序排序。目前支持：ASC、DESC。",
			},

			"inst_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "仅当实例为主实例或灾备实例时有效，可选值：slave，表示拉取从机的日志。",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "查询记录。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "SQL执行时间。",
						},
						"query_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Sql 执行时间（秒）。",
						},
						"sql_text": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SQL语句。",
						},
						"user_host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "客户地址。",
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "用户名。",
						},
						"database": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "数据库名称。",
						},
						"lock_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "锁定持续时间（秒）。",
						},
						"rows_examined": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "要扫描的行数。",
						},
						"rows_sent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "结果集中的行数。",
						},
						"sql_template": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SQL 模板。",
						},
						"md5": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sql语句的md5。",
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

func dataSourceTencentCloudMysqlSlowLogDataRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_slow_log_data.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var instanceId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("start_time"); v != nil {
		paramMap["StartTime"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("end_time"); v != nil {
		paramMap["EndTime"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("user_hosts"); ok {
		userHostsSet := v.(*schema.Set).List()
		paramMap["UserHosts"] = helper.InterfacesStringsPoint(userHostsSet)
	}

	if v, ok := d.GetOk("user_names"); ok {
		userNamesSet := v.(*schema.Set).List()
		paramMap["UserNames"] = helper.InterfacesStringsPoint(userNamesSet)
	}

	if v, ok := d.GetOk("data_bases"); ok {
		dataBasesSet := v.(*schema.Set).List()
		paramMap["DataBases"] = helper.InterfacesStringsPoint(dataBasesSet)
	}

	if v, ok := d.GetOk("sort_by"); ok {
		paramMap["SortBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("inst_type"); ok {
		paramMap["InstType"] = helper.String(v.(string))
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var items []*cdb.SlowLogItem
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlSlowLogDataByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		items = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(items))
	if items != nil {
		for _, slowLogItem := range items {
			slowLogItemMap := map[string]interface{}{}

			if slowLogItem.Timestamp != nil {
				slowLogItemMap["timestamp"] = slowLogItem.Timestamp
			}

			if slowLogItem.QueryTime != nil {
				slowLogItemMap["query_time"] = slowLogItem.QueryTime
			}

			if slowLogItem.SqlText != nil {
				slowLogItemMap["sql_text"] = slowLogItem.SqlText
			}

			if slowLogItem.UserHost != nil {
				slowLogItemMap["user_host"] = slowLogItem.UserHost
			}

			if slowLogItem.UserName != nil {
				slowLogItemMap["user_name"] = slowLogItem.UserName
			}

			if slowLogItem.Database != nil {
				slowLogItemMap["database"] = slowLogItem.Database
			}

			if slowLogItem.LockTime != nil {
				slowLogItemMap["lock_time"] = slowLogItem.LockTime
			}

			if slowLogItem.RowsExamined != nil {
				slowLogItemMap["rows_examined"] = slowLogItem.RowsExamined
			}

			if slowLogItem.RowsSent != nil {
				slowLogItemMap["rows_sent"] = slowLogItem.RowsSent
			}

			if slowLogItem.SqlTemplate != nil {
				slowLogItemMap["sql_template"] = slowLogItem.SqlTemplate
			}

			if slowLogItem.Md5 != nil {
				slowLogItemMap["md5"] = slowLogItem.Md5
			}

			tmpList = append(tmpList, slowLogItemMap)
		}

		_ = d.Set("items", tmpList)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
