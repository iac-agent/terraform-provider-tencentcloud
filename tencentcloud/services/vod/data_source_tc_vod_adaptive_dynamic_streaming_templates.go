package vod

import (
	"context"
	"fmt"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVodAdaptiveDynamicStreamingTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVodAdaptiveDynamicStreamingTemplatesRead,

		Schema: map[string]*schema.Schema{
			"definition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique ID filter of adaptive dynamic streaming template。",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Template 类型 filter. 有效值：`Preset`，`Custom`. `Preset`: preset template; `Custom`: custom template。",
			},
			"sub_app_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Subapplication ID in VOD. If you need to access a resource in a subapplication，enter the subapplication ID in this field; otherwise，leave it empty。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"template_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 adaptive dynamic streaming templates. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique ID adaptive dynamic streaming template。",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template 类型 filter. 有效值：`Preset`,`Custom`. `Preset`: preset template; `Custom`: custom template。",
						},
						"format": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Adaptive bitstream 格式",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "模板名称",
						},
						"drm_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DRM scheme 类型",
						},
						"disable_higher_video_bitrate": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否prohibit transcoding video from low bitrate to high bitrate. `false`: no，`true`: yes。",
						},
						"disable_higher_video_resolution": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否prohibit transcoding from low resolution to high resolution. `false`: no，`true`: yes。",
						},
						"comment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "模板描述",
						},
						"stream_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "列表 AdaptiveStreamTemplate parameter information of output substream for adaptive bitrate streaming。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"video": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Video parameter information。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"codec": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Video stream encoder. 有效值：`libx264`，`libx265`，`av1`.`libx264`: H.264，`libx265`: H.265，`av1`: AOMedia Video 1. Currently，a resolution within 640x480 must be specified for `H.265`. and the `av1` container only supports mp4。",
												},
												"fps": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Video frame rate in Hz. 取值范围：`[0，60]`. If the 值 is `0`，the frame rate will be the same as that of the 来源 video。",
												},
												"bitrate": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Bitrate of video stream in Kbps. 取值范围：`0` and `[128，35000]`. If the 值 is `0`，the bitrate of the video will be the same as that of the 来源 video。",
												},
												"resolution_adaptive": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Resolution adaption. 有效值：`true`,`false`. `true`: 已启用 In this case，`width` represents the long side of a video，while `height` the short side; `false`: 已禁用 In this case，`width` represents the width of a video，while `height` the height. Note: this field may return null，indicating that no valid values can be obtained。",
												},
												"width": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Maximum 值 of the width (or long side) of a video stream （像素）。 取值范围：`0` and `[128，4096]`. If both `width` and `height` are `0`，the resolution will be the same as that of the 来源 video; If `width` is `0`，but `height` is not `0`，`width` will be proportionally scaled; If `width` is not `0`，but `height` is `0`，`height` will be proportionally scaled; If both `width` and `height` are not `0`，the custom resolution will be used. Note: this field may return null，indicating that no valid values can be obtained。",
												},
												"height": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Maximum 值 of the height (or short side) of a video stream （像素）。 取值范围：`0` and `[128，4096]`. If both `width` and `height` are `0`，the resolution will be the same as that of the 来源 video; If `width` is `0`，but `height` is not `0`，`width` will be proportionally scaled; If `width` is not `0`，but `height` is `0`，`height` will be proportionally scaled; If both `width` and `height` are not `0`，the custom resolution will be used. Note: this field may return null，indicating that no valid values can be obtained。",
												},
												"fill_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Fill 类型 Fill refers to the way of processing a screenshot when its aspect ratio is different from that of the 来源 video. The following fill types are supported: `stretch`: stretch. The screenshot will be stretched frame by frame to match the aspect ratio of the 来源 video，which may make the screenshot shorter or longer; `black`: fill with black. This option retains the aspect ratio of the 来源 video for the screenshot and fills the unmatched area with black color blocks. Note: this field may return null，indicating that no valid values can be obtained。",
												},
											},
										},
									},
									"audio": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Audio parameter information。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"codec": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Audio stream encoder. Valid 值 are: `libfdk_aac` and `libmp3lame`。",
												},
												"bitrate": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Audio stream bitrate in Kbps. 取值范围：`0` and `[26，256]`. If the 值 is `0`，the bitrate of the audio stream will be the same as that of the original audio。",
												},
												"sample_rate": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Audio stream sample rate. 有效值：`32000`，`44100`，`48000`. Unit is HZ。",
												},
												"audio_channel": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: fmt.Sprintf("Audio channel system. Valid values: %s, %s, %s.", VOD_AUDIO_CHANNEL_MONO, VOD_AUDIO_CHANNEL_DUAL, VOD_AUDIO_CHANNEL_STEREO),
												},
											},
										},
									},
									"remove_audio": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否remove audio stream. `false`: no，`true`: yes。",
									},
								},
							},
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 of template in ISO date 格式",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "最后修改时间 of template in ISO date 格式",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudVodAdaptiveDynamicStreamingTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vod_adaptive_dynamic_streaming_templates.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	filter := make(map[string]interface{})
	if v, ok := d.GetOk("definition"); ok {
		filter["definitions"] = []string{v.(string)}
	}
	if v, ok := d.GetOk("type"); ok {
		filter["type"] = v.(string)
	}
	if v, ok := d.GetOk("sub_app_id"); ok {
		filter["sub_appid"] = v.(int)
	}

	vodService := VodService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	templates, err := vodService.DescribeAdaptiveDynamicStreamingTemplatesByFilter(ctx, filter)
	if err != nil {
		return err
	}

	templatesList := make([]map[string]interface{}, 0, len(templates))
	ids := make([]string, 0, len(templates))
	for _, item := range templates {
		templatesList = append(templatesList, func() map[string]interface{} {
			definitionStr := strconv.FormatUint(*item.Definition, 10)
			mapping := map[string]interface{}{
				"definition":                      definitionStr,
				"type":                            item.Type,
				"format":                          item.Format,
				"name":                            item.Name,
				"drm_type":                        item.DrmType,
				"disable_higher_video_bitrate":    *item.DisableHigherVideoBitrate == 1,
				"disable_higher_video_resolution": *item.DisableHigherVideoResolution == 1,
				"comment":                         item.Comment,
				"create_time":                     item.CreateTime,
				"update_time":                     item.UpdateTime,
			}
			var streamInfos = make([]interface{}, 0, len(item.StreamInfos))
			for _, v := range item.StreamInfos {
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
				})
			}
			mapping["stream_info"] = streamInfos
			ids = append(ids, definitionStr)
			return mapping
		}())
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("template_list", templatesList); e != nil {
		log.Printf("[CRITAL]%s provider set vod adaptive dynamic streaming template list fail, reason:%s ", logId, e.Error())
	}

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), templatesList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]", logId, output.(string), err.Error())
			return err
		}
	}

	return nil
}
