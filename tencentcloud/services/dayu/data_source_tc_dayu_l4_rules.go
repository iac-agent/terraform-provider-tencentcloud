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

func DataSourceTencentCloudDayuL4Rules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDayuL4RulesRead,
		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_RESOURCE_TYPE),
				Description:  "类型 资源 该 layer 4 规则 works 对于，有效 值 是 `bgpip`，`bgp`，`bgp-multip` 和 `net`。",
			},
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID 的 资源 该 layer 4 规则 works 对于。",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID 的 layer 4 规则 到 是 queried。",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 layer 4 规则 到 是 queried。",
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
						"source_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "来源 类型，`1` 对于 来源 的 主机，`2` 对于 来源 的 IP。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 规则。",
						},
						"s_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "来源 端口 的 layer 4 规则。",
						},
						"d_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "destination 端口 的 layer 4 规则。",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "协议 的 规则。",
						},
						"source_list": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "来源 IP 或 域名",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "权重 的 来源",
									},
								},
							},
							Description: "来源 列表 规则。",
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
						"health_check_timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "HTTP 状态 代码 `1` 表示 返回值 `1xx` 是 health. `2` 表示 返回值 `2xx` 是 health. `4` 表示 返回值 `3xx` 是 health. `8` 表示 返回值 `4xx` 是 health. `16` 表示 返回值 `5xx` 是 health. 如果 您 want 多个 返回 codes 到 indicate health，need 到 add corresponding 值。",
						},
						"session_switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicate 该 会话 将 keep 或 不。",
						},
						"session_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Session keep 时间，仅 有效 当 `session_switch` 是 true， 可用 值 ranges 从 1 到 300 和 单位 是 second。",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 4 layer 规则。",
						},
						"lb_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "LB 类型 规则，`1` 对于 权重 cycling 和 `2` 对于 IP hash。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDayuL4RulesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dayu_l4_rules.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DayuService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	resourceType := d.Get("resource_type").(string)
	resourceId := d.Get("resource_id").(string)
	ruleId := d.Get("rule_id").(string)
	name := d.Get("name").(string)

	rules := make([]*dayu.L4RuleEntry, 0)
	healths := make([]*dayu.L4RuleHealth, 0)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, hResult, _, err := service.DescribeL4Rules(ctx, resourceType, resourceId, name, ruleId)
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

	listItem := make(map[string]interface{})
	for k, rule := range rules {
		health := healths[k]
		listItem["name"] = *rule.RuleName
		listItem["protocol"] = *rule.Protocol
		listItem["s_port"] = int(*rule.SourcePort)
		listItem["d_port"] = int(*rule.VirtualPort)
		listItem["rule_id"] = *rule.RuleId
		listItem["lb_type"] = int(*rule.LbType)
		listItem["source_type"] = int(*rule.SourceType)
		listItem["session_time"] = int(*rule.KeepTime)
		listItem["session_switch"] = *rule.KeepEnable > 0
		listItem["source_list"] = flattenSourceList(rule.SourceList)
		if health.Enable != nil {
			listItem["health_check_switch"] = *health.Enable > 0
		}
		if health.TimeOut != nil {
			listItem["health_check_timeout"] = int(*health.TimeOut)
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
