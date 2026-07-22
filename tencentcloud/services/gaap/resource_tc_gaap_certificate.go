package gaap

import (
	"context"
	"errors"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudGaapCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapCertificateCreate,
		Read:   resourceTencentCloudGaapCertificateRead,
		Update: resourceTencentCloudGaapCertificateUpdate,
		Delete: resourceTencentCloudGaapCertificateDelete,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"BASIC", "CLIENT", "SERVER", "REALSERVER", "PROXY"}),
				Description:  "类型 证书. 有效 值: `BASIC`，`CLIENT`，`SERVER`，`REALSERVER` 和 `PROXY`. `BASIC` 表示 basic 证书; `CLIENT` 表示 客户端 CA 证书; `SERVER` 表示 服务器 SSL 证书; `REALSERVER` 表示 realserver CA 证书; `PROXY` 表示 proxy SSL 证书。",
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "内容 的 证书，和 URL 编码. 当 证书 是 basic authentication，使用 `用户:xxx 密码:xxx` 格式，其中 密码 是 encrypted 使用 `htpasswd` 或 `openssl`; 当 证书 是 `CA` 或 `SSL`， 格式 是 `pem`。",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "名称 证书。",
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "键 的 `SSL` 证书。",
			},

			// computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "创建时间 的 证书。",
			},
			"begin_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Beginning 时间 的 证书。",
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Ending 时间 的 证书。",
			},
			"issuer_cn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Issuer 名称 证书。",
			},
			"subject_cn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subject 名称 证书。",
			},
		},
	}
}

func resourceTencentCloudGaapCertificateCreate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_certificate.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	certificateType := gaapCertificateStringMap[d.Get("type").(string)]

	content := d.Get("content").(string)

	name := d.Get("name").(string)

	var key *string
	if rawKey, ok := d.GetOk("key"); ok {
		key = helper.String(rawKey.(string))
	}

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	id, err := service.createCertificate(ctx, certificateType, content, name, key)
	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceTencentCloudGaapCertificateRead(d, m)
}

func resourceTencentCloudGaapCertificateRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_certificate.read")()
	defer tccommon.InconsistentCheck(d, m)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	certificate, err := service.DescribeCertificateById(ctx, id)
	if err != nil {
		return err
	}

	if certificate == nil {
		d.SetId("")
		return nil
	}

	if certificate.CertificateType == nil {
		return errors.New("certificate type is nil")
	}
	if certType, ok := gaapCertificateIntMap[int(*certificate.CertificateType)]; ok {
		_ = d.Set("type", certType)
	} else {
		return fmt.Errorf("unknown certificate type %d", *certificate.CertificateType)
	}

	_ = d.Set("name", certificate.CertificateAlias)

	if certificate.CreateTime == nil {
		return errors.New("certificate create time is nil")
	}
	_ = d.Set("create_time", helper.FormatUnixTime(*certificate.CreateTime))

	if certificate.BeginTime != nil {
		_ = d.Set("begin_time", helper.FormatUnixTime(*certificate.BeginTime))
	}
	if certificate.EndTime != nil {
		_ = d.Set("end_time", helper.FormatUnixTime(*certificate.EndTime))
	}

	_ = d.Set("issuer_cn", certificate.IssuerCN)
	_ = d.Set("subject_cn", certificate.SubjectCN)

	return nil
}

func resourceTencentCloudGaapCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_certificate.update")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	name := d.Get("name").(string)

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	if err := service.ModifyCertificateName(ctx, id, name); err != nil {
		return err
	}

	return resourceTencentCloudGaapCertificateRead(d, m)
}

func resourceTencentCloudGaapCertificateDelete(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_certificate.delete")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	return service.DeleteCertificate(ctx, id)
}
