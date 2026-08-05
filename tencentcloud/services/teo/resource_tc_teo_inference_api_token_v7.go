package teo

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func ResourceTencentCloudTeoInferenceApiTokenV7() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoInferenceApiTokenV7Create,
		Read:   resourceTencentCloudTeoInferenceApiTokenV7Read,
		Delete: resourceTencentCloudTeoInferenceApiTokenV7Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The zone ID where the inference API token belongs to.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the inference API token, with a maximum length of 30 characters.",
			},
			"token_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique ID of the inference API token.",
			},
			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The content of the inference API token.",
			},
		},
	}
}

func resourceTencentCloudTeoInferenceApiTokenV7Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token_v7.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		zoneId string
		name   string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	request := teov20220901.NewCreateInferenceAPITokenRequest()
	request.ZoneId = &zoneId
	request.Name = &name

	var response *teov20220901.CreateInferenceAPITokenResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().CreateInferenceAPITokenWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo inference_api_token failed, reason:%+v", logId, err)
		return err
	}

	if response == nil || response.Response == nil || response.Response.TokenId == nil || *response.Response.TokenId == "" {
		log.Printf("[CRITAL]%s create teo inference_api_token failed, response is empty, logId: %s", logId, logId)
		return fmt.Errorf("create teo inference_api_token failed: empty response")
	}

	d.SetId(*response.Response.TokenId)

	return resourceTencentCloudTeoInferenceApiTokenV7Read(d, meta)
}

func resourceTencentCloudTeoInferenceApiTokenV7Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token_v7.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	tokenId := d.Id()

	var zoneId string
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	limit := int64(100)
	offset := int64(0)

	var targetToken *teov20220901.InferenceAPIToken
	for {
		request := teov20220901.NewDescribeInferenceAPITokensRequest()
		request.ZoneId = &zoneId
		request.Limit = &limit
		request.Offset = &offset

		var response *teov20220901.DescribeInferenceAPITokensResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().DescribeInferenceAPITokensWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read teo inference_api_token failed, reason:%+v", logId, err)
			return err
		}

		if response == nil || response.Response == nil || len(response.Response.Tokens) == 0 {
			break
		}

		for _, token := range response.Response.Tokens {
			if token.TokenId != nil && *token.TokenId == tokenId {
				targetToken = token
				break
			}
		}

		if targetToken != nil {
			break
		}

		if response.Response.TotalCount == nil || int64(len(response.Response.Tokens)) < limit || offset+limit >= *response.Response.TotalCount {
			break
		}

		offset += limit
	}

	if targetToken == nil {
		log.Printf("[WARN]%s resource tencentcloud_teo_inference_api_token_v7 token_id [%s] not found in zone_id [%s]", logId, tokenId, zoneId)
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	if targetToken.Name != nil {
		_ = d.Set("name", targetToken.Name)
	}

	if targetToken.TokenId != nil {
		_ = d.Set("token_id", targetToken.TokenId)
	}

	if targetToken.Content != nil {
		_ = d.Set("content", targetToken.Content)
	}

	return nil
}

func resourceTencentCloudTeoInferenceApiTokenV7Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token_v7.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	tokenId := d.Id()

	var zoneId string
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	request := teov20220901.NewDeleteInferenceAPITokenRequest()
	request.ZoneId = &zoneId
	request.TokenId = &tokenId

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().DeleteInferenceAPITokenWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete teo inference_api_token failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
