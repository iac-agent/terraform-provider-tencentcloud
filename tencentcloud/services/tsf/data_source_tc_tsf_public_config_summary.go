package tsf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTsfPublicConfigSummary() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfPublicConfigSummaryRead,
		Schema: map[string]*schema.Schema{
			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Query keyword for fuzzy search: configuration item 名称 如果未传入 in，the full set will be queried。",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "排序方式 time: creation_time; 排序方式 名称: config_name。",
			},

			"order_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Pass 0 for 升序 and 1 for 降序",
			},

			"config_tag_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "配置 标签列表",
			},

			"disable_program_auth_check": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "是否disable dataset authentication。",
			},

			"config_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "配置 Id List。",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Public 配置 Page Item。",
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
							Description: "配置 list。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration item ID.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"config_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration 名称注意：此字段可能返回 null，表示未找到有效值。",
									},
									"config_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration 版本 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"config_version_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration 版本 描述注意：此字段可能返回 null，表示未找到有效值。",
									},
									"config_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration 值注意：此字段可能返回 null，表示未找到有效值。",
									},
									"config_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "配置 类型 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"creation_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创建时间.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"application_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application ID.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"application_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application 名称 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"delete_flag": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Deletion flag，true: deletable; false: not deletable.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"last_update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Last 更新时间.注意：此字段可能返回 null，表示未找到有效值。",
									},
									"config_version_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Configure 版本 count.注意：此字段可能返回 null，表示未找到有效值。",
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

func dataSourceTencentCloudTsfPublicConfigSummaryRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tsf_public_config_summary.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("order_type"); v != nil {
		paramMap["OrderType"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("config_tag_list"); ok {
		configTagListSet := v.(*schema.Set).List()
		paramMap["ConfigTagList"] = helper.InterfacesStringsPoint(configTagListSet)
	}

	if v, _ := d.GetOk("disable_program_auth_check"); v != nil {
		paramMap["DisableProgramAuthCheck"] = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("config_id_list"); ok {
		configIdListSet := v.(*schema.Set).List()
		paramMap["ConfigIdList"] = helper.InterfacesStringsPoint(configIdListSet)
	}

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var config *tsf.TsfPageConfig
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfPublicConfigSummaryByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		config = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(config.Content))
	tsfPageConfigMap := map[string]interface{}{}
	if config != nil {
		if config.TotalCount != nil {
			tsfPageConfigMap["total_count"] = config.TotalCount
		}

		if config.Content != nil {
			contentList := []interface{}{}
			for _, content := range config.Content {
				contentMap := map[string]interface{}{}

				if content.ConfigId != nil {
					contentMap["config_id"] = content.ConfigId
				}

				if content.ConfigName != nil {
					contentMap["config_name"] = content.ConfigName
				}

				if content.ConfigVersion != nil {
					contentMap["config_version"] = content.ConfigVersion
				}

				if content.ConfigVersionDesc != nil {
					contentMap["config_version_desc"] = content.ConfigVersionDesc
				}

				if content.ConfigValue != nil {
					contentMap["config_value"] = content.ConfigValue
				}

				if content.ConfigType != nil {
					contentMap["config_type"] = content.ConfigType
				}

				if content.CreationTime != nil {
					contentMap["creation_time"] = content.CreationTime
				}

				if content.ApplicationId != nil {
					contentMap["application_id"] = content.ApplicationId
				}

				if content.ApplicationName != nil {
					contentMap["application_name"] = content.ApplicationName
				}

				if content.DeleteFlag != nil {
					contentMap["delete_flag"] = content.DeleteFlag
				}

				if content.LastUpdateTime != nil {
					contentMap["last_update_time"] = content.LastUpdateTime
				}

				if content.ConfigVersionCount != nil {
					contentMap["config_version_count"] = content.ConfigVersionCount
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.ConfigId)
			}

			tsfPageConfigMap["content"] = contentList
		}

		_ = d.Set("result", []interface{}{tsfPageConfigMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tsfPageConfigMap); e != nil {
			return e
		}
	}
	return nil
}
