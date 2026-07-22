package gaap

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudGaapHttpRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapHttpRulesRead,
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID layer7 listener to be queried。",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Forward 域名 of the layer7 listener to be queried。",
			},
			"path": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringPrefix("/"),
				Description:  "路径 of the forward rule to be queried。",
			},
			"forward_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Requested 主机 which is forwarded to the realserver by the listener to be queried。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// computed
			"rules": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "An information 列表 forward rule of the layer7 listeners. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID forward rule。",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID layer7 listener。",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Forward 域名 of the forward rule。",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "路径 of the forward rule。",
						},
						"realserver_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 realserver。",
						},
						"scheduler": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheduling policy of the forward rule。",
						},
						"health_check": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "表示是否health check is enable。",
						},
						"interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Interval of the health check。",
						},
						"connect_timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Timeout of the health check response。",
						},
						"health_check_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "路径 of health check。",
						},
						"health_check_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Method of the health check。",
						},
						"health_check_status_codes": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Description: "Return 代码 of confirmed normal。",
						},
						"forward_host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Requested 主机 which is forwarded to the realserver by the listener。",
						},
						"sni_switch": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ServerNameIndication (SNI) switch。",
						},
						"sni": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ServerNameIndication (SNI)。",
						},
						"realservers": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "An information 列表 GAAP realserver. Each element 包含following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID GAAP realserver。",
									},
									"ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IP of the GAAP realserver。",
									},
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "域名 of the GAAP realserver。",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "端口 of the GAAP realserver。",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Scheduling 权重",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "状态 GAAP realserver。",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudGaapHttpRulesRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_gaap_http_rules.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	listenerId := d.Get("listener_id").(string)

	var (
		domain      string
		path        *string
		forwardHost *string
		ids         []string
		rules       []map[string]interface{}
	)

	if raw, ok := d.GetOk("domain"); ok {
		domain = raw.(string)
	}
	if raw, ok := d.GetOk("path"); ok {
		path = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("forward_host"); ok {
		forwardHost = helper.String(raw.(string))
	}

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	domainRuleSets, err := service.DescribeDomains(ctx, listenerId, domain)
	if err != nil {
		return err
	}

	ids = make([]string, 0, len(domainRuleSets))
	rules = make([]map[string]interface{}, 0, len(domainRuleSets))

	for _, domainRule := range domainRuleSets {
		for _, rule := range domainRule.RuleSet {
			if nilFields := tccommon.CheckNil(rule, map[string]string{
				"RuleId":         "id",
				"Path":           "path",
				"RealServerType": "realserver type",
				"Scheduler":      "scheduler",
				"HealthCheck":    "health check",
				"CheckParams":    "health check params",
				"ForwardHost":    "forward host",
			}); len(nilFields) > 0 {
				return fmt.Errorf("rule %v are nil", nilFields)
			}

			checkParams := rule.CheckParams

			if nilFields := tccommon.CheckNil(checkParams, map[string]string{
				"DelayLoop":      "interval",
				"ConnectTimeout": "connect timeout",
				"Path":           "path",
				"Method":         "method",
			}); len(nilFields) > 0 {
				return fmt.Errorf("rule health check %v are nil", nilFields)
			}

			if len(checkParams.StatusCode) == 0 {
				return errors.New("rule health check status codes set is empty")
			}

			if path != nil && *rule.Path != *path {
				continue
			}

			if forwardHost != nil && *forwardHost != *rule.ForwardHost {
				continue
			}

			ids = append(ids, *rule.RuleId)

			m := map[string]interface{}{
				"id":                  *rule.RuleId,
				"listener_id":         listenerId,
				"domain":              *domainRule.Domain,
				"path":                *rule.Path,
				"realserver_type":     *rule.RealServerType,
				"scheduler":           *rule.Scheduler,
				"health_check":        *rule.HealthCheck == 1,
				"interval":            *checkParams.DelayLoop,
				"connect_timeout":     *checkParams.ConnectTimeout,
				"health_check_path":   *checkParams.Path,
				"health_check_method": *checkParams.Method,
				"forward_host":        *rule.ForwardHost,
				"sni_switch":          *rule.ServerNameIndicationSwitch,
				"sni":                 *rule.ServerNameIndication,
			}
			statusCodes := make([]int, 0, len(checkParams.StatusCode))
			for _, code := range checkParams.StatusCode {
				statusCodes = append(statusCodes, int(*code))
			}
			m["health_check_status_codes"] = statusCodes

			realservers := make([]map[string]interface{}, 0, len(rule.RealServerSet))
			for _, rs := range rule.RealServerSet {
				if nilFields := tccommon.CheckNil(rs, map[string]string{
					"RealServerId":     "id",
					"RealServerIP":     "ip or domain",
					"RealServerPort":   "port",
					"RealServerWeight": "weight",
					"RealServerStatus": "status",
				}); len(nilFields) > 0 {
					return fmt.Errorf("realserver %v are nil", nilFields)
				}

				realserver := map[string]interface{}{
					"id":     *rs.RealServerId,
					"port":   *rs.RealServerPort,
					"weight": *rs.RealServerWeight,
					"status": *rs.RealServerStatus,
				}

				if net.ParseIP(*rs.RealServerIP) == nil {
					realserver["domain"] = *rs.RealServerIP
				} else {
					realserver["ip"] = *rs.RealServerIP
				}

				realservers = append(realservers, realserver)
			}

			m["realservers"] = realservers

			rules = append(rules, m)
		}
	}

	_ = d.Set("rules", rules)
	d.SetId(helper.DataResourceIdsHash(ids))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), rules); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%v]",
				logId, output.(string), err)
			return err
		}
	}

	return nil
}
