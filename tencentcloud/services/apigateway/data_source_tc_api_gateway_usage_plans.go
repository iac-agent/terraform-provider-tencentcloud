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

func DataSourceTencentCloudAPIGatewayUsagePlans() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayUsagePlansRead,

		Schema: map[string]*schema.Schema{
			"usage_plan_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID usage plan。",
			},
			"usage_plan_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 usage plan。",
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
				Description: "A 列表 usage plans。",
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
						"usage_plan_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom usage plan 描述",
						},
						"max_request_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total 数量 requests allowed. Valid 值 formats: `-1`，`[1,99999999]`. The 默认值为 -1，which 表示no 限制",
						},
						"max_request_num_pre_sec": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "限制 of requests per second. Valid values formats: `-1`，`[1,2000]`. The 默认值为 -1，which 表示no 限制",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "最后修改时间 in the 格式 of `YYYY-MM-DDThh:mm:ssZ` according to ISO 8601 standard. UTC time is used。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 in the 格式 of `YYYY-MM-DDThh:mm:ssZ` according to ISO 8601 standard. UTC time is used。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAPIGatewayUsagePlansRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_api_gateway_usage_plans.read")

	var (
		logId                      = tccommon.GetLogId(tccommon.ContextNil)
		ctx                        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		apiGatewayService          = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		infos                      []*apigateway.UsagePlanStatusInfo
		list                       []map[string]interface{}
		usagePlanId, usagePlanName string
		err                        error
	)

	if v, ok := d.GetOk("usage_plan_id"); ok {
		usagePlanId = v.(string)
	}
	if v, ok := d.GetOk("usage_plan_name"); ok {
		usagePlanName = v.(string)
	}

	if err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		infos, err = apiGatewayService.DescribeUsagePlansStatus(ctx, usagePlanId, usagePlanName)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	for _, info := range infos {
		var infoMap = make(map[string]interface{}, 7)
		infoMap["usage_plan_id"] = info.UsagePlanId
		infoMap["usage_plan_name"] = info.UsagePlanName
		infoMap["usage_plan_desc"] = info.UsagePlanDesc
		infoMap["max_request_num"] = info.MaxRequestNum
		infoMap["max_request_num_pre_sec"] = info.MaxRequestNumPreSec
		infoMap["modify_time"] = info.ModifiedTime
		infoMap["create_time"] = info.CreatedTime

		list = append(list, infoMap)
	}

	if err = d.Set("list", list); err != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s", logId, err.Error())
		return err
	}

	d.SetId(strings.Join([]string{usagePlanId, usagePlanName}, tccommon.FILED_SP))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}
	return nil
}
