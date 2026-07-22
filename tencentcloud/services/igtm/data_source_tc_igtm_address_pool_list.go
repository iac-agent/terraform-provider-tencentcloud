package igtm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	igtmv20231024 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/igtm/v20231024"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudIgtmAddressPoolList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudIgtmAddressPoolListRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Alert 过滤器 conditions。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤字段名称，支持 列表 作为 follows:\n- PoolName: 地址 池 名称\n- MonitorId: Monitor ID. 此 是 必填 参数，failure 到 provide 将 cause interface 查询 failure。",
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
							Description: "是否enable fuzzy 查询，仅 支持 过滤字段名称 作为 域名\nWhen fuzzy 查询 是 已启用，最大 值 长度 是 1，otherwise 最大 值 长度 是 5. (Reserved 字段，currently 不 使用)。",
						},
					},
				},
			},

			"address_pool_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Resource 组 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pool_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "地址 池 ID。",
						},
						"pool_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地址 池 名称",
						},
						"addr_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地址 池 地址 类型: IPV4，IPV6，DOMAIN。",
						},
						"traffic_strategy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Traffic strategy: WEIGHT load balancing，ALL resolve all。",
						},
						"monitor_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Monitor ID。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OK normal，DOWN failure，WARN risk，UNKNOWN unknown。",
						},
						"address_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "地址 count。",
						},
						"monitor_group_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Probe point count。",
						},
						"monitor_task_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Detection 任务 count。",
						},
						"instance_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "实例 related 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
								},
							},
						},
						"address_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "地址 池 地址 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"addr": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "地址 值: 仅 支持 IPv4，IPv6 和 域名 名称 formats;\nLoopback addresses，reserved addresses，内部 网络 addresses 和 Tencent reserved 网络 segments 是 不 支持。",
									},
									"is_enable": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "是否enable: DISABLED 已禁用; ENABLED 已启用",
									},
									"address_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "地址 ID。",
									},
									"location": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "地址 名称",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "OK normal，DOWN failure，WARN risk，UNKNOWN detecting，UNMONITORED unknown。",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "权重，必填 当 流量 strategy 是 WEIGHT; 范围 1-100。",
									},
									"created_on": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创建时间。",
									},
									"updated_on": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "修改时间。",
									},
								},
							},
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "更新时间。",
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

func dataSourceTencentCloudIgtmAddressPoolListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_igtm_address_pool_list.read")()
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

	var respData []*igtmv20231024.AddressPool
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeIgtmAddressPoolListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	addressPoolSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, addressPoolSet := range respData {
			addressPoolSetMap := map[string]interface{}{}
			if addressPoolSet.PoolId != nil {
				addressPoolSetMap["pool_id"] = addressPoolSet.PoolId
			}

			if addressPoolSet.PoolName != nil {
				addressPoolSetMap["pool_name"] = addressPoolSet.PoolName
			}

			if addressPoolSet.AddrType != nil {
				addressPoolSetMap["addr_type"] = addressPoolSet.AddrType
			}

			if addressPoolSet.TrafficStrategy != nil {
				addressPoolSetMap["traffic_strategy"] = addressPoolSet.TrafficStrategy
			}

			if addressPoolSet.MonitorId != nil {
				addressPoolSetMap["monitor_id"] = addressPoolSet.MonitorId
			}

			if addressPoolSet.Status != nil {
				addressPoolSetMap["status"] = addressPoolSet.Status
			}

			if addressPoolSet.AddressNum != nil {
				addressPoolSetMap["address_num"] = addressPoolSet.AddressNum
			}

			if addressPoolSet.MonitorGroupNum != nil {
				addressPoolSetMap["monitor_group_num"] = addressPoolSet.MonitorGroupNum
			}

			if addressPoolSet.MonitorTaskNum != nil {
				addressPoolSetMap["monitor_task_num"] = addressPoolSet.MonitorTaskNum
			}

			instanceInfoList := make([]map[string]interface{}, 0, len(addressPoolSet.InstanceInfo))
			if addressPoolSet.InstanceInfo != nil {
				for _, instanceInfo := range addressPoolSet.InstanceInfo {
					instanceInfoMap := map[string]interface{}{}
					if instanceInfo.InstanceId != nil {
						instanceInfoMap["instance_id"] = instanceInfo.InstanceId
					}

					if instanceInfo.InstanceName != nil {
						instanceInfoMap["instance_name"] = instanceInfo.InstanceName
					}

					instanceInfoList = append(instanceInfoList, instanceInfoMap)
				}

				addressPoolSetMap["instance_info"] = instanceInfoList
			}

			addressSetList := make([]map[string]interface{}, 0, len(addressPoolSet.AddressSet))
			if addressPoolSet.AddressSet != nil {
				for _, addressSet := range addressPoolSet.AddressSet {
					addressSetMap := map[string]interface{}{}
					if addressSet.Addr != nil {
						addressSetMap["addr"] = addressSet.Addr
					}

					if addressSet.IsEnable != nil {
						addressSetMap["is_enable"] = addressSet.IsEnable
					}

					if addressSet.AddressId != nil {
						addressSetMap["address_id"] = addressSet.AddressId
					}

					if addressSet.Location != nil {
						addressSetMap["location"] = addressSet.Location
					}

					if addressSet.Status != nil {
						addressSetMap["status"] = addressSet.Status
					}

					if addressSet.Weight != nil {
						addressSetMap["weight"] = addressSet.Weight
					}

					if addressSet.CreatedOn != nil {
						addressSetMap["created_on"] = addressSet.CreatedOn
					}

					if addressSet.UpdatedOn != nil {
						addressSetMap["updated_on"] = addressSet.UpdatedOn
					}

					addressSetList = append(addressSetList, addressSetMap)
				}

				addressPoolSetMap["address_set"] = addressSetList
			}

			if addressPoolSet.CreatedOn != nil {
				addressPoolSetMap["created_on"] = addressPoolSet.CreatedOn
			}

			if addressPoolSet.UpdatedOn != nil {
				addressPoolSetMap["updated_on"] = addressPoolSet.UpdatedOn
			}

			addressPoolSetList = append(addressPoolSetList, addressPoolSetMap)
		}

		_ = d.Set("address_pool_set", addressPoolSetList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), addressPoolSetList); e != nil {
			return e
		}
	}

	return nil
}
