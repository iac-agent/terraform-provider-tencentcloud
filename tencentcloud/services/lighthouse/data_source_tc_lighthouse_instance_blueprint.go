package lighthouse

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudLighthouseInstanceBlueprint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseInstanceBlueprintRead,
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "实例 ID list，which currently can contain only one instance。",
			},

			"blueprint_instance_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Blueprint instance list information。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blueprint": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Blueprint instance information。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"blueprint_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Blueprint ID，which is the unique identifier of Blueprint。",
									},
									"display_title": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Blueprint title to be displayed。",
									},
									"display_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Blueprint 版本 to be displayed。",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Image 描述 information。",
									},
									"os_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "OS 名称",
									},
									"platform": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "OS 类型",
									},
									"platform_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "OS 类型，such as LINUX_UNIX and WINDOWS。",
									},
									"blueprint_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Blueprint 类型，such as APP_OS，PURE_OS，and PRIVATE。",
									},
									"image_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Blueprint picture URL",
									},
									"required_system_disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "System disk size 必填 by blueprint （GB）。",
									},
									"blueprint_state": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Blueprint 状态",
									},
									"created_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创建时间 according to ISO 8601 standard. UTC time is used. 格式 is YYYY-MM-DDThh:mm:ssZ。",
									},
									"blueprint_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Blueprint 名称",
									},
									"support_automation_tools": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否blueprint supports automation tools。",
									},
									"required_memory_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Memory size 必填 by blueprint （GB）。",
									},
									"image_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID Lighthouse blueprint shared from a CVM imageNote: this field may return null，indicating that no valid values can be obtained。",
									},
									"community_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "URL of official website of the open-来源 project。",
									},
									"guide_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Guide documentation URL",
									},
									"scene_id_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "数组 IDs of scenes associated with a blueprint注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"docker_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Docker 版本注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"software_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Software list。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Software 名称",
									},
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Software 版本",
									},
									"image_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Software picture URL",
									},
									"install_dir": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Software installation directory。",
									},
									"detail_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "列表 software details。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Unique detail 键",
												},
												"title": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Detail title。",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Detail 值",
												},
											},
										},
									},
								},
							},
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID",
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

func dataSourceTencentCloudLighthouseInstanceBlueprintRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_lighthouse_instance_blueprint.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		paramMap["InstanceIds"] = helper.InterfacesStringsPoint(instanceIdsSet)
	}

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var blueprintInstanceSet []*lighthouse.BlueprintInstance
	instanceIds := make([]string, 0)
	for _, instanceId := range d.Get("instance_ids").(*schema.Set).List() {
		instanceIds = append(instanceIds, instanceId.(string))
	}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseInstanceBlueprintByFilter(ctx, instanceIds)
		if e != nil {
			return tccommon.RetryError(e)
		}
		blueprintInstanceSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(blueprintInstanceSet))
	tmpList := make([]map[string]interface{}, 0, len(blueprintInstanceSet))

	if blueprintInstanceSet != nil {
		for _, blueprintInstance := range blueprintInstanceSet {
			blueprintInstanceMap := map[string]interface{}{}

			if blueprintInstance.Blueprint != nil {
				blueprintMap := map[string]interface{}{}

				if blueprintInstance.Blueprint.BlueprintId != nil {
					blueprintMap["blueprint_id"] = blueprintInstance.Blueprint.BlueprintId
				}

				if blueprintInstance.Blueprint.DisplayTitle != nil {
					blueprintMap["display_title"] = blueprintInstance.Blueprint.DisplayTitle
				}

				if blueprintInstance.Blueprint.DisplayVersion != nil {
					blueprintMap["display_version"] = blueprintInstance.Blueprint.DisplayVersion
				}

				if blueprintInstance.Blueprint.Description != nil {
					blueprintMap["description"] = blueprintInstance.Blueprint.Description
				}

				if blueprintInstance.Blueprint.OsName != nil {
					blueprintMap["os_name"] = blueprintInstance.Blueprint.OsName
				}

				if blueprintInstance.Blueprint.Platform != nil {
					blueprintMap["platform"] = blueprintInstance.Blueprint.Platform
				}

				if blueprintInstance.Blueprint.PlatformType != nil {
					blueprintMap["platform_type"] = blueprintInstance.Blueprint.PlatformType
				}

				if blueprintInstance.Blueprint.BlueprintType != nil {
					blueprintMap["blueprint_type"] = blueprintInstance.Blueprint.BlueprintType
				}

				if blueprintInstance.Blueprint.ImageUrl != nil {
					blueprintMap["image_url"] = blueprintInstance.Blueprint.ImageUrl
				}

				if blueprintInstance.Blueprint.RequiredSystemDiskSize != nil {
					blueprintMap["required_system_disk_size"] = blueprintInstance.Blueprint.RequiredSystemDiskSize
				}

				if blueprintInstance.Blueprint.BlueprintState != nil {
					blueprintMap["blueprint_state"] = blueprintInstance.Blueprint.BlueprintState
				}

				if blueprintInstance.Blueprint.CreatedTime != nil {
					blueprintMap["created_time"] = blueprintInstance.Blueprint.CreatedTime
				}

				if blueprintInstance.Blueprint.BlueprintName != nil {
					blueprintMap["blueprint_name"] = blueprintInstance.Blueprint.BlueprintName
				}

				if blueprintInstance.Blueprint.SupportAutomationTools != nil {
					blueprintMap["support_automation_tools"] = blueprintInstance.Blueprint.SupportAutomationTools
				}

				if blueprintInstance.Blueprint.RequiredMemorySize != nil {
					blueprintMap["required_memory_size"] = blueprintInstance.Blueprint.RequiredMemorySize
				}

				if blueprintInstance.Blueprint.ImageId != nil {
					blueprintMap["image_id"] = blueprintInstance.Blueprint.ImageId
				}

				if blueprintInstance.Blueprint.CommunityUrl != nil {
					blueprintMap["community_url"] = blueprintInstance.Blueprint.CommunityUrl
				}

				if blueprintInstance.Blueprint.GuideUrl != nil {
					blueprintMap["guide_url"] = blueprintInstance.Blueprint.GuideUrl
				}

				if blueprintInstance.Blueprint.SceneIdSet != nil {
					blueprintMap["scene_id_set"] = blueprintInstance.Blueprint.SceneIdSet
				}

				if blueprintInstance.Blueprint.DockerVersion != nil {
					blueprintMap["docker_version"] = blueprintInstance.Blueprint.DockerVersion
				}

				blueprintInstanceMap["blueprint"] = []map[string]interface{}{blueprintMap}
			}

			if blueprintInstance.SoftwareSet != nil {
				softwareSetList := make([]map[string]interface{}, 0)
				for _, softwareSet := range blueprintInstance.SoftwareSet {
					softwareSetMap := map[string]interface{}{}

					if softwareSet.Name != nil {
						softwareSetMap["name"] = softwareSet.Name
					}

					if softwareSet.Version != nil {
						softwareSetMap["version"] = softwareSet.Version
					}

					if softwareSet.ImageUrl != nil {
						softwareSetMap["image_url"] = softwareSet.ImageUrl
					}

					if softwareSet.InstallDir != nil {
						softwareSetMap["install_dir"] = softwareSet.InstallDir
					}

					if softwareSet.DetailSet != nil {
						detailSetList := make([]map[string]interface{}, 0)
						for _, detailSet := range softwareSet.DetailSet {
							detailSetMap := map[string]interface{}{}

							if detailSet.Key != nil {
								detailSetMap["key"] = detailSet.Key
							}

							if detailSet.Title != nil {
								detailSetMap["title"] = detailSet.Title
							}

							if detailSet.Value != nil {
								detailSetMap["value"] = detailSet.Value
							}

							detailSetList = append(detailSetList, detailSetMap)
						}

						softwareSetMap["detail_set"] = detailSetList
					}

					softwareSetList = append(softwareSetList, softwareSetMap)
				}

				blueprintInstanceMap["software_set"] = softwareSetList
			}

			if blueprintInstance.InstanceId != nil {
				blueprintInstanceMap["instance_id"] = blueprintInstance.InstanceId
			}

			ids = append(ids, *blueprintInstance.InstanceId)
			tmpList = append(tmpList, blueprintInstanceMap)
		}

		_ = d.Set("blueprint_instance_set", tmpList)
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
