package dlc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlcv20210125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDlcDataEngineSessionParameters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcDataEngineSessionParametersRead,
		Schema: map[string]*schema.Schema{
			"data_engine_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DataEngine ID。",
			},

			"data_engine_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Engine 名称 当 引擎 名称 是 指定， 名称 是 使用 first 到 obtain 配置。",
			},

			"data_engine_parameters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Engine Session Configuration List。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration ID。",
						},
						"child_image_version_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Minor 版本 镜像 ID。",
						},
						"engine_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "集群类型: SparkSQL/PrestoSQL/SparkBatch。",
						},
						"key_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter 键",
						},
						"key_description": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "描述 键",
						},
						"value_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "类型 值",
						},
						"value_length_limit": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Length 限制 的 值",
						},
						"value_regexp_limit": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Regular expression constraint 对于 值",
						},
						"value_default": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "默认值",
						},
						"is_public": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "是否为a 公有 版本: 1 对于 公有; 2 对于 私有。",
						},
						"parameter_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Configuration 类型: 1 对于 会话 配置 (默认值); 2 对于 common 配置; 3 对于 集群 配置",
						},
						"submit_method": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Submission 方法: 用户 或 BackGround。",
						},
						"operator": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "操作者",
						},
						"insert_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Insert 时间。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "更新时间。",
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

func dataSourceTencentCloudDlcDataEngineSessionParametersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_data_engine_session_parameters.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(nil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service      = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		dataEngineId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("data_engine_id"); ok {
		paramMap["DataEngineId"] = helper.String(v.(string))
		dataEngineId = v.(string)
	}

	if v, ok := d.GetOk("data_engine_name"); ok {
		paramMap["DataEngineName"] = helper.String(v.(string))
	}

	var respData []*dlcv20210125.DataEngineImageSessionParameter
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcDataEngineSessionParametersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	dataEngineParametersList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, dataEngineParameters := range respData {
			dataEngineParametersMap := map[string]interface{}{}
			if dataEngineParameters.ParameterId != nil {
				dataEngineParametersMap["parameter_id"] = dataEngineParameters.ParameterId
			}

			if dataEngineParameters.ChildImageVersionId != nil {
				dataEngineParametersMap["child_image_version_id"] = dataEngineParameters.ChildImageVersionId
			}

			if dataEngineParameters.EngineType != nil {
				dataEngineParametersMap["engine_type"] = dataEngineParameters.EngineType
			}

			if dataEngineParameters.KeyName != nil {
				dataEngineParametersMap["key_name"] = dataEngineParameters.KeyName
			}

			if dataEngineParameters.KeyDescription != nil {
				dataEngineParametersMap["key_description"] = dataEngineParameters.KeyDescription
			}

			if dataEngineParameters.ValueType != nil {
				dataEngineParametersMap["value_type"] = dataEngineParameters.ValueType
			}

			if dataEngineParameters.ValueLengthLimit != nil {
				dataEngineParametersMap["value_length_limit"] = dataEngineParameters.ValueLengthLimit
			}

			if dataEngineParameters.ValueRegexpLimit != nil {
				dataEngineParametersMap["value_regexp_limit"] = dataEngineParameters.ValueRegexpLimit
			}

			if dataEngineParameters.ValueDefault != nil {
				dataEngineParametersMap["value_default"] = dataEngineParameters.ValueDefault
			}

			if dataEngineParameters.IsPublic != nil {
				dataEngineParametersMap["is_public"] = dataEngineParameters.IsPublic
			}

			if dataEngineParameters.ParameterType != nil {
				dataEngineParametersMap["parameter_type"] = dataEngineParameters.ParameterType
			}

			if dataEngineParameters.SubmitMethod != nil {
				dataEngineParametersMap["submit_method"] = dataEngineParameters.SubmitMethod
			}

			if dataEngineParameters.Operator != nil {
				dataEngineParametersMap["operator"] = dataEngineParameters.Operator
			}

			if dataEngineParameters.InsertTime != nil {
				dataEngineParametersMap["insert_time"] = dataEngineParameters.InsertTime
			}

			if dataEngineParameters.UpdateTime != nil {
				dataEngineParametersMap["update_time"] = dataEngineParameters.UpdateTime
			}

			dataEngineParametersList = append(dataEngineParametersList, dataEngineParametersMap)
		}

		_ = d.Set("data_engine_parameters", dataEngineParametersList)
	}

	d.SetId(dataEngineId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataEngineParametersList); e != nil {
			return e
		}
	}

	return nil
}
