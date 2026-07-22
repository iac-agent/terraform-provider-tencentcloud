package vpc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpcRouteConflicts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcRouteConflictsRead,
		Schema: map[string]*schema.Schema{
			"route_table_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Routing 表 实例 ID，对于 示例:rtb-azd4dt1c。",
			},

			"destination_cidr_blocks": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "列表 conflicting destinations 到 check 对于。",
			},

			"route_conflict_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "路由 conflict 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"route_table_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "路由 表 ID。",
						},
						"destination_cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "destination cidr block。",
						},
						"conflict_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "路由 conflict 列表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destination_cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Destination Cidr Block，like 112.20.51.0/24。",
									},
									"gateway_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "next 网关 类型",
									},
									"gateway_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "next hop ID。",
									},
									"route_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "路由 ID。",
									},
									"route_description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "路由 描述",
									},
									"enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "如果 已启用",
									},
									"route_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "routr 类型",
									},
									"route_table_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "路由 表 ID。",
									},
									"destination_ipv6_cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Destination 的 Ipv6 Cidr Block。",
									},
									"route_item_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "唯一 策略 ID。",
									},
									"published_to_vbc": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "如果 published To ccn。",
									},
									"created_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创建时间。",
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
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudVpcRouteConflictsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_route_conflicts.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("route_table_id"); ok {
		paramMap["RouteTableId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("destination_cidr_blocks"); ok {
		destinationCidrBlocksSet := v.(*schema.Set).List()
		paramMap["DestinationCidrBlocks"] = helper.InterfacesStringsPoint(destinationCidrBlocksSet)
	}

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var routeConflictSet []*vpc.RouteConflict

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcRouteConflicts(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		routeConflictSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(routeConflictSet))
	tmpList := make([]map[string]interface{}, 0, len(routeConflictSet))

	if routeConflictSet != nil {
		for _, routeConflict := range routeConflictSet {
			routeConflictMap := map[string]interface{}{}

			if routeConflict.RouteTableId != nil {
				routeConflictMap["route_table_id"] = routeConflict.RouteTableId
			}

			if routeConflict.DestinationCidrBlock != nil {
				routeConflictMap["destination_cidr_block"] = routeConflict.DestinationCidrBlock
			}

			if routeConflict.ConflictSet != nil {
				conflictSetList := []interface{}{}
				for _, conflictSet := range routeConflict.ConflictSet {
					conflictSetMap := map[string]interface{}{}

					if conflictSet.DestinationCidrBlock != nil {
						conflictSetMap["destination_cidr_block"] = conflictSet.DestinationCidrBlock
					}

					if conflictSet.GatewayType != nil {
						conflictSetMap["gateway_type"] = conflictSet.GatewayType
					}

					if conflictSet.GatewayId != nil {
						conflictSetMap["gateway_id"] = conflictSet.GatewayId
					}

					if conflictSet.RouteId != nil {
						conflictSetMap["route_id"] = conflictSet.RouteId
					}

					if conflictSet.RouteDescription != nil {
						conflictSetMap["route_description"] = conflictSet.RouteDescription
					}

					if conflictSet.Enabled != nil {
						conflictSetMap["enabled"] = conflictSet.Enabled
					}

					if conflictSet.RouteType != nil {
						conflictSetMap["route_type"] = conflictSet.RouteType
					}

					if conflictSet.RouteTableId != nil {
						conflictSetMap["route_table_id"] = conflictSet.RouteTableId
					}

					if conflictSet.DestinationIpv6CidrBlock != nil {
						conflictSetMap["destination_ipv6_cidr_block"] = conflictSet.DestinationIpv6CidrBlock
					}

					if conflictSet.RouteItemId != nil {
						conflictSetMap["route_item_id"] = conflictSet.RouteItemId
					}

					if conflictSet.PublishedToVbc != nil {
						conflictSetMap["published_to_vbc"] = conflictSet.PublishedToVbc
					}

					if conflictSet.CreatedTime != nil {
						conflictSetMap["created_time"] = conflictSet.CreatedTime
					}

					conflictSetList = append(conflictSetList, conflictSetMap)
				}

				routeConflictMap["conflict_set"] = conflictSetList
			}

			ids = append(ids, *routeConflict.RouteTableId)
			tmpList = append(tmpList, routeConflictMap)
		}

		_ = d.Set("route_conflict_set", tmpList)
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
