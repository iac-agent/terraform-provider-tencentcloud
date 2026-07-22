package cdb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlRenewDbInstanceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlRenewDbInstanceOperationCreate,
		Read:   resourceTencentCloudMysqlRenewDbInstanceOperationRead,
		Delete: resourceTencentCloudMysqlRenewDbInstanceOperationDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "需要续费的实例ID，格式为：cdb-c1nl9rpv，与云数据库控制台页面显示的实例ID相同，可以使用【查询实例列表】(https://云.tencent.com/document/api/236/15872)。",
			},

			"time_span": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "续订时长，单位：月，可选值包括[1,2,3,4,5,6,7,8,9,10,11,12,24,36]。",
			},

			"modify_pay_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "如果您需要将按量付费实例续费为包年包月实例，则需要将该入参指定为“PREPAID”。",
			},

			"deal_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "交易 ID。",
			},

			"deadline_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "实例过期时间。",
			},
		},
	}
}

func resourceTencentCloudMysqlRenewDbInstanceOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_renew_db_instance_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = mysql.NewRenewDBInstanceRequest()
		response   = mysql.NewRenewDBInstanceResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("time_span"); v != nil {
		request.TimeSpan = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("modify_pay_type"); ok {
		request.ModifyPayType = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().RenewDBInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mysql renewDbInstanceOperation failed, reason:%+v", logId, err)
		return err
	}

	dealId := *response.Response.DealId
	d.SetId(instanceId)
	_ = d.Set("deal_id", dealId)

	return resourceTencentCloudMysqlRenewDbInstanceOperationRead(d, meta)
}

func resourceTencentCloudMysqlRenewDbInstanceOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_renew_db_instance_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	mysqlInfo, errRet := service.DescribeDBInstanceById(ctx, d.Id())
	if errRet != nil {
		return fmt.Errorf("Describe mysql instance fails, reaseon %s", errRet.Error())
	}

	if mysqlInfo == nil {
		d.SetId("")
		return nil
	}

	if mysqlInfo.DeadlineTime != nil {
		_ = d.Set("deadline_time", mysqlInfo.DeadlineTime)
	}

	return nil
}

func resourceTencentCloudMysqlRenewDbInstanceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_renew_db_instance_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
