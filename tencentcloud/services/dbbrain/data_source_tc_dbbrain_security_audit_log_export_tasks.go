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

func DataSourceTencentCloudDbbrainSecurityAuditLogExportTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainSecurityAuditLogExportTasksRead,
		Schema: map[string]*schema.Schema{
			"sec_audit_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "安全 audit 组 ID.",
			},

			"product": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "product, 可选 值 是 mysql.",
			},

			"async_request_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional:    true,
				Description: "async 请求 ID 列表.",
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "安全 audit 日志 export 任务 列表.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"async_request_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "async 请求 ID.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "start 时间.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "end 时间.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "create 时间.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态.",
						},
						"progress": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "任务 progress.",
						},
						"log_start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "日志 start 时间.",
						},
						"log_end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "日志 end 时间.",
						},
						"total_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "总数 大小 的 日志.",
						},
						"danger_levels": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "danger 级别 列表.",
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

func dataSourceTencentCloudDbbrainSecurityAuditLogExportTasksRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_security_audit_log_export_tasks.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var sag_id string

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("sec_audit_group_id"); ok {
		paramMap["sec_audit_group_id"] = helper.String(v.(string))
		sag_id = v.(string)
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["product"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("async_request_ids"); ok {
		async_request_idSet := v.(*schema.Set).List()
		tmpList := make([]*uint64, 0, len(async_request_idSet))
		for i := range async_request_idSet {
			async_request_id := async_request_idSet[i].(int)
			tmpList = append(tmpList, helper.IntUint64(async_request_id))
		}
		paramMap["async_request_ids"] = tmpList
	}

	dbbrainService := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var tasks []*dbbrain.SecLogExportTaskInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := dbbrainService.DescribeDbbrainSecurityAuditLogExportTasksByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		tasks = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Dbbrain tasks failed, reason:%+v", logId, err)
		return err
	}

	ids := make([]string, 0, len(tasks))
	taskList := make([]map[string]interface{}, 0, len(tasks))

	if tasks != nil {

		for _, task := range tasks {
			taskMap := map[string]interface{}{}
			if task.AsyncRequestId != nil {
				taskMap["async_request_id"] = task.AsyncRequestId
			}
			if task.StartTime != nil {
				taskMap["start_time"] = task.StartTime
			}
			if task.EndTime != nil {
				taskMap["end_time"] = task.EndTime
			}
			if task.CreateTime != nil {
				taskMap["create_time"] = task.CreateTime
			}
			if task.Status != nil {
				taskMap["status"] = task.Status
			}
			if task.Progress != nil {
				taskMap["progress"] = task.Progress
			}
			if task.LogStartTime != nil {
				taskMap["log_start_time"] = task.LogStartTime
			}
			if task.LogEndTime != nil {
				taskMap["log_end_time"] = task.LogEndTime
			}
			if task.TotalSize != nil {
				taskMap["total_size"] = task.TotalSize
			}
			if task.DangerLevels != nil {
				taskMap["danger_levels"] = task.DangerLevels
			}
			ids = append(ids, sag_id+tccommon.FILED_SP+helper.UInt64ToStr(*task.AsyncRequestId))
			taskList = append(taskList, taskMap)
		}
		d.SetId(helper.DataResourceIdsHash(ids))
		_ = d.Set("list", taskList)
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), taskList); e != nil {
			return e
		}
	}

	return nil
}
