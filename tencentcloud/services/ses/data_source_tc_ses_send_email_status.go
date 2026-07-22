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
				Description: "Date sent. 此 参数 为必填项. You 可以 仅 查询 sending 状态 对于 单个 date 在 时间。",
			},

			"message_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "MessageId 字段 返回 通过 SendMail API。",
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
							Description: "MessageId 字段 返回 通过 SendEmail API。",
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
							Description: "Tencent Cloud processing 状态: `0`: Successful. `1001`: Internal 系统 exception. `1002`: Internal 系统 exception. `1003`: Internal 系统 exception. `1003`: Internal 系统 exception. `1004`: Email sending timed out. `1005`: Internal 系统 exception. `1006`: You have sent too many emails 到 same 地址 在 short 周期 `1007`: email 地址 是 在 blocklist. `1008`: sender 域名 是 rejected 通过 recipient. `1009`: Internal 系统 exception. `1010`: daily email sending 限制 是 exceeded. `1011`: You have 无 权限 到 send 自定义 内容 Use template. `1013`: sender 域名 是 unsubscribed 从 通过 recipient. `2001`: No results 是 found. `3007`: 模板 ID 是 无效 或 template 是 unavailable. `3008`: sender 域名 是 temporarily blocked 通过 recipient 域名 `3009`: You have 无 权限 到 使用 此 template. `3010`: 格式 的 TemplateData 字段 是 incorrect. `3014`: email 不能 是 sent because sender 域名 是 不 verified. `3020`: recipient email 地址 是 在 blocklist. `3024`: Failed 到 precheck email 地址 格式 `3030`: Email sending 是 restricted temporarily due 到 high bounce 速率. `3033`: 账号 has insufficient balance 或 overdue payment。",
						},
						"deliver_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Recipient processing status0: Tencent Cloud has accepted 请求 和 added 它 到 send queue.1: email 是 delivered successfully. DeliverTime 表示time 当 email 是 delivered successfully.2: email 是 discarded. DeliverMessage 表示reason 对于 discarding.3: recipient&amp;#39;s ESP rejects email，probably because email 地址 does 不 exist 或 due 到 other reasons.8: email 是 delayed 通过 ESP. DeliverMessage 表示reason 对于 延迟",
						},
						"deliver_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 recipient processing 状态",
						},
						"request_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "时间戳 当 请求 arrives 在 Tencent Cloud。",
						},
						"deliver_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "时间戳 当 Tencent Cloud delivers email。",
						},
						"user_opened": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否recipient has opened email。",
						},
						"user_clicked": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否recipient has clicked links 在 email。",
						},
						"user_unsubscribed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否recipient has unsubscribed 从 email sent 通过 sender。",
						},
						"user_complainted": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否recipient has reported sender。",
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
