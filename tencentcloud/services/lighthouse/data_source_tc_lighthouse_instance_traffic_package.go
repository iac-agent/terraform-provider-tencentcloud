package lighthouse

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudLighthouseInstanceTrafficPackage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseInstanceTrafficPackageRead,
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "实例 ID 列表。",
			},

			"offset": {
				Optional:    true,
				Default:     0,
				Type:        schema.TypeInt,
				Description: "偏移量 默认值为 0。",
			},

			"limit": {
				Optional:    true,
				Default:     20,
				Type:        schema.TypeInt,
				Description: "数量 返回 results. 默认值为 20. Maximum 值 是 100。",
			},

			"instance_traffic_package_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 details 的 实例 流量 packages。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID",
						},
						"traffic_package_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "列表 流量 包 details。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"traffic_package_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Traffic packet ID。",
									},
									"traffic_used": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Traffic has been 使用 during effective 周期 的 流量 packet，在 bytes。",
									},
									"traffic_package_total": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "总数 流量 在 bytes during effective 周期 的 流量 packet。",
									},
									"traffic_package_remaining": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "remaining 流量 during effective 周期 的 流量 packet，在 bytes。",
									},
									"traffic_overflow": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "amount 的 流量 该 exceeds 配额 的 流量 packet during effective 周期 的 流量 packet，在 bytes。",
									},
									"start_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "开始时间 的 effective cycle 的 流量 packet. Expressed according 到 ISO8601 standard，和 使用 UTC 时间. 格式 是 YYYY-MM-DDThh:mm:ssZ。",
									},
									"end_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "结束时间 的 effective 周期 的 流量 packet. Expressed according 到 ISO8601 standard，和 使用 UTC 时间. 格式 是 YYYY-MM-DDThh:mm:ssZ。",
									},
									"deadline": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "过期时间 的 流量 包. Expressed according 到 ISO8601 standard，和 使用 UTC 时间. 格式 是 YYYY-MM-DDThh:mm:ssZ.。",
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
										Description: "Traffic packet 状态:" +
											"- `NETWORK_NORMAL`: normal." +
											"- `OVERDUE_NETWORK_DISABLED`: network disconnection due to arrears.",
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

func dataSourceTencentCloudLighthouseInstanceTrafficPackageRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_lighthouse_instance_traffic_package.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		instanceIds := make([]string, 0)
		for _, instanceId := range instanceIdsSet {
			instanceIds = append(instanceIds, instanceId.(string))
		}
		paramMap["instance_ids"] = instanceIds
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["offset"] = v.(int)
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["limit"] = v.(int)
	}

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var instanceTrafficPackageSet []*lighthouse.InstanceTrafficPackage

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseInstanceTrafficPackageByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instanceTrafficPackageSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceTrafficPackageSet))
	tmpList := make([]map[string]interface{}, 0, len(instanceTrafficPackageSet))

	if instanceTrafficPackageSet != nil {
		for _, instanceTrafficPackage := range instanceTrafficPackageSet {
			instanceTrafficPackageMap := map[string]interface{}{}

			if instanceTrafficPackage.InstanceId != nil {
				instanceTrafficPackageMap["instance_id"] = instanceTrafficPackage.InstanceId
			}

			if instanceTrafficPackage.TrafficPackageSet != nil {
				trafficPackageSetList := []map[string]interface{}{}
				for _, trafficPackageSet := range instanceTrafficPackage.TrafficPackageSet {
					trafficPackageSetMap := map[string]interface{}{}

					if trafficPackageSet.TrafficPackageId != nil {
						trafficPackageSetMap["traffic_package_id"] = trafficPackageSet.TrafficPackageId
					}

					if trafficPackageSet.TrafficUsed != nil {
						trafficPackageSetMap["traffic_used"] = trafficPackageSet.TrafficUsed
					}

					if trafficPackageSet.TrafficPackageTotal != nil {
						trafficPackageSetMap["traffic_package_total"] = trafficPackageSet.TrafficPackageTotal
					}

					if trafficPackageSet.TrafficPackageRemaining != nil {
						trafficPackageSetMap["traffic_package_remaining"] = trafficPackageSet.TrafficPackageRemaining
					}

					if trafficPackageSet.TrafficOverflow != nil {
						trafficPackageSetMap["traffic_overflow"] = trafficPackageSet.TrafficOverflow
					}

					if trafficPackageSet.StartTime != nil {
						trafficPackageSetMap["start_time"] = trafficPackageSet.StartTime
					}

					if trafficPackageSet.EndTime != nil {
						trafficPackageSetMap["end_time"] = trafficPackageSet.EndTime
					}

					if trafficPackageSet.Deadline != nil {
						trafficPackageSetMap["deadline"] = trafficPackageSet.Deadline
					}

					if trafficPackageSet.Status != nil {
						trafficPackageSetMap["status"] = trafficPackageSet.Status
					}

					trafficPackageSetList = append(trafficPackageSetList, trafficPackageSetMap)
				}

				instanceTrafficPackageMap["traffic_package_set"] = trafficPackageSetList
			}

			ids = append(ids, *instanceTrafficPackage.InstanceId)
			tmpList = append(tmpList, instanceTrafficPackageMap)
		}

		_ = d.Set("instance_traffic_package_set", tmpList)
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
