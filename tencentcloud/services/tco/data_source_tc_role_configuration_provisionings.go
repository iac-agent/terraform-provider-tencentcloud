package tco

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudRoleConfigurationProvisionings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRoleConfigurationProvisioningsRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Space ID。",
			},

			"role_configuration_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "权限配置 ID",
			},

			"target_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "类型 synchronized target 账号 of the Tencent Cloud Organization. ManagerUin: admin 账号; MemberUin: member 账号",
			},

			"target_uin": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "UIN of the synchronized target 账号 of the Tencent Cloud Organization。",
			},

			"deployment_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Deployed: Deployment succeeded; DeployedRequired: Redeployment 必填; DeployFailed: Deployment failed。",
			},

			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search by configuration 名称 is supported。",
			},

			"role_configuration_provisionings": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Department member 账号 list。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deployment_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Deployed: Deployment succeeded; DeployedRequired: Redeployment 必填; DeployFailed: Deployment failed。",
						},
						"role_configuration_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "权限配置 ID",
						},
						"role_configuration_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Permission configuration 名称",
						},
						"target_uin": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "UIN of the target 账号 of the Tencent Cloud Organization。",
						},
						"target_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "名称 target 账号 of the Tencent Cloud Organization。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "创建时间。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "修改时间。",
						},
						"target_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "类型 synchronized target 账号 of the Tencent Cloud Organization. ManagerUin: admin 账号; MemberUin: member 账号",
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

func dataSourceTencentCloudRoleConfigurationProvisioningsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_role_configuration_provisionings.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		paramMap["ZoneId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("role_configuration_id"); ok {
		paramMap["RoleConfigurationId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_type"); ok {
		paramMap["TargetType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("target_uin"); ok {
		paramMap["TargetUin"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("deployment_status"); ok {
		paramMap["DeploymentStatus"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter"); ok {
		paramMap["Filter"] = helper.String(v.(string))
	}

	var roleConfigurationProvisionings []*organization.RoleConfigurationProvisionings
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRoleConfigurationProvisioningsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		roleConfigurationProvisionings = result
		return nil
	})
	if err != nil {
		return err
	}

	var ids []string
	roleConfigurationProvisioningsList := make([]map[string]interface{}, 0, len(roleConfigurationProvisionings))
	if roleConfigurationProvisionings != nil {
		for _, roleConfigurationProvisionings := range roleConfigurationProvisionings {
			roleConfigurationProvisioningsMap := map[string]interface{}{}

			var roleConfigurationId string
			if roleConfigurationProvisionings.DeploymentStatus != nil {
				roleConfigurationProvisioningsMap["deployment_status"] = roleConfigurationProvisionings.DeploymentStatus
			}

			if roleConfigurationProvisionings.RoleConfigurationId != nil {
				roleConfigurationProvisioningsMap["role_configuration_id"] = roleConfigurationProvisionings.RoleConfigurationId
				roleConfigurationId = *roleConfigurationProvisionings.RoleConfigurationId
			}

			if roleConfigurationProvisionings.RoleConfigurationName != nil {
				roleConfigurationProvisioningsMap["role_configuration_name"] = roleConfigurationProvisionings.RoleConfigurationName
			}

			if roleConfigurationProvisionings.TargetUin != nil {
				roleConfigurationProvisioningsMap["target_uin"] = roleConfigurationProvisionings.TargetUin
			}

			if roleConfigurationProvisionings.TargetName != nil {
				roleConfigurationProvisioningsMap["target_name"] = roleConfigurationProvisionings.TargetName
			}

			if roleConfigurationProvisionings.CreateTime != nil {
				roleConfigurationProvisioningsMap["create_time"] = roleConfigurationProvisionings.CreateTime
			}

			if roleConfigurationProvisionings.UpdateTime != nil {
				roleConfigurationProvisioningsMap["update_time"] = roleConfigurationProvisionings.UpdateTime
			}

			if roleConfigurationProvisionings.TargetType != nil {
				roleConfigurationProvisioningsMap["target_type"] = roleConfigurationProvisionings.TargetType
			}

			ids = append(ids, roleConfigurationId)
			roleConfigurationProvisioningsList = append(roleConfigurationProvisioningsList, roleConfigurationProvisioningsMap)
		}

		_ = d.Set("role_configuration_provisionings", roleConfigurationProvisioningsList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
