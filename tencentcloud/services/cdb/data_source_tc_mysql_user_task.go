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

func DataSourceTencentCloudMysqlUserTask() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlUserTaskRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "实例ID，格式为：cdb-c1nl9rpv，与云数据库控制台页面显示的实例ID相同，可以通过【查询实例列表】（https://cloud.tencent.com/document/api/236/15872）接口获取输出参数中InstanceId字段的值。",
			},

			"async_request_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "异步任务请求ID，执行云数据库相关操作返回的AsyncRequestId。",
			},

			"task_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "任务类型。如果不传值，则查询所有任务类型。支持的值包括： `ROLLBACK` - 数据库回滚； `SQL OPERATION` - SQL 操作； `IMPORT DATA`-数据导入； `MODIFY PARAM`-参数设置； `INITIAL` - 初始化云数据库实例； `REBOOT` - 重新启动云数据库实例； `OPEN GTID` - 打开云数据库实例GTID； `UPGRADE RO` - 只读实例升级； `BATCH ROLLBACK` - 数据库批量回滚； `UPGRADE MASTER`-主升级； `DROP TABLES` - 删除云数据库表； `SWITCH DR TO MASTER` - 灾难恢复实例。",
			},

			"task_status": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "任务状态。如果不传值，则查询所有任务状态。支持的值包括： `UNDEFINED` - 未定义； `INITIAL` - 初始化； `RUNNING`-正在运行； `SUCCEED` - 执行成功； `FAILED` - 执行失败； `KILLED` - 终止； `已删除`-已删除； “暂停”- 暂停。",
			},

			"start_time_begin": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "第一个任务的开始时间，用于范围查询，时间格式如下：2017-12-31 10:40:01。",
			},

			"start_time_end": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "最后一个任务的开始时间，用于范围查询，时间格式如下：2017-12-31 10:40:01。",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "返回的实例任务信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "错误代码。",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "错误信息。",
						},
						"job_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例任务ID。",
						},
						"progress": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例任务进度。",
						},
						"task_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例任务状态，可能的值包括：UNDEFINED - 未定义；INITIAL - 初始化；RUNNING - 正在运行；SUCCEED - 执行成功；FAILED - 执行失败；KILLED - 终止；REMOVED - 已删除；PAUSED - 已暂停。WAITING - 等待（可取消）。",
						},
						"task_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例任务类型，可能的值包括：ROLLBACK - 数据库回滚；SQL OPERATION - SQL操作；IMPORT DATA - 数据导入；MODIFY PARAM - 参数设置；INITIAL - 初始化云数据库实例；REBOOT - 重启云数据库实例；OPEN GTID - 打开云数据库实例GTID；UPGRADE RO - 只读实例升级；BATCH ROLLBACK - 数据库批量回滚；UPGRADE MASTER - 主升级；DROP TABLES - 删除云数据库表；SWITCH DR TO MASTER - 灾难恢复实例。",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例任务开始时间。",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例任务结束时间。",
						},
						"instance_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "与任务关联的实例 ID。注意：该字段可能返回null，表示取不到有效值。",
						},
						"async_request_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "异步任务的请求ID。",
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

func dataSourceTencentCloudMysqlUserTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_user_task.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("async_request_id"); ok {
		paramMap["AsyncRequestId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_types"); ok {
		taskTypesSet := v.(*schema.Set).List()
		paramMap["TaskTypes"] = helper.InterfacesStringsPoint(taskTypesSet)
	}

	if v, ok := d.GetOk("task_status"); ok {
		taskStatusSet := v.(*schema.Set).List()
		paramMap["TaskStatus"] = helper.InterfacesStringsPoint(taskStatusSet)

	}

	if v, ok := d.GetOk("start_time_begin"); ok {
		paramMap["StartTimeBegin"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time_end"); ok {
		paramMap["StartTimeEnd"] = helper.String(v.(string))
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var items []*cdb.TaskDetail
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlUserTaskByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		items = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(items))
	tmpList := make([]map[string]interface{}, 0, len(items))
	if items != nil {
		for _, taskDetail := range items {
			taskDetailMap := map[string]interface{}{}

			if taskDetail.Code != nil {
				taskDetailMap["code"] = taskDetail.Code
			}

			if taskDetail.Message != nil {
				taskDetailMap["message"] = taskDetail.Message
			}

			if taskDetail.JobId != nil {
				taskDetailMap["job_id"] = taskDetail.JobId
			}

			if taskDetail.Progress != nil {
				taskDetailMap["progress"] = taskDetail.Progress
			}

			if taskDetail.TaskStatus != nil {
				taskDetailMap["task_status"] = taskDetail.TaskStatus
			}

			if taskDetail.TaskType != nil {
				taskDetailMap["task_type"] = taskDetail.TaskType
			}

			if taskDetail.StartTime != nil {
				taskDetailMap["start_time"] = taskDetail.StartTime
			}

			if taskDetail.EndTime != nil {
				taskDetailMap["end_time"] = taskDetail.EndTime
			}

			if taskDetail.InstanceIds != nil {
				taskDetailMap["instance_ids"] = taskDetail.InstanceIds
			}

			if taskDetail.AsyncRequestId != nil {
				taskDetailMap["async_request_id"] = taskDetail.AsyncRequestId
			}

			ids = append(ids, strconv.FormatInt(*taskDetail.JobId, 10))
			tmpList = append(tmpList, taskDetailMap)
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
