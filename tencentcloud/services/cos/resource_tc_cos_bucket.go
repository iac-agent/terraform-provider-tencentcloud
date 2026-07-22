package cos

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/beevik/etree"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func originPullRules() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"priority": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "优先级 的 源站-pull 规则，do 不 集合 same 值 对于 多个 规则。",
			},
			"sync_back_to_source": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Deprecated:  "It has been deprecated from version 1.81.196. Please use `back_to_source_mode` instead.",
				Description: "如果 `true`，COS 将 不 返回 3XX 状态 代码 当 pulling 数据 从 源站 服务器. Current 可用 可用区: ap-beijing，ap-shanghai，ap-singapore，ap-mumbai。",
			},
			"back_to_source_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"Proxy", "Mirror", "Redirect"}),
				Description:  "Back 到 来源 模式 Allow 值: Proxy，Mirror，Redirect。",
			},
			"http_redirect_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Redirect 代码 Effective 当 `back_to_source_mode` 是 `Redirect`. ex: 301，302，307. 默认为 302。",
			},
			"prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Triggers 源站-pull 规则 当 requested 文件 名称 matches 此 prefix。",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "协议 用于COS 到 访问 指定 源站 服务器. 可用 值 include `HTTP`，`HTTPS` 和 `FOLLOW`。",
			},
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Allows 仅 域名 名称 或 IP 地址 You 可以 optionally append 端口 数量 到 地址",
			},
			"follow_query_string": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "指定是否pass through COS 请求 查询 字符串 当 accessing 源站 服务器。",
			},
			"follow_redirection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "指定是否follow 3XX redirect 到 another 源站 服务器 到 pull 数据 从。",
			},
			//"copy_origin_data": {
			//	Type:		 schema.TypeBool,
			//	Optional: 	 true,
			//	Default:	 true,
			//	Description: "",
			//},
			"follow_http_headers": {
				Type:     schema.TypeSet,
				Optional: true,
				Set: func(i interface{}) int {
					return helper.HashString(i.(string))
				},
				Description: "指定pass through headers 当 accessing 源站 服务器。",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"custom_http_headers": {
				Type:        schema.TypeMap,
				Optional:    true,
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
	}
}

// x-cos-grant-* headers may conflict with xml acl body, we don't open up for now.
//func aclGrantHeaders() *schema.Schema {
//	return &schema.Schema{
//		Type:        schema.TypeMap,
//		Optional:    true,
//		Description: "ACL x-cos-grant-* headers 对于 多个 grand info",
//		Elem: &schema.Resource{
//			Schema: map[string]*schema.Schema{
//				"grant_read": {
//					Type:        schema.TypeString,
//					Optional:    true,
//					Description: "Allows grantee 到 read 存储桶; 格式: `ID=\"[OwnerUin]\"`.Use comma (,) 到 separate 多个 users, e.g `ID=\"100000000001\",ID=\"100000000002\"`",
//				},
//				"grant_write": {
//					Type:        schema.TypeString,
//					Optional:    true,
//					Description: "Allows grantee 到 write 到 存储桶; 格式: `ID=\"[OwnerUin]\"`.Use comma (,) 到 separate 多个 users, e.g `ID=\"100000000001\",ID=\"100000000002\"`",
//				},
//				"grant_read_acp": {
//					Type:        schema.TypeString,
//					Optional:    true,
//					Description: "Allows grantee 到 read ACL 的 存储桶; 格式: `ID=\"[OwnerUin]\"`.Use comma (,) 到 separate 多个 users, e.g `ID=\"100000000001\",ID=\"100000000002\"`",
//				},
//				"grant_write_acp": {
//					Type:        schema.TypeString,
//					Optional:    true,
//					Description: "Allows grantee 到 write ACL 的 存储桶; 格式: `ID=\"[OwnerUin]\"`.Use comma (,) 到 separate 多个 users, e.g `ID=\"100000000001\",ID=\"100000000002\"`",
//				},
//				"grant_full_control": {
//					Type:        schema.TypeString,
//					Optional:    true,
//					Description: "Grants 用户 full 权限 到 perform operations 在 存储桶; 格式: `ID=\"[OwnerUin]\"`.Use comma (,) 到 separate 多个 users, e.g `ID=\"100000000001\",ID=\"100000000002\"`",
//				},
//			},
//		},
//	}
//}

func ResourceTencentCloudCosBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCosBucketCreate,
		Read:   resourceTencentCloudCosBucketRead,
		Update: resourceTencentCloudCosBucketUpdate,
		Delete: resourceTencentCloudCosBucketDelete,
		Importer: &schema.ResourceImporter{
			State: helper.ImportWithDefaultValue(map[string]interface{}{
				"force_clean": false,
			}),
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateCosBucketName,
				Description:  "名称 存储桶 到 是 创建. 存储桶 格式 should 是 [自定义 名称]-[appid]，对于 示例 `mycos-1258798060`。",
			},
			"acl": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  s3.ObjectCannedACLPrivate,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{
					s3.ObjectCannedACLPrivate,
					s3.ObjectCannedACLPublicRead,
					s3.ObjectCannedACLPublicReadWrite,
				}),
				Description: "canned ACL 到 apply. 有效值：私有，公有-read，和 公有-read-write. 默认为 私有。",
			},
			"acl_body": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, olds, news string, d *schema.ResourceData) bool {
					return ACLBodyDiffFunc(olds, news, d)
				},
				DiffSuppressOnRefresh: true,
				ValidateFunc:          validateACLBody,
				Description:           "ACL XML 正文 对于 多个 grant info. NOTE: 此 argument 将 overwrite `acl`. Check https://intl.云.tencent.com/document/product/436/7737 对于 more detail。",
			},
			"encryption_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "服务器-side 加密 algorithm 到 使用. 有效 值 是 `AES256`，`KMS` 和 `SM4`。",
			},
			"kms_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "KMS Master 键 ID. 此 值 是 有效 仅 当 `encryption_algorithm` 是 集合 到 KMS. Set kms ID 到 指定 值 如果未指定， 默认值 kms ID 是 使用。",
			},
			"versioning_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable 存储桶 versioning. NOTE: `multi_az` 功能 是 true 对于 当前 存储桶，不能 disable 版本 control。",
			},
			"acceleration_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable 存储桶 acceleration。",
			},
			"force_clean": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Force cleanup all objects before delete 存储桶",
			},
			"replica_role": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"replica_rules", "versioning_enable"},
				Description:  "Request initiator identifier，格式: `qcs::cam::uin/<owneruin>:uin/<subuin>`. NOTE: 仅 `versioning_enable` 是 true 可以 configure 此 argument。",
			},
			"replica_rules": {
				Type:         schema.TypeList,
				Optional:     true,
				Description:  "列表 副本 规则. NOTE: 仅 `versioning_enable` 是 true 和 `replica_role` 集合 可以 configure 此 argument。",
				RequiredWith: []string{"replica_role", "versioning_enable"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "名称 特定 规则。",
						},
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "状态 identifier，可用值：`已启用`，`已禁用`。",
						},
						"priority": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Execution 优先级，用于handle scenarios 其中 目标 存储 buckets 是 same 和 多个 复制 规则 match same 对象. 注意: Supports setting positive integers 在 范围 的 1-1000. 优先级 值 的 different 规则 不能 是 duplicated. Storage 存储桶 复制 规则 必须 either all have 优先级 集合 或 all 不 have 优先级 集合. 当 all 规则 have 优先级 集合，overlapping prefixes 是 allowed 对于 different 规则 当 目标 存储 buckets 是 same. 当 different 规则 match same 对象， 规则 使用 smallest 优先级 值 将 是 triggered first. 当 none 的 规则 have 优先级 集合，overlapping prefixes 是 不 allowed 对于 different 规则。",
						},
						"prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Prefix matching 策略. Policies 不能 overlap; otherwise， 错误 将 是 返回. To match root directory，leave 此 参数 空。",
						},
						"filter": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "过滤器 objects 到 是 copied. 存储桶 功能 将 copy objects 该 match prefixes 和 标签 指定 在 过滤器 settings。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"and": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "当 filtering objects 到 是 copied，如果 both prefix 和 对象 标签 conditions 为必填项 simultaneously，或 如果 多个 对象 标签 conditions 是 needed，they 必须 是 enclosed 在 `And` statement。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"prefix": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "过滤器 objects 通过 prefix; 您 可以 指定at most 一个 prefix。",
												},
												"tag": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "当 filtering objects 到 是 copied，您 可以 使用 对象 标签 (多个 标签 是 支持) 作为 filtering criteria，使用 最大 的 10 标签 allowed. After adding 标签 作为 filtering criteria， `delete_marker_replication.状态` 选项 必须 是 集合 到 false。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "标签键",
															},
															"value": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "标签值",
															},
														},
													},
												},
											},
										},
									},
									"prefix": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "过滤器 objects 通过 prefix; 您 可以 指定at most 一个 prefix。",
									},
								},
							},
						},
						"destination_bucket": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Destination 存储桶 identifier，格式: `qcs::cos:<地域>::<bucketname-appid>`. NOTE: destination 存储桶 必须 启用 versioning。",
						},
						"destination_storage_class": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Storage class 的 destination，可用值：`Standard`，`Intelligent_Tiering`，`Standard_IA`. 默认为 following 当前 class 的 destination。",
						},
						"destination_encryption_kms_key_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "此 字段 必须 是 included 当 `source_selection_criteria.sse_kms_encrypted_objects.状态` 是 集合 到 已启用 It 是 用于指定KMS 键 用于KMS-encrypted objects copied 到 destination 存储桶",
						},
						"delete_marker_replication": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Synchronized deletion marker。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "是否synchronously delete 标签，支持 已禁用 或 已启用 默认值为 已启用，meaning 标签 将 是 删除 synchronously。",
									},
								},
							},
						},
						"source_selection_criteria": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "此 是 用于指定additional conditions 对于 objects 支持 通过 存储桶 复制 规则. Currently，仅 选项 到 replicate KMS-encrypted objects 是 支持。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sse_kms_encrypted_objects": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: "Choose 是否copy KMS-encrypted objects。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"status": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Choose 是否copy KMS encrypted objects; 支持 值 是 已启用 和 已禁用",
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
			"cors_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A 规则 的 Cross-Origin Resource Sharing (documented below)。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_origins": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "指定which origins 是 allowed。",
						},
						"allowed_methods": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "指定which methods 是 allowed. Can 是 `GET`，`PUT`，`POST`，`DELETE` 或 `HEAD`。",
						},
						"allowed_headers": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "指定which headers 是 allowed。",
						},
						"max_age_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "指定time （秒） 该 browser 可以 缓存 response 对于 preflight 请求。",
						},
						"expose_headers": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "指定expose 头部 在 response。",
						},
					},
				},
			},
			"origin_pull_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "存储桶 Origin-Pull settings。",
				Elem:        originPullRules(),
			},
			"origin_domain_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "存储桶 Origin 域名 settings。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "指定domain 主机",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "REST",
							Description: "指定origin 域名 类型，可用值：`REST`，`WEBSITE`，`ACCELERATE`，默认值：`REST`。",
						},
						"status": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "ENABLED",
							Description:  "域名 状态，默认值：`ENABLED`。",
							ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"ENABLED", "DISABLED"}),
						},
						//"force_replacement": {
						//	Type:		 schema.TypeString,
						//	Optional: 	 true,
						//	Description: "Specify 类型 到 replace exist 域名 resolve 记录.",
						//},
					},
				},
			},
			"lifecycle_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A 配置 的 对象 lifecycle management (documented below)。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "A 唯一 identifier 对于 规则. It 可以 是 up 到 255 字符。",
						},
						"filter_prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Object 键 prefix identifying 一个 或 more objects 到 其中 规则 applies。",
						},
						"transition": {
							Type:        schema.TypeSet,
							Optional:    true,
							Set:         transitionHash,
							Description: "指定a 周期 在 对象's transitions (documented below)。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: tccommon.ValidateCosBucketLifecycleTimestamp,
										Description:  "指定date after 其中 您 want corresponding 操作 到 take effect。",
									},
									"days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: tccommon.ValidateIntegerMin(0),
										Description:  "指定number 的 days after 对象 creation 当 特定 规则 操作 takes effect。",
									},
									"storage_class": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "指定storage class 到 其中 您 want 对象 到 transition. Available 值 include `STANDARD_IA`，`MAZ_STANDARD_IA`，`INTELLIGENT_TIERING`，`MAZ_INTELLIGENT_TIERING`，`ARCHIVE`，`DEEP_ARCHIVE`. For more 信息，please refer 到: https://云.tencent.com/document/product/436/33417。",
									},
								},
							},
						},
						"expiration": {
							Type:        schema.TypeSet,
							Optional:    true,
							Set:         expirationHash,
							MaxItems:    1,
							Description: "指定a 周期 在 对象's expire (documented below)。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: tccommon.ValidateCosBucketLifecycleTimestamp,
										Description:  "指定date after 其中 您 want corresponding 操作 到 take effect。",
									},
									"days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: tccommon.ValidateIntegerMin(0),
										Description:  "指定number 的 days after 对象 creation 当 特定 规则 操作 takes effect。",
									},
									"delete_marker": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "表示是否delete marker 的 expired 对象 将 是 removed。",
									},
								},
							},
						},
						"non_current_transition": {
							Type:        schema.TypeSet,
							Optional:    true,
							Set:         nonCurrentTransitionHash,
							Description: "指定a 周期 在 non 当前 对象's transitions。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"non_current_days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: tccommon.ValidateIntegerMin(0),
										Description:  "数量 days after non 当前 对象 creation 当 特定 规则 操作 takes effect。",
									},
									"storage_class": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "指定storage class 到 其中 您 want non 当前 对象 到 transition. Available 值 include `STANDARD_IA`，`MAZ_STANDARD_IA`，`INTELLIGENT_TIERING`，`MAZ_INTELLIGENT_TIERING`，`ARCHIVE`，`DEEP_ARCHIVE`. For more 信息，please refer 到: https://云.tencent.com/document/product/436/33417。",
									},
								},
							},
						},
						"non_current_expiration": {
							Type:        schema.TypeSet,
							Optional:    true,
							Set:         nonCurrentExpirationHash,
							MaxItems:    1,
							Description: "指定when non 当前 对象 versions shall expire。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"non_current_days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: tccommon.ValidateIntegerMin(0),
										Description:  "数量 days after non 当前 对象 creation 当 特定 规则 操作 takes effect. 最大 值 是 3650。",
									},
								},
							},
						},
						"abort_incomplete_multipart_upload": {
							Type:        schema.TypeSet,
							Optional:    true,
							Set:         abortIncompleteMultipartUploadHash,
							MaxItems:    1,
							Description: "Set 最大 时间 multipart upload 是 allowed 到 remain running。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days_after_initiation": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: tccommon.ValidateIntegerMin(1),
										Description:  "指定number 的 days after multipart upload starts 该 upload 必须 是 completed. 最大 值 是 3650。",
									},
								},
							},
						},
					},
				},
			},
			"website": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "A website 对象(documented below)。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index_document": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "COS 返回this 索引 document 当 requests 是 made 到 root 域名 或 any 的 subfolders。",
						},
						"error_document": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "An absolute 路径 到 document 到 返回 在 case 的 4XX 错误",
						},
						"redirect_all_requests_to": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"http", "https"}),
							Description:  "Redirects all 请求 configurations. 有效值：http，https. 默认为 `http`。",
						},
						"routing_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Routing 规则 配置. A RoutingRules 容器 可以 contain up 到 100 RoutingRule elements。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rules": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Routing 规则 列表。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"condition_error_code": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "指定error 代码 作为 match condition 对于 routing 规则. 有效值：仅 4xx 返回 codes，such 作为 403 或 404。",
												},
												"condition_prefix": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "指定object 键 prefix 作为 match condition 对于 routing 规则。",
												},
												"redirect_protocol": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "指定target 协议 对于 routing 规则. Only HTTPS 是 支持。",
												},
												"redirect_replace_key": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "指定target 对象 键 到 replace original 对象 键 在 请求。",
												},
												"redirect_replace_key_prefix": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "指定object 键 prefix 到 replace original prefix 在 请求. You 可以 集合 此 参数 仅 如果 condition 是 KeyPrefixEquals。",
												},
											},
										},
									},
								},
							},
						},
						"endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "`Endpoint` 的 静态 website。",
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 的 存储桶",
			},
			"log_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate 访问 日志 的 此 存储桶 到 是 saved 或 不. 默认为 `false`. 如果 集合 `true`， 访问 日志 将 是 saved 使用 `log_target_bucket`. To 启用 日志， full 访问 的 日志 服务 必须 是 granted. [Full Access 角色 Policy](https://intl.云.tencent.com/document/product/436/16920)。",
			},
			"log_target_bucket": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "目标 存储桶名称 其中 saves 访问 日志 的 此 存储桶 per 5 minutes. 日志 访问 文件 格式 是 `log_target_bucket`/`log_prefix`{YYYY}/{MM}/{DD}/{时间}_{random}_{索引}.gz. Only 有效 当 `log_enable` 是 `true`. 用户 必须 have full 访问 在 此 存储桶",
			},
			"log_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "prefix 日志 名称 其中 saves 访问 日志 的 此 存储桶 per 5 minutes. Eg. `MyLogPrefix/`. 日志 访问 文件 格式 是 `log_target_bucket`/`log_prefix`{YYYY}/{MM}/{DD}/{时间}_{random}_{索引}.gz. Only 有效 当 `log_enable` 是 `true`。",
			},
			"multi_az": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "表示是否create 存储桶 的 multi 可用 可用区",
			},
			"chdfs_ofs": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "表示是否create 存储桶 的 metadata acceleration. For more 信息，please refer 到 `https://www.tencentcloud.com/document/product/436/43305`。",
			},
			"enable_intelligent_tiering": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable intelligent tiering. NOTE: 当 intelligent tiering 配置 是 已启用，它 不能 是 turned 关闭 或 modified。",
			},
			"intelligent_tiering_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "指定limit 的 days 对于 standard-tier 数据 到 low-频率 数据 在 intelligent tiered 存储 配置，使用 可选 days 的 30，60，90. 默认值为 30。",
			},
			"intelligent_tiering_request_frequent": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "指定access 限制 对于 converting standard layer 数据 into low-频率 layer 数据 在 配置. 默认值为 once，其中 可以 是 使用 在 combination 使用 数量 days 到 achieve conversion effect. For 示例，如果 参数 是 集合 到 1 和 数量 访问 days 是 30，它 表示 该 objects 使用 less 比 一个 visit 在 30 consecutive days 将 是 reduced 从 standard layer 到 low 频率 layer。",
			},
			"intelligent_tiering_archiving_rule_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "列表 intelligent tiered 存储，archiving，和 deep archiving 规则. NOTE: 仅 `enable_intelligent_tiering` 是 true 可以 configure 此 argument。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "名称 intelligent tiering 规则 名称 列表 任务，使用 ID 集合 到 non-默认值 字符串，表示that 此 规则 是 conversion 规则 对于 archive 和 deep archive tiers。",
						},
						"status": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "表示是否intelligent tiering 规则 是 已启用 可能的值：已启用，已禁用 当 ID 是 `默认值`，仅 `已启用` 是 支持。",
						},
						"filter": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "指定configuration 信息 related 到 数据 transformation 在 intelligent tiered 存储 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"and": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "For filtering conditions，如果 both prefix 和 对象 标签 conditions 为必填项 simultaneously，they need 到 是 wrapped 使用 `And` 操作者",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"prefix": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "过滤器 objects 通过 prefix; 您 可以 指定at most 一个 prefix。",
												},
												"tag": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "当 filtering objects 对于 analysis，您 可以 使用 对象 标签 (多个 标签 是 支持) 作为 filtering criteria。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "标签键",
															},
															"value": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "标签值",
															},
														},
													},
												},
											},
										},
									},
									"prefix": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "过滤器 objects 通过 prefix; 您 可以 指定at most 一个 prefix。",
									},
									"tag": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "当 filtering objects 对于 analysis，您 可以 使用 对象 标签 (多个 标签 是 支持) 作为 filtering criteria。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "标签键",
												},
												"value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "标签值",
												},
											},
										},
									},
								},
							},
						},
						"tiering": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "指定configuration 信息 related 到 数据 transformation 在 intelligent tiered 存储 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_tier": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "当 `rule_id` 是 不 `默认值`，此 参数 是 用于指定archiving 或 deep archiving tier. possible 值 是: ARCHIVE_ACCESS，DEEP_ARCHIVE_ACCESS。",
									},
									"days": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "当 `rule_id` 是 不 集合 到 默认值，此 指定number 的 days after 其中 数据 是 transitioned 到 archive 或 deep archive tier 在 intelligent tiering 存储 配置. archive tier (ARCHIVE_ACCESS) 支持 范围 的 91 到 730 days. deep archive tier (DEEP_ARCHIVE_ACCESS) 支持 范围 的 180 到 730 days. Within same 规则， 数量 days 对于 deep archive tier 必须 是 greater 比 数量 days 对于 archive tier。",
									},
								},
							},
						},
					},
				},
			},
			"cdc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "CDC 集群 ID。",
			},
			"object_lock_configuration": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Object locking 配置. Once 已启用，此 功能 不能 是 已禁用",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Enable 对象 lock 配置。",
						},
						"rule": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Object locking 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Object lock 默认值 时长 (范围: 1-36500)。",
									},
									// "mode": {
									// 	Type:        schema.TypeString,
									// 	Optional:    true,
									// 	Description: "对象 lock 默认值 模式 仅 支持 enumerated 值 `COMPLIANCE`. 如果 此 字段 是 left blank, 它 defaults 到 `COMPLIANCE`.",
									// },
								},
							},
						},
					},
				},
			},
			//computed
			"cos_bucket_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL 的 此 COS 存储桶",
			},
		},
	}
}

func resourceTencentCloudCosBucketCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket.create")()

	var err error

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	bucket := d.Get("bucket").(string)
	acl := d.Get("acl").(string)
	role, roleOk := d.GetOk("replica_role")
	rule, ruleOk := d.GetOk("replica_rules")
	versioning := d.Get("versioning_enable").(bool)
	cdcId := d.Get("cdc_id").(string)

	if !versioning {
		if roleOk || role.(string) != "" {
			return fmt.Errorf("cannot configure role unless versioning enable")
		} else if ruleOk || len(rule.([]interface{})) > 0 {
			return fmt.Errorf("cannot configure replica rule unless versioning enable")
		}
	}

	cosService := CosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	useCosService, createOptions := getBucketPutOptions(d)

	if useCosService {
		// tencent
		err = cosService.TencentCosPutBucket(ctx, bucket, createOptions, cdcId)
	} else {
		// s3
		err = cosService.PutBucket(ctx, bucket, acl, cdcId)
	}
	if err != nil {
		return err
	}

	d.SetId(bucket)

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		if err := cosService.SetBucketTags(ctx, bucket, tags, cdcId); err != nil {
			return err
		}
	}

	return resourceTencentCloudCosBucketUpdate(d, meta)
}

func resourceTencentCloudCosBucketRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	bucket := d.Id()
	cosService := CosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	cdcId := d.Get("cdc_id").(string)
	code, header, err := cosService.TencentcloudHeadBucket(ctx, bucket, cdcId)
	if err != nil {
		if code == 404 {
			log.Printf("[WARN]%s bucket (%s) not found, error code (404)", logId, bucket)
			d.SetId("")
			return nil
		} else {
			return err
		}
	}

	if header != nil {
		if len(header["X-Cos-Bucket-Az-Type"]) > 0 && header["X-Cos-Bucket-Az-Type"][0] == "MAZ" {
			_ = d.Set("multi_az", true)
		}

		if len(header["X-Cos-Bucket-Arch"]) > 0 && header["X-Cos-Bucket-Arch"][0] == "OFS" {
			_ = d.Set("chdfs_ofs", true)
		}
	}

	ofs := d.Get("chdfs_ofs").(bool)
	cosDomain := meta.(tccommon.ProviderMeta).GetAPIV3Conn().CosDomain
	var cosBucketUrl string
	if cdcId == "" && cosDomain == "" {
		cosBucketUrl = fmt.Sprintf("%s.cos.%s.myqcloud.com", d.Id(), meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region)
	} else if cosDomain != "" {
		parsedURL, _ := url.Parse(cosDomain)
		parsedURL.Host = bucket + "." + parsedURL.Host
		cosBucketUrl = parsedURL.String()
	} else {
		cosBucketUrl = fmt.Sprintf("https://%s.%s.cos-cdc.%s.myqcloud.com", bucket, cdcId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region)
	}

	_ = d.Set("cos_bucket_url", cosBucketUrl)
	// set bucket in the import case
	if _, ok := d.GetOk("bucket"); !ok {
		_ = d.Set("bucket", d.Id())
	}

	if !ofs {
		// acl
		aclResult, err := cosService.GetBucketACL(ctx, bucket, cdcId)

		if err != nil {
			return err
		}

		aclBody, err := xml.Marshal(aclResult)
		if err != nil {
			return err
		}

		_ = d.Set("acl_body", string(aclBody))

		acl := GetBucketPublicACL(aclResult)

		_ = d.Set("acl", acl)

		if cdcId == "" && cosDomain == "" {
			originPullRules, err := cosService.GetBucketPullOrigin(ctx, bucket)
			if err != nil {
				return err
			}

			if err = d.Set("origin_pull_rules", originPullRules); err != nil {
				return fmt.Errorf("setting origin_pull_rules error: %v", err)
			}

			originDomainRules, err := cosService.GetBucketOriginDomain(ctx, bucket)
			if err != nil {
				return err
			}

			if err = d.Set("origin_domain_rules", originDomainRules); err != nil {
				return fmt.Errorf("setting origin_domain_rules error: %v", err)
			}

			replicaResult, err := cosService.GetBucketReplication(ctx, bucket, cdcId)
			if err != nil {
				return err
			}

			if replicaResult != nil {
				err := setBucketReplication(d, *replicaResult)
				if err != nil {
					return err
				}
			}
		}

		// read the website
		website, err := cosService.GetBucketWebsite(ctx, bucket, cdcId)
		if err != nil {
			return err
		}
		if len(website) > 0 && cosDomain == "" {
			// {bucket}.cos-website.{region}.myqcloud.com
			endPointUrl := fmt.Sprintf("%s.cos-website.%s.myqcloud.com", d.Id(), meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region)
			website[0]["endpoint"] = endPointUrl
		}
		if err = d.Set("website", website); err != nil {
			return fmt.Errorf("setting website error: %v", err)
		}

		// read the encryption algorithm
		encryption, kmsId, err := cosService.GetBucketEncryption(ctx, bucket, cdcId)
		if err != nil {
			return err
		}
		if err = d.Set("encryption_algorithm", encryption); err != nil {
			return fmt.Errorf("setting encryption error: %v", err)
		}
		if err = d.Set("kms_id", kmsId); err != nil {
			return fmt.Errorf("setting kms_id error: %v", err)
		}

		// read the versioning
		versioning, err := cosService.GetBucketVersioning(ctx, bucket, cdcId)
		if err != nil {
			return err
		}
		if err = d.Set("versioning_enable", versioning); err != nil {
			return fmt.Errorf("setting versioning_enable error: %v", err)
		}

		// read the acceleration
		acceleration, err := cosService.GetBucketAccleration(ctx, bucket, cdcId)
		if err != nil {
			return err
		}
		if err = d.Set("acceleration_enable", acceleration); err != nil {
			return fmt.Errorf("setting acceleration_enable error: %v", err)
		}
	}

	// read the cors
	corsRules, err := cosService.GetBucketCors(ctx, bucket, cdcId)
	if err != nil {
		return err
	}
	if err = d.Set("cors_rules", corsRules); err != nil {
		return fmt.Errorf("setting cors_rules error: %v", err)
	}

	// read the lifecycle
	lifecycleRules, err := cosService.GetBucketLifecycle(ctx, bucket, cdcId)
	if err != nil {
		return err
	}
	if err = d.Set("lifecycle_rules", lifecycleRules); err != nil {
		return fmt.Errorf("setting lifecycle_rules error: %v", err)
	}

	//read the log
	logEnable, logTargetBucket, logPrefix, err := cosService.GetBucketLogStatus(ctx, bucket, cdcId)
	if err != nil {
		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() != "UnSupportedLoggingRegion" {
				return err
			}
		}
	} else {
		_ = d.Set("log_enable", logEnable)
		_ = d.Set("log_target_bucket", logTargetBucket)
		_ = d.Set("log_prefix", logPrefix)
	}

	// read the tags
	tags, err := cosService.GetBucketTags(ctx, bucket, cdcId)
	if err != nil {
		return fmt.Errorf("get tags failed: %v", err)
	}
	if len(tags) > 0 {
		_ = d.Set("tags", tags)
	}

	//read intelligent tiering
	if !strings.Contains(cosBucketUrl, ".cos-cdc.") {
		result, err := cosService.BucketGetIntelligentTiering(ctx, bucket, cdcId)
		if err != nil {
			return fmt.Errorf("get intelligent tiering failed: %v", err)
		}
		if result != nil {
			if result.Status == "Enabled" {
				_ = d.Set("enable_intelligent_tiering", true)
			} else {
				_ = d.Set("enable_intelligent_tiering", false)
			}

			if result.Transition != nil {
				_ = d.Set("intelligent_tiering_days", result.Transition.Days)
				_ = d.Set("intelligent_tiering_request_frequent", result.Transition.RequestFrequent)
			}
		}

		//read intelligent tiering archiving rule list
		respData, err := cosService.BucketGetIntelligentTieringArchivingRuleList(ctx, bucket, cdcId)
		if err != nil {
			return fmt.Errorf("get intelligent tiering archiving rule list failed: %v", err)
		}

		if respData != nil {
			if respData.Configurations != nil && len(respData.Configurations) > 0 {
				intelligentTieringArchivingRules := make([]map[string]interface{}, 0, len(respData.Configurations))
				for _, config := range respData.Configurations {
					intelligentTieringArchivingRule := make(map[string]interface{})
					if config.Id == "default" {
						continue
					}

					intelligentTieringArchivingRule["rule_id"] = config.Id
					intelligentTieringArchivingRule["status"] = config.Status

					if config.Filter != nil {
						dMap := make(map[string]interface{})
						if config.Filter.And != nil {
							andMap := make(map[string]interface{}, 0)
							if config.Filter.And.Prefix != "" {
								andMap["prefix"] = config.Filter.And.Prefix
							}

							if config.Filter.And.Tag != nil && len(config.Filter.And.Tag) > 0 {
								tagList := make([]map[string]interface{}, 0, len(config.Filter.And.Tag))
								for _, item := range config.Filter.And.Tag {
									dMap := make(map[string]interface{})
									dMap["key"] = item.Key
									dMap["value"] = item.Value
									tagList = append(tagList, dMap)
								}

								andMap["tag"] = tagList
							}

							dMap["and"] = []interface{}{andMap}
						}

						if config.Filter.Prefix != "" {
							dMap["prefix"] = config.Filter.Prefix
						}

						if config.Filter.Tag != nil && len(config.Filter.Tag) > 0 {
							tagList := make([]map[string]interface{}, 0, len(config.Filter.Tag))
							for _, item := range config.Filter.Tag {
								dMap := make(map[string]interface{})
								dMap["key"] = item.Key
								dMap["value"] = item.Value
								tagList = append(tagList, dMap)
							}

							dMap["tag"] = tagList
						}

						intelligentTieringArchivingRule["filter"] = []interface{}{dMap}
					}

					if config.Tiering != nil && len(config.Tiering) > 0 {
						tieringList := make([]map[string]interface{}, 0, len(config.Tiering))
						for _, item := range config.Tiering {
							dMap := make(map[string]interface{})
							dMap["access_tier"] = item.AccessTier
							dMap["days"] = item.Days
							tieringList = append(tieringList, dMap)
						}

						intelligentTieringArchivingRule["tiering"] = tieringList
					}

					intelligentTieringArchivingRules = append(intelligentTieringArchivingRules, intelligentTieringArchivingRule)
				}

				_ = d.Set("intelligent_tiering_archiving_rule_list", intelligentTieringArchivingRules)
			}
		}
	}

	//read object lock config
	objLockData, err := cosService.BucketGetObjectLockConfiguration(ctx, bucket, cdcId)
	if err != nil {
		return fmt.Errorf("get object lock configuration failed: %v", err)
	}

	if objLockData != nil {
		dMap := make(map[string]interface{})
		if objLockData.ObjectLockEnabled == "Enabled" {
			dMap["enabled"] = true
		} else {
			dMap["enabled"] = false
		}

		if objLockData.Rule != nil {
			ruleList := make([]map[string]interface{}, 0, 1)
			ruleMap := make(map[string]interface{})
			ruleMap["days"] = objLockData.Rule.Days
			ruleList = append(ruleList, ruleMap)
			dMap["rule"] = ruleList
		}

		_ = d.Set("object_lock_configuration", []interface{}{dMap})
	}

	return nil
}

func resourceTencentCloudCosBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cosService := CosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	d.Partial(true)

	cdcId := d.Get("cdc_id").(string)
	if d.HasChange("enable_intelligent_tiering") || d.HasChange("intelligent_tiering_days") || d.HasChange("intelligent_tiering_request_frequent") {
		old, new := d.GetChange("enable_intelligent_tiering")
		if old.(bool) && !new.(bool) {
			return fmt.Errorf("enable_intelligent_tiering, intelligent_tiering_days and intelligent_tiering_request_frequent not support change!")
		}
		var transition cos.BucketIntelligentTieringTransition
		if v, ok := d.GetOk("intelligent_tiering_days"); ok {
			transition.Days = v.(int)
		} else {
			transition.Days = 30
		}
		if v, ok := d.GetOk("intelligent_tiering_request_frequent"); ok {
			transition.RequestFrequent = v.(int)
		} else {
			transition.RequestFrequent = 1
		}

		if v, ok := d.GetOk("enable_intelligent_tiering"); ok && v.(bool) {
			opt := &cos.BucketPutIntelligentTieringOptions{
				Status:     "Enabled",
				Transition: &transition,
			}
			err := cosService.BucketPutIntelligentTiering(ctx, d.Id(), opt, cdcId)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("intelligent_tiering_archiving_rule_list") {
		oldInterface, newInterface := d.GetChange("intelligent_tiering_archiving_rule_list")
		oldList := oldInterface.([]interface{})
		newList := newInterface.([]interface{})
		if len(oldList) > 0 {
			ruleIds := make([]string, 0, len(oldList))
			for _, item := range oldList {
				dMap := item.(map[string]interface{})
				if v, ok := dMap["rule_id"].(string); ok && v != "" && v != "default" {
					ruleIds = append(ruleIds, v)
				}
			}

			if len(ruleIds) > 0 {
				err := cosService.BucketDeleteIntelligentTieringArchivingRule(ctx, d.Id(), cdcId, ruleIds)
				if err != nil {
					return err
				}
			}
		}

		if len(newList) > 0 {
			rules := make([]*cos.BucketPutIntelligentTieringOptions, 0)
			for _, item := range newList {
				dMap := item.(map[string]interface{})
				rule := &cos.BucketPutIntelligentTieringOptions{}
				if v, ok := dMap["rule_id"].(string); ok && v != "" && v != "default" {
					rule.Id = v
				}

				if v, ok := dMap["status"].(string); ok && v != "" {
					rule.Status = v
				}

				if v, ok := dMap["filter"]; ok {
					for _, filterItem := range v.([]interface{}) {
						filterMap := filterItem.(map[string]interface{})
						filter := &cos.BucketIntelligentTieringFilter{}
						if v, ok := filterMap["and"]; ok {
							for _, andItem := range v.([]interface{}) {
								andMap := andItem.(map[string]interface{})
								and := &cos.BucketIntelligentTieringFilterAnd{}
								if v, ok := andMap["prefix"].(string); ok && v != "" {
									and.Prefix = v
								}

								if v, ok := andMap["tag"]; ok {
									for _, tag := range v.([]interface{}) {
										tagMap := tag.(map[string]interface{})
										tmpTag := cos.BucketTaggingTag{}
										if v, ok := tagMap["key"].(string); ok && v != "" {
											tmpTag.Key = v
										}

										if v, ok := tagMap["value"].(string); ok && v != "" {
											tmpTag.Value = v
										}

										and.Tag = append(and.Tag, &tmpTag)
									}
								}

								filter.And = and
							}
						}

						if v, ok := filterMap["prefix"].(string); ok && v != "" {
							filter.Prefix = v
						}

						if v, ok := filterMap["tag"]; ok {
							for _, tag := range v.([]interface{}) {
								tagMap := tag.(map[string]interface{})
								tmpTag := &cos.BucketTaggingTag{}
								if v, ok := tagMap["key"].(string); ok && v != "" {
									tmpTag.Key = v
								}

								if v, ok := tagMap["value"].(string); ok && v != "" {
									tmpTag.Value = v
								}

								filter.Tag = append(filter.Tag, tmpTag)
							}
						}

						rule.Filter = filter
					}
				}

				if v, ok := dMap["tiering"]; ok {
					for _, tieringItem := range v.([]interface{}) {
						tieringMap := tieringItem.(map[string]interface{})
						tiering := &cos.BucketIntelligentTieringTransition{}
						if v, ok := tieringMap["access_tier"].(string); ok && v != "" {
							tiering.AccessTier = v
						}

						if v, ok := tieringMap["days"].(int); ok {
							tiering.Days = v
						}

						rule.Tiering = append(rule.Tiering, tiering)
					}
				}

				rules = append(rules, rule)
			}

			if len(rules) > 0 {
				err := cosService.BucketPutIntelligentTieringArchivingRule(ctx, d.Id(), cdcId, rules)
				if err != nil {
					return err
				}
			}
		}
	}

	if d.HasChange("acl") {
		bucket := d.Get("bucket").(string)
		err := waitAclEnable(ctx, meta, bucket, cdcId)
		if err != nil {
			return err
		}

		err = resourceTencentCloudCosBucketAclUpdate(ctx, meta, d)
		if err != nil {
			return err
		}
	}

	if d.HasChange("acl_body") {
		body := d.Get("acl_body")
		bucket := d.Get("bucket").(string)
		err := waitAclEnable(ctx, meta, bucket, cdcId)
		if err != nil {
			return err
		}

		if err := resourceTencentCloudCosBucketOriginACLBodyUpdate(ctx, cosService, d); err != nil {
			return err
		}
		_ = d.Set("acl_body", body)
	}

	if d.HasChange("cors_rules") {
		err := resourceTencentCloudCosBucketCorsUpdate(ctx, meta, d)
		if err != nil {
			return err
		}

	}

	if d.HasChange("origin_pull_rules") {
		rules := d.Get("origin_pull_rules")
		err := resourceTencentCloudCosBucketOriginPullUpdate(ctx, cosService, d)
		if err != nil {
			return err
		}
		_ = d.Set("origin_pull_rules", rules)
	}

	if d.HasChange("origin_domain_rules") {
		rules := d.Get("origin_domain_rules")
		if err := resourceTencentCloudCosBucketOriginDomainUpdate(ctx, cosService, d); err != nil {
			return err
		}
		_ = d.Set("origin_domain_rules", rules)
	}

	if d.HasChange("lifecycle_rules") {
		err := resourceTencentCloudCosBucketLifecycleUpdate(ctx, meta, d)
		if err != nil {
			return err
		}

	}

	if d.HasChange("website") {
		err := resourceTencentCloudCosBucketWebsiteUpdate(ctx, meta, d)
		if err != nil {
			return err
		}

	}

	if d.HasChange("encryption_algorithm") || d.HasChange("kms_id") {
		err := resourceTencentCloudCosBucketEncryptionUpdate(ctx, meta, d)
		if err != nil {
			return err
		}
	}

	if d.HasChange("versioning_enable") {
		err := resourceTencentCloudCosBucketVersioningUpdate(ctx, meta, d)
		if err != nil {
			return err
		}

	}

	if d.HasChange("acceleration_enable") {
		err := resourceTencentCloudCosBucketAccelerationUpdate(ctx, meta, d)
		if err != nil {
			return err
		}

	}

	if d.HasChange("replica_role") || d.HasChange("replica_rules") {
		err := resourceTencentCloudCosBucketReplicaUpdate(ctx, cosService, d)

		if err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		bucket := d.Id()

		cosService := CosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		if err := cosService.SetBucketTags(ctx, bucket, helper.GetTags(d, "tags"), cdcId); err != nil {
			return err
		}

	}

	if d.HasChange("log_enable") || d.HasChange("log_target_bucket") || d.HasChange("log_prefix") {
		err := resourceTencentCloudCosBucketLogStatusUpdate(ctx, meta, d)
		if err != nil {
			return err
		}

	}

	if d.HasChange("object_lock_configuration") {
		if v, ok := d.GetOk("object_lock_configuration"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				objectLockConfig := &cos.BucketPutObjectLockOptions{}
				if v, ok := dMap["enabled"].(bool); ok {
					if v {
						objectLockConfig.ObjectLockEnabled = "Enabled"
					} else {
						objectLockConfig.ObjectLockEnabled = "Disabled"
					}
				}

				if v, ok := dMap["rule"]; ok {
					for _, ruleItem := range v.([]interface{}) {
						ruleMap := ruleItem.(map[string]interface{})
						rule := &cos.ObjectLockRule{}
						if v, ok := ruleMap["days"].(int); ok {
							rule.Days = v
						}

						objectLockConfig.Rule = rule
					}
				}

				err := cosService.BucketPutObjectLockConfiguration(ctx, d.Id(), cdcId, objectLockConfig)
				if err != nil {
					return err
				}
			}
		}
	}

	d.Partial(false)

	// wait for update cache
	// if not, the data may be outdated.
	time.Sleep(3 * time.Second)

	return resourceTencentCloudCosBucketRead(d, meta)
}

func waitAclEnable(ctx context.Context, meta interface{}, bucket string, cdcId string) error {
	logId := tccommon.GetLogId(ctx)
	cosService := CosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		aclResult, e := cosService.GetBucketACL(ctx, bucket, cdcId)
		if e != nil {
			if strings.Contains(e.Error(), "NoSuchBucket") {
				log.Printf("[CRITAL][retry]%s api[%s] because of bucket[%s] still on creating, need try again.\n", logId, "GetBucketACL", bucket)
				return resource.RetryableError(fmt.Errorf("[CRITAL][retry]%s api[%s] it still on creating, need try again.\n", logId, "GetBucketACL"))
			}
			log.Printf("[CRITAL]%s api[%s] fail when try to update acl, reason[%s]\n", logId, "GetBucketACL", e.Error())
			return resource.NonRetryableError(e)
		}

		if aclResult == nil {
			return resource.RetryableError(fmt.Errorf("[CRITAL][retry]%s api[%s] it still on creating, need try again.\n", logId, "GetBucketACL"))
		}
		return nil
	})
	return err
}

func resourceTencentCloudCosBucketDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	bucket := d.Id()
	forced := d.Get("force_clean").(bool)
	versioned := d.Get("versioning_enable").(bool)
	multiAz := d.Get("multi_az").(bool)
	cosService := CosService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	cdcId := d.Get("cdc_id").(string)
	err := cosService.DeleteBucket(ctx, bucket, forced, versioned, cdcId, multiAz)
	if err != nil {
		return err
	}

	// wait for bucket 404, means deleted
	err = resource.Retry(tccommon.ReadRetryTimeout*10, func() *resource.RetryError {
		code, _, e := cosService.TencentcloudHeadBucket(ctx, bucket, cdcId)
		if e != nil {
			if code == 404 {
				log.Printf("[WARN]%s bucket (%s) not found, error code (404)", logId, bucket)
				return nil
			} else {
				return resource.NonRetryableError(e)
			}
		}

		return resource.RetryableError(fmt.Errorf("Waiting for cos bucket [%s] deleted...", bucket))
	})

	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudCosBucketEncryptionUpdate(ctx context.Context, meta interface{}, d *schema.ResourceData) error {
	logId := tccommon.GetLogId(ctx)

	bucket := d.Get("bucket").(string)
	encryption := d.Get("encryption_algorithm").(string)
	kmsId := d.Get("kms_id").(string)
	cdcId := d.Get("cdc_id").(string)
	if encryption == "" {
		request := s3.DeleteBucketEncryptionInput{
			Bucket: aws.String(bucket),
		}
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).DeleteBucketEncryption(&request)

		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete bucket encryption", request.String(), err.Error())
			return fmt.Errorf("cos delete bucket error: %s, bucket: %s", err.Error(), bucket)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, "delete bucket encryption", request.String(), response.String())

		return nil
	}

	request := s3.PutBucketEncryptionInput{
		Bucket: aws.String(bucket),
	}
	request.ServerSideEncryptionConfiguration = &s3.ServerSideEncryptionConfiguration{}
	rules := make([]*s3.ServerSideEncryptionRule, 0)
	defaultRule := &s3.ServerSideEncryptionByDefault{
		SSEAlgorithm:   aws.String(encryption),
		KMSMasterKeyID: aws.String(kmsId),
	}
	rule := &s3.ServerSideEncryptionRule{
		ApplyServerSideEncryptionByDefault: defaultRule,
	}
	rules = append(rules, rule)
	request.ServerSideEncryptionConfiguration.Rules = rules
	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).PutBucketEncryption(&request)

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "put bucket encryption", request.String(), err.Error())
		return fmt.Errorf("cos put bucket encryption error: %s, bucket: %s", err.Error(), bucket)
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "put bucket encryption", request.String(), response.String())

	return nil
}

func resourceTencentCloudCosBucketVersioningUpdate(ctx context.Context, meta interface{}, d *schema.ResourceData) error {
	logId := tccommon.GetLogId(ctx)

	bucket := d.Get("bucket").(string)
	versioning := d.Get("versioning_enable").(bool)
	cdcId := d.Get("cdc_id").(string)
	status := "Suspended"
	if versioning {
		status = "Enabled"
	}
	request := s3.PutBucketVersioningInput{
		Bucket: aws.String(bucket),
		VersioningConfiguration: &s3.VersioningConfiguration{
			Status: aws.String(status),
		},
	}
	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).PutBucketVersioning(&request)

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "put bucket encryption", request.String(), err.Error())
		return fmt.Errorf("cos put bucket encryption error: %s, bucket: %s", err.Error(), bucket)
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "put bucket encryption", request.String(), response.String())

	return nil
}

func resourceTencentCloudCosBucketAccelerationUpdate(ctx context.Context, meta interface{}, d *schema.ResourceData) error {
	logId := tccommon.GetLogId(ctx)

	bucket := d.Get("bucket").(string)
	enabled := d.Get("acceleration_enable").(bool)
	cdcId := d.Get("cdc_id").(string)
	status := "Suspended"
	if enabled {
		status = "Enabled"
	}

	opt := &cos.BucketPutAccelerateOptions{
		Status: status,
	}
	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTencentCosClientNew(bucket, cdcId).Bucket.PutAccelerate(ctx, opt)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, status [%s], reason[%s]\n",
			logId, "put bucket acceleration", opt.Status, err.Error())
		return fmt.Errorf("cos put bucket acceleration error: %s, bucket: %s", err.Error(), bucket)
	}
	rb, _ := ioutil.ReadAll(response.Body)
	body, _ := json.Marshal(rb)
	log.Printf("[DEBUG]%s api[%s] success, status [%s], response body [%s]\n",
		logId, "put bucket acceleration", opt.Status, string(body))

	return err
}

func resourceTencentCloudCosBucketReplicaUpdate(ctx context.Context, service CosService, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	oldRole, newRole := d.GetChange("replica_role")
	oldRules, newRules := d.GetChange("replica_rules")
	cdcId := d.Get("cdc_id").(string)
	oldRuleLength := len(oldRules.([]interface{}))
	newRuleLength := len(newRules.([]interface{}))

	// check if remove
	if oldRole.(string) != "" && newRole.(string) == "" || oldRuleLength > 0 && newRuleLength == 0 {
		result, err := service.GetBucketReplication(ctx, bucket, cdcId)
		if err != nil {
			return err
		}

		if result != nil {
			err := service.DeleteBucketReplication(ctx, d.Get("bucket").(string), cdcId)
			if err != nil {
				return err
			}
		}
	} else if newRole.(string) != "" || newRuleLength > 0 {
		role, rules, _ := getBucketReplications(d)
		err := service.PutBucketReplication(ctx, d.Get("bucket").(string), role, rules, cdcId)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceTencentCloudCosBucketAclUpdate(ctx context.Context, meta interface{}, d *schema.ResourceData) error {
	logId := tccommon.GetLogId(ctx)

	bucket := d.Get("bucket").(string)
	acl := d.Get("acl").(string)
	cdcId := d.Get("cdc_id").(string)
	request := s3.PutBucketAclInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String(acl),
	}
	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).PutBucketAcl(&request)

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "put bucket acl", request.String(), err.Error())
		return fmt.Errorf("cos put bucket error: %s, bucket: %s", err.Error(), bucket)
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "put bucket acl", request.String(), response.String())

	return nil
}

func resourceTencentCloudCosBucketCorsUpdate(ctx context.Context, meta interface{}, d *schema.ResourceData) error {
	logId := tccommon.GetLogId(ctx)

	bucket := d.Get("bucket").(string)
	cors := d.Get("cors_rules").([]interface{})
	cdcId := d.Get("cdc_id").(string)

	if len(cors) == 0 {
		request := s3.DeleteBucketCorsInput{
			Bucket: aws.String(bucket),
		}
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).DeleteBucketCors(&request)

		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete bucket cors", request.String(), err.Error())
			return fmt.Errorf("cos delete bucket cors error: %s, bucket: %s", err.Error(), bucket)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, "delete bucket cors", request.String(), response.String())
	} else {
		rules := make([]*s3.CORSRule, 0, len(cors))
		for _, item := range cors {
			corsMap := item.(map[string]interface{})
			rule := &s3.CORSRule{}
			for k, v := range corsMap {
				if k == "max_age_seconds" {
					rule.MaxAgeSeconds = aws.Int64(int64(v.(int)))
				} else {
					vMap := make([]*string, len(v.([]interface{})))
					for i, value := range v.([]interface{}) {
						if str, ok := value.(string); ok {
							vMap[i] = aws.String(str)
						}
					}
					switch k {
					case "allowed_origins":
						rule.AllowedOrigins = vMap
					case "allowed_methods":
						rule.AllowedMethods = vMap
					case "allowed_headers":
						rule.AllowedHeaders = vMap
					case "expose_headers":
						rule.ExposeHeaders = vMap
					}
				}
			}
			rules = append(rules, rule)
		}
		request := s3.PutBucketCorsInput{
			Bucket: aws.String(bucket),
			CORSConfiguration: &s3.CORSConfiguration{
				CORSRules: rules,
			},
		}
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).PutBucketCors(&request)

		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "put bucket cors", request.String(), err.Error())
			return fmt.Errorf("cos put bucket cors error: %s, bucket: %s", err.Error(), bucket)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, "put bucket cors", request.String(), response.String())
	}
	return nil
}

func resourceTencentCloudCosBucketLifecycleUpdate(ctx context.Context, meta interface{}, d *schema.ResourceData) error {
	logId := tccommon.GetLogId(ctx)

	bucket := d.Get("bucket").(string)
	lifecycleRules := d.Get("lifecycle_rules").([]interface{})
	cdcId := d.Get("cdc_id").(string)
	if len(lifecycleRules) == 0 {
		request := s3.DeleteBucketLifecycleInput{
			Bucket: aws.String(bucket),
		}
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).DeleteBucketLifecycle(&request)

		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete bucket lifecycle", request.String(), err.Error())
			return fmt.Errorf("cos delete bucket lifecycle error: %s, bucket: %s", err.Error(), bucket)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, "delete bucket lifecycle", request.String(), response.String())
	} else {
		rules := make([]*s3.LifecycleRule, 0, len(lifecycleRules))
		for i, lifecycleRule := range lifecycleRules {
			r := lifecycleRule.(map[string]interface{})
			rule := &s3.LifecycleRule{}
			id, ok := r["id"].(string)
			if ok {
				rule.ID = &id
			}
			rule.Status = helper.String(s3.ExpirationStatusEnabled)
			prefix := r["filter_prefix"].(string)
			rule.Filter = &s3.LifecycleRuleFilter{
				Prefix: &prefix,
			}

			// Transitions
			transitions := d.Get(fmt.Sprintf("lifecycle_rules.%d.transition", i)).(*schema.Set).List()
			if len(transitions) > 0 {
				rule.Transitions = make([]*s3.Transition, 0, len(transitions))
				for _, transition := range transitions {
					transitionValue := transition.(map[string]interface{})
					t := &s3.Transition{}
					if val, ok := transitionValue["date"].(string); ok && val != "" {
						date, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", val))
						if err != nil {
							return fmt.Errorf("parsing cos bucket lifecycle transition date(%s) error: %s", val, err.Error())
						}
						t.Date = aws.Time(date)
					} else if val, ok := transitionValue["days"].(int); ok && val >= 0 {
						t.Days = aws.Int64(int64(val))
					}
					if val, ok := transitionValue["storage_class"].(string); ok && val != "" {
						t.StorageClass = aws.String(val)
					}

					rule.Transitions = append(rule.Transitions, t)
				}
			}

			// Expiration
			expirations := d.Get(fmt.Sprintf("lifecycle_rules.%d.expiration", i)).(*schema.Set).List()
			if len(expirations) > 0 {
				expiration := expirations[0].(map[string]interface{})
				e := &s3.LifecycleExpiration{}

				if val, ok := expiration["date"].(string); ok && val != "" {
					date, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", val))
					if err != nil {
						return fmt.Errorf("parsing cos bucket lifecycle expiration data(%s) error: %s", val, err.Error())
					}
					e.Date = aws.Time(date)
				} else if val, ok := expiration["days"].(int); ok && val > 0 {
					e.Days = aws.Int64(int64(val))
				}

				if val, ok := expiration["delete_marker"].(bool); ok && val {
					e.ExpiredObjectDeleteMarker = helper.Bool(true)
				}

				rule.Expiration = e
			}

			// Non Current Transitions
			nonCurrentTransitions := d.Get(fmt.Sprintf("lifecycle_rules.%d.non_current_transition", i)).(*schema.Set).List()
			if len(nonCurrentTransitions) > 0 {
				rule.NoncurrentVersionTransitions = make([]*s3.NoncurrentVersionTransition, 0, len(transitions))
				for _, transition := range nonCurrentTransitions {
					transitionValue := transition.(map[string]interface{})
					t := &s3.NoncurrentVersionTransition{}
					if val, ok := transitionValue["non_current_days"].(int); ok && val >= 0 {
						t.NoncurrentDays = aws.Int64(int64(val))
					}
					if val, ok := transitionValue["storage_class"].(string); ok && val != "" {
						t.StorageClass = aws.String(val)
					}

					rule.NoncurrentVersionTransitions = append(rule.NoncurrentVersionTransitions, t)
				}
			}

			// Non Current Expiration
			nonCurrentExpirations := d.Get(fmt.Sprintf("lifecycle_rules.%d.non_current_expiration", i)).(*schema.Set).List()
			if len(nonCurrentExpirations) > 0 {
				nonCurrentExpiration := nonCurrentExpirations[0].(map[string]interface{})
				e := &s3.NoncurrentVersionExpiration{}

				if val, ok := nonCurrentExpiration["non_current_days"].(int); ok && val > 0 {
					e.NoncurrentDays = aws.Int64(int64(val))
				}

				rule.NoncurrentVersionExpiration = e
			}

			// AbortIncompleteMultipartUpload
			abortIncompleteMultipartUploads := d.Get(fmt.Sprintf("lifecycle_rules.%d.abort_incomplete_multipart_upload", i)).(*schema.Set).List()
			if len(abortIncompleteMultipartUploads) > 0 {
				abortIncompleteMultipartUpload := abortIncompleteMultipartUploads[0].(map[string]interface{})
				e := &s3.AbortIncompleteMultipartUpload{}

				if val, ok := abortIncompleteMultipartUpload["days_after_initiation"].(int); ok && val > 0 {
					e.DaysAfterInitiation = aws.Int64(int64(val))
				}

				rule.AbortIncompleteMultipartUpload = e
			}
			rules = append(rules, rule)
		}

		request := s3.PutBucketLifecycleConfigurationInput{
			Bucket: aws.String(bucket),
			LifecycleConfiguration: &s3.BucketLifecycleConfiguration{
				Rules: rules,
			},
		}
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).PutBucketLifecycleConfiguration(&request)

		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "put bucket lifecycle", request.String(), err.Error())
			return fmt.Errorf("cos put bucket lifecycle error: %s, bucket: %s", err.Error(), bucket)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, "put bucket lifecycle", request.String(), response.String())
	}

	return nil
}

func resourceTencentCloudCosBucketWebsiteUpdate(ctx context.Context, meta interface{}, d *schema.ResourceData) error {
	logId := tccommon.GetLogId(ctx)

	bucket := d.Get("bucket").(string)
	website := d.Get("website").([]interface{})
	cdcId := d.Get("cdc_id").(string)

	if len(website) == 0 {
		request := s3.DeleteBucketWebsiteInput{
			Bucket: aws.String(bucket),
		}

		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).DeleteBucketWebsite(&request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete bucket website", request.String(), err.Error())
			return fmt.Errorf("cos delete bucket website error: %s, bucket: %s", err.Error(), bucket)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, "delete bucket website", request.String(), response.String())
	} else {
		if cdcId != "" {
			return fmt.Errorf("cdc cos not support set website.\n")
		}

		var w map[string]interface{}
		if website[0] != nil {
			w = website[0].(map[string]interface{})
		} else {
			w = make(map[string]interface{})
		}

		websiteConfiguration := cos.BucketPutWebsiteOptions{}
		if v, ok := w["index_document"].(string); ok && v != "" {
			websiteConfiguration.Index = v
		}

		if v, ok := w["error_document"].(string); ok && v != "" {
			websiteConfiguration.Error = &cos.ErrorDocument{
				Key: v,
			}
		}

		if v, ok := w["redirect_all_requests_to"].(string); ok && v != "" {
			websiteConfiguration.RedirectProtocol = &cos.RedirectRequestsProtocol{
				Protocol: v,
			}
		}

		if v, ok := w["routing_rules"]; ok {
			websiteRoutingRules := cos.WebsiteRoutingRules{}
			if len(v.([]interface{})) > 0 {
				for _, item := range v.([]interface{}) {
					if rules, ok := item.(map[string]interface{}); ok && rules != nil {
						if v, ok := rules["rules"]; ok {
							wbRules := []cos.WebsiteRoutingRule{}
							for _, rule := range v.([]interface{}) {
								if dMap, ok := rule.(map[string]interface{}); ok && rules != nil {
									wbRule := cos.WebsiteRoutingRule{}
									if v, ok := dMap["condition_error_code"].(string); ok && v != "" {
										wbRule.ConditionErrorCode = v
									}

									if v, ok := dMap["condition_prefix"].(string); ok && v != "" {
										wbRule.ConditionPrefix = v
									}

									if v, ok := dMap["redirect_protocol"].(string); ok && v != "" {
										wbRule.RedirectProtocol = v
									}

									if v, ok := dMap["redirect_replace_key"].(string); ok && v != "" {
										wbRule.RedirectReplaceKey = v
									}

									if v, ok := dMap["redirect_replace_key_prefix"].(string); ok && v != "" {
										wbRule.RedirectReplaceKeyPrefix = v
									}

									wbRules = append(wbRules, wbRule)
								}
							}

							websiteRoutingRules.Rules = wbRules
						}
					}
				}

				websiteConfiguration.RoutingRules = &websiteRoutingRules
			}
		}

		request := websiteConfiguration
		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTencentCosClientNew(bucket, cdcId).Bucket.PutWebsite(ctx, &request)
		if err != nil {
			return fmt.Errorf("cos put bucket website error: %s, bucket: %s", err.Error(), bucket)
		}

		reqBytes, _ := json.Marshal(request)
		respBytes, _ := json.Marshal(response.Response.Body)
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, "put bucket website", string(reqBytes), string(respBytes))
	}

	return nil
}

func resourceTencentCloudCosBucketLogStatusUpdate(ctx context.Context, meta interface{}, d *schema.ResourceData) error {
	logId := tccommon.GetLogId(ctx)

	bucket := d.Id()

	logSwitch := d.Get("log_enable").(bool)
	cdcId := d.Get("cdc_id").(string)
	if logSwitch {
		if d.HasChange("log_target_bucket") || d.HasChange("log_prefix") {
			targetBucket := d.Get("log_target_bucket").(string)
			logPrefix := d.Get("log_prefix").(string)
			//check
			if targetBucket == "" || logPrefix == "" {
				return fmt.Errorf("log_target_bucket and log_prefix should set valid value when log_enable is true")
			}

			//set log target bucket and prefix
			//grant are solved by the tencentcloud_cam_role_attachment resource
			request := &s3.PutBucketLoggingInput{
				Bucket: aws.String(bucket),
				BucketLoggingStatus: &s3.BucketLoggingStatus{
					LoggingEnabled: &s3.LoggingEnabled{
						TargetBucket: aws.String(targetBucket),
						TargetPrefix: aws.String(logPrefix),
					},
				},
			}

			resp, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).PutBucketLogging(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, "cos enable log error", request.String(), err.Error())
				return fmt.Errorf("cos enable log error: %s, bucket: %s", err.Error(), bucket)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], resp[%s]\n",
				logId, "cos enable log success", request.String(), resp.String())
		}
	} else {
		targetBucket := d.Get("log_target_bucket").(string)
		logPrefix := d.Get("log_prefix").(string)
		//check
		if targetBucket != "" || logPrefix != "" {
			return fmt.Errorf("log_target_bucket and log_prefix should set null when log_enable is false")
		}
		// set disabled, put empty request
		request := &s3.PutBucketLoggingInput{
			Bucket:              aws.String(bucket),
			BucketLoggingStatus: &s3.BucketLoggingStatus{},
		}

		resp, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosClientNew(cdcId).PutBucketLogging(request)
		if err != nil {
			return fmt.Errorf("cos disable log error: %s, bucket: %s", err.Error(), bucket)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], resp[%s]\n",
			logId, "cos enable log success", request.String(), resp.String())
	}

	return nil
}

func resourceTencentCloudCosBucketOriginACLBodyUpdate(ctx context.Context, service CosService, d *schema.ResourceData) error {
	logId := tccommon.GetLogId(ctx)
	aclHeader := ""
	aclBody := ""
	body, bodyOk := d.GetOk("acl_body")
	header, headerOk := d.GetOk("acl")
	bucket := d.Get("bucket").(string)
	cdcId := d.Get("cdc_id").(string)
	// If ACLXML update to empty, this will pass default header to delete verbose acl info
	if bodyOk {
		aclBody = body.(string)
	} else if headerOk {
		aclHeader = header.(string)
	} else {
		aclHeader = "private"
	}

	aclBodyOrderly, err := service.transACLBodyOrderly(ctx, aclBody)
	if err != nil {
		return fmt.Errorf("transfer ACL Body failed, reason:%v", err.Error())
	}

	log.Printf("[DEBUG]%s transACLBodyOrderly success, before:[\n%s\n], after:[\n%s\n]\n", logId, aclBody, aclBodyOrderly)

	if err = service.TencentCosPutBucketACLBody(ctx, bucket, aclBodyOrderly, aclHeader, cdcId); err != nil {
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, bucket:[%s]\n", logId, "put bucket acl body", bucket)

	return nil
}

func resourceTencentCloudCosBucketOriginPullUpdate(ctx context.Context, service CosService, d *schema.ResourceData) error {
	var rules []cos.BucketOriginRule
	v, ok := d.GetOk("origin_pull_rules")
	bucket := d.Get("bucket").(string)
	cdcId := d.Get("cdc_id").(string)
	if !ok {
		if err := service.DeleteBucketPullOrigin(ctx, bucket, cdcId); err != nil {
			return err
		}
		return nil
	}
	rulesRaw := v.([]interface{})
	for _, i := range rulesRaw {
		var (
			dMap = i.(map[string]interface{})
			item = &cos.BucketOriginRule{
				OriginCondition: &cos.BucketOriginCondition{
					HTTPStatusCode: "404",
				},
				OriginParameter: &cos.BucketOriginParameter{
					CopyOriginData: helper.Bool(true),
					HttpHeader:     &cos.BucketOriginHttpHeader{},
				},
				OriginInfo: &cos.BucketOriginInfo{
					FileInfo: &cos.BucketOriginFileInfo{
						PrefixDirective: false,
					},
				},
			}
		)

		// has deprecated
		if v, ok := dMap["sync_back_to_source"]; ok {
			if v.(bool) {
				item.OriginType = "Mirror"
			} else {
				item.OriginType = "Proxy"
			}
		}

		if v, ok := dMap["back_to_source_mode"].(string); ok && v != "" {
			item.OriginType = v
		}

		if v, ok := dMap["http_redirect_code"].(string); ok && v != "" {
			item.OriginParameter.HttpRedirectCode = v
		}

		if v, ok := dMap["priority"]; ok {
			item.RulePriority = v.(int)
		}
		if v, ok := dMap["prefix"]; ok {
			item.OriginCondition.Prefix = v.(string)
		}
		if v, ok := dMap["protocol"]; ok {
			item.OriginParameter.Protocol = v.(string)
		}
		if v, ok := dMap["host"]; ok {
			tmpHost := cos.BucketOriginHostInfo{}
			tmpHost.HostName = v.(string)
			item.OriginInfo.HostInfo = &tmpHost
		}
		if v, ok := dMap["follow_query_string"]; ok {
			item.OriginParameter.FollowQueryString = helper.Bool(v.(bool))
		}
		if v, ok := dMap["follow_redirection"]; ok {
			item.OriginParameter.FollowRedirection = helper.Bool(v.(bool))
		}
		//if v, ok := dMap["copy_origin_data"]; ok {
		//	item.OriginParameter.CopyOriginData = v.(bool)
		//}
		if v, ok := dMap["redirect_prefix"]; ok {
			value := v.(string)
			if value != "" {
				item.OriginInfo.FileInfo.PrefixDirective = true
			}
			item.OriginInfo.FileInfo.Prefix = value
		}
		if v, ok := dMap["redirect_suffix"]; ok {
			value := v.(string)
			if value != "" {
				item.OriginInfo.FileInfo.PrefixDirective = true
			}
			item.OriginInfo.FileInfo.Suffix = value
		}
		if v, ok := dMap["custom_http_headers"]; ok {
			var customHeaders []cos.OriginHttpHeader
			for key, val := range v.(map[string]interface{}) {
				customHeaders = append(customHeaders, cos.OriginHttpHeader{
					Key:   key,
					Value: val.(string),
				})
			}
			item.OriginParameter.HttpHeader.NewHttpHeaders = customHeaders
		}
		if v, ok := dMap["follow_http_headers"]; ok {
			var followHeaders []cos.OriginHttpHeader
			for _, item := range v.(*schema.Set).List() {
				header := cos.OriginHttpHeader{
					Key:   item.(string),
					Value: "",
				}
				followHeaders = append(followHeaders, header)
			}
			item.OriginParameter.HttpHeader.FollowHttpHeaders = followHeaders
		}
		rules = append(rules, *item)
	}

	if err := service.PutBucketPullOrigin(ctx, bucket, rules, cdcId); err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudCosBucketOriginDomainUpdate(ctx context.Context, service CosService, d *schema.ResourceData) error {
	v, ok := d.GetOk("origin_domain_rules")
	bucket := d.Get("bucket").(string)
	cdcId := d.Get("cdc_id").(string)
	if !ok {
		if err := service.DeleteBucketOriginDomain(ctx, bucket, cdcId); err != nil {
			return err
		}
		return nil
	}
	rules := v.([]interface{})
	domainRules := make([]cos.BucketDomainRule, 0)

	for _, rule := range rules {
		dMap := rule.(map[string]interface{})
		item := cos.BucketDomainRule{}
		if name, ok := dMap["domain"]; ok {
			item.Name = name.(string)
		}
		if status, ok := dMap["status"]; ok {
			item.Status = status.(string)
		}
		if domainType, ok := dMap["type"]; ok {
			item.Type = domainType.(string)
		}
		domainRules = append(domainRules, item)
	}

	if err := service.PutBucketOriginDomain(ctx, bucket, domainRules, cdcId); err != nil {
		return err
	}
	return nil
}

func getBucketPutOptions(d *schema.ResourceData) (useCosService bool, options *cos.BucketPutOptions) {
	opt := &cos.BucketPutOptions{
		XCosACL:              d.Get("acl").(string),
		XCosGrantRead:        "",
		XCosGrantWrite:       "",
		XCosGrantReadACP:     "",
		XCosGrantWriteACP:    "",
		XCosGrantFullControl: "",
	}
	grants, hasGrantHeaders := d.GetOk("grant_headers")
	maz, hasMAZ := d.GetOk("multi_az")
	ofs, hasOFS := d.GetOk("chdfs_ofs")

	if !hasGrantHeaders && !hasMAZ && !hasOFS {
		return false, opt
	}

	if hasGrantHeaders {
		headers := grants.(map[string]interface{})
		if v, ok := headers["grant_read"]; ok {
			opt.XCosGrantRead = v.(string)
		}
		if v, ok := headers["grant_write"]; ok {
			opt.XCosGrantWrite = v.(string)
		}
		if v, ok := headers["grant_read_acp"]; ok {
			opt.XCosGrantReadACP = v.(string)
		}
		if v, ok := headers["grant_write_acp"]; ok {
			opt.XCosGrantWriteACP = v.(string)
		}
		if v, ok := headers["grant_full_control"]; ok {
			opt.XCosGrantFullControl = v.(string)
		}
	}

	configOpt := cos.CreateBucketConfiguration{}
	if hasMAZ {
		if maz.(bool) {
			configOpt.BucketAZConfig = "MAZ"
			opt.CreateBucketConfiguration = &configOpt
		}
	}

	if hasOFS {
		if ofs.(bool) {
			configOpt.BucketArchConfig = "OFS"
			opt.CreateBucketConfiguration = &configOpt
		}
	}

	return true, opt
}

func expirationHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["date"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	return helper.HashString(buf.String())
}

func nonCurrentExpirationHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["non_current_days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	return helper.HashString(buf.String())
}

func abortIncompleteMultipartUploadHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["days_after_initiation"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	return helper.HashString(buf.String())
}

func transitionHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["date"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	if v, ok := m["storage_class"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	return helper.HashString(buf.String())
}

func nonCurrentTransitionHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["non_current_days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	if v, ok := m["storage_class"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	return helper.HashString(buf.String())
}

func getBucketReplications(d *schema.ResourceData) (role string, rules []cos.BucketReplicationRule, err error) {
	role = d.Get("replica_role").(string)
	replicaRules := d.Get("replica_rules").([]interface{})
	for i := range replicaRules {
		item := replicaRules[i].(map[string]interface{})
		rule := cos.BucketReplicationRule{
			Status: item["status"].(string),
			Destination: &cos.ReplicationDestination{
				Bucket: item["destination_bucket"].(string),
			},
		}
		if v, ok := item["priority"]; ok {
			rule.Priority = v.(int)
		}
		if v, ok := item["prefix"].(string); ok {
			rule.Prefix = v
		}
		if v, ok := item["filter"]; ok {
			for _, filterItem := range v.([]interface{}) {
				filterMap := filterItem.(map[string]interface{})
				filter := &cos.ReplicationFilter{}
				if v, ok := filterMap["and"]; ok {
					for _, andItem := range v.([]interface{}) {
						andMap := andItem.(map[string]interface{})
						and := &cos.ReplicationFilterAnd{}
						if v, ok := andMap["prefix"].(string); ok && v != "" {
							and.Prefix = v
						}

						if v, ok := andMap["tag"]; ok {
							for _, tag := range v.([]interface{}) {
								tagMap := tag.(map[string]interface{})
								tmpTag := cos.ObjectTaggingTag{}
								if v, ok := tagMap["key"].(string); ok && v != "" {
									tmpTag.Key = v
								}

								if v, ok := tagMap["value"].(string); ok && v != "" {
									tmpTag.Value = v
								}

								and.Tag = append(and.Tag, tmpTag)
							}
						}

						filter.And = and
					}
				}

				if v, ok := filterMap["prefix"].(string); ok && v != "" {
					filter.Prefix = v
				}

				rule.Filter = filter
			}
		}
		if v, ok := item["id"].(string); ok {
			rule.ID = v
		}
		if v, ok := item["destination_storage_class"].(string); ok && v != "" {
			rule.Destination.StorageClass = v
		}
		if v, ok := item["destination_encryption_kms_key_id"].(string); ok && v != "" {
			rule.Destination.EncryptionConfiguration = &cos.ReplicationEncryptionConfiguration{
				ReplicaKmsKeyID: v,
			}
		}
		if v, ok := item["delete_marker_replication"]; ok {
			for _, dmrConfig := range v.([]interface{}) {
				dMap := dmrConfig.(map[string]interface{})
				if vv, ok := dMap["status"].(string); ok {
					rule.DeleteMarkerReplication = &cos.DeleteMarkerReplication{
						Status: vv,
					}
				}
			}
		}
		if v, ok := item["source_selection_criteria"]; ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				if v, ok := dMap["sse_kms_encrypted_objects"]; ok {
					for _, item := range v.([]interface{}) {
						dMap := item.(map[string]interface{})
						rule.SourceSelectionCriteria = &cos.SourceSelectionCriteria{
							SseKmsEncryptedObjects: &cos.SseKmsEncryptedObjects{
								Status: dMap["status"].(string),
							},
						}
					}
				}
			}
		}
		rules = append(rules, rule)
	}
	return
}

func setBucketReplication(d *schema.ResourceData, result cos.GetBucketReplicationResult) (err error) {
	if result.Role != "" {
		_ = d.Set("replica_role", result.Role)
	}
	rules := make([]map[string]interface{}, 0)
	if len(result.Rule) > 0 {
		for i := range result.Rule {
			item := result.Rule[i]
			rule := map[string]interface{}{
				"status":                    item.Status,
				"destination_bucket":        item.Destination.Bucket,
				"destination_storage_class": item.Destination.StorageClass,
			}
			if item.ID != "" {
				rule["id"] = item.ID
			}
			if item.Priority != 0 {
				rule["priority"] = item.Priority
			}
			if item.Prefix != "" {
				rule["prefix"] = item.Prefix
			}
			if item.Filter != nil {
				filter := make([]map[string]interface{}, 0)
				filterMap := map[string]interface{}{}
				if item.Filter.And != nil {
					and := make([]map[string]interface{}, 0)
					andMap := map[string]interface{}{}
					if item.Filter.And.Prefix != "" {
						andMap["prefix"] = item.Filter.And.Prefix
					}
					if len(item.Filter.And.Tag) > 0 {
						tags := make([]map[string]interface{}, 0)
						for i := range item.Filter.And.Tag {
							tag := item.Filter.And.Tag[i]
							tagMap := map[string]interface{}{
								"key":   tag.Key,
								"value": tag.Value,
							}
							tags = append(tags, tagMap)
						}
						andMap["tag"] = tags
					}
					and = append(and, andMap)
					filterMap["and"] = and
				}
				if item.Filter.Prefix != "" {
					filterMap["prefix"] = item.Filter.Prefix
				}
				filter = append(filter, filterMap)
				rule["filter"] = filter
			}
			if item.Destination.EncryptionConfiguration != nil {
				rule["destination_encryption_kms_key_id"] = item.Destination.EncryptionConfiguration.ReplicaKmsKeyID
			}
			//
			var deleteMarkerReplicationMap = map[string]interface{}{
				"status": "Enabled",
			}
			if item.DeleteMarkerReplication != nil {
				if item.DeleteMarkerReplication.Status != "" {
					deleteMarkerReplicationMap["status"] = item.DeleteMarkerReplication.Status
				}
			}
			rule["delete_marker_replication"] = []interface{}{deleteMarkerReplicationMap}
			if item.SourceSelectionCriteria != nil {
				sourceSelectionCriteriaMap := map[string]interface{}{}
				if item.SourceSelectionCriteria.SseKmsEncryptedObjects != nil {
					sseKmsEncryptedObjectsMap := map[string]interface{}{
						"status": item.SourceSelectionCriteria.SseKmsEncryptedObjects.Status,
					}

					sourceSelectionCriteriaMap["sse_kms_encrypted_objects"] = []interface{}{sseKmsEncryptedObjectsMap}
				}

				rule["source_selection_criteria"] = []interface{}{sourceSelectionCriteriaMap}
			}
			rules = append(rules, rule)
		}
	}
	err = d.Set("replica_rules", rules)
	return
}

func ACLBodyDiffFunc(olds, news string, d *schema.ResourceData) (result bool) {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket.ACLBodyDiffFunc")()
	log.Printf("[DEBUG] ACLBodyDiffFunc called, before:[\n%s\n], after:[\n%s\n]\n", olds, news)

	oldDoc := etree.NewDocument()
	newDoc := etree.NewDocument()

	if err := oldDoc.ReadFromString(olds); err != nil {
		log.Printf("[CRITAL]read old xml from string error: %v", err)
		return false
	}

	if err := newDoc.ReadFromString(news); err != nil {
		log.Printf("[CRITAL]read new xml from string error: %v", err)
		return false
	}

	oldRoot := oldDoc.SelectElement("AccessControlPolicy")
	newRoot := newDoc.SelectElement("AccessControlPolicy")

	if oldRoot == nil || newRoot == nil {
		log.Println("[CRITAL]oldRoot or newRoot is nil: return false.")
		return false
	}

	oldOwner := oldRoot.SelectElement("Owner")
	newOwner := newRoot.SelectElement("Owner")

	if oldOwner == nil || newOwner == nil {
		log.Println("[CRITAL]oldOwner or newOwner is nil: return false.")
		return false
	}

	oldOwnerId := oldOwner.SelectElement("ID")
	oldOwnerName := oldOwner.SelectElement("DisplayName")
	newOwnerId := newOwner.SelectElement("ID")
	newOwnerName := newOwner.SelectElement("DisplayName")

	if oldOwnerId == nil || newOwnerId == nil {
		log.Println("[CRITAL]oldOwnerId or newOwnerId is nil: return false.")
		return false
	}

	// diff: Owner element
	if oldOwnerId.Text() != newOwnerId.Text() {
		log.Printf("[CRITAL]OwnerId[old:%s, new:%s] not equal: return false.\n", oldOwnerId.Text(), newOwnerId.Text())
		return false
	}

	// diff check: owner display name(if have)
	if oldOwnerName != nil {
		if newOwnerName == nil {
			log.Println("[CRITAL]newOwnerName is nil: return false.")
			return false
		}
		if oldOwnerName.Text() != newOwnerName.Text() {
			log.Printf("[CRITAL]OwnerName[old:%s, new:%s] not equal: return false.\n", oldOwnerName.Text(), newOwnerName.Text())
			return false
		}
	}

	// diff: ACL element
	oldGrantees := oldRoot.FindElements("//Grantee")
	newGrantees := newRoot.FindElements("//Grantee")
	// check count
	if len(oldGrantees) != len(newGrantees) {
		return false
	}
	// check content
	for _, oldGrantee := range oldGrantees {
		for _, attr := range oldGrantee.Attr {
			if attr.Key != "type" {
				// only need to handle the type attribute
				continue
			}
			// anonymous or real user
			oldGranteeType := attr.Value

			oldGranteeID := oldGrantee.SelectElement("ID")
			oldGranteeURI := oldGrantee.SelectElement("URI")
			oldGranteeDisplayName := oldGrantee.SelectElement("DisplayName")
			oldGrant := oldGrantee.Parent()
			oldGrantPermission := oldGrant.SelectElement("Permission")

			// find the new grant permission by specified grantee type
			result = false
			for _, newGrantee := range newRoot.FindElements(fmt.Sprintf("//Grantee[@type='%s']", oldGranteeType)) {
				newGranteeID := newGrantee.SelectElement("ID")
				newGranteeURI := newGrantee.SelectElement("URI")
				newGranteeDisplayName := newGrantee.SelectElement("DisplayName")
				newGrant := newGrantee.Parent()
				newGrantPermission := newGrant.SelectElement("Permission")

				// diff check: grantee id and name for real user
				if oldGranteeType == COS_ACL_GRANTEE_TYPE_USER {
					if oldGranteeID == nil || newGranteeID == nil {
						continue
					}
					if oldGranteeID.Text() != newGranteeID.Text() {
						continue
					}

					// diff check: grantee display name(if have)
					if oldGranteeDisplayName != nil {
						if newGranteeDisplayName == nil {
							continue
						}
						if oldGranteeDisplayName.Text() != newGranteeDisplayName.Text() {
							continue
						}
					}
				}

				// diff check: grantee uri for anonymous
				if oldGranteeType == COS_ACL_GRANTEE_TYPE_ANONYMOUS {
					if oldGranteeURI == nil || newGranteeURI == nil {
						continue
					}
					if oldGranteeURI.Text() != newGranteeURI.Text() {
						continue
					}
				}

				// diff check: permission
				if oldGrantPermission == nil || newGrantPermission == nil {
					continue
				}
				if oldGrantPermission.Text() != newGrantPermission.Text() {
					continue
				}

				// congrats! passed all diff checks for this grant.
				result = true

				var uid string
				if oldGranteeType == COS_ACL_GRANTEE_TYPE_USER {
					uid = oldGranteeID.Text()
				} else {
					uid = oldGranteeURI.Text()
				}
				log.Printf("[DEBUG] diff verification passed for grantee:[%s:%s]\n", oldGranteeType, uid)
			}
			if !result {
				return false
			}
		}
	}
	log.Printf("[DEBUG] Owner:%s's final equation result between old and new ACL is:[%v]\n", oldOwnerId.Text(), result)
	return result
}
