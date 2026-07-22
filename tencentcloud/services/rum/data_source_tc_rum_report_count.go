package rum

import (
	"context"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudRumReportCount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRumReportCountRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Start 时间 但 是 represented 使用 timestamp 在 秒.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "End 时间 但 是 represented 使用 timestamp 在 秒.",
			},

			"project_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},

			"report_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Report 类型, 空 是 meaning all 类型 count. `日志`:日志 报告 count, `pv`:pv 报告 count, `事件`:事件 报告 count, `speed`:speed 报告 count, `performance`:performance 报告 count, `自定义`:自定义 报告 count, `webvitals`:webvitals 报告 count, `miniProgramData`:miniProgramData 报告 count.",
			},

			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "实例 ID.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Return 值.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 save results.",
			},
		},
	}
}

func dataSourceTencentCloudRumReportCountRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_rum_report_count.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		startTime int
		endTime   int
	)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("start_time"); v != nil {
		startTime = v.(int)
		paramMap["StartTime"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("end_time"); v != nil {
		endTime = v.(int)
		paramMap["EndTime"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("project_id"); v != nil {
		paramMap["ID"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("report_type"); ok {
		paramMap["ReportType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceID"] = helper.String(v.(string))
	}

	service := RumService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var result *string
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeRumReportCountByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	var ids string
	if result != nil {
		ids = *result
		_ = d.Set("result", result)
	}

	d.SetId(helper.DataResourceIdsHash([]string{strconv.Itoa(startTime), strconv.Itoa(endTime), ids}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
