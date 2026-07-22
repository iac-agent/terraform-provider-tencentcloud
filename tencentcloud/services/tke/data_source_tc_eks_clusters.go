/*
Use this data source to query elastic kubernetes cluster resource (offlined).

~> **NOTE:**  This resource was offline and no longer supported.

# Example Usage

```

	data "tencentcloud_eks_clusters" "foo" {
	  cluster_id = "cls-xxxxxxxx"
	}

```
*/
package tke

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudEKSClusters() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource was offline and no longer supported.",
		Read:               dataSourceTencentCloudEKSClustersRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"cluster_name"},
				Description:   "ID cluster. Conflict with cluster_name，can not be set at the same time。",
				Optional:      true,
			},
			"cluster_name": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"cluster_id"},
				Optional:      true,
				Description:   "名称 cluster. Conflict with cluster_id，can not be set at the same time。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "EKS cluster list。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID cluster。",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 cluster。",
						},
						"cluster_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 cluster。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "私有网络 ID",
						},
						"subnet_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "子网 ID list。",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"k8s_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "EKS cluster kubernetes 版本",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "EKS 状态",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 of the clusters。",
						},
						"service_subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "子网 ID service。",
						},
						"dns_servers": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "列表 cluster custom DNS Server info。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DNS Server 域名 Empty 表示all 域名",
									},
									"servers": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "列表 DNS Server IP 地址",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"need_delete_cbs": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "表示是否delete CBS after EKS cluster remove。",
						},
						"enable_vpc_core_dns": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "表示是否enable dns in 用户 cluster，默认值为 `true`。",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "标签 of EKS cluster。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudEKSClustersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_eks_clusters.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := EksService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	var (
		id   string
		name string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		id = v.(string)
	}

	if v, ok := d.GetOk("cluster_name"); ok {
		name = v.(string)
	}

	tags := helper.GetTags(d, "tags")

	infos, err := service.DescribeEKSClusters(ctx, id, name)
	if err != nil && id == "" {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			infos, err = service.DescribeEKSClusters(ctx, id, name)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(infos))
	ids := make([]string, 0, len(infos))

LOOP:
	for _, info := range infos {
		if len(tags) > 0 {
			for k, v := range tags {
				if info.Tags[k] != v {
					continue LOOP
				}
			}
		}
		var infoMap = map[string]interface{}{
			"cluster_id":          info.ClusterId,
			"cluster_name":        info.ClusterName,
			"cluster_desc":        info.ClusterDesc,
			"vpc_id":              info.VpcId,
			"subnet_ids":          info.SubnetIds,
			"dns_servers":         info.DnsServers,
			"k8s_version":         info.K8SVersion,
			"status":              info.Status,
			"created_time":        info.CreatedTime,
			"service_subnet_id":   info.ServiceSubnetId,
			"need_delete_cbs":     info.NeedDeleteCbs,
			"enable_vpc_core_dns": info.EnableVpcCoreDNS,
		}

		list = append(list, infoMap)
		ids = append(ids, info.ClusterId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("list", list)
	if err != nil {
		log.Printf("[CRITAL]%s provider set tencentcloud_eks_clusters list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err = tccommon.WriteToFile(output.(string), list); err != nil {
			return err
		}
	}
	return nil
}
