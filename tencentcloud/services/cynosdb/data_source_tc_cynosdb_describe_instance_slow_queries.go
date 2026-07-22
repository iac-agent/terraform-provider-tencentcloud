package cynosdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCynosdbDescribeInstanceSlowQueries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbDescribeInstanceSlowQueriesRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "集群 ID。",
			},
			"start_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "开始时间。",
			},
			"end_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "结束时间。",
			},
			"binlogs": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Binlog列表注意：该字段可能返回null，表示无法获取到有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "二进制日志文件名。",
						},
						"file_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "文件大小（以字节为单位）。",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "最早交易时间。",
						},
						"finish_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "最晚交易时间。",
						},
						"binlog_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "二进制日志文件ID。",
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

func dataSourceTencentCloudCynosdbDescribeInstanceSlowQueriesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cynosdb_describe_instance_slow_queries.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service   = CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		binlogs   []*cynosdb.BinlogItem
		clusterId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbDescribeInstanceSlowQueriesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		binlogs = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(binlogs))

	if binlogs != nil {
		for _, binlogItem := range binlogs {
			binlogItemMap := map[string]interface{}{}

			if binlogItem.FileName != nil {
				binlogItemMap["file_name"] = binlogItem.FileName
			}

			if binlogItem.FileSize != nil {
				binlogItemMap["file_size"] = binlogItem.FileSize
			}

			if binlogItem.StartTime != nil {
				binlogItemMap["start_time"] = binlogItem.StartTime
			}

			if binlogItem.FinishTime != nil {
				binlogItemMap["finish_time"] = binlogItem.FinishTime
			}

			if binlogItem.BinlogId != nil {
				binlogItemMap["binlog_id"] = binlogItem.BinlogId
			}

			tmpList = append(tmpList, binlogItemMap)
		}

		_ = d.Set("binlogs", tmpList)
	}

	d.SetId(clusterId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
