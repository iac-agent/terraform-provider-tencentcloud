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

func ResourceTencentCloudCiMediaConcatTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaConcatTemplateCreate,
		Read:   resourceTencentCloudCiMediaConcatTemplateRead,
		Update: resourceTencentCloudCiMediaConcatTemplateUpdate,
		Delete: resourceTencentCloudCiMediaConcatTemplateDelete,
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

			"concat_template": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "stitching template。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"concat_fragment": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Package 格式",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Splicing 对象 地址",
									},
									"mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "节点 类型，`start`，`end`。",
									},
								},
							},
						},
						"audio": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "音频 参数， 目标 文件 does 不 require Audio 信息，need 到 集合 Audio.Remove 到 true。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"codec": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Codec 格式，值 aac，mp3。",
									},
									"samplerate": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Sampling Rate- 单位：Hz- 可选 11025，22050，32000，44100，48000，96000- Different packages，mp3 支持 different sampling rates，作为 shown 在 表 below。",
									},
									"bitrate": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Original 音频 bit 速率，单位: Kbps，取值范围：[8，1000]。",
									},
									"channels": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "数量 channels- 当 Codec 是 集合 到 aac，support 1，2，4，5，6，8- 当 Codec 是 集合 到 mp3，support 1，2。",
									},
								},
							},
						},
						"video": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "视频 信息，do 不 upload Video，其中 是 equivalent 到 deleting 视频 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"codec": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Codec 格式 `H.264`。",
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
									"bitrate": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Original 音频 bit 速率，单位: Kbps，取值范围：[8，1000]。",
									},
									"fps": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Frame 速率，取值范围：(0，60]，单位：fps。",
									},
									"crf": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Bit 速率-quality control factor，取值范围：(0，51]，如果 Crf 是 集合， setting 的 Bitrate 将 是 无效，当 Bitrate 是 空， 默认为 25。",
									},
									"remove": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "是否delete 来源 音频 流， 值 是 true，false。",
									},
									"rotate": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Rotation angle，取值范围：[0，360)，单位：degree。",
									},
								},
							},
						},
						"container": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Only splicing without transcoding。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"format": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Container 格式: mp4，flv，hls，ts，mp3，aac。",
									},
								},
							},
						},
						"audio_mix": {
							Type:        schema.TypeList,
							Optional:    true,
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
													Computed:    true,
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
				},
			},
		},
	}
}

func resourceTencentCloudCiMediaConcatTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_concat_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request = cos.CreateMediaConcatTemplateOptions{
			Tag: "Concat",
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

	if dMap, ok := helper.InterfacesHeadMap(d, "concat_template"); ok {
		concatTemplate := cos.ConcatTemplate{}
		if v, ok := dMap["concat_fragment"]; ok {
			for _, item := range v.([]interface{}) {
				concatFragmentMap := item.(map[string]interface{})
				concatFragment := cos.ConcatFragment{}
				if v, ok := concatFragmentMap["url"]; ok {
					concatFragment.Url = v.(string)
				}
				if v, ok := concatFragmentMap["mode"]; ok {
					concatFragment.Mode = v.(string)
				}
				concatTemplate.ConcatFragment = append(concatTemplate.ConcatFragment, concatFragment)
			}
		}
		if audioMap, ok := helper.InterfaceToMap(dMap, "audio"); ok {
			audio := cos.Audio{}
			if v, ok := audioMap["codec"]; ok {
				audio.Codec = v.(string)
			}
			if v, ok := audioMap["samplerate"]; ok {
				audio.Samplerate = v.(string)
			}
			if v, ok := audioMap["bitrate"]; ok {
				audio.Bitrate = v.(string)
			}
			if v, ok := audioMap["channels"]; ok {
				audio.Channels = v.(string)
			}
			concatTemplate.Audio = &audio
		}
		if videoMap, ok := helper.InterfaceToMap(dMap, "video"); ok {
			video := cos.Video{}
			if v, ok := videoMap["codec"]; ok {
				video.Codec = v.(string)
			}
			if v, ok := videoMap["width"]; ok {
				video.Width = v.(string)
			}
			if v, ok := videoMap["height"]; ok {
				video.Height = v.(string)
			}
			if v, ok := videoMap["bitrate"]; ok {
				video.Bitrate = v.(string)
			}
			if v, ok := videoMap["fps"]; ok {
				video.Fps = v.(string)
			}
			if v, ok := videoMap["crf"]; ok {
				video.Crf = v.(string)
			}
			if v, ok := videoMap["remove"]; ok {
				video.Remove = v.(string)
			}
			if v, ok := videoMap["rotate"]; ok {
				video.Rotate = v.(string)
			}
			concatTemplate.Video = &video
		}
		if containerMap, ok := helper.InterfaceToMap(dMap, "container"); ok {
			container := cos.Container{}
			if v, ok := containerMap["format"]; ok {
				container.Format = v.(string)
			}
			concatTemplate.Container = &container
		}
		if v, ok := dMap["audio_mix"]; ok {
			for _, item := range v.([]interface{}) {
				audioMixMap := item.(map[string]interface{})
				audioMix := cos.AudioMix{}
				if v, ok := audioMixMap["audio_source"]; ok {
					audioMix.AudioSource = v.(string)
				}
				if v, ok := audioMixMap["mix_mode"]; ok {
					audioMix.MixMode = v.(string)
				}
				if v, ok := audioMixMap["replace"]; ok {
					audioMix.Replace = v.(string)
				}
				if effectConfigMap, ok := helper.InterfaceToMap(audioMixMap, "effect_config"); ok {
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
				concatTemplate.AudioMixArray = append(concatTemplate.AudioMixArray, audioMix)
			}
		}
		request.ConcatTemplate = &concatTemplate
	}

	var response *cos.CreateMediaTemplateResult
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket).CI.CreateMediaConcatTemplate(ctx, &request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "CreateMediaConcatTemplate", request, result)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaConcatTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = response.Template.TemplateId
	d.SetId(bucket + tccommon.FILED_SP + templateId)

	return resourceTencentCloudCiMediaConcatTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaConcatTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_concat_template.read")()
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

	mediaConcatTemplate, err := service.DescribeCiMediaTemplateById(ctx, bucket, templateId)
	if err != nil {
		return err
	}

	if mediaConcatTemplate == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	_ = d.Set("bucket", bucket)

	if mediaConcatTemplate.Name != "" {
		_ = d.Set("name", mediaConcatTemplate.Name)
	}

	if mediaConcatTemplate.ConcatTemplate != nil {
		concatTemplateMap := map[string]interface{}{}

		if mediaConcatTemplate.ConcatTemplate.ConcatFragment != nil {
			concatFragmentList := []interface{}{}
			for _, concatFragment := range mediaConcatTemplate.ConcatTemplate.ConcatFragment {
				concatFragmentMap := map[string]interface{}{}
				if concatFragment.Url != "" {
					concatFragmentMap["url"] = concatFragment.Url
				}
				if concatFragment.Mode != "" {
					concatFragmentMap["mode"] = concatFragment.Mode
				}
				concatFragmentList = append(concatFragmentList, concatFragmentMap)
			}
			concatTemplateMap["concat_fragment"] = concatFragmentList
		}

		if mediaConcatTemplate.ConcatTemplate.Audio != nil {
			audioMap := map[string]interface{}{}

			if mediaConcatTemplate.ConcatTemplate.Audio.Codec != "" {
				audioMap["codec"] = mediaConcatTemplate.ConcatTemplate.Audio.Codec
			}

			if mediaConcatTemplate.ConcatTemplate.Audio.Samplerate != "" {
				audioMap["samplerate"] = mediaConcatTemplate.ConcatTemplate.Audio.Samplerate
			}

			if mediaConcatTemplate.ConcatTemplate.Audio.Bitrate != "" {
				audioMap["bitrate"] = mediaConcatTemplate.ConcatTemplate.Audio.Bitrate
			}

			if mediaConcatTemplate.ConcatTemplate.Audio.Channels != "" {
				audioMap["channels"] = mediaConcatTemplate.ConcatTemplate.Audio.Channels
			}

			concatTemplateMap["audio"] = []interface{}{audioMap}
		}

		if mediaConcatTemplate.ConcatTemplate.Video != nil {
			videoMap := map[string]interface{}{}

			if mediaConcatTemplate.ConcatTemplate.Video.Codec != "" {
				videoMap["codec"] = mediaConcatTemplate.ConcatTemplate.Video.Codec
			}

			if mediaConcatTemplate.ConcatTemplate.Video.Width != "" {
				videoMap["width"] = mediaConcatTemplate.ConcatTemplate.Video.Width
			}

			if mediaConcatTemplate.ConcatTemplate.Video.Height != "" {
				videoMap["height"] = mediaConcatTemplate.ConcatTemplate.Video.Height
			}

			if mediaConcatTemplate.ConcatTemplate.Video.Bitrate != "" {
				videoMap["bitrate"] = mediaConcatTemplate.ConcatTemplate.Video.Bitrate
			}

			if mediaConcatTemplate.ConcatTemplate.Video.Fps != "" {
				videoMap["fps"] = mediaConcatTemplate.ConcatTemplate.Video.Fps
			}

			if mediaConcatTemplate.ConcatTemplate.Video.Crf != "" {
				videoMap["crf"] = mediaConcatTemplate.ConcatTemplate.Video.Crf
			}

			if mediaConcatTemplate.ConcatTemplate.Video.Remove != "" {
				videoMap["remove"] = mediaConcatTemplate.ConcatTemplate.Video.Remove
			}

			if mediaConcatTemplate.ConcatTemplate.Video.Rotate != "" {
				videoMap["rotate"] = mediaConcatTemplate.ConcatTemplate.Video.Rotate
			}

			concatTemplateMap["video"] = []interface{}{videoMap}
		}

		if mediaConcatTemplate.ConcatTemplate.Container != nil {
			containerMap := map[string]interface{}{}

			if mediaConcatTemplate.ConcatTemplate.Container.Format != "" {
				containerMap["format"] = mediaConcatTemplate.ConcatTemplate.Container.Format
			}

			concatTemplateMap["container"] = []interface{}{containerMap}
		}

		if mediaConcatTemplate.ConcatTemplate.AudioMixArray != nil {
			audioMixList := []interface{}{}
			for _, audioMix := range mediaConcatTemplate.ConcatTemplate.AudioMixArray {
				audioMixMap := map[string]interface{}{}

				if audioMix.AudioSource != "" {
					audioMixMap["audio_source"] = audioMix.AudioSource
				}

				if audioMix.MixMode != "" {
					audioMixMap["mix_mode"] = audioMix.MixMode
				}

				if audioMix.Replace != "" {
					audioMixMap["replace"] = audioMix.Replace
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

					audioMixMap["effect_config"] = []interface{}{effectConfigMap}
				}

				audioMixList = append(audioMixList, audioMixMap)
			}

			concatTemplateMap["audio_mix"] = audioMixList
		}

		_ = d.Set("concat_template", []interface{}{concatTemplateMap})
	}

	return nil
}

func resourceTencentCloudCiMediaConcatTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_concat_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := cos.CreateMediaConcatTemplateOptions{
		Tag: "Concat",
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if d.HasChange("concat_template") {
		if dMap, ok := helper.InterfacesHeadMap(d, "concat_template"); ok {
			concatTemplate := cos.ConcatTemplate{}
			if v, ok := dMap["concat_fragment"]; ok {
				for _, item := range v.([]interface{}) {
					concatFragmentMap := item.(map[string]interface{})
					concatFragment := cos.ConcatFragment{}
					if v, ok := concatFragmentMap["url"]; ok {
						concatFragment.Url = v.(string)
					}
					if v, ok := concatFragmentMap["mode"]; ok {
						concatFragment.Mode = v.(string)
					}
					concatTemplate.ConcatFragment = append(concatTemplate.ConcatFragment, concatFragment)
				}
			}
			if audioMap, ok := helper.InterfaceToMap(dMap, "audio"); ok {
				audio := cos.Audio{}
				if v, ok := audioMap["codec"]; ok {
					audio.Codec = v.(string)
				}
				if v, ok := audioMap["samplerate"]; ok {
					audio.Samplerate = v.(string)
				}
				if v, ok := audioMap["bitrate"]; ok {
					audio.Bitrate = v.(string)
				}
				if v, ok := audioMap["channels"]; ok {
					audio.Channels = v.(string)
				}
				concatTemplate.Audio = &audio
			}
			if videoMap, ok := helper.InterfaceToMap(dMap, "video"); ok {
				video := cos.Video{}
				if v, ok := videoMap["codec"]; ok {
					video.Codec = v.(string)
				}
				if v, ok := videoMap["width"]; ok {
					video.Width = v.(string)
				}
				if v, ok := videoMap["height"]; ok {
					video.Height = v.(string)
				}
				if v, ok := videoMap["bitrate"]; ok {
					video.Bitrate = v.(string)
				}
				if v, ok := videoMap["fps"]; ok {
					video.Fps = v.(string)
				}
				if v, ok := videoMap["crf"]; ok {
					video.Crf = v.(string)
				}
				if v, ok := videoMap["remove"]; ok {
					video.Remove = v.(string)
				}
				if v, ok := videoMap["rotate"]; ok {
					video.Rotate = v.(string)
				}
				concatTemplate.Video = &video
			}
			if containerMap, ok := helper.InterfaceToMap(dMap, "container"); ok {
				container := cos.Container{}
				if v, ok := containerMap["format"]; ok {
					container.Format = v.(string)
				}
				concatTemplate.Container = &container
			}
			if v, ok := dMap["audio_mix"]; ok {
				for _, item := range v.([]interface{}) {
					audioMixMap := item.(map[string]interface{})
					audioMix := cos.AudioMix{}
					if v, ok := audioMixMap["audio_source"]; ok {
						audioMix.AudioSource = v.(string)
					}
					if v, ok := audioMixMap["mix_mode"]; ok {
						audioMix.MixMode = v.(string)
					}
					if v, ok := audioMixMap["replace"]; ok {
						audioMix.Replace = v.(string)
					}
					if effectConfigMap, ok := helper.InterfaceToMap(audioMixMap, "effect_config"); ok {
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
					concatTemplate.AudioMixArray = append(concatTemplate.AudioMixArray, audioMix)
				}
			}
			request.ConcatTemplate = &concatTemplate
		}
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket).CI.UpdateMediaConcatTemplate(ctx, &request, templateId)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaConcatTemplate", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaConcatTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaConcatTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaConcatTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_concat_template.delete")()
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
