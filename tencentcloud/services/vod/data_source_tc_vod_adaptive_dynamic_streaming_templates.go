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
				Description: "Unique ID 过滤器 的 adaptive 动态 streaming template。",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "模板 类型 过滤器. 有效值：`Preset`，`Custom`. `Preset`: preset template; `Custom`: 自定义 template。",
			},
			"sub_app_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Subapplication ID 在 VOD. 如果 您 need 到 访问 资源 在 subapplication，enter subapplication ID 在 此 字段; otherwise，leave 它 空。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"template_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 adaptive 动态 streaming templates. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique ID adaptive 动态 streaming template。",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "模板 类型 过滤器. 有效值：`Preset`,`Custom`. `Preset`: preset template; `Custom`: 自定义 template。",
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
							Description: "是否prohibit transcoding 视频 从 low bitrate 到 high bitrate. `false`: 无，`true`: yes。",
						},
						"disable_higher_video_resolution": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否prohibit transcoding 从 low resolution 到 high resolution. `false`: 无，`true`: yes。",
						},
						"comment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "模板描述",
						},
						"stream_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "列表 AdaptiveStreamTemplate 参数 信息 的 output substream 对于 adaptive bitrate streaming。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"video": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Video 参数 信息。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"codec": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Video 流 encoder. 有效值：`libx264`，`libx265`，`av1`.`libx264`: H.264，`libx265`: H.265，`av1`: AOMedia Video 1. Currently， resolution within 640x480 必须 是 指定 对于 `H.265`. 和 `av1` 容器 仅 支持 mp4。",
												},
												"fps": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Video frame 速率 在 Hz. 取值范围：`[0，60]`. 如果 值 是 `0`， frame 速率 将 是 same 作为 该 的 来源 视频。",
												},
												"bitrate": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Bitrate 的 视频 流 在 Kbps. 取值范围：`0` 和 `[128，35000]`. 如果 值 是 `0`， bitrate 的 视频 将 是 same 作为 该 的 来源 视频。",
												},
												"resolution_adaptive": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Resolution adaption. 有效值：`true`,`false`. `true`: 已启用 In 此 case，`宽度` 表示 long side 的 视频，while `高度` short side; `false`: 已禁用 In 此 case，`宽度` 表示 宽度 的 视频，while `高度` 高度. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"width": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Maximum 值 的 宽度 (或 long side) 的 视频 流 （像素）。 取值范围：`0` 和 `[128，4096]`. 如果 both `宽度` 和 `高度` 是 `0`， resolution 将 是 same 作为 该 的 来源 视频; 如果 `宽度` 是 `0`，但 `高度` 是 不 `0`，`宽度` 将 是 proportionally scaled; 如果 `宽度` 是 不 `0`，但 `高度` 是 `0`，`高度` 将 是 proportionally scaled; 如果 both `宽度` 和 `高度` 是 不 `0`， 自定义 resolution 将 是 使用. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"height": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Maximum 值 的 高度 (或 short side) 的 视频 流 （像素）。 取值范围：`0` 和 `[128，4096]`. 如果 both `宽度` 和 `高度` 是 `0`， resolution 将 是 same 作为 该 的 来源 视频; 如果 `宽度` 是 `0`，但 `高度` 是 不 `0`，`宽度` 将 是 proportionally scaled; 如果 `宽度` 是 不 `0`，但 `高度` 是 `0`，`高度` 将 是 proportionally scaled; 如果 both `宽度` 和 `高度` 是 不 `0`， 自定义 resolution 将 是 使用. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
												"fill_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Fill 类型 Fill refers 到 way 的 processing screenshot 当 its aspect ratio 是 different 从 该 的 来源 视频. following fill types 是 支持: `stretch`: stretch. screenshot 将 是 stretched frame 通过 frame 到 match aspect ratio 的 来源 视频，其中 可能 make screenshot shorter 或 longer; `black`: fill 使用 black. 此 选项 retains aspect ratio 的 来源 视频 对于 screenshot 和 fills unmatched area 使用 black color blocks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
												},
											},
										},
									},
									"audio": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Audio 参数 信息。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"codec": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Audio 流 encoder. 有效 值 是: `libfdk_aac` 和 `libmp3lame`。",
												},
												"bitrate": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Audio 流 bitrate 在 Kbps. 取值范围：`0` 和 `[26，256]`. 如果 值 是 `0`， bitrate 的 音频 流 将 是 same 作为 该 的 original 音频。",
												},
												"sample_rate": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Audio 流 sample 速率. 有效值：`32000`，`44100`，`48000`. Unit 是 HZ。",
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
										Description: "是否remove 音频 流. `false`: 无，`true`: yes。",
									},
								},
							},
						},
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
