package tco

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudOrganizationOrgShareUnitMembers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOrganizationOrgShareUnitMembersRead,
		Schema: map[string]*schema.Schema{
			"unit_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "共享单元 ID",
			},

			"area": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Shared 单位 area。",
			},

			"search_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search 对于 keywords. Support member Uin searches。",
			},

			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Shared 单位 member 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"share_member_uin": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Shared member Uin。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "创建时间。",
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

func dataSourceTencentCloudOrganizationOrgShareUnitMembersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_organization_org_share_unit_members.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("unit_id"); ok {
		paramMap["UnitId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("area"); ok {
		paramMap["Area"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_key"); ok {
		paramMap["SearchKey"] = helper.String(v.(string))
	}

	var respData []*organization.ShareUnitMember
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOrganizationOrgShareUnitMembersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	var ids []string
	itemsList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, items := range respData {
			itemsMap := map[string]interface{}{}

			var shareMemberUin int64
			if items.ShareMemberUin != nil {
				itemsMap["share_member_uin"] = items.ShareMemberUin
				shareMemberUin = *items.ShareMemberUin
			}

			if items.CreateTime != nil {
				itemsMap["create_time"] = items.CreateTime
			}

			ids = append(ids, helper.Int64ToStr(shareMemberUin))
			itemsList = append(itemsList, itemsMap)
		}

		_ = d.Set("items", itemsList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), itemsList); e != nil {
			return e
		}
	}

	return nil
}
