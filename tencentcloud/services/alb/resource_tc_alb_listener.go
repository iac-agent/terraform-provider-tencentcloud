package alb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
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
				Description: "ID of the ALB instance, in the format `alb-` followed by 8 alphanumeric characters.",
			},
			"listener_port": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Port used by the frontend of the load balancer instance. Valid values: 1~65535.",
			},
			"listener_protocol": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Listener protocol. Valid values: `HTTP`, `HTTPS` or `QUIC`.",
			},
			"default_actions": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "List of default forwarding rule actions. Currently a listener only supports adding 1 default forwarding rule action.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Forwarding action type. When creating a listener, the default forwarding action type only supports forwarding to target group.",
						},
						"target_group_config": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "Forwarding target group configuration. When creating a listener, the target group configuration in the forwarding action only supports a single target group.",
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
													Description: "Target group ID, in the format `lbtg-` followed by 8 alphanumeric characters.",
												},
												"weight": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Weight, valid range: [0, 100], default is 10.",
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
													Description: "Whether to enable sticky session, disabled by default.",
												},
												"timeout": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Timeout in seconds, valid range: 1~86400, default value: 1000.",
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
				Description: "List of CA certificate IDs configured for the listener. Currently a listener only supports adding 1 CA certificate. Required when ca_enabled is true.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ca_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable mutual authentication. Valid values: true (enabled), false (disabled, default).",
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
				Description: "Whether to enable Gzip compression. Valid values: true (default), false.",
			},
			"http2enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable HTTP/2. HTTP protocol defaults to false, HTTPS protocol defaults to true. Only HTTPS protocol supports this parameter.",
			},
			"idle_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Connection idle timeout in seconds. Valid range: 1~600. Default value: 15.",
			},
			"listener_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom listener name. Length is 1~255 characters; it can contain Chinese characters, letters, digits, hyphens (-), slashes (/), periods (.) and underscores (_).",
			},
			"request_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Connection request timeout in seconds. Valid values: 1~600. Default value: 60.",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Security policy ID, in the format `tls-` followed by 8 alphanumeric characters.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag list. Up to 20 tags are supported.",
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
							Description: "Whether to obtain the load balancer instance ID through the ALB-ID header field.",
						},
						"x_forwarded_for_client_src_port_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to obtain the client port accessing the load balancer instance through the X-Forwarded-Client-srcport header field.",
						},
						"x_forwarded_for_host_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable obtaining the client domain accessing the load balancer instance through the X-Forwarded-Host header field.",
						},
						"x_forwarded_for_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies how to handle the X-Forwarded-For (XFF) HTTP header field. Valid values: `append` (default), `remove`, `passthrough`.",
						},
						"x_forwarded_for_port_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to obtain the listener port of the load balancer instance through the X-Forwarded-Port header field.",
						},
						"x_forwarded_for_proto_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to obtain the listener protocol of the load balancer instance through the X-Forwarded-Proto header field.",
						},
						"x_tencent_client_idn_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to access the issuer of the client certificate $ssl_client_i_dn through the X-Tencent-Client-IDN header.",
						},
						"x_tencent_client_sdn_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to access the subject of the client certificate $ssl_client_s_dn through the X-Tencent-Client-SDN header.",
						},
						"x_tencent_client_serial_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to access the serial number of the client certificate $ssl_client_serial through the X-Tencent-Client-Serial header.",
						},
						"x_tencent_client_verify_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to access the verification result of the client certificate $ssl_client_verify through the X-Tencent-Client-Verify header.",
						},
					},
				},
			},
		},
	}
}

func BuildAlbListenerDefaultActionsForTest(d *schema.ResourceData) []interface{} {
	return buildAlbListenerDefaultActions(d)
}

func BuildAlbListenerXForwardedForConfigForTest(d *schema.ResourceData) map[string]interface{} {
	return buildAlbListenerXForwardedForConfig(d)
}

func buildAlbListenerDefaultActions(d *schema.ResourceData) []interface{} {
	defaultActions := []interface{}{}
	if v, ok := d.GetOk("default_actions"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			action := map[string]interface{}{
				"Type": dMap["type"].(string),
			}
			if tgc, ok := dMap["target_group_config"].([]interface{}); ok && len(tgc) > 0 {
				tgcMap := tgc[0].(map[string]interface{})
				targetGroupConfig := map[string]interface{}{}
				if tgs, ok := tgcMap["target_groups"].([]interface{}); ok && len(tgs) > 0 {
					targetGroups := []interface{}{}
					for _, tg := range tgs {
						tgMap := tg.(map[string]interface{})
						targetGroup := map[string]interface{}{
							"TargetGroupId": tgMap["target_group_id"].(string),
						}
						if weight, ok := tgMap["weight"]; ok && weight.(int) > 0 {
							targetGroup["Weight"] = weight.(int)
						}
						targetGroups = append(targetGroups, targetGroup)
					}
					targetGroupConfig["TargetGroups"] = targetGroups
				}
				if tgs, ok := tgcMap["target_group_sticky_session"].([]interface{}); ok && len(tgs) > 0 {
					tgsMap := tgs[0].(map[string]interface{})
					stickySession := map[string]interface{}{
						"Enabled": tgsMap["enabled"].(bool),
					}
					if timeout, ok := tgsMap["timeout"]; ok && timeout.(int) > 0 {
						stickySession["Timeout"] = timeout.(int)
					}
					targetGroupConfig["TargetGroupStickySession"] = stickySession
				}
				action["TargetGroupConfig"] = targetGroupConfig
			}
			defaultActions = append(defaultActions, action)
		}
	}
	return defaultActions
}

func buildAlbListenerXForwardedForConfig(d *schema.ResourceData) map[string]interface{} {
	if v, ok := d.GetOk("x_forwarded_for_config"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			xffConfig := map[string]interface{}{}
			if vv, ok := dMap["x_forwarded_for_alb_id_enabled"]; ok {
				xffConfig["XForwardedForAlbIdEnabled"] = vv.(bool)
			}
			if vv, ok := dMap["x_forwarded_for_client_src_port_enabled"]; ok {
				xffConfig["XForwardedForClientSrcPortEnabled"] = vv.(bool)
			}
			if vv, ok := dMap["x_forwarded_for_host_enabled"]; ok {
				xffConfig["XForwardedForHostEnabled"] = vv.(bool)
			}
			if vv, ok := dMap["x_forwarded_for_mode"]; ok && vv.(string) != "" {
				xffConfig["XForwardedForMode"] = vv.(string)
			}
			if vv, ok := dMap["x_forwarded_for_port_enabled"]; ok {
				xffConfig["XForwardedForPortEnabled"] = vv.(bool)
			}
			if vv, ok := dMap["x_forwarded_for_proto_enabled"]; ok {
				xffConfig["XForwardedForProtoEnabled"] = vv.(bool)
			}
			if vv, ok := dMap["x_tencent_client_idn_enabled"]; ok {
				xffConfig["XTencentClientIDNEnabled"] = vv.(bool)
			}
			if vv, ok := dMap["x_tencent_client_sdn_enabled"]; ok {
				xffConfig["XTencentClientSDNEnabled"] = vv.(bool)
			}
			if vv, ok := dMap["x_tencent_client_serial_enabled"]; ok {
				xffConfig["XTencentClientSerialEnabled"] = vv.(bool)
			}
			if vv, ok := dMap["x_tencent_client_verify_enabled"]; ok {
				xffConfig["XTencentClientVerifyEnabled"] = vv.(bool)
			}
			return xffConfig
		}
	}
	return nil
}

func resourceTencentCloudAlbListenerCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_listener.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	loadBalancerId := d.Get("load_balancer_id").(string)

	body := map[string]interface{}{}
	body["DefaultActions"] = buildAlbListenerDefaultActions(d)
	if v, ok := d.GetOkExists("listener_port"); ok {
		body["ListenerPort"] = v.(int)
	}
	if v, ok := d.GetOk("listener_protocol"); ok {
		body["ListenerProtocol"] = v.(string)
	}
	body["LoadBalancerId"] = loadBalancerId

	if v, ok := d.GetOk("ca_certificate_ids"); ok {
		ids := []string{}
		for _, item := range v.([]interface{}) {
			ids = append(ids, item.(string))
		}
		body["CaCertificateIds"] = ids
	}
	if v, ok := d.GetOkExists("ca_enabled"); ok {
		body["CaEnabled"] = v.(bool)
	}
	if v, ok := d.GetOk("certificate_ids"); ok {
		ids := []string{}
		for _, item := range v.([]interface{}) {
			ids = append(ids, item.(string))
		}
		body["CertificateIds"] = ids
	}
	if v, ok := d.GetOk("client_token"); ok {
		body["ClientToken"] = v.(string)
	}
	if v, ok := d.GetOkExists("gzip_enabled"); ok {
		body["GzipEnabled"] = v.(bool)
	}
	if v, ok := d.GetOkExists("http2enabled"); ok {
		body["Http2Enabled"] = v.(bool)
	}
	if v, ok := d.GetOkExists("idle_timeout"); ok {
		body["IdleTimeout"] = v.(int)
	}
	if v, ok := d.GetOk("listener_name"); ok {
		body["ListenerName"] = v.(string)
	}
	if v, ok := d.GetOkExists("request_timeout"); ok {
		body["RequestTimeout"] = v.(int)
	}
	if v, ok := d.GetOk("security_policy_id"); ok {
		body["SecurityPolicyId"] = v.(string)
	}
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagsList := []map[string]interface{}{}
		for k, v := range tags {
			tagsList = append(tagsList, map[string]interface{}{
				"TagKey":   k,
				"TagValue": fmt.Sprintf("%v", v),
			})
		}
		body["Tags"] = tagsList
	}
	if xff := buildAlbListenerXForwardedForConfig(d); xff != nil {
		body["XForwardedForConfig"] = xff
	}

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOmitNilClient("alb")
	request := tchttp.NewCommonRequest("alb", "2025-10-30", "CreateListener")
	err := request.SetActionParameters(body)
	if err != nil {
		return err
	}

	response := tchttp.NewCommonResponse()
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := client.Send(request, response)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), string(request.GetBody()), response.GetBody())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create alb listener failed, reason:%+v", logId, err)
		return err
	}

	var respData struct {
		Response struct {
			ListenerId string `json:"ListenerId"`
		} `json:"Response"`
	}
	if err = json.Unmarshal(response.GetBody(), &respData); err != nil {
		return fmt.Errorf("create alb listener failed, unmarshal response body failed: %v", err)
	}
	if respData.Response.ListenerId == "" {
		return fmt.Errorf("create alb listener failed, ListenerId is empty")
	}
	d.SetId(strings.Join([]string{loadBalancerId, respData.Response.ListenerId}, tccommon.FILED_SP))

	return resourceTencentCloudAlbListenerRead(d, meta)
}

func resourceTencentCloudAlbListenerRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_listener.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	loadBalancerId := idSplit[0]
	listenerId := idSplit[1]

	body := map[string]interface{}{
		"LoadBalancerId": loadBalancerId,
		"ListenerIds":    []string{listenerId},
	}

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOmitNilClient("alb")
	request := tchttp.NewCommonRequest("alb", "2025-10-30", "DescribeListeners")
	err := request.SetActionParameters(body)
	if err != nil {
		return err
	}

	response := tchttp.NewCommonResponse()
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		e := client.Send(request, response)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), string(request.GetBody()), response.GetBody())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s describe alb listener failed, reason:%+v", logId, err)
		return err
	}

	var respData struct {
		Response struct {
			Listeners []struct {
				CaEnable            bool                   `json:"CaEnable"`
				GzipEnabled         bool                   `json:"GzipEnabled"`
				Http2Enable         bool                   `json:"Http2Enable"`
				IdleTimeout         int64                  `json:"IdleTimeout"`
				ListenerId          string                 `json:"ListenerId"`
				ListenerName        string                 `json:"ListenerName"`
				ListenerPort        int64                  `json:"ListenerPort"`
				ListenerProtocol    string                 `json:"ListenerProtocol"`
				RequestTimeout      int64                  `json:"RequestTimeout"`
				TlsSecurityPolicyId string                 `json:"TlsSecurityPolicyId"`
				XForwardedForConfig map[string]interface{} `json:"XForwardedForConfig"`
			} `json:"Listeners"`
		} `json:"Response"`
	}
	if err = json.Unmarshal(response.GetBody(), &respData); err != nil {
		return fmt.Errorf("describe alb listener failed, unmarshal response body failed: %v", err)
	}
	if len(respData.Response.Listeners) == 0 {
		log.Printf("[WARN]%s resource `alb_listener` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}
	listener := respData.Response.Listeners[0]

	_ = d.Set("load_balancer_id", loadBalancerId)
	_ = d.Set("listener_port", listener.ListenerPort)
	_ = d.Set("listener_protocol", listener.ListenerProtocol)
	_ = d.Set("ca_enabled", listener.CaEnable)
	_ = d.Set("gzip_enabled", listener.GzipEnabled)
	_ = d.Set("http2enabled", listener.Http2Enable)
	_ = d.Set("idle_timeout", listener.IdleTimeout)
	_ = d.Set("listener_name", listener.ListenerName)
	_ = d.Set("request_timeout", listener.RequestTimeout)
	_ = d.Set("security_policy_id", listener.TlsSecurityPolicyId)

	if listener.XForwardedForConfig != nil {
		xffMap := map[string]interface{}{}
		if v, ok := listener.XForwardedForConfig["XForwardedForAlbIdEnabled"]; ok {
			xffMap["x_forwarded_for_alb_id_enabled"] = v
		}
		if v, ok := listener.XForwardedForConfig["XForwardedForClientSrcPortEnabled"]; ok {
			xffMap["x_forwarded_for_client_src_port_enabled"] = v
		}
		if v, ok := listener.XForwardedForConfig["XForwardedForHostEnabled"]; ok {
			xffMap["x_forwarded_for_host_enabled"] = v
		}
		if v, ok := listener.XForwardedForConfig["XForwardedForMode"]; ok {
			xffMap["x_forwarded_for_mode"] = v
		}
		if v, ok := listener.XForwardedForConfig["XForwardedForPortEnabled"]; ok {
			xffMap["x_forwarded_for_port_enabled"] = v
		}
		if v, ok := listener.XForwardedForConfig["XForwardedForProtoEnabled"]; ok {
			xffMap["x_forwarded_for_proto_enabled"] = v
		}
		if v, ok := listener.XForwardedForConfig["XTencentClientIDNEnabled"]; ok {
			xffMap["x_tencent_client_idn_enabled"] = v
		}
		if v, ok := listener.XForwardedForConfig["XTencentClientSDNEnabled"]; ok {
			xffMap["x_tencent_client_sdn_enabled"] = v
		}
		if v, ok := listener.XForwardedForConfig["XTencentClientSerialEnabled"]; ok {
			xffMap["x_tencent_client_serial_enabled"] = v
		}
		if v, ok := listener.XForwardedForConfig["XTencentClientVerifyEnabled"]; ok {
			xffMap["x_tencent_client_verify_enabled"] = v
		}
		_ = d.Set("x_forwarded_for_config", []interface{}{xffMap})
	}

	return nil
}

func resourceTencentCloudAlbListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_listener.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	loadBalancerId := idSplit[0]
	listenerId := idSplit[1]

	body := map[string]interface{}{
		"ListenerId":     listenerId,
		"LoadBalancerId": loadBalancerId,
	}

	if d.HasChange("default_actions") {
		body["DefaultActions"] = buildAlbListenerDefaultActions(d)
	}
	if d.HasChange("ca_certificate_ids") {
		if v, ok := d.GetOk("ca_certificate_ids"); ok {
			ids := []string{}
			for _, item := range v.([]interface{}) {
				ids = append(ids, item.(string))
			}
			body["CaCertificateIds"] = ids
		}
	}
	if d.HasChange("ca_enabled") {
		body["CaEnabled"] = d.Get("ca_enabled").(bool)
	}
	if d.HasChange("certificate_ids") {
		if v, ok := d.GetOk("certificate_ids"); ok {
			ids := []string{}
			for _, item := range v.([]interface{}) {
				ids = append(ids, item.(string))
			}
			body["CertificateIds"] = ids
		}
	}
	if d.HasChange("client_token") {
		if v, ok := d.GetOk("client_token"); ok {
			body["ClientToken"] = v.(string)
		}
	}
	if d.HasChange("gzip_enabled") {
		body["GzipEnabled"] = d.Get("gzip_enabled").(bool)
	}
	if d.HasChange("http2enabled") {
		body["Http2Enabled"] = d.Get("http2enabled").(bool)
	}
	if d.HasChange("idle_timeout") {
		body["IdleTimeout"] = d.Get("idle_timeout").(int)
	}
	if d.HasChange("listener_name") {
		if v, ok := d.GetOk("listener_name"); ok {
			body["ListenerName"] = v.(string)
		}
	}
	if d.HasChange("request_timeout") {
		body["RequestTimeout"] = d.Get("request_timeout").(int)
	}
	if d.HasChange("security_policy_id") {
		if v, ok := d.GetOk("security_policy_id"); ok {
			body["SecurityPolicyId"] = v.(string)
		}
	}
	if d.HasChange("x_forwarded_for_config") {
		if xff := buildAlbListenerXForwardedForConfig(d); xff != nil {
			body["XForwardedForConfig"] = xff
		}
	}

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOmitNilClient("alb")
	request := tchttp.NewCommonRequest("alb", "2025-10-30", "ModifyListenerAttributes")
	err := request.SetActionParameters(body)
	if err != nil {
		return err
	}

	response := tchttp.NewCommonResponse()
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := client.Send(request, response)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), string(request.GetBody()), response.GetBody())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update alb listener failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		resourceName := tccommon.BuildTagResourceName("alb", "listener", meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region, listenerId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudAlbListenerRead(d, meta)
}

func resourceTencentCloudAlbListenerDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_listener.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	loadBalancerId := idSplit[0]
	listenerId := idSplit[1]

	body := map[string]interface{}{
		"ListenerIds":    []string{listenerId},
		"LoadBalancerId": loadBalancerId,
	}
	if v, ok := d.GetOk("client_token"); ok {
		body["ClientToken"] = v.(string)
	}

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOmitNilClient("alb")
	request := tchttp.NewCommonRequest("alb", "2025-10-30", "DeleteListener")
	err := request.SetActionParameters(body)
	if err != nil {
		return err
	}

	response := tchttp.NewCommonResponse()
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := client.Send(request, response)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), string(request.GetBody()), response.GetBody())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete alb listener failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
