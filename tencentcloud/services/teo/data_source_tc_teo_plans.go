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
				Description: "<li>Filter conditions, the upper limit of Filters.Values is 20. The detailed filtering conditions are as follows: plan-type filter by plan type (plan-trial, plan-personal, plan-basic, plan-standard, plan-enterprise); plan-id filter by plan ID (form: edgeone-268z103ob0sx); area filter by plan acceleration region (mainland, overseas, global); status filter by plan status (normal, expiring-soon, expired, isolated). </li>",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter name.",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Filter value.",
						},
					},
				},
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting field, the values are: enable-time: effective time; expire-time: expiration time. If not filled in, the default value enable-time will be used.",
			},

			"direction": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting direction, the possible values are: asc: sort from small to large; desc: sort from large to small. If not filled in, the default value desc will be used.",
			},

			// computed
			"plans": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Plan list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plan_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plan type. Possible values are: plan-trial: Trial plan; plan-personal: Personal plan; plan-basic: Basic plan; plan-standard: Standard plan; plan-enterprise-v2: Enterprise plan; plan-enterprise-model-a: Enterprise Model A plan; plan-enterprise: Old Enterprise plan.",
						},
						"plan_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plan ID.",
						},
						"area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service area, the values are: mainland: Mainland China; overseas: Worldwide (excluding Mainland China); global: Worldwide (including Mainland China).",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Package status, the values are: normal: normal status; expiring-soon: about to expire; expired: expired; isolated: isolated; overdue-isolated: overdue isolated.",
						},
						"pay_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Payment type, possible values: 0: post-payment; 1: pre-payment.",
						},
						"zones_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Site information bound to the package, including site ID, site name, and site status.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Zone ID.",
									},
									"zone_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Zone name.",
									},
									"paused": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the site is disabled. The possible values are: false: not disabled; true: disabled.",
									},
								},
							},
						},
						"smart_request_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of intelligent acceleration requests in the package, unit: times.",
						},
						"vau_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "VAU specifications in the package, unit: piece.",
						},
						"acc_traffic_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The content acceleration traffic specifications in the package, unit: byte.",
						},
						"smart_traffic_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Smart acceleration traffic specifications within the package, unit: byte.",
						},
						"ddos_traffic_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "DDoS protection traffic specifications within the package, unit: bytes.",
						},
						"sec_traffic_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The security flow specification within the package, unit: byte.",
						},
						"sec_request_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of secure requests in the package, unit: times.",
						},
						"l4_traffic_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Layer 4 acceleration traffic specifications within the package, unit: byte.",
						},
						"cross_mlc_traffic_capacity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The optimized traffic specifications of the Chinese mainland network in the package, unit: bytes.",
						},
						"bindable": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether the package allows binding of new sites, the values are: true: allows binding of new sites; false: does not allow binding of new sites.",
						},
						"enabled_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The package effective time.",
						},
						"expired_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The expiration date of the package.",
						},
						"features": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "The functions supported by the package, such as ContentAcceleration, SmartAcceleration, L4, Waf, QUIC, CrossMLC, ProcessMedia, L4DDoS and L7DDoS specifications.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
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
