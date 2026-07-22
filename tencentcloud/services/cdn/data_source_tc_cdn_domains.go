package cdn

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCdnDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdnDomainsRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Acceleration 域名 名称",
			},
			"service_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SERVICE_TYPE),
				Description:  "Service 类型 acceleration 域名 名称 可用 值 include `web`，`download` 和 `media`。",
			},
			"full_url_cache": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "是否enable full-路径 缓存。",
			},
			"origin_pull_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_ORIGIN_PULL_PROTOCOL),
				Description:  "Origin-pull 协议 配置. 有效值：`http`，`https` 和 `follow`。",
			},
			"https_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_HTTPS_SWITCH),
				Description:  "HTTPS 配置. 有效值：`在`，`关闭` 和 `processing`。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"domain_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An 信息 列表 cdn 域名 Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 名称 ID。",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Acceleration 域名 名称",
						},
						"cname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CNAME 地址 的 域名 名称",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Acceleration 服务 状态",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "域名 名称 创建时间。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "最后修改时间 的 域名 名称",
						},
						"service_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service 类型 acceleration 域名 名称",
						},
						"area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Acceleration 地域",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "项目 CDN belongs 到。",
						},
						"full_url_cache": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否enable full-路径 缓存。",
						},
						"range_origin_switch": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sharding back 到 来源 配置 switch。",
						},
						"request_header": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Request 头部 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom 请求 头部 配置 switch。",
									},
									"header_rules": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Custom 请求 头部 配置 规则。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"header_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Http 头部 setting 方法。",
												},
												"header_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Http 头部 名称",
												},
												"header_value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Http 头部 值",
												},
												"rule_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Rule 类型",
												},
												"rule_paths": {
													Type:        schema.TypeList,
													Computed:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "Rule paths。",
												},
											},
										},
									},
								},
							},
						},
						"rule_cache": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Advanced 路径 缓存 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_paths": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "Rule paths。",
									},
									"rule_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule 类型",
									},
									"switch": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cache 配置 switch。",
									},
									"cache_time": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Cache 过期时间 setting， 单位 是 second。",
									},
									"compare_max_age": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Advanced 缓存 expiration 配置。",
									},
									"ignore_cache_control": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Force caching. After opening， 无-store 和 无-缓存 resources 返回 通过 源站 site 将 also 是 cached 在 accordance 使用 CacheRules 规则。",
									},
									"ignore_set_cookie": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Ignore Set-Cookie 头部 的 源站 site。",
									},
									"no_cache_switch": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cache 配置 switch。",
									},
									"re_validate": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Always check back 到 源站。",
									},
									"follow_origin_switch": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Follow 来源 station 配置 switch。",
									},
								},
							},
						},
						"origin": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Origin 服务器 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"origin_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Master 源站 服务器 类型",
									},
									"origin_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Master 源站 服务器 列表。",
									},
									"backup_origin_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Backup 源站 服务器 类型",
									},
									"backup_origin_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Backup 源站 服务器 列表。",
									},
									"backup_server_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "主机 头部 使用 当 accessing 备份 源站 服务器. 如果为空， ServerName 的 master 源站 服务器 将 是 使用 通过 默认值。",
									},
									"server_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "主机 头部 使用 当 accessing master 源站 服务器. 如果为空， acceleration 域名 名称 将 是 使用 通过 默认值。",
									},
									"cos_private_access": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "当 OriginType 是 COS，您 可以 指定if 访问 到 私有 buckets 是 allowed。",
									},
									"origin_pull_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Origin-pull 协议 配置。",
									},
								},
							},
						},
						"https_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "HTTPS acceleration 配置. It's 列表 和 consist 的 在 most 一个 item。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"https_switch": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "HTTPS 配置 switch。",
									},
									"http2_switch": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "HTTP2 配置 switch。",
									},
									"ocsp_stapling_switch": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "OCSP 配置 switch。",
									},
									"spdy_switch": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Spdy 配置 switch。",
									},
									"verify_client": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Client 证书 authentication 功能。",
									},
								},
							},
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "标签 的 cdn 域名",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCdnDomainsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cdn_domain.read")()
	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		domainConfigs []*cdn.DetailDomain
		err           error

		client     = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		region     = client.Region
		cdnService = CdnService{client: client}
		tagService = svctag.NewTagService(client)
	)

	var domainFilterMap = make(map[string]interface{}, 5)
	if v, ok := d.GetOk("domain"); ok {
		domainFilterMap["domain"] = v.(string)
	}
	if v, ok := d.GetOk("service_type"); ok {
		domainFilterMap["serviceType"] = v.(string)
	}
	if v, ok := d.GetOk("https_switch"); ok {
		domainFilterMap["httpsSwitch"] = v.(string)
	}
	if v, ok := d.GetOk("origin_pull_protocol"); ok {
		domainFilterMap["originPullProtocol"] = v.(string)
	}
	if v, ok := d.GetOkExists("full_url_cache"); ok {
		var value string
		if v.(bool) {
			value = "on"
		} else {
			value = "off"
		}

		domainFilterMap["fullUrlCache"] = value
	}

	domainConfigs, err = cdnService.DescribeDomainsConfigByFilters(ctx, domainFilterMap)
	if err != nil {
		log.Printf("[CRITAL]%s describeDomainsConfigByFilters fail, reason:%v ", logId, err)
		return err
	}

	cdnDomainList := make([]map[string]interface{}, 0, len(domainConfigs))
	ids := make([]string, 0, len(domainConfigs))
	for _, detailDomain := range domainConfigs {
		var fullUrlCache bool
		if detailDomain.CacheKey != nil && *detailDomain.CacheKey.FullUrlCache == CDN_SWITCH_ON {
			fullUrlCache = true
		}

		requestHeaders := make([]map[string]interface{}, 0, 1)
		requestHeader := make(map[string]interface{})
		requestHeader["switch"] = detailDomain.RequestHeader.Switch
		if len(detailDomain.RequestHeader.HeaderRules) > 0 {
			headerRules := make([]map[string]interface{}, len(detailDomain.RequestHeader.HeaderRules))
			headerRuleList := detailDomain.RequestHeader.HeaderRules
			for index, value := range headerRuleList {
				headerRule := make(map[string]interface{})
				headerRule["header_mode"] = value.HeaderMode
				headerRule["header_name"] = value.HeaderName
				headerRule["header_value"] = value.HeaderValue
				headerRule["rule_type"] = value.RuleType
				headerRule["rule_paths"] = value.RulePaths
				headerRules[index] = headerRule
			}
			requestHeader["header_rules"] = headerRules
		}
		requestHeaders = append(requestHeaders, requestHeader)

		ruleCaches := make([]map[string]interface{}, len(detailDomain.Cache.RuleCache))
		for index, value := range detailDomain.Cache.RuleCache {
			ruleCache := make(map[string]interface{})
			ruleCache["rule_paths"] = value.RulePaths
			ruleCache["rule_type"] = value.RuleType
			ruleCache["switch"] = value.CacheConfig.Cache.Switch
			ruleCache["cache_time"] = value.CacheConfig.Cache.CacheTime
			ruleCache["compare_max_age"] = value.CacheConfig.Cache.CompareMaxAge
			ruleCache["ignore_cache_control"] = value.CacheConfig.Cache.IgnoreCacheControl
			ruleCache["ignore_set_cookie"] = value.CacheConfig.Cache.IgnoreSetCookie
			ruleCache["no_cache_switch"] = value.CacheConfig.NoCache.Switch
			ruleCache["re_validate"] = value.CacheConfig.NoCache.Revalidate
			ruleCache["follow_origin_switch"] = value.CacheConfig.FollowOrigin.Switch
			ruleCaches[index] = ruleCache
		}

		origins := make([]map[string]interface{}, 0, 1)
		origin := make(map[string]interface{}, 8)
		origin["origin_type"] = detailDomain.Origin.OriginType
		origin["origin_list"] = detailDomain.Origin.Origins
		origin["backup_origin_type"] = detailDomain.Origin.BackupOriginType
		origin["backup_origin_list"] = detailDomain.Origin.BackupOrigins
		origin["backup_server_name"] = detailDomain.Origin.BackupServerName
		origin["server_name"] = detailDomain.Origin.ServerName
		origin["cos_private_access"] = detailDomain.Origin.CosPrivateAccess
		origin["origin_pull_protocol"] = detailDomain.Origin.OriginPullProtocol
		origins = append(origins, origin)

		httpsconfigs := make([]map[string]interface{}, 0, 1)
		if detailDomain.Https != nil {
			httpsConfig := make(map[string]interface{}, 7)
			httpsConfig["https_switch"] = detailDomain.Https.Switch
			httpsConfig["http2_switch"] = detailDomain.Https.Http2
			httpsConfig["ocsp_stapling_switch"] = detailDomain.Https.OcspStapling
			httpsConfig["spdy_switch"] = detailDomain.Https.Spdy
			httpsConfig["verify_client"] = detailDomain.Https.VerifyClient
			httpsconfigs = append(httpsconfigs, httpsConfig)
		}

		tags, errRet := tagService.DescribeResourceTags(ctx, CDN_SERVICE_NAME, CDN_RESOURCE_NAME_DOMAIN, region, *detailDomain.Domain)
		if errRet != nil {
			return errRet
		}

		mapping := map[string]interface{}{
			"id":                  detailDomain.ResourceId,
			"domain":              detailDomain.Domain,
			"cname":               detailDomain.Cname,
			"status":              detailDomain.Status,
			"create_time":         detailDomain.CreateTime,
			"update_time":         detailDomain.UpdateTime,
			"service_type":        detailDomain.ServiceType,
			"area":                detailDomain.Area,
			"project_id":          detailDomain.ProjectId,
			"full_url_cache":      fullUrlCache,
			"range_origin_switch": detailDomain.RangeOriginPull.Switch,
			"request_header":      requestHeaders,
			"rule_cache":          ruleCaches,
			"origin":              origins,
			"https_config":        httpsconfigs,
			"tags":                tags,
		}

		cdnDomainList = append(cdnDomainList, mapping)
		ids = append(ids, *detailDomain.ResourceId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("domain_list", cdnDomainList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set cdn domain list fail, reason:%v ", logId, err)
		return err
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), cdnDomainList); err != nil {
			return err
		}
	}
	return nil
}
