package ssl

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSslReplaceCertificateOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslReplaceCertificateCreate,
		Read:   resourceTencentCloudSslReplaceCertificateRead,
		Delete: resourceTencentCloudSslReplaceCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "证书 ID",
			},

			"valid_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Verification 类型: DNS_AUTO = automatic DNS verification (此 verification 类型 是 仅 支持 对于 域名 names 该 是 resolved 通过 Tencent Cloud 和 have normal resolution 状态)，DNS = manual DNS verification，FILE = 文件 verification。",
			},

			"csr_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "类型，默认值 Original. Available options: Original = original 证书 CSR，Upload = manual upload，Online = online generation。",
			},

			"csr_content": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CSR 内容",
			},

			"csr_key_password": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "KEY 密码",
			},

			"reason": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Reason 对于 reissue。",
			},

			"cert_csr_encrypt_algo": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CSR 加密 方法，可选: RSA，ECC，SM2. (Selectable 仅 如果 CsrType 是 Online)，默认为 RSA。",
			},

			"cert_csr_key_parameter": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "CSR 加密 参数，当 CsrEncryptAlgo 是 RSA，您 可以 choose 2048，4096，etc.，和 默认为 2048; 当 CsrEncryptAlgo 是 ECC，您 可以 choose prime256v1，secp384r1，etc.，和 默认为 prime256v1;。",
			},
		},
	}
}

func resourceTencentCloudSslReplaceCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_replace_certificate_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request       = ssl.NewReplaceCertificateRequest()
		response      = ssl.NewReplaceCertificateResponse()
		certificateId uint64
	)
	if v, ok := d.GetOk("certificate_id"); ok {
		request.CertificateId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("valid_type"); ok {
		request.ValidType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("csr_type"); ok {
		request.CsrType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("csr_content"); ok {
		request.CsrContent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("csr_key_password"); ok {
		request.CsrkeyPassword = helper.String(v.(string))
	}

	if v, ok := d.GetOk("reason"); ok {
		request.Reason = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cert_csr_encrypt_algo"); ok {
		request.CertCSREncryptAlgo = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cert_csr_key_parameter"); ok {
		request.CertCSRKeyParameter = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSSLCertificateClient().ReplaceCertificate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ssl replaceCertificate failed, reason:%+v", logId, err)
		return err
	}
	if response != nil && response.Response != nil && response.Response.CertificateId != nil {
		certificateId = helper.StrToUInt64(*response.Response.CertificateId)
	}

	d.SetId(helper.UInt64ToStr(certificateId))

	return resourceTencentCloudSslReplaceCertificateRead(d, meta)
}

func resourceTencentCloudSslReplaceCertificateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_replace_certificate_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSslReplaceCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_replace_certificate_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
