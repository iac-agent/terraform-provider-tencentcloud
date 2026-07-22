package dlc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"
)

func DataSourceTencentCloudDlcDescribeDataEngine() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcDescribeDataEngineRead,
		Schema: map[string]*schema.Schema{
			"data_engine_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Engine 名称",
			},

			"data_engine": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Data 引擎 details。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_engine_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine 名称",
						},
						"engine_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine 类型: spark/presto。",
						},
						"cluster_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster 资源类型 spark_private/presto_private/presto_cu/spark_cu。",
						},
						"quota_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Quota ID。",
						},
						"state": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Data 引擎 状态 -2 删除，-1 failed，0 initializing，1 suspended，2 running，3 ready 到 delete，和 4 deleting。",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "创建时间。",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "更新时间。",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cluster specifications。",
						},
						"mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Billing 模式: 0 shared 模式，1 pay-作为-您-go，和 2 monthly subscription。",
						},
						"min_clusters": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最小clusters。",
						},
						"max_clusters": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最大clusters。",
						},
						"auto_resume": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否recover automatically。",
						},
						"spend_after": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Automatic recovery 时间。",
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster IP 范围。",
						},
						"default_data_engine": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否为the 默认值 引擎。",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Returned 消息",
						},
						"data_engine_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine 唯一 ID。",
						},
						"sub_account_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "操作者",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "过期时间。",
						},
						"isolated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Isolation 时间。",
						},
						"reversal_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rectification 时间。",
						},
						"user_alias": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "用户名",
						},
						"tag_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "标签列表",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签键",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签值",
									},
								},
							},
						},
						"permissions": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Engine permissions。",
						},
						"auto_suspend": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否automatically suspend 集群，prepay 不 support。",
						},
						"crontab_resume_suspend": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Engine crontab resume 或 suspend strategy，仅 support: 0: Wait(默认值)，1: Kill。",
						},
						"crontab_resume_suspend_strategy": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Engine auto suspend strategy，当 AutoSuspend 是 true，CrontabResumeSuspend 必须 stop。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resume_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Scheduled pull-up 时间: For 示例: 8 o&amp;#39;clock 在 Monday 是 expressed 作为 1000000-08:00:00。",
									},
									"suspend_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Scheduled suspension 时间: For 示例: 20 o&amp;#39;clock 在 Monday 是 expressed 作为 1000000-20:00:00。",
									},
									"suspend_strategy": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Suspend 配置: 0 (默认值): wait 对于 任务 到 end before suspending，1: force suspend。",
									},
								},
							},
						},
						"engine_exec_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine exec 类型，仅 support SQL(默认值) 或 BATCH。",
						},
						"renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Automatic renewal flag，0，initial state，automatic renewal 是 不 performed 通过 默认值. 如果 用户 has prepaid non-stop 服务 privileges，automatic renewal 将 occur. 1: Automatic renewal. 2: Make 它 clear 该 there 将 是 无 automatic renewal。",
						},
						"auto_suspend_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cluster automatic suspension 时间，默认值 10 minutes。",
						},
						"network_connection_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Network 连接 配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Network 配置 ID。",
									},
									"associate_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network 配置 唯一 identifier。",
									},
									"house_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data 引擎 ID。",
									},
									"datasource_connection_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data 来源 ID (obsolete)。",
									},
									"state": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Network 配置 状态 (0-initialization，1-normal)。",
									},
									"create_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "创建时间。",
									},
									"update_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "更新时间。",
									},
									"appid": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "用户 appid。",
									},
									"house_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data 引擎 名称",
									},
									"datasource_connection_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network 配置 名称",
									},
									"network_connection_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Network 配置 类型",
									},
									"uin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "用户 uin。",
									},
									"sub_account_uin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "用户 sub uin。",
									},
									"network_connection_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network 配置 描述",
									},
									"datasource_connection_vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource vpcid。",
									},
									"datasource_connection_subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource subnetId。",
									},
									"datasource_connection_cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource 连接 cidr block。",
									},
									"datasource_connection_subnet_cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Datasource 连接 子网 cidr block。",
									},
								},
							},
						},
						"ui_u_r_l": {
							Type:        schema.TypeString,
							Computed:    true,
							Deprecated:  "It has been deprecated. Use `ui_url` instead.",
							Description: "Jump 地址 的 ui。",
						},
						"ui_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Jump 地址 的 ui。",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine 资源类型 不 match，仅 support: Standard_CU/Memory_CU(仅 BATCH ExecType)。",
						},
						"image_version_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine major 版本 ID。",
						},
						"child_image_version_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine Image 版本 ID。",
						},
						"image_version_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine 镜像 版本 名称",
						},
						"start_standby_cluster": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否enable 备份 集群。",
						},
						"elastic_switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "For spark Batch ExecType，yearly 和 monthly 集群 是否enable elasticity。",
						},
						"elastic_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "For spark Batch ExecType，yearly 和 monthly 集群 elastic 限制",
						},
						"default_house": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Is 它 默认值 引擎?。",
						},
						"max_concurrency": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最大concurrent tasks 在 单个 集群，默认值 5。",
						},
						"tolerable_queue_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Tolerable queuing 时间，默认值 0. scaling 可能 是 triggered 当 tasks 是 queued 对于 longer 比 tolerable 时间. 如果 此 参数 是 0，它 表示 该 容量 expansion 可能 是 triggered immediately once 任务 是 queued。",
						},
						"user_app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "用户 appid。",
						},
						"user_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "用户 uin。",
						},
						"session_resource_template": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "For spark Batch ExecType，集群 会话 资源 配置 template。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"driver_size": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Engine 驱动 大小 规格 仅 支持: small/medium/large/xlarge/m.small/m.medium/m.large/m.xlarge。",
									},
									"executor_size": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Engine executor 大小 规格 仅 支持: small/medium/large/xlarge/m.small/m.medium/m.large/m.xlarge。",
									},
									"executor_nums": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "指定number 的 executors. 最小 值 是 1 和 最大 值 是 less 比 集群 规格。",
									},
									"executor_max_numbers": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "指定executor max 数量 (在 动态 配置 scenario)， 最小 值 是 1，和 最大 值 是 less 比 集群 规格 (当 ExecutorMaxNumbers 是 less 比 ExecutorNums， 值 是 集合 到 ExecutorNums)。",
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

func dataSourceTencentCloudDlcDescribeDataEngineRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_describe_data_engine.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var dataEngineName string
	if v, ok := d.GetOk("data_engine_name"); ok {
		dataEngineName = v.(string)
	}

	service := DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var dataEngine *dlc.DataEngineInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcDataEngineByName(ctx, dataEngineName)
		if e != nil {
			return tccommon.RetryError(e)
		}
		dataEngine = result
		return nil
	})
	if err != nil {
		return err
	}
	dataEngineInfoMap := map[string]interface{}{}

	if dataEngine != nil {

		if dataEngine.DataEngineName != nil {
			dataEngineInfoMap["data_engine_name"] = dataEngine.DataEngineName
		}

		if dataEngine.EngineType != nil {
			dataEngineInfoMap["engine_type"] = dataEngine.EngineType
		}

		if dataEngine.ClusterType != nil {
			dataEngineInfoMap["cluster_type"] = dataEngine.ClusterType
		}

		if dataEngine.QuotaId != nil {
			dataEngineInfoMap["quota_id"] = dataEngine.QuotaId
		}

		if dataEngine.State != nil {
			dataEngineInfoMap["state"] = dataEngine.State
		}

		if dataEngine.CreateTime != nil {
			dataEngineInfoMap["create_time"] = dataEngine.CreateTime
		}

		if dataEngine.UpdateTime != nil {
			dataEngineInfoMap["update_time"] = dataEngine.UpdateTime
		}

		if dataEngine.Size != nil {
			dataEngineInfoMap["size"] = dataEngine.Size
		}

		if dataEngine.Mode != nil {
			dataEngineInfoMap["mode"] = dataEngine.Mode
		}

		if dataEngine.MinClusters != nil {
			dataEngineInfoMap["min_clusters"] = dataEngine.MinClusters
		}

		if dataEngine.MaxClusters != nil {
			dataEngineInfoMap["max_clusters"] = dataEngine.MaxClusters
		}

		if dataEngine.AutoResume != nil {
			dataEngineInfoMap["auto_resume"] = dataEngine.AutoResume
		}

		if dataEngine.SpendAfter != nil {
			dataEngineInfoMap["spend_after"] = dataEngine.SpendAfter
		}

		if dataEngine.CidrBlock != nil {
			dataEngineInfoMap["cidr_block"] = dataEngine.CidrBlock
		}

		if dataEngine.DefaultDataEngine != nil {
			dataEngineInfoMap["default_data_engine"] = dataEngine.DefaultDataEngine
		}

		if dataEngine.Message != nil {
			dataEngineInfoMap["message"] = dataEngine.Message
		}

		if dataEngine.DataEngineId != nil {
			dataEngineInfoMap["data_engine_id"] = dataEngine.DataEngineId
		}

		if dataEngine.SubAccountUin != nil {
			dataEngineInfoMap["sub_account_uin"] = dataEngine.SubAccountUin
		}

		if dataEngine.ExpireTime != nil {
			dataEngineInfoMap["expire_time"] = dataEngine.ExpireTime
		}

		if dataEngine.IsolatedTime != nil {
			dataEngineInfoMap["isolated_time"] = dataEngine.IsolatedTime
		}

		if dataEngine.ReversalTime != nil {
			dataEngineInfoMap["reversal_time"] = dataEngine.ReversalTime
		}

		if dataEngine.UserAlias != nil {
			dataEngineInfoMap["user_alias"] = dataEngine.UserAlias
		}

		if dataEngine.TagList != nil {
			var tagListList []interface{}
			for _, tagList := range dataEngine.TagList {
				tagListMap := map[string]interface{}{}

				if tagList.TagKey != nil {
					tagListMap["tag_key"] = tagList.TagKey
				}

				if tagList.TagValue != nil {
					tagListMap["tag_value"] = tagList.TagValue
				}

				tagListList = append(tagListList, tagListMap)
			}

			dataEngineInfoMap["tag_list"] = tagListList
		}

		if dataEngine.Permissions != nil {
			dataEngineInfoMap["permissions"] = dataEngine.Permissions
		}

		if dataEngine.AutoSuspend != nil {
			dataEngineInfoMap["auto_suspend"] = dataEngine.AutoSuspend
		}

		if dataEngine.CrontabResumeSuspend != nil {
			dataEngineInfoMap["crontab_resume_suspend"] = dataEngine.CrontabResumeSuspend
		}

		if dataEngine.CrontabResumeSuspendStrategy != nil {
			crontabResumeSuspendStrategyMap := map[string]interface{}{}

			if dataEngine.CrontabResumeSuspendStrategy.ResumeTime != nil {
				crontabResumeSuspendStrategyMap["resume_time"] = dataEngine.CrontabResumeSuspendStrategy.ResumeTime
			}

			if dataEngine.CrontabResumeSuspendStrategy.SuspendTime != nil {
				crontabResumeSuspendStrategyMap["suspend_time"] = dataEngine.CrontabResumeSuspendStrategy.SuspendTime
			}

			if dataEngine.CrontabResumeSuspendStrategy.SuspendStrategy != nil {
				crontabResumeSuspendStrategyMap["suspend_strategy"] = dataEngine.CrontabResumeSuspendStrategy.SuspendStrategy
			}

			dataEngineInfoMap["crontab_resume_suspend_strategy"] = []interface{}{crontabResumeSuspendStrategyMap}
		}

		if dataEngine.EngineExecType != nil {
			dataEngineInfoMap["engine_exec_type"] = dataEngine.EngineExecType
		}

		if dataEngine.RenewFlag != nil {
			dataEngineInfoMap["renew_flag"] = dataEngine.RenewFlag
		}

		if dataEngine.AutoSuspendTime != nil {
			dataEngineInfoMap["auto_suspend_time"] = dataEngine.AutoSuspendTime
		}

		if dataEngine.NetworkConnectionSet != nil {
			var networkConnectionSetList []interface{}
			for _, networkConnectionSet := range dataEngine.NetworkConnectionSet {
				networkConnectionSetMap := map[string]interface{}{}

				if networkConnectionSet.Id != nil {
					networkConnectionSetMap["id"] = networkConnectionSet.Id
				}

				if networkConnectionSet.AssociateId != nil {
					networkConnectionSetMap["associate_id"] = networkConnectionSet.AssociateId
				}

				if networkConnectionSet.HouseId != nil {
					networkConnectionSetMap["house_id"] = networkConnectionSet.HouseId
				}

				if networkConnectionSet.DatasourceConnectionId != nil {
					networkConnectionSetMap["datasource_connection_id"] = networkConnectionSet.DatasourceConnectionId
				}

				if networkConnectionSet.State != nil {
					networkConnectionSetMap["state"] = networkConnectionSet.State
				}

				if networkConnectionSet.CreateTime != nil {
					networkConnectionSetMap["create_time"] = networkConnectionSet.CreateTime
				}

				if networkConnectionSet.UpdateTime != nil {
					networkConnectionSetMap["update_time"] = networkConnectionSet.UpdateTime
				}

				if networkConnectionSet.Appid != nil {
					networkConnectionSetMap["appid"] = networkConnectionSet.Appid
				}

				if networkConnectionSet.HouseName != nil {
					networkConnectionSetMap["house_name"] = networkConnectionSet.HouseName
				}

				if networkConnectionSet.DatasourceConnectionName != nil {
					networkConnectionSetMap["datasource_connection_name"] = networkConnectionSet.DatasourceConnectionName
				}

				if networkConnectionSet.NetworkConnectionType != nil {
					networkConnectionSetMap["network_connection_type"] = networkConnectionSet.NetworkConnectionType
				}

				if networkConnectionSet.Uin != nil {
					networkConnectionSetMap["uin"] = networkConnectionSet.Uin
				}

				if networkConnectionSet.SubAccountUin != nil {
					networkConnectionSetMap["sub_account_uin"] = networkConnectionSet.SubAccountUin
				}

				if networkConnectionSet.NetworkConnectionDesc != nil {
					networkConnectionSetMap["network_connection_desc"] = networkConnectionSet.NetworkConnectionDesc
				}

				if networkConnectionSet.DatasourceConnectionVpcId != nil {
					networkConnectionSetMap["datasource_connection_vpc_id"] = networkConnectionSet.DatasourceConnectionVpcId
				}

				if networkConnectionSet.DatasourceConnectionSubnetId != nil {
					networkConnectionSetMap["datasource_connection_subnet_id"] = networkConnectionSet.DatasourceConnectionSubnetId
				}

				if networkConnectionSet.DatasourceConnectionCidrBlock != nil {
					networkConnectionSetMap["datasource_connection_cidr_block"] = networkConnectionSet.DatasourceConnectionCidrBlock
				}

				if networkConnectionSet.DatasourceConnectionSubnetCidrBlock != nil {
					networkConnectionSetMap["datasource_connection_subnet_cidr_block"] = networkConnectionSet.DatasourceConnectionSubnetCidrBlock
				}

				networkConnectionSetList = append(networkConnectionSetList, networkConnectionSetMap)
			}

			dataEngineInfoMap["network_connection_set"] = networkConnectionSetList
		}

		if dataEngine.UiURL != nil {
			dataEngineInfoMap["ui_u_r_l"] = dataEngine.UiURL
			dataEngineInfoMap["ui_url"] = dataEngine.UiURL
		}

		if dataEngine.ResourceType != nil {
			dataEngineInfoMap["resource_type"] = dataEngine.ResourceType
		}

		if dataEngine.ImageVersionId != nil {
			dataEngineInfoMap["image_version_id"] = dataEngine.ImageVersionId
		}

		if dataEngine.ChildImageVersionId != nil {
			dataEngineInfoMap["child_image_version_id"] = dataEngine.ChildImageVersionId
		}

		if dataEngine.ImageVersionName != nil {
			dataEngineInfoMap["image_version_name"] = dataEngine.ImageVersionName
		}

		if dataEngine.StartStandbyCluster != nil {
			dataEngineInfoMap["start_standby_cluster"] = dataEngine.StartStandbyCluster
		}

		if dataEngine.ElasticSwitch != nil {
			dataEngineInfoMap["elastic_switch"] = dataEngine.ElasticSwitch
		}

		if dataEngine.ElasticLimit != nil {
			dataEngineInfoMap["elastic_limit"] = dataEngine.ElasticLimit
		}

		if dataEngine.DefaultHouse != nil {
			dataEngineInfoMap["default_house"] = dataEngine.DefaultHouse
		}

		if dataEngine.MaxConcurrency != nil {
			dataEngineInfoMap["max_concurrency"] = dataEngine.MaxConcurrency
		}

		if dataEngine.TolerableQueueTime != nil {
			dataEngineInfoMap["tolerable_queue_time"] = dataEngine.TolerableQueueTime
		}

		if dataEngine.UserAppId != nil {
			dataEngineInfoMap["user_app_id"] = dataEngine.UserAppId
		}

		if dataEngine.UserUin != nil {
			dataEngineInfoMap["user_uin"] = dataEngine.UserUin
		}

		if dataEngine.SessionResourceTemplate != nil {
			sessionResourceTemplateMap := map[string]interface{}{}

			if dataEngine.SessionResourceTemplate.DriverSize != nil {
				sessionResourceTemplateMap["driver_size"] = dataEngine.SessionResourceTemplate.DriverSize
			}

			if dataEngine.SessionResourceTemplate.ExecutorSize != nil {
				sessionResourceTemplateMap["executor_size"] = dataEngine.SessionResourceTemplate.ExecutorSize
			}

			if dataEngine.SessionResourceTemplate.ExecutorNums != nil {
				sessionResourceTemplateMap["executor_nums"] = dataEngine.SessionResourceTemplate.ExecutorNums
			}

			if dataEngine.SessionResourceTemplate.ExecutorMaxNumbers != nil {
				sessionResourceTemplateMap["executor_max_numbers"] = dataEngine.SessionResourceTemplate.ExecutorMaxNumbers
			}

			dataEngineInfoMap["session_resource_template"] = []interface{}{sessionResourceTemplateMap}
		}

		_ = d.Set("data_engine", []interface{}{dataEngineInfoMap})
	}

	d.SetId(dataEngineName)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataEngineInfoMap); e != nil {
			return e
		}
	}
	return nil
}
