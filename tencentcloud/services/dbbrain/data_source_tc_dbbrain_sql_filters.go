package dbbrain

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbbrainSqlFilters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainSqlFiltersRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "实例 ID.",
			},

			"filter_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional:    true,
				Description: "过滤器 ID 列表.",
			},

			"statuses": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "状态 列表.",
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "sql 过滤器 列表.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "任务 ID.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "任务 状态, 可选 值 是 RUNNING, FINISHED, TERMINATED.",
						},
						"sql_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "sql 类型, 可选 值 是 SELECT, UPDATE, DELETE, INSERT, REPLACE.",
						},
						"origin_keys": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "源站 keys.",
						},
						"origin_rule": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "源站 规则.",
						},
						"rejected_sql_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "rejected sql count.",
						},
						"current_concurrency": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "当前 concurrency.",
						},
						"max_concurrency": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "maxmum concurrency.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "create 时间.",
						},
						"current_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "当前 时间.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "expire 时间.",
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

func dataSourceTencentCloudDbbrainSqlFiltersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_sql_filters.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter_ids"); ok {
		filter_idSet := v.(*schema.Set).List()
		tmpList := make([]*int64, 0, len(filter_idSet))
		for i := range filter_idSet {
			filter_id := filter_idSet[i].(int)
			tmpList = append(tmpList, helper.IntInt64(filter_id))
		}
		paramMap["filter_ids"] = tmpList
	}

	if v, ok := d.GetOk("statuses"); ok {
		statuseSet := v.(*schema.Set).List()
		tmpList := make([]*string, 0, len(statuseSet))
		for i := range statuseSet {
			status := statuseSet[i].(string)
			tmpList = append(tmpList, helper.String(status))
		}
		paramMap["statuses"] = tmpList
	}

	dbbrainService := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var items []*dbbrain.SQLFilter
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := dbbrainService.DescribeDbbrainSqlFiltersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		items = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Dbbrain items failed, reason:%+v", logId, err)
		return err
	}

	ids := make([]string, 0, len(items))
	itemList := make([]map[string]interface{}, 0, len(items))

	if items != nil {
		for _, item := range items {
			itemMap := map[string]interface{}{}
			if item.Id != nil {
				itemMap["id"] = item.Id
			}
			if item.Status != nil {
				itemMap["status"] = item.Status
			}
			if item.SqlType != nil {
				itemMap["sql_type"] = item.SqlType
			}
			if item.OriginKeys != nil {
				itemMap["origin_keys"] = item.OriginKeys
			}
			if item.OriginRule != nil {
				itemMap["origin_rule"] = item.OriginRule
			}
			if item.RejectedSqlCount != nil {
				itemMap["rejected_sql_count"] = item.RejectedSqlCount
			}
			if item.CurrentConcurrency != nil {
				itemMap["current_concurrency"] = item.CurrentConcurrency
			}
			if item.MaxConcurrency != nil {
				itemMap["max_concurrency"] = item.MaxConcurrency
			}
			if item.CreateTime != nil {
				itemMap["create_time"] = item.CreateTime
			}
			if item.CurrentTime != nil {
				itemMap["current_time"] = item.CurrentTime
			}
			if item.ExpireTime != nil {
				itemMap["expire_time"] = item.ExpireTime
			}
			ids = append(ids, helper.Int64ToStr(*item.Id))
			itemList = append(itemList, itemMap)
		}
		d.SetId(helper.DataResourceIdsHash(ids))
		_ = d.Set("list", itemList)
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), itemList); e != nil {
			return e
		}
	}

	return nil
}
