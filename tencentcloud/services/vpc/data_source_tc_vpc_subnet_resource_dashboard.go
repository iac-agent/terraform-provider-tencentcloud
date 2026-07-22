package vpc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpcSubnetResourceDashboard() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcSubnetResourceDashboardRead,
		Schema: map[string]*schema.Schema{
			"subnet_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "子网实例 ID，such 作为 `子网-f1xjkw1b`。",
			},

			"resource_statistics_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Information 的 resources 返回。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC 实例 ID，such 作为 vpc-f1xjkw1b。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "子网实例 ID，such 作为 `子网-bthucmmy`。",
						},
						"ip": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "总数 数量 使用 IP addresses。",
						},
						"resource_statistics_item_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Information 的 associated resources。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "资源类型，such 作为 CVM，ENI。",
									},
									"resource_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "资源名称",
									},
									"resource_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 resources。",
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

func dataSourceTencentCloudVpcSubnetResourceDashboardRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_subnet_resource_dashboard.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("subnet_ids"); ok {
		subnetIdsSet := v.(*schema.Set).List()
		paramMap["SubnetIds"] = helper.InterfacesStringsPoint(subnetIdsSet)
	}

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var resourceStatisticsSet []*vpc.ResourceStatistics

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcSubnetResourceDashboardByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		resourceStatisticsSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(resourceStatisticsSet))
	tmpList := make([]map[string]interface{}, 0, len(resourceStatisticsSet))

	if resourceStatisticsSet != nil {
		for _, resourceStatistics := range resourceStatisticsSet {
			resourceStatisticsMap := map[string]interface{}{}

			if resourceStatistics.VpcId != nil {
				resourceStatisticsMap["vpc_id"] = resourceStatistics.VpcId
			}

			if resourceStatistics.SubnetId != nil {
				resourceStatisticsMap["subnet_id"] = resourceStatistics.SubnetId
			}

			if resourceStatistics.Ip != nil {
				resourceStatisticsMap["ip"] = resourceStatistics.Ip
			}

			if resourceStatistics.ResourceStatisticsItemSet != nil {
				resourceStatisticsItemSetList := []interface{}{}
				for _, resourceStatisticsItemSet := range resourceStatistics.ResourceStatisticsItemSet {
					resourceStatisticsItemSetMap := map[string]interface{}{}

					if resourceStatisticsItemSet.ResourceType != nil {
						resourceStatisticsItemSetMap["resource_type"] = resourceStatisticsItemSet.ResourceType
					}

					if resourceStatisticsItemSet.ResourceName != nil {
						resourceStatisticsItemSetMap["resource_name"] = resourceStatisticsItemSet.ResourceName
					}

					if resourceStatisticsItemSet.ResourceCount != nil {
						resourceStatisticsItemSetMap["resource_count"] = resourceStatisticsItemSet.ResourceCount
					}

					resourceStatisticsItemSetList = append(resourceStatisticsItemSetList, resourceStatisticsItemSetMap)
				}

				resourceStatisticsMap["resource_statistics_item_set"] = resourceStatisticsItemSetList
			}

			ids = append(ids, *resourceStatistics.SubnetId)
			tmpList = append(tmpList, resourceStatisticsMap)
		}

		_ = d.Set("resource_statistics_set", tmpList)
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
