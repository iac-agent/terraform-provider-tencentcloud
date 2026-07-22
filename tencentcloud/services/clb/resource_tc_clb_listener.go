package clb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClbListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbListenerCreate,
		Read:   resourceTencentCloudClbListenerRead,
		Update: resourceTencentCloudClbListenerUpdate,
		Delete: resourceTencentCloudClbListenerDelete,
		Importer: &schema.ResourceImporter{
			State: helper.ImportWithDefaultValue(map[string]interface{}{
				"scheduler": CLB_LISTENER_SCHEDULER_WRR,
			}),
		},
		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description: "CLB的ID。",
			},
			"listener_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description: "CLB监听器名称，取值只能是汉字、英文字母、数字、下划线和连字符“-”。",
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 65535),
				Description: "CLB监听端口。",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CLB_LISTENER_PROTOCOL),
				Description: "侦听器内的协议类型。有效值：“TCP”、“UDP”、“HTTP”、“HTTPS”、“TCP_SSL”和“QUIC”。",
			},
			"health_check_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "是否开启健康检查。",
			},
			"health_check_time_out": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 60),
				Description: "健康检查响应超时。有效值范围：[2~60]秒。默认值为 2 秒。响应超时需要小于检查间隔。注意：仅支持 `TCP`、`UDP`、`TCP_SSL` 协议的监听器。",
			},
			"health_check_interval_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 300),
				Description: "健康检查的间隔时间。有效值范围：[2~300]秒。默认为 5 秒。注意：TCP/UDP/TCP_SSL监听可以直接配置，HTTP/HTTPS监听需要在tencentcloud_clb_listener_rule中配置。",
			},
			"health_check_health_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 10),
				Description: "健康检查的健康阈值，默认为`3`。如果连续3次健康检查返回成功结果，则判定后端云服务器健康。取值范围为2-10。注意：TCP/UDP/TCP_SSL监听可以直接配置，HTTP/HTTPS监听需要在tencentcloud_clb_listener_rule中配置。",
			},
			"health_check_unhealth_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 10),
				Description: "健康检查不健康阈值，默认为`3`。" +
					"If a success result is returned for the health check 3 consecutive times, " +
					"the CVM is identified as unhealthy. The value range is [2-10]. " +
					"NOTES: TCP/UDP/TCP_SSL listener allows direct configuration, " +
					"HTTP/HTTPS listener needs to be configured in `tencentcloud_clb_listener_rule`.",
			},
			"health_check_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(HEALTH_CHECK_TYPE),
				Description: "用于健康检查的协议。有效值：“CUSTOM”、“TCP”、“HTTP”、“HTTPS”、“PING”、“GRPC”。",
			},
			"health_check_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 65535),
				Description: "健康检查端口默认为后端服务的端口。" +
					"Unless you want to specify a specific port, it is recommended to leave it blank. " +
					"Only applicable to TCP/UDP listener.",
			},
			"health_check_http_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(HTTP_VERSION),
				Description: "后端服务的 HTTP 版本。当`health_check_type`的值为" +
					"the health check protocol is `HTTP`, this field is required. " +
					"Valid values: `HTTP/1.0`, `HTTP/1.1`.",
			},
			"health_check_http_code": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 31),
				Description: "TCP监听器的HTTP健康检查代码，有效值范围：[1~31]。当`health_check_type`的值为" +
					"the health check protocol is `HTTP`, this field is required. Valid values: `1`, `2`, `4`, `8`, `16`. " +
					"`1` means http_1xx, `2` means http_2xx, `4` means http_3xx, `8` means http_4xx, `16` means http_5xx." +
					"If you want multiple return codes to indicate health, need to add the corresponding values.",
			},
			"health_check_http_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "TCP监听的HTTP健康检查路径。",
			},
			"health_check_http_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "TCP监听的HTTP健康检查域。",
			},
			"health_check_http_method": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CLB_HTTP_METHOD),
				Description: "TCP监听的HTTP健康检查方法。有效值：“HEAD”、“GET”。",
			},
			"health_check_context_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CONTEX_TYPE),
				Description: "健康检查协议。当健康检查协议的“health_check_type”值为“CUSTOM”时，" +
					"this field is required, which represents the input format of the health check. " +
					"Valid values: `HEX`, `TEXT`.",
			},
			"health_check_send_context": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(0, 500),
				Description: "它代表健康检查发送的请求的内容。" +
					"When the value of `health_check_type` of the health check protocol is `CUSTOM`, " +
					"this field is required. Only visible ASCII characters are allowed and the maximum length is 500. " +
					"When `health_check_context_type` value is `HEX`, " +
					"the characters of SendContext and RecvContext can only be selected in `0123456789ABCDEF` " +
					"and the length must be even digits.",
			},
			"health_check_recv_context": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(0, 500),
				Description: "它代表健康检查返回的结果。" +
					"When the value of `health_check_type` of the health check protocol is `CUSTOM`, " +
					"this field is required. Only ASCII visible characters are allowed and the maximum length is 500. " +
					"When `health_check_context_type` value is `HEX`, " +
					"the characters of SendContext and RecvContext can only be selected in `0123456789ABCDEF` " +
					"and the length must be even digits.",
			},
			"health_source_ip_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description: "指定健康检查源IP的类型。 `0`（默认）：CLB VIP。 `1`：100.64 IP 范围。",
			},
			"certificate_ssl_mode": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"multi_cert_info"},
				ValidateFunc:  tccommon.ValidateAllowedStringValue(CERT_SSL_MODE),
				Description: "证书类型。有效值：“单向”、“相互”。注意：仅支持“HTTPS”和“TCP_SSL”协议的监听器，并且必须在可用时进行设置。",
			},
			"certificate_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"multi_cert_info"},
				Description: "服务器证书的 ID。注意：仅支持“HTTPS”和“TCP_SSL”协议的监听器，并且必须在可用时进行设置。",
			},
			"certificate_ca_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"multi_cert_info"},
				Description: "客户端证书的ID。注意：仅支持`HTTPS`和`TCP_SSL`协议的监听，且ssl模式为`MUTUAL`时必须设置。",
			},
			"multi_cert_info": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"certificate_ssl_mode", "certificate_id", "certificate_ca_id"},
				Description: "证书信息。您可以指定多个具有不同算法类型的服务器端证书。该参数仅适用于未开启SNI功能的HTTPS监听。不能同时指定Certificate 和MultiCertInfo。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ssl_mode": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CERT_SSL_MODE),
							Description: "认证类型。值：UNIDIRECTIONAL（单向身份验证）、MUTUAL（双向身份验证）。",
						},
						"cert_id_list": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "服务器证书 ID 列表。",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"session_expire_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(30, 3600),
				Description: "CLB 侦听器中的会话持续时间。注意：当调度程序指定为“WRR”时可用，当侦听器协议为“TCP_SSL”时不可用。注意：TCP/UDP/TCP_SSL监听可以直接配置，HTTP/HTTPS监听需要在tencentcloud_clb_listener_rule中配置。",
			},
			"scheduler": {
				Type:         schema.TypeString,
				Default:      CLB_LISTENER_SCHEDULER_WRR,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CLB_LISTENER_SCHEDULER),
				Description: "CLB监听的调度方式，可用值为'WRR'和'LEAST_CONN'。默认值为“WRR”。注意：`HTTP`和`HTTPS`协议的监听器还支持`IP Hash`方法。注意：TCP/UDP/TCP_SSL监听可以直接配置，HTTP/HTTPS监听需要在tencentcloud_clb_listener_rule中配置。",
			},
			"sni_switch": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Description: "指示是否启用 SNI，并且仅支持协议“HTTPS”。如果启用，则可以为tencentcloud_clb_listener_rule中的每条规则设置证书，否则所有规则都有证书。",
			},
			"target_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{CLB_TARGET_TYPE_NODE, CLB_TARGET_TYPE_TARGETGROUP, CLB_TARGET_TYPE_TARGETGROUP_V2}),
				Description: "后端目标类型。有效值：“NODE”、“TARGETGROUP”、“TARGETGROUP-V2”。 `NODE` 表示绑定普通节点，`TARGETGROUP` 表示绑定目标组。注意：TCP/UDP/TCP_SSL监听必须配置，HTTP/HTTPS监听需要在tencentcloud_clb_listener_rule中配置。",
			},
			"session_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{CLB_SESSION_TYPE_NORMAL, CLB_SESSION_TYPE_QUIC}),
				Description: "会话持续类型。有效值：`NORMAL`：默认会话持久类型； `QUIC_CID`：通过 QUIC 连接 ID 进行会话持久化。 “QUIC_CID”值只能在 UDP 侦听器中配置。如果不指定该字段，则使用默认的会话保持类型。",
			},
			"keepalive_enable": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "是否启用持久连接。该参数仅适用于HTTP和HTTPS监听。有效值：0（禁用；默认值）和 1（启用）。",
			},
			"end_port": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Computed:    true,
				Optional:    true,
				Description: "该参数用于指定结束端口，创建端口范围监听时需要此参数。输入“Ports”参数时只能传入一个成员，用于指定起始端口。如果您想尝试端口范围功能，请【提交工单】(https://console.cloud.tencent.com/workorder/category)。",
			},
			"h2c_switch": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Computed:    true,
				Optional:    true,
				Description: "启用H2C内网HTTP监听开关。",
			},
			"snat_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "是否启用SNAT。",
			},
			"deregister_target_rst": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "解绑真实服务器时是否发送TCP RST包给客户端。该参数仅适用于TCP监听。",
			},
			"idle_connect_timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "连接空闲超时时间（以秒为单位）。它仅适用于 TCP 侦听器。值范围：共享实例和专用实例为300-900；对于 LCU 支持的 CLB 实例，为 300-2000。默认为 900。设置长于 2000 秒的周期（最多 3600 秒）。请提交工单进行处理。",
			},
			"reschedule_target_zero_weight": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "重新调度功能，以权重0为开关，当后端服务器的权重设置为0时，触发重新调度。仅TCP/UDP监听支持。",
			},
			"reschedule_unhealthy": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "重新调度功能、健康检查异常开关。启用此开关会在后端服务器健康检查失败时触发重新调度。仅 TCP/UDP 侦听器支持。",
			},
			"reschedule_expand_target": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "重新调度功能是后端服务伸缩的开关，当后端服务器增加或减少时，会触发重新调度。仅受 TCP/UDP 侦听器支持。",
			},
			"reschedule_start_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "重新安排触发开始时间，值范围为 0 到 3600 秒。仅受 TCP/UDP 侦听器支持。",
			},
			"reschedule_interval": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "重新安排的触发持续时间，范围从 0 到 3600 秒。仅 TCP/UDP 侦听器支持。",
			},
			//computed
			"listener_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "该CLB监听器的ID。",
			},
		},
	}
}

func resourceTencentCloudClbListenerCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_listener.create")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	clbId := d.Get("clb_id").(string)
	listenerName := d.Get("listener_name").(string)
	request := clb.NewCreateListenerRequest()

	request.LoadBalancerId = helper.String(clbId)
	request.ListenerNames = []*string{&listenerName}

	port := int64(d.Get("port").(int))
	ports := []*int64{&port}
	request.Ports = ports
	protocol := d.Get("protocol").(string)
	request.Protocol = helper.String(protocol)

	healthSetFlag, healthCheck, healthErr := checkHealthCheckPara(ctx, d, protocol, HEALTH_APPLY_TYPE_LISTENER)
	if healthErr != nil {
		return healthErr
	}
	if healthSetFlag {
		request.HealthCheck = healthCheck
	}

	certificateSetFlag, certificateInput, certErr := checkCertificateInputPara(ctx, d, meta)

	if certErr != nil {
		return certErr
	}
	if certificateSetFlag {
		request.Certificate = certificateInput
	} else {
		if protocol == CLB_LISTENER_PROTOCOL_TCPSSL {
			return fmt.Errorf("[CHECK][CLB listener][Create] check: certificated need to be set when protocol is TCPSSL")
		}
	}

	multiCertificateSetFlag, multiCertInput, certErr := checkMultiCertificateInputPara(ctx, d, meta)
	if certErr != nil {
		return certErr
	}

	if multiCertificateSetFlag {
		request.MultiCertInfo = multiCertInput
	} else {
		if protocol == CLB_LISTENER_PROTOCOL_TCPSSL {
			return fmt.Errorf("[CHECK][CLB listener][Create] check: certificated need to be set when protocol is TCPSSL")
		}
	}

	scheduler := ""
	if v, ok := d.GetOk("scheduler"); ok {
		if v == CLB_LISTENER_SCHEDULER_IP_HASH {
			return fmt.Errorf("[CHECK][CLB listener][Create] check: Scheduler 'IP_HASH' can only be set with rule of listener HTTP/HTTPS")
		}
		scheduler = v.(string)
		request.Scheduler = helper.String(scheduler)
	}
	if v, ok := d.GetOk("target_type"); ok {
		targetType := v.(string)
		request.TargetType = &targetType
	} else if protocol == CLB_LISTENER_PROTOCOL_TCP || protocol == CLB_LISTENER_PROTOCOL_UDP ||
		protocol == CLB_LISTENER_PROTOCOL_TCPSSL || protocol == CLB_LISTENER_PROTOCOL_QUIC {
		targetType := CLB_TARGET_TYPE_NODE
		request.TargetType = &targetType
	}

	if v, ok := d.GetOk("session_expire_time"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP || protocol == CLB_LISTENER_PROTOCOL_UDP) {
			return fmt.Errorf("[CHECK][CLB listener][Create] check: session_expire_time can only be set with protocol TCP/UDP or rule of listener HTTP/HTTPS")
		}
		if scheduler != CLB_LISTENER_SCHEDULER_WRR && scheduler != "" {
			return fmt.Errorf("[CHECK][CLB listener][Create] check: session_expire_time can only be set when scheduler is WRR ")
		}
		vv := int64(v.(int))
		request.SessionExpireTime = &vv
	}
	if v, ok := d.GetOkExists("sni_switch"); ok {
		if protocol != CLB_LISTENER_PROTOCOL_HTTPS {
			return fmt.Errorf("[CHECK][CLB listener][Create] check: sni_switch can only be set with protocol HTTPS ")
		} else {
			vv := v.(bool)
			vvv := int64(0)
			if vv {
				vvv = 1
			} else {
				if !certificateSetFlag && !multiCertificateSetFlag {
					return fmt.Errorf("[CHECK][CLB listener][Create] check: certificated need to be set when protocol is HTTPS")
				}
			}
			request.SniSwitch = &vvv
		}
	}

	if v, ok := d.GetOk("session_type"); ok {
		request.SessionType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("keepalive_enable"); ok {
		request.KeepaliveEnable = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("end_port"); ok {
		request.EndPort = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("h2c_switch"); ok {
		request.H2cSwitch = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("snat_enable"); ok {
		request.SnatEnable = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("deregister_target_rst"); ok {
		request.DeregisterTargetRst = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("idle_connect_timeout"); ok {
		request.IdleConnectTimeout = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("reschedule_target_zero_weight"); ok {
		request.RescheduleTargetZeroWeight = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("reschedule_unhealthy"); ok {
		request.RescheduleUnhealthy = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("reschedule_expand_target"); ok {
		request.RescheduleExpandTarget = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("reschedule_start_time"); ok {
		request.RescheduleStartTime = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("reschedule_interval"); ok {
		request.RescheduleInterval = helper.IntInt64(v.(int))
	}

	var response *clb.CreateListenerResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().CreateListener(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

			if result == nil || result.Response == nil || result.Response.RequestId == nil {
				return resource.NonRetryableError(fmt.Errorf("Create CLB listener failed, Response si nil."))
			}

			requestId := *result.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
			if retryErr != nil {
				return resource.NonRetryableError(errors.WithStack(retryErr))
			}
		}

		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create CLB listener failed, reason:%+v", logId, err)
		return err
	}
	if response.Response.ListenerIds == nil || len(response.Response.ListenerIds) < 1 {
		return fmt.Errorf("[CHECK][CLB listener][Create] check: Response error, listener id is null")
	}
	listenerId := *response.Response.ListenerIds[0]

	//this ID style changes since terraform 1.47.0
	d.SetId(clbId + tccommon.FILED_SP + listenerId)
	return resourceTencentCloudClbListenerRead(d, meta)
}

func resourceTencentCloudClbListenerRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_listener.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	resourceId := d.Id()
	var listenerId = resourceId
	items := strings.Split(resourceId, tccommon.FILED_SP)
	itemLength := len(items)
	clbId := d.Get("clb_id").(string)

	if itemLength == 1 && clbId == "" {
		return fmt.Errorf("the old style listenerId %s does not support import, please use clbId#listenerId style", resourceId)
	} else if itemLength == 2 && clbId == "" {
		listenerId = items[1]
		clbId = items[0]
	} else if itemLength == 2 && clbId != "" {
		listenerId = items[1]
	} else if itemLength != 1 && itemLength != 2 {
		return fmt.Errorf("broken ID %s", resourceId)
	}

	var instance *clb.Listener
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := clbService.DescribeListenerById(ctx, listenerId, clbId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB listener failed, reason:%s\n ", logId, err.Error())
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("clb_id", clbId)
	_ = d.Set("listener_id", instance.ListenerId)
	_ = d.Set("port", instance.Port)
	_ = d.Set("protocol", instance.Protocol)
	if instance.ListenerName != nil {
		_ = d.Set("listener_name", instance.ListenerName)
	}
	if instance.TargetType != nil {
		_ = d.Set("target_type", instance.TargetType)
	}
	if instance.SessionExpireTime != nil {
		_ = d.Set("session_expire_time", instance.SessionExpireTime)
	}
	if *instance.Protocol == CLB_LISTENER_PROTOCOL_TCP || *instance.Protocol == CLB_LISTENER_PROTOCOL_TCPSSL ||
		*instance.Protocol == CLB_LISTENER_PROTOCOL_UDP || *instance.Protocol == CLB_LISTENER_PROTOCOL_QUIC {
		_ = d.Set("scheduler", instance.Scheduler)
	}
	_ = d.Set("sni_switch", *instance.SniSwitch > 0)

	//health check
	if instance.HealthCheck != nil {
		healthCheckSwitch := false
		if *instance.HealthCheck.HealthSwitch == int64(1) {
			healthCheckSwitch = true
		}
		_ = d.Set("health_check_switch", healthCheckSwitch)
		if instance.HealthCheck.IntervalTime != nil {
			_ = d.Set("health_check_interval_time", instance.HealthCheck.IntervalTime)
		}
		if instance.HealthCheck.TimeOut != nil {
			_ = d.Set("health_check_time_out", instance.HealthCheck.TimeOut)
		}
		if instance.HealthCheck.HealthNum != nil {
			_ = d.Set("health_check_health_num", instance.HealthCheck.HealthNum)
		}
		if instance.HealthCheck.UnHealthNum != nil {
			_ = d.Set("health_check_unhealth_num", instance.HealthCheck.UnHealthNum)
		}
		if instance.HealthCheck.CheckPort != nil {
			_ = d.Set("health_check_port", instance.HealthCheck.CheckPort)
		}
		if instance.HealthCheck.CheckType != nil {
			_ = d.Set("health_check_type", instance.HealthCheck.CheckType)
		}
		if instance.HealthCheck.HttpCode != nil {
			_ = d.Set("health_check_http_code", instance.HealthCheck.HttpCode)
		}
		if instance.HealthCheck.HttpCheckPath != nil {
			_ = d.Set("health_check_http_path", instance.HealthCheck.HttpCheckPath)
		}
		if instance.HealthCheck.HttpCheckDomain != nil {
			_ = d.Set("health_check_http_domain", instance.HealthCheck.HttpCheckDomain)
		}
		if instance.HealthCheck.HttpCheckMethod != nil {
			_ = d.Set("health_check_http_method", instance.HealthCheck.HttpCheckMethod)
		}
		if instance.HealthCheck.HttpVersion != nil {
			_ = d.Set("health_check_http_version", instance.HealthCheck.HttpVersion)
		}
		if instance.HealthCheck.ContextType != nil {
			_ = d.Set("health_check_context_type", instance.HealthCheck.ContextType)
		}
		if instance.HealthCheck.SendContext != nil {
			_ = d.Set("health_check_send_context", instance.HealthCheck.SendContext)
		}
		if instance.HealthCheck.RecvContext != nil {
			_ = d.Set("health_check_recv_context", instance.HealthCheck.RecvContext)
		}
		if instance.HealthCheck.SourceIpType != nil {
			_ = d.Set("health_source_ip_type", instance.HealthCheck.SourceIpType)
		}
	}

	if instance.Certificate != nil {
		// check single cert or multi cert
		if instance.Certificate.ExtCertIds != nil && len(instance.Certificate.ExtCertIds) > 0 {
			multiCertInfo := make([]map[string]interface{}, 0, 1)
			multiCert := make(map[string]interface{}, 0)
			certIds := make([]string, 0)
			if instance.Certificate.SSLMode != nil {
				multiCert["ssl_mode"] = *instance.Certificate.SSLMode
			}

			if instance.Certificate.CertId != nil {
				certIds = append(certIds, *instance.Certificate.CertId)
			}

			for _, item := range instance.Certificate.ExtCertIds {
				certIds = append(certIds, *item)
			}

			multiCert["cert_id_list"] = certIds
			multiCertInfo = append(multiCertInfo, multiCert)
			_ = d.Set("multi_cert_info", multiCertInfo)
		} else {
			_ = d.Set("certificate_ssl_mode", instance.Certificate.SSLMode)
			_ = d.Set("certificate_id", instance.Certificate.CertId)
			if instance.Certificate.CertCaId != nil {
				_ = d.Set("certificate_ca_id", instance.Certificate.CertCaId)
			}
		}
	}

	if instance.SessionType != nil {
		_ = d.Set("session_type", instance.SessionType)
	}
	if instance.KeepaliveEnable != nil {
		_ = d.Set("keepalive_enable", instance.KeepaliveEnable)
	}

	if instance.EndPort != nil {
		_ = d.Set("end_port", instance.EndPort)
	}

	if instance.AttrFlags != nil && len(instance.AttrFlags) > 0 {
		if tccommon.IsContains(helper.PStrings(instance.AttrFlags), "H2cSwitch") {
			_ = d.Set("h2c_switch", true)
		} else {
			_ = d.Set("h2c_switch", false)
		}

		if tccommon.IsContains(helper.PStrings(instance.AttrFlags), "SnatEnable") {
			_ = d.Set("snat_enable", true)
		} else {
			_ = d.Set("snat_enable", false)
		}
	} else {
		_ = d.Set("h2c_switch", false)
		_ = d.Set("snat_enable", false)
	}

	if instance.DeregisterTargetRst != nil {
		_ = d.Set("deregister_target_rst", instance.DeregisterTargetRst)
	}

	if instance.IdleConnectTimeout != nil {
		_ = d.Set("idle_connect_timeout", instance.IdleConnectTimeout)
	}

	_ = d.Set("reschedule_target_zero_weight", false)
	_ = d.Set("reschedule_unhealthy", false)
	_ = d.Set("reschedule_expand_target", false)
	if instance.AttrFlags != nil {
		for _, item := range instance.AttrFlags {
			if item != nil && *item == "RescheduleTargetZeroWeight" {
				_ = d.Set("reschedule_target_zero_weight", true)
			}

			if item != nil && *item == "RescheduleUnhealthy" {
				_ = d.Set("reschedule_unhealthy", true)
			}

			if item != nil && *item == "RescheduleExpandTarget" {
				_ = d.Set("reschedule_expand_target", true)
			}
		}
	}

	if instance.RescheduleStartTime != nil {
		_ = d.Set("reschedule_start_time", instance.RescheduleStartTime)
	}

	if instance.RescheduleInterval != nil {
		_ = d.Set("reschedule_interval", instance.RescheduleInterval)
	}

	return nil
}

func resourceTencentCloudClbListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_listener.update")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	resourceId := d.Id()
	items := strings.Split(resourceId, tccommon.FILED_SP)
	itemLength := len(items)
	listenerId := items[itemLength-1]
	clbId := d.Get("clb_id").(string)
	changed := false
	scheduler := ""
	listenerName := ""
	sessionExpireTime := 0
	protocol := d.Get("protocol").(string)

	request := clb.NewModifyListenerRequest()
	request.ListenerId = helper.String(listenerId)
	request.LoadBalancerId = helper.String(clbId)

	if d.HasChange("listener_name") {
		changed = true
		listenerName = d.Get("listener_name").(string)
		request.ListenerName = helper.String(listenerName)
	}

	if d.HasChange("scheduler") {
		changed = true
		scheduler = d.Get("scheduler").(string)
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP || protocol == CLB_LISTENER_PROTOCOL_UDP ||
			protocol == CLB_LISTENER_PROTOCOL_TCPSSL || protocol == CLB_LISTENER_PROTOCOL_QUIC) {
			return fmt.Errorf("[CHECK][CLB listener %s][Update] check: Scheduler can only be set with listener protocol TCP/UDP/TCP_SSL/QUIC or rule of listener HTTP/HTTPS", listenerId)
		}
		if scheduler == CLB_LISTENER_SCHEDULER_IP_HASH {
			return fmt.Errorf("[CHECK][CLB listener %s][Update] check: Scheduler 'IP_HASH' can only be set with rule of listener HTTP/HTTPS", listenerId)
		}
		request.Scheduler = helper.String(scheduler)
	}

	if d.HasChange("session_expire_time") {
		changed = true
		sessionExpireTime = d.Get("session_expire_time").(int)
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP || protocol == CLB_LISTENER_PROTOCOL_UDP) {
			return fmt.Errorf("[CHECK][CLB listener %s][Update] check: session_expire_time can only be set with protocol TCP/UDP or rule of listener HTTP/HTTPS", listenerId)
		}
		if scheduler != CLB_LISTENER_SCHEDULER_WRR && scheduler != "" {
			return fmt.Errorf("[CHECK][CLB listener %s][Update] check: session_expire_time can only be set when scheduler is WRR", listenerId)
		}
		sessionExpireTime64 := int64(sessionExpireTime)
		request.SessionExpireTime = &sessionExpireTime64
	}

	healthSetFlag, healthCheck, healthErr := checkHealthCheckPara(ctx, d, protocol, HEALTH_APPLY_TYPE_LISTENER)
	if healthErr != nil {
		return healthErr
	}
	if healthSetFlag {
		changed = true
		request.HealthCheck = healthCheck
	}

	certificateSetFlag, certificateInput, certErr := checkCertificateInputPara(ctx, d, meta)
	if certErr != nil {
		return certErr
	}
	if certificateSetFlag {
		changed = true
		request.Certificate = certificateInput
	}

	multiCertificateSetFlag, multiCertInput, certErr := checkMultiCertificateInputPara(ctx, d, meta)
	if certErr != nil {
		return certErr
	}

	if multiCertificateSetFlag {
		changed = true
		request.MultiCertInfo = multiCertInput
	}

	if d.HasChange("target_type") {
		changed = true
		targetType := d.Get("target_type").(string)
		request.TargetType = helper.String(targetType)
	}

	if d.HasChange("session_type") {
		changed = true
		sessionType := d.Get("session_type").(string)
		request.SessionType = helper.String(sessionType)
	}

	if d.HasChange("keepalive_enable") {
		changed = true
		keepaliveEnable := d.Get("keepalive_enable").(int)
		request.KeepaliveEnable = helper.IntInt64(keepaliveEnable)
	}

	if d.HasChange("snat_enable") {
		changed = true
		if v, ok := d.GetOkExists("snat_enable"); ok {
			request.SnatEnable = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("deregister_target_rst") {
		changed = true
		if v, ok := d.GetOkExists("deregister_target_rst"); ok {
			request.DeregisterTargetRst = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("idle_connect_timeout") {
		changed = true
		if v, ok := d.GetOkExists("idle_connect_timeout"); ok {
			request.IdleConnectTimeout = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("reschedule_target_zero_weight") {
		changed = true
		if v, ok := d.GetOkExists("reschedule_target_zero_weight"); ok {
			request.RescheduleTargetZeroWeight = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("reschedule_unhealthy") {
		changed = true
		if v, ok := d.GetOkExists("reschedule_unhealthy"); ok {
			request.RescheduleUnhealthy = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("reschedule_expand_target") {
		changed = true
		if v, ok := d.GetOkExists("reschedule_expand_target"); ok {
			request.RescheduleExpandTarget = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("reschedule_start_time") {
		changed = true
		if v, ok := d.GetOkExists("reschedule_start_time"); ok {
			request.RescheduleStartTime = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("reschedule_interval") {
		changed = true
		if v, ok := d.GetOkExists("reschedule_interval"); ok {
			request.RescheduleInterval = helper.IntInt64(v.(int))
		}
	}

	if changed {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyListener(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

				if response == nil || response.Response == nil || response.Response.RequestId == nil {
					return resource.NonRetryableError(fmt.Errorf("Modify CLB listener failed, Response si nil."))
				}

				requestId := *response.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
				if retryErr != nil {
					return resource.NonRetryableError(errors.WithStack(retryErr))
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update CLB listener failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudClbListenerRead(d, meta)
}

func resourceTencentCloudClbListenerDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_listener.delete")()
	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	resourceId := d.Id()
	items := strings.Split(resourceId, tccommon.FILED_SP)
	itemLength := len(items)
	listenerId := items[itemLength-1]
	clbId := d.Get("clb_id").(string)

	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := clbService.DeleteListenerById(ctx, clbId, listenerId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CLB listener failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
