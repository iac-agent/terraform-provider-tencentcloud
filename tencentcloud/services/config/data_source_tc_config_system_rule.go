package config

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func DataSourceTencentCloudConfigSystemRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudConfigSystemRuleRead,
		Schema: map[string]*schema.Schema{
			"identifier": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule unique identifier.",
			},

			"system_config_rule": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "System preset rule detail.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule unique identifier.",
						},
						"rule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule description.",
						},
						"risk_level": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Risk level. Valid values: 1 (high risk), 2 (medium risk), 3 (low risk).",
						},
						"service_function": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Corresponding service function.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last update time.",
						},
						"trigger_type": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Trigger type list.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"resource_type": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Supported resource type list.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"label": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Rule label list.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"reference_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of times this rule is referenced.",
						},
						"identifier_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule type.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudConfigSystemRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_config_system_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(nil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = ConfigService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		identifier = d.Get("identifier").(string)
	)

	respData, reqErr := service.DescribeConfigSystemRuleById(ctx, identifier)
	if reqErr != nil {
		return reqErr
	}

	if respData == nil {
		return fmt.Errorf("config system rule %s not found", identifier)
	}

	ruleList := flattenSystemConfigRuleList([]*configv20220802.SystemConfigRule{respData})
	_ = d.Set("system_config_rule", ruleList)

	d.SetId(identifier)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
