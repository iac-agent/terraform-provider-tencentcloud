package gaap

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudGaapCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapCertificatesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID 证书 到 是 queried。",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 证书 到 是 queried。",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"BASIC", "CLIENT", "SERVER", "REALSERVER", "PROXY"}),
				Description:  "类型 证书 到 是 queried. 有效值：`BASIC`，`CLIENT`，`SERVER`，`REALSERVER` 和 `PROXY`. `BASIC` 表示 basic 证书; `CLIENT` 表示 客户端 CA 证书; `SERVER` 表示 服务器 SSL 证书; `REALSERVER` 表示 realserver CA 证书; `PROXY` 表示 proxy SSL 证书。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// computed
			"certificates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An 信息 列表 证书. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 证书。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 证书。",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 证书。",
						},
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
				},
			},
		},
	}
}

func dataSourceTencentCloudGaapCertificatesRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_gaap_certificates.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		id              *string
		name            *string
		certificateType *int
		ids             []string
		certificates    []map[string]interface{}
	)

	if raw, ok := d.GetOk("id"); ok {
		id = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("name"); ok {
		name = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("type"); ok {
		certificateType = common.IntPtr(gaapCertificateStringMap[raw.(string)])
	}

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	respCertificates, err := service.DescribeCertificates(ctx, id, name, certificateType)
	if err != nil {
		return err
	}

	ids = make([]string, 0, len(respCertificates))
	certificates = make([]map[string]interface{}, 0, len(respCertificates))
	for _, certificate := range respCertificates {
		ids = append(ids, *certificate.CertificateId)

		var (
			certificateType string
			ok              bool
		)
		if certificateType, ok = gaapCertificateIntMap[int(*certificate.CertificateType)]; !ok {
			return fmt.Errorf("unknown certificate type %d", *certificate.CertificateType)
		}

		m := map[string]interface{}{
			"id":          *certificate.CertificateId,
			"name":        *certificate.CertificateAlias,
			"type":        certificateType,
			"create_time": helper.FormatUnixTime(*certificate.CreateTime),
		}

		if certificate.BeginTime != nil {
			m["begin_time"] = helper.FormatUnixTime(*certificate.BeginTime)
		}

		if certificate.EndTime != nil {
			m["end_time"] = helper.FormatUnixTime(*certificate.EndTime)
		}

		if certificate.IssuerCN != nil {
			m["issuer_cn"] = *certificate.IssuerCN
		}

		if certificate.SubjectCN != nil {
			m["subject_cn"] = *certificate.SubjectCN
		}

		certificates = append(certificates, m)
	}

	_ = d.Set("certificates", certificates)
	d.SetId(helper.DataResourceIdsHash(ids))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), certificates); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}

	return nil
}
