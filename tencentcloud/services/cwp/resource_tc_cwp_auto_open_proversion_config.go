package cwp

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cwpv20180228 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cwp/v20180228"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCwpAutoOpenProversionConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCwpAutoOpenProversionConfigCreate,
		Read:   resourceTencentCloudCwpAutoOpenProversionConfigRead,
		Update: resourceTencentCloudCwpAutoOpenProversionConfigUpdate,
		Delete: resourceTencentCloudCwpAutoOpenProversionConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Set auto-activation 状态\n<li>CLOSE: 关闭</li>\n<li>OPEN: 在</li>。",
			},

			"protect_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Enhanced Protection 模式 PROVERSION_POSTPAY Professional Edition - Pay-作为-您-go PROVERSION_PREPAY Professional Edition - Annual/Monthly Subscription FLAGSHIP_PREPAY Flagship Edition - Annual/Monthly Subscription。",
			},

			"auto_repurchase_switch": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Automatic purchase/expansion authorization switch，1 通过 默认值，0 对于 OFF，1 对于 ON。",
			},

			"auto_repurchase_renew_switch": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Auto-renewal 或 不 对于 auto-purchased orders，0 通过 默认值，0 对于 OFF，1 对于 ON。",
			},

			"repurchase_renew_switch": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "是否manually purchased 顺序 是 automatically renewed (默认为 0). 0 - 关闭; 1 -在。",
			},

			"auto_bind_rasp_switch": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Newly added machines 将 是 automatically bound 到 Rasp. 0: 已禁用，1: 已启用",
			},

			"auto_open_rasp_switch": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Newly added machines 将 have automatic Raspberry Pi protection 已启用 通过 默认值. (0: 已禁用，1: 已启用)。",
			},

			"auto_downgrade_switch": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Automatic scaling switch: 0 对于 关闭，1 对于 在。",
			},
		},
	}
}

func resourceTencentCloudCwpAutoOpenProversionConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cwp_auto_open_proversion_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	d.SetId(helper.BuildToken())

	return resourceTencentCloudCwpAutoOpenProversionConfigUpdate(d, meta)
}

func resourceTencentCloudCwpAutoOpenProversionConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cwp_auto_open_proversion_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CwpService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	respData, err := service.DescribeCwpAutoOpenProversionConfigById(ctx)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cwp_auto_open_proversion_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.AutoOpenStatus != nil {
		if *respData.AutoOpenStatus == true {
			_ = d.Set("status", "OPEN")
		} else {
			_ = d.Set("status", "CLOSE")
		}
	}

	if respData.ProtectType != nil {
		_ = d.Set("protect_type", respData.ProtectType)
	}

	if respData.AutoRepurchaseSwitch != nil {
		if *respData.AutoRepurchaseSwitch == true {
			_ = d.Set("auto_repurchase_switch", 1)
		} else {
			_ = d.Set("auto_repurchase_switch", 0)
		}
	}

	if respData.AutoRepurchaseRenewSwitch != nil {
		if *respData.AutoRepurchaseRenewSwitch == true {
			_ = d.Set("auto_repurchase_renew_switch", 1)
		} else {
			_ = d.Set("auto_repurchase_renew_switch", 0)
		}
	}

	if respData.RepurchaseRenewSwitch != nil {
		if *respData.RepurchaseRenewSwitch == true {
			_ = d.Set("repurchase_renew_switch", 1)
		} else {
			_ = d.Set("repurchase_renew_switch", 0)
		}
	}

	if respData.AutoBindRaspSwitch != nil {
		if *respData.AutoBindRaspSwitch == true {
			_ = d.Set("auto_bind_rasp_switch", 1)
		} else {
			_ = d.Set("auto_bind_rasp_switch", 0)
		}
	}

	if respData.AutoOpenRaspSwitch != nil {
		if *respData.AutoOpenRaspSwitch == true {
			_ = d.Set("auto_open_rasp_switch", 1)
		} else {
			_ = d.Set("auto_open_rasp_switch", 0)
		}
	}

	if respData.AutoDowngradeSwitch != nil {
		if *respData.AutoDowngradeSwitch == true {
			_ = d.Set("auto_downgrade_switch", 1)
		} else {
			_ = d.Set("auto_downgrade_switch", 0)
		}
	}

	return nil
}

func resourceTencentCloudCwpAutoOpenProversionConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cwp_auto_open_proversion_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = cwpv20180228.NewModifyAutoOpenProVersionConfigRequest()
	)

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protect_type"); ok {
		request.ProtectType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auto_repurchase_switch"); ok {
		request.AutoRepurchaseSwitch = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_repurchase_renew_switch"); ok {
		request.AutoRepurchaseRenewSwitch = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("repurchase_renew_switch"); ok {
		request.RepurchaseRenewSwitch = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_bind_rasp_switch"); ok {
		request.AutoBindRaspSwitch = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_open_rasp_switch"); ok {
		request.AutoOpenRaspSwitch = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_downgrade_switch"); ok {
		request.AutoDowngradeSwitch = helper.IntUint64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCwpClient().ModifyAutoOpenProVersionConfigWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update cwp auto open proversion config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudCwpAutoOpenProversionConfigRead(d, meta)
}

func resourceTencentCloudCwpAutoOpenProversionConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cwp_auto_open_proversion_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
