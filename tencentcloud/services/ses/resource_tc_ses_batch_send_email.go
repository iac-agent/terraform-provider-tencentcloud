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
				Description: "Sender 地址 Enter sender 地址 such 作为 noreply@mail.qcloud.com. To display sender 名称，enter 地址 在 following 格式:sender &amp;amp;lt;email 地址&amp;amp;gt;. For 示例:Tencent Cloud team &amp;amp;lt;noreply@mail.qcloud.com&amp;amp;gt;。",
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
				Description: "任务 类型 1: immediate; 2: scheduled; 3: recurring。",
			},

			"reply_to_addresses": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Reply-到 地址 You 可以 enter 有效 personal email 地址 该 可以 receive emails. 如果此参数为空，reply emails 将 fail 到 是 sent。",
			},

			"template": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "模板 当 emails 是 sent 使用 template。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "模板 ID 如果 您 do 不 have any template，please create 一个。",
						},
						"template_data": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Variable 参数 在 template. Please 使用 json.dump 到 格式 JSON 对象 into 字符串 类型The 对象 是 集合 的 键-值 pairs. Each 键 denotes variable，其中 是 represented 通过 {{键}}. 键 将 是 replaced 使用 correspondingvalue (represented 通过 {{值}}) 当 sending email.注意: 参数 值 不能 是 数据 的 complex 类型 such 作为 HTML.Example: {名称:xxx,age:xx}。",
						},
					},
				},
			},

			"attachments": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Attachment 参数 到 集合 当 您 need 到 send attachments. 此 参数 是 currently unavailable。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Attachment 名称，其中 不能 exceed 255 字符. Some attachment types 是 不 支持. For details，see [Attachment Types.](https://www.tencentcloud.com/document/product/1084/42373?has_map=1)。",
						},
						"content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Base64-encoded attachment 内容 You 可以 send attachments 的 up 到 4 MB 在 总数 大小.注意: TencentCloud API 支持 请求 packet 的 up 到 8 MB 在 大小，和 大小 的 attachmentcontent 将 increase 通过 1.5 times after Base64 编码. Therefore，您 need 到 keep 总数 大小 的 allattachments below 4 MB. 如果 entire 请求 exceeds 8 MB， API 将 返回 错误",
						},
					},
				},
			},

			"cycle_param": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Parameter 必填 对于 recurring sending 任务。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"begin_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "开始时间 的 任务。",
						},
						"interval_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "任务 recurrence 在 hours。",
						},
						"term_cycle": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "指定是否end cycle. 此 参数 是 用于update 任务. 有效值：0: No; 1: Yes。",
						},
					},
				},
			},

			"timed_param": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Parameter 必填 对于 scheduled sending 任务。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"begin_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "开始时间 的 scheduled sending 任务。",
						},
					},
				},
			},

			"unsubscribe": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Unsubscribe link 选项. 0: Do 不 add unsubscribe link; 1: English 2: Simplified Chinese; 3: Traditional Chinese; 4: Spanish; 5: French; 6: German; 7: Japanese; 8: Korean; 9: Arabic; 10: Thai。",
			},

			"ad_location": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "是否add ad 标签 0: Add 无 标签; 1: Add before subject; 2: Add after subject。",
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
