package tsf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTsfBusinessLogConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfBusinessLogConfigsRead,
		Schema: map[string]*schema.Schema{
			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "wild search word。",
			},

			"disable_program_auth_check": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Disable Program auth check or not。",
			},

			"config_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "配置 Id list。",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 business log configurations.注意：此字段可能返回 null，表示无法获取有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total Count.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Log configuration item list. 注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ConfigId。",
									},
									"config_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ConfigName。",
									},
									"config_path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Log 路径 of configuration item.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"config_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "描述 configuration item.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"config_tags": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "configuration 标签注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"config_pipeline": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Pipeline of configuration item.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"config_create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创建时间 of configuration item.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"config_update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "更新时间 of configuration item.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"config_schema": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "ParserSchema of configuration item.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"schema_type": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Schema 类型",
												},
												"schema_content": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "内容 of schema。",
												},
												"schema_date_format": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Schema 格式注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"schema_multiline_pattern": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Schema pattern of configuration item.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"schema_create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "创建时间 of configuration item.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"schema_pattern_layout": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "用户-defined parsing rules.注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"config_associated_groups": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "the associate group of 配置注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"group_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Group Id. 注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"group_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Group 名称 注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"application_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Application Id of Group. 注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"application_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Application 名称 注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"application_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Application 类型 注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"namespace_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Namespace ID to which the deployment group belongs.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"namespace_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Namespace 名称 to which the deployment group belongs.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"cluster_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "集群 ID to which the deployment group belongs.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"cluster_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Cluster 名称 to which the deployment group belongs.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"cluster_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "集群类型 to which the deployment group belongs.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"associated_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Time when the deployment group is associated with the log configuration.注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudTsfBusinessLogConfigsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tsf_business_log_configs.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("disable_program_auth_check"); v != nil {
		paramMap["DisableProgramAuthCheck"] = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("config_id_list"); ok {
		configIdListSet := v.(*schema.Set).List()
		paramMap["ConfigIdList"] = helper.InterfacesStringsPoint(configIdListSet)
	}

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var logConfig *tsf.TsfPageBusinessLogConfig
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfBusinessLogConfigsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		logConfig = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(logConfig.Content))
	tsfPageBusinessLogConfigMap := map[string]interface{}{}
	if logConfig != nil {
		if logConfig.TotalCount != nil {
			tsfPageBusinessLogConfigMap["total_count"] = logConfig.TotalCount
		}

		if logConfig.Content != nil {
			contentList := []interface{}{}
			for _, content := range logConfig.Content {
				contentMap := map[string]interface{}{}

				if content.ConfigId != nil {
					contentMap["config_id"] = content.ConfigId
				}

				if content.ConfigName != nil {
					contentMap["config_name"] = content.ConfigName
				}

				if content.ConfigPath != nil {
					contentMap["config_path"] = content.ConfigPath
				}

				if content.ConfigDesc != nil {
					contentMap["config_desc"] = content.ConfigDesc
				}

				if content.ConfigTags != nil {
					contentMap["config_tags"] = content.ConfigTags
				}

				if content.ConfigPipeline != nil {
					contentMap["config_pipeline"] = content.ConfigPipeline
				}

				if content.ConfigCreateTime != nil {
					contentMap["config_create_time"] = content.ConfigCreateTime
				}

				if content.ConfigUpdateTime != nil {
					contentMap["config_update_time"] = content.ConfigUpdateTime
				}

				if content.ConfigSchema != nil {
					configSchemaMap := map[string]interface{}{}

					if content.ConfigSchema.SchemaType != nil {
						configSchemaMap["schema_type"] = content.ConfigSchema.SchemaType
					}

					if content.ConfigSchema.SchemaContent != nil {
						configSchemaMap["schema_content"] = content.ConfigSchema.SchemaContent
					}

					if content.ConfigSchema.SchemaDateFormat != nil {
						configSchemaMap["schema_date_format"] = content.ConfigSchema.SchemaDateFormat
					}

					if content.ConfigSchema.SchemaMultilinePattern != nil {
						configSchemaMap["schema_multiline_pattern"] = content.ConfigSchema.SchemaMultilinePattern
					}

					if content.ConfigSchema.SchemaCreateTime != nil {
						configSchemaMap["schema_create_time"] = content.ConfigSchema.SchemaCreateTime
					}

					if content.ConfigSchema.SchemaPatternLayout != nil {
						configSchemaMap["schema_pattern_layout"] = content.ConfigSchema.SchemaPatternLayout
					}

					contentMap["config_schema"] = []interface{}{configSchemaMap}
				}

				if content.ConfigAssociatedGroups != nil {
					configAssociatedGroupsList := []interface{}{}
					for _, configAssociatedGroups := range content.ConfigAssociatedGroups {
						configAssociatedGroupsMap := map[string]interface{}{}

						if configAssociatedGroups.GroupId != nil {
							configAssociatedGroupsMap["group_id"] = configAssociatedGroups.GroupId
						}

						if configAssociatedGroups.GroupName != nil {
							configAssociatedGroupsMap["group_name"] = configAssociatedGroups.GroupName
						}

						if configAssociatedGroups.ApplicationId != nil {
							configAssociatedGroupsMap["application_id"] = configAssociatedGroups.ApplicationId
						}

						if configAssociatedGroups.ApplicationName != nil {
							configAssociatedGroupsMap["application_name"] = configAssociatedGroups.ApplicationName
						}

						if configAssociatedGroups.ApplicationType != nil {
							configAssociatedGroupsMap["application_type"] = configAssociatedGroups.ApplicationType
						}

						if configAssociatedGroups.NamespaceId != nil {
							configAssociatedGroupsMap["namespace_id"] = configAssociatedGroups.NamespaceId
						}

						if configAssociatedGroups.NamespaceName != nil {
							configAssociatedGroupsMap["namespace_name"] = configAssociatedGroups.NamespaceName
						}

						if configAssociatedGroups.ClusterId != nil {
							configAssociatedGroupsMap["cluster_id"] = configAssociatedGroups.ClusterId
						}

						if configAssociatedGroups.ClusterName != nil {
							configAssociatedGroupsMap["cluster_name"] = configAssociatedGroups.ClusterName
						}

						if configAssociatedGroups.ClusterType != nil {
							configAssociatedGroupsMap["cluster_type"] = configAssociatedGroups.ClusterType
						}

						if configAssociatedGroups.AssociatedTime != nil {
							configAssociatedGroupsMap["associated_time"] = configAssociatedGroups.AssociatedTime
						}

						configAssociatedGroupsList = append(configAssociatedGroupsList, configAssociatedGroupsMap)
					}

					contentMap["config_associated_groups"] = configAssociatedGroupsList
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.ConfigId)
			}

			tsfPageBusinessLogConfigMap["content"] = contentList
		}

		_ = d.Set("result", []interface{}{tsfPageBusinessLogConfigMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tsfPageBusinessLogConfigMap); e != nil {
			return e
		}
	}
	return nil
}
