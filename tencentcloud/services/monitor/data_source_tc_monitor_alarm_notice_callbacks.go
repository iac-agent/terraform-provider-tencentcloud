package monitor

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMonitorAlarmNoticeCallbacks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMonitorAlarmNoticeCallbacksRead,
		Schema: map[string]*schema.Schema{
			"url_notices": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Alarm callback 通知。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Callback URL (limited 到 256 字符)。",
						},
						"is_valid": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Verified 0=No 1=Yes。",
						},
						"validation_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Verification 代码",
						},
						"start_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 秒 starting 从 day 的 通知 开始时间。",
						},
						"end_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 秒 从 end 的 通知 day。",
						},
						"weekday": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "Notification 周期 1-7 表示 Monday 到 Sunday。",
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

func dataSourceTencentCloudMonitorAlarmNoticeCallbacksRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_alarm_notice_callbacks.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var urlNotices []*monitor.URLNotice
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMonitorAlarmNoticeCallbacksByFilter(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}
		urlNotices = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(urlNotices))
	tmpList := make([]map[string]interface{}, 0, len(urlNotices))

	if urlNotices != nil {
		for _, urlNotice := range urlNotices {
			urlNoticeMap := map[string]interface{}{}

			if urlNotice.URL != nil {
				urlNoticeMap["url"] = urlNotice.URL
			}

			if urlNotice.IsValid != nil {
				urlNoticeMap["is_valid"] = urlNotice.IsValid
			}

			if urlNotice.ValidationCode != nil {
				urlNoticeMap["validation_code"] = urlNotice.ValidationCode
			}

			if urlNotice.StartTime != nil {
				urlNoticeMap["start_time"] = urlNotice.StartTime
			}

			if urlNotice.EndTime != nil {
				urlNoticeMap["end_time"] = urlNotice.EndTime
			}

			if urlNotice.Weekday != nil {
				urlNoticeMap["weekday"] = urlNotice.Weekday
			}

			ids = append(ids, *urlNotice.URL)
			tmpList = append(tmpList, urlNoticeMap)
		}

		_ = d.Set("url_notices", tmpList)
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
