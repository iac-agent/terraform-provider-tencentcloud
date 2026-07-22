package mps

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMpsWatermarkTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsWatermarkTemplateCreate,
		Read:   resourceTencentCloudMpsWatermarkTemplateRead,
		Update: resourceTencentCloudMpsWatermarkTemplateUpdate,
		Delete: resourceTencentCloudMpsWatermarkTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Watermark 类型，可选 值:镜像，text，svg。",
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Watermark 模板名称，长度 限制: 64 字符。",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "模板描述 信息，长度 限制: 256 字符。",
			},

			"coordinate_origin": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Origin position，可选 值:TopLeft: 表示that 源站 的 coordinates 是 在 upper left corner 的 视频 镜像，和 源站 的 水印 是 upper left corner 的 picture 或 text.TopRight: 表示that 源站 的 coordinates 是 在 upper right corner 的 视频 镜像，和 源站 的 水印 是 在 upper right corner 的 picture 或 text.BottomLeft: 表示that 源站 的 coordinates 是 在 lower left corner 的 视频 镜像，和 源站 的 水印 是 lower left corner 的 picture 或 text.BottomRight: 表示that 源站 的 coordinates 是 在 lower right corner 的 视频 镜像，和 源站 的 水印 是 在 lower right corner 的 picture 或 text.默认值：TopLeft。",
			},

			"x_pos": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "horizontal position 的 源站 的 水印 从 源站 的 coordinates 的 视频 镜像. Support %，像素 two formats.当 字符串 结束 使用 %，它 表示 该 水印 XPos 指定a percentage 对于 视频 宽度，such 作为 10% 表示 该 XPos 是 10% 的 视频 宽度.当 字符串 结束 使用 像素，它 表示 该 水印 XPos 是 指定 pixel，such 作为 100px 表示 该 XPos 是 100 pixels.默认值：0px。",
			},

			"y_pos": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "vertical position 的 源站 的 水印 从 源站 的 coordinates 的 视频 镜像. Support %，像素 two formats.当 字符串 结束 使用 %，它 表示 该 水印 YPos 指定a percentage 对于 视频 高度，such 作为 10% 表示 该 YPos 是 10% 的 视频 高度.当 字符串 结束 使用 像素，它 表示 该 水印 YPos 是 指定 pixel，such 作为 100px 表示 该 YPos 是 100 pixels.默认值：0px。",
			},

			"image_template": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Image 水印 template，仅 当 类型 是 镜像，此 字段 为必填项 和 有效。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Watermark 镜像[Base64](https://tools.ietf.org/html/rfc4648) encoded 字符串. Support jpeg，png 镜像 格式",
						},
						"width": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "宽度 的 水印. Support %，像素 two formats:当 字符串 结束 使用 %，它 表示 该 水印 宽度 是 percentage 的 视频 宽度，such 作为 10% 表示 该 宽度 是 10% 的 视频 宽度.当 字符串 结束 使用 像素，它 表示 该 水印 宽度 单位 是 pixel，such 作为 100px 表示 该 宽度 是 100 pixels. 值 范围 是 [8，4096].默认值：10%。",
						},
						"height": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "高度 的 水印. Support %，像素 two formats:当 字符串 结束 使用 %，它 表示 该 水印 高度 是 percentage 大小 的 视频 高度，such 作为 10% 表示 该 高度 是 10% 的 视频 高度.当 字符串 结束 使用 像素，它 表示 该 水印 高度 单位 是 pixel，such 作为 100px 表示 该 高度 是 100 pixels. 值 范围 是 0 或 [8，4096].默认值：0px. 表示that 高度 是 scaled according 到 aspect ratio 的 original 水印 镜像。",
						},
						"repeat_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Watermark repeat 类型 Usage scenario: 水印 是 动态 镜像. Ranges:once: After 动态 水印 是 played，它 将 无 longer appear.repeat_last_frame: After 水印 是 played，stay 在 last frame.repeat: 水印 loops until end 的 视频 (默认值)。",
						},
					},
				},
			},

			"text_template": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Text 水印 template，仅 当 类型 是 text，此 字段 为必填项 和 有效。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"font_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Font 类型，currently 支持 two:simkai.ttf: 可以 support Chinese 和 English.arial.ttf: English 仅。",
						},
						"font_size": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Font 大小，格式: Npx，N 是 数量。",
						},
						"font_color": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Font color，格式: 0xRRGGBB，默认值：0xFFFFFF (white)。",
						},
						"font_alpha": {
							Type:        schema.TypeFloat,
							Required:    true,
							Description: "Text transparency，取值范围：(0，1].0: fully transparent.1: fully opaque.默认值：1。",
						},
					},
				},
			},

			"svg_template": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "SVG 水印 template，仅 当 类型 是 svg，此 字段 为必填项 和 有效。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"width": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "宽度 的 水印，支持 像素，%，W%，H%，S%，L% six formats.当 字符串 结束 使用 像素，它 表示 该 水印 宽度 单位 是 pixels，such 作为 100px 表示 该 宽度 是 100 pixels; 当 filling 0px 和 高度 是 不 0px，它 表示 该 宽度 的 水印 是 proportionally scaled according 到 original SVG 镜像; 当 both 宽度 和 高度 是 filled 当 0px，它 表示 该 宽度 的 水印 takes 宽度 的 original SVG 镜像.当 字符串 结束 使用 W%，它 表示 该 水印 宽度 是 percentage 的 视频 宽度，such 作为 10W% 表示 该 宽度 是 10% 的 视频 宽度.当 字符串 结束 使用 H%，它 表示 该 水印 宽度 是 percentage 的 视频 高度，such 作为 10H% 表示 该 宽度 是 10% 的 视频 高度.当 字符串 结束 使用 S%，它 表示 该 水印 宽度 是 percentage 大小 的 short side 的 视频，such 作为 10S% 表示 该 宽度 是 10% 的 short side 的 视频.当 字符串 结束 使用 L%，它 表示 该 水印 宽度 是 percentage 大小 的 long side 的 视频，such 作为 10L% 表示 该 宽度 是 10% 的 long side 的 视频.当 字符串 结束 使用 %，它 has same meaning 作为 W%.默认值：10W%。",
						},
						"height": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "高度 的 水印，支持 像素，W%，H%，S%，L% six formats:当 字符串 结束 使用 像素，它 表示 该 水印 高度 单位 是 pixels，such 作为 100px 表示 该 高度 是 100 pixels; 当 filling 0px 和 宽度 是 不 0px，它 表示 该 高度 的 水印 是 proportionally scaled according 到 original SVG 镜像; 当 both 宽度 和 高度 是 filled 当 0px，它 表示 该 高度 的 水印 takes 高度 的 original SVG 镜像.当 字符串 结束 使用 W%，它 表示 该 水印 高度 是 percentage 的 视频 宽度，such 作为 10W% 表示 该 高度 是 10% 的 视频 宽度.当 字符串 结束 使用 H%，它 表示 该 水印 高度 是 percentage 大小 的 视频 高度，such 作为 10H% 表示 该 高度 是 10% 的 视频 高度.当 字符串 结束 使用 S%，它 表示 该 水印 高度 是 percentage 大小 的 short side 的 视频，such 作为 10S% 表示 该 高度 是 10% 的 short side 的 视频.当 字符串 结束 使用 L%，它 表示 该 水印 高度 是 percentage 大小 的 long side 的 视频，such 作为 10L% 表示 该 高度 是 10% 的 long side 的 视频.当 字符串 结束 使用 %， meaning 是 same 作为 H%.默认值：0px。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMpsWatermarkTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_watermark_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = mps.NewCreateWatermarkTemplateRequest()
		response   = mps.NewCreateWatermarkTemplateResponse()
		definition int64
	)
	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	if v, ok := d.GetOk("coordinate_origin"); ok {
		request.CoordinateOrigin = helper.String(v.(string))
	}

	if v, ok := d.GetOk("x_pos"); ok {
		request.XPos = helper.String(v.(string))
	}

	if v, ok := d.GetOk("y_pos"); ok {
		request.YPos = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "image_template"); ok {
		imageWatermarkInput := mps.ImageWatermarkInput{}
		if v, ok := dMap["image_content"]; ok {
			imageWatermarkInput.ImageContent = helper.String(v.(string))
		}
		if v, ok := dMap["width"]; ok {
			imageWatermarkInput.Width = helper.String(v.(string))
		}
		if v, ok := dMap["height"]; ok {
			imageWatermarkInput.Height = helper.String(v.(string))
		}
		if v, ok := dMap["repeat_type"]; ok {
			imageWatermarkInput.RepeatType = helper.String(v.(string))
		}
		request.ImageTemplate = &imageWatermarkInput
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "text_template"); ok {
		textWatermarkTemplateInput := mps.TextWatermarkTemplateInput{}
		if v, ok := dMap["font_type"]; ok {
			textWatermarkTemplateInput.FontType = helper.String(v.(string))
		}
		if v, ok := dMap["font_size"]; ok {
			textWatermarkTemplateInput.FontSize = helper.String(v.(string))
		}
		if v, ok := dMap["font_color"]; ok {
			textWatermarkTemplateInput.FontColor = helper.String(v.(string))
		}
		if v, ok := dMap["font_alpha"]; ok {
			textWatermarkTemplateInput.FontAlpha = helper.Float64(v.(float64))
		}
		request.TextTemplate = &textWatermarkTemplateInput
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "svg_template"); ok {
		svgWatermarkInput := mps.SvgWatermarkInput{}
		if v, ok := dMap["width"]; ok {
			svgWatermarkInput.Width = helper.String(v.(string))
		}
		if v, ok := dMap["height"]; ok {
			svgWatermarkInput.Height = helper.String(v.(string))
		}
		request.SvgTemplate = &svgWatermarkInput
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().CreateWatermarkTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps watermarkTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition = *response.Response.Definition
	d.SetId(helper.Int64ToStr(definition))

	return resourceTencentCloudMpsWatermarkTemplateRead(d, meta)
}

func resourceTencentCloudMpsWatermarkTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_watermark_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	definition := d.Id()

	watermarkTemplate, err := service.DescribeMpsWatermarkTemplateById(ctx, definition)
	if err != nil {
		return err
	}

	if watermarkTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsWatermarkTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if watermarkTemplate.Type != nil {
		_ = d.Set("type", watermarkTemplate.Type)
	}

	if watermarkTemplate.Name != nil {
		_ = d.Set("name", watermarkTemplate.Name)
	}

	if watermarkTemplate.Comment != nil {
		_ = d.Set("comment", watermarkTemplate.Comment)
	}

	if watermarkTemplate.CoordinateOrigin != nil {
		_ = d.Set("coordinate_origin", watermarkTemplate.CoordinateOrigin)
	}

	if watermarkTemplate.XPos != nil {
		_ = d.Set("x_pos", watermarkTemplate.XPos)
	}

	if watermarkTemplate.YPos != nil {
		_ = d.Set("y_pos", watermarkTemplate.YPos)
	}

	if watermarkTemplate.ImageTemplate != nil {
		imageTemplateMap := map[string]interface{}{}

		if watermarkTemplate.ImageTemplate.ImageUrl != nil {
			url := watermarkTemplate.ImageTemplate.ImageUrl
			res, err := http.Get(*url)
			if err != nil {
				return err
			}
			content, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return err
			}
			base64Encode := base64.StdEncoding.EncodeToString(content)
			imageTemplateMap["image_content"] = base64Encode
		}

		if watermarkTemplate.ImageTemplate.Width != nil {
			imageTemplateMap["width"] = watermarkTemplate.ImageTemplate.Width
		}

		if watermarkTemplate.ImageTemplate.Height != nil {
			imageTemplateMap["height"] = watermarkTemplate.ImageTemplate.Height
		}

		if watermarkTemplate.ImageTemplate.RepeatType != nil {
			imageTemplateMap["repeat_type"] = watermarkTemplate.ImageTemplate.RepeatType
		}

		_ = d.Set("image_template", []interface{}{imageTemplateMap})
	}

	if watermarkTemplate.TextTemplate != nil {
		textTemplateMap := map[string]interface{}{}

		if watermarkTemplate.TextTemplate.FontType != nil {
			textTemplateMap["font_type"] = watermarkTemplate.TextTemplate.FontType
		}

		if watermarkTemplate.TextTemplate.FontSize != nil {
			textTemplateMap["font_size"] = watermarkTemplate.TextTemplate.FontSize
		}

		if watermarkTemplate.TextTemplate.FontColor != nil {
			textTemplateMap["font_color"] = watermarkTemplate.TextTemplate.FontColor
		}

		if watermarkTemplate.TextTemplate.FontAlpha != nil {
			textTemplateMap["font_alpha"] = watermarkTemplate.TextTemplate.FontAlpha
		}

		_ = d.Set("text_template", []interface{}{textTemplateMap})
	}

	if watermarkTemplate.SvgTemplate != nil {
		svgTemplateMap := map[string]interface{}{}

		if watermarkTemplate.SvgTemplate.Width != nil {
			svgTemplateMap["width"] = watermarkTemplate.SvgTemplate.Width
		}

		if watermarkTemplate.SvgTemplate.Height != nil {
			svgTemplateMap["height"] = watermarkTemplate.SvgTemplate.Height
		}

		_ = d.Set("svg_template", []interface{}{svgTemplateMap})
	}

	return nil
}

func resourceTencentCloudMpsWatermarkTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_watermark_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := mps.NewModifyWatermarkTemplateRequest()

	definition := d.Id()

	request.Definition = helper.StrToInt64Point(definition)

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}
	}

	if d.HasChange("coordinate_origin") {
		if v, ok := d.GetOk("coordinate_origin"); ok {
			request.CoordinateOrigin = helper.String(v.(string))
		}
	}

	if d.HasChange("x_pos") {
		if v, ok := d.GetOk("x_pos"); ok {
			request.XPos = helper.String(v.(string))
		}
	}

	if d.HasChange("y_pos") {
		if v, ok := d.GetOk("y_pos"); ok {
			request.YPos = helper.String(v.(string))
		}
	}

	if d.HasChange("image_template") {
		if dMap, ok := helper.InterfacesHeadMap(d, "image_template"); ok {
			imageWatermarkInput := mps.ImageWatermarkInputForUpdate{}
			if v, ok := dMap["image_content"]; ok {
				imageWatermarkInput.ImageContent = helper.String(v.(string))
			}
			if v, ok := dMap["width"]; ok {
				imageWatermarkInput.Width = helper.String(v.(string))
			}
			if v, ok := dMap["height"]; ok {
				imageWatermarkInput.Height = helper.String(v.(string))
			}
			if v, ok := dMap["repeat_type"]; ok {
				imageWatermarkInput.RepeatType = helper.String(v.(string))
			}
			request.ImageTemplate = &imageWatermarkInput
		}
	}

	if d.HasChange("text_template") {
		if dMap, ok := helper.InterfacesHeadMap(d, "text_template"); ok {
			textWatermarkTemplateInput := mps.TextWatermarkTemplateInputForUpdate{}
			if v, ok := dMap["font_type"]; ok {
				textWatermarkTemplateInput.FontType = helper.String(v.(string))
			}
			if v, ok := dMap["font_size"]; ok {
				textWatermarkTemplateInput.FontSize = helper.String(v.(string))
			}
			if v, ok := dMap["font_color"]; ok {
				textWatermarkTemplateInput.FontColor = helper.String(v.(string))
			}
			if v, ok := dMap["font_alpha"]; ok {
				textWatermarkTemplateInput.FontAlpha = helper.Float64(v.(float64))
			}
			request.TextTemplate = &textWatermarkTemplateInput
		}
	}

	if d.HasChange("svg_template") {
		if dMap, ok := helper.InterfacesHeadMap(d, "svg_template"); ok {
			svgWatermarkInput := mps.SvgWatermarkInputForUpdate{}
			if v, ok := dMap["width"]; ok {
				svgWatermarkInput.Width = helper.String(v.(string))
			}
			if v, ok := dMap["height"]; ok {
				svgWatermarkInput.Height = helper.String(v.(string))
			}
			request.SvgTemplate = &svgWatermarkInput
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().ModifyWatermarkTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mps watermarkTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMpsWatermarkTemplateRead(d, meta)
}

func resourceTencentCloudMpsWatermarkTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_watermark_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	definition := d.Id()

	if err := service.DeleteMpsWatermarkTemplateById(ctx, definition); err != nil {
		return err
	}

	return nil
}
