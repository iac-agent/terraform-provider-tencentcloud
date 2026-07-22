package apigateway

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudAPIGatewayApiUsagePlans() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAPIGatewayApiUsagePlanRead,
		Schema: map[string]*schema.Schema{
			"service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The unique ID service to be queried。",
			},
			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "API binding usage plan list.注意：此字段可能返回 null，表示无法获取有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service unique ID.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"service_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "服务名称注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"api_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API unique ID.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"api_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API 名称注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API 路径注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API method.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"usage_plan_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Use the unique ID plan.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"usage_plan_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Use the 名称 plan.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"usage_plan_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 usage plan.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"environment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Use the service environment bound by the plan.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"in_use_request_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The quota that has already been used.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"max_request_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Request total quota，-1 表示no 限制注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"max_request_num_pre_sec": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Request QPS upper 限制，-1 表示no 限制注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create a time using a schedule.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"modified_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Use the last 修改时间 of the plan.注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudAPIGatewayApiUsagePlanRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_api_gateway_api_usage_plans.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		result  []*apigateway.ApiUsagePlan
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("service_id"); ok {
		paramMap["ServiceId"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeAPIGatewayApiUsagePlanByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		result = response
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	if result != nil {
		apiUsagePlanListList := []interface{}{}
		for _, apiUsagePlanList := range result {
			apiUsagePlanListMap := map[string]interface{}{}

			if apiUsagePlanList.ServiceId != nil {
				apiUsagePlanListMap["service_id"] = apiUsagePlanList.ServiceId
			}

			if apiUsagePlanList.ServiceName != nil {
				apiUsagePlanListMap["service_name"] = apiUsagePlanList.ServiceName
			}

			if apiUsagePlanList.ApiId != nil {
				apiUsagePlanListMap["api_id"] = apiUsagePlanList.ApiId
			}

			if apiUsagePlanList.ApiName != nil {
				apiUsagePlanListMap["api_name"] = apiUsagePlanList.ApiName
			}

			if apiUsagePlanList.Path != nil {
				apiUsagePlanListMap["path"] = apiUsagePlanList.Path
			}

			if apiUsagePlanList.Method != nil {
				apiUsagePlanListMap["method"] = apiUsagePlanList.Method
			}

			if apiUsagePlanList.UsagePlanId != nil {
				apiUsagePlanListMap["usage_plan_id"] = apiUsagePlanList.UsagePlanId
			}

			if apiUsagePlanList.UsagePlanName != nil {
				apiUsagePlanListMap["usage_plan_name"] = apiUsagePlanList.UsagePlanName
			}

			if apiUsagePlanList.UsagePlanDesc != nil {
				apiUsagePlanListMap["usage_plan_desc"] = apiUsagePlanList.UsagePlanDesc
			}

			if apiUsagePlanList.Environment != nil {
				apiUsagePlanListMap["environment"] = apiUsagePlanList.Environment
			}

			if apiUsagePlanList.InUseRequestNum != nil {
				apiUsagePlanListMap["in_use_request_num"] = apiUsagePlanList.InUseRequestNum
			}

			if apiUsagePlanList.MaxRequestNum != nil {
				apiUsagePlanListMap["max_request_num"] = apiUsagePlanList.MaxRequestNum
			}

			if apiUsagePlanList.MaxRequestNumPreSec != nil {
				apiUsagePlanListMap["max_request_num_pre_sec"] = apiUsagePlanList.MaxRequestNumPreSec
			}

			if apiUsagePlanList.CreatedTime != nil {
				apiUsagePlanListMap["created_time"] = apiUsagePlanList.CreatedTime
			}

			if apiUsagePlanList.ModifiedTime != nil {
				apiUsagePlanListMap["modified_time"] = apiUsagePlanList.ModifiedTime
			}

			ids = append(ids, *apiUsagePlanList.UsagePlanId)
			apiUsagePlanListList = append(apiUsagePlanListList, apiUsagePlanListMap)
		}

		_ = d.Set("result", apiUsagePlanListList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
