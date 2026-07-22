package scf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudScfLayerVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfLayerVersionsRead,
		Schema: map[string]*schema.Schema{
			"layer_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Layer 名称",
			},

			"compatible_runtime": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Compatible runtimes。",
			},

			"layer_versions": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Layer 版本 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compatible_runtimes": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Runtime applicable 到 version注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"add_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "版本 description注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"license_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "License information注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"layer_version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "版本 数量。",
						},
						"layer_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Layer 名称",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current 状态 特定 layer 版本 For 有效 值，please see [here](https://intl.云.tencent.com/document/product/583/47175?from_cn_redirect=1#.E5.B1.82.EF.BC.88layer.EF.BC.89.E7.8A.B6.E6.80.81)。",
						},
						"stamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stamp注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudScfLayerVersionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_scf_layer_versions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("layer_name"); ok {
		paramMap["LayerName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("compatible_runtime"); ok {
		compatibleRuntimeSet := v.(*schema.Set).List()
		paramMap["CompatibleRuntime"] = helper.InterfacesStringsPoint(compatibleRuntimeSet)
	}

	service := ScfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var layerVersions []*scf.LayerVersionInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeScfLayerVersions(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		layerVersions = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(layerVersions))
	tmpList := make([]map[string]interface{}, 0, len(layerVersions))

	if layerVersions != nil {
		for _, layerVersionInfo := range layerVersions {
			layerVersionInfoMap := map[string]interface{}{}

			if layerVersionInfo.CompatibleRuntimes != nil {
				layerVersionInfoMap["compatible_runtimes"] = layerVersionInfo.CompatibleRuntimes
			}

			if layerVersionInfo.AddTime != nil {
				layerVersionInfoMap["add_time"] = layerVersionInfo.AddTime
			}

			if layerVersionInfo.Description != nil {
				layerVersionInfoMap["description"] = layerVersionInfo.Description
			}

			if layerVersionInfo.LicenseInfo != nil {
				layerVersionInfoMap["license_info"] = layerVersionInfo.LicenseInfo
			}

			if layerVersionInfo.LayerVersion != nil {
				layerVersionInfoMap["layer_version"] = layerVersionInfo.LayerVersion
			}

			if layerVersionInfo.LayerName != nil {
				layerVersionInfoMap["layer_name"] = layerVersionInfo.LayerName
			}

			if layerVersionInfo.Status != nil {
				layerVersionInfoMap["status"] = layerVersionInfo.Status
			}

			if layerVersionInfo.Stamp != nil {
				layerVersionInfoMap["stamp"] = layerVersionInfo.Stamp
			}

			ids = append(ids, *layerVersionInfo.LayerName)
			tmpList = append(tmpList, layerVersionInfoMap)
		}

		_ = d.Set("layer_versions", tmpList)
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
