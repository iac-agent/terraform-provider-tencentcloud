package scf

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudScfInvokeFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfInvokeFunctionCreate,
		Read:   resourceTencentCloudScfInvokeFunctionRead,
		Delete: resourceTencentCloudScfInvokeFunctionDelete,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Function 名称",
			},

			"invocation_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Fill 在 RequestResponse 对于 synchronized invocations (默认值 和 recommended) 和 Event 对于 asychronized invocations. 注意 该 对于 synchronized invocations， max 超时 周期 是 300s. Choose asychronized invocations 如果 必填 超时 周期 是 longer 比 300 秒. You 可以 also 使用 InvokeFunction 对于 synchronized invocations。",
			},

			"qualifier": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "版本 或 alias 的 triggered 函数. It 默认为 $LATEST。",
			},

			"client_context": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Function running 参数，其中 是 在 JSON 格式 最大 参数 大小 是 6 MB 对于 synchronized invocations 和 128KB 对于 asynchronized invocations. 此 字段 corresponds 到 事件 input 参数。",
			},

			"log_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Null 对于 async invocations。",
			},

			"namespace": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Namespace。",
			},

			"routing_key": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Traffic routing 配置 在 json 格式，e.g.，{k:v}. Please note 该 both k 和 v 必须 是 strings. Up 到 1024 bytes allowed。",
			},
		},
	}
}

func resourceTencentCloudScfInvokeFunctionCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_invoke_function.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = scf.NewInvokeRequest()
		response = scf.NewInvokeResponse()
	)
	if v, ok := d.GetOk("function_name"); ok {
		request.FunctionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("invocation_type"); ok {
		request.InvocationType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("qualifier"); ok {
		request.Qualifier = helper.String(v.(string))
	}

	if v, ok := d.GetOk("client_context"); ok {
		request.ClientContext = helper.String(v.(string))
	}

	if v, ok := d.GetOk("log_type"); ok {
		request.LogType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		request.Namespace = helper.String(v.(string))
	}

	if v, ok := d.GetOk("routing_key"); ok {
		request.RoutingKey = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseScfClient().Invoke(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate scf InvokeFunction failed, reason:%+v", logId, err)
		return err
	}

	functionRequestId := *response.Response.Result.FunctionRequestId

	d.SetId(functionRequestId)

	return resourceTencentCloudScfInvokeFunctionRead(d, meta)
}

func resourceTencentCloudScfInvokeFunctionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_invoke_function.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudScfInvokeFunctionDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_invoke_function.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
