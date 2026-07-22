package vpc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpcNetDetectStateCheck() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcNetDetectStateCheckRead,
		Schema: map[string]*schema.Schema{
			"detect_destination_ip": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "数组 detection destination IPv4 addresses，其中 包含at most two IP addresses。",
			},

			"next_hop_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "类型 next hop. Currently 支持 types 是:VPN: VPN 网关;DIRECTCONNECT: direct connect 网关;PEERCONNECTION: peering 连接;NAT: NAT 网关;NORMAL_CVM: normal CVM。",
			},

			"next_hop_destination": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "next-hop destination 网关. 值 是 related 到 NextHopType.如果 NextHopType 是 集合 到 VPN， 值 的 此 参数 是 VPN 网关 ID，such 作为 vpngw-12345678.如果 NextHopType 是 集合 到 DIRECTCONNECT， 值 的 此 参数 是 direct connect 网关 ID，such 作为 dcg-12345678.如果 NextHopType 是 集合 到 PEERCONNECTION， 值 的 此 参数 是 peering 连接 ID，such 作为 pcx-12345678.如果 NextHopType 是 集合 到 NAT， 值 的 此 参数 是 NAT 网关 ID，such 作为 nat-12345678.如果 NextHopType 是 集合 到 NORMAL_CVM， 值 的 此 参数 是 IPv4 地址 的 CVM，such 作为 10.0.0.12。",
			},

			"net_detect_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ID 网络 inspector 实例，e.g. netd-12345678. Enter 在 least 一个 的 此 参数，VpcId，SubnetId，和 NetDetectName. Use NetDetectId 如果 它 是 present。",
			},

			"vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ID `VPC` 实例，e.g. `vpc-12345678`，其中 是 使用 together 使用 SubnetId 和 NetDetectName. You should enter either 此 参数 或 NetDetectId，或 both. Use NetDetectId 如果 它 是 present。",
			},

			"subnet_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ID 子网 实例，e.g. `子网-12345678`，其中 是 使用 together 使用 VpcId 和 NetDetectName. You should enter either 此 参数 或 NetDetectId，或 both. Use NetDetectId 如果 它 是 present。",
			},

			"net_detect_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "名称 网络 inspector，up 到 60 bytes 在 长度. It 是 使用 together 使用 VpcId 和 NetDetectName. You should enter either 此 参数 或 NetDetectId，或 both. Use NetDetectId 如果 它 是 present。",
			},

			"net_detect_ip_state_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "数组 网络 detection verification results。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"detect_destination_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "destination IPv4 地址 的 网络 detection。",
						},
						"state": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "detection 结果0: successful;-1: 无 packet loss occurred during routing;-2: packet loss occurred 当 outbound 流量 是 blocked 通过 ACL;-3: packet loss occurred 当 inbound 流量 是 blocked 通过 ACL;-4: other errors。",
						},
						"delay": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "延迟. 单位：ms。",
						},
						"packet_loss_rate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "packet loss 速率。",
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

func dataSourceTencentCloudVpcNetDetectStateCheckRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_net_detect_state_check.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("detect_destination_ip"); ok {
		detectDestinationIpSet := v.(*schema.Set).List()
		paramMap["DetectDestinationIp"] = helper.InterfacesStringsPoint(detectDestinationIpSet)
	}

	if v, ok := d.GetOk("next_hop_type"); ok {
		paramMap["NextHopType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("next_hop_destination"); ok {
		paramMap["NextHopDestination"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("net_detect_id"); ok {
		paramMap["NetDetectId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		paramMap["VpcId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		paramMap["SubnetId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("net_detect_name"); ok {
		paramMap["NetDetectName"] = helper.String(v.(string))
	}

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var netDetectIpStateSet []*vpc.NetDetectIpState

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcNetDetectStateCheck(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		netDetectIpStateSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(netDetectIpStateSet))
	tmpList := make([]map[string]interface{}, 0, len(netDetectIpStateSet))

	if netDetectIpStateSet != nil {
		for _, netDetectIpState := range netDetectIpStateSet {
			netDetectIpStateMap := map[string]interface{}{}

			if netDetectIpState.DetectDestinationIp != nil {
				netDetectIpStateMap["detect_destination_ip"] = netDetectIpState.DetectDestinationIp
			}

			if netDetectIpState.State != nil {
				netDetectIpStateMap["state"] = netDetectIpState.State
			}

			if netDetectIpState.Delay != nil {
				netDetectIpStateMap["delay"] = netDetectIpState.Delay
			}

			if netDetectIpState.PacketLossRate != nil {
				netDetectIpStateMap["packet_loss_rate"] = netDetectIpState.PacketLossRate
			}

			ids = append(ids, *netDetectIpState.DetectDestinationIp)
			tmpList = append(tmpList, netDetectIpStateMap)
		}

		_ = d.Set("net_detect_ip_state_set", tmpList)
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
