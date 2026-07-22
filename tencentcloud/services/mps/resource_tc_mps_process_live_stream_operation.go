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
				Description: "Live 流 URL，其中 必须 是 live 流 文件 地址 RTMP，HLS，和 FLV 是 支持。",
			},

			"task_notify_config": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Event 通知 信息 的 任务，其中 是 用于指定live 流 processing 结果",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cmq_model": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CMQ model. There 是 two types: `Queue` 和 `Topic`. Currently，仅 `Queue` 是 支持。",
						},
						"cmq_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CMQ 地域，such 作为 `sh` 和 `bj`。",
						},
						"queue_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "此 参数 是 有效 当 model 是 `Queue`，indicating 名称 CMQ queue 对于 receiving 事件 notifications。",
						},
						"topic_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "此 参数 是 有效 当 model 是 `Topic`，indicating 名称 CMQ 主题 对于 receiving 事件 notifications。",
						},
						"notify_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "通知 类型，`CMQ` 通过 默认值. 如果 此 参数 是 集合 到 `URL`，HTTP callbacks 是 sent 到 URL 指定 通过 `NotifyUrl`.注意: 如果 您 do 不 pass 此 参数 或 pass 在 空 字符串，`CMQ` 将 是 使用. To 使用 different 通知 类型，指定this 参数 accordingly。",
						},
						"notify_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "HTTP callback URL，必填 如果 `NotifyType` 是 集合 到 `URL`。",
						},
					},
				},
			},

			"output_storage": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Target 存储桶 的 live 流 processing 输出文件 此 参数 为必填项 如果 文件 将 是 output。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "存储 类型 对于 media processing 输出文件 有效 值:`COS`: Tencent Cloud COS.`AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
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
				Description: "Target directory 的 live 流 processing 输出文件，such 作为 `/movie/201909/`. 如果此参数为空， `/` directory 将 是 使用。",
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

			"ai_analysis_task": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "AI 视频 intelligent analysis input 参数 types。",
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

			"ai_quality_control_task": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "参数 对于 视频 quality control 任务。",
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

			"session_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID 用于deduplication. 如果 there 是 请求 使用 same ID 在 last seven days， 当前 请求 将 返回 错误 ID 可以 contain up 到 50 字符. 如果此参数为空 或 空 字符串 是 entered，无 deduplication 将 是 performed。",
			},

			"session_context": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "来源 context 其中 是 用于pass through 用户 请求 信息. 任务 flow 状态 change callback 将 返回 值 的 此 字段. It 可以 contain up 到 1,000 字符。",
			},

			"schedule_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "scheme ID.注意 1: About `OutputStorage` 和 `OutputDir`:如果 output 存储 和 directory 是 指定 对于 subtask 的 scheme，those output settings 将 是 applied.如果 output 存储 和 directory 是 不 指定 对于 subtasks 的 scheme， output 参数 passed 在 `ProcessMedia` API 将 是 applied.注意 2: 如果 `TaskNotifyConfig` 是 指定， 指定 settings 将 是 使用 instead 的 默认值 callback settings 的 scheme。",
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
