package ckafka

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCkafkaConnectResource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaConnectResourceRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "连接 来源 类型",
			},

			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Keyword 对于 search。",
			},

			"offset": {
				Optional:    true,
				Default:     0,
				Type:        schema.TypeInt,
				Description: "Page 偏移量，默认为 0。",
			},

			"limit": {
				Optional:    true,
				Default:     20,
				Type:        schema.TypeInt,
				Description: "Return 数量， 默认为 20， 最大 是 100。",
			},

			"resource_region": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Keyword 查询 的 连接 来源，查询 连接 在 连接 management 列表 在 本地 地域 according 到 地域 (仅 support 连接 来源 containing 地域 input)。",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Connection 来源 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 连接 sources。",
						},
						"connect_resource_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Resource List。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resource ID。",
									},
									"resource_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "资源名称",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "描述",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "资源类型",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Resource 状态",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创建时间。",
									},
									"error_message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "错误 Messages。",
									},
									"datahub_task_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 Datahub tasks associated 使用 此 连接 来源",
									},
									"current_step": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "当前 step 的 连接 来源",
									},
									"task_progress": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Creation progress percentage。",
									},
									"step_list": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Step List。",
									},
									"dts_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Dts 配置，返回 当 类型 是 DTS。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Dts 端口",
												},
												"group_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ID Dts 消费者 组。",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "UserName 的 Dts 消费者 组。",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "密码 的 Dts 消费者 组。",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Dts ID。",
												},
												"topic": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Topic subscribed 通过 Dts。",
												},
												"is_update": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否update 到 associated Datahub 任务。",
												},
											},
										},
									},
									"mongo_db_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Mongo DB 配置，返回 当 类型 是 MONGODB。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "MongoDB 端口",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "用户名 的 连接 来源",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "密码 的 连接 来源",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 资源 的 连接 来源",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否connection 来源 是 自建集群",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 VIP 的 连接 来源，当 它 是 腾讯云实例，它 为必填项。",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "vpc ID 的 连接 来源，当 它 是 腾讯云实例，它 为必填项。",
												},
												"is_update": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否update 到 associated Datahub 任务。",
												},
											},
										},
									},
									"es_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Es 配置，返回 当 类型 是 ES。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "ES 端口",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "用户名 的 连接 来源",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "密码 的 连接 来源",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 资源 的 连接 来源",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否connection 来源 是 自建集群",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 VIP 的 连接 来源，当 它 是 腾讯云实例，它 为必填项。",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "vpc ID 的 连接 来源，当 它 是 腾讯云实例，它 为必填项。",
												},
												"is_update": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否update 到 associated Datahub 任务。",
												},
											},
										},
									},
									"clickhouse_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "ClickHouse 配置，返回 当 类型 是 CLICKHOUSE。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "ClickHouse 端口",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "用户名 的 连接 来源",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "密码 的 连接 来源",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 资源 的 连接 来源",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否connection 来源 是 自建集群",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 VIP 的 连接 来源，当 它 是 腾讯云实例，它 为必填项。",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "vpc ID 的 连接 来源，当 它 是 腾讯云实例，它 为必填项。",
												},
												"is_update": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否update 到 associated Datahub 任务。",
												},
											},
										},
									},
									"mysql_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Mysql 配置，返回 当 类型 是 MYSQL 或 TDSQL C MYSQL。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "MySQL 端口",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "用户名 的 连接 来源",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "密码 的 连接 来源",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "MySQL 实例 资源 的 连接 来源",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 VIP 的 连接 来源，当 它 是 腾讯云实例，它 为必填项。",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "vpc ID 的 连接 来源，当 它 是 腾讯云实例，它 为必填项。",
												},
												"is_update": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否update 到 associated Datahub 任务。",
												},
												"cluster_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "必填 当 类型 是 TDSQL C_MYSQL。",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Mysql 是否connection 来源 是 自建集群",
												},
											},
										},
									},
									"postgre_sql_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Postgresql 配置，返回 当 类型 是 POSTGRESQL 或 TDSQL C POSTGRESQL。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "PostgreSQL 端口",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "用户名 的 连接 来源",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "密码 的 连接 来源",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 资源 的 连接 来源",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 VIP 的 连接 来源，当 它 是 腾讯云实例，它 为必填项。",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "vpc ID 的 连接 来源，当 它 是 腾讯云实例，它 为必填项。",
												},
												"cluster_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "必填 当 类型 是 TDSQL C_POSTGRESQL。",
												},
												"is_update": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否update 到 associated Datahub 任务。",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否connection 来源 是 自建集群",
												},
											},
										},
									},
									"maria_db_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Mariadb 配置，返回 当 类型 是 MARIADB。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "MariaDB 端口",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "用户名 的 连接 来源",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "密码 的 连接 来源",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 资源 的 连接 来源",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 VIP 的 连接 来源，当 它 是 腾讯云实例，它 为必填项。",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "vpc ID 的 连接 来源，当 它 是 腾讯云实例，它 为必填项。",
												},
												"is_update": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否update 到 associated Datahub 任务。",
												},
											},
										},
									},
									"sql_server_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "SQL Server 配置，返回 当 类型 是 SQLSERVER。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "SQLServer 端口",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "用户名 的 连接 来源",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "密码 的 连接 来源",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 资源 的 连接 来源",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 VIP 的 连接 来源，当 它 是 腾讯云实例，它 为必填项。",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "vpc ID 的 连接 来源，当 它 是 腾讯云实例，它 为必填项。",
												},
												"is_update": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否update 到 associated Dip 任务。",
												},
											},
										},
									},
									"ctsdb_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Ctsdb 配置，返回 当 类型 是 CTSDB。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Ctsdb 端口",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Ctsdb VIP",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Ctsdb vpcId。",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "用户名 的 连接 来源",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "密码 的 连接 来源",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 资源 的 连接 来源",
												},
											},
										},
									},
									"doris_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Doris Configuration，返回 当 类型 是 DORIS。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Doris jdbc Load balancing 连接 端口，usually mapped 到 端口 9030 的 fe。",
												},
												"user_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "用户名 的 连接 来源",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "密码 的 连接 来源",
												},
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 资源 的 连接 来源",
												},
												"service_vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 VIP 的 连接 来源，当 它 是 腾讯云实例，它 为必填项。",
												},
												"uniq_vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "vpc ID 的 连接 来源，当 它 是 腾讯云实例，它 为必填项。",
												},
												"is_update": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否update 到 associated Datahub 任务。",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否connection 来源 是 自建集群",
												},
												"be_port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Doris's http load balancing 连接 端口，usually mapped 到 是's 8040 端口",
												},
											},
										},
									},
									"kafka_connect_param": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Kafka 配置，返回 当 类型 是 KAFKA。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"resource": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 资源 的 Kafka 连接 来源，必填 当 不 self-built。",
												},
												"self_built": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否为a 自建集群",
												},
												"is_update": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "是否update 到 associated Dip 任务。",
												},
												"broker_address": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Broker 地址 对于 Kafka 连接，必填 对于 self-build。",
												},
												"region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "实例 资源 地域 的 CKafka 连接 来源，必填 当 crossing regions。",
												},
											},
										},
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

func dataSourceTencentCloudCkafkaConnectResourceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ckafka_connect_resource.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("type"); ok {
		paramMap["type"] = v.(string)
	}

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["search_word"] = v.(string)
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["offset"] = v.(int)
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["limit"] = v.(int)
	}

	if v, ok := d.GetOk("resource_region"); ok {
		paramMap["resource_region"] = v.(string)
	}

	service := CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var result *ckafka.DescribeConnectResourcesResp

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeCkafkaConnectResourceByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	describeConnectResourcesRespMap := make(map[string]interface{})
	if result != nil {

		if result.TotalCount != nil {
			describeConnectResourcesRespMap["total_count"] = result.TotalCount
		}

		if result.ConnectResourceList != nil {
			connectResourceListList := []interface{}{}
			for _, connectResource := range result.ConnectResourceList {
				connectResourceListMap := map[string]interface{}{}

				if connectResource.ResourceId != nil {
					connectResourceListMap["resource_id"] = connectResource.ResourceId
				}

				if connectResource.ResourceName != nil {
					connectResourceListMap["resource_name"] = connectResource.ResourceName
				}

				if connectResource.Description != nil {
					connectResourceListMap["description"] = connectResource.Description
				}

				if connectResource.Type != nil {
					connectResourceListMap["type"] = connectResource.Type
				}

				if connectResource.Status != nil {
					connectResourceListMap["status"] = connectResource.Status
				}

				if connectResource.CreateTime != nil {
					connectResourceListMap["create_time"] = connectResource.CreateTime
				}

				if connectResource.ErrorMessage != nil {
					connectResourceListMap["error_message"] = connectResource.ErrorMessage
				}

				if connectResource.DatahubTaskCount != nil {
					connectResourceListMap["datahub_task_count"] = connectResource.DatahubTaskCount
				}

				if connectResource.CurrentStep != nil {
					connectResourceListMap["current_step"] = connectResource.CurrentStep
				}

				if connectResource.TaskProgress != nil {
					connectResourceListMap["task_progress"] = connectResource.TaskProgress
				}

				if connectResource.StepList != nil {
					connectResourceListMap["step_list"] = connectResource.StepList
				}

				if connectResource.DtsConnectParam != nil {
					dtsConnectParamMap := map[string]interface{}{}

					if connectResource.DtsConnectParam.Port != nil {
						dtsConnectParamMap["port"] = connectResource.DtsConnectParam.Port
					}

					if connectResource.DtsConnectParam.GroupId != nil {
						dtsConnectParamMap["group_id"] = connectResource.DtsConnectParam.GroupId
					}

					if connectResource.DtsConnectParam.UserName != nil {
						dtsConnectParamMap["user_name"] = connectResource.DtsConnectParam.UserName
					}

					if connectResource.DtsConnectParam.Password != nil {
						dtsConnectParamMap["password"] = connectResource.DtsConnectParam.Password
					}

					if connectResource.DtsConnectParam.Resource != nil {
						dtsConnectParamMap["resource"] = connectResource.DtsConnectParam.Resource
					}

					if connectResource.DtsConnectParam.Topic != nil {
						dtsConnectParamMap["topic"] = connectResource.DtsConnectParam.Topic
					}

					if connectResource.DtsConnectParam.IsUpdate != nil {
						dtsConnectParamMap["is_update"] = connectResource.DtsConnectParam.IsUpdate
					}

					connectResourceListMap["dts_connect_param"] = []interface{}{dtsConnectParamMap}
				}

				if connectResource.MongoDBConnectParam != nil {
					mongoDBConnectParamMap := map[string]interface{}{}

					if connectResource.MongoDBConnectParam.Port != nil {
						mongoDBConnectParamMap["port"] = connectResource.MongoDBConnectParam.Port
					}

					if connectResource.MongoDBConnectParam.UserName != nil {
						mongoDBConnectParamMap["user_name"] = connectResource.MongoDBConnectParam.UserName
					}

					if connectResource.MongoDBConnectParam.Password != nil {
						mongoDBConnectParamMap["password"] = connectResource.MongoDBConnectParam.Password
					}

					if connectResource.MongoDBConnectParam.Resource != nil {
						mongoDBConnectParamMap["resource"] = connectResource.MongoDBConnectParam.Resource
					}

					if connectResource.MongoDBConnectParam.SelfBuilt != nil {
						mongoDBConnectParamMap["self_built"] = connectResource.MongoDBConnectParam.SelfBuilt
					}

					if connectResource.MongoDBConnectParam.ServiceVip != nil {
						mongoDBConnectParamMap["service_vip"] = connectResource.MongoDBConnectParam.ServiceVip
					}

					if connectResource.MongoDBConnectParam.UniqVpcId != nil {
						mongoDBConnectParamMap["uniq_vpc_id"] = connectResource.MongoDBConnectParam.UniqVpcId
					}

					if connectResource.MongoDBConnectParam.IsUpdate != nil {
						mongoDBConnectParamMap["is_update"] = connectResource.MongoDBConnectParam.IsUpdate
					}

					connectResourceListMap["mongo_db_connect_param"] = []interface{}{mongoDBConnectParamMap}
				}

				if connectResource.EsConnectParam != nil {
					esConnectParamMap := map[string]interface{}{}

					if connectResource.EsConnectParam.Port != nil {
						esConnectParamMap["port"] = connectResource.EsConnectParam.Port
					}

					if connectResource.EsConnectParam.UserName != nil {
						esConnectParamMap["user_name"] = connectResource.EsConnectParam.UserName
					}

					if connectResource.EsConnectParam.Password != nil {
						esConnectParamMap["password"] = connectResource.EsConnectParam.Password
					}

					if connectResource.EsConnectParam.Resource != nil {
						esConnectParamMap["resource"] = connectResource.EsConnectParam.Resource
					}

					if connectResource.EsConnectParam.SelfBuilt != nil {
						esConnectParamMap["self_built"] = connectResource.EsConnectParam.SelfBuilt
					}

					if connectResource.EsConnectParam.ServiceVip != nil {
						esConnectParamMap["service_vip"] = connectResource.EsConnectParam.ServiceVip
					}

					if connectResource.EsConnectParam.UniqVpcId != nil {
						esConnectParamMap["uniq_vpc_id"] = connectResource.EsConnectParam.UniqVpcId
					}

					if connectResource.EsConnectParam.IsUpdate != nil {
						esConnectParamMap["is_update"] = connectResource.EsConnectParam.IsUpdate
					}

					connectResourceListMap["es_connect_param"] = []interface{}{esConnectParamMap}
				}

				if connectResource.ClickHouseConnectParam != nil {
					clickHouseConnectParamMap := map[string]interface{}{}

					if connectResource.ClickHouseConnectParam.Port != nil {
						clickHouseConnectParamMap["port"] = connectResource.ClickHouseConnectParam.Port
					}

					if connectResource.ClickHouseConnectParam.UserName != nil {
						clickHouseConnectParamMap["user_name"] = connectResource.ClickHouseConnectParam.UserName
					}

					if connectResource.ClickHouseConnectParam.Password != nil {
						clickHouseConnectParamMap["password"] = connectResource.ClickHouseConnectParam.Password
					}

					if connectResource.ClickHouseConnectParam.Resource != nil {
						clickHouseConnectParamMap["resource"] = connectResource.ClickHouseConnectParam.Resource
					}

					if connectResource.ClickHouseConnectParam.SelfBuilt != nil {
						clickHouseConnectParamMap["self_built"] = connectResource.ClickHouseConnectParam.SelfBuilt
					}

					if connectResource.ClickHouseConnectParam.ServiceVip != nil {
						clickHouseConnectParamMap["service_vip"] = connectResource.ClickHouseConnectParam.ServiceVip
					}

					if connectResource.ClickHouseConnectParam.UniqVpcId != nil {
						clickHouseConnectParamMap["uniq_vpc_id"] = connectResource.ClickHouseConnectParam.UniqVpcId
					}

					if connectResource.ClickHouseConnectParam.IsUpdate != nil {
						clickHouseConnectParamMap["is_update"] = connectResource.ClickHouseConnectParam.IsUpdate
					}

					connectResourceListMap["click_house_connect_param"] = []interface{}{clickHouseConnectParamMap}
				}

				if connectResource.MySQLConnectParam != nil {
					mySQLConnectParamMap := map[string]interface{}{}

					if connectResource.MySQLConnectParam.Port != nil {
						mySQLConnectParamMap["port"] = connectResource.MySQLConnectParam.Port
					}

					if connectResource.MySQLConnectParam.UserName != nil {
						mySQLConnectParamMap["user_name"] = connectResource.MySQLConnectParam.UserName
					}

					if connectResource.MySQLConnectParam.Password != nil {
						mySQLConnectParamMap["password"] = connectResource.MySQLConnectParam.Password
					}

					if connectResource.MySQLConnectParam.Resource != nil {
						mySQLConnectParamMap["resource"] = connectResource.MySQLConnectParam.Resource
					}

					if connectResource.MySQLConnectParam.ServiceVip != nil {
						mySQLConnectParamMap["service_vip"] = connectResource.MySQLConnectParam.ServiceVip
					}

					if connectResource.MySQLConnectParam.UniqVpcId != nil {
						mySQLConnectParamMap["uniq_vpc_id"] = connectResource.MySQLConnectParam.UniqVpcId
					}

					if connectResource.MySQLConnectParam.IsUpdate != nil {
						mySQLConnectParamMap["is_update"] = connectResource.MySQLConnectParam.IsUpdate
					}

					if connectResource.MySQLConnectParam.ClusterId != nil {
						mySQLConnectParamMap["cluster_id"] = connectResource.MySQLConnectParam.ClusterId
					}

					if connectResource.MySQLConnectParam.SelfBuilt != nil {
						mySQLConnectParamMap["self_built"] = connectResource.MySQLConnectParam.SelfBuilt
					}

					connectResourceListMap["mysql_connect_param"] = []interface{}{mySQLConnectParamMap}
				}

				if connectResource.PostgreSQLConnectParam != nil {
					postgreSQLConnectParamMap := map[string]interface{}{}

					if connectResource.PostgreSQLConnectParam.Port != nil {
						postgreSQLConnectParamMap["port"] = connectResource.PostgreSQLConnectParam.Port
					}

					if connectResource.PostgreSQLConnectParam.UserName != nil {
						postgreSQLConnectParamMap["user_name"] = connectResource.PostgreSQLConnectParam.UserName
					}

					if connectResource.PostgreSQLConnectParam.Password != nil {
						postgreSQLConnectParamMap["password"] = connectResource.PostgreSQLConnectParam.Password
					}

					if connectResource.PostgreSQLConnectParam.Resource != nil {
						postgreSQLConnectParamMap["resource"] = connectResource.PostgreSQLConnectParam.Resource
					}

					if connectResource.PostgreSQLConnectParam.ServiceVip != nil {
						postgreSQLConnectParamMap["service_vip"] = connectResource.PostgreSQLConnectParam.ServiceVip
					}

					if connectResource.PostgreSQLConnectParam.UniqVpcId != nil {
						postgreSQLConnectParamMap["uniq_vpc_id"] = connectResource.PostgreSQLConnectParam.UniqVpcId
					}

					if connectResource.PostgreSQLConnectParam.ClusterId != nil {
						postgreSQLConnectParamMap["cluster_id"] = connectResource.PostgreSQLConnectParam.ClusterId
					}

					if connectResource.PostgreSQLConnectParam.IsUpdate != nil {
						postgreSQLConnectParamMap["is_update"] = connectResource.PostgreSQLConnectParam.IsUpdate
					}

					if connectResource.PostgreSQLConnectParam.SelfBuilt != nil {
						postgreSQLConnectParamMap["self_built"] = connectResource.PostgreSQLConnectParam.SelfBuilt
					}

					connectResourceListMap["postgre_sql_connect_param"] = []interface{}{postgreSQLConnectParamMap}
				}

				if connectResource.MariaDBConnectParam != nil {
					mariaDBConnectParamMap := map[string]interface{}{}

					if connectResource.MariaDBConnectParam.Port != nil {
						mariaDBConnectParamMap["port"] = connectResource.MariaDBConnectParam.Port
					}

					if connectResource.MariaDBConnectParam.UserName != nil {
						mariaDBConnectParamMap["user_name"] = connectResource.MariaDBConnectParam.UserName
					}

					if connectResource.MariaDBConnectParam.Password != nil {
						mariaDBConnectParamMap["password"] = connectResource.MariaDBConnectParam.Password
					}

					if connectResource.MariaDBConnectParam.Resource != nil {
						mariaDBConnectParamMap["resource"] = connectResource.MariaDBConnectParam.Resource
					}

					if connectResource.MariaDBConnectParam.ServiceVip != nil {
						mariaDBConnectParamMap["service_vip"] = connectResource.MariaDBConnectParam.ServiceVip
					}

					if connectResource.MariaDBConnectParam.UniqVpcId != nil {
						mariaDBConnectParamMap["uniq_vpc_id"] = connectResource.MariaDBConnectParam.UniqVpcId
					}

					if connectResource.MariaDBConnectParam.IsUpdate != nil {
						mariaDBConnectParamMap["is_update"] = connectResource.MariaDBConnectParam.IsUpdate
					}

					connectResourceListMap["maria_db_connect_param"] = []interface{}{mariaDBConnectParamMap}
				}

				if connectResource.SQLServerConnectParam != nil {
					sQLServerConnectParamMap := map[string]interface{}{}

					if connectResource.SQLServerConnectParam.Port != nil {
						sQLServerConnectParamMap["port"] = connectResource.SQLServerConnectParam.Port
					}

					if connectResource.SQLServerConnectParam.UserName != nil {
						sQLServerConnectParamMap["user_name"] = connectResource.SQLServerConnectParam.UserName
					}

					if connectResource.SQLServerConnectParam.Password != nil {
						sQLServerConnectParamMap["password"] = connectResource.SQLServerConnectParam.Password
					}

					if connectResource.SQLServerConnectParam.Resource != nil {
						sQLServerConnectParamMap["resource"] = connectResource.SQLServerConnectParam.Resource
					}

					if connectResource.SQLServerConnectParam.ServiceVip != nil {
						sQLServerConnectParamMap["service_vip"] = connectResource.SQLServerConnectParam.ServiceVip
					}

					if connectResource.SQLServerConnectParam.UniqVpcId != nil {
						sQLServerConnectParamMap["uniq_vpc_id"] = connectResource.SQLServerConnectParam.UniqVpcId
					}

					if connectResource.SQLServerConnectParam.IsUpdate != nil {
						sQLServerConnectParamMap["is_update"] = connectResource.SQLServerConnectParam.IsUpdate
					}

					connectResourceListMap["sql_server_connect_param"] = []interface{}{sQLServerConnectParamMap}
				}

				if connectResource.CtsdbConnectParam != nil {
					ctsdbConnectParamMap := map[string]interface{}{}

					if connectResource.CtsdbConnectParam.Port != nil {
						ctsdbConnectParamMap["port"] = connectResource.CtsdbConnectParam.Port
					}

					if connectResource.CtsdbConnectParam.ServiceVip != nil {
						ctsdbConnectParamMap["service_vip"] = connectResource.CtsdbConnectParam.ServiceVip
					}

					if connectResource.CtsdbConnectParam.UniqVpcId != nil {
						ctsdbConnectParamMap["uniq_vpc_id"] = connectResource.CtsdbConnectParam.UniqVpcId
					}

					if connectResource.CtsdbConnectParam.UserName != nil {
						ctsdbConnectParamMap["user_name"] = connectResource.CtsdbConnectParam.UserName
					}

					if connectResource.CtsdbConnectParam.Password != nil {
						ctsdbConnectParamMap["password"] = connectResource.CtsdbConnectParam.Password
					}

					if connectResource.CtsdbConnectParam.Resource != nil {
						ctsdbConnectParamMap["resource"] = connectResource.CtsdbConnectParam.Resource
					}

					connectResourceListMap["ctsdb_connect_param"] = []interface{}{ctsdbConnectParamMap}
				}

				if connectResource.DorisConnectParam != nil {
					dorisConnectParamMap := map[string]interface{}{}

					if connectResource.DorisConnectParam.Port != nil {
						dorisConnectParamMap["port"] = connectResource.DorisConnectParam.Port
					}

					if connectResource.DorisConnectParam.UserName != nil {
						dorisConnectParamMap["user_name"] = connectResource.DorisConnectParam.UserName
					}

					if connectResource.DorisConnectParam.Password != nil {
						dorisConnectParamMap["password"] = connectResource.DorisConnectParam.Password
					}

					if connectResource.DorisConnectParam.Resource != nil {
						dorisConnectParamMap["resource"] = connectResource.DorisConnectParam.Resource
					}

					if connectResource.DorisConnectParam.ServiceVip != nil {
						dorisConnectParamMap["service_vip"] = connectResource.DorisConnectParam.ServiceVip
					}

					if connectResource.DorisConnectParam.UniqVpcId != nil {
						dorisConnectParamMap["uniq_vpc_id"] = connectResource.DorisConnectParam.UniqVpcId
					}

					if connectResource.DorisConnectParam.IsUpdate != nil {
						dorisConnectParamMap["is_update"] = connectResource.DorisConnectParam.IsUpdate
					}

					if connectResource.DorisConnectParam.SelfBuilt != nil {
						dorisConnectParamMap["self_built"] = connectResource.DorisConnectParam.SelfBuilt
					}

					if connectResource.DorisConnectParam.BePort != nil {
						dorisConnectParamMap["be_port"] = connectResource.DorisConnectParam.BePort
					}

					connectResourceListMap["doris_connect_param"] = []interface{}{dorisConnectParamMap}
				}

				if connectResource.KafkaConnectParam != nil {
					kafkaConnectParamMap := map[string]interface{}{}

					if connectResource.KafkaConnectParam.Resource != nil {
						kafkaConnectParamMap["resource"] = connectResource.KafkaConnectParam.Resource
					}

					if connectResource.KafkaConnectParam.SelfBuilt != nil {
						kafkaConnectParamMap["self_built"] = connectResource.KafkaConnectParam.SelfBuilt
					}

					if connectResource.KafkaConnectParam.IsUpdate != nil {
						kafkaConnectParamMap["is_update"] = connectResource.KafkaConnectParam.IsUpdate
					}

					if connectResource.KafkaConnectParam.BrokerAddress != nil {
						kafkaConnectParamMap["broker_address"] = connectResource.KafkaConnectParam.BrokerAddress
					}

					if connectResource.KafkaConnectParam.Region != nil {
						kafkaConnectParamMap["region"] = connectResource.KafkaConnectParam.Region
					}

					connectResourceListMap["kafka_connect_param"] = []interface{}{kafkaConnectParamMap}
				}

				connectResourceListList = append(connectResourceListList, connectResourceListMap)
				ids = append(ids, *connectResource.ResourceId)

			}

			describeConnectResourcesRespMap["connect_resource_list"] = connectResourceListList
		}

		_ = d.Set("result", []map[string]interface{}{describeConnectResourcesRespMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), describeConnectResourcesRespMap); e != nil {
			return e
		}
	}
	return nil
}
