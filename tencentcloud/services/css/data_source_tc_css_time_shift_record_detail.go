package css

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCssTimeShiftRecordDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCssTimeShiftRecordDetailRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Push 域名",
			},

			"app_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Push 路径",
			},

			"stream_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Stream 名称",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "starting 时间 的 查询 范围 是 指定 在 Unix 时间戳。",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "ending 时间 的 查询 范围 是 指定 在 Unix 时间戳。",
			},

			"domain_group": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "streaming 域名 belongs 到 组. 如果 there 是 无 域名 组 或 域名 组 是 空 字符串，它 可以 是 left blank。",
			},

			"trans_code_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "transcoding 模板 ID 可以 是 left blank 如果 它 是 0。",
			},

			"record_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "数组 时间-shift recording sessions.注意：此字段可能返回 null，表示未找到有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "identifier 对于 时间-shift recording 会话。",
						},
						"start_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "开始时间 的 recording 会话 是 指定 在 Unix 时间戳。",
						},
						"end_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "结束时间 的 recording 会话 是 指定 在 Unix 时间戳。",
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

func dataSourceTencentCloudCssTimeShiftRecordDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_css_time_shift_record_detail.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("domain"); ok {
		paramMap["Domain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("app_name"); ok {
		paramMap["AppName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("stream_name"); ok {
		paramMap["StreamName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("start_time"); ok {
		paramMap["StartTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("end_time"); ok {
		paramMap["EndTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("domain_group"); ok {
		paramMap["DomainGroup"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("trans_code_id"); ok {
		paramMap["TransCodeId"] = helper.IntUint64(v.(int))
	}

	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var recordList []*css.TimeShiftRecord
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCssTimeShiftRecordDetailByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		recordList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(recordList))
	tmpList := make([]map[string]interface{}, 0, len(recordList))

	if recordList != nil {
		for _, timeShiftRecord := range recordList {
			timeShiftRecordMap := map[string]interface{}{}

			if timeShiftRecord.Sid != nil {
				timeShiftRecordMap["sid"] = timeShiftRecord.Sid
			}

			if timeShiftRecord.StartTime != nil {
				timeShiftRecordMap["start_time"] = timeShiftRecord.StartTime
			}

			if timeShiftRecord.EndTime != nil {
				timeShiftRecordMap["end_time"] = timeShiftRecord.EndTime
			}

			ids = append(ids, *timeShiftRecord.Sid)
			tmpList = append(tmpList, timeShiftRecordMap)
		}

		_ = d.Set("record_list", tmpList)
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
