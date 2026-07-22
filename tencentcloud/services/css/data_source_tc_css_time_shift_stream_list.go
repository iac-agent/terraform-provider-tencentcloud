package css

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCssTimeShiftStreamList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCssTimeShiftStreamListRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The 开始时间，which must be a Unix 时间戳。",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The 结束时间，which must be a Unix 时间戳。",
			},

			"stream_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The stream 名称",
			},

			"domain": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The push 域名",
			},

			"domain_group": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The group the push 域名 belongs to。",
			},

			"total_size": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total 数量 records in the specified time 周期",
			},

			"stream_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The information of the streams.注意：此字段可能返回 null，表示无法获取有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The group the push 域名 belongs to.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The push 域名",
						},
						"app_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The push 路径",
						},
						"stream_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The stream 名称",
						},
						"start_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The stream 开始时间，which is a Unix 时间戳。",
						},
						"end_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The stream 结束时间 (for streams that ended before the time of query)，which is a Unix 时间戳。",
						},
						"trans_code_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The transcoding 模板 ID注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"stream_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The stream 类型 `0`: The original stream; `1`: The watermarked stream; `2`: The transcoded stream。",
						},
						"duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The storage duration (seconds) of the recording.注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudCssTimeShiftStreamListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_css_time_shift_stream_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOkExists("start_time"); ok {
		paramMap["StartTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("end_time"); ok {
		paramMap["EndTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("stream_name"); ok {
		paramMap["StreamName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain"); ok {
		paramMap["Domain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain_group"); ok {
		paramMap["DomainGroup"] = helper.String(v.(string))
	}

	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var timeShiftStreamList []*css.TimeShiftStreamInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCssTimeShiftStreamListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		timeShiftStreamList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(timeShiftStreamList))
	tmpList := make([]map[string]interface{}, 0, len(timeShiftStreamList))
	if timeShiftStreamList != nil {
		for _, timeShiftStreamInfo := range timeShiftStreamList {
			timeShiftStreamInfoMap := map[string]interface{}{}

			if timeShiftStreamInfo.DomainGroup != nil {
				timeShiftStreamInfoMap["domain_group"] = timeShiftStreamInfo.DomainGroup
			}

			if timeShiftStreamInfo.Domain != nil {
				timeShiftStreamInfoMap["domain"] = timeShiftStreamInfo.Domain
			}

			if timeShiftStreamInfo.AppName != nil {
				timeShiftStreamInfoMap["app_name"] = timeShiftStreamInfo.AppName
			}

			if timeShiftStreamInfo.StreamName != nil {
				timeShiftStreamInfoMap["stream_name"] = timeShiftStreamInfo.StreamName
			}

			if timeShiftStreamInfo.StartTime != nil {
				timeShiftStreamInfoMap["start_time"] = timeShiftStreamInfo.StartTime
			}

			if timeShiftStreamInfo.EndTime != nil {
				timeShiftStreamInfoMap["end_time"] = timeShiftStreamInfo.EndTime
			}

			if timeShiftStreamInfo.TransCodeId != nil {
				timeShiftStreamInfoMap["trans_code_id"] = timeShiftStreamInfo.TransCodeId
			}

			if timeShiftStreamInfo.StreamType != nil {
				timeShiftStreamInfoMap["stream_type"] = timeShiftStreamInfo.StreamType
			}

			if timeShiftStreamInfo.Duration != nil {
				timeShiftStreamInfoMap["duration"] = timeShiftStreamInfo.Duration
			}

			ids = append(ids, *timeShiftStreamInfo.Domain)
			tmpList = append(tmpList, timeShiftStreamInfoMap)
		}

		_ = d.Set("stream_list", tmpList)
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
