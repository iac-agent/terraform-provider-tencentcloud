package igtm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	igtmv20231024 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/igtm/v20231024"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudIgtmInstanceList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudIgtmInstanceListRead,
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
							Description: "过滤字段名称，支持 列表 作为 follows:\n- 实例 ID: IGTM 实例 ID.\n- 域名: IGTM 实例 域名\n- MonitorId: Monitor ID.\n- PoolId: Pool ID. 此 是 必填 参数，不 passing 它 将 cause interface 查询 failure。",
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

			"instance_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "实例 列表。",
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
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "资源 ID",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Business 域名",
						},
						"access_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cname 域名 访问 方法\nCUSTOM: Custom 访问 域名\nSYSTEM: System 访问 域名",
						},
						"access_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access 域名",
						},
						"access_sub_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access subdomain。",
						},
						"global_ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Global 记录 过期时间。",
						},
						"package_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Package 类型\nFREE: Free 版本\nSTANDARD: Standard 版本\nULTIMATE: Ultimate 版本",
						},
						"working_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 running 状态\nNORMAL: Healthy\nFAULTY: At risk\nDOWN: Down\nUNKNOWN: Unknown。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例状态，ENABLED: Normal，DISABLED: 已禁用",
						},
						"is_cname_configured": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether cname 访问: true accessed; false 不 accessed。",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "备注",
						},
						"strategy_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Strategy count。",
						},
						"address_pool_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Bound 地址 池 count。",
						},
						"monitor_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Bound 监控 count。",
						},
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
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 创建时间。",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 更新时间。",
						},
					},
				},
			},

			"system_access_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether 系统 域名 访问 是 支持: true 支持; false 不 支持。",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudIgtmInstanceListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_igtm_instance_list.read")()
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

	var (
		respData            []*igtmv20231024.Instance
		systemAccessEnabled *bool
	)
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, saEnabled, e := service.DescribeIgtmInstanceListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		systemAccessEnabled = saEnabled
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	instanceSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, instanceSet := range respData {
			instanceSetMap := map[string]interface{}{}
			if instanceSet.InstanceId != nil {
				instanceSetMap["instance_id"] = instanceSet.InstanceId
			}

			if instanceSet.InstanceName != nil {
				instanceSetMap["instance_name"] = instanceSet.InstanceName
			}

			if instanceSet.ResourceId != nil {
				instanceSetMap["resource_id"] = instanceSet.ResourceId
			}

			if instanceSet.Domain != nil {
				instanceSetMap["domain"] = instanceSet.Domain
			}

			if instanceSet.AccessType != nil {
				instanceSetMap["access_type"] = instanceSet.AccessType
			}

			if instanceSet.AccessDomain != nil {
				instanceSetMap["access_domain"] = instanceSet.AccessDomain
			}

			if instanceSet.AccessSubDomain != nil {
				instanceSetMap["access_sub_domain"] = instanceSet.AccessSubDomain
			}

			if instanceSet.GlobalTtl != nil {
				instanceSetMap["global_ttl"] = instanceSet.GlobalTtl
			}

			if instanceSet.PackageType != nil {
				instanceSetMap["package_type"] = instanceSet.PackageType
			}

			if instanceSet.WorkingStatus != nil {
				instanceSetMap["working_status"] = instanceSet.WorkingStatus
			}

			if instanceSet.Status != nil {
				instanceSetMap["status"] = instanceSet.Status
			}

			if instanceSet.IsCnameConfigured != nil {
				instanceSetMap["is_cname_configured"] = instanceSet.IsCnameConfigured
			}

			if instanceSet.Remark != nil {
				instanceSetMap["remark"] = instanceSet.Remark
			}

			if instanceSet.StrategyNum != nil {
				instanceSetMap["strategy_num"] = instanceSet.StrategyNum
			}

			if instanceSet.AddressPoolNum != nil {
				instanceSetMap["address_pool_num"] = instanceSet.AddressPoolNum
			}

			if instanceSet.MonitorNum != nil {
				instanceSetMap["monitor_num"] = instanceSet.MonitorNum
			}

			if instanceSet.PoolId != nil {
				instanceSetMap["pool_id"] = instanceSet.PoolId
			}

			if instanceSet.PoolName != nil {
				instanceSetMap["pool_name"] = instanceSet.PoolName
			}

			if instanceSet.CreatedOn != nil {
				instanceSetMap["created_on"] = instanceSet.CreatedOn
			}

			if instanceSet.UpdatedOn != nil {
				instanceSetMap["updated_on"] = instanceSet.UpdatedOn
			}

			instanceSetList = append(instanceSetList, instanceSetMap)
		}

		_ = d.Set("instance_set", instanceSetList)
	}

	if systemAccessEnabled != nil {
		_ = d.Set("system_access_enabled", systemAccessEnabled)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
