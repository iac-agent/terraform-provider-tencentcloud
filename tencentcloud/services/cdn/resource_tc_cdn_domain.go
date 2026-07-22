package cdn

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudCdnDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdnDomainCreate,
		Read:   resourceTencentCloudCdnDomainRead,
		Update: resourceTencentCloudCdnDomainUpdate,
		Delete: resourceTencentCloudCdnDomainDelete,
		Importer: &schema.ResourceImporter{
			//State: func(d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
			//	getDefaultSwitchOffMap := func() []interface{} {
			//		return []interface{}{
			//			map[string]interface{}{"switch": "off"},
			//		}
			//	}
			//	_ = d.Set("authentication", getDefaultSwitchOffMap())
			//	return []*schema.ResourceData{d}, nil
			//},
			State: helper.ImportWithDefaultValue(map[string]interface{}{
				"authentication": []interface{}{map[string]interface{}{
					"switch": "off",
				}},
				"cache_key": []interface{}{map[string]interface{}{
					"full_url_cache": "on",
				}},
				"full_url_cache": true,
			}),
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "名称 acceleration 域名",
			},
			"service_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SERVICE_TYPE),
				Description:  "Acceleration 域名 名称 服务 类型 `web`: 静态 acceleration，`download`: download acceleration，`media`: streaming media VOD acceleration，`hybrid`: hybrid acceleration，`动态`: 动态 acceleration。",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "项目 CDN belongs 到，默认为 0。",
			},
			"area": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_AREA),
				Description:  "域名 名称 acceleration 地域 `mainland`: acceleration inside mainland China，`overseas`: acceleration outside mainland China，`全局`: 全局 acceleration. Overseas acceleration 服务 必须 是 已启用 到 使用 overseas acceleration 和 全局 acceleration。",
			},
			"full_url_cache": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       true,
				ConflictsWith: []string{"cache_key"},
				Deprecated:    "Use `cache_key` -> `full_url_cache` instead.",
				Description:   "是否enable full-路径 缓存. 默认值为 `true`。",
			},
			"origin": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Origin 服务器 配置. It's 列表 和 consist 的 在 most 一个 item。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_ORIGIN_TYPE),
							Description: "Master 源站 服务器 类型. following types 是 支持: `域名`: Domain 名称, `domainv6`: IPv6 域名 名称, `cos`: COS 存储桶 地址, `third_party`: Third-party 对象 存储 源站," +
								"`igtm`: IGTM origin, `ip`: IP address, `ipv6`: One IPv6 address, `ip_ipv6`: Multiple IPv4 addresses and one IPv6 address, " +
								"`ip_domain`: IP addresses and domain names (only available to beta users), `ip_domainv6`: Multiple IPv4 addresses and one IPv6 domain name, `ipv6_domain`: Multiple IPv6 addresses and one domain name, `ipv6_domainv6`: Multiple IPv6 addresses and one IPv6 domain name, " +
								"`domain_domainv6`: Multiple IPv4 domain names and one IPv6 domain name, `ip_ipv6_domain`: Multiple IPv4 and IPv6 addresses and one domain name, `ip_ipv6_domainv6`: Multiple IPv4 and IPv6 addresses and one IPv6 domain name, `ip_domain_domainv6`: Multiple IPv4 addresses and IPv4 domain names and one IPv6 domain name, " +
								"`ipv6_domain_domainv6`: Multiple IPv4 domain names and IPv6 addresses and one IPv6 domain name, `ip_ipv6_domain_domainv6`: Multiple IPv4 and IPv6 addresses and IPv4 domain names and one IPv6 domain name.",
						},
						"origin_list": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Master 源站 服务器 列表. 有效 值 可以 是 ip 或 域名 名称 当 modifying 源站 服务器，您 need 到 enter corresponding `origin_type`。",
						},
						"server_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "主机 头部 使用 当 accessing master 源站 服务器. 如果为空， acceleration 域名 名称 将 是 使用 通过 默认值。",
						},
						"cos_private_access": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     CDN_SWITCH_OFF,
							Description: "当 OriginType 是 COS，您 可以 指定if 访问 到 私有 buckets 是 allowed. 有效 值 是 `在` 和 `关闭`. 和 默认值为 `关闭`。",
						},
						"origin_pull_protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_ORIGIN_PULL_PROTOCOL_HTTP,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_ORIGIN_PULL_PROTOCOL),
							Description:  "Origin-pull 协议 配置. `http`: forced HTTP 源站-pull，`follow`: 协议 follow 源站-pull，`https`: forced HTTPS 源站-pull. 此 仅 支持 源站 服务器 端口 443 对于 源站-pull。",
						},
						"backup_origin_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_BACKUP_ORIGIN_TYPE),
							Description:  "Backup 源站 服务器 类型，其中 支持 following types: `域名`: 域名 名称 类型，`ip`: IP 列表 使用 作为 源站 服务器，`ipv6_domain`: Multiple IPv6 addresses 和 一个 域名 名称，`ip_ipv6`: Multiple IPv4 addresses 和 一个 IPv6 地址，`ip_ipv6_domain`: Multiple IPv4 和 IPv6 addresses 和 一个 域名 名称",
						},
						"backup_origin_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Backup 源站 服务器 列表. 有效 值 可以 是 ip 或 域名 名称 当 modifying 备份 源站 服务器，您 need 到 enter corresponding `backup_origin_type`。",
						},
						"backup_server_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "主机 头部 使用 当 accessing 备份 源站 服务器. 如果为空， ServerName 的 master 源站 服务器 将 是 使用 通过 默认值。",
						},
						"origin_company": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Object 存储 back 到 来源 vendor. 必填 当 来源 station 类型 是 third-party 存储 来源 station (third_party). 可选 值 include following: `aws_s3`: AWS S3; `ali_oss`: Alibaba Cloud OSS; `hw_obs`: Huawei OBS; `qiniu_kodo`: Qiniu Cloud kodo; `others`: other vendors' 对象 存储，仅 支持 对象 存储 compatible 使用 AWS 签名 algorithm，such 作为 Tencent Cloud Financial 可用区 COS. Example 值: `hw_obs`。",
						},
					},
				},
			},
			"https_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "HTTPS acceleration 配置. It's 列表 和 consist 的 在 most 一个 item。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"https_switch": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "HTTPS 配置 switch. 有效 值 是 `在` 和 `关闭`。",
						},
						"http2_switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "HTTP2 配置 switch. 有效 值 是 `在` 和 `关闭`. 和 默认值为 `关闭`。",
						},
						"ocsp_stapling_switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "OCSP 配置 switch. 有效 值 是 `在` 和 `关闭`. 和 默认值为 `关闭`。",
						},
						"spdy_switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "Spdy 配置 switch. 有效 值 是 `在` 和 `关闭`. 和 默认值为 `关闭`. 此 参数 是 对于 white-列表 customer。",
						},
						"verify_client": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "Client 证书 authentication 功能. 有效 值 是 `在` 和 `关闭`. 和 默认值为 `关闭`。",
						},
						"server_certificate_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Server 证书 配置 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Server 证书 ID",
									},
									"certificate_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Server 证书 名称",
									},
									"certificate_content": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Server 证书 信息. 此 为必填项 当 uploading 外部 证书，其中 should contain 完整 证书 chain。",
									},
									"private_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Server 键 信息. 此 为必填项 当 uploading 外部 证书。",
									},
									"message": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Certificate 备注",
									},
									"deploy_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deploy 时间 的 服务器 证书。",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Expire 时间 的 服务器 证书。",
									},
								},
							},
						},
						"client_certificate_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Client 证书 配置 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Client 证书 名称",
									},
									"certificate_content": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Client Certificate PEM 格式，requires Base64 编码。",
									},
									"deploy_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Deploy 时间 的 客户端 证书。",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Expire 时间 的 客户端 证书。",
									},
								},
							},
						},
						"force_redirect": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "Configuration 的 forced HTTP 或 HTTPS redirects。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      CDN_SWITCH_OFF,
										ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
										Description:  "Forced redirect 配置 switch. 有效 值 是 `在` 和 `关闭`. 默认值为 `关闭`。",
									},
									"redirect_type": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      CDN_ORIGIN_PULL_PROTOCOL_HTTP,
										ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_FORCE_REDIRECT_TYPE),
										Description: "Forced redirect 类型. 有效 值 是 `http` 和 `https`. `http` 表示 forced redirect 从 HTTPS 到 HTTP, `https` 表示 forced redirect 从 HTTP 到 HTTPS." +
											"When `switch` setting `off`, this property does not need to be set or set to `http`. Default value is `http`.",
									},
									"redirect_status_code": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      302,
										ValidateFunc: tccommon.ValidateAllowedIntValue([]int{301, 302}),
										Description: "Forced redirect 状态 代码. 有效 值 是 `301` 和 `302`." +
											"When `switch` setting `off`, this property does not need to be set or set to `302`. Default value is `302`.",
									},
									"carry_headers": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "off",
										Description: "是否return newly added 头部 during force redirection. Values: `在`，`关闭`。",
									},
								},
							},
						},
						"tls_versions": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Tls 版本 settings，仅 support some Advanced 域名 names，support settings TLSv1，TLSV1.1，TLSV1.2，TLSv1.3，当 modifying 必须 open consecutive versions。",
						},
						"hsts": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "HSTS 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
										Description:  "HSTS 配置 switch. 有效 值 是 `在` 和 `关闭`。",
									},
									"max_age": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "MaxAge 值",
									},
									"include_sub_domains": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
										Description:  "是否include sub domains，值 `在` 和 `关闭`。",
									},
								},
							},
						},
					},
				},
			},
			"range_origin_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CDN_SWITCH_ON,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
				Description:  "Sharding back 到 来源 配置 switch. 有效 值 是 `在` 和 `关闭`. 默认值为 `在`。",
			},
			"ipv6_access_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CDN_SWITCH_OFF,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
				Description:  "ipv6 访问 配置 switch. Only 可用 当 area 集合 到 `mainland`. 有效 值 是 `在` 和 `关闭`. 默认值为 `关闭`。",
			},
			"follow_redirect_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CDN_SWITCH_OFF,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
				Description:  "301/302 redirect following switch，可用值：`在`，`关闭` (默认值)。",
			},
			"authentication": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "指定timestamp hotlink protection 配置，NOTE: 仅 一个 类型 可以 choose 对于 sub elements。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Authentication switching，可用值：`在`，`关闭`。",
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
						},
						"type_a": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "时间戳 hotlink protection 模式 A 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"secret_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "键 对于 签名 calculation. Only digits，upper 和 lower-case letters 是 allowed. Length 限制: 6-32 字符。",
									},
									"sign_param": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "签名 参数 名称 Only upper 和 lower-case letters，digits，和 underscores (_) 是 allowed. It 不能 start 使用 digit. Length 限制: 1-100 字符。",
									},
									"expire_time": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "签名 过期时间 在 second. 最大 值 是 630720000。",
									},
									"file_extensions": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "File extension 列表 settings determining 如果 authentication should 是 performed. NOTE: 如果 它 包含an asterisk (*)，此 表示all files。",
									},
									"filter_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "可用值：`whitelist` - all types apart 从 `file_extensions` 是 authenticated，`blacklist`: - 仅 types 在 `file_extensions` 是 authenticated。",
									},
									"backup_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "用于calculate 签名 6-32 字符. Only digits 和 letters 是 allowed。",
									},
								},
							},
						},
						"type_b": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "时间戳 hotlink protection 模式 B 配置. NOTE: according 到 upgrading 的 TencentCloud Platform，TypeB 是 unavailable 对于 now。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"secret_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "键 对于 签名 calculation. Only digits，upper 和 lower-case letters 是 allowed. Length 限制: 6-32 字符。",
									},
									"expire_time": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "签名 过期时间 在 second. 最大 值 是 630720000。",
									},
									"file_extensions": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "File extension 列表 settings determining 如果 authentication should 是 performed. NOTE: 如果 它 包含an asterisk (*)，此 表示all files。",
									},
									"filter_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "可用值：`whitelist` - all types apart 从 `file_extensions` 是 authenticated，`blacklist`: - 仅 types 在 `file_extensions` 是 authenticated。",
									},
									"backup_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "用于calculate 签名 6-32 字符. Only digits 和 letters 是 allowed。",
									},
								},
							},
						},
						"type_c": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "时间戳 hotlink protection 模式 C 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"secret_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "键 对于 签名 calculation. Only digits，upper 和 lower-case letters 是 allowed. Length 限制: 6-32 字符。",
									},
									"expire_time": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "签名 过期时间 在 second. 最大 值 是 630720000。",
									},
									"file_extensions": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "File extension 列表 settings determining 如果 authentication should 是 performed. NOTE: 如果 它 包含an asterisk (*)，此 表示all files。",
									},
									"filter_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "可用值：`whitelist` - all types apart 从 `file_extensions` 是 authenticated，`blacklist`: - 仅 types 在 `file_extensions` 是 authenticated。",
									},
									"time_format": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "时间戳 formation，可用值：`dec`，`hex`。",
									},
									"backup_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "用于calculate 签名 6-32 字符. Only digits 和 letters 是 allowed。",
									},
								},
							},
						},
						"type_d": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "时间戳 hotlink protection 模式 D 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"secret_key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "键 对于 签名 calculation. Only digits，upper 和 lower-case letters 是 allowed. Length 限制: 6-32 字符。",
									},
									"expire_time": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "签名 过期时间 在 second. 最大 值 是 630720000。",
									},
									"file_extensions": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "File extension 列表 settings determining 如果 authentication should 是 performed. NOTE: 如果 它 包含an asterisk (*)，此 表示all files。",
									},
									"filter_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "可用值：`whitelist` - all types apart 从 `file_extensions` 是 authenticated，`blacklist`: - 仅 types 在 `file_extensions` 是 authenticated。",
									},
									"time_param": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "时间戳 参数 名称 Only upper 和 lower-case letters，digits，和 underscores (_) 是 allowed. It 不能 start 使用 digit. Length 限制: 1-100 字符。",
									},
									"time_format": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "时间戳 formation，可用值：`dec`，`hex`。",
									},
									"backup_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "用于calculate 签名 6-32 字符. Only digits 和 letters 是 allowed。",
									},
								},
							},
						},
					},
				},
			},
			"rule_cache": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Advanced 路径 缓存 配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_paths": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Matching 内容 under corresponding 类型 的 CacheType: `all`: fill *, `文件`: fill 在 suffix 名称, such 作为 jpg, txt," +
								"`directory`: fill in the path, such as /xxx/test, `path`: fill in the absolute path, such as /xxx/test.html, `index`: fill /.",
						},
						"rule_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_RULE_TYPE_DEFAULT,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_RULE_TYPE),
							Description: "Rule 类型. following types 是 支持: `all`: all documents take effect, `文件`: 指定 文件 suffix takes effect," +
								"`directory`: the specified path takes effect, `path`: specify the absolute path to take effect, `index`: home page.",
						},
						"switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "Cache 配置 switch. 有效 值 是 `在` 和 `关闭`。",
						},
						"cache_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Cache 过期时间 setting， 单位 是 second， 最大 可以 是 集合 到 365 days。",
						},
						"compare_max_age": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description: "Advanced 缓存 expiration 配置. 当 它 是 turned 在, 它 将 compare max-age 值 返回 通过 源站 site 使用 缓存 expiration 时间 集合 在 CacheRules," +
								"and take the minimum value to cache at the node. Valid values are `on` and `off`. Default value is `off`.",
						},
						"ignore_cache_control": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description: "Force caching. After opening, 无-store 和 无-缓存 resources 返回 通过 源站 site 将 also 是 cached 在 accordance 使用 CacheRules" +
								"rules. Valid values are `on` and `off`. Default value is `off`.",
						},
						"ignore_set_cookie": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "Ignore Set-Cookie 头部 的 源站 site. 有效 值 是 `在` 和 `关闭`. 默认值为 `关闭`. 此 参数 是 对于 white-列表 customer。",
						},
						"no_cache_switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "Cache 配置 switch. 有效 值 是 `在` 和 `关闭`。",
						},
						"re_validate": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "Always check back 到 源站. 有效 值 是 `在` 和 `关闭`. 默认值为 `关闭`。",
						},
						"follow_origin_switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "Follow 来源 station 配置 switch. 有效 值 是 `在` 和 `关闭`。",
						},
						"heuristic_cache_switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "指定是否enable heuristic 缓存，仅 可用 while `follow_origin_switch` 已启用，值: `在`，`关闭` (Default)。",
						},
						"heuristic_cache_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "指定heuristic 缓存 时间 在 second，仅 可用 while `follow_origin_switch` 和 `heuristic_cache_switch` 已启用",
						},
					},
				},
			},
			"request_header": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Request 头部 配置. It's 列表 和 consist 的 在 most 一个 item。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      CDN_SWITCH_OFF,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "Custom 请求 头部 配置 switch. 有效 值 是 `在` 和 `关闭`. 和 默认值为 `关闭`。",
						},
						"header_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Custom 请求 头部 配置 规则。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"header_mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Http 头部 setting 方法. following types 是 支持: `集合`: sets 值 对于 existing 头部 参数， new 头部 参数，或 多个 头部 参数. Multiple 头部 参数 将 是 merged into 一个; `del`: deletes 头部 参数; `add`: adds 头部 参数. By 默认值，您 可以 repeat same 操作 到 add same 头部 参数，其中 可能 affect browser response. Please consider 集合 operation first。",
									},
									"header_name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: tccommon.ValidateStringLengthInRange(1, 100),
										Description:  "Http 头部 名称",
									},
									"header_value": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: tccommon.ValidateStringLengthInRange(1, 1000),
										Description:  "Http 头部 值，可选 当 模式 是 `del`，必填 当 模式 是 `add`/`集合`。",
									},
									"rule_type": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_HEADER_RULE),
										Description: "Rule 类型. following types 是 支持: `all`: all documents take effect, `文件`: 指定 文件 suffix takes effect," +
											"`directory`: the specified path takes effect, `path`: specify the absolute path to take effect.",
									},
									"rule_paths": {
										Type:     schema.TypeList,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Description: "Matching 内容 under corresponding 类型 的 CacheType: `all`: fill *, `文件`: fill 在 suffix 名称, such 作为 jpg, txt," +
											"`directory`: fill in the path, such as /xxx/test, `path`: fill in the absolute path, such as /xxx/test.html.",
									},
								},
							},
						},
					},
				},
			},
			// extensions
			"ip_filter": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "指定Ip 过滤器 configurations。",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch，可用值：`在`，`关闭` (默认值)。",
						},
						"filter_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IP `blacklist`/`whitelist` 类型",
						},
						"filters": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Ip 过滤器 列表，Supports IPs 在 X.X.X.X 格式，或 /8，/16，/24 格式 IP ranges. Up 到 50 allowlists 或 blocklists 可以 是 entered。",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"filter_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Ip 过滤器 规则，此 功能 是 仅 可用 到 selected beta customers。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filter_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Ip 过滤器 `blacklist`/`whitelist` 类型 过滤器 规则。",
									},
									"filters": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Ip 过滤器 规则 列表，支持 IPs 在 X.X.X.X 格式，或 /8，/16，/24 格式 IP ranges. Up 到 50 allowlists 或 blocklists 可以 是 entered。",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"rule_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Ip 过滤器 规则 类型 过滤器 规则，可用: `all`，`文件`，`directory`，`路径`。",
									},
									"rule_paths": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "内容 列表 对于 each `rule_type`: `*` 对于 `all`，文件 ext like `jpg` 对于 `文件`，`/dir/like/` 对于 `directory` 和 `/路径/索引html` 对于 `路径`。",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"return_code": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Return 代码，可用值：400-499。",
						},
					},
				},
			},
			"ip_freq_limit": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "指定Ip 频率 限制 configurations。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch，可用值：`在`，`关闭` (默认值)。",
						},
						"qps": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Sets limited 数量 requests per second，514 将 是 返回 对于 requests 该 exceed 限制",
						},
					},
				},
			},
			"status_code_cache": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "状态 代码 缓存 configurations。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch，可用值：`在`，`关闭` (默认值)。",
						},
						"cache_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 缓存 规则。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status_code": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "代码 的 状态 缓存. 可用值：`403`，`404`。",
									},
									"cache_time": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "状态 代码 缓存 过期时间 (在 秒)。",
									},
								},
							},
						},
					},
				},
			},
			"compression": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Smart 压缩 configurations。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch，可用值：`在`，`关闭` (默认值)。",
						},
						"compression_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 压缩 规则。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"compress": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Must 是 集合 作为 true，enables 压缩。",
									},
									"min_length": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "最小 文件 大小 到 触发器 压缩 (在 bytes)。",
									},
									"max_length": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "最大 文件 大小 到 触发器 压缩 (在 bytes)。",
									},
									"algorithms": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "列表 algorithms，可用: `gzip` 和 `brotli`。",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"file_extensions": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "列表 文件 extensions like `jpg`，`txt`。",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"rule_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Rule 类型，可用: `all`，`文件`，`directory`，`路径`，`contentType`。",
									},
									"rule_paths": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "列表 规则 paths 对于 each `rule_type`: `*` 对于 `all`，文件 ext like `jpg` 对于 `文件`，`/dir/like/` 对于 `directory` 和 `/路径/索引html` 对于 `路径`。",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"band_width_alert": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Bandwidth cap 配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch，可用值：`在`，`关闭` (默认值)。",
						},
						"bps_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "阈值 的 bps。",
						},
						"counter_measure": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Counter measure。",
						},
						"last_trigger_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last 触发器 时间。",
						},
						"alert_switch": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Switch alert。",
						},
						"alert_percentage": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Alert percentage。",
						},
						"last_trigger_time_overseas": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last 触发器 时间 的 overseas。",
						},
						"metric": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Metric。",
						},
						"statistic_item": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "指定statistic item 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Configuration switch，可用值：`在`，`关闭` (默认值)。",
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "类型 statistic item。",
									},
									"unblock_time": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Time 的 auto unblock。",
									},
									"bps_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "阈值 的 bps。",
									},
									"counter_measure": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Counter measure，值: `RETURN_404`，`RESOLVE_DNS_TO_ORIGIN`。",
									},
									"alert_switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Switch alert。",
									},
									"alert_percentage": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Alert percentage。",
									},
									"metric": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Metric。",
									},
									"cycle": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Cycle 的 checking 在 minutes，值 `60`，`1440`。",
									},
								},
							},
						},
					},
				},
			},
			"error_page": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "错误 页面 configurations。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch，可用值：`在`，`关闭` (默认值)。",
						},
						"page_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 错误 页面 规则。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status_code": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "状态 代码 的 错误 页面 规则。",
									},
									"redirect_code": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Redirect 代码 的 错误 页面 规则。",
									},
									"redirect_url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Redirect URL 的 错误 页面 规则。",
									},
								},
							},
						},
					},
				},
			},
			"response_header": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Response 头部 configurations。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch，可用值：`在`，`关闭` (默认值)。",
						},
						"header_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 response 头部 规则。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"header_mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Response 头部 模式",
									},
									"header_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "response 头部 名称 规则。",
									},
									"header_value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "response 头部 值 的 规则。",
									},
									"rule_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "response 规则 类型 规则。",
									},
									"rule_paths": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "response 规则 paths 的 规则。",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"downstream_capping": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Downstream capping 配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch，可用值：`在`，`关闭` (默认值)。",
						},
						"capping_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 capping 规则。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Capping 规则 类型",
									},
									"rule_paths": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "列表 capping 规则 路径",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"kbps_threshold": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Capping 规则 kbps 阈值。",
									},
								},
							},
						},
					},
				},
			},
			"response_header_cache_switch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Response 头部 缓存 switch，可用值：`在`，`关闭` (默认值)。",
			},
			"origin_pull_optimization": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Cross-border linkage optimization 配置. (此 功能 是 在 beta 和 不 generally 可用 yet)。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch，可用值：`在`，`关闭` (默认值)。",
						},
						"optimization_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Optimization 类型，值: `OVToCN` - Overseas 到 CN，`CNToOV` CN 到 Overseas。",
						},
					},
				},
			},
			"seo_switch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SEO switch，可用值：`在`，`关闭` (默认值)。",
			},
			"referer": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Referer 配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch，可用值：`在`，`关闭` (默认值)。",
						},
						"referer_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 referer 规则。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Referer 规则 类型",
									},
									"rule_paths": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Referer 规则 路径 列表。",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"referer_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Referer 类型",
									},
									"referers": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Referer 列表。",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"allow_empty": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "是否allow emptpy。",
									},
								},
							},
						},
					},
				},
			},
			"video_seek_switch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Video seek switch，可用值：`在`，`关闭` (默认值)。",
			},
			"max_age": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Browser 缓存 配置. (此 功能 是 在 beta 和 不 generally 可用 yet)。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch，可用值：`在`，`关闭` (默认值)。",
						},
						"max_age_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "列表 Max Age 规则 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_age_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "following types 是 支持: `all`: all documents take effect，`文件`: 指定 文件 suffix takes effect，`directory`: 指定 路径 takes effect，`路径`: 指定absolute 路径 到 take effect，`索引`: home 页面。",
									},
									"max_age_contents": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "列表 规则 paths 对于 each `max_age_type`: `*` 对于 `all`，文件 ext like `jpg` 对于 `文件`，`/dir/like/` 对于 `directory` 和 `/路径/索引html` 对于 `路径`。",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"max_age_time": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Max Age 时间 （秒）， 此 可以 集合 到 `0` 该 stands 对于 无 缓存。",
									},
									"follow_origin": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "是否follow 源站，值: `在`/`关闭`，如果 集合 到 `在`， `max_age_time` 将 是 ignored。",
									},
								},
							},
						},
					},
				},
			},
			"specific_config_mainland": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Specific 配置 对于 mainland，NOTE: Both specifying full schema 或 使用 它 是 superfluous，please 使用 云 api 参数 json passthroughs，check [Data Types](https://www.tencentcloud.com/document/api/228/31739#MainlandConfig) 对于 more details。",
				DiffSuppressFunc: helper.DiffSupressJSON,
			},
			"specific_config_overseas": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Specific 配置 对于 oversea，NOTE: Both specifying full schema 或 使用 它 是 superfluous，please 使用 云 api 参数 json passthroughs，check [Data Types](https://www.tencentcloud.com/document/api/228/31739#OverseaConfig) 对于 more details。",
				DiffSuppressFunc: helper.DiffSupressJSON,
			},
			"origin_pull_timeout": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Cross-border linkage optimization 配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connect_timeout": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "源站-pull 连接 超时 (在 秒). 有效 范围: 5-60。",
						},
						"receive_timeout": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "源站-pull receipt 超时 (在 秒). 有效 范围: 10-60。",
						},
					},
				},
			},
			"offline_cache_switch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Offline 缓存 switch，可用值：`在`，`关闭` (默认值)。",
			},
			"post_max_size": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Maximum post 大小 配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch，可用值：`在`，`关闭` (默认值)。",
						},
						"max_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Maximum 大小 （MB）， 值 范围 是 `[1，200]`。",
						},
					},
				},
			},
			"quic_switch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "QUIC switch，可用值：`在`，`关闭` (默认值)。",
			},
			"cache_key": {
				Optional:      true,
				Type:          schema.TypeList,
				MaxItems:      1,
				ConflictsWith: []string{"full_url_cache"},
				Description:   "Cache 键 配置 (Ignore Query String 配置). NOTE: All 的 `full_url_cache` 默认值为 `在`。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"full_url_cache": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     CDN_SWITCH_ON,
							Description: "是否enable full-路径 缓存，值 `在` (DEFAULT ON)，`关闭`。",
						},
						"ignore_case": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     CDN_SWITCH_OFF,
							Description: "指定是否cache 键 是 case sensitive。",
						},
						"query_string": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Request 参数 contained 在 CacheKey。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     CDN_SWITCH_OFF,
										Description: "是否use QueryString 作为 part 的 CacheKey，值 `在`，`关闭` (Default)。",
									},
									"reorder": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     CDN_SWITCH_OFF,
										Description: "是否sort again，值 `在`，`关闭` (Default)。",
									},
									"action": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Include/exclude 查询 参数. Values: `includeAll` (Default)，`excludeAll`，`includeCustom`，`excludeCustom`。",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "数组 included/excluded 查询 strings (separated 通过 `;`)。",
									},
								},
							},
						},
						"key_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "路径-特定 缓存 键 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_paths": {
										Type: schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Required:    true,
										Description: "列表 规则 paths 对于 each `key_rules`: `/` 对于 `索引`，文件 ext like `jpg` 对于 `文件`，`/dir/like/` 对于 `directory` 和 `/路径/索引html` 对于 `路径`。",
									},
									"rule_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Rule 类型，可用: `文件`，`directory`，`路径`，`索引`。",
									},
									"full_url_cache": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     CDN_SWITCH_ON,
										Description: "是否enable full-路径 缓存，值 `在` (DEFAULT ON)，`关闭`。",
									},
									"ignore_case": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     CDN_SWITCH_OFF,
										Description: "Whether caches 是 case insensitive。",
									},
									"query_string": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Required:    true,
										Description: "Request 参数 contained 在 CacheKey。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Optional:    true,
													Default:     CDN_SWITCH_OFF,
													Description: "是否use QueryString 作为 part 的 CacheKey，值 `在`，`关闭` (Default)。",
												},
												"action": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "指定key 规则 QS 操作，值: `includeCustom`，`excludeCustom`。",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Default:     "",
													Description: "数组 included/excluded 查询 strings (separated 通过 `;`)。",
												},
											},
										},
									},
									"rule_tag": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "指定rule 标签，默认值为 `用户`。",
									},
								},
							},
						},
					},
				},
			},
			"aws_private_access": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Access authentication 对于 S3 源站。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch，可用值：`在`，`关闭` (默认值)。",
						},
						"access_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Access ID。",
							Sensitive:   true,
						},
						"secret_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "键",
							Sensitive:   true,
						},
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "地域",
						},
						"bucket": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "存储桶",
						},
					},
				},
			},
			"oss_private_access": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Access authentication 对于 OSS 源站。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch，可用值：`在`，`关闭` (默认值)。",
						},
						"access_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Access ID。",
							Sensitive:   true,
						},
						"secret_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "键",
							Sensitive:   true,
						},
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "地域",
						},
						"bucket": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "存储桶",
						},
					},
				},
			},
			"hw_private_access": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Access authentication 对于 OBS 源站。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch，可用值：`在`，`关闭` (默认值)。",
						},
						"access_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Access ID。",
							Sensitive:   true,
						},
						"secret_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "键",
							Sensitive:   true,
						},
						"bucket": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "存储桶",
						},
					},
				},
			},
			"qn_private_access": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Access authentication 对于 OBS 源站。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch，可用值：`在`，`关闭` (默认值)。",
						},
						"access_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Access ID。",
							Sensitive:   true,
						},
						"secret_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "键",
							Sensitive:   true,
						},
					},
				},
			},
			"others_private_access": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Object 存储 back-到-来源 authentication 的 other vendors。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration switch，可用值：`在`，`关闭` (默认值)。",
						},
						"access_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Access ID。",
							Sensitive:   true,
						},
						"secret_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "键",
							Sensitive:   true,
						},
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "地域",
						},
						"bucket": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "存储桶",
						},
					},
				},
			},
			"https_billing": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "HTTPS 服务 是 已启用 通过 默认值 (此 是 paid 服务; please refer 到 billing 信息 和 product documentation 对于 details)。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "HTTPS 服务 配置 switch，possible 值 是: 在: 已启用 (默认值 setting)，将 incur charges; 关闭: 已禁用，将 block HTTPS requests。",
						},
					},
				},
			},
			"user_agent_filter": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "UserAgent blacklist/whitelist 配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "Configuration switch，有效 值 是 `在` 和 `关闭`。",
						},
						"filter_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "UA blacklist/whitelist effect 规则 列表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Rule 类型，有效值：`all`，`文件`，`directory`，`路径`。",
									},
									"rule_paths": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Rule paths。",
									},
									"user_agents": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "UserAgent 列表。",
									},
									"filter_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Blacklist 或 whitelist，有效值：`blacklist`，`whitelist`。",
									},
								},
							},
						},
					},
				},
			},
			"url_redirect": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "URL redirect 配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "Configuration switch，有效 值 是 `在` 和 `关闭`。",
						},
						"path_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "URL redirect 规则 列表，最大 10 规则。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"redirect_status_code": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: tccommon.ValidateAllowedIntValue([]int{301, 302}),
										Description:  "Redirect 状态 代码，有效值：`301`，`302`。",
									},
									"pattern": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "URL 路径 到 match，支持 wildcard `*`，max 长度 1024。",
									},
									"redirect_url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Target URL，必须 start 使用 `/`，max 长度 1024。",
									},
									"redirect_host": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Target 主机，必须 start 使用 `http://` 或 `https://`。",
									},
									"full_match": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "是否use full 路径 match。",
									},
								},
							},
						},
					},
				},
			},
			"origin_combine": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Origin combine 配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "Configuration switch，有效 值 是 `在` 和 `关闭`。",
						},
					},
				},
			},
			"range_origin_pull": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Range 源站 pull 配置 使用 路径-based 规则。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "Global 范围 源站 pull switch，有效 值 是 `在` 和 `关闭`。",
						},
						"range_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "路径-based 范围 源站 pull 规则。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
										Description:  "Rule switch，有效 值 是 `在` 和 `关闭`。",
									},
									"rule_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Rule 类型，有效值：`文件`，`directory`，`路径`。",
									},
									"rule_paths": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Rule paths。",
									},
								},
							},
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 的 cdn 域名",
			},

			// computed
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Acceleration 服务 状态",
			},
			"cname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CNAME 地址 的 域名 名称",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "创建时间 的 域名 名称",
			},
			"explicit_using_dry_run": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "用于validate 仅 通过 store arguments 到 请求 json 字符串 作为 expected，WARNING: 如果 集合 到 `true`，NO Cloud Api 将 是 invoked 但 store 作为 本地 数据，do 不 使用 此 argument unless 您 really know what 您 是 doing。",
			},
			"dry_run_create_result": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "用于store `dry_run` 请求 json。",
			},
			"access_port": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Access 端口 配置. 列表 ports 该 可以 是 accessed。",
			},
			"dry_run_update_result": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "用于store `dry_run` update 请求 json。",
			},
			"auto_guard": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Traffic anti-hotlinking protection 配置. 注意: Create API does 不 support 此 字段，它 将 是 集合 via Update API after creation。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "AutoGuard switch，有效 值 是 `在` 和 `关闭`。",
						},
						"filter_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "AutoGuard 过滤器 规则。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filter_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Block 类型 `forbidden`: block。",
									},
									"rule_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Block 规则 类型 `all`: all requests; `文件`: 文件 requests 使用 指定 suffix。",
									},
									"rule_paths": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Block 规则 paths。",
									},
								},
							},
						},
					},
				},
			},
			"geo_blocker": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Regional 访问 control 配置. 注意: Create API does 不 support 此 字段，它 将 是 集合 via Update API after creation。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CDN_SWITCH),
							Description:  "GeoBlocker switch，有效 值 是 `在` 和 `关闭`。",
						},
						"block_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "GeoBlocker block 规则。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"block_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Rule 类型 `whitelist`: whitelist; `blacklist`: blacklist。",
									},
									"rule_paths": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Rule paths。",
									},
									"rule_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Rule effective 类型 `all`: all; `directory`: directory。",
									},
									"districts": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Effective districts，e.g. `CN-HK`，`CN-BJ`，etc。",
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

func resourceTencentCloudCdnDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdn_domain.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cdnService := CdnService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	request := cdn.NewAddCdnDomainRequest()
	domain := d.Get("domain").(string)
	request.Domain = &domain
	request.ServiceType = helper.String(d.Get("service_type").(string))
	request.ProjectId = helper.IntInt64(d.Get("project_id").(int))
	if v, ok := d.GetOk("area"); ok {
		request.Area = helper.String(v.(string))
	}
	// Range Origin Pull
	request.RangeOriginPull = &cdn.RangeOriginPull{}
	request.RangeOriginPull.Switch = helper.String(d.Get("range_origin_switch").(string))
	// If range_origin_pull is configured, use its switch and range_rules
	if v, ok := helper.InterfacesHeadMap(d, "range_origin_pull"); ok {
		request.RangeOriginPull.Switch = helper.String(v["switch"].(string))
		if rules, ok := v["range_rules"].([]interface{}); ok && len(rules) > 0 {
			rangeRules := make([]*cdn.RangeOriginPullRule, 0, len(rules))
			for _, rule := range rules {
				ruleMap := rule.(map[string]interface{})
				rangeRule := &cdn.RangeOriginPullRule{
					Switch: helper.String(ruleMap["switch"].(string)),
				}
				if rv, ok := ruleMap["rule_type"].(string); ok && rv != "" {
					rangeRule.RuleType = &rv
				}
				if rv, ok := ruleMap["rule_paths"].([]interface{}); ok && len(rv) > 0 {
					rangeRule.RulePaths = helper.InterfacesStringsPoint(rv)
				}
				rangeRules = append(rangeRules, rangeRule)
			}
			request.RangeOriginPull.RangeRules = rangeRules
		}
	}

	if v, ok := d.GetOk("ipv6_access_switch"); ok {
		request.Ipv6Access = &cdn.Ipv6Access{
			Switch: helper.String(v.(string)),
		}
	}

	if v, ok := d.GetOk("follow_redirect_switch"); ok {
		request.FollowRedirect = &cdn.FollowRedirect{
			Switch: helper.String(v.(string)),
		}
	}

	// access_port - Note: AccessPort is only supported in UpdateDomainConfigRequest, not AddCdnDomainRequest
	// This will be handled in the Update function

	if v, ok := helper.InterfacesHeadMap(d, "authentication"); ok {
		switchOn := v["switch"].(string)
		request.Authentication = &cdn.Authentication{
			Switch: &switchOn,
		}

		if v, ok := v["type_a"].([]interface{}); ok && len(v) > 0 {
			var (
				item           = v[0].(map[string]interface{})
				secretKey      = item["secret_key"].(string)
				signParam      = item["sign_param"].(string)
				expireTime     = item["expire_time"].(int)
				fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
				filterType     = item["filter_type"].(string)
			)

			request.Authentication.TypeA = &cdn.AuthenticationTypeA{
				SecretKey:      &secretKey,
				SignParam:      &signParam,
				ExpireTime:     helper.IntInt64(expireTime),
				FileExtensions: fileExtensions,
				FilterType:     &filterType,
			}

			if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
				request.Authentication.TypeA.BackupSecretKey = &backupSecretKey
			}
		}

		if v, ok := v["type_b"].([]interface{}); ok && len(v) > 0 {
			var (
				item           = v[0].(map[string]interface{})
				secretKey      = item["secret_key"].(string)
				expireTime     = item["expire_time"].(int)
				fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
				filterType     = item["filter_type"].(string)
			)

			request.Authentication.TypeB = &cdn.AuthenticationTypeB{
				SecretKey:      &secretKey,
				ExpireTime:     helper.IntInt64(expireTime),
				FileExtensions: fileExtensions,
				FilterType:     &filterType,
			}

			if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
				request.Authentication.TypeB.BackupSecretKey = &backupSecretKey
			}
		}

		if v, ok := v["type_c"].([]interface{}); ok && len(v) > 0 {
			var (
				item           = v[0].(map[string]interface{})
				secretKey      = item["secret_key"].(string)
				expireTime     = item["expire_time"].(int)
				fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
				filterType     = item["filter_type"].(string)
			)

			request.Authentication.TypeC = &cdn.AuthenticationTypeC{
				SecretKey:      &secretKey,
				ExpireTime:     helper.IntInt64(expireTime),
				FileExtensions: fileExtensions,
				FilterType:     &filterType,
			}

			if timeFormat, ok := item["time_format"].(string); ok {
				request.Authentication.TypeC.TimeFormat = &timeFormat
			}

			if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
				request.Authentication.TypeC.BackupSecretKey = &backupSecretKey
			}
		}

		if v, ok := v["type_d"].([]interface{}); ok && len(v) > 0 {
			var (
				item           = v[0].(map[string]interface{})
				secretKey      = item["secret_key"].(string)
				expireTime     = item["expire_time"].(int)
				fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
				filterType     = item["filter_type"].(string)
				timeParam      = item["time_param"].(string)
			)

			request.Authentication.TypeD = &cdn.AuthenticationTypeD{
				SecretKey:      &secretKey,
				ExpireTime:     helper.IntInt64(expireTime),
				FileExtensions: fileExtensions,
				FilterType:     &filterType,
				TimeParam:      &timeParam,
			}

			if timeFormat, ok := item["time_format"].(string); ok {
				request.Authentication.TypeD.TimeFormat = &timeFormat
			}

			if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
				request.Authentication.TypeD.BackupSecretKey = &backupSecretKey
			}
		}
	}

	// rule_cache
	if v, ok := d.GetOk("rule_cache"); ok {
		ruleCache := v.([]interface{})
		var ruleCaches []*cdn.RuleCache
		for _, v := range ruleCache {
			re := &cdn.RuleCache{}
			ruleCacheMap := v.(map[string]interface{})
			rulePaths := ruleCacheMap["rule_paths"].([]interface{})
			rulePathList := make([]*string, 0, len(rulePaths))
			ruleType := ruleCacheMap["rule_type"].(string)
			if ruleType == CDN_RULE_TYPE_DEFAULT {
				rulePathList = append(rulePathList, helper.String(CDN_RULE_PATH))
			} else {
				for _, value := range rulePaths {
					rulePathList = append(rulePathList, helper.String(value.(string)))
				}
			}
			switchFlag := ruleCacheMap["switch"].(string)
			cacheTime := ruleCacheMap["cache_time"].(int)
			compareMaxAge := ruleCacheMap["compare_max_age"].(string)
			ignoreCacheControl := ruleCacheMap["ignore_cache_control"].(string)
			ignoreSetCookie := ruleCacheMap["ignore_set_cookie"].(string)
			noCacheSwitch := ruleCacheMap["no_cache_switch"].(string)
			reValidate := ruleCacheMap["re_validate"].(string)
			followOriginSwitch := ruleCacheMap["follow_origin_switch"].(string)
			ruleCacheConfig := &cdn.RuleCacheConfig{}
			cache := &cdn.CacheConfigCache{}
			noCache := &cdn.CacheConfigNoCache{}
			followOrigin := &cdn.CacheConfigFollowOrigin{}
			ruleCacheConfig.Cache = cache
			ruleCacheConfig.NoCache = noCache
			ruleCacheConfig.FollowOrigin = followOrigin
			re.CacheConfig = ruleCacheConfig
			re.RulePaths = rulePathList
			re.RuleType = &ruleType
			re.CacheConfig.Cache.Switch = &switchFlag
			re.CacheConfig.Cache.CacheTime = helper.IntInt64(cacheTime)
			re.CacheConfig.Cache.CompareMaxAge = &compareMaxAge
			re.CacheConfig.Cache.IgnoreCacheControl = &ignoreCacheControl
			re.CacheConfig.Cache.IgnoreSetCookie = &ignoreSetCookie
			re.CacheConfig.NoCache.Switch = &noCacheSwitch
			re.CacheConfig.NoCache.Revalidate = &reValidate
			re.CacheConfig.FollowOrigin.Switch = &followOriginSwitch
			heuristicCacheSwitch := ruleCacheMap["heuristic_cache_switch"].(string)
			heuristicCacheTime := ruleCacheMap["heuristic_cache_time"].(int)
			if heuristicCacheSwitch != "" {
				re.CacheConfig.FollowOrigin.HeuristicCache = &cdn.HeuristicCache{
					Switch:      &heuristicCacheSwitch,
					CacheConfig: &cdn.CacheConfig{},
				}
				if heuristicCacheTime > 0 {
					re.CacheConfig.FollowOrigin.HeuristicCache.CacheConfig.HeuristicCacheTimeSwitch = helper.String("on")
					re.CacheConfig.FollowOrigin.HeuristicCache.CacheConfig.HeuristicCacheTime = helper.IntInt64(heuristicCacheTime)
				} else {
					re.CacheConfig.FollowOrigin.HeuristicCache.CacheConfig.HeuristicCacheTimeSwitch = helper.String("off")
				}
			}
			ruleCaches = append(ruleCaches, re)
		}
		request.Cache = &cdn.Cache{}
		request.Cache.RuleCache = ruleCaches
	}

	if v, ok := d.GetOk("request_header"); ok {
		requestHeaders := v.([]interface{})
		requestHeader := requestHeaders[0].(map[string]interface{})
		headerRule := requestHeader["header_rules"].([]interface{})
		var headerRules []*cdn.HttpHeaderPathRule
		for _, value := range headerRule {
			hr := &cdn.HttpHeaderPathRule{}
			headerRuleMap := value.(map[string]interface{})
			headerMode := headerRuleMap["header_mode"].(string)
			headerName := headerRuleMap["header_name"].(string)
			headerValue := headerRuleMap["header_value"].(string)
			ruleType := headerRuleMap["rule_type"].(string)
			rulePaths := headerRuleMap["rule_paths"].([]interface{})
			rulePathList := make([]*string, 0, len(rulePaths))
			for _, value := range rulePaths {
				rulePathList = append(rulePathList, helper.String(value.(string)))
			}
			hr.HeaderMode = &headerMode
			hr.HeaderName = &headerName
			hr.HeaderValue = &headerValue
			hr.RuleType = &ruleType
			hr.RulePaths = rulePathList
			headerRules = append(headerRules, hr)
		}
		request.RequestHeader = &cdn.RequestHeader{}
		request.RequestHeader.Switch = helper.String(requestHeader["switch"].(string))
		request.RequestHeader.HeaderRules = headerRules
	}

	// origin
	origins := d.Get("origin").([]interface{})
	if len(origins) < 1 {
		return fmt.Errorf("origin is required")
	}
	origin := origins[0].(map[string]interface{})
	request.Origin = &cdn.Origin{}
	request.Origin.OriginType = helper.String(origin["origin_type"].(string))
	originList := origin["origin_list"].([]interface{})
	request.Origin.Origins = make([]*string, 0, len(originList))
	for _, item := range originList {
		request.Origin.Origins = append(request.Origin.Origins, helper.String(item.(string)))
	}
	if v := origin["server_name"]; v.(string) != "" {
		request.Origin.ServerName = helper.String(v.(string))
	}
	if v := origin["cos_private_access"]; v.(string) != "" {
		request.Origin.CosPrivateAccess = helper.String(v.(string))
	}
	if v := origin["origin_pull_protocol"]; v.(string) != "" {
		request.Origin.OriginPullProtocol = helper.String(v.(string))
	}
	if v := origin["backup_origin_type"]; v.(string) != "" {
		request.Origin.BackupOriginType = helper.String(v.(string))
	}
	if v := origin["backup_server_name"]; v.(string) != "" {
		request.Origin.BackupServerName = helper.String(v.(string))
	}
	if v := origin["backup_origin_list"]; len(v.([]interface{})) > 0 {
		backupOriginList := v.([]interface{})
		request.Origin.BackupOrigins = make([]*string, 0, len(backupOriginList))
		for _, item := range backupOriginList {
			request.Origin.BackupOrigins = append(request.Origin.BackupOrigins, helper.String(item.(string)))
		}
	}
	if v := origin["origin_company"]; v.(string) != "" {
		request.Origin.OriginCompany = helper.String(v.(string))
	}

	// https config
	if v, ok := d.GetOk("https_config"); ok {
		httpsConfigs := v.([]interface{})
		if len(httpsConfigs) > 0 {
			config := httpsConfigs[0].(map[string]interface{})
			request.Https = &cdn.Https{}
			request.Https.Switch = helper.String(config["https_switch"].(string))
			if v := config["http2_switch"]; v.(string) != "" {
				request.Https.Http2 = helper.String(v.(string))
			}
			request.Https.OcspStapling = helper.String(config["ocsp_stapling_switch"].(string))
			request.Https.Spdy = helper.String(config["spdy_switch"].(string))
			request.Https.VerifyClient = helper.String(config["verify_client"].(string))
			if v := config["server_certificate_config"]; len(v.([]interface{})) > 0 {
				serverCerts := v.([]interface{})
				if len(serverCerts) > 0 && serverCerts[0] != nil {
					serverCert := serverCerts[0].(map[string]interface{})
					request.Https.CertInfo = &cdn.ServerCert{}
					if v := serverCert["certificate_id"]; v.(string) != "" {
						request.Https.CertInfo.CertId = helper.String(v.(string))
					}
					if v := serverCert["certificate_content"]; v.(string) != "" {
						request.Https.CertInfo.Certificate = helper.String(v.(string))
					}
					if v := serverCert["private_key"]; v.(string) != "" {
						request.Https.CertInfo.PrivateKey = helper.String(v.(string))
					}
					if v := serverCert["message"]; v.(string) != "" {
						request.Https.CertInfo.Message = helper.String(v.(string))
					}
				}
			}
			if v := config["client_certificate_config"]; len(v.([]interface{})) > 0 {
				clientCerts := v.([]interface{})
				if len(clientCerts) > 0 && clientCerts[0] != nil {
					clientCert := clientCerts[0].(map[string]interface{})
					request.Https.ClientCertInfo = &cdn.ClientCert{}
					if v := clientCert["certificate_content"]; v.(string) != "" {
						request.Https.ClientCertInfo.Certificate = helper.String(v.(string))
					}
				}
			}
			if v, ok := config["force_redirect"]; ok {
				forceRedirect := v.([]interface{})
				if len(forceRedirect) > 0 && forceRedirect[0] != nil {
					var redirect cdn.ForceRedirect
					redirectMap := forceRedirect[0].(map[string]interface{})
					if sw := redirectMap["switch"]; sw.(string) != "" {
						redirect.Switch = helper.String(sw.(string))
					}
					if rt := redirectMap["redirect_type"]; rt.(string) != "" {
						redirect.RedirectType = helper.String(rt.(string))
					}
					if rsc := redirectMap["redirect_status_code"]; rsc.(int) != 0 {
						redirect.RedirectStatusCode = helper.Int64(int64(rsc.(int)))
					}
					if ch := redirectMap["carry_headers"]; ch.(string) != "" {
						redirect.CarryHeaders = helper.String(ch.(string))
					}
					request.ForceRedirect = &redirect
				}
			}
			if v, ok := config["tls_versions"]; ok {
				request.Https.TlsVersion = helper.InterfacesStringsPoint(v.([]interface{}))
			}
			// HSTS
			if v, ok := config["hsts"]; ok {
				hstsList := v.([]interface{})
				if len(hstsList) > 0 && hstsList[0] != nil {
					hstsMap := hstsList[0].(map[string]interface{})
					request.Https.Hsts = &cdn.Hsts{
						Switch: helper.String(hstsMap["switch"].(string)),
					}
					if maxAge, ok := hstsMap["max_age"].(int); ok && maxAge > 0 {
						request.Https.Hsts.MaxAge = helper.IntInt64(maxAge)
					}
					if includeSubDomains, ok := hstsMap["include_sub_domains"].(string); ok && includeSubDomains != "" {
						request.Https.Hsts.IncludeSubDomains = &includeSubDomains
					}
				}
			}
		}
	}

	// more added
	if v, ok := helper.InterfacesHeadMap(d, "ip_filter"); ok {
		vSwitch := v["switch"].(string)
		request.IpFilter = &cdn.IpFilter{
			Switch: &vSwitch,
		}
		if vv, ok := v["filter_type"].(string); ok {
			request.IpFilter.FilterType = &vv
		}
		if vv, ok := v["filters"].([]interface{}); ok {
			request.IpFilter.Filters = helper.InterfacesStringsPoint(vv)
		}

		//need white list func
		if vv, ok := v["filter_rules"].([]interface{}); ok && len(vv) > 0 {
			filterRules := make([]*cdn.IpFilterPathRule, 0)
			for i := range vv {
				item := vv[i].(map[string]interface{})
				rule := &cdn.IpFilterPathRule{}
				if rv, ok := item["filter_type"].(string); ok && rv != "" {
					rule.FilterType = &rv
				}
				if rv, ok := item["filters"].([]interface{}); ok && len(rv) > 0 {
					rule.Filters = helper.InterfacesStringsPoint(rv)
				}
				if rv, ok := item["rule_type"].(string); ok && rv != "" {
					rule.RuleType = &rv
				}
				if rv, ok := item["rule_paths"].([]interface{}); ok && len(rv) > 0 {
					rule.RulePaths = helper.InterfacesStringsPoint(rv)
				}
				filterRules = append(filterRules, rule)
			}
			request.IpFilter.FilterRules = filterRules
		}

		if vv, ok := v["return_code"].(int); ok {
			request.IpFilter.ReturnCode = helper.IntInt64(vv)
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "ip_freq_limit"); ok {
		vSwitch := v["switch"].(string)
		qps := v["qps"].(int)
		request.IpFreqLimit = &cdn.IpFreqLimit{
			Switch: &vSwitch,
			Qps:    helper.IntInt64(qps),
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "status_code_cache"); ok {
		vSwitch := v["switch"].(string)
		request.StatusCodeCache = &cdn.StatusCodeCache{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.StatusCodeCacheRule, 0)
		rules := v["cache_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.StatusCodeCacheRule{
				StatusCode: helper.String(item["status_code"].(string)),
			}
			if v, ok := item["cache_time"].(int); ok && v > 0 {
				rule.CacheTime = helper.IntInt64(v)
			}
			requestRules = append(requestRules, rule)
		}
		request.StatusCodeCache.CacheRules = requestRules
	}
	if v, ok := helper.InterfacesHeadMap(d, "compression"); ok {
		vSwitch := v["switch"].(string)
		request.Compression = &cdn.Compression{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.CompressionRule, 0)
		rules := v["compression_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			var (
				compress = item["compress"].(bool)
			)
			rule := &cdn.CompressionRule{
				Compress: &compress,
			}
			if v, ok := item["min_length"].(int); ok && v > 0 {
				rule.MinLength = helper.IntInt64(v)
			}
			if v, ok := item["max_length"].(int); ok && v > 0 {
				rule.MaxLength = helper.IntInt64(v)
			}
			if v, ok := item["algorithms"].([]interface{}); ok && len(v) > 0 {
				rule.Algorithms = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["file_extensions"].([]interface{}); ok && len(v) > 0 {
				rule.FileExtensions = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["rule_type"].(string); ok && v != "" {
				rule.RuleType = &v
			}
			if v, ok := item["rule_paths"].([]interface{}); ok && len(v) > 0 {
				rule.RulePaths = helper.InterfacesStringsPoint(v)
			}

			requestRules = append(requestRules, rule)
		}
		request.Compression.CompressionRules = requestRules

	}
	if v, ok := helper.InterfacesHeadMap(d, "band_width_alert"); ok {
		vSwitch := v["switch"].(string)
		request.BandwidthAlert = &cdn.BandwidthAlert{
			Switch: &vSwitch,
		}
		if v, ok := v["bps_threshold"].(int); ok && v > 0 {
			request.BandwidthAlert.BpsThreshold = helper.IntInt64(v)
		}
		if v, ok := v["counter_measure"].(string); ok && v != "" {
			request.BandwidthAlert.CounterMeasure = &v
		}
		//if v, ok := v["last_trigger_time"].(string); ok && v != "" {
		//	request.BandwidthAlert.LastTriggerTime = &v
		//}
		if v, ok := v["alert_switch"].(string); ok && v != "" {
			request.BandwidthAlert.AlertSwitch = &v
		}
		if v, ok := v["alert_percentage"].(int); ok && v > 0 {
			request.BandwidthAlert.AlertPercentage = helper.IntInt64(v)
		}
		//if v, ok := v["last_trigger_time_overseas"].(string); ok && v != "" {
		//	request.BandwidthAlert.LastTriggerTimeOverseas = &v
		//}
		if v, ok := v["metric"].(string); ok && v != "" {
			request.BandwidthAlert.Metric = &v
		}
		if statistic, ok := v["statistic_item"].([]interface{}); ok && len(statistic) > 0 {
			for i := range statistic {
				item := statistic[i].(map[string]interface{})
				vSwitch := item["switch"].(string)
				sItem := &cdn.StatisticItem{
					Switch: &vSwitch,
				}
				if vv, ok := item["type"].(string); ok && vv != "" {
					sItem.Type = &vv
				}
				if vv, ok := item["unblock_time"].(int); ok && vv != 0 {
					sItem.UnBlockTime = helper.IntUint64(vv)
				}
				if vv, ok := item["bps_threshold"].(int); ok && vv != 0 {
					sItem.BpsThreshold = helper.IntUint64(vv)
				}
				if vv, ok := item["counter_measure"].(string); ok && vv != "" {
					sItem.CounterMeasure = &vv
				}
				if vv, ok := item["alert_switch"].(string); ok && vv != "" {
					sItem.AlertSwitch = &vv
				}
				if vv, ok := item["alert_percentage"].(int); ok && vv != 0 {
					sItem.AlertPercentage = helper.IntUint64(vv)
				}
				if vv, ok := item["metric"].(string); ok && vv != "" {
					sItem.Metric = &vv
				}
				if vv, ok := item["cycle"].(int); ok && vv != 0 {
					sItem.BpsThreshold = helper.IntUint64(vv)
				}
				request.BandwidthAlert.StatisticItems = append(request.BandwidthAlert.StatisticItems, sItem)
			}

		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "error_page"); ok {
		vSwitch := v["switch"].(string)
		request.ErrorPage = &cdn.ErrorPage{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.ErrorPageRule, 0)
		rules := v["page_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.ErrorPageRule{}
			if v, ok := item["status_code"].(int); ok && v != 0 {
				rule.StatusCode = helper.IntInt64(v)
			}
			if v, ok := item["redirect_code"].(int); ok && v != 0 {
				rule.RedirectCode = helper.IntInt64(v)
			}
			if v, ok := item["redirect_url"].(string); ok && v != "" {
				rule.RedirectUrl = &v
			}
			requestRules = append(requestRules, rule)
		}
		request.ErrorPage.PageRules = requestRules
	}
	if v, ok := helper.InterfacesHeadMap(d, "response_header"); ok {
		vSwitch := v["switch"].(string)
		request.ResponseHeader = &cdn.ResponseHeader{
			Switch: &vSwitch,
		}
		responseRules := make([]*cdn.HttpHeaderPathRule, 0)
		rules := v["header_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.HttpHeaderPathRule{}
			if v, ok := item["header_mode"].(string); ok && v != "" {
				rule.HeaderMode = &v
			}
			if v, ok := item["header_name"].(string); ok && v != "" {
				rule.HeaderName = &v
			}
			if v, ok := item["header_value"].(string); ok && v != "" {
				rule.HeaderValue = &v
			}
			if v, ok := item["rule_type"].(string); ok && v != "" {
				rule.RuleType = &v
			}
			if v, ok := item["rule_paths"].([]interface{}); ok && len(v) > 0 {
				rule.RulePaths = helper.InterfacesStringsPoint(v)
			}
			responseRules = append(responseRules, rule)
		}
		request.ResponseHeader.HeaderRules = responseRules
	}
	if v, ok := helper.InterfacesHeadMap(d, "downstream_capping"); ok {
		vSwitch := v["switch"].(string)
		request.DownstreamCapping = &cdn.DownstreamCapping{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.CappingRule, 0)
		rules := v["capping_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.CappingRule{}
			if v, ok := item["rule_type"].(string); ok && v != "" {
				rule.RuleType = &v
			}
			if v, ok := item["rule_paths"].([]interface{}); ok && len(v) > 0 {
				rule.RulePaths = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["kbps_threshold"].(int); ok && v > 0 {
				rule.KBpsThreshold = helper.IntInt64(v)
			}
			requestRules = append(requestRules, rule)
		}
		request.DownstreamCapping.CappingRules = requestRules
	}
	if v, ok := d.GetOk("response_header_cache_switch"); ok {
		request.ResponseHeaderCache = &cdn.ResponseHeaderCache{
			Switch: helper.String(v.(string)),
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "origin_pull_optimization"); ok {
		vSwitch := v["switch"].(string)
		request.OriginPullOptimization = &cdn.OriginPullOptimization{
			Switch: &vSwitch,
		}
		if v, ok := v["optimization_type"].(string); ok && v != "" {
			request.OriginPullOptimization.OptimizationType = &v
		}
	}
	if v, ok := d.GetOk("seo_switch"); ok {
		request.Seo = &cdn.Seo{
			Switch: helper.String(v.(string)),
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "referer"); ok {
		vSwitch := v["switch"].(string)
		request.Referer = &cdn.Referer{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.RefererRule, 0)
		rules := v["referer_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.RefererRule{}
			if v, ok := item["rule_type"].(string); ok && v != "" {
				rule.RuleType = &v
			}
			if v, ok := item["rule_paths"].([]interface{}); ok && len(v) > 0 {
				rule.RulePaths = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["referer_type"].(string); ok && v != "" {
				rule.RefererType = &v
			}
			if v, ok := item["referers"].([]interface{}); ok && len(v) > 0 {
				rule.Referers = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["allow_empty"].(bool); ok {
				rule.AllowEmpty = &v
			}
			requestRules = append(requestRules, rule)
		}
		request.Referer.RefererRules = requestRules
	}
	if v, ok := d.GetOk("video_seek_switch"); ok {
		request.VideoSeek = &cdn.VideoSeek{
			Switch: helper.String(v.(string)),
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "max_age"); ok {
		vSwitch := v["switch"].(string)
		request.MaxAge = &cdn.MaxAge{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.MaxAgeRule, 0)
		rules := v["max_age_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.MaxAgeRule{}

			if v, ok := item["max_age_type"].(string); ok && v != "" {
				rule.MaxAgeType = &v
			}
			if v, ok := item["max_age_contents"].([]interface{}); ok && len(v) > 0 {
				rule.MaxAgeContents = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["max_age_time"].(int); ok && v > 0 {
				rule.MaxAgeTime = helper.IntInt64(v)
			}
			if v, ok := item["follow_origin"].(string); ok && v != "" {
				rule.FollowOrigin = &v
			}

			requestRules = append(requestRules, rule)
		}
		request.MaxAge.MaxAgeRules = requestRules
	}
	if v, ok := d.GetOk("specific_config_mainland"); ok && v.(string) != "" {
		request.SpecificConfig = &cdn.SpecificConfig{}
		configStruct := cdn.MainlandConfig{}
		err := json.Unmarshal([]byte(v.(string)), &configStruct)
		if err != nil {
			return fmt.Errorf("unmarshal specific_config_mainland fail: %s", err.Error())
		}
		request.SpecificConfig.Mainland = &configStruct
	}
	if v, ok := d.GetOk("specific_config_overseas"); ok && v.(string) != "" {
		if request.SpecificConfig == nil {
			request.SpecificConfig = &cdn.SpecificConfig{}
		}
		configStruct := cdn.OverseaConfig{}
		err := json.Unmarshal([]byte(v.(string)), &configStruct)
		if err != nil {
			return fmt.Errorf("unmarshal specific_config_overseas fail: %s", err.Error())
		}
		request.SpecificConfig.Overseas = &configStruct
	}
	if v, ok := helper.InterfacesHeadMap(d, "origin_pull_timeout"); ok {
		request.OriginPullTimeout = &cdn.OriginPullTimeout{}
		if vv, ok := v["connect_timeout"].(int); ok && vv > 0 {
			request.OriginPullTimeout.ConnectTimeout = helper.IntUint64(vv)
		}
		if vv, ok := v["receive_timeout"].(int); ok && vv > 0 {
			request.OriginPullTimeout.ReceiveTimeout = helper.IntUint64(vv)
		}
	}
	if v, ok := d.GetOk("offline_cache_switch"); ok {
		request.OfflineCache = &cdn.OfflineCache{
			Switch: helper.String(v.(string)),
		}
	}
	if v, ok := d.GetOk("quic_switch"); ok {
		request.Quic = &cdn.Quic{
			Switch: helper.String(v.(string)),
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "cache_key"); ok {
		request.CacheKey = &cdn.CacheKey{}
		if fuc := v["full_url_cache"].(string); fuc != "" {
			request.CacheKey.FullUrlCache = &fuc
		}
		if ic := v["ignore_case"].(string); ic != "" {
			request.CacheKey.IgnoreCase = &ic
		}
		if qs, ok := v["query_string"].([]interface{}); ok && len(qs) > 0 {
			if dMap, ok := qs[0].(map[string]interface{}); ok {
				qSwitch := dMap["switch"].(string)
				reorder := dMap["reorder"].(string)
				action := dMap["action"].(string)
				value := dMap["value"].(string)
				request.CacheKey.QueryString = &cdn.QueryStringKey{
					Switch:  &qSwitch,
					Reorder: &reorder,
					Action:  &action,
					Value:   &value,
				}
			}
		}
		if kr, ok := v["key_rules"].([]interface{}); ok {
			for i := range kr {
				rule, ok := kr[i].(map[string]interface{})
				if !ok {
					continue
				}
				ruleType := rule["rule_type"].(string)
				keyRule := &cdn.KeyRule{
					RuleType: &ruleType,
				}
				if vv := rule["full_url_cache"].(string); vv != "" {
					keyRule.FullUrlCache = &vv
				}
				if vv := rule["ignore_case"].(string); vv != "" {
					keyRule.IgnoreCase = &vv
				}
				if vv := rule["rule_tag"].(string); vv != "" {
					keyRule.RuleTag = &vv
				}
				if rp, ok := rule["rule_paths"].([]interface{}); ok {
					keyRule.RulePaths = helper.InterfacesStringsPoint(rp)
				}
				if qs, ok := rule["query_string"].([]interface{}); ok && len(qs) > 0 {
					if dMap, ok := qs[0].(map[string]interface{}); ok {
						vSwitch := dMap["switch"].(string)
						keyRule.QueryString = &cdn.RuleQueryString{
							Switch: &vSwitch,
						}
						if v := dMap["action"].(string); v != "" && vSwitch == "on" {
							keyRule.QueryString.Action = &v
						}
						if v := dMap["value"].(string); v != "" {
							keyRule.QueryString.Value = &v
						}
					}
				}
				request.CacheKey.KeyRules = append(request.CacheKey.KeyRules, keyRule)
			}
		}
	} else {
		fullUrlCache := d.Get("full_url_cache").(bool)
		request.CacheKey = &cdn.CacheKey{}
		if fullUrlCache {
			request.CacheKey.FullUrlCache = helper.String(CDN_SWITCH_ON)
		} else {
			request.CacheKey.FullUrlCache = helper.String(CDN_SWITCH_OFF)
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "aws_private_access"); ok {
		vSwitch := v["switch"].(string)
		request.AwsPrivateAccess = &cdn.AwsPrivateAccess{
			Switch: &vSwitch,
		}
		if v, ok := v["access_key"].(string); ok && v != "" {
			request.AwsPrivateAccess.AccessKey = &v
		}
		if v, ok := v["secret_key"].(string); ok && v != "" {
			request.AwsPrivateAccess.SecretKey = &v
		}
		if v, ok := v["region"].(string); ok && v != "" {
			request.AwsPrivateAccess.Region = &v
		}
		if v, ok := v["bucket"].(string); ok && v != "" {
			request.AwsPrivateAccess.Bucket = &v
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "oss_private_access"); ok {
		vSwitch := v["switch"].(string)
		request.OssPrivateAccess = &cdn.OssPrivateAccess{
			Switch: &vSwitch,
		}
		if v, ok := v["access_key"].(string); ok && v != "" {
			request.OssPrivateAccess.AccessKey = &v
		}
		if v, ok := v["secret_key"].(string); ok && v != "" {
			request.OssPrivateAccess.SecretKey = &v
		}
		if v, ok := v["region"].(string); ok && v != "" {
			request.OssPrivateAccess.Region = &v
		}
		if v, ok := v["bucket"].(string); ok && v != "" {
			request.OssPrivateAccess.Bucket = &v
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "hw_private_access"); ok {
		vSwitch := v["switch"].(string)
		request.HwPrivateAccess = &cdn.HwPrivateAccess{
			Switch: &vSwitch,
		}
		if v, ok := v["access_key"].(string); ok && v != "" {
			request.HwPrivateAccess.AccessKey = &v
		}
		if v, ok := v["secret_key"].(string); ok && v != "" {
			request.HwPrivateAccess.SecretKey = &v
		}
		if v, ok := v["bucket"].(string); ok && v != "" {
			request.HwPrivateAccess.Bucket = &v
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "qn_private_access"); ok {
		vSwitch := v["switch"].(string)
		request.QnPrivateAccess = &cdn.QnPrivateAccess{
			Switch: &vSwitch,
		}
		if v, ok := v["access_key"].(string); ok && v != "" {
			request.QnPrivateAccess.AccessKey = &v
		}
		if v, ok := v["secret_key"].(string); ok && v != "" {
			request.QnPrivateAccess.SecretKey = &v
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "others_private_access"); ok {
		vSwitch := v["switch"].(string)
		request.OthersPrivateAccess = &cdn.OthersPrivateAccess{
			Switch: &vSwitch,
		}
		if v, ok := v["access_key"].(string); ok && v != "" {
			request.OthersPrivateAccess.AccessKey = &v
		}
		if v, ok := v["secret_key"].(string); ok && v != "" {
			request.OthersPrivateAccess.SecretKey = &v
		}
		if v, ok := v["region"].(string); ok && v != "" {
			request.OthersPrivateAccess.Region = &v
		}
		if v, ok := v["bucket"].(string); ok && v != "" {
			request.OthersPrivateAccess.Bucket = &v
		}
	}

	if v, ok := helper.InterfacesHeadMap(d, "https_billing"); ok {
		vSwitch := v["switch"].(string)
		request.HttpsBilling = &cdn.HttpsBilling{
			Switch: &vSwitch,
		}
	}

	if v := d.Get("explicit_using_dry_run").(bool); v {
		d.SetId(domain)
		_ = d.Set("dry_run_create_result", request.ToJsonString())
		return nil
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().AddCdnDomain(request)
		if err != nil {
			if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == CDN_DOMAIN_CONFIG_ERROR || sdkErr.Code == CDN_HOST_EXISTS {
					return resource.NonRetryableError(err)
				}
			}
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(domain)

	time.Sleep(1 * time.Second)
	err = resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		domainConfig, err := cdnService.DescribeDomainsConfigByDomain(ctx, domain)
		if err != nil {
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		if *domainConfig.Status == CDN_DOMAIN_STATUS_PROCESSING {
			return resource.RetryableError(fmt.Errorf("domain status is still processing, retry..."))
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := updateCdnModifyOnlyParams(d, meta, ctx); err != nil {
		return err
	}

	// tags
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(client)
		region := client.Region
		resourceName := tccommon.BuildTagResourceName(CDN_SERVICE_NAME, CDN_RESOURCE_NAME_DOMAIN, region, domain)
		err := tagService.ModifyTags(ctx, resourceName, tags, nil)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudCdnDomainRead(d, meta)
}

func resourceTencentCloudCdnDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdn_domain.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	region := client.Region
	cdnService := CdnService{client: client}
	tagService := svctag.NewTagService(client)

	domain := d.Id()

	if v, ok := d.Get("explicit_using_dry_run").(bool); ok && v {
		d.SetId(domain)
		return nil
	}

	var domainConfig *cdn.DetailDomain
	var errRet error
	err := resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		domainConfig, errRet = cdnService.DescribeDomainsConfigByDomain(ctx, domain)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if domainConfig == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("domain", domain)
	_ = d.Set("service_type", domainConfig.ServiceType)
	_ = d.Set("project_id", domainConfig.ProjectId)
	_ = d.Set("area", domainConfig.Area)
	_ = d.Set("status", domainConfig.Status)
	_ = d.Set("create_time", domainConfig.CreateTime)
	_ = d.Set("cname", domainConfig.Cname)
	_ = d.Set("range_origin_switch", domainConfig.RangeOriginPull.Switch)

	if domainConfig.Ipv6Access != nil {
		_ = d.Set("ipv6_access_switch", domainConfig.Ipv6Access.Switch)
	}
	if domainConfig.FollowRedirect != nil {
		_ = d.Set("follow_redirect_switch", domainConfig.FollowRedirect.Switch)
	}

	origins := make([]map[string]interface{}, 0, 1)
	origin := make(map[string]interface{}, 8)
	origin["origin_type"] = domainConfig.Origin.OriginType
	origin["origin_list"] = domainConfig.Origin.Origins
	origin["server_name"] = domainConfig.Origin.ServerName
	origin["cos_private_access"] = domainConfig.Origin.CosPrivateAccess
	origin["origin_pull_protocol"] = domainConfig.Origin.OriginPullProtocol
	origin["backup_origin_type"] = domainConfig.Origin.BackupOriginType
	origin["backup_origin_list"] = domainConfig.Origin.BackupOrigins
	origin["backup_server_name"] = domainConfig.Origin.BackupServerName
	origin["origin_company"] = domainConfig.Origin.OriginCompany
	origins = append(origins, origin)
	_ = d.Set("origin", origins)

	if len(domainConfig.Cache.RuleCache) > 0 {
		ruleCaches := make([]map[string]interface{}, len(domainConfig.Cache.RuleCache))
		for index, value := range domainConfig.Cache.RuleCache {
			ruleCache := make(map[string]interface{})
			ruleCache["rule_paths"] = value.RulePaths
			ruleCache["rule_type"] = value.RuleType
			if value.CacheConfig == nil {
				continue
			}
			if value.CacheConfig.Cache != nil {
				ruleCache["switch"] = value.CacheConfig.Cache.Switch
				ruleCache["cache_time"] = value.CacheConfig.Cache.CacheTime
				ruleCache["compare_max_age"] = value.CacheConfig.Cache.CompareMaxAge
				ruleCache["ignore_cache_control"] = value.CacheConfig.Cache.IgnoreCacheControl
				ruleCache["ignore_set_cookie"] = value.CacheConfig.Cache.IgnoreSetCookie
			}
			if value.CacheConfig.NoCache != nil {
				ruleCache["no_cache_switch"] = value.CacheConfig.NoCache.Switch
				ruleCache["re_validate"] = value.CacheConfig.NoCache.Revalidate
			}
			if value.CacheConfig.FollowOrigin != nil {
				ruleCache["follow_origin_switch"] = value.CacheConfig.FollowOrigin.Switch
				if htc := value.CacheConfig.FollowOrigin.HeuristicCache; htc != nil {
					ruleCache["heuristic_cache_switch"] = htc.Switch
					if htc.CacheConfig != nil {
						ruleCache["heuristic_cache_time"] = htc.CacheConfig.HeuristicCacheTime
					}
				}
			}
			ruleCaches[index] = ruleCache
		}
		_ = d.Set("rule_cache", ruleCaches)
	}

	requestHeaders := make([]map[string]interface{}, 1)
	requestHeader := make(map[string]interface{})
	if domainConfig.RequestHeader != nil {
		requestHeader["switch"] = domainConfig.RequestHeader.Switch
		if len(domainConfig.RequestHeader.HeaderRules) > 0 {
			headerRules := make([]map[string]interface{}, len(domainConfig.RequestHeader.HeaderRules))
			headerRuleList := domainConfig.RequestHeader.HeaderRules
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
		requestHeaders[0] = requestHeader
		_ = d.Set("request_header", requestHeaders)
	}

	httpsConfigs := make([]map[string]interface{}, 0, 1)
	httpsConfig := make(map[string]interface{}, 8)
	httpsConfig["https_switch"] = domainConfig.Https.Switch
	httpsConfig["http2_switch"] = domainConfig.Https.Http2
	httpsConfig["ocsp_stapling_switch"] = domainConfig.Https.OcspStapling
	httpsConfig["spdy_switch"] = domainConfig.Https.Spdy
	httpsConfig["verify_client"] = domainConfig.Https.VerifyClient

	oldHttpsConfigs := make([]interface{}, 0)
	if _, ok := d.GetOk("https_config"); ok {
		oldHttpsConfigs = d.Get("https_config").([]interface{})
	}
	oldHttpsConfig := make(map[string]interface{})
	if len(oldHttpsConfigs) > 0 && oldHttpsConfigs[0] != nil {
		oldHttpsConfig = oldHttpsConfigs[0].(map[string]interface{})
	}
	oldServerConfigs := make([]interface{}, 0)
	if _, ok := oldHttpsConfig["server_certificate_config"]; ok {
		oldServerConfigs = oldHttpsConfig["server_certificate_config"].([]interface{})
	}
	oldServerConfig := make(map[string]interface{})
	if len(oldServerConfigs) > 0 && oldServerConfigs[0] != nil {
		oldServerConfig = oldServerConfigs[0].(map[string]interface{})
	}
	oldClientConfigs := make([]interface{}, 0)
	if _, ok := oldHttpsConfig["client_certificate_config"]; ok {
		oldClientConfigs = oldHttpsConfig["client_certificate_config"].([]interface{})
	}
	oldClientConfig := make(map[string]interface{})
	if len(oldClientConfigs) > 0 && oldClientConfigs[0] != nil {
		oldClientConfig = oldClientConfigs[0].(map[string]interface{})
	}

	if domainConfig.Https.CertInfo != nil && domainConfig.Https.CertInfo.CertName != nil {
		serverCertConfigs := make([]map[string]interface{}, 0, 1)
		serverCertConfig := make(map[string]interface{}, 5)
		serverCertConfig["certificate_id"] = domainConfig.Https.CertInfo.CertId
		serverCertConfig["certificate_name"] = domainConfig.Https.CertInfo.CertName
		serverCertConfig["certificate_content"] = oldServerConfig["certificate_content"]
		serverCertConfig["private_key"] = oldServerConfig["private_key"]
		serverCertConfig["message"] = domainConfig.Https.CertInfo.Message
		serverCertConfig["deploy_time"] = domainConfig.Https.CertInfo.DeployTime
		serverCertConfig["expire_time"] = domainConfig.Https.CertInfo.ExpireTime
		serverCertConfigs = append(serverCertConfigs, serverCertConfig)
		httpsConfig["server_certificate_config"] = serverCertConfigs
	}
	if domainConfig.Https.ClientCertInfo != nil && domainConfig.Https.ClientCertInfo.CertName != nil {
		clientCertConfigs := make([]map[string]interface{}, 0, 1)
		clientCertConfig := make(map[string]interface{}, 2)
		clientCertConfig["certificate_content"] = oldClientConfig["certificate_content"]
		clientCertConfig["certificate_name"] = domainConfig.Https.ClientCertInfo.CertName
		clientCertConfig["deploy_time"] = domainConfig.Https.ClientCertInfo.DeployTime
		clientCertConfig["expire_time"] = domainConfig.Https.ClientCertInfo.ExpireTime
		clientCertConfigs = append(clientCertConfigs, clientCertConfig)
		httpsConfig["client_certificate_config"] = clientCertConfigs
	}
	if domainConfig.ForceRedirect != nil {
		httpsConfig["force_redirect"] = []map[string]interface{}{
			{
				"switch":               domainConfig.ForceRedirect.Switch,
				"redirect_type":        domainConfig.ForceRedirect.RedirectType,
				"redirect_status_code": domainConfig.ForceRedirect.RedirectStatusCode,
				"carry_headers":        domainConfig.ForceRedirect.CarryHeaders,
			},
		}
	}
	if len(domainConfig.Https.TlsVersion) > 0 {
		tlsVersions := make([]string, 0)
		for _, tlsVersionItem := range domainConfig.Https.TlsVersion {
			tlsVersions = append(tlsVersions, *tlsVersionItem)
		}
		httpsConfig["tls_versions"] = tlsVersions
	}
	// HSTS
	if domainConfig.Https.Hsts != nil {
		hstsMap := map[string]interface{}{
			"switch": domainConfig.Https.Hsts.Switch,
		}
		if domainConfig.Https.Hsts.MaxAge != nil {
			hstsMap["max_age"] = domainConfig.Https.Hsts.MaxAge
		}
		if domainConfig.Https.Hsts.IncludeSubDomains != nil {
			hstsMap["include_sub_domains"] = domainConfig.Https.Hsts.IncludeSubDomains
		}
		httpsConfig["hsts"] = []interface{}{hstsMap}
	}
	httpsConfigs = append(httpsConfigs, httpsConfig)
	_ = d.Set("https_config", httpsConfigs)

	authRaw := d.Get("authentication").([]interface{})
	if authentication := domainConfig.Authentication; authentication != nil && len(authRaw) > 0 {
		auth := make(map[string]interface{})
		auth["switch"] = authentication.Switch
		if authType := authentication.TypeA; authType != nil {
			dMap := map[string]interface{}{
				"secret_key":        authType.SecretKey,
				"sign_param":        authType.SignParam,
				"expire_time":       authType.ExpireTime,
				"file_extensions":   authType.FileExtensions,
				"filter_type":       authType.FilterType,
				"backup_secret_key": authType.BackupSecretKey,
			}
			auth["type_a"] = []interface{}{dMap}
		}
		if authType := authentication.TypeB; authType != nil {
			dMap := map[string]interface{}{
				"secret_key":        authType.SecretKey,
				"expire_time":       authType.ExpireTime,
				"file_extensions":   authType.FileExtensions,
				"filter_type":       authType.FilterType,
				"backup_secret_key": authType.BackupSecretKey,
			}
			auth["type_b"] = []interface{}{dMap}
		}
		if authType := authentication.TypeC; authType != nil {
			dMap := map[string]interface{}{
				"secret_key":        authType.SecretKey,
				"expire_time":       authType.ExpireTime,
				"file_extensions":   authType.FileExtensions,
				"filter_type":       authType.FilterType,
				"time_format":       authType.TimeFormat,
				"backup_secret_key": authType.BackupSecretKey,
			}
			auth["type_c"] = []interface{}{dMap}
		}
		if authType := authentication.TypeD; authType != nil {
			dMap := map[string]interface{}{
				"secret_key":        authType.SecretKey,
				"expire_time":       authType.ExpireTime,
				"file_extensions":   authType.FileExtensions,
				"filter_type":       authType.FilterType,
				"time_param":        authType.TimeParam,
				"time_format":       authType.TimeFormat,
				"backup_secret_key": authType.BackupSecretKey,
			}
			auth["type_d"] = []interface{}{dMap}
		}
		_ = d.Set("authentication", []interface{}{auth})
	}

	// access_port
	if domainConfig.AccessPort != nil {
		portList := make([]interface{}, 0, len(domainConfig.AccessPort))
		for _, port := range domainConfig.AccessPort {
			if port != nil {
				portList = append(portList, int(*port))
			}
		}
		_ = d.Set("access_port", portList)
	}

	// auto_guard
	if domainConfig.AutoGuard != nil {
		autoGuard := map[string]interface{}{
			"switch": helper.PString(domainConfig.AutoGuard.Switch),
		}
		if domainConfig.AutoGuard.FilterRules != nil {
			filterRules := make([]interface{}, 0, len(domainConfig.AutoGuard.FilterRules))
			for _, rule := range domainConfig.AutoGuard.FilterRules {
				ruleMap := map[string]interface{}{
					"filter_type": helper.PString(rule.FilterType),
					"rule_type":   helper.PString(rule.RuleType),
					"rule_paths":  helper.StringsInterfaces(rule.RulePaths),
				}
				filterRules = append(filterRules, ruleMap)
			}
			autoGuard["filter_rules"] = filterRules
		}
		_ = d.Set("auto_guard", []interface{}{autoGuard})
	}

	// geo_blocker
	if domainConfig.GeoBlocker != nil {
		geoBlocker := map[string]interface{}{
			"switch": helper.PString(domainConfig.GeoBlocker.Switch),
		}
		if domainConfig.GeoBlocker.BlockRules != nil {
			blockRules := make([]interface{}, 0, len(domainConfig.GeoBlocker.BlockRules))
			for _, rule := range domainConfig.GeoBlocker.BlockRules {
				ruleMap := map[string]interface{}{
					"block_type": helper.PString(rule.BlockType),
					"rule_paths": helper.StringsInterfaces(rule.RulePaths),
					"rule_type":  helper.PString(rule.RuleType),
					"districts":  helper.StringsInterfaces(rule.Districts),
				}
				blockRules = append(blockRules, ruleMap)
			}
			geoBlocker["block_rules"] = blockRules
		}
		_ = d.Set("geo_blocker", []interface{}{geoBlocker})
	}

	dc := domainConfig

	if ok := checkCdnInfoWritable(d, "ip_filter", dc.IpFilter); ok {
		dMap := map[string]interface{}{
			"switch":      dc.IpFilter.Switch,
			"filter_type": dc.IpFilter.FilterType,
			"filters":     dc.IpFilter.Filters,
			"return_code": dc.IpFilter.ReturnCode,
		}
		if rules := dc.IpFilter.FilterRules; len(rules) > 0 {
			list := make([]map[string]interface{}, 0)
			for i := range rules {
				item := rules[i]
				rule := map[string]interface{}{
					"filter_type": item.FilterType,
					"filters":     item.Filters,
					"rule_type":   item.RuleType,
					"rule_paths":  item.RulePaths,
				}
				list = append(list, rule)
			}
			dMap["filter_rules"] = list
		}
		_ = helper.SetMapInterfaces(d, "ip_filter", dMap)
	}
	if ok := checkCdnInfoWritable(d, "ip_freq_limit", dc.IpFreqLimit); ok {
		dMap := map[string]interface{}{
			"switch": dc.IpFreqLimit.Switch,
			"qps":    dc.IpFreqLimit.Qps,
		}
		_ = helper.SetMapInterfaces(d, "ip_freq_limit", dMap)
	}
	if dc.StatusCodeCache != nil {
		dMap := map[string]interface{}{
			"switch": dc.StatusCodeCache.Switch,
		}
		if rules := dc.StatusCodeCache.CacheRules; len(rules) > 0 {
			list := make([]map[string]interface{}, 0)
			for i := range rules {
				item := rules[i]
				rule := map[string]interface{}{
					"status_code": item.StatusCode,
					"cache_time":  item.CacheTime,
				}
				list = append(list, rule)
			}
			dMap["cache_rules"] = list
		}
		_ = helper.SetMapInterfaces(d, "status_code_cache", dMap)
	}
	if ok := checkCdnInfoWritable(d, "compression", dc.Compression); ok {
		dMap := map[string]interface{}{
			"switch": dc.Compression.Switch,
		}
		if rules := dc.Compression.CompressionRules; len(rules) > 0 {
			list := make([]map[string]interface{}, 0)
			for i := range rules {
				item := rules[i]
				rule := map[string]interface{}{
					"compress":        item.Compress,
					"min_length":      item.MinLength,
					"max_length":      item.MaxLength,
					"algorithms":      item.Algorithms,
					"file_extensions": item.FileExtensions,
					"rule_type":       item.RuleType,
					"rule_paths":      item.RulePaths,
				}
				list = append(list, rule)
			}
			dMap["compression_rules"] = list
		}
		_ = helper.SetMapInterfaces(d, "compression", dMap)
	}
	if ok := checkCdnInfoWritable(d, "band_width_alert", dc.BandwidthAlert); ok {
		dMap := map[string]interface{}{
			"switch":                     dc.BandwidthAlert.Switch,
			"bps_threshold":              dc.BandwidthAlert.BpsThreshold,
			"counter_measure":            dc.BandwidthAlert.CounterMeasure,
			"last_trigger_time":          dc.BandwidthAlert.LastTriggerTime,
			"alert_switch":               dc.BandwidthAlert.AlertSwitch,
			"alert_percentage":           dc.BandwidthAlert.AlertPercentage,
			"last_trigger_time_overseas": dc.BandwidthAlert.LastTriggerTimeOverseas,
			"metric":                     dc.BandwidthAlert.Metric,
		}
		if si := dc.BandwidthAlert.StatisticItems; len(si) > 0 {
			rules := make([]interface{}, 0)
			for i := range si {
				item := si[i]
				rule := map[string]interface{}{
					"switch":           item.Switch,
					"type":             item.Type,
					"unblock_time":     item.UnBlockTime,
					"bps_threshold":    item.BpsThreshold,
					"counter_measure":  item.CounterMeasure,
					"alert_switch":     item.AlertSwitch,
					"alert_percentage": item.AlertPercentage,
					"metric":           item.Metric,
					"cycle":            item.Cycle,
				}
				rules = append(rules, rule)
			}
			dMap["statistic_item"] = rules
		}
		_ = helper.SetMapInterfaces(d, "band_width_alert", dMap)
	}
	if ok := checkCdnInfoWritable(d, "error_page", dc.ErrorPage); ok {
		dMap := map[string]interface{}{
			"switch": dc.ErrorPage.Switch,
		}
		if rules := dc.ErrorPage.PageRules; len(rules) > 0 {
			list := make([]map[string]interface{}, 0)
			for i := range rules {
				item := rules[i]
				rule := map[string]interface{}{
					"status_code":   item.StatusCode,
					"redirect_code": item.RedirectCode,
					"redirect_url":  item.RedirectUrl,
				}
				list = append(list, rule)
			}
			dMap["page_rules"] = list
		}
		_ = helper.SetMapInterfaces(d, "error_page", dMap)
	}
	if dc.ResponseHeader != nil {
		dMap := map[string]interface{}{
			"switch": dc.ResponseHeader.Switch,
		}
		if rules := dc.ResponseHeader.HeaderRules; len(rules) > 0 {
			list := make([]map[string]interface{}, 0)
			for i := range rules {
				item := rules[i]
				rule := map[string]interface{}{
					"header_mode":  item.HeaderMode,
					"header_name":  item.HeaderName,
					"header_value": item.HeaderValue,
					"rule_type":    item.RuleType,
					"rule_paths":   item.RulePaths,
				}
				list = append(list, rule)
			}
			dMap["header_rules"] = list
		}
		_ = helper.SetMapInterfaces(d, "response_header", dMap)
	}
	if ok := checkCdnInfoWritable(d, "downstream_capping", dc.DownstreamCapping); ok {
		dMap := map[string]interface{}{
			"switch": dc.DownstreamCapping.Switch,
		}
		if rules := dc.DownstreamCapping.CappingRules; len(rules) > 0 {
			list := make([]map[string]interface{}, 0)
			for i := range rules {
				item := rules[i]
				rule := map[string]interface{}{
					"rule_type":      item.RuleType,
					"rule_paths":     item.RulePaths,
					"kbps_threshold": item.KBpsThreshold,
				}
				list = append(list, rule)
			}
			dMap["capping_rules"] = list
		}
		_ = helper.SetMapInterfaces(d, "downstream_capping", dMap)
	}
	if dc.ResponseHeaderCache != nil {
		_ = d.Set("response_header_cache_switch", dc.ResponseHeaderCache.Switch)
	}
	if ok := checkCdnInfoWritable(d, "origin_pull_optimization", dc.OriginPullOptimization); ok {
		dMap := map[string]interface{}{
			"switch":            dc.OriginPullOptimization.Switch,
			"optimization_type": dc.OriginPullOptimization.OptimizationType,
		}
		_ = helper.SetMapInterfaces(d, "origin_pull_optimization", dMap)
	}
	if dc.Seo != nil {
		_ = d.Set("seo_switch", dc.Seo.Switch)
	}
	if ok := checkCdnInfoWritable(d, "referer", dc.Referer); ok {
		dMap := map[string]interface{}{
			"switch": dc.Referer.Switch,
		}
		if rules := dc.Referer.RefererRules; len(rules) > 0 {
			list := make([]map[string]interface{}, 0)
			for i := range rules {
				item := rules[i]
				rule := map[string]interface{}{
					"rule_type":    item.RuleType,
					"rule_paths":   item.RulePaths,
					"referer_type": item.RefererType,
					"referers":     item.Referers,
					"allow_empty":  item.AllowEmpty,
				}
				list = append(list, rule)
			}
			dMap["referer_rules"] = list
		}
		_ = helper.SetMapInterfaces(d, "referer", dMap)
	}
	if dc.VideoSeek != nil {
		_ = d.Set("video_seek_switch", dc.VideoSeek.Switch)
	}
	if ok := checkCdnInfoWritable(d, "max_age", dc.MaxAge); ok {
		dMap := map[string]interface{}{
			"switch": dc.MaxAge.Switch,
		}
		if rules := dc.MaxAge.MaxAgeRules; len(rules) > 0 {
			list := make([]map[string]interface{}, 0)
			for i := range rules {
				item := rules[i]
				rule := map[string]interface{}{
					"follow_origin":    item.FollowOrigin,
					"max_age_contents": item.MaxAgeContents,
					"max_age_type":     item.MaxAgeType,
					"max_age_time":     item.MaxAgeTime,
				}
				list = append(list, rule)
			}
			dMap["max_age_rules"] = list
		}
		_ = helper.SetMapInterfaces(d, "max_age", dMap)
	}
	if ok := checkCdnInfoWritable(d, "specific_config_mainland", dc.SpecificConfig); ok {
		specConfig, err := json.Marshal(dc.SpecificConfig.Mainland)
		if err == nil {
			_ = d.Set("specific_config_mainland", string(specConfig))
		}
	}
	if ok := checkCdnInfoWritable(d, "specific_config_overseas", dc.SpecificConfig); ok {
		specConfig, err := json.Marshal(dc.SpecificConfig.Overseas)
		if err == nil {
			_ = d.Set("specific_config_overseas", string(specConfig))
		}
	}
	if ok := checkCdnInfoWritable(d, "origin_pull_timeout", dc.OriginPullTimeout); ok {
		_ = helper.SetMapInterfaces(d, "origin_pull_timeout", map[string]interface{}{
			"connect_timeout": dc.OriginPullTimeout.ConnectTimeout,
			"receive_timeout": dc.OriginPullTimeout.ReceiveTimeout,
		})
	}
	if ok := checkCdnInfoWritable(d, "post_max_size", dc.PostMaxSize); ok {
		dMap := map[string]interface{}{
			"switch":   dc.PostMaxSize.Switch,
			"max_size": *dc.PostMaxSize.MaxSize / 1024 / 1024,
		}
		_ = helper.SetMapInterfaces(d, "post_max_size", dMap)
	}
	if ok := checkCdnInfoWritable(d, "cache_key", dc.CacheKey); ok {
		dMap := map[string]interface{}{
			"full_url_cache": dc.CacheKey.FullUrlCache,
			"ignore_case":    dc.CacheKey.IgnoreCase,
		}
		if qs := dc.CacheKey.QueryString; qs != nil {
			dMap["query_string"] = []interface{}{
				map[string]interface{}{
					"switch":  qs.Switch,
					"action":  qs.Action,
					"value":   qs.Value,
					"reorder": qs.Reorder,
				},
			}
		}
		if krs := dc.CacheKey.KeyRules; len(krs) > 0 {
			dMaps := make([]interface{}, 0)
			for i := range krs {
				kr := krs[i]
				dMap := map[string]interface{}{
					"full_url_cache": kr.FullUrlCache,
					"ignore_case":    kr.IgnoreCase,
				}
				if kr.RuleType != nil {
					dMap["rule_type"] = kr.RuleType
				}
				if len(kr.RulePaths) > 0 {
					dMap["rule_paths"] = helper.StringsInterfaces(kr.RulePaths)
				}
				if krqs := kr.QueryString; krqs != nil {
					dMap["query_string"] = []interface{}{
						map[string]interface{}{
							"value":  krqs.Value,
							"switch": krqs.Switch,
							"action": krqs.Action,
						},
					}
				}
				dMaps = append(dMaps, dMap)
			}
			dMap["key_rules"] = dMaps
		}
		_ = helper.SetMapInterfaces(d, "cache_key", dMap)
	} else if dc.CacheKey != nil && dc.CacheKey.FullUrlCache != nil {
		fullUrlCache := *dc.CacheKey.FullUrlCache == CDN_SWITCH_ON
		_ = d.Set("full_url_cache", fullUrlCache)
	}
	if dc.OfflineCache != nil {
		_ = d.Set("offline_cache_switch", dc.OfflineCache.Switch)
	}
	if dc.Quic != nil {
		_ = d.Set("quic_switch", dc.Quic.Switch)
	}
	if ok := checkCdnInfoWritable(d, "aws_private_access", dc.AwsPrivateAccess); ok {
		_ = helper.SetMapInterfaces(d, "aws_private_access", map[string]interface{}{
			"switch":     dc.AwsPrivateAccess.Switch,
			"access_key": dc.AwsPrivateAccess.AccessKey,
			"secret_key": dc.AwsPrivateAccess.SecretKey,
			"bucket":     dc.AwsPrivateAccess.Bucket,
			"region":     dc.AwsPrivateAccess.Region,
		})
	}
	if ok := checkCdnInfoWritable(d, "oss_private_access", dc.OssPrivateAccess); ok {
		_ = helper.SetMapInterfaces(d, "oss_private_access", map[string]interface{}{
			"switch":     dc.OssPrivateAccess.Switch,
			"access_key": dc.OssPrivateAccess.AccessKey,
			"secret_key": dc.OssPrivateAccess.SecretKey,
			"bucket":     dc.OssPrivateAccess.Bucket,
			"region":     dc.OssPrivateAccess.Region,
		})
	}
	if ok := checkCdnInfoWritable(d, "hw_private_access", dc.HwPrivateAccess); ok {
		_ = helper.SetMapInterfaces(d, "hw_private_access", map[string]interface{}{
			"switch":     dc.HwPrivateAccess.Switch,
			"access_key": dc.HwPrivateAccess.AccessKey,
			"secret_key": dc.HwPrivateAccess.SecretKey,
			"bucket":     dc.HwPrivateAccess.Bucket,
		})
	}
	if ok := checkCdnInfoWritable(d, "qn_private_access", dc.QnPrivateAccess); ok {
		_ = helper.SetMapInterfaces(d, "qn_private_access", map[string]interface{}{
			"switch":     dc.QnPrivateAccess.Switch,
			"access_key": dc.QnPrivateAccess.AccessKey,
			"secret_key": dc.QnPrivateAccess.SecretKey,
		})
	}
	if ok := checkCdnInfoWritable(d, "others_private_access", dc.OthersPrivateAccess); ok {
		_ = helper.SetMapInterfaces(d, "others_private_access", map[string]interface{}{
			"switch":     dc.OthersPrivateAccess.Switch,
			"access_key": dc.OthersPrivateAccess.AccessKey,
			"secret_key": dc.OthersPrivateAccess.SecretKey,
			"bucket":     dc.OthersPrivateAccess.Bucket,
			"region":     dc.OthersPrivateAccess.Region,
		})
	}
	if dc.HttpsBilling != nil {
		if dc.HttpsBilling.Switch != nil {
			tmpMap := map[string]interface{}{}
			tmpMap["switch"] = dc.HttpsBilling.Switch
			_ = d.Set("https_billing", []interface{}{tmpMap})
		}
	}

	// user_agent_filter
	if ok := checkCdnInfoWritable(d, "user_agent_filter", dc.UserAgentFilter); ok {
		dMap := map[string]interface{}{
			"switch": dc.UserAgentFilter.Switch,
		}
		if rules := dc.UserAgentFilter.FilterRules; len(rules) > 0 {
			list := make([]map[string]interface{}, 0, len(rules))
			for _, item := range rules {
				rule := map[string]interface{}{
					"rule_type":   item.RuleType,
					"rule_paths":  item.RulePaths,
					"user_agents": item.UserAgents,
					"filter_type": item.FilterType,
				}
				list = append(list, rule)
			}
			dMap["filter_rules"] = list
		}
		_ = helper.SetMapInterfaces(d, "user_agent_filter", dMap)
	}
	// url_redirect
	if ok := checkCdnInfoWritable(d, "url_redirect", dc.UrlRedirect); ok {
		dMap := map[string]interface{}{
			"switch": dc.UrlRedirect.Switch,
		}
		if rules := dc.UrlRedirect.PathRules; len(rules) > 0 {
			list := make([]map[string]interface{}, 0, len(rules))
			for _, item := range rules {
				rule := map[string]interface{}{
					"redirect_status_code": item.RedirectStatusCode,
					"pattern":              item.Pattern,
					"redirect_url":         item.RedirectUrl,
					"redirect_host":        item.RedirectHost,
					"full_match":           item.FullMatch,
				}
				list = append(list, rule)
			}
			dMap["path_rules"] = list
		}
		_ = helper.SetMapInterfaces(d, "url_redirect", dMap)
	}
	// origin_combine
	if ok := checkCdnInfoWritable(d, "origin_combine", dc.OriginCombine); ok {
		_ = helper.SetMapInterfaces(d, "origin_combine", map[string]interface{}{
			"switch": dc.OriginCombine.Switch,
		})
	}
	// range_origin_pull (per-path rules)
	if ok := checkCdnInfoWritable(d, "range_origin_pull", dc.RangeOriginPull); ok {
		dMap := map[string]interface{}{
			"switch": dc.RangeOriginPull.Switch,
		}
		if rules := dc.RangeOriginPull.RangeRules; len(rules) > 0 {
			list := make([]map[string]interface{}, 0, len(rules))
			for _, item := range rules {
				rule := map[string]interface{}{
					"switch":     item.Switch,
					"rule_type":  item.RuleType,
					"rule_paths": item.RulePaths,
				}
				list = append(list, rule)
			}
			dMap["range_rules"] = list
		}
		_ = helper.SetMapInterfaces(d, "range_origin_pull", dMap)
	}

	tags, errRet := tagService.DescribeResourceTags(ctx, CDN_SERVICE_NAME, CDN_RESOURCE_NAME_DOMAIN, region, domain)
	if errRet != nil {
		return errRet
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudCdnDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdn_domain.update")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	cdnService := CdnService{client: client}

	d.Partial(true)
	updateAttrs := make([]string, 0)

	domain := d.Id()
	request := cdn.NewUpdateDomainConfigRequest()
	request.Domain = &domain

	if d.HasChange("service_type") {
		request.ServiceType = helper.String(d.Get("service_type").(string))
		updateAttrs = append(updateAttrs, "service_type")
	}
	if d.HasChange("project_id") {
		request.ProjectId = helper.IntInt64(d.Get("project_id").(int))
		updateAttrs = append(updateAttrs, "project_id")
	}
	if d.HasChange("area") {
		request.Area = helper.String(d.Get("area").(string))
		updateAttrs = append(updateAttrs, "area")
	}
	if d.HasChange("range_origin_switch") {
		request.RangeOriginPull = &cdn.RangeOriginPull{}
		request.RangeOriginPull.Switch = helper.String(d.Get("range_origin_switch").(string))
		updateAttrs = append(updateAttrs, "range_origin_switch")
	}
	if d.HasChange("ipv6_access_switch") {
		request.Ipv6Access = &cdn.Ipv6Access{}
		request.Ipv6Access.Switch = helper.String(d.Get("ipv6_access_switch").(string))
		updateAttrs = append(updateAttrs, "ipv6_access_switch")
	}
	if d.HasChange("follow_redirect_switch") {
		request.FollowRedirect = &cdn.FollowRedirect{}
		request.FollowRedirect.Switch = helper.String(d.Get("follow_redirect_switch").(string))
		updateAttrs = append(updateAttrs, "follow_redirect_switch")
	}
	if d.HasChange("origin") {
		updateAttrs = append(updateAttrs, "origin")
		origins := d.Get("origin").([]interface{})
		if len(origins) < 1 {
			return fmt.Errorf("origin is required")
		}
		origin := origins[0].(map[string]interface{})
		request.Origin = &cdn.Origin{}
		request.Origin.OriginType = helper.String(origin["origin_type"].(string))
		originList := origin["origin_list"].([]interface{})
		request.Origin.Origins = make([]*string, 0, len(originList))
		for _, item := range originList {
			request.Origin.Origins = append(request.Origin.Origins, helper.String(item.(string)))
		}
		if v := origin["server_name"]; v.(string) != "" {
			request.Origin.ServerName = helper.String(v.(string))
		}
		if v := origin["cos_private_access"]; v.(string) != "" {
			request.Origin.CosPrivateAccess = helper.String(v.(string))
		}
		if v := origin["origin_pull_protocol"]; v.(string) != "" {
			request.Origin.OriginPullProtocol = helper.String(v.(string))
		}
		if v := origin["backup_origin_type"]; v.(string) != "" {
			request.Origin.BackupOriginType = helper.String(v.(string))
		}
		if v := origin["backup_server_name"]; v.(string) != "" {
			request.Origin.BackupServerName = helper.String(v.(string))
		}
		if v := origin["backup_origin_list"]; len(v.([]interface{})) > 0 {
			backupOriginList := v.([]interface{})
			request.Origin.BackupOrigins = make([]*string, 0, len(backupOriginList))
			for _, item := range backupOriginList {
				request.Origin.BackupOrigins = append(request.Origin.BackupOrigins, helper.String(item.(string)))
			}
		}
		if v := origin["origin_company"]; v.(string) != "" {
			request.Origin.OriginCompany = helper.String(v.(string))
		}
	}
	if d.HasChange("request_header") {
		updateAttrs = append(updateAttrs, "request_header")
		requestHeaders := d.Get("request_header").([]interface{})
		requestHeader := requestHeaders[0].(map[string]interface{})
		headerRule := requestHeader["header_rules"].([]interface{})
		var headerRules []*cdn.HttpHeaderPathRule
		for _, value := range headerRule {
			hr := &cdn.HttpHeaderPathRule{}
			headerRuleMap := value.(map[string]interface{})
			headerMode := headerRuleMap["header_mode"].(string)
			headerName := headerRuleMap["header_name"].(string)
			headerValue := headerRuleMap["header_value"].(string)
			ruleType := headerRuleMap["rule_type"].(string)
			rulePaths := headerRuleMap["rule_paths"].([]interface{})
			rulePathList := make([]*string, 0, len(rulePaths))
			for _, value := range rulePaths {
				rulePathList = append(rulePathList, helper.String(value.(string)))
			}
			hr.HeaderMode = &headerMode
			hr.HeaderName = &headerName
			hr.HeaderValue = &headerValue
			hr.RuleType = &ruleType
			hr.RulePaths = rulePathList
			headerRules = append(headerRules, hr)
		}
		request.RequestHeader = &cdn.RequestHeader{}
		request.RequestHeader.Switch = helper.String(requestHeader["switch"].(string))
		request.RequestHeader.HeaderRules = headerRules
	}
	if d.HasChange("rule_cache") {
		updateAttrs = append(updateAttrs, "rule_cache")
		ruleCache := d.Get("rule_cache").([]interface{})
		var ruleCaches []*cdn.RuleCache
		for _, v := range ruleCache {
			re := &cdn.RuleCache{}
			ruleCacheMap := v.(map[string]interface{})
			rulePaths := ruleCacheMap["rule_paths"].([]interface{})
			rulePathList := make([]*string, 0, len(rulePaths))
			ruleType := ruleCacheMap["rule_type"].(string)
			if ruleType == CDN_RULE_TYPE_DEFAULT {
				rulePathList = append(rulePathList, helper.String(CDN_RULE_PATH))
			} else {
				for _, value := range rulePaths {
					rulePathList = append(rulePathList, helper.String(value.(string)))
				}
			}
			switchFlag := ruleCacheMap["switch"].(string)
			cacheTime := ruleCacheMap["cache_time"].(int)
			compareMaxAge := ruleCacheMap["compare_max_age"].(string)
			ignoreCacheControl := ruleCacheMap["ignore_cache_control"].(string)
			ignoreSetCookie := ruleCacheMap["ignore_set_cookie"].(string)
			noCacheSwitch := ruleCacheMap["no_cache_switch"].(string)
			reValidate := ruleCacheMap["re_validate"].(string)
			followOriginSwitch := ruleCacheMap["follow_origin_switch"].(string)
			ruleCacheConfig := &cdn.RuleCacheConfig{}
			cache := &cdn.CacheConfigCache{}
			noCache := &cdn.CacheConfigNoCache{}
			followOrigin := &cdn.CacheConfigFollowOrigin{}
			ruleCacheConfig.Cache = cache
			ruleCacheConfig.NoCache = noCache
			ruleCacheConfig.FollowOrigin = followOrigin
			re.CacheConfig = ruleCacheConfig
			re.RulePaths = rulePathList
			re.RuleType = &ruleType
			re.CacheConfig.Cache.Switch = &switchFlag
			re.CacheConfig.Cache.CacheTime = helper.IntInt64(cacheTime)
			re.CacheConfig.Cache.CompareMaxAge = &compareMaxAge
			re.CacheConfig.Cache.IgnoreCacheControl = &ignoreCacheControl
			re.CacheConfig.Cache.IgnoreSetCookie = &ignoreSetCookie
			re.CacheConfig.NoCache.Switch = &noCacheSwitch
			re.CacheConfig.NoCache.Revalidate = &reValidate
			re.CacheConfig.FollowOrigin.Switch = &followOriginSwitch
			heuristicCacheSwitch := ruleCacheMap["heuristic_cache_switch"].(string)
			heuristicCacheTime := ruleCacheMap["heuristic_cache_time"].(int)
			if heuristicCacheSwitch != "" {
				re.CacheConfig.FollowOrigin.HeuristicCache = &cdn.HeuristicCache{
					Switch:      &heuristicCacheSwitch,
					CacheConfig: &cdn.CacheConfig{},
				}
				if heuristicCacheTime > 0 {
					re.CacheConfig.FollowOrigin.HeuristicCache.CacheConfig.HeuristicCacheTimeSwitch = helper.String("on")
					re.CacheConfig.FollowOrigin.HeuristicCache.CacheConfig.HeuristicCacheTime = helper.IntInt64(heuristicCacheTime)
				} else {
					re.CacheConfig.FollowOrigin.HeuristicCache.CacheConfig.HeuristicCacheTimeSwitch = helper.String("off")
				}
			}
			ruleCaches = append(ruleCaches, re)
		}
		request.Cache = &cdn.Cache{}
		request.Cache.RuleCache = ruleCaches
	}
	if d.HasChange("https_config") {
		updateAttrs = append(updateAttrs, "https_config")
		httpsConfigs := d.Get("https_config").([]interface{})
		if len(httpsConfigs) > 0 {
			config := httpsConfigs[0].(map[string]interface{})
			request.Https = &cdn.Https{}
			request.Https.Switch = helper.String(config["https_switch"].(string))
			if v := config["http2_switch"]; v.(string) != "" {
				request.Https.Http2 = helper.String(v.(string))
			}
			request.Https.OcspStapling = helper.String(config["ocsp_stapling_switch"].(string))
			request.Https.Spdy = helper.String(config["spdy_switch"].(string))
			request.Https.VerifyClient = helper.String(config["verify_client"].(string))
			if v := config["server_certificate_config"]; len(v.([]interface{})) > 0 {
				serverCerts := v.([]interface{})
				if len(serverCerts) > 0 && serverCerts[0] != nil {
					serverCert := serverCerts[0].(map[string]interface{})
					request.Https.CertInfo = &cdn.ServerCert{}
					if v := serverCert["certificate_id"]; v.(string) != "" {
						request.Https.CertInfo.CertId = helper.String(v.(string))
					}
					if v := serverCert["certificate_content"]; v.(string) != "" {
						request.Https.CertInfo.Certificate = helper.String(v.(string))
					}
					if v := serverCert["private_key"]; v.(string) != "" {
						request.Https.CertInfo.PrivateKey = helper.String(v.(string))
					}
					if v := serverCert["message"]; v.(string) != "" {
						request.Https.CertInfo.Message = helper.String(v.(string))
					}
				}
			}
			if v := config["client_certificate_config"]; len(v.([]interface{})) > 0 {
				clientCerts := v.([]interface{})
				if len(clientCerts) > 0 && clientCerts[0] != nil {
					clientCert := clientCerts[0].(map[string]interface{})
					request.Https.ClientCertInfo = &cdn.ClientCert{}
					if v := clientCert["certificate_content"]; v.(string) != "" {
						request.Https.ClientCertInfo.Certificate = helper.String(v.(string))
					}
				}
			}
			if v, ok := config["force_redirect"]; ok {
				forceRedirect := v.([]interface{})
				if len(forceRedirect) > 0 && forceRedirect[0] != nil {
					var redirect cdn.ForceRedirect
					redirectMap := forceRedirect[0].(map[string]interface{})
					if sw := redirectMap["switch"]; sw.(string) != "" {
						redirect.Switch = helper.String(sw.(string))
					}
					if rt := redirectMap["redirect_type"]; rt.(string) != "" {
						redirect.RedirectType = helper.String(rt.(string))
					}
					if rsc := redirectMap["redirect_status_code"]; rsc.(int) != 0 {
						redirect.RedirectStatusCode = helper.Int64(int64(rsc.(int)))
					}
					if ch := redirectMap["carry_headers"]; ch.(string) != "" {
						redirect.CarryHeaders = helper.String(ch.(string))
					}
					request.ForceRedirect = &redirect
				}
			}
			if v, ok := config["tls_versions"]; ok {
				request.Https.TlsVersion = helper.InterfacesStringsPoint(v.([]interface{}))
			}
			// HSTS
			if v, ok := config["hsts"]; ok {
				hstsList := v.([]interface{})
				if len(hstsList) > 0 && hstsList[0] != nil {
					hstsMap := hstsList[0].(map[string]interface{})
					request.Https.Hsts = &cdn.Hsts{
						Switch: helper.String(hstsMap["switch"].(string)),
					}
					if maxAge, ok := hstsMap["max_age"].(int); ok && maxAge > 0 {
						request.Https.Hsts.MaxAge = helper.IntInt64(maxAge)
					}
					if includeSubDomains, ok := hstsMap["include_sub_domains"].(string); ok && includeSubDomains != "" {
						request.Https.Hsts.IncludeSubDomains = &includeSubDomains
					}
				}
			}
		}
	}

	if d.HasChange("authentication") {
		updateAttrs = append(updateAttrs, "authentication")
		if v, ok := helper.InterfacesHeadMap(d, "authentication"); ok {
			switchOn := v["switch"].(string)
			request.Authentication = &cdn.Authentication{
				Switch: &switchOn,
			}

			if v, ok := v["type_a"].([]interface{}); ok && len(v) > 0 {
				var (
					item           = v[0].(map[string]interface{})
					secretKey      = item["secret_key"].(string)
					signParam      = item["sign_param"].(string)
					expireTime     = item["expire_time"].(int)
					fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
					filterType     = item["filter_type"].(string)
				)

				request.Authentication.TypeA = &cdn.AuthenticationTypeA{
					SecretKey:      &secretKey,
					SignParam:      &signParam,
					ExpireTime:     helper.IntInt64(expireTime),
					FileExtensions: fileExtensions,
					FilterType:     &filterType,
				}

				if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
					request.Authentication.TypeA.BackupSecretKey = &backupSecretKey
				}
			}

			if v, ok := v["type_b"].([]interface{}); ok && len(v) > 0 {
				var (
					item           = v[0].(map[string]interface{})
					secretKey      = item["secret_key"].(string)
					expireTime     = item["expire_time"].(int)
					fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
					filterType     = item["filter_type"].(string)
				)

				request.Authentication.TypeB = &cdn.AuthenticationTypeB{
					SecretKey:      &secretKey,
					ExpireTime:     helper.IntInt64(expireTime),
					FileExtensions: fileExtensions,
					FilterType:     &filterType,
				}

				if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
					request.Authentication.TypeB.BackupSecretKey = &backupSecretKey
				}
			}

			if v, ok := v["type_c"].([]interface{}); ok && len(v) > 0 {
				var (
					item           = v[0].(map[string]interface{})
					secretKey      = item["secret_key"].(string)
					expireTime     = item["expire_time"].(int)
					fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
					filterType     = item["filter_type"].(string)
				)

				request.Authentication.TypeC = &cdn.AuthenticationTypeC{
					SecretKey:      &secretKey,
					ExpireTime:     helper.IntInt64(expireTime),
					FileExtensions: fileExtensions,
					FilterType:     &filterType,
				}

				if timeFormat, ok := item["time_format"].(string); ok {
					request.Authentication.TypeC.TimeFormat = &timeFormat
				}

				if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
					request.Authentication.TypeC.BackupSecretKey = &backupSecretKey
				}
			}

			if v, ok := v["type_d"].([]interface{}); ok && len(v) > 0 {
				var (
					item           = v[0].(map[string]interface{})
					secretKey      = item["secret_key"].(string)
					expireTime     = item["expire_time"].(int)
					fileExtensions = helper.InterfacesStringsPoint(item["file_extensions"].([]interface{}))
					filterType     = item["filter_type"].(string)
					timeParam      = item["time_param"].(string)
				)

				request.Authentication.TypeD = &cdn.AuthenticationTypeD{
					SecretKey:      &secretKey,
					ExpireTime:     helper.IntInt64(expireTime),
					FileExtensions: fileExtensions,
					FilterType:     &filterType,
					TimeParam:      &timeParam,
				}

				if timeFormat, ok := item["time_format"].(string); ok {
					request.Authentication.TypeD.TimeFormat = &timeFormat
				}

				if backupSecretKey, ok := item["backup_secret_key"].(string); ok {
					request.Authentication.TypeD.BackupSecretKey = &backupSecretKey
				}
			}
		}
	}

	// access_port
	if d.HasChange("access_port") {
		updateAttrs = append(updateAttrs, "access_port")
		if v, ok := d.GetOk("access_port"); ok {
			ports := v.([]interface{})
			portList := make([]*int64, 0, len(ports))
			for _, port := range ports {
				portValue := int64(port.(int))
				portList = append(portList, &portValue)
			}
			request.AccessPort = portList
		}
	}

	// more added
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "ip_filter"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "ip_filter")
		vSwitch := v["switch"].(string)
		request.IpFilter = &cdn.IpFilter{
			Switch: &vSwitch,
		}
		if vv, ok := v["filter_type"].(string); ok {
			request.IpFilter.FilterType = &vv
		}
		if vv, ok := v["filters"].([]interface{}); ok {
			request.IpFilter.Filters = helper.InterfacesStringsPoint(vv)
		}

		//need white list func
		if vv, ok := v["filter_rules"].([]interface{}); ok && len(vv) > 0 {
			filterRules := make([]*cdn.IpFilterPathRule, 0)
			for i := range vv {
				item := vv[i].(map[string]interface{})
				rule := &cdn.IpFilterPathRule{}
				if rv, ok := item["filter_type"].(string); ok && rv != "" {
					rule.FilterType = &rv
				}
				if rv, ok := item["filters"].([]interface{}); ok && len(rv) > 0 {
					rule.Filters = helper.InterfacesStringsPoint(rv)
				}
				if rv, ok := item["rule_type"].(string); ok && rv != "" {
					rule.RuleType = &rv
				}
				if rv, ok := item["rule_paths"].([]interface{}); ok && len(rv) > 0 {
					rule.RulePaths = helper.InterfacesStringsPoint(rv)
				}
				filterRules = append(filterRules, rule)
			}
			request.IpFilter.FilterRules = filterRules
		}

		if vv, ok := v["return_code"].(int); ok {
			request.IpFilter.ReturnCode = helper.IntInt64(vv)
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "ip_freq_limit"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "ip_freq_limit")
		vSwitch := v["switch"].(string)
		qps := v["qps"].(int)
		request.IpFreqLimit = &cdn.IpFreqLimit{
			Switch: &vSwitch,
			Qps:    helper.IntInt64(qps),
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "status_code_cache"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "status_code_cache")
		vSwitch := v["switch"].(string)
		request.StatusCodeCache = &cdn.StatusCodeCache{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.StatusCodeCacheRule, 0)
		rules := v["cache_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.StatusCodeCacheRule{
				StatusCode: helper.String(item["status_code"].(string)),
			}
			if v, ok := item["cache_time"].(int); ok && v > 0 {
				rule.CacheTime = helper.IntInt64(v)
			}
			requestRules = append(requestRules, rule)
		}
		request.StatusCodeCache.CacheRules = requestRules
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "compression"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "compression")
		vSwitch := v["switch"].(string)
		request.Compression = &cdn.Compression{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.CompressionRule, 0)
		rules := v["compression_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			var (
				compress = item["compress"].(bool)
			)
			rule := &cdn.CompressionRule{
				Compress: &compress,
			}
			if v, ok := item["min_length"].(int); ok && v > 0 {
				rule.MinLength = helper.IntInt64(v)
			}
			if v, ok := item["max_length"].(int); ok && v > 0 {
				rule.MaxLength = helper.IntInt64(v)
			}
			if v, ok := item["algorithms"].([]interface{}); ok && len(v) > 0 {
				rule.Algorithms = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["file_extensions"].([]interface{}); ok && len(v) > 0 {
				rule.FileExtensions = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["rule_type"].(string); ok && v != "" {
				rule.RuleType = &v
			}
			if v, ok := item["rule_paths"].([]interface{}); ok && len(v) > 0 {
				rule.RulePaths = helper.InterfacesStringsPoint(v)
			}

			requestRules = append(requestRules, rule)
		}
		request.Compression.CompressionRules = requestRules

	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "band_width_alert"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "band_width_alert")
		vSwitch := v["switch"].(string)
		request.BandwidthAlert = &cdn.BandwidthAlert{
			Switch: &vSwitch,
		}
		if v, ok := v["bps_threshold"].(int); ok && v > 0 {
			request.BandwidthAlert.BpsThreshold = helper.IntInt64(v)
		}
		if v, ok := v["counter_measure"].(string); ok && v != "" {
			request.BandwidthAlert.CounterMeasure = &v
		}
		//if v, ok := v["last_trigger_time"].(string); ok && v != "" {
		//	request.BandwidthAlert.LastTriggerTime = &v
		//}
		if v, ok := v["alert_switch"].(string); ok && v != "" {
			request.BandwidthAlert.AlertSwitch = &v
		}
		if v, ok := v["alert_percentage"].(int); ok && v > 0 {
			request.BandwidthAlert.AlertPercentage = helper.IntInt64(v)
		}
		//if v, ok := v["last_trigger_time_overseas"].(string); ok && v != "" {
		//	request.BandwidthAlert.LastTriggerTimeOverseas = &v
		//}
		if v, ok := v["metric"].(string); ok && v != "" {
			request.BandwidthAlert.Metric = &v
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "error_page"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "error_page")
		vSwitch := v["switch"].(string)
		request.ErrorPage = &cdn.ErrorPage{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.ErrorPageRule, 0)
		rules := v["page_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.ErrorPageRule{}
			if v, ok := item["status_code"].(int); ok && v != 0 {
				rule.StatusCode = helper.IntInt64(v)
			}
			if v, ok := item["redirect_code"].(int); ok && v != 0 {
				rule.RedirectCode = helper.IntInt64(v)
			}
			if v, ok := item["redirect_url"].(string); ok && v != "" {
				rule.RedirectUrl = &v
			}
			requestRules = append(requestRules, rule)
		}
		request.ErrorPage.PageRules = requestRules
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "response_header"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "response_header")
		vSwitch := v["switch"].(string)
		request.ResponseHeader = &cdn.ResponseHeader{
			Switch: &vSwitch,
		}
		responseRules := make([]*cdn.HttpHeaderPathRule, 0)
		rules := v["header_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.HttpHeaderPathRule{}
			if v, ok := item["header_mode"].(string); ok && v != "" {
				rule.HeaderMode = &v
			}
			if v, ok := item["header_name"].(string); ok && v != "" {
				rule.HeaderName = &v
			}
			if v, ok := item["header_value"].(string); ok && v != "" {
				rule.HeaderValue = &v
			}
			if v, ok := item["rule_type"].(string); ok && v != "" {
				rule.RuleType = &v
			}
			if v, ok := item["rule_paths"].([]interface{}); ok && len(v) > 0 {
				rule.RulePaths = helper.InterfacesStringsPoint(v)
			}
			responseRules = append(responseRules, rule)
		}
		request.ResponseHeader.HeaderRules = responseRules
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "downstream_capping"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "downstream_capping")
		vSwitch := v["switch"].(string)
		request.DownstreamCapping = &cdn.DownstreamCapping{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.CappingRule, 0)
		rules := v["capping_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.CappingRule{}
			if v, ok := item["rule_type"].(string); ok && v != "" {
				rule.RuleType = &v
			}
			if v, ok := item["rule_paths"].([]interface{}); ok && len(v) > 0 {
				rule.RulePaths = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["kbps_threshold"].(int); ok && v > 0 {
				rule.KBpsThreshold = helper.IntInt64(v)
			}
			requestRules = append(requestRules, rule)
		}
		request.DownstreamCapping.CappingRules = requestRules
	}
	if d.HasChange("response_header_cache_switch") {
		updateAttrs = append(updateAttrs, "response_header_cache_switch")
		v, ok := d.Get("response_header_cache_switch").(string)
		if !ok || v == "" {
			v = "off"
		}
		request.ResponseHeaderCache = &cdn.ResponseHeaderCache{
			Switch: &v,
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "origin_pull_optimization"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "origin_pull_optimization")
		vSwitch := v["switch"].(string)
		request.OriginPullOptimization = &cdn.OriginPullOptimization{
			Switch: &vSwitch,
		}
		if v, ok := v["optimization_type"].(string); ok && v != "" {
			request.OriginPullOptimization.OptimizationType = &v
		}
	}
	if d.HasChange("seo_switch") {
		updateAttrs = append(updateAttrs, "seo_switch")
		v, ok := d.Get("seo_switch").(string)
		if !ok || v == "" {
			v = "off"
		}
		request.Seo = &cdn.Seo{
			Switch: &v,
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "referer"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "referer")
		vSwitch := v["switch"].(string)
		request.Referer = &cdn.Referer{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.RefererRule, 0)
		rules := v["referer_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.RefererRule{}
			if v, ok := item["rule_type"].(string); ok && v != "" {
				rule.RuleType = &v
			}
			if v, ok := item["rule_paths"].([]interface{}); ok && len(v) > 0 {
				rule.RulePaths = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["referer_type"].(string); ok && v != "" {
				rule.RefererType = &v
			}
			if v, ok := item["referers"].([]interface{}); ok && len(v) > 0 {
				rule.Referers = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["allow_empty"].(bool); ok {
				rule.AllowEmpty = &v
			}
			requestRules = append(requestRules, rule)
		}
		request.Referer.RefererRules = requestRules
	}
	if d.HasChange("video_seek_switch") {
		updateAttrs = append(updateAttrs, "video_seek_switch")
		v, ok := d.Get("video_seek_switch").(string)
		if !ok || v == "" {
			v = "off"
		}
		request.VideoSeek = &cdn.VideoSeek{
			Switch: &v,
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "max_age"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "max_age")
		vSwitch := v["switch"].(string)
		request.MaxAge = &cdn.MaxAge{
			Switch: &vSwitch,
		}
		requestRules := make([]*cdn.MaxAgeRule, 0)
		rules := v["max_age_rules"].([]interface{})
		for i := range rules {
			item := rules[i].(map[string]interface{})
			rule := &cdn.MaxAgeRule{}

			if v, ok := item["max_age_type"].(string); ok && v != "" {
				rule.MaxAgeType = &v
			}
			if v, ok := item["max_age_contents"].([]interface{}); ok && len(v) > 0 {
				rule.MaxAgeContents = helper.InterfacesStringsPoint(v)
			}
			if v, ok := item["max_age_time"].(int); ok && v > 0 {
				rule.MaxAgeTime = helper.IntInt64(v)
			}
			if v, ok := item["follow_origin"].(string); ok && v != "" {
				rule.FollowOrigin = &v
			}

			requestRules = append(requestRules, rule)
		}
		request.MaxAge.MaxAgeRules = requestRules
	}
	if v, ok := d.GetOk("specific_config_mainland"); ok && v.(string) != "" {
		updateAttrs = append(updateAttrs, "specific_config_mainland")
		request.SpecificConfig = &cdn.SpecificConfig{}
		configStruct := cdn.MainlandConfig{}
		err := json.Unmarshal([]byte(v.(string)), &configStruct)
		if err != nil {
			return fmt.Errorf("unmarshal specific_config_mainland fail: %s", err.Error())
		}
		request.SpecificConfig.Mainland = &configStruct
	}
	if v, ok := d.GetOk("specific_config_overseas"); ok && v.(string) != "" {
		updateAttrs = append(updateAttrs, "specific_config_overseas")
		if request.SpecificConfig == nil {
			request.SpecificConfig = &cdn.SpecificConfig{}
		}
		configStruct := cdn.OverseaConfig{}
		err := json.Unmarshal([]byte(v.(string)), &configStruct)
		if err != nil {
			return fmt.Errorf("unmarshal specific_config_overseas fail: %s", err.Error())
		}
		request.SpecificConfig.Overseas = &configStruct
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "origin_pull_timeout"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "origin_pull_timeout")
		request.OriginPullTimeout = &cdn.OriginPullTimeout{}
		if vv, ok := v["connect_timeout"].(int); ok && vv > 0 {
			request.OriginPullTimeout.ConnectTimeout = helper.IntUint64(vv)
		}
		if vv, ok := v["receive_timeout"].(int); ok && vv > 0 {
			request.OriginPullTimeout.ReceiveTimeout = helper.IntUint64(vv)
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "post_max_size"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "post_max_size")
		vSwitch := v["switch"].(string)
		maxSize := v["max_size"].(int)
		request.PostMaxSize = &cdn.PostSize{
			Switch: &vSwitch,
		}
		if maxSize > 0 {
			request.PostMaxSize.MaxSize = helper.IntInt64(maxSize * 1024 * 1024)
		}
	}
	if v, ok := d.GetOk("offline_cache_switch"); ok {
		request.OfflineCache = &cdn.OfflineCache{
			Switch: helper.String(v.(string)),
		}
	}
	if v, ok := d.GetOk("quic_switch"); ok {
		request.Quic = &cdn.Quic{
			Switch: helper.String(v.(string)),
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "cache_key"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "cache_key")
		request.CacheKey = &cdn.CacheKey{}
		if fuc := v["full_url_cache"].(string); fuc != "" {
			request.CacheKey.FullUrlCache = &fuc
		}
		if ic := v["ignore_case"].(string); ic != "" {
			request.CacheKey.IgnoreCase = &ic
		}
		if qs, ok := v["query_string"].([]interface{}); ok && len(qs) > 0 {
			if dMap, ok := qs[0].(map[string]interface{}); ok {
				qSwitch := dMap["switch"].(string)
				reorder := dMap["reorder"].(string)
				action := dMap["action"].(string)
				value := dMap["value"].(string)
				request.CacheKey.QueryString = &cdn.QueryStringKey{
					Switch:  &qSwitch,
					Reorder: &reorder,
					Action:  &action,
					Value:   &value,
				}
			}
		}
		if kr, ok := v["key_rules"].([]interface{}); ok {
			for i := range kr {
				rule, ok := kr[i].(map[string]interface{})
				if !ok {
					continue
				}
				ruleType := rule["rule_type"].(string)
				keyRule := &cdn.KeyRule{
					RuleType: &ruleType,
				}
				if vv := rule["full_url_cache"].(string); vv != "" {
					keyRule.FullUrlCache = &vv
				}
				if vv := rule["ignore_case"].(string); vv != "" {
					keyRule.IgnoreCase = &vv
				}
				if vv := rule["rule_tag"].(string); vv != "" {
					keyRule.RuleTag = &vv
				}
				if rp, ok := rule["rule_paths"].([]interface{}); ok {
					keyRule.RulePaths = helper.InterfacesStringsPoint(rp)
				}
				if qs, ok := rule["query_string"].([]interface{}); ok && len(qs) > 0 {
					if dMap, ok := qs[0].(map[string]interface{}); ok {
						vSwitch := dMap["switch"].(string)
						keyRule.QueryString = &cdn.RuleQueryString{
							Switch: &vSwitch,
						}
						if v := dMap["action"].(string); v != "" && vSwitch == "on" {
							keyRule.QueryString.Action = &v
						}
						if v := dMap["value"].(string); v != "" {
							keyRule.QueryString.Value = &v
						}
					}
				}
				request.CacheKey.KeyRules = append(request.CacheKey.KeyRules, keyRule)
			}
		}
	} else if d.HasChange("full_url_cache") {
		fullUrlCache := d.Get("full_url_cache").(bool)
		request.CacheKey = &cdn.CacheKey{}
		if fullUrlCache {
			request.CacheKey.FullUrlCache = helper.String(CDN_SWITCH_ON)
		} else {
			request.CacheKey.FullUrlCache = helper.String(CDN_SWITCH_OFF)
		}
		updateAttrs = append(updateAttrs, "full_url_cache")
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "aws_private_access"); ok && hasChanged {
		vSwitch := v["switch"].(string)
		request.AwsPrivateAccess = &cdn.AwsPrivateAccess{
			Switch: &vSwitch,
		}
		if v, ok := v["access_key"].(string); ok && v != "" {
			request.AwsPrivateAccess.AccessKey = &v
		}
		if v, ok := v["secret_key"].(string); ok && v != "" {
			request.AwsPrivateAccess.SecretKey = &v
		}
		if v, ok := v["region"].(string); ok && v != "" {
			request.AwsPrivateAccess.Region = &v
		}
		if v, ok := v["bucket"].(string); ok && v != "" {
			request.AwsPrivateAccess.Bucket = &v
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "oss_private_access"); ok && hasChanged {
		vSwitch := v["switch"].(string)
		request.OssPrivateAccess = &cdn.OssPrivateAccess{
			Switch: &vSwitch,
		}
		if v, ok := v["access_key"].(string); ok && v != "" {
			request.OssPrivateAccess.AccessKey = &v
		}
		if v, ok := v["secret_key"].(string); ok && v != "" {
			request.OssPrivateAccess.SecretKey = &v
		}
		if v, ok := v["region"].(string); ok && v != "" {
			request.OssPrivateAccess.Region = &v
		}
		if v, ok := v["bucket"].(string); ok && v != "" {
			request.OssPrivateAccess.Bucket = &v
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "hw_private_access"); ok && hasChanged {
		vSwitch := v["switch"].(string)
		request.HwPrivateAccess = &cdn.HwPrivateAccess{
			Switch: &vSwitch,
		}
		if v, ok := v["access_key"].(string); ok && v != "" {
			request.HwPrivateAccess.AccessKey = &v
		}
		if v, ok := v["secret_key"].(string); ok && v != "" {
			request.HwPrivateAccess.SecretKey = &v
		}
		if v, ok := v["bucket"].(string); ok && v != "" {
			request.HwPrivateAccess.Bucket = &v
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "qn_private_access"); ok && hasChanged {
		vSwitch := v["switch"].(string)
		request.QnPrivateAccess = &cdn.QnPrivateAccess{
			Switch: &vSwitch,
		}
		if v, ok := v["access_key"].(string); ok && v != "" {
			request.QnPrivateAccess.AccessKey = &v
		}
		if v, ok := v["secret_key"].(string); ok && v != "" {
			request.QnPrivateAccess.SecretKey = &v
		}
	}
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "others_private_access"); ok && hasChanged {
		vSwitch := v["switch"].(string)
		request.OthersPrivateAccess = &cdn.OthersPrivateAccess{
			Switch: &vSwitch,
		}
		if v, ok := v["access_key"].(string); ok && v != "" {
			request.OthersPrivateAccess.AccessKey = &v
		}
		if v, ok := v["secret_key"].(string); ok && v != "" {
			request.OthersPrivateAccess.SecretKey = &v
		}
		if v, ok := v["region"].(string); ok && v != "" {
			request.OthersPrivateAccess.Region = &v
		}
		if v, ok := v["bucket"].(string); ok && v != "" {
			request.OthersPrivateAccess.Bucket = &v
		}
	}
	if v, ok := helper.InterfacesHeadMap(d, "https_billing"); ok {
		updateAttrs = append(updateAttrs, "https_billing")
		vSwitch := v["switch"].(string)
		request.HttpsBilling = &cdn.HttpsBilling{
			Switch: &vSwitch,
		}
	}
	// user_agent_filter
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "user_agent_filter"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "user_agent_filter")
		vSwitch := v["switch"].(string)
		request.UserAgentFilter = &cdn.UserAgentFilter{
			Switch: &vSwitch,
		}
		if rules, ok := v["filter_rules"].([]interface{}); ok && len(rules) > 0 {
			filterRules := make([]*cdn.UserAgentFilterRule, 0, len(rules))
			for _, rule := range rules {
				ruleMap := rule.(map[string]interface{})
				filterRule := &cdn.UserAgentFilterRule{}
				if rv, ok := ruleMap["rule_type"].(string); ok && rv != "" {
					filterRule.RuleType = &rv
				}
				if rv, ok := ruleMap["rule_paths"].([]interface{}); ok && len(rv) > 0 {
					filterRule.RulePaths = helper.InterfacesStringsPoint(rv)
				}
				if rv, ok := ruleMap["user_agents"].([]interface{}); ok && len(rv) > 0 {
					filterRule.UserAgents = helper.InterfacesStringsPoint(rv)
				}
				if rv, ok := ruleMap["filter_type"].(string); ok && rv != "" {
					filterRule.FilterType = &rv
				}
				filterRules = append(filterRules, filterRule)
			}
			request.UserAgentFilter.FilterRules = filterRules
		}
	}
	// url_redirect
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "url_redirect"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "url_redirect")
		vSwitch := v["switch"].(string)
		request.UrlRedirect = &cdn.UrlRedirect{
			Switch: &vSwitch,
		}
		if rules, ok := v["path_rules"].([]interface{}); ok && len(rules) > 0 {
			pathRules := make([]*cdn.UrlRedirectRule, 0, len(rules))
			for _, rule := range rules {
				ruleMap := rule.(map[string]interface{})
				pathRule := &cdn.UrlRedirectRule{}
				if rv, ok := ruleMap["redirect_status_code"].(int); ok && rv > 0 {
					pathRule.RedirectStatusCode = helper.IntInt64(rv)
				}
				if rv, ok := ruleMap["pattern"].(string); ok && rv != "" {
					pathRule.Pattern = &rv
				}
				if rv, ok := ruleMap["redirect_url"].(string); ok && rv != "" {
					pathRule.RedirectUrl = &rv
				}
				if rv, ok := ruleMap["redirect_host"].(string); ok && rv != "" {
					pathRule.RedirectHost = &rv
				}
				if rv, ok := ruleMap["full_match"].(bool); ok {
					pathRule.FullMatch = &rv
				}
				pathRules = append(pathRules, pathRule)
			}
			request.UrlRedirect.PathRules = pathRules
		}
	}
	// origin_combine
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "origin_combine"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "origin_combine")
		vSwitch := v["switch"].(string)
		request.OriginCombine = &cdn.OriginCombine{
			Switch: &vSwitch,
		}
	}
	// range_origin_pull
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "range_origin_pull"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "range_origin_pull")
		vSwitch := v["switch"].(string)
		request.RangeOriginPull = &cdn.RangeOriginPull{
			Switch: &vSwitch,
		}
		if rules, ok := v["range_rules"].([]interface{}); ok && len(rules) > 0 {
			rangeRules := make([]*cdn.RangeOriginPullRule, 0, len(rules))
			for _, rule := range rules {
				ruleMap := rule.(map[string]interface{})
				rangeRule := &cdn.RangeOriginPullRule{
					Switch: helper.String(ruleMap["switch"].(string)),
				}
				if rv, ok := ruleMap["rule_type"].(string); ok && rv != "" {
					rangeRule.RuleType = &rv
				}
				if rv, ok := ruleMap["rule_paths"].([]interface{}); ok && len(rv) > 0 {
					rangeRule.RulePaths = helper.InterfacesStringsPoint(rv)
				}
				rangeRules = append(rangeRules, rangeRule)
			}
			request.RangeOriginPull.RangeRules = rangeRules
		}
	}
	// auto_guard
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "auto_guard"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "auto_guard")
		autoGuard := &cdn.AutoGuard{}
		if sw, ok := v["switch"].(string); ok && sw != "" {
			autoGuard.Switch = helper.String(sw)
		}
		if rules, ok := v["filter_rules"].([]interface{}); ok && len(rules) > 0 {
			filterRules := make([]*cdn.FilterRules, 0, len(rules))
			for _, r := range rules {
				ruleMap := r.(map[string]interface{})
				rule := &cdn.FilterRules{}
				if ft, ok := ruleMap["filter_type"].(string); ok && ft != "" {
					rule.FilterType = helper.String(ft)
				}
				if rt, ok := ruleMap["rule_type"].(string); ok && rt != "" {
					rule.RuleType = helper.String(rt)
				}
				if rp, ok := ruleMap["rule_paths"].([]interface{}); ok {
					rule.RulePaths = helper.InterfacesStringsPoint(rp)
				}
				filterRules = append(filterRules, rule)
			}
			autoGuard.FilterRules = filterRules
		}
		request.AutoGuard = autoGuard
	}
	// geo_blocker
	if v, ok, hasChanged := checkCdnHeadMapOkAndChanged(d, "geo_blocker"); ok && hasChanged {
		updateAttrs = append(updateAttrs, "geo_blocker")
		geoBlocker := &cdn.GeoBlocker{}
		if sw, ok := v["switch"].(string); ok && sw != "" {
			geoBlocker.Switch = helper.String(sw)
		}
		if rules, ok := v["block_rules"].([]interface{}); ok && len(rules) > 0 {
			blockRules := make([]*cdn.GeoBlockStrategy, 0, len(rules))
			for _, r := range rules {
				ruleMap := r.(map[string]interface{})
				rule := &cdn.GeoBlockStrategy{}
				if bt, ok := ruleMap["block_type"].(string); ok && bt != "" {
					rule.BlockType = helper.String(bt)
				}
				if rp, ok := ruleMap["rule_paths"].([]interface{}); ok {
					rule.RulePaths = helper.InterfacesStringsPoint(rp)
				}
				if rt, ok := ruleMap["rule_type"].(string); ok && rt != "" {
					rule.RuleType = helper.String(rt)
				}
				if ds, ok := ruleMap["districts"].([]interface{}); ok {
					rule.Districts = helper.InterfacesStringsPoint(ds)
				}
				blockRules = append(blockRules, rule)
			}
			geoBlocker.BlockRules = blockRules
		}
		request.GeoBlocker = geoBlocker
	}

	if v := d.Get("explicit_using_dry_run").(bool); v {
		_ = d.Set("dry_run_update_result", request.ToJsonString())
		return resourceTencentCloudCdnDomainRead(d, meta)
	}

	if len(updateAttrs) > 0 {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			_, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCdnClient().UpdateDomainConfig(request)
			if err != nil {
				if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
					if sdkErr.Code == CDN_DOMAIN_CONFIG_ERROR {
						return resource.NonRetryableError(err)
					}
				}
				return tccommon.RetryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}

		err = resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			domainConfig, err := cdnService.DescribeDomainsConfigByDomain(ctx, domain)
			if err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			if *domainConfig.Status == CDN_DOMAIN_STATUS_PROCESSING {
				return resource.RetryableError(fmt.Errorf("domain status is still processing, retry..."))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		tagService := svctag.NewTagService(client)
		region := client.Region
		resourceName := tccommon.BuildTagResourceName(CDN_SERVICE_NAME, CDN_RESOURCE_NAME_DOMAIN, region, domain)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}

	}

	d.Partial(false)

	return resourceTencentCloudCdnDomainRead(d, meta)
}

func resourceTencentCloudCdnDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdn_domain.delete")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	cdnService := CdnService{client: client}

	domain := d.Id()

	if v, ok := d.Get("explicit_using_dry_run").(bool); ok && v {
		d.SetId("")
		return nil
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(client)
		region := client.Region
		resourceName := tccommon.BuildTagResourceName(CDN_SERVICE_NAME, CDN_RESOURCE_NAME_DOMAIN, region, domain)
		deleteTags := make([]string, 0, len(tags))
		for key := range tags {
			deleteTags = append(deleteTags, key)
		}
		err := tagService.ModifyTags(ctx, resourceName, nil, deleteTags)
		if err != nil {
			return err
		}
	}

	var domainConfig *cdn.DetailDomain
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		domainConfig, errRet = cdnService.DescribeDomainsConfigByDomain(ctx, domain)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if domainConfig == nil {
		return nil
	}

	if *domainConfig.Status == CDN_DOMAIN_STATUS_ONLINE {
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			errRet = cdnService.StopDomain(ctx, domain)
			if errRet != nil {
				return tccommon.RetryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}

		err = resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
			domainConfig, err := cdnService.DescribeDomainsConfigByDomain(ctx, domain)
			if err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			if *domainConfig.Status == CDN_DOMAIN_STATUS_PROCESSING {
				return resource.RetryableError(fmt.Errorf("domain status is still processing, retry..."))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		errRet = cdnService.DeleteDomain(ctx, domain)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func updateCdnModifyOnlyParams(d *schema.ResourceData, meta interface{}, ctx context.Context) error {
	needUpdate := false

	domain := d.Id()
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := CdnService{client}
	request := cdn.NewUpdateDomainConfigRequest()
	request.Domain = &domain

	if v, ok := helper.InterfacesHeadMap(d, "post_max_size"); ok {
		needUpdate = true
		vSwitch := v["switch"].(string)
		maxSize := v["max_size"].(int)
		request.PostMaxSize = &cdn.PostSize{
			Switch: &vSwitch,
		}
		if maxSize > 0 {
			request.PostMaxSize.MaxSize = helper.IntInt64(maxSize * 1024 * 1024)
		}
	}

	// user_agent_filter - not supported by Create API, must be set via Update API
	if v, ok := helper.InterfacesHeadMap(d, "user_agent_filter"); ok {
		needUpdate = true
		vSwitch := v["switch"].(string)
		request.UserAgentFilter = &cdn.UserAgentFilter{
			Switch: &vSwitch,
		}
		if rules, ok := v["filter_rules"].([]interface{}); ok && len(rules) > 0 {
			filterRules := make([]*cdn.UserAgentFilterRule, 0, len(rules))
			for _, rule := range rules {
				ruleMap := rule.(map[string]interface{})
				filterRule := &cdn.UserAgentFilterRule{}
				if rv, ok := ruleMap["rule_type"].(string); ok && rv != "" {
					filterRule.RuleType = &rv
				}
				if rv, ok := ruleMap["rule_paths"].([]interface{}); ok && len(rv) > 0 {
					filterRule.RulePaths = helper.InterfacesStringsPoint(rv)
				}
				if rv, ok := ruleMap["user_agents"].([]interface{}); ok && len(rv) > 0 {
					filterRule.UserAgents = helper.InterfacesStringsPoint(rv)
				}
				if rv, ok := ruleMap["filter_type"].(string); ok && rv != "" {
					filterRule.FilterType = &rv
				}
				filterRules = append(filterRules, filterRule)
			}
			request.UserAgentFilter.FilterRules = filterRules
		}
	}

	// url_redirect - not supported by Create API, must be set via Update API
	if v, ok := helper.InterfacesHeadMap(d, "url_redirect"); ok {
		needUpdate = true
		vSwitch := v["switch"].(string)
		request.UrlRedirect = &cdn.UrlRedirect{
			Switch: &vSwitch,
		}
		if rules, ok := v["path_rules"].([]interface{}); ok && len(rules) > 0 {
			pathRules := make([]*cdn.UrlRedirectRule, 0, len(rules))
			for _, rule := range rules {
				ruleMap := rule.(map[string]interface{})
				pathRule := &cdn.UrlRedirectRule{}
				if rv, ok := ruleMap["redirect_status_code"].(int); ok && rv > 0 {
					pathRule.RedirectStatusCode = helper.IntInt64(rv)
				}
				if rv, ok := ruleMap["pattern"].(string); ok && rv != "" {
					pathRule.Pattern = &rv
				}
				if rv, ok := ruleMap["redirect_url"].(string); ok && rv != "" {
					pathRule.RedirectUrl = &rv
				}
				if rv, ok := ruleMap["redirect_host"].(string); ok && rv != "" {
					pathRule.RedirectHost = &rv
				}
				if rv, ok := ruleMap["full_match"].(bool); ok {
					pathRule.FullMatch = &rv
				}
				pathRules = append(pathRules, pathRule)
			}
			request.UrlRedirect.PathRules = pathRules
		}
	}

	// origin_combine - not supported by Create API, must be set via Update API
	if v, ok := helper.InterfacesHeadMap(d, "origin_combine"); ok {
		needUpdate = true
		vSwitch := v["switch"].(string)
		request.OriginCombine = &cdn.OriginCombine{
			Switch: &vSwitch,
		}
	}

	// access_port - not supported by Create API, must be set via Update API
	if v, ok := d.GetOk("access_port"); ok {
		needUpdate = true
		ports := v.([]interface{})
		portList := make([]*int64, 0, len(ports))
		for _, port := range ports {
			portValue := int64(port.(int))
			portList = append(portList, &portValue)
		}
		request.AccessPort = portList
	}

	if v, ok := helper.InterfacesHeadMap(d, "auto_guard"); ok {
		needUpdate = true
		autoGuard := &cdn.AutoGuard{}
		if sw, ok := v["switch"].(string); ok && sw != "" {
			autoGuard.Switch = helper.String(sw)
		}
		if rules, ok := v["filter_rules"].([]interface{}); ok && len(rules) > 0 {
			filterRules := make([]*cdn.FilterRules, 0, len(rules))
			for _, r := range rules {
				ruleMap := r.(map[string]interface{})
				rule := &cdn.FilterRules{}
				if ft, ok := ruleMap["filter_type"].(string); ok && ft != "" {
					rule.FilterType = helper.String(ft)
				}
				if rt, ok := ruleMap["rule_type"].(string); ok && rt != "" {
					rule.RuleType = helper.String(rt)
				}
				if rp, ok := ruleMap["rule_paths"].([]interface{}); ok {
					rule.RulePaths = helper.InterfacesStringsPoint(rp)
				}
				filterRules = append(filterRules, rule)
			}
			autoGuard.FilterRules = filterRules
		}
		request.AutoGuard = autoGuard
	}

	if v, ok := helper.InterfacesHeadMap(d, "geo_blocker"); ok {
		needUpdate = true
		geoBlocker := &cdn.GeoBlocker{}
		if sw, ok := v["switch"].(string); ok && sw != "" {
			geoBlocker.Switch = helper.String(sw)
		}
		if rules, ok := v["block_rules"].([]interface{}); ok && len(rules) > 0 {
			blockRules := make([]*cdn.GeoBlockStrategy, 0, len(rules))
			for _, r := range rules {
				ruleMap := r.(map[string]interface{})
				rule := &cdn.GeoBlockStrategy{}
				if bt, ok := ruleMap["block_type"].(string); ok && bt != "" {
					rule.BlockType = helper.String(bt)
				}
				if rp, ok := ruleMap["rule_paths"].([]interface{}); ok {
					rule.RulePaths = helper.InterfacesStringsPoint(rp)
				}
				if rt, ok := ruleMap["rule_type"].(string); ok && rt != "" {
					rule.RuleType = helper.String(rt)
				}
				if ds, ok := ruleMap["districts"].([]interface{}); ok {
					rule.Districts = helper.InterfacesStringsPoint(ds)
				}
				blockRules = append(blockRules, rule)
			}
			geoBlocker.BlockRules = blockRules
		}
		request.GeoBlocker = geoBlocker
	}

	if !needUpdate {
		return nil
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		err := service.UpdateDomainConfig(ctx, request)
		if err != nil {
			if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkErr.Code == CDN_DOMAIN_CONFIG_ERROR {
					return resource.NonRetryableError(err)
				}
			}
			return tccommon.RetryError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func checkCdnHeadMapOkAndChanged(d *schema.ResourceData, key string) (v map[string]interface{}, ok bool, changed bool) {
	changed = d.HasChange(key)
	v, ok = helper.InterfacesHeadMap(d, key)
	return
}

func checkCdnInfoWritable(d *schema.ResourceData, key string, val interface{}) bool {
	return val != nil
}
