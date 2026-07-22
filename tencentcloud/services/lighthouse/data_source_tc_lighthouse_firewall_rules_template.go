package lighthouse

import (
	"context"
	"encoding/json"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudLighthouseFirewallRulesTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseFirewallRulesTemplateRead,
		Schema: map[string]*schema.Schema{
			"firewall_rule_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Firewall 规则 details 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application 类型 有效 值 是 自定义，HTTP (80)，HTTPS (443)，Linux login (22)，Windows login (3389)，MySQL (3306)，SQL Server (1433)，all TCP ports，all UDP ports，Ping-ICMP，ALL。",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "协议 有效 值 是 TCP，UDP，ICMP，ALL。",
						},
						"port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "端口 有效 值 是 ALL，一个 单个 端口，多个 ports separated 通过 commas，或 端口 范围 indicated 通过 minus sign。",
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP 范围 或 IP (mutually exclusive). 默认值为 0.0.0.0/0，其中 表示all sources。",
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "有效 值 是 (ACCEPT，DROP). 默认值为 ACCEPT。",
						},
						"firewall_rule_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Firewall 规则 描述",
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

func dataSourceTencentCloudLighthouseFirewallRulesTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_lighthouse_firewall_rules_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var firewallRuleSet []*lighthouse.FirewallRuleInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseFirewallRulesTemplateByFilter(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}
		firewallRuleSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(firewallRuleSet))
	tmpList := make([]map[string]interface{}, 0, len(firewallRuleSet))

	if firewallRuleSet != nil {
		for _, firewallRuleInfo := range firewallRuleSet {
			firewallRuleInfoMap := map[string]interface{}{}

			if firewallRuleInfo.AppType != nil {
				firewallRuleInfoMap["app_type"] = firewallRuleInfo.AppType
			}

			if firewallRuleInfo.Protocol != nil {
				firewallRuleInfoMap["protocol"] = firewallRuleInfo.Protocol
			}

			if firewallRuleInfo.Port != nil {
				firewallRuleInfoMap["port"] = firewallRuleInfo.Port
			}

			if firewallRuleInfo.CidrBlock != nil {
				firewallRuleInfoMap["cidr_block"] = firewallRuleInfo.CidrBlock
			}

			if firewallRuleInfo.Action != nil {
				firewallRuleInfoMap["action"] = firewallRuleInfo.Action
			}

			if firewallRuleInfo.FirewallRuleDescription != nil {
				firewallRuleInfoMap["firewall_rule_description"] = firewallRuleInfo.FirewallRuleDescription
			}
			firewallRuleInfoJson, err := json.Marshal(*firewallRuleInfo)
			if err != nil {
				return err
			}
			ids = append(ids, string(firewallRuleInfoJson))
			tmpList = append(tmpList, firewallRuleInfoMap)
		}

		_ = d.Set("firewall_rule_set", tmpList)
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
