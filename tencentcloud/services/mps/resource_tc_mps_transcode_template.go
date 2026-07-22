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

func ResourceTencentCloudMpsTranscodeTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsTranscodeTemplateCreate,
		Read:   resourceTencentCloudMpsTranscodeTemplateRead,
		Update: resourceTencentCloudMpsTranscodeTemplateUpdate,
		Delete: resourceTencentCloudMpsTranscodeTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"container": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Encapsulation 格式，可选 值: mp4，flv，hls，mp3，flac，ogg，m4a. Among them，mp3，flac，ogg，m4a 是 pure 音频 files。",
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Transcoding 模板名称，长度 限制: 64 字符。",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "模板描述 信息，长度 限制: 256 字符。",
			},

			"remove_video": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "是否remove 视频 数据，值:0: reserved.1: remove.默认值：0。",
			},

			"remove_audio": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "是否remove 音频 数据，值:0: reserved.1: remove.默认值：0。",
			},

			"video_template": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Video 流 配置 参数，当 RemoveVideo 是 0，此 字段 为必填项。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"codec": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Encoding 格式 的 视频 流，可选 值:libx264: H.264 编码.libx265: H.265 编码.av1: AOMedia Video 1 编码.注意: Currently H.265 编码 必须 指定a resolution，和 它 needs 到 是 within 640*480.注意: av1 encoded containers currently 仅 support mp4。",
						},
						"fps": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Video frame 速率，取值范围：[0，100]，单位: Hz.当 值 是 0，它 表示 该 frame 速率 是 consistent 使用 original 视频.注意: 值 范围 对于 adaptive 代码 速率 是 [0，60]。",
						},
						"bitrate": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Bit 速率 的 视频 流，取值范围：0 和 [128，35000]，单位: kbps.当 值 是 0，它 表示 该 视频 bit 速率 是 consistent 使用 original 视频。",
						},
						"resolution_adaptive": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Adaptive resolution，可选 值:```open: open，在 此 时间，宽度 表示 long side 的 视频，高度 表示 short side 的 视频.close: close，在 此 时间，宽度 表示 宽度 的 视频，和 高度 表示 高度 的 视频.默认值：open.注意: In adaptive 模式，宽度 不能 是 smaller 比 高度。",
						},
						"width": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "最大 值 的 视频 流 宽度 (或 long side)，取值范围：0 和 [128，4096]，单位: 像素.当 宽度 和 高度 是 both 0， resolution 是 same.当 宽度 是 0 和 高度 是 不 0，宽度 是 scaled proportionally.当 宽度 是 不 0 和 高度 是 0，高度 是 scaled proportionally.当 both 宽度 和 高度 是 不 0， resolution 是 指定 通过 用户默认值：0。",
						},
						"height": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "最大 值 的 视频 流 高度 (或 short side)，取值范围：0 和 [128，4096]，单位: 像素.当 宽度 和 高度 是 both 0， resolution 是 same.当 宽度 是 0 和 高度 是 不 0，宽度 是 scaled proportionally.当 宽度 是 不 0 和 高度 是 0，高度 是 scaled proportionally.当 both 宽度 和 高度 是 不 0， resolution 是 指定 通过 用户默认值：0。",
						},
						"gop": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "间隔 between keyframe I frames，取值范围：0 和 [1，100000]，单位: 数量 frames.当 filling 0 或 不 filling， 系统 将 automatically 集合 gop 长度。",
						},
						"fill_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filling 方法，当 aspect ratio 的 视频 流 配置 是 inconsistent 使用 aspect ratio 的 original 视频， processing 方法 对于 transcoding 是 filling. 可选 filling 方法:stretch: Stretch，stretch each frame 到 fill entire screen，其中 可能 cause transcoded 视频 到 是 squashed 或 stretched.black: Leave black，keep aspect ratio 的 视频 unchanged，和 fill rest 的 edge 使用 black.white: Leave blank，keep aspect ratio 的 视频 unchanged，和 fill rest 的 edge 使用 white.gauss: Gaussian blur，keep aspect ratio 的 视频 unchanged，和 fill rest 的 edge 使用 Gaussian blur.默认值：black.注意: Adaptive 流 仅 支持 stretch，black。",
						},
						"vcrf": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Video constant bit 速率 control factor， 值 范围 是 [1，51].如果 此 参数 是 指定， 代码 速率 control 方法 的 CRF 将 是 用于transcoding ( 视频 代码 速率 将 无 longer take effect).如果 there 是 无 special requirement，它 是 不 recommended 到 指定this 参数。",
						},
					},
				},
			},

			"audio_template": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Audio 流 配置 参数，当 RemoveAudio 是 0，此 字段 为必填项。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"codec": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Encoding 格式 的 频率 流.当 outer 参数 Container 是 mp3， 可选 值 是:libmp3lame.当 outer 参数 Container 是 ogg 或 flac， 可选 值 是:flac.当 outer 参数 Container 是 m4a， 可选 值 是:libfdk_aac.libmp3lame.ac3.当 outer 参数 Container 是 mp4 或 flv， 可选 值 是:libfdk_aac: more suitable 对于 mp4.libmp3lame: more suitable 对于 flv.当 outer 参数 Container 是 hls， 可选 值 是:libfdk_aac.libmp3lame。",
						},
						"bitrate": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Bit 速率 的 音频 流，取值范围：0 和 [26，256]，单位: kbps.当 值 是 0，它 表示 该 音频 bit 速率 是 consistent 使用 original 音频。",
						},
						"sample_rate": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Sampling 速率 的 音频 流，可选 值32000.44100.48000.单位：Hz。",
						},
						"audio_channel": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Audio channel 模式，可选 值:`1: 单个 channel.2: Dual channel.6: Stereo.当 包 格式 的 media 是 音频 格式 (flac，ogg，mp3，m4a)， 数量 channels 是 不 allowed 到 是 集合 到 stereo.默认值：2。",
						},
					},
				},
			},

			"tehd_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Ultra-fast HD transcoding 参数。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Extremely high-definition 类型，可选 值:TEHD-100: Extreme HD-100.Not filling 表示 该 ultra-fast high-definition 是 不 已启用",
						},
						"max_video_bitrate": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "upper 限制 的 视频 bit 速率，其中 是 有效 当 类型 指定ultra-fast HD 类型Do 不 fill 在 或 fill 在 0 表示 该 there 是 无 upper 限制 在 视频 bit 速率。",
						},
					},
				},
			},

			"enhance_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Audio 和 视频 enhancement 配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"video_enhance": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Video Enhancement Configuration.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"frame_rate": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Interpolation frame 速率 配置.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability 配置 switch，可选 值: ON/OFF.默认值：ON。",
												},
												"fps": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Frame 速率，取值范围：[0，100]，单位: Hz.默认值：0.注意: For transcoding，此 参数 将 override Fps inside VideoTemplate.注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"super_resolution": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Super resolution 配置.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability 配置 switch，可选 值: ON/OFF.默认值：ON。",
												},
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "类型，可选 值:lq: super-resolution 对于 low-definition 视频 使用 more noise.hq: super resolution 对于 high-definition 视频.默认值：lq.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"size": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Super resolution 多个，可选 值:2: currently 仅 支持 2x super resolution.默认值：2.注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"hdr": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "HDR 配置.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability 配置 switch，可选 值: ON/OFF.默认值：ON。",
												},
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "类型，可选 值: HDR10/HLG.默认值：HDR10.注意: 编码 方法 的 视频 needs 到 是 libx265.注意: Video 编码 bit depth 是 10.注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"denoise": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Video Noise Reduction Configuration.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability 配置 switch，可选 值: ON/OFF.默认值：ON。",
												},
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "类型，可选 值: weak/strong.默认值：weak.注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"image_quality_enhance": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Comprehensive Enhanced Configuration.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability 配置 switch，可选 值: ON/OFF.默认值：ON。",
												},
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "类型，可选 值: weak/normal/strong.默认值：weak.注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"color_enhance": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Color Enhancement Configuration.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability 配置 switch，可选 值: ON/OFF.默认值：ON。",
												},
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "类型，可选 值: weak/normal/strong.默认值：weak.注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"sharp_enhance": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Detail Enhancement Configuration.注意：此字段可能返回 null，表示无法获取有效值。",
										Deprecated:  "It has been deprecated from version v1.82.67. Please do not use this again.",
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											return true
										},
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability 配置 switch，可选 值: ON/OFF.默认值：ON。",
												},
												"intensity": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "Intensity，取值范围：0.0~1.0.默认值：0.0.注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"face_enhance": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Face Enhancement Configuration.注意：此字段可能返回 null，表示无法获取有效值。",
										Deprecated:  "It has been deprecated from version v1.82.67. Please do not use this again.",
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											return true
										},
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability 配置 switch，可选 值: ON/OFF.默认值：ON。",
												},
												"intensity": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "Intensity，取值范围：0.0~1.0.默认值：0.0.注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"low_light_enhance": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Low Light Enhancement Configuration.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability 配置 switch，可选 值: ON/OFF.默认值：ON。",
												},
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "类型，可选 值: normal.默认值：normal.注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"scratch_repair": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "De-scratch 配置.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability 配置 switch，可选 值: ON/OFF.默认值：ON。",
												},
												"intensity": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "Intensity，取值范围：0.0~1.0.默认值：0.0.注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"artifact_repair": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "De-artifact (glitch) 配置.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Capability 配置 switch，可选 值: ON/OFF.默认值：ON。",
												},
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "类型，可选 值: weak/strong.默认值：weak.注意：此字段可能返回 null，表示无法获取有效值。",
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

func resourceTencentCloudMpsTranscodeTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_transcode_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = mps.NewCreateTranscodeTemplateRequest()
		response   = mps.NewCreateTranscodeTemplateResponse()
		definition int64
	)
	if v, ok := d.GetOk("container"); ok {
		request.Container = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	if v, _ := d.GetOk("remove_video"); v != nil {
		request.RemoveVideo = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("remove_audio"); v != nil {
		request.RemoveAudio = helper.IntInt64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "video_template"); ok {
		videoTemplateInfo := mps.VideoTemplateInfo{}
		if v, ok := dMap["codec"]; ok {
			videoTemplateInfo.Codec = helper.String(v.(string))
		}
		if v, ok := dMap["fps"]; ok {
			videoTemplateInfo.Fps = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["bitrate"]; ok {
			videoTemplateInfo.Bitrate = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["resolution_adaptive"]; ok {
			videoTemplateInfo.ResolutionAdaptive = helper.String(v.(string))
		}
		if v, ok := dMap["width"]; ok {
			videoTemplateInfo.Width = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["height"]; ok {
			videoTemplateInfo.Height = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["gop"]; ok {
			videoTemplateInfo.Gop = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["fill_type"]; ok {
			videoTemplateInfo.FillType = helper.String(v.(string))
		}
		if v, ok := dMap["vcrf"]; ok {
			videoTemplateInfo.Vcrf = helper.IntUint64(v.(int))
		}
		request.VideoTemplate = &videoTemplateInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "audio_template"); ok {
		audioTemplateInfo := mps.AudioTemplateInfo{}
		if v, ok := dMap["codec"]; ok {
			audioTemplateInfo.Codec = helper.String(v.(string))
		}
		if v, ok := dMap["bitrate"]; ok {
			audioTemplateInfo.Bitrate = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["sample_rate"]; ok {
			audioTemplateInfo.SampleRate = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["audio_channel"]; ok {
			audioTemplateInfo.AudioChannel = helper.IntInt64(v.(int))
		}
		request.AudioTemplate = &audioTemplateInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "tehd_config"); ok {
		tEHDConfig := mps.TEHDConfig{}
		if v, ok := dMap["type"]; ok {
			tEHDConfig.Type = helper.String(v.(string))
		}
		if v, ok := dMap["max_video_bitrate"]; ok {
			tEHDConfig.MaxVideoBitrate = helper.IntInt64(v.(int))
		}
		request.TEHDConfig = &tEHDConfig
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "enhance_config"); ok {
		enhanceConfig := mps.EnhanceConfig{}
		if videoEnhanceMap, ok := helper.InterfaceToMap(dMap, "video_enhance"); ok {
			videoEnhanceConfig := mps.VideoEnhanceConfig{}
			if frameRateMap, ok := helper.InterfaceToMap(videoEnhanceMap, "frame_rate"); ok {
				frameRateConfig := mps.FrameRateConfig{}
				if v, ok := frameRateMap["switch"]; ok {
					frameRateConfig.Switch = helper.String(v.(string))
				}
				if v, ok := frameRateMap["fps"]; ok {
					frameRateConfig.Fps = helper.IntUint64(v.(int))
				}
				videoEnhanceConfig.FrameRate = &frameRateConfig
			}
			if superResolutionMap, ok := helper.InterfaceToMap(videoEnhanceMap, "super_resolution"); ok {
				superResolutionConfig := mps.SuperResolutionConfig{}
				if v, ok := superResolutionMap["switch"]; ok {
					superResolutionConfig.Switch = helper.String(v.(string))
				}
				if v, ok := superResolutionMap["type"]; ok {
					superResolutionConfig.Type = helper.String(v.(string))
				}
				if v, ok := superResolutionMap["size"]; ok {
					superResolutionConfig.Size = helper.IntInt64(v.(int))
				}
				videoEnhanceConfig.SuperResolution = &superResolutionConfig
			}
			if hdrMap, ok := helper.InterfaceToMap(videoEnhanceMap, "hdr"); ok {
				hdrConfig := mps.HdrConfig{}
				if v, ok := hdrMap["switch"]; ok {
					hdrConfig.Switch = helper.String(v.(string))
				}
				if v, ok := hdrMap["type"]; ok {
					hdrConfig.Type = helper.String(v.(string))
				}
				videoEnhanceConfig.Hdr = &hdrConfig
			}
			if denoiseMap, ok := helper.InterfaceToMap(videoEnhanceMap, "denoise"); ok {
				videoDenoiseConfig := mps.VideoDenoiseConfig{}
				if v, ok := denoiseMap["switch"]; ok {
					videoDenoiseConfig.Switch = helper.String(v.(string))
				}
				if v, ok := denoiseMap["type"]; ok {
					videoDenoiseConfig.Type = helper.String(v.(string))
				}
				videoEnhanceConfig.Denoise = &videoDenoiseConfig
			}
			if imageQualityEnhanceMap, ok := helper.InterfaceToMap(videoEnhanceMap, "image_quality_enhance"); ok {
				imageQualityEnhanceConfig := mps.ImageQualityEnhanceConfig{}
				if v, ok := imageQualityEnhanceMap["switch"]; ok {
					imageQualityEnhanceConfig.Switch = helper.String(v.(string))
				}
				if v, ok := imageQualityEnhanceMap["type"]; ok {
					imageQualityEnhanceConfig.Type = helper.String(v.(string))
				}
				videoEnhanceConfig.ImageQualityEnhance = &imageQualityEnhanceConfig
			}
			if colorEnhanceMap, ok := helper.InterfaceToMap(videoEnhanceMap, "color_enhance"); ok {
				colorEnhanceConfig := mps.ColorEnhanceConfig{}
				if v, ok := colorEnhanceMap["switch"]; ok {
					colorEnhanceConfig.Switch = helper.String(v.(string))
				}
				if v, ok := colorEnhanceMap["type"]; ok {
					colorEnhanceConfig.Type = helper.String(v.(string))
				}
				videoEnhanceConfig.ColorEnhance = &colorEnhanceConfig
			}
			// if sharpEnhanceMap, ok := helper.InterfaceToMap(videoEnhanceMap, "sharp_enhance"); ok {
			// 	sharpEnhanceConfig := mps.SharpEnhanceConfig{}
			// 	if v, ok := sharpEnhanceMap["switch"]; ok {
			// 		sharpEnhanceConfig.Switch = helper.String(v.(string))
			// 	}
			// 	if v, ok := sharpEnhanceMap["intensity"]; ok {
			// 		sharpEnhanceConfig.Intensity = helper.Float64(v.(float64))
			// 	}
			// 	videoEnhanceConfig.SharpEnhance = &sharpEnhanceConfig
			// }
			// if faceEnhanceMap, ok := helper.InterfaceToMap(videoEnhanceMap, "face_enhance"); ok {
			// 	faceEnhanceConfig := mps.FaceEnhanceConfig{}
			// 	if v, ok := faceEnhanceMap["switch"]; ok {
			// 		faceEnhanceConfig.Switch = helper.String(v.(string))
			// 	}
			// 	if v, ok := faceEnhanceMap["intensity"]; ok {
			// 		faceEnhanceConfig.Intensity = helper.Float64(v.(float64))
			// 	}
			// 	videoEnhanceConfig.FaceEnhance = &faceEnhanceConfig
			// }
			if lowLightEnhanceMap, ok := helper.InterfaceToMap(videoEnhanceMap, "low_light_enhance"); ok {
				lowLightEnhanceConfig := mps.LowLightEnhanceConfig{}
				if v, ok := lowLightEnhanceMap["switch"]; ok {
					lowLightEnhanceConfig.Switch = helper.String(v.(string))
				}
				if v, ok := lowLightEnhanceMap["type"]; ok {
					lowLightEnhanceConfig.Type = helper.String(v.(string))
				}
				videoEnhanceConfig.LowLightEnhance = &lowLightEnhanceConfig
			}
			if scratchRepairMap, ok := helper.InterfaceToMap(videoEnhanceMap, "scratch_repair"); ok {
				scratchRepairConfig := mps.ScratchRepairConfig{}
				if v, ok := scratchRepairMap["switch"]; ok {
					scratchRepairConfig.Switch = helper.String(v.(string))
				}
				if v, ok := scratchRepairMap["intensity"]; ok {
					scratchRepairConfig.Intensity = helper.Float64(v.(float64))
				}
				videoEnhanceConfig.ScratchRepair = &scratchRepairConfig
			}
			if artifactRepairMap, ok := helper.InterfaceToMap(videoEnhanceMap, "artifact_repair"); ok {
				artifactRepairConfig := mps.ArtifactRepairConfig{}
				if v, ok := artifactRepairMap["switch"]; ok {
					artifactRepairConfig.Switch = helper.String(v.(string))
				}
				if v, ok := artifactRepairMap["type"]; ok {
					artifactRepairConfig.Type = helper.String(v.(string))
				}
				videoEnhanceConfig.ArtifactRepair = &artifactRepairConfig
			}
			enhanceConfig.VideoEnhance = &videoEnhanceConfig
		}
		request.EnhanceConfig = &enhanceConfig
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().CreateTranscodeTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps transcodeTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition = *response.Response.Definition
	d.SetId(helper.Int64ToStr(definition))

	return resourceTencentCloudMpsTranscodeTemplateRead(d, meta)
}

func resourceTencentCloudMpsTranscodeTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_transcode_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	definition := d.Id()

	transcodeTemplate, err := service.DescribeMpsTranscodeTemplateById(ctx, definition)
	if err != nil {
		return err
	}

	if transcodeTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsTranscodeTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if transcodeTemplate.Container != nil {
		_ = d.Set("container", transcodeTemplate.Container)
	}

	if transcodeTemplate.Name != nil {
		_ = d.Set("name", transcodeTemplate.Name)
	}

	if transcodeTemplate.Comment != nil {
		_ = d.Set("comment", transcodeTemplate.Comment)
	}

	if transcodeTemplate.RemoveVideo != nil {
		_ = d.Set("remove_video", transcodeTemplate.RemoveVideo)
	}

	if transcodeTemplate.RemoveAudio != nil {
		_ = d.Set("remove_audio", transcodeTemplate.RemoveAudio)
	}

	if transcodeTemplate.VideoTemplate != nil {
		videoTemplateMap := map[string]interface{}{}

		if transcodeTemplate.VideoTemplate.Codec != nil {
			videoTemplateMap["codec"] = transcodeTemplate.VideoTemplate.Codec
		}

		if transcodeTemplate.VideoTemplate.Fps != nil {
			videoTemplateMap["fps"] = transcodeTemplate.VideoTemplate.Fps
		}

		if transcodeTemplate.VideoTemplate.Bitrate != nil {
			videoTemplateMap["bitrate"] = transcodeTemplate.VideoTemplate.Bitrate
		}

		if transcodeTemplate.VideoTemplate.ResolutionAdaptive != nil {
			videoTemplateMap["resolution_adaptive"] = transcodeTemplate.VideoTemplate.ResolutionAdaptive
		}

		if transcodeTemplate.VideoTemplate.Width != nil {
			videoTemplateMap["width"] = transcodeTemplate.VideoTemplate.Width
		}

		if transcodeTemplate.VideoTemplate.Height != nil {
			videoTemplateMap["height"] = transcodeTemplate.VideoTemplate.Height
		}

		if transcodeTemplate.VideoTemplate.Gop != nil {
			videoTemplateMap["gop"] = transcodeTemplate.VideoTemplate.Gop
		}

		if transcodeTemplate.VideoTemplate.FillType != nil {
			videoTemplateMap["fill_type"] = transcodeTemplate.VideoTemplate.FillType
		}

		if transcodeTemplate.VideoTemplate.Vcrf != nil {
			videoTemplateMap["vcrf"] = transcodeTemplate.VideoTemplate.Vcrf
		}

		_ = d.Set("video_template", []interface{}{videoTemplateMap})
	}

	if transcodeTemplate.AudioTemplate != nil {
		audioTemplateMap := map[string]interface{}{}

		if transcodeTemplate.AudioTemplate.Codec != nil {
			audioTemplateMap["codec"] = transcodeTemplate.AudioTemplate.Codec
		}

		if transcodeTemplate.AudioTemplate.Bitrate != nil {
			audioTemplateMap["bitrate"] = transcodeTemplate.AudioTemplate.Bitrate
		}

		if transcodeTemplate.AudioTemplate.SampleRate != nil {
			audioTemplateMap["sample_rate"] = transcodeTemplate.AudioTemplate.SampleRate
		}

		if transcodeTemplate.AudioTemplate.AudioChannel != nil {
			audioTemplateMap["audio_channel"] = transcodeTemplate.AudioTemplate.AudioChannel
		}

		_ = d.Set("audio_template", []interface{}{audioTemplateMap})
	}

	if transcodeTemplate.TEHDConfig != nil {
		tEHDConfigMap := map[string]interface{}{}

		if transcodeTemplate.TEHDConfig.Type != nil {
			tEHDConfigMap["type"] = transcodeTemplate.TEHDConfig.Type
		}

		if transcodeTemplate.TEHDConfig.MaxVideoBitrate != nil {
			tEHDConfigMap["max_video_bitrate"] = transcodeTemplate.TEHDConfig.MaxVideoBitrate
		}

		_ = d.Set("tehd_config", []interface{}{tEHDConfigMap})
	}

	if transcodeTemplate.EnhanceConfig != nil {
		enhanceConfigMap := map[string]interface{}{}

		if transcodeTemplate.EnhanceConfig.VideoEnhance != nil {
			videoEnhanceMap := map[string]interface{}{}

			if transcodeTemplate.EnhanceConfig.VideoEnhance.FrameRate != nil {
				frameRateMap := map[string]interface{}{}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.FrameRate.Switch != nil {
					frameRateMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.FrameRate.Switch
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.FrameRate.Fps != nil {
					frameRateMap["fps"] = transcodeTemplate.EnhanceConfig.VideoEnhance.FrameRate.Fps
				}

				videoEnhanceMap["frame_rate"] = []interface{}{frameRateMap}
			}

			if transcodeTemplate.EnhanceConfig.VideoEnhance.SuperResolution != nil {
				superResolutionMap := map[string]interface{}{}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.SuperResolution.Switch != nil {
					superResolutionMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.SuperResolution.Switch
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.SuperResolution.Type != nil {
					superResolutionMap["type"] = transcodeTemplate.EnhanceConfig.VideoEnhance.SuperResolution.Type
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.SuperResolution.Size != nil {
					superResolutionMap["size"] = transcodeTemplate.EnhanceConfig.VideoEnhance.SuperResolution.Size
				}

				videoEnhanceMap["super_resolution"] = []interface{}{superResolutionMap}
			}

			if transcodeTemplate.EnhanceConfig.VideoEnhance.Hdr != nil {
				hdrMap := map[string]interface{}{}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.Hdr.Switch != nil {
					hdrMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.Hdr.Switch
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.Hdr.Type != nil {
					hdrMap["type"] = transcodeTemplate.EnhanceConfig.VideoEnhance.Hdr.Type
				}

				videoEnhanceMap["hdr"] = []interface{}{hdrMap}
			}

			if transcodeTemplate.EnhanceConfig.VideoEnhance.Denoise != nil {
				denoiseMap := map[string]interface{}{}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.Denoise.Switch != nil {
					denoiseMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.Denoise.Switch
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.Denoise.Type != nil {
					denoiseMap["type"] = transcodeTemplate.EnhanceConfig.VideoEnhance.Denoise.Type
				}

				videoEnhanceMap["denoise"] = []interface{}{denoiseMap}
			}

			if transcodeTemplate.EnhanceConfig.VideoEnhance.ImageQualityEnhance != nil {
				imageQualityEnhanceMap := map[string]interface{}{}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.ImageQualityEnhance.Switch != nil {
					imageQualityEnhanceMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.ImageQualityEnhance.Switch
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.ImageQualityEnhance.Type != nil {
					imageQualityEnhanceMap["type"] = transcodeTemplate.EnhanceConfig.VideoEnhance.ImageQualityEnhance.Type
				}

				videoEnhanceMap["image_quality_enhance"] = []interface{}{imageQualityEnhanceMap}
			}

			if transcodeTemplate.EnhanceConfig.VideoEnhance.ColorEnhance != nil {
				colorEnhanceMap := map[string]interface{}{}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.ColorEnhance.Switch != nil {
					colorEnhanceMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.ColorEnhance.Switch
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.ColorEnhance.Type != nil {
					colorEnhanceMap["type"] = transcodeTemplate.EnhanceConfig.VideoEnhance.ColorEnhance.Type
				}

				videoEnhanceMap["color_enhance"] = []interface{}{colorEnhanceMap}
			}

			// if transcodeTemplate.EnhanceConfig.VideoEnhance.SharpEnhance != nil {
			// 	sharpEnhanceMap := map[string]interface{}{}

			// 	if transcodeTemplate.EnhanceConfig.VideoEnhance.SharpEnhance.Switch != nil {
			// 		sharpEnhanceMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.SharpEnhance.Switch
			// 	}

			// 	if transcodeTemplate.EnhanceConfig.VideoEnhance.SharpEnhance.Intensity != nil {
			// 		sharpEnhanceMap["intensity"] = transcodeTemplate.EnhanceConfig.VideoEnhance.SharpEnhance.Intensity
			// 	}

			// 	videoEnhanceMap["sharp_enhance"] = []interface{}{sharpEnhanceMap}
			// }

			// if transcodeTemplate.EnhanceConfig.VideoEnhance.FaceEnhance != nil {
			// 	faceEnhanceMap := map[string]interface{}{}

			// 	if transcodeTemplate.EnhanceConfig.VideoEnhance.FaceEnhance.Switch != nil {
			// 		faceEnhanceMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.FaceEnhance.Switch
			// 	}

			// 	if transcodeTemplate.EnhanceConfig.VideoEnhance.FaceEnhance.Intensity != nil {
			// 		faceEnhanceMap["intensity"] = transcodeTemplate.EnhanceConfig.VideoEnhance.FaceEnhance.Intensity
			// 	}

			// 	videoEnhanceMap["face_enhance"] = []interface{}{faceEnhanceMap}
			// }

			if transcodeTemplate.EnhanceConfig.VideoEnhance.LowLightEnhance != nil {
				lowLightEnhanceMap := map[string]interface{}{}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.LowLightEnhance.Switch != nil {
					lowLightEnhanceMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.LowLightEnhance.Switch
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.LowLightEnhance.Type != nil {
					lowLightEnhanceMap["type"] = transcodeTemplate.EnhanceConfig.VideoEnhance.LowLightEnhance.Type
				}

				videoEnhanceMap["low_light_enhance"] = []interface{}{lowLightEnhanceMap}
			}

			if transcodeTemplate.EnhanceConfig.VideoEnhance.ScratchRepair != nil {
				scratchRepairMap := map[string]interface{}{}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.ScratchRepair.Switch != nil {
					scratchRepairMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.ScratchRepair.Switch
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.ScratchRepair.Intensity != nil {
					scratchRepairMap["intensity"] = transcodeTemplate.EnhanceConfig.VideoEnhance.ScratchRepair.Intensity
				}

				videoEnhanceMap["scratch_repair"] = []interface{}{scratchRepairMap}
			}

			if transcodeTemplate.EnhanceConfig.VideoEnhance.ArtifactRepair != nil {
				artifactRepairMap := map[string]interface{}{}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.ArtifactRepair.Switch != nil {
					artifactRepairMap["switch"] = transcodeTemplate.EnhanceConfig.VideoEnhance.ArtifactRepair.Switch
				}

				if transcodeTemplate.EnhanceConfig.VideoEnhance.ArtifactRepair.Type != nil {
					artifactRepairMap["type"] = transcodeTemplate.EnhanceConfig.VideoEnhance.ArtifactRepair.Type
				}

				videoEnhanceMap["artifact_repair"] = []interface{}{artifactRepairMap}
			}

			enhanceConfigMap["video_enhance"] = []interface{}{videoEnhanceMap}
		}

		_ = d.Set("enhance_config", []interface{}{enhanceConfigMap})
	}

	return nil
}

func resourceTencentCloudMpsTranscodeTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_transcode_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := mps.NewModifyTranscodeTemplateRequest()

	definition := d.Id()

	request.Definition = helper.StrToInt64Point(definition)

	if d.HasChange("container") {
		if v, ok := d.GetOk("container"); ok {
			request.Container = helper.String(v.(string))
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

	if d.HasChange("remove_video") {
		if v, _ := d.GetOk("remove_video"); v != nil {
			request.RemoveVideo = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("remove_audio") {
		if v, _ := d.GetOk("remove_audio"); v != nil {
			request.RemoveAudio = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("video_template") {
		if dMap, ok := helper.InterfacesHeadMap(d, "video_template"); ok {
			videoTemplateInfo := mps.VideoTemplateInfoForUpdate{}
			if v, ok := dMap["codec"]; ok {
				videoTemplateInfo.Codec = helper.String(v.(string))
			}
			if v, ok := dMap["fps"]; ok {
				videoTemplateInfo.Fps = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["bitrate"]; ok {
				videoTemplateInfo.Bitrate = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["resolution_adaptive"]; ok {
				videoTemplateInfo.ResolutionAdaptive = helper.String(v.(string))
			}
			if v, ok := dMap["width"]; ok {
				videoTemplateInfo.Width = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["height"]; ok {
				videoTemplateInfo.Height = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["gop"]; ok {
				videoTemplateInfo.Gop = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["fill_type"]; ok {
				videoTemplateInfo.FillType = helper.String(v.(string))
			}
			if v, ok := dMap["vcrf"]; ok {
				videoTemplateInfo.Vcrf = helper.IntUint64(v.(int))
			}
			request.VideoTemplate = &videoTemplateInfo
		}
	}

	if d.HasChange("audio_template") {
		if dMap, ok := helper.InterfacesHeadMap(d, "audio_template"); ok {
			audioTemplateInfo := mps.AudioTemplateInfoForUpdate{}
			if v, ok := dMap["codec"]; ok {
				audioTemplateInfo.Codec = helper.String(v.(string))
			}
			if v, ok := dMap["bitrate"]; ok {
				audioTemplateInfo.Bitrate = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["sample_rate"]; ok {
				audioTemplateInfo.SampleRate = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["audio_channel"]; ok {
				audioTemplateInfo.AudioChannel = helper.IntInt64(v.(int))
			}
			request.AudioTemplate = &audioTemplateInfo
		}
	}

	if d.HasChange("tehd_config") {
		if dMap, ok := helper.InterfacesHeadMap(d, "tehd_config"); ok {
			tEHDConfig := mps.TEHDConfigForUpdate{}
			if v, ok := dMap["type"]; ok {
				tEHDConfig.Type = helper.String(v.(string))
			}
			if v, ok := dMap["max_video_bitrate"]; ok {
				tEHDConfig.MaxVideoBitrate = helper.IntInt64(v.(int))
			}
			request.TEHDConfig = &tEHDConfig
		}
	}

	if d.HasChange("enhance_config") {
		if dMap, ok := helper.InterfacesHeadMap(d, "enhance_config"); ok {
			enhanceConfig := mps.EnhanceConfig{}
			if videoEnhanceMap, ok := helper.InterfaceToMap(dMap, "video_enhance"); ok {
				videoEnhanceConfig := mps.VideoEnhanceConfig{}
				if frameRateMap, ok := helper.InterfaceToMap(videoEnhanceMap, "frame_rate"); ok {
					frameRateConfig := mps.FrameRateConfig{}
					if v, ok := frameRateMap["switch"]; ok {
						frameRateConfig.Switch = helper.String(v.(string))
					}
					if v, ok := frameRateMap["fps"]; ok {
						frameRateConfig.Fps = helper.IntUint64(v.(int))
					}
					videoEnhanceConfig.FrameRate = &frameRateConfig
				}
				if superResolutionMap, ok := helper.InterfaceToMap(videoEnhanceMap, "super_resolution"); ok {
					superResolutionConfig := mps.SuperResolutionConfig{}
					if v, ok := superResolutionMap["switch"]; ok {
						superResolutionConfig.Switch = helper.String(v.(string))
					}
					if v, ok := superResolutionMap["type"]; ok {
						superResolutionConfig.Type = helper.String(v.(string))
					}
					if v, ok := superResolutionMap["size"]; ok {
						superResolutionConfig.Size = helper.IntInt64(v.(int))
					}
					videoEnhanceConfig.SuperResolution = &superResolutionConfig
				}
				if hdrMap, ok := helper.InterfaceToMap(videoEnhanceMap, "hdr"); ok {
					hdrConfig := mps.HdrConfig{}
					if v, ok := hdrMap["switch"]; ok {
						hdrConfig.Switch = helper.String(v.(string))
					}
					if v, ok := hdrMap["type"]; ok {
						hdrConfig.Type = helper.String(v.(string))
					}
					videoEnhanceConfig.Hdr = &hdrConfig
				}
				if denoiseMap, ok := helper.InterfaceToMap(videoEnhanceMap, "denoise"); ok {
					videoDenoiseConfig := mps.VideoDenoiseConfig{}
					if v, ok := denoiseMap["switch"]; ok {
						videoDenoiseConfig.Switch = helper.String(v.(string))
					}
					if v, ok := denoiseMap["type"]; ok {
						videoDenoiseConfig.Type = helper.String(v.(string))
					}
					videoEnhanceConfig.Denoise = &videoDenoiseConfig
				}
				if imageQualityEnhanceMap, ok := helper.InterfaceToMap(videoEnhanceMap, "image_quality_enhance"); ok {
					imageQualityEnhanceConfig := mps.ImageQualityEnhanceConfig{}
					if v, ok := imageQualityEnhanceMap["switch"]; ok {
						imageQualityEnhanceConfig.Switch = helper.String(v.(string))
					}
					if v, ok := imageQualityEnhanceMap["type"]; ok {
						imageQualityEnhanceConfig.Type = helper.String(v.(string))
					}
					videoEnhanceConfig.ImageQualityEnhance = &imageQualityEnhanceConfig
				}
				if colorEnhanceMap, ok := helper.InterfaceToMap(videoEnhanceMap, "color_enhance"); ok {
					colorEnhanceConfig := mps.ColorEnhanceConfig{}
					if v, ok := colorEnhanceMap["switch"]; ok {
						colorEnhanceConfig.Switch = helper.String(v.(string))
					}
					if v, ok := colorEnhanceMap["type"]; ok {
						colorEnhanceConfig.Type = helper.String(v.(string))
					}
					videoEnhanceConfig.ColorEnhance = &colorEnhanceConfig
				}
				// if sharpEnhanceMap, ok := helper.InterfaceToMap(videoEnhanceMap, "sharp_enhance"); ok {
				// 	sharpEnhanceConfig := mps.SharpEnhanceConfig{}
				// 	if v, ok := sharpEnhanceMap["switch"]; ok {
				// 		sharpEnhanceConfig.Switch = helper.String(v.(string))
				// 	}
				// 	if v, ok := sharpEnhanceMap["intensity"]; ok {
				// 		sharpEnhanceConfig.Intensity = helper.Float64(v.(float64))
				// 	}
				// 	videoEnhanceConfig.SharpEnhance = &sharpEnhanceConfig
				// }
				// if faceEnhanceMap, ok := helper.InterfaceToMap(videoEnhanceMap, "face_enhance"); ok {
				// 	faceEnhanceConfig := mps.FaceEnhanceConfig{}
				// 	if v, ok := faceEnhanceMap["switch"]; ok {
				// 		faceEnhanceConfig.Switch = helper.String(v.(string))
				// 	}
				// 	if v, ok := faceEnhanceMap["intensity"]; ok {
				// 		faceEnhanceConfig.Intensity = helper.Float64(v.(float64))
				// 	}
				// 	videoEnhanceConfig.FaceEnhance = &faceEnhanceConfig
				// }
				if lowLightEnhanceMap, ok := helper.InterfaceToMap(videoEnhanceMap, "low_light_enhance"); ok {
					lowLightEnhanceConfig := mps.LowLightEnhanceConfig{}
					if v, ok := lowLightEnhanceMap["switch"]; ok {
						lowLightEnhanceConfig.Switch = helper.String(v.(string))
					}
					if v, ok := lowLightEnhanceMap["type"]; ok {
						lowLightEnhanceConfig.Type = helper.String(v.(string))
					}
					videoEnhanceConfig.LowLightEnhance = &lowLightEnhanceConfig
				}
				if scratchRepairMap, ok := helper.InterfaceToMap(videoEnhanceMap, "scratch_repair"); ok {
					scratchRepairConfig := mps.ScratchRepairConfig{}
					if v, ok := scratchRepairMap["switch"]; ok {
						scratchRepairConfig.Switch = helper.String(v.(string))
					}
					if v, ok := scratchRepairMap["intensity"]; ok {
						scratchRepairConfig.Intensity = helper.Float64(v.(float64))
					}
					videoEnhanceConfig.ScratchRepair = &scratchRepairConfig
				}
				if artifactRepairMap, ok := helper.InterfaceToMap(videoEnhanceMap, "artifact_repair"); ok {
					artifactRepairConfig := mps.ArtifactRepairConfig{}
					if v, ok := artifactRepairMap["switch"]; ok {
						artifactRepairConfig.Switch = helper.String(v.(string))
					}
					if v, ok := artifactRepairMap["type"]; ok {
						artifactRepairConfig.Type = helper.String(v.(string))
					}
					videoEnhanceConfig.ArtifactRepair = &artifactRepairConfig
				}
				enhanceConfig.VideoEnhance = &videoEnhanceConfig
			}
			request.EnhanceConfig = &enhanceConfig
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().ModifyTranscodeTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mps transcodeTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMpsTranscodeTemplateRead(d, meta)
}

func resourceTencentCloudMpsTranscodeTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_transcode_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	definition := d.Id()

	if err := service.DeleteMpsTranscodeTemplateById(ctx, definition); err != nil {
		return err
	}

	return nil
}
