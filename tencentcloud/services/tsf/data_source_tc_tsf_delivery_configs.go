package tsf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTsfDeliveryConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfDeliveryConfigsRead,
		Schema: map[string]*schema.Schema{
			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "search word。",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "deploy group information about the deployment group associated with a delivery item.注意：此字段可能返回 null，表示未获取到有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "总数 注意：此字段可能返回 null，表示未获取到有效值。",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "内容 注意：此字段可能返回 null，表示未获取到有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "配置 id。",
									},
									"config_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "配置 名称",
									},
									"collect_path": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "harvest log 路径 注意：此字段可能返回 null，表示未获取到有效值。",
									},
									"groups": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Associated deployment group information.注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"group_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Group Id。",
												},
												"group_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Group 名称",
												},
												"cluster_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "集群类型",
												},
												"cluster_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "集群 ID 注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"cluster_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Cluster 名称 注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"namespace_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Namespace 名称 注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"associate_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Associate Time. 注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创建时间.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"kafka_v_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Kafka VIP 注意：此字段可能返回 null，表示未获取到有效值。",
									},
									"kafka_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "KafkaAddress refers to the 地址 of a Kafka server.注意：此字段可能返回 null，表示未获取到有效值。",
									},
									"kafka_v_port": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Kafka VPort. 注意：此字段可能返回 null，表示未获取到有效值。",
									},
									"topic": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Topic. 注意：此字段可能返回 null，表示未获取到有效值。",
									},
									"line_rule": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Line Rule for log. 注意：此字段可能返回 null，表示未获取到有效值。",
									},
									"custom_rule": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CustomRule 指定a custom line separator rule.注意：此字段可能返回 null，表示未获取到有效值。",
									},
									"enable_global_line_rule": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "表示是否a single row rule should be applied.注意：此字段可能返回 null，表示未获取到有效值。",
									},
									"enable_auth": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "whether use auth for kafka. 注意：此字段可能返回 null，表示未获取到有效值。",
									},
									"username": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "用户 名称 注意：此字段可能返回 null，表示未获取到有效值。",
									},
									"password": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "密码 注意：此字段可能返回 null，表示未获取到有效值。",
									},
									"kafka_infos": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Kafka Infos. 注意：此字段可能返回 null，表示未获取到有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"topic": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Kafka topic. 注意：此字段可能返回 null，表示未获取到有效值。",
												},
												"path": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "harvest log 路径 注意：此字段可能返回 null，表示未获取到有效值。",
												},
												"line_rule": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Line rule 指定type of line separator used in a file. It can have one of the following values: 默认值：The default line separator is 用于separate lines in the file. time: The lines in the file are separated based on time. custom: A custom line separator is used. In this case，the CustomRule field should be filled with the specific custom 值 注意：此字段可能返回 null，表示未获取到有效值。",
												},
												"custom_rule": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Custom Line Rule。",
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

func dataSourceTencentCloudTsfDeliveryConfigsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tsf_delivery_configs.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var bindGroups *tsf.DeliveryConfigBindGroups

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfDeliveryConfigsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		bindGroups = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(bindGroups.Content))
	deliveryConfigBindGroupsMap := map[string]interface{}{}
	if bindGroups != nil {
		if bindGroups.TotalCount != nil {
			deliveryConfigBindGroupsMap["total_count"] = bindGroups.TotalCount
		}

		if bindGroups.Content != nil {
			contentList := []interface{}{}
			for _, content := range bindGroups.Content {
				contentMap := map[string]interface{}{}

				if content.ConfigId != nil {
					contentMap["config_id"] = *content.ConfigId
				}

				if content.ConfigName != nil {
					contentMap["config_name"] = content.ConfigName
				}

				if content.CollectPath != nil {
					contentMap["collect_path"] = content.CollectPath
				}

				if content.Groups != nil {
					groupsList := []interface{}{}
					for _, groups := range content.Groups {
						groupsMap := map[string]interface{}{}

						if groups.GroupId != nil {
							groupsMap["group_id"] = groups.GroupId
						}

						if groups.GroupName != nil {
							groupsMap["group_name"] = groups.GroupName
						}

						if groups.ClusterType != nil {
							groupsMap["cluster_type"] = groups.ClusterType
						}

						if groups.ClusterId != nil {
							groupsMap["cluster_id"] = groups.ClusterId
						}

						if groups.ClusterName != nil {
							groupsMap["cluster_name"] = groups.ClusterName
						}

						if groups.NamespaceName != nil {
							groupsMap["namespace_name"] = groups.NamespaceName
						}

						if groups.AssociateTime != nil {
							groupsMap["associate_time"] = groups.AssociateTime
						}

						groupsList = append(groupsList, groupsMap)
					}

					contentMap["groups"] = groupsList
				}

				if content.CreateTime != nil {
					contentMap["create_time"] = content.CreateTime
				}

				if content.KafkaVIp != nil {
					contentMap["kafka_v_ip"] = content.KafkaVIp
				}

				if content.KafkaAddress != nil {
					contentMap["kafka_address"] = content.KafkaAddress
				}

				if content.KafkaVPort != nil {
					contentMap["kafka_v_port"] = content.KafkaVPort
				}

				if content.Topic != nil {
					contentMap["topic"] = content.Topic
				}

				if content.LineRule != nil {
					contentMap["line_rule"] = content.LineRule
				}

				if content.CustomRule != nil {
					contentMap["custom_rule"] = content.CustomRule
				}

				if content.EnableGlobalLineRule != nil {
					contentMap["enable_global_line_rule"] = content.EnableGlobalLineRule
				}

				if content.EnableAuth != nil {
					contentMap["enable_auth"] = content.EnableAuth
				}

				if content.Username != nil {
					contentMap["username"] = content.Username
				}

				if content.Password != nil {
					contentMap["password"] = content.Password
				}

				if content.KafkaInfos != nil {
					kafkaInfosList := []interface{}{}
					for _, kafkaInfos := range content.KafkaInfos {
						kafkaInfosMap := map[string]interface{}{}

						if kafkaInfos.Topic != nil {
							kafkaInfosMap["topic"] = kafkaInfos.Topic
						}

						if kafkaInfos.Path != nil {
							kafkaInfosMap["path"] = kafkaInfos.Path
						}

						if kafkaInfos.LineRule != nil {
							kafkaInfosMap["line_rule"] = kafkaInfos.LineRule
						}

						if kafkaInfos.CustomRule != nil {
							kafkaInfosMap["custom_rule"] = kafkaInfos.CustomRule
						}

						kafkaInfosList = append(kafkaInfosList, kafkaInfosMap)
					}

					contentMap["kafka_infos"] = kafkaInfosList
					ids = append(ids, *content.ConfigId)
				}

				contentList = append(contentList, contentMap)
			}

			deliveryConfigBindGroupsMap["content"] = contentList
		}

		_ = d.Set("result", []interface{}{deliveryConfigBindGroupsMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), deliveryConfigBindGroupsMap); e != nil {
			return e
		}
	}
	return nil
}
