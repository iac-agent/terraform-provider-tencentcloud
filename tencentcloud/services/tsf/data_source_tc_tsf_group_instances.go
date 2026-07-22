package tsf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTsfGroupInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfGroupInstancesRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "group id。",
			},

			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "search word。",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "顺序 term。",
			},

			"order_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "顺序 类型",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Machine information of the deployment group.注意：此字段可能返回 null，表示未找到有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total 数量 machine instances.注意：此字段可能返回 null，表示未找到有效值。",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "列表 machine instances.注意：此字段可能返回 null，表示未找到有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Machine instance ID.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Machine 名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"lan_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "内网 IP 地址注意：此字段可能返回 null，表示未找到有效值。",
									},
									"wan_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "公网 IP 地址注意：此字段可能返回 null，表示未找到有效值。",
									},
									"instance_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "描述注意：此字段可能返回 null，表示未找到有效值。",
									},
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "集群 ID注意：此字段可能返回 null，表示未找到有效值。",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "集群名称 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"instance_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VM 状态 For virtual machines，it 表示status of the virtual machine. For containers，it 表示status of the virtual machine where the pod is located.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"instance_available_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VM availability 状态 For virtual machines，it 表示是否virtual machine can be used as a resource. For containers，it 表示是否virtual machine can be 用于deploy pods.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"service_instance_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "状态 service instances under the service. For virtual machines，it 表示是否application is available and the agent 状态 For containers，it 表示status of the pod.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"count_in_tsf": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "表示是否this instance has been added to the TSF.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group id.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"application_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application id.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"application_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application 名称 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"instance_created_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创建时间 of the machine instance in CVM.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"instance_expired_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Expire time of the machine instance in CVM.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"instance_charge_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "machine 实例计费类型注意：此字段可能返回 null，表示未找到有效值。",
									},
									"instance_total_cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Total CPU information of the machine instance.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"instance_total_mem": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Total memory information of the machine instance.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"instance_used_cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "CPU information used by the machine instance.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"instance_used_mem": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Memory information used by the machine instance.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"instance_limit_cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "限制 CPU information of the machine instance.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"instance_limit_mem": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "限制 memory information of the machine instance.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"instance_pkg_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "instance pkg 版本注意：此字段可能返回 null，表示未找到有效值。",
									},
									"cluster_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "集群类型注意：此字段可能返回 null，表示未找到有效值。",
									},
									"restrict_state": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Business 状态 machine instance.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "更新时间.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"operation_state": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Execution 状态 instance.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"namespace_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace id.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"instance_zone_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance 可用区 ID注意：此字段可能返回 null，表示未找到有效值。",
									},
									"instance_import_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "InstanceImportMode import 模式注意：此字段可能返回 null，表示未找到有效值。",
									},
									"application_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application id.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"application_resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "application 资源 ID注意：此字段可能返回 null，表示未找到有效值。",
									},
									"service_sidecar_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Sidecar 状态注意：此字段可能返回 null，表示未找到有效值。",
									},
									"group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "组名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"namespace_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace 名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"reason": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Health checking reason.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"agent_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Agent 版本注意：此字段可能返回 null，表示未找到有效值。",
									},
									"node_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Container 主机 instance ID.注意：此字段可能返回 null，表示未找到有效值。",
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

func dataSourceTencentCloudTsfGroupInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tsf_group_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("group_id"); ok {
		paramMap["GroupId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("order_type"); v != nil {
		paramMap["OrderType"] = helper.IntInt64(v.(int))
	}

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var result *tsf.TsfPageInstance
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeTsfGroupInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result.Content))
	tsfPageInstanceMap := map[string]interface{}{}
	if result != nil {
		if result.TotalCount != nil {
			tsfPageInstanceMap["total_count"] = result.TotalCount
		}

		if result.Content != nil {
			contentList := []interface{}{}
			for _, content := range result.Content {
				contentMap := map[string]interface{}{}

				if content.InstanceId != nil {
					contentMap["instance_id"] = content.InstanceId
				}

				if content.InstanceName != nil {
					contentMap["instance_name"] = content.InstanceName
				}

				if content.LanIp != nil {
					contentMap["lan_ip"] = content.LanIp
				}

				if content.WanIp != nil {
					contentMap["wan_ip"] = content.WanIp
				}

				if content.InstanceDesc != nil {
					contentMap["instance_desc"] = content.InstanceDesc
				}

				if content.ClusterId != nil {
					contentMap["cluster_id"] = content.ClusterId
				}

				if content.ClusterName != nil {
					contentMap["cluster_name"] = content.ClusterName
				}

				if content.InstanceStatus != nil {
					contentMap["instance_status"] = content.InstanceStatus
				}

				if content.InstanceAvailableStatus != nil {
					contentMap["instance_available_status"] = content.InstanceAvailableStatus
				}

				if content.ServiceInstanceStatus != nil {
					contentMap["service_instance_status"] = content.ServiceInstanceStatus
				}

				if content.CountInTsf != nil {
					contentMap["count_in_tsf"] = content.CountInTsf
				}

				if content.GroupId != nil {
					contentMap["group_id"] = content.GroupId
				}

				if content.ApplicationId != nil {
					contentMap["application_id"] = content.ApplicationId
				}

				if content.ApplicationName != nil {
					contentMap["application_name"] = content.ApplicationName
				}

				if content.InstanceCreatedTime != nil {
					contentMap["instance_created_time"] = content.InstanceCreatedTime
				}

				if content.InstanceExpiredTime != nil {
					contentMap["instance_expired_time"] = content.InstanceExpiredTime
				}

				if content.InstanceChargeType != nil {
					contentMap["instance_charge_type"] = content.InstanceChargeType
				}

				if content.InstanceTotalCpu != nil {
					contentMap["instance_total_cpu"] = content.InstanceTotalCpu
				}

				if content.InstanceTotalMem != nil {
					contentMap["instance_total_mem"] = content.InstanceTotalMem
				}

				if content.InstanceUsedCpu != nil {
					contentMap["instance_used_cpu"] = content.InstanceUsedCpu
				}

				if content.InstanceUsedMem != nil {
					contentMap["instance_used_mem"] = content.InstanceUsedMem
				}

				if content.InstanceLimitCpu != nil {
					contentMap["instance_limit_cpu"] = content.InstanceLimitCpu
				}

				if content.InstanceLimitMem != nil {
					contentMap["instance_limit_mem"] = content.InstanceLimitMem
				}

				if content.InstancePkgVersion != nil {
					contentMap["instance_pkg_version"] = content.InstancePkgVersion
				}

				if content.ClusterType != nil {
					contentMap["cluster_type"] = content.ClusterType
				}

				if content.RestrictState != nil {
					contentMap["restrict_state"] = content.RestrictState
				}

				if content.UpdateTime != nil {
					contentMap["update_time"] = content.UpdateTime
				}

				if content.OperationState != nil {
					contentMap["operation_state"] = content.OperationState
				}

				if content.NamespaceId != nil {
					contentMap["namespace_id"] = content.NamespaceId
				}

				if content.InstanceZoneId != nil {
					contentMap["instance_zone_id"] = content.InstanceZoneId
				}

				if content.InstanceImportMode != nil {
					contentMap["instance_import_mode"] = content.InstanceImportMode
				}

				if content.ApplicationType != nil {
					contentMap["application_type"] = content.ApplicationType
				}

				if content.ApplicationResourceType != nil {
					contentMap["application_resource_type"] = content.ApplicationResourceType
				}

				if content.ServiceSidecarStatus != nil {
					contentMap["service_sidecar_status"] = content.ServiceSidecarStatus
				}

				if content.GroupName != nil {
					contentMap["group_name"] = content.GroupName
				}

				if content.NamespaceName != nil {
					contentMap["namespace_name"] = content.NamespaceName
				}

				if content.Reason != nil {
					contentMap["reason"] = content.Reason
				}

				if content.AgentVersion != nil {
					contentMap["agent_version"] = content.AgentVersion
				}

				if content.NodeInstanceId != nil {
					contentMap["node_instance_id"] = content.NodeInstanceId
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.InstanceId)
			}

			tsfPageInstanceMap["content"] = contentList
		}

		_ = d.Set("result", []interface{}{tsfPageInstanceMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tsfPageInstanceMap); e != nil {
			return e
		}
	}
	return nil
}
