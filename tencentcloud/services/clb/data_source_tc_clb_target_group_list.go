package clb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClbTargetGroupList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbTargetGroupListRead,
		Schema: map[string]*schema.Schema{
			"target_group_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "目标组 ID 数组。",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器数组，不包含 TargetGroupIds。有效值：TargetGroupVpcId 和 TargetGroupName。将首先使用目标组 ID。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤器名称。",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "过滤值数组。",
						},
					},
				},
			},

			"target_group_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "显示的目标群体的信息集。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "目标组 ID。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "目标组的 vpcid。",
						},
						"target_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "目标群体名称。",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "目标组的默认端口。注意：该字段可能返回null，表示取不到有效值。",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "目标组创建时间。",
						},
						"updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "目标群体修改时间。",
						},
						"associated_rule": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "关联规则数组。注意：该字段可能返回null，表示取不到有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"load_balancer_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "关联的CLB实例ID。",
									},
									"listener_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "关联监听器的ID。",
									},
									"location_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "关联的转发规则ID。注意：该字段可能返回null，表示取不到有效值。",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "关联侦听器的协议类型，例如 HTTP 或 TCP。",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "关联监听器的端口。",
									},
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "关联转发规则的域名。注意：该字段可能返回null，表示取不到有效值。",
									},
									"url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "关联的转发规则的URL。注意：该字段可能返回null，表示取不到有效值。",
									},
									"load_balancer_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CLB 实例名称。",
									},
									"listener_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "听众姓名。",
									},
								},
							},
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

func dataSourceTencentCloudClbTargetGroupListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clb_target_group_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("target_group_ids"); ok {
		targetGroupIdsSet := v.(*schema.Set).List()
		paramMap["TargetGroupIds"] = helper.InterfacesStringsPoint(targetGroupIdsSet)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*clb.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := clb.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["Filters"] = tmpSet
	}

	service := ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var targetGroupSet []*clb.TargetGroupInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClbTargetGroupListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		targetGroupSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(targetGroupSet))
	tmpList := make([]map[string]interface{}, 0, len(targetGroupSet))

	if targetGroupSet != nil {
		for _, targetGroupInfo := range targetGroupSet {
			targetGroupInfoMap := map[string]interface{}{}

			if targetGroupInfo.TargetGroupId != nil {
				targetGroupInfoMap["target_group_id"] = targetGroupInfo.TargetGroupId
			}

			if targetGroupInfo.VpcId != nil {
				targetGroupInfoMap["vpc_id"] = targetGroupInfo.VpcId
			}

			if targetGroupInfo.TargetGroupName != nil {
				targetGroupInfoMap["target_group_name"] = targetGroupInfo.TargetGroupName
			}

			if targetGroupInfo.Port != nil {
				targetGroupInfoMap["port"] = targetGroupInfo.Port
			}

			if targetGroupInfo.CreatedTime != nil {
				targetGroupInfoMap["created_time"] = targetGroupInfo.CreatedTime
			}

			if targetGroupInfo.UpdatedTime != nil {
				targetGroupInfoMap["updated_time"] = targetGroupInfo.UpdatedTime
			}

			if targetGroupInfo.AssociatedRule != nil {
				associatedRuleList := []interface{}{}
				for _, associatedRule := range targetGroupInfo.AssociatedRule {
					associatedRuleMap := map[string]interface{}{}

					if associatedRule.LoadBalancerId != nil {
						associatedRuleMap["load_balancer_id"] = associatedRule.LoadBalancerId
					}

					if associatedRule.ListenerId != nil {
						associatedRuleMap["listener_id"] = associatedRule.ListenerId
					}

					if associatedRule.LocationId != nil {
						associatedRuleMap["location_id"] = associatedRule.LocationId
					}

					if associatedRule.Protocol != nil {
						associatedRuleMap["protocol"] = associatedRule.Protocol
					}

					if associatedRule.Port != nil {
						associatedRuleMap["port"] = associatedRule.Port
					}

					if associatedRule.Domain != nil {
						associatedRuleMap["domain"] = associatedRule.Domain
					}

					if associatedRule.Url != nil {
						associatedRuleMap["url"] = associatedRule.Url
					}

					if associatedRule.LoadBalancerName != nil {
						associatedRuleMap["load_balancer_name"] = associatedRule.LoadBalancerName
					}

					if associatedRule.ListenerName != nil {
						associatedRuleMap["listener_name"] = associatedRule.ListenerName
					}

					associatedRuleList = append(associatedRuleList, associatedRuleMap)
				}

				targetGroupInfoMap["associated_rule"] = []interface{}{associatedRuleList}
			}

			ids = append(ids, *targetGroupInfo.TargetGroupId)
			tmpList = append(tmpList, targetGroupInfoMap)
		}

		_ = d.Set("target_group_set", tmpList)
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
