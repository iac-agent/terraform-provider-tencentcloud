package teo

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func DataSourceTencentCloudTeoDDoSAttackTopData() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoDDoSAttackTopDataRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Start time of the query range.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "End time of the query range. The query time range (EndTime - StartTime) must be less than or equal to 31 days.",
			},
			"metric_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The statistical metric to query. Valid values: `ddos_attackFlux_protocol`, `ddos_attackPackageNum_protocol`, `ddos_attackNum_attackType`, `ddos_attackNum_sregion`, `ddos_attackFlux_sip`, `ddos_attackFlux_sregion`.",
			},
			"zone_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of zone IDs to query. Up to 100 zone IDs. Use `*` to query all zones under the master account.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"policy_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of DDoS policy IDs to query. If not specified, all policy IDs are selected by default.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"attack_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Attack type filter. Valid values: `flood`, `icmpFlood`, `all`. Default is `all`.",
			},
			"protocol_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Protocol type filter. Valid values: `tcp`, `udp`, `all`. Default is `all`.",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Port number filter.",
			},
			"area": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Data area. Valid values: `overseas`, `mainland`. If not specified, the area is intelligently selected based on user location.",
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "DDoS attack Top data list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The dimension value of the Top query.",
						},
						"value": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of TopEntryValue items.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ranking entity name.",
									},
									"count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The ranking entity count.",
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
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTeoDDoSAttackTopDataRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_d_do_s_attack_top_data.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client()

	request := teo.NewDescribeDDoSAttackTopDataRequest()

	if v, ok := d.GetOk("start_time"); ok {
		request.StartTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("metric_name"); ok {
		request.MetricName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("zone_ids"); ok {
		zoneIdsSet := v.(*schema.Set).List()
		zoneIds := make([]*string, 0, len(zoneIdsSet))
		for _, item := range zoneIdsSet {
			zoneIds = append(zoneIds, helper.String(item.(string)))
		}
		request.ZoneIds = zoneIds
	}

	if v, ok := d.GetOk("policy_ids"); ok {
		policyIdsSet := v.(*schema.Set).List()
		policyIds := make([]*int64, 0, len(policyIdsSet))
		for _, item := range policyIdsSet {
			policyIds = append(policyIds, helper.IntInt64(item.(int)))
		}
		request.PolicyIds = policyIds
	}

	if v, ok := d.GetOk("attack_type"); ok {
		request.AttackType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol_type"); ok {
		request.ProtocolType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		request.Port = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("area"); ok {
		request.Area = helper.String(v.(string))
	}

	request.Limit = helper.IntInt64(100)

	ratelimit.Check(request.GetAction())

	var response *teo.DescribeDDoSAttackTopDataResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := client.DescribeDDoSAttackTopDataWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result == nil || result.Response == nil || len(result.Response.Data) == 0 {
			log.Printf("[DATASOURCE] tencentcloud_teo_d_do_s_attack_top_data read empty, skip SetId")
			return tccommon.RetryError(fmt.Errorf("DescribeDDoSAttackTopData returned empty response"))
		}
		response = result
		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	dataList := make([]map[string]interface{}, 0, len(response.Response.Data))
	for _, topEntry := range response.Response.Data {
		topEntryMap := map[string]interface{}{}

		if topEntry.Key != nil {
			topEntryMap["key"] = *topEntry.Key
		}

		valueList := make([]map[string]interface{}, 0, len(topEntry.Value))
		if topEntry.Value != nil {
			for _, topEntryValue := range topEntry.Value {
				valueMap := map[string]interface{}{}

				if topEntryValue.Name != nil {
					valueMap["name"] = *topEntryValue.Name
				}

				if topEntryValue.Count != nil {
					valueMap["count"] = int(*topEntryValue.Count)
				}

				valueList = append(valueList, valueMap)
			}
		}
		topEntryMap["value"] = valueList

		dataList = append(dataList, topEntryMap)
	}

	_ = d.Set("data", dataList)

	// Build ID from required parameters: startTime#endTime#metricName
	id := fmt.Sprintf("%s%s%s%s%s", *request.StartTime, tccommon.FILED_SP, *request.EndTime, tccommon.FILED_SP, *request.MetricName)
	d.SetId(id)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataList); e != nil {
			return e
		}
	}

	return nil
}
