package controlcenter

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	controlcenterv20230110 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/controlcenter/v20230110"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudControlcenterAccountFactoryBaselineItems() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudControlcenterAccountFactoryBaselineItemsRead,
		Schema: map[string]*schema.Schema{
			"baseline_items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "账号 factory baseline 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "指定unique identifier 对于 账号 factory baseline item，可以 仅 contain `english letters`，`digits`，和 `@,._[]-:()()[]+=.`，使用 长度 的 2-128 字符。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Baseline item 名称 指定a 唯一 名称 对于 功能 item. 支持 combination 的 english letters，numbers，chinese 字符，和 symbols @，&，_，[，]，-. 有效值：1-25 chinese 或 english 字符。",
						},
						"name_en": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Baseline item english 名称 指定a 唯一 名称 对于 baseline item. 支持 combination 的 english letters，digits，spaces，和 symbols @，&，_，[]，-. 有效值：1-64 english 字符。",
						},
						"weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Baseline item 权重 smaller 值， higher 权重 值 范围 equal 到 或 greater 比 0。",
						},
						"required": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "指定是否baseline item 为必填项 (1: 必填; 0: 可选)。",
						},
						"depends_on": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Baseline item dependency. 值 范围 的 N depends 在 count 的 other baseline items 它 relies 在。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Dependency 类型 有效值：LandingZoneSetUp 或 AccountFactorySetUp. LandingZoneSetUp refers 到 dependency 的 landingZone. AccountFactorySetUp refers 到 dependency 的 账号 factory。",
									},
									"identifier": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "指定unique identifier 对于 功能 item，可以 仅 contain `english letters`，`digits`，和 `@,._[]-:()()[]+=.`，使用 长度 的 2-128 字符。",
									},
								},
							},
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Baseline 描述，使用 长度 的 2 到 256 english 或 chinese 字符. 它 是 空 通过 默认值。",
						},
						"description_en": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Baseline item english 描述，使用 长度 的 2 到 1024 english 字符. 它 是 空 通过 默认值。",
						},
						"classify": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Baseline classification. 长度: 2-32 english 或 chinese 字符. 值 不能 是 空。",
						},
						"classify_en": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Baseline english classification，使用 长度 的 2-64 english 字符. 不能 是 空。",
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

func dataSourceTencentCloudControlcenterAccountFactoryBaselineItemsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_controlcenter_account_factory_baseline_items.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ControlcenterService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	var respData []*controlcenterv20230110.AccountFactoryItem
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeControlcenterAccountFactoryBaselineItemsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	baselineItemsList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, baselineItems := range respData {
			baselineItemsMap := map[string]interface{}{}
			if baselineItems.Identifier != nil {
				baselineItemsMap["identifier"] = baselineItems.Identifier
			}

			if baselineItems.Name != nil {
				baselineItemsMap["name"] = baselineItems.Name
			}

			if baselineItems.NameEn != nil {
				baselineItemsMap["name_en"] = baselineItems.NameEn
			}

			if baselineItems.Weight != nil {
				baselineItemsMap["weight"] = baselineItems.Weight
			}

			if baselineItems.Required != nil {
				baselineItemsMap["required"] = baselineItems.Required
			}

			dependsOnList := make([]map[string]interface{}, 0, len(baselineItems.DependsOn))
			if baselineItems.DependsOn != nil {
				for _, dependsOn := range baselineItems.DependsOn {
					dependsOnMap := map[string]interface{}{}
					if dependsOn.Type != nil {
						dependsOnMap["type"] = dependsOn.Type
					}

					if dependsOn.Identifier != nil {
						dependsOnMap["identifier"] = dependsOn.Identifier
					}

					dependsOnList = append(dependsOnList, dependsOnMap)
				}

				baselineItemsMap["depends_on"] = dependsOnList
			}

			if baselineItems.Description != nil {
				baselineItemsMap["description"] = baselineItems.Description
			}

			if baselineItems.DescriptionEn != nil {
				baselineItemsMap["description_en"] = baselineItems.DescriptionEn
			}

			if baselineItems.Classify != nil {
				baselineItemsMap["classify"] = baselineItems.Classify
			}

			if baselineItems.ClassifyEn != nil {
				baselineItemsMap["classify_en"] = baselineItems.ClassifyEn
			}

			baselineItemsList = append(baselineItemsList, baselineItemsMap)
		}

		_ = d.Set("baseline_items", baselineItemsList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), baselineItemsList); e != nil {
			return e
		}
	}

	return nil
}
