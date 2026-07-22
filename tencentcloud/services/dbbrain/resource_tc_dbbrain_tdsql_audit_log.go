package dbbrain

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDbbrainTdsqlAuditLog() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDbbrainTdsqlAuditLogCreate,
		Read:   resourceTencentCloudDbbrainTdsqlAuditLogRead,
		Delete: resourceTencentCloudDbbrainTdsqlAuditLogDelete,
		Schema: map[string]*schema.Schema{
			"product": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Service product 类型, 支持 值 include: dcdb - 云 数据库 Tdsql, mariadb - 云 数据库 MariaDB 对于 MariaDB..",
			},

			"node_request_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Consistent 使用 Product. For 示例: dcdb, mariadb.",
			},

			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "实例 ID.",
			},

			"start_time": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Start 时间, such 作为 `2019-09-10 12:13:14`.",
			},

			"end_time": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Deadline 时间, such 作为 `2019-09-11 10:13:14`.",
			},

			"filter": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "过滤器 conditions. Logs 可以 是 filtered according 到 过滤器 conditions 集合.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							ForceNew:    true,
							Description: "Client Address.",
						},
						"db_name": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							ForceNew:    true,
							Description: "Database 名称.",
						},
						"user": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							ForceNew:    true,
							Description: "Username.",
						},
						"sent_rows": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "Return 数量 的 rows. It 表示 到 过滤器 audit 日志 使用 数量 的 返回 rows greater 比 此 值.",
						},
						"affect_rows": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "Number 的 affected rows. Indicates filtering audit logs whose affected rows 是 greater 比 此 值.",
						},
						"exec_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "Execution 时间. 单位 是: us. It 表示 到 过滤器 audit logs whose execution 时间 是 greater 比 此 值.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDbbrainTdsqlAuditLogCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbbrain_tdsql_audit_log.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request        = dbbrain.NewCreateAuditLogFileRequest()
		response       = dbbrain.NewCreateAuditLogFileResponse()
		asyncRequestId string
		instanceId     string
		product        string
	)
	if v, ok := d.GetOk("product"); ok {
		request.Product = helper.String(v.(string))
		product = v.(string)
	}

	if v, ok := d.GetOk("node_request_type"); ok {
		request.NodeRequestType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("start_time"); ok {
		request.StartTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "filter"); ok {
		auditLogFilter := dbbrain.AuditLogFilter{}
		if v, ok := dMap["host"]; ok {
			hostSet := v.(*schema.Set).List()
			for i := range hostSet {
				host := hostSet[i].(string)
				auditLogFilter.Host = append(auditLogFilter.Host, &host)
			}
		}
		if v, ok := dMap["db_name"]; ok {
			dBNameSet := v.(*schema.Set).List()
			for i := range dBNameSet {
				dBName := dBNameSet[i].(string)
				auditLogFilter.DBName = append(auditLogFilter.DBName, &dBName)
			}
		}
		if v, ok := dMap["user"]; ok {
			userSet := v.(*schema.Set).List()
			for i := range userSet {
				user := userSet[i].(string)
				auditLogFilter.User = append(auditLogFilter.User, &user)
			}
		}
		if v, ok := dMap["sent_rows"]; ok {
			auditLogFilter.SentRows = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["affect_rows"]; ok {
			auditLogFilter.AffectRows = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["exec_time"]; ok {
			auditLogFilter.ExecTime = helper.IntInt64(v.(int))
		}
		request.Filter = &auditLogFilter
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDbbrainClient().CreateAuditLogFile(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dbbrain tdsqlAuditLog failed, reason:%+v", logId, err)
		return err
	}

	asyncRequestId = helper.Int64ToStr(*response.Response.AsyncRequestId)
	d.SetId(strings.Join([]string{asyncRequestId, instanceId, product}, tccommon.FILED_SP))

	return resourceTencentCloudDbbrainTdsqlAuditLogRead(d, meta)
}

func resourceTencentCloudDbbrainTdsqlAuditLogRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbbrain_tdsql_audit_log.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	asyncRequestId := idSplit[0]
	instanceId := idSplit[1]
	product := idSplit[2]

	tdsqlAuditLogs, err := service.DescribeDbbrainTdsqlAuditLogById(ctx, &asyncRequestId, instanceId, product)
	if err != nil {
		return err
	}

	if len(tdsqlAuditLogs) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `DbbrainTdsqlAuditLog` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	// tdsqlAuditLog := tdsqlAuditLogs[0]

	// if tdsqlAuditLog.CreateTime != nil {
	// 	_ = d.Set("start_time", tdsqlAuditLog.CreateTime)
	// }

	// if tdsqlAuditLog.FinishTime != nil {
	// 	_ = d.Set("end_time", tdsqlAuditLog.FinishTime)
	// }

	return nil
}

func resourceTencentCloudDbbrainTdsqlAuditLogDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dbbrain_tdsql_audit_log.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	asyncRequestId := idSplit[0]
	instanceId := idSplit[1]
	product := idSplit[2]

	if err := service.DeleteDbbrainTdsqlAuditLogById(ctx, asyncRequestId, instanceId, product); err != nil {
		return err
	}

	return nil
}
