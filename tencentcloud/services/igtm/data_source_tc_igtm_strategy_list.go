package igtm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	igtmv20231024 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/igtm/v20231024"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudIgtmStrategyList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudIgtmStrategyListRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "实例 ID",
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Strategy 过滤器 conditions: StrategyName: strategy 名称",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤字段名称，支持 列表 作为 follows:\n- 类型: main 资源类型，CDN.\n- instanceId: IGTM 实例 ID. 此 是 必填 参数，failure 到 pass 将 cause interface 查询 failure。",
						},
						"value": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "过滤器 字段 值。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"fuzzy": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "是否enable fuzzy 查询，仅 支持 过滤字段名称 作为 域名\nWhen fuzzy 查询 是 已启用，值 最大 长度 是 1，otherwise 值 最大 长度 是 5. (Reserved 字段，currently unused)。",
						},
					},
				},
			},

			"strategy_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Strategy 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Strategy 名称",
						},
						"source": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "地址 来源",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dns_line_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Resolution 请求来源 line ID。",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resolution 请求来源 line 名称",
									},
								},
							},
						},
						"strategy_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Strategy ID。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Health 状态: ok healthy，warn risk，down failure。",
						},
						"activate_main_pool_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Activated main 池 ID，null 表示 unknown。",
						},
						"activate_level": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Current activated 地址 池 级别，0 表示 fallback activated，null 表示 unknown。",
						},
						"active_pool_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current activated 地址 池 集合 类型: main main 池; fallback fallback 池。",
						},
						"active_traffic_strategy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current activated 地址 池 流量 strategy: all resolve all; 权重 load balancing。",
						},
						"monitor_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Monitor count。",
						},
						"is_enabled": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether 已启用: ENABLED 已启用; DISABLED 已禁用",
						},
						"keep_domain_records": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "是否retain lines: 已启用 retain，已禁用 不 retain，仅 retain 默认值 lines。",
						},
						"switch_pool_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheduling 模式: AUTO 默认值; PAUSE 仅 pause without switching。",
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

func dataSourceTencentCloudIgtmStrategyListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_igtm_strategy_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(nil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = IgtmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

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

	var respData []*igtmv20231024.Strategy
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeIgtmStrategyListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	strategySetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, strategySet := range respData {
			strategySetMap := map[string]interface{}{}
			if strategySet.InstanceId != nil {
				strategySetMap["instance_id"] = strategySet.InstanceId
			}

			if strategySet.Name != nil {
				strategySetMap["name"] = strategySet.Name
			}

			sourceList := make([]map[string]interface{}, 0, len(strategySet.Source))
			if strategySet.Source != nil {
				for _, source := range strategySet.Source {
					sourceMap := map[string]interface{}{}
					if source.DnsLineId != nil {
						sourceMap["dns_line_id"] = source.DnsLineId
					}

					if source.Name != nil {
						sourceMap["name"] = source.Name
					}

					sourceList = append(sourceList, sourceMap)
				}

				strategySetMap["source"] = sourceList
			}

			if strategySet.StrategyId != nil {
				strategySetMap["strategy_id"] = strategySet.StrategyId
			}

			if strategySet.Status != nil {
				strategySetMap["status"] = strategySet.Status
			}

			if strategySet.ActivateMainPoolId != nil {
				strategySetMap["activate_main_pool_id"] = strategySet.ActivateMainPoolId
			}

			if strategySet.ActivateLevel != nil {
				strategySetMap["activate_level"] = strategySet.ActivateLevel
			}

			if strategySet.ActivePoolType != nil {
				strategySetMap["active_pool_type"] = strategySet.ActivePoolType
			}

			if strategySet.ActiveTrafficStrategy != nil {
				strategySetMap["active_traffic_strategy"] = strategySet.ActiveTrafficStrategy
			}

			if strategySet.MonitorNum != nil {
				strategySetMap["monitor_num"] = strategySet.MonitorNum
			}

			if strategySet.IsEnabled != nil {
				strategySetMap["is_enabled"] = strategySet.IsEnabled
			}

			if strategySet.KeepDomainRecords != nil {
				strategySetMap["keep_domain_records"] = strategySet.KeepDomainRecords
			}

			if strategySet.SwitchPoolType != nil {
				strategySetMap["switch_pool_type"] = strategySet.SwitchPoolType
			}

			if strategySet.CreatedOn != nil {
				strategySetMap["created_on"] = strategySet.CreatedOn
			}

			if strategySet.UpdatedOn != nil {
				strategySetMap["updated_on"] = strategySet.UpdatedOn
			}

			strategySetList = append(strategySetList, strategySetMap)
		}

		_ = d.Set("strategy_set", strategySetList)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), strategySetList); e != nil {
			return e
		}
	}

	return nil
}
