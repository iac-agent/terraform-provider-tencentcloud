package teo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTeoPlans() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceTencentCloudTeoPlansRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 conditions， upper 限制 的 Filters. Values 是 20. detailed filtering conditions 是 作为 follows: <li>plan-类型<br>过滤器 according 到 [<strong>Package 类型</strong>]. <br>可选 types 是: <br>plan-trial: Trial Package; <br>plan-personal: Personal Package; <br>plan-basic: Basic Package; <br>plan-standard: Standard Package; <br>plan-enterprise: Enterprise Package. </li><li>plan-ID<br>过滤器 according 到 [<strong>Package ID</strong>]. 包 ID 是 在 form 的: edgeone-268z103ob0sx.</li><li>area<br>过滤器 according 到 [<strong>Package Acceleration 地域</strong>]. </li>Service area，可选 types 是: <br>mainland: Mainland China; <br>overseas: Global (excluding Mainland China); <br>全局: Global (包括 Mainland China).<br><li>状态<br>过滤器 通过 [<strong>Package 状态</strong>].<br> 可用 statuses 是:<br>normal: normal 状态;<br>expiring-soon: about 到 expire;<br>expired: expired;<br>isolated: isolated.</li>。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤名称",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "过滤值",
						},
					},
				},
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting 字段， 值 是: <li> 启用-时间: effective 时间; </li><li> expire-时间: 过期时间. </li> 如果未填写 在， 默认值 启用-时间 将 是 使用。",
			},

			"direction": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting direction， possible 值 是: <li>asc: sort 从 small 到 large; </li><li>desc: sort 从 large 到 small. </li>如果未填写 在， 默认值 desc 将 是 使用。",
			},

			// computed
			"plans": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Plan 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plan_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plan 类型 Possible 值 是: <li>plan-trial: Trial plan; </li><li>plan-personal: Personal plan; </li><li>plan-basic: Basic plan; </li><li>plan-standard: Standard plan; </li><li>plan-enterprise-v2: Enterprise plan; </li><li>plan-enterprise-model-: Enterprise Model A plan. </li><li>plan-enterprise: Old Enterprise plan. </li>。",
						},
						"plan_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plan ID。",
						},
						"area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service area， 值 是: <li>mainland: Mainland China; </li><li>overseas: Worldwide (excluding Mainland China); </li><li>全局: Worldwide (包括 Mainland China).</li>。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Package 状态， 值 是: <li>normal: normal 状态; </li><li>expiring-soon: about 到 expire; </li><li>expired: expired; </li><li>isolated: isolated; </li><li>overdue-isolated: overdue isolated.</li>。",
						},
						"pay_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Payment 类型，possible 值: <li>0: post-payment; </li><li>1: pre-payment.</li>。",
						},
						"zones_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Site 信息 bound 到 包，包括 site ID，站点名称，和 site 状态",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "可用区 ID",
									},
									"zone_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "可用区 名称",
									},
									"paused": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否site 是 已禁用 possible 值 是: <li>false: 不 已禁用; </li><li>true: 已禁用</li>。",
									},
								},
							},
						},
						"smart_request_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 intelligent acceleration requests 在 包，单位: times。",
						},
						"vau_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "VAU specifications 在 包，单位: piece。",
						},
						"acc_traffic_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "内容 acceleration 流量 specifications 在 包，单位: byte。",
						},
						"smart_traffic_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Smart acceleration 流量 specifications within 包，单位: byte。",
						},
						"ddos_traffic_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "DDoS protection 流量 specifications within 包，单位: bytes。",
						},
						"sec_traffic_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "安全 flow 规格 within 包，单位: byte。",
						},
						"sec_request_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 secure requests 在 包，单位: times。",
						},
						"l4_traffic_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Layer 4 acceleration 流量 specifications within 包，单位: byte。",
						},
						"cross_mlc_traffic_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "optimized 流量 specifications 的 Chinese mainland 网络 在 包，单位: bytes。",
						},
						"bindable": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "是否package allows binding 的 new sites， 值 是: <li>true: allows binding 的 new sites; </li><li>false: does 不 allow binding 的 new sites.</li>。",
						},
						"enabled_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "包 effective 时间。",
						},
						"expired_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "expiration date 的 包。",
						},
						"features": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "functions 支持 通过 包 have following 值: <li>ContentAcceleration: 内容 acceleration 函数; </li><li>SmartAcceleration: smart acceleration 函数; </li><li>L4: four-layer acceleration 函数; </li><li>Waf: advanced web protection; </li><li>QUIC: QUIC 函数; </li><li>CrossMLC: Chinese mainland 网络 optimization 函数; </li><li>ProcessMedia: media processing 函数; </li><li>L4DDoS: four-layer DDoS protection 函数; </li>L7DDoS 函数 将 仅 have 一个 的 following specifications <li>L7DDoS.CM30G; seven-layer DDoS protection 函数 - Chinese mainland 30G 最小 带宽 规格; </li><li>L7DDoS.CM60G; seven-layer DDoS protection 函数 - Chinese mainland 60G 最小 带宽 规格; </li> <li>L7DDoS.CM100G; Layer 7 DDoS protection 函数 - 100G guaranteed 带宽 对于 mainland China;</li><li>L7DDoS.Anycast300G; Layer 7 DDoS protection 函数 - 300G guaranteed 带宽 对于 Anycast outside mainland China;</li><li>L7DDoS.AnycastUnlimited; Layer 7 DDoS protection 函数 - unlimited full protection 对于 Anycast outside mainland China;</li><li>L7DDoS.CM30G_Anycast300G; Layer 7 DDoS protection 函数 - 30G guaranteed 带宽 对于 mainland China </li><li>L7DDoS.CM60G_Anycast300G; Layer 7 DDoS protection 函数 - 60G guaranteed 带宽 在 mainland China，300G guaranteed 带宽 在 anycast outside mainland China; </li><li>L7DDoS.CM100G_Anycast300G; Layer 7 DDoS protection 函数 - 100G guaranteed 带宽 在 mainland China，300G guaranteed 带宽 在 anycast outside mainland China; </li><li>L7DDoS.CM30G_AnycastUnlimited d; Layer 7 DDoS protection 函数 - 30G guaranteed 带宽 在 mainland China，unlimited Anycast protection outside mainland China; </li><li>L7DDoS.CM60G_AnycastUnlimited; Layer 7 DDoS protection 函数 - 60G guaranteed 带宽 在 mainland China，unlimited Anycast protection outside mainland China; </li><li>L7DDoS.CM100G_AnycastUnlimited; Layer 7 DDoS protection 函数 - 100G guaranteed 带宽 在 mainland China，unlimited Anycast protection outside mainland China; </li>。",
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

func DataSourceTencentCloudTeoPlansRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_plans.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*teo.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filter := teo.Filter{}
			filterMap := item.(map[string]interface{})
			if v, ok := filterMap["name"].(string); ok && v != "" {
				filter.Name = helper.String(v)
			}

			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["filters"] = tmpSet
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("direction"); ok {
		paramMap["Direction"] = helper.String(v.(string))
	}

	var plans []*teo.Plan
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTeoPlansByFilters(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		plans = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(plans))
	tmpList := make([]map[string]interface{}, 0, len(plans))
	if plans != nil {
		for _, plan := range plans {
			planMap := map[string]interface{}{}
			if plan.PlanType != nil {
				planMap["plan_type"] = plan.PlanType
			}

			if plan.PlanId != nil {
				planMap["plan_id"] = plan.PlanId
				ids = append(ids, *plan.PlanId)
			}

			if plan.Area != nil {
				planMap["area"] = plan.Area
			}

			if plan.Status != nil {
				planMap["status"] = plan.Status
			}

			if plan.PayMode != nil {
				planMap["pay_mode"] = plan.PayMode
			}

			if plan.ZonesInfo != nil {
				zonesInfoList := []interface{}{}
				for _, zonesInfo := range plan.ZonesInfo {
					zonesInfoMap := map[string]interface{}{}
					if zonesInfo.ZoneId != nil {
						zonesInfoMap["zone_id"] = zonesInfo.ZoneId
					}

					if zonesInfo.ZoneName != nil {
						zonesInfoMap["zone_name"] = zonesInfo.ZoneName
					}

					if zonesInfo.Paused != nil {
						zonesInfoMap["paused"] = zonesInfo.Paused
					}

					zonesInfoList = append(zonesInfoList, zonesInfoMap)
				}

				planMap["zones_info"] = zonesInfoList
			}

			if plan.SmartRequestCapacity != nil {
				planMap["smart_request_capacity"] = plan.SmartRequestCapacity
			}

			if plan.VAUCapacity != nil {
				planMap["vau_capacity"] = plan.VAUCapacity
			}

			if plan.AccTrafficCapacity != nil {
				planMap["acc_traffic_capacity"] = plan.AccTrafficCapacity
			}

			if plan.SmartTrafficCapacity != nil {
				planMap["smart_traffic_capacity"] = plan.SmartTrafficCapacity
			}

			if plan.DDoSTrafficCapacity != nil {
				planMap["ddos_traffic_capacity"] = plan.DDoSTrafficCapacity
			}

			if plan.SecTrafficCapacity != nil {
				planMap["sec_traffic_capacity"] = plan.SecTrafficCapacity
			}

			if plan.SecRequestCapacity != nil {
				planMap["sec_request_capacity"] = plan.SecRequestCapacity
			}

			if plan.L4TrafficCapacity != nil {
				planMap["l4_traffic_capacity"] = plan.L4TrafficCapacity
			}

			if plan.CrossMLCTrafficCapacity != nil {
				planMap["cross_mlc_traffic_capacity"] = plan.CrossMLCTrafficCapacity
			}

			if plan.Bindable != nil {
				planMap["bindable"] = plan.Bindable
			}

			if plan.EnabledTime != nil {
				planMap["enabled_time"] = plan.EnabledTime
			}

			if plan.ExpiredTime != nil {
				planMap["expired_time"] = plan.ExpiredTime
			}

			if plan.Features != nil {
				planMap["features"] = plan.Features
			}

			tmpList = append(tmpList, planMap)
		}

		_ = d.Set("plans", tmpList)
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
