package mps

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMpsManageTaskOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsManageTaskOperationCreate,
		Read:   resourceTencentCloudMpsManageTaskOperationRead,
		Delete: resourceTencentCloudMpsManageTaskOperationDelete,
		Schema: map[string]*schema.Schema{
			"operation_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "操作类型 有效 值:`Abort`: 任务 termination. Notice: 如果 任务 类型 是 live 流 processing (LiveStreamProcessTask)，tasks whose 任务 状态 是 `WAITING` 或 `PROCESSING` 可以 是 terminated.For other 任务 types，仅 tasks whose 任务 状态 是 `WAITING` 可以 是 terminated。",
			},

			"task_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Video processing 任务 ID。",
			},
		},
	}
}

func resourceTencentCloudMpsManageTaskOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_manage_task_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = mps.NewManageTaskRequest()
		taskId  string
	)
	if v, ok := d.GetOk("operation_type"); ok {
		request.OperationType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_id"); ok {
		request.TaskId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().ManageTask(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mps manageTaskOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(taskId)

	return resourceTencentCloudMpsManageTaskOperationRead(d, meta)
}

func resourceTencentCloudMpsManageTaskOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_manage_task_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMpsManageTaskOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_manage_task_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
