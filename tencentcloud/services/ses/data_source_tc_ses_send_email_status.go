package ses

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSesSendEmailStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSesSendEmailStatusRead,
		Schema: map[string]*schema.Schema{
			"request_date": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Date sent. This parameter 为必填项. You can only query the sending 状态 for a single date at a time。",
			},

			"message_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The MessageId field returned by the SendMail API。",
			},

			"to_email_address": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Recipient email 地址",
			},

			"email_status_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "状态 sent emails。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"message_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The MessageId field returned by the SendEmail API。",
						},
						"to_email_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Recipient email 地址",
						},
						"from_email_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sender email 地址",
						},
						"send_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Tencent Cloud processing 状态: `0`: Successful. `1001`: Internal system exception. `1002`: Internal system exception. `1003`: Internal system exception. `1003`: Internal system exception. `1004`: Email sending timed out. `1005`: Internal system exception. `1006`: You have sent too many emails to the same 地址 in a short 周期 `1007`: The email 地址 is in the blocklist. `1008`: The sender 域名 is rejected by the recipient. `1009`: Internal system exception. `1010`: The daily email sending 限制 is exceeded. `1011`: You have no permission to send custom 内容 Use a template. `1013`: The sender 域名 is unsubscribed from by the recipient. `2001`: No results were found. `3007`: The 模板 ID is invalid or the template is unavailable. `3008`: The sender 域名 is temporarily blocked by the recipient 域名 `3009`: You have no permission to use this template. `3010`: The 格式 of the TemplateData field is incorrect. `3014`: The email cannot be sent because the sender 域名 is not verified. `3020`: The recipient email 地址 is in the blocklist. `3024`: Failed to precheck the email 地址 格式 `3030`: Email sending is restricted temporarily due to a high bounce rate. `3033`: The 账号 has insufficient balance or overdue payment。",
						},
						"deliver_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Recipient processing status0: Tencent Cloud has accepted the request and added it to the send queue.1: The email is delivered successfully. DeliverTime 表示time when the email is delivered successfully.2: The email is discarded. DeliverMessage 表示reason for discarding.3: The recipient&amp;#39;s ESP rejects the email，probably because the email 地址 does not exist or due to other reasons.8: The email is delayed by the ESP. DeliverMessage 表示reason for 延迟",
						},
						"deliver_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 recipient processing 状态",
						},
						"request_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "时间戳 when the request arrives at Tencent Cloud。",
						},
						"deliver_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "时间戳 when Tencent Cloud delivers the email。",
						},
						"user_opened": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否recipient has opened the email。",
						},
						"user_clicked": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否recipient has clicked the links in the email。",
						},
						"user_unsubscribed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否recipient has unsubscribed from the email sent by the sender。",
						},
						"user_complainted": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否recipient has reported the sender。",
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

func dataSourceTencentCloudSesSendEmailStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ses_send_email_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("request_date"); ok {
		paramMap["RequestDate"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("message_id"); ok {
		paramMap["MessageId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("to_email_address"); ok {
		paramMap["ToEmailAddress"] = helper.String(v.(string))
	}

	service := SesService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var emailStatusList []*ses.SendEmailStatus

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSesSendEmailStatusByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		emailStatusList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(emailStatusList))
	tmpList := make([]map[string]interface{}, 0, len(emailStatusList))

	if emailStatusList != nil {
		for _, sendEmailStatus := range emailStatusList {
			sendEmailStatusMap := map[string]interface{}{}

			if sendEmailStatus.MessageId != nil {
				sendEmailStatusMap["message_id"] = sendEmailStatus.MessageId
			}

			if sendEmailStatus.ToEmailAddress != nil {
				sendEmailStatusMap["to_email_address"] = sendEmailStatus.ToEmailAddress
			}

			if sendEmailStatus.FromEmailAddress != nil {
				sendEmailStatusMap["from_email_address"] = sendEmailStatus.FromEmailAddress
			}

			if sendEmailStatus.SendStatus != nil {
				sendEmailStatusMap["send_status"] = sendEmailStatus.SendStatus
			}

			if sendEmailStatus.DeliverStatus != nil {
				sendEmailStatusMap["deliver_status"] = sendEmailStatus.DeliverStatus
			}

			if sendEmailStatus.DeliverMessage != nil {
				sendEmailStatusMap["deliver_message"] = sendEmailStatus.DeliverMessage
			}

			if sendEmailStatus.RequestTime != nil {
				sendEmailStatusMap["request_time"] = sendEmailStatus.RequestTime
			}

			if sendEmailStatus.DeliverTime != nil {
				sendEmailStatusMap["deliver_time"] = sendEmailStatus.DeliverTime
			}

			if sendEmailStatus.UserOpened != nil {
				sendEmailStatusMap["user_opened"] = sendEmailStatus.UserOpened
			}

			if sendEmailStatus.UserClicked != nil {
				sendEmailStatusMap["user_clicked"] = sendEmailStatus.UserClicked
			}

			if sendEmailStatus.UserUnsubscribed != nil {
				sendEmailStatusMap["user_unsubscribed"] = sendEmailStatus.UserUnsubscribed
			}

			if sendEmailStatus.UserComplainted != nil {
				sendEmailStatusMap["user_complainted"] = sendEmailStatus.UserComplainted
			}

			ids = append(ids, *sendEmailStatus.MessageId)
			tmpList = append(tmpList, sendEmailStatusMap)
		}

		_ = d.Set("email_status_list", tmpList)
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
