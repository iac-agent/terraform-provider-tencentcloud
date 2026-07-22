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
				Description:  "类型 resource that the layer 4 rule works for，valid values are `bgpip`，`bgp`，`bgp-multip` and `net`。",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "域名 of resource。",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "协议 of resource，值 range [`http`，`https`]。",
			},
			"ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Ip of the resource。",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Deprecated:  "It has been deprecated from version 1.81.21.",
				Description: "The page start 偏移量，默认为 `0`。",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Deprecated:  "It has been deprecated from version 1.81.21.",
				Description: "The 数量 pages，默认为 `10`。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 layer 4 rules. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keep_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Session hold time，（秒）。",
						},
						"lb_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Load balancing 模式，the 值 is [1 (weighted round-robin)]。",
						},
						"source_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Back-to-来源 IP or 域名 名称",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "权重 值，take 值 [0,100]。",
									},
								},
							},
							Description: "来源 列表 the rule。",
						},
						"keep_enable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Session keep switch，值 [0 (session keep closed)，1 (session keep open)]。",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 of resource。",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "协议 of resource，值 range [`http`，`https`]。",
						},
						"source_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Back-to-origin method，值 [1 (域名 名称 back-to-来源)，2 (IP back-to-来源)]。",
						},
						"https_to_http_enable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否enable the Https 协议 to use Http back-to-来源，take the 值 [0 (off)，1 (on)]，默认为 off。",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Rule 状态，值 [0 (rule configuration is successful)，1 (rule configuration is in effect)，2 (rule configuration fails)，3 (rule deletion is in effect)，5 (rule deletion fails)，6 (rule is waiting to be configured)，7 (rule pending deletion)，8 (rule pending configuration certificate)]。",
						},
						"cc_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CC protection 级别 of HTTPS 协议",
						},
						"cc_enable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CC protection 状态 HTTPS 协议，the 值 is [0 (off)，1 (on)]。",
						},
						"cc_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CC protection threshold of HTTPS 协议",
						},
						"region": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The area 代码",
						},
						"rule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule 描述",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "修改时间 of resource。",
						},
						"virtual_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Virtual 端口 of resource。",
						},
						"cc_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CC protection 状态，值 [0(off)，1(on)]。",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Ip of the resource。",
						},
						"ssl_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SSL ID resource。",
						},
						"cert_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The 来源 of the certificate。",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the resource。",
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
