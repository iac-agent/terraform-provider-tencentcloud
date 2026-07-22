package cynosdb

import (
	"context"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCynosdbParamTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbParamTemplatesRead,
		Schema: map[string]*schema.Schema{
			"engine_versions": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "数据库引擎版本号。",
			},

			"template_names": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "模板的名称列表。",
			},

			"template_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "模板的 ID 列表。",
			},

			"db_modes": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "数据库模式，可选值：NORMAL、SERVERLESS。",
			},

			"offset": {
				Optional:    true,
				Default:     0,
				Type:        schema.TypeInt,
				Description: "页面偏移量。",
			},

			"limit": {
				Optional:    true,
				Default:     10,
				Type:        schema.TypeInt,
				Description: "查询限制。",
			},

			"products": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "查询模板对应的产品类型。",
			},

			"template_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "模板类型。",
			},

			"engine_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "发动机类型。",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "返回结果的排序字段。",
			},

			"order_direction": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "按（升序、降序）排序。",
			},

			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "参数模板信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "模板的ID。",
						},
						"template_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "模板的名称。",
						},
						"template_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "模板的描述。",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "引擎版本。",
						},
						"db_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "数据库模式，可选值：NORMAL、SERVERLESS。",
						},
						"param_info_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "参数模板详情。注意：该字段可能返回null，表示取不到有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"current_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "当前值。",
									},
									"default": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "默认值。",
									},
									"enum_value": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "当参数类型为enum时，一组可选的值类型。注意：该字段可能返回null，表示取不到有效值。",
									},
									"max": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "参数类型为float/integer时的最大值。注意：该字段可能返回null，表示取不到有效值。",
									},
									"min": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "参数类型为float/integer时的最小值。注意：该字段可能返回null，表示取不到有效值。",
									},
									"param_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "参数名称。",
									},
									"need_reboot": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "是否重启。",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "参数说明。",
									},
									"param_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "参数类型：整数/浮点/字符串/枚举。",
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

func dataSourceTencentCloudCynosdbParamTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cynosdb_param_templates.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("engine_versions"); ok {
		engineVersionsSet := v.(*schema.Set).List()
		paramMap["engine_versions"] = helper.InterfacesStringsPoint(engineVersionsSet)
	}

	if v, ok := d.GetOk("template_names"); ok {
		templateNamesSet := v.(*schema.Set).List()
		paramMap["template_names"] = helper.InterfacesStringsPoint(templateNamesSet)
	}

	if v, ok := d.GetOk("template_ids"); ok {
		templateIds := make([]*int64, 0)
		for _, item := range v.(*schema.Set).List() {
			templateIds = append(templateIds, helper.IntInt64(item.(int)))
		}
		paramMap["template_ids"] = templateIds

	}

	if v, ok := d.GetOk("db_modes"); ok {
		dbModesSet := v.(*schema.Set).List()
		paramMap["db_modes"] = helper.InterfacesStringsPoint(dbModesSet)
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["offset"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["limit"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("products"); ok {
		productsSet := v.(*schema.Set).List()
		paramMap["products"] = helper.InterfacesStringsPoint(productsSet)
	}

	if v, ok := d.GetOk("template_types"); ok {
		templateTypesSet := v.(*schema.Set).List()
		paramMap["template_types"] = helper.InterfacesStringsPoint(templateTypesSet)
	}

	if v, ok := d.GetOk("engine_types"); ok {
		engineTypesSet := v.(*schema.Set).List()
		paramMap["engine_types"] = helper.InterfacesStringsPoint(engineTypesSet)
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["order_by"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_direction"); ok {
		paramMap["order_direction"] = helper.String(v.(string))
	}

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var items []*cynosdb.ParamTemplateListInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbParamTemplatesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		items = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(items))
	tmpList := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		ids = append(ids, strconv.FormatInt(*item.Id, 10))
		itemMap := make(map[string]interface{})
		itemMap["id"] = item.Id
		itemMap["template_name"] = item.TemplateName
		itemMap["template_description"] = item.TemplateDescription
		itemMap["engine_version"] = item.EngineVersion
		itemMap["db_mode"] = item.DbMode
		paramInfos := make([]map[string]interface{}, 0)
		if item.ParamInfoSet != nil {
			for _, paramInfo := range item.ParamInfoSet {
				paramInfoMap := make(map[string]interface{})
				paramInfoMap["current_value"] = paramInfo.CurrentValue
				paramInfoMap["default"] = paramInfo.Default
				enumValues := make([]string, 0)
				if paramInfo.EnumValue != nil {
					for _, enumValue := range paramInfo.EnumValue {
						enumValues = append(enumValues, *enumValue)
					}
				}
				paramInfoMap["enum_value"] = enumValues
				paramInfoMap["max"] = paramInfo.Max
				paramInfoMap["min"] = paramInfo.Min
				paramInfoMap["param_name"] = paramInfo.ParamName
				paramInfoMap["need_reboot"] = paramInfo.NeedReboot
				paramInfoMap["description"] = paramInfo.Description
				paramInfoMap["param_type"] = paramInfo.ParamType

				paramInfos = append(paramInfos, paramInfoMap)
			}
		}
		itemMap["param_info_set"] = paramInfos
		itemMap["id"] = item.Id
		tmpList = append(tmpList, itemMap)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("items", tmpList)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
