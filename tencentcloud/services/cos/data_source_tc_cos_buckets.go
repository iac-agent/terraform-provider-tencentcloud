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
				Description: "A prefix string to filter results by 存储桶名称",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 to filter 存储桶",
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
							Description: "存储桶名称，the 格式 likes `<存储桶>-<appid>`。",
						},
						"cors_rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A 列表 CORS rule configurations。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed_origins": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "指定which origins are allowed。",
									},
									"allowed_methods": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "指定which methods are allowed. Can be GET，PUT，POST，DELETE or HEAD。",
									},
									"allowed_headers": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "指定which headers are allowed。",
									},
									"max_age_seconds": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "指定time （秒） that browser can cache the response for a preflight request。",
									},
									"expose_headers": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "指定expose header in the response。",
									},
								},
							},
						},
						"lifecycle_rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The lifecycle configuration of a 存储桶",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filter_prefix": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Object 键 prefix identifying one or more objects to which the rule applies。",
									},
									"transition": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "指定a 周期 in the object's transitions。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"date": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "指定date after which you want the corresponding 操作 to take effect。",
												},
												"days": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "指定number of days after object creation when the specific rule 操作 takes effect。",
												},
												"storage_class": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "指定storage class to which you want the object to transition. Available values include STANDARD，STANDARD_IA and ARCHIVE。",
												},
											},
										},
									},
									"expiration": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "指定a 周期 in the object's expire。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"date": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "指定date after which you want the corresponding 操作 to take effect。",
												},
												"days": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "指定number of days after object creation when the specific rule 操作 takes effect。",
												},
											},
										},
									},
									"non_current_transition": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "指定when to transition objects of non current versions and the target storage class。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"non_current_days": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "数量 days after non current object creation when the specific rule 操作 takes effect。",
												},
												"storage_class": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "指定storage class to which you want the non current object to transition. Available values include STANDARD，STANDARD_IA and ARCHIVE。",
												},
											},
										},
									},
									"non_current_expiration": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "指定when non current object versions shall expire。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"non_current_days": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "数量 days after non current object creation when the specific rule 操作 takes effect. The maximum 值 is 3650。",
												},
											},
										},
									},
									"abort_incomplete_multipart_upload": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Set the maximum time a multipart upload is allowed to remain running。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"days_after_initiation": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "指定number of days after the multipart upload starts that the upload must be completed. The maximum 值 is 3650。",
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
							Description: "A 列表 one element containing configuration parameters used when the 存储桶 is used as a website。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"index_document": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "COS 返回this 索引 document when requests are made to the root 域名 or any of the subfolders。",
									},
									"error_document": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An absolute 路径 to the document to return in case of a 4XX 错误",
									},
								},
							},
						},
						"origin_pull_rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "存储桶 Origin-Pull rules。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"priority": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "优先级 of origin-pull rules，do not set the same 值 for multiple rules。",
									},
									"sync_back_to_source": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "If `true`，COS will not return 3XX 状态 代码 when pulling data from an origin server. Currently available 可用区: ap-beijing，ap-shanghai，ap-singapore，ap-mumbai。",
									},
									"back_to_source_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Back to 来源 模式 Allow 值: Proxy，Mirror，Redirect。",
									},
									"prefix": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Triggers the origin-pull rule when the requested file 名称 matches this prefix。",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "the 协议 用于COS to access the specified origin server. The available 值 include `HTTP`，`HTTPS` and `FOLLOW`。",
									},
									"host": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Allows only a 域名 名称 or IP 地址 You can optionally append a 端口 number to the 地址",
									},
									"follow_query_string": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "指定是否pass through COS request query string when accessing the origin server。",
									},
									"follow_redirection": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "指定是否follow 3XX redirect to another origin server to pull data from。",
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
										Description: "指定pass through headers when accessing the origin server。",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"custom_http_headers": {
										Type:        schema.TypeMap,
										Computed:    true,
										Description: "指定custom headers that you can add for COS to access your origin server。",
									},
									//"redirect_prefix": {
									//	Type:		schema.TypeString,
									//	Optional:   true,
									//	Description: "Prefix for the file to which a request is redirected when the origin-pull rule is triggered.",
									//},
									//"redirect_suffix": {
									//	Type:		schema.TypeString,
									//	Optional:   true,
									//	Description: "Suffix for the file to which a request is redirected when the origin-pull rule is triggered.",
									//},
								},
							},
						},
						"origin_domain_rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "存储桶 origin 域名 rules。",
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
							Description: "存储桶 access control configurations。",
						},
						"acl_body": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "存储桶 verbose acl configurations。",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The 标签 of a 存储桶",
						},
						"cos_bucket_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of this COS 存储桶",
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
