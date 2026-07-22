package wedata

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataIntegrationOfflineTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataIntegrationOfflineTaskCreate,
		Read:   resourceTencentCloudWedataIntegrationOfflineTaskRead,
		Update: resourceTencentCloudWedataIntegrationOfflineTaskUpdate,
		Delete: resourceTencentCloudWedataIntegrationOfflineTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			// OfflineTask
			"project_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "项目 ID",
			},
			"cycle_step": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Interval 时间 的 scheduling， 最小 值: 1。",
			},
			"delay_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "执行时间，单位 是 minutes，仅 可用 对于 day/week/month/year scheduling. For 示例，daily scheduling 是 executed once every day 在 02:00，和 delayTime 是 120 minutes。",
			},
			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Effective 结束时间， 格式 是 yyyy-MM-dd HH:mm:ss。",
			},
			"notes": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "描述 信息。",
			},
			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Effective 开始时间， 格式 是 yyyy-MM-dd HH:mm:ss。",
			},
			"task_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "任务 名称",
			},
			"task_action": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Scheduling 配置: flexible 周期 配置，仅 可用 对于 hourly/weekly/monthly/yearly scheduling. 如果 hourly 任务 是 指定 到 run 在 0:00，3:00 和 4:00 every day，它 是 0,3,4。",
			},
			"task_mode": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "任务 display 模式，0: canvas 模式，1: form 模式",
			},
			// IntegrationTask
			"task_info": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "任务 Information。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sync_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Synchronization 类型: 1. Whole 数据库 synchronization，2. Single 表 synchronization。",
						},
						"workflow_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "工作流 ID 到 其中 任务 belongs。",
						},
						"schedule_task_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "任务 scheduling ID (作业 ID such 作为 oceanus 或 us)。",
						},
						"task_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Inlong 任务 组 ID",
						},
						"creator_uin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "创建者 用户 ID。",
						},
						"operator_uin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "操作者 用户 ID。",
						},
						"owner_uin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "所有者 用户 ID。",
						},
						"app_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "用户 App ID。",
						},
						"status": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "任务 状态 1. Not started | 任务 initialization，2. 任务 starting，3. Running，4. Paused，5. 任务 stopping，6. Stopped，7. Execution failed，8. 删除，9. Locked，404. unknown 状态",
						},
						//"nodes": {
						//	Type:        schema.TypeList,
						//	Optional:    true,
						//	Description: "任务 Node Information.",
						//	Elem: &schema.Resource{
						//		Schema: map[string]*schema.Schema{
						//			"id": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Node ID.",
						//			},
						//			"task_id": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "任务 ID 到 其中 节点 belongs.",
						//			},
						//			"name": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Node Name.",
						//			},
						//			"node_type": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Node 类型: INPUT,OUTPUT,JOIN,FILTER,TRANSFORM.",
						//			},
						//			"data_source_type": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Data source 类型: MYSQL, POSTGRE, ORACLE, SQLSERVER, FTP, HIVE, HDFS, ICEBERG, KAFKA, HBASE, SPARK, TBASE, DB2, DM, GAUSSDB, GBASE, IMPALA, ES, S3_DATAINSIGHT, GREENPLUM, PHOENIX, SAP_HANA, SFTP, OCEANBASE, CLICKHOUSE, KUDU, VERTICA, REDIS, COS, DLC, DORIS, CKAFKA, DTS_KAFKA, S3, CDW, TDSQLC, TDSQL, MONGODB, SYBASE, REST_API, StarRocks, TCHOUSE_X.",
						//			},
						//			"description": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Node Description.",
						//			},
						//			"datasource_id": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Datasource ID.",
						//			},
						//			"config": {
						//				Type:        schema.TypeList,
						//				Optional:    true,
						//				Description: "Node 配置 信息.",
						//				Elem: &schema.Resource{
						//					Schema: map[string]*schema.Schema{
						//						"name": {
						//							Type:        schema.TypeString,
						//							Optional:    true,
						//							Description: "Configuration 名称.",
						//						},
						//						"value": {
						//							Type:        schema.TypeString,
						//							Optional:    true,
						//							Description: "Configuration 值.",
						//						},
						//					},
						//				},
						//			},
						//			"ext_config": {
						//				Type:        schema.TypeList,
						//				Optional:    true,
						//				Description: "Node extension 配置 信息.",
						//				Elem: &schema.Resource{
						//					Schema: map[string]*schema.Schema{
						//						"name": {
						//							Type:        schema.TypeString,
						//							Optional:    true,
						//							Description: "Configuration 名称.",
						//						},
						//						"value": {
						//							Type:        schema.TypeString,
						//							Optional:    true,
						//							Description: "Configuration 值.",
						//						},
						//					},
						//				},
						//			},
						//			"schema": {
						//				Type:        schema.TypeList,
						//				Optional:    true,
						//				Description: "Schema 信息.",
						//				Elem: &schema.Resource{
						//					Schema: map[string]*schema.Schema{
						//						"id": {
						//							Type:        schema.TypeString,
						//							Required:    true,
						//							Description: "Schema ID.",
						//						},
						//						"name": {
						//							Type:        schema.TypeString,
						//							Required:    true,
						//							Description: "Schema 名称.",
						//						},
						//						"type": {
						//							Type:        schema.TypeString,
						//							Required:    true,
						//							Description: "Schema 类型.",
						//						},
						//						"value": {
						//							Type:        schema.TypeString,
						//							Optional:    true,
						//							Description: "Schema 值.",
						//						},
						//						"properties": {
						//							Type:        schema.TypeList,
						//							Optional:    true,
						//							Description: "Schema extended attributes.",
						//							Elem: &schema.Resource{
						//								Schema: map[string]*schema.Schema{
						//									"name": {
						//										Type:        schema.TypeString,
						//										Optional:    true,
						//										Description: "Attributes 名称.",
						//									},
						//									"value": {
						//										Type:        schema.TypeString,
						//										Optional:    true,
						//										Description: "Attributes 值.",
						//									},
						//								},
						//							},
						//						},
						//						"alias": {
						//							Type:        schema.TypeString,
						//							Optional:    true,
						//							Description: "Schema alias.",
						//						},
						//						"comment": {
						//							Type:        schema.TypeString,
						//							Optional:    true,
						//							Description: "Schema comment.",
						//						},
						//					},
						//				},
						//			},
						//			"node_mapping": {
						//				Type:        schema.TypeList,
						//				MaxItems:    1,
						//				Optional:    true,
						//				Description: "Node mapping.",
						//				Elem: &schema.Resource{
						//					Schema: map[string]*schema.Schema{
						//						"source_id": {
						//							Type:        schema.TypeString,
						//							Optional:    true,
						//							Description: "Source 节点 ID.",
						//						},
						//						"sink_id": {
						//							Type:        schema.TypeString,
						//							Optional:    true,
						//							Description: "Sink 节点 ID.",
						//						},
						//						"source_schema": {
						//							Type:        schema.TypeList,
						//							Optional:    true,
						//							Description: "Source 节点 schema 信息.",
						//							Elem: &schema.Resource{
						//								Schema: map[string]*schema.Schema{
						//									"id": {
						//										Type:        schema.TypeString,
						//										Required:    true,
						//										Description: "Schema ID.",
						//									},
						//									"name": {
						//										Type:        schema.TypeString,
						//										Required:    true,
						//										Description: "Schema 名称.",
						//									},
						//									"type": {
						//										Type:        schema.TypeString,
						//										Required:    true,
						//										Description: "Schema 类型.",
						//									},
						//									"value": {
						//										Type:        schema.TypeString,
						//										Optional:    true,
						//										Description: "Schema 值.",
						//									},
						//									"properties": {
						//										Type:        schema.TypeList,
						//										Optional:    true,
						//										Description: "Schema extended attributes.",
						//										Elem: &schema.Resource{
						//											Schema: map[string]*schema.Schema{
						//												"name": {
						//													Type:        schema.TypeString,
						//													Optional:    true,
						//													Description: "Attributes 名称.",
						//												},
						//												"value": {
						//													Type:        schema.TypeString,
						//													Optional:    true,
						//													Description: "Attributes 值.",
						//												},
						//											},
						//										},
						//									},
						//									"alias": {
						//										Type:        schema.TypeString,
						//										Optional:    true,
						//										Description: "Schema alias.",
						//									},
						//									"comment": {
						//										Type:        schema.TypeString,
						//										Optional:    true,
						//										Description: "Schema comment.",
						//									},
						//								},
						//							},
						//						},
						//						"schema_mappings": {
						//							Type:        schema.TypeList,
						//							Optional:    true,
						//							Description: "Schema mapping 信息.",
						//							Elem: &schema.Resource{
						//								Schema: map[string]*schema.Schema{
						//									"source_schema_id": {
						//										Type:        schema.TypeString,
						//										Required:    true,
						//										Description: "Schema ID 从 source 节点.",
						//									},
						//									"sink_schema_id": {
						//										Type:        schema.TypeString,
						//										Required:    true,
						//										Description: "Schema ID 从 sink 节点.",
						//									},
						//								},
						//							},
						//						},
						//						"ext_config": {
						//							Type:        schema.TypeList,
						//							Optional:    true,
						//							Description: "Node extension 配置 信息.",
						//							Elem: &schema.Resource{
						//								Schema: map[string]*schema.Schema{
						//									"name": {
						//										Type:        schema.TypeString,
						//										Optional:    true,
						//										Description: "Configuration 名称.",
						//									},
						//									"value": {
						//										Type:        schema.TypeString,
						//										Optional:    true,
						//										Description: "Configuration 值.",
						//									},
						//								},
						//							},
						//						},
						//					},
						//				},
						//			},
						//			"app_id": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "User App ID.",
						//			},
						//			"project_id": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Project ID.",
						//			},
						//			"creator_uin": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Creator User ID.",
						//			},
						//			"operator_uin": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Operator User ID.",
						//			},
						//			"owner_uin": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Owner User ID.",
						//			},
						//			"create_time": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Create 时间.",
						//			},
						//			"update_time": {
						//				Type:        schema.TypeString,
						//				Optional:    true,
						//				Description: "Update 时间.",
						//			},
						//		},
						//	},
						//},
						"executor_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Executor 资源 ID。",
						},
						"config": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "任务 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Configuration 名称",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Configuration 值",
									},
								},
							},
						},
						"ext_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Node extension 配置 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Configuration 名称",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Configuration 值",
									},
								},
							},
						},
						"execute_context": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Execute context。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Configuration 名称",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Configuration 值",
									},
								},
							},
						},
						"mappings": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Node mapping。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "来源 节点 ID",
									},
									"sink_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Sink 节点 ID",
									},
									"source_schema": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "来源 节点 schema 信息。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Schema ID",
												},
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Schema 名称",
												},
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Schema 类型",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Schema 值",
												},
												"properties": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Schema extended attributes。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Attributes 名称",
															},
															"value": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Attributes 值",
															},
														},
													},
												},
												"alias": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Schema 别名",
												},
												"comment": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Schema 注释",
												},
											},
										},
									},
									"schema_mappings": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Schema mapping 信息。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"source_schema_id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Schema ID 从 来源 节点。",
												},
												"sink_schema_id": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Schema ID 从 sink 节点。",
												},
											},
										},
									},
									"ext_config": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Node extension 配置 信息。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Configuration 名称",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Configuration 值",
												},
											},
										},
									},
								},
							},
						},
						"task_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "任务 display 模式，0: canvas 模式，1: form 模式",
						},
						"incharge": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Incharge 用户",
						},
						"offline_task_add_entity": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Offline 任务 scheduling 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									//"workflow_name": {
									//	Type:        schema.TypeString,
									//	Optional:    true,
									//	Description: "名称 的 工作流 到 其中 任务 belongs.",
									//},
									//"dependency_workflow": {
									//	Type:        schema.TypeString,
									//	Optional:    true,
									//	Description: "Whether 到 support 工作流 dependencies: yes / 无, 默认值 值: 无.",
									//},
									//"start_time": {
									//	Type:        schema.TypeString,
									//	Optional:    true,
									//	Description: "Effective start 时间, 格式 是 yyyy-MM-dd HH:mm:ss.",
									//},
									//"end_time": {
									//	Type:        schema.TypeString,
									//	Optional:    true,
									//	Description: "Effective end 时间, 格式 是 yyyy-MM-dd HH:mm:ss.",
									//},
									"cycle_type": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Scheduling 类型，0: crontab 类型，1: minutes，2: hours，3: days，4: weeks，5: months，6: 一个-时间，7: 用户-driven，10: elastic 周期 (week)，11: elastic 周期 (month)，12: year，13: instant 触发器。",
									},
									//"cycle_step": {
									//	Type:        schema.TypeInt,
									//	Optional:    true,
									//	Description: "Interval 时间 的 scheduling, 最小 值: 1.",
									//},
									//"delay_time": {
									//	Type:        schema.TypeInt,
									//	Optional:    true,
									//	Description: "Execution 时间, 单位 是 minutes, 仅 可用 对于 day/week/month/year scheduling. For 示例, daily scheduling 是 executed once every day 在 02:00, 和 delayTime 是 120 minutes.",
									//},
									"crontab_expression": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Crontab expression。",
									},
									"retry_wait": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Retry waiting 时间，单位 是 minutes。",
									},
									"retriable": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "是否retry。",
									},
									"try_limit": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "数量 retries。",
									},
									//"run_priority": {
									//	Type:        schema.TypeInt,
									//	Optional:    true,
									//	Description: "任务 running 优先级.",
									//},
									//"product_name": {
									//	Type:        schema.TypeString,
									//	Optional:    true,
									//	Description: "Product 名称: DATA_INTEGRATION.",
									//},
									"self_depend": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Self-dependent 规则，1: Ordered serial 一个 在 时间，queued execution，2: Unordered serial 一个 在 时间，不 queued execution，3: Parallel，多个 在 once。",
									},
									//"task_action": {
									//	Type:        schema.TypeString,
									//	Optional:    true,
									//	Description: "Flexible cycle 配置, 如果 它 是 weekly 任务: 1 是 Sunday, 2 是 Monday, 3 是 Tuesday, 和 so 在. 如果 它 是 monthly 任务: 1, 表示 1st 和 3rd; L 表示 end 的 month.",
									//},
									"execution_end_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Scheduling execution 结束时间。",
									},
									"execution_start_time": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Scheduling execution 开始时间。",
									},
									//"task_auto_submit": {
									//	Type:        schema.TypeBool,
									//	Optional:    true,
									//	Description: "Whether 到 automatically submit.",
									//},
									//"instance_init_strategy": {
									//	Type:        schema.TypeString,
									//	Optional:    true,
									//	Description: "实例 initialization strategy.",
									//},
								},
							},
						},
						"executor_group_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Executor 组名称",
						},
						"in_long_manager_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "InLong manager URL",
						},
						"in_long_stream_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "InLong 流 ID。",
						},
						"in_long_manager_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "InLong manager 版本",
						},
						"data_proxy_url": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "Data proxy URL",
						},
						"submit": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "是否task 版本 has been submitted 对于 operation 和 maintenance。",
						},
						"input_datasource_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Input datasource 类型",
						},
						"output_datasource_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Output datasource 类型",
						},
						"num_records_in": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "数量 reads。",
						},
						"num_records_out": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "数量 writes。",
						},
						"reader_delay": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: "Read 延迟。",
						},
						"num_restarts": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Times 的 restarts。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "创建时间。",
						},
						"update_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "更新时间。",
						},
						"last_run_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "last 时间 任务 是 run。",
						},
						"stop_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "时间 任务 是 stopped。",
						},
						"has_version": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "是否task been submitted。",
						},
						"locked": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "是否task been locked。",
						},
						"locker": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "用户 locked 任务。",
						},
						"running_cu": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: "amount 的 resources consumed 通过 real-时间 任务。",
						},
						"task_alarm_regular_list": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "任务 告警 regular。",
						},
						"switch_resource": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Resource tiering 状态，0: 在 progress，1: successful，2: failed。",
						},
						"read_phase": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Reading stage，0: full amount，1: partial full amount，2: all incremental。",
						},
						"instance_version": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "实例 版本",
						},
					},
				},
			},
			// computed
			"task_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "任务 ID",
			},
		},
	}
}

func resourceTencentCloudWedataIntegrationOfflineTaskCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_integration_offline_task.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                        = tccommon.GetLogId(tccommon.ContextNil)
		request                      = wedata.NewCreateOfflineTaskRequest()
		response                     = wedata.NewCreateOfflineTaskResponse()
		modifyIntegrationTaskRequest = wedata.NewModifyIntegrationTaskRequest()
		projectId                    string
		taskId                       string
		taskName                     string
		notes                        string
		taskAction                   string
		startTime                    string
		endTime                      string
		cycleStep                    int
		delayTime                    int
	)

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOkExists("cycle_step"); ok {
		request.CycleStep = helper.IntInt64(v.(int))
		cycleStep = v.(int)
	}

	if v, ok := d.GetOkExists("delay_time"); ok {
		request.DelayTime = helper.IntInt64(v.(int))
		delayTime = v.(int)
	}

	if v, ok := d.GetOk("end_time"); ok {
		request.EndTime = helper.String(v.(string))
		endTime = v.(string)
	}

	if v, ok := d.GetOk("notes"); ok {
		request.Notes = helper.String(v.(string))
		notes = v.(string)
	}

	if v, ok := d.GetOk("start_time"); ok {
		request.StartTime = helper.String(v.(string))
		startTime = v.(string)
	}

	if v, ok := d.GetOk("task_name"); ok {
		request.TaskName = helper.String(v.(string))
		taskName = v.(string)
	}

	request.TypeId = helper.IntInt64(27)

	if v, ok := d.GetOk("task_action"); ok {
		request.TaskAction = helper.String(v.(string))
		taskAction = v.(string)
	}

	if v, ok := d.GetOk("task_mode"); ok {
		request.TaskMode = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataClient().CreateOfflineTask(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("wedata integrationOfflineTask not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create wedata integrationOfflineTask failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(strings.Join([]string{projectId, taskId}, tccommon.FILED_SP))

	// Create IntegrationTask

	if dMap, ok := helper.InterfacesHeadMap(d, "task_info"); ok {
		integrationTaskInfo := wedata.IntegrationTaskInfo{}
		integrationTaskInfo.ProjectId = &projectId
		integrationTaskInfo.TaskId = &taskId
		integrationTaskInfo.TaskName = &taskName
		integrationTaskInfo.Description = &notes
		integrationTaskInfo.TaskType = helper.IntInt64(202)

		if v, ok := dMap["sync_type"]; ok {
			integrationTaskInfo.SyncType = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["workflow_id"]; ok {
			integrationTaskInfo.WorkflowId = helper.String(v.(string))
		}

		if v, ok := dMap["schedule_task_id"]; ok {
			integrationTaskInfo.ScheduleTaskId = helper.String(v.(string))
		}

		if v, ok := dMap["task_group_id"]; ok {
			integrationTaskInfo.TaskGroupId = helper.String(v.(string))
		}

		if v, ok := dMap["creator_uin"]; ok {
			integrationTaskInfo.CreatorUin = helper.String(v.(string))
		}

		if v, ok := dMap["operator_uin"]; ok {
			integrationTaskInfo.OperatorUin = helper.String(v.(string))
		}

		if v, ok := dMap["owner_uin"]; ok {
			integrationTaskInfo.OwnerUin = helper.String(v.(string))
		}

		if v, ok := dMap["app_id"]; ok {
			integrationTaskInfo.AppId = helper.String(v.(string))
		}

		if v, ok := dMap["status"]; ok {
			integrationTaskInfo.Status = helper.IntInt64(v.(int))
		}

		//if v, ok := dMap["nodes"]; ok {
		//	for _, item := range v.([]interface{}) {
		//		nodesMap := item.(map[string]interface{})
		//		integrationNodeInfo := wedata.IntegrationNodeInfo{}
		//		if v, ok := nodesMap["id"]; ok {
		//			integrationNodeInfo.Id = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["task_id"]; ok {
		//			integrationNodeInfo.TaskId = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["name"]; ok {
		//			integrationNodeInfo.Name = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["node_type"]; ok {
		//			integrationNodeInfo.NodeType = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["data_source_type"]; ok {
		//			integrationNodeInfo.DataSourceType = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["description"]; ok {
		//			integrationNodeInfo.Description = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["datasource_id"]; ok {
		//			integrationNodeInfo.DatasourceId = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["config"]; ok {
		//			for _, item := range v.([]interface{}) {
		//				configMap := item.(map[string]interface{})
		//				recordField := wedata.RecordField{}
		//				if v, ok := configMap["name"]; ok {
		//					recordField.Name = helper.String(v.(string))
		//				}
		//
		//				if v, ok := configMap["value"]; ok {
		//					recordField.Value = helper.String(v.(string))
		//				}
		//
		//				integrationNodeInfo.Config = append(integrationNodeInfo.Config, &recordField)
		//			}
		//		}
		//
		//		if v, ok := nodesMap["ext_config"]; ok {
		//			for _, item := range v.([]interface{}) {
		//				extConfigMap := item.(map[string]interface{})
		//				recordField := wedata.RecordField{}
		//				if v, ok := extConfigMap["name"]; ok {
		//					recordField.Name = helper.String(v.(string))
		//				}
		//
		//				if v, ok := extConfigMap["value"]; ok {
		//					recordField.Value = helper.String(v.(string))
		//				}
		//
		//				integrationNodeInfo.ExtConfig = append(integrationNodeInfo.ExtConfig, &recordField)
		//			}
		//		}
		//
		//		if v, ok := nodesMap["schema"]; ok {
		//			for _, item := range v.([]interface{}) {
		//				schemaMap := item.(map[string]interface{})
		//				integrationNodeSchema := wedata.IntegrationNodeSchema{}
		//				if v, ok := schemaMap["id"]; ok {
		//					integrationNodeSchema.Id = helper.String(v.(string))
		//				}
		//
		//				if v, ok := schemaMap["name"]; ok {
		//					integrationNodeSchema.Name = helper.String(v.(string))
		//				}
		//
		//				if v, ok := schemaMap["type"]; ok {
		//					integrationNodeSchema.Type = helper.String(v.(string))
		//				}
		//
		//				if v, ok := schemaMap["value"]; ok {
		//					integrationNodeSchema.Value = helper.String(v.(string))
		//				}
		//
		//				if v, ok := schemaMap["properties"]; ok {
		//					for _, item := range v.([]interface{}) {
		//						propertiesMap := item.(map[string]interface{})
		//						recordField := wedata.RecordField{}
		//						if v, ok := propertiesMap["name"]; ok {
		//							recordField.Name = helper.String(v.(string))
		//						}
		//
		//						if v, ok := propertiesMap["value"]; ok {
		//							recordField.Value = helper.String(v.(string))
		//						}
		//
		//						integrationNodeSchema.Properties = append(integrationNodeSchema.Properties, &recordField)
		//					}
		//				}
		//
		//				if v, ok := schemaMap["alias"]; ok {
		//					integrationNodeSchema.Alias = helper.String(v.(string))
		//				}
		//
		//				if v, ok := schemaMap["comment"]; ok {
		//					integrationNodeSchema.Comment = helper.String(v.(string))
		//				}
		//
		//				integrationNodeInfo.Schema = append(integrationNodeInfo.Schema, &integrationNodeSchema)
		//			}
		//		}
		//		if nodeMappingMap, ok := helper.InterfaceToMap(nodesMap, "node_mapping"); ok {
		//			integrationNodeMapping := wedata.IntegrationNodeMapping{}
		//			if v, ok := nodeMappingMap["source_id"]; ok {
		//				integrationNodeMapping.SourceId = helper.String(v.(string))
		//			}
		//
		//			if v, ok := nodeMappingMap["sink_id"]; ok {
		//				integrationNodeMapping.SinkId = helper.String(v.(string))
		//			}
		//
		//			if v, ok := nodeMappingMap["source_schema"]; ok {
		//				for _, item := range v.([]interface{}) {
		//					sourceSchemaMap := item.(map[string]interface{})
		//					integrationNodeSchema := wedata.IntegrationNodeSchema{}
		//					if v, ok := sourceSchemaMap["id"]; ok {
		//						integrationNodeSchema.Id = helper.String(v.(string))
		//					}
		//
		//					if v, ok := sourceSchemaMap["name"]; ok {
		//						integrationNodeSchema.Name = helper.String(v.(string))
		//					}
		//
		//					if v, ok := sourceSchemaMap["type"]; ok {
		//						integrationNodeSchema.Type = helper.String(v.(string))
		//					}
		//
		//					if v, ok := sourceSchemaMap["value"]; ok {
		//						integrationNodeSchema.Value = helper.String(v.(string))
		//					}
		//
		//					if v, ok := sourceSchemaMap["properties"]; ok {
		//						for _, item := range v.([]interface{}) {
		//							propertiesMap := item.(map[string]interface{})
		//							recordField := wedata.RecordField{}
		//							if v, ok := propertiesMap["name"]; ok {
		//								recordField.Name = helper.String(v.(string))
		//							}
		//
		//							if v, ok := propertiesMap["value"]; ok {
		//								recordField.Value = helper.String(v.(string))
		//							}
		//
		//							integrationNodeSchema.Properties = append(integrationNodeSchema.Properties, &recordField)
		//						}
		//					}
		//
		//					if v, ok := sourceSchemaMap["alias"]; ok {
		//						integrationNodeSchema.Alias = helper.String(v.(string))
		//					}
		//
		//					if v, ok := sourceSchemaMap["comment"]; ok {
		//						integrationNodeSchema.Comment = helper.String(v.(string))
		//					}
		//
		//					integrationNodeMapping.SourceSchema = append(integrationNodeMapping.SourceSchema, &integrationNodeSchema)
		//				}
		//			}
		//
		//			if v, ok := nodeMappingMap["schema_mappings"]; ok {
		//				for _, item := range v.([]interface{}) {
		//					schemaMappingsMap := item.(map[string]interface{})
		//					integrationNodeSchemaMapping := wedata.IntegrationNodeSchemaMapping{}
		//					if v, ok := schemaMappingsMap["source_schema_id"]; ok {
		//						integrationNodeSchemaMapping.SourceSchemaId = helper.String(v.(string))
		//					}
		//
		//					if v, ok := schemaMappingsMap["sink_schema_id"]; ok {
		//						integrationNodeSchemaMapping.SinkSchemaId = helper.String(v.(string))
		//					}
		//
		//					integrationNodeMapping.SchemaMappings = append(integrationNodeMapping.SchemaMappings, &integrationNodeSchemaMapping)
		//				}
		//			}
		//
		//			if v, ok := nodeMappingMap["ext_config"]; ok {
		//				for _, item := range v.([]interface{}) {
		//					extConfigMap := item.(map[string]interface{})
		//					recordField := wedata.RecordField{}
		//					if v, ok := extConfigMap["name"]; ok {
		//						recordField.Name = helper.String(v.(string))
		//					}
		//
		//					if v, ok := extConfigMap["value"]; ok {
		//						recordField.Value = helper.String(v.(string))
		//					}
		//
		//					integrationNodeMapping.ExtConfig = append(integrationNodeMapping.ExtConfig, &recordField)
		//				}
		//			}
		//
		//			integrationNodeInfo.NodeMapping = &integrationNodeMapping
		//		}
		//
		//		if v, ok := nodesMap["app_id"]; ok {
		//			integrationNodeInfo.AppId = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["project_id"]; ok {
		//			integrationNodeInfo.ProjectId = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["creator_uin"]; ok {
		//			integrationNodeInfo.CreatorUin = helper.String(v.(string))
		//		}
		//		if v, ok := nodesMap["operator_uin"]; ok {
		//			integrationNodeInfo.OperatorUin = helper.String(v.(string))
		//		}
		//		if v, ok := nodesMap["owner_uin"]; ok {
		//			integrationNodeInfo.OwnerUin = helper.String(v.(string))
		//		}
		//		if v, ok := nodesMap["create_time"]; ok {
		//			integrationNodeInfo.CreateTime = helper.String(v.(string))
		//		}
		//		if v, ok := nodesMap["update_time"]; ok {
		//			integrationNodeInfo.UpdateTime = helper.String(v.(string))
		//		}
		//		integrationTaskInfo.Nodes = append(integrationTaskInfo.Nodes, &integrationNodeInfo)
		//	}
		//}

		if v, ok := dMap["executor_id"]; ok {
			integrationTaskInfo.ExecutorId = helper.String(v.(string))
		}

		if v, ok := dMap["config"]; ok {
			for _, item := range v.([]interface{}) {
				configMap := item.(map[string]interface{})
				recordField := wedata.RecordField{}
				if v, ok := configMap["name"]; ok {
					recordField.Name = helper.String(v.(string))
				}

				if v, ok := configMap["value"]; ok {
					recordField.Value = helper.String(v.(string))
				}

				integrationTaskInfo.Config = append(integrationTaskInfo.Config, &recordField)
			}
		}

		if v, ok := dMap["ext_config"]; ok {
			for _, item := range v.([]interface{}) {
				extConfigMap := item.(map[string]interface{})
				recordField := wedata.RecordField{}
				if v, ok := extConfigMap["name"]; ok {
					recordField.Name = helper.String(v.(string))
				}

				if v, ok := extConfigMap["value"]; ok {
					recordField.Value = helper.String(v.(string))
				}

				integrationTaskInfo.ExtConfig = append(integrationTaskInfo.ExtConfig, &recordField)
			}
		}

		if v, ok := dMap["execute_context"]; ok {
			for _, item := range v.([]interface{}) {
				executeContextMap := item.(map[string]interface{})
				recordField := wedata.RecordField{}
				if v, ok := executeContextMap["name"]; ok {
					recordField.Name = helper.String(v.(string))
				}

				if v, ok := executeContextMap["value"]; ok {
					recordField.Value = helper.String(v.(string))
				}

				integrationTaskInfo.ExecuteContext = append(integrationTaskInfo.ExecuteContext, &recordField)
			}
		}

		if v, ok := dMap["mappings"]; ok {
			for _, item := range v.([]interface{}) {
				mappingsMap := item.(map[string]interface{})
				integrationNodeMapping := wedata.IntegrationNodeMapping{}
				if v, ok := mappingsMap["source_id"]; ok {
					integrationNodeMapping.SourceId = helper.String(v.(string))
				}

				if v, ok := mappingsMap["sink_id"]; ok {
					integrationNodeMapping.SinkId = helper.String(v.(string))
				}

				if v, ok := mappingsMap["source_schema"]; ok {
					for _, item := range v.([]interface{}) {
						sourceSchemaMap := item.(map[string]interface{})
						integrationNodeSchema := wedata.IntegrationNodeSchema{}
						if v, ok := sourceSchemaMap["id"]; ok {
							integrationNodeSchema.Id = helper.String(v.(string))
						}

						if v, ok := sourceSchemaMap["name"]; ok {
							integrationNodeSchema.Name = helper.String(v.(string))
						}

						if v, ok := sourceSchemaMap["type"]; ok {
							integrationNodeSchema.Type = helper.String(v.(string))
						}

						if v, ok := sourceSchemaMap["value"]; ok {
							integrationNodeSchema.Value = helper.String(v.(string))
						}

						if v, ok := sourceSchemaMap["properties"]; ok {
							for _, item := range v.([]interface{}) {
								propertiesMap := item.(map[string]interface{})
								recordField := wedata.RecordField{}
								if v, ok := propertiesMap["name"]; ok {
									recordField.Name = helper.String(v.(string))
								}

								if v, ok := propertiesMap["value"]; ok {
									recordField.Value = helper.String(v.(string))
								}

								integrationNodeSchema.Properties = append(integrationNodeSchema.Properties, &recordField)
							}
						}

						if v, ok := sourceSchemaMap["alias"]; ok {
							integrationNodeSchema.Alias = helper.String(v.(string))
						}

						if v, ok := sourceSchemaMap["comment"]; ok {
							integrationNodeSchema.Comment = helper.String(v.(string))
						}

						integrationNodeMapping.SourceSchema = append(integrationNodeMapping.SourceSchema, &integrationNodeSchema)
					}
				}

				if v, ok := mappingsMap["schema_mappings"]; ok {
					for _, item := range v.([]interface{}) {
						schemaMappingsMap := item.(map[string]interface{})
						integrationNodeSchemaMapping := wedata.IntegrationNodeSchemaMapping{}
						if v, ok := schemaMappingsMap["source_schema_id"]; ok {
							integrationNodeSchemaMapping.SourceSchemaId = helper.String(v.(string))
						}

						if v, ok := schemaMappingsMap["sink_schema_id"]; ok {
							integrationNodeSchemaMapping.SinkSchemaId = helper.String(v.(string))
						}

						integrationNodeMapping.SchemaMappings = append(integrationNodeMapping.SchemaMappings, &integrationNodeSchemaMapping)
					}
				}

				if v, ok := mappingsMap["ext_config"]; ok {
					for _, item := range v.([]interface{}) {
						extConfigMap := item.(map[string]interface{})
						recordField := wedata.RecordField{}
						if v, ok := extConfigMap["name"]; ok {
							recordField.Name = helper.String(v.(string))
						}

						if v, ok := extConfigMap["value"]; ok {
							recordField.Value = helper.String(v.(string))
						}

						integrationNodeMapping.ExtConfig = append(integrationNodeMapping.ExtConfig, &recordField)
					}
				}

				integrationTaskInfo.Mappings = append(integrationTaskInfo.Mappings, &integrationNodeMapping)
			}
		}

		if v, ok := dMap["task_mode"]; ok {
			integrationTaskInfo.TaskMode = helper.String(v.(string))
		}

		if v, ok := dMap["incharge"]; ok {
			integrationTaskInfo.Incharge = helper.String(v.(string))
		}

		if offlineTaskAddEntityMap, ok := helper.InterfaceToMap(dMap, "offline_task_add_entity"); ok {
			offlineTaskAddParam := wedata.OfflineTaskAddParam{}
			//if v, ok := offlineTaskAddEntityMap["workflow_name"]; ok {
			//	offlineTaskAddParam.WorkflowName = helper.String(v.(string))
			//}
			//
			//if v, ok := offlineTaskAddEntityMap["dependency_workflow"]; ok {
			//	offlineTaskAddParam.DependencyWorkflow = helper.String(v.(string))
			//}

			offlineTaskAddParam.StartTime = &startTime
			offlineTaskAddParam.EndTime = &endTime
			offlineTaskAddParam.CycleStep = helper.IntUint64(cycleStep)
			offlineTaskAddParam.DelayTime = helper.IntUint64(delayTime)
			offlineTaskAddParam.TaskAction = &taskAction

			if v, ok := offlineTaskAddEntityMap["cycle_type"]; ok {
				offlineTaskAddParam.CycleType = helper.IntUint64(v.(int))
			}

			if v, ok := offlineTaskAddEntityMap["crontab_expression"]; ok {
				offlineTaskAddParam.CrontabExpression = helper.String(v.(string))
			}

			if v, ok := offlineTaskAddEntityMap["retry_wait"]; ok {
				offlineTaskAddParam.RetryWait = helper.IntUint64(v.(int))
			}

			if v, ok := offlineTaskAddEntityMap["retriable"]; ok {
				offlineTaskAddParam.Retriable = helper.IntUint64(v.(int))
			}

			if v, ok := offlineTaskAddEntityMap["try_limit"]; ok {
				offlineTaskAddParam.TryLimit = helper.IntUint64(v.(int))
			}

			//if v, ok := offlineTaskAddEntityMap["run_priority"]; ok {
			//	offlineTaskAddParam.RunPriority = helper.IntUint64(v.(int))
			//}
			//
			//if v, ok := offlineTaskAddEntityMap["product_name"]; ok {
			//	offlineTaskAddParam.ProductName = helper.String(v.(string))
			//}

			if v, ok := offlineTaskAddEntityMap["self_depend"]; ok {
				offlineTaskAddParam.SelfDepend = helper.IntUint64(v.(int))
			}

			if v, ok := offlineTaskAddEntityMap["execution_end_time"]; ok {
				offlineTaskAddParam.ExecutionEndTime = helper.String(v.(string))
			}

			if v, ok := offlineTaskAddEntityMap["execution_start_time"]; ok {
				offlineTaskAddParam.ExecutionStartTime = helper.String(v.(string))
			}

			//if v, ok := offlineTaskAddEntityMap["task_auto_submit"]; ok {
			//	offlineTaskAddParam.TaskAutoSubmit = helper.Bool(v.(bool))
			//}
			//
			//if v, ok := offlineTaskAddEntityMap["instance_init_strategy"]; ok {
			//	offlineTaskAddParam.InstanceInitStrategy = helper.String(v.(string))
			//}

			integrationTaskInfo.OfflineTaskAddEntity = &offlineTaskAddParam
		}

		if v, ok := dMap["executor_group_name"]; ok {
			integrationTaskInfo.ExecutorGroupName = helper.String(v.(string))
		}

		if v, ok := dMap["in_long_manager_url"]; ok {
			integrationTaskInfo.InLongManagerUrl = helper.String(v.(string))
		}

		if v, ok := dMap["in_long_stream_id"]; ok {
			integrationTaskInfo.InLongStreamId = helper.String(v.(string))
		}

		if v, ok := dMap["in_long_manager_version"]; ok {
			integrationTaskInfo.InLongManagerVersion = helper.String(v.(string))
		}

		if v, ok := dMap["data_proxy_url"]; ok {
			dataProxyUrlSet := v.(*schema.Set).List()
			for i := range dataProxyUrlSet {
				dataProxyUrl := dataProxyUrlSet[i].(string)
				integrationTaskInfo.DataProxyUrl = append(integrationTaskInfo.DataProxyUrl, &dataProxyUrl)
			}
		}

		if v, ok := dMap["submit"]; ok {
			integrationTaskInfo.Submit = helper.Bool(v.(bool))
		}

		if v, ok := dMap["input_datasource_type"]; ok {
			integrationTaskInfo.InputDatasourceType = helper.String(v.(string))
		}

		if v, ok := dMap["output_datasource_type"]; ok {
			integrationTaskInfo.OutputDatasourceType = helper.String(v.(string))
		}

		if v, ok := dMap["num_records_in"]; ok {
			integrationTaskInfo.NumRecordsIn = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["num_records_out"]; ok {
			integrationTaskInfo.NumRecordsOut = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["reader_delay"]; ok {
			integrationTaskInfo.ReaderDelay = helper.Float64(v.(float64))
		}

		if v, ok := dMap["num_restarts"]; ok {
			integrationTaskInfo.NumRestarts = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["create_time"]; ok {
			integrationTaskInfo.CreateTime = helper.String(v.(string))
		}

		if v, ok := dMap["update_time"]; ok {
			integrationTaskInfo.UpdateTime = helper.String(v.(string))
		}

		if v, ok := dMap["last_run_time"]; ok {
			integrationTaskInfo.LastRunTime = helper.String(v.(string))
		}

		if v, ok := dMap["stop_time"]; ok {
			integrationTaskInfo.StopTime = helper.String(v.(string))
		}

		if v, ok := dMap["has_version"]; ok {
			integrationTaskInfo.HasVersion = helper.Bool(v.(bool))
		}

		if v, ok := dMap["locked"]; ok {
			integrationTaskInfo.Locked = helper.Bool(v.(bool))
		}

		if v, ok := dMap["locker"]; ok {
			integrationTaskInfo.Locker = helper.String(v.(string))
		}

		if v, ok := dMap["running_cu"]; ok {
			integrationTaskInfo.RunningCu = helper.Float64(v.(float64))
		}

		if v, ok := dMap["task_alarm_regular_list"]; ok {
			taskAlarmRegularListSet := v.(*schema.Set).List()
			for i := range taskAlarmRegularListSet {
				taskAlarmRegularList := taskAlarmRegularListSet[i].(string)
				integrationTaskInfo.TaskAlarmRegularList = append(integrationTaskInfo.TaskAlarmRegularList, &taskAlarmRegularList)
			}
		}

		if v, ok := dMap["switch_resource"]; ok {
			integrationTaskInfo.SwitchResource = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["read_phase"]; ok {
			integrationTaskInfo.ReadPhase = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["instance_version"]; ok {
			integrationTaskInfo.InstanceVersion = helper.IntInt64(v.(int))
		}

		modifyIntegrationTaskRequest.TaskInfo = &integrationTaskInfo
	}

	modifyIntegrationTaskRequest.ProjectId = &projectId
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataClient().ModifyIntegrationTask(modifyIntegrationTaskRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyIntegrationTaskRequest.GetAction(), modifyIntegrationTaskRequest.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create wedata integration_real_time_task failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataIntegrationOfflineTaskRead(d, meta)
}

func resourceTencentCloudWedataIntegrationOfflineTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_integration_offline_task.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId := idSplit[0]
	taskId := idSplit[1]

	integrationOfflineTask, err := service.DescribeWedataIntegrationOfflineTaskById(ctx, projectId, taskId)
	if err != nil {
		return err
	}

	if integrationOfflineTask == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WedataIntegrationOfflineTask` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("project_id", projectId)
	_ = d.Set("task_id", taskId)

	if integrationOfflineTask.TaskInfo.OfflineTaskAddEntity.CycleStep != nil {
		_ = d.Set("cycle_step", integrationOfflineTask.TaskInfo.OfflineTaskAddEntity.CycleStep)
	}

	if integrationOfflineTask.TaskInfo.OfflineTaskAddEntity.DelayTime != nil {
		_ = d.Set("delay_time", integrationOfflineTask.TaskInfo.OfflineTaskAddEntity.DelayTime)
	}

	if integrationOfflineTask.TaskInfo.OfflineTaskAddEntity.EndTime != nil {
		_ = d.Set("end_time", integrationOfflineTask.TaskInfo.OfflineTaskAddEntity.EndTime)
	}

	if integrationOfflineTask.TaskInfo.OfflineTaskAddEntity.StartTime != nil {
		_ = d.Set("start_time", integrationOfflineTask.TaskInfo.OfflineTaskAddEntity.StartTime)
	}

	if integrationOfflineTask.TaskInfo.TaskName != nil {
		_ = d.Set("task_name", integrationOfflineTask.TaskInfo.TaskName)
	}

	if integrationOfflineTask.TaskInfo.OfflineTaskAddEntity.TaskAction != nil {
		_ = d.Set("task_action", integrationOfflineTask.TaskInfo.OfflineTaskAddEntity.TaskAction)
	}

	if integrationOfflineTask.TaskInfo.TaskMode != nil {
		_ = d.Set("task_mode", integrationOfflineTask.TaskInfo.TaskMode)
	}

	return nil
}

func resourceTencentCloudWedataIntegrationOfflineTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_integration_offline_task.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = wedata.NewModifyIntegrationTaskRequest()
		taskName   string
		notes      string
		taskAction string
		startTime  string
		endTime    string
		cycleStep  int
		delayTime  int
	)

	immutableArgs := []string{"project_id", "task_mode"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId := idSplit[0]
	taskId := idSplit[1]

	request.ProjectId = &projectId

	if v, ok := d.GetOkExists("cycle_step"); ok {
		cycleStep = v.(int)
	}

	if v, ok := d.GetOkExists("delay_time"); ok {
		delayTime = v.(int)
	}

	if v, ok := d.GetOk("end_time"); ok {
		endTime = v.(string)
	}

	if v, ok := d.GetOk("notes"); ok {
		notes = v.(string)
	}

	if v, ok := d.GetOk("start_time"); ok {
		startTime = v.(string)
	}

	if v, ok := d.GetOk("task_name"); ok {
		taskName = v.(string)
	}

	if v, ok := d.GetOk("task_action"); ok {
		taskAction = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "task_info"); ok {
		integrationTaskInfo := wedata.IntegrationTaskInfo{}
		integrationTaskInfo.ProjectId = &projectId
		integrationTaskInfo.TaskId = &taskId
		integrationTaskInfo.TaskName = &taskName
		integrationTaskInfo.Description = &notes
		integrationTaskInfo.TaskType = helper.IntInt64(202)

		if v, ok := dMap["sync_type"]; ok {
			integrationTaskInfo.SyncType = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["workflow_id"]; ok {
			integrationTaskInfo.WorkflowId = helper.String(v.(string))
		}

		if v, ok := dMap["schedule_task_id"]; ok {
			integrationTaskInfo.ScheduleTaskId = helper.String(v.(string))
		}

		if v, ok := dMap["task_group_id"]; ok {
			integrationTaskInfo.TaskGroupId = helper.String(v.(string))
		}

		if v, ok := dMap["creator_uin"]; ok {
			integrationTaskInfo.CreatorUin = helper.String(v.(string))
		}

		if v, ok := dMap["operator_uin"]; ok {
			integrationTaskInfo.OperatorUin = helper.String(v.(string))
		}

		if v, ok := dMap["owner_uin"]; ok {
			integrationTaskInfo.OwnerUin = helper.String(v.(string))
		}

		if v, ok := dMap["app_id"]; ok {
			integrationTaskInfo.AppId = helper.String(v.(string))
		}

		if v, ok := dMap["status"]; ok {
			integrationTaskInfo.Status = helper.IntInt64(v.(int))
		}

		//if v, ok := dMap["nodes"]; ok {
		//	for _, item := range v.([]interface{}) {
		//		nodesMap := item.(map[string]interface{})
		//		integrationNodeInfo := wedata.IntegrationNodeInfo{}
		//		if v, ok := nodesMap["id"]; ok {
		//			integrationNodeInfo.Id = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["task_id"]; ok {
		//			integrationNodeInfo.TaskId = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["name"]; ok {
		//			integrationNodeInfo.Name = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["node_type"]; ok {
		//			integrationNodeInfo.NodeType = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["data_source_type"]; ok {
		//			integrationNodeInfo.DataSourceType = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["description"]; ok {
		//			integrationNodeInfo.Description = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["datasource_id"]; ok {
		//			integrationNodeInfo.DatasourceId = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["config"]; ok {
		//			for _, item := range v.([]interface{}) {
		//				configMap := item.(map[string]interface{})
		//				recordField := wedata.RecordField{}
		//				if v, ok := configMap["name"]; ok {
		//					recordField.Name = helper.String(v.(string))
		//				}
		//
		//				if v, ok := configMap["value"]; ok {
		//					recordField.Value = helper.String(v.(string))
		//				}
		//
		//				integrationNodeInfo.Config = append(integrationNodeInfo.Config, &recordField)
		//			}
		//		}
		//
		//		if v, ok := nodesMap["ext_config"]; ok {
		//			for _, item := range v.([]interface{}) {
		//				extConfigMap := item.(map[string]interface{})
		//				recordField := wedata.RecordField{}
		//				if v, ok := extConfigMap["name"]; ok {
		//					recordField.Name = helper.String(v.(string))
		//				}
		//
		//				if v, ok := extConfigMap["value"]; ok {
		//					recordField.Value = helper.String(v.(string))
		//				}
		//
		//				integrationNodeInfo.ExtConfig = append(integrationNodeInfo.ExtConfig, &recordField)
		//			}
		//		}
		//
		//		if v, ok := nodesMap["schema"]; ok {
		//			for _, item := range v.([]interface{}) {
		//				schemaMap := item.(map[string]interface{})
		//				integrationNodeSchema := wedata.IntegrationNodeSchema{}
		//				if v, ok := schemaMap["id"]; ok {
		//					integrationNodeSchema.Id = helper.String(v.(string))
		//				}
		//
		//				if v, ok := schemaMap["name"]; ok {
		//					integrationNodeSchema.Name = helper.String(v.(string))
		//				}
		//
		//				if v, ok := schemaMap["type"]; ok {
		//					integrationNodeSchema.Type = helper.String(v.(string))
		//				}
		//
		//				if v, ok := schemaMap["value"]; ok {
		//					integrationNodeSchema.Value = helper.String(v.(string))
		//				}
		//
		//				if v, ok := schemaMap["properties"]; ok {
		//					for _, item := range v.([]interface{}) {
		//						propertiesMap := item.(map[string]interface{})
		//						recordField := wedata.RecordField{}
		//						if v, ok := propertiesMap["name"]; ok {
		//							recordField.Name = helper.String(v.(string))
		//						}
		//
		//						if v, ok := propertiesMap["value"]; ok {
		//							recordField.Value = helper.String(v.(string))
		//						}
		//
		//						integrationNodeSchema.Properties = append(integrationNodeSchema.Properties, &recordField)
		//					}
		//				}
		//
		//				if v, ok := schemaMap["alias"]; ok {
		//					integrationNodeSchema.Alias = helper.String(v.(string))
		//				}
		//
		//				if v, ok := schemaMap["comment"]; ok {
		//					integrationNodeSchema.Comment = helper.String(v.(string))
		//				}
		//
		//				integrationNodeInfo.Schema = append(integrationNodeInfo.Schema, &integrationNodeSchema)
		//			}
		//		}
		//		if nodeMappingMap, ok := helper.InterfaceToMap(nodesMap, "node_mapping"); ok {
		//			integrationNodeMapping := wedata.IntegrationNodeMapping{}
		//			if v, ok := nodeMappingMap["source_id"]; ok {
		//				integrationNodeMapping.SourceId = helper.String(v.(string))
		//			}
		//
		//			if v, ok := nodeMappingMap["sink_id"]; ok {
		//				integrationNodeMapping.SinkId = helper.String(v.(string))
		//			}
		//
		//			if v, ok := nodeMappingMap["source_schema"]; ok {
		//				for _, item := range v.([]interface{}) {
		//					sourceSchemaMap := item.(map[string]interface{})
		//					integrationNodeSchema := wedata.IntegrationNodeSchema{}
		//					if v, ok := sourceSchemaMap["id"]; ok {
		//						integrationNodeSchema.Id = helper.String(v.(string))
		//					}
		//
		//					if v, ok := sourceSchemaMap["name"]; ok {
		//						integrationNodeSchema.Name = helper.String(v.(string))
		//					}
		//
		//					if v, ok := sourceSchemaMap["type"]; ok {
		//						integrationNodeSchema.Type = helper.String(v.(string))
		//					}
		//
		//					if v, ok := sourceSchemaMap["value"]; ok {
		//						integrationNodeSchema.Value = helper.String(v.(string))
		//					}
		//
		//					if v, ok := sourceSchemaMap["properties"]; ok {
		//						for _, item := range v.([]interface{}) {
		//							propertiesMap := item.(map[string]interface{})
		//							recordField := wedata.RecordField{}
		//							if v, ok := propertiesMap["name"]; ok {
		//								recordField.Name = helper.String(v.(string))
		//							}
		//
		//							if v, ok := propertiesMap["value"]; ok {
		//								recordField.Value = helper.String(v.(string))
		//							}
		//
		//							integrationNodeSchema.Properties = append(integrationNodeSchema.Properties, &recordField)
		//						}
		//					}
		//
		//					if v, ok := sourceSchemaMap["alias"]; ok {
		//						integrationNodeSchema.Alias = helper.String(v.(string))
		//					}
		//
		//					if v, ok := sourceSchemaMap["comment"]; ok {
		//						integrationNodeSchema.Comment = helper.String(v.(string))
		//					}
		//
		//					integrationNodeMapping.SourceSchema = append(integrationNodeMapping.SourceSchema, &integrationNodeSchema)
		//				}
		//			}
		//
		//			if v, ok := nodeMappingMap["schema_mappings"]; ok {
		//				for _, item := range v.([]interface{}) {
		//					schemaMappingsMap := item.(map[string]interface{})
		//					integrationNodeSchemaMapping := wedata.IntegrationNodeSchemaMapping{}
		//					if v, ok := schemaMappingsMap["source_schema_id"]; ok {
		//						integrationNodeSchemaMapping.SourceSchemaId = helper.String(v.(string))
		//					}
		//
		//					if v, ok := schemaMappingsMap["sink_schema_id"]; ok {
		//						integrationNodeSchemaMapping.SinkSchemaId = helper.String(v.(string))
		//					}
		//
		//					integrationNodeMapping.SchemaMappings = append(integrationNodeMapping.SchemaMappings, &integrationNodeSchemaMapping)
		//				}
		//			}
		//
		//			if v, ok := nodeMappingMap["ext_config"]; ok {
		//				for _, item := range v.([]interface{}) {
		//					extConfigMap := item.(map[string]interface{})
		//					recordField := wedata.RecordField{}
		//					if v, ok := extConfigMap["name"]; ok {
		//						recordField.Name = helper.String(v.(string))
		//					}
		//
		//					if v, ok := extConfigMap["value"]; ok {
		//						recordField.Value = helper.String(v.(string))
		//					}
		//
		//					integrationNodeMapping.ExtConfig = append(integrationNodeMapping.ExtConfig, &recordField)
		//				}
		//			}
		//
		//			integrationNodeInfo.NodeMapping = &integrationNodeMapping
		//		}
		//
		//		if v, ok := nodesMap["app_id"]; ok {
		//			integrationNodeInfo.AppId = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["project_id"]; ok {
		//			integrationNodeInfo.ProjectId = helper.String(v.(string))
		//		}
		//
		//		if v, ok := nodesMap["creator_uin"]; ok {
		//			integrationNodeInfo.CreatorUin = helper.String(v.(string))
		//		}
		//		if v, ok := nodesMap["operator_uin"]; ok {
		//			integrationNodeInfo.OperatorUin = helper.String(v.(string))
		//		}
		//		if v, ok := nodesMap["owner_uin"]; ok {
		//			integrationNodeInfo.OwnerUin = helper.String(v.(string))
		//		}
		//		if v, ok := nodesMap["create_time"]; ok {
		//			integrationNodeInfo.CreateTime = helper.String(v.(string))
		//		}
		//		if v, ok := nodesMap["update_time"]; ok {
		//			integrationNodeInfo.UpdateTime = helper.String(v.(string))
		//		}
		//		integrationTaskInfo.Nodes = append(integrationTaskInfo.Nodes, &integrationNodeInfo)
		//	}
		//}

		if v, ok := dMap["executor_id"]; ok {
			integrationTaskInfo.ExecutorId = helper.String(v.(string))
		}

		if v, ok := dMap["config"]; ok {
			for _, item := range v.([]interface{}) {
				configMap := item.(map[string]interface{})
				recordField := wedata.RecordField{}
				if v, ok := configMap["name"]; ok {
					recordField.Name = helper.String(v.(string))
				}

				if v, ok := configMap["value"]; ok {
					recordField.Value = helper.String(v.(string))
				}

				integrationTaskInfo.Config = append(integrationTaskInfo.Config, &recordField)
			}
		}

		if v, ok := dMap["ext_config"]; ok {
			for _, item := range v.([]interface{}) {
				extConfigMap := item.(map[string]interface{})
				recordField := wedata.RecordField{}
				if v, ok := extConfigMap["name"]; ok {
					recordField.Name = helper.String(v.(string))
				}

				if v, ok := extConfigMap["value"]; ok {
					recordField.Value = helper.String(v.(string))
				}

				integrationTaskInfo.ExtConfig = append(integrationTaskInfo.ExtConfig, &recordField)
			}
		}

		if v, ok := dMap["execute_context"]; ok {
			for _, item := range v.([]interface{}) {
				executeContextMap := item.(map[string]interface{})
				recordField := wedata.RecordField{}
				if v, ok := executeContextMap["name"]; ok {
					recordField.Name = helper.String(v.(string))
				}

				if v, ok := executeContextMap["value"]; ok {
					recordField.Value = helper.String(v.(string))
				}

				integrationTaskInfo.ExecuteContext = append(integrationTaskInfo.ExecuteContext, &recordField)
			}
		}

		if v, ok := dMap["mappings"]; ok {
			for _, item := range v.([]interface{}) {
				mappingsMap := item.(map[string]interface{})
				integrationNodeMapping := wedata.IntegrationNodeMapping{}
				if v, ok := mappingsMap["source_id"]; ok {
					integrationNodeMapping.SourceId = helper.String(v.(string))
				}

				if v, ok := mappingsMap["sink_id"]; ok {
					integrationNodeMapping.SinkId = helper.String(v.(string))
				}

				if v, ok := mappingsMap["source_schema"]; ok {
					for _, item := range v.([]interface{}) {
						sourceSchemaMap := item.(map[string]interface{})
						integrationNodeSchema := wedata.IntegrationNodeSchema{}
						if v, ok := sourceSchemaMap["id"]; ok {
							integrationNodeSchema.Id = helper.String(v.(string))
						}

						if v, ok := sourceSchemaMap["name"]; ok {
							integrationNodeSchema.Name = helper.String(v.(string))
						}

						if v, ok := sourceSchemaMap["type"]; ok {
							integrationNodeSchema.Type = helper.String(v.(string))
						}

						if v, ok := sourceSchemaMap["value"]; ok {
							integrationNodeSchema.Value = helper.String(v.(string))
						}

						if v, ok := sourceSchemaMap["properties"]; ok {
							for _, item := range v.([]interface{}) {
								propertiesMap := item.(map[string]interface{})
								recordField := wedata.RecordField{}
								if v, ok := propertiesMap["name"]; ok {
									recordField.Name = helper.String(v.(string))
								}

								if v, ok := propertiesMap["value"]; ok {
									recordField.Value = helper.String(v.(string))
								}

								integrationNodeSchema.Properties = append(integrationNodeSchema.Properties, &recordField)
							}
						}

						if v, ok := sourceSchemaMap["alias"]; ok {
							integrationNodeSchema.Alias = helper.String(v.(string))
						}

						if v, ok := sourceSchemaMap["comment"]; ok {
							integrationNodeSchema.Comment = helper.String(v.(string))
						}

						integrationNodeMapping.SourceSchema = append(integrationNodeMapping.SourceSchema, &integrationNodeSchema)
					}
				}

				if v, ok := mappingsMap["schema_mappings"]; ok {
					for _, item := range v.([]interface{}) {
						schemaMappingsMap := item.(map[string]interface{})
						integrationNodeSchemaMapping := wedata.IntegrationNodeSchemaMapping{}
						if v, ok := schemaMappingsMap["source_schema_id"]; ok {
							integrationNodeSchemaMapping.SourceSchemaId = helper.String(v.(string))
						}

						if v, ok := schemaMappingsMap["sink_schema_id"]; ok {
							integrationNodeSchemaMapping.SinkSchemaId = helper.String(v.(string))
						}

						integrationNodeMapping.SchemaMappings = append(integrationNodeMapping.SchemaMappings, &integrationNodeSchemaMapping)
					}
				}

				if v, ok := mappingsMap["ext_config"]; ok {
					for _, item := range v.([]interface{}) {
						extConfigMap := item.(map[string]interface{})
						recordField := wedata.RecordField{}
						if v, ok := extConfigMap["name"]; ok {
							recordField.Name = helper.String(v.(string))
						}

						if v, ok := extConfigMap["value"]; ok {
							recordField.Value = helper.String(v.(string))
						}

						integrationNodeMapping.ExtConfig = append(integrationNodeMapping.ExtConfig, &recordField)
					}
				}

				integrationTaskInfo.Mappings = append(integrationTaskInfo.Mappings, &integrationNodeMapping)
			}
		}

		if v, ok := dMap["task_mode"]; ok {
			integrationTaskInfo.TaskMode = helper.String(v.(string))
		}

		if v, ok := dMap["incharge"]; ok {
			integrationTaskInfo.Incharge = helper.String(v.(string))
		}

		if offlineTaskAddEntityMap, ok := helper.InterfaceToMap(dMap, "offline_task_add_entity"); ok {
			offlineTaskAddParam := wedata.OfflineTaskAddParam{}
			//if v, ok := offlineTaskAddEntityMap["workflow_name"]; ok {
			//	offlineTaskAddParam.WorkflowName = helper.String(v.(string))
			//}
			//
			//if v, ok := offlineTaskAddEntityMap["dependency_workflow"]; ok {
			//	offlineTaskAddParam.DependencyWorkflow = helper.String(v.(string))
			//}

			offlineTaskAddParam.StartTime = &startTime
			offlineTaskAddParam.EndTime = &endTime
			offlineTaskAddParam.CycleStep = helper.IntUint64(cycleStep)
			offlineTaskAddParam.DelayTime = helper.IntUint64(delayTime)
			offlineTaskAddParam.TaskAction = &taskAction

			if v, ok := offlineTaskAddEntityMap["cycle_type"]; ok {
				offlineTaskAddParam.CycleType = helper.IntUint64(v.(int))
			}

			if v, ok := offlineTaskAddEntityMap["crontab_expression"]; ok {
				offlineTaskAddParam.CrontabExpression = helper.String(v.(string))
			}

			if v, ok := offlineTaskAddEntityMap["retry_wait"]; ok {
				offlineTaskAddParam.RetryWait = helper.IntUint64(v.(int))
			}

			if v, ok := offlineTaskAddEntityMap["retriable"]; ok {
				offlineTaskAddParam.Retriable = helper.IntUint64(v.(int))
			}

			if v, ok := offlineTaskAddEntityMap["try_limit"]; ok {
				offlineTaskAddParam.TryLimit = helper.IntUint64(v.(int))
			}

			//if v, ok := offlineTaskAddEntityMap["run_priority"]; ok {
			//	offlineTaskAddParam.RunPriority = helper.IntUint64(v.(int))
			//}
			//
			//if v, ok := offlineTaskAddEntityMap["product_name"]; ok {
			//	offlineTaskAddParam.ProductName = helper.String(v.(string))
			//}

			if v, ok := offlineTaskAddEntityMap["self_depend"]; ok {
				offlineTaskAddParam.SelfDepend = helper.IntUint64(v.(int))
			}

			if v, ok := offlineTaskAddEntityMap["execution_end_time"]; ok {
				offlineTaskAddParam.ExecutionEndTime = helper.String(v.(string))
			}

			if v, ok := offlineTaskAddEntityMap["execution_start_time"]; ok {
				offlineTaskAddParam.ExecutionStartTime = helper.String(v.(string))
			}

			//if v, ok := offlineTaskAddEntityMap["task_auto_submit"]; ok {
			//	offlineTaskAddParam.TaskAutoSubmit = helper.Bool(v.(bool))
			//}
			//
			//if v, ok := offlineTaskAddEntityMap["instance_init_strategy"]; ok {
			//	offlineTaskAddParam.InstanceInitStrategy = helper.String(v.(string))
			//}

			integrationTaskInfo.OfflineTaskAddEntity = &offlineTaskAddParam
		}

		if v, ok := dMap["executor_group_name"]; ok {
			integrationTaskInfo.ExecutorGroupName = helper.String(v.(string))
		}

		if v, ok := dMap["in_long_manager_url"]; ok {
			integrationTaskInfo.InLongManagerUrl = helper.String(v.(string))
		}

		if v, ok := dMap["in_long_stream_id"]; ok {
			integrationTaskInfo.InLongStreamId = helper.String(v.(string))
		}

		if v, ok := dMap["in_long_manager_version"]; ok {
			integrationTaskInfo.InLongManagerVersion = helper.String(v.(string))
		}

		if v, ok := dMap["data_proxy_url"]; ok {
			dataProxyUrlSet := v.(*schema.Set).List()
			for i := range dataProxyUrlSet {
				dataProxyUrl := dataProxyUrlSet[i].(string)
				integrationTaskInfo.DataProxyUrl = append(integrationTaskInfo.DataProxyUrl, &dataProxyUrl)
			}
		}

		if v, ok := dMap["submit"]; ok {
			integrationTaskInfo.Submit = helper.Bool(v.(bool))
		}

		if v, ok := dMap["input_datasource_type"]; ok {
			integrationTaskInfo.InputDatasourceType = helper.String(v.(string))
		}

		if v, ok := dMap["output_datasource_type"]; ok {
			integrationTaskInfo.OutputDatasourceType = helper.String(v.(string))
		}

		if v, ok := dMap["num_records_in"]; ok {
			integrationTaskInfo.NumRecordsIn = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["num_records_out"]; ok {
			integrationTaskInfo.NumRecordsOut = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["reader_delay"]; ok {
			integrationTaskInfo.ReaderDelay = helper.Float64(v.(float64))
		}

		if v, ok := dMap["num_restarts"]; ok {
			integrationTaskInfo.NumRestarts = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["create_time"]; ok {
			integrationTaskInfo.CreateTime = helper.String(v.(string))
		}

		if v, ok := dMap["update_time"]; ok {
			integrationTaskInfo.UpdateTime = helper.String(v.(string))
		}

		if v, ok := dMap["last_run_time"]; ok {
			integrationTaskInfo.LastRunTime = helper.String(v.(string))
		}

		if v, ok := dMap["stop_time"]; ok {
			integrationTaskInfo.StopTime = helper.String(v.(string))
		}

		if v, ok := dMap["has_version"]; ok {
			integrationTaskInfo.HasVersion = helper.Bool(v.(bool))
		}

		if v, ok := dMap["locked"]; ok {
			integrationTaskInfo.Locked = helper.Bool(v.(bool))
		}

		if v, ok := dMap["locker"]; ok {
			integrationTaskInfo.Locker = helper.String(v.(string))
		}

		if v, ok := dMap["running_cu"]; ok {
			integrationTaskInfo.RunningCu = helper.Float64(v.(float64))
		}

		if v, ok := dMap["task_alarm_regular_list"]; ok {
			taskAlarmRegularListSet := v.(*schema.Set).List()
			for i := range taskAlarmRegularListSet {
				taskAlarmRegularList := taskAlarmRegularListSet[i].(string)
				integrationTaskInfo.TaskAlarmRegularList = append(integrationTaskInfo.TaskAlarmRegularList, &taskAlarmRegularList)
			}
		}

		if v, ok := dMap["switch_resource"]; ok {
			integrationTaskInfo.SwitchResource = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["read_phase"]; ok {
			integrationTaskInfo.ReadPhase = helper.IntInt64(v.(int))
		}

		if v, ok := dMap["instance_version"]; ok {
			integrationTaskInfo.InstanceVersion = helper.IntInt64(v.(int))
		}

		request.TaskInfo = &integrationTaskInfo
	}

	request.ProjectId = &projectId
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataClient().ModifyIntegrationTask(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s modify wedata integration_real_time_task failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWedataIntegrationOfflineTaskRead(d, meta)
}

func resourceTencentCloudWedataIntegrationOfflineTaskDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_integration_offline_task.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	projectId := idSplit[0]
	taskId := idSplit[1]

	if err := service.DeleteWedataIntegrationOfflineTaskById(ctx, projectId, taskId); err != nil {
		return err
	}

	return nil
}
