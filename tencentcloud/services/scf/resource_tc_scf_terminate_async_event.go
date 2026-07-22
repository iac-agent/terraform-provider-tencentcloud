package scf

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudScfTerminateAsyncEvent() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfTerminateAsyncEventCreate,
		Read:   resourceTencentCloudScfTerminateAsyncEventRead,
		Delete: resourceTencentCloudScfTerminateAsyncEventDelete,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Function 名称",
			},

			"invoke_request_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Terminated invocation 请求 ID",
			},

			"namespace": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Namespace。",
			},

			"grace_shutdown": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "是否enable grace shutdown. 如果 它's true， SIGTERM signal 是 sent 到 指定 请求. See [Sending termination signal](https://www.tencentcloud.com/document/product/583/63969?from_cn_redirect=1#.E5.8F.91.E9.80.81.E7.BB.88.E6.AD.A2.E4.BF.A1.E5.8F.B7]. It's 集合 到 false 通过 默认值。",
			},
		},
	}
}

func resourceTencentCloudScfTerminateAsyncEventCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_terminate_async_event.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request         = scf.NewTerminateAsyncEventRequest()
		invokeRequestId string
	)
	if v, ok := d.GetOk("function_name"); ok {
		request.FunctionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("invoke_request_id"); ok {
		invokeRequestId = v.(string)
		request.InvokeRequestId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		request.Namespace = helper.String(v.(string))
	}

	if v, _ := d.GetOk("grace_shutdown"); v != nil {
		request.GraceShutdown = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseScfClient().TerminateAsyncEvent(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate scf terminateAsyncEvent failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(invokeRequestId)

	return resourceTencentCloudScfTerminateAsyncEventRead(d, meta)
}

func resourceTencentCloudScfTerminateAsyncEventRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_terminate_async_event.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudScfTerminateAsyncEventDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_terminate_async_event.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
