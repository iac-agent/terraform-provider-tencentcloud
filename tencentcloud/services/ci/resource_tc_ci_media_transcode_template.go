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

func ResourceTencentCloudCiMediaTranscodeTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaTranscodeTemplateCreate,
		Read:   resourceTencentCloudCiMediaTranscodeTemplateRead,
		Update: resourceTencentCloudCiMediaTranscodeTemplateUpdate,
		Delete: resourceTencentCloudCiMediaTranscodeTemplateDelete,
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
						"remove": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "是否delete 视频 流，true，false。",
						},
						"profile": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "编码 级别，Support baseline，main，high，auto- 当 Pixfmt 是 auto，此 参数 可以 仅 是 集合 到 auto，当 它 是 集合 到 other options， 参数 值 将 是 集合 到 auto- baseline: suitable 对于 mobile devices- main: suitable 对于 standard resolution devices- high: suitable 对于 high-resolution devices- Only H.264 支持 此 参数。",
						},
						"bitrate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Bit 速率 的 视频 输出文件，取值范围：[10，50000]，单位: Kbps，auto 表示 adaptive bit 速率。",
						},
						"crf": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Bit 速率-quality control factor，取值范围：(0，51]，如果 Crf 是 集合， setting 的 Bitrate 将 是 无效，当 Bitrate 是 空， 默认为 25。",
						},
						"gop": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "最大frames between 键 frames，取值范围：[1，100000]。",
						},
						"preset": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Video Algorithm Presets- H.264 支持 此 参数， 值 是 veryfast，fast，medium，slow，slower- VP8 支持 此 参数， 值 是 good，realtime- AV1 支持 此 参数， 值 是 5 (recommended 值)，4- H.265 和 VP9 do 不 support 此 参数。",
						},
						"bufsize": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "buffer 大小，取值范围：[1000，128000]，单位：Kb，此 参数 是 不 支持 当 Codec 是 VP8/VP9。",
						},
						"maxrate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Peak 视频 bit 速率，取值范围：[10，50000]，单位：Kbps，此 参数 是 不 支持 当 Codec 是 VP8/VP9。",
						},
						"pixfmt": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "视频 color 格式，H.264 support: yuv420p，yuv422p，yuv444p，yuvj420p，yuvj422p，yuvj444p，auto，H.265 support: yuv420p，yuv420p10le，auto，此 参数 是 不 支持 当 Codec 是 VP8/VP9/AV1。",
						},
						"long_short_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Adaptive 长度,true，false，此 参数 是 不 支持 当 Codec 是 VP8/VP9/AV1。",
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
							Optional:    true,
							Description: "Codec 格式，值 aac，mp3，flac，amr，Vorbis，opus，pcm_s16le。",
						},
						"samplerate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Sampling Rate- 单位：Hz- 可选 8000，11025，12000，16000，22050，24000，32000，44100，48000，88200，96000- Different packages，mp3 支持 different sampling rates，作为 shown 在 表 below- 当 Codec 是 集合 到 amr，仅 8000 是 支持- 当 Codec 是 集合 到 opus，它 支持 8000，16000，24000，48000。",
						},
						"bitrate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Original 音频 bit 速率，单位: Kbps，取值范围：[8，1000]。",
						},
						"channels": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "数量 channels- 当 Codec 是 集合 到 aac/flac，support 1，2，4，5，6，8- 当 Codec 是 集合 到 mp3/opus，support 1，2- 当 Codec 是 集合 到 Vorbis，仅 2 是 支持- 当 Codec 是 集合 到 amr，仅 1 是 支持- 当 Codec 是 集合 到 pcm_s16le，仅 1 和 2 是 支持- 当 encapsulation 格式 是 dash，8 是 不 支持。",
						},
						"remove": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "是否delete 来源 音频 流， 值 是 true，false。",
						},
						"keep_two_tracks": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Keep dual 音频 tracks， 值 是 true，false. 此 参数 是 无效 当 Video.Codec 是 H.265。",
						},
						"switch_track": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Convert track， 值 是 true，false. 此 参数 是 无效 当 Video.Codec 是 H.265。",
						},
						"sample_format": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Sampling bit 宽度- 当 Codec 是 集合 到 aac，support fltp- 当 Codec 是 集合 到 mp3，fltp，s16p，s32p 是 支持- 当 Codec 是 集合 到 flac，s16，s32，s16p，s32p 是 支持- 当 Codec 是 集合 到 amr，support s16，s16p- 当 Codec 是 集合 到 opus，support s16- 当 Codec 是 集合 到 pcm_s16le，support s16- 当 Codec 是 集合 到 Vorbis，support fltp- 此 参数 是 无效 当 Video.Codec 是 H.265。",
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
						"hls_encrypt": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "hls 加密 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_hls_encrypt": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "是否enable HLS 加密，support 加密 当 Container.格式 是 hls。",
									},
									"uri_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "HLS encrypted 键，此 参数 是 仅 meaningful 当 IsHlsEncrypt 是 true。",
									},
								},
							},
						},
					},
				},
			},

			"audio_mix": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "mixing 参数。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"audio_source": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "media 地址 的 音频 track 该 needs 到 是 mixed。",
						},
						"mix_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Mixing 模式 Repeat: background sound loop，Once: background sound 是 played once。",
						},
						"replace": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "是否replace original 音频 的 Input media 文件 使用 mixed 音频 track media。",
						},
						"effect_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Mix Fade Configuration。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_start_fadein": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "启用 fade 在。",
									},
									"start_fadein_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Fade 在 时长，greater 比 0，support floating point numbers。",
									},
									"enable_end_fadeout": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "启用 fade out。",
									},
									"end_fadeout_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "fade out 时间，greater 比 0，support floating point numbers。",
									},
									"enable_bgm_fade": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Enable bgm conversion fade 在。",
									},
									"bgm_fade_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "bgm transition fade-在 时长，support floating point numbers。",
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

func resourceTencentCloudCiMediaTranscodeTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_transcode_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request = cos.CreateMediaTranscodeTemplateOptions{
			Tag: "Transcode",
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
		video := cos.Video{}
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
		if v, ok := dMap["remove"]; ok {
			video.Remove = v.(string)
		}
		if v, ok := dMap["profile"]; ok {
			video.Profile = v.(string)
		}
		if v, ok := dMap["bitrate"]; ok {
			video.Bitrate = v.(string)
		}
		if v, ok := dMap["crf"]; ok {
			video.Crf = v.(string)
		}
		if v, ok := dMap["gop"]; ok {
			video.Gop = v.(string)
		}
		if v, ok := dMap["preset"]; ok {
			video.Preset = v.(string)
		}
		if v, ok := dMap["bufsize"]; ok {
			video.Bufsize = v.(string)
		}
		if v, ok := dMap["maxrate"]; ok {
			video.Maxrate = v.(string)
		}
		if v, ok := dMap["pixfmt"]; ok {
			video.Pixfmt = v.(string)
		}
		if v, ok := dMap["long_short_mode"]; ok {
			video.LongShortMode = v.(string)
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
		audio := cos.Audio{}
		if v, ok := dMap["codec"]; ok {
			audio.Codec = v.(string)
		}
		if v, ok := dMap["samplerate"]; ok {
			audio.Samplerate = v.(string)
		}
		if v, ok := dMap["bitrate"]; ok {
			audio.Bitrate = v.(string)
		}
		if v, ok := dMap["channels"]; ok {
			audio.Channels = v.(string)
		}
		if v, ok := dMap["remove"]; ok {
			audio.Remove = v.(string)
		}
		if v, ok := dMap["keep_two_tracks"]; ok {
			audio.KeepTwoTracks = v.(string)
		}
		if v, ok := dMap["switch_track"]; ok {
			audio.SwitchTrack = v.(string)
		}
		if v, ok := dMap["sample_format"]; ok {
			audio.SampleFormat = v.(string)
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
		if hlsEncryptMap, ok := helper.InterfaceToMap(dMap, "hls_encrypt"); ok {
			hlsEncrypt := cos.HlsEncrypt{}
			if v, ok := hlsEncryptMap["is_hls_encrypt"]; ok {
				if v.(string) == "true" {
					hlsEncrypt.IsHlsEncrypt = true
				} else {
					hlsEncrypt.IsHlsEncrypt = false
				}
			}
			if v, ok := hlsEncryptMap["uri_key"]; ok {
				hlsEncrypt.UriKey = v.(string)
			}
			transConfig.HlsEncrypt = &hlsEncrypt
		}
		request.TransConfig = &transConfig
	}

	if v, ok := d.GetOk("audio_mix"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			audioMix := cos.AudioMix{}
			if v, ok := dMap["audio_source"]; ok {
				audioMix.AudioSource = v.(string)
			}
			if v, ok := dMap["mix_mode"]; ok {
				audioMix.MixMode = v.(string)
			}
			if v, ok := dMap["replace"]; ok {
				audioMix.Replace = v.(string)
			}
			if effectConfigMap, ok := helper.InterfaceToMap(dMap, "effect_config"); ok {
				effectConfig := cos.EffectConfig{}
				if v, ok := effectConfigMap["enable_start_fadein"]; ok {
					effectConfig.EnableStartFadein = v.(string)
				}
				if v, ok := effectConfigMap["start_fadein_time"]; ok {
					effectConfig.StartFadeinTime = v.(string)
				}
				if v, ok := effectConfigMap["enable_end_fadeout"]; ok {
					effectConfig.EnableEndFadeout = v.(string)
				}
				if v, ok := effectConfigMap["end_fadeout_time"]; ok {
					effectConfig.EndFadeoutTime = v.(string)
				}
				if v, ok := effectConfigMap["enable_bgm_fade"]; ok {
					effectConfig.EnableBgmFade = v.(string)
				}
				if v, ok := effectConfigMap["bgm_fade_time"]; ok {
					effectConfig.BgmFadeTime = v.(string)
				}
				audioMix.EffectConfig = &effectConfig
			}
			request.AudioMixArray = append(request.AudioMixArray, audioMix)
		}
	}

	var response *cos.CreateMediaTemplateResult
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket).CI.CreateMediaTranscodeTemplate(ctx, &request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "CreateMediaTranscodeTemplate", request, result)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaTranscodeTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = response.Template.TemplateId
	d.SetId(bucket + tccommon.FILED_SP + templateId)

	return resourceTencentCloudCiMediaTranscodeTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaTranscodeTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_transcode_template.read")()
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
	mediaTranscodeTemplate := template.TransTpl
	if mediaTranscodeTemplate.Container != nil {
		containerMap := map[string]interface{}{}

		if mediaTranscodeTemplate.Container.Format != "" {
			containerMap["format"] = mediaTranscodeTemplate.Container.Format
		}

		if mediaTranscodeTemplate.Container.ClipConfig != nil {
			clipConfigMap := map[string]interface{}{}

			if mediaTranscodeTemplate.Container.ClipConfig.Duration != "" {
				clipConfigMap["duration"] = mediaTranscodeTemplate.Container.ClipConfig.Duration
			}

			containerMap["clip_config"] = []interface{}{clipConfigMap}
		}

		_ = d.Set("container", []interface{}{containerMap})
	}

	if mediaTranscodeTemplate.Video != nil {
		videoMap := map[string]interface{}{}

		if mediaTranscodeTemplate.Video.Codec != "" {
			videoMap["codec"] = mediaTranscodeTemplate.Video.Codec
		}

		if mediaTranscodeTemplate.Video.Width != "" {
			videoMap["width"] = mediaTranscodeTemplate.Video.Width
		}

		if mediaTranscodeTemplate.Video.Height != "" {
			videoMap["height"] = mediaTranscodeTemplate.Video.Height
		}

		if mediaTranscodeTemplate.Video.Fps != "" {
			videoMap["fps"] = mediaTranscodeTemplate.Video.Fps
		}

		if mediaTranscodeTemplate.Video.Remove != "" {
			videoMap["remove"] = mediaTranscodeTemplate.Video.Remove
		}

		if mediaTranscodeTemplate.Video.Profile != "" {
			videoMap["profile"] = mediaTranscodeTemplate.Video.Profile
		}

		if mediaTranscodeTemplate.Video.Bitrate != "" {
			videoMap["bitrate"] = mediaTranscodeTemplate.Video.Bitrate
		}

		if mediaTranscodeTemplate.Video.Crf != "" {
			videoMap["crf"] = mediaTranscodeTemplate.Video.Crf
		}

		if mediaTranscodeTemplate.Video.Gop != "" {
			videoMap["gop"] = mediaTranscodeTemplate.Video.Gop
		}

		if mediaTranscodeTemplate.Video.Preset != "" {
			videoMap["preset"] = mediaTranscodeTemplate.Video.Preset
		}

		if mediaTranscodeTemplate.Video.Bufsize != "" {
			videoMap["bufsize"] = mediaTranscodeTemplate.Video.Bufsize
		}

		if mediaTranscodeTemplate.Video.Maxrate != "" {
			videoMap["maxrate"] = mediaTranscodeTemplate.Video.Maxrate
		}

		if mediaTranscodeTemplate.Video.Pixfmt != "" {
			videoMap["pixfmt"] = mediaTranscodeTemplate.Video.Pixfmt
		}

		if mediaTranscodeTemplate.Video.LongShortMode != "" {
			videoMap["long_short_mode"] = mediaTranscodeTemplate.Video.LongShortMode
		}

		if mediaTranscodeTemplate.Video.Rotate != "" {
			videoMap["rotate"] = mediaTranscodeTemplate.Video.Rotate
		}

		_ = d.Set("video", []interface{}{videoMap})
	}

	if mediaTranscodeTemplate.TimeInterval != nil {
		timeIntervalMap := map[string]interface{}{}

		if mediaTranscodeTemplate.TimeInterval.Start != "" {
			timeIntervalMap["start"] = mediaTranscodeTemplate.TimeInterval.Start
		}

		if mediaTranscodeTemplate.TimeInterval.Duration != "" {
			timeIntervalMap["duration"] = mediaTranscodeTemplate.TimeInterval.Duration
		}

		_ = d.Set("time_interval", []interface{}{timeIntervalMap})
	}

	if mediaTranscodeTemplate.Audio != nil {
		audioMap := map[string]interface{}{}

		if mediaTranscodeTemplate.Audio.Codec != "" {
			audioMap["codec"] = mediaTranscodeTemplate.Audio.Codec
		}

		if mediaTranscodeTemplate.Audio.Samplerate != "" {
			audioMap["samplerate"] = mediaTranscodeTemplate.Audio.Samplerate
		}

		if mediaTranscodeTemplate.Audio.Bitrate != "" {
			audioMap["bitrate"] = mediaTranscodeTemplate.Audio.Bitrate
		}

		if mediaTranscodeTemplate.Audio.Channels != "" {
			audioMap["channels"] = mediaTranscodeTemplate.Audio.Channels
		}

		if mediaTranscodeTemplate.Audio.Remove != "" {
			audioMap["remove"] = mediaTranscodeTemplate.Audio.Remove
		}

		if mediaTranscodeTemplate.Audio.KeepTwoTracks != "" {
			audioMap["keep_two_tracks"] = mediaTranscodeTemplate.Audio.KeepTwoTracks
		}

		if mediaTranscodeTemplate.Audio.SwitchTrack != "" {
			audioMap["switch_track"] = mediaTranscodeTemplate.Audio.SwitchTrack
		}

		if mediaTranscodeTemplate.Audio.SampleFormat != "" {
			audioMap["sample_format"] = mediaTranscodeTemplate.Audio.SampleFormat
		}

		_ = d.Set("audio", []interface{}{audioMap})
	}

	if mediaTranscodeTemplate.TransConfig != nil {
		transConfigMap := map[string]interface{}{}

		if mediaTranscodeTemplate.TransConfig.AdjDarMethod != "" {
			transConfigMap["adj_dar_method"] = mediaTranscodeTemplate.TransConfig.AdjDarMethod
		}

		if mediaTranscodeTemplate.TransConfig.IsCheckReso != "" {
			transConfigMap["is_check_reso"] = mediaTranscodeTemplate.TransConfig.IsCheckReso
		}

		if mediaTranscodeTemplate.TransConfig.ResoAdjMethod != "" {
			transConfigMap["reso_adj_method"] = mediaTranscodeTemplate.TransConfig.ResoAdjMethod
		}

		if mediaTranscodeTemplate.TransConfig.IsCheckVideoBitrate != "" {
			transConfigMap["is_check_video_bitrate"] = mediaTranscodeTemplate.TransConfig.IsCheckVideoBitrate
		}

		if mediaTranscodeTemplate.TransConfig.VideoBitrateAdjMethod != "" {
			transConfigMap["video_bitrate_adj_method"] = mediaTranscodeTemplate.TransConfig.VideoBitrateAdjMethod
		}

		if mediaTranscodeTemplate.TransConfig.IsCheckAudioBitrate != "" {
			transConfigMap["is_check_audio_bitrate"] = mediaTranscodeTemplate.TransConfig.IsCheckAudioBitrate
		}

		if mediaTranscodeTemplate.TransConfig.AudioBitrateAdjMethod != "" {
			transConfigMap["audio_bitrate_adj_method"] = mediaTranscodeTemplate.TransConfig.AudioBitrateAdjMethod
		}

		if mediaTranscodeTemplate.TransConfig.DeleteMetadata != "" {
			transConfigMap["delete_metadata"] = mediaTranscodeTemplate.TransConfig.DeleteMetadata
		}

		if mediaTranscodeTemplate.TransConfig.IsHdr2Sdr != "" {
			transConfigMap["is_hdr2_sdr"] = mediaTranscodeTemplate.TransConfig.IsHdr2Sdr
		}

		if mediaTranscodeTemplate.TransConfig.HlsEncrypt != nil {
			hlsEncryptMap := map[string]interface{}{}

			hlsEncryptMap["is_hls_encrypt"] = fmt.Sprintf("%v", mediaTranscodeTemplate.TransConfig.HlsEncrypt.IsHlsEncrypt)

			if mediaTranscodeTemplate.TransConfig.HlsEncrypt.UriKey != "" {
				hlsEncryptMap["uri_key"] = mediaTranscodeTemplate.TransConfig.HlsEncrypt.UriKey
			}

			transConfigMap["hls_encrypt"] = []interface{}{hlsEncryptMap}
		}

		_ = d.Set("trans_config", []interface{}{transConfigMap})
	}

	if mediaTranscodeTemplate.AudioMixArray != nil {
		audioMixArrayList := []interface{}{}
		for _, audioMix := range mediaTranscodeTemplate.AudioMixArray {
			audioMixArrayMap := map[string]interface{}{}

			if audioMix.AudioSource != "" {
				audioMixArrayMap["audio_source"] = audioMix.AudioSource
			}

			if audioMix.MixMode != "" {
				audioMixArrayMap["mix_mode"] = audioMix.MixMode
			}

			if audioMix.Replace != "" {
				audioMixArrayMap["replace"] = audioMix.Replace
			}

			if audioMix.EffectConfig != nil {
				effectConfigMap := map[string]interface{}{}

				if audioMix.EffectConfig.EnableStartFadein != "" {
					effectConfigMap["enable_start_fadein"] = audioMix.EffectConfig.EnableStartFadein
				}

				if audioMix.EffectConfig.StartFadeinTime != "" {
					effectConfigMap["start_fadein_time"] = audioMix.EffectConfig.StartFadeinTime
				}

				if audioMix.EffectConfig.EnableEndFadeout != "" {
					effectConfigMap["enable_end_fadeout"] = audioMix.EffectConfig.EnableEndFadeout
				}

				if audioMix.EffectConfig.EndFadeoutTime != "" {
					effectConfigMap["end_fadeout_time"] = audioMix.EffectConfig.EndFadeoutTime
				}

				if audioMix.EffectConfig.EnableBgmFade != "" {
					effectConfigMap["enable_bgm_fade"] = audioMix.EffectConfig.EnableBgmFade
				}

				if audioMix.EffectConfig.BgmFadeTime != "" {
					effectConfigMap["bgm_fade_time"] = audioMix.EffectConfig.BgmFadeTime
				}

				audioMixArrayMap["effect_config"] = []interface{}{effectConfigMap}
			}

			audioMixArrayList = append(audioMixArrayList, audioMixArrayMap)
		}

		_ = d.Set("audio_mix_array", audioMixArrayList)

	}

	return nil
}

func resourceTencentCloudCiMediaTranscodeTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_transcode_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := cos.CreateMediaTranscodeTemplateOptions{
		Tag: "Transcode",
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

	if d.HasChange("video") {
		if dMap, ok := helper.InterfacesHeadMap(d, "video"); ok {
			video := cos.Video{}
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
			if v, ok := dMap["remove"]; ok {
				video.Remove = v.(string)
			}
			if v, ok := dMap["profile"]; ok {
				video.Profile = v.(string)
			}
			if v, ok := dMap["bitrate"]; ok {
				video.Bitrate = v.(string)
			}
			if v, ok := dMap["crf"]; ok {
				video.Crf = v.(string)
			}
			if v, ok := dMap["gop"]; ok {
				video.Gop = v.(string)
			}
			if v, ok := dMap["preset"]; ok {
				video.Preset = v.(string)
			}
			if v, ok := dMap["bufsize"]; ok {
				video.Bufsize = v.(string)
			}
			if v, ok := dMap["maxrate"]; ok {
				video.Maxrate = v.(string)
			}
			if v, ok := dMap["pixfmt"]; ok {
				video.Pixfmt = v.(string)
			}
			if v, ok := dMap["long_short_mode"]; ok {
				video.LongShortMode = v.(string)
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
			audio := cos.Audio{}
			if v, ok := dMap["codec"]; ok {
				audio.Codec = v.(string)
			}
			if v, ok := dMap["samplerate"]; ok {
				audio.Samplerate = v.(string)
			}
			if v, ok := dMap["bitrate"]; ok {
				audio.Bitrate = v.(string)
			}
			if v, ok := dMap["channels"]; ok {
				audio.Channels = v.(string)
			}
			if v, ok := dMap["remove"]; ok {
				audio.Remove = v.(string)
			}
			if v, ok := dMap["keep_two_tracks"]; ok {
				audio.KeepTwoTracks = v.(string)
			}
			if v, ok := dMap["switch_track"]; ok {
				audio.SwitchTrack = v.(string)
			}
			if v, ok := dMap["sample_format"]; ok {
				audio.SampleFormat = v.(string)
			}
			request.Audio = &audio
		}
	}

	if d.HasChange("trans_config") {
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
			if hlsEncryptMap, ok := helper.InterfaceToMap(dMap, "hls_encrypt"); ok {
				hlsEncrypt := cos.HlsEncrypt{}
				if v, ok := hlsEncryptMap["is_hls_encrypt"]; ok {
					if v.(string) == "true" {
						hlsEncrypt.IsHlsEncrypt = true
					} else {
						hlsEncrypt.IsHlsEncrypt = false
					}
				}
				if v, ok := hlsEncryptMap["uri_key"]; ok {
					hlsEncrypt.UriKey = v.(string)
				}
				transConfig.HlsEncrypt = &hlsEncrypt
			}
			request.TransConfig = &transConfig
		}
	}

	if d.HasChange("audio_mix") {
		if v, ok := d.GetOk("audio_mix"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				audioMix := cos.AudioMix{}
				if v, ok := dMap["audio_source"]; ok {
					audioMix.AudioSource = v.(string)
				}
				if v, ok := dMap["mix_mode"]; ok {
					audioMix.MixMode = v.(string)
				}
				if v, ok := dMap["replace"]; ok {
					audioMix.Replace = v.(string)
				}
				if effectConfigMap, ok := helper.InterfaceToMap(dMap, "effect_config"); ok {
					effectConfig := cos.EffectConfig{}
					if v, ok := effectConfigMap["enable_start_fadein"]; ok {
						effectConfig.EnableStartFadein = v.(string)
					}
					if v, ok := effectConfigMap["start_fadein_time"]; ok {
						effectConfig.StartFadeinTime = v.(string)
					}
					if v, ok := effectConfigMap["enable_end_fadeout"]; ok {
						effectConfig.EnableEndFadeout = v.(string)
					}
					if v, ok := effectConfigMap["end_fadeout_time"]; ok {
						effectConfig.EndFadeoutTime = v.(string)
					}
					if v, ok := effectConfigMap["enable_bgm_fade"]; ok {
						effectConfig.EnableBgmFade = v.(string)
					}
					if v, ok := effectConfigMap["bgm_fade_time"]; ok {
						effectConfig.BgmFadeTime = v.(string)
					}
					audioMix.EffectConfig = &effectConfig
				}
				request.AudioMixArray = append(request.AudioMixArray, audioMix)
			}
		}
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket).CI.UpdateMediaTranscodeTemplate(ctx, &request, templateId)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaTranscodeTemplate", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaTranscodeTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaTranscodeTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaTranscodeTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_transcode_template.delete")()
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
