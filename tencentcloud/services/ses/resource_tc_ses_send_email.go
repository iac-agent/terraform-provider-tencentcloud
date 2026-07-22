package ses

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSesSendEmail() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSesSendEmailCreate,
		Read:   resourceTencentCloudSesSendEmailRead,
		Delete: resourceTencentCloudSesSendEmailDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"from_email_address": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Sender 地址 Enter sender 地址，对于 示例，noreply@mail.qcloud.com.To display sender 名称，enter 地址 在 following 格式:Sender。",
			},

			"destination": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Recipient email addresses. You 可以 send email 到 up 到 50 recipients 在 时间. 注意: email 内容 将 display all recipient addresses. To send 一个-到-一个 emails 到 several recipients，please call API 多个 times 到 send emails。",
			},

			"subject": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Email subject。",
			},

			"reply_to_addresses": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Reply-到 地址 You 可以 enter 有效 personal email 地址 该 可以 receive emails. 如果此参数为空，reply emails 将 fail 到 是 sent。",
			},

			"cc": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Cc recipient email 地址，up 到 20 people 可以 是 copied。",
			},

			"bcc": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "email 地址 的 cc recipient 可以 support up 到 20 cc recipients。",
			},

			"template": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "模板 参数 对于 template-based sending. As Simple has been disused，模板 为必填项。",
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
				Description: "Parameters 对于 attachments 到 是 sent. TencentCloud API 支持 请求 packet 的 up 到 8 MB 在 大小,和 大小 的 attachment 内容 将 increase 通过 1.5 times after Base64 编码. Therefore,您 need 到 keep 总数 大小 的 all attachments below 4 MB. 如果 entire 请求 exceeds 8 MB, API 将 返回 错误",
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

			"unsubscribe": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Unsubscribe link 选项. 0: Do 不 add unsubscribe link; 1: English 2: Simplified Chinese; 3: Traditional Chinese; 4: Spanish; 5: French; 6: German; 7: Japanese; 8: Korean; 9: Arabic; 10: Thai。",
			},

			"trigger_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Email triggering 类型 0 (默认值): non-触发器-based，suitable 对于 marketing emails 和 non-immediate emails;1: 触发器-based，suitable 对于 immediate emails such 作为 emails containing verification codes.如果 大小 的 email exceeds 指定 值, 系统 将 automatically choose non-触发器-based 类型",
			},
		},
	}
}

func resourceTencentCloudSesSendEmailCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_send_email.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request   = ses.NewSendEmailRequest()
		response  = ses.NewSendEmailResponse()
		messageId string
	)
	if v, ok := d.GetOk("from_email_address"); ok {
		request.FromEmailAddress = helper.String(v.(string))
	}

	if v, ok := d.GetOk("destination"); ok {
		destinationSet := v.(*schema.Set).List()
		for i := range destinationSet {
			destination := destinationSet[i].(string)
			request.Destination = append(request.Destination, &destination)
		}
	}

	if v, ok := d.GetOk("subject"); ok {
		request.Subject = helper.String(v.(string))
	}

	if v, ok := d.GetOk("reply_to_addresses"); ok {
		request.ReplyToAddresses = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cc"); ok {
		ccSet := v.(*schema.Set).List()
		for i := range ccSet {
			cc := ccSet[i].(string)
			request.Cc = append(request.Cc, &cc)
		}
	}

	if v, ok := d.GetOk("bcc"); ok {
		bccSet := v.(*schema.Set).List()
		for i := range bccSet {
			bcc := bccSet[i].(string)
			request.Bcc = append(request.Bcc, &bcc)
		}
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

	if v, ok := d.GetOk("unsubscribe"); ok {
		request.Unsubscribe = helper.String(v.(string))
	}

	if v, _ := d.GetOk("trigger_type"); v != nil {
		request.TriggerType = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSesClient().SendEmail(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ses sendEmail failed, reason:%+v", logId, err)
		return err
	}

	messageId = *response.Response.MessageId
	d.SetId(messageId)

	return resourceTencentCloudSesSendEmailRead(d, meta)
}

func resourceTencentCloudSesSendEmailRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_send_email.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSesSendEmailDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_send_email.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
