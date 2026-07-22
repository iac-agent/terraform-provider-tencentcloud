package css

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCssMonitorReport() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCssMonitorReportRead,
		Schema: map[string]*schema.Schema{
			"monitor_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Monitor ID。",
			},

			"mps_result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "信息 about media processing 结果注意：此字段可能返回 null，表示未找到有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ai_asr_results": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "结果 的 intelligent speech recognition.注意：此字段可能返回 null，表示未找到有效值。",
						},
						"ai_ocr_results": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "结果 的 intelligent text recognition.注意：此字段可能返回 null，表示未找到有效值。",
						},
					},
				},
			},

			"diagnose_result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "信息 about media diagnostic 结果注意：此字段可能返回 null，表示未找到有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"stream_broken_results": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "信息 about 流 interruption.注意：此字段可能返回 null，表示未找到有效值。",
						},
						"low_frame_rate_results": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "信息 about low frame 速率.注意：此字段可能返回 null，表示未找到有效值。",
						},
						"stream_format_results": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "信息 about 流 格式 diagnosis.注意：此字段可能返回 null，表示未找到有效值。",
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

func dataSourceTencentCloudCssMonitorReportRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_css_monitor_report.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var monitorId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("monitor_id"); ok {
		monitorId = v.(string)
		paramMap["MonitorId"] = helper.String(v.(string))
	}

	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var mPSResult *css.DescribeMonitorReportResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCssMonitorReportByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		mPSResult = result
		return nil
	})
	if err != nil {
		return err
	}

	if mPSResult.MPSResult != nil {
		mPSResultMap := map[string]interface{}{}

		if mPSResult.MPSResult.AiAsrResults != nil {
			mPSResultMap["ai_asr_results"] = mPSResult.MPSResult.AiAsrResults
		}

		if mPSResult.MPSResult.AiOcrResults != nil {
			mPSResultMap["ai_ocr_results"] = mPSResult.MPSResult.AiOcrResults
		}

		_ = d.Set("mps_result", []interface{}{mPSResultMap})
	}

	if mPSResult.DiagnoseResult != nil {
		diagnoseResultMap := map[string]interface{}{}

		if mPSResult.DiagnoseResult.StreamBrokenResults != nil {
			diagnoseResultMap["stream_broken_results"] = mPSResult.DiagnoseResult.StreamBrokenResults
		}

		if mPSResult.DiagnoseResult.LowFrameRateResults != nil {
			diagnoseResultMap["low_frame_rate_results"] = mPSResult.DiagnoseResult.LowFrameRateResults
		}

		if mPSResult.DiagnoseResult.StreamFormatResults != nil {
			diagnoseResultMap["stream_format_results"] = mPSResult.DiagnoseResult.StreamFormatResults
		}

		_ = d.Set("diagnose_result", []interface{}{diagnoseResultMap})
	}

	d.SetId(helper.DataResourceIdsHash([]string{monitorId}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
