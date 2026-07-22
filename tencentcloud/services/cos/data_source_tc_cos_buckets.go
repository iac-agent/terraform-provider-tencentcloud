package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCosBuckets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCosBucketsRead,

		Schema: map[string]*schema.Schema{
			"bucket_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A prefix 字符串 到 过滤器 results 通过 存储桶名称",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 到 过滤器 存储桶",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// computed
			"bucket_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 存储桶 Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "存储桶名称， 格式 likes `<存储桶>-<appid>`。",
						},
						"cors_rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A 列表 CORS 规则 configurations。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed_origins": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "指定which origins 是 allowed。",
									},
									"allowed_methods": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "指定which methods 是 allowed. Can 是 GET，PUT，POST，DELETE 或 HEAD。",
									},
									"allowed_headers": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "指定which headers 是 allowed。",
									},
									"max_age_seconds": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "指定time （秒） 该 browser 可以 缓存 response 对于 preflight 请求。",
									},
									"expose_headers": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "指定expose 头部 在 response。",
									},
								},
							},
						},
						"lifecycle_rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "lifecycle 配置 的 存储桶",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filter_prefix": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Object 键 prefix identifying 一个 或 more objects 到 其中 规则 applies。",
									},
									"transition": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "指定a 周期 在 对象's transitions。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"date": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "指定date after 其中 您 want corresponding 操作 到 take effect。",
												},
												"days": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "指定number 的 days after 对象 creation 当 特定 规则 操作 takes effect。",
												},
												"storage_class": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "指定storage class 到 其中 您 want 对象 到 transition. Available 值 include STANDARD，STANDARD_IA 和 ARCHIVE。",
												},
											},
										},
									},
									"expiration": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "指定a 周期 在 对象's expire。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"date": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "指定date after 其中 您 want corresponding 操作 到 take effect。",
												},
												"days": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "指定number 的 days after 对象 creation 当 特定 规则 操作 takes effect。",
												},
											},
										},
									},
									"non_current_transition": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "指定when 到 transition objects 的 non 当前 versions 和 目标 存储 class。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"non_current_days": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "数量 days after non 当前 对象 creation 当 特定 规则 操作 takes effect。",
												},
												"storage_class": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "指定storage class 到 其中 您 want non 当前 对象 到 transition. Available 值 include STANDARD，STANDARD_IA 和 ARCHIVE。",
												},
											},
										},
									},
									"non_current_expiration": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "指定when non 当前 对象 versions shall expire。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"non_current_days": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "数量 days after non 当前 对象 creation 当 特定 规则 操作 takes effect. 最大 值 是 3650。",
												},
											},
										},
									},
									"abort_incomplete_multipart_upload": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Set 最大 时间 multipart upload 是 allowed 到 remain running。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"days_after_initiation": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "指定number 的 days after multipart upload starts 该 upload 必须 是 completed. 最大 值 是 3650。",
												},
											},
										},
									},
								},
							},
						},
						"website": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A 列表 一个 element containing 配置 参数 使用 当 存储桶 是 使用 作为 website。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"index_document": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "COS 返回this 索引 document 当 requests 是 made 到 root 域名 或 any 的 subfolders。",
									},
									"error_document": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An absolute 路径 到 document 到 返回 在 case 的 4XX 错误",
									},
								},
							},
						},
						"origin_pull_rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "存储桶 Origin-Pull 规则。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"priority": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "优先级 的 源站-pull 规则，do 不 集合 same 值 对于 多个 规则。",
									},
									"sync_back_to_source": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "如果 `true`，COS 将 不 返回 3XX 状态 代码 当 pulling 数据 从 源站 服务器. Currently 可用 可用区: ap-beijing，ap-shanghai，ap-singapore，ap-mumbai。",
									},
									"back_to_source_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Back 到 来源 模式 Allow 值: Proxy，Mirror，Redirect。",
									},
									"prefix": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Triggers 源站-pull 规则 当 requested 文件 名称 matches 此 prefix。",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "协议 用于COS 到 访问 指定 源站 服务器. 可用 值 include `HTTP`，`HTTPS` 和 `FOLLOW`。",
									},
									"host": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Allows 仅 域名 名称 或 IP 地址 You 可以 optionally append 端口 数量 到 地址",
									},
									"follow_query_string": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "指定是否pass through COS 请求 查询 字符串 当 accessing 源站 服务器。",
									},
									"follow_redirection": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "指定是否follow 3XX redirect 到 another 源站 服务器 到 pull 数据 从。",
									},
									//"copy_origin_data": {
									//	Type:		 schema.TypeBool,
									//	Optional: 	 true,
									//	Default:	 true,
									//	Description: "",
									//},
									"follow_http_headers": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "指定pass through headers 当 accessing 源站 服务器。",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"custom_http_headers": {
										Type:        schema.TypeMap,
										Computed:    true,
										Description: "指定custom headers 该 您 可以 add 对于 COS 到 访问 your 源站 服务器。",
									},
									//"redirect_prefix": {
									//	Type:		schema.TypeString,
									//	Optional:   true,
									//	Description: "Prefix 对于 文件 到 其中 请求 是 redirected 当 源站-pull 规则 是 triggered.",
									//},
									//"redirect_suffix": {
									//	Type:		schema.TypeString,
									//	Optional:   true,
									//	Description: "Suffix 对于 文件 到 其中 请求 是 redirected 当 源站-pull 规则 是 triggered.",
									//},
								},
							},
						},
						"origin_domain_rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "存储桶 源站 域名 规则。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "指定domain 主机",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "指定origin 域名 类型，可用值：`REST`，`WEBSITE`，`ACCELERATE`，默认值：`REST`。",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "域名 状态，默认值：`ENABLED`。",
									},
								},
							},
						},
						"acl": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "存储桶 访问 control configurations。",
						},
						"acl_body": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "存储桶 verbose acl configurations。",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "标签 的 存储桶",
						},
						"cos_bucket_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL 的 此 COS 存储桶",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCosBucketsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cos_buckets.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cosService := CosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	buckets, err := cosService.ListBuckets(ctx)
	if err != nil {
		return err
	}

	prefix := d.Get("bucket_prefix").(string)
	tags := helper.GetTags(d, "tags")

	bucketList := make([]map[string]interface{}, 0, len(buckets))

	for _, v := range buckets {
		bucket := make(map[string]interface{})
		if prefix != "" && !strings.HasPrefix(*v.Name, prefix) {
			continue
		}

		respTags, err := cosService.GetBucketTags(ctx, *v.Name, "")
		if err != nil {
			return err
		}

		var matchTags bool

		for k, v := range tags {
			if respTags[k] == v {
				matchTags = true
				break
			}
		}

		if len(tags) != 0 && !matchTags {
			continue
		}

		bucket["bucket"] = *v.Name

		corsRules, err := cosService.GetBucketCors(ctx, *v.Name, "")
		if err != nil {
			return err
		}
		bucket["cors_rules"] = corsRules

		lifecycleRules, err := cosService.GetDataSourceBucketLifecycle(ctx, *v.Name)
		if err != nil {
			return err
		}
		bucket["lifecycle_rules"] = lifecycleRules

		website, err := cosService.GetBucketWebsite(ctx, *v.Name, "")
		if err != nil {
			return err
		}
		bucket["website"] = website

		cosDomain := meta.(tccommon.ProviderMeta).GetAPIV3Conn().CosDomain
		if cosDomain == "" {
			originRules, err := cosService.GetBucketPullOrigin(ctx, *v.Name)
			if err != nil {
				return err
			}
			bucket["origin_pull_rules"] = originRules

			domainRules, err := cosService.GetBucketOriginDomain(ctx, *v.Name)
			if err == nil {
				bucket["origin_domain_rules"] = domainRules
			}
		}

		aclBody, err := cosService.GetBucketACL(ctx, *v.Name, "")

		if err != nil {
			return err
		}

		aclXML, err := xml.Marshal(aclBody)

		if err != nil {
			log.Printf("WARN: acl body marshal failed: %s", err.Error())
		} else {
			bucket["acl"] = GetBucketPublicACL(aclBody)
			bucket["acl_body"] = string(aclXML)
		}

		bucket["tags"] = respTags
		bucket["cos_bucket_url"] = fmt.Sprintf("%s.cos.%s.myqcloud.com", *v.Name, meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region)
		bucketList = append(bucketList, bucket)
	}

	ids := make([]string, 2)
	ids[0] = "bucketlist"
	ids[1] = prefix
	d.SetId(helper.DataResourceIdsHash(ids))
	if err := d.Set("bucket_list", bucketList); err != nil {
		return fmt.Errorf("setting bucket list error: %s", err.Error())
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err = tccommon.WriteToFile(output.(string), bucketList); err != nil {
			return err
		}
	}

	return nil
}
