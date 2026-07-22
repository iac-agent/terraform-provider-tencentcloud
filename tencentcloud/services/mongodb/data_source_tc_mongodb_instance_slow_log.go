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
				Description: "实例 ID，the 格式 is: cmgo-9d0p6umb.Same as the instance ID displayed in the cloud database console page。",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Slow log 开始时间，格式: yyyy-mm-dd hh:mm:ss，such as: 2019-06-01 10:00:00. The time intervalbetween the start and end of the query cannot exceed 24 hours,and only slow logs within the last 7 days are allowed to be queried。",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Slow log termination time，格式: yyyy-mm-dd hh:mm:ss，such as: 2019-06-02 12:00:00.The 时间间隔 between the start and end of the query cannot exceed 24 hours,and only slow logs within the last 7 days are allowed to be queried。",
			},

			"slow_ms": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Slow log 执行时间 threshold，return slow logs whose 执行时间 exceeds this threshold,the unit is milliseconds (ms)，and the minimum is 100 milliseconds。",
			},

			"format": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Slow log return 格式 By default，the original slow log 格式 is returned,and versions 4.4 and above can be set to json。",
			},

			"slow_logs": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "details of slow logs。",
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
