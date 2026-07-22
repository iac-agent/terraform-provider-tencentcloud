package dlc

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcRenewDataEngineOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcRenewDataEngineCreate,
		Read:   resourceTencentCloudDlcRenewDataEngineRead,
		Delete: resourceTencentCloudDlcRenewDataEngineDelete,
		Schema: map[string]*schema.Schema{
			"data_engine_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CU queue 名称",
			},

			"time_span": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Renewal 周期 在 months，其中 是 在 least 一个 month。",
			},

			"pay_mode": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Payment 类型 It 是 1 通过 默认值 和 是 prepaid。",
			},

			"time_unit": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Unit. It 是 m 通过 默认值，和 仅 m 可以 是 filled 在。",
			},

			"renew_flag": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Auto-renewal flag: 0 表示 initial 状态，和 there 是 无 automatic renewal 通过 默认值. 如果 用户 has privilege 到 retain services 使用 prepayment，there 将 是 automatic renewal. 1 表示 该 there 是 automatic renewal. 2 表示 该 there 是 surely 无 automatic renewal. 如果 它 是 不 指定， 参数 是 0 通过 默认值。",
			},
		},
	}
}

func resourceTencentCloudDlcRenewDataEngineCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_renew_data_engine_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		request        = dlc.NewRenewDataEngineRequest()
		dataEngineName string
	)

	if v, ok := d.GetOk("data_engine_name"); ok {
		dataEngineName = v.(string)
		request.DataEngineName = helper.String(v.(string))
	}

	if v, _ := d.GetOkExists("time_span"); v != nil {
		request.TimeSpan = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOkExists("pay_mode"); v != nil {
		request.PayMode = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("time_unit"); ok {
		request.TimeUnit = helper.String(v.(string))
	}

	if v, _ := d.GetOkExists("renew_flag"); v != nil {
		request.RenewFlag = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().RenewDataEngine(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate dlc renewDataEngine failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(dataEngineName)
	return resourceTencentCloudDlcRenewDataEngineRead(d, meta)
}

func resourceTencentCloudDlcRenewDataEngineRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_renew_data_engine_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDlcRenewDataEngineDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_renew_data_engine_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
