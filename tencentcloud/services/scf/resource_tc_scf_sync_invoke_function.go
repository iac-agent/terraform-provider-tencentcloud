package scf

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudScfSyncInvokeFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfSyncInvokeFunctionCreate,
		Read:   resourceTencentCloudScfSyncInvokeFunctionRead,
		Delete: resourceTencentCloudScfSyncInvokeFunctionDelete,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Function 名称",
			},

			"qualifier": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "版本 or alias of the function. It 默认为 $DEFAULT。",
			},

			"event": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Function running parameter，which is in the JSON 格式 Maximum parameter size is 6 MB. This field corresponds to event input parameter。",
			},

			"log_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Valid 值: None (default) or Tail. If the 值 is Tail，log in the response will contain the corresponding function execution log (up to 4KB)。",
			},

			"namespace": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Namespace. 默认为 used if it's left empty。",
			},

			"routing_key": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Traffic routing 配置 in json 格式，e.g.，{k:v}. Please note that both k and v must be strings. Up to 1024 bytes allowed。",
			},
		},
	}
}

func resourceTencentCloudScfSyncInvokeFunctionCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_sync_invoke_function.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request           = scf.NewInvokeFunctionRequest()
		response          = scf.NewInvokeFunctionResponse()
		functionRequestId string
	)
	if v, ok := d.GetOk("function_name"); ok {
		request.FunctionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("qualifier"); ok {
		request.Qualifier = helper.String(v.(string))
	}

	if v, ok := d.GetOk("event"); ok {
		request.Event = helper.String(v.(string))
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
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseScfClient().InvokeFunction(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate scf syncInvokeFunction failed, reason:%+v", logId, err)
		return err
	}

	functionRequestId = *response.Response.Result.FunctionRequestId
	d.SetId(functionRequestId)

	return resourceTencentCloudScfSyncInvokeFunctionRead(d, meta)
}

func resourceTencentCloudScfSyncInvokeFunctionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_sync_invoke_function.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudScfSyncInvokeFunctionDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_scf_sync_invoke_function.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
