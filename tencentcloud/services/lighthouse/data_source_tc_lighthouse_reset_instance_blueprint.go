package lighthouse

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudLighthouseResetInstanceBlueprint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseResetInstanceBlueprintRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "偏移量 默认值为 0。",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "数量 返回 results. 默认值为 20. Maximum 值 是 100。",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 listblueprint-idFilter 通过 镜像 ID.类型: StringRequired: noblueprint-typeFilter 通过 镜像 类型有效值：APP_OS: 应用 镜像; PURE_OS: 系统 镜像; PRIVATE: 自定义 imageType: StringRequired: noplatform-typeFilter 通过 镜像 平台 类型有效值：LINUX_UNIX: Linux 或 Unix; WINDOWS: WindowsType: StringRequired: noblueprint-nameFilter 通过 镜像 名称Type: StringRequired: noblueprint-stateFilter 通过 镜像 状态Type: StringRequired: noEach 请求 可以 contain up 到 10 Filters 和 5 过滤器.Values. BlueprintIds 和 Filters 不能 是 指定 在 same 时间。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "待过滤字段",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "过滤值 的 字段。",
						},
					},
				},
			},

			"reset_instance_blueprint_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 scene info。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blueprint_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Mirror details。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"blueprint_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Image ID，其中 是 唯一 identity 的 Blueprint。",
									},
									"display_title": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "mirror 镜像 shows title 到 公有。",
									},
									"display_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "镜像 shows 版本 到 公有。",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Mirror 描述 信息。",
									},
									"os_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Operating 系统 名称",
									},
									"platform": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Operating 系统 平台。",
									},
									"platform_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Operating 系统 平台 类型，such 作为 LINUX_UNIX，WINDOWS。",
									},
									"blueprint_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Image 类型，such 作为 APP_OS，PURE_OS，PRIVATE。",
									},
									"image_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Mirror 镜像 URL",
									},
									"required_system_disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "大小 的 系统 磁盘 必填 对于 镜像 (在 GB)。",
									},
									"blueprint_state": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Mirror 状态",
									},
									"created_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创建时间. Expressed according 到 ISO8601 standard，和 使用 UTC 时间. 格式 是 YYYY-MM-DDThh:mm:ssZ。",
									},
									"blueprint_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Mirror 名称",
									},
									"support_automation_tools": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否image 支持 automation helper。",
									},
									"required_memory_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Memory 必填 对于 mirroring (在 GB)。",
									},
									"image_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CVM 镜像 ID after sharing CVM 镜像 到 lightweight 应用 服务器。",
									},
									"community_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "official website Url。",
									},
									"guide_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Guide article Url。",
									},
									"scene_id_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "mirror association uses scene ID 列表。",
									},
									"docker_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Docker 版本 数量。",
									},
								},
							},
						},
						"is_resettable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否instance 镜像 可以 是 reset 到 目标 镜像。",
						},
						"non_resettable_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "信息 不能 是 reset. 当 mirror 可以 是 reset ''。",
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

func dataSourceTencentCloudLighthouseResetInstanceBlueprintRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_lighthouse_reset_instance_blueprint.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var instanceId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["instance_id"] = instanceId
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["offset"] = v.(int)
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["limit"] = v.(int)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*lighthouse.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := lighthouse.Filter{}
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
		paramMap["filters"] = tmpSet
	}

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var resetInstanceBlueprintSet []*lighthouse.ResetInstanceBlueprint

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseResetInstanceBlueprintByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		resetInstanceBlueprintSet = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(resetInstanceBlueprintSet))
	for _, resetInstanceBlueprint := range resetInstanceBlueprintSet {
		resetInstanceBlueprintMap := make(map[string]interface{})
		if resetInstanceBlueprint.BlueprintInfo != nil {
			blueprintInfo := make(map[string]interface{})

			if resetInstanceBlueprint.BlueprintInfo.BlueprintId != nil {
				blueprintInfo["blueprint_id"] = *resetInstanceBlueprint.BlueprintInfo.BlueprintId
			}
			if resetInstanceBlueprint.BlueprintInfo.DisplayTitle != nil {
				blueprintInfo["display_title"] = *resetInstanceBlueprint.BlueprintInfo.DisplayTitle
			}
			if resetInstanceBlueprint.BlueprintInfo.DisplayVersion != nil {
				blueprintInfo["display_version"] = *resetInstanceBlueprint.BlueprintInfo.DisplayVersion
			}
			if resetInstanceBlueprint.BlueprintInfo.Description != nil {
				blueprintInfo["description"] = *resetInstanceBlueprint.BlueprintInfo.Description
			}
			if resetInstanceBlueprint.BlueprintInfo.OsName != nil {
				blueprintInfo["os_name"] = *resetInstanceBlueprint.BlueprintInfo.OsName
			}
			if resetInstanceBlueprint.BlueprintInfo.Platform != nil {
				blueprintInfo["platform"] = *resetInstanceBlueprint.BlueprintInfo.Platform
			}
			if resetInstanceBlueprint.BlueprintInfo.PlatformType != nil {
				blueprintInfo["platform_type"] = *resetInstanceBlueprint.BlueprintInfo.PlatformType
			}
			if resetInstanceBlueprint.BlueprintInfo.BlueprintType != nil {
				blueprintInfo["blueprint_type"] = *resetInstanceBlueprint.BlueprintInfo.BlueprintType
			}
			if resetInstanceBlueprint.BlueprintInfo.ImageUrl != nil {
				blueprintInfo["image_url"] = *resetInstanceBlueprint.BlueprintInfo.ImageUrl
			}
			if resetInstanceBlueprint.BlueprintInfo.RequiredSystemDiskSize != nil {
				blueprintInfo["required_system_disk_size"] = *resetInstanceBlueprint.BlueprintInfo.RequiredSystemDiskSize
			}
			if resetInstanceBlueprint.BlueprintInfo.BlueprintState != nil {
				blueprintInfo["blueprint_state"] = *resetInstanceBlueprint.BlueprintInfo.BlueprintState
			}
			if resetInstanceBlueprint.BlueprintInfo.CreatedTime != nil {
				blueprintInfo["created_time"] = *resetInstanceBlueprint.BlueprintInfo.CreatedTime
			}
			if resetInstanceBlueprint.BlueprintInfo.BlueprintName != nil {
				blueprintInfo["blueprint_name"] = *resetInstanceBlueprint.BlueprintInfo.BlueprintName
			}
			if resetInstanceBlueprint.BlueprintInfo.SupportAutomationTools != nil {
				blueprintInfo["support_automation_tools"] = *resetInstanceBlueprint.BlueprintInfo.SupportAutomationTools
			}
			if resetInstanceBlueprint.BlueprintInfo.RequiredMemorySize != nil {
				blueprintInfo["required_memory_size"] = *resetInstanceBlueprint.BlueprintInfo.RequiredMemorySize
			}
			if resetInstanceBlueprint.BlueprintInfo.ImageId != nil {
				blueprintInfo["image_id"] = *resetInstanceBlueprint.BlueprintInfo.ImageId
			}
			if resetInstanceBlueprint.BlueprintInfo.CommunityUrl != nil {
				blueprintInfo["community_url"] = *resetInstanceBlueprint.BlueprintInfo.CommunityUrl
			}
			if resetInstanceBlueprint.BlueprintInfo.GuideUrl != nil {
				blueprintInfo["guide_url"] = *resetInstanceBlueprint.BlueprintInfo.GuideUrl
			}
			if resetInstanceBlueprint.BlueprintInfo.SceneIdSet != nil {
				sceneIds := make([]string, 0)
				for _, sceneId := range resetInstanceBlueprint.BlueprintInfo.SceneIdSet {
					sceneIds = append(sceneIds, *sceneId)
				}
				blueprintInfo["scene_id_set"] = sceneIds
			}
			if resetInstanceBlueprint.BlueprintInfo.DockerVersion != nil {
				blueprintInfo["docker_version"] = *resetInstanceBlueprint.BlueprintInfo.DockerVersion
			}
			resetInstanceBlueprintMap["blueprint_info"] = []map[string]interface{}{blueprintInfo}
		}
		if resetInstanceBlueprint.IsResettable != nil {
			resetInstanceBlueprintMap["is_resettable"] = *resetInstanceBlueprint.IsResettable
		}
		if resetInstanceBlueprint.NonResettableMessage != nil {
			resetInstanceBlueprintMap["non_resettable_message"] = *resetInstanceBlueprint.NonResettableMessage
		}
		tmpList = append(tmpList, resetInstanceBlueprintMap)
	}

	d.SetId(instanceId)
	_ = d.Set("reset_instance_blueprint_set", tmpList)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
