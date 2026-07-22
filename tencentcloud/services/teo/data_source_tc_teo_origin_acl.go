package teo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTeoOriginAcl() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoOriginAclRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "指定site ID。",
			},

			"origin_acl_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Describes binding relationship between l7 acceleration 域名/l4 proxy 实例 和 源站 服务器 IP 范围。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"l7_hosts": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "列表 L7 accelerated domains 该 启用 源站 ACLs. 此 字段 是 空 当 源站 protection 是 不 已启用",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"l4_proxy_ids": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "列表 L4 proxy 实例 该 启用 源站 ACLs. 此 字段 是 空 当 源站 protection 是 不 已启用",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"current_origin_acl": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Currently effective 源站 ACLs. 此 字段 是 空 当 源站 protection 是 不 已启用\nNote: 此 字段 可能 返回 null，其中 表示a failure 到 obtain 有效 值",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"entire_addresses": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "IP 范围 details.\nNote: 此 字段 可能 返回 null，其中 表示a failure 到 obtain 有效 值",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"i_pv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 子网",
													Deprecated:  "Field `i_pv4` has been deprecated from version 1.82.27. Use new field `ipv4` instead.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"i_pv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 子网",
													Deprecated:  "Field `i_pv6` has been deprecated from version 1.82.27. Use new field `ipv6` instead.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 子网",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 子网",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "版本 数量.\nNote: 此 字段 可能 返回 null，其中 表示a failure 到 obtain 有效 值",
									},
									"active_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "版本 effective 时间 在 UTC+8，following date 和 时间格式 的 ISO 8601 standard.\nNote: 此 字段 可能 返回 null，其中 表示a failure 到 obtain 有效 值",
									},
									"is_planed": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "此 参数 是 使用 到 记录 whether \"I've upgraded 到 lastest 版本\" 是 completed before 源站 ACLs 版本 是 effective. 有效 值:.\n- true: specifies 该 版本 是 effective 和 update 到 latest 版本 是 confirmed.\n- false: 当 版本 takes effect, confirmation 的 updating 到 latest 源站 ACLs 是 不 completed. IP 范围 是 forcibly 更新 到 latest 版本 在 backend. 当 此 参数 returns false, please confirm 在 时间 whether your 源站 服务器 firewall 配置 has been 更新 到 latest 版本 到 avoid 源站-pull failure.\nNote: 此 字段 可能 返回 null, 其中 indicates failure 到 obtain 有效 值.",
									},
								},
							},
						},
						"next_origin_acl": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "当 源站 ACLs 是 更新，此 字段 将 是 返回 使用 next 版本's 源站 IP 范围 到 take effect，包括 comparison 使用 当前 源站 IP 范围. 此 字段 是 空 如果 不 更新 或 源站 protection 是 不 已启用\nNote: 此 字段 可能 返回 null，其中 表示a failure 到 obtain 有效 值",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "版本 数量。",
									},
									"planned_active_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "版本 effective 时间，其中 adopts UTC+8 和 follows date 和 时间格式 的 ISO 8601 standard。",
									},
									"entire_addresses": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "IP 范围 details。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"i_pv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 子网",
													Deprecated:  "Field `i_pv4` has been deprecated from version 1.82.27. Use new field `ipv4` instead.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"i_pv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 子网",
													Deprecated:  "Field `i_pv6` has been deprecated from version 1.82.27. Use new field `ipv6` instead.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 子网",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 子网",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"added_addresses": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "latest 源站 IP 范围 newly-added compared 使用 源站 IP 范围 在 CurrentOrginACL。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"i_pv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 子网",
													Deprecated:  "Field `i_pv4` has been deprecated from version 1.82.27. Use new field `ipv4` instead.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"i_pv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 子网",
													Deprecated:  "Field `i_pv6` has been deprecated from version 1.82.27. Use new field `ipv6` instead.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 子网",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 子网",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"removed_addresses": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "latest 源站 IP 范围 删除 compared 使用 源站 IP 范围 在 CurrentOrginACL。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"i_pv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 子网",
													Deprecated:  "Field `i_pv4` has been deprecated from version 1.82.27. Use new field `ipv4` instead.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"i_pv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 子网",
													Deprecated:  "Field `i_pv6` has been deprecated from version 1.82.27. Use new field `ipv6` instead.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 子网",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 子网",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"no_change_addresses": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "latest 源站 IP 范围 是 unchanged compared 使用 源站 IP 范围 在 CurrentOrginACL。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"i_pv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 子网",
													Deprecated:  "Field `i_pv4` has been deprecated from version 1.82.27. Use new field `ipv4` instead.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"i_pv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 子网",
													Deprecated:  "Field `i_pv6` has been deprecated from version 1.82.27. Use new field `ipv6` instead.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv4": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv4 子网",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ipv6": {
													Type:        schema.TypeSet,
													Computed:    true,
													Description: "IPv6 子网",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Origin protection 状态 Vaild 值:\n- online: 在 effect;\n- offline: 已禁用;\n- updating: 配置 部署 在 progress。",
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

func dataSourceTencentCloudTeoOriginAclRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_origin_acl.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		zoneId  string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		paramMap["ZoneId"] = helper.String(v.(string))
		zoneId = v.(string)
	}

	var respData *teov20220901.DescribeOriginACLResponseParams
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTeoOriginAclByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	originACLInfoMap := map[string]interface{}{}
	if respData.OriginACLInfo != nil {
		if respData.OriginACLInfo.L7Hosts != nil {
			originACLInfoMap["l7_hosts"] = respData.OriginACLInfo.L7Hosts
		}

		if respData.OriginACLInfo.L4ProxyIds != nil {
			originACLInfoMap["l4_proxy_ids"] = respData.OriginACLInfo.L4ProxyIds
		}

		currentOriginACLMap := map[string]interface{}{}
		if respData.OriginACLInfo.CurrentOriginACL != nil {
			entireAddressesMap := map[string]interface{}{}
			if respData.OriginACLInfo.CurrentOriginACL.EntireAddresses != nil {
				if respData.OriginACLInfo.CurrentOriginACL.EntireAddresses.IPv4 != nil {
					entireAddressesMap["i_pv4"] = respData.OriginACLInfo.CurrentOriginACL.EntireAddresses.IPv4
					entireAddressesMap["ipv4"] = respData.OriginACLInfo.CurrentOriginACL.EntireAddresses.IPv4
				}

				if respData.OriginACLInfo.CurrentOriginACL.EntireAddresses.IPv6 != nil {
					entireAddressesMap["i_pv6"] = respData.OriginACLInfo.CurrentOriginACL.EntireAddresses.IPv6
					entireAddressesMap["ipv6"] = respData.OriginACLInfo.CurrentOriginACL.EntireAddresses.IPv6
				}

				currentOriginACLMap["entire_addresses"] = []interface{}{entireAddressesMap}
			}

			if respData.OriginACLInfo.CurrentOriginACL.Version != nil {
				currentOriginACLMap["version"] = respData.OriginACLInfo.CurrentOriginACL.Version
			}

			if respData.OriginACLInfo.CurrentOriginACL.ActiveTime != nil {
				currentOriginACLMap["active_time"] = respData.OriginACLInfo.CurrentOriginACL.ActiveTime
			}

			if respData.OriginACLInfo.CurrentOriginACL.IsPlaned != nil {
				currentOriginACLMap["is_planed"] = respData.OriginACLInfo.CurrentOriginACL.IsPlaned
			}

			originACLInfoMap["current_origin_acl"] = []interface{}{currentOriginACLMap}
		}

		nextOriginACLMap := map[string]interface{}{}
		if respData.OriginACLInfo.NextOriginACL != nil {
			if respData.OriginACLInfo.NextOriginACL.Version != nil {
				nextOriginACLMap["version"] = respData.OriginACLInfo.NextOriginACL.Version
			}

			if respData.OriginACLInfo.NextOriginACL.PlannedActiveTime != nil {
				nextOriginACLMap["planned_active_time"] = respData.OriginACLInfo.NextOriginACL.PlannedActiveTime
			}

			entireAddressesMap := map[string]interface{}{}
			if respData.OriginACLInfo.NextOriginACL.EntireAddresses != nil {
				if respData.OriginACLInfo.NextOriginACL.EntireAddresses.IPv4 != nil {
					entireAddressesMap["i_pv4"] = respData.OriginACLInfo.NextOriginACL.EntireAddresses.IPv4
					entireAddressesMap["ipv4"] = respData.OriginACLInfo.NextOriginACL.EntireAddresses.IPv4
				}

				if respData.OriginACLInfo.NextOriginACL.EntireAddresses.IPv6 != nil {
					entireAddressesMap["i_pv6"] = respData.OriginACLInfo.NextOriginACL.EntireAddresses.IPv6
					entireAddressesMap["ipv6"] = respData.OriginACLInfo.NextOriginACL.EntireAddresses.IPv6
				}

				nextOriginACLMap["entire_addresses"] = []interface{}{entireAddressesMap}
			}

			addedAddressesMap := map[string]interface{}{}
			if respData.OriginACLInfo.NextOriginACL.AddedAddresses != nil {
				if respData.OriginACLInfo.NextOriginACL.AddedAddresses.IPv4 != nil {
					addedAddressesMap["i_pv4"] = respData.OriginACLInfo.NextOriginACL.AddedAddresses.IPv4
					addedAddressesMap["ipv4"] = respData.OriginACLInfo.NextOriginACL.AddedAddresses.IPv4
				}

				if respData.OriginACLInfo.NextOriginACL.AddedAddresses.IPv6 != nil {
					addedAddressesMap["i_pv6"] = respData.OriginACLInfo.NextOriginACL.AddedAddresses.IPv6
					addedAddressesMap["ipv6"] = respData.OriginACLInfo.NextOriginACL.AddedAddresses.IPv6
				}

				nextOriginACLMap["added_addresses"] = []interface{}{addedAddressesMap}
			}

			removedAddressesMap := map[string]interface{}{}
			if respData.OriginACLInfo.NextOriginACL.RemovedAddresses != nil {
				if respData.OriginACLInfo.NextOriginACL.RemovedAddresses.IPv4 != nil {
					removedAddressesMap["i_pv4"] = respData.OriginACLInfo.NextOriginACL.RemovedAddresses.IPv4
					removedAddressesMap["ipv4"] = respData.OriginACLInfo.NextOriginACL.RemovedAddresses.IPv4
				}

				if respData.OriginACLInfo.NextOriginACL.RemovedAddresses.IPv6 != nil {
					removedAddressesMap["i_pv6"] = respData.OriginACLInfo.NextOriginACL.RemovedAddresses.IPv6
					removedAddressesMap["ipv6"] = respData.OriginACLInfo.NextOriginACL.RemovedAddresses.IPv6
				}

				nextOriginACLMap["removed_addresses"] = []interface{}{removedAddressesMap}
			}

			noChangeAddressesMap := map[string]interface{}{}
			if respData.OriginACLInfo.NextOriginACL.NoChangeAddresses != nil {
				if respData.OriginACLInfo.NextOriginACL.NoChangeAddresses.IPv4 != nil {
					noChangeAddressesMap["i_pv4"] = respData.OriginACLInfo.NextOriginACL.NoChangeAddresses.IPv4
					noChangeAddressesMap["ipv4"] = respData.OriginACLInfo.NextOriginACL.NoChangeAddresses.IPv4
				}

				if respData.OriginACLInfo.NextOriginACL.NoChangeAddresses.IPv6 != nil {
					noChangeAddressesMap["i_pv6"] = respData.OriginACLInfo.NextOriginACL.NoChangeAddresses.IPv6
					noChangeAddressesMap["ipv6"] = respData.OriginACLInfo.NextOriginACL.NoChangeAddresses.IPv6
				}

				nextOriginACLMap["no_change_addresses"] = []interface{}{noChangeAddressesMap}
			}

			originACLInfoMap["next_origin_acl"] = []interface{}{nextOriginACLMap}
		}

		if respData.OriginACLInfo.Status != nil {
			originACLInfoMap["status"] = respData.OriginACLInfo.Status
		}

		_ = d.Set("origin_acl_info", []interface{}{originACLInfoMap})
	}

	d.SetId(zoneId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), originACLInfoMap); e != nil {
			return e
		}
	}

	return nil
}
