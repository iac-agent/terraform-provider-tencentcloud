package ci

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCiMediaWatermarkTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaWatermarkTemplateCreate,
		Read:   resourceTencentCloudCiMediaWatermarkTemplateRead,
		Update: resourceTencentCloudCiMediaWatermarkTemplateUpdate,
		Delete: resourceTencentCloudCiMediaWatermarkTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "存储桶名称",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "模板名称 仅 支持 `Chinese`，`English`，`numbers`，`_`，`-` 和 `*`。",
			},

			"watermark": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "容器 格式",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Watermark 类型，Text: text 水印，Image: 镜像 水印。",
						},
						"pos": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Reference position，TopRight，TopLeft，BottomRight，BottomLeft，Left，Right，Top，Bottom，Center。",
						},
						"loc_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "偏移量 方法，Relativity: proportional，Absolute: fixed position。",
						},
						"dx": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Horizontal 偏移量，1: In picture 水印，如果 Background 是 true，当 locMode 是 Relativity，它 是 %，取值范围：[-300 0]; 当 locMode 是 Absolute，它 是 像素，取值范围：[-4096 0] ]，2: In picture 水印，如果 Background 是 false，当 locMode 是 Relativity，它 是 %，取值范围：[0 100]; 当 locMode 是 Absolute，它 是 像素，取值范围：[0 4096]，3: In text 水印，当 locMode 是 Relativity，它 是 %，取值范围：[0 100]; 当 locMode 是 Absolute，它 是 像素，取值范围：[0 4096]，4: 当 Pos 是 Top，Bottom 和 Center， 参数 是 无效。",
						},
						"dy": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Vertical 偏移量，1: In picture 水印，如果 Background 是 true，当 locMode 是 Relativity，它 是 %，取值范围：[-300 0]; 当 locMode 是 Absolute，它 是 像素，取值范围：[-4096 0] ],2: In picture 水印，如果 Background 是 false，当 locMode 是 Relativity，它 是 %，取值范围：[0 100]; 当 locMode 是 Absolute，它 是 像素，取值范围：[0 4096],3: In text 水印，当 locMode 是 Relativity，它 是 %，取值范围：[0 100]; 当 locMode 是 Absolute，它 是 像素，取值范围：[0 4096]，4: 当 Pos 是 Left，Right 和 Center， 参数 是 无效。",
						},
						"start_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Watermark 开始时间，1: [0 视频 时长]，2: 单位 是 second，3: support float 格式，execution accuracy 是 accurate 到 milliseconds。",
						},
						"end_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Watermark 结束时间，1: [0 视频 时长]，2: 单位 是 second，3: support float 格式，execution accuracy 是 accurate 到 milliseconds。",
						},

						"image": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Image 水印 节点。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "地址 的 水印 map (pass 在 after Urlencode 为必填项)。",
									},
									"mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Size 模式，Original: original 大小，Proportion: proportional，Fixed: fixed 大小。",
									},
									"width": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "宽度，1: 当 模式 是 Original，它 does 不 support setting 宽度 的 水印 镜像，2: 当 模式 是 Proportion， 单位 是 %， 值 范围 的 background 镜像: [100 300]; 值 范围 的 foreground 镜像: [1 100]，relative 到 Video 宽度，up 到 4096px，3: 当 模式 是 Fixed， 单位 是 像素，取值范围：[8，4096]，4: 如果 仅 宽度 是 集合，高度 是 calculated according 到 proportion 的 水印 镜像。",
									},
									"height": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "High，1: 当 模式 是 Original，它 does 不 support setting 宽度 的 水印 镜像，2: 当 模式 是 Proportion， 单位 是 %， 值 范围 的 background 镜像: [100 300]; 值 范围 的 foreground 镜像: [1 100]，relative 到 Video 宽度，up 到 4096px，3: 当 模式 是 Fixed， 单位 是 像素，取值范围：[8，4096]，4: 如果 仅 宽度 是 集合，高度 是 calculated according 到 proportion 的 水印 镜像。",
									},
									"transparency": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Transparency，取值范围：[1 100]，单位 %。",
									},
									"background": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "是否background 镜像。",
									},
								},
							},
						},
						"text": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Text Watermark Node。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"font_size": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Font 大小，取值范围：[5 100]，单位 像素。",
									},
									"font_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "font 类型",
									},
									"font_color": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Font color，格式: 0xRRGGBB。",
									},
									"transparency": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Transparency，取值范围：[1 100]，单位 %。",
									},
									"text": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Watermark 内容， 长度 does 不 exceed 64 字符，仅 支持 Chinese，English，numbers，_，- 和 *。",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCiMediaWatermarkTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_watermark_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request = cos.CreateMediaWatermarkTemplateOptions{
			Tag: "Watermark",
		}
		bucket     string
		templateId string
	)

	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	} else {
		return errors.New("get bucket failed!")
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "watermark"); ok {
		watermark := cos.Watermark{}
		if v, ok := dMap["type"]; ok {
			watermark.Type = v.(string)
		}
		if v, ok := dMap["pos"]; ok {
			watermark.Pos = v.(string)
		}
		if v, ok := dMap["loc_mode"]; ok {
			watermark.LocMode = v.(string)
		}
		if v, ok := dMap["dx"]; ok {
			watermark.Dx = v.(string)
		}
		if v, ok := dMap["dy"]; ok {
			watermark.Dy = v.(string)
		}
		if v, ok := dMap["start_time"]; ok {
			watermark.StartTime = v.(string)
		}
		if v, ok := dMap["end_time"]; ok {
			watermark.EndTime = v.(string)
		}
		if imageMap, ok := helper.InterfaceToMap(dMap, "image"); ok {
			image := cos.Image{}
			if v, ok := imageMap["url"]; ok {
				image.Url = v.(string)
			}
			if v, ok := imageMap["mode"]; ok {
				image.Mode = v.(string)
			}
			if v, ok := imageMap["width"]; ok {
				image.Width = v.(string)
			}
			if v, ok := imageMap["height"]; ok {
				image.Height = v.(string)
			}
			if v, ok := imageMap["transparency"]; ok {
				image.Transparency = v.(string)
			}
			if v, ok := imageMap["background"]; ok {
				image.Background = v.(string)
			}
			watermark.Image = &image
		}
		if textMap, ok := helper.InterfaceToMap(dMap, "text"); ok {
			text := cos.Text{}
			if v, ok := textMap["font_size"]; ok {
				text.FontSize = v.(string)
			}
			if v, ok := textMap["font_type"]; ok {
				text.FontType = v.(string)
			}
			if v, ok := textMap["font_color"]; ok {
				text.FontColor = v.(string)
			}
			if v, ok := textMap["transparency"]; ok {
				text.Transparency = v.(string)
			}
			if v, ok := textMap["text"]; ok {
				text.Text = v.(string)
			}
			watermark.Text = &text
		}
		request.Watermark = &watermark
	}

	var response *cos.CreateMediaTemplateResult
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket).CI.CreateMediaWatermarkTemplate(ctx, &request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "CreateMediaWatermarkTemplate", request, result)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaWatermarkTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = response.Template.TemplateId
	d.SetId(bucket + tccommon.FILED_SP + templateId)

	return resourceTencentCloudCiMediaWatermarkTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaWatermarkTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_watermark_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	mediaWatermarkTemplate, err := service.DescribeCiMediaTemplateById(ctx, bucket, templateId)
	if err != nil {
		return err
	}

	if mediaWatermarkTemplate == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	_ = d.Set("bucket", bucket)

	if mediaWatermarkTemplate.Name != "" {
		_ = d.Set("name", mediaWatermarkTemplate.Name)
	}

	if mediaWatermarkTemplate.Watermark != nil {
		watermarkMap := map[string]interface{}{}

		if mediaWatermarkTemplate.Watermark.Type != "" {
			watermarkMap["type"] = mediaWatermarkTemplate.Watermark.Type
		}

		if mediaWatermarkTemplate.Watermark.Pos != "" {
			watermarkMap["pos"] = mediaWatermarkTemplate.Watermark.Pos
		}

		if mediaWatermarkTemplate.Watermark.LocMode != "" {
			watermarkMap["loc_mode"] = mediaWatermarkTemplate.Watermark.LocMode
		}

		if mediaWatermarkTemplate.Watermark.Dx != "" {
			watermarkMap["dx"] = mediaWatermarkTemplate.Watermark.Dx
		}

		if mediaWatermarkTemplate.Watermark.Dy != "" {
			watermarkMap["dy"] = mediaWatermarkTemplate.Watermark.Dy
		}

		if mediaWatermarkTemplate.Watermark.StartTime != "" {
			watermarkMap["start_time"] = mediaWatermarkTemplate.Watermark.StartTime
		}

		if mediaWatermarkTemplate.Watermark.EndTime != "" {
			watermarkMap["end_time"] = mediaWatermarkTemplate.Watermark.EndTime
		}

		if mediaWatermarkTemplate.Watermark.Image != nil {
			imageMap := map[string]interface{}{}

			if mediaWatermarkTemplate.Watermark.Image.Url != "" {
				imageMap["url"] = mediaWatermarkTemplate.Watermark.Image.Url
			}

			if mediaWatermarkTemplate.Watermark.Image.Mode != "" {
				imageMap["mode"] = mediaWatermarkTemplate.Watermark.Image.Mode
			}

			if mediaWatermarkTemplate.Watermark.Image.Width != "" {
				imageMap["width"] = mediaWatermarkTemplate.Watermark.Image.Width
			}

			if mediaWatermarkTemplate.Watermark.Image.Height != "" {
				imageMap["height"] = mediaWatermarkTemplate.Watermark.Image.Height
			}

			if mediaWatermarkTemplate.Watermark.Image.Transparency != "" {
				imageMap["transparency"] = mediaWatermarkTemplate.Watermark.Image.Transparency
			}

			if mediaWatermarkTemplate.Watermark.Image.Background != "" {
				imageMap["background"] = mediaWatermarkTemplate.Watermark.Image.Background
			}

			watermarkMap["image"] = []interface{}{imageMap}
		}

		if mediaWatermarkTemplate.Watermark.Text != nil {
			textMap := map[string]interface{}{}

			if mediaWatermarkTemplate.Watermark.Text.FontSize != "" {
				textMap["font_size"] = mediaWatermarkTemplate.Watermark.Text.FontSize
			}

			if mediaWatermarkTemplate.Watermark.Text.FontType != "" {
				textMap["font_type"] = mediaWatermarkTemplate.Watermark.Text.FontType
			}

			if mediaWatermarkTemplate.Watermark.Text.FontColor != "" {
				textMap["font_color"] = mediaWatermarkTemplate.Watermark.Text.FontColor
			}

			if mediaWatermarkTemplate.Watermark.Text.Transparency != "" {
				textMap["transparency"] = mediaWatermarkTemplate.Watermark.Text.Transparency
			}

			if mediaWatermarkTemplate.Watermark.Text.Text != "" {
				textMap["text"] = mediaWatermarkTemplate.Watermark.Text.Text
			}

			watermarkMap["text"] = []interface{}{textMap}
		}

		_ = d.Set("watermark", []interface{}{watermarkMap})
	}

	return nil
}

func resourceTencentCloudCiMediaWatermarkTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_watermark_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := cos.CreateMediaWatermarkTemplateOptions{
		Tag: "Watermark",
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if d.HasChange("watermark") {
		if dMap, ok := helper.InterfacesHeadMap(d, "watermark"); ok {
			watermark := cos.Watermark{}
			if v, ok := dMap["type"]; ok {
				watermark.Type = v.(string)
			}
			if v, ok := dMap["pos"]; ok {
				watermark.Pos = v.(string)
			}
			if v, ok := dMap["loc_mode"]; ok {
				watermark.LocMode = v.(string)
			}
			if v, ok := dMap["dx"]; ok {
				watermark.Dx = v.(string)
			}
			if v, ok := dMap["dy"]; ok {
				watermark.Dy = v.(string)
			}
			if v, ok := dMap["start_time"]; ok {
				watermark.StartTime = v.(string)
			}
			if v, ok := dMap["end_time"]; ok {
				watermark.EndTime = v.(string)
			}

			if imageMap, ok := helper.InterfaceToMap(dMap, "image"); ok {
				image := cos.Image{}
				if v, ok := imageMap["url"]; ok {
					image.Url = v.(string)
				}
				if v, ok := imageMap["mode"]; ok {
					image.Mode = v.(string)
				}
				if v, ok := imageMap["width"]; ok {
					image.Width = v.(string)
				}
				if v, ok := imageMap["height"]; ok {
					image.Height = v.(string)
				}
				if v, ok := imageMap["transparency"]; ok {
					image.Transparency = v.(string)
				}
				if v, ok := imageMap["background"]; ok {
					image.Background = v.(string)
				}
				watermark.Image = &image
			}
			if textMap, ok := helper.InterfaceToMap(dMap, "text"); ok {
				text := cos.Text{}
				if v, ok := textMap["font_size"]; ok {
					text.FontSize = v.(string)
				}
				if v, ok := textMap["font_type"]; ok {
					text.FontType = v.(string)
				}
				if v, ok := textMap["font_color"]; ok {
					text.FontColor = v.(string)
				}
				if v, ok := textMap["transparency"]; ok {
					text.Transparency = v.(string)
				}
				if v, ok := textMap["text"]; ok {
					text.Text = v.(string)
				}
				watermark.Text = &text
			}
			request.Watermark = &watermark
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket).CI.UpdateMediaWatermarkTemplate(ctx, &request, templateId)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaWatermarkTemplate", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaWatermarkTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaWatermarkTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaWatermarkTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_watermark_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	if err := service.DeleteCiMediaTemplateById(ctx, bucket, templateId); err != nil {
		return err
	}

	return nil
}
