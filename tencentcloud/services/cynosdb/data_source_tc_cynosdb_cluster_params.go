package cynosdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCynosdbClusterParams() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbClusterParamsRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "集群的ID。",
			},

			"param_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "参数名称。",
			},

			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "实例参数列表。注意：该字段可能返回null，表示取不到有效值。",
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
							Description: "当参数为enum/字符串/bool时，可选值列表。注意：该字段可能返回null，表示取不到有效值。",
						},
						"max": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "参数类型为float/integer时的最大值。",
						},
						"min": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "参数类型为float/integer时的最小值。",
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
						"param_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "参数类型：整数/浮点/字符串/枚举/布尔。",
						},
						"match_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "当参数类型为字符串时，使用匹配类型、multiVal、regex。",
						},
						"match_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "匹配目标值，multiVal时，每个key除以`;`。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "参数说明。",
						},
						"is_global": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否为全局参数。注意：该字段可能返回null，表示取不到有效值。",
						},
						"is_func": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否是一个函数。注意：该字段可能返回null，表示取不到有效值。",
						},
						"func": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "函数注意：该字段可能返回null，表示取不到有效值。",
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

func dataSourceTencentCloudCynosdbClusterParamsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cynosdb_cluster_params.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var clusterId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("param_name"); ok {
		paramMap["param_name"] = v.(string)
	}

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var items []*cynosdb.ParamInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClusterParamsByFilter(ctx, clusterId, paramMap)
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
		ids = append(ids, *item.ParamName)
		itemMap := make(map[string]interface{})
		itemMap["current_value"] = item.CurrentValue
		itemMap["default"] = item.Default
		itemMap["max"] = item.Max
		itemMap["min"] = item.Min
		itemMap["param_name"] = item.ParamName
		itemMap["need_reboot"] = item.NeedReboot
		itemMap["param_type"] = item.ParamType
		itemMap["match_type"] = item.MatchType
		itemMap["match_value"] = item.MatchValue
		itemMap["description"] = item.Description
		itemMap["is_global"] = item.IsGlobal
		itemMap["is_func"] = item.IsFunc
		itemMap["func"] = item.Func
		enumValues := make([]string, 0)
		if item.EnumValue != nil {
			for _, enumValueItem := range item.EnumValue {
				enumValues = append(enumValues, *enumValueItem)
			}
			itemMap["enum_value"] = enumValues
		}
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
