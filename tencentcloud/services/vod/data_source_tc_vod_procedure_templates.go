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
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template 类型 filter. 有效值：`Preset`，`Custom`. `Preset`: preset template; `Custom`: custom template。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task flow 名称",
						},
						"comment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "模板描述",
						},
						"media_process_task": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Parameter of video processing task。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"transcode_task_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "列表 transcoding tasks. Note: this field may return null，indicating that no valid values can be obtained。",
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
													Description: "列表 up to `10` image or text watermarks. Note: this field may return null，indicating that no valid values can be obtained。",
													Elem:        VodWatermarkResource(),
												},
												"mosaic_list": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "列表 blurs. Up to 10 ones can be supported。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"coordinate_origin": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Origin position，which currently can only be: `TopLeft`: the origin of coordinates is in the top-left corner of the video，and the origin of the blur is in the top-left corner of the image or text。",
															},
															"x_pos": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The horizontal position of the origin of the blur relative to the origin of coordinates of the video. `%` and `px` formats are supported: If the string ends in `%`，the XPos of the blur will be the specified percentage of the video width; for example，10% means that XPos is 10% of the video width; If the string ends in `px`，the XPos of the blur will be the specified px; for example，100px means that XPos is 100 px。",
															},
															"y_pos": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Vertical position of the origin of blur relative to the origin of coordinates of video. `%` and `px` formats are supported: If the string ends in `%`，the YPos of the blur will be the specified percentage of the video height; for example，10% means that YPos is 10% of the video height; If the string ends in `px`，the YPos of the blur will be the specified px; for example，100px means that YPos is 100 px。",
															},
															"width": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Blur width. `%` and `px` formats are supported: If the string ends in `%`，the `width` of the blur will be the specified percentage of the video width; for example，10% means that `width` is 10% of the video width; If the string ends in `px`，the `width` of the blur will be in px; for example，100px means that Width is 100 px。",
															},
															"height": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Blur height. `%` and `px` formats are supported: If the string ends in `%`，the `height` of the blur will be the specified percentage of the video height; for example，10% means that Height is 10% of the video height; If the string ends in `px`，the `height` of the blur will be in px; for example，100px means that Height is 100 px。",
															},
															"start_time_offset": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "开始时间 偏移量 of blur （秒）。 如果此参数为空 or `0` is entered，the blur will appear upon the first video frame. 如果此参数为空 or `0` is entered，the blur will appear upon the first video frame; If this 值 is greater than `0` (e.g.，n)，the blur will appear at second n after the first video frame; If this 值 is smaller than `0` (e.g.，-n)，the blur will appear at second n before the last video frame。",
															},
															"end_time_offset": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "结束时间 偏移量 of blur （秒）。 如果此参数为空 or `0` is entered，the blur will exist till the last video frame; If this 值 is greater than `0` (e.g.，n)，the blur will exist till second n; If this 值 is smaller than `0` (e.g.，-n)，the blur will exist till second n before the last video frame。",
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
										Description: "列表 animated image generating tasks. Note: this field may return null，indicating that no valid values can be obtained。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Animated image generating 模板 ID",
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "开始时间 of animated image in video （秒）。",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "结束时间 of animated image in video （秒）。",
												},
											},
										},
									},
									"snapshot_by_time_offset_task_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "列表 time point screen capturing tasks. Note: this field may return null，indicating that no valid values can be obtained。",
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
													Description: "The 列表 screenshot time points. `s` and `%` formats are supported: When a time point string ends with `s`，its unit is second. For example，`3.5s` means the 3.5th second of the video; When a time point string ends with `%`，it is marked with corresponding percentage of the video duration. For example，`10%` means that the time point is at the 10% of the video entire duration。",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"watermark_list": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "列表 up to `10` image or text watermarks. Note: this field may return null，indicating that no valid values can be obtained。",
													Elem:        VodWatermarkResource(),
												},
											},
										},
									},
									"sample_snapshot_task_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "列表 sampled screen capturing tasks. Note: this field may return null，indicating that no valid values can be obtained。",
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
													Description: "列表 up to `10` image or text watermarks. Note: this field may return null，indicating that no valid values can be obtained。",
													Elem:        VodWatermarkResource(),
												},
											},
										},
									},
									"image_sprite_task_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "列表 image sprite generating tasks. Note: this field may return null，indicating that no valid values can be obtained。",
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
										Description: "列表 cover generating tasks. Note: this field may return null，indicating that no valid values can be obtained。",
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
													Description: "Screen capturing 模式 有效值：`Time`，`Percent`. `Time`: screen captures by time point，`Percent`: screen captures by percentage。",
												},
												"position_value": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Screenshot position: For time point screen capturing，this means to take a screenshot at a specified time point (in seconds) and use it as the cover. For percentage screen capturing，this 值 means to take a screenshot at a specified percentage of the video duration and use it as the cover。",
												},
												"watermark_list": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "列表 up to `10` image or text watermarks. Note: this field may return null，indicating that no valid values can be obtained。",
													Elem:        VodWatermarkResource(),
												},
											},
										},
									},
									"adaptive_dynamic_streaming_task_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "列表 adaptive bitrate streaming tasks. Note: this field may return null，indicating that no valid values can be obtained。",
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
													Description: "列表 up to `10` image or text watermarks. Note: this field may return null，indicating that no valid values can be obtained。",
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
