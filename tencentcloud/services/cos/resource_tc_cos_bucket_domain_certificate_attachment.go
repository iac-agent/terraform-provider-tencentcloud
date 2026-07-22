package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	cos "github.com/tencentyun/cos-go-sdk-v5"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCosBucketDomainCertificateAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCosBucketDomainCertificateAttachmentCreate,
		Read:   resourceTencentCloudCosBucketDomainCertificateAttachmentRead,
		// Update: resourceTencentCloudCosBucketDomainCertificateAttachmentUpdate,
		Delete: resourceTencentCloudCosBucketDomainCertificateAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateCosBucketName,
				Description:  "存储桶名称",
			},
			"domain_certificate": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				ForceNew:    true,
				Description: "证书 的 指定 doamin。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							ForceNew:    true,
							Description: "Certificate info。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cert_type": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Certificate 类型",
									},
									"custom_cert": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Required:    true,
										ForceNew:    true,
										Description: "Custom 证书。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cert_id": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "ID 证书。",
												},
												"cert": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Sensitive:   true,
													Description: "Public 键 的 证书。",
												},
												"private_key": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Sensitive:   true,
													Description: "Private 键 的 证书。",
												},
											},
										},
									},
								},
							},
						},
						"domain": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "名称 域名",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCosBucketDomainCertificateAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_domain_certificate_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId  = tccommon.GetLogId(tccommon.ContextNil)
		ctx    = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		bucket string
	)

	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	} else {
		return errors.New("get bucket failed!")
	}

	option := cos.BucketPutDomainCertificateOptions{}
	if dcMap, ok := helper.InterfacesHeadMap(d, "domain_certificate"); ok {
		if certMap, ok := helper.InterfaceToMap(dcMap, "certificate"); ok {
			certificateInfo := cos.BucketDomainCertificateInfo{}
			if v, ok := certMap["cert_type"]; ok {
				certificateInfo.CertType = v.(string)
			}

			if CustomCertMap, ok := helper.InterfaceToMap(certMap, "custom_cert"); ok {
				customCert := cos.BucketDomainCustomCert{}
				if v, ok := CustomCertMap["cert_id"]; ok {
					customCert.CertId = v.(string)
				}

				if v, ok := CustomCertMap["cert"]; ok {
					customCert.Cert = v.(string)
				}

				if v, ok := CustomCertMap["private_key"]; ok {
					customCert.PrivateKey = v.(string)
				}

				certificateInfo.CustomCert = &customCert
			}

			option.CertificateInfo = &certificateInfo
		}

		if v, ok := dcMap["domain"]; ok {
			option.DomainList = append(option.DomainList, v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTencentCosClient(bucket).Bucket.PutDomainCertificate(ctx, &option)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Create cos domain certificate failed, Response is nil."))
			}

			request, _ := xml.Marshal(option)
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, "PutDomainCertificate", request, result.Response.Body)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cos bucketDomainCertificate failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{bucket, option.DomainList[0]}, tccommon.FILED_SP))
	return nil
}

func resourceTencentCloudCosBucketDomainCertificateAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_domain_certificate_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = CosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id      = d.Id()
	)

	certResult, bucket, err := service.DescribeCosBucketDomainCertificate(ctx, id)
	log.Printf("[DEBUG] resource `bucketDomainCertificate certResult: %s`\n", certResult)
	if err != nil {
		return err
	}

	if certResult == nil {
		d.SetId("")
		return fmt.Errorf("resource `bucketDomainCertificate` %s does not exist", id)
	}

	_ = d.Set("bucket", bucket)

	return nil
}

func resourceTencentCloudCosBucketDomainCertificateAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_domain_certificate_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = CosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id      = d.Id()
	)

	if err := service.DeleteCosBucketDomainCertificate(ctx, id); err != nil {
		return err
	}

	return nil
}
