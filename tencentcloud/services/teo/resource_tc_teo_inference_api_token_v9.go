package teo

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoInferenceAPITokenV9() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoInferenceAPITokenV9Create,
		Read:   resourceTencentCloudTeoInferenceAPITokenV9Read,
		Delete: resourceTencentCloudTeoInferenceAPITokenV9Delete,

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
				Description: "The name of the inference API token, limited to 30 characters.",
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
				Description: "Inference API Token content. Only returned once during creation, subsequent queries will not return this value.",
			},
		},
	}
}

func resourceTencentCloudTeoInferenceAPITokenV9Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token_v9.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = teov20220901.NewCreateInferenceAPITokenRequest()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	var response *teov20220901.CreateInferenceAPITokenResponse
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateInferenceAPIToken(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo inference_api_token failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[CRUD]%s create teo inference_api_token, d.Id()=%s", logId, d.Id())

	if response == nil || response.Response == nil || response.Response.TokenId == nil || *response.Response.TokenId == "" {
		return fmt.Errorf("Create inference_api_token response is nil or TokenId is empty.")
	}

	tokenId := *response.Response.TokenId
	d.SetId(tokenId)

	if response.Response.Content != nil {
		_ = d.Set("content", response.Response.Content)
	}

	return resourceTencentCloudTeoInferenceAPITokenV9Read(d, meta)
}

func resourceTencentCloudTeoInferenceAPITokenV9Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token_v9.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = teov20220901.NewDescribeInferenceAPITokensRequest()
		tokenId = d.Id()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
	}

	limit := int64(100)
	request.Limit = &limit

	var response *teov20220901.DescribeInferenceAPITokensResponse
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeInferenceAPITokens(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read teo inference_api_token failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response == nil || response.Response == nil || len(response.Response.Tokens) == 0 {
		log.Printf("[CRUD]%s resource `teo_inference_api_token_v9` [%s] not found, response is empty, d.SetId(\"\")", logId, d.Id())
		d.SetId("")
		return nil
	}

	var foundToken *teov20220901.InferenceAPIToken
	for _, token := range response.Response.Tokens {
		if token.TokenId != nil && *token.TokenId == tokenId {
			foundToken = token
			break
		}
	}

	if foundToken == nil {
		log.Printf("[CRUD]%s resource `teo_inference_api_token_v9` [%s] not found in response, d.SetId(\"\")", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", request.ZoneId)

	if foundToken.Name != nil {
		_ = d.Set("name", foundToken.Name)
	}

	if foundToken.TokenId != nil {
		_ = d.Set("token_id", foundToken.TokenId)
	}

	// Content is only returned once during creation, do not overwrite in Read

	return nil
}

func resourceTencentCloudTeoInferenceAPITokenV9Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token_v9.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = teov20220901.NewDeleteInferenceAPITokenRequest()
		tokenId = d.Id()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
	}

	request.TokenId = &tokenId

	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteInferenceAPIToken(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo inference_api_token failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
