package monitor

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMonitorAlarmAllNamespaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMonitorAlarmAllNamespacesRead,
		Schema: map[string]*schema.Schema{
			"scene_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Currently，仅 ST_ALARM=告警 类型 是 filtered based 在 usage scenarios。",
			},

			"module": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Fixed 值，作为 `监控`。",
			},

			"monitor_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "过滤器 based 在 监控 类型，do 不 fill 在 默认值，check all types MT_QCE=云 product 监控。",
			},

			"ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "过滤器 based 在 ID 的 命名空间 without filling 在 默认值 查询 对于 all。",
			},

			"qce_namespaces_new": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Types 的 告警 strategies 对于 云 products。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace labeling。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace 名称",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace 值",
						},
						"product_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Product 名称",
						},
						"config": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Configuration 信息。",
						},
						"available_regions": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "列表 支持 regions。",
						},
						"sort_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sort ID。",
						},
						"dashboard_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique representation 在 仪表盘。",
						},
					},
				},
			},

			"custom_namespaces_new": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Other 告警 strategy types 是 currently 不 支持。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace labeling。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace 名称",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace 值",
						},
						"product_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Product 名称",
						},
						"config": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Configuration 信息。",
						},
						"available_regions": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "列表 支持 regions。",
						},
						"sort_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sort ID。",
						},
						"dashboard_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique representation 在 仪表盘。",
						},
					},
				},
			},

			"common_namespaces": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "General 告警 strategy types (包括: 应用 performance 监控，front-end performance 监控，云 dial testing)。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace labeling。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace 名称",
						},
						"monitor_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Monitoring 类型",
						},
						"dimensions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Dimension Information。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Dimension 键 identifier，backend English 名称",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Dimension 键 名称，Chinese 和 English frontend display 名称",
									},
									"is_required": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "必填 或 不。",
									},
									"operators": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "列表 支持 operators。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "操作者 identification。",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "操作者 Display 名称",
												},
											},
										},
									},
									"is_multiple": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Do 您 support 多个 selections。",
									},
									"is_mutable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Can I modify 它 after creation。",
									},
									"is_visible": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否display 到 users。",
									},
									"can_filter_policy": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Can 它 是 用于filter 策略 列表。",
									},
									"can_filter_history": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Can 它 是 用于filter 告警 history。",
									},
									"can_group_by": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Can 它 是 使用 作为 aggregation dimension。",
									},
									"must_group_by": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Must 它 是 使用 作为 aggregation dimension。",
									},
									"show_value_replace": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "键 到 replace 在 front-end translation。",
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

func dataSourceTencentCloudMonitorAlarmAllNamespacesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_alarm_all_namespaces.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("scene_type"); ok {
		paramMap["SceneType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("module"); ok {
		paramMap["Module"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("monitor_types"); ok {
		monitorTypesSet := v.(*schema.Set).List()
		paramMap["MonitorTypes"] = helper.InterfacesStringsPoint(monitorTypesSet)
	}

	if v, ok := d.GetOk("ids"); ok {
		idsSet := v.(*schema.Set).List()
		paramMap["Ids"] = helper.InterfacesStringsPoint(idsSet)
	}

	service := MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var qceNamespacesNew []*monitor.CommonNamespace
	var customNamespacesNew []*monitor.CommonNamespace
	var commonNamespaces []*monitor.CommonNamespaceNew
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		qce, custom, common, e := service.DescribeMonitorAlarmAllNamespacesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		qceNamespacesNew = qce
		customNamespacesNew = custom
		commonNamespaces = common
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	if qceNamespacesNew != nil {
		tmpList := make([]map[string]interface{}, 0)
		for _, commonNamespace := range qceNamespacesNew {
			commonNamespaceMap := map[string]interface{}{}

			if commonNamespace.Id != nil {
				commonNamespaceMap["id"] = commonNamespace.Id
			}

			if commonNamespace.Name != nil {
				commonNamespaceMap["name"] = commonNamespace.Name
			}

			if commonNamespace.Value != nil {
				commonNamespaceMap["value"] = commonNamespace.Value
			}

			if commonNamespace.ProductName != nil {
				commonNamespaceMap["product_name"] = commonNamespace.ProductName
			}

			if commonNamespace.Config != nil {
				commonNamespaceMap["config"] = commonNamespace.Config
			}

			if commonNamespace.AvailableRegions != nil {
				commonNamespaceMap["available_regions"] = commonNamespace.AvailableRegions
			}

			if commonNamespace.SortId != nil {
				commonNamespaceMap["sort_id"] = commonNamespace.SortId
			}

			if commonNamespace.DashboardId != nil {
				commonNamespaceMap["dashboard_id"] = commonNamespace.DashboardId
			}

			ids = append(ids, *commonNamespace.Id)
			tmpList = append(tmpList, commonNamespaceMap)
		}

		_ = d.Set("qce_namespaces_new", tmpList)
	}

	if customNamespacesNew != nil {
		tmpList := make([]map[string]interface{}, 0)
		for _, commonNamespace := range customNamespacesNew {
			commonNamespaceMap := map[string]interface{}{}

			if commonNamespace.Id != nil {
				commonNamespaceMap["id"] = commonNamespace.Id
			}

			if commonNamespace.Name != nil {
				commonNamespaceMap["name"] = commonNamespace.Name
			}

			if commonNamespace.Value != nil {
				commonNamespaceMap["value"] = commonNamespace.Value
			}

			if commonNamespace.ProductName != nil {
				commonNamespaceMap["product_name"] = commonNamespace.ProductName
			}

			if commonNamespace.Config != nil {
				commonNamespaceMap["config"] = commonNamespace.Config
			}

			if commonNamespace.AvailableRegions != nil {
				commonNamespaceMap["available_regions"] = commonNamespace.AvailableRegions
			}

			if commonNamespace.SortId != nil {
				commonNamespaceMap["sort_id"] = commonNamespace.SortId
			}

			if commonNamespace.DashboardId != nil {
				commonNamespaceMap["dashboard_id"] = commonNamespace.DashboardId
			}

			ids = append(ids, *commonNamespace.Id)
			tmpList = append(tmpList, commonNamespaceMap)
		}

		_ = d.Set("custom_namespaces_new", tmpList)
	}

	if commonNamespaces != nil {
		tmpList := make([]map[string]interface{}, 0)
		for _, commonNamespaceNew := range commonNamespaces {
			commonNamespaceNewMap := map[string]interface{}{}

			if commonNamespaceNew.Id != nil {
				commonNamespaceNewMap["id"] = commonNamespaceNew.Id
			}

			if commonNamespaceNew.Name != nil {
				commonNamespaceNewMap["name"] = commonNamespaceNew.Name
			}

			if commonNamespaceNew.MonitorType != nil {
				commonNamespaceNewMap["monitor_type"] = commonNamespaceNew.MonitorType
			}

			if commonNamespaceNew.Dimensions != nil {
				dimensionsList := []interface{}{}
				for _, dimensions := range commonNamespaceNew.Dimensions {
					dimensionsMap := map[string]interface{}{}

					if dimensions.Key != nil {
						dimensionsMap["key"] = dimensions.Key
					}

					if dimensions.Name != nil {
						dimensionsMap["name"] = dimensions.Name
					}

					if dimensions.IsRequired != nil {
						dimensionsMap["is_required"] = dimensions.IsRequired
					}

					if dimensions.Operators != nil {
						operatorsList := []interface{}{}
						for _, operators := range dimensions.Operators {
							operatorsMap := map[string]interface{}{}

							if operators.Id != nil {
								operatorsMap["id"] = operators.Id
							}

							if operators.Name != nil {
								operatorsMap["name"] = operators.Name
							}

							operatorsList = append(operatorsList, operatorsMap)
						}

						dimensionsMap["operators"] = operatorsList
					}

					if dimensions.IsMultiple != nil {
						dimensionsMap["is_multiple"] = dimensions.IsMultiple
					}

					if dimensions.IsMutable != nil {
						dimensionsMap["is_mutable"] = dimensions.IsMutable
					}

					if dimensions.IsVisible != nil {
						dimensionsMap["is_visible"] = dimensions.IsVisible
					}

					if dimensions.CanFilterPolicy != nil {
						dimensionsMap["can_filter_policy"] = dimensions.CanFilterPolicy
					}

					if dimensions.CanFilterHistory != nil {
						dimensionsMap["can_filter_history"] = dimensions.CanFilterHistory
					}

					if dimensions.CanGroupBy != nil {
						dimensionsMap["can_group_by"] = dimensions.CanGroupBy
					}

					if dimensions.MustGroupBy != nil {
						dimensionsMap["must_group_by"] = dimensions.MustGroupBy
					}

					if dimensions.ShowValueReplace != nil {
						dimensionsMap["show_value_replace"] = dimensions.ShowValueReplace
					}

					dimensionsList = append(dimensionsList, dimensionsMap)
				}

				commonNamespaceNewMap["dimensions"] = dimensionsList
			}

			ids = append(ids, *commonNamespaceNew.Id)
			tmpList = append(tmpList, commonNamespaceNewMap)
		}

		_ = d.Set("common_namespaces", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
