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

func ResourceTencentCloudCiMediaTranscodeProTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaTranscodeProTemplateCreate,
		Read:   resourceTencentCloudCiMediaTranscodeProTemplateRead,
		Update: resourceTencentCloudCiMediaTranscodeProTemplateUpdate,
		Delete: resourceTencentCloudCiMediaTranscodeProTemplateDelete,
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
						"clip_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Fragment 配置，有效 当 格式 是 hls 和 dash。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"duration": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Fragmentation 时长，默认值 5s。",
									},
								},
							},
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
							Optional:    true,
							Description: "Codec 格式，默认值：`H.264`，当 格式 是 WebM，它 是 VP8，取值范围：`H.264`，`H.265`，`VP8`，`VP9`，`AV1`。",
						},
						"profile": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "编码 级别，Support baseline，main，high，auto- 当 Pixfmt 是 auto，此 参数 可以 仅 是 集合 到 auto，当 它 是 集合 到 other options， 参数 值 将 是 集合 到 auto- baseline: suitable 对于 mobile devices- main: suitable 对于 standard resolution devices- high: suitable 对于 high-resolution devices- Only H.264 支持 此 参数。",
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
						"interlaced": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "字段 pattern。",
						},
						"fps": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Frame 速率，取值范围：(0，60]，单位：fps。",
						},
						"bitrate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Bit 速率 的 视频 输出文件，取值范围：[10，50000]，单位: Kbps，auto 表示 adaptive bit 速率。",
						},
						"rotate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rotation angle，取值范围：[0，360)，单位：degree。",
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

			"audio": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Audio 信息，do 不 transmit Audio，其中 是 equivalent 到 deleting 音频 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"codec": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Codec 格式，值 aac，mp3，flac，amr，Vorbis，opus，pcm_s16le。",
						},
						"remove": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "是否delete 来源 音频 流， 值 是 true，false。",
						},
					},
				},
			},

			"trans_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "transcoding 配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"adj_dar_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resolution adjustment 方法，值 scale，crop，pad，none，当 aspect ratio 的 output 视频 是 different 从 original 视频，adjust resolution accordingly according 到 此 参数。",
						},
						"is_check_reso": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "是否check resolution，当 它 是 false，transcode according 到 配置 参数。",
						},
						"reso_adj_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resolution adjustment 模式，值 0，1; 0 表示 使用 original 视频 resolution; 1 表示 返回 transcoding failed，Take effect 当 IsCheckReso 是 true。",
						},
						"is_check_video_bitrate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "是否check 视频 代码 速率，当 它 是 false，transcode according 到 配置 参数。",
						},
						"video_bitrate_adj_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Video bit 速率 adjustment 方法，值 0，1; 当 output 视频 bit 速率 是 greater 比 original 视频 bit 速率，0 表示 使用 original 视频 bit 速率; 1 表示 返回 transcoding failed，Take effect 当 IsCheckVideoBitrate 是 true。",
						},
						"is_check_audio_bitrate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "是否check 音频 代码 速率，true，false，当 false，transcode according 到 配置 参数。",
						},
						"audio_bitrate_adj_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Audio bit 速率 adjustment 模式，值 0，1; 当 output 音频 bit 速率 是 greater 比 original 音频 bit 速率，0 表示 使用 original 音频 bit 速率; 1 表示 返回 transcoding failed，Take effect 当 IsCheckAudioBitrate 是 true。",
						},
						"delete_metadata": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "是否delete MetaData 信息 在 文件，true，false，当 false，keep 来源 文件 信息。",
						},
						"is_hdr2_sdr": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "是否enable HDR 到 SDR true，false。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCiMediaTranscodeProTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_transcode_pro_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request = cos.CreateMediaTranscodeProTemplateOptions{
			Tag: "TranscodePro",
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
		if clipConfigMap, ok := helper.InterfaceToMap(dMap, "clip_config"); ok {
			clipConfig := cos.ClipConfig{}
			if v, ok := clipConfigMap["duration"]; ok {
				clipConfig.Duration = v.(string)
			}
			container.ClipConfig = &clipConfig
		}
		request.Container = &container
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "video"); ok {
		video := cos.TranscodeProVideo{}
		if v, ok := dMap["codec"]; ok {
			video.Codec = v.(string)
		}
		if v, ok := dMap["profile"]; ok {
			video.Profile = v.(string)
		}
		if v, ok := dMap["width"]; ok {
			video.Width = v.(string)
		}
		if v, ok := dMap["height"]; ok {
			video.Height = v.(string)
		}
		if v, ok := dMap["interlaced"]; ok {
			video.Interlaced = v.(string)
		}
		if v, ok := dMap["fps"]; ok {
			video.Fps = v.(string)
		}
		if v, ok := dMap["bitrate"]; ok {
			video.Bitrate = v.(string)
		}
		if v, ok := dMap["rotate"]; ok {
			video.Rotate = v.(string)
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

	if dMap, ok := helper.InterfacesHeadMap(d, "audio"); ok {
		audio := cos.TranscodeProAudio{}
		if v, ok := dMap["codec"]; ok {
			audio.Codec = v.(string)
		}
		if v, ok := dMap["remove"]; ok {
			audio.Remove = v.(string)
		}
		request.Audio = &audio
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "trans_config"); ok {
		transConfig := cos.TransConfig{}
		if v, ok := dMap["adj_dar_method"]; ok {
			transConfig.AdjDarMethod = v.(string)
		}
		if v, ok := dMap["is_check_reso"]; ok {
			transConfig.IsCheckReso = v.(string)
		}
		if v, ok := dMap["reso_adj_method"]; ok {
			transConfig.ResoAdjMethod = v.(string)
		}
		if v, ok := dMap["is_check_video_bitrate"]; ok {
			transConfig.IsCheckVideoBitrate = v.(string)
		}
		if v, ok := dMap["video_bitrate_adj_method"]; ok {
			transConfig.VideoBitrateAdjMethod = v.(string)
		}
		if v, ok := dMap["is_check_audio_bitrate"]; ok {
			transConfig.IsCheckAudioBitrate = v.(string)
		}
		if v, ok := dMap["audio_bitrate_adj_method"]; ok {
			transConfig.AudioBitrateAdjMethod = v.(string)
		}
		if v, ok := dMap["delete_metadata"]; ok {
			transConfig.DeleteMetadata = v.(string)
		}
		if v, ok := dMap["is_hdr2_sdr"]; ok {
			transConfig.IsHdr2Sdr = v.(string)
		}
		request.TransConfig = &transConfig
	}

	var response *cos.CreateMediaTemplateResult
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket).CI.CreateMediaTranscodeProTemplate(ctx, &request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "CreateMediaTranscodeProTemplate", request, result)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaTranscodeProTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = response.Template.TemplateId
	d.SetId(bucket + tccommon.FILED_SP + templateId)

	return resourceTencentCloudCiMediaTranscodeProTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaTranscodeProTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_transcode_pro_template.read")()
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

	if template.Name != "" {
		_ = d.Set("name", template.Name)
	}
	mediaTranscodeProTemplate := template.TransProTpl
	if mediaTranscodeProTemplate != nil {
		if mediaTranscodeProTemplate.Container != nil {
			containerMap := map[string]interface{}{}

			if mediaTranscodeProTemplate.Container.Format != "" {
				containerMap["format"] = mediaTranscodeProTemplate.Container.Format
			}

			if mediaTranscodeProTemplate.Container.ClipConfig != nil {
				clipConfigMap := map[string]interface{}{}

				if mediaTranscodeProTemplate.Container.ClipConfig.Duration != "" {
					clipConfigMap["duration"] = mediaTranscodeProTemplate.Container.ClipConfig.Duration
				}

				containerMap["clip_config"] = []interface{}{clipConfigMap}
			}

			_ = d.Set("container", []interface{}{containerMap})
		}

		if mediaTranscodeProTemplate.Video != nil {
			videoMap := map[string]interface{}{}

			if mediaTranscodeProTemplate.Video.Codec != "" {
				videoMap["codec"] = mediaTranscodeProTemplate.Video.Codec
			}

			if mediaTranscodeProTemplate.Video.Profile != "" {
				videoMap["profile"] = mediaTranscodeProTemplate.Video.Profile
			}

			if mediaTranscodeProTemplate.Video.Width != "" {
				videoMap["width"] = mediaTranscodeProTemplate.Video.Width
			}

			if mediaTranscodeProTemplate.Video.Height != "" {
				videoMap["height"] = mediaTranscodeProTemplate.Video.Height
			}

			// if mediaTranscodeProTemplate.Video.Interlaced != "" {
			// 	videoMap["interlaced"] = mediaTranscodeProTemplate.Video.Interlaced
			// }

			if mediaTranscodeProTemplate.Video.Fps != "" {
				videoMap["fps"] = mediaTranscodeProTemplate.Video.Fps
			}

			if mediaTranscodeProTemplate.Video.Bitrate != "" {
				videoMap["bitrate"] = mediaTranscodeProTemplate.Video.Bitrate
			}

			if mediaTranscodeProTemplate.Video.Rotate != "" {
				videoMap["rotate"] = mediaTranscodeProTemplate.Video.Rotate
			}

			_ = d.Set("video", []interface{}{videoMap})
		}

		if mediaTranscodeProTemplate.TimeInterval != nil {
			timeIntervalMap := map[string]interface{}{}

			if mediaTranscodeProTemplate.TimeInterval.Start != "" {
				timeIntervalMap["start"] = mediaTranscodeProTemplate.TimeInterval.Start
			}

			if mediaTranscodeProTemplate.TimeInterval.Duration != "" {
				timeIntervalMap["duration"] = mediaTranscodeProTemplate.TimeInterval.Duration
			}

			_ = d.Set("time_interval", []interface{}{timeIntervalMap})
		}

		if mediaTranscodeProTemplate.Audio != nil {
			audioMap := map[string]interface{}{}

			if mediaTranscodeProTemplate.Audio.Codec != "" {
				audioMap["codec"] = mediaTranscodeProTemplate.Audio.Codec
			}

			if mediaTranscodeProTemplate.Audio.Remove != "" {
				audioMap["remove"] = mediaTranscodeProTemplate.Audio.Remove
			}

			_ = d.Set("audio", []interface{}{audioMap})
		}

		if mediaTranscodeProTemplate.TransConfig != nil {
			transConfigMap := map[string]interface{}{}

			if mediaTranscodeProTemplate.TransConfig.AdjDarMethod != "" {
				transConfigMap["adj_dar_method"] = mediaTranscodeProTemplate.TransConfig.AdjDarMethod
			}

			if mediaTranscodeProTemplate.TransConfig.IsCheckReso != "" {
				transConfigMap["is_check_reso"] = mediaTranscodeProTemplate.TransConfig.IsCheckReso
			}

			if mediaTranscodeProTemplate.TransConfig.ResoAdjMethod != "" {
				transConfigMap["reso_adj_method"] = mediaTranscodeProTemplate.TransConfig.ResoAdjMethod
			}

			if mediaTranscodeProTemplate.TransConfig.IsCheckVideoBitrate != "" {
				transConfigMap["is_check_video_bitrate"] = mediaTranscodeProTemplate.TransConfig.IsCheckVideoBitrate
			}

			if mediaTranscodeProTemplate.TransConfig.VideoBitrateAdjMethod != "" {
				transConfigMap["video_bitrate_adj_method"] = mediaTranscodeProTemplate.TransConfig.VideoBitrateAdjMethod
			}

			if mediaTranscodeProTemplate.TransConfig.IsCheckAudioBitrate != "" {
				transConfigMap["is_check_audio_bitrate"] = mediaTranscodeProTemplate.TransConfig.IsCheckAudioBitrate
			}

			if mediaTranscodeProTemplate.TransConfig.AudioBitrateAdjMethod != "" {
				transConfigMap["audio_bitrate_adj_method"] = mediaTranscodeProTemplate.TransConfig.AudioBitrateAdjMethod
			}

			if mediaTranscodeProTemplate.TransConfig.DeleteMetadata != "" {
				transConfigMap["delete_metadata"] = mediaTranscodeProTemplate.TransConfig.DeleteMetadata
			}

			if mediaTranscodeProTemplate.TransConfig.IsHdr2Sdr != "" {
				transConfigMap["is_hdr2_sdr"] = mediaTranscodeProTemplate.TransConfig.IsHdr2Sdr
			}

			_ = d.Set("trans_config", []interface{}{transConfigMap})
		}
	}

	return nil
}

func resourceTencentCloudCiMediaTranscodeProTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_transcode_pro_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := cos.CreateMediaTranscodeProTemplateOptions{
		Tag: "TranscodePro",
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

	if dMap, ok := helper.InterfacesHeadMap(d, "container"); ok {
		container := cos.Container{}
		if v, ok := dMap["format"]; ok {
			container.Format = v.(string)
		}
		if clipConfigMap, ok := helper.InterfaceToMap(dMap, "clip_config"); ok {
			clipConfig := cos.ClipConfig{}
			if v, ok := clipConfigMap["duration"]; ok {
				clipConfig.Duration = v.(string)
			}
			container.ClipConfig = &clipConfig
		}
		request.Container = &container
	}

	if d.HasChange("video") {
		if dMap, ok := helper.InterfacesHeadMap(d, "video"); ok {
			video := cos.TranscodeProVideo{}
			if v, ok := dMap["codec"]; ok {
				video.Codec = v.(string)
			}
			if v, ok := dMap["profile"]; ok {
				video.Profile = v.(string)
			}
			if v, ok := dMap["width"]; ok {
				video.Width = v.(string)
			}
			if v, ok := dMap["height"]; ok {
				video.Height = v.(string)
			}
			if v, ok := dMap["interlaced"]; ok {
				video.Interlaced = v.(string)
			}
			if v, ok := dMap["fps"]; ok {
				video.Fps = v.(string)
			}
			if v, ok := dMap["bitrate"]; ok {
				video.Bitrate = v.(string)
			}
			if v, ok := dMap["rotate"]; ok {
				video.Rotate = v.(string)
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

	if d.HasChange("audio") {
		if dMap, ok := helper.InterfacesHeadMap(d, "audio"); ok {
			audio := cos.TranscodeProAudio{}
			if v, ok := dMap["codec"]; ok {
				audio.Codec = v.(string)
			}
			if v, ok := dMap["remove"]; ok {
				audio.Remove = v.(string)
			}
			request.Audio = &audio
		}
	}

	if d.HasChange("vidtrans_configeo") {
		if dMap, ok := helper.InterfacesHeadMap(d, "trans_config"); ok {
			transConfig := cos.TransConfig{}
			if v, ok := dMap["adj_dar_method"]; ok {
				transConfig.AdjDarMethod = v.(string)
			}
			if v, ok := dMap["is_check_reso"]; ok {
				transConfig.IsCheckReso = v.(string)
			}
			if v, ok := dMap["reso_adj_method"]; ok {
				transConfig.ResoAdjMethod = v.(string)
			}
			if v, ok := dMap["is_check_video_bitrate"]; ok {
				transConfig.IsCheckVideoBitrate = v.(string)
			}
			if v, ok := dMap["video_bitrate_adj_method"]; ok {
				transConfig.VideoBitrateAdjMethod = v.(string)
			}
			if v, ok := dMap["is_check_audio_bitrate"]; ok {
				transConfig.IsCheckAudioBitrate = v.(string)
			}
			if v, ok := dMap["audio_bitrate_adj_method"]; ok {
				transConfig.AudioBitrateAdjMethod = v.(string)
			}
			if v, ok := dMap["delete_metadata"]; ok {
				transConfig.DeleteMetadata = v.(string)
			}
			if v, ok := dMap["is_hdr2_sdr"]; ok {
				transConfig.IsHdr2Sdr = v.(string)
			}
			request.TransConfig = &transConfig
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket).CI.UpdateMediaTranscodeProTemplate(ctx, &request, templateId)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaTranscodeProTemplate", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaTranscodeProTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaTranscodeProTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaTranscodeProTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_transcode_pro_template.delete")()
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
