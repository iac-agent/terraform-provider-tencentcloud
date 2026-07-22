package cdwch

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clickhouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClickhouseInstanceNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClickhouseInstanceNodesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID",
			},

			"node_role": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Cluster 角色 类型，默认为 `数据` 数据 节点。",
			},

			"display_policy": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Display strategy，display all 当 All。",
			},

			"force_all": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "当 true，返回all nodes，该 是， 限制 是 infinitely large。",
			},

			"instance_nodes_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Total 数量 实例 nodes。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP 地址",
						},
						"spec": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Model，such 作为 S1。",
						},
						"core": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CPU 核数",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory 大小。",
						},
						"disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk 类型",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk 大小。",
						},
						"cluster": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 clickhouse 集群 到 其中 它 belongs。",
						},
						"node_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Group 信息 到 其中 节点 belongs。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group 名称",
									},
									"shard_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Fragmented variable 名称",
									},
									"replica_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Copy variable 名称",
									},
								},
							},
						},
						"rip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC IP。",
						},
						"is_ch_proxy": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "当 true，它 表示that chproxy process has been deployed 在 节点。",
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

func dataSourceTencentCloudClickhouseInstanceNodesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clickhouse_instance_nodes.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("node_role"); ok {
		paramMap["NodeRole"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("display_policy"); ok {
		paramMap["DisplayPolicy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("force_all"); ok {
		paramMap["ForceAll"] = helper.Bool(v.(bool))
	}

	service := CdwchService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var instanceNodesList []*clickhouse.InstanceNode
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClickhouseInstanceNodesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instanceNodesList = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(instanceNodesList))

	if instanceNodesList != nil {
		for _, instanceNode := range instanceNodesList {
			instanceNodeMap := map[string]interface{}{}

			if instanceNode.Ip != nil {
				instanceNodeMap["ip"] = instanceNode.Ip
			}

			if instanceNode.Spec != nil {
				instanceNodeMap["spec"] = instanceNode.Spec
			}

			if instanceNode.Core != nil {
				instanceNodeMap["core"] = instanceNode.Core
			}

			if instanceNode.Memory != nil {
				instanceNodeMap["memory"] = instanceNode.Memory
			}

			if instanceNode.DiskType != nil {
				instanceNodeMap["disk_type"] = instanceNode.DiskType
			}

			if instanceNode.DiskSize != nil {
				instanceNodeMap["disk_size"] = instanceNode.DiskSize
			}

			if instanceNode.Cluster != nil {
				instanceNodeMap["cluster"] = instanceNode.Cluster
			}

			if instanceNode.NodeGroups != nil {
				var nodeGroupsList []interface{}
				for _, nodeGroups := range instanceNode.NodeGroups {
					nodeGroupsMap := map[string]interface{}{}

					if nodeGroups.GroupName != nil {
						nodeGroupsMap["group_name"] = nodeGroups.GroupName
					}

					if nodeGroups.ShardName != nil {
						nodeGroupsMap["shard_name"] = nodeGroups.ShardName
					}

					if nodeGroups.ReplicaName != nil {
						nodeGroupsMap["replica_name"] = nodeGroups.ReplicaName
					}

					nodeGroupsList = append(nodeGroupsList, nodeGroupsMap)
				}

				instanceNodeMap["node_groups"] = nodeGroupsList
			}

			if instanceNode.Rip != nil {
				instanceNodeMap["rip"] = instanceNode.Rip
			}

			if instanceNode.IsCHProxy != nil {
				instanceNodeMap["is_ch_proxy"] = instanceNode.IsCHProxy
			}

			tmpList = append(tmpList, instanceNodeMap)
		}

		_ = d.Set("instance_nodes_list", tmpList)
	}

	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
