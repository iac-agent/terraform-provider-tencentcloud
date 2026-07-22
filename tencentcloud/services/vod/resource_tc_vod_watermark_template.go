package vod

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVodWatermarkTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVodWatermarkTemplateCreate,
		Read:   resourceTencentCloudVodWatermarkTemplateRead,
		Update: resourceTencentCloudVodWatermarkTemplateUpdate,
		Delete: resourceTencentCloudVodWatermarkTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Watermarking 类型 有效值：镜像: 镜像 水印; text: text 水印; svg: SVG 水印。",
			},

			"sub_app_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "VOD [应用](https://intl.云.tencent.com/document/product/266/14574) ID. For customers who activate VOD 服务 从 December 25，2023，如果 they want 到 访问 resources 在 VOD 应用 (whether 它's 默认值 应用 或 newly 创建 一个)，they 必须 fill 在 此 字段 使用 应用 ID。",
			},

			"name": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Watermarking 模板名称 Length 限制: 64 字符。",
			},

			"comment": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "模板描述 Length 限制: 256 字符。",
			},

			"coordinate_origin": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Origin position. 有效值：TopLeft: 源站 的 coordinates 是 在 top-left corner 的 视频，和 源站 的 水印 是 在 top-left corner 的 镜像 或 text; TopRight: 源站 的 coordinates 是 在 top-right corner 的 视频，和 源站 的 水印 是 在 top-right corner 的 镜像 或 text; BottomLeft: 源站 的 coordinates 是 在 bottom-left corner 的 视频，和 源站 的 水印 是 在 bottom-left corner 的 镜像 或 text; BottomRight: 源站 的 coordinates 是 在 bottom-right corner 的 视频，和 源站 的 水印 是 在 bottom-right corner 的 镜像 或 text.默认值：TopLeft。",
			},

			"x_pos": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "horizontal position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `XPos` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `XPos` 是 10% 的 视频 宽度; 如果 字符串 结束 在 像素， `XPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `XPos` 是 100 像素.默认值：0 像素。",
			},

			"y_pos": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "vertical position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `YPos` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `YPos` 是 10% 的 视频 高度; 如果 字符串 结束 在 像素， `YPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `YPos` 是 100 像素.默认值：0 像素。",
			},

			"image_template": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Image watermarking template. 此 字段 为必填项 当 `类型` 是 `镜像` 和 是 无效 当 `类型` 是 `text`。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "[Base64](https://tools.ietf.org/html/rfc4648) encoded 字符串 的 水印 镜像. Only JPEG，PNG，和 GIF images 是 支持。",
						},
						"width": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Watermark 宽度. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `宽度` 的 水印 将 是 指定 percentage 的 视频 宽度. For 示例，`10%` 表示 该 `宽度` 是 10% 的 视频 宽度; 如果 字符串 结束 在 像素， `宽度` 的 水印 将 是 在 pixels. For 示例，`100px` 表示 该 `宽度` 是 100 pixels. 取值范围：[8，4096]. 默认值：10%。",
						},
						"height": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Watermark 高度. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `高度` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `高度` 是 10% 的 视频 高度; 如果 字符串 结束 在 像素， `高度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `高度` 是 100 像素. 有效值：0 或 [8,4096]. 默认值：0 像素，其中 表示 该 `高度` 将 是 proportionally scaled according 到 aspect ratio 的 original 水印 镜像。",
						},
						"repeat_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Repeat 类型 animated 水印. 有效值：once: 无 longer appears after 水印 playback 结束. repeat_last_frame: stays 在 last frame after 水印 playback 结束. repeat (默认值): repeats playback until 视频 结束。",
						},
						"transparency": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Image 水印 transparency: 0: completely opaque 100: completely transparent 默认值：0。",
						},
					},
				},
			},

			"text_template": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Text watermarking template. 此 字段 为必填项 当 `类型` 是 `text` 和 是 无效 当 `类型` 是 `镜像`。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"font_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Font 类型 Currently，two types 是 支持: simkai.ttf: both Chinese 和 English 是 支持; arial.ttf: 仅 English 是 支持。",
						},
						"font_size": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Font 大小 在 Npx 格式 其中 N 是 numeric 值",
						},
						"font_color": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Font color 在 0xRRGGBB 格式 默认值：0xFFFFFF (white)。",
						},
						"font_alpha": {
							Type:        schema.TypeFloat,
							Required:    true,
							Description: "Text transparency. 取值范围：(0，1] 0: completely transparent 1: completely opaque 默认值：1。",
						},
					},
				},
			},

			"svg_template": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "SVG watermarking template. 此 字段 为必填项 当 `类型` 是 `svg` 和 是 无效 当 `类型` 是 `镜像` 或 `text`。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"width": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Watermark 宽度，其中 支持 six formats 的 像素，%，W%，H%，S%，和 L%: 如果 字符串 结束 在 像素， `宽度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `宽度` 是 100 像素; 如果 `0px` 是 entered 和 `高度` 是 不 `0px`， 水印 宽度 将 是 proportionally scaled based 在 来源 SVG 镜像; 如果 `0px` 是 entered 对于 both `宽度` 和 `高度`， 水印 宽度 将 是 宽度 的 来源 SVG 镜像; 如果 字符串 结束 在 `W%`， `宽度` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10W%` 表示 该 `宽度` 是 10% 的 视频 宽度; 如果 字符串 结束 在 `H%`， `宽度` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10H%` 表示 该 `宽度` 是 10% 的 视频 高度; 如果 字符串 结束 在 `S%`， `宽度` 的 水印 将 是 指定 percentage 的 short side 的 视频; 对于 示例，`10S%` 表示 该 `宽度` 是 10% 的 short side 的 视频; 如果 字符串 结束 在 `L%`， `宽度` 的 水印 将 是 指定 percentage 的 long side 的 视频; 对于 示例，`10L%` 表示 该 `宽度` 是 10% 的 long side 的 视频; 如果 字符串 结束 在 %， meaning 是 same 作为 `W%`. 默认值：10W%。",
						},
						"height": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Watermark 高度，其中 支持 six formats 的 像素，%，W%，H%，S%，和 L%: 如果 字符串 结束 在 像素， `高度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `高度` 是 100 像素; 如果 `0px` 是 entered 和 `宽度` 是 不 `0px`， 水印 高度 将 是 proportionally scaled based 在 来源 SVG 镜像; 如果 `0px` 是 entered 对于 both `宽度` 和 `高度`， 水印 高度 将 是 高度 的 来源 SVG 镜像; 如果 字符串 结束 在 `W%`， `高度` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10W%` 表示 该 `高度` 是 10% 的 视频 宽度; 如果 字符串 结束 在 `H%`， `高度` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10H%` 表示 该 `高度` 是 10% 的 视频 高度; 如果 字符串 结束 在 `S%`， `高度` 的 水印 将 是 指定 percentage 的 short side 的 视频; 对于 示例，`10S%` 表示 该 `高度` 是 10% 的 short side 的 视频; 如果 字符串 结束 在 `L%`， `高度` 的 水印 将 是 指定 percentage 的 long side 的 视频; 对于 示例，`10L%` 表示 该 `高度` 是 10% 的 long side 的 视频; 如果 字符串 结束 在 %， meaning 是 same 作为 `H%`. 默认值：0 像素。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudVodWatermarkTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_watermark_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = vod.NewCreateWatermarkTemplateRequest()
		response = vod.NewCreateWatermarkTemplateResponse()
		subAppId string
	)
	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("sub_app_id"); ok {
		subAppId = helper.IntToStr(v.(int))
		request.SubAppId = helper.IntUint64(v.(int))
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
		imageWatermarkInput := vod.ImageWatermarkInput{}
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
		if v, ok := dMap["transparency"]; ok {
			imageWatermarkInput.Transparency = helper.IntInt64(v.(int))
		}
		request.ImageTemplate = &imageWatermarkInput
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "text_template"); ok {
		textWatermarkTemplateInput := vod.TextWatermarkTemplateInput{}
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
		svgWatermarkInput := vod.SvgWatermarkInput{}
		if v, ok := dMap["width"]; ok {
			svgWatermarkInput.Width = helper.String(v.(string))
		}
		if v, ok := dMap["height"]; ok {
			svgWatermarkInput.Height = helper.String(v.(string))
		}
		request.SvgTemplate = &svgWatermarkInput
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().CreateWatermarkTemplate(request)
		if e != nil {
			if sdkError, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "FailedOperation" && sdkError.Message == "invalid vod user" {
					return resource.RetryableError(e)
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), e.Error())
			return resource.NonRetryableError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vod watermarkTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition := *response.Response.Definition
	d.SetId(subAppId + tccommon.FILED_SP + helper.Int64ToStr(definition))

	return resourceTencentCloudVodWatermarkTemplateRead(d, meta)
}

func resourceTencentCloudVodWatermarkTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_watermark_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("watermark template id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[0]
	definition := idSplit[1]

	watermarkTemplate, err := service.DescribeVodWatermarkTemplateById(ctx, *helper.StrToUint64Point(subAppId), helper.StrToInt64(definition))
	if err != nil {
		return err
	}

	if watermarkTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VodWatermarkTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("sub_app_id", helper.StrToInt(subAppId))

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
			imageContentResp, err := http.Get(*watermarkTemplate.ImageTemplate.ImageUrl)
			if err != nil {
				return err
			}
			content, err := ioutil.ReadAll(imageContentResp.Body)
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

		if watermarkTemplate.ImageTemplate.Transparency != nil {
			imageTemplateMap["transparency"] = watermarkTemplate.ImageTemplate.Transparency
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

func resourceTencentCloudVodWatermarkTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_watermark_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vod.NewModifyWatermarkTemplateRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("watermark template id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[0]
	definition := idSplit[1]

	request.SubAppId = helper.StrToUint64Point(subAppId)
	request.Definition = helper.StrToInt64Point(definition)

	immutableArgs := []string{"type", "sub_app_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

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
			imageWatermarkInput := vod.ImageWatermarkInputForUpdate{}
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
			if v, ok := dMap["transparency"]; ok {
				imageWatermarkInput.Transparency = helper.IntInt64(v.(int))
			}
			request.ImageTemplate = &imageWatermarkInput
		}
	}

	if d.HasChange("text_template") {
		if dMap, ok := helper.InterfacesHeadMap(d, "text_template"); ok {
			textWatermarkTemplateInput := vod.TextWatermarkTemplateInputForUpdate{}
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
			svgWatermarkInput := vod.SvgWatermarkInputForUpdate{}
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
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().ModifyWatermarkTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update vod watermarkTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudVodWatermarkTemplateRead(d, meta)
}

func resourceTencentCloudVodWatermarkTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_watermark_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("watermark template id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[0]
	definition := idSplit[1]

	if err := service.DeleteVodWatermarkTemplateById(ctx, helper.StrToUInt64(subAppId), helper.StrToInt64(definition)); err != nil {
		return err
	}

	return nil
}
