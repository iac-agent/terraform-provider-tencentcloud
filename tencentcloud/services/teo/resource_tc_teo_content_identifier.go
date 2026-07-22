package teo

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoContentIdentifier() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoContentIdentifierCreate,
		Read:   resourceTencentCloudTeoContentIdentifierRead,
		Update: resourceTencentCloudTeoContentIdentifierUpdate,
		Delete: resourceTencentCloudTeoContentIdentifierDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "描述 内容 identifier，长度 限制 的 up 到 20 字符。",
			},

			"plan_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target plan ID 到 是 bound，可用 仅 对于 enterprise edition. <li>如果 there 是 already plan under your 账号，go 到 [plan management](https://console.云.tencent.com/edgeone/包) 到 get plan ID 和 directly bind 内容 identifier 到 plan;</li><li>如果 您 do 不 have plan 到 bind，please purchase enterprise edition plan first.</li>。",
			},

			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "标签 的 内容 identifier. 此 参数 是 用于authority control. 到 create 标签，go 到 [标签 console](https://console.云.tencent.com/标签/taglist)。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "标签键\n注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "标签值\n注意：此字段可能返回 null，表示无法获取有效值。",
						},
					},
				},
			},

			// computed
			"content_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "内容 identifier ID。",
			},

			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "创建时间，其中 是 在 Coordinated Universal Time (UTC) 和 follows ISO 8601 date 和 时间格式.。",
			},

			"modified_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "时间 的 latest update，在 Coordinated Universal Time (UTC)，following ISO 8601 date 和 时间格式.。",
			},
		},
	}
}

func resourceTencentCloudTeoContentIdentifierCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_content_identifier.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = teov20220901.NewCreateContentIdentifierRequest()
		response = teov20220901.NewCreateContentIdentifierResponse()
	)

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("plan_id"); ok {
		request.PlanId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			tagsMap := item.(map[string]interface{})
			tag := teov20220901.Tag{}
			if v, ok := tagsMap["tag_key"].(string); ok && v != "" {
				tag.TagKey = helper.String(v)
			}

			if v, ok := tagsMap["tag_value"].(string); ok && v != "" {
				tag.TagValue = helper.String(v)
			}

			request.Tags = append(request.Tags, &tag)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateContentIdentifierWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo content identifier failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo content identifier failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.ContentId == nil {
		return fmt.Errorf("ContentId is nil.")
	}

	d.SetId(*response.Response.ContentId)
	return resourceTencentCloudTeoContentIdentifierRead(d, meta)
}

func resourceTencentCloudTeoContentIdentifierRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_content_identifier.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		contentId = d.Id()
	)

	respData, err := service.DescribeTeoContentIdentifierById(ctx, contentId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `teo_content_identifier` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.PlanId != nil {
		_ = d.Set("plan_id", respData.PlanId)
	}

	if respData.Tags != nil {
		tagsList := make([]map[string]interface{}, 0, len(respData.Tags))
		for _, tags := range respData.Tags {
			tagsMap := map[string]interface{}{}
			if tags.TagKey != nil {
				tagsMap["tag_key"] = tags.TagKey
			}

			if tags.TagValue != nil {
				tagsMap["tag_value"] = tags.TagValue
			}

			tagsList = append(tagsList, tagsMap)
		}

		_ = d.Set("tags", tagsList)
	}

	if respData.ContentId != nil {
		_ = d.Set("content_id", respData.ContentId)
	}

	if respData.CreatedOn != nil {
		_ = d.Set("created_on", respData.CreatedOn)
	}

	if respData.ModifiedOn != nil {
		_ = d.Set("modified_on", respData.ModifiedOn)
	}

	return nil
}

func resourceTencentCloudTeoContentIdentifierUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_content_identifier.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		contentId = d.Id()
	)

	immutableArgs := []string{"plan_id", "tags"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("description") {
		request := teov20220901.NewModifyContentIdentifierRequest()
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		request.ContentId = &contentId
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyContentIdentifierWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update teo content identifier failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudTeoContentIdentifierRead(d, meta)
}

func resourceTencentCloudTeoContentIdentifierDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_content_identifier.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = teov20220901.NewDeleteContentIdentifierRequest()
		contentId = d.Id()
	)

	request.ContentId = &contentId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteContentIdentifierWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo content identifier failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
