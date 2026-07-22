package mongodb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMongodbInstanceSlowLog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMongodbInstanceSlowLogRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID， 格式 是: cmgo-9d0p6umb.Same 作为 实例 ID displayed 在 云 数据库 console 页面。",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Slow 日志 开始时间，格式: yyyy-mm-dd hh:mm:ss，such 作为: 2019-06-01 10:00:00. 时间 intervalbetween start 和 end 的 查询 不能 exceed 24 hours,和 仅 slow logs within last 7 days 是 allowed 到 是 queried。",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Slow 日志 termination 时间，格式: yyyy-mm-dd hh:mm:ss，such 作为: 2019-06-02 12:00:00. 时间间隔 between start 和 end 的 查询 不能 exceed 24 hours,和 仅 slow logs within last 7 days 是 allowed 到 是 queried。",
			},

			"slow_ms": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Slow 日志 执行时间 阈值，返回 slow logs whose 执行时间 exceeds 此 阈值, 单位 是 milliseconds (ms)，和 最小 是 100 milliseconds。",
			},

			"format": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Slow 日志 返回 格式 By 默认值， original slow 日志 格式 是 返回,和 versions 4.4 和 above 可以 是 集合 到 json。",
			},

			"slow_logs": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "details 的 slow logs。",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudMongodbInstanceSlowLogRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mongodb_instance_slow_log.read")()
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

	if v, _ := d.GetOk("slow_ms"); v != nil {
		paramMap["slow_ms"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("format"); ok {
		paramMap["format"] = helper.String(v.(string))
	}

	service := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var slowLogs []*string

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMongodbInstanceSlowLogByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		slowLogs = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(slowLogs))
	if slowLogs != nil {
		_ = d.Set("slow_logs", slowLogs)
	}

	for _, slowLog := range slowLogs {
		ids = append(ids, *slowLog)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), slowLogs); e != nil {
			return e
		}
	}
	return nil
}
