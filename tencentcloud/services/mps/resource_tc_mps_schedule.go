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

func ResourceTencentCloudMpsSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsScheduleCreate,
		Read:   resourceTencentCloudMpsScheduleRead,
		Update: resourceTencentCloudMpsScheduleUpdate,
		Delete: resourceTencentCloudMpsScheduleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"schedule_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "scheme 名称 (max 128 字符). 此 名称 should 是 唯一 across your 账号",
			},

			"trigger": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "触发器 的 scheme. 如果 文件 是 uploaded 到 指定 存储桶， scheme 将 是 triggered。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "触发器 类型 有效值：`CosFileUpload`: Tencent Cloud COS 触发器. `AwsS3FileUpload`: AWS S3 触发器. Currently，此 类型 是 仅 支持 对于 transcoding tasks 和 schemes (不 支持 对于 workflows)。",
						},
						"cos_file_upload_trigger": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "此 参数 为必填项 和 有效 当 `类型` 是 `CosFileUpload`，indicating COS 触发器 规则.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "名称 COS 存储桶 bound 到 工作流，such 作为 `TopRankVideo-125xxx88`。",
									},
									"region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "地域 的 COS 存储桶 bound 到 工作流，such 作为 `ap-chongiqng`。",
									},
									"dir": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Input 路径 directory bound 到 工作流，such 作为 `/movie/201907/`. 如果此参数为空， `/` root directory 将 是 使用。",
									},
									"formats": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "格式 列表 files 该 可以 触发器 工作流，such 作为 [mp4，flv，mov]. 如果此参数为空，files 在 all formats 可以 触发器 工作流。",
									},
								},
							},
						},
						"aws_s3_file_upload_trigger": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "AWS S3 触发器. 此 参数 是 有效 和 必填 如果 `类型` 是 `AwsS3FileUpload`.注意: Currently， 键 对于 AWS S3 存储桶， 触发器 SQS queue，和 callback SQS queue 必须 是 same.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"s3_bucket": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "AWS S3 存储桶 bound 到 scheme。",
									},
									"s3_region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "地域 的 AWS S3 存储桶",
									},
									"dir": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "存储桶 directory bound. It 必须 是 absolute 路径 该 starts 和 结束 使用 `/`，such 作为 `/movie/201907/`. 如果 您 do 不 指定this， root directory 将 是 bound.	。",
									},
									"formats": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "文件 formats 该 将 触发器 scheme，such 作为 [mp4，flv，mov]. 如果 您 do 不 指定this， upload 的 files 在 any 格式 将 触发器 scheme.	。",
									},
									"s3_secret_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "键 ID AWS S3 存储桶注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"s3_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "键 的 AWS S3 存储桶注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"aws_sqs": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "SQS queue 的 AWS S3 存储桶Note: queue 必须 是 在 same 地域 作为 存储桶注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"sqs_region": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "地域 的 SQS queue。",
												},
												"sqs_queue_name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "名称 SQS queue。",
												},
												"s3_secret_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "键 ID 必填 到 read 从/write 到 SQS queue。",
												},
												"s3_secret_key": {
													Type:        schema.TypeString,
													Optional:    true,
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
				Required:    true,
				Type:        schema.TypeList,
				Description: "subtasks 的 scheme。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"activity_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "subtask 类型 `input`: start. `output`: end. `操作-trans`: Transcoding. `操作-samplesnapshot`: Sampled screencapturing. `操作-AIAnalysis`: 内容 analysis. `操作-AIRecognition`: 内容 recognition. `操作-aiReview`: 内容 moderation. `操作-animated-graphics`: Animated screenshot generation. `操作-镜像-sprite`: Image sprite generation. `操作-snapshotByTimeOffset`: Time point screencapturing. `操作-adaptive-substream`: Adaptive bitrate streaming.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"reardrive_index": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Optional:    true,
							Computed:    true,
							Description: "indexes 的 subsequent actions. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"activity_para": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "参数 的 subtask.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"transcode_task": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "A transcoding 任务。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "ID 视频 transcoding template。",
												},
												"raw_parameter": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Computed:    true,
													Description: "Custom 视频 transcoding 参数，其中 是 有效 如果 `Definition` 是 0.此 参数 是 使用 在 highly customized scenarios. We recommend 您 使用 `Definition` 到 指定transcoding 参数 preferably。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"container": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Container. 有效值：mp4; flv; hls; mp3; flac; ogg; m4a. Among them，mp3，flac，ogg，和 m4a 是 对于 音频 files。",
															},
															"remove_video": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "是否remove 视频 数据. 有效值：0: retain; 1: remove.默认值：0。",
															},
															"remove_audio": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "是否remove 音频 数据. 有效值：0: retain; 1: remove.默认值：0。",
															},
															"video_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Computed:    true,
																Description: "Video 流 配置 参数. 此 字段 为必填项 当 `RemoveVideo` 是 0。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"codec": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "视频 codec. 有效值：`libx264`: H.264 `libx265`: H.265 `av1`: AOMedia Video 1Note: You 必须 指定a resolution (不 higher 比 640 x 480) 如果 H.265 codec 是 使用.注意: You 可以 仅 使用 AOMedia Video 1 codec 对于 MP4 files。",
																		},
																		"fps": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "视频 frame 速率 (Hz). 取值范围：[0，100].如果 值 是 0， frame 速率 将 是 same 作为 该 的 来源 视频.注意: For adaptive bitrate streaming， 值 范围 的 此 参数 是 [0，60]。",
																		},
																		"bitrate": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "视频 bitrate (Kbps). 取值范围：0 和 [128，35000].如果 值 是 0， bitrate 的 视频 将 是 same 作为 该 的 来源 视频。",
																		},
																		"resolution_adaptive": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Resolution adaption. 有效值：open: 已启用 当 resolution adaption 是 已启用，`宽度` 表示long side 的 视频，while `高度` 表示short side. close: 已禁用 当 resolution adaption 是 已禁用，`宽度` 表示width 的 视频，while `高度` 表示height.默认值：open.注意: 当 resolution adaption 是 已启用，`宽度` 不能 是 smaller 比 `高度`。",
																		},
																		"width": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Maximum 值 的 宽度 (或 long side) 的 视频 流 （像素）。 取值范围：0 和 [128，4,096]. 如果 both `宽度` 和 `高度` 是 0， resolution 将 是 same 作为 该 的 来源 视频; 如果 `宽度` 是 0，但 `高度` 是 不 0，`宽度` 将 是 proportionally scaled; 如果 `宽度` 是 不 0，但 `高度` 是 0，`高度` 将 是 proportionally scaled; 如果 both `宽度` 和 `高度` 是 不 0， 自定义 resolution 将 是 使用.默认值：0。",
																		},
																		"height": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Maximum 值 的 高度 (或 short side) 的 视频 流 （像素）。 取值范围：0 和 [128，4,096]. 如果 both `宽度` 和 `高度` 是 0， resolution 将 是 same 作为 该 的 来源 视频; 如果 `宽度` 是 0，但 `高度` 是 不 0，`宽度` 将 是 proportionally scaled; 如果 `宽度` 是 不 0，但 `高度` 是 0，`高度` 将 是 proportionally scaled; 如果 both `宽度` 和 `高度` 是 不 0， 自定义 resolution 将 是 使用.默认值：0。",
																		},
																		"gop": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Frame 间隔 between I keyframes. 取值范围：0 和 [1,100000].如果 此 参数 是 0 或 left 空， 系统 将 automatically 集合 GOP 长度。",
																		},
																		"fill_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "fill 模式，其中 表示how 视频 是 resized 当 视频's original aspect ratio 是 different 从 目标 aspect ratio. 有效值：stretch: Stretch 镜像 frame 通过 frame 到 fill entire screen. 视频 镜像 可能 become squashed 或 stretched after transcoding. black: Keep 镜像&#39;s original aspect ratio 和 fill blank space 使用 black bars. white: Keep 镜像's original aspect ratio 和 fill blank space 使用 white bars. gauss: Keep 镜像's original aspect ratio 和 apply Gaussian blur 到 blank space.默认值：black.注意: Only `stretch` 和 `black` 是 支持 对于 adaptive bitrate streaming。",
																		},
																		"vcrf": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "control factor 的 视频 constant bitrate. 取值范围：[1，51]如果 此 参数 是 指定，CRF ( bitrate control 方法) 将 是 用于transcoding. (Video bitrate 将 无 longer take effect.)It 是 不 recommended 到 指定this 参数 如果 there 是 无 special requirements。",
																		},
																	},
																},
															},
															"audio_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Computed:    true,
																Description: "Audio 流 配置 参数. 此 字段 为必填项 当 `RemoveAudio` 是 0。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"codec": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Audio 流 codec.当 outer `Container` 参数 是 `mp3`， 有效 值 是: libmp3lame.当 outer `Container` 参数 是 `ogg` 或 `flac`， 有效 值 是: flac.当 outer `Container` 参数 是 `m4a`， 有效 值 include: libfdk_aac; libmp3lame; ac3.当 outer `Container` 参数 是 `mp4` 或 `flv`， 有效 值 include: libfdk_aac: more suitable 对于 mp4; libmp3lame: more suitable 对于 flv.当 outer `Container` 参数 是 `hls`， 有效 值 include: libfdk_aac; libmp3lame。",
																		},
																		"bitrate": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Audio 流 bitrate 在 Kbps. 取值范围：0 和 [26，256].如果 值 是 0， bitrate 的 音频 流 将 是 same 作为 该 的 original 音频。",
																		},
																		"sample_rate": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Audio 流 sample 速率. 有效值：32,000 44,100 48,000In Hz。",
																		},
																		"audio_channel": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Audio channel 系统. 有效值：1: Mono 2: Dual 6: StereoWhen media 是 packaged 在 音频 格式 (FLAC，OGG，MP3，M4A)， sound channel 不能 是 集合 到 stereo.默认值：2。",
																		},
																	},
																},
															},
															"tehd_config": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Computed:    true,
																Description: "TESHD transcoding 参数。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "TESHD 类型 有效值：TEHD-100: TESHD-100.如果此参数为空，TESHD 将 不 是 已启用",
																		},
																		"max_video_bitrate": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Maximum bitrate，其中 是 有效 当 `类型` 是 `TESHD`.如果此参数为空 或 0 是 entered，there 将 是 无 upper 限制 对于 bitrate。",
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
													Description: "Video transcoding 自定义 参数，其中 是 有效 当 `Definition` 是 不 0.当 any 参数 在 此 structure 是 entered，they 将 是 用于override corresponding 参数 在 templates.此 参数 是 使用 在 highly customized scenarios. We recommend 您 仅 使用 `Definition` 到 指定transcoding 参数.注意: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 是 found。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"container": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Container 格式 有效值：mp4，flv，hls，mp3，flac，ogg，和 m4a; mp3，flac，ogg，和 m4a 是 formats 的 音频 files。",
															},
															"remove_video": {
																Type:        schema.TypeInt,
																Optional:    true,
																Computed:    true,
																Description: "是否remove 视频 数据. 有效值：0: retain 1: remove。",
															},
															"remove_audio": {
																Type:        schema.TypeInt,
																Optional:    true,
																Computed:    true,
																Description: "是否remove 音频 数据. 有效值：0: retain 1: remove。",
															},
															"video_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Computed:    true,
																Description: "Video 流 配置 参数。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"codec": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "视频 codec. 有效值：libx264: H.264 libx265: H.265 av1: AOMedia Video 1Note: You 必须 指定a resolution (不 higher 比 640 x 480) 如果 H.265 codec 是 使用.注意: You 可以 仅 使用 AOMedia Video 1 codec 对于 MP4 files。",
																		},
																		"fps": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Computed:    true,
																			Description: "Video frame 速率 在 Hz. 取值范围：[0，100].如果 值 是 0， frame 速率 将 是 same 作为 该 的 来源 视频。",
																		},
																		"bitrate": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Computed:    true,
																			Description: "Bitrate 的 视频 流 在 Kbps. 取值范围：0 和 [128，35,000].如果 值 是 0， bitrate 的 视频 将 是 same 作为 该 的 来源 视频。",
																		},
																		"resolution_adaptive": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Resolution adaption. 有效值：open: 已启用 当 resolution adaption 是 已启用，`宽度` 表示long side 的 视频，while `高度` 表示short side. close: 已禁用 当 resolution adaption 是 已禁用，`宽度` 表示width 的 视频，while `高度` 表示height.注意: 当 resolution adaption 是 已启用，`宽度` 不能 是 smaller 比 `高度`。",
																		},
																		"width": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Computed:    true,
																			Description: "Maximum 值 的 宽度 (或 long side) 的 视频 流 （像素）。 取值范围：0 和 [128，4,096]. 如果 both `宽度` 和 `高度` 是 0， resolution 将 是 same 作为 该 的 来源 视频; 如果 `宽度` 是 0，但 `高度` 是 不 0，`宽度` 将 是 proportionally scaled; 如果 `宽度` 是 不 0，但 `高度` 是 0，`高度` 将 是 proportionally scaled; 如果 both `宽度` 和 `高度` 是 不 0， 自定义 resolution 将 是 使用。",
																		},
																		"height": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Computed:    true,
																			Description: "Maximum 值 的 高度 (或 short side) 的 视频 流 （像素）。 取值范围：0 和 [128，4,096]。",
																		},
																		"gop": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Computed:    true,
																			Description: "Frame 间隔 between I keyframes. 取值范围：0 和 [1,100000]. 如果 此 参数 是 0， 系统 将 automatically 集合 GOP 长度。",
																		},
																		"fill_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Fill 类型 Fill refers 到 way 的 processing screenshot 当 its aspect ratio 是 different 从 该 的 来源 视频. following fill types 是 支持: stretch: stretch. screenshot 将 是 stretched frame 通过 frame 到 match aspect ratio 的 来源 视频，其中 可能 make screenshot shorter 或 longer; black: fill 使用 black. 此 选项 retains aspect ratio 的 来源 视频 对于 screenshot 和 fills unmatched area 使用 black color blocks. white: fill 使用 white. 此 选项 retains aspect ratio 的 来源 视频 对于 screenshot 和 fills unmatched area 使用 white color blocks. gauss: fill 使用 Gaussian blur. 此 选项 retains aspect ratio 的 来源 视频 对于 screenshot 和 fills unmatched area 使用 Gaussian blur。",
																		},
																		"vcrf": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Computed:    true,
																			Description: "control factor 的 视频 constant bitrate. 取值范围：[0，51]. 此 参数 将 是 已禁用 如果 您 enter `0`.It 是 不 recommended 到 指定this 参数 如果 there 是 无 special requirements。",
																		},
																		"content_adapt_stream": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "是否enable adaptive 编码. 有效值：0: Disable 1: Enable默认值：0. 如果 此 参数 是 集合 到 `1`，多个 streams 使用 different resolutions 和 bitrates 将 是 generated automatically. highest resolution，bitrate，和 quality 的 streams 是 determined 通过 值 的 `宽度` 和 `高度`，`Bitrate`，和 `Vcrf` 在 `VideoTemplate` respectively. 如果 these 参数 是 不 集合 在 `VideoTemplate`， highest resolution generated 将 是 same 作为 该 的 来源 视频，和 highest 视频 quality 将 是 close 到 VMAF 95. To 使用 此 参数 或 learn about billing details 的 adaptive 编码，please contact your sales rep。",
																		},
																	},
																},
															},
															"audio_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Computed:    true,
																Description: "Audio 流 配置 参数。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"codec": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Audio 流 codec.当 outer `Container` 参数 是 `mp3`， 有效 值 是: libmp3lame.当 outer `Container` 参数 是 `ogg` 或 `flac`， 有效 值 是: flac.当 outer `Container` 参数 是 `m4a`， 有效 值 include: libfdk_aac; libmp3lame; ac3.当 outer `Container` 参数 是 `mp4` 或 `flv`， 有效 值 include: libfdk_aac: More suitable 对于 mp4; libmp3lame: More suitable 对于 flv; mp2.当 outer `Container` 参数 是 `hls`， 有效 值 include: libfdk_aac; libmp3lame。",
																		},
																		"bitrate": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Computed:    true,
																			Description: "Audio 流 bitrate 在 Kbps. 取值范围：0 和 [26，256]. 如果 值 是 0， bitrate 的 音频 流 将 是 same 作为 该 的 original 音频。",
																		},
																		"sample_rate": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Computed:    true,
																			Description: "Audio 流 sample 速率. 有效值：32,000 44,100 48,000In Hz。",
																		},
																		"audio_channel": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Computed:    true,
																			Description: "Audio channel 系统. 有效值：1: Mono 2: Dual 6: StereoWhen media 是 packaged 在 音频 格式 (FLAC，OGG，MP3，M4A)， sound channel 不能 是 集合 到 stereo。",
																		},
																		"stream_selects": {
																			Type: schema.TypeSet,
																			Elem: &schema.Schema{
																				Type: schema.TypeInt,
																			},
																			Optional:    true,
																			Description: "音频 tracks 到 retain. All 音频 tracks 是 retained 通过 默认值。",
																		},
																	},
																},
															},
															"tehd_config": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Computed:    true,
																Description: "TESHD transcoding 参数。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "TESHD 类型 有效值：TEHD-100: TESHD-100.如果 此 参数 是 left blank，无 modification 将 是 made。",
																		},
																		"max_video_bitrate": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Computed:    true,
																			Description: "Maximum bitrate. 如果此参数为空，无 modification 将 是 made。",
																		},
																	},
																},
															},
															"subtitle_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "subtitle settings。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"path": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "URL 的 subtitles 到 add 到 视频。",
																		},
																		"stream_index": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "subtitle track 到 add 到 视频. 如果 both `路径` 和 `StreamIndex` 是 指定，`路径` 将 是 使用. You need 到 指定at least 一个 的 two 参数。",
																		},
																		"font_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "font 类型 有效值：`hei.ttf` `song.ttf` `simkai.ttf` `arial.ttf` (对于 English 仅). 默认为 `hei.ttf`。",
																		},
																		"font_size": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "font 大小 (pixels). 如果 此 是 不 指定， font 大小 在 subtitle 文件 将 是 使用。",
																		},
																		"font_color": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "font color 在 0xRRGGBB 格式 默认值：0xFFFFFF (white)。",
																		},
																		"font_alpha": {
																			Type:        schema.TypeFloat,
																			Optional:    true,
																			Description: "text transparency. 取值范围：0-1. 0: Completely transparent 1: Completely opaque默认值：1。",
																		},
																	},
																},
															},
															"addon_audio_stream": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "信息 的 外部 音频 track 到 add.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "input 类型 有效值：`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
																		},
																		"cos_input_info": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "信息 的 COS 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `COS`。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "COS 存储桶 的 对象 到 process，such 作为 `TopRankVideo-125xxx88`。",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "地域 的 COS 存储桶，such 作为 `ap-chongqing`。",
																					},
																					"object": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "路径 的 对象 到 process，such 作为 `/movie/201907/WildAnimal.mov`。",
																					},
																				},
																			},
																		},
																		"url_input_info": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "URL 的 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"url": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "URL 的 视频。",
																					},
																				},
																			},
																		},
																		"s3_input_info": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "信息 的 AWS S3 对象 processed. 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "AWS S3 存储桶",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "地域 的 AWS S3 存储桶",
																					},
																					"s3_object": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "路径 的 AWS S3 对象。",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "键 ID 必填 到 访问 AWS S3 对象。",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Optional:    true,
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
																Optional:    true,
																Description: "Transcoding extension 字段.注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"add_on_subtitles": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Subtitle files 到 insert.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "inserting 类型 有效值：`subtitle-流`:Insert title track. `close-caption-708`:CEA-708 subtitle encode 到 SEI frame. `close-caption-608`:CEA-608 subtitle encode 到 SEI frame. 注意：此字段可能返回 null，表示无法获取有效值。",
																		},
																		"subtitle": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "Subtitle 文件.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "input 类型 有效值： `COS`:A COS 存储桶 地址 `URL`:A URL `AWS-S3`:An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "信息 的 COS 对象 到 process. 此 参数 是 有效 和 必填 当 类型 是 COS。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "COS 存储桶 的 对象 到 process，such 作为 TopRankVideo-125xxx88。",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "地域 的 COS 存储桶，such 作为 ap-chongqing。",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "路径 的 对象 到 process，such 作为 /movie/201907/WildAnimal.mov。",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "URL 的 对象 到 process. 此 参数 是 有效 和 必填 当 类型 是 URL注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "URL 的 视频。",
																								},
																							},
																						},
																					},
																					"s3_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "信息 的 AWS S3 对象 processed. 此 参数 为必填项 如果 类型 是 AWS-S3.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"s3_bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "S3 存储桶注意：此字段可能返回 null，表示无法获取有效值。",
																								},
																								"s3_region": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "地域 的 AWS S3 存储桶，support: us-east-1 eu-west-3注意：此字段可能返回 null，表示无法获取有效值。",
																								},
																								"s3_object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "路径 的 AWS S3 对象.注意：此字段可能返回 null，表示无法获取有效值。",
																								},
																								"s3_secret_id": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "键 ID 必填 到 访问 AWS S3 对象.注意：此字段可能返回 null，表示无法获取有效值。",
																								},
																								"s3_secret_key": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "键 必填 到 访问 AWS S3 对象.注意：此字段可能返回 null，表示无法获取有效值。",
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
													Optional:    true,
													Description: "列表 up 到 10 镜像 或 text watermarks.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "ID watermarking template。",
															},
															"raw_parameter": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Computed:    true,
																Description: "Custom 水印 参数，其中 是 有效 如果 `Definition` 是 0.此 参数 是 使用 在 highly customized scenarios. We recommend 您 使用 `Definition` 到 指定watermark 参数 preferably.Custom 水印 参数 是 不 可用 对于 screenshot。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Watermark 类型 有效值：镜像: 镜像 水印。",
																		},
																		"coordinate_origin": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Origin position，其中 currently 可以 仅 是: TopLeft: 源站 的 coordinates 是 在 top-left corner 的 视频，和 源站 的 水印 是 在 top-left corner 的 镜像 或 text.默认值：TopLeft。",
																		},
																		"x_pos": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "horizontal position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `XPos` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `XPos` 是 10% 的 视频 宽度; 如果 字符串 结束 在 像素， `XPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `XPos` 是 100 像素.默认值：0 像素。",
																		},
																		"y_pos": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "vertical position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `YPos` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `YPos` 是 10% 的 视频 高度; 如果 字符串 结束 在 像素， `YPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `YPos` 是 100 像素.默认值：0 像素。",
																		},
																		"image_template": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "Image 水印 template. 此 字段 为必填项 当 `类型` 是 `镜像` 和 是 无效 当 `类型` 是 `text`。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"image_content": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Required:    true,
																						Description: "Input 内容 的 水印 镜像. JPEG 和 PNG images 是 支持。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"type": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "input 类型 有效值：`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
																								},
																								"cos_input_info": {
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "信息 的 COS 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `COS`。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"bucket": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "COS 存储桶 的 对象 到 process，such 作为 `TopRankVideo-125xxx88`。",
																											},
																											"region": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "地域 的 COS 存储桶，such 作为 `ap-chongqing`。",
																											},
																											"object": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "路径 的 对象 到 process，such 作为 `/movie/201907/WildAnimal.mov`。",
																											},
																										},
																									},
																								},
																								"url_input_info": {
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "URL 的 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"url": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "URL 的 视频。",
																											},
																										},
																									},
																								},
																								"s3_input_info": {
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "信息 的 AWS S3 对象 processed. 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"s3_bucket": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "AWS S3 存储桶",
																											},
																											"s3_region": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "地域 的 AWS S3 存储桶",
																											},
																											"s3_object": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "路径 的 AWS S3 对象。",
																											},
																											"s3_secret_id": {
																												Type:        schema.TypeString,
																												Optional:    true,
																												Description: "键 ID 必填 到 访问 AWS S3 对象。",
																											},
																											"s3_secret_key": {
																												Type:        schema.TypeString,
																												Optional:    true,
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
																						Optional:    true,
																						Description: "Watermark 宽度. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `宽度` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `宽度` 是 10% 的 视频 宽度; 如果 字符串 结束 在 像素， `宽度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `宽度` 是 100 像素.默认值：10%。",
																					},
																					"height": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Watermark 高度. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `高度` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `高度` 是 10% 的 视频 高度; 如果 字符串 结束 在 像素， `高度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `高度` 是 100 像素.默认值：0 像素，其中 表示 该 `高度` 将 是 proportionally scaled according 到 aspect ratio 的 original 水印 镜像。",
																					},
																					"repeat_type": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Repeat 类型 animated 水印. 有效值：`once`: 无 longer appears after 水印 playback 结束. `repeat_last_frame`: stays 在 last frame after 水印 playback 结束. `repeat` (默认值): repeats playback until 视频 结束。",
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
																Description: "Text 内容 的 up 到 100 字符. 此 字段 为必填项 仅 当 水印 类型 是 text.Text 水印 是 不 可用 对于 screenshot。",
															},
															"svg_content": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "SVG 内容 的 up 到 2,000,000 字符. 此 字段 为必填项 仅 当 水印 类型 是 `SVG`.SVG 水印 是 不 可用 对于 screenshot。",
															},
															"start_time_offset": {
																Type:        schema.TypeFloat,
																Optional:    true,
																Description: "开始时间 偏移量 的 水印 （秒）。 如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame. 如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame; 如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 appear 在 second n after first 视频 frame; 如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 appear 在 second n before last 视频 frame。",
															},
															"end_time_offset": {
																Type:        schema.TypeFloat,
																Optional:    true,
																Description: "结束时间 偏移量 的 水印 （秒）。 如果此参数为空 或 0 是 entered， 水印 将 exist till last 视频 frame; 如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 exist till second n; 如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 exist till second n before last 视频 frame。",
															},
														},
													},
												},
												"mosaic_set": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "列表 blurs. Up 到 10 ones 可以 是 支持。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"coordinate_origin": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Origin position，其中 currently 可以 仅 是: TopLeft: 源站 的 coordinates 是 在 top-left corner 的 视频，和 源站 的 blur 是 在 top-left corner 的 镜像 或 text.默认值：TopLeft。",
															},
															"x_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "horizontal position 的 源站 的 blur relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `XPos` 的 blur 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `XPos` 是 10% 的 视频 宽度; 如果 字符串 结束 在 像素， `XPos` 的 blur 将 是 指定 像素; 对于 示例，`100px` 表示 该 `XPos` 是 100 像素.默认值：0 像素。",
															},
															"y_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Vertical position 的 源站 的 blur relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `YPos` 的 blur 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `YPos` 是 10% 的 视频 高度; 如果 字符串 结束 在 像素， `YPos` 的 blur 将 是 指定 像素; 对于 示例，`100px` 表示 该 `YPos` 是 100 像素.默认值：0 像素。",
															},
															"width": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Blur 宽度. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `宽度` 的 blur 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `宽度` 是 10% 的 视频 宽度; 如果 字符串 结束 在 像素， `宽度` 的 blur 将 是 在 像素; 对于 示例，`100px` 表示 该 `宽度` 是 100 像素.默认值：10%。",
															},
															"height": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Blur 高度. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `高度` 的 blur 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `高度` 是 10% 的 视频 高度; 如果 字符串 结束 在 像素， `高度` 的 blur 将 是 在 像素; 对于 示例，`100px` 表示 该 `高度` 是 100 像素.默认值：10%。",
															},
															"start_time_offset": {
																Type:        schema.TypeFloat,
																Optional:    true,
																Description: "开始时间 偏移量 的 blur （秒）。 如果此参数为空 或 0 是 entered， blur 将 appear upon first 视频 frame. 如果此参数为空 或 0 是 entered， blur 将 appear upon first 视频 frame; 如果 此 值 是 greater 比 0 (e.g.，n)， blur 将 appear 在 second n after first 视频 frame; 如果 此 值 是 smaller 比 0 (e.g.，-n)， blur 将 appear 在 second n before last 视频 frame。",
															},
															"end_time_offset": {
																Type:        schema.TypeFloat,
																Optional:    true,
																Description: "结束时间 偏移量 的 blur （秒）。 如果此参数为空 或 0 是 entered， blur 将 exist till last 视频 frame; 如果 此 值 是 greater 比 0 (e.g.，n)， blur 将 exist till second n; 如果 此 值 是 smaller 比 0 (e.g.，-n)， blur 将 exist till second n before last 视频 frame。",
															},
														},
													},
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "开始时间 偏移量 的 transcoded 视频，（秒）。 如果此参数为空 或 集合 到 0， transcoded 视频 将 start 在 same 时间 作为 original 视频. 如果 此 参数 是 集合 到 positive 数量 (n 对于 示例)， transcoded 视频 将 start 在 nth second 的 original 视频. 如果 此 参数 是 集合 到 negative 数量 (-n 对于 示例)， transcoded 视频 将 start 在 nth second before end 的 original 视频。",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "结束时间 偏移量 的 transcoded 视频，（秒）。 如果此参数为空 或 集合 到 0， transcoded 视频 将 end 在 same 时间 作为 original 视频. 如果 此 参数 是 集合 到 positive 数量 (n 对于 示例)， transcoded 视频 将 end 在 nth second 的 original 视频. 如果 此 参数 是 集合 到 negative 数量 (-n 对于 示例)， transcoded 视频 将 end 在 nth second before end 的 original 视频。",
												},
												"output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Target 存储桶 的 输出文件 如果此参数为空， `OutputStorage` 值 的 upper 文件夹 将 是 inherited.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "存储 类型 对于 media processing 输出文件 有效值：`COS`: Tencent Cloud COS `AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
															},
															"cos_output_storage": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "location 到 save output 对象 在 COS. 此 参数 是 有效 和 必填 当 `类型` 是 COS.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bucket": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "存储桶 到 其中 输出文件 的 media processing 是 saved，such 作为 `TopRankVideo-125xxx88`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																		},
																		"region": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "地域 的 输出存储桶，such 作为 `ap-chongqing`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																		},
																	},
																},
															},
															"s3_output_storage": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "AWS S3 存储桶 到 save 输出文件 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"s3_bucket": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "AWS S3 存储桶",
																		},
																		"s3_region": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "地域 的 AWS S3 存储桶",
																		},
																		"s3_secret_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "键 ID 必填 到 upload files 到 AWS S3 对象。",
																		},
																		"s3_secret_key": {
																			Type:        schema.TypeString,
																			Optional:    true,
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
													Optional:    true,
													Description: "路径 到 primary 输出文件，其中 可以 是 relative 路径 或 absolute 路径 如果此参数为空， following relative 路径 将 是 使用 通过 默认值：`{inputName}_transcode_{definition}.{格式}`。",
												},
												"segment_object_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "路径 到 输出文件 part ( 路径 到 ts during transcoding 到 HLS)，其中 可以 仅 是 relative 路径 如果此参数为空， following relative 路径 将 是 使用 通过 默认值：`{inputName}_transcode_{definition}_{数量}.{格式}`。",
												},
												"object_number_format": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Rule 的 `{数量}` variable 在 输出路径 after transcoding.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"initial_value": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Start 值 的 `{数量}` variable. 默认值：0。",
															},
															"increment": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Increment 的 `{数量}` variable. 默认值：1。",
															},
															"min_length": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "最小长度the `{数量}` variable. A placeholder 将 是 使用 如果 variable 长度 是 below 最小 requirement. 默认值：1。",
															},
															"place_holder": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Placeholder 使用 当 `{数量}` variable 长度 是 below 最小 requirement. 默认值：0。",
															},
														},
													},
												},
												"head_tail_parameter": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Opening 和 closing credits parametersNote: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 是 found。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"head_set": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Opening credits 列表。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "input 类型 有效值：`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
																		},
																		"cos_input_info": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "信息 的 COS 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `COS`。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "COS 存储桶 的 对象 到 process，such 作为 `TopRankVideo-125xxx88`。",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "地域 的 COS 存储桶，such 作为 `ap-chongqing`。",
																					},
																					"object": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "路径 的 对象 到 process，such 作为 `/movie/201907/WildAnimal.mov`。",
																					},
																				},
																			},
																		},
																		"url_input_info": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "URL 的 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"url": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "URL 的 视频。",
																					},
																				},
																			},
																		},
																		"s3_input_info": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "信息 的 AWS S3 对象 processed. 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "AWS S3 存储桶",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "地域 的 AWS S3 存储桶",
																					},
																					"s3_object": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "路径 的 AWS S3 对象。",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "键 ID 必填 到 访问 AWS S3 对象。",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Optional:    true,
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
																Optional:    true,
																Description: "Closing credits 列表。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "input 类型 有效值：`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
																		},
																		"cos_input_info": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "信息 的 COS 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `COS`。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "COS 存储桶 的 对象 到 process，such 作为 `TopRankVideo-125xxx88`。",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "地域 的 COS 存储桶，such 作为 `ap-chongqing`。",
																					},
																					"object": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "路径 的 对象 到 process，such 作为 `/movie/201907/WildAnimal.mov`。",
																					},
																				},
																			},
																		},
																		"url_input_info": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "URL 的 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"url": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "URL 的 视频。",
																					},
																				},
																			},
																		},
																		"s3_input_info": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "信息 的 AWS S3 对象 processed. 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "AWS S3 存储桶",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "地域 的 AWS S3 存储桶",
																					},
																					"s3_object": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "路径 的 AWS S3 对象。",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "键 ID 必填 到 访问 AWS S3 对象。",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Optional:    true,
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
										MaxItems:    1,
										Optional:    true,
										Description: "An animated screenshot generation 任务。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
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
												"output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Target 存储桶 的 generated animated 镜像 文件. 如果此参数为空， `OutputStorage` 值 的 upper 文件夹 将 是 inherited.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "存储 类型 对于 media processing 输出文件 有效值：`COS`: Tencent Cloud COS `AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
															},
															"cos_output_storage": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "location 到 save output 对象 在 COS. 此 参数 是 有效 和 必填 当 `类型` 是 COS.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bucket": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "存储桶 到 其中 输出文件 的 media processing 是 saved，such 作为 `TopRankVideo-125xxx88`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																		},
																		"region": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "地域 的 输出存储桶，such 作为 `ap-chongqing`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																		},
																	},
																},
															},
															"s3_output_storage": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "AWS S3 存储桶 到 save 输出文件 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"s3_bucket": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "AWS S3 存储桶",
																		},
																		"s3_region": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "地域 的 AWS S3 存储桶",
																		},
																		"s3_secret_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "键 ID 必填 到 upload files 到 AWS S3 对象。",
																		},
																		"s3_secret_key": {
																			Type:        schema.TypeString,
																			Optional:    true,
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
													Optional:    true,
													Description: "输出路径 到 generated animated 镜像 文件，其中 可以 是 relative 路径 或 absolute 路径 如果此参数为空， following relative 路径 将 是 使用 通过 默认值：`{inputName}_animatedGraphic_{definition}.{格式}`。",
												},
											},
										},
									},
									"snapshot_by_time_offset_task": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "A 时间 point screencapturing 任务。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "ID 时间 point screencapturing template。",
												},
												"ext_time_offset_set": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Description: "列表 screenshot 时间 points 在 格式 的 `s` 或 `%`: 如果 字符串 结束 在 `s`，它 表示 该 时间 point 是 在 秒; 对于 示例，`3.5s` 表示 该 时间 point 是 3.5th second; 如果 字符串 结束 在 `%`，它 表示 该 时间 point 是 指定 percentage 的 视频 时长; 对于 示例，`10%` 表示 该 时间 point 是 10% 的 视频 时长。",
												},
												"watermark_set": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "列表 up 到 10 镜像 或 text watermarks.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "ID watermarking template。",
															},
															"raw_parameter": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Computed:    true,
																Description: "Custom 水印 参数，其中 是 有效 如果 `Definition` 是 0.此 参数 是 使用 在 highly customized scenarios. We recommend 您 使用 `Definition` 到 指定watermark 参数 preferably.Custom 水印 参数 是 不 可用 对于 screenshot。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Watermark 类型 有效值：镜像: 镜像 水印。",
																		},
																		"coordinate_origin": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Origin position，其中 currently 可以 仅 是: TopLeft: 源站 的 coordinates 是 在 top-left corner 的 视频，和 源站 的 水印 是 在 top-left corner 的 镜像 或 text.默认值：TopLeft。",
																		},
																		"x_pos": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "horizontal position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `XPos` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `XPos` 是 10% 的 视频 宽度; 如果 字符串 结束 在 像素， `XPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `XPos` 是 100 像素.默认值：0 像素。",
																		},
																		"y_pos": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "vertical position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `YPos` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `YPos` 是 10% 的 视频 高度; 如果 字符串 结束 在 像素， `YPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `YPos` 是 100 像素.默认值：0 像素。",
																		},
																		"image_template": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "Image 水印 template. 此 字段 为必填项 当 `类型` 是 `镜像` 和 是 无效 当 `类型` 是 `text`。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"image_content": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Required:    true,
																						Description: "Input 内容 的 水印 镜像. JPEG 和 PNG images 是 支持。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"type": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "input 类型 有效值：`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
																								},
																								"cos_input_info": {
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "信息 的 COS 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `COS`。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"bucket": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "COS 存储桶 的 对象 到 process，such 作为 `TopRankVideo-125xxx88`。",
																											},
																											"region": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "地域 的 COS 存储桶，such 作为 `ap-chongqing`。",
																											},
																											"object": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "路径 的 对象 到 process，such 作为 `/movie/201907/WildAnimal.mov`。",
																											},
																										},
																									},
																								},
																								"url_input_info": {
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "URL 的 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"url": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "URL 的 视频。",
																											},
																										},
																									},
																								},
																								"s3_input_info": {
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "信息 的 AWS S3 对象 processed. 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"s3_bucket": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "AWS S3 存储桶",
																											},
																											"s3_region": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "地域 的 AWS S3 存储桶",
																											},
																											"s3_object": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "路径 的 AWS S3 对象。",
																											},
																											"s3_secret_id": {
																												Type:        schema.TypeString,
																												Optional:    true,
																												Description: "键 ID 必填 到 访问 AWS S3 对象。",
																											},
																											"s3_secret_key": {
																												Type:        schema.TypeString,
																												Optional:    true,
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
																						Optional:    true,
																						Description: "Watermark 宽度. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `宽度` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `宽度` 是 10% 的 视频 宽度; 如果 字符串 结束 在 像素， `宽度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `宽度` 是 100 像素.默认值：10%。",
																					},
																					"height": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Watermark 高度. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `高度` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `高度` 是 10% 的 视频 高度; 如果 字符串 结束 在 像素， `高度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `高度` 是 100 像素.默认值：0 像素，其中 表示 该 `高度` 将 是 proportionally scaled according 到 aspect ratio 的 original 水印 镜像。",
																					},
																					"repeat_type": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Repeat 类型 animated 水印. 有效值：`once`: 无 longer appears after 水印 playback 结束. `repeat_last_frame`: stays 在 last frame after 水印 playback 结束. `repeat` (默认值): repeats playback until 视频 结束。",
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
																Description: "Text 内容 的 up 到 100 字符. 此 字段 为必填项 仅 当 水印 类型 是 text.Text 水印 是 不 可用 对于 screenshot。",
															},
															"svg_content": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "SVG 内容 的 up 到 2,000,000 字符. 此 字段 为必填项 仅 当 水印 类型 是 `SVG`.SVG 水印 是 不 可用 对于 screenshot。",
															},
															"start_time_offset": {
																Type:        schema.TypeFloat,
																Optional:    true,
																Description: "开始时间 偏移量 的 水印 （秒）。 如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame. 如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame; 如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 appear 在 second n after first 视频 frame; 如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 appear 在 second n before last 视频 frame。",
															},
															"end_time_offset": {
																Type:        schema.TypeFloat,
																Optional:    true,
																Description: "结束时间 偏移量 的 水印 （秒）。 如果此参数为空 或 0 是 entered， 水印 将 exist till last 视频 frame; 如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 exist till second n; 如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 exist till second n before last 视频 frame。",
															},
														},
													},
												},
												"output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Target 存储桶 的 generated 时间 point screenshot 文件. 如果此参数为空， `OutputStorage` 值 的 upper 文件夹 将 是 inherited.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "存储 类型 对于 media processing 输出文件 有效值：`COS`: Tencent Cloud COS `AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
															},
															"cos_output_storage": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "location 到 save output 对象 在 COS. 此 参数 是 有效 和 必填 当 `类型` 是 COS.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bucket": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "存储桶 到 其中 输出文件 的 media processing 是 saved，such 作为 `TopRankVideo-125xxx88`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																		},
																		"region": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "地域 的 输出存储桶，such 作为 `ap-chongqing`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																		},
																	},
																},
															},
															"s3_output_storage": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "AWS S3 存储桶 到 save 输出文件 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"s3_bucket": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "AWS S3 存储桶",
																		},
																		"s3_region": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "地域 的 AWS S3 存储桶",
																		},
																		"s3_secret_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "键 ID 必填 到 upload files 到 AWS S3 对象。",
																		},
																		"s3_secret_key": {
																			Type:        schema.TypeString,
																			Optional:    true,
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
													Optional:    true,
													Description: "输出路径 到 generated 时间 point screenshot，其中 可以 是 relative 路径 或 absolute 路径 如果此参数为空， following relative 路径 将 是 使用 通过 默认值：`{inputName}_snapshotByTimeOffset_{definition}_{数量}.{格式}`。",
												},
												"object_number_format": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Rule 的 `{数量}` variable 在 时间 point screenshot 输出路径注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"initial_value": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Start 值 的 `{数量}` variable. 默认值：0。",
															},
															"increment": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Increment 的 `{数量}` variable. 默认值：1。",
															},
															"min_length": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "最小长度the `{数量}` variable. A placeholder 将 是 使用 如果 variable 长度 是 below 最小 requirement. 默认值：1。",
															},
															"place_holder": {
																Type:        schema.TypeString,
																Optional:    true,
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
										MaxItems:    1,
										Optional:    true,
										Description: "A sampled screencapturing 任务。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Sampled screencapturing 模板 ID",
												},
												"watermark_set": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "列表 up 到 10 镜像 或 text watermarks.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "ID watermarking template。",
															},
															"raw_parameter": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Computed:    true,
																Description: "Custom 水印 参数，其中 是 有效 如果 `Definition` 是 0.此 参数 是 使用 在 highly customized scenarios. We recommend 您 使用 `Definition` 到 指定watermark 参数 preferably.Custom 水印 参数 是 不 可用 对于 screenshot。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Watermark 类型 有效值：镜像: 镜像 水印。",
																		},
																		"coordinate_origin": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Origin position，其中 currently 可以 仅 是: TopLeft: 源站 的 coordinates 是 在 top-left corner 的 视频，和 源站 的 水印 是 在 top-left corner 的 镜像 或 text.默认值：TopLeft。",
																		},
																		"x_pos": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "horizontal position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `XPos` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `XPos` 是 10% 的 视频 宽度; 如果 字符串 结束 在 像素， `XPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `XPos` 是 100 像素.默认值：0 像素。",
																		},
																		"y_pos": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "vertical position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `YPos` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `YPos` 是 10% 的 视频 高度; 如果 字符串 结束 在 像素， `YPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `YPos` 是 100 像素.默认值：0 像素。",
																		},
																		"image_template": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "Image 水印 template. 此 字段 为必填项 当 `类型` 是 `镜像` 和 是 无效 当 `类型` 是 `text`。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"image_content": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Required:    true,
																						Description: "Input 内容 的 水印 镜像. JPEG 和 PNG images 是 支持。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"type": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "input 类型 有效值：`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
																								},
																								"cos_input_info": {
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "信息 的 COS 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `COS`。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"bucket": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "COS 存储桶 的 对象 到 process，such 作为 `TopRankVideo-125xxx88`。",
																											},
																											"region": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "地域 的 COS 存储桶，such 作为 `ap-chongqing`。",
																											},
																											"object": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "路径 的 对象 到 process，such 作为 `/movie/201907/WildAnimal.mov`。",
																											},
																										},
																									},
																								},
																								"url_input_info": {
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "URL 的 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"url": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "URL 的 视频。",
																											},
																										},
																									},
																								},
																								"s3_input_info": {
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "信息 的 AWS S3 对象 processed. 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"s3_bucket": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "AWS S3 存储桶",
																											},
																											"s3_region": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "地域 的 AWS S3 存储桶",
																											},
																											"s3_object": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "路径 的 AWS S3 对象。",
																											},
																											"s3_secret_id": {
																												Type:        schema.TypeString,
																												Optional:    true,
																												Description: "键 ID 必填 到 访问 AWS S3 对象。",
																											},
																											"s3_secret_key": {
																												Type:        schema.TypeString,
																												Optional:    true,
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
																						Optional:    true,
																						Description: "Watermark 宽度. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `宽度` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `宽度` 是 10% 的 视频 宽度; 如果 字符串 结束 在 像素， `宽度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `宽度` 是 100 像素.默认值：10%。",
																					},
																					"height": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Watermark 高度. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `高度` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `高度` 是 10% 的 视频 高度; 如果 字符串 结束 在 像素， `高度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `高度` 是 100 像素.默认值：0 像素，其中 表示 该 `高度` 将 是 proportionally scaled according 到 aspect ratio 的 original 水印 镜像。",
																					},
																					"repeat_type": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Repeat 类型 animated 水印. 有效值：`once`: 无 longer appears after 水印 playback 结束. `repeat_last_frame`: stays 在 last frame after 水印 playback 结束. `repeat` (默认值): repeats playback until 视频 结束。",
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
																Description: "Text 内容 的 up 到 100 字符. 此 字段 为必填项 仅 当 水印 类型 是 text.Text 水印 是 不 可用 对于 screenshot。",
															},
															"svg_content": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "SVG 内容 的 up 到 2,000,000 字符. 此 字段 为必填项 仅 当 水印 类型 是 `SVG`.SVG 水印 是 不 可用 对于 screenshot。",
															},
															"start_time_offset": {
																Type:        schema.TypeFloat,
																Optional:    true,
																Description: "开始时间 偏移量 的 水印 （秒）。 如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame. 如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame; 如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 appear 在 second n after first 视频 frame; 如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 appear 在 second n before last 视频 frame。",
															},
															"end_time_offset": {
																Type:        schema.TypeFloat,
																Optional:    true,
																Description: "结束时间 偏移量 的 水印 （秒）。 如果此参数为空 或 0 是 entered， 水印 将 exist till last 视频 frame; 如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 exist till second n; 如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 exist till second n before last 视频 frame。",
															},
														},
													},
												},
												"output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Target 存储桶 的 sampled screenshot. 如果此参数为空， `OutputStorage` 值 的 upper 文件夹 将 是 inherited.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "存储 类型 对于 media processing 输出文件 有效值：`COS`: Tencent Cloud COS `AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
															},
															"cos_output_storage": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "location 到 save output 对象 在 COS. 此 参数 是 有效 和 必填 当 `类型` 是 COS.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bucket": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "存储桶 到 其中 输出文件 的 media processing 是 saved，such 作为 `TopRankVideo-125xxx88`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																		},
																		"region": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "地域 的 输出存储桶，such 作为 `ap-chongqing`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																		},
																	},
																},
															},
															"s3_output_storage": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "AWS S3 存储桶 到 save 输出文件 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"s3_bucket": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "AWS S3 存储桶",
																		},
																		"s3_region": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "地域 的 AWS S3 存储桶",
																		},
																		"s3_secret_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "键 ID 必填 到 upload files 到 AWS S3 对象。",
																		},
																		"s3_secret_key": {
																			Type:        schema.TypeString,
																			Optional:    true,
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
													Optional:    true,
													Description: "输出路径 到 generated sampled screenshot，其中 可以 是 relative 路径 或 absolute 路径 如果此参数为空， following relative 路径 将 是 使用 通过 默认值：`{inputName}_sampleSnapshot_{definition}_{数量}.{格式}`。",
												},
												"object_number_format": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Rule 的 `{数量}` variable 在 sampled screenshot 输出路径注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"initial_value": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Start 值 的 `{数量}` variable. 默认值：0。",
															},
															"increment": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Increment 的 `{数量}` variable. 默认值：1。",
															},
															"min_length": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "最小长度the `{数量}` variable. A placeholder 将 是 使用 如果 variable 长度 是 below 最小 requirement. 默认值：1。",
															},
															"place_holder": {
																Type:        schema.TypeString,
																Optional:    true,
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
										MaxItems:    1,
										Optional:    true,
										Description: "An 镜像 sprite generation 任务。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "ID 镜像 sprite generating template。",
												},
												"output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Target 存储桶 的 generated 镜像 sprite. 如果此参数为空， `OutputStorage` 值 的 upper 文件夹 将 是 inherited.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "存储 类型 对于 media processing 输出文件 有效值：`COS`: Tencent Cloud COS `AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
															},
															"cos_output_storage": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "location 到 save output 对象 在 COS. 此 参数 是 有效 和 必填 当 `类型` 是 COS.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bucket": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "存储桶 到 其中 输出文件 的 media processing 是 saved，such 作为 `TopRankVideo-125xxx88`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																		},
																		"region": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "地域 的 输出存储桶，such 作为 `ap-chongqing`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																		},
																	},
																},
															},
															"s3_output_storage": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "AWS S3 存储桶 到 save 输出文件 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"s3_bucket": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "AWS S3 存储桶",
																		},
																		"s3_region": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "地域 的 AWS S3 存储桶",
																		},
																		"s3_secret_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "键 ID 必填 到 upload files 到 AWS S3 对象。",
																		},
																		"s3_secret_key": {
																			Type:        schema.TypeString,
																			Optional:    true,
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
													Optional:    true,
													Description: "输出路径 到 generated 镜像 sprite 文件，其中 可以 是 relative 路径 或 absolute 路径 如果此参数为空， following relative 路径 将 是 使用 通过 默认值：`{inputName}_imageSprite_{definition}_{数量}.{格式}`。",
												},
												"web_vtt_object_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "输出路径 到 WebVTT 文件 after 镜像 sprite 是 generated，其中 可以 仅 是 relative 路径 如果此参数为空， following relative 路径 将 是 使用 通过 默认值：`{inputName}_imageSprite_{definition}.{格式}`。",
												},
												"object_number_format": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Rule 的 `{数量}` variable 在 镜像 sprite 输出路径注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"initial_value": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Start 值 的 `{数量}` variable. 默认值：0。",
															},
															"increment": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Increment 的 `{数量}` variable. 默认值：1。",
															},
															"min_length": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "最小长度the `{数量}` variable. A placeholder 将 是 使用 如果 variable 长度 是 below 最小 requirement. 默认值：1。",
															},
															"place_holder": {
																Type:        schema.TypeString,
																Optional:    true,
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
										MaxItems:    1,
										Optional:    true,
										Description: "An adaptive bitrate streaming 任务。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Adaptive bitrate streaming 模板 ID",
												},
												"watermark_set": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "列表 up 到 10 镜像 或 text watermarks。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "ID watermarking template。",
															},
															"raw_parameter": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Computed:    true,
																Description: "Custom 水印 参数，其中 是 有效 如果 `Definition` 是 0.此 参数 是 使用 在 highly customized scenarios. We recommend 您 使用 `Definition` 到 指定watermark 参数 preferably.Custom 水印 参数 是 不 可用 对于 screenshot。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Watermark 类型 有效值：镜像: 镜像 水印。",
																		},
																		"coordinate_origin": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Origin position，其中 currently 可以 仅 是: TopLeft: 源站 的 coordinates 是 在 top-left corner 的 视频，和 源站 的 水印 是 在 top-left corner 的 镜像 或 text.默认值：TopLeft。",
																		},
																		"x_pos": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "horizontal position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `XPos` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `XPos` 是 10% 的 视频 宽度; 如果 字符串 结束 在 像素， `XPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `XPos` 是 100 像素.默认值：0 像素。",
																		},
																		"y_pos": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "vertical position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `YPos` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `YPos` 是 10% 的 视频 高度; 如果 字符串 结束 在 像素， `YPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `YPos` 是 100 像素.默认值：0 像素。",
																		},
																		"image_template": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "Image 水印 template. 此 字段 为必填项 当 `类型` 是 `镜像` 和 是 无效 当 `类型` 是 `text`。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"image_content": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Required:    true,
																						Description: "Input 内容 的 水印 镜像. JPEG 和 PNG images 是 支持。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"type": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "input 类型 有效值：`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
																								},
																								"cos_input_info": {
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "信息 的 COS 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `COS`。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"bucket": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "COS 存储桶 的 对象 到 process，such 作为 `TopRankVideo-125xxx88`。",
																											},
																											"region": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "地域 的 COS 存储桶，such 作为 `ap-chongqing`。",
																											},
																											"object": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "路径 的 对象 到 process，such 作为 `/movie/201907/WildAnimal.mov`。",
																											},
																										},
																									},
																								},
																								"url_input_info": {
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "URL 的 对象 到 process. 此 参数 是 有效 和 必填 当 `类型` 是 `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"url": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "URL 的 视频。",
																											},
																										},
																									},
																								},
																								"s3_input_info": {
																									Type:        schema.TypeList,
																									MaxItems:    1,
																									Optional:    true,
																									Description: "信息 的 AWS S3 对象 processed. 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"s3_bucket": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "AWS S3 存储桶",
																											},
																											"s3_region": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "地域 的 AWS S3 存储桶",
																											},
																											"s3_object": {
																												Type:        schema.TypeString,
																												Required:    true,
																												Description: "路径 的 AWS S3 对象。",
																											},
																											"s3_secret_id": {
																												Type:        schema.TypeString,
																												Optional:    true,
																												Description: "键 ID 必填 到 访问 AWS S3 对象。",
																											},
																											"s3_secret_key": {
																												Type:        schema.TypeString,
																												Optional:    true,
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
																						Optional:    true,
																						Description: "Watermark 宽度. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `宽度` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `宽度` 是 10% 的 视频 宽度; 如果 字符串 结束 在 像素， `宽度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `宽度` 是 100 像素.默认值：10%。",
																					},
																					"height": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Watermark 高度. % 和 像素 formats 是 支持: 如果 字符串 结束 在 %， `高度` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `高度` 是 10% 的 视频 高度; 如果 字符串 结束 在 像素， `高度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `高度` 是 100 像素.默认值：0 像素，其中 表示 该 `高度` 将 是 proportionally scaled according 到 aspect ratio 的 original 水印 镜像。",
																					},
																					"repeat_type": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Repeat 类型 animated 水印. 有效值：`once`: 无 longer appears after 水印 playback 结束. `repeat_last_frame`: stays 在 last frame after 水印 playback 结束. `repeat` (默认值): repeats playback until 视频 结束。",
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
																Description: "Text 内容 的 up 到 100 字符. 此 字段 为必填项 仅 当 水印 类型 是 text.Text 水印 是 不 可用 对于 screenshot。",
															},
															"svg_content": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "SVG 内容 的 up 到 2,000,000 字符. 此 字段 为必填项 仅 当 水印 类型 是 `SVG`.SVG 水印 是 不 可用 对于 screenshot。",
															},
															"start_time_offset": {
																Type:        schema.TypeFloat,
																Optional:    true,
																Description: "开始时间 偏移量 的 水印 （秒）。 如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame. 如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame; 如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 appear 在 second n after first 视频 frame; 如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 appear 在 second n before last 视频 frame。",
															},
															"end_time_offset": {
																Type:        schema.TypeFloat,
																Optional:    true,
																Description: "结束时间 偏移量 的 水印 （秒）。 如果此参数为空 或 0 是 entered， 水印 将 exist till last 视频 frame; 如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 exist till second n; 如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 exist till second n before last 视频 frame。",
															},
														},
													},
												},
												"output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Target 存储桶 的 输出文件 after being transcoded 到 adaptive bitrate streaming. 如果此参数为空， `OutputStorage` 值 的 upper 文件夹 将 是 inherited.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "存储 类型 对于 media processing 输出文件 有效值：`COS`: Tencent Cloud COS `AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
															},
															"cos_output_storage": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "location 到 save output 对象 在 COS. 此 参数 是 有效 和 必填 当 `类型` 是 COS.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bucket": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "存储桶 到 其中 输出文件 的 media processing 是 saved，such 作为 `TopRankVideo-125xxx88`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																		},
																		"region": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "地域 的 输出存储桶，such 作为 `ap-chongqing`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
																		},
																	},
																},
															},
															"s3_output_storage": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "AWS S3 存储桶 到 save 输出文件 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"s3_bucket": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "AWS S3 存储桶",
																		},
																		"s3_region": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "地域 的 AWS S3 存储桶",
																		},
																		"s3_secret_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "键 ID 必填 到 upload files 到 AWS S3 对象。",
																		},
																		"s3_secret_key": {
																			Type:        schema.TypeString,
																			Optional:    true,
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
													Optional:    true,
													Description: "relative 或 absolute 输出路径 的 manifest 文件 after being transcoded 到 adaptive bitrate streaming. 如果此参数为空， relative 路径 在 following 格式 将 是 使用 通过 默认值：`{inputName}_adaptiveDynamicStreaming_{definition}.{格式}`。",
												},
												"sub_stream_object_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "relative 输出路径 的 substream 文件 after being transcoded 到 adaptive bitrate streaming. 如果此参数为空， relative 路径 在 following 格式 将 是 使用 通过 默认值：`{inputName}_adaptiveDynamicStreaming_{definition}_{subStreamNumber}.{格式}`。",
												},
												"segment_object_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "relative 输出路径 的 segment 文件 after being transcoded 到 adaptive bitrate streaming (在 HLS 格式 仅). 如果此参数为空， relative 路径 在 following 格式 将 是 使用 通过 默认值：`{inputName}_adaptiveDynamicStreaming_{definition}_{subStreamNumber}_{segmentNumber}.{格式}`。",
												},
												"add_on_subtitles": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Subtitle files 到 insert.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "inserting 类型 有效值：subtitle-流:Insert title track close-caption-708:CEA-708 subtitle encode 到 SEI frame close-caption-608:CEA-608 subtitle encode 到 SEI frame注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"subtitle": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Subtitle 文件.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "input 类型 有效值： COS:A COS 存储桶 地址 URL:A URL AWS-S3:An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
																		},
																		"cos_input_info": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "信息 的 COS 对象 到 process. 此 参数 是 有效 和 必填 当 类型 是 COS。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "COS 存储桶 的 对象 到 process，such 作为 TopRankVideo-125xxx88。",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "地域 的 COS 存储桶，such 作为 ap-chongqing。",
																					},
																					"object": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "路径 的 对象 到 process，such 作为 /movie/201907/WildAnimal.mov。",
																					},
																				},
																			},
																		},
																		"url_input_info": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "URL 的 对象 到 process. 此 参数 是 有效 和 必填 当 类型 是 URL注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"url": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "URL 的 视频。",
																					},
																				},
																			},
																		},
																		"s3_input_info": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "信息 的 AWS S3 对象 processed. 此 参数 为必填项 如果 类型 是 AWS-S3. 注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "S3 存储桶注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "地域 的 AWS S3 存储桶，support: us-east-1 eu-west-3注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																					"s3_object": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "路径 的 AWS S3 对象.注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "键 ID 必填 到 访问 AWS S3 对象.注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "键 必填 到 访问 AWS S3 对象.注意：此字段可能返回 null，表示无法获取有效值。",
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
										MaxItems:    1,
										Optional:    true,
										Description: "A 内容 moderation 任务。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Video 内容 audit 模板 ID",
												},
											},
										},
									},
									"ai_analysis_task": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "A 内容 analysis 任务。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Video 内容 analysis 模板 ID",
												},
												"extended_parameter": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "An extended 参数，whose 值 是 stringfied JSON.注意: 此 参数 是 对于 customers 使用 special requirements. It needs 到 是 customized offline.注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"ai_recognition_task": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "A 内容 recognition 任务。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"definition": {
													Type:        schema.TypeInt,
													Required:    true,
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
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "存储桶 到 save 输出文件 如果 您 do 不 指定this 参数， 存储桶 在 `Trigger` 将 是 使用。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "存储 类型 对于 media processing 输出文件 有效值：`COS`: Tencent Cloud COS `AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
						},
						"cos_output_storage": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "location 到 save output 对象 在 COS. 此 参数 是 有效 和 必填 当 `类型` 是 COS.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "存储桶 到 其中 输出文件 的 media processing 是 saved，such 作为 `TopRankVideo-125xxx88`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
									},
									"region": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "地域 的 输出存储桶，such 作为 `ap-chongqing`. 如果此参数为空， 值 的 upper layer 将 是 inherited。",
									},
								},
							},
						},
						"s3_output_storage": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "AWS S3 存储桶 到 save 输出文件 此 参数 为必填项 如果 `类型` 是 `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"s3_bucket": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "AWS S3 存储桶",
									},
									"s3_region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "地域 的 AWS S3 存储桶",
									},
									"s3_secret_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "键 ID 必填 到 upload files 到 AWS S3 对象。",
									},
									"s3_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "键 必填 到 upload files 到 AWS S3 对象。",
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
				Description: "directory 到 save media processing 输出文件，其中 必须 start 和 end 使用 `/`，such 作为 `/movie/201907/`.如果 您 do 不 指定this， 文件 将 是 saved 到 触发器 directory。",
			},

			"task_notify_config": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "通知 配置. 如果 您 do 不 指定this 参数，notifications 将 不 是 sent。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cmq_model": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CMQ 或 TDMQ-CMQ model. 有效值：Queue，Topic。",
						},
						"cmq_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CMQ 或 TDMQ-CMQ 地域，such 作为 `sh` (Shanghai) 或 `bj` (Beijing)。",
						},
						"topic_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CMQ 或 TDMQ-CMQ 主题 到 receive notifications. 此 参数 是 有效 当 `CmqModel` 是 `Topic`。",
						},
						"queue_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CMQ 或 TDMQ-CMQ queue 到 receive notifications. 此 参数 是 有效 当 `CmqModel` 是 `Queue`。",
						},
						"notify_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Workflow 通知 方法. 有效值：Finish，Change. 如果此参数为空，`Finish` 将 是 使用。",
						},
						"notify_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "通知 类型 有效值：`CMQ`: 此 值 是 无 longer 使用. Please 使用 `TDMQ-CMQ` instead. `TDMQ-CMQ`: 消息 queue `URL`: 如果 `NotifyType` 是 集合 到 `URL`，HTTP callbacks 是 sent 到 URL 指定 通过 `NotifyUrl`. HTTP 和 JSON 是 用于the callbacks. packet 包含response 参数 的 `ParseNotification` API. `SCF`: 此 通知 类型 是 不 recommended. You need 到 configure 它 在 SCF console. `AWS-SQS`: AWS queue. 此 类型 是 仅 支持 对于 AWS tasks，和 queue 必须 是 在 same 地域 作为 AWS 存储桶Note: 如果 您 do 不 pass 此 参数 或 pass 在 空 字符串，`CMQ` 将 是 使用. To 使用 different 通知 类型，指定this 参数 accordingly。",
						},
						"notify_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "HTTP callback URL，必填 如果 `NotifyType` 是 集合 到 `URL`。",
						},
						"aws_sqs": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "AWS SQS queue. 此 参数 为必填项 如果 `NotifyType` 是 `AWS-SQS`.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sqs_region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "地域 的 SQS queue。",
									},
									"sqs_queue_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "名称 SQS queue。",
									},
									"s3_secret_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "键 ID 必填 到 read 从/write 到 SQS queue。",
									},
									"s3_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "键 必填 到 read 从/write 到 SQS queue。",
									},
								},
							},
						},
					},
				},
			},

			"resource_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "资源 ID，您 need 到 ensure 该 corresponding 资源 是 open. 默认为 账号 main 资源 ID。",
			},
		},
	}
}

func resourceTencentCloudMpsScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_schedule.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = mps.NewCreateScheduleRequest()
		response   = mps.NewCreateScheduleResponse()
		scheduleId string
	)
	if v, ok := d.GetOk("schedule_name"); ok {
		request.ScheduleName = helper.String(v.(string))
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
					if formatsSet[i] != nil {
						formats := formatsSet[i].(string)
						cosFileUploadTrigger.Formats = append(cosFileUploadTrigger.Formats, &formats)
					}
				}
			}
			workflowTrigger.CosFileUploadTrigger = &cosFileUploadTrigger
		}
		if awsS3FileUploadTriggerMap, ok := helper.InterfaceToMap(dMap, "aws_s3_file_upload_trigger"); ok {
			awsS3FileUploadTrigger := mps.AwsS3FileUploadTrigger{}
			if v, ok := awsS3FileUploadTriggerMap["s3_bucket"]; ok {
				awsS3FileUploadTrigger.S3Bucket = helper.String(v.(string))
			}
			if v, ok := awsS3FileUploadTriggerMap["s3_region"]; ok {
				awsS3FileUploadTrigger.S3Region = helper.String(v.(string))
			}
			if v, ok := awsS3FileUploadTriggerMap["dir"]; ok {
				awsS3FileUploadTrigger.Dir = helper.String(v.(string))
			}
			if v, ok := awsS3FileUploadTriggerMap["formats"]; ok {
				formatsSet := v.(*schema.Set).List()
				for i := range formatsSet {
					if formatsSet[i] != nil {
						formats := formatsSet[i].(string)
						awsS3FileUploadTrigger.Formats = append(awsS3FileUploadTrigger.Formats, &formats)
					}
				}
			}
			if v, ok := awsS3FileUploadTriggerMap["s3_secret_id"]; ok {
				awsS3FileUploadTrigger.S3SecretId = helper.String(v.(string))
			}
			if v, ok := awsS3FileUploadTriggerMap["s3_secret_key"]; ok {
				awsS3FileUploadTrigger.S3SecretKey = helper.String(v.(string))
			}
			if awsSQSMap, ok := helper.InterfaceToMap(awsS3FileUploadTriggerMap, "aws_sqs"); ok {
				awsSQS := mps.AwsSQS{}
				if v, ok := awsSQSMap["sqs_region"]; ok {
					awsSQS.SQSRegion = helper.String(v.(string))
				}
				if v, ok := awsSQSMap["sqs_queue_name"]; ok {
					awsSQS.SQSQueueName = helper.String(v.(string))
				}
				if v, ok := awsSQSMap["s3_secret_id"]; ok {
					awsSQS.S3SecretId = helper.String(v.(string))
				}
				if v, ok := awsSQSMap["s3_secret_key"]; ok {
					awsSQS.S3SecretKey = helper.String(v.(string))
				}
				awsS3FileUploadTrigger.AwsSQS = &awsSQS
			}
			workflowTrigger.AwsS3FileUploadTrigger = &awsS3FileUploadTrigger
		}
		request.Trigger = &workflowTrigger
	}

	if v, ok := d.GetOk("activities"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			activity := mps.Activity{}
			if v, ok := dMap["activity_type"]; ok {
				activity.ActivityType = helper.String(v.(string))
			}
			if v, ok := dMap["reardrive_index"]; ok {
				reardriveIndexSet := v.(*schema.Set).List()
				for i := range reardriveIndexSet {
					reardriveIndex := reardriveIndexSet[i].(int)
					activity.ReardriveIndex = append(activity.ReardriveIndex, helper.IntInt64(reardriveIndex))
				}
			}
			if activityParaMap, ok := helper.InterfaceToMap(dMap, "activity_para"); ok {
				activityPara := mps.ActivityPara{}
				if transcodeTaskMap, ok := helper.InterfaceToMap(activityParaMap, "transcode_task"); ok {
					transcodeTaskInput := mps.TranscodeTaskInput{}
					if v, ok := transcodeTaskMap["definition"]; ok {
						transcodeTaskInput.Definition = helper.IntUint64(v.(int))
					}
					if rawParameterMap, ok := helper.InterfaceToMap(transcodeTaskMap, "raw_parameter"); ok {
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
					if overrideParameterMap, ok := helper.InterfaceToMap(transcodeTaskMap, "override_parameter"); ok {
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
						if v, ok := overrideParameterMap["addon_audio_stream"]; ok {
							for _, item := range v.([]interface{}) {
								addonAudioStreamMap := item.(map[string]interface{})
								mediaInputInfo := mps.MediaInputInfo{}
								if v, ok := addonAudioStreamMap["type"]; ok {
									mediaInputInfo.Type = helper.String(v.(string))
								}
								if cosInputInfoMap, ok := helper.InterfaceToMap(addonAudioStreamMap, "cos_input_info"); ok {
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
								if urlInputInfoMap, ok := helper.InterfaceToMap(addonAudioStreamMap, "url_input_info"); ok {
									urlInputInfo := mps.UrlInputInfo{}
									if v, ok := urlInputInfoMap["url"]; ok {
										urlInputInfo.Url = helper.String(v.(string))
									}
									mediaInputInfo.UrlInputInfo = &urlInputInfo
								}
								if s3InputInfoMap, ok := helper.InterfaceToMap(addonAudioStreamMap, "s3_input_info"); ok {
									s3InputInfo := mps.S3InputInfo{}
									if v, ok := s3InputInfoMap["s3_bucket"]; ok {
										s3InputInfo.S3Bucket = helper.String(v.(string))
									}
									if v, ok := s3InputInfoMap["s3_region"]; ok {
										s3InputInfo.S3Region = helper.String(v.(string))
									}
									if v, ok := s3InputInfoMap["s3_object"]; ok {
										s3InputInfo.S3Object = helper.String(v.(string))
									}
									if v, ok := s3InputInfoMap["s3_secret_id"]; ok {
										s3InputInfo.S3SecretId = helper.String(v.(string))
									}
									if v, ok := s3InputInfoMap["s3_secret_key"]; ok {
										s3InputInfo.S3SecretKey = helper.String(v.(string))
									}
									mediaInputInfo.S3InputInfo = &s3InputInfo
								}
								overrideTranscodeParameter.AddonAudioStream = append(overrideTranscodeParameter.AddonAudioStream, &mediaInputInfo)
							}
						}
						if v, ok := overrideParameterMap["std_ext_info"]; ok {
							overrideTranscodeParameter.StdExtInfo = helper.String(v.(string))
						}
						if v, ok := overrideParameterMap["add_on_subtitles"]; ok {
							for _, item := range v.([]interface{}) {
								addOnSubtitlesMap := item.(map[string]interface{})
								addOnSubtitle := mps.AddOnSubtitle{}
								if v, ok := addOnSubtitlesMap["type"]; ok {
									addOnSubtitle.Type = helper.String(v.(string))
								}
								if subtitleMap, ok := helper.InterfaceToMap(addOnSubtitlesMap, "subtitle"); ok {
									mediaInputInfo := mps.MediaInputInfo{}
									if v, ok := subtitleMap["type"]; ok {
										mediaInputInfo.Type = helper.String(v.(string))
									}
									if cosInputInfoMap, ok := helper.InterfaceToMap(subtitleMap, "cos_input_info"); ok {
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
									if urlInputInfoMap, ok := helper.InterfaceToMap(subtitleMap, "url_input_info"); ok {
										urlInputInfo := mps.UrlInputInfo{}
										if v, ok := urlInputInfoMap["url"]; ok {
											urlInputInfo.Url = helper.String(v.(string))
										}
										mediaInputInfo.UrlInputInfo = &urlInputInfo
									}
									if s3InputInfoMap, ok := helper.InterfaceToMap(subtitleMap, "s3_input_info"); ok {
										s3InputInfo := mps.S3InputInfo{}
										if v, ok := s3InputInfoMap["s3_bucket"]; ok {
											s3InputInfo.S3Bucket = helper.String(v.(string))
										}
										if v, ok := s3InputInfoMap["s3_region"]; ok {
											s3InputInfo.S3Region = helper.String(v.(string))
										}
										if v, ok := s3InputInfoMap["s3_object"]; ok {
											s3InputInfo.S3Object = helper.String(v.(string))
										}
										if v, ok := s3InputInfoMap["s3_secret_id"]; ok {
											s3InputInfo.S3SecretId = helper.String(v.(string))
										}
										if v, ok := s3InputInfoMap["s3_secret_key"]; ok {
											s3InputInfo.S3SecretKey = helper.String(v.(string))
										}
										mediaInputInfo.S3InputInfo = &s3InputInfo
									}
									addOnSubtitle.Subtitle = &mediaInputInfo
								}
								overrideTranscodeParameter.AddOnSubtitles = append(overrideTranscodeParameter.AddOnSubtitles, &addOnSubtitle)
							}
						}
						transcodeTaskInput.OverrideParameter = &overrideTranscodeParameter
					}
					if v, ok := transcodeTaskMap["watermark_set"]; ok {
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
										if s3InputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "s3_input_info"); ok {
											s3InputInfo := mps.S3InputInfo{}
											if v, ok := s3InputInfoMap["s3_bucket"]; ok {
												s3InputInfo.S3Bucket = helper.String(v.(string))
											}
											if v, ok := s3InputInfoMap["s3_region"]; ok {
												s3InputInfo.S3Region = helper.String(v.(string))
											}
											if v, ok := s3InputInfoMap["s3_object"]; ok {
												s3InputInfo.S3Object = helper.String(v.(string))
											}
											if v, ok := s3InputInfoMap["s3_secret_id"]; ok {
												s3InputInfo.S3SecretId = helper.String(v.(string))
											}
											if v, ok := s3InputInfoMap["s3_secret_key"]; ok {
												s3InputInfo.S3SecretKey = helper.String(v.(string))
											}
											mediaInputInfo.S3InputInfo = &s3InputInfo
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
					if v, ok := transcodeTaskMap["mosaic_set"]; ok {
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
					if v, ok := transcodeTaskMap["start_time_offset"]; ok {
						transcodeTaskInput.StartTimeOffset = helper.Float64(v.(float64))
					}
					if v, ok := transcodeTaskMap["end_time_offset"]; ok {
						transcodeTaskInput.EndTimeOffset = helper.Float64(v.(float64))
					}
					if outputStorageMap, ok := helper.InterfaceToMap(transcodeTaskMap, "output_storage"); ok {
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
						if s3OutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "s3_output_storage"); ok {
							s3OutputStorage := mps.S3OutputStorage{}
							if v, ok := s3OutputStorageMap["s3_bucket"]; ok {
								s3OutputStorage.S3Bucket = helper.String(v.(string))
							}
							if v, ok := s3OutputStorageMap["s3_region"]; ok {
								s3OutputStorage.S3Region = helper.String(v.(string))
							}
							if v, ok := s3OutputStorageMap["s3_secret_id"]; ok {
								s3OutputStorage.S3SecretId = helper.String(v.(string))
							}
							if v, ok := s3OutputStorageMap["s3_secret_key"]; ok {
								s3OutputStorage.S3SecretKey = helper.String(v.(string))
							}
							taskOutputStorage.S3OutputStorage = &s3OutputStorage
						}
						transcodeTaskInput.OutputStorage = &taskOutputStorage
					}
					if v, ok := transcodeTaskMap["output_object_path"]; ok {
						transcodeTaskInput.OutputObjectPath = helper.String(v.(string))
					}
					if v, ok := transcodeTaskMap["segment_object_name"]; ok {
						transcodeTaskInput.SegmentObjectName = helper.String(v.(string))
					}
					if objectNumberFormatMap, ok := helper.InterfaceToMap(transcodeTaskMap, "object_number_format"); ok {
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
					if headTailParameterMap, ok := helper.InterfaceToMap(transcodeTaskMap, "head_tail_parameter"); ok {
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
								if s3InputInfoMap, ok := helper.InterfaceToMap(headSetMap, "s3_input_info"); ok {
									s3InputInfo := mps.S3InputInfo{}
									if v, ok := s3InputInfoMap["s3_bucket"]; ok {
										s3InputInfo.S3Bucket = helper.String(v.(string))
									}
									if v, ok := s3InputInfoMap["s3_region"]; ok {
										s3InputInfo.S3Region = helper.String(v.(string))
									}
									if v, ok := s3InputInfoMap["s3_object"]; ok {
										s3InputInfo.S3Object = helper.String(v.(string))
									}
									if v, ok := s3InputInfoMap["s3_secret_id"]; ok {
										s3InputInfo.S3SecretId = helper.String(v.(string))
									}
									if v, ok := s3InputInfoMap["s3_secret_key"]; ok {
										s3InputInfo.S3SecretKey = helper.String(v.(string))
									}
									mediaInputInfo.S3InputInfo = &s3InputInfo
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
								if s3InputInfoMap, ok := helper.InterfaceToMap(tailSetMap, "s3_input_info"); ok {
									s3InputInfo := mps.S3InputInfo{}
									if v, ok := s3InputInfoMap["s3_bucket"]; ok {
										s3InputInfo.S3Bucket = helper.String(v.(string))
									}
									if v, ok := s3InputInfoMap["s3_region"]; ok {
										s3InputInfo.S3Region = helper.String(v.(string))
									}
									if v, ok := s3InputInfoMap["s3_object"]; ok {
										s3InputInfo.S3Object = helper.String(v.(string))
									}
									if v, ok := s3InputInfoMap["s3_secret_id"]; ok {
										s3InputInfo.S3SecretId = helper.String(v.(string))
									}
									if v, ok := s3InputInfoMap["s3_secret_key"]; ok {
										s3InputInfo.S3SecretKey = helper.String(v.(string))
									}
									mediaInputInfo.S3InputInfo = &s3InputInfo
								}
								headTailParameter.TailSet = append(headTailParameter.TailSet, &mediaInputInfo)
							}
						}
						transcodeTaskInput.HeadTailParameter = &headTailParameter
					}
					activityPara.TranscodeTask = &transcodeTaskInput
				}
				if animatedGraphicTaskMap, ok := helper.InterfaceToMap(activityParaMap, "animated_graphic_task"); ok {
					animatedGraphicTaskInput := mps.AnimatedGraphicTaskInput{}
					if v, ok := animatedGraphicTaskMap["definition"]; ok {
						animatedGraphicTaskInput.Definition = helper.IntUint64(v.(int))
					}
					if v, ok := animatedGraphicTaskMap["start_time_offset"]; ok {
						animatedGraphicTaskInput.StartTimeOffset = helper.Float64(v.(float64))
					}
					if v, ok := animatedGraphicTaskMap["end_time_offset"]; ok {
						animatedGraphicTaskInput.EndTimeOffset = helper.Float64(v.(float64))
					}
					if outputStorageMap, ok := helper.InterfaceToMap(animatedGraphicTaskMap, "output_storage"); ok {
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
						if s3OutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "s3_output_storage"); ok {
							s3OutputStorage := mps.S3OutputStorage{}
							if v, ok := s3OutputStorageMap["s3_bucket"]; ok {
								s3OutputStorage.S3Bucket = helper.String(v.(string))
							}
							if v, ok := s3OutputStorageMap["s3_region"]; ok {
								s3OutputStorage.S3Region = helper.String(v.(string))
							}
							if v, ok := s3OutputStorageMap["s3_secret_id"]; ok {
								s3OutputStorage.S3SecretId = helper.String(v.(string))
							}
							if v, ok := s3OutputStorageMap["s3_secret_key"]; ok {
								s3OutputStorage.S3SecretKey = helper.String(v.(string))
							}
							taskOutputStorage.S3OutputStorage = &s3OutputStorage
						}
						animatedGraphicTaskInput.OutputStorage = &taskOutputStorage
					}
					if v, ok := animatedGraphicTaskMap["output_object_path"]; ok {
						animatedGraphicTaskInput.OutputObjectPath = helper.String(v.(string))
					}
					activityPara.AnimatedGraphicTask = &animatedGraphicTaskInput
				}
				if snapshotByTimeOffsetTaskMap, ok := helper.InterfaceToMap(activityParaMap, "snapshot_by_time_offset_task"); ok {
					snapshotByTimeOffsetTaskInput := mps.SnapshotByTimeOffsetTaskInput{}
					if v, ok := snapshotByTimeOffsetTaskMap["definition"]; ok {
						snapshotByTimeOffsetTaskInput.Definition = helper.IntUint64(v.(int))
					}
					if v, ok := snapshotByTimeOffsetTaskMap["ext_time_offset_set"]; ok {
						extTimeOffsetSetSet := v.(*schema.Set).List()
						for i := range extTimeOffsetSetSet {
							if extTimeOffsetSetSet[i] != nil {
								extTimeOffsetSet := extTimeOffsetSetSet[i].(string)
								snapshotByTimeOffsetTaskInput.ExtTimeOffsetSet = append(snapshotByTimeOffsetTaskInput.ExtTimeOffsetSet, &extTimeOffsetSet)
							}
						}
					}

					if v, ok := snapshotByTimeOffsetTaskMap["watermark_set"]; ok {
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
										if s3InputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "s3_input_info"); ok {
											s3InputInfo := mps.S3InputInfo{}
											if v, ok := s3InputInfoMap["s3_bucket"]; ok {
												s3InputInfo.S3Bucket = helper.String(v.(string))
											}
											if v, ok := s3InputInfoMap["s3_region"]; ok {
												s3InputInfo.S3Region = helper.String(v.(string))
											}
											if v, ok := s3InputInfoMap["s3_object"]; ok {
												s3InputInfo.S3Object = helper.String(v.(string))
											}
											if v, ok := s3InputInfoMap["s3_secret_id"]; ok {
												s3InputInfo.S3SecretId = helper.String(v.(string))
											}
											if v, ok := s3InputInfoMap["s3_secret_key"]; ok {
												s3InputInfo.S3SecretKey = helper.String(v.(string))
											}
											mediaInputInfo.S3InputInfo = &s3InputInfo
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
					if outputStorageMap, ok := helper.InterfaceToMap(snapshotByTimeOffsetTaskMap, "output_storage"); ok {
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
						if s3OutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "s3_output_storage"); ok {
							s3OutputStorage := mps.S3OutputStorage{}
							if v, ok := s3OutputStorageMap["s3_bucket"]; ok {
								s3OutputStorage.S3Bucket = helper.String(v.(string))
							}
							if v, ok := s3OutputStorageMap["s3_region"]; ok {
								s3OutputStorage.S3Region = helper.String(v.(string))
							}
							if v, ok := s3OutputStorageMap["s3_secret_id"]; ok {
								s3OutputStorage.S3SecretId = helper.String(v.(string))
							}
							if v, ok := s3OutputStorageMap["s3_secret_key"]; ok {
								s3OutputStorage.S3SecretKey = helper.String(v.(string))
							}
							taskOutputStorage.S3OutputStorage = &s3OutputStorage
						}
						snapshotByTimeOffsetTaskInput.OutputStorage = &taskOutputStorage
					}
					if v, ok := snapshotByTimeOffsetTaskMap["output_object_path"]; ok {
						snapshotByTimeOffsetTaskInput.OutputObjectPath = helper.String(v.(string))
					}
					if objectNumberFormatMap, ok := helper.InterfaceToMap(snapshotByTimeOffsetTaskMap, "object_number_format"); ok {
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
					activityPara.SnapshotByTimeOffsetTask = &snapshotByTimeOffsetTaskInput
				}
				if sampleSnapshotTaskMap, ok := helper.InterfaceToMap(activityParaMap, "sample_snapshot_task"); ok {
					sampleSnapshotTaskInput := mps.SampleSnapshotTaskInput{}
					if v, ok := sampleSnapshotTaskMap["definition"]; ok {
						sampleSnapshotTaskInput.Definition = helper.IntUint64(v.(int))
					}
					if v, ok := sampleSnapshotTaskMap["watermark_set"]; ok {
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
										if s3InputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "s3_input_info"); ok {
											s3InputInfo := mps.S3InputInfo{}
											if v, ok := s3InputInfoMap["s3_bucket"]; ok {
												s3InputInfo.S3Bucket = helper.String(v.(string))
											}
											if v, ok := s3InputInfoMap["s3_region"]; ok {
												s3InputInfo.S3Region = helper.String(v.(string))
											}
											if v, ok := s3InputInfoMap["s3_object"]; ok {
												s3InputInfo.S3Object = helper.String(v.(string))
											}
											if v, ok := s3InputInfoMap["s3_secret_id"]; ok {
												s3InputInfo.S3SecretId = helper.String(v.(string))
											}
											if v, ok := s3InputInfoMap["s3_secret_key"]; ok {
												s3InputInfo.S3SecretKey = helper.String(v.(string))
											}
											mediaInputInfo.S3InputInfo = &s3InputInfo
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
					if outputStorageMap, ok := helper.InterfaceToMap(sampleSnapshotTaskMap, "output_storage"); ok {
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
						if s3OutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "s3_output_storage"); ok {
							s3OutputStorage := mps.S3OutputStorage{}
							if v, ok := s3OutputStorageMap["s3_bucket"]; ok {
								s3OutputStorage.S3Bucket = helper.String(v.(string))
							}
							if v, ok := s3OutputStorageMap["s3_region"]; ok {
								s3OutputStorage.S3Region = helper.String(v.(string))
							}
							if v, ok := s3OutputStorageMap["s3_secret_id"]; ok {
								s3OutputStorage.S3SecretId = helper.String(v.(string))
							}
							if v, ok := s3OutputStorageMap["s3_secret_key"]; ok {
								s3OutputStorage.S3SecretKey = helper.String(v.(string))
							}
							taskOutputStorage.S3OutputStorage = &s3OutputStorage
						}
						sampleSnapshotTaskInput.OutputStorage = &taskOutputStorage
					}
					if v, ok := sampleSnapshotTaskMap["output_object_path"]; ok {
						sampleSnapshotTaskInput.OutputObjectPath = helper.String(v.(string))
					}
					if objectNumberFormatMap, ok := helper.InterfaceToMap(sampleSnapshotTaskMap, "object_number_format"); ok {
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
					activityPara.SampleSnapshotTask = &sampleSnapshotTaskInput
				}
				if imageSpriteTaskMap, ok := helper.InterfaceToMap(activityParaMap, "image_sprite_task"); ok {
					imageSpriteTaskInput := mps.ImageSpriteTaskInput{}
					if v, ok := imageSpriteTaskMap["definition"]; ok {
						imageSpriteTaskInput.Definition = helper.IntUint64(v.(int))
					}
					if outputStorageMap, ok := helper.InterfaceToMap(imageSpriteTaskMap, "output_storage"); ok {
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
						if s3OutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "s3_output_storage"); ok {
							s3OutputStorage := mps.S3OutputStorage{}
							if v, ok := s3OutputStorageMap["s3_bucket"]; ok {
								s3OutputStorage.S3Bucket = helper.String(v.(string))
							}
							if v, ok := s3OutputStorageMap["s3_region"]; ok {
								s3OutputStorage.S3Region = helper.String(v.(string))
							}
							if v, ok := s3OutputStorageMap["s3_secret_id"]; ok {
								s3OutputStorage.S3SecretId = helper.String(v.(string))
							}
							if v, ok := s3OutputStorageMap["s3_secret_key"]; ok {
								s3OutputStorage.S3SecretKey = helper.String(v.(string))
							}
							taskOutputStorage.S3OutputStorage = &s3OutputStorage
						}
						imageSpriteTaskInput.OutputStorage = &taskOutputStorage
					}
					if v, ok := imageSpriteTaskMap["output_object_path"]; ok {
						imageSpriteTaskInput.OutputObjectPath = helper.String(v.(string))
					}
					if v, ok := imageSpriteTaskMap["web_vtt_object_name"]; ok {
						imageSpriteTaskInput.WebVttObjectName = helper.String(v.(string))
					}
					if objectNumberFormatMap, ok := helper.InterfaceToMap(imageSpriteTaskMap, "object_number_format"); ok {
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
					activityPara.ImageSpriteTask = &imageSpriteTaskInput
				}
				if adaptiveDynamicStreamingTaskMap, ok := helper.InterfaceToMap(activityParaMap, "adaptive_dynamic_streaming_task"); ok {
					adaptiveDynamicStreamingTaskInput := mps.AdaptiveDynamicStreamingTaskInput{}
					if v, ok := adaptiveDynamicStreamingTaskMap["definition"]; ok {
						adaptiveDynamicStreamingTaskInput.Definition = helper.IntUint64(v.(int))
					}
					if v, ok := adaptiveDynamicStreamingTaskMap["watermark_set"]; ok {
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
										if s3InputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "s3_input_info"); ok {
											s3InputInfo := mps.S3InputInfo{}
											if v, ok := s3InputInfoMap["s3_bucket"]; ok {
												s3InputInfo.S3Bucket = helper.String(v.(string))
											}
											if v, ok := s3InputInfoMap["s3_region"]; ok {
												s3InputInfo.S3Region = helper.String(v.(string))
											}
											if v, ok := s3InputInfoMap["s3_object"]; ok {
												s3InputInfo.S3Object = helper.String(v.(string))
											}
											if v, ok := s3InputInfoMap["s3_secret_id"]; ok {
												s3InputInfo.S3SecretId = helper.String(v.(string))
											}
											if v, ok := s3InputInfoMap["s3_secret_key"]; ok {
												s3InputInfo.S3SecretKey = helper.String(v.(string))
											}
											mediaInputInfo.S3InputInfo = &s3InputInfo
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
					if outputStorageMap, ok := helper.InterfaceToMap(adaptiveDynamicStreamingTaskMap, "output_storage"); ok {
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
						if s3OutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "s3_output_storage"); ok {
							s3OutputStorage := mps.S3OutputStorage{}
							if v, ok := s3OutputStorageMap["s3_bucket"]; ok {
								s3OutputStorage.S3Bucket = helper.String(v.(string))
							}
							if v, ok := s3OutputStorageMap["s3_region"]; ok {
								s3OutputStorage.S3Region = helper.String(v.(string))
							}
							if v, ok := s3OutputStorageMap["s3_secret_id"]; ok {
								s3OutputStorage.S3SecretId = helper.String(v.(string))
							}
							if v, ok := s3OutputStorageMap["s3_secret_key"]; ok {
								s3OutputStorage.S3SecretKey = helper.String(v.(string))
							}
							taskOutputStorage.S3OutputStorage = &s3OutputStorage
						}
						adaptiveDynamicStreamingTaskInput.OutputStorage = &taskOutputStorage
					}
					if v, ok := adaptiveDynamicStreamingTaskMap["output_object_path"]; ok {
						adaptiveDynamicStreamingTaskInput.OutputObjectPath = helper.String(v.(string))
					}
					if v, ok := adaptiveDynamicStreamingTaskMap["sub_stream_object_name"]; ok {
						adaptiveDynamicStreamingTaskInput.SubStreamObjectName = helper.String(v.(string))
					}
					if v, ok := adaptiveDynamicStreamingTaskMap["segment_object_name"]; ok {
						adaptiveDynamicStreamingTaskInput.SegmentObjectName = helper.String(v.(string))
					}
					if v, ok := adaptiveDynamicStreamingTaskMap["add_on_subtitles"]; ok {
						for _, item := range v.([]interface{}) {
							addOnSubtitlesMap := item.(map[string]interface{})
							addOnSubtitle := mps.AddOnSubtitle{}
							if v, ok := addOnSubtitlesMap["type"]; ok {
								addOnSubtitle.Type = helper.String(v.(string))
							}
							if subtitleMap, ok := helper.InterfaceToMap(addOnSubtitlesMap, "subtitle"); ok {
								mediaInputInfo := mps.MediaInputInfo{}
								if v, ok := subtitleMap["type"]; ok {
									mediaInputInfo.Type = helper.String(v.(string))
								}
								if cosInputInfoMap, ok := helper.InterfaceToMap(subtitleMap, "cos_input_info"); ok {
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
								if urlInputInfoMap, ok := helper.InterfaceToMap(subtitleMap, "url_input_info"); ok {
									urlInputInfo := mps.UrlInputInfo{}
									if v, ok := urlInputInfoMap["url"]; ok {
										urlInputInfo.Url = helper.String(v.(string))
									}
									mediaInputInfo.UrlInputInfo = &urlInputInfo
								}
								if s3InputInfoMap, ok := helper.InterfaceToMap(subtitleMap, "s3_input_info"); ok {
									s3InputInfo := mps.S3InputInfo{}
									if v, ok := s3InputInfoMap["s3_bucket"]; ok {
										s3InputInfo.S3Bucket = helper.String(v.(string))
									}
									if v, ok := s3InputInfoMap["s3_region"]; ok {
										s3InputInfo.S3Region = helper.String(v.(string))
									}
									if v, ok := s3InputInfoMap["s3_object"]; ok {
										s3InputInfo.S3Object = helper.String(v.(string))
									}
									if v, ok := s3InputInfoMap["s3_secret_id"]; ok {
										s3InputInfo.S3SecretId = helper.String(v.(string))
									}
									if v, ok := s3InputInfoMap["s3_secret_key"]; ok {
										s3InputInfo.S3SecretKey = helper.String(v.(string))
									}
									mediaInputInfo.S3InputInfo = &s3InputInfo
								}
								addOnSubtitle.Subtitle = &mediaInputInfo
							}
							adaptiveDynamicStreamingTaskInput.AddOnSubtitles = append(adaptiveDynamicStreamingTaskInput.AddOnSubtitles, &addOnSubtitle)
						}
					}
					activityPara.AdaptiveDynamicStreamingTask = &adaptiveDynamicStreamingTaskInput
				}
				if aiContentReviewTaskMap, ok := helper.InterfaceToMap(activityParaMap, "ai_content_review_task"); ok {
					aiContentReviewTaskInput := mps.AiContentReviewTaskInput{}
					if v, ok := aiContentReviewTaskMap["definition"]; ok {
						aiContentReviewTaskInput.Definition = helper.IntUint64(v.(int))
					}
					activityPara.AiContentReviewTask = &aiContentReviewTaskInput
				}
				if aiAnalysisTaskMap, ok := helper.InterfaceToMap(activityParaMap, "ai_analysis_task"); ok {
					aiAnalysisTaskInput := mps.AiAnalysisTaskInput{}
					if v, ok := aiAnalysisTaskMap["definition"]; ok {
						aiAnalysisTaskInput.Definition = helper.IntUint64(v.(int))
					}
					if v, ok := aiAnalysisTaskMap["extended_parameter"]; ok {
						aiAnalysisTaskInput.ExtendedParameter = helper.String(v.(string))
					}
					activityPara.AiAnalysisTask = &aiAnalysisTaskInput
				}
				if aiRecognitionTaskMap, ok := helper.InterfaceToMap(activityParaMap, "ai_recognition_task"); ok {
					aiRecognitionTaskInput := mps.AiRecognitionTaskInput{}
					if v, ok := aiRecognitionTaskMap["definition"]; ok {
						aiRecognitionTaskInput.Definition = helper.IntUint64(v.(int))
					}
					activityPara.AiRecognitionTask = &aiRecognitionTaskInput
				}
				activity.ActivityPara = &activityPara
			}
			request.Activities = append(request.Activities, &activity)
		}
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
		if s3OutputStorageMap, ok := helper.InterfaceToMap(dMap, "s3_output_storage"); ok {
			s3OutputStorage := mps.S3OutputStorage{}
			if v, ok := s3OutputStorageMap["s3_bucket"]; ok {
				s3OutputStorage.S3Bucket = helper.String(v.(string))
			}
			if v, ok := s3OutputStorageMap["s3_region"]; ok {
				s3OutputStorage.S3Region = helper.String(v.(string))
			}
			if v, ok := s3OutputStorageMap["s3_secret_id"]; ok {
				s3OutputStorage.S3SecretId = helper.String(v.(string))
			}
			if v, ok := s3OutputStorageMap["s3_secret_key"]; ok {
				s3OutputStorage.S3SecretKey = helper.String(v.(string))
			}
			taskOutputStorage.S3OutputStorage = &s3OutputStorage
		}
		request.OutputStorage = &taskOutputStorage
	}

	if v, ok := d.GetOk("output_dir"); ok {
		request.OutputDir = helper.String(v.(string))
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
		if awsSQSMap, ok := helper.InterfaceToMap(dMap, "aws_sqs"); ok {
			awsSQS := mps.AwsSQS{}
			if v, ok := awsSQSMap["sqs_region"]; ok {
				awsSQS.SQSRegion = helper.String(v.(string))
			}
			if v, ok := awsSQSMap["sqs_queue_name"]; ok {
				awsSQS.SQSQueueName = helper.String(v.(string))
			}
			if v, ok := awsSQSMap["s3_secret_id"]; ok {
				awsSQS.S3SecretId = helper.String(v.(string))
			}
			if v, ok := awsSQSMap["s3_secret_key"]; ok {
				awsSQS.S3SecretKey = helper.String(v.(string))
			}
			taskNotifyConfig.AwsSQS = &awsSQS
		}
		request.TaskNotifyConfig = &taskNotifyConfig
	}

	if v, ok := d.GetOk("resource_id"); ok {
		request.ResourceId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().CreateSchedule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps schedule failed, reason:%+v", logId, err)
		return err
	}

	scheduleId = helper.Int64ToStr(*response.Response.ScheduleId)
	d.SetId(scheduleId)

	return resourceTencentCloudMpsScheduleRead(d, meta)
}

func resourceTencentCloudMpsScheduleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_schedule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	scheduleId := d.Id()

	schedules, err := service.DescribeMpsScheduleById(ctx, &scheduleId)
	if err != nil {
		return err
	}

	if len(schedules) == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsSchedule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	schedule := schedules[0]

	if schedule.ScheduleName != nil {
		_ = d.Set("schedule_name", schedule.ScheduleName)
	}

	if schedule.Trigger != nil {
		triggerMap := map[string]interface{}{}

		if schedule.Trigger.Type != nil {
			triggerMap["type"] = schedule.Trigger.Type
		}

		if schedule.Trigger.CosFileUploadTrigger != nil {
			cosFileUploadTriggerMap := map[string]interface{}{}

			if schedule.Trigger.CosFileUploadTrigger.Bucket != nil {
				cosFileUploadTriggerMap["bucket"] = schedule.Trigger.CosFileUploadTrigger.Bucket
			}

			if schedule.Trigger.CosFileUploadTrigger.Region != nil {
				cosFileUploadTriggerMap["region"] = schedule.Trigger.CosFileUploadTrigger.Region
			}

			if schedule.Trigger.CosFileUploadTrigger.Dir != nil {
				cosFileUploadTriggerMap["dir"] = schedule.Trigger.CosFileUploadTrigger.Dir
			}

			if schedule.Trigger.CosFileUploadTrigger.Formats != nil {
				cosFileUploadTriggerMap["formats"] = schedule.Trigger.CosFileUploadTrigger.Formats
			}

			triggerMap["cos_file_upload_trigger"] = []interface{}{cosFileUploadTriggerMap}
		}

		if schedule.Trigger.AwsS3FileUploadTrigger != nil {
			awsS3FileUploadTriggerMap := map[string]interface{}{}

			if schedule.Trigger.AwsS3FileUploadTrigger.S3Bucket != nil {
				awsS3FileUploadTriggerMap["s3_bucket"] = schedule.Trigger.AwsS3FileUploadTrigger.S3Bucket
			}

			if schedule.Trigger.AwsS3FileUploadTrigger.S3Region != nil {
				awsS3FileUploadTriggerMap["s3_region"] = schedule.Trigger.AwsS3FileUploadTrigger.S3Region
			}

			if schedule.Trigger.AwsS3FileUploadTrigger.Dir != nil {
				awsS3FileUploadTriggerMap["dir"] = schedule.Trigger.AwsS3FileUploadTrigger.Dir
			}

			if schedule.Trigger.AwsS3FileUploadTrigger.Formats != nil {
				awsS3FileUploadTriggerMap["formats"] = schedule.Trigger.AwsS3FileUploadTrigger.Formats
			}

			if schedule.Trigger.AwsS3FileUploadTrigger.S3SecretId != nil {
				awsS3FileUploadTriggerMap["s3_secret_id"] = schedule.Trigger.AwsS3FileUploadTrigger.S3SecretId
			}

			if schedule.Trigger.AwsS3FileUploadTrigger.S3SecretKey != nil {
				awsS3FileUploadTriggerMap["s3_secret_key"] = schedule.Trigger.AwsS3FileUploadTrigger.S3SecretKey
			}

			if schedule.Trigger.AwsS3FileUploadTrigger.AwsSQS != nil {
				awsSQSMap := map[string]interface{}{}

				if schedule.Trigger.AwsS3FileUploadTrigger.AwsSQS.SQSRegion != nil {
					awsSQSMap["sqs_region"] = schedule.Trigger.AwsS3FileUploadTrigger.AwsSQS.SQSRegion
				}

				if schedule.Trigger.AwsS3FileUploadTrigger.AwsSQS.SQSQueueName != nil {
					awsSQSMap["sqs_queue_name"] = schedule.Trigger.AwsS3FileUploadTrigger.AwsSQS.SQSQueueName
				}

				if schedule.Trigger.AwsS3FileUploadTrigger.AwsSQS.S3SecretId != nil {
					awsSQSMap["s3_secret_id"] = schedule.Trigger.AwsS3FileUploadTrigger.AwsSQS.S3SecretId
				}

				if schedule.Trigger.AwsS3FileUploadTrigger.AwsSQS.S3SecretKey != nil {
					awsSQSMap["s3_secret_key"] = schedule.Trigger.AwsS3FileUploadTrigger.AwsSQS.S3SecretKey
				}

				awsS3FileUploadTriggerMap["aws_sqs"] = []interface{}{awsSQSMap}
			}

			triggerMap["aws_s3_file_upload_trigger"] = []interface{}{awsS3FileUploadTriggerMap}
		}

		_ = d.Set("trigger", []interface{}{triggerMap})
	}

	if schedule.Activities != nil {
		activitiesList := []interface{}{}
		for _, activity := range schedule.Activities {
			activitiesMap := map[string]interface{}{}

			if activity.ActivityType != nil {
				activitiesMap["activity_type"] = activity.ActivityType
			}

			if activity.ReardriveIndex != nil {
				activitiesMap["reardrive_index"] = activity.ReardriveIndex
			}

			if activity.ActivityPara != nil {
				activityParaMap := map[string]interface{}{}

				if activity.ActivityPara.TranscodeTask != nil {
					transcodeTaskMap := map[string]interface{}{}

					if activity.ActivityPara.TranscodeTask.Definition != nil {
						transcodeTaskMap["definition"] = activity.ActivityPara.TranscodeTask.Definition
					}

					if activity.ActivityPara.TranscodeTask.RawParameter != nil {
						rawParameterMap := map[string]interface{}{}

						if activity.ActivityPara.TranscodeTask.RawParameter.Container != nil {
							rawParameterMap["container"] = activity.ActivityPara.TranscodeTask.RawParameter.Container
						}

						if activity.ActivityPara.TranscodeTask.RawParameter.RemoveVideo != nil {
							rawParameterMap["remove_video"] = activity.ActivityPara.TranscodeTask.RawParameter.RemoveVideo
						}

						if activity.ActivityPara.TranscodeTask.RawParameter.RemoveAudio != nil {
							rawParameterMap["remove_audio"] = activity.ActivityPara.TranscodeTask.RawParameter.RemoveAudio
						}

						if activity.ActivityPara.TranscodeTask.RawParameter.VideoTemplate != nil {
							videoTemplateMap := map[string]interface{}{}

							if activity.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Codec != nil {
								videoTemplateMap["codec"] = activity.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Codec
							}

							if activity.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Fps != nil {
								videoTemplateMap["fps"] = activity.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Fps
							}

							if activity.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Bitrate != nil {
								videoTemplateMap["bitrate"] = activity.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Bitrate
							}

							if activity.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.ResolutionAdaptive != nil {
								videoTemplateMap["resolution_adaptive"] = activity.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.ResolutionAdaptive
							}

							if activity.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Width != nil {
								videoTemplateMap["width"] = activity.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Width
							}

							if activity.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Height != nil {
								videoTemplateMap["height"] = activity.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Height
							}

							if activity.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Gop != nil {
								videoTemplateMap["gop"] = activity.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Gop
							}

							if activity.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.FillType != nil {
								videoTemplateMap["fill_type"] = activity.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.FillType
							}

							if activity.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Vcrf != nil {
								videoTemplateMap["vcrf"] = activity.ActivityPara.TranscodeTask.RawParameter.VideoTemplate.Vcrf
							}

							rawParameterMap["video_template"] = []interface{}{videoTemplateMap}
						}

						if activity.ActivityPara.TranscodeTask.RawParameter.AudioTemplate != nil {
							audioTemplateMap := map[string]interface{}{}

							if activity.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.Codec != nil {
								audioTemplateMap["codec"] = activity.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.Codec
							}

							if activity.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.Bitrate != nil {
								audioTemplateMap["bitrate"] = activity.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.Bitrate
							}

							if activity.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.SampleRate != nil {
								audioTemplateMap["sample_rate"] = activity.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.SampleRate
							}

							if activity.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.AudioChannel != nil {
								audioTemplateMap["audio_channel"] = activity.ActivityPara.TranscodeTask.RawParameter.AudioTemplate.AudioChannel
							}

							rawParameterMap["audio_template"] = []interface{}{audioTemplateMap}
						}

						if activity.ActivityPara.TranscodeTask.RawParameter.TEHDConfig != nil {
							tEHDConfigMap := map[string]interface{}{}

							if activity.ActivityPara.TranscodeTask.RawParameter.TEHDConfig.Type != nil {
								tEHDConfigMap["type"] = activity.ActivityPara.TranscodeTask.RawParameter.TEHDConfig.Type
							}

							if activity.ActivityPara.TranscodeTask.RawParameter.TEHDConfig.MaxVideoBitrate != nil {
								tEHDConfigMap["max_video_bitrate"] = activity.ActivityPara.TranscodeTask.RawParameter.TEHDConfig.MaxVideoBitrate
							}

							rawParameterMap["tehd_config"] = []interface{}{tEHDConfigMap}
						}

						transcodeTaskMap["raw_parameter"] = []interface{}{rawParameterMap}
					}

					if activity.ActivityPara.TranscodeTask.OverrideParameter != nil {
						overrideParameterMap := map[string]interface{}{}

						if activity.ActivityPara.TranscodeTask.OverrideParameter.Container != nil {
							overrideParameterMap["container"] = activity.ActivityPara.TranscodeTask.OverrideParameter.Container
						}

						if activity.ActivityPara.TranscodeTask.OverrideParameter.RemoveVideo != nil {
							overrideParameterMap["remove_video"] = activity.ActivityPara.TranscodeTask.OverrideParameter.RemoveVideo
						}

						if activity.ActivityPara.TranscodeTask.OverrideParameter.RemoveAudio != nil {
							overrideParameterMap["remove_audio"] = activity.ActivityPara.TranscodeTask.OverrideParameter.RemoveAudio
						}

						if activity.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate != nil {
							videoTemplateMap := map[string]interface{}{}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Codec != nil {
								videoTemplateMap["codec"] = activity.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Codec
							}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Fps != nil {
								videoTemplateMap["fps"] = activity.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Fps
							}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Bitrate != nil {
								videoTemplateMap["bitrate"] = activity.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Bitrate
							}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.ResolutionAdaptive != nil {
								videoTemplateMap["resolution_adaptive"] = activity.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.ResolutionAdaptive
							}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Width != nil {
								videoTemplateMap["width"] = activity.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Width
							}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Height != nil {
								videoTemplateMap["height"] = activity.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Height
							}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Gop != nil {
								videoTemplateMap["gop"] = activity.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Gop
							}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.FillType != nil {
								videoTemplateMap["fill_type"] = activity.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.FillType
							}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Vcrf != nil {
								videoTemplateMap["vcrf"] = activity.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.Vcrf
							}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.ContentAdaptStream != nil {
								videoTemplateMap["content_adapt_stream"] = activity.ActivityPara.TranscodeTask.OverrideParameter.VideoTemplate.ContentAdaptStream
							}

							overrideParameterMap["video_template"] = []interface{}{videoTemplateMap}
						}

						if activity.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate != nil {
							audioTemplateMap := map[string]interface{}{}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.Codec != nil {
								audioTemplateMap["codec"] = activity.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.Codec
							}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.Bitrate != nil {
								audioTemplateMap["bitrate"] = activity.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.Bitrate
							}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.SampleRate != nil {
								audioTemplateMap["sample_rate"] = activity.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.SampleRate
							}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.AudioChannel != nil {
								audioTemplateMap["audio_channel"] = activity.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.AudioChannel
							}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.StreamSelects != nil {
								audioTemplateMap["stream_selects"] = activity.ActivityPara.TranscodeTask.OverrideParameter.AudioTemplate.StreamSelects
							}

							overrideParameterMap["audio_template"] = []interface{}{audioTemplateMap}
						}

						if activity.ActivityPara.TranscodeTask.OverrideParameter.TEHDConfig != nil {
							tEHDConfigMap := map[string]interface{}{}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.TEHDConfig.Type != nil {
								tEHDConfigMap["type"] = activity.ActivityPara.TranscodeTask.OverrideParameter.TEHDConfig.Type
							}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.TEHDConfig.MaxVideoBitrate != nil {
								tEHDConfigMap["max_video_bitrate"] = activity.ActivityPara.TranscodeTask.OverrideParameter.TEHDConfig.MaxVideoBitrate
							}

							overrideParameterMap["tehd_config"] = []interface{}{tEHDConfigMap}
						}

						if activity.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate != nil {
							subtitleTemplateMap := map[string]interface{}{}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.Path != nil {
								subtitleTemplateMap["path"] = activity.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.Path
							}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.StreamIndex != nil {
								subtitleTemplateMap["stream_index"] = activity.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.StreamIndex
							}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontType != nil {
								subtitleTemplateMap["font_type"] = activity.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontType
							}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontSize != nil {
								subtitleTemplateMap["font_size"] = activity.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontSize
							}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontColor != nil {
								subtitleTemplateMap["font_color"] = activity.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontColor
							}

							if activity.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontAlpha != nil {
								subtitleTemplateMap["font_alpha"] = activity.ActivityPara.TranscodeTask.OverrideParameter.SubtitleTemplate.FontAlpha
							}

							overrideParameterMap["subtitle_template"] = []interface{}{subtitleTemplateMap}
						}

						if activity.ActivityPara.TranscodeTask.OverrideParameter.AddonAudioStream != nil {
							addonAudioStreamList := []interface{}{}
							for _, addonAudioStream := range activity.ActivityPara.TranscodeTask.OverrideParameter.AddonAudioStream {
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

						if activity.ActivityPara.TranscodeTask.OverrideParameter.StdExtInfo != nil {
							overrideParameterMap["std_ext_info"] = activity.ActivityPara.TranscodeTask.OverrideParameter.StdExtInfo
						}

						if activity.ActivityPara.TranscodeTask.OverrideParameter.AddOnSubtitles != nil {
							addOnSubtitlesList := []interface{}{}
							for _, addOnSubtitles := range activity.ActivityPara.TranscodeTask.OverrideParameter.AddOnSubtitles {
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

					if activity.ActivityPara.TranscodeTask.WatermarkSet != nil {
						watermarkSetList := []interface{}{}
						for _, watermarkSet := range activity.ActivityPara.TranscodeTask.WatermarkSet {
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

					if activity.ActivityPara.TranscodeTask.MosaicSet != nil {
						mosaicSetList := []interface{}{}
						for _, mosaicSet := range activity.ActivityPara.TranscodeTask.MosaicSet {
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

					if activity.ActivityPara.TranscodeTask.StartTimeOffset != nil {
						transcodeTaskMap["start_time_offset"] = activity.ActivityPara.TranscodeTask.StartTimeOffset
					}

					if activity.ActivityPara.TranscodeTask.EndTimeOffset != nil {
						transcodeTaskMap["end_time_offset"] = activity.ActivityPara.TranscodeTask.EndTimeOffset
					}

					if activity.ActivityPara.TranscodeTask.OutputStorage != nil {
						outputStorageMap := map[string]interface{}{}

						if activity.ActivityPara.TranscodeTask.OutputStorage.Type != nil {
							outputStorageMap["type"] = activity.ActivityPara.TranscodeTask.OutputStorage.Type
						}

						if activity.ActivityPara.TranscodeTask.OutputStorage.CosOutputStorage != nil {
							cosOutputStorageMap := map[string]interface{}{}

							if activity.ActivityPara.TranscodeTask.OutputStorage.CosOutputStorage.Bucket != nil {
								cosOutputStorageMap["bucket"] = activity.ActivityPara.TranscodeTask.OutputStorage.CosOutputStorage.Bucket
							}

							if activity.ActivityPara.TranscodeTask.OutputStorage.CosOutputStorage.Region != nil {
								cosOutputStorageMap["region"] = activity.ActivityPara.TranscodeTask.OutputStorage.CosOutputStorage.Region
							}

							outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
						}

						if activity.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage != nil {
							s3OutputStorageMap := map[string]interface{}{}

							if activity.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3Bucket != nil {
								s3OutputStorageMap["s3_bucket"] = activity.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3Bucket
							}

							if activity.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3Region != nil {
								s3OutputStorageMap["s3_region"] = activity.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3Region
							}

							if activity.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3SecretId != nil {
								s3OutputStorageMap["s3_secret_id"] = activity.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3SecretId
							}

							if activity.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3SecretKey != nil {
								s3OutputStorageMap["s3_secret_key"] = activity.ActivityPara.TranscodeTask.OutputStorage.S3OutputStorage.S3SecretKey
							}

							outputStorageMap["s3_output_storage"] = []interface{}{s3OutputStorageMap}
						}

						transcodeTaskMap["output_storage"] = []interface{}{outputStorageMap}
					}

					if activity.ActivityPara.TranscodeTask.OutputObjectPath != nil {
						transcodeTaskMap["output_object_path"] = activity.ActivityPara.TranscodeTask.OutputObjectPath
					}

					if activity.ActivityPara.TranscodeTask.SegmentObjectName != nil {
						transcodeTaskMap["segment_object_name"] = activity.ActivityPara.TranscodeTask.SegmentObjectName
					}

					if activity.ActivityPara.TranscodeTask.ObjectNumberFormat != nil {
						objectNumberFormatMap := map[string]interface{}{}

						if activity.ActivityPara.TranscodeTask.ObjectNumberFormat.InitialValue != nil {
							objectNumberFormatMap["initial_value"] = activity.ActivityPara.TranscodeTask.ObjectNumberFormat.InitialValue
						}

						if activity.ActivityPara.TranscodeTask.ObjectNumberFormat.Increment != nil {
							objectNumberFormatMap["increment"] = activity.ActivityPara.TranscodeTask.ObjectNumberFormat.Increment
						}

						if activity.ActivityPara.TranscodeTask.ObjectNumberFormat.MinLength != nil {
							objectNumberFormatMap["min_length"] = activity.ActivityPara.TranscodeTask.ObjectNumberFormat.MinLength
						}

						if activity.ActivityPara.TranscodeTask.ObjectNumberFormat.PlaceHolder != nil {
							objectNumberFormatMap["place_holder"] = activity.ActivityPara.TranscodeTask.ObjectNumberFormat.PlaceHolder
						}

						transcodeTaskMap["object_number_format"] = []interface{}{objectNumberFormatMap}
					}

					if activity.ActivityPara.TranscodeTask.HeadTailParameter != nil {
						headTailParameterMap := map[string]interface{}{}

						if activity.ActivityPara.TranscodeTask.HeadTailParameter.HeadSet != nil {
							headSetList := []interface{}{}
							for _, headSet := range activity.ActivityPara.TranscodeTask.HeadTailParameter.HeadSet {
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

						if activity.ActivityPara.TranscodeTask.HeadTailParameter.TailSet != nil {
							tailSetList := []interface{}{}
							for _, tailSet := range activity.ActivityPara.TranscodeTask.HeadTailParameter.TailSet {
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

				if activity.ActivityPara.AnimatedGraphicTask != nil {
					animatedGraphicTaskMap := map[string]interface{}{}

					if activity.ActivityPara.AnimatedGraphicTask.Definition != nil {
						animatedGraphicTaskMap["definition"] = activity.ActivityPara.AnimatedGraphicTask.Definition
					}

					if activity.ActivityPara.AnimatedGraphicTask.StartTimeOffset != nil {
						animatedGraphicTaskMap["start_time_offset"] = activity.ActivityPara.AnimatedGraphicTask.StartTimeOffset
					}

					if activity.ActivityPara.AnimatedGraphicTask.EndTimeOffset != nil {
						animatedGraphicTaskMap["end_time_offset"] = activity.ActivityPara.AnimatedGraphicTask.EndTimeOffset
					}

					if activity.ActivityPara.AnimatedGraphicTask.OutputStorage != nil {
						outputStorageMap := map[string]interface{}{}

						if activity.ActivityPara.AnimatedGraphicTask.OutputStorage.Type != nil {
							outputStorageMap["type"] = activity.ActivityPara.AnimatedGraphicTask.OutputStorage.Type
						}

						if activity.ActivityPara.AnimatedGraphicTask.OutputStorage.CosOutputStorage != nil {
							cosOutputStorageMap := map[string]interface{}{}

							if activity.ActivityPara.AnimatedGraphicTask.OutputStorage.CosOutputStorage.Bucket != nil {
								cosOutputStorageMap["bucket"] = activity.ActivityPara.AnimatedGraphicTask.OutputStorage.CosOutputStorage.Bucket
							}

							if activity.ActivityPara.AnimatedGraphicTask.OutputStorage.CosOutputStorage.Region != nil {
								cosOutputStorageMap["region"] = activity.ActivityPara.AnimatedGraphicTask.OutputStorage.CosOutputStorage.Region
							}

							outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
						}

						if activity.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage != nil {
							s3OutputStorageMap := map[string]interface{}{}

							if activity.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3Bucket != nil {
								s3OutputStorageMap["s3_bucket"] = activity.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3Bucket
							}

							if activity.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3Region != nil {
								s3OutputStorageMap["s3_region"] = activity.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3Region
							}

							if activity.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3SecretId != nil {
								s3OutputStorageMap["s3_secret_id"] = activity.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3SecretId
							}

							if activity.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3SecretKey != nil {
								s3OutputStorageMap["s3_secret_key"] = activity.ActivityPara.AnimatedGraphicTask.OutputStorage.S3OutputStorage.S3SecretKey
							}

							outputStorageMap["s3_output_storage"] = []interface{}{s3OutputStorageMap}
						}

						animatedGraphicTaskMap["output_storage"] = []interface{}{outputStorageMap}
					}

					if activity.ActivityPara.AnimatedGraphicTask.OutputObjectPath != nil {
						animatedGraphicTaskMap["output_object_path"] = activity.ActivityPara.AnimatedGraphicTask.OutputObjectPath
					}

					activityParaMap["animated_graphic_task"] = []interface{}{animatedGraphicTaskMap}
				}

				if activity.ActivityPara.SnapshotByTimeOffsetTask != nil {
					snapshotByTimeOffsetTaskMap := map[string]interface{}{}

					if activity.ActivityPara.SnapshotByTimeOffsetTask.Definition != nil {
						snapshotByTimeOffsetTaskMap["definition"] = activity.ActivityPara.SnapshotByTimeOffsetTask.Definition
					}

					if activity.ActivityPara.SnapshotByTimeOffsetTask.ExtTimeOffsetSet != nil {
						snapshotByTimeOffsetTaskMap["ext_time_offset_set"] = activity.ActivityPara.SnapshotByTimeOffsetTask.ExtTimeOffsetSet
					}

					if activity.ActivityPara.SnapshotByTimeOffsetTask.WatermarkSet != nil {
						watermarkSetList := []interface{}{}
						for _, watermarkSet := range activity.ActivityPara.SnapshotByTimeOffsetTask.WatermarkSet {
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

					if activity.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage != nil {
						outputStorageMap := map[string]interface{}{}

						if activity.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.Type != nil {
							outputStorageMap["type"] = activity.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.Type
						}

						if activity.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.CosOutputStorage != nil {
							cosOutputStorageMap := map[string]interface{}{}

							if activity.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.CosOutputStorage.Bucket != nil {
								cosOutputStorageMap["bucket"] = activity.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.CosOutputStorage.Bucket
							}

							if activity.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.CosOutputStorage.Region != nil {
								cosOutputStorageMap["region"] = activity.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.CosOutputStorage.Region
							}

							outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
						}

						if activity.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage != nil {
							s3OutputStorageMap := map[string]interface{}{}

							if activity.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3Bucket != nil {
								s3OutputStorageMap["s3_bucket"] = activity.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3Bucket
							}

							if activity.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3Region != nil {
								s3OutputStorageMap["s3_region"] = activity.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3Region
							}

							if activity.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3SecretId != nil {
								s3OutputStorageMap["s3_secret_id"] = activity.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3SecretId
							}

							if activity.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3SecretKey != nil {
								s3OutputStorageMap["s3_secret_key"] = activity.ActivityPara.SnapshotByTimeOffsetTask.OutputStorage.S3OutputStorage.S3SecretKey
							}

							outputStorageMap["s3_output_storage"] = []interface{}{s3OutputStorageMap}
						}

						snapshotByTimeOffsetTaskMap["output_storage"] = []interface{}{outputStorageMap}
					}

					if activity.ActivityPara.SnapshotByTimeOffsetTask.OutputObjectPath != nil {
						snapshotByTimeOffsetTaskMap["output_object_path"] = activity.ActivityPara.SnapshotByTimeOffsetTask.OutputObjectPath
					}

					if activity.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat != nil {
						objectNumberFormatMap := map[string]interface{}{}

						if activity.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.InitialValue != nil {
							objectNumberFormatMap["initial_value"] = activity.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.InitialValue
						}

						if activity.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.Increment != nil {
							objectNumberFormatMap["increment"] = activity.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.Increment
						}

						if activity.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.MinLength != nil {
							objectNumberFormatMap["min_length"] = activity.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.MinLength
						}

						if activity.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.PlaceHolder != nil {
							objectNumberFormatMap["place_holder"] = activity.ActivityPara.SnapshotByTimeOffsetTask.ObjectNumberFormat.PlaceHolder
						}

						snapshotByTimeOffsetTaskMap["object_number_format"] = []interface{}{objectNumberFormatMap}
					}

					activityParaMap["snapshot_by_time_offset_task"] = []interface{}{snapshotByTimeOffsetTaskMap}
				}

				if activity.ActivityPara.SampleSnapshotTask != nil {
					sampleSnapshotTaskMap := map[string]interface{}{}

					if activity.ActivityPara.SampleSnapshotTask.Definition != nil {
						sampleSnapshotTaskMap["definition"] = activity.ActivityPara.SampleSnapshotTask.Definition
					}

					if activity.ActivityPara.SampleSnapshotTask.WatermarkSet != nil {
						watermarkSetList := []interface{}{}
						for _, watermarkSet := range activity.ActivityPara.SampleSnapshotTask.WatermarkSet {
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

					if activity.ActivityPara.SampleSnapshotTask.OutputStorage != nil {
						outputStorageMap := map[string]interface{}{}

						if activity.ActivityPara.SampleSnapshotTask.OutputStorage.Type != nil {
							outputStorageMap["type"] = activity.ActivityPara.SampleSnapshotTask.OutputStorage.Type
						}

						if activity.ActivityPara.SampleSnapshotTask.OutputStorage.CosOutputStorage != nil {
							cosOutputStorageMap := map[string]interface{}{}

							if activity.ActivityPara.SampleSnapshotTask.OutputStorage.CosOutputStorage.Bucket != nil {
								cosOutputStorageMap["bucket"] = activity.ActivityPara.SampleSnapshotTask.OutputStorage.CosOutputStorage.Bucket
							}

							if activity.ActivityPara.SampleSnapshotTask.OutputStorage.CosOutputStorage.Region != nil {
								cosOutputStorageMap["region"] = activity.ActivityPara.SampleSnapshotTask.OutputStorage.CosOutputStorage.Region
							}

							outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
						}

						if activity.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage != nil {
							s3OutputStorageMap := map[string]interface{}{}

							if activity.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3Bucket != nil {
								s3OutputStorageMap["s3_bucket"] = activity.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3Bucket
							}

							if activity.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3Region != nil {
								s3OutputStorageMap["s3_region"] = activity.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3Region
							}

							if activity.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3SecretId != nil {
								s3OutputStorageMap["s3_secret_id"] = activity.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3SecretId
							}

							if activity.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3SecretKey != nil {
								s3OutputStorageMap["s3_secret_key"] = activity.ActivityPara.SampleSnapshotTask.OutputStorage.S3OutputStorage.S3SecretKey
							}

							outputStorageMap["s3_output_storage"] = []interface{}{s3OutputStorageMap}
						}

						sampleSnapshotTaskMap["output_storage"] = []interface{}{outputStorageMap}
					}

					if activity.ActivityPara.SampleSnapshotTask.OutputObjectPath != nil {
						sampleSnapshotTaskMap["output_object_path"] = activity.ActivityPara.SampleSnapshotTask.OutputObjectPath
					}

					if activity.ActivityPara.SampleSnapshotTask.ObjectNumberFormat != nil {
						objectNumberFormatMap := map[string]interface{}{}

						if activity.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.InitialValue != nil {
							objectNumberFormatMap["initial_value"] = activity.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.InitialValue
						}

						if activity.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.Increment != nil {
							objectNumberFormatMap["increment"] = activity.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.Increment
						}

						if activity.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.MinLength != nil {
							objectNumberFormatMap["min_length"] = activity.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.MinLength
						}

						if activity.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.PlaceHolder != nil {
							objectNumberFormatMap["place_holder"] = activity.ActivityPara.SampleSnapshotTask.ObjectNumberFormat.PlaceHolder
						}

						sampleSnapshotTaskMap["object_number_format"] = []interface{}{objectNumberFormatMap}
					}

					activityParaMap["sample_snapshot_task"] = []interface{}{sampleSnapshotTaskMap}
				}

				if activity.ActivityPara.ImageSpriteTask != nil {
					imageSpriteTaskMap := map[string]interface{}{}

					if activity.ActivityPara.ImageSpriteTask.Definition != nil {
						imageSpriteTaskMap["definition"] = activity.ActivityPara.ImageSpriteTask.Definition
					}

					if activity.ActivityPara.ImageSpriteTask.OutputStorage != nil {
						outputStorageMap := map[string]interface{}{}

						if activity.ActivityPara.ImageSpriteTask.OutputStorage.Type != nil {
							outputStorageMap["type"] = activity.ActivityPara.ImageSpriteTask.OutputStorage.Type
						}

						if activity.ActivityPara.ImageSpriteTask.OutputStorage.CosOutputStorage != nil {
							cosOutputStorageMap := map[string]interface{}{}

							if activity.ActivityPara.ImageSpriteTask.OutputStorage.CosOutputStorage.Bucket != nil {
								cosOutputStorageMap["bucket"] = activity.ActivityPara.ImageSpriteTask.OutputStorage.CosOutputStorage.Bucket
							}

							if activity.ActivityPara.ImageSpriteTask.OutputStorage.CosOutputStorage.Region != nil {
								cosOutputStorageMap["region"] = activity.ActivityPara.ImageSpriteTask.OutputStorage.CosOutputStorage.Region
							}

							outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
						}

						if activity.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage != nil {
							s3OutputStorageMap := map[string]interface{}{}

							if activity.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3Bucket != nil {
								s3OutputStorageMap["s3_bucket"] = activity.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3Bucket
							}

							if activity.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3Region != nil {
								s3OutputStorageMap["s3_region"] = activity.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3Region
							}

							if activity.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3SecretId != nil {
								s3OutputStorageMap["s3_secret_id"] = activity.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3SecretId
							}

							if activity.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3SecretKey != nil {
								s3OutputStorageMap["s3_secret_key"] = activity.ActivityPara.ImageSpriteTask.OutputStorage.S3OutputStorage.S3SecretKey
							}

							outputStorageMap["s3_output_storage"] = []interface{}{s3OutputStorageMap}
						}

						imageSpriteTaskMap["output_storage"] = []interface{}{outputStorageMap}
					}

					if activity.ActivityPara.ImageSpriteTask.OutputObjectPath != nil {
						imageSpriteTaskMap["output_object_path"] = activity.ActivityPara.ImageSpriteTask.OutputObjectPath
					}

					if activity.ActivityPara.ImageSpriteTask.WebVttObjectName != nil {
						imageSpriteTaskMap["web_vtt_object_name"] = activity.ActivityPara.ImageSpriteTask.WebVttObjectName
					}

					if activity.ActivityPara.ImageSpriteTask.ObjectNumberFormat != nil {
						objectNumberFormatMap := map[string]interface{}{}

						if activity.ActivityPara.ImageSpriteTask.ObjectNumberFormat.InitialValue != nil {
							objectNumberFormatMap["initial_value"] = activity.ActivityPara.ImageSpriteTask.ObjectNumberFormat.InitialValue
						}

						if activity.ActivityPara.ImageSpriteTask.ObjectNumberFormat.Increment != nil {
							objectNumberFormatMap["increment"] = activity.ActivityPara.ImageSpriteTask.ObjectNumberFormat.Increment
						}

						if activity.ActivityPara.ImageSpriteTask.ObjectNumberFormat.MinLength != nil {
							objectNumberFormatMap["min_length"] = activity.ActivityPara.ImageSpriteTask.ObjectNumberFormat.MinLength
						}

						if activity.ActivityPara.ImageSpriteTask.ObjectNumberFormat.PlaceHolder != nil {
							objectNumberFormatMap["place_holder"] = activity.ActivityPara.ImageSpriteTask.ObjectNumberFormat.PlaceHolder
						}

						imageSpriteTaskMap["object_number_format"] = []interface{}{objectNumberFormatMap}
					}

					activityParaMap["image_sprite_task"] = []interface{}{imageSpriteTaskMap}
				}

				if activity.ActivityPara.AdaptiveDynamicStreamingTask != nil {
					adaptiveDynamicStreamingTaskMap := map[string]interface{}{}

					if activity.ActivityPara.AdaptiveDynamicStreamingTask.Definition != nil {
						adaptiveDynamicStreamingTaskMap["definition"] = activity.ActivityPara.AdaptiveDynamicStreamingTask.Definition
					}

					if activity.ActivityPara.AdaptiveDynamicStreamingTask.WatermarkSet != nil {
						watermarkSetList := []interface{}{}
						for _, watermarkSet := range activity.ActivityPara.AdaptiveDynamicStreamingTask.WatermarkSet {
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

					if activity.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage != nil {
						outputStorageMap := map[string]interface{}{}

						if activity.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.Type != nil {
							outputStorageMap["type"] = activity.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.Type
						}

						if activity.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.CosOutputStorage != nil {
							cosOutputStorageMap := map[string]interface{}{}

							if activity.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.CosOutputStorage.Bucket != nil {
								cosOutputStorageMap["bucket"] = activity.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.CosOutputStorage.Bucket
							}

							if activity.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.CosOutputStorage.Region != nil {
								cosOutputStorageMap["region"] = activity.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.CosOutputStorage.Region
							}

							outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
						}

						if activity.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage != nil {
							s3OutputStorageMap := map[string]interface{}{}

							if activity.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3Bucket != nil {
								s3OutputStorageMap["s3_bucket"] = activity.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3Bucket
							}

							if activity.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3Region != nil {
								s3OutputStorageMap["s3_region"] = activity.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3Region
							}

							if activity.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3SecretId != nil {
								s3OutputStorageMap["s3_secret_id"] = activity.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3SecretId
							}

							if activity.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3SecretKey != nil {
								s3OutputStorageMap["s3_secret_key"] = activity.ActivityPara.AdaptiveDynamicStreamingTask.OutputStorage.S3OutputStorage.S3SecretKey
							}

							outputStorageMap["s3_output_storage"] = []interface{}{s3OutputStorageMap}
						}

						adaptiveDynamicStreamingTaskMap["output_storage"] = []interface{}{outputStorageMap}
					}

					if activity.ActivityPara.AdaptiveDynamicStreamingTask.OutputObjectPath != nil {
						adaptiveDynamicStreamingTaskMap["output_object_path"] = activity.ActivityPara.AdaptiveDynamicStreamingTask.OutputObjectPath
					}

					if activity.ActivityPara.AdaptiveDynamicStreamingTask.SubStreamObjectName != nil {
						adaptiveDynamicStreamingTaskMap["sub_stream_object_name"] = activity.ActivityPara.AdaptiveDynamicStreamingTask.SubStreamObjectName
					}

					if activity.ActivityPara.AdaptiveDynamicStreamingTask.SegmentObjectName != nil {
						adaptiveDynamicStreamingTaskMap["segment_object_name"] = activity.ActivityPara.AdaptiveDynamicStreamingTask.SegmentObjectName
					}

					if activity.ActivityPara.AdaptiveDynamicStreamingTask.AddOnSubtitles != nil {
						addOnSubtitlesList := []interface{}{}
						for _, addOnSubtitles := range activity.ActivityPara.AdaptiveDynamicStreamingTask.AddOnSubtitles {
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

				if activity.ActivityPara.AiContentReviewTask != nil {
					aiContentReviewTaskMap := map[string]interface{}{}

					if activity.ActivityPara.AiContentReviewTask.Definition != nil {
						aiContentReviewTaskMap["definition"] = activity.ActivityPara.AiContentReviewTask.Definition
					}

					activityParaMap["ai_content_review_task"] = []interface{}{aiContentReviewTaskMap}
				}

				if activity.ActivityPara.AiAnalysisTask != nil {
					aiAnalysisTaskMap := map[string]interface{}{}

					if activity.ActivityPara.AiAnalysisTask.Definition != nil {
						aiAnalysisTaskMap["definition"] = activity.ActivityPara.AiAnalysisTask.Definition
					}

					if activity.ActivityPara.AiAnalysisTask.ExtendedParameter != nil {
						aiAnalysisTaskMap["extended_parameter"] = activity.ActivityPara.AiAnalysisTask.ExtendedParameter
					}

					activityParaMap["ai_analysis_task"] = []interface{}{aiAnalysisTaskMap}
				}

				if activity.ActivityPara.AiRecognitionTask != nil {
					aiRecognitionTaskMap := map[string]interface{}{}

					if activity.ActivityPara.AiRecognitionTask.Definition != nil {
						aiRecognitionTaskMap["definition"] = activity.ActivityPara.AiRecognitionTask.Definition
					}

					activityParaMap["ai_recognition_task"] = []interface{}{aiRecognitionTaskMap}
				}

				activitiesMap["activity_para"] = []interface{}{activityParaMap}
			}

			activitiesList = append(activitiesList, activitiesMap)
		}

		_ = d.Set("activities", activitiesList)

	}

	if schedule.OutputStorage != nil {
		outputStorageMap := map[string]interface{}{}

		if schedule.OutputStorage.Type != nil {
			outputStorageMap["type"] = schedule.OutputStorage.Type
		}

		if schedule.OutputStorage.CosOutputStorage != nil {
			cosOutputStorageMap := map[string]interface{}{}

			if schedule.OutputStorage.CosOutputStorage.Bucket != nil {
				cosOutputStorageMap["bucket"] = schedule.OutputStorage.CosOutputStorage.Bucket
			}

			if schedule.OutputStorage.CosOutputStorage.Region != nil {
				cosOutputStorageMap["region"] = schedule.OutputStorage.CosOutputStorage.Region
			}

			outputStorageMap["cos_output_storage"] = []interface{}{cosOutputStorageMap}
		}

		if schedule.OutputStorage.S3OutputStorage != nil {
			s3OutputStorageMap := map[string]interface{}{}

			if schedule.OutputStorage.S3OutputStorage.S3Bucket != nil {
				s3OutputStorageMap["s3_bucket"] = schedule.OutputStorage.S3OutputStorage.S3Bucket
			}

			if schedule.OutputStorage.S3OutputStorage.S3Region != nil {
				s3OutputStorageMap["s3_region"] = schedule.OutputStorage.S3OutputStorage.S3Region
			}

			if schedule.OutputStorage.S3OutputStorage.S3SecretId != nil {
				s3OutputStorageMap["s3_secret_id"] = schedule.OutputStorage.S3OutputStorage.S3SecretId
			}

			if schedule.OutputStorage.S3OutputStorage.S3SecretKey != nil {
				s3OutputStorageMap["s3_secret_key"] = schedule.OutputStorage.S3OutputStorage.S3SecretKey
			}

			outputStorageMap["s3_output_storage"] = []interface{}{s3OutputStorageMap}
		}

		_ = d.Set("output_storage", []interface{}{outputStorageMap})
	}

	if schedule.OutputDir != nil {
		_ = d.Set("output_dir", schedule.OutputDir)
	}

	if schedule.TaskNotifyConfig != nil {
		taskNotifyConfigMap := map[string]interface{}{}

		if schedule.TaskNotifyConfig.CmqModel != nil {
			taskNotifyConfigMap["cmq_model"] = schedule.TaskNotifyConfig.CmqModel
		}

		if schedule.TaskNotifyConfig.CmqRegion != nil {
			taskNotifyConfigMap["cmq_region"] = schedule.TaskNotifyConfig.CmqRegion
		}

		if schedule.TaskNotifyConfig.TopicName != nil {
			taskNotifyConfigMap["topic_name"] = schedule.TaskNotifyConfig.TopicName
		}

		if schedule.TaskNotifyConfig.QueueName != nil {
			taskNotifyConfigMap["queue_name"] = schedule.TaskNotifyConfig.QueueName
		}

		if schedule.TaskNotifyConfig.NotifyMode != nil {
			taskNotifyConfigMap["notify_mode"] = schedule.TaskNotifyConfig.NotifyMode
		}

		if schedule.TaskNotifyConfig.NotifyType != nil {
			taskNotifyConfigMap["notify_type"] = schedule.TaskNotifyConfig.NotifyType
		}

		if schedule.TaskNotifyConfig.NotifyUrl != nil {
			taskNotifyConfigMap["notify_url"] = schedule.TaskNotifyConfig.NotifyUrl
		}

		if schedule.TaskNotifyConfig.AwsSQS != nil {
			awsSQSMap := map[string]interface{}{}

			if schedule.TaskNotifyConfig.AwsSQS.SQSRegion != nil {
				awsSQSMap["sqs_region"] = schedule.TaskNotifyConfig.AwsSQS.SQSRegion
			}

			if schedule.TaskNotifyConfig.AwsSQS.SQSQueueName != nil {
				awsSQSMap["sqs_queue_name"] = schedule.TaskNotifyConfig.AwsSQS.SQSQueueName
			}

			if schedule.TaskNotifyConfig.AwsSQS.S3SecretId != nil {
				awsSQSMap["s3_secret_id"] = schedule.TaskNotifyConfig.AwsSQS.S3SecretId
			}

			if schedule.TaskNotifyConfig.AwsSQS.S3SecretKey != nil {
				awsSQSMap["s3_secret_key"] = schedule.TaskNotifyConfig.AwsSQS.S3SecretKey
			}

			taskNotifyConfigMap["aws_sqs"] = []interface{}{awsSQSMap}
		}

		_ = d.Set("task_notify_config", []interface{}{taskNotifyConfigMap})
	}

	if schedule.ResourceId != nil {
		_ = d.Set("resource_id", schedule.ResourceId)
	}

	return nil
}

func resourceTencentCloudMpsScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_schedule.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := mps.NewModifyScheduleRequest()

	scheduleId := d.Id()

	request.ScheduleId = helper.StrToInt64Point(scheduleId)

	if d.HasChange("schedule_name") {
		if v, ok := d.GetOk("schedule_name"); ok {
			request.ScheduleName = helper.String(v.(string))
		}
	}

	if d.HasChange("trigger") {
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
						if formatsSet[i] != nil {
							formats := formatsSet[i].(string)
							cosFileUploadTrigger.Formats = append(cosFileUploadTrigger.Formats, &formats)
						}
					}
				}
				workflowTrigger.CosFileUploadTrigger = &cosFileUploadTrigger
			}
			if awsS3FileUploadTriggerMap, ok := helper.InterfaceToMap(dMap, "aws_s3_file_upload_trigger"); ok {
				awsS3FileUploadTrigger := mps.AwsS3FileUploadTrigger{}
				if v, ok := awsS3FileUploadTriggerMap["s3_bucket"]; ok {
					awsS3FileUploadTrigger.S3Bucket = helper.String(v.(string))
				}
				if v, ok := awsS3FileUploadTriggerMap["s3_region"]; ok {
					awsS3FileUploadTrigger.S3Region = helper.String(v.(string))
				}
				if v, ok := awsS3FileUploadTriggerMap["dir"]; ok {
					awsS3FileUploadTrigger.Dir = helper.String(v.(string))
				}
				if v, ok := awsS3FileUploadTriggerMap["formats"]; ok {
					formatsSet := v.(*schema.Set).List()
					for i := range formatsSet {
						if formatsSet[i] != nil {
							formats := formatsSet[i].(string)
							awsS3FileUploadTrigger.Formats = append(awsS3FileUploadTrigger.Formats, &formats)
						}
					}
				}
				if v, ok := awsS3FileUploadTriggerMap["s3_secret_id"]; ok {
					awsS3FileUploadTrigger.S3SecretId = helper.String(v.(string))
				}
				if v, ok := awsS3FileUploadTriggerMap["s3_secret_key"]; ok {
					awsS3FileUploadTrigger.S3SecretKey = helper.String(v.(string))
				}
				if awsSQSMap, ok := helper.InterfaceToMap(awsS3FileUploadTriggerMap, "aws_sqs"); ok {
					awsSQS := mps.AwsSQS{}
					if v, ok := awsSQSMap["sqs_region"]; ok {
						awsSQS.SQSRegion = helper.String(v.(string))
					}
					if v, ok := awsSQSMap["sqs_queue_name"]; ok {
						awsSQS.SQSQueueName = helper.String(v.(string))
					}
					if v, ok := awsSQSMap["s3_secret_id"]; ok {
						awsSQS.S3SecretId = helper.String(v.(string))
					}
					if v, ok := awsSQSMap["s3_secret_key"]; ok {
						awsSQS.S3SecretKey = helper.String(v.(string))
					}
					awsS3FileUploadTrigger.AwsSQS = &awsSQS
				}
				workflowTrigger.AwsS3FileUploadTrigger = &awsS3FileUploadTrigger
			}
			request.Trigger = &workflowTrigger
		}
	}

	if d.HasChange("activities") {
		if v, ok := d.GetOk("activities"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				activity := mps.Activity{}
				if v, ok := dMap["activity_type"]; ok {
					activity.ActivityType = helper.String(v.(string))
				}
				if v, ok := dMap["reardrive_index"]; ok {
					reardriveIndexSet := v.(*schema.Set).List()
					for i := range reardriveIndexSet {
						reardriveIndex := reardriveIndexSet[i].(int)
						activity.ReardriveIndex = append(activity.ReardriveIndex, helper.IntInt64(reardriveIndex))
					}
				}
				if activityParaMap, ok := helper.InterfaceToMap(dMap, "activity_para"); ok {
					activityPara := mps.ActivityPara{}
					if transcodeTaskMap, ok := helper.InterfaceToMap(activityParaMap, "transcode_task"); ok {
						transcodeTaskInput := mps.TranscodeTaskInput{}
						if v, ok := transcodeTaskMap["definition"]; ok {
							transcodeTaskInput.Definition = helper.IntUint64(v.(int))
						}
						if rawParameterMap, ok := helper.InterfaceToMap(transcodeTaskMap, "raw_parameter"); ok {
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
						if overrideParameterMap, ok := helper.InterfaceToMap(transcodeTaskMap, "override_parameter"); ok {
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
							if v, ok := overrideParameterMap["addon_audio_stream"]; ok {
								for _, item := range v.([]interface{}) {
									addonAudioStreamMap := item.(map[string]interface{})
									mediaInputInfo := mps.MediaInputInfo{}
									if v, ok := addonAudioStreamMap["type"]; ok {
										mediaInputInfo.Type = helper.String(v.(string))
									}
									if cosInputInfoMap, ok := helper.InterfaceToMap(addonAudioStreamMap, "cos_input_info"); ok {
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
									if urlInputInfoMap, ok := helper.InterfaceToMap(addonAudioStreamMap, "url_input_info"); ok {
										urlInputInfo := mps.UrlInputInfo{}
										if v, ok := urlInputInfoMap["url"]; ok {
											urlInputInfo.Url = helper.String(v.(string))
										}
										mediaInputInfo.UrlInputInfo = &urlInputInfo
									}
									if s3InputInfoMap, ok := helper.InterfaceToMap(addonAudioStreamMap, "s3_input_info"); ok {
										s3InputInfo := mps.S3InputInfo{}
										if v, ok := s3InputInfoMap["s3_bucket"]; ok {
											s3InputInfo.S3Bucket = helper.String(v.(string))
										}
										if v, ok := s3InputInfoMap["s3_region"]; ok {
											s3InputInfo.S3Region = helper.String(v.(string))
										}
										if v, ok := s3InputInfoMap["s3_object"]; ok {
											s3InputInfo.S3Object = helper.String(v.(string))
										}
										if v, ok := s3InputInfoMap["s3_secret_id"]; ok {
											s3InputInfo.S3SecretId = helper.String(v.(string))
										}
										if v, ok := s3InputInfoMap["s3_secret_key"]; ok {
											s3InputInfo.S3SecretKey = helper.String(v.(string))
										}
										mediaInputInfo.S3InputInfo = &s3InputInfo
									}
									overrideTranscodeParameter.AddonAudioStream = append(overrideTranscodeParameter.AddonAudioStream, &mediaInputInfo)
								}
							}
							if v, ok := overrideParameterMap["std_ext_info"]; ok {
								overrideTranscodeParameter.StdExtInfo = helper.String(v.(string))
							}
							if v, ok := overrideParameterMap["add_on_subtitles"]; ok {
								for _, item := range v.([]interface{}) {
									addOnSubtitlesMap := item.(map[string]interface{})
									addOnSubtitle := mps.AddOnSubtitle{}
									if v, ok := addOnSubtitlesMap["type"]; ok {
										addOnSubtitle.Type = helper.String(v.(string))
									}
									if subtitleMap, ok := helper.InterfaceToMap(addOnSubtitlesMap, "subtitle"); ok {
										mediaInputInfo := mps.MediaInputInfo{}
										if v, ok := subtitleMap["type"]; ok {
											mediaInputInfo.Type = helper.String(v.(string))
										}
										if cosInputInfoMap, ok := helper.InterfaceToMap(subtitleMap, "cos_input_info"); ok {
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
										if urlInputInfoMap, ok := helper.InterfaceToMap(subtitleMap, "url_input_info"); ok {
											urlInputInfo := mps.UrlInputInfo{}
											if v, ok := urlInputInfoMap["url"]; ok {
												urlInputInfo.Url = helper.String(v.(string))
											}
											mediaInputInfo.UrlInputInfo = &urlInputInfo
										}
										if s3InputInfoMap, ok := helper.InterfaceToMap(subtitleMap, "s3_input_info"); ok {
											s3InputInfo := mps.S3InputInfo{}
											if v, ok := s3InputInfoMap["s3_bucket"]; ok {
												s3InputInfo.S3Bucket = helper.String(v.(string))
											}
											if v, ok := s3InputInfoMap["s3_region"]; ok {
												s3InputInfo.S3Region = helper.String(v.(string))
											}
											if v, ok := s3InputInfoMap["s3_object"]; ok {
												s3InputInfo.S3Object = helper.String(v.(string))
											}
											if v, ok := s3InputInfoMap["s3_secret_id"]; ok {
												s3InputInfo.S3SecretId = helper.String(v.(string))
											}
											if v, ok := s3InputInfoMap["s3_secret_key"]; ok {
												s3InputInfo.S3SecretKey = helper.String(v.(string))
											}
											mediaInputInfo.S3InputInfo = &s3InputInfo
										}
										addOnSubtitle.Subtitle = &mediaInputInfo
									}
									overrideTranscodeParameter.AddOnSubtitles = append(overrideTranscodeParameter.AddOnSubtitles, &addOnSubtitle)
								}
							}
							transcodeTaskInput.OverrideParameter = &overrideTranscodeParameter
						}
						if v, ok := transcodeTaskMap["watermark_set"]; ok {
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
											if s3InputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "s3_input_info"); ok {
												s3InputInfo := mps.S3InputInfo{}
												if v, ok := s3InputInfoMap["s3_bucket"]; ok {
													s3InputInfo.S3Bucket = helper.String(v.(string))
												}
												if v, ok := s3InputInfoMap["s3_region"]; ok {
													s3InputInfo.S3Region = helper.String(v.(string))
												}
												if v, ok := s3InputInfoMap["s3_object"]; ok {
													s3InputInfo.S3Object = helper.String(v.(string))
												}
												if v, ok := s3InputInfoMap["s3_secret_id"]; ok {
													s3InputInfo.S3SecretId = helper.String(v.(string))
												}
												if v, ok := s3InputInfoMap["s3_secret_key"]; ok {
													s3InputInfo.S3SecretKey = helper.String(v.(string))
												}
												mediaInputInfo.S3InputInfo = &s3InputInfo
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
						if v, ok := transcodeTaskMap["mosaic_set"]; ok {
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
						if v, ok := transcodeTaskMap["start_time_offset"]; ok {
							transcodeTaskInput.StartTimeOffset = helper.Float64(v.(float64))
						}
						if v, ok := transcodeTaskMap["end_time_offset"]; ok {
							transcodeTaskInput.EndTimeOffset = helper.Float64(v.(float64))
						}
						if outputStorageMap, ok := helper.InterfaceToMap(transcodeTaskMap, "output_storage"); ok {
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
							if s3OutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "s3_output_storage"); ok {
								s3OutputStorage := mps.S3OutputStorage{}
								if v, ok := s3OutputStorageMap["s3_bucket"]; ok {
									s3OutputStorage.S3Bucket = helper.String(v.(string))
								}
								if v, ok := s3OutputStorageMap["s3_region"]; ok {
									s3OutputStorage.S3Region = helper.String(v.(string))
								}
								if v, ok := s3OutputStorageMap["s3_secret_id"]; ok {
									s3OutputStorage.S3SecretId = helper.String(v.(string))
								}
								if v, ok := s3OutputStorageMap["s3_secret_key"]; ok {
									s3OutputStorage.S3SecretKey = helper.String(v.(string))
								}
								taskOutputStorage.S3OutputStorage = &s3OutputStorage
							}
							transcodeTaskInput.OutputStorage = &taskOutputStorage
						}
						if v, ok := transcodeTaskMap["output_object_path"]; ok {
							transcodeTaskInput.OutputObjectPath = helper.String(v.(string))
						}
						if v, ok := transcodeTaskMap["segment_object_name"]; ok {
							transcodeTaskInput.SegmentObjectName = helper.String(v.(string))
						}
						if objectNumberFormatMap, ok := helper.InterfaceToMap(transcodeTaskMap, "object_number_format"); ok {
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
						if headTailParameterMap, ok := helper.InterfaceToMap(transcodeTaskMap, "head_tail_parameter"); ok {
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
									if s3InputInfoMap, ok := helper.InterfaceToMap(headSetMap, "s3_input_info"); ok {
										s3InputInfo := mps.S3InputInfo{}
										if v, ok := s3InputInfoMap["s3_bucket"]; ok {
											s3InputInfo.S3Bucket = helper.String(v.(string))
										}
										if v, ok := s3InputInfoMap["s3_region"]; ok {
											s3InputInfo.S3Region = helper.String(v.(string))
										}
										if v, ok := s3InputInfoMap["s3_object"]; ok {
											s3InputInfo.S3Object = helper.String(v.(string))
										}
										if v, ok := s3InputInfoMap["s3_secret_id"]; ok {
											s3InputInfo.S3SecretId = helper.String(v.(string))
										}
										if v, ok := s3InputInfoMap["s3_secret_key"]; ok {
											s3InputInfo.S3SecretKey = helper.String(v.(string))
										}
										mediaInputInfo.S3InputInfo = &s3InputInfo
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
									if s3InputInfoMap, ok := helper.InterfaceToMap(tailSetMap, "s3_input_info"); ok {
										s3InputInfo := mps.S3InputInfo{}
										if v, ok := s3InputInfoMap["s3_bucket"]; ok {
											s3InputInfo.S3Bucket = helper.String(v.(string))
										}
										if v, ok := s3InputInfoMap["s3_region"]; ok {
											s3InputInfo.S3Region = helper.String(v.(string))
										}
										if v, ok := s3InputInfoMap["s3_object"]; ok {
											s3InputInfo.S3Object = helper.String(v.(string))
										}
										if v, ok := s3InputInfoMap["s3_secret_id"]; ok {
											s3InputInfo.S3SecretId = helper.String(v.(string))
										}
										if v, ok := s3InputInfoMap["s3_secret_key"]; ok {
											s3InputInfo.S3SecretKey = helper.String(v.(string))
										}
										mediaInputInfo.S3InputInfo = &s3InputInfo
									}
									headTailParameter.TailSet = append(headTailParameter.TailSet, &mediaInputInfo)
								}
							}
							transcodeTaskInput.HeadTailParameter = &headTailParameter
						}
						activityPara.TranscodeTask = &transcodeTaskInput
					}
					if animatedGraphicTaskMap, ok := helper.InterfaceToMap(activityParaMap, "animated_graphic_task"); ok {
						animatedGraphicTaskInput := mps.AnimatedGraphicTaskInput{}
						if v, ok := animatedGraphicTaskMap["definition"]; ok {
							animatedGraphicTaskInput.Definition = helper.IntUint64(v.(int))
						}
						if v, ok := animatedGraphicTaskMap["start_time_offset"]; ok {
							animatedGraphicTaskInput.StartTimeOffset = helper.Float64(v.(float64))
						}
						if v, ok := animatedGraphicTaskMap["end_time_offset"]; ok {
							animatedGraphicTaskInput.EndTimeOffset = helper.Float64(v.(float64))
						}
						if outputStorageMap, ok := helper.InterfaceToMap(animatedGraphicTaskMap, "output_storage"); ok {
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
							if s3OutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "s3_output_storage"); ok {
								s3OutputStorage := mps.S3OutputStorage{}
								if v, ok := s3OutputStorageMap["s3_bucket"]; ok {
									s3OutputStorage.S3Bucket = helper.String(v.(string))
								}
								if v, ok := s3OutputStorageMap["s3_region"]; ok {
									s3OutputStorage.S3Region = helper.String(v.(string))
								}
								if v, ok := s3OutputStorageMap["s3_secret_id"]; ok {
									s3OutputStorage.S3SecretId = helper.String(v.(string))
								}
								if v, ok := s3OutputStorageMap["s3_secret_key"]; ok {
									s3OutputStorage.S3SecretKey = helper.String(v.(string))
								}
								taskOutputStorage.S3OutputStorage = &s3OutputStorage
							}
							animatedGraphicTaskInput.OutputStorage = &taskOutputStorage
						}
						if v, ok := animatedGraphicTaskMap["output_object_path"]; ok {
							animatedGraphicTaskInput.OutputObjectPath = helper.String(v.(string))
						}
						activityPara.AnimatedGraphicTask = &animatedGraphicTaskInput
					}
					if snapshotByTimeOffsetTaskMap, ok := helper.InterfaceToMap(activityParaMap, "snapshot_by_time_offset_task"); ok {
						snapshotByTimeOffsetTaskInput := mps.SnapshotByTimeOffsetTaskInput{}
						if v, ok := snapshotByTimeOffsetTaskMap["definition"]; ok {
							snapshotByTimeOffsetTaskInput.Definition = helper.IntUint64(v.(int))
						}
						if v, ok := snapshotByTimeOffsetTaskMap["ext_time_offset_set"]; ok {
							extTimeOffsetSetSet := v.(*schema.Set).List()
							for i := range extTimeOffsetSetSet {
								if extTimeOffsetSetSet[i] != nil {
									extTimeOffsetSet := extTimeOffsetSetSet[i].(string)
									snapshotByTimeOffsetTaskInput.ExtTimeOffsetSet = append(snapshotByTimeOffsetTaskInput.ExtTimeOffsetSet, &extTimeOffsetSet)
								}
							}
						}

						if v, ok := snapshotByTimeOffsetTaskMap["watermark_set"]; ok {
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
											if s3InputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "s3_input_info"); ok {
												s3InputInfo := mps.S3InputInfo{}
												if v, ok := s3InputInfoMap["s3_bucket"]; ok {
													s3InputInfo.S3Bucket = helper.String(v.(string))
												}
												if v, ok := s3InputInfoMap["s3_region"]; ok {
													s3InputInfo.S3Region = helper.String(v.(string))
												}
												if v, ok := s3InputInfoMap["s3_object"]; ok {
													s3InputInfo.S3Object = helper.String(v.(string))
												}
												if v, ok := s3InputInfoMap["s3_secret_id"]; ok {
													s3InputInfo.S3SecretId = helper.String(v.(string))
												}
												if v, ok := s3InputInfoMap["s3_secret_key"]; ok {
													s3InputInfo.S3SecretKey = helper.String(v.(string))
												}
												mediaInputInfo.S3InputInfo = &s3InputInfo
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
						if outputStorageMap, ok := helper.InterfaceToMap(snapshotByTimeOffsetTaskMap, "output_storage"); ok {
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
							if s3OutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "s3_output_storage"); ok {
								s3OutputStorage := mps.S3OutputStorage{}
								if v, ok := s3OutputStorageMap["s3_bucket"]; ok {
									s3OutputStorage.S3Bucket = helper.String(v.(string))
								}
								if v, ok := s3OutputStorageMap["s3_region"]; ok {
									s3OutputStorage.S3Region = helper.String(v.(string))
								}
								if v, ok := s3OutputStorageMap["s3_secret_id"]; ok {
									s3OutputStorage.S3SecretId = helper.String(v.(string))
								}
								if v, ok := s3OutputStorageMap["s3_secret_key"]; ok {
									s3OutputStorage.S3SecretKey = helper.String(v.(string))
								}
								taskOutputStorage.S3OutputStorage = &s3OutputStorage
							}
							snapshotByTimeOffsetTaskInput.OutputStorage = &taskOutputStorage
						}
						if v, ok := snapshotByTimeOffsetTaskMap["output_object_path"]; ok {
							snapshotByTimeOffsetTaskInput.OutputObjectPath = helper.String(v.(string))
						}
						if objectNumberFormatMap, ok := helper.InterfaceToMap(snapshotByTimeOffsetTaskMap, "object_number_format"); ok {
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
						activityPara.SnapshotByTimeOffsetTask = &snapshotByTimeOffsetTaskInput
					}
					if sampleSnapshotTaskMap, ok := helper.InterfaceToMap(activityParaMap, "sample_snapshot_task"); ok {
						sampleSnapshotTaskInput := mps.SampleSnapshotTaskInput{}
						if v, ok := sampleSnapshotTaskMap["definition"]; ok {
							sampleSnapshotTaskInput.Definition = helper.IntUint64(v.(int))
						}
						if v, ok := sampleSnapshotTaskMap["watermark_set"]; ok {
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
											if s3InputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "s3_input_info"); ok {
												s3InputInfo := mps.S3InputInfo{}
												if v, ok := s3InputInfoMap["s3_bucket"]; ok {
													s3InputInfo.S3Bucket = helper.String(v.(string))
												}
												if v, ok := s3InputInfoMap["s3_region"]; ok {
													s3InputInfo.S3Region = helper.String(v.(string))
												}
												if v, ok := s3InputInfoMap["s3_object"]; ok {
													s3InputInfo.S3Object = helper.String(v.(string))
												}
												if v, ok := s3InputInfoMap["s3_secret_id"]; ok {
													s3InputInfo.S3SecretId = helper.String(v.(string))
												}
												if v, ok := s3InputInfoMap["s3_secret_key"]; ok {
													s3InputInfo.S3SecretKey = helper.String(v.(string))
												}
												mediaInputInfo.S3InputInfo = &s3InputInfo
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
						if outputStorageMap, ok := helper.InterfaceToMap(sampleSnapshotTaskMap, "output_storage"); ok {
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
							if s3OutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "s3_output_storage"); ok {
								s3OutputStorage := mps.S3OutputStorage{}
								if v, ok := s3OutputStorageMap["s3_bucket"]; ok {
									s3OutputStorage.S3Bucket = helper.String(v.(string))
								}
								if v, ok := s3OutputStorageMap["s3_region"]; ok {
									s3OutputStorage.S3Region = helper.String(v.(string))
								}
								if v, ok := s3OutputStorageMap["s3_secret_id"]; ok {
									s3OutputStorage.S3SecretId = helper.String(v.(string))
								}
								if v, ok := s3OutputStorageMap["s3_secret_key"]; ok {
									s3OutputStorage.S3SecretKey = helper.String(v.(string))
								}
								taskOutputStorage.S3OutputStorage = &s3OutputStorage
							}
							sampleSnapshotTaskInput.OutputStorage = &taskOutputStorage
						}
						if v, ok := sampleSnapshotTaskMap["output_object_path"]; ok {
							sampleSnapshotTaskInput.OutputObjectPath = helper.String(v.(string))
						}
						if objectNumberFormatMap, ok := helper.InterfaceToMap(sampleSnapshotTaskMap, "object_number_format"); ok {
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
						activityPara.SampleSnapshotTask = &sampleSnapshotTaskInput
					}
					if imageSpriteTaskMap, ok := helper.InterfaceToMap(activityParaMap, "image_sprite_task"); ok {
						imageSpriteTaskInput := mps.ImageSpriteTaskInput{}
						if v, ok := imageSpriteTaskMap["definition"]; ok {
							imageSpriteTaskInput.Definition = helper.IntUint64(v.(int))
						}
						if outputStorageMap, ok := helper.InterfaceToMap(imageSpriteTaskMap, "output_storage"); ok {
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
							if s3OutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "s3_output_storage"); ok {
								s3OutputStorage := mps.S3OutputStorage{}
								if v, ok := s3OutputStorageMap["s3_bucket"]; ok {
									s3OutputStorage.S3Bucket = helper.String(v.(string))
								}
								if v, ok := s3OutputStorageMap["s3_region"]; ok {
									s3OutputStorage.S3Region = helper.String(v.(string))
								}
								if v, ok := s3OutputStorageMap["s3_secret_id"]; ok {
									s3OutputStorage.S3SecretId = helper.String(v.(string))
								}
								if v, ok := s3OutputStorageMap["s3_secret_key"]; ok {
									s3OutputStorage.S3SecretKey = helper.String(v.(string))
								}
								taskOutputStorage.S3OutputStorage = &s3OutputStorage
							}
							imageSpriteTaskInput.OutputStorage = &taskOutputStorage
						}
						if v, ok := imageSpriteTaskMap["output_object_path"]; ok {
							imageSpriteTaskInput.OutputObjectPath = helper.String(v.(string))
						}
						if v, ok := imageSpriteTaskMap["web_vtt_object_name"]; ok {
							imageSpriteTaskInput.WebVttObjectName = helper.String(v.(string))
						}
						if objectNumberFormatMap, ok := helper.InterfaceToMap(imageSpriteTaskMap, "object_number_format"); ok {
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
						activityPara.ImageSpriteTask = &imageSpriteTaskInput
					}
					if adaptiveDynamicStreamingTaskMap, ok := helper.InterfaceToMap(activityParaMap, "adaptive_dynamic_streaming_task"); ok {
						adaptiveDynamicStreamingTaskInput := mps.AdaptiveDynamicStreamingTaskInput{}
						if v, ok := adaptiveDynamicStreamingTaskMap["definition"]; ok {
							adaptiveDynamicStreamingTaskInput.Definition = helper.IntUint64(v.(int))
						}
						if v, ok := adaptiveDynamicStreamingTaskMap["watermark_set"]; ok {
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
											if s3InputInfoMap, ok := helper.InterfaceToMap(imageContentMap, "s3_input_info"); ok {
												s3InputInfo := mps.S3InputInfo{}
												if v, ok := s3InputInfoMap["s3_bucket"]; ok {
													s3InputInfo.S3Bucket = helper.String(v.(string))
												}
												if v, ok := s3InputInfoMap["s3_region"]; ok {
													s3InputInfo.S3Region = helper.String(v.(string))
												}
												if v, ok := s3InputInfoMap["s3_object"]; ok {
													s3InputInfo.S3Object = helper.String(v.(string))
												}
												if v, ok := s3InputInfoMap["s3_secret_id"]; ok {
													s3InputInfo.S3SecretId = helper.String(v.(string))
												}
												if v, ok := s3InputInfoMap["s3_secret_key"]; ok {
													s3InputInfo.S3SecretKey = helper.String(v.(string))
												}
												mediaInputInfo.S3InputInfo = &s3InputInfo
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
						if outputStorageMap, ok := helper.InterfaceToMap(adaptiveDynamicStreamingTaskMap, "output_storage"); ok {
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
							if s3OutputStorageMap, ok := helper.InterfaceToMap(outputStorageMap, "s3_output_storage"); ok {
								s3OutputStorage := mps.S3OutputStorage{}
								if v, ok := s3OutputStorageMap["s3_bucket"]; ok {
									s3OutputStorage.S3Bucket = helper.String(v.(string))
								}
								if v, ok := s3OutputStorageMap["s3_region"]; ok {
									s3OutputStorage.S3Region = helper.String(v.(string))
								}
								if v, ok := s3OutputStorageMap["s3_secret_id"]; ok {
									s3OutputStorage.S3SecretId = helper.String(v.(string))
								}
								if v, ok := s3OutputStorageMap["s3_secret_key"]; ok {
									s3OutputStorage.S3SecretKey = helper.String(v.(string))
								}
								taskOutputStorage.S3OutputStorage = &s3OutputStorage
							}
							adaptiveDynamicStreamingTaskInput.OutputStorage = &taskOutputStorage
						}
						if v, ok := adaptiveDynamicStreamingTaskMap["output_object_path"]; ok {
							adaptiveDynamicStreamingTaskInput.OutputObjectPath = helper.String(v.(string))
						}
						if v, ok := adaptiveDynamicStreamingTaskMap["sub_stream_object_name"]; ok {
							adaptiveDynamicStreamingTaskInput.SubStreamObjectName = helper.String(v.(string))
						}
						if v, ok := adaptiveDynamicStreamingTaskMap["segment_object_name"]; ok {
							adaptiveDynamicStreamingTaskInput.SegmentObjectName = helper.String(v.(string))
						}
						if v, ok := adaptiveDynamicStreamingTaskMap["add_on_subtitles"]; ok {
							for _, item := range v.([]interface{}) {
								addOnSubtitlesMap := item.(map[string]interface{})
								addOnSubtitle := mps.AddOnSubtitle{}
								if v, ok := addOnSubtitlesMap["type"]; ok {
									addOnSubtitle.Type = helper.String(v.(string))
								}
								if subtitleMap, ok := helper.InterfaceToMap(addOnSubtitlesMap, "subtitle"); ok {
									mediaInputInfo := mps.MediaInputInfo{}
									if v, ok := subtitleMap["type"]; ok {
										mediaInputInfo.Type = helper.String(v.(string))
									}
									if cosInputInfoMap, ok := helper.InterfaceToMap(subtitleMap, "cos_input_info"); ok {
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
									if urlInputInfoMap, ok := helper.InterfaceToMap(subtitleMap, "url_input_info"); ok {
										urlInputInfo := mps.UrlInputInfo{}
										if v, ok := urlInputInfoMap["url"]; ok {
											urlInputInfo.Url = helper.String(v.(string))
										}
										mediaInputInfo.UrlInputInfo = &urlInputInfo
									}
									if s3InputInfoMap, ok := helper.InterfaceToMap(subtitleMap, "s3_input_info"); ok {
										s3InputInfo := mps.S3InputInfo{}
										if v, ok := s3InputInfoMap["s3_bucket"]; ok {
											s3InputInfo.S3Bucket = helper.String(v.(string))
										}
										if v, ok := s3InputInfoMap["s3_region"]; ok {
											s3InputInfo.S3Region = helper.String(v.(string))
										}
										if v, ok := s3InputInfoMap["s3_object"]; ok {
											s3InputInfo.S3Object = helper.String(v.(string))
										}
										if v, ok := s3InputInfoMap["s3_secret_id"]; ok {
											s3InputInfo.S3SecretId = helper.String(v.(string))
										}
										if v, ok := s3InputInfoMap["s3_secret_key"]; ok {
											s3InputInfo.S3SecretKey = helper.String(v.(string))
										}
										mediaInputInfo.S3InputInfo = &s3InputInfo
									}
									addOnSubtitle.Subtitle = &mediaInputInfo
								}
								adaptiveDynamicStreamingTaskInput.AddOnSubtitles = append(adaptiveDynamicStreamingTaskInput.AddOnSubtitles, &addOnSubtitle)
							}
						}
						activityPara.AdaptiveDynamicStreamingTask = &adaptiveDynamicStreamingTaskInput
					}
					if aiContentReviewTaskMap, ok := helper.InterfaceToMap(activityParaMap, "ai_content_review_task"); ok {
						aiContentReviewTaskInput := mps.AiContentReviewTaskInput{}
						if v, ok := aiContentReviewTaskMap["definition"]; ok {
							aiContentReviewTaskInput.Definition = helper.IntUint64(v.(int))
						}
						activityPara.AiContentReviewTask = &aiContentReviewTaskInput
					}
					if aiAnalysisTaskMap, ok := helper.InterfaceToMap(activityParaMap, "ai_analysis_task"); ok {
						aiAnalysisTaskInput := mps.AiAnalysisTaskInput{}
						if v, ok := aiAnalysisTaskMap["definition"]; ok {
							aiAnalysisTaskInput.Definition = helper.IntUint64(v.(int))
						}
						if v, ok := aiAnalysisTaskMap["extended_parameter"]; ok {
							aiAnalysisTaskInput.ExtendedParameter = helper.String(v.(string))
						}
						activityPara.AiAnalysisTask = &aiAnalysisTaskInput
					}
					if aiRecognitionTaskMap, ok := helper.InterfaceToMap(activityParaMap, "ai_recognition_task"); ok {
						aiRecognitionTaskInput := mps.AiRecognitionTaskInput{}
						if v, ok := aiRecognitionTaskMap["definition"]; ok {
							aiRecognitionTaskInput.Definition = helper.IntUint64(v.(int))
						}
						activityPara.AiRecognitionTask = &aiRecognitionTaskInput
					}
					activity.ActivityPara = &activityPara
				}
				request.Activities = append(request.Activities, &activity)
			}
		}
	}

	if d.HasChange("output_storage") {
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
			if s3OutputStorageMap, ok := helper.InterfaceToMap(dMap, "s3_output_storage"); ok {
				s3OutputStorage := mps.S3OutputStorage{}
				if v, ok := s3OutputStorageMap["s3_bucket"]; ok {
					s3OutputStorage.S3Bucket = helper.String(v.(string))
				}
				if v, ok := s3OutputStorageMap["s3_region"]; ok {
					s3OutputStorage.S3Region = helper.String(v.(string))
				}
				if v, ok := s3OutputStorageMap["s3_secret_id"]; ok {
					s3OutputStorage.S3SecretId = helper.String(v.(string))
				}
				if v, ok := s3OutputStorageMap["s3_secret_key"]; ok {
					s3OutputStorage.S3SecretKey = helper.String(v.(string))
				}
				taskOutputStorage.S3OutputStorage = &s3OutputStorage
			}
			request.OutputStorage = &taskOutputStorage
		}
	}

	if d.HasChange("output_dir") {
		if v, ok := d.GetOk("output_dir"); ok {
			request.OutputDir = helper.String(v.(string))
		}
	}

	if d.HasChange("task_notify_config") {
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
			if awsSQSMap, ok := helper.InterfaceToMap(dMap, "aws_sqs"); ok {
				awsSQS := mps.AwsSQS{}
				if v, ok := awsSQSMap["sqs_region"]; ok {
					awsSQS.SQSRegion = helper.String(v.(string))
				}
				if v, ok := awsSQSMap["sqs_queue_name"]; ok {
					awsSQS.SQSQueueName = helper.String(v.(string))
				}
				if v, ok := awsSQSMap["s3_secret_id"]; ok {
					awsSQS.S3SecretId = helper.String(v.(string))
				}
				if v, ok := awsSQSMap["s3_secret_key"]; ok {
					awsSQS.S3SecretKey = helper.String(v.(string))
				}
				taskNotifyConfig.AwsSQS = &awsSQS
			}
			request.TaskNotifyConfig = &taskNotifyConfig
		}
	}

	if d.HasChange("resource_id") {
		if v, ok := d.GetOk("resource_id"); ok {
			request.ResourceId = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().ModifySchedule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mps schedule failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMpsScheduleRead(d, meta)
}

func resourceTencentCloudMpsScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_schedule.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	scheduleId := d.Id()

	if err := service.DeleteMpsScheduleById(ctx, scheduleId); err != nil {
		return err
	}

	return nil
}
