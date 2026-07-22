package ses

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSesBatchSendEmail() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSesBatchSendEmailCreate,
		Read:   resourceTencentCloudSesBatchSendEmailRead,
		Delete: resourceTencentCloudSesBatchSendEmailDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"from_email_address": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Sender 地址 Enter a sender 地址 such as noreply@mail.qcloud.com. To display the sender 名称，enter the 地址 in the following 格式:sender &amp;amp;lt;email 地址&amp;amp;gt;. For example:Tencent Cloud team &amp;amp;lt;noreply@mail.qcloud.com&amp;amp;gt;。",
			},

			"receiver_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Recipient 组 ID",
			},

			"subject": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Email subject。",
			},

			"task_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Task 类型 1: immediate; 2: scheduled; 3: recurring。",
			},

			"reply_to_addresses": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Reply-to 地址 You can enter a valid personal email 地址 that can receive emails. 如果此参数为空，reply emails will fail to be sent。",
			},

			"template": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Template when emails are sent using a template。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "模板 ID If you do not have any template，please create one。",
						},
						"template_data": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Variable parameters in the template. Please use json.dump to 格式 the JSON object into a string 类型The object is a set of 键-值 pairs. Each 键 denotes a variable，which is represented by {{键}}. The 键 will be replaced with the correspondingvalue (represented by {{值}}) when sending the email.Note: The parameter 值 cannot be data of a complex 类型 such as HTML.Example: {名称:xxx,age:xx}。",
						},
					},
				},
			},

			"attachments": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Attachment parameters to set when you need to send attachments. This parameter is currently unavailable。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Attachment 名称，which cannot exceed 255 characters. Some attachment types are not supported. For details，see [Attachment Types.](https://www.tencentcloud.com/document/product/1084/42373?has_map=1)。",
						},
						"content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Base64-encoded attachment 内容 You can send attachments of up to 4 MB in the total size.Note: The TencentCloud API supports a request packet of up to 8 MB in size，and the size of the attachmentcontent will increase by 1.5 times after Base64 encoding. Therefore，you need to keep the total size of allattachments below 4 MB. If the entire request exceeds 8 MB，the API will return an 错误",
						},
					},
				},
			},

			"cycle_param": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Parameter 必填 for a recurring sending task。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"begin_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "开始时间 of the task。",
						},
						"interval_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Task recurrence in hours。",
						},
						"term_cycle": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "指定是否end the cycle. This parameter is 用于update the task. 有效值：0: No; 1: Yes。",
						},
					},
				},
			},

			"timed_param": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Parameter 必填 for a scheduled sending task。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"begin_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "开始时间 of a scheduled sending task。",
						},
					},
				},
			},

			"unsubscribe": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Unsubscribe link option.  0: Do not add unsubscribe link; 1: English 2: Simplified Chinese;  3: Traditional Chinese; 4: Spanish; 5: French;  6: German; 7: Japanese; 8: Korean;  9: Arabic; 10: Thai。",
			},

			"ad_location": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "是否add an ad 标签 0: Add no 标签; 1: Add before the subject; 2: Add after the subject。",
			},
		},
	}
}

func resourceTencentCloudSesBatchSendEmailCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_batch_send_email.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = ses.NewBatchSendEmailRequest()
		response = ses.NewBatchSendEmailResponse()
		taskId   uint64
	)
	if v, ok := d.GetOk("from_email_address"); ok {
		request.FromEmailAddress = helper.String(v.(string))
	}

	if v, _ := d.GetOk("receiver_id"); v != nil {
		request.ReceiverId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("subject"); ok {
		request.Subject = helper.String(v.(string))
	}

	if v, _ := d.GetOk("task_type"); v != nil {
		request.TaskType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("reply_to_addresses"); ok {
		request.ReplyToAddresses = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "template"); ok {
		template := ses.Template{}
		if v, ok := dMap["template_id"]; ok {
			template.TemplateID = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["template_data"]; ok {
			template.TemplateData = helper.String(v.(string))
		}
		request.Template = &template
	}

	if v, ok := d.GetOk("attachments"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			attachment := ses.Attachment{}
			if v, ok := dMap["file_name"]; ok {
				attachment.FileName = helper.String(v.(string))
			}
			if v, ok := dMap["content"]; ok {
				attachment.Content = helper.String(v.(string))
			}
			request.Attachments = append(request.Attachments, &attachment)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "cycle_param"); ok {
		cycleEmailParam := ses.CycleEmailParam{}
		if v, ok := dMap["begin_time"]; ok {
			cycleEmailParam.BeginTime = helper.String(v.(string))
		}
		if v, ok := dMap["interval_time"]; ok {
			cycleEmailParam.IntervalTime = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["term_cycle"]; ok {
			cycleEmailParam.TermCycle = helper.IntUint64(v.(int))
		}
		request.CycleParam = &cycleEmailParam
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "timed_param"); ok {
		timedEmailParam := ses.TimedEmailParam{}
		if v, ok := dMap["begin_time"]; ok {
			timedEmailParam.BeginTime = helper.String(v.(string))
		}
		request.TimedParam = &timedEmailParam
	}

	if v, ok := d.GetOk("unsubscribe"); ok {
		request.Unsubscribe = helper.String(v.(string))
	}

	if v, _ := d.GetOk("ad_location"); v != nil {
		request.ADLocation = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSesClient().BatchSendEmail(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ses batchSendEmail failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(helper.UInt64ToStr(taskId))

	return resourceTencentCloudSesBatchSendEmailRead(d, meta)
}

func resourceTencentCloudSesBatchSendEmailRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_batch_send_email.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSesBatchSendEmailDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_batch_send_email.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
