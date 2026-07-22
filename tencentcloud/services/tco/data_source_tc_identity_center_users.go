package tco

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudIdentityCenterUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudIdentityCenterUsersRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Space ID。",
			},

			"user_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用户 状态: 已启用，已禁用",
			},

			"user_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用户 类型 Manual: manually 创建; Synchronized: externally imported。",
			},

			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "过滤器 criterion，其中 currently 仅 支持 用户名，email 地址，userId，和 描述",
			},

			"filter_groups": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Filtered 用户 组. IsSelected=1 将 是 返回 对于 sub-用户 associated 使用 此 用户 组。",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"sort_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sorting 字段，其中 currently 仅 支持 CreateTime. 默认为 CreateTime 字段。",
			},

			"sort_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sorting 类型 Desc: 降序; Asc: 升序 It should 是 集合 along 使用 SortField。",
			},

			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "用户 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Queried 用户名",
						},
						"first_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "First 名称 用户",
						},
						"last_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Last 名称 用户",
						},
						"display_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Display 名称 用户",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "用户 描述",
						},
						"email": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Email 地址 的 用户，其中 必须 是 唯一 within directory。",
						},
						"user_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "用户 状态 有效值：已启用，已禁用",
						},
						"user_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "用户 类型 Manual: manually 创建; Synchronized: externally imported。",
						},
						"user_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "用户 ID。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "创建时间 的 用户",
						},
						"update_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "修改时间 的 用户",
						},
						"is_selected": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether selected。",
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

func dataSourceTencentCloudIdentityCenterUsersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_identity_center_users.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		paramMap["ZoneId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_status"); ok {
		paramMap["UserStatus"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_type"); ok {
		paramMap["UserType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter"); ok {
		paramMap["Filter"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter_groups"); ok {
		filterGroupsList := []*string{}
		filterGroupsSet := v.(*schema.Set).List()
		for i := range filterGroupsSet {
			filterGroups := filterGroupsSet[i].(string)
			filterGroupsList = append(filterGroupsList, helper.String(filterGroups))
		}
		paramMap["FilterGroups"] = filterGroupsList
	}

	if v, ok := d.GetOk("sort_field"); ok {
		paramMap["SortField"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_type"); ok {
		paramMap["SortType"] = helper.String(v.(string))
	}

	var users []*organization.UserInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeIdentityCenterUsersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		users = result
		return nil
	})
	if err != nil {
		return err
	}

	usersList := make([]map[string]interface{}, 0, len(users))
	ids := make([]string, 0, len(users))
	for _, user := range users {
		usersMap := map[string]interface{}{}

		if user.UserName != nil {
			usersMap["user_name"] = user.UserName
		}

		if user.FirstName != nil {
			usersMap["first_name"] = user.FirstName
		}

		if user.LastName != nil {
			usersMap["last_name"] = user.LastName
		}

		if user.DisplayName != nil {
			usersMap["display_name"] = user.DisplayName
		}

		if user.Description != nil {
			usersMap["description"] = user.Description
		}

		if user.Email != nil {
			usersMap["email"] = user.Email
		}

		if user.UserStatus != nil {
			usersMap["user_status"] = user.UserStatus
		}

		if user.UserType != nil {
			usersMap["user_type"] = user.UserType
		}

		if user.UserId != nil {
			usersMap["user_id"] = user.UserId
			ids = append(ids, *user.UserId)
		}

		if user.CreateTime != nil {
			usersMap["create_time"] = user.CreateTime
		}

		if user.UpdateTime != nil {
			usersMap["update_time"] = user.UpdateTime
		}

		if user.IsSelected != nil {
			usersMap["is_selected"] = user.IsSelected
		}

		usersList = append(usersList, usersMap)

		_ = d.Set("users", usersList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), usersList); e != nil {
			return e
		}
	}

	return nil
}
