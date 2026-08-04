package teo

import (
	"fmt"
	"log"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudTeoInferenceAPITokenOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoInferenceAPITokenOperationCreate,
		Read:   resourceTencentCloudTeoInferenceAPITokenOperationRead,
		Delete: resourceTencentCloudTeoInferenceAPITokenOperationDelete,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Inference API Token name, length limit does not exceed 30 characters.",
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
		},
	}
}

func resourceTencentCloudTeoInferenceAPITokenOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	zoneId := d.Get("zone_id").(string)
	name := d.Get("name").(string)

	if zoneId == "" {
		return fmt.Errorf("zone_id is required")
	}

	if name == "" {
		return fmt.Errorf("name is required")
	}

	request := teov20220901.NewCreateInferenceAPITokenRequest()
	request.ZoneId = helper.String(zoneId)
	request.Name = helper.String(name)

	var response *teov20220901.CreateInferenceAPITokenResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateInferenceAPIToken(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		return err
	}

	if response == nil || response.Response == nil {
		return fmt.Errorf("CreateInferenceAPIToken API returned empty response")
	}

	d.SetId(zoneId)

	if response.Response.TokenId != nil {
		_ = d.Set("token_id", *response.Response.TokenId)
	}

	if response.Response.Content != nil {
		_ = d.Set("content", *response.Response.Content)
	}

	return resourceTencentCloudTeoInferenceAPITokenOperationRead(d, meta)
}

func resourceTencentCloudTeoInferenceAPITokenOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTeoInferenceAPITokenOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
