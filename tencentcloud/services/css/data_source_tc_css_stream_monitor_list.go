package css

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCssStreamMonitorList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCssStreamMonitorListRead,
		Schema: map[string]*schema.Schema{

			"live_stream_monitors": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 live 流 监控 tasks.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"monitor_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Monitoring 任务 ID.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
						},
						"monitor_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Monitoring 任务 名称 Up 到 128 bytes.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
						},
						"output_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Monitoring 任务 output 信息.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"output_stream_width": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "宽度 的 output 流 在 pixels 对于 监控 任务. 范围 是 [1，1920]. It 是 recommended 到 是 在 least 100 pixels.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
									},
									"output_stream_height": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "高度 的 output 流 在 pixels 对于 监控 任务. 范围 是 [1，1080]. It 是 recommended 到 是 在 least 100 pixels.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
									},
									"output_stream_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "名称 output 流 对于 监控 任务.如果未指定， 系统 将 generate 名称 automatically. 名称 should 是 within 256 bytes 和 可以 仅 contain letters，numbers，`-`，`_`，和 `.` 字符.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
									},
									"output_domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "playback 域名 对于 监控 任务.It should 是 within 128 bytes 和 可以 仅 是 filled 使用 已启用 playback 域名Note: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
									},
									"output_app": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "playback 路径 对于 监控 任务.It should 是 within 32 bytes 和 可以 仅 contain letters，numbers，`-`，`_`，和 `.` 字符.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
									},
								},
							},
						},
						"input_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "input 流 信息 对于 监控 任务.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"input_stream_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "名称 input 流 对于 监控 任务.It should 是 within 256 bytes 和 可以 仅 contain letters，numbers，`-`，`_`，和 `.` 字符.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
									},
									"input_domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "push 域名 对于 input 流 到 是 monitored.It should 是 within 128 bytes 和 可以 仅 是 filled 使用 已启用 push 域名Note: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
									},
									"input_app": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "push 路径 对于 input 流 到 是 monitored.It should 是 within 32 bytes 和 可以 仅 contain letters，numbers，`-`，`_`，和 `.` 字符.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
									},
									"input_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "push URL 对于 input 流 到 是 monitored. In most cases，此 参数 不是必填项.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "描述 监控 任务.It should 是 within 256 bytes.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "状态 监控 任务. 0: Represents idle. 1: Represents 监控 在 progress.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
						},
						"start_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "last 开始时间 的 监控 任务，在 Unix 时间戳 格式Note: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
						},
						"stop_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "last stop 时间 的 监控 任务，在 Unix 时间戳 格式Note: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "创建时间 的 监控 任务，在 Unix 时间戳 格式Note: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "更新时间 的 监控 任务，在 Unix 时间戳 格式Note: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
						},
						"notify_policy": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "通知 策略 对于 监控 events.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"notify_policy_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "类型 通知 策略: Range [0,1] 0: Represents 无 通知 策略 是 使用. 1: Represents 使用 的 全局 callback 策略，其中 all events 是 notified 到 CallbackUrl.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
									},
									"callback_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "callback URL 对于 notifications. It should 是 的 长度 [0,512] 和 仅 support URLs 使用 http 和 https types.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
									},
								},
							},
						},
						"audible_input_index_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "列表 input indices 对于 output 音频.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
						},
						"ai_asr_input_index_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "列表 input indices 对于 enabling intelligent speech recognition.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
						},
						"check_stream_broken": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否enable 流 disconnection detection.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
						},
						"check_stream_low_frame_rate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否enable low frame 速率 detection.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
						},
						"asr_language": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "语言 对于 intelligent speech recognition:0: Disabled1: Chinese2: English3: Japanese4: KoreanNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
						},
						"ocr_language": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "语言 对于 intelligent text recognition:0: Disabled1: Chinese 和 EnglishNote: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
						},
						"ai_ocr_input_index_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "列表 input indices 对于 enabling intelligent text recognition.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
						},
						"allow_monitor_report": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否store 监控 events 在 监控 报告 和 allow querying 的 监控 报告.注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
						},
						"ai_format_diagnose": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否enable 格式 diagnosis. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 是 可用。",
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

func dataSourceTencentCloudCssStreamMonitorListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_css_stream_monitor_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})

	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var streamMonitorList []*css.LiveStreamMonitorInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCssStreamMonitorListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		streamMonitorList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(streamMonitorList))
	tmpList := make([]map[string]interface{}, 0, len(streamMonitorList))
	if streamMonitorList != nil {
		for _, liveStreamMonitorInfo := range streamMonitorList {
			liveStreamMonitorInfoMap := map[string]interface{}{}

			if liveStreamMonitorInfo.MonitorId != nil {
				liveStreamMonitorInfoMap["monitor_id"] = liveStreamMonitorInfo.MonitorId
			}

			if liveStreamMonitorInfo.MonitorName != nil {
				liveStreamMonitorInfoMap["monitor_name"] = liveStreamMonitorInfo.MonitorName
			}

			if liveStreamMonitorInfo.OutputInfo != nil {
				outputInfoMap := map[string]interface{}{}

				if liveStreamMonitorInfo.OutputInfo.OutputStreamWidth != nil {
					outputInfoMap["output_stream_width"] = liveStreamMonitorInfo.OutputInfo.OutputStreamWidth
				}

				if liveStreamMonitorInfo.OutputInfo.OutputStreamHeight != nil {
					outputInfoMap["output_stream_height"] = liveStreamMonitorInfo.OutputInfo.OutputStreamHeight
				}

				if liveStreamMonitorInfo.OutputInfo.OutputStreamName != nil {
					outputInfoMap["output_stream_name"] = liveStreamMonitorInfo.OutputInfo.OutputStreamName
				}

				if liveStreamMonitorInfo.OutputInfo.OutputDomain != nil {
					outputInfoMap["output_domain"] = liveStreamMonitorInfo.OutputInfo.OutputDomain
				}

				if liveStreamMonitorInfo.OutputInfo.OutputApp != nil {
					outputInfoMap["output_app"] = liveStreamMonitorInfo.OutputInfo.OutputApp
				}

				liveStreamMonitorInfoMap["output_info"] = []interface{}{outputInfoMap}
			}

			if liveStreamMonitorInfo.InputList != nil {
				inputListList := []interface{}{}
				for _, inputList := range liveStreamMonitorInfo.InputList {
					inputListMap := map[string]interface{}{}

					if inputList.InputStreamName != nil {
						inputListMap["input_stream_name"] = inputList.InputStreamName
					}

					if inputList.InputDomain != nil {
						inputListMap["input_domain"] = inputList.InputDomain
					}

					if inputList.InputApp != nil {
						inputListMap["input_app"] = inputList.InputApp
					}

					if inputList.InputUrl != nil {
						inputListMap["input_url"] = inputList.InputUrl
					}

					if inputList.Description != nil {
						inputListMap["description"] = inputList.Description
					}

					inputListList = append(inputListList, inputListMap)
				}

				liveStreamMonitorInfoMap["input_list"] = inputListList
			}

			if liveStreamMonitorInfo.Status != nil {
				liveStreamMonitorInfoMap["status"] = liveStreamMonitorInfo.Status
			}

			if liveStreamMonitorInfo.StartTime != nil {
				liveStreamMonitorInfoMap["start_time"] = liveStreamMonitorInfo.StartTime
			}

			if liveStreamMonitorInfo.StopTime != nil {
				liveStreamMonitorInfoMap["stop_time"] = liveStreamMonitorInfo.StopTime
			}

			if liveStreamMonitorInfo.CreateTime != nil {
				liveStreamMonitorInfoMap["create_time"] = liveStreamMonitorInfo.CreateTime
			}

			if liveStreamMonitorInfo.UpdateTime != nil {
				liveStreamMonitorInfoMap["update_time"] = liveStreamMonitorInfo.UpdateTime
			}

			if liveStreamMonitorInfo.NotifyPolicy != nil {
				notifyPolicyMap := map[string]interface{}{}

				if liveStreamMonitorInfo.NotifyPolicy.NotifyPolicyType != nil {
					notifyPolicyMap["notify_policy_type"] = liveStreamMonitorInfo.NotifyPolicy.NotifyPolicyType
				}

				if liveStreamMonitorInfo.NotifyPolicy.CallbackUrl != nil {
					notifyPolicyMap["callback_url"] = liveStreamMonitorInfo.NotifyPolicy.CallbackUrl
				}

				liveStreamMonitorInfoMap["notify_policy"] = []interface{}{notifyPolicyMap}
			}

			if liveStreamMonitorInfo.AudibleInputIndexList != nil {
				liveStreamMonitorInfoMap["audible_input_index_list"] = liveStreamMonitorInfo.AudibleInputIndexList
			}

			if liveStreamMonitorInfo.AiAsrInputIndexList != nil {
				liveStreamMonitorInfoMap["ai_asr_input_index_list"] = liveStreamMonitorInfo.AiAsrInputIndexList
			}

			if liveStreamMonitorInfo.CheckStreamBroken != nil {
				liveStreamMonitorInfoMap["check_stream_broken"] = liveStreamMonitorInfo.CheckStreamBroken
			}

			if liveStreamMonitorInfo.CheckStreamLowFrameRate != nil {
				liveStreamMonitorInfoMap["check_stream_low_frame_rate"] = liveStreamMonitorInfo.CheckStreamLowFrameRate
			}

			if liveStreamMonitorInfo.AsrLanguage != nil {
				liveStreamMonitorInfoMap["asr_language"] = liveStreamMonitorInfo.AsrLanguage
			}

			if liveStreamMonitorInfo.OcrLanguage != nil {
				liveStreamMonitorInfoMap["ocr_language"] = liveStreamMonitorInfo.OcrLanguage
			}

			if liveStreamMonitorInfo.AiOcrInputIndexList != nil {
				liveStreamMonitorInfoMap["ai_ocr_input_index_list"] = liveStreamMonitorInfo.AiOcrInputIndexList
			}

			if liveStreamMonitorInfo.AllowMonitorReport != nil {
				liveStreamMonitorInfoMap["allow_monitor_report"] = liveStreamMonitorInfo.AllowMonitorReport
			}

			if liveStreamMonitorInfo.AiFormatDiagnose != nil {
				liveStreamMonitorInfoMap["ai_format_diagnose"] = liveStreamMonitorInfo.AiFormatDiagnose
			}

			ids = append(ids, *liveStreamMonitorInfo.MonitorId)
			tmpList = append(tmpList, liveStreamMonitorInfoMap)
		}

		_ = d.Set("live_stream_monitors", tmpList)
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
