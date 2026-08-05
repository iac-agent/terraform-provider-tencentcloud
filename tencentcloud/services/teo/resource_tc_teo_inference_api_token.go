package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoInferenceAPIToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoInferenceAPITokenCreate,
		Read:   resourceTencentCloudTeoInferenceAPITokenRead,
		Delete: resourceTencentCloudTeoInferenceAPITokenDelete,
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

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Inference API Token name, with a length limit of no more than 30 characters.",
			},

			"token_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Inference API Token ID.",
			},

			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Inference API Token content.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time, in ISO date format.",
			},
		},
	}
}

func resourceTencentCloudTeoInferenceAPITokenCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = teov20220901.NewCreateInferenceAPITokenRequest()
		response = teov20220901.NewCreateInferenceAPITokenResponse()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateInferenceAPITokenWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo inference api token failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo inference api token failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.TokenId == nil {
		return fmt.Errorf("TokenId is nil.")
	}

	d.SetId(*response.Response.TokenId)

	return resourceTencentCloudTeoInferenceAPITokenRead(d, meta)
}

func resourceTencentCloudTeoInferenceAPITokenRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		tokenId = d.Id()
	)

	zoneId := ""
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}
	if zoneId == "" {
		// For import, the ID is in zone_id#token_id format
		parts := strings.Split(tokenId, tccommon.FILED_SP)
		if len(parts) == 2 {
			zoneId = parts[0]
			tokenId = parts[1]
			_ = d.Set("zone_id", zoneId)
			d.SetId(tokenId)
		}
	}

	respData, err := service.DescribeTeoInferenceAPITokenById(ctx, zoneId, tokenId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `teo_inference_api_token` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.TokenId != nil {
		_ = d.Set("token_id", respData.TokenId)
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Content != nil {
		_ = d.Set("content", respData.Content)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	return nil
}

func resourceTencentCloudTeoInferenceAPITokenDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewDeleteInferenceAPITokenRequest()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
	}

	request.TokenId = helper.String(d.Id())

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteInferenceAPITokenWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo inference api token failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
