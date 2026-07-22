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

func DataSourceTencentCloudDayuL7Rules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDayuL7RulesRead,
		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_RESOURCE_TYPE_HTTPS),
				Description:  "类型 资源 该 layer 7 规则 works 对于，有效 值 是 `bgpip`。",
			},
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID 的 资源 该 layer 7 规则 works 对于。",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID 的 layer 7 规则 到 是 queried。",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "域名 的 layer 7 规则 到 是 queried。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 layer 7 规则. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 该 7 layer 规则 works 对于。",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "协议 的 规则。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 规则。",
						},
						"switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicate 规则 将 take effect 或 不。",
						},
						"source_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "来源 类型，1 对于 来源 的 主机，2 对于 来源 的 ip。",
						},
						"source_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type:        schema.TypeString,
								Description: "来源 ip 或 域名",
							},
							Description: "来源 列表 规则。",
						},
						"ssl_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SSL ID。",
						},
						"health_check_switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "表示是否health check 是 已启用",
						},
						"health_check_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Interval 时间 的 health check。",
						},
						"health_check_health_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Health 阈值 的 health check。",
						},
						"health_check_unhealth_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unhealthy 阈值 的 health check。",
						},
						"health_check_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "HTTP 状态 代码 `1` 表示 返回值 `1xx` 是 health. `2` 表示 返回值 `2xx` 是 health. `4` 表示 返回值 `3xx` 是 health. `8` 表示 返回值 `4xx` 是 health. `16` 表示 返回值 `5xx` 是 health. 如果 您 want 多个 返回 codes 到 indicate health，need 到 add corresponding 值。",
						},
						"health_check_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "路径 的 health check。",
						},
						"health_check_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Methods 的 health check。",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 的 7 layer 规则。",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "状态 规则. `0` 对于 create/modify success，`2` 对于 create/modify fail，`3` 对于 delete success，`5` 对于 waiting 到 是 创建/modified，`7` 对于 waiting 到 是 删除 和 `8` 对于 waiting 到 get SSL ID。",
						},
						"threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Threshold 的 规则。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDayuL7RulesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dayu_l7_rules.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DayuService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	resourceType := d.Get("resource_type").(string)
	resourceId := d.Get("resource_id").(string)
	ruleId := d.Get("rule_id").(string)
	domain := d.Get("domain").(string)
	protocol := ""

	rules := make([]*dayu.L7RuleEntry, 0)
	healths := make([]*dayu.L7RuleHealth, 0)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, hResult, _, err := service.DescribeL7Rules(ctx, resourceType, resourceId, domain, ruleId, protocol)
		if err != nil {
			return tccommon.RetryError(err)
		}
		rules = result
		healths = hResult
		return nil
	})
	if err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(rules))
	ids := make([]string, 0, len(rules))

	for k, rule := range rules {
		listItem := make(map[string]interface{})

		listItem["name"] = *rule.RuleName
		listItem["domain"] = *rule.Domain
		listItem["ssl_id"] = *rule.SSLId
		listItem["rule_id"] = *rule.RuleId
		listItem["protocol"] = *rule.Protocol
		listItem["source_type"] = int(*rule.SourceType)
		listItem["status"] = int(*rule.Status)
		listItem["threshold"] = int(*rule.CCThreshold)

		if *rule.Protocol == DAYU_L7_RULE_PROTOCOL_HTTPS {
			listItem["switch"] = *rule.CCEnable > 0
		} else {
			listItem["switch"] = *rule.CCStatus > 0
		}
		sourceList := make([]*string, 0, len(rule.SourceList))
		for _, v := range rule.SourceList {
			sourceList = append(sourceList, v.Source)
		}
		listItem["source_list"] = helper.StringsInterfaces(sourceList)

		if k < len(healths) {
			health := healths[k]

			if health.Enable != nil {
				listItem["health_check_switch"] = *health.Enable > 0
			}
			if health.Url != nil {
				listItem["health_check_path"] = *health.Url
			}
			if health.StatusCode != nil {
				listItem["health_check_code"] = int(*health.StatusCode)
			}
			if health.Interval != nil {
				listItem["health_check_interval"] = int(*health.Interval)
			}
			if health.KickNum != nil {
				listItem["health_check_unhealth_num"] = int(*health.KickNum)
			}
			if health.AliveNum != nil {
				listItem["health_check_health_num"] = int(*health.AliveNum)
			}
			if health.Method != nil {
				listItem["health_check_method"] = *health.Method
			}
		}
		list = append(list, listItem)
		ids = append(ids, listItem["rule_id"].(string))
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
