package tsf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTsfApplicationConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfApplicationConfigRead,
		Schema: map[string]*schema.Schema{
			"application_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Application ID，查询 all 当 不 提供。",
			},

			"config_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Configuration ID，查询 all 使用 higher 优先级 当 不 提供。",
			},

			"config_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Configuration ID 列表，查询 all 使用 lower 优先级 当 不 提供。",
			},

			"config_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Configuration 名称，precise 查询，查询 all 当 不 提供。",
			},

			"config_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Configuration 版本，precise 查询，查询 all 当 不 提供。",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Paginated 配置 item 列表. 注意：此字段可能返回 null，表示无法获取有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "TsfPageConfig。",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Configuration item 列表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration ID. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"config_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration 名称 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"config_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration 版本 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"config_version_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration 版本 描述 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"config_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration 值 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"config_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration 类型 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"creation_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CreationTime. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"application_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "应用 ID. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"application_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "应用 ID. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"delete_flag": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "delete flag，true: allow delete; false: delete prohibit。",
									},
									"last_update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "last 更新时间. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"config_version_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "配置 版本 count. 注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudTsfApplicationConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tsf_application_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("application_id"); ok {
		paramMap["ApplicationId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_id"); ok {
		paramMap["ConfigId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_id_list"); ok {
		configIdListSet := v.(*schema.Set).List()
		paramMap["ConfigIdList"] = helper.InterfacesStringsPoint(configIdListSet)
	}

	if v, ok := d.GetOk("config_name"); ok {
		paramMap["ConfigName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("config_version"); ok {
		paramMap["ConfigVersion"] = helper.String(v.(string))
	}

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var config *tsf.TsfPageConfig
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfApplicationConfigByFilter(ctx, paramMap)
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

		err = d.Set("result", []interface{}{tsfPageConfigMap})
		if err != nil {
			return err
		}
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
