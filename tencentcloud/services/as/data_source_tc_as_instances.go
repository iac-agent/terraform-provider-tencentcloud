package as

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudAsInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAsInstancesRead,
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "实例 ID 云 服务器 (CVM) 到 是 queried. 限制 是 100 per 请求。",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 conditions. 如果 there 是 多个 Filters， relationship between Filters 是 logical AND (AND) relationship. 如果 there 是 多个 Values 在 same 过滤器， relationship between Values under same 过滤器 是 logical OR (OR) relationship。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Fields 到 是 filtered. 有效 names: `实例-ID`: Filters 通过 实例 ID，`auto-scaling-组-ID`: 过滤器 通过 scaling 组 ID",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "值 的 字段。",
						},
					},
				},
			},

			"instance_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 实例 details。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID",
						},
						"auto_scaling_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Auto scaling 组 ID",
						},
						"auto_scaling_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Auto scaling 组名称",
						},
						"launch_configuration_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "启动配置 ID",
						},
						"launch_configuration_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Launch 配置 名称",
						},
						"life_cycle_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Life cycle state. Please refer 到 link 对于 字段 值 details: https://云.tencent.com/document/api/377/20453#实例。",
						},
						"health_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Health 状态， 有效 值 是 HEALTHY 和 UNHEALTHY。",
						},
						"protected_from_scale_in": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable scale 在 protection。",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Available 可用区",
						},
						"creation_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "有效值：`AUTO_CREATION`，`MANUAL_ATTACHING`。",
						},
						"add_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "时间 当 实例 joined 组。",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例类型",
						},
						"version_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "版本 ID。",
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

func dataSourceTencentCloudAsInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_as_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		paramMap["InstanceIds"] = helper.InterfacesStringsPoint(instanceIdsSet)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*as.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := as.Filter{}
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
		paramMap["filters"] = tmpSet
	}

	service := AsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var instanceList []*as.Instance

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAsInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instanceList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceList))
	tmpList := make([]map[string]interface{}, 0, len(instanceList))

	if instanceList != nil {
		for _, instance := range instanceList {
			instanceMap := map[string]interface{}{}

			if instance.InstanceId != nil {
				instanceMap["instance_id"] = instance.InstanceId
			}

			if instance.AutoScalingGroupId != nil {
				instanceMap["auto_scaling_group_id"] = instance.AutoScalingGroupId
			}

			if instance.AutoScalingGroupName != nil {
				instanceMap["auto_scaling_group_name"] = instance.AutoScalingGroupName
			}

			if instance.LaunchConfigurationId != nil {
				instanceMap["launch_configuration_id"] = instance.LaunchConfigurationId
			}

			if instance.LaunchConfigurationName != nil {
				instanceMap["launch_configuration_name"] = instance.LaunchConfigurationName
			}

			if instance.LifeCycleState != nil {
				instanceMap["life_cycle_state"] = instance.LifeCycleState
			}

			if instance.HealthStatus != nil {
				instanceMap["health_status"] = instance.HealthStatus
			}

			if instance.ProtectedFromScaleIn != nil {
				instanceMap["protected_from_scale_in"] = instance.ProtectedFromScaleIn
			}

			if instance.Zone != nil {
				instanceMap["zone"] = instance.Zone
			}

			if instance.CreationType != nil {
				instanceMap["creation_type"] = instance.CreationType
			}

			if instance.AddTime != nil {
				instanceMap["add_time"] = instance.AddTime
			}

			if instance.InstanceType != nil {
				instanceMap["instance_type"] = instance.InstanceType
			}

			if instance.VersionNumber != nil {
				instanceMap["version_number"] = instance.VersionNumber
			}

			ids = append(ids, *instance.InstanceId)
			tmpList = append(tmpList, instanceMap)
		}

		_ = d.Set("instance_list", tmpList)
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
