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

func ResourceTencentCloudMpsWorkflow() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsWorkflowCreate,
		Read:   resourceTencentCloudMpsWorkflowRead,
		Update: resourceTencentCloudMpsWorkflowUpdate,
		Delete: resourceTencentCloudMpsWorkflowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"workflow_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Workflow 名称，up 到 128 字符. 名称 是 唯一 对于 same 用户",
			},

			"trigger": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "触发器 规则 bound 到 工作流，当 uploaded 视频 hits 规则 到 此 对象， 工作流 将 是 triggered。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "类型 触发器，currently 仅 支持 CosFileUpload。",
						},
						"cos_file_upload_trigger": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Mandatory 和 有效 当 类型 是 CosFileUpload， 规则 是 triggered 对于 COS.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "名称 COS 存储桶 bound 到 工作流。",
									},
									"region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "park 到 其中 COS 存储桶 bound 到 工作流 belongs。",
									},
									"dir": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "input 路径 directory 的 工作流 binding 必须 是 absolute 路径，该 是，start 和 end 使用 `/`。",
									},
									"formats": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Computed:    true,
										Description: "A 列表 文件 formats 该 是 allowed 到 是 triggered 通过 工作流，如果未填写 在，它 表示 该 files 的 all formats 可以 触发器 工作流。",
									},
								},
							},
						},
					},
				},
			},

			"output_storage": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "File output 存储 location 对于 media processing. 如果 left blank， 存储 location 在 Trigger 将 是 inherited。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "类型 media processing output 对象 存储 location，now 仅 支持 COS。",
						},
						"cos_output_storage": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "有效 当 类型 是 COS，此 item 为必填项，indicating media processing COS output location.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "目标 存储桶 名称 文件 output generated 通过 media processing，如果未填写，它 表示 upper layer。",
									},
									"region": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "park 的 目标 存储桶 对于 output 的 文件 generated 通过 media processing. 如果未填写，它 表示 inheriting 从 upper layer。",
									},
								},
							},
						},
					},
				},
			},

			"output_dir": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "目标 directory 的 输出文件 generated 通过 media processing，如果未填写，它 表示 该 它 是 consistent 使用 directory 其中 触发器 文件 是 located。",
			},

			"media_process_task": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Media Processing 类型 任务 Parameters。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"transcode_task_set": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Video Transcoding 任务 List。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Video Transcoding 模板 ID",
									},
									"raw_parameter": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Computed:    true,
										Description: "Video transcoding 自定义 参数，有效 当 Definition 是 filled 使用 0.此 参数 是 使用 在 highly customized scenarios. It 是 recommended 该 您 使用 Definition 到 指定transcoding 参数 first.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"container": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Encapsulation 格式，可选 值: mp4，flv，hls，mp3，flac，ogg，m4a. Among them，mp3，flac，ogg，m4a 是 pure 音频 files。",
												},
												"remove_video": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "是否remove 视频 数据，值:0: reserved.1: remove.默认值：0。",
												},
												"remove_audio": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "是否remove 音频 数据，值:0: reserved.1: remove.默认值：0。",
												},
												"video_template": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
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
																Description: "Filling 方法，当 aspect ratio 的 视频 流 配置 是 inconsistent 使用 aspect ratio 的 original 视频， processing 方法 对于 transcoding 是 filling. 可选 filling 方法:stretch: Stretch，stretch each frame 到 fill entire screen，其中 可能 cause transcoded 视频 到 是 squashed 或 stretched;.black: Leave black，keep aspect ratio 的 视频 unchanged，和 fill rest 的 edge 使用 black.white: Leave blank，keep aspect ratio 的 视频 unchanged，和 fill rest 的 edge 使用 white.gauss: Gaussian blur，keep aspect ratio 的 视频 unchanged，和 fill rest 的 edge 使用 Gaussian blur.默认值：black.注意: Adaptive 流 仅 支持 stretch，black。",
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
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
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
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
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
											},
										},
									},
									"override_parameter": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Video transcoding 自定义 参数，有效 当 Definition 是 不 filled 使用 0.当 some transcoding 参数 在 此 structure 是 filled 在， 参数 在 transcoding template 将 是 overwritten 使用 filled 参数.此 参数 是 使用 在 highly customized scenarios，它 是 recommended 该 您 仅 使用 Definition 到 指定transcoding 参数.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"container": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Encapsulation 格式，可选 值: mp4，flv，hls，mp3，flac，ogg，m4a. Among them，mp3，flac，ogg，m4a 是 pure 音频 files。",
												},
												"remove_video": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "是否remove 视频 数据，值:0: reserved.1: remove。",
												},
												"remove_audio": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "是否remove 音频 数据，值:0: reserved.1: remove。",
												},
												"video_template": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Video streaming 配置 参数。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"codec": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Encoding 格式 的 视频 流，可选 值:libx264: H.264 编码.libx265: H.265 编码.av1: AOMedia Video 1 编码.注意: Currently H.265 编码 必须 指定a resolution，和 它 needs 到 是 within 640*480.注意: av1 encoded containers currently 仅 support mp4。",
															},
															"fps": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Video frame 速率，取值范围：[0，100]，单位: Hz.当 值 是 0，它 表示 该 frame 速率 是 consistent 使用 original 视频。",
															},
															"bitrate": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Bit 速率 的 视频 流，取值范围：0 和 [128，35000]，单位: kbps.当 值 是 0，它 表示 该 视频 bit 速率 是 consistent 使用 original 视频。",
															},
															"resolution_adaptive": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Adaptive resolution，可选 值:```open: open，在 此 时间，宽度 表示 long side 的 视频，高度 表示 short side 的 视频.close: close，在 此 时间，宽度 表示 宽度 的 视频，和 高度 表示 高度 的 视频.注意: In adaptive 模式，宽度 不能 是 smaller 比 高度。",
															},
															"width": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "最大 值 的 视频 流 宽度 (或 long side)，取值范围：0 和 [128，4096]，单位: 像素.当 宽度 和 高度 是 both 0， resolution 是 same.当 宽度 是 0 和 高度 是 不 0，宽度 是 scaled proportionally.当 宽度 是 不 0 和 高度 是 0，高度 是 scaled proportionally.当 both 宽度 和 高度 是 不 0， resolution 是 指定 通过 用户",
															},
															"height": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "最大 值 的 视频 流 高度 (或 short side)，取值范围：0 和 [128，4096]，单位: 像素。",
															},
															"gop": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "间隔 between keyframe I frames，取值范围：0 和 [1，100000]，单位: 数量 frames.当 filling 0 或 不 filling， 系统 将 automatically 集合 gop 长度。",
															},
															"fill_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Filling 方法，当 aspect ratio 的 视频 流 配置 是 inconsistent 使用 aspect ratio 的 original 视频， processing 方法 对于 transcoding 是 filling;. 可选 filling 方法:stretch: Stretch，stretch each frame 到 fill entire screen，其中 可能 cause transcoded 视频 到 是 squashed 或 stretched; black: Leave black，keep aspect ratio 的 视频 unchanged，和 fill rest 的 edge 使用 black.white: Leave blank，keep aspect ratio 的 视频 unchanged，和 fill rest 的 edge 使用 white.gauss: Gaussian blur，keep aspect ratio 的 视频 unchanged，和 fill rest 的 edge 使用 Gaussian blur。",
															},
															"vcrf": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Video constant bit 速率 control factor， 值 范围 是 [1，51]，Fill 在 0 到 disable 此 参数.如果 there 是 无 special requirement，它 是 不 recommended 到 指定this 参数。",
															},
															"content_adapt_stream": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "内容 Adaptive Encoding. 可选 值:0: 不 open.1: open.默认值：0.当 此 参数 是 turned 在，多个 代码 streams 使用 different resolutions 和 different bit rates 将 是 adaptively generated. 宽度 和 高度 的 VideoTemplate 是 最大 resolutions among 多个 代码 streams，和 bit rates 在 VideoTemplate 是 多个 代码 rates. highest bit 速率 在 流， vcrf 在 VideoTemplate 是 highest quality among 多个 bit streams. 当 resolution，bit 速率 和 vcrf 是 不 集合， highest resolution generated 通过 ContentAdaptStream 参数 是 resolution 的 视频 来源，和 视频 quality 是 close 到 vmaf95. To 启用 此 参数 或 learn about billing details，please contact your Tencent Cloud Business。",
															},
														},
													},
												},
												"audio_template": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Audio 流 配置 参数。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"codec": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Encoding 格式 的 频率 流.当 outer 参数 Container 是 mp3， 可选 值 是:libmp3lame.当 outer 参数 Container 是 ogg 或 flac， 可选 值 是:flac.当 outer 参数 Container 是 m4a， 可选 值 是:libfdk_aac.libmp3lame.ac3.当 outer 参数 Container 是 mp4 或 flv， 可选 值 是:libfdk_aac: more suitable 对于 mp4.libmp3lame: more suitable 对于 flv.当 outer 参数 Container 是 hls， 可选 值 是:libfdk_aac.libmp3lame。",
															},
															"bitrate": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Bit 速率 的 视频 流，取值范围：0 和 [128，35000]，单位: kbps.当 值 是 0，它 表示 该 视频 bit 速率 是 consistent 使用 original 视频。",
															},
															"sample_rate": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Sampling 速率 的 音频 流，可选 值32000.44100.48000.单位：Hz。",
															},
															"audio_channel": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Audio channel 模式，可选 值:`1: 单个 channel.2: Dual channel.6: Stereo.当 包 格式 的 media 是 音频 格式 (flac，ogg，mp3，m4a)， 数量 channels 是 不 allowed 到 是 集合 到 stereo。",
															},
															"stream_selects": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeInt,
																},
																Optional:    true,
																Description: "指定audio track 到 preserve 对于 output. 默认为 到 keep all sources。",
															},
														},
													},
												},
												"tehd_config": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Ultra-fast HD transcoding 参数。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Extremely high-definition 类型，可选 值:TEHD-100: Extreme HD-100.Not filling 表示 该 ultra-fast high-definition 是 不 已启用",
															},
															"max_video_bitrate": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "upper 限制 的 视频 bit 速率，No filling 表示 无 modification。",
															},
														},
													},
												},
												"subtitle_template": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Subtitle Stream Configuration Parameters。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"path": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "地址 的 subtitle 文件 到 是 compressed into 视频。",
															},
															"stream_index": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "指定subtitle track 到 是 compressed into 视频. 如果 there 是 指定 路径， 路径 has higher 优先级 路径 和 StreamIndex 指定at least 一个。",
															},
															"font_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Font 类型hei.ttf，song.ttf，simkai.ttf，arial.ttf.默认值：hei.ttf。",
															},
															"font_size": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Font 大小，格式: Npx，N 是 值，如果未指定， subtitle 文件 shall prevail。",
															},
															"font_color": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Font color，格式: 0xRRGGBB，默认值：0xFFFFFF (white)。",
															},
															"font_alpha": {
																Type:        schema.TypeFloat,
																Optional:    true,
																Description: "Text transparency，取值范围：(0，1].0: fully transparent.1: fully opaque.默认值：1。",
															},
														},
													},
												},
											},
										},
									},
									"watermark_set": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Watermark 列表，support 多个 pictures 或 text watermarks，up 到 10.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Watermark 模板 ID",
												},
												"raw_parameter": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Computed:    true,
													Description: "Watermark 自定义 参数，有效 当 Definition 是 filled 使用 0.此 参数 是 使用 在 highly customized scenarios，它 是 recommended 该 您 使用 Definition 到 指定watermark 参数 first.Watermark 自定义 参数 do 不 support screenshot watermarking。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Watermark 类型，可选 值:镜像: 镜像 水印。",
															},
															"coordinate_origin": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Origin position，currently 仅 支持:TopLeft: 表示that 源站 的 coordinates 是 在 upper left corner 的 视频 镜像，和 源站 的 水印 是 upper left corner 的 picture 或 text.默认值：TopLeft。",
															},
															"x_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "horizontal position 的 源站 的 水印 从 源站 的 coordinates 的 视频 镜像. Support %，像素 two formats:当 字符串 结束 使用 %，它 表示 该 水印 XPos 指定a percentage 对于 视频 宽度，such 作为 10% 表示 该 XPos 是 10% 的 视频 宽度.当 字符串 结束 使用 像素，它 表示 该 水印 XPos 是 指定 pixel，such 作为 100px 表示 该 XPos 是 100 pixels.默认值：0px。",
															},
															"y_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "vertical position 的 源站 的 水印 从 源站 的 coordinates 的 视频 镜像. Support %，像素 two formats:当 字符串 结束 使用 %，它 表示 该 水印 YPos 指定a percentage 对于 视频 高度，such 作为 10% 表示 该 YPos 是 10% 的 视频 高度.当 字符串 结束 使用 像素，它 表示 该 水印 YPos 是 指定 pixel，such 作为 100px 表示 该 YPos 是 100 pixels.默认值：0px。",
															},
															"image_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Image 水印 template，当 类型 是 镜像，此 字段 为必填项. 当 类型 是 text，此 字段 是 无效。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"image_content": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Required:    true,
																			Description: "input 内容 的 水印 镜像. Support jpeg，png 镜像 格式",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Enter 类型 来源 对象，其中 支持 COS 和 URL",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "有效 当 类型 是 COS，此 item 为必填项，indicating media processing COS 对象 信息。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "名称 COS 存储桶 其中 media processing 对象 文件 是 located。",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "park 到 其中 COS 存储桶 其中 media processing 目标 文件 resides belongs。",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "Input 路径 对于 media processing 对象 files。",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "有效 当 类型 是 URL，此 item 为必填项，indicating media processing URL 对象 信息.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "Video URL",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"width": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "宽度 的 水印. Support %，像素 two formats:当 字符串 结束 使用 %，它 表示 该 水印 宽度 是 percentage 的 视频 宽度，such 作为 10% 表示 该 宽度 是 10% 的 视频 宽度.当 字符串 结束 使用 像素，它 表示 该 水印 宽度 单位 是 pixels，such 作为 100px 表示 该 宽度 是 100 pixels.默认值：10%。",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "高度 的 水印. Support %，像素 two formats:当 字符串 结束 使用 %，它 表示 该 水印 高度 是 percentage 大小 的 视频 高度，such 作为 10% 表示 该 高度 是 10% 的 视频 高度.当 字符串 结束 使用 像素，它 表示 该 水印 高度 单位 是 pixel，such 作为 100px 表示 该 高度 是 100 pixels.默认值：0px，indicating 该 高度 是 scaled according 到 aspect ratio 的 original 水印 镜像。",
																		},
																		"repeat_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Watermark repeat 类型 Usage scenario: 水印 是 动态 镜像. Ranges.once: After 动态 水印 是 played，它 将 无 longer appear.repeat_last_frame: After 水印 是 played，stay 在 last frame.repeat: 水印 loops until end 的 视频 (默认值)。",
																		},
																	},
																},
															},
														},
													},
												},
												"text_content": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Text 内容， 长度 does 不 exceed 100 字符. Fill 在 仅 当 水印 类型 是 text 水印.Text 水印 does 不 support screenshot watermarking。",
												},
												"svg_content": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "SVG 内容 长度 不能 exceed 2000000 字符. Fill 在 仅 如果 水印 类型 是 SVG 水印.SVG 水印 does 不 support screenshot watermarking。",
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "开始时间 偏移量 的 水印，单位: second. Do 不 fill 在 或 fill 在 0，其中 表示 该 水印 将 start 到 appear 当 screen appears.Do 不 fill 在 或 fill 在 0，其中 表示 水印 将 appear 从 beginning 的 screen.当 值 是 greater 比 0 (assumed 到 是 n)，它 表示 该 水印 appears 从 nth second 的 screen.当 值 是 less 比 0 (assumed 到 是 -n)，它 表示 该 水印 starts 到 appear n 秒 before end 的 screen。",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "结束时间 偏移量 的 水印，单位: second.Do 不 fill 在 或 fill 在 0，indicating 该 水印 lasts until end 的 screen.当 值 是 greater 比 0 (assumed 到 是 n)，它 表示 该 水印 lasts until nth second 和 disappears.当 值 是 less 比 0 (assumed 到 是 -n)，它 表示 该 水印 lasts until 它 disappears n 秒 before end 的 screen。",
												},
											},
										},
									},
									"mosaic_set": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Mosaic 列表，up 到 10 sheets 可以 是 支持。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"coordinate_origin": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Origin position，currently 仅 支持:TopLeft: 表示that coordinate 源站 是 located 在 upper left corner 的 视频 镜像，和 源站 的 mosaic 是 upper left corner 的 picture 或 text默认值：TopLeft。",
												},
												"x_pos": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "horizontal position 的 源站 的 水印 从 源站 的 coordinates 的 视频 镜像. Support %，像素 two formats:当 字符串 结束 使用 %，它 表示 该 水印 XPos 指定a percentage 对于 视频 宽度，such 作为 10% 表示 该 XPos 是 10% 的 视频 宽度.当 字符串 结束 使用 像素，它 表示 该 水印 XPos 是 指定 pixel，such 作为 100px 表示 该 XPos 是 100 pixels.默认值：0px。",
												},
												"y_pos": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "vertical position 的 源站 的 水印 从 源站 的 coordinates 的 视频 镜像. Support %，像素 two formats:当 字符串 结束 使用 %，它 表示 该 水印 YPos 指定a percentage 对于 视频 高度，such 作为 10% 表示 该 YPos 是 10% 的 视频 高度.当 字符串 结束 使用 像素，它 表示 该 水印 YPos 是 指定 pixel，such 作为 100px 表示 该 YPos 是 100 pixels.默认值：0px。",
												},
												"width": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "宽度 的 mosaic. Support %，像素 two formats:当 字符串 结束 使用 %，它 表示 该 mosaic 宽度 是 percentage 大小 的 视频 宽度，such 作为 10% 表示 该 宽度 是 10% 的 视频 宽度. 字符串 结束 使用 像素，indicating 该 mosaic 宽度 单位 是 pixels，such 作为 100px 表示that 宽度 是 100 pixels.默认值：10%。",
												},
												"height": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "高度 的 mosaic. Support %，像素 two formats.当 字符串 结束 使用 %，它 表示 该 mosaic 高度 是 percentage 大小 的 视频 高度，such 作为 10% 表示 该 高度 是 10% 的 视频 高度.当 字符串 结束 使用 像素，它 表示 该 mosaic 高度 单位 是 pixel，such 作为 100px 表示 该 高度 是 100 pixels.默认值：10%。",
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "开始时间 偏移量 的 mosaic，单位: second. Do 不 fill 或 fill 在 0，其中 表示 该 mosaic 将 start 到 appear 当 screen appears.Fill 在 或 fill 在 0，其中 表示 该 mosaic 将 appear 从 beginning 的 screen.当 值 是 greater 比 0 (assumed 到 是 n)，它 表示 该 mosaic appears 从 nth second 的 screen.当 值 是 less 比 0 (assumed 到 是 -n)，它 表示 该 mosaic starts 到 appear n 秒 before end 的 screen。",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "结束时间 偏移量 的 mosaic，单位: second.Fill 在 或 fill 在 0，indicating 该 mosaic continues until end 的 screen.当 值 是 greater 比 0 (assumed 到 是 n)，它 表示 该 mosaic lasts until nth second 和 disappears.当 值 是 less 比 0 (assumed 到 是 -n)，它 表示 该 mosaic lasts until 它 disappears n 秒 before end 的 screen。",
												},
											},
										},
									},
									"start_time_offset": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "开始时间 偏移量 的 transcoded 视频，单位: second.Do 不 fill 在 或 fill 在 0，indicating 该 transcoded 视频 starts 从 beginning 的 original 视频.当 值 是 greater 比 0 (assumed 到 是 n)，它 表示 该 transcoded 视频 starts 从 nth second position 的 original 视频.当 值 是 less 比 0 (assumed 到 是 -n)，它 表示 该 transcoded 视频 starts 从 position n 秒 before end 的 original 视频。",
									},
									"end_time_offset": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "结束时间 偏移量 的 视频 after transcoding，单位: second.Do 不 fill 在 或 fill 在 0，indicating 该 transcoded 视频 continues until end 的 original 视频.当 值 是 greater 比 0 (assumed 到 是 n)，它 表示 该 transcoded 视频 lasts until nth second 的 original 视频 和 terminates.当 值 是 less 比 0 (assumed 到 是 -n)，它 表示 该 transcoded 视频 lasts until n 秒 before end 的 original 视频。",
									},
									"output_storage": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "目标 存储 的 transcoded 文件，如果未填写，它 将 inherit OutputStorage 值 的 upper layer.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "类型 media processing output 对象 存储 location，now 仅 支持 COS。",
												},
												"cos_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "有效 当 类型 是 COS，此 item 为必填项，indicating media processing COS output location.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bucket": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "目标 存储桶 名称 文件 output generated 通过 media processing，如果未填写，它 表示 upper layer。",
															},
															"region": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "park 的 目标 存储桶 对于 output 的 文件 generated 通过 media processing. 如果未填写，它 表示 inheriting 从 upper layer。",
															},
														},
													},
												},
											},
										},
									},
									"output_object_path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "输出路径 的 main 文件 after transcoding 可以 是 relative 路径 或 absolute 路径 如果未填写， 默认为 relative 路径: {inputName}_transcode_{definition}.{格式}。",
									},
									"segment_object_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "输出路径 的 transcoded fragment 文件 ( 路径 的 ts 当 transcoding HLS)，可以 仅 是 relative 路径 如果未填写， 默认值 是: `{inputName}_transcode_{definition}_{数量}.{格式}。",
									},
									"object_number_format": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Rules 对于 `{数量}` variable 在 输出路径 after transcoding.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"initial_value": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "starting 值 的 `{数量}` variable， 默认为 0。",
												},
												"increment": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "growth step 的 `{数量}` variable， 默认为 1。",
												},
												"min_length": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "最小长度the `{数量}` variable，如果 insufficient，placeholders 将 是 filled. 默认为 1。",
												},
												"place_holder": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "当 长度 的 `{数量}` variable 是 insufficient， placeholder 是 added. 默认为 0。",
												},
											},
										},
									},
									"head_tail_parameter": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Opening 和 ending 参数.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"head_set": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Title 列表。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Enter 类型 来源 对象，其中 支持 COS 和 URL",
															},
															"cos_input_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "有效 当 类型 是 COS，此 item 为必填项，indicating media processing COS 对象 信息。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bucket": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "名称 COS 存储桶 其中 media processing 对象 文件 是 located。",
																		},
																		"region": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "park 到 其中 COS 存储桶 其中 media processing 目标 文件 resides belongs。",
																		},
																		"object": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Input 路径 对于 media processing 对象 files。",
																		},
																	},
																},
															},
															"url_input_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "有效 当 类型 是 URL，此 item 为必填项，indicating media processing URL 对象 信息.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Video URL",
																		},
																	},
																},
															},
														},
													},
												},
												"tail_set": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Ending List。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Enter 类型 来源 对象，其中 支持 COS 和 URL",
															},
															"cos_input_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "有效 当 类型 是 COS，此 item 为必填项，indicating media processing COS 对象 信息。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bucket": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "名称 COS 存储桶 其中 media processing 对象 文件 是 located。",
																		},
																		"region": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "park 到 其中 COS 存储桶 其中 media processing 目标 文件 resides belongs。",
																		},
																		"object": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Input 路径 对于 media processing 对象 files。",
																		},
																	},
																},
															},
															"url_input_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "有效 当 类型 是 URL，此 item 为必填项，indicating media processing URL 对象 信息.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Video URL",
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
							},
						},
						"animated_graphic_task_set": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Video Rotation Map 任务 List。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Video turntable template ID。",
									},
									"start_time_offset": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: "开始时间 的 animation 在 视频，（秒）。",
									},
									"end_time_offset": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: "结束时间 的 animation 在 视频，（秒）。",
									},
									"output_storage": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "目标 存储 的 transcoded 文件，如果未填写，它 将 inherit OutputStorage 值 的 upper layer.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "类型 media processing output 对象 存储 location，now 仅 支持 COS。",
												},
												"cos_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "有效 当 类型 是 COS，此 item 为必填项，indicating media processing COS output location.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bucket": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "目标 存储桶 名称 文件 output generated 通过 media processing，如果未填写，它 表示 upper layer。",
															},
															"region": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "park 的 目标 存储桶 对于 output 的 文件 generated 通过 media processing. 如果未填写，它 表示 inheriting 从 upper layer。",
															},
														},
													},
												},
											},
										},
									},
									"output_object_path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "输出路径 的 文件 after rotating 镜像，其中 可以 是 relative 路径 或 absolute 路径 如果未填写， 默认为 relative 路径: {inputName}_animatedGraphic_{definition}.{格式}。",
									},
								},
							},
						},
						"snapshot_by_time_offset_task_set": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Screenshot 任务 列表 视频 according 到 时间 point。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Specified 时间 point screenshot 模板 ID",
									},
									"ext_time_offset_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Screenshot 时间 point 列表， 时间 point 支持 two formats: s 和 %:;当 字符串 结束 使用 s，它 表示 该 时间 point 是 （秒）， such 作为 3.5s 表示 该 时间 point 是 3.5th second.当 字符串 结束 使用 %，它 表示 该 时间 point 是 percentage 的 视频 时长，such 作为 10% 表示 该 时间 point 是 first 10% 的 时间 在 视频。",
									},
									"time_offset_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeFloat,
										},
										Optional:    true,
										Description: "Screenshot 时间 point 列表， 单位 是 &lt;font color=red&gt;秒&lt;/font&gt;. 此 参数 是 无 longer recommended，它 是 recommended 该 您 使用 ExtTimeOffsetSet 参数。",
									},
									"watermark_set": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Watermark 列表，support 多个 pictures 或 text watermarks，up 到 10。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Watermark 模板 ID",
												},
												"raw_parameter": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Computed:    true,
													Description: "Watermark 自定义 参数，有效 当 Definition 是 filled 使用 0.此 参数 是 使用 在 highly customized scenarios，它 是 recommended 该 您 使用 Definition 到 指定watermark 参数 first.Watermark 自定义 参数 do 不 support screenshot watermarking。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Watermark 类型，可选 值:镜像: 镜像 水印。",
															},
															"coordinate_origin": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Origin position，currently 仅 支持:TopLeft: 表示that 源站 的 coordinates 是 在 upper left corner 的 视频 镜像，和 源站 的 水印 是 upper left corner 的 picture 或 text.默认值：TopLeft。",
															},
															"x_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "horizontal position 的 源站 的 水印 从 源站 的 coordinates 的 视频 镜像. Support %，像素 two formats:当 字符串 结束 使用 %，它 表示 该 水印 XPos 指定a percentage 对于 视频 宽度，such 作为 10% 表示 该 XPos 是 10% 的 视频 宽度.当 字符串 结束 使用 像素，它 表示 该 水印 XPos 是 指定 pixel，such 作为 100px 表示 该 XPos 是 100 pixels.默认值：0px。",
															},
															"y_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "vertical position 的 源站 的 水印 从 源站 的 coordinates 的 视频 镜像. Support %，像素 two formats:当 字符串 结束 使用 %，它 表示 该 水印 YPos 指定a percentage 对于 视频 高度，such 作为 10% 表示 该 YPos 是 10% 的 视频 高度.当 字符串 结束 使用 像素，它 表示 该 水印 YPos 是 指定 pixel，such 作为 100px 表示 该 YPos 是 100 pixels.默认值：0px。",
															},
															"image_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Image 水印 template，当 类型 是 镜像，此 字段 为必填项. 当 类型 是 text，此 字段 是 无效。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"image_content": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Required:    true,
																			Description: "input 内容 的 水印 镜像. Support jpeg，png 镜像 格式",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Enter 类型 来源 对象，其中 支持 COS 和 URL",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "有效 当 类型 是 COS，此 item 为必填项，indicating media processing COS 对象 信息。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "名称 COS 存储桶 其中 media processing 对象 文件 是 located。",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "park 到 其中 COS 存储桶 其中 media processing 目标 文件 resides belongs。",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "Input 路径 对于 media processing 对象 files。",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "有效 当 类型 是 URL，此 item 为必填项，indicating media processing URL 对象 信息.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "Video URL",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"width": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "宽度 的 水印. Support %，像素 two formats:当 字符串 结束 使用 %，它 表示 该 水印 宽度 是 percentage 的 视频 宽度，such 作为 10% 表示 该 宽度 是 10% 的 视频 宽度.当 字符串 结束 使用 像素，它 表示 该 水印 宽度 单位 是 pixels，such 作为 100px 表示 该 宽度 是 100 pixels.默认值：10%。",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "高度 的 水印. Support %，像素 two formats:当 字符串 结束 使用 %，它 表示 该 水印 高度 是 percentage 大小 的 视频 高度，such 作为 10% 表示 该 高度 是 10% 的 视频 高度.当 字符串 结束 使用 像素，它 表示 该 水印 高度 单位 是 pixel，such 作为 100px 表示 该 高度 是 100 pixels.默认值：0px，indicating 该 高度 是 scaled according 到 aspect ratio 的 original 水印 镜像。",
																		},
																		"repeat_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Watermark repeat 类型 Usage scenario: 水印 是 动态 镜像. Ranges.once: After 动态 水印 是 played，它 将 无 longer appear.repeat_last_frame: After 水印 是 played，stay 在 last frame.repeat: 水印 loops until end 的 视频 (默认值)。",
																		},
																	},
																},
															},
														},
													},
												},
												"text_content": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Text 内容， 长度 does 不 exceed 100 字符. Fill 在 仅 当 水印 类型 是 text 水印.Text 水印 does 不 support screenshot watermarking。",
												},
												"svg_content": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "SVG 内容 长度 不能 exceed 2000000 字符. Fill 在 仅 如果 水印 类型 是 SVG 水印.SVG 水印 does 不 support screenshot watermarking。",
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "开始时间 偏移量 的 水印，单位: second. Do 不 fill 在 或 fill 在 0，其中 表示 该 水印 将 start 到 appear 当 screen appears.Do 不 fill 在 或 fill 在 0，其中 表示 水印 将 appear 从 beginning 的 screen.当 值 是 greater 比 0 (assumed 到 是 n)，它 表示 该 水印 appears 从 nth second 的 screen.当 值 是 less 比 0 (assumed 到 是 -n)，它 表示 该 水印 starts 到 appear n 秒 before end 的 screen。",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "结束时间 偏移量 的 水印，单位: second.Do 不 fill 在 或 fill 在 0，indicating 该 水印 lasts until end 的 screen.当 值 是 greater 比 0 (assumed 到 是 n)，它 表示 该 水印 lasts until nth second 和 disappears.当 值 是 less 比 0 (assumed 到 是 -n)，它 表示 该 水印 lasts until 它 disappears n 秒 before end 的 screen。",
												},
											},
										},
									},
									"output_storage": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "目标 存储 的 文件 after screenshot 在 时间 point，如果未填写，它 将 inherit OutputStorage 值 的 upper layer.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "类型 media processing output 对象 存储 location，now 仅 支持 COS。",
												},
												"cos_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "有效 当 类型 是 COS，此 item 为必填项，indicating media processing COS output location.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bucket": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "目标 存储桶 名称 文件 output generated 通过 media processing，如果未填写，它 表示 upper layer。",
															},
															"region": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "park 的 目标 存储桶 对于 output 的 文件 generated 通过 media processing. 如果未填写，它 表示 inheriting 从 upper layer。",
															},
														},
													},
												},
											},
										},
									},
									"output_object_path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "输出路径 的 picture 文件 after 快照 在 时间 point 可以 是 relative 路径 或 absolute 路径 如果未填写， 默认为 relative 路径: `{inputName}_snapshotByTimeOffset_{definition}_{数量}.{格式}`。",
									},
									"object_number_format": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Rules 对于 `{数量}` variable 在 输出路径 after screenshot 在 时间 point.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"initial_value": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "starting 值 的 `{数量}` variable， 默认为 0。",
												},
												"increment": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "growth step 的 `{数量}` variable，默认为 1。",
												},
												"min_length": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "最小长度the `{数量}` variable，如果 insufficient，placeholders 将 是 filled. 默认为 1。",
												},
												"place_holder": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "当 长度 的 `{数量}` variable 是 insufficient， placeholder 是 added. 默认为 0。",
												},
											},
										},
									},
								},
							},
						},
						"sample_snapshot_task_set": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Screenshot 任务 列表 对于 视频 sampling。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Sample screenshot 模板 ID",
									},
									"watermark_set": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Watermark 列表，support 多个 pictures 或 text watermarks，up 到 10。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Watermark 模板 ID",
												},
												"raw_parameter": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Computed:    true,
													Description: "Watermark 自定义 参数，有效 当 Definition 是 filled 使用 0.此 参数 是 使用 在 highly customized scenarios，它 是 recommended 该 您 使用 Definition 到 指定watermark 参数 first.Watermark 自定义 参数 do 不 support screenshot watermarking。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Watermark 类型，可选 值:镜像: 镜像 水印。",
															},
															"coordinate_origin": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Origin position，currently 仅 支持:TopLeft: 表示that 源站 的 coordinates 是 在 upper left corner 的 视频 镜像，和 源站 的 水印 是 upper left corner 的 picture 或 text.默认值：TopLeft。",
															},
															"x_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "horizontal position 的 源站 的 水印 从 源站 的 coordinates 的 视频 镜像. Support %，像素 two formats:当 字符串 结束 使用 %，它 表示 该 水印 XPos 指定a percentage 对于 视频 宽度，such 作为 10% 表示 该 XPos 是 10% 的 视频 宽度.当 字符串 结束 使用 像素，它 表示 该 水印 XPos 是 指定 pixel，such 作为 100px 表示 该 XPos 是 100 pixels.默认值：0px。",
															},
															"y_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "vertical position 的 源站 的 水印 从 源站 的 coordinates 的 视频 镜像. Support %，像素 two formats:当 字符串 结束 使用 %，它 表示 该 水印 YPos 指定a percentage 对于 视频 高度，such 作为 10% 表示 该 YPos 是 10% 的 视频 高度.当 字符串 结束 使用 像素，它 表示 该 水印 YPos 是 指定 pixel，such 作为 100px 表示 该 YPos 是 100 pixels.默认值：0px。",
															},
															"image_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Image 水印 template，当 类型 是 镜像，此 字段 为必填项. 当 类型 是 text，此 字段 是 无效。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"image_content": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Required:    true,
																			Description: "input 内容 的 水印 镜像. Support jpeg，png 镜像 格式",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Enter 类型 来源 对象，其中 支持 COS 和 URL",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "有效 当 类型 是 COS，此 item 为必填项，indicating media processing COS 对象 信息。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "名称 COS 存储桶 其中 media processing 对象 文件 是 located。",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "park 到 其中 COS 存储桶 其中 media processing 目标 文件 resides belongs。",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "Input 路径 对于 media processing 对象 files。",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "有效 当 类型 是 URL，此 item 为必填项，indicating media processing URL 对象 信息.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "Video URL",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"width": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "宽度 的 水印. Support %，像素 two formats:当 字符串 结束 使用 %，它 表示 该 水印 宽度 是 percentage 的 视频 宽度，such 作为 10% 表示 该 宽度 是 10% 的 视频 宽度.当 字符串 结束 使用 像素，它 表示 该 水印 宽度 单位 是 pixels，such 作为 100px 表示 该 宽度 是 100 pixels.默认值：10%。",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "高度 的 水印. Support %，像素 two formats:当 字符串 结束 使用 %，它 表示 该 水印 高度 是 percentage 大小 的 视频 高度，such 作为 10% 表示 该 高度 是 10% 的 视频 高度.当 字符串 结束 使用 像素，它 表示 该 水印 高度 单位 是 pixel，such 作为 100px 表示 该 高度 是 100 pixels.默认值：0px，indicating 该 高度 是 scaled according 到 aspect ratio 的 original 水印 镜像。",
																		},
																		"repeat_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Watermark repeat 类型 Usage scenario: 水印 是 动态 镜像. Ranges.once: After 动态 水印 是 played，它 将 无 longer appear.repeat_last_frame: After 水印 是 played，stay 在 last frame.repeat: 水印 loops until end 的 视频 (默认值)。",
																		},
																	},
																},
															},
														},
													},
												},
												"text_content": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Text 内容， 长度 does 不 exceed 100 字符. Fill 在 仅 当 水印 类型 是 text 水印.Text 水印 does 不 support screenshot watermarking。",
												},
												"svg_content": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "SVG 内容 长度 不能 exceed 2000000 字符. Fill 在 仅 如果 水印 类型 是 SVG 水印.SVG 水印 does 不 support screenshot watermarking。",
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "开始时间 偏移量 的 水印，单位: second. Do 不 fill 在 或 fill 在 0，其中 表示 该 水印 将 start 到 appear 当 screen appears.Do 不 fill 在 或 fill 在 0，其中 表示 水印 将 appear 从 beginning 的 screen.当 值 是 greater 比 0 (assumed 到 是 n)，它 表示 该 水印 appears 从 nth second 的 screen.当 值 是 less 比 0 (assumed 到 是 -n)，它 表示 该 水印 starts 到 appear n 秒 before end 的 screen。",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "结束时间 偏移量 的 水印，单位: second.Do 不 fill 在 或 fill 在 0，indicating 该 水印 lasts until end 的 screen.当 值 是 greater 比 0 (assumed 到 是 n)，它 表示 该 水印 lasts until nth second 和 disappears.当 值 是 less 比 0 (assumed 到 是 -n)，它 表示 该 水印 lasts until 它 disappears n 秒 before end 的 screen。",
												},
											},
										},
									},
									"output_storage": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "目标 存储 的 文件 after screenshot 在 时间 point，如果未填写，它 将 inherit OutputStorage 值 的 upper layer.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "类型 media processing output 对象 存储 location，now 仅 支持 COS。",
												},
												"cos_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "有效 当 类型 是 COS，此 item 为必填项，indicating media processing COS output location.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bucket": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "目标 存储桶 名称 文件 output generated 通过 media processing，如果未填写，它 表示 upper layer。",
															},
															"region": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "park 的 目标 存储桶 对于 output 的 文件 generated 通过 media processing. 如果未填写，它 表示 inheriting 从 upper layer。",
															},
														},
													},
												},
											},
										},
									},
									"output_object_path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "输出路径 的 镜像 文件 after sampling screenshot，其中 可以 是 relative 路径 或 absolute 路径 如果未填写， 默认为 relative 路径: `{inputName}_sampleSnapshot_{definition}_{数量}.{格式}`。",
									},
									"object_number_format": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Rules 对于 `{数量}` variable 在 输出路径 after sampling screenshot.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"initial_value": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "starting 值 的 `{数量}` variable， 默认为 0。",
												},
												"increment": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "growth step 的 `{数量}` variable， 默认为 1。",
												},
												"min_length": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "最小长度the `{数量}` variable，如果 insufficient，placeholders 将 是 filled. 默认为 1。",
												},
												"place_holder": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "当 长度 的 `{数量}` variable 是 insufficient， placeholder 是 added. 默认为 0。",
												},
											},
										},
									},
								},
							},
						},
						"image_sprite_task_set": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Sprite 镜像 capture 任务 列表 对于 视频。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Sprite Illustration 模板 ID",
									},
									"output_storage": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "目标 存储 的 文件 after sprite 镜像 是 intercepted，如果未填写，它 将 inherit OutputStorage 值 的 upper layer.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "类型 media processing output 对象 存储 location，now 仅 支持 COS。",
												},
												"cos_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "有效 当 类型 是 COS，此 item 为必填项，indicating media processing COS output location.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bucket": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "目标 存储桶 名称 文件 output generated 通过 media processing，如果未填写，它 表示 upper layer。",
															},
															"region": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "park 的 目标 存储桶 对于 output 的 文件 generated 通过 media processing. 如果未填写，它 表示 inheriting 从 upper layer。",
															},
														},
													},
												},
											},
										},
									},
									"output_object_path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "After capturing sprite 镜像， 输出路径 的 sprite 镜像 文件 可以 是 relative 路径 或 absolute 路径 如果未填写， 默认为 relative 路径: `{inputName}_imageSprite_{definition}_{数量}.{格式}`。",
									},
									"web_vtt_object_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "After capturing sprite 镜像， 输出路径 的 Web VTT 文件 可以 仅 是 relative 路径 如果未填写， 默认为 relative 路径: `{inputName}_imageSprite_{definition}.{格式}`。",
									},
									"object_number_format": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Rules 对于 `{数量}` variable 在 输出路径 after intercepting Sprite 镜像.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"initial_value": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "starting 值 的 `{数量}` variable， 默认为 0。",
												},
												"increment": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "growth step 的 `{数量}` variable， 默认为 1。",
												},
												"min_length": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "最小长度the `{数量}` variable，如果 insufficient，placeholders 将 是 filled. 默认为 1。",
												},
												"place_holder": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "当 长度 的 `{数量}` variable 是 insufficient， placeholder 是 added. 默认为 0。",
												},
											},
										},
									},
								},
							},
						},
						"adaptive_dynamic_streaming_task_set": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Transfer Adaptive 代码 Stream 任务 List。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Transfer Adaptive 代码 Stream 模板 ID",
									},
									"watermark_set": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Watermark 列表，support 多个 pictures 或 text watermarks，up 到 10。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Watermark 模板 ID",
												},
												"raw_parameter": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Computed:    true,
													Description: "Watermark 自定义 参数，有效 当 Definition 是 filled 使用 0.此 参数 是 使用 在 highly customized scenarios，它 是 recommended 该 您 使用 Definition 到 指定watermark 参数 first.Watermark 自定义 参数 do 不 support screenshot watermarking。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Watermark 类型，可选 值:镜像: 镜像 水印。",
															},
															"coordinate_origin": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Origin position，currently 仅 支持:TopLeft: 表示that 源站 的 coordinates 是 在 upper left corner 的 视频 镜像，和 源站 的 水印 是 upper left corner 的 picture 或 text.默认值：TopLeft。",
															},
															"x_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "horizontal position 的 源站 的 水印 从 源站 的 coordinates 的 视频 镜像. Support %，像素 two formats:当 字符串 结束 使用 %，它 表示 该 水印 XPos 指定a percentage 对于 视频 宽度，such 作为 10% 表示 该 XPos 是 10% 的 视频 宽度.当 字符串 结束 使用 像素，它 表示 该 水印 XPos 是 指定 pixel，such 作为 100px 表示 该 XPos 是 100 pixels.默认值：0px。",
															},
															"y_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "vertical position 的 源站 的 水印 从 源站 的 coordinates 的 视频 镜像. Support %，像素 two formats:当 字符串 结束 使用 %，它 表示 该 水印 YPos 指定a percentage 对于 视频 高度，such 作为 10% 表示 该 YPos 是 10% 的 视频 高度.当 字符串 结束 使用 像素，它 表示 该 水印 YPos 是 指定 pixel，such 作为 100px 表示 该 YPos 是 100 pixels.默认值：0px。",
															},
															"image_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Image 水印 template，当 类型 是 镜像，此 字段 为必填项. 当 类型 是 text，此 字段 是 无效。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"image_content": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Required:    true,
																			Description: "input 内容 的 水印 镜像. Support jpeg，png 镜像 格式",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Enter 类型 来源 对象，其中 支持 COS 和 URL",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "有效 当 类型 是 COS，此 item 为必填项，indicating media processing COS 对象 信息。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "名称 COS 存储桶 其中 media processing 对象 文件 是 located。",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "park 到 其中 COS 存储桶 其中 media processing 目标 文件 resides belongs。",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "Input 路径 对于 media processing 对象 files。",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "有效 当 类型 是 URL，此 item 为必填项，indicating media processing URL 对象 信息.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "Video URL",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"width": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "宽度 的 水印. Support %，像素 two formats:当 字符串 结束 使用 %，它 表示 该 水印 宽度 是 percentage 的 视频 宽度，such 作为 10% 表示 该 宽度 是 10% 的 视频 宽度.当 字符串 结束 使用 像素，它 表示 该 水印 宽度 单位 是 pixels，such 作为 100px 表示 该 宽度 是 100 pixels.默认值：10%。",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "高度 的 水印. Support %，像素 two formats:当 字符串 结束 使用 %，它 表示 该 水印 高度 是 percentage 大小 的 视频 高度，such 作为 10% 表示 该 高度 是 10% 的 视频 高度.当 字符串 结束 使用 像素，它 表示 该 水印 高度 单位 是 pixel，such 作为 100px 表示 该 高度 是 100 pixels.默认值：0px，indicating 该 高度 是 scaled according 到 aspect ratio 的 original 水印 镜像。",
																		},
																		"repeat_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Watermark repeat 类型 Usage scenario: 水印 是 动态 镜像. Ranges.once: After 动态 水印 是 played，它 将 无 longer appear.repeat_last_frame: After 水印 是 played，stay 在 last frame.repeat: 水印 loops until end 的 视频 (默认值)。",
																		},
																	},
																},
															},
														},
													},
												},
												"text_content": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Text 内容， 长度 does 不 exceed 100 字符. Fill 在 仅 当 水印 类型 是 text 水印.Text 水印 does 不 support screenshot watermarking。",
												},
												"svg_content": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "SVG 内容 长度 不能 exceed 2000000 字符. Fill 在 仅 如果 水印 类型 是 SVG 水印.SVG 水印 does 不 support screenshot watermarking。",
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "开始时间 偏移量 的 水印，单位: second. Do 不 fill 在 或 fill 在 0，其中 表示 该 水印 将 start 到 appear 当 screen appears.Do 不 fill 在 或 fill 在 0，其中 表示 水印 将 appear 从 beginning 的 screen.当 值 是 greater 比 0 (assumed 到 是 n)，它 表示 该 水印 appears 从 nth second 的 screen.当 值 是 less 比 0 (assumed 到 是 -n)，它 表示 该 水印 starts 到 appear n 秒 before end 的 screen。",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "结束时间 偏移量 的 水印，单位: second.Do 不 fill 在 或 fill 在 0，indicating 该 水印 lasts until end 的 screen.当 值 是 greater 比 0 (assumed 到 是 n)，它 表示 该 水印 lasts until nth second 和 disappears.当 值 是 less 比 0 (assumed 到 是 -n)，它 表示 该 水印 lasts until 它 disappears n 秒 before end 的 screen。",
												},
											},
										},
									},
									"output_storage": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "目标 存储 的 文件 after converting 到 adaptive 代码 流，如果未填写，它 将 inherit OutputStorage 值 的 upper layer.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "类型 media processing output 对象 存储 location，now 仅 支持 COS。",
												},
												"cos_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "有效 当 类型 是 COS，此 item 为必填项，indicating media processing COS output location.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bucket": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "目标 存储桶 名称 文件 output generated 通过 media processing，如果未填写，它 表示 upper layer。",
															},
															"region": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "park 的 目标 存储桶 对于 output 的 文件 generated 通过 media processing. 如果未填写，它 表示 inheriting 从 upper layer。",
															},
														},
													},
												},
											},
										},
									},
									"output_object_path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "After converting 到 adaptive 流， 输出路径 的 manifest 文件 可以 是 relative 路径 或 absolute 路径 如果未填写， 默认为 relative 路径: `{inputName}_adaptiveDynamicStreaming_{definition}.{格式}`。",
									},
									"sub_stream_object_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "After converting 到 adaptive 流， 输出路径 的 sub-流 文件 可以 仅 是 relative 路径 如果未填写， 默认为 relative 路径: {inputName}_adaptiveDynamicStreaming_{definition}_{subStreamNumber}.{格式}`。",
									},
									"segment_object_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "After converting 到 adaptive 流 (仅 HLS)， 输出路径 的 fragmented 文件 可以 仅 是 relative 路径 如果未填写， 默认为 relative 路径: `{inputName}_adaptiveDynamicStreaming_{definition}_{subStreamNumber}_{segmentNumber}.{格式}`。",
									},
								},
							},
						},
					},
				},
			},

			"ai_content_review_task": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Video 内容 Moderation 类型 任务 Parameters。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Video 内容 Review 模板 ID",
						},
					},
				},
			},

			"ai_analysis_task": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Video 内容 Analysis 类型 任务 Parameters。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Video 内容 Analysis 模板 ID",
						},
						"extended_parameter": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Extension 参数 whose 值 是 serialized json 字符串.注意: 此 参数 是 customized demand 参数，其中 requires offline docking.注意：此字段可能返回 null，表示无法获取有效值。",
						},
					},
				},
			},

			"ai_recognition_task": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Video 内容 recognition 类型 任务 参数。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Video Intelligent Recognition 模板 ID",
						},
					},
				},
			},

			"task_notify_config": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "事件 通知 配置 的 任务，如果 它 是 不 filled，它 表示 该 事件 通知 将 不 是 获取。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cmq_model": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CMQ 或 TDMQ-CMQ model，there 是 two kinds 的 Queue 和 Topic。",
						},
						"cmq_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "地域 的 CMQ 或 TDMQ-CMQ，such 作为 sh，bj，etc。",
						},
						"topic_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "有效 当 model 是 Topic，indicating 主题 名称 CMQ 或 TDMQ-CMQ 该 receives 事件 notifications。",
						},
						"queue_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "有效 当 model 是 Queue，indicating queue 名称 CMQ 或 TDMQ-CMQ 该 receives 事件 通知。",
						},
						"notify_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "模式 的 工作流 通知， possible 值 是 Finish 和 Change，leaving blank 表示 Finish。",
						},
						"notify_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Notification 类型，可选 值:CMQ: offline，它 是 recommended 到 switch 到 TDMQ-CMQ.TDMQ-CMQ: 消息 queue.URL: 当 URL 是 指定， HTTP callback 是 pushed 到 地址 指定 通过 NotifyUrl， callback 协议 是 http+json，和 包 正文 内容 是 same 作为 output 参数 的 parsing 事件 通知 interface.SCF: 不 recommended，additional 配置 的 SCF 在 console 为必填项.注意: CMQ 是 默认值 当 不 filled 或 空，如果 您 need 到 使用 other types，您 need 到 fill 在 corresponding 类型 值",
						},
						"notify_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "HTTP callback 地址，必填 当 NotifyType 是 URL",
						},
					},
				},
			},

			"task_priority": {
				Optional:    true,
				Type:        schema.TypeInt,
				Default:     0,
				Description: "优先级 的 工作流， larger 值， higher 优先级， 值 范围 是 -10 到 10，和 blank 表示 0。",
			},
		},
	}
}

func resourceTencentCloudMpsWorkflowCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_workflow.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = mps.NewCreateWorkflowRequest()
		response   = mps.NewCreateWorkflowResponse()
		workflowId int64
	)
	if v, ok := d.GetOk("workflow_name"); ok {
		request.WorkflowName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "trigger"); ok {
		workflowTrigger := mps.WorkflowTrigger{}
		if v, ok := dMap["type"]; ok {
			workflowTrigger.Type = helper.String(v.(string))
		}
		if cosFileUploadTriggerMap, ok := helper.InterfaceToMap(dMap, "cos_file_upload_trigger"); ok {
			cosFileUploadTrigger := mps.CosFileUploadTrigger{}
			if v, ok := cosFileUploadTriggerMap["bucket"]; ok {
				cosFileUploadTrigger.Bucket = helper.String(v.(string))
			}
			if v, ok := cosFileUploadTriggerMap["region"]; ok {
				cosFileUploadTrigger.Region = helper.String(v.(string))
			}
			if v, ok := cosFileUploadTriggerMap["dir"]; ok {
				cosFileUploadTrigger.Dir = helper.String(v.(string))
			}
			if v, ok := cosFileUploadTriggerMap["formats"]; ok {
				formatsSet := v.(*schema.Set).List()
				for i := range formatsSet {
					formats := formatsSet[i].(string)
					cosFileUploadTrigger.Formats = append(cosFileUploadTrigger.Formats, &formats)
				}
			}
			workflowTrigger.CosFileUploadTrigger = &cosFileUploadTrigger
		}
		request.Trigger = &workflowTrigger
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "output_storage"); ok {
		taskOutputStorage := mps.TaskOutputStorage{}
		if v, ok := dMap["type"]; ok {
			taskOutputStorage.Type = helper.String(v.(string))
		}
		if cosOutputStorageMap, ok := helper.InterfaceToMap(dMap, "cos_output_storage"); ok {
			cosOutputStorage := mps.CosOutputStorage{}
			if v, ok := cosOutputStorageMap["bucket"]; ok {
				cosOutputStorage.Bucket = helper.String(v.(string))
			}
			if v, ok := cosOutputStorageMap["region"]; ok {
				cosOutputStorage.Region = helper.String(v.(string))
			}
			taskOutputStorage.CosOutputStorage = &cosOutputStorage
		}
		request.OutputStorage = &taskOutputStorage
	}

	if v, ok := d.GetOk("output_dir"); ok {
		request.OutputDir = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "media_process_task"); ok {
		mediaProcessTaskInput := mps.MediaProcessTaskInput{}
		if v, ok := dMap["transcode_task_set"]; ok {
			for _, item := range v.([]interface{}) {
				transcodeTaskSetMap := item.(map[string]interface{})
				transcodeTaskInput := mps.TranscodeTaskInput{}
				if v, ok := transcodeTaskSetMap["definition"]; ok {
					transcodeTaskInput.Definition = helper.IntUint64(v.(int))
				}
				if rawParameterMap, ok := helper.InterfaceToMap(transcodeTaskSetMap, "raw_parameter"); ok {
					rawTranscodeParameter := mps.RawTranscodeParameter{}
					if v, ok := rawParameterMap["container"]; ok {
						rawTranscodeParameter.Container = helper.String(v.(string))
					}
					if v, ok := rawParameterMap["remove_video"]; ok {
						rawTranscodeParameter.RemoveVideo = helper.IntInt64(v.(int))
					}
					if v, ok := rawParameterMap["remove_audio"]; ok {
						rawTranscodeParameter.RemoveAudio = helper.IntInt64(v.(int))
					}
					if videoTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "video_template"); ok {
						videoTemplateInfo := mps.VideoTemplateInfo{}
						if v, ok := videoTemplateMap["codec"]; ok {
							videoTemplateInfo.Codec = helper.String(v.(string))
						}
						if v, ok := videoTemplateMap["fps"]; ok {
							videoTemplateInfo.Fps = helper.IntInt64(v.(int))
						}
						if v, ok := videoTemplateMap["bitrate"]; ok {
							videoTemplateInfo.Bitrate = helper.IntInt64(v.(int))
						}
						if v, ok := videoTemplateMap["resolution_adaptive"]; ok {
							videoTemplateInfo.ResolutionAdaptive = helper.String(v.(string))
						}
						if v, ok := videoTemplateMap["width"]; ok {
							videoTemplateInfo.Width = helper.IntUint64(v.(int))
						}
						if v, ok := videoTemplateMap["height"]; ok {
							videoTemplateInfo.Height = helper.IntUint64(v.(int))
						}
						if v, ok := videoTemplateMap["gop"]; ok {
							videoTemplateInfo.Gop = helper.IntUint64(v.(int))
						}
						if v, ok := videoTemplateMap["fill_type"]; ok {
							videoTemplateInfo.FillType = helper.String(v.(string))
						}
						if v, ok := videoTemplateMap["vcrf"]; ok {
							videoTemplateInfo.Vcrf = helper.IntUint64(v.(int))
						}
						rawTranscodeParameter.VideoTemplate = &videoTemplateInfo
					}
					if audioTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "audio_template"); ok {
						audioTemplateInfo := mps.AudioTemplateInfo{}
						if v, ok := audioTemplateMap["codec"]; ok {
							audioTemplateInfo.Codec = helper.String(v.(string))
						}
						if v, ok := audioTemplateMap["bitrate"]; ok {
							audioTemplateInfo.Bitrate = helper.IntInt64(v.(int))
						}
						if v, ok := audioTemplateMap["sample_rate"]; ok {
							audioTemplateInfo.SampleRate = helper.IntUint64(v.(int))
						}
						if v, ok := audioTemplateMap["audio_channel"]; ok {
							audioTemplateInfo.AudioChannel = helper.IntInt64(v.(int))
						}
						rawTranscodeParameter.AudioTemplate = &audioTemplateInfo
					}
					if tEHDConfigMap, ok := helper.InterfaceToMap(rawParameterMap, "tehd_config"); ok {
						tEHDConfig := mps.TEHDConfig{}
						if v, ok := tEHDConfigMap["type"]; ok {
							tEHDConfig.Type = helper.String(v.(string))
						}
						if v, ok := tEHDConfigMap["max_video_bitrate"]; ok {
							tEHDConfig.MaxVideoBitrate = helper.IntInt64(v.(int))
						}
						rawTranscodeParameter.TEHDConfig = &tEHDConfig
					}
					transcodeTaskInput.RawParameter = &rawTranscodeParameter
				}
				if overrideParameterMap, ok := helper.InterfaceToMap(transcodeTaskSetMap, "override_parameter"); ok {
					overrideTranscodeParameter := mps.OverrideTranscodeParameter{}
					if v, ok := overrideParameterMap["container"]; ok {
						overrideTranscodeParameter.Container = helper.String(v.(string))
					}
					if v, ok := overrideParameterMap["remove_video"]; ok {
						overrideTranscodeParameter.RemoveVideo = helper.IntUint64(v.(int))
					}
					if v, ok := overrideParameterMap["remove_audio"]; ok {
						overrideTranscodeParameter.RemoveAudio = helper.IntUint64(v.(int))
					}
					if videoTemplateMap, ok := helper.InterfaceToMap(overrideParameterMap, "video_template"); ok {
						videoTemplateInfoForUpdate := mps.VideoTemplateInfoForUpdate{}
						if v, ok := videoTemplateMap["codec"]; ok {
							videoTemplateInfoForUpdate.Codec = helper.String(v.(string))
						}
						if v, ok := videoTemplateMap["fps"]; ok {
							videoTemplateInfoForUpdate.Fps = helper.IntInt64(v.(int))
						}
						if v, ok := videoTemplateMap["bitrate"]; ok {
							videoTemplateInfoForUpdate.Bitrate = helper.IntInt64(v.(int))
						}
						if v, ok := videoTemplateMap["resolution_adaptive"]; ok {
							videoTemplateInfoForUpdate.ResolutionAdaptive = helper.String(v.(string))
						}
						if v, ok := videoTemplateMap["width"]; ok {
							videoTemplateInfoForUpdate.Width = helper.IntUint64(v.(int))
						}
						if v, ok := videoTemplateMap["height"]; ok {
							videoTemplateInfoForUpdate.Height = helper.IntUint64(v.(int))
						}
						if v, ok := videoTemplateMap["gop"]; ok {
							videoTemplateInfoForUpdate.Gop = helper.IntUint64(v.(int))
						}
						if v, ok := videoTemplateMap["fill_type"]; ok {
							videoTemplateInfoForUpdate.FillType = helper.String(v.(string))
						}
						if v, ok := videoTemplateMap["vcrf"]; ok {
							videoTemplateInfoForUpdate.Vcrf = helper.IntUint64(v.(int))
						}
						if v, ok := videoTemplateMap["content_adapt_stream"]; ok {
							videoTemplateInfoForUpdate.ContentAdaptStream = helper.IntUint64(v.(int))
						}
						overrideTranscodeParameter.VideoTemplate = &videoTemplateInfoForUpdate
					}
					if audioTemplateMap, ok := helper.InterfaceToMap(overrideParameterMap, "audio_template"); ok {
						audioTemplateInfoForUpdate := mps.AudioTemplateInfoForUpdate{}
						if v, ok := audioTemplateMap["codec"]; ok {
							audioTemplateInfoForUpdate.Codec = helper.String(v.(string))
						}
						if v, ok := audioTemplateMap["bitrate"]; ok {
							audioTemplateInfoForUpdate.Bitrate = helper.IntInt64(v.(int))
						}
						if v, ok := audioTemplateMap["sample_rate"]; ok {
							audioTemplateInfoForUpdate.SampleRate = helper.IntUint64(v.(int))
						}
						if v, ok := audioTemplateMap["audio_channel"]; ok {
							audioTemplateInfoForUpdate.AudioChannel = helper.IntInt64(v.(int))
						}
						if v, ok := audioTemplateMap["stream_selects"]; ok {
							streamSelectsSet := v.(*schema.Set).List()
							for i := range streamSelectsSet {
								streamSelects := streamSelectsSet[i].(int)
								audioTemplateInfoForUpdate.StreamSelects = append(audioTemplateInfoForUpdate.StreamSelects, helper.IntInt64(streamSelects))
							}
						}
						overrideTranscodeParameter.AudioTemplate = &audioTemplateInfoForUpdate
					}
					if tEHDConfigMap, ok := helper.InterfaceToMap(overrideParameterMap, "tehd_config"); ok {
						tEHDConfigForUpdate := mps.TEHDConfigForUpdate{}
						if v, ok := tEHDConfigMap["type"]; ok {
							tEHDConfigForUpdate.Type = helper.String(v.(string))
						}
						if v, ok := tEHDConfigMap["max_video_bitrate"]; ok {
							tEHDConfigForUpdate.MaxVideoBitrate = helper.IntInt64(v.(int))
						}
						overrideTranscodeParameter.TEHDConfig = &tEHDConfigForUpdate
					}
					if subtitleTemplateMap, ok := helper.InterfaceToMap(overrideParameterMap, "subtitle_template"); ok {
						subtitleTemplate := mps.SubtitleTemplate{}
						if v, ok := subtitleTemplateMap["path"]; ok {
							subtitleTemplate.Path = helper.String(v.(string))
						}
						if v, ok := subtitleTemplateMap["stream_index"]; ok {
							subtitleTemplate.StreamIndex = helper.IntInt64(v.(int))
						}
						if v, ok := subtitleTemplateMap["font_type"]; ok {
							subtitleTemplate.FontType = helper.String(v.(string))
						}
						if v, ok := subtitleTemplateMap["font_size"]; ok {
							subtitleTemplate.FontSize = helper.String(v.(string))
						}
						if v, ok := subtitleTemplateMap["font_color"]; ok {
							subtitleTemplate.FontColor = helper.String(v.(string))
						}
						if v, ok := subtitleTemplateMap["font_alpha"]; ok {
							subtitleTemplate.FontAlpha = helper.Float64(v.(float64))
						}
						overrideTranscodeParameter.SubtitleTemplate = &subtitleTemplate
					}
					transcodeTaskInput.OverrideParameter = &overrideTranscodeParameter
				}
				if v, ok := transcodeTaskSetMap["watermark_set"]; ok {
					for _, item := range v.([]interface{}) {
						watermarkSetMap := item.(map[string]interface{})
						watermarkInput := mps.WatermarkInput{}
						if v, ok := watermarkSetMap["definition"]; ok {
							watermarkInput.Definition = helper.IntUint64(v.(int))
						}
						if rawParameterMap, ok := helper.InterfaceToMap(watermarkSetMap, "raw_parameter"); ok {
							rawWatermarkParameter := mps.RawWatermarkParameter{}
							if v, ok := rawParameterMap["type"]; ok {
								rawWatermarkParameter.Type = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["coordinate_origin"]; ok {
								rawWatermarkParameter.CoordinateOrigin = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["x_pos"]; ok {
								rawWatermarkParameter.XPos = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["y_pos"]; ok {
								rawWatermarkParameter.YPos = helper.String(v.(string))
							}
							if imageTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "image_template"); ok {
								rawImageWatermarkInput := mps.RawImageWatermarkInput{}
								if imageContentMap, ok := helper.InterfaceToMap(imageTemplateMap, "image_content"); ok {
									mediaInputInfo := mps.MediaInputInfo{}
									if v, ok := imageContentMap["type"]; ok {
										mediaInputInfo.Type = helper.String(v.(string))
									}
									if cosInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "cos_input_info"); ok {
										cosInputInfo := mps.CosInputInfo{}
										if v, ok := cosInputInfoMap["bucket"]; ok {
											cosInputInfo.Bucket = helper.String(v.(string))
										}
										if v, ok := cosInputInfoMap["region"]; ok {
											cosInputInfo.Region = helper.String(v.(string))
										}
										if v, ok := cosInputInfoMap["object"]; ok {
											cosInputInfo.Object = helper.String(v.(string))
										}
										mediaInputInfo.CosInputInfo = &cosInputInfo
									}
									if urlInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "url_input_info"); ok {
										urlInputInfo := mps.UrlInputInfo{}
										if v, ok := urlInputInfoMap["url"]; ok {
											urlInputInfo.Url = helper.String(v.(string))
										}
										mediaInputInfo.UrlInputInfo = &urlInputInfo
									}
									rawImageWatermarkInput.ImageContent = &mediaInputInfo
								}
								if v, ok := imageTemplateMap["width"]; ok {
									rawImageWatermarkInput.Width = helper.String(v.(string))
								}
								if v, ok := imageTemplateMap["height"]; ok {
									rawImageWatermarkInput.Height = helper.String(v.(string))
								}
								if v, ok := imageTemplateMap["repeat_type"]; ok {
									rawImageWatermarkInput.RepeatType = helper.String(v.(string))
								}
								rawWatermarkParameter.ImageTemplate = &rawImageWatermarkInput
							}
							watermarkInput.RawParameter = &rawWatermarkParameter
						}
						if v, ok := watermarkSetMap["text_content"]; ok {
							watermarkInput.TextContent = helper.String(v.(string))
						}
						if v, ok := watermarkSetMap["svg_content"]; ok {
							watermarkInput.SvgContent = helper.String(v.(string))
						}
						if v, ok := watermarkSetMap["start_time_offset"]; ok {
							watermarkInput.StartTimeOffset = helper.Float64(v.(float64))
						}
						if v, ok := watermarkSetMap["end_time_offset"]; ok {
							watermarkInput.EndTimeOffset = helper.Float64(v.(float64))
						}
						transcodeTaskInput.WatermarkSet = append(transcodeTaskInput.WatermarkSet, &watermarkInput)
					}
				}
				if v, ok := transcodeTaskSetMap["mosaic_set"]; ok {
					for _, item := range v.([]interface{}) {
						mosaicSetMap := item.(map[string]interface{})
						mosaicInput := mps.MosaicInput{}
						if v, ok := mosaicSetMap["coordinate_origin"]; ok {
							mosaicInput.CoordinateOrigin = helper.String(v.(string))
						}
						if v, ok := mosaicSetMap["x_pos"]; ok {
							mosaicInput.XPos = helper.String(v.(string))
						}
						if v, ok := mosaicSetMap["y_pos"]; ok {
							mosaicInput.YPos = helper.String(v.(string))
						}
						if v, ok := mosaicSetMap["width"]; ok {
							mosaicInput.Width = helper.String(v.(string))
						}
						if v, ok := mosaicSetMap["height"]; ok {
							mosaicInput.Height = helper.String(v.(string))
						}
						if v, ok := mosaicSetMap["start_time_offset"]; ok {
							mosaicInput.StartTimeOffset = helper.Float64(v.(float64))
						}
						if v, ok := mosaicSetMap["end_time_offset"]; ok {
							mosaicInput.EndTimeOffset = helper.Float64(v.(float64))
						}
						transcodeTaskInput.MosaicSet = append(transcodeTaskInput.MosaicSet, &mosaicInput)
					}
				}
				if v, ok := transcodeTaskSetMap["start_time_offset"]; ok {
					transcodeTaskInput.StartTimeOffset = helper.Float64(v.(float64))
				}
				if v, ok := transcodeTaskSetMap["end_time_offset"]; ok {
					transcodeTaskInput.EndTimeOffset = helper.Float64(v.(float64))
				}
				if outputStorageMap, ok := helper.InterfaceToMap(transcodeTaskSetMap, "output_storage"); ok {
					taskOutputStorage := mps.TaskOutputStorage{}
					if v, ok := outputStorageMap["type"]; ok {
						taskOutputStorage.Type = helper.String(v.(string))
					}
					if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
						cosOutputStorage := mps.CosOutputStorage{}
						if v, ok := cosOutputStorageMap["bucket"]; ok {
							cosOutputStorage.Bucket = helper.String(v.(string))
						}
						if v, ok := cosOutputStorageMap["region"]; ok {
							cosOutputStorage.Region = helper.String(v.(string))
						}
						taskOutputStorage.CosOutputStorage = &cosOutputStorage
					}
					transcodeTaskInput.OutputStorage = &taskOutputStorage
				}
				if v, ok := transcodeTaskSetMap["output_object_path"]; ok {
					transcodeTaskInput.OutputObjectPath = helper.String(v.(string))
				}
				if v, ok := transcodeTaskSetMap["segment_object_name"]; ok {
					transcodeTaskInput.SegmentObjectName = helper.String(v.(string))
				}
				if objectNumberFormatMap, ok := helper.InterfaceToMap(transcodeTaskSetMap, "object_number_format"); ok {
					numberFormat := mps.NumberFormat{}
					if v, ok := objectNumberFormatMap["initial_value"]; ok {
						numberFormat.InitialValue = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["increment"]; ok {
						numberFormat.Increment = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["min_length"]; ok {
						numberFormat.MinLength = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["place_holder"]; ok {
						numberFormat.PlaceHolder = helper.String(v.(string))
					}
					transcodeTaskInput.ObjectNumberFormat = &numberFormat
				}
				if headTailParameterMap, ok := helper.InterfaceToMap(transcodeTaskSetMap, "head_tail_parameter"); ok {
					headTailParameter := mps.HeadTailParameter{}
					if v, ok := headTailParameterMap["head_set"]; ok {
						for _, item := range v.([]interface{}) {
							headSetMap := item.(map[string]interface{})
							mediaInputInfo := mps.MediaInputInfo{}
							if v, ok := headSetMap["type"]; ok {
								mediaInputInfo.Type = helper.String(v.(string))
							}
							if cosInputInfoMap, ok := helper.InterfaceToMap(headSetMap, "cos_input_info"); ok {
								cosInputInfo := mps.CosInputInfo{}
								if v, ok := cosInputInfoMap["bucket"]; ok {
									cosInputInfo.Bucket = helper.String(v.(string))
								}
								if v, ok := cosInputInfoMap["region"]; ok {
									cosInputInfo.Region = helper.String(v.(string))
								}
								if v, ok := cosInputInfoMap["object"]; ok {
									cosInputInfo.Object = helper.String(v.(string))
								}
								mediaInputInfo.CosInputInfo = &cosInputInfo
							}
							if urlInputInfoMap, ok := helper.InterfaceToMap(headSetMap, "url_input_info"); ok {
								urlInputInfo := mps.UrlInputInfo{}
								if v, ok := urlInputInfoMap["url"]; ok {
									urlInputInfo.Url = helper.String(v.(string))
								}
								mediaInputInfo.UrlInputInfo = &urlInputInfo
							}
							headTailParameter.HeadSet = append(headTailParameter.HeadSet, &mediaInputInfo)
						}
					}
					if v, ok := headTailParameterMap["tail_set"]; ok {
						for _, item := range v.([]interface{}) {
							tailSetMap := item.(map[string]interface{})
							mediaInputInfo := mps.MediaInputInfo{}
							if v, ok := tailSetMap["type"]; ok {
								mediaInputInfo.Type = helper.String(v.(string))
							}
							if cosInputInfoMap, ok := helper.InterfaceToMap(tailSetMap, "cos_input_info"); ok {
								cosInputInfo := mps.CosInputInfo{}
								if v, ok := cosInputInfoMap["bucket"]; ok {
									cosInputInfo.Bucket = helper.String(v.(string))
								}
								if v, ok := cosInputInfoMap["region"]; ok {
									cosInputInfo.Region = helper.String(v.(string))
								}
								if v, ok := cosInputInfoMap["object"]; ok {
									cosInputInfo.Object = helper.String(v.(string))
								}
								mediaInputInfo.CosInputInfo = &cosInputInfo
							}
							if urlInputInfoMap, ok := helper.InterfaceToMap(tailSetMap, "url_input_info"); ok {
								urlInputInfo := mps.UrlInputInfo{}
								if v, ok := urlInputInfoMap["url"]; ok {
									urlInputInfo.Url = helper.String(v.(string))
								}
								mediaInputInfo.UrlInputInfo = &urlInputInfo
							}
							headTailParameter.TailSet = append(headTailParameter.TailSet, &mediaInputInfo)
						}
					}
					transcodeTaskInput.HeadTailParameter = &headTailParameter
				}
				mediaProcessTaskInput.TranscodeTaskSet = append(mediaProcessTaskInput.TranscodeTaskSet, &transcodeTaskInput)
			}
		}
		if v, ok := dMap["animated_graphic_task_set"]; ok {
			for _, item := range v.([]interface{}) {
				animatedGraphicTaskSetMap := item.(map[string]interface{})
				animatedGraphicTaskInput := mps.AnimatedGraphicTaskInput{}
				if v, ok := animatedGraphicTaskSetMap["definition"]; ok {
					animatedGraphicTaskInput.Definition = helper.IntUint64(v.(int))
				}
				if v, ok := animatedGraphicTaskSetMap["start_time_offset"]; ok {
					animatedGraphicTaskInput.StartTimeOffset = helper.Float64(v.(float64))
				}
				if v, ok := animatedGraphicTaskSetMap["end_time_offset"]; ok {
					animatedGraphicTaskInput.EndTimeOffset = helper.Float64(v.(float64))
				}
				if outputStorageMap, ok := helper.InterfaceToMap(animatedGraphicTaskSetMap, "output_storage"); ok {
					taskOutputStorage := mps.TaskOutputStorage{}
					if v, ok := outputStorageMap["type"]; ok {
						taskOutputStorage.Type = helper.String(v.(string))
					}
					if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
						cosOutputStorage := mps.CosOutputStorage{}
						if v, ok := cosOutputStorageMap["bucket"]; ok {
							cosOutputStorage.Bucket = helper.String(v.(string))
						}
						if v, ok := cosOutputStorageMap["region"]; ok {
							cosOutputStorage.Region = helper.String(v.(string))
						}
						taskOutputStorage.CosOutputStorage = &cosOutputStorage
					}
					animatedGraphicTaskInput.OutputStorage = &taskOutputStorage
				}
				if v, ok := animatedGraphicTaskSetMap["output_object_path"]; ok {
					animatedGraphicTaskInput.OutputObjectPath = helper.String(v.(string))
				}
				mediaProcessTaskInput.AnimatedGraphicTaskSet = append(mediaProcessTaskInput.AnimatedGraphicTaskSet, &animatedGraphicTaskInput)
			}
		}
		if v, ok := dMap["snapshot_by_time_offset_task_set"]; ok {
			for _, item := range v.([]interface{}) {
				snapshotByTimeOffsetTaskSetMap := item.(map[string]interface{})
				snapshotByTimeOffsetTaskInput := mps.SnapshotByTimeOffsetTaskInput{}
				if v, ok := snapshotByTimeOffsetTaskSetMap["definition"]; ok {
					snapshotByTimeOffsetTaskInput.Definition = helper.IntUint64(v.(int))
				}
				if v, ok := snapshotByTimeOffsetTaskSetMap["ext_time_offset_set"]; ok {
					extTimeOffsetSetSet := v.(*schema.Set).List()
					for i := range extTimeOffsetSetSet {
						extTimeOffsetSet := extTimeOffsetSetSet[i].(string)
						snapshotByTimeOffsetTaskInput.ExtTimeOffsetSet = append(snapshotByTimeOffsetTaskInput.ExtTimeOffsetSet, &extTimeOffsetSet)
					}
				}
				if v, _ := d.GetOk("time_offset_set"); v != nil {
					timeOffsetSetSet := v.(*schema.Set).List()
					for i := range timeOffsetSetSet {
						timeOffsetSet := timeOffsetSetSet[i].(float64)
						snapshotByTimeOffsetTaskInput.TimeOffsetSet = append(snapshotByTimeOffsetTaskInput.TimeOffsetSet, &timeOffsetSet)
					}
				}

				if v, ok := snapshotByTimeOffsetTaskSetMap["watermark_set"]; ok {
					for _, item := range v.([]interface{}) {
						watermarkSetMap := item.(map[string]interface{})
						watermarkInput := mps.WatermarkInput{}
						if v, ok := watermarkSetMap["definition"]; ok {
							watermarkInput.Definition = helper.IntUint64(v.(int))
						}
						if rawParameterMap, ok := helper.InterfaceToMap(watermarkSetMap, "raw_parameter"); ok {
							rawWatermarkParameter := mps.RawWatermarkParameter{}
							if v, ok := rawParameterMap["type"]; ok {
								rawWatermarkParameter.Type = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["coordinate_origin"]; ok {
								rawWatermarkParameter.CoordinateOrigin = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["x_pos"]; ok {
								rawWatermarkParameter.XPos = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["y_pos"]; ok {
								rawWatermarkParameter.YPos = helper.String(v.(string))
							}
							if imageTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "image_template"); ok {
								rawImageWatermarkInput := mps.RawImageWatermarkInput{}
								if imageContentMap, ok := helper.InterfaceToMap(imageTemplateMap, "image_content"); ok {
									mediaInputInfo := mps.MediaInputInfo{}
									if v, ok := imageContentMap["type"]; ok {
										mediaInputInfo.Type = helper.String(v.(string))
									}
									if cosInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "cos_input_info"); ok {
										cosInputInfo := mps.CosInputInfo{}
										if v, ok := cosInputInfoMap["bucket"]; ok {
											cosInputInfo.Bucket = helper.String(v.(string))
										}
										if v, ok := cosInputInfoMap["region"]; ok {
											cosInputInfo.Region = helper.String(v.(string))
										}
										if v, ok := cosInputInfoMap["object"]; ok {
											cosInputInfo.Object = helper.String(v.(string))
										}
										mediaInputInfo.CosInputInfo = &cosInputInfo
									}
									if urlInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "url_input_info"); ok {
										urlInputInfo := mps.UrlInputInfo{}
										if v, ok := urlInputInfoMap["url"]; ok {
											urlInputInfo.Url = helper.String(v.(string))
										}
										mediaInputInfo.UrlInputInfo = &urlInputInfo
									}
									rawImageWatermarkInput.ImageContent = &mediaInputInfo
								}
								if v, ok := imageTemplateMap["width"]; ok {
									rawImageWatermarkInput.Width = helper.String(v.(string))
								}
								if v, ok := imageTemplateMap["height"]; ok {
									rawImageWatermarkInput.Height = helper.String(v.(string))
								}
								if v, ok := imageTemplateMap["repeat_type"]; ok {
									rawImageWatermarkInput.RepeatType = helper.String(v.(string))
								}
								rawWatermarkParameter.ImageTemplate = &rawImageWatermarkInput
							}
							watermarkInput.RawParameter = &rawWatermarkParameter
						}
						if v, ok := watermarkSetMap["text_content"]; ok {
							watermarkInput.TextContent = helper.String(v.(string))
						}
						if v, ok := watermarkSetMap["svg_content"]; ok {
							watermarkInput.SvgContent = helper.String(v.(string))
						}
						if v, ok := watermarkSetMap["start_time_offset"]; ok {
							watermarkInput.StartTimeOffset = helper.Float64(v.(float64))
						}
						if v, ok := watermarkSetMap["end_time_offset"]; ok {
							watermarkInput.EndTimeOffset = helper.Float64(v.(float64))
						}
						snapshotByTimeOffsetTaskInput.WatermarkSet = append(snapshotByTimeOffsetTaskInput.WatermarkSet, &watermarkInput)
					}
				}
				if outputStorageMap, ok := helper.InterfaceToMap(snapshotByTimeOffsetTaskSetMap, "output_storage"); ok {
					taskOutputStorage := mps.TaskOutputStorage{}
					if v, ok := outputStorageMap["type"]; ok {
						taskOutputStorage.Type = helper.String(v.(string))
					}
					if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
						cosOutputStorage := mps.CosOutputStorage{}
						if v, ok := cosOutputStorageMap["bucket"]; ok {
							cosOutputStorage.Bucket = helper.String(v.(string))
						}
						if v, ok := cosOutputStorageMap["region"]; ok {
							cosOutputStorage.Region = helper.String(v.(string))
						}
						taskOutputStorage.CosOutputStorage = &cosOutputStorage
					}
					snapshotByTimeOffsetTaskInput.OutputStorage = &taskOutputStorage
				}
				if v, ok := snapshotByTimeOffsetTaskSetMap["output_object_path"]; ok {
					snapshotByTimeOffsetTaskInput.OutputObjectPath = helper.String(v.(string))
				}
				if objectNumberFormatMap, ok := helper.InterfaceToMap(snapshotByTimeOffsetTaskSetMap, "object_number_format"); ok {
					numberFormat := mps.NumberFormat{}
					if v, ok := objectNumberFormatMap["initial_value"]; ok {
						numberFormat.InitialValue = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["increment"]; ok {
						numberFormat.Increment = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["min_length"]; ok {
						numberFormat.MinLength = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["place_holder"]; ok {
						numberFormat.PlaceHolder = helper.String(v.(string))
					}
					snapshotByTimeOffsetTaskInput.ObjectNumberFormat = &numberFormat
				}
				mediaProcessTaskInput.SnapshotByTimeOffsetTaskSet = append(mediaProcessTaskInput.SnapshotByTimeOffsetTaskSet, &snapshotByTimeOffsetTaskInput)
			}
		}
		if v, ok := dMap["sample_snapshot_task_set"]; ok {
			for _, item := range v.([]interface{}) {
				sampleSnapshotTaskSetMap := item.(map[string]interface{})
				sampleSnapshotTaskInput := mps.SampleSnapshotTaskInput{}
				if v, ok := sampleSnapshotTaskSetMap["definition"]; ok {
					sampleSnapshotTaskInput.Definition = helper.IntUint64(v.(int))
				}
				if v, ok := sampleSnapshotTaskSetMap["watermark_set"]; ok {
					for _, item := range v.([]interface{}) {
						watermarkSetMap := item.(map[string]interface{})
						watermarkInput := mps.WatermarkInput{}
						if v, ok := watermarkSetMap["definition"]; ok {
							watermarkInput.Definition = helper.IntUint64(v.(int))
						}
						if rawParameterMap, ok := helper.InterfaceToMap(watermarkSetMap, "raw_parameter"); ok {
							rawWatermarkParameter := mps.RawWatermarkParameter{}
							if v, ok := rawParameterMap["type"]; ok {
								rawWatermarkParameter.Type = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["coordinate_origin"]; ok {
								rawWatermarkParameter.CoordinateOrigin = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["x_pos"]; ok {
								rawWatermarkParameter.XPos = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["y_pos"]; ok {
								rawWatermarkParameter.YPos = helper.String(v.(string))
							}
							if imageTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "image_template"); ok {
								rawImageWatermarkInput := mps.RawImageWatermarkInput{}
								if imageContentMap, ok := helper.InterfaceToMap(imageTemplateMap, "image_content"); ok {
									mediaInputInfo := mps.MediaInputInfo{}
									if v, ok := imageContentMap["type"]; ok {
										mediaInputInfo.Type = helper.String(v.(string))
									}
									if cosInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "cos_input_info"); ok {
										cosInputInfo := mps.CosInputInfo{}
										if v, ok := cosInputInfoMap["bucket"]; ok {
											cosInputInfo.Bucket = helper.String(v.(string))
										}
										if v, ok := cosInputInfoMap["region"]; ok {
											cosInputInfo.Region = helper.String(v.(string))
										}
										if v, ok := cosInputInfoMap["object"]; ok {
											cosInputInfo.Object = helper.String(v.(string))
										}
										mediaInputInfo.CosInputInfo = &cosInputInfo
									}
									if urlInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "url_input_info"); ok {
										urlInputInfo := mps.UrlInputInfo{}
										if v, ok := urlInputInfoMap["url"]; ok {
											urlInputInfo.Url = helper.String(v.(string))
										}
										mediaInputInfo.UrlInputInfo = &urlInputInfo
									}
									rawImageWatermarkInput.ImageContent = &mediaInputInfo
								}
								if v, ok := imageTemplateMap["width"]; ok {
									rawImageWatermarkInput.Width = helper.String(v.(string))
								}
								if v, ok := imageTemplateMap["height"]; ok {
									rawImageWatermarkInput.Height = helper.String(v.(string))
								}
								if v, ok := imageTemplateMap["repeat_type"]; ok {
									rawImageWatermarkInput.RepeatType = helper.String(v.(string))
								}
								rawWatermarkParameter.ImageTemplate = &rawImageWatermarkInput
							}
							watermarkInput.RawParameter = &rawWatermarkParameter
						}
						if v, ok := watermarkSetMap["text_content"]; ok {
							watermarkInput.TextContent = helper.String(v.(string))
						}
						if v, ok := watermarkSetMap["svg_content"]; ok {
							watermarkInput.SvgContent = helper.String(v.(string))
						}
						if v, ok := watermarkSetMap["start_time_offset"]; ok {
							watermarkInput.StartTimeOffset = helper.Float64(v.(float64))
						}
						if v, ok := watermarkSetMap["end_time_offset"]; ok {
							watermarkInput.EndTimeOffset = helper.Float64(v.(float64))
						}
						sampleSnapshotTaskInput.WatermarkSet = append(sampleSnapshotTaskInput.WatermarkSet, &watermarkInput)
					}
				}
				if outputStorageMap, ok := helper.InterfaceToMap(sampleSnapshotTaskSetMap, "output_storage"); ok {
					taskOutputStorage := mps.TaskOutputStorage{}
					if v, ok := outputStorageMap["type"]; ok {
						taskOutputStorage.Type = helper.String(v.(string))
					}
					if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
						cosOutputStorage := mps.CosOutputStorage{}
						if v, ok := cosOutputStorageMap["bucket"]; ok {
							cosOutputStorage.Bucket = helper.String(v.(string))
						}
						if v, ok := cosOutputStorageMap["region"]; ok {
							cosOutputStorage.Region = helper.String(v.(string))
						}
						taskOutputStorage.CosOutputStorage = &cosOutputStorage
					}
					sampleSnapshotTaskInput.OutputStorage = &taskOutputStorage
				}
				if v, ok := sampleSnapshotTaskSetMap["output_object_path"]; ok {
					sampleSnapshotTaskInput.OutputObjectPath = helper.String(v.(string))
				}
				if objectNumberFormatMap, ok := helper.InterfaceToMap(sampleSnapshotTaskSetMap, "object_number_format"); ok {
					numberFormat := mps.NumberFormat{}
					if v, ok := objectNumberFormatMap["initial_value"]; ok {
						numberFormat.InitialValue = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["increment"]; ok {
						numberFormat.Increment = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["min_length"]; ok {
						numberFormat.MinLength = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["place_holder"]; ok {
						numberFormat.PlaceHolder = helper.String(v.(string))
					}
					sampleSnapshotTaskInput.ObjectNumberFormat = &numberFormat
				}
				mediaProcessTaskInput.SampleSnapshotTaskSet = append(mediaProcessTaskInput.SampleSnapshotTaskSet, &sampleSnapshotTaskInput)
			}
		}
		if v, ok := dMap["image_sprite_task_set"]; ok {
			for _, item := range v.([]interface{}) {
				imageSpriteTaskSetMap := item.(map[string]interface{})
				imageSpriteTaskInput := mps.ImageSpriteTaskInput{}
				if v, ok := imageSpriteTaskSetMap["definition"]; ok {
					imageSpriteTaskInput.Definition = helper.IntUint64(v.(int))
				}
				if outputStorageMap, ok := helper.InterfaceToMap(imageSpriteTaskSetMap, "output_storage"); ok {
					taskOutputStorage := mps.TaskOutputStorage{}
					if v, ok := outputStorageMap["type"]; ok {
						taskOutputStorage.Type = helper.String(v.(string))
					}
					if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
						cosOutputStorage := mps.CosOutputStorage{}
						if v, ok := cosOutputStorageMap["bucket"]; ok {
							cosOutputStorage.Bucket = helper.String(v.(string))
						}
						if v, ok := cosOutputStorageMap["region"]; ok {
							cosOutputStorage.Region = helper.String(v.(string))
						}
						taskOutputStorage.CosOutputStorage = &cosOutputStorage
					}
					imageSpriteTaskInput.OutputStorage = &taskOutputStorage
				}
				if v, ok := imageSpriteTaskSetMap["output_object_path"]; ok {
					imageSpriteTaskInput.OutputObjectPath = helper.String(v.(string))
				}
				if v, ok := imageSpriteTaskSetMap["web_vtt_object_name"]; ok {
					imageSpriteTaskInput.WebVttObjectName = helper.String(v.(string))
				}
				if objectNumberFormatMap, ok := helper.InterfaceToMap(imageSpriteTaskSetMap, "object_number_format"); ok {
					numberFormat := mps.NumberFormat{}
					if v, ok := objectNumberFormatMap["initial_value"]; ok {
						numberFormat.InitialValue = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["increment"]; ok {
						numberFormat.Increment = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["min_length"]; ok {
						numberFormat.MinLength = helper.IntUint64(v.(int))
					}
					if v, ok := objectNumberFormatMap["place_holder"]; ok {
						numberFormat.PlaceHolder = helper.String(v.(string))
					}
					imageSpriteTaskInput.ObjectNumberFormat = &numberFormat
				}
				mediaProcessTaskInput.ImageSpriteTaskSet = append(mediaProcessTaskInput.ImageSpriteTaskSet, &imageSpriteTaskInput)
			}
		}
		if v, ok := dMap["adaptive_dynamic_streaming_task_set"]; ok {
			for _, item := range v.([]interface{}) {
				adaptiveDynamicStreamingTaskSetMap := item.(map[string]interface{})
				adaptiveDynamicStreamingTaskInput := mps.AdaptiveDynamicStreamingTaskInput{}
				if v, ok := adaptiveDynamicStreamingTaskSetMap["definition"]; ok {
					adaptiveDynamicStreamingTaskInput.Definition = helper.IntUint64(v.(int))
				}
				if v, ok := adaptiveDynamicStreamingTaskSetMap["watermark_set"]; ok {
					for _, item := range v.([]interface{}) {
						watermarkSetMap := item.(map[string]interface{})
						watermarkInput := mps.WatermarkInput{}
						if v, ok := watermarkSetMap["definition"]; ok {
							watermarkInput.Definition = helper.IntUint64(v.(int))
						}
						if rawParameterMap, ok := helper.InterfaceToMap(watermarkSetMap, "raw_parameter"); ok {
							rawWatermarkParameter := mps.RawWatermarkParameter{}
							if v, ok := rawParameterMap["type"]; ok {
								rawWatermarkParameter.Type = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["coordinate_origin"]; ok {
								rawWatermarkParameter.CoordinateOrigin = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["x_pos"]; ok {
								rawWatermarkParameter.XPos = helper.String(v.(string))
							}
							if v, ok := rawParameterMap["y_pos"]; ok {
								rawWatermarkParameter.YPos = helper.String(v.(string))
							}
							if imageTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "image_template"); ok {
								rawImageWatermarkInput := mps.RawImageWatermarkInput{}
								if imageContentMap, ok := helper.InterfaceToMap(imageTemplateMap, "image_content"); ok {
									mediaInputInfo := mps.MediaInputInfo{}
									if v, ok := imageContentMap["type"]; ok {
										mediaInputInfo.Type = helper.String(v.(string))
									}
									if cosInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "cos_input_info"); ok {
										cosInputInfo := mps.CosInputInfo{}
										if v, ok := cosInputInfoMap["bucket"]; ok {
											cosInputInfo.Bucket = helper.String(v.(string))
										}
										if v, ok := cosInputInfoMap["region"]; ok {
											cosInputInfo.Region = helper.String(v.(string))
										}
										if v, ok := cosInputInfoMap["object"]; ok {
											cosInputInfo.Object = helper.String(v.(string))
										}
										mediaInputInfo.CosInputInfo = &cosInputInfo
									}
									if urlInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "url_input_info"); ok {
										urlInputInfo := mps.UrlInputInfo{}
										if v, ok := urlInputInfoMap["url"]; ok {
											urlInputInfo.Url = helper.String(v.(string))
										}
										mediaInputInfo.UrlInputInfo = &urlInputInfo
									}
									rawImageWatermarkInput.ImageContent = &mediaInputInfo
								}
								if v, ok := imageTemplateMap["width"]; ok {
									rawImageWatermarkInput.Width = helper.String(v.(string))
								}
								if v, ok := imageTemplateMap["height"]; ok {
									rawImageWatermarkInput.Height = helper.String(v.(string))
								}
								if v, ok := imageTemplateMap["repeat_type"]; ok {
									rawImageWatermarkInput.RepeatType = helper.String(v.(string))
								}
								rawWatermarkParameter.ImageTemplate = &rawImageWatermarkInput
							}
							watermarkInput.RawParameter = &rawWatermarkParameter
						}
						if v, ok := watermarkSetMap["text_content"]; ok {
							watermarkInput.TextContent = helper.String(v.(string))
						}
						if v, ok := watermarkSetMap["svg_content"]; ok {
							watermarkInput.SvgContent = helper.String(v.(string))
						}
						if v, ok := watermarkSetMap["start_time_offset"]; ok {
							watermarkInput.StartTimeOffset = helper.Float64(v.(float64))
						}
						if v, ok := watermarkSetMap["end_time_offset"]; ok {
							watermarkInput.EndTimeOffset = helper.Float64(v.(float64))
						}
						adaptiveDynamicStreamingTaskInput.WatermarkSet = append(adaptiveDynamicStreamingTaskInput.WatermarkSet, &watermarkInput)
					}
				}
				if outputStorageMap, ok := helper.InterfaceToMap(adaptiveDynamicStreamingTaskSetMap, "output_storage"); ok {
					taskOutputStorage := mps.TaskOutputStorage{}
					if v, ok := outputStorageMap["type"]; ok {
						taskOutputStorage.Type = helper.String(v.(string))
					}
					if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
						cosOutputStorage := mps.CosOutputStorage{}
						if v, ok := cosOutputStorageMap["bucket"]; ok {
							cosOutputStorage.Bucket = helper.String(v.(string))
						}
						if v, ok := cosOutputStorageMap["region"]; ok {
							cosOutputStorage.Region = helper.String(v.(string))
						}
						taskOutputStorage.CosOutputStorage = &cosOutputStorage
					}
					adaptiveDynamicStreamingTaskInput.OutputStorage = &taskOutputStorage
				}
				if v, ok := adaptiveDynamicStreamingTaskSetMap["output_object_path"]; ok {
					adaptiveDynamicStreamingTaskInput.OutputObjectPath = helper.String(v.(string))
				}
				if v, ok := adaptiveDynamicStreamingTaskSetMap["sub_stream_object_name"]; ok {
					adaptiveDynamicStreamingTaskInput.SubStreamObjectName = helper.String(v.(string))
				}
				if v, ok := adaptiveDynamicStreamingTaskSetMap["segment_object_name"]; ok {
					adaptiveDynamicStreamingTaskInput.SegmentObjectName = helper.String(v.(string))
				}
				mediaProcessTaskInput.AdaptiveDynamicStreamingTaskSet = append(mediaProcessTaskInput.AdaptiveDynamicStreamingTaskSet, &adaptiveDynamicStreamingTaskInput)
			}
		}
		request.MediaProcessTask = &mediaProcessTaskInput
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "ai_content_review_task"); ok {
		aiContentReviewTaskInput := mps.AiContentReviewTaskInput{}
		if v, ok := dMap["definition"]; ok {
			aiContentReviewTaskInput.Definition = helper.IntUint64(v.(int))
		}
		request.AiContentReviewTask = &aiContentReviewTaskInput
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "ai_analysis_task"); ok {
		aiAnalysisTaskInput := mps.AiAnalysisTaskInput{}
		if v, ok := dMap["definition"]; ok {
			aiAnalysisTaskInput.Definition = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["extended_parameter"]; ok {
			aiAnalysisTaskInput.ExtendedParameter = helper.String(v.(string))
		}
		request.AiAnalysisTask = &aiAnalysisTaskInput
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "ai_recognition_task"); ok {
		aiRecognitionTaskInput := mps.AiRecognitionTaskInput{}
		if v, ok := dMap["definition"]; ok {
			aiRecognitionTaskInput.Definition = helper.IntUint64(v.(int))
		}
		request.AiRecognitionTask = &aiRecognitionTaskInput
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "task_notify_config"); ok {
		taskNotifyConfig := mps.TaskNotifyConfig{}
		if v, ok := dMap["cmq_model"]; ok {
			taskNotifyConfig.CmqModel = helper.String(v.(string))
		}
		if v, ok := dMap["cmq_region"]; ok {
			taskNotifyConfig.CmqRegion = helper.String(v.(string))
		}
		if v, ok := dMap["topic_name"]; ok {
			taskNotifyConfig.TopicName = helper.String(v.(string))
		}
		if v, ok := dMap["queue_name"]; ok {
			taskNotifyConfig.QueueName = helper.String(v.(string))
		}
		if v, ok := dMap["notify_mode"]; ok {
			taskNotifyConfig.NotifyMode = helper.String(v.(string))
		}
		if v, ok := dMap["notify_type"]; ok {
			taskNotifyConfig.NotifyType = helper.String(v.(string))
		}
		if v, ok := dMap["notify_url"]; ok {
			taskNotifyConfig.NotifyUrl = helper.String(v.(string))
		}
		request.TaskNotifyConfig = &taskNotifyConfig
	}

	if v, _ := d.GetOk("task_priority"); v != nil {
		request.TaskPriority = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().CreateWorkflow(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps workflow failed, reason:%+v", logId, err)
		return err
	}

	workflowId = *response.Response.WorkflowId
	d.SetId(helper.Int64ToStr(workflowId))

	return resourceTencentCloudMpsWorkflowRead(d, meta)
}

func resourceTencentCloudMpsWorkflowRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_workflow.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	workflowId := d.Id()

	workflow, err := service.DescribeMpsWorkflowById(ctx, workflowId)
	if err != nil {
		return err
	}

	if workflow == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsWorkflow` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if workflow.WorkflowName != nil {
		_ = d.Set("workflow_name", workflow.WorkflowName)
	}

	if workflow.Trigger != nil {
		triggerMap := map[string]interface{}{}

		if workflow.Trigger.Type != nil {
			triggerMap["type"] = workflow.Trigger.Type
		}

		if workflow.Trigger.CosFileUploadTrigger != nil {
			cosFileUploadTriggerMap := map[string]interface{}{}

			if workflow.Trigger.CosFileUploadTrigger.Bucket != nil {
				cosFileUploadTriggerMap["bucket"] = workflow.Trigger.CosFileUploadTrigger.Bucket
			}

			if workflow.Trigger.CosFileUploadTrigger.Region != nil {
				cosFileUploadTriggerMap["region"] = workflow.Trigger.CosFileUploadTrigger.Region
			}

			if workflow.Trigger.CosFileUploadTrigger.Dir != nil {
				cosFileUploadTriggerMap["dir"] = workflow.Trigger.CosFileUploadTrigger.Dir
			}

			if workflow.Trigger.CosFileUploadTrigger.Formats != nil {
				cosFileUploadTriggerMap["formats"] = workflow.Trigger.CosFileUploadTrigger.Formats
			}

			triggerMap["cos_file_upload_trigger"] = []interface{}{cosFileUploadTriggerMap}
		}

		_ = d.Set("trigger", []interface{}{triggerMap})
	}

	if workflow.OutputStorage != nil {
		outputStorageMap := map[string]interface{}{}

		if workflow.OutputStorage.Type != nil {
			outputStorageMap["type"] = workflow.OutputStorage.Type
		}

		if workflow.OutputStorage.CosOutputStorage != nil {
			cosOutputStorageMap := map[string]interface{}{}

			if workflow.OutputStorage.CosOutputStorage.Bucket != nil {
				cosOutputStorageMap["bucket"] = workflow.OutputStorage.CosOutputStorage.Bucket
			}

			if workflow.OutputStorage.CosOutputStorage.Region != nil {
				cosOutputStorageMap["region"] = workflow.OutputStorage.CosOutputStorage.Region
			}

			outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
		}

		_ = d.Set("output_storage", []interface{}{outputStorageMap})
	}

	if workflow.OutputDir != nil {
		_ = d.Set("output_dir", workflow.OutputDir)
	}

	if workflow.MediaProcessTask != nil {
		mediaProcessTaskMap := map[string]interface{}{}

		if workflow.MediaProcessTask.TranscodeTaskSet != nil {
			transcodeTaskSetList := []interface{}{}
			for _, transcodeTaskSet := range workflow.MediaProcessTask.TranscodeTaskSet {
				transcodeTaskSetMap := map[string]interface{}{}

				if transcodeTaskSet.Definition != nil {
					transcodeTaskSetMap["definition"] = transcodeTaskSet.Definition
				}

				if transcodeTaskSet.RawParameter != nil {
					rawParameterMap := map[string]interface{}{}

					if transcodeTaskSet.RawParameter.Container != nil {
						rawParameterMap["container"] = transcodeTaskSet.RawParameter.Container
					}

					if transcodeTaskSet.RawParameter.RemoveVideo != nil {
						rawParameterMap["remove_video"] = transcodeTaskSet.RawParameter.RemoveVideo
					}

					if transcodeTaskSet.RawParameter.RemoveAudio != nil {
						rawParameterMap["remove_audio"] = transcodeTaskSet.RawParameter.RemoveAudio
					}

					if transcodeTaskSet.RawParameter.VideoTemplate != nil {
						videoTemplateMap := map[string]interface{}{}

						if transcodeTaskSet.RawParameter.VideoTemplate.Codec != nil {
							videoTemplateMap["codec"] = transcodeTaskSet.RawParameter.VideoTemplate.Codec
						}

						if transcodeTaskSet.RawParameter.VideoTemplate.Fps != nil {
							videoTemplateMap["fps"] = transcodeTaskSet.RawParameter.VideoTemplate.Fps
						}

						if transcodeTaskSet.RawParameter.VideoTemplate.Bitrate != nil {
							videoTemplateMap["bitrate"] = transcodeTaskSet.RawParameter.VideoTemplate.Bitrate
						}

						if transcodeTaskSet.RawParameter.VideoTemplate.ResolutionAdaptive != nil {
							videoTemplateMap["resolution_adaptive"] = transcodeTaskSet.RawParameter.VideoTemplate.ResolutionAdaptive
						}

						if transcodeTaskSet.RawParameter.VideoTemplate.Width != nil {
							videoTemplateMap["width"] = transcodeTaskSet.RawParameter.VideoTemplate.Width
						}

						if transcodeTaskSet.RawParameter.VideoTemplate.Height != nil {
							videoTemplateMap["height"] = transcodeTaskSet.RawParameter.VideoTemplate.Height
						}

						if transcodeTaskSet.RawParameter.VideoTemplate.Gop != nil {
							videoTemplateMap["gop"] = transcodeTaskSet.RawParameter.VideoTemplate.Gop
						}

						if transcodeTaskSet.RawParameter.VideoTemplate.FillType != nil {
							videoTemplateMap["fill_type"] = transcodeTaskSet.RawParameter.VideoTemplate.FillType
						}

						if transcodeTaskSet.RawParameter.VideoTemplate.Vcrf != nil {
							videoTemplateMap["vcrf"] = transcodeTaskSet.RawParameter.VideoTemplate.Vcrf
						}

						rawParameterMap["video_template"] = []interface{}{videoTemplateMap}
					}

					if transcodeTaskSet.RawParameter.AudioTemplate != nil {
						audioTemplateMap := map[string]interface{}{}

						if transcodeTaskSet.RawParameter.AudioTemplate.Codec != nil {
							audioTemplateMap["codec"] = transcodeTaskSet.RawParameter.AudioTemplate.Codec
						}

						if transcodeTaskSet.RawParameter.AudioTemplate.Bitrate != nil {
							audioTemplateMap["bitrate"] = transcodeTaskSet.RawParameter.AudioTemplate.Bitrate
						}

						if transcodeTaskSet.RawParameter.AudioTemplate.SampleRate != nil {
							audioTemplateMap["sample_rate"] = transcodeTaskSet.RawParameter.AudioTemplate.SampleRate
						}

						if transcodeTaskSet.RawParameter.AudioTemplate.AudioChannel != nil {
							audioTemplateMap["audio_channel"] = transcodeTaskSet.RawParameter.AudioTemplate.AudioChannel
						}

						rawParameterMap["audio_template"] = []interface{}{audioTemplateMap}
					}

					if transcodeTaskSet.RawParameter.TEHDConfig != nil {
						tEHDConfigMap := map[string]interface{}{}

						if transcodeTaskSet.RawParameter.TEHDConfig.Type != nil {
							tEHDConfigMap["type"] = transcodeTaskSet.RawParameter.TEHDConfig.Type
						}

						if transcodeTaskSet.RawParameter.TEHDConfig.MaxVideoBitrate != nil {
							tEHDConfigMap["max_video_bitrate"] = transcodeTaskSet.RawParameter.TEHDConfig.MaxVideoBitrate
						}

						rawParameterMap["tehd_config"] = []interface{}{tEHDConfigMap}
					}

					transcodeTaskSetMap["raw_parameter"] = []interface{}{rawParameterMap}
				}

				if transcodeTaskSet.OverrideParameter != nil {
					overrideParameterMap := map[string]interface{}{}

					if transcodeTaskSet.OverrideParameter.Container != nil {
						overrideParameterMap["container"] = transcodeTaskSet.OverrideParameter.Container
					}

					if transcodeTaskSet.OverrideParameter.RemoveVideo != nil {
						overrideParameterMap["remove_video"] = transcodeTaskSet.OverrideParameter.RemoveVideo
					}

					if transcodeTaskSet.OverrideParameter.RemoveAudio != nil {
						overrideParameterMap["remove_audio"] = transcodeTaskSet.OverrideParameter.RemoveAudio
					}

					if transcodeTaskSet.OverrideParameter.VideoTemplate != nil {
						videoTemplateMap := map[string]interface{}{}

						if transcodeTaskSet.OverrideParameter.VideoTemplate.Codec != nil {
							videoTemplateMap["codec"] = transcodeTaskSet.OverrideParameter.VideoTemplate.Codec
						}

						if transcodeTaskSet.OverrideParameter.VideoTemplate.Fps != nil {
							videoTemplateMap["fps"] = transcodeTaskSet.OverrideParameter.VideoTemplate.Fps
						}

						if transcodeTaskSet.OverrideParameter.VideoTemplate.Bitrate != nil {
							videoTemplateMap["bitrate"] = transcodeTaskSet.OverrideParameter.VideoTemplate.Bitrate
						}

						if transcodeTaskSet.OverrideParameter.VideoTemplate.ResolutionAdaptive != nil {
							videoTemplateMap["resolution_adaptive"] = transcodeTaskSet.OverrideParameter.VideoTemplate.ResolutionAdaptive
						}

						if transcodeTaskSet.OverrideParameter.VideoTemplate.Width != nil {
							videoTemplateMap["width"] = transcodeTaskSet.OverrideParameter.VideoTemplate.Width
						}

						if transcodeTaskSet.OverrideParameter.VideoTemplate.Height != nil {
							videoTemplateMap["height"] = transcodeTaskSet.OverrideParameter.VideoTemplate.Height
						}

						if transcodeTaskSet.OverrideParameter.VideoTemplate.Gop != nil {
							videoTemplateMap["gop"] = transcodeTaskSet.OverrideParameter.VideoTemplate.Gop
						}

						if transcodeTaskSet.OverrideParameter.VideoTemplate.FillType != nil {
							videoTemplateMap["fill_type"] = transcodeTaskSet.OverrideParameter.VideoTemplate.FillType
						}

						if transcodeTaskSet.OverrideParameter.VideoTemplate.Vcrf != nil {
							videoTemplateMap["vcrf"] = transcodeTaskSet.OverrideParameter.VideoTemplate.Vcrf
						}

						if transcodeTaskSet.OverrideParameter.VideoTemplate.ContentAdaptStream != nil {
							videoTemplateMap["content_adapt_stream"] = transcodeTaskSet.OverrideParameter.VideoTemplate.ContentAdaptStream
						}

						overrideParameterMap["video_template"] = []interface{}{videoTemplateMap}
					}

					if transcodeTaskSet.OverrideParameter.AudioTemplate != nil {
						audioTemplateMap := map[string]interface{}{}

						if transcodeTaskSet.OverrideParameter.AudioTemplate.Codec != nil {
							audioTemplateMap["codec"] = transcodeTaskSet.OverrideParameter.AudioTemplate.Codec
						}

						if transcodeTaskSet.OverrideParameter.AudioTemplate.Bitrate != nil {
							audioTemplateMap["bitrate"] = transcodeTaskSet.OverrideParameter.AudioTemplate.Bitrate
						}

						if transcodeTaskSet.OverrideParameter.AudioTemplate.SampleRate != nil {
							audioTemplateMap["sample_rate"] = transcodeTaskSet.OverrideParameter.AudioTemplate.SampleRate
						}

						if transcodeTaskSet.OverrideParameter.AudioTemplate.AudioChannel != nil {
							audioTemplateMap["audio_channel"] = transcodeTaskSet.OverrideParameter.AudioTemplate.AudioChannel
						}

						if transcodeTaskSet.OverrideParameter.AudioTemplate.StreamSelects != nil {
							audioTemplateMap["stream_selects"] = transcodeTaskSet.OverrideParameter.AudioTemplate.StreamSelects
						}

						overrideParameterMap["audio_template"] = []interface{}{audioTemplateMap}
					}

					if transcodeTaskSet.OverrideParameter.TEHDConfig != nil {
						tEHDConfigMap := map[string]interface{}{}

						if transcodeTaskSet.OverrideParameter.TEHDConfig.Type != nil {
							tEHDConfigMap["type"] = transcodeTaskSet.OverrideParameter.TEHDConfig.Type
						}

						if transcodeTaskSet.OverrideParameter.TEHDConfig.MaxVideoBitrate != nil {
							tEHDConfigMap["max_video_bitrate"] = transcodeTaskSet.OverrideParameter.TEHDConfig.MaxVideoBitrate
						}

						overrideParameterMap["tehd_config"] = []interface{}{tEHDConfigMap}
					}

					if transcodeTaskSet.OverrideParameter.SubtitleTemplate != nil {
						subtitleTemplateMap := map[string]interface{}{}

						if transcodeTaskSet.OverrideParameter.SubtitleTemplate.Path != nil {
							subtitleTemplateMap["path"] = transcodeTaskSet.OverrideParameter.SubtitleTemplate.Path
						}

						if transcodeTaskSet.OverrideParameter.SubtitleTemplate.StreamIndex != nil {
							subtitleTemplateMap["stream_index"] = transcodeTaskSet.OverrideParameter.SubtitleTemplate.StreamIndex
						}

						if transcodeTaskSet.OverrideParameter.SubtitleTemplate.FontType != nil {
							subtitleTemplateMap["font_type"] = transcodeTaskSet.OverrideParameter.SubtitleTemplate.FontType
						}

						if transcodeTaskSet.OverrideParameter.SubtitleTemplate.FontSize != nil {
							subtitleTemplateMap["font_size"] = transcodeTaskSet.OverrideParameter.SubtitleTemplate.FontSize
						}

						if transcodeTaskSet.OverrideParameter.SubtitleTemplate.FontColor != nil {
							subtitleTemplateMap["font_color"] = transcodeTaskSet.OverrideParameter.SubtitleTemplate.FontColor
						}

						if transcodeTaskSet.OverrideParameter.SubtitleTemplate.FontAlpha != nil {
							subtitleTemplateMap["font_alpha"] = transcodeTaskSet.OverrideParameter.SubtitleTemplate.FontAlpha
						}

						overrideParameterMap["subtitle_template"] = []interface{}{subtitleTemplateMap}
					}

					transcodeTaskSetMap["override_parameter"] = []interface{}{overrideParameterMap}
				}

				if transcodeTaskSet.WatermarkSet != nil {
					watermarkSetList := []interface{}{}
					for _, watermarkSet := range transcodeTaskSet.WatermarkSet {
						watermarkSetMap := map[string]interface{}{}

						if watermarkSet.Definition != nil {
							watermarkSetMap["definition"] = watermarkSet.Definition
						}

						if watermarkSet.RawParameter != nil {
							rawParameterMap := map[string]interface{}{}

							if watermarkSet.RawParameter.Type != nil {
								rawParameterMap["type"] = watermarkSet.RawParameter.Type
							}

							if watermarkSet.RawParameter.CoordinateOrigin != nil {
								rawParameterMap["coordinate_origin"] = watermarkSet.RawParameter.CoordinateOrigin
							}

							if watermarkSet.RawParameter.XPos != nil {
								rawParameterMap["x_pos"] = watermarkSet.RawParameter.XPos
							}

							if watermarkSet.RawParameter.YPos != nil {
								rawParameterMap["y_pos"] = watermarkSet.RawParameter.YPos
							}

							if watermarkSet.RawParameter.ImageTemplate != nil {
								imageTemplateMap := map[string]interface{}{}

								if watermarkSet.RawParameter.ImageTemplate.ImageContent != nil {
									imageContentMap := map[string]interface{}{}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.Type != nil {
										imageContentMap["type"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.Type
									}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo != nil {
										cosInputInfoMap := map[string]interface{}{}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Bucket != nil {
											cosInputInfoMap["bucket"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Bucket
										}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Region != nil {
											cosInputInfoMap["region"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Region
										}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Object != nil {
											cosInputInfoMap["object"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Object
										}

										imageContentMap["cos_input_info"] = []interface{}{cosInputInfoMap}
									}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo != nil {
										urlInputInfoMap := map[string]interface{}{}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo.Url != nil {
											urlInputInfoMap["url"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo.Url
										}

										imageContentMap["url_input_info"] = []interface{}{urlInputInfoMap}
									}

									imageTemplateMap["image_content"] = []interface{}{imageContentMap}
								}

								if watermarkSet.RawParameter.ImageTemplate.Width != nil {
									imageTemplateMap["width"] = watermarkSet.RawParameter.ImageTemplate.Width
								}

								if watermarkSet.RawParameter.ImageTemplate.Height != nil {
									imageTemplateMap["height"] = watermarkSet.RawParameter.ImageTemplate.Height
								}

								if watermarkSet.RawParameter.ImageTemplate.RepeatType != nil {
									imageTemplateMap["repeat_type"] = watermarkSet.RawParameter.ImageTemplate.RepeatType
								}

								rawParameterMap["image_template"] = []interface{}{imageTemplateMap}
							}

							watermarkSetMap["raw_parameter"] = []interface{}{rawParameterMap}
						}

						if watermarkSet.TextContent != nil {
							watermarkSetMap["text_content"] = watermarkSet.TextContent
						}

						if watermarkSet.SvgContent != nil {
							watermarkSetMap["svg_content"] = watermarkSet.SvgContent
						}

						if watermarkSet.StartTimeOffset != nil {
							watermarkSetMap["start_time_offset"] = watermarkSet.StartTimeOffset
						}

						if watermarkSet.EndTimeOffset != nil {
							watermarkSetMap["end_time_offset"] = watermarkSet.EndTimeOffset
						}

						watermarkSetList = append(watermarkSetList, watermarkSetMap)
					}

					transcodeTaskSetMap["watermark_set"] = watermarkSetList
				}

				if transcodeTaskSet.MosaicSet != nil {
					mosaicSetList := []interface{}{}
					for _, mosaicSet := range transcodeTaskSet.MosaicSet {
						mosaicSetMap := map[string]interface{}{}

						if mosaicSet.CoordinateOrigin != nil {
							mosaicSetMap["coordinate_origin"] = mosaicSet.CoordinateOrigin
						}

						if mosaicSet.XPos != nil {
							mosaicSetMap["x_pos"] = mosaicSet.XPos
						}

						if mosaicSet.YPos != nil {
							mosaicSetMap["y_pos"] = mosaicSet.YPos
						}

						if mosaicSet.Width != nil {
							mosaicSetMap["width"] = mosaicSet.Width
						}

						if mosaicSet.Height != nil {
							mosaicSetMap["height"] = mosaicSet.Height
						}

						if mosaicSet.StartTimeOffset != nil {
							mosaicSetMap["start_time_offset"] = mosaicSet.StartTimeOffset
						}

						if mosaicSet.EndTimeOffset != nil {
							mosaicSetMap["end_time_offset"] = mosaicSet.EndTimeOffset
						}

						mosaicSetList = append(mosaicSetList, mosaicSetMap)
					}

					transcodeTaskSetMap["mosaic_set"] = mosaicSetList
				}

				if transcodeTaskSet.StartTimeOffset != nil {
					transcodeTaskSetMap["start_time_offset"] = transcodeTaskSet.StartTimeOffset
				}

				if transcodeTaskSet.EndTimeOffset != nil {
					transcodeTaskSetMap["end_time_offset"] = transcodeTaskSet.EndTimeOffset
				}

				if transcodeTaskSet.OutputStorage != nil {
					outputStorageMap := map[string]interface{}{}

					if transcodeTaskSet.OutputStorage.Type != nil {
						outputStorageMap["type"] = transcodeTaskSet.OutputStorage.Type
					}

					if transcodeTaskSet.OutputStorage.CosOutputStorage != nil {
						cosOutputStorageMap := map[string]interface{}{}

						if transcodeTaskSet.OutputStorage.CosOutputStorage.Bucket != nil {
							cosOutputStorageMap["bucket"] = transcodeTaskSet.OutputStorage.CosOutputStorage.Bucket
						}

						if transcodeTaskSet.OutputStorage.CosOutputStorage.Region != nil {
							cosOutputStorageMap["region"] = transcodeTaskSet.OutputStorage.CosOutputStorage.Region
						}

						outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
					}

					transcodeTaskSetMap["output_storage"] = []interface{}{outputStorageMap}
				}

				if transcodeTaskSet.OutputObjectPath != nil {
					transcodeTaskSetMap["output_object_path"] = transcodeTaskSet.OutputObjectPath
				}

				if transcodeTaskSet.SegmentObjectName != nil {
					transcodeTaskSetMap["segment_object_name"] = transcodeTaskSet.SegmentObjectName
				}

				if transcodeTaskSet.ObjectNumberFormat != nil {
					objectNumberFormatMap := map[string]interface{}{}

					if transcodeTaskSet.ObjectNumberFormat.InitialValue != nil {
						objectNumberFormatMap["initial_value"] = transcodeTaskSet.ObjectNumberFormat.InitialValue
					}

					if transcodeTaskSet.ObjectNumberFormat.Increment != nil {
						objectNumberFormatMap["increment"] = transcodeTaskSet.ObjectNumberFormat.Increment
					}

					if transcodeTaskSet.ObjectNumberFormat.MinLength != nil {
						objectNumberFormatMap["min_length"] = transcodeTaskSet.ObjectNumberFormat.MinLength
					}

					if transcodeTaskSet.ObjectNumberFormat.PlaceHolder != nil {
						objectNumberFormatMap["place_holder"] = transcodeTaskSet.ObjectNumberFormat.PlaceHolder
					}

					transcodeTaskSetMap["object_number_format"] = []interface{}{objectNumberFormatMap}
				}

				if transcodeTaskSet.HeadTailParameter != nil {
					headTailParameterMap := map[string]interface{}{}

					if transcodeTaskSet.HeadTailParameter.HeadSet != nil {
						headSetList := []interface{}{}
						for _, headSet := range transcodeTaskSet.HeadTailParameter.HeadSet {
							headSetMap := map[string]interface{}{}

							if headSet.Type != nil {
								headSetMap["type"] = headSet.Type
							}

							if headSet.CosInputInfo != nil {
								cosInputInfoMap := map[string]interface{}{}

								if headSet.CosInputInfo.Bucket != nil {
									cosInputInfoMap["bucket"] = headSet.CosInputInfo.Bucket
								}

								if headSet.CosInputInfo.Region != nil {
									cosInputInfoMap["region"] = headSet.CosInputInfo.Region
								}

								if headSet.CosInputInfo.Object != nil {
									cosInputInfoMap["object"] = headSet.CosInputInfo.Object
								}

								headSetMap["cos_input_info"] = []interface{}{cosInputInfoMap}
							}

							if headSet.UrlInputInfo != nil {
								urlInputInfoMap := map[string]interface{}{}

								if headSet.UrlInputInfo.Url != nil {
									urlInputInfoMap["url"] = headSet.UrlInputInfo.Url
								}

								headSetMap["url_input_info"] = []interface{}{urlInputInfoMap}
							}

							headSetList = append(headSetList, headSetMap)
						}

						headTailParameterMap["head_set"] = headSetList
					}

					if transcodeTaskSet.HeadTailParameter.TailSet != nil {
						tailSetList := []interface{}{}
						for _, tailSet := range transcodeTaskSet.HeadTailParameter.TailSet {
							tailSetMap := map[string]interface{}{}

							if tailSet.Type != nil {
								tailSetMap["type"] = tailSet.Type
							}

							if tailSet.CosInputInfo != nil {
								cosInputInfoMap := map[string]interface{}{}

								if tailSet.CosInputInfo.Bucket != nil {
									cosInputInfoMap["bucket"] = tailSet.CosInputInfo.Bucket
								}

								if tailSet.CosInputInfo.Region != nil {
									cosInputInfoMap["region"] = tailSet.CosInputInfo.Region
								}

								if tailSet.CosInputInfo.Object != nil {
									cosInputInfoMap["object"] = tailSet.CosInputInfo.Object
								}

								tailSetMap["cos_input_info"] = []interface{}{cosInputInfoMap}
							}

							if tailSet.UrlInputInfo != nil {
								urlInputInfoMap := map[string]interface{}{}

								if tailSet.UrlInputInfo.Url != nil {
									urlInputInfoMap["url"] = tailSet.UrlInputInfo.Url
								}

								tailSetMap["url_input_info"] = []interface{}{urlInputInfoMap}
							}

							tailSetList = append(tailSetList, tailSetMap)
						}

						headTailParameterMap["tail_set"] = tailSetList
					}

					transcodeTaskSetMap["head_tail_parameter"] = []interface{}{headTailParameterMap}
				}

				transcodeTaskSetList = append(transcodeTaskSetList, transcodeTaskSetMap)
			}

			mediaProcessTaskMap["transcode_task_set"] = transcodeTaskSetList
		}

		if workflow.MediaProcessTask.AnimatedGraphicTaskSet != nil {
			animatedGraphicTaskSetList := []interface{}{}
			for _, animatedGraphicTaskSet := range workflow.MediaProcessTask.AnimatedGraphicTaskSet {
				animatedGraphicTaskSetMap := map[string]interface{}{}

				if animatedGraphicTaskSet.Definition != nil {
					animatedGraphicTaskSetMap["definition"] = animatedGraphicTaskSet.Definition
				}

				if animatedGraphicTaskSet.StartTimeOffset != nil {
					animatedGraphicTaskSetMap["start_time_offset"] = animatedGraphicTaskSet.StartTimeOffset
				}

				if animatedGraphicTaskSet.EndTimeOffset != nil {
					animatedGraphicTaskSetMap["end_time_offset"] = animatedGraphicTaskSet.EndTimeOffset
				}

				if animatedGraphicTaskSet.OutputStorage != nil {
					outputStorageMap := map[string]interface{}{}

					if animatedGraphicTaskSet.OutputStorage.Type != nil {
						outputStorageMap["type"] = animatedGraphicTaskSet.OutputStorage.Type
					}

					if animatedGraphicTaskSet.OutputStorage.CosOutputStorage != nil {
						cosOutputStorageMap := map[string]interface{}{}

						if animatedGraphicTaskSet.OutputStorage.CosOutputStorage.Bucket != nil {
							cosOutputStorageMap["bucket"] = animatedGraphicTaskSet.OutputStorage.CosOutputStorage.Bucket
						}

						if animatedGraphicTaskSet.OutputStorage.CosOutputStorage.Region != nil {
							cosOutputStorageMap["region"] = animatedGraphicTaskSet.OutputStorage.CosOutputStorage.Region
						}

						outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
					}

					animatedGraphicTaskSetMap["output_storage"] = []interface{}{outputStorageMap}
				}

				if animatedGraphicTaskSet.OutputObjectPath != nil {
					animatedGraphicTaskSetMap["output_object_path"] = animatedGraphicTaskSet.OutputObjectPath
				}

				animatedGraphicTaskSetList = append(animatedGraphicTaskSetList, animatedGraphicTaskSetMap)
			}

			mediaProcessTaskMap["animated_graphic_task_set"] = animatedGraphicTaskSetList
		}

		if workflow.MediaProcessTask.SnapshotByTimeOffsetTaskSet != nil {
			snapshotByTimeOffsetTaskSetList := []interface{}{}
			for _, snapshotByTimeOffsetTaskSet := range workflow.MediaProcessTask.SnapshotByTimeOffsetTaskSet {
				snapshotByTimeOffsetTaskSetMap := map[string]interface{}{}

				if snapshotByTimeOffsetTaskSet.Definition != nil {
					snapshotByTimeOffsetTaskSetMap["definition"] = snapshotByTimeOffsetTaskSet.Definition
				}

				if snapshotByTimeOffsetTaskSet.ExtTimeOffsetSet != nil {
					snapshotByTimeOffsetTaskSetMap["ext_time_offset_set"] = snapshotByTimeOffsetTaskSet.ExtTimeOffsetSet
				}

				if snapshotByTimeOffsetTaskSet.TimeOffsetSet != nil {
					snapshotByTimeOffsetTaskSetMap["time_offset_set"] = snapshotByTimeOffsetTaskSet.TimeOffsetSet
				}

				if snapshotByTimeOffsetTaskSet.WatermarkSet != nil {
					watermarkSetList := []interface{}{}
					for _, watermarkSet := range snapshotByTimeOffsetTaskSet.WatermarkSet {
						watermarkSetMap := map[string]interface{}{}

						if watermarkSet.Definition != nil {
							watermarkSetMap["definition"] = watermarkSet.Definition
						}

						if watermarkSet.RawParameter != nil {
							rawParameterMap := map[string]interface{}{}

							if watermarkSet.RawParameter.Type != nil {
								rawParameterMap["type"] = watermarkSet.RawParameter.Type
							}

							if watermarkSet.RawParameter.CoordinateOrigin != nil {
								rawParameterMap["coordinate_origin"] = watermarkSet.RawParameter.CoordinateOrigin
							}

							if watermarkSet.RawParameter.XPos != nil {
								rawParameterMap["x_pos"] = watermarkSet.RawParameter.XPos
							}

							if watermarkSet.RawParameter.YPos != nil {
								rawParameterMap["y_pos"] = watermarkSet.RawParameter.YPos
							}

							if watermarkSet.RawParameter.ImageTemplate != nil {
								imageTemplateMap := map[string]interface{}{}

								if watermarkSet.RawParameter.ImageTemplate.ImageContent != nil {
									imageContentMap := map[string]interface{}{}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.Type != nil {
										imageContentMap["type"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.Type
									}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo != nil {
										cosInputInfoMap := map[string]interface{}{}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Bucket != nil {
											cosInputInfoMap["bucket"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Bucket
										}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Region != nil {
											cosInputInfoMap["region"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Region
										}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Object != nil {
											cosInputInfoMap["object"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Object
										}

										imageContentMap["cos_input_info"] = []interface{}{cosInputInfoMap}
									}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo != nil {
										urlInputInfoMap := map[string]interface{}{}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo.Url != nil {
											urlInputInfoMap["url"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo.Url
										}

										imageContentMap["url_input_info"] = []interface{}{urlInputInfoMap}
									}

									imageTemplateMap["image_content"] = []interface{}{imageContentMap}
								}

								if watermarkSet.RawParameter.ImageTemplate.Width != nil {
									imageTemplateMap["width"] = watermarkSet.RawParameter.ImageTemplate.Width
								}

								if watermarkSet.RawParameter.ImageTemplate.Height != nil {
									imageTemplateMap["height"] = watermarkSet.RawParameter.ImageTemplate.Height
								}

								if watermarkSet.RawParameter.ImageTemplate.RepeatType != nil {
									imageTemplateMap["repeat_type"] = watermarkSet.RawParameter.ImageTemplate.RepeatType
								}

								rawParameterMap["image_template"] = []interface{}{imageTemplateMap}
							}

							watermarkSetMap["raw_parameter"] = []interface{}{rawParameterMap}
						}

						if watermarkSet.TextContent != nil {
							watermarkSetMap["text_content"] = watermarkSet.TextContent
						}

						if watermarkSet.SvgContent != nil {
							watermarkSetMap["svg_content"] = watermarkSet.SvgContent
						}

						if watermarkSet.StartTimeOffset != nil {
							watermarkSetMap["start_time_offset"] = watermarkSet.StartTimeOffset
						}

						if watermarkSet.EndTimeOffset != nil {
							watermarkSetMap["end_time_offset"] = watermarkSet.EndTimeOffset
						}

						watermarkSetList = append(watermarkSetList, watermarkSetMap)
					}

					snapshotByTimeOffsetTaskSetMap["watermark_set"] = watermarkSetList
				}

				if snapshotByTimeOffsetTaskSet.OutputStorage != nil {
					outputStorageMap := map[string]interface{}{}

					if snapshotByTimeOffsetTaskSet.OutputStorage.Type != nil {
						outputStorageMap["type"] = snapshotByTimeOffsetTaskSet.OutputStorage.Type
					}

					if snapshotByTimeOffsetTaskSet.OutputStorage.CosOutputStorage != nil {
						cosOutputStorageMap := map[string]interface{}{}

						if snapshotByTimeOffsetTaskSet.OutputStorage.CosOutputStorage.Bucket != nil {
							cosOutputStorageMap["bucket"] = snapshotByTimeOffsetTaskSet.OutputStorage.CosOutputStorage.Bucket
						}

						if snapshotByTimeOffsetTaskSet.OutputStorage.CosOutputStorage.Region != nil {
							cosOutputStorageMap["region"] = snapshotByTimeOffsetTaskSet.OutputStorage.CosOutputStorage.Region
						}

						outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
					}

					snapshotByTimeOffsetTaskSetMap["output_storage"] = []interface{}{outputStorageMap}
				}

				if snapshotByTimeOffsetTaskSet.OutputObjectPath != nil {
					snapshotByTimeOffsetTaskSetMap["output_object_path"] = snapshotByTimeOffsetTaskSet.OutputObjectPath
				}

				if snapshotByTimeOffsetTaskSet.ObjectNumberFormat != nil {
					objectNumberFormatMap := map[string]interface{}{}

					if snapshotByTimeOffsetTaskSet.ObjectNumberFormat.InitialValue != nil {
						objectNumberFormatMap["initial_value"] = snapshotByTimeOffsetTaskSet.ObjectNumberFormat.InitialValue
					}

					if snapshotByTimeOffsetTaskSet.ObjectNumberFormat.Increment != nil {
						objectNumberFormatMap["increment"] = snapshotByTimeOffsetTaskSet.ObjectNumberFormat.Increment
					}

					if snapshotByTimeOffsetTaskSet.ObjectNumberFormat.MinLength != nil {
						objectNumberFormatMap["min_length"] = snapshotByTimeOffsetTaskSet.ObjectNumberFormat.MinLength
					}

					if snapshotByTimeOffsetTaskSet.ObjectNumberFormat.PlaceHolder != nil {
						objectNumberFormatMap["place_holder"] = snapshotByTimeOffsetTaskSet.ObjectNumberFormat.PlaceHolder
					}

					snapshotByTimeOffsetTaskSetMap["object_number_format"] = []interface{}{objectNumberFormatMap}
				}

				snapshotByTimeOffsetTaskSetList = append(snapshotByTimeOffsetTaskSetList, snapshotByTimeOffsetTaskSetMap)
			}

			mediaProcessTaskMap["snapshot_by_time_offset_task_set"] = snapshotByTimeOffsetTaskSetList
		}

		if workflow.MediaProcessTask.SampleSnapshotTaskSet != nil {
			sampleSnapshotTaskSetList := []interface{}{}
			for _, sampleSnapshotTaskSet := range workflow.MediaProcessTask.SampleSnapshotTaskSet {
				sampleSnapshotTaskSetMap := map[string]interface{}{}

				if sampleSnapshotTaskSet.Definition != nil {
					sampleSnapshotTaskSetMap["definition"] = sampleSnapshotTaskSet.Definition
				}

				if sampleSnapshotTaskSet.WatermarkSet != nil {
					watermarkSetList := []interface{}{}
					for _, watermarkSet := range sampleSnapshotTaskSet.WatermarkSet {
						watermarkSetMap := map[string]interface{}{}

						if watermarkSet.Definition != nil {
							watermarkSetMap["definition"] = watermarkSet.Definition
						}

						if watermarkSet.RawParameter != nil {
							rawParameterMap := map[string]interface{}{}

							if watermarkSet.RawParameter.Type != nil {
								rawParameterMap["type"] = watermarkSet.RawParameter.Type
							}

							if watermarkSet.RawParameter.CoordinateOrigin != nil {
								rawParameterMap["coordinate_origin"] = watermarkSet.RawParameter.CoordinateOrigin
							}

							if watermarkSet.RawParameter.XPos != nil {
								rawParameterMap["x_pos"] = watermarkSet.RawParameter.XPos
							}

							if watermarkSet.RawParameter.YPos != nil {
								rawParameterMap["y_pos"] = watermarkSet.RawParameter.YPos
							}

							if watermarkSet.RawParameter.ImageTemplate != nil {
								imageTemplateMap := map[string]interface{}{}

								if watermarkSet.RawParameter.ImageTemplate.ImageContent != nil {
									imageContentMap := map[string]interface{}{}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.Type != nil {
										imageContentMap["type"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.Type
									}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo != nil {
										cosInputInfoMap := map[string]interface{}{}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Bucket != nil {
											cosInputInfoMap["bucket"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Bucket
										}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Region != nil {
											cosInputInfoMap["region"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Region
										}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Object != nil {
											cosInputInfoMap["object"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Object
										}

										imageContentMap["cos_input_info"] = []interface{}{cosInputInfoMap}
									}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo != nil {
										urlInputInfoMap := map[string]interface{}{}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo.Url != nil {
											urlInputInfoMap["url"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo.Url
										}

										imageContentMap["url_input_info"] = []interface{}{urlInputInfoMap}
									}

									imageTemplateMap["image_content"] = []interface{}{imageContentMap}
								}

								if watermarkSet.RawParameter.ImageTemplate.Width != nil {
									imageTemplateMap["width"] = watermarkSet.RawParameter.ImageTemplate.Width
								}

								if watermarkSet.RawParameter.ImageTemplate.Height != nil {
									imageTemplateMap["height"] = watermarkSet.RawParameter.ImageTemplate.Height
								}

								if watermarkSet.RawParameter.ImageTemplate.RepeatType != nil {
									imageTemplateMap["repeat_type"] = watermarkSet.RawParameter.ImageTemplate.RepeatType
								}

								rawParameterMap["image_template"] = []interface{}{imageTemplateMap}
							}

							watermarkSetMap["raw_parameter"] = []interface{}{rawParameterMap}
						}

						if watermarkSet.TextContent != nil {
							watermarkSetMap["text_content"] = watermarkSet.TextContent
						}

						if watermarkSet.SvgContent != nil {
							watermarkSetMap["svg_content"] = watermarkSet.SvgContent
						}

						if watermarkSet.StartTimeOffset != nil {
							watermarkSetMap["start_time_offset"] = watermarkSet.StartTimeOffset
						}

						if watermarkSet.EndTimeOffset != nil {
							watermarkSetMap["end_time_offset"] = watermarkSet.EndTimeOffset
						}

						watermarkSetList = append(watermarkSetList, watermarkSetMap)
					}

					sampleSnapshotTaskSetMap["watermark_set"] = watermarkSetList
				}

				if sampleSnapshotTaskSet.OutputStorage != nil {
					outputStorageMap := map[string]interface{}{}

					if sampleSnapshotTaskSet.OutputStorage.Type != nil {
						outputStorageMap["type"] = sampleSnapshotTaskSet.OutputStorage.Type
					}

					if sampleSnapshotTaskSet.OutputStorage.CosOutputStorage != nil {
						cosOutputStorageMap := map[string]interface{}{}

						if sampleSnapshotTaskSet.OutputStorage.CosOutputStorage.Bucket != nil {
							cosOutputStorageMap["bucket"] = sampleSnapshotTaskSet.OutputStorage.CosOutputStorage.Bucket
						}

						if sampleSnapshotTaskSet.OutputStorage.CosOutputStorage.Region != nil {
							cosOutputStorageMap["region"] = sampleSnapshotTaskSet.OutputStorage.CosOutputStorage.Region
						}

						outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
					}

					sampleSnapshotTaskSetMap["output_storage"] = []interface{}{outputStorageMap}
				}

				if sampleSnapshotTaskSet.OutputObjectPath != nil {
					sampleSnapshotTaskSetMap["output_object_path"] = sampleSnapshotTaskSet.OutputObjectPath
				}

				if sampleSnapshotTaskSet.ObjectNumberFormat != nil {
					objectNumberFormatMap := map[string]interface{}{}

					if sampleSnapshotTaskSet.ObjectNumberFormat.InitialValue != nil {
						objectNumberFormatMap["initial_value"] = sampleSnapshotTaskSet.ObjectNumberFormat.InitialValue
					}

					if sampleSnapshotTaskSet.ObjectNumberFormat.Increment != nil {
						objectNumberFormatMap["increment"] = sampleSnapshotTaskSet.ObjectNumberFormat.Increment
					}

					if sampleSnapshotTaskSet.ObjectNumberFormat.MinLength != nil {
						objectNumberFormatMap["min_length"] = sampleSnapshotTaskSet.ObjectNumberFormat.MinLength
					}

					if sampleSnapshotTaskSet.ObjectNumberFormat.PlaceHolder != nil {
						objectNumberFormatMap["place_holder"] = sampleSnapshotTaskSet.ObjectNumberFormat.PlaceHolder
					}

					sampleSnapshotTaskSetMap["object_number_format"] = []interface{}{objectNumberFormatMap}
				}

				sampleSnapshotTaskSetList = append(sampleSnapshotTaskSetList, sampleSnapshotTaskSetMap)
			}

			mediaProcessTaskMap["sample_snapshot_task_set"] = sampleSnapshotTaskSetList
		}

		if workflow.MediaProcessTask.ImageSpriteTaskSet != nil {
			imageSpriteTaskSetList := []interface{}{}
			for _, imageSpriteTaskSet := range workflow.MediaProcessTask.ImageSpriteTaskSet {
				imageSpriteTaskSetMap := map[string]interface{}{}

				if imageSpriteTaskSet.Definition != nil {
					imageSpriteTaskSetMap["definition"] = imageSpriteTaskSet.Definition
				}

				if imageSpriteTaskSet.OutputStorage != nil {
					outputStorageMap := map[string]interface{}{}

					if imageSpriteTaskSet.OutputStorage.Type != nil {
						outputStorageMap["type"] = imageSpriteTaskSet.OutputStorage.Type
					}

					if imageSpriteTaskSet.OutputStorage.CosOutputStorage != nil {
						cosOutputStorageMap := map[string]interface{}{}

						if imageSpriteTaskSet.OutputStorage.CosOutputStorage.Bucket != nil {
							cosOutputStorageMap["bucket"] = imageSpriteTaskSet.OutputStorage.CosOutputStorage.Bucket
						}

						if imageSpriteTaskSet.OutputStorage.CosOutputStorage.Region != nil {
							cosOutputStorageMap["region"] = imageSpriteTaskSet.OutputStorage.CosOutputStorage.Region
						}

						outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
					}

					imageSpriteTaskSetMap["output_storage"] = []interface{}{outputStorageMap}
				}

				if imageSpriteTaskSet.OutputObjectPath != nil {
					imageSpriteTaskSetMap["output_object_path"] = imageSpriteTaskSet.OutputObjectPath
				}

				if imageSpriteTaskSet.WebVttObjectName != nil {
					imageSpriteTaskSetMap["web_vtt_object_name"] = imageSpriteTaskSet.WebVttObjectName
				}

				if imageSpriteTaskSet.ObjectNumberFormat != nil {
					objectNumberFormatMap := map[string]interface{}{}

					if imageSpriteTaskSet.ObjectNumberFormat.InitialValue != nil {
						objectNumberFormatMap["initial_value"] = imageSpriteTaskSet.ObjectNumberFormat.InitialValue
					}

					if imageSpriteTaskSet.ObjectNumberFormat.Increment != nil {
						objectNumberFormatMap["increment"] = imageSpriteTaskSet.ObjectNumberFormat.Increment
					}

					if imageSpriteTaskSet.ObjectNumberFormat.MinLength != nil {
						objectNumberFormatMap["min_length"] = imageSpriteTaskSet.ObjectNumberFormat.MinLength
					}

					if imageSpriteTaskSet.ObjectNumberFormat.PlaceHolder != nil {
						objectNumberFormatMap["place_holder"] = imageSpriteTaskSet.ObjectNumberFormat.PlaceHolder
					}

					imageSpriteTaskSetMap["object_number_format"] = []interface{}{objectNumberFormatMap}
				}

				imageSpriteTaskSetList = append(imageSpriteTaskSetList, imageSpriteTaskSetMap)
			}

			mediaProcessTaskMap["image_sprite_task_set"] = imageSpriteTaskSetList
		}

		if workflow.MediaProcessTask.AdaptiveDynamicStreamingTaskSet != nil {
			adaptiveDynamicStreamingTaskSetList := []interface{}{}
			for _, adaptiveDynamicStreamingTaskSet := range workflow.MediaProcessTask.AdaptiveDynamicStreamingTaskSet {
				adaptiveDynamicStreamingTaskSetMap := map[string]interface{}{}

				if adaptiveDynamicStreamingTaskSet.Definition != nil {
					adaptiveDynamicStreamingTaskSetMap["definition"] = adaptiveDynamicStreamingTaskSet.Definition
				}

				if adaptiveDynamicStreamingTaskSet.WatermarkSet != nil {
					watermarkSetList := []interface{}{}
					for _, watermarkSet := range adaptiveDynamicStreamingTaskSet.WatermarkSet {
						watermarkSetMap := map[string]interface{}{}

						if watermarkSet.Definition != nil {
							watermarkSetMap["definition"] = watermarkSet.Definition
						}

						if watermarkSet.RawParameter != nil {
							rawParameterMap := map[string]interface{}{}

							if watermarkSet.RawParameter.Type != nil {
								rawParameterMap["type"] = watermarkSet.RawParameter.Type
							}

							if watermarkSet.RawParameter.CoordinateOrigin != nil {
								rawParameterMap["coordinate_origin"] = watermarkSet.RawParameter.CoordinateOrigin
							}

							if watermarkSet.RawParameter.XPos != nil {
								rawParameterMap["x_pos"] = watermarkSet.RawParameter.XPos
							}

							if watermarkSet.RawParameter.YPos != nil {
								rawParameterMap["y_pos"] = watermarkSet.RawParameter.YPos
							}

							if watermarkSet.RawParameter.ImageTemplate != nil {
								imageTemplateMap := map[string]interface{}{}

								if watermarkSet.RawParameter.ImageTemplate.ImageContent != nil {
									imageContentMap := map[string]interface{}{}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.Type != nil {
										imageContentMap["type"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.Type
									}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo != nil {
										cosInputInfoMap := map[string]interface{}{}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Bucket != nil {
											cosInputInfoMap["bucket"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Bucket
										}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Region != nil {
											cosInputInfoMap["region"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Region
										}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Object != nil {
											cosInputInfoMap["object"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.CosInputInfo.Object
										}

										imageContentMap["cos_input_info"] = []interface{}{cosInputInfoMap}
									}

									if watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo != nil {
										urlInputInfoMap := map[string]interface{}{}

										if watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo.Url != nil {
											urlInputInfoMap["url"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.UrlInputInfo.Url
										}

										imageContentMap["url_input_info"] = []interface{}{urlInputInfoMap}
									}

									imageTemplateMap["image_content"] = []interface{}{imageContentMap}
								}

								if watermarkSet.RawParameter.ImageTemplate.Width != nil {
									imageTemplateMap["width"] = watermarkSet.RawParameter.ImageTemplate.Width
								}

								if watermarkSet.RawParameter.ImageTemplate.Height != nil {
									imageTemplateMap["height"] = watermarkSet.RawParameter.ImageTemplate.Height
								}

								if watermarkSet.RawParameter.ImageTemplate.RepeatType != nil {
									imageTemplateMap["repeat_type"] = watermarkSet.RawParameter.ImageTemplate.RepeatType
								}

								rawParameterMap["image_template"] = []interface{}{imageTemplateMap}
							}

							watermarkSetMap["raw_parameter"] = []interface{}{rawParameterMap}
						}

						if watermarkSet.TextContent != nil {
							watermarkSetMap["text_content"] = watermarkSet.TextContent
						}

						if watermarkSet.SvgContent != nil {
							watermarkSetMap["svg_content"] = watermarkSet.SvgContent
						}

						if watermarkSet.StartTimeOffset != nil {
							watermarkSetMap["start_time_offset"] = watermarkSet.StartTimeOffset
						}

						if watermarkSet.EndTimeOffset != nil {
							watermarkSetMap["end_time_offset"] = watermarkSet.EndTimeOffset
						}

						watermarkSetList = append(watermarkSetList, watermarkSetMap)
					}

					adaptiveDynamicStreamingTaskSetMap["watermark_set"] = watermarkSetList
				}

				if adaptiveDynamicStreamingTaskSet.OutputStorage != nil {
					outputStorageMap := map[string]interface{}{}

					if adaptiveDynamicStreamingTaskSet.OutputStorage.Type != nil {
						outputStorageMap["type"] = adaptiveDynamicStreamingTaskSet.OutputStorage.Type
					}

					if adaptiveDynamicStreamingTaskSet.OutputStorage.CosOutputStorage != nil {
						cosOutputStorageMap := map[string]interface{}{}

						if adaptiveDynamicStreamingTaskSet.OutputStorage.CosOutputStorage.Bucket != nil {
							cosOutputStorageMap["bucket"] = adaptiveDynamicStreamingTaskSet.OutputStorage.CosOutputStorage.Bucket
						}

						if adaptiveDynamicStreamingTaskSet.OutputStorage.CosOutputStorage.Region != nil {
							cosOutputStorageMap["region"] = adaptiveDynamicStreamingTaskSet.OutputStorage.CosOutputStorage.Region
						}

						outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
					}

					adaptiveDynamicStreamingTaskSetMap["output_storage"] = []interface{}{outputStorageMap}
				}

				if adaptiveDynamicStreamingTaskSet.OutputObjectPath != nil {
					adaptiveDynamicStreamingTaskSetMap["output_object_path"] = adaptiveDynamicStreamingTaskSet.OutputObjectPath
				}

				if adaptiveDynamicStreamingTaskSet.SubStreamObjectName != nil {
					adaptiveDynamicStreamingTaskSetMap["sub_stream_object_name"] = adaptiveDynamicStreamingTaskSet.SubStreamObjectName
				}

				if adaptiveDynamicStreamingTaskSet.SegmentObjectName != nil {
					adaptiveDynamicStreamingTaskSetMap["segment_object_name"] = adaptiveDynamicStreamingTaskSet.SegmentObjectName
				}

				adaptiveDynamicStreamingTaskSetList = append(adaptiveDynamicStreamingTaskSetList, adaptiveDynamicStreamingTaskSetMap)
			}

			mediaProcessTaskMap["adaptive_dynamic_streaming_task_set"] = adaptiveDynamicStreamingTaskSetList
		}

		_ = d.Set("media_process_task", []interface{}{mediaProcessTaskMap})
	}

	if workflow.AiContentReviewTask != nil {
		aiContentReviewTaskMap := map[string]interface{}{}

		if workflow.AiContentReviewTask.Definition != nil {
			aiContentReviewTaskMap["definition"] = workflow.AiContentReviewTask.Definition
		}

		_ = d.Set("ai_content_review_task", []interface{}{aiContentReviewTaskMap})
	}

	if workflow.AiAnalysisTask != nil {
		aiAnalysisTaskMap := map[string]interface{}{}

		if workflow.AiAnalysisTask.Definition != nil {
			aiAnalysisTaskMap["definition"] = workflow.AiAnalysisTask.Definition
		}

		if workflow.AiAnalysisTask.ExtendedParameter != nil {
			aiAnalysisTaskMap["extended_parameter"] = workflow.AiAnalysisTask.ExtendedParameter
		}

		_ = d.Set("ai_analysis_task", []interface{}{aiAnalysisTaskMap})
	}

	if workflow.AiRecognitionTask != nil {
		aiRecognitionTaskMap := map[string]interface{}{}

		if workflow.AiRecognitionTask.Definition != nil {
			aiRecognitionTaskMap["definition"] = workflow.AiRecognitionTask.Definition
		}

		_ = d.Set("ai_recognition_task", []interface{}{aiRecognitionTaskMap})
	}

	if workflow.TaskNotifyConfig != nil {
		taskNotifyConfigMap := map[string]interface{}{}

		if workflow.TaskNotifyConfig.CmqModel != nil {
			taskNotifyConfigMap["cmq_model"] = workflow.TaskNotifyConfig.CmqModel
		}

		if workflow.TaskNotifyConfig.CmqRegion != nil {
			taskNotifyConfigMap["cmq_region"] = workflow.TaskNotifyConfig.CmqRegion
		}

		if workflow.TaskNotifyConfig.TopicName != nil {
			taskNotifyConfigMap["topic_name"] = workflow.TaskNotifyConfig.TopicName
		}

		if workflow.TaskNotifyConfig.QueueName != nil {
			taskNotifyConfigMap["queue_name"] = workflow.TaskNotifyConfig.QueueName
		}

		if workflow.TaskNotifyConfig.NotifyMode != nil {
			taskNotifyConfigMap["notify_mode"] = workflow.TaskNotifyConfig.NotifyMode
		}

		if workflow.TaskNotifyConfig.NotifyType != nil {
			taskNotifyConfigMap["notify_type"] = workflow.TaskNotifyConfig.NotifyType
		}

		if workflow.TaskNotifyConfig.NotifyUrl != nil {
			taskNotifyConfigMap["notify_url"] = workflow.TaskNotifyConfig.NotifyUrl
		}

		_ = d.Set("task_notify_config", []interface{}{taskNotifyConfigMap})
	}

	if workflow.TaskPriority != nil {
		_ = d.Set("task_priority", workflow.TaskPriority)
	}

	return nil
}

func resourceTencentCloudMpsWorkflowUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_workflow.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := mps.NewResetWorkflowRequest()

	workflowId := d.Id()

	request.WorkflowId = helper.StrToInt64Point(workflowId)

	isChanged := false

	mutableArgs := []string{
		"workflow_name", "trigger", "output_storage",
		"output_dir", "media_process_task", "ai_content_review_task",
		"ai_analysis_task", "ai_recognition_task", "task_notify_config", "task_priority",
	}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			isChanged = true
			break
		}
	}

	if isChanged {
		if v, ok := d.GetOk("workflow_name"); ok {
			request.WorkflowName = helper.String(v.(string))
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "trigger"); ok {
			workflowTrigger := mps.WorkflowTrigger{}
			if v, ok := dMap["type"]; ok {
				workflowTrigger.Type = helper.String(v.(string))
			}
			if cosFileUploadTriggerMap, ok := helper.InterfaceToMap(dMap, "cos_file_upload_trigger"); ok {
				cosFileUploadTrigger := mps.CosFileUploadTrigger{}
				if v, ok := cosFileUploadTriggerMap["bucket"]; ok {
					cosFileUploadTrigger.Bucket = helper.String(v.(string))
				}
				if v, ok := cosFileUploadTriggerMap["region"]; ok {
					cosFileUploadTrigger.Region = helper.String(v.(string))
				}
				if v, ok := cosFileUploadTriggerMap["dir"]; ok {
					cosFileUploadTrigger.Dir = helper.String(v.(string))
				}
				if v, ok := cosFileUploadTriggerMap["formats"]; ok {
					formatsSet := v.(*schema.Set).List()
					for i := range formatsSet {
						formats := formatsSet[i].(string)
						cosFileUploadTrigger.Formats = append(cosFileUploadTrigger.Formats, &formats)
					}
				}
				workflowTrigger.CosFileUploadTrigger = &cosFileUploadTrigger
			}
			request.Trigger = &workflowTrigger
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "output_storage"); ok {
			taskOutputStorage := mps.TaskOutputStorage{}
			if v, ok := dMap["type"]; ok {
				taskOutputStorage.Type = helper.String(v.(string))
			}
			if cosOutputStorageMap, ok := helper.InterfaceToMap(dMap, "cos_output_storage"); ok {
				cosOutputStorage := mps.CosOutputStorage{}
				if v, ok := cosOutputStorageMap["bucket"]; ok {
					cosOutputStorage.Bucket = helper.String(v.(string))
				}
				if v, ok := cosOutputStorageMap["region"]; ok {
					cosOutputStorage.Region = helper.String(v.(string))
				}
				taskOutputStorage.CosOutputStorage = &cosOutputStorage
			}
			request.OutputStorage = &taskOutputStorage
		}

		if v, ok := d.GetOk("output_dir"); ok {
			request.OutputDir = helper.String(v.(string))
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "media_process_task"); ok {
			mediaProcessTaskInput := mps.MediaProcessTaskInput{}
			if v, ok := dMap["transcode_task_set"]; ok {
				for _, item := range v.([]interface{}) {
					transcodeTaskSetMap := item.(map[string]interface{})
					transcodeTaskInput := mps.TranscodeTaskInput{}
					if v, ok := transcodeTaskSetMap["definition"]; ok {
						transcodeTaskInput.Definition = helper.IntUint64(v.(int))
					}
					if rawParameterMap, ok := helper.InterfaceToMap(transcodeTaskSetMap, "raw_parameter"); ok {
						rawTranscodeParameter := mps.RawTranscodeParameter{}
						if v, ok := rawParameterMap["container"]; ok {
							rawTranscodeParameter.Container = helper.String(v.(string))
						}
						if v, ok := rawParameterMap["remove_video"]; ok {
							rawTranscodeParameter.RemoveVideo = helper.IntInt64(v.(int))
						}
						if v, ok := rawParameterMap["remove_audio"]; ok {
							rawTranscodeParameter.RemoveAudio = helper.IntInt64(v.(int))
						}
						if videoTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "video_template"); ok {
							videoTemplateInfo := mps.VideoTemplateInfo{}
							if v, ok := videoTemplateMap["codec"]; ok {
								videoTemplateInfo.Codec = helper.String(v.(string))
							}
							if v, ok := videoTemplateMap["fps"]; ok {
								videoTemplateInfo.Fps = helper.IntInt64(v.(int))
							}
							if v, ok := videoTemplateMap["bitrate"]; ok {
								videoTemplateInfo.Bitrate = helper.IntInt64(v.(int))
							}
							if v, ok := videoTemplateMap["resolution_adaptive"]; ok {
								videoTemplateInfo.ResolutionAdaptive = helper.String(v.(string))
							}
							if v, ok := videoTemplateMap["width"]; ok {
								videoTemplateInfo.Width = helper.IntUint64(v.(int))
							}
							if v, ok := videoTemplateMap["height"]; ok {
								videoTemplateInfo.Height = helper.IntUint64(v.(int))
							}
							if v, ok := videoTemplateMap["gop"]; ok {
								videoTemplateInfo.Gop = helper.IntUint64(v.(int))
							}
							if v, ok := videoTemplateMap["fill_type"]; ok {
								videoTemplateInfo.FillType = helper.String(v.(string))
							}
							if v, ok := videoTemplateMap["vcrf"]; ok {
								videoTemplateInfo.Vcrf = helper.IntUint64(v.(int))
							}
							rawTranscodeParameter.VideoTemplate = &videoTemplateInfo
						}
						if audioTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "audio_template"); ok {
							audioTemplateInfo := mps.AudioTemplateInfo{}
							if v, ok := audioTemplateMap["codec"]; ok {
								audioTemplateInfo.Codec = helper.String(v.(string))
							}
							if v, ok := audioTemplateMap["bitrate"]; ok {
								audioTemplateInfo.Bitrate = helper.IntInt64(v.(int))
							}
							if v, ok := audioTemplateMap["sample_rate"]; ok {
								audioTemplateInfo.SampleRate = helper.IntUint64(v.(int))
							}
							if v, ok := audioTemplateMap["audio_channel"]; ok {
								audioTemplateInfo.AudioChannel = helper.IntInt64(v.(int))
							}
							rawTranscodeParameter.AudioTemplate = &audioTemplateInfo
						}
						if tEHDConfigMap, ok := helper.InterfaceToMap(rawParameterMap, "tehd_config"); ok {
							tEHDConfig := mps.TEHDConfig{}
							if v, ok := tEHDConfigMap["type"]; ok {
								tEHDConfig.Type = helper.String(v.(string))
							}
							if v, ok := tEHDConfigMap["max_video_bitrate"]; ok {
								tEHDConfig.MaxVideoBitrate = helper.IntInt64(v.(int))
							}
							rawTranscodeParameter.TEHDConfig = &tEHDConfig
						}
						transcodeTaskInput.RawParameter = &rawTranscodeParameter
					}
					if overrideParameterMap, ok := helper.InterfaceToMap(transcodeTaskSetMap, "override_parameter"); ok {
						overrideTranscodeParameter := mps.OverrideTranscodeParameter{}
						if v, ok := overrideParameterMap["container"]; ok {
							overrideTranscodeParameter.Container = helper.String(v.(string))
						}
						if v, ok := overrideParameterMap["remove_video"]; ok {
							overrideTranscodeParameter.RemoveVideo = helper.IntUint64(v.(int))
						}
						if v, ok := overrideParameterMap["remove_audio"]; ok {
							overrideTranscodeParameter.RemoveAudio = helper.IntUint64(v.(int))
						}
						if videoTemplateMap, ok := helper.InterfaceToMap(overrideParameterMap, "video_template"); ok {
							videoTemplateInfoForUpdate := mps.VideoTemplateInfoForUpdate{}
							if v, ok := videoTemplateMap["codec"]; ok {
								videoTemplateInfoForUpdate.Codec = helper.String(v.(string))
							}
							if v, ok := videoTemplateMap["fps"]; ok {
								videoTemplateInfoForUpdate.Fps = helper.IntInt64(v.(int))
							}
							if v, ok := videoTemplateMap["bitrate"]; ok {
								videoTemplateInfoForUpdate.Bitrate = helper.IntInt64(v.(int))
							}
							if v, ok := videoTemplateMap["resolution_adaptive"]; ok {
								videoTemplateInfoForUpdate.ResolutionAdaptive = helper.String(v.(string))
							}
							if v, ok := videoTemplateMap["width"]; ok {
								videoTemplateInfoForUpdate.Width = helper.IntUint64(v.(int))
							}
							if v, ok := videoTemplateMap["height"]; ok {
								videoTemplateInfoForUpdate.Height = helper.IntUint64(v.(int))
							}
							if v, ok := videoTemplateMap["gop"]; ok {
								videoTemplateInfoForUpdate.Gop = helper.IntUint64(v.(int))
							}
							if v, ok := videoTemplateMap["fill_type"]; ok {
								videoTemplateInfoForUpdate.FillType = helper.String(v.(string))
							}
							if v, ok := videoTemplateMap["vcrf"]; ok {
								videoTemplateInfoForUpdate.Vcrf = helper.IntUint64(v.(int))
							}
							if v, ok := videoTemplateMap["content_adapt_stream"]; ok {
								videoTemplateInfoForUpdate.ContentAdaptStream = helper.IntUint64(v.(int))
							}
							overrideTranscodeParameter.VideoTemplate = &videoTemplateInfoForUpdate
						}
						if audioTemplateMap, ok := helper.InterfaceToMap(overrideParameterMap, "audio_template"); ok {
							audioTemplateInfoForUpdate := mps.AudioTemplateInfoForUpdate{}
							if v, ok := audioTemplateMap["codec"]; ok {
								audioTemplateInfoForUpdate.Codec = helper.String(v.(string))
							}
							if v, ok := audioTemplateMap["bitrate"]; ok {
								audioTemplateInfoForUpdate.Bitrate = helper.IntInt64(v.(int))
							}
							if v, ok := audioTemplateMap["sample_rate"]; ok {
								audioTemplateInfoForUpdate.SampleRate = helper.IntUint64(v.(int))
							}
							if v, ok := audioTemplateMap["audio_channel"]; ok {
								audioTemplateInfoForUpdate.AudioChannel = helper.IntInt64(v.(int))
							}
							if v, ok := audioTemplateMap["stream_selects"]; ok {
								streamSelectsSet := v.(*schema.Set).List()
								for i := range streamSelectsSet {
									streamSelects := streamSelectsSet[i].(int)
									audioTemplateInfoForUpdate.StreamSelects = append(audioTemplateInfoForUpdate.StreamSelects, helper.IntInt64(streamSelects))
								}
							}
							overrideTranscodeParameter.AudioTemplate = &audioTemplateInfoForUpdate
						}
						if tEHDConfigMap, ok := helper.InterfaceToMap(overrideParameterMap, "tehd_config"); ok {
							tEHDConfigForUpdate := mps.TEHDConfigForUpdate{}
							if v, ok := tEHDConfigMap["type"]; ok {
								tEHDConfigForUpdate.Type = helper.String(v.(string))
							}
							if v, ok := tEHDConfigMap["max_video_bitrate"]; ok {
								tEHDConfigForUpdate.MaxVideoBitrate = helper.IntInt64(v.(int))
							}
							overrideTranscodeParameter.TEHDConfig = &tEHDConfigForUpdate
						}
						if subtitleTemplateMap, ok := helper.InterfaceToMap(overrideParameterMap, "subtitle_template"); ok {
							subtitleTemplate := mps.SubtitleTemplate{}
							if v, ok := subtitleTemplateMap["path"]; ok {
								subtitleTemplate.Path = helper.String(v.(string))
							}
							if v, ok := subtitleTemplateMap["stream_index"]; ok {
								subtitleTemplate.StreamIndex = helper.IntInt64(v.(int))
							}
							if v, ok := subtitleTemplateMap["font_type"]; ok {
								subtitleTemplate.FontType = helper.String(v.(string))
							}
							if v, ok := subtitleTemplateMap["font_size"]; ok {
								subtitleTemplate.FontSize = helper.String(v.(string))
							}
							if v, ok := subtitleTemplateMap["font_color"]; ok {
								subtitleTemplate.FontColor = helper.String(v.(string))
							}
							if v, ok := subtitleTemplateMap["font_alpha"]; ok {
								subtitleTemplate.FontAlpha = helper.Float64(v.(float64))
							}
							overrideTranscodeParameter.SubtitleTemplate = &subtitleTemplate
						}
						transcodeTaskInput.OverrideParameter = &overrideTranscodeParameter
					}
					if v, ok := transcodeTaskSetMap["watermark_set"]; ok {
						for _, item := range v.([]interface{}) {
							watermarkSetMap := item.(map[string]interface{})
							watermarkInput := mps.WatermarkInput{}
							if v, ok := watermarkSetMap["definition"]; ok {
								watermarkInput.Definition = helper.IntUint64(v.(int))
							}
							if rawParameterMap, ok := helper.InterfaceToMap(watermarkSetMap, "raw_parameter"); ok {
								rawWatermarkParameter := mps.RawWatermarkParameter{}
								if v, ok := rawParameterMap["type"]; ok {
									rawWatermarkParameter.Type = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["coordinate_origin"]; ok {
									rawWatermarkParameter.CoordinateOrigin = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["x_pos"]; ok {
									rawWatermarkParameter.XPos = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["y_pos"]; ok {
									rawWatermarkParameter.YPos = helper.String(v.(string))
								}
								if imageTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "image_template"); ok {
									rawImageWatermarkInput := mps.RawImageWatermarkInput{}
									if imageContentMap, ok := helper.InterfaceToMap(imageTemplateMap, "image_content"); ok {
										mediaInputInfo := mps.MediaInputInfo{}
										if v, ok := imageContentMap["type"]; ok {
											mediaInputInfo.Type = helper.String(v.(string))
										}
										if cosInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "cos_input_info"); ok {
											cosInputInfo := mps.CosInputInfo{}
											if v, ok := cosInputInfoMap["bucket"]; ok {
												cosInputInfo.Bucket = helper.String(v.(string))
											}
											if v, ok := cosInputInfoMap["region"]; ok {
												cosInputInfo.Region = helper.String(v.(string))
											}
											if v, ok := cosInputInfoMap["object"]; ok {
												cosInputInfo.Object = helper.String(v.(string))
											}
											mediaInputInfo.CosInputInfo = &cosInputInfo
										}
										if urlInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "url_input_info"); ok {
											urlInputInfo := mps.UrlInputInfo{}
											if v, ok := urlInputInfoMap["url"]; ok {
												urlInputInfo.Url = helper.String(v.(string))
											}
											mediaInputInfo.UrlInputInfo = &urlInputInfo
										}
										rawImageWatermarkInput.ImageContent = &mediaInputInfo
									}
									if v, ok := imageTemplateMap["width"]; ok {
										rawImageWatermarkInput.Width = helper.String(v.(string))
									}
									if v, ok := imageTemplateMap["height"]; ok {
										rawImageWatermarkInput.Height = helper.String(v.(string))
									}
									if v, ok := imageTemplateMap["repeat_type"]; ok {
										rawImageWatermarkInput.RepeatType = helper.String(v.(string))
									}
									rawWatermarkParameter.ImageTemplate = &rawImageWatermarkInput
								}
								watermarkInput.RawParameter = &rawWatermarkParameter
							}
							if v, ok := watermarkSetMap["text_content"]; ok {
								watermarkInput.TextContent = helper.String(v.(string))
							}
							if v, ok := watermarkSetMap["svg_content"]; ok {
								watermarkInput.SvgContent = helper.String(v.(string))
							}
							if v, ok := watermarkSetMap["start_time_offset"]; ok {
								watermarkInput.StartTimeOffset = helper.Float64(v.(float64))
							}
							if v, ok := watermarkSetMap["end_time_offset"]; ok {
								watermarkInput.EndTimeOffset = helper.Float64(v.(float64))
							}
							transcodeTaskInput.WatermarkSet = append(transcodeTaskInput.WatermarkSet, &watermarkInput)
						}
					}
					if v, ok := transcodeTaskSetMap["mosaic_set"]; ok {
						for _, item := range v.([]interface{}) {
							mosaicSetMap := item.(map[string]interface{})
							mosaicInput := mps.MosaicInput{}
							if v, ok := mosaicSetMap["coordinate_origin"]; ok {
								mosaicInput.CoordinateOrigin = helper.String(v.(string))
							}
							if v, ok := mosaicSetMap["x_pos"]; ok {
								mosaicInput.XPos = helper.String(v.(string))
							}
							if v, ok := mosaicSetMap["y_pos"]; ok {
								mosaicInput.YPos = helper.String(v.(string))
							}
							if v, ok := mosaicSetMap["width"]; ok {
								mosaicInput.Width = helper.String(v.(string))
							}
							if v, ok := mosaicSetMap["height"]; ok {
								mosaicInput.Height = helper.String(v.(string))
							}
							if v, ok := mosaicSetMap["start_time_offset"]; ok {
								mosaicInput.StartTimeOffset = helper.Float64(v.(float64))
							}
							if v, ok := mosaicSetMap["end_time_offset"]; ok {
								mosaicInput.EndTimeOffset = helper.Float64(v.(float64))
							}
							transcodeTaskInput.MosaicSet = append(transcodeTaskInput.MosaicSet, &mosaicInput)
						}
					}
					if v, ok := transcodeTaskSetMap["start_time_offset"]; ok {
						transcodeTaskInput.StartTimeOffset = helper.Float64(v.(float64))
					}
					if v, ok := transcodeTaskSetMap["end_time_offset"]; ok {
						transcodeTaskInput.EndTimeOffset = helper.Float64(v.(float64))
					}
					if outputStorageMap, ok := helper.InterfaceToMap(transcodeTaskSetMap, "output_storage"); ok {
						taskOutputStorage := mps.TaskOutputStorage{}
						if v, ok := outputStorageMap["type"]; ok {
							taskOutputStorage.Type = helper.String(v.(string))
						}
						if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
							cosOutputStorage := mps.CosOutputStorage{}
							if v, ok := cosOutputStorageMap["bucket"]; ok {
								cosOutputStorage.Bucket = helper.String(v.(string))
							}
							if v, ok := cosOutputStorageMap["region"]; ok {
								cosOutputStorage.Region = helper.String(v.(string))
							}
							taskOutputStorage.CosOutputStorage = &cosOutputStorage
						}
						transcodeTaskInput.OutputStorage = &taskOutputStorage
					}
					if v, ok := transcodeTaskSetMap["output_object_path"]; ok {
						transcodeTaskInput.OutputObjectPath = helper.String(v.(string))
					}
					if v, ok := transcodeTaskSetMap["segment_object_name"]; ok {
						transcodeTaskInput.SegmentObjectName = helper.String(v.(string))
					}
					if objectNumberFormatMap, ok := helper.InterfaceToMap(transcodeTaskSetMap, "object_number_format"); ok {
						numberFormat := mps.NumberFormat{}
						if v, ok := objectNumberFormatMap["initial_value"]; ok {
							numberFormat.InitialValue = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["increment"]; ok {
							numberFormat.Increment = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["min_length"]; ok {
							numberFormat.MinLength = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["place_holder"]; ok {
							numberFormat.PlaceHolder = helper.String(v.(string))
						}
						transcodeTaskInput.ObjectNumberFormat = &numberFormat
					}
					if headTailParameterMap, ok := helper.InterfaceToMap(transcodeTaskSetMap, "head_tail_parameter"); ok {
						headTailParameter := mps.HeadTailParameter{}
						if v, ok := headTailParameterMap["head_set"]; ok {
							for _, item := range v.([]interface{}) {
								headSetMap := item.(map[string]interface{})
								mediaInputInfo := mps.MediaInputInfo{}
								if v, ok := headSetMap["type"]; ok {
									mediaInputInfo.Type = helper.String(v.(string))
								}
								if cosInputInfoMap, ok := helper.InterfaceToMap(headSetMap, "cos_input_info"); ok {
									cosInputInfo := mps.CosInputInfo{}
									if v, ok := cosInputInfoMap["bucket"]; ok {
										cosInputInfo.Bucket = helper.String(v.(string))
									}
									if v, ok := cosInputInfoMap["region"]; ok {
										cosInputInfo.Region = helper.String(v.(string))
									}
									if v, ok := cosInputInfoMap["object"]; ok {
										cosInputInfo.Object = helper.String(v.(string))
									}
									mediaInputInfo.CosInputInfo = &cosInputInfo
								}
								if urlInputInfoMap, ok := helper.InterfaceToMap(headSetMap, "url_input_info"); ok {
									urlInputInfo := mps.UrlInputInfo{}
									if v, ok := urlInputInfoMap["url"]; ok {
										urlInputInfo.Url = helper.String(v.(string))
									}
									mediaInputInfo.UrlInputInfo = &urlInputInfo
								}
								headTailParameter.HeadSet = append(headTailParameter.HeadSet, &mediaInputInfo)
							}
						}
						if v, ok := headTailParameterMap["tail_set"]; ok {
							for _, item := range v.([]interface{}) {
								tailSetMap := item.(map[string]interface{})
								mediaInputInfo := mps.MediaInputInfo{}
								if v, ok := tailSetMap["type"]; ok {
									mediaInputInfo.Type = helper.String(v.(string))
								}
								if cosInputInfoMap, ok := helper.InterfaceToMap(tailSetMap, "cos_input_info"); ok {
									cosInputInfo := mps.CosInputInfo{}
									if v, ok := cosInputInfoMap["bucket"]; ok {
										cosInputInfo.Bucket = helper.String(v.(string))
									}
									if v, ok := cosInputInfoMap["region"]; ok {
										cosInputInfo.Region = helper.String(v.(string))
									}
									if v, ok := cosInputInfoMap["object"]; ok {
										cosInputInfo.Object = helper.String(v.(string))
									}
									mediaInputInfo.CosInputInfo = &cosInputInfo
								}
								if urlInputInfoMap, ok := helper.InterfaceToMap(tailSetMap, "url_input_info"); ok {
									urlInputInfo := mps.UrlInputInfo{}
									if v, ok := urlInputInfoMap["url"]; ok {
										urlInputInfo.Url = helper.String(v.(string))
									}
									mediaInputInfo.UrlInputInfo = &urlInputInfo
								}
								headTailParameter.TailSet = append(headTailParameter.TailSet, &mediaInputInfo)
							}
						}
						transcodeTaskInput.HeadTailParameter = &headTailParameter
					}
					mediaProcessTaskInput.TranscodeTaskSet = append(mediaProcessTaskInput.TranscodeTaskSet, &transcodeTaskInput)
				}
			}
			if v, ok := dMap["animated_graphic_task_set"]; ok {
				for _, item := range v.([]interface{}) {
					animatedGraphicTaskSetMap := item.(map[string]interface{})
					animatedGraphicTaskInput := mps.AnimatedGraphicTaskInput{}
					if v, ok := animatedGraphicTaskSetMap["definition"]; ok {
						animatedGraphicTaskInput.Definition = helper.IntUint64(v.(int))
					}
					if v, ok := animatedGraphicTaskSetMap["start_time_offset"]; ok {
						animatedGraphicTaskInput.StartTimeOffset = helper.Float64(v.(float64))
					}
					if v, ok := animatedGraphicTaskSetMap["end_time_offset"]; ok {
						animatedGraphicTaskInput.EndTimeOffset = helper.Float64(v.(float64))
					}
					if outputStorageMap, ok := helper.InterfaceToMap(animatedGraphicTaskSetMap, "output_storage"); ok {
						taskOutputStorage := mps.TaskOutputStorage{}
						if v, ok := outputStorageMap["type"]; ok {
							taskOutputStorage.Type = helper.String(v.(string))
						}
						if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
							cosOutputStorage := mps.CosOutputStorage{}
							if v, ok := cosOutputStorageMap["bucket"]; ok {
								cosOutputStorage.Bucket = helper.String(v.(string))
							}
							if v, ok := cosOutputStorageMap["region"]; ok {
								cosOutputStorage.Region = helper.String(v.(string))
							}
							taskOutputStorage.CosOutputStorage = &cosOutputStorage
						}
						animatedGraphicTaskInput.OutputStorage = &taskOutputStorage
					}
					if v, ok := animatedGraphicTaskSetMap["output_object_path"]; ok {
						animatedGraphicTaskInput.OutputObjectPath = helper.String(v.(string))
					}
					mediaProcessTaskInput.AnimatedGraphicTaskSet = append(mediaProcessTaskInput.AnimatedGraphicTaskSet, &animatedGraphicTaskInput)
				}
			}
			if v, ok := dMap["snapshot_by_time_offset_task_set"]; ok {
				for _, item := range v.([]interface{}) {
					snapshotByTimeOffsetTaskSetMap := item.(map[string]interface{})
					snapshotByTimeOffsetTaskInput := mps.SnapshotByTimeOffsetTaskInput{}
					if v, ok := snapshotByTimeOffsetTaskSetMap["definition"]; ok {
						snapshotByTimeOffsetTaskInput.Definition = helper.IntUint64(v.(int))
					}
					if v, ok := snapshotByTimeOffsetTaskSetMap["ext_time_offset_set"]; ok {
						extTimeOffsetSetSet := v.(*schema.Set).List()
						for i := range extTimeOffsetSetSet {
							extTimeOffsetSet := extTimeOffsetSetSet[i].(string)
							snapshotByTimeOffsetTaskInput.ExtTimeOffsetSet = append(snapshotByTimeOffsetTaskInput.ExtTimeOffsetSet, &extTimeOffsetSet)
						}
					}
					if v, _ := d.GetOk("time_offset_set"); v != nil {
						timeOffsetSetSet := v.(*schema.Set).List()
						for i := range timeOffsetSetSet {
							timeOffsetSet := timeOffsetSetSet[i].(float64)
							snapshotByTimeOffsetTaskInput.TimeOffsetSet = append(snapshotByTimeOffsetTaskInput.TimeOffsetSet, &timeOffsetSet)
						}
					}

					if v, ok := snapshotByTimeOffsetTaskSetMap["watermark_set"]; ok {
						for _, item := range v.([]interface{}) {
							watermarkSetMap := item.(map[string]interface{})
							watermarkInput := mps.WatermarkInput{}
							if v, ok := watermarkSetMap["definition"]; ok {
								watermarkInput.Definition = helper.IntUint64(v.(int))
							}
							if rawParameterMap, ok := helper.InterfaceToMap(watermarkSetMap, "raw_parameter"); ok {
								rawWatermarkParameter := mps.RawWatermarkParameter{}
								if v, ok := rawParameterMap["type"]; ok {
									rawWatermarkParameter.Type = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["coordinate_origin"]; ok {
									rawWatermarkParameter.CoordinateOrigin = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["x_pos"]; ok {
									rawWatermarkParameter.XPos = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["y_pos"]; ok {
									rawWatermarkParameter.YPos = helper.String(v.(string))
								}
								if imageTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "image_template"); ok {
									rawImageWatermarkInput := mps.RawImageWatermarkInput{}
									if imageContentMap, ok := helper.InterfaceToMap(imageTemplateMap, "image_content"); ok {
										mediaInputInfo := mps.MediaInputInfo{}
										if v, ok := imageContentMap["type"]; ok {
											mediaInputInfo.Type = helper.String(v.(string))
										}
										if cosInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "cos_input_info"); ok {
											cosInputInfo := mps.CosInputInfo{}
											if v, ok := cosInputInfoMap["bucket"]; ok {
												cosInputInfo.Bucket = helper.String(v.(string))
											}
											if v, ok := cosInputInfoMap["region"]; ok {
												cosInputInfo.Region = helper.String(v.(string))
											}
											if v, ok := cosInputInfoMap["object"]; ok {
												cosInputInfo.Object = helper.String(v.(string))
											}
											mediaInputInfo.CosInputInfo = &cosInputInfo
										}
										if urlInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "url_input_info"); ok {
											urlInputInfo := mps.UrlInputInfo{}
											if v, ok := urlInputInfoMap["url"]; ok {
												urlInputInfo.Url = helper.String(v.(string))
											}
											mediaInputInfo.UrlInputInfo = &urlInputInfo
										}
										rawImageWatermarkInput.ImageContent = &mediaInputInfo
									}
									if v, ok := imageTemplateMap["width"]; ok {
										rawImageWatermarkInput.Width = helper.String(v.(string))
									}
									if v, ok := imageTemplateMap["height"]; ok {
										rawImageWatermarkInput.Height = helper.String(v.(string))
									}
									if v, ok := imageTemplateMap["repeat_type"]; ok {
										rawImageWatermarkInput.RepeatType = helper.String(v.(string))
									}
									rawWatermarkParameter.ImageTemplate = &rawImageWatermarkInput
								}
								watermarkInput.RawParameter = &rawWatermarkParameter
							}
							if v, ok := watermarkSetMap["text_content"]; ok {
								watermarkInput.TextContent = helper.String(v.(string))
							}
							if v, ok := watermarkSetMap["svg_content"]; ok {
								watermarkInput.SvgContent = helper.String(v.(string))
							}
							if v, ok := watermarkSetMap["start_time_offset"]; ok {
								watermarkInput.StartTimeOffset = helper.Float64(v.(float64))
							}
							if v, ok := watermarkSetMap["end_time_offset"]; ok {
								watermarkInput.EndTimeOffset = helper.Float64(v.(float64))
							}
							snapshotByTimeOffsetTaskInput.WatermarkSet = append(snapshotByTimeOffsetTaskInput.WatermarkSet, &watermarkInput)
						}
					}
					if outputStorageMap, ok := helper.InterfaceToMap(snapshotByTimeOffsetTaskSetMap, "output_storage"); ok {
						taskOutputStorage := mps.TaskOutputStorage{}
						if v, ok := outputStorageMap["type"]; ok {
							taskOutputStorage.Type = helper.String(v.(string))
						}
						if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
							cosOutputStorage := mps.CosOutputStorage{}
							if v, ok := cosOutputStorageMap["bucket"]; ok {
								cosOutputStorage.Bucket = helper.String(v.(string))
							}
							if v, ok := cosOutputStorageMap["region"]; ok {
								cosOutputStorage.Region = helper.String(v.(string))
							}
							taskOutputStorage.CosOutputStorage = &cosOutputStorage
						}
						snapshotByTimeOffsetTaskInput.OutputStorage = &taskOutputStorage
					}
					if v, ok := snapshotByTimeOffsetTaskSetMap["output_object_path"]; ok {
						snapshotByTimeOffsetTaskInput.OutputObjectPath = helper.String(v.(string))
					}
					if objectNumberFormatMap, ok := helper.InterfaceToMap(snapshotByTimeOffsetTaskSetMap, "object_number_format"); ok {
						numberFormat := mps.NumberFormat{}
						if v, ok := objectNumberFormatMap["initial_value"]; ok {
							numberFormat.InitialValue = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["increment"]; ok {
							numberFormat.Increment = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["min_length"]; ok {
							numberFormat.MinLength = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["place_holder"]; ok {
							numberFormat.PlaceHolder = helper.String(v.(string))
						}
						snapshotByTimeOffsetTaskInput.ObjectNumberFormat = &numberFormat
					}
					mediaProcessTaskInput.SnapshotByTimeOffsetTaskSet = append(mediaProcessTaskInput.SnapshotByTimeOffsetTaskSet, &snapshotByTimeOffsetTaskInput)
				}
			}
			if v, ok := dMap["sample_snapshot_task_set"]; ok {
				for _, item := range v.([]interface{}) {
					sampleSnapshotTaskSetMap := item.(map[string]interface{})
					sampleSnapshotTaskInput := mps.SampleSnapshotTaskInput{}
					if v, ok := sampleSnapshotTaskSetMap["definition"]; ok {
						sampleSnapshotTaskInput.Definition = helper.IntUint64(v.(int))
					}
					if v, ok := sampleSnapshotTaskSetMap["watermark_set"]; ok {
						for _, item := range v.([]interface{}) {
							watermarkSetMap := item.(map[string]interface{})
							watermarkInput := mps.WatermarkInput{}
							if v, ok := watermarkSetMap["definition"]; ok {
								watermarkInput.Definition = helper.IntUint64(v.(int))
							}
							if rawParameterMap, ok := helper.InterfaceToMap(watermarkSetMap, "raw_parameter"); ok {
								rawWatermarkParameter := mps.RawWatermarkParameter{}
								if v, ok := rawParameterMap["type"]; ok {
									rawWatermarkParameter.Type = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["coordinate_origin"]; ok {
									rawWatermarkParameter.CoordinateOrigin = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["x_pos"]; ok {
									rawWatermarkParameter.XPos = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["y_pos"]; ok {
									rawWatermarkParameter.YPos = helper.String(v.(string))
								}
								if imageTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "image_template"); ok {
									rawImageWatermarkInput := mps.RawImageWatermarkInput{}
									if imageContentMap, ok := helper.InterfaceToMap(imageTemplateMap, "image_content"); ok {
										mediaInputInfo := mps.MediaInputInfo{}
										if v, ok := imageContentMap["type"]; ok {
											mediaInputInfo.Type = helper.String(v.(string))
										}
										if cosInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "cos_input_info"); ok {
											cosInputInfo := mps.CosInputInfo{}
											if v, ok := cosInputInfoMap["bucket"]; ok {
												cosInputInfo.Bucket = helper.String(v.(string))
											}
											if v, ok := cosInputInfoMap["region"]; ok {
												cosInputInfo.Region = helper.String(v.(string))
											}
											if v, ok := cosInputInfoMap["object"]; ok {
												cosInputInfo.Object = helper.String(v.(string))
											}
											mediaInputInfo.CosInputInfo = &cosInputInfo
										}
										if urlInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "url_input_info"); ok {
											urlInputInfo := mps.UrlInputInfo{}
											if v, ok := urlInputInfoMap["url"]; ok {
												urlInputInfo.Url = helper.String(v.(string))
											}
											mediaInputInfo.UrlInputInfo = &urlInputInfo
										}
										rawImageWatermarkInput.ImageContent = &mediaInputInfo
									}
									if v, ok := imageTemplateMap["width"]; ok {
										rawImageWatermarkInput.Width = helper.String(v.(string))
									}
									if v, ok := imageTemplateMap["height"]; ok {
										rawImageWatermarkInput.Height = helper.String(v.(string))
									}
									if v, ok := imageTemplateMap["repeat_type"]; ok {
										rawImageWatermarkInput.RepeatType = helper.String(v.(string))
									}
									rawWatermarkParameter.ImageTemplate = &rawImageWatermarkInput
								}
								watermarkInput.RawParameter = &rawWatermarkParameter
							}
							if v, ok := watermarkSetMap["text_content"]; ok {
								watermarkInput.TextContent = helper.String(v.(string))
							}
							if v, ok := watermarkSetMap["svg_content"]; ok {
								watermarkInput.SvgContent = helper.String(v.(string))
							}
							if v, ok := watermarkSetMap["start_time_offset"]; ok {
								watermarkInput.StartTimeOffset = helper.Float64(v.(float64))
							}
							if v, ok := watermarkSetMap["end_time_offset"]; ok {
								watermarkInput.EndTimeOffset = helper.Float64(v.(float64))
							}
							sampleSnapshotTaskInput.WatermarkSet = append(sampleSnapshotTaskInput.WatermarkSet, &watermarkInput)
						}
					}
					if outputStorageMap, ok := helper.InterfaceToMap(sampleSnapshotTaskSetMap, "output_storage"); ok {
						taskOutputStorage := mps.TaskOutputStorage{}
						if v, ok := outputStorageMap["type"]; ok {
							taskOutputStorage.Type = helper.String(v.(string))
						}
						if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
							cosOutputStorage := mps.CosOutputStorage{}
							if v, ok := cosOutputStorageMap["bucket"]; ok {
								cosOutputStorage.Bucket = helper.String(v.(string))
							}
							if v, ok := cosOutputStorageMap["region"]; ok {
								cosOutputStorage.Region = helper.String(v.(string))
							}
							taskOutputStorage.CosOutputStorage = &cosOutputStorage
						}
						sampleSnapshotTaskInput.OutputStorage = &taskOutputStorage
					}
					if v, ok := sampleSnapshotTaskSetMap["output_object_path"]; ok {
						sampleSnapshotTaskInput.OutputObjectPath = helper.String(v.(string))
					}
					if objectNumberFormatMap, ok := helper.InterfaceToMap(sampleSnapshotTaskSetMap, "object_number_format"); ok {
						numberFormat := mps.NumberFormat{}
						if v, ok := objectNumberFormatMap["initial_value"]; ok {
							numberFormat.InitialValue = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["increment"]; ok {
							numberFormat.Increment = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["min_length"]; ok {
							numberFormat.MinLength = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["place_holder"]; ok {
							numberFormat.PlaceHolder = helper.String(v.(string))
						}
						sampleSnapshotTaskInput.ObjectNumberFormat = &numberFormat
					}
					mediaProcessTaskInput.SampleSnapshotTaskSet = append(mediaProcessTaskInput.SampleSnapshotTaskSet, &sampleSnapshotTaskInput)
				}
			}
			if v, ok := dMap["image_sprite_task_set"]; ok {
				for _, item := range v.([]interface{}) {
					imageSpriteTaskSetMap := item.(map[string]interface{})
					imageSpriteTaskInput := mps.ImageSpriteTaskInput{}
					if v, ok := imageSpriteTaskSetMap["definition"]; ok {
						imageSpriteTaskInput.Definition = helper.IntUint64(v.(int))
					}
					if outputStorageMap, ok := helper.InterfaceToMap(imageSpriteTaskSetMap, "output_storage"); ok {
						taskOutputStorage := mps.TaskOutputStorage{}
						if v, ok := outputStorageMap["type"]; ok {
							taskOutputStorage.Type = helper.String(v.(string))
						}
						if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
							cosOutputStorage := mps.CosOutputStorage{}
							if v, ok := cosOutputStorageMap["bucket"]; ok {
								cosOutputStorage.Bucket = helper.String(v.(string))
							}
							if v, ok := cosOutputStorageMap["region"]; ok {
								cosOutputStorage.Region = helper.String(v.(string))
							}
							taskOutputStorage.CosOutputStorage = &cosOutputStorage
						}
						imageSpriteTaskInput.OutputStorage = &taskOutputStorage
					}
					if v, ok := imageSpriteTaskSetMap["output_object_path"]; ok {
						imageSpriteTaskInput.OutputObjectPath = helper.String(v.(string))
					}
					if v, ok := imageSpriteTaskSetMap["web_vtt_object_name"]; ok {
						imageSpriteTaskInput.WebVttObjectName = helper.String(v.(string))
					}
					if objectNumberFormatMap, ok := helper.InterfaceToMap(imageSpriteTaskSetMap, "object_number_format"); ok {
						numberFormat := mps.NumberFormat{}
						if v, ok := objectNumberFormatMap["initial_value"]; ok {
							numberFormat.InitialValue = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["increment"]; ok {
							numberFormat.Increment = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["min_length"]; ok {
							numberFormat.MinLength = helper.IntUint64(v.(int))
						}
						if v, ok := objectNumberFormatMap["place_holder"]; ok {
							numberFormat.PlaceHolder = helper.String(v.(string))
						}
						imageSpriteTaskInput.ObjectNumberFormat = &numberFormat
					}
					mediaProcessTaskInput.ImageSpriteTaskSet = append(mediaProcessTaskInput.ImageSpriteTaskSet, &imageSpriteTaskInput)
				}
			}
			if v, ok := dMap["adaptive_dynamic_streaming_task_set"]; ok {
				for _, item := range v.([]interface{}) {
					adaptiveDynamicStreamingTaskSetMap := item.(map[string]interface{})
					adaptiveDynamicStreamingTaskInput := mps.AdaptiveDynamicStreamingTaskInput{}
					if v, ok := adaptiveDynamicStreamingTaskSetMap["definition"]; ok {
						adaptiveDynamicStreamingTaskInput.Definition = helper.IntUint64(v.(int))
					}
					if v, ok := adaptiveDynamicStreamingTaskSetMap["watermark_set"]; ok {
						for _, item := range v.([]interface{}) {
							watermarkSetMap := item.(map[string]interface{})
							watermarkInput := mps.WatermarkInput{}
							if v, ok := watermarkSetMap["definition"]; ok {
								watermarkInput.Definition = helper.IntUint64(v.(int))
							}
							if rawParameterMap, ok := helper.InterfaceToMap(watermarkSetMap, "raw_parameter"); ok {
								rawWatermarkParameter := mps.RawWatermarkParameter{}
								if v, ok := rawParameterMap["type"]; ok {
									rawWatermarkParameter.Type = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["coordinate_origin"]; ok {
									rawWatermarkParameter.CoordinateOrigin = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["x_pos"]; ok {
									rawWatermarkParameter.XPos = helper.String(v.(string))
								}
								if v, ok := rawParameterMap["y_pos"]; ok {
									rawWatermarkParameter.YPos = helper.String(v.(string))
								}
								if imageTemplateMap, ok := helper.InterfaceToMap(rawParameterMap, "image_template"); ok {
									rawImageWatermarkInput := mps.RawImageWatermarkInput{}
									if imageContentMap, ok := helper.InterfaceToMap(imageTemplateMap, "image_content"); ok {
										mediaInputInfo := mps.MediaInputInfo{}
										if v, ok := imageContentMap["type"]; ok {
											mediaInputInfo.Type = helper.String(v.(string))
										}
										if cosInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "cos_input_info"); ok {
											cosInputInfo := mps.CosInputInfo{}
											if v, ok := cosInputInfoMap["bucket"]; ok {
												cosInputInfo.Bucket = helper.String(v.(string))
											}
											if v, ok := cosInputInfoMap["region"]; ok {
												cosInputInfo.Region = helper.String(v.(string))
											}
											if v, ok := cosInputInfoMap["object"]; ok {
												cosInputInfo.Object = helper.String(v.(string))
											}
											mediaInputInfo.CosInputInfo = &cosInputInfo
										}
										if urlInputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "url_input_info"); ok {
											urlInputInfo := mps.UrlInputInfo{}
											if v, ok := urlInputInfoMap["url"]; ok {
												urlInputInfo.Url = helper.String(v.(string))
											}
											mediaInputInfo.UrlInputInfo = &urlInputInfo
										}
										rawImageWatermarkInput.ImageContent = &mediaInputInfo
									}
									if v, ok := imageTemplateMap["width"]; ok {
										rawImageWatermarkInput.Width = helper.String(v.(string))
									}
									if v, ok := imageTemplateMap["height"]; ok {
										rawImageWatermarkInput.Height = helper.String(v.(string))
									}
									if v, ok := imageTemplateMap["repeat_type"]; ok {
										rawImageWatermarkInput.RepeatType = helper.String(v.(string))
									}
									rawWatermarkParameter.ImageTemplate = &rawImageWatermarkInput
								}
								watermarkInput.RawParameter = &rawWatermarkParameter
							}
							if v, ok := watermarkSetMap["text_content"]; ok {
								watermarkInput.TextContent = helper.String(v.(string))
							}
							if v, ok := watermarkSetMap["svg_content"]; ok {
								watermarkInput.SvgContent = helper.String(v.(string))
							}
							if v, ok := watermarkSetMap["start_time_offset"]; ok {
								watermarkInput.StartTimeOffset = helper.Float64(v.(float64))
							}
							if v, ok := watermarkSetMap["end_time_offset"]; ok {
								watermarkInput.EndTimeOffset = helper.Float64(v.(float64))
							}
							adaptiveDynamicStreamingTaskInput.WatermarkSet = append(adaptiveDynamicStreamingTaskInput.WatermarkSet, &watermarkInput)
						}
					}
					if outputStorageMap, ok := helper.InterfaceToMap(adaptiveDynamicStreamingTaskSetMap, "output_storage"); ok {
						taskOutputStorage := mps.TaskOutputStorage{}
						if v, ok := outputStorageMap["type"]; ok {
							taskOutputStorage.Type = helper.String(v.(string))
						}
						if cosOutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "cos_output_storage"); ok {
							cosOutputStorage := mps.CosOutputStorage{}
							if v, ok := cosOutputStorageMap["bucket"]; ok {
								cosOutputStorage.Bucket = helper.String(v.(string))
							}
							if v, ok := cosOutputStorageMap["region"]; ok {
								cosOutputStorage.Region = helper.String(v.(string))
							}
							taskOutputStorage.CosOutputStorage = &cosOutputStorage
						}
						adaptiveDynamicStreamingTaskInput.OutputStorage = &taskOutputStorage
					}
					if v, ok := adaptiveDynamicStreamingTaskSetMap["output_object_path"]; ok {
						adaptiveDynamicStreamingTaskInput.OutputObjectPath = helper.String(v.(string))
					}
					if v, ok := adaptiveDynamicStreamingTaskSetMap["sub_stream_object_name"]; ok {
						adaptiveDynamicStreamingTaskInput.SubStreamObjectName = helper.String(v.(string))
					}
					if v, ok := adaptiveDynamicStreamingTaskSetMap["segment_object_name"]; ok {
						adaptiveDynamicStreamingTaskInput.SegmentObjectName = helper.String(v.(string))
					}
					mediaProcessTaskInput.AdaptiveDynamicStreamingTaskSet = append(mediaProcessTaskInput.AdaptiveDynamicStreamingTaskSet, &adaptiveDynamicStreamingTaskInput)
				}
			}
			request.MediaProcessTask = &mediaProcessTaskInput
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "ai_content_review_task"); ok {
			aiContentReviewTaskInput := mps.AiContentReviewTaskInput{}
			if v, ok := dMap["definition"]; ok {
				aiContentReviewTaskInput.Definition = helper.IntUint64(v.(int))
			}
			request.AiContentReviewTask = &aiContentReviewTaskInput
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "ai_analysis_task"); ok {
			aiAnalysisTaskInput := mps.AiAnalysisTaskInput{}
			if v, ok := dMap["definition"]; ok {
				aiAnalysisTaskInput.Definition = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["extended_parameter"]; ok {
				aiAnalysisTaskInput.ExtendedParameter = helper.String(v.(string))
			}
			request.AiAnalysisTask = &aiAnalysisTaskInput
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "ai_recognition_task"); ok {
			aiRecognitionTaskInput := mps.AiRecognitionTaskInput{}
			if v, ok := dMap["definition"]; ok {
				aiRecognitionTaskInput.Definition = helper.IntUint64(v.(int))
			}
			request.AiRecognitionTask = &aiRecognitionTaskInput
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "task_notify_config"); ok {
			taskNotifyConfig := mps.TaskNotifyConfig{}
			if v, ok := dMap["cmq_model"]; ok {
				taskNotifyConfig.CmqModel = helper.String(v.(string))
			}
			if v, ok := dMap["cmq_region"]; ok {
				taskNotifyConfig.CmqRegion = helper.String(v.(string))
			}
			if v, ok := dMap["topic_name"]; ok {
				taskNotifyConfig.TopicName = helper.String(v.(string))
			}
			if v, ok := dMap["queue_name"]; ok {
				taskNotifyConfig.QueueName = helper.String(v.(string))
			}
			if v, ok := dMap["notify_mode"]; ok {
				taskNotifyConfig.NotifyMode = helper.String(v.(string))
			}
			if v, ok := dMap["notify_type"]; ok {
				taskNotifyConfig.NotifyType = helper.String(v.(string))
			}
			if v, ok := dMap["notify_url"]; ok {
				taskNotifyConfig.NotifyUrl = helper.String(v.(string))
			}
			request.TaskNotifyConfig = &taskNotifyConfig
		}

		if v, _ := d.GetOk("task_priority"); v != nil {
			request.TaskPriority = helper.IntInt64(v.(int))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().ResetWorkflow(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update mps workflow failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudMpsWorkflowRead(d, meta)
}

func resourceTencentCloudMpsWorkflowDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_workflow.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	workflowId := d.Id()

	if err := service.DeleteMpsWorkflowById(ctx, workflowId); err != nil {
		return err
	}

	return nil
}
