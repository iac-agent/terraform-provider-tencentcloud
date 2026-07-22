package cynosdb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCynosdbClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbClustersRead,

		Schema: map[string]*schema.Schema{
			"db_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CynosDB 类型，可选值包括`MYSQL`、`POSTGRESQL`。",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "需要查询的集群ID。",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "需要查询的项目ID。",
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "需要查询的集群名称。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"cluster_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "集群列表。每个元素包含以下属性：",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "项目 ID。",
						},
						"available_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CynosDB集群的可用区。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "专有网络ID。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "该VPC内子网的ID。",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CynosDB集群的端口。",
						},
						"db_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CynosDB 类型，可用值包括“MYSQL”。",
						},
						"db_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CynosDB 的版本，与 db_type 相关。对于“MYSQL”，可用值为“5.7”。",
						},
						"cluster_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CynosDB集群实例的存储限制，单位为GB。",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CynosDB 集群的名称。",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CynosDB集群的ID。",
						},
						// payment
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例的收费类型。有效值为“PREPAID”和“POSTPAID_BY_HOUR”。默认值为“POSTPAID_BY_HOUR”。",
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "自动更新标志。有效值为“0”(MANUAL_RENEW)、“1”(AUTO_RENEW)。仅适用于 PREPAID 集群。",
						},
						"cluster_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cynosdb 集群的状态。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CynosDB 集群的创建时间。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCynosdbClustersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cynosdb_clusters.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	params := make(map[string]string)
	if v, ok := d.GetOk("cluster_id"); ok {
		params["ClusterId"] = v.(string)
	}
	if v, ok := d.GetOk("cluster_name"); ok {
		params["ClusterName"] = v.(string)
	}
	if v, ok := d.GetOkExists("project_id"); ok {
		params["ProjectId"] = fmt.Sprintf("%d", v.(int))
	}
	clusterType := ""
	if v, ok := d.GetOk("cluster_type"); ok {
		clusterType = v.(string)
	}

	cynosdbService := CynosdbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	var clusters []*cynosdb.CynosdbCluster
	var err error
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		clusters, err = cynosdbService.DescribeClusters(ctx, params)
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read cynosdb clusters failed, reason:%s\n ", logId, err.Error())
		return err
	}

	ids := make([]string, 0, len(clusters))
	clusterList := make([]map[string]interface{}, 0, len(clusters))
	for _, cluster := range clusters {
		if clusterType != "" && clusterType != *cluster.DbType {
			continue
		}
		mapping := map[string]interface{}{
			"cluster_id":      cluster.ClusterId,
			"cluster_name":    cluster.ClusterName,
			"cluster_limit":   cluster.StorageLimit,
			"db_type":         cluster.DbType,
			"available_zone":  cluster.Zone,
			"project_id":      cluster.ProjectID,
			"create_time":     cluster.CreateTime,
			"cluster_status":  cluster.Status,
			"auto_renew_flag": cluster.RenewFlag,
			"port":            cluster.Vport,
			"vpc_id":          cluster.VpcId,
			"subnet_id":       cluster.SubnetId,
			"db_version":      cluster.DbVersion,
			"charge_type":     CYNOSDB_CHARGE_TYPE[*cluster.PayMode],
		}

		clusterList = append(clusterList, mapping)
		ids = append(ids, *cluster.ClusterId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("cluster_list", clusterList)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err = tccommon.WriteToFile(output.(string), clusterList); err != nil {
			return err
		}
	}

	return nil
}
