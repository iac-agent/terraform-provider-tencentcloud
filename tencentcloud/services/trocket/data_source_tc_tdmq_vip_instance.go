package trocket

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctdmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmq"

	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTdmqVipInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqVipInstanceRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "集群 ID",
			},
			// computed
			"cluster_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Cluster information。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "集群 ID",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster 名称",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "创建时间，（毫秒）。",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster 描述 information注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"public_end_point": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public network access 地址",
						},
						"vpc_end_point": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC access 地址",
						},
						"support_namespace_endpoint": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether namespace access points are supported注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"vpcs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "VPC information注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "私有网络 ID",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet Id。",
									},
								},
							},
						},
						"is_vip": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否为a dedicated instance注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"rocket_mq_flag": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Rocketmq cluster identification注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Billing 状态，1 means normal，2 means stopped，3 means destroyed注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"isolate_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Overdue suspension time，in milliseconds注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"http_public_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "HTTP 协议 public network access address注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"http_vpc_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "HTTP 协议 VPC access address注意：此字段可能返回 null，表示无法获取有效值。",
						},
					},
				},
			},
			"instance_config": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Cluster configuration。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_tps_per_namespace": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Single namespace TPS upper 限制",
						},
						"max_namespace_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最大namespaces。",
						},
						"used_namespace_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 used namespaces。",
						},
						"max_topic_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最大topics。",
						},
						"used_topic_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The 数量 topics used。",
						},
						"max_group_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最大groups。",
						},
						"used_group_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 used groups。",
						},
						"config_display": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "集群类型",
						},
						"node_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 cluster nodes。",
						},
						"node_distribution": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Node distribution。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Availability 可用区",
									},
									"zone_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Availability 可用区 ID",
									},
									"node_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "节点数量",
									},
								},
							},
						},
						"topic_distribution": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Topic distribution。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"topic_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Topic 类型",
									},
									"count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 topics。",
									},
								},
							},
						},
						"max_queues_per_topic": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最大queues per topic注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudTdmqVipInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tdmq_vip_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service     = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		vipInstance *tdmq.DescribeRocketMQVipInstanceDetailResponseParams
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTdmqVipInstanceByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil {
			return nil
		}

		vipInstance = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0)
	if vipInstance.ClusterInfo != nil {
		rocketMQClusterInfoMap := map[string]interface{}{}

		if vipInstance.ClusterInfo.ClusterId != nil {
			rocketMQClusterInfoMap["cluster_id"] = vipInstance.ClusterInfo.ClusterId
		}

		if vipInstance.ClusterInfo.ClusterName != nil {
			rocketMQClusterInfoMap["cluster_name"] = vipInstance.ClusterInfo.ClusterName
		}

		if vipInstance.ClusterInfo.Region != nil {
			rocketMQClusterInfoMap["region"] = vipInstance.ClusterInfo.Region
		}

		if vipInstance.ClusterInfo.CreateTime != nil {
			rocketMQClusterInfoMap["create_time"] = vipInstance.ClusterInfo.CreateTime
		}

		if vipInstance.ClusterInfo.Remark != nil {
			rocketMQClusterInfoMap["remark"] = vipInstance.ClusterInfo.Remark
		}

		if vipInstance.ClusterInfo.PublicEndPoint != nil {
			rocketMQClusterInfoMap["public_end_point"] = vipInstance.ClusterInfo.PublicEndPoint
		}

		if vipInstance.ClusterInfo.VpcEndPoint != nil {
			rocketMQClusterInfoMap["vpc_end_point"] = vipInstance.ClusterInfo.VpcEndPoint
		}

		if vipInstance.ClusterInfo.SupportNamespaceEndpoint != nil {
			rocketMQClusterInfoMap["support_namespace_endpoint"] = vipInstance.ClusterInfo.SupportNamespaceEndpoint
		}

		if vipInstance.ClusterInfo.Vpcs != nil {
			vpcsList := []interface{}{}
			for _, vpcs := range vipInstance.ClusterInfo.Vpcs {
				vpcsMap := map[string]interface{}{}

				if vpcs.VpcId != nil {
					vpcsMap["vpc_id"] = vpcs.VpcId
				}

				if vpcs.SubnetId != nil {
					vpcsMap["subnet_id"] = vpcs.SubnetId
				}

				vpcsList = append(vpcsList, vpcsMap)
			}

			rocketMQClusterInfoMap["vpcs"] = []interface{}{vpcsList}
		}

		if vipInstance.ClusterInfo.IsVip != nil {
			rocketMQClusterInfoMap["is_vip"] = vipInstance.ClusterInfo.IsVip
		}

		if vipInstance.ClusterInfo.RocketMQFlag != nil {
			rocketMQClusterInfoMap["rocket_mq_flag"] = vipInstance.ClusterInfo.RocketMQFlag
		}

		if vipInstance.ClusterInfo.Status != nil {
			rocketMQClusterInfoMap["status"] = vipInstance.ClusterInfo.Status
		}

		if vipInstance.ClusterInfo.IsolateTime != nil {
			rocketMQClusterInfoMap["isolate_time"] = vipInstance.ClusterInfo.IsolateTime
		}

		if vipInstance.ClusterInfo.HttpPublicEndpoint != nil {
			rocketMQClusterInfoMap["http_public_endpoint"] = vipInstance.ClusterInfo.HttpPublicEndpoint
		}

		if vipInstance.ClusterInfo.HttpVpcEndpoint != nil {
			rocketMQClusterInfoMap["http_vpc_endpoint"] = vipInstance.ClusterInfo.HttpVpcEndpoint
		}

		ids = append(ids, *vipInstance.ClusterInfo.ClusterId)
		_ = d.Set("cluster_info", rocketMQClusterInfoMap)
	}

	if vipInstance.InstanceConfig != nil {
		rocketMQInstanceConfigMap := map[string]interface{}{}

		if vipInstance.InstanceConfig.MaxTpsPerNamespace != nil {
			rocketMQInstanceConfigMap["max_tps_per_namespace"] = vipInstance.InstanceConfig.MaxTpsPerNamespace
		}

		if vipInstance.InstanceConfig.MaxNamespaceNum != nil {
			rocketMQInstanceConfigMap["max_namespace_num"] = vipInstance.InstanceConfig.MaxNamespaceNum
		}

		if vipInstance.InstanceConfig.UsedNamespaceNum != nil {
			rocketMQInstanceConfigMap["used_namespace_num"] = vipInstance.InstanceConfig.UsedNamespaceNum
		}

		if vipInstance.InstanceConfig.MaxTopicNum != nil {
			rocketMQInstanceConfigMap["max_topic_num"] = vipInstance.InstanceConfig.MaxTopicNum
		}

		if vipInstance.InstanceConfig.UsedTopicNum != nil {
			rocketMQInstanceConfigMap["used_topic_num"] = vipInstance.InstanceConfig.UsedTopicNum
		}

		if vipInstance.InstanceConfig.MaxGroupNum != nil {
			rocketMQInstanceConfigMap["max_group_num"] = vipInstance.InstanceConfig.MaxGroupNum
		}

		if vipInstance.InstanceConfig.UsedGroupNum != nil {
			rocketMQInstanceConfigMap["used_group_num"] = vipInstance.InstanceConfig.UsedGroupNum
		}

		if vipInstance.InstanceConfig.ConfigDisplay != nil {
			rocketMQInstanceConfigMap["config_display"] = vipInstance.InstanceConfig.ConfigDisplay
		}

		if vipInstance.InstanceConfig.NodeCount != nil {
			rocketMQInstanceConfigMap["node_count"] = vipInstance.InstanceConfig.NodeCount
		}

		if vipInstance.InstanceConfig.NodeDistribution != nil {
			nodeDistributionList := []interface{}{}
			for _, nodeDistribution := range vipInstance.InstanceConfig.NodeDistribution {
				nodeDistributionMap := map[string]interface{}{}

				if nodeDistribution.ZoneName != nil {
					nodeDistributionMap["zone_name"] = nodeDistribution.ZoneName
				}

				if nodeDistribution.ZoneId != nil {
					nodeDistributionMap["zone_id"] = nodeDistribution.ZoneId
				}

				if nodeDistribution.NodeCount != nil {
					nodeDistributionMap["node_count"] = nodeDistribution.NodeCount
				}

				nodeDistributionList = append(nodeDistributionList, nodeDistributionMap)
			}

			rocketMQInstanceConfigMap["node_distribution"] = []interface{}{nodeDistributionList}
		}

		if vipInstance.InstanceConfig.TopicDistribution != nil {
			topicDistributionList := []interface{}{}
			for _, topicDistribution := range vipInstance.InstanceConfig.TopicDistribution {
				topicDistributionMap := map[string]interface{}{}

				if topicDistribution.TopicType != nil {
					topicDistributionMap["topic_type"] = topicDistribution.TopicType
				}

				if topicDistribution.Count != nil {
					topicDistributionMap["count"] = topicDistribution.Count
				}

				topicDistributionList = append(topicDistributionList, topicDistributionMap)
			}

			rocketMQInstanceConfigMap["topic_distribution"] = []interface{}{topicDistributionList}
		}

		if vipInstance.InstanceConfig.MaxQueuesPerTopic != nil {
			rocketMQInstanceConfigMap["max_queues_per_topic"] = vipInstance.InstanceConfig.MaxQueuesPerTopic
		}

		_ = d.Set("instance_config", rocketMQInstanceConfigMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
