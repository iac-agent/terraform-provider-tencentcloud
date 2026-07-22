package cynosdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCynosdbAuditLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbAuditLogsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例ID。",
			},
			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "开始时间，格式：2017-07-12 10:29:20。",
			},
			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "结束时间格式为2017-07-12 10:29:20。",
			},
			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "排序方式。支持的值包括：ASC - 升序、DESC - 降序。",
			},
			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "对字段进行排序。支持的值包括：时间戳-时间戳； &amp;#39;效果行&amp;#39; - 影响行数； &amp;#39;执行时间&amp;#39; - 执行时间。",
			},
			"filter": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "过滤条件。您可以根据设置的过滤条件过滤日志。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "客户地址。",
						},
						"user": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "用户名。",
						},
						"db_name": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "数据库名称。",
						},
						"table_name": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "表名。",
						},
						"policy_name": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "审核策略名称。",
						},
						"sql": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SQL 语句。支持模糊匹配。",
						},
						"sql_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SQL 类型。目前支持：SELECT、Insert、UPDATE、DELETE、CREATE、DROP、ALT、SET、REPLACE、EXECUTE。",
						},
						"exec_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "执行时间。单位：毫秒表示过滤器执行时间大于该值的审核日志。",
						},
						"affect_rows": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "影响行数。指示过滤影响行数大于此值的审核日志。",
						},
						"sql_types": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "SQL 类型。支持多种类型同时查询。目前支持：SELECT、Insert、UPDATE、DELETE、CREATE、DROP、ALT、SET、REPLACE、EXECUTE。",
						},
						"sqls": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "SQL 语句。支持传递多个SQL语句。",
						},
						"sent_rows": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "返回行数。",
						},
						"thread_id": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "线程 ID。",
						},
					},
				},
			},
			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "审核日志详细信息。注意：该字段可能返回null，表示无法获取到有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"affect_rows": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "影响行数。",
						},
						"err_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "错误代码。",
						},
						"sql_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SQL 类型。",
						},
						"table_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "表名。",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例名称。",
						},
						"policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "审核策略名称。",
						},
						"db_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "数据库名称。",
						},
						"sql": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SQL 语句。",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "客户地址。",
						},
						"user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "用户名。",
						},
						"exec_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "执行时间。",
						},
						"timestamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "时间戳。",
						},
						"sent_rows": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "发送的行数。",
						},
						"thread_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "执行线程ID。",
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

func dataSourceTencentCloudCynosdbAuditLogsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cynosdb_audit_logs.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		items      []*cynosdb.AuditLog
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "filter"); ok {
		auditLogFilter := cynosdb.AuditLogFilter{}
		if v, ok := dMap["host"]; ok {
			hostSet := v.(*schema.Set).List()
			auditLogFilter.Host = helper.InterfacesStringsPoint(hostSet)
		}
		if v, ok := dMap["user"]; ok {
			userSet := v.(*schema.Set).List()
			auditLogFilter.User = helper.InterfacesStringsPoint(userSet)
		}
		if v, ok := dMap["db_name"]; ok {
			dBNameSet := v.(*schema.Set).List()
			auditLogFilter.DBName = helper.InterfacesStringsPoint(dBNameSet)
		}
		if v, ok := dMap["table_name"]; ok {
			tableNameSet := v.(*schema.Set).List()
			auditLogFilter.TableName = helper.InterfacesStringsPoint(tableNameSet)
		}
		if v, ok := dMap["policy_name"]; ok {
			policyNameSet := v.(*schema.Set).List()
			auditLogFilter.PolicyName = helper.InterfacesStringsPoint(policyNameSet)
		}
		if v, ok := dMap["sql"]; ok {
			auditLogFilter.Sql = helper.String(v.(string))
		}
		if v, ok := dMap["sql_type"]; ok {
			auditLogFilter.SqlType = helper.String(v.(string))
		}
		if v, ok := dMap["exec_time"]; ok {
			auditLogFilter.ExecTime = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["affect_rows"]; ok {
			auditLogFilter.AffectRows = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["sql_types"]; ok {
			sqlTypesSet := v.(*schema.Set).List()
			auditLogFilter.SqlTypes = helper.InterfacesStringsPoint(sqlTypesSet)
		}
		if v, ok := dMap["sqls"]; ok {
			sqlsSet := v.(*schema.Set).List()
			auditLogFilter.Sqls = helper.InterfacesStringsPoint(sqlsSet)
		}
		if v, ok := dMap["sent_rows"]; ok {
			auditLogFilter.SentRows = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["thread_id"]; ok {
			threadIdSet := v.(*schema.Set).List()
			auditLogFilter.ThreadId = helper.InterfacesStringsPoint(threadIdSet)
		}
		paramMap["filter"] = &auditLogFilter
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbAuditLogsByFilter(ctx, paramMap)
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
		for _, auditLog := range items {
			auditLogMap := map[string]interface{}{}

			if auditLog.AffectRows != nil {
				auditLogMap["affect_rows"] = auditLog.AffectRows
			}

			if auditLog.ErrCode != nil {
				auditLogMap["err_code"] = auditLog.ErrCode
			}

			if auditLog.SqlType != nil {
				auditLogMap["sql_type"] = auditLog.SqlType
			}

			if auditLog.TableName != nil {
				auditLogMap["table_name"] = auditLog.TableName
			}

			if auditLog.InstanceName != nil {
				auditLogMap["instance_name"] = auditLog.InstanceName
			}

			if auditLog.PolicyName != nil {
				auditLogMap["policy_name"] = auditLog.PolicyName
			}

			if auditLog.DBName != nil {
				auditLogMap["db_name"] = auditLog.DBName
			}

			if auditLog.Sql != nil {
				auditLogMap["sql"] = auditLog.Sql
			}

			if auditLog.Host != nil {
				auditLogMap["host"] = auditLog.Host
			}

			if auditLog.User != nil {
				auditLogMap["user"] = auditLog.User
			}

			if auditLog.ExecTime != nil {
				auditLogMap["exec_time"] = auditLog.ExecTime
			}

			if auditLog.Timestamp != nil {
				auditLogMap["timestamp"] = auditLog.Timestamp
			}

			if auditLog.SentRows != nil {
				auditLogMap["sent_rows"] = auditLog.SentRows
			}

			if auditLog.ThreadId != nil {
				auditLogMap["thread_id"] = auditLog.ThreadId
			}

			tmpList = append(tmpList, auditLogMap)
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
