package vpc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpcNetDetectStates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcNetDetectStatesRead,
		Schema: map[string]*schema.Schema{
			"net_detect_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "数组 网络 detection 实例 `IDs`，such 作为 [`netd-12345678`]。",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 conditions. `NetDetectIds` 和 `Filters` 不能 是 指定 在 same 时间.net-detect-ID - String - (过滤器 condition) 网络 detection 实例 ID，such 作为 netd-12345678。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "attribute 名称 如果 more 比 一个 过滤器 exists， logical relation between these Filters 是 `AND`。",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Attribute 值 如果 多个 值 exist 在 一个 过滤器， logical relationship between these 值 是 `OR`. For `bool` 参数， 有效 值 include `TRUE` 和 `FALSE`。",
						},
					},
				},
			},

			"net_detect_state_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "数组 网络 detection verification results 该 meet requirements.注意：此字段可能返回 null，表示无法获取有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"net_detect_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 网络 detection 实例，such 作为 netd-12345678。",
						},
						"net_detect_ip_state_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "数组 网络 detection destination IP verification results。",
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

func dataSourceTencentCloudVpcNetDetectStatesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_net_detect_states.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("net_detect_ids"); ok {
		netDetectIdsSet := v.(*schema.Set).List()
		paramMap["NetDetectIds"] = helper.InterfacesStringsPoint(netDetectIdsSet)
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

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var netDetectStateSet []*vpc.NetDetectState

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcNetDetectStatesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		netDetectStateSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(netDetectStateSet))
	tmpList := make([]map[string]interface{}, 0, len(netDetectStateSet))

	if netDetectStateSet != nil {
		for _, netDetectState := range netDetectStateSet {
			netDetectStateMap := map[string]interface{}{}

			if netDetectState.NetDetectId != nil {
				netDetectStateMap["net_detect_id"] = netDetectState.NetDetectId
			}

			if netDetectState.NetDetectIpStateSet != nil {
				netDetectIpStateSetList := []interface{}{}
				for _, netDetectIpStateSet := range netDetectState.NetDetectIpStateSet {
					netDetectIpStateSetMap := map[string]interface{}{}

					if netDetectIpStateSet.DetectDestinationIp != nil {
						netDetectIpStateSetMap["detect_destination_ip"] = netDetectIpStateSet.DetectDestinationIp
					}

					if netDetectIpStateSet.State != nil {
						netDetectIpStateSetMap["state"] = netDetectIpStateSet.State
					}

					if netDetectIpStateSet.Delay != nil {
						netDetectIpStateSetMap["delay"] = netDetectIpStateSet.Delay
					}

					if netDetectIpStateSet.PacketLossRate != nil {
						netDetectIpStateSetMap["packet_loss_rate"] = netDetectIpStateSet.PacketLossRate
					}

					netDetectIpStateSetList = append(netDetectIpStateSetList, netDetectIpStateSetMap)
				}

				netDetectStateMap["net_detect_ip_state_set"] = netDetectIpStateSetList
			}

			ids = append(ids, *netDetectState.NetDetectId)
			tmpList = append(tmpList, netDetectStateMap)
		}

		_ = d.Set("net_detect_state_set", tmpList)
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
