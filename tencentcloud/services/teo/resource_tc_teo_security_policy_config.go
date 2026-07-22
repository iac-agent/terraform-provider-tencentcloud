package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoSecurityPolicyConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoSecurityPolicyConfigCreate,
		Read:   resourceTencentCloudTeoSecurityPolicyConfigRead,
		Update: resourceTencentCloudTeoSecurityPolicyConfigUpdate,
		Delete: resourceTencentCloudTeoSecurityPolicyConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "可用区 ID",
			},

			"security_policy": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Security 策略 配置. 它 是 recommended 到 使用 对于 自定义 policies 和 managed 规则 configurations 的 Web protection. 它 支持 configuring 安全 policies 使用 expression grammar。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"custom_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Custom 规则 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Deprecated:  "It has been deprecated from version 1.81.184. Please use `precise_match_rules` or `basic_access_rules` instead.",
										Description: "列表 自定义 规则 definitions. <br>当 modifying Web protection 配置 使用 ModifySecurityPolicy: <br> - 如果 Rules 参数 是 不 指定 或 参数 长度 的 Rules 是 zero: clear all 自定义 规则 configurations. <br> - 如果 参数 值 的 CustomRules 在 SecurityPolicy 参数 是 不 指定: keep existing 自定义 规则 配置 without modification。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "名称 自定义 规则。",
												},
												"condition": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "特定 内容 的 自定义 规则 必须 comply 使用 expression grammar. please refer 到 product document 对于 detailed specifications。",
												},
												"action": {
													Type:        schema.TypeList,
													Required:    true,
													MaxItems:    1,
													Description: "Execution actions 对于 自定义 规则. 名称 参数 值 的 SecurityAction 支持: <li>Deny: block;</li> <li>Monitor: observe;</li> <li>ReturnCustomPage: block 使用 指定 页面;</li> <li>Redirect: Redirect 到 URL;</li> <li>BlockIP: IP blocking;</li> <li>JSChallenge: JavaScript challenge;</li> <li>ManagedChallenge: managed challenge;</li> <li>Allow: Allow.</li>。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specific actions 对于 safe execution. 有效 值:.\n<li>Deny: block</li> <li>Monitor: Monitor</li> <li>ReturnCustomPage: 使用 指定 页面 到 block</li> <li>Redirect: Redirect 到 URL</li> <li>BlockIP: IP block</li> <li>JSChallenge: JavaScript challenge</li> <li>ManagedChallenge: managed challenge</li> <li>已禁用: 已禁用</li> <li>Allow: Allow</li>。",
															},
															"block_ip_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 BlockIP。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"duration": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Penalty 时长 对于 blocking ips. 支持 units: <li>s: second，值 范围 1-120;</li> <li>m: minute，值 范围 1-120;</li> <li>h: hour，值 范围 1-48.</li>。",
																		},
																	},
																},
															},
															"return_custom_page_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 ReturnCustomPage。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"response_code": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Response 状态 代码",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Response 自定义 页面 ID。",
																		},
																	},
																},
															},
															"redirect_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 Redirect。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Redirect URL",
																		},
																	},
																},
															},
														},
													},
												},
												"enabled": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "表示是否custom 规则 是 已启用 有效值：<li>在: 已启用</li> <li>关闭: 已禁用</li>。",
												},
												"id": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "ID 自定义 规则. <br> 规则 ID 支持 different 规则 配置 operations: <br> - add new 规则: ID 是 空 或 ID 参数 是 不 指定; <br> - modify existing 规则: 指定rule ID 该 needs 到 是 更新/modified; <br> - delete existing 规则: existing Rules 不 included 在 Rules 列表 CustomRules 参数 将 是 删除。",
												},
												"rule_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "类型 自定义 规则. 有效值：<li>BasicAccessRule: basic 访问 control;</li> <li>PreciseMatchRule: exact matching 规则，默认值;</li> <li>ManagedAccessRule: expert customized 规则，对于 output 仅.</li> 默认值为 PreciseMatchRule。",
												},
												"priority": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Customizes 优先级 的 规则. 取值范围：0-100. 它 默认为 0. 仅 支持 `rule_type` 是 `PreciseMatchRule`。",
												},
											},
										},
									},
									"precise_match_rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "列表 自定义 规则 definitions. <br>当 modifying Web protection 配置 使用 ModifySecurityPolicy: <br> - 如果 Rules 参数 是 不 指定 或 参数 长度 的 Rules 是 zero: clear all 自定义 规则 configurations. <br> - 如果 参数 值 的 CustomRules 在 SecurityPolicy 参数 是 不 指定: keep existing 自定义 规则 配置 without modification。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "名称 自定义 规则。",
												},
												"condition": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "特定 内容 的 自定义 规则 必须 comply 使用 expression grammar. please refer 到 product document 对于 detailed specifications。",
												},
												"action": {
													Type:        schema.TypeList,
													Required:    true,
													MaxItems:    1,
													Description: "Execution actions 对于 自定义 规则. 名称 参数 值 的 SecurityAction 支持: <li>Deny: block;</li> <li>Monitor: observe;</li> <li>ReturnCustomPage: block 使用 指定 页面;</li> <li>Redirect: Redirect 到 URL;</li> <li>BlockIP: IP blocking;</li> <li>JSChallenge: JavaScript challenge;</li> <li>ManagedChallenge: managed challenge;</li> <li>Allow: Allow.</li>。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specific actions 对于 safe execution. 有效 值:.\n<li>Deny: block</li> <li>Monitor: Monitor</li> <li>ReturnCustomPage: 使用 指定 页面 到 block</li> <li>Redirect: Redirect 到 URL</li> <li>BlockIP: IP block</li> <li>JSChallenge: JavaScript challenge</li> <li>ManagedChallenge: managed challenge</li> <li>已禁用: 已禁用</li> <li>Allow: Allow</li>。",
															},
															"block_ip_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 BlockIP。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"duration": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Penalty 时长 对于 blocking ips. 支持 units: <li>s: second，值 范围 1-120;</li> <li>m: minute，值 范围 1-120;</li> <li>h: hour，值 范围 1-48.</li>。",
																		},
																	},
																},
															},
															"return_custom_page_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 ReturnCustomPage。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"response_code": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Response 状态 代码",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Response 自定义 页面 ID。",
																		},
																	},
																},
															},
															"redirect_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 Redirect。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Redirect URL",
																		},
																	},
																},
															},
														},
													},
												},
												"enabled": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "表示是否custom 规则 是 已启用 有效值：<li>在: 已启用</li> <li>关闭: 已禁用</li>。",
												},
												"id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ID 自定义 规则. <br> 规则 ID 支持 different 规则 配置 operations: <br> - add new 规则: ID 是 空 或 ID 参数 是 不 指定; <br> - modify existing 规则: 指定rule ID 该 needs 到 是 更新/modified; <br> - delete existing 规则: existing Rules 不 included 在 Rules 列表 CustomRules 参数 将 是 删除。",
												},
												"rule_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "类型 自定义 规则. 有效值：<li>BasicAccessRule: basic 访问 control;</li> <li>PreciseMatchRule: exact matching 规则，默认值;</li> <li>ManagedAccessRule: expert customized 规则，对于 output 仅.</li> 默认值为 PreciseMatchRule。",
												},
												"priority": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Customizes 优先级 的 规则. 取值范围：0-100. 它 默认为 0. 仅 支持 `rule_type` 是 `PreciseMatchRule`。",
												},
											},
										},
									},
									"basic_access_rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "列表 自定义 规则 definitions. <br>当 modifying Web protection 配置 使用 ModifySecurityPolicy: <br> - 如果 Rules 参数 是 不 指定 或 参数 长度 的 Rules 是 zero: clear all 自定义 规则 configurations. <br> - 如果 参数 值 的 CustomRules 在 SecurityPolicy 参数 是 不 指定: keep existing 自定义 规则 配置 without modification。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "名称 自定义 规则。",
												},
												"condition": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "特定 内容 的 自定义 规则 必须 comply 使用 expression grammar. please refer 到 product document 对于 detailed specifications。",
												},
												"action": {
													Type:        schema.TypeList,
													Required:    true,
													MaxItems:    1,
													Description: "Execution actions 对于 自定义 规则. 名称 参数 值 的 SecurityAction 支持: <li>Deny: block;</li> <li>Monitor: observe;</li> <li>ReturnCustomPage: block 使用 指定 页面;</li> <li>Redirect: Redirect 到 URL;</li> <li>BlockIP: IP blocking;</li> <li>JSChallenge: JavaScript challenge;</li> <li>ManagedChallenge: managed challenge;</li> <li>Allow: Allow.</li>。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specific actions 对于 safe execution. 有效 值:.\n<li>Deny: block</li> <li>Monitor: Monitor</li> <li>ReturnCustomPage: 使用 指定 页面 到 block</li> <li>Redirect: Redirect 到 URL</li> <li>BlockIP: IP block</li> <li>JSChallenge: JavaScript challenge</li> <li>ManagedChallenge: managed challenge</li> <li>已禁用: 已禁用</li> <li>Allow: Allow</li>。",
															},
															"block_ip_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 BlockIP。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"duration": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Penalty 时长 对于 blocking ips. 支持 units: <li>s: second，值 范围 1-120;</li> <li>m: minute，值 范围 1-120;</li> <li>h: hour，值 范围 1-48.</li>。",
																		},
																	},
																},
															},
															"return_custom_page_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 ReturnCustomPage。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"response_code": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Response 状态 代码",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Response 自定义 页面 ID。",
																		},
																	},
																},
															},
															"redirect_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 Redirect。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Redirect URL",
																		},
																	},
																},
															},
														},
													},
												},
												"enabled": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "表示是否custom 规则 是 已启用 有效值：<li>在: 已启用</li> <li>关闭: 已禁用</li>。",
												},
												"id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ID 自定义 规则. <br> 规则 ID 支持 different 规则 配置 operations: <br> - add new 规则: ID 是 空 或 ID 参数 是 不 指定; <br> - modify existing 规则: 指定rule ID 该 needs 到 是 更新/modified; <br> - delete existing 规则: existing Rules 不 included 在 Rules 列表 CustomRules 参数 将 是 删除。",
												},
												"rule_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "类型 自定义 规则. 有效值：<li>BasicAccessRule: basic 访问 control;</li> <li>PreciseMatchRule: exact matching 规则，默认值;</li> <li>ManagedAccessRule: expert customized 规则，对于 output 仅.</li> 默认值为 PreciseMatchRule。",
												},
												"priority": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Customizes 优先级 的 规则. 取值范围：0-100. 它 默认为 0. 仅 支持 `rule_type` 是 `PreciseMatchRule`。",
												},
											},
										},
									},
								},
							},
						},
						"managed_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "Managed 规则 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "表示是否managed 规则 是 已启用 有效值：<li>在: 已启用 all managed 规则 take effect 作为 已配置;</li> <li>关闭: 已禁用 all managed 规则 do 不 take effect.</li>。",
									},
									"detection_only": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "表示是否evaluation 模式 是 已启用 它 是 有效 仅 当 已启用 参数 是 集合 到 在. 有效值：<li>在: 已启用 all managed 规则 take effect 在 observation 模式</li> <li>关闭: 已禁用 all managed 规则 take effect according 到 actual 配置.</li>。",
									},
									"semantic_analysis": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "是否managed 规则 semantic analysis 选项 是 已启用 是 有效 仅 当 已启用 参数 是 在. 有效值：<li>在: 启用. perform semantic analysis 在 requests before processing them;</li> <li>关闭: disable. process requests directly without semantic analysis.</li> <br/>默认值 关闭。",
									},
									"auto_update": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Managed 规则 automatic update 选项。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"auto_update_to_latest_version": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "表示是否enable automatic update 到 latest 版本 有效值：<li>在: 已启用</li> <li>关闭: 已禁用</li>。",
												},
												"ruleset_version": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "currently 使用 版本，在 格式 compliant 使用 ISO 8601 standard，such 作为 2023-12-21T12:00:32Z. 它 是 空 通过 默认值 和 是 仅 output 参数。",
												},
											},
										},
									},
									"managed_rule_groups": {
										Type:        schema.TypeSet,
										Optional:    true,
										Computed:    true,
										Description: "Configuration 的 managed 规则 组. 如果 此 structure 是 passed 作为 空 数组 或 GroupId 是 不 included 在 列表，它 将 是 processed based 在 默认值 方法。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"group_id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "组名称 的 managed 规则. 如果 规则 组 对于 配置 是 不 指定，它 将 是 processed based 在 默认值 配置. refer 到 product documentation 对于 特定 值 的 GroupId。",
												},
												"sensitivity_level": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Protection 级别 的 managed 规则 组. 有效值：<li>loose: lenient，仅 包含ultra-high risk 规则. 在 此 point，configure 操作，和 RuleActions 配置 是 无效;</li> <li>normal: normal，包含ultra-high risk 和 high-risk 规则. 在 此 point，configure 操作，和 RuleActions 配置 是 无效;</li> <li>strict: strict，包含ultra-high risk，high-risk 和 medium-risk 规则. 在 此 point，configure 操作，和 RuleActions 配置 是 无效;</li> <li>extreme: super strict，包含ultra-high risk，high-risk，medium-risk 和 low-risk 规则. 在 此 point，configure 操作，和 RuleActions 配置 是 无效;</li> <li>自定义: 自定义，refined strategy. configure disposal 方法 对于 each individual 规则. 在 此 point， 操作 字段 是 无效. 使用 RuleActions 到 configure refined strategy 对于 each individual 规则.</li>。",
												},
												"action": {
													Type:        schema.TypeList,
													Required:    true,
													MaxItems:    1,
													Description: "Handling actions 对于 managed 规则 groups. 名称 参数 值 的 SecurityAction 支持: <li>Deny: block 和 respond 使用 interception 页面;</li> <li>Monitor: observe，do 不 process requests 和 记录 安全 events 在 logs;</li> <li>已禁用: 不 已启用，do 不 scan requests 和 skip 此 规则.</li>。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specific actions 对于 safe execution. 有效 值:.\n<li>Deny: block</li> <li>Monitor: Monitor</li> <li>ReturnCustomPage: 使用 指定 页面 到 block</li> <li>Redirect: Redirect 到 URL</li> <li>BlockIP: IP block</li> <li>JSChallenge: JavaScript challenge</li> <li>ManagedChallenge: managed challenge</li> <li>已禁用: 已禁用</li> <li>Allow: Allow</li>。",
															},
															"block_ip_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 BlockIP。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"duration": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Penalty 时长 对于 blocking ips. 支持 units: <li>s: second，值 范围 1-120;</li> <li>m: minute，值 范围 1-120;</li> <li>h: hour，值 范围 1-48.</li>。",
																		},
																	},
																},
															},
															"return_custom_page_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 ReturnCustomPage。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"response_code": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Response 状态 代码",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Response 自定义 页面 ID。",
																		},
																	},
																},
															},
															"redirect_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 Redirect。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Redirect URL",
																		},
																	},
																},
															},
														},
													},
												},
												"rule_actions": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Specific 配置 的 规则 items under managed 规则 组. 配置 是 effective 仅 当 SensitivityLevel 是 自定义。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"rule_id": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specific items under managed 规则 组，其中 是 用于rewrite 配置 内容 的 此 individual 规则 item. refer 到 product documentation 对于 details。",
															},
															"action": {
																Type:        schema.TypeList,
																Required:    true,
																MaxItems:    1,
																Description: "指定handling 操作 对于 managed 规则 item 在 RuleId. 名称 参数 值 的 SecurityAction 支持: <li>Deny: block 和 respond 使用 interception 页面;</li> <li>Monitor: observe，do 不 process 请求 和 记录 安全 事件 在 logs;</li> <li>已禁用: 已禁用，do 不 scan 请求 和 skip 此 规则.</li>。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specific actions 对于 safe execution. 有效 值:.\n<li>Deny: block</li> <li>Monitor: Monitor</li> <li>ReturnCustomPage: 使用 指定 页面 到 block</li> <li>Redirect: Redirect 到 URL</li> <li>BlockIP: IP block</li> <li>JSChallenge: JavaScript challenge</li> <li>ManagedChallenge: managed challenge</li> <li>已禁用: 已禁用</li> <li>Allow: Allow</li>。",
																		},
																		"block_ip_action_parameters": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "Additional 参数 当 名称 是 BlockIP。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"duration": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Penalty 时长 对于 blocking ips. 支持 units: <li>s: second，值 范围 1-120;</li> <li>m: minute，值 范围 1-120;</li> <li>h: hour，值 范围 1-48.</li>。",
																					},
																				},
																			},
																		},
																		"return_custom_page_action_parameters": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "Additional 参数 当 名称 是 ReturnCustomPage。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"response_code": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Response 状态 代码",
																					},
																					"error_page_id": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Response 自定义 页面 ID。",
																					},
																				},
																			},
																		},
																		"redirect_action_parameters": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "Additional 参数 当 名称 是 Redirect。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"url": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Redirect URL",
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
												"meta_data": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Managed 规则 组 信息，对于 output 仅。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"group_detail": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Managed 规则 组 描述，对于 output 仅。",
															},
															"group_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Managed 规则 组名称，对于 output 仅。",
															},
															"rule_details": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "All sub-规则 信息 under 当前 managed 规则 组，对于 output 仅。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"rule_id": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Managed 规则 ID。",
																		},
																		"risk_level": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Protection 级别 的 managed 规则. 有效值：<li>low: low risk. 此 规则 has relatively low risk 和 是 applicable 到 访问 scenarios 在 very strict control 环境. 此 级别 的 规则 可能 generate considerable false alarms.</li> <li>medium: medium risk. 此 表示 risk 的 此 规则 是 normal 和 它 是 suitable 对于 protection scenarios 使用 stricter requirements.</li> <li>high: high risk. 此 表示that risk 的 此 规则 是 relatively high 和 它 将 不 generate false alarms 在 most scenarios.</li> <li>extreme: ultra-high risk. 此 表示 该 risk 的 此 规则 是 extremely high 和 它 将 不 generate false alarms basically.</li>。",
																		},
																		"description": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Rule 描述",
																		},
																		"tags": {
																			Type:        schema.TypeSet,
																			Computed:    true,
																			Description: "Rule 标签 some types 的 规则 do 不 have 标签",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"rule_version": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Rule ownership 版本",
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
						"http_ddos_protection": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "HTTP DDOS protection 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"adaptive_frequency_control": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Specific 配置 的 adaptive 频率 control。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Whether adaptive 频率 control 是 已启用 possible 值 是: <li>在: 已启用; </li><li>关闭: 已禁用 </li>。",
												},
												"sensitivity": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "restriction 级别 的 adaptive 频率 control. 当 已启用 是 在，此 字段 为必填项. 值 是: <li>Loose: loose; </li><li>Moderate: moderate; </li><li>Strict: strict. </li>。",
												},
												"action": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "handling 方法 的 adaptive 频率 control. 当 已启用 是 在，此 字段 为必填项. SecurityAction's 名称 值 支持: <li>Monitor: Observe; </li><li>Deny: Intercept; </li><li>Challenge: Challenge，其中 ChallengeActionParameters.名称 仅 支持 JSChallenge. </li>。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "特定 操作 的 安全 execution. 值 是:\n<li>Deny: intercept，block 请求 到 访问 site resources;</li>\n<li>Monitor: observe，仅 记录 logs;</li>\n<li>Redirect: redirect 到 URL;</li>\n<li>已禁用: 已禁用，do 不 启用 指定 规则;</li>\n<li>Allow: allow 访问，但 延迟 processing requests;</li>\n<li>Challenge: challenge，respond 到 challenge 内容;</li>\n<li>BlockIP: 到 是 abandoned，IP ban;</li>\n<li>ReturnCustomPage: 到 是 abandoned，使用 指定 页面 到 intercept;</li>\n<li>JSChallenge: 到 是 abandoned，JavaScript challenge;</li>\n<li>ManagedChallenge: 到 是 abandoned，managed challenge.</li>。",
															},
															"deny_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 Deny。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"block_ip": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "是否extend blocking 的 来源 IP. possible 值 是:\n<li>在: 在;</li>\n<li>关闭: 关闭.</li>\nWhen 已启用， 客户端 IP 该 triggers 规则 将 是 blocked continuously. 当 此 选项 是 已启用， BlockIpDuration 参数 必须 是 指定 在 same 时间.\nNote: 此 选项 不能 是 已启用 在 same 时间 作为 ReturnCustomPage 或 Stall options。",
																		},
																		"block_ip_duration": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "当 BlockIP 是 在， IP blocking 时长。",
																		},
																		"return_custom_page": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "是否use 自定义 pages. possible 值 是:\n<li>在: 在;</li>\n<li>关闭: 关闭.</li>\nAfter enabling，使用 自定义 页面 内容 到 intercept (respond 到) requests. 当 enabling 此 选项，您 必须 指定ResponseCode 和 ErrorPageId 参数 在 same 时间.\nNote: 此 选项 不能 是 已启用 在 same 时间 作为 BlockIp 或 Stall options。",
																		},
																		"response_code": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Customize 状态 代码 的 页面。",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "PageId 的 自定义 页面。",
																		},
																		"stall": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "是否ignore 请求来源 suspension. 值 是:\n<li>在: Enable;</li>\n<li>关闭: Disable.</li>\nAfter enabling，它 将 无 longer respond 到 requests 在 当前 连接 会话 和 将 不 actively disconnect. It 是 用于fight against crawlers 和 consume 客户端 连接 resources.\nNote: 此 选项 不能 是 已启用 在 same 时间 作为 BlockIp 或 ReturnCustomPage options。",
																		},
																	},
																},
															},
															"redirect_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 Redirect。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "URL 到 redirect。",
																		},
																	},
																},
															},
															"challenge_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 Challenge。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"challenge_option": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "特定 challenge 操作 到 是 executed safely. possible 值 是: <li> InterstitialChallenge: interstitial challenge; </li><li> InlineChallenge: embedded challenge; </li><li> JSChallenge: JavaScript challenge; </li><li> ManagedChallenge: managed challenge. </li>。",
																		},
																		"interval": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "时间间隔 对于 repeating challenge. 当 名称 是 InterstitialChallenge/InlineChallenge，此 字段 为必填项. 默认值为 300s. Supported units 是: <li>s: 秒，值 范围 1 到 60; </li><li>m: minutes，值 范围 1 到 60; </li><li>h: hours，值 范围 1 到 24. </li>。",
																		},
																		"attester_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Client authentication 方法 ID. 此 字段 为必填项 当 名称 是 InterstitialChallenge/InlineChallenge。",
																		},
																	},
																},
															},
															"block_ip_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "待废弃， additional 参数 当 名称 是 BlockIP。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"duration": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "penalty 时长 对于 banning IP. Supported units 是: <li>s: 秒，值 范围 1 到 120; </li><li>m: minutes，值 范围 1 到 120; </li><li>h: hours，值 范围 1 到 48. </li>。",
																		},
																	},
																},
															},
															"return_custom_page_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "待废弃， additional 参数 当 名称 是 ReturnCustomPage。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"response_code": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Response 状态 代码",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "自定义 页面 ID response。",
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
									"client_filtering": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Specific 配置 的 intelligent 客户端 filtering。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Whether smart 客户端 filtering 是 已启用 possible 值 是: <li>在: 已启用; </li><li>关闭: 已禁用 </li>。",
												},
												"action": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "方法 的 intelligent 客户端 filtering. 当 已启用 是 在，此 字段 为必填项. SecurityAction 名称 值 支持: <li>Monitor: Observe; </li><li>Deny: Intercept; </li><li>Challenge: Challenge，其中 ChallengeActionParameters.名称 仅 支持 JSChallenge. </li>。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "特定 操作 的 安全 execution. 值 是:\n<li>Deny: intercept，block 请求 到 访问 site resources;</li>\n<li>Monitor: observe，仅 记录 logs;</li>\n<li>Redirect: redirect 到 URL;</li>\n<li>已禁用: 已禁用，do 不 启用 指定 规则;</li>\n<li>Allow: allow 访问，但 延迟 processing requests;</li>\n<li>Challenge: challenge，respond 到 challenge 内容;</li>\n<li>BlockIP: 到 是 abandoned，IP ban;</li>\n<li>ReturnCustomPage: 到 是 abandoned，使用 指定 页面 到 intercept;</li>\n<li>JSChallenge: 到 是 abandoned，JavaScript challenge;</li>\n<li>ManagedChallenge: 到 是 abandoned，managed challenge.</li>。",
															},
															"deny_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 Deny。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"block_ip": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "是否extend blocking 的 来源 IP. possible 值 是:\n<li>在: 在;</li>\n<li>关闭: 关闭.</li>\nWhen 已启用， 客户端 IP 该 triggers 规则 将 是 blocked continuously. 当 此 选项 是 已启用， BlockIpDuration 参数 必须 是 指定 在 same 时间.\nNote: 此 选项 不能 是 已启用 在 same 时间 作为 ReturnCustomPage 或 Stall options。",
																		},
																		"block_ip_duration": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "当 BlockIP 是 在， IP blocking 时长。",
																		},
																		"return_custom_page": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "是否use 自定义 pages. possible 值 是:\n<li>在: 在;</li>\n<li>关闭: 关闭.</li>\nAfter enabling，使用 自定义 页面 内容 到 intercept (respond 到) requests. 当 enabling 此 选项，您 必须 指定ResponseCode 和 ErrorPageId 参数 在 same 时间.\nNote: 此 选项 不能 是 已启用 在 same 时间 作为 BlockIp 或 Stall options。",
																		},
																		"response_code": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Customize 状态 代码 的 页面。",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "PageId 的 自定义 页面。",
																		},
																		"stall": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "是否ignore 请求来源 suspension. 值 是:\n<li>在: Enable;</li>\n<li>关闭: Disable.</li>\nAfter enabling，它 将 无 longer respond 到 requests 在 当前 连接 会话 和 将 不 actively disconnect. It 是 用于fight against crawlers 和 consume 客户端 连接 resources.\nNote: 此 选项 不能 是 已启用 在 same 时间 作为 BlockIp 或 ReturnCustomPage options。",
																		},
																	},
																},
															},
															"redirect_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 Redirect。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "URL 到 redirect。",
																		},
																	},
																},
															},
															"challenge_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 Challenge。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"challenge_option": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "特定 challenge 操作 到 是 executed safely. possible 值 是: <li> InterstitialChallenge: interstitial challenge; </li><li> InlineChallenge: embedded challenge; </li><li> JSChallenge: JavaScript challenge; </li><li> ManagedChallenge: managed challenge. </li>。",
																		},
																		"interval": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "时间间隔 对于 repeating challenge. 当 名称 是 InterstitialChallenge/InlineChallenge，此 字段 为必填项. 默认值为 300s. Supported units 是: <li>s: 秒，值 范围 1 到 60; </li><li>m: minutes，值 范围 1 到 60; </li><li>h: hours，值 范围 1 到 24. </li>。",
																		},
																		"attester_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Client authentication 方法 ID. 此 字段 为必填项 当 名称 是 InterstitialChallenge/InlineChallenge。",
																		},
																	},
																},
															},
															"block_ip_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "待废弃， additional 参数 当 名称 是 BlockIP。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"duration": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "penalty 时长 对于 banning IP. Supported units 是: <li>s: 秒，值 范围 1 到 120; </li><li>m: minutes，值 范围 1 到 120; </li><li>h: hours，值 范围 1 到 48. </li>。",
																		},
																	},
																},
															},
															"return_custom_page_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "待废弃， additional 参数 当 名称 是 ReturnCustomPage。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"response_code": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Response 状态 代码",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "自定义 页面 ID response。",
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
									"bandwidth_abuse_defense": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Specific 配置 的 流量 fraud prevention。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "是否anti-theft 功能 (仅 applicable 到 mainland China) 是 已启用 possible 值 是: <li>在: 已启用; </li><li>关闭: 已禁用 </li>。",
												},
												"action": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "方法 对于 preventing 流量 fraud (仅 applicable 到 mainland China). 当 已启用 是 在，此 字段 为必填项. SecurityAction 名称 值 支持: <li>Monitor: Observe; </li><li>Deny: Intercept; </li><li>Challenge: Challenge，其中 ChallengeActionParameters.名称 仅 支持 JSChallenge. </li>。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "特定 操作 的 安全 execution. 值 是:\n<li>Deny: intercept，block 请求 到 访问 site resources;</li>\n<li>Monitor: observe，仅 记录 logs;</li>\n<li>Redirect: redirect 到 URL;</li>\n<li>已禁用: 已禁用，do 不 启用 指定 规则;</li>\n<li>Allow: allow 访问，但 延迟 processing requests;</li>\n<li>Challenge: challenge，respond 到 challenge 内容;</li>\n<li>BlockIP: 到 是 abandoned，IP ban;</li>\n<li>ReturnCustomPage: 到 是 abandoned，使用 指定 页面 到 intercept;</li>\n<li>JSChallenge: 到 是 abandoned，JavaScript challenge;</li>\n<li>ManagedChallenge: 到 是 abandoned，managed challenge.</li>。",
															},
															"deny_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 Deny。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"block_ip": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "是否extend blocking 的 来源 IP. possible 值 是:\n<li>在: 在;</li>\n<li>关闭: 关闭.</li>\nWhen 已启用， 客户端 IP 该 triggers 规则 将 是 blocked continuously. 当 此 选项 是 已启用， BlockIpDuration 参数 必须 是 指定 在 same 时间.\nNote: 此 选项 不能 是 已启用 在 same 时间 作为 ReturnCustomPage 或 Stall options。",
																		},
																		"block_ip_duration": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "当 BlockIP 是 在， IP blocking 时长。",
																		},
																		"return_custom_page": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "是否use 自定义 pages. possible 值 是:\n<li>在: 在;</li>\n<li>关闭: 关闭.</li>\nAfter enabling，使用 自定义 页面 内容 到 intercept (respond 到) requests. 当 enabling 此 选项，您 必须 指定ResponseCode 和 ErrorPageId 参数 在 same 时间.\nNote: 此 选项 不能 是 已启用 在 same 时间 作为 BlockIp 或 Stall options。",
																		},
																		"response_code": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Customize 状态 代码 的 页面。",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "PageId 的 自定义 页面。",
																		},
																		"stall": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "是否ignore 请求来源 suspension. 值 是:\n<li>在: Enable;</li>\n<li>关闭: Disable.</li>\nAfter enabling，它 将 无 longer respond 到 requests 在 当前 连接 会话 和 将 不 actively disconnect. It 是 用于fight against crawlers 和 consume 客户端 连接 resources.\nNote: 此 选项 不能 是 已启用 在 same 时间 作为 BlockIp 或 ReturnCustomPage options。",
																		},
																	},
																},
															},
															"redirect_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 Redirect。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "URL 到 redirect。",
																		},
																	},
																},
															},
															"challenge_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 Challenge。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"challenge_option": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "特定 challenge 操作 到 是 executed safely. possible 值 是: <li> InterstitialChallenge: interstitial challenge; </li><li> InlineChallenge: embedded challenge; </li><li> JSChallenge: JavaScript challenge; </li><li> ManagedChallenge: managed challenge. </li>。",
																		},
																		"interval": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "时间间隔 对于 repeating challenge. 当 名称 是 InterstitialChallenge/InlineChallenge，此 字段 为必填项. 默认值为 300s. Supported units 是: <li>s: 秒，值 范围 1 到 60; </li><li>m: minutes，值 范围 1 到 60; </li><li>h: hours，值 范围 1 到 24. </li>。",
																		},
																		"attester_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Client authentication 方法 ID. 此 字段 为必填项 当 名称 是 InterstitialChallenge/InlineChallenge。",
																		},
																	},
																},
															},
															"block_ip_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "待废弃， additional 参数 当 名称 是 BlockIP。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"duration": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "penalty 时长 对于 banning IP. Supported units 是: <li>s: 秒，值 范围 1 到 120; </li><li>m: minutes，值 范围 1 到 120; </li><li>h: hours，值 范围 1 到 48. </li>。",
																		},
																	},
																},
															},
															"return_custom_page_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "待废弃， additional 参数 当 名称 是 ReturnCustomPage。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"response_code": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Response 状态 代码",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "自定义 页面 ID response。",
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
									"slow_attack_defense": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Specific 配置 的 slow attack protection。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Whether slow attack protection 是 已启用 possible 值 是: <li>在: 已启用; </li><li>关闭: 已禁用 </li>。",
												},
												"action": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "handling 方法 的 slow attack protection. 当 已启用 是 在，此 字段 为必填项. SecurityAction 名称 值 支持: <li>Monitor: Observe; </li><li>Deny: Intercept; </li>。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "特定 操作 的 安全 execution. 值 是:\n<li>Deny: intercept，block 请求 到 访问 site resources;</li>\n<li>Monitor: observe，仅 记录 logs;</li>\n<li>Redirect: redirect 到 URL;</li>\n<li>已禁用: 已禁用，do 不 启用 指定 规则;</li>\n<li>Allow: allow 访问，但 延迟 processing requests;</li>\n<li>Challenge: challenge，respond 到 challenge 内容;</li>\n<li>BlockIP: 到 是 abandoned，IP ban;</li>\n<li>ReturnCustomPage: 到 是 abandoned，使用 指定 页面 到 intercept;</li>\n<li>JSChallenge: 到 是 abandoned，JavaScript challenge;</li>\n<li>ManagedChallenge: 到 是 abandoned，managed challenge.</li>。",
															},
															"deny_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 Deny。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"block_ip": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "是否extend blocking 的 来源 IP. possible 值 是:\n<li>在: 在;</li>\n<li>关闭: 关闭.</li>\nWhen 已启用， 客户端 IP 该 triggers 规则 将 是 blocked continuously. 当 此 选项 是 已启用， BlockIpDuration 参数 必须 是 指定 在 same 时间.\nNote: 此 选项 不能 是 已启用 在 same 时间 作为 ReturnCustomPage 或 Stall options。",
																		},
																		"block_ip_duration": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "当 BlockIP 是 在， IP blocking 时长。",
																		},
																		"return_custom_page": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "是否use 自定义 pages. possible 值 是:\n<li>在: 在;</li>\n<li>关闭: 关闭.</li>\nAfter enabling，使用 自定义 页面 内容 到 intercept (respond 到) requests. 当 enabling 此 选项，您 必须 指定ResponseCode 和 ErrorPageId 参数 在 same 时间.\nNote: 此 选项 不能 是 已启用 在 same 时间 作为 BlockIp 或 Stall options。",
																		},
																		"response_code": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Customize 状态 代码 的 页面。",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "PageId 的 自定义 页面。",
																		},
																		"stall": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "是否ignore 请求来源 suspension. 值 是:\n<li>在: Enable;</li>\n<li>关闭: Disable.</li>\nAfter enabling，它 将 无 longer respond 到 requests 在 当前 连接 会话 和 将 不 actively disconnect. It 是 用于fight against crawlers 和 consume 客户端 连接 resources.\nNote: 此 选项 不能 是 已启用 在 same 时间 作为 BlockIp 或 ReturnCustomPage options。",
																		},
																	},
																},
															},
															"redirect_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 Redirect。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "URL 到 redirect。",
																		},
																	},
																},
															},
															"challenge_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 Challenge。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"challenge_option": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "特定 challenge 操作 到 是 executed safely. possible 值 是: <li> InterstitialChallenge: interstitial challenge; </li><li> InlineChallenge: embedded challenge; </li><li> JSChallenge: JavaScript challenge; </li><li> ManagedChallenge: managed challenge. </li>。",
																		},
																		"interval": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "时间间隔 对于 repeating challenge. 当 名称 是 InterstitialChallenge/InlineChallenge，此 字段 为必填项. 默认值为 300s. Supported units 是: <li>s: 秒，值 范围 1 到 60; </li><li>m: minutes，值 范围 1 到 60; </li><li>h: hours，值 范围 1 到 24. </li>。",
																		},
																		"attester_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Client authentication 方法 ID. 此 字段 为必填项 当 名称 是 InterstitialChallenge/InlineChallenge。",
																		},
																	},
																},
															},
															"block_ip_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "待废弃， additional 参数 当 名称 是 BlockIP。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"duration": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "penalty 时长 对于 banning IP. Supported units 是: <li>s: 秒，值 范围 1 到 120; </li><li>m: minutes，值 范围 1 到 120; </li><li>h: hours，值 范围 1 到 48. </li>。",
																		},
																	},
																},
															},
															"return_custom_page_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "待废弃， additional 参数 当 名称 是 ReturnCustomPage。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"response_code": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Response 状态 代码",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "自定义 页面 ID response。",
																		},
																	},
																},
															},
														},
													},
												},
												"minimal_request_body_transfer_rate": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Specific 配置 的 最小 速率 阈值 对于 text transmission. 此 字段 为必填项 当 已启用 是 在。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"minimal_avg_transfer_rate_threshold": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Minimum text transmission 速率 阈值. 单位 仅 支持 bps。",
															},
															"counting_period": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "最小 text transmission 速率 统计 时间 范围， possible 值 是: <li>10s: 10 秒; </li><li>30s: 30 秒; </li><li>60s: 60 秒; </li><li>120s: 120 秒. </li>。",
															},
															"enabled": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "是否text transmission 最小 速率 阈值 是 已启用 possible 值 是: <li>在: 已启用; </li><li>关闭: 已禁用 </li>。",
															},
														},
													},
												},
												"request_body_transfer_timeout": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Specific 配置 的 text transmission 超时. 当 已启用 是 在，此 字段 为必填项。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"idle_timeout": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "text transmission 超时 周期 是 between 5 和 120，和 单位 仅 支持 秒 (s)。",
															},
															"enabled": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "是否text transmission 超时 是 已启用 possible 值 是: <li>在: 已启用; </li><li>关闭: 已禁用 </li>。",
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
						"rate_limiting_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Rate limiting 规则 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "A 列表 precise 速率 limiting definitions. 当 使用 ModifySecurityPolicy 到 modify Web protection 配置: <br> <li> 如果 Rules 参数 是 不 指定，或 Rules 参数 长度 是 zero: clear all precise 速率 limiting configurations. </li>. <li> 如果 RateLimitingRules 参数 值 是 不 指定 在 SecurityPolicy 参数: keep existing 自定义 规则 配置 和 do 不 modify 它. </li>。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "ID precise 速率 限制 <br> 规则 ID 可以 support different 规则 配置 operations: <br> <li> <b>Add</b> new 规则: ID 是 空 或 ID 参数 是 不 指定; </li><li> <b>Modify</b> existing 规则: 指定rule ID 到 是 更新/modified; </li><li> <b>Delete</b> existing 规则: 在 RateLimitingRules 参数， existing 规则 不 included 在 Rules 列表 将 是 删除. </li>。",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "名称 precise 速率 限制",
												},
												"condition": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "特定 内容 的 precise 速率 限制 必须 conform 到 expression syntax. For detailed specifications，see product documentation。",
												},
												"count_by": {
													Type:        schema.TypeSet,
													Optional:    true,
													Description: "matching 方法 的 速率 阈值 请求 功能. 当 已启用 是 在，此 字段 为必填项. <br /><br />当 there 是 多个 conditions，多个 conditions 将 是 combined 对于 statistical calculation. 数量 conditions 不能 exceed 5. possible 值 是: <br/><li><b>http.请求.ip</b>: 客户端 IP; </li><li><b>http.请求.xff_header_ip</b>: 客户端 IP (matching XFF 头部 first); </li><li><b>http.请求.uri.路径</b>: requested 访问 路径; </li><li><b>http.请求.cookies['会话']</b>: cookie named 会话，其中 会话 可以 是 replaced 通过 参数 您 specify; </li><li><b>http.请求.headers['用户-agent']</b>: HTTP 头部 named 用户-agent，其中 用户-agent 可以 是 replaced 通过 参数 您 specify; </li><li><b>http.请求.ja3</b>: requested JA3 fingerprint; </li><li><b>http.请求.uri.查询['测试']</b>: URL 查询 参数 named 测试，其中 测试 可以 是 replaced 通过 参数 您 specify. </li>。",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"max_request_threshold": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "cumulative 数量 interceptions within 时间 范围 的 precise 速率 限制，ranging 从 1 到 100000。",
												},
												"counting_period": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "statistical 时间 window， possible 值 是: <li>1s: 1 second; </li><li>5s: 5 秒; </li><li>10s: 10 秒; </li><li>20s: 20 秒; </li><li>30s: 30 秒; </li><li>40s: 40 秒; </li><li>50s: 50 秒; </li><li>1m: 1 minute; </li><li>2m: 2 minutes; </li><li>5m: 5 minutes; </li><li>10m: 10 minutes; </li><li>1h: 1 hour. </li>。",
												},
												"action_duration": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "操作 时长 的 操作 支持 units 是: <li>s: 秒，使用 值 的 1 到 120; </li><li>m: minutes，使用 值 的 1 到 120; </li><li>h: hours，使用 值 的 1 到 48; </li><li>d: days，使用 值 的 1 到 30. </li>。",
												},
												"action": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "precise 速率 限制 handling 方法. 值 是: <li>Monitor: Observe; </li><li>Deny: Intercept，其中 DenyActionParameters.名称 支持 Deny 和 ReturnCustomPage; </li><li>Challenge: Challenge，其中 ChallengeActionParameters.名称 支持 JSChallenge 和 ManagedChallenge; </li><li>Redirect: Redirect 到 URL; </li>。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "特定 操作 的 安全 execution. 值 是:\n<li>Deny: intercept，block 请求 到 访问 site resources;</li>\n<li>Monitor: observe，仅 记录 logs;</li>\n<li>Redirect: redirect 到 URL;</li>\n<li>已禁用: 已禁用，do 不 启用 指定 规则;</li>\n<li>Allow: allow 访问，但 延迟 processing requests;</li>\n<li>Challenge: challenge，respond 到 challenge 内容;</li>\n<li>BlockIP: 到 是 abandoned，IP ban;</li>\n<li>ReturnCustomPage: 到 是 abandoned，使用 指定 页面 到 intercept;</li>\n<li>JSChallenge: 到 是 abandoned，JavaScript challenge;</li>\n<li>ManagedChallenge: 到 是 abandoned，managed challenge.</li>。",
															},
															"deny_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 Deny。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"block_ip": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "是否extend blocking 的 来源 IP. possible 值 是:\n<li>在: 在;</li>\n<li>关闭: 关闭.</li>\nWhen 已启用， 客户端 IP 该 triggers 规则 将 是 blocked continuously. 当 此 选项 是 已启用， BlockIpDuration 参数 必须 是 指定 在 same 时间.\nNote: 此 选项 不能 是 已启用 在 same 时间 作为 ReturnCustomPage 或 Stall options。",
																		},
																		"block_ip_duration": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "当 BlockIP 是 在， IP blocking 时长。",
																		},
																		"return_custom_page": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "是否use 自定义 pages. possible 值 是:\n<li>在: 在;</li>\n<li>关闭: 关闭.</li>\nAfter enabling，使用 自定义 页面 内容 到 intercept (respond 到) requests. 当 enabling 此 选项，您 必须 指定ResponseCode 和 ErrorPageId 参数 在 same 时间.\nNote: 此 选项 不能 是 已启用 在 same 时间 作为 BlockIp 或 Stall options。",
																		},
																		"response_code": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Customize 状态 代码 的 页面。",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "PageId 的 自定义 页面。",
																		},
																		"stall": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "是否ignore 请求来源 suspension. 值 是:\n<li>在: Enable;</li>\n<li>关闭: Disable.</li>\nAfter enabling，它 将 无 longer respond 到 requests 在 当前 连接 会话 和 将 不 actively disconnect. It 是 用于fight against crawlers 和 consume 客户端 连接 resources.\nNote: 此 选项 不能 是 已启用 在 same 时间 作为 BlockIp 或 ReturnCustomPage options。",
																		},
																	},
																},
															},
															"redirect_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 Redirect。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "URL 到 redirect。",
																		},
																	},
																},
															},
															"challenge_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Additional 参数 当 名称 是 Challenge。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"challenge_option": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "特定 challenge 操作 到 是 executed safely. possible 值 是: <li> InterstitialChallenge: interstitial challenge; </li><li> InlineChallenge: embedded challenge; </li><li> JSChallenge: JavaScript challenge; </li><li> ManagedChallenge: managed challenge. </li>。",
																		},
																		"interval": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "时间间隔 对于 repeating challenge. 当 名称 是 InterstitialChallenge/InlineChallenge，此 字段 为必填项. 默认值为 300s. Supported units 是: <li>s: 秒，值 范围 1 到 60; </li><li>m: minutes，值 范围 1 到 60; </li><li>h: hours，值 范围 1 到 24. </li>。",
																		},
																		"attester_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Client authentication 方法 ID. 此 字段 为必填项 当 名称 是 InterstitialChallenge/InlineChallenge。",
																		},
																	},
																},
															},
															"block_ip_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "待废弃， additional 参数 当 名称 是 BlockIP。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"duration": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "penalty 时长 对于 banning IP. Supported units 是: <li>s: 秒，值 范围 1 到 120; </li><li>m: minutes，值 范围 1 到 120; </li><li>h: hours，值 范围 1 到 48. </li>。",
																		},
																	},
																},
															},
															"return_custom_page_action_parameters": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "待废弃， additional 参数 当 名称 是 ReturnCustomPage。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"response_code": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Response 状态 代码",
																		},
																		"error_page_id": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "自定义 页面 ID response。",
																		},
																	},
																},
															},
														},
													},
												},
												"priority": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "优先级 的 precise 速率 limiting ranges 从 0 到 100，和 默认为 0。",
												},
												"enabled": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "是否precise 速率 限制 规则 是 已启用 possible 值 是: <li>在: 已启用; </li><li>关闭: 已禁用 </li>。",
												},
											},
										},
									},
								},
							},
						},
						"exception_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Exception 规则 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Definition 列表 exception 规则. 当 使用 ModifySecurityPolicy 到 modify Web protection 配置: <li>如果 Rules 参数 是 不 指定，或 长度 的 Rules 参数 是 zero: clear all exception 规则 configurations. </li>.<li>如果 ExceptionRules 参数 值 是 不 指定 在 SecurityPolicy 参数: keep existing exception 规则 configurations 和 do 不 modify them. </li>。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "ID exception 规则. <br> 规则 ID 可以 support different 规则 配置 operations: <br> <li> <b>Add</b> new 规则: ID 是 空 或 ID 参数 是 不 指定; </li><li> <b>Modify</b> existing 规则: 指定rule ID 到 是 更新/modified; </li><li> <b>Delete</b> existing 规则: 在 ExceptionRules 参数， existing 规则 不 included 在 Rules 列表 将 是 删除. </li>。",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "名称 exception 规则。",
												},
												"condition": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "特定 内容 的 exception 规则 必须 comply 使用 expression syntax. For detailed specifications，see product documentation。",
												},
												"skip_scope": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Exception 规则 execution options， 值 是: <li>WebSecurityModules: 指定security protection 模块 对于 exception 规则. </li>.<li>ManagedRules: 指定managed 规则. </li>。",
												},
												"skip_option": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "特定 类型 skipped 请求. possible 值 是: <li>SkipOnAllRequestFields: skip all requests; </li><li>SkipOnSpecifiedRequestFields: skip 指定 请求 字段. </li>. 此 选项 是 仅 有效 当 SkipScope 是 ManagedRules。",
												},
												"web_security_modules_for_exception": {
													Type:        schema.TypeSet,
													Optional:    true,
													Description: "指定security protection 模块 对于 exception 规则. It 是 有效 仅 当 SkipScope 是 WebSecurityModules. possible 值 是: <li>websec-mod-managed-规则: managed 规则; </li><li>websec-mod-速率-limiting: 速率 limiting; </li><li>websec-mod-自定义-规则: 自定义 规则; </li><li>websec-mod-adaptive-control: adaptive 频率 control，intelligent 客户端 filtering，slow attack protection，流量 theft protection; </li><li>websec-mod-bot: Bot management. </li>。",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"managed_rules_for_exception": {
													Type:        schema.TypeSet,
													Optional:    true,
													Description: "指定specific managed 规则 对于 exception 规则. 此 是 仅 有效 当 SkipScope 是 ManagedRules 和 ManagedRuleGroupsForException 不能 是 指定。",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"managed_rule_groups_for_exception": {
													Type:        schema.TypeSet,
													Optional:    true,
													Description: "指定managed 规则 组 对于 exception 规则. 此 是 仅 有效 当 SkipScope 是 ManagedRules 和 ManagedRulesForException 不能 是 指定。",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"request_fields_for_exception": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "指定specific 配置 的 exception 规则 到 skip 指定 请求 字段. 此 是 仅 有效 当 SkipScope 是 ManagedRules 和 SkipOption 是 SkipOnSpecifiedRequestFields。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"scope": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specific 字段 到 skip. Supported 值:<br/>\n<li>正文.json: JSON 请求 内容; 在 此 case, Condition 支持 键 和 值, 和 TargetField 支持 键 和 值, 对于 示例, { \"Scope\": \"正文.json\", \"Condition\": \"\", \"TargetField\": \"键\" }, 其中 表示 该 all 参数 的 JSON 请求 内容 skip WAF scanning;</li>\n<li style=\"margin-top:5px\">cookie: Cookie; 在 此 case, Condition 支持 键 和 值, 和 TargetField 支持 键 和 值, 对于 示例, { \"Scope\": \"cookie\", \"Condition\": \"${键} 在 ['account-ID'] 和 ${值} like ['prefix-*']\", \"TargetField\": \"值\" }, 其中 表示 该 Cookie 参数 名称 是 equal 到 account-ID 和 参数 值 wildcard matches prefix-* 到 skip WAF scanning;</li>\n<li style=\"margin-top:5px\">头部: HTTP 头部 参数; Condition 支持 键 和 值, TargetField 支持 键 和 值, 对于 示例 { \"Scope\": \"头部\", \"Condition\": \"${键} like ['x-auth-*']\", \"TargetField\": \"值\" }, 其中 表示 该 头部 参数 名称 wildcard matches x-auth-* 和 skips WAF scanning; </li>\n<li style=\"margin-top:5px\">uri.查询: URL encoded 内容/查询 参数; Condition 支持 键 和 值, TargetField 支持 键 和 值, 对于 示例 { \"Scope\": \"uri.查询\", \"Condition\": \"${键} 在 ['action'] 和 ${值} 在 ['upload', 'delete']\", \"TargetField\": \"值\" }, 其中 表示 该 参数 名称 的 URL encoded 内容/查询 参数 是 equal 到 action And 参数 值 是 equal 到 upload 或 delete 到 skip WAF scanning;</li>\n<li style=\"margin-top:5px\">uri: 请求 路径 URI; 在 此 case, Condition 必须 是 空, TargetField 支持 查询, 路径, fullpath, 对于 示例, { \"Scope\": \"uri\", \"Condition\": \"\", \"TargetField\": \"查询\" }, indicating 该 请求 路径 URI 仅 查询 参数 skip WAF scanning;</li>\n<li style=\"margin-top:5px\">正文: 请求 正文 内容. In 此 case, Condition 必须 是 空, TargetField 支持 fullbody 和 multipart, 对于 示例, { \"Scope\": \"正文\", \"Condition\": \"\", \"TargetField\": \"fullbody\" }, indicating 该 请求 正文 内容 是 完整 请求 正文 和 skips WAF scanning;</li>.",
															},
															"condition": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "expression 的 特定 字段 到 是 skipped 必须 conform 到 expression syntax. <br />\nCondition 支持 expression 配置 syntax: <li> Written according 到 matching condition expression syntax 的 规则，supporting references 到 键 和 值 </li>.<li> Supports 在，like operators，和 和 logical combinations. </li>.\nFor 示例: <li>${键} 在 ['x-trace-ID']: 参数 名称 是 equal 到 x-trace-ID. </li>.<li>${键} 在 ['x-trace-ID'] 和 ${值} like ['Bearer *']: 参数 名称 是 equal 到 x-trace-ID 和 参数 值 wildcard matches Bearer *. </li>。",
															},
															"target_field": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "当 范围 参数 uses different 值， 支持 值 在 TargetField expression 是 作为 follows:\n<li> 正文.json: 支持 键 和 值</li>\n<li> cookie: 支持 键 和 值</li>\n<li> 头部: 支持 键 和 值</li>\n<li> uri.查询: 支持 键 和 值</li>\n<li> uri: 支持 路径，查询 和 fullpath</li>\n<li> 正文: 支持 fullbody 和 multipart</li>。",
															},
														},
													},
												},
												"enabled": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "是否exception 规则 是 已启用 值 是: <li>在: 已启用</li><li>关闭: 已禁用</li>。",
												},
											},
										},
									},
								},
							},
						},
						"bot_management": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Bot management 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Whether Bot management 是 已启用 有效值：<li>在: 已启用</li> <li>关闭: 已禁用</li>。",
									},
									"custom_rules": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Bot 自定义 规则 配置。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rules": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "列表 Bot 自定义 规则。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Rule ID. add new 规则 当 ID 是 空; modify existing 规则 当 ID 是 指定; existing 规则 不 included 在 列表 将 是 删除。",
															},
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Rule 名称",
															},
															"enabled": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "是否rule 是 已启用 有效值：<li>在: 已启用</li> <li>关闭: 已禁用</li>。",
															},
															"priority": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "优先级 的 规则. 取值范围：1-100，默认值 50。",
															},
															"condition": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Match condition expression. 必须 comply 使用 expression grammar。",
															},
															"action": {
																Type:        schema.TypeList,
																Required:    true,
																Description: "Weighted 操作 列表. all weights 必须 sum 到 100。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"security_action": {
																			Type:        schema.TypeList,
																			Required:    true,
																			MaxItems:    1,
																			Description: "Security 操作 配置。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"name": {
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "操作 名称 有效值：Deny，Monitor，Allow，Challenge，已禁用，Redirect，Trans。",
																					},
																					"deny_action_parameters": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						MaxItems:    1,
																						Description: "Additional 参数 当 名称 是 Deny。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"block_ip": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "是否enable IP ban extension. 有效值：在/关闭。",
																								},
																								"block_ip_duration": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "IP ban 时长. 必填 当 BlockIp 是 在。",
																								},
																								"return_custom_page": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "是否use 自定义 页面. 有效值：在/关闭。",
																								},
																								"response_code": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "Custom 页面 response 状态 代码",
																								},
																								"error_page_id": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "Custom 页面 ID。",
																								},
																								"stall": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "是否stall 连接. 有效值：在/关闭。",
																								},
																							},
																						},
																					},
																					"redirect_action_parameters": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						MaxItems:    1,
																						Description: "Additional 参数 当 名称 是 Redirect。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"url": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "Redirect 目标 URL",
																								},
																							},
																						},
																					},
																					"allow_action_parameters": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						MaxItems:    1,
																						Description: "Additional 参数 当 名称 是 Allow。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"min_delay_time": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "Minimum 延迟 response 时间。",
																								},
																								"max_delay_time": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "Maximum 延迟 response 时间。",
																								},
																							},
																						},
																					},
																					"challenge_action_parameters": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						MaxItems:    1,
																						Description: "Additional 参数 当 名称 是 Challenge。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"challenge_option": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "Challenge 类型 有效值：InterstitialChallenge，InlineChallenge，JSChallenge，ManagedChallenge。",
																								},
																								"interval": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "Repeat challenge 间隔。",
																								},
																								"attester_id": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "Client attestation 方法 ID。",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"weight": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "操作 权重 取值范围：10-100，必须 是 multiples 的 10，sum 的 all weights 必须 equal 100。",
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
									"basic_bot_settings": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Basic Bot settings 配置 该 applies 到 all domains。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"source_idc": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "IDC 来源 IP 配置。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"base_action": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Default 操作 对于 IDC 来源 requests。",
																Elem:        securityActionSchemaForBotManagement(),
															},
															"action_overrides": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "操作 overrides 对于 特定 IDC 来源 规则 IDs。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ids": {
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: "Rule IDs 到 override。",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"action": {
																			Type:        schema.TypeList,
																			Required:    true,
																			MaxItems:    1,
																			Description: "Override 操作",
																			Elem:        securityActionSchemaForBotManagement(),
																		},
																	},
																},
															},
														},
													},
												},
												"search_engine_bots": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Search 引擎 bot 配置。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"base_action": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Default 操作 对于 search 引擎 bot requests。",
																Elem:        securityActionSchemaForBotManagement(),
															},
															"action_overrides": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "操作 overrides 对于 特定 search 引擎 bot 规则 IDs。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ids": {
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: "Rule IDs 到 override。",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"action": {
																			Type:        schema.TypeList,
																			Required:    true,
																			MaxItems:    1,
																			Description: "Override 操作",
																			Elem:        securityActionSchemaForBotManagement(),
																		},
																	},
																},
															},
														},
													},
												},
												"known_bot_categories": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Known Bot category 配置。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"base_action": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Default 操作 对于 known bot category requests。",
																Elem:        securityActionSchemaForBotManagement(),
															},
															"action_overrides": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "操作 overrides 对于 特定 known bot category 规则 IDs。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"ids": {
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: "Rule IDs 到 override。",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"action": {
																			Type:        schema.TypeList,
																			Required:    true,
																			MaxItems:    1,
																			Description: "Override 操作",
																			Elem:        securityActionSchemaForBotManagement(),
																		},
																	},
																},
															},
														},
													},
												},
												"ip_reputation": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "IP threat intelligence 配置。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"enabled": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Whether IP reputation 是 已启用 有效值：在/关闭。",
															},
															"ip_reputation_group": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "IP reputation 组 配置。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"base_action": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "Default 操作 对于 IP reputation requests。",
																			Elem:        securityActionSchemaForBotManagement(),
																		},
																		"action_overrides": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "操作 overrides 对于 特定 IP reputation 规则 IDs。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"ids": {
																						Type:        schema.TypeList,
																						Required:    true,
																						Description: "Rule IDs 到 override。",
																						Elem:        &schema.Schema{Type: schema.TypeString},
																					},
																					"action": {
																						Type:        schema.TypeList,
																						Required:    true,
																						MaxItems:    1,
																						Description: "Override 操作",
																						Elem:        securityActionSchemaForBotManagement(),
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
												"bot_intelligence": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Bot intelligence analysis 配置。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"enabled": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Whether Bot intelligence 是 已启用 有效值：在/关闭。",
															},
															"id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Rule ID (output 仅)。",
															},
															"bot_ratings": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "Bot rating-based 操作 配置。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"high_risk_bot_requests_action": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "操作 对于 high risk Bot requests。",
																			Elem:        securityActionSchemaForBotManagement(),
																		},
																		"likely_bot_requests_action": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "操作 对于 likely Bot requests。",
																			Elem:        securityActionSchemaForBotManagement(),
																		},
																		"verified_bot_requests_action": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "操作 对于 verified Bot requests。",
																			Elem:        securityActionSchemaForBotManagement(),
																		},
																		"human_requests_action": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "操作 对于 human requests。",
																			Elem:        securityActionSchemaForBotManagement(),
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
									"client_attestation_rules": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Client attestation 规则 配置 (beta 功能，requires support ticket)。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rules": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "列表 客户端 attestation 规则。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Rule ID。",
															},
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Rule 名称",
															},
															"enabled": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "是否rule 是 已启用 有效值：在/关闭。",
															},
															"priority": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "优先级 的 规则. 取值范围：0-100，默认值 0。",
															},
															"condition": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Match condition expression。",
															},
															"attester_id": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Client attestation 选项 ID。",
															},
															"device_profiles": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Device profiles 配置。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"client_type": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Client 类型 有效值：iOS，Android，WebView，WeChatMiniProgram。",
																		},
																		"high_risk_min_score": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Minimum score 对于 high risk. 取值范围：1-99，默认值 50。",
																		},
																		"high_risk_request_action": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "操作 对于 high risk requests。",
																			Elem:        securityActionSchemaForBotManagement(),
																		},
																		"medium_risk_min_score": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Minimum score 对于 medium risk. 取值范围：1-99，默认值 15。",
																		},
																		"medium_risk_request_action": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "操作 对于 medium risk requests。",
																			Elem:        securityActionSchemaForBotManagement(),
																		},
																	},
																},
															},
															"invalid_attestation_action": {
																Type:        schema.TypeList,
																Optional:    true,
																MaxItems:    1,
																Description: "操作 当 attestation 是 无效。",
																Elem:        securityActionSchemaForBotManagement(),
															},
														},
													},
												},
											},
										},
									},
									"browser_impersonation_detection": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Browser impersonation detection 配置。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rules": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "列表 browser impersonation detection 规则。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Rule ID。",
															},
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Rule 名称",
															},
															"enabled": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "是否rule 是 已启用 有效值：在/关闭。",
															},
															"condition": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Match condition expression。",
															},
															"action": {
																Type:        schema.TypeList,
																Required:    true,
																MaxItems:    1,
																Description: "Browser impersonation detection 操作",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bot_session_validation": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "Bot 会话 validation 配置。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"issue_new_bot_session_cookie": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "是否issue new Bot 会话 cookie. 有效值：在 (update+validate)，关闭 (validate 仅)。",
																					},
																					"max_new_session_trigger_config": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						MaxItems:    1,
																						Description: "Trigger 阈值 对于 cookie renewal。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"max_new_session_count_interval": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "Time window 对于 new 会话 count。",
																								},
																								"max_new_session_count_threshold": {
																									Type:        schema.TypeInt,
																									Optional:    true,
																									Description: "Cumulative count 阈值 对于 new sessions。",
																								},
																							},
																						},
																					},
																					"session_expired_action": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						MaxItems:    1,
																						Description: "操作 对于 expired sessions。",
																						Elem:        securityActionSchemaForBotManagement(),
																					},
																					"session_invalid_action": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						MaxItems:    1,
																						Description: "操作 对于 无效 sessions。",
																						Elem:        securityActionSchemaForBotManagement(),
																					},
																					"session_rate_control": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						MaxItems:    1,
																						Description: "Session 速率 control 配置。",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"enabled": {
																									Type:        schema.TypeString,
																									Optional:    true,
																									Description: "Whether 会话 速率 control 是 已启用 有效值：在/关闭。",
																								},
																								"high_rate_session_action": {
																									Type:        schema.TypeList,
																									Optional:    true,
																									MaxItems:    1,
																									Description: "操作 对于 high 速率 sessions。",
																									Elem:        securityActionSchemaForBotManagement(),
																								},
																								"mid_rate_session_action": {
																									Type:        schema.TypeList,
																									Optional:    true,
																									MaxItems:    1,
																									Description: "操作 对于 medium 速率 sessions。",
																									Elem:        securityActionSchemaForBotManagement(),
																								},
																								"low_rate_session_action": {
																									Type:        schema.TypeList,
																									Optional:    true,
																									MaxItems:    1,
																									Description: "操作 对于 low 速率 sessions。",
																									Elem:        securityActionSchemaForBotManagement(),
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"client_behavior_detection": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			MaxItems:    1,
																			Description: "Client behavior detection 配置。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"crypto_challenge_intensity": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Proof-的-work intensity. 有效值：low，medium，high。",
																					},
																					"crypto_challenge_delay_before": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Execution 延迟 before challenge. 有效 范围: 0ms-1000ms。",
																					},
																					"max_challenge_count_interval": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Threshold 时间 window。",
																					},
																					"max_challenge_count_threshold": {
																						Type:        schema.TypeInt,
																						Optional:    true,
																						Description: "Threshold cumulative count. 取值范围：1-100000000。",
																					},
																					"challenge_not_finished_action": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						MaxItems:    1,
																						Description: "操作 当 JS challenge 不 finished。",
																						Elem:        securityActionSchemaForBotManagement(),
																					},
																					"challenge_timeout_action": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						MaxItems:    1,
																						Description: "操作 在 challenge 超时。",
																						Elem:        securityActionSchemaForBotManagement(),
																					},
																					"bot_client_action": {
																						Type:        schema.TypeList,
																						Optional:    true,
																						MaxItems:    1,
																						Description: "操作 对于 Bot 客户端。",
																						Elem:        securityActionSchemaForBotManagement(),
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
								},
							},
						},
					},
				},
			},

			"entity": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"ZoneDefaultPolicy", "Template", "Host"}),
				Description:  "Security 策略 类型 following 参数 值 可以 是 使用: <li>ZoneDefaultPolicy: 用于指定a site-级别 策略;</li> <li>模板: 用于指定a 策略 模板. 您 need 到 simultaneously 指定TemplateId 参数;</li> <li>主机: 用于指定a 域名-级别 策略 (note: 当 使用 域名 名称 到 指定a dns 服务 策略，仅 dns services 或 策略 templates 该 have applied 域名-级别 策略 是 支持).</li>。",
			},

			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "指定specified 域名 当 Entity 参数 值 是 主机，使用 域名-级别 策略 指定 通过 此 参数. 对于 示例: 使用 www.示例.com 到 configure 域名-级别 策略 的 域名",
			},

			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "指定policy 模板 ID 使用 此 参数 到 指定ID 的 策略 模板 当 Entity 参数 值 是 模板。",
			},
		},
	}
}

func securityActionSchemaForBotManagement() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "操作 名称 有效值：Deny，Monitor，Allow，Challenge，已禁用，Redirect，Trans。",
			},
			"deny_action_parameters": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Additional 参数 当 名称 是 Deny。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"block_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "是否enable IP ban extension. 有效值：在/关闭。",
						},
						"block_ip_duration": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IP ban 时长。",
						},
						"return_custom_page": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "是否use 自定义 页面. 有效值：在/关闭。",
						},
						"response_code": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom 页面 response 状态 代码",
						},
						"error_page_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom 页面 ID。",
						},
						"stall": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "是否stall 连接. 有效值：在/关闭。",
						},
					},
				},
			},
			"redirect_action_parameters": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Additional 参数 当 名称 是 Redirect。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Redirect 目标 URL",
						},
					},
				},
			},
			"allow_action_parameters": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Additional 参数 当 名称 是 Allow。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_delay_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Minimum 延迟 response 时间。",
						},
						"max_delay_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Maximum 延迟 response 时间。",
						},
					},
				},
			},
			"challenge_action_parameters": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Additional 参数 当 名称 是 Challenge。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"challenge_option": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Challenge 类型 有效值：InterstitialChallenge，InlineChallenge，JSChallenge，ManagedChallenge。",
						},
						"interval": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Repeat challenge 间隔。",
						},
						"attester_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Client attestation 方法 ID。",
						},
					},
				},
			},
		},
	}
}

func flattenSecurityActionForBotManagement(action *teov20220901.SecurityAction) map[string]interface{} {
	actionMap := map[string]interface{}{}
	if action == nil {
		return actionMap
	}
	if action.Name != nil {
		actionMap["name"] = action.Name
	}
	if action.DenyActionParameters != nil {
		denyMap := map[string]interface{}{}
		if action.DenyActionParameters.BlockIp != nil {
			denyMap["block_ip"] = action.DenyActionParameters.BlockIp
		}
		if action.DenyActionParameters.BlockIpDuration != nil {
			denyMap["block_ip_duration"] = action.DenyActionParameters.BlockIpDuration
		}
		if action.DenyActionParameters.ReturnCustomPage != nil {
			denyMap["return_custom_page"] = action.DenyActionParameters.ReturnCustomPage
		}
		if action.DenyActionParameters.ResponseCode != nil {
			denyMap["response_code"] = action.DenyActionParameters.ResponseCode
		}
		if action.DenyActionParameters.ErrorPageId != nil {
			denyMap["error_page_id"] = action.DenyActionParameters.ErrorPageId
		}
		if action.DenyActionParameters.Stall != nil {
			denyMap["stall"] = action.DenyActionParameters.Stall
		}
		actionMap["deny_action_parameters"] = []interface{}{denyMap}
	}
	if action.RedirectActionParameters != nil {
		redirectMap := map[string]interface{}{}
		if action.RedirectActionParameters.URL != nil {
			redirectMap["url"] = action.RedirectActionParameters.URL
		}
		actionMap["redirect_action_parameters"] = []interface{}{redirectMap}
	}
	if action.AllowActionParameters != nil {
		allowMap := map[string]interface{}{}
		if action.AllowActionParameters.MinDelayTime != nil {
			allowMap["min_delay_time"] = action.AllowActionParameters.MinDelayTime
		}
		if action.AllowActionParameters.MaxDelayTime != nil {
			allowMap["max_delay_time"] = action.AllowActionParameters.MaxDelayTime
		}
		actionMap["allow_action_parameters"] = []interface{}{allowMap}
	}
	if action.ChallengeActionParameters != nil {
		challengeMap := map[string]interface{}{}
		if action.ChallengeActionParameters.ChallengeOption != nil {
			challengeMap["challenge_option"] = action.ChallengeActionParameters.ChallengeOption
		}
		if action.ChallengeActionParameters.Interval != nil {
			challengeMap["interval"] = action.ChallengeActionParameters.Interval
		}
		if action.ChallengeActionParameters.AttesterId != nil {
			challengeMap["attester_id"] = action.ChallengeActionParameters.AttesterId
		}
		actionMap["challenge_action_parameters"] = []interface{}{challengeMap}
	}
	return actionMap
}

func flattenActionOverrides(overrides []*teov20220901.BotManagementActionOverrides) []interface{} {
	result := make([]interface{}, 0, len(overrides))
	for _, override := range overrides {
		overrideMap := map[string]interface{}{}
		if override.Ids != nil {
			idsList := make([]interface{}, 0, len(override.Ids))
			for _, id := range override.Ids {
				if id != nil {
					idsList = append(idsList, *id)
				}
			}
			overrideMap["ids"] = idsList
		}
		if override.Action != nil {
			overrideMap["action"] = []interface{}{flattenSecurityActionForBotManagement(override.Action)}
		}
		result = append(result, overrideMap)
	}
	return result
}

func flattenBaseActionAndOverrides(baseAction *teov20220901.SecurityAction, actionOverrides []*teov20220901.BotManagementActionOverrides) map[string]interface{} {
	result := map[string]interface{}{}
	if baseAction != nil {
		result["base_action"] = []interface{}{flattenSecurityActionForBotManagement(baseAction)}
	}
	if actionOverrides != nil {
		result["action_overrides"] = flattenActionOverrides(actionOverrides)
	}
	return result
}

func flattenBotManagement(botManagement *teov20220901.BotManagement) map[string]interface{} {
	botManagementMap := map[string]interface{}{}
	if botManagement == nil {
		return botManagementMap
	}
	if botManagement.Enabled != nil {
		botManagementMap["enabled"] = botManagement.Enabled
	}
	if botManagement.CustomRules != nil {
		customRulesMap := map[string]interface{}{}
		if botManagement.CustomRules.Rules != nil {
			rulesList := make([]interface{}, 0, len(botManagement.CustomRules.Rules))
			for _, rule := range botManagement.CustomRules.Rules {
				ruleMap := map[string]interface{}{}
				if rule.Id != nil {
					ruleMap["id"] = rule.Id
				}
				if rule.Name != nil {
					ruleMap["name"] = rule.Name
				}
				if rule.Enabled != nil {
					ruleMap["enabled"] = rule.Enabled
				}
				if rule.Priority != nil {
					ruleMap["priority"] = rule.Priority
				}
				if rule.Condition != nil {
					ruleMap["condition"] = rule.Condition
				}
				if rule.Action != nil {
					actionList := make([]interface{}, 0, len(rule.Action))
					for _, weightedAction := range rule.Action {
						weightedActionMap := map[string]interface{}{}
						if weightedAction.SecurityAction != nil {
							weightedActionMap["security_action"] = []interface{}{flattenSecurityActionForBotManagement(weightedAction.SecurityAction)}
						}
						if weightedAction.Weight != nil {
							weightedActionMap["weight"] = weightedAction.Weight
						}
						actionList = append(actionList, weightedActionMap)
					}
					ruleMap["action"] = actionList
				}
				rulesList = append(rulesList, ruleMap)
			}
			customRulesMap["rules"] = rulesList
		}
		botManagementMap["custom_rules"] = []interface{}{customRulesMap}
	}
	if botManagement.BasicBotSettings != nil {
		basicBotSettingsMap := map[string]interface{}{}
		if botManagement.BasicBotSettings.SourceIDC != nil {
			sourceIDCMap := flattenBaseActionAndOverrides(botManagement.BasicBotSettings.SourceIDC.BaseAction, botManagement.BasicBotSettings.SourceIDC.BotManagementActionOverrides)
			basicBotSettingsMap["source_idc"] = []interface{}{sourceIDCMap}
		}
		if botManagement.BasicBotSettings.SearchEngineBots != nil {
			searchEngineBotsMap := flattenBaseActionAndOverrides(botManagement.BasicBotSettings.SearchEngineBots.BaseAction, botManagement.BasicBotSettings.SearchEngineBots.BotManagementActionOverrides)
			basicBotSettingsMap["search_engine_bots"] = []interface{}{searchEngineBotsMap}
		}
		if botManagement.BasicBotSettings.KnownBotCategories != nil {
			knownBotCategoriesMap := flattenBaseActionAndOverrides(botManagement.BasicBotSettings.KnownBotCategories.BaseAction, botManagement.BasicBotSettings.KnownBotCategories.BotManagementActionOverrides)
			basicBotSettingsMap["known_bot_categories"] = []interface{}{knownBotCategoriesMap}
		}
		if botManagement.BasicBotSettings.IPReputation != nil {
			ipReputationMap := map[string]interface{}{}
			if botManagement.BasicBotSettings.IPReputation.Enabled != nil {
				ipReputationMap["enabled"] = botManagement.BasicBotSettings.IPReputation.Enabled
			}
			if botManagement.BasicBotSettings.IPReputation.IPReputationGroup != nil {
				ipReputationGroupMap := flattenBaseActionAndOverrides(
					botManagement.BasicBotSettings.IPReputation.IPReputationGroup.BaseAction,
					botManagement.BasicBotSettings.IPReputation.IPReputationGroup.BotManagementActionOverrides,
				)
				ipReputationMap["ip_reputation_group"] = []interface{}{ipReputationGroupMap}
			}
			basicBotSettingsMap["ip_reputation"] = []interface{}{ipReputationMap}
		}
		if botManagement.BasicBotSettings.BotIntelligence != nil {
			botIntelligenceMap := map[string]interface{}{}
			if botManagement.BasicBotSettings.BotIntelligence.Enabled != nil {
				botIntelligenceMap["enabled"] = botManagement.BasicBotSettings.BotIntelligence.Enabled
			}
			if botManagement.BasicBotSettings.BotIntelligence.Id != nil {
				botIntelligenceMap["id"] = botManagement.BasicBotSettings.BotIntelligence.Id
			}
			if botManagement.BasicBotSettings.BotIntelligence.BotRatings != nil {
				botRatingsMap := map[string]interface{}{}
				if botManagement.BasicBotSettings.BotIntelligence.BotRatings.HighRiskBotRequestsAction != nil {
					botRatingsMap["high_risk_bot_requests_action"] = []interface{}{flattenSecurityActionForBotManagement(botManagement.BasicBotSettings.BotIntelligence.BotRatings.HighRiskBotRequestsAction)}
				}
				if botManagement.BasicBotSettings.BotIntelligence.BotRatings.LikelyBotRequestsAction != nil {
					botRatingsMap["likely_bot_requests_action"] = []interface{}{flattenSecurityActionForBotManagement(botManagement.BasicBotSettings.BotIntelligence.BotRatings.LikelyBotRequestsAction)}
				}
				if botManagement.BasicBotSettings.BotIntelligence.BotRatings.VerifiedBotRequestsAction != nil {
					botRatingsMap["verified_bot_requests_action"] = []interface{}{flattenSecurityActionForBotManagement(botManagement.BasicBotSettings.BotIntelligence.BotRatings.VerifiedBotRequestsAction)}
				}
				if botManagement.BasicBotSettings.BotIntelligence.BotRatings.HumanRequestsAction != nil {
					botRatingsMap["human_requests_action"] = []interface{}{flattenSecurityActionForBotManagement(botManagement.BasicBotSettings.BotIntelligence.BotRatings.HumanRequestsAction)}
				}
				botIntelligenceMap["bot_ratings"] = []interface{}{botRatingsMap}
			}
			basicBotSettingsMap["bot_intelligence"] = []interface{}{botIntelligenceMap}
		}
		botManagementMap["basic_bot_settings"] = []interface{}{basicBotSettingsMap}
	}
	if botManagement.ClientAttestationRules != nil {
		clientAttestationRulesMap := map[string]interface{}{}
		if botManagement.ClientAttestationRules.Rules != nil {
			rulesList := make([]interface{}, 0, len(botManagement.ClientAttestationRules.Rules))
			for _, rule := range botManagement.ClientAttestationRules.Rules {
				ruleMap := map[string]interface{}{}
				if rule.Id != nil {
					ruleMap["id"] = rule.Id
				}
				if rule.Name != nil {
					ruleMap["name"] = rule.Name
				}
				if rule.Enabled != nil {
					ruleMap["enabled"] = rule.Enabled
				}
				if rule.Priority != nil {
					ruleMap["priority"] = rule.Priority
				}
				if rule.Condition != nil {
					ruleMap["condition"] = rule.Condition
				}
				if rule.AttesterId != nil {
					ruleMap["attester_id"] = rule.AttesterId
				}
				if rule.DeviceProfiles != nil {
					deviceProfilesList := make([]interface{}, 0, len(rule.DeviceProfiles))
					for _, profile := range rule.DeviceProfiles {
						profileMap := map[string]interface{}{}
						if profile.ClientType != nil {
							profileMap["client_type"] = profile.ClientType
						}
						if profile.HighRiskMinScore != nil {
							profileMap["high_risk_min_score"] = profile.HighRiskMinScore
						}
						if profile.HighRiskRequestAction != nil {
							profileMap["high_risk_request_action"] = []interface{}{flattenSecurityActionForBotManagement(profile.HighRiskRequestAction)}
						}
						if profile.MediumRiskMinScore != nil {
							profileMap["medium_risk_min_score"] = profile.MediumRiskMinScore
						}
						if profile.MediumRiskRequestAction != nil {
							profileMap["medium_risk_request_action"] = []interface{}{flattenSecurityActionForBotManagement(profile.MediumRiskRequestAction)}
						}
						deviceProfilesList = append(deviceProfilesList, profileMap)
					}
					ruleMap["device_profiles"] = deviceProfilesList
				}
				if rule.InvalidAttestationAction != nil {
					ruleMap["invalid_attestation_action"] = []interface{}{flattenSecurityActionForBotManagement(rule.InvalidAttestationAction)}
				}
				rulesList = append(rulesList, ruleMap)
			}
			clientAttestationRulesMap["rules"] = rulesList
		}
		botManagementMap["client_attestation_rules"] = []interface{}{clientAttestationRulesMap}
	}
	if botManagement.BrowserImpersonationDetection != nil {
		browserImpersonationDetectionMap := map[string]interface{}{}
		if botManagement.BrowserImpersonationDetection.Rules != nil {
			rulesList := make([]interface{}, 0, len(botManagement.BrowserImpersonationDetection.Rules))
			for _, rule := range botManagement.BrowserImpersonationDetection.Rules {
				ruleMap := map[string]interface{}{}
				if rule.Id != nil {
					ruleMap["id"] = rule.Id
				}
				if rule.Name != nil {
					ruleMap["name"] = rule.Name
				}
				if rule.Enabled != nil {
					ruleMap["enabled"] = rule.Enabled
				}
				if rule.Condition != nil {
					ruleMap["condition"] = rule.Condition
				}
				if rule.Action != nil {
					actionMap := map[string]interface{}{}
					if rule.Action.BotSessionValidation != nil {
						bsvMap := map[string]interface{}{}
						if rule.Action.BotSessionValidation.IssueNewBotSessionCookie != nil {
							bsvMap["issue_new_bot_session_cookie"] = rule.Action.BotSessionValidation.IssueNewBotSessionCookie
						}
						if rule.Action.BotSessionValidation.MaxNewSessionTriggerConfig != nil {
							triggerConfigMap := map[string]interface{}{}
							if rule.Action.BotSessionValidation.MaxNewSessionTriggerConfig.MaxNewSessionCountInterval != nil {
								triggerConfigMap["max_new_session_count_interval"] = rule.Action.BotSessionValidation.MaxNewSessionTriggerConfig.MaxNewSessionCountInterval
							}
							if rule.Action.BotSessionValidation.MaxNewSessionTriggerConfig.MaxNewSessionCountThreshold != nil {
								triggerConfigMap["max_new_session_count_threshold"] = rule.Action.BotSessionValidation.MaxNewSessionTriggerConfig.MaxNewSessionCountThreshold
							}
							bsvMap["max_new_session_trigger_config"] = []interface{}{triggerConfigMap}
						}
						if rule.Action.BotSessionValidation.SessionExpiredAction != nil {
							bsvMap["session_expired_action"] = []interface{}{flattenSecurityActionForBotManagement(rule.Action.BotSessionValidation.SessionExpiredAction)}
						}
						if rule.Action.BotSessionValidation.SessionInvalidAction != nil {
							bsvMap["session_invalid_action"] = []interface{}{flattenSecurityActionForBotManagement(rule.Action.BotSessionValidation.SessionInvalidAction)}
						}
						if rule.Action.BotSessionValidation.SessionRateControl != nil {
							srcMap := map[string]interface{}{}
							if rule.Action.BotSessionValidation.SessionRateControl.Enabled != nil {
								srcMap["enabled"] = rule.Action.BotSessionValidation.SessionRateControl.Enabled
							}
							if rule.Action.BotSessionValidation.SessionRateControl.HighRateSessionAction != nil {
								srcMap["high_rate_session_action"] = []interface{}{flattenSecurityActionForBotManagement(rule.Action.BotSessionValidation.SessionRateControl.HighRateSessionAction)}
							}
							if rule.Action.BotSessionValidation.SessionRateControl.MidRateSessionAction != nil {
								srcMap["mid_rate_session_action"] = []interface{}{flattenSecurityActionForBotManagement(rule.Action.BotSessionValidation.SessionRateControl.MidRateSessionAction)}
							}
							if rule.Action.BotSessionValidation.SessionRateControl.LowRateSessionAction != nil {
								srcMap["low_rate_session_action"] = []interface{}{flattenSecurityActionForBotManagement(rule.Action.BotSessionValidation.SessionRateControl.LowRateSessionAction)}
							}
							bsvMap["session_rate_control"] = []interface{}{srcMap}
						}
						actionMap["bot_session_validation"] = []interface{}{bsvMap}
					}
					if rule.Action.ClientBehaviorDetection != nil {
						cbdMap := map[string]interface{}{}
						if rule.Action.ClientBehaviorDetection.CryptoChallengeIntensity != nil {
							cbdMap["crypto_challenge_intensity"] = rule.Action.ClientBehaviorDetection.CryptoChallengeIntensity
						}
						if rule.Action.ClientBehaviorDetection.CryptoChallengeDelayBefore != nil {
							cbdMap["crypto_challenge_delay_before"] = rule.Action.ClientBehaviorDetection.CryptoChallengeDelayBefore
						}
						if rule.Action.ClientBehaviorDetection.MaxChallengeCountInterval != nil {
							cbdMap["max_challenge_count_interval"] = rule.Action.ClientBehaviorDetection.MaxChallengeCountInterval
						}
						if rule.Action.ClientBehaviorDetection.MaxChallengeCountThreshold != nil {
							cbdMap["max_challenge_count_threshold"] = rule.Action.ClientBehaviorDetection.MaxChallengeCountThreshold
						}
						if rule.Action.ClientBehaviorDetection.ChallengeNotFinishedAction != nil {
							cbdMap["challenge_not_finished_action"] = []interface{}{flattenSecurityActionForBotManagement(rule.Action.ClientBehaviorDetection.ChallengeNotFinishedAction)}
						}
						if rule.Action.ClientBehaviorDetection.ChallengeTimeoutAction != nil {
							cbdMap["challenge_timeout_action"] = []interface{}{flattenSecurityActionForBotManagement(rule.Action.ClientBehaviorDetection.ChallengeTimeoutAction)}
						}
						if rule.Action.ClientBehaviorDetection.BotClientAction != nil {
							cbdMap["bot_client_action"] = []interface{}{flattenSecurityActionForBotManagement(rule.Action.ClientBehaviorDetection.BotClientAction)}
						}
						actionMap["client_behavior_detection"] = []interface{}{cbdMap}
					}
					ruleMap["action"] = []interface{}{actionMap}
				}
				rulesList = append(rulesList, ruleMap)
			}
			browserImpersonationDetectionMap["rules"] = rulesList
		}
		botManagementMap["browser_impersonation_detection"] = []interface{}{browserImpersonationDetectionMap}
	}
	return botManagementMap
}

func expandSecurityActionForBotManagement(actionMap map[string]interface{}) *teov20220901.SecurityAction {
	securityAction := teov20220901.SecurityAction{}
	if v, ok := actionMap["name"].(string); ok && v != "" {
		securityAction.Name = helper.String(v)
	}
	if denyMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["deny_action_parameters"]); ok {
		denyActionParameters := teov20220901.DenyActionParameters{}
		if v, ok := denyMap["block_ip"].(string); ok && v != "" {
			denyActionParameters.BlockIp = helper.String(v)
		}
		if v, ok := denyMap["block_ip_duration"].(string); ok && v != "" {
			denyActionParameters.BlockIpDuration = helper.String(v)
		}
		if v, ok := denyMap["return_custom_page"].(string); ok && v != "" {
			denyActionParameters.ReturnCustomPage = helper.String(v)
		}
		if v, ok := denyMap["response_code"].(string); ok && v != "" {
			denyActionParameters.ResponseCode = helper.String(v)
		}
		if v, ok := denyMap["error_page_id"].(string); ok && v != "" {
			denyActionParameters.ErrorPageId = helper.String(v)
		}
		if v, ok := denyMap["stall"].(string); ok && v != "" {
			denyActionParameters.Stall = helper.String(v)
		}
		securityAction.DenyActionParameters = &denyActionParameters
	}
	if redirectMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["redirect_action_parameters"]); ok {
		redirectActionParameters := teov20220901.RedirectActionParameters{}
		if v, ok := redirectMap["url"].(string); ok && v != "" {
			redirectActionParameters.URL = helper.String(v)
		}
		securityAction.RedirectActionParameters = &redirectActionParameters
	}
	if allowMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["allow_action_parameters"]); ok {
		allowActionParameters := teov20220901.AllowActionParameters{}
		if v, ok := allowMap["min_delay_time"].(string); ok && v != "" {
			allowActionParameters.MinDelayTime = helper.String(v)
		}
		if v, ok := allowMap["max_delay_time"].(string); ok && v != "" {
			allowActionParameters.MaxDelayTime = helper.String(v)
		}
		securityAction.AllowActionParameters = &allowActionParameters
	}
	if challengeMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["challenge_action_parameters"]); ok {
		challengeActionParameters := teov20220901.ChallengeActionParameters{}
		if v, ok := challengeMap["challenge_option"].(string); ok && v != "" {
			challengeActionParameters.ChallengeOption = helper.String(v)
		}
		if v, ok := challengeMap["interval"].(string); ok && v != "" {
			challengeActionParameters.Interval = helper.String(v)
		}
		if v, ok := challengeMap["attester_id"].(string); ok && v != "" {
			challengeActionParameters.AttesterId = helper.String(v)
		}
		securityAction.ChallengeActionParameters = &challengeActionParameters
	}
	return &securityAction
}

func expandActionOverrides(overridesList []interface{}) []*teov20220901.BotManagementActionOverrides {
	result := make([]*teov20220901.BotManagementActionOverrides, 0, len(overridesList))
	for _, item := range overridesList {
		overrideMap := item.(map[string]interface{})
		override := teov20220901.BotManagementActionOverrides{}
		if v, ok := overrideMap["ids"]; ok {
			idsList := v.([]interface{})
			for _, id := range idsList {
				if id != nil {
					idStr := id.(string)
					override.Ids = append(override.Ids, &idStr)
				}
			}
		}
		if actionMap, ok := helper.ConvertInterfacesHeadToMap(overrideMap["action"]); ok {
			override.Action = expandSecurityActionForBotManagement(actionMap)
		}
		result = append(result, &override)
	}
	return result
}

func expandBotManagement(botManagementMap map[string]interface{}) *teov20220901.BotManagement {
	botManagement := teov20220901.BotManagement{}
	if v, ok := botManagementMap["enabled"].(string); ok && v != "" {
		botManagement.Enabled = helper.String(v)
	}
	if customRulesMap, ok := helper.ConvertInterfacesHeadToMap(botManagementMap["custom_rules"]); ok {
		customRules := teov20220901.BotManagementCustomRules{}
		if v, ok := customRulesMap["rules"]; ok {
			for _, item := range v.([]interface{}) {
				ruleMap := item.(map[string]interface{})
				rule := teov20220901.BotManagementCustomRule{}
				if v, ok := ruleMap["id"].(string); ok && v != "" {
					rule.Id = helper.String(v)
				}
				if v, ok := ruleMap["name"].(string); ok && v != "" {
					rule.Name = helper.String(v)
				}
				if v, ok := ruleMap["enabled"].(string); ok && v != "" {
					rule.Enabled = helper.String(v)
				}
				if v, ok := ruleMap["priority"].(int); ok {
					rule.Priority = helper.IntInt64(v)
				}
				if v, ok := ruleMap["condition"].(string); ok && v != "" {
					rule.Condition = helper.String(v)
				}
				if v, ok := ruleMap["action"]; ok {
					actionList := v.([]interface{})
					for _, actionItem := range actionList {
						weightedActionMap := actionItem.(map[string]interface{})
						weightedAction := teov20220901.SecurityWeightedAction{}
						if saMap, ok := helper.ConvertInterfacesHeadToMap(weightedActionMap["security_action"]); ok {
							weightedAction.SecurityAction = expandSecurityActionForBotManagement(saMap)
						}
						if v, ok := weightedActionMap["weight"].(int); ok {
							weightedAction.Weight = helper.IntInt64(v)
						}
						rule.Action = append(rule.Action, &weightedAction)
					}
				}
				customRules.Rules = append(customRules.Rules, &rule)
			}
		}
		botManagement.CustomRules = &customRules
	}
	if basicBotSettingsMap, ok := helper.ConvertInterfacesHeadToMap(botManagementMap["basic_bot_settings"]); ok {
		basicBotSettings := teov20220901.BasicBotSettings{}
		if sourceIDCMap, ok := helper.ConvertInterfacesHeadToMap(basicBotSettingsMap["source_idc"]); ok {
			sourceIDC := teov20220901.SourceIDC{}
			if baseActionMap, ok := helper.ConvertInterfacesHeadToMap(sourceIDCMap["base_action"]); ok {
				sourceIDC.BaseAction = expandSecurityActionForBotManagement(baseActionMap)
			}
			if v, ok := sourceIDCMap["action_overrides"]; ok {
				sourceIDC.BotManagementActionOverrides = expandActionOverrides(v.([]interface{}))
			}
			basicBotSettings.SourceIDC = &sourceIDC
		}
		if searchEngineBotsMap, ok := helper.ConvertInterfacesHeadToMap(basicBotSettingsMap["search_engine_bots"]); ok {
			searchEngineBots := teov20220901.SearchEngineBots{}
			if baseActionMap, ok := helper.ConvertInterfacesHeadToMap(searchEngineBotsMap["base_action"]); ok {
				searchEngineBots.BaseAction = expandSecurityActionForBotManagement(baseActionMap)
			}
			if v, ok := searchEngineBotsMap["action_overrides"]; ok {
				searchEngineBots.BotManagementActionOverrides = expandActionOverrides(v.([]interface{}))
			}
			basicBotSettings.SearchEngineBots = &searchEngineBots
		}
		if knownBotCategoriesMap, ok := helper.ConvertInterfacesHeadToMap(basicBotSettingsMap["known_bot_categories"]); ok {
			knownBotCategories := teov20220901.KnownBotCategories{}
			if baseActionMap, ok := helper.ConvertInterfacesHeadToMap(knownBotCategoriesMap["base_action"]); ok {
				knownBotCategories.BaseAction = expandSecurityActionForBotManagement(baseActionMap)
			}
			if v, ok := knownBotCategoriesMap["action_overrides"]; ok {
				knownBotCategories.BotManagementActionOverrides = expandActionOverrides(v.([]interface{}))
			}
			basicBotSettings.KnownBotCategories = &knownBotCategories
		}
		if ipReputationMap, ok := helper.ConvertInterfacesHeadToMap(basicBotSettingsMap["ip_reputation"]); ok {
			ipReputation := teov20220901.IPReputation{}
			if v, ok := ipReputationMap["enabled"].(string); ok && v != "" {
				ipReputation.Enabled = helper.String(v)
			}
			if ipReputationGroupMap, ok := helper.ConvertInterfacesHeadToMap(ipReputationMap["ip_reputation_group"]); ok {
				ipReputationGroup := teov20220901.IPReputationGroup{}
				if baseActionMap, ok := helper.ConvertInterfacesHeadToMap(ipReputationGroupMap["base_action"]); ok {
					ipReputationGroup.BaseAction = expandSecurityActionForBotManagement(baseActionMap)
				}
				if v, ok := ipReputationGroupMap["action_overrides"]; ok {
					ipReputationGroup.BotManagementActionOverrides = expandActionOverrides(v.([]interface{}))
				}
				ipReputation.IPReputationGroup = &ipReputationGroup
			}
			basicBotSettings.IPReputation = &ipReputation
		}
		if botIntelligenceMap, ok := helper.ConvertInterfacesHeadToMap(basicBotSettingsMap["bot_intelligence"]); ok {
			botIntelligence := teov20220901.BotIntelligence{}
			if v, ok := botIntelligenceMap["enabled"].(string); ok && v != "" {
				botIntelligence.Enabled = helper.String(v)
			}
			if botRatingsMap, ok := helper.ConvertInterfacesHeadToMap(botIntelligenceMap["bot_ratings"]); ok {
				botRatings := teov20220901.BotRatings{}
				if actionMap, ok := helper.ConvertInterfacesHeadToMap(botRatingsMap["high_risk_bot_requests_action"]); ok {
					botRatings.HighRiskBotRequestsAction = expandSecurityActionForBotManagement(actionMap)
				}
				if actionMap, ok := helper.ConvertInterfacesHeadToMap(botRatingsMap["likely_bot_requests_action"]); ok {
					botRatings.LikelyBotRequestsAction = expandSecurityActionForBotManagement(actionMap)
				}
				if actionMap, ok := helper.ConvertInterfacesHeadToMap(botRatingsMap["verified_bot_requests_action"]); ok {
					botRatings.VerifiedBotRequestsAction = expandSecurityActionForBotManagement(actionMap)
				}
				if actionMap, ok := helper.ConvertInterfacesHeadToMap(botRatingsMap["human_requests_action"]); ok {
					botRatings.HumanRequestsAction = expandSecurityActionForBotManagement(actionMap)
				}
				botIntelligence.BotRatings = &botRatings
			}
			basicBotSettings.BotIntelligence = &botIntelligence
		}
		botManagement.BasicBotSettings = &basicBotSettings
	}
	if clientAttestationRulesMap, ok := helper.ConvertInterfacesHeadToMap(botManagementMap["client_attestation_rules"]); ok {
		clientAttestationRules := teov20220901.ClientAttestationRules{}
		if v, ok := clientAttestationRulesMap["rules"]; ok {
			for _, item := range v.([]interface{}) {
				ruleMap := item.(map[string]interface{})
				rule := teov20220901.ClientAttestationRule{}
				if v, ok := ruleMap["id"].(string); ok && v != "" {
					rule.Id = helper.String(v)
				}
				if v, ok := ruleMap["name"].(string); ok && v != "" {
					rule.Name = helper.String(v)
				}
				if v, ok := ruleMap["enabled"].(string); ok && v != "" {
					rule.Enabled = helper.String(v)
				}
				if v, ok := ruleMap["priority"].(int); ok {
					rule.Priority = helper.IntUint64(v)
				}
				if v, ok := ruleMap["condition"].(string); ok && v != "" {
					rule.Condition = helper.String(v)
				}
				if v, ok := ruleMap["attester_id"].(string); ok && v != "" {
					rule.AttesterId = helper.String(v)
				}
				if v, ok := ruleMap["device_profiles"]; ok {
					for _, profileItem := range v.([]interface{}) {
						profileMap := profileItem.(map[string]interface{})
						profile := teov20220901.DeviceProfile{}
						if v, ok := profileMap["client_type"].(string); ok && v != "" {
							profile.ClientType = helper.String(v)
						}
						if v, ok := profileMap["high_risk_min_score"].(int); ok {
							profile.HighRiskMinScore = helper.IntUint64(v)
						}
						if actionMap, ok := helper.ConvertInterfacesHeadToMap(profileMap["high_risk_request_action"]); ok {
							profile.HighRiskRequestAction = expandSecurityActionForBotManagement(actionMap)
						}
						if v, ok := profileMap["medium_risk_min_score"].(int); ok {
							profile.MediumRiskMinScore = helper.IntUint64(v)
						}
						if actionMap, ok := helper.ConvertInterfacesHeadToMap(profileMap["medium_risk_request_action"]); ok {
							profile.MediumRiskRequestAction = expandSecurityActionForBotManagement(actionMap)
						}
						rule.DeviceProfiles = append(rule.DeviceProfiles, &profile)
					}
				}
				if actionMap, ok := helper.ConvertInterfacesHeadToMap(ruleMap["invalid_attestation_action"]); ok {
					rule.InvalidAttestationAction = expandSecurityActionForBotManagement(actionMap)
				}
				clientAttestationRules.Rules = append(clientAttestationRules.Rules, &rule)
			}
		}
		botManagement.ClientAttestationRules = &clientAttestationRules
	}
	if browserImpersonationDetectionMap, ok := helper.ConvertInterfacesHeadToMap(botManagementMap["browser_impersonation_detection"]); ok {
		browserImpersonationDetection := teov20220901.BrowserImpersonationDetection{}
		if v, ok := browserImpersonationDetectionMap["rules"]; ok {
			for _, item := range v.([]interface{}) {
				ruleMap := item.(map[string]interface{})
				rule := teov20220901.BrowserImpersonationDetectionRule{}
				if v, ok := ruleMap["id"].(string); ok && v != "" {
					rule.Id = helper.String(v)
				}
				if v, ok := ruleMap["name"].(string); ok && v != "" {
					rule.Name = helper.String(v)
				}
				if v, ok := ruleMap["enabled"].(string); ok && v != "" {
					rule.Enabled = helper.String(v)
				}
				if v, ok := ruleMap["condition"].(string); ok && v != "" {
					rule.Condition = helper.String(v)
				}
				if actionMap, ok := helper.ConvertInterfacesHeadToMap(ruleMap["action"]); ok {
					bidAction := teov20220901.BrowserImpersonationDetectionAction{}
					if bsvMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["bot_session_validation"]); ok {
						bsv := teov20220901.BotSessionValidation{}
						if v, ok := bsvMap["issue_new_bot_session_cookie"].(string); ok && v != "" {
							bsv.IssueNewBotSessionCookie = helper.String(v)
						}
						if triggerConfigMap, ok := helper.ConvertInterfacesHeadToMap(bsvMap["max_new_session_trigger_config"]); ok {
							triggerConfig := teov20220901.MaxNewSessionTriggerConfig{}
							if v, ok := triggerConfigMap["max_new_session_count_interval"].(string); ok && v != "" {
								triggerConfig.MaxNewSessionCountInterval = helper.String(v)
							}
							if v, ok := triggerConfigMap["max_new_session_count_threshold"].(int); ok {
								triggerConfig.MaxNewSessionCountThreshold = helper.IntInt64(v)
							}
							bsv.MaxNewSessionTriggerConfig = &triggerConfig
						}
						if actionMap, ok := helper.ConvertInterfacesHeadToMap(bsvMap["session_expired_action"]); ok {
							bsv.SessionExpiredAction = expandSecurityActionForBotManagement(actionMap)
						}
						if actionMap, ok := helper.ConvertInterfacesHeadToMap(bsvMap["session_invalid_action"]); ok {
							bsv.SessionInvalidAction = expandSecurityActionForBotManagement(actionMap)
						}
						if srcMap, ok := helper.ConvertInterfacesHeadToMap(bsvMap["session_rate_control"]); ok {
							src := teov20220901.SessionRateControl{}
							if v, ok := srcMap["enabled"].(string); ok && v != "" {
								src.Enabled = helper.String(v)
							}
							if actionMap, ok := helper.ConvertInterfacesHeadToMap(srcMap["high_rate_session_action"]); ok {
								src.HighRateSessionAction = expandSecurityActionForBotManagement(actionMap)
							}
							if actionMap, ok := helper.ConvertInterfacesHeadToMap(srcMap["mid_rate_session_action"]); ok {
								src.MidRateSessionAction = expandSecurityActionForBotManagement(actionMap)
							}
							if actionMap, ok := helper.ConvertInterfacesHeadToMap(srcMap["low_rate_session_action"]); ok {
								src.LowRateSessionAction = expandSecurityActionForBotManagement(actionMap)
							}
							bsv.SessionRateControl = &src
						}
						bidAction.BotSessionValidation = &bsv
					}
					if cbdMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["client_behavior_detection"]); ok {
						cbd := teov20220901.ClientBehaviorDetection{}
						if v, ok := cbdMap["crypto_challenge_intensity"].(string); ok && v != "" {
							cbd.CryptoChallengeIntensity = helper.String(v)
						}
						if v, ok := cbdMap["crypto_challenge_delay_before"].(string); ok && v != "" {
							cbd.CryptoChallengeDelayBefore = helper.String(v)
						}
						if v, ok := cbdMap["max_challenge_count_interval"].(string); ok && v != "" {
							cbd.MaxChallengeCountInterval = helper.String(v)
						}
						if v, ok := cbdMap["max_challenge_count_threshold"].(int); ok {
							cbd.MaxChallengeCountThreshold = helper.IntInt64(v)
						}
						if actionMap, ok := helper.ConvertInterfacesHeadToMap(cbdMap["challenge_not_finished_action"]); ok {
							cbd.ChallengeNotFinishedAction = expandSecurityActionForBotManagement(actionMap)
						}
						if actionMap, ok := helper.ConvertInterfacesHeadToMap(cbdMap["challenge_timeout_action"]); ok {
							cbd.ChallengeTimeoutAction = expandSecurityActionForBotManagement(actionMap)
						}
						if actionMap, ok := helper.ConvertInterfacesHeadToMap(cbdMap["bot_client_action"]); ok {
							cbd.BotClientAction = expandSecurityActionForBotManagement(actionMap)
						}
						bidAction.ClientBehaviorDetection = &cbd
					}
					rule.Action = &bidAction
				}
				browserImpersonationDetection.Rules = append(browserImpersonationDetection.Rules, &rule)
			}
		}
		botManagement.BrowserImpersonationDetection = &browserImpersonationDetection
	}
	return &botManagement
}

func resourceTencentCloudTeoSecurityPolicyConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_policy_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		zoneId     string
		entity     string
		host       string
		templateId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("entity"); ok {
		entity = v.(string)
	}

	if v, ok := d.GetOk("host"); ok {
		host = v.(string)
	}

	if v, ok := d.GetOk("template_id"); ok {
		templateId = v.(string)
	}

	if entity == "ZoneDefaultPolicy" && host == "" && templateId == "" {
		d.SetId(strings.Join([]string{zoneId, entity}, tccommon.FILED_SP))
	} else if entity == "Host" && host != "" && templateId == "" {
		d.SetId(strings.Join([]string{zoneId, entity, host}, tccommon.FILED_SP))
	} else if entity == "Template" && host == "" && templateId != "" {
		d.SetId(strings.Join([]string{zoneId, entity, templateId}, tccommon.FILED_SP))
	} else {
		return fmt.Errorf("If `entity` is `ZoneDefaultPolicy`, Please do not set `host` and `template_id`; If `entity` is `Host`, Only support set `host`; If `entity` is `Template`, Only support set `template_id`.")
	}

	return resourceTencentCloudTeoSecurityPolicyConfigUpdate(d, meta)
}

func resourceTencentCloudTeoSecurityPolicyConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_policy_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		zoneId     string
		entity     string
		host       string
		templateId string
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if !(len(idSplit) == 2 || len(idSplit) == 3) {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId = idSplit[0]
	entity = idSplit[1]
	if entity == "ZoneDefaultPolicy" && len(idSplit) == 2 {

	} else if entity == "Host" && len(idSplit) == 3 {
		host = idSplit[2]
	} else if entity == "Template" && len(idSplit) == 3 {
		templateId = idSplit[2]
	} else {
		return fmt.Errorf("`entity` is illegal, %s.", entity)
	}

	respData, err := service.DescribeTeoSecurityPolicyConfigById(ctx, zoneId, entity, host, templateId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `teo_security_policy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("entity", entity)
	_ = d.Set("host", host)
	_ = d.Set("template_id", templateId)

	securityPolicyList := make([]map[string]interface{}, 0, 1)
	securityPolicyMap := map[string]interface{}{}
	if respData.CustomRules != nil {
		customRulesMap := map[string]interface{}{}
		preciseMatchRulesList := make([]map[string]interface{}, 0, len(respData.CustomRules.Rules))
		basicAccessRulesList := make([]map[string]interface{}, 0, len(respData.CustomRules.Rules))
		if respData.CustomRules.Rules != nil {
			for _, rules := range respData.CustomRules.Rules {
				rulesMap := map[string]interface{}{}
				ruleType := ""
				if rules.Name != nil {
					rulesMap["name"] = rules.Name
				}

				if rules.Condition != nil {
					rulesMap["condition"] = rules.Condition
				}

				actionMap := map[string]interface{}{}
				if rules.Action != nil {
					if rules.Action.Name != nil {
						actionMap["name"] = rules.Action.Name
					}

					blockIPActionParametersMap := map[string]interface{}{}
					if rules.Action.BlockIPActionParameters != nil {
						if rules.Action.BlockIPActionParameters.Duration != nil {
							blockIPActionParametersMap["duration"] = rules.Action.BlockIPActionParameters.Duration
						}

						actionMap["block_ip_action_parameters"] = []interface{}{blockIPActionParametersMap}
					}

					returnCustomPageActionParametersMap := map[string]interface{}{}
					if rules.Action.ReturnCustomPageActionParameters != nil {
						if rules.Action.ReturnCustomPageActionParameters.ResponseCode != nil {
							returnCustomPageActionParametersMap["response_code"] = rules.Action.ReturnCustomPageActionParameters.ResponseCode
						}

						if rules.Action.ReturnCustomPageActionParameters.ErrorPageId != nil {
							returnCustomPageActionParametersMap["error_page_id"] = rules.Action.ReturnCustomPageActionParameters.ErrorPageId
						}

						actionMap["return_custom_page_action_parameters"] = []interface{}{returnCustomPageActionParametersMap}
					}

					redirectActionParametersMap := map[string]interface{}{}
					if rules.Action.RedirectActionParameters != nil {
						if rules.Action.RedirectActionParameters.URL != nil {
							redirectActionParametersMap["url"] = rules.Action.RedirectActionParameters.URL
						}

						actionMap["redirect_action_parameters"] = []interface{}{redirectActionParametersMap}
					}

					rulesMap["action"] = []interface{}{actionMap}
				}

				if rules.Enabled != nil {
					rulesMap["enabled"] = rules.Enabled
				}

				if rules.Id != nil {
					rulesMap["id"] = rules.Id
				}

				if rules.RuleType != nil {
					rulesMap["rule_type"] = rules.RuleType
					ruleType = *rules.RuleType
				}

				if rules.Priority != nil {
					rulesMap["priority"] = rules.Priority
				}

				if ruleType == "PreciseMatchRule" {
					preciseMatchRulesList = append(preciseMatchRulesList, rulesMap)
				} else if ruleType == "BasicAccessRule" {
					basicAccessRulesList = append(basicAccessRulesList, rulesMap)
				} else {
					continue
				}
			}

			if len(preciseMatchRulesList) > 0 {
				customRulesMap["precise_match_rules"] = preciseMatchRulesList
			}

			if len(basicAccessRulesList) > 0 {
				customRulesMap["basic_access_rules"] = basicAccessRulesList
			}

			if len(preciseMatchRulesList) > 0 || len(basicAccessRulesList) > 0 {
				securityPolicyMap["custom_rules"] = []interface{}{customRulesMap}
			}
		}
	}

	if respData.ManagedRules != nil {
		managedRulesMap := map[string]interface{}{}
		if respData.ManagedRules.Enabled != nil {
			managedRulesMap["enabled"] = respData.ManagedRules.Enabled
		}

		if respData.ManagedRules.DetectionOnly != nil {
			managedRulesMap["detection_only"] = respData.ManagedRules.DetectionOnly
		}

		if respData.ManagedRules.SemanticAnalysis != nil {
			managedRulesMap["semantic_analysis"] = respData.ManagedRules.SemanticAnalysis
		}

		if respData.ManagedRules.AutoUpdate != nil {
			autoUpdateMap := map[string]interface{}{}
			if respData.ManagedRules.AutoUpdate.AutoUpdateToLatestVersion != nil {
				autoUpdateMap["auto_update_to_latest_version"] = respData.ManagedRules.AutoUpdate.AutoUpdateToLatestVersion
			}

			if respData.ManagedRules.AutoUpdate.RulesetVersion != nil {
				autoUpdateMap["ruleset_version"] = respData.ManagedRules.AutoUpdate.RulesetVersion
			}

			managedRulesMap["auto_update"] = []interface{}{autoUpdateMap}
		}

		if respData.ManagedRules.ManagedRuleGroups != nil {
			managedRuleGroupsList := make([]map[string]interface{}, 0, len(respData.ManagedRules.ManagedRuleGroups))
			for _, managedRuleGroups := range respData.ManagedRules.ManagedRuleGroups {
				managedRuleGroupsMap := map[string]interface{}{}

				if managedRuleGroups.GroupId != nil {
					managedRuleGroupsMap["group_id"] = managedRuleGroups.GroupId
				}

				if managedRuleGroups.SensitivityLevel != nil {
					managedRuleGroupsMap["sensitivity_level"] = managedRuleGroups.SensitivityLevel
				}

				if managedRuleGroups.Action != nil {
					actionMap := map[string]interface{}{}
					if managedRuleGroups.Action.Name != nil {
						actionMap["name"] = managedRuleGroups.Action.Name
					}

					blockIPActionParametersMap := map[string]interface{}{}
					if managedRuleGroups.Action.BlockIPActionParameters != nil {
						if managedRuleGroups.Action.BlockIPActionParameters.Duration != nil {
							blockIPActionParametersMap["duration"] = managedRuleGroups.Action.BlockIPActionParameters.Duration
						}

						actionMap["block_ip_action_parameters"] = []interface{}{blockIPActionParametersMap}
					}

					if managedRuleGroups.Action.ReturnCustomPageActionParameters != nil {
						returnCustomPageActionParametersMap := map[string]interface{}{}
						if managedRuleGroups.Action.ReturnCustomPageActionParameters.ResponseCode != nil {
							returnCustomPageActionParametersMap["response_code"] = managedRuleGroups.Action.ReturnCustomPageActionParameters.ResponseCode
						}

						if managedRuleGroups.Action.ReturnCustomPageActionParameters.ErrorPageId != nil {
							returnCustomPageActionParametersMap["error_page_id"] = managedRuleGroups.Action.ReturnCustomPageActionParameters.ErrorPageId
						}

						actionMap["return_custom_page_action_parameters"] = []interface{}{returnCustomPageActionParametersMap}
					}

					if managedRuleGroups.Action.RedirectActionParameters != nil {
						redirectActionParametersMap := map[string]interface{}{}
						if managedRuleGroups.Action.RedirectActionParameters.URL != nil {
							redirectActionParametersMap["url"] = managedRuleGroups.Action.RedirectActionParameters.URL
						}

						actionMap["redirect_action_parameters"] = []interface{}{redirectActionParametersMap}
					}

					managedRuleGroupsMap["action"] = []interface{}{actionMap}
				}

				if managedRuleGroups.RuleActions != nil {
					ruleActionsList := make([]map[string]interface{}, 0, len(managedRuleGroups.RuleActions))
					for _, ruleActions := range managedRuleGroups.RuleActions {
						ruleActionsMap := map[string]interface{}{}
						if ruleActions.RuleId != nil {
							ruleActionsMap["rule_id"] = ruleActions.RuleId
						}

						if ruleActions.Action != nil {
							actionMap := map[string]interface{}{}
							if ruleActions.Action.Name != nil {
								actionMap["name"] = ruleActions.Action.Name
							}

							if ruleActions.Action.BlockIPActionParameters != nil {
								blockIPActionParametersMap := map[string]interface{}{}
								if ruleActions.Action.BlockIPActionParameters.Duration != nil {
									blockIPActionParametersMap["duration"] = ruleActions.Action.BlockIPActionParameters.Duration
								}

								actionMap["block_ip_action_parameters"] = []interface{}{blockIPActionParametersMap}
							}

							if ruleActions.Action.ReturnCustomPageActionParameters != nil {
								returnCustomPageActionParametersMap := map[string]interface{}{}
								if ruleActions.Action.ReturnCustomPageActionParameters.ResponseCode != nil {
									returnCustomPageActionParametersMap["response_code"] = ruleActions.Action.ReturnCustomPageActionParameters.ResponseCode
								}

								if ruleActions.Action.ReturnCustomPageActionParameters.ErrorPageId != nil {
									returnCustomPageActionParametersMap["error_page_id"] = ruleActions.Action.ReturnCustomPageActionParameters.ErrorPageId
								}

								actionMap["return_custom_page_action_parameters"] = []interface{}{returnCustomPageActionParametersMap}
							}

							if ruleActions.Action.RedirectActionParameters != nil {
								redirectActionParametersMap := map[string]interface{}{}
								if ruleActions.Action.RedirectActionParameters.URL != nil {
									redirectActionParametersMap["url"] = ruleActions.Action.RedirectActionParameters.URL
								}

								actionMap["redirect_action_parameters"] = []interface{}{redirectActionParametersMap}
							}

							ruleActionsMap["action"] = []interface{}{actionMap}
						}

						ruleActionsList = append(ruleActionsList, ruleActionsMap)
					}

					managedRuleGroupsMap["rule_actions"] = ruleActionsList
				}

				if managedRuleGroups.MetaData != nil {
					metaDataMap := map[string]interface{}{}
					if managedRuleGroups.MetaData.GroupDetail != nil {
						metaDataMap["group_detail"] = managedRuleGroups.MetaData.GroupDetail
					}

					if managedRuleGroups.MetaData.GroupName != nil {
						metaDataMap["group_name"] = managedRuleGroups.MetaData.GroupName
					}

					if managedRuleGroups.MetaData.RuleDetails != nil {
						ruleDetailsList := make([]map[string]interface{}, 0, len(managedRuleGroups.MetaData.RuleDetails))
						for _, ruleDetails := range managedRuleGroups.MetaData.RuleDetails {
							ruleDetailsMap := map[string]interface{}{}
							if ruleDetails.RuleId != nil {
								ruleDetailsMap["rule_id"] = ruleDetails.RuleId
							}

							if ruleDetails.RiskLevel != nil {
								ruleDetailsMap["risk_level"] = ruleDetails.RiskLevel
							}

							if ruleDetails.Description != nil {
								ruleDetailsMap["description"] = ruleDetails.Description
							}

							if ruleDetails.Tags != nil {
								ruleDetailsMap["tags"] = ruleDetails.Tags
							}

							if ruleDetails.RuleVersion != nil {
								ruleDetailsMap["rule_version"] = ruleDetails.RuleVersion
							}

							ruleDetailsList = append(ruleDetailsList, ruleDetailsMap)
						}

						metaDataMap["rule_details"] = ruleDetailsList
					}

					managedRuleGroupsMap["meta_data"] = []interface{}{metaDataMap}
				}

				managedRuleGroupsList = append(managedRuleGroupsList, managedRuleGroupsMap)
			}

			managedRulesMap["managed_rule_groups"] = managedRuleGroupsList
		}

		securityPolicyMap["managed_rules"] = []interface{}{managedRulesMap}
	}

	if respData.HttpDDoSProtection != nil {
		httpDDoSProtectionMap := map[string]interface{}{}

		if respData.HttpDDoSProtection.AdaptiveFrequencyControl != nil {
			adaptiveFrequencyControlMap := map[string]interface{}{}

			if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Enabled != nil {
				adaptiveFrequencyControlMap["enabled"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Enabled
			}

			if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Sensitivity != nil {
				adaptiveFrequencyControlMap["sensitivity"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Sensitivity
			}

			if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action != nil {
				actionMap := map[string]interface{}{}

				if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.Name != nil {
					actionMap["name"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.Name
				}

				if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters != nil {
					denyActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.BlockIp != nil {
						denyActionParametersMap["block_ip"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.BlockIp
					}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.BlockIpDuration != nil {
						denyActionParametersMap["block_ip_duration"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.BlockIpDuration
					}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.ReturnCustomPage != nil {
						denyActionParametersMap["return_custom_page"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.ReturnCustomPage
					}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.ResponseCode != nil {
						denyActionParametersMap["response_code"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.ResponseCode
					}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.ErrorPageId != nil {
						denyActionParametersMap["error_page_id"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.ErrorPageId
					}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.Stall != nil {
						denyActionParametersMap["stall"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.DenyActionParameters.Stall
					}

					actionMap["deny_action_parameters"] = []interface{}{denyActionParametersMap}
				}

				if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.RedirectActionParameters != nil {
					redirectActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.RedirectActionParameters.URL != nil {
						redirectActionParametersMap["url"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.RedirectActionParameters.URL
					}

					actionMap["redirect_action_parameters"] = []interface{}{redirectActionParametersMap}
				}

				if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ChallengeActionParameters != nil {
					challengeActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ChallengeActionParameters.ChallengeOption != nil {
						challengeActionParametersMap["challenge_option"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ChallengeActionParameters.ChallengeOption
					}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ChallengeActionParameters.Interval != nil {
						challengeActionParametersMap["interval"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ChallengeActionParameters.Interval
					}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ChallengeActionParameters.AttesterId != nil {
						challengeActionParametersMap["attester_id"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ChallengeActionParameters.AttesterId
					}

					actionMap["challenge_action_parameters"] = []interface{}{challengeActionParametersMap}
				}

				if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.BlockIPActionParameters != nil {
					blockIPActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.BlockIPActionParameters.Duration != nil {
						blockIPActionParametersMap["duration"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.BlockIPActionParameters.Duration
					}

					actionMap["block_i_p_action_parameters"] = []interface{}{blockIPActionParametersMap}
				}

				if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ReturnCustomPageActionParameters != nil {
					returnCustomPageActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ReturnCustomPageActionParameters.ResponseCode != nil {
						returnCustomPageActionParametersMap["response_code"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ReturnCustomPageActionParameters.ResponseCode
					}

					if respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ReturnCustomPageActionParameters.ErrorPageId != nil {
						returnCustomPageActionParametersMap["error_page_id"] = respData.HttpDDoSProtection.AdaptiveFrequencyControl.Action.ReturnCustomPageActionParameters.ErrorPageId
					}

					actionMap["return_custom_page_action_parameters"] = []interface{}{returnCustomPageActionParametersMap}
				}

				adaptiveFrequencyControlMap["action"] = []interface{}{actionMap}
			}

			httpDDoSProtectionMap["adaptive_frequency_control"] = []interface{}{adaptiveFrequencyControlMap}
		}

		if respData.HttpDDoSProtection.ClientFiltering != nil {
			clientFilteringMap := map[string]interface{}{}

			if respData.HttpDDoSProtection.ClientFiltering.Enabled != nil {
				clientFilteringMap["enabled"] = respData.HttpDDoSProtection.ClientFiltering.Enabled
			}

			if respData.HttpDDoSProtection.ClientFiltering.Action != nil {
				actionMap := map[string]interface{}{}

				if respData.HttpDDoSProtection.ClientFiltering.Action.Name != nil {
					actionMap["name"] = respData.HttpDDoSProtection.ClientFiltering.Action.Name
				}

				if respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters != nil {
					denyActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.BlockIp != nil {
						denyActionParametersMap["block_ip"] = respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.BlockIp
					}

					if respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.BlockIpDuration != nil {
						denyActionParametersMap["block_ip_duration"] = respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.BlockIpDuration
					}

					if respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.ReturnCustomPage != nil {
						denyActionParametersMap["return_custom_page"] = respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.ReturnCustomPage
					}

					if respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.ResponseCode != nil {
						denyActionParametersMap["response_code"] = respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.ResponseCode
					}

					if respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.ErrorPageId != nil {
						denyActionParametersMap["error_page_id"] = respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.ErrorPageId
					}

					if respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.Stall != nil {
						denyActionParametersMap["stall"] = respData.HttpDDoSProtection.ClientFiltering.Action.DenyActionParameters.Stall
					}

					actionMap["deny_action_parameters"] = []interface{}{denyActionParametersMap}
				}

				if respData.HttpDDoSProtection.ClientFiltering.Action.RedirectActionParameters != nil {
					redirectActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.ClientFiltering.Action.RedirectActionParameters.URL != nil {
						redirectActionParametersMap["url"] = respData.HttpDDoSProtection.ClientFiltering.Action.RedirectActionParameters.URL
					}

					actionMap["redirect_action_parameters"] = []interface{}{redirectActionParametersMap}
				}

				if respData.HttpDDoSProtection.ClientFiltering.Action.ChallengeActionParameters != nil {
					challengeActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.ClientFiltering.Action.ChallengeActionParameters.ChallengeOption != nil {
						challengeActionParametersMap["challenge_option"] = respData.HttpDDoSProtection.ClientFiltering.Action.ChallengeActionParameters.ChallengeOption
					}

					if respData.HttpDDoSProtection.ClientFiltering.Action.ChallengeActionParameters.Interval != nil {
						challengeActionParametersMap["interval"] = respData.HttpDDoSProtection.ClientFiltering.Action.ChallengeActionParameters.Interval
					}

					if respData.HttpDDoSProtection.ClientFiltering.Action.ChallengeActionParameters.AttesterId != nil {
						challengeActionParametersMap["attester_id"] = respData.HttpDDoSProtection.ClientFiltering.Action.ChallengeActionParameters.AttesterId
					}

					actionMap["challenge_action_parameters"] = []interface{}{challengeActionParametersMap}
				}

				if respData.HttpDDoSProtection.ClientFiltering.Action.BlockIPActionParameters != nil {
					blockIPActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.ClientFiltering.Action.BlockIPActionParameters.Duration != nil {
						blockIPActionParametersMap["duration"] = respData.HttpDDoSProtection.ClientFiltering.Action.BlockIPActionParameters.Duration
					}

					actionMap["block_ip_action_parameters"] = []interface{}{blockIPActionParametersMap}
				}

				if respData.HttpDDoSProtection.ClientFiltering.Action.ReturnCustomPageActionParameters != nil {
					returnCustomPageActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.ClientFiltering.Action.ReturnCustomPageActionParameters.ResponseCode != nil {
						returnCustomPageActionParametersMap["response_code"] = respData.HttpDDoSProtection.ClientFiltering.Action.ReturnCustomPageActionParameters.ResponseCode
					}

					if respData.HttpDDoSProtection.ClientFiltering.Action.ReturnCustomPageActionParameters.ErrorPageId != nil {
						returnCustomPageActionParametersMap["error_page_id"] = respData.HttpDDoSProtection.ClientFiltering.Action.ReturnCustomPageActionParameters.ErrorPageId
					}

					actionMap["return_custom_page_action_parameters"] = []interface{}{returnCustomPageActionParametersMap}
				}

				clientFilteringMap["action"] = []interface{}{actionMap}
			}

			httpDDoSProtectionMap["client_filtering"] = []interface{}{clientFilteringMap}
		}

		if respData.HttpDDoSProtection.BandwidthAbuseDefense != nil {
			bandwidthAbuseDefenseMap := map[string]interface{}{}

			if respData.HttpDDoSProtection.BandwidthAbuseDefense.Enabled != nil {
				bandwidthAbuseDefenseMap["enabled"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Enabled
			}

			if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action != nil {
				actionMap := map[string]interface{}{}

				if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.Name != nil {
					actionMap["name"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.Name
				}

				if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters != nil {
					denyActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.BlockIp != nil {
						denyActionParametersMap["block_ip"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.BlockIp
					}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.BlockIpDuration != nil {
						denyActionParametersMap["block_ip_duration"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.BlockIpDuration
					}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.ReturnCustomPage != nil {
						denyActionParametersMap["return_custom_page"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.ReturnCustomPage
					}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.ResponseCode != nil {
						denyActionParametersMap["response_code"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.ResponseCode
					}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.ErrorPageId != nil {
						denyActionParametersMap["error_page_id"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.ErrorPageId
					}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.Stall != nil {
						denyActionParametersMap["stall"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.DenyActionParameters.Stall
					}

					actionMap["deny_action_parameters"] = []interface{}{denyActionParametersMap}
				}

				if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.RedirectActionParameters != nil {
					redirectActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.RedirectActionParameters.URL != nil {
						redirectActionParametersMap["url"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.RedirectActionParameters.URL
					}

					actionMap["redirect_action_parameters"] = []interface{}{redirectActionParametersMap}
				}

				if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ChallengeActionParameters != nil {
					challengeActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ChallengeActionParameters.ChallengeOption != nil {
						challengeActionParametersMap["challenge_option"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ChallengeActionParameters.ChallengeOption
					}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ChallengeActionParameters.Interval != nil {
						challengeActionParametersMap["interval"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ChallengeActionParameters.Interval
					}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ChallengeActionParameters.AttesterId != nil {
						challengeActionParametersMap["attester_id"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ChallengeActionParameters.AttesterId
					}

					actionMap["challenge_action_parameters"] = []interface{}{challengeActionParametersMap}
				}

				if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.BlockIPActionParameters != nil {
					blockIPActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.BlockIPActionParameters.Duration != nil {
						blockIPActionParametersMap["duration"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.BlockIPActionParameters.Duration
					}

					actionMap["block_ip_action_parameters"] = []interface{}{blockIPActionParametersMap}
				}

				if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ReturnCustomPageActionParameters != nil {
					returnCustomPageActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ReturnCustomPageActionParameters.ResponseCode != nil {
						returnCustomPageActionParametersMap["response_code"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ReturnCustomPageActionParameters.ResponseCode
					}

					if respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ReturnCustomPageActionParameters.ErrorPageId != nil {
						returnCustomPageActionParametersMap["error_page_id"] = respData.HttpDDoSProtection.BandwidthAbuseDefense.Action.ReturnCustomPageActionParameters.ErrorPageId
					}

					actionMap["return_custom_page_action_parameters"] = []interface{}{returnCustomPageActionParametersMap}
				}

				bandwidthAbuseDefenseMap["action"] = []interface{}{actionMap}
			}

			httpDDoSProtectionMap["bandwidth_abuse_defense"] = []interface{}{bandwidthAbuseDefenseMap}
		}

		if respData.HttpDDoSProtection.SlowAttackDefense != nil {
			slowAttackDefenseMap := map[string]interface{}{}

			if respData.HttpDDoSProtection.SlowAttackDefense.Enabled != nil {
				slowAttackDefenseMap["enabled"] = respData.HttpDDoSProtection.SlowAttackDefense.Enabled
			}

			if respData.HttpDDoSProtection.SlowAttackDefense.Action != nil {
				actionMap := map[string]interface{}{}

				if respData.HttpDDoSProtection.SlowAttackDefense.Action.Name != nil {
					actionMap["name"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.Name
				}

				if respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters != nil {
					denyActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.BlockIp != nil {
						denyActionParametersMap["block_ip"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.BlockIp
					}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.BlockIpDuration != nil {
						denyActionParametersMap["block_ip_duration"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.BlockIpDuration
					}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.ReturnCustomPage != nil {
						denyActionParametersMap["return_custom_page"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.ReturnCustomPage
					}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.ResponseCode != nil {
						denyActionParametersMap["response_code"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.ResponseCode
					}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.ErrorPageId != nil {
						denyActionParametersMap["error_page_id"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.ErrorPageId
					}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.Stall != nil {
						denyActionParametersMap["stall"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.DenyActionParameters.Stall
					}

					actionMap["deny_action_parameters"] = []interface{}{denyActionParametersMap}
				}

				if respData.HttpDDoSProtection.SlowAttackDefense.Action.RedirectActionParameters != nil {
					redirectActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.RedirectActionParameters.URL != nil {
						redirectActionParametersMap["url"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.RedirectActionParameters.URL
					}

					actionMap["redirect_action_parameters"] = []interface{}{redirectActionParametersMap}
				}

				if respData.HttpDDoSProtection.SlowAttackDefense.Action.ChallengeActionParameters != nil {
					challengeActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.ChallengeActionParameters.ChallengeOption != nil {
						challengeActionParametersMap["challenge_option"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.ChallengeActionParameters.ChallengeOption
					}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.ChallengeActionParameters.Interval != nil {
						challengeActionParametersMap["interval"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.ChallengeActionParameters.Interval
					}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.ChallengeActionParameters.AttesterId != nil {
						challengeActionParametersMap["attester_id"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.ChallengeActionParameters.AttesterId
					}

					actionMap["challenge_action_parameters"] = []interface{}{challengeActionParametersMap}
				}

				if respData.HttpDDoSProtection.SlowAttackDefense.Action.BlockIPActionParameters != nil {
					blockIPActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.BlockIPActionParameters.Duration != nil {
						blockIPActionParametersMap["duration"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.BlockIPActionParameters.Duration
					}

					actionMap["block_ip_action_parameters"] = []interface{}{blockIPActionParametersMap}
				}

				if respData.HttpDDoSProtection.SlowAttackDefense.Action.ReturnCustomPageActionParameters != nil {
					returnCustomPageActionParametersMap := map[string]interface{}{}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.ReturnCustomPageActionParameters.ResponseCode != nil {
						returnCustomPageActionParametersMap["response_code"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.ReturnCustomPageActionParameters.ResponseCode
					}

					if respData.HttpDDoSProtection.SlowAttackDefense.Action.ReturnCustomPageActionParameters.ErrorPageId != nil {
						returnCustomPageActionParametersMap["error_page_id"] = respData.HttpDDoSProtection.SlowAttackDefense.Action.ReturnCustomPageActionParameters.ErrorPageId
					}

					actionMap["return_custom_page_action_parameters"] = []interface{}{returnCustomPageActionParametersMap}
				}

				slowAttackDefenseMap["action"] = []interface{}{actionMap}
			}

			if respData.HttpDDoSProtection.SlowAttackDefense.MinimalRequestBodyTransferRate != nil {
				minimalRequestBodyTransferRateMap := map[string]interface{}{}

				if respData.HttpDDoSProtection.SlowAttackDefense.MinimalRequestBodyTransferRate.MinimalAvgTransferRateThreshold != nil {
					minimalRequestBodyTransferRateMap["minimal_avg_transfer_rate_threshold"] = respData.HttpDDoSProtection.SlowAttackDefense.MinimalRequestBodyTransferRate.MinimalAvgTransferRateThreshold
				}

				if respData.HttpDDoSProtection.SlowAttackDefense.MinimalRequestBodyTransferRate.CountingPeriod != nil {
					minimalRequestBodyTransferRateMap["counting_period"] = respData.HttpDDoSProtection.SlowAttackDefense.MinimalRequestBodyTransferRate.CountingPeriod
				}

				if respData.HttpDDoSProtection.SlowAttackDefense.MinimalRequestBodyTransferRate.Enabled != nil {
					minimalRequestBodyTransferRateMap["enabled"] = respData.HttpDDoSProtection.SlowAttackDefense.MinimalRequestBodyTransferRate.Enabled
				}

				slowAttackDefenseMap["minimal_request_body_transfer_rate"] = []interface{}{minimalRequestBodyTransferRateMap}
			}

			if respData.HttpDDoSProtection.SlowAttackDefense.RequestBodyTransferTimeout != nil {
				requestBodyTransferTimeoutMap := map[string]interface{}{}

				if respData.HttpDDoSProtection.SlowAttackDefense.RequestBodyTransferTimeout.IdleTimeout != nil {
					requestBodyTransferTimeoutMap["idle_timeout"] = respData.HttpDDoSProtection.SlowAttackDefense.RequestBodyTransferTimeout.IdleTimeout
				}

				if respData.HttpDDoSProtection.SlowAttackDefense.RequestBodyTransferTimeout.Enabled != nil {
					requestBodyTransferTimeoutMap["enabled"] = respData.HttpDDoSProtection.SlowAttackDefense.RequestBodyTransferTimeout.Enabled
				}

				slowAttackDefenseMap["request_body_transfer_timeout"] = []interface{}{requestBodyTransferTimeoutMap}
			}

			httpDDoSProtectionMap["slow_attack_defense"] = []interface{}{slowAttackDefenseMap}
		}

		securityPolicyMap["http_ddos_protection"] = []interface{}{httpDDoSProtectionMap}
	}

	if respData.RateLimitingRules != nil {
		rateLimitingRulesMap := map[string]interface{}{}

		if respData.RateLimitingRules.Rules != nil {
			rulesList := []interface{}{}
			for _, rules := range respData.RateLimitingRules.Rules {
				rulesMap := map[string]interface{}{}

				if rules.Id != nil {
					rulesMap["id"] = rules.Id
				}

				if rules.Name != nil {
					rulesMap["name"] = rules.Name
				}

				if rules.Condition != nil {
					rulesMap["condition"] = rules.Condition
				}

				if rules.CountBy != nil {
					rulesMap["count_by"] = rules.CountBy
				}

				if rules.MaxRequestThreshold != nil {
					rulesMap["max_request_threshold"] = rules.MaxRequestThreshold
				}

				if rules.CountingPeriod != nil {
					rulesMap["counting_period"] = rules.CountingPeriod
				}

				if rules.ActionDuration != nil {
					rulesMap["action_duration"] = rules.ActionDuration
				}

				if rules.Action != nil {
					actionMap := map[string]interface{}{}

					if rules.Action.Name != nil {
						actionMap["name"] = rules.Action.Name
					}

					if rules.Action.DenyActionParameters != nil {
						denyActionParametersMap := map[string]interface{}{}

						if rules.Action.DenyActionParameters.BlockIp != nil {
							denyActionParametersMap["block_ip"] = rules.Action.DenyActionParameters.BlockIp
						}

						if rules.Action.DenyActionParameters.BlockIpDuration != nil {
							denyActionParametersMap["block_ip_duration"] = rules.Action.DenyActionParameters.BlockIpDuration
						}

						if rules.Action.DenyActionParameters.ReturnCustomPage != nil {
							denyActionParametersMap["return_custom_page"] = rules.Action.DenyActionParameters.ReturnCustomPage
						}

						if rules.Action.DenyActionParameters.ResponseCode != nil {
							denyActionParametersMap["response_code"] = rules.Action.DenyActionParameters.ResponseCode
						}

						if rules.Action.DenyActionParameters.ErrorPageId != nil {
							denyActionParametersMap["error_page_id"] = rules.Action.DenyActionParameters.ErrorPageId
						}

						if rules.Action.DenyActionParameters.Stall != nil {
							denyActionParametersMap["stall"] = rules.Action.DenyActionParameters.Stall
						}

						actionMap["deny_action_parameters"] = []interface{}{denyActionParametersMap}
					}

					if rules.Action.RedirectActionParameters != nil {
						redirectActionParametersMap := map[string]interface{}{}

						if rules.Action.RedirectActionParameters.URL != nil {
							redirectActionParametersMap["url"] = rules.Action.RedirectActionParameters.URL
						}

						actionMap["redirect_action_parameters"] = []interface{}{redirectActionParametersMap}
					}

					if rules.Action.ChallengeActionParameters != nil {
						challengeActionParametersMap := map[string]interface{}{}

						if rules.Action.ChallengeActionParameters.ChallengeOption != nil {
							challengeActionParametersMap["challenge_option"] = rules.Action.ChallengeActionParameters.ChallengeOption
						}

						if rules.Action.ChallengeActionParameters.Interval != nil {
							challengeActionParametersMap["interval"] = rules.Action.ChallengeActionParameters.Interval
						}

						if rules.Action.ChallengeActionParameters.AttesterId != nil {
							challengeActionParametersMap["attester_id"] = rules.Action.ChallengeActionParameters.AttesterId
						}

						actionMap["challenge_action_parameters"] = []interface{}{challengeActionParametersMap}
					}

					if rules.Action.BlockIPActionParameters != nil {
						blockIPActionParametersMap := map[string]interface{}{}

						if rules.Action.BlockIPActionParameters.Duration != nil {
							blockIPActionParametersMap["duration"] = rules.Action.BlockIPActionParameters.Duration
						}

						actionMap["block_ip_action_parameters"] = []interface{}{blockIPActionParametersMap}
					}

					if rules.Action.ReturnCustomPageActionParameters != nil {
						returnCustomPageActionParametersMap := map[string]interface{}{}

						if rules.Action.ReturnCustomPageActionParameters.ResponseCode != nil {
							returnCustomPageActionParametersMap["response_code"] = rules.Action.ReturnCustomPageActionParameters.ResponseCode
						}

						if rules.Action.ReturnCustomPageActionParameters.ErrorPageId != nil {
							returnCustomPageActionParametersMap["error_page_id"] = rules.Action.ReturnCustomPageActionParameters.ErrorPageId
						}

						actionMap["return_custom_page_action_parameters"] = []interface{}{returnCustomPageActionParametersMap}
					}

					rulesMap["action"] = []interface{}{actionMap}
				}

				if rules.Priority != nil {
					rulesMap["priority"] = rules.Priority
				}

				if rules.Enabled != nil {
					rulesMap["enabled"] = rules.Enabled
				}

				rulesList = append(rulesList, rulesMap)
			}

			if len(rulesList) > 0 {
				rateLimitingRulesMap["rules"] = rulesList
				securityPolicyMap["rate_limiting_rules"] = []interface{}{rateLimitingRulesMap}
			}
		}
	}

	if respData.ExceptionRules != nil {
		exceptionRulesMap := map[string]interface{}{}

		if respData.ExceptionRules.Rules != nil {
			rulesList := []interface{}{}
			for _, rules := range respData.ExceptionRules.Rules {
				rulesMap := map[string]interface{}{}

				if rules.Id != nil {
					rulesMap["id"] = rules.Id
				}

				if rules.Name != nil {
					rulesMap["name"] = rules.Name
				}

				if rules.Condition != nil {
					rulesMap["condition"] = rules.Condition
				}

				if rules.SkipScope != nil {
					rulesMap["skip_scope"] = rules.SkipScope
				}

				if rules.SkipOption != nil {
					rulesMap["skip_option"] = rules.SkipOption
				}

				if rules.WebSecurityModulesForException != nil {
					rulesMap["web_security_modules_for_exception"] = rules.WebSecurityModulesForException
				}

				if rules.ManagedRulesForException != nil {
					rulesMap["managed_rules_for_exception"] = rules.ManagedRulesForException
				}

				if rules.ManagedRuleGroupsForException != nil {
					rulesMap["managed_rule_groups_for_exception"] = rules.ManagedRuleGroupsForException
				}

				if rules.RequestFieldsForException != nil {
					requestFieldsForExceptionList := []interface{}{}
					for _, requestFieldsForException := range rules.RequestFieldsForException {
						requestFieldsForExceptionMap := map[string]interface{}{}

						if requestFieldsForException.Scope != nil {
							requestFieldsForExceptionMap["scope"] = requestFieldsForException.Scope
						}

						if requestFieldsForException.Condition != nil {
							requestFieldsForExceptionMap["condition"] = requestFieldsForException.Condition
						}

						if requestFieldsForException.TargetField != nil {
							requestFieldsForExceptionMap["target_field"] = requestFieldsForException.TargetField
						}

						requestFieldsForExceptionList = append(requestFieldsForExceptionList, requestFieldsForExceptionMap)
					}

					rulesMap["request_fields_for_exception"] = requestFieldsForExceptionList
				}

				if rules.Enabled != nil {
					rulesMap["enabled"] = rules.Enabled
				}

				rulesList = append(rulesList, rulesMap)
			}

			if len(rulesList) > 0 {
				exceptionRulesMap["rules"] = rulesList
				securityPolicyMap["exception_rules"] = []interface{}{exceptionRulesMap}
			}
		}
	}

	if respData.BotManagement != nil {
		botManagementMap := flattenBotManagement(respData.BotManagement)
		securityPolicyMap["bot_management"] = []interface{}{botManagementMap}
	}

	securityPolicyList = append(securityPolicyList, securityPolicyMap)
	_ = d.Set("security_policy", securityPolicyList)
	return nil
}

func resourceTencentCloudTeoSecurityPolicyConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_policy_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = teov20220901.NewModifySecurityPolicyRequest()
		zoneId     string
		entity     string
		host       string
		templateId string
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if !(len(idSplit) == 2 || len(idSplit) == 3) {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId = idSplit[0]
	entity = idSplit[1]
	if entity == "ZoneDefaultPolicy" && len(idSplit) == 2 {

	} else if entity == "Host" && len(idSplit) == 3 {
		host = idSplit[2]
	} else if entity == "Template" && len(idSplit) == 3 {
		templateId = idSplit[2]
	} else {
		return fmt.Errorf("`entity` is illegal, %s.", entity)
	}

	request.ZoneId = &zoneId
	request.Entity = &entity
	request.TemplateId = &templateId
	request.Host = &host
	request.SecurityConfig = &teov20220901.SecurityConfig{
		RateLimitConfig: &teov20220901.RateLimitConfig{
			RateLimitUserRules: []*teov20220901.RateLimitUserRule{},
			Switch:             helper.String("on"),
		},
	}

	if securityPolicyMap, ok := helper.InterfacesHeadMap(d, "security_policy"); ok {
		securityPolicy := teov20220901.SecurityPolicy{}
		if customRulesMap, ok := helper.ConvertInterfacesHeadToMap(securityPolicyMap["custom_rules"]); ok {
			customRules := teov20220901.CustomRules{}
			if v, ok := customRulesMap["rules"]; ok {
				if len(v.([]interface{})) > 0 {
					return fmt.Errorf("`rules` has been deprecated from version 1.81.184. Please use `precise_match_rules` or `basic_access_rules` instead.")
				}
			}

			if v, ok := customRulesMap["precise_match_rules"]; ok {
				for _, item := range v.([]interface{}) {
					rulesMap := item.(map[string]interface{})
					customRule := teov20220901.CustomRule{}
					if v, ok := rulesMap["name"].(string); ok && v != "" {
						customRule.Name = helper.String(v)
					}

					if v, ok := rulesMap["condition"].(string); ok && v != "" {
						customRule.Condition = helper.String(v)
					}

					if actionMap, ok := helper.ConvertInterfacesHeadToMap(rulesMap["action"]); ok {
						securityAction := teov20220901.SecurityAction{}
						if v, ok := actionMap["name"].(string); ok && v != "" {
							securityAction.Name = helper.String(v)
						}

						if blockIPActionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["block_ip_action_parameters"]); ok {
							blockIPActionParameters := teov20220901.BlockIPActionParameters{}
							if v, ok := blockIPActionParametersMap["duration"].(string); ok && v != "" {
								blockIPActionParameters.Duration = helper.String(v)
							}

							securityAction.BlockIPActionParameters = &blockIPActionParameters
						}

						if returnCustomPageActionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["return_custom_page_action_parameters"]); ok {
							returnCustomPageActionParameters := teov20220901.ReturnCustomPageActionParameters{}
							if v, ok := returnCustomPageActionParametersMap["response_code"].(string); ok && v != "" {
								returnCustomPageActionParameters.ResponseCode = helper.String(v)
							}

							if v, ok := returnCustomPageActionParametersMap["error_page_id"].(string); ok && v != "" {
								returnCustomPageActionParameters.ErrorPageId = helper.String(v)
							}

							securityAction.ReturnCustomPageActionParameters = &returnCustomPageActionParameters
						}

						if redirectActionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["redirect_action_parameters"]); ok {
							redirectActionParameters := teov20220901.RedirectActionParameters{}
							if v, ok := redirectActionParametersMap["url"].(string); ok && v != "" {
								redirectActionParameters.URL = helper.String(v)
							}

							securityAction.RedirectActionParameters = &redirectActionParameters
						}

						customRule.Action = &securityAction
					}

					if v, ok := rulesMap["enabled"].(string); ok && v != "" {
						customRule.Enabled = helper.String(v)
					}

					if v, ok := rulesMap["id"].(string); ok && v != "" {
						customRule.Id = helper.String(v)
					}

					customRule.RuleType = helper.String("PreciseMatchRule")

					if v, ok := rulesMap["priority"].(int); ok {
						customRule.Priority = helper.IntInt64(v)
					}

					customRules.Rules = append(customRules.Rules, &customRule)
				}
			}

			if v, ok := customRulesMap["basic_access_rules"]; ok {
				for _, item := range v.([]interface{}) {
					rulesMap := item.(map[string]interface{})
					customRule := teov20220901.CustomRule{}
					if v, ok := rulesMap["name"].(string); ok && v != "" {
						customRule.Name = helper.String(v)
					}

					if v, ok := rulesMap["condition"].(string); ok && v != "" {
						customRule.Condition = helper.String(v)
					}

					if actionMap, ok := helper.ConvertInterfacesHeadToMap(rulesMap["action"]); ok {
						securityAction := teov20220901.SecurityAction{}
						if v, ok := actionMap["name"].(string); ok && v != "" {
							securityAction.Name = helper.String(v)
						}

						if blockIPActionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["block_ip_action_parameters"]); ok {
							blockIPActionParameters := teov20220901.BlockIPActionParameters{}
							if v, ok := blockIPActionParametersMap["duration"].(string); ok && v != "" {
								blockIPActionParameters.Duration = helper.String(v)
							}

							securityAction.BlockIPActionParameters = &blockIPActionParameters
						}

						if returnCustomPageActionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["return_custom_page_action_parameters"]); ok {
							returnCustomPageActionParameters := teov20220901.ReturnCustomPageActionParameters{}
							if v, ok := returnCustomPageActionParametersMap["response_code"].(string); ok && v != "" {
								returnCustomPageActionParameters.ResponseCode = helper.String(v)
							}

							if v, ok := returnCustomPageActionParametersMap["error_page_id"].(string); ok && v != "" {
								returnCustomPageActionParameters.ErrorPageId = helper.String(v)
							}

							securityAction.ReturnCustomPageActionParameters = &returnCustomPageActionParameters
						}

						if redirectActionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["redirect_action_parameters"]); ok {
							redirectActionParameters := teov20220901.RedirectActionParameters{}
							if v, ok := redirectActionParametersMap["url"].(string); ok && v != "" {
								redirectActionParameters.URL = helper.String(v)
							}

							securityAction.RedirectActionParameters = &redirectActionParameters
						}

						customRule.Action = &securityAction
					}

					if v, ok := rulesMap["enabled"].(string); ok && v != "" {
						customRule.Enabled = helper.String(v)
					}

					if v, ok := rulesMap["id"].(string); ok && v != "" {
						customRule.Id = helper.String(v)
					}

					customRule.RuleType = helper.String("BasicAccessRule")

					if v, ok := rulesMap["priority"].(int); ok {
						customRule.Priority = helper.IntInt64(v)
					}

					customRules.Rules = append(customRules.Rules, &customRule)
				}
			}

			securityPolicy.CustomRules = &customRules
		} else {
			securityPolicy.CustomRules = &teov20220901.CustomRules{
				Rules: []*teov20220901.CustomRule{},
			}
		}

		if managedRulesMap, ok := helper.ConvertInterfacesHeadToMap(securityPolicyMap["managed_rules"]); ok {
			managedRules := teov20220901.ManagedRules{}
			if v, ok := managedRulesMap["enabled"].(string); ok && v != "" {
				managedRules.Enabled = helper.String(v)
			}

			if v, ok := managedRulesMap["detection_only"].(string); ok && v != "" {
				managedRules.DetectionOnly = helper.String(v)
			}

			if v, ok := managedRulesMap["semantic_analysis"].(string); ok && v != "" {
				managedRules.SemanticAnalysis = helper.String(v)
			}

			if autoUpdateMap, ok := helper.ConvertInterfacesHeadToMap(managedRulesMap["auto_update"]); ok {
				managedRuleAutoUpdate := teov20220901.ManagedRuleAutoUpdate{}
				if v, ok := autoUpdateMap["auto_update_to_latest_version"].(string); ok && v != "" {
					managedRuleAutoUpdate.AutoUpdateToLatestVersion = helper.String(v)
				}

				managedRules.AutoUpdate = &managedRuleAutoUpdate
			}

			if v, ok := managedRulesMap["managed_rule_groups"]; ok {
				for _, item := range v.(*schema.Set).List() {
					managedRuleGroupsMap := item.(map[string]interface{})
					managedRuleGroup := teov20220901.ManagedRuleGroup{}
					if v, ok := managedRuleGroupsMap["group_id"].(string); ok && v != "" {
						managedRuleGroup.GroupId = helper.String(v)
					}

					if v, ok := managedRuleGroupsMap["sensitivity_level"].(string); ok && v != "" {
						managedRuleGroup.SensitivityLevel = helper.String(v)
					}

					if actionMap, ok := helper.ConvertInterfacesHeadToMap(managedRuleGroupsMap["action"]); ok {
						securityAction2 := teov20220901.SecurityAction{}
						if v, ok := actionMap["name"].(string); ok && v != "" {
							securityAction2.Name = helper.String(v)
						}

						if blockIPActionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["block_ip_action_parameters"]); ok {
							blockIPActionParameters2 := teov20220901.BlockIPActionParameters{}
							if v, ok := blockIPActionParametersMap["duration"].(string); ok && v != "" {
								blockIPActionParameters2.Duration = helper.String(v)
							}

							securityAction2.BlockIPActionParameters = &blockIPActionParameters2
						}

						if returnCustomPageActionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["return_custom_page_action_parameters"]); ok {
							returnCustomPageActionParameters2 := teov20220901.ReturnCustomPageActionParameters{}
							if v, ok := returnCustomPageActionParametersMap["response_code"].(string); ok && v != "" {
								returnCustomPageActionParameters2.ResponseCode = helper.String(v)
							}

							if v, ok := returnCustomPageActionParametersMap["error_page_id"].(string); ok && v != "" {
								returnCustomPageActionParameters2.ErrorPageId = helper.String(v)
							}

							securityAction2.ReturnCustomPageActionParameters = &returnCustomPageActionParameters2
						}

						if redirectActionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["redirect_action_parameters"]); ok {
							redirectActionParameters2 := teov20220901.RedirectActionParameters{}
							if v, ok := redirectActionParametersMap["url"].(string); ok && v != "" {
								redirectActionParameters2.URL = helper.String(v)
							}

							securityAction2.RedirectActionParameters = &redirectActionParameters2
						}

						managedRuleGroup.Action = &securityAction2
					}

					if v, ok := managedRuleGroupsMap["rule_actions"]; ok {
						for _, item := range v.([]interface{}) {
							ruleActionsMap := item.(map[string]interface{})
							managedRuleAction := teov20220901.ManagedRuleAction{}
							if v, ok := ruleActionsMap["rule_id"].(string); ok && v != "" {
								managedRuleAction.RuleId = helper.String(v)
							}

							if actionMap, ok := helper.ConvertInterfacesHeadToMap(ruleActionsMap["action"]); ok {
								securityAction3 := teov20220901.SecurityAction{}
								if v, ok := actionMap["name"].(string); ok && v != "" {
									securityAction3.Name = helper.String(v)
								}

								if blockIPActionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["block_ip_action_parameters"]); ok {
									blockIPActionParameters3 := teov20220901.BlockIPActionParameters{}
									if v, ok := blockIPActionParametersMap["duration"].(string); ok && v != "" {
										blockIPActionParameters3.Duration = helper.String(v)
									}

									securityAction3.BlockIPActionParameters = &blockIPActionParameters3
								}

								if returnCustomPageActionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["return_custom_page_action_parameters"]); ok {
									returnCustomPageActionParameters3 := teov20220901.ReturnCustomPageActionParameters{}
									if v, ok := returnCustomPageActionParametersMap["response_code"].(string); ok && v != "" {
										returnCustomPageActionParameters3.ResponseCode = helper.String(v)
									}

									if v, ok := returnCustomPageActionParametersMap["error_page_id"].(string); ok && v != "" {
										returnCustomPageActionParameters3.ErrorPageId = helper.String(v)
									}

									securityAction3.ReturnCustomPageActionParameters = &returnCustomPageActionParameters3
								}

								if redirectActionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionMap["redirect_action_parameters"]); ok {
									redirectActionParameters3 := teov20220901.RedirectActionParameters{}
									if v, ok := redirectActionParametersMap["url"].(string); ok && v != "" {
										redirectActionParameters3.URL = helper.String(v)
									}

									securityAction3.RedirectActionParameters = &redirectActionParameters3
								}

								managedRuleAction.Action = &securityAction3
							}

							managedRuleGroup.RuleActions = append(managedRuleGroup.RuleActions, &managedRuleAction)
						}
					}

					managedRules.ManagedRuleGroups = append(managedRules.ManagedRuleGroups, &managedRuleGroup)
				}
			}

			securityPolicy.ManagedRules = &managedRules
		}

		if httpDDoSProtectionMap, ok := helper.ConvertInterfacesHeadToMap(securityPolicyMap["http_ddos_protection"]); ok {
			httpDDoSProtection := teov20220901.HttpDDoSProtection{}
			if adaptiveFrequencyControlMap, ok := helper.InterfaceToMap(httpDDoSProtectionMap, "adaptive_frequency_control"); ok {
				adaptiveFrequencyControl := teov20220901.AdaptiveFrequencyControl{}
				if v, ok := adaptiveFrequencyControlMap["enabled"]; ok {
					adaptiveFrequencyControl.Enabled = helper.String(v.(string))
				}

				if v, ok := adaptiveFrequencyControlMap["sensitivity"]; ok {
					adaptiveFrequencyControl.Sensitivity = helper.String(v.(string))
				}

				if actionMap, ok := helper.InterfaceToMap(adaptiveFrequencyControlMap, "action"); ok {
					securityAction := teov20220901.SecurityAction{}
					if v, ok := actionMap["name"]; ok {
						securityAction.Name = helper.String(v.(string))
					}

					if denyActionParametersMap, ok := helper.InterfaceToMap(actionMap, "deny_action_parameters"); ok {
						denyActionParameters := teov20220901.DenyActionParameters{}
						if v, ok := denyActionParametersMap["block_ip"]; ok {
							denyActionParameters.BlockIp = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["block_ip_duration"]; ok {
							denyActionParameters.BlockIpDuration = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["return_custom_page"]; ok {
							denyActionParameters.ReturnCustomPage = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["response_code"]; ok {
							denyActionParameters.ResponseCode = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["error_page_id"]; ok {
							denyActionParameters.ErrorPageId = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["stall"]; ok {
							denyActionParameters.Stall = helper.String(v.(string))
						}

						securityAction.DenyActionParameters = &denyActionParameters
					}

					if redirectActionParametersMap, ok := helper.InterfaceToMap(actionMap, "redirect_action_parameters"); ok {
						redirectActionParameters := teov20220901.RedirectActionParameters{}
						if v, ok := redirectActionParametersMap["url"]; ok {
							redirectActionParameters.URL = helper.String(v.(string))
						}

						securityAction.RedirectActionParameters = &redirectActionParameters
					}

					if challengeActionParametersMap, ok := helper.InterfaceToMap(actionMap, "challenge_action_parameters"); ok {
						challengeActionParameters := teov20220901.ChallengeActionParameters{}
						if v, ok := challengeActionParametersMap["challenge_option"]; ok {
							challengeActionParameters.ChallengeOption = helper.String(v.(string))
						}

						if v, ok := challengeActionParametersMap["interval"]; ok {
							challengeActionParameters.Interval = helper.String(v.(string))
						}

						if v, ok := challengeActionParametersMap["attester_id"]; ok {
							challengeActionParameters.AttesterId = helper.String(v.(string))
						}

						securityAction.ChallengeActionParameters = &challengeActionParameters
					}

					if blockIPActionParametersMap, ok := helper.InterfaceToMap(actionMap, "block_ip_action_parameters"); ok {
						blockIPActionParameters := teov20220901.BlockIPActionParameters{}
						if v, ok := blockIPActionParametersMap["duration"]; ok {
							blockIPActionParameters.Duration = helper.String(v.(string))
						}

						securityAction.BlockIPActionParameters = &blockIPActionParameters
					}

					if returnCustomPageActionParametersMap, ok := helper.InterfaceToMap(actionMap, "return_custom_page_action_parameters"); ok {
						returnCustomPageActionParameters := teov20220901.ReturnCustomPageActionParameters{}
						if v, ok := returnCustomPageActionParametersMap["response_code"]; ok {
							returnCustomPageActionParameters.ResponseCode = helper.String(v.(string))
						}

						if v, ok := returnCustomPageActionParametersMap["error_page_id"]; ok {
							returnCustomPageActionParameters.ErrorPageId = helper.String(v.(string))
						}

						securityAction.ReturnCustomPageActionParameters = &returnCustomPageActionParameters
					}

					adaptiveFrequencyControl.Action = &securityAction
				}

				httpDDoSProtection.AdaptiveFrequencyControl = &adaptiveFrequencyControl
			}

			if clientFilteringMap, ok := helper.InterfaceToMap(httpDDoSProtectionMap, "client_filtering"); ok {
				clientFiltering := teov20220901.ClientFiltering{}
				if v, ok := clientFilteringMap["enabled"]; ok {
					clientFiltering.Enabled = helper.String(v.(string))
				}

				if actionMap, ok := helper.InterfaceToMap(clientFilteringMap, "action"); ok {
					securityAction := teov20220901.SecurityAction{}
					if v, ok := actionMap["name"]; ok {
						securityAction.Name = helper.String(v.(string))
					}

					if denyActionParametersMap, ok := helper.InterfaceToMap(actionMap, "deny_action_parameters"); ok {
						denyActionParameters := teov20220901.DenyActionParameters{}
						if v, ok := denyActionParametersMap["block_ip"]; ok {
							denyActionParameters.BlockIp = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["block_ip_duration"]; ok {
							denyActionParameters.BlockIpDuration = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["return_custom_page"]; ok {
							denyActionParameters.ReturnCustomPage = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["response_code"]; ok {
							denyActionParameters.ResponseCode = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["error_page_id"]; ok {
							denyActionParameters.ErrorPageId = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["stall"]; ok {
							denyActionParameters.Stall = helper.String(v.(string))
						}

						securityAction.DenyActionParameters = &denyActionParameters
					}

					if redirectActionParametersMap, ok := helper.InterfaceToMap(actionMap, "redirect_action_parameters"); ok {
						redirectActionParameters := teov20220901.RedirectActionParameters{}
						if v, ok := redirectActionParametersMap["url"]; ok {
							redirectActionParameters.URL = helper.String(v.(string))
						}

						securityAction.RedirectActionParameters = &redirectActionParameters
					}

					if challengeActionParametersMap, ok := helper.InterfaceToMap(actionMap, "challenge_action_parameters"); ok {
						challengeActionParameters := teov20220901.ChallengeActionParameters{}
						if v, ok := challengeActionParametersMap["challenge_option"]; ok {
							challengeActionParameters.ChallengeOption = helper.String(v.(string))
						}

						if v, ok := challengeActionParametersMap["interval"]; ok {
							challengeActionParameters.Interval = helper.String(v.(string))
						}

						if v, ok := challengeActionParametersMap["attester_id"]; ok {
							challengeActionParameters.AttesterId = helper.String(v.(string))
						}

						securityAction.ChallengeActionParameters = &challengeActionParameters
					}

					if blockIPActionParametersMap, ok := helper.InterfaceToMap(actionMap, "block_ip_action_parameters"); ok {
						blockIPActionParameters := teov20220901.BlockIPActionParameters{}
						if v, ok := blockIPActionParametersMap["duration"]; ok {
							blockIPActionParameters.Duration = helper.String(v.(string))
						}

						securityAction.BlockIPActionParameters = &blockIPActionParameters
					}

					if returnCustomPageActionParametersMap, ok := helper.InterfaceToMap(actionMap, "return_custom_page_action_parameters"); ok {
						returnCustomPageActionParameters := teov20220901.ReturnCustomPageActionParameters{}
						if v, ok := returnCustomPageActionParametersMap["response_code"]; ok {
							returnCustomPageActionParameters.ResponseCode = helper.String(v.(string))
						}

						if v, ok := returnCustomPageActionParametersMap["error_page_id"]; ok {
							returnCustomPageActionParameters.ErrorPageId = helper.String(v.(string))
						}

						securityAction.ReturnCustomPageActionParameters = &returnCustomPageActionParameters
					}

					clientFiltering.Action = &securityAction
				}

				httpDDoSProtection.ClientFiltering = &clientFiltering
			}

			if bandwidthAbuseDefenseMap, ok := helper.InterfaceToMap(httpDDoSProtectionMap, "bandwidth_abuse_defense"); ok {
				bandwidthAbuseDefense := teov20220901.BandwidthAbuseDefense{}
				if v, ok := bandwidthAbuseDefenseMap["enabled"]; ok {
					bandwidthAbuseDefense.Enabled = helper.String(v.(string))
				}

				if actionMap, ok := helper.InterfaceToMap(bandwidthAbuseDefenseMap, "action"); ok {
					securityAction := teov20220901.SecurityAction{}
					if v, ok := actionMap["name"]; ok {
						securityAction.Name = helper.String(v.(string))
					}

					if denyActionParametersMap, ok := helper.InterfaceToMap(actionMap, "deny_action_parameters"); ok {
						denyActionParameters := teov20220901.DenyActionParameters{}
						if v, ok := denyActionParametersMap["block_ip"]; ok {
							denyActionParameters.BlockIp = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["block_ip_duration"]; ok {
							denyActionParameters.BlockIpDuration = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["return_custom_page"]; ok {
							denyActionParameters.ReturnCustomPage = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["response_code"]; ok {
							denyActionParameters.ResponseCode = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["error_page_id"]; ok {
							denyActionParameters.ErrorPageId = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["stall"]; ok {
							denyActionParameters.Stall = helper.String(v.(string))
						}

						securityAction.DenyActionParameters = &denyActionParameters
					}

					if redirectActionParametersMap, ok := helper.InterfaceToMap(actionMap, "redirect_action_parameters"); ok {
						redirectActionParameters := teov20220901.RedirectActionParameters{}
						if v, ok := redirectActionParametersMap["url"]; ok {
							redirectActionParameters.URL = helper.String(v.(string))
						}

						securityAction.RedirectActionParameters = &redirectActionParameters
					}

					if challengeActionParametersMap, ok := helper.InterfaceToMap(actionMap, "challenge_action_parameters"); ok {
						challengeActionParameters := teov20220901.ChallengeActionParameters{}
						if v, ok := challengeActionParametersMap["challenge_option"]; ok {
							challengeActionParameters.ChallengeOption = helper.String(v.(string))
						}

						if v, ok := challengeActionParametersMap["interval"]; ok {
							challengeActionParameters.Interval = helper.String(v.(string))
						}

						if v, ok := challengeActionParametersMap["attester_id"]; ok {
							challengeActionParameters.AttesterId = helper.String(v.(string))
						}

						securityAction.ChallengeActionParameters = &challengeActionParameters
					}

					if blockIPActionParametersMap, ok := helper.InterfaceToMap(actionMap, "block_ip_action_parameters"); ok {
						blockIPActionParameters := teov20220901.BlockIPActionParameters{}
						if v, ok := blockIPActionParametersMap["duration"]; ok {
							blockIPActionParameters.Duration = helper.String(v.(string))
						}

						securityAction.BlockIPActionParameters = &blockIPActionParameters
					}

					if returnCustomPageActionParametersMap, ok := helper.InterfaceToMap(actionMap, "return_custom_page_action_parameters"); ok {
						returnCustomPageActionParameters := teov20220901.ReturnCustomPageActionParameters{}
						if v, ok := returnCustomPageActionParametersMap["response_code"]; ok {
							returnCustomPageActionParameters.ResponseCode = helper.String(v.(string))
						}

						if v, ok := returnCustomPageActionParametersMap["error_page_id"]; ok {
							returnCustomPageActionParameters.ErrorPageId = helper.String(v.(string))
						}

						securityAction.ReturnCustomPageActionParameters = &returnCustomPageActionParameters
					}
					bandwidthAbuseDefense.Action = &securityAction
				}

				httpDDoSProtection.BandwidthAbuseDefense = &bandwidthAbuseDefense
			}

			if slowAttackDefenseMap, ok := helper.InterfaceToMap(httpDDoSProtectionMap, "slow_attack_defense"); ok {
				slowAttackDefense := teov20220901.SlowAttackDefense{}
				if v, ok := slowAttackDefenseMap["enabled"]; ok {
					slowAttackDefense.Enabled = helper.String(v.(string))
				}

				if actionMap, ok := helper.InterfaceToMap(slowAttackDefenseMap, "action"); ok {
					securityAction := teov20220901.SecurityAction{}
					if v, ok := actionMap["name"]; ok {
						securityAction.Name = helper.String(v.(string))
					}

					if denyActionParametersMap, ok := helper.InterfaceToMap(actionMap, "deny_action_parameters"); ok {
						denyActionParameters := teov20220901.DenyActionParameters{}
						if v, ok := denyActionParametersMap["block_ip"]; ok {
							denyActionParameters.BlockIp = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["block_ip_duration"]; ok {
							denyActionParameters.BlockIpDuration = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["return_custom_page"]; ok {
							denyActionParameters.ReturnCustomPage = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["response_code"]; ok {
							denyActionParameters.ResponseCode = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["error_page_id"]; ok {
							denyActionParameters.ErrorPageId = helper.String(v.(string))
						}

						if v, ok := denyActionParametersMap["stall"]; ok {
							denyActionParameters.Stall = helper.String(v.(string))
						}

						securityAction.DenyActionParameters = &denyActionParameters
					}

					if redirectActionParametersMap, ok := helper.InterfaceToMap(actionMap, "redirect_action_parameters"); ok {
						redirectActionParameters := teov20220901.RedirectActionParameters{}
						if v, ok := redirectActionParametersMap["url"]; ok {
							redirectActionParameters.URL = helper.String(v.(string))
						}

						securityAction.RedirectActionParameters = &redirectActionParameters
					}

					if challengeActionParametersMap, ok := helper.InterfaceToMap(actionMap, "challenge_action_parameters"); ok {
						challengeActionParameters := teov20220901.ChallengeActionParameters{}
						if v, ok := challengeActionParametersMap["challenge_option"]; ok {
							challengeActionParameters.ChallengeOption = helper.String(v.(string))
						}

						if v, ok := challengeActionParametersMap["interval"]; ok {
							challengeActionParameters.Interval = helper.String(v.(string))
						}

						if v, ok := challengeActionParametersMap["attester_id"]; ok {
							challengeActionParameters.AttesterId = helper.String(v.(string))
						}

						securityAction.ChallengeActionParameters = &challengeActionParameters
					}

					if blockIPActionParametersMap, ok := helper.InterfaceToMap(actionMap, "block_ip_action_parameters"); ok {
						blockIPActionParameters := teov20220901.BlockIPActionParameters{}
						if v, ok := blockIPActionParametersMap["duration"]; ok {
							blockIPActionParameters.Duration = helper.String(v.(string))
						}

						securityAction.BlockIPActionParameters = &blockIPActionParameters
					}

					if returnCustomPageActionParametersMap, ok := helper.InterfaceToMap(actionMap, "return_custom_page_action_parameters"); ok {
						returnCustomPageActionParameters := teov20220901.ReturnCustomPageActionParameters{}
						if v, ok := returnCustomPageActionParametersMap["response_code"]; ok {
							returnCustomPageActionParameters.ResponseCode = helper.String(v.(string))
						}

						if v, ok := returnCustomPageActionParametersMap["error_page_id"]; ok {
							returnCustomPageActionParameters.ErrorPageId = helper.String(v.(string))
						}

						securityAction.ReturnCustomPageActionParameters = &returnCustomPageActionParameters
					}

					slowAttackDefense.Action = &securityAction
				}

				if minimalRequestBodyTransferRateMap, ok := helper.InterfaceToMap(slowAttackDefenseMap, "minimal_request_body_transfer_rate"); ok {
					minimalRequestBodyTransferRate := teov20220901.MinimalRequestBodyTransferRate{}
					if v, ok := minimalRequestBodyTransferRateMap["minimal_avg_transfer_rate_threshold"]; ok {
						minimalRequestBodyTransferRate.MinimalAvgTransferRateThreshold = helper.String(v.(string))
					}

					if v, ok := minimalRequestBodyTransferRateMap["counting_period"]; ok {
						minimalRequestBodyTransferRate.CountingPeriod = helper.String(v.(string))
					}

					if v, ok := minimalRequestBodyTransferRateMap["enabled"]; ok {
						minimalRequestBodyTransferRate.Enabled = helper.String(v.(string))
					}

					slowAttackDefense.MinimalRequestBodyTransferRate = &minimalRequestBodyTransferRate
				}

				if requestBodyTransferTimeoutMap, ok := helper.InterfaceToMap(slowAttackDefenseMap, "request_body_transfer_timeout"); ok {
					requestBodyTransferTimeout := teov20220901.RequestBodyTransferTimeout{}
					if v, ok := requestBodyTransferTimeoutMap["idle_timeout"]; ok {
						requestBodyTransferTimeout.IdleTimeout = helper.String(v.(string))
					}

					if v, ok := requestBodyTransferTimeoutMap["enabled"]; ok {
						requestBodyTransferTimeout.Enabled = helper.String(v.(string))
					}

					slowAttackDefense.RequestBodyTransferTimeout = &requestBodyTransferTimeout
				}

				httpDDoSProtection.SlowAttackDefense = &slowAttackDefense
			}

			securityPolicy.HttpDDoSProtection = &httpDDoSProtection
		}

		if rateLimitingRulesMap, ok := helper.ConvertInterfacesHeadToMap(securityPolicyMap["rate_limiting_rules"]); ok {
			rateLimitingRules := teov20220901.RateLimitingRules{}
			if v, ok := rateLimitingRulesMap["rules"]; ok {
				for _, item := range v.([]interface{}) {
					rulesMap := item.(map[string]interface{})
					rateLimitingRule := teov20220901.RateLimitingRule{}
					if v, ok := rulesMap["id"]; ok {
						rateLimitingRule.Id = helper.String(v.(string))
					}

					if v, ok := rulesMap["name"]; ok {
						rateLimitingRule.Name = helper.String(v.(string))
					}

					if v, ok := rulesMap["condition"]; ok {
						rateLimitingRule.Condition = helper.String(v.(string))
					}

					if v, ok := rulesMap["count_by"]; ok {
						countBySet := v.(*schema.Set).List()
						for i := range countBySet {
							if countBySet[i] != nil {
								countBy := countBySet[i].(string)
								rateLimitingRule.CountBy = append(rateLimitingRule.CountBy, &countBy)
							}
						}
					}

					if v, ok := rulesMap["max_request_threshold"]; ok {
						rateLimitingRule.MaxRequestThreshold = helper.IntInt64(v.(int))
					}

					if v, ok := rulesMap["counting_period"]; ok {
						rateLimitingRule.CountingPeriod = helper.String(v.(string))
					}

					if v, ok := rulesMap["action_duration"]; ok {
						rateLimitingRule.ActionDuration = helper.String(v.(string))
					}

					if actionMap, ok := helper.InterfaceToMap(rulesMap, "action"); ok {
						securityAction := teov20220901.SecurityAction{}
						if v, ok := actionMap["name"]; ok {
							securityAction.Name = helper.String(v.(string))
						}

						if denyActionParametersMap, ok := helper.InterfaceToMap(actionMap, "deny_action_parameters"); ok {
							denyActionParameters := teov20220901.DenyActionParameters{}
							if v, ok := denyActionParametersMap["block_ip"]; ok {
								denyActionParameters.BlockIp = helper.String(v.(string))
							}

							if v, ok := denyActionParametersMap["block_ip_duration"]; ok {
								denyActionParameters.BlockIpDuration = helper.String(v.(string))
							}

							if v, ok := denyActionParametersMap["return_custom_page"]; ok {
								denyActionParameters.ReturnCustomPage = helper.String(v.(string))
							}

							if v, ok := denyActionParametersMap["response_code"]; ok {
								denyActionParameters.ResponseCode = helper.String(v.(string))
							}

							if v, ok := denyActionParametersMap["error_page_id"]; ok {
								denyActionParameters.ErrorPageId = helper.String(v.(string))
							}

							if v, ok := denyActionParametersMap["stall"]; ok {
								denyActionParameters.Stall = helper.String(v.(string))
							}

							securityAction.DenyActionParameters = &denyActionParameters
						}

						if redirectActionParametersMap, ok := helper.InterfaceToMap(actionMap, "redirect_action_parameters"); ok {
							redirectActionParameters := teov20220901.RedirectActionParameters{}
							if v, ok := redirectActionParametersMap["url"]; ok {
								redirectActionParameters.URL = helper.String(v.(string))
							}

							securityAction.RedirectActionParameters = &redirectActionParameters
						}

						if challengeActionParametersMap, ok := helper.InterfaceToMap(actionMap, "challenge_action_parameters"); ok {
							challengeActionParameters := teov20220901.ChallengeActionParameters{}
							if v, ok := challengeActionParametersMap["challenge_option"]; ok {
								challengeActionParameters.ChallengeOption = helper.String(v.(string))
							}

							if v, ok := challengeActionParametersMap["interval"]; ok {
								challengeActionParameters.Interval = helper.String(v.(string))
							}

							if v, ok := challengeActionParametersMap["attester_id"]; ok {
								challengeActionParameters.AttesterId = helper.String(v.(string))
							}

							securityAction.ChallengeActionParameters = &challengeActionParameters
						}

						if blockIPActionParametersMap, ok := helper.InterfaceToMap(actionMap, "block_ip_action_parameters"); ok {
							blockIPActionParameters := teov20220901.BlockIPActionParameters{}
							if v, ok := blockIPActionParametersMap["duration"]; ok {
								blockIPActionParameters.Duration = helper.String(v.(string))
							}

							securityAction.BlockIPActionParameters = &blockIPActionParameters
						}

						if returnCustomPageActionParametersMap, ok := helper.InterfaceToMap(actionMap, "return_custom_page_action_parameters"); ok {
							returnCustomPageActionParameters := teov20220901.ReturnCustomPageActionParameters{}
							if v, ok := returnCustomPageActionParametersMap["response_code"]; ok {
								returnCustomPageActionParameters.ResponseCode = helper.String(v.(string))
							}

							if v, ok := returnCustomPageActionParametersMap["error_page_id"]; ok {
								returnCustomPageActionParameters.ErrorPageId = helper.String(v.(string))
							}

							securityAction.ReturnCustomPageActionParameters = &returnCustomPageActionParameters
						}

						rateLimitingRule.Action = &securityAction
					}

					if v, ok := rulesMap["priority"]; ok {
						rateLimitingRule.Priority = helper.IntInt64(v.(int))
					}

					if v, ok := rulesMap["enabled"]; ok {
						rateLimitingRule.Enabled = helper.String(v.(string))
					}

					rateLimitingRules.Rules = append(rateLimitingRules.Rules, &rateLimitingRule)
				}

				securityPolicy.RateLimitingRules = &rateLimitingRules
			}
		}

		if exceptionRulesMap, ok := helper.ConvertInterfacesHeadToMap(securityPolicyMap["exception_rules"]); ok {
			exceptionRules := teov20220901.ExceptionRules{}
			if v, ok := exceptionRulesMap["rules"]; ok {
				for _, item := range v.([]interface{}) {
					rulesMap := item.(map[string]interface{})
					exceptionRule := teov20220901.ExceptionRule{}
					if v, ok := rulesMap["id"]; ok {
						exceptionRule.Id = helper.String(v.(string))
					}

					if v, ok := rulesMap["name"]; ok {
						exceptionRule.Name = helper.String(v.(string))
					}

					if v, ok := rulesMap["condition"]; ok {
						exceptionRule.Condition = helper.String(v.(string))
					}

					if v, ok := rulesMap["skip_scope"]; ok {
						exceptionRule.SkipScope = helper.String(v.(string))
					}

					if v, ok := rulesMap["skip_option"]; ok {
						exceptionRule.SkipOption = helper.String(v.(string))
					}

					if v, ok := rulesMap["web_security_modules_for_exception"]; ok {
						webSecurityModulesForExceptionSet := v.(*schema.Set).List()
						for i := range webSecurityModulesForExceptionSet {
							if webSecurityModulesForExceptionSet[i] != nil {
								webSecurityModulesForException := webSecurityModulesForExceptionSet[i].(string)
								exceptionRule.WebSecurityModulesForException = append(exceptionRule.WebSecurityModulesForException, &webSecurityModulesForException)
							}
						}
					}

					if v, ok := rulesMap["managed_rules_for_exception"]; ok {
						managedRulesForExceptionSet := v.(*schema.Set).List()
						for i := range managedRulesForExceptionSet {
							if managedRulesForExceptionSet[i] != nil {
								managedRulesForException := managedRulesForExceptionSet[i].(string)
								exceptionRule.ManagedRulesForException = append(exceptionRule.ManagedRulesForException, &managedRulesForException)
							}
						}
					}

					if v, ok := rulesMap["managed_rule_groups_for_exception"]; ok {
						managedRuleGroupsForExceptionSet := v.(*schema.Set).List()
						for i := range managedRuleGroupsForExceptionSet {
							if managedRuleGroupsForExceptionSet[i] != nil {
								managedRuleGroupsForException := managedRuleGroupsForExceptionSet[i].(string)
								exceptionRule.ManagedRuleGroupsForException = append(exceptionRule.ManagedRuleGroupsForException, &managedRuleGroupsForException)
							}
						}
					}

					if v, ok := rulesMap["request_fields_for_exception"]; ok {
						for _, item := range v.([]interface{}) {
							requestFieldsForExceptionMap := item.(map[string]interface{})
							requestFieldsForException := teov20220901.RequestFieldsForException{}
							if v, ok := requestFieldsForExceptionMap["scope"]; ok {
								requestFieldsForException.Scope = helper.String(v.(string))
							}

							if v, ok := requestFieldsForExceptionMap["condition"]; ok {
								requestFieldsForException.Condition = helper.String(v.(string))
							}

							if v, ok := requestFieldsForExceptionMap["target_field"]; ok {
								requestFieldsForException.TargetField = helper.String(v.(string))
							}

							exceptionRule.RequestFieldsForException = append(exceptionRule.RequestFieldsForException, &requestFieldsForException)
						}
					}

					if v, ok := rulesMap["enabled"]; ok {
						exceptionRule.Enabled = helper.String(v.(string))
					}

					exceptionRules.Rules = append(exceptionRules.Rules, &exceptionRule)
				}

				securityPolicy.ExceptionRules = &exceptionRules
			}
		} else {
			securityPolicy.ExceptionRules = &teov20220901.ExceptionRules{
				Rules: []*teov20220901.ExceptionRule{},
			}
		}

		if botManagementMap, ok := helper.ConvertInterfacesHeadToMap(securityPolicyMap["bot_management"]); ok {
			securityPolicy.BotManagement = expandBotManagement(botManagementMap)
		}

		request.SecurityPolicy = &securityPolicy
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifySecurityPolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo security policy failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s modify teo security policy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudTeoSecurityPolicyConfigRead(d, meta)
}

func resourceTencentCloudTeoSecurityPolicyConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_policy_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
