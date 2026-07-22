package dbbrain

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbbrainSlowLogUserHostStats() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainSlowLogUserHostStatsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start 时间 的 查询 范围, 时间 格式 such 作为: 2019-09-10 12:13:14.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "EndTime 时间 的 查询 范围, 时间 格式 such 作为: 2019-09-10 12:13:14.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Types 的 服务 products, 支持 值:`mysql` - Cloud Database MySQL; `cynosdb` - Cloud Database TDSQL-C 对于 MySQL, defaults 到 `mysql`.",
			},

			"md5": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "MD5 值 的 SOL template.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Detailed 列表 的 slow 日志 proportion 对于 each source 地址.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "source 地址.",
						},
						"ratio": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "ratio 的 数量 的 slow logs 的 source 地址 到 总数, 在 %.",
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 的 slow logs 对于 此 source 地址.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainSlowLogUserHostStatsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_slow_log_user_host_stats.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var id string
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
		id = v.(string)
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["start_time"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["end_time"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["product"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("md5"); ok {
		paramMap["md5"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var items []*dbbrain.SlowLogHost

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainSlowLogUserHostStatsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		items = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(items))

	if items != nil {
		for _, slowLogHost := range items {
			slowLogHostMap := map[string]interface{}{}

			if slowLogHost.UserHost != nil {
				slowLogHostMap["user_host"] = slowLogHost.UserHost
			}

			if slowLogHost.Ratio != nil {
				slowLogHostMap["ratio"] = slowLogHost.Ratio
			}

			if slowLogHost.Count != nil {
				slowLogHostMap["count"] = slowLogHost.Count
			}

			tmpList = append(tmpList, slowLogHostMap)
		}

		_ = d.Set("items", tmpList)
	}

	d.SetId(id)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
