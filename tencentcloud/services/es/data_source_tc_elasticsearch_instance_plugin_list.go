package es

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	elasticsearch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudElasticsearchInstancePluginList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudElasticsearchInstancePluginListRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "顺序 field. 有效值：`pluginName`。",
			},

			"order_by_type": {
				Optional: true,
				Type:     schema.TypeString,
				Description: "Order type. Valid values:\n" +
					"- asc: Ascending asc\n" +
					"- desc: Descending Desc.",
			},

			"plugin_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Plugin 类型 有效值：`0`: System plugin。",
			},

			"plugin_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Plugin information list。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plugin_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plugin 名称",
						},
						"plugin_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plugin 版本",
						},
						"plugin_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plugin 描述",
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "Plugin status. Valid values:\n" +
								"- `-2` has been uninstalled\n" +
								"- `-1` has been installed in\n" +
								"- `0` installation.",
						},
						"removable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否plug-in can be uninstalled。",
						},
						"plugin_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Plugin 类型 有效值：`0`: System plugin。",
						},
						"plugin_update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plugin 更新时间。",
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

func dataSourceTencentCloudElasticsearchInstancePluginListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_elasticsearch_instance_plugin_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by_type"); ok {
		paramMap["OrderByType"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("plugin_type"); v != nil {
		paramMap["PluginType"] = helper.IntInt64(v.(int))
	}

	service := ElasticsearchService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var pluginList []*elasticsearch.DescribeInstancePluginInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeElasticsearchInstancePluginListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		pluginList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(pluginList))
	tmpList := make([]map[string]interface{}, 0, len(pluginList))

	if pluginList != nil {
		for _, describeInstancePluginInfo := range pluginList {
			describeInstancePluginInfoMap := map[string]interface{}{}

			if describeInstancePluginInfo.PluginName != nil {
				describeInstancePluginInfoMap["plugin_name"] = describeInstancePluginInfo.PluginName
			}

			if describeInstancePluginInfo.PluginVersion != nil {
				describeInstancePluginInfoMap["plugin_version"] = describeInstancePluginInfo.PluginVersion
			}

			if describeInstancePluginInfo.PluginDesc != nil {
				describeInstancePluginInfoMap["plugin_desc"] = describeInstancePluginInfo.PluginDesc
			}

			if describeInstancePluginInfo.Status != nil {
				describeInstancePluginInfoMap["status"] = describeInstancePluginInfo.Status
			}

			if describeInstancePluginInfo.Removable != nil {
				describeInstancePluginInfoMap["removable"] = describeInstancePluginInfo.Removable
			}

			if describeInstancePluginInfo.PluginType != nil {
				describeInstancePluginInfoMap["plugin_type"] = describeInstancePluginInfo.PluginType
			}

			if describeInstancePluginInfo.PluginUpdateTime != nil {
				describeInstancePluginInfoMap["plugin_update_time"] = describeInstancePluginInfo.PluginUpdateTime
			}

			ids = append(ids, *describeInstancePluginInfo.PluginName)
			tmpList = append(tmpList, describeInstancePluginInfoMap)
		}

		_ = d.Set("plugin_list", tmpList)
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
