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

func DataSourceTencentCloudAPIGatewayUsagePlanEnvironments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudUsagePlanEnvironmentRead,

		Schema: map[string]*schema.Schema{
			"usage_plan_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID usage plan 到 是 queried。",
			},
			"bind_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      API_GATEWAY_TYPE_SERVICE,
				ValidateFunc: tccommon.ValidateAllowedStringValue(API_GATEWAY_TYPES),
				Description:  "Binding 类型 有效值：`API`，`SERVICE`. 默认值：`SERVICE`。",
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
				Description: "A 列表 usage plan binding details。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "服务 ID",
						},
						"service_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "服务名称",
						},
						"api_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API ID，此 值 是 空 如果 attach 服务。",
						},
						"api_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API 名称，此 值 是 空 如果 attach 服务。",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API 路径，此 值 是 空 如果 attach 服务。",
						},
						"method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API 方法，此 值 是 空 如果 attach 服务。",
						},
						"environment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "环境 名称",
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
	}
}

func dataSourceTencentCloudUsagePlanEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_api_gateway_usage_plans.read")

	var (
		logId             = tccommon.GetLogId(tccommon.ContextNil)
		ctx               = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		usagePlanId       = d.Get("usage_plan_id").(string)
		bindType          = d.Get("bind_type").(string)
		infos             []*apigateway.UsagePlanEnvironment
		list              []map[string]interface{}
		err               error
	)

	if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		infos, err = apiGatewayService.DescribeUsagePlanEnvironments(ctx, usagePlanId, bindType)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	for _, info := range infos {
		list = append(list, map[string]interface{}{
			"service_id":   info.ServiceId,
			"service_name": info.ServiceName,
			"api_id":       info.ApiId,
			"api_name":     info.ApiName,
			"path":         info.Path,
			"method":       info.Method,
			"environment":  info.Environment,
			"modify_time":  info.ModifiedTime,
			"create_time":  info.CreatedTime,
		})
	}

	if err = d.Set("list", list); err != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
		return err
	}

	d.SetId(strings.Join([]string{usagePlanId, bindType}, tccommon.FILED_SP))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}
	return nil
}
