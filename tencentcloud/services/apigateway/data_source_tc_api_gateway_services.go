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

func DataSourceTencentCloudAPIGatewayServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayServicesRead,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "服务名称 对于 查询。",
			},
			"service_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "服务 ID 对于 查询。",
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
				Description: "A 列表 services。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom 服务 ID",
						},
						"service_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom 服务名称",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service frontend 请求 类型 有效值：`http`，`https`，`http&https`。",
						},
						"service_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom 服务 描述",
						},
						"exclusive_set_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Deprecated:  "It has been deprecated from version 1.81.9.",
							Description: "Self-deployed 集群名称，其中 是 用于指定self-deployed 集群 其中 服务 是 到 是 创建。",
						},
						"net_type": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Description: "网络类型列表，用于指定支持的网络类型。" +
								"Valid values: `INNER`, `OUTER`. " +
								"`INNER` indicates access over private network, and `OUTER` indicates access over public network.",
						},
						"ip_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP 版本 数量。",
						},
						"internal_sub_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Private 网络 访问 sub-域名 名称",
						},
						"outer_sub_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public 网络 访问 subdomain 名称",
						},
						"inner_http_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "端口 数量 对于 http 访问 over 私有 网络。",
						},
						"inner_https_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "端口 数量 对于 https 访问 over 私有 网络。",
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
						"usage_plan_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A 列表 attach usage plans. Each element 包含following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"usage_plan_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID usage plan。",
									},
									"usage_plan_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "名称 usage plan。",
									},
									"bind_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Binding 类型",
									},
									"api_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID API。",
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

func dataSourceTencentCloudAPIGatewayServicesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_api_gateway_services.read")()

	var (
		logId                  = tccommon.GetLogId(tccommon.ContextNil)
		ctx                    = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService      = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		services               []*apigateway.Service
		serviceName, serviceId string
		has                    bool
		err                    error
	)

	if v, ok := d.GetOk("service_name"); ok {
		serviceName = v.(string)
	}
	if v, ok := d.GetOk("service_id"); ok {
		serviceId = v.(string)
	}

	if outErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		services, err = apiGatewayService.DescribeServicesStatus(ctx, serviceId, serviceName)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); outErr != nil {
		return outErr
	}

	list := make([]map[string]interface{}, 0, len(services))

	for _, service := range services {
		var info apigateway.DescribeServiceResponse
		if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			info, has, err = apiGatewayService.DescribeService(ctx, *service.ServiceId)
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

		var plans []*apigateway.ApiUsagePlan

		var planList = make([]map[string]interface{}, 0, len(info.Response.ApiIdStatusSet))
		var hasContains = make(map[string]bool, len(info.Response.ApiIdStatusSet))

		//from service
		if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			plans, err = apiGatewayService.DescribeServiceUsagePlan(ctx, *service.ServiceId)
			if err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			return nil
		}); err != nil {
			return err
		}

		for _, item := range plans {
			if hasContains[*item.UsagePlanId] {
				continue
			}
			hasContains[*item.UsagePlanId] = true
			planList = append(
				planList, map[string]interface{}{
					"usage_plan_id":   item.UsagePlanId,
					"usage_plan_name": item.UsagePlanName,
					"bind_type":       API_GATEWAY_TYPE_SERVICE,
					"api_id":          "",
				})
		}

		//from api
		if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			plans, err = apiGatewayService.DescribeApiUsagePlan(ctx, *service.ServiceId)
			if err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
		for _, item := range plans {
			planList = append(
				planList, map[string]interface{}{
					"usage_plan_id":   item.UsagePlanId,
					"usage_plan_name": item.UsagePlanName,
					"bind_type":       API_GATEWAY_TYPE_API,
					"api_id":          item.ApiId,
				})
		}

		list = append(list, map[string]interface{}{
			"service_id":   info.Response.ServiceId,
			"service_name": info.Response.ServiceName,
			"protocol":     info.Response.Protocol,
			"service_desc": info.Response.ServiceDesc,
			//"exclusive_set_name":  info.Response.ExclusiveSetName,
			"ip_version":          info.Response.IpVersion,
			"net_type":            info.Response.NetTypes,
			"internal_sub_domain": info.Response.InternalSubDomain,
			"outer_sub_domain":    info.Response.OuterSubDomain,
			"inner_http_port":     info.Response.InnerHttpPort,
			"inner_https_port":    info.Response.InnerHttpsPort,
			"modify_time":         info.Response.ModifiedTime,
			"create_time":         info.Response.CreatedTime,
			"usage_plan_list":     planList,
		})
	}

	if err = d.Set("list", list); err != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
		return err
	}

	d.SetId(strings.Join([]string{serviceName, serviceId}, tccommon.FILED_SP))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}
	return nil
}
