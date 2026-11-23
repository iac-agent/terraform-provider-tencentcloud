package alb

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	albv20251030 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/alb/v20251030"

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
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the load balancer instance, in the format of `alb-` followed by 8 alphanumeric characters.",
			},
			"listener_port": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Port used by the frontend of the load balancer instance. Range: 1~65535.",
			},
			"listener_protocol": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Listener protocol. Valid values: `HTTP`, `HTTPS`, `QUIC`.",
			},
			"listener_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom listener name. Length: 1~255 characters. It can contain Chinese characters, letters, digits, `-`, `/`, `.` and `_`.",
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
							Description: "Forwarding action type. When creating a listener, only forwarding to target group is supported.",
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
													Description: "Target group ID, in the format of `lbtg-` followed by 8 alphanumeric characters.",
												},
												"weight": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Weight. Range: [0, 100]. Default is 10.",
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
													Description: "Whether to enable sticky session. Default is false.",
												},
												"timeout": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Timeout in seconds. Range: 1~86400. Default is 1000.",
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
			"certificate_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Server certificate ID list.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ca_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable mutual authentication. Default is false.",
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
				Description: "Security policy ID, in the format of `tls-` followed by 8 alphanumeric characters.",
			},
			"idle_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Connection idle timeout in seconds. Range: 1~600. Default is 15.",
			},
			"request_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Connection request timeout in seconds. Range: 1~600. Default is 60.",
			},
			"gzip_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable Gzip compression. Default is true.",
			},
			"http2enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable HTTP/2. HTTP protocol defaults to false, HTTPS protocol defaults to true. Only supported by HTTPS protocol.",
			},
			"x_forwarded_for_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "X-Forwarded-For configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"x_forwarded_for_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies how to handle the X-Forwarded-For (XFF) HTTP header field. Valid values: `append` (default), `remove`, `passthrough`.",
						},
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
		},
	}
}

func buildAlbDefaultActions(d *schema.ResourceData) []*albv20251030.DefaultAction {
	if len(d.Get("default_actions").([]interface{})) == 0 {
		return nil
	}
	actions := make([]*albv20251030.DefaultAction, 0, 1)
	actionMap, ok := helper.InterfacesHeadMap(d, "default_actions")
	if !ok {
		return nil
	}
	action := albv20251030.DefaultAction{}
	if v, ok := actionMap["type"]; ok {
		action.Type = helper.String(v.(string))
	}
	if tgConfigList, ok := helper.InterfacesHeadMap(d, "default_actions.0.target_group_config"); ok {
		tgConfig := albv20251030.TargetGroupConfig{}
		if tgList, ok := tgConfigList["target_groups"].([]interface{}); ok {
			tgs := make([]*albv20251030.TargetGroupTuple, 0, len(tgList))
			for _, tg := range tgList {
				tgMap, mapOk := tg.(map[string]interface{})
				if !mapOk {
					continue
				}
				tgt := albv20251030.TargetGroupTuple{}
				if v, ok := tgMap["target_group_id"]; ok {
					tgt.TargetGroupId = helper.String(v.(string))
				}
				if v, ok := tgMap["weight"]; ok {
					tgt.Weight = helper.IntInt64(v.(int))
				}
				tgs = append(tgs, &tgt)
			}
			tgConfig.TargetGroups = tgs
		}
		if stickyMap, ok := helper.InterfacesHeadMap(d, "default_actions.0.target_group_config.0.target_group_sticky_session"); ok {
			sticky := albv20251030.TargetGroupStickySession{}
			if v, ok := stickyMap["enabled"]; ok {
				sticky.Enabled = helper.Bool(v.(bool))
			}
			if v, ok := stickyMap["timeout"]; ok {
				sticky.Timeout = helper.IntInt64(v.(int))
			}
			tgConfig.TargetGroupStickySession = &sticky
		}
		action.TargetGroupConfig = &tgConfig
	}
	actions = append(actions, &action)
	return actions
}

func buildAlbXForwardedForConfig(d *schema.ResourceData) *albv20251030.XForwardedForConfig {
	xffMap, ok := helper.InterfacesHeadMap(d, "x_forwarded_for_config")
	if !ok {
		return nil
	}
	xff := albv20251030.XForwardedForConfig{}
	if v, ok := xffMap["x_forwarded_for_mode"]; ok {
		xff.XForwardedForMode = helper.String(v.(string))
	}
	if v, ok := xffMap["x_forwarded_for_alb_id_enabled"]; ok {
		xff.XForwardedForAlbIdEnabled = helper.Bool(v.(bool))
	}
	if v, ok := xffMap["x_forwarded_for_client_src_port_enabled"]; ok {
		xff.XForwardedForClientSrcPortEnabled = helper.Bool(v.(bool))
	}
	if v, ok := xffMap["x_forwarded_for_host_enabled"]; ok {
		xff.XForwardedForHostEnabled = helper.Bool(v.(bool))
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

func buildAlbTags(d *schema.ResourceData) []*albv20251030.TagInfo {
	tagList, ok := d.Get("tags").([]interface{})
	if !ok || len(tagList) == 0 {
		return nil
	}
	tags := make([]*albv20251030.TagInfo, 0, len(tagList))
	for _, tag := range tagList {
		tagMap, mapOk := tag.(map[string]interface{})
		if !mapOk {
			continue
		}
		tagInfo := albv20251030.TagInfo{}
		if v, ok := tagMap["tag_key"]; ok {
			tagInfo.TagKey = helper.String(v.(string))
		}
		if v, ok := tagMap["tag_value"]; ok {
			tagInfo.TagValue = helper.String(v.(string))
		}
		tags = append(tags, &tagInfo)
	}
	return tags
}

func resourceTencentCloudAlbListenerCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_listener.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		request  = albv20251030.NewCreateListenerRequest()
		response = albv20251030.NewCreateListenerResponse()
	)

	if v, ok := d.GetOk("load_balancer_id"); ok {
		request.LoadBalancerId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("listener_port"); ok {
		request.ListenerPort = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("listener_protocol"); ok {
		request.ListenerProtocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("listener_name"); ok {
		request.ListenerName = helper.String(v.(string))
	}

	request.DefaultActions = buildAlbDefaultActions(d)

	if v, ok := d.GetOk("certificate_ids"); ok {
		request.CertificateIds = helper.Strings(v.([]string))
	}

	if v, ok := d.GetOkExists("ca_enabled"); ok {
		request.CaEnabled = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("ca_certificate_ids"); ok {
		request.CaCertificateIds = helper.Strings(v.([]string))
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

	request.XForwardedForConfig = buildAlbXForwardedForConfig(d)

	request.Tags = buildAlbTags(d)

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

	if response.Response == nil || response.Response.ListenerId == nil {
		return fmt.Errorf("alb listener id is nil")
	}
	d.SetId(*response.Response.ListenerId)

	return resourceTencentCloudAlbListenerRead(d, meta)
}

func resourceTencentCloudAlbListenerRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_listener.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := AlbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	listenerId := d.Id()

	respData, err := service.DescribeAlbListenerById(ctx, listenerId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `alb_listener` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if respData.ListenerId != nil {
		listenerId = *respData.ListenerId
	}

	if respData.ListenerPort != nil {
		_ = d.Set("listener_port", respData.ListenerPort)
	}

	if respData.ListenerProtocol != nil {
		_ = d.Set("listener_protocol", respData.ListenerProtocol)
	}

	if respData.ListenerName != nil {
		_ = d.Set("listener_name", respData.ListenerName)
	}

	if respData.IdleTimeout != nil {
		_ = d.Set("idle_timeout", respData.IdleTimeout)
	}

	if respData.RequestTimeout != nil {
		_ = d.Set("request_timeout", respData.RequestTimeout)
	}

	if respData.GzipEnabled != nil {
		_ = d.Set("gzip_enabled", respData.GzipEnabled)
	}

	if respData.Http2Enable != nil {
		_ = d.Set("http2enabled", respData.Http2Enable)
	}

	if respData.CaEnable != nil {
		_ = d.Set("ca_enabled", respData.CaEnable)
	}

	if respData.TlsSecurityPolicyId != nil {
		_ = d.Set("security_policy_id", respData.TlsSecurityPolicyId)
	}

	if respData.XForwardedForConfig != nil {
		xff := make(map[string]interface{})
		if respData.XForwardedForConfig.XForwardedForMode != nil {
			xff["x_forwarded_for_mode"] = respData.XForwardedForConfig.XForwardedForMode
		}
		if respData.XForwardedForConfig.XForwardedForAlbIdEnabled != nil {
			xff["x_forwarded_for_alb_id_enabled"] = respData.XForwardedForConfig.XForwardedForAlbIdEnabled
		}
		if respData.XForwardedForConfig.XForwardedForClientSrcPortEnabled != nil {
			xff["x_forwarded_for_client_src_port_enabled"] = respData.XForwardedForConfig.XForwardedForClientSrcPortEnabled
		}
		if respData.XForwardedForConfig.XForwardedForHostEnabled != nil {
			xff["x_forwarded_for_host_enabled"] = respData.XForwardedForConfig.XForwardedForHostEnabled
		}
		if respData.XForwardedForConfig.XForwardedForPortEnabled != nil {
			xff["x_forwarded_for_port_enabled"] = respData.XForwardedForConfig.XForwardedForPortEnabled
		}
		if respData.XForwardedForConfig.XForwardedForProtoEnabled != nil {
			xff["x_forwarded_for_proto_enabled"] = respData.XForwardedForConfig.XForwardedForProtoEnabled
		}
		if respData.XForwardedForConfig.XTencentClientIDNEnabled != nil {
			xff["x_tencent_client_idn_enabled"] = respData.XForwardedForConfig.XTencentClientIDNEnabled
		}
		if respData.XForwardedForConfig.XTencentClientSDNEnabled != nil {
			xff["x_tencent_client_sdn_enabled"] = respData.XForwardedForConfig.XTencentClientSDNEnabled
		}
		if respData.XForwardedForConfig.XTencentClientSerialEnabled != nil {
			xff["x_tencent_client_serial_enabled"] = respData.XForwardedForConfig.XTencentClientSerialEnabled
		}
		if respData.XForwardedForConfig.XTencentClientVerifyEnabled != nil {
			xff["x_tencent_client_verify_enabled"] = respData.XForwardedForConfig.XTencentClientVerifyEnabled
		}
		_ = d.Set("x_forwarded_for_config", []interface{}{xff})
	}

	if respData.Tags != nil {
		tags := make([]interface{}, 0, len(respData.Tags))
		for _, tag := range respData.Tags {
			tagMap := make(map[string]interface{})
			if tag.TagKey != nil {
				tagMap["tag_key"] = tag.TagKey
			}
			if tag.TagValue != nil {
				tagMap["tag_value"] = tag.TagValue
			}
			tags = append(tags, tagMap)
		}
		_ = d.Set("tags", tags)
	}

	return nil
}

func resourceTencentCloudAlbListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_listener.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	listenerId := d.Id()

	loadBalancerId := ""
	if v, ok := d.GetOk("load_balancer_id"); ok {
		loadBalancerId = v.(string)
	}

	request := albv20251030.NewModifyListenerAttributesRequest()
	request.ListenerId = helper.String(listenerId)
	request.LoadBalancerId = helper.String(loadBalancerId)

	needChange := false
	mutableArgs := []string{"listener_name", "default_actions", "certificate_ids", "ca_enabled", "ca_certificate_ids", "security_policy_id", "idle_timeout", "request_timeout", "gzip_enabled", "http2enabled", "x_forwarded_for_config"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		if v, ok := d.GetOk("listener_name"); ok {
			request.ListenerName = helper.String(v.(string))
		}

		request.DefaultActions = buildAlbDefaultActions(d)

		if v, ok := d.GetOk("certificate_ids"); ok {
			request.CertificateIds = helper.Strings(v.([]string))
		}

		if v, ok := d.GetOkExists("ca_enabled"); ok {
			request.CaEnabled = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOk("ca_certificate_ids"); ok {
			request.CaCertificateIds = helper.Strings(v.([]string))
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

		request.XForwardedForConfig = buildAlbXForwardedForConfig(d)

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
	}

	return resourceTencentCloudAlbListenerRead(d, meta)
}

func resourceTencentCloudAlbListenerDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_alb_listener.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	listenerId := d.Id()

	loadBalancerId := ""
	if v, ok := d.GetOk("load_balancer_id"); ok {
		loadBalancerId = v.(string)
	}

	var (
		request  = albv20251030.NewDeleteListenerRequest()
		response = albv20251030.NewDeleteListenerResponse()
	)

	request.ListenerIds = helper.Strings([]string{listenerId})
	request.LoadBalancerId = helper.String(loadBalancerId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAlbV20251030Client().DeleteListenerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete alb listener failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil
}
