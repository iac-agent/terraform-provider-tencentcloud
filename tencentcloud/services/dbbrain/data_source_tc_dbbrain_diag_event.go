package dbbrain

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbbrainDiagEvent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainDiagEventRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "isntance ID.",
			},

			"event_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Event ID. Obtain 它 through `Get 实例 Diagnosis History DescribeDBDiagHistory`.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product 类型, 支持 值 include: `mysql` - 云 数据库 MySQL, `cynosdb` - 云 数据库 CynosDB 对于 MySQL, 默认值 是 `mysql`.",
			},

			"diag_item": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "diagnostic item.",
			},

			"diag_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Diagnostic 类型.",
			},

			"explanation": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Diagnostic 事件 details, output 是 空 如果 there 是 无 additional explanatory 信息.",
			},

			"outline": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Diagnostic summary.",
			},

			"problem": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Diagnosed problem.",
			},

			"severity": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "severity. severity 是 divided into 5 levels, according 到 degree 的 impact 从 high 到 low: 1: Fatal, 2: Serious, 3: Warning, 4: Prompt, 5: Healthy.",
			},

			"start_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Starting 时间.",
			},

			"suggestions": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "A diagnostic suggestion, 或 空 如果 there 是 无 suggestion.",
			},

			"metric": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "reserved text. 注意: 此 字段 可能 返回 null, indicating 该 无 有效 值 可以 是 获取.",
			},

			"end_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "End Time.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainDiagEventRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_diag_event.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var id string

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
		id = v.(string)
	}

	if v, _ := d.GetOk("event_id"); v != nil {
		paramMap["event_id"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["product"] = helper.String(v.(string))
	}

	var result *dbbrain.DescribeDBDiagEventResponseParams
	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		var e error
		result, e = service.DescribeDbbrainDiagEventByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if result != nil {
		if result.DiagItem != nil {
			_ = d.Set("diag_item", result.DiagItem)
		}

		if result.DiagType != nil {
			_ = d.Set("diag_type", result.DiagType)
		}

		if result.EventId != nil {
			_ = d.Set("event_id", result.EventId)
		}

		if result.Explanation != nil {
			_ = d.Set("explanation", result.Explanation)
		}

		if result.Outline != nil {
			_ = d.Set("outline", result.Outline)
		}

		if result.Problem != nil {
			_ = d.Set("problem", result.Problem)
		}

		if result.Severity != nil {
			_ = d.Set("severity", result.Severity)
		}

		if result.StartTime != nil {
			_ = d.Set("start_time", result.StartTime)
		}

		if result.Suggestions != nil {
			_ = d.Set("suggestions", result.Suggestions)
		}

		if result.Metric != nil {
			_ = d.Set("metric", result.Metric)
		}

		if result.EndTime != nil {
			_ = d.Set("end_time", result.EndTime)
		}

	}

	d.SetId(id)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
