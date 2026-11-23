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
				Description: "ID of the ALB instance, in the format of `alb-xxxxxxxx`.",
			},
			"listener_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 65535),
				Description:  "Front-end port of the ALB instance. Value range: 1-65535.",
			},
			"listener_protocol": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Listener protocol. Valid values: `HTTP`, `HTTPS` and `QUIC`.",
			},
			"default_actions": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "List of default forwarding rule actions. Currently, only 1 default forwarding rule action is supported per listener.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type of the forwarding action. When creating a listener, only forwarding to target group is supported.",
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
													Description: "Target group ID, in the format of `lbtg-xxxxxxxx`.",
												},
												"weight": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Weight. Value range: [0, 100]. Default is 10.",
												},
											},
										},
									},
									"target_group_sticky_session": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Sticky session among target groups.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:        schema.TypeBool,
													Required:    true,
													Description: "Whether to enable sticky session. Default is disabled.",
												},
												"timeout": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: tccommon.ValidateIntegerInRange(1, 86400),
													Description:  "Timeout in seconds. Value range: 1-86400. Default is 1000.",
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
				MaxItems:    1,
				Description: "List of CA certificate IDs of the listener. Currently, only 1 CA certificate is supported. Required when `ca_enabled` is `true`.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ca_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable mutual authentication. Default is `false`.",
			},
			"certificate_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of server certificate IDs.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"gzip_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable Gzip compression. Default is `true`.",
			},
			"http2enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable HTTP/2. HTTP protocol defaults to `false`, HTTPS protocol defaults to `true`. Only HTTPS protocol supports this parameter.",
			},
			"idle_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 600),
				Description:  "Idle connection timeout in seconds. Value range: 1-600. Default is 15.",
			},
			"listener_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom listener name. Length is 1-255 characters, which can contain Chinese characters, letters, digits, hyphens (-), slashes (/), periods (.) and underscores (_).",
			},
			"request_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 600),
				Description:  "Request timeout in seconds. Value range: 1-600. Default is 60. If the backend server does not respond within the timeout, the ALB aborts waiting and returns HTTP 504 to the client.",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Security policy ID, in the format of `tls-xxxxxxxx`.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    20,
				Description: "List of tags. Up to 20 tags are supported.",
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
							Description: "Whether to get the ALB instance ID via the ALB-ID header.",
						},
						"x_forwarded_for_client_src_port_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to get the client port via the X-Forwarded-Client-srcport header.",
						},
						"x_forwarded_for_host_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to get the client domain via the X-Forwarded-Host header.",
						},
						"x_forwarded_for_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "How to handle the X-Forwarded-For (XFF) HTTP header. Valid values: `append` (default, append the real client IP to the X-Forwarded-For header), `remove` (remove the X-Forwarded-For header), `passthrough` (pass the X-Forwarded-For header through unchanged).",
						},
						"x_forwarded_for_port_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to get the listener port via the X-Forwarded-Port header.",
						},
						"x_forwarded_for_proto_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to get the listener protocol via the X-Forwarded-Proto header.",
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
		},
	}
}

func buildAlbDefaultActions(d *schema.ResourceData) []*alb.DefaultAction {
	var actions []*alb.DefaultAction
	if v, ok := d.GetOk("default_actions"); ok {
		for _, item := range v.([]interface{}) {
			actionMap := item.(map[string]interface{})
			action := alb.DefaultAction{}
			if v, ok := actionMap["type"].(string); ok && v != "" {
				action.Type = helper.String(v)
			}
			if items, ok := actionMap["target_group_config"].([]interface{}); ok {
				for _, item := range items {
					configMap := item.(map[string]interface{})
					config := alb.TargetGroupConfig{}
					if tgItems, ok := configMap["target_groups"].([]interface{}); ok {
						for _, item := range tgItems {
							tgMap := item.(map[string]interface{})
							tg := alb.TargetGroupTuple{}
							if v, ok := tgMap["target_group_id"].(string); ok && v != "" {
								tg.TargetGroupId = helper.String(v)
							}
							if v, ok := tgMap["weight"].(int); ok {
								tg.Weight = helper.IntInt64(v)
							}
							config.TargetGroups = append(config.TargetGroups, &tg)
						}
					}
					if ssItems, ok := configMap["target_group_sticky_session"].([]interface{}); ok {
						for _, item := range ssItems {
							ssMap := item.(map[string]interface{})
							ss := alb.TargetGroupStickySession{}
							if v, ok := ssMap["enabled"].(bool); ok {
								ss.Enabled = helper.Bool(v)
							}
							if v, ok := ssMap["timeout"].(int); ok {
								ss.Timeout = helper.IntInt64(v)
							}
							config.TargetGroupStickySession = &ss
						}
					}
					action.TargetGroupConfig = &config
				}
			}
			actions = append(actions, &action)
		}
	}
	return actions
}

func buildAlbXForwardedForConfig(d *schema.ResourceData) *alb.XForwardedForConfig {
	if v, ok := d.GetOk("x_forwarded_for_config"); ok {
		for _, item := range v.([]interface{}) {
			configMap := item.(map[string]interface{})
			config := alb.XForwardedForConfig{}
			if v, ok := configMap["x_forwarded_for_alb_id_enabled"].(bool); ok {
				config.XForwardedForAlbIdEnabled = helper.Bool(v)
			}
			if v, ok := configMap["x_forwarded_for_client_src_port_enabled"].(bool); ok {
				config.XForwardedForClientSrcPortEnabled = helper.Bool(v)
			}
			if v, ok := configMap["x_forwarded_for_host_enabled"].(bool); ok {
				config.XForwardedForHostEnabled = helper.Bool(v)
			}
			if v, ok := configMap["x_forwarded_for_mode"].(string); ok && v != "" {
				config.XForwardedForMode = helper.String(v)
			}
			if v, ok := configMap["x_forwarded_for_port_enabled"].(bool); ok {
				config.XForwardedForPortEnabled = helper.Bool(v)
			}
			if v, ok := configMap["x_forwarded_for_proto_enabled"].(bool); ok {
				config.XForwardedForProtoEnabled = helper.Bool(v)
			}
			if v, ok := configMap["x_tencent_client_idn_enabled"].(bool); ok {
				config.XTencentClientIDNEnabled = helper.Bool(v)
			}
			if v, ok := configMap["x_tencent_client_sdn_enabled"].(bool); ok {
				config.XTencentClientSDNEnabled = helper.Bool(v)
			}
			if v, ok := configMap["x_tencent_client_serial_enabled"].(bool); ok {
				config.XTencentClientSerialEnabled = helper.Bool(v)
			}
			if v, ok := configMap["x_tencent_client_verify_enabled"].(bool); ok {
				config.XTencentClientVerifyEnabled = helper.Bool(v)
			}
			return &config
		}
	}
	return nil
}

func buildAlbCertificateIds(d *schema.ResourceData) []*string {
	var ids []*string
	if v, ok := d.GetOk("certificate_ids"); ok {
		for _, item := range v.([]interface{}) {
			if s, ok := item.(string); ok && s != "" {
				ids = append(ids, helper.String(s))
			}
		}
	}
	return ids
}

func buildAlbCaCertificateIds(d *schema.ResourceData) []*string {
	var ids []*string
	if v, ok := d.GetOk("ca_certificate_ids"); ok {
		for _, item := range v.([]interface{}) {
			if s, ok := item.(string); ok && s != "" {
				ids = append(ids, helper.String(s))
			}
		}
	}
	return ids
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
		return fmt.Errorf("default_actions must be set")
	}

	if v, ok := d.GetOkExists("ca_enabled"); ok {
		request.CaEnabled = helper.Bool(v.(bool))
	}
	request.CaCertificateIds = buildAlbCaCertificateIds(d)
	request.CertificateIds = buildAlbCertificateIds(d)
	if v, ok := d.GetOkExists("gzip_enabled"); ok {
		request.GzipEnabled = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("http2enabled"); ok {
		request.Http2Enabled = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("idle_timeout"); ok {
		request.IdleTimeout = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("listener_name"); ok {
		request.ListenerName = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("request_timeout"); ok {
		request.RequestTimeout = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("security_policy_id"); ok {
		request.SecurityPolicyId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			tagMap := item.(map[string]interface{})
			tag := alb.TagInfo{}
			if key, ok := tagMap["tag_key"].(string); ok && key != "" {
				tag.TagKey = helper.String(key)
			}
			if value, ok := tagMap["tag_value"].(string); ok && value != "" {
				tag.TagValue = helper.String(value)
			}
			request.Tags = append(request.Tags, &tag)
		}
	}
	if config := buildAlbXForwardedForConfig(d); config != nil {
		request.XForwardedForConfig = config
	}

	var response *alb.CreateListenerResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAlbClient().CreateListenerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if result == nil || result.Response == nil || result.Response.ListenerId == nil {
			return resource.NonRetryableError(fmt.Errorf("create ALB listener failed, Response is nil"))
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ALB listener failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(*response.Response.ListenerId)

	return resourceTencentCloudAlbListenerRead(d, meta)
}

func resourceTencentCloudAlbListenerRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_listener.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	listenerId := items[len(items)-1]
	loadBalancerId := d.Get("load_balancer_id").(string)
	if len(items) == 2 {
		listenerId = items[1]
		loadBalancerId = items[0]
	}
	if loadBalancerId == "" {
		return fmt.Errorf("load_balancer_id is required for ALB listener import, please use loadBalancerId#listenerId style")
	}

	request := alb.NewDescribeListenersRequest()
	request.LoadBalancerId = helper.String(loadBalancerId)
	request.ListenerIds = []*string{helper.String(listenerId)}

	var instance *alb.ListenerOutput
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAlbClient().DescribeListenersWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("describe ALB listeners failed, Response is nil"))
		}
		if len(result.Response.Listeners) > 0 {
			instance = result.Response.Listeners[0]
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read ALB listener failed, reason:%+v", logId, err)
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	d.SetId(loadBalancerId + tccommon.FILED_SP + listenerId)
	_ = d.Set("load_balancer_id", loadBalancerId)
	if instance.ListenerId != nil {
		_ = d.Set("listener_id", instance.ListenerId)
	}
	if instance.ListenerPort != nil {
		_ = d.Set("listener_port", instance.ListenerPort)
	}
	if instance.ListenerProtocol != nil {
		_ = d.Set("listener_protocol", instance.ListenerProtocol)
	}
	if instance.ListenerName != nil {
		_ = d.Set("listener_name", instance.ListenerName)
	}
	if instance.CaEnable != nil {
		_ = d.Set("ca_enabled", instance.CaEnable)
	}
	if instance.GzipEnabled != nil {
		_ = d.Set("gzip_enabled", instance.GzipEnabled)
	}
	if instance.Http2Enable != nil {
		_ = d.Set("http2enabled", instance.Http2Enable)
	}
	if instance.IdleTimeout != nil {
		_ = d.Set("idle_timeout", instance.IdleTimeout)
	}
	if instance.RequestTimeout != nil {
		_ = d.Set("request_timeout", instance.RequestTimeout)
	}
	if instance.TlsSecurityPolicyId != nil {
		_ = d.Set("security_policy_id", instance.TlsSecurityPolicyId)
	}
	if instance.Tags != nil {
		var tags []interface{}
		for _, tag := range instance.Tags {
			tagMap := map[string]interface{}{}
			if tag.TagKey != nil {
				tagMap["tag_key"] = *tag.TagKey
			}
			if tag.TagValue != nil {
				tagMap["tag_value"] = *tag.TagValue
			}
			tags = append(tags, tagMap)
		}
		_ = d.Set("tags", tags)
	}
	if instance.XForwardedForConfig != nil {
		config := instance.XForwardedForConfig
		configMap := map[string]interface{}{
			"x_forwarded_for_alb_id_enabled":          helper.PBool(config.XForwardedForAlbIdEnabled),
			"x_forwarded_for_client_src_port_enabled": helper.PBool(config.XForwardedForClientSrcPortEnabled),
			"x_forwarded_for_host_enabled":            helper.PBool(config.XForwardedForHostEnabled),
			"x_forwarded_for_port_enabled":            helper.PBool(config.XForwardedForPortEnabled),
			"x_forwarded_for_proto_enabled":           helper.PBool(config.XForwardedForProtoEnabled),
			"x_tencent_client_idn_enabled":            helper.PBool(config.XTencentClientIDNEnabled),
			"x_tencent_client_sdn_enabled":            helper.PBool(config.XTencentClientSDNEnabled),
			"x_tencent_client_serial_enabled":         helper.PBool(config.XTencentClientSerialEnabled),
			"x_tencent_client_verify_enabled":         helper.PBool(config.XTencentClientVerifyEnabled),
		}
		if config.XForwardedForMode != nil {
			configMap["x_forwarded_for_mode"] = *config.XForwardedForMode
		}
		_ = d.Set("x_forwarded_for_config", []interface{}{configMap})
	}

	return nil
}

func resourceTencentCloudAlbListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_listener.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	listenerId := items[len(items)-1]
	loadBalancerId := d.Get("load_balancer_id").(string)
	if len(items) == 2 {
		listenerId = items[1]
		loadBalancerId = items[0]
	}

	if d.HasChangeExcept("load_balancer_id") && d.HasChangeExcept("listener_port") && d.HasChangeExcept("listener_protocol") {
		request := alb.NewModifyListenerAttributesRequest()
		request.ListenerId = helper.String(listenerId)
		request.LoadBalancerId = helper.String(loadBalancerId)
		if d.HasChange("default_actions") {
			request.DefaultActions = buildAlbDefaultActions(d)
		}
		if d.HasChange("ca_certificate_ids") {
			request.CaCertificateIds = buildAlbCaCertificateIds(d)
		}
		if d.HasChange("ca_enabled") {
			if v, ok := d.GetOkExists("ca_enabled"); ok {
				request.CaEnabled = helper.Bool(v.(bool))
			}
		}
		if d.HasChange("certificate_ids") {
			request.CertificateIds = buildAlbCertificateIds(d)
		}
		if d.HasChange("gzip_enabled") {
			if v, ok := d.GetOkExists("gzip_enabled"); ok {
				request.GzipEnabled = helper.Bool(v.(bool))
			}
		}
		if d.HasChange("http2enabled") {
			if v, ok := d.GetOkExists("http2enabled"); ok {
				request.Http2Enabled = helper.Bool(v.(bool))
			}
		}
		if d.HasChange("idle_timeout") {
			if v, ok := d.GetOkExists("idle_timeout"); ok {
				request.IdleTimeout = helper.IntUint64(v.(int))
			}
		}
		if d.HasChange("listener_name") {
			if v, ok := d.GetOk("listener_name"); ok {
				request.ListenerName = helper.String(v.(string))
			}
		}
		if d.HasChange("request_timeout") {
			if v, ok := d.GetOkExists("request_timeout"); ok {
				request.RequestTimeout = helper.IntUint64(v.(int))
			}
		}
		if d.HasChange("security_policy_id") {
			if v, ok := d.GetOk("security_policy_id"); ok {
				request.SecurityPolicyId = helper.String(v.(string))
			}
		}
		if d.HasChange("x_forwarded_for_config") {
			request.XForwardedForConfig = buildAlbXForwardedForConfig(d)
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAlbClient().ModifyListenerAttributesWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("modify ALB listener attributes failed, Response is nil"))
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update ALB listener failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudAlbListenerRead(d, meta)
}

func resourceTencentCloudAlbListenerDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_listener.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	listenerId := items[len(items)-1]
	loadBalancerId := d.Get("load_balancer_id").(string)
	if len(items) == 2 {
		listenerId = items[1]
		loadBalancerId = items[0]
	}

	request := alb.NewDeleteListenerRequest()
	request.LoadBalancerId = helper.String(loadBalancerId)
	request.ListenerIds = []*string{helper.String(listenerId)}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAlbClient().DeleteListenerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("delete ALB listener failed, Response is nil"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete ALB listener failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
