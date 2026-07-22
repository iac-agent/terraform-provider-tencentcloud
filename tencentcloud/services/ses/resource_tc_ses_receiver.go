package ses

import (
	"context"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSesReceiver() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSesReceiverCreate,
		Read:   resourceTencentCloudSesReceiverRead,
		Delete: resourceTencentCloudSesReceiverDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"receivers_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Recipient 组名称",
			},

			"desc": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Recipient 组 描述",
			},

			"data": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeSet,
				Description: "Recipient email 和 template 参数 在 数组 格式 数量 recipients 是 limited 到 within 20,000. 如果 there 是 对象 在 `数据` 列表 该 inputs `template_data`，then other objects 是 also 必填",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Recipient email addresses。",
						},
						"template_data": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Variable 参数 在 template，please 使用 json.dump 到 格式 JSON 对象 作为 字符串 类型 对象 是 集合 的 键-值 pairs，其中 each 键 表示 variable 在 template，和 variables 在 template 是 represented 通过 {{键}}，和 corresponding 值 将 是 replaced 使用 {{值}} 当 sent.注意: Parameter 值 不能 是 complex 数据 such 作为 HTML. 总数 长度 的 TemplateData ( entire JSON structure) should 是 less 比 800 bytes。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudSesReceiverCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_receiver.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = ses.NewCreateReceiverRequest()
		response   = ses.NewCreateReceiverResponse()
		receiverId uint64
	)
	if v, ok := d.GetOk("receivers_name"); ok {
		request.ReceiversName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("desc"); ok {
		request.Desc = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSesClient().CreateReceiver(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ses Receiver failed, reason:%+v", logId, err)
		return err
	}

	receiverId = *response.Response.ReceiverId

	if v, ok := d.GetOk("data"); ok {
		datas := v.(*schema.Set).List()
		dataList := make([]*ses.ReceiverInputData, 0, len(datas))
		emilList := make([]*string, 0, len(datas))
		for _, item := range datas {
			var email string
			var templateData string
			dMap := item.(map[string]interface{})
			if v, ok := dMap["email"]; ok {
				email = v.(string)
			}
			if v, ok := dMap["template_data"]; ok {
				templateData = v.(string)
			}

			if templateData != "" {
				receiver := ses.ReceiverInputData{}
				receiver.Email = &email
				receiver.TemplateData = &templateData
				dataList = append(dataList, &receiver)
			} else {
				emilList = append(emilList, &email)
			}
		}

		if len(dataList) > 0 {
			request := ses.NewCreateReceiverDetailWithDataRequest()
			request.ReceiverId = &receiverId
			request.Datas = dataList
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSesClient().CreateReceiverDetailWithData(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s create ses receiverDetail failed, reason:%+v", logId, err)
				return err
			}
		}
		if len(emilList) > 0 {
			request := ses.NewCreateReceiverDetailRequest()
			request.ReceiverId = &receiverId
			request.Emails = emilList
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSesClient().CreateReceiverDetail(request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s create ses receiverDetail failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	d.SetId(strconv.Itoa(int(receiverId)))

	return resourceTencentCloudSesReceiverRead(d, meta)
}

func resourceTencentCloudSesReceiverRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_receiver.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := SesService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	receiverId := d.Id()
	receiver, err := service.DescribeSesReceiverById(ctx, receiverId)
	if err != nil {
		return err
	}

	if receiver == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SesReceiver` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if receiver.ReceiversName != nil {
		_ = d.Set("receivers_name", receiver.ReceiversName)
	}

	if receiver.Desc != nil {
		_ = d.Set("desc", receiver.Desc)
	}

	receiverData, err := service.DescribeSesReceiverDetailById(ctx, receiverId)
	if err != nil {
		return err
	}
	if receiverData != nil {
		dataList := []interface{}{}
		for _, data := range receiverData {
			dataMap := map[string]interface{}{}

			if data.Email != nil {
				dataMap["email"] = data.Email
			}

			if data.TemplateData != nil {
				dataMap["template_data"] = data.TemplateData
			}

			dataList = append(dataList, dataMap)
		}

		_ = d.Set("data", dataList)
	}

	return nil
}

func resourceTencentCloudSesReceiverDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ses_receiver.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := SesService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	receiverId := d.Id()

	if err := service.DeleteSesReceiverById(ctx, receiverId); err != nil {
		return err
	}

	return nil
}
