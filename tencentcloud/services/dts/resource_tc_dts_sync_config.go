package dts

import (
	"context"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDtsSyncConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsSyncConfigCreate,
		Read:   resourceTencentCloudDtsSyncConfigRead,
		Update: resourceTencentCloudDtsSyncConfigUpdate,
		Delete: resourceTencentCloudDtsSyncConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Synchronization 实例 ID (i.e. identifies synchronization 作业)。",
			},

			"src_access_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "来源 访问 类型，cdb (云 数据库)，cvm (云 主机 self-built)，vpc (私有 网络)，extranet (外部 网络)，vpncloud (vpn 访问)，dcg (dedicated line 访问)，ccn (云 networking )，intranet (self-developed 云)，noProxy，note 该 特定 可选 值 depends 在 当前 link。",
			},

			"dst_access_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Target end 访问 类型，cdb (云 数据库)，cvm (云 主机 self-built)，vpc (私有 网络)，extranet (外部 网络)，vpncloud (vpn 访问)，dcg (dedicated line 访问)，ccn (云 networking )，intranet (self-developed 云)，noProxy，note 该 特定 可选 值 depends 在 当前 link。",
			},

			"options": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Sync 任务 Options。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"init_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Synchronous initialization options，Data (full 数据 initialization)，Structure (structure initialization)，Full (full 数据 和 structure initialization，默认值)，None (incremental 仅). 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"deal_of_exist_same_table": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "processing 的 表 使用 same 名称，ReportErrorAfterCheck (pre-check 和 报告 错误，默认值)，InitializeAfterDelete (delete 和 re-initialize)，ExecuteAfterIgnore (ignore 和 continue 到 execute). 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"conflict_handle_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Conflict handling options，ReportError (错误 报告， 默认值)，Ignore (ignore)，Cover (cover)，ConditionCover (condition coverage). 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"add_additional_column": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "是否add additional columns. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"op_types": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "DML 和 DDL options 到 是 synchronized，Insert (insert operation)，Update (update operation)，Delete (delete operation)，DDL (structure synchronization)，leave blank (不 selected)，PartialDDL (自定义，work 使用 DdlOptions). 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"conflict_handle_option": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Detailed options 对于 conflict handling，such 作为 conditional rows 和 conditional actions 在 conditional overrides. 注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"condition_column": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Columns covered 通过 condition. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"condition_operator": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Conditional Override Operation. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"condition_order_in_src_and_dst": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Conditional Override 优先级 Processing. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"ddl_options": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "DDL synchronization options，specifically describe 其中 DDLs 到 synchronize. 注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ddl_object": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Ddl 类型，such 作为 Database，Table，View，索引，etc. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"ddl_value": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "特定 值 的 ddl， possible 值 对于 Database [Create,Drop,Alter]. possible 值 对于 Table [Create,Drop,Alter,Truncate,Rename]. possible 值 对于 View[Create,Drop].For possible 值 的 索引 [Create，Drop]. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"rate_limit_option": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "任务 speed 限制 信息\n注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"current_dump_thread": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "数量 full export threads currently 在 effect. 值 的 此 字段 可以 是 adjusted 当 configuring 任务. 注意: 如果 它 是 不 集合 或 集合 到 0，它 表示 当前值 是 maintained. 最大 值 是 16.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"default_dump_thread": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "默认值 数量 full export threads. 此 字段 是 仅 meaningful 在 output 参数.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"current_dump_rps": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "full export Rps currently 在 effect. 值 的 此 字段 可以 是 adjusted 当 configuring 任务. 注意: 如果 它 是 不 集合 或 集合 到 0，它 表示 当前值 是 maintained. 最大 值 是 50,000,000.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"default_dump_rps": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "默认值 full export Rps. 此 字段 是 仅 meaningful 在 output 参数.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"current_load_thread": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "数量 full import threads currently 在 effect. 值 的 此 字段 可以 是 adjusted 当 configuring 任务. 注意: 如果 它 是 不 集合 或 集合 到 0，它 表示 当前值 是 maintained. 最大 值 是 16.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"default_load_thread": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "默认值 数量 full import threads. 此 字段 是 仅 meaningful 在 output 参数.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"current_load_rps": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "full import Rps currently 在 effect. 值 的 此 字段 可以 是 adjusted 当 configuring 任务. 注意: 如果 它 是 不 集合 或 集合 到 0，它 表示 当前值 是 maintained. 最大 值 是 50,000,000.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"default_load_rps": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "默认值 full import Rps. 此 字段 是 仅 meaningful 在 output 参数.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"current_sinker_thread": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "数量 incremental import threads currently 在 effect. 值 的 此 字段 可以 是 adjusted 当 configuring 任务. 注意: 如果 它 是 不 集合 或 集合 到 0，它 表示 当前值 是 maintained. 最大 值 是 128.注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"default_sinker_thread": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "默认值 数量 incremental import threads. 此 字段 是 仅 meaningful 在 output 参数.\n注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"has_user_set_rate_limit": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "enum:\"无\"/\"yes\", 无：用户未设置限速；是：已设置速度限制。该字段仅在输出参数中有意义。注意：该字段可能返回null，表示取不到有效值。",
									},
								},
							},
						},
					},
				},
			},

			"objects": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Synchronize 数据库 表 对象 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Migration 对象 类型 Partial (partial 对象). 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"databases": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Synchronization 对象，不 null 当 模式 是 Partial. 注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"db_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "名称 库 该 needs 到 是 migrated 或 synchronized. 此 item 为必填项 当 ObjectMode 是 Partial. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"new_db_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "名称 库 after 迁移 或 synchronization，其中 是 same 作为 来源 库 通过 默认值. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"db_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "DB selection 模式: All (对于 all objects under 当前 对象)，Partial (对于 some objects)，当 模式 是 Partial，此 item 为必填项. 注意 该 synchronization 的 advanced objects does 不 depend 在 此 值 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"schema_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Migrated 或 synchronized schema注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"new_schema_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Schema 名称 after 迁移 或 synchronization. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"table_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Table selection 模式: All (对于 all objects under 当前 对象)，Partial (对于 some objects)，此 item 为必填项 当 DBMode 是 Partial. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"tables": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "A collection 的 表 graph objects，当 TableMode 是 Partial，此 item needs 到 是 filled 在. 注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"table_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Table 名称 注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"new_table_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "New 表 名称 注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"filter_condition": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "过滤器 condition. 注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"view_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "View selection 模式: All 是 all view objects under 当前 对象，Partial 是 part 的 view objects. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"views": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "View 对象 collection，当 ViewMode 是 Partial，此 item needs 到 是 filled 在. 注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"view_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "View 名称 注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"new_view_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "New view 名称 注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"function_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Select 模式 到 是 synchronized，Partial 是 part，all 是 entire selection. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"functions": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Computed:    true,
										Description: "必填 当 FunctionMode 值 是 Partial. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"procedure_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Select 模式 到 是 synchronized，Partial 是 part，All 是 whole selection. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"procedures": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Computed:    true,
										Description: "必填 当 值 的 ProcedureMode 是 Partial. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"trigger_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Trigger 迁移 模式，all (对于 all objects under 当前 对象)，partial (partial objects). 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"triggers": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Computed:    true,
										Description: "当 TriggerMode 是 partial，指定name 的 触发器 到 是 migrated. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"event_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Event 迁移 模式，all (对于 all objects under 当前 对象)，partial (partial objects). 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"events": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Computed:    true,
										Description: "当 EventMode 是 partial，指定name 的 事件 到 是 migrated. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"advanced_objects": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Computed:    true,
							Description: "For advanced 对象 types，such 作为 函数 和 procedure，当 advanced 对象 needs 到 是 synchronized， initialization 类型 必须 include structure initialization 类型，该 是， 值 的 Options.InitType 字段 是 Structure 或 Full. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"online_ddl": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "OnlineDDL 类型 注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "状态",
									},
								},
							},
						},
					},
				},
			},

			"job_name": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Sync 作业名称",
			},

			"job_mode": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "enumeration 值 是 liteMode 和 fullMode，corresponding 到 lite 模式 或 normal 模式 respectively。",
			},

			"run_mode": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Operation 模式，such 作为: Immediate (表示immediate operation， 默认值为 此 值)，Timed (表示scheduled operation)。",
			},

			"expect_run_time": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Expected 开始时间，当 值 的 RunMode 是 Timed，此 值 为必填项，such 作为: 2006-01-02 15:04:05。",
			},

			"src_info": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "来源 信息，单个-节点 数据库 使用。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "english 名称 地域 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"role": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "节点 类型 tdsql mysql 版本， enumeration 值 是 proxy，集合. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"db_kernel": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Database kernel 类型，用于distinguish different kernels 在 tdsql: percona，mariadb，mysql. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Database 实例 ID 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IP 地址 的 实例，其中 为必填项 当 访问 类型 是 non-cdb. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "实例端口，此 item 为必填项 当 访问 类型 是 non-cdb. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"user": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "用户名，必填 对于 实例 该 require 用户名 和 密码 authentication 对于 访问. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Computed:    true,
							Description: "密码，必填 对于 实例 该 require 用户名 和 密码 authentication 对于 访问. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"db_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Database 名称，当 数据库 是 cdwpg，它 needs 到 是 提供. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Private 网络 ID，其中 为必填项 对于 访问 methods 的 私有 网络，leased line，和 VPN. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "子网 ID under 私有 网络，此 item 为必填项 对于 私有 网络，leased line，和 VPN 访问 methods. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"cvm_instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CVM 实例 short ID，其中 是 same 作为 实例 ID displayed 在 云 服务器 console 页面. 如果 它 是 self-built 实例 的 CVM，此 字段 needs 到 是 passed. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"uniq_dcg_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Leased line 网关 ID，其中 为必填项 对于 leased line 访问 类型 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"uniq_vpn_gw_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "VPN 网关 ID，其中 为必填项 对于 VPN 访问 类型 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"ccn_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cloud networking ID，其中 为必填项 对于 云 networking 访问 类型 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"supplier": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cloud vendor 类型，当 实例 是 RDS 实例，fill 在 aliyun，在 other cases fill 在 others， 默认为 others. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Database 版本，有效 仅 当 实例 是 RDS 实例，ignored 通过 other 实例， 格式 是: 5.6 或 5.7， 默认为 5.6. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"account": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "账号 到 其中 实例 belongs. 此 字段 为必填项 如果 它 是 cross-账号 实例. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"account_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "账号 到 其中 资源 belongs 是 空 或 self (表示 resources within 此 账号)，other (表示 cross-账号 resources). 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"account_role": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "角色 during cross-账号 synchronization，仅 [-zA-Z0-9-_]+ 是 allowed，如果 它 是 cross-账号 实例，此 字段 为必填项. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"role_external_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "External 角色 ID. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"tmp_secret_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Temporary 键 ID，必填 如果 它 是 cross-账号 实例. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"tmp_secret_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Temporary 键 键，必填 如果 它 是 cross-账号 实例. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"tmp_token": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Temporary 令牌，必填 如果 它 是 cross-账号 实例. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"encrypt_conn": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "是否use encrypted transmission，UnEncrypted 表示 不 到 使用 encrypted transmission，Encrypted 表示 到 使用 encrypted transmission， 默认为 UnEncrypted. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"database_net_env": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "网络 环境 到 其中 数据库 belongs. It 为必填项 当 AccessType 是 Cloud Network (CCN). `UserIDC` 表示 用户 IDC. `TencentVPC` 表示 Tencent Cloud VPC. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
					},
				},
			},

			"dst_info": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Target 信息，单个-节点 数据库 使用。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "english 名称 地域 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"role": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "节点 类型 tdsql mysql 版本， enumeration 值 是 proxy，集合. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"db_kernel": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Database kernel 类型，用于distinguish different kernels 在 tdsql: percona，mariadb，mysql. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Database 实例 ID 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IP 地址 的 实例，其中 为必填项 当 访问 类型 是 non-cdb. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "实例端口，此 item 为必填项 当 访问 类型 是 non-cdb. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"user": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "用户名，必填 对于 实例 该 require 用户名 和 密码 authentication 对于 访问. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Computed:    true,
							Description: "密码，必填 对于 实例 该 require 用户名 和 密码 authentication 对于 访问. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"db_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Database 名称，当 数据库 是 cdwpg，它 needs 到 是 提供. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Private 网络 ID，其中 为必填项 对于 访问 methods 的 私有 网络，leased line，和 VPN. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "子网 ID under 私有 网络，此 item 为必填项 对于 私有 网络，leased line，和 VPN 访问 methods. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"cvm_instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CVM 实例 short ID，其中 是 same 作为 实例 ID displayed 在 云 服务器 console 页面. 如果 它 是 self-built 实例 的 CVM，此 字段 needs 到 是 passed. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"uniq_dcg_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Leased line 网关 ID，其中 为必填项 对于 leased line 访问 类型 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"uniq_vpn_gw_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "VPN 网关 ID，其中 为必填项 对于 VPN 访问 类型 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"ccn_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cloud networking ID，其中 为必填项 对于 云 networking 访问 类型 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"supplier": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cloud vendor 类型，当 实例 是 RDS 实例，fill 在 aliyun，在 other cases fill 在 others， 默认为 others. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Database 版本，有效 仅 当 实例 是 RDS 实例，ignored 通过 other 实例， 格式 是: 5.6 或 5.7， 默认为 5.6. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"account": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "账号 到 其中 实例 belongs. 此 字段 为必填项 如果 它 是 cross-账号 实例. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"account_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "账号 到 其中 资源 belongs 是 空 或 self (表示 resources within 此 账号)，other (表示 cross-账号 resources). 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"account_role": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "角色 during cross-账号 synchronization，仅 [-zA-Z0-9-_]+ 是 allowed，如果 它 是 cross-账号 实例，此 字段 为必填项. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"role_external_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "External 角色 ID. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"tmp_secret_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Temporary 键 ID，必填 如果 它 是 cross-账号 实例. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"tmp_secret_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Temporary 键 键，必填 如果 它 是 cross-账号 实例. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"tmp_token": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Temporary 令牌，必填 如果 它 是 cross-账号 实例. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"encrypt_conn": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "是否use encrypted transmission，UnEncrypted 表示 不 到 使用 encrypted transmission，Encrypted 表示 到 使用 encrypted transmission， 默认为 UnEncrypted. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"database_net_env": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "网络 环境 到 其中 数据库 belongs. It 为必填项 当 AccessType 是 Cloud Network (CCN). `UserIDC` 表示 用户 IDC. `TencentVPC` 表示 Tencent Cloud VPC. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
					},
				},
			},

			"auto_retry_time_range_minutes": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "时间 周期 的 automatic retry，可以 是 集合 从 5 到 720 minutes，0 表示 无 retry。",
			},
		},
	}
}

func resourceTencentCloudDtsSyncConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_sync_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var jobId string
	if v, ok := d.GetOk("job_id"); ok {
		jobId = v.(string)
	}
	d.SetId(jobId)

	return resourceTencentCloudDtsSyncConfigUpdate(d, meta)
}

func resourceTencentCloudDtsSyncConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_sync_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DtsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	jobId := d.Id()

	syncConfig, err := service.DescribeDtsSyncConfigById(ctx, jobId)
	if err != nil {
		return err
	}

	if syncConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DtsSyncConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("job_id", jobId)

	if syncConfig.SrcAccessType != nil {
		_ = d.Set("src_access_type", syncConfig.SrcAccessType)
	}

	if syncConfig.DstAccessType != nil {
		_ = d.Set("dst_access_type", syncConfig.DstAccessType)
	}

	if syncConfig.Options != nil {
		optionsMap := map[string]interface{}{}

		if syncConfig.Options.InitType != nil {
			optionsMap["init_type"] = syncConfig.Options.InitType
		}

		if syncConfig.Options.DealOfExistSameTable != nil {
			optionsMap["deal_of_exist_same_table"] = syncConfig.Options.DealOfExistSameTable
		}

		if syncConfig.Options.ConflictHandleType != nil {
			optionsMap["conflict_handle_type"] = syncConfig.Options.ConflictHandleType
		}

		if syncConfig.Options.AddAdditionalColumn != nil {
			optionsMap["add_additional_column"] = syncConfig.Options.AddAdditionalColumn
		}

		if syncConfig.Options.OpTypes != nil {
			optionsMap["op_types"] = helper.StringsInterfaces(syncConfig.Options.OpTypes)
		}

		if syncConfig.Options.ConflictHandleOption != nil {
			conflictHandleOptionMap := map[string]interface{}{}

			if syncConfig.Options.ConflictHandleOption.ConditionColumn != nil {
				conflictHandleOptionMap["condition_column"] = syncConfig.Options.ConflictHandleOption.ConditionColumn
			}

			if syncConfig.Options.ConflictHandleOption.ConditionOperator != nil {
				conflictHandleOptionMap["condition_operator"] = syncConfig.Options.ConflictHandleOption.ConditionOperator
			}

			if syncConfig.Options.ConflictHandleOption.ConditionOrderInSrcAndDst != nil {
				conflictHandleOptionMap["condition_order_in_src_and_dst"] = syncConfig.Options.ConflictHandleOption.ConditionOrderInSrcAndDst
			}

			optionsMap["conflict_handle_option"] = []interface{}{conflictHandleOptionMap}
		}

		if syncConfig.Options.DdlOptions != nil {
			ddlOptionsList := []interface{}{}
			for _, ddlOptions := range syncConfig.Options.DdlOptions {
				ddlOptionsMap := map[string]interface{}{}

				if ddlOptions.DdlObject != nil {
					ddlOptionsMap["ddl_object"] = ddlOptions.DdlObject
				}

				if ddlOptions.DdlValue != nil {
					ddlOptionsMap["ddl_value"] = ddlOptions.DdlValue
				}

				ddlOptionsList = append(ddlOptionsList, ddlOptionsMap)
			}

			optionsMap["ddl_options"] = ddlOptionsList
		}

		if syncConfig.Options.RateLimitOption != nil {
			rateLimitOptionMap := map[string]interface{}{}
			if syncConfig.Options.RateLimitOption.CurrentDumpThread != nil {
				rateLimitOptionMap["current_dump_thread"] = syncConfig.Options.RateLimitOption.CurrentDumpThread
			}

			if syncConfig.Options.RateLimitOption.DefaultDumpThread != nil {
				rateLimitOptionMap["default_dump_thread"] = syncConfig.Options.RateLimitOption.DefaultDumpThread
			}

			if syncConfig.Options.RateLimitOption.CurrentDumpRps != nil {
				rateLimitOptionMap["current_dump_rps"] = syncConfig.Options.RateLimitOption.CurrentDumpRps
			}

			if syncConfig.Options.RateLimitOption.DefaultDumpRps != nil {
				rateLimitOptionMap["default_dump_rps"] = syncConfig.Options.RateLimitOption.DefaultDumpRps
			}

			if syncConfig.Options.RateLimitOption.CurrentLoadThread != nil {
				rateLimitOptionMap["current_load_thread"] = syncConfig.Options.RateLimitOption.CurrentLoadThread
			}

			if syncConfig.Options.RateLimitOption.DefaultLoadThread != nil {
				rateLimitOptionMap["default_load_thread"] = syncConfig.Options.RateLimitOption.DefaultLoadThread
			}

			if syncConfig.Options.RateLimitOption.CurrentLoadRps != nil {
				rateLimitOptionMap["current_load_rps"] = syncConfig.Options.RateLimitOption.CurrentLoadRps
			}

			if syncConfig.Options.RateLimitOption.DefaultLoadRps != nil {
				rateLimitOptionMap["default_load_rps"] = syncConfig.Options.RateLimitOption.DefaultLoadRps
			}

			if syncConfig.Options.RateLimitOption.CurrentSinkerThread != nil {
				rateLimitOptionMap["current_sinker_thread"] = syncConfig.Options.RateLimitOption.CurrentSinkerThread
			}

			if syncConfig.Options.RateLimitOption.DefaultSinkerThread != nil {
				rateLimitOptionMap["default_sinker_thread"] = syncConfig.Options.RateLimitOption.DefaultSinkerThread
			}

			if syncConfig.Options.RateLimitOption.HasUserSetRateLimit != nil {
				rateLimitOptionMap["has_user_set_rate_limit"] = syncConfig.Options.RateLimitOption.HasUserSetRateLimit
			}

			optionsMap["rate_limit_option"] = []interface{}{rateLimitOptionMap}
		}

		_ = d.Set("options", []interface{}{optionsMap})
	}

	if syncConfig.Objects != nil {
		objectsMap := map[string]interface{}{}

		if syncConfig.Objects.Mode != nil {
			objectsMap["mode"] = syncConfig.Objects.Mode
		}

		if syncConfig.Objects.Databases != nil {
			databasesList := []interface{}{}
			for _, databases := range syncConfig.Objects.Databases {
				databasesMap := map[string]interface{}{}

				if databases.DbName != nil {
					databasesMap["db_name"] = databases.DbName
				}

				if databases.NewDbName != nil {
					databasesMap["new_db_name"] = databases.NewDbName
				}

				if databases.DbMode != nil {
					databasesMap["db_mode"] = databases.DbMode
				}

				if databases.SchemaName != nil {
					databasesMap["schema_name"] = databases.SchemaName
				}

				if databases.NewSchemaName != nil {
					databasesMap["new_schema_name"] = databases.NewSchemaName
				}

				if databases.TableMode != nil {
					databasesMap["table_mode"] = databases.TableMode
				}

				if databases.Tables != nil {
					tablesList := []interface{}{}
					for _, tables := range databases.Tables {
						tablesMap := map[string]interface{}{}

						if tables.TableName != nil {
							tablesMap["table_name"] = tables.TableName
						}

						if tables.NewTableName != nil {
							tablesMap["new_table_name"] = tables.NewTableName
						}

						if tables.FilterCondition != nil {
							tablesMap["filter_condition"] = tables.FilterCondition
						}

						tablesList = append(tablesList, tablesMap)
					}

					databasesMap["tables"] = tablesList
				}

				if databases.ViewMode != nil {
					databasesMap["view_mode"] = databases.ViewMode
				}

				if databases.Views != nil {
					viewsList := []interface{}{}
					for _, views := range databases.Views {
						viewsMap := map[string]interface{}{}

						if views.ViewName != nil {
							viewsMap["view_name"] = views.ViewName
						}

						if views.NewViewName != nil {
							viewsMap["new_view_name"] = views.NewViewName
						}

						viewsList = append(viewsList, viewsMap)
					}

					databasesMap["views"] = viewsList
				}

				if databases.FunctionMode != nil {
					databasesMap["function_mode"] = databases.FunctionMode
				}

				if databases.Functions != nil {
					databasesMap["functions"] = helper.StringsInterfaces(databases.Functions)
				}

				if databases.ProcedureMode != nil {
					databasesMap["procedure_mode"] = databases.ProcedureMode
				}

				if databases.Procedures != nil {
					databasesMap["procedures"] = helper.StringsInterfaces(databases.Procedures)
				}

				if databases.TriggerMode != nil {
					databasesMap["trigger_mode"] = databases.TriggerMode
				}

				if databases.Triggers != nil {
					databasesMap["triggers"] = helper.StringsInterfaces(databases.Triggers)
				}

				if databases.EventMode != nil {
					databasesMap["event_mode"] = databases.EventMode
				}

				if databases.Events != nil {
					databasesMap["events"] = helper.StringsInterfaces(databases.Events)
				}

				databasesList = append(databasesList, databasesMap)
			}

			objectsMap["databases"] = databasesList
		}

		if syncConfig.Objects.AdvancedObjects != nil {
			objectsMap["advanced_objects"] = helper.StringsInterfaces(syncConfig.Objects.AdvancedObjects)
		}

		if syncConfig.Objects.OnlineDDL != nil {
			onlineDDLMap := map[string]interface{}{}

			if syncConfig.Objects.OnlineDDL.Status != nil {
				onlineDDLMap["status"] = syncConfig.Objects.OnlineDDL.Status
			}

			objectsMap["online_ddl"] = []interface{}{onlineDDLMap}
		}

		_ = d.Set("objects", []interface{}{objectsMap})
	}

	if syncConfig.JobName != nil {
		_ = d.Set("job_name", syncConfig.JobName)
	}

	// if syncConfig.JobMode != nil {
	// 	_ = d.Set("job_mode", syncConfig.JobMode)
	// }

	if syncConfig.RunMode != nil {
		_ = d.Set("run_mode", syncConfig.RunMode)
	}

	if syncConfig.ExpectRunTime != nil {
		_ = d.Set("expect_run_time", syncConfig.ExpectRunTime)
	}

	if syncConfig.SrcInfo != nil {
		srcInfoMap := map[string]interface{}{}

		if syncConfig.SrcInfo.Region != nil {
			srcInfoMap["region"] = syncConfig.SrcInfo.Region
		}

		if syncConfig.SrcInfo.Role != nil {
			srcInfoMap["role"] = syncConfig.SrcInfo.Role
		}

		if syncConfig.SrcInfo.DbKernel != nil {
			srcInfoMap["db_kernel"] = syncConfig.SrcInfo.DbKernel
		}

		if syncConfig.SrcInfo.InstanceId != nil {
			srcInfoMap["instance_id"] = syncConfig.SrcInfo.InstanceId
		}

		if syncConfig.SrcInfo.Ip != nil {
			srcInfoMap["ip"] = syncConfig.SrcInfo.Ip
		}

		if syncConfig.SrcInfo.Port != nil {
			srcInfoMap["port"] = syncConfig.SrcInfo.Port
		}

		if syncConfig.SrcInfo.User != nil {
			srcInfoMap["user"] = syncConfig.SrcInfo.User
		}

		if syncConfig.SrcInfo.DbName != nil {
			srcInfoMap["db_name"] = syncConfig.SrcInfo.DbName
		}

		if syncConfig.SrcInfo.VpcId != nil {
			srcInfoMap["vpc_id"] = syncConfig.SrcInfo.VpcId
		}

		if syncConfig.SrcInfo.SubnetId != nil {
			srcInfoMap["subnet_id"] = syncConfig.SrcInfo.SubnetId
		}

		if syncConfig.SrcInfo.CvmInstanceId != nil {
			srcInfoMap["cvm_instance_id"] = syncConfig.SrcInfo.CvmInstanceId
		}

		if syncConfig.SrcInfo.UniqDcgId != nil {
			srcInfoMap["uniq_dcg_id"] = syncConfig.SrcInfo.UniqDcgId
		}

		if syncConfig.SrcInfo.UniqVpnGwId != nil {
			srcInfoMap["uniq_vpn_gw_id"] = syncConfig.SrcInfo.UniqVpnGwId
		}

		if syncConfig.SrcInfo.CcnId != nil {
			srcInfoMap["ccn_id"] = syncConfig.SrcInfo.CcnId
		}

		if syncConfig.SrcInfo.Supplier != nil {
			srcInfoMap["supplier"] = syncConfig.SrcInfo.Supplier
		}

		if syncConfig.SrcInfo.EngineVersion != nil {
			srcInfoMap["engine_version"] = syncConfig.SrcInfo.EngineVersion
		}

		if syncConfig.SrcInfo.Account != nil {
			srcInfoMap["account"] = syncConfig.SrcInfo.Account
		}

		if syncConfig.SrcInfo.AccountMode != nil {
			srcInfoMap["account_mode"] = syncConfig.SrcInfo.AccountMode
		}

		if syncConfig.SrcInfo.AccountRole != nil {
			srcInfoMap["account_role"] = syncConfig.SrcInfo.AccountRole
		}

		if syncConfig.SrcInfo.RoleExternalId != nil {
			srcInfoMap["role_external_id"] = syncConfig.SrcInfo.RoleExternalId
		}

		if syncConfig.SrcInfo.TmpSecretId != nil {
			srcInfoMap["tmp_secret_id"] = syncConfig.SrcInfo.TmpSecretId
		}

		if syncConfig.SrcInfo.TmpSecretKey != nil {
			srcInfoMap["tmp_secret_key"] = syncConfig.SrcInfo.TmpSecretKey
		}

		if syncConfig.SrcInfo.TmpToken != nil {
			srcInfoMap["tmp_token"] = syncConfig.SrcInfo.TmpToken
		}

		if syncConfig.SrcInfo.EncryptConn != nil {
			srcInfoMap["encrypt_conn"] = syncConfig.SrcInfo.EncryptConn
		}

		if syncConfig.SrcInfo.DatabaseNetEnv != nil {
			srcInfoMap["database_net_env"] = syncConfig.SrcInfo.DatabaseNetEnv
		}

		// reset the password due to the describe api always return an empty string
		password := syncConfig.SrcInfo.Password
		if password != nil && *password != "" {
			srcInfoMap["password"] = password
		} else {
			key := "src_info.0.password"
			if v, ok := d.GetOk(key); ok {
				srcInfoMap["password"] = helper.String(v.(string))
			}
		}

		_ = d.Set("src_info", []interface{}{srcInfoMap})
	}

	if syncConfig.DstInfo != nil {
		dstInfoMap := map[string]interface{}{}

		if syncConfig.DstInfo.Region != nil {
			dstInfoMap["region"] = syncConfig.DstInfo.Region
		}

		if syncConfig.DstInfo.Role != nil {
			dstInfoMap["role"] = syncConfig.DstInfo.Role
		}

		if syncConfig.DstInfo.DbKernel != nil {
			dstInfoMap["db_kernel"] = syncConfig.DstInfo.DbKernel
		}

		if syncConfig.DstInfo.InstanceId != nil {
			dstInfoMap["instance_id"] = syncConfig.DstInfo.InstanceId
		}

		if syncConfig.DstInfo.Ip != nil {
			dstInfoMap["ip"] = syncConfig.DstInfo.Ip
		}

		if syncConfig.DstInfo.Port != nil {
			dstInfoMap["port"] = syncConfig.DstInfo.Port
		}

		if syncConfig.DstInfo.User != nil {
			dstInfoMap["user"] = syncConfig.DstInfo.User
		}

		if syncConfig.DstInfo.DbName != nil {
			dstInfoMap["db_name"] = syncConfig.DstInfo.DbName
		}

		if syncConfig.DstInfo.VpcId != nil {
			dstInfoMap["vpc_id"] = syncConfig.DstInfo.VpcId
		}

		if syncConfig.DstInfo.SubnetId != nil {
			dstInfoMap["subnet_id"] = syncConfig.DstInfo.SubnetId
		}

		if syncConfig.DstInfo.CvmInstanceId != nil {
			dstInfoMap["cvm_instance_id"] = syncConfig.DstInfo.CvmInstanceId
		}

		if syncConfig.DstInfo.UniqDcgId != nil {
			dstInfoMap["uniq_dcg_id"] = syncConfig.DstInfo.UniqDcgId
		}

		if syncConfig.DstInfo.UniqVpnGwId != nil {
			dstInfoMap["uniq_vpn_gw_id"] = syncConfig.DstInfo.UniqVpnGwId
		}

		if syncConfig.DstInfo.CcnId != nil {
			dstInfoMap["ccn_id"] = syncConfig.DstInfo.CcnId
		}

		if syncConfig.DstInfo.Supplier != nil {
			dstInfoMap["supplier"] = syncConfig.DstInfo.Supplier
		}

		if syncConfig.DstInfo.EngineVersion != nil {
			dstInfoMap["engine_version"] = syncConfig.DstInfo.EngineVersion
		}

		if syncConfig.DstInfo.Account != nil {
			dstInfoMap["account"] = syncConfig.DstInfo.Account
		}

		if syncConfig.DstInfo.AccountMode != nil {
			dstInfoMap["account_mode"] = syncConfig.DstInfo.AccountMode
		}

		if syncConfig.DstInfo.AccountRole != nil {
			dstInfoMap["account_role"] = syncConfig.DstInfo.AccountRole
		}

		if syncConfig.DstInfo.RoleExternalId != nil {
			dstInfoMap["role_external_id"] = syncConfig.DstInfo.RoleExternalId
		}

		if syncConfig.DstInfo.TmpSecretId != nil {
			dstInfoMap["tmp_secret_id"] = syncConfig.DstInfo.TmpSecretId
		}

		if syncConfig.DstInfo.TmpSecretKey != nil {
			dstInfoMap["tmp_secret_key"] = syncConfig.DstInfo.TmpSecretKey
		}

		if syncConfig.DstInfo.TmpToken != nil {
			dstInfoMap["tmp_token"] = syncConfig.DstInfo.TmpToken
		}

		if syncConfig.DstInfo.EncryptConn != nil {
			dstInfoMap["encrypt_conn"] = syncConfig.DstInfo.EncryptConn
		}

		// reset the password due to the describe api always return an empty string
		password := syncConfig.SrcInfo.Password
		if password != nil && *password != "" {
			dstInfoMap["password"] = password
		} else {
			key := "dst_info.0.password"
			if v, ok := d.GetOk(key); ok {
				dstInfoMap["password"] = helper.String(v.(string))
			}
		}

		_ = d.Set("dst_info", []interface{}{dstInfoMap})
	}

	if syncConfig.AutoRetryTimeRangeMinutes != nil {
		_ = d.Set("auto_retry_time_range_minutes", syncConfig.AutoRetryTimeRangeMinutes)
	}

	return nil
}

func resourceTencentCloudDtsSyncConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_sync_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := dts.NewConfigureSyncJobRequest()

	jobId := d.Id()

	request.JobId = &jobId

	if v, ok := d.GetOk("src_access_type"); ok {
		request.SrcAccessType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_access_type"); ok {
		request.DstAccessType = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "options"); ok {
		options := dts.Options{}
		if v, ok := dMap["init_type"]; ok {
			options.InitType = helper.String(v.(string))
		}
		if v, ok := dMap["deal_of_exist_same_table"]; ok {
			options.DealOfExistSameTable = helper.String(v.(string))
		}
		if v, ok := dMap["conflict_handle_type"]; ok {
			options.ConflictHandleType = helper.String(v.(string))
		}
		if v, ok := dMap["add_additional_column"]; ok {
			options.AddAdditionalColumn = helper.Bool(v.(bool))
		}
		if v, ok := dMap["op_types"]; ok {
			opTypesSet := v.(*schema.Set).List()
			for i := range opTypesSet {
				if opTypesSet[i] != nil {
					opTypes := opTypesSet[i].(string)
					options.OpTypes = append(options.OpTypes, &opTypes)
				}
			}
		}
		if conflictHandleOptionMap, ok := helper.InterfaceToMap(dMap, "conflict_handle_option"); ok {
			conflictHandleOption := dts.ConflictHandleOption{}
			if v, ok := conflictHandleOptionMap["condition_column"]; ok {
				conflictHandleOption.ConditionColumn = helper.String(v.(string))
			}
			if v, ok := conflictHandleOptionMap["condition_operator"]; ok {
				conflictHandleOption.ConditionOperator = helper.String(v.(string))
			}
			if v, ok := conflictHandleOptionMap["condition_order_in_src_and_dst"]; ok {
				conflictHandleOption.ConditionOrderInSrcAndDst = helper.String(v.(string))
			}
			options.ConflictHandleOption = &conflictHandleOption
		}
		if v, ok := dMap["ddl_options"]; ok {
			for _, item := range v.([]interface{}) {
				ddlOptionsMap := item.(map[string]interface{})
				ddlOption := dts.DdlOption{}
				if v, ok := ddlOptionsMap["ddl_object"]; ok {
					ddlOption.DdlObject = helper.String(v.(string))
				}
				if v, ok := ddlOptionsMap["ddl_value"]; ok {
					ddlValueSet := v.(*schema.Set).List()
					for i := range ddlValueSet {
						if ddlValueSet[i] != nil {
							ddlValue := ddlValueSet[i].(string)
							ddlOption.DdlValue = append(ddlOption.DdlValue, &ddlValue)
						}
					}
				}
				options.DdlOptions = append(options.DdlOptions, &ddlOption)
			}
		}
		if rateLimitOptionMap, ok := helper.InterfaceToMap(dMap, "rate_limit_option"); ok {
			rateLimitOption := dts.RateLimitOption{}
			if v, ok := rateLimitOptionMap["current_dump_thread"].(int); ok {
				rateLimitOption.CurrentDumpThread = helper.IntInt64(v)
			}

			if v, ok := rateLimitOptionMap["default_dump_thread"].(int); ok {
				rateLimitOption.DefaultDumpThread = helper.IntInt64(v)
			}

			if v, ok := rateLimitOptionMap["current_dump_rps"].(int); ok {
				rateLimitOption.CurrentDumpRps = helper.IntInt64(v)
			}

			if v, ok := rateLimitOptionMap["default_dump_rps"].(int); ok {
				rateLimitOption.DefaultDumpRps = helper.IntInt64(v)
			}

			if v, ok := rateLimitOptionMap["current_load_thread"].(int); ok {
				rateLimitOption.CurrentLoadThread = helper.IntInt64(v)
			}

			if v, ok := rateLimitOptionMap["default_load_thread"].(int); ok {
				rateLimitOption.DefaultLoadThread = helper.IntInt64(v)
			}

			if v, ok := rateLimitOptionMap["current_load_rps"].(int); ok {
				rateLimitOption.CurrentLoadRps = helper.IntInt64(v)
			}

			if v, ok := rateLimitOptionMap["default_load_rps"].(int); ok {
				rateLimitOption.DefaultLoadRps = helper.IntInt64(v)
			}

			if v, ok := rateLimitOptionMap["current_sinker_thread"].(int); ok {
				rateLimitOption.CurrentSinkerThread = helper.IntInt64(v)
			}

			if v, ok := rateLimitOptionMap["default_sinker_thread"].(int); ok {
				rateLimitOption.DefaultSinkerThread = helper.IntInt64(v)
			}

			if v, ok := rateLimitOptionMap["has_user_set_rate_limit"].(string); ok && v != "" {
				rateLimitOption.HasUserSetRateLimit = helper.String(v)
			}

			options.RateLimitOption = &rateLimitOption
		}

		request.Options = &options
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "objects"); ok {
		objects := dts.Objects{}
		if v, ok := dMap["mode"]; ok {
			objects.Mode = helper.String(v.(string))
		}
		if v, ok := dMap["databases"]; ok {
			for _, item := range v.([]interface{}) {
				databasesMap := item.(map[string]interface{})
				database := dts.Database{}
				if v, ok := databasesMap["db_name"]; ok {
					database.DbName = helper.String(v.(string))
				}
				if v, ok := databasesMap["new_db_name"]; ok {
					database.NewDbName = helper.String(v.(string))
				}
				if v, ok := databasesMap["db_mode"]; ok {
					database.DbMode = helper.String(v.(string))
				}
				if v, ok := databasesMap["schema_name"]; ok {
					database.SchemaName = helper.String(v.(string))
				}
				if v, ok := databasesMap["new_schema_name"]; ok {
					database.NewSchemaName = helper.String(v.(string))
				}
				if v, ok := databasesMap["table_mode"]; ok {
					database.TableMode = helper.String(v.(string))
				}
				if v, ok := databasesMap["tables"]; ok {
					for _, item := range v.([]interface{}) {
						tablesMap := item.(map[string]interface{})
						table := dts.Table{}
						if v, ok := tablesMap["table_name"]; ok {
							table.TableName = helper.String(v.(string))
						}
						if v, ok := tablesMap["new_table_name"]; ok {
							table.NewTableName = helper.String(v.(string))
						}
						if v, ok := tablesMap["filter_condition"]; ok {
							table.FilterCondition = helper.String(v.(string))
						}
						database.Tables = append(database.Tables, &table)
					}
				}
				if v, ok := databasesMap["view_mode"]; ok {
					database.ViewMode = helper.String(v.(string))
				}
				if v, ok := databasesMap["views"]; ok {
					for _, item := range v.([]interface{}) {
						viewsMap := item.(map[string]interface{})
						view := dts.View{}
						if v, ok := viewsMap["view_name"]; ok {
							view.ViewName = helper.String(v.(string))
						}
						if v, ok := viewsMap["new_view_name"]; ok {
							view.NewViewName = helper.String(v.(string))
						}
						database.Views = append(database.Views, &view)
					}
				}
				if v, ok := databasesMap["function_mode"]; ok {
					database.FunctionMode = helper.String(v.(string))
				}
				if v, ok := databasesMap["functions"]; ok {
					functionsSet := v.(*schema.Set).List()
					for i := range functionsSet {
						if functionsSet[i] != nil {
							functions := functionsSet[i].(string)
							database.Functions = append(database.Functions, &functions)
						}
					}
				}
				if v, ok := databasesMap["procedure_mode"]; ok {
					database.ProcedureMode = helper.String(v.(string))
				}
				if v, ok := databasesMap["procedures"]; ok {
					proceduresSet := v.(*schema.Set).List()
					for i := range proceduresSet {
						if proceduresSet[i] != nil {
							procedures := proceduresSet[i].(string)
							database.Procedures = append(database.Procedures, &procedures)
						}
					}
				}
				if v, ok := databasesMap["trigger_mode"]; ok {
					database.TriggerMode = helper.String(v.(string))
				}
				if v, ok := databasesMap["triggers"]; ok {
					triggersSet := v.(*schema.Set).List()
					for i := range triggersSet {
						if triggersSet[i] != nil {
							triggers := triggersSet[i].(string)
							database.Triggers = append(database.Triggers, &triggers)
						}
					}
				}
				if v, ok := databasesMap["event_mode"]; ok {
					database.EventMode = helper.String(v.(string))
				}
				if v, ok := databasesMap["events"]; ok {
					eventsSet := v.(*schema.Set).List()
					for i := range eventsSet {
						if eventsSet[i] != nil {
							events := eventsSet[i].(string)
							database.Events = append(database.Events, &events)
						}
					}
				}
				objects.Databases = append(objects.Databases, &database)
			}
		}
		if v, ok := dMap["advanced_objects"]; ok {
			advancedObjectsSet := v.(*schema.Set).List()
			for i := range advancedObjectsSet {
				if advancedObjectsSet[i] != nil {
					advancedObjects := advancedObjectsSet[i].(string)
					objects.AdvancedObjects = append(objects.AdvancedObjects, &advancedObjects)
				}
			}
		}
		if onlineDDLMap, ok := helper.InterfaceToMap(dMap, "online_ddl"); ok {
			onlineDDL := dts.OnlineDDL{}
			if v, ok := onlineDDLMap["status"]; ok {
				onlineDDL.Status = helper.String(v.(string))
			}
			objects.OnlineDDL = &onlineDDL
		}
		request.Objects = &objects
	}

	if v, ok := d.GetOk("job_name"); ok {
		request.JobName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("job_mode"); ok {
		request.JobMode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("run_mode"); ok {
		request.RunMode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("expect_run_time"); ok {
		request.ExpectRunTime = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "src_info"); ok {
		endpoint := dts.Endpoint{}
		if v, ok := dMap["region"]; ok {
			endpoint.Region = helper.String(v.(string))
		}
		if v, ok := dMap["role"]; ok {
			endpoint.Role = helper.String(v.(string))
		}
		if v, ok := dMap["db_kernel"]; ok {
			endpoint.DbKernel = helper.String(v.(string))
		}
		if v, ok := dMap["instance_id"]; ok {
			endpoint.InstanceId = helper.String(v.(string))
		}
		if v, ok := dMap["ip"]; ok {
			endpoint.Ip = helper.String(v.(string))
		}
		if v, ok := dMap["port"]; ok {
			endpoint.Port = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["user"]; ok {
			endpoint.User = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			endpoint.Password = helper.String(v.(string))
		}
		if v, ok := dMap["db_name"]; ok {
			endpoint.DbName = helper.String(v.(string))
		}
		if v, ok := dMap["vpc_id"]; ok {
			endpoint.VpcId = helper.String(v.(string))
		}
		if v, ok := dMap["subnet_id"]; ok {
			endpoint.SubnetId = helper.String(v.(string))
		}
		if v, ok := dMap["cvm_instance_id"]; ok {
			endpoint.CvmInstanceId = helper.String(v.(string))
		}
		if v, ok := dMap["uniq_dcg_id"]; ok {
			endpoint.UniqDcgId = helper.String(v.(string))
		}
		if v, ok := dMap["uniq_vpn_gw_id"]; ok {
			endpoint.UniqVpnGwId = helper.String(v.(string))
		}
		if v, ok := dMap["ccn_id"]; ok {
			endpoint.CcnId = helper.String(v.(string))
		}
		if v, ok := dMap["supplier"]; ok {
			endpoint.Supplier = helper.String(v.(string))
		}
		if v, ok := dMap["engine_version"]; ok {
			endpoint.EngineVersion = helper.String(v.(string))
		}
		if v, ok := dMap["account"]; ok {
			endpoint.Account = helper.String(v.(string))
		}
		if v, ok := dMap["account_mode"]; ok {
			endpoint.AccountMode = helper.String(v.(string))
		}
		if v, ok := dMap["account_role"]; ok {
			endpoint.AccountRole = helper.String(v.(string))
		}
		if v, ok := dMap["role_external_id"]; ok {
			endpoint.RoleExternalId = helper.String(v.(string))
		}
		if v, ok := dMap["tmp_secret_id"]; ok {
			endpoint.TmpSecretId = helper.String(v.(string))
		}
		if v, ok := dMap["tmp_secret_key"]; ok {
			endpoint.TmpSecretKey = helper.String(v.(string))
		}
		if v, ok := dMap["tmp_token"]; ok {
			endpoint.TmpToken = helper.String(v.(string))
		}
		if v, ok := dMap["encrypt_conn"]; ok {
			endpoint.EncryptConn = helper.String(v.(string))
		}
		if v, ok := dMap["database_net_env"]; ok {
			endpoint.DatabaseNetEnv = helper.String(v.(string))
		}
		request.SrcInfo = &endpoint
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "dst_info"); ok {
		endpoint := dts.Endpoint{}
		if v, ok := dMap["region"]; ok {
			endpoint.Region = helper.String(v.(string))
		}
		if v, ok := dMap["role"]; ok {
			endpoint.Role = helper.String(v.(string))
		}
		if v, ok := dMap["db_kernel"]; ok {
			endpoint.DbKernel = helper.String(v.(string))
		}
		if v, ok := dMap["instance_id"]; ok {
			endpoint.InstanceId = helper.String(v.(string))
		}
		if v, ok := dMap["ip"]; ok {
			endpoint.Ip = helper.String(v.(string))
		}
		if v, ok := dMap["port"]; ok {
			endpoint.Port = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["user"]; ok {
			endpoint.User = helper.String(v.(string))
		}
		if v, ok := dMap["password"]; ok {
			endpoint.Password = helper.String(v.(string))
		}
		if v, ok := dMap["db_name"]; ok {
			endpoint.DbName = helper.String(v.(string))
		}
		if v, ok := dMap["vpc_id"]; ok {
			endpoint.VpcId = helper.String(v.(string))
		}
		if v, ok := dMap["subnet_id"]; ok {
			endpoint.SubnetId = helper.String(v.(string))
		}
		if v, ok := dMap["cvm_instance_id"]; ok {
			endpoint.CvmInstanceId = helper.String(v.(string))
		}
		if v, ok := dMap["uniq_dcg_id"]; ok {
			endpoint.UniqDcgId = helper.String(v.(string))
		}
		if v, ok := dMap["uniq_vpn_gw_id"]; ok {
			endpoint.UniqVpnGwId = helper.String(v.(string))
		}
		if v, ok := dMap["ccn_id"]; ok {
			endpoint.CcnId = helper.String(v.(string))
		}
		if v, ok := dMap["supplier"]; ok {
			endpoint.Supplier = helper.String(v.(string))
		}
		if v, ok := dMap["engine_version"]; ok {
			endpoint.EngineVersion = helper.String(v.(string))
		}
		if v, ok := dMap["account"]; ok {
			endpoint.Account = helper.String(v.(string))
		}
		if v, ok := dMap["account_mode"]; ok {
			endpoint.AccountMode = helper.String(v.(string))
		}
		if v, ok := dMap["account_role"]; ok {
			endpoint.AccountRole = helper.String(v.(string))
		}
		if v, ok := dMap["role_external_id"]; ok {
			endpoint.RoleExternalId = helper.String(v.(string))
		}
		if v, ok := dMap["tmp_secret_id"]; ok {
			endpoint.TmpSecretId = helper.String(v.(string))
		}
		if v, ok := dMap["tmp_secret_key"]; ok {
			endpoint.TmpSecretKey = helper.String(v.(string))
		}
		if v, ok := dMap["tmp_token"]; ok {
			endpoint.TmpToken = helper.String(v.(string))
		}
		if v, ok := dMap["encrypt_conn"]; ok {
			endpoint.EncryptConn = helper.String(v.(string))
		}
		if v, ok := dMap["database_net_env"]; ok {
			endpoint.DatabaseNetEnv = helper.String(v.(string))
		}
		request.DstInfo = &endpoint
	}

	if v, ok := d.GetOkExists("auto_retry_time_range_minutes"); ok {
		request.AutoRetryTimeRangeMinutes = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDtsClient().ConfigureSyncJob(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dts syncConfig failed, reason:%+v", logId, err)
		return err
	}

	service := DtsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"Initialized"}, tccommon.ReadRetryTimeout, time.Second, service.DtsSyncJobStateRefreshFunc(d.Id(), "", []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudDtsSyncConfigRead(d, meta)
}

func resourceTencentCloudDtsSyncConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_sync_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
