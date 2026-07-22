package audit

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cloudaudit "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudAuditEvents() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAuditEventsRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Start 时间戳 （秒） (不能 是 90 days after 当前 时间)。",
			},

			"end_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "End 时间戳 （秒） ( 时间 范围 对于 查询 是 less 比 30 days)。",
			},

			"max_results": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Max 数量 返回 logs (up 到 50)。",
			},

			"lookup_attributes": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Search condition. 有效 值: `RequestId`, `EventName`, `ActionType` (write/read), `PrincipalId` (sub-account), `ResourceType`, `ResourceName`, `AccessKeyId`, `SensitiveAction`, `ApiErrorCode`, `CamErrorCode`, 和 `Tags` (Format 的 AttributeValue: [{\"键\":\"*\",\"值\":\"*\"}]).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attribute_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "有效值：RequestId，EventName，ReadOnly，用户名，ResourceType，ResourceName，AccessKeyId，和 EventId\nNote: `null` 可能 是 返回 对于 此 字段，indicating 该 无 有效 值 可以 是 获取。",
						},
						"attribute_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "值 的 `AttributeValue`\nNote: `null` 可能 是 返回 对于 此 字段，indicating 该 无 有效 值 可以 是 获取。",
						},
					},
				},
			},

			"is_return_location": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "是否return IP location. `1`: yes，`0`: 无。",
			},

			"events": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Logset. 注意: `null` 可能 是 返回 对于 此 字段，indicating 该 无 有效 值 可以 是 获取。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Log ID。",
						},
						"username": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "用户名",
						},
						"event_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Event Time。",
						},
						"cloud_audit_event": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Log details。",
						},
						"resource_type_cn": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "描述 资源类型 在 Chinese (please 使用 此 字段 作为 必填; 如果 您 是 使用 other languages，ignore 此 字段)。",
						},
						"error_code": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Authentication 错误码",
						},
						"event_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "事件名称",
						},
						"secret_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "证书 ID\nNote: `null` 可能 是 返回 对于 此 字段，indicating 该 无 有效 值 可以 是 获取。",
						},
						"event_source": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "请求来源",
						},
						"request_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "请求 ID",
						},
						"resource_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resource 地域",
						},
						"account_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Root 账号 ID。",
						},
						"source_ip_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "来源 IP\nNote: `null` 可能 是 返回 对于 此 字段，indicating 该 无 有效 值 可以 是 获取。",
						},
						"event_name_cn": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "描述 事件名称 在 Chinese (please 使用 此 字段 作为 必填; 如果 您 是 使用 other languages，ignore 此 字段)。",
						},
						"resources": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Resource pair。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "资源类型",
									},
									"resource_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "资源名称\nNote: `null` 可能 是 返回 对于 此 字段，indicating 该 无 有效 值 可以 是 获取。",
									},
								},
							},
						},
						"event_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Event 地域",
						},
						"location": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IP location。",
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

func dataSourceTencentCloudAuditEventsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_audit_event.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := AuditService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})

	if v, ok := d.GetOkExists("max_results"); ok {
		paramMap["MaxResults"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("start_time"); ok {
		paramMap["StartTime"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("end_time"); ok {
		paramMap["EndTime"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("lookup_attributes"); ok {
		lookupAttributesSet := v.([]interface{})
		tmpSet := make([]*cloudaudit.LookupAttribute, 0, len(lookupAttributesSet))
		for _, item := range lookupAttributesSet {
			lookupAttributesMap := item.(map[string]interface{})
			lookupAttribute := cloudaudit.LookupAttribute{}
			if v, ok := lookupAttributesMap["attribute_key"]; ok {
				lookupAttribute.AttributeKey = helper.String(v.(string))
			}
			if v, ok := lookupAttributesMap["attribute_value"]; ok {
				lookupAttribute.AttributeValue = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &lookupAttribute)
		}
		paramMap["LookupAttributes"] = tmpSet
	}

	if v, ok := d.GetOkExists("is_return_location"); ok {
		paramMap["IsReturnLocation"] = helper.IntUint64(v.(int))
	}

	var respData []*cloudaudit.Event
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAuditEventByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(respData))
	eventsList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, events := range respData {
			eventsMap := map[string]interface{}{}
			println(*events.EventId)
			if events.EventId != nil {
				eventsMap["event_id"] = events.EventId
			}

			if events.Username != nil {
				eventsMap["username"] = events.Username
			}

			if events.EventTime != nil {
				eventsMap["event_time"] = events.EventTime
			}

			if events.CloudAuditEvent != nil {
				eventsMap["cloud_audit_event"] = events.CloudAuditEvent
			}

			if events.ResourceTypeCn != nil {
				eventsMap["resource_type_cn"] = events.ResourceTypeCn
			}

			if events.ErrorCode != nil {
				eventsMap["error_code"] = events.ErrorCode
			}

			if events.EventName != nil {
				eventsMap["event_name"] = events.EventName
			}

			if events.SecretId != nil {
				eventsMap["secret_id"] = events.SecretId
			}

			if events.EventSource != nil {
				eventsMap["event_source"] = events.EventSource
			}

			if events.RequestID != nil {
				eventsMap["request_id"] = events.RequestID
			}

			if events.ResourceRegion != nil {
				eventsMap["resource_region"] = events.ResourceRegion
			}

			if events.AccountID != nil {
				eventsMap["account_id"] = events.AccountID
			}

			if events.SourceIPAddress != nil {
				eventsMap["source_ip_address"] = events.SourceIPAddress
			}

			if events.EventNameCn != nil {
				eventsMap["event_name_cn"] = events.EventNameCn
			}

			resourcesMap := map[string]interface{}{}

			if events.Resources != nil {
				if events.Resources.ResourceType != nil {
					resourcesMap["resource_type"] = events.Resources.ResourceType
				}

				if events.Resources.ResourceName != nil {
					resourcesMap["resource_name"] = events.Resources.ResourceName
				}

				eventsMap["resources"] = []interface{}{resourcesMap}
			}

			if events.EventRegion != nil {
				eventsMap["event_region"] = events.EventRegion
			}

			if events.Location != nil {
				eventsMap["location"] = events.Location
			}

			eventsList = append(eventsList, eventsMap)
			ids = append(ids, *events.EventId)
		}

		_ = d.Set("events", eventsList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), eventsList); e != nil {
			return e
		}
	}

	return nil
}
