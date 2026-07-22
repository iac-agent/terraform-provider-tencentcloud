package ckafka

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCkafkaDatahubTask() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaDatahubTaskRead,
		Schema: map[string]*schema.Schema{
			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "search 键",
			},

			"target_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Destination 类型 dump。",
			},

			"task_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Task 类型，SOURCE|SINK。",
			},

			"source_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The 来源 类型",
			},

			"resource": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Resource。",
			},

			"task_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Datahub task information list。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "task ID。",
						},
						"task_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "TaskName。",
						},
						"task_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "TaskType，SOURCE|SINK。",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "状态，-1 failed to create，0 to create，1 to run，2 to delete，3 to deleted，4 to delete failed，5 to pause，6 to pause，7 to pause，8 to resume，9 to resume failed。",
						},
						"source_resource": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "data resource。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "资源类型",
									},
									"kafka_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "ckafka configuration，必填 when 类型 is KAFKA。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否cluster is built by yourself instead of cloud product。",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "instance resource。",
												},
												"topic": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Topic 名称，use `,` when more than 1 topic。",
												},
												"offset_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "偏移量 类型，from beginning:earliest，from latest:latest，from specific time:时间戳。",
												},
												"start_time": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "when 偏移量 类型 时间戳 为必填项。",
												},
												"resource_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例名称",
												},
												"zone_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "可用区 ID",
												},
												"topic_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Topic Id。",
												},
												"partition_num": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "the partition num of the topic。",
												},
												"enable_toleration": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "enable dead letter queue。",
												},
												"qps_limit": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Qps(query per seconds) 限制",
												},
												"table_mappings": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "maps of table to topic，必填 when multi topic is selected。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"database": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "database 名称",
															},
															"table": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "table 名称",
															},
															"topic": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Topic 名称",
															},
															"topic_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Topic ID。",
															},
														},
													},
												},
												"use_table_mapping": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否use multi table。",
												},
												"use_auto_create_topic": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Does the used topic need to be automatically created (currently only supports SOURCE inflow tasks，if you do not use to distribute to multiple topics，you need to fill in the topic 名称 that needs to be automatically created in the Topic field)。",
												},
												"compression_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "是否compress when writing to the Topic，if it is not 已启用，fill in none，if it is 已启用，fill in open。",
												},
												"msg_multiple": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "1 来源 topic 消息 is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka)。",
												},
												"connector_sync_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ConnectorSyncType。",
												},
												"keep_partition": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "KeepPartition。",
												},
											},
										},
									},
									"event_bus_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "EB configuration，必填 when 类型 is EB。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "资源类型 EB_COS/EB_ES/EB_CLS。",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否为a 自建集群",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 ID",
												},
												"namespace": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SCF namespace。",
												},
												"function_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SCF function 名称",
												},
												"qualifier": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SCF 版本 and alias。",
												},
											},
										},
									},
									"mongo_db_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "MongoDB 配置，必填 when 类型 is MONGODB。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MongoDB database 名称",
												},
												"collection": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MongoDB collection。",
												},
												"copy_existing": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否copy the stock data，the default parameter is true。",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "资源 ID",
												},
												"ip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Mongo DB connection ip。",
												},
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "MongoDB connection 端口",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MongoDB database 用户 名称",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MongoDB database 密码",
												},
												"listening_event": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Listening event 类型，if it is empty，it means select all. Values include insert，update，replace，delete，invalidate，drop，dropdatabase，rename，used between multiple types，separated by commas。",
												},
												"read_preference": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Master-slave 优先级，default master node。",
												},
												"pipeline": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "aggregation pipeline。",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否为a 自建集群",
												},
											},
										},
									},
									"es_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Es configuration，必填 when 类型 is ES。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Resource。",
												},
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Es connection 端口",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Es UserName。",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Es 密码",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否为a 自建集群",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 VIP",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "instance 私有网络 ID",
												},
												"drop_invalid_message": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether Es discards the 消息 of parsing failure。",
												},
												"index": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Es 索引 名称",
												},
												"date_format": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Es date suffix。",
												},
												"content_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "键 for data in non-json 格式",
												},
												"drop_invalid_json_message": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether Es discards messages in non-json 格式",
												},
												"document_id_field": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The field 名称 document ID 值 dumped into Es。",
												},
												"index_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Es custom 索引 名称 类型，STRING，JSONPATH，the 默认为 STRING。",
												},
												"drop_cls": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "When the member parameter Drop Invalid 消息 To Cls is set to true，the Drop Invalid 消息 parameter is invalid。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"drop_invalid_message_to_cls": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否deliver to cls。",
															},
															"drop_cls_region": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The 地域 where the cls is delivered。",
															},
															"drop_cls_owneruin": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Delivery 账号 of cls。",
															},
															"drop_cls_topic_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "topic of cls。",
															},
															"drop_cls_log_set": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "cls log set。",
															},
														},
													},
												},
												"database_primary_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "When the 消息 dumped to ES is the binlog of Database，if you need to synchronize database operations，that is，fill in the primary 键 of the database table when adding，deleting，and modifying operations to ES。",
												},
												"drop_dlq": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "dead letter queue。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "类型，DLQ dead letter queue，IGNORE_ERROR|DROP。",
															},
															"kafka_param": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Ckafka 类型 dlq。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"self_built": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "是否为a 自建集群",
																		},
																		"resource": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "资源 ID",
																		},
																		"topic": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Topic 名称，multiple separated by,。",
																		},
																		"offset_type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "偏移量 类型，initial position earliest，latest position latest，time point position 时间戳。",
																		},
																		"start_time": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "It must be passed when the 偏移量 类型 is 时间戳，and the time stamp is passed，accurate to the second。",
																		},
																		"resource_name": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "资源 ID 名称",
																		},
																		"zone_id": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "可用区 ID",
																		},
																		"topic_id": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Topic Id。",
																		},
																		"partition_num": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Partition num。",
																		},
																		"enable_toleration": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Enable the fault-tolerant instance and enable the dead-letter queue。",
																		},
																		"qps_limit": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "QPS 限制",
																		},
																		"table_mappings": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The route from Table to Topic must be passed when the Distribute to multiple topics switch is turned on。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"database": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "database 名称",
																					},
																					"table": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Table 名称，multiple tables，separated by (commas)。",
																					},
																					"topic": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Topic 名称",
																					},
																					"topic_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Topic ID。",
																					},
																				},
																			},
																		},
																		"use_table_mapping": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Distribute to multiple topics switch，the 默认为 false。",
																		},
																		"use_auto_create_topic": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "是否used topic need to be automatically created (currently only supports SOURCE inflow tasks，if you do not use to distribute to multiple topics，you need to fill in the topic 名称 that needs to be automatically created in the Topic field)。",
																		},
																		"compression_type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "是否compress when writing to the Topic，if it is not 已启用，fill in none，if it is 已启用，fill in open。",
																		},
																		"msg_multiple": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "1 来源 topic 消息 is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka)。",
																		},
																		"connector_sync_type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "ConnectorSyncType。",
																		},
																		"keep_partition": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "KeepPartition。",
																		},
																	},
																},
															},
															"retry_interval": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "重试间隔",
															},
															"max_retry_attempts": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "重试次数",
															},
															"topic_param": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "DIP Topic 类型 dead letter queue。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"resource": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The topic 名称 topic sold separately。",
																		},
																		"offset_type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "偏移量 类型，initial position earliest，latest position latest，time point position 时间戳。",
																		},
																		"start_time": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "It must be passed when the 偏移量 类型 is 时间戳，and the time stamp is passed，accurate to the second。",
																		},
																		"topic_id": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "TopicId。",
																		},
																		"compression_type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "是否perform compression when writing a topic，if it is not 已启用，fill in none，if it is 已启用，you can choose one of gzip，snappy，lz4 to fill in。",
																		},
																		"use_auto_create_topic": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "是否used topic need to be automatically created (currently only supports SOURCE inflow tasks)。",
																		},
																		"msg_multiple": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "1 来源 topic 消息 is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka)。",
																		},
																	},
																},
															},
															"dlq_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "dlq 类型，CKAFKA|TOPIC。",
															},
														},
													},
												},
											},
										},
									},
									"tdw_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Tdw configuration，必填 when 类型 is TDW。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bid": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Tdw bid。",
												},
												"tid": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Tdw tid。",
												},
												"is_domestic": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "default true。",
												},
												"tdw_host": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "TDW 地址，defalt tl-tdbank-tdmanager.tencent-distribute.com。",
												},
												"tdw_port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "TDW 端口，default 8099。",
												},
											},
										},
									},
									"dts_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Dts configuration，必填 when 类型 is DTS。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Dts instance Id。",
												},
												"ip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Dts connection ip。",
												},
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Dts connection 端口",
												},
												"topic": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Dts topic。",
												},
												"group_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Dts consumer group Id。",
												},
												"group_user": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Dts 账号",
												},
												"group_password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Dts consumer group passwd。",
												},
												"tran_sql": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "False to synchronize the original data，true to synchronize the parsed json 格式 data，the 默认为 true。",
												},
											},
										},
									},
									"click_house_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "ClickHouse 配置，类型 CLICKHOUSE requierd。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cluster": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ClickHouse cluster。",
												},
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ClickHouse database 名称",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ClickHouse table。",
												},
												"schema": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "ClickHouse schema。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"column_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "column 名称",
															},
															"json_key": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The json 键 名称 corresponding to this column。",
															},
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "类型 table column。",
															},
															"allow_null": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否column item is allowed to be empty。",
															},
														},
													},
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "资源 ID",
												},
												"ip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ClickHouse ip。",
												},
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "ClickHouse 端口",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ClickHouse 用户 名称",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ClickHouse passwd。",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 VIP",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "instance 私有网络 ID",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否为a 自建集群",
												},
												"drop_invalid_message": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether ClickHouse discards the 消息 that fails to parse，the 默认为 true。",
												},
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ClickHouse type, emr-clickhouse: emr; cdw-clickhouse: cdwch; selfBuilt: \"\".",
												},
												"drop_cls": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "When the member parameter Drop Invalid 消息 To Cls is set to true，the Drop Invalid 消息 parameter is invalid。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"drop_invalid_message_to_cls": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否deliver to cls。",
															},
															"drop_cls_region": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "cls 地域",
															},
															"drop_cls_owneruin": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "cls 账号",
															},
															"drop_cls_topic_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "cls topicId。",
															},
															"drop_cls_log_set": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "CLS 日志集 ID",
															},
														},
													},
												},
											},
										},
									},
									"cls_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Cls configuration，必填 when 类型 is CLS。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"decode_json": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否produced information is in json 格式",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "cls id。",
												},
												"log_set": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "LogSet id。",
												},
												"content_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "必填 when Decode Json is false。",
												},
												"time_field": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "指定content of a field in the 消息 as the time of the cls log. The 格式 of the field 内容 needs to be a second-级别 时间戳。",
												},
											},
										},
									},
									"cos_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Cos configuration，必填 when 类型 is COS。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bucket_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "cos 存储桶名称",
												},
												"region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "地域 代码",
												},
												"object_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ObjectKey。",
												},
												"aggregate_batch_size": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The size of aggregated messages MB。",
												},
												"aggregate_interval": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "时间间隔。",
												},
												"format_output_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The file 格式 after 消息 aggregation csv|json。",
												},
												"object_key_prefix": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Dumped object directory prefix。",
												},
												"directory_time_format": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Partition 格式 formatted according to strptime time。",
												},
											},
										},
									},
									"my_sql_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "MySQL configuration，必填 when 类型 is MYSQL。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MySQL database 名称，* is the whole database。",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 MySQL data table， is the non-system table in all the monitored databases，which can be separated by，to monitor multiple data tables，but the data table needs to be filled in the 格式 of data database 名称data table 名称，when a regular expression needs to be filled in，the 格式 is data database 名称data table 名称",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MySQL connection Id。",
												},
												"snapshot_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "是否Copy inventory information (schema_only does not copy，initial full amount)，the 默认为 initial。",
												},
												"ddl_topic": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The Topic that stores the Ddl information of My SQL，if it is empty，it will not be stored by default。",
												},
												"data_source_monitor_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "TABLE 表示that the read item is a table，QUERY 表示that the read item is a query。",
												},
												"data_source_monitor_resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "When DataMonitorMode=TABLE，pass in the Table that needs to be read; when DataMonitorMode=QUERY，pass in the query sql statement that needs to be read。",
												},
												"data_source_increment_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "TIMESTAMP 表示that the incremental column is of 时间戳 类型，INCREMENT 表示that the incremental column is of self-incrementing id 类型&#39;。",
												},
												"data_source_increment_column": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 column to be monitored。",
												},
												"data_source_start_from": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "HEAD means copy stock + incremental data，TAIL means copy only incremental data。",
												},
												"data_target_insert_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "INSERT means insert using Insert 模式，UPSERT means insert using Upsert 模式",
												},
												"data_target_primary_key_field": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "When DataInsertMode=UPSERT，pass in the primary 键 that the current upsert depends on。",
												},
												"data_target_record_mapping": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Mapping relationship between tables and messages。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"json_key": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The 键 名称 消息",
															},
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "消息类型",
															},
															"allow_null": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否message is allowed to be empty。",
															},
															"column_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Corresponding mapping column 名称",
															},
															"extra_info": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Database table extra fields。",
															},
															"column_size": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "current column size。",
															},
															"decimal_digits": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "current column precision。",
															},
															"auto_increment": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否为an auto-increment column。",
															},
															"default_value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Database table default parameters。",
															},
														},
													},
												},
												"topic_regex": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Regular expression for routing events to specific topics，默认为 (.*)。",
												},
												"topic_replacement": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "TopicRegex，$1，$2。",
												},
												"key_columns": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "格式 library1.table1 field 1,field 2;library 2.table2 field 2，between tables; (semicolon) separated，between fields，(comma) separated. The table that is not specified 默认为 the primary 键 of the table。",
												},
												"drop_invalid_message": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否discard messages that fail to parse，the 默认为 true。",
												},
												"drop_cls": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "When the member parameter Drop Invalid 消息 To Cls is set to true，the Drop Invalid 消息 parameter is invalid。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"drop_invalid_message_to_cls": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否deliver to cls。",
															},
															"drop_cls_region": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The 地域 where the cls is delivered。",
															},
															"drop_cls_owneruin": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "账号",
															},
															"drop_cls_topic_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "cls topic。",
															},
															"drop_cls_log_set": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "CLS 日志集 ID",
															},
														},
													},
												},
												"output_format": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "output 格式，DEFAULT，CANAL_1，CANAL_2。",
												},
												"is_table_prefix": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "When the Table input is a prefix，the 值 of this item is true，otherwise it is false。",
												},
												"include_content_changes": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "If the 值 is all，DDL data and DML data will also be written to the selected topic; if the 值 is dml，only DML data will be written to the selected topic。",
												},
												"include_query": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "If the 值 is true，and the 值 of the binlog rows query log events configuration item in My SQL is ON，the data flowing into the topic 包含original SQL statement; if the 值 is false，the data flowing into the topic does not contain Original SQL statement。",
												},
												"record_with_schema": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "If the 值 is true，the 消息 will carry the schema corresponding to the 消息 structure，if the 值 is false，it will not carry。",
												},
												"signal_database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "database 名称 signal table。",
												},
												"is_table_regular": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否input table is a regular expression，if this option and Is Table Prefix are true at the same time，the judgment 优先级 of this option is higher than Is Table Prefix。",
												},
											},
										},
									},
									"postgre_sql_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "PostgreSQL configuration，必填 when 类型 is POSTGRESQL or TDSQL C_POSTGRESQL。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "PostgreSQL database 名称",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "PostgreSQL tableName，* is the non-system table in all the monitored databases，you can use，to monitor multiple data tables，but the data table needs to be filled in the 格式 of Schema 名称Data table 名称，and you need to fill in a regular expression When，the 格式 is Schema 名称data table 名称",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "PostgreSQL connection Id。",
												},
												"plugin_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "(decoderbufs/pgoutput)，default decoderbufs。",
												},
												"snapshot_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "never|initial，default initial。",
												},
												"data_format": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Upstream data 格式 (JSON|Debezium)，必填 when the database synchronization 模式 matches the default field。",
												},
												"data_target_insert_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "INSERT means insert using Insert 模式，UPSERT means insert using Upsert 模式",
												},
												"data_target_primary_key_field": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "When DataInsertMode=UPSERT，pass in the primary 键 that the current upsert depends on。",
												},
												"data_target_record_mapping": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Mapping relationship between tables and messages。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"json_key": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The 键 名称 消息",
															},
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "消息类型",
															},
															"allow_null": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否message is allowed to be empty。",
															},
															"column_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Column 名称",
															},
															"extra_info": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Database table extra fields。",
															},
															"column_size": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "current ColumnSize。",
															},
															"decimal_digits": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "current Column DecimalDigits。",
															},
															"auto_increment": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否为an auto-increment column。",
															},
															"default_value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Database table default parameters。",
															},
														},
													},
												},
												"drop_invalid_message": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否discard messages that fail to parse，the 默认为 true。",
												},
												"is_table_regular": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否input table is a regular expression。",
												},
												"key_columns": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "格式  library1.table1:field 1,field2;library2.table2:field2，between tables; (semicolon) separated，between fields，(comma) separated. The table that is not specified 默认为 the primary 键 of the table。",
												},
												"record_with_schema": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "If the 值 is true，the 消息 will carry the schema corresponding to the 消息 structure，if the 值 is false，it will not carry。",
												},
											},
										},
									},
									"topic_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Topic configuration，必填 when 类型 is Topic。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The topic 名称 topic sold separately。",
												},
												"offset_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "偏移量 类型，initial position earliest，latest position latest，time point position 时间戳。",
												},
												"start_time": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "It must be passed when the 偏移量 类型 is 时间戳，and the time stamp is passed，accurate to the second。",
												},
												"topic_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Topic TopicId。",
												},
												"compression_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "是否perform compression when writing a topic，if it is not 已启用，fill in none，if it is 已启用，you can choose one of gzip，snappy，lz4 to fill in。",
												},
												"use_auto_create_topic": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否used topic need to be automatically created (currently only supports SOURCE inflow tasks)。",
												},
												"msg_multiple": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "1 来源 topic 消息 is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka)。",
												},
											},
										},
									},
									"maria_db_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "MariaDB configuration，必填 when 类型 is MARIADB。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MariaDB database 名称，* for all database。",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MariaDB db 名称，is the non-system table in all the monitored databases，you can use to monitor multiple data tables，but the data table needs to be filled in the 格式 of data database 名称data table 名称",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MariaDB connection Id。",
												},
												"snapshot_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "schema_only|initial，default initial。",
												},
												"key_columns": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "格式  library 1. table 1: field 1，field 2; library 2. table 2: field 2，between tables; (semicolon) separated，between fields，(comma) separated. The table that is not specified 默认为 the primary 键 of the table。",
												},
												"is_table_prefix": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "When the Table input is a prefix，the 值 of this item is true，otherwise it is false。",
												},
												"output_format": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "output 格式，DEFAULT，CANAL_1，CANAL_2。",
												},
												"include_content_changes": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "If the 值 is all，DDL data and DML data will also be written to the selected topic; if the 值 is dml，only DML data will be written to the selected topic。",
												},
												"include_query": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "If the 值 is true，and the 值 of the binlog rows query log events configuration item in My SQL is ON，the data flowing into the topic 包含original SQL statement; if the 值 is false，the data flowing into the topic does not contain Original SQL statement。",
												},
												"record_with_schema": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "If the 值 is true，the 消息 will carry the schema corresponding to the 消息 structure，if the 值 is false，it will not carry。",
												},
											},
										},
									},
									"sql_server_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "SQLServer configuration，必填 when 类型 is SQLSERVER。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SQLServer database 名称",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SQLServer table is the non-system table in all the monitored databases，you can use，to monitor multiple data tables，but the data table needs to be filled in the 格式 of data database 名称data table 名称",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SQLServer connection Id。",
												},
												"snapshot_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "schema_only|initial default initial。",
												},
											},
										},
									},
									"ctsdb_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Ctsdb configuration，必填 when 类型 is CTSDB。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "资源 ID",
												},
												"ctsdb_metric": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Ctsdb metric。",
												},
											},
										},
									},
									"scf_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Scf configuration，必填 when 类型 is SCF。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"function_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SCF function 名称",
												},
												"namespace": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SCF cloud function namespace，the 默认为 default。",
												},
												"qualifier": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SCF cloud function 版本 and alias，the 默认为 DEFAULT。",
												},
												"batch_size": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The 最大messages sent in each batch，the 默认为 1000。",
												},
												"max_retries": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The 数量 retries after the SCF call fails，the 默认为 5。",
												},
											},
										},
									},
								},
							},
						},
						"target_resource": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Target Resource。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "资源类型",
									},
									"kafka_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "ckafka configuration，必填 when 类型 is KAFKA。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否cluster is built by yourself instead of cloud product。",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "instance resource。",
												},
												"topic": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Topic 名称，use `,` when more than 1 topic。",
												},
												"offset_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "偏移量 类型，from beginning:earliest，from latest:latest，from specific time:时间戳。",
												},
												"start_time": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "when 偏移量 类型 时间戳 为必填项。",
												},
												"resource_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例名称",
												},
												"zone_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "可用区 ID",
												},
												"topic_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Topic Id。",
												},
												"partition_num": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "the partition num of the topic。",
												},
												"enable_toleration": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "enable dead letter queue。",
												},
												"qps_limit": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Qps(query per seconds) 限制",
												},
												"table_mappings": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "maps of table to topic，必填 when multi topic is selected。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"database": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "database 名称",
															},
															"table": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "table 名称",
															},
															"topic": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Topic 名称",
															},
															"topic_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Topic ID。",
															},
														},
													},
												},
												"use_table_mapping": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否use multi table。",
												},
												"use_auto_create_topic": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Does the used topic need to be automatically created (currently only supports SOURCE inflow tasks，if you do not use to distribute to multiple topics，you need to fill in the topic 名称 that needs to be automatically created in the Topic field)。",
												},
												"compression_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "是否compress when writing to the Topic，if it is not 已启用，fill in none，if it is 已启用，fill in open。",
												},
												"msg_multiple": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "1 来源 topic 消息 is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka)。",
												},
												"connector_sync_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ConnectorSyncType。",
												},
												"keep_partition": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "KeepPartition。",
												},
											},
										},
									},
									"event_bus_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "EB configuration，必填 when 类型 is EB。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "资源类型 EB_COS/EB_ES/EB_CLS。",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否为a 自建集群",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 ID",
												},
												"namespace": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SCF namespace。",
												},
												"function_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SCF function 名称",
												},
												"qualifier": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SCF 版本 and alias。",
												},
											},
										},
									},
									"mongo_db_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "MongoDB 配置，必填 when 类型 is MONGODB。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MongoDB database 名称",
												},
												"collection": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MongoDB collection。",
												},
												"copy_existing": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否copy the stock data，the default parameter is true。",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "资源 ID",
												},
												"ip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Mongo DB connection ip。",
												},
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "MongoDB connection 端口",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MongoDB database 用户 名称",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MongoDB database 密码",
												},
												"listening_event": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Listening event 类型，if it is empty，it means select all. Values include insert，update，replace，delete，invalidate，drop，dropdatabase，rename，used between multiple types，separated by commas。",
												},
												"read_preference": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Master-slave 优先级，default master node。",
												},
												"pipeline": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "aggregation pipeline。",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否为a 自建集群",
												},
											},
										},
									},
									"es_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Es configuration，必填 when 类型 is ES。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Resource。",
												},
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Es connection 端口",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Es UserName。",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Es 密码",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否为a 自建集群",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 VIP",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "instance 私有网络 ID",
												},
												"drop_invalid_message": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether Es discards the 消息 of parsing failure。",
												},
												"index": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Es 索引 名称",
												},
												"date_format": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Es date suffix。",
												},
												"content_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "键 for data in non-json 格式",
												},
												"drop_invalid_json_message": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether Es discards messages in non-json 格式",
												},
												"document_id_field": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The field 名称 document ID 值 dumped into Es。",
												},
												"index_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Es custom 索引 名称 类型，STRING，JSONPATH，the 默认为 STRING。",
												},
												"drop_cls": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "When the member parameter Drop Invalid 消息 To Cls is set to true，the Drop Invalid 消息 parameter is invalid。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"drop_invalid_message_to_cls": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否deliver to cls。",
															},
															"drop_cls_region": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The 地域 where the cls is delivered。",
															},
															"drop_cls_owneruin": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Delivery 账号 of cls。",
															},
															"drop_cls_topic_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "topic of cls。",
															},
															"drop_cls_log_set": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "cls log set。",
															},
														},
													},
												},
												"database_primary_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "When the 消息 dumped to ES is the binlog of Database，if you need to synchronize database operations，that is，fill in the primary 键 of the database table when adding，deleting，and modifying operations to ES。",
												},
												"drop_dlq": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "dead letter queue。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "类型，DLQ dead letter queue，IGNORE_ERROR|DROP。",
															},
															"kafka_param": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Ckafka 类型 dlq。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"self_built": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "是否为a 自建集群",
																		},
																		"resource": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "资源 ID",
																		},
																		"topic": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Topic 名称，multiple separated by。",
																		},
																		"offset_type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "偏移量 类型，initial position earliest，latest position latest，time point position 时间戳。",
																		},
																		"start_time": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "It must be passed when the 偏移量 类型 is 时间戳，and the time stamp is passed，accurate to the second。",
																		},
																		"resource_name": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "资源 ID 名称",
																		},
																		"zone_id": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "可用区 ID",
																		},
																		"topic_id": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Topic Id。",
																		},
																		"partition_num": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Partition num。",
																		},
																		"enable_toleration": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Enable the fault-tolerant instance and enable the dead-letter queue。",
																		},
																		"qps_limit": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "QPS 限制",
																		},
																		"table_mappings": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The route from Table to Topic must be passed when the Distribute to multiple topics switch is turned on。",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"database": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "database 名称",
																					},
																					"table": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Table 名称，multiple tables，separated by (commas)。",
																					},
																					"topic": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Topic 名称",
																					},
																					"topic_id": {
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Topic ID。",
																					},
																				},
																			},
																		},
																		"use_table_mapping": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Distribute to multiple topics switch，the 默认为 false。",
																		},
																		"use_auto_create_topic": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "是否used topic need to be automatically created (currently only supports SOURCE inflow tasks，if you do not use to distribute to multiple topics，you need to fill in the topic 名称 that needs to be automatically created in the Topic field)。",
																		},
																		"compression_type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "是否compress when writing to the Topic，if it is not 已启用，fill in none，if it is 已启用，fill in open。",
																		},
																		"msg_multiple": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "1 来源 topic 消息 is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka)。",
																		},
																		"connector_sync_type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "ConnectorSyncType。",
																		},
																		"keep_partition": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "KeepPartition。",
																		},
																	},
																},
															},
															"retry_interval": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "重试间隔",
															},
															"max_retry_attempts": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "重试次数",
															},
															"topic_param": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "DIP Topic 类型 dead letter queue。",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"resource": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The topic 名称 topic sold separately。",
																		},
																		"offset_type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "偏移量 类型，initial position earliest，latest position latest，time point position 时间戳。",
																		},
																		"start_time": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "It must be passed when the 偏移量 类型 is 时间戳，and the time stamp is passed，accurate to the second。",
																		},
																		"topic_id": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "TopicId。",
																		},
																		"compression_type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "是否perform compression when writing a topic，if it is not 已启用，fill in none，if it is 已启用，you can choose one of gzip，snappy，lz4 to fill in。",
																		},
																		"use_auto_create_topic": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "是否used topic need to be automatically created (currently only supports SOURCE inflow tasks)。",
																		},
																		"msg_multiple": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "1 来源 topic 消息 is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka)。",
																		},
																	},
																},
															},
															"dlq_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "dlq 类型，CKAFKA|TOPIC。",
															},
														},
													},
												},
											},
										},
									},
									"tdw_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Tdw configuration，必填 when 类型 is TDW。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bid": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Tdw bid。",
												},
												"tid": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Tdw tid。",
												},
												"is_domestic": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "default true。",
												},
												"tdw_host": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "TDW 地址，defalt tl-tdbank-tdmanager.tencent-distribute.com。",
												},
												"tdw_port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "TDW 端口，default 8099。",
												},
											},
										},
									},
									"dts_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Dts configuration，必填 when 类型 is DTS。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Dts instance Id。",
												},
												"ip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Dts connection ip。",
												},
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Dts connection 端口",
												},
												"topic": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Dts topic。",
												},
												"group_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Dts consumer group Id。",
												},
												"group_user": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Dts 账号",
												},
												"group_password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Dts consumer group passwd。",
												},
												"tran_sql": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "False to synchronize the original data，true to synchronize the parsed json 格式 data，the 默认为 true。",
												},
											},
										},
									},
									"click_house_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "ClickHouse 配置，类型 CLICKHOUSE requierd。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cluster": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ClickHouse cluster。",
												},
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ClickHouse database 名称",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ClickHouse table。",
												},
												"schema": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "ClickHouse schema。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"column_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "column 名称",
															},
															"json_key": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The json 键 名称 corresponding to this column。",
															},
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "类型 table column。",
															},
															"allow_null": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否column item is allowed to be empty。",
															},
														},
													},
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "资源 ID",
												},
												"ip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ClickHouse ip。",
												},
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "ClickHouse 端口",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ClickHouse 用户 名称",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ClickHouse passwd。",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 VIP",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "instance 私有网络 ID",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否为a 自建集群",
												},
												"drop_invalid_message": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether ClickHouse discards the 消息 that fails to parse，the 默认为 true。",
												},
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ClickHouse type, emr-clickhouse: emr; cdw-clickhouse: cdwch; selfBuilt: \"\".",
												},
												"drop_cls": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "When the member parameter Drop Invalid 消息 To Cls is set to true，the Drop Invalid 消息 parameter is invalid。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"drop_invalid_message_to_cls": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否deliver to cls。",
															},
															"drop_cls_region": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "cls 地域",
															},
															"drop_cls_owneruin": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "cls 账号",
															},
															"drop_cls_topic_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "cls topicId。",
															},
															"drop_cls_log_set": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "CLS 日志集 ID",
															},
														},
													},
												},
											},
										},
									},
									"cls_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Cls configuration，必填 when 类型 is CLS。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"decode_json": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否produced information is in json 格式",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "cls id。",
												},
												"log_set": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "LogSet id。",
												},
												"content_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "必填 when Decode Json is false。",
												},
												"time_field": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "指定content of a field in the 消息 as the time of the cls log. The 格式 of the field 内容 needs to be a second-级别 时间戳。",
												},
											},
										},
									},
									"cos_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Cos configuration，必填 when 类型 is COS。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bucket_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "cos 存储桶名称",
												},
												"region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "地域 代码",
												},
												"object_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ObjectKey。",
												},
												"aggregate_batch_size": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The size of aggregated messages MB。",
												},
												"aggregate_interval": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "时间间隔。",
												},
												"format_output_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The file 格式 after 消息 aggregation csv|json。",
												},
												"object_key_prefix": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Dumped object directory prefix。",
												},
												"directory_time_format": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Partition 格式 formatted according to strptime time。",
												},
											},
										},
									},
									"my_sql_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "MySQL configuration，必填 when 类型 is MYSQL。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MySQL database 名称，* is the whole database。",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 名称 MySQL data table is the non-system table in all the monitored databases，which can be separated by，to monitor multiple data tables，but the data table needs to be filled in the 格式 of data database 名称data table 名称，when a regular expression needs to be filled in，the 格式 is data database 名称data table 名称",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MySQL connection Id。",
												},
												"snapshot_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "是否Copy inventory information (schema_only does not copy，initial full amount)，the 默认为 initial。",
												},
												"ddl_topic": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The Topic that stores the Ddl information of My SQL，if it is empty，it will not be stored by default。",
												},
												"data_source_monitor_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "TABLE 表示that the read item is a table，QUERY 表示that the read item is a query。",
												},
												"data_source_monitor_resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "When DataMonitorMode=TABLE，pass in the Table that needs to be read; when DataMonitorMode=QUERY，pass in the query sql statement that needs to be read。",
												},
												"data_source_increment_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "TIMESTAMP 表示that the incremental column is of 时间戳 类型，INCREMENT 表示that the incremental column is of self-incrementing id 类型",
												},
												"data_source_increment_column": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "the 名称 column to be monitored。",
												},
												"data_source_start_from": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "HEAD means copy stock + incremental data，TAIL means copy only incremental data。",
												},
												"data_target_insert_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "INSERT means insert using Insert 模式，UPSERT means insert using Upsert 模式",
												},
												"data_target_primary_key_field": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "When DataInsertMode=UPSERT，pass in the primary 键 that the current upsert depends on。",
												},
												"data_target_record_mapping": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Mapping relationship between tables and messages。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"json_key": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The 键 名称 消息",
															},
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "消息类型",
															},
															"allow_null": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否message is allowed to be empty。",
															},
															"column_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Corresponding mapping column 名称",
															},
															"extra_info": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Database table extra fields。",
															},
															"column_size": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "current column size。",
															},
															"decimal_digits": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "current column precision。",
															},
															"auto_increment": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否为an auto-increment column。",
															},
															"default_value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Database table default parameters。",
															},
														},
													},
												},
												"topic_regex": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Regular expression for routing events to specific topics，默认为 (.*)。",
												},
												"topic_replacement": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "TopicRegex，$1，$2。",
												},
												"key_columns": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "格式 library1.table1 field 1,field 2;library 2.table2 field 2，between tables; (semicolon) separated，between fields，(comma) separated. The table that is not specified 默认为 the primary 键 of the table。",
												},
												"drop_invalid_message": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否discard messages that fail to parse，the 默认为 true。",
												},
												"drop_cls": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "When the member parameter Drop Invalid 消息 To Cls is set to true，the Drop Invalid 消息 parameter is invalid。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"drop_invalid_message_to_cls": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否deliver to cls。",
															},
															"drop_cls_region": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The 地域 where the cls is delivered。",
															},
															"drop_cls_owneruin": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "账号",
															},
															"drop_cls_topic_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "cls topic。",
															},
															"drop_cls_log_set": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "CLS 日志集 ID",
															},
														},
													},
												},
												"output_format": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "output 格式，DEFAULT，CANAL_1，CANAL_2。",
												},
												"is_table_prefix": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "When the Table input is a prefix，the 值 of this item is true，otherwise it is false。",
												},
												"include_content_changes": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "If the 值 is all，DDL data and DML data will also be written to the selected topic; if the 值 is dml，only DML data will be written to the selected topic。",
												},
												"include_query": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "If the 值 is true，and the 值 of the binlog rows query log events configuration item in My SQL is ON，the data flowing into the topic 包含original SQL statement; if the 值 is false，the data flowing into the topic does not contain Original SQL statement。",
												},
												"record_with_schema": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "If the 值 is true，the 消息 will carry the schema corresponding to the 消息 structure，if the 值 is false，it will not carry。",
												},
												"signal_database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "database 名称 signal table。",
												},
												"is_table_regular": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否input table is a regular expression，if this option and Is Table Prefix are true at the same time，the judgment 优先级 of this option is higher than Is Table Prefix。",
												},
											},
										},
									},
									"postgre_sql_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "PostgreSQL configuration，必填 when 类型 is POSTGRESQL or TDSQL C_POSTGRESQL。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "PostgreSQL database 名称",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "PostgreSQL tableName，is the non-system table in all the monitored databases，you can use，to monitor multiple data tables，but the data table needs to be filled in the 格式 of Schema 名称Data table 名称，and you need to fill in a regular expression When，the 格式 is Schema 名称data table 名称",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "PostgreSQL connection Id。",
												},
												"plugin_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "(decoderbufs/pgoutput)，default decoderbufs。",
												},
												"snapshot_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "never|initial，default initial。",
												},
												"data_format": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Upstream data 格式 (JSON|Debezium)，必填 when the database synchronization 模式 matches the default field。",
												},
												"data_target_insert_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "INSERT means insert using Insert 模式，UPSERT means insert using Upsert 模式",
												},
												"data_target_primary_key_field": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "When DataInsertMode=UPSERT，pass in the primary 键 that the current upsert depends on。",
												},
												"data_target_record_mapping": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Mapping relationship between tables and messages。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"json_key": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The 键 名称 消息",
															},
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "消息类型",
															},
															"allow_null": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否message is allowed to be empty。",
															},
															"column_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Column 名称",
															},
															"extra_info": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Database table extra fields。",
															},
															"column_size": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "current ColumnSize。",
															},
															"decimal_digits": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "current Column DecimalDigits。",
															},
															"auto_increment": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "是否为an auto-increment column。",
															},
															"default_value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Database table default parameters。",
															},
														},
													},
												},
												"drop_invalid_message": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否discard messages that fail to parse，the 默认为 true。",
												},
												"is_table_regular": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否input table is a regular expression。",
												},
												"key_columns": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "格式  library1.table1:field 1,field2;library2.table2:field2，between tables; (semicolon) separated，between fields，(comma) separated. The table that is not specified 默认为 the primary 键 of the table。",
												},
												"record_with_schema": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "If the 值 is true，the 消息 will carry the schema corresponding to the 消息 structure，if the 值 is false，it will not carry。",
												},
											},
										},
									},
									"topic_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Topic configuration，必填 when 类型 is Topic。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The topic 名称 topic sold separately。",
												},
												"offset_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "偏移量 类型，initial position earliest，latest position latest，time point position 时间戳。",
												},
												"start_time": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "It must be passed when the 偏移量 类型 is 时间戳，and the time stamp is passed，accurate to the second。",
												},
												"topic_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Topic TopicId。",
												},
												"compression_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "是否perform compression when writing a topic，if it is not 已启用，fill in none，if it is 已启用，you can choose one of gzip，snappy，lz4 to fill in。",
												},
												"use_auto_create_topic": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否used topic need to be automatically created (currently only supports SOURCE inflow tasks)。",
												},
												"msg_multiple": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "1 来源 topic 消息 is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka)。",
												},
											},
										},
									},
									"maria_db_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "MariaDB configuration，必填 when 类型 is MARIADB。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MariaDB database 名称，* for all database。",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MariaDB db 名称，is the non-system table in all the monitored databases，you can use，to monitor multiple data tables，but the data table needs to be filled in the 格式 of data database 名称data table 名称",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MariaDB connection Id。",
												},
												"snapshot_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "schema_only|initial，default initial。",
												},
												"key_columns": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "格式  library 1. table 1: field 1，field 2; library 2. table 2: field 2，between tables; (semicolon) separated，between fields，(comma) separated. The table that is not specified 默认为 the primary 键 of the table。",
												},
												"is_table_prefix": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "When the Table input is a prefix，the 值 of this item is true，otherwise it is false。",
												},
												"output_format": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "output 格式，DEFAULT，CANAL_1，CANAL_2。",
												},
												"include_content_changes": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "If the 值 is all，DDL data and DML data will also be written to the selected topic; if the 值 is dml，only DML data will be written to the selected topic。",
												},
												"include_query": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "If the 值 is true，and the 值 of the binlog rows query log events configuration item in My SQL is ON，the data flowing into the topic 包含original SQL statement; if the 值 is false，the data flowing into the topic does not contain Original SQL statement。",
												},
												"record_with_schema": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "If the 值 is true，the 消息 will carry the schema corresponding to the 消息 structure，if the 值 is false，it will not carry。",
												},
											},
										},
									},
									"sql_server_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "SQLServer configuration，必填 when 类型 is SQLSERVER。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SQLServer database 名称",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SQLServer table，is the non-system table in all the monitored databases，you can use，to monitor multiple data tables，but the data table needs to be filled in the 格式 of data database 名称data table 名称",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SQLServer connection Id。",
												},
												"snapshot_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "schema_only|initial default initial。",
												},
											},
										},
									},
									"ctsdb_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Ctsdb configuration，必填 when 类型 is CTSDB。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "资源 ID",
												},
												"ctsdb_metric": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Ctsdb metric。",
												},
											},
										},
									},
									"scf_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Scf configuration，必填 when 类型 is SCF。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"function_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SCF function 名称",
												},
												"namespace": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SCF cloud function namespace，the 默认为 default。",
												},
												"qualifier": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "SCF cloud function 版本 and alias，the 默认为 DEFAULT。",
												},
												"batch_size": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The 最大messages sent in each batch，the 默认为 1000。",
												},
												"max_retries": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The 数量 retries after the SCF call fails，the 默认为 5。",
												},
											},
										},
									},
								},
							},
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CreateTime。",
						},
						"error_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ErrorMessage。",
						},
						"task_progress": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Creation progress percentage。",
						},
						"task_current_step": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task Current Step。",
						},
						"datahub_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Datahub Id。",
						},
						"step_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "StepList。",
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

func dataSourceTencentCloudCkafkaDatahubTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ckafka_datahub_task.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["search_word"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_type"); ok {
		paramMap["target_type"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_type"); ok {
		paramMap["task_type"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_type"); ok {
		paramMap["source_type"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource"); ok {
		paramMap["resource"] = helper.String(v.(string))
	}

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var datahubTaskInfos []*ckafka.DatahubTaskInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCkafkaDatahubTaskByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		datahubTaskInfos = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(datahubTaskInfos))
	taskList := make([]map[string]interface{}, 0, len(datahubTaskInfos))
	for _, task := range datahubTaskInfos {
		taskMap := map[string]interface{}{}
		if task.TaskId != nil {
			taskMap["task_id"] = task.TaskId
			ids = append(ids, *task.TaskId)
		}

		if task.TaskName != nil {
			taskMap["task_name"] = task.TaskName
		}

		if task.TaskType != nil {
			taskMap["task_type"] = task.TaskType
		}

		if task.Status != nil {
			taskMap["status"] = task.Status
		}

		if task.SourceResource != nil {
			sourceResourceMap := map[string]interface{}{}

			if task.SourceResource.Type != nil {
				sourceResourceMap["type"] = task.SourceResource.Type
			}

			if task.SourceResource.KafkaParam != nil {
				kafkaParamMap := map[string]interface{}{}

				if task.SourceResource.KafkaParam.SelfBuilt != nil {
					kafkaParamMap["self_built"] = task.SourceResource.KafkaParam.SelfBuilt
				}

				if task.SourceResource.KafkaParam.Resource != nil {
					kafkaParamMap["resource"] = task.SourceResource.KafkaParam.Resource
				}

				if task.SourceResource.KafkaParam.Topic != nil {
					kafkaParamMap["topic"] = task.SourceResource.KafkaParam.Topic
				}

				if task.SourceResource.KafkaParam.OffsetType != nil {
					kafkaParamMap["offset_type"] = task.SourceResource.KafkaParam.OffsetType
				}

				if task.SourceResource.KafkaParam.StartTime != nil {
					kafkaParamMap["start_time"] = task.SourceResource.KafkaParam.StartTime
				}

				if task.SourceResource.KafkaParam.ResourceName != nil {
					kafkaParamMap["resource_name"] = task.SourceResource.KafkaParam.ResourceName
				}

				if task.SourceResource.KafkaParam.ZoneId != nil {
					kafkaParamMap["zone_id"] = task.SourceResource.KafkaParam.ZoneId
				}

				if task.SourceResource.KafkaParam.TopicId != nil {
					kafkaParamMap["topic_id"] = task.SourceResource.KafkaParam.TopicId
				}

				if task.SourceResource.KafkaParam.PartitionNum != nil {
					kafkaParamMap["partition_num"] = task.SourceResource.KafkaParam.PartitionNum
				}

				if task.SourceResource.KafkaParam.EnableToleration != nil {
					kafkaParamMap["enable_toleration"] = task.SourceResource.KafkaParam.EnableToleration
				}

				if task.SourceResource.KafkaParam.QpsLimit != nil {
					kafkaParamMap["qps_limit"] = task.SourceResource.KafkaParam.QpsLimit
				}

				if task.SourceResource.KafkaParam.TableMappings != nil {
					tableMappingsList := []interface{}{}
					for _, tableMappings := range task.SourceResource.KafkaParam.TableMappings {
						tableMappingsMap := map[string]interface{}{}

						if tableMappings.Database != nil {
							tableMappingsMap["database"] = tableMappings.Database
						}

						if tableMappings.Table != nil {
							tableMappingsMap["table"] = tableMappings.Table
						}

						if tableMappings.Topic != nil {
							tableMappingsMap["topic"] = tableMappings.Topic
						}

						if tableMappings.TopicId != nil {
							tableMappingsMap["topic_id"] = tableMappings.TopicId
						}

						tableMappingsList = append(tableMappingsList, tableMappingsMap)
					}

					kafkaParamMap["table_mappings"] = []interface{}{tableMappingsList}
				}

				if task.SourceResource.KafkaParam.UseTableMapping != nil {
					kafkaParamMap["use_table_mapping"] = task.SourceResource.KafkaParam.UseTableMapping
				}

				if task.SourceResource.KafkaParam.UseAutoCreateTopic != nil {
					kafkaParamMap["use_auto_create_topic"] = task.SourceResource.KafkaParam.UseAutoCreateTopic
				}

				if task.SourceResource.KafkaParam.CompressionType != nil {
					kafkaParamMap["compression_type"] = task.SourceResource.KafkaParam.CompressionType
				}

				if task.SourceResource.KafkaParam.MsgMultiple != nil {
					kafkaParamMap["msg_multiple"] = task.SourceResource.KafkaParam.MsgMultiple
				}

				if task.SourceResource.KafkaParam.ConnectorSyncType != nil {
					kafkaParamMap["connector_sync_type"] = task.SourceResource.KafkaParam.ConnectorSyncType
				}

				if task.SourceResource.KafkaParam.KeepPartition != nil {
					kafkaParamMap["keep_partition"] = task.SourceResource.KafkaParam.KeepPartition
				}

				sourceResourceMap["kafka_param"] = []interface{}{kafkaParamMap}
			}

			if task.SourceResource.EventBusParam != nil {
				eventBusParamMap := map[string]interface{}{}

				if task.SourceResource.EventBusParam.Type != nil {
					eventBusParamMap["type"] = task.SourceResource.EventBusParam.Type
				}

				if task.SourceResource.EventBusParam.SelfBuilt != nil {
					eventBusParamMap["self_built"] = task.SourceResource.EventBusParam.SelfBuilt
				}

				if task.SourceResource.EventBusParam.Resource != nil {
					eventBusParamMap["resource"] = task.SourceResource.EventBusParam.Resource
				}

				if task.SourceResource.EventBusParam.Namespace != nil {
					eventBusParamMap["namespace"] = task.SourceResource.EventBusParam.Namespace
				}

				if task.SourceResource.EventBusParam.FunctionName != nil {
					eventBusParamMap["function_name"] = task.SourceResource.EventBusParam.FunctionName
				}

				if task.SourceResource.EventBusParam.Qualifier != nil {
					eventBusParamMap["qualifier"] = task.SourceResource.EventBusParam.Qualifier
				}

				sourceResourceMap["event_bus_param"] = []interface{}{eventBusParamMap}
			}

			if task.SourceResource.MongoDBParam != nil {
				mongoDBParamMap := map[string]interface{}{}

				if task.SourceResource.MongoDBParam.Database != nil {
					mongoDBParamMap["database"] = task.SourceResource.MongoDBParam.Database
				}

				if task.SourceResource.MongoDBParam.Collection != nil {
					mongoDBParamMap["collection"] = task.SourceResource.MongoDBParam.Collection
				}

				if task.SourceResource.MongoDBParam.CopyExisting != nil {
					mongoDBParamMap["copy_existing"] = task.SourceResource.MongoDBParam.CopyExisting
				}

				if task.SourceResource.MongoDBParam.Resource != nil {
					mongoDBParamMap["resource"] = task.SourceResource.MongoDBParam.Resource
				}

				if task.SourceResource.MongoDBParam.Ip != nil {
					mongoDBParamMap["ip"] = task.SourceResource.MongoDBParam.Ip
				}

				if task.SourceResource.MongoDBParam.Port != nil {
					mongoDBParamMap["port"] = task.SourceResource.MongoDBParam.Port
				}

				if task.SourceResource.MongoDBParam.UserName != nil {
					mongoDBParamMap["user_name"] = task.SourceResource.MongoDBParam.UserName
				}

				if task.SourceResource.MongoDBParam.Password != nil {
					mongoDBParamMap["password"] = task.SourceResource.MongoDBParam.Password
				}

				if task.SourceResource.MongoDBParam.ListeningEvent != nil {
					mongoDBParamMap["listening_event"] = task.SourceResource.MongoDBParam.ListeningEvent
				}

				if task.SourceResource.MongoDBParam.ReadPreference != nil {
					mongoDBParamMap["read_preference"] = task.SourceResource.MongoDBParam.ReadPreference
				}

				if task.SourceResource.MongoDBParam.Pipeline != nil {
					mongoDBParamMap["pipeline"] = task.SourceResource.MongoDBParam.Pipeline
				}

				if task.SourceResource.MongoDBParam.SelfBuilt != nil {
					mongoDBParamMap["self_built"] = task.SourceResource.MongoDBParam.SelfBuilt
				}

				sourceResourceMap["mongo_db_param"] = []interface{}{mongoDBParamMap}
			}

			if task.SourceResource.EsParam != nil {
				esParamMap := map[string]interface{}{}

				if task.SourceResource.EsParam.Resource != nil {
					esParamMap["resource"] = task.SourceResource.EsParam.Resource
				}

				if task.SourceResource.EsParam.Port != nil {
					esParamMap["port"] = task.SourceResource.EsParam.Port
				}

				if task.SourceResource.EsParam.UserName != nil {
					esParamMap["user_name"] = task.SourceResource.EsParam.UserName
				}

				if task.SourceResource.EsParam.Password != nil {
					esParamMap["password"] = task.SourceResource.EsParam.Password
				}

				if task.SourceResource.EsParam.SelfBuilt != nil {
					esParamMap["self_built"] = task.SourceResource.EsParam.SelfBuilt
				}

				if task.SourceResource.EsParam.ServiceVip != nil {
					esParamMap["service_vip"] = task.SourceResource.EsParam.ServiceVip
				}

				if task.SourceResource.EsParam.UniqVpcId != nil {
					esParamMap["uniq_vpc_id"] = task.SourceResource.EsParam.UniqVpcId
				}

				if task.SourceResource.EsParam.DropInvalidMessage != nil {
					esParamMap["drop_invalid_message"] = task.SourceResource.EsParam.DropInvalidMessage
				}

				if task.SourceResource.EsParam.Index != nil {
					esParamMap["index"] = task.SourceResource.EsParam.Index
				}

				if task.SourceResource.EsParam.DateFormat != nil {
					esParamMap["date_format"] = task.SourceResource.EsParam.DateFormat
				}

				if task.SourceResource.EsParam.ContentKey != nil {
					esParamMap["content_key"] = task.SourceResource.EsParam.ContentKey
				}

				if task.SourceResource.EsParam.DropInvalidJsonMessage != nil {
					esParamMap["drop_invalid_json_message"] = task.SourceResource.EsParam.DropInvalidJsonMessage
				}

				if task.SourceResource.EsParam.DocumentIdField != nil {
					esParamMap["document_id_field"] = task.SourceResource.EsParam.DocumentIdField
				}

				if task.SourceResource.EsParam.IndexType != nil {
					esParamMap["index_type"] = task.SourceResource.EsParam.IndexType
				}

				if task.SourceResource.EsParam.DropCls != nil {
					dropClsMap := map[string]interface{}{}

					if task.SourceResource.EsParam.DropCls.DropInvalidMessageToCls != nil {
						dropClsMap["drop_invalid_message_to_cls"] = task.SourceResource.EsParam.DropCls.DropInvalidMessageToCls
					}

					if task.SourceResource.EsParam.DropCls.DropClsRegion != nil {
						dropClsMap["drop_cls_region"] = task.SourceResource.EsParam.DropCls.DropClsRegion
					}

					if task.SourceResource.EsParam.DropCls.DropClsOwneruin != nil {
						dropClsMap["drop_cls_owneruin"] = task.SourceResource.EsParam.DropCls.DropClsOwneruin
					}

					if task.SourceResource.EsParam.DropCls.DropClsTopicId != nil {
						dropClsMap["drop_cls_topic_id"] = task.SourceResource.EsParam.DropCls.DropClsTopicId
					}

					if task.SourceResource.EsParam.DropCls.DropClsLogSet != nil {
						dropClsMap["drop_cls_log_set"] = task.SourceResource.EsParam.DropCls.DropClsLogSet
					}

					esParamMap["drop_cls"] = []interface{}{dropClsMap}
				}

				if task.SourceResource.EsParam.DatabasePrimaryKey != nil {
					esParamMap["database_primary_key"] = task.SourceResource.EsParam.DatabasePrimaryKey
				}

				if task.SourceResource.EsParam.DropDlq != nil {
					dropDlqMap := map[string]interface{}{}

					if task.SourceResource.EsParam.DropDlq.Type != nil {
						dropDlqMap["type"] = task.SourceResource.EsParam.DropDlq.Type
					}

					if task.SourceResource.EsParam.DropDlq.KafkaParam != nil {
						kafkaParamMap := map[string]interface{}{}

						if task.SourceResource.EsParam.DropDlq.KafkaParam.SelfBuilt != nil {
							kafkaParamMap["self_built"] = task.SourceResource.EsParam.DropDlq.KafkaParam.SelfBuilt
						}

						if task.SourceResource.EsParam.DropDlq.KafkaParam.Resource != nil {
							kafkaParamMap["resource"] = task.SourceResource.EsParam.DropDlq.KafkaParam.Resource
						}

						if task.SourceResource.EsParam.DropDlq.KafkaParam.Topic != nil {
							kafkaParamMap["topic"] = task.SourceResource.EsParam.DropDlq.KafkaParam.Topic
						}

						if task.SourceResource.EsParam.DropDlq.KafkaParam.OffsetType != nil {
							kafkaParamMap["offset_type"] = task.SourceResource.EsParam.DropDlq.KafkaParam.OffsetType
						}

						if task.SourceResource.EsParam.DropDlq.KafkaParam.StartTime != nil {
							kafkaParamMap["start_time"] = task.SourceResource.EsParam.DropDlq.KafkaParam.StartTime
						}

						if task.SourceResource.EsParam.DropDlq.KafkaParam.ResourceName != nil {
							kafkaParamMap["resource_name"] = task.SourceResource.EsParam.DropDlq.KafkaParam.ResourceName
						}

						if task.SourceResource.EsParam.DropDlq.KafkaParam.ZoneId != nil {
							kafkaParamMap["zone_id"] = task.SourceResource.EsParam.DropDlq.KafkaParam.ZoneId
						}

						if task.SourceResource.EsParam.DropDlq.KafkaParam.TopicId != nil {
							kafkaParamMap["topic_id"] = task.SourceResource.EsParam.DropDlq.KafkaParam.TopicId
						}

						if task.SourceResource.EsParam.DropDlq.KafkaParam.PartitionNum != nil {
							kafkaParamMap["partition_num"] = task.SourceResource.EsParam.DropDlq.KafkaParam.PartitionNum
						}

						if task.SourceResource.EsParam.DropDlq.KafkaParam.EnableToleration != nil {
							kafkaParamMap["enable_toleration"] = task.SourceResource.EsParam.DropDlq.KafkaParam.EnableToleration
						}

						if task.SourceResource.EsParam.DropDlq.KafkaParam.QpsLimit != nil {
							kafkaParamMap["qps_limit"] = task.SourceResource.EsParam.DropDlq.KafkaParam.QpsLimit
						}

						if task.SourceResource.EsParam.DropDlq.KafkaParam.TableMappings != nil {
							tableMappingsList := []interface{}{}
							for _, tableMappings := range task.SourceResource.EsParam.DropDlq.KafkaParam.TableMappings {
								tableMappingsMap := map[string]interface{}{}

								if tableMappings.Database != nil {
									tableMappingsMap["database"] = tableMappings.Database
								}

								if tableMappings.Table != nil {
									tableMappingsMap["table"] = tableMappings.Table
								}

								if tableMappings.Topic != nil {
									tableMappingsMap["topic"] = tableMappings.Topic
								}

								if tableMappings.TopicId != nil {
									tableMappingsMap["topic_id"] = tableMappings.TopicId
								}

								tableMappingsList = append(tableMappingsList, tableMappingsMap)
							}

							kafkaParamMap["table_mappings"] = []interface{}{tableMappingsList}
						}

						if task.SourceResource.EsParam.DropDlq.KafkaParam.UseTableMapping != nil {
							kafkaParamMap["use_table_mapping"] = task.SourceResource.EsParam.DropDlq.KafkaParam.UseTableMapping
						}

						if task.SourceResource.EsParam.DropDlq.KafkaParam.UseAutoCreateTopic != nil {
							kafkaParamMap["use_auto_create_topic"] = task.SourceResource.EsParam.DropDlq.KafkaParam.UseAutoCreateTopic
						}

						if task.SourceResource.EsParam.DropDlq.KafkaParam.CompressionType != nil {
							kafkaParamMap["compression_type"] = task.SourceResource.EsParam.DropDlq.KafkaParam.CompressionType
						}

						if task.SourceResource.EsParam.DropDlq.KafkaParam.MsgMultiple != nil {
							kafkaParamMap["msg_multiple"] = task.SourceResource.EsParam.DropDlq.KafkaParam.MsgMultiple
						}

						if task.SourceResource.EsParam.DropDlq.KafkaParam.ConnectorSyncType != nil {
							kafkaParamMap["connector_sync_type"] = task.SourceResource.EsParam.DropDlq.KafkaParam.ConnectorSyncType
						}

						if task.SourceResource.EsParam.DropDlq.KafkaParam.KeepPartition != nil {
							kafkaParamMap["keep_partition"] = task.SourceResource.EsParam.DropDlq.KafkaParam.KeepPartition
						}

						dropDlqMap["kafka_param"] = []interface{}{kafkaParamMap}
					}

					if task.SourceResource.EsParam.DropDlq.RetryInterval != nil {
						dropDlqMap["retry_interval"] = task.SourceResource.EsParam.DropDlq.RetryInterval
					}

					if task.SourceResource.EsParam.DropDlq.MaxRetryAttempts != nil {
						dropDlqMap["max_retry_attempts"] = task.SourceResource.EsParam.DropDlq.MaxRetryAttempts
					}

					if task.SourceResource.EsParam.DropDlq.TopicParam != nil {
						topicParamMap := map[string]interface{}{}

						if task.SourceResource.EsParam.DropDlq.TopicParam.Resource != nil {
							topicParamMap["resource"] = task.SourceResource.EsParam.DropDlq.TopicParam.Resource
						}

						if task.SourceResource.EsParam.DropDlq.TopicParam.OffsetType != nil {
							topicParamMap["offset_type"] = task.SourceResource.EsParam.DropDlq.TopicParam.OffsetType
						}

						if task.SourceResource.EsParam.DropDlq.TopicParam.StartTime != nil {
							topicParamMap["start_time"] = task.SourceResource.EsParam.DropDlq.TopicParam.StartTime
						}

						if task.SourceResource.EsParam.DropDlq.TopicParam.TopicId != nil {
							topicParamMap["topic_id"] = task.SourceResource.EsParam.DropDlq.TopicParam.TopicId
						}

						if task.SourceResource.EsParam.DropDlq.TopicParam.CompressionType != nil {
							topicParamMap["compression_type"] = task.SourceResource.EsParam.DropDlq.TopicParam.CompressionType
						}

						if task.SourceResource.EsParam.DropDlq.TopicParam.UseAutoCreateTopic != nil {
							topicParamMap["use_auto_create_topic"] = task.SourceResource.EsParam.DropDlq.TopicParam.UseAutoCreateTopic
						}

						if task.SourceResource.EsParam.DropDlq.TopicParam.MsgMultiple != nil {
							topicParamMap["msg_multiple"] = task.SourceResource.EsParam.DropDlq.TopicParam.MsgMultiple
						}

						dropDlqMap["topic_param"] = []interface{}{topicParamMap}
					}

					if task.SourceResource.EsParam.DropDlq.DlqType != nil {
						dropDlqMap["dlq_type"] = task.SourceResource.EsParam.DropDlq.DlqType
					}

					esParamMap["drop_dlq"] = []interface{}{dropDlqMap}
				}

				sourceResourceMap["es_param"] = []interface{}{esParamMap}
			}

			if task.SourceResource.TdwParam != nil {
				tdwParamMap := map[string]interface{}{}

				if task.SourceResource.TdwParam.Bid != nil {
					tdwParamMap["bid"] = task.SourceResource.TdwParam.Bid
				}

				if task.SourceResource.TdwParam.Tid != nil {
					tdwParamMap["tid"] = task.SourceResource.TdwParam.Tid
				}

				if task.SourceResource.TdwParam.IsDomestic != nil {
					tdwParamMap["is_domestic"] = task.SourceResource.TdwParam.IsDomestic
				}

				if task.SourceResource.TdwParam.TdwHost != nil {
					tdwParamMap["tdw_host"] = task.SourceResource.TdwParam.TdwHost
				}

				if task.SourceResource.TdwParam.TdwPort != nil {
					tdwParamMap["tdw_port"] = task.SourceResource.TdwParam.TdwPort
				}

				sourceResourceMap["tdw_param"] = []interface{}{tdwParamMap}
			}

			if task.SourceResource.DtsParam != nil {
				dtsParamMap := map[string]interface{}{}

				if task.SourceResource.DtsParam.Resource != nil {
					dtsParamMap["resource"] = task.SourceResource.DtsParam.Resource
				}

				if task.SourceResource.DtsParam.Ip != nil {
					dtsParamMap["ip"] = task.SourceResource.DtsParam.Ip
				}

				if task.SourceResource.DtsParam.Port != nil {
					dtsParamMap["port"] = task.SourceResource.DtsParam.Port
				}

				if task.SourceResource.DtsParam.Topic != nil {
					dtsParamMap["topic"] = task.SourceResource.DtsParam.Topic
				}

				if task.SourceResource.DtsParam.GroupId != nil {
					dtsParamMap["group_id"] = task.SourceResource.DtsParam.GroupId
				}

				if task.SourceResource.DtsParam.GroupUser != nil {
					dtsParamMap["group_user"] = task.SourceResource.DtsParam.GroupUser
				}

				if task.SourceResource.DtsParam.GroupPassword != nil {
					dtsParamMap["group_password"] = task.SourceResource.DtsParam.GroupPassword
				}

				if task.SourceResource.DtsParam.TranSql != nil {
					dtsParamMap["tran_sql"] = task.SourceResource.DtsParam.TranSql
				}

				sourceResourceMap["dts_param"] = []interface{}{dtsParamMap}
			}

			if task.SourceResource.ClickHouseParam != nil {
				clickHouseParamMap := map[string]interface{}{}

				if task.SourceResource.ClickHouseParam.Cluster != nil {
					clickHouseParamMap["cluster"] = task.SourceResource.ClickHouseParam.Cluster
				}

				if task.SourceResource.ClickHouseParam.Database != nil {
					clickHouseParamMap["database"] = task.SourceResource.ClickHouseParam.Database
				}

				if task.SourceResource.ClickHouseParam.Table != nil {
					clickHouseParamMap["table"] = task.SourceResource.ClickHouseParam.Table
				}

				if task.SourceResource.ClickHouseParam.Schema != nil {
					schemaList := []interface{}{}
					for _, schema := range task.SourceResource.ClickHouseParam.Schema {
						schemaMap := map[string]interface{}{}

						if schema.ColumnName != nil {
							schemaMap["column_name"] = schema.ColumnName
						}

						if schema.JsonKey != nil {
							schemaMap["json_key"] = schema.JsonKey
						}

						if schema.Type != nil {
							schemaMap["type"] = schema.Type
						}

						if schema.AllowNull != nil {
							schemaMap["allow_null"] = schema.AllowNull
						}

						schemaList = append(schemaList, schemaMap)
					}

					clickHouseParamMap["schema"] = []interface{}{schemaList}
				}

				if task.SourceResource.ClickHouseParam.Resource != nil {
					clickHouseParamMap["resource"] = task.SourceResource.ClickHouseParam.Resource
				}

				if task.SourceResource.ClickHouseParam.Ip != nil {
					clickHouseParamMap["ip"] = task.SourceResource.ClickHouseParam.Ip
				}

				if task.SourceResource.ClickHouseParam.Port != nil {
					clickHouseParamMap["port"] = task.SourceResource.ClickHouseParam.Port
				}

				if task.SourceResource.ClickHouseParam.UserName != nil {
					clickHouseParamMap["user_name"] = task.SourceResource.ClickHouseParam.UserName
				}

				if task.SourceResource.ClickHouseParam.Password != nil {
					clickHouseParamMap["password"] = task.SourceResource.ClickHouseParam.Password
				}

				if task.SourceResource.ClickHouseParam.ServiceVip != nil {
					clickHouseParamMap["service_vip"] = task.SourceResource.ClickHouseParam.ServiceVip
				}

				if task.SourceResource.ClickHouseParam.UniqVpcId != nil {
					clickHouseParamMap["uniq_vpc_id"] = task.SourceResource.ClickHouseParam.UniqVpcId
				}

				if task.SourceResource.ClickHouseParam.SelfBuilt != nil {
					clickHouseParamMap["self_built"] = task.SourceResource.ClickHouseParam.SelfBuilt
				}

				if task.SourceResource.ClickHouseParam.DropInvalidMessage != nil {
					clickHouseParamMap["drop_invalid_message"] = task.SourceResource.ClickHouseParam.DropInvalidMessage
				}

				if task.SourceResource.ClickHouseParam.Type != nil {
					clickHouseParamMap["type"] = task.SourceResource.ClickHouseParam.Type
				}

				if task.SourceResource.ClickHouseParam.DropCls != nil {
					dropClsMap := map[string]interface{}{}

					if task.SourceResource.ClickHouseParam.DropCls.DropInvalidMessageToCls != nil {
						dropClsMap["drop_invalid_message_to_cls"] = task.SourceResource.ClickHouseParam.DropCls.DropInvalidMessageToCls
					}

					if task.SourceResource.ClickHouseParam.DropCls.DropClsRegion != nil {
						dropClsMap["drop_cls_region"] = task.SourceResource.ClickHouseParam.DropCls.DropClsRegion
					}

					if task.SourceResource.ClickHouseParam.DropCls.DropClsOwneruin != nil {
						dropClsMap["drop_cls_owneruin"] = task.SourceResource.ClickHouseParam.DropCls.DropClsOwneruin
					}

					if task.SourceResource.ClickHouseParam.DropCls.DropClsTopicId != nil {
						dropClsMap["drop_cls_topic_id"] = task.SourceResource.ClickHouseParam.DropCls.DropClsTopicId
					}

					if task.SourceResource.ClickHouseParam.DropCls.DropClsLogSet != nil {
						dropClsMap["drop_cls_log_set"] = task.SourceResource.ClickHouseParam.DropCls.DropClsLogSet
					}

					clickHouseParamMap["drop_cls"] = []interface{}{dropClsMap}
				}

				sourceResourceMap["click_house_param"] = []interface{}{clickHouseParamMap}
			}

			if task.SourceResource.ClsParam != nil {
				clsParamMap := map[string]interface{}{}

				if task.SourceResource.ClsParam.DecodeJson != nil {
					clsParamMap["decode_json"] = task.SourceResource.ClsParam.DecodeJson
				}

				if task.SourceResource.ClsParam.Resource != nil {
					clsParamMap["resource"] = task.SourceResource.ClsParam.Resource
				}

				if task.SourceResource.ClsParam.LogSet != nil {
					clsParamMap["log_set"] = task.SourceResource.ClsParam.LogSet
				}

				if task.SourceResource.ClsParam.ContentKey != nil {
					clsParamMap["content_key"] = task.SourceResource.ClsParam.ContentKey
				}

				if task.SourceResource.ClsParam.TimeField != nil {
					clsParamMap["time_field"] = task.SourceResource.ClsParam.TimeField
				}

				sourceResourceMap["cls_param"] = []interface{}{clsParamMap}
			}

			if task.SourceResource.CosParam != nil {
				cosParamMap := map[string]interface{}{}

				if task.SourceResource.CosParam.BucketName != nil {
					cosParamMap["bucket_name"] = task.SourceResource.CosParam.BucketName
				}

				if task.SourceResource.CosParam.Region != nil {
					cosParamMap["region"] = task.SourceResource.CosParam.Region
				}

				if task.SourceResource.CosParam.ObjectKey != nil {
					cosParamMap["object_key"] = task.SourceResource.CosParam.ObjectKey
				}

				if task.SourceResource.CosParam.AggregateBatchSize != nil {
					cosParamMap["aggregate_batch_size"] = task.SourceResource.CosParam.AggregateBatchSize
				}

				if task.SourceResource.CosParam.AggregateInterval != nil {
					cosParamMap["aggregate_interval"] = task.SourceResource.CosParam.AggregateInterval
				}

				if task.SourceResource.CosParam.FormatOutputType != nil {
					cosParamMap["format_output_type"] = task.SourceResource.CosParam.FormatOutputType
				}

				if task.SourceResource.CosParam.ObjectKeyPrefix != nil {
					cosParamMap["object_key_prefix"] = task.SourceResource.CosParam.ObjectKeyPrefix
				}

				if task.SourceResource.CosParam.DirectoryTimeFormat != nil {
					cosParamMap["directory_time_format"] = task.SourceResource.CosParam.DirectoryTimeFormat
				}

				sourceResourceMap["cos_param"] = []interface{}{cosParamMap}
			}

			if task.SourceResource.MySQLParam != nil {
				mySQLParamMap := map[string]interface{}{}

				if task.SourceResource.MySQLParam.Database != nil {
					mySQLParamMap["database"] = task.SourceResource.MySQLParam.Database
				}

				if task.SourceResource.MySQLParam.Table != nil {
					mySQLParamMap["table"] = task.SourceResource.MySQLParam.Table
				}

				if task.SourceResource.MySQLParam.Resource != nil {
					mySQLParamMap["resource"] = task.SourceResource.MySQLParam.Resource
				}

				if task.SourceResource.MySQLParam.SnapshotMode != nil {
					mySQLParamMap["snapshot_mode"] = task.SourceResource.MySQLParam.SnapshotMode
				}

				if task.SourceResource.MySQLParam.DdlTopic != nil {
					mySQLParamMap["ddl_topic"] = task.SourceResource.MySQLParam.DdlTopic
				}

				if task.SourceResource.MySQLParam.DataSourceMonitorMode != nil {
					mySQLParamMap["data_source_monitor_mode"] = task.SourceResource.MySQLParam.DataSourceMonitorMode
				}

				if task.SourceResource.MySQLParam.DataSourceMonitorResource != nil {
					mySQLParamMap["data_source_monitor_resource"] = task.SourceResource.MySQLParam.DataSourceMonitorResource
				}

				if task.SourceResource.MySQLParam.DataSourceIncrementMode != nil {
					mySQLParamMap["data_source_increment_mode"] = task.SourceResource.MySQLParam.DataSourceIncrementMode
				}

				if task.SourceResource.MySQLParam.DataSourceIncrementColumn != nil {
					mySQLParamMap["data_source_increment_column"] = task.SourceResource.MySQLParam.DataSourceIncrementColumn
				}

				if task.SourceResource.MySQLParam.DataSourceStartFrom != nil {
					mySQLParamMap["data_source_start_from"] = task.SourceResource.MySQLParam.DataSourceStartFrom
				}

				if task.SourceResource.MySQLParam.DataTargetInsertMode != nil {
					mySQLParamMap["data_target_insert_mode"] = task.SourceResource.MySQLParam.DataTargetInsertMode
				}

				if task.SourceResource.MySQLParam.DataTargetPrimaryKeyField != nil {
					mySQLParamMap["data_target_primary_key_field"] = task.SourceResource.MySQLParam.DataTargetPrimaryKeyField
				}

				if task.SourceResource.MySQLParam.DataTargetRecordMapping != nil {
					dataTargetRecordMappingList := []interface{}{}
					for _, dataTargetRecordMapping := range task.SourceResource.MySQLParam.DataTargetRecordMapping {
						dataTargetRecordMappingMap := map[string]interface{}{}

						if dataTargetRecordMapping.JsonKey != nil {
							dataTargetRecordMappingMap["json_key"] = dataTargetRecordMapping.JsonKey
						}

						if dataTargetRecordMapping.Type != nil {
							dataTargetRecordMappingMap["type"] = dataTargetRecordMapping.Type
						}

						if dataTargetRecordMapping.AllowNull != nil {
							dataTargetRecordMappingMap["allow_null"] = dataTargetRecordMapping.AllowNull
						}

						if dataTargetRecordMapping.ColumnName != nil {
							dataTargetRecordMappingMap["column_name"] = dataTargetRecordMapping.ColumnName
						}

						if dataTargetRecordMapping.ExtraInfo != nil {
							dataTargetRecordMappingMap["extra_info"] = dataTargetRecordMapping.ExtraInfo
						}

						if dataTargetRecordMapping.ColumnSize != nil {
							dataTargetRecordMappingMap["column_size"] = dataTargetRecordMapping.ColumnSize
						}

						if dataTargetRecordMapping.DecimalDigits != nil {
							dataTargetRecordMappingMap["decimal_digits"] = dataTargetRecordMapping.DecimalDigits
						}

						if dataTargetRecordMapping.AutoIncrement != nil {
							dataTargetRecordMappingMap["auto_increment"] = dataTargetRecordMapping.AutoIncrement
						}

						if dataTargetRecordMapping.DefaultValue != nil {
							dataTargetRecordMappingMap["default_value"] = dataTargetRecordMapping.DefaultValue
						}

						dataTargetRecordMappingList = append(dataTargetRecordMappingList, dataTargetRecordMappingMap)
					}

					mySQLParamMap["data_target_record_mapping"] = []interface{}{dataTargetRecordMappingList}
				}

				if task.SourceResource.MySQLParam.TopicRegex != nil {
					mySQLParamMap["topic_regex"] = task.SourceResource.MySQLParam.TopicRegex
				}

				if task.SourceResource.MySQLParam.TopicReplacement != nil {
					mySQLParamMap["topic_replacement"] = task.SourceResource.MySQLParam.TopicReplacement
				}

				if task.SourceResource.MySQLParam.KeyColumns != nil {
					mySQLParamMap["key_columns"] = task.SourceResource.MySQLParam.KeyColumns
				}

				if task.SourceResource.MySQLParam.DropInvalidMessage != nil {
					mySQLParamMap["drop_invalid_message"] = task.SourceResource.MySQLParam.DropInvalidMessage
				}

				if task.SourceResource.MySQLParam.DropCls != nil {
					dropClsMap := map[string]interface{}{}

					if task.SourceResource.MySQLParam.DropCls.DropInvalidMessageToCls != nil {
						dropClsMap["drop_invalid_message_to_cls"] = task.SourceResource.MySQLParam.DropCls.DropInvalidMessageToCls
					}

					if task.SourceResource.MySQLParam.DropCls.DropClsRegion != nil {
						dropClsMap["drop_cls_region"] = task.SourceResource.MySQLParam.DropCls.DropClsRegion
					}

					if task.SourceResource.MySQLParam.DropCls.DropClsOwneruin != nil {
						dropClsMap["drop_cls_owneruin"] = task.SourceResource.MySQLParam.DropCls.DropClsOwneruin
					}

					if task.SourceResource.MySQLParam.DropCls.DropClsTopicId != nil {
						dropClsMap["drop_cls_topic_id"] = task.SourceResource.MySQLParam.DropCls.DropClsTopicId
					}

					if task.SourceResource.MySQLParam.DropCls.DropClsLogSet != nil {
						dropClsMap["drop_cls_log_set"] = task.SourceResource.MySQLParam.DropCls.DropClsLogSet
					}

					mySQLParamMap["drop_cls"] = []interface{}{dropClsMap}
				}

				if task.SourceResource.MySQLParam.OutputFormat != nil {
					mySQLParamMap["output_format"] = task.SourceResource.MySQLParam.OutputFormat
				}

				if task.SourceResource.MySQLParam.IsTablePrefix != nil {
					mySQLParamMap["is_table_prefix"] = task.SourceResource.MySQLParam.IsTablePrefix
				}

				if task.SourceResource.MySQLParam.IncludeContentChanges != nil {
					mySQLParamMap["include_content_changes"] = task.SourceResource.MySQLParam.IncludeContentChanges
				}

				if task.SourceResource.MySQLParam.IncludeQuery != nil {
					mySQLParamMap["include_query"] = task.SourceResource.MySQLParam.IncludeQuery
				}

				if task.SourceResource.MySQLParam.RecordWithSchema != nil {
					mySQLParamMap["record_with_schema"] = task.SourceResource.MySQLParam.RecordWithSchema
				}

				if task.SourceResource.MySQLParam.SignalDatabase != nil {
					mySQLParamMap["signal_database"] = task.SourceResource.MySQLParam.SignalDatabase
				}

				if task.SourceResource.MySQLParam.IsTableRegular != nil {
					mySQLParamMap["is_table_regular"] = task.SourceResource.MySQLParam.IsTableRegular
				}

				sourceResourceMap["my_sql_param"] = []interface{}{mySQLParamMap}
			}

			if task.SourceResource.PostgreSQLParam != nil {
				postgreSQLParamMap := map[string]interface{}{}

				if task.SourceResource.PostgreSQLParam.Database != nil {
					postgreSQLParamMap["database"] = task.SourceResource.PostgreSQLParam.Database
				}

				if task.SourceResource.PostgreSQLParam.Table != nil {
					postgreSQLParamMap["table"] = task.SourceResource.PostgreSQLParam.Table
				}

				if task.SourceResource.PostgreSQLParam.Resource != nil {
					postgreSQLParamMap["resource"] = task.SourceResource.PostgreSQLParam.Resource
				}

				if task.SourceResource.PostgreSQLParam.PluginName != nil {
					postgreSQLParamMap["plugin_name"] = task.SourceResource.PostgreSQLParam.PluginName
				}

				if task.SourceResource.PostgreSQLParam.SnapshotMode != nil {
					postgreSQLParamMap["snapshot_mode"] = task.SourceResource.PostgreSQLParam.SnapshotMode
				}

				if task.SourceResource.PostgreSQLParam.DataFormat != nil {
					postgreSQLParamMap["data_format"] = task.SourceResource.PostgreSQLParam.DataFormat
				}

				if task.SourceResource.PostgreSQLParam.DataTargetInsertMode != nil {
					postgreSQLParamMap["data_target_insert_mode"] = task.SourceResource.PostgreSQLParam.DataTargetInsertMode
				}

				if task.SourceResource.PostgreSQLParam.DataTargetPrimaryKeyField != nil {
					postgreSQLParamMap["data_target_primary_key_field"] = task.SourceResource.PostgreSQLParam.DataTargetPrimaryKeyField
				}

				if task.SourceResource.PostgreSQLParam.DataTargetRecordMapping != nil {
					dataTargetRecordMappingList := []interface{}{}
					for _, dataTargetRecordMapping := range task.SourceResource.PostgreSQLParam.DataTargetRecordMapping {
						dataTargetRecordMappingMap := map[string]interface{}{}

						if dataTargetRecordMapping.JsonKey != nil {
							dataTargetRecordMappingMap["json_key"] = dataTargetRecordMapping.JsonKey
						}

						if dataTargetRecordMapping.Type != nil {
							dataTargetRecordMappingMap["type"] = dataTargetRecordMapping.Type
						}

						if dataTargetRecordMapping.AllowNull != nil {
							dataTargetRecordMappingMap["allow_null"] = dataTargetRecordMapping.AllowNull
						}

						if dataTargetRecordMapping.ColumnName != nil {
							dataTargetRecordMappingMap["column_name"] = dataTargetRecordMapping.ColumnName
						}

						if dataTargetRecordMapping.ExtraInfo != nil {
							dataTargetRecordMappingMap["extra_info"] = dataTargetRecordMapping.ExtraInfo
						}

						if dataTargetRecordMapping.ColumnSize != nil {
							dataTargetRecordMappingMap["column_size"] = dataTargetRecordMapping.ColumnSize
						}

						if dataTargetRecordMapping.DecimalDigits != nil {
							dataTargetRecordMappingMap["decimal_digits"] = dataTargetRecordMapping.DecimalDigits
						}

						if dataTargetRecordMapping.AutoIncrement != nil {
							dataTargetRecordMappingMap["auto_increment"] = dataTargetRecordMapping.AutoIncrement
						}

						if dataTargetRecordMapping.DefaultValue != nil {
							dataTargetRecordMappingMap["default_value"] = dataTargetRecordMapping.DefaultValue
						}

						dataTargetRecordMappingList = append(dataTargetRecordMappingList, dataTargetRecordMappingMap)
					}

					postgreSQLParamMap["data_target_record_mapping"] = []interface{}{dataTargetRecordMappingList}
				}

				if task.SourceResource.PostgreSQLParam.DropInvalidMessage != nil {
					postgreSQLParamMap["drop_invalid_message"] = task.SourceResource.PostgreSQLParam.DropInvalidMessage
				}

				if task.SourceResource.PostgreSQLParam.IsTableRegular != nil {
					postgreSQLParamMap["is_table_regular"] = task.SourceResource.PostgreSQLParam.IsTableRegular
				}

				if task.SourceResource.PostgreSQLParam.KeyColumns != nil {
					postgreSQLParamMap["key_columns"] = task.SourceResource.PostgreSQLParam.KeyColumns
				}

				if task.SourceResource.PostgreSQLParam.RecordWithSchema != nil {
					postgreSQLParamMap["record_with_schema"] = task.SourceResource.PostgreSQLParam.RecordWithSchema
				}

				sourceResourceMap["postgre_sql_param"] = []interface{}{postgreSQLParamMap}
			}

			if task.SourceResource.TopicParam != nil {
				topicParamMap := map[string]interface{}{}

				if task.SourceResource.TopicParam.Resource != nil {
					topicParamMap["resource"] = task.SourceResource.TopicParam.Resource
				}

				if task.SourceResource.TopicParam.OffsetType != nil {
					topicParamMap["offset_type"] = task.SourceResource.TopicParam.OffsetType
				}

				if task.SourceResource.TopicParam.StartTime != nil {
					topicParamMap["start_time"] = task.SourceResource.TopicParam.StartTime
				}

				if task.SourceResource.TopicParam.TopicId != nil {
					topicParamMap["topic_id"] = task.SourceResource.TopicParam.TopicId
				}

				if task.SourceResource.TopicParam.CompressionType != nil {
					topicParamMap["compression_type"] = task.SourceResource.TopicParam.CompressionType
				}

				if task.SourceResource.TopicParam.UseAutoCreateTopic != nil {
					topicParamMap["use_auto_create_topic"] = task.SourceResource.TopicParam.UseAutoCreateTopic
				}

				if task.SourceResource.TopicParam.MsgMultiple != nil {
					topicParamMap["msg_multiple"] = task.SourceResource.TopicParam.MsgMultiple
				}

				sourceResourceMap["topic_param"] = []interface{}{topicParamMap}
			}

			if task.SourceResource.MariaDBParam != nil {
				mariaDBParamMap := map[string]interface{}{}

				if task.SourceResource.MariaDBParam.Database != nil {
					mariaDBParamMap["database"] = task.SourceResource.MariaDBParam.Database
				}

				if task.SourceResource.MariaDBParam.Table != nil {
					mariaDBParamMap["table"] = task.SourceResource.MariaDBParam.Table
				}

				if task.SourceResource.MariaDBParam.Resource != nil {
					mariaDBParamMap["resource"] = task.SourceResource.MariaDBParam.Resource
				}

				if task.SourceResource.MariaDBParam.SnapshotMode != nil {
					mariaDBParamMap["snapshot_mode"] = task.SourceResource.MariaDBParam.SnapshotMode
				}

				if task.SourceResource.MariaDBParam.KeyColumns != nil {
					mariaDBParamMap["key_columns"] = task.SourceResource.MariaDBParam.KeyColumns
				}

				if task.SourceResource.MariaDBParam.IsTablePrefix != nil {
					mariaDBParamMap["is_table_prefix"] = task.SourceResource.MariaDBParam.IsTablePrefix
				}

				if task.SourceResource.MariaDBParam.OutputFormat != nil {
					mariaDBParamMap["output_format"] = task.SourceResource.MariaDBParam.OutputFormat
				}

				if task.SourceResource.MariaDBParam.IncludeContentChanges != nil {
					mariaDBParamMap["include_content_changes"] = task.SourceResource.MariaDBParam.IncludeContentChanges
				}

				if task.SourceResource.MariaDBParam.IncludeQuery != nil {
					mariaDBParamMap["include_query"] = task.SourceResource.MariaDBParam.IncludeQuery
				}

				if task.SourceResource.MariaDBParam.RecordWithSchema != nil {
					mariaDBParamMap["record_with_schema"] = task.SourceResource.MariaDBParam.RecordWithSchema
				}

				sourceResourceMap["maria_db_param"] = []interface{}{mariaDBParamMap}
			}

			if task.SourceResource.SQLServerParam != nil {
				sQLServerParamMap := map[string]interface{}{}

				if task.SourceResource.SQLServerParam.Database != nil {
					sQLServerParamMap["database"] = task.SourceResource.SQLServerParam.Database
				}

				if task.SourceResource.SQLServerParam.Table != nil {
					sQLServerParamMap["table"] = task.SourceResource.SQLServerParam.Table
				}

				if task.SourceResource.SQLServerParam.Resource != nil {
					sQLServerParamMap["resource"] = task.SourceResource.SQLServerParam.Resource
				}

				if task.SourceResource.SQLServerParam.SnapshotMode != nil {
					sQLServerParamMap["snapshot_mode"] = task.SourceResource.SQLServerParam.SnapshotMode
				}

				sourceResourceMap["sql_server_param"] = []interface{}{sQLServerParamMap}
			}

			if task.SourceResource.CtsdbParam != nil {
				ctsdbParamMap := map[string]interface{}{}

				if task.SourceResource.CtsdbParam.Resource != nil {
					ctsdbParamMap["resource"] = task.SourceResource.CtsdbParam.Resource
				}

				if task.SourceResource.CtsdbParam.CtsdbMetric != nil {
					ctsdbParamMap["ctsdb_metric"] = task.SourceResource.CtsdbParam.CtsdbMetric
				}

				sourceResourceMap["ctsdb_param"] = []interface{}{ctsdbParamMap}
			}

			if task.SourceResource.ScfParam != nil {
				scfParamMap := map[string]interface{}{}

				if task.SourceResource.ScfParam.FunctionName != nil {
					scfParamMap["function_name"] = task.SourceResource.ScfParam.FunctionName
				}

				if task.SourceResource.ScfParam.Namespace != nil {
					scfParamMap["namespace"] = task.SourceResource.ScfParam.Namespace
				}

				if task.SourceResource.ScfParam.Qualifier != nil {
					scfParamMap["qualifier"] = task.SourceResource.ScfParam.Qualifier
				}

				if task.SourceResource.ScfParam.BatchSize != nil {
					scfParamMap["batch_size"] = task.SourceResource.ScfParam.BatchSize
				}

				if task.SourceResource.ScfParam.MaxRetries != nil {
					scfParamMap["max_retries"] = task.SourceResource.ScfParam.MaxRetries
				}

				sourceResourceMap["scf_param"] = []interface{}{scfParamMap}
			}

			taskMap["source_resource"] = []interface{}{sourceResourceMap}
		}

		if task.TargetResource != nil {
			targetResourceMap := map[string]interface{}{}

			if task.TargetResource.Type != nil {
				targetResourceMap["type"] = task.TargetResource.Type
			}

			if task.TargetResource.KafkaParam != nil {
				kafkaParamMap := map[string]interface{}{}

				if task.TargetResource.KafkaParam.SelfBuilt != nil {
					kafkaParamMap["self_built"] = task.TargetResource.KafkaParam.SelfBuilt
				}

				if task.TargetResource.KafkaParam.Resource != nil {
					kafkaParamMap["resource"] = task.TargetResource.KafkaParam.Resource
				}

				if task.TargetResource.KafkaParam.Topic != nil {
					kafkaParamMap["topic"] = task.TargetResource.KafkaParam.Topic
				}

				if task.TargetResource.KafkaParam.OffsetType != nil {
					kafkaParamMap["offset_type"] = task.TargetResource.KafkaParam.OffsetType
				}

				if task.TargetResource.KafkaParam.StartTime != nil {
					kafkaParamMap["start_time"] = task.TargetResource.KafkaParam.StartTime
				}

				if task.TargetResource.KafkaParam.ResourceName != nil {
					kafkaParamMap["resource_name"] = task.TargetResource.KafkaParam.ResourceName
				}

				if task.TargetResource.KafkaParam.ZoneId != nil {
					kafkaParamMap["zone_id"] = task.TargetResource.KafkaParam.ZoneId
				}

				if task.TargetResource.KafkaParam.TopicId != nil {
					kafkaParamMap["topic_id"] = task.TargetResource.KafkaParam.TopicId
				}

				if task.TargetResource.KafkaParam.PartitionNum != nil {
					kafkaParamMap["partition_num"] = task.TargetResource.KafkaParam.PartitionNum
				}

				if task.TargetResource.KafkaParam.EnableToleration != nil {
					kafkaParamMap["enable_toleration"] = task.TargetResource.KafkaParam.EnableToleration
				}

				if task.TargetResource.KafkaParam.QpsLimit != nil {
					kafkaParamMap["qps_limit"] = task.TargetResource.KafkaParam.QpsLimit
				}

				if task.TargetResource.KafkaParam.TableMappings != nil {
					tableMappingsList := []interface{}{}
					for _, tableMappings := range task.TargetResource.KafkaParam.TableMappings {
						tableMappingsMap := map[string]interface{}{}

						if tableMappings.Database != nil {
							tableMappingsMap["database"] = tableMappings.Database
						}

						if tableMappings.Table != nil {
							tableMappingsMap["table"] = tableMappings.Table
						}

						if tableMappings.Topic != nil {
							tableMappingsMap["topic"] = tableMappings.Topic
						}

						if tableMappings.TopicId != nil {
							tableMappingsMap["topic_id"] = tableMappings.TopicId
						}

						tableMappingsList = append(tableMappingsList, tableMappingsMap)
					}

					kafkaParamMap["table_mappings"] = []interface{}{tableMappingsList}
				}

				if task.TargetResource.KafkaParam.UseTableMapping != nil {
					kafkaParamMap["use_table_mapping"] = task.TargetResource.KafkaParam.UseTableMapping
				}

				if task.TargetResource.KafkaParam.UseAutoCreateTopic != nil {
					kafkaParamMap["use_auto_create_topic"] = task.TargetResource.KafkaParam.UseAutoCreateTopic
				}

				if task.TargetResource.KafkaParam.CompressionType != nil {
					kafkaParamMap["compression_type"] = task.TargetResource.KafkaParam.CompressionType
				}

				if task.TargetResource.KafkaParam.MsgMultiple != nil {
					kafkaParamMap["msg_multiple"] = task.TargetResource.KafkaParam.MsgMultiple
				}

				if task.TargetResource.KafkaParam.ConnectorSyncType != nil {
					kafkaParamMap["connector_sync_type"] = task.TargetResource.KafkaParam.ConnectorSyncType
				}

				if task.TargetResource.KafkaParam.KeepPartition != nil {
					kafkaParamMap["keep_partition"] = task.TargetResource.KafkaParam.KeepPartition
				}

				targetResourceMap["kafka_param"] = []interface{}{kafkaParamMap}
			}

			if task.TargetResource.EventBusParam != nil {
				eventBusParamMap := map[string]interface{}{}

				if task.TargetResource.EventBusParam.Type != nil {
					eventBusParamMap["type"] = task.TargetResource.EventBusParam.Type
				}

				if task.TargetResource.EventBusParam.SelfBuilt != nil {
					eventBusParamMap["self_built"] = task.TargetResource.EventBusParam.SelfBuilt
				}

				if task.TargetResource.EventBusParam.Resource != nil {
					eventBusParamMap["resource"] = task.TargetResource.EventBusParam.Resource
				}

				if task.TargetResource.EventBusParam.Namespace != nil {
					eventBusParamMap["namespace"] = task.TargetResource.EventBusParam.Namespace
				}

				if task.TargetResource.EventBusParam.FunctionName != nil {
					eventBusParamMap["function_name"] = task.TargetResource.EventBusParam.FunctionName
				}

				if task.TargetResource.EventBusParam.Qualifier != nil {
					eventBusParamMap["qualifier"] = task.TargetResource.EventBusParam.Qualifier
				}

				targetResourceMap["event_bus_param"] = []interface{}{eventBusParamMap}
			}

			if task.TargetResource.MongoDBParam != nil {
				mongoDBParamMap := map[string]interface{}{}

				if task.TargetResource.MongoDBParam.Database != nil {
					mongoDBParamMap["database"] = task.TargetResource.MongoDBParam.Database
				}

				if task.TargetResource.MongoDBParam.Collection != nil {
					mongoDBParamMap["collection"] = task.TargetResource.MongoDBParam.Collection
				}

				if task.TargetResource.MongoDBParam.CopyExisting != nil {
					mongoDBParamMap["copy_existing"] = task.TargetResource.MongoDBParam.CopyExisting
				}

				if task.TargetResource.MongoDBParam.Resource != nil {
					mongoDBParamMap["resource"] = task.TargetResource.MongoDBParam.Resource
				}

				if task.TargetResource.MongoDBParam.Ip != nil {
					mongoDBParamMap["ip"] = task.TargetResource.MongoDBParam.Ip
				}

				if task.TargetResource.MongoDBParam.Port != nil {
					mongoDBParamMap["port"] = task.TargetResource.MongoDBParam.Port
				}

				if task.TargetResource.MongoDBParam.UserName != nil {
					mongoDBParamMap["user_name"] = task.TargetResource.MongoDBParam.UserName
				}

				if task.TargetResource.MongoDBParam.Password != nil {
					mongoDBParamMap["password"] = task.TargetResource.MongoDBParam.Password
				}

				if task.TargetResource.MongoDBParam.ListeningEvent != nil {
					mongoDBParamMap["listening_event"] = task.TargetResource.MongoDBParam.ListeningEvent
				}

				if task.TargetResource.MongoDBParam.ReadPreference != nil {
					mongoDBParamMap["read_preference"] = task.TargetResource.MongoDBParam.ReadPreference
				}

				if task.TargetResource.MongoDBParam.Pipeline != nil {
					mongoDBParamMap["pipeline"] = task.TargetResource.MongoDBParam.Pipeline
				}

				if task.TargetResource.MongoDBParam.SelfBuilt != nil {
					mongoDBParamMap["self_built"] = task.TargetResource.MongoDBParam.SelfBuilt
				}

				targetResourceMap["mongo_db_param"] = []interface{}{mongoDBParamMap}
			}

			if task.TargetResource.EsParam != nil {
				esParamMap := map[string]interface{}{}

				if task.TargetResource.EsParam.Resource != nil {
					esParamMap["resource"] = task.TargetResource.EsParam.Resource
				}

				if task.TargetResource.EsParam.Port != nil {
					esParamMap["port"] = task.TargetResource.EsParam.Port
				}

				if task.TargetResource.EsParam.UserName != nil {
					esParamMap["user_name"] = task.TargetResource.EsParam.UserName
				}

				if task.TargetResource.EsParam.Password != nil {
					esParamMap["password"] = task.TargetResource.EsParam.Password
				}

				if task.TargetResource.EsParam.SelfBuilt != nil {
					esParamMap["self_built"] = task.TargetResource.EsParam.SelfBuilt
				}

				if task.TargetResource.EsParam.ServiceVip != nil {
					esParamMap["service_vip"] = task.TargetResource.EsParam.ServiceVip
				}

				if task.TargetResource.EsParam.UniqVpcId != nil {
					esParamMap["uniq_vpc_id"] = task.TargetResource.EsParam.UniqVpcId
				}

				if task.TargetResource.EsParam.DropInvalidMessage != nil {
					esParamMap["drop_invalid_message"] = task.TargetResource.EsParam.DropInvalidMessage
				}

				if task.TargetResource.EsParam.Index != nil {
					esParamMap["index"] = task.TargetResource.EsParam.Index
				}

				if task.TargetResource.EsParam.DateFormat != nil {
					esParamMap["date_format"] = task.TargetResource.EsParam.DateFormat
				}

				if task.TargetResource.EsParam.ContentKey != nil {
					esParamMap["content_key"] = task.TargetResource.EsParam.ContentKey
				}

				if task.TargetResource.EsParam.DropInvalidJsonMessage != nil {
					esParamMap["drop_invalid_json_message"] = task.TargetResource.EsParam.DropInvalidJsonMessage
				}

				if task.TargetResource.EsParam.DocumentIdField != nil {
					esParamMap["document_id_field"] = task.TargetResource.EsParam.DocumentIdField
				}

				if task.TargetResource.EsParam.IndexType != nil {
					esParamMap["index_type"] = task.TargetResource.EsParam.IndexType
				}

				if task.TargetResource.EsParam.DropCls != nil {
					dropClsMap := map[string]interface{}{}

					if task.TargetResource.EsParam.DropCls.DropInvalidMessageToCls != nil {
						dropClsMap["drop_invalid_message_to_cls"] = task.TargetResource.EsParam.DropCls.DropInvalidMessageToCls
					}

					if task.TargetResource.EsParam.DropCls.DropClsRegion != nil {
						dropClsMap["drop_cls_region"] = task.TargetResource.EsParam.DropCls.DropClsRegion
					}

					if task.TargetResource.EsParam.DropCls.DropClsOwneruin != nil {
						dropClsMap["drop_cls_owneruin"] = task.TargetResource.EsParam.DropCls.DropClsOwneruin
					}

					if task.TargetResource.EsParam.DropCls.DropClsTopicId != nil {
						dropClsMap["drop_cls_topic_id"] = task.TargetResource.EsParam.DropCls.DropClsTopicId
					}

					if task.TargetResource.EsParam.DropCls.DropClsLogSet != nil {
						dropClsMap["drop_cls_log_set"] = task.TargetResource.EsParam.DropCls.DropClsLogSet
					}

					esParamMap["drop_cls"] = []interface{}{dropClsMap}
				}

				if task.TargetResource.EsParam.DatabasePrimaryKey != nil {
					esParamMap["database_primary_key"] = task.TargetResource.EsParam.DatabasePrimaryKey
				}

				if task.TargetResource.EsParam.DropDlq != nil {
					dropDlqMap := map[string]interface{}{}

					if task.TargetResource.EsParam.DropDlq.Type != nil {
						dropDlqMap["type"] = task.TargetResource.EsParam.DropDlq.Type
					}

					if task.TargetResource.EsParam.DropDlq.KafkaParam != nil {
						kafkaParamMap := map[string]interface{}{}

						if task.TargetResource.EsParam.DropDlq.KafkaParam.SelfBuilt != nil {
							kafkaParamMap["self_built"] = task.TargetResource.EsParam.DropDlq.KafkaParam.SelfBuilt
						}

						if task.TargetResource.EsParam.DropDlq.KafkaParam.Resource != nil {
							kafkaParamMap["resource"] = task.TargetResource.EsParam.DropDlq.KafkaParam.Resource
						}

						if task.TargetResource.EsParam.DropDlq.KafkaParam.Topic != nil {
							kafkaParamMap["topic"] = task.TargetResource.EsParam.DropDlq.KafkaParam.Topic
						}

						if task.TargetResource.EsParam.DropDlq.KafkaParam.OffsetType != nil {
							kafkaParamMap["offset_type"] = task.TargetResource.EsParam.DropDlq.KafkaParam.OffsetType
						}

						if task.TargetResource.EsParam.DropDlq.KafkaParam.StartTime != nil {
							kafkaParamMap["start_time"] = task.TargetResource.EsParam.DropDlq.KafkaParam.StartTime
						}

						if task.TargetResource.EsParam.DropDlq.KafkaParam.ResourceName != nil {
							kafkaParamMap["resource_name"] = task.TargetResource.EsParam.DropDlq.KafkaParam.ResourceName
						}

						if task.TargetResource.EsParam.DropDlq.KafkaParam.ZoneId != nil {
							kafkaParamMap["zone_id"] = task.TargetResource.EsParam.DropDlq.KafkaParam.ZoneId
						}

						if task.TargetResource.EsParam.DropDlq.KafkaParam.TopicId != nil {
							kafkaParamMap["topic_id"] = task.TargetResource.EsParam.DropDlq.KafkaParam.TopicId
						}

						if task.TargetResource.EsParam.DropDlq.KafkaParam.PartitionNum != nil {
							kafkaParamMap["partition_num"] = task.TargetResource.EsParam.DropDlq.KafkaParam.PartitionNum
						}

						if task.TargetResource.EsParam.DropDlq.KafkaParam.EnableToleration != nil {
							kafkaParamMap["enable_toleration"] = task.TargetResource.EsParam.DropDlq.KafkaParam.EnableToleration
						}

						if task.TargetResource.EsParam.DropDlq.KafkaParam.QpsLimit != nil {
							kafkaParamMap["qps_limit"] = task.TargetResource.EsParam.DropDlq.KafkaParam.QpsLimit
						}

						if task.TargetResource.EsParam.DropDlq.KafkaParam.TableMappings != nil {
							tableMappingsList := []interface{}{}
							for _, tableMappings := range task.TargetResource.EsParam.DropDlq.KafkaParam.TableMappings {
								tableMappingsMap := map[string]interface{}{}

								if tableMappings.Database != nil {
									tableMappingsMap["database"] = tableMappings.Database
								}

								if tableMappings.Table != nil {
									tableMappingsMap["table"] = tableMappings.Table
								}

								if tableMappings.Topic != nil {
									tableMappingsMap["topic"] = tableMappings.Topic
								}

								if tableMappings.TopicId != nil {
									tableMappingsMap["topic_id"] = tableMappings.TopicId
								}

								tableMappingsList = append(tableMappingsList, tableMappingsMap)
							}

							kafkaParamMap["table_mappings"] = []interface{}{tableMappingsList}
						}

						if task.TargetResource.EsParam.DropDlq.KafkaParam.UseTableMapping != nil {
							kafkaParamMap["use_table_mapping"] = task.TargetResource.EsParam.DropDlq.KafkaParam.UseTableMapping
						}

						if task.TargetResource.EsParam.DropDlq.KafkaParam.UseAutoCreateTopic != nil {
							kafkaParamMap["use_auto_create_topic"] = task.TargetResource.EsParam.DropDlq.KafkaParam.UseAutoCreateTopic
						}

						if task.TargetResource.EsParam.DropDlq.KafkaParam.CompressionType != nil {
							kafkaParamMap["compression_type"] = task.TargetResource.EsParam.DropDlq.KafkaParam.CompressionType
						}

						if task.TargetResource.EsParam.DropDlq.KafkaParam.MsgMultiple != nil {
							kafkaParamMap["msg_multiple"] = task.TargetResource.EsParam.DropDlq.KafkaParam.MsgMultiple
						}

						if task.TargetResource.EsParam.DropDlq.KafkaParam.ConnectorSyncType != nil {
							kafkaParamMap["connector_sync_type"] = task.TargetResource.EsParam.DropDlq.KafkaParam.ConnectorSyncType
						}

						if task.TargetResource.EsParam.DropDlq.KafkaParam.KeepPartition != nil {
							kafkaParamMap["keep_partition"] = task.TargetResource.EsParam.DropDlq.KafkaParam.KeepPartition
						}

						dropDlqMap["kafka_param"] = []interface{}{kafkaParamMap}
					}

					if task.TargetResource.EsParam.DropDlq.RetryInterval != nil {
						dropDlqMap["retry_interval"] = task.TargetResource.EsParam.DropDlq.RetryInterval
					}

					if task.TargetResource.EsParam.DropDlq.MaxRetryAttempts != nil {
						dropDlqMap["max_retry_attempts"] = task.TargetResource.EsParam.DropDlq.MaxRetryAttempts
					}

					if task.TargetResource.EsParam.DropDlq.TopicParam != nil {
						topicParamMap := map[string]interface{}{}

						if task.TargetResource.EsParam.DropDlq.TopicParam.Resource != nil {
							topicParamMap["resource"] = task.TargetResource.EsParam.DropDlq.TopicParam.Resource
						}

						if task.TargetResource.EsParam.DropDlq.TopicParam.OffsetType != nil {
							topicParamMap["offset_type"] = task.TargetResource.EsParam.DropDlq.TopicParam.OffsetType
						}

						if task.TargetResource.EsParam.DropDlq.TopicParam.StartTime != nil {
							topicParamMap["start_time"] = task.TargetResource.EsParam.DropDlq.TopicParam.StartTime
						}

						if task.TargetResource.EsParam.DropDlq.TopicParam.TopicId != nil {
							topicParamMap["topic_id"] = task.TargetResource.EsParam.DropDlq.TopicParam.TopicId
						}

						if task.TargetResource.EsParam.DropDlq.TopicParam.CompressionType != nil {
							topicParamMap["compression_type"] = task.TargetResource.EsParam.DropDlq.TopicParam.CompressionType
						}

						if task.TargetResource.EsParam.DropDlq.TopicParam.UseAutoCreateTopic != nil {
							topicParamMap["use_auto_create_topic"] = task.TargetResource.EsParam.DropDlq.TopicParam.UseAutoCreateTopic
						}

						if task.TargetResource.EsParam.DropDlq.TopicParam.MsgMultiple != nil {
							topicParamMap["msg_multiple"] = task.TargetResource.EsParam.DropDlq.TopicParam.MsgMultiple
						}

						dropDlqMap["topic_param"] = []interface{}{topicParamMap}
					}

					if task.TargetResource.EsParam.DropDlq.DlqType != nil {
						dropDlqMap["dlq_type"] = task.TargetResource.EsParam.DropDlq.DlqType
					}

					esParamMap["drop_dlq"] = []interface{}{dropDlqMap}
				}

				targetResourceMap["es_param"] = []interface{}{esParamMap}
			}

			if task.TargetResource.TdwParam != nil {
				tdwParamMap := map[string]interface{}{}

				if task.TargetResource.TdwParam.Bid != nil {
					tdwParamMap["bid"] = task.TargetResource.TdwParam.Bid
				}

				if task.TargetResource.TdwParam.Tid != nil {
					tdwParamMap["tid"] = task.TargetResource.TdwParam.Tid
				}

				if task.TargetResource.TdwParam.IsDomestic != nil {
					tdwParamMap["is_domestic"] = task.TargetResource.TdwParam.IsDomestic
				}

				if task.TargetResource.TdwParam.TdwHost != nil {
					tdwParamMap["tdw_host"] = task.TargetResource.TdwParam.TdwHost
				}

				if task.TargetResource.TdwParam.TdwPort != nil {
					tdwParamMap["tdw_port"] = task.TargetResource.TdwParam.TdwPort
				}

				targetResourceMap["tdw_param"] = []interface{}{tdwParamMap}
			}

			if task.TargetResource.DtsParam != nil {
				dtsParamMap := map[string]interface{}{}

				if task.TargetResource.DtsParam.Resource != nil {
					dtsParamMap["resource"] = task.TargetResource.DtsParam.Resource
				}

				if task.TargetResource.DtsParam.Ip != nil {
					dtsParamMap["ip"] = task.TargetResource.DtsParam.Ip
				}

				if task.TargetResource.DtsParam.Port != nil {
					dtsParamMap["port"] = task.TargetResource.DtsParam.Port
				}

				if task.TargetResource.DtsParam.Topic != nil {
					dtsParamMap["topic"] = task.TargetResource.DtsParam.Topic
				}

				if task.TargetResource.DtsParam.GroupId != nil {
					dtsParamMap["group_id"] = task.TargetResource.DtsParam.GroupId
				}

				if task.TargetResource.DtsParam.GroupUser != nil {
					dtsParamMap["group_user"] = task.TargetResource.DtsParam.GroupUser
				}

				if task.TargetResource.DtsParam.GroupPassword != nil {
					dtsParamMap["group_password"] = task.TargetResource.DtsParam.GroupPassword
				}

				if task.TargetResource.DtsParam.TranSql != nil {
					dtsParamMap["tran_sql"] = task.TargetResource.DtsParam.TranSql
				}

				targetResourceMap["dts_param"] = []interface{}{dtsParamMap}
			}

			if task.TargetResource.ClickHouseParam != nil {
				clickHouseParamMap := map[string]interface{}{}

				if task.TargetResource.ClickHouseParam.Cluster != nil {
					clickHouseParamMap["cluster"] = task.TargetResource.ClickHouseParam.Cluster
				}

				if task.TargetResource.ClickHouseParam.Database != nil {
					clickHouseParamMap["database"] = task.TargetResource.ClickHouseParam.Database
				}

				if task.TargetResource.ClickHouseParam.Table != nil {
					clickHouseParamMap["table"] = task.TargetResource.ClickHouseParam.Table
				}

				if task.TargetResource.ClickHouseParam.Schema != nil {
					schemaList := []interface{}{}
					for _, schema := range task.TargetResource.ClickHouseParam.Schema {
						schemaMap := map[string]interface{}{}

						if schema.ColumnName != nil {
							schemaMap["column_name"] = schema.ColumnName
						}

						if schema.JsonKey != nil {
							schemaMap["json_key"] = schema.JsonKey
						}

						if schema.Type != nil {
							schemaMap["type"] = schema.Type
						}

						if schema.AllowNull != nil {
							schemaMap["allow_null"] = schema.AllowNull
						}

						schemaList = append(schemaList, schemaMap)
					}

					clickHouseParamMap["schema"] = []interface{}{schemaList}
				}

				if task.TargetResource.ClickHouseParam.Resource != nil {
					clickHouseParamMap["resource"] = task.TargetResource.ClickHouseParam.Resource
				}

				if task.TargetResource.ClickHouseParam.Ip != nil {
					clickHouseParamMap["ip"] = task.TargetResource.ClickHouseParam.Ip
				}

				if task.TargetResource.ClickHouseParam.Port != nil {
					clickHouseParamMap["port"] = task.TargetResource.ClickHouseParam.Port
				}

				if task.TargetResource.ClickHouseParam.UserName != nil {
					clickHouseParamMap["user_name"] = task.TargetResource.ClickHouseParam.UserName
				}

				if task.TargetResource.ClickHouseParam.Password != nil {
					clickHouseParamMap["password"] = task.TargetResource.ClickHouseParam.Password
				}

				if task.TargetResource.ClickHouseParam.ServiceVip != nil {
					clickHouseParamMap["service_vip"] = task.TargetResource.ClickHouseParam.ServiceVip
				}

				if task.TargetResource.ClickHouseParam.UniqVpcId != nil {
					clickHouseParamMap["uniq_vpc_id"] = task.TargetResource.ClickHouseParam.UniqVpcId
				}

				if task.TargetResource.ClickHouseParam.SelfBuilt != nil {
					clickHouseParamMap["self_built"] = task.TargetResource.ClickHouseParam.SelfBuilt
				}

				if task.TargetResource.ClickHouseParam.DropInvalidMessage != nil {
					clickHouseParamMap["drop_invalid_message"] = task.TargetResource.ClickHouseParam.DropInvalidMessage
				}

				if task.TargetResource.ClickHouseParam.Type != nil {
					clickHouseParamMap["type"] = task.TargetResource.ClickHouseParam.Type
				}

				if task.TargetResource.ClickHouseParam.DropCls != nil {
					dropClsMap := map[string]interface{}{}

					if task.TargetResource.ClickHouseParam.DropCls.DropInvalidMessageToCls != nil {
						dropClsMap["drop_invalid_message_to_cls"] = task.TargetResource.ClickHouseParam.DropCls.DropInvalidMessageToCls
					}

					if task.TargetResource.ClickHouseParam.DropCls.DropClsRegion != nil {
						dropClsMap["drop_cls_region"] = task.TargetResource.ClickHouseParam.DropCls.DropClsRegion
					}

					if task.TargetResource.ClickHouseParam.DropCls.DropClsOwneruin != nil {
						dropClsMap["drop_cls_owneruin"] = task.TargetResource.ClickHouseParam.DropCls.DropClsOwneruin
					}

					if task.TargetResource.ClickHouseParam.DropCls.DropClsTopicId != nil {
						dropClsMap["drop_cls_topic_id"] = task.TargetResource.ClickHouseParam.DropCls.DropClsTopicId
					}

					if task.TargetResource.ClickHouseParam.DropCls.DropClsLogSet != nil {
						dropClsMap["drop_cls_log_set"] = task.TargetResource.ClickHouseParam.DropCls.DropClsLogSet
					}

					clickHouseParamMap["drop_cls"] = []interface{}{dropClsMap}
				}

				targetResourceMap["click_house_param"] = []interface{}{clickHouseParamMap}
			}

			if task.TargetResource.ClsParam != nil {
				clsParamMap := map[string]interface{}{}

				if task.TargetResource.ClsParam.DecodeJson != nil {
					clsParamMap["decode_json"] = task.TargetResource.ClsParam.DecodeJson
				}

				if task.TargetResource.ClsParam.Resource != nil {
					clsParamMap["resource"] = task.TargetResource.ClsParam.Resource
				}

				if task.TargetResource.ClsParam.LogSet != nil {
					clsParamMap["log_set"] = task.TargetResource.ClsParam.LogSet
				}

				if task.TargetResource.ClsParam.ContentKey != nil {
					clsParamMap["content_key"] = task.TargetResource.ClsParam.ContentKey
				}

				if task.TargetResource.ClsParam.TimeField != nil {
					clsParamMap["time_field"] = task.TargetResource.ClsParam.TimeField
				}

				targetResourceMap["cls_param"] = []interface{}{clsParamMap}
			}

			if task.TargetResource.CosParam != nil {
				cosParamMap := map[string]interface{}{}

				if task.TargetResource.CosParam.BucketName != nil {
					cosParamMap["bucket_name"] = task.TargetResource.CosParam.BucketName
				}

				if task.TargetResource.CosParam.Region != nil {
					cosParamMap["region"] = task.TargetResource.CosParam.Region
				}

				if task.TargetResource.CosParam.ObjectKey != nil {
					cosParamMap["object_key"] = task.TargetResource.CosParam.ObjectKey
				}

				if task.TargetResource.CosParam.AggregateBatchSize != nil {
					cosParamMap["aggregate_batch_size"] = task.TargetResource.CosParam.AggregateBatchSize
				}

				if task.TargetResource.CosParam.AggregateInterval != nil {
					cosParamMap["aggregate_interval"] = task.TargetResource.CosParam.AggregateInterval
				}

				if task.TargetResource.CosParam.FormatOutputType != nil {
					cosParamMap["format_output_type"] = task.TargetResource.CosParam.FormatOutputType
				}

				if task.TargetResource.CosParam.ObjectKeyPrefix != nil {
					cosParamMap["object_key_prefix"] = task.TargetResource.CosParam.ObjectKeyPrefix
				}

				if task.TargetResource.CosParam.DirectoryTimeFormat != nil {
					cosParamMap["directory_time_format"] = task.TargetResource.CosParam.DirectoryTimeFormat
				}

				targetResourceMap["cos_param"] = []interface{}{cosParamMap}
			}

			if task.TargetResource.MySQLParam != nil {
				mySQLParamMap := map[string]interface{}{}

				if task.TargetResource.MySQLParam.Database != nil {
					mySQLParamMap["database"] = task.TargetResource.MySQLParam.Database
				}

				if task.TargetResource.MySQLParam.Table != nil {
					mySQLParamMap["table"] = task.TargetResource.MySQLParam.Table
				}

				if task.TargetResource.MySQLParam.Resource != nil {
					mySQLParamMap["resource"] = task.TargetResource.MySQLParam.Resource
				}

				if task.TargetResource.MySQLParam.SnapshotMode != nil {
					mySQLParamMap["snapshot_mode"] = task.TargetResource.MySQLParam.SnapshotMode
				}

				if task.TargetResource.MySQLParam.DdlTopic != nil {
					mySQLParamMap["ddl_topic"] = task.TargetResource.MySQLParam.DdlTopic
				}

				if task.TargetResource.MySQLParam.DataSourceMonitorMode != nil {
					mySQLParamMap["data_source_monitor_mode"] = task.TargetResource.MySQLParam.DataSourceMonitorMode
				}

				if task.TargetResource.MySQLParam.DataSourceMonitorResource != nil {
					mySQLParamMap["data_source_monitor_resource"] = task.TargetResource.MySQLParam.DataSourceMonitorResource
				}

				if task.TargetResource.MySQLParam.DataSourceIncrementMode != nil {
					mySQLParamMap["data_source_increment_mode"] = task.TargetResource.MySQLParam.DataSourceIncrementMode
				}

				if task.TargetResource.MySQLParam.DataSourceIncrementColumn != nil {
					mySQLParamMap["data_source_increment_column"] = task.TargetResource.MySQLParam.DataSourceIncrementColumn
				}

				if task.TargetResource.MySQLParam.DataSourceStartFrom != nil {
					mySQLParamMap["data_source_start_from"] = task.TargetResource.MySQLParam.DataSourceStartFrom
				}

				if task.TargetResource.MySQLParam.DataTargetInsertMode != nil {
					mySQLParamMap["data_target_insert_mode"] = task.TargetResource.MySQLParam.DataTargetInsertMode
				}

				if task.TargetResource.MySQLParam.DataTargetPrimaryKeyField != nil {
					mySQLParamMap["data_target_primary_key_field"] = task.TargetResource.MySQLParam.DataTargetPrimaryKeyField
				}

				if task.TargetResource.MySQLParam.DataTargetRecordMapping != nil {
					dataTargetRecordMappingList := []interface{}{}
					for _, dataTargetRecordMapping := range task.TargetResource.MySQLParam.DataTargetRecordMapping {
						dataTargetRecordMappingMap := map[string]interface{}{}

						if dataTargetRecordMapping.JsonKey != nil {
							dataTargetRecordMappingMap["json_key"] = dataTargetRecordMapping.JsonKey
						}

						if dataTargetRecordMapping.Type != nil {
							dataTargetRecordMappingMap["type"] = dataTargetRecordMapping.Type
						}

						if dataTargetRecordMapping.AllowNull != nil {
							dataTargetRecordMappingMap["allow_null"] = dataTargetRecordMapping.AllowNull
						}

						if dataTargetRecordMapping.ColumnName != nil {
							dataTargetRecordMappingMap["column_name"] = dataTargetRecordMapping.ColumnName
						}

						if dataTargetRecordMapping.ExtraInfo != nil {
							dataTargetRecordMappingMap["extra_info"] = dataTargetRecordMapping.ExtraInfo
						}

						if dataTargetRecordMapping.ColumnSize != nil {
							dataTargetRecordMappingMap["column_size"] = dataTargetRecordMapping.ColumnSize
						}

						if dataTargetRecordMapping.DecimalDigits != nil {
							dataTargetRecordMappingMap["decimal_digits"] = dataTargetRecordMapping.DecimalDigits
						}

						if dataTargetRecordMapping.AutoIncrement != nil {
							dataTargetRecordMappingMap["auto_increment"] = dataTargetRecordMapping.AutoIncrement
						}

						if dataTargetRecordMapping.DefaultValue != nil {
							dataTargetRecordMappingMap["default_value"] = dataTargetRecordMapping.DefaultValue
						}

						dataTargetRecordMappingList = append(dataTargetRecordMappingList, dataTargetRecordMappingMap)
					}

					mySQLParamMap["data_target_record_mapping"] = []interface{}{dataTargetRecordMappingList}
				}

				if task.TargetResource.MySQLParam.TopicRegex != nil {
					mySQLParamMap["topic_regex"] = task.TargetResource.MySQLParam.TopicRegex
				}

				if task.TargetResource.MySQLParam.TopicReplacement != nil {
					mySQLParamMap["topic_replacement"] = task.TargetResource.MySQLParam.TopicReplacement
				}

				if task.TargetResource.MySQLParam.KeyColumns != nil {
					mySQLParamMap["key_columns"] = task.TargetResource.MySQLParam.KeyColumns
				}

				if task.TargetResource.MySQLParam.DropInvalidMessage != nil {
					mySQLParamMap["drop_invalid_message"] = task.TargetResource.MySQLParam.DropInvalidMessage
				}

				if task.TargetResource.MySQLParam.DropCls != nil {
					dropClsMap := map[string]interface{}{}

					if task.TargetResource.MySQLParam.DropCls.DropInvalidMessageToCls != nil {
						dropClsMap["drop_invalid_message_to_cls"] = task.TargetResource.MySQLParam.DropCls.DropInvalidMessageToCls
					}

					if task.TargetResource.MySQLParam.DropCls.DropClsRegion != nil {
						dropClsMap["drop_cls_region"] = task.TargetResource.MySQLParam.DropCls.DropClsRegion
					}

					if task.TargetResource.MySQLParam.DropCls.DropClsOwneruin != nil {
						dropClsMap["drop_cls_owneruin"] = task.TargetResource.MySQLParam.DropCls.DropClsOwneruin
					}

					if task.TargetResource.MySQLParam.DropCls.DropClsTopicId != nil {
						dropClsMap["drop_cls_topic_id"] = task.TargetResource.MySQLParam.DropCls.DropClsTopicId
					}

					if task.TargetResource.MySQLParam.DropCls.DropClsLogSet != nil {
						dropClsMap["drop_cls_log_set"] = task.TargetResource.MySQLParam.DropCls.DropClsLogSet
					}

					mySQLParamMap["drop_cls"] = []interface{}{dropClsMap}
				}

				if task.TargetResource.MySQLParam.OutputFormat != nil {
					mySQLParamMap["output_format"] = task.TargetResource.MySQLParam.OutputFormat
				}

				if task.TargetResource.MySQLParam.IsTablePrefix != nil {
					mySQLParamMap["is_table_prefix"] = task.TargetResource.MySQLParam.IsTablePrefix
				}

				if task.TargetResource.MySQLParam.IncludeContentChanges != nil {
					mySQLParamMap["include_content_changes"] = task.TargetResource.MySQLParam.IncludeContentChanges
				}

				if task.TargetResource.MySQLParam.IncludeQuery != nil {
					mySQLParamMap["include_query"] = task.TargetResource.MySQLParam.IncludeQuery
				}

				if task.TargetResource.MySQLParam.RecordWithSchema != nil {
					mySQLParamMap["record_with_schema"] = task.TargetResource.MySQLParam.RecordWithSchema
				}

				if task.TargetResource.MySQLParam.SignalDatabase != nil {
					mySQLParamMap["signal_database"] = task.TargetResource.MySQLParam.SignalDatabase
				}

				if task.TargetResource.MySQLParam.IsTableRegular != nil {
					mySQLParamMap["is_table_regular"] = task.TargetResource.MySQLParam.IsTableRegular
				}

				targetResourceMap["my_sql_param"] = []interface{}{mySQLParamMap}
			}

			if task.TargetResource.PostgreSQLParam != nil {
				postgreSQLParamMap := map[string]interface{}{}

				if task.TargetResource.PostgreSQLParam.Database != nil {
					postgreSQLParamMap["database"] = task.TargetResource.PostgreSQLParam.Database
				}

				if task.TargetResource.PostgreSQLParam.Table != nil {
					postgreSQLParamMap["table"] = task.TargetResource.PostgreSQLParam.Table
				}

				if task.TargetResource.PostgreSQLParam.Resource != nil {
					postgreSQLParamMap["resource"] = task.TargetResource.PostgreSQLParam.Resource
				}

				if task.TargetResource.PostgreSQLParam.PluginName != nil {
					postgreSQLParamMap["plugin_name"] = task.TargetResource.PostgreSQLParam.PluginName
				}

				if task.TargetResource.PostgreSQLParam.SnapshotMode != nil {
					postgreSQLParamMap["snapshot_mode"] = task.TargetResource.PostgreSQLParam.SnapshotMode
				}

				if task.TargetResource.PostgreSQLParam.DataFormat != nil {
					postgreSQLParamMap["data_format"] = task.TargetResource.PostgreSQLParam.DataFormat
				}

				if task.TargetResource.PostgreSQLParam.DataTargetInsertMode != nil {
					postgreSQLParamMap["data_target_insert_mode"] = task.TargetResource.PostgreSQLParam.DataTargetInsertMode
				}

				if task.TargetResource.PostgreSQLParam.DataTargetPrimaryKeyField != nil {
					postgreSQLParamMap["data_target_primary_key_field"] = task.TargetResource.PostgreSQLParam.DataTargetPrimaryKeyField
				}

				if task.TargetResource.PostgreSQLParam.DataTargetRecordMapping != nil {
					dataTargetRecordMappingList := []interface{}{}
					for _, dataTargetRecordMapping := range task.TargetResource.PostgreSQLParam.DataTargetRecordMapping {
						dataTargetRecordMappingMap := map[string]interface{}{}

						if dataTargetRecordMapping.JsonKey != nil {
							dataTargetRecordMappingMap["json_key"] = dataTargetRecordMapping.JsonKey
						}

						if dataTargetRecordMapping.Type != nil {
							dataTargetRecordMappingMap["type"] = dataTargetRecordMapping.Type
						}

						if dataTargetRecordMapping.AllowNull != nil {
							dataTargetRecordMappingMap["allow_null"] = dataTargetRecordMapping.AllowNull
						}

						if dataTargetRecordMapping.ColumnName != nil {
							dataTargetRecordMappingMap["column_name"] = dataTargetRecordMapping.ColumnName
						}

						if dataTargetRecordMapping.ExtraInfo != nil {
							dataTargetRecordMappingMap["extra_info"] = dataTargetRecordMapping.ExtraInfo
						}

						if dataTargetRecordMapping.ColumnSize != nil {
							dataTargetRecordMappingMap["column_size"] = dataTargetRecordMapping.ColumnSize
						}

						if dataTargetRecordMapping.DecimalDigits != nil {
							dataTargetRecordMappingMap["decimal_digits"] = dataTargetRecordMapping.DecimalDigits
						}

						if dataTargetRecordMapping.AutoIncrement != nil {
							dataTargetRecordMappingMap["auto_increment"] = dataTargetRecordMapping.AutoIncrement
						}

						if dataTargetRecordMapping.DefaultValue != nil {
							dataTargetRecordMappingMap["default_value"] = dataTargetRecordMapping.DefaultValue
						}

						dataTargetRecordMappingList = append(dataTargetRecordMappingList, dataTargetRecordMappingMap)
					}

					postgreSQLParamMap["data_target_record_mapping"] = []interface{}{dataTargetRecordMappingList}
				}

				if task.TargetResource.PostgreSQLParam.DropInvalidMessage != nil {
					postgreSQLParamMap["drop_invalid_message"] = task.TargetResource.PostgreSQLParam.DropInvalidMessage
				}

				if task.TargetResource.PostgreSQLParam.IsTableRegular != nil {
					postgreSQLParamMap["is_table_regular"] = task.TargetResource.PostgreSQLParam.IsTableRegular
				}

				if task.TargetResource.PostgreSQLParam.KeyColumns != nil {
					postgreSQLParamMap["key_columns"] = task.TargetResource.PostgreSQLParam.KeyColumns
				}

				if task.TargetResource.PostgreSQLParam.RecordWithSchema != nil {
					postgreSQLParamMap["record_with_schema"] = task.TargetResource.PostgreSQLParam.RecordWithSchema
				}

				targetResourceMap["postgre_sql_param"] = []interface{}{postgreSQLParamMap}
			}

			if task.TargetResource.TopicParam != nil {
				topicParamMap := map[string]interface{}{}

				if task.TargetResource.TopicParam.Resource != nil {
					topicParamMap["resource"] = task.TargetResource.TopicParam.Resource
				}

				if task.TargetResource.TopicParam.OffsetType != nil {
					topicParamMap["offset_type"] = task.TargetResource.TopicParam.OffsetType
				}

				if task.TargetResource.TopicParam.StartTime != nil {
					topicParamMap["start_time"] = task.TargetResource.TopicParam.StartTime
				}

				if task.TargetResource.TopicParam.TopicId != nil {
					topicParamMap["topic_id"] = task.TargetResource.TopicParam.TopicId
				}

				if task.TargetResource.TopicParam.CompressionType != nil {
					topicParamMap["compression_type"] = task.TargetResource.TopicParam.CompressionType
				}

				if task.TargetResource.TopicParam.UseAutoCreateTopic != nil {
					topicParamMap["use_auto_create_topic"] = task.TargetResource.TopicParam.UseAutoCreateTopic
				}

				if task.TargetResource.TopicParam.MsgMultiple != nil {
					topicParamMap["msg_multiple"] = task.TargetResource.TopicParam.MsgMultiple
				}

				targetResourceMap["topic_param"] = []interface{}{topicParamMap}
			}

			if task.TargetResource.MariaDBParam != nil {
				mariaDBParamMap := map[string]interface{}{}

				if task.TargetResource.MariaDBParam.Database != nil {
					mariaDBParamMap["database"] = task.TargetResource.MariaDBParam.Database
				}

				if task.TargetResource.MariaDBParam.Table != nil {
					mariaDBParamMap["table"] = task.TargetResource.MariaDBParam.Table
				}

				if task.TargetResource.MariaDBParam.Resource != nil {
					mariaDBParamMap["resource"] = task.TargetResource.MariaDBParam.Resource
				}

				if task.TargetResource.MariaDBParam.SnapshotMode != nil {
					mariaDBParamMap["snapshot_mode"] = task.TargetResource.MariaDBParam.SnapshotMode
				}

				if task.TargetResource.MariaDBParam.KeyColumns != nil {
					mariaDBParamMap["key_columns"] = task.TargetResource.MariaDBParam.KeyColumns
				}

				if task.TargetResource.MariaDBParam.IsTablePrefix != nil {
					mariaDBParamMap["is_table_prefix"] = task.TargetResource.MariaDBParam.IsTablePrefix
				}

				if task.TargetResource.MariaDBParam.OutputFormat != nil {
					mariaDBParamMap["output_format"] = task.TargetResource.MariaDBParam.OutputFormat
				}

				if task.TargetResource.MariaDBParam.IncludeContentChanges != nil {
					mariaDBParamMap["include_content_changes"] = task.TargetResource.MariaDBParam.IncludeContentChanges
				}

				if task.TargetResource.MariaDBParam.IncludeQuery != nil {
					mariaDBParamMap["include_query"] = task.TargetResource.MariaDBParam.IncludeQuery
				}

				if task.TargetResource.MariaDBParam.RecordWithSchema != nil {
					mariaDBParamMap["record_with_schema"] = task.TargetResource.MariaDBParam.RecordWithSchema
				}

				targetResourceMap["maria_db_param"] = []interface{}{mariaDBParamMap}
			}

			if task.TargetResource.SQLServerParam != nil {
				sQLServerParamMap := map[string]interface{}{}

				if task.TargetResource.SQLServerParam.Database != nil {
					sQLServerParamMap["database"] = task.TargetResource.SQLServerParam.Database
				}

				if task.TargetResource.SQLServerParam.Table != nil {
					sQLServerParamMap["table"] = task.TargetResource.SQLServerParam.Table
				}

				if task.TargetResource.SQLServerParam.Resource != nil {
					sQLServerParamMap["resource"] = task.TargetResource.SQLServerParam.Resource
				}

				if task.TargetResource.SQLServerParam.SnapshotMode != nil {
					sQLServerParamMap["snapshot_mode"] = task.TargetResource.SQLServerParam.SnapshotMode
				}

				targetResourceMap["sql_server_param"] = []interface{}{sQLServerParamMap}
			}

			if task.TargetResource.CtsdbParam != nil {
				ctsdbParamMap := map[string]interface{}{}

				if task.TargetResource.CtsdbParam.Resource != nil {
					ctsdbParamMap["resource"] = task.TargetResource.CtsdbParam.Resource
				}

				if task.TargetResource.CtsdbParam.CtsdbMetric != nil {
					ctsdbParamMap["ctsdb_metric"] = task.TargetResource.CtsdbParam.CtsdbMetric
				}

				targetResourceMap["ctsdb_param"] = []interface{}{ctsdbParamMap}
			}

			if task.TargetResource.ScfParam != nil {
				scfParamMap := map[string]interface{}{}

				if task.TargetResource.ScfParam.FunctionName != nil {
					scfParamMap["function_name"] = task.TargetResource.ScfParam.FunctionName
				}

				if task.TargetResource.ScfParam.Namespace != nil {
					scfParamMap["namespace"] = task.TargetResource.ScfParam.Namespace
				}

				if task.TargetResource.ScfParam.Qualifier != nil {
					scfParamMap["qualifier"] = task.TargetResource.ScfParam.Qualifier
				}

				if task.TargetResource.ScfParam.BatchSize != nil {
					scfParamMap["batch_size"] = task.TargetResource.ScfParam.BatchSize
				}

				if task.TargetResource.ScfParam.MaxRetries != nil {
					scfParamMap["max_retries"] = task.TargetResource.ScfParam.MaxRetries
				}

				targetResourceMap["scf_param"] = []interface{}{scfParamMap}
			}

			taskMap["target_resource"] = []interface{}{targetResourceMap}
		}

		if task.CreateTime != nil {
			taskMap["create_time"] = task.CreateTime
		}

		if task.ErrorMessage != nil {
			taskMap["error_message"] = task.ErrorMessage
		}

		if task.TaskProgress != nil {
			taskMap["task_progress"] = task.TaskProgress
		}

		if task.TaskCurrentStep != nil {
			taskMap["task_current_step"] = task.TaskCurrentStep
		}

		if task.DatahubId != nil {
			taskMap["datahub_id"] = task.DatahubId
		}

		if task.StepList != nil {
			stepList := make([]string, 0)
			for _, v := range task.StepList {
				stepList = append(stepList, *v)
			}
			taskMap["step_list"] = stepList
		}

		taskList = append(taskList, taskMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("task_list", taskList)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), taskList); e != nil {
			return e
		}
	}
	return nil
}
