package mps

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMpsWithdrawsWatermarkOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsWithdrawsWatermarkOperationCreate,
		Read:   resourceTencentCloudMpsWithdrawsWatermarkOperationRead,
		Delete: resourceTencentCloudMpsWithdrawsWatermarkOperationDelete,
		Schema: map[string]*schema.Schema{
			"input_info": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Input information of file for metadata getting。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The input 类型 有效值：`COS`: A COS 存储桶 地址  `URL`: A URL  `AWS-S3`: An AWS S3 存储桶 地址 Currently，this 类型 is only supported for transcoding tasks.。",
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
							Description: "The notification 类型 有效值： `CMQ`: This 值 is no longer used. Please use `TDMQ-CMQ` instead.  `TDMQ-CMQ`: 消息 queue  `URL`: If `NotifyType` is set to `URL`，HTTP callbacks are sent to the URL specified by `NotifyUrl`. HTTP and JSON are 用于the callbacks. The packet 包含response parameters of the `ParseNotification` API.  `SCF`: This notification 类型 is not recommended. You need to configure it in the SCF console.  `AWS-SQS`: AWS queue. This 类型 is only supported for AWS tasks，and the queue must be in the same 地域 as the AWS 存储桶 Note: If you do not pass this parameter or pass in an empty string，`CMQ` will be used. To use a different notification 类型，指定this parameter accordingly。",
						},
						"notify_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "HTTP callback URL，必填 if `NotifyType` is set to `URL`。",
						},
						"aws_sqs": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The AWS SQS queue. This parameter 为必填项 if `NotifyType` is `AWS-SQS`.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sqs_region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The 地域 of the SQS queue。",
									},
									"sqs_queue_name": {
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

			"session_context": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The 来源 context which is 用于pass through the 用户 request information. The task flow 状态 change callback will return the 值 of this field。",
			},
		},
	}
}

func resourceTencentCloudMpsWithdrawsWatermarkOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_withdraws_watermark_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = mps.NewWithdrawsWatermarkRequest()
		response = mps.NewWithdrawsWatermarkResponse()
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
		if awsSqsMap, ok := helper.InterfaceToMap(dMap, "aws_sqs"); ok {
			awsSqs := mps.AwsSQS{}
			if v, ok := awsSqsMap["sqs_region"]; ok {
				awsSqs.SQSRegion = helper.String(v.(string))
			}
			if v, ok := awsSqsMap["sqs_queue_name"]; ok {
				awsSqs.SQSQueueName = helper.String(v.(string))
			}
			if v, ok := awsSqsMap["s3_secret_id"]; ok {
				awsSqs.S3SecretId = helper.String(v.(string))
			}
			if v, ok := awsSqsMap["s3_secret_key"]; ok {
				awsSqs.S3SecretKey = helper.String(v.(string))
			}
			taskNotifyConfig.AwsSQS = &awsSqs
		}
		request.TaskNotifyConfig = &taskNotifyConfig
	}

	if v, ok := d.GetOk("session_context"); ok {
		request.SessionContext = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().WithdrawsWatermark(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mps withdrawsWatermarkOperation failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(taskId)

	return resourceTencentCloudMpsWithdrawsWatermarkOperationRead(d, meta)
}

func resourceTencentCloudMpsWithdrawsWatermarkOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_withdraws_watermark_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMpsWithdrawsWatermarkOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_withdraws_watermark_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
