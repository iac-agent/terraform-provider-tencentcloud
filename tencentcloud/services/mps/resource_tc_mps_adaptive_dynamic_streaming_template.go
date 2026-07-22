package mps

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMpsAdaptiveDynamicStreamingTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsAdaptiveDynamicStreamingTemplateCreate,
		Read:   resourceTencentCloudMpsAdaptiveDynamicStreamingTemplateRead,
		Update: resourceTencentCloudMpsAdaptiveDynamicStreamingTemplateUpdate,
		Delete: resourceTencentCloudMpsAdaptiveDynamicStreamingTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"format": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Adaptive transcoding 格式，值 范围:HLS，MPEG-DASH。",
			},

			"stream_infos": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Convert adaptive 代码 流 到 output sub-流 参数 信息，和 output up 到 10 sub-streams.注意: frame 速率 的 each sub-流 必须 是 consistent; 如果 不， frame 速率 的 first sub-流 是 使用 作为 output frame 速率。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"video": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Video 参数 信息。",
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
										Description: "Adaptive resolution，可选 值:open: At 此 时间，宽度 表示 long side 的 视频，高度 表示 short side 的 视频.close: At 此 point，宽度 表示 宽度 的 视频，和 高度 表示 高度 的 视频.默认值：open.注意: In adaptive 模式，宽度 不能 是 smaller 比 高度。",
									},
									"width": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "最大 值 的 宽度 (或 long side) 的 视频 streaming，取值范围：0 和 [128，4096]，单位: 像素.当 宽度 和 高度 是 both 0， resolution 是 same.当 宽度 是 0 和 高度 是 不 0，宽度 是 scaled proportionally.当 宽度 是 不 0 和 高度 是 0，高度 是 scaled proportionally.当 both 宽度 和 高度 是 不 0， resolution 是 指定 通过 用户默认值：0。",
									},
									"height": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "最大 值 的 高度 (或 short side) 的 视频 streaming，取值范围：0 和 [128，4096]，单位: 像素.当 宽度 和 高度 是 both 0， resolution 是 same.当 宽度 是 0 和 高度 是 不 0，宽度 是 scaled proportionally.当 宽度 是 不 0 和 高度 是 0，高度 是 scaled proportionally.当 both 宽度 和 高度 是 不 0， resolution 是 指定 通过 用户默认值：0。",
									},
									"gop": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "间隔 between keyframe I frames，取值范围：0 和 [1，100000]，单位: 数量 frames.当 filling 0 或 不 filling， 系统 将 automatically 集合 gop 长度。",
									},
									"fill_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Filling 类型，当 aspect ratio 的 视频 流 配置 是 inconsistent 使用 aspect ratio 的 original 视频， processing 方法 对于 transcoding 是 filling. 可选 filling 类型:stretch: Stretching，stretching each frame 到 fill entire screen，其中 可能 cause transcoded 视频 到 是 squashed 或 stretched.black: Leave black，keep 视频 aspect ratio unchanged，和 fill rest 的 edge 使用 black.white: Leave blank，keep aspect ratio 的 视频，和 fill rest 的 edge 使用 white.gauss: Gaussian blur，keep aspect ratio 的 视频 unchanged，和 使用 Gaussian blur 对于 rest 的 edge.默认值：black.注意: Adaptive 流 仅 支持 stretch，black。",
									},
									"vcrf": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Video constant bit 速率 control factor， 值 范围 是 [1，51].如果 此 参数 是 指定， 代码 速率 control 方法 的 CRF 将 是 用于transcoding ( 视频 代码 速率 将 无 longer take effect).如果 there 是 无 special requirement，它 是 不 recommended 到 指定this 参数。",
									},
								},
							},
						},
						"audio": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Audio 参数 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"codec": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Encoding 格式 的 音频 流.当 outer 参数 Container 是 mp3， 可选 值 是:libmp3lame.当 outer 参数 Container 是 ogg 或 flac， 可选 值 是:flac.当 outer 参数 Container 是 m4a， 可选 值 是:libfdk_aac.libmp3lame.ac3.当 outer 参数 Container 是 mp4 或 flv， 可选 值 是:libfdk_aac: more suitable 对于 mp4.libmp3lame: more suitable 对于 flv.当 outer 参数 Container 是 hls， 可选 值 是:libfdk_aac.libmp3lame。",
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
						"remove_audio": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "是否remove 音频 流，值:0: reserved.1: remove。",
						},
						"remove_video": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "是否remove 视频 流，值:0: reserved.1: remove。",
						},
					},
				},
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "模板名称，长度 限制: 64 字符。",
			},

			"disable_higher_video_bitrate": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "是否prohibit 视频 从 low bit 速率 到 high bit 速率，值 范围:0: 无.1: yes.默认值：0。",
			},

			"disable_higher_video_resolution": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "是否prohibit conversion 的 视频 resolution 到 high resolution，值 范围:0: 无.1: yes.默认值：0。",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "模板描述 信息，长度 限制: 256 字符。",
			},

			"pure_audio": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "表示是否it 是 音频-仅. 0 表示 视频 template，1 表示 音频-仅 template.\nWhen 值 是 1.\n1. StreamInfos.N.RemoveVideo=1\n2. StreamInfos.N.RemoveAudio=0\n3. StreamInfos.N.Video.Codec=copy\nWhen 值 是 0.\n1. StreamInfos.N.Video.Codec 不能 是 copy.\n2. StreamInfos.N.Video.Fps 不能 是 null.\nNote: 此 值 仅 distinguishes template types. 任务 uses 值 的 RemoveAudio 和 RemoveVideo。",
			},

			"segment_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Segment 类型 有效值：\nts-segment: HLS+TS segment\nts-byterange: HLS+TS byte 范围\nmp4-segment: HLS+MP4 segment\nmp4-byterange: HLS/DASH+MP4 byte 范围\nts-packed-音频: HLS+TS+Packed Audio segment\nmp4-packed-音频: HLS+MP4+Packed Audio segment\nts-ts-segment: HLS+TS+TS segment\nts-ts-byterange: HLS+TS+TS byte 范围\nmp4-mp4-segment: HLS+MP4+MP4 segment\nmp4-mp4-byterange: HLS/DASH+MP4+MP4 byte 范围\nts-packed-音频-byterange: HLS+TS+Packed Audio byte 范围\nmp4-packed-音频-byterange: HLS+MP4+Packed Audio byte 范围.\n 默认值：ts-segment. 注意: segment 格式 对于 adaptive bitrate streaming 是 determined 通过 此 字段. For DASH 格式，SegmentType 可以 仅 是 mp4-byterange 或 mp4-mp4-byterange。",
			},
		},
	}
}

func resourceTencentCloudMpsAdaptiveDynamicStreamingTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_adaptive_dynamic_streaming_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = mps.NewCreateAdaptiveDynamicStreamingTemplateRequest()
		response   = mps.NewCreateAdaptiveDynamicStreamingTemplateResponse()
		definition uint64
	)

	if v, ok := d.GetOk("format"); ok {
		request.Format = helper.String(v.(string))
	}

	if v, ok := d.GetOk("stream_infos"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			adaptiveStreamTemplate := mps.AdaptiveStreamTemplate{}
			if videoMap, ok := helper.InterfaceToMap(dMap, "video"); ok {
				videoTemplateInfo := mps.VideoTemplateInfo{}
				if v, ok := videoMap["codec"]; ok {
					videoTemplateInfo.Codec = helper.String(v.(string))
				}

				if v, ok := videoMap["fps"]; ok {
					videoTemplateInfo.Fps = helper.IntInt64(v.(int))
				}

				if v, ok := videoMap["bitrate"]; ok {
					videoTemplateInfo.Bitrate = helper.IntInt64(v.(int))
				}

				if v, ok := videoMap["resolution_adaptive"]; ok {
					videoTemplateInfo.ResolutionAdaptive = helper.String(v.(string))
				}

				if v, ok := videoMap["width"]; ok {
					videoTemplateInfo.Width = helper.IntUint64(v.(int))
				}

				if v, ok := videoMap["height"]; ok {
					videoTemplateInfo.Height = helper.IntUint64(v.(int))
				}

				if v, ok := videoMap["gop"]; ok {
					videoTemplateInfo.Gop = helper.IntUint64(v.(int))
				}

				if v, ok := videoMap["fill_type"]; ok {
					videoTemplateInfo.FillType = helper.String(v.(string))
				}

				if v, ok := videoMap["vcrf"]; ok {
					videoTemplateInfo.Vcrf = helper.IntUint64(v.(int))
				}

				adaptiveStreamTemplate.Video = &videoTemplateInfo
			}

			if audioMap, ok := helper.InterfaceToMap(dMap, "audio"); ok {
				audioTemplateInfo := mps.AudioTemplateInfo{}
				if v, ok := audioMap["codec"]; ok {
					audioTemplateInfo.Codec = helper.String(v.(string))
				}

				if v, ok := audioMap["bitrate"]; ok {
					audioTemplateInfo.Bitrate = helper.IntInt64(v.(int))
				}

				if v, ok := audioMap["sample_rate"]; ok {
					audioTemplateInfo.SampleRate = helper.IntUint64(v.(int))
				}

				if v, ok := audioMap["audio_channel"]; ok {
					audioTemplateInfo.AudioChannel = helper.IntInt64(v.(int))
				}

				adaptiveStreamTemplate.Audio = &audioTemplateInfo
			}

			if v, ok := dMap["remove_audio"]; ok {
				adaptiveStreamTemplate.RemoveAudio = helper.IntUint64(v.(int))
			}

			if v, ok := dMap["remove_video"]; ok {
				adaptiveStreamTemplate.RemoveVideo = helper.IntUint64(v.(int))
			}

			request.StreamInfos = append(request.StreamInfos, &adaptiveStreamTemplate)
		}
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("disable_higher_video_bitrate"); ok {
		request.DisableHigherVideoBitrate = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("disable_higher_video_resolution"); ok {
		request.DisableHigherVideoResolution = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("pure_audio"); ok {
		request.PureAudio = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("segment_type"); ok {
		request.SegmentType = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().CreateAdaptiveDynamicStreamingTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create mps adaptive dynamic streaming template failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mps adaptiveDynamicStreamingTemplate failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.Definition == nil {
		return fmt.Errorf("Definition is nil.")
	}

	definition = *response.Response.Definition
	d.SetId(helper.UInt64ToStr(definition))

	return resourceTencentCloudMpsAdaptiveDynamicStreamingTemplateRead(d, meta)
}

func resourceTencentCloudMpsAdaptiveDynamicStreamingTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_adaptive_dynamic_streaming_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		definition = d.Id()
	)

	adaptiveDynamicStreamingTemplate, err := service.DescribeMpsAdaptiveDynamicStreamingTemplateById(ctx, definition)
	if err != nil {
		return err
	}

	if adaptiveDynamicStreamingTemplate == nil {
		log.Printf("[WARN]%s resource `tencentcloud_mps_adaptive_dynamic_streaming_template` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if adaptiveDynamicStreamingTemplate.Format != nil {
		_ = d.Set("format", adaptiveDynamicStreamingTemplate.Format)
	}

	if adaptiveDynamicStreamingTemplate.StreamInfos != nil {
		streamInfosList := []interface{}{}
		for _, streamInfos := range adaptiveDynamicStreamingTemplate.StreamInfos {
			streamInfosMap := map[string]interface{}{}
			if streamInfos.Video != nil {
				videoMap := map[string]interface{}{}
				if streamInfos.Video.Codec != nil {
					videoMap["codec"] = streamInfos.Video.Codec
				}

				if streamInfos.Video.Fps != nil {
					videoMap["fps"] = streamInfos.Video.Fps
				}

				if streamInfos.Video.Bitrate != nil {
					videoMap["bitrate"] = streamInfos.Video.Bitrate
				}

				if streamInfos.Video.ResolutionAdaptive != nil {
					videoMap["resolution_adaptive"] = streamInfos.Video.ResolutionAdaptive
				}

				if streamInfos.Video.Width != nil {
					videoMap["width"] = streamInfos.Video.Width
				}

				if streamInfos.Video.Height != nil {
					videoMap["height"] = streamInfos.Video.Height
				}

				if streamInfos.Video.Gop != nil {
					videoMap["gop"] = streamInfos.Video.Gop
				}

				if streamInfos.Video.FillType != nil {
					videoMap["fill_type"] = streamInfos.Video.FillType
				}

				if streamInfos.Video.Vcrf != nil {
					videoMap["vcrf"] = streamInfos.Video.Vcrf
				}

				streamInfosMap["video"] = []interface{}{videoMap}
			}

			if streamInfos.Audio != nil {
				audioMap := map[string]interface{}{}
				if streamInfos.Audio.Codec != nil {
					audioMap["codec"] = streamInfos.Audio.Codec
				}

				if streamInfos.Audio.Bitrate != nil {
					audioMap["bitrate"] = streamInfos.Audio.Bitrate
				}

				if streamInfos.Audio.SampleRate != nil {
					audioMap["sample_rate"] = streamInfos.Audio.SampleRate
				}

				if streamInfos.Audio.AudioChannel != nil {
					audioMap["audio_channel"] = streamInfos.Audio.AudioChannel
				}

				streamInfosMap["audio"] = []interface{}{audioMap}
			}

			if streamInfos.RemoveAudio != nil {
				streamInfosMap["remove_audio"] = streamInfos.RemoveAudio
			}

			if streamInfos.RemoveVideo != nil {
				streamInfosMap["remove_video"] = streamInfos.RemoveVideo
			}

			streamInfosList = append(streamInfosList, streamInfosMap)
		}

		_ = d.Set("stream_infos", streamInfosList)

	}

	if adaptiveDynamicStreamingTemplate.Name != nil {
		_ = d.Set("name", adaptiveDynamicStreamingTemplate.Name)
	}

	if adaptiveDynamicStreamingTemplate.DisableHigherVideoBitrate != nil {
		_ = d.Set("disable_higher_video_bitrate", adaptiveDynamicStreamingTemplate.DisableHigherVideoBitrate)
	}

	if adaptiveDynamicStreamingTemplate.DisableHigherVideoResolution != nil {
		_ = d.Set("disable_higher_video_resolution", adaptiveDynamicStreamingTemplate.DisableHigherVideoResolution)
	}

	if adaptiveDynamicStreamingTemplate.Comment != nil {
		_ = d.Set("comment", adaptiveDynamicStreamingTemplate.Comment)
	}

	if adaptiveDynamicStreamingTemplate.PureAudio != nil {
		_ = d.Set("pure_audio", adaptiveDynamicStreamingTemplate.PureAudio)
	}

	if adaptiveDynamicStreamingTemplate.SegmentType != nil {
		_ = d.Set("segment_type", adaptiveDynamicStreamingTemplate.SegmentType)
	}

	return nil
}

func resourceTencentCloudMpsAdaptiveDynamicStreamingTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_adaptive_dynamic_streaming_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = mps.NewModifyAdaptiveDynamicStreamingTemplateRequest()
		definition = d.Id()
	)

	needChange := false
	mutableArgs := []string{"format", "stream_infos", "name", "disable_higher_video_bitrate", "disable_higher_video_resolution", "comment", "pure_audio", "segment_type"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		if v, ok := d.GetOk("format"); ok {
			request.Format = helper.String(v.(string))
		}

		if v, ok := d.GetOk("stream_infos"); ok {
			for _, item := range v.([]interface{}) {
				adaptiveStreamTemplateMap := item.(map[string]interface{})
				adaptiveStreamTemplate := mps.AdaptiveStreamTemplate{}
				if videoMap, ok := helper.InterfaceToMap(adaptiveStreamTemplateMap, "video"); ok {
					videoTemplateInfo := mps.VideoTemplateInfo{}
					if v, ok := videoMap["codec"]; ok {
						videoTemplateInfo.Codec = helper.String(v.(string))
					}

					if v, ok := videoMap["fps"]; ok {
						videoTemplateInfo.Fps = helper.IntInt64(v.(int))
					}

					if v, ok := videoMap["bitrate"]; ok {
						videoTemplateInfo.Bitrate = helper.IntInt64(v.(int))
					}

					if v, ok := videoMap["resolution_adaptive"]; ok {
						videoTemplateInfo.ResolutionAdaptive = helper.String(v.(string))
					}

					if v, ok := videoMap["width"]; ok {
						videoTemplateInfo.Width = helper.IntUint64(v.(int))
					}

					if v, ok := videoMap["height"]; ok {
						videoTemplateInfo.Height = helper.IntUint64(v.(int))
					}

					if v, ok := videoMap["gop"]; ok {
						videoTemplateInfo.Gop = helper.IntUint64(v.(int))
					}

					if v, ok := videoMap["fill_type"]; ok {
						videoTemplateInfo.FillType = helper.String(v.(string))
					}

					if v, ok := videoMap["vcrf"]; ok {
						videoTemplateInfo.Vcrf = helper.IntUint64(v.(int))
					}

					adaptiveStreamTemplate.Video = &videoTemplateInfo
				}

				if audioMap, ok := helper.InterfaceToMap(adaptiveStreamTemplateMap, "audio"); ok {
					audioTemplateInfo := mps.AudioTemplateInfo{}
					if v, ok := audioMap["codec"]; ok {
						audioTemplateInfo.Codec = helper.String(v.(string))
					}

					if v, ok := audioMap["bitrate"]; ok {
						audioTemplateInfo.Bitrate = helper.IntInt64(v.(int))
					}

					if v, ok := audioMap["sample_rate"]; ok {
						audioTemplateInfo.SampleRate = helper.IntUint64(v.(int))
					}

					if v, ok := audioMap["audio_channel"]; ok {
						audioTemplateInfo.AudioChannel = helper.IntInt64(v.(int))
					}

					adaptiveStreamTemplate.Audio = &audioTemplateInfo
				}

				if v, ok := adaptiveStreamTemplateMap["remove_audio"]; ok {
					adaptiveStreamTemplate.RemoveAudio = helper.IntUint64(v.(int))
				}

				if v, ok := adaptiveStreamTemplateMap["remove_video"]; ok {
					adaptiveStreamTemplate.RemoveVideo = helper.IntUint64(v.(int))
				}

				request.StreamInfos = append(request.StreamInfos, &adaptiveStreamTemplate)
			}
		}

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("disable_higher_video_bitrate"); ok {
			request.DisableHigherVideoBitrate = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("disable_higher_video_resolution"); ok {
			request.DisableHigherVideoResolution = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("pure_audio"); ok {
			request.PureAudio = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("segment_type"); ok {
			request.SegmentType = helper.String(v.(string))
		}

		request.Definition = helper.StrToUint64Point(definition)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().ModifyAdaptiveDynamicStreamingTemplate(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update mps adaptiveDynamicStreamingTemplate failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudMpsAdaptiveDynamicStreamingTemplateRead(d, meta)
}

func resourceTencentCloudMpsAdaptiveDynamicStreamingTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_adaptive_dynamic_streaming_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		definition = d.Id()
	)

	if err := service.DeleteMpsAdaptiveDynamicStreamingTemplateById(ctx, definition); err != nil {
		return err
	}

	return nil
}
