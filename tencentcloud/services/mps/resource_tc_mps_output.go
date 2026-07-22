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

func ResourceTencentCloudMpsOutput() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsOutputCreate,
		Read:   resourceTencentCloudMpsOutputRead,
		Update: resourceTencentCloudMpsOutputUpdate,
		Delete: resourceTencentCloudMpsOutputDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"flow_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Flow ID。",
			},

			"output": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Output 配置 的 transport 流。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"output_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "名称 output。",
						},
						"description": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Output 描述",
						},
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Output 协议，可选 [SRT|RTP|RTMP|RTMP_PULL]。",
						},
						"output_region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "输出地域",
						},
						"srt_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "配置 的 output SRT。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destinations": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "目标 地址 的 relay 为必填项 当 模式 是 CALLER，和 仅 一个 组 可以 是 filled 在。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Output IP。",
												},
												"port": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "output 端口",
												},
											},
										},
									},
									"stream_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "relay 流 ID SRT. You 可以 choose uppercase 和 lowercase letters，numbers 和 special 字符 (.#!:&amp;,=_-). 长度 是 0~512。",
									},
									"latency": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "总数 延迟 的 relaying SRT， 默认为 0， 单位 是 ms，和 范围 是 [0，3000]。",
									},
									"recv_latency": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "reception 延迟 的 relay SRT， 默认为 120， 单位 是 ms， 范围 是 [0，3000]。",
									},
									"peer_latency": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "peer 延迟 的 relaying SRT， 默认为 0， 单位 是 ms，和 范围 是 [0，3000]。",
									},
									"peer_idle_timeout": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "peer idle 超时 对于 relaying SRT， 默认为 5000， 单位 是 ms，和 范围 是 [1000，10000]。",
									},
									"passphrase": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "加密 键 对于 relaying SRT，其中 是 空 通过 默认值，indicating 无 加密. Only ascii 代码 值 可以 是 filled 在，和 长度 是 [10，79]。",
									},
									"pb_key_len": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "键 长度 的 relay SRT， 默认为 0，可选 [0|16|24|32]。",
									},
									"mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "SRT 模式，可选 [LISTENER|CALLER]，默认为 CALLER。",
									},
								},
							},
						},
						"rtmp_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Output RTMP 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destinations": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "目标 地址 的 relay 可以 是 filled 在 1~2。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "relayed URL， 格式 是: rtmp://域名/live。",
												},
												"stream_key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "relayed StreamKey，在 格式: 流?键=值",
												},
											},
										},
									},
									"chunk_size": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "RTMP Chunk 大小，范围 是 [4096，40960]。",
									},
								},
							},
						},
						"rtp_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Output RTP 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destinations": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "目标 地址 的 relay 可以 是 filled 在 1~2。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "目标 IP 的 relay。",
												},
												"port": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Destination 端口 对于 relays。",
												},
											},
										},
									},
									"fec": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "You 可以 仅 fill 在 none。",
									},
									"idle_timeout": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Idle 超时，单位 ms。",
									},
								},
							},
						},
						"allow_ip_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Computed:    true,
							Description: "IP whitelist 列表， 格式 是 CIDR，such 作为 0.0.0.0/0. 当 协议 是 RTMP_PULL，它 是 有效，和 如果 它 是 空，它 表示 该 客户端 IP 是 不 limited。",
						},
						"max_concurrent": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "最大concurrent pull streams， 最大 是 4，和 默认为 4. Only SRT 或 RTMP_PULL 可以 集合 此 参数。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMpsOutputCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_output.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = mps.NewCreateStreamLinkOutputInfoRequest()
		response = mps.NewCreateStreamLinkOutputInfoResponse()
		outputId string
		flowId   string
		protocol string
	)
	if v, ok := d.GetOk("flow_id"); ok {
		request.FlowId = helper.String(v.(string))
		flowId = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "output"); ok {
		createOutputInfo := mps.CreateOutputInfo{}
		if v, ok := dMap["output_name"]; ok {
			createOutputInfo.OutputName = helper.String(v.(string))
		}
		if v, ok := dMap["description"]; ok {
			createOutputInfo.Description = helper.String(v.(string))
		}
		if v, ok := dMap["protocol"]; ok {
			createOutputInfo.Protocol = helper.String(v.(string))
			protocol = v.(string)
		}
		if v, ok := dMap["output_region"]; ok {
			createOutputInfo.OutputRegion = helper.String(v.(string))
		}
		if protocol == PROTOCOL_SRT {
			if sRTSettingsMap, ok := helper.InterfaceToMap(dMap, "srt_settings"); ok {
				createOutputSRTSettings := mps.CreateOutputSRTSettings{}
				if v, ok := sRTSettingsMap["destinations"]; ok {
					for _, item := range v.([]interface{}) {
						destinationsMap := item.(map[string]interface{})
						createOutputSRTSettingsDestinations := mps.CreateOutputSRTSettingsDestinations{}
						if v, ok := destinationsMap["ip"]; ok {
							createOutputSRTSettingsDestinations.Ip = helper.String(v.(string))
						}
						if v, ok := destinationsMap["port"]; ok {
							createOutputSRTSettingsDestinations.Port = helper.IntInt64(v.(int))
						}
						createOutputSRTSettings.Destinations = append(createOutputSRTSettings.Destinations, &createOutputSRTSettingsDestinations)
					}
				}
				if v, ok := sRTSettingsMap["stream_id"]; ok {
					createOutputSRTSettings.StreamId = helper.String(v.(string))
				}
				if v, ok := sRTSettingsMap["latency"]; ok {
					createOutputSRTSettings.Latency = helper.IntInt64(v.(int))
				}
				if v, ok := sRTSettingsMap["recv_latency"]; ok {
					createOutputSRTSettings.RecvLatency = helper.IntInt64(v.(int))
				}
				if v, ok := sRTSettingsMap["peer_latency"]; ok {
					createOutputSRTSettings.PeerLatency = helper.IntInt64(v.(int))
				}
				if v, ok := sRTSettingsMap["peer_idle_timeout"]; ok {
					createOutputSRTSettings.PeerIdleTimeout = helper.IntInt64(v.(int))
				}
				if v, ok := sRTSettingsMap["passphrase"]; ok {
					createOutputSRTSettings.Passphrase = helper.String(v.(string))
				}
				if v, ok := sRTSettingsMap["pb_key_len"]; ok {
					createOutputSRTSettings.PbKeyLen = helper.IntInt64(v.(int))
				}
				if v, ok := sRTSettingsMap["mode"]; ok {
					createOutputSRTSettings.Mode = helper.String(v.(string))
				}
				createOutputInfo.SRTSettings = &createOutputSRTSettings
			}
		}
		if protocol == PROTOCOL_RTMP || protocol == PROTOCOL_RTMP_PULL {
			if rTMPSettingsMap, ok := helper.InterfaceToMap(dMap, "rtmp_settings"); ok {
				createOutputRTMPSettings := mps.CreateOutputRTMPSettings{}
				if v, ok := rTMPSettingsMap["destinations"]; ok {
					for _, item := range v.([]interface{}) {
						destinationsMap := item.(map[string]interface{})
						createOutputRtmpSettingsDestinations := mps.CreateOutputRtmpSettingsDestinations{}
						if v, ok := destinationsMap["url"]; ok {
							createOutputRtmpSettingsDestinations.Url = helper.String(v.(string))
						}
						if v, ok := destinationsMap["stream_key"]; ok {
							createOutputRtmpSettingsDestinations.StreamKey = helper.String(v.(string))
						}
						createOutputRTMPSettings.Destinations = append(createOutputRTMPSettings.Destinations, &createOutputRtmpSettingsDestinations)
					}
				}
				if v, ok := rTMPSettingsMap["chunk_size"]; ok {
					createOutputRTMPSettings.ChunkSize = helper.IntInt64(v.(int))
				}
				createOutputInfo.RTMPSettings = &createOutputRTMPSettings
			}
		}
		if protocol == PROTOCOL_RTP {
			if rTPSettingsMap, ok := helper.InterfaceToMap(dMap, "rtp_settings"); ok {
				createOutputInfoRTPSettings := mps.CreateOutputInfoRTPSettings{}
				if v, ok := rTPSettingsMap["destinations"]; ok {
					for _, item := range v.([]interface{}) {
						destinationsMap := item.(map[string]interface{})
						createOutputRTPSettingsDestinations := mps.CreateOutputRTPSettingsDestinations{}
						if v, ok := destinationsMap["ip"]; ok {
							createOutputRTPSettingsDestinations.Ip = helper.String(v.(string))
						}
						if v, ok := destinationsMap["port"]; ok {
							createOutputRTPSettingsDestinations.Port = helper.IntInt64(v.(int))
						}
						createOutputInfoRTPSettings.Destinations = append(createOutputInfoRTPSettings.Destinations, &createOutputRTPSettingsDestinations)
					}
				}
				if v, ok := rTPSettingsMap["fec"]; ok {
					createOutputInfoRTPSettings.FEC = helper.String(v.(string))
				}
				if v, ok := rTPSettingsMap["idle_timeout"]; ok {
					createOutputInfoRTPSettings.IdleTimeout = helper.IntInt64(v.(int))
				}
				createOutputInfo.RTPSettings = &createOutputInfoRTPSettings
			}
		}
		if v, ok := dMap["allow_ip_list"]; ok {
			allowIpListSet := v.(*schema.Set).List()
			for i := range allowIpListSet {
				if allowIpListSet[i] != nil {
					allowIpList := allowIpListSet[i].(string)
					createOutputInfo.AllowIpList = append(createOutputInfo.AllowIpList, &allowIpList)
				}
			}
		}
		if protocol == PROTOCOL_SRT || protocol == PROTOCOL_RTMP_PULL {
			if v, ok := dMap["max_concurrent"]; ok {
				createOutputInfo.MaxConcurrent = helper.IntUint64(v.(int))
			}
		}
		request.Output = &createOutputInfo
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().CreateStreamLinkOutputInfo(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps output failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.Info != nil {
		outputId = *response.Response.Info.OutputId
	}

	d.SetId(strings.Join([]string{flowId, outputId}, tccommon.FILED_SP))

	return resourceTencentCloudMpsOutputRead(d, meta)
}

func resourceTencentCloudMpsOutputRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_output.read")()
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
	outputId := idSplit[1]

	output, err := service.DescribeMpsOutputById(ctx, flowId, outputId)
	if err != nil {
		return err
	}

	if output == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource Mps output group [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("flow_id", flowId)

	outputMap := map[string]interface{}{}

	if output.OutputName != nil {
		outputMap["output_name"] = output.OutputName
	}

	if output.Description != nil {
		outputMap["description"] = output.Description
	}

	if output.Protocol != nil {
		outputMap["protocol"] = output.Protocol
		protocol = *output.Protocol
	}

	if output.OutputRegion != nil {
		outputMap["output_region"] = output.OutputRegion
	}

	if protocol == PROTOCOL_SRT && output.SRTSettings != nil {
		sRTSettingsMap := map[string]interface{}{}

		if output.SRTSettings.Destinations != nil {
			destinationsList := []interface{}{}
			for _, destinations := range output.SRTSettings.Destinations {
				destinationsMap := map[string]interface{}{}

				if destinations.Ip != nil {
					destinationsMap["ip"] = destinations.Ip
				}

				if destinations.Port != nil {
					destinationsMap["port"] = destinations.Port
				}

				destinationsList = append(destinationsList, destinationsMap)
			}

			sRTSettingsMap["destinations"] = destinationsList
		}

		if output.SRTSettings.StreamId != nil {
			sRTSettingsMap["stream_id"] = output.SRTSettings.StreamId
		}

		if output.SRTSettings.Latency != nil {
			sRTSettingsMap["latency"] = output.SRTSettings.Latency
		}

		if output.SRTSettings.RecvLatency != nil {
			sRTSettingsMap["recv_latency"] = output.SRTSettings.RecvLatency
		}

		if output.SRTSettings.PeerLatency != nil {
			sRTSettingsMap["peer_latency"] = output.SRTSettings.PeerLatency
		}

		if output.SRTSettings.PeerIdleTimeout != nil {
			sRTSettingsMap["peer_idle_timeout"] = output.SRTSettings.PeerIdleTimeout
		}

		if output.SRTSettings.Passphrase != nil {
			sRTSettingsMap["passphrase"] = output.SRTSettings.Passphrase
		}

		if output.SRTSettings.PbKeyLen != nil {
			sRTSettingsMap["pb_key_len"] = output.SRTSettings.PbKeyLen
		}

		if output.SRTSettings.Mode != nil {
			sRTSettingsMap["mode"] = output.SRTSettings.Mode
		}

		outputMap["srt_settings"] = []interface{}{sRTSettingsMap}
	}

	if (protocol == PROTOCOL_RTMP || protocol == PROTOCOL_RTMP_PULL) && output.RTMPSettings != nil {
		rTMPSettingsMap := map[string]interface{}{}

		if output.RTMPSettings.Destinations != nil {
			destinationsList := []interface{}{}
			for _, destinations := range output.RTMPSettings.Destinations {
				destinationsMap := map[string]interface{}{}

				if destinations.Url != nil {
					destinationsMap["url"] = destinations.Url
				}

				if destinations.StreamKey != nil {
					destinationsMap["stream_key"] = destinations.StreamKey
				}

				destinationsList = append(destinationsList, destinationsMap)
			}

			rTMPSettingsMap["destinations"] = destinationsList
		}

		if output.RTMPSettings.ChunkSize != nil {
			rTMPSettingsMap["chunk_size"] = output.RTMPSettings.ChunkSize
		}

		outputMap["rtmp_settings"] = []interface{}{rTMPSettingsMap}
	}

	if protocol == PROTOCOL_RTP && output.RTPSettings != nil {
		rTPSettingsMap := map[string]interface{}{}

		if output.RTPSettings.Destinations != nil {
			destinationsList := []interface{}{}
			for _, destinations := range output.RTPSettings.Destinations {
				destinationsMap := map[string]interface{}{}

				if destinations.Ip != nil {
					destinationsMap["ip"] = destinations.Ip
				}

				if destinations.Port != nil {
					destinationsMap["port"] = destinations.Port
				}

				destinationsList = append(destinationsList, destinationsMap)
			}

			rTPSettingsMap["destinations"] = destinationsList
		}

		if output.RTPSettings.FEC != nil {
			rTPSettingsMap["fec"] = output.RTPSettings.FEC
		}

		if output.RTPSettings.IdleTimeout != nil {
			rTPSettingsMap["idle_timeout"] = output.RTPSettings.IdleTimeout
		}

		outputMap["rtp_settings"] = []interface{}{rTPSettingsMap}
	}

	if output.AllowIpList != nil {
		outputMap["allow_ip_list"] = output.AllowIpList
	}

	if output.MaxConcurrent != nil {
		outputMap["max_concurrent"] = output.MaxConcurrent
	}

	_ = d.Set("output", []interface{}{outputMap})

	return nil
}

func resourceTencentCloudMpsOutputUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_output.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := mps.NewModifyStreamLinkOutputInfoRequest()
	var protocol string

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	flowId := idSplit[0]
	outputId := idSplit[1]

	request.FlowId = &flowId

	immutableArgs := []string{"flow_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("output") {
		if dMap, ok := helper.InterfacesHeadMap(d, "output"); ok {
			modifyOutputInfo := mps.ModifyOutputInfo{}
			modifyOutputInfo.OutputId = &outputId
			if v, ok := dMap["output_name"]; ok {
				modifyOutputInfo.OutputName = helper.String(v.(string))
			}
			if v, ok := dMap["description"]; ok {
				modifyOutputInfo.Description = helper.String(v.(string))
			}
			if v, ok := dMap["protocol"]; ok {
				modifyOutputInfo.Protocol = helper.String(v.(string))
				protocol = v.(string)
			}
			if sRTSettingsMap, ok := helper.InterfaceToMap(dMap, "srt_settings"); ok {
				modifyOutputSRTSettings := mps.CreateOutputSRTSettings{}
				if v, ok := sRTSettingsMap["destinations"]; ok {
					for _, item := range v.([]interface{}) {
						destinationsMap := item.(map[string]interface{})
						modifyOutputSRTSettingsDestinations := mps.CreateOutputSRTSettingsDestinations{}
						if v, ok := destinationsMap["ip"]; ok {
							modifyOutputSRTSettingsDestinations.Ip = helper.String(v.(string))
						}
						if v, ok := destinationsMap["port"]; ok {
							modifyOutputSRTSettingsDestinations.Port = helper.IntInt64(v.(int))
						}
						modifyOutputSRTSettings.Destinations = append(modifyOutputSRTSettings.Destinations, &modifyOutputSRTSettingsDestinations)
					}
				}
				if v, ok := sRTSettingsMap["stream_id"]; ok {
					modifyOutputSRTSettings.StreamId = helper.String(v.(string))
				}
				if v, ok := sRTSettingsMap["latency"]; ok {
					modifyOutputSRTSettings.Latency = helper.IntInt64(v.(int))
				}
				if v, ok := sRTSettingsMap["recv_latency"]; ok {
					modifyOutputSRTSettings.RecvLatency = helper.IntInt64(v.(int))
				}
				if v, ok := sRTSettingsMap["peer_latency"]; ok {
					modifyOutputSRTSettings.PeerLatency = helper.IntInt64(v.(int))
				}
				if v, ok := sRTSettingsMap["peer_idle_timeout"]; ok {
					modifyOutputSRTSettings.PeerIdleTimeout = helper.IntInt64(v.(int))
				}
				if v, ok := sRTSettingsMap["passphrase"]; ok {
					modifyOutputSRTSettings.Passphrase = helper.String(v.(string))
				}
				if v, ok := sRTSettingsMap["pb_key_len"]; ok {
					modifyOutputSRTSettings.PbKeyLen = helper.IntInt64(v.(int))
				}
				if v, ok := sRTSettingsMap["mode"]; ok {
					modifyOutputSRTSettings.Mode = helper.String(v.(string))
				}
				modifyOutputInfo.SRTSettings = &modifyOutputSRTSettings
			}
			if rTMPSettingsMap, ok := helper.InterfaceToMap(dMap, "rtmp_settings"); ok {
				createOutputRTMPSettings := mps.CreateOutputRTMPSettings{}
				if v, ok := rTMPSettingsMap["destinations"]; ok {
					for _, item := range v.([]interface{}) {
						destinationsMap := item.(map[string]interface{})
						createOutputRtmpSettingsDestinations := mps.CreateOutputRtmpSettingsDestinations{}
						if v, ok := destinationsMap["url"]; ok {
							createOutputRtmpSettingsDestinations.Url = helper.String(v.(string))
						}
						if v, ok := destinationsMap["stream_key"]; ok {
							createOutputRtmpSettingsDestinations.StreamKey = helper.String(v.(string))
						}
						createOutputRTMPSettings.Destinations = append(createOutputRTMPSettings.Destinations, &createOutputRtmpSettingsDestinations)
					}
				}
				if v, ok := rTMPSettingsMap["chunk_size"]; ok {
					createOutputRTMPSettings.ChunkSize = helper.IntInt64(v.(int))
				}
				modifyOutputInfo.RTMPSettings = &createOutputRTMPSettings
			}
			if rTPSettingsMap, ok := helper.InterfaceToMap(dMap, "rtp_settings"); ok {
				createOutputInfoRTPSettings := mps.CreateOutputInfoRTPSettings{}
				if v, ok := rTPSettingsMap["destinations"]; ok {
					for _, item := range v.([]interface{}) {
						destinationsMap := item.(map[string]interface{})
						createOutputRTPSettingsDestinations := mps.CreateOutputRTPSettingsDestinations{}
						if v, ok := destinationsMap["ip"]; ok {
							createOutputRTPSettingsDestinations.Ip = helper.String(v.(string))
						}
						if v, ok := destinationsMap["port"]; ok {
							createOutputRTPSettingsDestinations.Port = helper.IntInt64(v.(int))
						}
						createOutputInfoRTPSettings.Destinations = append(createOutputInfoRTPSettings.Destinations, &createOutputRTPSettingsDestinations)
					}
				}
				if v, ok := rTPSettingsMap["fec"]; ok {
					createOutputInfoRTPSettings.FEC = helper.String(v.(string))
				}
				if v, ok := rTPSettingsMap["idle_timeout"]; ok {
					createOutputInfoRTPSettings.IdleTimeout = helper.IntInt64(v.(int))
				}
				modifyOutputInfo.RTPSettings = &createOutputInfoRTPSettings
			}
			if v, ok := dMap["allow_ip_list"]; ok {
				allowIpListSet := v.(*schema.Set).List()
				for i := range allowIpListSet {
					if allowIpListSet[i] != nil {
						allowIpList := allowIpListSet[i].(string)
						modifyOutputInfo.AllowIpList = append(modifyOutputInfo.AllowIpList, &allowIpList)
					}
				}
			}
			if protocol == PROTOCOL_SRT || protocol == PROTOCOL_RTMP_PULL {
				if v, ok := dMap["max_concurrent"]; ok {
					modifyOutputInfo.MaxConcurrent = helper.IntUint64(v.(int))
				}
			}
			request.Output = &modifyOutputInfo
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().ModifyStreamLinkOutputInfo(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mps output failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMpsOutputRead(d, meta)
}

func resourceTencentCloudMpsOutputDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_output.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	flowId := idSplit[0]
	outputId := idSplit[1]

	if err := service.DeleteMpsOutputById(ctx, flowId, outputId); err != nil {
		return err
	}

	return nil
}
