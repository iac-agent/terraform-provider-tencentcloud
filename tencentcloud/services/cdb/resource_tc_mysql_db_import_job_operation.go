package cdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlDbImportJobOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlDbImportJobOperationCreate,
		Read:   resourceTencentCloudMysqlDbImportJobOperationRead,
		Delete: resourceTencentCloudMysqlDbImportJobOperationDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "实例ID与云数据库控制台页面显示的实例ID一致，格式为：cdb-c1nl9rpv。",
			},

			"user": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "云数据库的用户名。",
			},

			"file_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "文件名。该文件是指用户上传到腾讯云的文件，仅支持.sql文件。",
			},

			"password": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "云数据库实例用户帐号的密码。",
			},

			"db_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "导入的目标数据库名称，如果不传则表示没有指定数据库。",
			},

			"cos_url": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "导入的目标数据库名称，如果不传则表示没有指定数据库。",
			},

			"async_request_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "异步任务的请求ID。",
			},
		},
	}
}

func resourceTencentCloudMysqlDbImportJobOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_db_import_job_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request    = mysql.NewCreateDBImportJobRequest()
		response   = mysql.NewCreateDBImportJobResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user"); ok {
		request.User = helper.String(v.(string))
	}

	if v, ok := d.GetOk("file_name"); ok {
		request.FileName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_name"); ok {
		request.DbName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cos_url"); ok {
		request.CosUrl = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().CreateDBImportJob(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mysql dbImportJob failed, reason:%+v", logId, err)
		return err
	}

	asyncRequestId := *response.Response.AsyncRequestId
	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		taskStatus, message, err := service.DescribeAsyncRequestInfo(ctx, asyncRequestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
			return nil
		}
		if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("%s create dbImportJob  status is %s", instanceId, taskStatus))
		}
		err = fmt.Errorf("%s create dbImportJob status is %s,we won't wait for it finish ,it show message:%s", instanceId, taskStatus, message)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dbImportJob fail, reason:%s\n ", logId, err.Error())
		return err
	}

	d.SetId(instanceId + tccommon.FILED_SP + asyncRequestId)

	return resourceTencentCloudMysqlDbImportJobOperationRead(d, meta)
}

func resourceTencentCloudMysqlDbImportJobOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_db_import_job_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	asyncRequestId := idSplit[1]

	dbImportJob, err := service.DescribeMysqlDbImportJobById(ctx, instanceId, asyncRequestId)
	if err != nil {
		return err
	}

	if dbImportJob == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlDbImportJob` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if dbImportJob.InstanceId != nil {
		_ = d.Set("instance_id", dbImportJob.InstanceId)
	}

	if dbImportJob.FileName != nil {
		_ = d.Set("file_name", dbImportJob.FileName)
	}

	if dbImportJob.DbName != nil {
		_ = d.Set("db_name", dbImportJob.DbName)
	}

	if dbImportJob.AsyncRequestId != nil {
		_ = d.Set("async_request_id", dbImportJob.AsyncRequestId)
	}

	return nil
}

func resourceTencentCloudMysqlDbImportJobOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_db_import_job_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
