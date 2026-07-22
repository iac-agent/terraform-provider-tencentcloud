package apigateway

import (
	"context"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudApiGatewayApiAppApi() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudApiGatewayApiAppApiRead,
		Schema: map[string]*schema.Schema{
			"service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "唯一 ID 服务 其中 API resides。",
			},
			"api_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "API interface 唯一 ID。",
			},
			"api_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Api 地域",
			},
			// computed
			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "API details。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "唯一 ID 服务 其中 API resides。",
						},
						"service_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 服务 其中 API resides。",
						},
						"service_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A 描述 服务 其中 API resides。",
						},
						"api_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API interface 唯一 ID。",
						},
						"api_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 API interface。",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间，expressed 在 accordance 使用 ISO8601 standard 和 使用 UTC 时间. 格式 是: YYYY-MM-DDThh:mm:ssZ。",
						},
						"modified_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last 修改时间，expressed 在 accordance 使用 ISO8601 standard 和 使用 UTC 时间. 格式 是: YYYY-MM-DDThh:mm:ssZ。",
						},
						"api_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 API interface。",
						},
						"api_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API 类型 Possible 值 是 NORMAL (normal API) 和 TSF (microservice API)。",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "front-end 请求 类型 API，such 作为 HTTP 或 HTTPS 或 HTTP 和 HTTPS。",
						},
						"auth_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API authentication 类型 Possible 值 是 SECRET (键 pair authentication)，NONE (authentication-free)，和 OAUTH。",
						},
						"api_business_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 OAUTH API. Possible 值 是 NORMAL (Business API)，OAUTH (Authorization API)。",
						},
						"auth_relation_api_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OAUTH 唯一 ID authorization API associated 使用 business API。",
						},
						"oauth_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "OAUTH 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"public_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Public 键，用于verify 用户 令牌",
									},
									"token_location": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "令牌 delivery position。",
									},
									"login_redirect_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Redirect 地址，用于guide users 到 日志 在。",
									},
								},
							},
						},
						"is_debug_after_charge": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否debug after purchase (参数 reserved 在 云 market)。",
						},
						"request_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "requested frontend 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API 路径，such 作为 /路径",
									},
									"method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API 请求 方法，such 作为 GET。",
									},
								},
							},
						},
						"response_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Return 类型",
						},
						"response_success_example": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom response 配置 successful response 示例。",
						},
						"response_fail_example": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom response 配置 failure response 示例。",
						},
						"response_error_codes": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "用户-defined 错误码 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Custom response 配置 错误码",
									},
									"msg": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom response 配置 错误信息",
									},
									"desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom response 配置 错误码 备注",
									},
									"converted_code": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Custom 错误码 conversion。",
									},
									"need_convert": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否为necessary 到 启用 错误码 conversion。",
									},
								},
							},
						},
						"request_parameters": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Front-end 请求 参数。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API front-end 参数 名称",
									},
									"position": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "front-end 参数 position 的 API，such 作为 头部. Currently 支持 头部，查询，路径",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API front-end 参数 类型，such 作为 String，int。",
									},
									"default_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API front-end 参数 默认值",
									},
									"required": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "。",
									},
									"desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API front-end 参数 备注",
									},
								},
							},
						},
						"service_timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "backend 服务 超时 的 API，（秒）。",
						},
						"service_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "backend 服务 类型 API. Possible 值 是 HTTP，MOCK，TSF，CLB，SCF，WEBSOCKET，和 TARGET (内部 testing)。",
						},
						"service_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Backend 服务 配置 对于 API。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"product": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Backend 类型 It takes effect 当 vpc 是 已启用 Currently 支持 types 是 clb，cvm 和 upstream。",
									},
									"uniq_vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "唯一 ID vpc。",
									},
									"url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API&amp;#39;s backend 服务 URL 如果 ServiceType 是 HTTP，此 参数 必须 是 passed。",
									},
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API backend 服务 路径，such 作为 /路径 如果 ServiceType 是 HTTP，此 参数 为必填项. front-end 和 back-end paths 可以 是 different。",
									},
									"method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API backend 服务 请求 方法，such 作为 GET. 如果 ServiceType 是 HTTP，此 参数 为必填项. front-end 和 back-end methods 可以 是 different。",
									},
									"upstream_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Only 必填 当 binding vpc channel。",
									},
								},
							},
						},
						"service_parameters": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "API backend 服务 参数。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "backend 服务 参数 名称 API. 此 参数 将 是 使用 仅 如果 ServiceType 是 HTTP. front-end 和 back-end 参数 names 可以 是 different。",
									},
									"position": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "backend 服务 参数 location 的 API，such 作为 head. 此 参数 是 仅 使用 如果 ServiceType 是 HTTP. front-end 和 back-end 参数 positions 可以 是 已配置 differently。",
									},
									"relevant_request_parameter_position": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "front-end 参数 position corresponding 到 back-end 服务 参数 的 API，such 作为 head. 此 参数 是 仅 使用 如果 ServiceType 是 HTTP。",
									},
									"relevant_request_parameter_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "front-end 参数 名称 corresponding 到 back-end 服务 参数 的 API. 此 参数 是 仅 使用 如果 ServiceType 是 HTTP。",
									},
									"default_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Default 值 对于 APIs backend 服务 参数. 此 参数 是 仅 使用 如果 ServiceType 是 HTTP。",
									},
									"relevant_request_parameter_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "备注 在 backend 服务 参数 的 API. 此 参数 是 仅 使用 如果 ServiceType 是 HTTP。",
									},
								},
							},
						},
						"constant_parameters": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Constant 参数。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Constant 参数 名称 此 参数 是 仅 使用 如果 ServiceType 是 HTTP。",
									},
									"desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Constant 参数 描述 此 参数 是 仅 使用 如果 ServiceType 是 HTTP。",
									},
									"position": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Constant 参数 position. 此 参数 是 仅 使用 如果 ServiceType 是 HTTP。",
									},
									"default_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Constant 参数 默认值 此 参数 是 仅 使用 如果 ServiceType 是 HTTP。",
									},
								},
							},
						},
						"service_mock_return_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "APIs backend Mock 返回information. 如果 ServiceType 是 Mock，此 参数 必须 是 passed。",
						},
						"service_scf_function_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf 函数 名称 Effective 当 backend 类型 是 SCF。",
						},
						"service_scf_function_namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf 函数 命名空间. Effective 当 backend 类型 是 SCF。",
						},
						"service_scf_function_qualifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf 函数 版本 Effective 当 backend 类型 是 SCF。",
						},
						"service_scf_is_integrated_response": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否enable integrated response。",
						},
						"service_websocket_register_function_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf websocket registration 函数 命名空间. 有效 当 front-end 类型 是 WEBSOCKET 和 back-end 类型 是 SCF。",
						},
						"service_websocket_register_function_namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf websocket registration 函数 命名空间. 有效 当 front-end 类型 是 WEBSOCKET 和 back-end 类型 是 SCF。",
						},
						"service_websocket_register_function_qualifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf websocket transfer 函数 版本 有效 当 front-end 类型 是 WEBSOCKET 和 back-end 类型 是 SCF。",
						},
						"service_websocket_cleanup_function_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf websocket cleaning 函数. 有效 当 front-end 类型 是 WEBSOCKET 和 back-end 类型 是 SCF。",
						},
						"service_websocket_cleanup_function_namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf websocket cleanup 函数 命名空间. 有效 当 front-end 类型 是 WEBSOCKET 和 back-end 类型 是 SCF。",
						},
						"service_websocket_cleanup_function_qualifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf websocket cleanup 函数 版本 有效 当 front-end 类型 是 WEBSOCKET 和 back-end 类型 是 SCF。",
						},
						"internal_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "WEBSOCKET pushback 地址",
						},
						"service_websocket_transport_function_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf websocket transfer 函数. 有效 当 front-end 类型 是 WEBSOCKET 和 back-end 类型 是 SCF。",
						},
						"service_websocket_transport_function_namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf websocket transfer 函数 命名空间. 有效 当 front-end 类型 是 WEBSOCKET 和 back-end 类型 是 SCF。",
						},
						"service_websocket_transport_function_qualifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scf websocket transfer 函数 版本 有效 当 front-end 类型 是 WEBSOCKET 和 back-end 类型 是 SCF。",
						},
						"micro_services": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "API binding microservice 列表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Microservice 集群 ID。",
									},
									"namespace_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Microservice 命名空间 ID。",
									},
									"micro_service_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Microservice 名称",
									},
								},
							},
						},
						"micro_services_info": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "Microservice 信息 details。",
						},
						"service_tsf_load_balance_conf": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Load balancing 配置 对于 microservices。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_load_balance": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否enable load balancing。",
									},
									"method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Load balancing 方法。",
									},
									"session_stick_required": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否enable 会话 persistence。",
									},
									"session_stick_timeout": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Session retention 超时。",
									},
								},
							},
						},
						"service_tsf_health_check_conf": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Health check 配置 对于 microservices。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_health_check": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否enable health check。",
									},
									"request_volume_threshold": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Health check 阈值。",
									},
									"sleep_window_in_milliseconds": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Window 大小。",
									},
									"error_threshold_percentage": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Threshold percentage。",
									},
								},
							},
						},
						"enable_cors": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否enable cross-域名",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "API binding 标签 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "键 的 标签",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "值 的 note。",
									},
								},
							},
						},
						"environments": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "API published 环境 信息。",
						},
						"is_base64_encoded": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否enable Base64 编码 将 仅 take effect 当 backend 是 scf。",
						},
						"is_base64_trigger": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否enable Base64-encoded 头部 triggering 将 仅 take effect 当 backend 是 scf。",
						},
						"base64_encoded_trigger_rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Header triggers 规则，和 总数 数量 规则 does 不 exceed 10。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Header 对于 编码 triggering，可选 值 Accept 和 Content_Type correspond 到 Accept 和 内容-类型 在 actual 数据 flow 请求 头部。",
									},
									"value": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "An 数组 可选 值 对于 头部 triggered 通过 编码. 最大 字符串 长度 的 数组 element 是 40. elements 可以 include numbers，English letters 和 special 字符. 可选 值 对于 special 字符 是: `.` `+` ` *` `-` `/` `_` For 示例 [ 应用/x-vpeg005，应用/xhtml+xml，应用/vnd.ms -项目，应用/vnd.rn-rn_music_package] etc. 是 all legal。",
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
				Description: "用于save apiAppApis。",
			},
		},
	}
}

func dataSourceTencentCloudApiGatewayApiAppApiRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_api_gateway_api_app_api.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		apiAppApi  *apigateway.ApiInfo
		service_id string
		api_id     string
		api_region string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("service_id"); ok {
		paramMap["ServiceId"] = helper.String(v.(string))
		service_id = v.(string)
	}

	if v, ok := d.GetOk("api_id"); ok {
		paramMap["APIId"] = helper.String(v.(string))
		api_id = v.(string)
	}

	if v, ok := d.GetOk("api_region"); ok {
		paramMap["ApiRegion"] = helper.String(v.(string))
		api_region = v.(string)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeApiGatewayApiAppApiByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		apiAppApi = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0)
	if apiAppApi != nil {
		apiInfoMap := map[string]interface{}{}

		if apiAppApi.ServiceId != nil {
			apiInfoMap["service_id"] = apiAppApi.ServiceId
		}

		if apiAppApi.ServiceName != nil {
			apiInfoMap["service_name"] = apiAppApi.ServiceName
		}

		if apiAppApi.ServiceDesc != nil {
			apiInfoMap["service_desc"] = apiAppApi.ServiceDesc
		}

		if apiAppApi.ApiId != nil {
			apiInfoMap["api_id"] = apiAppApi.ApiId
		}

		if apiAppApi.ApiDesc != nil {
			apiInfoMap["api_desc"] = apiAppApi.ApiDesc
		}

		if apiAppApi.CreatedTime != nil {
			apiInfoMap["created_time"] = apiAppApi.CreatedTime
		}

		if apiAppApi.ModifiedTime != nil {
			apiInfoMap["modified_time"] = apiAppApi.ModifiedTime
		}

		if apiAppApi.ApiName != nil {
			apiInfoMap["api_name"] = apiAppApi.ApiName
		}

		if apiAppApi.ApiType != nil {
			apiInfoMap["api_type"] = apiAppApi.ApiType
		}

		if apiAppApi.Protocol != nil {
			apiInfoMap["protocol"] = apiAppApi.Protocol
		}

		if apiAppApi.AuthType != nil {
			apiInfoMap["auth_type"] = apiAppApi.AuthType
		}

		if apiAppApi.ApiBusinessType != nil {
			apiInfoMap["api_business_type"] = apiAppApi.ApiBusinessType
		}

		if apiAppApi.AuthRelationApiId != nil {
			apiInfoMap["auth_relation_api_id"] = apiAppApi.AuthRelationApiId
		}

		if apiAppApi.OauthConfig != nil {
			oauthConfigMap := map[string]interface{}{}

			if apiAppApi.OauthConfig.PublicKey != nil {
				oauthConfigMap["public_key"] = apiAppApi.OauthConfig.PublicKey
			}

			if apiAppApi.OauthConfig.TokenLocation != nil {
				oauthConfigMap["token_location"] = apiAppApi.OauthConfig.TokenLocation
			}

			if apiAppApi.OauthConfig.LoginRedirectUrl != nil {
				oauthConfigMap["login_redirect_url"] = apiAppApi.OauthConfig.LoginRedirectUrl
			}

			apiInfoMap["oauth_config"] = []interface{}{oauthConfigMap}
		}

		if apiAppApi.IsDebugAfterCharge != nil {
			apiInfoMap["is_debug_after_charge"] = apiAppApi.IsDebugAfterCharge
		}

		if apiAppApi.RequestConfig != nil {
			requestConfigMap := map[string]interface{}{}

			if apiAppApi.RequestConfig.Path != nil {
				requestConfigMap["path"] = apiAppApi.RequestConfig.Path
			}

			if apiAppApi.RequestConfig.Method != nil {
				requestConfigMap["method"] = apiAppApi.RequestConfig.Method
			}

			apiInfoMap["request_config"] = []interface{}{requestConfigMap}
		}

		if apiAppApi.ResponseType != nil {
			apiInfoMap["response_type"] = apiAppApi.ResponseType
		}

		if apiAppApi.ResponseSuccessExample != nil {
			apiInfoMap["response_success_example"] = apiAppApi.ResponseSuccessExample
		}

		if apiAppApi.ResponseFailExample != nil {
			apiInfoMap["response_fail_example"] = apiAppApi.ResponseFailExample
		}

		if apiAppApi.ResponseErrorCodes != nil {
			responseErrorCodesList := []interface{}{}
			for _, responseErrorCodes := range apiAppApi.ResponseErrorCodes {
				responseErrorCodesMap := map[string]interface{}{}

				if responseErrorCodes.Code != nil {
					responseErrorCodesMap["code"] = responseErrorCodes.Code
				}

				if responseErrorCodes.Msg != nil {
					responseErrorCodesMap["msg"] = responseErrorCodes.Msg
				}

				if responseErrorCodes.Desc != nil {
					responseErrorCodesMap["desc"] = responseErrorCodes.Desc
				}

				if responseErrorCodes.ConvertedCode != nil {
					responseErrorCodesMap["converted_code"] = responseErrorCodes.ConvertedCode
				}

				if responseErrorCodes.NeedConvert != nil {
					responseErrorCodesMap["need_convert"] = responseErrorCodes.NeedConvert
				}

				responseErrorCodesList = append(responseErrorCodesList, responseErrorCodesMap)
			}

			apiInfoMap["response_error_codes"] = responseErrorCodesList
		}

		if apiAppApi.RequestParameters != nil {
			requestParametersList := []interface{}{}
			for _, requestParameters := range apiAppApi.RequestParameters {
				requestParametersMap := map[string]interface{}{}

				if requestParameters.Name != nil {
					requestParametersMap["name"] = requestParameters.Name
				}

				if requestParameters.Position != nil {
					requestParametersMap["position"] = requestParameters.Position
				}

				if requestParameters.Type != nil {
					requestParametersMap["type"] = requestParameters.Type
				}

				if requestParameters.DefaultValue != nil {
					requestParametersMap["default_value"] = requestParameters.DefaultValue
				}

				if requestParameters.Required != nil {
					requestParametersMap["required"] = requestParameters.Required
				}

				if requestParameters.Desc != nil {
					requestParametersMap["desc"] = requestParameters.Desc
				}

				requestParametersList = append(requestParametersList, requestParametersMap)
			}

			apiInfoMap["request_parameters"] = requestParametersList
		}

		if apiAppApi.ServiceTimeout != nil {
			apiInfoMap["service_timeout"] = apiAppApi.ServiceTimeout
		}

		if apiAppApi.ServiceType != nil {
			apiInfoMap["service_type"] = apiAppApi.ServiceType
		}

		if apiAppApi.ServiceConfig != nil {
			serviceConfigMap := map[string]interface{}{}

			if apiAppApi.ServiceConfig.Product != nil {
				serviceConfigMap["product"] = apiAppApi.ServiceConfig.Product
			}

			if apiAppApi.ServiceConfig.UniqVpcId != nil {
				serviceConfigMap["uniq_vpc_id"] = apiAppApi.ServiceConfig.UniqVpcId
			}

			if apiAppApi.ServiceConfig.Url != nil {
				serviceConfigMap["url"] = apiAppApi.ServiceConfig.Url
			}

			if apiAppApi.ServiceConfig.Path != nil {
				serviceConfigMap["path"] = apiAppApi.ServiceConfig.Path
			}

			if apiAppApi.ServiceConfig.Method != nil {
				serviceConfigMap["method"] = apiAppApi.ServiceConfig.Method
			}

			if apiAppApi.ServiceConfig.UpstreamId != nil {
				serviceConfigMap["upstream_id"] = apiAppApi.ServiceConfig.UpstreamId
			}

			apiInfoMap["service_config"] = []interface{}{serviceConfigMap}
		}

		if apiAppApi.ServiceParameters != nil {
			serviceParametersList := []interface{}{}
			for _, serviceParameters := range apiAppApi.ServiceParameters {
				serviceParametersMap := map[string]interface{}{}

				if serviceParameters.Name != nil {
					serviceParametersMap["name"] = serviceParameters.Name
				}

				if serviceParameters.Position != nil {
					serviceParametersMap["position"] = serviceParameters.Position
				}

				if serviceParameters.RelevantRequestParameterPosition != nil {
					serviceParametersMap["relevant_request_parameter_position"] = serviceParameters.RelevantRequestParameterPosition
				}

				if serviceParameters.RelevantRequestParameterName != nil {
					serviceParametersMap["relevant_request_parameter_name"] = serviceParameters.RelevantRequestParameterName
				}

				if serviceParameters.DefaultValue != nil {
					serviceParametersMap["default_value"] = serviceParameters.DefaultValue
				}

				if serviceParameters.RelevantRequestParameterDesc != nil {
					serviceParametersMap["relevant_request_parameter_desc"] = serviceParameters.RelevantRequestParameterDesc
				}

				serviceParametersList = append(serviceParametersList, serviceParametersMap)
			}

			apiInfoMap["service_parameters"] = serviceParametersList
		}

		if apiAppApi.ConstantParameters != nil {
			constantParametersList := []interface{}{}
			for _, constantParameters := range apiAppApi.ConstantParameters {
				constantParametersMap := map[string]interface{}{}

				if constantParameters.Name != nil {
					constantParametersMap["name"] = constantParameters.Name
				}

				if constantParameters.Desc != nil {
					constantParametersMap["desc"] = constantParameters.Desc
				}

				if constantParameters.Position != nil {
					constantParametersMap["position"] = constantParameters.Position
				}

				if constantParameters.DefaultValue != nil {
					constantParametersMap["default_value"] = constantParameters.DefaultValue
				}

				constantParametersList = append(constantParametersList, constantParametersMap)
			}

			apiInfoMap["constant_parameters"] = constantParametersList
		}

		if apiAppApi.ServiceMockReturnMessage != nil {
			apiInfoMap["service_mock_return_message"] = apiAppApi.ServiceMockReturnMessage
		}

		if apiAppApi.ServiceScfFunctionName != nil {
			apiInfoMap["service_scf_function_name"] = apiAppApi.ServiceScfFunctionName
		}

		if apiAppApi.ServiceScfFunctionNamespace != nil {
			apiInfoMap["service_scf_function_namespace"] = apiAppApi.ServiceScfFunctionNamespace
		}

		if apiAppApi.ServiceScfFunctionQualifier != nil {
			apiInfoMap["service_scf_function_qualifier"] = apiAppApi.ServiceScfFunctionQualifier
		}

		if apiAppApi.ServiceScfIsIntegratedResponse != nil {
			apiInfoMap["service_scf_is_integrated_response"] = apiAppApi.ServiceScfIsIntegratedResponse
		}

		if apiAppApi.ServiceWebsocketRegisterFunctionName != nil {
			apiInfoMap["service_websocket_register_function_name"] = apiAppApi.ServiceWebsocketRegisterFunctionName
		}

		if apiAppApi.ServiceWebsocketRegisterFunctionNamespace != nil {
			apiInfoMap["service_websocket_register_function_namespace"] = apiAppApi.ServiceWebsocketRegisterFunctionNamespace
		}

		if apiAppApi.ServiceWebsocketRegisterFunctionQualifier != nil {
			apiInfoMap["service_websocket_register_function_qualifier"] = apiAppApi.ServiceWebsocketRegisterFunctionQualifier
		}

		if apiAppApi.ServiceWebsocketCleanupFunctionName != nil {
			apiInfoMap["service_websocket_cleanup_function_name"] = apiAppApi.ServiceWebsocketCleanupFunctionName
		}

		if apiAppApi.ServiceWebsocketCleanupFunctionNamespace != nil {
			apiInfoMap["service_websocket_cleanup_function_namespace"] = apiAppApi.ServiceWebsocketCleanupFunctionNamespace
		}

		if apiAppApi.ServiceWebsocketCleanupFunctionQualifier != nil {
			apiInfoMap["service_websocket_cleanup_function_qualifier"] = apiAppApi.ServiceWebsocketCleanupFunctionQualifier
		}

		if apiAppApi.InternalDomain != nil {
			apiInfoMap["internal_domain"] = apiAppApi.InternalDomain
		}

		if apiAppApi.ServiceWebsocketTransportFunctionName != nil {
			apiInfoMap["service_websocket_transport_function_name"] = apiAppApi.ServiceWebsocketTransportFunctionName
		}

		if apiAppApi.ServiceWebsocketTransportFunctionNamespace != nil {
			apiInfoMap["service_websocket_transport_function_namespace"] = apiAppApi.ServiceWebsocketTransportFunctionNamespace
		}

		if apiAppApi.ServiceWebsocketTransportFunctionQualifier != nil {
			apiInfoMap["service_websocket_transport_function_qualifier"] = apiAppApi.ServiceWebsocketTransportFunctionQualifier
		}

		if apiAppApi.MicroServices != nil {
			microServicesList := []interface{}{}
			for _, microServices := range apiAppApi.MicroServices {
				microServicesMap := map[string]interface{}{}

				if microServices.ClusterId != nil {
					microServicesMap["cluster_id"] = microServices.ClusterId
				}

				if microServices.NamespaceId != nil {
					microServicesMap["namespace_id"] = microServices.NamespaceId
				}

				if microServices.MicroServiceName != nil {
					microServicesMap["micro_service_name"] = microServices.MicroServiceName
				}

				microServicesList = append(microServicesList, microServicesMap)
			}

			apiInfoMap["micro_services"] = microServicesList
		}

		if apiAppApi.MicroServicesInfo != nil {
			apiInfoMap["micro_services_info"] = apiAppApi.MicroServicesInfo
		}

		if apiAppApi.ServiceTsfLoadBalanceConf != nil {
			serviceTsfLoadBalanceConfMap := map[string]interface{}{}

			if apiAppApi.ServiceTsfLoadBalanceConf.IsLoadBalance != nil {
				serviceTsfLoadBalanceConfMap["is_load_balance"] = apiAppApi.ServiceTsfLoadBalanceConf.IsLoadBalance
			}

			if apiAppApi.ServiceTsfLoadBalanceConf.Method != nil {
				serviceTsfLoadBalanceConfMap["method"] = apiAppApi.ServiceTsfLoadBalanceConf.Method
			}

			if apiAppApi.ServiceTsfLoadBalanceConf.SessionStickRequired != nil {
				serviceTsfLoadBalanceConfMap["session_stick_required"] = apiAppApi.ServiceTsfLoadBalanceConf.SessionStickRequired
			}

			if apiAppApi.ServiceTsfLoadBalanceConf.SessionStickTimeout != nil {
				serviceTsfLoadBalanceConfMap["session_stick_timeout"] = apiAppApi.ServiceTsfLoadBalanceConf.SessionStickTimeout
			}

			apiInfoMap["service_tsf_load_balance_conf"] = []interface{}{serviceTsfLoadBalanceConfMap}
		}

		if apiAppApi.ServiceTsfHealthCheckConf != nil {
			serviceTsfHealthCheckConfMap := map[string]interface{}{}

			if apiAppApi.ServiceTsfHealthCheckConf.IsHealthCheck != nil {
				serviceTsfHealthCheckConfMap["is_health_check"] = apiAppApi.ServiceTsfHealthCheckConf.IsHealthCheck
			}

			if apiAppApi.ServiceTsfHealthCheckConf.RequestVolumeThreshold != nil {
				serviceTsfHealthCheckConfMap["request_volume_threshold"] = apiAppApi.ServiceTsfHealthCheckConf.RequestVolumeThreshold
			}

			if apiAppApi.ServiceTsfHealthCheckConf.SleepWindowInMilliseconds != nil {
				serviceTsfHealthCheckConfMap["sleep_window_in_milliseconds"] = apiAppApi.ServiceTsfHealthCheckConf.SleepWindowInMilliseconds
			}

			if apiAppApi.ServiceTsfHealthCheckConf.ErrorThresholdPercentage != nil {
				serviceTsfHealthCheckConfMap["error_threshold_percentage"] = apiAppApi.ServiceTsfHealthCheckConf.ErrorThresholdPercentage
			}

			apiInfoMap["service_tsf_health_check_conf"] = []interface{}{serviceTsfHealthCheckConfMap}
		}

		if apiAppApi.EnableCORS != nil {
			apiInfoMap["enable_cors"] = apiAppApi.EnableCORS
		}

		if apiAppApi.Tags != nil {
			tagsList := []interface{}{}
			for _, tags := range apiAppApi.Tags {
				tagsMap := map[string]interface{}{}

				if tags.Key != nil {
					tagsMap["key"] = tags.Key
				}

				if tags.Value != nil {
					tagsMap["value"] = tags.Value
				}

				tagsList = append(tagsList, tagsMap)
			}

			apiInfoMap["tags"] = tagsList
		}

		if apiAppApi.Environments != nil {
			apiInfoMap["environments"] = apiAppApi.Environments
		}

		if apiAppApi.IsBase64Encoded != nil {
			apiInfoMap["is_base64_encoded"] = apiAppApi.IsBase64Encoded
		}

		if apiAppApi.IsBase64Trigger != nil {
			apiInfoMap["is_base64_trigger"] = apiAppApi.IsBase64Trigger
		}

		if apiAppApi.Base64EncodedTriggerRules != nil {
			base64EncodedTriggerRulesList := []interface{}{}
			for _, base64EncodedTriggerRules := range apiAppApi.Base64EncodedTriggerRules {
				base64EncodedTriggerRulesMap := map[string]interface{}{}

				if base64EncodedTriggerRules.Name != nil {
					base64EncodedTriggerRulesMap["name"] = base64EncodedTriggerRules.Name
				}

				if base64EncodedTriggerRules.Value != nil {
					base64EncodedTriggerRulesMap["value"] = base64EncodedTriggerRules.Value
				}

				base64EncodedTriggerRulesList = append(base64EncodedTriggerRulesList, base64EncodedTriggerRulesMap)
			}

			apiInfoMap["base64_encoded_trigger_rules"] = base64EncodedTriggerRulesList
		}

		tmpList = append(tmpList, apiInfoMap)
		_ = d.Set("result", tmpList)
	}

	d.SetId(strings.Join([]string{service_id, api_id, api_region}, tccommon.FILED_SP))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
