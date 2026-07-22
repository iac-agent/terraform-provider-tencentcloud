package tsf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTsfCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfClusterRead,
		Schema: map[string]*schema.Schema{
			"cluster_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "集群 ID 列表 到 是 queried，如果未填写 在 或 passed，all 内容 将 是 queried。",
			},

			"cluster_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "类型 集群 到 是 queried，如果 left blank 或 不 passed，all 内容 将 是 queried. C: 容器，V: virtual machine。",
			},

			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "过滤器 通过 keywords 对于 Cluster ID 或 名称",
			},

			"disable_program_auth_check": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "是否disable dataset authentication。",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "TSF 集群 pagination 对象. 注意: 此 字段 可能 返回 null，indicating 无 有效 值",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total 数量 items. 注意：此字段可能返回 null，表示未找到有效值。",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Cluster 列表. 注意: 此 字段 可能 返回 null，indicating 无 有效 值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "集群 ID 注意: 此 字段 可能 返回 null，indicating 无 有效 值",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "集群名称 注意: 此 字段 可能 返回 null，indicating 无 有效 值",
									},
									"cluster_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster 描述 注意: 此 字段 可能 返回 null，indicating 无 有效 值",
									},
									"cluster_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "集群类型 注意: 此 字段 可能 返回 null，indicating 无 有效 值",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Private 网络 ID 集群. 注意: 此 字段 可能 返回 null，indicating 无 有效 值",
									},
									"cluster_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "集群 状态 注意: 此 字段 可能 返回 null，indicating 无 有效 值",
									},
									"cluster_cidr": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "集群 CIDR. 注意: 此 字段 可能 返回 null，indicating 无 有效 值",
									},
									"cluster_total_cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Total CPU 的 集群，单位: cores. 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"cluster_total_mem": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Total 内存 的 集群，单位: G. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
									},
									"cluster_used_cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Used CPU 的 集群，在 cores. 此 字段 可能 返回 null，indicating 无 有效 值",
									},
									"cluster_used_mem": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Cluster 使用 内存 （GB）。 此 字段 可能 返回 null，indicating 无 有效 值",
									},
									"instance_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Cluster 实例 数量. 此 字段 可能 返回 null，indicating 无 有效 值",
									},
									"run_instance_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Cluster running 实例 数量. 此 字段 可能 返回 null，indicating 无 有效 值",
									},
									"normal_instance_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Cluster normal 实例 数量. 此 字段 可能 返回 null，indicating 无 有效 值",
									},
									"delete_flag": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Deletion 标签: true 表示 它 可以 是 删除，false 表示 它 不能 是 删除. 注意: 此 字段 可能 返回 null，indicating 无 有效 值",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CreationTime. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "last 更新时间. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"tsf_region_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "地域 ID TSF. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"tsf_region_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "地域名称 的 TSF. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"tsf_zone_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "可用区 ID 的 TSF. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"tsf_zone_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "可用区 名称 TSF. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"delete_flag_reason": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Reason why 集群 不能 是 删除. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"cluster_limit_cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Maximum CPU 限制 的 集群，在 cores. 此 字段 可能 返回 null，indicating 该 无 有效 值 是 found。",
									},
									"cluster_limit_mem": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Cluster 最大 内存 限制 （GB）。 此 字段 可能 返回 null，indicating 该 无 有效 值 是 found。",
									},
									"run_service_instance_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 可用 服务 实例 在 集群. 注意: 此 字段 可能 返回 null，indicating 无 有效 值",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster 子网 ID. 注意: 此 字段 可能 返回 null，indicating 无 有效 值。",
									},
									"operation_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Control 信息 返回 到 frontend. 此 字段 可能 返回 null，indicating 无 有效 值",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"init": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Control 信息 的 initialization button 返回 到 front end. 注意: 此 字段 可能 返回 null，indicating 无 有效 值",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"disabled_reason": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Reason 对于 不 displaying. 注意: 此 字段 可能 返回 null，indicating 无 有效 值",
															},
															"enabled": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "availability 的 button (whether 它 是 clickable) 可能 返回 null indicating 该 信息 是 不 可用。",
															},
															"supported": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否display button. 注意：此字段可能返回 null，表示未找到有效值。",
															},
														},
													},
												},
												"add_instance": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Add 实例 button control 信息，注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"disabled_reason": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "reason why 此 button 是 不 displayed，可能 返回 null 如果 不 applicable。",
															},
															"enabled": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否button 是 clickable. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
															},
															"supported": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否button 是 clickable. 注意：此字段可能返回 null，表示未找到有效值。",
															},
														},
													},
												},
												"destroy": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Control 信息 对于 destroying machine，可能 返回 null 如果 无 有效 值 是 获取。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"disabled_reason": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "reason why 此 button 是 不 displayed，可能 返回 null 如果 不 applicable。",
															},
															"enabled": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否button 是 clickable. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 获取。",
															},
															"supported": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否button 是 clickable. 注意：此字段可能返回 null，表示未找到有效值。",
															},
														},
													},
												},
											},
										},
									},
									"cluster_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "集群 版本，可能 返回 null 如果 不 applicable。",
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

func dataSourceTencentCloudTsfClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tsf_cluster.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id_list"); ok {
		clusterIdListSet := v.(*schema.Set).List()
		paramMap["ClusterIdList"] = helper.InterfacesStringsPoint(clusterIdListSet)
	}

	if v, ok := d.GetOk("cluster_type"); ok {
		paramMap["ClusterType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("disable_program_auth_check"); v != nil {
		paramMap["DisableProgramAuthCheck"] = helper.Bool(v.(bool))
	}

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var cluster *tsf.TsfPageCluster
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfClusterByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		cluster = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(cluster.Content))
	tsfPageClusterMap := map[string]interface{}{}
	if cluster != nil {
		if cluster.TotalCount != nil {
			tsfPageClusterMap["total_count"] = cluster.TotalCount
		}

		if cluster.Content != nil {
			contentList := []interface{}{}
			for _, content := range cluster.Content {
				contentMap := map[string]interface{}{}

				if content.ClusterId != nil {
					contentMap["cluster_id"] = content.ClusterId
				}

				if content.ClusterName != nil {
					contentMap["cluster_name"] = content.ClusterName
				}

				if content.ClusterDesc != nil {
					contentMap["cluster_desc"] = content.ClusterDesc
				}

				if content.ClusterType != nil {
					contentMap["cluster_type"] = content.ClusterType
				}

				if content.VpcId != nil {
					contentMap["vpc_id"] = content.VpcId
				}

				if content.ClusterStatus != nil {
					contentMap["cluster_status"] = content.ClusterStatus
				}

				if content.ClusterCIDR != nil {
					contentMap["cluster_cidr"] = content.ClusterCIDR
				}

				if content.ClusterTotalCpu != nil {
					contentMap["cluster_total_cpu"] = content.ClusterTotalCpu
				}

				if content.ClusterTotalMem != nil {
					contentMap["cluster_total_mem"] = content.ClusterTotalMem
				}

				if content.ClusterUsedCpu != nil {
					contentMap["cluster_used_cpu"] = content.ClusterUsedCpu
				}

				if content.ClusterUsedMem != nil {
					contentMap["cluster_used_mem"] = content.ClusterUsedMem
				}

				if content.InstanceCount != nil {
					contentMap["instance_count"] = content.InstanceCount
				}

				if content.RunInstanceCount != nil {
					contentMap["run_instance_count"] = content.RunInstanceCount
				}

				if content.NormalInstanceCount != nil {
					contentMap["normal_instance_count"] = content.NormalInstanceCount
				}

				if content.DeleteFlag != nil {
					contentMap["delete_flag"] = content.DeleteFlag
				}

				if content.CreateTime != nil {
					contentMap["create_time"] = content.CreateTime
				}

				if content.UpdateTime != nil {
					contentMap["update_time"] = content.UpdateTime
				}

				if content.TsfRegionId != nil {
					contentMap["tsf_region_id"] = content.TsfRegionId
				}

				if content.TsfRegionName != nil {
					contentMap["tsf_region_name"] = content.TsfRegionName
				}

				if content.TsfZoneId != nil {
					contentMap["tsf_zone_id"] = content.TsfZoneId
				}

				if content.TsfZoneName != nil {
					contentMap["tsf_zone_name"] = content.TsfZoneName
				}

				if content.DeleteFlagReason != nil {
					contentMap["delete_flag_reason"] = content.DeleteFlagReason
				}

				if content.ClusterLimitCpu != nil {
					contentMap["cluster_limit_cpu"] = content.ClusterLimitCpu
				}

				if content.ClusterLimitMem != nil {
					contentMap["cluster_limit_mem"] = content.ClusterLimitMem
				}

				if content.RunServiceInstanceCount != nil {
					contentMap["run_service_instance_count"] = content.RunServiceInstanceCount
				}

				if content.SubnetId != nil {
					contentMap["subnet_id"] = content.SubnetId
				}

				if content.OperationInfo != nil {
					operationInfoMap := map[string]interface{}{}

					if content.OperationInfo.Init != nil {
						initMap := map[string]interface{}{}

						if content.OperationInfo.Init.DisabledReason != nil {
							initMap["disabled_reason"] = content.OperationInfo.Init.DisabledReason
						}

						if content.OperationInfo.Init.Enabled != nil {
							initMap["enabled"] = content.OperationInfo.Init.Enabled
						}

						if content.OperationInfo.Init.Supported != nil {
							initMap["supported"] = content.OperationInfo.Init.Supported
						}

						operationInfoMap["init"] = []interface{}{initMap}
					}

					if content.OperationInfo.AddInstance != nil {
						addInstanceMap := map[string]interface{}{}

						if content.OperationInfo.AddInstance.DisabledReason != nil {
							addInstanceMap["disabled_reason"] = content.OperationInfo.AddInstance.DisabledReason
						}

						if content.OperationInfo.AddInstance.Enabled != nil {
							addInstanceMap["enabled"] = content.OperationInfo.AddInstance.Enabled
						}

						if content.OperationInfo.AddInstance.Supported != nil {
							addInstanceMap["supported"] = content.OperationInfo.AddInstance.Supported
						}

						operationInfoMap["add_instance"] = []interface{}{addInstanceMap}
					}

					if content.OperationInfo.Destroy != nil {
						destroyMap := map[string]interface{}{}

						if content.OperationInfo.Destroy.DisabledReason != nil {
							destroyMap["disabled_reason"] = content.OperationInfo.Destroy.DisabledReason
						}

						if content.OperationInfo.Destroy.Enabled != nil {
							destroyMap["enabled"] = content.OperationInfo.Destroy.Enabled
						}

						if content.OperationInfo.Destroy.Supported != nil {
							destroyMap["supported"] = content.OperationInfo.Destroy.Supported
						}

						operationInfoMap["destroy"] = []interface{}{destroyMap}
					}

					contentMap["operation_info"] = []interface{}{operationInfoMap}
				}

				if content.ClusterVersion != nil {
					contentMap["cluster_version"] = content.ClusterVersion
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.ClusterId)
			}

			tsfPageClusterMap["content"] = contentList
		}

		_ = d.Set("result", []interface{}{tsfPageClusterMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tsfPageClusterMap); e != nil {
			return e
		}
	}
	return nil
}
