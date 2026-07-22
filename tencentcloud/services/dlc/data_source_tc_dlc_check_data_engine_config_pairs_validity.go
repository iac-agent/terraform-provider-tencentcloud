package dlc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDlcCheckDataEngineConfigPairsValidity() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcCheckDataEngineConfigPairsValidityRead,
		Schema: map[string]*schema.Schema{
			"child_image_version_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ID minor 版本 的 引擎。",
			},

			"data_engine_config_pairs": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "用户-defined 参数。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_item": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration item。",
						},
						"config_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration 值",
						},
					},
				},
			},

			"image_version_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ID major 版本 的 引擎. 如果 there 是 ID minor 版本，仅 ID minor 版本 needs 到 是 input. 如果 不， latest ID minor 版本 under major 版本 将 是 acquired。",
			},

			"is_available": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Parameter validity: true: 有效，false: 在 least 一个 无效 参数 exists。",
			},

			"unavailable_config": {
				Computed:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Invalid 参数 集合。",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudDlcCheckDataEngineConfigPairsValidityRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_check_data_engine_config_pairs_validity.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var childImageVersionId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("child_image_version_id"); ok {
		childImageVersionId = v.(string)
		paramMap["ChildImageVersionId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("data_engine_config_pairs"); ok {
		dataEngineConfigPairsSet := v.([]interface{})
		tmpSet := make([]*dlc.DataEngineConfigPair, 0, len(dataEngineConfigPairsSet))

		for _, item := range dataEngineConfigPairsSet {
			dataEngineConfigPair := dlc.DataEngineConfigPair{}
			dataEngineConfigPairMap := item.(map[string]interface{})

			if v, ok := dataEngineConfigPairMap["config_item"]; ok {
				dataEngineConfigPair.ConfigItem = helper.String(v.(string))
			}
			if v, ok := dataEngineConfigPairMap["config_value"]; ok {
				dataEngineConfigPair.ConfigValue = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &dataEngineConfigPair)
		}
		paramMap["data_engine_config_pairs"] = tmpSet
	}

	if v, ok := d.GetOk("image_version_id"); ok {
		paramMap["ImageVersionId"] = helper.String(v.(string))
	}

	service := DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var data *dlc.CheckDataEngineConfigPairsValidityResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcCheckDataEngineConfigPairsValidityByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		data = result
		return nil
	})
	if err != nil {
		return err
	}
	result := make(map[string]interface{})
	if data.IsAvailable != nil {
		_ = d.Set("is_available", data.IsAvailable)
		result["is_available"] = data.IsAvailable
	}

	if data.UnavailableConfig != nil {
		_ = d.Set("unavailable_config", data.UnavailableConfig)
		result["unavailable_config"] = data.UnavailableConfig
	}

	d.SetId(childImageVersionId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
