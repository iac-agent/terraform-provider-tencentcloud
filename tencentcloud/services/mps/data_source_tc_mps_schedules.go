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
				Description: "The IDs of the schemes to query. Array length 限制: 100。",
			},

			"trigger_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The trigger 类型 Valid values:`CosFileUpload`: The scheme is triggered when a file is uploaded to Tencent Cloud Object Storage (COS).`AwsS3FileUpload`: The scheme is triggered when a file is uploaded to AWS S3.If you do not 指定this parameter or leave it empty，all schemes will be returned regardless of the trigger 类型",
			},

			"status": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The scheme 状态 Valid values:`已启用`，`已禁用`. If you do not 指定this parameter，all schemes will be returned regardless of the 状态",
			},

			"schedule_info_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The information of the schemes。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The scheme ID。",
						},
						"schedule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The scheme 名称注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The scheme 状态 Valid values:`已启用``已禁用`注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"trigger": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The trigger of the scheme.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The trigger 类型 Valid values:`CosFileUpload`: Tencent Cloud COS trigger.`AwsS3FileUpload`: AWS S3 trigger. Currently，this 类型 is only supported for transcoding tasks and schemes (not supported for workflows)。",
									},
									"cos_file_upload_trigger": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "This parameter 为必填项 and valid when `类型` is `CosFileUpload`，indicating the COS trigger rule.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bucket": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 COS 存储桶 bound to a workflow，such as `TopRankVideo-125xxx88`。",
												},
												"region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "地域 of the COS 存储桶 bound to a workflow，such as `ap-chongiqng`。",
												},
												"dir": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Input 路径 directory bound to a workflow，such as `/movie/201907/`. 如果此参数为空，the `/` root directory will be used。",
												},
												"formats": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "格式 列表 files that can trigger a workflow，such as [mp4，flv，mov]. 如果此参数为空，files in all formats can trigger the workflow。",
												},
											},
										},
									},
									"aws_s3_file_upload_trigger": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The AWS S3 trigger. This parameter is valid and 必填 if `类型` is `AwsS3FileUpload`.Note: Currently，the 键 for the AWS S3 存储桶，the trigger SQS queue，and the callback SQS queue must be the same.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"s3_bucket": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The AWS S3 存储桶 bound to the scheme。",
												},
												"s3_region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 地域 of the AWS S3 存储桶",
												},
												"dir": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 存储桶 directory bound. It must be an absolute 路径 that starts and ends with `/`，such as `/movie/201907/`. If you do not 指定this，the root directory will be bound.	。",
												},
												"formats": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "The file formats that will trigger the scheme，such as [mp4，flv，mov]. If you do not 指定this，the upload of files in any 格式 will trigger the scheme.	。",
												},
												"s3_secret_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 键 ID AWS S3 存储桶注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"s3_secret_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 键 of the AWS S3 存储桶注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"aws_sqs": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The SQS queue of the AWS S3 存储桶Note: The queue must be in the same 地域 as the 存储桶注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"sqs_region": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The 地域 of the SQS queue。",
															},
															"sqs_queue_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The 名称 SQS queue。",
															},
															"s3_secret_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The 键 ID 必填 to read from/write to the SQS queue。",
															},
															"s3_secret_key": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The 键 必填 to read from/write to the SQS queue。",
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
							Description: "The subtasks of the scheme.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"activity_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The subtask 类型`input`: The start.`output`: The end.`操作-trans`: Transcoding.`操作-samplesnapshot`: Sampled screencapturing.`操作-AIAnalysis`: 内容 analysis.`操作-AIRecognition`: 内容 recognition.`操作-aiReview`: 内容 moderation.`操作-animated-graphics`: Animated screenshot generation.`操作-image-sprite`: Image sprite generation.`操作-snapshotByTimeOffset`: Time point screencapturing.`操作-adaptive-substream`: Adaptive bitrate streaming.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"reardrive_index": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
										Computed:    true,
										Description: "The indexes of the subsequent actions.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"activity_para": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The parameters of a subtask.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"transcode_task": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "A transcoding task。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "ID video transcoding template。",
															},
															"raw_parameter": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Custom video transcoding parameter，which is valid if `Definition` is 0.This parameter is used in highly customized scenarios. We recommend you use `Definition` to 指定transcoding parameter preferably。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"container": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Container. 有效值：mp4; flv; hls; mp3; flac; ogg; m4a. Among them，mp3，flac，ogg，and m4a are for audio files。",
																		},
																		"remove_video": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "是否remove video data. Valid values:0: retain;1: remove.默认值：0。",
																		},
																		"remove_audio": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "是否remove audio data. Valid values:0: retain;1: remove.默认值：0。",
																		},
																		"video_template": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Video stream configuration parameter. This field 为必填项 when `RemoveVideo` is 0。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"codec": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The video codec. Valid values:`libx264`: H.264`libx265`: H.265`av1`: AOMedia Video 1Note: You must 指定a resolution (not higher than 640 x 480) if the H.265 codec is used.Note: You can only use the AOMedia Video 1 codec for MP4 files。",
																					},
																					"fps": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "The video frame rate (Hz). 取值范围：[0，100].If the 值 is 0，the frame rate will be the same as that of the 来源 video.Note: For adaptive bitrate streaming，the 值 range of this parameter is [0，60]。",
																					},
																					"bitrate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "The video bitrate (Kbps). 取值范围：0 and [128，35000].If the 值 is 0，the bitrate of the video will be the same as that of the 来源 video。",
																					},
																					"resolution_adaptive": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Resolution adaption. Valid values:open: 已启用 When resolution adaption is 已启用，`Width` 表示long side of a video，while `Height` 表示short side.close: 已禁用 When resolution adaption is 已禁用，`Width` 表示width of a video，while `Height` 表示height.默认值：open.Note: When resolution adaption is 已启用，`Width` cannot be smaller than `Height`。",
																					},
																					"width": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Maximum 值 of the width (or long side) of a video stream （像素）。 取值范围：0 and [128，4,096].If both `Width` and `Height` are 0，the resolution will be the same as that of the 来源 video;If `Width` is 0，but `Height` is not 0，`Width` will be proportionally scaled;If `Width` is not 0，but `Height` is 0，`Height` will be proportionally scaled;If both `Width` and `Height` are not 0，the custom resolution will be used.默认值：0。",
																					},
																					"height": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Maximum 值 of the height (or short side) of a video stream （像素）。 取值范围：0 and [128，4,096].If both `Width` and `Height` are 0，the resolution will be the same as that of the 来源 video;If `Width` is 0，but `Height` is not 0，`Width` will be proportionally scaled;If `Width` is not 0，but `Height` is 0，`Height` will be proportionally scaled;If both `Width` and `Height` are not 0，the custom resolution will be used.默认值：0。",
																					},
																					"gop": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Frame interval between I keyframes. 取值范围：0 and [1,100000].If this parameter is 0 or left empty，the system will automatically set the GOP length。",
																					},
																					"fill_type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The fill 模式，which 表示how a video is resized when the video&#39;s original aspect ratio is different from the target aspect ratio. Valid values:stretch: Stretch the image frame by frame to fill the entire screen. The video image may become squashed or stretched after transcoding.black: Keep the image&#39;s original aspect ratio and fill the blank space with black bars.white: Keep the image&#39;s original aspect ratio and fill the blank space with white bars.gauss: Keep the image&#39;s original aspect ratio and apply Gaussian blur to the blank space.默认值：black.Note: Only `stretch` and `black` are supported for adaptive bitrate streaming。",
																					},
																					"vcrf": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "The control factor of video constant bitrate. 取值范围：[1，51]If this parameter is specified，CRF (a bitrate control method) will be 用于transcoding. (Video bitrate will no longer take effect.)It is not recommended to 指定this parameter if there are no special requirements。",
																					},
																				},
																			},
																		},
																		"audio_template": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Audio stream configuration parameter. This field 为必填项 when `RemoveAudio` is 0。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"codec": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Audio stream codec.When the outer `Container` parameter is `mp3`，the valid 值 is:libmp3lame.When the outer `Container` parameter is `ogg` or `flac`，the valid 值 is:flac.When the outer `Container` parameter is `m4a`，the valid values include:libfdk_aac;libmp3lame;ac3.When the outer `Container` parameter is `mp4` or `flv`，the valid values include:libfdk_aac: more suitable for mp4;libmp3lame: more suitable for flv.When the outer `Container` parameter is `hls`，the valid values include:libfdk_aac;libmp3lame。",
																					},
																					"bitrate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Audio stream bitrate in Kbps. 取值范围：0 and [26，256].If the 值 is 0，the bitrate of the audio stream will be the same as that of the original audio。",
																					},
																					"sample_rate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Audio stream sample rate. Valid values:32,00044,10048,000In Hz。",
																					},
																					"audio_channel": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Audio channel system. Valid values:1: Mono2: Dual6: StereoWhen the media is packaged in audio 格式 (FLAC，OGG，MP3，M4A)，the sound channel cannot be set to stereo.默认值：2。",
																					},
																				},
																			},
																		},
																		"tehd_config": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "TESHD transcoding parameter。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "TESHD 类型 Valid values:`TEHD-100`: TESHD-100. 如果此参数为空，TESHD will not be 已启用",
																					},
																					"max_video_bitrate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Maximum bitrate，which is valid when `类型` is `TESHD`. 如果此参数为空 or 0 is entered，there will be no upper 限制 for bitrate。",
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
																Description: "Video transcoding custom parameter，which is valid when `Definition` is not 0.When any parameters in this structure are entered，they will be 用于override corresponding parameters in templates.This parameter is used in highly customized scenarios. We recommend you only use `Definition` to 指定transcoding parameter.Note: this field may return `null`，indicating that no valid 值 was found。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"container": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Container 格式 有效值：mp4，flv，hls，mp3，flac，ogg，and m4a; mp3，flac，ogg，and m4a are formats of audio files。",
																		},
																		"remove_video": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "是否remove video data. Valid values:0: retain1: remove。",
																		},
																		"remove_audio": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "是否remove audio data. Valid values:0: retain1: remove。",
																		},
																		"video_template": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Video stream configuration parameter。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"codec": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The video codec. Valid values:libx264: H.264libx265: H.265av1: AOMedia Video 1Note: You must 指定a resolution (not higher than 640 x 480) if the H.265 codec is used.Note: You can only use the AOMedia Video 1 codec for MP4 files。",
																					},
																					"fps": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Video frame rate in Hz. 取值范围：[0，100].If the 值 is 0，the frame rate will be the same as that of the 来源 video。",
																					},
																					"bitrate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Bitrate of a video stream in Kbps. 取值范围：0 and [128，35,000].If the 值 is 0，the bitrate of the video will be the same as that of the 来源 video。",
																					},
																					"resolution_adaptive": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Resolution adaption. Valid values:open: 已启用 When resolution adaption is 已启用，`Width` 表示long side of a video，while `Height` 表示short side.close: 已禁用 When resolution adaption is 已禁用，`Width` 表示width of a video，while `Height` 表示height.Note: When resolution adaption is 已启用，`Width` cannot be smaller than `Height`。",
																					},
																					"width": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Maximum 值 of the width (or long side) of a video stream （像素）。 取值范围：0 and [128，4,096].If both `Width` and `Height` are 0，the resolution will be the same as that of the 来源 video;If `Width` is 0，but `Height` is not 0，`Width` will be proportionally scaled;If `Width` is not 0，but `Height` is 0，`Height` will be proportionally scaled;If both `Width` and `Height` are not 0，the custom resolution will be used。",
																					},
																					"height": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Maximum 值 of the height (or short side) of a video stream （像素）。 取值范围：0 and [128，4,096]。",
																					},
																					"gop": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Frame interval between I keyframes. 取值范围：0 and [1,100000]. If this parameter is 0，the system will automatically set the GOP length。",
																					},
																					"fill_type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Fill 类型 Fill refers to the way of processing a screenshot when its aspect ratio is different from that of the 来源 video. The following fill types are supported: stretch: stretch. The screenshot will be stretched frame by frame to match the aspect ratio of the 来源 video，which may make the screenshot shorter or longer;black: fill with black. This option retains the aspect ratio of the 来源 video for the screenshot and fills the unmatched area with black color blocks.white: fill with white. This option retains the aspect ratio of the 来源 video for the screenshot and fills the unmatched area with white color blocks.gauss: fill with Gaussian blur. This option retains the aspect ratio of the 来源 video for the screenshot and fills the unmatched area with Gaussian blur。",
																					},
																					"vcrf": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "The control factor of video constant bitrate. 取值范围：[0，51]. This parameter will be 已禁用 if you enter `0`.It is not recommended to 指定this parameter if there are no special requirements。",
																					},
																					"content_adapt_stream": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "是否enable adaptive encoding. Valid values:0: Disable1: Enable默认值：0. If this parameter is set to `1`，multiple streams with different resolutions and bitrates will be generated automatically. The highest resolution，bitrate，and quality of the streams are determined by the values of `width` and `height`，`Bitrate`，and `Vcrf` in `VideoTemplate` respectively. If these parameters are not set in `VideoTemplate`，the highest resolution generated will be the same as that of the 来源 video，and the highest video quality will be close to VMAF 95. To use this parameter or learn about the billing details of adaptive encoding，please contact your sales rep。",
																					},
																				},
																			},
																		},
																		"audio_template": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Audio stream configuration parameter。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"codec": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Audio stream codec.When the outer `Container` parameter is `mp3`，the valid 值 is:libmp3lame.When the outer `Container` parameter is `ogg` or `flac`，the valid 值 is:flac.When the outer `Container` parameter is `m4a`，the valid values include:libfdk_aac;libmp3lame;ac3.When the outer `Container` parameter is `mp4` or `flv`，the valid values include:libfdk_aac: More suitable for mp4;libmp3lame: More suitable for flv;mp2.When the outer `Container` parameter is `hls`，the valid values include:libfdk_aac;libmp3lame。",
																					},
																					"bitrate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Audio stream bitrate in Kbps. 取值范围：0 and [26，256]. If the 值 is 0，the bitrate of the audio stream will be the same as that of the original audio。",
																					},
																					"sample_rate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Audio stream sample rate. Valid values:32,00044,10048,000In Hz。",
																					},
																					"audio_channel": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Audio channel system. Valid values:1: Mono2: Dual6: StereoWhen the media is packaged in audio 格式 (FLAC，OGG，MP3，M4A)，the sound channel cannot be set to stereo。",
																					},
																					"stream_selects": {
																						Type: schema.TypeSet,
																						Elem: &schema.Schema{
																							Type: schema.TypeInt,
																						},
																						Computed:    true,
																						Description: "The audio tracks to retain. All audio tracks are retained by default。",
																					},
																				},
																			},
																		},
																		"tehd_config": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The TSC transcoding parameters.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The TSC 类型 Valid values:`TEHD-100`: TSC-100 (video TSC). `TEHD-200`: TSC-200 (audio TSC). If this parameter is left blank，no modification will be made.注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																					"max_video_bitrate": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "The maximum video bitrate. If this parameter is not specified，no modifications will be made.注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																				},
																			},
																		},
																		"subtitle_template": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The subtitle settings.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"path": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The URL of the subtitles to add to the video.注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																					"stream_index": {
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "The subtitle track to add to the video. If both `路径` and `StreamIndex` are specified，`路径` will be used. You need to 指定at least one of the two parameters.注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																					"font_type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The font. Valid values:`hei.ttf`: Heiti.`song.ttf`: Songti.`simkai.ttf`: Kaiti.`arial.ttf`: Arial.The 默认为 `hei.ttf`.注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																					"font_size": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The font size (pixels). If this is not specified，the font size in the subtitle file will be used.注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																					"font_color": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The font color in 0xRRGGBB 格式 默认值：0xFFFFFF (white).注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																					"font_alpha": {
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "The text transparency. 取值范围：0-1.`0`: Fully transparent.`1`: Fully opaque.默认值：1.注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																				},
																			},
																		},
																		"addon_audio_stream": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The information of the external audio track to add.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The input 类型 Valid values:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，this 类型 is only supported for transcoding tasks。",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The information of the COS object to process. This parameter is valid and 必填 when `类型` is `COS`。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The COS 存储桶 of the object to process，such as `TopRankVideo-125xxx88`。",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 地域 of the COS 存储桶，such as `ap-chongqing`。",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 路径 of the object to process，such as `/movie/201907/WildAnimal.mov`。",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The URL of the object to process. This parameter is valid and 必填 when `类型` is `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "URL of a video。",
																								},
																							},
																						},
																					},
																					"s3_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The information of the AWS S3 object processed. This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"s3_bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The AWS S3 存储桶",
																								},
																								"s3_region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 地域 of the AWS S3 存储桶",
																								},
																								"s3_object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 路径 of the AWS S3 object。",
																								},
																								"s3_secret_id": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 键 ID 必填 to access the AWS S3 object。",
																								},
																								"s3_secret_key": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 键 必填 to access the AWS S3 object。",
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
																			Description: "An extended field for transcoding.注意：此字段可能返回 null，表示无法获取有效值。",
																		},
																		"add_on_subtitles": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The subtitle file to add.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 模式 Valid values:`subtitle-stream`: Add a subtitle track.`close-caption-708`: Embed EA-708 subtitles in SEI frames.`close-caption-608`: Embed CEA-608 subtitles in SEI frames.注意：此字段可能返回 null，表示无法获取有效值。",
																					},
																					"subtitle": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The subtitle file.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"type": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The input 类型 Valid values:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，this 类型 is only supported for transcoding tasks。",
																								},
																								"cos_input_info": {
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "The information of the COS object to process. This parameter is valid and 必填 when `类型` is `COS`。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"bucket": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The COS 存储桶 of the object to process，such as `TopRankVideo-125xxx88`。",
																											},
																											"region": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The 地域 of the COS 存储桶，such as `ap-chongqing`。",
																											},
																											"object": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The 路径 of the object to process，such as `/movie/201907/WildAnimal.mov`。",
																											},
																										},
																									},
																								},
																								"url_input_info": {
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "The URL of the object to process. This parameter is valid and 必填 when `类型` is `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"url": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "URL of a video。",
																											},
																										},
																									},
																								},
																								"s3_input_info": {
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "The information of the AWS S3 object processed. This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"s3_bucket": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The AWS S3 存储桶",
																											},
																											"s3_region": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The 地域 of the AWS S3 存储桶",
																											},
																											"s3_object": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The 路径 of the AWS S3 object。",
																											},
																											"s3_secret_id": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The 键 ID 必填 to access the AWS S3 object。",
																											},
																											"s3_secret_key": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The 键 必填 to access the AWS S3 object。",
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
																Description: "列表 up to 10 image or text watermarks.注意：此字段可能返回 null，表示无法获取有效值。",
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
																			Description: "Custom watermark parameter，which is valid if `Definition` is 0.This parameter is used in highly customized scenarios. We recommend you use `Definition` to 指定watermark parameter preferably.Custom watermark parameter is not available for screenshot。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Watermark 类型 Valid values:image: image watermark。",
																					},
																					"coordinate_origin": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Origin position，which currently can only be:TopLeft: the origin of coordinates is in the top-left corner of the video，and the origin of the watermark is in the top-left corner of the image or text.默认值：TopLeft。",
																					},
																					"x_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The horizontal position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %，the `XPos` of the watermark will be the specified percentage of the video width; for example，`10%` means that `XPos` is 10% of the video width;If the string ends in px，the `XPos` of the watermark will be the specified px; for example，`100px` means that `XPos` is 100 px.默认值：0 px。",
																					},
																					"y_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The vertical position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %，the `YPos` of the watermark will be the specified percentage of the video height; for example，`10%` means that `YPos` is 10% of the video height;If the string ends in px，the `YPos` of the watermark will be the specified px; for example，`100px` means that `YPos` is 100 px.默认值：0 px。",
																					},
																					"image_template": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Image watermark template. This field 为必填项 when `类型` is `image` and is invalid when `类型` is `text`。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"image_content": {
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Input 内容 of watermark image. JPEG and PNG images are supported。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"type": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The input 类型 Valid values:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，this 类型 is only supported for transcoding tasks。",
																											},
																											"cos_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The information of the COS object to process. This parameter is valid and 必填 when `类型` is `COS`。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The COS 存储桶 of the object to process，such as `TopRankVideo-125xxx88`。",
																														},
																														"region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 地域 of the COS 存储桶，such as `ap-chongqing`。",
																														},
																														"object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 路径 of the object to process，such as `/movie/201907/WildAnimal.mov`。",
																														},
																													},
																												},
																											},
																											"url_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The URL of the object to process. This parameter is valid and 必填 when `类型` is `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"url": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "URL of a video。",
																														},
																													},
																												},
																											},
																											"s3_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The information of the AWS S3 object processed. This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"s3_bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The AWS S3 存储桶",
																														},
																														"s3_region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 地域 of the AWS S3 存储桶",
																														},
																														"s3_object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 路径 of the AWS S3 object。",
																														},
																														"s3_secret_id": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 键 ID 必填 to access the AWS S3 object。",
																														},
																														"s3_secret_key": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 键 必填 to access the AWS S3 object。",
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
																									Description: "Watermark width. % and px formats are supported:If the string ends in %，the `Width` of the watermark will be the specified percentage of the video width; for example，`10%` means that `Width` is 10% of the video width;If the string ends in px，the `Width` of the watermark will be in px; for example，`100px` means that `Width` is 100 px.默认值：10%。",
																								},
																								"height": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Watermark height. % and px formats are supported:If the string ends in %，the `Height` of the watermark will be the specified percentage of the video height; for example，`10%` means that `Height` is 10% of the video height;If the string ends in px，the `Height` of the watermark will be in px; for example，`100px` means that `Height` is 100 px.默认值：0 px，which means that `Height` will be proportionally scaled according to the aspect ratio of the original watermark image。",
																								},
																								"repeat_type": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Repeat 类型 an animated watermark. Valid values:`once`: no longer appears after watermark playback ends.`repeat_last_frame`: stays on the last frame after watermark playback ends.`repeat` (default): repeats the playback until the video ends。",
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
																			Description: "Text 内容 of up to 100 characters. This field 为必填项 only when the watermark 类型 is text.Text watermark is not available for screenshot。",
																		},
																		"svg_content": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "SVG 内容 of up to 2,000,000 characters. This field 为必填项 only when the watermark 类型 is `SVG`.SVG watermark is not available for screenshot。",
																		},
																		"start_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "开始时间 偏移量 of a watermark （秒）。 如果此参数为空 or 0 is entered，the watermark will appear upon the first video frame.如果此参数为空 or 0 is entered，the watermark will appear upon the first video frame;If this 值 is greater than 0 (e.g.，n)，the watermark will appear at second n after the first video frame;If this 值 is smaller than 0 (e.g.，-n)，the watermark will appear at second n before the last video frame。",
																		},
																		"end_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "结束时间 偏移量 of a watermark （秒）。如果此参数为空 or 0 is entered，the watermark will exist till the last video frame;If this 值 is greater than 0 (e.g.，n)，the watermark will exist till second n;If this 值 is smaller than 0 (e.g.，-n)，the watermark will exist till second n before the last video frame。",
																		},
																	},
																},
															},
															"mosaic_set": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "列表 blurs. Up to 10 ones can be supported。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"coordinate_origin": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Origin position，which currently can only be:TopLeft: the origin of coordinates is in the top-left corner of the video，and the origin of the blur is in the top-left corner of the image or text.默认值：TopLeft。",
																		},
																		"x_pos": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The horizontal position of the origin of the blur relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %，the `XPos` of the blur will be the specified percentage of the video width; for example，`10%` means that `XPos` is 10% of the video width;If the string ends in px，the `XPos` of the blur will be the specified px; for example，`100px` means that `XPos` is 100 px.默认值：0 px。",
																		},
																		"y_pos": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Vertical position of the origin of blur relative to the origin of coordinates of video. % and px formats are supported:If the string ends in %，the `YPos` of the blur will be the specified percentage of the video height; for example，`10%` means that `YPos` is 10% of the video height;If the string ends in px，the `YPos` of the blur will be the specified px; for example，`100px` means that `YPos` is 100 px.默认值：0 px。",
																		},
																		"width": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Blur width. % and px formats are supported:If the string ends in %，the `Width` of the blur will be the specified percentage of the video width; for example，`10%` means that `Width` is 10% of the video width;If the string ends in px，the `Width` of the blur will be in px; for example，`100px` means that `Width` is 100 px.默认值：10%。",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Blur height. % and px formats are supported:If the string ends in %，the `Height` of the blur will be the specified percentage of the video height; for example，`10%` means that `Height` is 10% of the video height;If the string ends in px，the `Height` of the blur will be in px; for example，`100px` means that `Height` is 100 px.默认值：10%。",
																		},
																		"start_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "开始时间 偏移量 of blur （秒）。 如果此参数为空 or 0 is entered，the blur will appear upon the first video frame.如果此参数为空 or 0 is entered，the blur will appear upon the first video frame;If this 值 is greater than 0 (e.g.，n)，the blur will appear at second n after the first video frame;If this 值 is smaller than 0 (e.g.，-n)，the blur will appear at second n before the last video frame。",
																		},
																		"end_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "结束时间 偏移量 of blur （秒）。如果此参数为空 or 0 is entered，the blur will exist till the last video frame;If this 值 is greater than 0 (e.g.，n)，the blur will exist till second n;If this 值 is smaller than 0 (e.g.，-n)，the blur will exist till second n before the last video frame。",
																		},
																	},
																},
															},
															"start_time_offset": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "开始时间 偏移量 of a transcoded video，（秒）。如果此参数为空 or set to 0，the transcoded video will start at the same time as the original video.If this parameter is set to a positive number (n for example)，the transcoded video will start at the nth second of the original video.If this parameter is set to a negative number (-n for example)，the transcoded video will start at the nth second before the end of the original video。",
															},
															"end_time_offset": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "结束时间 偏移量 of a transcoded video，（秒）。如果此参数为空 or set to 0，the transcoded video will end at the same time as the original video.If this parameter is set to a positive number (n for example)，the transcoded video will end at the nth second of the original video.If this parameter is set to a negative number (-n for example)，the transcoded video will end at the nth second before the end of the original video。",
															},
															"output_storage": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Target 存储桶 of an 输出文件 如果此参数为空，the `OutputStorage` 值 of the upper folder will be inherited.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The storage 类型 for a media processing 输出文件 Valid values:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This 类型 is only supported for AWS tasks，and the 输出存储桶 must be in the same 地域 as the 存储桶 of the 来源 file。",
																		},
																		"cos_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The location to save the output object in COS. This parameter is valid and 必填 when `类型` is COS.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 存储桶 to which the 输出文件 of media processing is saved，such as `TopRankVideo-125xxx88`. 如果此参数为空，the 值 of the upper layer will be inherited。",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 地域 of the 输出存储桶，such as `ap-chongqing`. 如果此参数为空，the 值 of the upper layer will be inherited。",
																					},
																				},
																			},
																		},
																		"s3_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The AWS S3 存储桶 to save the 输出文件 This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The AWS S3 存储桶",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 地域 of the AWS S3 存储桶",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 键 ID 必填 to upload files to the AWS S3 object。",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 键 必填 to upload files to the AWS S3 object。",
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
																Description: "路径 to a primary 输出文件，which can be a relative 路径 or an absolute 路径 如果此参数为空，the following relative 路径 will be used by 默认值：`{inputName}_transcode_{definition}.{格式}`。",
															},
															"segment_object_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "路径 to an 输出文件 part (the 路径 to ts during transcoding to HLS)，which can only be a relative 路径 如果此参数为空，the following relative 路径 will be used by 默认值：`{inputName}_transcode_{definition}_{number}.{格式}`。",
															},
															"object_number_format": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Rule of the `{number}` variable in the 输出路径 after transcoding.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"initial_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Start 值 of the `{number}` variable. 默认值：0。",
																		},
																		"increment": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Increment of the `{number}` variable. 默认值：1。",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "最小长度the `{number}` variable. A placeholder will be used if the variable length is below the minimum requirement. 默认值：1。",
																		},
																		"place_holder": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Placeholder used when the `{number}` variable length is below the minimum requirement. 默认值：0。",
																		},
																	},
																},
															},
															"head_tail_parameter": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Opening and closing credits parametersNote: this field may return `null`，indicating that no valid 值 was found。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"head_set": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Opening credits list。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The input 类型 Valid values:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，this 类型 is only supported for transcoding tasks。",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The information of the COS object to process. This parameter is valid and 必填 when `类型` is `COS`。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The COS 存储桶 of the object to process，such as `TopRankVideo-125xxx88`。",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 地域 of the COS 存储桶，such as `ap-chongqing`。",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 路径 of the object to process，such as `/movie/201907/WildAnimal.mov`。",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The URL of the object to process. This parameter is valid and 必填 when `类型` is `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "URL of a video。",
																								},
																							},
																						},
																					},
																					"s3_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The information of the AWS S3 object processed. This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"s3_bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The AWS S3 存储桶",
																								},
																								"s3_region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 地域 of the AWS S3 存储桶",
																								},
																								"s3_object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 路径 of the AWS S3 object。",
																								},
																								"s3_secret_id": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 键 ID 必填 to access the AWS S3 object。",
																								},
																								"s3_secret_key": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 键 必填 to access the AWS S3 object。",
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
																			Description: "Closing credits list。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The input 类型 Valid values:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，this 类型 is only supported for transcoding tasks。",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The information of the COS object to process. This parameter is valid and 必填 when `类型` is `COS`。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The COS 存储桶 of the object to process，such as `TopRankVideo-125xxx88`。",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 地域 of the COS 存储桶，such as `ap-chongqing`。",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 路径 of the object to process，such as `/movie/201907/WildAnimal.mov`。",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The URL of the object to process. This parameter is valid and 必填 when `类型` is `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "URL of a video。",
																								},
																							},
																						},
																					},
																					"s3_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The information of the AWS S3 object processed. This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"s3_bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The AWS S3 存储桶",
																								},
																								"s3_region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 地域 of the AWS S3 存储桶",
																								},
																								"s3_object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 路径 of the AWS S3 object。",
																								},
																								"s3_secret_id": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 键 ID 必填 to access the AWS S3 object。",
																								},
																								"s3_secret_key": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 键 必填 to access the AWS S3 object。",
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
													Description: "An animated screenshot generation task。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Animated image generating 模板 ID",
															},
															"start_time_offset": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "开始时间 of an animated image in a video （秒）。",
															},
															"end_time_offset": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "结束时间 of an animated image in a video （秒）。",
															},
															"output_storage": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Target 存储桶 of a generated animated image file. 如果此参数为空，the `OutputStorage` 值 of the upper folder will be inherited.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The storage 类型 for a media processing 输出文件 Valid values:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This 类型 is only supported for AWS tasks，and the 输出存储桶 must be in the same 地域 as the 存储桶 of the 来源 file。",
																		},
																		"cos_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The location to save the output object in COS. This parameter is valid and 必填 when `类型` is COS.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 存储桶 to which the 输出文件 of media processing is saved，such as `TopRankVideo-125xxx88`. 如果此参数为空，the 值 of the upper layer will be inherited。",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 地域 of the 输出存储桶，such as `ap-chongqing`. 如果此参数为空，the 值 of the upper layer will be inherited。",
																					},
																				},
																			},
																		},
																		"s3_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The AWS S3 存储桶 to save the 输出文件 This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The AWS S3 存储桶",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 地域 of the AWS S3 存储桶",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 键 ID 必填 to upload files to the AWS S3 object。",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 键 必填 to upload files to the AWS S3 object。",
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
																Description: "输出路径 to a generated animated image file，which can be a relative 路径 or an absolute 路径 如果此参数为空，the following relative 路径 will be used by 默认值：`{inputName}_animatedGraphic_{definition}.{格式}`。",
															},
														},
													},
												},
												"snapshot_by_time_offset_task": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "A time point screencapturing task。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "ID time point screencapturing template。",
															},
															"ext_time_offset_set": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Computed:    true,
																Description: "列表 screenshot time points in the 格式 of `s` or `%`:If the string ends in `s`，it means that the time point is in seconds; for example，`3.5s` means that the time point is the 3.5th second;If the string ends in `%`，it means that the time point is the specified percentage of the video duration; for example，`10%` means that the time point is 10% of the video duration。",
															},
															"watermark_set": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "列表 up to 10 image or text watermarks.注意：此字段可能返回 null，表示无法获取有效值。",
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
																			Description: "Custom watermark parameter，which is valid if `Definition` is 0.This parameter is used in highly customized scenarios. We recommend you use `Definition` to 指定watermark parameter preferably.Custom watermark parameter is not available for screenshot。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Watermark 类型 Valid values:image: image watermark。",
																					},
																					"coordinate_origin": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Origin position，which currently can only be:TopLeft: the origin of coordinates is in the top-left corner of the video，and the origin of the watermark is in the top-left corner of the image or text.默认值：TopLeft。",
																					},
																					"x_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The horizontal position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %，the `XPos` of the watermark will be the specified percentage of the video width; for example，`10%` means that `XPos` is 10% of the video width;If the string ends in px，the `XPos` of the watermark will be the specified px; for example，`100px` means that `XPos` is 100 px.默认值：0 px。",
																					},
																					"y_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The vertical position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %，the `YPos` of the watermark will be the specified percentage of the video height; for example，`10%` means that `YPos` is 10% of the video height;If the string ends in px，the `YPos` of the watermark will be the specified px; for example，`100px` means that `YPos` is 100 px.默认值：0 px。",
																					},
																					"image_template": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Image watermark template. This field 为必填项 when `类型` is `image` and is invalid when `类型` is `text`。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"image_content": {
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Input 内容 of watermark image. JPEG and PNG images are supported。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"type": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The input 类型 Valid values:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，this 类型 is only supported for transcoding tasks。",
																											},
																											"cos_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The information of the COS object to process. This parameter is valid and 必填 when `类型` is `COS`。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The COS 存储桶 of the object to process，such as `TopRankVideo-125xxx88`。",
																														},
																														"region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 地域 of the COS 存储桶，such as `ap-chongqing`。",
																														},
																														"object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 路径 of the object to process，such as `/movie/201907/WildAnimal.mov`。",
																														},
																													},
																												},
																											},
																											"url_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The URL of the object to process. This parameter is valid and 必填 when `类型` is `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"url": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "URL of a video。",
																														},
																													},
																												},
																											},
																											"s3_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The information of the AWS S3 object processed. This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"s3_bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The AWS S3 存储桶",
																														},
																														"s3_region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 地域 of the AWS S3 存储桶",
																														},
																														"s3_object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 路径 of the AWS S3 object。",
																														},
																														"s3_secret_id": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 键 ID 必填 to access the AWS S3 object。",
																														},
																														"s3_secret_key": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 键 必填 to access the AWS S3 object。",
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
																									Description: "Watermark width. % and px formats are supported:If the string ends in %，the `Width` of the watermark will be the specified percentage of the video width; for example，`10%` means that `Width` is 10% of the video width;If the string ends in px，the `Width` of the watermark will be in px; for example，`100px` means that `Width` is 100 px.默认值：10%。",
																								},
																								"height": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Watermark height. % and px formats are supported:If the string ends in %，the `Height` of the watermark will be the specified percentage of the video height; for example，`10%` means that `Height` is 10% of the video height;If the string ends in px，the `Height` of the watermark will be in px; for example，`100px` means that `Height` is 100 px.默认值：0 px，which means that `Height` will be proportionally scaled according to the aspect ratio of the original watermark image。",
																								},
																								"repeat_type": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Repeat 类型 an animated watermark. Valid values:`once`: no longer appears after watermark playback ends.`repeat_last_frame`: stays on the last frame after watermark playback ends.`repeat` (default): repeats the playback until the video ends。",
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
																			Description: "Text 内容 of up to 100 characters. This field 为必填项 only when the watermark 类型 is text.Text watermark is not available for screenshot。",
																		},
																		"svg_content": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "SVG 内容 of up to 2,000,000 characters. This field 为必填项 only when the watermark 类型 is `SVG`.SVG watermark is not available for screenshot。",
																		},
																		"start_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "开始时间 偏移量 of a watermark （秒）。 如果此参数为空 or 0 is entered，the watermark will appear upon the first video frame.如果此参数为空 or 0 is entered，the watermark will appear upon the first video frame;If this 值 is greater than 0 (e.g.，n)，the watermark will appear at second n after the first video frame;If this 值 is smaller than 0 (e.g.，-n)，the watermark will appear at second n before the last video frame。",
																		},
																		"end_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "结束时间 偏移量 of a watermark （秒）。如果此参数为空 or 0 is entered，the watermark will exist till the last video frame;If this 值 is greater than 0 (e.g.，n)，the watermark will exist till second n;If this 值 is smaller than 0 (e.g.，-n)，the watermark will exist till second n before the last video frame。",
																		},
																	},
																},
															},
															"output_storage": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Target 存储桶 of a generated time point screenshot file. 如果此参数为空，the `OutputStorage` 值 of the upper folder will be inherited.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The storage 类型 for a media processing 输出文件 Valid values:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This 类型 is only supported for AWS tasks，and the 输出存储桶 must be in the same 地域 as the 存储桶 of the 来源 file。",
																		},
																		"cos_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The location to save the output object in COS. This parameter is valid and 必填 when `类型` is COS.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 存储桶 to which the 输出文件 of media processing is saved，such as `TopRankVideo-125xxx88`. 如果此参数为空，the 值 of the upper layer will be inherited。",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 地域 of the 输出存储桶，such as `ap-chongqing`. 如果此参数为空，the 值 of the upper layer will be inherited。",
																					},
																				},
																			},
																		},
																		"s3_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The AWS S3 存储桶 to save the 输出文件 This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The AWS S3 存储桶",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 地域 of the AWS S3 存储桶",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 键 ID 必填 to upload files to the AWS S3 object。",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 键 必填 to upload files to the AWS S3 object。",
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
																Description: "输出路径 to a generated time point screenshot，which can be a relative 路径 or an absolute 路径 如果此参数为空，the following relative 路径 will be used by 默认值：`{inputName}_snapshotByTimeOffset_{definition}_{number}.{格式}`。",
															},
															"object_number_format": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Rule of the `{number}` variable in the time point screenshot 输出路径注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"initial_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Start 值 of the `{number}` variable. 默认值：0。",
																		},
																		"increment": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Increment of the `{number}` variable. 默认值：1。",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "最小长度the `{number}` variable. A placeholder will be used if the variable length is below the minimum requirement. 默认值：1。",
																		},
																		"place_holder": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Placeholder used when the `{number}` variable length is below the minimum requirement. 默认值：0。",
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
													Description: "A sampled screencapturing task。",
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
																Description: "列表 up to 10 image or text watermarks.注意：此字段可能返回 null，表示无法获取有效值。",
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
																			Description: "Custom watermark parameter，which is valid if `Definition` is 0.This parameter is used in highly customized scenarios. We recommend you use `Definition` to 指定watermark parameter preferably.Custom watermark parameter is not available for screenshot。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Watermark 类型 Valid values:image: image watermark。",
																					},
																					"coordinate_origin": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Origin position，which currently can only be:TopLeft: the origin of coordinates is in the top-left corner of the video，and the origin of the watermark is in the top-left corner of the image or text.默认值：TopLeft。",
																					},
																					"x_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The horizontal position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %，the `XPos` of the watermark will be the specified percentage of the video width; for example，`10%` means that `XPos` is 10% of the video width;If the string ends in px，the `XPos` of the watermark will be the specified px; for example，`100px` means that `XPos` is 100 px.默认值：0 px。",
																					},
																					"y_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The vertical position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %，the `YPos` of the watermark will be the specified percentage of the video height; for example，`10%` means that `YPos` is 10% of the video height;If the string ends in px，the `YPos` of the watermark will be the specified px; for example，`100px` means that `YPos` is 100 px.默认值：0 px。",
																					},
																					"image_template": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Image watermark template. This field 为必填项 when `类型` is `image` and is invalid when `类型` is `text`。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"image_content": {
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Input 内容 of watermark image. JPEG and PNG images are supported。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"type": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The input 类型 Valid values:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，this 类型 is only supported for transcoding tasks。",
																											},
																											"cos_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The information of the COS object to process. This parameter is valid and 必填 when `类型` is `COS`。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The COS 存储桶 of the object to process，such as `TopRankVideo-125xxx88`。",
																														},
																														"region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 地域 of the COS 存储桶，such as `ap-chongqing`。",
																														},
																														"object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 路径 of the object to process，such as `/movie/201907/WildAnimal.mov`。",
																														},
																													},
																												},
																											},
																											"url_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The URL of the object to process. This parameter is valid and 必填 when `类型` is `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"url": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "URL of a video。",
																														},
																													},
																												},
																											},
																											"s3_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The information of the AWS S3 object processed. This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"s3_bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The AWS S3 存储桶",
																														},
																														"s3_region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 地域 of the AWS S3 存储桶",
																														},
																														"s3_object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 路径 of the AWS S3 object。",
																														},
																														"s3_secret_id": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 键 ID 必填 to access the AWS S3 object。",
																														},
																														"s3_secret_key": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 键 必填 to access the AWS S3 object。",
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
																									Description: "Watermark width. % and px formats are supported:If the string ends in %，the `Width` of the watermark will be the specified percentage of the video width; for example，`10%` means that `Width` is 10% of the video width;If the string ends in px，the `Width` of the watermark will be in px; for example，`100px` means that `Width` is 100 px.默认值：10%。",
																								},
																								"height": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Watermark height. % and px formats are supported:If the string ends in %，the `Height` of the watermark will be the specified percentage of the video height; for example，`10%` means that `Height` is 10% of the video height;If the string ends in px，the `Height` of the watermark will be in px; for example，`100px` means that `Height` is 100 px.默认值：0 px，which means that `Height` will be proportionally scaled according to the aspect ratio of the original watermark image。",
																								},
																								"repeat_type": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Repeat 类型 an animated watermark. Valid values:`once`: no longer appears after watermark playback ends.`repeat_last_frame`: stays on the last frame after watermark playback ends.`repeat` (default): repeats the playback until the video ends。",
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
																			Description: "Text 内容 of up to 100 characters. This field 为必填项 only when the watermark 类型 is text.Text watermark is not available for screenshot。",
																		},
																		"svg_content": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "SVG 内容 of up to 2,000,000 characters. This field 为必填项 only when the watermark 类型 is `SVG`.SVG watermark is not available for screenshot。",
																		},
																		"start_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "开始时间 偏移量 of a watermark （秒）。 如果此参数为空 or 0 is entered，the watermark will appear upon the first video frame.如果此参数为空 or 0 is entered，the watermark will appear upon the first video frame;If this 值 is greater than 0 (e.g.，n)，the watermark will appear at second n after the first video frame;If this 值 is smaller than 0 (e.g.，-n)，the watermark will appear at second n before the last video frame。",
																		},
																		"end_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "结束时间 偏移量 of a watermark （秒）。如果此参数为空 or 0 is entered，the watermark will exist till the last video frame;If this 值 is greater than 0 (e.g.，n)，the watermark will exist till second n;If this 值 is smaller than 0 (e.g.，-n)，the watermark will exist till second n before the last video frame。",
																		},
																	},
																},
															},
															"output_storage": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Target 存储桶 of a sampled screenshot. 如果此参数为空，the `OutputStorage` 值 of the upper folder will be inherited.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The storage 类型 for a media processing 输出文件 Valid values:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This 类型 is only supported for AWS tasks，and the 输出存储桶 must be in the same 地域 as the 存储桶 of the 来源 file。",
																		},
																		"cos_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The location to save the output object in COS. This parameter is valid and 必填 when `类型` is COS.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 存储桶 to which the 输出文件 of media processing is saved，such as `TopRankVideo-125xxx88`. 如果此参数为空，the 值 of the upper layer will be inherited。",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 地域 of the 输出存储桶，such as `ap-chongqing`. 如果此参数为空，the 值 of the upper layer will be inherited。",
																					},
																				},
																			},
																		},
																		"s3_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The AWS S3 存储桶 to save the 输出文件 This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The AWS S3 存储桶",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 地域 of the AWS S3 存储桶",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 键 ID 必填 to upload files to the AWS S3 object。",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 键 必填 to upload files to the AWS S3 object。",
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
																Description: "输出路径 to a generated sampled screenshot，which can be a relative 路径 or an absolute 路径 如果此参数为空，the following relative 路径 will be used by 默认值：`{inputName}_sampleSnapshot_{definition}_{number}.{格式}`。",
															},
															"object_number_format": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Rule of the `{number}` variable in the sampled screenshot 输出路径注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"initial_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Start 值 of the `{number}` variable. 默认值：0。",
																		},
																		"increment": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Increment of the `{number}` variable. 默认值：1。",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "最小长度the `{number}` variable. A placeholder will be used if the variable length is below the minimum requirement. 默认值：1。",
																		},
																		"place_holder": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Placeholder used when the `{number}` variable length is below the minimum requirement. 默认值：0。",
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
													Description: "An image sprite generation task。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "ID an image sprite generating template。",
															},
															"output_storage": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Target 存储桶 of a generated image sprite. 如果此参数为空，the `OutputStorage` 值 of the upper folder will be inherited.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The storage 类型 for a media processing 输出文件 Valid values:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This 类型 is only supported for AWS tasks，and the 输出存储桶 must be in the same 地域 as the 存储桶 of the 来源 file。",
																		},
																		"cos_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The location to save the output object in COS. This parameter is valid and 必填 when `类型` is COS.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 存储桶 to which the 输出文件 of media processing is saved，such as `TopRankVideo-125xxx88`. 如果此参数为空，the 值 of the upper layer will be inherited。",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 地域 of the 输出存储桶，such as `ap-chongqing`. 如果此参数为空，the 值 of the upper layer will be inherited。",
																					},
																				},
																			},
																		},
																		"s3_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The AWS S3 存储桶 to save the 输出文件 This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The AWS S3 存储桶",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 地域 of the AWS S3 存储桶",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 键 ID 必填 to upload files to the AWS S3 object。",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 键 必填 to upload files to the AWS S3 object。",
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
																Description: "输出路径 to a generated image sprite file，which can be a relative 路径 or an absolute 路径 如果此参数为空，the following relative 路径 will be used by 默认值：`{inputName}_imageSprite_{definition}_{number}.{格式}`。",
															},
															"web_vtt_object_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "输出路径 to the WebVTT file after an image sprite is generated，which can only be a relative 路径 如果此参数为空，the following relative 路径 will be used by 默认值：`{inputName}_imageSprite_{definition}.{格式}`。",
															},
															"object_number_format": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Rule of the `{number}` variable in the image sprite 输出路径注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"initial_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Start 值 of the `{number}` variable. 默认值：0。",
																		},
																		"increment": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Increment of the `{number}` variable. 默认值：1。",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "最小长度the `{number}` variable. A placeholder will be used if the variable length is below the minimum requirement. 默认值：1。",
																		},
																		"place_holder": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Placeholder used when the `{number}` variable length is below the minimum requirement. 默认值：0。",
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
													Description: "An adaptive bitrate streaming task。",
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
																Description: "列表 up to 10 image or text watermarks。",
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
																			Description: "Custom watermark parameter，which is valid if `Definition` is 0.This parameter is used in highly customized scenarios. We recommend you use `Definition` to 指定watermark parameter preferably.Custom watermark parameter is not available for screenshot。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Watermark 类型 Valid values:image: image watermark。",
																					},
																					"coordinate_origin": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Origin position，which currently can only be:TopLeft: the origin of coordinates is in the top-left corner of the video，and the origin of the watermark is in the top-left corner of the image or text.默认值：TopLeft。",
																					},
																					"x_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The horizontal position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %，the `XPos` of the watermark will be the specified percentage of the video width; for example，`10%` means that `XPos` is 10% of the video width;If the string ends in px，the `XPos` of the watermark will be the specified px; for example，`100px` means that `XPos` is 100 px.默认值：0 px。",
																					},
																					"y_pos": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The vertical position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %，the `YPos` of the watermark will be the specified percentage of the video height; for example，`10%` means that `YPos` is 10% of the video height;If the string ends in px，the `YPos` of the watermark will be the specified px; for example，`100px` means that `YPos` is 100 px.默认值：0 px。",
																					},
																					"image_template": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Image watermark template. This field 为必填项 when `类型` is `image` and is invalid when `类型` is `text`。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"image_content": {
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Input 内容 of watermark image. JPEG and PNG images are supported。",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"type": {
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The input 类型 Valid values:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，this 类型 is only supported for transcoding tasks。",
																											},
																											"cos_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The information of the COS object to process. This parameter is valid and 必填 when `类型` is `COS`。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The COS 存储桶 of the object to process，such as `TopRankVideo-125xxx88`。",
																														},
																														"region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 地域 of the COS 存储桶，such as `ap-chongqing`。",
																														},
																														"object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 路径 of the object to process，such as `/movie/201907/WildAnimal.mov`。",
																														},
																													},
																												},
																											},
																											"url_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The URL of the object to process. This parameter is valid and 必填 when `类型` is `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"url": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "URL of a video。",
																														},
																													},
																												},
																											},
																											"s3_input_info": {
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The information of the AWS S3 object processed. This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"s3_bucket": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The AWS S3 存储桶",
																														},
																														"s3_region": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 地域 of the AWS S3 存储桶",
																														},
																														"s3_object": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 路径 of the AWS S3 object。",
																														},
																														"s3_secret_id": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 键 ID 必填 to access the AWS S3 object。",
																														},
																														"s3_secret_key": {
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The 键 必填 to access the AWS S3 object。",
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
																									Description: "Watermark width. % and px formats are supported:If the string ends in %，the `Width` of the watermark will be the specified percentage of the video width; for example，`10%` means that `Width` is 10% of the video width;If the string ends in px，the `Width` of the watermark will be in px; for example，`100px` means that `Width` is 100 px.默认值：10%。",
																								},
																								"height": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Watermark height. % and px formats are supported:If the string ends in %，the `Height` of the watermark will be the specified percentage of the video height; for example，`10%` means that `Height` is 10% of the video height;If the string ends in px，the `Height` of the watermark will be in px; for example，`100px` means that `Height` is 100 px.默认值：0 px，which means that `Height` will be proportionally scaled according to the aspect ratio of the original watermark image。",
																								},
																								"repeat_type": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Repeat 类型 an animated watermark. Valid values:`once`: no longer appears after watermark playback ends.`repeat_last_frame`: stays on the last frame after watermark playback ends.`repeat` (default): repeats the playback until the video ends。",
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
																			Description: "Text 内容 of up to 100 characters. This field 为必填项 only when the watermark 类型 is text.Text watermark is not available for screenshot。",
																		},
																		"svg_content": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "SVG 内容 of up to 2,000,000 characters. This field 为必填项 only when the watermark 类型 is `SVG`.SVG watermark is not available for screenshot。",
																		},
																		"start_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "开始时间 偏移量 of a watermark （秒）。 如果此参数为空 or 0 is entered，the watermark will appear upon the first video frame.如果此参数为空 or 0 is entered，the watermark will appear upon the first video frame;If this 值 is greater than 0 (e.g.，n)，the watermark will appear at second n after the first video frame;If this 值 is smaller than 0 (e.g.，-n)，the watermark will appear at second n before the last video frame。",
																		},
																		"end_time_offset": {
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "结束时间 偏移量 of a watermark （秒）。如果此参数为空 or 0 is entered，the watermark will exist till the last video frame;If this 值 is greater than 0 (e.g.，n)，the watermark will exist till second n;If this 值 is smaller than 0 (e.g.，-n)，the watermark will exist till second n before the last video frame。",
																		},
																	},
																},
															},
															"output_storage": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Target 存储桶 of an 输出文件 after being transcoded to adaptive bitrate streaming. 如果此参数为空，the `OutputStorage` 值 of the upper folder will be inherited.Note: this field may return null，indicating that no valid values can be obtained。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The storage 类型 for a media processing 输出文件 Valid values:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This 类型 is only supported for AWS tasks，and the 输出存储桶 must be in the same 地域 as the 存储桶 of the 来源 file。",
																		},
																		"cos_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The location to save the output object in COS. This parameter is valid and 必填 when `类型` is COS.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 存储桶 to which the 输出文件 of media processing is saved，such as `TopRankVideo-125xxx88`. 如果此参数为空，the 值 of the upper layer will be inherited。",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 地域 of the 输出存储桶，such as `ap-chongqing`. 如果此参数为空，the 值 of the upper layer will be inherited。",
																					},
																				},
																			},
																		},
																		"s3_output_storage": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The AWS S3 存储桶 to save the 输出文件 This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The AWS S3 存储桶",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 地域 of the AWS S3 存储桶",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 键 ID 必填 to upload files to the AWS S3 object。",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The 键 必填 to upload files to the AWS S3 object。",
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
																Description: "The relative or absolute 输出路径 of the manifest file after being transcoded to adaptive bitrate streaming. 如果此参数为空，a relative 路径 in the following 格式 will be used by 默认值：`{inputName}_adaptiveDynamicStreaming_{definition}.{格式}`。",
															},
															"sub_stream_object_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The relative 输出路径 of the substream file after being transcoded to adaptive bitrate streaming. 如果此参数为空，a relative 路径 in the following 格式 will be used by 默认值：`{inputName}_adaptiveDynamicStreaming_{definition}_{subStreamNumber}.{格式}`。",
															},
															"segment_object_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The relative 输出路径 of the segment file after being transcoded to adaptive bitrate streaming (in HLS 格式 only). 如果此参数为空，a relative 路径 in the following 格式 will be used by 默认值：`{inputName}_adaptiveDynamicStreaming_{definition}_{subStreamNumber}_{segmentNumber}.{格式}`。",
															},
															"add_on_subtitles": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The subtitle file to add.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The 模式 Valid values:`subtitle-stream`: Add a subtitle track.`close-caption-708`: Embed EA-708 subtitles in SEI frames.`close-caption-608`: Embed CEA-608 subtitles in SEI frames.注意：此字段可能返回 null，表示无法获取有效值。",
																		},
																		"subtitle": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The subtitle file.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The input 类型 Valid values:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，this 类型 is only supported for transcoding tasks。",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The information of the COS object to process. This parameter is valid and 必填 when `类型` is `COS`。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The COS 存储桶 of the object to process，such as `TopRankVideo-125xxx88`。",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 地域 of the COS 存储桶，such as `ap-chongqing`。",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 路径 of the object to process，such as `/movie/201907/WildAnimal.mov`。",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The URL of the object to process. This parameter is valid and 必填 when `类型` is `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "URL of a video。",
																								},
																							},
																						},
																					},
																					"s3_input_info": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The information of the AWS S3 object processed. This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"s3_bucket": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The AWS S3 存储桶",
																								},
																								"s3_region": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 地域 of the AWS S3 存储桶",
																								},
																								"s3_object": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 路径 of the AWS S3 object。",
																								},
																								"s3_secret_id": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 键 ID 必填 to access the AWS S3 object。",
																								},
																								"s3_secret_key": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The 键 必填 to access the AWS S3 object。",
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
													Description: "A 内容 moderation task。",
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
													Description: "A 内容 analysis task。",
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
																Description: "An extended parameter，whose 值 is a stringfied JSON.Note: This parameter is for customers with special requirements. It needs to be customized offline.注意：此字段可能返回 null，表示无法获取有效值。",
															},
														},
													},
												},
												"ai_recognition_task": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "A 内容 recognition task。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"definition": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Intelligent video recognition 模板 ID",
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
							Description: "The 存储桶 to save the 输出文件注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The storage 类型 for a media processing 输出文件 Valid values:`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This 类型 is only supported for AWS tasks，and the 输出存储桶 must be in the same 地域 as the 存储桶 of the 来源 file。",
									},
									"cos_output_storage": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The location to save the output object in COS. This parameter is valid and 必填 when `类型` is COS.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bucket": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 存储桶 to which the 输出文件 of media processing is saved，such as `TopRankVideo-125xxx88`. 如果此参数为空，the 值 of the upper layer will be inherited。",
												},
												"region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 地域 of the 输出存储桶，such as `ap-chongqing`. 如果此参数为空，the 值 of the upper layer will be inherited。",
												},
											},
										},
									},
									"s3_output_storage": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The AWS S3 存储桶 to save the 输出文件 This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"s3_bucket": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The AWS S3 存储桶",
												},
												"s3_region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 地域 of the AWS S3 存储桶",
												},
												"s3_secret_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 键 ID 必填 to upload files to the AWS S3 object。",
												},
												"s3_secret_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 键 必填 to upload files to the AWS S3 object。",
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
							Description: "The directory to save the 输出文件注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"task_notify_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The notification configuration.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cmq_model": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CMQ or TDMQ-CMQ model. 有效值：Queue，Topic。",
									},
									"cmq_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CMQ or TDMQ-CMQ 地域，such as `sh` (Shanghai) or `bj` (Beijing)。",
									},
									"topic_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CMQ or TDMQ-CMQ topic to receive notifications. This parameter is valid when `CmqModel` is `Topic`。",
									},
									"queue_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CMQ or TDMQ-CMQ queue to receive notifications. This parameter is valid when `CmqModel` is `Queue`。",
									},
									"notify_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Workflow notification method. 有效值：Finish，Change. 如果此参数为空，`Finish` will be used。",
									},
									"notify_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The notification 类型 Valid values:`CMQ`: This 值 is no longer used. Please use `TDMQ-CMQ` instead.`TDMQ-CMQ`: 消息 queue`URL`: If `NotifyType` is set to `URL`，HTTP callbacks are sent to the URL specified by `NotifyUrl`. HTTP and JSON are 用于the callbacks. The packet 包含response parameters of the `ParseNotification` API.`SCF`: This notification 类型 is not recommended. You need to configure it in the SCF console.`AWS-SQS`: AWS queue. This 类型 is only supported for AWS tasks，and the queue must be in the same 地域 as the AWS 存储桶Note: If you do not pass this parameter or pass in an empty string，`CMQ` will be used. To use a different notification 类型，指定this parameter accordingly。",
									},
									"notify_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "HTTP callback URL，必填 if `NotifyType` is set to `URL`。",
									},
									"aws_sqs": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The AWS SQS queue. This parameter 为必填项 if `NotifyType` is `AWS-SQS`.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"sqs_region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 地域 of the SQS queue。",
												},
												"sqs_queue_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 SQS queue。",
												},
												"s3_secret_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 键 ID 必填 to read from/write to the SQS queue。",
												},
												"s3_secret_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 键 必填 to read from/write to the SQS queue。",
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
							Description: "The 创建时间 in [ISO date 格式](https://intl.cloud.tencent.com/document/product/862/37710?from_cn_redirect=1#52).注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The last updated time in [ISO date 格式](https://intl.cloud.tencent.com/document/product/862/37710?from_cn_redirect=1#52).注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource ID. If there is no associated resource ID，fill it with the 账号's main resource ID。",
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
