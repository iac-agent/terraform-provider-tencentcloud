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

func DataSourceTencentCloudClbListenerRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbListenerRulesRead,

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "需要查询的CLB ID。",
			},
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "需要查询的CLB监听ID。",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "需要查询的转发规则ID。",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "需要查询的转发规则的域名。",
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "需要查询的转发规则的url。",
			},
			"scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CLB_LISTENER_SCHEDULER),
				Description: "CLB监听转发规则的调度方式，可用值包括`WRR`、`IP HASH`和`LEAST_CONN`。默认为“WRR”。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"rule_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "监听器的转发规则列表。每个元素包含以下属性：",
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
						"domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "规则的域。",
						},
						"url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "规则的 URL。",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "规则的 ID。",
						},
						"health_check_switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否开启健康检查。",
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
							Description: "健康检查的不健康阈值，默认为`3`。如果健康检查连续3次返回成功，则判定该云服务器不健康。取值范围为2-10。注意：TCP/UDP/TCP_SSL监听可以直接配置，HTTP/HTTPS监听需要在tencentcloud_clb_listener_rule中配置。",
						},
						"health_check_http_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "HTTP 状态代码。默认值为 31，取值范围为 1-31。 1表示返回值'1xx'是健康。 2表示返回值'2xx'是健康。 4 表示返回值“3xx”是健康状况。 8表示返回值4xx是健康值。 16表示返回值“5xx”是健康。如果想要多个返回码来指示健康状况，需要添加相应的值。注意：“TCP”监听器的“HTTP”健康检查仅支持指定一个健康检查状态码。注意：仅支持“HTTP”和“HTTPS”协议的侦听器。",
						},
						"health_check_http_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "健康检查路径。注意：仅支持“HTTPS”和“HTTP”协议的侦听器。",
						},
						"health_check_http_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "健康检查域名。注意：仅支持“HTTPS”和“HTTP”协议的侦听器。",
						},
						"health_check_http_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "健康检查方法。注意：仅支持“HTTPS”和“HTTP”协议的侦听器。默认为'HEAD'，可用值包括'HEAD'和'GET'。",
						},
						"certificate_ssl_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SSL 模式类型，可用值包括“UNIDIRECTIONAL”、“MUTUAL”。注：仅支持“HTTPS”和“TCP_SSL”协议的侦听器。",
						},
						"certificate_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "服务器证书的 ID。注意：仅支持“HTTPS”和“TCP_SSL”协议的侦听器。",
						},
						"certificate_ca_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "客户端证书的ID。注意：仅支持“HTTPS”和“TCP_SSL”协议的侦听器。",
						},
						"session_expire_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CLB 侦听器中的会话持续时间。注意：当调度程序指定为“WRR”时可用。注意：TCP/UDP/TCP_SSL监听可以直接配置，HTTP/HTTPS监听需要在tencentcloud_clb_listener_rule中配置。",
						},
						"scheduler": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB监听的调度方式，可用值包括'WRR'、'IP_HASH'和'LEAST_CONN'。默认值为“WRR”。注意：TCP/UDP/TCP_SSL监听可以直接配置，HTTP/HTTPS监听需要在tencentcloud_clb_listener_rule中配置。",
						},
						"http2_switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否设置HTTP2协议。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudClbListenerRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clb_listener_rules.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	listenerId := d.Get("listener_id").(string)
	checkErr := ListenerIdCheck(listenerId)
	if checkErr != nil {
		return checkErr
	}
	clbId := d.Get("clb_id").(string)
	params := make(map[string]string)
	params["clb_id"] = clbId
	params["listener_id"] = listenerId
	if v, ok := d.GetOk("clb_id"); ok {
		params["clb_id"] = v.(string)
	}
	if v, ok := d.GetOk("scheduler"); ok {
		params["scheduler"] = v.(string)
	}
	if v, ok := d.GetOk("rule_id"); ok {
		params["rule_id"] = v.(string)
		checkErr := RuleIdCheck(params["rule_id"])
		if checkErr != nil {
			return checkErr
		}
	}
	if v, ok := d.GetOk("domain"); ok {
		params["domain"] = v.(string)
	}
	if v, ok := d.GetOk("url"); ok {
		params["url"] = v.(string)
	}

	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var rules []*clb.RuleOutput
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := clbService.DescribeRulesByFilter(ctx, params)
		if e != nil {
			return tccommon.RetryError(e)
		}
		rules = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB listener rules failed, reason:%+v", logId, err)
		return err
	}
	ruleList := make([]map[string]interface{}, 0, len(rules))
	ids := make([]string, 0, len(rules))
	for _, rule := range rules {
		mapping := map[string]interface{}{
			"clb_id":              clbId,
			"listener_id":         listenerId,
			"rule_id":             rule.LocationId,
			"domain":              rule.Domain,
			"url":                 rule.Url,
			"session_expire_time": rule.SessionExpireTime,
			"scheduler":           rule.Scheduler,
			"http2_switch":        rule.Http2,
		}
		if rule.HealthCheck != nil {
			healthCheckSwitch := false
			if *rule.HealthCheck.HealthSwitch == int64(1) {
				healthCheckSwitch = true
			}
			mapping["health_check_switch"] = healthCheckSwitch
			mapping["health_check_interval_time"] = *rule.HealthCheck.IntervalTime
			mapping["health_check_health_num"] = *rule.HealthCheck.HealthNum
			mapping["health_check_unhealth_num"] = *rule.HealthCheck.UnHealthNum
			mapping["health_check_http_code"] = *rule.HealthCheck.HttpCode
			mapping["health_check_http_method"] = *rule.HealthCheck.HttpCheckMethod
			mapping["health_check_http_domain"] = *rule.HealthCheck.HttpCheckDomain
			mapping["health_check_http_path"] = *rule.HealthCheck.HttpCheckPath
		}
		if rule.Certificate != nil {
			mapping["certificate_ssl_mode"] = *rule.Certificate.SSLMode
			mapping["certificate_id"] = *rule.Certificate.CertId
			if rule.Certificate.CertCaId != nil {
				mapping["certificate_ca_id"] = *rule.Certificate.CertCaId
			}
		}
		ruleList = append(ruleList, mapping)
		ids = append(ids, *rule.LocationId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("rule_list", ruleList); e != nil {
		log.Printf("[CRITAL]%s provider set CLB listener rule list fail, reason:%+v", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), ruleList); e != nil {
			return e
		}
	}

	return nil
}
