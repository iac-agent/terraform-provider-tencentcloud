package clb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClbExclusiveClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbExclusiveClustersRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤查询可用区资源列表，具体如下： cluster-类型 - 字符串 - 必填：否 - （过滤条件）按集群类型过滤，如TGW。 cluster-id - String - 必填：否 - （过滤条件）按集群 ID 过滤，如 tgw-xxxxxxxx。 cluster-名称 - String - 必填：否 - （过滤条件）按集群名称过滤，如test-xxxxxx。 cluster-标签 - String - 必填：否 - （过滤条件）按集群标签过滤，如 TAG-xxxxx。 VIP - String - 必填：否 - （过滤条件）按集群中的vip进行过滤，如x.x.x.x。 network - 字符串 - 必需：否 - （过滤条件）按集群网络类型过滤，例如公共或专用。 可用区 - 字符串 - 必填：否 - （过滤条件）按集群区域过滤，如 ap-guangzhou-1。 isp - 字符串 - 必需：否 - （过滤条件）按 TGW 集群 isp 类型过滤，例如 BGP。 loadblancer-id - 字符串 - 必填：否 - （过滤条件）按集群中的 loadblancer-id 进行过滤，如 lb-xxxxxxxx。",
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
							Description: "过滤值数组。",
						},
					},
				},
			},

			"cluster_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "集群列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "集群ID。",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "集群名称。",
						},
						"cluster_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "集群类型：TGW、STGW、VPCGW。",
						},
						"cluster_tag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "专用第 7 层标签。注意：该字段可能返回null，表示取不到有效值。",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "。",
						},
						"network": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "集群网络类型。",
						},
						"max_conn": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最大连接数。",
						},
						"max_in_flow": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最大传入带宽。",
						},
						"max_in_pkg": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最大传入数据包。",
						},
						"max_out_flow": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最大输出带宽。",
						},
						"max_out_pkg": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最大输出数据包。",
						},
						"max_new_conn": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最大新连接数。",
						},
						"http_max_new_conn": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "新 http 连接的最大数量。",
						},
						"https_max_new_conn": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "新 https 连接的最大数量。",
						},
						"http_qps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Http Qps。",
						},
						"https_qps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "https Qps。",
						},
						"resource_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "集群中的资源总数。",
						},
						"idle_resource_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "集群中可用资源的总数。",
						},
						"load_balance_director_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "集群中转发器的总数。",
						},
						"isp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ISP：BGP、CMCC、CUCC、CTCC、内部。",
						},
						"clusters_zone": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "集群所在可用区。",
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
						"clusters_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "集群版本。",
						},
						"disaster_recovery_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "集群容灾类型：SINGLE-ZONE、DISASTER-RECOVERY、MUTUAL-DISASTER-RECOVERY。",
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

func dataSourceTencentCloudClbExclusiveClustersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clb_exclusive_clusters.read")()
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

	var clusterSet []*clb.Cluster

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClbExclusiveClustersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		clusterSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(clusterSet))
	tmpList := make([]map[string]interface{}, 0, len(clusterSet))

	if clusterSet != nil {
		for _, cluster := range clusterSet {
			clusterMap := map[string]interface{}{}

			if cluster.ClusterId != nil {
				clusterMap["cluster_id"] = cluster.ClusterId
			}

			if cluster.ClusterName != nil {
				clusterMap["cluster_name"] = cluster.ClusterName
			}

			if cluster.ClusterType != nil {
				clusterMap["cluster_type"] = cluster.ClusterType
			}

			if cluster.ClusterTag != nil {
				clusterMap["cluster_tag"] = cluster.ClusterTag
			}

			if cluster.Zone != nil {
				clusterMap["zone"] = cluster.Zone
			}

			if cluster.Network != nil {
				clusterMap["network"] = cluster.Network
			}

			if cluster.MaxConn != nil {
				clusterMap["max_conn"] = cluster.MaxConn
			}

			if cluster.MaxInFlow != nil {
				clusterMap["max_in_flow"] = cluster.MaxInFlow
			}

			if cluster.MaxInPkg != nil {
				clusterMap["max_in_pkg"] = cluster.MaxInPkg
			}

			if cluster.MaxOutFlow != nil {
				clusterMap["max_out_flow"] = cluster.MaxOutFlow
			}

			if cluster.MaxOutPkg != nil {
				clusterMap["max_out_pkg"] = cluster.MaxOutPkg
			}

			if cluster.MaxNewConn != nil {
				clusterMap["max_new_conn"] = cluster.MaxNewConn
			}

			if cluster.HTTPMaxNewConn != nil {
				clusterMap["http_max_new_conn"] = cluster.HTTPMaxNewConn
			}

			if cluster.HTTPSMaxNewConn != nil {
				clusterMap["https_max_new_conn"] = cluster.HTTPSMaxNewConn
			}

			if cluster.HTTPQps != nil {
				clusterMap["http_qps"] = cluster.HTTPQps
			}

			if cluster.HTTPSQps != nil {
				clusterMap["https_qps"] = cluster.HTTPSQps
			}

			if cluster.ResourceCount != nil {
				clusterMap["resource_count"] = cluster.ResourceCount
			}

			if cluster.IdleResourceCount != nil {
				clusterMap["idle_resource_count"] = cluster.IdleResourceCount
			}

			if cluster.LoadBalanceDirectorCount != nil {
				clusterMap["load_balance_director_count"] = cluster.LoadBalanceDirectorCount
			}

			if cluster.Isp != nil {
				clusterMap["isp"] = cluster.Isp
			}

			if cluster.ClustersZone != nil {
				clustersZoneMap := map[string]interface{}{}

				if cluster.ClustersZone.MasterZone != nil {
					clustersZoneMap["master_zone"] = cluster.ClustersZone.MasterZone
				}

				if cluster.ClustersZone.SlaveZone != nil {
					clustersZoneMap["slave_zone"] = cluster.ClustersZone.SlaveZone
				}

				clusterMap["clusters_zone"] = []interface{}{clustersZoneMap}
			}

			if cluster.ClustersVersion != nil {
				clusterMap["clusters_version"] = cluster.ClustersVersion
			}

			if cluster.DisasterRecoveryType != nil {
				clusterMap["disaster_recovery_type"] = cluster.DisasterRecoveryType
			}

			ids = append(ids, *cluster.ClusterId)
			tmpList = append(tmpList, clusterMap)
		}

		_ = d.Set("cluster_set", tmpList)
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
