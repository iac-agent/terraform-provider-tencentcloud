package sqlserver

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSqlserverDescHaLog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverDescHaLogRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID.",
			},
			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start 时间 (yyyy-MM-dd HH:mm:ss).",
			},
			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End 时间 (yyyy-MM-dd HH:mm:ss).",
			},
			"switch_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Switching 模式 0-系统 automatically switches, 1-manual switch, 如果 不 filled 在, all 将 是 checked 通过 默认值.",
			},
			"switch_log": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Master/Slave switching 日志.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Switch 事件 ID 注意: 此 字段 可能 返回 null, indicating 该 无 有效 值 可以 是 获取.",
						},
						"switch_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Switching 模式 0-系统 automatic switching, 1-manual switching 注意: 此 字段 可能 返回 null, indicating 该 无 有效 值 可以 是 获取.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Switch start 时间 注意: 此 字段 可能 返回 null, indicating 该 无 有效 值 可以 是 获取.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Switch end 时间 注意: 此 字段 可能 返回 null, indicating 该 无 有效 值 可以 是 获取.",
						},
						"reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Machine failure causes automatic switching 注意: 此 字段 可能 返回 null, indicating 该 无 有效 值 可以 是 获取.",
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

func dataSourceTencentCloudSqlserverDescHaLogRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_sqlserver_desc_ha_log.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		switchLogs []*sqlserver.SwitchLog
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("switch_type"); ok {
		paramMap["SwitchType"] = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverDescHaLogByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		switchLogs = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(switchLogs))

	if switchLogs != nil {
		for _, switchLog := range switchLogs {
			switchLogMap := map[string]interface{}{}

			if switchLog.EventId != nil {
				switchLogMap["event_id"] = switchLog.EventId
			}

			if switchLog.SwitchType != nil {
				switchLogMap["switch_type"] = switchLog.SwitchType
			}

			if switchLog.StartTime != nil {
				switchLogMap["start_time"] = switchLog.StartTime
			}

			if switchLog.EndTime != nil {
				switchLogMap["end_time"] = switchLog.EndTime
			}

			if switchLog.Reason != nil {
				switchLogMap["reason"] = switchLog.Reason
			}

			tmpList = append(tmpList, switchLogMap)
		}

		_ = d.Set("switch_log", tmpList)
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
