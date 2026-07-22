package dbbrain

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbbrainDiagHistory() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainDiagHistoryRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start 时间, such 作为 `2019-09-10 12:13:14`.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End 时间, such 作为 `2019-09-11 12:13:14`, 间隔 between end 时间 和 start 时间 可以 是 up 到 2 days.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product 类型, 支持 值 include: `mysql` - 云 数据库 MySQL, `cynosdb` - 云 数据库 CynosDB 对于 MySQL, 默认值 是 `mysql`.",
			},

			"events": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Event 描述.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"diag_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Diagnostic 类型.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End Time.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "start Time.",
						},
						"event_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Event 唯一 ID.",
						},
						"severity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "severity. severity 是 divided into 5 levels, according 到 degree 的 impact 从 high 到 low: 1: Fatal, 2: Serious, 3: Warning, 4: Prompt, 5: Healthy.",
						},
						"outline": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Diagnostic summary.",
						},
						"diag_item": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description 的 diagnostic item.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID.",
						},
						"metric": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "reserved text. 注意: 此 字段 可能 返回 null, indicating 该 无 有效 值 可以 是 获取.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainDiagHistoryRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_diag_history.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["start_time"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["end_time"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var events []*dbbrain.DiagHistoryEventItem

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainDiagHistoryByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		events = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(events))
	tmpList := make([]map[string]interface{}, 0, len(events))

	if events != nil {
		for _, diagHistoryEventItem := range events {
			diagHistoryEventItemMap := map[string]interface{}{}

			if diagHistoryEventItem.DiagType != nil {
				diagHistoryEventItemMap["diag_type"] = diagHistoryEventItem.DiagType
			}

			if diagHistoryEventItem.EndTime != nil {
				diagHistoryEventItemMap["end_time"] = diagHistoryEventItem.EndTime
			}

			if diagHistoryEventItem.StartTime != nil {
				diagHistoryEventItemMap["start_time"] = diagHistoryEventItem.StartTime
			}

			if diagHistoryEventItem.EventId != nil {
				diagHistoryEventItemMap["event_id"] = diagHistoryEventItem.EventId
			}

			if diagHistoryEventItem.Severity != nil {
				diagHistoryEventItemMap["severity"] = diagHistoryEventItem.Severity
			}

			if diagHistoryEventItem.Outline != nil {
				diagHistoryEventItemMap["outline"] = diagHistoryEventItem.Outline
			}

			if diagHistoryEventItem.DiagItem != nil {
				diagHistoryEventItemMap["diag_item"] = diagHistoryEventItem.DiagItem
			}

			if diagHistoryEventItem.InstanceId != nil {
				diagHistoryEventItemMap["instance_id"] = diagHistoryEventItem.InstanceId
			}

			if diagHistoryEventItem.Metric != nil {
				diagHistoryEventItemMap["metric"] = diagHistoryEventItem.Metric
			}

			if diagHistoryEventItem.Region != nil {
				diagHistoryEventItemMap["region"] = diagHistoryEventItem.Region
			}

			ids = append(ids, *diagHistoryEventItem.InstanceId)
			tmpList = append(tmpList, diagHistoryEventItemMap)
		}

		_ = d.Set("events", tmpList)
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
