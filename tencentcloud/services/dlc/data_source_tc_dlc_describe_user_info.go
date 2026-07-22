package dlc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDlcDescribeUserInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcDescribeUserInfoRead,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "用户 ID。",
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "类型 queried information. Group: working group; DataAuth: data permission; EngineAuth: engine permission。",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter criteria that are queriedWhen the 类型 is Group，the fuzzy search is supported as the 键 is workgroup-名称When the 类型 is DataAuth，the keys supported are:policy-类型: types of permissions;policy-来源: data sources;data-名称: fuzzy search of the database and table.When the 类型 is EngineAuth，the keys supported are:policy-类型: types of permissions;policy-来源: data sources;engine-名称: fuzzy search of the database and table。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Attribute 名称 If more than one filter exists，the logical relationship between these filters is `OR`。",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Attribute 值 If multiple values exist in one filter，the logical relationship between these values is `OR`。",
						},
					},
				},
			},

			"sort_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort fields.When the 类型 is Group，the create-time and group-名称 are supported.When the 类型 is DataAuth，create-time is supported.When the 类型 is EngineAuth，create-time is supported。",
			},

			"sorting": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting methods: desc means in 顺序; asc means in reverse 顺序; it is asc by default。",
			},

			"user_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Detailed 用户 information注意：此字段可能返回 null，表示无法获取有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "用户 ID注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Types of returned information. Group: returned information about the working group where the current 用户 is; DataAuth: returned information about the current 用户&amp;#39;s data permission; EngineAuth: returned information about the current 用户&amp;#39;s engine permission注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"user_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Types of users. ADMIN: administrators; COMMON: general users注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"user_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "用户 description注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"data_policy_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of data permission information注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Collection of policies注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target database. `*` represents all databases in the current catalog. To grant admin permissions，it must be `*`; to grant data connection permissions，it must be null; to grant other permissions，it can be any database。",
												},
												"catalog": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target data 来源 To grant admin permission，it must be `*` (all resources at this 级别); to grant data 来源 and database permissions，it must be `COSDataCatalog` or `*`; to grant table permissions，it can be a custom data 来源; if it is left empty，`DataLakeCatalog` is used. Note: To grant permissions on a custom data 来源，the permissions that can be managed in the Data Lake Compute console are subsets of the 账号 permissions granted when you connect the data 来源 to the console。",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target table. `*` represents all tables in the current database. To grant admin permissions，it must be `*`; to grant data connection and database permissions，it must be null; to grant other permissions，it can be any table。",
												},
												"operation": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The target permissions，which vary by permission 级别 Admin: `ALL` (default); data connection: `CREATE`; database: `ALL`，`CREATE`，`ALTER`，and `DROP`; table: `ALL`，`SELECT`，`INSERT`，`ALTER`，`DELETE`，`DROP`，and `UPDATE`. Note: For table permissions，if a data 来源 other than `COSDataCatalog` is specified，only the `SELECT` permission can be granted here。",
												},
												"policy_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission 类型 有效值：`ADMIN`，`DATASOURCE`，`DATABASE`，`TABLE`，`VIEW`，`FUNCTION`，`COLUMN`，and `ENGINE`. Note: If it is left empty，`ADMIN` is used。",
												},
												"function": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target function. `*` represents all functions in the current catalog. To grant admin permissions，it must be `*`; to grant data connection permissions，it must be null; to grant other permissions，it can be any function.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"view": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target view. `*` represents all views in the current database. To grant admin permissions，it must be `*`; to grant data connection and database permissions，it must be null; to grant other permissions，it can be any view.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"column": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target column. `*` represents all columns. To grant admin permissions，it must be `*`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"data_engine": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target data engine. `*` represents all engines. To grant admin permissions，it must be `*`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"re_auth": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否grantee is allowed to further grant the permissions. 有效值：`false` (default) and `true` (the grantee can grant permissions gained here to other sub-users).注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"source": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission 来源，which 不是必填项 when input parameters are passed in. 有效值：`USER` (from the 用户) and `WORKGROUP` (from one or more associated work groups).注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The grant 模式，which 不是必填项 as an input parameter. 有效值：`COMMON` and `SENIOR`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 操作者，which 不是必填项 as an input parameter.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission policy 创建时间，which 不是必填项 as an input parameter.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"source_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The ID work group，which applies only when the 值 of the `来源` field is `WORKGROUP`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"source_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 work group，which applies only when the 值 of the `来源` field is `WORKGROUP`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The policy ID.注意：此字段可能返回 null，表示无法获取有效值。",
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
							Description: "Collection of engine permissions注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Collection of policies注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target database. `*` represents all databases in the current catalog. To grant admin permissions，it must be `*`; to grant data connection permissions，it must be null; to grant other permissions，it can be any database。",
												},
												"catalog": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target data 来源 To grant admin permission，it must be `*` (all resources at this 级别); to grant data 来源 and database permissions，it must be `COSDataCatalog` or `*`; to grant table permissions，it can be a custom data 来源; if it is left empty，`DataLakeCatalog` is used. Note: To grant permissions on a custom data 来源，the permissions that can be managed in the Data Lake Compute console are subsets of the 账号 permissions granted when you connect the data 来源 to the console。",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target table. `*` represents all tables in the current database. To grant admin permissions，it must be `*`; to grant data connection and database permissions，it must be null; to grant other permissions，it can be any table。",
												},
												"operation": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The target permissions，which vary by permission 级别 Admin: `ALL` (default); data connection: `CREATE`; database: `ALL`，`CREATE`，`ALTER`，and `DROP`; table: `ALL`，`SELECT`，`INSERT`，`ALTER`，`DELETE`，`DROP`，and `UPDATE`. Note: For table permissions，if a data 来源 other than `COSDataCatalog` is specified，only the `SELECT` permission can be granted here。",
												},
												"policy_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission 类型 有效值：`ADMIN`，`DATASOURCE`，`DATABASE`，`TABLE`，`VIEW`，`FUNCTION`，`COLUMN`，and `ENGINE`. Note: If it is left empty，`ADMIN` is used。",
												},
												"function": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target function. `*` represents all functions in the current catalog. To grant admin permissions，it must be `*`; to grant data connection permissions，it must be null; to grant other permissions，it can be any function.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"view": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target view. `*` represents all views in the current database. To grant admin permissions，it must be `*`; to grant data connection and database permissions，it must be null; to grant other permissions，it can be any view.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"column": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target column. `*` represents all columns. To grant admin permissions，it must be `*`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"data_engine": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target data engine. `*` represents all engines. To grant admin permissions，it must be `*`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"re_auth": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否grantee is allowed to further grant the permissions. 有效值：`false` (default) and `true` (the grantee can grant permissions gained here to other sub-users).注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"source": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission 来源，which 不是必填项 when input parameters are passed in. 有效值：`USER` (from the 用户) and `WORKGROUP` (from one or more associated work groups).注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The grant 模式，which 不是必填项 as an input parameter. 有效值：`COMMON` and `SENIOR`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 操作者，which 不是必填项 as an input parameter.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission policy 创建时间，which 不是必填项 as an input parameter.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"source_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The ID work group，which applies only when the 值 of the `来源` field is `WORKGROUP`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"source_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 work group，which applies only when the 值 of the `来源` field is `WORKGROUP`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The policy ID.注意：此字段可能返回 null，表示无法获取有效值。",
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
						"work_group_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Information about collections of working groups bound to the user注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"work_group_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Collection of working group information注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"work_group_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Unique ID working group。",
												},
												"work_group_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Working 组名称",
												},
												"work_group_description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Working group description注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"creator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "创建者",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 创建时间 of the working group，e.g. at 16:19:32 on Jul 28，2021。",
												},
											},
										},
									},
									"total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total working groups注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"user_alias": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "用户 alias注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"row_filter_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of filtered rows注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Collection of policies注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target database. `*` represents all databases in the current catalog. To grant admin permissions，it must be `*`; to grant data connection permissions，it must be null; to grant other permissions，it can be any database。",
												},
												"catalog": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target data 来源 To grant admin permission，it must be `*` (all resources at this 级别); to grant data 来源 and database permissions，it must be `COSDataCatalog` or `*`; to grant table permissions，it can be a custom data 来源; if it is left empty，`DataLakeCatalog` is used. Note: To grant permissions on a custom data 来源，the permissions that can be managed in the Data Lake Compute console are subsets of the 账号 permissions granted when you connect the data 来源 to the console。",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target table. `*` represents all tables in the current database. To grant admin permissions，it must be `*`; to grant data connection and database permissions，it must be null; to grant other permissions，it can be any table。",
												},
												"operation": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The target permissions，which vary by permission 级别 Admin: `ALL` (default); data connection: `CREATE`; database: `ALL`，`CREATE`，`ALTER`，and `DROP`; table: `ALL`，`SELECT`，`INSERT`，`ALTER`，`DELETE`，`DROP`，and `UPDATE`. Note: For table permissions，if a data 来源 other than `COSDataCatalog` is specified，only the `SELECT` permission can be granted here。",
												},
												"policy_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission 类型 有效值：`ADMIN`，`DATASOURCE`，`DATABASE`，`TABLE`，`VIEW`，`FUNCTION`，`COLUMN`，and `ENGINE`. Note: If it is left empty，`ADMIN` is used。",
												},
												"function": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target function. `*` represents all functions in the current catalog. To grant admin permissions，it must be `*`; to grant data connection permissions，it must be null; to grant other permissions，it can be any function.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"view": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target view. `*` represents all views in the current database. To grant admin permissions，it must be `*`; to grant data connection and database permissions，it must be null; to grant other permissions，it can be any view.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"column": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target column. `*` represents all columns. To grant admin permissions，it must be `*`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"data_engine": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target data engine. `*` represents all engines. To grant admin permissions，it must be `*`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"re_auth": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否grantee is allowed to further grant the permissions. 有效值：`false` (default) and `true` (the grantee can grant permissions gained here to other sub-users).注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"source": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission 来源，which 不是必填项 when input parameters are passed in. 有效值：`USER` (from the 用户) and `WORKGROUP` (from one or more associated work groups).注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The grant 模式，which 不是必填项 as an input parameter. 有效值：`COMMON` and `SENIOR`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 操作者，which 不是必填项 as an input parameter.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission policy 创建时间，which 不是必填项 as an input parameter.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"source_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The ID work group，which applies only when the 值 of the `来源` field is `WORKGROUP`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"source_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 work group，which applies only when the 值 of the `来源` field is `WORKGROUP`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The policy ID.注意：此字段可能返回 null，表示无法获取有效值。",
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
						"account_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "账号 类型",
						},
						"catalog_policy_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Collection of catalog permissions注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Collection of policies注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target database. `*` represents all databases in the current catalog. To grant admin permissions，it must be `*`; to grant data connection permissions，it must be null; to grant other permissions，it can be any database。",
												},
												"catalog": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target data 来源 To grant admin permission，it must be `*` (all resources at this 级别); to grant data 来源 and database permissions，it must be `COSDataCatalog` or `*`; to grant table permissions，it can be a custom data 来源; if it is left empty，`DataLakeCatalog` is used. Note: To grant permissions on a custom data 来源，the permissions that can be managed in the Data Lake Compute console are subsets of the 账号 permissions granted when you connect the data 来源 to the console。",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target table. `*` represents all tables in the current database. To grant admin permissions，it must be `*`; to grant data connection and database permissions，it must be null; to grant other permissions，it can be any table。",
												},
												"operation": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The target permissions，which vary by permission 级别 Admin: `ALL` (default); data connection: `CREATE`; database: `ALL`，`CREATE`，`ALTER`，and `DROP`; table: `ALL`，`SELECT`，`INSERT`，`ALTER`，`DELETE`，`DROP`，and `UPDATE`. Note: For table permissions，if a data 来源 other than `COSDataCatalog` is specified，only the `SELECT` permission can be granted here。",
												},
												"policy_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission 类型 有效值：`ADMIN`，`DATASOURCE`，`DATABASE`，`TABLE`，`VIEW`，`FUNCTION`，`COLUMN`，and `ENGINE`. Note: If it is left empty，`ADMIN` is used。",
												},
												"function": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target function. `*` represents all functions in the current catalog. To grant admin permissions，it must be `*`; to grant data connection permissions，it must be null; to grant other permissions，it can be any function.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"view": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target view. `*` represents all views in the current database. To grant admin permissions，it must be `*`; to grant data connection and database permissions，it must be null; to grant other permissions，it can be any view.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"column": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target column. `*` represents all columns. To grant admin permissions，it must be `*`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"data_engine": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 target data engine. `*` represents all engines. To grant admin permissions，it must be `*`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"re_auth": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否grantee is allowed to further grant the permissions. 有效值：`false` (default) and `true` (the grantee can grant permissions gained here to other sub-users).注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"source": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission 来源，which 不是必填项 when input parameters are passed in. 有效值：`USER` (from the 用户) and `WORKGROUP` (from one or more associated work groups).注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The grant 模式，which 不是必填项 as an input parameter. 有效值：`COMMON` and `SENIOR`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 操作者，which 不是必填项 as an input parameter.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"create_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission policy 创建时间，which 不是必填项 as an input parameter.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"source_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The ID work group，which applies only when the 值 of the `来源` field is `WORKGROUP`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"source_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 work group，which applies only when the 值 of the `来源` field is `WORKGROUP`.注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The policy ID.注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudDlcDescribeUserInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_describe_user_info.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var userId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("user_id"); ok {
		userId = v.(string)
		paramMap["UserId"] = helper.String(v.(string))
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

	var userInfo *dlc.UserDetailInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcDescribeUserInfoByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		userInfo = result
		return nil
	})
	if err != nil {
		return err
	}

	userDetailInfoMap := map[string]interface{}{}

	if userInfo != nil {

		if userInfo.UserId != nil {
			userDetailInfoMap["user_id"] = userInfo.UserId
		}

		if userInfo.Type != nil {
			userDetailInfoMap["type"] = userInfo.Type
		}

		if userInfo.UserType != nil {
			userDetailInfoMap["user_type"] = userInfo.UserType
		}

		if userInfo.UserDescription != nil {
			userDetailInfoMap["user_description"] = userInfo.UserDescription
		}

		if userInfo.DataPolicyInfo != nil {
			dataPolicyInfoMap := map[string]interface{}{}

			if userInfo.DataPolicyInfo.PolicySet != nil {
				var policySetList []interface{}
				for _, policySet := range userInfo.DataPolicyInfo.PolicySet {
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

			if userInfo.DataPolicyInfo.TotalCount != nil {
				dataPolicyInfoMap["total_count"] = userInfo.DataPolicyInfo.TotalCount
			}

			userDetailInfoMap["data_policy_info"] = []interface{}{dataPolicyInfoMap}
		}

		if userInfo.EnginePolicyInfo != nil {
			enginePolicyInfoMap := map[string]interface{}{}

			if userInfo.EnginePolicyInfo.PolicySet != nil {
				var policySetList []interface{}
				for _, policySet := range userInfo.EnginePolicyInfo.PolicySet {
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

			if userInfo.EnginePolicyInfo.TotalCount != nil {
				enginePolicyInfoMap["total_count"] = userInfo.EnginePolicyInfo.TotalCount
			}

			userDetailInfoMap["engine_policy_info"] = []interface{}{enginePolicyInfoMap}
		}

		if userInfo.WorkGroupInfo != nil {
			workGroupInfoMap := map[string]interface{}{}

			if userInfo.WorkGroupInfo.WorkGroupSet != nil {
				var workGroupSetList []interface{}
				for _, workGroupSet := range userInfo.WorkGroupInfo.WorkGroupSet {
					workGroupSetMap := map[string]interface{}{}

					if workGroupSet.WorkGroupId != nil {
						workGroupSetMap["work_group_id"] = workGroupSet.WorkGroupId
					}

					if workGroupSet.WorkGroupName != nil {
						workGroupSetMap["work_group_name"] = workGroupSet.WorkGroupName
					}

					if workGroupSet.WorkGroupDescription != nil {
						workGroupSetMap["work_group_description"] = workGroupSet.WorkGroupDescription
					}

					if workGroupSet.Creator != nil {
						workGroupSetMap["creator"] = workGroupSet.Creator
					}

					if workGroupSet.CreateTime != nil {
						workGroupSetMap["create_time"] = workGroupSet.CreateTime
					}

					workGroupSetList = append(workGroupSetList, workGroupSetMap)
				}

				workGroupInfoMap["work_group_set"] = workGroupSetList
			}

			if userInfo.WorkGroupInfo.TotalCount != nil {
				workGroupInfoMap["total_count"] = userInfo.WorkGroupInfo.TotalCount
			}

			userDetailInfoMap["work_group_info"] = []interface{}{workGroupInfoMap}
		}

		if userInfo.UserAlias != nil {
			userDetailInfoMap["user_alias"] = userInfo.UserAlias
		}

		if userInfo.RowFilterInfo != nil {
			rowFilterInfoMap := map[string]interface{}{}

			if userInfo.RowFilterInfo.PolicySet != nil {
				var policySetList []interface{}
				for _, policySet := range userInfo.RowFilterInfo.PolicySet {
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

			if userInfo.RowFilterInfo.TotalCount != nil {
				rowFilterInfoMap["total_count"] = userInfo.RowFilterInfo.TotalCount
			}

			userDetailInfoMap["row_filter_info"] = []interface{}{rowFilterInfoMap}
		}

		if userInfo.AccountType != nil {
			userDetailInfoMap["account_type"] = userInfo.AccountType
		}

		if userInfo.CatalogPolicyInfo != nil {
			catalogPolicyInfoMap := map[string]interface{}{}

			if userInfo.CatalogPolicyInfo.PolicySet != nil {
				policySetList := []interface{}{}
				for _, policySet := range userInfo.CatalogPolicyInfo.PolicySet {
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

				catalogPolicyInfoMap["policy_set"] = policySetList
			}

			if userInfo.CatalogPolicyInfo.TotalCount != nil {
				catalogPolicyInfoMap["total_count"] = userInfo.CatalogPolicyInfo.TotalCount
			}

			userDetailInfoMap["catalog_policy_info"] = []interface{}{catalogPolicyInfoMap}
		}

		_ = d.Set("user_info", []interface{}{userDetailInfoMap})
	}

	d.SetId(userId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), userDetailInfoMap); e != nil {
			return e
		}
	}
	return nil
}
