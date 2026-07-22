package mps

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMpsProcessLiveStreamOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsProcessLiveStreamOperationCreate,
		Read:   resourceTencentCloudMpsProcessLiveStreamOperationRead,
		Delete: resourceTencentCloudMpsProcessLiveStreamOperationDelete,
		Schema: map[string]*schema.Schema{
			"url": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Live stream URL，which must be a live stream file 地址 RTMP，HLS，and FLV are supported。",
			},

			"task_notify_config": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Event notification information of a task，which is 用于指定live stream processing 结果",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cmq_model": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CMQ model. There are two types: `Queue` and `Topic`. Currently，only `Queue` is supported。",
						},
						"cmq_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CMQ 地域，such as `sh` and `bj`。",
						},
						"queue_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This parameter is valid when the model is `Queue`，indicating the 名称 CMQ queue for receiving event notifications。",
						},
						"topic_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This parameter is valid when the model is `Topic`，indicating the 名称 CMQ topic for receiving event notifications。",
						},
						"notify_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The notification 类型，`CMQ` by default. If this parameter is set to `URL`，HTTP callbacks are sent to the URL specified by `NotifyUrl`.Note: If you do not pass this parameter or pass in an empty string，`CMQ` will be used. To use a different notification 类型，指定this parameter accordingly。",
						},
						"notify_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "HTTP callback URL，必填 if `NotifyType` is set to `URL`。",
						},
					},
				},
			},

			"output_storage": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Target 存储桶 of a live stream processing 输出文件 This parameter 为必填项 if a file will be output。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The storage 类型 for a media processing 输出文件 Valid values:`COS`: Tencent Cloud COS.`AWS-S3`: AWS S3. This 类型 is only supported for AWS tasks，and the 输出存储桶 must be in the same 地域 as the 存储桶 of the 来源 file。",
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
				Description: "Target directory of a live stream processing 输出文件，such as `/movie/201909/`. 如果此参数为空，the `/` directory will be used。",
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

			"ai_recognition_task": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "类型 parameter of video 内容 recognition task。",
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

			"ai_analysis_task": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "AI video intelligent analysis input parameter types。",
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

			"ai_quality_control_task": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The parameters for a video quality control task。",
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

			"session_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID 用于deduplication. If there was a request with the same ID in the last seven days，the current request will return an 错误 The ID can contain up to 50 characters. 如果此参数为空 or an empty string is entered，no deduplication will be performed。",
			},

			"session_context": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The 来源 context which is 用于pass through the 用户 request information. The task flow 状态 change callback will return the 值 of this field. It can contain up to 1,000 characters。",
			},

			"schedule_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "The scheme ID.Note 1: About `OutputStorage` and `OutputDir`:If an output storage and directory are specified for a subtask of the scheme，those output settings will be applied.If an output storage and directory are not specified for the subtasks of a scheme，the output parameters passed in the `ProcessMedia` API will be applied.Note 2: If `TaskNotifyConfig` is specified，the specified settings will be used instead of the default callback settings of the scheme。",
			},
		},
	}
}

func resourceTencentCloudMpsProcessLiveStreamOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_process_live_stream_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = mps.NewProcessLiveStreamRequest()
		response = mps.NewProcessLiveStreamResponse()
		taskId   string
	)
	if v, ok := d.GetOk("url"); ok {
		request.Url = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "task_notify_config"); ok {
		liveStreamTaskNotifyConfig := mps.LiveStreamTaskNotifyConfig{}
		if v, ok := dMap["cmq_model"]; ok {
			liveStreamTaskNotifyConfig.CmqModel = helper.String(v.(string))
		}
		if v, ok := dMap["cmq_region"]; ok {
			liveStreamTaskNotifyConfig.CmqRegion = helper.String(v.(string))
		}
		if v, ok := dMap["queue_name"]; ok {
			liveStreamTaskNotifyConfig.QueueName = helper.String(v.(string))
		}
		if v, ok := dMap["topic_name"]; ok {
			liveStreamTaskNotifyConfig.TopicName = helper.String(v.(string))
		}
		if v, ok := dMap["notify_type"]; ok {
			liveStreamTaskNotifyConfig.NotifyType = helper.String(v.(string))
		}
		if v, ok := dMap["notify_url"]; ok {
			liveStreamTaskNotifyConfig.NotifyUrl = helper.String(v.(string))
		}
		request.TaskNotifyConfig = &liveStreamTaskNotifyConfig
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

	if dMap, ok := helper.InterfacesHeadMap(d, "ai_content_review_task"); ok {
		aiContentReviewTaskInput := mps.AiContentReviewTaskInput{}
		if v, ok := dMap["definition"]; ok {
			aiContentReviewTaskInput.Definition = helper.IntUint64(v.(int))
		}
		request.AiContentReviewTask = &aiContentReviewTaskInput
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "ai_recognition_task"); ok {
		aiRecognitionTaskInput := mps.AiRecognitionTaskInput{}
		if v, ok := dMap["definition"]; ok {
			aiRecognitionTaskInput.Definition = helper.IntUint64(v.(int))
		}
		request.AiRecognitionTask = &aiRecognitionTaskInput
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

	if v, ok := d.GetOk("session_id"); ok {
		request.SessionId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("session_context"); ok {
		request.SessionContext = helper.String(v.(string))
	}

	if v, _ := d.GetOk("schedule_id"); v != nil {
		request.ScheduleId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().ProcessLiveStream(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mps processLiveStreamOperation failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(taskId)

	return resourceTencentCloudMpsProcessLiveStreamOperationRead(d, meta)
}

func resourceTencentCloudMpsProcessLiveStreamOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_process_live_stream_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMpsProcessLiveStreamOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_process_live_stream_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
