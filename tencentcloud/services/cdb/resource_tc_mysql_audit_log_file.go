package cdb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlAuditLogFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlAuditLogFileCreate,
		Read:   resourceTencentCloudMysqlAuditLogFileRead,
		Delete: resourceTencentCloudMysqlAuditLogFileDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "实例的ID。",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "开始时间。",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "结束时间。",
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "排序方式。支持的值为：“ASC”- 升序，“DESC”- 降序。",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "排序字段。支持的值包括：`时间戳` - 时间戳； `affectRows` - 受影响的行； `execTime` - 执行时间。",
			},

			"filter": {
				Optional:    true,
				Type:        schema.TypeList,
				ForceNew:    true,
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
							Description: "客户地址。",
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
							Description: "数据库名称。",
						},
						"table_name": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "表名。",
						},
						"policy_name": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "政策名称。",
						},
						"sql": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SQL 语句。支持模糊匹配。",
						},
						"sql_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SQL 类型。目前支持：选择、插入、更新、删除、创建、删除、更改、设置、替换、执行。",
						},
						"exec_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "执行时间。单位是：毫秒。表示过滤执行时间大于该值的审计日志。",
						},
						"affect_rows": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "影响行数。表示过滤受影响行数大于该值的审计日志。",
						},
						"sql_types": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "SQL 类型。支持多种类型同时查询。目前支持：选择、插入、更新、删除、创建、删除、更改、设置、替换、执行。",
						},
						"sqls": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "SQL 语句。支持传递多个sql语句。",
						},
					},
				},
			},

			"file_size": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "文件大小（KB）。",
			},

			"download_url": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "下载网址。",
			},
		},
	}
}

func resourceTencentCloudMysqlAuditLogFileCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_audit_log_file.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = mysql.NewCreateAuditLogFileRequest()
		response   = mysql.NewCreateAuditLogFileResponse()
		instanceId string
		fileName   string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
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
		auditLogFilter := mysql.AuditLogFilter{}
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
		request.Filter = &auditLogFilter
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().CreateAuditLogFile(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mysql auditLogFile failed, reason:%+v", logId, err)
		return err
	}

	fileName = *response.Response.FileName
	d.SetId(instanceId + tccommon.FILED_SP + fileName)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf(
		[]string{"creating"},
		[]string{"success", "failed"},
		1*tccommon.ReadRetryTimeout,
		time.Second,
		service.MysqlAuditLogFileStateRefreshFunc(instanceId, fileName, []string{}),
	)

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudMysqlAuditLogFileRead(d, meta)
}

func resourceTencentCloudMysqlAuditLogFileRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_audit_log_file.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	fileName := idSplit[1]

	auditLogFile, err := service.DescribeMysqlAuditLogFileById(ctx, instanceId, fileName)
	if err != nil {
		return err
	}

	if auditLogFile == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlAuditLogFile` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if auditLogFile.FileSize != nil {
		_ = d.Set("file_size", auditLogFile.FileSize)
	}

	if auditLogFile.DownloadUrl != nil {
		_ = d.Set("download_url", auditLogFile.DownloadUrl)
	}

	return nil
}

func resourceTencentCloudMysqlAuditLogFileDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_audit_log_file.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	fileName := idSplit[1]

	if err := service.DeleteMysqlAuditLogFileById(ctx, instanceId, fileName); err != nil {
		return err
	}

	return nil
}
