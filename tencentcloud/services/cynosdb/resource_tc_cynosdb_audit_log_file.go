package cynosdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudCynosdbAuditLogFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbAuditLogFileCreate,
		Read:   resourceTencentCloudCynosdbAuditLogFileRead,
		Delete: resourceTencentCloudCynosdbAuditLogFileDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "实例的ID。",
			},

			"start_time": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "开始时间。",
			},

			"end_time": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "结束时间。",
			},

			"order": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "排序方式。支持的值为：“ASC”- 升序，“DESC”- 降序。",
			},

			"order_by": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "排序字段。支持的值为：\n`时间戳` - 时间戳\n`affectRows` - 受影响的行\n`execTime` - 执行时间。",
			},

			"filter": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "过滤条件。可以根据设置的过滤条件对日志进行过滤。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "客户端主机。",
						},
						"user": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "用户名。",
						},
						"db_name": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "数据库的名称。",
						},
						"table_name": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "表的名称。",
						},
						"policy_name": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "审核策略的名称。",
						},
						"sql": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SQL 语句。支持模糊匹配。",
						},
						"sql_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SQL 类型。目前支持：SELECT、INSERT、UPDATE、DELETE、CREATE、DROP、ALTER、SET、REPLACE、EXECUTE。",
						},
						"exec_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "执行时间。单位是：毫秒。表示过滤执行时间大于该值的审计日志。",
						},
						"affect_rows": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "影响行数。表示过滤影响行数大于该值的审计日志。",
						},
						"sql_types": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "SQL 类型。支持多种类型同时查询。目前支持：SELECT、INSERT、UPDATE、DELETE、CREATE、DROP、ALTER、SET、REPLACE、EXECUTE。",
						},
						"sqls": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "SQL 语句。支持传递多个sql语句。",
						},
						"sent_rows": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "返回行数。",
						},
						"thread_id": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "线程的ID。",
						},
					},
				},
			},
			// computed
			"file_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "审核日志文件名。",
			},
			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "审核日志文件创建时间。格式为2019-03-20 17:09:13。",
			},
			"file_size": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "文件大小，单位为KB。",
			},
			"download_url": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "审计日志的下载地址。",
			},
			"err_msg": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "错误信息。",
			},
		},
	}
}

func resourceTencentCloudCynosdbAuditLogFileCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_audit_log_file.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = cynosdb.NewCreateAuditLogFileRequest()
		response   = cynosdb.NewCreateAuditLogFileResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOk("start_time"); ok {
		request.StartTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order"); ok {
		request.Order = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		request.OrderBy = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "filter"); ok {
		auditLogFilter := cynosdb.AuditLogFilter{}
		if v, ok := dMap["host"]; ok {
			hostSet := v.(*schema.Set).List()
			for i := range hostSet {
				host := hostSet[i].(string)
				auditLogFilter.Host = append(auditLogFilter.Host, &host)
			}
		}
		if v, ok := dMap["user"]; ok {
			userSet := v.(*schema.Set).List()
			for i := range userSet {
				user := userSet[i].(string)
				auditLogFilter.User = append(auditLogFilter.User, &user)
			}
		}
		if v, ok := dMap["db_name"]; ok {
			dBNameSet := v.(*schema.Set).List()
			for i := range dBNameSet {
				dBName := dBNameSet[i].(string)
				auditLogFilter.DBName = append(auditLogFilter.DBName, &dBName)
			}
		}
		if v, ok := dMap["table_name"]; ok {
			tableNameSet := v.(*schema.Set).List()
			for i := range tableNameSet {
				tableName := tableNameSet[i].(string)
				auditLogFilter.TableName = append(auditLogFilter.TableName, &tableName)
			}
		}
		if v, ok := dMap["policy_name"]; ok {
			policyNameSet := v.(*schema.Set).List()
			for i := range policyNameSet {
				policyName := policyNameSet[i].(string)
				auditLogFilter.PolicyName = append(auditLogFilter.PolicyName, &policyName)
			}
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
			for i := range sqlTypesSet {
				sqlTypes := sqlTypesSet[i].(string)
				auditLogFilter.SqlTypes = append(auditLogFilter.SqlTypes, &sqlTypes)
			}
		}
		if v, ok := dMap["sqls"]; ok {
			sqlsSet := v.(*schema.Set).List()
			for i := range sqlsSet {
				sqls := sqlsSet[i].(string)
				auditLogFilter.Sqls = append(auditLogFilter.Sqls, &sqls)
			}
		}
		if v, ok := dMap["sent_rows"]; ok {
			auditLogFilter.SentRows = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["thread_id"]; ok {
			threadIdSet := v.(*schema.Set).List()
			for i := range threadIdSet {
				threadId := threadIdSet[i].(string)
				auditLogFilter.ThreadId = append(auditLogFilter.ThreadId, &threadId)
			}
		}
		request.Filter = &auditLogFilter
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().CreateAuditLogFile(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb auditLogFile failed, reason:%+v", logId, err)
		return err
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		request := cynosdb.NewDescribeAuditLogFilesRequest()
		request.InstanceId = helper.String(instanceId)
		request.FileName = response.Response.FileName
		ratelimit.Check(request.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().DescribeAuditLogFiles(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if len(result.Response.Items) > 0 && *result.Response.Items[0].Status == "success" {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("%s not ready", *response.Response.FileName))
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb auditLogFile failed, reason:%+v", logId, err)
		return err
	}

	auditLogFileId := strings.Join([]string{instanceId, *response.Response.FileName}, tccommon.FILED_SP)
	d.SetId(auditLogFileId)

	return resourceTencentCloudCynosdbAuditLogFileRead(d, meta)
}

func resourceTencentCloudCynosdbAuditLogFileRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_audit_log_file.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	fileName := idSplit[1]

	auditLogFile, err := service.DescribeCynosdbAuditLogFileById(ctx, instanceId, fileName)
	if err != nil {
		return err
	}

	if auditLogFile == nil {
		d.SetId("")
		return fmt.Errorf("resource `CynosdbAuditLogFile` %s does not exist", d.Id())
	}

	_ = d.Set("file_name", *auditLogFile.FileName)
	_ = d.Set("create_time", *auditLogFile.CreateTime)
	_ = d.Set("file_size", *auditLogFile.FileSize)
	_ = d.Set("download_url", *auditLogFile.DownloadUrl)
	_ = d.Set("err_msg", *auditLogFile.ErrMsg)

	return nil
}

func resourceTencentCloudCynosdbAuditLogFileDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_audit_log_file.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	fileName := idSplit[1]

	if err := service.DeleteCynosdbAuditLogFileById(ctx, instanceId, fileName); err != nil {
		return err
	}

	return nil
}
