package clb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClbTargetHealth() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbTargetHealthRead,
		Schema: map[string]*schema.Schema{
			"load_balancer_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "需要查询的CLB实例ID列表。",
			},

			"load_balancers": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "CLB实例列表。注意：该字段可能返回null，表示取不到有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB实例ID。",
						},
						"load_balancer_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB 实例名称。注意：该字段可能返回null，表示取不到有效值。",
						},
						"listeners": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "听众名单。注意：该字段可能返回null，表示取不到有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"listener_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "听众 ID。",
									},
									"listener_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "听众姓名。注意：该字段可能返回null，表示取不到有效值。",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "监听器协议。",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "侦听器端口。",
									},
									"rules": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "列表 forwarding 规则 的 listener.注意：该字段可能返回null，表示取不到有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"location_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "转发规则ID。",
												},
												"domain": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "转发规则的域名。注意：该字段可能返回null，表示取不到有效值。",
												},
												"url": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "转发规则 URL。注意：该字段可能返回null，表示取不到有效值。",
												},
												"targets": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "该规则绑定的真实服务器的健康状态。注意：该字段可能返回null，表示取不到有效值。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"ip": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "目标的私有IP。",
															},
															"port": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "绑定到目标的端口。",
															},
															"health_status": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "目前的健康状况。真实：健康；假：不健康。",
															},
															"target_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "目标的实例ID，例如ins-12345678。",
															},
															"health_status_detail": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "有关当前健康状况的详细信息。活着：健康；死：例外；未知：检查未启动/正在检查/未知状态。",
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
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudClbTargetHealthRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clb_target_health.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("load_balancer_ids"); ok {
		loadBalancerIdsSet := v.(*schema.Set).List()
		paramMap["LoadBalancerIds"] = helper.InterfacesStringsPoint(loadBalancerIdsSet)
	}

	service := ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var loadBalancers []*clb.LoadBalancerHealth

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClbTargetHealthByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		loadBalancers = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(loadBalancers))
	tmpList := make([]map[string]interface{}, 0, len(loadBalancers))

	if loadBalancers != nil {
		for _, loadBalancerHealth := range loadBalancers {
			loadBalancerHealthMap := map[string]interface{}{}

			if loadBalancerHealth.LoadBalancerId != nil {
				loadBalancerHealthMap["load_balancer_id"] = loadBalancerHealth.LoadBalancerId
			}

			if loadBalancerHealth.LoadBalancerName != nil {
				loadBalancerHealthMap["load_balancer_name"] = loadBalancerHealth.LoadBalancerName
			}

			if loadBalancerHealth.Listeners != nil {
				listenersList := []interface{}{}
				for _, listeners := range loadBalancerHealth.Listeners {
					listenersMap := map[string]interface{}{}

					if listeners.ListenerId != nil {
						listenersMap["listener_id"] = listeners.ListenerId
					}

					if listeners.ListenerName != nil {
						listenersMap["listener_name"] = listeners.ListenerName
					}

					if listeners.Protocol != nil {
						listenersMap["protocol"] = listeners.Protocol
					}

					if listeners.Port != nil {
						listenersMap["port"] = listeners.Port
					}

					if listeners.Rules != nil {
						rulesList := []interface{}{}
						for _, rules := range listeners.Rules {
							rulesMap := map[string]interface{}{}

							if rules.LocationId != nil {
								rulesMap["location_id"] = rules.LocationId
							}

							if rules.Domain != nil {
								rulesMap["domain"] = rules.Domain
							}

							if rules.Url != nil {
								rulesMap["url"] = rules.Url
							}

							if rules.Targets != nil {
								targetsList := []interface{}{}
								for _, targets := range rules.Targets {
									targetsMap := map[string]interface{}{}

									if targets.IP != nil {
										targetsMap["ip"] = targets.IP
									}

									if targets.Port != nil {
										targetsMap["port"] = targets.Port
									}

									if targets.HealthStatus != nil {
										targetsMap["health_status"] = targets.HealthStatus
									}

									if targets.TargetId != nil {
										targetsMap["target_id"] = targets.TargetId
									}

									if targets.HealthStatusDetail != nil {
										targetsMap["health_status_detail"] = targets.HealthStatusDetail
									}

									targetsList = append(targetsList, targetsMap)
								}

								rulesMap["targets"] = targetsList
							}

							rulesList = append(rulesList, rulesMap)
						}

						listenersMap["rules"] = rulesList
					}

					listenersList = append(listenersList, listenersMap)
				}

				loadBalancerHealthMap["listeners"] = listenersList
			}

			ids = append(ids, *loadBalancerHealth.LoadBalancerId)
			tmpList = append(tmpList, loadBalancerHealthMap)
		}

		_ = d.Set("load_balancers", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
