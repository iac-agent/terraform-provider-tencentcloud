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

func DataSourceTencentCloudDayuCCHttpsPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDayuCCHttpsPoliciesRead,
		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_RESOURCE_TYPE),
				Description:  "类型 resource that the CC https policy works for，valid 值 is `bgpip`。",
			},
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the resource that the CC https policy works for。",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the CC https policy to be queried。",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 20),
				Description:  "名称 CC https policy to be queried。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 CC https policies. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID resource that the CC self-define https policy works for。",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 resource that the CC self-define https policy works for。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 CC self-define https policy。",
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "操作 模式",
						},
						"switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicate the CC self-define https policy takes effect or not。",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 that the CC self-define https policy works for。",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule ID 域名 that the CC self-define https policy works for。",
						},
						"rule_list": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"skey": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "键 of the rule。",
									},
									"operator": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "操作者 of the rule。",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule 值",
									},
								},
							},
							Description: "Rule 列表 the CC self-define https policy。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 of the CC self-define https policy。",
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the CC self-define https policy。",
						},
						"ip_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Ip of the CC self-define https policy。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDayuCCHttpsPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dayu_cc_https_policies.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DayuService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	resourceType := d.Get("resource_type").(string)
	resourceId := d.Get("resource_id").(string)
	policyId := d.Get("policy_id").(string)
	name := d.Get("name").(string)

	ccPolicies := make([]*dayu.CCPolicy, 0)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, _, err := service.DescribeCCSelfdefinePolicies(ctx, resourceType, resourceId, name, policyId)
		if err != nil {
			return tccommon.RetryError(err)
		}
		ccPolicies = result
		return nil
	})
	if err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(ccPolicies))
	ids := make([]string, 0, len(ccPolicies))

	listItem := make(map[string]interface{})
	for _, policy := range ccPolicies {
		listItem["name"] = *policy.Name
		listItem["create_time"] = *policy.CreateTime
		listItem["policy_id"] = *policy.SetId
		listItem["switch"] = *policy.Switch > 0
		listItem["ip_list"] = helper.StringsInterfaces(policy.IpList)
		listItem["action"] = *policy.ExeMode
		listItem["rule_list"] = flattenCCRuleList(policy.RuleList)
		listItem["rule_id"] = *policy.RuleId
		listItem["domain"] = *policy.Domain

		list = append(list, listItem)
		ids = append(ids, *policy.SetId)
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
