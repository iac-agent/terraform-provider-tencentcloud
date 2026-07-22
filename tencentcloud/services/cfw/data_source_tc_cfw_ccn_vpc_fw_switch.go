package cfw

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfwv20190904 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCfwCcnVpcFwSwitch() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCfwCcnVpcFwSwitchRead,
		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CCN ID。",
			},

			"interconnect_pairs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Interconnect pair 配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_a": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Group A。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例 ID",
									},
									"instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例类型 such 作为 VPC 或 DIRECTCONNECT。",
									},
									"instance_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "地域 其中 实例 是 located。",
									},
									"access_cidr_mode": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Network segment 模式 对于 accessing firewall: 0-无 访问，1-访问 all 网络 segments associated 使用 实例，2-访问 用户-defined 网络 segments。",
									},
									"access_cidr_list": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "列表 网络 segments 对于 accessing firewall。",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"group_b": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Group B。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例 ID",
									},
									"instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例类型 such 作为 VPC 或 DIRECTCONNECT。",
									},
									"instance_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "地域 其中 实例 是 located。",
									},
									"access_cidr_mode": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Network segment 模式 对于 accessing firewall: 0-无 访问，1-访问 all 网络 segments associated 使用 实例，2-访问 用户-defined 网络 segments。",
									},
									"access_cidr_list": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "列表 网络 segments 对于 accessing firewall。",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"interconnect_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Interconnect 模式: \"CrossConnect\": cross interconnect (each 实例 在 组 A interconnects 使用 each 实例 在 组 B), \"FullMesh\": full mesh (组 A 内容 是 identical 到 组 B, equivalent 到 pairwise interconnection within 组).",
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

func dataSourceTencentCloudCfwCcnVpcFwSwitchRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cfw_ccn_vpc_fw_switch.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		ccnId   string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("ccn_id"); ok {
		paramMap["CcnId"] = helper.String(v.(string))
		ccnId = v.(string)
	}

	var respData []*cfwv20190904.InterconnectPair
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCfwCcnVpcFwSwitchByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	interconnectPairsList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, interconnectPairs := range respData {
			interconnectPairsMap := map[string]interface{}{}
			groupAList := make([]map[string]interface{}, 0, len(interconnectPairs.GroupA))
			if interconnectPairs.GroupA != nil {
				for _, groupA := range interconnectPairs.GroupA {
					groupAMap := map[string]interface{}{}
					if groupA.InstanceId != nil {
						groupAMap["instance_id"] = groupA.InstanceId
					}

					if groupA.InstanceType != nil {
						groupAMap["instance_type"] = groupA.InstanceType
					}

					if groupA.InstanceRegion != nil {
						groupAMap["instance_region"] = groupA.InstanceRegion
					}

					if groupA.AccessCidrMode != nil {
						groupAMap["access_cidr_mode"] = groupA.AccessCidrMode
					}

					if groupA.AccessCidrList != nil {
						groupAMap["access_cidr_list"] = groupA.AccessCidrList
					}

					groupAList = append(groupAList, groupAMap)
				}

				interconnectPairsMap["group_a"] = groupAList
			}

			groupBList := make([]map[string]interface{}, 0, len(interconnectPairs.GroupB))
			if interconnectPairs.GroupB != nil {
				for _, groupB := range interconnectPairs.GroupB {
					groupBMap := map[string]interface{}{}
					if groupB.InstanceId != nil {
						groupBMap["instance_id"] = groupB.InstanceId
					}

					if groupB.InstanceType != nil {
						groupBMap["instance_type"] = groupB.InstanceType
					}

					if groupB.InstanceRegion != nil {
						groupBMap["instance_region"] = groupB.InstanceRegion
					}

					if groupB.AccessCidrMode != nil {
						groupBMap["access_cidr_mode"] = groupB.AccessCidrMode
					}

					if groupB.AccessCidrList != nil {
						groupBMap["access_cidr_list"] = groupB.AccessCidrList
					}

					groupBList = append(groupBList, groupBMap)
				}

				interconnectPairsMap["group_b"] = groupBList
			}

			if interconnectPairs.InterconnectMode != nil {
				interconnectPairsMap["interconnect_mode"] = interconnectPairs.InterconnectMode
			}

			interconnectPairsList = append(interconnectPairsList, interconnectPairsMap)
		}

		_ = d.Set("interconnect_pairs", interconnectPairsList)
	}

	d.SetId(ccnId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), interconnectPairsList); e != nil {
			return e
		}
	}

	return nil
}
