package cos

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCosBatch() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudCosBatchRead,
		Create: resourceTencentCloudCosBatchCreate,
		Update: resourceTencentCloudCosBatchUpdate,
		Delete: resourceTencentCloudCosBatchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"uin": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Uin。",
			},
			"appid": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Appid。",
			},
			"confirmation_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "是否confirm before performing 任务. 默认为 false。",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Mission 描述 如果 您 已配置 此 信息 当 您 创建 任务， 内容 是 返回. 描述 长度 ranges 从 0 到 256 bytes。",
			},
			"priority": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Mission 优先级 higher 值， higher 优先级 的 任务. 优先级 值 范围 从 0 到 2147483647。",
			},
			"role_arn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "COS 资源 identifier，其中 是 用于identify 角色 您 创建. You need 此 资源 identifier 到 verify your identity。",
			},
			"manifest": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "列表 objects 到 是 processed。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "location 信息 的 列表 objects。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"etag": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "指定etag 的 对象 列表. Length 1-1024 bytes。",
									},
									"object_arn": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "指定unique 资源 identifier 的 对象 manifest，其中 是 1-1024 bytes long。",
									},
									"object_version_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "指定version 的 对象 manifest ID，其中 是 1-1024 bytes long。",
									},
								},
							},
						},
						"spec": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "格式 信息 该 describes 列表 objects. 如果 它 是 CSV 文件，此 element describes 字段 contained 在 manifest。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fields": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Describes 字段 contained 在 listing，其中 您 need 到 使用 到 指定CSV 文件 字段 当 格式 是 COSBatchOperations_CSV_V1. Legal 字段 是: Ignore，存储桶，键，VersionId。",
									},
									"format": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "指定format 信息 对于 列表 objects. Legal 字段 是: COSBatchOperations_CSV_V1，COSInventoryReport_CSV_V1。",
									},
								},
							},
						},
					},
				},
			},
			"operation": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Select 操作 到 是 performed 在 objects 在 manifest 文件。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cos_put_object_copy": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "指定specific 参数 对于 batch copy operation 在 objects 在 列表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_control_directive": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "此 element specifies how ACL 是 copied. 有效 值:\n" +
											"- Copy: inherits the source object ACL\n" +
											"- Replaced: replace source ACL\n" +
											"- Add: add a new ACL based on the source ACL.",
									},
									"access_control_grants": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Controls 特定 访问 到 对象。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"display_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "用户 名称",
												},
												"identifier": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "用户 ID (UIN) 在 qcs 格式 For 示例: qcs::cam::uin/100000000001:uin/100000000001。",
												},
												"type_identifier": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "指定type 的 Identifier. Currently，仅 用户 ID 是 支持. Enumerated 值: ID。",
												},
												"permission": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "指定a 权限 到 是 granted. Enumerated 值: READ,WRITE,FULL_CONTROL。",
												},
											},
										},
									},
									"canned_access_control_list": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Defines ACL 属性 的 对象. 有效值：私有，公有-read。",
									},
									"prefix_replace": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "指定是否prefix 的 来源 对象 needs 到 是 replaced. A 值 的 true 表示replacement 对象 prefix，其中 needs 到 是 使用 使用 <ResourcesPrefix> 和 <TargetKeyPrefix>. 默认值：false。",
									},
									"resources_prefix": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "此 字段 是 有效 仅 当 < PrefixReplace > 值 是 true. 指定source 对象 prefix 到 是 replaced，和 replacement directory should end 使用 `/`. Can 是 空 使用 最大长度1024 bytes。",
									},
									"target_key_prefix": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "此 字段 是 有效 仅 当 <PrefixReplace> 值 是 true. 此 值 表示 replaced prefix，和 replacement directory should end 使用 /. Can 是 空 使用 最大长度1024 bytes。",
									},
									"modified_since_constraint": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "当 对象 是 modified after 指定 时间， operation 是 performed，otherwise 412 是 返回。",
									},
									"unmodified_since_constraint": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "当 对象 has 不 been modified after 指定 时间， operation 是 performed，otherwise 412 是 返回。",
									},
									"metadata_directive": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "此 element 指定是否copy 对象 metadata 从 来源 对象 或 replace 它 使用 metadata 在 < NewObjectMetadata > element. 有效值：Copy，Replaced，Add. Copy: inherit 来源 对象 metadata; Replaced: replace 来源 metadata; Add: add new metadata based 在 来源 metadata。",
									},
									"new_object_metadata": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Configure metadata 对于 对象。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cache_control": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "caching instructions defined 在 RFC 2616 是 saved 作为 对象 metadata。",
												},
												"content_disposition": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "文件 名称 defined 在 RFC 2616 是 saved 作为 对象 metadata。",
												},
												"content_encoding": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "编码 格式 defined 在 RFC 2616 是 saved 作为 对象 metadata。",
												},
												"content_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "内容 types defined 在 RFC 2616 是 saved 作为 对象 metadata。",
												},
												"http_expires_date": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "缓存 过期时间 defined 在 RFC 2616 是 saved 作为 对象 metadata。",
												},
												"sse_algorithm": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Server 加密 algorithm. Currently，仅 AES256 是 支持。",
												},
												"user_metadata": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Includes 用户-defined metadata。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "键",
															},
															"value": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "值",
															},
														},
													},
												},
											},
										},
									},
									"tagging_directive": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "此 element 指定是否copy 对象 标签 从 来源 对象 或 replace 它 使用 标签 在 < NewObjectTagging > element. 有效值：Copy，Replaced，Add. Copy: inherits 来源 对象 标签; Replaced: replaces 来源 标签; Add: adds new 标签 based 在 来源 标签",
									},
									"new_object_tagging": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "标签 的 配置 对象，其中 必须 是 指定 当 < TaggingDirective > 值 是 Replace 或 Add。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "键",
												},
												"value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "值",
												},
											},
										},
									},
									"storage_class": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Sets 存储 级别 的 对象. Enumerated 值: STANDARD,STANDARD_IA. 默认值：STANDARD。",
									},
									"target_resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Sets 目标 存储桶 对于 Copy. Use qcs 到 specify，对于 示例，qcs::cos:ap-chengdu:uid/1250000000:examplebucket-1250000000。",
									},
								},
							},
						},
						"cos_initiate_restore_object": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "指定specific 参数 对于 batch 恢复 operation 对于 archive 存储 类型 objects 在 inventory。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"expiration_in_days": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Sets 数量 days after 其中 copy 将 是 automatically expired 和 删除， 整数 在 范围 的 1-365。",
									},
									"job_tier": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Select archive recovery model. 可用值：Bulk，Standard。",
									},
								},
							},
						},
					},
				},
			},
			"report": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "任务 completion 报告。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Delivery 存储桶 对于 任务 completion reports。",
						},
						"enabled": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "是否output 任务 completion 报告。",
						},
						"format": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "任务 completion 报告 格式 信息. Legal 值: Report_CSV_V1。",
						},
						"prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Prefix 信息 对于 任务 completion 报告. Length 0-256 bytes。",
						},
						"report_scope": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "任务 completion 报告 任务 信息 该 needs 到 是 recorded 到 determine 是否record execution 信息 的 all operations 或 信息 的 failed operations. Legal 值: AllTasks，FailedTasksOnly。",
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Current 状态 的 任务.\n" +
					"Legal parameter values include Active, Cancelled, Cancelling, Complete, Completing, Failed, Failing, New, Paused, Pausing, Preparing, Ready, Suspended.\n" +
					"For Update status, when you move a task to the Ready state, COS will assume that you have confirmed the task and will perform it. When you move a task to the Cancelled state, COS cancels the task. Optional parameters include: Ready, Cancelled.",
			},
			"job_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job ID。",
			},
		},
	}
}

func resourceTencentCloudCosBatchRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_batch.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	uin := idSplit[0]
	appid, _ := strconv.Atoi(idSplit[1])
	jobId := idSplit[2]
	headers := &cos.BatchRequestHeaders{
		XCosAppid: appid,
	}

	result, response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosBatchClient(uin).Batch.DescribeJob(ctx, jobId, headers)
	responseBody, _ := json.Marshal(response.Body)
	if err != nil {
		log.Printf("[DEBUG]%s api[DescribeJob] success, request body [%s], response body [%s], err: [%s]\n", logId, jobId, responseBody, err.Error())
		return err
	}
	if result == nil || result.Job == nil {
		return fmt.Errorf("DescribeJob response is nil!")
	}
	confirmationRequired, err := strconv.ParseBool(result.Job.ConfirmationRequired)
	if err != nil {
		return err
	}
	_ = d.Set("uin", uin)
	_ = d.Set("appid", appid)
	_ = d.Set("job_id", jobId)
	_ = d.Set("confirmation_required", confirmationRequired)
	_ = d.Set("description", result.Job.Description)
	_ = d.Set("priority", result.Job.Priority)
	_ = d.Set("role_arn", result.Job.RoleArn)
	manifestResult := make(map[string]interface{})
	locationResult := make(map[string]interface{})
	specResult := make(map[string]interface{})
	manifest := result.Job.Manifest
	location := manifest.Location
	spec := manifest.Spec
	locationResult["etag"] = location.ETag
	locationResult["object_arn"] = location.ObjectArn
	locationResult["object_version_id"] = location.ObjectVersionId
	manifestResult["location"] = []interface{}{locationResult}
	specResult["fields"] = spec.Fields
	specResult["format"] = spec.Format
	manifestResult["spec"] = []interface{}{specResult}
	_ = d.Set("manifest", []interface{}{manifestResult})

	operationResult := make(map[string]interface{})
	if result.Job.Operation.PutObjectCopy != nil {
		putObjectCopyResult := make(map[string]interface{})
		putObjectCopy := result.Job.Operation.PutObjectCopy
		putObjectCopyResult["access_control_directive"] = putObjectCopy.AccessControlDirective
		accessControlGrants := putObjectCopy.AccessControlGrants
		if accessControlGrants != nil {
			accessControlGrantsResult := make(map[string]interface{})
			accessControlGrantsResult["display_name"] = accessControlGrants.COSGrants.Grantee.DisplayName
			accessControlGrantsResult["identifier"] = accessControlGrants.COSGrants.Grantee.Identifier
			accessControlGrantsResult["type_identifier"] = accessControlGrants.COSGrants.Grantee.TypeIdentifier
			accessControlGrantsResult["permission"] = accessControlGrants.COSGrants.Permission
			putObjectCopyResult["access_control_grants"] = []interface{}{accessControlGrantsResult}

		}

		putObjectCopyResult["canned_access_control_list"] = putObjectCopy.CannedAccessControlList
		putObjectCopyResult["prefix_replace"] = putObjectCopy.PrefixReplace
		putObjectCopyResult["resources_prefix"] = putObjectCopy.ResourcesPrefix
		putObjectCopyResult["target_key_prefix"] = putObjectCopy.TargetKeyPrefix
		putObjectCopyResult["modified_since_constraint"] = putObjectCopy.ModifiedSinceConstraint
		putObjectCopyResult["unmodified_since_constraint"] = putObjectCopy.UnModifiedSinceConstraint
		putObjectCopyResult["metadata_directive"] = putObjectCopy.MetadataDirective

		newObjectMetadata := putObjectCopy.NewObjectMetadata
		if newObjectMetadata != nil {
			newObjectMetadataResult := make(map[string]interface{})
			newObjectMetadataResult["cache_control"] = newObjectMetadata.CacheControl
			newObjectMetadataResult["content_disposition"] = newObjectMetadata.ContentDisposition
			newObjectMetadataResult["content_encoding"] = newObjectMetadata.ContentEncoding
			newObjectMetadataResult["content_type"] = newObjectMetadata.ContentType
			newObjectMetadataResult["http_expires_date"] = newObjectMetadata.HttpExpiresDate
			newObjectMetadataResult["sse_algorithm"] = newObjectMetadata.SSEAlgorithm
			userMetadataResult := make([]interface{}, 0)
			userMetadata := newObjectMetadata.UserMetadata
			for _, item := range userMetadata {
				userMetadataResult = append(userMetadataResult, map[string]interface{}{
					"key":   item.Key,
					"value": item.Value,
				})
			}
			newObjectMetadataResult["user_metadata"] = userMetadataResult
			putObjectCopyResult["new_object_metadata"] = []interface{}{newObjectMetadataResult}
		}

		putObjectCopyResult["tagging_directive"] = putObjectCopy.TaggingDirective
		if putObjectCopy.NewObjectTagging != nil {
			cosTagResult := make([]interface{}, 0)
			for _, item := range putObjectCopy.NewObjectTagging.COSTag {
				cosTagResult = append(cosTagResult, map[string]interface{}{
					"key":   item.Key,
					"value": item.Value,
				})
			}
			putObjectCopyResult["new_object_tagging"] = cosTagResult
		}

		putObjectCopyResult["storage_class"] = putObjectCopy.StorageClass
		putObjectCopyResult["target_resource"] = putObjectCopy.TargetResource

		operationResult["cos_put_object_copy"] = []interface{}{putObjectCopyResult}
	}
	if result.Job.Operation.RestoreObject != nil {
		restoreObjectResult := make(map[string]interface{})
		restoreObject := result.Job.Operation.RestoreObject
		restoreObjectResult["expiration_in_days"] = restoreObject.ExpirationInDays
		restoreObjectResult["job_tier"] = restoreObject.JobTier
		operationResult["cos_initiate_restore_object"] = []interface{}{restoreObjectResult}
	}

	_ = d.Set("operation", []interface{}{operationResult})
	if result.Job.Report != nil {
		report := result.Job.Report
		reportResult := make(map[string]interface{})
		reportResult["bucket"] = report.Bucket
		reportResult["enabled"] = report.Enabled
		reportResult["format"] = report.Format
		reportResult["prefix"] = report.Prefix
		reportResult["report_scope"] = report.ReportScope
		_ = d.Set("report", []interface{}{reportResult})
	}

	_ = d.Set("status", result.Job.Status)
	return nil
}

func resourceTencentCloudCosBatchCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_batch.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	uin := d.Get("uin").(string)
	opt := &cos.BatchCreateJobOptions{
		ClientRequestToken: uuid.New().String(),
		Priority:           d.Get("priority").(int),
		RoleArn:            d.Get("role_arn").(string),
	}
	if v, ok := d.GetOk("confirmation_required"); ok && v.(bool) {
		opt.ConfirmationRequired = "true"
	} else {
		opt.ConfirmationRequired = "false"
	}

	if v, ok := d.GetOk("description"); ok {
		opt.Description = v.(string)
	}

	if manifestMap, ok := helper.InterfacesHeadMap(d, "manifest"); ok {
		batchJobManifest := &cos.BatchJobManifest{}
		locationMap := manifestMap["location"].([]interface{})[0].(map[string]interface{})
		specMap := manifestMap["spec"].([]interface{})[0].(map[string]interface{})
		location := &cos.BatchJobManifestLocation{
			ETag:      locationMap["etag"].(string),
			ObjectArn: locationMap["object_arn"].(string),
		}
		if v, ok := locationMap["object_version_id"]; ok {
			location.ObjectVersionId = v.(string)
		}
		batchJobManifest.Location = location
		spec := &cos.BatchJobManifestSpec{
			Format: specMap["format"].(string),
		}
		batchJobManifest.Spec = spec
		if v, ok := specMap["fields"]; ok {
			fields := make([]string, 0)
			for _, item := range v.([]interface{}) {
				fields = append(fields, item.(string))
			}
			spec.Fields = fields
		}
		opt.Manifest = batchJobManifest
	}

	if operationMap, ok := helper.InterfacesHeadMap(d, "operation"); ok {
		operation := &cos.BatchJobOperation{}
		if v, ok := operationMap["cos_put_object_copy"]; ok {
			cosPutObjectCopy := v.([]interface{})[0].(map[string]interface{})
			putObjectCopy := &cos.BatchJobOperationCopy{}

			if v, ok := cosPutObjectCopy["access_control_directive"]; ok {
				putObjectCopy.AccessControlDirective = v.(string)
			}
			if v, ok := cosPutObjectCopy["access_control_grants"]; ok && len(v.([]interface{})) > 0 {
				accessControlGrantMap := v.([]interface{})[0].(map[string]interface{})
				grantee := &cos.BatchGrantee{}
				grant := &cos.BatchCOSGrant{}
				if v, ok := accessControlGrantMap["display_name"]; ok {
					grantee.DisplayName = v.(string)
				}
				if v, ok := accessControlGrantMap["identifier"]; ok {
					grantee.Identifier = v.(string)
				}
				if v, ok := accessControlGrantMap["type_identifier"]; ok {
					grantee.TypeIdentifier = v.(string)
				}
				grant.Grantee = grantee
				if v, ok := accessControlGrantMap["permission"]; ok {
					grant.Permission = v.(string)
				}
				putObjectCopy.AccessControlGrants = &cos.BatchAccessControlGrants{
					COSGrants: grant,
				}
			}
			if v, ok := cosPutObjectCopy["canned_access_control_list"]; ok {
				putObjectCopy.CannedAccessControlList = v.(string)
			}
			if v, ok := cosPutObjectCopy["prefix_replace"]; ok {
				putObjectCopy.PrefixReplace = v.(bool)
			}
			if v, ok := cosPutObjectCopy["resources_prefix"]; ok {
				putObjectCopy.ResourcesPrefix = v.(string)
			}
			if v, ok := cosPutObjectCopy["target_key_prefix"]; ok {
				putObjectCopy.TargetKeyPrefix = v.(string)
			}
			if v, ok := cosPutObjectCopy["metadata_directive"]; ok {
				putObjectCopy.MetadataDirective = v.(string)
			}
			if v, ok := cosPutObjectCopy["modified_since_constraint"]; ok {
				putObjectCopy.ModifiedSinceConstraint = int64(v.(int))
			}
			if v, ok := cosPutObjectCopy["unmodified_since_constraint"]; ok {
				putObjectCopy.UnModifiedSinceConstraint = int64(v.(int))
			}

			if v, ok := cosPutObjectCopy["new_object_metadata"]; ok && len(v.([]interface{})) > 0 {
				newObjectMetadataMap := v.([]interface{})[0].(map[string]interface{})
				newObjectMetadata := &cos.BatchNewObjectMetadata{}
				if v, ok := newObjectMetadataMap["cache_control"]; ok {
					newObjectMetadata.CacheControl = v.(string)
				}
				if v, ok := newObjectMetadataMap["content_disposition"]; ok {
					newObjectMetadata.ContentDisposition = v.(string)
				}
				if v, ok := newObjectMetadataMap["content_encoding"]; ok {
					newObjectMetadata.ContentEncoding = v.(string)
				}
				if v, ok := newObjectMetadataMap["content_type"]; ok {
					newObjectMetadata.ContentType = v.(string)
				}
				if v, ok := newObjectMetadataMap["http_expires_date"]; ok {
					newObjectMetadata.HttpExpiresDate = v.(string)
				}
				if v, ok := newObjectMetadataMap["sse_algorithm"]; ok {
					newObjectMetadata.SSEAlgorithm = v.(string)
				}
				if v, ok := newObjectMetadataMap["user_metadata"]; ok {
					newObjectMetadata.UserMetadata = make([]cos.BatchMetadata, 0)
					for _, userMetadataItem := range v.([]interface{}) {
						userMetadataItemMap := userMetadataItem.(map[string]interface{})
						batchMetadata := cos.BatchMetadata{
							Key:   userMetadataItemMap["key"].(string),
							Value: userMetadataItemMap["value"].(string),
						}
						newObjectMetadata.UserMetadata = append(newObjectMetadata.UserMetadata, batchMetadata)
					}
				}
				putObjectCopy.NewObjectMetadata = newObjectMetadata
			}

			if v, ok := cosPutObjectCopy["tagging_directive"]; ok {
				putObjectCopy.TaggingDirective = v.(string)
			}
			if v, ok := cosPutObjectCopy["new_object_tagging"]; ok {
				newObjectTaggings := v.([]interface{})
				cosTags := make([]cos.BatchCOSTag, 0)
				for _, item := range newObjectTaggings {
					tag := item.(map[string]interface{})
					cosTags = append(cosTags, cos.BatchCOSTag{
						Key:   tag["key"].(string),
						Value: tag["value"].(string),
					})
				}
				putObjectCopy.NewObjectTagging = &cos.BatchNewObjectTagging{COSTag: cosTags}
			}
			if v, ok := cosPutObjectCopy["storage_class"]; ok {
				putObjectCopy.StorageClass = v.(string)
			}
			if v, ok := cosPutObjectCopy["target_resource"]; ok {
				putObjectCopy.TargetResource = v.(string)
			}

			operation.PutObjectCopy = putObjectCopy

		}

		if v, ok := operationMap["cos_initiate_restore_object"]; ok && len(v.([]interface{})) > 0 {
			restoreObject := &cos.BatchInitiateRestoreObject{}
			cosInitiateRestoreObject := v.([]interface{})[0].(map[string]interface{})
			if v, ok := cosInitiateRestoreObject["expiration_in_days"]; ok {
				restoreObject.ExpirationInDays = v.(int)
			}
			if v, ok := cosInitiateRestoreObject["job_tier"]; ok {
				restoreObject.JobTier = v.(string)
			}
			operation.RestoreObject = restoreObject
		}
		opt.Operation = operation
	}

	if reportMap, ok := helper.InterfacesHeadMap(d, "report"); ok {
		batchJobReport := &cos.BatchJobReport{}
		if v, ok := reportMap["bucket"]; ok {
			batchJobReport.Bucket = v.(string)
		}
		if v, ok := reportMap["enabled"]; ok {
			batchJobReport.Enabled = v.(string)
		}
		if v, ok := reportMap["format"]; ok {
			batchJobReport.Format = v.(string)
		}
		if v, ok := reportMap["prefix"]; ok {
			batchJobReport.Prefix = v.(string)
		}
		if v, ok := reportMap["report_scope"]; ok {
			batchJobReport.ReportScope = v.(string)
		}
		opt.Report = batchJobReport
	}
	appid := d.Get("appid").(int)
	headers := &cos.BatchRequestHeaders{
		XCosAppid: appid,
	}
	var batchCreateJobResult *cos.BatchCreateJobResult
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		req, _ := json.Marshal(opt)
		result, response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosBatchClient(uin).Batch.CreateJob(ctx, opt, headers)
		responseBody, _ := json.Marshal(response.Body)
		log.Printf("[DEBUG]%s api[CreateJob], request body [%s], response body [%s]\n", logId, req, responseBody)
		if e != nil {
			return tccommon.RetryError(e)
		}

		batchCreateJobResult = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create job failed, reason:%+v", logId, err)
		return err
	}
	if v, ok := d.GetOk("status"); ok {
		opt := &cos.BatchUpdateStatusOptions{
			JobId:              batchCreateJobResult.JobId,
			RequestedJobStatus: v.(string),
		}
		req, _ := json.Marshal(opt)
		_, response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosBatchClient(uin).Batch.UpdateJobStatus(ctx, opt, headers)
		responseBody, _ := json.Marshal(response.Body)
		if err != nil {
			log.Printf("[DEBUG]%s api[UpdateJobStatus] error, request body [%s], response body [%s], err: [%s]\n", logId, req, responseBody, err.Error())
			return err
		}
	}
	d.SetId(uin + tccommon.FILED_SP + strconv.Itoa(appid) + tccommon.FILED_SP + batchCreateJobResult.JobId)
	return resourceTencentCloudCosBatchRead(d, meta)
}

func resourceTencentCloudCosBatchUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_batch.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	uin := idSplit[0]
	appid, _ := strconv.Atoi(idSplit[1])
	jobId := idSplit[2]
	headers := &cos.BatchRequestHeaders{
		XCosAppid: appid,
	}
	if d.HasChange("priority") {
		opt := &cos.BatchUpdatePriorityOptions{
			JobId:    jobId,
			Priority: d.Get("priority").(int),
		}
		req, _ := json.Marshal(opt)
		_, response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosBatchClient(uin).Batch.UpdateJobPriority(ctx, opt, headers)
		responseBody, _ := json.Marshal(response.Body)
		if err != nil {
			log.Printf("[DEBUG]%s api[UpdateJobPriority] error, request body [%s], response body [%s], err: [%s]\n", logId, req, responseBody, err.Error())
			return err
		}
	}
	if d.HasChange("status") {
		opt := &cos.BatchUpdateStatusOptions{
			JobId:              jobId,
			RequestedJobStatus: d.Get("status").(string),
		}
		req, _ := json.Marshal(opt)
		_, response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosBatchClient(uin).Batch.UpdateJobStatus(ctx, opt, headers)
		responseBody, _ := json.Marshal(response.Body)
		if err != nil {
			log.Printf("[DEBUG]%s api[UpdateJobStatus] error, request body [%s], response body [%s], err: [%s]\n", logId, req, responseBody, err.Error())
			return err
		}
	}
	return resourceTencentCloudCosBatchRead(d, meta)
}

func resourceTencentCloudCosBatchDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_batch.delete")()
	defer tccommon.InconsistentCheck(d, meta)()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	uin := idSplit[0]
	appid, _ := strconv.Atoi(idSplit[1])
	jobId := idSplit[2]
	headers := &cos.BatchRequestHeaders{
		XCosAppid: appid,
	}
	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosBatchClient(uin).Batch.DeleteJob(ctx, jobId, headers)
	responseBody, _ := json.Marshal(response.Body)
	if err != nil {
		log.Printf("[DEBUG]%s api[DeleteJob] success, response body [%s], err: [%s]\n", logId, responseBody, err.Error())
		return err
	}
	return nil
}
