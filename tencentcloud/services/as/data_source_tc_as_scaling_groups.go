package as

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudAsScalingGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAsScalingGroupRead,

		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A specified scaling 组 ID 用于query。",
			},
			"configuration_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter results by 启动配置 ID",
			},
			"scaling_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A scaling 组名称 用于query。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 用于query。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// computed
			"scaling_group_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 scaling group. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scaling_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Auto scaling 组 ID",
						},
						"scaling_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Auto scaling 组名称",
						},
						"configuration_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "启动配置 ID",
						},
						"max_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The 最大CVM instances。",
						},
						"min_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The 最小CVM instances。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID vpc with which the instance is associated。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID project to which the scaling group belongs. 默认值为 0。",
						},
						"subnet_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A 列表 subnet IDs。",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"zones": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A 列表 available zones。",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"default_cooldown": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Default cooldown time of scaling group。",
						},
						"desired_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The desired 数量 CVM instances。",
						},
						"load_balancer_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A 列表 traditional clb ids which the CVM instances attached to。",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"forward_load_balancers": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A 列表 application clb。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"load_balancer_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID available load balancers。",
									},
									"listener_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Listener ID for application load balancers。",
									},
									"location_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID forwarding rules。",
									},
									"target_attribute": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Attribute 列表 target rules。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "端口 number。",
												},
												"weight": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "权重",
												},
											},
										},
									},
								},
							},
						},
						"termination_policies": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A policy 用于select a CVM instance to be terminated from the scaling group。",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"retry_policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A retry policy can be used when a creation fails。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current 状态 a scaling group。",
						},
						"instance_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 instance。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time when the AS group was created。",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "标签 of the scaling group。",
						},
						"multi_zone_subnet_policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Multi 可用区 or subnet strategy，有效值：PRIORITY and EQUALITY。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAsScalingGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_as_scaling_groups.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	scalingGroupId := ""
	configurationId := ""
	scalingGroupName := ""
	if v, ok := d.GetOk("scaling_group_id"); ok {
		scalingGroupId = v.(string)
	}
	if v, ok := d.GetOk("configuration_id"); ok {
		configurationId = v.(string)
	}
	if v, ok := d.GetOk("scaling_group_name"); ok {
		scalingGroupName = v.(string)
	}

	tags := helper.GetTags(d, "tags")

	scalingGroups, err := asService.DescribeAutoScalingGroupByFilter(ctx, scalingGroupId, configurationId, scalingGroupName, tags)
	if err != nil {
		return err
	}

	scalingGroupList := make([]map[string]interface{}, 0, len(scalingGroups))
	for _, scalingGroup := range scalingGroups {
		tags := make(map[string]string, len(scalingGroup.Tags))
		for _, tag := range scalingGroup.Tags {
			tags[*tag.Key] = *tag.Value
		}

		mapping := map[string]interface{}{
			"scaling_group_id":         scalingGroup.AutoScalingGroupId,
			"scaling_group_name":       scalingGroup.AutoScalingGroupName,
			"configuration_id":         scalingGroup.LaunchConfigurationId,
			"status":                   scalingGroup.AutoScalingGroupStatus,
			"instance_count":           scalingGroup.InstanceCount,
			"max_size":                 scalingGroup.MaxSize,
			"min_size":                 scalingGroup.MinSize,
			"vpc_id":                   scalingGroup.VpcId,
			"subnet_ids":               helper.StringsInterfaces(scalingGroup.SubnetIdSet),
			"zones":                    helper.StringsInterfaces(scalingGroup.ZoneSet),
			"default_cooldown":         scalingGroup.DefaultCooldown,
			"desired_capacity":         scalingGroup.DesiredCapacity,
			"load_balancer_ids":        helper.StringsInterfaces(scalingGroup.LoadBalancerIdSet),
			"termination_policies":     helper.StringsInterfaces(scalingGroup.TerminationPolicySet),
			"retry_policy":             scalingGroup.RetryPolicy,
			"create_time":              scalingGroup.CreatedTime,
			"tags":                     tags,
			"multi_zone_subnet_policy": scalingGroup.MultiZoneSubnetPolicy,
		}
		if scalingGroup.ForwardLoadBalancerSet != nil && len(scalingGroup.ForwardLoadBalancerSet) > 0 {
			forwardLoadBalancers := make([]map[string]interface{}, 0, len(scalingGroup.ForwardLoadBalancerSet))
			for _, v := range scalingGroup.ForwardLoadBalancerSet {
				targetAttributes := make([]map[string]interface{}, 0, len(v.TargetAttributes))
				for _, vv := range v.TargetAttributes {
					targetAttribute := map[string]interface{}{
						"port":   vv.Port,
						"weight": vv.Weight,
					}
					targetAttributes = append(targetAttributes, targetAttribute)
				}
				forwardLoadBalancer := map[string]interface{}{
					"load_balancer_id": v.LoadBalancerId,
					"listener_id":      v.ListenerId,
					"target_attribute": targetAttributes,
					"location_id":      v.LocationId,
				}
				forwardLoadBalancers = append(forwardLoadBalancers, forwardLoadBalancer)
			}
			mapping["forward_load_balancers"] = forwardLoadBalancers
		}
		scalingGroupList = append(scalingGroupList, mapping)
	}

	d.SetId("ScalingGroupList" + scalingGroupId + scalingGroupName + configurationId)
	err = d.Set("scaling_group_list", scalingGroupList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set scaling group list fail, reason:%s\n ", logId, err.Error())
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), scalingGroupList); err != nil {
			return err
		}
	}

	return nil
}
