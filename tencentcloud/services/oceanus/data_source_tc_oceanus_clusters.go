package oceanus

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudOceanusClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOceanusClustersRead,
		Schema: map[string]*schema.Schema{
			"cluster_ids": {
				Optional:    true,
				Type:        schema.TypeSet,
				MaxItems:    100,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Query 一个 或 more clusters 通过 their ID. 最大clusters 该 可以 是 queried 在 once 是 100。",
			},
			"order_type": {
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateAllowedIntValue(CLUSTER_ORDER_TYPE),
				Description:  "sorting 规则 的 集群 信息 results. Possible 值 是 1 (排序方式 时间 在 降序)，2 (排序方式 时间 在 升序)，和 3 (排序方式 状态)。",
			},
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "filtering 规则。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "字段 到 是 filtered。",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "filtering 值 的 字段。",
						},
					},
				},
			},
			"work_space_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Workspace SerialId。",
			},
			"cluster_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Cluster 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 集群。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 集群。",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 其中 集群 是 located。",
						},
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "用户 AppID。",
						},
						"owner_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "main 账号 UIN。",
						},
						"creator_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建者 UIN。",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "状态 集群. Possible 值 是 1 (uninitialized)，3 (initializing)，和 2 (running)。",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A 描述 集群。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "时间 当 集群 是 创建。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "时间 的 last operation 在 集群。",
						},
						"cu_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 CUs。",
						},
						"cu_mem": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "内存 规格 的 CU。",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "availability 可用区",
						},
						"status_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态 描述",
						},
						"ccns": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "网络。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID VPC。",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID 子网。",
									},
									"ccn_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID Cloud Connect Network (CCN)，such 作为 ccn-rahigzjd。",
									},
								},
							},
						},
						"net_environment_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "网络。",
						},
						"free_cu_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 free CUs。",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "标签 bound 到 集群.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签键注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签值注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"isolated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "时间 当 集群 是 isolated. 如果 集群 has 不 been isolated，此 字段 将 是 -.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "过期时间 的 集群. 如果 集群 does 不 have 过期时间，此 字段 将 是 -.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"seconds_until_expiry": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "数量 秒 until 集群 expires. 如果 集群 does 不 have 过期时间，此 字段 将 是 -.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "auto-renewal flag. 0 表示default state ( 用户 has 不 集合 它，其中 是 initial state; 如果 用户 has 已启用 prepaid non-stop privilege， 集群 将 是 automatically renewed)，1 表示automatic renewal，和 2 表示no automatic renewal (集合 通过 用户).注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"default_cos_bucket": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "默认值 COS 存储桶 的 集群.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"cls_log_set": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLS logset 的 集群.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"cls_topic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLS 主题 ID 集群.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"cls_log_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 CLS logset 的 集群.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"cls_topic_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 CLS 主题 的 集群.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"version": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "版本 信息 的 集群.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"flink": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Flink 版本 的 集群.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"supported_flink": {
										Type:        schema.TypeSet,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Computed:    true,
										Description: "Flink versions 支持 通过 集群.注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"free_cu": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "数量 free CUs 在 granularity 级别注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"default_log_collect_conf": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "默认值 日志 collection 配置 的 集群.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"customized_dns_enabled": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "值: 0 - 不 集合，1 - 集合，2 - 不 allowed 到 集合.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"correlations": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Space 信息.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_group_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "集群 ID",
									},
									"cluster_group_serial_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster SerialId。",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "集群名称",
									},
									"work_space_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Workspace SerialId。",
									},
									"work_space_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Workspace 名称",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Binding 状态 2 - bound，1 - unbound。",
									},
									"project_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "项目 ID",
									},
									"project_id_str": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "项目 ID 在 字符串 格式注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"running_cu": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Running CU.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"pay_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "0 - postpaid，1 - prepaid.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"is_need_manage_node": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Front-end distinguishes 是否cluster needs 2CU logic，because historical clusters do 不 need 到 是 changed. 默认为 1. All new clusters need 到 是 changed.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"cluster_sessions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Session 集群 信息.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{},
							},
						},
						"arch_generation": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "V3 版本 = 2.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"cluster_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "0: TKE，1: EKS.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"orders": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "顺序 信息.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "1 - create，2 - renew，3 - scale.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"auto_renew_flag": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "1 - auto-renewal.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"operate_uin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "UIN 的 操作者注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"compute_cu": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 CUs 在 final 集群.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"order_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "时间 的 顺序注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"sql_gateways": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Gateway 信息.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"serial_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Unique identifier.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"flink_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Flink kernel 版本注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "状态 1 - stopped，2 - starting，3 - started，4 - start failed，5 - stopping.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"creator_uin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创建者注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"resource_refs": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Reference resources.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"workspace_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Unique identifier 的 space。",
												},
												"resource_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Unique identifier 的 资源。",
												},
												"version": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "版本 数量。",
												},
												"type": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Reference 类型 0: 用户 资源.注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"cu_spec": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "CU 规格.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创建时间.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "更新时间.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"properties": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Configuration 参数.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "键 的 系统 配置。",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "值 的 系统 配置。",
												},
											},
										},
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

func dataSourceTencentCloudOceanusClustersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_oceanus_clusters.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = OceanusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		clusterSet []*oceanus.Cluster
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_ids"); ok {
		clusterIdsSet := v.(*schema.Set).List()
		paramMap["ClusterIds"] = helper.InterfacesStringsPoint(clusterIdsSet)
	}

	if v, ok := d.GetOkExists("order_type"); ok {
		paramMap["OrderType"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*oceanus.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := oceanus.Filter{}
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

	if v, ok := d.GetOk("work_space_id"); ok {
		paramMap["WorkSpaceId"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOceanusClustersByFilter(ctx, paramMap)
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

			if cluster.Name != nil {
				clusterMap["name"] = cluster.Name
			}

			if cluster.Region != nil {
				clusterMap["region"] = cluster.Region
			}

			if cluster.AppId != nil {
				clusterMap["app_id"] = cluster.AppId
			}

			if cluster.OwnerUin != nil {
				clusterMap["owner_uin"] = cluster.OwnerUin
			}

			if cluster.CreatorUin != nil {
				clusterMap["creator_uin"] = cluster.CreatorUin
			}

			if cluster.Status != nil {
				clusterMap["status"] = cluster.Status
			}

			if cluster.Remark != nil {
				clusterMap["remark"] = cluster.Remark
			}

			if cluster.CreateTime != nil {
				clusterMap["create_time"] = cluster.CreateTime
			}

			if cluster.UpdateTime != nil {
				clusterMap["update_time"] = cluster.UpdateTime
			}

			if cluster.CuNum != nil {
				clusterMap["cu_num"] = cluster.CuNum
			}

			if cluster.CuMem != nil {
				clusterMap["cu_mem"] = cluster.CuMem
			}

			if cluster.Zone != nil {
				clusterMap["zone"] = cluster.Zone
			}

			if cluster.StatusDesc != nil {
				clusterMap["status_desc"] = cluster.StatusDesc
			}

			if cluster.CCNs != nil {
				CCNsList := []interface{}{}
				for _, cCNs := range cluster.CCNs {
					cCNsMap := map[string]interface{}{}

					if cCNs.VpcId != nil {
						cCNsMap["vpc_id"] = cCNs.VpcId
					}

					if cCNs.SubnetId != nil {
						cCNsMap["subnet_id"] = cCNs.SubnetId
					}

					if cCNs.CcnId != nil {
						cCNsMap["ccn_id"] = cCNs.CcnId
					}

					CCNsList = append(CCNsList, cCNsMap)
				}

				clusterMap["ccns"] = CCNsList
			}

			if cluster.NetEnvironmentType != nil {
				clusterMap["net_environment_type"] = cluster.NetEnvironmentType
			}

			if cluster.FreeCuNum != nil {
				clusterMap["free_cu_num"] = cluster.FreeCuNum
			}

			if cluster.Tags != nil {
				tagsList := []interface{}{}
				for _, tags := range cluster.Tags {
					tagsMap := map[string]interface{}{}

					if tags.TagKey != nil {
						tagsMap["tag_key"] = tags.TagKey
					}

					if tags.TagValue != nil {
						tagsMap["tag_value"] = tags.TagValue
					}

					tagsList = append(tagsList, tagsMap)
				}

				clusterMap["tags"] = tagsList
			}

			if cluster.IsolatedTime != nil {
				clusterMap["isolated_time"] = cluster.IsolatedTime
			}

			if cluster.ExpireTime != nil {
				clusterMap["expire_time"] = cluster.ExpireTime
			}

			if cluster.SecondsUntilExpiry != nil {
				clusterMap["seconds_until_expiry"] = cluster.SecondsUntilExpiry
			}

			if cluster.AutoRenewFlag != nil {
				clusterMap["auto_renew_flag"] = cluster.AutoRenewFlag
			}

			if cluster.DefaultCOSBucket != nil {
				clusterMap["default_cos_bucket"] = cluster.DefaultCOSBucket
			}

			if cluster.CLSLogSet != nil {
				clusterMap["cls_log_set"] = cluster.CLSLogSet
			}

			if cluster.CLSTopicId != nil {
				clusterMap["cls_topic_id"] = cluster.CLSTopicId
			}

			if cluster.CLSLogName != nil {
				clusterMap["cls_log_name"] = cluster.CLSLogName
			}

			if cluster.CLSTopicName != nil {
				clusterMap["cls_topic_name"] = cluster.CLSTopicName
			}

			if cluster.Version != nil {
				versionMap := map[string]interface{}{}

				if cluster.Version.Flink != nil {
					versionMap["flink"] = cluster.Version.Flink
				}

				if cluster.Version.SupportedFlink != nil {
					versionMap["supported_flink"] = cluster.Version.SupportedFlink
				}

				clusterMap["version"] = []interface{}{versionMap}
			}

			if cluster.FreeCu != nil {
				clusterMap["free_cu"] = cluster.FreeCu
			}

			if cluster.DefaultLogCollectConf != nil {
				clusterMap["default_log_collect_conf"] = cluster.DefaultLogCollectConf
			}

			if cluster.CustomizedDNSEnabled != nil {
				clusterMap["customized_dns_enabled"] = cluster.CustomizedDNSEnabled
			}

			if cluster.Correlations != nil {
				correlationsList := []interface{}{}
				for _, correlations := range cluster.Correlations {
					correlationsMap := map[string]interface{}{}

					if correlations.ClusterGroupId != nil {
						correlationsMap["cluster_group_id"] = correlations.ClusterGroupId
					}

					if correlations.ClusterGroupSerialId != nil {
						correlationsMap["cluster_group_serial_id"] = correlations.ClusterGroupSerialId
					}

					if correlations.ClusterName != nil {
						correlationsMap["cluster_name"] = correlations.ClusterName
					}

					if correlations.WorkSpaceId != nil {
						correlationsMap["work_space_id"] = correlations.WorkSpaceId
					}

					if correlations.WorkSpaceName != nil {
						correlationsMap["work_space_name"] = correlations.WorkSpaceName
					}

					if correlations.Status != nil {
						correlationsMap["status"] = correlations.Status
					}

					if correlations.ProjectId != nil {
						correlationsMap["project_id"] = correlations.ProjectId
					}

					if correlations.ProjectIdStr != nil {
						correlationsMap["project_id_str"] = correlations.ProjectIdStr
					}

					correlationsList = append(correlationsList, correlationsMap)
				}

				clusterMap["correlations"] = correlationsList
			}

			if cluster.RunningCu != nil {
				clusterMap["running_cu"] = cluster.RunningCu
			}

			if cluster.PayMode != nil {
				clusterMap["pay_mode"] = cluster.PayMode
			}

			if cluster.IsNeedManageNode != nil {
				clusterMap["is_need_manage_node"] = cluster.IsNeedManageNode
			}

			if cluster.ClusterSessions != nil {
				tmpList = make([]map[string]interface{}, 0, len(cluster.ClusterSessions))
				//for _, item := range cluster.ClusterSessions {
				//	sessionMap := map[string]interface{}{}
				//	if item != nil {
				//
				//	}
				//
				//	tmpList = append(tmpList, sessionMap)
				//}

				clusterMap["cluster_sessions"] = tmpList
			}

			if cluster.ArchGeneration != nil {
				clusterMap["arch_generation"] = cluster.ArchGeneration
			}

			if cluster.ClusterType != nil {
				clusterMap["cluster_type"] = cluster.ClusterType
			}

			if cluster.Orders != nil {
				ordersList := []interface{}{}
				for _, orders := range cluster.Orders {
					ordersMap := map[string]interface{}{}

					if orders.Type != nil {
						ordersMap["type"] = orders.Type
					}

					if orders.AutoRenewFlag != nil {
						ordersMap["auto_renew_flag"] = orders.AutoRenewFlag
					}

					if orders.OperateUin != nil {
						ordersMap["operate_uin"] = orders.OperateUin
					}

					if orders.ComputeCu != nil {
						ordersMap["compute_cu"] = orders.ComputeCu
					}

					if orders.OrderTime != nil {
						ordersMap["order_time"] = orders.OrderTime
					}

					ordersList = append(ordersList, ordersMap)
				}

				clusterMap["orders"] = ordersList
			}

			if cluster.SqlGateways != nil {
				sqlGatewaysList := []interface{}{}
				for _, sqlGateways := range cluster.SqlGateways {
					sqlGatewaysMap := map[string]interface{}{}

					if sqlGateways.SerialId != nil {
						sqlGatewaysMap["serial_id"] = sqlGateways.SerialId
					}

					if sqlGateways.FlinkVersion != nil {
						sqlGatewaysMap["flink_version"] = sqlGateways.FlinkVersion
					}

					if sqlGateways.Status != nil {
						sqlGatewaysMap["status"] = sqlGateways.Status
					}

					if sqlGateways.CreatorUin != nil {
						sqlGatewaysMap["creator_uin"] = sqlGateways.CreatorUin
					}

					if sqlGateways.ResourceRefs != nil {
						resourceRefsList := []interface{}{}
						for _, resourceRefs := range sqlGateways.ResourceRefs {
							resourceRefsMap := map[string]interface{}{}

							if resourceRefs.WorkspaceId != nil {
								resourceRefsMap["workspace_id"] = resourceRefs.WorkspaceId
							}

							if resourceRefs.ResourceId != nil {
								resourceRefsMap["resource_id"] = resourceRefs.ResourceId
							}

							if resourceRefs.Version != nil {
								resourceRefsMap["version"] = resourceRefs.Version
							}

							if resourceRefs.Type != nil {
								resourceRefsMap["type"] = resourceRefs.Type
							}

							resourceRefsList = append(resourceRefsList, resourceRefsMap)
						}

						sqlGatewaysMap["resource_refs"] = resourceRefsList
					}

					if sqlGateways.CuSpec != nil {
						sqlGatewaysMap["cu_spec"] = sqlGateways.CuSpec
					}

					if sqlGateways.CreateTime != nil {
						sqlGatewaysMap["create_time"] = sqlGateways.CreateTime
					}

					if sqlGateways.UpdateTime != nil {
						sqlGatewaysMap["update_time"] = sqlGateways.UpdateTime
					}

					if sqlGateways.Properties != nil {
						propertiesList := []interface{}{}
						for _, properties := range sqlGateways.Properties {
							propertiesMap := map[string]interface{}{}

							if properties.Key != nil {
								propertiesMap["key"] = properties.Key
							}

							if properties.Value != nil {
								propertiesMap["value"] = properties.Value
							}

							propertiesList = append(propertiesList, propertiesMap)
						}

						sqlGatewaysMap["properties"] = propertiesList
					}

					sqlGatewaysList = append(sqlGatewaysList, sqlGatewaysMap)
				}

				clusterMap["sql_gateways"] = sqlGatewaysList
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
