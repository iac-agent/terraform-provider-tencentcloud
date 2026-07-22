package dayu

import (
	"context"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDayuDdosPolicyCases() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDayuDdosPolicyCasesRead,
		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_RESOURCE_TYPE),
				Description:  "类型 资源 该 DDoS 策略 case works 对于，有效 值 是 `bgpip`，`bgp`，`bgp-multip` 和 `net`。",
			},
			"scene_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID DDoS 策略 case 到 是 查询。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 DDoS 策略 cases. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 资源 该 DDoS 策略 case works 对于。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 DDoS 策略 case。",
						},
						"platform_types": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type:        schema.TypeString,
								Description: "Platform 的 DDoS 策略 case。",
							},
							Computed:    true,
							Description: "Platform 集合 的 DDoS 策略 case。",
						},
						"app_protocols": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type:        schema.TypeString,
								Description: "App 协议 的 DDoS 策略 case。",
							},
							Computed:    true,
							Description: "App 协议 集合 的 DDoS 策略 case。",
						},
						"app_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "App 类型 DDoS 策略 case。",
						},
						"tcp_start_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start 端口 的 TCP 服务。",
						},
						"tcp_end_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End 端口 的 TCP 服务。",
						},
						"udp_start_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start 端口 的 UDP 服务。",
						},
						"udp_end_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End 端口 的 UDP 服务。",
						},
						"has_abroad": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicate 是否service involves overseas 或 不。",
						},
						"has_initiate_tcp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicate 是否service actively initiates TCP requests 或 不。",
						},
						"has_initiate_udp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicate 是否actively initiate UDP requests 或 不。",
						},
						"has_vpn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicate 是否service involves VPN 服务 或 不。",
						},
						"peer_tcp_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "端口 该 actively initiates TCP requests。",
						},
						"peer_udp_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "端口 该 actively initiates UDP requests。",
						},
						"tcp_footprint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "fixed 签名 的 TCP 协议 load。",
						},
						"udp_footprint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "fixed 签名 的 TCP 协议 load。",
						},
						"web_api_urls": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type:        schema.TypeString,
								Description: "Web API URL",
							},
							Computed:    true,
							Description: "Web API URL 集合。",
						},
						"min_tcp_package_len": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "最小长度TCP 消息 包。",
						},
						"max_tcp_package_len": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "max 长度 的 TCP 消息 包。",
						},
						"min_udp_package_len": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "最小长度UDP 消息 包。",
						},
						"max_udp_package_len": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "max 长度 的 UDP 消息 包。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 DDoS 策略 case。",
						},
						"scene_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID DDoS 策略 case。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDayuDdosPolicyCasesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dayu_ddos_policy_cases.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DayuService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	resourceType := d.Get("resource_type").(string)
	sceneId := d.Get("scene_id").(string)

	var ddosPolicyCase dayu.KeyValueRecord
	has := false
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, flag, err := service.DescribeDdosPolicyCase(ctx, resourceType, sceneId)
		if err != nil {
			return tccommon.RetryError(err)
		}
		ddosPolicyCase = result
		has = flag
		return nil
	})

	if err != nil {
		return err
	}

	var list []map[string]interface{}
	var ids []string
	if has {
		list = make([]map[string]interface{}, 0, 1)
		ids = make([]string, 0, 1)
	}
	listItem := make(map[string]interface{})
	for _, record := range ddosPolicyCase.Record {
		key := *record.Key
		if key == "CaseName" {
			listItem["name"] = *record.Value
		}
		if key == "HasInitiateTcp" {

			listItem["has_initiate_tcp"] = *record.Value
		}
		if key == "HasInitiateUdp" {
			listItem["has_initiate_udp"] = *record.Value
		}
		if key == "HasVPN" {
			listItem["has_vpn"] = *record.Value
		}
		if key == "PeerTcpPort" {
			listItem["peer_tcp_port"] = *record.Value
		}
		if key == "PeerUdpPort" {
			listItem["peer_udp_port"] = *record.Value
		}
		if key == "TcpFootprint" {
			listItem["tcp_footprint"] = *record.Value
		}
		if key == "UdpFootprint" {
			listItem["udp_footprint"] = *record.Value
		}
		if key == "HasAbroad" {
			_ = d.Set("has_abroad", *record.Value)
			listItem["has_abroad"] = *record.Value
		}
		if key == "TcpSportStart" {
			listItem["tcp_start_port"] = *record.Value
		}
		if key == "TcpSportEnd" {
			listItem["tcp_end_port"] = *record.Value
		}
		if key == "UdpSportStart" {
			listItem["udp_start_port"] = *record.Value
		}
		if key == "UdpSportEnd" {
			listItem["udp_end_port"] = *record.Value
		}
		if key == "MaxUdpPackageLen" {
			listItem["max_udp_package_len"] = *record.Value
		}
		if key == "MinUdpPackageLen" {
			listItem["min_udp_package_len"] = *record.Value
		}
		if key == "MaxTcpPackageLen" {
			listItem["max_tcp_package_len"] = *record.Value
		}
		if key == "MinTcpPackageLen" {
			listItem["min_tcp_package_len"] = *record.Value
		}
		if key == "AppType" {
			listItem["app_type"] = *record.Value
		}
		if key == "AppProtocols" {
			listItem["app_protocols"] = strings.Split(*record.Value, ";")
		}
		if key == "WebApiUrl" {
			listItem["web_api_urls"] = strings.Split(*record.Value, ";")
		}
		if key == "PlatformTypes" {
			listItem["platform_types"] = strings.Split(*record.Value, ";")
		}
		if key == "Id" {
			listItem["scene_id"] = *record.Value
		}
		if key == "CreateTime" {
			listItem["create_time"] = *record.Value
		}
	}
	list = append(list, listItem)
	ids = append(ids, listItem["scene_id"].(string))

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
