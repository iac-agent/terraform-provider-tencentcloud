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
		Update: resourceTencentCloudTeoInferenceAPITokenUpdate,
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
		zoneId   string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateInferenceAPITokenWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create inference_api_token failed, Response is nil"))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create inference_api_token failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.TokenId == nil || *response.Response.TokenId == "" {
		return fmt.Errorf("Create inference_api_token returned empty TokenId")
	}

	tokenId := *response.Response.TokenId
	d.SetId(strings.Join([]string{zoneId, tokenId}, tccommon.FILED_SP))

	return resourceTencentCloudTeoInferenceAPITokenRead(d, meta)
}

func resourceTencentCloudTeoInferenceAPITokenRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	tokenId := idSplit[1]

	_ = d.Set("zone_id", zoneId)

	respData, err := service.DescribeTeoInferenceAPITokenById(ctx, zoneId, tokenId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] inference_api_token id=%s", d.Id())
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

func resourceTencentCloudTeoInferenceAPITokenUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{"zone_id", "name"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	return resourceTencentCloudTeoInferenceAPITokenRead(d, meta)
}

func resourceTencentCloudTeoInferenceAPITokenDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewDeleteInferenceAPITokenRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	tokenId := idSplit[1]

	request.ZoneId = helper.String(zoneId)
	request.TokenId = helper.String(tokenId)

	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteInferenceAPITokenWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete inference_api_token failed, Response is nil"))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete inference_api_token failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
