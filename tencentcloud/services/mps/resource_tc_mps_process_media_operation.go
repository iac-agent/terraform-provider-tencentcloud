package mps

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMpsProcessMediaOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsProcessMediaOperationCreate,
		Read:   resourceTencentCloudMpsProcessMediaOperationRead,
		Delete: resourceTencentCloudMpsProcessMediaOperationDelete,
		Schema: map[string]*schema.Schema{
			"input_info": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "信息 的 文件 到 process。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "input 类型 有效 值:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
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

			"output_storage": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "存储 location 的 media processing 输出文件 如果此参数为空， 存储 location 在 `InputInfo` 将 是 inherited。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "存储 类型 对于 media processing 输出文件 有效 值:`COS`: Tencent Cloud COS`&gt;AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
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
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "directory 到 save media processing 输出文件，其中 必须 start 和 end 使用 `/`，such 作为 `/movie/201907/`.如果 您 do 不 指定this 参数， 文件 将 是 saved 到 directory 指定 在 `InputInfo`。",
			},

			"schedule_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "scheme ID.注意 1: About `OutputStorage` 和 `OutputDir`如果 output 存储 和 directory 是 指定 对于 subtask 的 scheme，those output settings 将 是 applied.如果 output 存储 和 directory 是 不 指定 对于 subtasks 的 scheme， output 参数 passed 在 `ProcessMedia` API 将 是 applied.注意 2: 如果 `TaskNotifyConfig` 是 指定， 指定 settings 将 是 使用 instead 的 默认值 callback settings 的 scheme.注意 3: 触发器 已配置 对于 scheme 是 对于 automatically starting scheme. It stops working 当 您 manually call 此 API 到 start scheme。",
			},

			"media_process_task": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "media processing 参数 到 使用。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"transcode_task_set": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 transcoding tasks。",
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
													Description: "是否remove 视频 数据. 有效 值:0: retain;1: remove.默认值：0。",
												},
												"remove_audio": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "是否remove 音频 数据. 有效 值:0: retain;1: remove.默认值：0。",
												},
												"video_template": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Video 流 配置 参数. 此 字段 为必填项 当 `RemoveVideo` 是 0。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"codec": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "视频 codec. 有效 值:`libx264`: H.264`libx265`: H.265`av1`: AOMedia Video 1Note: You 必须 指定a resolution (不 higher 比 640 x 480) 如果 H.265 codec 是 使用.注意: You 可以 仅 使用 AOMedia Video 1 codec 对于 MP4 files。",
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
																Description: "Resolution adaption. 有效 值:open: 已启用 当 resolution adaption 是 已启用，`宽度` 表示long side 的 视频，while `高度` 表示short side.close: 已禁用 当 resolution adaption 是 已禁用，`宽度` 表示width 的 视频，while `高度` 表示height.默认值：open.注意: 当 resolution adaption 是 已启用，`宽度` 不能 是 smaller 比 `高度`。",
															},
															"width": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum 值 的 宽度 (或 long side) 的 视频 流 （像素）。 取值范围：0 和 [128，4,096].如果 both `宽度` 和 `高度` 是 0， resolution 将 是 same 作为 该 的 来源 视频;如果 `宽度` 是 0，但 `高度` 是 不 0，`宽度` 将 是 proportionally scaled;如果 `宽度` 是 不 0，但 `高度` 是 0，`高度` 将 是 proportionally scaled;如果 both `宽度` 和 `高度` 是 不 0， 自定义 resolution 将 是 使用.默认值：0。",
															},
															"height": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum 值 的 高度 (或 short side) 的 视频 流 （像素）。 取值范围：0 和 [128，4,096].如果 both `宽度` 和 `高度` 是 0， resolution 将 是 same 作为 该 的 来源 视频;如果 `宽度` 是 0，但 `高度` 是 不 0，`宽度` 将 是 proportionally scaled;如果 `宽度` 是 不 0，但 `高度` 是 0，`高度` 将 是 proportionally scaled;如果 both `宽度` 和 `高度` 是 不 0， 自定义 resolution 将 是 使用.默认值：0。",
															},
															"gop": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Frame 间隔 between I keyframes. 取值范围：0 和 [1,100000].如果 此 参数 是 0 或 left 空， 系统 将 automatically 集合 GOP 长度。",
															},
															"fill_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "fill 模式，其中 表示how 视频 是 resized 当 视频's original aspect ratio 是 different 从 目标 aspect ratio. 有效 值:stretch: Stretch 镜像 frame 通过 frame 到 fill entire screen. 视频 镜像 可能 become squashed 或 stretched after transcoding.black: Keep 镜像&#39;s original aspect ratio 和 fill blank space 使用 black bars.white: Keep 镜像's original aspect ratio 和 fill blank space 使用 white bars.gauss: Keep 镜像's original aspect ratio 和 apply Gaussian blur 到 blank space.默认值：black.注意: Only `stretch` 和 `black` 是 支持 对于 adaptive bitrate streaming。",
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
													Description: "Audio 流 配置 参数. 此 字段 为必填项 当 `RemoveAudio` 是 0。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"codec": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Audio 流 codec.当 outer `Container` 参数 是 `mp3`， 有效 值 是:libmp3lame.当 outer `Container` 参数 是 `ogg` 或 `flac`， 有效 值 是:flac.当 outer `Container` 参数 是 `m4a`， 有效 值 include:libfdk_aac;libmp3lame;ac3.当 outer `Container` 参数 是 `mp4` 或 `flv`， 有效 值 include:libfdk_aac: more suitable 对于 mp4;libmp3lame: more suitable 对于 flv.当 outer `Container` 参数 是 `hls`， 有效 值 include:libfdk_aac;libmp3lame。",
															},
															"bitrate": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Audio 流 bitrate 在 Kbps. 取值范围：0 和 [26，256].如果 值 是 0， bitrate 的 音频 流 将 是 same 作为 该 的 original 音频。",
															},
															"sample_rate": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Audio 流 sample 速率. 有效 值:32,00044,10048,000In Hz。",
															},
															"audio_channel": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Audio channel 系统. 有效 值:1: Mono2: Dual6: StereoWhen media 是 packaged 在 音频 格式 (FLAC，OGG，MP3，M4A)， sound channel 不能 是 集合 到 stereo.默认值：2。",
															},
														},
													},
												},
												"tehd_config": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "TESHD transcoding 参数。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "TESHD 类型 有效 值:TEHD-100: TESHD-100.如果此参数为空，TESHD 将 不 是 已启用",
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
													Description: "是否remove 视频 数据. 有效 值:0: retain1: remove。",
												},
												"remove_audio": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "是否remove 音频 数据. 有效 值:0: retain1: remove。",
												},
												"video_template": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Video 流 配置 参数。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"codec": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "视频 codec. 有效 值:libx264: H.264libx265: H.265av1: AOMedia Video 1Note: You 必须 指定a resolution (不 higher 比 640 x 480) 如果 H.265 codec 是 使用.注意: You 可以 仅 使用 AOMedia Video 1 codec 对于 MP4 files。",
															},
															"fps": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Video frame 速率 在 Hz. 取值范围：[0，100].如果 值 是 0， frame 速率 将 是 same 作为 该 的 来源 视频。",
															},
															"bitrate": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Bitrate 的 视频 流 在 Kbps. 取值范围：0 和 [128，35,000].如果 值 是 0， bitrate 的 视频 将 是 same 作为 该 的 来源 视频。",
															},
															"resolution_adaptive": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Resolution adaption. 有效 值:open: 已启用 当 resolution adaption 是 已启用，`宽度` 表示long side 的 视频，while `高度` 表示short side.close: 已禁用 当 resolution adaption 是 已禁用，`宽度` 表示width 的 视频，while `高度` 表示height.注意: 当 resolution adaption 是 已启用，`宽度` 不能 是 smaller 比 `高度`。",
															},
															"width": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum 值 的 宽度 (或 long side) 的 视频 流 （像素）。 取值范围：0 和 [128，4,096].如果 both `宽度` 和 `高度` 是 0， resolution 将 是 same 作为 该 的 来源 视频;如果 `宽度` 是 0，但 `高度` 是 不 0，`宽度` 将 是 proportionally scaled;如果 `宽度` 是 不 0，但 `高度` 是 0，`高度` 将 是 proportionally scaled;如果 both `宽度` 和 `高度` 是 不 0， 自定义 resolution 将 是 使用。",
															},
															"height": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum 值 的 高度 (或 short side) 的 视频 流 （像素）。 取值范围：0 和 [128，4,096]。",
															},
															"gop": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Frame 间隔 between I keyframes. 取值范围：0 和 [1,100000]. 如果 此 参数 是 0， 系统 将 automatically 集合 GOP 长度。",
															},
															"fill_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Fill 类型 Fill refers 到 way 的 processing screenshot 当 its aspect ratio 是 different 从 该 的 来源 视频. following fill types 是 支持: stretch: stretch. screenshot 将 是 stretched frame 通过 frame 到 match aspect ratio 的 来源 视频，其中 可能 make screenshot shorter 或 longer;black: fill 使用 black. 此 选项 retains aspect ratio 的 来源 视频 对于 screenshot 和 fills unmatched area 使用 black color blocks.white: fill 使用 white. 此 选项 retains aspect ratio 的 来源 视频 对于 screenshot 和 fills unmatched area 使用 white color blocks.gauss: fill 使用 Gaussian blur. 此 选项 retains aspect ratio 的 来源 视频 对于 screenshot 和 fills unmatched area 使用 Gaussian blur。",
															},
															"vcrf": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "control factor 的 视频 constant bitrate. 取值范围：[0，51]. 此 参数 将 是 已禁用 如果 您 enter `0`.It 是 不 recommended 到 指定this 参数 如果 there 是 无 special requirements。",
															},
															"content_adapt_stream": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "是否enable adaptive 编码. 有效 值:0: Disable1: Enable默认值：0. 如果 此 参数 是 集合 到 `1`，多个 streams 使用 different resolutions 和 bitrates 将 是 generated automatically. highest resolution，bitrate，和 quality 的 streams 是 determined 通过 值 的 `宽度` 和 `高度`，`Bitrate`，和 `Vcrf` 在 `VideoTemplate` respectively. 如果 these 参数 是 不 集合 在 `VideoTemplate`， highest resolution generated 将 是 same 作为 该 的 来源 视频，和 highest 视频 quality 将 是 close 到 VMAF 95. To 使用 此 参数 或 learn about billing details 的 adaptive 编码，please contact your sales rep。",
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
																Description: "Audio 流 codec.当 outer `Container` 参数 是 `mp3`， 有效 值 是:libmp3lame.当 outer `Container` 参数 是 `ogg` 或 `flac`， 有效 值 是:flac.当 outer `Container` 参数 是 `m4a`， 有效 值 include:libfdk_aac;libmp3lame;ac3.当 outer `Container` 参数 是 `mp4` 或 `flv`， 有效 值 include:libfdk_aac: More suitable 对于 mp4;libmp3lame: More suitable 对于 flv;mp2.当 outer `Container` 参数 是 `hls`， 有效 值 include:libfdk_aac;libmp3lame。",
															},
															"bitrate": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Audio 流 bitrate 在 Kbps. 取值范围：0 和 [26，256]. 如果 值 是 0， bitrate 的 音频 流 将 是 same 作为 该 的 original 音频。",
															},
															"sample_rate": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Audio 流 sample 速率. 有效 值:32,00044,10048,000In Hz。",
															},
															"audio_channel": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Audio channel 系统. 有效 值:1: Mono2: Dual6: StereoWhen media 是 packaged 在 音频 格式 (FLAC，OGG，MP3，M4A)， sound channel 不能 是 集合 到 stereo。",
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
													Description: "TSC transcoding 参数.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "TSC 类型 有效 值:`TEHD-100`: TSC-100 (视频 TSC). `TEHD-200`: TSC-200 (音频 TSC). 如果 此 参数 是 left blank，无 modification 将 是 made.注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"max_video_bitrate": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "最大 视频 bitrate. 如果 此 参数 是 不 指定，无 modifications 将 是 made.注意：此字段可能返回 null，表示无法获取有效值。",
															},
														},
													},
												},
												"subtitle_template": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "subtitle settings.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"path": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "URL 的 subtitles 到 add 到 视频.注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"stream_index": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "subtitle track 到 add 到 视频. 如果 both `路径` 和 `StreamIndex` 是 指定，`路径` 将 是 使用. You need 到 指定at least 一个 的 two 参数.注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"font_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "font. 有效 值:`hei.ttf`: Heiti.`song.ttf`: Songti.`simkai.ttf`: Kaiti.`arial.ttf`: Arial. 默认为 `hei.ttf`.注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"font_size": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "font 大小 (pixels). 如果 此 是 不 指定， font 大小 在 subtitle 文件 将 是 使用.注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"font_color": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "font color 在 0xRRGGBB 格式 默认值：0xFFFFFF (white).注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"font_alpha": {
																Type:        schema.TypeFloat,
																Optional:    true,
																Description: "text transparency. 取值范围：0-1.`0`: Fully transparent.`1`: Fully opaque.默认值：1.注意：此字段可能返回 null，表示无法获取有效值。",
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
																Description: "input 类型 有效 值:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
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
													Description: "An extended 字段 对于 transcoding.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"add_on_subtitles": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "subtitle 文件 到 add.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "模式 有效 值:`subtitle-流`: Add subtitle track.`close-caption-708`: Embed CEA-708 subtitles 在 SEI frames.`close-caption-608`: Embed CEA-608 subtitles 在 SEI frames.注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"subtitle": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "subtitle 文件.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "input 类型 有效 值:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
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
													Description: "Custom 水印 参数，其中 是 有效 如果 `Definition` 是 0.此 参数 是 使用 在 highly customized scenarios. We recommend 您 使用 `Definition` 到 指定watermark 参数 preferably.Custom 水印 参数 是 不 可用 对于 screenshot。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Watermark 类型 有效 值:镜像: 镜像 水印。",
															},
															"coordinate_origin": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Origin position，其中 currently 可以 仅 是:TopLeft: 源站 的 coordinates 是 在 top-left corner 的 视频，和 源站 的 水印 是 在 top-left corner 的 镜像 或 text.默认值：TopLeft。",
															},
															"x_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "horizontal position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `XPos` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `XPos` 是 10% 的 视频 宽度;如果 字符串 结束 在 像素， `XPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `XPos` 是 100 像素.默认值：0 像素。",
															},
															"y_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "vertical position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `YPos` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `YPos` 是 10% 的 视频 高度;如果 字符串 结束 在 像素， `YPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `YPos` 是 100 像素.默认值：0 像素。",
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
																						Description: "input 类型 有效 值:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
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
																			Description: "Watermark 宽度. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `宽度` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `宽度` 是 10% 的 视频 宽度;如果 字符串 结束 在 像素， `宽度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `宽度` 是 100 像素.默认值：10%。",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Watermark 高度. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `高度` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `高度` 是 10% 的 视频 高度;如果 字符串 结束 在 像素， `高度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `高度` 是 100 像素.默认值：0 像素，其中 表示 该 `高度` 将 是 proportionally scaled according 到 aspect ratio 的 original 水印 镜像。",
																		},
																		"repeat_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
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
													Description: "开始时间 偏移量 的 水印 （秒）。 如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame.如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame;如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 appear 在 second n after first 视频 frame;如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 appear 在 second n before last 视频 frame。",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "结束时间 偏移量 的 水印 （秒）。如果此参数为空 或 0 是 entered， 水印 将 exist till last 视频 frame;如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 exist till second n;如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 exist till second n before last 视频 frame。",
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
													Description: "Origin position，其中 currently 可以 仅 是:TopLeft: 源站 的 coordinates 是 在 top-left corner 的 视频，和 源站 的 blur 是 在 top-left corner 的 镜像 或 text.默认值：TopLeft。",
												},
												"x_pos": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "horizontal position 的 源站 的 blur relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `XPos` 的 blur 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `XPos` 是 10% 的 视频 宽度;如果 字符串 结束 在 像素， `XPos` 的 blur 将 是 指定 像素; 对于 示例，`100px` 表示 该 `XPos` 是 100 像素.默认值：0 像素。",
												},
												"y_pos": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Vertical position 的 源站 的 blur relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `YPos` 的 blur 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `YPos` 是 10% 的 视频 高度;如果 字符串 结束 在 像素， `YPos` 的 blur 将 是 指定 像素; 对于 示例，`100px` 表示 该 `YPos` 是 100 像素.默认值：0 像素。",
												},
												"width": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Blur 宽度. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `宽度` 的 blur 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `宽度` 是 10% 的 视频 宽度;如果 字符串 结束 在 像素， `宽度` 的 blur 将 是 在 像素; 对于 示例，`100px` 表示 该 `宽度` 是 100 像素.默认值：10%。",
												},
												"height": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Blur 高度. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `高度` 的 blur 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `高度` 是 10% 的 视频 高度;如果 字符串 结束 在 像素， `高度` 的 blur 将 是 在 像素; 对于 示例，`100px` 表示 该 `高度` 是 100 像素.默认值：10%。",
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "开始时间 偏移量 的 blur （秒）。 如果此参数为空 或 0 是 entered， blur 将 appear upon first 视频 frame.如果此参数为空 或 0 是 entered， blur 将 appear upon first 视频 frame;如果 此 值 是 greater 比 0 (e.g.，n)， blur 将 appear 在 second n after first 视频 frame;如果 此 值 是 smaller 比 0 (e.g.，-n)， blur 将 appear 在 second n before last 视频 frame。",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "结束时间 偏移量 的 blur （秒）。如果此参数为空 或 0 是 entered， blur 将 exist till last 视频 frame;如果 此 值 是 greater 比 0 (e.g.，n)， blur 将 exist till second n;如果 此 值 是 smaller 比 0 (e.g.，-n)， blur 将 exist till second n before last 视频 frame。",
												},
											},
										},
									},
									"start_time_offset": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "开始时间 偏移量 的 transcoded 视频，（秒）。如果此参数为空 或 集合 到 0， transcoded 视频 将 start 在 same 时间 作为 original 视频.如果 此 参数 是 集合 到 positive 数量 (n 对于 示例)， transcoded 视频 将 start 在 nth second 的 original 视频.如果 此 参数 是 集合 到 negative 数量 (-n 对于 示例)， transcoded 视频 将 start 在 nth second before end 的 original 视频。",
									},
									"end_time_offset": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "结束时间 偏移量 的 transcoded 视频，（秒）。如果此参数为空 或 集合 到 0， transcoded 视频 将 end 在 same 时间 作为 original 视频.如果 此 参数 是 集合 到 positive 数量 (n 对于 示例)， transcoded 视频 将 end 在 nth second 的 original 视频.如果 此 参数 是 集合 到 negative 数量 (-n 对于 示例)， transcoded 视频 将 end 在 nth second before end 的 original 视频。",
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
													Description: "存储 类型 对于 media processing 输出文件 有效 值:`COS`: Tencent Cloud COS`&gt;AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
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
																Description: "input 类型 有效 值:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
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
																Description: "input 类型 有效 值:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
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
						"animated_graphic_task_set": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 animated 镜像 generating tasks。",
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
													Description: "存储 类型 对于 media processing 输出文件 有效 值:`COS`: Tencent Cloud COS`&gt;AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
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
						"snapshot_by_time_offset_task_set": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 时间 point screencapturing tasks。",
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
										Description: "列表 screenshot 时间 points 在 格式 的 `s` 或 `%`:如果 字符串 结束 在 `s`，它 表示 该 时间 point 是 在 秒; 对于 示例，`3.5s` 表示 该 时间 point 是 3.5th second;如果 字符串 结束 在 `%`，它 表示 该 时间 point 是 指定 percentage 的 视频 时长; 对于 示例，`10%` 表示 该 时间 point 是 10% 的 视频 时长。",
									},
									"time_offset_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeFloat,
										},
										Optional:    true,
										Description: "列表 时间 points 的 screenshots 在 &lt;font color=red&gt;秒&lt;/font&gt;。",
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
													Description: "Custom 水印 参数，其中 是 有效 如果 `Definition` 是 0.此 参数 是 使用 在 highly customized scenarios. We recommend 您 使用 `Definition` 到 指定watermark 参数 preferably.Custom 水印 参数 是 不 可用 对于 screenshot。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Watermark 类型 有效 值:镜像: 镜像 水印。",
															},
															"coordinate_origin": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Origin position，其中 currently 可以 仅 是:TopLeft: 源站 的 coordinates 是 在 top-left corner 的 视频，和 源站 的 水印 是 在 top-left corner 的 镜像 或 text.默认值：TopLeft。",
															},
															"x_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "horizontal position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `XPos` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `XPos` 是 10% 的 视频 宽度;如果 字符串 结束 在 像素， `XPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `XPos` 是 100 像素.默认值：0 像素。",
															},
															"y_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "vertical position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `YPos` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `YPos` 是 10% 的 视频 高度;如果 字符串 结束 在 像素， `YPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `YPos` 是 100 像素.默认值：0 像素。",
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
																						Description: "input 类型 有效 值:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
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
																			Description: "Watermark 宽度. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `宽度` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `宽度` 是 10% 的 视频 宽度;如果 字符串 结束 在 像素， `宽度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `宽度` 是 100 像素.默认值：10%。",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Watermark 高度. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `高度` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `高度` 是 10% 的 视频 高度;如果 字符串 结束 在 像素， `高度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `高度` 是 100 像素.默认值：0 像素，其中 表示 该 `高度` 将 是 proportionally scaled according 到 aspect ratio 的 original 水印 镜像。",
																		},
																		"repeat_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
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
													Description: "开始时间 偏移量 的 水印 （秒）。 如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame.如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame;如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 appear 在 second n after first 视频 frame;如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 appear 在 second n before last 视频 frame。",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "结束时间 偏移量 的 水印 （秒）。如果此参数为空 或 0 是 entered， 水印 将 exist till last 视频 frame;如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 exist till second n;如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 exist till second n before last 视频 frame。",
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
													Description: "存储 类型 对于 media processing 输出文件 有效 值:`COS`: Tencent Cloud COS`&gt;AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
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
						"sample_snapshot_task_set": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 sampled screencapturing tasks。",
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
													Description: "Custom 水印 参数，其中 是 有效 如果 `Definition` 是 0.此 参数 是 使用 在 highly customized scenarios. We recommend 您 使用 `Definition` 到 指定watermark 参数 preferably.Custom 水印 参数 是 不 可用 对于 screenshot。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Watermark 类型 有效 值:镜像: 镜像 水印。",
															},
															"coordinate_origin": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Origin position，其中 currently 可以 仅 是:TopLeft: 源站 的 coordinates 是 在 top-left corner 的 视频，和 源站 的 水印 是 在 top-left corner 的 镜像 或 text.默认值：TopLeft。",
															},
															"x_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "horizontal position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `XPos` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `XPos` 是 10% 的 视频 宽度;如果 字符串 结束 在 像素， `XPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `XPos` 是 100 像素.默认值：0 像素。",
															},
															"y_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "vertical position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `YPos` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `YPos` 是 10% 的 视频 高度;如果 字符串 结束 在 像素， `YPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `YPos` 是 100 像素.默认值：0 像素。",
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
																						Description: "input 类型 有效 值:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
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
																			Description: "Watermark 宽度. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `宽度` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `宽度` 是 10% 的 视频 宽度;如果 字符串 结束 在 像素， `宽度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `宽度` 是 100 像素.默认值：10%。",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Watermark 高度. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `高度` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `高度` 是 10% 的 视频 高度;如果 字符串 结束 在 像素， `高度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `高度` 是 100 像素.默认值：0 像素，其中 表示 该 `高度` 将 是 proportionally scaled according 到 aspect ratio 的 original 水印 镜像。",
																		},
																		"repeat_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
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
													Description: "开始时间 偏移量 的 水印 （秒）。 如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame.如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame;如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 appear 在 second n after first 视频 frame;如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 appear 在 second n before last 视频 frame。",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "结束时间 偏移量 的 水印 （秒）。如果此参数为空 或 0 是 entered， 水印 将 exist till last 视频 frame;如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 exist till second n;如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 exist till second n before last 视频 frame。",
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
													Description: "存储 类型 对于 media processing 输出文件 有效 值:`COS`: Tencent Cloud COS`&gt;AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
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
						"image_sprite_task_set": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 镜像 sprite generating tasks。",
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
													Description: "存储 类型 对于 media processing 输出文件 有效 值:`COS`: Tencent Cloud COS`&gt;AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
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
						"adaptive_dynamic_streaming_task_set": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 adaptive bitrate streaming tasks。",
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
													Description: "Custom 水印 参数，其中 是 有效 如果 `Definition` 是 0.此 参数 是 使用 在 highly customized scenarios. We recommend 您 使用 `Definition` 到 指定watermark 参数 preferably.Custom 水印 参数 是 不 可用 对于 screenshot。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Watermark 类型 有效 值:镜像: 镜像 水印。",
															},
															"coordinate_origin": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Origin position，其中 currently 可以 仅 是:TopLeft: 源站 的 coordinates 是 在 top-left corner 的 视频，和 源站 的 水印 是 在 top-left corner 的 镜像 或 text.默认值：TopLeft。",
															},
															"x_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "horizontal position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `XPos` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `XPos` 是 10% 的 视频 宽度;如果 字符串 结束 在 像素， `XPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `XPos` 是 100 像素.默认值：0 像素。",
															},
															"y_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "vertical position 的 源站 的 水印 relative 到 源站 的 coordinates 的 视频. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `YPos` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `YPos` 是 10% 的 视频 高度;如果 字符串 结束 在 像素， `YPos` 的 水印 将 是 指定 像素; 对于 示例，`100px` 表示 该 `YPos` 是 100 像素.默认值：0 像素。",
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
																						Description: "input 类型 有效 值:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
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
																			Description: "Watermark 宽度. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `宽度` 的 水印 将 是 指定 percentage 的 视频 宽度; 对于 示例，`10%` 表示 该 `宽度` 是 10% 的 视频 宽度;如果 字符串 结束 在 像素， `宽度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `宽度` 是 100 像素.默认值：10%。",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Watermark 高度. % 和 像素 formats 是 支持:如果 字符串 结束 在 %， `高度` 的 水印 将 是 指定 percentage 的 视频 高度; 对于 示例，`10%` 表示 该 `高度` 是 10% 的 视频 高度;如果 字符串 结束 在 像素， `高度` 的 水印 将 是 在 像素; 对于 示例，`100px` 表示 该 `高度` 是 100 像素.默认值：0 像素，其中 表示 该 `高度` 将 是 proportionally scaled according 到 aspect ratio 的 original 水印 镜像。",
																		},
																		"repeat_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
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
													Description: "开始时间 偏移量 的 水印 （秒）。 如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame.如果此参数为空 或 0 是 entered， 水印 将 appear upon first 视频 frame;如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 appear 在 second n after first 视频 frame;如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 appear 在 second n before last 视频 frame。",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "结束时间 偏移量 的 水印 （秒）。如果此参数为空 或 0 是 entered， 水印 将 exist till last 视频 frame;如果 此 值 是 greater 比 0 (e.g.，n)， 水印 将 exist till second n;如果 此 值 是 smaller 比 0 (e.g.，-n)， 水印 将 exist till second n before last 视频 frame。",
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
													Description: "存储 类型 对于 media processing 输出文件 有效 值:`COS`: Tencent Cloud COS`&gt;AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
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
										Description: "subtitle 文件 到 add.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "模式 有效 值:`subtitle-流`: Add subtitle track.`close-caption-708`: Embed CEA-708 subtitles 在 SEI frames.`close-caption-608`: Embed CEA-608 subtitles 在 SEI frames.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"subtitle": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "subtitle 文件.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "input 类型 有效 值:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，此 类型 是 仅 支持 对于 transcoding tasks。",
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
					},
				},
			},

			"ai_content_review_task": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "类型 参数 的 视频 内容 audit 任务。",
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
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Video 内容 analysis 任务 参数。",
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
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "类型 参数 的 视频 内容 recognition 任务。",
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

			"ai_quality_control_task": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "参数 的 quality control 任务。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "ID quality control template.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"channel_ext_para": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "channel extension 参数，其中 是 serialized JSON 字符串.注意：此字段可能返回 null，表示无法获取有效值。",
						},
					},
				},
			},

			"task_notify_config": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Event 通知 信息 的 任务. 如果此参数为空，无 事件 notifications 将 是 获取。",
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
							Description: "Workflow 通知 方法. 有效值：Finish，Change. 如果此参数为空，`Finish` 将 是 使用。",
						},
						"notify_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "通知 类型 有效 值:`CMQ`: 此 值 是 无 longer 使用. Please 使用 `TDMQ-CMQ` instead.`TDMQ-CMQ`: 消息 queue`URL`: 如果 `NotifyType` 是 集合 到 `URL`，HTTP callbacks 是 sent 到 URL 指定 通过 `NotifyUrl`. HTTP 和 JSON 是 用于the callbacks. packet 包含response 参数 的 `ParseNotification` API.`SCF`: 此 通知 类型 是 不 recommended. You need 到 configure 它 在 SCF console.`AWS-SQS`: AWS queue. 此 类型 是 仅 支持 对于 AWS tasks，和 queue 必须 是 在 same 地域 作为 AWS 存储桶&lt;font color=red&gt;注意: 如果 您 do 不 pass 此 参数 或 pass 在 空 字符串，`CMQ` 将 是 使用. To 使用 different 通知 类型，指定this 参数 accordingly.&lt;/font&gt;。",
						},
						"notify_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "HTTP callback URL，必填 如果 `NotifyType` 是 集合 到 `URL`。",
						},
						"aws_sqa": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "AWS SQS queue. 此 参数 为必填项 如果 `NotifyType` 是 `AWS-SQS`.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sqa_region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "地域 的 SQS queue。",
									},
									"sqa_queue_name": {
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

			"tasks_priority": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "任务 flow 优先级 higher 值， higher 优先级 取值范围：[-10，10]. 如果此参数为空，0 将 是 使用。",
			},

			"session_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID 用于deduplication. 如果 there 是 请求 使用 same ID 在 last three days， 当前 请求 将 返回 错误 ID 可以 contain up 到 50 字符. 如果此参数为空 或 空 字符串 是 entered，无 deduplication 将 是 performed。",
			},

			"session_context": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "来源 context 其中 是 用于pass through 用户 请求 信息. 任务 flow 状态 change callback 将 返回 值 的 此 字段. It 可以 contain up 到 1,000 字符。",
			},

			"task_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "任务 类型 `Online` (默认值): A 任务 该 是 executed immediately. `Offline`: A 任务 该 是 executed 当 系统 是 idle (within three days 通过 默认值)。",
			},
		},
	}
}

func resourceTencentCloudMpsProcessMediaOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_process_media_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = mps.NewProcessMediaRequest()
		response = mps.NewProcessMediaResponse()
		taskId   string
	)
	if dMap, ok := helper.InterfacesHeadMap(d, "input_info"); ok {
		mediaInputInfo := mps.MediaInputInfo{}
		if v, ok := dMap["type"]; ok {
			mediaInputInfo.Type = helper.String(v.(string))
		}
		if cosInputInfoMap, ok := helper.InterfaceToMap(dMap, "cos_input_info"); ok {
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
		if urlInputInfoMap, ok := helper.InterfaceToMap(dMap, "url_input_info"); ok {
			urlInputInfo := mps.UrlInputInfo{}
			if v, ok := urlInputInfoMap["url"]; ok {
				urlInputInfo.Url = helper.String(v.(string))
			}
			mediaInputInfo.UrlInputInfo = &urlInputInfo
		}
		if s3InputInfoMap, ok := helper.InterfaceToMap(dMap, "s3_input_info"); ok {
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
		request.InputInfo = &mediaInputInfo
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

	if v, ok := d.GetOkExists("schedule_id"); v != nil && ok {
		request.ScheduleId = helper.IntInt64(v.(int))
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
						if extTimeOffsetSetSet[i] != nil {
							extTimeOffsetSet := extTimeOffsetSetSet[i].(string)
							snapshotByTimeOffsetTaskInput.ExtTimeOffsetSet = append(snapshotByTimeOffsetTaskInput.ExtTimeOffsetSet, &extTimeOffsetSet)
						}
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
				if v, ok := adaptiveDynamicStreamingTaskSetMap["output_object_path"]; ok {
					adaptiveDynamicStreamingTaskInput.OutputObjectPath = helper.String(v.(string))
				}
				if v, ok := adaptiveDynamicStreamingTaskSetMap["sub_stream_object_name"]; ok {
					adaptiveDynamicStreamingTaskInput.SubStreamObjectName = helper.String(v.(string))
				}
				if v, ok := adaptiveDynamicStreamingTaskSetMap["segment_object_name"]; ok {
					adaptiveDynamicStreamingTaskInput.SegmentObjectName = helper.String(v.(string))
				}
				if v, ok := adaptiveDynamicStreamingTaskSetMap["add_on_subtitles"]; ok {
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

	if dMap, ok := helper.InterfacesHeadMap(d, "ai_quality_control_task"); ok {
		aiQualityControlTaskInput := mps.AiQualityControlTaskInput{}
		if v, ok := dMap["definition"]; ok {
			aiQualityControlTaskInput.Definition = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["channel_ext_para"]; ok {
			aiQualityControlTaskInput.ChannelExtPara = helper.String(v.(string))
		}
		request.AiQualityControlTask = &aiQualityControlTaskInput
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
		if awsSQSMap, ok := helper.InterfaceToMap(dMap, "aws_sqa"); ok {
			awsSQS := mps.AwsSQS{}
			if v, ok := awsSQSMap["sqa_region"]; ok {
				awsSQS.SQSRegion = helper.String(v.(string))
			}
			if v, ok := awsSQSMap["sqa_queue_name"]; ok {
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

	if v, ok := d.GetOkExists("tasks_priority"); v != nil && ok {
		request.TasksPriority = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("session_id"); ok {
		request.SessionId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("session_context"); ok {
		request.SessionContext = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_type"); ok {
		request.TaskType = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().ProcessMedia(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mps processMediaOperation failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(taskId)

	return resourceTencentCloudMpsProcessMediaOperationRead(d, meta)
}

func resourceTencentCloudMpsProcessMediaOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_process_media_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMpsProcessMediaOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_process_media_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
