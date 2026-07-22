package waf

import (
	"context"
	"strconv"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWafAttackLogList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafAttackLogListRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "域名 对于 查询，all 域名 使用 all。",
			},
			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "开始时间。",
			},
			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "结束时间。",
			},
			"query_count": {
				Optional:    true,
				Type:        schema.TypeInt,
				Default:     10,
				Description: "数量 queries，默认为 10，最大 的 100。",
			},
			"page": {
				Optional:    true,
				Type:        schema.TypeInt,
				Default:     0,
				Description: "数量 pages，starting 从 0 通过 默认值。",
			},
			"query_string": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Lucene grammar。",
			},
			"sort": {
				Optional:    true,
				Type:        schema.TypeString,
				Default:     "desc",
				Description: "Default desc，support desc，asc。",
			},
			"data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Attack 日志 数组。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "detail 的 attack 日志。",
						},
						"file_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Useless。",
						},
						"source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Useless。",
						},
						"time_stamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time 字符串。",
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

func dataSourceTencentCloudWafAttackLogListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_waf_attack_log_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service       = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		attackLogList []*waf.AttackLogInfo
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("domain"); ok {
		paramMap["Domain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("query_count"); ok {
		paramMap["Count"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("query_string"); ok {
		paramMap["QueryString"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort"); ok {
		paramMap["Sort"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("page"); ok {
		paramMap["Page"] = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafAttackLogListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		attackLogList = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(attackLogList))

	if attackLogList != nil {
		for _, attackLogInfo := range attackLogList {
			attackLogInfoMap := map[string]interface{}{}

			if attackLogInfo.Content != nil {
				attackLogInfoMap["content"] = attackLogInfo.Content
			}

			if attackLogInfo.FileName != nil {
				attackLogInfoMap["file_name"] = attackLogInfo.FileName
			}

			if attackLogInfo.Source != nil {
				attackLogInfoMap["source"] = attackLogInfo.Source
			}

			if attackLogInfo.TimeStamp != nil {
				attackLogInfoMap["time_stamp"] = attackLogInfo.TimeStamp
			}

			tmpList = append(tmpList, attackLogInfoMap)
		}

		_ = d.Set("data", tmpList)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
