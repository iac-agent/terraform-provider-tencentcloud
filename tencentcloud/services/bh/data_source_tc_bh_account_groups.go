package bh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bhv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudBhAccountGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudBhAccountGroupsRead,
		Schema: map[string]*schema.Schema{
			"deep_in": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "是否recursively 查询，0 对于 non-recursive，1 对于 recursive。",
			},

			"parent_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Parent 账号 组 ID，默认值 0，查询 all groups under root 账号 组。",
			},

			"group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "账号 组名称，fuzzy 查询。",
			},

			"page_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Get 数据 从 其中 页面。",
			},

			"account_group_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "账号 组 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "账号 组 ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "账号 组名称",
						},
						"id_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "账号 组 ID 路径",
						},
						"name_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "账号 组名称 路径",
						},
						"parent_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Parent 账号 组 ID",
						},
						"source": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "账号 组 来源",
						},
						"user_total": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total 数量 users under 账号 组。",
						},
						"is_leaf": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否为a leaf 节点。",
						},
						"import_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "账号 组 import 类型",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "账号 组 描述",
						},
						"parent_org_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parent 来源 账号 organization ID. 当 使用 third-party import 用户 sources，记录 组 ID 此 组 在 来源 organization structure。",
						},
						"org_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "来源 账号 organization ID. 当 使用 third-party import 用户 sources，记录 组 ID 此 组 在 来源 organization structure。",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否account 组 has been connected，0 表示 不 connected，1 表示 connected。",
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

func dataSourceTencentCloudBhAccountGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_bh_account_groups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = BhService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOkExists("deep_in"); ok {
		paramMap["DeepIn"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("parent_id"); ok {
		paramMap["ParentId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("group_name"); ok {
		paramMap["GroupName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("page_num"); ok {
		paramMap["PageNum"] = helper.IntInt64(v.(int))
	}

	var respData []*bhv20230418.AccountGroup
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeBhAccountGroupsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	accountGroupSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, accountGroupSet := range respData {
			accountGroupSetMap := map[string]interface{}{}
			if accountGroupSet.Id != nil {
				accountGroupSetMap["id"] = accountGroupSet.Id
			}

			if accountGroupSet.Name != nil {
				accountGroupSetMap["name"] = accountGroupSet.Name
			}

			if accountGroupSet.IdPath != nil {
				accountGroupSetMap["id_path"] = accountGroupSet.IdPath
			}

			if accountGroupSet.NamePath != nil {
				accountGroupSetMap["name_path"] = accountGroupSet.NamePath
			}

			if accountGroupSet.ParentId != nil {
				accountGroupSetMap["parent_id"] = accountGroupSet.ParentId
			}

			if accountGroupSet.Source != nil {
				accountGroupSetMap["source"] = accountGroupSet.Source
			}

			if accountGroupSet.UserTotal != nil {
				accountGroupSetMap["user_total"] = accountGroupSet.UserTotal
			}

			if accountGroupSet.IsLeaf != nil {
				accountGroupSetMap["is_leaf"] = accountGroupSet.IsLeaf
			}

			if accountGroupSet.ImportType != nil {
				accountGroupSetMap["import_type"] = accountGroupSet.ImportType
			}

			if accountGroupSet.Description != nil {
				accountGroupSetMap["description"] = accountGroupSet.Description
			}

			if accountGroupSet.ParentOrgId != nil {
				accountGroupSetMap["parent_org_id"] = accountGroupSet.ParentOrgId
			}

			if accountGroupSet.OrgId != nil {
				accountGroupSetMap["org_id"] = accountGroupSet.OrgId
			}

			if accountGroupSet.Status != nil {
				accountGroupSetMap["status"] = accountGroupSet.Status
			}

			accountGroupSetList = append(accountGroupSetList, accountGroupSetMap)
		}

		_ = d.Set("account_group_set", accountGroupSetList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), accountGroupSetList); e != nil {
			return e
		}
	}

	return nil
}
