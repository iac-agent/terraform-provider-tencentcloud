package lighthouse

import (
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudLighthouseRenewDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseRenewDiskCreate,
		Read:   resourceTencentCloudLighthouseRenewDiskRead,
		Delete: resourceTencentCloudLighthouseRenewDiskDelete,
		Schema: map[string]*schema.Schema{
			"disk_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "列表 磁盘 ID。",
			},

			"renew_disk_charge_prepaid": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Renew 云 hard 磁盘 subscription related 参数 settings。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Renewal 周期",
						},
						"renew_flag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Automatic renewal falg. 值:NOTIFY_AND_AUTO_RENEW: Notice expires 和 auto-renews.NOTIFY_AND_MANUAL_RENEW: Notification expires without automatic renewal，users need 到 manually renew.DISABLE_NOTIFY_AND_AUTO_RENEW: No automatic renewal 和 无 通知.默认值：NOTIFY_AND_MANUAL_RENEW. 如果 此 参数 是 指定 作为 NOTIFY_AND_AUTO_RENEW， 磁盘 将 是 automatically renewed monthly 当 账号 balance 是 sufficient。",
						},
						"time_unit": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "newly purchased 单位. 默认值：m。",
						},
						"cur_instance_deadline": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Current 实例 过期时间. Such 作为 2018-01-01 00:00:00. Specifying 此 参数 可以 align 过期时间 的 实例 attached 到 磁盘. One 的 此 参数 和 周期 必须 是 指定，和 不能 是 指定 在 same 时间。",
						},
					},
				},
			},

			"auto_voucher": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "是否automatically 使用 voucher. Not 使用 通过 默认值。",
			},
		},
	}
}

func resourceTencentCloudLighthouseRenewDiskCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_renew_disk.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = lighthouse.NewRenewDisksRequest()
		diskId  string
	)
	if v, ok := d.GetOk("disk_id"); ok {
		diskId = v.(string)
		request.DiskIds = []*string{&diskId}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "renew_disk_charge_prepaid"); ok {
		renewDiskChargePrepaid := lighthouse.RenewDiskChargePrepaid{}
		if v, ok := dMap["period"]; ok {
			renewDiskChargePrepaid.Period = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["renew_flag"]; ok {
			renewDiskChargePrepaid.RenewFlag = helper.String(v.(string))
		}
		if v, ok := dMap["time_unit"]; ok {
			renewDiskChargePrepaid.TimeUnit = helper.String(v.(string))
		}
		if v, ok := dMap["cur_instance_deadline"]; ok {
			renewDiskChargePrepaid.CurInstanceDeadline = helper.String(v.(string))
		}
		request.RenewDiskChargePrepaid = &renewDiskChargePrepaid
	}

	if v, _ := d.GetOk("auto_voucher"); v != nil {
		request.AutoVoucher = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().RenewDisks(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate lighthouse renewDisks failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(diskId)

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*tccommon.ReadRetryTimeout, time.Second, service.LighthouseDiskLatestOperationRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseRenewDiskRead(d, meta)
}

func resourceTencentCloudLighthouseRenewDiskRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_renew_disk.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudLighthouseRenewDiskDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_renew_disk.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
