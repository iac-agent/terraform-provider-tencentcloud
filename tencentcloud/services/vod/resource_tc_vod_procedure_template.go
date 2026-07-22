package vod

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudVodProcedureTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVodProcedureTemplateCreate,
		Read:   resourceTencentCloudVodProcedureTemplateRead,
		Update: resourceTencentCloudVodProcedureTemplateUpdate,
		Delete: resourceTencentCloudVodProcedureTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 20),
				Description:  "任务 flow 名称 (up 到 20 字符)。",
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
			"media_process_task": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				MinItems:    1,
				Description: "Parameter 的 视频 processing 任务。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"transcode_task_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 transcoding tasks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Video transcoding 模板 ID",
									},
									"watermark_list": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    10,
										Description: "列表 up 到 `10` 镜像 或 text watermarks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
										Elem:        VodWatermarkResource(),
									},
									"mosaic_list": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    10,
										Description: "列表 blurs. Up 到 10 ones 可以 是 支持。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"coordinate_origin": {
													Type:        schema.TypeString,
													Optional:    true,
													Default:     "TopLeft",
													Description: "Origin position，其中 currently 可以 仅 是: `TopLeft`: 源站 的 coordinates 是 在 top-left corner 的 视频，和 源站 的 blur 是 在 top-left corner 的 镜像 或 text. 默认值：TopLeft。",
												},
												"x_pos": {
													Type:        schema.TypeString,
													Optional:    true,
													Default:     "0px",
													Description: "horizontal position 的 源站 的 blur relative 到 源站 的 coordinates 的 视频. `%` 和 `像素` formats 是 支持: 如果 字符串 结束 在 `%`， XPos 的 blur 将 是 指定 percentage 的 视频 宽度; 对于 示例，10% 表示 该 XPos 是 10% 的 视频 宽度; 如果 字符串 结束 在 `像素`， XPos 的 blur 将 是 指定 像素; 对于 示例，100px 表示 该 XPos 是 100 像素. 默认值：`0px`。",
												},
												"y_pos": {
													Type:        schema.TypeString,
													Optional:    true,
													Default:     "0px",
													Description: "Vertical position 的 源站 的 blur relative 到 源站 的 coordinates 的 视频. `%` 和 `像素` formats 是 支持: 如果 字符串 结束 在 `%`， YPos 的 blur 将 是 指定 percentage 的 视频 高度; 对于 示例，10% 表示 该 YPos 是 10% 的 视频 高度; 如果 字符串 结束 在 `像素`， YPos 的 blur 将 是 指定 像素; 对于 示例，100px 表示 该 YPos 是 100 像素. 默认值：`0px`。",
												},
												"width": {
													Type:        schema.TypeString,
													Optional:    true,
													Default:     "10%",
													Description: "Blur 宽度. `%` 和 `像素` formats 是 支持: 如果 字符串 结束 在 `%`， `宽度` 的 blur 将 是 指定 percentage 的 视频 宽度; 对于 示例，10% 表示 该 `宽度` 是 10% 的 视频 宽度; 如果 字符串 结束 在 `像素`， `宽度` 的 blur 将 是 在 像素; 对于 示例，100px 表示 该 宽度 是 100 像素. 默认值：`10%`。",
												},
												"height": {
													Type:        schema.TypeString,
													Optional:    true,
													Default:     "10%",
													Description: "Blur 高度. `%` 和 `像素` formats 是 支持: 如果 字符串 结束 在 `%`， `高度` 的 blur 将 是 指定 percentage 的 视频 高度; 对于 示例，10% 表示 该 高度 是 10% 的 视频 高度; 如果 字符串 结束 在 `像素`， `高度` 的 blur 将 是 在 像素; 对于 示例，100px 表示 该 高度 是 100 像素. 默认值：`10%`。",
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "开始时间 偏移量 的 blur （秒）。 如果此参数为空 或 `0` 是 entered， blur 将 appear upon first 视频 frame. 如果此参数为空 或 `0` 是 entered， blur 将 appear upon first 视频 frame; 如果 此 值 是 greater 比 `0` (e.g.，n)， blur 将 appear 在 second n after first 视频 frame; 如果 此 值 是 smaller 比 `0` (e.g.，-n)， blur 将 appear 在 second n before last 视频 frame。",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "结束时间 偏移量 的 blur （秒）。 如果此参数为空 或 `0` 是 entered， blur 将 exist till last 视频 frame; 如果 此 值 是 greater 比 `0` (e.g.，n)， blur 将 exist till second n; 如果 此 值 是 smaller 比 `0` (e.g.，-n)， blur 将 exist till second n before last 视频 frame。",
												},
											},
										},
									},
									"trace_watermark": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: "Digital 水印。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "是否use digital watermarks. 此 参数 为必填项. 有效值：ON，OFF。",
												},
											},
										},
									},
									"copy_right_watermark": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: "opyright 水印。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"text": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Copyright 信息，最大 长度 是 200 字符。",
												},
											},
										},
									},
									"head_tail_list": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "列表 视频 opening/closing credits 配置 template IDs. You 可以 enter up 到 10 IDs。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Video opening/closing credits 配置 模板 ID",
												},
											},
										},
									},
									"start_time_offset": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Computed:    true,
										Description: "开始时间 偏移量 的 blur （秒）。 如果此参数为空 或 `0` 是 entered， blur 将 appear upon first 视频 frame. 如果此参数为空 或 `0` 是 entered， blur 将 appear upon first 视频 frame; 如果 此 值 是 greater 比 `0` (e.g.，n)， blur 将 appear 在 second n after first 视频 frame; 如果 此 值 是 smaller 比 `0` (e.g.，-n)， blur 将 appear 在 second n before last 视频 frame。",
									},
									"end_time_offset": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Computed:    true,
										Description: "结束时间 偏移量 的 blur （秒）。 如果此参数为空 或 `0` 是 entered， blur 将 exist till last 视频 frame; 如果 此 值 是 greater 比 `0` (e.g.，n)， blur 将 exist till second n; 如果 此 值 是 smaller 比 `0` (e.g.，-n)， blur 将 exist till second n before last 视频 frame。",
									},
								},
							},
						},
						"animated_graphic_task_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 animated 镜像 generating tasks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Animated 镜像 generating 模板 ID",
									},
									"start_time_offset": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: "开始时间 的 animated 镜像 在 视频 （秒）。",
									},
									"end_time_offset": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: "结束时间 的 animated 镜像 在 视频 （秒）。",
									},
								},
							},
						},
						"snapshot_by_time_offset_task_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 时间 point screen capturing tasks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Time point screen capturing 模板 ID",
									},
									"ext_time_offset_list": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "列表 screenshot 时间 points. `s` 和 `%` formats 是 支持: 当 时间 point 字符串 结束 使用 `s`，its 单位 是 second. For 示例，`3.5s` 表示 3.5th second 的 视频; 当 时间 point 字符串 结束 使用 `%`，它 是 marked 使用 corresponding percentage 的 视频 时长. For 示例，`10%` 表示 该 时间 point 是 在 10% 的 视频 entire 时长。",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"watermark_list": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    10,
										Description: "列表 up 到 `10` 镜像 或 text watermarks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
										Elem:        VodWatermarkResource(),
									},
									"time_offset_list": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "列表 时间 points 对于 screencapturing （毫秒）。 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
										Elem: &schema.Schema{
											Type: schema.TypeFloat,
										},
									},
								},
							},
						},
						"sample_snapshot_task_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 sampled screen capturing tasks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Sampled screen capturing 模板 ID",
									},
									"watermark_list": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    10,
										Description: "列表 up 到 `10` 镜像 或 text watermarks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
										Elem:        VodWatermarkResource(),
									},
								},
							},
						},
						"image_sprite_task_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 镜像 sprite generating tasks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Image sprite generating 模板 ID",
									},
								},
							},
						},
						"cover_by_snapshot_task_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 cover generating tasks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Time point screen capturing 模板 ID",
									},
									"position_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"Time", "Percent"}),
										Description:  "Screen capturing 模式 有效值：`Time`，`Percent`. `Time`: screen captures 通过 时间 point，`Percent`: screen captures 通过 percentage。",
									},
									"position_value": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: "Screenshot position: For 时间 point screen capturing，此 表示 到 take screenshot 在 指定 时间 point (在 秒) 和 使用 它 作为 cover. For percentage screen capturing，此 值 表示 到 take screenshot 在 指定 percentage 的 视频 时长 和 使用 它 作为 cover。",
									},
									"watermark_list": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    10,
										Description: "列表 up 到 `10` 镜像 或 text watermarks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
										Elem:        VodWatermarkResource(),
									},
								},
							},
						},
						"adaptive_dynamic_streaming_task_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 adaptive bitrate streaming tasks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Adaptive bitrate streaming 模板 ID",
									},
									"watermark_list": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    10,
										Description: "列表 up 到 `10` 镜像 或 text watermarks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
										Elem:        VodWatermarkResource(),
									},
									"subtitle_list": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Subtitle 列表，element 是 subtitle ID，support 多个 subtitles，up 到 16。",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"ai_analysis_task": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Parameter 的 AI-based 内容 analysis 任务。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Video 内容 analysis 模板 ID",
						},
					},
				},
			},
			"ai_recognition_task": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "类型 参数 的 AI-based 内容 recognition 任务。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Intelligent 视频 recognition 模板 ID",
						},
					},
				},
			},
			"review_audio_video_task": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "类型 参数 的 AI-based 内容 recognition 任务。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"review_contents": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "类型 的 moderated 内容. 有效 值:\n" +
								"- `Media`: The original audio/video;\n" +
								"- `Cover`: Thumbnails.",
						},
						"definition": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Review template。",
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
			"type": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "模板 类型, 值 范围:\n" +
					"- Preset: system preset template;\n" +
					"- Custom: user-defined templates.",
			},
		},
	}
}

func genWatermarkList(item map[string]interface{}) (list []*vod.WatermarkInput) {
	if _, ok := item["watermark_list"]; !ok {
		return nil
	}
	waterV := item["watermark_list"].([]interface{})
	list = make([]*vod.WatermarkInput, 0, len(waterV))
	for _, water := range waterV {
		waterVV := water.(map[string]interface{})
		list = append(list, &vod.WatermarkInput{
			Definition: func(str string) *uint64 {
				idUint, _ := strconv.ParseUint(str, 0, 64)
				return &idUint
			}(waterVV["definition"].(string)),
			TextContent: func() *string {
				if content, ok := waterVV["text_content"]; !ok {
					return nil
				} else {
					return helper.String(content.(string))
				}
			}(),
			SvgContent: func() *string {
				if content, ok := waterVV["svg_content"]; !ok {
					return nil
				} else {
					return helper.String(content.(string))
				}
			}(),
			StartTimeOffset: func() *float64 {
				if content, ok := waterVV["start_time_offset"]; !ok {
					return nil
				} else {
					return helper.Float64(content.(float64))
				}
			}(),
			EndTimeOffset: func() *float64 {
				if content, ok := waterVV["end_time_offset"]; !ok {
					return nil
				} else {
					return helper.Float64(content.(float64))
				}
			}(),
		})
	}
	return list
}

func genMosaicList(item map[string]interface{}) (list []*vod.MosaicInput) {
	if _, ok := item["mosaic_list"]; !ok {
		return nil
	}
	mosaicV := item["mosaic_list"].([]interface{})
	list = make([]*vod.MosaicInput, 0, len(mosaicV))
	for _, mosaic := range mosaicV {
		mosaicVV := mosaic.(map[string]interface{})
		list = append(list, &vod.MosaicInput{
			CoordinateOrigin: func() *string {
				if content, ok := mosaicVV["coordinate_origin"]; !ok {
					return nil
				} else {
					return helper.String(content.(string))
				}
			}(),
			XPos: func() *string {
				if content, ok := mosaicVV["x_pos"]; !ok {
					return nil
				} else {
					return helper.String(content.(string))
				}
			}(),
			YPos: func() *string {
				if content, ok := mosaicVV["y_pos"]; !ok {
					return nil
				} else {
					return helper.String(content.(string))
				}
			}(),
			Width: func() *string {
				if content, ok := mosaicVV["width"]; !ok {
					return nil
				} else {
					return helper.String(content.(string))
				}
			}(),
			Height: func() *string {
				if content, ok := mosaicVV["height"]; !ok {
					return nil
				} else {
					return helper.String(content.(string))
				}
			}(),
			StartTimeOffset: func() *float64 {
				if content, ok := mosaicVV["start_time_offset"]; !ok {
					return nil
				} else {
					return helper.Float64(content.(float64))
				}
			}(),
			EndTimeOffset: func() *float64 {
				if content, ok := mosaicVV["end_time_offset"]; !ok {
					return nil
				} else {
					return helper.Float64(content.(float64))
				}
			}(),
		})
	}
	return list
}

func generateMediaProcessTask(d *schema.ResourceData) (mediaReq *vod.MediaProcessTaskInput) {
	mediaReq = &vod.MediaProcessTaskInput{}
	mediaProcessTask := d.Get("media_process_task").([]interface{})[0].(map[string]interface{})
	// transcode_task_list
	if transcodeTask, ok := mediaProcessTask["transcode_task_list"]; ok {
		transcodeTaskList := transcodeTask.([]interface{})
		transcodeReq := make([]*vod.TranscodeTaskInput, 0, len(transcodeTaskList))
		for _, itemV := range transcodeTaskList {
			item := itemV.(map[string]interface{})
			transcodeReq = append(transcodeReq, &vod.TranscodeTaskInput{
				Definition: func(str string) *uint64 {
					idUint, _ := strconv.ParseUint(str, 0, 64)
					return &idUint
				}(item["definition"].(string)),
				WatermarkSet: genWatermarkList(item),
				MosaicSet:    genMosaicList(item),
				TraceWatermark: func() *vod.TraceWatermarkInput {
					if _, ok := item["trace_watermark"]; !ok {
						return nil
					}
					traceWatermarks := item["trace_watermark"].([]interface{})
					if len(traceWatermarks) <= 0 {
						return nil
					}
					traceWatermarkInput := &vod.TraceWatermarkInput{}
					traceWatermarkItem := traceWatermarks[0].(map[string]interface{})
					if v, ok := traceWatermarkItem["switch"]; ok {
						traceWatermarkInput.Switch = helper.String(v.(string))
					}
					return traceWatermarkInput
				}(),
				CopyRightWatermark: func() *vod.CopyRightWatermarkInput {
					if _, ok := item["copy_right_watermark"]; !ok {
						return nil
					}
					copyRightWatermarks := item["copy_right_watermark"].([]interface{})
					if len(copyRightWatermarks) <= 0 {
						return nil
					}

					copyRightWatermarkInput := &vod.CopyRightWatermarkInput{}
					copyRightWatermarkItem := copyRightWatermarks[0].(map[string]interface{})
					if vv, ok := copyRightWatermarkItem["text"]; ok {
						copyRightWatermarkInput = &vod.CopyRightWatermarkInput{
							Text: helper.String(vv.(string)),
						}
					}
					return copyRightWatermarkInput
				}(),
				HeadTailSet: func() (list []*vod.HeadTailTaskInput) {
					if _, ok := item["head_tail_list"]; !ok {
						return
					}
					headTailSet := item["head_tail_list"].([]interface{})
					list = make([]*vod.HeadTailTaskInput, 0, len(headTailSet))
					for _, headTail := range headTailSet {
						headTailMap := headTail.(map[string]interface{})
						list = append(list, &vod.HeadTailTaskInput{
							Definition: func(str string) *int64 {
								definition, _ := strconv.ParseInt(str, 0, 64)
								return &definition
							}(headTailMap["definition"].(string)),
						})
					}
					return
				}(),
				StartTimeOffset: func() *float64 {
					if _, ok := item["start_time_offset"]; !ok {
						return nil
					}
					return helper.Float64(item["start_time_offset"].(float64))
				}(),
				EndTimeOffset: func() *float64 {
					if _, ok := item["end_time_offset"]; !ok {
						return nil
					}
					return helper.Float64(item["end_time_offset"].(float64))
				}(),
			})
		}
		mediaReq.TranscodeTaskSet = transcodeReq
	}
	// animated_graphic_task_list
	if animateTask, ok := mediaProcessTask["animated_graphic_task_list"]; ok {
		animateTaskList := animateTask.([]interface{})
		animateReq := make([]*vod.AnimatedGraphicTaskInput, 0, len(animateTaskList))
		for _, itemV := range animateTaskList {
			item := itemV.(map[string]interface{})
			animateReq = append(animateReq, &vod.AnimatedGraphicTaskInput{
				Definition: func(str string) *uint64 {
					idUint, _ := strconv.ParseUint(str, 0, 64)
					return &idUint
				}(item["definition"].(string)),
				StartTimeOffset: helper.Float64(item["start_time_offset"].(float64)),
				EndTimeOffset:   helper.Float64(item["end_time_offset"].(float64)),
			})
		}
		mediaReq.AnimatedGraphicTaskSet = animateReq
	}
	// snapshot_by_time_offset_task_list
	if snapshotTask, ok := mediaProcessTask["snapshot_by_time_offset_task_list"]; ok {
		snapshotTaskList := snapshotTask.([]interface{})
		snapshotReq := make([]*vod.SnapshotByTimeOffsetTaskInput, 0, len(snapshotTaskList))
		for _, itemV := range snapshotTaskList {
			item := itemV.(map[string]interface{})
			snapshotReq = append(snapshotReq, &vod.SnapshotByTimeOffsetTaskInput{
				Definition: func(str string) *uint64 {
					idUint, _ := strconv.ParseUint(str, 0, 64)
					return &idUint
				}(item["definition"].(string)),
				WatermarkSet: genWatermarkList(item),
				ExtTimeOffsetSet: func() (list []*string) {
					if _, ok := item["ext_time_offset_list"]; !ok {
						return nil
					}
					extTimeV := item["ext_time_offset_list"].([]interface{})
					list = make([]*string, 0, len(extTimeV))
					for _, extTimeVV := range extTimeV {
						list = append(list, helper.String(extTimeVV.(string)))
					}
					return list
				}(),
				TimeOffsetSet: func() (list []*float64) {
					if _, ok := item["time_offset_list"]; !ok {
						return nil
					}
					timeOffsetSet := item["time_offset_list"].([]interface{})
					list = make([]*float64, 0, len(timeOffsetSet))
					for _, timeOffset := range timeOffsetSet {
						list = append(list, helper.Float64(timeOffset.(float64)))
					}
					return list
				}(),
			})
		}
		mediaReq.SnapshotByTimeOffsetTaskSet = snapshotReq
	}
	// sample_snapshot_task_list
	if sampleTask, ok := mediaProcessTask["sample_snapshot_task_list"]; ok {
		sampleTaskList := sampleTask.([]interface{})
		sampleReq := make([]*vod.SampleSnapshotTaskInput, 0, len(sampleTaskList))
		for _, itemV := range sampleTaskList {
			item := itemV.(map[string]interface{})
			sampleReq = append(sampleReq, &vod.SampleSnapshotTaskInput{
				Definition: func(str string) *uint64 {
					idUint, _ := strconv.ParseUint(str, 0, 64)
					return &idUint
				}(item["definition"].(string)),
				WatermarkSet: genWatermarkList(item),
			})
		}
		mediaReq.SampleSnapshotTaskSet = sampleReq
	}
	// image_sprite_task_list
	if imageTask, ok := mediaProcessTask["image_sprite_task_list"]; ok {
		imageTaskList := imageTask.([]interface{})
		imageReq := make([]*vod.ImageSpriteTaskInput, 0, len(imageTaskList))
		for _, itemV := range imageTaskList {
			item := itemV.(map[string]interface{})
			imageReq = append(imageReq, &vod.ImageSpriteTaskInput{
				Definition: func(str string) *uint64 {
					idUint, _ := strconv.ParseUint(str, 0, 64)
					return &idUint
				}(item["definition"].(string)),
			})
		}
		mediaReq.ImageSpriteTaskSet = imageReq
	}
	// cover_by_snapshot_task_list
	if coverTask, ok := mediaProcessTask["cover_by_snapshot_task_list"]; ok {
		coverTaskList := coverTask.([]interface{})
		coverReq := make([]*vod.CoverBySnapshotTaskInput, 0, len(coverTaskList))
		for _, itemV := range coverTaskList {
			item := itemV.(map[string]interface{})
			coverReq = append(coverReq, &vod.CoverBySnapshotTaskInput{
				Definition: func(str string) *uint64 {
					idUint, _ := strconv.ParseUint(str, 0, 64)
					return &idUint
				}(item["definition"].(string)),
				WatermarkSet:  genWatermarkList(item),
				PositionType:  helper.String(item["iposition_type"].(string)),
				PositionValue: helper.Float64(item["position_value"].(float64)),
			})
		}
		mediaReq.CoverBySnapshotTaskSet = coverReq
	}
	// adaptive_dynamic_streaming_task_list
	if adaptiveTask, ok := mediaProcessTask["adaptive_dynamic_streaming_task_list"]; ok {
		adaptiveTaskList := adaptiveTask.([]interface{})
		adaptiveReq := make([]*vod.AdaptiveDynamicStreamingTaskInput, 0, len(adaptiveTaskList))
		for _, itemV := range adaptiveTaskList {
			item := itemV.(map[string]interface{})
			adaptiveReq = append(adaptiveReq, &vod.AdaptiveDynamicStreamingTaskInput{
				Definition: func(str string) *uint64 {
					idUint, _ := strconv.ParseUint(str, 0, 64)
					return &idUint
				}(item["definition"].(string)),
				WatermarkSet: genWatermarkList(item),
				SubtitleSet: func() (list []*string) {
					if _, ok := item["subtitle_list"]; !ok {
						return nil
					}
					subtitleList := item["subtitle_list"].([]interface{})
					list = make([]*string, 0, len(subtitleList))
					for _, subTitle := range subtitleList {
						list = append(list, helper.String(subTitle.(string)))
					}
					return list
				}(),
			})
		}
		mediaReq.AdaptiveDynamicStreamingTaskSet = adaptiveReq
	}

	return mediaReq
}

func resourceTencentCloudVodProcedureTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_procedure_template.create")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = vod.NewCreateProcedureTemplateRequest()
	)

	name := d.Get("name").(string)
	request.Name = helper.String(name)
	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	resourceId := name
	if v, ok := d.GetOk("sub_app_id"); ok {
		subAppId := v.(int)
		resourceId += tccommon.FILED_SP
		resourceId += helper.IntToStr(subAppId)
		request.SubAppId = helper.IntUint64(subAppId)
	}
	if _, ok := d.GetOk("media_process_task"); ok {
		mediaReq := generateMediaProcessTask(d)
		request.MediaProcessTask = mediaReq
	}
	//ai_analysis_task
	if aiAnalysisTask, ok := d.GetOk("ai_analysis_task"); ok {
		aiAnalysisTaskList := aiAnalysisTask.([]interface{})
		aiAnalysisTaskItem := aiAnalysisTaskList[0].(map[string]interface{})

		request.AiAnalysisTask = &vod.AiAnalysisTaskInput{
			Definition: helper.StrToUint64Point(aiAnalysisTaskItem["definition"].(string)),
		}
	}
	//ai_recognition_task
	if aiRecognitionTask, ok := d.GetOk("ai_recognition_task"); ok {
		aiRecognitionTaskList := aiRecognitionTask.([]interface{})
		aiRecognitionTaskItem := aiRecognitionTaskList[0].(map[string]interface{})

		request.AiRecognitionTask = &vod.AiRecognitionTaskInput{
			Definition: helper.StrToUint64Point(aiRecognitionTaskItem["definition"].(string)),
		}
	}
	//review_audio_video_task
	if reviewAudioVideoTask, ok := d.GetOk("review_audio_video_task"); ok {
		reviewAudioVideoTaskList := reviewAudioVideoTask.([]interface{})
		reviewAudioVideoTaskItem := reviewAudioVideoTaskList[0].(map[string]interface{})

		request.ReviewAudioVideoTask = &vod.ProcedureReviewAudioVideoTaskInput{
			Definition: helper.StrToUint64Point(reviewAudioVideoTaskItem["definition"].(string)),
			ReviewContents: func() (list []*string) {
				if _, ok := reviewAudioVideoTaskItem["review_contents"]; !ok {
					return
				}
				reviewContentList := reviewAudioVideoTaskItem["review_contents"].([]interface{})
				list = make([]*string, 0, len(reviewContentList))
				for _, reviewContent := range reviewContentList {
					list = append(list, helper.String(reviewContent.(string)))
				}
				return
			}(),
		}
	}
	var err error
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().CreateProcedureTemplate(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(resourceId)

	return resourceTencentCloudVodProcedureTemplateRead(d, meta)
}

func resourceTencentCloudVodProcedureTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_procedure_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		name       string
		subAppId   int
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		id         = d.Id()
		client     = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		vodService = VodService{client: client}
	)
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) == 2 {
		name = idSplit[0]
		subAppId = helper.StrToInt(idSplit[1])
	} else {
		name = id
	}
	template, has, err := vodService.DescribeProcedureTemplatesById(ctx, name, subAppId)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", template.Name)
	_ = d.Set("type", template.Type)
	_ = d.Set("comment", template.Comment)
	_ = d.Set("create_time", template.CreateTime)
	_ = d.Set("update_time", template.UpdateTime)
	if subAppId != 0 {
		_ = d.Set("sub_app_id", subAppId)
	}

	mediaProcessTaskElem := make(map[string]interface{})
	if template.MediaProcessTask != nil {
		// transcode_task_list
		if template.MediaProcessTask.TranscodeTaskSet != nil {
			list := make([]map[string]interface{}, 0, len(template.MediaProcessTask.TranscodeTaskSet))
			for _, item := range template.MediaProcessTask.TranscodeTaskSet {
				list = append(list, map[string]interface{}{
					"definition": strconv.FormatUint(*item.Definition, 10),
					"watermark_list": func() interface{} {
						if item.WatermarkSet == nil {
							return nil
						}
						waterList := make([]map[string]interface{}, 0, len(item.WatermarkSet))
						for _, waterV := range item.WatermarkSet {
							waterList = append(waterList, map[string]interface{}{
								"definition":        strconv.FormatUint(*waterV.Definition, 10),
								"text_content":      waterV.TextContent,
								"svg_content":       waterV.SvgContent,
								"start_time_offset": waterV.StartTimeOffset,
								"end_time_offset":   waterV.EndTimeOffset,
							})
						}
						return waterList
					}(),
					"mosaic_list": func() interface{} {
						if item.MosaicSet == nil {
							return nil
						}
						mosaicList := make([]map[string]interface{}, 0, len(item.MosaicSet))
						for _, mosaicV := range item.MosaicSet {
							mosaicList = append(mosaicList, map[string]interface{}{
								"coordinate_origin": mosaicV.CoordinateOrigin,
								"x_pos":             mosaicV.XPos,
								"y_pos":             mosaicV.YPos,
								"width":             mosaicV.Width,
								"height":            mosaicV.Height,
								"start_time_offset": mosaicV.StartTimeOffset,
								"end_time_offset":   mosaicV.EndTimeOffset,
							})
						}
						return mosaicList
					}(),
					"trace_watermark": func() interface{} {
						if item.TraceWatermark == nil {
							return nil
						}
						traceWatermark := map[string]interface{}{
							"switch": item.TraceWatermark.Switch,
						}
						return []interface{}{traceWatermark}
					}(),
					"copy_right_watermark": func() interface{} {
						if item.CopyRightWatermark == nil {
							return nil
						}
						copyRightWatermark := map[string]interface{}{
							"text": item.CopyRightWatermark.Text,
						}
						return []interface{}{copyRightWatermark}
					}(),
					"head_tail_list": func() interface{} {
						if item.HeadTailSet == nil {
							return nil
						}
						headTailList := make([]interface{}, 0, len(item.HeadTailSet))
						for _, headTail := range item.HeadTailSet {
							headTailList = append(headTailList, map[string]interface{}{
								"definition": headTail.Definition,
							})
						}
						return headTailList
					}(),
					"start_time_offset": func() *float64 {
						if item.StartTimeOffset == nil {
							return nil
						}
						return item.StartTimeOffset
					}(),
					"end_time_offset": func() *float64 {
						if item.EndTimeOffset == nil {
							return nil
						}
						return item.EndTimeOffset
					}(),
				})
			}
			mediaProcessTaskElem["transcode_task_list"] = list
		}
		// animated_graphic_task_list
		if template.MediaProcessTask.AnimatedGraphicTaskSet != nil {
			list := make([]map[string]interface{}, 0, len(template.MediaProcessTask.AnimatedGraphicTaskSet))
			for _, item := range template.MediaProcessTask.AnimatedGraphicTaskSet {
				list = append(list, map[string]interface{}{
					"definition":        strconv.FormatUint(*item.Definition, 10),
					"start_time_offset": item.StartTimeOffset,
					"end_time_offset":   item.EndTimeOffset,
				})
			}
			mediaProcessTaskElem["animated_graphic_task_list"] = list
		}
		// snapshot_by_time_offset_task_list
		if template.MediaProcessTask.SnapshotByTimeOffsetTaskSet != nil {
			list := make([]map[string]interface{}, 0, len(template.MediaProcessTask.SnapshotByTimeOffsetTaskSet))
			for _, item := range template.MediaProcessTask.SnapshotByTimeOffsetTaskSet {
				list = append(list, map[string]interface{}{
					"definition": strconv.FormatUint(*item.Definition, 10),
					"watermark_list": func() interface{} {
						if item.WatermarkSet == nil {
							return nil
						}
						waterList := make([]map[string]interface{}, 0, len(item.WatermarkSet))
						for _, waterV := range item.WatermarkSet {
							waterList = append(waterList, map[string]interface{}{
								"definition":        strconv.FormatUint(*waterV.Definition, 10),
								"text_content":      waterV.TextContent,
								"svg_content":       waterV.SvgContent,
								"start_time_offset": waterV.StartTimeOffset,
								"end_time_offset":   waterV.EndTimeOffset,
							})
						}
						return waterList
					}(),
					"ext_time_offset_list": func() interface{} {
						if item.ExtTimeOffsetSet == nil {
							return nil
						}
						extList := make([]interface{}, 0, len(item.ExtTimeOffsetSet))
						for _, extV := range item.ExtTimeOffsetSet {
							extList = append(extList, extV)
						}
						return extList
					}(),
					"time_offset_list": func() interface{} {
						if item.TimeOffsetSet == nil {
							return nil
						}
						timeOffsetList := make([]interface{}, 0, len(item.TimeOffsetSet))
						for _, timeOffset := range item.TimeOffsetSet {
							timeOffsetList = append(timeOffsetList, timeOffset)
						}
						return timeOffsetList
					}(),
				})
			}
			mediaProcessTaskElem["snapshot_by_time_offset_task_list"] = list
		}
		// sample_snapshot_task_list
		if template.MediaProcessTask.SampleSnapshotTaskSet != nil {
			list := make([]map[string]interface{}, 0, len(template.MediaProcessTask.SampleSnapshotTaskSet))
			for _, item := range template.MediaProcessTask.SampleSnapshotTaskSet {
				list = append(list, map[string]interface{}{
					"definition": strconv.FormatUint(*item.Definition, 10),
					"watermark_list": func() interface{} {
						if item.WatermarkSet == nil {
							return nil
						}
						waterList := make([]map[string]interface{}, 0, len(item.WatermarkSet))
						for _, waterV := range item.WatermarkSet {
							waterList = append(waterList, map[string]interface{}{
								"definition":        strconv.FormatUint(*waterV.Definition, 10),
								"text_content":      waterV.TextContent,
								"svg_content":       waterV.SvgContent,
								"start_time_offset": waterV.StartTimeOffset,
								"end_time_offset":   waterV.EndTimeOffset,
							})
						}
						return waterList
					}(),
				})
			}
			mediaProcessTaskElem["sample_snapshot_task_list"] = list
		}
		// image_sprite_task_list
		if template.MediaProcessTask.ImageSpriteTaskSet != nil {
			list := make([]map[string]interface{}, 0, len(template.MediaProcessTask.ImageSpriteTaskSet))
			for _, item := range template.MediaProcessTask.ImageSpriteTaskSet {
				list = append(list, map[string]interface{}{
					"definition": strconv.FormatUint(*item.Definition, 10),
				})
			}
			mediaProcessTaskElem["image_sprite_task_list"] = list
		}
		// cover_by_snapshot_task_list
		if template.MediaProcessTask.CoverBySnapshotTaskSet != nil {
			list := make([]map[string]interface{}, 0, len(template.MediaProcessTask.CoverBySnapshotTaskSet))
			for _, item := range template.MediaProcessTask.CoverBySnapshotTaskSet {
				list = append(list, map[string]interface{}{
					"definition": strconv.FormatUint(*item.Definition, 10),
					"watermark_list": func() interface{} {
						if item.WatermarkSet == nil {
							return nil
						}
						waterList := make([]map[string]interface{}, 0, len(item.WatermarkSet))
						for _, waterV := range item.WatermarkSet {
							waterList = append(waterList, map[string]interface{}{
								"definition":        strconv.FormatUint(*waterV.Definition, 10),
								"text_content":      waterV.TextContent,
								"svg_content":       waterV.SvgContent,
								"start_time_offset": waterV.StartTimeOffset,
								"end_time_offset":   waterV.EndTimeOffset,
							})
						}
						return waterList
					}(),
					"position_type":  item.PositionType,
					"position_value": item.PositionValue,
				})
			}
			mediaProcessTaskElem["cover_by_snapshot_task_list"] = list
		}
		// adaptive_dynamic_streaming_task_list
		if template.MediaProcessTask.AdaptiveDynamicStreamingTaskSet != nil {
			list := make([]map[string]interface{}, 0, len(template.MediaProcessTask.AdaptiveDynamicStreamingTaskSet))
			for _, item := range template.MediaProcessTask.AdaptiveDynamicStreamingTaskSet {
				list = append(list, map[string]interface{}{
					"definition": strconv.FormatUint(*item.Definition, 10),
					"watermark_list": func() interface{} {
						if item.WatermarkSet == nil {
							return nil
						}
						waterList := make([]map[string]interface{}, 0, len(item.WatermarkSet))
						for _, waterV := range item.WatermarkSet {
							waterList = append(waterList, map[string]interface{}{
								"definition":        strconv.FormatUint(*waterV.Definition, 10),
								"text_content":      waterV.TextContent,
								"svg_content":       waterV.SvgContent,
								"start_time_offset": waterV.StartTimeOffset,
								"end_time_offset":   waterV.EndTimeOffset,
							})
						}
						return waterList
					}(),
					"subtitle_list": func() interface{} {
						if item.SubtitleSet == nil {
							return nil
						}
						subtitleList := make([]interface{}, 0, len(item.SubtitleSet))
						for _, subtitleV := range item.SubtitleSet {
							subtitleList = append(subtitleList, subtitleV)
						}
						return subtitleList
					}(),
				})
			}
			mediaProcessTaskElem["adaptive_dynamic_streaming_task_list"] = list
		}

		_ = d.Set("media_process_task", []interface{}{mediaProcessTaskElem})
	}
	aiAnalysisTask := make(map[string]interface{})
	if template.AiAnalysisTask != nil && template.AiAnalysisTask.Definition != nil {
		aiAnalysisTask["definition"] = template.AiAnalysisTask.Definition
	}
	_ = d.Set("ai_analysis_task", []interface{}{aiAnalysisTask})
	aiRecognitionTask := make(map[string]interface{})
	if template.AiRecognitionTask != nil && template.AiRecognitionTask.Definition != nil {
		aiRecognitionTask["definition"] = template.AiRecognitionTask.Definition
	}
	_ = d.Set("ai_recognition_task", []interface{}{aiRecognitionTask})
	reviewAudioVideoTask := make(map[string]interface{})
	if template.ReviewAudioVideoTask != nil {
		if template.ReviewAudioVideoTask.Definition != nil {
			reviewAudioVideoTask["definition"] = template.ReviewAudioVideoTask.Definition
		}
		if template.ReviewAudioVideoTask.ReviewContents != nil {
			reviewContentList := make([]string, 0, len(template.ReviewAudioVideoTask.ReviewContents))
			for _, revireviewContent := range template.ReviewAudioVideoTask.ReviewContents {
				reviewContentList = append(reviewContentList, *revireviewContent)
			}
			reviewAudioVideoTask["review_contents"] = reviewContentList
		}
	}
	_ = d.Set("review_audio_video_task", []interface{}{reviewAudioVideoTask})
	return nil
}

func resourceTencentCloudVodProcedureTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_procedure_template.update")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = vod.NewResetProcedureTemplateRequest()
		id         = d.Id()
		changeFlag = false
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) == 2 {
		request.Name = helper.String(idSplit[0])
		subAppId := helper.StrToInt(idSplit[1])
		request.SubAppId = helper.IntUint64(subAppId)
	} else {
		request.Name = &id
		if v, ok := d.GetOk("sub_app_id"); ok {
			request.SubAppId = helper.IntUint64(v.(int))
		}
	}

	immutableArgs := []string{"sub_app_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("comment") {
		changeFlag = true
		request.Comment = helper.String(d.Get("comment").(string))
	}

	if d.HasChange("media_process_task") {
		changeFlag = true
		mediaReq := generateMediaProcessTask(d)
		request.MediaProcessTask = mediaReq
	}

	if d.HasChange("ai_analysis_task") {
		changeFlag = true
		if aiAnalysisTask, ok := d.GetOk("ai_analysis_task"); ok {
			aiAnalysisTaskList := aiAnalysisTask.([]interface{})
			aiAnalysisTaskItem := aiAnalysisTaskList[0].(map[string]interface{})

			request.AiAnalysisTask = &vod.AiAnalysisTaskInput{
				Definition: helper.StrToUint64Point(aiAnalysisTaskItem["definition"].(string)),
			}
		}
	}

	if d.HasChange("ai_recognition_task") {
		changeFlag = true
		if aiRecognitionTask, ok := d.GetOk("ai_recognition_task"); ok {
			aiRecognitionTaskList := aiRecognitionTask.([]interface{})
			aiRecognitionTaskItem := aiRecognitionTaskList[0].(map[string]interface{})

			request.AiRecognitionTask = &vod.AiRecognitionTaskInput{
				Definition: helper.StrToUint64Point(aiRecognitionTaskItem["definition"].(string)),
			}
		}
	}
	if d.HasChange("review_audio_video_task") {
		changeFlag = true
		if reviewAudioVideoTask, ok := d.GetOk("review_audio_video_task"); ok {
			reviewAudioVideoTaskList := reviewAudioVideoTask.([]interface{})
			reviewAudioVideoTaskItem := reviewAudioVideoTaskList[0].(map[string]interface{})
			request.ReviewAudioVideoTask = &vod.ProcedureReviewAudioVideoTaskInput{
				Definition: helper.StrToUint64Point(reviewAudioVideoTaskItem["definition"].(string)),
				ReviewContents: func() (list []*string) {
					if _, ok := reviewAudioVideoTaskItem["review_contents"]; !ok {
						return
					}
					reviewContentList := reviewAudioVideoTaskItem["review_contents"].([]interface{})
					list = make([]*string, 0, len(reviewContentList))
					for _, reviewContent := range reviewContentList {
						list = append(list, helper.String(reviewContent.(string)))
					}
					return
				}(),
			}
		}
	}

	if changeFlag {
		var err error
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			_, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().ResetProcedureTemplate(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
				return tccommon.RetryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}

		return resourceTencentCloudVodProcedureTemplateRead(d, meta)
	}

	return nil
}

func resourceTencentCloudVodProcedureTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_procedure_template.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	idSplit := strings.Split(id, tccommon.FILED_SP)
	var (
		name     string
		subAppId int
	)
	if len(idSplit) == 2 {
		name = idSplit[0]
		subAppId = helper.StrToInt(idSplit[1])
	} else {
		name = id
		if v, ok := d.GetOk("sub_app_id"); ok {
			subAppId = v.(int)
		}
	}
	vodService := VodService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	if err := vodService.DeleteProcedureTemplate(ctx, name, uint64(subAppId)); err != nil {
		return err
	}

	return nil
}
