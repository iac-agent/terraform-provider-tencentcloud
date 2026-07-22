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

func DataSourceTencentCloudContainerClusterInstances() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source has been deprecated in Terraform TencentCloud provider version 1.16.0. Please use `tencentcloud_kubernetes_clusters` instead.",
		Read:               dataSourceTencentCloudContainerClusterInstancesRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "An ID identify 集群，like cls-xxxxxx。",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "An int variable describe how many 实例 在 返回 在 most。",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "数量 实例。",
			},
			"nodes": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "An 信息 列表 kubernetes 实例。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"abnormal_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Describe reason 当 节点 是 在 abnormal state(如果 它 是)。",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Describe cpu 的 节点。",
						},
						"mem": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Describe 内存 的 节点。",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An ID identify 节点，提供 通过 cvm。",
						},
						"is_normal": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Describe 是否node 是 normal。",
						},
						"wan_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Describe WAN IP 的 节点。",
						},
						"lan_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Describe LAN IP 的 节点。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudContainerClusterInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_container_cluster_instances.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := tke.NewDescribeClusterInstancesRequest()
	if clusterId, ok := d.GetOkExists("cluster_id"); ok {
		request.ClusterId = common.StringPtr(clusterId.(string))
	}

	if limit, ok := d.GetOkExists("limit"); ok {
		request.Limit = common.Int64Ptr(limit.(int64))
	}

	var response *tke.DescribeClusterInstancesResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeClient().DescribeClusterInstances(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s DescribeClusterInstances failed, reason:%s\n ", logId, err.Error())
		return err
	}

	nodes := make([]map[string]interface{}, 0, *response.Response.TotalCount)
	ids := make([]string, 0, *response.Response.TotalCount)
	for _, node := range response.Response.InstanceSet {
		ids = append(ids, *node.InstanceId)

		nodeInfo := make(map[string]interface{})
		nodeInfo["instance_id"] = *node.InstanceId
		nodeInfo["abnormal_reason"] = *node.FailedReason
		nodeInfo["wan_ip"] = ""
		nodeInfo["lan_ip"] = ""
		nodeInfo["cpu"] = 0
		nodeInfo["mem"] = 0
		if *node.InstanceState == "failed" {
			nodeInfo["is_normal"] = 0
		} else {
			nodeInfo["is_normal"] = 1
		}

		describeInstancesreq := cvm.NewDescribeInstancesRequest()
		describeInstancesreq.InstanceIds = []*string{node.InstanceId}
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

		if len(describeInstancesResponse.Response.InstanceSet) > 0 {
			nodeInfo["cpu"] = *describeInstancesResponse.Response.InstanceSet[0].CPU
			nodeInfo["mem"] = *describeInstancesResponse.Response.InstanceSet[0].Memory
			if len(describeInstancesResponse.Response.InstanceSet[0].PublicIpAddresses) > 0 {
				nodeInfo["wan_ip"] = *describeInstancesResponse.Response.InstanceSet[0].PublicIpAddresses[0]
			}
			if len(describeInstancesResponse.Response.InstanceSet[0].PrivateIpAddresses) > 0 {
				nodeInfo["lan_ip"] = *describeInstancesResponse.Response.InstanceSet[0].PrivateIpAddresses[0]
			}
		}

		nodes = append(nodes, nodeInfo)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("nodes", nodes)
	_ = d.Set("total_count", *response.Response.TotalCount)

	return nil
}
