package tco

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudOrganizationOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationOrganizationCreate,
		Read:   resourceTencentCloudOrganizationOrganizationRead,
		Update: resourceTencentCloudOrganizationOrganizationUpdate,
		Delete: resourceTencentCloudOrganizationOrganizationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"org_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Enterprise organization ID.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
			},

			"root_node_name": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Root 节点名称",
			},

			"host_uin": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "创建者 Uin.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
			},

			"nick_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "创建者 nickname.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
			},

			"org_type": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Enterprise organization 类型Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
			},

			"is_manager": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "是否organize administrator.Yes: true，无: falseNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
			},

			"org_policy_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Strategy 类型Financial Management: FinancialNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
			},

			"org_policy_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Strategic 名称Note: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
			},

			"org_permission": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 membership authority 的 members.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Permissions ID。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Permission 名称",
						},
					},
				},
			},

			"root_node_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Organize root 节点 IDNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Organize 创建时间.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
			},

			"join_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Members join 时间.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
			},

			"is_allow_quit": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "是否members 是 allowed 到 withdraw.Allow: Allow，不 allowed: DENIEDNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
			},

			"pay_uin": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "UIN 在 behalf 的 payer.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
			},

			"pay_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "名称 payment.注意: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
			},

			"is_assign_manager": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether trusted 服务 administrator.Yes: true，无: falseNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
			},

			"is_auth_manager": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "是否real -名称 subject administrator.Yes: true，无: falseNote: 此 字段 可能 返回 NULL，indicating 该 有效 值 不能 是 获取。",
			},
		},
	}
}

func resourceTencentCloudOrganizationOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request  = organization.NewCreateOrganizationRequest()
		response = organization.NewCreateOrganizationResponse()
		orgId    uint64
	)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().CreateOrganization(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create organization organization failed, reason:%+v", logId, err)
		return err
	}

	orgId = *response.Response.OrgId

	if v, ok := d.GetOk("root_node_name"); ok {
		service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		innerErr := service.UpdateOrganizationRootNodeName(ctx, orgId, v.(string))
		if innerErr != nil {
			return innerErr
		}
	}

	d.SetId(helper.UInt64ToStr(orgId))

	return resourceTencentCloudOrganizationOrganizationRead(d, meta)
}

func resourceTencentCloudOrganizationOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	organization, err := service.DescribeOrganizationOrganizationById(ctx)
	if err != nil {
		return err
	}

	if organization == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OrganizationOrganization` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if organization.OrgId != nil {
		_ = d.Set("org_id", organization.OrgId)
	}

	if organization.HostUin != nil {
		_ = d.Set("host_uin", organization.HostUin)
	}

	if organization.NickName != nil {
		_ = d.Set("nick_name", organization.NickName)
	}

	if organization.OrgType != nil {
		_ = d.Set("org_type", organization.OrgType)
	}

	if organization.IsManager != nil {
		_ = d.Set("is_manager", organization.IsManager)
	}

	if organization.OrgPolicyType != nil {
		_ = d.Set("org_policy_type", organization.OrgPolicyType)
	}

	if organization.OrgPolicyName != nil {
		_ = d.Set("org_policy_name", organization.OrgPolicyName)
	}

	if organization.OrgPermission != nil {
		var orgPermissionList []interface{}
		for _, orgPermission := range organization.OrgPermission {
			orgPermissionMap := map[string]interface{}{}

			if orgPermission.Id != nil {
				orgPermissionMap["id"] = orgPermission.Id
			}

			if orgPermission.Name != nil {
				orgPermissionMap["name"] = orgPermission.Name
			}

			orgPermissionList = append(orgPermissionList, orgPermissionMap)
		}

		_ = d.Set("org_permission", orgPermissionList)

	}

	if organization.RootNodeId != nil {
		_ = d.Set("root_node_id", organization.RootNodeId)
	}

	if organization.CreateTime != nil {
		_ = d.Set("create_time", organization.CreateTime)
	}

	if organization.JoinTime != nil {
		_ = d.Set("join_time", organization.JoinTime)
	}

	if organization.IsAllowQuit != nil {
		_ = d.Set("is_allow_quit", organization.IsAllowQuit)
	}

	if organization.PayUin != nil {
		_ = d.Set("pay_uin", organization.PayUin)
	}

	if organization.PayName != nil {
		_ = d.Set("pay_name", organization.PayName)
	}

	if organization.IsAssignManager != nil {
		_ = d.Set("is_assign_manager", organization.IsAssignManager)
	}

	if organization.IsAuthManager != nil {
		_ = d.Set("is_auth_manager", organization.IsAuthManager)
	}

	if organization.RootNodeId != nil {
		orgNode, err := service.DescribeOrganizationOrgNode(ctx, helper.Int64ToStr(*organization.RootNodeId))
		if err != nil {
			return err
		}
		_ = d.Set("root_node_name", orgNode.Name)
	}
	return nil
}

func resourceTencentCloudOrganizationOrganizationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	orgIdString := d.Id()
	orgId := helper.StrToUInt64(orgIdString)
	if d.HasChange("root_node_name") {
		if v, ok := d.GetOk("root_node_name"); ok {
			service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
			innerErr := service.UpdateOrganizationRootNodeName(ctx, orgId, v.(string))
			if innerErr != nil {
				return innerErr
			}
		}
	}

	return resourceTencentCloudOrganizationOrganizationRead(d, meta)
}

func resourceTencentCloudOrganizationOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	if err := service.DeleteOrganizationOrganizationById(ctx); err != nil {
		return err
	}

	return nil
}
