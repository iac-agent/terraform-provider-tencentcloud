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

func ResourceTencentCloudMysqlSwitchMasterSlaveOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlSwitchMasterSlaveOperationCreate,
		Read:   resourceTencentCloudMysqlSwitchMasterSlaveOperationRead,
		Delete: resourceTencentCloudMysqlSwitchMasterSlaveOperationDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "实例 ID。",
			},

			"dst_slave": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "目标实例。可能的值： `first` - 第一个备用； “第二”- 第二次待机。默认值为“first”，仅多可用区实例支持将其设置为“second”。",
			},

			"force_switch": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "是否强制切换。默认值为 False。请注意，如果将强制开关设置为True，则实例存在数据丢失的风险，请谨慎使用。",
			},

			"wait_switch": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "是否在时间窗口内切换。默认为False，即在时间窗口内不切换。注意，如果ForceSwitch参数设置为True，则该参数不会生效。",
			},
		},
	}
}

func resourceTencentCloudMysqlSwitchMasterSlaveOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_switch_master_slave_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request    = mysql.NewSwitchDBInstanceMasterSlaveRequest()
		response   = mysql.NewSwitchDBInstanceMasterSlaveResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_slave"); ok {
		request.DstSlave = helper.String(v.(string))
	}

	if v, _ := d.GetOk("force_switch"); v != nil {
		request.ForceSwitch = helper.Bool(v.(bool))
	}

	if v, _ := d.GetOk("wait_switch"); v != nil {
		request.WaitSwitch = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().SwitchDBInstanceMasterSlave(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mysql switchMasterSlaveOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

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
			return resource.RetryableError(fmt.Errorf("%s operate mysql switchMasterSlaveOperation status is %s", instanceId, taskStatus))
		}
		err = fmt.Errorf("%s operate mysql switchMasterSlaveOperation status is %s,we won't wait for it finish ,it show message:%s", instanceId, taskStatus, message)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate mysql switchMasterSlaveOperation fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudMysqlSwitchMasterSlaveOperationRead(d, meta)
}

func resourceTencentCloudMysqlSwitchMasterSlaveOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_switch_master_slave_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMysqlSwitchMasterSlaveOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_switch_master_slave_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
