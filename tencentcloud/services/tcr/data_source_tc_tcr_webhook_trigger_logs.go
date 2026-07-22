package tcr

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTcrWebhookTriggerLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcrWebhookTriggerLogsRead,
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID。",
			},

			"namespace": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "命名空间。",
			},

			"trigger_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "触发器 ID。",
			},

			"logs": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "日志 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "日志 ID。",
						},
						"trigger_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "触发器 ID。",
						},
						"event_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "事件 类型",
						},
						"notify_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "通知 类型",
						},
						"detail": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "webhook 触发器 detail。",
						},
						"creation_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "更新时间。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态",
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签描述列表",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudTcrWebhookTriggerLogsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tcr_webhook_trigger_logs.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("registry_id"); ok {
		paramMap["registry_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		paramMap["namespace"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("trigger_id"); ok {
		paramMap["trigger_id"] = helper.IntInt64(v.(int))
	}

	service := TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var logs []*tcr.WebhookTriggerLog

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTcrWebhookTriggerLogByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		logs = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(logs))
	tmpList := make([]map[string]interface{}, 0, len(logs))

	if logs != nil {
		for _, webhookTriggerLog := range logs {
			webhookTriggerLogMap := map[string]interface{}{}

			if webhookTriggerLog.Id != nil {
				webhookTriggerLogMap["id"] = webhookTriggerLog.Id
			}

			if webhookTriggerLog.TriggerId != nil {
				webhookTriggerLogMap["trigger_id"] = webhookTriggerLog.TriggerId
			}

			if webhookTriggerLog.EventType != nil {
				webhookTriggerLogMap["event_type"] = webhookTriggerLog.EventType
			}

			if webhookTriggerLog.NotifyType != nil {
				webhookTriggerLogMap["notify_type"] = webhookTriggerLog.NotifyType
			}

			if webhookTriggerLog.Detail != nil {
				webhookTriggerLogMap["detail"] = webhookTriggerLog.Detail
			}

			if webhookTriggerLog.CreationTime != nil {
				webhookTriggerLogMap["creation_time"] = webhookTriggerLog.CreationTime
			}

			if webhookTriggerLog.UpdateTime != nil {
				webhookTriggerLogMap["update_time"] = webhookTriggerLog.UpdateTime
			}

			if webhookTriggerLog.Status != nil {
				webhookTriggerLogMap["status"] = webhookTriggerLog.Status
			}

			ids = append(ids, helper.Int64ToStr(*webhookTriggerLog.Id))
			tmpList = append(tmpList, webhookTriggerLogMap)
		}

		_ = d.Set("logs", tmpList)
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
