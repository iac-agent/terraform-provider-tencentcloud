package alb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	alb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAlbListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAlbListenerCreate,
		Read:   resourceTencentCloudAlbListenerRead,
		Update: resourceTencentCloudAlbListenerUpdate,
		Delete: resourceTencentCloudAlbListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the ALB load balancer instance, formatted as `alb-` followed by 8 alphanumeric characters.",
			},
			"listener_protocol": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Listener protocol. Valid values: `HTTP`, `HTTPS`, `QUIC`.",
			},
			"listener_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 65535),
				Description:  "Front-end port of the load balancer instance. Valid value range: [1-65535].",
			},
			"default_actions": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Default forwarding rule action list. Currently only 1 default forwarding rule action is supported per listener.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Forwarding action type. Only forwarding to target group is supported when creating a listener.",
						},
						"target_group_config": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "Forwarding target group configuration. Only a single target group is supported when creating a listener.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_groups": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Target group list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"target_group_id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Target group ID, formatted as `lbtg-` followed by 8 alphanumeric characters.",
												},
												"weight": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: tccommon.ValidateIntegerInRange(0, 100),
													Description:  "Weight, valid value range: [0-100]. Default to 10.",
												},
											},
										},
									},
									"target_group_sticky_session": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Sticky session across target groups.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:        schema.TypeBool,
													Required:    true,
													Description: "Whether to enable sticky session. Default to false.",
												},
												"timeout": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: tccommon.ValidateIntegerInRange(1, 86400),
													Description:  "Timeout in seconds. Valid value range: [1-86400]. Default to 1000.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"listener_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom listener name. Length of 1-255 characters, can contain Chinese characters, letters, digits, hyphen (-), slash (/), period (.) and underscore (_).",
			},
			"certificate_ids": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"ca_certificate_ids"},
				Description:   "Server certificate ID list.",
				Elem:          &schema.Schema{Type: schema.TypeString},
			},
			"ca_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable mutual authentication. Default to false.",
			},
			"ca_certificate_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "CA certificate ID list of the listener. Currently only 1 CA certificate is supported. Required when `ca_enabled` is true.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Security policy ID, formatted as `tls-` followed by 8 alphanumeric characters.",
			},
			"idle_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 600),
				Description:  "Connection idle timeout in seconds. Valid value range: [1-600]. Default to 15.",
			},
			"request_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 600),
				Description:  "Connection request timeout in seconds. Valid value range: [1-600]. Default to 60.",
			},
			"gzip_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable Gzip compression. Default to true.",
			},
			"http2enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable HTTP/2. HTTP protocol defaults to false, HTTPS protocol defaults to true. Only supported by HTTPS protocol.",
			},
			"x_forwarded_for_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "X-Forwarded-For configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"x_forwarded_for_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies how to process the X-Forwarded-For header. Valid values: `append`, `remove`, `passthrough`.",
						},
						"x_forwarded_for_alb_id_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to obtain the load balancer instance ID via the ALB-ID header.",
						},
						"x_forwarded_for_client_src_port_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to obtain the client source port via the X-Forwarded-Client-srcport header.",
						},
						"x_forwarded_for_host_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to obtain the client domain via the X-Forwarded-Host header.",
						},
						"x_forwarded_for_port_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to obtain the listener port via the X-Forwarded-Port header.",
						},
						"x_forwarded_for_proto_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to obtain the listener protocol via the X-Forwarded-Proto header.",
						},
						"x_tencent_client_idn_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to access the issuer of the client certificate $ssl_client_i_dn via the X-Tencent-Client-IDN header.",
						},
						"x_tencent_client_sdn_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to access the subject of the client certificate $ssl_client_s_dn via the X-Tencent-Client-SDN header.",
						},
						"x_tencent_client_serial_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to access the serial number of the client certificate $ssl_client_serial via the X-Tencent-Client-Serial header.",
						},
						"x_tencent_client_verify_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to access the verification result of the client certificate $ssl_client_verify via the X-Tencent-Client-Verify header.",
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    20,
				Description: "Tag list. Up to 20 tags are supported.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},
			"client_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Client token used to guarantee request idempotency. Only ASCII characters are supported.",
			},
		},
	}
}

func buildAlbDefaultActions(d *schema.ResourceData) []*alb.DefaultAction {
	actions := make([]*alb.DefaultAction, 0, 1)
	actionList, ok := d.GetOk("default_actions")
	if !ok {
		return actions
	}
	for _, item := range actionList.([]interface{}) {
		actionMap, actionOk := item.(map[string]interface{})
		if !actionOk {
			continue
		}
		action := &alb.DefaultAction{}
		if v, ok := actionMap["type"]; ok {
			action.Type = helper.String(v.(string))
		}
		configList, configOk := actionMap["target_group_config"].([]interface{})
		if configOk && len(configList) > 0 {
			configMap, mOk := configList[0].(map[string]interface{})
			if mOk {
				config := &alb.TargetGroupConfig{}
				groupList, groupOk := configMap["target_groups"].([]interface{})
				if groupOk && len(groupList) > 0 {
					groups := make([]*alb.TargetGroupTuple, 0, len(groupList))
					for _, groupItem := range groupList {
						groupMap, gOk := groupItem.(map[string]interface{})
						if !gOk {
							continue
						}
						group := &alb.TargetGroupTuple{}
						if v, ok := groupMap["target_group_id"]; ok {
							group.TargetGroupId = helper.String(v.(string))
						}
						if v, ok := groupMap["weight"]; ok && v.(int) > 0 {
							group.Weight = helper.IntInt64(v.(int))
						}
						groups = append(groups, group)
					}
					config.TargetGroups = groups
				}
				stickyList, stickyOk := configMap["target_group_sticky_session"].([]interface{})
				if stickyOk && len(stickyList) > 0 {
					stickyMap, sOk := stickyList[0].(map[string]interface{})
					if sOk {
						sticky := &alb.TargetGroupStickySession{}
						if v, ok := stickyMap["enabled"]; ok {
							sticky.Enabled = helper.Bool(v.(bool))
						}
						if v, ok := stickyMap["timeout"]; ok && v.(int) > 0 {
							sticky.Timeout = helper.IntInt64(v.(int))
						}
						config.TargetGroupStickySession = sticky
					}
				}
				action.TargetGroupConfig = config
			}
		}
		actions = append(actions, action)
	}
	return actions
}

func buildAlbXForwardedForConfig(d *schema.ResourceData) *alb.XForwardedForConfig {
	configList, ok := d.GetOk("x_forwarded_for_config")
	if !ok {
		return nil
	}
	list := configList.([]interface{})
	if len(list) == 0 {
		return nil
	}
	configMap, mOk := list[0].(map[string]interface{})
	if !mOk {
		return nil
	}
	config := &alb.XForwardedForConfig{}
	if v, ok := configMap["x_forwarded_for_mode"]; ok && v.(string) != "" {
		config.XForwardedForMode = helper.String(v.(string))
	}
	if v, ok := configMap["x_forwarded_for_alb_id_enabled"]; ok {
		config.XForwardedForAlbIdEnabled = helper.Bool(v.(bool))
	}
	if v, ok := configMap["x_forwarded_for_client_src_port_enabled"]; ok {
		config.XForwardedForClientSrcPortEnabled = helper.Bool(v.(bool))
	}
	if v, ok := configMap["x_forwarded_for_host_enabled"]; ok {
		config.XForwardedForHostEnabled = helper.Bool(v.(bool))
	}
	if v, ok := configMap["x_forwarded_for_port_enabled"]; ok {
		config.XForwardedForPortEnabled = helper.Bool(v.(bool))
	}
	if v, ok := configMap["x_forwarded_for_proto_enabled"]; ok {
		config.XForwardedForProtoEnabled = helper.Bool(v.(bool))
	}
	if v, ok := configMap["x_tencent_client_idn_enabled"]; ok {
		config.XTencentClientIDNEnabled = helper.Bool(v.(bool))
	}
	if v, ok := configMap["x_tencent_client_sdn_enabled"]; ok {
		config.XTencentClientSDNEnabled = helper.Bool(v.(bool))
	}
	if v, ok := configMap["x_tencent_client_serial_enabled"]; ok {
		config.XTencentClientSerialEnabled = helper.Bool(v.(bool))
	}
	if v, ok := configMap["x_tencent_client_verify_enabled"]; ok {
		config.XTencentClientVerifyEnabled = helper.Bool(v.(bool))
	}
	return config
}

func buildAlbTags(d *schema.ResourceData) []*alb.TagInfo {
	tagList, ok := d.GetOk("tags")
	if !ok {
		return nil
	}
	tags := make([]*alb.TagInfo, 0, len(tagList.([]interface{})))
	for _, item := range tagList.([]interface{}) {
		tagMap, mOk := item.(map[string]interface{})
		if !mOk {
			continue
		}
		tag := &alb.TagInfo{}
		if v, ok := tagMap["tag_key"]; ok {
			tag.TagKey = helper.String(v.(string))
		}
		if v, ok := tagMap["tag_value"]; ok {
			tag.TagValue = helper.String(v.(string))
		}
		tags = append(tags, tag)
	}
	return tags
}

func resourceTencentCloudAlbListenerCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_listener.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := alb.NewCreateListenerRequest()
	request.LoadBalancerId = helper.String(d.Get("load_balancer_id").(string))
	request.ListenerPort = helper.IntUint64(d.Get("listener_port").(int))
	request.ListenerProtocol = helper.String(d.Get("listener_protocol").(string))
	request.DefaultActions = buildAlbDefaultActions(d)
	if len(request.DefaultActions) == 0 {
		return fmt.Errorf("[CHECK][ALB listener][Create] check: default_actions is required")
	}

	if v, ok := d.GetOk("listener_name"); ok {
		request.ListenerName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("certificate_ids"); ok {
		certificateIds := make([]*string, 0, len(v.([]interface{})))
		for _, item := range v.([]interface{}) {
			certificateIds = append(certificateIds, helper.String(item.(string)))
		}
		request.CertificateIds = certificateIds
	}
	if v, ok := d.GetOkExists("ca_enabled"); ok {
		request.CaEnabled = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOk("ca_certificate_ids"); ok {
		caCertificateIds := make([]*string, 0, len(v.([]interface{})))
		for _, item := range v.([]interface{}) {
			caCertificateIds = append(caCertificateIds, helper.String(item.(string)))
		}
		request.CaCertificateIds = caCertificateIds
	}
	if v, ok := d.GetOk("security_policy_id"); ok {
		request.SecurityPolicyId = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("idle_timeout"); ok {
		request.IdleTimeout = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOkExists("request_timeout"); ok {
		request.RequestTimeout = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOkExists("gzip_enabled"); ok {
		request.GzipEnabled = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("http2enabled"); ok {
		request.Http2Enabled = helper.Bool(v.(bool))
	}
	if config := buildAlbXForwardedForConfig(d); config != nil {
		request.XForwardedForConfig = config
	}
	if tags := buildAlbTags(d); tags != nil {
		request.Tags = tags
	}
	if v, ok := d.GetOk("client_token"); ok {
		request.ClientToken = helper.String(v.(string))
	}

	var response *alb.CreateListenerResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAlbClient().CreateListenerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ALB listener failed, reason:%+v\n", logId, err)
		return err
	}
	if response == nil || response.Response == nil || response.Response.ListenerId == nil {
		return fmt.Errorf("[CHECK][ALB listener][Create] check: response error, listener id is null")
	}
	listenerId := *response.Response.ListenerId
	d.SetId(*request.LoadBalancerId + tccommon.FILED_SP + listenerId)

	return resourceTencentCloudAlbListenerRead(d, meta)
}

func resourceTencentCloudAlbListenerRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_listener.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) != 2 {
		return fmt.Errorf("broken ID %s, expected format load_balancer_id#listener_id", d.Id())
	}
	loadBalancerId := items[0]
	listenerId := items[1]

	request := alb.NewDescribeListenersRequest()
	request.LoadBalancerId = helper.String(loadBalancerId)
	request.ListenerIds = []*string{helper.String(listenerId)}

	var listener *alb.ListenerOutput
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAlbClient().DescribeListenersWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe ALB listeners failed, Response is nil."))
		}
		if len(result.Response.Listeners) > 0 {
			listener = result.Response.Listeners[0]
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read ALB listener failed, reason:%+v\n", logId, err)
		return err
	}

	if listener == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("load_balancer_id", loadBalancerId)
	if listener.ListenerId != nil {
		_ = d.Set("listener_id", listener.ListenerId)
	}
	if listener.ListenerPort != nil {
		_ = d.Set("listener_port", int(*listener.ListenerPort))
	}
	if listener.ListenerProtocol != nil {
		_ = d.Set("listener_protocol", listener.ListenerProtocol)
	}
	if listener.ListenerName != nil {
		_ = d.Set("listener_name", listener.ListenerName)
	}
	if listener.GzipEnabled != nil {
		_ = d.Set("gzip_enabled", listener.GzipEnabled)
	}
	if listener.Http2Enable != nil {
		_ = d.Set("http2enabled", listener.Http2Enable)
	}
	if listener.IdleTimeout != nil {
		_ = d.Set("idle_timeout", int(*listener.IdleTimeout))
	}
	if listener.RequestTimeout != nil {
		_ = d.Set("request_timeout", int(*listener.RequestTimeout))
	}
	if listener.CaEnable != nil {
		_ = d.Set("ca_enabled", listener.CaEnable)
	}
	if listener.TlsSecurityPolicyId != nil {
		_ = d.Set("security_policy_id", listener.TlsSecurityPolicyId)
	}
	if listener.XForwardedForConfig != nil {
		config := make(map[string]interface{}, 0)
		if listener.XForwardedForConfig.XForwardedForMode != nil {
			config["x_forwarded_for_mode"] = *listener.XForwardedForConfig.XForwardedForMode
		}
		if listener.XForwardedForConfig.XForwardedForAlbIdEnabled != nil {
			config["x_forwarded_for_alb_id_enabled"] = *listener.XForwardedForConfig.XForwardedForAlbIdEnabled
		}
		if listener.XForwardedForConfig.XForwardedForClientSrcPortEnabled != nil {
			config["x_forwarded_for_client_src_port_enabled"] = *listener.XForwardedForConfig.XForwardedForClientSrcPortEnabled
		}
		if listener.XForwardedForConfig.XForwardedForHostEnabled != nil {
			config["x_forwarded_for_host_enabled"] = *listener.XForwardedForConfig.XForwardedForHostEnabled
		}
		if listener.XForwardedForConfig.XForwardedForPortEnabled != nil {
			config["x_forwarded_for_port_enabled"] = *listener.XForwardedForConfig.XForwardedForPortEnabled
		}
		if listener.XForwardedForConfig.XForwardedForProtoEnabled != nil {
			config["x_forwarded_for_proto_enabled"] = *listener.XForwardedForConfig.XForwardedForProtoEnabled
		}
		if listener.XForwardedForConfig.XTencentClientIDNEnabled != nil {
			config["x_tencent_client_idn_enabled"] = *listener.XForwardedForConfig.XTencentClientIDNEnabled
		}
		if listener.XForwardedForConfig.XTencentClientSDNEnabled != nil {
			config["x_tencent_client_sdn_enabled"] = *listener.XForwardedForConfig.XTencentClientSDNEnabled
		}
		if listener.XForwardedForConfig.XTencentClientSerialEnabled != nil {
			config["x_tencent_client_serial_enabled"] = *listener.XForwardedForConfig.XTencentClientSerialEnabled
		}
		if listener.XForwardedForConfig.XTencentClientVerifyEnabled != nil {
			config["x_tencent_client_verify_enabled"] = *listener.XForwardedForConfig.XTencentClientVerifyEnabled
		}
		configs := make([]map[string]interface{}, 0, 1)
		configs = append(configs, config)
		_ = d.Set("x_forwarded_for_config", configs)
	}
	if listener.Tags != nil {
		tagList := make([]map[string]interface{}, 0, len(listener.Tags))
		for _, tag := range listener.Tags {
			tagMap := make(map[string]interface{}, 0)
			if tag.TagKey != nil {
				tagMap["tag_key"] = *tag.TagKey
			}
			if tag.TagValue != nil {
				tagMap["tag_value"] = *tag.TagValue
			}
			tagList = append(tagList, tagMap)
		}
		_ = d.Set("tags", tagList)
	}

	return nil
}

func resourceTencentCloudAlbListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_listener.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) != 2 {
		return fmt.Errorf("broken ID %s, expected format load_balancer_id#listener_id", d.Id())
	}
	loadBalancerId := items[0]
	listenerId := items[1]

	request := alb.NewModifyListenerAttributesRequest()
	request.ListenerId = helper.String(listenerId)
	request.LoadBalancerId = helper.String(loadBalancerId)
	changed := false

	if d.HasChange("listener_name") {
		changed = true
		request.ListenerName = helper.String(d.Get("listener_name").(string))
	}
	if d.HasChange("default_actions") {
		changed = true
		request.DefaultActions = buildAlbDefaultActions(d)
	}
	if d.HasChange("certificate_ids") {
		changed = true
		certificateIds := make([]*string, 0, 1)
		if v, ok := d.GetOk("certificate_ids"); ok {
			for _, item := range v.([]interface{}) {
				certificateIds = append(certificateIds, helper.String(item.(string)))
			}
		}
		request.CertificateIds = certificateIds
	}
	if d.HasChange("ca_enabled") {
		changed = true
		request.CaEnabled = helper.Bool(d.Get("ca_enabled").(bool))
	}
	if d.HasChange("ca_certificate_ids") {
		changed = true
		caCertificateIds := make([]*string, 0, 1)
		if v, ok := d.GetOk("ca_certificate_ids"); ok {
			for _, item := range v.([]interface{}) {
				caCertificateIds = append(caCertificateIds, helper.String(item.(string)))
			}
		}
		request.CaCertificateIds = caCertificateIds
	}
	if d.HasChange("security_policy_id") {
		changed = true
		request.SecurityPolicyId = helper.String(d.Get("security_policy_id").(string))
	}
	if d.HasChange("idle_timeout") {
		changed = true
		request.IdleTimeout = helper.IntUint64(d.Get("idle_timeout").(int))
	}
	if d.HasChange("request_timeout") {
		changed = true
		request.RequestTimeout = helper.IntUint64(d.Get("request_timeout").(int))
	}
	if d.HasChange("gzip_enabled") {
		changed = true
		request.GzipEnabled = helper.Bool(d.Get("gzip_enabled").(bool))
	}
	if d.HasChange("http2enabled") {
		changed = true
		request.Http2Enabled = helper.Bool(d.Get("http2enabled").(bool))
	}
	if d.HasChange("x_forwarded_for_config") {
		changed = true
		request.XForwardedForConfig = buildAlbXForwardedForConfig(d)
	}

	if !changed {
		return nil
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAlbClient().ModifyListenerAttributesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update ALB listener failed, reason:%+v\n", logId, err)
		return err
	}

	return resourceTencentCloudAlbListenerRead(d, meta)
}

func resourceTencentCloudAlbListenerDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_listener.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) != 2 {
		return fmt.Errorf("broken ID %s, expected format load_balancer_id#listener_id", d.Id())
	}
	loadBalancerId := items[0]
	listenerId := items[1]

	request := alb.NewDeleteListenerRequest()
	request.LoadBalancerId = helper.String(loadBalancerId)
	request.ListenerIds = []*string{helper.String(listenerId)}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAlbClient().DeleteListenerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete ALB listener failed, reason:%+v\n", logId, err)
		return err
	}

	describeRequest := alb.NewDescribeListenersRequest()
	describeRequest.LoadBalancerId = helper.String(loadBalancerId)
	describeRequest.ListenerIds = []*string{helper.String(listenerId)}
	return resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAlbClient().DescribeListenersWithContext(ctx, describeRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result != nil && result.Response != nil && len(result.Response.Listeners) > 0 {
			return resource.RetryableError(fmt.Errorf("ALB listener %s is still deleting", listenerId))
		}
		return nil
	})
}
