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
				Description:  "类型 资源 该 DDoS 策略 works 对于，有效 值 是 `bgpip`，`bgp`，`bgp-multip` 和 `net`。",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID DDoS 策略 到 是 查询。",
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
							Description: "名称 DDoS 策略。",
						},
						"drop_options": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"drop_tcp": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate 是否drop TCP 协议 或 不。",
									},
									"drop_udp": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate 到 drop UDP 协议 或 不。",
									},
									"drop_icmp": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate 是否drop ICMP 协议 或 不。",
									},
									"drop_other": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate 是否drop other protocols(exclude TCP/UDP/ICMP) 或 不。",
									},
									"drop_abroad": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Indicate 是否drop abroad 流量 或 不。",
									},
									"check_sync_conn": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate 是否check null 连接 或 不。",
									},
									"d_new_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "限制 的 new connections based 在 destination IP。",
									},
									"d_conn_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "限制 的 concurrent connections based 在 destination IP.", //?
									},
									"s_conn_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "限制 的 concurrent connections based 在 来源 IP。",
									},
									"s_new_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "限制 的 new connections based 在 来源 IP。",
									},
									"bad_conn_threshold": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 new connections based 在 destination IP 该 触发器 suppression 的 connections。",
									},
									"null_conn_enable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate 到 启用 null 连接 或 不。",
									},
									"conn_timeout": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Connection 超时 的 abnormal 连接 check。",
									},
									"syn_rate": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "percentage 的 syn 在 ack 的 abnormal 连接 check。",
									},
									"syn_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "限制 的 syn 的 abnormal 连接 check。",
									},
									"tcp_mbps_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "限制 的 TCP 流量。",
									},
									"udp_mbps_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "限制 的 UDP 流量 速率。",
									},
									"icmp_mbps_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "限制 的 ICMP 流量 速率。",
									},
									"other_mbps_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "限制 的 other protocols(exclude TCP/UDP/ICMP) 流量 速率。",
									},
								},
							},
							Description: "Option 列表 abnormal check 的 DDoS 策略。",
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
										Description: "操作 的 端口 到 take。",
									},
									"kind": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "类型 forbidden 端口，和 有效 值 是 0，1，2. 0 对于 destination 端口，1 对于 来源 端口 和 2 对于 both destination 和 来源 posts。",
									},
								},
							},
							Description: "端口 limits 的 abnormal check 的 DDoS 策略。",
						},
						"black_ips": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: tccommon.ValidateIp,
							},
							Optional:    true,
							Description: "Black ip 列表。",
						},
						"white_ips": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: tccommon.ValidateIp,
							},
							Optional:    true,
							Description: "White ip 列表。",
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
										Description: "Start 端口 的 destination。",
									},
									"d_end_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "End 端口 的 destination。",
									},
									"s_start_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Start 端口 的 来源",
									},
									"s_end_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "End 端口 的 来源",
									},
									"pkt_length_min": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "最小长度the packet。",
									},
									"pkt_length_max": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "max 长度 的 packet。",
									},
									"match_begin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Indicate 是否check load 或 不。",
									},
									"match_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Match 类型",
									},
									"match_str": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "键 word 或 regular expression。",
									},
									"depth": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "depth 的 match。",
									},
									"offset": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "偏移量 的 match。",
									},
									"is_include": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate 是否include 键 word/regular expression 或 不。",
									},
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "操作 的 端口 到 take。",
									},
								},
							},
							Description: "消息 过滤器 options 列表。",
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
										Description: "端口 范围 的 TCP。",
									},
									"udp_port_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "端口 范围 的 TCP。",
									},
									"offset": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "偏移量 的 水印。",
									},
									"auto_remove": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate 是否auto-remove 水印 或 不。",
									},
									"open_switch": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate 是否open 水印 或 不。",
									},
								},
							},
							Description: "Watermark 策略 options，和 仅 support 一个 水印 策略 在 most。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 DDoS 策略。",
						},
						"scene_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 的 策略 case 该 DDoS 策略 works 对于。",
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 的 策略。",
						},
						"watermark_key": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID 的 水印。",
									},
									"content": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "内容 的 水印。",
									},
									"open_switch": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicate 是否auto-remove 水印 或 不。",
									},
									"create_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "创建时间 的 水印。",
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
