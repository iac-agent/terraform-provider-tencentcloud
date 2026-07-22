package apigateway

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"
)

func DataSourceTencentCloudAPIGatewayCustomerDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayCustomerDomainRead,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The 服务 ID",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			//Computed
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Service custom 域名 名称 list。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 名称",
						},
						"is_status_on": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "域名 名称 resolution 状态 有效值：`true`，`false`. `true` means normal parsing，`false` means parsing failed。",
						},
						"certificate_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The 证书 ID",
						},
						"is_default_mapping": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否use default 路径 mapping. 有效值：`true`，`false`. `true` means to use default 路径 mapping，`false` means to use custom 路径 mapping。",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom 域名 名称 agreement 类型",
						},
						"net_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network 类型",
						},
						"path_mappings": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "域名 名称 mapping 路径 and environment list。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The 域名 mapping 路径",
									},
									"environment": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Release environment。",
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

func dataSourceTencentCloudAPIGatewayCustomerDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_api_gateway_customer_domains.read")

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		serviceId         = d.Get("service_id").(string)
		infos             []*apigateway.DomainSetList
		list              []map[string]interface{}
		err               error
	)
	if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		infos, err = apiGatewayService.DescribeServiceSubDomains(ctx, serviceId)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	for _, info := range infos {
		var (
			pathMapping []map[string]interface{}
			status      bool
		)
		if !*info.IsDefaultMapping && *info.DomainName != "" {
			var mappings *apigateway.ServiceSubDomainMappings
			mappings, err = apiGatewayService.DescribeServiceSubDomainMappings(ctx, serviceId, *info.DomainName)
			if err != nil {
				return err
			}

			for _, v := range mappings.PathMappingSet {
				pathMapping = append(pathMapping, map[string]interface{}{
					"path":        v.Path,
					"environment": v.Environment,
				})
			}
		}
		if *info.Status == 1 {
			status = true
		}
		list = append(list, map[string]interface{}{
			"domain_name":        info.DomainName,
			"is_status_on":       status,
			"certificate_id":     info.CertificateId,
			"is_default_mapping": info.IsDefaultMapping,
			"protocol":           info.Protocol,
			"net_type":           info.NetType,
			"path_mappings":      pathMapping,
		})
	}

	if err = d.Set("list", list); err != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
		return err
	}

	d.SetId(serviceId)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}
	return nil
}
