package igtm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	igtmv20231024 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/igtm/v20231024"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudIgtmDetectTaskPackageList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudIgtmDetectTaskPackageListRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Detect 任务 过滤器 conditions。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤字段名称，支持 列表 作为 follows:\n- ResourceId: detect 任务 资源 ID\n- PeriodStart: 最小 过期时间.\n- PeriodEnd: 最大 过期时间。",
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

			"task_package_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Detect 任务 包 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "资源 ID",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "资源类型\nTASK Detect 任务。",
						},
						"quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Quota。",
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
							Description: "状态\nENABLED: Normal\nISOLATED: Isolated\nDESTROYED: Destroyed\nREFUNDED: Refunded。",
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
						"group": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Detect 任务 类型: 100 系统 setting; 200 billing; 300 management 系统; 110D 监控 迁移 free 任务; 120 disaster recovery switch 任务。",
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

func dataSourceTencentCloudIgtmDetectTaskPackageListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_igtm_detect_task_package_list.read")()
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

	var respData []*igtmv20231024.DetectTaskPackage
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeIgtmDetectTaskPackageListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	taskPackageSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, taskPackageSet := range respData {
			taskPackageSetMap := map[string]interface{}{}
			if taskPackageSet.ResourceId != nil {
				taskPackageSetMap["resource_id"] = taskPackageSet.ResourceId
			}

			if taskPackageSet.ResourceType != nil {
				taskPackageSetMap["resource_type"] = taskPackageSet.ResourceType
			}

			if taskPackageSet.Quota != nil {
				taskPackageSetMap["quota"] = taskPackageSet.Quota
			}

			if taskPackageSet.CurrentDeadline != nil {
				taskPackageSetMap["current_deadline"] = taskPackageSet.CurrentDeadline
			}

			if taskPackageSet.CreateTime != nil {
				taskPackageSetMap["create_time"] = taskPackageSet.CreateTime
			}

			if taskPackageSet.IsExpire != nil {
				taskPackageSetMap["is_expire"] = taskPackageSet.IsExpire
			}

			if taskPackageSet.Status != nil {
				taskPackageSetMap["status"] = taskPackageSet.Status
			}

			if taskPackageSet.AutoRenewFlag != nil {
				taskPackageSetMap["auto_renew_flag"] = taskPackageSet.AutoRenewFlag
			}

			if taskPackageSet.Remark != nil {
				taskPackageSetMap["remark"] = taskPackageSet.Remark
			}

			costItemListList := make([]map[string]interface{}, 0, len(taskPackageSet.CostItemList))
			if taskPackageSet.CostItemList != nil {
				for _, costItemList := range taskPackageSet.CostItemList {
					costItemListMap := map[string]interface{}{}
					if costItemList.CostName != nil {
						costItemListMap["cost_name"] = costItemList.CostName
					}

					if costItemList.CostValue != nil {
						costItemListMap["cost_value"] = costItemList.CostValue
					}

					costItemListList = append(costItemListList, costItemListMap)
				}

				taskPackageSetMap["cost_item_list"] = costItemListList
			}

			if taskPackageSet.Group != nil {
				taskPackageSetMap["group"] = taskPackageSet.Group
			}

			taskPackageSetList = append(taskPackageSetList, taskPackageSetMap)
		}

		_ = d.Set("task_package_set", taskPackageSetList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), taskPackageSetList); e != nil {
			return e
		}
	}

	return nil
}
