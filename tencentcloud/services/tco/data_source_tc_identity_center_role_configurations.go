package tco

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudIdentityCenterRoleConfigurations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudIdentityCenterRoleConfigurationsRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Space ID。",
			},

			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "过滤器 criteria, 其中 是 case insensitive. Currently, 仅 RoleConfigurationName 是 支持 和 仅 eq (Equals) 和 sw (Start With) 是 支持. Example: 过滤器 = \"RoleConfigurationName, 仅 sw 测试\" 表示 querying all 权限 configurations starting 使用 测试. 过滤器 = \"RoleConfigurationName, 仅 eq TestRoleConfiguration\" 表示 querying 权限 配置 named TestRoleConfiguration.",
			},

			"filter_targets": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Check 是否member 账号 has been 已配置 使用 permissions. 如果 已配置，返回 IsSelected: true; otherwise，返回 false。",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"principal_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "UserId 的 authorized 用户 或 GroupId 的 authorized 用户 组，其中 必须 是 集合 together 使用 input 参数 FilterTargets。",
			},

			"role_configurations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Permission 配置 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_configuration_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "权限配置 ID",
						},
						"role_configuration_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Permission 配置 名称",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Permission 配置 描述",
						},
						"session_duration": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Session 时长. It 表示maximum 会话 时长 当 CIC users 使用 访问 配置 到 访问 member accounts.\n单位：秒。",
						},
						"relay_state": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Initial 访问 页面. It 表示initial 访问 页面 URL 当 CIC users 使用 访问 配置 到 访问 member accounts。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "创建时间 的 权限 配置。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "更新时间 的 权限 配置。",
						},
						"is_selected": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "如果 input 参数 FilterTargets 是 提供，check 是否member 账号 has been 已配置 使用 permissions. 如果 已配置，返回 true; otherwise，返回 false。",
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

func dataSourceTencentCloudIdentityCenterRoleConfigurationsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_identity_center_role_configurations.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		paramMap["ZoneId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter"); ok {
		paramMap["Filter"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter_targets"); ok {
		filterTargetsList := []*int64{}
		filterTargetsSet := v.(*schema.Set).List()
		for i := range filterTargetsSet {
			filterTargets := filterTargetsSet[i].(int)
			filterTargetsList = append(filterTargetsList, helper.IntInt64(filterTargets))
		}
		paramMap["FilterTargets"] = filterTargetsList
	}

	if v, ok := d.GetOk("principal_id"); ok {
		paramMap["PrincipalId"] = helper.String(v.(string))
	}

	var roleConfigurations []*organization.RoleConfiguration
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeIdentityCenterRoleConfigurationsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		roleConfigurations = result
		return nil
	})
	if err != nil {
		return err
	}

	roleConfigurationsList := make([]map[string]interface{}, 0, len(roleConfigurations))
	ids := make([]string, 0, len(roleConfigurations))
	for _, roleConfiguration := range roleConfigurations {
		roleConfigurationsMap := map[string]interface{}{}

		if roleConfiguration.RoleConfigurationId != nil {
			roleConfigurationsMap["role_configuration_id"] = roleConfiguration.RoleConfigurationId
			ids = append(ids, *roleConfiguration.RoleConfigurationId)
		}

		if roleConfiguration.RoleConfigurationName != nil {
			roleConfigurationsMap["role_configuration_name"] = roleConfiguration.RoleConfigurationName
		}

		if roleConfiguration.Description != nil {
			roleConfigurationsMap["description"] = roleConfiguration.Description
		}

		if roleConfiguration.SessionDuration != nil {
			roleConfigurationsMap["session_duration"] = roleConfiguration.SessionDuration
		}

		if roleConfiguration.RelayState != nil {
			roleConfigurationsMap["relay_state"] = roleConfiguration.RelayState
		}

		if roleConfiguration.CreateTime != nil {
			roleConfigurationsMap["create_time"] = roleConfiguration.CreateTime
		}

		if roleConfiguration.UpdateTime != nil {
			roleConfigurationsMap["update_time"] = roleConfiguration.UpdateTime
		}

		if roleConfiguration.IsSelected != nil {
			roleConfigurationsMap["is_selected"] = roleConfiguration.IsSelected
		}

		roleConfigurationsList = append(roleConfigurationsList, roleConfigurationsMap)
	}

	_ = d.Set("role_configurations", roleConfigurationsList)

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), roleConfigurationsList); e != nil {
			return e
		}
	}

	return nil
}
