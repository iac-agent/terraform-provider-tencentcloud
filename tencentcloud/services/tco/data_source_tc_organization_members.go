package tco

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudOrganizationMembers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOrganizationMembersRead,
		Schema: map[string]*schema.Schema{
			"lang": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "有效值：`en` (Tencent Cloud International); `zh` (Tencent Cloud)。",
			},

			"search_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Search 通过 member 名称 或 ID。",
			},

			"auth_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Entity 名称",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Abbreviation 的 trusted 服务，其中 为必填项 during querying trusted 服务 admin。",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Member 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"member_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Member UIN注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member name注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"member_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member 类型 有效值：`Invite` (invited); `Create` (创建).注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"org_policy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Relationship 策略 type注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"org_policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Relationship 策略 name注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"org_permission": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Relationship 策略 permission注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Permission ID。",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Permission 名称",
									},
								},
							},
						},
						"node_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node ID注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"node_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node name注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remarks注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"is_allow_quit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "是否member 是 allowed 到 leave. 有效值：`Allow`，`Denied`.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"pay_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Payer UIN注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"pay_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Payer name注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"org_identity": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Management identity注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"identity_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Identity ID.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"identity_alias_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Identity 名称注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"bind_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security 信息 binding 状态 有效值：`Unbound`，`有效`，`Success`，`Failed`.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"permission_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member 权限 状态 有效值：`Confirmed`，`UnConfirmed`.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"nick_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tencent Cloud nickname. 注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudOrganizationMembersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_organization_members.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("lang"); ok {
		paramMap["Lang"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_key"); ok {
		paramMap["SearchKey"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("auth_name"); ok {
		paramMap["AuthName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var items []*organization.OrgMember

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOrganizationMembersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		items = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(items))
	tmpList := make([]map[string]interface{}, 0, len(items))

	if items != nil {
		for _, orgMember := range items {
			orgMemberMap := map[string]interface{}{}

			if orgMember.MemberUin != nil {
				orgMemberMap["member_uin"] = orgMember.MemberUin
			}

			if orgMember.Name != nil {
				orgMemberMap["name"] = orgMember.Name
			}

			if orgMember.MemberType != nil {
				orgMemberMap["member_type"] = orgMember.MemberType
			}

			if orgMember.OrgPolicyType != nil {
				orgMemberMap["org_policy_type"] = orgMember.OrgPolicyType
			}

			if orgMember.OrgPolicyName != nil {
				orgMemberMap["org_policy_name"] = orgMember.OrgPolicyName
			}

			if orgMember.OrgPermission != nil {
				orgPermissionList := []interface{}{}
				for _, orgPermission := range orgMember.OrgPermission {
					orgPermissionMap := map[string]interface{}{}

					if orgPermission.Id != nil {
						orgPermissionMap["id"] = orgPermission.Id
					}

					if orgPermission.Name != nil {
						orgPermissionMap["name"] = orgPermission.Name
					}

					orgPermissionList = append(orgPermissionList, orgPermissionMap)
				}

				orgMemberMap["org_permission"] = orgPermissionList
			}

			if orgMember.NodeId != nil {
				orgMemberMap["node_id"] = orgMember.NodeId
			}

			if orgMember.NodeName != nil {
				orgMemberMap["node_name"] = orgMember.NodeName
			}

			if orgMember.Remark != nil {
				orgMemberMap["remark"] = orgMember.Remark
			}

			if orgMember.CreateTime != nil {
				orgMemberMap["create_time"] = orgMember.CreateTime
			}

			if orgMember.UpdateTime != nil {
				orgMemberMap["update_time"] = orgMember.UpdateTime
			}

			if orgMember.IsAllowQuit != nil {
				orgMemberMap["is_allow_quit"] = orgMember.IsAllowQuit
			}

			if orgMember.PayUin != nil {
				orgMemberMap["pay_uin"] = orgMember.PayUin
			}

			if orgMember.PayName != nil {
				orgMemberMap["pay_name"] = orgMember.PayName
			}

			if orgMember.OrgIdentity != nil {
				orgIdentityList := []interface{}{}
				for _, orgIdentity := range orgMember.OrgIdentity {
					orgIdentityMap := map[string]interface{}{}

					if orgIdentity.IdentityId != nil {
						orgIdentityMap["identity_id"] = orgIdentity.IdentityId
					}

					if orgIdentity.IdentityAliasName != nil {
						orgIdentityMap["identity_alias_name"] = orgIdentity.IdentityAliasName
					}

					orgIdentityList = append(orgIdentityList, orgIdentityMap)
				}

				orgMemberMap["org_identity"] = orgIdentityList
			}

			if orgMember.BindStatus != nil {
				orgMemberMap["bind_status"] = orgMember.BindStatus
			}

			if orgMember.PermissionStatus != nil {
				orgMemberMap["permission_status"] = orgMember.PermissionStatus
			}

			if orgMember.NickName != nil {
				orgMemberMap["nick_name"] = orgMember.NickName
			}

			ids = append(ids, *orgMember.Name)
			tmpList = append(tmpList, orgMemberMap)
		}

		_ = d.Set("items", tmpList)
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
