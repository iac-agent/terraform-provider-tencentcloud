package as

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudAsLastActivity() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAsLastActivityRead,
		Schema: map[string]*schema.Schema{
			"auto_scaling_group_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "ID 列表 auto scaling 组。",
			},

			"exclude_cancelled_activity": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Exclude cancellation 类型 activities 当 querying. 默认值为 false，indicating 该 cancellation 类型 activities 是 不 excluded。",
			},

			"activity_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Information 集合 的 eligible scaling activities. Scaling groups without scaling activities 是 不 返回. For 示例，如果 there 是 50 auto scaling 组 IDs 但 仅 45 records 是 返回，它 表示that 5 的 auto scaling groups do 不 have scaling activities。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_scaling_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Auto scaling 组 ID",
						},
						"activity_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scaling activity ID。",
						},
						"activity_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 scaling activity. 取值范围：SCALE_OUT，SCALE_IN，ATTACH_INSTANCES，REMOVE_INSTANCES，DETACH_INSTANCES，TERMINATE_INSTANCES_UNEXPECTEDLY，REPLACE_UNHEALTHY_INSTANCE，START_INSTANCES，STOP_INSTANCES，INVOKE_COMMAND。",
						},
						"status_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scaling activity 状态 取值范围：INIT，RUNNING，SUCCESSFUL，PARTIALLY_SUCCESSFUL，FAILED，CANCELLED。",
						},
						"status_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 scaling activity 状态",
						},
						"cause": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cause 的 scaling activity。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 scaling activity。",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "开始时间 的 scaling activity。",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "结束时间 的 scaling activity。",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 scaling activity。",
						},
						"activity_related_instance_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Information 集合 的 实例 related 到 scaling activity。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例 ID",
									},
									"instance_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "状态 实例 在 scaling activity. 取值范围：INIT，RUNNING，SUCCESSFUL，FAILED。",
									},
								},
							},
						},
						"status_message_simplified": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Brief 描述 scaling activity 状态",
						},
						"lifecycle_action_result_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "结果 的 lifecycle hook 操作 在 scaling activity。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lifecycle_hook_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID lifecycle hook。",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID 实例。",
									},
									"invocation_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Execution 任务 ID. You 可以 查询 结果 通过 使用 DescribeInvocations API 的 TAT。",
									},
									"invoke_command_result": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "结果 的 command invocation，取值范围：SUCCESSFUL，FAILED，NONE。",
									},
									"notification_result": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Notification 结果，其中 表示whether 它 是 successful 到 notify CMQ/TDMQ，取值范围：SUCCESSFUL，FAILED，NONE。",
									},
									"lifecycle_action_result": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "结果 的 lifecycle hook 操作，取值范围：CONTINUE，ABANDON。",
									},
									"result_reason": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Reason 的 结果，取值范围：HEARTBEAT_TIMEOUT: Heartbeat timed out. setting 的 DefaultResult 是 使用. NOTIFICATION_FAILURE: Failed 到 send 通知. setting 的 DefaultResult 是 使用. CALL_INTERFACE: Calls CompleteLifecycleAction 到 集合 结果 ANOTHER_ACTION_ABANDON: It has been 集合 到 ABANDON 通过 another operation. COMMAND_CALL_FAILURE: Failed 到 call command. DefaultResult 是 applied. COMMAND_EXEC_FINISH: Command completed COMMAND_CALL_FAILURE: Failed 到 execute command. DefaultResult 是 applied. COMMAND_EXEC_RESULT_CHECK_FAILURE: Failed 到 check command 结果 DefaultResult 是 applied。",
									},
								},
							},
						},
						"detailed_status_message_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Detailed 描述 scaling activity 状态",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "错误 类型",
									},
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "AZ 信息。",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例 ID",
									},
									"instance_charge_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例 billing 模式",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "子网 ID",
									},
									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "错误信息",
									},
									"instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例类型",
									},
								},
							},
						},
						"invocation_result_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "结果 的 command execution。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例 ID 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"invocation_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Execution activity ID. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"invocation_task_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Execution 任务 ID. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"command_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "命令 ID 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"task_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Execution 状态 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"error_message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Execution exception 信息. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
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

func dataSourceTencentCloudAsLastActivityRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_as_last_activity.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = AsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("auto_scaling_group_ids"); ok {
		autoScalingGroupIdsSet := v.(*schema.Set).List()
		paramMap["AutoScalingGroupIds"] = helper.InterfacesStringsPoint(autoScalingGroupIdsSet)
	}

	if v, ok := d.GetOk("exclude_cancelled_activity"); ok {
		paramMap["ExcludeCancelledActivity"] = helper.Bool(v.(bool))
	}

	var activitySet []*as.Activity
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAsLastActivity(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		activitySet = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(activitySet))
	tmpList := make([]map[string]interface{}, 0, len(activitySet))
	if activitySet != nil {
		for _, activity := range activitySet {
			activityMap := map[string]interface{}{}
			if activity.AutoScalingGroupId != nil {
				activityMap["auto_scaling_group_id"] = activity.AutoScalingGroupId
			}

			if activity.ActivityId != nil {
				activityMap["activity_id"] = activity.ActivityId
			}

			if activity.ActivityType != nil {
				activityMap["activity_type"] = activity.ActivityType
			}

			if activity.StatusCode != nil {
				activityMap["status_code"] = activity.StatusCode
			}

			if activity.StatusMessage != nil {
				activityMap["status_message"] = activity.StatusMessage
			}

			if activity.Cause != nil {
				activityMap["cause"] = activity.Cause
			}

			if activity.Description != nil {
				activityMap["description"] = activity.Description
			}

			if activity.StartTime != nil {
				activityMap["start_time"] = activity.StartTime
			}

			if activity.EndTime != nil {
				activityMap["end_time"] = activity.EndTime
			}

			if activity.CreatedTime != nil {
				activityMap["created_time"] = activity.CreatedTime
			}

			if activity.ActivityRelatedInstanceSet != nil {
				activityRelatedInstanceSetList := []interface{}{}
				for _, activityRelatedInstanceSet := range activity.ActivityRelatedInstanceSet {
					activityRelatedInstanceSetMap := map[string]interface{}{}

					if activityRelatedInstanceSet.InstanceId != nil {
						activityRelatedInstanceSetMap["instance_id"] = activityRelatedInstanceSet.InstanceId
					}

					if activityRelatedInstanceSet.InstanceStatus != nil {
						activityRelatedInstanceSetMap["instance_status"] = activityRelatedInstanceSet.InstanceStatus
					}

					activityRelatedInstanceSetList = append(activityRelatedInstanceSetList, activityRelatedInstanceSetMap)
				}

				activityMap["activity_related_instance_set"] = activityRelatedInstanceSetList
			}

			if activity.StatusMessageSimplified != nil {
				activityMap["status_message_simplified"] = activity.StatusMessageSimplified
			}

			if activity.LifecycleActionResultSet != nil {
				lifecycleActionResultSetList := []interface{}{}
				for _, lifecycleActionResultSet := range activity.LifecycleActionResultSet {
					lifecycleActionResultSetMap := map[string]interface{}{}

					if lifecycleActionResultSet.LifecycleHookId != nil {
						lifecycleActionResultSetMap["lifecycle_hook_id"] = lifecycleActionResultSet.LifecycleHookId
					}

					if lifecycleActionResultSet.InstanceId != nil {
						lifecycleActionResultSetMap["instance_id"] = lifecycleActionResultSet.InstanceId
					}

					if lifecycleActionResultSet.InvocationId != nil {
						lifecycleActionResultSetMap["invocation_id"] = lifecycleActionResultSet.InvocationId
					}

					if lifecycleActionResultSet.InvokeCommandResult != nil {
						lifecycleActionResultSetMap["invoke_command_result"] = lifecycleActionResultSet.InvokeCommandResult
					}

					if lifecycleActionResultSet.NotificationResult != nil {
						lifecycleActionResultSetMap["notification_result"] = lifecycleActionResultSet.NotificationResult
					}

					if lifecycleActionResultSet.LifecycleActionResult != nil {
						lifecycleActionResultSetMap["lifecycle_action_result"] = lifecycleActionResultSet.LifecycleActionResult
					}

					if lifecycleActionResultSet.ResultReason != nil {
						lifecycleActionResultSetMap["result_reason"] = lifecycleActionResultSet.ResultReason
					}

					lifecycleActionResultSetList = append(lifecycleActionResultSetList, lifecycleActionResultSetMap)
				}

				activityMap["lifecycle_action_result_set"] = lifecycleActionResultSetList
			}

			if activity.DetailedStatusMessageSet != nil {
				detailedStatusMessageSetList := []interface{}{}
				for _, detailedStatusMessageSet := range activity.DetailedStatusMessageSet {
					detailedStatusMessageSetMap := map[string]interface{}{}

					if detailedStatusMessageSet.Code != nil {
						detailedStatusMessageSetMap["code"] = detailedStatusMessageSet.Code
					}

					if detailedStatusMessageSet.Zone != nil {
						detailedStatusMessageSetMap["zone"] = detailedStatusMessageSet.Zone
					}

					if detailedStatusMessageSet.InstanceId != nil {
						detailedStatusMessageSetMap["instance_id"] = detailedStatusMessageSet.InstanceId
					}

					if detailedStatusMessageSet.InstanceChargeType != nil {
						detailedStatusMessageSetMap["instance_charge_type"] = detailedStatusMessageSet.InstanceChargeType
					}

					if detailedStatusMessageSet.SubnetId != nil {
						detailedStatusMessageSetMap["subnet_id"] = detailedStatusMessageSet.SubnetId
					}

					if detailedStatusMessageSet.Message != nil {
						detailedStatusMessageSetMap["message"] = detailedStatusMessageSet.Message
					}

					if detailedStatusMessageSet.InstanceType != nil {
						detailedStatusMessageSetMap["instance_type"] = detailedStatusMessageSet.InstanceType
					}

					detailedStatusMessageSetList = append(detailedStatusMessageSetList, detailedStatusMessageSetMap)
				}

				activityMap["detailed_status_message_set"] = detailedStatusMessageSetList
			}

			if activity.InvocationResultSet != nil {
				invocationResultSetList := []interface{}{}
				for _, invocationResultSet := range activity.InvocationResultSet {
					invocationResultSetMap := map[string]interface{}{}

					if invocationResultSet.InstanceId != nil {
						invocationResultSetMap["instance_id"] = invocationResultSet.InstanceId
					}

					if invocationResultSet.InvocationId != nil {
						invocationResultSetMap["invocation_id"] = invocationResultSet.InvocationId
					}

					if invocationResultSet.InvocationTaskId != nil {
						invocationResultSetMap["invocation_task_id"] = invocationResultSet.InvocationTaskId
					}

					if invocationResultSet.CommandId != nil {
						invocationResultSetMap["command_id"] = invocationResultSet.CommandId
					}

					if invocationResultSet.TaskStatus != nil {
						invocationResultSetMap["task_status"] = invocationResultSet.TaskStatus
					}

					if invocationResultSet.ErrorMessage != nil {
						invocationResultSetMap["error_message"] = invocationResultSet.ErrorMessage
					}

					invocationResultSetList = append(invocationResultSetList, invocationResultSetMap)
				}

				activityMap["invocation_result_set"] = invocationResultSetList
			}

			ids = append(ids, *activity.AutoScalingGroupId)
			tmpList = append(tmpList, activityMap)
		}

		_ = d.Set("activity_set", tmpList)
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
