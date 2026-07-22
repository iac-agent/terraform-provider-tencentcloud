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
				Description: "Filter by 角色 of the CMK 创建者 `0` - created by 用户，`1` - created by cloud product. 默认值为 `0`。",
			},
			"order_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "顺序 to sort the CMK 创建时间. `0` - desc，`1` - asc. 默认值为 `0`。",
			},
			"key_state": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Filter by state of CMK. `0` - all CMKs are queried，`1` - only 已启用 CMKs are queried，`2` - only 已禁用 CMKs are queried，`3` - only PendingDelete CMKs are queried，`4` - only PendingImport CMKs are queried，`5` - only Archived CMKs are queried。",
			},
			"search_key_alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Words 用于match the results，and the words can be: key_id and alias。",
			},
			"origin": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     KMS_ORIGIN_ALL,
				Description: "Filter by origin of CMK. `TENCENT_KMS` - CMK created by KMS，`EXTERNAL` - CMK imported by 用户，`ALL` - all CMKs. 默认值为 `ALL`。",
			},
			"key_usage": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     KMS_KEY_USAGE_ENCRYPT_DECRYPT,
				Description: "Filter by usage of CMK. Available values include `ALL`，`ENCRYPT_DECRYPT`，`ASYMMETRIC_DECRYPT_RSA_2048`，`ASYMMETRIC_DECRYPT_SM2`，`ASYMMETRIC_SIGN_VERIFY_SM2`，`ASYMMETRIC_SIGN_VERIFY_RSA_2048`，`ASYMMETRIC_SIGN_VERIFY_ECC`. 默认值为 `ENCRYPT_DECRYPT`。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 to filter CMK。",
			},
			"hsm_cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The HSM cluster ID corresponding to KMS Advanced Edition (only valid for KMS Exclusive/Managed Edition service instances)。",
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
							Description: "创建时间 of CMK。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 CMK。",
						},
						"key_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State of CMK。",
						},
						"key_usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Usage of CMK。",
						},
						"creator_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Uin of CMK 创建者",
						},
						"key_rotation_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "指定是否enable 键 rotation。",
						},
						"owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建者 of CMK。",
						},
						"next_rotate_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Next rotate time of CMK when key_rotation_enabled is true。",
						},
						"deletion_date": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Delete time of CMK。",
						},
						"origin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Origin of CMK. `TENCENT_KMS` - CMK created by KMS，`EXTERNAL` - CMK imported by 用户",
						},
						"valid_to": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Valid when origin is `EXTERNAL`，it means the effective date of the 键 material。",
						},
						"hsm_cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The HSM cluster ID corresponding to KMS Advanced Edition (only valid for KMS Exclusive/Managed Edition service instances)。",
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
