package clb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClbClusterResources() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbClusterResourcesRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤查询集群的条件。 cluster-id - String - 必填：否 - （过滤条件）按集群 ID 过滤，如 tgw-12345678。 VIP - 字符串 - 必填：否 - （过滤条件）按负载均衡器 VIP 过滤，例如 192.168.0.1。 loadblancer-id - 字符串 - 必填：否 - （过滤条件）按 loadblancer ID 过滤，例如 lbl-12345678。 idle - String - 必填：否 - （过滤条件）过滤条件 负载均衡是否空闲，如True、False。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤器名称。",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "过滤值。",
						},
					},
				},
			},

			"cluster_resource_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "集群资源集。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "集群 ID。",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "贵宾。",
						},
						"load_balancer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "负载平衡 ID。",
						},
						"idle": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "是不是闲着呢。",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "集群名称。",
						},
						"isp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "等离子。",
						},
						"clusters_zone": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "集群区。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"master_zone": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "集群所在的可用性主域。",
									},
									"slave_zone": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "集群所在的可用从区。",
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

func dataSourceTencentCloudClbClusterResourcesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clb_cluster_resources.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*clb.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := clb.Filter{}
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
		paramMap["Filters"] = tmpSet
	}

	service := ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var clusterResourceSet []*clb.ClusterResource

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClbClusterResourcesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		clusterResourceSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(clusterResourceSet))
	tmpList := make([]map[string]interface{}, 0, len(clusterResourceSet))

	if clusterResourceSet != nil {
		for _, clusterResource := range clusterResourceSet {
			clusterResourceMap := map[string]interface{}{}

			if clusterResource.ClusterId != nil {
				clusterResourceMap["cluster_id"] = clusterResource.ClusterId
			}

			if clusterResource.Vip != nil {
				clusterResourceMap["vip"] = clusterResource.Vip
			}

			if clusterResource.LoadBalancerId != nil {
				clusterResourceMap["load_balancer_id"] = clusterResource.LoadBalancerId
			}

			if clusterResource.Idle != nil {
				clusterResourceMap["idle"] = clusterResource.Idle
			}

			if clusterResource.ClusterName != nil {
				clusterResourceMap["cluster_name"] = clusterResource.ClusterName
			}

			if clusterResource.Isp != nil {
				clusterResourceMap["isp"] = clusterResource.Isp
			}

			if clusterResource.ClustersZone != nil {
				clustersZoneMap := map[string]interface{}{}

				if clusterResource.ClustersZone.MasterZone != nil {
					clustersZoneMap["master_zone"] = clusterResource.ClustersZone.MasterZone
				}

				if clusterResource.ClustersZone.SlaveZone != nil {
					clustersZoneMap["slave_zone"] = clusterResource.ClustersZone.SlaveZone
				}

				clusterResourceMap["clusters_zone"] = []interface{}{clustersZoneMap}
			}

			ids = append(ids, *clusterResource.ClusterId)
			tmpList = append(tmpList, clusterResourceMap)
		}

		_ = d.Set("cluster_resource_set", tmpList)
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
