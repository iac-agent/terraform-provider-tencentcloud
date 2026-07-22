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

func ResourceTencentCloudCiMediaAnimationTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaAnimationTemplateCreate,
		Read:   resourceTencentCloudCiMediaAnimationTemplateRead,
		Update: resourceTencentCloudCiMediaAnimationTemplateUpdate,
		Delete: resourceTencentCloudCiMediaAnimationTemplateDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
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

			"container": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "容器 格式",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"format": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Package 格式",
						},
					},
				},
			},

			"video": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "视频 信息，do 不 upload Video，其中 是 equivalent 到 deleting 视频 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"codec": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Codec 格式 `gif`，`webp`。",
						},
						"width": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "宽度，取值范围：[128，4096]，单位：像素，如果 仅 宽度 是 集合，高度 是 calculated according 到 original ratio 的 视频，必须 是 even。",
						},
						"height": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "High，取值范围：[128，4096]，单位：像素，如果 仅 高度 是 集合，宽度 是 calculated according 到 original ratio 的 视频，必须 是 even。",
						},
						"fps": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Frame 速率，取值范围：(0，60]，单位：fps。",
						},
						"animate_only_keep_key_frame": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "GIFs 是 kept 仅 Keyframe，优先级: AnimateFramesPerSecond &gt; AnimateOnlyKeepKeyFrame &gt; AnimateTimeIntervalOfFrame。",
						},
						"animate_time_interval_of_frame": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Animation frame extraction every 时间，(0，视频 时长]，Animation frame extraction 时间间隔，如果 TimeInterval.Duration 是 集合，它 是 less 比 此 值",
						},
						"animate_frames_per_second": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Animation per second frame 数量，优先级: AnimateFramesPerSecond &gt; AnimateOnlyKeepKeyFrame &gt; AnimateTimeIntervalOfFrame。",
						},
						"quality": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Set relative quality，[1，100)，webp 镜像 quality setting takes effect，gif has 无 quality 参数。",
						},
					},
				},
			},

			"time_interval": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "时间间隔。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Starting 时间，[0 视频 时长]，（秒）， Support float 格式， execution accuracy 是 accurate 到 milliseconds。",
						},
						"duration": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "时长，[0 视频 时长]，（秒）， Support float 格式， execution accuracy 是 accurate 到 milliseconds。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCiMediaAnimationTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_animation_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request = cos.CreateMediaAnimationTemplateOptions{
			Tag: "Animation",
		}
		templateId string
		bucket     string
	)

	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	} else {
		return errors.New("get bucket failed!")
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "container"); ok {
		container := cos.Container{}
		if v, ok := dMap["format"]; ok {
			container.Format = v.(string)
		}
		request.Container = &container
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "video"); ok {
		video := cos.AnimationVideo{}
		if v, ok := dMap["codec"]; ok {
			video.Codec = v.(string)
		}
		if v, ok := dMap["width"]; ok {
			video.Width = v.(string)
		}
		if v, ok := dMap["height"]; ok {
			video.Height = v.(string)
		}
		if v, ok := dMap["fps"]; ok {
			video.Fps = v.(string)
		}
		if v, ok := dMap["animate_only_keep_key_frame"]; ok {
			video.AnimateOnlyKeepKeyFrame = v.(string)
		}
		if v, ok := dMap["animate_time_interval_of_frame"]; ok {
			video.AnimateTimeIntervalOfFrame = v.(string)
		}
		if v, ok := dMap["animate_frames_per_second"]; ok {
			video.AnimateFramesPerSecond = v.(string)
		}
		if v, ok := dMap["quality"]; ok {
			video.Quality = v.(string)
		}
		request.Video = &video
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "time_interval"); ok {
		timeInterval := cos.TimeInterval{}
		if v, ok := dMap["start"]; ok {
			timeInterval.Start = v.(string)
		}
		if v, ok := dMap["duration"]; ok {
			timeInterval.Duration = v.(string)
		}
		request.TimeInterval = &timeInterval
	}

	var response *cos.CreateMediaTemplateResult
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket).CI.CreateMediaAnimationTemplate(ctx, &request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%+v], response body [%+v]\n", logId, "CreateMediaAnimationTemplate", request, result)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaAnimationTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = response.Template.TemplateId
	d.SetId(bucket + tccommon.FILED_SP + templateId)

	return resourceTencentCloudCiMediaAnimationTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaAnimationTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_animation_template.read")()
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

	template, err := service.DescribeCiMediaTemplateById(ctx, bucket, templateId)
	if err != nil {
		return err
	}

	if template == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	_ = d.Set("bucket", bucket)

	if template.Name != "" {
		_ = d.Set("name", template.Name)
	}

	mediaAnimationTemplate := template.TransTpl
	if mediaAnimationTemplate != nil {
		if mediaAnimationTemplate.Container != nil {
			containerMap := map[string]interface{}{}

			if mediaAnimationTemplate.Container.Format != "" {
				containerMap["format"] = mediaAnimationTemplate.Container.Format
			}

			_ = d.Set("container", []interface{}{containerMap})
		}

		if mediaAnimationTemplate.Video != nil {
			videoMap := map[string]interface{}{}

			if mediaAnimationTemplate.Video.Codec != "" {
				videoMap["codec"] = mediaAnimationTemplate.Video.Codec
			}

			if mediaAnimationTemplate.Video.Width != "" {
				videoMap["width"] = mediaAnimationTemplate.Video.Width
			}

			if mediaAnimationTemplate.Video.Height != "" {
				videoMap["height"] = mediaAnimationTemplate.Video.Height
			}

			if mediaAnimationTemplate.Video.Fps != "" {
				videoMap["fps"] = mediaAnimationTemplate.Video.Fps
			}

			// if mediaAnimationTemplate.Video.AnimateOnlyKeepKeyFrame != "" {
			// 	videoMap["animate_only_keep_key_frame"] = mediaAnimationTemplate.Video.AnimateOnlyKeepKeyFrame
			// }

			// if mediaAnimationTemplate.Video.AnimateTimeIntervalOfFrame != "" {
			// 	videoMap["animate_time_interval_of_frame"] = mediaAnimationTemplate.Video.AnimateTimeIntervalOfFrame
			// }

			// if mediaAnimationTemplate.Video.AnimateFramesPerSecond != "" {
			// 	videoMap["animate_frames_per_second"] = mediaAnimationTemplate.Video.AnimateFramesPerSecond
			// }

			// if mediaAnimationTemplate.Video.Quality != "" {
			// 	videoMap["quality"] = mediaAnimationTemplate.Video.Quality
			// }

			err = d.Set("video", []interface{}{videoMap})
			if err != nil {
				return err
			}
		}

		if mediaAnimationTemplate.TimeInterval != nil {
			timeIntervalMap := map[string]interface{}{}

			if mediaAnimationTemplate.TimeInterval.Start != "" {
				timeIntervalMap["start"] = mediaAnimationTemplate.TimeInterval.Start
			}

			if mediaAnimationTemplate.TimeInterval.Duration != "" {
				timeIntervalMap["duration"] = mediaAnimationTemplate.TimeInterval.Duration
			}

			err = d.Set("time_interval", []interface{}{timeIntervalMap})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceTencentCloudCiMediaAnimationTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_animation_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := cos.CreateMediaAnimationTemplateOptions{
		Tag: "Animation",
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "container"); ok {
		container := cos.Container{}
		if v, ok := dMap["format"]; ok {
			container.Format = v.(string)
		}
		request.Container = &container
	}

	if d.HasChange("video") {
		if dMap, ok := helper.InterfacesHeadMap(d, "video"); ok {
			video := cos.AnimationVideo{}
			if v, ok := dMap["codec"]; ok {
				video.Codec = v.(string)
			}
			if v, ok := dMap["width"]; ok {
				video.Width = v.(string)
			}
			if v, ok := dMap["height"]; ok {
				video.Height = v.(string)
			}
			if v, ok := dMap["fps"]; ok {
				video.Fps = v.(string)
			}
			if v, ok := dMap["animate_only_keep_key_frame"]; ok {
				video.AnimateOnlyKeepKeyFrame = v.(string)
			}
			if v, ok := dMap["animate_time_interval_of_frame"]; ok {
				video.AnimateTimeIntervalOfFrame = v.(string)
			}
			if v, ok := dMap["animate_frames_per_second"]; ok {
				video.AnimateFramesPerSecond = v.(string)
			}
			if v, ok := dMap["quality"]; ok {
				video.Quality = v.(string)
			}
			request.Video = &video
		}
	}

	if d.HasChange("time_interval") {
		if dMap, ok := helper.InterfacesHeadMap(d, "time_interval"); ok {
			timeInterval := cos.TimeInterval{}
			if v, ok := dMap["start"]; ok {
				timeInterval.Start = v.(string)
			}
			if v, ok := dMap["duration"]; ok {
				timeInterval.Duration = v.(string)
			}
			request.TimeInterval = &timeInterval
		}
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket).CI.UpdateMediaAnimationTemplate(ctx, &request, templateId)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaAnimationTemplate", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaAnimationTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaAnimationTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaAnimationTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_animation_template.delete")()
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
