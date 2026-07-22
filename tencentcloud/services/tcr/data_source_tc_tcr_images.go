package tcr

import (
	"context"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTcrImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcrImagesRead,
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID",
			},

			"namespace_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "命名空间 名称",
			},

			"repository_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "repository 名称",
			},

			"image_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "镜像 版本 名称，默认为 fuzzy match。",
			},

			"digest": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "指定image digest 对于 lookup。",
			},

			"exact_match": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "指定whether 它 是 exact match，true 是 exact match，和 不 filled 是 fuzzy match。",
			},

			"image_info_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "容器 镜像 信息 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"digest": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "hash 值",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "镜像 大小 (单位: byte)。",
						},
						"image_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "标签 名称",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "更新时间。",
						},
						"kind": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "产品类型,note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
						},
						"kms_signature": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "kms 签名 信息,note: 此 字段 可能 返回 null，indicating 该 无 有效 值 可以 是 获取。",
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

func dataSourceTencentCloudTcrImagesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tcr_images.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		registryId    string
		namespaceName string
		repoName      string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("registry_id"); ok {
		paramMap["registry_id"] = helper.String(v.(string))
		registryId = v.(string)
	}

	if v, ok := d.GetOk("namespace_name"); ok {
		paramMap["namespace_name"] = helper.String(v.(string))
		namespaceName = v.(string)
	}

	if v, ok := d.GetOk("repository_name"); ok {
		paramMap["repository_name"] = helper.String(v.(string))
		repoName = v.(string)
	}

	if v, ok := d.GetOk("image_version"); ok {
		paramMap["image_version"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("digest"); ok {
		paramMap["digest"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("exact_match"); v != nil {
		paramMap["exact_match"] = helper.Bool(v.(bool))
	}

	service := TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var imageInfoList []*tcr.TcrImageInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTcrImagesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		imageInfoList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(imageInfoList))
	tmpList := make([]map[string]interface{}, 0, len(imageInfoList))

	if imageInfoList != nil {
		for _, tcrImageInfo := range imageInfoList {
			tcrImageInfoMap := map[string]interface{}{}

			if tcrImageInfo.Digest != nil {
				tcrImageInfoMap["digest"] = tcrImageInfo.Digest
			}

			if tcrImageInfo.Size != nil {
				tcrImageInfoMap["size"] = tcrImageInfo.Size
			}

			if tcrImageInfo.ImageVersion != nil {
				tcrImageInfoMap["image_version"] = tcrImageInfo.ImageVersion
			}

			if tcrImageInfo.UpdateTime != nil {
				tcrImageInfoMap["update_time"] = tcrImageInfo.UpdateTime
			}

			if tcrImageInfo.Kind != nil {
				tcrImageInfoMap["kind"] = tcrImageInfo.Kind
			}

			if tcrImageInfo.KmsSignature != nil {
				tcrImageInfoMap["kms_signature"] = tcrImageInfo.KmsSignature
			}

			ids = append(ids, strings.Join([]string{registryId, namespaceName, repoName, *tcrImageInfo.ImageVersion}, tccommon.FILED_SP))
			tmpList = append(tmpList, tcrImageInfoMap)
		}

		_ = d.Set("image_info_list", tmpList)
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
