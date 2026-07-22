package emr

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudEmrCvmQuota() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEmrCvmQuotaRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "EMR 集群 ID.",
			},

			"zone_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Zone ID.",
			},

			"post_paid_quota_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Postpaid 配额 列表 注意: 此 字段 可能 返回 null, indicating 该 无 有效 值 可以 是 获取.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"used_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Used 配额 注意: 此 字段 可能 返回 null, indicating 该 有效 值 不能 是 获取.",
						},
						"remaining_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Residual 配额 注意: 此 字段 可能 返回 null, indicating 该 有效 值 不能 是 获取.",
						},
						"total_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total 配额 注意: 此 字段 可能 返回 null, indicating 该 有效 值 不能 是 获取.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Available area 注意: 此 字段 可能 返回 null, indicating 该 有效 值 不能 是 获取.",
						},
					},
				},
			},

			"spot_paid_quota_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Biding 实例 配额 列表 注意: 此 字段 可能 返回 null, indicating 该 有效 值 不能 是 获取.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"used_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Used 配额 注意: 此 字段 可能 返回 null, indicating 该 有效 值 不能 是 获取.",
						},
						"remaining_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Residual 配额 注意: 此 字段 可能 返回 null, indicating 该 有效 值 不能 是 获取.",
						},
						"total_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total 配额 注意: 此 字段 可能 返回 null, indicating 该 有效 值 不能 是 获取.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Available area 注意: 此 字段 可能 返回 null, indicating 该 有效 值 不能 是 获取.",
						},
					},
				},
			},

			"eks_quota_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Eks 配额 注意: 此 字段 可能 返回 null, indicating 该 有效 值 不能 是 获取.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "specifications 的 marketable 资源 是 作为 follows: `TASK`, `CORE`, `MASTER`, `ROUTER`.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cpu cores.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory quantity (单位: GB).",
						},
						"number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies 最大 数量 的 resources 该 可以 是 applied 对于.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 save results.",
			},
		},
	}
}

func dataSourceTencentCloudEmrCvmQuotaRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_emr_cvm_quota.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var clusterId string

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("zone_id"); ok {
		paramMap["ZoneId"] = helper.IntInt64(v.(int))
	}

	service := EMRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var cvmQuota *emr.DescribeCvmQuotaResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeEmrCvmQuotaByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		cvmQuota = result
		return nil
	})
	if err != nil {
		return err
	}

	//ids := make([]string, 0, len(postPaidQuotaSet))
	tmpList := make([]map[string]interface{}, 0)

	if cvmQuota.PostPaidQuotaSet != nil {
		tmpList := make([]map[string]interface{}, 0, len(cvmQuota.PostPaidQuotaSet))

		for _, quotaEntity := range cvmQuota.PostPaidQuotaSet {
			quotaEntityMap := map[string]interface{}{}

			if quotaEntity.UsedQuota != nil {
				quotaEntityMap["used_quota"] = quotaEntity.UsedQuota
			}

			if quotaEntity.RemainingQuota != nil {
				quotaEntityMap["remaining_quota"] = quotaEntity.RemainingQuota
			}

			if quotaEntity.TotalQuota != nil {
				quotaEntityMap["total_quota"] = quotaEntity.TotalQuota
			}

			if quotaEntity.Zone != nil {
				quotaEntityMap["zone"] = quotaEntity.Zone
			}

			tmpList = append(tmpList, quotaEntityMap)
		}
		_ = d.Set("post_paid_quota_set", tmpList)
	}

	if cvmQuota.SpotPaidQuotaSet != nil {
		tmpList := make([]map[string]interface{}, 0, len(cvmQuota.SpotPaidQuotaSet))

		for _, quotaEntity := range cvmQuota.SpotPaidQuotaSet {
			quotaEntityMap := map[string]interface{}{}

			if quotaEntity.UsedQuota != nil {
				quotaEntityMap["used_quota"] = quotaEntity.UsedQuota
			}

			if quotaEntity.RemainingQuota != nil {
				quotaEntityMap["remaining_quota"] = quotaEntity.RemainingQuota
			}

			if quotaEntity.TotalQuota != nil {
				quotaEntityMap["total_quota"] = quotaEntity.TotalQuota
			}

			if quotaEntity.Zone != nil {
				quotaEntityMap["zone"] = quotaEntity.Zone
			}
			tmpList = append(tmpList, quotaEntityMap)
		}

		_ = d.Set("spot_paid_quota_set", tmpList)
	}

	if cvmQuota.EksQuotaSet != nil {
		tmpList := make([]map[string]interface{}, 0, len(cvmQuota.EksQuotaSet))

		for _, podSaleSpec := range cvmQuota.EksQuotaSet {
			podSaleSpecMap := map[string]interface{}{}

			if podSaleSpec.NodeType != nil {
				podSaleSpecMap["node_type"] = podSaleSpec.NodeType
			}

			if podSaleSpec.Cpu != nil {
				podSaleSpecMap["cpu"] = podSaleSpec.Cpu
			}

			if podSaleSpec.Memory != nil {
				podSaleSpecMap["memory"] = podSaleSpec.Memory
			}

			if podSaleSpec.Number != nil {
				podSaleSpecMap["number"] = podSaleSpec.Number
			}

			tmpList = append(tmpList, podSaleSpecMap)
		}

		_ = d.Set("eks_quota_set", tmpList)
	}

	d.SetId(clusterId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
