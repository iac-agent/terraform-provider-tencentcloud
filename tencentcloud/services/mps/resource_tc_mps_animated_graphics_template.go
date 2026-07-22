package mps

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMpsAnimatedGraphicsTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsAnimatedGraphicsTemplateCreate,
		Read:   resourceTencentCloudMpsAnimatedGraphicsTemplateRead,
		Update: resourceTencentCloudMpsAnimatedGraphicsTemplateUpdate,
		Delete: resourceTencentCloudMpsAnimatedGraphicsTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"fps": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Frame 速率，取值范围：[1，30]，单位: Hz。",
			},

			"width": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "最大 值 的 animation 宽度 (或 long side)，取值范围：0 和 [128，4096]，单位: 像素.当 宽度 和 高度 是 both 0， resolution 是 same.当 宽度 是 0 和 高度 是 不 0，宽度 是 scaled proportionally.当 宽度 是 不 0 和 高度 是 0，高度 是 scaled proportionally.当 both 宽度 和 高度 是 不 0， resolution 是 指定 通过 用户默认值：0。",
			},

			"height": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "最大 值 的 animation 高度 (或 short side)，取值范围：0 和 [128，4096]，单位: 像素.当 宽度 和 高度 是 both 0， resolution 是 same.当 宽度 是 0 和 高度 是 不 0，宽度 是 scaled proportionally.当 宽度 是 不 0 和 高度 是 0，高度 是 scaled proportionally.当 both 宽度 和 高度 是 不 0， resolution 是 指定 通过 用户默认值：0。",
			},

			"resolution_adaptive": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Adaptive resolution，可选 值:open: At 此 时间，宽度 表示 long side 的 视频，高度 表示 short side 的 视频.close: At 此 point，宽度 表示 宽度 的 视频，和 高度 表示 高度 的 视频.默认值：open。",
			},

			"format": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Animation 格式， 值 是 gif 和 webp. 默认为 gif。",
			},

			"quality": {
				Optional:    true,
				Type:        schema.TypeFloat,
				Description: "Image quality，取值范围：[1，100]，默认值为 75。",
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Rotation diagram 模板名称，长度 限制: 64 字符。",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "模板描述 信息，长度 限制: 256 字符。",
			},
		},
	}
}

func resourceTencentCloudMpsAnimatedGraphicsTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_animated_graphics_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = mps.NewCreateAnimatedGraphicsTemplateRequest()
		response   = mps.NewCreateAnimatedGraphicsTemplateResponse()
		definition uint64
	)
	if v, ok := d.GetOkExists("fps"); ok {
		request.Fps = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("width"); ok {
		request.Width = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("height"); ok {
		request.Height = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("resolution_adaptive"); ok {
		request.ResolutionAdaptive = helper.String(v.(string))
	}

	if v, ok := d.GetOk("format"); ok {
		request.Format = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("quality"); ok {
		request.Quality = helper.Float64(v.(float64))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().CreateAnimatedGraphicsTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps animatedGraphicsTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition = *response.Response.Definition
	d.SetId(helper.UInt64ToStr(definition))

	return resourceTencentCloudMpsAnimatedGraphicsTemplateRead(d, meta)
}

func resourceTencentCloudMpsAnimatedGraphicsTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_animated_graphics_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	definition := d.Id()

	animatedGraphicsTemplate, err := service.DescribeMpsAnimatedGraphicsTemplateById(ctx, definition)
	if err != nil {
		return err
	}

	if animatedGraphicsTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsAnimatedGraphicsTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if animatedGraphicsTemplate.Fps != nil {
		_ = d.Set("fps", animatedGraphicsTemplate.Fps)
	}

	if animatedGraphicsTemplate.Width != nil {
		_ = d.Set("width", animatedGraphicsTemplate.Width)
	}

	if animatedGraphicsTemplate.Height != nil {
		_ = d.Set("height", animatedGraphicsTemplate.Height)
	}

	if animatedGraphicsTemplate.ResolutionAdaptive != nil {
		_ = d.Set("resolution_adaptive", animatedGraphicsTemplate.ResolutionAdaptive)
	}

	if animatedGraphicsTemplate.Format != nil {
		_ = d.Set("format", animatedGraphicsTemplate.Format)
	}

	if animatedGraphicsTemplate.Quality != nil {
		_ = d.Set("quality", animatedGraphicsTemplate.Quality)
	}

	if animatedGraphicsTemplate.Name != nil {
		_ = d.Set("name", animatedGraphicsTemplate.Name)
	}

	if animatedGraphicsTemplate.Comment != nil {
		_ = d.Set("comment", animatedGraphicsTemplate.Comment)
	}

	return nil
}

func resourceTencentCloudMpsAnimatedGraphicsTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_animated_graphics_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := mps.NewModifyAnimatedGraphicsTemplateRequest()

	definition := d.Id()

	request.Definition = helper.StrToUint64Point(definition)

	mutableArgs := []string{"fps", "width", "height", "resolution_adaptive", "format", "quality", "name", "comment"}

	needChange := false

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {

		if v, ok := d.GetOkExists("fps"); ok {
			request.Fps = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("width"); ok {
			request.Width = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("height"); ok {
			request.Height = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("resolution_adaptive"); ok {
			request.ResolutionAdaptive = helper.String(v.(string))
		}

		if v, ok := d.GetOk("format"); ok {
			request.Format = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("quality"); ok {
			request.Quality = helper.Float64(v.(float64))
		}

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().ModifyAnimatedGraphicsTemplate(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mps animatedGraphicsTemplate failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudMpsAnimatedGraphicsTemplateRead(d, meta)
}

func resourceTencentCloudMpsAnimatedGraphicsTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_animated_graphics_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	definition := d.Id()

	if err := service.DeleteMpsAnimatedGraphicsTemplateById(ctx, definition); err != nil {
		return err
	}

	return nil
}
