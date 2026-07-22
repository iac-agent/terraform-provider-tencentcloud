package ssl

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSslDescribeCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeCertificateRead,
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "证书 ID",
			},
			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "结果 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"owner_uin": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "账号 UIN.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"project_id": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "项目 IDNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"from": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Certificate 来源: Trustasia,uploadNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"certificate_type": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Certificate 类型: CA = CA 证书，SVR = 服务器 证书.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"package_type": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Types 的 Certificate Package: 1 = Geotrust DV SSL CA -G3，2 = Trustasia TLS RSA CA，3 = SecureSite Enhanced Enterprise Edition (EV Pro)，4 = SecureSite enhanced (EV)，5 = SecureSite Enterprise Professional Edition (OVPro)，6 = SecureSite Enterprise (OV)，7 = SecureSite Enterprise (OV) compatriots，8 = Geotrust enhanced 类型 (EV)，9 = Geotrust Enterprise (OV)，10 = Geotrust Enterprise (OV) pass,11 = Trustasia 域名 Multi -域名 SSL 证书，12 = Trustasia 域名 model (DV) passing，13 = Trustasia Enterprise Passing Character (OV) SSL 证书 (D3)，14 = Trustasia Enterprise (OV) SSL 证书 (D3)，15= Trustasia Enterprise Multi -域名 名称 (OV) SSL 证书 (D3)，16 = Trustasia enhanced (EV) SSL 证书 (D3)，17 = Trustasia enhanced multi -域名 名称 (EV) SSL 证书 (D3)，18 = GlobalSign enterprise 类型 enterprise 类型(OV) SSL 证书，19 = GlobalSign Enterprise 类型 -类型 STL Certificate，20 = GlobalSign enhanced (EV) SSL 证书，21 = Trustasia Enterprise Tongzhi Multi -域名 名称 (OV) SSL 证书 (D3)，22 = GlobalSignignMulti -域名 名称 (OV) SSL 证书，23 = GlobalSign Enterprise 类型 -类型 multi -域名 名称 (OV) SSL 证书，24 = GlobalSign enhanced multi -域名 名称 (EV) SSL 证书.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"product_zh_name": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Certificate issuer 名称Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"domain": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "域名 名称Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"alias": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "备注 名称Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"status": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "= Submitted 信息，到 是 uploaded 到 confirmation letter，9 = Certificate 是 revoked，10 = revoked，11 = Re -issuance，12 = Upload 和 revoke confirmation letter.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"status_msg": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "状态 信息.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"verify_type": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Verification 类型: DNS_AUTO = Automatic DNS verification，DNS = manual DNS verification，文件 = 文件 verification，email = email verification.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"vulnerability_status": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Vulnerability scanning 状态Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"cert_begin_time": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Certificate takes effect 时间.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"cert_end_time": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "证书 是 无效 时间.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"validity_period": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Validity 周期: 单位 (month).注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"insert_time": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "应用 时间.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"order_id": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "顺序 ID.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"certificate_extra": {
							Computed:    true,
							Type:        schema.TypeList,
							Description: "Certificate extension 信息.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain_number": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Certificate 可以 是 已配置 在 数量 域名 names.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"origin_certificate_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Original 证书 IDNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"replaced_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Re -issue original ID 证书.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"replaced_for": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Re -issue new ID.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"renew_order": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "New 顺序 证书 IDNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"s_m_cert": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Is 它 national secret certificateNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"company_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "类型 company. 注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
								},
							},
						},

						"dv_auth_detail": {
							Computed:    true,
							Type:        schema.TypeList,
							Description: "DV certification 信息.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dv_auth_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DV certification 键Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"dv_auth_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DV certification 值Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"dv_auth_domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DV authentication 值 域名 名称Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"dv_auth_path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DV authentication 值 路径Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"dv_auth_key_sub_domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DV certification sub -域名 名称Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"dv_auths": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "DV certification 信息.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"dv_auth_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "DV certification 键Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
												},
												"dv_auth_value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "DV certification 值Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
												},
												"dv_auth_domain": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "DV authentication 值 域名 名称Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
												},
												"dv_auth_path": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "DV authentication 值 路径Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
												},
												"dv_auth_sub_domain": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "DV certification sub -域名 名称,注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
												},
												"dv_auth_verify_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "DV certification 类型Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
												},
											},
										},
									},
								},
							},
						},

						"vulnerability_report": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Vulnerability scanning evaluation 报告.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"package_type_name": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "Certificate 类型 名称Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"status_name": {
							Computed:    true,
							Type:        schema.TypeString,
							Description: "状态 描述Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"subject_alt_name": {
							Computed: true,
							Type:     schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "证书 包含multiple 域名 names (containing main 域名 名称).注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"is_vip": {
							Computed:    true,
							Type:        schema.TypeBool,
							Description: "是否为a VIP customer.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"is_wildcard": {
							Computed:    true,
							Type:        schema.TypeBool,
							Description: "是否为a pan -域名 证书 证书.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"is_dv": {
							Computed:    true,
							Type:        schema.TypeBool,
							Description: "是否为the DV 版本Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"is_vulnerability": {
							Computed:    true,
							Type:        schema.TypeBool,
							Description: "是否vulnerability scanning 函数 是 已启用Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"renew_able": {
							Computed:    true,
							Type:        schema.TypeBool,
							Description: "Whether 您 可以 issue 证书.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"submitted_data": {
							Computed:    true,
							Type:        schema.TypeList,
							Description: "Submitted 信息 信息.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"csr_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CSR 类型，(online = online CSR，PARSE = paste CSR).注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"csr_content": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CSR 内容Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"certificate_domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "域名 信息.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"domain_list": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "DNS 信息.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"key_password": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Private 键 密码Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"organization_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Enterprise 或 单位 名称Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"organization_division": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "department.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"organization_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "地址Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"organization_country": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "nation.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"organization_city": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "city.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"organization_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Province.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"postal_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Postal 代码Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"phone_area_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Local 地域 代码Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"phone_number": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Landline 数量.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"admin_first_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Administrator 名称Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"admin_last_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "surname 的 administrator.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"admin_phone_num": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Administrator phone 数量.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"admin_email": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Administrator mailbox 地址Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"admin_position": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Administrator position.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"contact_first_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Contact 名称Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"contact_last_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Contact surname.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"contact_number": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Contact phone 数量.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"contact_email": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Contact mailbox 地址,注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"contact_position": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Contact position.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"verify_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Verification 类型Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
								},
							},
						},

						"deployable": {
							Computed:    true,
							Type:        schema.TypeBool,
							Description: "Whether 它 可以 是 deployed.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"c_a_encrypt_algorithms": {
							Computed: true,
							Type:     schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "All 加密 methods 的 CA certificateNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"c_a_common_names": {
							Computed: true,
							Type:     schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "All general names 的 CA certificateNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"c_a_end_times": {
							Computed: true,
							Type:     schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "CA 证书 all maturity timeNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},

						"dv_revoke_auth_detail": {
							Computed:    true,
							Type:        schema.TypeList,
							Description: "DV 证书 revoking verification valueNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dv_auth_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DV certification 键Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"dv_auth_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DV certification 值Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"dv_auth_domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DV authentication 值 域名 名称Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"dv_auth_path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DV authentication 值 路径Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"dv_auth_sub_domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DV certification sub -域名 名称,注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
									"dv_auth_verify_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DV certification 类型Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
									},
								},
							},
						},
					},
				}},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudSslDescribeCertificateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ssl_describe_certificate.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := SslService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	responese := ssl.DescribeCertificateResponseParams{}
	CertificateId := d.Get("certificate_id").(string)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSslDescribeCertificateByID(ctx, CertificateId)
		if e != nil {
			if sdkErr := helper.UnwarpSDKError(e); sdkErr != nil && tccommon.IsContains("LimitExceeded", sdkErr.Code) {
				return resource.RetryableError(e)
			}
			return tccommon.RetryError(e)
		}
		responese = *result
		return nil
	})
	if err != nil {
		return err
	}
	sslResponseMap := map[string]interface{}{}
	if responese.OwnerUin != nil {
		sslResponseMap["owner_uin"] = responese.OwnerUin
	}

	if responese.ProjectId != nil {
		sslResponseMap["project_id"] = responese.ProjectId
	}

	if responese.From != nil {
		sslResponseMap["from"] = responese.From
	}

	if responese.CertificateType != nil {
		sslResponseMap["certificate_type"] = responese.CertificateType
	}

	if responese.PackageType != nil {
		sslResponseMap["package_type"] = responese.PackageType
	}

	if responese.ProductZhName != nil {
		sslResponseMap["product_zh_name"] = responese.ProductZhName
	}

	if responese.Domain != nil {
		sslResponseMap["domain"] = responese.Domain
	}

	if responese.Alias != nil {
		sslResponseMap["alias"] = responese.Alias
	}

	if responese.Status != nil {
		sslResponseMap["status"] = responese.Status
	}

	if responese.StatusMsg != nil {
		sslResponseMap["status_msg"] = responese.StatusMsg
	}

	if responese.VerifyType != nil {
		sslResponseMap["verify_type"] = responese.VerifyType
	}

	if responese.VulnerabilityStatus != nil {
		sslResponseMap["vulnerability_status"] = responese.VulnerabilityStatus
	}

	if responese.CertBeginTime != nil {
		sslResponseMap["cert_begin_time"] = responese.CertBeginTime
	}

	if responese.CertEndTime != nil {
		sslResponseMap["cert_end_time"] = responese.CertEndTime
	}

	if responese.ValidityPeriod != nil {
		sslResponseMap["validity_period"] = responese.ValidityPeriod
	}

	if responese.InsertTime != nil {
		sslResponseMap["insert_time"] = responese.InsertTime
	}

	if responese.OrderId != nil {
		sslResponseMap["order_id"] = responese.OrderId
	}

	if responese.CertificateExtra != nil {
		certificateExtraMap := map[string]interface{}{}

		if responese.CertificateExtra.DomainNumber != nil {
			certificateExtraMap["domain_number"] = responese.CertificateExtra.DomainNumber
		}

		if responese.CertificateExtra.OriginCertificateId != nil {
			certificateExtraMap["origin_certificate_id"] = responese.CertificateExtra.OriginCertificateId
		}

		if responese.CertificateExtra.ReplacedBy != nil {
			certificateExtraMap["replaced_by"] = responese.CertificateExtra.ReplacedBy
		}

		if responese.CertificateExtra.ReplacedFor != nil {
			certificateExtraMap["replaced_for"] = responese.CertificateExtra.ReplacedFor
		}

		if responese.CertificateExtra.RenewOrder != nil {
			certificateExtraMap["renew_order"] = responese.CertificateExtra.RenewOrder
		}

		if responese.CertificateExtra.SMCert != nil {
			certificateExtraMap["s_m_cert"] = responese.CertificateExtra.SMCert
		}

		if responese.CertificateExtra.CompanyType != nil {
			certificateExtraMap["company_type"] = responese.CertificateExtra.CompanyType
		}

		sslResponseMap["certificate_extra"] = []interface{}{certificateExtraMap}
	}

	if responese.DvAuthDetail != nil {
		DvAuthDetailMap := map[string]interface{}{}

		if responese.DvAuthDetail.DvAuthKey != nil {
			DvAuthDetailMap["dv_auth_key"] = responese.DvAuthDetail.DvAuthKey
		}

		if responese.DvAuthDetail.DvAuthValue != nil {
			DvAuthDetailMap["dv_auth_value"] = responese.DvAuthDetail.DvAuthValue
		}

		if responese.DvAuthDetail.DvAuthDomain != nil {
			DvAuthDetailMap["dv_auth_domain"] = responese.DvAuthDetail.DvAuthDomain
		}

		if responese.DvAuthDetail.DvAuthPath != nil {
			DvAuthDetailMap["dv_auth_path"] = responese.DvAuthDetail.DvAuthPath
		}

		if responese.DvAuthDetail.DvAuthKeySubDomain != nil {
			DvAuthDetailMap["dv_auth_key_sub_domain"] = responese.DvAuthDetail.DvAuthKeySubDomain
		}

		if responese.DvAuthDetail.DvAuths != nil {
			dvAuthsList := []interface{}{}
			for _, dvAuths := range responese.DvAuthDetail.DvAuths {
				dvAuthsMap := map[string]interface{}{}

				if dvAuths.DvAuthKey != nil {
					dvAuthsMap["dv_auth_key"] = dvAuths.DvAuthKey
				}

				if dvAuths.DvAuthValue != nil {
					dvAuthsMap["dv_auth_value"] = dvAuths.DvAuthValue
				}

				if dvAuths.DvAuthDomain != nil {
					dvAuthsMap["dv_auth_domain"] = dvAuths.DvAuthDomain
				}

				if dvAuths.DvAuthPath != nil {
					dvAuthsMap["dv_auth_path"] = dvAuths.DvAuthPath
				}

				if dvAuths.DvAuthSubDomain != nil {
					dvAuthsMap["dv_auth_sub_domain"] = dvAuths.DvAuthSubDomain
				}

				if dvAuths.DvAuthVerifyType != nil {
					dvAuthsMap["dv_auth_verify_type"] = dvAuths.DvAuthVerifyType
				}

				dvAuthsList = append(dvAuthsList, dvAuthsMap)
			}

			DvAuthDetailMap["dv_auths"] = []interface{}{dvAuthsList}
		}

		sslResponseMap["dv_auth_detail"] = []interface{}{DvAuthDetailMap}
	}

	if responese.VulnerabilityReport != nil {
		sslResponseMap["vulnerability_report"] = responese.VulnerabilityReport
	}

	if responese.PackageTypeName != nil {
		sslResponseMap["package_type_name"] = responese.PackageTypeName
	}

	if responese.StatusName != nil {
		sslResponseMap["status_name"] = responese.StatusName
	}

	if responese.SubjectAltName != nil {
		sslResponseMap["subject_alt_name"] = responese.SubjectAltName
	}

	if responese.IsVip != nil {
		sslResponseMap["is_vip"] = responese.IsVip
	}

	if responese.IsWildcard != nil {
		sslResponseMap["is_wildcard"] = responese.IsWildcard
	}

	if responese.IsDv != nil {
		sslResponseMap["is_dv"] = responese.IsDv
	}

	if responese.IsVulnerability != nil {
		sslResponseMap["is_vulnerability"] = responese.IsVulnerability
	}

	if responese.RenewAble != nil {
		sslResponseMap["renew_able"] = responese.RenewAble
	}

	if responese.SubmittedData != nil {
		submittedDataMap := map[string]interface{}{}

		if responese.SubmittedData.CsrType != nil {
			submittedDataMap["csr_type"] = responese.SubmittedData.CsrType
		}

		if responese.SubmittedData.CsrContent != nil {
			submittedDataMap["csr_content"] = responese.SubmittedData.CsrContent
		}

		if responese.SubmittedData.CertificateDomain != nil {
			submittedDataMap["certificate_domain"] = responese.SubmittedData.CertificateDomain
		}

		if responese.SubmittedData.DomainList != nil {
			submittedDataMap["domain_list"] = responese.SubmittedData.DomainList
		}

		if responese.SubmittedData.KeyPassword != nil {
			submittedDataMap["key_password"] = responese.SubmittedData.KeyPassword
		}

		if responese.SubmittedData.OrganizationName != nil {
			submittedDataMap["organization_name"] = responese.SubmittedData.OrganizationName
		}

		if responese.SubmittedData.OrganizationDivision != nil {
			submittedDataMap["organization_division"] = responese.SubmittedData.OrganizationDivision
		}

		if responese.SubmittedData.OrganizationAddress != nil {
			submittedDataMap["organization_address"] = responese.SubmittedData.OrganizationAddress
		}

		if responese.SubmittedData.OrganizationCountry != nil {
			submittedDataMap["organization_country"] = responese.SubmittedData.OrganizationCountry
		}

		if responese.SubmittedData.OrganizationCity != nil {
			submittedDataMap["organization_city"] = responese.SubmittedData.OrganizationCity
		}

		if responese.SubmittedData.OrganizationRegion != nil {
			submittedDataMap["organization_region"] = responese.SubmittedData.OrganizationRegion
		}

		if responese.SubmittedData.PostalCode != nil {
			submittedDataMap["postal_code"] = responese.SubmittedData.PostalCode
		}

		if responese.SubmittedData.PhoneAreaCode != nil {
			submittedDataMap["phone_area_code"] = responese.SubmittedData.PhoneAreaCode
		}

		if responese.SubmittedData.PhoneNumber != nil {
			submittedDataMap["phone_number"] = responese.SubmittedData.PhoneNumber
		}

		if responese.SubmittedData.AdminFirstName != nil {
			submittedDataMap["admin_first_name"] = responese.SubmittedData.AdminFirstName
		}

		if responese.SubmittedData.AdminLastName != nil {
			submittedDataMap["admin_last_name"] = responese.SubmittedData.AdminLastName
		}

		if responese.SubmittedData.AdminPhoneNum != nil {
			submittedDataMap["admin_phone_num"] = responese.SubmittedData.AdminPhoneNum
		}

		if responese.SubmittedData.AdminEmail != nil {
			submittedDataMap["admin_email"] = responese.SubmittedData.AdminEmail
		}

		if responese.SubmittedData.AdminPosition != nil {
			submittedDataMap["admin_position"] = responese.SubmittedData.AdminPosition
		}

		if responese.SubmittedData.ContactFirstName != nil {
			submittedDataMap["contact_first_name"] = responese.SubmittedData.ContactFirstName
		}

		if responese.SubmittedData.ContactLastName != nil {
			submittedDataMap["contact_last_name"] = responese.SubmittedData.ContactLastName
		}

		if responese.SubmittedData.ContactNumber != nil {
			submittedDataMap["contact_number"] = responese.SubmittedData.ContactNumber
		}

		if responese.SubmittedData.ContactEmail != nil {
			submittedDataMap["contact_email"] = responese.SubmittedData.ContactEmail
		}

		if responese.SubmittedData.ContactPosition != nil {
			submittedDataMap["contact_position"] = responese.SubmittedData.ContactPosition
		}

		if responese.SubmittedData.VerifyType != nil {
			submittedDataMap["verify_type"] = responese.SubmittedData.VerifyType
		}

		sslResponseMap["submitted_data"] = []interface{}{submittedDataMap}
	}

	if responese.Deployable != nil {
		sslResponseMap["deployable"] = responese.Deployable
	}

	if responese.CAEncryptAlgorithms != nil {
		sslResponseMap["c_a_encrypt_algorithms"] = responese.CAEncryptAlgorithms
	}

	if responese.CACommonNames != nil {
		sslResponseMap["c_a_common_names"] = responese.CACommonNames
	}

	if responese.CAEndTimes != nil {
		sslResponseMap["c_a_end_times"] = responese.CAEndTimes
	}

	if responese.DvRevokeAuthDetail != nil {
		tmpList := []interface{}{}
		for _, dvAuths := range responese.DvRevokeAuthDetail {
			dvAuthsMap := map[string]interface{}{}

			if dvAuths.DvAuthKey != nil {
				dvAuthsMap["dv_auth_key"] = dvAuths.DvAuthKey
			}

			if dvAuths.DvAuthValue != nil {
				dvAuthsMap["dv_auth_value"] = dvAuths.DvAuthValue
			}

			if dvAuths.DvAuthDomain != nil {
				dvAuthsMap["dv_auth_domain"] = dvAuths.DvAuthDomain
			}

			if dvAuths.DvAuthPath != nil {
				dvAuthsMap["dv_auth_path"] = dvAuths.DvAuthPath
			}

			if dvAuths.DvAuthSubDomain != nil {
				dvAuthsMap["dv_auth_sub_domain"] = dvAuths.DvAuthSubDomain
			}

			if dvAuths.DvAuthVerifyType != nil {
				dvAuthsMap["dv_auth_verify_type"] = dvAuths.DvAuthVerifyType
			}

			tmpList = append(tmpList, dvAuthsMap)
		}

		sslResponseMap["dv_revoke_auth_detail"] = tmpList
	}
	_ = d.Set("result", []interface{}{sslResponseMap})
	d.SetId(CertificateId)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), sslResponseMap); e != nil {
			return e
		}
	}
	return nil
}
