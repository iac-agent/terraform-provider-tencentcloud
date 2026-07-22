package monitor

import (
	"context"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMonitorAlarmConditionsTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMonitorAlarmConditionsTemplateRead,
		Schema: map[string]*schema.Schema{
			"module": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Fixed 值，作为&amp;amp;#39; 监控 &amp;amp;#39;。",
			},

			"view_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "View 名称，composed 的 [DescribeAllNamespaces]( https://云.tencent.com/document/product/248/48683 )Obtain. For 云 product 监控，retrieve QceNamespacesNew. N.ID 参数 从 interface，such 作为 cvm_ Device。",
			},

			"group_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "过滤器 queries based 在 触发器 condition template names。",
			},

			"group_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "过滤器 queries based 在 触发器 condition 模板 ID",
			},

			"update_time_order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "指定sorting 方法 通过 更新时间，asc=ascending，desc=descending。",
			},

			"policy_count_order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "指定sorting 方法 based 在 数量 binding policies，asc=ascending，desc=descending。",
			},

			"template_group_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "模板 List。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"conditions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Indicator 告警 规则。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"alarm_notify_period": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Alarm 通知 频率。",
									},
									"alarm_notify_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Predefined repeated 通知 strategy (0- 告警 仅 once，1- exponential 告警，2- 连接 告警)。",
									},
									"calc_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Detection 方法。",
									},
									"calc_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Detection 值",
									},
									"continue_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Duration （秒）。",
									},
									"metric_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Indicator ID。",
									},
									"metric_display_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Indicator display 名称 (外部)。",
									},
									"period": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Cycle。",
									},
									"rule_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Rule ID。",
									},
									"unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Indicator 单位。",
									},
									"is_advanced": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "是否为an advanced indicator，0: No; 1: Yes。",
									},
									"is_open": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "是否activate advanced indicators，0: No; 1: Yes。",
									},
									"product_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Product ID。",
									},
								},
							},
						},
						"event_conditions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Event 告警 规则。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"alarm_notify_period": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Alarm 通知 频率。",
									},
									"alarm_notify_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Predefined repeated 通知 strategy (0- 告警 仅 once，1- exponential 告警，2- 连接 告警)。",
									},
									"event_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "事件 ID",
									},
									"event_display_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Event Display 名称 (External)。",
									},
									"rule_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule ID。",
									},
								},
							},
						},
						"policy_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Associate Alert Policy Group。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"can_set_default": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Can 它 是 集合 作为 默认值 告警 strategy。",
									},
									"group_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Alarm Policy 组 ID",
									},
									"group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Alarm Policy Group 名称",
									},
									"insert_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "创建时间。",
									},
									"is_default": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Is 它 默认值 告警 策略。",
									},
									"enable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Alarm Policy Enable 状态",
									},
									"last_edit_uin": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Last modified 通过 UIN。",
									},
									"no_shielded_instance_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 unshielded 实例。",
									},
									"parent_group_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Parent Policy 组 ID",
									},
									"project_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "项目 ID",
									},
									"receiver_infos": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Alarm receiving 对象 信息。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"end_time": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Effective 周期 结束时间。",
												},
												"need_send_notice": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Do 您 need 到 send 通知。",
												},
												"notify_way": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "Alarm reception channel。",
												},
												"person_interval": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Telephone 告警 到 personal 间隔 (秒)。",
												},
												"receiver_group_list": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Computed:    true,
													Description: "消息 receiving 组 列表。",
												},
												"receiver_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Receiver 类型",
												},
												"receiver_user_list": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Computed:    true,
													Description: "Recipient 列表. 列表 recipient IDs queried through 平台 interface。",
												},
												"recover_notify": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "Alarm recovery 通知 方法。",
												},
												"round_interval": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Telephone 告警 间隔 per round (秒)。",
												},
												"round_number": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "数量 phone 告警 rounds。",
												},
												"send_for": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "Timing 的 telephone 告警 通知. 可选 OCCUR (通知 during 告警)，RECOVER (通知 during recovery)。",
												},
												"start_time": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Effective 周期 开始时间。",
												},
												"uid_list": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Computed:    true,
													Description: "Telephone 告警 receiver uid。",
												},
											},
										},
									},
									"remark": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "备注",
									},
									"update_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "修改时间。",
									},
									"total_instance_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total 数量 bound 实例。",
									},
									"view_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "View。",
									},
									"is_union_rule": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Is 它 relationship 规则 使用。",
									},
								},
							},
						},
						"group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "模板 Policy 组 ID",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "模板 Policy Group 名称",
						},
						"insert_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "创建时间。",
						},
						"last_edit_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Last modified 通过 UIN。",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "备注",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "更新时间。",
						},
						"view_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "View。",
						},
						"is_union_rule": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Is 它 relationship 使用。",
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

func dataSourceTencentCloudMonitorAlarmConditionsTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_alarm_conditions_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("module"); ok {
		paramMap["Module"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("view_name"); ok {
		paramMap["ViewName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_name"); ok {
		paramMap["GroupName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_id"); ok {
		paramMap["GroupID"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("update_time_order"); ok {
		paramMap["UpdateTimeOrder"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("policy_count_order"); ok {
		paramMap["PolicyCountOrder"] = helper.String(v.(string))
	}

	service := MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var templateGroupList []*monitor.TemplateGroup

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMonitorAlarmConditionsTemplateByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		templateGroupList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(templateGroupList))
	tmpList := make([]map[string]interface{}, 0, len(templateGroupList))

	if templateGroupList != nil {
		for _, templateGroup := range templateGroupList {
			templateGroupMap := map[string]interface{}{}

			if templateGroup.Conditions != nil {
				conditionsList := []interface{}{}
				for _, conditions := range templateGroup.Conditions {
					conditionsMap := map[string]interface{}{}

					if conditions.AlarmNotifyPeriod != nil {
						conditionsMap["alarm_notify_period"] = conditions.AlarmNotifyPeriod
					}

					if conditions.AlarmNotifyType != nil {
						conditionsMap["alarm_notify_type"] = conditions.AlarmNotifyType
					}

					if conditions.CalcType != nil {
						conditionsMap["calc_type"] = conditions.CalcType
					}

					if conditions.CalcValue != nil {
						conditionsMap["calc_value"] = conditions.CalcValue
					}

					if conditions.ContinueTime != nil {
						conditionsMap["continue_time"] = conditions.ContinueTime
					}

					if conditions.MetricID != nil {
						conditionsMap["metric_id"] = conditions.MetricID
					}

					if conditions.MetricDisplayName != nil {
						conditionsMap["metric_display_name"] = conditions.MetricDisplayName
					}

					if conditions.Period != nil {
						conditionsMap["period"] = conditions.Period
					}

					if conditions.RuleID != nil {
						conditionsMap["rule_id"] = conditions.RuleID
					}

					if conditions.Unit != nil {
						conditionsMap["unit"] = conditions.Unit
					}

					if conditions.IsAdvanced != nil {
						conditionsMap["is_advanced"] = conditions.IsAdvanced
					}

					if conditions.IsOpen != nil {
						conditionsMap["is_open"] = conditions.IsOpen
					}

					if conditions.ProductId != nil {
						conditionsMap["product_id"] = conditions.ProductId
					}

					conditionsList = append(conditionsList, conditionsMap)
				}

				templateGroupMap["conditions"] = conditionsList
			}

			if templateGroup.EventConditions != nil {
				eventConditionsList := []interface{}{}
				for _, eventConditions := range templateGroup.EventConditions {
					eventConditionsMap := map[string]interface{}{}

					if eventConditions.AlarmNotifyPeriod != nil {
						eventConditionsMap["alarm_notify_period"] = eventConditions.AlarmNotifyPeriod
					}

					if eventConditions.AlarmNotifyType != nil {
						eventConditionsMap["alarm_notify_type"] = eventConditions.AlarmNotifyType
					}

					if eventConditions.EventID != nil {
						eventConditionsMap["event_id"] = eventConditions.EventID
					}

					if eventConditions.EventDisplayName != nil {
						eventConditionsMap["event_display_name"] = eventConditions.EventDisplayName
					}

					if eventConditions.RuleID != nil {
						eventConditionsMap["rule_id"] = eventConditions.RuleID
					}

					eventConditionsList = append(eventConditionsList, eventConditionsMap)
				}

				templateGroupMap["event_conditions"] = eventConditionsList
			}

			if templateGroup.PolicyGroups != nil {
				policyGroupsList := []interface{}{}
				for _, policyGroups := range templateGroup.PolicyGroups {
					policyGroupsMap := map[string]interface{}{}

					if policyGroups.CanSetDefault != nil {
						policyGroupsMap["can_set_default"] = policyGroups.CanSetDefault
					}

					if policyGroups.GroupID != nil {
						policyGroupsMap["group_id"] = policyGroups.GroupID
					}

					if policyGroups.GroupName != nil {
						policyGroupsMap["group_name"] = policyGroups.GroupName
					}

					if policyGroups.InsertTime != nil {
						policyGroupsMap["insert_time"] = policyGroups.InsertTime
					}

					if policyGroups.IsDefault != nil {
						policyGroupsMap["is_default"] = policyGroups.IsDefault
					}

					if policyGroups.Enable != nil {
						policyGroupsMap["enable"] = policyGroups.Enable
					}

					if policyGroups.LastEditUin != nil {
						policyGroupsMap["last_edit_uin"] = policyGroups.LastEditUin
					}

					if policyGroups.NoShieldedInstanceCount != nil {
						policyGroupsMap["no_shielded_instance_count"] = policyGroups.NoShieldedInstanceCount
					}

					if policyGroups.ParentGroupID != nil {
						policyGroupsMap["parent_group_id"] = policyGroups.ParentGroupID
					}

					if policyGroups.ProjectID != nil {
						policyGroupsMap["project_id"] = policyGroups.ProjectID
					}

					if policyGroups.ReceiverInfos != nil {
						receiverInfosList := []interface{}{}
						for _, receiverInfos := range policyGroups.ReceiverInfos {
							receiverInfosMap := map[string]interface{}{}

							if receiverInfos.EndTime != nil {
								receiverInfosMap["end_time"] = receiverInfos.EndTime
							}

							if receiverInfos.NeedSendNotice != nil {
								receiverInfosMap["need_send_notice"] = receiverInfos.NeedSendNotice
							}

							if receiverInfos.NotifyWay != nil {
								receiverInfosMap["notify_way"] = receiverInfos.NotifyWay
							}

							if receiverInfos.PersonInterval != nil {
								receiverInfosMap["person_interval"] = receiverInfos.PersonInterval
							}

							if receiverInfos.ReceiverGroupList != nil {
								receiverInfosMap["receiver_group_list"] = receiverInfos.ReceiverGroupList
							}

							if receiverInfos.ReceiverType != nil {
								receiverInfosMap["receiver_type"] = receiverInfos.ReceiverType
							}

							if receiverInfos.ReceiverUserList != nil {
								receiverInfosMap["receiver_user_list"] = receiverInfos.ReceiverUserList
							}

							if receiverInfos.RecoverNotify != nil {
								receiverInfosMap["recover_notify"] = receiverInfos.RecoverNotify
							}

							if receiverInfos.RoundInterval != nil {
								receiverInfosMap["round_interval"] = receiverInfos.RoundInterval
							}

							if receiverInfos.RoundNumber != nil {
								receiverInfosMap["round_number"] = receiverInfos.RoundNumber
							}

							if receiverInfos.SendFor != nil {
								receiverInfosMap["send_for"] = receiverInfos.SendFor
							}

							if receiverInfos.StartTime != nil {
								receiverInfosMap["start_time"] = receiverInfos.StartTime
							}

							if receiverInfos.UIDList != nil {
								receiverInfosMap["uid_list"] = receiverInfos.UIDList
							}

							receiverInfosList = append(receiverInfosList, receiverInfosMap)
						}

						policyGroupsMap["receiver_infos"] = receiverInfosList
					}

					if policyGroups.Remark != nil {
						policyGroupsMap["remark"] = policyGroups.Remark
					}

					if policyGroups.UpdateTime != nil {
						policyGroupsMap["update_time"] = policyGroups.UpdateTime
					}

					if policyGroups.TotalInstanceCount != nil {
						policyGroupsMap["total_instance_count"] = policyGroups.TotalInstanceCount
					}

					if policyGroups.ViewName != nil {
						policyGroupsMap["view_name"] = policyGroups.ViewName
					}

					if policyGroups.IsUnionRule != nil {
						policyGroupsMap["is_union_rule"] = policyGroups.IsUnionRule
					}

					policyGroupsList = append(policyGroupsList, policyGroupsMap)
				}

				templateGroupMap["policy_groups"] = policyGroupsList
			}

			if templateGroup.GroupID != nil {
				templateGroupMap["group_id"] = templateGroup.GroupID
			}

			if templateGroup.GroupName != nil {
				templateGroupMap["group_name"] = templateGroup.GroupName
			}

			if templateGroup.InsertTime != nil {
				templateGroupMap["insert_time"] = templateGroup.InsertTime
			}

			if templateGroup.LastEditUin != nil {
				templateGroupMap["last_edit_uin"] = templateGroup.LastEditUin
			}

			if templateGroup.Remark != nil {
				templateGroupMap["remark"] = templateGroup.Remark
			}

			if templateGroup.UpdateTime != nil {
				templateGroupMap["update_time"] = templateGroup.UpdateTime
			}

			if templateGroup.ViewName != nil {
				templateGroupMap["view_name"] = templateGroup.ViewName
			}

			if templateGroup.IsUnionRule != nil {
				templateGroupMap["is_union_rule"] = templateGroup.IsUnionRule
			}

			ids = append(ids, strconv.Itoa(int(*templateGroup.GroupID)))
			tmpList = append(tmpList, templateGroupMap)
		}

		_ = d.Set("template_group_list", tmpList)
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
