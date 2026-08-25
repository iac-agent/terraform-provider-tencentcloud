package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func ResourceTencentCloudTeoInferenceApiToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoInferenceApiTokenCreate,
		Read:   resourceTencentCloudTeoInferenceApiTokenRead,
		Update: resourceTencentCloudTeoInferenceApiTokenUpdate,
		Delete: resourceTencentCloudTeoInferenceApiTokenDelete,
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
				Description: "Name of the inference API token, up to 30 characters.",
			},

			"token_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Inference API token ID.",
			},

			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Inference API token content.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time, in ISO 8601 date format.",
			},
		},
	}
}

func resourceTencentCloudTeoInferenceApiTokenCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		zoneId  string
		name    string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	token, err := service.AddTeoInferenceApiToken(ctx, zoneId, name)
	if err != nil {
		log.Printf("[CRITAL]%s create teo_inference_api_token failed, reason:%+v", logId, err)
		return err
	}

	if token == nil {
		return fmt.Errorf("Create teo_inference_api_token failed, Response is nil.")
	}

	if token.TokenId == nil {
		log.Printf("[CRITAL]%s create teo_inference_api_token failed, TokenId is nil, logId=%s", logId, logId)
		return fmt.Errorf("Create teo_inference_api_token failed, TokenId is nil.")
	}

	if *token.TokenId == "" {
		log.Printf("[CRITAL]%s create teo_inference_api_token failed, TokenId is empty, logId=%s", logId, logId)
		return fmt.Errorf("Create teo_inference_api_token failed, TokenId is empty.")
	}

	d.SetId(strings.Join([]string{zoneId, *token.TokenId}, tccommon.FILED_SP))
	return resourceTencentCloudTeoInferenceApiTokenRead(d, meta)
}

func resourceTencentCloudTeoInferenceApiTokenRead(d *schema.ResourceData, meta interface{}) error {
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

	token, err := service.DescribeTeoInferenceApiTokenById(ctx, zoneId, tokenId)
	if err != nil {
		return err
	}

	if token == nil {
		log.Printf("[CRUD] teo_inference_api_token id=%s", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	if token.Name != nil {
		_ = d.Set("name", token.Name)
	}

	if token.TokenId != nil {
		_ = d.Set("token_id", token.TokenId)
	}

	if token.Content != nil {
		_ = d.Set("content", token.Content)
	}

	if token.CreateTime != nil {
		_ = d.Set("create_time", token.CreateTime)
	}

	return nil
}

func resourceTencentCloudTeoInferenceApiTokenUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{"name"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	return resourceTencentCloudTeoInferenceApiTokenRead(d, meta)
}

func resourceTencentCloudTeoInferenceApiTokenDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_api_token.delete")()
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

	err := service.DeleteTeoInferenceApiToken(ctx, zoneId, tokenId)
	if err != nil {
		log.Printf("[CRITAL]%s delete teo_inference_api_token failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
