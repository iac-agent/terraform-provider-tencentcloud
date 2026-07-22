package vpc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpcGatewayFlowMonitorDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcGatewayFlowMonitorDetailRead,
		Schema: map[string]*schema.Schema{
			"time_point": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "point 在 时间. 此 表示details 的 此 minute 将 是 queried. For 示例，在 `2019-02-28 18:15:20`，details 在 `18:15` 将 是 queried。",
			},

			"vpn_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "实例 ID VPN 网关，such 作为 `vpn-ltjahce6`。",
			},

			"direct_connect_gateway_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "实例 ID Direct Connect 网关，such 作为 `dcg-ltjahce6`。",
			},

			"peering_connection_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "实例 ID peering 连接，such 作为 `pcx-ltjahce6`。",
			},

			"nat_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "实例 ID NAT 网关，such 作为 `nat-ltjahce6`。",
			},

			"order_field": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "顺序 字段 支持 `InPkg`，`OutPkg`，`InTraffic`，和 `OutTraffic`。",
			},

			"order_direction": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "顺序 methods. Ascending: `ASC`，Descending: `DESC`。",
			},

			"gateway_flow_monitor_detail_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "网关 流量 监控 details。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"private_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Origin `IP`。",
						},
						"in_pkg": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Inbound packets。",
						},
						"out_pkg": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Outbound packets。",
						},
						"in_traffic": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Inbound 流量，在 Byte。",
						},
						"out_traffic": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Outbound 流量，在 Byte。",
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

func dataSourceTencentCloudVpcGatewayFlowMonitorDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_gateway_flow_monitor_detail.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("time_point"); ok {
		paramMap["TimePoint"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpn_id"); ok {
		paramMap["VpnId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("direct_connect_gateway_id"); ok {
		paramMap["DirectConnectGatewayId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("peering_connection_id"); ok {
		paramMap["PeeringConnectionId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("nat_id"); ok {
		paramMap["NatId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_field"); ok {
		paramMap["OrderField"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_direction"); ok {
		paramMap["OrderDirection"] = helper.String(v.(string))
	}

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var gatewayFlowMonitorDetailSet []*vpc.GatewayFlowMonitorDetail

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcGatewayFlowMonitorDetailByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		gatewayFlowMonitorDetailSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(gatewayFlowMonitorDetailSet))
	tmpList := make([]map[string]interface{}, 0, len(gatewayFlowMonitorDetailSet))

	if gatewayFlowMonitorDetailSet != nil {
		for _, gatewayFlowMonitorDetail := range gatewayFlowMonitorDetailSet {
			gatewayFlowMonitorDetailMap := map[string]interface{}{}

			if gatewayFlowMonitorDetail.PrivateIpAddress != nil {
				gatewayFlowMonitorDetailMap["private_ip_address"] = gatewayFlowMonitorDetail.PrivateIpAddress
			}

			if gatewayFlowMonitorDetail.InPkg != nil {
				gatewayFlowMonitorDetailMap["in_pkg"] = gatewayFlowMonitorDetail.InPkg
			}

			if gatewayFlowMonitorDetail.OutPkg != nil {
				gatewayFlowMonitorDetailMap["out_pkg"] = gatewayFlowMonitorDetail.OutPkg
			}

			if gatewayFlowMonitorDetail.InTraffic != nil {
				gatewayFlowMonitorDetailMap["in_traffic"] = gatewayFlowMonitorDetail.InTraffic
			}

			if gatewayFlowMonitorDetail.OutTraffic != nil {
				gatewayFlowMonitorDetailMap["out_traffic"] = gatewayFlowMonitorDetail.OutTraffic
			}

			ids = append(ids, *gatewayFlowMonitorDetail.PrivateIpAddress)
			tmpList = append(tmpList, gatewayFlowMonitorDetailMap)
		}

		_ = d.Set("gateway_flow_monitor_detail_set", tmpList)
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
