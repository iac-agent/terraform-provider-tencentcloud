package gaap

import (
	"context"
	"errors"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudGaapLayer4Listener() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapLayer4ListenerCreate,
		Read:   resourceTencentCloudGaapLayer4ListenerRead,
		Update: resourceTencentCloudGaapLayer4ListenerUpdate,
		Delete: resourceTencentCloudGaapLayer4ListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"TCP", "UDP"}),
				ForceNew:     true,
				Description:  "协议 的 layer4 listener. 有效 值: `TCP` 和 `UDP`。",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 30),
				Description:  "名称 layer4 listener， 最大 长度 是 30。",
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidatePort,
				ForceNew:     true,
				Description:  "端口 的 layer4 listener。",
			},
			"scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "rr",
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"rr", "wrr", "lc"}),
				Description:  "Scheduling 策略 的 layer4 listener，默认值为 `rr`. 有效 值: `rr`，`wrr` 和 `lc`。",
			},
			"realserver_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"IP", "DOMAIN"}),
				ForceNew:     true,
				Description:  "类型 realserver. 有效 值: `IP` 和 `DOMAIN`. NOTES: 当 `协议` 是 指定 作为 `TCP` 和 `scheduler` 是 指定 作为 `wrr`， item 可以 仅 是 集合 到 `IP`。",
			},
			"proxy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID GAAP proxy。",
			},
			"health_check": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "表示是否health check 是 启用，默认值为 `false`。",
			},
			"interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ValidateFunc: tccommon.ValidateIntegerInRange(5, 300),
				Description:  "Interval 的 health check，默认值为 5s。",
			},
			"connect_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 60),
				Description:  "Timeout 的 health check response，should less 比 间隔，默认值为 2s. NOTES: Require less 比 `间隔`。",
			},
			"healthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 10),
				Description:  "Health 阈值，其中 表示how many consecutive inspections 是 successful， 来源 station 是 determined 到 是 healthy. Range 从 1 到 10. 默认值为 1。",
			},
			"unhealthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 10),
				Description:  "Unhealthy 阈值，其中 表示how many consecutive check failures 来源 station 是 considered unhealthy. Range 从 1 到 10. 默认值为 1。",
			},
			"check_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"PORT", "PING"}),
				Description:  "UDP 源站 服务器 health 类型 PORT 表示 check 端口，和 PING 表示 PING。",
			},
			"check_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "UDP 源站 station health check probe 端口",
			},
			"context_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"TEXT"}),
				Description:  "UDP 来源 station health check 端口 probe 消息类型: TEXT 表示 text. Only 使用 当 health check 类型 是 PORT。",
			},
			"send_context": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "UDP 来源 服务器 health check 端口 detection sends messages. Only 使用 当 health check 类型 是 PORT。",
			},
			"recv_context": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "UDP 来源 服务器 health check 端口 detects received messages. Only 使用 当 health check 类型 是 PORT。",
			},
			"client_ip_method": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      0,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description:  "way listener gets 客户端 IP，0 对于 TOA，1 对于 Proxy 协议，默认值为 0. NOTES: Only 支持 listeners 的 `TCP` 协议",
			},
			"realserver_bind_set": {
				Type:     schema.TypeSet,
				Optional: true,
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					return helper.HashString(fmt.Sprintf("%s-%s-%d-%d", m["id"].(string), m["ip"].(string), m["port"].(int), m["weight"].(int)))
				},
				Description: "An 信息 列表 GAAP realserver。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID GAAP realserver。",
						},
						"ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP 的 GAAP realserver。",
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: tccommon.ValidatePort,
							Description:  "端口 的 GAAP realserver。",
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: tccommon.ValidateIntegerInRange(1, 100),
							Description:  "Scheduling 权重，默认值为 `1`. 范围 的 值 是 [1,100]。",
						},
					},
				},
			},

			// computed
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "状态 layer4 listener。",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "创建时间 的 layer4 listener。",
			},
		},
	}
}

func resourceTencentCloudGaapLayer4ListenerCreate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_layer4_listener.create")()
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	protocol := d.Get("protocol").(string)
	name := d.Get("name").(string)
	port := d.Get("port").(int)
	scheduler := d.Get("scheduler").(string)
	realserverType := d.Get("realserver_type").(string)

	if protocol == "TCP" && realserverType == "DOMAIN" && scheduler == "wrr" {
		return errors.New("TCP listener DOMAIN realserver type doesn't support wrr scheduler")
	}

	proxyId := d.Get("proxy_id").(string)
	healthCheck := d.Get("health_check").(bool)

	interval := d.Get("interval").(int)
	connectTimeout := d.Get("connect_timeout").(int)
	healthyThreshold := d.Get("healthy_threshold").(int)
	unhealthyThreshold := d.Get("unhealthy_threshold").(int)

	if connectTimeout >= interval {
		return errors.New("connect_timeout must be less than interval")
	}
	clientIPMethod := d.Get("client_ip_method").(int)

	var realservers []gaapRealserverBind
	if raw, ok := d.GetOk("realserver_bind_set"); ok {
		list := raw.(*schema.Set).List()
		realservers = make([]gaapRealserverBind, 0, len(list))
		for _, v := range list {
			m := v.(map[string]interface{})

			if scheduler == "rr" && m["weight"].(int) != 1 {
				return errors.New("when scheduler is rr, realserver weight should be 1 or null")
			}

			realservers = append(realservers, gaapRealserverBind{
				id:     m["id"].(string),
				ip:     m["ip"].(string),
				port:   m["port"].(int),
				weight: m["weight"].(int),
			})
		}
	}

	var (
		id  string
		err error
	)

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	switch protocol {
	case "TCP":
		id, err = service.CreateTCPListener(ctx, name, scheduler, realserverType, proxyId, port, interval, connectTimeout, clientIPMethod, healthyThreshold, unhealthyThreshold, healthCheck)
		if err != nil {
			return err
		}

	case "UDP":
		optionalParams := make(map[string]interface{})
		if v, ok := d.GetOk("check_type"); ok {
			optionalParams["check_type"] = v.(string)
		} else {
			if healthCheck {
				return errors.New("Must set check_type when enable health_check.")
			}
		}
		if v, ok := d.GetOk("check_port"); ok {
			optionalParams["check_port"] = v.(int)
		}
		if v, ok := d.GetOk("context_type"); ok {
			optionalParams["context_type"] = v.(string)
		}
		if v, ok := d.GetOk("send_context"); ok {
			optionalParams["send_context"] = v.(string)
		}
		if v, ok := d.GetOk("recv_context"); ok {
			optionalParams["recv_context"] = v.(string)
		}
		id, err = service.CreateUDPListener(ctx, name, scheduler, realserverType, proxyId, port, interval, connectTimeout, healthyThreshold, unhealthyThreshold, healthCheck, optionalParams)
		if err != nil {
			return err
		}
	}

	d.SetId(id)

	if len(realservers) > 0 {
		if err := service.BindLayer4ListenerRealservers(ctx, id, protocol, proxyId, realservers); err != nil {
			return err
		}
	}

	return resourceTencentCloudGaapLayer4ListenerRead(d, m)
}

func resourceTencentCloudGaapLayer4ListenerRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_layer4_listener.read")()
	defer tccommon.InconsistentCheck(d, m)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()

	var (
		protocol           string
		name               *string
		port               *uint64
		scheduler          *string
		realServerType     *string
		healthCheck        *bool
		interval           *uint64
		connectTimeout     *uint64
		status             *uint64
		createTime         string
		realservers        []map[string]interface{}
		clientIpMethod     *uint64
		proxyId            *string
		healthyThreshold   *uint64
		unhealthyThreshold *uint64
		checkType          *string
		checkPort          *int64
		contextType        *string
		sendContext        *string
		recvContext        *string
	)

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	tcpListeners, err := service.DescribeTCPListeners(ctx, nil, &id, nil, nil)
	if err != nil {
		return err
	}

	udpListeners, err := service.DescribeUDPListeners(ctx, nil, &id, nil, nil)
	if err != nil {
		return err
	}

	switch {
	case len(tcpListeners) > 0:
		protocol = "TCP"

		listener := tcpListeners[0]
		clientIpMethod = listener.ClientIPMethod
		name = listener.ListenerName
		port = listener.Port
		scheduler = listener.Scheduler
		realServerType = listener.RealServerType
		proxyId = listener.ProxyId

		if listener.HealthCheck == nil {
			return errors.New("listener health check is nil")
		}
		healthCheck = helper.Bool(*listener.HealthCheck == 1)

		interval = listener.DelayLoop
		connectTimeout = listener.ConnectTimeout
		healthyThreshold = listener.HealthyThreshold
		unhealthyThreshold = listener.UnhealthyThreshold

		if len(listener.RealServerSet) > 0 {
			realservers = make([]map[string]interface{}, 0, len(listener.RealServerSet))
			for _, rs := range listener.RealServerSet {
				realservers = append(realservers, map[string]interface{}{
					"id":     rs.RealServerId,
					"ip":     rs.RealServerIP,
					"port":   rs.RealServerPort,
					"weight": rs.RealServerWeight,
				})
			}
		}

		status = listener.ListenerStatus

		if listener.CreateTime == nil {
			return errors.New("listener create time is nil")
		}
		createTime = helper.FormatUnixTime(*listener.CreateTime)

	case len(udpListeners) > 0:
		protocol = "UDP"

		listener := udpListeners[0]

		name = listener.ListenerName
		port = listener.Port
		scheduler = listener.Scheduler
		realServerType = listener.RealServerType
		proxyId = listener.ProxyId

		healthCheck = helper.Bool(*listener.HealthCheck == 1)
		connectTimeout = listener.ConnectTimeout
		interval = listener.DelayLoop
		healthyThreshold = listener.HealthyThreshold
		unhealthyThreshold = listener.UnhealthyThreshold
		checkType = listener.CheckType
		checkPort = listener.CheckPort
		contextType = listener.ContextType
		sendContext = listener.SendContext
		recvContext = listener.RecvContext
		if len(listener.RealServerSet) > 0 {
			realservers = make([]map[string]interface{}, 0, len(listener.RealServerSet))
			for _, rs := range listener.RealServerSet {
				realservers = append(realservers, map[string]interface{}{
					"id":     rs.RealServerId,
					"ip":     rs.RealServerIP,
					"port":   rs.RealServerPort,
					"weight": rs.RealServerWeight,
				})
			}
		}

		status = listener.ListenerStatus

		if listener.CreateTime == nil {
			return errors.New("listener create time is nil")
		}
		createTime = helper.FormatUnixTime(*listener.CreateTime)

	default:
		d.SetId("")
		return nil
	}

	_ = d.Set("protocol", protocol)
	_ = d.Set("name", name)
	_ = d.Set("port", port)
	_ = d.Set("scheduler", scheduler)
	_ = d.Set("realserver_type", realServerType)
	_ = d.Set("health_check", healthCheck)
	_ = d.Set("interval", interval)
	_ = d.Set("connect_timeout", connectTimeout)
	_ = d.Set("realserver_bind_set", realservers)
	_ = d.Set("status", status)
	_ = d.Set("create_time", createTime)
	_ = d.Set("proxy_id", proxyId)
	_ = d.Set("healthy_threshold", healthyThreshold)
	_ = d.Set("unhealthy_threshold", unhealthyThreshold)
	if protocol == "TCP" {
		_ = d.Set("client_ip_method", clientIpMethod)
	}
	if protocol == "UDP" {
		_ = d.Set("check_type", checkType)
		_ = d.Set("check_port", checkPort)
		_ = d.Set("context_type", contextType)
		_ = d.Set("send_context", sendContext)
		_ = d.Set("recv_context", recvContext)
	}
	return nil
}

func resourceTencentCloudGaapLayer4ListenerUpdate(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_layer4_listener.update")()
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	protocol := d.Get("protocol").(string)
	proxyId := d.Get("proxy_id").(string)
	var (
		name               *string
		scheduler          *string
		healthCheck        *bool
		interval           int
		connectTimeout     int
		attrChange         []string
		healthyThreshold   int
		unhealthyThreshold int
	)

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	d.Partial(true)

	if d.HasChange("name") {
		attrChange = append(attrChange, "name")
	}
	name = helper.String(d.Get("name").(string))

	if d.HasChange("scheduler") {
		attrChange = append(attrChange, "scheduler")
	}
	scheduler = helper.String(d.Get("scheduler").(string))

	if d.HasChange("health_check") {
		attrChange = append(attrChange, "health_check")
	}
	healthCheck = helper.Bool(d.Get("health_check").(bool))

	if d.HasChange("interval") {
		attrChange = append(attrChange, "interval")
	}
	interval = d.Get("interval").(int)

	if d.HasChange("connect_timeout") {
		attrChange = append(attrChange, "connect_timeout")
	}
	connectTimeout = d.Get("connect_timeout").(int)

	if d.HasChange("healthy_threshold") {
		attrChange = append(attrChange, "healthy_threshold")
	}
	healthyThreshold = d.Get("healthy_threshold").(int)

	if d.HasChange("unhealthy_threshold") {
		attrChange = append(attrChange, "unhealthy_threshold")
	}
	unhealthyThreshold = d.Get("unhealthy_threshold").(int)

	optionalParams := make(map[string]interface{})
	if d.HasChange("check_type") {
		attrChange = append(attrChange, "check_type")
	}
	if d.HasChange("check_port") {
		attrChange = append(attrChange, "check_port")
	}
	if d.HasChange("context_type") {
		attrChange = append(attrChange, "context_type")
	}
	if d.HasChange("send_context") {
		attrChange = append(attrChange, "send_context")
	}
	if d.HasChange("recv_context") {
		attrChange = append(attrChange, "recv_context")
	}
	if v, ok := d.GetOk("check_type"); ok {
		optionalParams["check_type"] = v.(string)
	}
	if v, ok := d.GetOk("check_port"); ok {
		optionalParams["check_port"] = v.(int)
	}
	if v, ok := d.GetOk("context_type"); ok {
		optionalParams["context_type"] = v.(string)
	}
	if v, ok := d.GetOk("send_context"); ok {
		optionalParams["send_context"] = v.(string)
	}
	if v, ok := d.GetOk("recv_context"); ok {
		optionalParams["recv_context"] = v.(string)
	}
	// only check for TCP listener
	if protocol == "TCP" && connectTimeout >= interval {
		return errors.New("connect_timeout must be less than interval")
	}

	if len(attrChange) > 0 {
		switch protocol {
		case "TCP":
			if err := service.ModifyTCPListenerAttribute(ctx, proxyId, id, name, scheduler, healthCheck, interval, connectTimeout, healthyThreshold, unhealthyThreshold); err != nil {
				return err
			}

		case "UDP":
			if err := service.ModifyUDPListenerAttribute(ctx, proxyId, id, name, scheduler, connectTimeout, interval, healthyThreshold, unhealthyThreshold, healthCheck, optionalParams); err != nil {
				return err
			}
		}
	}

	if d.HasChange("realserver_bind_set") {
		list := d.Get("realserver_bind_set").(*schema.Set).List()
		realservers := make([]gaapRealserverBind, 0, len(list))
		for _, v := range list {
			m := v.(map[string]interface{})
			realservers = append(realservers, gaapRealserverBind{
				id:     m["id"].(string),
				ip:     m["ip"].(string),
				port:   m["port"].(int),
				weight: m["weight"].(int),
			})
		}

		if err := service.BindLayer4ListenerRealservers(ctx, id, protocol, proxyId, realservers); err != nil {
			return err
		}

	}

	d.Partial(false)

	return resourceTencentCloudGaapLayer4ListenerRead(d, m)
}

func resourceTencentCloudGaapLayer4ListenerDelete(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_gaap_layer4_listener.delete")()
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	protocol := d.Get("protocol").(string)
	proxyId := d.Get("proxy_id").(string)

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	return service.DeleteLayer4Listener(ctx, id, proxyId, protocol)
}
