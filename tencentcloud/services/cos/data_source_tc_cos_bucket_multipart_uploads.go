package cos

import (
	"context"
	"encoding/json"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCosBucketMultipartUploads() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCosBucketMultipartUploadsRead,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "存储桶",
			},
			"delimiter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "delimiter 是 symbol，和 Object 名称 包含Object between 指定 prefix 和 first occurrence 的 delimiter 字符 作为 集合 的 elements: common prefix. 如果 there 是 无 prefix，start 从 beginning 的 路径",
			},
			"encoding_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "指定encoding 格式 的 返回值 Legal 值: URL",
			},
			"prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "返回 Object 键 必须 是 prefixed 使用 Prefix. 注意 该 当 使用 prefix 查询， 返回 键 still 包含Prefix。",
			},
			"uploads": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information 对于 each Upload。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 Object。",
						},
						"upload_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mark ID 此 multipart upload。",
						},
						"storage_class": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "用于represent 存储 级别 的 chunk. Enumerated 值: STANDARD,STANDARD_IA,ARCHIVE。",
						},
						"initiated": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "starting 时间 的 multipart upload。",
						},
						"owner": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Information 用于represent 所有者 的 these chunks。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "用户's 唯一 CAM identity ID。",
									},
									"display_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Abbreviation 对于 用户 identity ID (UIN)。",
									},
								},
							},
						},
						"initiator": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "用于represent 信息 的 initiator 的 此 upload。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "用户's 唯一 CAM identity ID。",
									},
									"display_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Abbreviation 对于 用户 identity ID (UIN)。",
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

func dataSourceTencentCloudCosBucketMultipartUploadsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cos_bucket_multipart_uploads.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	bucket := d.Get("bucket").(string)
	multipartUploads := make([]map[string]interface{}, 0)
	opt := &cos.ListMultipartUploadsOptions{}
	if v, ok := d.GetOk("delimiter"); ok {
		opt.Delimiter = v.(string)
	}
	if v, ok := d.GetOk("encoding_type"); ok {
		opt.EncodingType = v.(string)
	}
	if v, ok := d.GetOk("prefix"); ok {
		opt.Prefix = v.(string)
	}
	ids := make([]string, 0)
	for {
		result, response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTencentCosClient(bucket).Bucket.ListMultipartUploads(ctx, opt)
		responseBody, _ := json.Marshal(response.Body)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG]%s api[ListMultipartUploads] success, response body [%s]\n", logId, responseBody)
		for _, item := range result.Uploads {
			itemMap := make(map[string]interface{})
			itemMap["key"] = item.Key
			itemMap["upload_id"] = item.UploadID
			itemMap["initiated"] = item.Initiated
			itemMap["storage_class"] = item.StorageClass
			if item.Owner != nil {
				owner := map[string]interface{}{
					"display_name": item.Owner.DisplayName,
					"id":           item.Owner.ID,
				}
				itemMap["owner"] = []map[string]interface{}{owner}
			}
			if item.Initiator != nil {
				initiator := map[string]interface{}{
					"display_name": item.Initiator.DisplayName,
					"id":           item.Initiator.ID,
				}
				itemMap["initiator"] = []map[string]interface{}{initiator}
			}
			ids = append(ids, item.UploadID)
			multipartUploads = append(multipartUploads, itemMap)
		}
		if result.IsTruncated {
			opt.KeyMarker = result.KeyMarker
			opt.UploadIDMarker = result.UploadIDMarker
		} else {
			break
		}
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("uploads", multipartUploads)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), multipartUploads); err != nil {
			return err
		}
	}

	return nil
}
