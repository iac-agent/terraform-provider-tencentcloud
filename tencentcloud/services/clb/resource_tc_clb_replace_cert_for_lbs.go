package clb

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClbReplaceCertForLbs() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbReplaceCertForLbsCreate,
		Read:   resourceTencentCloudClbReplaceCertForLbsRead,
		Delete: resourceTencentCloudClbReplaceCertForLbsDelete,
		Schema: map[string]*schema.Schema{
			"old_certificate_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "待替换的证书ID，可以是服务器证书，也可以是客户端证书。",
			},

			"certificate": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "新证书的内容等信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ssl_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "认证类型。取值范围：UNIDIRECTIONAL（单向认证）、MUTUAL（相互认证）。",
						},
						"cert_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "服务器证书ID。如果将此参数留空，则必须上传证书，包括 CertContent、CertKey 和 CertName。",
						},
						"cert_ca_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "客户端证书ID。当监听器采用双向认证（即SSLMode=mutual）时，如果该参数留空，则必须上传客户端证书，包括CertCaContent和CertCaName。",
						},
						"cert_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "上传的服务器证书名称。如果没有CertId，则该参数为必填项。",
						},
						"cert_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "上传的服务器证书的密钥。如果没有CertId，则该参数为必填项。",
						},
						"cert_content": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "上传的服务器证书内容。如果没有CertId，则该参数为必填项。",
						},
						"cert_ca_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "上传的客户端CA证书的名称。当SSLMode=mutual时，如果没有CertCaId，则该参数为必填项。",
						},
						"cert_ca_content": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "上传的客户端证书的内容。当SSLMode=mutual时，如果没有CertCaId，则该参数为必填项。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClbReplaceCertForLbsCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_replace_cert_for_lbs.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request          = clb.NewReplaceCertForLoadBalancersRequest()
		oldCertificateId string
	)
	if v, ok := d.GetOk("old_certificate_id"); ok {
		oldCertificateId = v.(string)
		request.OldCertificateId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "certificate"); ok {
		certificateInput := clb.CertificateInput{}
		if v, ok := dMap["ssl_mode"]; ok {
			certificateInput.SSLMode = helper.String(v.(string))
		}
		if v, ok := dMap["cert_id"]; ok {
			certificateInput.CertId = helper.String(v.(string))
		}
		if v, ok := dMap["cert_ca_id"]; ok {
			certificateInput.CertCaId = helper.String(v.(string))
		}
		if v, ok := dMap["cert_name"]; ok {
			certificateInput.CertName = helper.String(v.(string))
		}
		if v, ok := dMap["cert_key"]; ok {
			certificateInput.CertKey = helper.String(v.(string))
		}
		if v, ok := dMap["cert_content"]; ok {
			certificateInput.CertContent = helper.String(v.(string))
		}
		if v, ok := dMap["cert_ca_name"]; ok {
			certificateInput.CertCaName = helper.String(v.(string))
		}
		if v, ok := dMap["cert_ca_content"]; ok {
			certificateInput.CertCaContent = helper.String(v.(string))
		}
		request.Certificate = &certificateInput
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ReplaceCertForLoadBalancers(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate clb replaceCertForLbs failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(oldCertificateId)

	return resourceTencentCloudClbReplaceCertForLbsRead(d, meta)
}

func resourceTencentCloudClbReplaceCertForLbsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_replace_cert_for_lbs.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudClbReplaceCertForLbsDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_replace_cert_for_lbs.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
