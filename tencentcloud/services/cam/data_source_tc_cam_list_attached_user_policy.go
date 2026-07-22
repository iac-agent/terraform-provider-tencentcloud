package cam

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCamListAttachedUserPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamListAttachedUserPolicyRead,
		Schema: map[string]*schema.Schema{
			"target_uin": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Target 用户 ID。",
			},

			"attach_type": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "0: Return direct association and group association policies，1: Only return direct association policies，2: Only return group association policies。",
			},

			"strategy_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Policy 类型",
			},

			"keyword": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Search Keywords。",
			},

			"policy_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Policy List Data。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy ID。",
						},
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
						"add_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"strategy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy 类型 (1 represents custom policy，2 represents preset policy)。",
						},
						"create_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation 模式 (1 represents policies created by product or project permissions，others represent policies created by policy syntax)。",
						},
						"groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Associated information with group注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "组 ID",
									},
									"group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Group 名称",
									},
								},
							},
						},
						"deactived": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Has it been taken offline (0: No 1: Yes)注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"deactived_detail": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "列表 offline products注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudCamListAttachedUserPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cam_list_attached_user_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOkExists("target_uin"); v != nil {
		paramMap["TargetUin"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOkExists("attach_type"); v != nil {
		paramMap["AttachType"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOkExists("strategy_type"); v != nil {
		paramMap["StrategyType"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("keyword"); ok {
		paramMap["Keyword"] = helper.String(v.(string))
	}

	service := CamService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var policyList []*cam.AttachedUserPolicy

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCamListAttachedUserPolicyByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		policyList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(policyList))
	tmpList := make([]map[string]interface{}, 0, len(policyList))

	if policyList != nil {
		for _, attachedUserPolicy := range policyList {
			attachedUserPolicyMap := map[string]interface{}{}

			if attachedUserPolicy.PolicyId != nil {
				attachedUserPolicyMap["policy_id"] = attachedUserPolicy.PolicyId
			}

			if attachedUserPolicy.PolicyName != nil {
				attachedUserPolicyMap["policy_name"] = attachedUserPolicy.PolicyName
			}

			if attachedUserPolicy.Description != nil {
				attachedUserPolicyMap["description"] = attachedUserPolicy.Description
			}

			if attachedUserPolicy.AddTime != nil {
				attachedUserPolicyMap["add_time"] = attachedUserPolicy.AddTime
			}

			if attachedUserPolicy.StrategyType != nil {
				attachedUserPolicyMap["strategy_type"] = attachedUserPolicy.StrategyType
			}

			if attachedUserPolicy.CreateMode != nil {
				attachedUserPolicyMap["create_mode"] = attachedUserPolicy.CreateMode
			}

			if attachedUserPolicy.Groups != nil {
				groupsList := []interface{}{}
				for _, groups := range attachedUserPolicy.Groups {
					groupsMap := map[string]interface{}{}

					if groups.GroupId != nil {
						groupsMap["group_id"] = groups.GroupId
					}

					if groups.GroupName != nil {
						groupsMap["group_name"] = groups.GroupName
					}

					groupsList = append(groupsList, groupsMap)
				}

				attachedUserPolicyMap["groups"] = groupsList
			}

			if attachedUserPolicy.Deactived != nil {
				attachedUserPolicyMap["deactived"] = attachedUserPolicy.Deactived
			}

			if attachedUserPolicy.DeactivedDetail != nil {
				attachedUserPolicyMap["deactived_detail"] = attachedUserPolicy.DeactivedDetail
			}

			ids = append(ids, *attachedUserPolicy.PolicyId)
			tmpList = append(tmpList, attachedUserPolicyMap)
		}

		_ = d.Set("policy_list", tmpList)
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
