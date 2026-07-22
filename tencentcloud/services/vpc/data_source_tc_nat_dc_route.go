package vpc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudNatDcRoute() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudNatDcRouteRead,
		Schema: map[string]*schema.Schema{
			"nat_gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Unique identifier 的 Nat Gateway。",
			},

			"vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Unique identifier 的 Vpc。",
			},

			"nat_direct_connect_gateway_route_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Data 的 路由。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IPv4 CIDR 的 子网。",
						},
						"gateway_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 next-hop 网关，有效值：DIRECTCONNECT。",
						},
						"gateway_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 的 next-hop 网关。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 路由。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "更新时间 的 路由。",
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

func dataSourceTencentCloudNatDcRouteRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_nat_dc_route.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("nat_gateway_id"); ok {
		paramMap["NatGatewayId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		paramMap["VpcId"] = helper.String(v.(string))
	}

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var natDirectConnectGatewayRouteSet []*vpc.NatDirectConnectGatewayRoute

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeNatDcRouteByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		natDirectConnectGatewayRouteSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(natDirectConnectGatewayRouteSet))
	tmpList := make([]map[string]interface{}, 0, len(natDirectConnectGatewayRouteSet))

	if natDirectConnectGatewayRouteSet != nil {
		for _, natDirectConnectGatewayRoute := range natDirectConnectGatewayRouteSet {
			natDirectConnectGatewayRouteMap := map[string]interface{}{}

			if natDirectConnectGatewayRoute.DestinationCidrBlock != nil {
				natDirectConnectGatewayRouteMap["destination_cidr_block"] = natDirectConnectGatewayRoute.DestinationCidrBlock
			}

			if natDirectConnectGatewayRoute.GatewayType != nil {
				natDirectConnectGatewayRouteMap["gateway_type"] = natDirectConnectGatewayRoute.GatewayType
			}

			if natDirectConnectGatewayRoute.GatewayId != nil {
				natDirectConnectGatewayRouteMap["gateway_id"] = natDirectConnectGatewayRoute.GatewayId
			}

			if natDirectConnectGatewayRoute.CreateTime != nil {
				natDirectConnectGatewayRouteMap["create_time"] = natDirectConnectGatewayRoute.CreateTime
			}

			if natDirectConnectGatewayRoute.UpdateTime != nil {
				natDirectConnectGatewayRouteMap["update_time"] = natDirectConnectGatewayRoute.UpdateTime
			}

			ids = append(ids, *natDirectConnectGatewayRoute.GatewayId)
			tmpList = append(tmpList, natDirectConnectGatewayRouteMap)
		}

		_ = d.Set("nat_direct_connect_gateway_route_set", tmpList)
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
