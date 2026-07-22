package kms

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudKmsWhiteBoxKeyDetails() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsWhiteBoxKeyDetailsRead,
		Schema: map[string]*schema.Schema{
			"key_status": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "过滤器 condition: 状态 键，0: 已禁用，1: 已启用",
			},
			"key_infos": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 white box 键 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 algorithm 使用 通过 键",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "键 创建时间，Unix 时间戳。",
						},
						"decrypt_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "White box decryption 键，base64 encoded。",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "资源 ID，格式: creatorUin/$creatorUin/$keyId。",
						},
						"key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Globally 唯一 identifier 对于 white box 键",
						},
						"creator_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "创建者",
						},
						"alias": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "As alias 对于 键 该 是 easier 到 identify 和 easier 到 understand，它 不能 是 空 和 是 combination 的 1-60 alphanumeric 字符 - _. first character 必须 是 letter 或 数量. It 不能 是 repeated。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 键",
						},
						"encrypt_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "White box 加密 键，base64 encoded。",
						},
						"owner_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "创建者",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态 white box 键， 值 是: 已启用 | 已禁用",
						},
						"device_fingerprint_bind": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is there device fingerprint bound 到 当前 键?。",
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

func dataSourceTencentCloudKmsWhiteBoxKeyDetailsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_kms_white_box_key_details.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service         = KmsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		whiteBoxKeyInfo []*kms.WhiteboxKeyInfo
	)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("key_status"); v != nil {
		paramMap["KeyStatus"] = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsWhiteBoxKeyDetailsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		whiteBoxKeyInfo = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(whiteBoxKeyInfo))
	tmpList := make([]map[string]interface{}, 0, len(whiteBoxKeyInfo))

	if whiteBoxKeyInfo != nil {
		for _, whiteBoxKey := range whiteBoxKeyInfo {
			whiteBoxKeyInfoMap := map[string]interface{}{}

			if whiteBoxKey.Algorithm != nil {
				whiteBoxKeyInfoMap["algorithm"] = whiteBoxKey.Algorithm
			}

			if whiteBoxKey.CreateTime != nil {
				whiteBoxKeyInfoMap["create_time"] = whiteBoxKey.CreateTime
			}

			if whiteBoxKey.DecryptKey != nil {
				whiteBoxKeyInfoMap["decrypt_key"] = whiteBoxKey.DecryptKey
			}

			if whiteBoxKey.ResourceId != nil {
				whiteBoxKeyInfoMap["resource_id"] = whiteBoxKey.ResourceId
			}

			if whiteBoxKey.KeyId != nil {
				whiteBoxKeyInfoMap["key_id"] = whiteBoxKey.KeyId
			}

			if whiteBoxKey.CreatorUin != nil {
				whiteBoxKeyInfoMap["creator_uin"] = whiteBoxKey.CreatorUin
			}

			if whiteBoxKey.Alias != nil {
				whiteBoxKeyInfoMap["alias"] = whiteBoxKey.Alias
			}

			if whiteBoxKey.Description != nil {
				whiteBoxKeyInfoMap["description"] = whiteBoxKey.Description
			}

			if whiteBoxKey.EncryptKey != nil {
				whiteBoxKeyInfoMap["encrypt_key"] = whiteBoxKey.EncryptKey
			}

			if whiteBoxKey.OwnerUin != nil {
				whiteBoxKeyInfoMap["owner_uin"] = whiteBoxKey.OwnerUin
			}

			if whiteBoxKey.Status != nil {
				whiteBoxKeyInfoMap["status"] = whiteBoxKey.Status
			}

			if whiteBoxKey.DeviceFingerprintBind != nil {
				whiteBoxKeyInfoMap["device_fingerprint_bind"] = whiteBoxKey.DeviceFingerprintBind
			}

			ids = append(ids, *whiteBoxKey.KeyId)
			tmpList = append(tmpList, whiteBoxKeyInfoMap)
		}

		_ = d.Set("key_infos", tmpList)
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
