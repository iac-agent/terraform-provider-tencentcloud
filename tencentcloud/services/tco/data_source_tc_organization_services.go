package tco

import (
	"context"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudOrganizationServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOrganizationServicesRead,
		Schema: map[string]*schema.Schema{
			"search_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Keyword 对于 search 通过 名称",
			},
			// computed
			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Organization 服务 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Organization 服务 ID 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"product_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Organization 服务 product 名称 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"is_assign": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否support delegation. 有效值：1 (yes)，2 (无). 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Organization 服务 描述 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"member_num": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "数量 当前 delegated admins. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"document": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Help documentation. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"console_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Console 路径 的 organization 服务 product. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"is_usage_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否access usage 状态 有效值：1 (yes)，2 (无). 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"can_assign_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "限制 对于 数量 delegated admins. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"product": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Organization 服务 product identifier. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"service_grant": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否support organization 服务 authorization. 有效值：1 (yes)，2 (无). 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"grant_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enabling 状态 organization 服务 authorization. 此 字段 是 有效 当 ServiceGrant 是 1. 有效值：已启用，已禁用 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"is_set_management_scope": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否support setting delegated management 范围 有效值：1 (yes)，2 (无).\n注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudOrganizationServicesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_organization_services.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		items   []*organization.OrganizationServiceAssign
	)

	paramMap := make(map[string]interface{})

	if v, ok := d.GetOk("search_key"); ok {
		paramMap["SearchKey"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOrganizationServicesByFilter(ctx, paramMap)
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
		for _, item := range items {
			orgServiceAssignMap := map[string]interface{}{}

			if item.ServiceId != nil {
				orgServiceAssignMap["service_id"] = item.ServiceId
				serviceIdStr := strconv.FormatUint(*item.ServiceId, 10)
				ids = append(ids, serviceIdStr)
			}

			if item.ProductName != nil {
				orgServiceAssignMap["product_name"] = item.ProductName
			}

			if item.IsAssign != nil {
				orgServiceAssignMap["is_assign"] = item.IsAssign
			}

			if item.Description != nil {
				orgServiceAssignMap["description"] = item.Description
			}

			if item.MemberNum != nil {
				orgServiceAssignMap["member_num"] = item.MemberNum
			}

			if item.ConsoleUrl != nil {
				orgServiceAssignMap["console_url"] = item.ConsoleUrl
			}

			if item.IsUsageStatus != nil {
				orgServiceAssignMap["is_usage_status"] = item.IsUsageStatus
			}

			if item.CanAssignCount != nil {
				orgServiceAssignMap["can_assign_count"] = item.CanAssignCount
			}

			if item.Product != nil {
				orgServiceAssignMap["product"] = item.Product
			}

			if item.ServiceGrant != nil {
				orgServiceAssignMap["service_grant"] = item.ServiceGrant
			}

			if item.GrantStatus != nil {
				orgServiceAssignMap["grant_status"] = item.GrantStatus
			}

			if item.IsSetManagementScope != nil {
				orgServiceAssignMap["is_set_management_scope"] = item.IsSetManagementScope
			}

			tmpList = append(tmpList, orgServiceAssignMap)
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
