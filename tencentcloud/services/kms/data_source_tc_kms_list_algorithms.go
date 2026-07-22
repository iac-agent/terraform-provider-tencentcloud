package kms

import (
	"context"
	"strconv"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
)

func DataSourceTencentCloudKmsListAlgorithms() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsListAlgorithmsRead,
		Schema: map[string]*schema.Schema{
			"symmetric_algorithms": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Symmetric 加密 algorithms 支持 在 此 地域",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "键 usage。",
						},
						"algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Algorithm。",
						},
					},
				},
			},
			"asymmetric_algorithms": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Asymmetric 加密 algorithms 支持 在 此 地域",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "键 usage。",
						},
						"algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Algorithm。",
						},
					},
				},
			},
			"asymmetric_sign_verify_algorithms": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Asymmetric 签名 verification algorithms 支持 在 此 地域",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "键 usage。",
						},
						"algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Algorithm。",
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

func dataSourceTencentCloudKmsListAlgorithmsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_kms_list_algorithms.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service        = KmsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		listAlgorithms *kms.ListAlgorithmsResponseParams
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsListAlgorithmsByFilter(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}

		listAlgorithms = result
		return nil
	})

	if err != nil {
		return err
	}

	if listAlgorithms.SymmetricAlgorithms != nil {
		tmpList := make([]map[string]interface{}, 0, len(listAlgorithms.SymmetricAlgorithms))
		for _, item := range listAlgorithms.SymmetricAlgorithms {
			itemMap := map[string]interface{}{}
			if item.KeyUsage != nil {
				itemMap["key_usage"] = item.KeyUsage
			}

			if item.Algorithm != nil {
				itemMap["algorithm"] = item.Algorithm
			}

			tmpList = append(tmpList, itemMap)
		}

		_ = d.Set("symmetric_algorithms", tmpList)
	}

	if listAlgorithms.AsymmetricAlgorithms != nil {
		tmpList := make([]map[string]interface{}, 0, len(listAlgorithms.AsymmetricAlgorithms))
		for _, item := range listAlgorithms.AsymmetricAlgorithms {
			itemMap := map[string]interface{}{}
			if item.KeyUsage != nil {
				itemMap["key_usage"] = item.KeyUsage
			}

			if item.Algorithm != nil {
				itemMap["algorithm"] = item.Algorithm
			}

			tmpList = append(tmpList, itemMap)
		}

		_ = d.Set("asymmetric_algorithms", tmpList)
	}

	if listAlgorithms.AsymmetricSignVerifyAlgorithms != nil {
		tmpList := make([]map[string]interface{}, 0, len(listAlgorithms.AsymmetricSignVerifyAlgorithms))
		for _, item := range listAlgorithms.AsymmetricSignVerifyAlgorithms {
			itemMap := map[string]interface{}{}
			if item.KeyUsage != nil {
				itemMap["key_usage"] = item.KeyUsage
			}

			if item.Algorithm != nil {
				itemMap["algorithm"] = item.Algorithm
			}

			tmpList = append(tmpList, itemMap)
		}

		_ = d.Set("asymmetric_sign_verify_algorithms", tmpList)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
