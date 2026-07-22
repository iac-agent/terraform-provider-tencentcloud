package ccn

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCcnCrossBorderRegionBandwidthLimits() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCcnCrossBorderRegionBandwidthLimitsRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 condition. Currently，仅 一个 值 是 支持. 支持 字段，1)来源-地域， 值 是 like ap-guangzhou; 2)destination-地域， 值 是 like ap-shanghai; 3)ccn-ids,云 网络 ID 数组， 值 是 like ccn-12345678; 4)用户-账号-ID,用户 账号 ID， 值 是 like 12345678。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "attribute 名称",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "值 的 字段。",
						},
					},
				},
			},

			"ccn_bandwidth_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Info 的 cross 地域 ccn 实例。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ccn_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ccn ID。",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"expired_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "过期时间。",
						},
						"region_flow_control_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 的 RegionFlowControl。",
						},
						"renew_flag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "续费标识",
						},
						"ccn_region_bandwidth_limit": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "带宽 限制 的 cross 地域",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "来源 地域，such 作为 &#39;ap-shanghai&#39;。",
									},
									"destination_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "destination 地域，such 作为。",
									},
									"bandwidth_limit": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "带宽 列表(Mbps)。",
									},
								},
							},
						},
						"market_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "market ID。",
						},
						"user_account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "用户 账号 ID。",
						},
						"is_cross_border": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "如果 cross 地域",
						},
						"is_security_lock": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "`true` 表示 locked。",
						},
						"instance_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "`POSTPAID` 或 `PREPAID`。",
						},
						"update_time": {
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

func dataSourceTencentCloudCcnCrossBorderRegionBandwidthLimitsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ccn_cross_border_region_bandwidth_limits.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*vpc.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := vpc.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var ccnBandwidthSet []*vpc.CcnBandwidth

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcCcnRegionBandwidthLimitsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		ccnBandwidthSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(ccnBandwidthSet))
	tmpList := make([]map[string]interface{}, 0, len(ccnBandwidthSet))

	if ccnBandwidthSet != nil {
		for _, ccnBandwidth := range ccnBandwidthSet {
			ccnBandwidthMap := map[string]interface{}{}

			if ccnBandwidth.CcnId != nil {
				ccnBandwidthMap["ccn_id"] = ccnBandwidth.CcnId
			}

			if ccnBandwidth.CreatedTime != nil {
				ccnBandwidthMap["created_time"] = ccnBandwidth.CreatedTime
			}

			if ccnBandwidth.ExpiredTime != nil {
				ccnBandwidthMap["expired_time"] = ccnBandwidth.ExpiredTime
			}

			if ccnBandwidth.RegionFlowControlId != nil {
				ccnBandwidthMap["region_flow_control_id"] = ccnBandwidth.RegionFlowControlId
			}

			if ccnBandwidth.RenewFlag != nil {
				ccnBandwidthMap["renew_flag"] = ccnBandwidth.RenewFlag
			}

			if ccnBandwidth.CcnRegionBandwidthLimit != nil {
				ccnRegionBandwidthLimitMap := map[string]interface{}{}

				if ccnBandwidth.CcnRegionBandwidthLimit.SourceRegion != nil {
					ccnRegionBandwidthLimitMap["source_region"] = ccnBandwidth.CcnRegionBandwidthLimit.SourceRegion
				}

				if ccnBandwidth.CcnRegionBandwidthLimit.DestinationRegion != nil {
					ccnRegionBandwidthLimitMap["destination_region"] = ccnBandwidth.CcnRegionBandwidthLimit.DestinationRegion
				}

				if ccnBandwidth.CcnRegionBandwidthLimit.BandwidthLimit != nil {
					ccnRegionBandwidthLimitMap["bandwidth_limit"] = ccnBandwidth.CcnRegionBandwidthLimit.BandwidthLimit
				}

				ccnBandwidthMap["ccn_region_bandwidth_limit"] = []interface{}{ccnRegionBandwidthLimitMap}
			}

			if ccnBandwidth.MarketId != nil {
				ccnBandwidthMap["market_id"] = ccnBandwidth.MarketId
			}

			if ccnBandwidth.UserAccountID != nil {
				ccnBandwidthMap["user_account_id"] = ccnBandwidth.UserAccountID
			}

			if ccnBandwidth.IsCrossBorder != nil {
				ccnBandwidthMap["is_cross_border"] = ccnBandwidth.IsCrossBorder
			}

			if ccnBandwidth.IsSecurityLock != nil {
				ccnBandwidthMap["is_security_lock"] = ccnBandwidth.IsSecurityLock
			}

			if ccnBandwidth.InstanceChargeType != nil {
				ccnBandwidthMap["instance_charge_type"] = ccnBandwidth.InstanceChargeType
			}

			if ccnBandwidth.UpdateTime != nil {
				ccnBandwidthMap["update_time"] = ccnBandwidth.UpdateTime
			}

			ids = append(ids, *ccnBandwidth.CcnId)
			tmpList = append(tmpList, ccnBandwidthMap)
		}

		_ = d.Set("ccn_bandwidth_set", tmpList)
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
