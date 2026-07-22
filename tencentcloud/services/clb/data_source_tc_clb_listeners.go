package clb

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClbListeners() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbListenersRead,

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "需要查询的CLB ID。",
			},
			"listener_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "需要查询的监听者id。",
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CLB_LISTENER_PROTOCOL),
				Description: "侦听器内的协议类型，可用值包括“TCP”、“UDP”、“HTTP”、“HTTPS”和“TCP_SSL”。",
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 65535),
				Description: "CLB监听端口。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"listener_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "云负载均衡器的监听器列表。每个元素包含以下属性：",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"clb_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB的ID。",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "监听者ID。",
						},
						"listener_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB 侦听器的名称。",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "侦听器的协议。可用值包括“HTTP”、“HTTPS”、“TCP”、“UDP”、“TCP_SSL”。",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CLB监听端口。",
						},
						"health_check_switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否开启健康检查。",
						},
						"health_check_time_out": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "健康检查响应超时。取值范围为 2-60 秒，默认为“2”秒。响应超时需要小于检查间隔。注意： TCP/UDP/TCP_SSL 侦听器允许直接配置。",
						},
						"health_check_interval_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "健康检查的间隔时间。值范围为 2-300 秒，默认为“5”秒。注意：TCP/UDP/TCP_SSL监听可以直接配置，HTTP/HTTPS监听需要在tencentcloud_clb_listener_rule中配置。",
						},
						"health_check_health_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "健康检查的健康阈值，默认为`3`。如果健康检查连续3次返回成功，则判定该云服务器健康。取值范围为2-10。注意：TCP/UDP/TCP_SSL监听可以直接配置，HTTP/HTTPS监听需要在tencentcloud_clb_listener_rule中配置。",
						},
						"health_check_unhealth_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "健康检查不健康阈值，默认为`3`。如果健康检查连续3次返回成功，则判定该云服务器不健康。取值范围为2-10。注意：TCP/UDP/TCP_SSL监听可以直接配置，HTTP/HTTPS监听需要在tencentcloud_clb_listener_rule中配置。",
						},
						"health_check_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "用于健康检查的协议。",
						},
						"health_check_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "健康检查端口是后端服务的端口。",
						},
						"health_check_http_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "后端服务的 HTTP 版本。",
						},
						"health_check_http_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "TCP监听的HTTP健康检查代码。",
						},
						"health_check_http_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "TCP监听的HTTP健康检查路径。",
						},
						"health_check_http_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "TCP监听的HTTP健康检查域。",
						},
						"health_check_http_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "TCP监听的HTTP健康检查方法。",
						},
						"health_check_context_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "健康检查协议。",
						},
						"health_check_send_context": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "它代表健康检查发送的请求的内容。",
						},
						"health_check_recv_context": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "它代表健康检查返回的结果。",
						},
						"certificate_ssl_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "证书类型，可用值包括“UNIDIRECTIONAL”、“MUTUAL”。注意：仅支持“HTTPS”和“TCP_SSL”协议的监听器，并且必须在可用时进行设置。",
						},
						"certificate_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "服务器证书的 ID。当协议为“HTTPS”或“TCP_SSL”时必须设置。注意：仅受“HTTPS”和“TCP_SSL”协议的侦听器支持，并且必须在可用时进行设置。",
						},
						"certificate_ca_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "客户端证书的 ID。当 SSLMode 为“mutual”时必须设置。注意：仅受“HTTPS”和“TCP_SSL”协议的侦听器支持。",
						},
						"session_expire_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CLB 侦听器中的会话持续时间。注意：TCP/UDP/TCP_SSL监听可以直接配置，HTTP/HTTPS监听需要在tencentcloud_clb_listener_rule中配置。",
						},
						"scheduler": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB监听的调度方式，可选值为`WRR`和`LEAST_CONN`。默认为“WRR”。注意：“HTTP”和“HTTPS”协议的监听器还支持“IP HASH”方法。注意：TCP/UDP/TCP_SSL监听可以直接配置，HTTP/HTTPS监听需要在tencentcloud_clb_listener_rule中配置。",
						},
						"sni_switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "指示是否启用SNI。注意：仅受“HTTPS”协议支持。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudClbListenersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clb_listeners.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	clbId := d.Get("clb_id").(string)

	params := make(map[string]interface{})
	params["clb_id"] = clbId
	if v, ok := d.GetOk("listener_id"); ok {
		listenerId := v.(string)
		params["listener_id"] = listenerId
		checkErr := ListenerIdCheck(listenerId)
		if checkErr != nil {
			return checkErr
		}
	}
	if v, ok := d.GetOk("port"); ok {
		params["port"] = v.(int)
	}
	if v, ok := d.GetOk("protocol"); ok {
		params["protocol"] = v.(string)
	}

	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var listeners []*clb.Listener
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := clbService.DescribeListenersByFilter(ctx, params)
		if e != nil {
			return tccommon.RetryError(e)
		}
		listeners = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB listeners failed, reason:%+v", logId, err)
		return err
	}
	listenerList := make([]map[string]interface{}, 0, len(listeners))
	ids := make([]string, 0, len(listeners))
	for _, listener := range listeners {
		mapping := map[string]interface{}{
			"clb_id":        clbId,
			"listener_id":   listener.ListenerId,
			"listener_name": listener.ListenerName,
			"protocol":      listener.Protocol,
			"port":          listener.Port,
		}
		if listener.SessionExpireTime != nil {
			mapping["session_expire_time"] = listener.SessionExpireTime
		}
		if listener.SniSwitch != nil {
			sniSwitch := false
			if *listener.SniSwitch == int64(1) {
				sniSwitch = true
			}
			mapping["sni_switch"] = sniSwitch
		}
		mapping["scheduler"] = listener.Scheduler
		if listener.HealthCheck != nil {
			health_check_switch := false
			if *listener.HealthCheck.HealthSwitch == int64(1) {
				health_check_switch = true
			}
			mapping["health_check_switch"] = health_check_switch
			mapping["health_check_time_out"] = listener.HealthCheck.TimeOut
			mapping["health_check_interval_time"] = listener.HealthCheck.IntervalTime
			mapping["health_check_health_num"] = listener.HealthCheck.HealthNum
			mapping["health_check_unhealth_num"] = listener.HealthCheck.UnHealthNum
			mapping["health_check_http_code"] = listener.HealthCheck.HttpCode
			mapping["health_check_http_path"] = listener.HealthCheck.HttpCheckPath
			mapping["health_check_http_domain"] = listener.HealthCheck.HttpCheckDomain
			mapping["health_check_http_method"] = listener.HealthCheck.HttpCheckMethod
			mapping["health_check_http_version"] = listener.HealthCheck.HttpVersion
			mapping["health_check_context_type"] = listener.HealthCheck.ContextType
			mapping["health_check_send_context"] = listener.HealthCheck.SendContext
			mapping["health_check_recv_context"] = listener.HealthCheck.RecvContext
			mapping["health_check_type"] = listener.HealthCheck.CheckType
			mapping["health_check_port"] = listener.HealthCheck.CheckPort
		}
		if listener.Certificate != nil {
			mapping["certificate_ssl_mode"] = listener.Certificate.SSLMode
			mapping["certificate_id"] = listener.Certificate.CertId
			if listener.Certificate.CertCaId != nil {
				mapping["certificate_ca_id"] = listener.Certificate.CertCaId
			}
		}
		listenerList = append(listenerList, mapping)
		ids = append(ids, *listener.ListenerId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("listener_list", listenerList); e != nil {
		log.Printf("[CRITAL]%s provider set CLB listener list fail, reason:%+v", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), listenerList); e != nil {
			return e
		}
	}

	return nil
}
