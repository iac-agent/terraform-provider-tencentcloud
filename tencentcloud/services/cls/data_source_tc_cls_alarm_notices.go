package cls

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClsAlarmNotices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClsAlarmNoticesRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "过滤条件。最多 10 个过滤器，每个过滤器最多有 5 个值。同一过滤器内的多个值使用 OR 逻辑，多个过滤器使用 AND 逻辑。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤字段名称。支持的值：“名称”（告警通知组名称）、“alarmNoticeId”（告警通知 ID）、“uid”（接收方用户 ID）、“groupId”（接收方用户组 ID）、“deliverFlag”（下发状态：1-未启用、2-启用、3-异常）。",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "过滤字段值。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"has_alarm_shield_count": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "是否查询报警屏蔽计数统计。默认为 false。",
			},

			"alarm_notices": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "报警通知配置列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "报警通知名称。",
						},
						"alarm_notice_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "报警通知ID。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创作时间。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "最后更新时间。",
						},
						"notice_receivers": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "通知接收者列表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"receiver_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "接收器类型。可以是 Uin 或 Group。",
									},
									"receiver_ids": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "接收者 ID。",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"receiver_channels": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "通知渠道（电子邮件、短信、微信、电话）。",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"start_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "允许的通知开始时间。",
									},
									"end_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "允许的通知结束时间。",
									},
									"index": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "索引顺序。",
									},
								},
							},
						},
						"web_callbacks": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Webhook 回调列表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "回调网址。",
									},
									"callback_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "回调类型。 WeCom 或 Http 或钉钉或 Lark 或 Webhook。",
									},
									"method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "HTTP 方法。获取或发布。",
									},
									"headers": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "请求标头。",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"body": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "请求正文。",
									},
									"index": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "索引顺序。",
									},
								},
							},
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "标签列表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签键。",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签值。",
									},
								},
							},
						},
						"jump_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "报警回调跳转域。",
						},
						"notice_rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "通知规则列表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"notice_receivers": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "请通知接收者此规则。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"receiver_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "接收器类型。",
												},
												"receiver_ids": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "接收者 ID。",
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
												"receiver_channels": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "通知渠道。",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"start_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "开始时间。",
												},
												"end_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "结束时间。",
												},
												"index": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "指数。",
												},
											},
										},
									},
									"web_callbacks": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "此规则的 Webhook 回调。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "回调网址。",
												},
												"callback_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "回调类型。",
												},
												"method": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "HTTP 方法。",
												},
												"headers": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "标头。",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"body": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "身体。",
												},
												"index": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "指数。",
												},
											},
										},
									},
									"repeat_interval": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "重复间隔以分钟为单位。",
									},
									"time_range_start": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "有效开始时间（24 小时格式 HH:mm:ss）。",
									},
									"time_range_end": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "有效结束时间（24 小时格式 HH:mm:ss）。",
									},
									"notify_way": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "通知方式。",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"receiver_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "接收器类型。",
									},
									"day_of_week": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "一周中的天数（0-6，0 是星期日）。",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"jump_domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "跳转域。",
									},
								},
							},
						},
						"deliver_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "发货状态（0：已发货，1：未发货）。",
						},
						"deliver_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "发送标志（1：未使能，2：使能，3：异常）。",
						},
						"alarm_shield_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "报警屏蔽状态（0：不屏蔽，1：屏蔽）。",
						},
						"alarm_shield_count": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "报警屏蔽计数统计。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "屏蔽告警总数。",
									},
								},
							},
						},

						"callback_prioritize": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "webhook 回调是否优先。",
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

func dataSourceTencentCloudClsAlarmNoticesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cls_alarm_notices.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*cls.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := cls.Filter{}
			if v, ok := filtersMap["key"].(string); ok && v != "" {
				filter.Key = helper.String(v)
			}

			if v, ok := filtersMap["values"]; ok {
				valueSet := v.(*schema.Set).List()
				for i := range valueSet {
					value := valueSet[i].(string)
					filter.Values = append(filter.Values, helper.String(value))
				}
			}
			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	if v, ok := d.GetOkExists("has_alarm_shield_count"); ok {
		paramMap["HasAlarmShieldCount"] = helper.Bool(v.(bool))
	}

	var respData []*cls.AlarmNotice

	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClsAlarmNoticesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	alarmNoticesList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, alarmNotice := range respData {
			alarmNoticeMap := map[string]interface{}{}

			if alarmNotice.Name != nil {
				alarmNoticeMap["name"] = alarmNotice.Name
			}

			if alarmNotice.AlarmNoticeId != nil {
				alarmNoticeMap["alarm_notice_id"] = alarmNotice.AlarmNoticeId
			}

			if alarmNotice.CreateTime != nil {
				alarmNoticeMap["create_time"] = alarmNotice.CreateTime
			}

			if alarmNotice.UpdateTime != nil {
				alarmNoticeMap["update_time"] = alarmNotice.UpdateTime
			}

			if alarmNotice.NoticeReceivers != nil {
				noticeReceiversList := []interface{}{}
				for _, noticeReceiver := range alarmNotice.NoticeReceivers {
					noticeReceiverMap := map[string]interface{}{}

					if noticeReceiver.ReceiverType != nil {
						noticeReceiverMap["receiver_type"] = noticeReceiver.ReceiverType
					}

					if noticeReceiver.ReceiverIds != nil {
						noticeReceiverMap["receiver_ids"] = noticeReceiver.ReceiverIds
					}

					if noticeReceiver.ReceiverChannels != nil {
						noticeReceiverMap["receiver_channels"] = noticeReceiver.ReceiverChannels
					}

					if noticeReceiver.StartTime != nil {
						noticeReceiverMap["start_time"] = noticeReceiver.StartTime
					}

					if noticeReceiver.EndTime != nil {
						noticeReceiverMap["end_time"] = noticeReceiver.EndTime
					}

					if noticeReceiver.Index != nil {
						noticeReceiverMap["index"] = noticeReceiver.Index
					}

					noticeReceiversList = append(noticeReceiversList, noticeReceiverMap)
				}

				alarmNoticeMap["notice_receivers"] = noticeReceiversList
			}

			if alarmNotice.WebCallbacks != nil {
				webCallbacksList := []interface{}{}
				for _, webCallback := range alarmNotice.WebCallbacks {
					webCallbackMap := map[string]interface{}{}

					if webCallback.Url != nil {
						webCallbackMap["url"] = webCallback.Url
					}

					if webCallback.CallbackType != nil {
						webCallbackMap["callback_type"] = webCallback.CallbackType
					}

					if webCallback.Method != nil {
						webCallbackMap["method"] = webCallback.Method
					}

					if webCallback.Headers != nil {
						webCallbackMap["headers"] = webCallback.Headers
					}

					if webCallback.Body != nil {
						webCallbackMap["body"] = webCallback.Body
					}

					if webCallback.Index != nil {
						webCallbackMap["index"] = webCallback.Index
					}

					webCallbacksList = append(webCallbacksList, webCallbackMap)
				}

				alarmNoticeMap["web_callbacks"] = webCallbacksList
			}

			if alarmNotice.Tags != nil {
				tagsList := []interface{}{}
				for _, tag := range alarmNotice.Tags {
					tagMap := map[string]interface{}{}

					if tag.Key != nil {
						tagMap["key"] = tag.Key
					}

					if tag.Value != nil {
						tagMap["value"] = tag.Value
					}

					tagsList = append(tagsList, tagMap)
				}

				alarmNoticeMap["tags"] = tagsList
			}

			if alarmNotice.JumpDomain != nil {
				alarmNoticeMap["jump_domain"] = alarmNotice.JumpDomain
			}

			if alarmNotice.NoticeRules != nil {
				noticeRulesList := []interface{}{}
				for _, noticeRule := range alarmNotice.NoticeRules {
					noticeRuleMap := map[string]interface{}{}

					if noticeRule.NoticeReceivers != nil {
						ruleNoticeReceiversList := []interface{}{}
						for _, receiver := range noticeRule.NoticeReceivers {
							receiverMap := map[string]interface{}{}

							if receiver.ReceiverType != nil {
								receiverMap["receiver_type"] = receiver.ReceiverType
							}

							if receiver.ReceiverIds != nil {
								receiverMap["receiver_ids"] = receiver.ReceiverIds
							}

							if receiver.ReceiverChannels != nil {
								receiverMap["receiver_channels"] = receiver.ReceiverChannels
							}

							if receiver.StartTime != nil {
								receiverMap["start_time"] = receiver.StartTime
							}

							if receiver.EndTime != nil {
								receiverMap["end_time"] = receiver.EndTime
							}

							if receiver.Index != nil {
								receiverMap["index"] = receiver.Index
							}

							ruleNoticeReceiversList = append(ruleNoticeReceiversList, receiverMap)
						}

						noticeRuleMap["notice_receivers"] = ruleNoticeReceiversList
					}

					if noticeRule.WebCallbacks != nil {
						ruleWebCallbacksList := []interface{}{}
						for _, callback := range noticeRule.WebCallbacks {
							callbackMap := map[string]interface{}{}

							if callback.Url != nil {
								callbackMap["url"] = callback.Url
							}

							if callback.CallbackType != nil {
								callbackMap["callback_type"] = callback.CallbackType
							}

							if callback.Method != nil {
								callbackMap["method"] = callback.Method
							}

							if callback.Headers != nil {
								callbackMap["headers"] = callback.Headers
							}

							if callback.Body != nil {
								callbackMap["body"] = callback.Body
							}

							if callback.Index != nil {
								callbackMap["index"] = callback.Index
							}

							ruleWebCallbacksList = append(ruleWebCallbacksList, callbackMap)
						}

						noticeRuleMap["web_callbacks"] = ruleWebCallbacksList
					}

					noticeRulesList = append(noticeRulesList, noticeRuleMap)
				}

				alarmNoticeMap["notice_rules"] = noticeRulesList
			}

			if alarmNotice.DeliverStatus != nil {
				alarmNoticeMap["deliver_status"] = alarmNotice.DeliverStatus
			}

			if alarmNotice.DeliverFlag != nil {
				alarmNoticeMap["deliver_flag"] = alarmNotice.DeliverFlag
			}

			if alarmNotice.AlarmShieldStatus != nil {
				alarmNoticeMap["alarm_shield_status"] = alarmNotice.AlarmShieldStatus
			}

			if alarmNotice.AlarmShieldCount != nil {
				alarmShieldCountList := []interface{}{}
				alarmShieldCountMap := map[string]interface{}{}

				if alarmNotice.AlarmShieldCount.TotalCount != nil {
					alarmShieldCountMap["total_count"] = alarmNotice.AlarmShieldCount.TotalCount
				}

				alarmShieldCountList = append(alarmShieldCountList, alarmShieldCountMap)
				alarmNoticeMap["alarm_shield_count"] = alarmShieldCountList
			}

			if alarmNotice.CallbackPrioritize != nil {
				alarmNoticeMap["callback_prioritize"] = alarmNotice.CallbackPrioritize
			}

			alarmNoticesList = append(alarmNoticesList, alarmNoticeMap)
		}

		_ = d.Set("alarm_notices", alarmNoticesList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
