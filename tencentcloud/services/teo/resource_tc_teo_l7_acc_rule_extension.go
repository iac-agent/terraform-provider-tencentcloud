package teo

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func TencentTeoL7RuleBranchBasicInfo(depth int) map[string]*schema.Schema {
	schemaMap := map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Match condition. https://www.tencentcloud.com/document/product/1145/54759。",
		},
		"actions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Sub-Rule branch. 此 列表 currently 支持 filling 在 仅 一个 规则; 多个 entries 是 无效。",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Operation 名称 名称 needs 到 correspond 到 参数 structure，对于 示例，如果 名称=Cache，CacheParameters 为必填项.\n- `Cache`: Node 缓存 TTL;\n- `CacheKey`: Custom Cache 键;\n- `CachePrefresh`: Cache pre-refresh;\n- `AccessURLRedirect`: Access URL redirection;\n- `UpstreamURLRewrite`: Back-到-源站 URL rewrite;\n- `QUIC`: QUIC;\n- `WebSocket`: WebSocket;\n- `Authentication`: 令牌 authentication;\n- `MaxAge`: Browser 缓存 TTL;\n- `StatusCodeCache`: 状态 代码 缓存 TTL;\n- `OfflineCache`: Offline 缓存;\n- `SmartRouting`: Smart acceleration;\n- `RangeOriginPull`: Segment back-到-源站;\n- `UpstreamHTTP2`: HTTP2 back-到-源站;\n- `HostHeader`: 主机 Header rewrite;\n- `ForceRedirectHTTPS`: Access 协议 forced HTTPS jump 配置;\n- `OriginPullProtocol`: Back-到-源站 HTTPS;\n- `Compression`: Smart 压缩 配置;\n- `HSTS`: HSTS;\n- `ClientIPHeader`: Header 信息 配置 对于 storing 客户端 请求 IP;\n- `OCSPStapling`: OCSP stapling;\n- `HTTP2`: HTTP2 Access;\n- `PostMaxSize`: POST 请求 upload 文件 streaming 最大 限制 配置;\n- `ClientIPCountry`: Carry 客户端 IP 地域 信息 当 returning 到 来源;\n- `UpstreamFollowRedirect`: Return 到 来源 follow redirection 参数 配置;\n- `UpstreamRequest`: Return 到 来源 请求 参数;\n- `TLSConfig`: SSL/TLS 安全;\n- `ModifyOrigin`: Modify 来源 station;\n- `HTTPUpstreamTimeout`: Seven-layer 返回 到 来源 超时 配置;\n- `HttpResponse`: HTTP response;\n- `ErrorPage`: Custom 错误 页面;\n- `ModifyResponseHeader`: Modify HTTP 节点 response 头部;\n- `ModifyRequestHeader`: Modify HTTP 节点 请求 头部;\n- `ResponseSpeedLimit`: Single 连接 download speed 限制\n- `SetContentIdentifierParameters`: Set 内容 identifier。",
					},
					"cache_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Node 缓存 ttl 配置 参数. 当 名称 是 缓存，此 参数 为必填项。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"follow_origin": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "Cache follows 源站 服务器. 如果未指定，此 配置 是 不 集合. 仅 一个 的 followorigin，nocache，或 customtime 可以 have switch 集合 到 在。",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"switch": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "是否enable 配置 的 following 源站 服务器. 有效值：`在`: Enable; `关闭`: Disable。",
											},
											"default_cache": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "是否cache 当 源站 服务器 does 不 返回 缓存-control 头部. 此 字段 为必填项 当 switch 是 在; 当 switch 是 关闭，此 字段 不是必填项 和 将 是 ineffective 如果 filled. 有效值：On: 缓存; Off: do 不 缓存。",
											},
											"default_cache_strategy": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "是否use 默认值 caching 策略 当 源站 服务器 does 不 返回 缓存-control 头部. 此 字段 为必填项 当 defaultcache 是 集合 到 在; otherwise，它 是 ineffective. 当 defaultcachetime 是 不 0，此 字段 should 是 关闭. 有效值：在: 使用 默认值 caching 策略. 关闭: do 不 使用 默认值 caching 策略。",
											},
											"default_cache_time": {
												Type:        schema.TypeInt,
												Optional:    true,
												Description: "默认值 缓存 时间 （秒） 当 源站 服务器 does 不 返回 缓存-control 头部. 值 ranges 从 0 到 315360000. 此 字段 为必填项 当 defaultcache 是 集合 到 在; otherwise，它 是 ineffective. 当 defaultcachestrategy 是 在，此 字段 should 是 0。",
											},
										},
									},
								},
								"no_cache": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "No 缓存. 如果未指定，此 配置 是 不 集合. 仅 一个 的 followorigin，nocache，或 customtime 可以 have switch 集合 到 在。",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"switch": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "是否enable 无-缓存 配置. 有效值：`在`: Enable; `关闭`: Disable。",
											},
										},
									},
								},
								"custom_time": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "Custom 缓存 时间. 如果未指定，此 配置 是 不 集合. 仅 一个 的 followorigin，nocache，或 customtime 可以 have switch 集合 到 在。",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"switch": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Custom 缓存 时间 switch. 值: `在`: Enable; `关闭`: Disable。",
											},
											"ignore_cache_control": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Ignore 源站 服务器 cachecontrol switch. 值: `在`: Enable; `关闭`: Disable。",
											},
											"cache_time": {
												Type:        schema.TypeInt,
												Optional:    true,
												Description: "Custom 缓存 时间 值，单位: 秒. 取值范围：0-315360000。",
											},
										},
									},
								},
							},
						},
					},
					"cache_key_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Custom 缓存 键 配置 参数. 当 名称 是 cachekey，此 参数 为必填项。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"full_url_cache": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Switch 对于 retaining 完整 查询 字符串. 值: 在: 启用; 关闭: disable。",
								},
								"query_string": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "Configuration 参数 对于 retaining 查询 字符串. 此 字段 和 fullurlcache 必须 是 集合 simultaneously，但 不能 both 是 在。",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"switch": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Query 字符串 retain/ignore 指定 参数 switch. 有效值：在: 启用; 关闭: disable。",
											},
											"action": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Actions 到 retain/ignore 指定 参数 在 查询 字符串. 值: `includeCustom`: retain partial 参数. `excludeCustom`: ignore partial 参数.note: 此 字段 为必填项 当 switch 是 在. 当 switch 是 关闭，此 字段 不是必填项 和 将 不 take effect 如果 filled。",
											},
											"values": {
												Type:        schema.TypeList,
												Optional:    true,
												Description: "A 列表 参数 names 到 keep/ignore 在 查询 字符串。",
												Elem: &schema.Schema{
													Type: schema.TypeString,
												},
											},
										},
									},
								},
								"ignore_case": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Switch 对于 ignoring case. 值: 启用; 关闭: disable.note: 在 least 一个 的 fullurlcache，ignorecase，头部，scheme，或 cookie 必须 是 已配置。",
								},
								"header": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "HTTP 请求 头部 配置 参数. 在 least 一个 的 following configurations 必须 是 集合: fullurlcache，ignorecase，头部，scheme，cookie。",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"switch": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "是否enable 功能. 值: 在: 启用; 关闭: disable。",
											},
											"values": {
												Type:        schema.TypeList,
												Optional:    true,
												Description: "Custom 缓存 键 http 请求 头部 列表. note: 此 字段 为必填项 当 switch 是 在; 当 switch 是 关闭，此 字段 不是必填项 和 将 不 take effect 如果 filled。",
												Elem: &schema.Schema{
													Type: schema.TypeString,
												},
											},
										},
									},
								},
								"scheme": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Request 协议 switch. 有效值：在: 启用; 关闭: disable。",
								},
								"cookie": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "Cookie 配置 参数. 在 least 一个 的 following configurations 必须 是 集合: fullurlcache，ignorecase，头部，scheme，cookie。",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"switch": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "是否enable 功能. 值: 在: 启用; 关闭: disable。",
											},
											"action": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Cache 操作 值: full: retain all; ignore: ignore all; includeCustom: retain partial 参数; excludeCustom: ignore partial 参数. note: 当 switch 是 在，此 字段 为必填项. 当 switch 是 关闭，此 字段 不是必填项 和 将 不 take effect 如果 filled。",
											},
											"values": {
												Type:        schema.TypeList,
												Optional:    true,
												Description: "Custom 缓存 键 cookie 名称 列表。",
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
					"cache_prefresh_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "缓存 prefresh 配置 参数. 此 参数 为必填项 当 名称 是 cacheprefresh。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "是否enable 缓存 prefresh. 值: 启用; 关闭: disable。",
								},
								"cache_time_percent": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "Prefresh 间隔 集合 作为 percentage 的 节点 缓存 时间. 取值范围：1-99. note: 此 字段 为必填项 当 switch 是 在; 当 switch 是 关闭，此 字段 不是必填项 和 将 不 take effect 如果 filled。",
								},
							},
						},
					},
					"access_url_redirect_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "访问 URL redirection 配置 参数. 此 参数 为必填项 当 名称 是 accessurlredirect。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"status_code": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "状态 代码 有效值：301，302，303，307，308。",
								},
								"protocol": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Target 请求 协议 有效值：http: 目标 请求 协议 http; https: 目标 请求 协议 https; follow: follow 请求。",
								},
								"host_name": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "Target hostname。",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"action": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Target hostname 配置，有效值：follow: follow 请求; 自定义: 自定义。",
											},
											"value": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Custom 值 对于 目标 hostname，最大 长度 是 1024。",
											},
										},
									},
								},
								"url_path": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "Target 路径",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"action": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "操作 到 是 executed. 值: follow: follow 请求; 自定义: 自定义; regex: regular expression matching。",
											},
											"regex": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Regular expression matching expression，长度 范围 是 1-1024. note: 当 操作 是 regex，此 字段 为必填项; 当 操作 是 follow 或 自定义，此 字段 不是必填项 和 将 不 take effect 如果 filled。",
											},
											"value": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Redirect 目标 URL，长度 范围 是 1-1024.note: 当 操作 是 regex 或 自定义，此 字段 为必填项; 当 操作 是 follow，此 字段 不是必填项 和 将 不 take effect 如果 filled。",
											},
										},
									},
								},
								"query_string": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "Carry 查询 参数。",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"action": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "操作 到 是 executed. 值: full: retain all; ignore: ignore all。",
											},
										},
									},
								},
							},
						},
					},
					"upstream_url_rewrite_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "源站-pull URL rewrite 配置 参数. 此 参数 为必填项 当 名称 是 upstreamurlrewrite。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"type": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Origin-Pull URL rewriting 类型，仅 路径 是 支持。",
								},
								"action": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Origin-Pull URL rewrite 操作 有效值：replace: replace 路径 prefix; addPrefix: add 路径 prefix; rmvPrefix: remove 路径 prefix。",
								},
								"value": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Origin-Pull URL rewrite 值，最大 长度 1024，必须 start 使用 /.note: 当 操作 是 addprefix，它 不能 end 使用 /; 当 操作 是 rmvprefix，* 不能 是 present。",
								},
								"regex": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Origin URL Rewrite uses regular expression 对于 matching 完整 路径 It 必须 conform 到 Google RE2 规格 和 have 长度 范围 的 1 到 1024. 此 字段 为必填项 当 操作 是 regexReplace; otherwise，它 为可选项。",
								},
							},
						},
					},
					"quic_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "quic 配置 参数. 此 参数 为必填项 当 名称 是 quic。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "是否enable quic. 值: 在: 启用; 关闭: disable。",
								},
							},
						},
					},
					"web_socket_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "websocket 配置 参数. 此 参数 为必填项 当 名称 是 websocket。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "是否enable websocket 连接 超时. 值: 在: 使用 超时 作为 websocket 超时;; 关闭: 平台 still 支持 websocket connections，使用 系统 默认值 超时 的 15 秒。",
								},
								"timeout": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "Timeout，单位: 秒. 最大 超时 是 120 秒。",
								},
							},
						},
					},
					"authentication_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "令牌 authentication 配置 参数. 此 参数 为必填项 当 名称 是 authentication。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"auth_type": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Authentication 类型 有效 值:\n- `TypeA`: authentication 方法 类型，对于 特定 meaning please refer 到 authentication 方法 . https://www.tencentcloud.com/document/product/1145/62475;\n- `TypeB`: authentication 方法 b 类型，对于 特定 meaning please refer 到 authentication 方法 b. https://www.tencentcloud.com/document/product/1145/62476;\n- `TypeC`: authentication 方法 c 类型，对于 特定 meaning please refer 到 authentication 方法 c. https://www.tencentcloud.com/document/product/1145/62477;\n- `TypeD`: authentication 方法 d 类型，对于 特定 meaning please refer 到 authentication 方法 d. https://www.tencentcloud.com/document/product/1145/62478;\n- `TypeVOD`: authentication 方法 v 类型，对于 特定 meaning please refer 到 authentication 方法 v. https://www.tencentcloud.com/document/product/1145/62479。",
								},
								"secret_key": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "primary authentication 键 consists 的 6-40 uppercase 和 lowercase english letters 或 digits, 和 不能 contain \" 和 $.",
								},
								"timeout": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "Validity 周期 的 authentication url, 在 秒, 值 范围: 1-630720000. 使用 到 determine 如果 客户端 访问 请求 has expired: 如果 当前 时间 exceeds \"timestamp + validity 周期\", 它 是 expired 请求, 和 403 是 返回 directly. 如果 当前 时间 does 不 exceed \"timestamp + validity 周期\", 请求 是 不 expired, 和 md5 字符串 是 further validated. note: 当 authtype 是 一个 的 typea, typeb, typec, 或 typed, 此 字段 是 必填.",
								},
								"backup_secret_key": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "备份 authentication 键 consists 的 6-40 uppercase 和 lowercase english letters 或 digits, 和 不能 contain \" 和 $.",
								},
								"auth_param": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Authentication 参数 名称 节点 将 validate 值 corresponding 到 此 参数 名称 consists 的 1-100 uppercase 和 lowercase letters，numbers，或 underscores.note: 此 字段 为必填项 当 authtype 是 either typea 或 typed。",
								},
								"time_param": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Authentication 时间戳. 它 不能 是 same 作为 值 的 authparam 字段.note: 此 字段 为必填项 当 authtype 是 typed。",
								},
								"time_format": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Authentication 时间格式. 值: dec: decimal; hex: hexadecimal。",
								},
							},
						},
					},
					"max_age_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Browser 缓存 ttl 配置 参数. 此 参数 为必填项 当 名称 是 maxage。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"follow_origin": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "指定是否follow 源站 服务器 缓存-control 配置，使用 following 值: 在: follow 源站 服务器 和 ignore 字段 cachetime; 关闭: do 不 follow 源站 服务器 和 apply 字段 cachetime。",
								},
								"cache_time": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "Custom 缓存 时间 值，单位: 秒. 取值范围：0-315360000. note: 当 followorigin 是 关闭，它 表示 不 following 源站 服务器 和 使用 cachetime 到 集合 缓存 时间; otherwise，此 字段 将 不 take effect。",
								},
							},
						},
					},
					"status_code_cache_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "状态 代码 缓存 ttl 配置 参数. 此 参数 为必填项 当 名称 是 statuscodecache。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"status_code_cache_params": {
									Type:        schema.TypeList,
									Optional:    true,
									Description: "状态 代码 缓存 ttl。",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"status_code": {
												Type:        schema.TypeInt,
												Optional:    true,
												Description: "状态 代码 有效值：400，401，403，404，405，407，414，500，501，502，503，504，509，514。",
											},
											"cache_time": {
												Type:        schema.TypeInt,
												Optional:    true,
												Description: "Cache 时间 值 （秒）。 取值范围：0-31536000。",
											},
										},
									},
								},
							},
						},
					},
					"offline_cache_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Offline 缓存 配置 参数. 此 参数 为必填项 当 名称 是 offlinecache。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "是否enable offline caching. 值: 在: 启用; Off: disable。",
								},
							},
						},
					},
					"smart_routing_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Smart acceleration 配置 参数. 此 参数 为必填项 当 名称 是 smartrouting。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "是否enable smart acceleration. 值: 在: 启用; Off: disable。",
								},
							},
						},
					},
					"range_origin_pull_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Shard 来源 retrieval 配置 参数. 此 参数 为必填项 当 名称 是 集合 到 rangeoriginpull。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "是否enable 范围 gets. 值 是: 在: 启用; Off: disable。",
								},
							},
						},
					},
					"upstream_http2_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "HTTP2 源站-pull 配置 参数. 此 参数 为必填项 当 名称 是 集合 到 upstreamhttp2。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "是否enable http2 源站-pull. 有效值：在: 启用; 关闭: disable。",
								},
							},
						},
					},
					"host_header_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "主机 头部 rewrite 配置 参数. 此 参数 为必填项 当 名称 是 集合 到 hostheader。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"action": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "操作 到 是 executed. 值: followOrigin: follow 源站 服务器 域名 名称; 自定义: 自定义。",
								},
								"server_name": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "主机 头部 rewrite requires 完整 域名 名称 note: 此 字段 为必填项 当 switch 是 在; 当 switch 是 关闭，此 字段 不是必填项 和 any 值 将 是 ignored。",
								},
							},
						},
					},
					"force_redirect_https_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Force https redirect 配置 参数. 此 参数 为必填项 当 名称 是 集合 到 forceredirecthttps。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "是否enable forced redirect 配置 switch. 值: 在: 启用; 关闭: disable。",
								},
								"redirect_status_code": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "Redirection 状态 代码 此 字段 为必填项 当 switch 是 在; otherwise，它 是 不 effective. 有效值：301: 301 redirect; 302: 302 redirect。",
								},
							},
						},
					},
					"origin_pull_protocol_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Back-到-源站 HTTPS 配置 参数. 此 参数 为必填项 当 名称 值 是 `OriginPullProtocol`。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"protocol": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Back-到-源站 协议 配置. Possible 值 是: `http`: 使用 HTTP 协议 对于 back-到-源站; `https`: 使用 HTTPS 协议 对于 back-到-源站; `follow`: follow 协议",
								},
							},
						},
					},
					"compression_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Intelligent 压缩 配置. 此 参数 为必填项 当 名称 是 集合 到 压缩。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "是否enable smart 压缩. 值: 在: 启用; 关闭: disable。",
								},
								"algorithms": {
									Type:        schema.TypeList,
									Optional:    true,
									Description: "Supported 压缩 algorithm 列表. 此 字段 为必填项 当 switch 是 在; otherwise，它 是 不 effective. 有效值：brotli: brotli algorithm; gzip: gzip algorithm。",
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
							},
						},
					},
					"hsts_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "HSTS 配置 参数. 此 参数 为必填项 当 名称 是 hsts。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "是否enable hsts. 值: 在: 启用; 关闭: disable。",
								},
								"timeout": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "Cache hsts 头部 时间，单位: 秒. 取值范围：1-31536000. note: 此 字段 为必填项 当 switch 是 在; 当 switch 是 关闭，此 字段 不是必填项 和 将 不 take effect 如果 filled。",
								},
								"include_sub_domains": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "是否allow other subdomains 到 inherit same hsts 头部. 值: 在: allows other subdomains 到 inherit same hsts 头部; 关闭: does 不 allow other subdomains 到 inherit same hsts 头部. note: 当 switch 是 在，此 字段 为必填项; 当 switch 是 关闭，此 字段 不是必填项 和 将 不 take effect 如果 filled。",
								},
								"preload": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "是否allow browser 到 preload hsts 头部. 有效值：在: allows browser 到 preload hsts 头部; 关闭: does 不 allow browser 到 preload hsts 头部. note: 当 switch 是 在，此 字段 为必填项; 当 switch 是 关闭，此 字段 不是必填项 和 将 不 take effect 如果 filled。",
								},
							},
						},
					},
					"client_ip_header_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Client ip 头部 配置 对于 storing 客户端 请求 ip 信息. 此 参数 为必填项 当 名称 是 clientipheader。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "是否enable 配置. 值: 在: 启用; 关闭: disable。",
								},
								"header_name": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "名称 请求 头部 containing 客户端 ip 地址 对于 源站-pull. 当 switch 是 在，此 参数 为必填项. x-forwarded-对于 是 不 allowed 对于 此 参数。",
								},
							},
						},
					},
					"ocsp_stapling_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "OCSP stapling 配置 参数. 此 参数 为必填项 当 名称 是 集合 到 ocspstapling。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "是否enable ocsp stapling 配置 switch. 值: 在: 启用; 关闭: disable。",
								},
							},
						},
					},
					"http2_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "HTTP2 访问 配置 参数. 此 参数 为必填项 当 名称 是 http2。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "是否enable http2 访问. 值: 在: 启用; 关闭: disable。",
								},
							},
						},
					},
					"post_max_size_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Maximum 大小 配置 对于 文件 streaming upload via post 请求. 此 参数 为必填项 当 名称 是 postmaxsize。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "是否enable post 请求 文件 upload 限制，在 bytes (默认值 限制: 32 * 2^20 bytes). 有效值：在: 启用 限制; 关闭: disable 限制",
								},
								"max_size": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "Maximum 大小 的 文件 uploaded 对于 streaming via post 请求，在 bytes. 取值范围：1 * 2^20 bytes 到 500 * 2^20 bytes。",
								},
							},
						},
					},
					"client_ip_country_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Configuration 参数 对于 carrying 地域 信息 的 客户端 ip during 源站-pull. 此 参数 为必填项 当 名称 是 集合 到 clientipcountry。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "是否enable 配置. 值: 在: 启用; 关闭: disable。",
								},
								"header_name": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "名称 请求 头部 该 包含client ip 地域 它 是 有效 当 switch=在. 默认值 eo-客户端-ipcountry 是 使用 当 它 是 不 指定。",
								},
							},
						},
					},
					"upstream_follow_redirect_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Configuration 参数 对于 following redirects during 源站-pull. 此 参数 为必填项 当 名称 是 集合 到 upstreamfollowredirect。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "是否enable 源站-pull 到 follow redirection 配置. 值: 在: 启用; 关闭: disable。",
								},
								"max_times": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "最大redirects. 取值范围：1-5. 注意: 此 字段 为必填项 当 switch 是 在; 当 switch 是 关闭，此 字段 不是必填项 和 将 不 take effect 如果 filled。",
								},
							},
						},
					},
					"upstream_request_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Configuration 参数 对于 源站-pull 请求. 此 参数 为必填项 当 名称 是 集合 到 upstreamrequest。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"query_string": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "Query 字符串 配置. 可选 如果 不 提供，它 将 不 是 已配置。",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"switch": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "是否enable 源站-pull 请求 参数 查询 字符串. 值: 在: 启用; 关闭: disable。",
											},
											"action": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Query 字符串 模式 此 参数 为必填项 当 switch 是 在. 值: full: retain all; ignore: ignore all; includeCustom: retain partial 参数; excludeCustom: ignore partial 参数。",
											},
											"values": {
												Type:        schema.TypeList,
												Optional:    true,
												Description: "指定parameter 值. 此 参数 takes effect 仅 当 查询 字符串 模式 操作 是 includecustom 或 excludecustom，和 是 用于指定parameters 到 是 reserved 或 ignored. up 到 10 参数 是 支持。",
												Elem: &schema.Schema{
													Type: schema.TypeString,
												},
											},
										},
									},
								},
								"cookie": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "Cookie 配置. 可选 如果 不 提供，它 将 不 是 已配置。",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"switch": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "是否enable 源站-pull 请求 参数 cookie. 有效值：在: 启用; 关闭: disable。",
											},
											"action": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "Origin-Pull 请求 参数 cookie 模式 此 参数 为必填项 当 switch 是 在. 有效值：full: retain all; ignore: ignore all; includeCustom: retain partial 参数; excludeCustom: ignore partial 参数。",
											},
											"values": {
												Type:        schema.TypeList,
												Optional:    true,
												Description: "指定parameter 值. 此 参数 takes effect 仅 当 查询 字符串 模式 操作 是 includecustom 或 excludecustom，和 是 用于指定parameters 到 是 reserved 或 ignored. up 到 10 参数 是 支持。",
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
					"tls_config_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "SSL/TLS 安全 配置 参数. 此 参数 为必填项 当 名称 是 集合 到 tlsconfig。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"version": {
									Type:        schema.TypeList,
									Optional:    true,
									Description: "TLS 版本 在 least 一个 必须 是 指定. 如果 多个 versions 是 指定，they 必须 是 consecutive，e.g.，启用 tls1，1.1，1.2，和 1.3. 它 是 不 allowed 到 启用 仅 1 和 1.2 while disabling 1.1. 有效值：tlsv1: tlsv1 版本; `tlsv1.1`: tlsv1.1 版本; `tlsv1.2`: tlsv1.2 版本; `tlsv1.3`: tlsv1.3 版本",
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								"cipher_suite": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Cipher suite. 对于 detailed 信息，please refer 到 tls versions 和 cipher suites 描述，https://www.tencentcloud.com/document/product/1145/54154?has_map=1. 有效值：loose-v2023: loose-v2023 cipher suite; general-v2023: general-v2023 cipher suite; strict-v2023: strict-v2023 cipher suite。",
								},
							},
						},
					},
					"modify_origin_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Configuration 参数 对于 modifying 源站 服务器. 此 参数 为必填项 当 名称 是 集合 到 modifyorigin。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"origin_type": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "源站 类型 值: IPDomain: ipv4，ipv6，或 域名 名称 类型 源站 服务器; OriginGroup: 源站 服务器 组 类型 源站 服务器; LoadBalance: 云 load balancer (clb)，此 功能 是 在 beta 测试. 到 使用 它，please 提交 ticket 或 contact smart customer 服务; COS: tencent 云 COS 源站 服务器; AWSS3: all 对象 存储 源站 servers 该 support aws s3 协议",
								},
								"origin": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Origin 服务器 地址，其中 varies according 到 值 的 origintype: 当 origintype = ipdomain，fill 在 ipv4 地址， ipv6 地址，或 域名 名称; 当 origintype = cos，please fill 在 访问 域名 名称 COS 存储桶; 当 origintype = awss3，fill 在 访问 域名 名称 s3 存储桶; 当 origintype = origingroup，fill 在 源站 服务器 组 ID; 当 origintype = loadbalance，fill 在 云 load balancer 实例 ID 此 功能 是 currently 仅 可用 到 allowlist。",
								},
								"origin_protocol": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Origin-Pull 协议 配置. 此 参数 为必填项 当 origintype 是 ipdomain，origingroup，或 loadbalance. 有效值：Http: 使用 http 协议; Https: 使用 https 协议; Follow: follow 协议",
								},
								"http_origin_port": {
									Type:         schema.TypeInt,
									Optional:     true,
									ValidateFunc: tccommon.ValidateIntegerInRange(1, 65535),
									Description:  "Ports 对于 http 源站-pull requests. 取值范围：1-65535. 此 参数 takes effect 仅 当 源站-pull 协议 originprotocol 是 http 或 follow。",
								},
								"https_origin_port": {
									Type:         schema.TypeInt,
									Optional:     true,
									ValidateFunc: tccommon.ValidateIntegerInRange(1, 65535),
									Description:  "Ports 对于 https 源站-pull requests. 取值范围：1-65535. 此 参数 takes effect 仅 当 源站-pull 协议 originprotocol 是 https 或 follow。",
								},
								"private_access": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Whether 访问 到 私有 对象 存储 源站 服务器 是 allowed. 此 参数 是 有效 仅 当 源站 服务器 类型 origintype 是 COS 或 awss3. 有效值：在: 启用 私有 authentication; 关闭: disable 私有 authentication. 如果未指定， 默认值为 关闭。",
								},
								"private_parameters": {
									Type:        schema.TypeList,
									Optional:    true,
									MaxItems:    1,
									Description: "Private authentication 参数. 此 参数 是 有效 仅 当 origintype = awss3 和 privateaccess = 在。",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"access_key_id": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "Authentication 参数 访问 键 ID。",
											},
											"secret_access_key": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "Authentication 参数 secret 访问 键",
											},
											"signature_version": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "Authentication 版本 值: v2: v2 版本; v4: v4 版本",
											},
											"region": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "地域 的 存储桶",
											},
										},
									},
								},
							},
						},
					},
					"http_upstream_timeout_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Configuration 的 layer 7 源站 超时. 此 参数 为必填项 当 名称 是 httpupstreamtimeout。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"response_timeout": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "HTTP response 超时 （秒）。 取值范围：5-600。",
								},
							},
						},
					},
					"http_response_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "HTTP response 配置 参数. 此 参数 为必填项 当 名称 是 httpresponse。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"status_code": {
									Type:        schema.TypeInt,
									Optional:    true,
									Description: "Response 状态 代码 支持 2xx，4xx，5xx，excluding 499，514，101，301，302，303，509，520-599。",
								},
								"response_page": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Response 页面 ID。",
								},
							},
						},
					},
					"error_page_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Custom 错误 页面 配置 参数. 此 参数 为必填项 当 名称 是 errorpage。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"error_page_params": {
									Type:        schema.TypeList,
									Optional:    true,
									Description: "Custom 错误 页面 配置 列表。",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"status_code": {
												Type:        schema.TypeInt,
												Required:    true,
												Description: "状态 代码 支持 值 是 400，403，404，405，414，416，451，500，501，502，503，504。",
											},
											"redirect_url": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "Redirect URL requires full redirect 路径，such 作为 https://www.测试.com/错误html。",
											},
										},
									},
								},
							},
						},
					},
					"modify_response_header_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Modify http 节点 response 头部 配置 参数. 此 参数 为必填项 当 名称 是 modifyresponseheader。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"header_actions": {
									Type:        schema.TypeList,
									Optional:    true,
									Description: "HTTP 源站-pull 头部 规则 列表。",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"action": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "HTTP 头部 setting methods. 有效值：集合: sets 值 对于 existing 头部 参数; del: deletes 头部 参数; add: adds 头部 参数。",
											},
											"name": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "HTTP 头部 名称",
											},
											"value": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "HTTP 头部 值 此 参数 为必填项 当 操作 是 集合 到 集合 或 add; 它 为可选项 当 操作 是 集合 到 del。",
											},
										},
									},
								},
							},
						},
					},
					"modify_request_header_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Modify http 节点 请求 头部 配置 参数. 此 参数 为必填项 当 名称 是 modifyrequestheader。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"header_actions": {
									Type:        schema.TypeList,
									Optional:    true,
									Description: "列表 http 头部 setting 规则。",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"action": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "HTTP 头部 setting methods. 有效值：集合: sets 值 对于 existing 头部 参数; del: deletes 头部 参数; add: adds 头部 参数。",
											},
											"name": {
												Type:        schema.TypeString,
												Required:    true,
												Description: "HTTP 头部 名称",
											},
											"value": {
												Type:        schema.TypeString,
												Optional:    true,
												Description: "HTTP 头部 值 此 参数 为必填项 当 操作 是 集合 到 集合 或 add; 它 为可选项 当 操作 是 集合 到 del。",
											},
										},
									},
								},
							},
						},
					},
					"response_speed_limit_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Single 连接 download speed 限制 配置 参数. 此 参数 为必填项 当 名称 是 responsespeedlimit。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"mode": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Download 速率 限制 模式 有效值：LimitUponDownload: 速率 限制 throughout download process; LimitAfterSpecificBytesDownloaded: 速率 限制 after downloading 特定 bytes 在 full speed; LimitAfterSpecificSecondsDownloaded: start speed 限制 after downloading 在 full speed 对于 特定 时长。",
								},
								"max_speed": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "Rate-Limiting 值，在 kb/s. enter numerical 值 到 指定rate 限制",
								},
								"start_at": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "Rate-Limiting start 值，其中 可以 是 download 大小 或 指定 时长，在 kb 或 s. 此 参数 为必填项 当 模式 是 集合 到 limitafterspecificbytesdownloaded 或 limitafterspecificsecondsdownloaded. enter numerical 值 到 指定download 大小 或 时长。",
								},
							},
						},
					},
					"set_content_identifier_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "内容 identification 配置 参数. 此 参数 为必填项 当 名称 是 httpresponse。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"content_identifier": {
									Type:        schema.TypeString,
									Optional:    true,
									Description: "内容 identifier ID。",
								},
							},
						},
					},
					"content_compression_parameters": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "内容 压缩 配置 参数. 此 参数 为必填项 当 `名称` 参数 是 集合 到 `ContentCompression`. 此 参数 uses whitelist 函数; please contact Tencent Cloud engineers 如果 needed。",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"switch": {
									Type:        schema.TypeString,
									Required:    true,
									Description: "内容 压缩 配置 switch，possible 值 是: 在: 已启用; 关闭: 已禁用 当 Switch 是 集合 到 `在`，both Brotli 和 gzip 压缩 algorithms 将 是 支持。",
								},
							},
						},
					},
				},
			},
		},
	}

	if depth < 8 {
		schemaMap["sub_rules"] = &schema.Schema{
			Type:        schema.TypeList,
			Optional:    true,
			Description: "列表 sub-规则. 多个 规则 exist 在 此 列表 和 是 executed sequentially 从 top 到 bottom. note: subrules 和 actions 不能 both 是 空. currently，仅 一个 layer 的 subrules 是 支持。",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"branches": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Sub-规则 branch。",
						Elem: &schema.Resource{
							Schema: TencentTeoL7RuleBranchBasicInfo(depth + 1),
						},
					},
					"description": {
						Type:        schema.TypeList,
						Optional:    true,
						Description: "Rule comments。",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		}
	}

	return schemaMap
}

func resourceTencentCloudTeoL7AccRuleGetBranchs(rulesMap map[string]interface{}) []*teo.RuleBranch {
	ruleBranchs := []*teo.RuleBranch{}
	if v, ok := rulesMap["branches"]; ok {
		for _, item := range v.([]interface{}) {
			branchesMap := item.(map[string]interface{})
			ruleBranch := teov20220901.RuleBranch{}
			if v, ok := branchesMap["condition"].(string); ok && v != "" {
				ruleBranch.Condition = helper.String(v)
			}
			if v, ok := branchesMap["actions"]; ok {
				for _, item := range v.([]interface{}) {
					actionsMap := item.(map[string]interface{})
					ruleEngineAction := teov20220901.RuleEngineAction{}
					if v, ok := actionsMap["name"].(string); ok && v != "" {
						ruleEngineAction.Name = helper.String(v)
					}
					if cacheParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["cache_parameters"]); ok {
						cacheParameters := teov20220901.CacheParameters{}
						if followOriginMap, ok := helper.ConvertInterfacesHeadToMap(cacheParametersMap["follow_origin"]); ok {
							followOrigin := teov20220901.FollowOrigin{}
							if v, ok := followOriginMap["switch"].(string); ok && v != "" {
								followOrigin.Switch = helper.String(v)
							}
							if v, ok := followOriginMap["default_cache"].(string); ok && v != "" {
								followOrigin.DefaultCache = helper.String(v)
							}
							if v, ok := followOriginMap["default_cache_strategy"].(string); ok && v != "" {
								followOrigin.DefaultCacheStrategy = helper.String(v)
							}
							if v, ok := followOriginMap["default_cache_time"].(int); ok {
								followOrigin.DefaultCacheTime = helper.IntInt64(v)
							}
							cacheParameters.FollowOrigin = &followOrigin
						}
						if noCacheMap, ok := helper.ConvertInterfacesHeadToMap(cacheParametersMap["no_cache"]); ok {
							noCache := teov20220901.NoCache{}
							if v, ok := noCacheMap["switch"].(string); ok && v != "" {
								noCache.Switch = helper.String(v)
							}
							cacheParameters.NoCache = &noCache
						}
						if customTimeMap, ok := helper.ConvertInterfacesHeadToMap(cacheParametersMap["custom_time"]); ok {
							customTime := teov20220901.CustomTime{}
							if v, ok := customTimeMap["switch"].(string); ok && v != "" {
								customTime.Switch = helper.String(v)
							}
							if v, ok := customTimeMap["ignore_cache_control"].(string); ok && v != "" {
								customTime.IgnoreCacheControl = helper.String(v)
							}
							if v, ok := customTimeMap["cache_time"].(int); ok {
								customTime.CacheTime = helper.IntInt64(v)
							}
							cacheParameters.CustomTime = &customTime
						}
						ruleEngineAction.CacheParameters = &cacheParameters
					}
					if cacheKeyParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["cache_key_parameters"]); ok {
						cacheKeyParameters := teov20220901.CacheKeyParameters{}
						if v, ok := cacheKeyParametersMap["full_url_cache"].(string); ok && v != "" {
							cacheKeyParameters.FullURLCache = helper.String(v)
						}
						if queryStringMap, ok := helper.ConvertInterfacesHeadToMap(cacheKeyParametersMap["query_string"]); ok {
							cacheKeyQueryString := teov20220901.CacheKeyQueryString{}
							if v, ok := queryStringMap["switch"].(string); ok && v != "" {
								cacheKeyQueryString.Switch = helper.String(v)
							}
							if v, ok := queryStringMap["action"].(string); ok && v != "" {
								cacheKeyQueryString.Action = helper.String(v)
							}
							if v, ok := queryStringMap["values"]; ok {
								valuesSet := v.([]interface{})
								for i := range valuesSet {
									values := valuesSet[i].(string)
									cacheKeyQueryString.Values = append(cacheKeyQueryString.Values, helper.String(values))
								}
							}
							cacheKeyParameters.QueryString = &cacheKeyQueryString
						}
						if v, ok := cacheKeyParametersMap["ignore_case"].(string); ok && v != "" {
							cacheKeyParameters.IgnoreCase = helper.String(v)
						}
						if headerMap, ok := helper.ConvertInterfacesHeadToMap(cacheKeyParametersMap["header"]); ok {
							cacheKeyHeader := teov20220901.CacheKeyHeader{}
							if v, ok := headerMap["switch"].(string); ok && v != "" {
								cacheKeyHeader.Switch = helper.String(v)
							}
							if v, ok := headerMap["values"]; ok {
								valuesSet := v.([]interface{})
								for i := range valuesSet {
									values := valuesSet[i].(string)
									cacheKeyHeader.Values = append(cacheKeyHeader.Values, helper.String(values))
								}
							}
							cacheKeyParameters.Header = &cacheKeyHeader
						}
						if v, ok := cacheKeyParametersMap["scheme"].(string); ok && v != "" {
							cacheKeyParameters.Scheme = helper.String(v)
						}
						if cookieMap, ok := helper.ConvertInterfacesHeadToMap(cacheKeyParametersMap["cookie"]); ok {
							cacheKeyCookie := teov20220901.CacheKeyCookie{}
							if v, ok := cookieMap["switch"].(string); ok && v != "" {
								cacheKeyCookie.Switch = helper.String(v)
							}
							if v, ok := cookieMap["action"].(string); ok && v != "" {
								cacheKeyCookie.Action = helper.String(v)
							}
							if v, ok := cookieMap["values"]; ok {
								valuesSet := v.([]interface{})
								for i := range valuesSet {
									values := valuesSet[i].(string)
									cacheKeyCookie.Values = append(cacheKeyCookie.Values, helper.String(values))
								}
							}
							cacheKeyParameters.Cookie = &cacheKeyCookie
						}
						ruleEngineAction.CacheKeyParameters = &cacheKeyParameters
					}
					if cachePrefreshParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["cache_prefresh_parameters"]); ok {
						cachePrefreshParameters := teov20220901.CachePrefreshParameters{}
						if v, ok := cachePrefreshParametersMap["switch"].(string); ok && v != "" {
							cachePrefreshParameters.Switch = helper.String(v)
						}
						if v, ok := cachePrefreshParametersMap["cache_time_percent"].(int); ok {
							cachePrefreshParameters.CacheTimePercent = helper.IntInt64(v)
						}
						ruleEngineAction.CachePrefreshParameters = &cachePrefreshParameters
					}
					if accessURLRedirectParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["access_url_redirect_parameters"]); ok {
						accessURLRedirectParameters := teov20220901.AccessURLRedirectParameters{}
						if v, ok := accessURLRedirectParametersMap["status_code"].(int); ok {
							accessURLRedirectParameters.StatusCode = helper.IntInt64(v)
						}
						if v, ok := accessURLRedirectParametersMap["protocol"].(string); ok && v != "" {
							accessURLRedirectParameters.Protocol = helper.String(v)
						}
						if hostNameMap, ok := helper.ConvertInterfacesHeadToMap(accessURLRedirectParametersMap["host_name"]); ok {
							hostName := teov20220901.HostName{}
							if v, ok := hostNameMap["action"].(string); ok && v != "" {
								hostName.Action = helper.String(v)
							}
							if v, ok := hostNameMap["value"].(string); ok && v != "" {
								hostName.Value = helper.String(v)
							}
							accessURLRedirectParameters.HostName = &hostName
						}
						if uRLPathMap, ok := helper.ConvertInterfacesHeadToMap(accessURLRedirectParametersMap["url_path"]); ok {
							uRLPath := teov20220901.URLPath{}
							if v, ok := uRLPathMap["action"].(string); ok && v != "" {
								uRLPath.Action = helper.String(v)
							}
							if v, ok := uRLPathMap["regex"].(string); ok && v != "" {
								uRLPath.Regex = helper.String(v)
							}
							if v, ok := uRLPathMap["value"].(string); ok && v != "" {
								uRLPath.Value = helper.String(v)
							}
							accessURLRedirectParameters.URLPath = &uRLPath
						}
						if queryStringMap, ok := helper.ConvertInterfacesHeadToMap(accessURLRedirectParametersMap["query_string"]); ok {
							accessURLRedirectQueryString := teov20220901.AccessURLRedirectQueryString{}
							if v, ok := queryStringMap["action"].(string); ok && v != "" {
								accessURLRedirectQueryString.Action = helper.String(v)
							}
							accessURLRedirectParameters.QueryString = &accessURLRedirectQueryString
						}
						ruleEngineAction.AccessURLRedirectParameters = &accessURLRedirectParameters
					}
					if upstreamURLRewriteParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["upstream_url_rewrite_parameters"]); ok {
						upstreamURLRewriteParameters := teov20220901.UpstreamURLRewriteParameters{}
						if v, ok := upstreamURLRewriteParametersMap["type"].(string); ok && v != "" {
							upstreamURLRewriteParameters.Type = helper.String(v)
						}
						if v, ok := upstreamURLRewriteParametersMap["action"].(string); ok && v != "" {
							upstreamURLRewriteParameters.Action = helper.String(v)
						}
						if v, ok := upstreamURLRewriteParametersMap["value"].(string); ok && v != "" {
							upstreamURLRewriteParameters.Value = helper.String(v)
						}
						if v, ok := upstreamURLRewriteParametersMap["regex"].(string); ok && v != "" {
							upstreamURLRewriteParameters.Regex = helper.String(v)
						}
						ruleEngineAction.UpstreamURLRewriteParameters = &upstreamURLRewriteParameters
					}
					if qUICParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["quic_parameters"]); ok {
						qUICParameters := teov20220901.QUICParameters{}
						if v, ok := qUICParametersMap["switch"].(string); ok && v != "" {
							qUICParameters.Switch = helper.String(v)
						}
						ruleEngineAction.QUICParameters = &qUICParameters
					}
					if webSocketParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["web_socket_parameters"]); ok {
						webSocketParameters := teov20220901.WebSocketParameters{}
						if v, ok := webSocketParametersMap["switch"].(string); ok && v != "" {
							webSocketParameters.Switch = helper.String(v)
						}
						if v, ok := webSocketParametersMap["timeout"].(int); ok {
							webSocketParameters.Timeout = helper.IntInt64(v)
						}
						ruleEngineAction.WebSocketParameters = &webSocketParameters
					}
					if authenticationParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["authentication_parameters"]); ok {
						authenticationParameters := teov20220901.AuthenticationParameters{}
						if v, ok := authenticationParametersMap["auth_type"].(string); ok && v != "" {
							authenticationParameters.AuthType = helper.String(v)
						}
						if v, ok := authenticationParametersMap["secret_key"].(string); ok && v != "" {
							authenticationParameters.SecretKey = helper.String(v)
						}
						if v, ok := authenticationParametersMap["timeout"].(int); ok {
							authenticationParameters.Timeout = helper.IntInt64(v)
						}
						if v, ok := authenticationParametersMap["backup_secret_key"].(string); ok && v != "" {
							authenticationParameters.BackupSecretKey = helper.String(v)
						}
						if v, ok := authenticationParametersMap["auth_param"].(string); ok && v != "" {
							authenticationParameters.AuthParam = helper.String(v)
						}
						if v, ok := authenticationParametersMap["time_param"].(string); ok && v != "" {
							authenticationParameters.TimeParam = helper.String(v)
						}
						if v, ok := authenticationParametersMap["time_format"].(string); ok && v != "" {
							authenticationParameters.TimeFormat = helper.String(v)
						}
						ruleEngineAction.AuthenticationParameters = &authenticationParameters
					}
					if maxAgeParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["max_age_parameters"]); ok {
						maxAgeParameters := teov20220901.MaxAgeParameters{}
						if v, ok := maxAgeParametersMap["follow_origin"].(string); ok && v != "" {
							maxAgeParameters.FollowOrigin = helper.String(v)
						}
						if v, ok := maxAgeParametersMap["cache_time"].(int); ok {
							maxAgeParameters.CacheTime = helper.IntInt64(v)
						}
						ruleEngineAction.MaxAgeParameters = &maxAgeParameters
					}
					if statusCodeCacheParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["status_code_cache_parameters"]); ok {
						statusCodeCacheParameters := teov20220901.StatusCodeCacheParameters{}
						if v, ok := statusCodeCacheParametersMap["status_code_cache_params"]; ok {
							for _, item := range v.([]interface{}) {
								statusCodeCacheParamsMap := item.(map[string]interface{})
								statusCodeCacheParam := teov20220901.StatusCodeCacheParam{}
								if v, ok := statusCodeCacheParamsMap["status_code"].(int); ok {
									statusCodeCacheParam.StatusCode = helper.IntInt64(v)
								}
								if v, ok := statusCodeCacheParamsMap["cache_time"].(int); ok {
									statusCodeCacheParam.CacheTime = helper.IntInt64(v)
								}
								statusCodeCacheParameters.StatusCodeCacheParams = append(statusCodeCacheParameters.StatusCodeCacheParams, &statusCodeCacheParam)
							}
						}
						ruleEngineAction.StatusCodeCacheParameters = &statusCodeCacheParameters
					}
					if offlineCacheParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["offline_cache_parameters"]); ok {
						offlineCacheParameters := teov20220901.OfflineCacheParameters{}
						if v, ok := offlineCacheParametersMap["switch"].(string); ok && v != "" {
							offlineCacheParameters.Switch = helper.String(v)
						}
						ruleEngineAction.OfflineCacheParameters = &offlineCacheParameters
					}
					if smartRoutingParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["smart_routing_parameters"]); ok {
						smartRoutingParameters := teov20220901.SmartRoutingParameters{}
						if v, ok := smartRoutingParametersMap["switch"].(string); ok && v != "" {
							smartRoutingParameters.Switch = helper.String(v)
						}
						ruleEngineAction.SmartRoutingParameters = &smartRoutingParameters
					}
					if rangeOriginPullParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["range_origin_pull_parameters"]); ok {
						rangeOriginPullParameters := teov20220901.RangeOriginPullParameters{}
						if v, ok := rangeOriginPullParametersMap["switch"].(string); ok && v != "" {
							rangeOriginPullParameters.Switch = helper.String(v)
						}
						ruleEngineAction.RangeOriginPullParameters = &rangeOriginPullParameters
					}
					if upstreamHTTP2ParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["upstream_http2_parameters"]); ok {
						upstreamHTTP2Parameters := teov20220901.UpstreamHTTP2Parameters{}
						if v, ok := upstreamHTTP2ParametersMap["switch"].(string); ok && v != "" {
							upstreamHTTP2Parameters.Switch = helper.String(v)
						}
						ruleEngineAction.UpstreamHTTP2Parameters = &upstreamHTTP2Parameters
					}
					if hostHeaderParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["host_header_parameters"]); ok {
						hostHeaderParameters := teov20220901.HostHeaderParameters{}
						if v, ok := hostHeaderParametersMap["action"].(string); ok && v != "" {
							hostHeaderParameters.Action = helper.String(v)
						}
						if v, ok := hostHeaderParametersMap["server_name"].(string); ok && v != "" {
							hostHeaderParameters.ServerName = helper.String(v)
						}
						ruleEngineAction.HostHeaderParameters = &hostHeaderParameters
					}
					if forceRedirectHTTPSParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["force_redirect_https_parameters"]); ok {
						forceRedirectHTTPSParameters := teov20220901.ForceRedirectHTTPSParameters{}
						if v, ok := forceRedirectHTTPSParametersMap["switch"].(string); ok && v != "" {
							forceRedirectHTTPSParameters.Switch = helper.String(v)
						}
						if v, ok := forceRedirectHTTPSParametersMap["redirect_status_code"].(int); ok {
							forceRedirectHTTPSParameters.RedirectStatusCode = helper.IntInt64(v)
						}
						ruleEngineAction.ForceRedirectHTTPSParameters = &forceRedirectHTTPSParameters
					}
					if originPullProtocolParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["origin_pull_protocol_parameters"]); ok {
						originPullProtocolParameters := teov20220901.OriginPullProtocolParameters{}
						if v, ok := originPullProtocolParametersMap["protocol"].(string); ok && v != "" {
							originPullProtocolParameters.Protocol = helper.String(v)
						}
						ruleEngineAction.OriginPullProtocolParameters = &originPullProtocolParameters
					}
					if compressionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["compression_parameters"]); ok {
						compressionParameters := teov20220901.CompressionParameters{}
						if v, ok := compressionParametersMap["switch"].(string); ok && v != "" {
							compressionParameters.Switch = helper.String(v)
						}
						if v, ok := compressionParametersMap["algorithms"]; ok {
							algorithmsSet := v.([]interface{})
							for i := range algorithmsSet {
								algorithms := algorithmsSet[i].(string)
								compressionParameters.Algorithms = append(compressionParameters.Algorithms, helper.String(algorithms))
							}
						}
						ruleEngineAction.CompressionParameters = &compressionParameters
					}
					if hSTSParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["hsts_parameters"]); ok {
						hSTSParameters := teov20220901.HSTSParameters{}
						if v, ok := hSTSParametersMap["switch"].(string); ok && v != "" {
							hSTSParameters.Switch = helper.String(v)
						}
						if v, ok := hSTSParametersMap["timeout"].(int); ok {
							hSTSParameters.Timeout = helper.IntInt64(v)
						}
						if v, ok := hSTSParametersMap["include_sub_domains"].(string); ok && v != "" {
							hSTSParameters.IncludeSubDomains = helper.String(v)
						}
						if v, ok := hSTSParametersMap["preload"].(string); ok && v != "" {
							hSTSParameters.Preload = helper.String(v)
						}
						ruleEngineAction.HSTSParameters = &hSTSParameters
					}
					if clientIPHeaderParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["client_ip_header_parameters"]); ok {
						clientIPHeaderParameters := teov20220901.ClientIPHeaderParameters{}
						if v, ok := clientIPHeaderParametersMap["switch"].(string); ok && v != "" {
							clientIPHeaderParameters.Switch = helper.String(v)
						}
						if v, ok := clientIPHeaderParametersMap["header_name"].(string); ok && v != "" {
							clientIPHeaderParameters.HeaderName = helper.String(v)
						}
						ruleEngineAction.ClientIPHeaderParameters = &clientIPHeaderParameters
					}
					if oCSPStaplingParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["ocsp_stapling_parameters"]); ok {
						oCSPStaplingParameters := teov20220901.OCSPStaplingParameters{}
						if v, ok := oCSPStaplingParametersMap["switch"].(string); ok && v != "" {
							oCSPStaplingParameters.Switch = helper.String(v)
						}
						ruleEngineAction.OCSPStaplingParameters = &oCSPStaplingParameters
					}
					if hTTP2ParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["http2_parameters"]); ok {
						hTTP2Parameters := teov20220901.HTTP2Parameters{}
						if v, ok := hTTP2ParametersMap["switch"].(string); ok && v != "" {
							hTTP2Parameters.Switch = helper.String(v)
						}
						ruleEngineAction.HTTP2Parameters = &hTTP2Parameters
					}
					if postMaxSizeParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["post_max_size_parameters"]); ok {
						postMaxSizeParameters := teov20220901.PostMaxSizeParameters{}
						if v, ok := postMaxSizeParametersMap["switch"].(string); ok && v != "" {
							postMaxSizeParameters.Switch = helper.String(v)
						}
						if v, ok := postMaxSizeParametersMap["max_size"].(int); ok {
							postMaxSizeParameters.MaxSize = helper.IntInt64(v)
						}
						ruleEngineAction.PostMaxSizeParameters = &postMaxSizeParameters
					}
					if clientIPCountryParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["client_ip_country_parameters"]); ok {
						clientIPCountryParameters := teov20220901.ClientIPCountryParameters{}
						if v, ok := clientIPCountryParametersMap["switch"].(string); ok && v != "" {
							clientIPCountryParameters.Switch = helper.String(v)
						}
						if v, ok := clientIPCountryParametersMap["header_name"].(string); ok && v != "" {
							clientIPCountryParameters.HeaderName = helper.String(v)
						}
						ruleEngineAction.ClientIPCountryParameters = &clientIPCountryParameters
					}
					if upstreamFollowRedirectParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["upstream_follow_redirect_parameters"]); ok {
						upstreamFollowRedirectParameters := teov20220901.UpstreamFollowRedirectParameters{}
						if v, ok := upstreamFollowRedirectParametersMap["switch"].(string); ok && v != "" {
							upstreamFollowRedirectParameters.Switch = helper.String(v)
						}
						if v, ok := upstreamFollowRedirectParametersMap["max_times"].(int); ok {
							upstreamFollowRedirectParameters.MaxTimes = helper.IntInt64(v)
						}
						ruleEngineAction.UpstreamFollowRedirectParameters = &upstreamFollowRedirectParameters
					}
					if upstreamRequestParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["upstream_request_parameters"]); ok {
						upstreamRequestParameters := teov20220901.UpstreamRequestParameters{}
						if queryStringMap, ok := helper.ConvertInterfacesHeadToMap(upstreamRequestParametersMap["query_string"]); ok {
							upstreamRequestQueryString := teov20220901.UpstreamRequestQueryString{}
							if v, ok := queryStringMap["switch"].(string); ok && v != "" {
								upstreamRequestQueryString.Switch = helper.String(v)
							}
							if v, ok := queryStringMap["action"].(string); ok && v != "" {
								upstreamRequestQueryString.Action = helper.String(v)
							}
							if v, ok := queryStringMap["values"]; ok {
								valuesSet := v.([]interface{})
								for i := range valuesSet {
									values := valuesSet[i].(string)
									upstreamRequestQueryString.Values = append(upstreamRequestQueryString.Values, helper.String(values))
								}
							}
							upstreamRequestParameters.QueryString = &upstreamRequestQueryString
						}
						if cookieMap, ok := helper.ConvertInterfacesHeadToMap(upstreamRequestParametersMap["cookie"]); ok {
							upstreamRequestCookie := teov20220901.UpstreamRequestCookie{}
							if v, ok := cookieMap["switch"].(string); ok && v != "" {
								upstreamRequestCookie.Switch = helper.String(v)
							}
							if v, ok := cookieMap["action"].(string); ok && v != "" {
								upstreamRequestCookie.Action = helper.String(v)
							}
							if v, ok := cookieMap["values"]; ok {
								valuesSet := v.([]interface{})
								for i := range valuesSet {
									values := valuesSet[i].(string)
									upstreamRequestCookie.Values = append(upstreamRequestCookie.Values, helper.String(values))
								}
							}
							upstreamRequestParameters.Cookie = &upstreamRequestCookie
						}
						ruleEngineAction.UpstreamRequestParameters = &upstreamRequestParameters
					}
					if tLSConfigParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["tls_config_parameters"]); ok {
						tLSConfigParameters := teov20220901.TLSConfigParameters{}
						if v, ok := tLSConfigParametersMap["version"]; ok {
							versionSet := v.([]interface{})
							for i := range versionSet {
								version := versionSet[i].(string)
								tLSConfigParameters.Version = append(tLSConfigParameters.Version, helper.String(version))
							}
						}
						if v, ok := tLSConfigParametersMap["cipher_suite"].(string); ok && v != "" {
							tLSConfigParameters.CipherSuite = helper.String(v)
						}
						ruleEngineAction.TLSConfigParameters = &tLSConfigParameters
					}
					if modifyOriginParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["modify_origin_parameters"]); ok {
						modifyOriginParameters := teov20220901.ModifyOriginParameters{}
						if v, ok := modifyOriginParametersMap["origin_type"].(string); ok && v != "" {
							modifyOriginParameters.OriginType = helper.String(v)
						}
						if v, ok := modifyOriginParametersMap["origin"].(string); ok && v != "" {
							modifyOriginParameters.Origin = helper.String(v)
						}
						if v, ok := modifyOriginParametersMap["origin_protocol"].(string); ok && v != "" {
							modifyOriginParameters.OriginProtocol = helper.String(v)
						}
						if v, ok := modifyOriginParametersMap["http_origin_port"].(int); ok && v != 0 {
							modifyOriginParameters.HTTPOriginPort = helper.IntInt64(v)
						}
						if v, ok := modifyOriginParametersMap["https_origin_port"].(int); ok && v != 0 {
							modifyOriginParameters.HTTPSOriginPort = helper.IntInt64(v)
						}
						if v, ok := modifyOriginParametersMap["private_access"].(string); ok && v != "" {
							modifyOriginParameters.PrivateAccess = helper.String(v)
						}
						if privateParametersMap, ok := helper.ConvertInterfacesHeadToMap(modifyOriginParametersMap["private_parameters"]); ok {
							originPrivateParameters := teov20220901.OriginPrivateParameters{}
							if v, ok := privateParametersMap["access_key_id"].(string); ok && v != "" {
								originPrivateParameters.AccessKeyId = helper.String(v)
							}
							if v, ok := privateParametersMap["secret_access_key"].(string); ok && v != "" {
								originPrivateParameters.SecretAccessKey = helper.String(v)
							}
							if v, ok := privateParametersMap["signature_version"].(string); ok && v != "" {
								originPrivateParameters.SignatureVersion = helper.String(v)
							}
							if v, ok := privateParametersMap["region"].(string); ok && v != "" {
								originPrivateParameters.Region = helper.String(v)
							}
							modifyOriginParameters.PrivateParameters = &originPrivateParameters
						}
						ruleEngineAction.ModifyOriginParameters = &modifyOriginParameters
					}
					if hTTPUpstreamTimeoutParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["http_upstream_timeout_parameters"]); ok {
						hTTPUpstreamTimeoutParameters := teov20220901.HTTPUpstreamTimeoutParameters{}
						if v, ok := hTTPUpstreamTimeoutParametersMap["response_timeout"].(int); ok {
							hTTPUpstreamTimeoutParameters.ResponseTimeout = helper.IntInt64(v)
						}
						ruleEngineAction.HTTPUpstreamTimeoutParameters = &hTTPUpstreamTimeoutParameters
					}
					if httpResponseParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["http_response_parameters"]); ok {
						hTTPResponseParameters := teov20220901.HTTPResponseParameters{}
						if v, ok := httpResponseParametersMap["status_code"].(int); ok {
							hTTPResponseParameters.StatusCode = helper.IntInt64(v)
						}
						if v, ok := httpResponseParametersMap["response_page"].(string); ok && v != "" {
							hTTPResponseParameters.ResponsePage = helper.String(v)
						}
						ruleEngineAction.HttpResponseParameters = &hTTPResponseParameters
					}
					if errorPageParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["error_page_parameters"]); ok {
						errorPageParameters := teov20220901.ErrorPageParameters{}
						if v, ok := errorPageParametersMap["error_page_params"]; ok {
							for _, item := range v.([]interface{}) {
								errorPageParamsMap := item.(map[string]interface{})
								errorPage := teov20220901.ErrorPage{}
								if v, ok := errorPageParamsMap["status_code"].(int); ok {
									errorPage.StatusCode = helper.IntInt64(v)
								}
								if v, ok := errorPageParamsMap["redirect_url"].(string); ok && v != "" {
									errorPage.RedirectURL = helper.String(v)
								}
								errorPageParameters.ErrorPageParams = append(errorPageParameters.ErrorPageParams, &errorPage)
							}
						}
						ruleEngineAction.ErrorPageParameters = &errorPageParameters
					}
					if modifyResponseHeaderParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["modify_response_header_parameters"]); ok {
						modifyResponseHeaderParameters := teov20220901.ModifyResponseHeaderParameters{}
						if v, ok := modifyResponseHeaderParametersMap["header_actions"]; ok {
							for _, item := range v.([]interface{}) {
								headerActionsMap := item.(map[string]interface{})
								headerAction := teov20220901.HeaderAction{}
								if v, ok := headerActionsMap["action"].(string); ok && v != "" {
									headerAction.Action = helper.String(v)
								}
								if v, ok := headerActionsMap["name"].(string); ok && v != "" {
									headerAction.Name = helper.String(v)
								}
								if v, ok := headerActionsMap["value"].(string); ok && v != "" {
									headerAction.Value = helper.String(v)
								}
								modifyResponseHeaderParameters.HeaderActions = append(modifyResponseHeaderParameters.HeaderActions, &headerAction)
							}
						}
						ruleEngineAction.ModifyResponseHeaderParameters = &modifyResponseHeaderParameters
					}
					if modifyRequestHeaderParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["modify_request_header_parameters"]); ok {
						modifyRequestHeaderParameters := teov20220901.ModifyRequestHeaderParameters{}
						if v, ok := modifyRequestHeaderParametersMap["header_actions"]; ok {
							for _, item := range v.([]interface{}) {
								headerActionsMap := item.(map[string]interface{})
								headerAction := teov20220901.HeaderAction{}
								if v, ok := headerActionsMap["action"].(string); ok && v != "" {
									headerAction.Action = helper.String(v)
								}
								if v, ok := headerActionsMap["name"].(string); ok && v != "" {
									headerAction.Name = helper.String(v)
								}
								if v, ok := headerActionsMap["value"].(string); ok && v != "" {
									headerAction.Value = helper.String(v)
								}
								modifyRequestHeaderParameters.HeaderActions = append(modifyRequestHeaderParameters.HeaderActions, &headerAction)
							}
						}
						ruleEngineAction.ModifyRequestHeaderParameters = &modifyRequestHeaderParameters
					}
					if responseSpeedLimitParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["response_speed_limit_parameters"]); ok {
						responseSpeedLimitParameters := teov20220901.ResponseSpeedLimitParameters{}
						if v, ok := responseSpeedLimitParametersMap["mode"].(string); ok && v != "" {
							responseSpeedLimitParameters.Mode = helper.String(v)
						}
						if v, ok := responseSpeedLimitParametersMap["max_speed"].(string); ok && v != "" {
							responseSpeedLimitParameters.MaxSpeed = helper.String(v)
						}
						if v, ok := responseSpeedLimitParametersMap["start_at"].(string); ok && v != "" {
							responseSpeedLimitParameters.StartAt = helper.String(v)
						}
						ruleEngineAction.ResponseSpeedLimitParameters = &responseSpeedLimitParameters
					}
					if setContentIdentifierParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["set_content_identifier_parameters"]); ok {
						setContentIdentifierParameters := teov20220901.SetContentIdentifierParameters{}
						if v, ok := setContentIdentifierParametersMap["content_identifier"].(string); ok && v != "" {
							setContentIdentifierParameters.ContentIdentifier = helper.String(v)
						}
						ruleEngineAction.SetContentIdentifierParameters = &setContentIdentifierParameters
					}
					if contentCompressionParametersMap, ok := helper.ConvertInterfacesHeadToMap(actionsMap["content_compression_parameters"]); ok {
						contentCompressionParameters := teov20220901.ContentCompressionParameters{}
						if v, ok := contentCompressionParametersMap["switch"].(string); ok && v != "" {
							contentCompressionParameters.Switch = helper.String(v)
						}
						ruleEngineAction.ContentCompressionParameters = &contentCompressionParameters
					}
					ruleBranch.Actions = append(ruleBranch.Actions, &ruleEngineAction)
				}
			}

			if v, ok := branchesMap["sub_rules"]; ok {
				for _, item := range v.([]interface{}) {
					subRulesMap := item.(map[string]interface{})
					ruleEngineSubRule := teov20220901.RuleEngineSubRule{}
					if _, ok := subRulesMap["branches"]; ok {
						branchs := resourceTencentCloudTeoL7AccRuleGetBranchs(subRulesMap)
						ruleEngineSubRule.Branches = branchs
					}
					if v, ok := subRulesMap["description"]; ok {
						descriptionSet := v.([]interface{})
						for i := range descriptionSet {
							description := descriptionSet[i].(string)
							ruleEngineSubRule.Description = append(ruleEngineSubRule.Description, helper.String(description))
						}
					}
					ruleBranch.SubRules = append(ruleBranch.SubRules, &ruleEngineSubRule)
				}
			}
			ruleBranchs = append(ruleBranchs, &ruleBranch)
		}
	}

	return ruleBranchs
}

func resourceTencentCloudTeoL7AccRuleSetBranchs(ruleBranches []*teo.RuleBranch) []map[string]interface{} {
	branchesList := make([]map[string]interface{}, 0, len(ruleBranches))
	if len(ruleBranches) > 0 {
		for _, branches := range ruleBranches {
			branchesMap := map[string]interface{}{}

			if branches.Condition != nil {
				branchesMap["condition"] = branches.Condition
			}

			actionsList := make([]map[string]interface{}, 0, len(branches.Actions))
			if branches.Actions != nil {
				for _, actions := range branches.Actions {
					actionsMap := map[string]interface{}{}

					if actions.Name != nil {
						actionsMap["name"] = actions.Name
					}

					cacheParametersMap := map[string]interface{}{}

					if actions.CacheParameters != nil {
						followOriginMap := map[string]interface{}{}

						if actions.CacheParameters.FollowOrigin != nil {
							if actions.CacheParameters.FollowOrigin.Switch != nil {
								followOriginMap["switch"] = actions.CacheParameters.FollowOrigin.Switch
							}

							if actions.CacheParameters.FollowOrigin.DefaultCache != nil {
								followOriginMap["default_cache"] = actions.CacheParameters.FollowOrigin.DefaultCache
							}

							if actions.CacheParameters.FollowOrigin.DefaultCacheStrategy != nil {
								followOriginMap["default_cache_strategy"] = actions.CacheParameters.FollowOrigin.DefaultCacheStrategy
							}

							if actions.CacheParameters.FollowOrigin.DefaultCacheTime != nil {
								followOriginMap["default_cache_time"] = actions.CacheParameters.FollowOrigin.DefaultCacheTime
							}

							cacheParametersMap["follow_origin"] = []interface{}{followOriginMap}
						}

						noCacheMap := map[string]interface{}{}

						if actions.CacheParameters.NoCache != nil {
							if actions.CacheParameters.NoCache.Switch != nil {
								noCacheMap["switch"] = actions.CacheParameters.NoCache.Switch
							}

							cacheParametersMap["no_cache"] = []interface{}{noCacheMap}
						}

						customTimeMap := map[string]interface{}{}

						if actions.CacheParameters.CustomTime != nil {
							if actions.CacheParameters.CustomTime.Switch != nil {
								customTimeMap["switch"] = actions.CacheParameters.CustomTime.Switch
							}

							if actions.CacheParameters.CustomTime.IgnoreCacheControl != nil {
								customTimeMap["ignore_cache_control"] = actions.CacheParameters.CustomTime.IgnoreCacheControl
							}

							if actions.CacheParameters.CustomTime.CacheTime != nil {
								customTimeMap["cache_time"] = actions.CacheParameters.CustomTime.CacheTime
							}

							cacheParametersMap["custom_time"] = []interface{}{customTimeMap}
						}

						actionsMap["cache_parameters"] = []interface{}{cacheParametersMap}
					}

					cacheKeyParametersMap := map[string]interface{}{}

					if actions.CacheKeyParameters != nil {
						if actions.CacheKeyParameters.FullURLCache != nil {
							cacheKeyParametersMap["full_url_cache"] = actions.CacheKeyParameters.FullURLCache
						}

						queryStringMap := map[string]interface{}{}

						if actions.CacheKeyParameters.QueryString != nil {
							if actions.CacheKeyParameters.QueryString.Switch != nil {
								queryStringMap["switch"] = actions.CacheKeyParameters.QueryString.Switch
							}

							if actions.CacheKeyParameters.QueryString.Action != nil {
								queryStringMap["action"] = actions.CacheKeyParameters.QueryString.Action
							}

							if actions.CacheKeyParameters.QueryString.Values != nil {
								queryStringMap["values"] = actions.CacheKeyParameters.QueryString.Values
							}

							cacheKeyParametersMap["query_string"] = []interface{}{queryStringMap}
						}

						if actions.CacheKeyParameters.IgnoreCase != nil {
							cacheKeyParametersMap["ignore_case"] = actions.CacheKeyParameters.IgnoreCase
						}

						headerMap := map[string]interface{}{}

						if actions.CacheKeyParameters.Header != nil {
							if actions.CacheKeyParameters.Header.Switch != nil {
								headerMap["switch"] = actions.CacheKeyParameters.Header.Switch
							}

							if actions.CacheKeyParameters.Header.Values != nil {
								headerMap["values"] = actions.CacheKeyParameters.Header.Values
							}

							cacheKeyParametersMap["header"] = []interface{}{headerMap}
						}

						if actions.CacheKeyParameters.Scheme != nil {
							cacheKeyParametersMap["scheme"] = actions.CacheKeyParameters.Scheme
						}

						cookieMap := map[string]interface{}{}

						if actions.CacheKeyParameters.Cookie != nil {
							if actions.CacheKeyParameters.Cookie.Switch != nil {
								cookieMap["switch"] = actions.CacheKeyParameters.Cookie.Switch
							}

							if actions.CacheKeyParameters.Cookie.Action != nil {
								cookieMap["action"] = actions.CacheKeyParameters.Cookie.Action
							}

							if actions.CacheKeyParameters.Cookie.Values != nil {
								cookieMap["values"] = actions.CacheKeyParameters.Cookie.Values
							}

							cacheKeyParametersMap["cookie"] = []interface{}{cookieMap}
						}

						actionsMap["cache_key_parameters"] = []interface{}{cacheKeyParametersMap}
					}

					cachePrefreshParametersMap := map[string]interface{}{}

					if actions.CachePrefreshParameters != nil {
						if actions.CachePrefreshParameters.Switch != nil {
							cachePrefreshParametersMap["switch"] = actions.CachePrefreshParameters.Switch
						}

						if actions.CachePrefreshParameters.CacheTimePercent != nil {
							cachePrefreshParametersMap["cache_time_percent"] = actions.CachePrefreshParameters.CacheTimePercent
						}

						actionsMap["cache_prefresh_parameters"] = []interface{}{cachePrefreshParametersMap}
					}

					accessURLRedirectParametersMap := map[string]interface{}{}

					if actions.AccessURLRedirectParameters != nil {
						if actions.AccessURLRedirectParameters.StatusCode != nil {
							accessURLRedirectParametersMap["status_code"] = actions.AccessURLRedirectParameters.StatusCode
						}

						if actions.AccessURLRedirectParameters.Protocol != nil {
							accessURLRedirectParametersMap["protocol"] = actions.AccessURLRedirectParameters.Protocol
						}

						hostNameMap := map[string]interface{}{}

						if actions.AccessURLRedirectParameters.HostName != nil {
							if actions.AccessURLRedirectParameters.HostName.Action != nil {
								hostNameMap["action"] = actions.AccessURLRedirectParameters.HostName.Action
							}

							if actions.AccessURLRedirectParameters.HostName.Value != nil {
								hostNameMap["value"] = actions.AccessURLRedirectParameters.HostName.Value
							}

							accessURLRedirectParametersMap["host_name"] = []interface{}{hostNameMap}
						}

						uRLPathMap := map[string]interface{}{}

						if actions.AccessURLRedirectParameters.URLPath != nil {
							if actions.AccessURLRedirectParameters.URLPath.Action != nil {
								uRLPathMap["action"] = actions.AccessURLRedirectParameters.URLPath.Action
							}

							if actions.AccessURLRedirectParameters.URLPath.Regex != nil {
								uRLPathMap["regex"] = actions.AccessURLRedirectParameters.URLPath.Regex
							}

							if actions.AccessURLRedirectParameters.URLPath.Value != nil {
								uRLPathMap["value"] = actions.AccessURLRedirectParameters.URLPath.Value
							}

							accessURLRedirectParametersMap["url_path"] = []interface{}{uRLPathMap}
						}

						queryStringMap := map[string]interface{}{}

						if actions.AccessURLRedirectParameters.QueryString != nil {
							if actions.AccessURLRedirectParameters.QueryString.Action != nil {
								queryStringMap["action"] = actions.AccessURLRedirectParameters.QueryString.Action
							}

							accessURLRedirectParametersMap["query_string"] = []interface{}{queryStringMap}
						}

						actionsMap["access_url_redirect_parameters"] = []interface{}{accessURLRedirectParametersMap}
					}

					upstreamURLRewriteParametersMap := map[string]interface{}{}

					if actions.UpstreamURLRewriteParameters != nil {
						if actions.UpstreamURLRewriteParameters.Type != nil {
							upstreamURLRewriteParametersMap["type"] = actions.UpstreamURLRewriteParameters.Type
						}

						if actions.UpstreamURLRewriteParameters.Action != nil {
							upstreamURLRewriteParametersMap["action"] = actions.UpstreamURLRewriteParameters.Action
						}

						if actions.UpstreamURLRewriteParameters.Value != nil {
							upstreamURLRewriteParametersMap["value"] = actions.UpstreamURLRewriteParameters.Value
						}

						if actions.UpstreamURLRewriteParameters.Regex != nil {
							upstreamURLRewriteParametersMap["regex"] = actions.UpstreamURLRewriteParameters.Regex

						}

						actionsMap["upstream_url_rewrite_parameters"] = []interface{}{upstreamURLRewriteParametersMap}
					}

					qUICParametersMap := map[string]interface{}{}

					if actions.QUICParameters != nil {
						if actions.QUICParameters.Switch != nil {
							qUICParametersMap["switch"] = actions.QUICParameters.Switch
						}

						actionsMap["quic_parameters"] = []interface{}{qUICParametersMap}
					}

					webSocketParametersMap := map[string]interface{}{}

					if actions.WebSocketParameters != nil {
						if actions.WebSocketParameters.Switch != nil {
							webSocketParametersMap["switch"] = actions.WebSocketParameters.Switch
						}

						if actions.WebSocketParameters.Timeout != nil {
							webSocketParametersMap["timeout"] = actions.WebSocketParameters.Timeout
						}

						actionsMap["web_socket_parameters"] = []interface{}{webSocketParametersMap}
					}

					authenticationParametersMap := map[string]interface{}{}

					if actions.AuthenticationParameters != nil {
						if actions.AuthenticationParameters.AuthType != nil {
							authenticationParametersMap["auth_type"] = actions.AuthenticationParameters.AuthType
						}

						if actions.AuthenticationParameters.SecretKey != nil {
							authenticationParametersMap["secret_key"] = actions.AuthenticationParameters.SecretKey
						}

						if actions.AuthenticationParameters.Timeout != nil {
							authenticationParametersMap["timeout"] = actions.AuthenticationParameters.Timeout
						}

						if actions.AuthenticationParameters.BackupSecretKey != nil {
							authenticationParametersMap["backup_secret_key"] = actions.AuthenticationParameters.BackupSecretKey
						}

						if actions.AuthenticationParameters.AuthParam != nil {
							authenticationParametersMap["auth_param"] = actions.AuthenticationParameters.AuthParam
						}

						if actions.AuthenticationParameters.TimeParam != nil {
							authenticationParametersMap["time_param"] = actions.AuthenticationParameters.TimeParam
						}

						if actions.AuthenticationParameters.TimeFormat != nil {
							authenticationParametersMap["time_format"] = actions.AuthenticationParameters.TimeFormat
						}

						actionsMap["authentication_parameters"] = []interface{}{authenticationParametersMap}
					}

					maxAgeParametersMap := map[string]interface{}{}

					if actions.MaxAgeParameters != nil {
						if actions.MaxAgeParameters.FollowOrigin != nil {
							maxAgeParametersMap["follow_origin"] = actions.MaxAgeParameters.FollowOrigin
						}

						if actions.MaxAgeParameters.CacheTime != nil {
							maxAgeParametersMap["cache_time"] = actions.MaxAgeParameters.CacheTime
						}

						actionsMap["max_age_parameters"] = []interface{}{maxAgeParametersMap}
					}

					statusCodeCacheParametersMap := map[string]interface{}{}

					if actions.StatusCodeCacheParameters != nil {
						statusCodeCacheParamsList := make([]map[string]interface{}, 0, len(actions.StatusCodeCacheParameters.StatusCodeCacheParams))
						if actions.StatusCodeCacheParameters.StatusCodeCacheParams != nil {
							for _, statusCodeCacheParams := range actions.StatusCodeCacheParameters.StatusCodeCacheParams {
								statusCodeCacheParamsMap := map[string]interface{}{}

								if statusCodeCacheParams.StatusCode != nil {
									statusCodeCacheParamsMap["status_code"] = statusCodeCacheParams.StatusCode
								}

								if statusCodeCacheParams.CacheTime != nil {
									statusCodeCacheParamsMap["cache_time"] = statusCodeCacheParams.CacheTime
								}

								statusCodeCacheParamsList = append(statusCodeCacheParamsList, statusCodeCacheParamsMap)
							}

							statusCodeCacheParametersMap["status_code_cache_params"] = statusCodeCacheParamsList
						}
						actionsMap["status_code_cache_parameters"] = []interface{}{statusCodeCacheParametersMap}
					}

					offlineCacheParametersMap := map[string]interface{}{}

					if actions.OfflineCacheParameters != nil {
						if actions.OfflineCacheParameters.Switch != nil {
							offlineCacheParametersMap["switch"] = actions.OfflineCacheParameters.Switch
						}

						actionsMap["offline_cache_parameters"] = []interface{}{offlineCacheParametersMap}
					}

					smartRoutingParametersMap := map[string]interface{}{}

					if actions.SmartRoutingParameters != nil {
						if actions.SmartRoutingParameters.Switch != nil {
							smartRoutingParametersMap["switch"] = actions.SmartRoutingParameters.Switch
						}

						actionsMap["smart_routing_parameters"] = []interface{}{smartRoutingParametersMap}
					}

					rangeOriginPullParametersMap := map[string]interface{}{}

					if actions.RangeOriginPullParameters != nil {
						if actions.RangeOriginPullParameters.Switch != nil {
							rangeOriginPullParametersMap["switch"] = actions.RangeOriginPullParameters.Switch
						}

						actionsMap["range_origin_pull_parameters"] = []interface{}{rangeOriginPullParametersMap}
					}

					upstreamHTTP2ParametersMap := map[string]interface{}{}

					if actions.UpstreamHTTP2Parameters != nil {
						if actions.UpstreamHTTP2Parameters.Switch != nil {
							upstreamHTTP2ParametersMap["switch"] = actions.UpstreamHTTP2Parameters.Switch
						}

						actionsMap["upstream_http2_parameters"] = []interface{}{upstreamHTTP2ParametersMap}
					}

					hostHeaderParametersMap := map[string]interface{}{}

					if actions.HostHeaderParameters != nil {
						if actions.HostHeaderParameters.Action != nil {
							hostHeaderParametersMap["action"] = actions.HostHeaderParameters.Action
						}

						if actions.HostHeaderParameters.ServerName != nil {
							hostHeaderParametersMap["server_name"] = actions.HostHeaderParameters.ServerName
						}

						actionsMap["host_header_parameters"] = []interface{}{hostHeaderParametersMap}
					}

					forceRedirectHTTPSParametersMap := map[string]interface{}{}

					if actions.ForceRedirectHTTPSParameters != nil {
						if actions.ForceRedirectHTTPSParameters.Switch != nil {
							forceRedirectHTTPSParametersMap["switch"] = actions.ForceRedirectHTTPSParameters.Switch
						}

						if actions.ForceRedirectHTTPSParameters.RedirectStatusCode != nil {
							forceRedirectHTTPSParametersMap["redirect_status_code"] = actions.ForceRedirectHTTPSParameters.RedirectStatusCode
						}

						actionsMap["force_redirect_https_parameters"] = []interface{}{forceRedirectHTTPSParametersMap}
					}

					originPullProtocolParametersMap := map[string]interface{}{}

					if actions.OriginPullProtocolParameters != nil {
						if actions.OriginPullProtocolParameters.Protocol != nil {
							originPullProtocolParametersMap["protocol"] = actions.OriginPullProtocolParameters.Protocol
						}

						actionsMap["origin_pull_protocol_parameters"] = []interface{}{originPullProtocolParametersMap}
					}

					compressionParametersMap := map[string]interface{}{}

					if actions.CompressionParameters != nil {
						if actions.CompressionParameters.Switch != nil {
							compressionParametersMap["switch"] = actions.CompressionParameters.Switch
						}

						if actions.CompressionParameters.Algorithms != nil {
							compressionParametersMap["algorithms"] = actions.CompressionParameters.Algorithms
						}

						actionsMap["compression_parameters"] = []interface{}{compressionParametersMap}
					}

					hSTSParametersMap := map[string]interface{}{}

					if actions.HSTSParameters != nil {
						if actions.HSTSParameters.Switch != nil {
							hSTSParametersMap["switch"] = actions.HSTSParameters.Switch
						}

						if actions.HSTSParameters.Timeout != nil {
							hSTSParametersMap["timeout"] = actions.HSTSParameters.Timeout
						}

						if actions.HSTSParameters.IncludeSubDomains != nil {
							hSTSParametersMap["include_sub_domains"] = actions.HSTSParameters.IncludeSubDomains
						}

						if actions.HSTSParameters.Preload != nil {
							hSTSParametersMap["preload"] = actions.HSTSParameters.Preload
						}

						actionsMap["hsts_parameters"] = []interface{}{hSTSParametersMap}
					}

					clientIPHeaderParametersMap := map[string]interface{}{}

					if actions.ClientIPHeaderParameters != nil {
						if actions.ClientIPHeaderParameters.Switch != nil {
							clientIPHeaderParametersMap["switch"] = actions.ClientIPHeaderParameters.Switch
						}

						if actions.ClientIPHeaderParameters.HeaderName != nil {
							clientIPHeaderParametersMap["header_name"] = actions.ClientIPHeaderParameters.HeaderName
						}

						actionsMap["client_ip_header_parameters"] = []interface{}{clientIPHeaderParametersMap}
					}

					oCSPStaplingParametersMap := map[string]interface{}{}

					if actions.OCSPStaplingParameters != nil {
						if actions.OCSPStaplingParameters.Switch != nil {
							oCSPStaplingParametersMap["switch"] = actions.OCSPStaplingParameters.Switch
						}

						actionsMap["ocsp_stapling_parameters"] = []interface{}{oCSPStaplingParametersMap}
					}

					hTTP2ParametersMap := map[string]interface{}{}

					if actions.HTTP2Parameters != nil {
						if actions.HTTP2Parameters.Switch != nil {
							hTTP2ParametersMap["switch"] = actions.HTTP2Parameters.Switch
						}

						actionsMap["http2_parameters"] = []interface{}{hTTP2ParametersMap}
					}

					postMaxSizeParametersMap := map[string]interface{}{}

					if actions.PostMaxSizeParameters != nil {
						if actions.PostMaxSizeParameters.Switch != nil {
							postMaxSizeParametersMap["switch"] = actions.PostMaxSizeParameters.Switch
						}

						if actions.PostMaxSizeParameters.MaxSize != nil {
							postMaxSizeParametersMap["max_size"] = actions.PostMaxSizeParameters.MaxSize
						}

						actionsMap["post_max_size_parameters"] = []interface{}{postMaxSizeParametersMap}
					}

					clientIPCountryParametersMap := map[string]interface{}{}

					if actions.ClientIPCountryParameters != nil {
						if actions.ClientIPCountryParameters.Switch != nil {
							clientIPCountryParametersMap["switch"] = actions.ClientIPCountryParameters.Switch
						}

						if actions.ClientIPCountryParameters.HeaderName != nil {
							clientIPCountryParametersMap["header_name"] = actions.ClientIPCountryParameters.HeaderName
						}

						actionsMap["client_ip_country_parameters"] = []interface{}{clientIPCountryParametersMap}
					}

					upstreamFollowRedirectParametersMap := map[string]interface{}{}

					if actions.UpstreamFollowRedirectParameters != nil {
						if actions.UpstreamFollowRedirectParameters.Switch != nil {
							upstreamFollowRedirectParametersMap["switch"] = actions.UpstreamFollowRedirectParameters.Switch
						}

						if actions.UpstreamFollowRedirectParameters.MaxTimes != nil {
							upstreamFollowRedirectParametersMap["max_times"] = actions.UpstreamFollowRedirectParameters.MaxTimes
						}

						actionsMap["upstream_follow_redirect_parameters"] = []interface{}{upstreamFollowRedirectParametersMap}
					}

					upstreamRequestParametersMap := map[string]interface{}{}

					if actions.UpstreamRequestParameters != nil {
						queryStringMap := map[string]interface{}{}

						if actions.UpstreamRequestParameters.QueryString != nil {
							if actions.UpstreamRequestParameters.QueryString.Switch != nil {
								queryStringMap["switch"] = actions.UpstreamRequestParameters.QueryString.Switch
							}

							if actions.UpstreamRequestParameters.QueryString.Action != nil {
								queryStringMap["action"] = actions.UpstreamRequestParameters.QueryString.Action
							}

							if actions.UpstreamRequestParameters.QueryString.Values != nil {
								queryStringMap["values"] = actions.UpstreamRequestParameters.QueryString.Values
							}

							upstreamRequestParametersMap["query_string"] = []interface{}{queryStringMap}
						}

						cookieMap := map[string]interface{}{}

						if actions.UpstreamRequestParameters.Cookie != nil {
							if actions.UpstreamRequestParameters.Cookie.Switch != nil {
								cookieMap["switch"] = actions.UpstreamRequestParameters.Cookie.Switch
							}

							if actions.UpstreamRequestParameters.Cookie.Action != nil {
								cookieMap["action"] = actions.UpstreamRequestParameters.Cookie.Action
							}

							if actions.UpstreamRequestParameters.Cookie.Values != nil {
								cookieMap["values"] = actions.UpstreamRequestParameters.Cookie.Values
							}

							upstreamRequestParametersMap["cookie"] = []interface{}{cookieMap}
						}

						actionsMap["upstream_request_parameters"] = []interface{}{upstreamRequestParametersMap}
					}

					tLSConfigParametersMap := map[string]interface{}{}

					if actions.TLSConfigParameters != nil {
						if actions.TLSConfigParameters.Version != nil {
							tLSConfigParametersMap["version"] = actions.TLSConfigParameters.Version
						}

						if actions.TLSConfigParameters.CipherSuite != nil {
							tLSConfigParametersMap["cipher_suite"] = actions.TLSConfigParameters.CipherSuite
						}

						actionsMap["tls_config_parameters"] = []interface{}{tLSConfigParametersMap}
					}

					modifyOriginParametersMap := map[string]interface{}{}

					if actions.ModifyOriginParameters != nil {
						if actions.ModifyOriginParameters.OriginType != nil {
							modifyOriginParametersMap["origin_type"] = actions.ModifyOriginParameters.OriginType
						}

						if actions.ModifyOriginParameters.Origin != nil {
							modifyOriginParametersMap["origin"] = actions.ModifyOriginParameters.Origin
						}

						if actions.ModifyOriginParameters.OriginProtocol != nil {
							modifyOriginParametersMap["origin_protocol"] = actions.ModifyOriginParameters.OriginProtocol
						}

						if actions.ModifyOriginParameters.HTTPOriginPort != nil {
							modifyOriginParametersMap["http_origin_port"] = actions.ModifyOriginParameters.HTTPOriginPort
						}

						if actions.ModifyOriginParameters.HTTPSOriginPort != nil {
							modifyOriginParametersMap["https_origin_port"] = actions.ModifyOriginParameters.HTTPSOriginPort
						}

						if actions.ModifyOriginParameters.PrivateAccess != nil {
							modifyOriginParametersMap["private_access"] = actions.ModifyOriginParameters.PrivateAccess
						}

						privateParametersMap := map[string]interface{}{}

						if actions.ModifyOriginParameters.PrivateParameters != nil {
							if actions.ModifyOriginParameters.PrivateParameters.AccessKeyId != nil {
								privateParametersMap["access_key_id"] = actions.ModifyOriginParameters.PrivateParameters.AccessKeyId
							}

							if actions.ModifyOriginParameters.PrivateParameters.SecretAccessKey != nil {
								privateParametersMap["secret_access_key"] = actions.ModifyOriginParameters.PrivateParameters.SecretAccessKey
							}

							if actions.ModifyOriginParameters.PrivateParameters.SignatureVersion != nil {
								privateParametersMap["signature_version"] = actions.ModifyOriginParameters.PrivateParameters.SignatureVersion
							}

							if actions.ModifyOriginParameters.PrivateParameters.Region != nil {
								privateParametersMap["region"] = actions.ModifyOriginParameters.PrivateParameters.Region
							}

							modifyOriginParametersMap["private_parameters"] = []interface{}{privateParametersMap}
						}

						actionsMap["modify_origin_parameters"] = []interface{}{modifyOriginParametersMap}
					}

					hTTPUpstreamTimeoutParametersMap := map[string]interface{}{}

					if actions.HTTPUpstreamTimeoutParameters != nil {
						if actions.HTTPUpstreamTimeoutParameters.ResponseTimeout != nil {
							hTTPUpstreamTimeoutParametersMap["response_timeout"] = actions.HTTPUpstreamTimeoutParameters.ResponseTimeout
						}

						actionsMap["http_upstream_timeout_parameters"] = []interface{}{hTTPUpstreamTimeoutParametersMap}
					}

					httpResponseParametersMap := map[string]interface{}{}

					if actions.HttpResponseParameters != nil {
						if actions.HttpResponseParameters.StatusCode != nil {
							httpResponseParametersMap["status_code"] = actions.HttpResponseParameters.StatusCode
						}

						if actions.HttpResponseParameters.ResponsePage != nil {
							httpResponseParametersMap["response_page"] = actions.HttpResponseParameters.ResponsePage
						}

						actionsMap["http_response_parameters"] = []interface{}{httpResponseParametersMap}
					}

					errorPageParametersMap := map[string]interface{}{}

					if actions.ErrorPageParameters != nil {
						errorPageParamsList := make([]map[string]interface{}, 0, len(actions.ErrorPageParameters.ErrorPageParams))
						if actions.ErrorPageParameters.ErrorPageParams != nil {
							for _, errorPageParams := range actions.ErrorPageParameters.ErrorPageParams {
								errorPageParamsMap := map[string]interface{}{}

								if errorPageParams.StatusCode != nil {
									errorPageParamsMap["status_code"] = errorPageParams.StatusCode
								}

								if errorPageParams.RedirectURL != nil {
									errorPageParamsMap["redirect_url"] = errorPageParams.RedirectURL
								}

								errorPageParamsList = append(errorPageParamsList, errorPageParamsMap)
							}

							errorPageParametersMap["error_page_params"] = errorPageParamsList
						}
						actionsMap["error_page_parameters"] = []interface{}{errorPageParametersMap}
					}

					modifyResponseHeaderParametersMap := map[string]interface{}{}

					if actions.ModifyResponseHeaderParameters != nil {
						headerActionsList := make([]map[string]interface{}, 0, len(actions.ModifyResponseHeaderParameters.HeaderActions))
						if actions.ModifyResponseHeaderParameters.HeaderActions != nil {
							for _, headerActions := range actions.ModifyResponseHeaderParameters.HeaderActions {
								headerActionsMap := map[string]interface{}{}

								if headerActions.Action != nil {
									headerActionsMap["action"] = headerActions.Action
								}

								if headerActions.Name != nil {
									headerActionsMap["name"] = headerActions.Name
								}

								if headerActions.Value != nil {
									headerActionsMap["value"] = headerActions.Value
								}

								headerActionsList = append(headerActionsList, headerActionsMap)
							}

							modifyResponseHeaderParametersMap["header_actions"] = headerActionsList
						}
						actionsMap["modify_response_header_parameters"] = []interface{}{modifyResponseHeaderParametersMap}
					}

					modifyRequestHeaderParametersMap := map[string]interface{}{}

					if actions.ModifyRequestHeaderParameters != nil {
						headerActionsList := make([]map[string]interface{}, 0, len(actions.ModifyRequestHeaderParameters.HeaderActions))
						if actions.ModifyRequestHeaderParameters.HeaderActions != nil {
							for _, headerActions := range actions.ModifyRequestHeaderParameters.HeaderActions {
								headerActionsMap := map[string]interface{}{}

								if headerActions.Action != nil {
									headerActionsMap["action"] = headerActions.Action
								}

								if headerActions.Name != nil {
									headerActionsMap["name"] = headerActions.Name
								}

								if headerActions.Value != nil {
									headerActionsMap["value"] = headerActions.Value
								}

								headerActionsList = append(headerActionsList, headerActionsMap)
							}

							modifyRequestHeaderParametersMap["header_actions"] = headerActionsList
						}
						actionsMap["modify_request_header_parameters"] = []interface{}{modifyRequestHeaderParametersMap}
					}

					responseSpeedLimitParametersMap := map[string]interface{}{}

					if actions.ResponseSpeedLimitParameters != nil {
						if actions.ResponseSpeedLimitParameters.Mode != nil {
							responseSpeedLimitParametersMap["mode"] = actions.ResponseSpeedLimitParameters.Mode
						}

						if actions.ResponseSpeedLimitParameters.MaxSpeed != nil {
							responseSpeedLimitParametersMap["max_speed"] = actions.ResponseSpeedLimitParameters.MaxSpeed
						}

						if actions.ResponseSpeedLimitParameters.StartAt != nil {
							responseSpeedLimitParametersMap["start_at"] = actions.ResponseSpeedLimitParameters.StartAt
						}

						actionsMap["response_speed_limit_parameters"] = []interface{}{responseSpeedLimitParametersMap}
					}

					setContentIdentifierParametersMap := map[string]interface{}{}

					if actions.SetContentIdentifierParameters != nil {
						if actions.SetContentIdentifierParameters.ContentIdentifier != nil {
							setContentIdentifierParametersMap["content_identifier"] = actions.SetContentIdentifierParameters.ContentIdentifier
						}

						actionsMap["set_content_identifier_parameters"] = []interface{}{setContentIdentifierParametersMap}
					}

					contentCompressionParametersMap := map[string]interface{}{}
					if actions.ContentCompressionParameters != nil {
						if actions.ContentCompressionParameters.Switch != nil {
							contentCompressionParametersMap["switch"] = actions.ContentCompressionParameters.Switch
						}

						actionsMap["content_compression_parameters"] = []interface{}{contentCompressionParametersMap}
					}

					actionsList = append(actionsList, actionsMap)
				}

				branchesMap["actions"] = actionsList
			}

			subRulesList := make([]map[string]interface{}, 0, len(branches.SubRules))
			if branches.SubRules != nil {
				for _, subRules := range branches.SubRules {
					subRulesMap := map[string]interface{}{}

					if subRules.Branches != nil {
						subRulesMap["branches"] = resourceTencentCloudTeoL7AccRuleSetBranchs(subRules.Branches)
					}

					if subRules.Description != nil {
						subRulesMap["description"] = subRules.Description
					}

					subRulesList = append(subRulesList, subRulesMap)
				}

				branchesMap["sub_rules"] = subRulesList
			}

			branchesList = append(branchesList, branchesMap)
		}
	}
	return branchesList
}

func resourceTencentCloudTeoL7AccRuleContent(rules []*teo.RuleEngineItem) (string, error) {
	type Content struct {
		FormatVersion string                `json:"FormatVersion,omitempty"`
		Rules         []*teo.RuleEngineItem `json:"Rules,omitempty"`
	}
	content := Content{
		FormatVersion: "1.0",
		Rules:         rules,
	}
	contentBytes, err := json.Marshal(content)
	if err != nil {
		return "", err
	}
	return string(contentBytes), nil
}
