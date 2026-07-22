package mps

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMpsSchedules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMpsSchedulesRead,
		Schema: map[string]*schema.Schema{
			"schedule_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "IDs 的 schemes 到 查询. Array 长度 限制: 100。",
			},

			"trigger_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "触发器 类型 有效 值:`CosFileUpload`: scheme 是 triggered 当 文件 是 uploaded 到 Tencent Cloud Object Storage (COS).`AwsS3FileUpload`: scheme 是 triggered 当 文件 是 uploaded 到 AWS S3.如果 您 do 不 指定this 参数 或 leave 它 空，all schemes 将 是 返回 regardless 的 触发器 类型",
			},

			"status": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "scheme 状态 有效 值:`已启用`，`已禁用`. 如果 您 do 不 指定this 参数，all schemes 将 是 返回 regardless 的 状态",
			},

			"schedule_info_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "信息 的 schemes。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "scheme ID。",
						},
						"schedule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "scheme 名称注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "scheme 状态 有效 值:`已启用``已禁用`注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"trigger": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "触发器 的 scheme.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "触发器 类型 有效 值:`CosFileUpload`: Tencent Cloud COS 触发器.`AwsS3FileUpload`: AWS S3 触发器. Currently，此 类型 是 仅 支持 对于 transcoding tasks 和 schemes (不 支持 对于 workflows)。",
									},
									"cos_file_upload_trigger": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "此 参数 为必填项 和 有效 当 `类型` 是 `CosFileUpload`，indicating COS 触发器 规则.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bucket": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 COS 存储桶 bound 到 工作流，such 作为 `TopRankVideo-125xxx88`。",
												},
												"region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "地域 的 COS 存储桶 bound 到 工作流，such 作为 `ap-chongiqng`。",
												},
												"dir": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Input 路径 directory bound 到 工作流，such 作为 `/movie/201907/`. 如果此参数为空， `/` root directory 将 是 使用。",
												},
												"formats": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "格式 列表 files 该 可以 触发器 工作流，such 作为 [mp4，flv，mov]. 如果此参数为空，files 在 all formats 可以 触发器 工作流。",
												},
											},
										},
									},
									"aws_s3_file_upload_trigger": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "AWS S3 触发器. 此 参数 是 有效 和 必填 如果 `类型` 是 `AwsS3FileUpload`.注意: Currently， 键 对于 AWS S3 存储桶， 触发器 SQS queue，和 callback SQS queue 必须 是 same.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"s3_bucket": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "AWS S3 存储桶 bound 到 scheme。",
												},
												"s3_region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "地域 的 AWS S3 存储桶",
												},
												"dir": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "存储桶 directory bound. It 必须 是 absolute 路径 该 starts 和 结束 使用 `/`，such 作为 `/movie/201907/`. 如果 您 do 不 指定this， root directory 将 是 bound.	。",
												},
												"formats": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "文件 formats 该 将 触发器 scheme，such 作为 [mp4，flv，mov]. 如果 您 do 不 指定this， upload 的 files 在 any 格式 将 触发器 scheme.	。",
												},
												"s3_secret_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "键 ID AWS S3 存储桶注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"s3_secret_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "键 的 AWS S3 存储桶注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"aws_sqs": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "SQS queue 的 AWS S3 存储桶Note: queue 必须 是 在 same 地域 作为 存储桶注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"sqs_region": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "地域 的 SQS queue。",
															},
															"sqs_queue_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "名称 SQS queue。",
															},
															"s3_secret_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "键 ID 必填 到 read 从/write 到 SQS queue。",
															},
															"s3_secret_key": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "键 必填 到 read 从/write 到 SQS queue。",
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
						"activities": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "subtasks 的 scheme.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"activity_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "subtask 类型`input`: start.`output`: end.`操作-trans`: Transcoding.`操作-samplesnapshot`: Sampled screencapturing.`操作-AIAnalysis`: 内容 analysis.`操作-AIRecognition`: 内容 recognition.`操作-aiReview`: 内容 moderation.`操作-animated-graphics`: Animated screenshot generation.`操作-镜像-sprite`: Image sprite generation.`操作-snapshotByTimeOffset`: Time point screencapturing.`操作-adaptive-substream`: Adaptive bitrate streaming.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"reardrive_index": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
										Computed:    true,
										Description: "indexes 的 subsequent actions.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"activity_para": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "参数 的 subtask.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"transcode_task": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "A transcoding 任务。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "ID 视频 transcoding template。",
															},
															"raw_parameter": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Custom 视频 transcoding 参数，其中 是 有效 如果 `Definition` 是 0.此 参数 是 使用 在 highly customized scenarios. We recommend 您 使用 `Definition` 到 指定transcoding 参数 preferably。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"container": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Container. 有效值：mp4; flv; hls; mp3; flac; ogg; m4a. Among them，mp3，flac，ogg，和 m4a 是 对于 音频 files。",
																		},
																		"remove_video": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "是否remove 视频 数据. 有效 值:0: retain;1: remove.默认值：0。",
																		},
																		"remove_audio": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "是否remove 音频 数据. 有效 值:0: retain;1: remove.默认值：0。",
																		},
																		"video_template": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Video 流 配置 参数. 此 字段 为必填项 当 `RemoveVideo` 是 0。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"codec": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "视频 codec. 有效 值:`libx264`: H.264`libx265`: H.265`av1`: AOMedia Video 1Note: You 必须 指定a resolution (不 higher 比 640 x 480) 如果 H.265 codec 是 使用.注意: You 可以 仅 使用 AOMedia Video 1 codec 对于 MP4 files。",
																					},
																					"fps": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "视频 frame 速率 (Hz). 取值范围：[0，100].如果 值 是 0， frame 速率 将 是 same 作为 该 的 来源 视频.注意: For adaptive bitrate streaming， 值 范围 的 此 参数 是 [0，60]。",
																					},
																					"bitrate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "视频 bitrate (Kbps). 取值范围：0 和 [128，35000].如果 值 是 0， bitrate 的 视频 将 是 same 作为 该 的 来源 视频。",
																					},
																					"resolution_adaptive": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Resolution adaption. 有效 值:open: 已启用 当 resolution adaption 是 已启用，`宽度` 表示long side 的 视频，while `高度` 表示short side.close: 已禁用 当 resolution adaption 是 已禁用，`宽度` 表示width 的 视频，while `高度` 表示height.默认值：open.注意: 当 resolution adaption 是 已启用，`宽度` 不能 是 smaller 比 `高度`。",
																					},
																					"width": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Maximum 值 的 宽度 (或 long side) 的 视频 流 （像素）。 取值范围：0 和 [128，4,096].如果 both `宽度` 和 `高度` 是 0， resolution 将 是 same 作为 该 的 来源 视频;如果 `宽度` 是 0，但 `高度` 是 不 0，`宽度` 将 是 proportionally scaled;如果 `宽度` 是 不 0，但 `高度` 是 0，`高度` 将 是 proportionally scaled;如果 both `宽度` 和 `高度` 是 不 0， 自定义 resolution 将 是 使用.默认值：0。",
																					},
																					"height": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Maximum 值 的 高度 (或 short side) 的 视频 流 （像素）。 取值范围：0 和 [128，4,096].如果 both `宽度` 和 `高度` 是 0， resolution 将 是 same 作为 该 的 来源 视频;如果 `宽度` 是 0，但 `高度` 是 不 0，`宽度` 将 是 proportionally scaled;如果 `宽度` 是 不 0，但 `高度` 是 0，`高度` 将 是 proportionally scaled;如果 both `宽度` 和 `高度` 是 不 0， 自定义 resolution 将 是 使用.默认值：0。",
																					},
																					"gop": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Frame 间隔 between I keyframes. 取值范围：0 和 [1,100000].如果 此 参数 是 0 或 left 空， 系统 将 automatically 集合 GOP 长度。",
																					},
																					"fill_type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "fill 模式，其中 表示how 视频 是 resized 当 视频&#39;s original aspect ratio 是 different 从 目标 aspect ratio. 有效 值:stretch: Stretch 镜像 frame 通过 frame 到 fill entire screen. 视频 镜像 可能 become squashed 或 stretched after transcoding.black: Keep 镜像&#39;s original aspect ratio 和 fill blank space 使用 black bars.white: Keep 镜像&#39;s original aspect ratio 和 fill blank space 使用 white bars.gauss: Keep 镜像&#39;s original aspect ratio 和 apply Gaussian blur 到 blank space.默认值：black.注意: Only `stretch` 和 `black` 是 支持 对于 adaptive bitrate streaming。",
																					},
																					"vcrf": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "control factor 的 视频 constant bitrate. 取值范围：[1，51]如果 此 参数 是 指定，CRF ( bitrate control 方法) 将 是 用于transcoding. (Video bitrate 将 无 longer take effect.)It 是 不 recommended 到 指定this 参数 如果 there 是 无 special requirements。",
																					},
																				},
																			},
																		},
																		"audio_template": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Audio 流 配置 参数. 此 字段 为必填项 当 `RemoveAudio` 是 0。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"codec": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Audio 流 codec.当 outer `Container` 参数 是 `mp3`， 有效 值 是:libmp3lame.当 outer `Container` 参数 是 `ogg` 或 `flac`， 有效 值 是:flac.当 outer `Container` 参数 是 `m4a`， 有效 值 include:libfdk_aac;libmp3lame;ac3.当 outer `Container` 参数 是 `mp4` 或 `flv`， 有效 值 include:libfdk_aac: more suitable 对于 mp4;libmp3lame: more suitable 对于 flv.当 outer `Container` 参数 是 `hls`， 有效 值 include:libfdk_aac;libmp3lame。",
																					},
																					"bitrate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Audio 流 bitrate 在 Kbps. 取值范围：0 和 [26，256].如果 值 是 0， bitrate 的 音频 流 将 是 same 作为 该 的 original 音频。",
																					},
																					"sample_rate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Audio 流 sample 速率. 有效 值:32,00044,10048,000In Hz。",
																					},
																					"audio_channel": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Audio channel 系统. 有效 值:1: Mono2: Dual6: StereoWhen media 是 packaged 在 音频 格式 (FLAC，OGG，MP3，M4A)， sound channel 不能 是 集合 到 stereo.默认值：2。",
																					},
																				},
																			},
																		},
																		"tehd_config": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "TESHD transcoding 参数。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "TESHD 类型 有效 值:`TEHD-100`: TESHD-100. 如果此参数为空，TESHD 将 不 是 已启用",
																					},
																					"max_video_bitrate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Maximum bitrate，其中 是 有效 当 `类型` 是 `TESHD`. 如果此参数为空 或 0 是 entered，there 将 是 无 upper 限制 对于 bitrate。",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"override_parameter": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Video transcoding 自定义 参数，其中 是 有效 当 `Definition` 是 不 0.当 any 参数 在 此 structure 是 entered，they 将 是 用于override corresponding 参数 在 templates.此 参数 是 使用 在 highly customized scenarios. We recommend 您 仅 使用 `Definition` 到 指定transcoding 参数.注意: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 是 found。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"container": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Container 格式 有效值：mp4，flv，hls，mp3，flac，ogg，和 m4a; mp3，flac，ogg，和 m4a 是 formats 的 音频 files。",
																		},
																		"remove_video": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "是否remove 视频 数据. 有效 值:0: retain1: remove。",
																		},
																		"remove_audio": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "是否remove 音频 数据. 有效 值:0: retain1: remove。",
																		},
																		"video_template": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Video 流 配置 参数。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"codec": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "视频 codec. 有效 值:libx264: H.264libx265: H.265av1: AOMedia Video 1Note: You 必须 指定a resolution (不 higher 比 640 x 480) 如果 H.265 codec 是 使用.注意: You 可以 仅 使用 AOMedia Video 1 codec 对于 MP4 files。",
																					},
																					"fps": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Video frame 速率 在 Hz. 取值范围：[0，100].如果 值 是 0， frame 速率 将 是 same 作为 该 的 来源 视频。",
																					},
																					"bitrate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Bitrate 的 视频 流 在 Kbps. 取值范围：0 和 [128，35,000].如果 值 是 0， bitrate 的 视频 将 是 same 作为 该 的 来源 视频。",
																					},
																					"resolution_adaptive": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Resolution adaption. 有效 值:open: 已启用 当 resolution adaption 是 已启用，`宽度` 表示long side 的 视频，while `高度` 表示short side.close: 已禁用 当 resolution adaption 是 已禁用，`宽度` 表示width 的 视频，while `高度` 表示height.注意: 当 resolution adaption 是 已启用，`宽度` 不能 是 smaller 比 `高度`。",
																					},
																					"width": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Maximum 值 的 宽度 (或 long side) 的 视频 流 （像素）。 取值范围：0 和 [128，4,096].如果 both `宽度` 和 `高度` 是 0， resolution 将 是 same 作为 该 的 来源 视频;如果 `宽度` 是 0，但 `高度` 是 不 0，`宽度` 将 是 proportionally scaled;如果 `宽度` 是 不 0，但 `高度` 是 0，`高度` 将 是 proportionally scaled;如果 both `宽度` 和 `高度` 是 不 0， 自定义 resolution 将 是 使用。",
																					},
																					"height": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Maximum 值 的 高度 (或 short side) 的 视频 流 （像素）。 取值范围：0 和 [128，4,096]。",
																					},
																					"gop": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Frame 间隔 between I keyframes. 取值范围：0 和 [1,100000]. 如果 此 参数 是 0， 系统 将 automatically 集合 GOP 长度。",
																					},
																					"fill_type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Fill 类型 Fill refers 到 way 的 processing screenshot 当 its aspect ratio 是 different 从 该 的 来源 视频. following fill types 是 支持: stretch: stretch. screenshot 将 是 stretched frame 通过 frame 到 match aspect ratio 的 来源 视频，其中 可能 make screenshot shorter 或 longer;black: fill 使用 black. 此 选项 retains aspect ratio 的 来源 视频 对于 screenshot 和 fills unmatched area 使用 black color blocks.white: fill 使用 white. 此 选项 retains aspect ratio 的 来源 视频 对于 screenshot 和 fills unmatched area 使用 white color blocks.gauss: fill 使用 Gaussian blur. 此 选项 retains aspect ratio 的 来源 视频 对于 screenshot 和 fills unmatched area 使用 Gaussian blur。",
																					},
																					"vcrf": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "control factor 的 视频 constant bitrate. 取值范围：[0，51]. 此 参数 将 是 已禁用 如果 您 enter `0`.It 是 不 recommended 到 指定this 参数 如果 there 是 无 special requirements。",
																					},
																					"content_adapt_stream": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "是否enable adaptive 编码. 有效 值:0: Disable1: Enable默认值：0. 如果 此 参数 是 集合 到 `1`，多个 streams 使用 different resolutions 和 bitrates 将 是 generated automatically. highest resolution，bitrate，和 quality 的 streams 是 determined 通过 值 的 `宽度` 和 `高度`，`Bitrate`，和 `Vcrf` 在 `VideoTemplate` respectively. 如果 these 参数 是 不 集合 在 `VideoTemplate`， highest resolution generated 将 是 same 作为 该 的 来源 视频，和 highest 视频 quality 将 是 close 到 VMAF 95. To 使用 此 参数 或 learn about billing details 的 adaptive 编码，please contact your sales rep。",
																					},
																				},
																			},
																		},
																		"audio_template": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Audio 流 配置 参数。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"codec": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Audio 流 codec.当 outer `Container` 参数 是 `mp3`， 有效 值 是:libmp3lame.当 outer `Container` 参数 是 `ogg` 或 `flac`， 有效 值 是:flac.当 outer `Container` 参数 是 `m4a`， 有效 值 include:libfdk_aac;libmp3lame;ac3.当 outer `Container` 参数 是 `mp4` 或 `flv`， 有效 值 include:libfdk_aac: More suitable 对于 mp4;libmp3lame: More suitable 对于 flv;mp2.当 outer `Container` 参数 是 `hls`， 有效 值 include:libfdk_aac;libmp3lame。",
																					},
																					"bitrate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Audio 流 bitrate 在 Kbps. 取值范围：0 和 [26，256]. 如果 值 是 0， bitrate 的 音频 流 将 是 same 作为 该 的 original 音频。",
																					},
																					"sample_rate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Audio 流 sample 速率. 有效 值:32,00044,10048,000In Hz。",
																					},
																					"audio_channel": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Audio channel 系统. 有效 值:1: Mono2: Dual6: StereoWhen media 是 packaged 在 音频 格式 (FLAC，OGG，MP3，M4A)， sound channel 不能 是 集合 到 stereo。",
																					},
																					"stream_selects": {
																						Type: schema.TypeSet,
																						Elem: &schema.Schema{
																							Type: schema.TypeInt,
																						},
																						Computed:    true,
																						Description: "音频 tracks 到 retain. All 音频 tracks 是 retained 通过 默认值。",
																					},
																				},
																			},
																		},
																		"tehd_config": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "TSC transcoding 参数.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "TSC 类型 有效 值:`TEHD-100`: TSC-100 (视频 TSC). `TEHD-200`: TSC-200 (音频 TSC). 如果 此 参数 是 left blank，无 modification 将 是 made.注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																					"max_video_bitrate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "最大 视频 bitrate. 如果 此 参数 是 不 指定，无 modifications 将 是 made.注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																				},
																			},
																		},
																		"subtitle_template": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "subtitle settings.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"path": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "URL 的 subtitles 到 add 到 视频.注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																					"stream_index": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "subtitle track 到 add 到 视频. 如果 both `路径` 和 `StreamIndex` 是 指定，`路径` 将 是 使用. You need 到 指定at least 一个 的 two 参数.注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																					"font_type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "font. 有效 值:`hei.ttf`: Heiti.`song.ttf`: Songti.`simkai.ttf`: Kaiti.`arial.ttf`: Arial. 默认为 `hei.ttf`.注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																					"font_size": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "font 大小 (pixels). 如果 此 是 不 指定， font 大小 在 subtitle 文件 将 是 使用.注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																					"font_color": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "font color 在 0xRRGGBB 格式 默认值：0xFFFFFF (white).注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																					"font_alpha": {
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "text transparency. 取值范围：0-1.`0`: Fully transparent.`1`: Fully opaque.默认值：1.注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																				},
																			},
																		},
																		"addon_audio_stream": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "信息 的 外部 音频 track 到 add.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "input 类型 有效 值:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "信息 的 COS 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `COS`。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "COS 存储桶 的 对象 到 process，such 作为 `TopRankVideo-125xxx88`。",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "地域 的 COS 存储桶，such 作为 `ap-chongqing`。",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "路径 的 对象 到 process，such 作为 `/movie/201907/WildAnimal.mov`。",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "URL 的 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "URL 的 视频。",
																								},
																							},
																						},
																					},
																					"s3_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "信息 的 AWS S3 对象 processed. 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"s3_bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "AWS S3 存储桶",
																								},
																								"s3_region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "地域 的 AWS S3 存储桶",
																								},
																								"s3_object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "路径 的 AWS S3 对象。",
																								},
																								"s3_secret_id": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "键 ID 必填 到 访问 AWS S3 对象。",
																								},
																								"s3_secret_key": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "键 必填 到 访问 AWS S3 对象。",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"std_ext_info": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "An extended 字段 对于 transcoding.注意：此字段可能返回 null，表示无法获取有效值。",
																		},
																		"add_on_subtitles": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "subtitle 文件 到 add.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "模式 有效 值:`subtitle-流`: Add subtitle track.`close-caption-708`: Embed EA-708 subtitles 在 SEI frames.`close-caption-608`: Embed CEA-608 subtitles 在 SEI frames.注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																					"subtitle": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "subtitle 文件.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"type": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "input 类型 有效 值:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
																								},
																								"cos_input_info": {
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "信息 的 COS 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `COS`。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"bucket": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "COS 存储桶 的 对象 到 process，such 作为 `TopRankVideo-125xxx88`。",
																											},
																											"region": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "地域 的 COS 存储桶，such 作为 `ap-chongqing`。",
																											},
																											"object": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "路径 的 对象 到 process，such 作为 `/movie/201907/WildAnimal.mov`。",
																											},
																										},
																									},
																								},
																								"url_input_info": {
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "URL 的 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"url": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "URL 的 视频。",
																											},
																										},
																									},
																								},
																								"s3_input_info": {
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "信息 的 AWS S3 对象 processed. 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"s3_bucket": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "AWS S3 存储桶",
																											},
																											"s3_region": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "地域 的 AWS S3 存储桶",
																											},
																											"s3_object": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "路径 的 AWS S3 对象。",
																											},
																											"s3_secret_id": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "键 ID 必填 到 访问 AWS S3 对象。",
																											},
																											"s3_secret_key": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "键 必填 到 访问 AWS S3 对象。",
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
															"watermark_set": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "列表 up 到 10 镜像 或 text watermarks.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"definition": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "ID watermarking template。",
																		},
																		"raw_parameter": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Custom 水印 参数，其中 是 有效 如果 `Definition` 是 0.此 参数 是 使用 在 highly customized scenarios. We recommend 您 使用 `Definition` 到 指定watermark 参数 preferably.Custom 水印 参数 是 不 可用 对于 screenshot。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Watermark 类型 有效 值:镜像: 镜像 水印。",
																					},
																					"coordinate_origin": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Origin position，其中 currently 可以 仅 是:TopLeft: 源站 的 coordinates 是 在 top-left corner 的 视频，和 源站 的 水印 是 在 top-left corner 的 镜像 或 text.默认值：TopLeft。",
																					},
																					"x_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "horizontal position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `XPos` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `XPos` 是 10% 的 视频 宽度;如果 字符串 结束 在 像素， `XPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `XPos` 是 100 像素.默认值：0 像素。",
																					},
																					"y_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "vertical position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `YPos` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `YPos` 是 10% 的 视频 高度;如果 字符串 结束 在 像素， `YPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `YPos` 是 100 像素.默认值：0 像素。",
																					},
																					"image_template": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Image 水印 template. 此 字段 为必填项 当 `类型` 是 `镜像` 和 是 无效 当 `类型` 是 `text`。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"image_content": {
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Input 内容 的 水印 镜像. JPEG 和 PNG images 是 支持。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"type": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "input 类型 有效 值:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
																											},
																											"cos_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "信息 的 COS 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `COS`。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "COS 存储桶 的 对象 到 process，such 作为 `TopRankVideo-125xxx88`。",
																														},
																														"region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "地域 的 COS 存储桶，such 作为 `ap-chongqing`。",
																														},
																														"object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "路径 的 对象 到 process，such 作为 `/movie/201907/WildAnimal.mov`。",
																														},
																													},
																												},
																											},
																											"url_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "URL 的 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"url": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "URL 的 视频。",
																														},
																													},
																												},
																											},
																											"s3_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "信息 的 AWS S3 对象 processed. 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"s3_bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "AWS S3 存储桶",
																														},
																														"s3_region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "地域 的 AWS S3 存储桶",
																														},
																														"s3_object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "路径 的 AWS S3 对象。",
																														},
																														"s3_secret_id": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "键 ID 必填 到 访问 AWS S3 对象。",
																														},
																														"s3_secret_key": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "键 必填 到 访问 AWS S3 对象。",
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																								"width": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Watermark 宽度. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `宽度` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `宽度` 是 10% 的 视频 宽度;如果 字符串 结束 在 像素， `宽度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `宽度` 是 100 像素.默认值：10%。",
																								},
																								"height": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Watermark 高度. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `高度` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `高度` 是 10% 的 视频 高度;如果 字符串 结束 在 像素， `高度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `高度` 是 100 像素.默认值：0 像素，其中 表示 该 `高度` 将 是 proportionally scaled according 到 aspect ratio 的 original 水印 镜像。",
																								},
																								"repeat_type": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Repeat 类型 animated 水印. 有效 值:`once`: 无 longer appears after 水印 playback 结束.`repeat_last_frame`: stays 在 last frame after 水印 playback 结束.`repeat` (默认值): repeats playback until 视频 结束。",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"text_content": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Text 内容 的 up 到 100 字符. 此 字段 为必填项 仅 当 水印 类型 是 text.Text 水印 是 不 可用 对于 screenshot。",
																		},
																		"svg_content": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "SVG 内容 的 up 到 2,000,000 字符. 此 字段 为必填项 仅 当 水印 类型 是 `SVG`.SVG 水印 是 不 可用 对于 screenshot。",
																		},
																		"start_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "开始时间 偏移量 的 水印 （秒）。 如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame.如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame;如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 appear 在 second n after first 视频 frame;如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 appear 在 second n before last 视频 frame。",
																		},
																		"end_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "结束时间 偏移量 的 水印 （秒）。如果此参数为空 或 0 是 entered， 水印 将 exist till last 视频 frame;如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 exist till second n;如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 exist till second n before last 视频 frame。",
																		},
																	},
																},
															},
															"mosaic_set": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "列表 blurs. Up 到 10 ones 可以 是 支持。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"coordinate_origin": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Origin position，其中 currently 可以 仅 是:TopLeft: 源站 的 coordinates 是 在 top-left corner 的 视频，和 源站 的 blur 是 在 top-left corner 的 镜像 或 text.默认值：TopLeft。",
																		},
																		"x_pos": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "horizontal position 的 源站 的 blur relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `XPos` 的 blur 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `XPos` 是 10% 的 视频 宽度;如果 字符串 结束 在 像素， `XPos` 的 blur 将 是 指定 像素; 对于 示例，`100px` 表示 该 `XPos` 是 100 像素.默认值：0 像素。",
																		},
																		"y_pos": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Vertical position 的 源站 的 blur relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `YPos` 的 blur 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `YPos` 是 10% 的 视频 高度;如果 字符串 结束 在 像素， `YPos` 的 blur 将 是 指定 像素; 对于 示例，`100px` 表示 该 `YPos` 是 100 像素.默认值：0 像素。",
																		},
																		"width": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Blur 宽度. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `宽度` 的 blur 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `宽度` 是 10% 的 视频 宽度;如果 字符串 结束 在 像素， `宽度` 的 blur 将 是 在 像素; 对于 示例，`100px` 表示 该 `宽度` 是 100 像素.默认值：10%。",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Blur 高度. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `高度` 的 blur 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `高度` 是 10% 的 视频 高度;如果 字符串 结束 在 像素， `高度` 的 blur 将 是 在 像素; 对于 示例，`100px` 表示 该 `高度` 是 100 像素.默认值：10%。",
																		},
																		"start_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "开始时间 偏移量 的 blur （秒）。 如果此参数为空 或 0 是 entered， blur 将 appear upon first 视频 frame.如果此参数为空 或 0 是 entered， blur 将 appear upon first 视频 frame;如果 此 值 是 greater 比 0 (e.g.，n)， blur 将 appear 在 second n after first 视频 frame;如果 此 值 是 smaller 比 0 (e.g.，-n)， blur 将 appear 在 second n before last 视频 frame。",
																		},
																		"end_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "结束时间 偏移量 的 blur （秒）。如果此参数为空 或 0 是 entered， blur 将 exist till last 视频 frame;如果 此 值 是 greater 比 0 (e.g.，n)， blur 将 exist till second n;如果 此 值 是 smaller 比 0 (e.g.，-n)， blur 将 exist till second n before last 视频 frame。",
																		},
																	},
																},
															},
															"start_time_offset": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "开始时间 偏移量 的 transcoded 视频，（秒）。如果此参数为空 或 集合 到 0， transcoded 视频 将 start 在 same 时间 作为 original 视频.如果 此 参数 是 集合 到 positive 数量 (n 对于 示例)， transcoded 视频 将 start 在 nth second 的 original 视频.如果 此 参数 是 集合 到 negative 数量 (-n 对于 示例)， transcoded 视频 将 start 在 nth second before end 的 original 视频。",
															},
															"end_time_offset": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "结束时间 偏移量 的 transcoded 视频，（秒）。如果此参数为空 或 集合 到 0， transcoded 视频 将 end 在 same 时间 作为 original 视频.如果 此 参数 是 集合 到 positive 数量 (n 对于 示例)， transcoded 视频 将 end 在 nth second 的 original 视频.如果 此 参数 是 集合 到 negative 数量 (-n 对于 示例)， transcoded 视频 将 end 在 nth second before end 的 original 视频。",
															},
															"output_storage": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Target 存储桶 的 输出文件 如果此参数为空， `OutputStorage` 值 的 upper 文件夹 将 是 inherited.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "存储 类型 对于 media processing 输出文件 有效 值:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
																		},
																		"cos_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "location 到 save output 对象 在 COS. 此 参数 是 有效 和 必填 当 `类型` 是 COS.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "存储桶 到 其中 输出文件 的 media processing 是 saved，such 作为 `TopRankVideo-125xxx88`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "地域 的 输出存储桶，such 作为 `ap-chongqing`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																					},
																				},
																			},
																		},
																		"s3_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "AWS S3 存储桶 到 save 输出文件 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "AWS S3 存储桶",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "地域 的 AWS S3 存储桶",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "键 ID 必填 到 upload files 到 AWS S3 对象。",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "键 必填 到 upload files 到 AWS S3 对象。",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"output_object_path": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "路径 到 primary 输出文件，其中 可以 是 relative 路径 或 absolute 路径 如果此参数为空， following relative 路径 将 是 使用 通过 默认值：`{inputName}_transcode_{definition}.{格式}`。",
															},
															"segment_object_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "路径 到 输出文件 part ( 路径 到 ts during transcoding 到 HLS)，其中 可以 仅 是 relative 路径 如果此参数为空， following relative 路径 将 是 使用 通过 默认值：`{inputName}_transcode_{definition}_{数量}.{格式}`。",
															},
															"object_number_format": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Rule 的 `{数量}` variable 在 输出路径 after transcoding.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"initial_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Start 值 的 `{数量}` variable. 默认值：0。",
																		},
																		"increment": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Increment 的 `{数量}` variable. 默认值：1。",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "最小长度the `{数量}` variable. A placeholder 将 是 使用 如果 variable 长度 是 below 最小 requirement. 默认值：1。",
																		},
																		"place_holder": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Placeholder 使用 当 `{数量}` variable 长度 是 below 最小 requirement. 默认值：0。",
																		},
																	},
																},
															},
															"head_tail_parameter": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Opening 和 closing credits parametersNote: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 是 found。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"head_set": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Opening credits 列表。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "input 类型 有效 值:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "信息 的 COS 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `COS`。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "COS 存储桶 的 对象 到 process，such 作为 `TopRankVideo-125xxx88`。",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "地域 的 COS 存储桶，such 作为 `ap-chongqing`。",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "路径 的 对象 到 process，such 作为 `/movie/201907/WildAnimal.mov`。",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "URL 的 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "URL 的 视频。",
																								},
																							},
																						},
																					},
																					"s3_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "信息 的 AWS S3 对象 processed. 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"s3_bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "AWS S3 存储桶",
																								},
																								"s3_region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "地域 的 AWS S3 存储桶",
																								},
																								"s3_object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "路径 的 AWS S3 对象。",
																								},
																								"s3_secret_id": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "键 ID 必填 到 访问 AWS S3 对象。",
																								},
																								"s3_secret_key": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "键 必填 到 访问 AWS S3 对象。",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"tail_set": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Closing credits 列表。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "input 类型 有效 值:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "信息 的 COS 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `COS`。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "COS 存储桶 的 对象 到 process，such 作为 `TopRankVideo-125xxx88`。",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "地域 的 COS 存储桶，such 作为 `ap-chongqing`。",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "路径 的 对象 到 process，such 作为 `/movie/201907/WildAnimal.mov`。",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "URL 的 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "URL 的 视频。",
																								},
																							},
																						},
																					},
																					"s3_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "信息 的 AWS S3 对象 processed. 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"s3_bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "AWS S3 存储桶",
																								},
																								"s3_region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "地域 的 AWS S3 存储桶",
																								},
																								"s3_object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "路径 的 AWS S3 对象。",
																								},
																								"s3_secret_id": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "键 ID 必填 到 访问 AWS S3 对象。",
																								},
																								"s3_secret_key": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "键 必填 到 访问 AWS S3 对象。",
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
												"animated_graphic_task": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "An animated screenshot generation 任务。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
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
															"output_storage": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Target 存储桶 的 generated animated 镜像 文件. 如果此参数为空， `OutputStorage` 值 的 upper 文件夹 将 是 inherited.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "存储 类型 对于 media processing 输出文件 有效 值:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
																		},
																		"cos_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "location 到 save output 对象 在 COS. 此 参数 是 有效 和 必填 当 `类型` 是 COS.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "存储桶 到 其中 输出文件 的 media processing 是 saved，such 作为 `TopRankVideo-125xxx88`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "地域 的 输出存储桶，such 作为 `ap-chongqing`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																					},
																				},
																			},
																		},
																		"s3_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "AWS S3 存储桶 到 save 输出文件 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "AWS S3 存储桶",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "地域 的 AWS S3 存储桶",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "键 ID 必填 到 upload files 到 AWS S3 对象。",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "键 必填 到 upload files 到 AWS S3 对象。",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"output_object_path": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "输出路径 到 generated animated 镜像 文件，其中 可以 是 relative 路径 或 absolute 路径 如果此参数为空， following relative 路径 将 是 使用 通过 默认值：`{inputName}_animatedGraphic_{definition}.{格式}`。",
															},
														},
													},
												},
												"snapshot_by_time_offset_task": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "A 时间 point screencapturing 任务。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "ID 时间 point screencapturing template。",
															},
															"ext_time_offset_set": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Computed:    true,
																Description: "列表 screenshot 时间 points 在 格式 的 `s` 或 `%`:如果 字符串 结束 在 `s`，它 表示 该 时间 point 是 在 秒; 对于 示例，`3.5s` 表示 该 时间 point 是 3.5th second;如果 字符串 结束 在 `%`，它 表示 该 时间 point 是 指定 percentage 的 视频 时长; 对于 示例，`10%` 表示 该 时间 point 是 10% 的 视频 时长。",
															},
															"watermark_set": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "列表 up 到 10 镜像 或 text watermarks.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"definition": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "ID watermarking template。",
																		},
																		"raw_parameter": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Custom 水印 参数，其中 是 有效 如果 `Definition` 是 0.此 参数 是 使用 在 highly customized scenarios. We recommend 您 使用 `Definition` 到 指定watermark 参数 preferably.Custom 水印 参数 是 不 可用 对于 screenshot。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Watermark 类型 有效 值:镜像: 镜像 水印。",
																					},
																					"coordinate_origin": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Origin position，其中 currently 可以 仅 是:TopLeft: 源站 的 coordinates 是 在 top-left corner 的 视频，和 源站 的 水印 是 在 top-left corner 的 镜像 或 text.默认值：TopLeft。",
																					},
																					"x_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "horizontal position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `XPos` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `XPos` 是 10% 的 视频 宽度;如果 字符串 结束 在 像素， `XPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `XPos` 是 100 像素.默认值：0 像素。",
																					},
																					"y_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "vertical position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `YPos` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `YPos` 是 10% 的 视频 高度;如果 字符串 结束 在 像素， `YPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `YPos` 是 100 像素.默认值：0 像素。",
																					},
																					"image_template": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Image 水印 template. 此 字段 为必填项 当 `类型` 是 `镜像` 和 是 无效 当 `类型` 是 `text`。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"image_content": {
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Input 内容 的 水印 镜像. JPEG 和 PNG images 是 支持。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"type": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "input 类型 有效 值:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
																											},
																											"cos_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "信息 的 COS 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `COS`。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "COS 存储桶 的 对象 到 process，such 作为 `TopRankVideo-125xxx88`。",
																														},
																														"region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "地域 的 COS 存储桶，such 作为 `ap-chongqing`。",
																														},
																														"object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "路径 的 对象 到 process，such 作为 `/movie/201907/WildAnimal.mov`。",
																														},
																													},
																												},
																											},
																											"url_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "URL 的 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"url": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "URL 的 视频。",
																														},
																													},
																												},
																											},
																											"s3_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "信息 的 AWS S3 对象 processed. 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"s3_bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "AWS S3 存储桶",
																														},
																														"s3_region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "地域 的 AWS S3 存储桶",
																														},
																														"s3_object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "路径 的 AWS S3 对象。",
																														},
																														"s3_secret_id": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "键 ID 必填 到 访问 AWS S3 对象。",
																														},
																														"s3_secret_key": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "键 必填 到 访问 AWS S3 对象。",
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																								"width": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Watermark 宽度. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `宽度` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `宽度` 是 10% 的 视频 宽度;如果 字符串 结束 在 像素， `宽度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `宽度` 是 100 像素.默认值：10%。",
																								},
																								"height": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Watermark 高度. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `高度` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `高度` 是 10% 的 视频 高度;如果 字符串 结束 在 像素， `高度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `高度` 是 100 像素.默认值：0 像素，其中 表示 该 `高度` 将 是 proportionally scaled according 到 aspect ratio 的 original 水印 镜像。",
																								},
																								"repeat_type": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Repeat 类型 animated 水印. 有效 值:`once`: 无 longer appears after 水印 playback 结束.`repeat_last_frame`: stays 在 last frame after 水印 playback 结束.`repeat` (默认值): repeats playback until 视频 结束。",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"text_content": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Text 内容 的 up 到 100 字符. 此 字段 为必填项 仅 当 水印 类型 是 text.Text 水印 是 不 可用 对于 screenshot。",
																		},
																		"svg_content": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "SVG 内容 的 up 到 2,000,000 字符. 此 字段 为必填项 仅 当 水印 类型 是 `SVG`.SVG 水印 是 不 可用 对于 screenshot。",
																		},
																		"start_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "开始时间 偏移量 的 水印 （秒）。 如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame.如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame;如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 appear 在 second n after first 视频 frame;如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 appear 在 second n before last 视频 frame。",
																		},
																		"end_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "结束时间 偏移量 的 水印 （秒）。如果此参数为空 或 0 是 entered， 水印 将 exist till last 视频 frame;如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 exist till second n;如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 exist till second n before last 视频 frame。",
																		},
																	},
																},
															},
															"output_storage": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Target 存储桶 的 generated 时间 point screenshot 文件. 如果此参数为空， `OutputStorage` 值 的 upper 文件夹 将 是 inherited.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "存储 类型 对于 media processing 输出文件 有效 值:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
																		},
																		"cos_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "location 到 save output 对象 在 COS. 此 参数 是 有效 和 必填 当 `类型` 是 COS.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "存储桶 到 其中 输出文件 的 media processing 是 saved，such 作为 `TopRankVideo-125xxx88`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "地域 的 输出存储桶，such 作为 `ap-chongqing`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																					},
																				},
																			},
																		},
																		"s3_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "AWS S3 存储桶 到 save 输出文件 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "AWS S3 存储桶",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "地域 的 AWS S3 存储桶",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "键 ID 必填 到 upload files 到 AWS S3 对象。",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "键 必填 到 upload files 到 AWS S3 对象。",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"output_object_path": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "输出路径 到 generated 时间 point screenshot，其中 可以 是 relative 路径 或 absolute 路径 如果此参数为空， following relative 路径 将 是 使用 通过 默认值：`{inputName}_snapshotByTimeOffset_{definition}_{数量}.{格式}`。",
															},
															"object_number_format": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Rule 的 `{数量}` variable 在 时间 point screenshot 输出路径注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"initial_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Start 值 的 `{数量}` variable. 默认值：0。",
																		},
																		"increment": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Increment 的 `{数量}` variable. 默认值：1。",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "最小长度the `{数量}` variable. A placeholder 将 是 使用 如果 variable 长度 是 below 最小 requirement. 默认值：1。",
																		},
																		"place_holder": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Placeholder 使用 当 `{数量}` variable 长度 是 below 最小 requirement. 默认值：0。",
																		},
																	},
																},
															},
														},
													},
												},
												"sample_snapshot_task": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "A sampled screencapturing 任务。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Sampled screencapturing 模板 ID",
															},
															"watermark_set": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "列表 up 到 10 镜像 或 text watermarks.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"definition": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "ID watermarking template。",
																		},
																		"raw_parameter": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Custom 水印 参数，其中 是 有效 如果 `Definition` 是 0.此 参数 是 使用 在 highly customized scenarios. We recommend 您 使用 `Definition` 到 指定watermark 参数 preferably.Custom 水印 参数 是 不 可用 对于 screenshot。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Watermark 类型 有效 值:镜像: 镜像 水印。",
																					},
																					"coordinate_origin": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Origin position，其中 currently 可以 仅 是:TopLeft: 源站 的 coordinates 是 在 top-left corner 的 视频，和 源站 的 水印 是 在 top-left corner 的 镜像 或 text.默认值：TopLeft。",
																					},
																					"x_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "horizontal position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `XPos` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `XPos` 是 10% 的 视频 宽度;如果 字符串 结束 在 像素， `XPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `XPos` 是 100 像素.默认值：0 像素。",
																					},
																					"y_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "vertical position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `YPos` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `YPos` 是 10% 的 视频 高度;如果 字符串 结束 在 像素， `YPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `YPos` 是 100 像素.默认值：0 像素。",
																					},
																					"image_template": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Image 水印 template. 此 字段 为必填项 当 `类型` 是 `镜像` 和 是 无效 当 `类型` 是 `text`。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"image_content": {
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Input 内容 的 水印 镜像. JPEG 和 PNG images 是 支持。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"type": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "input 类型 有效 值:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
																											},
																											"cos_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "信息 的 COS 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `COS`。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "COS 存储桶 的 对象 到 process，such 作为 `TopRankVideo-125xxx88`。",
																														},
																														"region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "地域 的 COS 存储桶，such 作为 `ap-chongqing`。",
																														},
																														"object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "路径 的 对象 到 process，such 作为 `/movie/201907/WildAnimal.mov`。",
																														},
																													},
																												},
																											},
																											"url_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "URL 的 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"url": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "URL 的 视频。",
																														},
																													},
																												},
																											},
																											"s3_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "信息 的 AWS S3 对象 processed. 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"s3_bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "AWS S3 存储桶",
																														},
																														"s3_region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "地域 的 AWS S3 存储桶",
																														},
																														"s3_object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "路径 的 AWS S3 对象。",
																														},
																														"s3_secret_id": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "键 ID 必填 到 访问 AWS S3 对象。",
																														},
																														"s3_secret_key": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "键 必填 到 访问 AWS S3 对象。",
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																								"width": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Watermark 宽度. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `宽度` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `宽度` 是 10% 的 视频 宽度;如果 字符串 结束 在 像素， `宽度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `宽度` 是 100 像素.默认值：10%。",
																								},
																								"height": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Watermark 高度. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `高度` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `高度` 是 10% 的 视频 高度;如果 字符串 结束 在 像素， `高度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `高度` 是 100 像素.默认值：0 像素，其中 表示 该 `高度` 将 是 proportionally scaled according 到 aspect ratio 的 original 水印 镜像。",
																								},
																								"repeat_type": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Repeat 类型 animated 水印. 有效 值:`once`: 无 longer appears after 水印 playback 结束.`repeat_last_frame`: stays 在 last frame after 水印 playback 结束.`repeat` (默认值): repeats playback until 视频 结束。",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"text_content": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Text 内容 的 up 到 100 字符. 此 字段 为必填项 仅 当 水印 类型 是 text.Text 水印 是 不 可用 对于 screenshot。",
																		},
																		"svg_content": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "SVG 内容 的 up 到 2,000,000 字符. 此 字段 为必填项 仅 当 水印 类型 是 `SVG`.SVG 水印 是 不 可用 对于 screenshot。",
																		},
																		"start_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "开始时间 偏移量 的 水印 （秒）。 如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame.如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame;如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 appear 在 second n after first 视频 frame;如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 appear 在 second n before last 视频 frame。",
																		},
																		"end_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "结束时间 偏移量 的 水印 （秒）。如果此参数为空 或 0 是 entered， 水印 将 exist till last 视频 frame;如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 exist till second n;如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 exist till second n before last 视频 frame。",
																		},
																	},
																},
															},
															"output_storage": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Target 存储桶 的 sampled screenshot. 如果此参数为空， `OutputStorage` 值 的 upper 文件夹 将 是 inherited.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "存储 类型 对于 media processing 输出文件 有效 值:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
																		},
																		"cos_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "location 到 save output 对象 在 COS. 此 参数 是 有效 和 必填 当 `类型` 是 COS.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "存储桶 到 其中 输出文件 的 media processing 是 saved，such 作为 `TopRankVideo-125xxx88`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "地域 的 输出存储桶，such 作为 `ap-chongqing`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																					},
																				},
																			},
																		},
																		"s3_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "AWS S3 存储桶 到 save 输出文件 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "AWS S3 存储桶",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "地域 的 AWS S3 存储桶",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "键 ID 必填 到 upload files 到 AWS S3 对象。",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "键 必填 到 upload files 到 AWS S3 对象。",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"output_object_path": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "输出路径 到 generated sampled screenshot，其中 可以 是 relative 路径 或 absolute 路径 如果此参数为空， following relative 路径 将 是 使用 通过 默认值：`{inputName}_sampleSnapshot_{definition}_{数量}.{格式}`。",
															},
															"object_number_format": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Rule 的 `{数量}` variable 在 sampled screenshot 输出路径注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"initial_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Start 值 的 `{数量}` variable. 默认值：0。",
																		},
																		"increment": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Increment 的 `{数量}` variable. 默认值：1。",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "最小长度the `{数量}` variable. A placeholder 将 是 使用 如果 variable 长度 是 below 最小 requirement. 默认值：1。",
																		},
																		"place_holder": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Placeholder 使用 当 `{数量}` variable 长度 是 below 最小 requirement. 默认值：0。",
																		},
																	},
																},
															},
														},
													},
												},
												"image_sprite_task": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "An 镜像 sprite generation 任务。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "ID 镜像 sprite generating template。",
															},
															"output_storage": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Target 存储桶 的 generated 镜像 sprite. 如果此参数为空， `OutputStorage` 值 的 upper 文件夹 将 是 inherited.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "存储 类型 对于 media processing 输出文件 有效 值:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
																		},
																		"cos_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "location 到 save output 对象 在 COS. 此 参数 是 有效 和 必填 当 `类型` 是 COS.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "存储桶 到 其中 输出文件 的 media processing 是 saved，such 作为 `TopRankVideo-125xxx88`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "地域 的 输出存储桶，such 作为 `ap-chongqing`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																					},
																				},
																			},
																		},
																		"s3_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "AWS S3 存储桶 到 save 输出文件 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "AWS S3 存储桶",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "地域 的 AWS S3 存储桶",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "键 ID 必填 到 upload files 到 AWS S3 对象。",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "键 必填 到 upload files 到 AWS S3 对象。",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"output_object_path": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "输出路径 到 generated 镜像 sprite 文件，其中 可以 是 relative 路径 或 absolute 路径 如果此参数为空， following relative 路径 将 是 使用 通过 默认值：`{inputName}_imageSprite_{definition}_{数量}.{格式}`。",
															},
															"web_vtt_object_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "输出路径 到 WebVTT 文件 after 镜像 sprite 是 generated，其中 可以 仅 是 relative 路径 如果此参数为空， following relative 路径 将 是 使用 通过 默认值：`{inputName}_imageSprite_{definition}.{格式}`。",
															},
															"object_number_format": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Rule 的 `{数量}` variable 在 镜像 sprite 输出路径注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"initial_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Start 值 的 `{数量}` variable. 默认值：0。",
																		},
																		"increment": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Increment 的 `{数量}` variable. 默认值：1。",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "最小长度the `{数量}` variable. A placeholder 将 是 使用 如果 variable 长度 是 below 最小 requirement. 默认值：1。",
																		},
																		"place_holder": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Placeholder 使用 当 `{数量}` variable 长度 是 below 最小 requirement. 默认值：0。",
																		},
																	},
																},
															},
														},
													},
												},
												"adaptive_dynamic_streaming_task": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "An adaptive bitrate streaming 任务。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Adaptive bitrate streaming 模板 ID",
															},
															"watermark_set": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "列表 up 到 10 镜像 或 text watermarks。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"definition": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "ID watermarking template。",
																		},
																		"raw_parameter": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Custom 水印 参数，其中 是 有效 如果 `Definition` 是 0.此 参数 是 使用 在 highly customized scenarios. We recommend 您 使用 `Definition` 到 指定watermark 参数 preferably.Custom 水印 参数 是 不 可用 对于 screenshot。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Watermark 类型 有效 值:镜像: 镜像 水印。",
																					},
																					"coordinate_origin": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Origin position，其中 currently 可以 仅 是:TopLeft: 源站 的 coordinates 是 在 top-left corner 的 视频，和 源站 的 水印 是 在 top-left corner 的 镜像 或 text.默认值：TopLeft。",
																					},
																					"x_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "horizontal position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `XPos` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `XPos` 是 10% 的 视频 宽度;如果 字符串 结束 在 像素， `XPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `XPos` 是 100 像素.默认值：0 像素。",
																					},
																					"y_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "vertical position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `YPos` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `YPos` 是 10% 的 视频 高度;如果 字符串 结束 在 像素， `YPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `YPos` 是 100 像素.默认值：0 像素。",
																					},
																					"image_template": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Image 水印 template. 此 字段 为必填项 当 `类型` 是 `镜像` 和 是 无效 当 `类型` 是 `text`。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"image_content": {
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Input 内容 的 水印 镜像. JPEG 和 PNG images 是 支持。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"type": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "input 类型 有效 值:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
																											},
																											"cos_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "信息 的 COS 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `COS`。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "COS 存储桶 的 对象 到 process，such 作为 `TopRankVideo-125xxx88`。",
																														},
																														"region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "地域 的 COS 存储桶，such 作为 `ap-chongqing`。",
																														},
																														"object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "路径 的 对象 到 process，such 作为 `/movie/201907/WildAnimal.mov`。",
																														},
																													},
																												},
																											},
																											"url_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "URL 的 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"url": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "URL 的 视频。",
																														},
																													},
																												},
																											},
																											"s3_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "信息 的 AWS S3 对象 processed. 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"s3_bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "AWS S3 存储桶",
																														},
																														"s3_region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "地域 的 AWS S3 存储桶",
																														},
																														"s3_object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "路径 的 AWS S3 对象。",
																														},
																														"s3_secret_id": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "键 ID 必填 到 访问 AWS S3 对象。",
																														},
																														"s3_secret_key": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "键 必填 到 访问 AWS S3 对象。",
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																								"width": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Watermark 宽度. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `宽度` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `宽度` 是 10% 的 视频 宽度;如果 字符串 结束 在 像素， `宽度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `宽度` 是 100 像素.默认值：10%。",
																								},
																								"height": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Watermark 高度. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `高度` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `高度` 是 10% 的 视频 高度;如果 字符串 结束 在 像素， `高度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `高度` 是 100 像素.默认值：0 像素，其中 表示 该 `高度` 将 是 proportionally scaled according 到 aspect ratio 的 original 水印 镜像。",
																								},
																								"repeat_type": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Repeat 类型 animated 水印. 有效 值:`once`: 无 longer appears after 水印 playback 结束.`repeat_last_frame`: stays 在 last frame after 水印 playback 结束.`repeat` (默认值): repeats playback until 视频 结束。",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"text_content": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Text 内容 的 up 到 100 字符. 此 字段 为必填项 仅 当 水印 类型 是 text.Text 水印 是 不 可用 对于 screenshot。",
																		},
																		"svg_content": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "SVG 内容 的 up 到 2,000,000 字符. 此 字段 为必填项 仅 当 水印 类型 是 `SVG`.SVG 水印 是 不 可用 对于 screenshot。",
																		},
																		"start_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "开始时间 偏移量 的 水印 （秒）。 如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame.如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame;如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 appear 在 second n after first 视频 frame;如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 appear 在 second n before last 视频 frame。",
																		},
																		"end_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "结束时间 偏移量 的 水印 （秒）。如果此参数为空 或 0 是 entered， 水印 将 exist till last 视频 frame;如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 exist till second n;如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 exist till second n before last 视频 frame。",
																		},
																	},
																},
															},
															"output_storage": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Target 存储桶 的 输出文件 after being transcoded 到 adaptive bitrate streaming. 如果此参数为空， `OutputStorage` 值 的 upper 文件夹 将 是 inherited.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "存储 类型 对于 media processing 输出文件 有效 值:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
																		},
																		"cos_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "location 到 save output 对象 在 COS. 此 参数 是 有效 和 必填 当 `类型` 是 COS.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "存储桶 到 其中 输出文件 的 media processing 是 saved，such 作为 `TopRankVideo-125xxx88`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "地域 的 输出存储桶，such 作为 `ap-chongqing`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																					},
																				},
																			},
																		},
																		"s3_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "AWS S3 存储桶 到 save 输出文件 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "AWS S3 存储桶",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "地域 的 AWS S3 存储桶",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "键 ID 必填 到 upload files 到 AWS S3 对象。",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "键 必填 到 upload files 到 AWS S3 对象。",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"output_object_path": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "relative 或 absolute 输出路径 的 manifest 文件 after being transcoded 到 adaptive bitrate streaming. 如果此参数为空， relative 路径 在 following 格式 将 是 使用 通过 默认值：`{inputName}_adaptiveDynamicStreaming_{definition}.{格式}`。",
															},
															"sub_stream_object_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "relative 输出路径 的 substream 文件 after being transcoded 到 adaptive bitrate streaming. 如果此参数为空， relative 路径 在 following 格式 将 是 使用 通过 默认值：`{inputName}_adaptiveDynamicStreaming_{definition}_{subStreamNumber}.{格式}`。",
															},
															"segment_object_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "relative 输出路径 的 segment 文件 after being transcoded 到 adaptive bitrate streaming (在 HLS 格式 仅). 如果此参数为空， relative 路径 在 following 格式 将 是 使用 通过 默认值：`{inputName}_adaptiveDynamicStreaming_{definition}_{subStreamNumber}_{segmentNumber}.{格式}`。",
															},
															"add_on_subtitles": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "subtitle 文件 到 add.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "模式 有效 值:`subtitle-流`: Add subtitle track.`close-caption-708`: Embed EA-708 subtitles 在 SEI frames.`close-caption-608`: Embed CEA-608 subtitles 在 SEI frames.注意：此字段可能返回 null，表示无法获取有效值。",
																		},
																		"subtitle": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "subtitle 文件.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "input 类型 有效 值:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "信息 的 COS 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `COS`。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "COS 存储桶 的 对象 到 process，such 作为 `TopRankVideo-125xxx88`。",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "地域 的 COS 存储桶，such 作为 `ap-chongqing`。",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "路径 的 对象 到 process，such 作为 `/movie/201907/WildAnimal.mov`。",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "URL 的 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "URL 的 视频。",
																								},
																							},
																						},
																					},
																					"s3_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "信息 的 AWS S3 对象 processed. 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"s3_bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "AWS S3 存储桶",
																								},
																								"s3_region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "地域 的 AWS S3 存储桶",
																								},
																								"s3_object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "路径 的 AWS S3 对象。",
																								},
																								"s3_secret_id": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "键 ID 必填 到 访问 AWS S3 对象。",
																								},
																								"s3_secret_key": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "键 必填 到 访问 AWS S3 对象。",
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
												"ai_content_review_task": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "A 内容 moderation 任务。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Video 内容 audit 模板 ID",
															},
														},
													},
												},
												"ai_analysis_task": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "A 内容 analysis 任务。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Video 内容 analysis 模板 ID",
															},
															"extended_parameter": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "An extended 参数，whose 值 是 stringfied JSON.注意: 此 参数 是 对于 customers 使用 special requirements. It needs 到 是 customized offline.注意：此字段可能返回 null，表示无法获取有效值。",
															},
														},
													},
												},
												"ai_recognition_task": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "A 内容 recognition 任务。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Intelligent 视频 recognition 模板 ID",
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
						"output_storage": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "存储桶 到 save 输出文件注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "存储 类型 对于 media processing 输出文件 有效 值:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
									},
									"cos_output_storage": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "location 到 save output 对象 在 COS. 此 参数 是 有效 和 必填 当 `类型` 是 COS.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bucket": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "存储桶 到 其中 输出文件 的 media processing 是 saved，such 作为 `TopRankVideo-125xxx88`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
												},
												"region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "地域 的 输出存储桶，such 作为 `ap-chongqing`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
												},
											},
										},
									},
									"s3_output_storage": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "AWS S3 存储桶 到 save 输出文件 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"s3_bucket": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "AWS S3 存储桶",
												},
												"s3_region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "地域 的 AWS S3 存储桶",
												},
												"s3_secret_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "键 ID 必填 到 upload files 到 AWS S3 对象。",
												},
												"s3_secret_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "键 必填 到 upload files 到 AWS S3 对象。",
												},
											},
										},
									},
								},
							},
						},
						"output_dir": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "directory 到 save 输出文件注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"task_notify_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "通知 配置.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cmq_model": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CMQ 或 TDMQ-CMQ model. 有效值：Queue，Topic。",
									},
									"cmq_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CMQ 或 TDMQ-CMQ 地域，such 作为 `sh` (Shanghai) 或 `bj` (Beijing)。",
									},
									"topic_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CMQ 或 TDMQ-CMQ 主题 到 receive notifications. 此 参数 是 有效 当 `CmqModel` 是 `Topic`。",
									},
									"queue_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CMQ 或 TDMQ-CMQ queue 到 receive notifications. 此 参数 是 有效 当 `CmqModel` 是 `Queue`。",
									},
									"notify_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Workflow 通知 方法. 有效值：Finish，Change. 如果此参数为空，`Finish` 将 是 使用。",
									},
									"notify_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "通知 类型 有效 值:`CMQ`: 此 值 是 无 longer 使用. Please 使用 `TDMQ-CMQ` instead.`TDMQ-CMQ`: 消息 queue`URL`: 如果 `NotifyType` 是 集合 到 `URL`，HTTP callbacks 是 sent 到 URL 指定 通过 `NotifyUrl`. HTTP 和 JSON 是 用于the callbacks. packet 包含response 参数 的 `ParseNotification` API.`SCF`: 此 通知 类型 是 不 recommended. You need 到 configure 它 在 SCF console.`AWS-SQS`: AWS queue. 此 类型 是 仅 支持 对于 AWS tasks，和 queue 必须 是 在 same 地域 作为 AWS 存储桶Note: 如果 您 do 不 pass 此 参数 或 pass 在 空 字符串，`CMQ` 将 是 使用. To 使用 different 通知 类型，指定this 参数 accordingly。",
									},
									"notify_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "HTTP callback URL，必填 如果 `NotifyType` 是 集合 到 `URL`。",
									},
									"aws_sqs": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "AWS SQS queue. 此 参数 为必填项 如果 `NotifyType` 是 `AWS-SQS`.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"sqs_region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "地域 的 SQS queue。",
												},
												"sqs_queue_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 SQS queue。",
												},
												"s3_secret_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "键 ID 必填 到 read 从/write 到 SQS queue。",
												},
												"s3_secret_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "键 必填 到 read 从/write 到 SQS queue。",
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
							Description: "创建时间 在 [ISO date 格式](https://intl.云.tencent.com/document/product/862/37710?from_cn_redirect=1#52).注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "last 更新 时间 在 [ISO date 格式](https://intl.云.tencent.com/document/product/862/37710?from_cn_redirect=1#52).注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "资源 ID. 如果 there 是 无 associated 资源 ID，fill 它 使用 账号's main 资源 ID。",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudMpsSchedulesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mps_schedules.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("schedule_ids"); ok {
		scheduleIdsSet := v.(*schema.Set).List()
		scheduleIdList := []interface{}{}
		for i := range scheduleIdsSet {
			scheduleIds := scheduleIdsSet[i].(int)
			scheduleIdList = append(scheduleIdList, scheduleIds)
		}
		paramMap["ScheduleIds"] = scheduleIdList
	}

	if v, ok := d.GetOk("trigger_type"); ok {
		paramMap["TriggerType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		paramMap["Status"] = helper.String(v.(string))
	}

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var scheduleInfoSet []*mps.SchedulesInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMpsSchedulesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		scheduleInfoSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(scheduleInfoSet))
	tmpList := make([]map[string]interface{}, 0, len(scheduleInfoSet))

	if scheduleInfoSet != nil {
		for _, schedulesInfo := range scheduleInfoSet {
			schedulesInfoMap := map[string]interface{}{}

			if schedulesInfo.ScheduleId != nil {
				schedulesInfoMap["schedule_id"] = schedulesInfo.ScheduleId
			}

			if schedulesInfo.ScheduleName != nil {
				schedulesInfoMap["schedule_name"] = schedulesInfo.ScheduleName
			}

			if schedulesInfo.Status != nil {
				schedulesInfoMap["status"] = schedulesInfo.Status
			}

			if schedulesInfo.Trigger != nil {
				triggerMap := map[string]interface{}{}

				if schedulesInfo.Trigger.Type != nil {
					triggerMap["type"] = schedulesInfo.Trigger.Type
				}

				if schedulesInfo.Trigger.CosFileUploadTrigger != nil {
					cosFileUploadTriggerMap := map[string]interface{}{}

					if schedulesInfo.Trigger.CosFileUploadTrigger.Bucket != nil {
						cosFileUploadTriggerMap["bucket"] = schedulesInfo.Trigger.CosFileUploadTrigger.Bucket
					}

					if schedulesInfo.Trigger.CosFileUploadTrigger.Region != nil {
						cosFileUploadTriggerMap["region"] = schedulesInfo.Trigger.CosFileUploadTrigger.Region
					}

					if schedulesInfo.Trigger.CosFileUploadTrigger.Dir != nil {
						cosFileUploadTriggerMap["dir"] = schedulesInfo.Trigger.CosFileUploadTrigger.Dir
					}

					if schedulesInfo.Trigger.CosFileUploadTrigger.Formats != nil {
						cosFileUploadTriggerMap["formats"] = schedulesInfo.Trigger.CosFileUploadTrigger.Formats
					}

					triggerMap["cos_file_upload_trigger"] = []interface{}{cosFileUploadTriggerMap}
				}

				if schedulesInfo.Trigger.AwsS3FileUploadTrigger != nil {
					awsS3FileUploadTriggerMap := map[string]interface{}{}

					if schedulesInfo.Trigger.AwsS3FileUploadTrigger.S3Bucket != nil {
						awsS3FileUploadTriggerMap["s3_bucket"] = schedulesInfo.Trigger.AwsS3FileUploadTrigger.S3Bucket
					}

					if schedulesInfo.Trigger.AwsS3FileUploadTrigger.S3Region != nil {
						awsS3FileUploadTriggerMap["s3_region"] = schedulesInfo.Trigger.AwsS3FileUploadTrigger.S3Region
					}

					if schedulesInfo.Trigger.AwsS3FileUploadTrigger.Dir != nil {
						awsS3FileUploadTriggerMap["dir"] = schedulesInfo.Trigger.AwsS3FileUploadTrigger.Dir
					}

					if schedulesInfo.Trigger.AwsS3FileUploadTrigger.Formats != nil {
						awsS3FileUploadTriggerMap["formats"] = schedulesInfo.Trigger.AwsS3FileUploadTrigger.Formats
					}

					if schedulesInfo.Trigger.AwsS3FileUploadTrigger.S3SecretId != nil {
						awsS3FileUploadTriggerMap["s3_secret_id"] = schedulesInfo.Trigger.AwsS3FileUploadTrigger.S3SecretId
					}

					if schedulesInfo.Trigger.AwsS3FileUploadTrigger.S3SecretKey != nil {
						awsS3FileUploadTriggerMap["s3_secret_key"] = schedulesInfo.Trigger.AwsS3FileUploadTrigger.S3SecretKey
					}

					if schedulesInfo.Trigger.AwsS3FileUploadTrigger.AwsSQS != nil {
						awsSQSMap := map[string]interface{}{}

						if schedulesInfo.Trigger.AwsS3FileUploadTrigger.AwsSQS.SQSRegion != nil {
							awsSQSMap["sqs_region"] = schedulesInfo.Trigger.AwsS3FileUploadTrigger.AwsSQS.SQSRegion
						}

						if schedulesInfo.Trigger.AwsS3FileUploadTrigger.AwsSQS.SQSQueueName != nil {
							awsSQSMap["sqs_queue_name"] = schedulesInfo.Trigger.AwsS3FileUploadTrigger.AwsSQS.SQSQueueName
						}

						if schedulesInfo.Trigger.AwsS3FileUploadTrigger.AwsSQS.S3SecretId != nil {
							awsSQSMap["s3_secret_id"] = schedulesInfo.Trigger.AwsS3FileUploadTrigger.AwsSQS.S3SecretId
						}

						if schedulesInfo.Trigger.AwsS3FileUploadTrigger.AwsSQS.S3SecretKey != nil {
							awsSQSMap["s3_secret_key"] = schedulesInfo.Trigger.AwsS3FileUploadTrigger.AwsSQS.S3SecretKey
						}

						awsS3FileUploadTriggerMap["aws_sqs"] = []interface{}{awsSQSMap}
					}

					triggerMap["aws_s3_file_upload_trigger"] = []interface{}{awsS3FileUploadTriggerMap}
				}

				schedulesInfoMap["trigger"] = []interface{}{triggerMap}
			}

			if schedulesInfo.Activities != nil {
				activitiesList := []interface{}{}
				for _, activities := range schedulesInfo.Activities {
					activitiesMap := map[string]interface{}{}

					if activities.ActivityType != nil {
						activitiesMap["activity_type"] = activities.ActivityType
					}

					if activities.ReardriveIndex != nil {
						activitiesMap["reardrive_index"] = activities.ReardriveIndex
					}

					if activities.ActivityPara != nil {
						activityParaMap := map[string]interface{}{}

						if activities.ActivityPara.TranscodeTask != nil {
							transcodeTaskMap := map[string]interface{}{}

							if activities.ActivityPara.TranscodeTask.Definition != nil {
								transcodeTaskMap["definition"] = activities.ActivityPara.TranscodeTask.Definition
							}

							if activities.ActivityPara.TranscodeTask.RawParameter != nil {
								rawParameterMap := map[string]interface{}{}

								if activities.ActivityPara.TranscodeTask.RawParameter.Container != nil {
									rawParameterMap["container"] = activities.ActivityPara.TranscodeTask.RawParameter.Container
								}

								if activities.ActivityPara.TranscodeTask.RawParameter.RemoveVideo != nil {
									rawParameterMap["remove_video"] = activities.ActivityPara.TranscodeTask.RawParameter.RemoveVideo
								}

								if activities.ActivityPara.TranscodeTask.RawParameter.RemoveAudio != nil {
									rawParameterMap["remove_audio"] = activities.ActivityPara.TranscodeTask.RawParameter.RemoveAudio
								}

								if activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate != nil {
									videoTemplateMap := map[string]interface{}{}

									if activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Codec != nil {
										videoTemplateMap["codec"] = activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Codec
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Fps != nil {
										videoTemplateMap["fps"] = activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Fps
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Bitrate != nil {
										videoTemplateMap["bitrate"] = activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Bitrate
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.ResolutionAdaptive != nil {
										videoTemplateMap["resolution_adaptive"] = activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.ResolutionAdaptive
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Width != nil {
										videoTemplateMap["width"] = activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Width
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Height != nil {
										videoTemplateMap["height"] = activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Height
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Gop != nil {
										videoTemplateMap["gop"] = activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Gop
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.FillType != nil {
										videoTemplateMap["fill_type"] = activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.FillType
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Vcrf != nil {
										videoTemplateMap["vcrf"] = activities.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Vcrf
									}

									rawParameterMap["video_template"] = []interface{}{videoTemplateMap}
								}

								if activities.ActivityPara.TranscodeTask.RawParameter.AudioTemplate != nil {
									audioTemplateMap := map[string]interface{}{}

									if activities.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.Codec != nil {
										audioTemplateMap["codec"] = activities.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.Codec
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.Bitrate != nil {
										audioTemplateMap["bitrate"] = activities.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.Bitrate
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.SampleRate != nil {
										audioTemplateMap["sample_rate"] = activities.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.SampleRate
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.AudioChannel != nil {
										audioTemplateMap["audio_channel"] = activities.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.AudioChannel
									}

									rawParameterMap["audio_template"] = []interface{}{audioTemplateMap}
								}

								if activities.ActivityPara.TranscodeTask.RawParameter.TEHDConfig != nil {
									tEHDConfigMap := map[string]interface{}{}

									if activities.ActivityPara.TranscodeTask.RawParameter.TEHDConfig.Type != nil {
										tEHDConfigMap["type"] = activities.ActivityPara.TranscodeTask.RawParameter.TEHDConfig.Type
									}

									if activities.ActivityPara.TranscodeTask.RawParameter.TEHDConfig.MaxVideoBitrate != nil {
										tEHDConfigMap["max_video_bitrate"] = activities.ActivityPara.TranscodeTask.RawParameter.TEHDConfig.MaxVideoBitrate
									}

									rawParameterMap["tehd_config"] = []interface{}{tEHDConfigMap}
								}

								transcodeTaskMap["raw_parameter"] = []interface{}{rawParameterMap}
							}

							if activities.ActivityPara.TranscodeTask.OverrideParameter != nil {
								overrideParameterMap := map[string]interface{}{}

								if activities.ActivityPara.TranscodeTask.OverrideParameter.Container != nil {
									overrideParameterMap["container"] = activities.ActivityPara.TranscodeTask.OverrideParameter.Container
								}

								if activities.ActivityPara.TranscodeTask.OverrideParameter.RemoveVideo != nil {
									overrideParameterMap["remove_video"] = activities.ActivityPara.TranscodeTask.OverrideParameter.RemoveVideo
								}

								if activities.ActivityPara.TranscodeTask.OverrideParameter.RemoveAudio != nil {
									overrideParameterMap["remove_audio"] = activities.ActivityPara.TranscodeTask.OverrideParameter.RemoveAudio
								}

								if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate != nil {
									videoTemplateMap := map[string]interface{}{}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Codec != nil {
										videoTemplateMap["codec"] = activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Codec
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Fps != nil {
										videoTemplateMap["fps"] = activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Fps
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Bitrate != nil {
										videoTemplateMap["bitrate"] = activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Bitrate
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.ResolutionAdaptive != nil {
										videoTemplateMap["resolution_adaptive"] = activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.ResolutionAdaptive
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Width != nil {
										videoTemplateMap["width"] = activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Width
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Height != nil {
										videoTemplateMap["height"] = activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Height
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Gop != nil {
										videoTemplateMap["gop"] = activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Gop
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.FillType != nil {
										videoTemplateMap["fill_type"] = activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.FillType
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Vcrf != nil {
										videoTemplateMap["vcrf"] = activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Vcrf
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.ContentAdaptStream != nil {
										videoTemplateMap["content_adapt_stream"] = activities.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.ContentAdaptStream
									}

									overrideParameterMap["video_template"] = []interface{}{videoTemplateMap}
								}

								if activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate != nil {
									audioTemplateMap := map[string]interface{}{}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.Codec != nil {
										audioTemplateMap["codec"] = activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.Codec
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.Bitrate != nil {
										audioTemplateMap["bitrate"] = activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.Bitrate
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.SampleRate != nil {
										audioTemplateMap["sample_rate"] = activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.SampleRate
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.AudioChannel != nil {
										audioTemplateMap["audio_channel"] = activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.AudioChannel
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.StreamSelects != nil {
										audioTemplateMap["stream_selects"] = activities.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.StreamSelects
									}

									overrideParameterMap["audio_template"] = []interface{}{audioTemplateMap}
								}

								if activities.ActivityPara.TranscodeTask.OverrideParameter.TEHDConfig != nil {
									tEHDConfigMap := map[string]interface{}{}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.TEHDConfig.Type != nil {
										tEHDConfigMap["type"] = activities.ActivityPara.TranscodeTask.OverrideParameter.TEHDConfig.Type
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.TEHDConfig.MaxVideoBitrate != nil {
										tEHDConfigMap["max_video_bitrate"] = activities.ActivityPara.TranscodeTask.OverrideParameter.TEHDConfig.MaxVideoBitrate
									}

									overrideParameterMap["tehd_config"] = []interface{}{tEHDConfigMap}
								}

								if activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate != nil {
									subtitleTemplateMap := map[string]interface{}{}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.Path != nil {
										subtitleTemplateMap["path"] = activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.Path
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.StreamIndex != nil {
										subtitleTemplateMap["stream_index"] = activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.StreamIndex
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontType != nil {
										subtitleTemplateMap["font_type"] = activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontType
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontSize != nil {
										subtitleTemplateMap["font_size"] = activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontSize
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontColor != nil {
										subtitleTemplateMap["font_color"] = activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontColor
									}

									if activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontAlpha != nil {
										subtitleTemplateMap["font_alpha"] = activities.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontAlpha
									}

									overrideParameterMap["subtitle_template"] = []interface{}{subtitleTemplateMap}
								}

								if activities.ActivityPara.TranscodeTask.OverrideParameter.AddonAudioStream != nil {
									addonAudioStreamList := []interface{}{}
									for _, addonAudioStream := range activities.ActivityPara.TranscodeTask.OverrideParameter.AddonAudioStream {
										addonAudioStreamMap := map[string]interface{}{}

										if addonAudioStream.Type != nil {
											addonAudioStreamMap["type"] = addonAudioStream.Type
										}

										if addonAudioStream.CosInputInfo != nil {
											cosInputInfoMap := map[string]interface{}{}

											if addonAudioStream.CosInputInfo.Bucket != nil {
												cosInputInfoMap["bucket"] = addonAudioStream.CosInputInfo.Bucket
											}

											if addonAudioStream.CosInputInfo.Region != nil {
												cosInputInfoMap["region"] = addonAudioStream.CosInputInfo.Region
											}

											if addonAudioStream.CosInputInfo.Object != nil {
												cosInputInfoMap["object"] = addonAudioStream.CosInputInfo.Object
											}

											addonAudioStreamMap["cos_input_info"] = []interface{}{cosInputInfoMap}
										}

										if addonAudioStream.UrlInputInfo != nil {
											urlInputInfoMap := map[string]interface{}{}

											if addonAudioStream.UrlInputInfo.Url != nil {
												urlInputInfoMap["url"] = addonAudioStream.UrlInputInfo.Url
											}

											addonAudioStreamMap["url_input_info"] = []interface{}{urlInputInfoMap}
										}

										if addonAudioStream.S3InputInfo != nil {
											s3InputInfoMap := map[string]interface{}{}

											if addonAudioStream.S3InputInfo.S3Bucket != nil {
												s3InputInfoMap["s3_bucket"] = addonAudioStream.S3InputInfo.S3Bucket
											}

											if addonAudioStream.S3InputInfo.S3Region != nil {
												s3InputInfoMap["s3_region"] = addonAudioStream.S3InputInfo.S3Region
											}

											if addonAudioStream.S3InputInfo.S3Object != nil {
												s3InputInfoMap["s3_object"] = addonAudioStream.S3InputInfo.S3Object
											}

											if addonAudioStream.S3InputInfo.S3SecretId != nil {
												s3InputInfoMap["s3_secret_id"] = addonAudioStream.S3InputInfo.S3SecretId
											}

											if addonAudioStream.S3InputInfo.S3SecretKey != nil {
												s3InputInfoMap["s3_secret_key"] = addonAudioStream.S3InputInfo.S3SecretKey
											}

											addonAudioStreamMap["s3_input_info"] = []interface{}{s3InputInfoMap}
										}

										addonAudioStreamList = append(addonAudioStreamList, addonAudioStreamMap)
									}

									overrideParameterMap["addon_audio_stream"] = addonAudioStreamList
								}

								if activities.ActivityPara.TranscodeTask.OverrideParameter.StdExtInfo != nil {
									overrideParameterMap["std_ext_info"] = activities.ActivityPara.TranscodeTask.OverrideParameter.StdExtInfo
								}

								if activities.ActivityPara.TranscodeTask.OverrideParameter.AddOnSubtitles != nil {
									addOnSubtitlesList := []interface{}{}
									for _, addOnSubtitles := range activities.ActivityPara.TranscodeTask.OverrideParameter.AddOnSubtitles {
										addOnSubtitlesMap := map[string]interface{}{}

										if addOnSubtitles.Type != nil {
											addOnSubtitlesMap["type"] = addOnSubtitles.Type
										}

										if addOnSubtitles.Subtitle != nil {
											subtitleMap := map[string]interface{}{}

											if addOnSubtitles.Subtitle.Type != nil {
												subtitleMap["type"] = addOnSubtitles.Subtitle.Type
											}

											if addOnSubtitles.Subtitle.CosInputInfo != nil {
												cosInputInfoMap := map[string]interface{}{}

												if addOnSubtitles.Subtitle.CosInputInfo.Bucket != nil {
													cosInputInfoMap["bucket"] = addOnSubtitles.Subtitle.CosInputInfo.Bucket
												}

												if addOnSubtitles.Subtitle.CosInputInfo.Region != nil {
													cosInputInfoMap["region"] = addOnSubtitles.Subtitle.CosInputInfo.Region
												}

												if addOnSubtitles.Subtitle.CosInputInfo.Object != nil {
													cosInputInfoMap["object"] = addOnSubtitles.Subtitle.CosInputInfo.Object
												}

												subtitleMap["cos_input_info"] = []interface{}{cosInputInfoMap}
											}

											if addOnSubtitles.Subtitle.UrlInputInfo != nil {
												urlInputInfoMap := map[string]interface{}{}

												if addOnSubtitles.Subtitle.UrlInputInfo.Url != nil {
													urlInputInfoMap["url"] = addOnSubtitles.Subtitle.UrlInputInfo.Url
												}

												subtitleMap["url_input_info"] = []interface{}{urlInputInfoMap}
											}

											if addOnSubtitles.Subtitle.S3InputInfo != nil {
												s3InputInfoMap := map[string]interface{}{}

												if addOnSubtitles.Subtitle.S3InputInfo.S3Bucket != nil {
													s3InputInfoMap["s3_bucket"] = addOnSubtitles.Subtitle.S3InputInfo.S3Bucket
												}

												if addOnSubtitles.Subtitle.S3InputInfo.S3Region != nil {
													s3InputInfoMap["s3_region"] = addOnSubtitles.Subtitle.S3InputInfo.S3Region
												}

												if addOnSubtitles.Subtitle.S3InputInfo.S3Object != nil {
													s3InputInfoMap["s3_object"] = addOnSubtitles.Subtitle.S3InputInfo.S3Object
												}

												if addOnSubtitles.Subtitle.S3InputInfo.S3SecretId != nil {
													s3InputInfoMap["s3_secret_id"] = addOnSubtitles.Subtitle.S3InputInfo.S3SecretId
												}

												if addOnSubtitles.Subtitle.S3InputInfo.S3SecretKey != nil {
													s3InputInfoMap["s3_secret_key"] = addOnSubtitles.Subtitle.S3InputInfo.S3SecretKey
												}

												subtitleMap["s3_input_info"] = []interface{}{s3InputInfoMap}
											}

											addOnSubtitlesMap["subtitle"] = []interface{}{subtitleMap}
										}

										addOnSubtitlesList = append(addOnSubtitlesList, addOnSubtitlesMap)
									}

									overrideParameterMap["add_on_subtitles"] = addOnSubtitlesList
								}

								transcodeTaskMap["override_parameter"] = []interface{}{overrideParameterMap}
							}

							if activities.ActivityPara.TranscodeTask.WatermarkSet != nil {
								watermarkSetList := []interface{}{}
								for _, watermarkSet := range activities.ActivityPara.TranscodeTask.WatermarkSet {
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

												if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo != nil {
													s3InputInfoMap := map[string]interface{}{}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Bucket != nil {
														s3InputInfoMap["s3_bucket"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Bucket
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Region != nil {
														s3InputInfoMap["s3_region"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Region
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Object != nil {
														s3InputInfoMap["s3_object"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Object
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretId != nil {
														s3InputInfoMap["s3_secret_id"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretId
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretKey != nil {
														s3InputInfoMap["s3_secret_key"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretKey
													}

													imageContentMap["s3_input_info"] = []interface{}{s3InputInfoMap}
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

								transcodeTaskMap["watermark_set"] = watermarkSetList
							}

							if activities.ActivityPara.TranscodeTask.MosaicSet != nil {
								mosaicSetList := []interface{}{}
								for _, mosaicSet := range activities.ActivityPara.TranscodeTask.MosaicSet {
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

								transcodeTaskMap["mosaic_set"] = mosaicSetList
							}

							if activities.ActivityPara.TranscodeTask.StartTimeOffset != nil {
								transcodeTaskMap["start_time_offset"] = activities.ActivityPara.TranscodeTask.StartTimeOffset
							}

							if activities.ActivityPara.TranscodeTask.EndTimeOffset != nil {
								transcodeTaskMap["end_time_offset"] = activities.ActivityPara.TranscodeTask.EndTimeOffset
							}

							if activities.ActivityPara.TranscodeTask.OutputStorage != nil {
								outputStorageMap := map[string]interface{}{}

								if activities.ActivityPara.TranscodeTask.OutputStorage.Type != nil {
									outputStorageMap["type"] = activities.ActivityPara.TranscodeTask.OutputStorage.Type
								}

								if activities.ActivityPara.TranscodeTask.OutputStorage.CosOutputStorage != nil {
									cosOutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.TranscodeTask.OutputStorage.CosOutputStorage.Bucket != nil {
										cosOutputStorageMap["bucket"] = activities.ActivityPara.TranscodeTask.OutputStorage.CosOutputStorage.Bucket
									}

									if activities.ActivityPara.TranscodeTask.OutputStorage.CosOutputStorage.Region != nil {
										cosOutputStorageMap["region"] = activities.ActivityPara.TranscodeTask.OutputStorage.CosOutputStorage.Region
									}

									outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
								}

								if activities.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage != nil {
									s3OutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3Bucket != nil {
										s3OutputStorageMap["s3_bucket"] = activities.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3Bucket
									}

									if activities.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3Region != nil {
										s3OutputStorageMap["s3_region"] = activities.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3Region
									}

									if activities.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3SecretId != nil {
										s3OutputStorageMap["s3_secret_id"] = activities.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3SecretId
									}

									if activities.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3SecretKey != nil {
										s3OutputStorageMap["s3_secret_key"] = activities.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3SecretKey
									}

									outputStorageMap["s3_output_storage"] = []interface{}{s3OutputStorageMap}
								}

								transcodeTaskMap["output_storage"] = []interface{}{outputStorageMap}
							}

							if activities.ActivityPara.TranscodeTask.OutputObjectPath != nil {
								transcodeTaskMap["output_object_path"] = activities.ActivityPara.TranscodeTask.OutputObjectPath
							}

							if activities.ActivityPara.TranscodeTask.SegmentObjectName != nil {
								transcodeTaskMap["segment_object_name"] = activities.ActivityPara.TranscodeTask.SegmentObjectName
							}

							if activities.ActivityPara.TranscodeTask.ObjectNumberFormat != nil {
								objectNumberFormatMap := map[string]interface{}{}

								if activities.ActivityPara.TranscodeTask.ObjectNumberFormat.InitialValue != nil {
									objectNumberFormatMap["initial_value"] = activities.ActivityPara.TranscodeTask.ObjectNumberFormat.InitialValue
								}

								if activities.ActivityPara.TranscodeTask.ObjectNumberFormat.Increment != nil {
									objectNumberFormatMap["increment"] = activities.ActivityPara.TranscodeTask.ObjectNumberFormat.Increment
								}

								if activities.ActivityPara.TranscodeTask.ObjectNumberFormat.MinLength != nil {
									objectNumberFormatMap["min_length"] = activities.ActivityPara.TranscodeTask.ObjectNumberFormat.MinLength
								}

								if activities.ActivityPara.TranscodeTask.ObjectNumberFormat.PlaceHolder != nil {
									objectNumberFormatMap["place_holder"] = activities.ActivityPara.TranscodeTask.ObjectNumberFormat.PlaceHolder
								}

								transcodeTaskMap["object_number_format"] = []interface{}{objectNumberFormatMap}
							}

							if activities.ActivityPara.TranscodeTask.HeadTailParameter != nil {
								headTailParameterMap := map[string]interface{}{}

								if activities.ActivityPara.TranscodeTask.HeadTailParameter.HeadSet != nil {
									headSetList := []interface{}{}
									for _, headSet := range activities.ActivityPara.TranscodeTask.HeadTailParameter.HeadSet {
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

										if headSet.S3InputInfo != nil {
											s3InputInfoMap := map[string]interface{}{}

											if headSet.S3InputInfo.S3Bucket != nil {
												s3InputInfoMap["s3_bucket"] = headSet.S3InputInfo.S3Bucket
											}

											if headSet.S3InputInfo.S3Region != nil {
												s3InputInfoMap["s3_region"] = headSet.S3InputInfo.S3Region
											}

											if headSet.S3InputInfo.S3Object != nil {
												s3InputInfoMap["s3_object"] = headSet.S3InputInfo.S3Object
											}

											if headSet.S3InputInfo.S3SecretId != nil {
												s3InputInfoMap["s3_secret_id"] = headSet.S3InputInfo.S3SecretId
											}

											if headSet.S3InputInfo.S3SecretKey != nil {
												s3InputInfoMap["s3_secret_key"] = headSet.S3InputInfo.S3SecretKey
											}

											headSetMap["s3_input_info"] = []interface{}{s3InputInfoMap}
										}

										headSetList = append(headSetList, headSetMap)
									}

									headTailParameterMap["head_set"] = headSetList
								}

								if activities.ActivityPara.TranscodeTask.HeadTailParameter.TailSet != nil {
									tailSetList := []interface{}{}
									for _, tailSet := range activities.ActivityPara.TranscodeTask.HeadTailParameter.TailSet {
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

										if tailSet.S3InputInfo != nil {
											s3InputInfoMap := map[string]interface{}{}

											if tailSet.S3InputInfo.S3Bucket != nil {
												s3InputInfoMap["s3_bucket"] = tailSet.S3InputInfo.S3Bucket
											}

											if tailSet.S3InputInfo.S3Region != nil {
												s3InputInfoMap["s3_region"] = tailSet.S3InputInfo.S3Region
											}

											if tailSet.S3InputInfo.S3Object != nil {
												s3InputInfoMap["s3_object"] = tailSet.S3InputInfo.S3Object
											}

											if tailSet.S3InputInfo.S3SecretId != nil {
												s3InputInfoMap["s3_secret_id"] = tailSet.S3InputInfo.S3SecretId
											}

											if tailSet.S3InputInfo.S3SecretKey != nil {
												s3InputInfoMap["s3_secret_key"] = tailSet.S3InputInfo.S3SecretKey
											}

											tailSetMap["s3_input_info"] = []interface{}{s3InputInfoMap}
										}

										tailSetList = append(tailSetList, tailSetMap)
									}

									headTailParameterMap["tail_set"] = tailSetList
								}

								transcodeTaskMap["head_tail_parameter"] = []interface{}{headTailParameterMap}
							}

							activityParaMap["transcode_task"] = []interface{}{transcodeTaskMap}
						}

						if activities.ActivityPara.AnimatedGraphicTask != nil {
							animatedGraphicTaskMap := map[string]interface{}{}

							if activities.ActivityPara.AnimatedGraphicTask.Definition != nil {
								animatedGraphicTaskMap["definition"] = activities.ActivityPara.AnimatedGraphicTask.Definition
							}

							if activities.ActivityPara.AnimatedGraphicTask.StartTimeOffset != nil {
								animatedGraphicTaskMap["start_time_offset"] = activities.ActivityPara.AnimatedGraphicTask.StartTimeOffset
							}

							if activities.ActivityPara.AnimatedGraphicTask.EndTimeOffset != nil {
								animatedGraphicTaskMap["end_time_offset"] = activities.ActivityPara.AnimatedGraphicTask.EndTimeOffset
							}

							if activities.ActivityPara.AnimatedGraphicTask.OutputStorage != nil {
								outputStorageMap := map[string]interface{}{}

								if activities.ActivityPara.AnimatedGraphicTask.OutputStorage.Type != nil {
									outputStorageMap["type"] = activities.ActivityPara.AnimatedGraphicTask.OutputStorage.Type
								}

								if activities.ActivityPara.AnimatedGraphicTask.OutputStorage.CosOutputStorage != nil {
									cosOutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.AnimatedGraphicTask.OutputStorage.CosOutputStorage.Bucket != nil {
										cosOutputStorageMap["bucket"] = activities.ActivityPara.AnimatedGraphicTask.OutputStorage.CosOutputStorage.Bucket
									}

									if activities.ActivityPara.AnimatedGraphicTask.OutputStorage.CosOutputStorage.Region != nil {
										cosOutputStorageMap["region"] = activities.ActivityPara.AnimatedGraphicTask.OutputStorage.CosOutputStorage.Region
									}

									outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
								}

								if activities.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage != nil {
									s3OutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3Bucket != nil {
										s3OutputStorageMap["s3_bucket"] = activities.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3Bucket
									}

									if activities.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3Region != nil {
										s3OutputStorageMap["s3_region"] = activities.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3Region
									}

									if activities.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3SecretId != nil {
										s3OutputStorageMap["s3_secret_id"] = activities.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3SecretId
									}

									if activities.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3SecretKey != nil {
										s3OutputStorageMap["s3_secret_key"] = activities.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3SecretKey
									}

									outputStorageMap["s3_output_storage"] = []interface{}{s3OutputStorageMap}
								}

								animatedGraphicTaskMap["output_storage"] = []interface{}{outputStorageMap}
							}

							if activities.ActivityPara.AnimatedGraphicTask.OutputObjectPath != nil {
								animatedGraphicTaskMap["output_object_path"] = activities.ActivityPara.AnimatedGraphicTask.OutputObjectPath
							}

							activityParaMap["animated_graphic_task"] = []interface{}{animatedGraphicTaskMap}
						}

						if activities.ActivityPara.SnapshotByTimeOffsetTask != nil {
							snapshotByTimeOffsetTaskMap := map[string]interface{}{}

							if activities.ActivityPara.SnapshotByTimeOffsetTask.Definition != nil {
								snapshotByTimeOffsetTaskMap["definition"] = activities.ActivityPara.SnapshotByTimeOffsetTask.Definition
							}

							if activities.ActivityPara.SnapshotByTimeOffsetTask.ExtTimeOffsetSet != nil {
								snapshotByTimeOffsetTaskMap["ext_time_offset_set"] = activities.ActivityPara.SnapshotByTimeOffsetTask.ExtTimeOffsetSet
							}

							if activities.ActivityPara.SnapshotByTimeOffsetTask.WatermarkSet != nil {
								watermarkSetList := []interface{}{}
								for _, watermarkSet := range activities.ActivityPara.SnapshotByTimeOffsetTask.WatermarkSet {
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

												if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo != nil {
													s3InputInfoMap := map[string]interface{}{}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Bucket != nil {
														s3InputInfoMap["s3_bucket"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Bucket
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Region != nil {
														s3InputInfoMap["s3_region"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Region
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Object != nil {
														s3InputInfoMap["s3_object"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Object
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretId != nil {
														s3InputInfoMap["s3_secret_id"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretId
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretKey != nil {
														s3InputInfoMap["s3_secret_key"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretKey
													}

													imageContentMap["s3_input_info"] = []interface{}{s3InputInfoMap}
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

								snapshotByTimeOffsetTaskMap["watermark_set"] = watermarkSetList
							}

							if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage != nil {
								outputStorageMap := map[string]interface{}{}

								if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.Type != nil {
									outputStorageMap["type"] = activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.Type
								}

								if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.CosOutputStorage != nil {
									cosOutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.CosOutputStorage.Bucket != nil {
										cosOutputStorageMap["bucket"] = activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.CosOutputStorage.Bucket
									}

									if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.CosOutputStorage.Region != nil {
										cosOutputStorageMap["region"] = activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.CosOutputStorage.Region
									}

									outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
								}

								if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage != nil {
									s3OutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3Bucket != nil {
										s3OutputStorageMap["s3_bucket"] = activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3Bucket
									}

									if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3Region != nil {
										s3OutputStorageMap["s3_region"] = activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3Region
									}

									if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3SecretId != nil {
										s3OutputStorageMap["s3_secret_id"] = activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3SecretId
									}

									if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3SecretKey != nil {
										s3OutputStorageMap["s3_secret_key"] = activities.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3SecretKey
									}

									outputStorageMap["s3_output_storage"] = []interface{}{s3OutputStorageMap}
								}

								snapshotByTimeOffsetTaskMap["output_storage"] = []interface{}{outputStorageMap}
							}

							if activities.ActivityPara.SnapshotByTimeOffsetTask.OutputObjectPath != nil {
								snapshotByTimeOffsetTaskMap["output_object_path"] = activities.ActivityPara.SnapshotByTimeOffsetTask.OutputObjectPath
							}

							if activities.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat != nil {
								objectNumberFormatMap := map[string]interface{}{}

								if activities.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.InitialValue != nil {
									objectNumberFormatMap["initial_value"] = activities.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.InitialValue
								}

								if activities.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.Increment != nil {
									objectNumberFormatMap["increment"] = activities.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.Increment
								}

								if activities.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.MinLength != nil {
									objectNumberFormatMap["min_length"] = activities.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.MinLength
								}

								if activities.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.PlaceHolder != nil {
									objectNumberFormatMap["place_holder"] = activities.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.PlaceHolder
								}

								snapshotByTimeOffsetTaskMap["object_number_format"] = []interface{}{objectNumberFormatMap}
							}

							activityParaMap["snapshot_by_time_offset_task"] = []interface{}{snapshotByTimeOffsetTaskMap}
						}

						if activities.ActivityPara.SampleSnapshotTask != nil {
							sampleSnapshotTaskMap := map[string]interface{}{}

							if activities.ActivityPara.SampleSnapshotTask.Definition != nil {
								sampleSnapshotTaskMap["definition"] = activities.ActivityPara.SampleSnapshotTask.Definition
							}

							if activities.ActivityPara.SampleSnapshotTask.WatermarkSet != nil {
								watermarkSetList := []interface{}{}
								for _, watermarkSet := range activities.ActivityPara.SampleSnapshotTask.WatermarkSet {
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

												if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo != nil {
													s3InputInfoMap := map[string]interface{}{}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Bucket != nil {
														s3InputInfoMap["s3_bucket"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Bucket
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Region != nil {
														s3InputInfoMap["s3_region"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Region
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Object != nil {
														s3InputInfoMap["s3_object"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Object
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretId != nil {
														s3InputInfoMap["s3_secret_id"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretId
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretKey != nil {
														s3InputInfoMap["s3_secret_key"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretKey
													}

													imageContentMap["s3_input_info"] = []interface{}{s3InputInfoMap}
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

								sampleSnapshotTaskMap["watermark_set"] = watermarkSetList
							}

							if activities.ActivityPara.SampleSnapshotTask.OutputStorage != nil {
								outputStorageMap := map[string]interface{}{}

								if activities.ActivityPara.SampleSnapshotTask.OutputStorage.Type != nil {
									outputStorageMap["type"] = activities.ActivityPara.SampleSnapshotTask.OutputStorage.Type
								}

								if activities.ActivityPara.SampleSnapshotTask.OutputStorage.CosOutputStorage != nil {
									cosOutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.SampleSnapshotTask.OutputStorage.CosOutputStorage.Bucket != nil {
										cosOutputStorageMap["bucket"] = activities.ActivityPara.SampleSnapshotTask.OutputStorage.CosOutputStorage.Bucket
									}

									if activities.ActivityPara.SampleSnapshotTask.OutputStorage.CosOutputStorage.Region != nil {
										cosOutputStorageMap["region"] = activities.ActivityPara.SampleSnapshotTask.OutputStorage.CosOutputStorage.Region
									}

									outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
								}

								if activities.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage != nil {
									s3OutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3Bucket != nil {
										s3OutputStorageMap["s3_bucket"] = activities.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3Bucket
									}

									if activities.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3Region != nil {
										s3OutputStorageMap["s3_region"] = activities.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3Region
									}

									if activities.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3SecretId != nil {
										s3OutputStorageMap["s3_secret_id"] = activities.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3SecretId
									}

									if activities.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3SecretKey != nil {
										s3OutputStorageMap["s3_secret_key"] = activities.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3SecretKey
									}

									outputStorageMap["s3_output_storage"] = []interface{}{s3OutputStorageMap}
								}

								sampleSnapshotTaskMap["output_storage"] = []interface{}{outputStorageMap}
							}

							if activities.ActivityPara.SampleSnapshotTask.OutputObjectPath != nil {
								sampleSnapshotTaskMap["output_object_path"] = activities.ActivityPara.SampleSnapshotTask.OutputObjectPath
							}

							if activities.ActivityPara.SampleSnapshotTask.ObjectNumberFormat != nil {
								objectNumberFormatMap := map[string]interface{}{}

								if activities.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.InitialValue != nil {
									objectNumberFormatMap["initial_value"] = activities.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.InitialValue
								}

								if activities.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.Increment != nil {
									objectNumberFormatMap["increment"] = activities.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.Increment
								}

								if activities.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.MinLength != nil {
									objectNumberFormatMap["min_length"] = activities.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.MinLength
								}

								if activities.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.PlaceHolder != nil {
									objectNumberFormatMap["place_holder"] = activities.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.PlaceHolder
								}

								sampleSnapshotTaskMap["object_number_format"] = []interface{}{objectNumberFormatMap}
							}

							activityParaMap["sample_snapshot_task"] = []interface{}{sampleSnapshotTaskMap}
						}

						if activities.ActivityPara.ImageSpriteTask != nil {
							imageSpriteTaskMap := map[string]interface{}{}

							if activities.ActivityPara.ImageSpriteTask.Definition != nil {
								imageSpriteTaskMap["definition"] = activities.ActivityPara.ImageSpriteTask.Definition
							}

							if activities.ActivityPara.ImageSpriteTask.OutputStorage != nil {
								outputStorageMap := map[string]interface{}{}

								if activities.ActivityPara.ImageSpriteTask.OutputStorage.Type != nil {
									outputStorageMap["type"] = activities.ActivityPara.ImageSpriteTask.OutputStorage.Type
								}

								if activities.ActivityPara.ImageSpriteTask.OutputStorage.CosOutputStorage != nil {
									cosOutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.ImageSpriteTask.OutputStorage.CosOutputStorage.Bucket != nil {
										cosOutputStorageMap["bucket"] = activities.ActivityPara.ImageSpriteTask.OutputStorage.CosOutputStorage.Bucket
									}

									if activities.ActivityPara.ImageSpriteTask.OutputStorage.CosOutputStorage.Region != nil {
										cosOutputStorageMap["region"] = activities.ActivityPara.ImageSpriteTask.OutputStorage.CosOutputStorage.Region
									}

									outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
								}

								if activities.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage != nil {
									s3OutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3Bucket != nil {
										s3OutputStorageMap["s3_bucket"] = activities.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3Bucket
									}

									if activities.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3Region != nil {
										s3OutputStorageMap["s3_region"] = activities.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3Region
									}

									if activities.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3SecretId != nil {
										s3OutputStorageMap["s3_secret_id"] = activities.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3SecretId
									}

									if activities.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3SecretKey != nil {
										s3OutputStorageMap["s3_secret_key"] = activities.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3SecretKey
									}

									outputStorageMap["s3_output_storage"] = []interface{}{s3OutputStorageMap}
								}

								imageSpriteTaskMap["output_storage"] = []interface{}{outputStorageMap}
							}

							if activities.ActivityPara.ImageSpriteTask.OutputObjectPath != nil {
								imageSpriteTaskMap["output_object_path"] = activities.ActivityPara.ImageSpriteTask.OutputObjectPath
							}

							if activities.ActivityPara.ImageSpriteTask.WebVttObjectName != nil {
								imageSpriteTaskMap["web_vtt_object_name"] = activities.ActivityPara.ImageSpriteTask.WebVttObjectName
							}

							if activities.ActivityPara.ImageSpriteTask.ObjectNumberFormat != nil {
								objectNumberFormatMap := map[string]interface{}{}

								if activities.ActivityPara.ImageSpriteTask.ObjectNumberFormat.InitialValue != nil {
									objectNumberFormatMap["initial_value"] = activities.ActivityPara.ImageSpriteTask.ObjectNumberFormat.InitialValue
								}

								if activities.ActivityPara.ImageSpriteTask.ObjectNumberFormat.Increment != nil {
									objectNumberFormatMap["increment"] = activities.ActivityPara.ImageSpriteTask.ObjectNumberFormat.Increment
								}

								if activities.ActivityPara.ImageSpriteTask.ObjectNumberFormat.MinLength != nil {
									objectNumberFormatMap["min_length"] = activities.ActivityPara.ImageSpriteTask.ObjectNumberFormat.MinLength
								}

								if activities.ActivityPara.ImageSpriteTask.ObjectNumberFormat.PlaceHolder != nil {
									objectNumberFormatMap["place_holder"] = activities.ActivityPara.ImageSpriteTask.ObjectNumberFormat.PlaceHolder
								}

								imageSpriteTaskMap["object_number_format"] = []interface{}{objectNumberFormatMap}
							}

							activityParaMap["image_sprite_task"] = []interface{}{imageSpriteTaskMap}
						}

						if activities.ActivityPara.AdaptiveDynamicStreamingTask != nil {
							adaptiveDynamicStreamingTaskMap := map[string]interface{}{}

							if activities.ActivityPara.AdaptiveDynamicStreamingTask.Definition != nil {
								adaptiveDynamicStreamingTaskMap["definition"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.Definition
							}

							if activities.ActivityPara.AdaptiveDynamicStreamingTask.WatermarkSet != nil {
								watermarkSetList := []interface{}{}
								for _, watermarkSet := range activities.ActivityPara.AdaptiveDynamicStreamingTask.WatermarkSet {
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

												if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo != nil {
													s3InputInfoMap := map[string]interface{}{}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Bucket != nil {
														s3InputInfoMap["s3_bucket"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Bucket
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Region != nil {
														s3InputInfoMap["s3_region"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Region
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Object != nil {
														s3InputInfoMap["s3_object"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3Object
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretId != nil {
														s3InputInfoMap["s3_secret_id"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretId
													}

													if watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretKey != nil {
														s3InputInfoMap["s3_secret_key"] = watermarkSet.RawParameter.ImageTemplate.ImageContent.S3InputInfo.S3SecretKey
													}

													imageContentMap["s3_input_info"] = []interface{}{s3InputInfoMap}
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

								adaptiveDynamicStreamingTaskMap["watermark_set"] = watermarkSetList
							}

							if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage != nil {
								outputStorageMap := map[string]interface{}{}

								if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.Type != nil {
									outputStorageMap["type"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.Type
								}

								if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.CosOutputStorage != nil {
									cosOutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.CosOutputStorage.Bucket != nil {
										cosOutputStorageMap["bucket"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.CosOutputStorage.Bucket
									}

									if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.CosOutputStorage.Region != nil {
										cosOutputStorageMap["region"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.CosOutputStorage.Region
									}

									outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
								}

								if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage != nil {
									s3OutputStorageMap := map[string]interface{}{}

									if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3Bucket != nil {
										s3OutputStorageMap["s3_bucket"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3Bucket
									}

									if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3Region != nil {
										s3OutputStorageMap["s3_region"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3Region
									}

									if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3SecretId != nil {
										s3OutputStorageMap["s3_secret_id"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3SecretId
									}

									if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3SecretKey != nil {
										s3OutputStorageMap["s3_secret_key"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3SecretKey
									}

									outputStorageMap["s3_output_storage"] = []interface{}{s3OutputStorageMap}
								}

								adaptiveDynamicStreamingTaskMap["output_storage"] = []interface{}{outputStorageMap}
							}

							if activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputObjectPath != nil {
								adaptiveDynamicStreamingTaskMap["output_object_path"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.OutputObjectPath
							}

							if activities.ActivityPara.AdaptiveDynamicStreamingTask.SubStreamObjectName != nil {
								adaptiveDynamicStreamingTaskMap["sub_stream_object_name"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.SubStreamObjectName
							}

							if activities.ActivityPara.AdaptiveDynamicStreamingTask.SegmentObjectName != nil {
								adaptiveDynamicStreamingTaskMap["segment_object_name"] = activities.ActivityPara.AdaptiveDynamicStreamingTask.SegmentObjectName
							}

							if activities.ActivityPara.AdaptiveDynamicStreamingTask.AddOnSubtitles != nil {
								addOnSubtitlesList := []interface{}{}
								for _, addOnSubtitles := range activities.ActivityPara.AdaptiveDynamicStreamingTask.AddOnSubtitles {
									addOnSubtitlesMap := map[string]interface{}{}

									if addOnSubtitles.Type != nil {
										addOnSubtitlesMap["type"] = addOnSubtitles.Type
									}

									if addOnSubtitles.Subtitle != nil {
										subtitleMap := map[string]interface{}{}

										if addOnSubtitles.Subtitle.Type != nil {
											subtitleMap["type"] = addOnSubtitles.Subtitle.Type
										}

										if addOnSubtitles.Subtitle.CosInputInfo != nil {
											cosInputInfoMap := map[string]interface{}{}

											if addOnSubtitles.Subtitle.CosInputInfo.Bucket != nil {
												cosInputInfoMap["bucket"] = addOnSubtitles.Subtitle.CosInputInfo.Bucket
											}

											if addOnSubtitles.Subtitle.CosInputInfo.Region != nil {
												cosInputInfoMap["region"] = addOnSubtitles.Subtitle.CosInputInfo.Region
											}

											if addOnSubtitles.Subtitle.CosInputInfo.Object != nil {
												cosInputInfoMap["object"] = addOnSubtitles.Subtitle.CosInputInfo.Object
											}

											subtitleMap["cos_input_info"] = []interface{}{cosInputInfoMap}
										}

										if addOnSubtitles.Subtitle.UrlInputInfo != nil {
											urlInputInfoMap := map[string]interface{}{}

											if addOnSubtitles.Subtitle.UrlInputInfo.Url != nil {
												urlInputInfoMap["url"] = addOnSubtitles.Subtitle.UrlInputInfo.Url
											}

											subtitleMap["url_input_info"] = []interface{}{urlInputInfoMap}
										}

										if addOnSubtitles.Subtitle.S3InputInfo != nil {
											s3InputInfoMap := map[string]interface{}{}

											if addOnSubtitles.Subtitle.S3InputInfo.S3Bucket != nil {
												s3InputInfoMap["s3_bucket"] = addOnSubtitles.Subtitle.S3InputInfo.S3Bucket
											}

											if addOnSubtitles.Subtitle.S3InputInfo.S3Region != nil {
												s3InputInfoMap["s3_region"] = addOnSubtitles.Subtitle.S3InputInfo.S3Region
											}

											if addOnSubtitles.Subtitle.S3InputInfo.S3Object != nil {
												s3InputInfoMap["s3_object"] = addOnSubtitles.Subtitle.S3InputInfo.S3Object
											}

											if addOnSubtitles.Subtitle.S3InputInfo.S3SecretId != nil {
												s3InputInfoMap["s3_secret_id"] = addOnSubtitles.Subtitle.S3InputInfo.S3SecretId
											}

											if addOnSubtitles.Subtitle.S3InputInfo.S3SecretKey != nil {
												s3InputInfoMap["s3_secret_key"] = addOnSubtitles.Subtitle.S3InputInfo.S3SecretKey
											}

											subtitleMap["s3_input_info"] = []interface{}{s3InputInfoMap}
										}

										addOnSubtitlesMap["subtitle"] = []interface{}{subtitleMap}
									}

									addOnSubtitlesList = append(addOnSubtitlesList, addOnSubtitlesMap)
								}

								adaptiveDynamicStreamingTaskMap["add_on_subtitles"] = addOnSubtitlesList
							}

							activityParaMap["adaptive_dynamic_streaming_task"] = []interface{}{adaptiveDynamicStreamingTaskMap}
						}

						if activities.ActivityPara.AiContentReviewTask != nil {
							aiContentReviewTaskMap := map[string]interface{}{}

							if activities.ActivityPara.AiContentReviewTask.Definition != nil {
								aiContentReviewTaskMap["definition"] = activities.ActivityPara.AiContentReviewTask.Definition
							}

							activityParaMap["ai_content_review_task"] = []interface{}{aiContentReviewTaskMap}
						}

						if activities.ActivityPara.AiAnalysisTask != nil {
							aiAnalysisTaskMap := map[string]interface{}{}

							if activities.ActivityPara.AiAnalysisTask.Definition != nil {
								aiAnalysisTaskMap["definition"] = activities.ActivityPara.AiAnalysisTask.Definition
							}

							if activities.ActivityPara.AiAnalysisTask.ExtendedParameter != nil {
								aiAnalysisTaskMap["extended_parameter"] = activities.ActivityPara.AiAnalysisTask.ExtendedParameter
							}

							activityParaMap["ai_analysis_task"] = []interface{}{aiAnalysisTaskMap}
						}

						if activities.ActivityPara.AiRecognitionTask != nil {
							aiRecognitionTaskMap := map[string]interface{}{}

							if activities.ActivityPara.AiRecognitionTask.Definition != nil {
								aiRecognitionTaskMap["definition"] = activities.ActivityPara.AiRecognitionTask.Definition
							}

							activityParaMap["ai_recognition_task"] = []interface{}{aiRecognitionTaskMap}
						}

						activitiesMap["activity_para"] = []interface{}{activityParaMap}
					}

					activitiesList = append(activitiesList, activitiesMap)
				}

				schedulesInfoMap["activities"] = activitiesList
			}

			if schedulesInfo.OutputStorage != nil {
				outputStorageMap := map[string]interface{}{}

				if schedulesInfo.OutputStorage.Type != nil {
					outputStorageMap["type"] = schedulesInfo.OutputStorage.Type
				}

				if schedulesInfo.OutputStorage.CosOutputStorage != nil {
					cosOutputStorageMap := map[string]interface{}{}

					if schedulesInfo.OutputStorage.CosOutputStorage.Bucket != nil {
						cosOutputStorageMap["bucket"] = schedulesInfo.OutputStorage.CosOutputStorage.Bucket
					}

					if schedulesInfo.OutputStorage.CosOutputStorage.Region != nil {
						cosOutputStorageMap["region"] = schedulesInfo.OutputStorage.CosOutputStorage.Region
					}

					outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
				}

				if schedulesInfo.OutputStorage.S3OutputStorage != nil {
					s3OutputStorageMap := map[string]interface{}{}

					if schedulesInfo.OutputStorage.S3OutputStorage.S3Bucket != nil {
						s3OutputStorageMap["s3_bucket"] = schedulesInfo.OutputStorage.S3OutputStorage.S3Bucket
					}

					if schedulesInfo.OutputStorage.S3OutputStorage.S3Region != nil {
						s3OutputStorageMap["s3_region"] = schedulesInfo.OutputStorage.S3OutputStorage.S3Region
					}

					if schedulesInfo.OutputStorage.S3OutputStorage.S3SecretId != nil {
						s3OutputStorageMap["s3_secret_id"] = schedulesInfo.OutputStorage.S3OutputStorage.S3SecretId
					}

					if schedulesInfo.OutputStorage.S3OutputStorage.S3SecretKey != nil {
						s3OutputStorageMap["s3_secret_key"] = schedulesInfo.OutputStorage.S3OutputStorage.S3SecretKey
					}

					outputStorageMap["s3_output_storage"] = []interface{}{s3OutputStorageMap}
				}

				schedulesInfoMap["output_storage"] = []interface{}{outputStorageMap}
			}

			if schedulesInfo.OutputDir != nil {
				schedulesInfoMap["output_dir"] = schedulesInfo.OutputDir
			}

			if schedulesInfo.TaskNotifyConfig != nil {
				taskNotifyConfigMap := map[string]interface{}{}

				if schedulesInfo.TaskNotifyConfig.CmqModel != nil {
					taskNotifyConfigMap["cmq_model"] = schedulesInfo.TaskNotifyConfig.CmqModel
				}

				if schedulesInfo.TaskNotifyConfig.CmqRegion != nil {
					taskNotifyConfigMap["cmq_region"] = schedulesInfo.TaskNotifyConfig.CmqRegion
				}

				if schedulesInfo.TaskNotifyConfig.TopicName != nil {
					taskNotifyConfigMap["topic_name"] = schedulesInfo.TaskNotifyConfig.TopicName
				}

				if schedulesInfo.TaskNotifyConfig.QueueName != nil {
					taskNotifyConfigMap["queue_name"] = schedulesInfo.TaskNotifyConfig.QueueName
				}

				if schedulesInfo.TaskNotifyConfig.NotifyMode != nil {
					taskNotifyConfigMap["notify_mode"] = schedulesInfo.TaskNotifyConfig.NotifyMode
				}

				if schedulesInfo.TaskNotifyConfig.NotifyType != nil {
					taskNotifyConfigMap["notify_type"] = schedulesInfo.TaskNotifyConfig.NotifyType
				}

				if schedulesInfo.TaskNotifyConfig.NotifyUrl != nil {
					taskNotifyConfigMap["notify_url"] = schedulesInfo.TaskNotifyConfig.NotifyUrl
				}

				if schedulesInfo.TaskNotifyConfig.AwsSQS != nil {
					awsSQSMap := map[string]interface{}{}

					if schedulesInfo.TaskNotifyConfig.AwsSQS.SQSRegion != nil {
						awsSQSMap["sqs_region"] = schedulesInfo.TaskNotifyConfig.AwsSQS.SQSRegion
					}

					if schedulesInfo.TaskNotifyConfig.AwsSQS.SQSQueueName != nil {
						awsSQSMap["sqs_queue_name"] = schedulesInfo.TaskNotifyConfig.AwsSQS.SQSQueueName
					}

					if schedulesInfo.TaskNotifyConfig.AwsSQS.S3SecretId != nil {
						awsSQSMap["s3_secret_id"] = schedulesInfo.TaskNotifyConfig.AwsSQS.S3SecretId
					}

					if schedulesInfo.TaskNotifyConfig.AwsSQS.S3SecretKey != nil {
						awsSQSMap["s3_secret_key"] = schedulesInfo.TaskNotifyConfig.AwsSQS.S3SecretKey
					}

					taskNotifyConfigMap["aws_sqs"] = []interface{}{awsSQSMap}
				}

				schedulesInfoMap["task_notify_config"] = []interface{}{taskNotifyConfigMap}
			}

			if schedulesInfo.CreateTime != nil {
				schedulesInfoMap["create_time"] = schedulesInfo.CreateTime
			}

			if schedulesInfo.UpdateTime != nil {
				schedulesInfoMap["update_time"] = schedulesInfo.UpdateTime
			}

			if schedulesInfo.ResourceId != nil {
				schedulesInfoMap["resource_id"] = schedulesInfo.ResourceId
			}

			ids = append(ids, helper.Int64ToStr(*schedulesInfo.ScheduleId))
			tmpList = append(tmpList, schedulesInfoMap)
		}

		_ = d.Set("schedule_info_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
