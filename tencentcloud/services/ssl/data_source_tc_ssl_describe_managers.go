package ssl

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSslDescribeManagers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslDescribeManagersRead,
		Schema: map[string]*schema.Schema{
			"company_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "公司 ID",
			},

			"manager_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Manager&amp;#39;s 名称 (将 是 abandoned)，please 使用 Searchkey。",
			},

			"manager_mail": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Vague 查询 manager email (将 是 abandoned)，please 使用 Searchkey。",
			},

			"status": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "过滤器 according 到 状态 manager，和 值 是 可用&amp;#39;None&amp;#39; Unable 到 提交 review&amp;#39;Audit&amp;#39;，Asian Credit Review&amp;#39;Caaudit&amp;#39; CA review&amp;#39;OK&amp;#39; has been reviewed&amp;#39;Invalid&amp;#39; review failed&amp;#39;Expiring&amp;#39; 是 about 到 expire&amp;#39;Expired&amp;#39; expired。",
			},

			"search_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Manager&amp;#39;s surname/Manager 名称/mailbox/department precise matching。",
			},

			"managers": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "公司管理员 List。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态: Audit: OK during review: review passed inValid: expired expiRing: 是 about 到 expire Expired: expired。",
						},
						"manager_first_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Manager 名称",
						},
						"manager_last_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Manager 名称",
						},
						"manager_position": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Manager position。",
						},
						"manager_phone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Manager phone call。",
						},
						"manager_mail": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Manager mailbox。",
						},
						"manager_department": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Administrator department。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation timeNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},
						"domain_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 administrators。",
						},
						"cert_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 administrative certificates。",
						},
						"manager_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Manager ID。",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Examine validity expiration timeNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},
						"submit_audit_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "last 时间 review timeNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
						},
						"verify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Examination timeNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
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

func dataSourceTencentCloudSslDescribeManagersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ssl_describe_managers.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var companyId string
	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("company_id"); v != nil {
		companyId = helper.IntToStr(v.(int))
		paramMap["CompanyId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("manager_name"); ok {
		paramMap["ManagerName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("manager_mail"); ok {
		paramMap["ManagerMail"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		paramMap["Status"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_key"); ok {
		paramMap["SearchKey"] = helper.String(v.(string))
	}

	service := SslService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var managers []*ssl.ManagerInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSslDescribeManagersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		managers = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(managers))

	if managers != nil {
		for _, managerInfo := range managers {
			managerInfoMap := map[string]interface{}{}

			if managerInfo.Status != nil {
				managerInfoMap["status"] = managerInfo.Status
			}

			if managerInfo.ManagerFirstName != nil {
				managerInfoMap["manager_first_name"] = managerInfo.ManagerFirstName
			}

			if managerInfo.ManagerLastName != nil {
				managerInfoMap["manager_last_name"] = managerInfo.ManagerLastName
			}

			if managerInfo.ManagerPosition != nil {
				managerInfoMap["manager_position"] = managerInfo.ManagerPosition
			}

			if managerInfo.ManagerPhone != nil {
				managerInfoMap["manager_phone"] = managerInfo.ManagerPhone
			}

			if managerInfo.ManagerMail != nil {
				managerInfoMap["manager_mail"] = managerInfo.ManagerMail
			}

			if managerInfo.ManagerDepartment != nil {
				managerInfoMap["manager_department"] = managerInfo.ManagerDepartment
			}

			if managerInfo.CreateTime != nil {
				managerInfoMap["create_time"] = managerInfo.CreateTime
			}

			if managerInfo.DomainCount != nil {
				managerInfoMap["domain_count"] = managerInfo.DomainCount
			}

			if managerInfo.CertCount != nil {
				managerInfoMap["cert_count"] = managerInfo.CertCount
			}

			if managerInfo.ManagerId != nil {
				managerInfoMap["manager_id"] = managerInfo.ManagerId
			}

			if managerInfo.ExpireTime != nil {
				managerInfoMap["expire_time"] = managerInfo.ExpireTime
			}

			if managerInfo.SubmitAuditTime != nil {
				managerInfoMap["submit_audit_time"] = managerInfo.SubmitAuditTime
			}

			if managerInfo.VerifyTime != nil {
				managerInfoMap["verify_time"] = managerInfo.VerifyTime
			}

			tmpList = append(tmpList, managerInfoMap)
		}

		_ = d.Set("managers", tmpList)
	}

	d.SetId(companyId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
