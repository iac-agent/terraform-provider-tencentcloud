package css

import (
	"context"
	"fmt"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCssTimeshiftTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssTimeshiftTemplateCreate,
		Read:   resourceTencentCloudCssTimeshiftTemplateRead,
		Update: resourceTencentCloudCssTimeshiftTemplateUpdate,
		Delete: resourceTencentCloudCssTimeshiftTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"template_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "模板名称Maximum 长度: 255 bytes.Only letters，numbers，underscores，和 hyphens 是 支持。",
			},

			"duration": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "时间 shifting 时长.单位：Second。",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "模板描述Only letters，numbers，underscores，和 hyphens 是 支持。",
			},

			"area": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "地域`Mainland`: Chinese mainland.`Overseas`: Outside Chinese mainland.默认值：`Mainland`。",
			},

			"item_duration": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "segment 大小.取值范围：3-10.单位：Second.默认值：5。",
			},

			"remove_watermark": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "是否remove watermarks.如果 您 pass 在 `true`， original 流 将 是 recorded.默认值：`false`。",
			},

			"transcode_template_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "transcoding template IDs.此 API works 仅 如果 `RemoveWatermark` 是 `false`。",
			},
		},
	}
}

func resourceTencentCloudCssTimeshiftTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_timeshift_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = css.NewCreateLiveTimeShiftTemplateRequest()
		response   = css.NewCreateLiveTimeShiftTemplateResponse()
		templateId int64
	)
	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("duration"); ok {
		request.Duration = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("area"); ok {
		request.Area = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("item_duration"); ok {
		request.ItemDuration = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("remove_watermark"); ok {
		request.RemoveWatermark = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("transcode_template_ids"); ok {
		transcodeTemplateIdsSet := v.(*schema.Set).List()
		for i := range transcodeTemplateIdsSet {
			transcodeTemplateIds := transcodeTemplateIdsSet[i].(int)
			request.TranscodeTemplateIds = append(request.TranscodeTemplateIds, helper.IntInt64(transcodeTemplateIds))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCssClient().CreateLiveTimeShiftTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create css timeshiftTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(helper.Int64ToStr(templateId))

	return resourceTencentCloudCssTimeshiftTemplateRead(d, meta)
}

func resourceTencentCloudCssTimeshiftTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_timeshift_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	templateId := d.Id()
	templateIdInt64, err := strconv.ParseInt(templateId, 10, 64)
	if err != nil {
		return fmt.Errorf("TemplateId format type error: %s", err.Error())
	}

	timeshiftTemplate, err := service.DescribeCssTimeshiftTemplateById(ctx, templateIdInt64)
	if err != nil {
		return err
	}

	if timeshiftTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssTimeshiftTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if timeshiftTemplate.TemplateName != nil {
		_ = d.Set("template_name", timeshiftTemplate.TemplateName)
	}

	if timeshiftTemplate.Duration != nil {
		_ = d.Set("duration", timeshiftTemplate.Duration)
	}

	if timeshiftTemplate.Description != nil {
		_ = d.Set("description", timeshiftTemplate.Description)
	}

	if timeshiftTemplate.Area != nil {
		_ = d.Set("area", timeshiftTemplate.Area)
	}

	if timeshiftTemplate.ItemDuration != nil {
		_ = d.Set("item_duration", timeshiftTemplate.ItemDuration)
	}

	if timeshiftTemplate.RemoveWatermark != nil {
		_ = d.Set("remove_watermark", timeshiftTemplate.RemoveWatermark)
	}

	if timeshiftTemplate.TranscodeTemplateIds != nil {
		_ = d.Set("transcode_template_ids", timeshiftTemplate.TranscodeTemplateIds)
	}

	return nil
}

func resourceTencentCloudCssTimeshiftTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_timeshift_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := css.NewModifyLiveTimeShiftTemplateRequest()

	templateId := d.Id()
	templateIdInt64, _ := strconv.ParseInt(templateId, 10, 64)
	templateIdUint := uint64(templateIdInt64)

	request.TemplateId = &templateIdUint

	if d.HasChange("template_name") {
		if v, ok := d.GetOk("template_name"); ok {
			request.TemplateName = helper.String(v.(string))
		}
	}

	if d.HasChange("duration") {
		if v, ok := d.GetOkExists("duration"); ok {
			request.Duration = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("area") {
		if v, ok := d.GetOk("area"); ok {
			request.Area = helper.String(v.(string))
		}
	}

	if d.HasChange("item_duration") {
		if v, ok := d.GetOkExists("item_duration"); ok {
			request.ItemDuration = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("remove_watermark") {
		if v, ok := d.GetOkExists("remove_watermark"); ok {
			request.RemoveWatermark = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("transcode_template_ids") {
		if v, ok := d.GetOk("transcode_template_ids"); ok {
			transcodeTemplateIdsSet := v.(*schema.Set).List()
			for i := range transcodeTemplateIdsSet {
				transcodeTemplateIds := transcodeTemplateIdsSet[i].(int)
				request.TranscodeTemplateIds = append(request.TranscodeTemplateIds, helper.IntInt64(transcodeTemplateIds))
			}
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCssClient().ModifyLiveTimeShiftTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update css timeshiftTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCssTimeshiftTemplateRead(d, meta)
}

func resourceTencentCloudCssTimeshiftTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_timeshift_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	templateId := d.Id()
	templateIdInt64, _ := strconv.ParseInt(templateId, 10, 64)

	if err := service.DeleteCssTimeshiftTemplateById(ctx, templateIdInt64); err != nil {
		return err
	}

	return nil
}
