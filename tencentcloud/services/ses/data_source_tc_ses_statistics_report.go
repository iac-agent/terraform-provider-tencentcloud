package ses

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSesStatisticsReport() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSesStatisticsReportRead,
		Schema: map[string]*schema.Schema{
			"start_date": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start date。",
			},

			"end_date": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End date。",
			},

			"domain": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sender 域名",
			},

			"receiving_mailbox_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Recipient 地址 类型，对于 示例，gmail.com。",
			},

			"daily_volumes": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Daily email sending 统计。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"send_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"request_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 email requests。",
						},
						"accepted_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 email requests accepted 通过 Tencent Cloud。",
						},
						"delivered_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 delivered emails。",
						},
						"opened_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 users (deduplicated) who opened emails。",
						},
						"clicked_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 recipients who clicked 在 links 在 emails。",
						},
						"bounce_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 bounced emails。",
						},
						"unsubscribe_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 users who canceled subscriptions. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
					},
				},
			},

			"overall_volume": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Overall email sending 统计。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"send_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"request_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 email requests。",
						},
						"accepted_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 email requests accepted 通过 Tencent Cloud。",
						},
						"delivered_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 delivered emails。",
						},
						"opened_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 users (deduplicated) who opened emails。",
						},
						"clicked_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 recipients who clicked 在 links 在 emails。",
						},
						"bounce_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 bounced emails。",
						},
						"unsubscribe_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 users who canceled subscriptions. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
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

func dataSourceTencentCloudSesStatisticsReportRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ses_statistics_report.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("start_date"); ok {
		paramMap["StartDate"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_date"); ok {
		paramMap["EndDate"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain"); ok {
		paramMap["Domain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("receiving_mailbox_type"); ok {
		paramMap["ReceivingMailboxType"] = helper.String(v.(string))
	}

	service := SesService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var statisticsReport *ses.GetStatisticsReportResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSesStatisticsReportByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		statisticsReport = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(statisticsReport.DailyVolumes))
	tmpList := make([]map[string]interface{}, 0, len(statisticsReport.DailyVolumes))

	if statisticsReport.DailyVolumes != nil {
		for _, volume := range statisticsReport.DailyVolumes {
			volumeMap := map[string]interface{}{}

			if volume.SendDate != nil {
				volumeMap["send_date"] = volume.SendDate
			}

			if volume.RequestCount != nil {
				volumeMap["request_count"] = volume.RequestCount
			}

			if volume.AcceptedCount != nil {
				volumeMap["accepted_count"] = volume.AcceptedCount
			}

			if volume.DeliveredCount != nil {
				volumeMap["delivered_count"] = volume.DeliveredCount
			}

			if volume.OpenedCount != nil {
				volumeMap["opened_count"] = volume.OpenedCount
			}

			if volume.ClickedCount != nil {
				volumeMap["clicked_count"] = volume.ClickedCount
			}

			if volume.BounceCount != nil {
				volumeMap["bounce_count"] = volume.BounceCount
			}

			if volume.UnsubscribeCount != nil {
				volumeMap["unsubscribe_count"] = volume.UnsubscribeCount
			}

			ids = append(ids, *volume.SendDate)
			tmpList = append(tmpList, volumeMap)
		}

		_ = d.Set("daily_volumes", tmpList)
	}

	if statisticsReport.OverallVolume != nil {
		overallVolume := statisticsReport.OverallVolume
		volumeMap := map[string]interface{}{}

		if overallVolume.SendDate != nil {
			volumeMap["send_date"] = overallVolume.SendDate
		}

		if overallVolume.RequestCount != nil {
			volumeMap["request_count"] = overallVolume.RequestCount
		}

		if overallVolume.AcceptedCount != nil {
			volumeMap["accepted_count"] = overallVolume.AcceptedCount
		}

		if overallVolume.DeliveredCount != nil {
			volumeMap["delivered_count"] = overallVolume.DeliveredCount
		}

		if overallVolume.OpenedCount != nil {
			volumeMap["opened_count"] = overallVolume.OpenedCount
		}

		if overallVolume.ClickedCount != nil {
			volumeMap["clicked_count"] = overallVolume.ClickedCount
		}

		if overallVolume.BounceCount != nil {
			volumeMap["bounce_count"] = overallVolume.BounceCount
		}

		if overallVolume.UnsubscribeCount != nil {
			volumeMap["unsubscribe_count"] = overallVolume.UnsubscribeCount
		}

		_ = d.Set("overall_volume", []interface{}{volumeMap})
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
