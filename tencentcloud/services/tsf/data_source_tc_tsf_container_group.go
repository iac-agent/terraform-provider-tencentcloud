package tsf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTsfContainerGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfContainerGroupRead,
		Schema: map[string]*schema.Schema{
			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "search word，support 组名称",
			},

			"application_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ApplicationId，必填",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "sorting 字段. By 默认值，它 是 createTime 字段. Supports ID，名称，createTime。",
			},

			"order_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "sorting 顺序 By 默认值，它 是 1，indicating 降序 0 表示ascending 顺序，和 1 表示descending 顺序",
			},

			"cluster_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID。",
			},

			"namespace_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Namespace ID。",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "结果 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "列表 部署 groups.注意：此字段可能返回 null，表示未找到有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group ID.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "组名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创建时间.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"server": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Image 服务器.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"repo_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Image 名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"tag_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Image 版本 名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster ID.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "集群名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"namespace_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace ID.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"namespace_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Namespace 名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"cpu_request": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "initial amount 的 CPU，corresponding 到 K8S 请求.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"cpu_limit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "最大 amount 的 CPU，corresponding 到 K8S 限制注意：此字段可能返回 null，表示未找到有效值。",
									},
									"mem_request": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "initial amount 的 内存 allocated 在 MiB，corresponding 到 K8S 请求.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"mem_limit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "最大 amount 的 内存 allocated 在 MiB，corresponding 到 K8S 限制注意：此字段可能返回 null，表示未找到有效值。",
									},
									"alias": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group 描述注意：此字段可能返回 null，表示未找到有效值。",
									},
									"kube_inject_enable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "值 的 KubeInjectEnable.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"updated_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Update 类型注意：此字段可能返回 null，表示未找到有效值。",
									},
								},
							},
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "总数",
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

func dataSourceTencentCloudTsfContainerGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tsf__container_group.read")()
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

	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_id"); ok {
		paramMap["NamespaceId"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var result *tsf.ContainGroupResult
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeTsfDescriptionContainerGroupByFilter(ctx, paramMap)
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
	containGroupResultMap := map[string]interface{}{}
	if result != nil {

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

				if content.CreateTime != nil {
					contentMap["create_time"] = content.CreateTime
				}

				if content.Server != nil {
					contentMap["server"] = content.Server
				}

				if content.RepoName != nil {
					contentMap["repo_name"] = content.RepoName
				}

				if content.TagName != nil {
					contentMap["tag_name"] = content.TagName
				}

				if content.ClusterId != nil {
					contentMap["cluster_id"] = content.ClusterId
				}

				if content.ClusterName != nil {
					contentMap["cluster_name"] = content.ClusterName
				}

				if content.NamespaceId != nil {
					contentMap["namespace_id"] = content.NamespaceId
				}

				if content.NamespaceName != nil {
					contentMap["namespace_name"] = content.NamespaceName
				}

				if content.CpuRequest != nil {
					contentMap["cpu_request"] = content.CpuRequest
				}

				if content.CpuLimit != nil {
					contentMap["cpu_limit"] = content.CpuLimit
				}

				if content.MemRequest != nil {
					contentMap["mem_request"] = content.MemRequest
				}

				if content.MemLimit != nil {
					contentMap["mem_limit"] = content.MemLimit
				}

				if content.Alias != nil {
					contentMap["alias"] = content.Alias
				}

				if content.KubeInjectEnable != nil {
					contentMap["kube_inject_enable"] = content.KubeInjectEnable
				}

				if content.UpdatedTime != nil {
					contentMap["updated_time"] = content.UpdatedTime
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.GroupId)
			}

			containGroupResultMap["content"] = contentList
		}

		if result.TotalCount != nil {
			containGroupResultMap["total_count"] = result.TotalCount
		}

		_ = d.Set("result", []interface{}{containGroupResultMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), containGroupResultMap); e != nil {
			return e
		}
	}
	return nil
}
