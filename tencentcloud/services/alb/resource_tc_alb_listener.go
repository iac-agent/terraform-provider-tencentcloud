package alb

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	alb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAlbListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAlbListenerCreate,
		Read:   resourceTencentCloudAlbListenerRead,
		Update: resourceTencentCloudAlbListenerUpdate,
		Delete: resourceTencentCloudAlbListenerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceTencentCloudAlbListenerImport,
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "ID of the load balancer instance, format is alb- followed by 8 alphanumeric characters.",
			},
			"listener_port": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
				Description: "Port used by the frontend of the load balancer instance. Valid value range: 1~65535.",
			},
			"listener_protocol": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Listener protocol. Valid values: HTTP, HTTPS or QUIC.",
			},
			"default_actions": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "List of default forwarding rule actions. Currently a listener supports only 1 default forwarding rule action.",
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
							Description: "Forwarding target group configuration. Only a single target group is supported.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_groups": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "List of target groups.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"target_group_id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Target group ID, format is lbtg- followed by 8 alphanumeric characters.",
												},
												"weight": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Weight, valid value range: [0, 100], default is 10.",
												},
											},
										},
									},
									"target_group_sticky_session": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Sticky session between target groups.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:        schema.TypeBool,
													Required:    true,
													Description: "Whether to enable sticky session, default is disabled.",
												},
												"timeout": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Timeout in seconds, valid value range: 1~86400, default is 1000.",
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
			"ca_certificate_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of CA certificate IDs configured for the listener. Currently only 1 CA certificate is supported. Required when ca_enabled is true.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ca_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable mutual authentication. true: enable. false (default): disable.",
			},
			"certificate_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of server certificate IDs.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"client_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Client token used to ensure request idempotency. Only ASCII characters are supported.",
			},
			"gzip_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable Gzip compression. true (default): yes. false: no.",
			},
			"http2enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable HTTP/2. Default false for HTTP protocol, default true for HTTPS protocol. Only HTTPS protocol supports this parameter.",
			},
			"idle_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Connection idle timeout in seconds. Valid value range: 1~600. Default is 15.",
			},
			"listener_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom listener name. Length is 1~255 characters, can contain Chinese characters, letters, digits, hyphens (-), slashes (/), periods (.) and underscores (_).",
			},
			"request_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Connection request timeout in seconds. Valid value range: 1~600. Default is 60.",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Security policy ID, format is tls- followed by 8 alphanumeric characters.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    20,
				Description: "List of tags. Up to 20 supported.",
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
			"x_forwarded_for_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "X-Forwarded-For configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"x_forwarded_for_alb_id_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to obtain the load balancer instance ID via the ALB-ID header field.",
						},
						"x_forwarded_for_client_src_port_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to obtain the client port via the X-Forwarded-Client-srcport header field.",
						},
						"x_forwarded_for_host_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to obtain the client domain via the X-Forwarded-Host header field.",
						},
						"x_forwarded_for_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies how to handle the X-Forwarded-For (XFF) HTTP header field. Valid values: append, remove, passthrough.",
						},
						"x_forwarded_for_port_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to obtain the listener port via the X-Forwarded-Port header field.",
						},
						"x_forwarded_for_proto_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to obtain the listener protocol via the X-Forwarded-Proto header field.",
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
			"listener_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Listener ID, format is lst- followed by 8 alphanumeric characters.",
			},
			"listener_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Listener status. Valid values: Active, Provisioning, Configuring, ProvisionFailed.",
			},
		},
	}
}

func buildAlbDefaultActions(d *schema.ResourceData) []*alb.DefaultAction {
	var actions []*alb.DefaultAction
	actionsRaw, ok := d.GetOk("default_actions")
	if !ok {
		return actions
	}
	actionsList := actionsRaw.([]interface{})
	for _, actionItem := range actionsList {
		actionMap := actionItem.(map[string]interface{})
		action := alb.DefaultAction{}
		if v, ok := actionMap["type"]; ok {
			action.Type = helper.String(v.(string))
		}
		if v, ok := actionMap["target_group_config"]; ok {
			configList := v.([]interface{})
			if len(configList) > 0 {
				configMap := configList[0].(map[string]interface{})
				config := alb.TargetGroupConfig{}
				if v, ok := configMap["target_groups"]; ok {
					for _, tgItem := range v.([]interface{}) {
						tgMap := tgItem.(map[string]interface{})
						tg := alb.TargetGroupTuple{}
						if v, ok := tgMap["target_group_id"]; ok {
							tg.TargetGroupId = helper.String(v.(string))
						}
						if v, ok := tgMap["weight"]; ok {
							if v.(int) > 0 {
								tg.Weight = helper.Int64(int64(v.(int)))
							}
						}
						config.TargetGroups = append(config.TargetGroups, &tg)
					}
				}
				if v, ok := configMap["target_group_sticky_session"]; ok {
					sessionList := v.([]interface{})
					if len(sessionList) > 0 {
						sessionMap := sessionList[0].(map[string]interface{})
						session := alb.TargetGroupStickySession{}
						if v, ok := sessionMap["enabled"]; ok {
							session.Enabled = helper.Bool(v.(bool))
						}
						if v, ok := sessionMap["timeout"]; ok {
							if v.(int) > 0 {
								session.Timeout = helper.Int64(int64(v.(int)))
							}
						}
						config.TargetGroupStickySession = &session
					}
				}
				action.TargetGroupConfig = &config
			}
		}
		actions = append(actions, &action)
	}
	return actions
}

func buildAlbXForwardedForConfig(d *schema.ResourceData) *alb.XForwardedForConfig {
	xffRaw, ok := d.GetOk("x_forwarded_for_config")
	if !ok {
		return nil
	}
	xffList := xffRaw.([]interface{})
	if len(xffList) == 0 {
		return nil
	}
	xffMap := xffList[0].(map[string]interface{})
	xff := alb.XForwardedForConfig{}
	if v, ok := xffMap["x_forwarded_for_alb_id_enabled"]; ok {
		xff.XForwardedForAlbIdEnabled = helper.Bool(v.(bool))
	}
	if v, ok := xffMap["x_forwarded_for_client_src_port_enabled"]; ok {
		xff.XForwardedForClientSrcPortEnabled = helper.Bool(v.(bool))
	}
	if v, ok := xffMap["x_forwarded_for_host_enabled"]; ok {
		xff.XForwardedForHostEnabled = helper.Bool(v.(bool))
	}
	if v, ok := xffMap["x_forwarded_for_mode"]; ok {
		if s := v.(string); s != "" {
			xff.XForwardedForMode = helper.String(s)
		}
	}
	if v, ok := xffMap["x_forwarded_for_port_enabled"]; ok {
		xff.XForwardedForPortEnabled = helper.Bool(v.(bool))
	}
	if v, ok := xffMap["x_forwarded_for_proto_enabled"]; ok {
		xff.XForwardedForProtoEnabled = helper.Bool(v.(bool))
	}
	if v, ok := xffMap["x_tencent_client_idn_enabled"]; ok {
		xff.XTencentClientIDNEnabled = helper.Bool(v.(bool))
	}
	if v, ok := xffMap["x_tencent_client_sdn_enabled"]; ok {
		xff.XTencentClientSDNEnabled = helper.Bool(v.(bool))
	}
	if v, ok := xffMap["x_tencent_client_serial_enabled"]; ok {
		xff.XTencentClientSerialEnabled = helper.Bool(v.(bool))
	}
	if v, ok := xffMap["x_tencent_client_verify_enabled"]; ok {
		xff.XTencentClientVerifyEnabled = helper.Bool(v.(bool))
	}
	return &xff
}

func setAlbDefaultActions(d *schema.ResourceData, actions []*alb.DefaultAction) error {
	var result []interface{}
	for _, action := range actions {
		actionMap := map[string]interface{}{
			"type": "",
		}
		if action.Type != nil {
			actionMap["type"] = *action.Type
		}
		if action.TargetGroupConfig != nil {
			tgConfigMap := map[string]interface{}{
				"target_groups": []interface{}{},
			}
			var tgList []interface{}
			for _, tg := range action.TargetGroupConfig.TargetGroups {
				tgMap := map[string]interface{}{}
				if tg.TargetGroupId != nil {
					tgMap["target_group_id"] = *tg.TargetGroupId
				}
				if tg.Weight != nil {
					tgMap["weight"] = *tg.Weight
				}
				tgList = append(tgList, tgMap)
			}
			tgConfigMap["target_groups"] = tgList
			if action.TargetGroupConfig.TargetGroupStickySession != nil {
				sessionMap := map[string]interface{}{}
				if action.TargetGroupConfig.TargetGroupStickySession.Enabled != nil {
					sessionMap["enabled"] = *action.TargetGroupConfig.TargetGroupStickySession.Enabled
				}
				if action.TargetGroupConfig.TargetGroupStickySession.Timeout != nil {
					sessionMap["timeout"] = *action.TargetGroupConfig.TargetGroupStickySession.Timeout
				}
				tgConfigMap["target_group_sticky_session"] = []interface{}{sessionMap}
			}
			actionMap["target_group_config"] = []interface{}{tgConfigMap}
		}
		result = append(result, actionMap)
	}
	if len(result) > 0 {
		return d.Set("default_actions", result)
	}
	return nil
}

func setAlbXForwardedForConfig(d *schema.ResourceData, xff *alb.XForwardedForConfig) error {
	if xff == nil {
		return nil
	}
	xffMap := map[string]interface{}{}
	if xff.XForwardedForAlbIdEnabled != nil {
		xffMap["x_forwarded_for_alb_id_enabled"] = *xff.XForwardedForAlbIdEnabled
	}
	if xff.XForwardedForClientSrcPortEnabled != nil {
		xffMap["x_forwarded_for_client_src_port_enabled"] = *xff.XForwardedForClientSrcPortEnabled
	}
	if xff.XForwardedForHostEnabled != nil {
		xffMap["x_forwarded_for_host_enabled"] = *xff.XForwardedForHostEnabled
	}
	if xff.XForwardedForMode != nil {
		xffMap["x_forwarded_for_mode"] = *xff.XForwardedForMode
	}
	if xff.XForwardedForPortEnabled != nil {
		xffMap["x_forwarded_for_port_enabled"] = *xff.XForwardedForPortEnabled
	}
	if xff.XForwardedForProtoEnabled != nil {
		xffMap["x_forwarded_for_proto_enabled"] = *xff.XForwardedForProtoEnabled
	}
	if xff.XTencentClientIDNEnabled != nil {
		xffMap["x_tencent_client_idn_enabled"] = *xff.XTencentClientIDNEnabled
	}
	if xff.XTencentClientSDNEnabled != nil {
		xffMap["x_tencent_client_sdn_enabled"] = *xff.XTencentClientSDNEnabled
	}
	if xff.XTencentClientSerialEnabled != nil {
		xffMap["x_tencent_client_serial_enabled"] = *xff.XTencentClientSerialEnabled
	}
	if xff.XTencentClientVerifyEnabled != nil {
		xffMap["x_tencent_client_verify_enabled"] = *xff.XTencentClientVerifyEnabled
	}
	return d.Set("x_forwarded_for_config", []interface{}{xffMap})
}

func resourceTencentCloudAlbListenerCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_listener.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		request  = alb.NewCreateListenerRequest()
		response = alb.NewCreateListenerResponse()
	)

	request.LoadBalancerId = helper.String(d.Get("load_balancer_id").(string))
	request.ListenerPort = helper.Uint64(uint64(d.Get("listener_port").(int)))
	request.ListenerProtocol = helper.String(d.Get("listener_protocol").(string))

	if actions := buildAlbDefaultActions(d); len(actions) > 0 {
		request.DefaultActions = actions
	}

	if v, ok := d.GetOk("ca_certificate_ids"); ok {
		request.CaCertificateIds = helper.Strings(v.([]string))
	}
	if v, ok := d.GetOkExists("ca_enabled"); ok {
		request.CaEnabled = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOk("certificate_ids"); ok {
		request.CertificateIds = helper.Strings(v.([]string))
	}
	if v, ok := d.GetOk("client_token"); ok {
		request.ClientToken = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("gzip_enabled"); ok {
		request.GzipEnabled = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("http2enabled"); ok {
		request.Http2Enabled = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOk("idle_timeout"); ok {
		request.IdleTimeout = helper.Uint64(uint64(v.(int)))
	}
	if v, ok := d.GetOk("listener_name"); ok {
		request.ListenerName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("request_timeout"); ok {
		request.RequestTimeout = helper.Uint64(uint64(v.(int)))
	}
	if v, ok := d.GetOk("security_policy_id"); ok {
		request.SecurityPolicyId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			tagsMap := item.(map[string]interface{})
			tagInfo := alb.TagInfo{}
			if v, ok := tagsMap["tag_key"]; ok {
				tagInfo.TagKey = helper.String(v.(string))
			}
			if v, ok := tagsMap["tag_value"]; ok {
				tagInfo.TagValue = helper.String(v.(string))
			}
			request.Tags = append(request.Tags, &tagInfo)
		}
	}
	if xff := buildAlbXForwardedForConfig(d); xff != nil {
		request.XForwardedForConfig = xff
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAlbV20251030Client().CreateListenerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create alb listener failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.ListenerId == nil {
		return fmt.Errorf("listener id not returned by CreateListener")
	}
	d.SetId(*response.Response.ListenerId)

	return resourceTencentCloudAlbListenerRead(d, meta)
}

func resourceTencentCloudAlbListenerRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_listener.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	loadBalancerId := d.Get("load_balancer_id").(string)

	svc := NewAlbService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	listener, err := svc.DescribeListenerById(ctx, loadBalancerId, d.Id())
	if err != nil {
		return err
	}

	if listener == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `alb_listener` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if listener.ListenerPort != nil {
		_ = d.Set("listener_port", *listener.ListenerPort)
	}
	if listener.ListenerProtocol != nil {
		_ = d.Set("listener_protocol", *listener.ListenerProtocol)
	}
	if listener.ListenerId != nil {
		_ = d.Set("listener_id", *listener.ListenerId)
	}
	if listener.ListenerStatus != nil {
		_ = d.Set("listener_status", *listener.ListenerStatus)
	}
	if listener.ListenerName != nil {
		_ = d.Set("listener_name", *listener.ListenerName)
	}
	if listener.GzipEnabled != nil {
		_ = d.Set("gzip_enabled", *listener.GzipEnabled)
	}
	if listener.Http2Enable != nil {
		_ = d.Set("http2enabled", *listener.Http2Enable)
	}
	if listener.IdleTimeout != nil {
		_ = d.Set("idle_timeout", *listener.IdleTimeout)
	}
	if listener.RequestTimeout != nil {
		_ = d.Set("request_timeout", *listener.RequestTimeout)
	}
	if listener.TlsSecurityPolicyId != nil {
		_ = d.Set("security_policy_id", *listener.TlsSecurityPolicyId)
	}
	if listener.CaEnable != nil {
		_ = d.Set("ca_enabled", *listener.CaEnable)
	}
	if err = setAlbXForwardedForConfig(d, listener.XForwardedForConfig); err != nil {
		return fmt.Errorf("set x_forwarded_for_config failed, reason:%+v", err)
	}

	var tagList []interface{}
	for _, tag := range listener.Tags {
		tagMap := map[string]interface{}{}
		if tag.TagKey != nil {
			tagMap["tag_key"] = *tag.TagKey
		}
		if tag.TagValue != nil {
			tagMap["tag_value"] = *tag.TagValue
		}
		tagList = append(tagList, tagMap)
	}
	if len(tagList) > 0 {
		if err = d.Set("tags", tagList); err != nil {
			return fmt.Errorf("set tags failed, reason:%+v", err)
		}
	}

	return nil
}

func resourceTencentCloudAlbListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_listener.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	request := alb.NewModifyListenerAttributesRequest()

	request.LoadBalancerId = helper.String(d.Get("load_balancer_id").(string))
	request.ListenerId = helper.String(d.Id())

	if d.HasChange("listener_name") {
		request.ListenerName = helper.String(d.Get("listener_name").(string))
	}
	if d.HasChange("default_actions") {
		if actions := buildAlbDefaultActions(d); len(actions) > 0 {
			request.DefaultActions = actions
		}
	}
	if d.HasChange("ca_certificate_ids") {
		v, ok := d.GetOk("ca_certificate_ids")
		if ok {
			request.CaCertificateIds = helper.Strings(v.([]string))
		}
	}
	if d.HasChange("ca_enabled") {
		request.CaEnabled = helper.Bool(d.Get("ca_enabled").(bool))
	}
	if d.HasChange("certificate_ids") {
		v, ok := d.GetOk("certificate_ids")
		if ok {
			request.CertificateIds = helper.Strings(v.([]string))
		}
	}
	if d.HasChange("client_token") {
		if v, ok := d.GetOk("client_token"); ok {
			request.ClientToken = helper.String(v.(string))
		}
	}
	if d.HasChange("gzip_enabled") {
		request.GzipEnabled = helper.Bool(d.Get("gzip_enabled").(bool))
	}
	if d.HasChange("http2enabled") {
		request.Http2Enabled = helper.Bool(d.Get("http2enabled").(bool))
	}
	if d.HasChange("idle_timeout") {
		request.IdleTimeout = helper.Uint64(uint64(d.Get("idle_timeout").(int)))
	}
	if d.HasChange("request_timeout") {
		request.RequestTimeout = helper.Uint64(uint64(d.Get("request_timeout").(int)))
	}
	if d.HasChange("security_policy_id") {
		if v, ok := d.GetOk("security_policy_id"); ok {
			request.SecurityPolicyId = helper.String(v.(string))
		}
	}
	if d.HasChange("x_forwarded_for_config") {
		if xff := buildAlbXForwardedForConfig(d); xff != nil {
			request.XForwardedForConfig = xff
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAlbV20251030Client().ModifyListenerAttributesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update alb listener failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudAlbListenerRead(d, meta)
}

func resourceTencentCloudAlbListenerDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_listener.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	request := alb.NewDeleteListenerRequest()

	request.LoadBalancerId = helper.String(d.Get("load_balancer_id").(string))
	request.ListenerIds = helper.Strings([]string{d.Id()})

	if v, ok := d.GetOk("client_token"); ok {
		request.ClientToken = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAlbV20251030Client().DeleteListenerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete alb listener failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

func resourceTencentCloudAlbListenerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import id, expected format: <load_balancer_id>:<listener_id>")
	}
	_ = d.Set("load_balancer_id", parts[0])
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, nil
}
