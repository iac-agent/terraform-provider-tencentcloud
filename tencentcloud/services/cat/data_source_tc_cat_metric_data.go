package cat

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCatMetricData() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCatMetricDataRead,
		Schema: map[string]*schema.Schema{
			"analyze_task_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Analysis 的 任务 类型，支持 types: `AnalyzeTaskType_Network`: 网络 quality，`AnalyzeTaskType_Browse`: 页面 performance，`AnalyzeTaskType_Transport`: 端口 performance，`AnalyzeTaskType_UploadDownload`: 文件 transport，`AnalyzeTaskType_MediaStream`: audiovisual experience。",
			},

			"metric_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Metric 类型，metrics queries 是 passed 使用 gauge 通过 默认值。",
			},

			"field": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Detailed 字段 的 metrics，指定 metrics 可以 是 passed 或 aggregate metrics，such 作为 avg(ping_time) 表示 entire 延迟",
			},

			"filter": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "过滤器 conditions 可以 是 passed 作为 单个 过滤器 或 多个 参数 concatenated together。",
			},

			"group_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Aggregation 时间，such 作为 1m，1d，30d，和 so 在。",
			},

			"filters": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Multiple condition filtering，支持 combining 多个 filtering conditions 对于 查询。",
			},

			"metric_set": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Return JSON 字符串。",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudCatMetricDataRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cat_metric_data.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("analyze_task_type"); ok {
		paramMap["AnalyzeTaskType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("metric_type"); ok {
		paramMap["MetricType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("field"); ok {
		paramMap["Field"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter"); ok {
		paramMap["Filter"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_by"); ok {
		paramMap["GroupBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.(*schema.Set).List()
		paramMap["Filters"] = helper.InterfacesStringsPoint(filtersSet)
	}

	service := CatService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var metric *cat.DescribeProbeMetricDataResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCatMetricDataByFilter(ctx, paramMap)
		if e != nil {
			if sdkError, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "FailedOperation.DbQueryFailed" {
					return resource.NonRetryableError(e)
				}
			}
			return tccommon.RetryError(e)
		}
		metric = result
		return nil
	})
	if err != nil {
		return err
	}

	var metricSet string
	if metric != nil && metric.MetricSet != nil {
		metricSet = *metric.MetricSet
		_ = d.Set("metric_set", metric.MetricSet)
	}

	d.SetId(helper.DataResourceIdsHash([]string{metricSet}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), metricSet); e != nil {
			return e
		}
	}
	return nil
}
