package apigateway

import (
	"context"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

func DataSourceTencentCloudAPIGatewayAPIs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayAPIsRead,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "服务 ID for query。",
			},
			"api_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom API 名称",
			},
			"api_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Created API ID。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			// Computed values.
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 APIs。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Which service this API belongs. Refer to resource `tencentcloud_api_gateway_service`。",
						},
						"api_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom API 名称",
						},
						"api_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom API 描述",
						},
						"auth_type": {
							Type:     schema.TypeString,
							Computed: true,
							Description: "API 认证类型。有效值：“秘密”、“无”。" +
								"`SECRET` means key pair authentication, `NONE` means no authentication.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API frontend request 类型，such as `HTTP`,`WEBSOCKET`。",
						},
						"enable_cors": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否enable CORS。",
						},
						"request_config_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Request frontend 路径 configuration. Like `/用户/getinfo`。",
						},
						"request_config_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Request frontend method configuration. Like `GET`,`POST`,`PUT`,`DELETE`,`HEAD`,`ANY`。",
						},
						"request_parameters": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Frontend request parameters。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Parameter 名称",
									},
									"position": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Parameter location。",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Parameter 类型",
									},
									"desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Parameter 描述",
									},
									"default_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Parameter 默认值",
									},
									"required": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "If this parameter 必填",
									},
								},
							},
						},
						"service_config_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API backend service 类型",
						},
						"service_config_timeout": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "API backend service timeout 周期 （秒）。",
						},
						"service_config_product": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backend 类型 This parameter takes effect when VPC is 已启用 Currently，only `clb` is supported。",
						},
						"service_config_vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique 私有网络 ID",
						},
						"service_config_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API backend service URL This parameter 为必填项 when `service_config_type` is `HTTP`。",
						},
						"service_config_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API backend service 路径，such as /路径 If `service_config_type` is `HTTP`，this parameter will be 必填 The frontend `request_config_path` and backend 路径 `service_config_path` can be different。",
						},
						"service_config_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API backend service request method，such as `GET`. If `service_config_type` is `HTTP`，this parameter will be 必填 The frontend `request_config_method` and backend method `service_config_method` can be different。",
						},
						"service_config_scf_function_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SCF function 名称 This parameter takes effect when `service_config_type` is `SCF`。",
						},
						"service_config_scf_function_namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SCF function namespace. This parameter takes effect when  `service_config_type` is `SCF`。",
						},
						"service_config_scf_function_qualifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SCF function 版本 This parameter takes effect when `service_config_type`  is `SCF`。",
						},
						"service_config_mock_return_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Returned information of API backend mocking。",
						},
						"response_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Return 类型",
						},
						"response_success_example": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Successful response sample of custom response configuration。",
						},
						"response_fail_example": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Response failure sample of custom response configuration。",
						},
						"response_error_codes": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Custom 错误码 configuration. Must keep at least one after set。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Custom response configuration 错误码",
									},
									"msg": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom response configuration 错误信息",
									},
									"desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Parameter 描述",
									},
									"converted_code": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Custom 错误码 conversion。",
									},
									"need_convert": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否enable 错误码 conversion. 默认值：`false`。",
									},
								},
							},
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "最后修改时间 in the 格式 of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 in the 格式 of YYYY-MM-DDThh:mm:ssZ according to ISO 8601 standard. UTC time is used。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAPIGatewayAPIsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_api_gateway_apis.read")()

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		apiName           = d.Get("api_name").(string)
		apiId             = d.Get("api_id").(string)
		serviceId         = d.Get("service_id").(string)
		apiSet            []*apigateway.DescribeApisStatusResultApiIdStatusSetInfo
		err               error
	)

	if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		apiSet, err = apiGatewayService.DescribeApisStatus(ctx, serviceId, apiName, apiId)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(apiSet))
	for _, apiKey := range apiSet {
		var (
			info apigateway.ApiInfo
			has  bool
			item = make(map[string]interface{})
		)
		if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			info, has, err = apiGatewayService.DescribeApi(ctx, *apiKey.ServiceId, *apiKey.ApiId)
			if err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
		if !has {
			continue
		}

		item["service_id"] = info.ServiceId
		item["api_name"] = info.ApiName
		item["api_desc"] = info.ApiDesc
		item["auth_type"] = info.AuthType
		item["protocol"] = info.Protocol
		item["enable_cors"] = info.EnableCORS
		item["response_type"] = info.ResponseType
		item["response_success_example"] = info.ResponseSuccessExample
		item["response_fail_example"] = info.ResponseFailExample
		item["service_config_type"] = info.ServiceType
		item["service_config_timeout"] = info.ServiceTimeout
		item["service_config_scf_function_name"] = info.ServiceScfFunctionName
		item["service_config_scf_function_namespace"] = info.ServiceScfFunctionNamespace
		item["service_config_scf_function_qualifier"] = info.ServiceScfFunctionQualifier
		item["service_config_mock_return_message"] = info.ServiceMockReturnMessage
		item["modify_time"] = info.ModifiedTime
		item["create_time"] = info.CreatedTime

		if info.RequestConfig != nil {
			item["request_config_path"] = info.RequestConfig.Path
			item["request_config_method"] = info.RequestConfig.Method
		} else {
			item["request_config_path"] = ""
			item["request_config_method"] = ""
		}

		paramList := make([]map[string]interface{}, 0, len(info.RequestParameters))
		if info.RequestParameters != nil {
			for _, param := range info.RequestParameters {
				paramList = append(paramList, map[string]interface{}{
					"name":          param.Name,
					"position":      param.Position,
					"type":          param.Type,
					"desc":          param.Desc,
					"default_value": param.DefaultValue,
					"required":      param.Required,
				})
			}
		}
		item["request_parameters"] = paramList

		if info.ServiceConfig != nil {
			item["service_config_product"] = info.ServiceConfig.Product
			item["service_config_vpc_id"] = info.ServiceConfig.UniqVpcId
			item["service_config_url"] = info.ServiceConfig.Url
			item["service_config_path"] = info.ServiceConfig.Path
			item["service_config_method"] = info.ServiceConfig.Method
		} else {
			item["service_config_product"] = ""
			item["service_config_vpc_id"] = ""
			item["service_config_url"] = ""
			item["service_config_path"] = ""
			item["service_config_method"] = ""
		}

		codeList := make([]map[string]interface{}, 0, len(info.ResponseErrorCodes))
		if info.ResponseErrorCodes != nil {
			for _, code := range info.ResponseErrorCodes {
				codeList = append(codeList, map[string]interface{}{
					"code":           code.Code,
					"msg":            code.Msg,
					"desc":           code.Desc,
					"converted_code": code.ConvertedCode,
					"need_convert":   code.NeedConvert,
				})
			}
		}
		item["response_error_codes"] = codeList

		list = append(list, item)
	}

	if err = d.Set("list", list); err != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
		return err
	}

	d.SetId(strings.Join([]string{apiName, apiId}, tccommon.FILED_SP))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}
	return nil
}
