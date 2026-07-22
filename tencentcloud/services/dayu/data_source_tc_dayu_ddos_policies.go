package dayu

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDayuDdosPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDayuDdosPoliciesRead,
		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_RESOURCE_TYPE),
				Description:  "类型 resource that the DDoS policy works for，valid values are `bgpip`，`bgp`，`bgp-multip` and `net`。",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID DDoS policy to be query。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 DDoS policies. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 DDoS policy。",
						},
						"drop_options": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"drop_tcp": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate 是否drop TCP 协议 or not。",
									},
									"drop_udp": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate to drop UDP 协议 or not。",
									},
									"drop_icmp": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate 是否drop ICMP 协议 or not。",
									},
									"drop_other": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate 是否drop other protocols(exclude TCP/UDP/ICMP) or not。",
									},
									"drop_abroad": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Indicate 是否drop abroad traffic or not。",
									},
									"check_sync_conn": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate 是否check null connection or not。",
									},
									"d_new_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The 限制 of new connections based on destination IP。",
									},
									"d_conn_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The limit of concurrent connections based on destination IP.", //?
									},
									"s_conn_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The 限制 of concurrent connections based on 来源 IP。",
									},
									"s_new_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The 限制 of new connections based on 来源 IP。",
									},
									"bad_conn_threshold": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The 数量 new connections based on destination IP that trigger suppression of connections。",
									},
									"null_conn_enable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate to enable null connection or not。",
									},
									"conn_timeout": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Connection timeout of abnormal connection check。",
									},
									"syn_rate": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The percentage of syn in ack of abnormal connection check。",
									},
									"syn_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The 限制 of syn of abnormal connection check。",
									},
									"tcp_mbps_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The 限制 of TCP traffic。",
									},
									"udp_mbps_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The 限制 of UDP traffic rate。",
									},
									"icmp_mbps_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The 限制 of ICMP traffic rate。",
									},
									"other_mbps_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The 限制 of other protocols(exclude TCP/UDP/ICMP) traffic rate。",
									},
								},
							},
							Description: "Option 列表 abnormal check of the DDoS policy。",
						},
						"port_filters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "协议",
									},
									"start_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Start 端口",
									},
									"end_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "End 端口",
									},
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "操作 of 端口 to take。",
									},
									"kind": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "类型 forbidden 端口，and valid values are 0，1，2. 0 for destination 端口，1 for 来源 端口 and 2 for both destination and 来源 posts。",
									},
								},
							},
							Description: "端口 limits of abnormal check of the DDoS policy。",
						},
						"black_ips": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: tccommon.ValidateIp,
							},
							Optional:    true,
							Description: "Black ip list。",
						},
						"white_ips": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: tccommon.ValidateIp,
							},
							Optional:    true,
							Description: "White ip list。",
						},
						"packet_filters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "协议",
									},
									"d_start_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Start 端口 of the destination。",
									},
									"d_end_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "End 端口 of the destination。",
									},
									"s_start_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Start 端口 of the 来源",
									},
									"s_end_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "End 端口 of the 来源",
									},
									"pkt_length_min": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The 最小长度the packet。",
									},
									"pkt_length_max": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The max length of the packet。",
									},
									"match_begin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Indicate 是否check load or not。",
									},
									"match_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Match 类型",
									},
									"match_str": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The 键 word or regular expression。",
									},
									"depth": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The depth of match。",
									},
									"offset": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The 偏移量 of match。",
									},
									"is_include": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate 是否include the 键 word/regular expression or not。",
									},
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "操作 of 端口 to take。",
									},
								},
							},
							Description: "消息 filter options list。",
						},
						"watermark_filters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tcp_port_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "端口 range of TCP。",
									},
									"udp_port_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "端口 range of TCP。",
									},
									"offset": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The 偏移量 of watermark。",
									},
									"auto_remove": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate 是否auto-remove the watermark or not。",
									},
									"open_switch": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate 是否open watermark or not。",
									},
								},
							},
							Description: "Watermark policy options，and only support one watermark policy at most。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 of the DDoS policy。",
						},
						"scene_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of policy case that the DDoS policy works for。",
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of policy。",
						},
						"watermark_key": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Id of the watermark。",
									},
									"content": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "内容 of the watermark。",
									},
									"open_switch": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate 是否auto-remove the watermark or not。",
									},
									"create_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "创建时间 of the watermark。",
									},
								},
							},
							Description: "Watermark 内容",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDayuDdosPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dayu_ddos_policies.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DayuService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	resourceType := d.Get("resource_type").(string)
	policyId := d.Get("policy_id").(string)

	policies := make([]*dayu.DDosPolicy, 0)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, err := service.DescribeDdosPolicies(ctx, resourceType, policyId)
		if err != nil {
			return tccommon.RetryError(err)
		}
		policies = result
		return nil
	})

	if err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(policies))
	ids := make([]string, 0, len(policies))
	for _, ddosPolicy := range policies {
		listItem := make(map[string]interface{})
		listItem["drop_options"] = flattenDdosDropOptionList([]*dayu.DDoSPolicyDropOption{ddosPolicy.DropOptions})
		listItem["port_filters"] = flattenDdosPortLimitList(ddosPolicy.PortLimits)
		listItem["packet_filters"] = flattenDdosPacketFilterList(ddosPolicy.PacketFilters)
		listItem["black_ips"], listItem["white_ips"] = flattenIpBlackWhiteList(ddosPolicy.IpBlackWhiteLists)
		listItem["watermark_filters"] = flattenWaterPrintPolicyList(ddosPolicy.WaterPrint)
		listItem["create_time"] = *ddosPolicy.CreateTime
		listItem["name"] = *ddosPolicy.PolicyName
		listItem["policy_id"] = *ddosPolicy.PolicyId
		listItem["scene_id"] = *ddosPolicy.SceneId
		listItem["watermark_key"] = flattenWaterPrintKeyList(ddosPolicy.WaterKey)
		list = append(list, listItem)
		ids = append(ids, *ddosPolicy.PolicyId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}
	return nil

}
