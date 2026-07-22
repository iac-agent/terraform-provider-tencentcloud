package tco

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudIdentityCenterGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudIdentityCenterGroupsRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Space ID。",
			},

			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "过滤器 criterion. Format: <Attribute> <Operator> <Value>, case-insensitive. Currently, <Attribute> 支持 仅 GroupName, 和 <Operator> 支持 仅 eq (Equals) 和 sw (Start With). For 示例, 过滤器 = \"GroupName sw 测试\" indicates querying all 用户 groups 使用 names starting 使用 测试; 过滤器 = \"GroupName eq testgroup\" indicates querying 用户 组 使用 名称 testgroup.",
			},

			"group_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用户 组 类型 Manual: manually 创建; Synchronized: externally imported。",
			},

			"filter_users": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Filtered 用户 IsSelected=1 将 是 返回 对于 用户 组 associated 使用 此 用户",
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

			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "用户 组 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "用户 组名称",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "用户 组 描述",
						},
						"create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "创建时间 的 用户 组。",
						},
						"group_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "用户 组 类型 Manual: manually 创建; Synchronized: externally imported。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "修改时间 的 用户 组。",
						},
						"group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "用户 组 ID",
						},
						"member_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "数量 组 members。",
						},
						"is_selected": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "如果 input 参数 FilterUsers 是 提供，返回 true 当 用户 是 在 用户 组; otherwise，返回 false。",
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

func dataSourceTencentCloudIdentityCenterGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_identity_center_groups.read")()
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

	if v, ok := d.GetOk("group_type"); ok {
		paramMap["GroupType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter_users"); ok {
		filterUsersList := []*string{}
		filterUsersSet := v.(*schema.Set).List()
		for i := range filterUsersSet {
			filterUsers := filterUsersSet[i].(string)
			filterUsersList = append(filterUsersList, helper.String(filterUsers))
		}
		paramMap["FilterUsers"] = filterUsersList
	}

	if v, ok := d.GetOk("sort_field"); ok {
		paramMap["SortField"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_type"); ok {
		paramMap["SortType"] = helper.String(v.(string))
	}

	var groups []*organization.GroupInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeIdentityCenterGroupsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		groups = result
		return nil
	})
	if err != nil {
		return err
	}

	groupsList := make([]map[string]interface{}, 0, len(groups))
	ids := make([]string, 0, len(groups))
	for _, group := range groups {
		groupsMap := map[string]interface{}{}

		if group.GroupName != nil {
			groupsMap["group_name"] = group.GroupName
		}

		if group.Description != nil {
			groupsMap["description"] = group.Description
		}

		if group.CreateTime != nil {
			groupsMap["create_time"] = group.CreateTime
		}

		if group.GroupType != nil {
			groupsMap["group_type"] = group.GroupType
		}

		if group.UpdateTime != nil {
			groupsMap["update_time"] = group.UpdateTime
		}

		if group.GroupId != nil {
			groupsMap["group_id"] = group.GroupId
			ids = append(ids, *group.GroupId)
		}

		if group.MemberCount != nil {
			groupsMap["member_count"] = group.MemberCount
		}

		if group.IsSelected != nil {
			groupsMap["is_selected"] = group.IsSelected
		}

		groupsList = append(groupsList, groupsMap)
	}

	_ = d.Set("groups", groupsList)

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), groupsList); e != nil {
			return e
		}
	}

	return nil
}
