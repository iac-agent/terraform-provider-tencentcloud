package dlc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlcv20210125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDlcDataEngineNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcDataEngineNetworkRead,
		Schema: map[string]*schema.Schema{
			"sort_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort Field。",
			},

			"sorting": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "排序顺序，asc 或 desc。",
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "过滤器 conditions 是 可选，引擎-网络-ID--引擎 网络 ID，引擎-网络-state--引擎 网络 状态",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Attribute 名称，如果 there 是 多个 filters， relationship between filters 是 logical OR relationship。",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Attribute 值，如果 there 是 多个 值， relationship between 值 是 logical OR relationship。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"engine_networks_infos": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Engine 网络 信息 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine_network_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Engine 网络 名称",
						},
						"engine_network_state": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Engine 网络 状态，0--initialized，2--可用，-1--删除。",
						},
						"engine_network_cidr": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Engine 网络 CIDR。",
						},
						"engine_network_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Engine 网络 ID。",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "创建时间。",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "更新时间。",
						},
						"private_link_number": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "数量 私有 links。",
						},
						"engine_number": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "数量 engines。",
						},
						"gate_way_info": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Gateway 信息 列表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"gateway_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "网关 ID",
									},
									"gateway_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Gateway 名称",
									},
									"size": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Gateway 大小。",
									},
									"state": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Gateway 状态: -1--Failed，-2--Deleted，0--Init,1--Pause，2--running，3--ToBeDeleted，4--Deleting，5--Pausing，6--Resuming，7--Isolating，8--Isolated，9--Renewing，10--Modifying，11--Modified，12--Restoring，13--Restored，14--ToBeRestored。",
									},
									"pay_mode": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "付费模式",
									},
									"mode": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Gateway 模式",
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

func dataSourceTencentCloudDlcDataEngineNetworkRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_data_engine_network.read")()
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
				valuesSet := v.([]interface{})
				for i := range valuesSet {
					values := valuesSet[i].(string)
					filter.Values = append(filter.Values, helper.String(values))
				}
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	var respData []*dlcv20210125.EngineNetworkInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcDataEngineNetworkByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	engineNetworkInfosList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, engineNetworkInfos := range respData {
			engineNetworkInfosMap := map[string]interface{}{}
			if engineNetworkInfos.EngineNetworkName != nil {
				engineNetworkInfosMap["engine_network_name"] = engineNetworkInfos.EngineNetworkName
			}

			if engineNetworkInfos.EngineNetworkState != nil {
				engineNetworkInfosMap["engine_network_state"] = engineNetworkInfos.EngineNetworkState
			}

			if engineNetworkInfos.EngineNetworkCidr != nil {
				engineNetworkInfosMap["engine_network_cidr"] = engineNetworkInfos.EngineNetworkCidr
			}

			if engineNetworkInfos.EngineNetworkId != nil {
				engineNetworkInfosMap["engine_network_id"] = engineNetworkInfos.EngineNetworkId
			}

			if engineNetworkInfos.CreateTime != nil {
				engineNetworkInfosMap["create_time"] = engineNetworkInfos.CreateTime
			}

			if engineNetworkInfos.UpdateTime != nil {
				engineNetworkInfosMap["update_time"] = engineNetworkInfos.UpdateTime
			}

			if engineNetworkInfos.PrivateLinkNumber != nil {
				engineNetworkInfosMap["private_link_number"] = engineNetworkInfos.PrivateLinkNumber
			}

			if engineNetworkInfos.EngineNumber != nil {
				engineNetworkInfosMap["engine_number"] = engineNetworkInfos.EngineNumber
			}

			gateWayInfoList := make([]map[string]interface{}, 0, len(engineNetworkInfos.GateWayInfo))
			if engineNetworkInfos.GateWayInfo != nil {
				for _, gateWayInfo := range engineNetworkInfos.GateWayInfo {
					gateWayInfoMap := map[string]interface{}{}

					if gateWayInfo.GatewayId != nil {
						gateWayInfoMap["gateway_id"] = gateWayInfo.GatewayId
					}

					if gateWayInfo.GatewayName != nil {
						gateWayInfoMap["gateway_name"] = gateWayInfo.GatewayName
					}

					if gateWayInfo.Size != nil {
						gateWayInfoMap["size"] = gateWayInfo.Size
					}

					if gateWayInfo.State != nil {
						gateWayInfoMap["state"] = gateWayInfo.State
					}

					if gateWayInfo.PayMode != nil {
						gateWayInfoMap["pay_mode"] = gateWayInfo.PayMode
					}

					if gateWayInfo.Mode != nil {
						gateWayInfoMap["mode"] = gateWayInfo.Mode
					}

					gateWayInfoList = append(gateWayInfoList, gateWayInfoMap)
				}

				engineNetworkInfosMap["gate_way_info"] = gateWayInfoList
			}

			engineNetworkInfosList = append(engineNetworkInfosList, engineNetworkInfosMap)
		}

		_ = d.Set("engine_network_infos", engineNetworkInfosList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), engineNetworkInfosList); e != nil {
			return e
		}
	}

	return nil
}
