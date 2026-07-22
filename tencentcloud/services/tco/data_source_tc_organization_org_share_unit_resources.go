package tco

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudOrganizationOrgShareUnitResources() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOrganizationOrgShareUnitResourcesRead,
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
				Description: "Search 对于 keywords. Support product 资源 ID search。",
			},

			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Shared 资源类型",
			},

			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Shared 单位 资源 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Shared 资源 ID。",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Shared 资源类型",
						},
						"create_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "创建时间。",
						},
						"product_resource_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Product 资源 ID。",
						},
						"shared_member_num": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "数量 shared 单位 members。",
						},
						"shared_member_use_num": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "数量 shared 单位 members 在 使用。",
						},
						"share_manager_uin": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Sharing administrator OwnerUin。",
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

func dataSourceTencentCloudOrganizationOrgShareUnitResourcesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_organization_org_share_unit_resources.read")()
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

	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	var respData []*organization.ShareUnitResource
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOrganizationOrgShareUnitResourcesByFilter(ctx, paramMap)
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

			var resourceId string
			if items.ResourceId != nil {
				itemsMap["resource_id"] = items.ResourceId
				resourceId = *items.ResourceId
			}

			if items.Type != nil {
				itemsMap["type"] = items.Type
			}

			if items.CreateTime != nil {
				itemsMap["create_time"] = items.CreateTime
			}

			if items.ProductResourceId != nil {
				itemsMap["product_resource_id"] = items.ProductResourceId
			}

			if items.SharedMemberNum != nil {
				itemsMap["shared_member_num"] = items.SharedMemberNum
			}

			if items.SharedMemberUseNum != nil {
				itemsMap["shared_member_use_num"] = items.SharedMemberUseNum
			}

			if items.ShareManagerUin != nil {
				itemsMap["share_manager_uin"] = items.ShareManagerUin
			}

			ids = append(ids, resourceId)
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
