package postgresql

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPostgresqlDefaultParameters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresqlDefaultParametersRead,
		Schema: map[string]*schema.Schema{
			"db_major_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "major 数据库 版本 数量，such 作为 11，12，13。",
			},

			"db_engine": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Database 引擎，such 作为 postgresql，mssql_compatible。",
			},

			"param_info_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Parameter information注意：此字段可能返回 null，表示无法获取有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Parameter IDNote: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter nameNote: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
						},
						"param_value_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "值 类型 参数. 有效值：`整数`，`real` (floating-point)，`bool`，`enum`，`mutil_enum` (此 类型 参数 可以 是 集合 到 多个 enumerated 值).For `整数` 或 `real` 参数， `Min` 字段 表示 最小 值 和 `Max` 字段 最大 值 For `bool` 参数， 有效 值 include `true` 和 `false`; For `enum` 或 `mutil_enum` 参数， `EnumValue` 字段 表示 有效 值.注意: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
						},
						"unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unit 的 参数 值 如果 参数 has 无 单位，此 字段 将 返回 null.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "默认值 的 参数，其中 是 返回 作为 stringNote: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
						},
						"current_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "当前值 的 参数，其中 是 返回 作为 stringNote: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
						},
						"max": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "最大 值 的 `整数` 或 `real` parameterNote: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
						},
						"enum_value": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "值 范围 的 enum parameterNote: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
						},
						"min": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "最小 值 的 `整数` 或 `real` parameterNote: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
						},
						"param_description_ch": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter 描述 在 ChineseNote: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
						},
						"param_description_en": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter 描述 在 EnglishNote: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
						},
						"need_reboot": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否restart 实例 对于 modified 参数 到 take effect. 有效值：`true` (yes)，`false` (无)注意: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
						},
						"classification_cn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter category 在 ChineseNote: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
						},
						"classification_en": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter category 在 EnglishNote: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
						},
						"spec_related": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否parameter 是 related 到 specifications. 有效值：`true` (yes)，`false` (无)注意: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
						},
						"advanced": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否为a 键 参数. 有效值：`true` (yes，和 modifying 它 可能 affect 实例 performance)，`false` (无)注意: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
						},
						"last_modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "最后修改时间 的 parameterNote: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
						},
						"standby_related": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Primary-standby constraint. 有效值：`0` (无 constraint)，`1` ( 参数 值 的 standby 服务器 必须 是 greater 比 该 的 primary 服务器)，`2` ( 参数 值 的 primary 服务器 必须 是 greater 比 该 的 standby 服务器.)注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"version_relation_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Associated 参数 版本 信息，其中 refers 到 detailed 参数 信息 的 kernel 版本注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Parameter name注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"db_kernel_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "kernel 版本 该 corresponds 到 参数 information注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Default 参数 值 under kernel 版本 和 规格 的 instance注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Unit 的 参数 值 如果 参数 has 无 单位，此 字段 将 返回 null.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"max": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "最大 值 的 `整数` 或 `real` parameter注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"min": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "最小 值 的 `整数` 或 `real` parameter注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"enum_value": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "值 范围 的 enum parameter注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"spec_relation_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Associated 参数 规格 信息，其中 refers 到 detailed 参数 信息 的 specifications.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Parameter name注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"memory": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "规格 该 corresponds 到 参数 information注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "默认值 参数 值 under 此 specification注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Unit 的 参数 值 如果 参数 has 无 单位，此 字段 将 返回 null.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"max": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "最大 值 的 `整数` 或 `real` parameter注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"min": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "最小 值 的 `整数` 或 `real` parameter注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"enum_value": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "值 范围 的 enum parameter注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudPostgresqlDefaultParametersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_postgresql_default_parameters.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("db_major_version"); ok {
		paramMap["DBMajorVersion"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_engine"); ok {
		paramMap["DBEngine"] = helper.String(v.(string))
	}

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var paramInfoSet []*postgresql.ParamInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePostgresqlDefaultParametersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		paramInfoSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(paramInfoSet))
	tmpList := make([]map[string]interface{}, 0, len(paramInfoSet))

	if paramInfoSet != nil {
		for _, paramInfo := range paramInfoSet {
			paramInfoMap := map[string]interface{}{}

			if paramInfo.ID != nil {
				paramInfoMap["id"] = paramInfo.ID
			}

			if paramInfo.Name != nil {
				paramInfoMap["name"] = paramInfo.Name
			}

			if paramInfo.ParamValueType != nil {
				paramInfoMap["param_value_type"] = paramInfo.ParamValueType
			}

			if paramInfo.Unit != nil {
				paramInfoMap["unit"] = paramInfo.Unit
			}

			if paramInfo.DefaultValue != nil {
				paramInfoMap["default_value"] = paramInfo.DefaultValue
			}

			if paramInfo.CurrentValue != nil {
				paramInfoMap["current_value"] = paramInfo.CurrentValue
			}

			if paramInfo.Max != nil {
				paramInfoMap["max"] = paramInfo.Max
			}

			if paramInfo.EnumValue != nil {
				paramInfoMap["enum_value"] = paramInfo.EnumValue
			}

			if paramInfo.Min != nil {
				paramInfoMap["min"] = paramInfo.Min
			}

			if paramInfo.ParamDescriptionCH != nil {
				paramInfoMap["param_description_ch"] = paramInfo.ParamDescriptionCH
			}

			if paramInfo.ParamDescriptionEN != nil {
				paramInfoMap["param_description_en"] = paramInfo.ParamDescriptionEN
			}

			if paramInfo.NeedReboot != nil {
				paramInfoMap["need_reboot"] = paramInfo.NeedReboot
			}

			if paramInfo.ClassificationCN != nil {
				paramInfoMap["classification_cn"] = paramInfo.ClassificationCN
			}

			if paramInfo.ClassificationEN != nil {
				paramInfoMap["classification_en"] = paramInfo.ClassificationEN
			}

			if paramInfo.SpecRelated != nil {
				paramInfoMap["spec_related"] = paramInfo.SpecRelated
			}

			if paramInfo.Advanced != nil {
				paramInfoMap["advanced"] = paramInfo.Advanced
			}

			if paramInfo.LastModifyTime != nil {
				paramInfoMap["last_modify_time"] = paramInfo.LastModifyTime
			}

			if paramInfo.StandbyRelated != nil {
				paramInfoMap["standby_related"] = paramInfo.StandbyRelated
			}

			if paramInfo.VersionRelationSet != nil {
				versionRelationSetList := []interface{}{}
				for _, versionRelationSet := range paramInfo.VersionRelationSet {
					versionRelationSetMap := map[string]interface{}{}

					if versionRelationSet.Name != nil {
						versionRelationSetMap["name"] = versionRelationSet.Name
					}

					if versionRelationSet.DBKernelVersion != nil {
						versionRelationSetMap["db_kernel_version"] = versionRelationSet.DBKernelVersion
					}

					if versionRelationSet.Value != nil {
						versionRelationSetMap["value"] = versionRelationSet.Value
					}

					if versionRelationSet.Unit != nil {
						versionRelationSetMap["unit"] = versionRelationSet.Unit
					}

					if versionRelationSet.Max != nil {
						versionRelationSetMap["max"] = versionRelationSet.Max
					}

					if versionRelationSet.Min != nil {
						versionRelationSetMap["min"] = versionRelationSet.Min
					}

					if versionRelationSet.EnumValue != nil {
						versionRelationSetMap["enum_value"] = versionRelationSet.EnumValue
					}

					versionRelationSetList = append(versionRelationSetList, versionRelationSetMap)
				}

				paramInfoMap["version_relation_set"] = versionRelationSetList
			}

			if paramInfo.SpecRelationSet != nil {
				specRelationSetList := []interface{}{}
				for _, specRelationSet := range paramInfo.SpecRelationSet {
					specRelationSetMap := map[string]interface{}{}

					if specRelationSet.Name != nil {
						specRelationSetMap["name"] = specRelationSet.Name
					}

					if specRelationSet.Memory != nil {
						specRelationSetMap["memory"] = specRelationSet.Memory
					}

					if specRelationSet.Value != nil {
						specRelationSetMap["value"] = specRelationSet.Value
					}

					if specRelationSet.Unit != nil {
						specRelationSetMap["unit"] = specRelationSet.Unit
					}

					if specRelationSet.Max != nil {
						specRelationSetMap["max"] = specRelationSet.Max
					}

					if specRelationSet.Min != nil {
						specRelationSetMap["min"] = specRelationSet.Min
					}

					if specRelationSet.EnumValue != nil {
						specRelationSetMap["enum_value"] = specRelationSet.EnumValue
					}

					specRelationSetList = append(specRelationSetList, specRelationSetMap)
				}

				paramInfoMap["spec_relation_set"] = specRelationSetList
			}

			ids = append(ids, helper.Int64ToStr(*paramInfo.ID))
			tmpList = append(tmpList, paramInfoMap)
		}

		_ = d.Set("param_info_set", tmpList)
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
