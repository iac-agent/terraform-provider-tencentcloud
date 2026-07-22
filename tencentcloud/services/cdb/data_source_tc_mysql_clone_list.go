package cdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMysqlCloneList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlCloneListRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "查询指定源实例的克隆任务列表。",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "克隆任务列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"src_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "克隆任务的源实例Id。",
						},
						"dst_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "克隆任务新生成的实例 ID。",
						},
						"clone_job_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "克隆任务对应的任务列表id。",
						},
						"rollback_strategy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "克隆实例使用的策略包括以下几种：timepoint：指定时间点回滚，backupset：指定备份文件回滚。",
						},
						"rollback_target_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "克隆实例回滚的时间点。",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "任务开始时间。",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "任务结束时间。",
						},
						"task_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "任务状态，包括以下状态：initial、running、wait_complete、success、failed。",
						},
						"new_region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "克隆实例所在地域的Id。",
						},
						"src_region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "源实例所在地域的Id。",
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

func dataSourceTencentCloudMysqlCloneListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_clone_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := ""
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var cloneList []*cdb.CloneItem
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlCloneListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		cloneList = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(cloneList))
	if cloneList != nil {
		for _, cloneItem := range cloneList {
			cloneItemMap := map[string]interface{}{}

			if cloneItem.SrcInstanceId != nil {
				cloneItemMap["src_instance_id"] = cloneItem.SrcInstanceId
			}

			if cloneItem.DstInstanceId != nil {
				cloneItemMap["dst_instance_id"] = cloneItem.DstInstanceId
			}

			if cloneItem.CloneJobId != nil {
				cloneItemMap["clone_job_id"] = cloneItem.CloneJobId
			}

			if cloneItem.RollbackStrategy != nil {
				cloneItemMap["rollback_strategy"] = cloneItem.RollbackStrategy
			}

			if cloneItem.RollbackTargetTime != nil {
				cloneItemMap["rollback_target_time"] = cloneItem.RollbackTargetTime
			}

			if cloneItem.StartTime != nil {
				cloneItemMap["start_time"] = cloneItem.StartTime
			}

			if cloneItem.EndTime != nil {
				cloneItemMap["end_time"] = cloneItem.EndTime
			}

			if cloneItem.TaskStatus != nil {
				cloneItemMap["task_status"] = cloneItem.TaskStatus
			}

			// if cloneItem.NewRegionId != nil {
			// 	cloneItemMap["new_region_id"] = cloneItem.NewRegionId
			// }

			// if cloneItem.SrcRegionId != nil {
			// 	cloneItemMap["src_region_id"] = cloneItem.SrcRegionId
			// }

			tmpList = append(tmpList, cloneItemMap)
		}

		err = d.Set("items", tmpList)
		if err != nil {
			return err
		}
	}

	d.SetId(helper.DataResourceIdsHash([]string{instanceId}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
