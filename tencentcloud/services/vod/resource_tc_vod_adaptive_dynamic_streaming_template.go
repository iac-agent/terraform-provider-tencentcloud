package vod

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudVodAdaptiveDynamicStreamingTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVodAdaptiveDynamicStreamingTemplateCreate,
		Read:   resourceTencentCloudVodAdaptiveDynamicStreamingTemplateRead,
		Update: resourceTencentCloudVodAdaptiveDynamicStreamingTemplateUpdate,
		Delete: resourceTencentCloudVodAdaptiveDynamicStreamingTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"format": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Adaptive bitstream 格式 有效值：`HLS`。",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 64),
				Description:  "模板名称 Length 限制: 64 字符。",
			},
			"drm_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "DRM scheme 类型 有效值：`SimpleAES`. 如果 此 字段 是 空 字符串，DRM 将 不 是 performed 在 视频。",
			},
			"disable_higher_video_bitrate": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "是否prohibit transcoding 视频 从 low bitrate 到 high bitrate. 有效值：`false`,`true`. `false`: 无，`true`: yes. 默认值：`false`。",
			},
			"disable_higher_video_resolution": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "是否prohibit transcoding 从 low resolution 到 high resolution. 有效值：`false`,`true`. `false`: 无，`true`: yes. 默认值：`false`。",
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 256),
				Description:  "模板描述 Length 限制: 256 字符。",
			},
			"sub_app_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "VOD [应用](https://intl.云.tencent.com/document/product/266/14574) ID. For customers who activate VOD 服务 从 December 25，2023，如果 they want 到 访问 resources 在 VOD 应用 (whether 它's 默认值 应用 或 newly 创建 一个)，they 必须 fill 在 此 字段 使用 应用 ID。",
			},
			"stream_info": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "列表 AdaptiveStreamTemplate 参数 信息 的 output substream 对于 adaptive bitrate streaming. Up 到 10 substreams 可以 是 output. 注意: frame 速率 的 all substreams 必须 是 same; otherwise， frame 速率 的 first substream 将 是 使用 作为 output frame 速率。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"video": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							MinItems:    1,
							Description: "Video 参数 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"codec": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Video 流 encoder. 有效值：`libx264`,`libx265`,`av1`. `libx264`: H.264，`libx265`: H.265，`av1`: AOMedia Video 1. Currently， resolution within 640x480 必须 是 指定 对于 `H.265`. 和 `av1` 容器 仅 支持 mp4。",
									},
									"fps": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: tccommon.ValidateIntegerInRange(0, 60),
										Description:  "Video frame 速率 在 Hz. 取值范围：`[0，60]`. 如果 值 是 `0`， frame 速率 将 是 same 作为 该 的 来源 视频。",
									},
									"bitrate": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Bitrate 的 视频 流 在 Kbps. 取值范围：`0` 和 `[128，35000]`. 如果 值 是 `0`， bitrate 的 视频 将 是 same 作为 该 的 来源 视频。",
									},
									"resolution_adaptive": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     true,
										Description: "Resolution adaption. 有效值：`true`,`false`. `true`: 已启用 In 此 case，`宽度` 表示 long side 的 视频，while `高度` short side; `false`: 已禁用 In 此 case，`宽度` 表示 宽度 的 视频，while `高度` 高度. 默认值：`true`. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
									},
									"width": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     0,
										Description: "Maximum 值 的 宽度 (或 long side) 的 视频 流 （像素）。 取值范围：`0` 和 `[128，4096]`. 如果 both `宽度` 和 `高度` 是 `0`， resolution 将 是 same 作为 该 的 来源 视频; 如果 `宽度` 是 `0`，但 `高度` 是 不 `0`，`宽度` 将 是 proportionally scaled; 如果 `宽度` 是 不 `0`，但 `高度` 是 `0`，`高度` 将 是 proportionally scaled; 如果 both `宽度` 和 `高度` 是 不 `0`， 自定义 resolution 将 是 使用. 默认值：`0`. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
									},
									"height": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     0,
										Description: "Maximum 值 的 高度 (或 short side) 的 视频 流 （像素）。 取值范围：`0` 和 `[128，4096]`. 如果 both `宽度` 和 `高度` 是 `0`， resolution 将 是 same 作为 该 的 来源 视频; 如果 `宽度` 是 `0`，但 `高度` 是 不 `0`，`宽度` 将 是 proportionally scaled; 如果 `宽度` 是 不 `0`，但 `高度` 是 `0`，`高度` 将 是 proportionally scaled; 如果 both `宽度` 和 `高度` 是 不 `0`， 自定义 resolution 将 是 使用. 默认值：`0`. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
									},
									"fill_type": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "black",
										ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"stretch", "black"}),
										Description:  "Fill 类型 Fill refers 到 way 的 processing screenshot 当 its aspect ratio 是 different 从 该 的 来源 视频. following fill types 是 支持: `stretch`: stretch. screenshot 将 是 stretched frame 通过 frame 到 match aspect ratio 的 来源 视频，其中 可能 make screenshot shorter 或 longer; `black`: fill 使用 black. 此 选项 retains aspect ratio 的 来源 视频 对于 screenshot 和 fills unmatched area 使用 black color blocks. 默认值：black. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
									},
									"vcrf": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
										Description: "Video constant bit 速率 control factor, 值 范围 是 [1,51].\n" +
											"Note:\n" +
											"- If this parameter is specified, the bitrate control method of CRF will be used for transcoding (the video bitrate will no longer take effect);\n" +
											"- This field is required when the video stream encoding format is H.266. The recommended value is 28;\n" +
											"- If there are no special requirements, it is not recommended to specify this parameter.",
									},
									"gop": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Interval between Keyframe I frames，取值范围：0 和 [1，100000]，单位: 数量 frames. 当 您 fill 在 0 或 leave 它 空， gop 长度 是 automatically 集合。",
									},
									"preserve_hdr_switch": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										Description: "Whether transcoding output still maintains HDR 当 original 视频 是 HDR (High Dynamic Range). Value 范围:\n" +
											"- ON: if the original file is HDR, the transcoding output remains HDR;, otherwise the transcoding output is SDR (Standard Dynamic Range);\n" +
											"- OFF: regardless of whether the original file is HDR or SDR, the transcoding output is SDR;\n" +
											"Default value: OFF.",
									},
									"codec_tag": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										Description: "Encoding label, 有效 仅 如果 编码 格式 的 视频 流 是 H.265 编码. Available 值:\n" +
											"- hvc1: stands for hvc1 tag;\n" +
											"- hev1: stands for the hev1 tag;\n" +
											"Default value: hvc1.",
									},
								},
							},
						},
						"audio": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							MinItems:    1,
							Description: "Audio 参数 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"codec": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Audio 流 encoder. 有效 值 是: `libfdk_aac` 和 `libmp3lame`. while `libfdk_aac` 是 recommended。",
									},
									"bitrate": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Audio 流 bitrate 在 Kbps. 取值范围：`0` 和 `[26，256]`. 如果 值 是 `0`， bitrate 的 音频 流 将 是 same 作为 该 的 original 音频。",
									},
									"sample_rate": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Audio 流 sample 速率. 有效值：`32000`，`44100`，`48000`Hz。",
									},
									"audio_channel": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     VOD_AUDIO_CHANNEL_DUAL,
										Description: fmt.Sprintf("Audio channel system. Valid values: %s, %s, %s. Default value: %s.", VOD_AUDIO_CHANNEL_MONO, VOD_AUDIO_CHANNEL_DUAL, VOD_AUDIO_CHANNEL_STEREO, VOD_AUDIO_CHANNEL_DUAL),
									},
								},
							},
						},
						"remove_audio": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "是否remove 音频 流. 有效值：`false`: 无，`true`: yes. `false` 通过 默认值。",
						},
						"remove_video": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "是否remove 视频 流. 有效值：`false`: 无，`true`: yes. `false` 通过 默认值。",
						},
						"tehd_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							MinItems:    1,
							Description: "Extremely fast HD transcoding 参数。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
										Description: "Extreme high-speed HD 类型, 可用 值:\n" +
											"- TEHD-100: super high definition-100th;\n" +
											"- OFF: turn off Ultra High definition.",
									},
									"max_video_bitrate": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Video bitrate 限制，其中 是 有效 当 类型 指定extreme speed HD 类型 如果 您 leave 它 空 或 enter 0，there 是 无 视频 bitrate 限制",
									},
								},
							},
						},
					},
				},
			},
			"segment_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Segment 类型, 有效 当 Format 是 HLS, 可选 值:\n" +
					"- ts: ts segment;\n" +
					"- fmp4: fmp4 segment;\n" +
					"Default value: ts.",
			},
			// computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "创建时间 的 template 在 ISO date 格式",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "最后修改时间 的 template 在 ISO date 格式",
			},
		},
	}
}

func resourceTencentCloudVodAdaptiveDynamicStreamingTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_adaptive_dynamic_streaming_template.create")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = vod.NewCreateAdaptiveDynamicStreamingTemplateRequest()
	)

	if v, ok := d.GetOk("segment_type"); ok {
		request.SegmentType = helper.String(v.(string))
	}
	request.Format = helper.String(d.Get("format").(string))
	request.Name = helper.String(d.Get("name").(string))
	if v, ok := d.GetOk("drm_type"); ok {
		request.DrmType = helper.String(v.(string))
	}
	request.DisableHigherVideoBitrate = helper.Uint64(DISABLE_HIGHER_VIDEO_BITRATE_TO_UNINT[d.Get("disable_higher_video_bitrate").(bool)])
	request.DisableHigherVideoResolution = helper.Uint64(DISABLE_HIGHER_VIDEO_RESOLUTION_TO_UNINT[d.Get("disable_higher_video_resolution").(bool)])
	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}
	var resourceId string
	if v, ok := d.GetOk("sub_app_id"); ok {
		subAppId := v.(int)
		resourceId += helper.IntToStr(subAppId)
		resourceId += tccommon.FILED_SP
		request.SubAppId = helper.IntUint64(subAppId)
	}
	streamInfos := d.Get("stream_info").([]interface{})
	request.StreamInfos = make([]*vod.AdaptiveStreamTemplate, 0, len(streamInfos))
	for _, item := range streamInfos {
		v := item.(map[string]interface{})
		video := v["video"].([]interface{})[0].(map[string]interface{})
		audio := v["audio"].([]interface{})[0].(map[string]interface{})
		rAudio := REMOVE_AUDIO_TO_UNINT[v["remove_audio"].(bool)]
		videoTemplateInfo := &vod.VideoTemplateInfo{
			Codec:              helper.String(video["codec"].(string)),
			Fps:                helper.IntUint64(video["fps"].(int)),
			Bitrate:            helper.IntUint64(video["bitrate"].(int)),
			ResolutionAdaptive: helper.String(RESOLUTION_ADAPTIVE_TO_STRING[video["resolution_adaptive"].(bool)]),
			Width:              helper.IntUint64(video["width"].(int)),
			Height:             helper.IntUint64(video["height"].(int)),
			FillType:           helper.String(video["fill_type"].(string)),
		}
		var rVideo uint64
		if v, ok := video["remove_video"]; ok && v.(bool) {
			rVideo = REMOVE_AUDIO_TO_UNINT[v.(bool)]
		}
		if v, ok := video["vcrf"]; ok && v.(int) != 0 {
			videoTemplateInfo.Vcrf = helper.IntUint64(v.(int))
		}
		if v, ok := video["gop"]; ok {
			videoTemplateInfo.Gop = helper.IntUint64(v.(int))
		}
		if v, ok := video["preserve_hdr_switch"]; ok && v.(string) != "" {
			videoTemplateInfo.PreserveHDRSwitch = helper.String(v.(string))
		}
		if v, ok := video["codec_tag"]; ok && v.(string) != "" {
			videoTemplateInfo.CodecTag = helper.String(v.(string))
		}

		var tehdConfig map[string]interface{}
		if len(v["tehd_config"].([]interface{})) > 0 {
			tehdConfig = v["tehd_config"].([]interface{})[0].(map[string]interface{})
		}
		request.StreamInfos = append(request.StreamInfos, &vod.AdaptiveStreamTemplate{

			Video: videoTemplateInfo,
			Audio: &vod.AudioTemplateInfo{
				Codec:        helper.String(audio["codec"].(string)),
				Bitrate:      helper.IntUint64(audio["bitrate"].(int)),
				SampleRate:   helper.IntUint64(audio["sample_rate"].(int)),
				AudioChannel: helper.Int64(VOD_AUDIO_CHANNEL_TYPE_TO_INT[audio["audio_channel"].(string)]),
			},
			RemoveAudio: &rAudio,
			RemoveVideo: &rVideo,
			TEHDConfig: func() *vod.TEHDConfig {
				if tehdConfig == nil {
					return nil
				}
				tehd := &vod.TEHDConfig{
					Type: helper.String(tehdConfig["type"].(string)),
				}
				if v, ok := tehdConfig["max_video_bitrate"]; ok {
					tehd.MaxVideoBitrate = helper.IntUint64(v.(int))
				}
				return tehd
			}(),
		})
	}

	var response *vod.CreateAdaptiveDynamicStreamingTemplateResponse
	var err error
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().CreateAdaptiveDynamicStreamingTemplate(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "FailedOperation" && sdkError.Message == "invalid vod user" {
					return resource.RetryableError(err)
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), strconv.ErrRange.Error())
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if response == nil || response.Response == nil {
		return fmt.Errorf("for vod adaptive dynamic streaming template creation, response is nil")
	}
	resourceId += strconv.FormatUint(*response.Response.Definition, 10)
	d.SetId(resourceId)

	return resourceTencentCloudVodAdaptiveDynamicStreamingTemplateRead(d, meta)
}

func resourceTencentCloudVodAdaptiveDynamicStreamingTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_adaptive_dynamic_streaming_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		subAppId   int
		definition string
		client     = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		vodService = VodService{client: client}
	)
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) == 2 {
		subAppId = helper.StrToInt(idSplit[0])
		definition = idSplit[1]
	} else {
		definition = d.Id()
	}
	// waiting for refreshing cache
	time.Sleep(30 * time.Second)
	template, has, err := vodService.DescribeAdaptiveDynamicStreamingTemplatesById(ctx, definition, subAppId)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("format", template.Format)
	_ = d.Set("name", template.Name)
	_ = d.Set("drm_type", template.DrmType)
	_ = d.Set("disable_higher_video_bitrate", *template.DisableHigherVideoBitrate == 1)
	_ = d.Set("disable_higher_video_resolution", *template.DisableHigherVideoResolution == 1)
	_ = d.Set("comment", template.Comment)
	_ = d.Set("create_time", template.CreateTime)
	_ = d.Set("update_time", template.UpdateTime)
	_ = d.Set("segment_type", template.SegmentType)

	var streamInfos = make([]interface{}, 0, len(template.StreamInfos))
	for _, v := range template.StreamInfos {
		streamInfos = append(streamInfos, map[string]interface{}{
			"video": []map[string]interface{}{
				{
					"codec":               v.Video.Codec,
					"fps":                 v.Video.Fps,
					"bitrate":             v.Video.Bitrate,
					"resolution_adaptive": *v.Video.ResolutionAdaptive == "open",
					"width":               v.Video.Width,
					"height":              v.Video.Height,
					"fill_type":           v.Video.FillType,
					"vcrf":                v.Video.Vcrf,
					"gop":                 v.Video.Gop,
					"preserve_hdr_switch": v.Video.PreserveHDRSwitch,
					"codec_tag":           v.Video.CodecTag,
				},
			},
			"audio": []map[string]interface{}{
				{
					"codec":         v.Audio.Codec,
					"bitrate":       v.Audio.Bitrate,
					"sample_rate":   v.Audio.SampleRate,
					"audio_channel": VOD_AUDIO_CHANNEL_TYPE_TO_STRING[*v.Audio.AudioChannel],
				},
			},
			"remove_audio": *v.RemoveAudio == 1,
			"remove_video": *v.RemoveVideo == 1,
			"tehd_config": func() []map[string]interface{} {
				if v.TEHDConfig == nil {
					return nil
				}
				return []map[string]interface{}{
					{
						"type":              v.TEHDConfig.Type,
						"max_video_bitrate": v.TEHDConfig.MaxVideoBitrate,
					},
				}
			}(),
		})
	}
	_ = d.Set("stream_info", streamInfos)
	if subAppId != 0 {
		_ = d.Set("sub_app_id", subAppId)
	}

	return nil
}

func resourceTencentCloudVodAdaptiveDynamicStreamingTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_adaptive_dynamic_streaming_template.update")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = vod.NewModifyAdaptiveDynamicStreamingTemplateRequest()
		changeFlag = false
		subAppId   int
		definition string
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) == 2 {
		subAppId = helper.StrToInt(idSplit[0])
		definition = idSplit[1]
		request.SubAppId = helper.IntUint64(subAppId)
	} else {
		definition = d.Id()
		if v, ok := d.GetOk("sub_app_id"); ok {
			request.SubAppId = helper.IntUint64(v.(int))
		}
	}

	request.Definition = helper.StrToUint64Point(definition)

	immutableArgs := []string{"sub_app_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("format") {
		changeFlag = true
		request.Format = helper.String(d.Get("format").(string))
	}
	if d.HasChange("name") {
		changeFlag = true
		request.Name = helper.String(d.Get("name").(string))
	}
	if d.HasChange("disable_higher_video_bitrate") {
		changeFlag = true
		request.DisableHigherVideoBitrate = helper.Uint64(DISABLE_HIGHER_VIDEO_BITRATE_TO_UNINT[d.Get("disable_higher_video_bitrate").(bool)])
	}
	if d.HasChange("disable_higher_video_resolution") {
		changeFlag = true
		request.DisableHigherVideoResolution = helper.Uint64(DISABLE_HIGHER_VIDEO_RESOLUTION_TO_UNINT[d.Get("disable_higher_video_resolution").(bool)])
	}
	if d.HasChange("comment") {
		changeFlag = true
		request.Comment = helper.String(d.Get("comment").(string))
	}
	if d.HasChange("stream_info") {
		changeFlag = true
		streamInfos := d.Get("stream_info").([]interface{})
		request.StreamInfos = make([]*vod.AdaptiveStreamTemplate, 0, len(streamInfos))
		for _, item := range streamInfos {
			v := item.(map[string]interface{})
			video := v["video"].([]interface{})[0].(map[string]interface{})
			audio := v["audio"].([]interface{})[0].(map[string]interface{})
			var tehdConfig map[string]interface{}
			if len(v["tehd_config"].([]interface{})) > 0 {
				tehdConfig = v["tehd_config"].([]interface{})[0].(map[string]interface{})
			}
			rAudio := REMOVE_AUDIO_TO_UNINT[v["remove_audio"].(bool)]
			var rVideo uint64
			if v, ok := video["remove_video"]; ok && v.(bool) {
				rVideo = REMOVE_AUDIO_TO_UNINT[v.(bool)]
			}
			request.StreamInfos = append(request.StreamInfos, &vod.AdaptiveStreamTemplate{
				Video: &vod.VideoTemplateInfo{
					Codec:              helper.String(video["codec"].(string)),
					Fps:                helper.IntUint64(video["fps"].(int)),
					Bitrate:            helper.IntUint64(video["bitrate"].(int)),
					ResolutionAdaptive: helper.String(RESOLUTION_ADAPTIVE_TO_STRING[video["resolution_adaptive"].(bool)]),
					Width: func(width int) *uint64 {
						if width == 0 {
							return nil
						}
						return helper.IntUint64(width)
					}(video["width"].(int)),
					Height: func(height int) *uint64 {
						if height == 0 {
							return nil
						}
						return helper.IntUint64(height)
					}(video["height"].(int)),
					FillType: helper.String(video["fill_type"].(string)),
					Vcrf: func() *uint64 {
						if v, ok := video["vcrf"]; !ok || v.(int) == 0 {
							return nil
						}
						return helper.IntUint64(video["vcrf"].(int))
					}(),
					Gop: func() *uint64 {
						if _, ok := video["gop"]; !ok {
							return nil
						}
						return helper.IntUint64(video["gop"].(int))
					}(),
					PreserveHDRSwitch: func() *string {
						if v, ok := video["preserve_hdr_switch"]; !ok || v.(string) == "" {
							return nil
						}
						return helper.String(video["preserve_hdr_switch"].(string))
					}(),
					CodecTag: func() *string {
						if v, ok := video["codec_tag"]; !ok || v.(string) == "" {
							return nil
						}
						return helper.String(video["codec_tag"].(string))
					}(),
				},
				Audio: &vod.AudioTemplateInfo{
					Codec:        helper.String(audio["codec"].(string)),
					Bitrate:      helper.IntUint64(audio["bitrate"].(int)),
					SampleRate:   helper.IntUint64(audio["sample_rate"].(int)),
					AudioChannel: helper.Int64(VOD_AUDIO_CHANNEL_TYPE_TO_INT[audio["audio_channel"].(string)]),
				},
				RemoveAudio: &rAudio,
				RemoveVideo: &rVideo,
				TEHDConfig: func() *vod.TEHDConfig {
					if tehdConfig == nil {
						return nil
					}
					tehd := &vod.TEHDConfig{
						Type: helper.String(tehdConfig["type"].(string)),
					}
					if v, ok := tehdConfig["max_video_bitrate"]; ok {
						tehd.MaxVideoBitrate = helper.IntUint64(v.(int))
					}
					return tehd
				}(),
			})
		}
	}

	if changeFlag {
		var err error
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			_, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().ModifyAdaptiveDynamicStreamingTemplate(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
				return tccommon.RetryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}

		return resourceTencentCloudVodAdaptiveDynamicStreamingTemplateRead(d, meta)
	}

	return nil
}

func resourceTencentCloudVodAdaptiveDynamicStreamingTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_adaptive_dynamic_streaming_template.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		subAppId   int
		definition string
	)
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) == 2 {
		subAppId = helper.StrToInt(idSplit[0])
		definition = idSplit[1]
	} else {
		definition = d.Id()
		if v, ok := d.GetOk("sub_app_id"); ok {
			subAppId = v.(int)
		}
	}
	vodService := VodService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	if err := vodService.DeleteAdaptiveDynamicStreamingTemplate(ctx, definition, uint64(subAppId)); err != nil {
		return err
	}

	return nil
}
