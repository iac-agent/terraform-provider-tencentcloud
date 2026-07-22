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
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClbListenerRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbListenerRuleCreate,
		Read:   resourceTencentCloudClbListenerRuleRead,
		Update: resourceTencentCloudClbListenerRuleUpdate,
		Delete: resourceTencentCloudClbListenerRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CLB监听器ID。",
			},
			"clb_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CLB实例ID。",
			},
			"domain": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"domains"},
				ExactlyOneOf:  []string{"domain", "domains"},
				Description: "监听规则的域名。单域规则传递到“域名”，多域规则传递到“domains”。",
			},
			"domains": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"domain"},
				ExactlyOneOf:  []string{"domain", "domains"},
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description: "监听规则的域名列表。单域规则传递到“域名”，多域规则传递到“domains”。",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "侦听器规则的 URL。",
			},
			"health_check_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "是否开启健康检查。",
			},
			"health_check_interval_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 300),
				Description: "健康检查的间隔时间。有效值范围：（2~300）秒。默认值为“5”秒。注意：TCP/UDP/TCP_SSL监听可以直接配置，HTTP/HTTPS监听需要在tencentcloud_clb_listener_rule中配置。",
			},
			"health_check_health_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 10),
				Description: "健康检查的健康阈值，默认为`3`。如果连续3次健康检查返回成功结果，则说明转发正常。取值范围为[2-10]。注意：TCP/UDP/TCP_SSL监听可以直接配置，HTTP/HTTPS监听需要在tencentcloud_clb_listener_rule中配置。",
			},
			"health_check_unhealth_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 10),
				Description: "健康检查不健康阈值，默认为`3`。如果连续3次返回不健康结果，则表明转发异常。取值范围为[2-10]。  注意：TCP/UDP/TCP_SSL监听可以直接配置，HTTP/HTTPS监听需要在tencentcloud_clb_listener_rule中配置。",
			},
			"health_check_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "自定义检测相关参数。健康检查端口，默认为后端服务的端口，除非要指定特定端口，建议留空。 （仅适用于 TCP/UDP 侦听器）。",
			},
			"health_check_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(HEALTH_CHECK_TYPE),
				Description: "健康检查类型。有效值为“CUSTOM”、“PING”、“TCP”、“HTTP”、“HTTPS”、“GRPC”、“GRPCS”。",
			},
			"health_check_time_out": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 60),
				Description: "健康检查超时。取值范围为[2-60]（秒）。",
			},
			"health_check_http_code": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 31),
				Description: "HTTP 状态代码。默认值为 31。有效值范围：[1~31]。 `1表示返回值'1xx'是健康。 ‘2’表示返回值‘2xx’是健康状况。 ‘4’表示返回值‘3xx’是健康状况。 ‘8’表示返回值‘4xx’是健康状况。 16表示返回值“5xx”是健康。如果想要多个返回码来指示健康状况，需要添加相应的值。注意：“TCP”监听器的“HTTP”健康检查仅支持指定一个健康检查状态码。注意：仅支持“HTTP”和“HTTPS”协议的侦听器。",
			},
			"health_check_http_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "健康检查路径。注意：仅支持“HTTP”和“HTTPS”协议的监听器。",
			},
			"health_check_http_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "健康检查域名。注意：仅支持“HTTP”和“HTTPS”协议的监听器。",
			},
			"health_check_http_method": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CLB_HTTP_METHOD),
				Description: "健康检查方法。注意：仅支持“HTTP”和“HTTPS”协议的监听器。默认为“HEAD”，可用值为“HEAD”和“GET”。",
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
				ForceNew:      true,
				ConflictsWith: []string{"multi_cert_info"},
				ValidateFunc:  tccommon.ValidateAllowedStringValue(CERT_SSL_MODE),
				Description: "证书类型。有效值：“单向”、“相互”。注意：仅支持HTTPS协议的监听。",
			},
			"certificate_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"multi_cert_info"},
				Description: "服务器证书的 ID。注意：仅支持HTTPS协议的监听。",
			},
			"certificate_ca_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"multi_cert_info"},
				Description: "客户端证书的ID。注意：仅支持HTTPS协议的监听。",
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
							ForceNew:     true,
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
				Description: "CLB 侦听器中的会话持续时间。注意：当调度程序指定为“WRR”时可用，当侦听器协议为“TCP_SSL”时不可用。  注意：TCP/UDP/TCP_SSL监听可以直接配置，HTTP/HTTPS监听需要在tencentcloud_clb_listener_rule中配置。",
			},
			"http2_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "是否应用HTTP2.0协议。",
			},
			"scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CLB_LISTENER_SCHEDULER),
				Description: "CLB监听规则的调度方法。有效值：“WRR”、“IP HASH”、“LEAST_CONN”。当“target_type”不是“TARGETGROUP-V2”时，默认值为“WRR”。  注意：TCP/UDP/TCP_SSL监听可以直接配置，HTTP/HTTPS监听需要在tencentcloud_clb_listener_rule中配置。",
			},
			"target_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      CLB_TARGET_TYPE_NODE,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{CLB_TARGET_TYPE_NODE, CLB_TARGET_TYPE_TARGETGROUP, CLB_TARGET_TYPE_TARGETGROUP_V2}),
				Description: "后端目标类型。有效值：“NODE”、“TARGETGROUP”、“TARGETGROUP-V2”。 `NODE` 表示绑定普通节点，`TARGETGROUP` 表示绑定目标组。",
			},
			"forward_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"HTTP", "HTTPS", "GRPC", "GRPCS", "TRPC"}),
				Description: "CLB实例与真实服务器之间的转发协议。有效值：“HTTP”、“HTTPS”、“GRPC”、“GRPCS”、“TRPC”。默认为“HTTP”。",
			},
			"quic": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "是否启用QUIC。注意：QUIC只能对HTTPS域名启用。",
			},
			"oauth": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "OAuth 配置信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"oauth_enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "启用或禁用身份验证。 True：启用；错误：禁用。",
						},
						"oauth_failure_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "所有 IAP 失败后，请求将被拒绝或释放。绕过：通过；拒绝：拒绝。",
						},
					},
				},
			},
			//computed
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "该CLB监听规则ID。",
			},
		},
	}
}

func resourceTencentCloudClbListenerRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_listener_rule.create")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	listenerId := d.Get("listener_id").(string)
	checkErr := ListenerIdCheck(listenerId)
	if checkErr != nil {
		return checkErr
	}
	clbId := d.Get("clb_id").(string)
	protocol := ""
	//get listener protocol
	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, e := clbService.DescribeListenerById(ctx, listenerId, clbId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if instance != nil {
			protocol = *(instance.Protocol)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s get CLB listener failed, reason:%+v", logId, err)
		return err
	}

	if !(protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
		return fmt.Errorf("[CHECK][CLB listener rule][Create] check: The rule can only be created/modified with listeners of protocol HTTP/HTTPS")
	}
	request := clb.NewCreateRuleRequest()
	request.LoadBalancerId = helper.String(clbId)
	request.ListenerId = helper.String(listenerId)

	//rule set
	var (
		rule    clb.RuleInput
		domain  string
		domains []*string
	)

	if v, ok := d.GetOk("domain"); ok {
		rule.Domain = helper.String(v.(string))
		domain = v.(string)
	}

	if v, ok := d.GetOk("domains"); ok {
		tmpDomains := v.(*schema.Set).List()
		domains = make([]*string, 0, len(tmpDomains))
		for _, value := range tmpDomains {
			tmpDomain := value.(string)
			domains = append(domains, &tmpDomain)
		}

		rule.Domains = domains
	}

	url := d.Get("url").(string)
	rule.Url = helper.String(url)
	rule.TargetType = helper.String(d.Get("target_type").(string))
	if v, ok := d.GetOk("forward_type"); ok {
		rule.ForwardType = helper.String(v.(string))
	}
	scheduler := ""
	targetType := d.Get("target_type").(string)
	if v, ok := d.GetOk("scheduler"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
			return fmt.Errorf("[CHECK][CLB listener rule][Create] check: Scheduler can only be set with listener protocol TCP/UDP/TCP_SSL or rule of listener HTTP/HTTPS")
		}

		scheduler = v.(string)
		rule.Scheduler = helper.String(scheduler)
	} else {
		// Apply default value WRR only when target_type is not TARGETGROUP-V2
		if targetType != CLB_TARGET_TYPE_TARGETGROUP_V2 {
			scheduler = CLB_LISTENER_SCHEDULER_WRR
			rule.Scheduler = helper.String(scheduler)
		}
	}

	if v, ok := d.GetOkExists("session_expire_time"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
			return fmt.Errorf("[CHECK][CLB listener rule][Create] check: session_expire_time can only be set with protocol TCP/UDP or rule of listener HTTP/HTTPS")
		}
		if scheduler != CLB_LISTENER_SCHEDULER_WRR && scheduler != "" {
			return fmt.Errorf("[CHECK][CLB listener rule][Create] check: session_expire_time can only be set when scheduler is WRR ")
		}
		vv := int64(v.(int))
		rule.SessionExpireTime = &vv
	}
	healthSetFlag, healthCheck, healthErr := checkHealthCheckPara(ctx, d, protocol, HEALTH_APPLY_TYPE_RULE)
	if healthErr != nil {
		return healthErr
	}
	if healthSetFlag {
		rule.HealthCheck = healthCheck
	}

	certificateSetFlag, certificateInput, certErr := checkCertificateInputPara(ctx, d, meta)
	if certErr != nil {
		return certErr
	}
	if certificateSetFlag {
		if !(protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
			return fmt.Errorf("[CHECK][CLB listener rule][Create] check: certificate para can only be set with rule of linstener with protocol 'HTTPS'")
		}
		rule.Certificate = certificateInput
	}

	multiCertificateSetFlag, multiCertInput, certErr := checkMultiCertificateInputPara(ctx, d, meta)
	if certErr != nil {
		return certErr
	}

	if multiCertificateSetFlag {
		rule.MultiCertInfo = multiCertInput
	} else {
		if protocol == CLB_LISTENER_PROTOCOL_TCPSSL {
			return fmt.Errorf("[CHECK][CLB listener][Create] check: certificated need to be set when protocol is HTTPS")
		}
	}

	if v, ok := d.GetOkExists("quic"); ok {
		rule.Quic = helper.Bool(v.(bool))
	}

	request.Rules = []*clb.RuleInput{&rule}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		requestId := ""
		response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().CreateRule(request)
		if e != nil {
			if err := processRetryErrMsg(e); err != nil {
				return err
			}
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			if response == nil || response.Response == nil || response.Response.RequestId == nil {
				return resource.NonRetryableError(fmt.Errorf("Create CLB listener rule failed, Response is nil."))
			}

			requestId = *response.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
			if retryErr != nil {
				return resource.NonRetryableError(errors.WithStack(retryErr))
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create CLB listener rule failed, reason:%+v", logId, err)
		return err
	}

	locationId := ""
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ruleInstance, ruleErr := clbService.DescribeRuleByPara(ctx, clbId, listenerId, domain, url, domains)
		if ruleErr != nil {
			return tccommon.RetryError(errors.WithStack(ruleErr))
		}

		if ruleInstance == nil || ruleInstance.LocationId == nil {
			return resource.NonRetryableError(fmt.Errorf("read CLB listener rule failed, Response is nil."))
		}

		locationId = *ruleInstance.LocationId
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB listener rule failed, reason:%+v", logId, err)
		return err
	}

	//this ID style changes since terraform 1.47.0
	d.SetId(strings.Join([]string{clbId, listenerId, locationId}, tccommon.FILED_SP))

	// set http2
	if v, ok := d.GetOkExists("http2_switch"); ok {
		http2Switch := v.(bool)
		domainRequest := clb.NewModifyDomainAttributesRequest()
		domainRequest.Http2 = &http2Switch
		domainRequest.LoadBalancerId = &clbId
		domainRequest.ListenerId = &listenerId
		if domain != "" {
			domainRequest.Domain = &domain
		} else {
			domainRequest.NewDomains = domains
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyDomainAttributes(domainRequest)
			if e != nil {
				if sdkError, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkError.Code == "FailedOperation.ResourceInOperating" {
						return resource.RetryableError(e)
					}
				}

				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
				if response == nil || response.Response == nil || response.Response.RequestId == nil {
					return resource.NonRetryableError(fmt.Errorf("Modify domain attributes failed, Response is nil."))
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
			log.Printf("[CRITAL]%s update CLB listener rule failed, reason:%+v", logId, err)
			return err
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "oauth"); ok {
		modifyRuleRequest := clb.NewModifyRuleRequest()
		modifyRuleRequest.ListenerId = helper.String(listenerId)
		modifyRuleRequest.LoadBalancerId = helper.String(clbId)
		modifyRuleRequest.LocationId = helper.String(locationId)
		oauth := &clb.OAuth{}
		if v, ok := dMap["oauth_enable"]; ok {
			oauth.OAuthEnable = helper.Bool(v.(bool))
		}
		if v, ok := dMap["oauth_failure_status"]; ok {
			oauth.OAuthFailureStatus = helper.String(v.(string))
		}
		modifyRuleRequest.OAuth = oauth
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyRule(modifyRuleRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
				if response == nil || response.Response == nil || response.Response.RequestId == nil {
					return resource.NonRetryableError(fmt.Errorf("Modify rule failed, Response is nil."))
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
			log.Printf("[CRITAL]%s update CLB listener rule failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudClbListenerRuleRead(d, meta)
}

func resourceTencentCloudClbListenerRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_listener_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	resourceId := d.Id()
	var locationId = resourceId
	items := strings.Split(resourceId, tccommon.FILED_SP)
	itemLength := len(items)
	clbId := d.Get("clb_id").(string)
	listenerId := d.Get("listener_id").(string)
	checkErr := ListenerIdCheck(listenerId)
	if checkErr != nil {
		return checkErr
	}
	if itemLength == 1 && clbId == "" {
		return fmt.Errorf("The old style listenerId %s does not support import, please use clbId#listenerId style", resourceId)
	} else if itemLength == 3 && clbId == "" {
		locationId = items[2]
		listenerId = items[1]
		clbId = items[0]
	} else if itemLength == 3 && clbId != "" {
		locationId = items[2]
		listenerId = items[1]
	} else if itemLength != 1 && itemLength != 3 {
		return fmt.Errorf("broken ID %s", resourceId)
	}

	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	//this function is not supported by api, need to be travelled
	filter := map[string]string{"rule_id": locationId, "listener_id": listenerId, "clb_id": clbId}
	var instances []*clb.RuleOutput
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := clbService.DescribeRulesByFilter(ctx, filter)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instances = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB listener rule failed, reason:%+v", logId, err)
		return err
	}

	if len(instances) == 0 {
		d.SetId("")
		return nil
	}

	instance := instances[0]
	_ = d.Set("clb_id", clbId)
	_ = d.Set("listener_id", listenerId)
	if instance.Domain != nil {
		_ = d.Set("domain", instance.Domain)
	}

	if instance.Domains != nil {
		_ = d.Set("domains", helper.StringsInterfaces(instance.Domains))
	}

	_ = d.Set("rule_id", instance.LocationId)
	_ = d.Set("url", instance.Url)
	_ = d.Set("scheduler", instance.Scheduler)
	_ = d.Set("session_expire_time", instance.SessionExpireTime)
	_ = d.Set("target_type", instance.TargetType)
	_ = d.Set("forward_type", instance.ForwardType)
	_ = d.Set("http2_switch", instance.Http2)

	if instance.QuicStatus != nil {
		if *instance.QuicStatus == "QUIC_ACTIVE" {
			_ = d.Set("quic", true)
		} else {
			_ = d.Set("quic", false)
		}
	}

	//health check
	if instance.HealthCheck != nil {
		health_check_switch := false
		if *instance.HealthCheck.HealthSwitch == int64(1) {
			health_check_switch = true
		}
		_ = d.Set("health_check_switch", health_check_switch)
		_ = d.Set("health_check_interval_time", instance.HealthCheck.IntervalTime)
		_ = d.Set("health_check_health_num", instance.HealthCheck.HealthNum)
		_ = d.Set("health_check_unhealth_num", instance.HealthCheck.UnHealthNum)
		_ = d.Set("health_check_http_method", helper.String(strings.ToUpper(*instance.HealthCheck.HttpCheckMethod)))
		_ = d.Set("health_check_http_domain", instance.HealthCheck.HttpCheckDomain)
		_ = d.Set("health_check_http_path", instance.HealthCheck.HttpCheckPath)
		_ = d.Set("health_check_http_code", instance.HealthCheck.HttpCode)
		if instance.HealthCheck.CheckPort != nil {
			_ = d.Set("health_check_port", instance.HealthCheck.CheckPort)
		}
		_ = d.Set("health_check_type", instance.HealthCheck.CheckType)
		_ = d.Set("health_check_time_out", instance.HealthCheck.TimeOut)
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

	if instance.OAuth != nil {
		oath := make(map[string]interface{})
		if instance.OAuth.OAuthEnable != nil {
			oath["oauth_enable"] = instance.OAuth.OAuthEnable
		}
		if instance.OAuth.OAuthFailureStatus != nil {
			oath["oauth_failure_status"] = instance.OAuth.OAuthFailureStatus
		}
		_ = d.Set("oauth", []interface{}{oath})
	}

	return nil
}

func resourceTencentCloudClbListenerRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_listener_rule.update")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	resourceId := d.Id()
	items := strings.Split(resourceId, tccommon.FILED_SP)
	itemLength := len(items)
	locationId := items[itemLength-1]
	listenerId := d.Get("listener_id").(string)
	checkErr := ListenerIdCheck(listenerId)
	if checkErr != nil {
		return checkErr
	}
	clbId := d.Get("clb_id").(string)
	protocol := ""
	//get listener protocol
	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		instance, e := clbService.DescribeListenerById(ctx, listenerId, clbId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		protocol = *(instance.Protocol)
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s get CLB listener failed, reason:%s\n ", logId, err.Error())
		return err
	}

	changed := false
	url := ""
	scheduler := ""
	sessionExpireTime := 0

	request := clb.NewModifyRuleRequest()
	request.ListenerId = helper.String(listenerId)
	request.LoadBalancerId = helper.String(clbId)
	request.LocationId = helper.String(locationId)
	if d.HasChange("url") {
		changed = true
		url = d.Get("url").(string)
		request.Url = helper.String(url)
	}
	if d.HasChange("oauth") {
		changed = true
		if dMap, ok := helper.InterfacesHeadMap(d, "oauth"); ok {
			oauth := &clb.OAuth{}
			if v, ok := dMap["oauth_enable"]; ok {
				oauth.OAuthEnable = helper.Bool(v.(bool))
			}
			if v, ok := dMap["oauth_failure_status"]; ok {
				oauth.OAuthFailureStatus = helper.String(v.(string))
			}
			request.OAuth = oauth
		}
	}

	if d.HasChange("forward_type") {
		changed = true
		request.ForwardType = helper.String(d.Get("forward_type").(string))
	}

	if d.HasChange("scheduler") {
		changed = true
		scheduler = d.Get("scheduler").(string)
		targetType := d.Get("target_type").(string)
		if !(protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
			return fmt.Errorf("[CHECK][CLB listener rule %s][Update] check: Scheduler can only be set with listener protocol TCP/UDP/TCP_SSL or rule of listener HTTP/HTTPS", locationId)
		}
		// Only set scheduler when target_type is not TARGETGROUP-V2
		if targetType != CLB_TARGET_TYPE_TARGETGROUP_V2 {
			request.Scheduler = helper.String(scheduler)
		}
	}

	if d.HasChange("session_expire_time") {
		changed = true
		sessionExpireTime = d.Get("session_expire_time").(int)
		if !(protocol == CLB_LISTENER_PROTOCOL_HTTP || protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
			return fmt.Errorf("[CHECK][CLB listener rule %s][Update] check: session_expire_time can only be set with protocol TCP/UDP or rule of listener HTTP/HTTPS", locationId)
		}
		if scheduler != CLB_LISTENER_SCHEDULER_WRR && scheduler != "" {
			return fmt.Errorf("[CHECK][CLB listener rule %s][Update] check: session_expire_time can only be set when scheduler is WRR", locationId)
		}
		sessionExpireTime64 := int64(sessionExpireTime)
		request.SessionExpireTime = &sessionExpireTime64
	}

	healthSetFlag, healthCheck, healthErr := checkHealthCheckPara(ctx, d, protocol, HEALTH_APPLY_TYPE_RULE)
	if healthErr != nil {
		return healthErr
	}

	if healthSetFlag {
		changed = true
		request.HealthCheck = healthCheck
	}

	if changed {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyRule(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
				if response == nil || response.Response == nil || response.Response.RequestId == nil {
					return resource.NonRetryableError(fmt.Errorf("Modify rule failed, Response is nil."))
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
			log.Printf("[CRITAL]%s update CLB listener rule failed, reason:%+v", logId, err)
			return err
		}
	}

	//modify domain and ssl
	domainChanged := false
	domainRequest := clb.NewModifyDomainAttributesRequest()
	if d.HasChange("domain") {
		old, new := d.GetChange("domain")
		domainChanged = true
		domainRequest.Domain = helper.String(old.(string))
		domainRequest.NewDomain = helper.String(new.(string))
	} else if d.HasChange("domains") {
		old, new := d.GetChange("domains")
		domainChanged = true
		oldDomains := old.(*schema.Set).List()
		newDomains := new.(*schema.Set).List()

		if len(oldDomains) < 1 || len(newDomains) < 1 {
			return fmt.Errorf("Params `domains` cant not be empty.")
		}

		domainRequest.Domain = helper.String(oldDomains[0].(string))
		tmpDomains := make([]*string, 0, len(newDomains))
		for _, value := range newDomains {
			domain := value.(string)
			tmpDomains = append(tmpDomains, &domain)
		}

		domainRequest.NewDomains = tmpDomains
	} else {
		domainRequest.Domain = helper.String(d.Get("domain").(string))
	}

	if d.HasChange("certificate_id") || d.HasChange("certificate_ca_id ") || d.HasChange("certificate_ssl_mode") {
		domainChanged = true
		certificateSetFlag, certificateInput, certErr := checkCertificateInputPara(ctx, d, meta)
		if certErr != nil {
			return certErr
		}
		if certificateSetFlag {
			if !(protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
				return fmt.Errorf("[CHECK][CLB listener rule][Create] check: certificate para can only be set with rule of linstener with protocol 'HTTPS'")
			}
			domainRequest.Certificate = certificateInput
		}
	}

	if d.HasChange("multi_cert_info") {
		domainChanged = true
		multiCertificateSetFlag, multiCertInput, certErr := checkMultiCertificateInputPara(ctx, d, meta)
		if certErr != nil {
			return certErr
		}

		if multiCertificateSetFlag {
			domainRequest.MultiCertInfo = multiCertInput
		} else {
			if protocol == CLB_LISTENER_PROTOCOL_TCPSSL {
				return fmt.Errorf("[CHECK][CLB listener][Create] check: certificated need to be set when protocol is HTTPS")
			}
		}
	}

	if d.HasChange("http2_switch") {
		if v, ok := d.GetOkExists("http2_switch"); ok {
			if !(protocol == CLB_LISTENER_PROTOCOL_HTTPS) {
				return fmt.Errorf("[CHECK][CLB listener rule][Create] check: certificate para can only be set with rule of linstener with protocol 'HTTPS'")
			}
			domainChanged = true
			domainRequest.Http2 = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("quic") {
		domainChanged = true
		if v, ok := d.GetOkExists("quic"); ok {
			domainRequest.Quic = helper.Bool(v.(bool))
		}
	}

	if domainChanged {
		domainRequest.ListenerId = &listenerId
		domainRequest.LoadBalancerId = &clbId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyDomainAttributes(domainRequest)
			if e != nil {
				if err := processRetryErrMsg(e); err != nil {
					return err
				}
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, domainRequest.GetAction(), domainRequest.ToJsonString(), response.ToJsonString())
				if response == nil || response.Response == nil || response.Response.RequestId == nil {
					return resource.NonRetryableError(fmt.Errorf("Modify domain attributes failed, Response is nil."))
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
			log.Printf("[CRITAL]%s update CLB listener rule failed, reason:%+v", logId, err)
			return err
		}
	}

	return nil
}

func resourceTencentCloudClbListenerRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_listener_rule.delete")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	resourceId := d.Id()
	items := strings.Split(resourceId, tccommon.FILED_SP)
	itemLength := len(items)
	locationId := items[itemLength-1]
	listenerId := d.Get("listener_id").(string)
	checkErr := ListenerIdCheck(listenerId)
	if checkErr != nil {
		return checkErr
	}
	clbId := d.Get("clb_id").(string)

	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := clbService.DeleteRuleById(ctx, clbId, listenerId, locationId)
		if e != nil {
			if err := processRetryErrMsg(e); err != nil {
				return err
			}
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CLB listener rule failed, reason:%+v", logId, err)
		return err
	}
	return nil
}
