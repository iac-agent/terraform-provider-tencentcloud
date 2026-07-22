package kms

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudKmsDescribeKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsDescribeKeysRead,
		Schema: map[string]*schema.Schema{
			"key_ids": {
				Required:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Query ID 列表 CMK，batch 查询 支持 up 到 100 KeyIds 在 时间。",
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

func dataSourceTencentCloudKmsDescribeKeysRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_kms_describe_keys.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service     = KmsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		keyMetadata []*kms.KeyMetadata
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("key_ids"); ok {
		keyIdsSet := v.(*schema.Set).List()
		paramMap["KeyIds"] = helper.InterfacesStringsPoint(keyIdsSet)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsKeyListsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		keyMetadata = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(keyMetadata))
	tmpList := make([]map[string]interface{}, 0, len(keyMetadata))

	if keyMetadata != nil {
		for _, key := range keyMetadata {
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
			}

			tmpList = append(tmpList, mapping)
			ids = append(ids, *key.KeyId)
		}

		_ = d.Set("key_list", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
