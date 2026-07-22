package vod

import (
	"context"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVodProcedureTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVodProcedureTemplatesRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 procedure template。",
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
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "模板 类型 过滤器. 有效值：`Preset`，`Custom`. `Preset`: preset template; `Custom`: 自定义 template。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "任务 flow 名称",
						},
						"comment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "模板描述",
						},
						"media_process_task": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Parameter 的 视频 processing 任务。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"transcode_task_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "列表 transcoding tasks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Video transcoding 模板 ID",
												},
												"watermark_list": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "列表 up 到 `10` 镜像 或 text watermarks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
													Elem:        VodWatermarkResource(),
												},
												"mosaic_list": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "列表 blurs. Up 到 10 ones 可以 是 支持。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"coordinate_origin": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Origin position，其中 currently 可以 仅 是: `TopLeft`: 源站 的 coordinates 是 在 top-left corner 的 视频，和 源站 的 blur 是 在 top-left corner 的 镜像 或 text。",
															},
															"x_pos": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "horizontal position 的 源站 的 blur relative 到 源站 的 coordinates 的 视频. `%` 和 `像素` formats 是 支持: 如果 字符串 结束 在 `%`， XPos 的 blur 将 是 指定 percentage 的 视频 宽度; 对于 示例，10% 表示 该 XPos 是 10% 的 视频 宽度; 如果 字符串 结束 在 `像素`， XPos 的 blur 将 是 指定 像素; 对于 示例，100px 表示 该 XPos 是 100 像素。",
															},
															"y_pos": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Vertical position 的 源站 的 blur relative 到 源站 的 coordinates 的 视频. `%` 和 `像素` formats 是 支持: 如果 字符串 结束 在 `%`， YPos 的 blur 将 是 指定 percentage 的 视频 高度; 对于 示例，10% 表示 该 YPos 是 10% 的 视频 高度; 如果 字符串 结束 在 `像素`， YPos 的 blur 将 是 指定 像素; 对于 示例，100px 表示 该 YPos 是 100 像素。",
															},
															"width": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Blur 宽度. `%` 和 `像素` formats 是 支持: 如果 字符串 结束 在 `%`， `宽度` 的 blur 将 是 指定 percentage 的 视频 宽度; 对于 示例，10% 表示 该 `宽度` 是 10% 的 视频 宽度; 如果 字符串 结束 在 `像素`， `宽度` 的 blur 将 是 在 像素; 对于 示例，100px 表示 该 宽度 是 100 像素。",
															},
															"height": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Blur 高度. `%` 和 `像素` formats 是 支持: 如果 字符串 结束 在 `%`， `高度` 的 blur 将 是 指定 percentage 的 视频 高度; 对于 示例，10% 表示 该 高度 是 10% 的 视频 高度; 如果 字符串 结束 在 `像素`， `高度` 的 blur 将 是 在 像素; 对于 示例，100px 表示 该 高度 是 100 像素。",
															},
															"start_time_offset": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "开始时间 偏移量 的 blur （秒）。 如果此参数为空 或 `0` 是 entered， blur 将 appear upon first 视频 frame. 如果此参数为空 或 `0` 是 entered， blur 将 appear upon first 视频 frame; 如果 此 值 是 greater 比 `0` (e.g.，n)， blur 将 appear 在 second n after first 视频 frame; 如果 此 值 是 smaller 比 `0` (e.g.，-n)， blur 将 appear 在 second n before last 视频 frame。",
															},
															"end_time_offset": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "结束时间 偏移量 的 blur （秒）。 如果此参数为空 或 `0` 是 entered， blur 将 exist till last 视频 frame; 如果 此 值 是 greater 比 `0` (e.g.，n)， blur 将 exist till second n; 如果 此 值 是 smaller 比 `0` (e.g.，-n)， blur 将 exist till second n before last 视频 frame。",
															},
														},
													},
												},
											},
										},
									},
									"animated_graphic_task_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "列表 animated 镜像 generating tasks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Animated 镜像 generating 模板 ID",
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "开始时间 的 animated 镜像 在 视频 （秒）。",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "结束时间 的 animated 镜像 在 视频 （秒）。",
												},
											},
										},
									},
									"snapshot_by_time_offset_task_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "列表 时间 point screen capturing tasks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Time point screen capturing 模板 ID",
												},
												"ext_time_offset_list": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "列表 screenshot 时间 points. `s` 和 `%` formats 是 支持: 当 时间 point 字符串 结束 使用 `s`，its 单位 是 second. For 示例，`3.5s` 表示 3.5th second 的 视频; 当 时间 point 字符串 结束 使用 `%`，它 是 marked 使用 corresponding percentage 的 视频 时长. For 示例，`10%` 表示 该 时间 point 是 在 10% 的 视频 entire 时长。",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"watermark_list": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "列表 up 到 `10` 镜像 或 text watermarks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
													Elem:        VodWatermarkResource(),
												},
											},
										},
									},
									"sample_snapshot_task_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "列表 sampled screen capturing tasks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Sampled screen capturing 模板 ID",
												},
												"watermark_list": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "列表 up 到 `10` 镜像 或 text watermarks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
													Elem:        VodWatermarkResource(),
												},
											},
										},
									},
									"image_sprite_task_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "列表 镜像 sprite generating tasks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Image sprite generating 模板 ID",
												},
											},
										},
									},
									"cover_by_snapshot_task_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "列表 cover generating tasks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Time point screen capturing 模板 ID",
												},
												"position_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Screen capturing 模式 有效值：`Time`，`Percent`. `Time`: screen captures 通过 时间 point，`Percent`: screen captures 通过 percentage。",
												},
												"position_value": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Screenshot position: For 时间 point screen capturing，此 表示 到 take screenshot 在 指定 时间 point (在 秒) 和 使用 它 作为 cover. For percentage screen capturing，此 值 表示 到 take screenshot 在 指定 percentage 的 视频 时长 和 使用 它 作为 cover。",
												},
												"watermark_list": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "列表 up 到 `10` 镜像 或 text watermarks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
													Elem:        VodWatermarkResource(),
												},
											},
										},
									},
									"adaptive_dynamic_streaming_task_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "列表 adaptive bitrate streaming tasks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Adaptive bitrate streaming 模板 ID",
												},
												"watermark_list": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "列表 up 到 `10` 镜像 或 text watermarks. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
													Elem:        VodWatermarkResource(),
												},
											},
										},
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

func dataSourceTencentCloudVodProcedureTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vod_procedure_templates.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	filter := make(map[string]interface{})
	if v, ok := d.GetOk("name"); ok {
		filter["name"] = []string{v.(string)}
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
	templates, err := vodService.DescribeProcedureTemplatesByFilter(ctx, filter)
	if err != nil {
		return err
	}

	templatesList := make([]map[string]interface{}, 0, len(templates))
	ids := make([]string, 0, len(templates))
	for _, templateItem := range templates {
		templatesList = append(templatesList, func() map[string]interface{} {
			mapping := map[string]interface{}{
				"type":        templateItem.Type,
				"name":        templateItem.Name,
				"comment":     templateItem.Comment,
				"create_time": templateItem.CreateTime,
				"update_time": templateItem.UpdateTime,
			}
			mediaProcessTaskElem := make(map[string]interface{})
			if templateItem.MediaProcessTask != nil {
				// transcode_task_list
				if templateItem.MediaProcessTask.TranscodeTaskSet != nil {
					list := make([]map[string]interface{}, 0, len(templateItem.MediaProcessTask.TranscodeTaskSet))
					for _, item := range templateItem.MediaProcessTask.TranscodeTaskSet {
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
						})
					}
					mediaProcessTaskElem["transcode_task_list"] = list
				}
				// animated_graphic_task_list
				if templateItem.MediaProcessTask.AnimatedGraphicTaskSet != nil {
					list := make([]map[string]interface{}, 0, len(templateItem.MediaProcessTask.AnimatedGraphicTaskSet))
					for _, item := range templateItem.MediaProcessTask.AnimatedGraphicTaskSet {
						list = append(list, map[string]interface{}{
							"definition":        strconv.FormatUint(*item.Definition, 10),
							"start_time_offset": item.StartTimeOffset,
							"end_time_offset":   item.EndTimeOffset,
						})
					}
					mediaProcessTaskElem["animated_graphic_task_list"] = list
				}
				// snapshot_by_time_offset_task_list
				if templateItem.MediaProcessTask.SnapshotByTimeOffsetTaskSet != nil {
					list := make([]map[string]interface{}, 0, len(templateItem.MediaProcessTask.SnapshotByTimeOffsetTaskSet))
					for _, item := range templateItem.MediaProcessTask.SnapshotByTimeOffsetTaskSet {
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
						})
					}
					mediaProcessTaskElem["snapshot_by_time_offset_task_list"] = list
				}
				// sample_snapshot_task_list
				if templateItem.MediaProcessTask.SampleSnapshotTaskSet != nil {
					list := make([]map[string]interface{}, 0, len(templateItem.MediaProcessTask.SampleSnapshotTaskSet))
					for _, item := range templateItem.MediaProcessTask.SampleSnapshotTaskSet {
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
				if templateItem.MediaProcessTask.ImageSpriteTaskSet != nil {
					list := make([]map[string]interface{}, 0, len(templateItem.MediaProcessTask.ImageSpriteTaskSet))
					for _, item := range templateItem.MediaProcessTask.ImageSpriteTaskSet {
						list = append(list, map[string]interface{}{
							"definition": strconv.FormatUint(*item.Definition, 10),
						})
					}
					mediaProcessTaskElem["image_sprite_task_list"] = list
				}
				// cover_by_snapshot_task_list
				if templateItem.MediaProcessTask.CoverBySnapshotTaskSet != nil {
					list := make([]map[string]interface{}, 0, len(templateItem.MediaProcessTask.CoverBySnapshotTaskSet))
					for _, item := range templateItem.MediaProcessTask.CoverBySnapshotTaskSet {
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
				if templateItem.MediaProcessTask.AdaptiveDynamicStreamingTaskSet != nil {
					list := make([]map[string]interface{}, 0, len(templateItem.MediaProcessTask.AdaptiveDynamicStreamingTaskSet))
					for _, item := range templateItem.MediaProcessTask.AdaptiveDynamicStreamingTaskSet {
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
					mediaProcessTaskElem["adaptive_dynamic_streaming_task_list"] = list
				}
			}
			mapping["media_process_task"] = []interface{}{mediaProcessTaskElem}
			ids = append(ids, *templateItem.Name)
			return mapping
		}())
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("template_list", templatesList); e != nil {
		log.Printf("[CRITAL]%s provider set procedure template list fail, reason:%s ", logId, e.Error())
	}

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), templatesList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]", logId, output.(string), err.Error())
			return err
		}
	}

	return nil
}
