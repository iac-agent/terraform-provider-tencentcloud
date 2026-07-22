package dlc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDlcDescribeWorkGroupInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcDescribeWorkGroupInfoRead,
		Schema: map[string]*schema.Schema{
			"work_group_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Working 组 ID",
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Types 的 queried 信息. 用户: 用户 信息; DataAuth: 数据 permissions; EngineAuth: 引擎 permissions。",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 criteria 该 是 queriedWhen 类型 是 用户， fuzzy search 是 支持 作为 键 是 用户-名称When 类型 是 DataAuth， keys 支持 是:策略-类型: types 的 permissions;策略-来源: 数据 sources;数据-名称: fuzzy search 的 数据库 和 表.当 类型 是 EngineAuth， keys 支持 是:策略-类型: types 的 permissions;策略-来源: 数据 sources;引擎-名称: fuzzy search 的 数据库 和 表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Attribute 名称 如果 more 比 一个 过滤器 exists， logical relationship between these filters 是 `OR`。",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Attribute 值 如果 多个 值 exist 在 一个 过滤器， logical relationship between these 值 是 `OR`。",
						},
					},
				},
			},

			"sort_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort 字段.当 类型 是 用户，create-时间 和 用户-名称 是 支持.当 类型 是 DataAuth，create-时间 是 支持.当 类型 是 EngineAuth，create-时间 是 支持。",
			},

			"sorting": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting methods: desc 表示 在 顺序; asc 表示 在 reverse 顺序; 它 是 asc 通过 默认值。",
			},

			"work_group_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Details about working groups注意：此字段可能返回 null，表示无法获取有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"work_group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Working 组 ID注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"work_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Working 组 name注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Types 的 信息 included. 用户: 用户 信息; DataAuth: 数据 permissions; EngineAuth: 引擎 permissions注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"user_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection 的 users bound 到 working groups注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Collection 的 用户 information注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "用户 ID 其中 matches sub-用户 UIN 在 CAM side。",
												},
												"user_description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "用户 descriptionNote: 返回 值 的 此 字段 可能 是 null，indicating 该 无 有效 值 是 获取。",
												},
												"creator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "创建者 的 当前 用户",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "创建时间 的 当前 用户，e.g. 16:19:32，July 28，2021。",
												},
												"user_alias": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "用户 alias。",
												},
											},
										},
									},
									"total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total users注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"data_policy_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection 的 数据 permissions注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Collection 的 policies注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 目标 数据库. `*` 表示 all databases 在 当前 catalog. To grant admin permissions，它 必须 是 `*`; 到 grant 数据 连接 permissions，它 必须 是 null; 到 grant other permissions，它 可以 是 any 数据库。",
												},
												"catalog": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 目标 数据 来源 To grant admin 权限，它 必须 是 `*` (all resources 在 此 级别); 到 grant 数据 来源 和 数据库 permissions，它 必须 是 `COSDataCatalog` 或 `*`; 到 grant 表 permissions，它 可以 是 自定义 数据 来源; 如果 它 是 left 空，`DataLakeCatalog` 是 使用. 注意: To grant permissions 在 自定义 数据 来源， permissions 该 可以 是 managed 在 Data Lake Compute console 是 subsets 的 账号 permissions granted 当 您 connect 数据 来源 到 console。",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 目标 表. `*` 表示 all tables 在 当前 数据库. To grant admin permissions，它 必须 是 `*`; 到 grant 数据 连接 和 数据库 permissions，它 必须 是 null; 到 grant other permissions，它 可以 是 any 表。",
												},
												"operation": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "目标 permissions，其中 vary 通过 权限 级别 Admin: `ALL` (默认值); 数据 连接: `CREATE`; 数据库: `ALL`，`CREATE`，`ALTER`，和 `DROP`; 表: `ALL`，`SELECT`，`INSERT`，`ALTER`，`DELETE`，`DROP`，和 `UPDATE`. 注意: For 表 permissions，如果 数据 来源 other 比 `COSDataCatalog` 是 指定，仅 `SELECT` 权限 可以 是 granted here。",
												},
												"policy_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "权限 类型 有效值：`ADMIN`，`DATASOURCE`，`DATABASE`，`TABLE`，`VIEW`，`FUNCTION`，`COLUMN`，和 `ENGINE`. 注意: 如果 它 是 left 空，`ADMIN` 是 使用。",
												},
												"function": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 目标 函数. `*` 表示 all functions 在 当前 catalog. To grant admin permissions，它 必须 是 `*`; 到 grant 数据 连接 permissions，它 必须 是 null; 到 grant other permissions，它 可以 是 any 函数.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"view": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 目标 view. `*` 表示 all views 在 当前 数据库. To grant admin permissions，它 必须 是 `*`; 到 grant 数据 连接 和 数据库 permissions，它 必须 是 null; 到 grant other permissions，它 可以 是 any view.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"column": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 目标 列. `*` 表示 all columns. To grant admin permissions，它 必须 是 `*`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"data_engine": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 目标 数据 引擎. `*` 表示 all engines. To grant admin permissions，它 必须 是 `*`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"re_auth": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否grantee 是 allowed 到 further grant permissions. 有效值：`false` (默认值) 和 `true` ( grantee 可以 grant permissions gained here 到 other sub-users).注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"source": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "权限 来源，其中 不是必填项 当 input 参数 是 passed 在. 有效值：`USER` (从 用户) 和 `WORKGROUP` (从 一个 或 more associated work groups).注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "grant 模式，其中 不是必填项 作为 input 参数. 有效值：`COMMON` 和 `SENIOR`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "操作者，其中 不是必填项 作为 input 参数.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "权限 策略 创建时间，其中 不是必填项 作为 input 参数.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"source_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "ID work 组，其中 applies 仅 当 值 的 `来源` 字段 是 `WORKGROUP`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"source_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 work 组，其中 applies 仅 当 值 的 `来源` 字段 是 `WORKGROUP`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "策略 ID.注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total policies注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"engine_policy_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection 的 引擎 permissions注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Collection 的 policies注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 目标 数据库. `*` 表示 all databases 在 当前 catalog. To grant admin permissions，它 必须 是 `*`; 到 grant 数据 连接 permissions，它 必须 是 null; 到 grant other permissions，它 可以 是 any 数据库。",
												},
												"catalog": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 目标 数据 来源 To grant admin 权限，它 必须 是 `*` (all resources 在 此 级别); 到 grant 数据 来源 和 数据库 permissions，它 必须 是 `COSDataCatalog` 或 `*`; 到 grant 表 permissions，它 可以 是 自定义 数据 来源; 如果 它 是 left 空，`DataLakeCatalog` 是 使用. 注意: To grant permissions 在 自定义 数据 来源， permissions 该 可以 是 managed 在 Data Lake Compute console 是 subsets 的 账号 permissions granted 当 您 connect 数据 来源 到 console。",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 目标 表. `*` 表示 all tables 在 当前 数据库. To grant admin permissions，它 必须 是 `*`; 到 grant 数据 连接 和 数据库 permissions，它 必须 是 null; 到 grant other permissions，它 可以 是 any 表。",
												},
												"operation": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "目标 permissions，其中 vary 通过 权限 级别 Admin: `ALL` (默认值); 数据 连接: `CREATE`; 数据库: `ALL`，`CREATE`，`ALTER`，和 `DROP`; 表: `ALL`，`SELECT`，`INSERT`，`ALTER`，`DELETE`，`DROP`，和 `UPDATE`. 注意: For 表 permissions，如果 数据 来源 other 比 `COSDataCatalog` 是 指定，仅 `SELECT` 权限 可以 是 granted here。",
												},
												"policy_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "权限 类型 有效值：`ADMIN`，`DATASOURCE`，`DATABASE`，`TABLE`，`VIEW`，`FUNCTION`，`COLUMN`，和 `ENGINE`. 注意: 如果 它 是 left 空，`ADMIN` 是 使用。",
												},
												"function": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 目标 函数. `*` 表示 all functions 在 当前 catalog. To grant admin permissions，它 必须 是 `*`; 到 grant 数据 连接 permissions，它 必须 是 null; 到 grant other permissions，它 可以 是 any 函数.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"view": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 目标 view. `*` 表示 all views 在 当前 数据库. To grant admin permissions，它 必须 是 `*`; 到 grant 数据 连接 和 数据库 permissions，它 必须 是 null; 到 grant other permissions，它 可以 是 any view.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"column": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 目标 列. `*` 表示 all columns. To grant admin permissions，它 必须 是 `*`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"data_engine": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 目标 数据 引擎. `*` 表示 all engines. To grant admin permissions，它 必须 是 `*`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"re_auth": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否grantee 是 allowed 到 further grant permissions. 有效值：`false` (默认值) 和 `true` ( grantee 可以 grant permissions gained here 到 other sub-users).注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"source": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "权限 来源，其中 不是必填项 当 input 参数 是 passed 在. 有效值：`USER` (从 用户) 和 `WORKGROUP` (从 一个 或 more associated work groups).注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "grant 模式，其中 不是必填项 作为 input 参数. 有效值：`COMMON` 和 `SENIOR`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "操作者，其中 不是必填项 作为 input 参数.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "权限 策略 创建时间，其中 不是必填项 作为 input 参数.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"source_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "ID work 组，其中 applies 仅 当 值 的 `来源` 字段 是 `WORKGROUP`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"source_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 work 组，其中 applies 仅 当 值 的 `来源` 字段 是 `WORKGROUP`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "策略 ID.注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total policies注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"work_group_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Working 组 description注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"row_filter_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection 的 信息 about filtered rows注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Collection 的 policies注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 目标 数据库. `*` 表示 all databases 在 当前 catalog. To grant admin permissions，它 必须 是 `*`; 到 grant 数据 连接 permissions，它 必须 是 null; 到 grant other permissions，它 可以 是 any 数据库。",
												},
												"catalog": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 目标 数据 来源 To grant admin 权限，它 必须 是 `*` (all resources 在 此 级别); 到 grant 数据 来源 和 数据库 permissions，它 必须 是 `COSDataCatalog` 或 `*`; 到 grant 表 permissions，它 可以 是 自定义 数据 来源; 如果 它 是 left 空，`DataLakeCatalog` 是 使用. 注意: To grant permissions 在 自定义 数据 来源， permissions 该 可以 是 managed 在 Data Lake Compute console 是 subsets 的 账号 permissions granted 当 您 connect 数据 来源 到 console。",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 目标 表. `*` 表示 all tables 在 当前 数据库. To grant admin permissions，它 必须 是 `*`; 到 grant 数据 连接 和 数据库 permissions，它 必须 是 null; 到 grant other permissions，它 可以 是 any 表。",
												},
												"operation": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "目标 permissions，其中 vary 通过 权限 级别 Admin: `ALL` (默认值); 数据 连接: `CREATE`; 数据库: `ALL`，`CREATE`，`ALTER`，和 `DROP`; 表: `ALL`，`SELECT`，`INSERT`，`ALTER`，`DELETE`，`DROP`，和 `UPDATE`. 注意: For 表 permissions，如果 数据 来源 other 比 `COSDataCatalog` 是 指定，仅 `SELECT` 权限 可以 是 granted here。",
												},
												"policy_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "权限 类型 有效值：`ADMIN`，`DATASOURCE`，`DATABASE`，`TABLE`，`VIEW`，`FUNCTION`，`COLUMN`，和 `ENGINE`. 注意: 如果 它 是 left 空，`ADMIN` 是 使用。",
												},
												"function": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 目标 函数. `*` 表示 all functions 在 当前 catalog. To grant admin permissions，它 必须 是 `*`; 到 grant 数据 连接 permissions，它 必须 是 null; 到 grant other permissions，它 可以 是 any 函数.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"view": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 目标 view. `*` 表示 all views 在 当前 数据库. To grant admin permissions，它 必须 是 `*`; 到 grant 数据 连接 和 数据库 permissions，它 必须 是 null; 到 grant other permissions，它 可以 是 any view.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"column": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 目标 列. `*` 表示 all columns. To grant admin permissions，它 必须 是 `*`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"data_engine": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 目标 数据 引擎. `*` 表示 all engines. To grant admin permissions，它 必须 是 `*`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"re_auth": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否grantee 是 allowed 到 further grant permissions. 有效值：`false` (默认值) 和 `true` ( grantee 可以 grant permissions gained here 到 other sub-users).注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"source": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "权限 来源，其中 不是必填项 当 input 参数 是 passed 在. 有效值：`USER` (从 用户) 和 `WORKGROUP` (从 一个 或 more associated work groups).注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "grant 模式，其中 不是必填项 作为 input 参数. 有效值：`COMMON` 和 `SENIOR`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "操作者，其中 不是必填项 作为 input 参数.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "权限 策略 创建时间，其中 不是必填项 作为 input 参数.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"source_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "ID work 组，其中 applies 仅 当 值 的 `来源` 字段 是 `WORKGROUP`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"source_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "名称 work 组，其中 applies 仅 当 值 的 `来源` 字段 是 `WORKGROUP`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "策略 ID.注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total policies注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
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

func dataSourceTencentCloudDlcDescribeWorkGroupInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_describe_work_group_info.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOkExists("work_group_id"); v != nil {
		paramMap["WorkGroupId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*dlc.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := dlc.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["Filters"] = tmpSet
	}

	if v, ok := d.GetOk("sort_by"); ok {
		paramMap["SortBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sorting"); ok {
		paramMap["Sorting"] = helper.String(v.(string))
	}

	service := DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var workGroupInfo *dlc.WorkGroupDetailInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcDescribeWorkGroupInfoByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		workGroupInfo = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	workGroupDetailInfoMap := map[string]interface{}{}

	if workGroupInfo != nil {

		if workGroupInfo.WorkGroupId != nil {
			workGroupDetailInfoMap["work_group_id"] = workGroupInfo.WorkGroupId
		}

		if workGroupInfo.WorkGroupName != nil {
			workGroupDetailInfoMap["work_group_name"] = workGroupInfo.WorkGroupName
		}

		if workGroupInfo.Type != nil {
			workGroupDetailInfoMap["type"] = workGroupInfo.Type
		}

		if workGroupInfo.UserInfo != nil {
			userInfoMap := map[string]interface{}{}

			if workGroupInfo.UserInfo.UserSet != nil {
				var userSetList []interface{}
				for _, userSet := range workGroupInfo.UserInfo.UserSet {
					userSetMap := map[string]interface{}{}

					if userSet.UserId != nil {
						userSetMap["user_id"] = userSet.UserId
					}

					if userSet.UserDescription != nil {
						userSetMap["user_description"] = userSet.UserDescription
					}

					if userSet.Creator != nil {
						userSetMap["creator"] = userSet.Creator
					}

					if userSet.CreateTime != nil {
						userSetMap["create_time"] = userSet.CreateTime
					}

					if userSet.UserAlias != nil {
						userSetMap["user_alias"] = userSet.UserAlias
					}

					userSetList = append(userSetList, userSetMap)
				}

				userInfoMap["user_set"] = userSetList
			}

			if workGroupInfo.UserInfo.TotalCount != nil {
				userInfoMap["total_count"] = workGroupInfo.UserInfo.TotalCount
			}

			workGroupDetailInfoMap["user_info"] = []interface{}{userInfoMap}
		}

		if workGroupInfo.DataPolicyInfo != nil {
			dataPolicyInfoMap := map[string]interface{}{}

			if workGroupInfo.DataPolicyInfo.PolicySet != nil {
				var policySetList []interface{}
				for _, policySet := range workGroupInfo.DataPolicyInfo.PolicySet {
					policySetMap := map[string]interface{}{}

					if policySet.Database != nil {
						policySetMap["database"] = policySet.Database
					}

					if policySet.Catalog != nil {
						policySetMap["catalog"] = policySet.Catalog
					}

					if policySet.Table != nil {
						policySetMap["table"] = policySet.Table
					}

					if policySet.Operation != nil {
						policySetMap["operation"] = policySet.Operation
					}

					if policySet.PolicyType != nil {
						policySetMap["policy_type"] = policySet.PolicyType
					}

					if policySet.Function != nil {
						policySetMap["function"] = policySet.Function
					}

					if policySet.View != nil {
						policySetMap["view"] = policySet.View
					}

					if policySet.Column != nil {
						policySetMap["column"] = policySet.Column
					}

					if policySet.DataEngine != nil {
						policySetMap["data_engine"] = policySet.DataEngine
					}

					if policySet.ReAuth != nil {
						policySetMap["re_auth"] = policySet.ReAuth
					}

					if policySet.Source != nil {
						policySetMap["source"] = policySet.Source
					}

					if policySet.Mode != nil {
						policySetMap["mode"] = policySet.Mode
					}

					if policySet.Operator != nil {
						policySetMap["operator"] = policySet.Operator
					}

					if policySet.CreateTime != nil {
						policySetMap["create_time"] = policySet.CreateTime
					}

					if policySet.SourceId != nil {
						policySetMap["source_id"] = policySet.SourceId
					}

					if policySet.SourceName != nil {
						policySetMap["source_name"] = policySet.SourceName
					}

					if policySet.Id != nil {
						policySetMap["id"] = policySet.Id
					}

					policySetList = append(policySetList, policySetMap)
				}

				dataPolicyInfoMap["policy_set"] = policySetList
			}

			if workGroupInfo.DataPolicyInfo.TotalCount != nil {
				dataPolicyInfoMap["total_count"] = workGroupInfo.DataPolicyInfo.TotalCount
			}

			workGroupDetailInfoMap["data_policy_info"] = []interface{}{dataPolicyInfoMap}
		}

		if workGroupInfo.EnginePolicyInfo != nil {
			enginePolicyInfoMap := map[string]interface{}{}

			if workGroupInfo.EnginePolicyInfo.PolicySet != nil {
				var policySetList []interface{}
				for _, policySet := range workGroupInfo.EnginePolicyInfo.PolicySet {
					policySetMap := map[string]interface{}{}

					if policySet.Database != nil {
						policySetMap["database"] = policySet.Database
					}

					if policySet.Catalog != nil {
						policySetMap["catalog"] = policySet.Catalog
					}

					if policySet.Table != nil {
						policySetMap["table"] = policySet.Table
					}

					if policySet.Operation != nil {
						policySetMap["operation"] = policySet.Operation
					}

					if policySet.PolicyType != nil {
						policySetMap["policy_type"] = policySet.PolicyType
					}

					if policySet.Function != nil {
						policySetMap["function"] = policySet.Function
					}

					if policySet.View != nil {
						policySetMap["view"] = policySet.View
					}

					if policySet.Column != nil {
						policySetMap["column"] = policySet.Column
					}

					if policySet.DataEngine != nil {
						policySetMap["data_engine"] = policySet.DataEngine
					}

					if policySet.ReAuth != nil {
						policySetMap["re_auth"] = policySet.ReAuth
					}

					if policySet.Source != nil {
						policySetMap["source"] = policySet.Source
					}

					if policySet.Mode != nil {
						policySetMap["mode"] = policySet.Mode
					}

					if policySet.Operator != nil {
						policySetMap["operator"] = policySet.Operator
					}

					if policySet.CreateTime != nil {
						policySetMap["create_time"] = policySet.CreateTime
					}

					if policySet.SourceId != nil {
						policySetMap["source_id"] = policySet.SourceId
					}

					if policySet.SourceName != nil {
						policySetMap["source_name"] = policySet.SourceName
					}

					if policySet.Id != nil {
						policySetMap["id"] = policySet.Id
					}

					policySetList = append(policySetList, policySetMap)
				}

				enginePolicyInfoMap["policy_set"] = policySetList
			}

			if workGroupInfo.EnginePolicyInfo.TotalCount != nil {
				enginePolicyInfoMap["total_count"] = workGroupInfo.EnginePolicyInfo.TotalCount
			}

			workGroupDetailInfoMap["engine_policy_info"] = []interface{}{enginePolicyInfoMap}
		}

		if workGroupInfo.WorkGroupDescription != nil {
			workGroupDetailInfoMap["work_group_description"] = workGroupInfo.WorkGroupDescription
		}

		if workGroupInfo.RowFilterInfo != nil {
			rowFilterInfoMap := map[string]interface{}{}

			if workGroupInfo.RowFilterInfo.PolicySet != nil {
				var policySetList []interface{}
				for _, policySet := range workGroupInfo.RowFilterInfo.PolicySet {
					policySetMap := map[string]interface{}{}

					if policySet.Database != nil {
						policySetMap["database"] = policySet.Database
					}

					if policySet.Catalog != nil {
						policySetMap["catalog"] = policySet.Catalog
					}

					if policySet.Table != nil {
						policySetMap["table"] = policySet.Table
					}

					if policySet.Operation != nil {
						policySetMap["operation"] = policySet.Operation
					}

					if policySet.PolicyType != nil {
						policySetMap["policy_type"] = policySet.PolicyType
					}

					if policySet.Function != nil {
						policySetMap["function"] = policySet.Function
					}

					if policySet.View != nil {
						policySetMap["view"] = policySet.View
					}

					if policySet.Column != nil {
						policySetMap["column"] = policySet.Column
					}

					if policySet.DataEngine != nil {
						policySetMap["data_engine"] = policySet.DataEngine
					}

					if policySet.ReAuth != nil {
						policySetMap["re_auth"] = policySet.ReAuth
					}

					if policySet.Source != nil {
						policySetMap["source"] = policySet.Source
					}

					if policySet.Mode != nil {
						policySetMap["mode"] = policySet.Mode
					}

					if policySet.Operator != nil {
						policySetMap["operator"] = policySet.Operator
					}

					if policySet.CreateTime != nil {
						policySetMap["create_time"] = policySet.CreateTime
					}

					if policySet.SourceId != nil {
						policySetMap["source_id"] = policySet.SourceId
					}

					if policySet.SourceName != nil {
						policySetMap["source_name"] = policySet.SourceName
					}

					if policySet.Id != nil {
						policySetMap["id"] = policySet.Id
					}

					policySetList = append(policySetList, policySetMap)
				}

				rowFilterInfoMap["policy_set"] = policySetList
			}

			if workGroupInfo.RowFilterInfo.TotalCount != nil {
				rowFilterInfoMap["total_count"] = workGroupInfo.RowFilterInfo.TotalCount
			}

			workGroupDetailInfoMap["row_filter_info"] = []interface{}{rowFilterInfoMap}
		}

		ids = append(ids, helper.Int64ToStr(*workGroupInfo.WorkGroupId))
		_ = d.Set("work_group_info", []interface{}{workGroupDetailInfoMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), workGroupDetailInfoMap); e != nil {
			return e
		}
	}
	return nil
}
