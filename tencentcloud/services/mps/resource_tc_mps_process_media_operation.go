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
				Description: "The information of the file to process。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The input 类型 Valid values:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，this 类型 is only supported for transcoding tasks。",
						},
						"cos_input_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The information of the COS object to process. This parameter is valid and 必填 when `类型` is `COS`。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The COS 存储桶 of the object to process，such as `TopRankVideo-125xxx88`。",
									},
									"region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The 地域 of the COS 存储桶，such as `ap-chongqing`。",
									},
									"object": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The 路径 of the object to process，such as `/movie/201907/WildAnimal.mov`。",
									},
								},
							},
						},
						"url_input_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The URL of the object to process. This parameter is valid and 必填 when `类型` is `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "URL of a video。",
									},
								},
							},
						},
						"s3_input_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The information of the AWS S3 object processed. This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"s3_bucket": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The AWS S3 存储桶",
									},
									"s3_region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The 地域 of the AWS S3 存储桶",
									},
									"s3_object": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The 路径 of the AWS S3 object。",
									},
									"s3_secret_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The 键 ID 必填 to access the AWS S3 object。",
									},
									"s3_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The 键 必填 to access the AWS S3 object。",
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
				Description: "The storage location of the media processing 输出文件 如果此参数为空，the storage location in `InputInfo` will be inherited。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The storage 类型 for a media processing 输出文件 Valid values:`COS`: Tencent Cloud COS`&gt;AWS-S3`: AWS S3. This 类型 is only supported for AWS tasks，and the 输出存储桶 must be in the same 地域 as the 存储桶 of the 来源 file。",
						},
						"cos_output_storage": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The location to save the output object in COS. This parameter is valid and 必填 when `类型` is COS.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The 存储桶 to which the 输出文件 of media processing is saved，such as `TopRankVideo-125xxx88`. 如果此参数为空，the 值 of the upper layer will be inherited。",
									},
									"region": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The 地域 of the 输出存储桶，such as `ap-chongqing`. 如果此参数为空，the 值 of the upper layer will be inherited。",
									},
								},
							},
						},
						"s3_output_storage": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The AWS S3 存储桶 to save the 输出文件 This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"s3_bucket": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The AWS S3 存储桶",
									},
									"s3_region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The 地域 of the AWS S3 存储桶",
									},
									"s3_secret_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The 键 ID 必填 to upload files to the AWS S3 object。",
									},
									"s3_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The 键 必填 to upload files to the AWS S3 object。",
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
				Description: "The directory to save the media processing 输出文件，which must start and end with `/`，such as `/movie/201907/`.If you do not 指定this parameter，the file will be saved to the directory specified in `InputInfo`。",
			},

			"schedule_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "The scheme ID.Note 1: About `OutputStorage` and `OutputDir`If an output storage and directory are specified for a subtask of the scheme，those output settings will be applied.If an output storage and directory are not specified for the subtasks of a scheme，the output parameters passed in the `ProcessMedia` API will be applied.Note 2: If `TaskNotifyConfig` is specified，the specified settings will be used instead of the default callback settings of the scheme.Note 3: The trigger configured for a scheme is for automatically starting a scheme. It stops working when you manually call this API to start a scheme。",
			},

			"media_process_task": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The media processing parameters to use。",
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
										Description: "ID video transcoding template。",
									},
									"raw_parameter": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Custom video transcoding parameter，which is valid if `Definition` is 0.This parameter is used in highly customized scenarios. We recommend you use `Definition` to 指定transcoding parameter preferably。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"container": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Container. 有效值：mp4; flv; hls; mp3; flac; ogg; m4a. Among them，mp3，flac，ogg，and m4a are for audio files。",
												},
												"remove_video": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "是否remove video data. Valid values:0: retain;1: remove.默认值：0。",
												},
												"remove_audio": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "是否remove audio data. Valid values:0: retain;1: remove.默认值：0。",
												},
												"video_template": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Video stream configuration parameter. This field 为必填项 when `RemoveVideo` is 0。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"codec": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The video codec. Valid values:`libx264`: H.264`libx265`: H.265`av1`: AOMedia Video 1Note: You must 指定a resolution (not higher than 640 x 480) if the H.265 codec is used.Note: You can only use the AOMedia Video 1 codec for MP4 files。",
															},
															"fps": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The video frame rate (Hz). 取值范围：[0，100].If the 值 is 0，the frame rate will be the same as that of the 来源 video.Note: For adaptive bitrate streaming，the 值 range of this parameter is [0，60]。",
															},
															"bitrate": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The video bitrate (Kbps). 取值范围：0 and [128，35000].If the 值 is 0，the bitrate of the video will be the same as that of the 来源 video。",
															},
															"resolution_adaptive": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Resolution adaption. Valid values:open: 已启用 When resolution adaption is 已启用，`Width` 表示long side of a video，while `Height` 表示short side.close: 已禁用 When resolution adaption is 已禁用，`Width` 表示width of a video，while `Height` 表示height.默认值：open.Note: When resolution adaption is 已启用，`Width` cannot be smaller than `Height`。",
															},
															"width": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum 值 of the width (or long side) of a video stream （像素）。 取值范围：0 and [128，4,096].If both `Width` and `Height` are 0，the resolution will be the same as that of the 来源 video;If `Width` is 0，but `Height` is not 0，`Width` will be proportionally scaled;If `Width` is not 0，but `Height` is 0，`Height` will be proportionally scaled;If both `Width` and `Height` are not 0，the custom resolution will be used.默认值：0。",
															},
															"height": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum 值 of the height (or short side) of a video stream （像素）。 取值范围：0 and [128，4,096].If both `Width` and `Height` are 0，the resolution will be the same as that of the 来源 video;If `Width` is 0，but `Height` is not 0，`Width` will be proportionally scaled;If `Width` is not 0，but `Height` is 0，`Height` will be proportionally scaled;If both `Width` and `Height` are not 0，the custom resolution will be used.默认值：0。",
															},
															"gop": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Frame interval between I keyframes. 取值范围：0 and [1,100000].If this parameter is 0 or left empty，the system will automatically set the GOP length。",
															},
															"fill_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The fill 模式，which 表示how a video is resized when the video's original aspect ratio is different from the target aspect ratio. Valid values:stretch: Stretch the image frame by frame to fill the entire screen. The video image may become squashed or stretched after transcoding.black: Keep the image&#39;s original aspect ratio and fill the blank space with black bars.white: Keep the image's original aspect ratio and fill the blank space with white bars.gauss: Keep the image's original aspect ratio and apply Gaussian blur to the blank space.默认值：black.Note: Only `stretch` and `black` are supported for adaptive bitrate streaming。",
															},
															"vcrf": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The control factor of video constant bitrate. 取值范围：[1，51]If this parameter is specified，CRF (a bitrate control method) will be 用于transcoding. (Video bitrate will no longer take effect.)It is not recommended to 指定this parameter if there are no special requirements。",
															},
														},
													},
												},
												"audio_template": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Audio stream configuration parameter. This field 为必填项 when `RemoveAudio` is 0。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"codec": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Audio stream codec.When the outer `Container` parameter is `mp3`，the valid 值 is:libmp3lame.When the outer `Container` parameter is `ogg` or `flac`，the valid 值 is:flac.When the outer `Container` parameter is `m4a`，the valid values include:libfdk_aac;libmp3lame;ac3.When the outer `Container` parameter is `mp4` or `flv`，the valid values include:libfdk_aac: more suitable for mp4;libmp3lame: more suitable for flv.When the outer `Container` parameter is `hls`，the valid values include:libfdk_aac;libmp3lame。",
															},
															"bitrate": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Audio stream bitrate in Kbps. 取值范围：0 and [26，256].If the 值 is 0，the bitrate of the audio stream will be the same as that of the original audio。",
															},
															"sample_rate": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Audio stream sample rate. Valid values:32,00044,10048,000In Hz。",
															},
															"audio_channel": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Audio channel system. Valid values:1: Mono2: Dual6: StereoWhen the media is packaged in audio 格式 (FLAC，OGG，MP3，M4A)，the sound channel cannot be set to stereo.默认值：2。",
															},
														},
													},
												},
												"tehd_config": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "TESHD transcoding parameter。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "TESHD 类型 Valid values:TEHD-100: TESHD-100.如果此参数为空，TESHD will not be 已启用",
															},
															"max_video_bitrate": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum bitrate，which is valid when `类型` is `TESHD`.如果此参数为空 or 0 is entered，there will be no upper 限制 for bitrate。",
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
										Description: "Video transcoding custom parameter，which is valid when `Definition` is not 0.When any parameters in this structure are entered，they will be 用于override corresponding parameters in templates.This parameter is used in highly customized scenarios. We recommend you only use `Definition` to 指定transcoding parameter.Note: this field may return `null`，indicating that no valid 值 was found。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"container": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Container 格式 有效值：mp4，flv，hls，mp3，flac，ogg，and m4a; mp3，flac，ogg，and m4a are formats of audio files。",
												},
												"remove_video": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "是否remove video data. Valid values:0: retain1: remove。",
												},
												"remove_audio": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "是否remove audio data. Valid values:0: retain1: remove。",
												},
												"video_template": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Video stream configuration parameter。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"codec": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The video codec. Valid values:libx264: H.264libx265: H.265av1: AOMedia Video 1Note: You must 指定a resolution (not higher than 640 x 480) if the H.265 codec is used.Note: You can only use the AOMedia Video 1 codec for MP4 files。",
															},
															"fps": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Video frame rate in Hz. 取值范围：[0，100].If the 值 is 0，the frame rate will be the same as that of the 来源 video。",
															},
															"bitrate": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Bitrate of a video stream in Kbps. 取值范围：0 and [128，35,000].If the 值 is 0，the bitrate of the video will be the same as that of the 来源 video。",
															},
															"resolution_adaptive": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Resolution adaption. Valid values:open: 已启用 When resolution adaption is 已启用，`Width` 表示long side of a video，while `Height` 表示short side.close: 已禁用 When resolution adaption is 已禁用，`Width` 表示width of a video，while `Height` 表示height.Note: When resolution adaption is 已启用，`Width` cannot be smaller than `Height`。",
															},
															"width": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum 值 of the width (or long side) of a video stream （像素）。 取值范围：0 and [128，4,096].If both `Width` and `Height` are 0，the resolution will be the same as that of the 来源 video;If `Width` is 0，but `Height` is not 0，`Width` will be proportionally scaled;If `Width` is not 0，but `Height` is 0，`Height` will be proportionally scaled;If both `Width` and `Height` are not 0，the custom resolution will be used。",
															},
															"height": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum 值 of the height (or short side) of a video stream （像素）。 取值范围：0 and [128，4,096]。",
															},
															"gop": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Frame interval between I keyframes. 取值范围：0 and [1,100000]. If this parameter is 0，the system will automatically set the GOP length。",
															},
															"fill_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Fill 类型 Fill refers to the way of processing a screenshot when its aspect ratio is different from that of the 来源 video. The following fill types are supported: stretch: stretch. The screenshot will be stretched frame by frame to match the aspect ratio of the 来源 video，which may make the screenshot shorter or longer;black: fill with black. This option retains the aspect ratio of the 来源 video for the screenshot and fills the unmatched area with black color blocks.white: fill with white. This option retains the aspect ratio of the 来源 video for the screenshot and fills the unmatched area with white color blocks.gauss: fill with Gaussian blur. This option retains the aspect ratio of the 来源 video for the screenshot and fills the unmatched area with Gaussian blur。",
															},
															"vcrf": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The control factor of video constant bitrate. 取值范围：[0，51]. This parameter will be 已禁用 if you enter `0`.It is not recommended to 指定this parameter if there are no special requirements。",
															},
															"content_adapt_stream": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "是否enable adaptive encoding. Valid values:0: Disable1: Enable默认值：0. If this parameter is set to `1`，multiple streams with different resolutions and bitrates will be generated automatically. The highest resolution，bitrate，and quality of the streams are determined by the values of `width` and `height`，`Bitrate`，and `Vcrf` in `VideoTemplate` respectively. If these parameters are not set in `VideoTemplate`，the highest resolution generated will be the same as that of the 来源 video，and the highest video quality will be close to VMAF 95. To use this parameter or learn about the billing details of adaptive encoding，please contact your sales rep。",
															},
														},
													},
												},
												"audio_template": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Audio stream configuration parameter。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"codec": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Audio stream codec.When the outer `Container` parameter is `mp3`，the valid 值 is:libmp3lame.When the outer `Container` parameter is `ogg` or `flac`，the valid 值 is:flac.When the outer `Container` parameter is `m4a`，the valid values include:libfdk_aac;libmp3lame;ac3.When the outer `Container` parameter is `mp4` or `flv`，the valid values include:libfdk_aac: More suitable for mp4;libmp3lame: More suitable for flv;mp2.When the outer `Container` parameter is `hls`，the valid values include:libfdk_aac;libmp3lame。",
															},
															"bitrate": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Audio stream bitrate in Kbps. 取值范围：0 and [26，256]. If the 值 is 0，the bitrate of the audio stream will be the same as that of the original audio。",
															},
															"sample_rate": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Audio stream sample rate. Valid values:32,00044,10048,000In Hz。",
															},
															"audio_channel": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Audio channel system. Valid values:1: Mono2: Dual6: StereoWhen the media is packaged in audio 格式 (FLAC，OGG，MP3，M4A)，the sound channel cannot be set to stereo。",
															},
															"stream_selects": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeInt,
																},
																Optional:    true,
																Description: "The audio tracks to retain. All audio tracks are retained by default。",
															},
														},
													},
												},
												"tehd_config": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The TSC transcoding parameters.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The TSC 类型 Valid values:`TEHD-100`: TSC-100 (video TSC). `TEHD-200`: TSC-200 (audio TSC). If this parameter is left blank，no modification will be made.注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"max_video_bitrate": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The maximum video bitrate. If this parameter is not specified，no modifications will be made.注意：此字段可能返回 null，表示无法获取有效值。",
															},
														},
													},
												},
												"subtitle_template": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The subtitle settings.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"path": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The URL of the subtitles to add to the video.注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"stream_index": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "The subtitle track to add to the video. If both `路径` and `StreamIndex` are specified，`路径` will be used. You need to 指定at least one of the two parameters.注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"font_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The font. Valid values:`hei.ttf`: Heiti.`song.ttf`: Songti.`simkai.ttf`: Kaiti.`arial.ttf`: Arial.The 默认为 `hei.ttf`.注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"font_size": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The font size (pixels). If this is not specified，the font size in the subtitle file will be used.注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"font_color": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The font color in 0xRRGGBB 格式 默认值：0xFFFFFF (white).注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"font_alpha": {
																Type:        schema.TypeFloat,
																Optional:    true,
																Description: "The text transparency. 取值范围：0-1.`0`: Fully transparent.`1`: Fully opaque.默认值：1.注意：此字段可能返回 null，表示无法获取有效值。",
															},
														},
													},
												},
												"addon_audio_stream": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "The information of the external audio track to add.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The input 类型 Valid values:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，this 类型 is only supported for transcoding tasks。",
															},
															"cos_input_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The information of the COS object to process. This parameter is valid and 必填 when `类型` is `COS`。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bucket": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The COS 存储桶 of the object to process，such as `TopRankVideo-125xxx88`。",
																		},
																		"region": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The 地域 of the COS 存储桶，such as `ap-chongqing`。",
																		},
																		"object": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The 路径 of the object to process，such as `/movie/201907/WildAnimal.mov`。",
																		},
																	},
																},
															},
															"url_input_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The URL of the object to process. This parameter is valid and 必填 when `类型` is `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "URL of a video。",
																		},
																	},
																},
															},
															"s3_input_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The information of the AWS S3 object processed. This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"s3_bucket": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The AWS S3 存储桶",
																		},
																		"s3_region": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The 地域 of the AWS S3 存储桶",
																		},
																		"s3_object": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The 路径 of the AWS S3 object。",
																		},
																		"s3_secret_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The 键 ID 必填 to access the AWS S3 object。",
																		},
																		"s3_secret_key": {
																			Type:        schema.TypeString,
																			Optional:    true,
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
													Optional:    true,
													Description: "An extended field for transcoding.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"add_on_subtitles": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "The subtitle file to add.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The 模式 Valid values:`subtitle-stream`: Add a subtitle track.`close-caption-708`: Embed CEA-708 subtitles in SEI frames.`close-caption-608`: Embed CEA-608 subtitles in SEI frames.注意：此字段可能返回 null，表示无法获取有效值。",
															},
															"subtitle": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The subtitle file.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The input 类型 Valid values:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，this 类型 is only supported for transcoding tasks。",
																		},
																		"cos_input_info": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "The information of the COS object to process. This parameter is valid and 必填 when `类型` is `COS`。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bucket": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "The COS 存储桶 of the object to process，such as `TopRankVideo-125xxx88`。",
																					},
																					"region": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "The 地域 of the COS 存储桶，such as `ap-chongqing`。",
																					},
																					"object": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "The 路径 of the object to process，such as `/movie/201907/WildAnimal.mov`。",
																					},
																				},
																			},
																		},
																		"url_input_info": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "The URL of the object to process. This parameter is valid and 必填 when `类型` is `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"url": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "URL of a video。",
																					},
																				},
																			},
																		},
																		"s3_input_info": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "The information of the AWS S3 object processed. This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"s3_bucket": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "The AWS S3 存储桶",
																					},
																					"s3_region": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "The 地域 of the AWS S3 存储桶",
																					},
																					"s3_object": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "The 路径 of the AWS S3 object。",
																					},
																					"s3_secret_id": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "The 键 ID 必填 to access the AWS S3 object。",
																					},
																					"s3_secret_key": {
																						Type:        schema.TypeString,
																						Optional:    true,
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
										Optional:    true,
										Description: "列表 up to 10 image or text watermarks.注意：此字段可能返回 null，表示无法获取有效值。",
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
													Description: "Custom watermark parameter，which is valid if `Definition` is 0.This parameter is used in highly customized scenarios. We recommend you use `Definition` to 指定watermark parameter preferably.Custom watermark parameter is not available for screenshot。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Watermark 类型 Valid values:image: image watermark。",
															},
															"coordinate_origin": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Origin position，which currently can only be:TopLeft: the origin of coordinates is in the top-left corner of the video，and the origin of the watermark is in the top-left corner of the image or text.默认值：TopLeft。",
															},
															"x_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The horizontal position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %，the `XPos` of the watermark will be the specified percentage of the video width; for example，`10%` means that `XPos` is 10% of the video width;If the string ends in px，the `XPos` of the watermark will be the specified px; for example，`100px` means that `XPos` is 100 px.默认值：0 px。",
															},
															"y_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The vertical position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %，the `YPos` of the watermark will be the specified percentage of the video height; for example，`10%` means that `YPos` is 10% of the video height;If the string ends in px，the `YPos` of the watermark will be the specified px; for example，`100px` means that `YPos` is 100 px.默认值：0 px。",
															},
															"image_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Image watermark template. This field 为必填项 when `类型` is `image` and is invalid when `类型` is `text`。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"image_content": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Required:    true,
																			Description: "Input 内容 of watermark image. JPEG and PNG images are supported。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "The input 类型 Valid values:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，this 类型 is only supported for transcoding tasks。",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "The information of the COS object to process. This parameter is valid and 必填 when `类型` is `COS`。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The COS 存储桶 of the object to process，such as `TopRankVideo-125xxx88`。",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The 地域 of the COS 存储桶，such as `ap-chongqing`。",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The 路径 of the object to process，such as `/movie/201907/WildAnimal.mov`。",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "The URL of the object to process. This parameter is valid and 必填 when `类型` is `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "URL of a video。",
																								},
																							},
																						},
																					},
																					"s3_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "The information of the AWS S3 object processed. This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"s3_bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The AWS S3 存储桶",
																								},
																								"s3_region": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The 地域 of the AWS S3 存储桶",
																								},
																								"s3_object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The 路径 of the AWS S3 object。",
																								},
																								"s3_secret_id": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "The 键 ID 必填 to access the AWS S3 object。",
																								},
																								"s3_secret_key": {
																									Type:        schema.TypeString,
																									Optional:    true,
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
																			Optional:    true,
																			Description: "Watermark width. % and px formats are supported:If the string ends in %，the `Width` of the watermark will be the specified percentage of the video width; for example，`10%` means that `Width` is 10% of the video width;If the string ends in px，the `Width` of the watermark will be in px; for example，`100px` means that `Width` is 100 px.默认值：10%。",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Watermark height. % and px formats are supported:If the string ends in %，the `Height` of the watermark will be the specified percentage of the video height; for example，`10%` means that `Height` is 10% of the video height;If the string ends in px，the `Height` of the watermark will be in px; for example，`100px` means that `Height` is 100 px.默认值：0 px，which means that `Height` will be proportionally scaled according to the aspect ratio of the original watermark image。",
																		},
																		"repeat_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
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
													Optional:    true,
													Description: "Text 内容 of up to 100 characters. This field 为必填项 only when the watermark 类型 is text.Text watermark is not available for screenshot。",
												},
												"svg_content": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "SVG 内容 of up to 2,000,000 characters. This field 为必填项 only when the watermark 类型 is `SVG`.SVG watermark is not available for screenshot。",
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "开始时间 偏移量 of a watermark （秒）。 如果此参数为空 or 0 is entered，the watermark will appear upon the first video frame.如果此参数为空 or 0 is entered，the watermark will appear upon the first video frame;If this 值 is greater than 0 (e.g.，n)，the watermark will appear at second n after the first video frame;If this 值 is smaller than 0 (e.g.，-n)，the watermark will appear at second n before the last video frame。",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "结束时间 偏移量 of a watermark （秒）。如果此参数为空 or 0 is entered，the watermark will exist till the last video frame;If this 值 is greater than 0 (e.g.，n)，the watermark will exist till second n;If this 值 is smaller than 0 (e.g.，-n)，the watermark will exist till second n before the last video frame。",
												},
											},
										},
									},
									"mosaic_set": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "列表 blurs. Up to 10 ones can be supported。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"coordinate_origin": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Origin position，which currently can only be:TopLeft: the origin of coordinates is in the top-left corner of the video，and the origin of the blur is in the top-left corner of the image or text.默认值：TopLeft。",
												},
												"x_pos": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The horizontal position of the origin of the blur relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %，the `XPos` of the blur will be the specified percentage of the video width; for example，`10%` means that `XPos` is 10% of the video width;If the string ends in px，the `XPos` of the blur will be the specified px; for example，`100px` means that `XPos` is 100 px.默认值：0 px。",
												},
												"y_pos": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Vertical position of the origin of blur relative to the origin of coordinates of video. % and px formats are supported:If the string ends in %，the `YPos` of the blur will be the specified percentage of the video height; for example，`10%` means that `YPos` is 10% of the video height;If the string ends in px，the `YPos` of the blur will be the specified px; for example，`100px` means that `YPos` is 100 px.默认值：0 px。",
												},
												"width": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Blur width. % and px formats are supported:If the string ends in %，the `Width` of the blur will be the specified percentage of the video width; for example，`10%` means that `Width` is 10% of the video width;If the string ends in px，the `Width` of the blur will be in px; for example，`100px` means that `Width` is 100 px.默认值：10%。",
												},
												"height": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Blur height. % and px formats are supported:If the string ends in %，the `Height` of the blur will be the specified percentage of the video height; for example，`10%` means that `Height` is 10% of the video height;If the string ends in px，the `Height` of the blur will be in px; for example，`100px` means that `Height` is 100 px.默认值：10%。",
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "开始时间 偏移量 of blur （秒）。 如果此参数为空 or 0 is entered，the blur will appear upon the first video frame.如果此参数为空 or 0 is entered，the blur will appear upon the first video frame;If this 值 is greater than 0 (e.g.，n)，the blur will appear at second n after the first video frame;If this 值 is smaller than 0 (e.g.，-n)，the blur will appear at second n before the last video frame。",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "结束时间 偏移量 of blur （秒）。如果此参数为空 or 0 is entered，the blur will exist till the last video frame;If this 值 is greater than 0 (e.g.，n)，the blur will exist till second n;If this 值 is smaller than 0 (e.g.，-n)，the blur will exist till second n before the last video frame。",
												},
											},
										},
									},
									"start_time_offset": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "开始时间 偏移量 of a transcoded video，（秒）。如果此参数为空 or set to 0，the transcoded video will start at the same time as the original video.If this parameter is set to a positive number (n for example)，the transcoded video will start at the nth second of the original video.If this parameter is set to a negative number (-n for example)，the transcoded video will start at the nth second before the end of the original video。",
									},
									"end_time_offset": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "结束时间 偏移量 of a transcoded video，（秒）。如果此参数为空 or set to 0，the transcoded video will end at the same time as the original video.If this parameter is set to a positive number (n for example)，the transcoded video will end at the nth second of the original video.If this parameter is set to a negative number (-n for example)，the transcoded video will end at the nth second before the end of the original video。",
									},
									"output_storage": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Target 存储桶 of an 输出文件 如果此参数为空，the `OutputStorage` 值 of the upper folder will be inherited.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The storage 类型 for a media processing 输出文件 Valid values:`COS`: Tencent Cloud COS`&gt;AWS-S3`: AWS S3. This 类型 is only supported for AWS tasks，and the 输出存储桶 must be in the same 地域 as the 存储桶 of the 来源 file。",
												},
												"cos_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The location to save the output object in COS. This parameter is valid and 必填 when `类型` is COS.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bucket": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The 存储桶 to which the 输出文件 of media processing is saved，such as `TopRankVideo-125xxx88`. 如果此参数为空，the 值 of the upper layer will be inherited。",
															},
															"region": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The 地域 of the 输出存储桶，such as `ap-chongqing`. 如果此参数为空，the 值 of the upper layer will be inherited。",
															},
														},
													},
												},
												"s3_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The AWS S3 存储桶 to save the 输出文件 This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"s3_bucket": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The AWS S3 存储桶",
															},
															"s3_region": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The 地域 of the AWS S3 存储桶",
															},
															"s3_secret_id": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The 键 ID 必填 to upload files to the AWS S3 object。",
															},
															"s3_secret_key": {
																Type:        schema.TypeString,
																Optional:    true,
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
										Optional:    true,
										Description: "路径 to a primary 输出文件，which can be a relative 路径 or an absolute 路径 如果此参数为空，the following relative 路径 will be used by 默认值：`{inputName}_transcode_{definition}.{格式}`。",
									},
									"segment_object_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "路径 to an 输出文件 part (the 路径 to ts during transcoding to HLS)，which can only be a relative 路径 如果此参数为空，the following relative 路径 will be used by 默认值：`{inputName}_transcode_{definition}_{number}.{格式}`。",
									},
									"object_number_format": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Rule of the `{number}` variable in the 输出路径 after transcoding.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"initial_value": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Start 值 of the `{number}` variable. 默认值：0。",
												},
												"increment": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Increment of the `{number}` variable. 默认值：1。",
												},
												"min_length": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "最小长度the `{number}` variable. A placeholder will be used if the variable length is below the minimum requirement. 默认值：1。",
												},
												"place_holder": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Placeholder used when the `{number}` variable length is below the minimum requirement. 默认值：0。",
												},
											},
										},
									},
									"head_tail_parameter": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Opening and closing credits parametersNote: this field may return `null`，indicating that no valid 值 was found。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"head_set": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Opening credits list。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The input 类型 Valid values:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，this 类型 is only supported for transcoding tasks。",
															},
															"cos_input_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The information of the COS object to process. This parameter is valid and 必填 when `类型` is `COS`。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bucket": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The COS 存储桶 of the object to process，such as `TopRankVideo-125xxx88`。",
																		},
																		"region": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The 地域 of the COS 存储桶，such as `ap-chongqing`。",
																		},
																		"object": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The 路径 of the object to process，such as `/movie/201907/WildAnimal.mov`。",
																		},
																	},
																},
															},
															"url_input_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The URL of the object to process. This parameter is valid and 必填 when `类型` is `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "URL of a video。",
																		},
																	},
																},
															},
															"s3_input_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The information of the AWS S3 object processed. This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"s3_bucket": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The AWS S3 存储桶",
																		},
																		"s3_region": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The 地域 of the AWS S3 存储桶",
																		},
																		"s3_object": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The 路径 of the AWS S3 object。",
																		},
																		"s3_secret_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The 键 ID 必填 to access the AWS S3 object。",
																		},
																		"s3_secret_key": {
																			Type:        schema.TypeString,
																			Optional:    true,
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
													Optional:    true,
													Description: "Closing credits list。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The input 类型 Valid values:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，this 类型 is only supported for transcoding tasks。",
															},
															"cos_input_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The information of the COS object to process. This parameter is valid and 必填 when `类型` is `COS`。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bucket": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The COS 存储桶 of the object to process，such as `TopRankVideo-125xxx88`。",
																		},
																		"region": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The 地域 of the COS 存储桶，such as `ap-chongqing`。",
																		},
																		"object": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The 路径 of the object to process，such as `/movie/201907/WildAnimal.mov`。",
																		},
																	},
																},
															},
															"url_input_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The URL of the object to process. This parameter is valid and 必填 when `类型` is `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "URL of a video。",
																		},
																	},
																},
															},
															"s3_input_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The information of the AWS S3 object processed. This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"s3_bucket": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The AWS S3 存储桶",
																		},
																		"s3_region": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The 地域 of the AWS S3 存储桶",
																		},
																		"s3_object": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The 路径 of the AWS S3 object。",
																		},
																		"s3_secret_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The 键 ID 必填 to access the AWS S3 object。",
																		},
																		"s3_secret_key": {
																			Type:        schema.TypeString,
																			Optional:    true,
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
						"animated_graphic_task_set": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 animated image generating tasks。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Animated image generating 模板 ID",
									},
									"start_time_offset": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: "开始时间 of an animated image in a video （秒）。",
									},
									"end_time_offset": {
										Type:        schema.TypeFloat,
										Required:    true,
										Description: "结束时间 of an animated image in a video （秒）。",
									},
									"output_storage": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Target 存储桶 of a generated animated image file. 如果此参数为空，the `OutputStorage` 值 of the upper folder will be inherited.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The storage 类型 for a media processing 输出文件 Valid values:`COS`: Tencent Cloud COS`&gt;AWS-S3`: AWS S3. This 类型 is only supported for AWS tasks，and the 输出存储桶 must be in the same 地域 as the 存储桶 of the 来源 file。",
												},
												"cos_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The location to save the output object in COS. This parameter is valid and 必填 when `类型` is COS.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bucket": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The 存储桶 to which the 输出文件 of media processing is saved，such as `TopRankVideo-125xxx88`. 如果此参数为空，the 值 of the upper layer will be inherited。",
															},
															"region": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The 地域 of the 输出存储桶，such as `ap-chongqing`. 如果此参数为空，the 值 of the upper layer will be inherited。",
															},
														},
													},
												},
												"s3_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The AWS S3 存储桶 to save the 输出文件 This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"s3_bucket": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The AWS S3 存储桶",
															},
															"s3_region": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The 地域 of the AWS S3 存储桶",
															},
															"s3_secret_id": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The 键 ID 必填 to upload files to the AWS S3 object。",
															},
															"s3_secret_key": {
																Type:        schema.TypeString,
																Optional:    true,
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
										Optional:    true,
										Description: "输出路径 to a generated animated image file，which can be a relative 路径 or an absolute 路径 如果此参数为空，the following relative 路径 will be used by 默认值：`{inputName}_animatedGraphic_{definition}.{格式}`。",
									},
								},
							},
						},
						"snapshot_by_time_offset_task_set": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 time point screencapturing tasks。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "ID time point screencapturing template。",
									},
									"ext_time_offset_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "列表 screenshot time points in the 格式 of `s` or `%`:If the string ends in `s`，it means that the time point is in seconds; for example，`3.5s` means that the time point is the 3.5th second;If the string ends in `%`，it means that the time point is the specified percentage of the video duration; for example，`10%` means that the time point is 10% of the video duration。",
									},
									"time_offset_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeFloat,
										},
										Optional:    true,
										Description: "列表 time points of screenshots in &lt;font color=red&gt;seconds&lt;/font&gt;。",
									},
									"watermark_set": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "列表 up to 10 image or text watermarks.注意：此字段可能返回 null，表示无法获取有效值。",
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
													Description: "Custom watermark parameter，which is valid if `Definition` is 0.This parameter is used in highly customized scenarios. We recommend you use `Definition` to 指定watermark parameter preferably.Custom watermark parameter is not available for screenshot。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Watermark 类型 Valid values:image: image watermark。",
															},
															"coordinate_origin": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Origin position，which currently can only be:TopLeft: the origin of coordinates is in the top-left corner of the video，and the origin of the watermark is in the top-left corner of the image or text.默认值：TopLeft。",
															},
															"x_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The horizontal position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %，the `XPos` of the watermark will be the specified percentage of the video width; for example，`10%` means that `XPos` is 10% of the video width;If the string ends in px，the `XPos` of the watermark will be the specified px; for example，`100px` means that `XPos` is 100 px.默认值：0 px。",
															},
															"y_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The vertical position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %，the `YPos` of the watermark will be the specified percentage of the video height; for example，`10%` means that `YPos` is 10% of the video height;If the string ends in px，the `YPos` of the watermark will be the specified px; for example，`100px` means that `YPos` is 100 px.默认值：0 px。",
															},
															"image_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Image watermark template. This field 为必填项 when `类型` is `image` and is invalid when `类型` is `text`。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"image_content": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Required:    true,
																			Description: "Input 内容 of watermark image. JPEG and PNG images are supported。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "The input 类型 Valid values:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，this 类型 is only supported for transcoding tasks。",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "The information of the COS object to process. This parameter is valid and 必填 when `类型` is `COS`。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The COS 存储桶 of the object to process，such as `TopRankVideo-125xxx88`。",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The 地域 of the COS 存储桶，such as `ap-chongqing`。",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The 路径 of the object to process，such as `/movie/201907/WildAnimal.mov`。",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "The URL of the object to process. This parameter is valid and 必填 when `类型` is `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "URL of a video。",
																								},
																							},
																						},
																					},
																					"s3_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "The information of the AWS S3 object processed. This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"s3_bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The AWS S3 存储桶",
																								},
																								"s3_region": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The 地域 of the AWS S3 存储桶",
																								},
																								"s3_object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The 路径 of the AWS S3 object。",
																								},
																								"s3_secret_id": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "The 键 ID 必填 to access the AWS S3 object。",
																								},
																								"s3_secret_key": {
																									Type:        schema.TypeString,
																									Optional:    true,
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
																			Optional:    true,
																			Description: "Watermark width. % and px formats are supported:If the string ends in %，the `Width` of the watermark will be the specified percentage of the video width; for example，`10%` means that `Width` is 10% of the video width;If the string ends in px，the `Width` of the watermark will be in px; for example，`100px` means that `Width` is 100 px.默认值：10%。",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Watermark height. % and px formats are supported:If the string ends in %，the `Height` of the watermark will be the specified percentage of the video height; for example，`10%` means that `Height` is 10% of the video height;If the string ends in px，the `Height` of the watermark will be in px; for example，`100px` means that `Height` is 100 px.默认值：0 px，which means that `Height` will be proportionally scaled according to the aspect ratio of the original watermark image。",
																		},
																		"repeat_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
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
													Optional:    true,
													Description: "Text 内容 of up to 100 characters. This field 为必填项 only when the watermark 类型 is text.Text watermark is not available for screenshot。",
												},
												"svg_content": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "SVG 内容 of up to 2,000,000 characters. This field 为必填项 only when the watermark 类型 is `SVG`.SVG watermark is not available for screenshot。",
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "开始时间 偏移量 of a watermark （秒）。 如果此参数为空 or 0 is entered，the watermark will appear upon the first video frame.如果此参数为空 or 0 is entered，the watermark will appear upon the first video frame;If this 值 is greater than 0 (e.g.，n)，the watermark will appear at second n after the first video frame;If this 值 is smaller than 0 (e.g.，-n)，the watermark will appear at second n before the last video frame。",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "结束时间 偏移量 of a watermark （秒）。如果此参数为空 or 0 is entered，the watermark will exist till the last video frame;If this 值 is greater than 0 (e.g.，n)，the watermark will exist till second n;If this 值 is smaller than 0 (e.g.，-n)，the watermark will exist till second n before the last video frame。",
												},
											},
										},
									},
									"output_storage": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Target 存储桶 of a generated time point screenshot file. 如果此参数为空，the `OutputStorage` 值 of the upper folder will be inherited.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The storage 类型 for a media processing 输出文件 Valid values:`COS`: Tencent Cloud COS`&gt;AWS-S3`: AWS S3. This 类型 is only supported for AWS tasks，and the 输出存储桶 must be in the same 地域 as the 存储桶 of the 来源 file。",
												},
												"cos_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The location to save the output object in COS. This parameter is valid and 必填 when `类型` is COS.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bucket": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The 存储桶 to which the 输出文件 of media processing is saved，such as `TopRankVideo-125xxx88`. 如果此参数为空，the 值 of the upper layer will be inherited。",
															},
															"region": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The 地域 of the 输出存储桶，such as `ap-chongqing`. 如果此参数为空，the 值 of the upper layer will be inherited。",
															},
														},
													},
												},
												"s3_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The AWS S3 存储桶 to save the 输出文件 This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"s3_bucket": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The AWS S3 存储桶",
															},
															"s3_region": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The 地域 of the AWS S3 存储桶",
															},
															"s3_secret_id": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The 键 ID 必填 to upload files to the AWS S3 object。",
															},
															"s3_secret_key": {
																Type:        schema.TypeString,
																Optional:    true,
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
										Optional:    true,
										Description: "输出路径 to a generated time point screenshot，which can be a relative 路径 or an absolute 路径 如果此参数为空，the following relative 路径 will be used by 默认值：`{inputName}_snapshotByTimeOffset_{definition}_{number}.{格式}`。",
									},
									"object_number_format": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Rule of the `{number}` variable in the time point screenshot 输出路径注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"initial_value": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Start 值 of the `{number}` variable. 默认值：0。",
												},
												"increment": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Increment of the `{number}` variable. 默认值：1。",
												},
												"min_length": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "最小长度the `{number}` variable. A placeholder will be used if the variable length is below the minimum requirement. 默认值：1。",
												},
												"place_holder": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Placeholder used when the `{number}` variable length is below the minimum requirement. 默认值：0。",
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
										Description: "列表 up to 10 image or text watermarks.注意：此字段可能返回 null，表示无法获取有效值。",
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
													Description: "Custom watermark parameter，which is valid if `Definition` is 0.This parameter is used in highly customized scenarios. We recommend you use `Definition` to 指定watermark parameter preferably.Custom watermark parameter is not available for screenshot。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Watermark 类型 Valid values:image: image watermark。",
															},
															"coordinate_origin": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Origin position，which currently can only be:TopLeft: the origin of coordinates is in the top-left corner of the video，and the origin of the watermark is in the top-left corner of the image or text.默认值：TopLeft。",
															},
															"x_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The horizontal position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %，the `XPos` of the watermark will be the specified percentage of the video width; for example，`10%` means that `XPos` is 10% of the video width;If the string ends in px，the `XPos` of the watermark will be the specified px; for example，`100px` means that `XPos` is 100 px.默认值：0 px。",
															},
															"y_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The vertical position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %，the `YPos` of the watermark will be the specified percentage of the video height; for example，`10%` means that `YPos` is 10% of the video height;If the string ends in px，the `YPos` of the watermark will be the specified px; for example，`100px` means that `YPos` is 100 px.默认值：0 px。",
															},
															"image_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Image watermark template. This field 为必填项 when `类型` is `image` and is invalid when `类型` is `text`。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"image_content": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Required:    true,
																			Description: "Input 内容 of watermark image. JPEG and PNG images are supported。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "The input 类型 Valid values:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，this 类型 is only supported for transcoding tasks。",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "The information of the COS object to process. This parameter is valid and 必填 when `类型` is `COS`。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The COS 存储桶 of the object to process，such as `TopRankVideo-125xxx88`。",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The 地域 of the COS 存储桶，such as `ap-chongqing`。",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The 路径 of the object to process，such as `/movie/201907/WildAnimal.mov`。",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "The URL of the object to process. This parameter is valid and 必填 when `类型` is `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "URL of a video。",
																								},
																							},
																						},
																					},
																					"s3_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "The information of the AWS S3 object processed. This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"s3_bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The AWS S3 存储桶",
																								},
																								"s3_region": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The 地域 of the AWS S3 存储桶",
																								},
																								"s3_object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The 路径 of the AWS S3 object。",
																								},
																								"s3_secret_id": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "The 键 ID 必填 to access the AWS S3 object。",
																								},
																								"s3_secret_key": {
																									Type:        schema.TypeString,
																									Optional:    true,
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
																			Optional:    true,
																			Description: "Watermark width. % and px formats are supported:If the string ends in %，the `Width` of the watermark will be the specified percentage of the video width; for example，`10%` means that `Width` is 10% of the video width;If the string ends in px，the `Width` of the watermark will be in px; for example，`100px` means that `Width` is 100 px.默认值：10%。",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Watermark height. % and px formats are supported:If the string ends in %，the `Height` of the watermark will be the specified percentage of the video height; for example，`10%` means that `Height` is 10% of the video height;If the string ends in px，the `Height` of the watermark will be in px; for example，`100px` means that `Height` is 100 px.默认值：0 px，which means that `Height` will be proportionally scaled according to the aspect ratio of the original watermark image。",
																		},
																		"repeat_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
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
													Optional:    true,
													Description: "Text 内容 of up to 100 characters. This field 为必填项 only when the watermark 类型 is text.Text watermark is not available for screenshot。",
												},
												"svg_content": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "SVG 内容 of up to 2,000,000 characters. This field 为必填项 only when the watermark 类型 is `SVG`.SVG watermark is not available for screenshot。",
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "开始时间 偏移量 of a watermark （秒）。 如果此参数为空 or 0 is entered，the watermark will appear upon the first video frame.如果此参数为空 or 0 is entered，the watermark will appear upon the first video frame;If this 值 is greater than 0 (e.g.，n)，the watermark will appear at second n after the first video frame;If this 值 is smaller than 0 (e.g.，-n)，the watermark will appear at second n before the last video frame。",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "结束时间 偏移量 of a watermark （秒）。如果此参数为空 or 0 is entered，the watermark will exist till the last video frame;If this 值 is greater than 0 (e.g.，n)，the watermark will exist till second n;If this 值 is smaller than 0 (e.g.，-n)，the watermark will exist till second n before the last video frame。",
												},
											},
										},
									},
									"output_storage": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Target 存储桶 of a sampled screenshot. 如果此参数为空，the `OutputStorage` 值 of the upper folder will be inherited.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The storage 类型 for a media processing 输出文件 Valid values:`COS`: Tencent Cloud COS`&gt;AWS-S3`: AWS S3. This 类型 is only supported for AWS tasks，and the 输出存储桶 must be in the same 地域 as the 存储桶 of the 来源 file。",
												},
												"cos_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The location to save the output object in COS. This parameter is valid and 必填 when `类型` is COS.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bucket": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The 存储桶 to which the 输出文件 of media processing is saved，such as `TopRankVideo-125xxx88`. 如果此参数为空，the 值 of the upper layer will be inherited。",
															},
															"region": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The 地域 of the 输出存储桶，such as `ap-chongqing`. 如果此参数为空，the 值 of the upper layer will be inherited。",
															},
														},
													},
												},
												"s3_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The AWS S3 存储桶 to save the 输出文件 This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"s3_bucket": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The AWS S3 存储桶",
															},
															"s3_region": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The 地域 of the AWS S3 存储桶",
															},
															"s3_secret_id": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The 键 ID 必填 to upload files to the AWS S3 object。",
															},
															"s3_secret_key": {
																Type:        schema.TypeString,
																Optional:    true,
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
										Optional:    true,
										Description: "输出路径 to a generated sampled screenshot，which can be a relative 路径 or an absolute 路径 如果此参数为空，the following relative 路径 will be used by 默认值：`{inputName}_sampleSnapshot_{definition}_{number}.{格式}`。",
									},
									"object_number_format": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Rule of the `{number}` variable in the sampled screenshot 输出路径注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"initial_value": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Start 值 of the `{number}` variable. 默认值：0。",
												},
												"increment": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Increment of the `{number}` variable. 默认值：1。",
												},
												"min_length": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "最小长度the `{number}` variable. A placeholder will be used if the variable length is below the minimum requirement. 默认值：1。",
												},
												"place_holder": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Placeholder used when the `{number}` variable length is below the minimum requirement. 默认值：0。",
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
							Description: "列表 image sprite generating tasks。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "ID an image sprite generating template。",
									},
									"output_storage": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Target 存储桶 of a generated image sprite. 如果此参数为空，the `OutputStorage` 值 of the upper folder will be inherited.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The storage 类型 for a media processing 输出文件 Valid values:`COS`: Tencent Cloud COS`&gt;AWS-S3`: AWS S3. This 类型 is only supported for AWS tasks，and the 输出存储桶 must be in the same 地域 as the 存储桶 of the 来源 file。",
												},
												"cos_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The location to save the output object in COS. This parameter is valid and 必填 when `类型` is COS.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bucket": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The 存储桶 to which the 输出文件 of media processing is saved，such as `TopRankVideo-125xxx88`. 如果此参数为空，the 值 of the upper layer will be inherited。",
															},
															"region": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The 地域 of the 输出存储桶，such as `ap-chongqing`. 如果此参数为空，the 值 of the upper layer will be inherited。",
															},
														},
													},
												},
												"s3_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The AWS S3 存储桶 to save the 输出文件 This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"s3_bucket": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The AWS S3 存储桶",
															},
															"s3_region": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The 地域 of the AWS S3 存储桶",
															},
															"s3_secret_id": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The 键 ID 必填 to upload files to the AWS S3 object。",
															},
															"s3_secret_key": {
																Type:        schema.TypeString,
																Optional:    true,
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
										Optional:    true,
										Description: "输出路径 to a generated image sprite file，which can be a relative 路径 or an absolute 路径 如果此参数为空，the following relative 路径 will be used by 默认值：`{inputName}_imageSprite_{definition}_{number}.{格式}`。",
									},
									"web_vtt_object_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "输出路径 to the WebVTT file after an image sprite is generated，which can only be a relative 路径 如果此参数为空，the following relative 路径 will be used by 默认值：`{inputName}_imageSprite_{definition}.{格式}`。",
									},
									"object_number_format": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Rule of the `{number}` variable in the image sprite 输出路径注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"initial_value": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Start 值 of the `{number}` variable. 默认值：0。",
												},
												"increment": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Increment of the `{number}` variable. 默认值：1。",
												},
												"min_length": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "最小长度the `{number}` variable. A placeholder will be used if the variable length is below the minimum requirement. 默认值：1。",
												},
												"place_holder": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Placeholder used when the `{number}` variable length is below the minimum requirement. 默认值：0。",
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
										Description: "列表 up to 10 image or text watermarks。",
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
													Description: "Custom watermark parameter，which is valid if `Definition` is 0.This parameter is used in highly customized scenarios. We recommend you use `Definition` to 指定watermark parameter preferably.Custom watermark parameter is not available for screenshot。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Watermark 类型 Valid values:image: image watermark。",
															},
															"coordinate_origin": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Origin position，which currently can only be:TopLeft: the origin of coordinates is in the top-left corner of the video，and the origin of the watermark is in the top-left corner of the image or text.默认值：TopLeft。",
															},
															"x_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The horizontal position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %，the `XPos` of the watermark will be the specified percentage of the video width; for example，`10%` means that `XPos` is 10% of the video width;If the string ends in px，the `XPos` of the watermark will be the specified px; for example，`100px` means that `XPos` is 100 px.默认值：0 px。",
															},
															"y_pos": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The vertical position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported:If the string ends in %，the `YPos` of the watermark will be the specified percentage of the video height; for example，`10%` means that `YPos` is 10% of the video height;If the string ends in px，the `YPos` of the watermark will be the specified px; for example，`100px` means that `YPos` is 100 px.默认值：0 px。",
															},
															"image_template": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Image watermark template. This field 为必填项 when `类型` is `image` and is invalid when `类型` is `text`。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"image_content": {
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Required:    true,
																			Description: "Input 内容 of watermark image. JPEG and PNG images are supported。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"type": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "The input 类型 Valid values:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，this 类型 is only supported for transcoding tasks。",
																					},
																					"cos_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "The information of the COS object to process. This parameter is valid and 必填 when `类型` is `COS`。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The COS 存储桶 of the object to process，such as `TopRankVideo-125xxx88`。",
																								},
																								"region": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The 地域 of the COS 存储桶，such as `ap-chongqing`。",
																								},
																								"object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The 路径 of the object to process，such as `/movie/201907/WildAnimal.mov`。",
																								},
																							},
																						},
																					},
																					"url_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "The URL of the object to process. This parameter is valid and 必填 when `类型` is `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "URL of a video。",
																								},
																							},
																						},
																					},
																					"s3_input_info": {
																						Type:        schema.TypeList,
																						MaxItems:    1,
																						Optional:    true,
																						Description: "The information of the AWS S3 object processed. This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"s3_bucket": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The AWS S3 存储桶",
																								},
																								"s3_region": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The 地域 of the AWS S3 存储桶",
																								},
																								"s3_object": {
																									Type:        schema.TypeString,
																									Required:    true,
																									Description: "The 路径 of the AWS S3 object。",
																								},
																								"s3_secret_id": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "The 键 ID 必填 to access the AWS S3 object。",
																								},
																								"s3_secret_key": {
																									Type:        schema.TypeString,
																									Optional:    true,
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
																			Optional:    true,
																			Description: "Watermark width. % and px formats are supported:If the string ends in %，the `Width` of the watermark will be the specified percentage of the video width; for example，`10%` means that `Width` is 10% of the video width;If the string ends in px，the `Width` of the watermark will be in px; for example，`100px` means that `Width` is 100 px.默认值：10%。",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Watermark height. % and px formats are supported:If the string ends in %，the `Height` of the watermark will be the specified percentage of the video height; for example，`10%` means that `Height` is 10% of the video height;If the string ends in px，the `Height` of the watermark will be in px; for example，`100px` means that `Height` is 100 px.默认值：0 px，which means that `Height` will be proportionally scaled according to the aspect ratio of the original watermark image。",
																		},
																		"repeat_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
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
													Optional:    true,
													Description: "Text 内容 of up to 100 characters. This field 为必填项 only when the watermark 类型 is text.Text watermark is not available for screenshot。",
												},
												"svg_content": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "SVG 内容 of up to 2,000,000 characters. This field 为必填项 only when the watermark 类型 is `SVG`.SVG watermark is not available for screenshot。",
												},
												"start_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "开始时间 偏移量 of a watermark （秒）。 如果此参数为空 or 0 is entered，the watermark will appear upon the first video frame.如果此参数为空 or 0 is entered，the watermark will appear upon the first video frame;If this 值 is greater than 0 (e.g.，n)，the watermark will appear at second n after the first video frame;If this 值 is smaller than 0 (e.g.，-n)，the watermark will appear at second n before the last video frame。",
												},
												"end_time_offset": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "结束时间 偏移量 of a watermark （秒）。如果此参数为空 or 0 is entered，the watermark will exist till the last video frame;If this 值 is greater than 0 (e.g.，n)，the watermark will exist till second n;If this 值 is smaller than 0 (e.g.，-n)，the watermark will exist till second n before the last video frame。",
												},
											},
										},
									},
									"output_storage": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Target 存储桶 of an 输出文件 after being transcoded to adaptive bitrate streaming. 如果此参数为空，the `OutputStorage` 值 of the upper folder will be inherited.Note: this field may return null，indicating that no valid values can be obtained。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The storage 类型 for a media processing 输出文件 Valid values:`COS`: Tencent Cloud COS`&gt;AWS-S3`: AWS S3. This 类型 is only supported for AWS tasks，and the 输出存储桶 must be in the same 地域 as the 存储桶 of the 来源 file。",
												},
												"cos_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The location to save the output object in COS. This parameter is valid and 必填 when `类型` is COS.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bucket": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The 存储桶 to which the 输出文件 of media processing is saved，such as `TopRankVideo-125xxx88`. 如果此参数为空，the 值 of the upper layer will be inherited。",
															},
															"region": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The 地域 of the 输出存储桶，such as `ap-chongqing`. 如果此参数为空，the 值 of the upper layer will be inherited。",
															},
														},
													},
												},
												"s3_output_storage": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The AWS S3 存储桶 to save the 输出文件 This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"s3_bucket": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The AWS S3 存储桶",
															},
															"s3_region": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The 地域 of the AWS S3 存储桶",
															},
															"s3_secret_id": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The 键 ID 必填 to upload files to the AWS S3 object。",
															},
															"s3_secret_key": {
																Type:        schema.TypeString,
																Optional:    true,
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
										Optional:    true,
										Description: "The relative or absolute 输出路径 of the manifest file after being transcoded to adaptive bitrate streaming. 如果此参数为空，a relative 路径 in the following 格式 will be used by 默认值：`{inputName}_adaptiveDynamicStreaming_{definition}.{格式}`。",
									},
									"sub_stream_object_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The relative 输出路径 of the substream file after being transcoded to adaptive bitrate streaming. 如果此参数为空，a relative 路径 in the following 格式 will be used by 默认值：`{inputName}_adaptiveDynamicStreaming_{definition}_{subStreamNumber}.{格式}`。",
									},
									"segment_object_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The relative 输出路径 of the segment file after being transcoded to adaptive bitrate streaming (in HLS 格式 only). 如果此参数为空，a relative 路径 in the following 格式 will be used by 默认值：`{inputName}_adaptiveDynamicStreaming_{definition}_{subStreamNumber}_{segmentNumber}.{格式}`。",
									},
									"add_on_subtitles": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The subtitle file to add.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The 模式 Valid values:`subtitle-stream`: Add a subtitle track.`close-caption-708`: Embed CEA-708 subtitles in SEI frames.`close-caption-608`: Embed CEA-608 subtitles in SEI frames.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"subtitle": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "The subtitle file.注意：此字段可能返回 null，表示无法获取有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "The input 类型 Valid values:`COS`: A COS 存储桶 地址 `URL`: A URL `AWS-S3`: An AWS S3 存储桶 地址 Currently，this 类型 is only supported for transcoding tasks。",
															},
															"cos_input_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The information of the COS object to process. This parameter is valid and 必填 when `类型` is `COS`。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bucket": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The COS 存储桶 of the object to process，such as `TopRankVideo-125xxx88`。",
																		},
																		"region": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The 地域 of the COS 存储桶，such as `ap-chongqing`。",
																		},
																		"object": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The 路径 of the object to process，such as `/movie/201907/WildAnimal.mov`。",
																		},
																	},
																},
															},
															"url_input_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The URL of the object to process. This parameter is valid and 必填 when `类型` is `URL`.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "URL of a video。",
																		},
																	},
																},
															},
															"s3_input_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "The information of the AWS S3 object processed. This parameter 为必填项 if `类型` is `AWS-S3`.注意：此字段可能返回 null，表示无法获取有效值。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"s3_bucket": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The AWS S3 存储桶",
																		},
																		"s3_region": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The 地域 of the AWS S3 存储桶",
																		},
																		"s3_object": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The 路径 of the AWS S3 object。",
																		},
																		"s3_secret_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The 键 ID 必填 to access the AWS S3 object。",
																		},
																		"s3_secret_key": {
																			Type:        schema.TypeString,
																			Optional:    true,
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
					},
				},
			},

			"ai_content_review_task": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "类型 parameter of a video 内容 audit task。",
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
				Description: "Video 内容 analysis task parameter。",
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
							Description: "An extended parameter，whose 值 is a stringfied JSON.Note: This parameter is for customers with special requirements. It needs to be customized offline.注意：此字段可能返回 null，表示无法获取有效值。",
						},
					},
				},
			},

			"ai_recognition_task": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "类型 parameter of a video 内容 recognition task。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Intelligent video recognition 模板 ID",
						},
					},
				},
			},

			"ai_quality_control_task": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The parameters of a quality control task。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The ID quality control template.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"channel_ext_para": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The channel extension parameter，which is a serialized JSON string.注意：此字段可能返回 null，表示无法获取有效值。",
						},
					},
				},
			},

			"task_notify_config": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Event notification information of a task. 如果此参数为空，no event notifications will be obtained。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cmq_model": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The CMQ or TDMQ-CMQ model. 有效值：Queue，Topic。",
						},
						"cmq_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The CMQ or TDMQ-CMQ 地域，such as `sh` (Shanghai) or `bj` (Beijing)。",
						},
						"topic_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The CMQ or TDMQ-CMQ topic to receive notifications. This parameter is valid when `CmqModel` is `Topic`。",
						},
						"queue_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The CMQ or TDMQ-CMQ queue to receive notifications. This parameter is valid when `CmqModel` is `Queue`。",
						},
						"notify_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Workflow notification method. 有效值：Finish，Change. 如果此参数为空，`Finish` will be used。",
						},
						"notify_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The notification 类型 Valid values:`CMQ`: This 值 is no longer used. Please use `TDMQ-CMQ` instead.`TDMQ-CMQ`: 消息 queue`URL`: If `NotifyType` is set to `URL`，HTTP callbacks are sent to the URL specified by `NotifyUrl`. HTTP and JSON are 用于the callbacks. The packet 包含response parameters of the `ParseNotification` API.`SCF`: This notification 类型 is not recommended. You need to configure it in the SCF console.`AWS-SQS`: AWS queue. This 类型 is only supported for AWS tasks，and the queue must be in the same 地域 as the AWS 存储桶&lt;font color=red&gt;Note: If you do not pass this parameter or pass in an empty string，`CMQ` will be used. To use a different notification 类型，指定this parameter accordingly.&lt;/font&gt;。",
						},
						"notify_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "HTTP callback URL，必填 if `NotifyType` is set to `URL`。",
						},
						"aws_sqa": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The AWS SQS queue. This parameter 为必填项 if `NotifyType` is `AWS-SQS`.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sqa_region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The 地域 of the SQS queue。",
									},
									"sqa_queue_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The 名称 SQS queue。",
									},
									"s3_secret_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The 键 ID 必填 to read from/write to the SQS queue。",
									},
									"s3_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The 键 必填 to read from/write to the SQS queue。",
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
				Description: "Task flow 优先级 The higher the 值，the higher the 优先级 取值范围：[-10，10]. 如果此参数为空，0 will be used。",
			},

			"session_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID 用于deduplication. If there was a request with the same ID in the last three days，the current request will return an 错误 The ID can contain up to 50 characters. 如果此参数为空 or an empty string is entered，no deduplication will be performed。",
			},

			"session_context": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The 来源 context which is 用于pass through the 用户 request information. The task flow 状态 change callback will return the 值 of this field. It can contain up to 1,000 characters。",
			},

			"task_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The task 类型 `Online` (default): A task that is executed immediately. `Offline`: A task that is executed when the system is idle (within three days by default)。",
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
