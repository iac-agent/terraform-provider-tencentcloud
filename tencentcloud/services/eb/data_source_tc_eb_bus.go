package eb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudEbBus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEbBusRead,
		Schema: map[string]*schema.Schema{
			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "According 到 其中 字段 到 sort 返回 results， following 字段 是 支持: `created_at` (创建时间)，`updated_at` (修改时间)。",
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Return results 在 ascending 或 降序，可选 值 ASC (ascending) 和 DESC (descending)。",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 conditions. upper 限制 的 Filters per 请求 是 10，和 upper 限制 的 过滤器.Values 5。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "一个或多个过滤值",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "名称 过滤器 键",
						},
					},
				},
			},

			"event_buses": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "事件 集合 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mod_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "更新时间。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event 集合 描述，unlimited character 类型，描述 within 200 字符。",
						},
						"add_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"event_bus_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event 集合 名称，其中 可以 仅 contain letters，numbers，underscores，hyphens，starts 使用 letter 和 结束 使用 数量 或 letter，2~60 字符。",
						},
						"event_bus_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "事件 bus ID。",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "事件 bus 类型",
						},
						"pay_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Billing 模式，note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"connection_briefs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Connector basic 信息，note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Connector 类型，note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Connector 状态，note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
									},
								},
							},
						},
						"target_briefs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Target brief 信息，note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Target ID。",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Target 类型",
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

func dataSourceTencentCloudEbBusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_eb_bus.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*eb.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := eb.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["Filters"] = tmpSet
	}

	service := EbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var eventBuses []*eb.EventBus

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeEbBusByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		eventBuses = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(eventBuses))
	tmpList := make([]map[string]interface{}, 0, len(eventBuses))

	if eventBuses != nil {
		for _, eventBus := range eventBuses {
			eventBusMap := map[string]interface{}{}

			if eventBus.ModTime != nil {
				eventBusMap["mod_time"] = eventBus.ModTime
			}

			if eventBus.Description != nil {
				eventBusMap["description"] = eventBus.Description
			}

			if eventBus.AddTime != nil {
				eventBusMap["add_time"] = eventBus.AddTime
			}

			if eventBus.EventBusName != nil {
				eventBusMap["event_bus_name"] = eventBus.EventBusName
			}

			if eventBus.EventBusId != nil {
				eventBusMap["event_bus_id"] = eventBus.EventBusId
			}

			if eventBus.Type != nil {
				eventBusMap["type"] = eventBus.Type
			}

			if eventBus.PayMode != nil {
				eventBusMap["pay_mode"] = eventBus.PayMode
			}

			if eventBus.ConnectionBriefs != nil {
				connectionBriefsList := []interface{}{}
				for _, connectionBriefs := range eventBus.ConnectionBriefs {
					connectionBriefsMap := map[string]interface{}{}

					if connectionBriefs.Type != nil {
						connectionBriefsMap["type"] = connectionBriefs.Type
					}

					if connectionBriefs.Status != nil {
						connectionBriefsMap["status"] = connectionBriefs.Status
					}

					connectionBriefsList = append(connectionBriefsList, connectionBriefsMap)
				}

				eventBusMap["connection_briefs"] = connectionBriefsList
			}

			if eventBus.TargetBriefs != nil {
				targetBriefsList := []interface{}{}
				for _, targetBriefs := range eventBus.TargetBriefs {
					targetBriefsMap := map[string]interface{}{}

					if targetBriefs.TargetId != nil {
						targetBriefsMap["target_id"] = targetBriefs.TargetId
					}

					if targetBriefs.Type != nil {
						targetBriefsMap["type"] = targetBriefs.Type
					}

					targetBriefsList = append(targetBriefsList, targetBriefsMap)
				}

				eventBusMap["target_briefs"] = targetBriefsList
			}

			ids = append(ids, *eventBus.EventBusId)
			tmpList = append(tmpList, eventBusMap)
		}

		_ = d.Set("event_buses", tmpList)
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
