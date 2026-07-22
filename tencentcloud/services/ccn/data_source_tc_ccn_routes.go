package ccn

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCcnRoutes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCcnRoutesRead,

		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID CCN 到 是 queried。",
			},
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 conditions。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "待过滤字段 Support `路由-ID`，`cidr-block`，`实例-类型`，`实例-地域`，`实例-ID`，`路由-表-ID`。",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "过滤值 的 字段。",
						},
					},
				},
			},
			// Computed
			"route_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "CCN 路由 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"route_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "路由 ID。",
						},
						"destination_cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Destination。",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Next hop 类型 (associated 实例类型)，all types: VPC，DIRECTCONNECT。",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Next jump (associated 实例 ID)。",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Next jump (associated 实例名称)。",
						},
						"instance_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Next jump (associated 实例 地域)。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "更新时间。",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is routing 已启用",
						},
						"instance_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "UIN (root 账号) 到 其中 associated 实例 belongs。",
						},
						"extra_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Extension 状态 routing。",
						},
						"is_bgp": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is 它 动态 routing。",
						},
						"route_priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Routing 优先级",
						},
						"instance_extra_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Next hop extension 名称 (associated 实例 extension 名称)。",
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

func dataSourceTencentCloudCcnRoutesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ccn_routes.read")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("ccn_id"); ok {
		paramMap["CcnId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*vpc.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filter := vpc.Filter{}
			filterMap := item.(map[string]interface{})
			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}

			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	var routeSet []*vpc.CcnRoute
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcDescribeCcnRoutesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		routeSet = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(routeSet))
	tmpList := make([]map[string]interface{}, 0, len(routeSet))

	if routeSet != nil {
		for _, route := range routeSet {
			tmpMap := make(map[string]interface{})
			if route.RouteId != nil {
				tmpMap["route_id"] = route.RouteId
			}

			if route.DestinationCidrBlock != nil {
				tmpMap["destination_cidr_block"] = route.DestinationCidrBlock
			}

			if route.InstanceType != nil {
				tmpMap["instance_type"] = route.InstanceType
			}

			if route.InstanceId != nil {
				tmpMap["instance_id"] = route.InstanceId
			}

			if route.InstanceName != nil {
				tmpMap["instance_name"] = route.InstanceName
			}

			if route.InstanceRegion != nil {
				tmpMap["instance_region"] = route.InstanceRegion
			}

			if route.UpdateTime != nil {
				tmpMap["update_time"] = route.UpdateTime
			}

			if route.Enabled != nil {
				tmpMap["enabled"] = route.Enabled
			}

			if route.InstanceUin != nil {
				tmpMap["instance_uin"] = route.InstanceUin
			}

			if route.ExtraState != nil {
				tmpMap["extra_state"] = route.ExtraState
			}

			if route.IsBgp != nil {
				tmpMap["is_bgp"] = route.IsBgp
			}

			if route.RoutePriority != nil {
				tmpMap["route_priority"] = route.RoutePriority
			}

			if route.InstanceExtraName != nil {
				tmpMap["instance_extra_name"] = route.InstanceExtraName
			}

			ids = append(ids, *route.RouteId)
			tmpList = append(tmpList, tmpMap)
		}

		_ = d.Set("route_list", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
