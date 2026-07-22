package apigateway

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudAPIGatewayApiAppService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayApiAppServicesRead,
		Schema: map[string]*schema.Schema{
			"service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The unique ID service to be queried。",
			},
			"api_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Territory to which the service belongs。",
			},
			// computed
			"api_id_status_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "API list.注意：此字段可能返回 null，表示无法获取有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service unique ID。",
						},
						"api_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API unique ID。",
						},
						"api_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API Description注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API PATH。",
						},
						"method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API METHOD。",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service 创建时间。",
						},
						"modified_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service 修改时间。",
						},
						"api_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API 名称注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"uniq_vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC unique ID.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"api_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API 类型注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API 协议注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"is_debug_after_charge": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否debug after purchase.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"auth_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Authorization 类型注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"api_business_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API business 类型注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"auth_relation_api_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique ID association authorization API.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"oauth_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "OAuth configuration information.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"public_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Public 键，用于verify 用户 tokens。",
									},
									"token_location": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "令牌 passes the position。",
									},
									"login_redirect_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Redirect 地址，用于guide users in login operations。",
									},
								},
							},
						},
						"token_location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OAuth2.0 API request，令牌 storage location.注意：此字段可能返回 null，表示无法获取有效值。",
						},
					},
				},
			},
			"api_total_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Total 数量 APIs.注意：此字段可能返回 null，表示无法获取有效值。",
			},
			"available_environments": {
				Computed:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "列表 service environments.注意：此字段可能返回 null，表示无法获取有效值。",
			},
			"created_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Service 创建时间.注意：此字段可能返回 null，表示无法获取有效值。",
			},
			"inner_http_port": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Internal network access HTTP service 端口 number。",
			},
			"inner_https_port": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Internal network access https 端口 number。",
			},
			"internal_sub_domain": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Intranet access sub 域名 名称",
			},
			"ip_version": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "IP 版本注意：此字段可能返回 null，表示无法获取有效值。",
			},
			"modified_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Service 修改时间.注意：此字段可能返回 null，表示无法获取有效值。",
			},
			"net_types": {
				Computed:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A 列表 network types，where INNER represents internal network access and OUTER represents external network access。",
			},
			"outer_sub_domain": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "External network access sub 域名 名称",
			},
			"protocol": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Service support 协议，可选 values are http，https，and http&amp;amp;https。",
			},
			"service_desc": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Service 描述注意：此字段可能返回 null，表示无法获取有效值。",
			},
			"service_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "服务名称注意：此字段可能返回 null，表示无法获取有效值。",
			},
			"set_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Reserved fields.注意：此字段可能返回 null，表示无法获取有效值。",
			},
			"usage_plan_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Use a plan array.注意：此字段可能返回 null，表示无法获取有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"environment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Environment 名称",
						},
						"usage_plan_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Use a unique ID for the plan。",
						},
						"usage_plan_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Use the plan 名称",
						},
						"usage_plan_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Use plan 描述注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"max_request_num_pre_sec": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Use plan qps，-1 表示no restrictions。",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Use planned time。",
						},
						"modified_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Use the schedule to modify the time。",
						},
					},
				},
			},
			"usage_plan_total_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Total 数量 usage plans.注意：此字段可能返回 null，表示无法获取有效值。",
			},
			"user_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The 用户 类型 this service.注意：此字段可能返回 null，表示无法获取有效值。",
			},
			//"tags": {
			//	Type:        schema.TypeMap,
			//	Optional:    true,
			//	Description: "标签描述列表。",
			//},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudAPIGatewayApiAppServicesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_api_gateway_api_app_services.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service       = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		apiAppService *apigateway.DescribeServiceForApiAppResponseParams
		serviceId     string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("service_id"); ok {
		paramMap["ServiceId"] = helper.String(v.(string))
		serviceId = v.(string)
	}

	if v, ok := d.GetOk("api_region"); ok {
		paramMap["ApiRegion"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeAPIGatewayApiAppServiceByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		apiAppService = response
		return nil
	})

	if err != nil {
		return err
	}

	if apiAppService.ApiIdStatusSet != nil {
		tmpList := make([]map[string]interface{}, 0)
		for _, item := range apiAppService.ApiIdStatusSet {
			setMap := map[string]interface{}{}
			if item.ServiceId != nil {
				setMap["service_id"] = item.ServiceId
			}

			if item.ApiId != nil {
				setMap["api_id"] = item.ApiId
			}

			if item.ApiDesc != nil {
				setMap["api_desc"] = item.ApiDesc
			}

			if item.Path != nil {
				setMap["path"] = item.Path
			}

			if item.Method != nil {
				setMap["method"] = item.Method
			}

			if item.CreatedTime != nil {
				setMap["created_time"] = item.CreatedTime
			}

			if item.ModifiedTime != nil {
				setMap["modified_time"] = item.ModifiedTime
			}

			if item.ApiName != nil {
				setMap["api_name"] = item.ApiName
			}

			if item.UniqVpcId != nil {
				setMap["uniq_vpc_id"] = item.UniqVpcId
			}

			if item.ApiType != nil {
				setMap["api_type"] = item.ApiType
			}

			if item.Protocol != nil {
				setMap["protocol"] = item.Protocol
			}

			if item.IsDebugAfterCharge != nil {
				setMap["is_debug_after_charge"] = item.IsDebugAfterCharge
			}

			if item.AuthType != nil {
				setMap["auth_type"] = item.AuthType
			}

			if item.ApiBusinessType != nil {
				setMap["api_business_type"] = item.ApiBusinessType
			}

			if item.AuthRelationApiId != nil {
				setMap["auth_relation_api_id"] = item.AuthRelationApiId
			}

			if item.OauthConfig != nil {
				confList := make([]map[string]interface{}, 0)
				confMap := map[string]interface{}{}

				if item.OauthConfig.PublicKey != nil {
					confMap["public_key"] = item.OauthConfig.PublicKey
				}

				if item.OauthConfig.TokenLocation != nil {
					confMap["token_location"] = item.OauthConfig.TokenLocation
				}

				if item.OauthConfig.LoginRedirectUrl != nil {
					confMap["login_redirect_url"] = item.OauthConfig.LoginRedirectUrl
				}

				confList = append(confList, confMap)
				setMap["oauth_config"] = confList
			}

			if item.TokenLocation != nil {
				setMap["token_location"] = item.TokenLocation
			}

			tmpList = append(tmpList, setMap)
		}

		_ = d.Set("api_id_status_set", tmpList)
	}

	if apiAppService.ApiTotalCount != nil {
		_ = d.Set("api_total_count", apiAppService.ApiTotalCount)
	}

	if apiAppService.AvailableEnvironments != nil {
		_ = d.Set("available_environments", apiAppService.AvailableEnvironments)
	}

	if apiAppService.CreatedTime != nil {
		_ = d.Set("created_time", apiAppService.CreatedTime)
	}

	if apiAppService.InnerHttpPort != nil {
		_ = d.Set("inner_http_port", apiAppService.InnerHttpPort)
	}

	if apiAppService.InnerHttpsPort != nil {
		_ = d.Set("inner_https_port", apiAppService.InnerHttpsPort)
	}

	if apiAppService.InternalSubDomain != nil {
		_ = d.Set("internal_sub_domain", apiAppService.InternalSubDomain)
	}

	if apiAppService.IpVersion != nil {
		_ = d.Set("ip_version", apiAppService.IpVersion)
	}

	if apiAppService.ModifiedTime != nil {
		_ = d.Set("modified_time", apiAppService.ModifiedTime)
	}

	if apiAppService.NetTypes != nil {
		_ = d.Set("net_types", apiAppService.NetTypes)
	}

	if apiAppService.OuterSubDomain != nil {
		_ = d.Set("outer_sub_domain", apiAppService.OuterSubDomain)
	}

	if apiAppService.Protocol != nil {
		_ = d.Set("protocol", apiAppService.Protocol)
	}

	if apiAppService.ServiceDesc != nil {
		_ = d.Set("service_desc", apiAppService.ServiceDesc)
	}

	if apiAppService.ServiceName != nil {
		_ = d.Set("service_name", apiAppService.ServiceName)
	}

	if apiAppService.SetId != nil {
		_ = d.Set("set_id", apiAppService.SetId)
	}

	if apiAppService.UsagePlanList != nil {
		tmpList := make([]map[string]interface{}, 0)
		for _, usagePlan := range apiAppService.UsagePlanList {
			usagePlanMap := map[string]interface{}{}

			if usagePlan.Environment != nil {
				usagePlanMap["environment"] = usagePlan.Environment
			}

			if usagePlan.UsagePlanId != nil {
				usagePlanMap["usage_plan_id"] = usagePlan.UsagePlanId
			}

			if usagePlan.UsagePlanName != nil {
				usagePlanMap["usage_plan_name"] = usagePlan.UsagePlanName
			}

			if usagePlan.UsagePlanDesc != nil {
				usagePlanMap["usage_plan_desc"] = usagePlan.UsagePlanDesc
			}

			if usagePlan.MaxRequestNumPreSec != nil {
				usagePlanMap["max_request_num_pre_sec"] = usagePlan.MaxRequestNumPreSec
			}

			if usagePlan.CreatedTime != nil {
				usagePlanMap["created_time"] = usagePlan.CreatedTime
			}

			if usagePlan.ModifiedTime != nil {
				usagePlanMap["modified_time"] = usagePlan.ModifiedTime
			}

			tmpList = append(tmpList, usagePlanMap)
		}

		_ = d.Set("usage_plan_list", tmpList)
	}

	if apiAppService.UsagePlanTotalCount != nil {
		_ = d.Set("usage_plan_total_count", apiAppService.UsagePlanTotalCount)
	}

	if apiAppService.UserType != nil {
		_ = d.Set("user_type", apiAppService.UserType)
	}

	if apiAppService.ServiceDesc != nil {
		_ = d.Set("service_desc", apiAppService.ServiceDesc)
	}

	d.SetId(serviceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
