package dlc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlcv20210125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDlcStandardEngineResourceGroupConfigInformation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcStandardEngineResourceGroupConfigInformationRead,
		Schema: map[string]*schema.Schema{
			"sort_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort Field。",
			},

			"sorting": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Ascending 或 descending。",
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "过滤器 conditions 是 可选，引擎-资源-组-ID 或 引擎-ID。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Attribute 名称 如果 there 是 多个 filters， relationship between filters 是 logical OR relationship。",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Attribute 值，如果 there 是 多个 Values 在 same 过滤器， relationship between Values under same 过滤器 是 logical OR relationship。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"standard_engine_resource_group_config_infos": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Standard 引擎 资源 组，配置 related 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_group_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Engine 资源 组 ID。",
						},
						"data_engine_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Engine ID。",
						},
						"static_config_pairs": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Static 参数 的 资源 组，其中 require restarting 资源 组 到 take effect。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_item": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Configuration items。",
									},
									"config_value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Configuration item 值",
									},
								},
							},
						},
						"dynamic_config_pairs": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Dynamic 参数 的 资源 组，effective 在 next 任务。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_item": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Configuration items。",
									},
									"config_value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Configuration item 值",
									},
								},
							},
						},
						"create_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "创建时间。",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Required:    true,
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

func dataSourceTencentCloudDlcStandardEngineResourceGroupConfigInformationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_standard_engine_resource_group_config_information.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("sort_by"); ok {
		paramMap["SortBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sorting"); ok {
		paramMap["Sorting"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*dlcv20210125.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := dlcv20210125.Filter{}
			if v, ok := filtersMap["name"].(string); ok && v != "" {
				filter.Name = helper.String(v)
			}

			if v, ok := filtersMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				for i := range valuesSet {
					values := valuesSet[i].(string)
					filter.Values = append(filter.Values, helper.String(values))
				}
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	var respData *dlcv20210125.DescribeStandardEngineResourceGroupConfigInfoResponseParams
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcStandardEngineResourceGroupConfigInformationByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	standardEngineResourceGroupConfigInfosList := make([]map[string]interface{}, 0, len(respData.StandardEngineResourceGroupConfigInfos))
	if respData.StandardEngineResourceGroupConfigInfos != nil {
		for _, standardEngineResourceGroupConfigInfos := range respData.StandardEngineResourceGroupConfigInfos {
			standardEngineResourceGroupConfigInfosMap := map[string]interface{}{}

			if standardEngineResourceGroupConfigInfos.ResourceGroupId != nil {
				standardEngineResourceGroupConfigInfosMap["resource_group_id"] = standardEngineResourceGroupConfigInfos.ResourceGroupId
			}

			if standardEngineResourceGroupConfigInfos.DataEngineId != nil {
				standardEngineResourceGroupConfigInfosMap["data_engine_id"] = standardEngineResourceGroupConfigInfos.DataEngineId
			}

			staticConfigPairsList := make([]map[string]interface{}, 0, len(standardEngineResourceGroupConfigInfos.StaticConfigPairs))
			if standardEngineResourceGroupConfigInfos.StaticConfigPairs != nil {
				for _, staticConfigPairs := range standardEngineResourceGroupConfigInfos.StaticConfigPairs {
					staticConfigPairsMap := map[string]interface{}{}
					if staticConfigPairs.ConfigItem != nil {
						staticConfigPairsMap["config_item"] = staticConfigPairs.ConfigItem
					}

					if staticConfigPairs.ConfigValue != nil {
						staticConfigPairsMap["config_value"] = staticConfigPairs.ConfigValue
					}

					staticConfigPairsList = append(staticConfigPairsList, staticConfigPairsMap)
				}

				standardEngineResourceGroupConfigInfosMap["static_config_pairs"] = staticConfigPairsList
			}

			dynamicConfigPairsList := make([]map[string]interface{}, 0, len(standardEngineResourceGroupConfigInfos.DynamicConfigPairs))
			if standardEngineResourceGroupConfigInfos.DynamicConfigPairs != nil {
				for _, dynamicConfigPairs := range standardEngineResourceGroupConfigInfos.DynamicConfigPairs {
					dynamicConfigPairsMap := map[string]interface{}{}
					if dynamicConfigPairs.ConfigItem != nil {
						dynamicConfigPairsMap["config_item"] = dynamicConfigPairs.ConfigItem
					}

					if dynamicConfigPairs.ConfigValue != nil {
						dynamicConfigPairsMap["config_value"] = dynamicConfigPairs.ConfigValue
					}

					dynamicConfigPairsList = append(dynamicConfigPairsList, dynamicConfigPairsMap)
				}

				standardEngineResourceGroupConfigInfosMap["dynamic_config_pairs"] = dynamicConfigPairsList
			}

			if standardEngineResourceGroupConfigInfos.CreateTime != nil {
				standardEngineResourceGroupConfigInfosMap["create_time"] = standardEngineResourceGroupConfigInfos.CreateTime
			}

			if standardEngineResourceGroupConfigInfos.UpdateTime != nil {
				standardEngineResourceGroupConfigInfosMap["update_time"] = standardEngineResourceGroupConfigInfos.UpdateTime
			}

			standardEngineResourceGroupConfigInfosList = append(standardEngineResourceGroupConfigInfosList, standardEngineResourceGroupConfigInfosMap)
		}

		_ = d.Set("standard_engine_resource_group_config_infos", standardEngineResourceGroupConfigInfosList)
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
