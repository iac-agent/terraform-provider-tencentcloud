package dayuv2

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcdayu "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dayu"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDayuL4RulesV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDayuL4RulesReadV2,
		Schema: map[string]*schema.Schema{
			"business": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(svcdayu.DAYU_RESOURCE_TYPE),
				Description:  "类型 资源 该 layer 4 规则 works 对于，有效 值 是 `bgpip`，`bgp`，`bgp-multip` 和 `net`。",
			},
			"virtual_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Virtual 端口 的 资源。",
			},
			"ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Ip 的 资源。",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 layer 4 规则. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "协议 的 规则。",
						},
						"source_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "来源 端口 的 layer 4 规则。",
						},
						"virtual_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "virtual 端口 的 layer 4 规则。",
						},
						"keeptime": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "keeptime 的 layer 4 规则。",
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
						"keep_enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "会话 hold switch。",
						},
						"source_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "来源 类型，`1` 对于 来源 的 主机，`2` 对于 来源 的 IP。",
						},
						"rule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 规则。",
						},
						"remove_switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Remove 水印 state。",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule 修改时间。",
						},
						"region": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Corresponding regional 信息。",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Bind 资源 IP 信息。",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Bind 资源 ID 信息。",
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

func dataSourceTencentCloudDayuL4RulesReadV2(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dayu_l4_rules_v2.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcdayu.NewDayuService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	business := d.Get("business").(string)
	extendParams := make(map[string]interface{})
	if v, ok := d.GetOk("ip"); ok {
		extendParams["ip"] = v.(string)
	}
	if v, ok := d.GetOk("virtual_port"); ok {
		extendParams["virtual_port"] = v.(int)
	}

	rules := make([]*dayu.NewL4RuleEntry, 0)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, err := service.DescribeNewL4Rules(ctx, business, extendParams)

		if err != nil {
			return tccommon.RetryError(err)
		}
		rules = result
		return nil
	})

	if err != nil {
		return err
	}
	list := make([]map[string]interface{}, 0)
	for _, rule := range rules {
		tmpRule := make(map[string]interface{})
		tmpRule["protocol"] = *rule.Protocol
		tmpRule["source_port"] = *rule.SourcePort
		tmpRule["virtual_port"] = *rule.VirtualPort
		tmpRule["keeptime"] = *rule.KeepEnable
		tmpSourceList := make([]map[string]interface{}, 0)
		for _, source := range rule.SourceList {
			tmpSource := make(map[string]interface{})
			tmpSource["source"] = *source.Source
			tmpSource["weight"] = *source.Weight
			tmpSourceList = append(tmpSourceList, tmpSource)
		}
		tmpRule["source_list"] = tmpSourceList
		tmpRule["rule_id"] = *rule.RuleId
		tmpRule["lb_type"] = *rule.LbType
		tmpRule["keep_enable"] = *rule.KeepEnable == 1
		tmpRule["source_type"] = *rule.SourceType
		tmpRule["rule_name"] = *rule.RuleName
		tmpRule["remove_switch"] = *rule.RemoveSwitch == 1
		tmpRule["modify_time"] = *rule.ModifyTime
		tmpRule["region"] = *rule.Region
		tmpRule["ip"] = *rule.Ip
		tmpRule["id"] = *rule.Id
		list = append(list, tmpRule)
	}
	ids := make([]string, 0, len(list))
	for _, listItem := range list {
		ids = append(ids, listItem["rule_id"].(string))
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("list", list); e != nil {
		log.Printf("[CRITAL]%s provider set rules fail, reason:%s\n", logId, e.Error())
		return e
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}
	return nil

}
