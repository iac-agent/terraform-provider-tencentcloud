package cam

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	camv20190116 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func DataSourceTencentCloudCamPolicyDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamPolicyDetailRead,
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Policy ID。",
			},

			"policy_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Policy detail 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy 名称",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy 描述",
						},
						"type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Policy 类型 1 表示 自定义 策略，2 表示 preset 策略。",
						},
						"add_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time 策略 是 创建。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time 策略 是 last 更新。",
						},
						"policy_document": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy document。",
						},
						"preset_alias": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Preset 策略 alias. 注意: 此 字段 可能 返回 null。",
						},
						"is_service_linked_role_policy": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否policy 是 服务-linked 角色 策略. 0 表示 无，1 表示 yes。",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "标签 associated 使用 策略。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签键",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签值",
									},
								},
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

func dataSourceTencentCloudCamPolicyDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cam_policy_detail.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CamService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	policyIdRaw := d.Get("policy_id").(int)
	policyId := uint64(policyIdRaw)
	paramMap["PolicyId"] = &policyId

	var respData *camv20190116.GetPolicyResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCamPolicyDetailByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	policyInfoMap := map[string]interface{}{}

	if respData != nil {
		if respData.PolicyName != nil {
			policyInfoMap["policy_name"] = respData.PolicyName
		}

		if respData.Description != nil {
			policyInfoMap["description"] = respData.Description
		}

		if respData.Type != nil {
			policyInfoMap["type"] = int(*respData.Type)
		}

		if respData.AddTime != nil {
			policyInfoMap["add_time"] = respData.AddTime
		}

		if respData.UpdateTime != nil {
			policyInfoMap["update_time"] = respData.UpdateTime
		}

		if respData.PolicyDocument != nil {
			policyInfoMap["policy_document"] = respData.PolicyDocument
		}

		if respData.PresetAlias != nil {
			policyInfoMap["preset_alias"] = respData.PresetAlias
		}

		if respData.IsServiceLinkedRolePolicy != nil {
			policyInfoMap["is_service_linked_role_policy"] = int(*respData.IsServiceLinkedRolePolicy)
		}

		tagsList := make([]map[string]interface{}, 0, len(respData.Tags))
		for _, tag := range respData.Tags {
			tagMap := map[string]interface{}{}
			if tag.Key != nil {
				tagMap["key"] = tag.Key
			}

			if tag.Value != nil {
				tagMap["value"] = tag.Value
			}

			tagsList = append(tagsList, tagMap)
		}

		policyInfoMap["tags"] = tagsList
		_ = d.Set("policy_info", []interface{}{policyInfoMap})
	}

	d.SetId(strconv.Itoa(policyIdRaw))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), policyInfoMap); e != nil {
			return e
		}
	}

	return nil
}
