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
				Description: "过滤器 conditions。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤字段名称，支持 列表 作为 follows:\n- 实例 ID: 实例 ID.\n- InstanceName: 实例名称\n- ResourceId: 包 ID.\n- PackageType: 包 类型 此 是 必填 参数，不 passing 它 将 cause interface 查询 failure。",
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
							Description: "是否enable fuzzy 查询，仅 支持 过滤字段名称 作为 域名\nWhen fuzzy 查询 是 已启用，最大 值 长度 是 1，otherwise 最大 值 长度 是 5. (Reserved 字段，不 currently 使用)。",
						},
					},
				},
			},

			"is_used": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether 使用: 0 不 使用 1 使用。",
			},

			"instance_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "实例 包 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 包 资源 ID。",
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
							Description: "Whether expired 0 无 1 yes。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例状态\nENABLED: Normal\nDISABLED: 已禁用",
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether auto-renew 0 无 1 yes。",
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
							Description: "Minimum check 间隔 时间 s。",
						},
						"min_global_ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum TTL s。",
						},
						"traffic_strategy": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Traffic strategy 类型: ALL 返回 all，WEIGHT 权重",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"schedule_strategy": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Strategy 类型: LOCATION 调度 通过 geographic location，DELAY 调度 通过 延迟",
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
