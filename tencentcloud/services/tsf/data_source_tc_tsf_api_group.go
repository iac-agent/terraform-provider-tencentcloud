package tsf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTsfApiGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfApiGroupRead,
		Schema: map[string]*schema.Schema{
			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "search word。",
			},

			"group_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Group 类型 ms: Microservice 组; 外部: External API 组。",
			},

			"auth_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Authentication 类型 secret: Secret 键 authentication; none: No authentication。",
			},

			"status": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Publishing 状态 drafted: Not published. released: Published。",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting 字段: created_time 或 group_context。",
			},

			"order_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Sorting 类型: 0 (ASC) 或 1 (DESC)。",
			},

			"gateway_instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Gateway 实例 ID",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Pagination structure.注意：此字段可能返回 null，表示无法获取有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "记录 count。",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Api 组 info。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Api Group ID.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Api Group 名称注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"group_context": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Api Group Context.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"auth_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authentication 类型 secret: 键 authentication; none: 无 authentication.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Release 状态 drafted: 不 released. released: released.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"created_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group 创建时间.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"updated_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group 创建时间，such 作为: 2019-06-20 15:51:28.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"binded_gateway_deploy_groups": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "网关 组 bind 使用 api 组 列表。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"deploy_group_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Gateway 部署 组 bound 到 API 组.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"deploy_group_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Deploy 组名称注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"application_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Application ID.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"application_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Application 名称注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"application_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Application 名称注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"group_status": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Application category: V: virtual machine 应用，C: 容器 应用. 注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"cluster_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "集群类型，C: 容器，V: virtual machine.注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"api_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "api count。",
									},
									"acl_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "数量 APIs.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "描述注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"group_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group 类型注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"gateway_instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway 实例 类型注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"gateway_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway 实例 ID注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"namespace_name_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace 名称 键注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"service_name_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "键 值 的 microservice 名称 参数.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"namespace_name_key_position": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace 参数 location，路径，头部，或 查询，默认为 路径 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"service_name_key_position": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Microservice 名称 参数 location，路径，头部，或 查询，默认为 路径注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudTsfApiGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tsf_api_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_type"); ok {
		paramMap["GroupType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("auth_type"); ok {
		paramMap["AuthType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("status"); ok {
		paramMap["Status"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("order_type"); v != nil {
		paramMap["OrderType"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("gateway_instance_id"); ok {
		paramMap["GatewayInstanceId"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var apiGroupInfo *tsf.TsfPageApiGroupInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfApiGroupByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		apiGroupInfo = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(apiGroupInfo.Content))
	tsfPageApiGroupInfoMap := map[string]interface{}{}
	if apiGroupInfo != nil {

		if apiGroupInfo.TotalCount != nil {
			tsfPageApiGroupInfoMap["total_count"] = apiGroupInfo.TotalCount
		}

		if apiGroupInfo.Content != nil {
			contentList := []interface{}{}
			for _, content := range apiGroupInfo.Content {
				contentMap := map[string]interface{}{}

				if content.GroupId != nil {
					contentMap["group_id"] = content.GroupId
				}

				if content.GroupName != nil {
					contentMap["group_name"] = content.GroupName
				}

				if content.GroupContext != nil {
					contentMap["group_context"] = content.GroupContext
				}

				if content.AuthType != nil {
					contentMap["auth_type"] = content.AuthType
				}

				if content.Status != nil {
					contentMap["status"] = content.Status
				}

				if content.CreatedTime != nil {
					contentMap["created_time"] = content.CreatedTime
				}

				if content.UpdatedTime != nil {
					contentMap["updated_time"] = content.UpdatedTime
				}

				if content.BindedGatewayDeployGroups != nil {
					bindedGatewayDeployGroupsList := []interface{}{}
					for _, bindedGatewayDeployGroups := range content.BindedGatewayDeployGroups {
						bindedGatewayDeployGroupsMap := map[string]interface{}{}

						if bindedGatewayDeployGroups.DeployGroupId != nil {
							bindedGatewayDeployGroupsMap["deploy_group_id"] = bindedGatewayDeployGroups.DeployGroupId
						}

						if bindedGatewayDeployGroups.DeployGroupName != nil {
							bindedGatewayDeployGroupsMap["deploy_group_name"] = bindedGatewayDeployGroups.DeployGroupName
						}

						if bindedGatewayDeployGroups.ApplicationId != nil {
							bindedGatewayDeployGroupsMap["application_id"] = bindedGatewayDeployGroups.ApplicationId
						}

						if bindedGatewayDeployGroups.ApplicationName != nil {
							bindedGatewayDeployGroupsMap["application_name"] = bindedGatewayDeployGroups.ApplicationName
						}

						if bindedGatewayDeployGroups.ApplicationType != nil {
							bindedGatewayDeployGroupsMap["application_type"] = bindedGatewayDeployGroups.ApplicationType
						}

						if bindedGatewayDeployGroups.GroupStatus != nil {
							bindedGatewayDeployGroupsMap["group_status"] = bindedGatewayDeployGroups.GroupStatus
						}

						if bindedGatewayDeployGroups.ClusterType != nil {
							bindedGatewayDeployGroupsMap["cluster_type"] = bindedGatewayDeployGroups.ClusterType
						}

						bindedGatewayDeployGroupsList = append(bindedGatewayDeployGroupsList, bindedGatewayDeployGroupsMap)
					}

					contentMap["binded_gateway_deploy_groups"] = bindedGatewayDeployGroupsList
				}

				if content.ApiCount != nil {
					contentMap["api_count"] = content.ApiCount
				}

				if content.AclMode != nil {
					contentMap["acl_mode"] = content.AclMode
				}

				if content.Description != nil {
					contentMap["description"] = content.Description
				}

				if content.GroupType != nil {
					contentMap["group_type"] = content.GroupType
				}

				if content.GatewayInstanceType != nil {
					contentMap["gateway_instance_type"] = content.GatewayInstanceType
				}

				if content.GatewayInstanceId != nil {
					contentMap["gateway_instance_id"] = content.GatewayInstanceId
				}

				if content.NamespaceNameKey != nil {
					contentMap["namespace_name_key"] = content.NamespaceNameKey
				}

				if content.ServiceNameKey != nil {
					contentMap["service_name_key"] = content.ServiceNameKey
				}

				if content.NamespaceNameKeyPosition != nil {
					contentMap["namespace_name_key_position"] = content.NamespaceNameKeyPosition
				}

				if content.ServiceNameKeyPosition != nil {
					contentMap["service_name_key_position"] = content.ServiceNameKeyPosition
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.GroupId)
			}

			tsfPageApiGroupInfoMap["content"] = contentList
		}

		_ = d.Set("result", []interface{}{tsfPageApiGroupInfoMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tsfPageApiGroupInfoMap); e != nil {
			return e
		}
	}
	return nil
}
