package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoInferenceDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoInferenceDomainCreate,
		Read:   resourceTencentCloudTeoInferenceDomainRead,
		Update: resourceTencentCloudTeoInferenceDomainUpdate,
		Delete: resourceTencentCloudTeoInferenceDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID.",
			},
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Inference service ID.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Inference service domain.",
			},
			"auth_switch": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Inference task request authentication switch. Valid values: `Off`, `On`. Default value: `On`.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Domain status. Valid values: `Online`, `Process`, `Offline`, `Init`, `Forbidden`.",
			},
			"cname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CNAME address.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modification time.",
			},
		},
	}
}

func resourceTencentCloudTeoInferenceDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_domain.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(
			context.Background(), logId, d, meta)
		request = teo.NewCreateInferenceDomainRequest()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_id"); ok {
		request.ServiceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain"); ok {
		request.Domain = helper.String(v.(string))
	}

	if v, ok := d.GetOk("auth_switch"); ok {
		request.AuthSwitch = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().CreateInferenceDomainWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo inference domain failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo inference domain failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{
		d.Get("zone_id").(string),
		d.Get("service_id").(string),
		d.Get("domain").(string),
	}, tccommon.FILED_SP))

	return resourceTencentCloudTeoInferenceDomainRead(d, meta)
}

func resourceTencentCloudTeoInferenceDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_domain.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(
			context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	serviceId := idSplit[1]
	domain := idSplit[2]

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("service_id", serviceId)
	_ = d.Set("domain", domain)

	respData, err := service.DescribeTeoInferenceDomainById(ctx, zoneId, serviceId, domain)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `teo_inference_domain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Domain != nil {
		_ = d.Set("domain", respData.Domain)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.AuthSwitch != nil {
		_ = d.Set("auth_switch", respData.AuthSwitch)
	}

	if respData.Cname != nil {
		_ = d.Set("cname", respData.Cname)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}

	return nil
}

func resourceTencentCloudTeoInferenceDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_domain.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(
			context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	serviceId := idSplit[1]
	domain := idSplit[2]

	if d.HasChange("auth_switch") {
		request := teo.NewCreateInferenceDomainRequest()
		request.ZoneId = helper.String(zoneId)
		request.ServiceId = helper.String(serviceId)
		request.Domain = helper.String(domain)
		if v, ok := d.GetOk("auth_switch"); ok {
			request.AuthSwitch = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().CreateInferenceDomainWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update teo inference domain auth switch failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTeoInferenceDomainRead(d, meta)
}

func resourceTencentCloudTeoInferenceDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_domain.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(
			context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	serviceId := idSplit[1]
	domain := idSplit[2]

	request := teo.NewOperateInferenceDomainRequest()
	request.ZoneId = helper.String(zoneId)
	request.ServiceId = helper.String(serviceId)
	request.Domain = helper.String(domain)
	request.Operation = helper.String("Delete")

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().OperateInferenceDomainWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete teo inference domain failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete teo inference domain failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
