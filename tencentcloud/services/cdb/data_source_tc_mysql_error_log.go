package cdb

import (
	"context"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMysqlErrorLog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlErrorLogRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID。",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "开始时间戳。例如 1585142640。",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "结束时间戳。例如 1585142640。",
			},

			"key_words": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "要匹配的关键字列表，最多支持 15 个关键字。",
			},

			"inst_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "仅当实例为主实例或灾备实例时有效，可选值：slave，表示拉取从机的日志。",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "记录返回。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "错误发生的时间。",
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "错误详细信息。",
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

func dataSourceTencentCloudMysqlErrorLogRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_error_log.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("start_time"); v != nil {
		paramMap["StartTime"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("end_time"); v != nil {
		paramMap["EndTime"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("key_words"); ok {
		keyWordsSet := v.(*schema.Set).List()
		paramMap["KeyWords"] = helper.InterfacesStringsPoint(keyWordsSet)
	}

	if v, ok := d.GetOk("inst_type"); ok {
		paramMap["InstType"] = helper.String(v.(string))
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var result []*cdb.ErrlogItem
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeMysqlErrorLogByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	tmpList := make([]map[string]interface{}, 0, len(result))
	if result != nil {
		for _, errlogItem := range result {
			errlogItemMap := map[string]interface{}{}

			if errlogItem.Timestamp != nil {
				errlogItemMap["timestamp"] = errlogItem.Timestamp
			}

			if errlogItem.Content != nil {
				errlogItemMap["content"] = errlogItem.Content
			}

			ids = append(ids, strconv.FormatUint(*errlogItem.Timestamp, 10))
			tmpList = append(tmpList, errlogItemMap)
		}

		_ = d.Set("items", tmpList)
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
