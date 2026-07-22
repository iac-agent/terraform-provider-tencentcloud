package scf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudScfTriggers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfTriggersRead,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Function 名称",
			},

			"namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Namespace. 默认值：默认值。",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "表示by 其中 字段 到 sort 返回 results. 有效值：add_time，mod_time. 默认值：mod_time。",
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "表示是否returned results 是 sorted 在 ascending 或 降序 有效值：ASC，DESC. 默认值：DESC。",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "* Qualifier:Function 版本，alias。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Fields 到 是 filtered. Up 到 10 conditions allowed.Values 的 名称: VpcId，SubnetId，ClsTopicId，ClsLogsetId，角色，CfsId，CfsMountInsId，Eip. Values 限制: 1.名称 options: 状态，Runtime，FunctionType，PublicNetStatus，AsyncRunEnable，TraceEnable. Values 限制: 20.当 名称 是 Runtime，CustomImage refers 到 镜像 类型 函数。",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "过滤器 值 的 字段。",
						},
					},
				},
			},

			"triggers": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Trigger 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否enable。",
						},
						"qualifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Function 版本 或 alias。",
						},
						"trigger_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Trigger 名称",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Trigger 类型",
						},
						"trigger_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Detailed 配置 的 触发器。",
						},
						"available_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "是否trigger 是 可用。",
						},
						"custom_argument": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom parameterNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"add_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Trigger 创建时间。",
						},
						"mod_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Trigger 最后修改时间。",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Minimum 资源 ID 触发器。",
						},
						"bind_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Trigger-Function binding 状态",
						},
						"trigger_attribute": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Trigger 类型 Two-way 表示 该 触发器 可以 是 manipulated 在 both consoles，while 一个-way 表示 该 触发器 可以 是 创建 仅 在 SCF Console。",
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

func dataSourceTencentCloudScfTriggersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_scf_triggers.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("function_name"); ok {
		paramMap["FunctionName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		paramMap["Namespace"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*scf.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := scf.Filter{}
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
		paramMap["Filters"] = tmpSet
	}

	service := ScfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var triggers []*scf.TriggerInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeScfTriggersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		triggers = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(triggers))
	tmpList := make([]map[string]interface{}, 0, len(triggers))

	if triggers != nil {
		for _, triggerInfo := range triggers {
			triggerInfoMap := map[string]interface{}{}

			if triggerInfo.Enable != nil {
				triggerInfoMap["enable"] = triggerInfo.Enable
			}

			if triggerInfo.Qualifier != nil {
				triggerInfoMap["qualifier"] = triggerInfo.Qualifier
			}

			if triggerInfo.TriggerName != nil {
				triggerInfoMap["trigger_name"] = triggerInfo.TriggerName
			}

			if triggerInfo.Type != nil {
				triggerInfoMap["type"] = triggerInfo.Type
			}

			if triggerInfo.TriggerDesc != nil {
				triggerInfoMap["trigger_desc"] = triggerInfo.TriggerDesc
			}

			if triggerInfo.AvailableStatus != nil {
				triggerInfoMap["available_status"] = triggerInfo.AvailableStatus
			}

			if triggerInfo.CustomArgument != nil {
				triggerInfoMap["custom_argument"] = triggerInfo.CustomArgument
			}

			if triggerInfo.AddTime != nil {
				triggerInfoMap["add_time"] = triggerInfo.AddTime
			}

			if triggerInfo.ModTime != nil {
				triggerInfoMap["mod_time"] = triggerInfo.ModTime
			}

			if triggerInfo.ResourceId != nil {
				triggerInfoMap["resource_id"] = triggerInfo.ResourceId
			}

			if triggerInfo.BindStatus != nil {
				triggerInfoMap["bind_status"] = triggerInfo.BindStatus
			}

			if triggerInfo.TriggerAttribute != nil {
				triggerInfoMap["trigger_attribute"] = triggerInfo.TriggerAttribute
			}

			ids = append(ids, *triggerInfo.TriggerName)
			tmpList = append(tmpList, triggerInfoMap)
		}

		_ = d.Set("triggers", tmpList)
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
