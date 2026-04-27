package teo

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoMultiPathGatewayLine() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoMultiPathGatewayLineCreate,
		Read:   resourceTencentCloudTeoMultiPathGatewayLineRead,
		Update: resourceTencentCloudTeoMultiPathGatewayLineUpdate,
		Delete: resourceTencentCloudTeoMultiPathGatewayLineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Multi-path gateway ID.",
			},
			"line_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Line type. Valid values: direct, proxy, custom.",
			},
			"line_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Line address, format is ip:port.",
			},
			"proxy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "L4 proxy instance ID. Required when line_type is proxy.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Forwarding rule ID. Required when line_type is proxy.",
			},
			"line_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Line ID.",
			},
		},
	}
}

func resourceTencentCloudTeoMultiPathGatewayLineCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_multi_path_gateway_line.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		zoneId    = d.Get("zone_id").(string)
		gatewayId = d.Get("gateway_id").(string)
		response  *teo.CreateMultiPathGatewayLineResponse
	)

	request := teo.NewCreateMultiPathGatewayLineRequest()
	request.ZoneId = helper.String(zoneId)
	request.GatewayId = helper.String(gatewayId)

	if v, ok := d.GetOk("line_type"); ok {
		request.LineType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("line_address"); ok {
		request.LineAddress = helper.String(v.(string))
	}

	if v, ok := d.GetOk("proxy_id"); ok {
		request.ProxyId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_id"); ok {
		request.RuleId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateMultiPathGatewayLine(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil || result.Response.LineId == nil {
			return resource.NonRetryableError(fmt.Errorf("create teo multi-path gateway line failed, response is nil"))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s create teo multi-path gateway line failed, reason: %s", logId, err.Error())
		return err
	}

	lineId := *response.Response.LineId
	_ = d.Set("line_id", lineId)

	d.SetId(zoneId + tccommon.FILED_SP + gatewayId + tccommon.FILED_SP + lineId)

	return resourceTencentCloudTeoMultiPathGatewayLineRead(d, meta)
}

func resourceTencentCloudTeoMultiPathGatewayLineRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_multi_path_gateway_line.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	zoneId := d.Get("zone_id").(string)
	gatewayId := d.Get("gateway_id").(string)
	lineId := d.Get("line_id").(string)

	if zoneId == "" || gatewayId == "" || lineId == "" {
		// Fallback to parsing from d.Id() for import case
		idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("resource id is broken, id is %s", d.Id())
		}
		zoneId = idSplit[0]
		gatewayId = idSplit[1]
		lineId = idSplit[2]
	}

	request := teo.NewDescribeMultiPathGatewayLineRequest()
	request.ZoneId = helper.String(zoneId)
	request.GatewayId = helper.String(gatewayId)
	request.LineId = helper.String(lineId)

	var line *teo.MultiPathGatewayLine
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeMultiPathGatewayLine(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil || response.Response.Line == nil {
			return resource.NonRetryableError(fmt.Errorf("line not found"))
		}

		line = response.Response.Line
		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s read teo multi-path gateway line failed, reason: %s", logId, err.Error())
		return err
	}

	if line == nil {
		log.Printf("[CRITICAL]%s read teo multi-path gateway line failed, reason: line(id:%s) not found", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("gateway_id", gatewayId)

	if line.LineId != nil {
		_ = d.Set("line_id", line.LineId)
	}
	if line.LineType != nil {
		_ = d.Set("line_type", line.LineType)
	}
	if line.LineAddress != nil {
		_ = d.Set("line_address", line.LineAddress)
	}
	if line.ProxyId != nil {
		_ = d.Set("proxy_id", line.ProxyId)
	}
	if line.RuleId != nil {
		_ = d.Set("rule_id", line.RuleId)
	}

	return nil
}

func resourceTencentCloudTeoMultiPathGatewayLineUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_multi_path_gateway_line.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	zoneId := d.Get("zone_id").(string)
	gatewayId := d.Get("gateway_id").(string)
	lineId := d.Get("line_id").(string)

	request := teo.NewModifyMultiPathGatewayLineRequest()
	request.ZoneId = helper.String(zoneId)
	request.GatewayId = helper.String(gatewayId)
	request.LineId = helper.String(lineId)

	if v, ok := d.GetOk("line_type"); ok {
		request.LineType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("line_address"); ok {
		request.LineAddress = helper.String(v.(string))
	}

	if v, ok := d.GetOk("proxy_id"); ok {
		request.ProxyId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_id"); ok {
		request.RuleId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyMultiPathGatewayLine(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s update teo multi-path gateway line failed, reason: %s", logId, err.Error())
		return err
	}

	return resourceTencentCloudTeoMultiPathGatewayLineRead(d, meta)
}

func resourceTencentCloudTeoMultiPathGatewayLineDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_multi_path_gateway_line.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	zoneId := d.Get("zone_id").(string)
	gatewayId := d.Get("gateway_id").(string)
	lineId := d.Get("line_id").(string)

	request := teo.NewDeleteMultiPathGatewayLineRequest()
	request.ZoneId = helper.String(zoneId)
	request.GatewayId = helper.String(gatewayId)
	request.LineId = helper.String(lineId)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteMultiPathGatewayLine(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if err != nil {
		log.Printf("[CRITICAL]%s delete teo multi-path gateway line failed, reason: %s", logId, err.Error())
		return err
	}

	return nil
}
