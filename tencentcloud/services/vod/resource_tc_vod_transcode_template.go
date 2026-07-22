package vod

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVodTranscodeTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVodTranscodeTemplateCreate,
		Read:   resourceTencentCloudVodTranscodeTemplateRead,
		Update: resourceTencentCloudVodTranscodeTemplateUpdate,
		Delete: resourceTencentCloudVodTranscodeTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"container": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "容器 格式 有效值：`mp4`，`flv`，`hls`，`mp3`，`flac`，`ogg`，`m4a`，`wav` ( `mp3`，`flac`，`ogg`，`m4a`，和 `wav` 是 音频 文件 formats)。",
			},

			"sub_app_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "VOD [应用](https://intl.云.tencent.com/document/product/266/14574) ID. For customers who activate VOD 服务 从 December 25，2023，如果 they want 到 访问 resources 在 VOD 应用 (whether 它's 默认值 应用 或 newly 创建 一个)，they 必须 fill 在 此 字段 使用 应用 ID。",
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Transcoding 模板名称 Length 限制: 64 字符。",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "模板描述 Length 限制: 256 字符。",
			},

			"remove_video": {
				Optional: true,
				Type:     schema.TypeInt,
				Description: "Whether 到 remove 视频 数据. 有效 值:\n" +
					"- 0: retain\n" +
					"- 1: remove\n" +
					"Default value: 0.",
			},

			"remove_audio": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "是否remove 音频 数据. 有效 值:0: retain 1: remove 默认值：0。",
			},

			"video_template": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Video 流 配置 参数. 此 字段 为必填项 当 `RemoveVideo` 是 0。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"codec": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "视频 codec. 有效 值:libx264: H.264; libx265: H.265; av1: AOMedia Video 1; H.266: H.266. AOMedia Video 1 和 H.266 codecs 可以 仅 是 用于MP4 files. Only CRF 是 支持 对于 H.266 currently。",
						},
						"fps": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Video frame 速率 在 Hz. 取值范围：[0,100].如果 值 是 0， frame 速率 将 是 same 作为 该 的 来源 视频。",
						},
						"bitrate": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Bitrate 的 视频 流 在 Kbps. 取值范围：0 和 [128，35,000].如果 值 是 0， bitrate 的 视频 将 是 same 作为 该 的 来源 视频。",
						},
						"resolution_adaptive": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resolution adaption. 有效 值:open: 已启用 In 此 case，`宽度` 表示 long side 的 视频，while `高度` short side;close: 已禁用 In 此 case，`宽度` 表示 宽度 的 视频，while `高度` 高度.默认值：open.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"width": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "最大 视频 宽度 (或 long side) 在 pixels. 取值范围：0 和 [128，8192].如果 both `宽度` 和 `高度` 是 0， output resolution 将 是 same 作为 该 的 来源 视频.如果 `宽度` 是 0 和 `高度` 是 不， 视频 宽度 将 是 proportionally scaled.如果 `宽度` 是 不 0 和 `高度` 是， 视频 高度 将 是 proportionally scaled.如果 neither `宽度` nor `高度` 是 0， 指定 宽度 和 高度 将 是 使用.默认值：0。",
						},
						"height": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "最大 视频 高度 (或 short side) 在 pixels. 取值范围：0 和 [128，8192].如果 both `宽度` 和 `高度` 是 0， output resolution 将 是 same 作为 该 的 来源 视频.如果 `宽度` 是 0 和 `高度` 是 不， 视频 宽度 将 是 proportionally scaled.如果 `宽度` 是 不 0 和 `高度` 是， 视频 高度 将 是 proportionally scaled.如果 neither `宽度` nor `高度` 是 0， 指定 宽度 和 高度 将 是 使用.默认值：0。",
						},
						"fill_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Fill 类型， way 的 processing screenshot 当 已配置 aspect ratio 是 different 从 该 的 来源 视频. 有效 值:stretch: stretches 视频 镜像 frame 通过 frame 到 fill screen. 视频 镜像 可能 become squashed 或 stretched after transcoding.black: fills uncovered area 使用 black color，without changing 镜像&#39;s aspect ratio.white: fills uncovered area 使用 white color，without changing 镜像&#39;s aspect ratio.gauss: applies Gaussian blur 到 uncovered area，without changing 镜像&#39;s aspect ratio.默认值：black。",
						},
						"vcrf": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "视频 constant 速率 factor (CRF). 取值范围：1-51.如果 此 参数 是 指定，CRF 编码 将 是 使用 和 bitrate 参数 将 是 ignored.如果 `Codec` 是 `H.266`，此 参数 为必填项 (`28` 是 recommended).We don't recommend 使用 此 参数 unless 您 have special requirements。",
						},
						"gop": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "I-frame 间隔 在 frames. 有效值：0 和 1-100000.当 此 参数 是 集合 到 0 或 left 空，`Gop` 将 是 automatically 集合。",
						},
						"preserve_hdr_switch": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "是否output HDR (high 动态 范围) 视频 如果 来源 视频 是 HDR. 有效 值:ON: 如果 来源 视频 是 HDR，output HDR 视频; 如果 不，output SDR (standard 动态 范围) 视频.OFF: Output SDR 视频 regardless 的 是否source 视频 是 HDR.默认值：OFF。",
						},
						"codec_tag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "codec 标签 此 参数 是 有效 仅 如果 H.265 codec 是 使用. 有效 值:hvc1hev1默认值：hvc1。",
						},
					},
				},
			},

			"audio_template": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Audio 流 配置 参数. 此 字段 为必填项 当 `RemoveAudio` 是 0。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"codec": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "音频 codec.如果 `Container` 是 `mp3`， 有效 值 是:`libmp3lame`如果 `Container` 是 `ogg` 或 `flac`， 有效 值 是:`flac`如果 `Container` 是 `m4a`， 有效 值 是:`libfdk_aac``libmp3lame``ac3`如果 `Container` 是 `mp4` 或 `flv`， 有效 值 是:`libfdk_aac` (Recommended 对于 MP4)`libmp3lame` (Recommended 对于 FLV)`mp2`如果 `Container` 是 `hls`， 有效 值 是:`libfdk_aac`如果 `格式` 是 `HLS` 或 `MPEG-DASH`， 有效 值 是:`libfdk_aac`如果 `Container` 是 `wav`， 有效 值 是:`pcm16`。",
						},
						"bitrate": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Audio 流 bitrate 在 Kbps. 取值范围：0 和 [26，256].如果 值 是 0， bitrate 的 音频 流 将 是 same 作为 该 的 original 音频。",
						},
						"sample_rate": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "音频 sample 速率. 有效 值:`16000` (有效 仅 如果 `Codec` 是 `pcm16`)`32000``44100``48000`单位：Hz。",
						},
						"audio_channel": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Audio channel 系统. 有效 值:1: mono-channel2: dual-channel6: stereoYou 不能 集合 sound channel 作为 stereo 对于 media files 在 容器 formats 对于 audios (FLAC，OGG，MP3，M4A).默认值：2。",
						},
					},
				},
			},

			"tehd_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "TESHD transcoding 参数。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "TESHD transcoding 类型 有效值：TEHD-100，OFF (默认值)。",
						},
						"max_video_bitrate": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Maximum bitrate，其中 是 有效 当 `类型` 是 `TESHD`.如果 此 参数 是 left blank 或 0 是 entered，there 将 是 无 upper 限制 对于 bitrate。",
						},
					},
				},
			},

			"segment_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "segment 类型 此 参数 是 有效 仅 如果 `Container` 是 `hls`. 有效值：`ts`: TS segment; `fmp4`: fMP4 segment 默认值：`ts`。",
			},
		},
	}
}

func resourceTencentCloudVodTranscodeTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_transcode_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = vod.NewCreateTranscodeTemplateRequest()
		response = vod.NewCreateTranscodeTemplateResponse()
		subAppId string
	)
	if v, ok := d.GetOk("container"); ok {
		request.Container = helper.String(v.(string))
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

	if v, ok := d.GetOkExists("remove_video"); ok {
		request.RemoveVideo = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("remove_audio"); ok {
		request.RemoveAudio = helper.IntInt64(v.(int))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "video_template"); ok {
		videoTemplateInfo := vod.VideoTemplateInfo{}
		if v, ok := dMap["codec"]; ok {
			videoTemplateInfo.Codec = helper.String(v.(string))
		}
		if v, ok := dMap["fps"]; ok {
			videoTemplateInfo.Fps = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["bitrate"]; ok {
			videoTemplateInfo.Bitrate = helper.IntUint64(v.(int))
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
		if v, ok := dMap["fill_type"]; ok {
			videoTemplateInfo.FillType = helper.String(v.(string))
		}
		if v, ok := dMap["vcrf"]; ok && v.(int) != 0 {
			videoTemplateInfo.Vcrf = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["gop"]; ok {
			videoTemplateInfo.Gop = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["preserve_hdr_switch"]; ok {
			videoTemplateInfo.PreserveHDRSwitch = helper.String(v.(string))
		}
		if v, ok := dMap["codec_tag"]; ok {
			videoTemplateInfo.CodecTag = helper.String(v.(string))
		}
		request.VideoTemplate = &videoTemplateInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "audio_template"); ok {
		audioTemplateInfo := vod.AudioTemplateInfo{}
		if v, ok := dMap["codec"]; ok {
			audioTemplateInfo.Codec = helper.String(v.(string))
		}
		if v, ok := dMap["bitrate"]; ok {
			audioTemplateInfo.Bitrate = helper.IntUint64(v.(int))
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
		tEHDConfig := vod.TEHDConfig{}
		if v, ok := dMap["type"]; ok {
			tEHDConfig.Type = helper.String(v.(string))
		}
		if v, ok := dMap["max_video_bitrate"]; ok {
			tEHDConfig.MaxVideoBitrate = helper.IntUint64(v.(int))
		}
		request.TEHDConfig = &tEHDConfig
	}

	if v, ok := d.GetOk("segment_type"); ok {
		request.SegmentType = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().CreateTranscodeTemplate(request)
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
		log.Printf("[CRITAL]%s create vod transcodeTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition := *response.Response.Definition
	d.SetId(subAppId + tccommon.FILED_SP + helper.Int64ToStr(definition))

	return resourceTencentCloudVodTranscodeTemplateRead(d, meta)
}

func resourceTencentCloudVodTranscodeTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_transcode_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("transcode template id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[0]
	definition := idSplit[1]

	transcodeTemplate, err := service.DescribeVodTranscodeTemplateById(ctx, helper.StrToUInt64(subAppId), helper.StrToInt64(definition))
	if err != nil {
		return err
	}

	if transcodeTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VodTranscodeTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if transcodeTemplate.Container != nil {
		_ = d.Set("container", transcodeTemplate.Container)
	}

	_ = d.Set("sub_app_id", helper.StrToInt(subAppId))

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

		if transcodeTemplate.VideoTemplate.FillType != nil {
			videoTemplateMap["fill_type"] = transcodeTemplate.VideoTemplate.FillType
		}

		if transcodeTemplate.VideoTemplate.Vcrf != nil && *transcodeTemplate.VideoTemplate.Vcrf != 0 {
			videoTemplateMap["vcrf"] = transcodeTemplate.VideoTemplate.Vcrf
		}

		if transcodeTemplate.VideoTemplate.Gop != nil {
			videoTemplateMap["gop"] = transcodeTemplate.VideoTemplate.Gop
		}

		if transcodeTemplate.VideoTemplate.PreserveHDRSwitch != nil {
			videoTemplateMap["preserve_hdr_switch"] = transcodeTemplate.VideoTemplate.PreserveHDRSwitch
		}

		if transcodeTemplate.VideoTemplate.CodecTag != nil {
			videoTemplateMap["codec_tag"] = transcodeTemplate.VideoTemplate.CodecTag
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

	if transcodeTemplate.SegmentType != nil {
		_ = d.Set("segment_type", transcodeTemplate.SegmentType)
	}

	return nil
}

func resourceTencentCloudVodTranscodeTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_transcode_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vod.NewModifyTranscodeTemplateRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("transcode template id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[0]
	definition := idSplit[1]

	request.SubAppId = helper.StrToUint64Point(subAppId)
	request.Definition = helper.StrToInt64Point(definition)

	immutableArgs := []string{"sub_app_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("container") || d.HasChange("name") || d.HasChange("comment") || d.HasChange("remove_video") || d.HasChange("remove_audio") || d.HasChange("video_template") || d.HasChange("audio_template") || d.HasChange("tehd_config") || d.HasChange("segment_type") {
		if v, ok := d.GetOk("container"); ok {
			request.Container = helper.String(v.(string))
		}
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}
		if v, ok := d.GetOkExists("remove_video"); ok {
			request.RemoveVideo = helper.IntInt64(v.(int))
		}
		if v, ok := d.GetOkExists("remove_audio"); ok {
			request.RemoveAudio = helper.IntInt64(v.(int))
		}
		if dMap, ok := helper.InterfacesHeadMap(d, "video_template"); ok {
			videoTemplateInfo := vod.VideoTemplateInfoForUpdate{}
			if v, ok := dMap["codec"]; ok {
				videoTemplateInfo.Codec = helper.String(v.(string))
			}
			if v, ok := dMap["fps"]; ok {
				videoTemplateInfo.Fps = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["bitrate"]; ok {
				videoTemplateInfo.Bitrate = helper.IntUint64(v.(int))
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
			if v, ok := dMap["fill_type"]; ok {
				videoTemplateInfo.FillType = helper.String(v.(string))
			}
			if v, ok := dMap["vcrf"]; ok {
				videoTemplateInfo.Vcrf = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["gop"]; ok {
				videoTemplateInfo.Gop = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["preserve_hdr_switch"]; ok {
				videoTemplateInfo.PreserveHDRSwitch = helper.String(v.(string))
			}
			if v, ok := dMap["codec_tag"]; ok {
				videoTemplateInfo.CodecTag = helper.String(v.(string))
			}
			request.VideoTemplate = &videoTemplateInfo
		}
		if dMap, ok := helper.InterfacesHeadMap(d, "audio_template"); ok {
			audioTemplateInfo := vod.AudioTemplateInfoForUpdate{}
			if v, ok := dMap["codec"]; ok {
				audioTemplateInfo.Codec = helper.String(v.(string))
			}
			if v, ok := dMap["bitrate"]; ok {
				audioTemplateInfo.Bitrate = helper.IntUint64(v.(int))
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
			tEHDConfig := vod.TEHDConfigForUpdate{}
			if v, ok := dMap["type"]; ok {
				tEHDConfig.Type = helper.String(v.(string))
			}
			if v, ok := dMap["max_video_bitrate"]; ok {
				tEHDConfig.MaxVideoBitrate = helper.IntUint64(v.(int))
			}
			request.TEHDConfig = &tEHDConfig
		}
		if v, ok := d.GetOk("segment_type"); ok {
			request.SegmentType = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().ModifyTranscodeTemplate(request)
		if e != nil {
			return resource.RetryableError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update vod transcodeTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudVodTranscodeTemplateRead(d, meta)
}

func resourceTencentCloudVodTranscodeTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_transcode_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("transcode template id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[0]
	definition := idSplit[1]

	if err := service.DeleteVodTranscodeTemplateById(ctx, helper.StrToUInt64(subAppId), helper.StrToInt64(definition)); err != nil {
		return err
	}

	return nil
}
