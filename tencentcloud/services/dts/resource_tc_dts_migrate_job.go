package dts

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDtsMigrateJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDtsMigrateJobCreate,
		Read:   resourceTencentCloudDtsMigrateJobRead,
		Update: resourceTencentCloudDtsMigrateJobUpdate,
		Delete: resourceTencentCloudDtsMigrateJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"service_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Migrate 服务 ID 从 `tencentcloud_dts_migrate_service`。",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "任务 状态 有效值：创建(Created)，checking (Checking)，checkPass (Check passed)，checkNotPass (Check 不 passed)，readyRun (Ready 对于 running)，running (Running)，readyComplete (Preparation completed)，success (Successful)，failed (Failed)，stopping (Stopping)，completing (Completing)，pausing (Pausing)，manualPaused (Paused)。",
			},

			// for modify operation
			"run_mode": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Running 模式 有效值：immediate，timed。",
			},

			"migrate_option": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Migration 作业 配置 options，用于describe how 任务 performs 迁移。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_table": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Migration 对象 选项，您 need 到 tell 迁移 服务 其中 库 表 objects 到 migrate。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"object_mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Migration 对象 类型 有效值：all，partial。",
									},
									"databases": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Migration 对象，其中 为必填项 如果 ObjectMode 是 partial。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"db_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "名称 数据库 到 是 migrated 或 synced，其中 为必填项 如果 ObjectMode 是 partial。",
												},
												"new_db_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "名称 数据库 after 迁移 或 sync，其中 是 same 作为 来源 数据库 名称 通过 默认值。",
												},
												"schema_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "schema 到 是 migrated 或 synced。",
												},
												"new_schema_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "名称 schema after 迁移 或 sync。",
												},
												"db_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Database selection 模式，其中 为必填项 如果 ObjectMode 是 partial. 有效值：all，partial。",
												},
												"schema_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Schema selection 模式 有效值：all，partial。",
												},
												"table_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Table selection 模式，其中 为必填项 如果 DBMode 是 partial. 有效值：all，partial。",
												},
												"tables": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "集合 的 表 objects，其中 为必填项 如果 TableMode 是 partial。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"table_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "名称 migrated 表，其中 是 case-sensitive。",
															},
															"new_table_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "New 名称 migrated 表. 此 参数 为必填项 当 TableEditMode 是 rename. It 是 mutually exclusive 使用 TmpTables.。",
															},
															"tmp_tables": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Optional:    true,
																Computed:    true,
																Description: "temp tables 到 是 migrated. 此 参数 是 mutually exclusive 使用 NewTableName. It 是 有效 仅 当 已配置 迁移 objects 是 表-级别 ones 和 TableEditMode 是 pt. To migrate temp tables generated 当 pt-osc 或 other tools 是 使用 during 迁移 process，您 必须 configure 此 参数 first. For 示例，如果 您 want 到 perform pt-osc operation 在 表 named 't1'，configure 此 参数 作为 ['_t1_new','_t1_old']; 到 perform gh-ost operation 在 t1，configure 它 作为 ['_t1_ghc','_t1_gho','_t1_del']. Temp tables generated 通过 pt-osc 和 gh-ost operations 可以 是 已配置 在 same 时间。",
															},
															"table_edit_mode": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Table editing 类型 有效值：rename (表 mapping); pt (additional 表 sync)。",
															},
														},
													},
												},
												"view_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "View selection 模式 有效值：all，partial。",
												},
												"views": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "集合 的 view objects，其中 为必填项 如果 ViewMode 是 partial。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"view_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "View 名称",
															},
															"new_view_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "View 名称 after 迁移。",
															},
														},
													},
												},
												"role_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "角色 selection 模式，其中 是 exclusive 到 PostgreSQL. 有效值：all，partial。",
												},
												"roles": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "角色，其中 是 exclusive 到 PostgreSQL 和 必填 如果 RoleMode 是 partial。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"role_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "角色 名称",
															},
															"new_role_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "角色 名称 after 迁移。",
															},
														},
													},
												},
												"function_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Sync 模式 有效值：partial，all。",
												},
												"trigger_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Sync 模式 有效值：partial，all。",
												},
												"event_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Sync 模式 有效值：partial，all。",
												},
												"procedure_mode": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Sync 模式 有效值：partial，all。",
												},
												"functions": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Computed:    true,
													Description: "此 参数 为必填项 如果 FunctionMode 是 partial。",
												},
												"procedures": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Computed:    true,
													Description: "此 参数 为必填项 如果 ProcedureMode 是 partial。",
												},
												"events": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Computed:    true,
													Description: "此 参数 为必填项 如果 EventMode 是 partial。",
												},
												"triggers": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Optional:    true,
													Computed:    true,
													Description: "此 参数 为必填项 如果 TriggerMode 是 partial。",
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
										Description: "Advanced 对象 types，such 作为 触发器，函数，procedure，事件. 注意: 如果 您 want 到 migrate 和 synchronize advanced objects， corresponding advanced 对象 类型 should 是 included 在 此 配置。",
									},
								},
							},
						},
						"migrate_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Migration 类型 有效值：full，structure，fullAndIncrement. 默认值：fullAndIncrement。",
						},
						"consistency": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Data consistency check 选项. Data consistency check 是 已禁用 通过 默认值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Data consistency check 类型 有效值：full，noCheck，notConfigured。",
									},
								},
							},
						},
						"is_migrate_account": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "是否migrate accounts。",
						},
						"is_override_root": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "是否use Root 账号 在 来源 数据库 到 overwrite 该 在 目标 数据库. 有效值：false，true. For 数据库/表 或 structural 迁移，您 should 指定false. 注意 该 此 参数 takes effect 仅 对于 OldDTS。",
						},
						"is_dst_read_only": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "是否set 目标 数据库 到 read-仅 during 迁移，其中 takes effect 仅 对于 MySQL databases. 有效值：true，false. 默认值：false。",
						},
						"extra_attr": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Additional 信息. You 可以 集合 additional 参数 对于 certain 数据库 types。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Option 键",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Option 值",
									},
								},
							},
						},
					},
				},
			},

			"src_info": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "来源 实例 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "实例 地域",
						},
						"access_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instances 网络 访问 类型 有效值：extranet (公有 网络); ipv6 (公有 IPv6); cvm (self-build 在 CVM); dcg (Direct Connect); vpncloud (VPN 访问); cdb (数据库); ccn (CCN); intranet (intranet); vpc (VPC). 注意 该 有效 值 是 subject 到 当前 link。",
						},
						"database_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database 类型，such 作为 mysql，redis，mongodb，postgresql，mariadb，和 percona。",
						},
						"node_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Node 类型，空 或 simple 表示a general 节点，集群 表示a 集群 节点; 对于 mongo services，有效值：replicaset (mongodb 副本 集合)，standalone (mongodb 单个 节点)，集群 (mongodb 集群); 对于 redis 实例，有效值：空 或 simple (单个 节点)，集群 (集群)，集群-缓存 (缓存 集群)，集群-proxy (proxy 集群)。",
						},
						"info": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Database 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Node 角色 在 distributed 数据库，such 作为 mongos 节点 在 MongoDB。",
									},
									"db_kernel": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Kernel 版本，such 作为 different kernel versions 的 MariaDB。",
									},
									"host": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "实例 IP 地址，其中 为必填项 对于 following 访问 types: 公有 网络，Direct Connect，VPN，CCN，intranet，和 VPC。",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "实例端口，其中 为必填项 对于 following 访问 types: 公有 网络，self-build 在 CVM，Direct Connect，VPN，CCN，intranet，和 VPC。",
									},
									"user": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "实例 用户名",
									},
									"password": {
										Type:        schema.TypeString,
										Optional:    true,
										Sensitive:   true,
										Description: "实例 密码",
									},
									"cvm_instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Short CVM 实例 ID 在 格式 的 ins-olgl39y8，其中 为必填项 如果 访问 类型 是 cvm. It 是 same 作为 实例 ID displayed 在 CVM console。",
									},
									"uniq_vpn_gw_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "VPN 网关 ID 在 格式 的 vpngw-9ghexg7q，其中 为必填项 如果 访问 类型 是 vpncloud。",
									},
									"uniq_dcg_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Direct Connect 网关 ID 在 格式 的 dcg-0rxtqqxb，其中 为必填项 如果 访问 类型 是 dcg。",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Database 实例 ID 在 格式 的 cdb-powiqx8q，其中 为必填项 如果 访问 类型 是 cdb。",
									},
									"ccn_gw_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "CCN 实例 ID such 作为 ccn-afp6kltc。",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "私有网络 ID 在 格式 的 vpc-92jblxto，其中 为必填项 如果 访问 类型 是 vpc，vpncloud，ccn，或 dcg。",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ID 子网 在 VPC 在 格式 的 子网-3paxmkdz，其中 为必填项 如果 访问 类型 是 vpc，vpncloud，ccn，或 dcg。",
									},
									"engine_version": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Database 版本 在 格式 的 5.6 或 5.7，其中 takes effect 仅 如果 实例 是 RDS 实例. 默认值：5.6。",
									},
									"account": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "实例 账号",
									},
									"account_role": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "角色 用于cross-账号 迁移，其中 可以 contain [-zA-Z0-9-_]+。",
									},
									"account_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "账号 到 其中 资源 belongs. 有效值：空 或 self ( 当前 账号); other (another 账号)。",
									},
									"tmp_secret_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Temporary SecretId，您 可以 obtain temporary 键 通过 GetFederationToken。",
									},
									"tmp_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Temporary SecretKey，您 可以 obtain temporary 键 通过 GetFederationToken。",
									},
									"tmp_token": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Temporary 令牌，您 可以 obtain temporary 键 通过 GetFederationToken。",
									},
								},
							},
						},
						"supplier": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "实例 服务 provider，such 作为 `aliyun` 和 `others`。",
						},
						"extra_attr": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "For MongoDB，您 可以 define following 参数: ['AuthDatabase':'admin'，'AuthFlag': '1'，'AuthMechanism':'SCRAM-SHA-1']。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Option 键",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Option 值",
									},
								},
							},
						},
					},
				},
			},

			"dst_info": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Target 数据库 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "实例 地域",
						},
						"access_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instances 网络 访问 类型 有效值：extranet (公有 网络); ipv6 (公有 IPv6); cvm (self-build 在 CVM); dcg (Direct Connect); vpncloud (VPN 访问); cdb (数据库); ccn (CCN); intranet (intranet); vpc (VPC). 注意 该 有效 值 是 subject 到 当前 link。",
						},
						"database_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database 类型，such 作为 mysql，redis，mongodb，postgresql，mariadb，和 percona。",
						},
						"node_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Node 类型，空 或 simple 表示a general 节点，集群 表示a 集群 节点; 对于 mongo services，有效值：replicaset (mongodb 副本 集合)，standalone (mongodb 单个 节点)，集群 (mongodb 集群); 对于 redis 实例，有效值：空 或 simple (单个 节点)，集群 (集群)，集群-缓存 (缓存 集群)，集群-proxy (proxy 集群)。",
						},
						"info": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Database 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Node 角色 在 distributed 数据库，such 作为 mongos 节点 在 MongoDB。",
									},
									"db_kernel": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Kernel 版本，such 作为 different kernel versions 的 MariaDB。",
									},
									"host": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "实例 IP 地址，其中 为必填项 对于 following 访问 types: 公有 网络，Direct Connect，VPN，CCN，intranet，和 VPC。",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "实例端口，其中 为必填项 对于 following 访问 types: 公有 网络，self-build 在 CVM，Direct Connect，VPN，CCN，intranet，和 VPC。",
									},
									"user": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "实例 用户名",
									},
									"password": {
										Type:        schema.TypeString,
										Optional:    true,
										Sensitive:   true,
										Description: "实例 密码",
									},
									"cvm_instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Short CVM 实例 ID 在 格式 的 ins-olgl39y8，其中 为必填项 如果 访问 类型 是 cvm. It 是 same 作为 实例 ID displayed 在 CVM console。",
									},
									"uniq_vpn_gw_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "VPN 网关 ID 在 格式 的 vpngw-9ghexg7q，其中 为必填项 如果 访问 类型 是 vpncloud。",
									},
									"uniq_dcg_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Direct Connect 网关 ID 在 格式 的 dcg-0rxtqqxb，其中 为必填项 如果 访问 类型 是 dcg。",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Database 实例 ID 在 格式 的 cdb-powiqx8q，其中 为必填项 如果 访问 类型 是 cdb。",
									},
									"ccn_gw_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "CCN 实例 ID such 作为 ccn-afp6kltc。",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "私有网络 ID 在 格式 的 vpc-92jblxto，其中 为必填项 如果 访问 类型 是 vpc，vpncloud，ccn，或 dcg。",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "ID 子网 在 VPC 在 格式 的 子网-3paxmkdz，其中 为必填项 如果 访问 类型 是 vpc，vpncloud，ccn，或 dcg。",
									},
									"engine_version": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Database 版本 在 格式 的 5.6 或 5.7，其中 takes effect 仅 如果 实例 是 RDS 实例. 默认值：5.6。",
									},
									"account": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "实例 账号",
									},
									"account_role": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "角色 用于cross-账号 迁移，其中 可以 contain [-zA-Z0-9-_]+。",
									},
									"account_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "账号 到 其中 资源 belongs. 有效值：空 或 self ( 当前 账号); other (another 账号)。",
									},
									"tmp_secret_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Temporary SecretId，您 可以 obtain temporary 键 通过 GetFederationToken。",
									},
									"tmp_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Temporary SecretKey，您 可以 obtain temporary 键 通过 GetFederationToken。",
									},
									"tmp_token": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Temporary 令牌，您 可以 obtain temporary 键 通过 GetFederationToken。",
									},
								},
							},
						},
						"supplier": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "实例 服务 provider，such 作为 `aliyun` 和 `others`。",
						},
						"extra_attr": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "For MongoDB，您 可以 define following 参数: ['AuthDatabase':'admin','AuthFlag': '1'，'AuthMechanism':'SCRAM-SHA-1']。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Option 键",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Option 值",
									},
								},
							},
						},
					},
				},
			},

			"expect_run_time": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Expected 开始时间 在 格式 的 `2006-01-02 15:04:05`，其中 为必填项 如果 RunMode 是 timed。",
			},

			"auto_retry_time_range_minutes": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "automatic retry 时间 周期 可以 是 集合 从 5 到 720 minutes，使用 0 indicating 无 retry。",
			},
		},
	}
}

func resourceTencentCloudDtsMigrateJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_migrate_job.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		tcClient  = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		service   = DtsService{client: tcClient}
		conf      *resource.StateChangeConf
		serviceId string
	)

	if v, ok := d.GetOk("service_id"); ok {
		serviceId = v.(string)
	}

	// case "modify":
	err := handleModifyMigrate(d, tcClient, logId, serviceId)
	if err != nil {
		return err
	}

	conf = tccommon.BuildStateChangeConf([]string{}, []string{"created"}, 3*tccommon.ReadRetryTimeout, time.Second, service.DtsMigrateJobStateRefreshFunc(serviceId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	// case "check":
	err = handleCheckMigrate(d, tcClient, logId, serviceId)
	if err != nil {
		return err
	}

	conf = tccommon.BuildStateChangeConf([]string{}, []string{"checkPass", "checkNotPass"}, 3*tccommon.ReadRetryTimeout, time.Second, service.DtsMigrateCheckConfigStateRefreshFunc(serviceId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	d.SetId(serviceId)
	return resourceTencentCloudDtsMigrateJobRead(d, meta)
}

func resourceTencentCloudDtsMigrateJobUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_migrate_job.update")()
	defer tccommon.InconsistentCheck(d, meta)()
	logId := tccommon.GetLogId(tccommon.ContextNil)

	log.Printf("[DEBUG]%s tencentcloud_dts_migrate_job.update in. id:[%s]\n", logId, d.Id())

	return resourceTencentCloudDtsMigrateJobCreate(d, meta)
}

func resourceTencentCloudDtsMigrateJobRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_migrate_job.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DtsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	jobId := d.Id()
	log.Printf("[DEBUG]%s tencentcloud_dts_migrate_job.read trying to call DescribeDtsMigrateJobById. jobId:[%s]\n", logId, jobId)
	migrateJob, err := service.DescribeDtsMigrateJobById(ctx, jobId)
	if err != nil {
		return err
	}

	if migrateJob == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	if migrateJob.JobId != nil {
		_ = d.Set("service_id", migrateJob.JobId)
	}

	if migrateJob.Status != nil {
		_ = d.Set("status", migrateJob.Status)
	}

	// for modify operation
	if migrateJob.RunMode != nil {
		_ = d.Set("run_mode", migrateJob.RunMode)
	}

	if migrateJob.MigrateOption != nil {
		migrateOptionMap := make(map[string]interface{})

		if migrateJob.MigrateOption.DatabaseTable != nil {
			databaseTableMap := make(map[string]interface{})

			if migrateJob.MigrateOption.DatabaseTable.ObjectMode != nil {
				databaseTableMap["object_mode"] = migrateJob.MigrateOption.DatabaseTable.ObjectMode
			}

			if migrateJob.MigrateOption.DatabaseTable.Databases != nil {
				databasesList := make([]interface{}, 0, len(migrateJob.MigrateOption.DatabaseTable.Databases))
				for _, databases := range migrateJob.MigrateOption.DatabaseTable.Databases {
					databasesMap := make(map[string]interface{})

					if databases.DbName != nil {
						databasesMap["db_name"] = databases.DbName
					}

					if databases.NewDbName != nil {
						databasesMap["new_db_name"] = databases.NewDbName
					}

					if databases.SchemaName != nil {
						databasesMap["schema_name"] = databases.SchemaName
					}

					if databases.NewSchemaName != nil {
						databasesMap["new_schema_name"] = databases.NewSchemaName
					}

					if databases.DBMode != nil {
						databasesMap["db_mode"] = databases.DBMode
					}

					if databases.SchemaMode != nil {
						databasesMap["schema_mode"] = databases.SchemaMode
					}

					if databases.TableMode != nil {
						databasesMap["table_mode"] = databases.TableMode
					}

					if databases.Tables != nil {
						tablesList := make([]interface{}, 0, len(databases.Tables))
						for _, tables := range databases.Tables {
							tablesMap := make(map[string]interface{})

							if tables.TableName != nil {
								tablesMap["table_name"] = tables.TableName
							}

							if tables.NewTableName != nil {
								tablesMap["new_table_name"] = tables.NewTableName
							}

							if tables.TmpTables != nil {
								tablesMap["tmp_tables"] = tables.TmpTables
							}

							if tables.TableEditMode != nil {
								tablesMap["table_edit_mode"] = tables.TableEditMode
							}

							tablesList = append(tablesList, tablesMap)
						}

						databasesMap["tables"] = tablesList
					}

					if databases.ViewMode != nil {
						databasesMap["view_mode"] = databases.ViewMode
					}

					if databases.Views != nil {
						viewsList := make([]interface{}, 0, len(databases.Views))
						for _, views := range databases.Views {
							viewsMap := make(map[string]interface{})

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

					if databases.RoleMode != nil {
						databasesMap["role_mode"] = databases.RoleMode
					}

					if databases.Roles != nil {
						rolesList := make([]interface{}, 0, len(databases.Roles))
						for _, roles := range databases.Roles {
							rolesMap := make(map[string]interface{})

							if roles.RoleName != nil {
								rolesMap["role_name"] = roles.RoleName
							}

							if roles.NewRoleName != nil {
								rolesMap["new_role_name"] = roles.NewRoleName
							}

							rolesList = append(rolesList, rolesMap)
						}

						databasesMap["roles"] = rolesList
					}

					if databases.FunctionMode != nil {
						databasesMap["function_mode"] = databases.FunctionMode
					}

					if databases.TriggerMode != nil {
						databasesMap["trigger_mode"] = databases.TriggerMode
					}

					if databases.EventMode != nil {
						databasesMap["event_mode"] = databases.EventMode
					}

					if databases.ProcedureMode != nil {
						databasesMap["procedure_mode"] = databases.ProcedureMode
					}

					log.Printf("[DEBUG]%s read databases.Functions:[%v],len[%d]", logId, databases.Functions, len(databases.Functions))
					for _, fun := range databases.Functions {
						log.Printf("[DEBUG]%s read databases.Functions: iterate fun:[%s]", logId, *fun)
					}

					if databases.Functions != nil {
						databasesMap["functions"] = databases.Functions
						log.Printf("[DEBUG]%s read databases.Functions: i'm in. databasesMap:[%v]", logId, databasesMap["functions"])
					}

					if databases.Procedures != nil {
						databasesMap["procedures"] = databases.Procedures
					}

					if databases.Events != nil {
						databasesMap["events"] = databases.Events
					}

					if databases.Triggers != nil {
						databasesMap["triggers"] = databases.Triggers
					}

					databasesList = append(databasesList, databasesMap)
				}

				// databaseTableMap["databases"] = []interface{}{databasesList}
				databaseTableMap["databases"] = databasesList
			}

			if migrateJob.MigrateOption.DatabaseTable.AdvancedObjects != nil {
				databaseTableMap["advanced_objects"] = migrateJob.MigrateOption.DatabaseTable.AdvancedObjects
			}

			migrateOptionMap["database_table"] = []interface{}{databaseTableMap}
		}

		if migrateJob.MigrateOption.MigrateType != nil {
			migrateOptionMap["migrate_type"] = migrateJob.MigrateOption.MigrateType
		}

		log.Printf("[DEBUG]%s read  migrateJob.MigrateOption.Consistency:[%v]", logId, migrateJob.MigrateOption.Consistency)
		if migrateJob.MigrateOption.Consistency != nil {
			consistencyMap := make(map[string]interface{})

			mode := migrateJob.MigrateOption.Consistency.Mode
			if mode != nil && *mode != "" {
				consistencyMap["mode"] = migrateJob.MigrateOption.Consistency.Mode
			}
			migrateOptionMap["consistency"] = []interface{}{consistencyMap}
		}

		if migrateJob.MigrateOption.IsMigrateAccount != nil {
			migrateOptionMap["is_migrate_account"] = migrateJob.MigrateOption.IsMigrateAccount
		}

		if migrateJob.MigrateOption.IsOverrideRoot != nil {
			migrateOptionMap["is_override_root"] = migrateJob.MigrateOption.IsOverrideRoot
		}

		if migrateJob.MigrateOption.IsDstReadOnly != nil {
			migrateOptionMap["is_dst_read_only"] = migrateJob.MigrateOption.IsDstReadOnly
		}

		if migrateJob.MigrateOption.ExtraAttr != nil {
			extraAttrList := make([]interface{}, 0, len(migrateJob.MigrateOption.ExtraAttr))
			for _, extraAttr := range migrateJob.MigrateOption.ExtraAttr {
				extraAttrMap := make(map[string]interface{})

				if extraAttr.Key != nil {
					extraAttrMap["key"] = extraAttr.Key
				}

				if extraAttr.Value != nil {
					extraAttrMap["value"] = extraAttr.Value
				}

				extraAttrList = append(extraAttrList, extraAttrMap)
			}

			migrateOptionMap["extra_attr"] = extraAttrList
		}

		_ = d.Set("migrate_option", []interface{}{migrateOptionMap})
	}

	if migrateJob.SrcInfo != nil {
		srcInfoMap := make(map[string]interface{})

		if migrateJob.SrcInfo.Region != nil {
			srcInfoMap["region"] = migrateJob.SrcInfo.Region
		}

		if migrateJob.SrcInfo.AccessType != nil {
			srcInfoMap["access_type"] = migrateJob.SrcInfo.AccessType
		}

		if migrateJob.SrcInfo.DatabaseType != nil {
			srcInfoMap["database_type"] = migrateJob.SrcInfo.DatabaseType
		}

		if migrateJob.SrcInfo.NodeType != nil {
			srcInfoMap["node_type"] = migrateJob.SrcInfo.NodeType
		}

		if migrateJob.SrcInfo.Info != nil {
			infoList := make([]interface{}, 0, len(migrateJob.SrcInfo.Info))
			for i, info := range migrateJob.SrcInfo.Info {
				infoMap := make(map[string]interface{})

				if info.Password == nil || *info.Password == "" {
					//reset password
					key := fmt.Sprintf("src_info.0.info.%v.password", i)
					if v, ok := d.GetOk(key); ok {
						infoMap["password"] = helper.String(v.(string))
						log.Printf("[DEBUG]%s set src_info.0.info.%v.password:[key:%s]", logId, i, key)
					}
				} else {
					infoMap["password"] = info.Password
				}

				if info.Role != nil {
					infoMap["role"] = info.Role
				}

				if info.DbKernel != nil {
					infoMap["db_kernel"] = info.DbKernel
				}

				if info.Host != nil {
					infoMap["host"] = info.Host
				}

				if info.Port != nil {
					infoMap["port"] = info.Port
				}

				if info.User != nil {
					infoMap["user"] = info.User
				}

				if info.CvmInstanceId != nil {
					infoMap["cvm_instance_id"] = info.CvmInstanceId
				}

				if info.UniqVpnGwId != nil {
					infoMap["uniq_vpn_gw_id"] = info.UniqVpnGwId
				}

				if info.UniqDcgId != nil {
					infoMap["uniq_dcg_id"] = info.UniqDcgId
				}

				if info.InstanceId != nil {
					infoMap["instance_id"] = info.InstanceId
				}

				if info.CcnGwId != nil {
					infoMap["ccn_gw_id"] = info.CcnGwId
				}

				if info.VpcId != nil {
					infoMap["vpc_id"] = info.VpcId
				}

				if info.SubnetId != nil {
					infoMap["subnet_id"] = info.SubnetId
				}

				log.Printf("[DEBUG]%s read  migrateJob.SrcInfo.Info.EngineVersion:[%v,%s]", logId, info.EngineVersion, *info.EngineVersion)
				if info.EngineVersion != nil {
					infoMap["engine_version"] = info.EngineVersion
				}

				if info.Account != nil {
					infoMap["account"] = info.Account
				}

				if info.AccountRole != nil {
					infoMap["account_role"] = info.AccountRole
				}

				if info.AccountMode != nil {
					infoMap["account_mode"] = info.AccountMode
				}

				if info.TmpSecretId != nil {
					infoMap["tmp_secret_id"] = info.TmpSecretId
				}

				if info.TmpSecretKey != nil {
					infoMap["tmp_secret_key"] = info.TmpSecretKey
				}

				if info.TmpToken != nil {
					infoMap["tmp_token"] = info.TmpToken
				}

				infoList = append(infoList, infoMap)
			}

			srcInfoMap["info"] = infoList
		}

		if migrateJob.SrcInfo.Supplier != nil {
			srcInfoMap["supplier"] = migrateJob.SrcInfo.Supplier
		}

		if migrateJob.SrcInfo.ExtraAttr != nil {
			extraAttrList := make([]interface{}, 0, len(migrateJob.SrcInfo.ExtraAttr))
			for _, extraAttr := range migrateJob.SrcInfo.ExtraAttr {
				extraAttrMap := make(map[string]interface{})

				if extraAttr.Key != nil {
					extraAttrMap["key"] = extraAttr.Key
				}

				if extraAttr.Value != nil {
					extraAttrMap["value"] = extraAttr.Value
				}

				extraAttrList = append(extraAttrList, extraAttrMap)
			}

			srcInfoMap["extra_attr"] = extraAttrList
		}

		_ = d.Set("src_info", []interface{}{srcInfoMap})
	}

	if migrateJob.DstInfo != nil {
		dstInfoMap := make(map[string]interface{})
		if migrateJob.DstInfo.Region != nil {
			dstInfoMap["region"] = migrateJob.DstInfo.Region
		}

		if migrateJob.DstInfo.AccessType != nil {
			dstInfoMap["access_type"] = migrateJob.DstInfo.AccessType
		}

		if migrateJob.DstInfo.DatabaseType != nil {
			dstInfoMap["database_type"] = migrateJob.DstInfo.DatabaseType
		}

		if migrateJob.DstInfo.NodeType != nil {
			dstInfoMap["node_type"] = migrateJob.DstInfo.NodeType
		}

		log.Printf("[DEBUG]%s read migrateJob.DstInfo.Info :[%v], len:[%v]", logId, migrateJob.DstInfo.Info, len(migrateJob.DstInfo.Info))
		if migrateJob.DstInfo.Info != nil {
			infoList := make([]interface{}, 0, len(migrateJob.DstInfo.Info))
			for i, info := range migrateJob.DstInfo.Info {
				infoMap := make(map[string]interface{})

				if info.Password == nil || *info.Password == "" {
					//reset password
					key := fmt.Sprintf("dst_info.0.info.%v.password", i)
					if v, ok := d.GetOk(key); ok {
						infoMap["password"] = helper.String(v.(string))
						log.Printf("[DEBUG]%s set dst_info.0.info.%v.password:[key:%s]", logId, i, key)
					}
				} else {
					infoMap["password"] = info.Password
				}

				if info.Role != nil {
					infoMap["role"] = info.Role
				}

				if info.DbKernel != nil {
					infoMap["db_kernel"] = info.DbKernel
				}

				if info.Host != nil {
					infoMap["host"] = info.Host
				}

				if info.Port != nil {
					infoMap["port"] = info.Port
				}

				if info.User != nil {
					infoMap["user"] = info.User
				}

				if info.CvmInstanceId != nil {
					infoMap["cvm_instance_id"] = info.CvmInstanceId
				}

				if info.UniqVpnGwId != nil {
					infoMap["uniq_vpn_gw_id"] = info.UniqVpnGwId
				}

				if info.UniqDcgId != nil {
					infoMap["uniq_dcg_id"] = info.UniqDcgId
				}

				if info.InstanceId != nil {
					infoMap["instance_id"] = info.InstanceId
				}

				if info.CcnGwId != nil {
					infoMap["ccn_gw_id"] = info.CcnGwId
				}

				if info.VpcId != nil {
					infoMap["vpc_id"] = info.VpcId
				}

				if info.SubnetId != nil {
					infoMap["subnet_id"] = info.SubnetId
				}

				log.Printf("[DEBUG]%s read  migrateJob.DstInfo.Info.EngineVersion:[%v,%s]", logId, info.EngineVersion, *info.EngineVersion)
				if d.HasChange("engine_version") {
					if info.EngineVersion != nil {
						infoMap["engine_version"] = info.EngineVersion
					}
				}

				if info.Account != nil {
					infoMap["account"] = info.Account
				}

				if info.AccountRole != nil {
					infoMap["account_role"] = info.AccountRole
				}

				if info.AccountMode != nil {
					infoMap["account_mode"] = info.AccountMode
				}

				if info.TmpSecretId != nil {
					infoMap["tmp_secret_id"] = info.TmpSecretId
				}

				if info.TmpSecretKey != nil {
					infoMap["tmp_secret_key"] = info.TmpSecretKey
				}

				if info.TmpToken != nil {
					infoMap["tmp_token"] = info.TmpToken
				}

				infoList = append(infoList, infoMap)
			}

			dstInfoMap["info"] = infoList
		}

		log.Printf("[DEBUG]%s read migrateJob.DstInfo.Supplier :[%s]", logId, *migrateJob.DstInfo.Supplier)
		if migrateJob.DstInfo.Supplier != nil {
			dstInfoMap["supplier"] = migrateJob.DstInfo.Supplier
		}

		if migrateJob.DstInfo.ExtraAttr != nil {
			extraAttrList := make([]interface{}, 0, len(migrateJob.DstInfo.ExtraAttr))
			for _, extraAttr := range migrateJob.DstInfo.ExtraAttr {
				extraAttrMap := make(map[string]interface{})

				if extraAttr.Key != nil {
					extraAttrMap["key"] = extraAttr.Key
				}

				if extraAttr.Value != nil {
					extraAttrMap["value"] = extraAttr.Value
				}

				extraAttrList = append(extraAttrList, extraAttrMap)
			}

			dstInfoMap["extra_attr"] = extraAttrList
		}

		_ = d.Set("dst_info", []interface{}{dstInfoMap})
	}

	if migrateJob.ExpectRunTime != nil {
		_ = d.Set("expect_run_time", migrateJob.ExpectRunTime)
	}

	return nil
}

func handleModifyMigrate(d *schema.ResourceData, tcClient *connectivity.TencentCloudClient, logId, jobId string) error {
	configMigrationJobRequest := dts.NewModifyMigrationJobRequest()
	configMigrationJobRequest.JobId = helper.String(jobId)

	if v, ok := d.GetOk("run_mode"); ok {
		configMigrationJobRequest.RunMode = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "migrate_option"); ok {
		migrateOption := dts.MigrateOption{}
		if databaseTableMap, ok := helper.InterfaceToMap(dMap, "database_table"); ok {
			databaseTableObject := dts.DatabaseTableObject{}
			if v, ok := databaseTableMap["object_mode"]; ok && v.(string) != "" {
				databaseTableObject.ObjectMode = helper.String(v.(string))
			}
			if v, ok := databaseTableMap["databases"]; ok {
				for _, item := range v.([]interface{}) {
					databasesMap := item.(map[string]interface{})
					dBItem := dts.DBItem{}
					if v, ok := databasesMap["db_name"]; ok && v.(string) != "" {
						dBItem.DbName = helper.String(v.(string))
					}
					if v, ok := databasesMap["new_db_name"]; ok && v.(string) != "" {
						dBItem.NewDbName = helper.String(v.(string))
					}
					if v, ok := databasesMap["schema_name"]; ok && v.(string) != "" {
						dBItem.SchemaName = helper.String(v.(string))
					}
					if v, ok := databasesMap["new_schema_name"]; ok && v.(string) != "" {
						dBItem.NewSchemaName = helper.String(v.(string))
					}
					if v, ok := databasesMap["db_mode"]; ok && v.(string) != "" {
						dBItem.DBMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["schema_mode"]; ok && v.(string) != "" {
						dBItem.SchemaMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["table_mode"]; ok && v.(string) != "" {
						dBItem.TableMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["tables"]; ok {
						for _, item := range v.([]interface{}) {
							tablesMap := item.(map[string]interface{})
							tableItem := dts.TableItem{}
							if v, ok := tablesMap["table_name"]; ok && v.(string) != "" {
								tableItem.TableName = helper.String(v.(string))
							}
							if v, ok := tablesMap["new_table_name"]; ok && v.(string) != "" {
								tableItem.NewTableName = helper.String(v.(string))
							}
							if v, ok := tablesMap["tmp_tables"]; ok {
								tmpTablesSet := v.(*schema.Set).List()
								for i := range tmpTablesSet {
									tmpTables := tmpTablesSet[i].(string)
									tableItem.TmpTables = append(tableItem.TmpTables, &tmpTables)
								}
							}
							if v, ok := tablesMap["table_edit_mode"]; ok && v.(string) != "" {
								tableItem.TableEditMode = helper.String(v.(string))
							}
							dBItem.Tables = append(dBItem.Tables, &tableItem)
						}
					}
					if v, ok := databasesMap["view_mode"]; ok && v.(string) != "" {
						dBItem.ViewMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["views"]; ok {
						for _, item := range v.([]interface{}) {
							viewsMap := item.(map[string]interface{})
							viewItem := dts.ViewItem{}
							if v, ok := viewsMap["view_name"]; ok && v.(string) != "" {
								viewItem.ViewName = helper.String(v.(string))
							}
							if v, ok := viewsMap["new_view_name"]; ok && v.(string) != "" {
								viewItem.NewViewName = helper.String(v.(string))
							}
							dBItem.Views = append(dBItem.Views, &viewItem)
						}
					}
					if v, ok := databasesMap["role_mode"]; ok && v.(string) != "" {
						dBItem.RoleMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["roles"]; ok {
						for _, item := range v.([]interface{}) {
							rolesMap := item.(map[string]interface{})
							roleItem := dts.RoleItem{}
							if v, ok := rolesMap["role_name"]; ok && v.(string) != "" {
								roleItem.RoleName = helper.String(v.(string))
							}
							if v, ok := rolesMap["new_role_name"]; ok && v.(string) != "" {
								roleItem.NewRoleName = helper.String(v.(string))
							}
							dBItem.Roles = append(dBItem.Roles, &roleItem)
						}
					}
					if v, ok := databasesMap["function_mode"]; ok && v.(string) != "" {
						dBItem.FunctionMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["trigger_mode"]; ok && v.(string) != "" {
						dBItem.TriggerMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["event_mode"]; ok && v.(string) != "" {
						dBItem.EventMode = helper.String(v.(string))
					}
					if v, ok := databasesMap["procedure_mode"]; ok && v.(string) != "" {
						dBItem.ProcedureMode = helper.String(v.(string))
					}
					log.Printf("[DEBUG]%s modify databases.Functions: databasesMap[\"functions\"]:[%v]", logId, databasesMap["functions"])
					if v, ok := databasesMap["functions"]; ok {
						functionsSet := v.(*schema.Set).List()
						log.Printf("[DEBUG]%s modify databases.Functions: i'm in. functionsSet:[%v]", logId, functionsSet)
						for _, funcc := range functionsSet {
							functions := funcc.(*string)
							dBItem.Functions = append(dBItem.Functions, functions)
							log.Printf("[DEBUG]%s modify databases.Functions: iterate functions:[%s]", logId, *functions)
						}
					}
					if v, ok := databasesMap["procedures"]; ok {
						proceduresSet := v.(*schema.Set).List()
						for _, proc := range proceduresSet {
							procedures := proc.(string)
							dBItem.Procedures = append(dBItem.Procedures, &procedures)
						}
					}
					if v, ok := databasesMap["events"]; ok {
						eventsSet := v.(*schema.Set).List()
						for i := range eventsSet {
							events := eventsSet[i].(string)
							dBItem.Events = append(dBItem.Events, &events)
						}
					}
					if v, ok := databasesMap["triggers"]; ok {
						triggersSet := v.(*schema.Set).List()
						for i := range triggersSet {
							triggers := triggersSet[i].(string)
							dBItem.Triggers = append(dBItem.Triggers, &triggers)
						}
					}
					databaseTableObject.Databases = append(databaseTableObject.Databases, &dBItem)
				}
			}
			if v, ok := databaseTableMap["advanced_objects"]; ok {
				advancedObjectsSet := v.(*schema.Set).List()
				for i := range advancedObjectsSet {
					advancedObjects := advancedObjectsSet[i].(string)
					databaseTableObject.AdvancedObjects = append(databaseTableObject.AdvancedObjects, &advancedObjects)
				}
			}
			migrateOption.DatabaseTable = &databaseTableObject
		}
		if v, ok := dMap["migrate_type"]; ok && v.(string) != "" {
			migrateOption.MigrateType = helper.String(v.(string))
		}
		log.Printf("[DEBUG]%s update  migrateJob.MigrateOption.Consistency dMap(consistency):[%v]", logId, dMap["consistency"])
		if consistencyMap, ok := helper.InterfaceToMap(dMap, "consistency"); ok {
			log.Printf("[DEBUG]%s update  migrateJob.MigrateOption.Consistency:[%v]", logId, consistencyMap)
			consistencyOption := dts.ConsistencyOption{}
			if v, ok := consistencyMap["mode"]; ok && v.(string) != "" {
				consistencyOption.Mode = helper.String(v.(string))
			}
			migrateOption.Consistency = &consistencyOption
		}
		if v, ok := dMap["is_migrate_account"]; ok {
			migrateOption.IsMigrateAccount = helper.Bool(v.(bool))
		}
		if v, ok := dMap["is_override_root"]; ok {
			migrateOption.IsOverrideRoot = helper.Bool(v.(bool))
		}
		if v, ok := dMap["is_dst_read_only"]; ok {
			migrateOption.IsDstReadOnly = helper.Bool(v.(bool))
		}
		if v, ok := dMap["extra_attr"]; ok {
			for _, item := range v.([]interface{}) {
				extraAttrMap := item.(map[string]interface{})
				keyValuePairOption := dts.KeyValuePairOption{}
				if v, ok := extraAttrMap["key"]; ok && v.(string) != "" {
					keyValuePairOption.Key = helper.String(v.(string))
				}
				if v, ok := extraAttrMap["value"]; ok && v.(string) != "" {
					keyValuePairOption.Value = helper.String(v.(string))
				}
				migrateOption.ExtraAttr = append(migrateOption.ExtraAttr, &keyValuePairOption)
			}
		}
		configMigrationJobRequest.MigrateOption = &migrateOption
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "src_info"); ok {
		dBEndpointInfo := dts.DBEndpointInfo{}
		if v, ok := dMap["region"]; ok && v.(string) != "" {
			dBEndpointInfo.Region = helper.String(v.(string))
		}
		if v, ok := dMap["access_type"]; ok && v.(string) != "" {
			dBEndpointInfo.AccessType = helper.String(v.(string))
		}
		if v, ok := dMap["database_type"]; ok && v.(string) != "" {
			dBEndpointInfo.DatabaseType = helper.String(v.(string))
		}
		if v, ok := dMap["node_type"]; ok && v.(string) != "" {
			dBEndpointInfo.NodeType = helper.String(v.(string))
		}
		if v, ok := dMap["info"]; ok {
			for _, item := range v.([]interface{}) {
				srcInfoMap := item.(map[string]interface{})
				dBInfo := dts.DBInfo{}
				if v, ok := srcInfoMap["role"]; ok && v.(string) != "" {
					dBInfo.Role = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["db_kernel"]; ok && v.(string) != "" {
					dBInfo.DbKernel = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["host"]; ok && v.(string) != "" {
					dBInfo.Host = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["port"]; ok {
					dBInfo.Port = helper.IntUint64(v.(int))
				}
				if v, ok := srcInfoMap["user"]; ok && v.(string) != "" {
					dBInfo.User = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["password"]; ok && v.(string) != "" {
					dBInfo.Password = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["cvm_instance_id"]; ok && v.(string) != "" {
					dBInfo.CvmInstanceId = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["uniq_vpn_gw_id"]; ok && v.(string) != "" {
					dBInfo.UniqVpnGwId = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["uniq_dcg_id"]; ok && v.(string) != "" {
					dBInfo.UniqDcgId = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["instance_id"]; ok && v.(string) != "" {
					dBInfo.InstanceId = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["ccn_gw_id"]; ok && v.(string) != "" {
					dBInfo.CcnGwId = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["vpc_id"]; ok && v.(string) != "" {
					dBInfo.VpcId = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["subnet_id"]; ok && v.(string) != "" {
					dBInfo.SubnetId = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["engine_version"]; ok && v.(string) != "" {
					dBInfo.EngineVersion = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["account"]; ok && v.(string) != "" {
					dBInfo.Account = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["account_role"]; ok && v.(string) != "" {
					dBInfo.AccountRole = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["account_mode"]; ok && v.(string) != "" {
					dBInfo.AccountMode = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["tmp_secret_id"]; ok && v.(string) != "" {
					dBInfo.TmpSecretId = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["tmp_secret_key"]; ok && v.(string) != "" {
					dBInfo.TmpSecretKey = helper.String(v.(string))
				}
				if v, ok := srcInfoMap["tmp_token"]; ok && v.(string) != "" {
					dBInfo.TmpToken = helper.String(v.(string))
				}
				dBEndpointInfo.Info = append(dBEndpointInfo.Info, &dBInfo)
			}
		}
		if v, ok := dMap["supplier"]; ok && v.(string) != "" {
			dBEndpointInfo.Supplier = helper.String(v.(string))
		}
		if v, ok := dMap["extra_attr"]; ok {
			for _, item := range v.([]interface{}) {
				extraAttrMap := item.(map[string]interface{})
				keyValuePairOption := dts.KeyValuePairOption{}
				if v, ok := extraAttrMap["key"]; ok && v.(string) != "" {
					keyValuePairOption.Key = helper.String(v.(string))
				}
				if v, ok := extraAttrMap["value"]; ok && v.(string) != "" {
					keyValuePairOption.Value = helper.String(v.(string))
				}
				dBEndpointInfo.ExtraAttr = append(dBEndpointInfo.ExtraAttr, &keyValuePairOption)
			}
		}
		configMigrationJobRequest.SrcInfo = &dBEndpointInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "dst_info"); ok {
		dBEndpointInfo := dts.DBEndpointInfo{}
		if v, ok := dMap["region"]; ok && v.(string) != "" {
			dBEndpointInfo.Region = helper.String(v.(string))
		}
		if v, ok := dMap["access_type"]; ok && v.(string) != "" {
			dBEndpointInfo.AccessType = helper.String(v.(string))
		}
		if v, ok := dMap["database_type"]; ok && v.(string) != "" {
			dBEndpointInfo.DatabaseType = helper.String(v.(string))
		}
		if v, ok := dMap["node_type"]; ok && v.(string) != "" {
			dBEndpointInfo.NodeType = helper.String(v.(string))
		}
		if v, ok := dMap["info"]; ok {
			for _, item := range v.([]interface{}) {
				dstInfoMap := item.(map[string]interface{})
				dBInfo := dts.DBInfo{}
				if v, ok := dstInfoMap["role"]; ok && v.(string) != "" {
					dBInfo.Role = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["db_kernel"]; ok && v.(string) != "" {
					dBInfo.DbKernel = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["host"]; ok && v.(string) != "" {
					dBInfo.Host = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["port"]; ok {
					dBInfo.Port = helper.IntUint64(v.(int))
				}
				if v, ok := dstInfoMap["user"]; ok && v.(string) != "" {
					dBInfo.User = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["password"]; ok && v.(string) != "" {
					dBInfo.Password = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["cvm_instance_id"]; ok && v.(string) != "" {
					dBInfo.CvmInstanceId = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["uniq_vpn_gw_id"]; ok && v.(string) != "" {
					dBInfo.UniqVpnGwId = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["uniq_dcg_id"]; ok && v.(string) != "" {
					dBInfo.UniqDcgId = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["instance_id"]; ok && v.(string) != "" {
					dBInfo.InstanceId = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["ccn_gw_id"]; ok && v.(string) != "" {
					dBInfo.CcnGwId = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["vpc_id"]; ok && v.(string) != "" {
					dBInfo.VpcId = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["subnet_id"]; ok && v.(string) != "" {
					dBInfo.SubnetId = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["engine_version"]; ok && v.(string) != "" {
					dBInfo.EngineVersion = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["account"]; ok && v.(string) != "" {
					dBInfo.Account = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["account_role"]; ok && v.(string) != "" {
					dBInfo.AccountRole = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["account_mode"]; ok && v.(string) != "" {
					dBInfo.AccountMode = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["tmp_secret_id"]; ok && v.(string) != "" {
					dBInfo.TmpSecretId = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["tmp_secret_key"]; ok && v.(string) != "" {
					dBInfo.TmpSecretKey = helper.String(v.(string))
				}
				if v, ok := dstInfoMap["tmp_token"]; ok && v.(string) != "" {
					dBInfo.TmpToken = helper.String(v.(string))
				}
				dBEndpointInfo.Info = append(dBEndpointInfo.Info, &dBInfo)
			}
		}
		if v, ok := dMap["supplier"]; ok && v.(string) != "" {
			dBEndpointInfo.Supplier = helper.String(v.(string))
		}
		if v, ok := dMap["extra_attr"]; ok {
			for _, item := range v.([]interface{}) {
				extraAttrMap := item.(map[string]interface{})
				keyValuePairOption := dts.KeyValuePairOption{}
				if v, ok := extraAttrMap["key"]; ok && v.(string) != "" {
					keyValuePairOption.Key = helper.String(v.(string))
				}
				if v, ok := extraAttrMap["value"]; ok && v.(string) != "" {
					keyValuePairOption.Value = helper.String(v.(string))
				}
				dBEndpointInfo.ExtraAttr = append(dBEndpointInfo.ExtraAttr, &keyValuePairOption)
			}
		}
		configMigrationJobRequest.DstInfo = &dBEndpointInfo
	}

	if v, ok := d.GetOk("expect_run_time"); ok && v.(string) != "" {
		configMigrationJobRequest.ExpectRunTime = helper.String(v.(string))
	}

	if v, _ := d.GetOk("auto_retry_time_range_minutes"); v != nil {
		configMigrationJobRequest.AutoRetryTimeRangeMinutes = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := tcClient.UseDtsClient().ModifyMigrationJob(configMigrationJobRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, configMigrationJobRequest.GetAction(), configMigrationJobRequest.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dts migrateJob failed, reason:%+v", logId, err)
		return err
	}
	return nil
}

func handleCheckMigrate(d *schema.ResourceData, tcClient *connectivity.TencentCloudClient, logId, jobId string) error {
	checkMigrateJobRequest := dts.NewCreateMigrateCheckJobRequest()
	checkMigrateJobRequest.JobId = helper.String(jobId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := tcClient.UseDtsClient().CreateMigrateCheckJob(checkMigrateJobRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, checkMigrateJobRequest.GetAction(), checkMigrateJobRequest.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s check dts migrateJob failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

// reserve for future implementation
// func handleResumeMigrate(d *schema.ResourceData, tcClient *connectivity.TencentCloudClient, logId, jobId string) error {
// 	resumeMigrateJobRequest := dts.NewResumeMigrateJobRequest()
// 	resumeMigrateJobRequest.JobId = helper.String(jobId)
// 	service := DtsService{client: tcClient}

// 	if d.HasChange("resume_option") {
// 		if v, ok := d.GetOk("resume_option"); ok {
// 			resumeMigrateJobRequest.ResumeOption = helper.String(v.(string))
// 		}
// 	}

// 	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
// 		result, e := tcClient.UseDtsClient().ResumeMigrateJob(resumeMigrateJobRequest)
// 		if e != nil {
// 			return tccommon.RetryError(e)
// 		} else {
// 			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, resumeMigrateJobRequest.GetAction(), resumeMigrateJobRequest.ToJsonString(), result.ToJsonString())
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		log.Printf("[CRITAL]%s resume dts migrateJob failed, reason:%+v", logId, err)
// 		return err
// 	}

// 	conf := tccommon.BuildStateChangeConf([]string{}, []string{"readyComplete", "success", "failed"}, 3*tccommon.ReadRetryTimeout, time.Second, service.DtsMigrateJobStateRefreshFunc(jobId, []string{}))
// 	if _, e := conf.WaitForState(); e != nil {
// 		return e
// 	}

// 	return nil
// }

// func handleCompleteMigrate(d *schema.ResourceData, tcClient *connectivity.TencentCloudClient, logId, jobId string) error {
// 	completeMigrateJobRequest := dts.NewCompleteMigrateJobRequest()
// 	completeMigrateJobRequest.JobId = helper.String(jobId)
// 	service := DtsService{client: tcClient}

// 	if d.HasChange("complete_mode") {
// 		if v, ok := d.GetOk("complete_mode"); ok {
// 			completeMigrateJobRequest.CompleteMode = helper.String(v.(string))
// 		}
// 	}

// 	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
// 		result, e := tcClient.UseDtsClient().CompleteMigrateJob(completeMigrateJobRequest)
// 		if e != nil {
// 			return tccommon.RetryError(e)
// 		} else {
// 			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, completeMigrateJobRequest.GetAction(), completeMigrateJobRequest.ToJsonString(), result.ToJsonString())
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		log.Printf("[CRITAL]%s complete dts migrateJob failed, reason:%+v", logId, err)
// 		return err
// 	}

// 	conf := tccommon.BuildStateChangeConf([]string{}, []string{"success", "error", "failed"}, 3*tccommon.ReadRetryTimeout, time.Second, service.DtsMigrateJobStateRefreshFunc(jobId, []string{}))
// 	if _, e := conf.WaitForState(); e != nil {
// 		return e
// 	}

// 	return nil
// }

// func handleCompareMigrate(d *schema.ResourceData, tcClient *connectivity.TencentCloudClient, logId, jobId string) error {
// 	startCompareRequest := dts.NewStartCompareRequest()
// 	startCompareRequest.JobId = helper.String(jobId)

// 	if d.HasChange("compare_task_id") {
// 		if v, ok := d.GetOk("compare_task_id"); ok {
// 			startCompareRequest.CompareTaskId = helper.String(v.(string))
// 		}
// 	}

// 	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
// 		result, e := tcClient.UseDtsClient().StartCompare(startCompareRequest)
// 		if e != nil {
// 			return tccommon.RetryError(e)
// 		} else {
// 			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, startCompareRequest.GetAction(), startCompareRequest.ToJsonString(), result.ToJsonString())
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		log.Printf("[CRITAL]%s compare dts migrate job failed, reason:%+v", logId, err)
// 		return err
// 	}

// 	return nil
// }

func resourceTencentCloudDtsMigrateJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dts_migrate_job.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
