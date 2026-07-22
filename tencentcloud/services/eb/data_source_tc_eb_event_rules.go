package eb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudEbEventRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEbEventRulesRead,
		Schema: map[string]*schema.Schema{
			"event_bus_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "事件 bus ID。",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "According 到 其中 字段 到 sort 返回 results， following 字段 是 支持: AddTime (创建时间)，ModTime (修改时间)。",
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Return results 在 ascending 或 降序，可选 值 ASC (ascending) 和 DESC (descending)。",
			},

			"rules": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Event 规则 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态",
						},
						"mod_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "修改时间。",
						},
						"enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "启用 switch。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "规则 ID。",
						},
						"add_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"event_bus_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "事件 bus ID。",
						},
						"rule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "规则 名称",
						},
						"targets": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Target brief 信息，note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "目标 ID。",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "目标 类型",
									},
								},
							},
						},
						"dead_letter_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "dlq 规则 集合 通过 规则. It 可能 是 null. 注意: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dispose_method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Support three modes 的 dlq，discarding，ignoring errors 和 continuing 到 pass，corresponding 到: DLQ，DROP，IGNORE_ERROR。",
									},
									"ckafka_delivery_params": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "After setting DLQ 模式，此 选项 为必填项. 错误信息 将 是 delivered 到 corresponding kafka 主题 注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"topic_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ckafka 主题 名称",
												},
												"resource_description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ckafka 资源 qcs six-segment。",
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

func dataSourceTencentCloudEbEventRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_eb_event_rules.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("event_bus_id"); ok {
		paramMap["EventBusId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	service := EbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var rules []*eb.Rule

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeEbEventRulesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		rules = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(rules))
	tmpList := make([]map[string]interface{}, 0, len(rules))

	if rules != nil {
		for _, rule := range rules {
			ruleMap := map[string]interface{}{}

			if rule.Status != nil {
				ruleMap["status"] = rule.Status
			}

			if rule.ModTime != nil {
				ruleMap["mod_time"] = rule.ModTime
			}

			if rule.Enable != nil {
				ruleMap["enable"] = rule.Enable
			}

			if rule.Description != nil {
				ruleMap["description"] = rule.Description
			}

			if rule.RuleId != nil {
				ruleMap["rule_id"] = rule.RuleId
			}

			if rule.AddTime != nil {
				ruleMap["add_time"] = rule.AddTime
			}

			if rule.EventBusId != nil {
				ruleMap["event_bus_id"] = rule.EventBusId
			}

			if rule.RuleName != nil {
				ruleMap["rule_name"] = rule.RuleName
			}

			if rule.Targets != nil {
				targetsList := []interface{}{}
				for _, targets := range rule.Targets {
					targetsMap := map[string]interface{}{}

					if targets.TargetId != nil {
						targetsMap["target_id"] = targets.TargetId
					}

					if targets.Type != nil {
						targetsMap["type"] = targets.Type
					}

					targetsList = append(targetsList, targetsMap)
				}

				ruleMap["targets"] = []interface{}{targetsList}
			}

			if rule.DeadLetterConfig != nil {
				deadLetterConfigMap := map[string]interface{}{}

				if rule.DeadLetterConfig.DisposeMethod != nil {
					deadLetterConfigMap["dispose_method"] = rule.DeadLetterConfig.DisposeMethod
				}

				if rule.DeadLetterConfig.CkafkaDeliveryParams != nil {
					ckafkaDeliveryParamsMap := map[string]interface{}{}

					if rule.DeadLetterConfig.CkafkaDeliveryParams.TopicName != nil {
						ckafkaDeliveryParamsMap["topic_name"] = rule.DeadLetterConfig.CkafkaDeliveryParams.TopicName
					}

					if rule.DeadLetterConfig.CkafkaDeliveryParams.ResourceDescription != nil {
						ckafkaDeliveryParamsMap["resource_description"] = rule.DeadLetterConfig.CkafkaDeliveryParams.ResourceDescription
					}

					deadLetterConfigMap["ckafka_delivery_params"] = []interface{}{ckafkaDeliveryParamsMap}
				}

				ruleMap["dead_letter_config"] = []interface{}{deadLetterConfigMap}
			}

			ids = append(ids, *rule.EventBusId)
			tmpList = append(tmpList, ruleMap)
		}

		_ = d.Set("rules", tmpList)
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
