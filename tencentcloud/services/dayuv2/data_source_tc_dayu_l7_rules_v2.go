package dayuv2

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcdayu "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dayu"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDayuL7RulesV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDayuL7RulesReadV2,
		Schema: map[string]*schema.Schema{
			"business": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(svcdayu.DAYU_RESOURCE_TYPE),
				Description:  "类型 资源 该 layer 4 规则 works 对于，有效 值 是 `bgpip`，`bgp`，`bgp-multip` 和 `net`。",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "域名 的 资源。",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "协议 的 资源，值 范围 [`http`，`https`]。",
			},
			"ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Ip 的 资源。",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Deprecated:  "It has been deprecated from version 1.81.21.",
				Description: "页面 start 偏移量，默认为 `0`。",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Deprecated:  "It has been deprecated from version 1.81.21.",
				Description: "数量 pages，默认为 `10`。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 layer 4 规则. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keep_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Session hold 时间，（秒）。",
						},
						"lb_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Load balancing 模式， 值 是 [1 (weighted round-robin)]。",
						},
						"source_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Back-到-来源 IP 或 域名 名称",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "权重 值，take 值 [0,100]。",
									},
								},
							},
							Description: "来源 列表 规则。",
						},
						"keep_enable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Session keep switch，值 [0 (会话 keep closed)，1 (会话 keep open)]。",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 的 资源。",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "协议 的 资源，值 范围 [`http`，`https`]。",
						},
						"source_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Back-到-源站 方法，值 [1 (域名 名称 back-到-来源)，2 (IP back-到-来源)]。",
						},
						"https_to_http_enable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否enable Https 协议 到 使用 Http back-到-来源，take 值 [0 (关闭)，1 (在)]，默认为 关闭。",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Rule 状态，值 [0 (规则 配置 是 successful)，1 (规则 配置 是 在 effect)，2 (规则 配置 fails)，3 (规则 deletion 是 在 effect)，5 (规则 deletion fails)，6 (规则 是 waiting 到 是 已配置)，7 (规则 pending deletion)，8 (规则 pending 配置 证书)]。",
						},
						"cc_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CC protection 级别 的 HTTPS 协议",
						},
						"cc_enable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CC protection 状态 HTTPS 协议， 值 是 [0 (关闭)，1 (在)]。",
						},
						"cc_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CC protection 阈值 的 HTTPS 协议",
						},
						"region": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "area 代码",
						},
						"rule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule 描述",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "修改时间 的 资源。",
						},
						"virtual_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Virtual 端口 的 资源。",
						},
						"cc_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CC protection 状态，值 [0(关闭)，1(在)]。",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Ip 的 资源。",
						},
						"ssl_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SSL ID 资源。",
						},
						"cert_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "来源 的 证书。",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 的 资源。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDayuL7RulesReadV2(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dayu_l4_rules_v2.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcdayu.NewDayuService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	business := d.Get("business").(string)
	domain := d.Get("domain").(string)
	protocol := d.Get("protocol").(string)
	ip := d.Get("ip").(string)

	extendParams := make(map[string]interface{})
	extendParams["domain"] = domain
	extendParams["protocol"] = protocol
	extendParams["ip"] = ip

	rules, _, err := service.DescribeL7RulesV2(ctx, business, extendParams)
	if err != nil {
		return err
	}

	resultList := make([]map[string]interface{}, 0)
	for _, rule := range rules {
		tmpResultItem := make(map[string]interface{})
		tmpResultItem["keep_time"] = *rule.KeepTime
		tmpResultItem["lb_type"] = *rule.LbType
		sourceList := make([]map[string]interface{}, 0)
		for _, source := range rule.SourceList {
			tmpSource := make(map[string]interface{})
			tmpSource["source"] = *source.Source
			tmpSource["weight"] = *source.Weight
			sourceList = append(sourceList, tmpSource)
		}
		tmpResultItem["source_list"] = sourceList
		tmpResultItem["keep_enable"] = *rule.KeepEnable
		tmpResultItem["domain"] = *rule.Domain
		tmpResultItem["protocol"] = *rule.Protocol
		tmpResultItem["source_type"] = *rule.SourceType
		tmpResultItem["https_to_http_enable"] = *rule.HttpsToHttpEnable
		tmpResultItem["status"] = *rule.Status
		tmpResultItem["cc_level"] = *rule.CCLevel
		tmpResultItem["cc_enable"] = *rule.CCEnable
		tmpResultItem["cc_threshold"] = *rule.CCThreshold
		tmpResultItem["region"] = *rule.Region
		tmpResultItem["rule_name"] = *rule.RuleName
		tmpResultItem["modify_time"] = *rule.ModifyTime
		tmpResultItem["virtual_port"] = *rule.VirtualPort
		tmpResultItem["cc_status"] = *rule.CCStatus
		tmpResultItem["ip"] = *rule.Ip
		tmpResultItem["cert_type"] = *rule.CertType
		tmpResultItem["id"] = *rule.Id
		resultList = append(resultList, tmpResultItem)
	}
	ids := make([]string, 0, len(resultList))
	for _, listItem := range resultList {
		ids = append(ids, listItem["id"].(string))
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("list", resultList); e != nil {
		log.Printf("[CRITAL]%s provider set rules fail, reason:%s\n", logId, e.Error())
		return e
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), resultList)
	}
	return nil
}
