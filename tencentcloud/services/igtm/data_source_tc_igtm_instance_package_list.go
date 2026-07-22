package igtm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	igtmv20231024 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/igtm/v20231024"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudIgtmInstancePackageList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudIgtmInstancePackageListRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤字段名称，supported list as follows:\n- 实例 ID: instance ID.\n- InstanceName: 实例名称\n- ResourceId: package ID.\n- PackageType: package 类型 This is a 必填 parameter，not passing it will cause interface query failure。",
						},
						"value": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "过滤字段值",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"fuzzy": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "是否enable fuzzy query，only supports 过滤字段名称 as 域名\nWhen fuzzy query is 已启用，maximum 值 length is 1，otherwise maximum 值 length is 5. (Reserved field，not currently used)。",
						},
					},
				},
			},

			"is_used": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether used: 0 not used 1 used。",
			},

			"instance_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance package list。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance package resource ID。",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例名称",
						},
						"package_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Package 类型\nFREE: Free 版本\nSTANDARD: Standard 版本\nULTIMATE: Ultimate 版本",
						},
						"current_deadline": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Package 过期时间。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Package 创建时间。",
						},
						"is_expire": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether expired 0 no 1 yes。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例状态\nENABLED: Normal\nDISABLED: 已禁用",
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether auto-renew 0 no 1 yes。",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "备注",
						},
						"cost_item_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Billing item。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cost_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Billing item 名称",
									},
									"cost_value": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Billing item 值",
									},
								},
							},
						},
						"min_check_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum check interval time s。",
						},
						"min_global_ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum TTL s。",
						},
						"traffic_strategy": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Traffic strategy 类型: ALL return all，WEIGHT 权重",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"schedule_strategy": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Strategy 类型: LOCATION schedule by geographic location，DELAY schedule by 延迟",
							Elem: &schema.Schema{
								Type: schema.TypeString,
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

func dataSourceTencentCloudIgtmInstancePackageListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_igtm_instance_package_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = IgtmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*igtmv20231024.ResourceFilter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			resourceFilter := igtmv20231024.ResourceFilter{}
			if v, ok := filtersMap["name"].(string); ok && v != "" {
				resourceFilter.Name = helper.String(v)
			}

			if v, ok := filtersMap["value"]; ok {
				valueSet := v.(*schema.Set).List()
				for i := range valueSet {
					value := valueSet[i].(string)
					resourceFilter.Value = append(resourceFilter.Value, helper.String(value))
				}
			}

			if v, ok := filtersMap["fuzzy"].(bool); ok {
				resourceFilter.Fuzzy = helper.Bool(v)
			}

			tmpSet = append(tmpSet, &resourceFilter)
		}

		paramMap["Filters"] = tmpSet
	}

	if v, ok := d.GetOkExists("is_used"); ok {
		paramMap["IsUsed"] = helper.IntUint64(v.(int))
	}

	var respData []*igtmv20231024.InstancePackage
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeIgtmInstancePackageListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	instanceSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, instanceSet := range respData {
			instanceSetMap := map[string]interface{}{}
			if instanceSet.ResourceId != nil {
				instanceSetMap["resource_id"] = instanceSet.ResourceId
			}

			if instanceSet.InstanceId != nil {
				instanceSetMap["instance_id"] = instanceSet.InstanceId
			}

			if instanceSet.InstanceName != nil {
				instanceSetMap["instance_name"] = instanceSet.InstanceName
			}

			if instanceSet.PackageType != nil {
				instanceSetMap["package_type"] = instanceSet.PackageType
			}

			if instanceSet.CurrentDeadline != nil {
				instanceSetMap["current_deadline"] = instanceSet.CurrentDeadline
			}

			if instanceSet.CreateTime != nil {
				instanceSetMap["create_time"] = instanceSet.CreateTime
			}

			if instanceSet.IsExpire != nil {
				instanceSetMap["is_expire"] = instanceSet.IsExpire
			}

			if instanceSet.Status != nil {
				instanceSetMap["status"] = instanceSet.Status
			}

			if instanceSet.AutoRenewFlag != nil {
				instanceSetMap["auto_renew_flag"] = instanceSet.AutoRenewFlag
			}

			if instanceSet.Remark != nil {
				instanceSetMap["remark"] = instanceSet.Remark
			}

			costItemListList := make([]map[string]interface{}, 0, len(instanceSet.CostItemList))
			if instanceSet.CostItemList != nil {
				for _, costItemList := range instanceSet.CostItemList {
					costItemListMap := map[string]interface{}{}
					if costItemList.CostName != nil {
						costItemListMap["cost_name"] = costItemList.CostName
					}

					if costItemList.CostValue != nil {
						costItemListMap["cost_value"] = costItemList.CostValue
					}

					costItemListList = append(costItemListList, costItemListMap)
				}

				instanceSetMap["cost_item_list"] = costItemListList
			}

			if instanceSet.MinCheckInterval != nil {
				instanceSetMap["min_check_interval"] = instanceSet.MinCheckInterval
			}

			if instanceSet.MinGlobalTtl != nil {
				instanceSetMap["min_global_ttl"] = instanceSet.MinGlobalTtl
			}

			if instanceSet.TrafficStrategy != nil {
				instanceSetMap["traffic_strategy"] = instanceSet.TrafficStrategy
			}

			if instanceSet.ScheduleStrategy != nil {
				instanceSetMap["schedule_strategy"] = instanceSet.ScheduleStrategy
			}

			instanceSetList = append(instanceSetList, instanceSetMap)
		}

		_ = d.Set("instance_set", instanceSetList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), instanceSetList); e != nil {
			return e
		}
	}

	return nil
}
