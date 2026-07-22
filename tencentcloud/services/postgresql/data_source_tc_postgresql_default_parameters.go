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
				Description: "The major database 版本 number，such as 11，12，13。",
			},

			"db_engine": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Database engine，such as postgresql，mssql_compatible。",
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
							Description: "Parameter IDNote: this field may return `null`，indicating that no valid values can be obtained。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter nameNote: this field may return `null`，indicating that no valid values can be obtained。",
						},
						"param_value_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "值 类型 parameter. 有效值：`integer`，`real` (floating-point)，`bool`，`enum`，`mutil_enum` (this 类型 parameter can be set to multiple enumerated values).For an `integer` or `real` parameter，the `Min` field represents the minimum 值 and the `Max` field the maximum 值 For a `bool` parameter，the valid values include `true` and `false`; For an `enum` or `mutil_enum` parameter，the `EnumValue` field represents the valid values.Note: this field may return `null`，indicating that no valid values can be obtained。",
						},
						"unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unit of the parameter 值 If the parameter has no unit，this field will return null.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"default_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "默认值 of the parameter，which is returned as a stringNote: this field may return `null`，indicating that no valid values can be obtained。",
						},
						"current_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "当前值 of the parameter，which is returned as a stringNote: this field may return `null`，indicating that no valid values can be obtained。",
						},
						"max": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The maximum 值 of the `integer` or `real` parameterNote: this field may return `null`，indicating that no valid values can be obtained。",
						},
						"enum_value": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "值 range of the enum parameterNote: this field may return `null`，indicating that no valid values can be obtained。",
						},
						"min": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The minimum 值 of the `integer` or `real` parameterNote: this field may return `null`，indicating that no valid values can be obtained。",
						},
						"param_description_ch": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter 描述 in ChineseNote: this field may return `null`，indicating that no valid values can be obtained。",
						},
						"param_description_en": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter 描述 in EnglishNote: this field may return `null`，indicating that no valid values can be obtained。",
						},
						"need_reboot": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否restart the instance for the modified parameter to take effect. 有效值：`true` (yes)，`false` (no)Note: this field may return `null`，indicating that no valid values can be obtained。",
						},
						"classification_cn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter category in ChineseNote: this field may return `null`，indicating that no valid values can be obtained。",
						},
						"classification_en": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter category in EnglishNote: this field may return `null`，indicating that no valid values can be obtained。",
						},
						"spec_related": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否parameter is related to specifications. 有效值：`true` (yes)，`false` (no)Note: this field may return `null`，indicating that no valid values can be obtained。",
						},
						"advanced": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否为a 键 parameter. 有效值：`true` (yes，and modifying it may affect instance performance)，`false` (no)Note: this field may return `null`，indicating that no valid values can be obtained。",
						},
						"last_modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The 最后修改时间 of the parameterNote: this field may return `null`，indicating that no valid values can be obtained。",
						},
						"standby_related": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Primary-standby constraint. 有效值：`0` (no constraint)，`1` (The parameter 值 of the standby server must be greater than that of the primary server)，`2` (The parameter 值 of the primary server must be greater than that of the standby server.)注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"version_relation_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Associated parameter 版本 information，which refers to the detailed parameter information of the kernel 版本注意：此字段可能返回 null，表示无法获取有效值。",
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
										Description: "The kernel 版本 that corresponds to the parameter information注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Default parameter 值 under the kernel 版本 and specification of the instance注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Unit of the parameter 值 If the parameter has no unit，this field will return null.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"max": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "The maximum 值 of the `integer` or `real` parameter注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"min": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "The minimum 值 of the `integer` or `real` parameter注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"enum_value": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "值 range of the enum parameter注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"spec_relation_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Associated parameter specification information，which refers to the detailed parameter information of the specifications.注意：此字段可能返回 null，表示无法获取有效值。",
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
										Description: "The specification that corresponds to the parameter information注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The default parameter 值 under this specification注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Unit of the parameter 值 If the parameter has no unit，this field will return null.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"max": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "The maximum 值 of the `integer` or `real` parameter注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"min": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "The minimum 值 of the `integer` or `real` parameter注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"enum_value": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "值 range of the enum parameter注意：此字段可能返回 null，表示无法获取有效值。",
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
