package sms

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSmsSign() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudSmsSignRead,
		Create: resourceTencentCloudSmsSignCreate,
		Update: resourceTencentCloudSmsSignUpdate,
		Delete: resourceTencentCloudSmsSignDelete,
		Schema: map[string]*schema.Schema{
			"sign_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Sms sign 名称，唯一。",
			},

			"sign_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Sms sign 类型: 0，1，2，3，4，5，6。",
			},

			"document_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "DocumentType 是 用于enterprise authentication，或 website，app authentication，etc. DocumentType: 0，1，2，3，4，5，6，7，8。",
			},

			"international": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "是否为Global SMS: 0: Mainland China SMS; 1: Global SMS。",
			},

			"sign_purpose": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "签名 purpose: 0: 对于 personal 使用; 1: 对于 others。",
			},

			"proof_image": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "You should Base64-encode 镜像 的 identity 证书 corresponding 到 签名 first，remove prefix 数据:镜像/jpeg;base64，从 resulted 字符串，和 then 使用 它 作为 值 的 此 参数。",
			},

			"commission_image": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Power 的 attorney，其中 should 是 submitted 如果 SignPurpose 是 对于 使用 通过 others. You should Base64-encode 镜像 first，remove prefix 数据:镜像/jpeg;base64，从 resulted 字符串，和 then 使用 它 作为 值 的 此 参数. 注意: 此 字段 将 take effect 仅 当 SignPurpose 是 1 (对于 用户 通过 others)。",
			},

			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "签名 应用 备注",
			},
		},
	}
}

func resourceTencentCloudSmsSignCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sms_sign.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request       = sms.NewAddSmsSignRequest()
		response      *sms.AddSmsSignResponse
		signId        uint64
		international int
	)

	if v, ok := d.GetOk("sign_name"); ok {
		request.SignName = helper.String(v.(string))
	}

	if v, _ := d.GetOk("sign_type"); v != nil {
		request.SignType = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("document_type"); v != nil {
		request.DocumentType = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("international"); v != nil {
		international = v.(int)
		request.International = helper.IntUint64(international)
	}

	if v, _ := d.GetOk("sign_purpose"); v != nil {
		request.SignPurpose = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("proof_image"); ok {
		request.ProofImage = helper.String(v.(string))
	}

	if v, ok := d.GetOk("commission_image"); ok {
		request.CommissionImage = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSmsClient().AddSmsSign(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sms sign failed, reason:%+v", logId, err)
		return err
	}

	signId = *response.Response.AddSignStatus.SignId
	d.SetId(helper.UInt64ToStr(signId) + tccommon.FILED_SP + strconv.Itoa(international))
	return resourceTencentCloudSmsSignRead(d, meta)
}

func resourceTencentCloudSmsSignRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sms_sign.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := SmsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	signId := idSplit[0]
	international := idSplit[1]

	sign, err := service.DescribeSmsSign(ctx, signId, international)

	if err != nil {
		return err
	}

	if sign == nil {
		d.SetId("")
		return fmt.Errorf("resource `sign` %s does not exist", signId)
	}

	if sign.SignName != nil {
		_ = d.Set("sign_name", sign.SignName)
	}

	if sign.International != nil {
		_ = d.Set("international", sign.International)
	}

	return nil
}

func resourceTencentCloudSmsSignUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sms_sign.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := sms.NewModifySmsSignRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	signId := idSplit[0]

	request.SignId = helper.Uint64(helper.StrToUInt64(signId))

	if v, ok := d.GetOk("sign_name"); ok {
		request.SignName = helper.String(v.(string))
	}

	if v, _ := d.GetOk("sign_type"); v != nil {
		request.SignType = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("document_type"); v != nil {
		request.DocumentType = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("international"); v != nil {
		request.International = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("sign_purpose"); v != nil {
		request.SignPurpose = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("proof_image"); ok {
		request.ProofImage = helper.String(v.(string))
	}

	if v, ok := d.GetOk("commission_image"); ok {
		request.CommissionImage = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSmsClient().ModifySmsSign(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create sms sign failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSmsSignRead(d, meta)
}

func resourceTencentCloudSmsSignDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sms_sign.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := SmsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	signId := idSplit[0]

	if err := service.DeleteSmsSignById(ctx, signId); err != nil {
		return err
	}

	return nil
}
