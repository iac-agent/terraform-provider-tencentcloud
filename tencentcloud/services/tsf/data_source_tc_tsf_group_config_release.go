package tsf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTsfGroupConfigRelease() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfGroupConfigReleaseRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "groupId。",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Information related to the deployment group release.注意：此字段可能返回 null，表示未找到有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"package_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Package Id.注意：此字段可能返回 null，表示未找到有效值。",
						},
						"package_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Package 名称注意：此字段可能返回 null，表示未找到有效值。",
						},
						"package_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Package 版本注意：此字段可能返回 null，表示未找到有效值。",
						},
						"repo_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "image 名称注意：此字段可能返回 null，表示未找到有效值。",
						},
						"tag_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "image 标签 名称注意：此字段可能返回 null，表示未找到有效值。",
						},
						"public_config_release_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Release public 配置 list。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_release_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release ID.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"config_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item  ID.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"config_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item 名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"config_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration 版本注意：此字段可能返回 null，表示未找到有效值。",
									},
									"release_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release time.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release 组 ID注意：此字段可能返回 null，表示未找到有效值。",
									},
									"group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release 组名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"namespace_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release namespace ID.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"namespace_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release namespace 名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release cluster ID.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release 集群名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"release_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release 描述注意：此字段可能返回 null，表示未找到有效值。",
									},
									"application_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release application ID.注意：此字段可能返回 null，表示未找到有效值。",
									},
								},
							},
						},
						"config_release_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Configuration item release list.注意：此字段可能返回 null，表示未找到有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_release_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release ID.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"config_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release 配置 ID.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"config_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release 配置 名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"config_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release 配置 版本注意：此字段可能返回 null，表示未找到有效值。",
									},
									"release_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release time.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release 配置 组 ID注意：此字段可能返回 null，表示未找到有效值。",
									},
									"group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release 配置 组名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"namespace_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release namespace ID.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"namespace_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release namespace 名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release cluster ID.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release 集群名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"release_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release 描述注意：此字段可能返回 null，表示未找到有效值。",
									},
									"application_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release 配置 ID.注意：此字段可能返回 null，表示未找到有效值。",
									},
								},
							},
						},
						"file_config_release_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "File configuration item release list.注意：此字段可能返回 null，表示未找到有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_release_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release ID.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"config_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release 配置 ID.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"config_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release 配置 名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"config_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release 配置 版本注意：此字段可能返回 null，表示未找到有效值。",
									},
									"release_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release time.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release 配置 组 ID注意：此字段可能返回 null，表示未找到有效值。",
									},
									"group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release 配置 组名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"namespace_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release namespace ID.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"namespace_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release namespace 名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release cluster ID.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release 集群名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"release_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item release 描述注意：此字段可能返回 null，表示未找到有效值。",
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

func dataSourceTencentCloudTsfGroupConfigReleaseRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tsf_group_config_release.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var groupId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
		paramMap["GroupId"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var result *tsf.GroupRelease
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeTsfGroupConfigReleaseByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	groupReleaseMap := map[string]interface{}{}
	if result != nil {
		if result.PackageId != nil {
			groupReleaseMap["package_id"] = result.PackageId
		}

		if result.PackageName != nil {
			groupReleaseMap["package_name"] = result.PackageName
		}

		if result.PackageVersion != nil {
			groupReleaseMap["package_version"] = result.PackageVersion
		}

		if result.RepoName != nil {
			groupReleaseMap["repo_name"] = result.RepoName
		}

		if result.TagName != nil {
			groupReleaseMap["tag_name"] = result.TagName
		}

		if result.PublicConfigReleaseList != nil {
			publicConfigReleaseListList := []interface{}{}
			for _, publicConfigReleaseList := range result.PublicConfigReleaseList {
				publicConfigReleaseListMap := map[string]interface{}{}

				if publicConfigReleaseList.ConfigReleaseId != nil {
					publicConfigReleaseListMap["config_release_id"] = publicConfigReleaseList.ConfigReleaseId
				}

				if publicConfigReleaseList.ConfigId != nil {
					publicConfigReleaseListMap["config_id"] = publicConfigReleaseList.ConfigId
				}

				if publicConfigReleaseList.ConfigName != nil {
					publicConfigReleaseListMap["config_name"] = publicConfigReleaseList.ConfigName
				}

				if publicConfigReleaseList.ConfigVersion != nil {
					publicConfigReleaseListMap["config_version"] = publicConfigReleaseList.ConfigVersion
				}

				if publicConfigReleaseList.ReleaseTime != nil {
					publicConfigReleaseListMap["release_time"] = publicConfigReleaseList.ReleaseTime
				}

				if publicConfigReleaseList.GroupId != nil {
					publicConfigReleaseListMap["group_id"] = publicConfigReleaseList.GroupId
				}

				if publicConfigReleaseList.GroupName != nil {
					publicConfigReleaseListMap["group_name"] = publicConfigReleaseList.GroupName
				}

				if publicConfigReleaseList.NamespaceId != nil {
					publicConfigReleaseListMap["namespace_id"] = publicConfigReleaseList.NamespaceId
				}

				if publicConfigReleaseList.NamespaceName != nil {
					publicConfigReleaseListMap["namespace_name"] = publicConfigReleaseList.NamespaceName
				}

				if publicConfigReleaseList.ClusterId != nil {
					publicConfigReleaseListMap["cluster_id"] = publicConfigReleaseList.ClusterId
				}

				if publicConfigReleaseList.ClusterName != nil {
					publicConfigReleaseListMap["cluster_name"] = publicConfigReleaseList.ClusterName
				}

				if publicConfigReleaseList.ReleaseDesc != nil {
					publicConfigReleaseListMap["release_desc"] = publicConfigReleaseList.ReleaseDesc
				}

				if publicConfigReleaseList.ApplicationId != nil {
					publicConfigReleaseListMap["application_id"] = publicConfigReleaseList.ApplicationId
				}

				publicConfigReleaseListList = append(publicConfigReleaseListList, publicConfigReleaseListMap)
			}

			groupReleaseMap["public_config_release_list"] = publicConfigReleaseListList
		}

		if result.ConfigReleaseList != nil {
			configReleaseListList := []interface{}{}
			for _, configReleaseList := range result.ConfigReleaseList {
				configReleaseListMap := map[string]interface{}{}

				if configReleaseList.ConfigReleaseId != nil {
					configReleaseListMap["config_release_id"] = configReleaseList.ConfigReleaseId
				}

				if configReleaseList.ConfigId != nil {
					configReleaseListMap["config_id"] = configReleaseList.ConfigId
				}

				if configReleaseList.ConfigName != nil {
					configReleaseListMap["config_name"] = configReleaseList.ConfigName
				}

				if configReleaseList.ConfigVersion != nil {
					configReleaseListMap["config_version"] = configReleaseList.ConfigVersion
				}

				if configReleaseList.ReleaseTime != nil {
					configReleaseListMap["release_time"] = configReleaseList.ReleaseTime
				}

				if configReleaseList.GroupId != nil {
					configReleaseListMap["group_id"] = configReleaseList.GroupId
				}

				if configReleaseList.GroupName != nil {
					configReleaseListMap["group_name"] = configReleaseList.GroupName
				}

				if configReleaseList.NamespaceId != nil {
					configReleaseListMap["namespace_id"] = configReleaseList.NamespaceId
				}

				if configReleaseList.NamespaceName != nil {
					configReleaseListMap["namespace_name"] = configReleaseList.NamespaceName
				}

				if configReleaseList.ClusterId != nil {
					configReleaseListMap["cluster_id"] = configReleaseList.ClusterId
				}

				if configReleaseList.ClusterName != nil {
					configReleaseListMap["cluster_name"] = configReleaseList.ClusterName
				}

				if configReleaseList.ReleaseDesc != nil {
					configReleaseListMap["release_desc"] = configReleaseList.ReleaseDesc
				}

				if configReleaseList.ApplicationId != nil {
					configReleaseListMap["application_id"] = configReleaseList.ApplicationId
				}

				configReleaseListList = append(configReleaseListList, configReleaseListMap)
			}

			groupReleaseMap["config_release_list"] = configReleaseListList
		}

		if result.FileConfigReleaseList != nil {
			fileConfigReleaseListList := []interface{}{}
			for _, fileConfigReleaseList := range result.FileConfigReleaseList {
				fileConfigReleaseListMap := map[string]interface{}{}

				if fileConfigReleaseList.ConfigReleaseId != nil {
					fileConfigReleaseListMap["config_release_id"] = fileConfigReleaseList.ConfigReleaseId
				}

				if fileConfigReleaseList.ConfigId != nil {
					fileConfigReleaseListMap["config_id"] = fileConfigReleaseList.ConfigId
				}

				if fileConfigReleaseList.ConfigName != nil {
					fileConfigReleaseListMap["config_name"] = fileConfigReleaseList.ConfigName
				}

				if fileConfigReleaseList.ConfigVersion != nil {
					fileConfigReleaseListMap["config_version"] = fileConfigReleaseList.ConfigVersion
				}

				if fileConfigReleaseList.ReleaseTime != nil {
					fileConfigReleaseListMap["release_time"] = fileConfigReleaseList.ReleaseTime
				}

				if fileConfigReleaseList.GroupId != nil {
					fileConfigReleaseListMap["group_id"] = fileConfigReleaseList.GroupId
				}

				if fileConfigReleaseList.GroupName != nil {
					fileConfigReleaseListMap["group_name"] = fileConfigReleaseList.GroupName
				}

				if fileConfigReleaseList.NamespaceId != nil {
					fileConfigReleaseListMap["namespace_id"] = fileConfigReleaseList.NamespaceId
				}

				if fileConfigReleaseList.NamespaceName != nil {
					fileConfigReleaseListMap["namespace_name"] = fileConfigReleaseList.NamespaceName
				}

				if fileConfigReleaseList.ClusterId != nil {
					fileConfigReleaseListMap["cluster_id"] = fileConfigReleaseList.ClusterId
				}

				if fileConfigReleaseList.ClusterName != nil {
					fileConfigReleaseListMap["cluster_name"] = fileConfigReleaseList.ClusterName
				}

				if fileConfigReleaseList.ReleaseDesc != nil {
					fileConfigReleaseListMap["release_desc"] = fileConfigReleaseList.ReleaseDesc
				}

				fileConfigReleaseListList = append(fileConfigReleaseListList, fileConfigReleaseListMap)
			}

			groupReleaseMap["file_config_release_list"] = fileConfigReleaseListList
		}

		_ = d.Set("result", []interface{}{groupReleaseMap})
	}

	d.SetId(helper.DataResourceIdsHash([]string{groupId}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), groupReleaseMap); e != nil {
			return e
		}
	}
	return nil
}
