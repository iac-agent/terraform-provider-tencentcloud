package apigateway

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiGateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudApiGatewayImportOpenApi() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudApiGatewayImportOpenApiCreate,
		Read:   resourceTencentCloudApiGatewayImportOpenApiRead,
		Delete: resourceTencentCloudApiGatewayImportOpenApiDelete,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "唯一 ID 服务 其中 API 是 located。",
			},
			"content": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "OpenAPI 正文 内容",
			},
			"encode_type": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				Default:      IMPORT_OPEN_API_ENCODE_TYPE_YAML,
				ValidateFunc: tccommon.ValidateAllowedStringValue(IMPORT_OPEN_API_ENCODE_TYPE),
				Description:  "内容 格式 可以 仅 是 YAML 或 JSON，和 默认为 YAML。",
			},
			"content_version": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Default:     "OpenAPI",
				Description: "内容 版本 默认为 OpenAPI 和 currently 仅 支持 OpenAPI。",
			},
			// Computed
			"api_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Custom Api ID。",
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
			"api_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "API 类型，支持 NORMAL (regular API) 和 TSF (microservice API)，默认为 NORMAL。",
			},
			"auth_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API authentication 类型 Support SECRET (键 Pair Authentication)，NONE (Authentication Exemption)，OAUTH，APP (Application Authentication). 默认为 NONE。",
			},
			"protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API frontend 请求 类型 有效值：`HTTP`，`WEBSOCKET`. 默认值：`HTTP`。",
			},
			"enable_cors": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "是否enable CORS. 默认值：`true`。",
			},
			"request_config_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Request frontend 路径 配置. Like `/用户/getinfo`。",
			},
			"request_config_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Request frontend 方法 配置. 有效值：`GET`,`POST`,`PUT`,`DELETE`,`HEAD`,`ANY`. 默认值：`GET`。",
			},
			"constant_parameters": {
				Computed:    true,
				Type:        schema.TypeSet,
				Description: "Constant 参数。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Constant 参数 名称 此 参数 是 仅 使用 当 ServiceType 是 HTTP.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Constant 参数 描述 此 参数 是 仅 使用 当 ServiceType 是 HTTP.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"position": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Constant 参数 position. 此 参数 是 仅 使用 当 ServiceType 是 HTTP.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "默认值 对于 constant 参数. 此 参数 是 仅 使用 当 ServiceType 是 HTTP.注意：此字段可能返回 null，表示无法获取有效值。",
						},
					},
				},
			},
			"request_parameters": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Frontend 请求 参数。",
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
							Description: "如果 此 参数 必填 默认值：`false`。",
						},
					},
				},
			},
			"micro_services": {
				Computed:    true,
				Type:        schema.TypeSet,
				Description: "API bound microservice 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Micro 服务 集群。",
						},
						"namespace_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Microservice 命名空间。",
						},
						"micro_service_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Microservice 名称",
						},
					},
				},
			},
			"service_tsf_load_balance_conf": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Load balancing 配置 对于 microservices。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_load_balance": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is load balancing 已启用注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Load balancing 方法.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"session_stick_required": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否enable 会话 persistence.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"session_stick_timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Session hold 超时.注意：此字段可能返回 null，表示无法获取有效值。",
						},
					},
				},
			},
			"service_tsf_health_check_conf": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Health check 配置 对于 microservices。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_health_check": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否initiate health check.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"request_volume_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Health check 阈值.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"sleep_window_in_milliseconds": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Window 大小.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"error_threshold_percentage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Threshold percentage.注意：此字段可能返回 null，表示无法获取有效值。",
						},
					},
				},
			},
			"api_business_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "当 `auth_type` 是 OAUTH，此 字段 是 有效，NORMAL: Business API，OAUTH: Authorization API。",
			},
			"service_config_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "backend 服务 类型 API. Supports HTTP，MOCK，TSF，SCF，WEBSOCKET，COS，TARGET (内部 testing)。",
			},
			"service_config_timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "API backend 服务 超时 周期 （秒）。 默认值：`5`。",
			},
			"service_config_product": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Backend 类型 Effective 当 enabling vpc，currently 支持 types 是 clb，cvm，和 upstream。",
			},
			"service_config_vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique 私有网络 ID",
			},
			"service_config_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "backend 服务 URL 的 API. 如果 ServiceType 是 HTTP，此 参数 必须 是 passed。",
			},
			"service_config_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API backend 服务 路径，such 作为 /路径 如果 `service_config_type` 是 `HTTP`，此 参数 将 是 必填 frontend `request_config_path` 和 backend 路径 `service_config_path` 可以 是 different。",
			},
			"service_config_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API backend 服务 请求 方法，such 作为 `GET`. 如果 `service_config_type` 是 `HTTP`，此 参数 将 是 必填 frontend `request_config_method` 和 backend 方法 `service_config_method` 可以 是 different。",
			},
			"service_config_upstream_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Only 必填 当 binding 到 VPC channels注意：此字段可能返回 null，表示无法获取有效值。",
			},
			"service_config_cos_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "API backend COS 配置. 如果 ServiceType 是 COS，then 此 参数 必须 是 passed.注意：此字段可能返回 null，表示无法获取有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API calls backend COS 方法，和 可选 值 对于 front-end 请求 方法 和 操作 是:GET: GetObjectPUT: PutObjectPOST: PostObject，AppendObjectHEAD: HeadObjectDELETE: DeleteObject.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"bucket_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "存储桶 名称 API backend COS.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"authorization": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "API calls 签名 switch 的 backend COS，其中 默认为 false.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"path_match_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "路径 matching 模式 对于 API backend COS，可选 值:BackEndPath: Backend 路径 matchingFullPath: Full 路径 MatchingThe 默认值 是: BackEndPath注意：此字段可能返回 null，表示无法获取有效值。",
						},
					},
				},
			},
			"service_config_scf_function_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCF 函数 名称 此 参数 takes effect 当 `service_config_type` 是 `SCF`。",
			},
			"service_config_scf_function_namespace": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCF 函数 命名空间. 此 参数 takes effect 当 `service_config_type` 是 `SCF`。",
			},
			"service_config_scf_function_qualifier": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SCF 函数 版本 此 参数 takes effect 当 `service_config_type` 是 `SCF`。",
			},
			"service_config_scf_function_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scf 函数 类型 Effective 当 backend 类型 是 SCF. Support Event Triggering (EVENT) 和 HTTP Direct Cloud Function (HTTP)。",
			},
			"service_config_scf_is_integrated_response": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "是否enable response integration. Effective 当 backend 类型 是 SCF。",
			},
			"service_config_mock_return_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Returned 信息 的 API backend mocking. 此 参数 为必填项 当 `service_config_type` 是 `MOCK`。",
			},
			"service_config_websocket_register_function_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket registration 函数. It takes effect 当 当前 end 类型 是 WEBSOCKET 和 backend 类型 是 SCF。",
			},
			"service_config_websocket_cleanup_function_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket cleaning 函数. It takes effect 当 当前 end 类型 是 WEBSOCKET 和 backend 类型 是 SCF。",
			},
			"service_config_websocket_transport_function_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket transfer 函数. It takes effect 当 当前 end 类型 是 WEBSOCKET 和 backend 类型 是 SCF。",
			},
			"service_config_websocket_register_function_namespace": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket registers 函数 namespaces. It takes effect 当 当前 end 类型 是 WEBSOCKET 和 backend 类型 是 SCF。",
			},
			"service_config_websocket_register_function_qualifier": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket transfer 函数 版本 It takes effect 当 当前 end 类型 是 WEBSOCKET 和 backend 类型 是 SCF。",
			},
			"service_config_websocket_transport_function_namespace": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket transfer 函数 命名空间. It takes effect 当 当前 end 类型 是 WEBSOCKET 和 backend 类型 是 SCF。",
			},
			"service_config_websocket_transport_function_qualifier": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket transfer 函数 版本 It takes effect 当 当前 end 类型 是 WEBSOCKET 和 backend 类型 是 SCF。",
			},
			"service_config_websocket_cleanup_function_namespace": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket cleans up 函数 命名空间. It takes effect 当 当前 end 类型 是 WEBSOCKET 和 backend 类型 是 SCF。",
			},
			"service_config_websocket_cleanup_function_qualifier": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scf websocket cleaning 函数 版本 It takes effect 当 当前 end 类型 是 WEBSOCKET 和 backend 类型 是 SCF。",
			},
			"is_debug_after_charge": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Charge after starting debugging. (Cloud Market Reserved Fields)。",
			},
			"is_delete_response_error_codes": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Do 您 want 到 delete 自定义 response 配置 错误码? 如果 它 是 不 passed 或 False 是 passed，它 将 不 是 删除. 如果 True 是 passed，all 自定义 response 配置 错误 codes 对于 此 API 将 是 删除。",
			},
			"response_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Return 类型 有效值：`HTML`，`JSON`，`TEXT`，`BINARY`，`XML`. 默认值：`HTML`。",
			},
			"response_success_example": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Successful response sample 的 自定义 response 配置。",
			},
			"response_fail_example": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Response failure sample 的 自定义 response 配置。",
			},
			"response_error_codes": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Custom 错误码 配置. Must keep 在 least 一个 after 集合。",
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
			"auth_relation_api_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "唯一 ID associated authorization API takes effect 当 AuthType 是 OAUTH 和 ApiBusinessType 是 NORMAL. 唯一 ID oauth2.0 authorized API 该 identifies business API binding。",
			},
			"service_parameters": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "backend 服务 参数 的 API。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "backend 服务 参数 名称 API. 此 参数 是 仅 使用 当 ServiceType 是 HTTP. front 和 rear 参数 names 可以 是 different.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"position": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "backend 服务 参数 location 的 API，such 作为 head. 此 参数 是 仅 使用 当 ServiceType 是 HTTP. 参数 positions 在 front 和 rear 结束 可以 是 已配置 differently.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"relevant_request_parameter_position": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "location 的 front-end 参数 corresponding 到 backend 服务 参数 的 API，such 作为 head. 此 参数 是 仅 使用 当 ServiceType 是 HTTP.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"relevant_request_parameter_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 front-end 参数 corresponding 到 backend 服务 参数 的 API. 此 参数 是 仅 使用 当 ServiceType 是 HTTP.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "默认值 对于 backend 服务 参数 的 API. 此 参数 是 仅 使用 当 ServiceType 是 HTTP.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"relevant_request_parameter_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "备注 在 backend 服务 参数 的 API. 此 参数 是 仅 使用 当 ServiceType 是 HTTP.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"relevant_request_parameter_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "backend 服务 参数 类型 API. 此 参数 是 仅 使用 当 ServiceType 是 HTTP.注意：此字段可能返回 null，表示无法获取有效值。",
						},
					},
				},
			},
			"oauth_config": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "OAuth 配置. Effective 当 AuthType 是 OAUTH。",
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
							Description: "令牌 passes position。",
						},
						"login_redirect_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Redirect 地址，用于guide users 在 login operations。",
						},
					},
				},
			},
			"is_base64_encoded": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "是否enable Base64 编码 将 仅 take effect 当 backend 是 scf。",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "最后修改时间 在 格式 的 YYYY-MM-DDThh:mm:ssZ according 到 ISO 8601 standard. UTC 时间 是 使用。",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "创建时间 在 格式 的 YYYY-MM-DDThh:mm:ssZ according 到 ISO 8601 standard. UTC 时间 是 使用。",
			},
		},
	}
}

func resourceTencentCloudApiGatewayImportOpenApiCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_import_open_api.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		request   = apiGateway.NewImportOpenApiRequest()
		response  = apiGateway.NewImportOpenApiResponse()
		serviceId string
		apiId     string
	)

	if v, ok := d.GetOk("service_id"); ok {
		request.ServiceId = helper.String(v.(string))
		serviceId = v.(string)
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	if v, ok := d.GetOk("encode_type"); ok {
		request.EncodeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("content_version"); ok {
		request.ContentVersion = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAPIGatewayClient().ImportOpenApi(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("apiGateway importOpenApi not exists")
			return resource.NonRetryableError(e)
		}

		if *result.Response.Result.ApiSet[0].Status == "success" {
			response = result
			return nil
		}

		return resource.RetryableError(fmt.Errorf("create apiGateway importOpenApi is running, status: %s", *result.Response.Result.ApiSet[0].Status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s create apiGateway importOpenApi failed, reason:%+v", logId, err)
		return err
	}

	apiId = *response.Response.Result.ApiSet[0].ApiId
	d.SetId(strings.Join([]string{serviceId, apiId}, tccommon.FILED_SP))
	return resourceTencentCloudApiGatewayImportOpenApiRead(d, meta)
}

func resourceTencentCloudApiGatewayImportOpenApiRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_import_open_api.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	serviceId := idSplit[0]
	apiId := idSplit[1]

	info, err := service.DescribeApiGatewayImportOpenApiById(ctx, serviceId, apiId)
	if err != nil {
		return err
	}

	if info == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ApiGatewayImportOpenApi` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("service_id", info.ServiceId)
	_ = d.Set("api_id", info.ApiId)
	_ = d.Set("api_name", info.ApiName)
	_ = d.Set("api_desc", info.ApiDesc)
	_ = d.Set("api_type", info.ApiType)
	_ = d.Set("auth_type", info.AuthType)
	_ = d.Set("protocol", info.Protocol)
	_ = d.Set("request_config_path", info.RequestConfig.Path)
	_ = d.Set("request_config_method", info.RequestConfig.Method)
	_ = d.Set("enable_cors", info.EnableCORS)

	if info.ConstantParameters != nil {
		constantParametersList := []interface{}{}
		for _, constantParameters := range info.ConstantParameters {
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

		_ = d.Set("constant_parameters", constantParametersList)
	}

	if info.RequestParameters != nil {
		list := make([]map[string]interface{}, 0, len(info.RequestParameters))
		for _, param := range info.RequestParameters {
			list = append(list, map[string]interface{}{
				"name":          param.Name,
				"position":      param.Position,
				"type":          param.Type,
				"desc":          param.Desc,
				"default_value": param.DefaultValue,
				"required":      param.Required,
			})
		}
		_ = d.Set("request_parameters", list)
	}

	if info.MicroServices != nil {
		microServicesList := []interface{}{}
		for _, microServices := range info.MicroServices {
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

		_ = d.Set("micro_services", microServicesList)

	}

	if info.ServiceTsfLoadBalanceConf != nil {
		ServiceTsfLoadBalanceConfMap := map[string]interface{}{}

		if info.ServiceTsfLoadBalanceConf.IsLoadBalance != nil {
			ServiceTsfLoadBalanceConfMap["is_load_balance"] = info.ServiceTsfLoadBalanceConf.IsLoadBalance
		}

		if info.ServiceTsfLoadBalanceConf.Method != nil {
			ServiceTsfLoadBalanceConfMap["method"] = info.ServiceTsfLoadBalanceConf.Method
		}

		if info.ServiceTsfLoadBalanceConf.SessionStickRequired != nil {
			ServiceTsfLoadBalanceConfMap["session_stick_required"] = info.ServiceTsfLoadBalanceConf.SessionStickRequired
		}

		if info.ServiceTsfLoadBalanceConf.SessionStickTimeout != nil {
			ServiceTsfLoadBalanceConfMap["session_stick_timeout"] = info.ServiceTsfLoadBalanceConf.SessionStickTimeout
		}

		_ = d.Set("service_tsf_load_balance_conf", []interface{}{ServiceTsfLoadBalanceConfMap})
	} else {
		_ = d.Set("service_tsf_load_balance_conf", []interface{}{})
	}

	if info.ServiceTsfHealthCheckConf != nil {
		serviceTsfHealthCheckConfMap := map[string]interface{}{}

		if info.ServiceTsfHealthCheckConf.IsHealthCheck != nil {
			serviceTsfHealthCheckConfMap["is_health_check"] = info.ServiceTsfHealthCheckConf.IsHealthCheck
		}

		if info.ServiceTsfHealthCheckConf.RequestVolumeThreshold != nil {
			serviceTsfHealthCheckConfMap["request_volume_threshold"] = info.ServiceTsfHealthCheckConf.RequestVolumeThreshold
		}

		if info.ServiceTsfHealthCheckConf.SleepWindowInMilliseconds != nil {
			serviceTsfHealthCheckConfMap["sleep_window_in_milliseconds"] = info.ServiceTsfHealthCheckConf.SleepWindowInMilliseconds
		}

		if info.ServiceTsfHealthCheckConf.ErrorThresholdPercentage != nil {
			serviceTsfHealthCheckConfMap["error_threshold_percentage"] = info.ServiceTsfHealthCheckConf.ErrorThresholdPercentage
		}

		_ = d.Set("service_tsf_health_check_conf", []interface{}{serviceTsfHealthCheckConfMap})
	} else {
		_ = d.Set("service_tsf_health_check_conf", []interface{}{})
	}

	if info.ApiBusinessType != nil {
		_ = d.Set("api_business_type", info.ApiBusinessType)
	}

	_ = d.Set("service_config_type", info.ServiceType)
	_ = d.Set("service_config_timeout", info.ServiceTimeout)
	if info.ServiceConfig != nil {
		if info.ServiceConfig.Product != nil {
			_ = d.Set("service_config_product", info.ServiceConfig.Product)
		}

		if info.ServiceConfig.UniqVpcId != nil {
			_ = d.Set("service_config_vpc_id", info.ServiceConfig.UniqVpcId)
		}

		if info.ServiceConfig.Url != nil {
			_ = d.Set("service_config_url", info.ServiceConfig.Url)
		}

		if info.ServiceConfig.Path != nil {
			_ = d.Set("service_config_path", info.ServiceConfig.Path)
		}

		if info.ServiceConfig.Method != nil {
			_ = d.Set("service_config_method", info.ServiceConfig.Method)
		}

		if info.ServiceConfig.UpstreamId != nil {
			_ = d.Set("service_config_upstream_id", info.ServiceConfig.UpstreamId)
		}

		if info.ServiceConfig.CosConfig != nil {
			cosConfigMap := map[string]interface{}{}

			if info.ServiceConfig.CosConfig.Action != nil {
				cosConfigMap["action"] = info.ServiceConfig.CosConfig.Action
			}

			if info.ServiceConfig.CosConfig.BucketName != nil {
				cosConfigMap["bucket_name"] = info.ServiceConfig.CosConfig.BucketName
			}

			if info.ServiceConfig.CosConfig.Authorization != nil {
				cosConfigMap["authorization"] = info.ServiceConfig.CosConfig.Authorization
			}

			if info.ServiceConfig.CosConfig.PathMatchMode != nil {
				cosConfigMap["path_match_mode"] = info.ServiceConfig.CosConfig.PathMatchMode
			}

			_ = d.Set("service_config_cos_config", []interface{}{cosConfigMap})
		} else {
			_ = d.Set("service_config_cos_config", []interface{}{})
		}
	}

	_ = d.Set("service_config_scf_function_name", info.ServiceScfFunctionName)
	_ = d.Set("service_config_scf_function_namespace", info.ServiceScfFunctionNamespace)
	_ = d.Set("service_config_scf_function_qualifier", info.ServiceScfFunctionQualifier)
	_ = d.Set("service_config_scf_is_integrated_response", info.ServiceScfIsIntegratedResponse)
	_ = d.Set("service_config_mock_return_message", info.ServiceMockReturnMessage)

	_ = d.Set("service_config_websocket_register_function_name", info.ServiceWebsocketRegisterFunctionName)
	_ = d.Set("service_config_websocket_cleanup_function_name", info.ServiceWebsocketCleanupFunctionName)
	_ = d.Set("service_config_websocket_transport_function_name", info.ServiceWebsocketTransportFunctionName)
	_ = d.Set("service_config_websocket_register_function_namespace", info.ServiceWebsocketRegisterFunctionNamespace)
	_ = d.Set("service_config_websocket_register_function_qualifier", info.ServiceWebsocketRegisterFunctionQualifier)
	_ = d.Set("service_config_websocket_transport_function_namespace", info.ServiceWebsocketTransportFunctionNamespace)
	_ = d.Set("service_config_websocket_transport_function_qualifier", info.ServiceWebsocketTransportFunctionQualifier)
	_ = d.Set("service_config_websocket_cleanup_function_namespace", info.ServiceWebsocketCleanupFunctionNamespace)
	_ = d.Set("service_config_websocket_cleanup_function_qualifier", info.ServiceWebsocketCleanupFunctionQualifier)

	_ = d.Set("is_debug_after_charge", info.IsDebugAfterCharge)
	_ = d.Set("response_type", info.ResponseType)
	_ = d.Set("response_success_example", info.ResponseSuccessExample)
	_ = d.Set("response_fail_example", info.ResponseFailExample)
	_ = d.Set("auth_relation_api_id", info.AuthRelationApiId)
	_ = d.Set("update_time", info.ModifiedTime)
	_ = d.Set("create_time", info.CreatedTime)

	if info.ServiceParameters != nil {
		serviceParametersList := []interface{}{}
		for _, serviceParameters := range info.ServiceParameters {
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

		_ = d.Set("service_parameters", serviceParametersList)

	}

	if info.OauthConfig != nil {
		oauthConfigMap := map[string]interface{}{}

		if info.OauthConfig.PublicKey != nil {
			oauthConfigMap["public_key"] = info.OauthConfig.PublicKey
		}

		if info.OauthConfig.TokenLocation != nil {
			oauthConfigMap["token_location"] = info.OauthConfig.TokenLocation
		}

		if info.OauthConfig.LoginRedirectUrl != nil {
			oauthConfigMap["login_redirect_url"] = info.OauthConfig.LoginRedirectUrl
		}

		_ = d.Set("oauth_config", []interface{}{oauthConfigMap})
	} else {
		_ = d.Set("oauth_config", []interface{}{})
	}

	if info.ResponseErrorCodes != nil {
		list := make([]map[string]interface{}, 0, len(info.ResponseErrorCodes))
		for _, code := range info.ResponseErrorCodes {
			list = append(list, map[string]interface{}{
				"code":           code.Code,
				"msg":            code.Msg,
				"desc":           code.Desc,
				"converted_code": code.ConvertedCode,
				"need_convert":   code.NeedConvert,
			})
		}
		_ = d.Set("response_error_codes", list)
	}

	_ = d.Set("is_base64_encoded", info.IsBase64Encoded)

	return nil
}

func resourceTencentCloudApiGatewayImportOpenApiDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_import_open_api.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	serviceId := idSplit[0]
	apiId := idSplit[1]

	if err := service.DeleteApiGatewayImportOpenApiById(ctx, serviceId, apiId); err != nil {
		return err
	}

	return nil
}
