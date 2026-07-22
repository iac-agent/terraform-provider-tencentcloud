package config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudConfigCompliancePacks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudConfigCompliancePacksRead,
		Schema: map[string]*schema.Schema{
			"compliance_pack_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Compliance pack 名称 对于 filtering。",
			},

			"risk_level": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "风险等级 列表 对于 filtering. 有效值：1 (high risk)，2 (medium risk)，3 (low risk)。",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Compliance pack 状态 对于 filtering. 有效值：ACTIVE，NO_ACTIVE。",
			},

			"compliance_result": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Compliance 结果 列表 对于 filtering. 有效值：COMPLIANT，NON_COMPLIANT。",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"order_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort 类型 有效值：desc (descending)，asc (ascending)。",
			},

			"compliance_pack_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Compliance pack 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compliance_pack_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Compliance pack ID。",
						},
						"compliance_pack_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Compliance pack 名称",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Compliance pack 状态 有效值：ACTIVE，NO_ACTIVE。",
						},
						"risk_level": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "风险等级 有效值：1 (high risk)，2 (medium risk)，3 (low risk)。",
						},
						"compliance_result": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Compliance 结果 有效值：COMPLIANT，NON_COMPLIANT。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Compliance pack 描述",
						},
						"rule_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 规则 在 compliance pack。",
						},
						"no_compliant_names": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "列表 non-compliant 规则 names。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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

func dataSourceTencentCloudConfigCompliancePacksRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_config_compliance_packs.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ConfigService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := buildCompliancePacksParamMap(d)

	respData, reqErr := service.DescribeConfigCompliancePacksByFilter(ctx, paramMap)
	if reqErr != nil {
		return reqErr
	}

	packList := flattenConfigCompliancePackList(respData)
	_ = d.Set("compliance_pack_list", packList)

	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}

func buildCompliancePacksParamMap(d *schema.ResourceData) map[string]interface{} {
	paramMap := make(map[string]interface{})

	if v, ok := d.GetOk("compliance_pack_name"); ok {
		paramMap["CompliancePackName"] = v.(string)
	}

	if v, ok := d.GetOk("risk_level"); ok {
		rawList := v.([]interface{})
		riskLevels := make([]*uint64, 0, len(rawList))
		for _, item := range rawList {
			val := uint64(item.(int))
			riskLevels = append(riskLevels, &val)
		}

		paramMap["RiskLevel"] = riskLevels
	}

	if v, ok := d.GetOk("status"); ok {
		paramMap["Status"] = v.(string)
	}

	if v, ok := d.GetOk("compliance_result"); ok {
		rawList := v.([]interface{})
		results := make([]*string, 0, len(rawList))
		for _, item := range rawList {
			val := item.(string)
			results = append(results, &val)
		}

		paramMap["ComplianceResult"] = results
	}

	if v, ok := d.GetOk("order_type"); ok {
		paramMap["OrderType"] = v.(string)
	}

	return paramMap
}

func flattenConfigCompliancePackList(items []*configv20220802.ConfigCompliancePack) []map[string]interface{} {
	packList := make([]map[string]interface{}, 0, len(items))
	for _, pack := range items {
		packMap := map[string]interface{}{}

		if pack.CompliancePackId != nil {
			packMap["compliance_pack_id"] = pack.CompliancePackId
		}

		if pack.CompliancePackName != nil {
			packMap["compliance_pack_name"] = pack.CompliancePackName
		}

		if pack.Status != nil {
			packMap["status"] = pack.Status
		}

		if pack.RiskLevel != nil {
			packMap["risk_level"] = int(*pack.RiskLevel)
		}

		if pack.ComplianceResult != nil {
			packMap["compliance_result"] = pack.ComplianceResult
		}

		if pack.CreateTime != nil {
			packMap["create_time"] = pack.CreateTime
		}

		if pack.Description != nil {
			packMap["description"] = pack.Description
		}

		if pack.RuleCount != nil {
			packMap["rule_count"] = int(*pack.RuleCount)
		}

		if pack.NoCompliantNames != nil {
			names := make([]string, 0, len(pack.NoCompliantNames))
			for _, name := range pack.NoCompliantNames {
				if name != nil {
					names = append(names, *name)
				}
			}

			packMap["no_compliant_names"] = names
		}

		packList = append(packList, packMap)
	}

	return packList
}
