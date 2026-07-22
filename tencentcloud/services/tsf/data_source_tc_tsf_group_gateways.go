package tsf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTsfGroupGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfGroupGatewaysRead,
		Schema: map[string]*schema.Schema{
			"gateway_deploy_group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "gateway group Id。",
			},

			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "search word。",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "api group information。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "总数",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "api group Info。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "api group id.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "api 组名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"group_context": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "api group context.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"auth_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authentication 类型 secret: 键 authentication; none: no authentication.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Release 状态 drafted: not released. released: released.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"created_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group 创建时间，such as: 2019-06-20 15:51:28.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"updated_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group 更新时间，such as: 2019-06-20 15:51:28.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"binded_gateway_deploy_groups": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Gateway deployment group bound to the API group.注意：此字段可能返回 null，表示未找到有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"deploy_group_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Gateway deployment 组 ID注意：此字段可能返回 null，表示未找到有效值。",
												},
												"deploy_group_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Gateway deployment 组名称注意：此字段可能返回 null，表示未找到有效值。",
												},
												"application_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "application ID.注意：此字段可能返回 null，表示未找到有效值。",
												},
												"application_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "application 名称注意：此字段可能返回 null，表示未找到有效值。",
												},
												"application_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Application category: V: virtual machine application，C: container application.注意：此字段可能返回 null，表示未找到有效值。",
												},
												"group_status": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Application 状态 deployment group，with possible values: Running，Waiting，Paused，Updating，RollingBack，Abnormal，Unknown.注意：此字段可能返回 null，表示未找到有效值。",
												},
												"cluster_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "集群类型，with possible values: C: container，V: virtual machine.注意：此字段可能返回 null，表示未找到有效值。",
												},
											},
										},
									},
									"api_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 APIs.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"acl_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ACL 类型 for accessing the group.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "描述注意：此字段可能返回 null，表示未找到有效值。",
									},
									"group_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group 类型 ms: microservice group; external: external API group.This field may return null，which means no valid 值 was found。",
									},
									"gateway_instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway 实例类型注意：此字段可能返回 null，表示未找到有效值。",
									},
									"gateway_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway instance ID.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"namespace_name_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace parameter 键注意：此字段可能返回 null，表示未找到有效值。",
									},
									"service_name_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Microservice 名称 parameter 键注意：此字段可能返回 null，表示未找到有效值。",
									},
									"namespace_name_key_position": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace parameter location，路径，header，or query. The 默认为 路径注意：此字段可能返回 null，表示未找到有效值。",
									},
									"service_name_key_position": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Microservice 名称 parameter location，路径，header，or query. The 默认为 路径注意：此字段可能返回 null，表示未找到有效值。",
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

func dataSourceTencentCloudTsfGroupGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tsf_group_gateways.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("gateway_deploy_group_id"); ok {
		paramMap["GatewayDeployGroupId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var result *tsf.TsfPageApiGroupInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeTsfGroupGatewaysByFilter(ctx, paramMap)
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
	tsfPageApiGroupInfoMap := map[string]interface{}{}
	if result != nil {

		if result.TotalCount != nil {
			tsfPageApiGroupInfoMap["total_count"] = result.TotalCount
		}

		if result.Content != nil {
			contentList := []interface{}{}
			for _, content := range result.Content {
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
