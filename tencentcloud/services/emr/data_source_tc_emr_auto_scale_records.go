package emr

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudEmrAutoScaleRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEmrAutoScaleRecordsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "EMR 集群 ID.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Record filtering 参数, currently 仅 `StartTime`, `EndTime` 和 `StrategyName` 是 支持. `StartTime` 和 `EndTime` support 时间 格式 的 2006-01-02 15:04:05 或 2006/01/02 15:04:05.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Key. 注意: 此 字段 可能 返回 null, indicating 该 无 有效 值 可以 是 获取.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value. 注意: 此 字段 可能 返回 null, indicating 该 无 有效 值 可以 是 获取.",
						},
					},
				},
			},

			"record_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Record 列表.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"strategy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule 名称 的 expanding 和 shrinking 容量.",
						},
						"scale_action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "`SCALE_OUT` 和 `SCALE_IN` respectively represent expanding 和 shrinking 容量.",
						},
						"action_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "`SUCCESS`, `FAILED`, `PART_SUCCESS`, `IN_PROCESS`.",
						},
						"action_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Process Trigger Time.",
						},
						"scale_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scalability-related Description.",
						},
						"expect_scale_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Effective 仅 当 ScaleAction 是 SCALE_OUT.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Process End Time.",
						},
						"strategy_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Strategy Type, 1 对于 Load scaling, 2 对于 Time scaling.",
						},
						"spec_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specification 信息 使用 当 expanding 容量.",
						},
						"compensate_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Compensation 和 expansion, 0 表示 无 start, 1 表示 start. 注意: 此 字段 可能 返回 null, indicating 该 无 有效 值 可以 是 获取.",
						},
						"compensate_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Compensation Times 注意: 此 字段 可能 返回 null, indicating 该 无 有效 值 可以 是 获取.",
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

func dataSourceTencentCloudEmrAutoScaleRecordsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_emr_auto_scale_records.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var instanceId string

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*emr.KeyValue, 0, len(filtersSet))

		for _, item := range filtersSet {
			keyValue := emr.KeyValue{}
			keyValueMap := item.(map[string]interface{})

			if v, ok := keyValueMap["key"]; ok {
				keyValue.Key = helper.String(v.(string))
			}
			if v, ok := keyValueMap["value"]; ok {
				keyValue.Value = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &keyValue)
		}
		paramMap["Filters"] = tmpSet
	}

	service := EMRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var recordList []*emr.AutoScaleRecord
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeEmrAutoScaleRecordsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		recordList = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(recordList))

	if recordList != nil {
		for _, autoScaleRecord := range recordList {
			autoScaleRecordMap := map[string]interface{}{}

			if autoScaleRecord.StrategyName != nil {
				autoScaleRecordMap["strategy_name"] = autoScaleRecord.StrategyName
			}

			if autoScaleRecord.ScaleAction != nil {
				autoScaleRecordMap["scale_action"] = autoScaleRecord.ScaleAction
			}

			if autoScaleRecord.ActionStatus != nil {
				autoScaleRecordMap["action_status"] = autoScaleRecord.ActionStatus
			}

			if autoScaleRecord.ActionTime != nil {
				autoScaleRecordMap["action_time"] = autoScaleRecord.ActionTime
			}

			if autoScaleRecord.ScaleInfo != nil {
				autoScaleRecordMap["scale_info"] = autoScaleRecord.ScaleInfo
			}

			if autoScaleRecord.ExpectScaleNum != nil {
				autoScaleRecordMap["expect_scale_num"] = autoScaleRecord.ExpectScaleNum
			}

			if autoScaleRecord.EndTime != nil {
				autoScaleRecordMap["end_time"] = autoScaleRecord.EndTime
			}

			if autoScaleRecord.StrategyType != nil {
				autoScaleRecordMap["strategy_type"] = autoScaleRecord.StrategyType
			}

			if autoScaleRecord.SpecInfo != nil {
				autoScaleRecordMap["spec_info"] = autoScaleRecord.SpecInfo
			}

			if autoScaleRecord.CompensateFlag != nil {
				autoScaleRecordMap["compensate_flag"] = autoScaleRecord.CompensateFlag
			}

			if autoScaleRecord.CompensateCount != nil {
				autoScaleRecordMap["compensate_count"] = autoScaleRecord.CompensateCount
			}
			tmpList = append(tmpList, autoScaleRecordMap)
		}

		_ = d.Set("record_list", tmpList)
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
