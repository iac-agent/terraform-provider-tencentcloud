package cvm

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCvmRenewHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmRenewHostCreate,
		Read:   resourceTencentCloudCvmRenewHostRead,
		Delete: resourceTencentCloudCvmRenewHostDelete,
		Schema: map[string]*schema.Schema{
			"host_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CDH 实例 ID。",
			},

			"host_charge_prepaid": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Prepaid 模式，该 是，yearly 和 monthly subscription related 参数 settings. Through 此 参数，您 可以 指定attributes such 作为 purchase 时长 的 Subscription 实例 和 是否set automatic renewal. 如果 payment 模式 的 指定 实例 是 prepaid，此 参数 必须 是 passed。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "时长 的 purchasing 实例，单位: month. 取值范围：1，2，3，4，5，6，7，8，9，10，11，12，24，36。",
						},
						"renew_flag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "自动续费标识 有效 值:&lt;br&gt;&lt;li&gt;NOTIFY_AND_AUTO_RENEW: notify upon expiration 和 renew automatically&lt;br&gt;&lt;li&gt;NOTIFY_AND_MANUAL_RENEW: notify upon expiration 但 do 不 renew automatically&lt;br&gt;&lt;li&gt;DISABLE_NOTIFY_AND_MANUAL_RENEW: neither notify upon expiration nor renew automatically&lt;br&gt;&lt;br&gt;默认值：NOTIFY_AND_AUTO_RENEW。如果 此 参数 是 指定 作为 NOTIFY_AND_AUTO_RENEW， 实例 将 是 automatically renewed 在 monthly basis 如果 账号 balance 是 sufficient。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCvmRenewHostCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_renew_host.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = cvm.NewRenewHostsRequest()
	)
	hostId := d.Get("host_id").(string)
	request.HostIds = []*string{&hostId}

	if dMap, ok := helper.InterfacesHeadMap(d, "host_charge_prepaid"); ok {
		chargePrepaid := cvm.ChargePrepaid{}
		if v, ok := dMap["period"]; ok {
			chargePrepaid.Period = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["renew_flag"]; ok {
			chargePrepaid.RenewFlag = helper.String(v.(string))
		}
		request.HostChargePrepaid = &chargePrepaid
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().RenewHosts(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm renewHost failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(hostId)

	return resourceTencentCloudCvmRenewHostRead(d, meta)
}

func resourceTencentCloudCvmRenewHostRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_renew_host.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmRenewHostDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_renew_host.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
