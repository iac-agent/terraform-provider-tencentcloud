package mps

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMpsEditMediaOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsEditMediaOperationCreate,
		Read:   resourceTencentCloudMpsEditMediaOperationRead,
		Delete: resourceTencentCloudMpsEditMediaOperationDelete,
		Schema: map[string]*schema.Schema{
			"file_infos": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Information 的 input 视频 文件。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"input_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Video input 信息。",
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
						"start_time_offset": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: "开始时间 偏移量 的 视频 clipping （秒）。",
						},
						"end_time_offset": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: "结束时间 偏移量 的 视频 clipping （秒）。",
						},
					},
				},
			},

			"output_storage": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "存储 location 的 media processing 输出文件",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "存储 类型 对于 media processing 输出文件 有效值：`COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. 此 类型 是 仅 支持 对于 AWS tasks，和 输出存储桶 必须 是 在 same 地域 作为 存储桶 的 来源 文件。",
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
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "路径 到 save media processing 输出文件",
			},

			"output_config": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Configuration 对于 output files 的 视频 editing。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "格式 有效值：`mp4` (默认值)，`hls`，`mov`，`flv`，`avi`。",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "editing 模式 有效 值 是 `normal` 和 `fast`. 默认为 `normal`，其中 表示precise editing。",
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
							Description: "通知 类型 有效值：`CMQ`: 此 值 是 无 longer 使用. Please 使用 `TDMQ-CMQ` instead. `TDMQ-CMQ`: 消息 queue. `URL`: 如果 `NotifyType` 是 集合 到 `URL`，HTTP callbacks 是 sent 到 URL 指定 通过 `NotifyUrl`. HTTP 和 JSON 是 用于the callbacks. packet 包含response 参数 的 `ParseNotification` API. `SCF`: 此 通知 类型 是 不 recommended. You need 到 configure 它 在 SCF console. `AWS-SQS`: AWS queue. 此 类型 是 仅 支持 对于 AWS tasks，和 queue 必须 是 在 same 地域 作为 AWS 存储桶 如果 您 do 不 pass 此 参数 或 pass 在 空 字符串，`CMQ` 将 是 使用. To 使用 different 通知 类型，指定this 参数 accordingly。",
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

			"tasks_priority": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "任务 优先级 higher 值， higher 优先级 取值范围：[-10,10]. 如果此参数为空，0 将 是 使用。",
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
		},
	}
}

func resourceTencentCloudMpsEditMediaOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_edit_media_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = mps.NewEditMediaRequest()
		response = mps.NewEditMediaResponse()
		taskId   string
	)
	if v, ok := d.GetOk("file_infos"); ok {
		for _, item := range v.([]interface{}) {
			editMediaFileInfo := mps.EditMediaFileInfo{}
			dMap := item.(map[string]interface{})
			if inputInfoMap, ok := helper.InterfaceToMap(dMap, "input_info"); ok {
				mediaInputInfo := mps.MediaInputInfo{}
				if v, ok := inputInfoMap["type"]; ok {
					mediaInputInfo.Type = helper.String(v.(string))
				}
				if cosInputInfoMap, ok := helper.InterfaceToMap(inputInfoMap, "cos_input_info"); ok {
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
				if urlInputInfoMap, ok := helper.InterfaceToMap(inputInfoMap, "url_input_info"); ok {
					urlInputInfo := mps.UrlInputInfo{}
					if v, ok := urlInputInfoMap["url"]; ok {
						urlInputInfo.Url = helper.String(v.(string))
					}
					mediaInputInfo.UrlInputInfo = &urlInputInfo
				}
				if s3InputInfoMap, ok := helper.InterfaceToMap(inputInfoMap, "s3_input_info"); ok {
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
				editMediaFileInfo.InputInfo = &mediaInputInfo
			}
			if v, ok := dMap["start_time_offset"]; ok {
				editMediaFileInfo.StartTimeOffset = helper.Float64(v.(float64))
			}
			if v, ok := dMap["end_time_offset"]; ok {
				editMediaFileInfo.EndTimeOffset = helper.Float64(v.(float64))
			}
			request.FileInfos = append(request.FileInfos, &editMediaFileInfo)
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

	if v, ok := d.GetOk("output_object_path"); ok {
		request.OutputObjectPath = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "output_config"); ok {
		editMediaOutputConfig := mps.EditMediaOutputConfig{}
		if v, ok := dMap["container"]; ok {
			editMediaOutputConfig.Container = helper.String(v.(string))
		}
		if v, ok := dMap["type"]; ok {
			editMediaOutputConfig.Type = helper.String(v.(string))
		}
		request.OutputConfig = &editMediaOutputConfig
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

	if v, _ := d.GetOk("tasks_priority"); v != nil {
		request.TasksPriority = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("session_id"); ok {
		request.SessionId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("session_context"); ok {
		request.SessionContext = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().EditMedia(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mps editMediaOperation failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(taskId)

	return resourceTencentCloudMpsEditMediaOperationRead(d, meta)
}

func resourceTencentCloudMpsEditMediaOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_edit_media_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMpsEditMediaOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_edit_media_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
