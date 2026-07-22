package tke

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudContainerClusters() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source has been deprecated in Terraform TencentCloud provider version 1.16.0. Please use `tencentcloud_kubernetes_clusters` instead.",
		Read:               dataSourceTencentCloudContainerClustersRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An ID identify 集群，like `cls-xxxxxx`。",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "An int variable describe how many 集群 在 返回 在 most。",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "数量 clusters。",
			},
			"clusters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An 信息 列表 kubernetes clusters。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An ID identify 集群，like `cls-xxxxxx`。",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 集群。",
						},
						"security_certification_authority": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Describe 证书 字符串 needed 对于 使用 kubectl 到 访问 到 kubernetes。",
						},
						"security_cluster_external_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Describe 地址 needed 对于 使用 kubectl 到 访问 到 kubernetes。",
						},
						"security_username": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Describe 用户名 needed 对于 使用 kubectl 到 访问 到 kubernetes。",
						},
						"security_password": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Describe 密码 needed 对于 使用 kubectl 到 访问 到 kubernetes。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 集群。",
						},
						"kubernetes_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Describe running kubernetes 版本 在 集群。",
						},
						"nodes_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Describe how many 集群 实例 在 集群。",
						},
						"nodes_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Describe 当前 状态 实例 在 集群。",
						},
						"total_cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Describe 总数 cpu 的 each 实例 在 集群。",
						},
						"total_mem": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Describe 总数 内存 的 each 实例 在 集群。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudContainerClustersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_container_clusters.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := tke.NewDescribeClustersRequest()
	if clusterId, ok := d.GetOkExists("cluster_id"); ok {
		request.ClusterIds = []*string{common.StringPtr(clusterId.(string))}
	}

	if limit, ok := d.GetOkExists("limit"); ok {
		request.Limit = common.Int64Ptr(limit.(int64))
	}

	var response *tke.DescribeClustersResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeClient().DescribeClusters(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s DescribeClusters failed, reason:%s\n ", logId, err.Error())
		return err
	}

	ids := make([]string, 0, *response.Response.TotalCount)
	clustersList := make([]map[string]interface{}, 0, *response.Response.TotalCount)
	for _, cluster := range response.Response.Clusters {
		ids = append(ids, *cluster.ClusterId)

		clusterInfo := make(map[string]interface{}, 1)
		clusterInfo["cluster_id"] = *cluster.ClusterId
		clusterInfo["cluster_name"] = *cluster.ClusterName
		clusterInfo["description"] = *cluster.ClusterDescription
		clusterInfo["kubernetes_version"] = *cluster.ClusterVersion
		clusterInfo["nodes_num"] = *cluster.ClusterNodeNum
		clusterInfo["nodes_status"] = *cluster.ClusterStatus
		clusterInfo["total_cpu"] = int64(0)
		clusterInfo["total_mem"] = int64(0)

		describeClusterInstancesreq := tke.NewDescribeClusterInstancesRequest()
		describeClusterInstancesreq.ClusterId = cluster.ClusterId
		describeClusterInstancesreq.Limit = common.Int64Ptr(100)
		var describeClusterInstancesResponse *tke.DescribeClusterInstancesResponse
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeClient().DescribeClusterInstances(describeClusterInstancesreq)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, describeClusterInstancesreq.GetAction(), describeClusterInstancesreq.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			describeClusterInstancesResponse = result
			return nil
		})
		if err != nil {
			continue
		}

		instanceIds := []*string{}
		for _, v := range describeClusterInstancesResponse.Response.InstanceSet {
			instanceIds = append(instanceIds, v.InstanceId)
		}

		if len(instanceIds) > 0 {
			describeInstancesreq := cvm.NewDescribeInstancesRequest()
			describeInstancesreq.InstanceIds = instanceIds
			var describeInstancesResponse *cvm.DescribeInstancesResponse
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().DescribeInstances(describeInstancesreq)
				if e != nil {
					log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
						logId, describeInstancesreq.GetAction(), describeInstancesreq.ToJsonString(), e.Error())
					return tccommon.RetryError(e)
				}
				describeInstancesResponse = result
				return nil
			})
			if err != nil {
				log.Printf("[CRITAL]%s DescribeInstances failed, reason:%s\n ", logId, err.Error())
				return err
			}

			for _, v := range describeInstancesResponse.Response.InstanceSet {
				clusterInfo["total_cpu"] = clusterInfo["total_cpu"].(int64) + *v.CPU
				clusterInfo["total_mem"] = clusterInfo["total_mem"].(int64) + *v.Memory
			}
		}

		describeClusterSecurityreq := tke.NewDescribeClusterSecurityRequest()
		describeClusterSecurityreq.ClusterId = cluster.ClusterId
		var securityResponse *tke.DescribeClusterSecurityResponse
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeClient().DescribeClusterSecurity(describeClusterSecurityreq)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, describeClusterSecurityreq.GetAction(), describeClusterSecurityreq.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			securityResponse = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s DescribeClusterSecurity failed, reason:%s\n ", logId, err.Error())
			return err
		}

		clusterInfo["security_certification_authority"] = *securityResponse.Response.CertificationAuthority
		clusterInfo["security_cluster_external_endpoint"] = *securityResponse.Response.ClusterExternalEndpoint
		clusterInfo["security_username"] = *securityResponse.Response.UserName
		clusterInfo["security_password"] = *securityResponse.Response.Password
		clustersList = append(clustersList, clusterInfo)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("clusters", clustersList)
	_ = d.Set("total_count", *response.Response.TotalCount)

	return nil
}
