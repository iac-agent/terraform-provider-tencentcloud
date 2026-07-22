package tse

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTseGatewayCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTseGatewayCertificatesRead,
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "网关 ID",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 conditions，有效 值: `BindDomain`，`名称`。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "过滤名称",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "过滤值",
						},
					},
				},
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "结果",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "总数 注意：此字段可能返回 null，表示有效值不可用。",
						},
						"certificates_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Certificate 列表 网关. 注意：此字段可能返回 null，表示有效值不可用。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Certificate 名称 注意：此字段可能返回 null，表示有效值不可用。",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "证书 ID 注意：此字段可能返回 null，表示有效值不可用。",
									},
									"bind_domains": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Domains 的 binding. 注意：此字段可能返回 null，表示有效值不可用。",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "状态 证书. Reference 值:- expired- 活跃 注意：此字段可能返回 null，表示有效值不可用。",
									},
									"crt": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Pem 格式 的 证书. 注意：此字段可能返回 null，表示有效值不可用。",
									},
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Private 键 的 证书. 注意：此字段可能返回 null，表示有效值不可用。",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "过期时间 的 证书. 注意：此字段可能返回 null，表示有效值不可用。",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Upload 时间 的 证书. 注意：此字段可能返回 null，表示有效值不可用。",
									},
									"issue_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Issuance 时间 的 certificate注意：此字段可能返回 null，表示有效值不可用。",
									},
									"cert_source": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "来源 的 证书. Reference 值:- native. 来源: konga- ssl. 来源: ssl 平台. 注意：此字段可能返回 null，表示有效值不可用。",
									},
									"cert_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "证书 ID ssl 平台. 注意：此字段可能返回 null，表示有效值不可用。",
									},
								},
							},
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

func dataSourceTencentCloudTseGatewayCertificatesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tse_gateway_certificates.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("gateway_id"); ok {
		paramMap["GatewayId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*tse.ListFilter, 0, len(filtersSet))

		for _, item := range filtersSet {
			listFilter := tse.ListFilter{}
			listFilterMap := item.(map[string]interface{})

			if v, ok := listFilterMap["key"]; ok {
				listFilter.Key = helper.String(v.(string))
			}
			if v, ok := listFilterMap["value"]; ok {
				listFilter.Value = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &listFilter)
		}
		paramMap["filters"] = tmpSet
	}

	service := TseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var result *tse.KongCertificatesList
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeTseGatewayCertificatesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	var ids []string
	kongCertificatesListMap := map[string]interface{}{}
	if result != nil {
		ids = make([]string, 0, *result.Total)
		if result.Total != nil {
			kongCertificatesListMap["total"] = result.Total
		}

		if result.CertificatesList != nil {
			certificatesListList := []interface{}{}
			for _, certificatesList := range result.CertificatesList {
				certificatesListMap := map[string]interface{}{}

				if certificatesList.Name != nil {
					certificatesListMap["name"] = certificatesList.Name
				}

				if certificatesList.Id != nil {
					certificatesListMap["id"] = certificatesList.Id
				}

				if certificatesList.BindDomains != nil {
					certificatesListMap["bind_domains"] = certificatesList.BindDomains
				}

				if certificatesList.Status != nil {
					certificatesListMap["status"] = certificatesList.Status
				}

				if certificatesList.Crt != nil {
					certificatesListMap["crt"] = certificatesList.Crt
				}

				if certificatesList.Key != nil {
					certificatesListMap["key"] = certificatesList.Key
				}

				if certificatesList.ExpireTime != nil {
					certificatesListMap["expire_time"] = certificatesList.ExpireTime
				}

				if certificatesList.CreateTime != nil {
					certificatesListMap["create_time"] = certificatesList.CreateTime
				}

				if certificatesList.IssueTime != nil {
					certificatesListMap["issue_time"] = certificatesList.IssueTime
				}

				if certificatesList.CertSource != nil {
					certificatesListMap["cert_source"] = certificatesList.CertSource
				}

				if certificatesList.CertId != nil {
					certificatesListMap["cert_id"] = certificatesList.CertId
				}

				certificatesListList = append(certificatesListList, certificatesListMap)
				ids = append(ids, *certificatesList.Id)
			}

			kongCertificatesListMap["certificates_list"] = certificatesListList
		}

		_ = d.Set("result", []interface{}{kongCertificatesListMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), kongCertificatesListMap); e != nil {
			return e
		}
	}
	return nil
}
