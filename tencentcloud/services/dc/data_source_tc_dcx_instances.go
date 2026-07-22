package dc

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTencentCloudDcxInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcxInstancesRead,

		Schema: map[string]*schema.Schema{
			"dcx_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID dedicated tunnels 到 是 queried。",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 dedicated tunnels 到 是 queried。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// Computed values
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information 列表 dedicated tunnels。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dcx_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID dedicated tunnel。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 dedicated tunnel。",
						},
						"network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 网络. 有效值：`VPC`，`BMVPC` 和 `CCN`. 默认值为 `VPC`。",
						},
						"dcg_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID DC Gateway. Currently 仅 new 在 console。",
						},
						"network_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 的 dedicated tunnel。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID VPC 或 BMVPC。",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Bandwidth 的 DC。",
						},
						"route_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 路由. 有效值：`BGP` 和 `STATIC`. 默认值为 `BGP`。",
						},
						"bgp_asn": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "BGP ASN 的 用户",
						},
						"bgp_auth_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "BGP 键 的 用户",
						},
						"route_filter_prefixes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Static 路由， 网络 地址 的 用户 IDC。",
						},
						"vlan": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Vlan 的 dedicated tunnels. 有效 值 ranges: [0-3000]. `0` 表示 该 仅 一个 tunnel 可以 是 创建 对于 physical connect。",
						},
						"tencent_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Interconnect IP 的 DC within Tencent。",
						},
						"customer_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Interconnect IP 的 DC within 客户端。",
						},
						"dc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID DC。",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State 的 dedicated tunnels. 有效值：`PENDING`，`ALLOCATING`，`ALLOCATED`，`ALTERING`，`DELETING`，`DELETED`，`COMFIRMING` 和 `REJECTED`。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 资源。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDcxInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dcx_instances.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		id   = ""
		name = ""
	)
	if temp, ok := d.GetOk("dcx_id"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			id = tempStr
		}
	}
	if temp, ok := d.GetOk("name"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			name = tempStr
		}
	}

	var infos, err = service.DescribeDirectConnectTunnels(ctx, id, name)

	if err != nil {
		return err
	}
	var instanceList = make([]map[string]interface{}, 0, len(infos))

	for _, item := range infos {

		var infoMap = make(map[string]interface{})
		infoMap["dcx_id"] = *item.DirectConnectTunnelId
		infoMap["name"] = *item.DirectConnectTunnelName
		infoMap["network_type"] = strings.ToUpper(service.strPt2str(item.NetworkType))

		infoMap["network_region"] = service.strPt2str(item.NetworkRegion)
		infoMap["vpc_id"] = service.strPt2str(item.VpcId)
		infoMap["bandwidth"] = service.int64Pt2int64(item.Bandwidth)

		infoMap["route_type"] = strings.ToUpper(service.strPt2str(item.RouteType))

		if item.BgpPeer == nil {
			infoMap["bgp_asn"] = 0
			infoMap["bgp_auth_key"] = ""
		} else {
			infoMap["bgp_asn"] = service.int64Pt2int64(item.BgpPeer.Asn)
			infoMap["bgp_auth_key"] = service.strPt2str(item.BgpPeer.AuthKey)
		}

		infoMap["vlan"] = service.int64Pt2int64(item.Vlan)
		infoMap["tencent_address"] = service.strPt2str(item.TencentAddress)
		infoMap["customer_address"] = service.strPt2str(item.CustomerAddress)
		infoMap["dcg_id"] = service.strPt2str(item.DirectConnectGatewayId)

		infoMap["dc_id"] = service.strPt2str(item.DirectConnectId)
		infoMap["state"] = strings.ToUpper(service.strPt2str(item.State))
		infoMap["create_time"] = service.strPt2str(item.CreatedTime)

		var routeFilterPrefixes = make([]string, 0, len(item.RouteFilterPrefixes))
		for _, v := range item.RouteFilterPrefixes {
			if v.Cidr != nil {
				routeFilterPrefixes = append(routeFilterPrefixes, *v.Cidr)
			}
		}
		infoMap["route_filter_prefixes"] = routeFilterPrefixes

		instanceList = append(instanceList, infoMap)
	}

	if err := d.Set("instance_list", instanceList); err != nil {
		log.Printf("[CRITAL]%s provider set  dcx instances fail, reason:%s\n ", logId, err.Error())
		return err
	}

	m := md5.New()
	_, err = m.Write([]byte(id + "_" + name))
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%x", m.Sum(nil)))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), instanceList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}
	return nil
}
