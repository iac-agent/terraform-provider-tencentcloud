package kms

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudKmsKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsKeysRead,
		Schema: map[string]*schema.Schema{
			"role": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "过滤器 通过 角色 的 CMK 创建者 `0` - 创建 通过 用户，`1` - 创建 通过 云 product. 默认值为 `0`。",
			},
			"order_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "顺序 到 sort CMK 创建时间. `0` - desc，`1` - asc. 默认值为 `0`。",
			},
			"key_state": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "过滤器 通过 state 的 CMK. `0` - all CMKs 是 queried，`1` - 仅 已启用 CMKs 是 queried，`2` - 仅 已禁用 CMKs 是 queried，`3` - 仅 PendingDelete CMKs 是 queried，`4` - 仅 PendingImport CMKs 是 queried，`5` - 仅 Archived CMKs 是 queried。",
			},
			"search_key_alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Words 用于match results，和 words 可以 是: key_id 和 alias。",
			},
			"origin": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     KMS_ORIGIN_ALL,
				Description: "过滤器 通过 源站 的 CMK. `TENCENT_KMS` - CMK 创建 通过 KMS，`EXTERNAL` - CMK imported 通过 用户，`ALL` - all CMKs. 默认值为 `ALL`。",
			},
			"key_usage": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     KMS_KEY_USAGE_ENCRYPT_DECRYPT,
				Description: "过滤器 通过 usage 的 CMK. Available 值 include `ALL`，`ENCRYPT_DECRYPT`，`ASYMMETRIC_DECRYPT_RSA_2048`，`ASYMMETRIC_DECRYPT_SM2`，`ASYMMETRIC_SIGN_VERIFY_SM2`，`ASYMMETRIC_SIGN_VERIFY_RSA_2048`，`ASYMMETRIC_SIGN_VERIFY_ECC`. 默认值为 `ENCRYPT_DECRYPT`。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 到 过滤器 CMK。",
			},
			"hsm_cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "HSM 集群 ID corresponding 到 KMS Advanced Edition (仅 有效 对于 KMS Exclusive/Managed Edition 服务 实例)。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"key_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 KMS keys。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID CMK。",
						},
						"alias": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 CMK。",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "创建时间 的 CMK。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 CMK。",
						},
						"key_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State 的 CMK。",
						},
						"key_usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Usage 的 CMK。",
						},
						"creator_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Uin 的 CMK 创建者",
						},
						"key_rotation_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "指定是否enable 键 rotation。",
						},
						"owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建者 的 CMK。",
						},
						"next_rotate_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Next rotate 时间 的 CMK 当 key_rotation_enabled 是 true。",
						},
						"deletion_date": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Delete 时间 的 CMK。",
						},
						"origin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Origin 的 CMK. `TENCENT_KMS` - CMK 创建 通过 KMS，`EXTERNAL` - CMK imported 通过 用户",
						},
						"valid_to": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "有效 当 源站 是 `EXTERNAL`，它 表示 effective date 的 键 material。",
						},
						"hsm_cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "HSM 集群 ID corresponding 到 KMS Advanced Edition (仅 有效 对于 KMS Exclusive/Managed Edition 服务 实例)。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudKmsKeysRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_kms_keys.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	param := make(map[string]interface{})
	if v, ok := d.GetOk("role"); ok {
		param["role"] = v.(int)
	}
	if v, ok := d.GetOk("order_type"); ok {
		param["order_type"] = v.(int)
	}
	if v, ok := d.GetOk("key_state"); ok {
		keyState := v.(int)
		param["key_state"] = uint64(keyState)
	}
	if v, ok := d.GetOk("search_key_alias"); ok {
		param["search_key_alias"] = v.(string)
	}
	if v, ok := d.GetOk("origin"); ok {
		param["origin"] = v.(string)
	}
	if v, ok := d.GetOk("key_usage"); ok {
		param["key_usage"] = v.(string)
	}
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		param["tag_filter"] = tags
	}
	if v, ok := d.GetOk("hsm_cluster_id"); ok {
		param["hsm_cluster_id"] = v.(string)
	}

	kmsService := KmsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var keys []*kms.KeyMetadata
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := kmsService.DescribeKeysByFilter(ctx, param)
		if e != nil {
			return tccommon.RetryError(e)
		}
		keys = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read KMS keys failed, reason:%+v", logId, err)
		return err
	}
	keyList := make([]map[string]interface{}, 0, len(keys))
	ids := make([]string, 0, len(keys))
	for _, key := range keys {
		mapping := map[string]interface{}{
			"key_id":               key.KeyId,
			"alias":                key.Alias,
			"create_time":          key.CreateTime,
			"description":          key.Description,
			"key_state":            key.KeyState,
			"key_usage":            key.KeyUsage,
			"creator_uin":          key.CreatorUin,
			"key_rotation_enabled": key.KeyRotationEnabled,
			"owner":                key.Owner,
			"next_rotate_time":     key.NextRotateTime,
			"deletion_date":        key.DeletionDate,
			"origin":               key.Origin,
			"valid_to":             key.ValidTo,
			"hsm_cluster_id":       key.HsmClusterId,
		}

		keyList = append(keyList, mapping)
		ids = append(ids, *key.KeyId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("key_list", keyList); e != nil {
		log.Printf("[CRITAL]%s provider set KMS key list fail, reason:%+v", logId, e)
		return e
	}
	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), keyList)
	}
	return nil
}
