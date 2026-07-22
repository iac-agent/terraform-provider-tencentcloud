package cdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMysqlSlowLog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlSlowLogRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例ID，格式为：cdb-c1nl9rpv。与云数据库控制台页面显示的实例ID相同。",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "满足查询条件的慢查询日志详情。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "备份文件名。",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "备份文件大小，单位：Byte。",
						},
						"date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "备份快照时间，时间格式：2016-03-17 02:10:37。",
						},
						"intranet_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "内网下载地址。",
						},
						"internet_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "外网下载地址。",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "日志具体类型，可能的值：slowlog - 慢速日志。",
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

func dataSourceTencentCloudMysqlSlowLogRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_slow_log.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var instanceId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var items []*cdb.SlowLogInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlSlowLogByFilter(ctx, paramMap)
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
		for _, slowLogInfo := range items {
			slowLogInfoMap := map[string]interface{}{}

			if slowLogInfo.Name != nil {
				slowLogInfoMap["name"] = slowLogInfo.Name
			}

			if slowLogInfo.Size != nil {
				slowLogInfoMap["size"] = slowLogInfo.Size
			}

			if slowLogInfo.Date != nil {
				slowLogInfoMap["date"] = slowLogInfo.Date
			}

			if slowLogInfo.IntranetUrl != nil {
				slowLogInfoMap["intranet_url"] = slowLogInfo.IntranetUrl
			}

			if slowLogInfo.InternetUrl != nil {
				slowLogInfoMap["internet_url"] = slowLogInfo.InternetUrl
			}

			if slowLogInfo.Type != nil {
				slowLogInfoMap["type"] = slowLogInfo.Type
			}

			tmpList = append(tmpList, slowLogInfoMap)
		}

		_ = d.Set("items", tmpList)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
