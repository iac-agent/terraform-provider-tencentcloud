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
				Description: "Inference API Token name, cannot exceed 30 characters.",
			},

			"token_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Inference API Token ID.",
			},

			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Inference API Token content.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time. The time is in Coordinated Universal Time (UTC) and follows the date and time format specified by the ISO 8601 standard.",
			},

			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Pagination query offset. Default value: 0.",
			},

			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     20,
				Description: "Pagination query limit. Default value: 20, maximum value: 100.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of Inference API Tokens.",
			},
		},
	}
}

func resourceTencentCloudTeoInferenceAPITokenCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewCreateInferenceAPITokenRequest()
		zoneId  string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	var response *teov20220901.CreateInferenceAPITokenResponse
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateInferenceAPITokenWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo inference_api_token failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create teo inference_api_token, logId: %s, id: %s", logId, logId, d.Id())

	if response == nil || response.Response == nil || response.Response.TokenId == nil {
		return fmt.Errorf("create teo inference_api_token failed, TokenId is nil")
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
		request = teov20220901.NewDescribeInferenceAPITokensRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	tokenId := idSplit[1]

	request.ZoneId = helper.String(zoneId)

	if v, ok := d.GetOk("offset"); ok {
		request.Offset = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("limit"); ok {
		request.Limit = helper.Int64(int64(v.(int)))
	}

	var response *teov20220901.DescribeInferenceAPITokensResponse
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DescribeInferenceAPITokensWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read teo inference_api_token failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response == nil || response.Response == nil || len(response.Response.Tokens) == 0 {
		log.Printf("[CRUD] teo inference_api_token id=%s", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	if response.Response.TotalCount != nil {
		_ = d.Set("total_count", response.Response.TotalCount)
	}

	var target *teov20220901.InferenceAPIToken
	for _, item := range response.Response.Tokens {
		if item.TokenId != nil && *item.TokenId == tokenId {
			target = item
			break
		}
	}

	if target == nil {
		log.Printf("[CRUD] teo inference_api_token id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if target.TokenId != nil {
		_ = d.Set("token_id", target.TokenId)
	}

	if target.Name != nil {
		_ = d.Set("name", target.Name)
	}

	if target.Content != nil {
		_ = d.Set("content", target.Content)
	}

	if target.CreateTime != nil {
		_ = d.Set("create_time", target.CreateTime)
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

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteInferenceAPITokenWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo inference_api_token failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
