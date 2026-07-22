package cdc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdc/v20201214"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func ResourceTencentCloudCdcSite() *schema.Resource {
	return &schema.Resource{
		Create: ResourceTencentCloudCdcSiteCreate,
		Read:   ResourceTencentCloudCdcSiteRead,
		Update: ResourceTencentCloudCdcSiteUpdate,
		Delete: ResourceTencentCloudCdcSiteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site 名称",
			},
			"country": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site Country。",
			},
			"province": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site Province。",
			},
			"city": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site City。",
			},
			"address_line": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site Detail 地址",
			},
			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Site 描述",
			},
			//"note": {
			//	Optional:    true,
			//	Type:        schema.TypeString,
			//	Description: "Site 注意.",
			//},
			"fiber_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Site Fiber 类型 Using optical fiber 类型 到 connect CDC device 到 网络 SM(Single-模式) 或 MM(Multi-模式) fibers 是 可用。",
			},
			"optical_standard": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Site Optical Standard. Optical standard 用于connect CDC device 到 网络 此 字段 depends 在 uplink speed，optical fiber 类型，和 distance 到 upstream equipment. Allow 值: `SM`，`MM`。",
			},
			"power_connectors": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Site Power Connectors. Example: 380VAC3P。",
			},
			"power_feed_drop": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Site Power Feed Drop. Whether power 是 supplied 从 above 或 below rack. Allow 值: `UP`，`DOWN`。",
			},
			"max_weight": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Site Max 权重 容量 (KG)。",
			},
			"power_draw_kva": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Site Power DrawKva (KW)。",
			},
			"uplink_speed_gbps": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Uplink speed 从 网络 到 Tencent Cloud 地域",
			},
			"uplink_count": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "数量 uplinks 使用 通过 each CDC device (2 devices per rack) 当 connected 到 网络。",
			},
			"condition_requirement": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "是否following environmental conditions 是 met: n1. There 是 无 material requirements 或 acceptance standard 在 site 该 将 affect delivery 和 installation 的 CDC device. n2. following conditions 是 met 对于 finalized rack positions: Temperature ranges 从 41 到 104 degrees F (5 到 40 degrees C). Humidity ranges 从 10 degrees F (-12 degrees C) 到 70 degrees F (21 degrees C) 和 relative humidity ranges 从 8% RH 到 80% RH. Air flows 从 front 到 back 在 rack position 和 there 是 sufficient air 在 CFM (cubic feet per minute). air quantity 在 CFM 必须 是 145.8 times power consumption (在 KVA) 的 CDC。",
			},
			"dimension_requirement": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "是否following dimension conditions 是 met: Your loading dock 可以 accommodate 一个 rack 容器 (H x W x D = 94 x 54 x 48). You 可以 provide clear 路由 从 delivery point 的 your rack (H x W x D = 80 x 24 x 48) 到 its final installation location. You should consider platforms，corridors，doors，turns，ramps，freight elevators 作为 well 作为 other 访问 restrictions 当 measuring depth. There shall 是 48 或 greater front clearance 和 24 或 greater rear clearance 其中 CDC 是 finally installed。",
			},
			"redundant_networking": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether redundant upstream equipment (switch 或 router) 是 提供 so 该 both 网络 devices 可以 是 connected 到 网络。",
			},
			//"postal_code": {
			//	Optional:    true,
			//	Type:        schema.TypeInt,
			//	Description: "Postal 代码 的 site area.",
			//},
			"optional_address_line": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Detailed 地址 的 site area (到 是 added)。",
			},
			"need_help": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether 您 need help 从 Tencent Cloud 对于 rack installation。",
			},
			"redundant_power": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether there 是 power redundancy。",
			},
			"breaker_requirement": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether there 是 upstream circuit breaker。",
			},
		},
	}
}

func ResourceTencentCloudCdcSiteCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdc_site.create")()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = cdc.NewCreateSiteRequest()
		response = cdc.NewCreateSiteResponse()
		siteId   string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("country"); ok {
		request.Country = helper.String(v.(string))
	}

	if v, ok := d.GetOk("province"); ok {
		request.Province = helper.String(v.(string))
	}

	if v, ok := d.GetOk("city"); ok {
		request.City = helper.String(v.(string))
	}

	if v, ok := d.GetOk("address_line"); ok {
		request.AddressLine = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	//if v, ok := d.GetOk("note"); ok {
	//	request.Note = helper.String(v.(string))
	//}

	if v, ok := d.GetOk("fiber_type"); ok {
		request.FiberType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("optical_standard"); ok {
		request.OpticalStandard = helper.String(v.(string))
	}

	if v, ok := d.GetOk("power_connectors"); ok {
		request.PowerConnectors = helper.String(v.(string))
	}

	if v, ok := d.GetOk("power_feed_drop"); ok {
		request.PowerFeedDrop = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("max_weight"); ok {
		request.MaxWeight = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("power_draw_kva"); ok {
		request.PowerDrawKva = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("uplink_speed_gbps"); ok {
		request.UplinkSpeedGbps = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("uplink_count"); ok {
		request.UplinkCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("condition_requirement"); ok {
		request.ConditionRequirement = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("dimension_requirement"); ok {
		request.DimensionRequirement = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("redundant_networking"); ok {
		request.RedundantNetworking = helper.Bool(v.(bool))
	}

	//if v, ok := d.GetOkExists("postal_code"); ok {
	//	request.PostalCode = helper.IntInt64(v.(int))
	//}

	if v, ok := d.GetOk("optional_address_line"); ok {
		request.OptionalAddressLine = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("need_help"); ok {
		request.NeedHelp = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("redundant_power"); ok {
		request.RedundantPower = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("breaker_requirement"); ok {
		request.BreakerRequirement = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCdcClient().CreateSite(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("create cdc site failed")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cdc site failed, reason:%+v", logId, err)
		return err
	}

	siteId = *response.Response.SiteId
	d.SetId(siteId)

	return ResourceTencentCloudCdcSiteRead(d, meta)
}

func ResourceTencentCloudCdcSiteRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdc_site.read")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = CdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		siteId  = d.Id()
	)

	siteDetail, err := service.DescribeCdcSiteDetailById(ctx, siteId)
	if err != nil {
		return err
	}

	if siteDetail == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdcSite` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if siteDetail.Name != nil {
		_ = d.Set("name", siteDetail.Name)
	}

	if siteDetail.Country != nil {
		_ = d.Set("country", siteDetail.Country)
	}

	if siteDetail.Province != nil {
		_ = d.Set("province", siteDetail.Province)
	}

	if siteDetail.City != nil {
		_ = d.Set("city", siteDetail.City)
	}

	if siteDetail.AddressLine != nil {
		_ = d.Set("address_line", siteDetail.AddressLine)
	}

	if siteDetail.Description != nil {
		_ = d.Set("description", siteDetail.Description)
	}

	//if siteDetail.Note != nil {
	//	_ = d.Set("note", siteDetail.Note)
	//}

	if siteDetail.FiberType != nil {
		_ = d.Set("fiber_type", siteDetail.FiberType)
	}

	if siteDetail.OpticalStandard != nil {
		_ = d.Set("optical_standard", siteDetail.OpticalStandard)
	}

	if siteDetail.PowerConnectors != nil {
		_ = d.Set("power_connectors", siteDetail.PowerConnectors)
	}

	if siteDetail.PowerFeedDrop != nil {
		_ = d.Set("power_feed_drop", siteDetail.PowerFeedDrop)
	}

	if siteDetail.MaxWeight != nil {
		_ = d.Set("max_weight", siteDetail.MaxWeight)
	}

	if siteDetail.PowerDrawKva != nil {
		_ = d.Set("power_draw_kva", siteDetail.PowerDrawKva)
	}

	if siteDetail.UplinkSpeedGbps != nil {
		_ = d.Set("uplink_speed_gbps", siteDetail.UplinkSpeedGbps)
	}

	if siteDetail.UplinkCount != nil {
		_ = d.Set("uplink_count", siteDetail.UplinkCount)
	}

	if siteDetail.ConditionRequirement != nil {
		_ = d.Set("condition_requirement", siteDetail.ConditionRequirement)
	}

	if siteDetail.DimensionRequirement != nil {
		_ = d.Set("dimension_requirement", siteDetail.DimensionRequirement)
	}

	if siteDetail.RedundantNetworking != nil {
		_ = d.Set("redundant_networking", siteDetail.RedundantNetworking)
	}

	//if siteDetail.PostalCode != nil {
	//	_ = d.Set("postal_code", siteDetail.PostalCode)
	//}

	if siteDetail.OptionalAddressLine != nil {
		_ = d.Set("optional_address_line", siteDetail.OptionalAddressLine)
	}

	if siteDetail.NeedHelp != nil {
		_ = d.Set("need_help", siteDetail.NeedHelp)
	}

	if siteDetail.RedundantPower != nil {
		_ = d.Set("redundant_power", siteDetail.RedundantPower)
	}

	if siteDetail.BreakerRequirement != nil {
		_ = d.Set("breaker_requirement", siteDetail.BreakerRequirement)
	}

	return nil
}

func ResourceTencentCloudCdcSiteUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdc_site.update")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = cdc.NewModifySiteInfoRequest()
		siteId  = d.Id()
	)

	immutableArgs := []string{"fiber_type", "optical_standard", "power_connectors", "power_feed_drop", "max_weight", "power_draw_kva", "uplink_speed_gbps", "uplink_count", "condition_requirement", "dimension_requirement", "redundant_networking", "optional_address_line", "need_help", "redundant_power", "breaker_requirement"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.SiteId = &siteId
	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("country") {
		if v, ok := d.GetOk("country"); ok {
			request.Country = helper.String(v.(string))
		}
	}

	if d.HasChange("province") {
		if v, ok := d.GetOk("province"); ok {
			request.Province = helper.String(v.(string))
		}
	}

	if d.HasChange("city") {
		if v, ok := d.GetOk("city"); ok {
			request.City = helper.String(v.(string))
		}
	}

	if d.HasChange("address_line") {
		if v, ok := d.GetOk("address_line"); ok {
			request.AddressLine = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	//if d.HasChange("note") {
	//	if v, ok := d.GetOk("note"); ok {
	//		request.Note = helper.String(v.(string))
	//	}
	//}

	//if d.HasChange("postal_code") {
	//	if v, ok := d.GetOkExists("postal_code"); ok {
	//		request.PostalCode = helper.String(v.(string))
	//	}
	//}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCdcClient().ModifySiteInfo(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cdc site failed, reason:%+v", logId, err)
		return err
	}

	return ResourceTencentCloudCdcSiteRead(d, meta)
}

func ResourceTencentCloudCdcSiteDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdc_site.delete")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = CdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		siteId  = d.Id()
	)

	if err := service.DeleteCdcSiteById(ctx, siteId); err != nil {
		return err
	}

	return nil
}
