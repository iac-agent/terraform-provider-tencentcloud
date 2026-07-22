package dayu

import (
	"context"
	"errors"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"
)

func ResourceTencentCloudDayuDdosPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuDdosPolicyCreate,
		Read:   resourceTencentCloudDayuDdosPolicyRead,
		Update: resourceTencentCloudDayuDdosPolicyUpdate,
		Delete: resourceTencentCloudDayuDdosPolicyDelete,

		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_RESOURCE_TYPE),
				ForceNew:     true,
				Description:  "类型 资源 该 DDoS 策略 works 对于. 有效值：`bgpip`，`bgp`，`bgp-multip` 和 `net`。",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 32),
				Description:  "名称 DDoS 策略. Length should between 1 和 32。",
			},
			"drop_options": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"drop_tcp": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicate 是否drop TCP 协议 或 不。",
						},
						"drop_udp": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicate 到 drop UDP 协议 或 不。",
						},
						"drop_icmp": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicate 是否drop ICMP 协议 或 不。",
						},
						"drop_other": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicate 是否drop other protocols(exclude TCP/UDP/ICMP) 或 不。",
						},
						"drop_abroad": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicate 是否drop abroad 流量 或 不。",
						},
						"check_sync_conn": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicate 是否check null 连接 或 不。",
						},
						"d_new_limit": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 4294967295),
							Description:  "限制 的 new connections based 在 destination IP. 有效 值 ranges: (0~4294967295)。",
						},
						"d_conn_limit": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 4294967295),
							Description:  "限制 的 concurrent connections based 在 destination IP. 有效 值 ranges: (0~4294967295)。",
						},
						"s_new_limit": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 4294967295),
							Description:  "限制 的 new connections based 在 来源 IP. 有效 值 ranges: (0~4294967295)。",
						},
						"s_conn_limit": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 4294967295),
							Description:  "限制 的 concurrent connections based 在 来源 IP. 有效 值 ranges: (0~4294967295)。",
						},
						"bad_conn_threshold": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 4294967295),
							Description:  "数量 new connections based 在 destination IP 该 触发器 suppression 的 connections. 有效 值 ranges: (0~4294967295)。",
						},
						"null_conn_enable": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Indicate 到 启用 null 连接 或 不。",
						},
						"conn_timeout": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 65535),
							Description:  "Connection 超时 的 abnormal 连接 check. 有效 值 ranges: (0~65535)。",
						},
						"syn_rate": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 100),
							Description:  "percentage 的 syn 在 ack 的 abnormal 连接 check. 有效 值 ranges: (0~100)。",
						},
						"syn_limit": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 100),
							Description:  "限制 的 syn 的 abnormal 连接 check. 有效 值 ranges: (0~100)。",
						},
						"tcp_mbps_limit": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 4294967295),
							Description:  "限制 的 TCP 流量. 有效 值 ranges: (0~4294967295)(Mbps)。",
						},
						"udp_mbps_limit": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 4294967295),
							Description:  "限制 的 UDP 流量 速率. 有效 值 ranges: (0~4294967295)(Mbps)。",
						},
						"icmp_mbps_limit": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 4294967295),
							Description:  "限制 的 ICMP 流量 速率. 有效 值 ranges: (0~4294967295)(Mbps)。",
						},
						"other_mbps_limit": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 4294967295),
							Description:  "限制 的 other protocols(exclude TCP/UDP/ICMP) 流量 速率. 有效 值 ranges: (0~4294967295)(Mbps)。",
						},
					},
				},
				Description: "Option 列表 abnormal check 的 DDos 策略，should 集合 在 least 一个 策略。",
			},
			"port_filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_PROTOCOL),
							Description:  "协议 有效 值 是 `tcp`，`udp`，`icmp`，`all`。",
						},
						"start_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: tccommon.ValidatePort,
							Description:  "Start 端口 有效 值 ranges: (0~65535)。",
						},
						"end_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      65535,
							ValidateFunc: tccommon.ValidatePort,
							Description:  "End 端口 有效 值 ranges: (0~65535). It 必须 是 greater 比 `start_port`。",
						},
						"action": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_PORT_ACTION),
							Description:  "操作 的 端口 到 take. 有效值：`drop`，`transmit`。",
						},
						"kind": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1, 2}),
							Description:  "类型 forbidden 端口 有效值：`0`，`1`，`2`. `0` 对于 destination ports make effect，`1` 对于 来源 ports make effect. `2` 对于 both destination 和 来源 ports。",
						},
					},
				},
				Description: "端口 limits 的 abnormal check 的 DDos 策略。",
			},
			"black_ips": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: tccommon.ValidateIp,
				},
				Optional:    true,
				Description: "Black IP 列表。",
			},
			"white_ips": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: tccommon.ValidateIp,
				},
				Optional:    true,
				Description: "White IP 列表。",
			},
			"packet_filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_PROTOCOL),
							Description:  "协议 有效值：`tcp`，`udp`，`icmp`，`all`。",
						},
						"d_start_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: tccommon.ValidatePort,
							Description:  "Start 端口 的 destination. 有效 值 ranges: (0~65535)。",
						},
						"d_end_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: tccommon.ValidatePort,
							Description:  "End 端口 的 destination. 有效 值 ranges: (0~65535). It 必须 是 greater 比 `d_start_port`。",
						},
						"s_start_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: tccommon.ValidatePort,
							Description:  "Start 端口 的 来源 有效 值 ranges: (0~65535)。",
						},
						"s_end_port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: tccommon.ValidatePort,
							Description:  "End 端口 的 来源 有效 值 ranges: (0~65535). It 必须 是 greater 比 `s_start_port`。",
						},
						"pkt_length_min": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 1500),
							Description:  "最小长度the packet. 有效 值 ranges: (0~1500)(Mbps)。",
						},
						"pkt_length_max": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 1500),
							Description:  "max 长度 的 packet. 有效 值 ranges: (0~1500)(Mbps). It 必须 是 greater 比 `pkt_length_min`。",
						},
						"match_begin": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_MATCH_SWITCH),
							Description:  "Indicate 是否check load 或 不，`begin_l5` 表示 到 match 和 `no_match` 表示 不。",
						},
						"match_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_MATCH_TYPE),
							Description:  "Match 类型 有效值：`sunday` 和 `pcre`. `sunday` 表示 键 word match while `pcre` 表示 regular match。",
						},
						"match_str": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "键 word 或 regular expression。",
						},
						"depth": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 1500),
							Description:  "depth 的 match. 有效 值 ranges: (0~1500)。",
						},
						"offset": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 1500),
							Description:  "偏移量 的 match. 有效 值 ranges: (0~1500)。",
						},
						"is_include": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Indicate 是否include 键 word/regular expression 或 不。",
						},
						"action": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_PACKET_ACTION),
							Description:  "操作 的 端口 到 take. 有效值：`drop`，`drop_black`,`drop_rst`,`drop_black_rst`,`transmit`.`drop`(drop packet)，`drop_black`(drop packet 和 black ip),`drop_rst`(drop packet 和 disconnect),`drop_black_rst`(drop packet，black ip 和 disconnect),`transmit`(transmit packet)。",
						},
					},
				},
				Description: "消息 过滤器 options 列表。",
			},
			"watermark_filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tcp_port_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: tccommon.ValidatePortRange,
							},
							Description: "端口 范围 的 TCP， 格式 是 like `2000-3000`。",
						},
						"udp_port_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: tccommon.ValidatePortRange,
							},
							Description: "端口 范围 的 TCP， 格式 是 like `2000-3000`。",
						},
						"offset": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: tccommon.ValidateIntegerInRange(0, 100),
							Description:  "偏移量 的 水印. 有效 值 ranges: (0~1500)。",
						},
						"auto_remove": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicate 是否auto-remove 水印 或 不。",
						},
						"open_switch": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicate 是否open 水印 或 不. It muse 是 集合 `true` 当 any 字段 的 水印 是 集合。",
						},
					},
				},
				Description: "Watermark 策略 options，和 仅 support 一个 水印 策略 在 most。",
			},
			//computed
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
	}
}

func resourceTencentCloudDayuDdosPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_ddos_policy.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	resourceType := d.Get("resource_type").(string)
	name := d.Get("name").(string)
	//set DDosPolicyDropOption
	dropMapping := d.Get("drop_options").([]interface{})
	ddosPolicyDropOption, _ := setDdosPolicyDropOption(dropMapping)

	//set DDoSPolicyPortLimit
	portMapping := d.Get("port_filters").([]interface{})
	ddosPolicyPortLimit, lErr := setDdosPolicyPortLimit(portMapping)

	if lErr != nil {
		return lErr
	}

	//set IpBlackWhite
	blackIps := d.Get("black_ips").(*schema.Set).List()
	whiteIps := d.Get("white_ips").(*schema.Set).List()
	ipBlackWhite, ipErr := setIpBlackWhite(blackIps, whiteIps)

	if ipErr != nil {
		return ipErr
	}
	//set DDoSPolicyPacketFilter
	packetFilterMapping := d.Get("packet_filters").([]interface{})
	ddosPacketFilter, pErr := setDdosPolicyPacketFilter(packetFilterMapping)
	if pErr != nil {
		return pErr
	}

	//set WaterPrintPolicy
	waterPrintMapping := d.Get("watermark_filters").([]interface{})
	waterPrintPolicy, _ := setWaterPrintPolicy(waterPrintMapping)

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	policyId := ""

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := dayuService.CreateDdosPolicy(ctx, resourceType, name, ddosPolicyDropOption, ddosPolicyPortLimit, ipBlackWhite, ddosPacketFilter, waterPrintPolicy)
		if e != nil {
			return tccommon.RetryError(e)
		}
		policyId = result
		return nil
	})

	if err != nil {
		return err
	}

	d.SetId(resourceType + tccommon.FILED_SP + policyId)

	return resourceTencentCloudDayuDdosPolicyRead(d, meta)
}

func resourceTencentCloudDayuDdosPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_ddos_policy.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of DDoS policy")
	}
	resourceType := items[0]
	policyId := items[1]

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	ddosPolicy, has, err := dayuService.DescribeDdosPolicy(ctx, resourceType, policyId)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ddosPolicy, has, err = dayuService.DescribeDdosPolicy(ctx, resourceType, policyId)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}
	_ = d.Set("drop_options", flattenDdosDropOptionList([]*dayu.DDoSPolicyDropOption{ddosPolicy.DropOptions}))
	_ = d.Set("port_filters", flattenDdosPortLimitList(ddosPolicy.PortLimits))
	_ = d.Set("packet_filters", flattenDdosPacketFilterList(ddosPolicy.PacketFilters))
	blackIps, whiteIps := flattenIpBlackWhiteList(ddosPolicy.IpBlackWhiteLists)
	_ = d.Set("black_ips", blackIps)
	_ = d.Set("white_ips", whiteIps)
	_ = d.Set("watermark_filters", flattenWaterPrintPolicyList(ddosPolicy.WaterPrint))
	_ = d.Set("create_time", ddosPolicy.CreateTime)
	_ = d.Set("name", ddosPolicy.PolicyName)
	_ = d.Set("scene_id", ddosPolicy.SceneId)
	_ = d.Set("policy_id", ddosPolicy.PolicyId)
	_ = d.Set("watermark_key", flattenWaterPrintKeyList(ddosPolicy.WaterKey))

	return nil
}

func resourceTencentCloudDayuDdosPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_ddos_policy.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of DDoS policy")
	}
	resourceType := items[0]
	policyId := items[1]
	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	d.Partial(true)

	if d.HasChange("name") {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := dayuService.ModifyDdosPolicyName(ctx, resourceType, policyId, d.Get("name").(string))
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})
		if err != nil {
			return err
		}

	}

	if d.HasChange("watermark_filters") || d.HasChange("white_ips") || d.HasChange("black_ips") || d.HasChange("packet_filters") || d.HasChange("port_filters") || d.HasChange("drop_options") {

		//set DDosPolicyDropOption
		dropMapping := d.Get("drop_options").([]interface{})
		ddosPolicyDropOption, _ := setDdosPolicyDropOption(dropMapping)

		//set DDoSPolicyPortLimit
		portMapping := d.Get("port_filters").([]interface{})
		ddosPolicyPortLimit, lErr := setDdosPolicyPortLimit(portMapping)

		if lErr != nil {
			return lErr
		}

		//set IpBlackWhite
		blackIps := d.Get("black_ips").(*schema.Set).List()
		whiteIps := d.Get("white_ips").(*schema.Set).List()
		ipBlackWhite, ipErr := setIpBlackWhite(blackIps, whiteIps)

		if ipErr != nil {
			return ipErr
		}
		//set DDoSPolicyPacketFilter
		packetFilterMapping := d.Get("packet_filters").([]interface{})
		ddosPacketFilter, pErr := setDdosPolicyPacketFilter(packetFilterMapping)
		if pErr != nil {
			return pErr
		}

		//set WaterPrintPolicy
		waterPrintMapping := d.Get("watermark_filters").([]interface{})
		waterPrintPolicy, _ := setWaterPrintPolicy(waterPrintMapping)

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := dayuService.ModifyDdosPolicy(ctx, resourceType, policyId, ddosPolicyDropOption, ddosPolicyPortLimit, ipBlackWhite, ddosPacketFilter, waterPrintPolicy)
			if e != nil {
				return tccommon.RetryError(e)
			}
			return nil
		})
		if err != nil {
			return err
		}

	}

	d.Partial(false)

	return resourceTencentCloudDayuDdosPolicyRead(d, meta)
}

func resourceTencentCloudDayuDdosPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_ddos_policy.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 2 {
		return fmt.Errorf("broken ID of DDoS policy")
	}
	resourceType := items[0]
	policyId := items[1]

	dayuService := DayuService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := dayuService.DeleteDdosPolicy(ctx, resourceType, policyId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	_, has, err := dayuService.DescribeDdosPolicy(ctx, resourceType, policyId)
	if err != nil || has {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, has, err = dayuService.DescribeDdosPolicy(ctx, resourceType, policyId)
			if err != nil {
				return tccommon.RetryError(err)
			}

			if has {
				err = fmt.Errorf("delete DDoS policy fail, DDoS policy still exist from sdk DescribeDDosPolicy")
				return resource.RetryableError(err)
			}

			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		return nil
	} else {
		return errors.New("delete DDoS policy fail, DDoS policy still exist from sdk DescribeDDosPolicy")
	}
}
