package tsf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTsfGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfGroupsRead,
		Schema: map[string]*schema.Schema{
			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "searchWord，support groupName。",
			},

			"application_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "applicationId。",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "sort term。",
			},

			"order_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "顺序 类型，0 desc，1 asc。",
			},

			"namespace_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "命名空间 ID。",
			},

			"cluster_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "clusterId。",
			},

			"group_resource_type_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Group resourceType 列表。",
			},

			"status": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "组 状态 过滤器，`Running`: running，`Unknown`: unknown。",
			},

			"group_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "组 ID 列表。",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Pagination 信息 的 virtual machine 部署 组.注意：此字段可能返回 null，表示未找到有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "总数 virtual machine 部署 组. 注意：此字段可能返回 null，表示未找到有效值。",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Virtual machine 部署 组 列表. 注意：此字段可能返回 null，表示未找到有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "组 ID 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "组 ID 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"application_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application 类型 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"group_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group 描述 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group 更新时间. 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "集群 ID 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"startup_parameters": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group start up Parameters. 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"namespace_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace ID. 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Create Time. 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "集群名称 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"application_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application ID. 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"application_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application 名称 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"namespace_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace 名称 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"microservice_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Microservice 类型 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"group_resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group 资源类型 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"updated_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "更新时间. 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"deploy_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group 描述 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"alias": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group alias. 注意：此字段可能返回 null，表示未找到有效值。",
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

func dataSourceTencentCloudTsfGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tsf_groups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("application_id"); ok {
		paramMap["ApplicationId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("order_type"); v != nil {
		paramMap["OrderType"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("namespace_id"); ok {
		paramMap["NamespaceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_resource_type_list"); ok {
		groupResourceTypeListSet := v.(*schema.Set).List()
		paramMap["GroupResourceTypeList"] = helper.InterfacesStringsPoint(groupResourceTypeListSet)
	}

	if v, ok := d.GetOk("status"); ok {
		paramMap["Status"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_id_list"); ok {
		groupIdListSet := v.(*schema.Set).List()
		paramMap["GroupIdList"] = helper.InterfacesStringsPoint(groupIdListSet)
	}

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var result *tsf.TsfPageVmGroup
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeTsfGroupsByFilter(ctx, paramMap)
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
	tsfPageVmGroupMap := map[string]interface{}{}
	if result != nil {
		if result.TotalCount != nil {
			tsfPageVmGroupMap["total_count"] = result.TotalCount
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

				if content.ApplicationType != nil {
					contentMap["application_type"] = content.ApplicationType
				}

				if content.GroupDesc != nil {
					contentMap["group_desc"] = content.GroupDesc
				}

				if content.UpdateTime != nil {
					contentMap["update_time"] = content.UpdateTime
				}

				if content.ClusterId != nil {
					contentMap["cluster_id"] = content.ClusterId
				}

				if content.StartupParameters != nil {
					contentMap["startup_parameters"] = content.StartupParameters
				}

				if content.NamespaceId != nil {
					contentMap["namespace_id"] = content.NamespaceId
				}

				if content.CreateTime != nil {
					contentMap["create_time"] = content.CreateTime
				}

				if content.ClusterName != nil {
					contentMap["cluster_name"] = content.ClusterName
				}

				if content.ApplicationId != nil {
					contentMap["application_id"] = content.ApplicationId
				}

				if content.ApplicationName != nil {
					contentMap["application_name"] = content.ApplicationName
				}

				if content.NamespaceName != nil {
					contentMap["namespace_name"] = content.NamespaceName
				}

				if content.MicroserviceType != nil {
					contentMap["microservice_type"] = content.MicroserviceType
				}

				if content.GroupResourceType != nil {
					contentMap["group_resource_type"] = content.GroupResourceType
				}

				if content.UpdatedTime != nil {
					contentMap["updated_time"] = content.UpdatedTime
				}

				if content.DeployDesc != nil {
					contentMap["deploy_desc"] = content.DeployDesc
				}

				if content.Alias != nil {
					contentMap["alias"] = content.Alias
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.GroupId)
			}

			tsfPageVmGroupMap["content"] = contentList
		}

		_ = d.Set("result", []interface{}{tsfPageVmGroupMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tsfPageVmGroupMap); e != nil {
			return e
		}
	}
	return nil
}
