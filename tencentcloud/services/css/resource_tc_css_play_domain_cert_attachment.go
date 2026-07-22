package css

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCssPlayDomainCertAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssPlayDomainCertAttachmentCreate,
		Read:   resourceTencentCloudCssPlayDomainCertAttachmentRead,
		Delete: resourceTencentCloudCssPlayDomainCertAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_info": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "playback domains 到 bind 和 是否enable HTTPS 对于 them. 如果 `CloudCertId` 是 unspecified，和 域名 是 already bound 使用 证书，此 API 将 仅 update HTTPS 配置 的 域名",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "域名 名称",
						},
						"status": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "是否enable https 规则 对于 域名 名称 1: 启用，0: 已禁用，-1: remain unchanged。",
						},
					},
				},
			},
			"cloud_cert_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Tencent 云 ssl 证书 ID. Refer 到 `tencentcloud_ssl_certificate` 到 create 或 obtain 资源 ID。",
			},
			"certificate_alias": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "证书 备注 Synonymous 使用 CertName。",
			},
			"cert_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "证书 类型 0: Self-owned 证书，1: Tencent Cloud ssl managed 证书。",
			},
			"cert_expire_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "证书 过期时间。",
			},
			"cert_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "证书 ID",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "时间 当 规则 是 last 更新。",
			},
		},
	}
}

func resourceTencentCloudCssPlayDomainCertAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_play_domain_cert_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request     = css.NewModifyLiveDomainCertBindingsRequest()
		response    = css.NewModifyLiveDomainCertBindingsResponse()
		cloudCertId string
		domainName  string
	)

	if v, ok := d.GetOk("cloud_cert_id"); ok {
		cloudCertId = v.(string)
		request.CloudCertId = helper.String(cloudCertId)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "domain_info"); ok {
		info := css.LiveCertDomainInfo{}
		if v, ok := dMap["domain_name"]; ok {
			domainName = v.(string)
			info.DomainName = helper.String(domainName)
		}
		if v, ok := dMap["status"]; ok {
			info.Status = helper.IntInt64(v.(int))
		}
		request.DomainInfos = append(request.DomainInfos, &info)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCssClient().ModifyLiveDomainCertBindings(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create css playDomainCertAttachment failed, error reason: %+v", logId, err)
		return err
	}

	if len(response.Response.Errors) > 0 {
		return fmt.Errorf("[CRITAL]%s create css playDomainCertAttachment failed, reason: response.Response.Errors[%+v]", logId, response.Response.Errors)
	}

	d.SetId(strings.Join([]string{domainName, cloudCertId}, tccommon.FILED_SP))

	return resourceTencentCloudCssPlayDomainCertAttachmentRead(d, meta)
}

func resourceTencentCloudCssPlayDomainCertAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_play_domain_cert_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	domainName := idSplit[0]
	cloudCertId := idSplit[1]

	playDomainCertAttachment, err := service.DescribeCssPlayDomainCertAttachmentById(ctx, domainName, cloudCertId)
	if err != nil {
		return err
	}

	if playDomainCertAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssPlayDomainCertAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if playDomainCertAttachment.CloudCertId != nil {
		_ = d.Set("cloud_cert_id", playDomainCertAttachment.CloudCertId)
	}

	domainInfosMap := map[string]interface{}{}
	if playDomainCertAttachment.DomainName != nil {
		domainInfosMap["domain_name"] = playDomainCertAttachment.DomainName
	}

	if playDomainCertAttachment.Status != nil {
		domainInfosMap["status"] = playDomainCertAttachment.Status
	}
	_ = d.Set("domain_info", []interface{}{domainInfosMap})

	if playDomainCertAttachment.CertificateAlias != nil {
		_ = d.Set("certificate_alias", playDomainCertAttachment.CertificateAlias)
	}

	if playDomainCertAttachment.CertType != nil {
		_ = d.Set("cert_type", playDomainCertAttachment.CertType)
	}

	if playDomainCertAttachment.CertExpireTime != nil {
		_ = d.Set("cert_expire_time", playDomainCertAttachment.CertExpireTime)
	}

	if playDomainCertAttachment.CertId != nil {
		_ = d.Set("cert_id", playDomainCertAttachment.CertId)
	}

	if playDomainCertAttachment.UpdateTime != nil {
		_ = d.Set("update_time", playDomainCertAttachment.UpdateTime)
	}

	return nil
}

func resourceTencentCloudCssPlayDomainCertAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_play_domain_cert_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	domainName := idSplit[0]

	if err := service.DeleteCssPlayDomainCertAttachmentById(ctx, domainName); err != nil {
		return err
	}

	return nil
}
