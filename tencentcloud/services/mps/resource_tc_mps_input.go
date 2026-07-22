package mps

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMpsInput() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsInputCreate,
		Read:   resourceTencentCloudMpsInputRead,
		Update: resourceTencentCloudMpsInputUpdate,
		Delete: resourceTencentCloudMpsInputDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"flow_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Flow ID。",
			},

			"input_group": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "input 组 对于 input. Only support 一个 组 对于 一个 `tencentcloud_mps_input`. Use `for_each` 到 create 多个 inputs Scenario。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"input_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "input 名称，您 可以 fill 在 uppercase 和 lowercase letters，numbers 和 underscores，和 长度 是 [1，32]。",
						},
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Input 协议，可选 [SRT|RTP|RTMP|RTMP_PULL]。",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "input 描述 使用 长度 的 [0，255]。",
						},
						"allow_ip_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Computed:    true,
							Description: "input IP whitelist， 格式 是 CIDR。",
						},
						"srt_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "input SRT 配置 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "SRT 模式，可选 [LISTENER|CALLER]，默认为 LISTENER。",
									},
									"stream_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Stream ID，可选 uppercase 和 lowercase letters，numbers 和 special 字符 (.#!:&amp;,=_-)，长度 0~512. Specific 格式 可以 refer 到:https://github.com/Haivision/srt/blob/master/docs/features/访问-control.md#standard-keys。",
									},
									"latency": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "延迟，默认值 0，单位 ms，范围 [0，3000]。",
									},
									"recv_latency": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Receiving 延迟，默认为 120，单位 ms，范围 是 [0，3000]。",
									},
									"peer_latency": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Peer 延迟， 默认为 0， 单位 是 ms，和 范围 是 [0，3000]。",
									},
									"peer_idle_timeout": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Peer 超时，默认为 5000，单位 ms，范围 是 [1000，10000]。",
									},
									"passphrase": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "decryption 键，其中 是 空 通过 默认值，表示 无 加密. Only ascii 代码 值 可以 是 filled 在，和 长度 是 [10，79]。",
									},
									"pb_key_len": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "键 长度，默认为 0，可选 [0|16|24|32]。",
									},
									"source_addresses": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "SRT peer 地址，必填 当 模式 是 CALLER，和 仅 1 集合 可以 是 filled 在。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Peer IP。",
												},
												"port": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Peer 端口",
												},
											},
										},
									},
								},
							},
						},
						"rtp_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Input RTP 配置 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fec": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "默认为 &#39;none&#39;，可选 值[&#39;none&#39;]。",
									},
									"idle_timeout": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Idle 超时， 默认为 5000， 单位 是 ms，和 范围 是 [1000，10000]。",
									},
								},
							},
						},
						"fail_over": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "活跃/standby switch 的 input，[OPEN|CLOSE] 为可选项，和 默认为 CLOSE。",
						},
						"rtmp_pull_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Input RTMP_PULL 配置 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_addresses": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "来源 site 地址 的 RTMP 来源 site，there 可以 仅 是 一个。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"tc_url": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "TcUrl 地址 的 RTMP 来源 服务器。",
												},
												"stream_key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "StreamKey 信息 的 RTMP 来源 site。",
												},
											},
										},
									},
								},
							},
						},
						"rtsp_pull_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Input RTSP_PULL 配置 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_addresses": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "来源 site 地址 的 RTSP 来源 site，there 可以 仅 是 一个。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "URL 地址 的 RTSP 来源 site。",
												},
											},
										},
									},
								},
							},
						},
						"hls_pull_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Input HLS_PULL 配置 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_addresses": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "There 是 仅 一个 源站 地址 的 HLS 源站 station。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "URL 地址 的 HLS 源站 site。",
												},
											},
										},
									},
								},
							},
						},
						"resilient_stream": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "延迟 broadcast smooth streaming 配置 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "是否enable delayed broadcast smooth spit 流，true 是 已启用，false 是 不 已启用，和 默认为 不 已启用 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"buffer_time": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "延迟 时间，（秒）， currently 支持 范围 的 10 到 300 秒. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMpsInputCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_input.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = mps.NewCreateStreamLinkInputRequest()
		response = mps.NewCreateStreamLinkInputResponse()
		inputId  string
		flowId   string
		protocol string
	)
	if v, ok := d.GetOk("flow_id"); ok {
		request.FlowId = helper.String(v.(string))
		flowId = v.(string)
	}

	if v, ok := d.GetOk("input_group"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			createInput := mps.CreateInput{}
			if v, ok := dMap["input_name"]; ok {
				createInput.InputName = helper.String(v.(string))
			}
			if v, ok := dMap["protocol"]; ok {
				createInput.Protocol = helper.String(v.(string))
				protocol = v.(string)
			}
			if v, ok := dMap["description"]; ok {
				createInput.Description = helper.String(v.(string))
			}
			if v, ok := dMap["allow_ip_list"]; ok {
				allowIpListSet := v.(*schema.Set).List()
				for i := range allowIpListSet {
					if allowIpListSet[i] != nil {
						allowIpList := allowIpListSet[i].(string)
						createInput.AllowIpList = append(createInput.AllowIpList, &allowIpList)
					}
				}
			}
			if protocol == PROTOCOL_SRT {
				if sRTSettingsMap, ok := helper.InterfaceToMap(dMap, "srt_settings"); ok {
					createInputSRTSettings := mps.CreateInputSRTSettings{}
					if v, ok := sRTSettingsMap["mode"]; ok {
						createInputSRTSettings.Mode = helper.String(v.(string))
					}
					if v, ok := sRTSettingsMap["stream_id"]; ok {
						createInputSRTSettings.StreamId = helper.String(v.(string))
					}
					if v, ok := sRTSettingsMap["latency"]; ok {
						createInputSRTSettings.Latency = helper.IntInt64(v.(int))
					}
					if v, ok := sRTSettingsMap["recv_latency"]; ok {
						createInputSRTSettings.RecvLatency = helper.IntInt64(v.(int))
					}
					if v, ok := sRTSettingsMap["peer_latency"]; ok {
						createInputSRTSettings.PeerLatency = helper.IntInt64(v.(int))
					}
					if v, ok := sRTSettingsMap["peer_idle_timeout"]; ok {
						createInputSRTSettings.PeerIdleTimeout = helper.IntInt64(v.(int))
					}
					if v, ok := sRTSettingsMap["passphrase"]; ok {
						createInputSRTSettings.Passphrase = helper.String(v.(string))
					}
					if v, ok := sRTSettingsMap["pb_key_len"]; ok {
						createInputSRTSettings.PbKeyLen = helper.IntInt64(v.(int))
					}
					if v, ok := sRTSettingsMap["source_addresses"]; ok {
						for _, item := range v.([]interface{}) {
							sourceAddressesMap := item.(map[string]interface{})
							sRTSourceAddressReq := mps.SRTSourceAddressReq{}
							if v, ok := sourceAddressesMap["ip"]; ok {
								sRTSourceAddressReq.Ip = helper.String(v.(string))
							}
							if v, ok := sourceAddressesMap["port"]; ok {
								sRTSourceAddressReq.Port = helper.IntInt64(v.(int))
							}
							createInputSRTSettings.SourceAddresses = append(createInputSRTSettings.SourceAddresses, &sRTSourceAddressReq)
						}
					}
					createInput.SRTSettings = &createInputSRTSettings
				}
			}
			if protocol == PROTOCOL_RTP {
				if rTPSettingsMap, ok := helper.InterfaceToMap(dMap, "rtp_settings"); ok {
					createInputRTPSettings := mps.CreateInputRTPSettings{}
					if v, ok := rTPSettingsMap["fec"]; ok {
						createInputRTPSettings.FEC = helper.String(v.(string))
					}
					if v, ok := rTPSettingsMap["idle_timeout"]; ok {
						createInputRTPSettings.IdleTimeout = helper.IntInt64(v.(int))
					}
					createInput.RTPSettings = &createInputRTPSettings
				}
			}
			if v, ok := dMap["fail_over"]; ok {
				createInput.FailOver = helper.String(v.(string))
			}
			if protocol == PROTOCOL_RTMP || protocol == PROTOCOL_RTMP_PULL {
				if rTMPPullSettingsMap, ok := helper.InterfaceToMap(dMap, "rtmp_pull_settings"); ok {
					createInputRTMPPullSettings := mps.CreateInputRTMPPullSettings{}
					if v, ok := rTMPPullSettingsMap["source_addresses"]; ok {
						for _, item := range v.([]interface{}) {
							sourceAddressesMap := item.(map[string]interface{})
							rTMPPullSourceAddress := mps.RTMPPullSourceAddress{}
							if v, ok := sourceAddressesMap["tc_url"]; ok {
								rTMPPullSourceAddress.TcUrl = helper.String(v.(string))
							}
							if v, ok := sourceAddressesMap["stream_key"]; ok {
								rTMPPullSourceAddress.StreamKey = helper.String(v.(string))
							}
							createInputRTMPPullSettings.SourceAddresses = append(createInputRTMPPullSettings.SourceAddresses, &rTMPPullSourceAddress)
						}
					}
					createInput.RTMPPullSettings = &createInputRTMPPullSettings
				}
			}
			if protocol == PROTOCOL_RTSP_PULL {
				if rTSPPullSettingsMap, ok := helper.InterfaceToMap(dMap, "rtsp_pull_settings"); ok {
					createInputRTSPPullSettings := mps.CreateInputRTSPPullSettings{}
					if v, ok := rTSPPullSettingsMap["source_addresses"]; ok {
						for _, item := range v.([]interface{}) {
							sourceAddressesMap := item.(map[string]interface{})
							rTSPPullSourceAddress := mps.RTSPPullSourceAddress{}
							if v, ok := sourceAddressesMap["url"]; ok {
								rTSPPullSourceAddress.Url = helper.String(v.(string))
							}
							createInputRTSPPullSettings.SourceAddresses = append(createInputRTSPPullSettings.SourceAddresses, &rTSPPullSourceAddress)
						}
					}
					createInput.RTSPPullSettings = &createInputRTSPPullSettings
				}
			}
			if protocol == PROTOCOL_HLS || protocol == PROTOCOL_HLS_PULL {
				if hLSPullSettingsMap, ok := helper.InterfaceToMap(dMap, "hls_pull_settings"); ok {
					createInputHLSPullSettings := mps.CreateInputHLSPullSettings{}
					if v, ok := hLSPullSettingsMap["source_addresses"]; ok {
						for _, item := range v.([]interface{}) {
							sourceAddressesMap := item.(map[string]interface{})
							hLSPullSourceAddress := mps.HLSPullSourceAddress{}
							if v, ok := sourceAddressesMap["url"]; ok {
								hLSPullSourceAddress.Url = helper.String(v.(string))
							}
							createInputHLSPullSettings.SourceAddresses = append(createInputHLSPullSettings.SourceAddresses, &hLSPullSourceAddress)
						}
					}
					createInput.HLSPullSettings = &createInputHLSPullSettings
				}
			}
			if resilientStreamMap, ok := helper.InterfaceToMap(dMap, "resilient_stream"); ok {
				resilientStreamConf := mps.ResilientStreamConf{}
				if v, ok := resilientStreamMap["enable"]; ok {
					resilientStreamConf.Enable = helper.Bool(v.(bool))
				}
				if v, ok := resilientStreamMap["buffer_time"]; ok {
					resilientStreamConf.BufferTime = helper.IntUint64(v.(int))
				}
				createInput.ResilientStream = &resilientStreamConf
			}
			request.InputGroup = append(request.InputGroup, &createInput)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().CreateStreamLinkInput(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps input failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.Info != nil && len(response.Response.Info.InputGroup) > 0 {
		inputId = *response.Response.Info.InputGroup[0].InputId
	}

	d.SetId(strings.Join([]string{flowId, inputId}, tccommon.FILED_SP))

	return resourceTencentCloudMpsInputRead(d, meta)
}

func resourceTencentCloudMpsInputRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_input.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var protocol string

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	flowId := idSplit[0]
	inputId := idSplit[1]

	input, err := service.DescribeMpsInputById(ctx, flowId, inputId)
	if err != nil {
		return err
	}

	if input == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource Mps input group [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("flow_id", flowId)

	inputGroupMap := map[string]interface{}{}

	if input.InputName != nil {
		inputGroupMap["input_name"] = input.InputName
	}

	if input.Protocol != nil {
		inputGroupMap["protocol"] = input.Protocol
		protocol = *input.Protocol
	}

	if input.Description != nil {
		inputGroupMap["description"] = input.Description
	}

	if input.AllowIpList != nil {
		inputGroupMap["allow_ip_list"] = input.AllowIpList
	}

	if protocol == PROTOCOL_SRT && input.SRTSettings != nil {
		sRTSettingsMap := map[string]interface{}{}

		if input.SRTSettings.Mode != nil {
			sRTSettingsMap["mode"] = input.SRTSettings.Mode
		}

		if input.SRTSettings.StreamId != nil {
			sRTSettingsMap["stream_id"] = input.SRTSettings.StreamId
		}

		if input.SRTSettings.Latency != nil {
			sRTSettingsMap["latency"] = input.SRTSettings.Latency
		}

		if input.SRTSettings.RecvLatency != nil {
			sRTSettingsMap["recv_latency"] = input.SRTSettings.RecvLatency
		}

		// if input.SRTSettings.PeerLatency != nil {
		// 	sRTSettingsMap["peer_latency"] = input.SRTSettings.PeerLatency
		// }
		//cannot be imported
		if input.SRTSettings.PeerLatency != nil {
			index := fmt.Sprintf("input_group.%d.srt_settings.0.peer_latency", 0)
			oldValue := d.Get(index).(int)
			if *input.SRTSettings.PeerLatency == 0 {
				// need fix: the SDK has bug that cannot return the real value for peer_latency.
				sRTSettingsMap["peer_latency"] = helper.IntInt64(oldValue)
			} else {
				sRTSettingsMap["peer_latency"] = input.SRTSettings.PeerLatency
			}
		}

		if input.SRTSettings.PeerIdleTimeout != nil {
			sRTSettingsMap["peer_idle_timeout"] = input.SRTSettings.PeerIdleTimeout
		}

		if input.SRTSettings.Passphrase != nil {
			sRTSettingsMap["passphrase"] = input.SRTSettings.Passphrase
		}

		if input.SRTSettings.PbKeyLen != nil {
			sRTSettingsMap["pb_key_len"] = input.SRTSettings.PbKeyLen
		}

		if input.SRTSettings.SourceAddresses != nil {
			sourceAddressesList := []interface{}{}
			for _, sourceAddresses := range input.SRTSettings.SourceAddresses {
				sourceAddressesMap := map[string]interface{}{}

				if sourceAddresses.Ip != nil {
					sourceAddressesMap["ip"] = sourceAddresses.Ip
				}

				if sourceAddresses.Port != nil {
					sourceAddressesMap["port"] = sourceAddresses.Port
				}

				sourceAddressesList = append(sourceAddressesList, sourceAddressesMap)
			}

			sRTSettingsMap["source_addresses"] = sourceAddressesList
		}

		inputGroupMap["srt_settings"] = []interface{}{sRTSettingsMap}
	}

	if protocol == PROTOCOL_RTP && input.RTPSettings != nil {
		rTPSettingsMap := map[string]interface{}{}

		if input.RTPSettings.FEC != nil {
			rTPSettingsMap["fec"] = input.RTPSettings.FEC
		}

		if input.RTPSettings.IdleTimeout != nil {
			rTPSettingsMap["idle_timeout"] = input.RTPSettings.IdleTimeout
		}

		inputGroupMap["rtp_settings"] = []interface{}{rTPSettingsMap}
	}

	if input.FailOver != nil {
		inputGroupMap["fail_over"] = input.FailOver
	}

	if (protocol == PROTOCOL_RTMP || protocol == PROTOCOL_RTMP_PULL) && input.RTMPPullSettings != nil {
		rTMPPullSettingsMap := map[string]interface{}{}

		if input.RTMPPullSettings.SourceAddresses != nil {
			sourceAddressesList := []interface{}{}
			for _, sourceAddresses := range input.RTMPPullSettings.SourceAddresses {
				sourceAddressesMap := map[string]interface{}{}

				if sourceAddresses.TcUrl != nil {
					sourceAddressesMap["tc_url"] = sourceAddresses.TcUrl
				}

				if sourceAddresses.StreamKey != nil {
					sourceAddressesMap["stream_key"] = sourceAddresses.StreamKey
				}

				sourceAddressesList = append(sourceAddressesList, sourceAddressesMap)
			}

			rTMPPullSettingsMap["source_addresses"] = sourceAddressesList
		}

		inputGroupMap["rtmp_pull_settings"] = []interface{}{rTMPPullSettingsMap}
	}

	if protocol == PROTOCOL_RTSP_PULL && input.RTSPPullSettings != nil {
		rTSPPullSettingsMap := map[string]interface{}{}

		if input.RTSPPullSettings.SourceAddresses != nil {
			sourceAddressesList := []interface{}{}
			for _, sourceAddresses := range input.RTSPPullSettings.SourceAddresses {
				sourceAddressesMap := map[string]interface{}{}

				if sourceAddresses.Url != nil {
					sourceAddressesMap["url"] = sourceAddresses.Url
				}

				sourceAddressesList = append(sourceAddressesList, sourceAddressesMap)
			}

			rTSPPullSettingsMap["source_addresses"] = sourceAddressesList
		}

		inputGroupMap["rtsp_pull_settings"] = []interface{}{rTSPPullSettingsMap}
	}

	if (protocol == PROTOCOL_HLS || protocol == PROTOCOL_HLS_PULL) && input.HLSPullSettings != nil {
		hLSPullSettingsMap := map[string]interface{}{}

		if input.HLSPullSettings.SourceAddresses != nil {
			sourceAddressesList := []interface{}{}
			for _, sourceAddresses := range input.HLSPullSettings.SourceAddresses {
				sourceAddressesMap := map[string]interface{}{}

				if sourceAddresses.Url != nil {
					sourceAddressesMap["url"] = sourceAddresses.Url
				}

				sourceAddressesList = append(sourceAddressesList, sourceAddressesMap)
			}

			hLSPullSettingsMap["source_addresses"] = sourceAddressesList
		}

		inputGroupMap["hls_pull_settings"] = []interface{}{hLSPullSettingsMap}
	}

	if input.ResilientStream != nil {
		resilientStreamMap := map[string]interface{}{}

		if input.ResilientStream.Enable != nil {
			resilientStreamMap["enable"] = input.ResilientStream.Enable
		}

		if input.ResilientStream.BufferTime != nil {
			resilientStreamMap["buffer_time"] = input.ResilientStream.BufferTime
		}

		inputGroupMap["resilient_stream"] = []interface{}{resilientStreamMap}
	}

	_ = d.Set("input_group", []interface{}{inputGroupMap})

	return nil
}

func resourceTencentCloudMpsInputUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_input.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := mps.NewModifyStreamLinkInputRequest()
	var protocol string

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	flowId := idSplit[0]
	inputId := idSplit[1]

	request.FlowId = &flowId

	immutableArgs := []string{"flow_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("input_group") {
		if v, ok := d.GetOk("input_group"); ok {
			for _, item := range v.([]interface{}) {
				modifyInput := mps.ModifyInput{}
				modifyInput.InputId = &inputId
				dMap := item.(map[string]interface{})
				if v, ok := dMap["input_name"]; ok {
					modifyInput.InputName = helper.String(v.(string))
				}
				if v, ok := dMap["protocol"]; ok {
					modifyInput.Protocol = helper.String(v.(string))
					protocol = v.(string)
				}
				if v, ok := dMap["description"]; ok {
					modifyInput.Description = helper.String(v.(string))
				}
				if v, ok := dMap["allow_ip_list"]; ok {
					allowIpListSet := v.(*schema.Set).List()
					for i := range allowIpListSet {
						if allowIpListSet[i] != nil {
							allowIpList := allowIpListSet[i].(string)
							modifyInput.AllowIpList = append(modifyInput.AllowIpList, &allowIpList)
						}
					}
				}
				if protocol == PROTOCOL_SRT {
					if sRTSettingsMap, ok := helper.InterfaceToMap(dMap, "srt_settings"); ok {
						createInputSRTSettings := mps.CreateInputSRTSettings{}
						if v, ok := sRTSettingsMap["mode"]; ok {
							createInputSRTSettings.Mode = helper.String(v.(string))
						}
						if v, ok := sRTSettingsMap["stream_id"]; ok {
							createInputSRTSettings.StreamId = helper.String(v.(string))
						}
						if v, ok := sRTSettingsMap["latency"]; ok {
							createInputSRTSettings.Latency = helper.IntInt64(v.(int))
						}
						if v, ok := sRTSettingsMap["recv_latency"]; ok {
							createInputSRTSettings.RecvLatency = helper.IntInt64(v.(int))
						}
						if v, ok := sRTSettingsMap["peer_latency"]; ok {
							createInputSRTSettings.PeerLatency = helper.IntInt64(v.(int))
						}
						if v, ok := sRTSettingsMap["peer_idle_timeout"]; ok {
							createInputSRTSettings.PeerIdleTimeout = helper.IntInt64(v.(int))
						}
						if v, ok := sRTSettingsMap["passphrase"]; ok {
							createInputSRTSettings.Passphrase = helper.String(v.(string))
						}
						if v, ok := sRTSettingsMap["pb_key_len"]; ok {
							createInputSRTSettings.PbKeyLen = helper.IntInt64(v.(int))
						}
						if v, ok := sRTSettingsMap["source_addresses"]; ok {
							for _, item := range v.([]interface{}) {
								sourceAddressesMap := item.(map[string]interface{})
								sRTSourceAddressReq := mps.SRTSourceAddressReq{}
								if v, ok := sourceAddressesMap["ip"]; ok {
									sRTSourceAddressReq.Ip = helper.String(v.(string))
								}
								if v, ok := sourceAddressesMap["port"]; ok {
									sRTSourceAddressReq.Port = helper.IntInt64(v.(int))
								}
								createInputSRTSettings.SourceAddresses = append(createInputSRTSettings.SourceAddresses, &sRTSourceAddressReq)
							}
						}
						modifyInput.SRTSettings = &createInputSRTSettings
					}
				}
				if protocol == PROTOCOL_RTP {
					if rTPSettingsMap, ok := helper.InterfaceToMap(dMap, "rtp_settings"); ok {
						createInputRTPSettings := mps.CreateInputRTPSettings{}
						if v, ok := rTPSettingsMap["fec"]; ok {
							createInputRTPSettings.FEC = helper.String(v.(string))
						}
						if v, ok := rTPSettingsMap["idle_timeout"]; ok {
							createInputRTPSettings.IdleTimeout = helper.IntInt64(v.(int))
						}
						modifyInput.RTPSettings = &createInputRTPSettings
					}
				}
				if v, ok := dMap["fail_over"]; ok {
					modifyInput.FailOver = helper.String(v.(string))
				}
				if protocol == PROTOCOL_RTMP_PULL {
					if rTMPPullSettingsMap, ok := helper.InterfaceToMap(dMap, "rtmp_pull_settings"); ok {
						createInputRTMPPullSettings := mps.CreateInputRTMPPullSettings{}
						if v, ok := rTMPPullSettingsMap["source_addresses"]; ok {
							for _, item := range v.([]interface{}) {
								sourceAddressesMap := item.(map[string]interface{})
								rTMPPullSourceAddress := mps.RTMPPullSourceAddress{}
								if v, ok := sourceAddressesMap["tc_url"]; ok {
									rTMPPullSourceAddress.TcUrl = helper.String(v.(string))
								}
								if v, ok := sourceAddressesMap["stream_key"]; ok {
									rTMPPullSourceAddress.StreamKey = helper.String(v.(string))
								}
								createInputRTMPPullSettings.SourceAddresses = append(createInputRTMPPullSettings.SourceAddresses, &rTMPPullSourceAddress)
							}
						}
						modifyInput.RTMPPullSettings = &createInputRTMPPullSettings
					}
				}
				if protocol == PROTOCOL_RTSP_PULL {
					if rTSPPullSettingsMap, ok := helper.InterfaceToMap(dMap, "rtsp_pull_settings"); ok {
						createInputRTSPPullSettings := mps.CreateInputRTSPPullSettings{}
						if v, ok := rTSPPullSettingsMap["source_addresses"]; ok {
							for _, item := range v.([]interface{}) {
								sourceAddressesMap := item.(map[string]interface{})
								rTSPPullSourceAddress := mps.RTSPPullSourceAddress{}
								if v, ok := sourceAddressesMap["url"]; ok {
									rTSPPullSourceAddress.Url = helper.String(v.(string))
								}
								createInputRTSPPullSettings.SourceAddresses = append(createInputRTSPPullSettings.SourceAddresses, &rTSPPullSourceAddress)
							}
						}
						modifyInput.RTSPPullSettings = &createInputRTSPPullSettings
					}
				}
				if hLSPullSettingsMap, ok := helper.InterfaceToMap(dMap, "hls_pull_settings"); ok {
					createInputHLSPullSettings := mps.CreateInputHLSPullSettings{}
					if v, ok := hLSPullSettingsMap["source_addresses"]; ok {
						for _, item := range v.([]interface{}) {
							sourceAddressesMap := item.(map[string]interface{})
							hLSPullSourceAddress := mps.HLSPullSourceAddress{}
							if v, ok := sourceAddressesMap["url"]; ok {
								hLSPullSourceAddress.Url = helper.String(v.(string))
							}
							createInputHLSPullSettings.SourceAddresses = append(createInputHLSPullSettings.SourceAddresses, &hLSPullSourceAddress)
						}
					}
					modifyInput.HLSPullSettings = &createInputHLSPullSettings
				}
				if resilientStreamMap, ok := helper.InterfaceToMap(dMap, "resilient_stream"); ok {
					resilientStreamConf := mps.ResilientStreamConf{}
					if v, ok := resilientStreamMap["enable"]; ok {
						resilientStreamConf.Enable = helper.Bool(v.(bool))
					}
					if v, ok := resilientStreamMap["buffer_time"]; ok {
						resilientStreamConf.BufferTime = helper.IntUint64(v.(int))
					}
					modifyInput.ResilientStream = &resilientStreamConf
				}
				//  modify api only support to modify one input one time
				request.Input = &modifyInput
			}
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().ModifyStreamLinkInput(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mps input failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMpsInputRead(d, meta)
}

func resourceTencentCloudMpsInputDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_input.delete")()
	defer tccommon.InconsistentCheck(d, meta)()
	// deleted through `tencentcloud_mps_flow`
	return nil
}
