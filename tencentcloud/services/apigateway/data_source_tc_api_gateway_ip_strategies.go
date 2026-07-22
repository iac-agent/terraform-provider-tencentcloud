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

func DataSourceTencentCloudAPIGatewayIpStrategy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayIpStrategyRead,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "服务 ID 到 是 queried。",
			},
			"strategy_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 IP 策略。",
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
				Description: "A 列表 strategy。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"strategy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "strategy ID。",
						},
						"strategy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 strategy。",
						},
						"strategy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 strategy。",
						},
						"ip_list": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "列表 IP。",
						},
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "服务 ID",
						},
						"bind_api_total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 API bound 到 strategy。",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "最后修改时间 在 格式 的 YYYY-MM-DDThh:mm:ssZ according 到 ISO 8601 standard. UTC 时间 是 使用。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 在 格式 的 YYYY-MM-DDThh:mm:ssZ according 到 ISO 8601 standard. UTC 时间 是 使用。",
						},
						"attach_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "列表 bound API details。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "服务 ID",
									},
									"api_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API ID。",
									},
									"api_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API interface 描述",
									},
									"api_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API 名称",
									},
									"vpc_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "私有网络 ID",
									},
									"uniq_vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC 唯一 ID。",
									},
									"api_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API 类型 有效值：`NORMAL`，`TSF`. `NORMAL` 表示 common API，`TSF` 表示 microservice API。",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API 协议",
									},
									"auth_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API authentication 类型 有效值：`SECRET`，`NONE`，`OAUTH`. `SECRET` 表示 键 pair authentication，`NONE` 表示 无 authentication。",
									},
									"api_business_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "类型 oauth API. 此 字段 是 有效 当 `auth_type` 是 `OAUTH`，和 值 是 `NORMAL` (business API) 和 `OAUTH` (authorization API)。",
									},
									"auth_relation_api_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "唯一 ID associated authorization API，其中 takes effect 当 authType 是 `OAUTH` 和 `ApiBusinessType` 是 normal. Identifies 唯一 ID oauth2.0 authorization API bound 到 business API。",
									},
									"tags": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Computed:    true,
										Description: "标签 信息 associated 使用 API。",
									},
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API 路径",
									},
									"method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "API 请求 方法。",
									},
									"relation_business_api_ids": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Computed:    true,
										Description: "列表 business API associated 使用 authorized API。",
									},
									"oauth_config": {
										Type:        schema.TypeMap,
										Computed:    true,
										Description: "OAUTH 配置 信息. It takes effect 当 authType 是 `OAUTH`。",
									},
									"modify_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "最后修改时间 在 格式 的 `YYYY-MM-DDThh:mm:ssZ` according 到 ISO 8601 standard. UTC 时间 是 使用。",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创建时间 在 格式 的 `YYYY-MM-DDThh:mm:ssZ` according 到 ISO 8601 standard. UTC 时间 是 使用。",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAPIGatewayIpStrategyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_api_gateway_ip_strategy.read")

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		serviceId         = d.Get("service_id").(string)
		infos             []*apigateway.IPStrategy
		list              []map[string]interface{}
		strategyName      string
		err               error
	)
	if v, ok := d.GetOk("strategy_name"); ok {
		strategyName = v.(string)
	}

	if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		infos, err = apiGatewayService.DescribeIPStrategysStatus(ctx, serviceId, strategyName)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	for _, info := range infos {
		var attachListInfo []map[string]interface{}

		for _, env := range API_GATEWAY_SERVICE_ENVS {
			var strategy *apigateway.IPStrategy
			if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				strategy, err = apiGatewayService.DescribeIPStrategies(ctx, serviceId, *info.StrategyId, env)
				if err != nil {
					return tccommon.RetryError(err, tccommon.InternalError)
				}
				return nil
			}); err != nil {
				return err
			}

			for _, api := range strategy.BindApis {
				attachListInfo = append(attachListInfo, map[string]interface{}{
					"service_id":                api.ServiceId,
					"api_id":                    api.ApiId,
					"api_desc":                  api.ApiDesc,
					"api_name":                  api.ApiName,
					"vpc_id":                    api.VpcId,
					"uniq_vpc_id":               api.UniqVpcId,
					"api_type":                  api.ApiType,
					"protocol":                  api.Protocol,
					"auth_type":                 api.AuthType,
					"api_business_type":         api.ApiBusinessType,
					"auth_relation_api_id":      api.AuthRelationApiId,
					"tags":                      api.Tags,
					"path":                      api.Path,
					"method":                    api.Method,
					"relation_business_api_ids": api.RelationBuniessApiIds,
					"oauth_config":              flattenOauthConfigMappings(api.OauthConfig),
					"modify_time":               api.ModifiedTime,
					"create_time":               api.CreatedTime,
				})
			}
		}

		infoMap := map[string]interface{}{
			"strategy_id":          info.StrategyId,
			"strategy_name":        info.StrategyName,
			"strategy_type":        info.StrategyType,
			"ip_list":              info.StrategyData,
			"service_id":           info.ServiceId,
			"bind_api_total_count": info.BindApiTotalCount,
			"modify_time":          info.ModifiedTime,
			"create_time":          info.CreatedTime,
			"attach_list":          attachListInfo,
		}

		list = append(list, infoMap)
	}

	if err = d.Set("list", list); err != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
		return err
	}

	d.SetId(strings.Join([]string{serviceId, strategyName}, tccommon.FILED_SP))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}
	return nil
}
