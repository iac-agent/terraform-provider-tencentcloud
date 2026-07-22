package postgresql

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPostgresqlReadonlyGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresqlReadonlyGroupsRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 condition. primary ID 必须 是 指定 在 格式 的 db-master-实例-ID 到 过滤器 results，或 else null 将 是 返回。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "过滤名称",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "一个或多个过滤值",
						},
					},
				},
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting criterion. 有效 值:ROGroupId，CreateTime，名称",
			},

			"order_by_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting 顺序 有效 值:desc，asc。",
			},

			"read_only_group_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 read-仅 groups。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"read_only_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "read-仅 组 id注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"read_only_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "read-仅 组 name注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "项目 id注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"master_db_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "master 实例 id注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"min_delay_eliminate_reserve": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum 数量 Reserved Instances注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"max_replay_latency": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "延迟 space 大小 阈值。",
						},
						"replay_latency_eliminate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "延迟 大小 switch。",
						},
						"max_replay_lag": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "延迟 时间 大小 阈值。",
						},
						"replay_lag_eliminate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "延迟 时间 switch。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "virtual 网络 ID。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "子网-id注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 ID",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 ID",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "state。",
						},
						"read_only_db_instance_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "实例 details。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "地域 到 其中 实例 belongs，such 作为: ap-guangzhou，corresponding 到 地域 字段 的 RegionSet。",
									},
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Availability 可用区 到 其中 实例 belongs，such 作为: ap-guangzhou-3，corresponding 到 可用区 字段 的 ZoneSet。",
									},
									"project_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "项目 ID。",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "私有 网络 ID。",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "子网 ID。",
									},
									"db_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例 ID。",
									},
									"db_instance_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例名称",
									},
									"db_instance_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例状态，respectively: applying (applying)，init (到 是 initialized)，initing (initializing)，running (running)，limited run (limited run)，isolated (isolated)，recycling (recycling )，recycled (recycled)，作业 running (任务 execution)，offline (offline)，migrating (迁移)，expanding (expanding)，waitSwitch (waiting 对于 switching)，switching (switching)，readonly (read-仅 )，restarting (restarting)，网络 changing (网络 changing)，upgrading (kernel 版本 upgrade)。",
									},
									"db_instance_memory": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "内存 大小 allocated 通过 实例，单位: GB。",
									},
									"db_instance_storage": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "大小 的 存储 space allocated 通过 实例，单位: GB。",
									},
									"db_instance_cpu": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 CPUs allocated 通过 实例。",
									},
									"db_instance_class": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "sales 规格 ID。",
									},
									"db_instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例类型， types 是: 1. primary (primary 实例); 2. readonly (read-仅 实例); 3. guard (disaster recovery 实例); 4. temp (temporary 实例)。",
									},
									"db_instance_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例 版本，currently 仅 支持 standard (dual machine high availability 版本，一个 master 和 一个 slave)。",
									},
									"db_charset": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例 DB character 集合。",
									},
									"db_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "PostgreSQL 版本",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例 创建时间。",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "时间 当 实例 performed last update。",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例 过期时间。",
									},
									"isolated_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例 isolation 时间。",
									},
									"pay_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "billing 模式，1. prepaid (subscription，prepaid); 2. postpaid (billing 通过 卷，postpaid)。",
									},
									"auto_renew": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "auto-renew，1: auto-renew，0: 无 auto-renew。",
									},
									"db_instance_net_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "实例 网络 连接 信息。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "DNS 域名 名称",
												},
												"ip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "IP 地址",
												},
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "连接 端口 地址",
												},
												"net_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "网络 类型，1. inner (intranet 地址 的 basic 网络); 2. 私有 (intranet 地址 的 私有 网络); 3. 公有 (extranet 地址 的 basic 网络 或 私有 网络);。",
												},
												"status": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "网络 连接 状态，1. initing (unopened); 2. opened (opened); 3. closed (closed); 4. opening (opening); 5. closing (closed);。",
												},
												"vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "私有 网络 ID注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "子网 ID注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"protocol_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "协议 类型 对于 connecting 到 数据库，currently 支持: postgresql，mssql (MSSQL compatible syntax)注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "machine 类型",
									},
									"app_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "用户&#39;s AppId。",
									},
									"uid": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Uid 的 实例。",
									},
									"support_ipv6": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "是否instance 支持 Ipv6，1: support，0: 不 support。",
									},
									"tag_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "标签 信息 bound 到 instance注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"tag_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "标签 键",
												},
												"tag_value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "标签值",
												},
											},
										},
									},
									"master_db_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Master 实例 信息，仅 返回 当 实例 是 read-only注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"read_only_instance_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 read-仅 instances注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"status_in_readonly_group": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "状态 read-仅 实例 在 read-仅 group注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"offline_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "offline time注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"db_kernel_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Database kernel version注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"network_access_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "实例 网络 信息 列表 (此 字段 是 obsolete)注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"resource_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Network 资源 ID，实例 ID 或 RO 组 id注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"resource_type": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "资源类型，1-实例 2-RO group注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "私有 网络 ID注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"vip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "IPV4 address注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"vip6": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "IPV6 address注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"vport": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "访问 port注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "子网 ID注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"vpc_status": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Network 状态，1-applying，2-使用，3-deleting，4-deleted注意：此字段可能返回 null，表示无法获取有效值。",
												},
											},
										},
									},
									"db_major_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "PostgreSQL major version注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"db_node_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "实例 节点 information注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"role": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Node 类型， 值 可以 是:Primary，representing primary 节点;Standby，stands 对于 standby 节点。",
												},
												"zone": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Availability 可用区 其中 节点 是 located，such 作为 ap-guangzhou-1。",
												},
											},
										},
									},
									"is_support_t_d_e": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "是否instance 支持 TDE 数据 加密 0: 不 支持，1: supported注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"db_engine": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Database 引擎 该 支持:1. postgresql (云 数据库 PostgreSQL);2. mssql_compatible (MSSQL compatible - 云 数据库 PostgreSQL);注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"db_engine_config": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration 信息 对于 数据库 engine注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"rebalance": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "automatic load balancing switch。",
						},
						"db_instance_net_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "网络 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DNS 域名 名称",
									},
									"ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IP 地址",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "连接 端口 地址",
									},
									"net_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "网络 类型，1. inner (intranet 地址 的 basic 网络); 2. 私有 (intranet 地址 的 私有 网络); 3. 公有 (extranet 地址 的 basic 网络 或 私有 网络);。",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "网络 连接 状态，1. initing (unopened); 2. opened (opened); 3. closed (closed); 4. opening (opening); 5. closing (closed);。",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "私有 网络 ID注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "子网 ID注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"protocol_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "协议 类型 对于 connecting 到 数据库，currently 支持: postgresql，mssql (MSSQL compatible syntax)注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"network_access_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Read-仅 列表 组 网络 信息 (此 字段 是 obsolete)注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network 资源 ID，实例 ID 或 RO 组 id注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"resource_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "资源类型，1-实例 2-RO group注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "私有 网络 ID注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"vip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IPV4 address注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"vip6": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IPV6 address注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"vport": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "访问 port注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "子网 ID注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"vpc_status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Network 状态，1-applying，2-使用，3-deleting，4-deleted注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudPostgresqlReadonlyGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_postgresql_readonly_groups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*postgresql.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := postgresql.Filter{}
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

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by_type"); ok {
		paramMap["OrderByType"] = helper.String(v.(string))
	}

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var readOnlyGroupList []*postgresql.ReadOnlyGroup

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePostgresqlReadonlyGroupsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		readOnlyGroupList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(readOnlyGroupList))
	tmpList := make([]map[string]interface{}, 0, len(readOnlyGroupList))

	if readOnlyGroupList != nil {
		for _, readOnlyGroup := range readOnlyGroupList {
			readOnlyGroupMap := map[string]interface{}{}

			if readOnlyGroup.ReadOnlyGroupId != nil {
				readOnlyGroupMap["read_only_group_id"] = readOnlyGroup.ReadOnlyGroupId
			}

			if readOnlyGroup.ReadOnlyGroupName != nil {
				readOnlyGroupMap["read_only_group_name"] = readOnlyGroup.ReadOnlyGroupName
			}

			if readOnlyGroup.ProjectId != nil {
				readOnlyGroupMap["project_id"] = readOnlyGroup.ProjectId
			}

			if readOnlyGroup.MasterDBInstanceId != nil {
				readOnlyGroupMap["master_db_instance_id"] = readOnlyGroup.MasterDBInstanceId
			}

			if readOnlyGroup.MinDelayEliminateReserve != nil {
				readOnlyGroupMap["min_delay_eliminate_reserve"] = readOnlyGroup.MinDelayEliminateReserve
			}

			if readOnlyGroup.MaxReplayLatency != nil {
				readOnlyGroupMap["max_replay_latency"] = readOnlyGroup.MaxReplayLatency
			}

			if readOnlyGroup.ReplayLatencyEliminate != nil {
				readOnlyGroupMap["replay_latency_eliminate"] = readOnlyGroup.ReplayLatencyEliminate
			}

			if readOnlyGroup.MaxReplayLag != nil {
				readOnlyGroupMap["max_replay_lag"] = readOnlyGroup.MaxReplayLag
			}

			if readOnlyGroup.ReplayLagEliminate != nil {
				readOnlyGroupMap["replay_lag_eliminate"] = readOnlyGroup.ReplayLagEliminate
			}

			if readOnlyGroup.VpcId != nil {
				readOnlyGroupMap["vpc_id"] = readOnlyGroup.VpcId
			}

			if readOnlyGroup.SubnetId != nil {
				readOnlyGroupMap["subnet_id"] = readOnlyGroup.SubnetId
			}

			if readOnlyGroup.Region != nil {
				readOnlyGroupMap["region"] = readOnlyGroup.Region
			}

			if readOnlyGroup.Zone != nil {
				readOnlyGroupMap["zone"] = readOnlyGroup.Zone
			}

			if readOnlyGroup.Status != nil {
				readOnlyGroupMap["status"] = readOnlyGroup.Status
			}

			if readOnlyGroup.ReadOnlyDBInstanceList != nil {
				readOnlyDbInstanceListList := []interface{}{}
				for _, readOnlyDbInstanceList := range readOnlyGroup.ReadOnlyDBInstanceList {
					readOnlyDbInstanceListMap := map[string]interface{}{}

					if readOnlyDbInstanceList.Region != nil {
						readOnlyDbInstanceListMap["region"] = readOnlyDbInstanceList.Region
					}

					if readOnlyDbInstanceList.Zone != nil {
						readOnlyDbInstanceListMap["zone"] = readOnlyDbInstanceList.Zone
					}

					if readOnlyDbInstanceList.ProjectId != nil {
						readOnlyDbInstanceListMap["project_id"] = readOnlyDbInstanceList.ProjectId
					}

					if readOnlyDbInstanceList.VpcId != nil {
						readOnlyDbInstanceListMap["vpc_id"] = readOnlyDbInstanceList.VpcId
					}

					if readOnlyDbInstanceList.SubnetId != nil {
						readOnlyDbInstanceListMap["subnet_id"] = readOnlyDbInstanceList.SubnetId
					}

					if readOnlyDbInstanceList.DBInstanceId != nil {
						readOnlyDbInstanceListMap["db_instance_id"] = readOnlyDbInstanceList.DBInstanceId
					}

					if readOnlyDbInstanceList.DBInstanceName != nil {
						readOnlyDbInstanceListMap["db_instance_name"] = readOnlyDbInstanceList.DBInstanceName
					}

					if readOnlyDbInstanceList.DBInstanceStatus != nil {
						readOnlyDbInstanceListMap["db_instance_status"] = readOnlyDbInstanceList.DBInstanceStatus
					}

					if readOnlyDbInstanceList.DBInstanceMemory != nil {
						readOnlyDbInstanceListMap["db_instance_memory"] = readOnlyDbInstanceList.DBInstanceMemory
					}

					if readOnlyDbInstanceList.DBInstanceStorage != nil {
						readOnlyDbInstanceListMap["db_instance_storage"] = readOnlyDbInstanceList.DBInstanceStorage
					}

					if readOnlyDbInstanceList.DBInstanceCpu != nil {
						readOnlyDbInstanceListMap["db_instance_cpu"] = readOnlyDbInstanceList.DBInstanceCpu
					}

					if readOnlyDbInstanceList.DBInstanceClass != nil {
						readOnlyDbInstanceListMap["db_instance_class"] = readOnlyDbInstanceList.DBInstanceClass
					}

					if readOnlyDbInstanceList.DBInstanceType != nil {
						readOnlyDbInstanceListMap["db_instance_type"] = readOnlyDbInstanceList.DBInstanceType
					}

					if readOnlyDbInstanceList.DBInstanceVersion != nil {
						readOnlyDbInstanceListMap["db_instance_version"] = readOnlyDbInstanceList.DBInstanceVersion
					}

					if readOnlyDbInstanceList.DBCharset != nil {
						readOnlyDbInstanceListMap["db_charset"] = readOnlyDbInstanceList.DBCharset
					}

					if readOnlyDbInstanceList.DBVersion != nil {
						readOnlyDbInstanceListMap["db_version"] = readOnlyDbInstanceList.DBVersion
					}

					if readOnlyDbInstanceList.CreateTime != nil {
						readOnlyDbInstanceListMap["create_time"] = readOnlyDbInstanceList.CreateTime
					}

					if readOnlyDbInstanceList.UpdateTime != nil {
						readOnlyDbInstanceListMap["update_time"] = readOnlyDbInstanceList.UpdateTime
					}

					if readOnlyDbInstanceList.ExpireTime != nil {
						readOnlyDbInstanceListMap["expire_time"] = readOnlyDbInstanceList.ExpireTime
					}

					if readOnlyDbInstanceList.IsolatedTime != nil {
						readOnlyDbInstanceListMap["isolated_time"] = readOnlyDbInstanceList.IsolatedTime
					}

					if readOnlyDbInstanceList.PayType != nil {
						readOnlyDbInstanceListMap["pay_type"] = readOnlyDbInstanceList.PayType
					}

					if readOnlyDbInstanceList.AutoRenew != nil {
						readOnlyDbInstanceListMap["auto_renew"] = readOnlyDbInstanceList.AutoRenew
					}

					if readOnlyDbInstanceList.DBInstanceNetInfo != nil {
						dbInstanceNetInfoList := []interface{}{}
						for _, dbInstanceNetInfo := range readOnlyDbInstanceList.DBInstanceNetInfo {
							dbInstanceNetInfoMap := map[string]interface{}{}

							if dbInstanceNetInfo.Address != nil {
								dbInstanceNetInfoMap["address"] = dbInstanceNetInfo.Address
							}

							if dbInstanceNetInfo.Ip != nil {
								dbInstanceNetInfoMap["ip"] = dbInstanceNetInfo.Ip
							}

							if dbInstanceNetInfo.Port != nil {
								dbInstanceNetInfoMap["port"] = dbInstanceNetInfo.Port
							}

							if dbInstanceNetInfo.NetType != nil {
								dbInstanceNetInfoMap["net_type"] = dbInstanceNetInfo.NetType
							}

							if dbInstanceNetInfo.Status != nil {
								dbInstanceNetInfoMap["status"] = dbInstanceNetInfo.Status
							}

							if dbInstanceNetInfo.VpcId != nil {
								dbInstanceNetInfoMap["vpc_id"] = dbInstanceNetInfo.VpcId
							}

							if dbInstanceNetInfo.SubnetId != nil {
								dbInstanceNetInfoMap["subnet_id"] = dbInstanceNetInfo.SubnetId
							}

							if dbInstanceNetInfo.ProtocolType != nil {
								dbInstanceNetInfoMap["protocol_type"] = dbInstanceNetInfo.ProtocolType
							}

							dbInstanceNetInfoList = append(dbInstanceNetInfoList, dbInstanceNetInfoMap)
						}

						readOnlyDbInstanceListMap["db_instance_net_info"] = dbInstanceNetInfoList
					}

					if readOnlyDbInstanceList.Type != nil {
						readOnlyDbInstanceListMap["type"] = readOnlyDbInstanceList.Type
					}

					if readOnlyDbInstanceList.AppId != nil {
						readOnlyDbInstanceListMap["app_id"] = readOnlyDbInstanceList.AppId
					}

					if readOnlyDbInstanceList.Uid != nil {
						readOnlyDbInstanceListMap["uid"] = readOnlyDbInstanceList.Uid
					}

					if readOnlyDbInstanceList.SupportIpv6 != nil {
						readOnlyDbInstanceListMap["support_ipv6"] = readOnlyDbInstanceList.SupportIpv6
					}

					if readOnlyDbInstanceList.TagList != nil {
						tagListList := []interface{}{}
						for _, tagList := range readOnlyDbInstanceList.TagList {
							tagListMap := map[string]interface{}{}

							if tagList.TagKey != nil {
								tagListMap["tag_key"] = tagList.TagKey
							}

							if tagList.TagValue != nil {
								tagListMap["tag_value"] = tagList.TagValue
							}

							tagListList = append(tagListList, tagListMap)
						}

						readOnlyDbInstanceListMap["tag_list"] = tagListList
					}

					if readOnlyDbInstanceList.MasterDBInstanceId != nil {
						readOnlyDbInstanceListMap["master_db_instance_id"] = readOnlyDbInstanceList.MasterDBInstanceId
					}

					if readOnlyDbInstanceList.ReadOnlyInstanceNum != nil {
						readOnlyDbInstanceListMap["read_only_instance_num"] = readOnlyDbInstanceList.ReadOnlyInstanceNum
					}

					if readOnlyDbInstanceList.StatusInReadonlyGroup != nil {
						readOnlyDbInstanceListMap["status_in_readonly_group"] = readOnlyDbInstanceList.StatusInReadonlyGroup
					}

					if readOnlyDbInstanceList.OfflineTime != nil {
						readOnlyDbInstanceListMap["offline_time"] = readOnlyDbInstanceList.OfflineTime
					}

					if readOnlyDbInstanceList.DBKernelVersion != nil {
						readOnlyDbInstanceListMap["db_kernel_version"] = readOnlyDbInstanceList.DBKernelVersion
					}

					if readOnlyDbInstanceList.NetworkAccessList != nil {
						networkAccessListList := []interface{}{}
						for _, networkAccessList := range readOnlyDbInstanceList.NetworkAccessList {
							networkAccessListMap := map[string]interface{}{}

							if networkAccessList.ResourceId != nil {
								networkAccessListMap["resource_id"] = networkAccessList.ResourceId
							}

							if networkAccessList.ResourceType != nil {
								networkAccessListMap["resource_type"] = networkAccessList.ResourceType
							}

							if networkAccessList.VpcId != nil {
								networkAccessListMap["vpc_id"] = networkAccessList.VpcId
							}

							if networkAccessList.Vip != nil {
								networkAccessListMap["vip"] = networkAccessList.Vip
							}

							if networkAccessList.Vip6 != nil {
								networkAccessListMap["vip6"] = networkAccessList.Vip6
							}

							if networkAccessList.Vport != nil {
								networkAccessListMap["vport"] = networkAccessList.Vport
							}

							if networkAccessList.SubnetId != nil {
								networkAccessListMap["subnet_id"] = networkAccessList.SubnetId
							}

							if networkAccessList.VpcStatus != nil {
								networkAccessListMap["vpc_status"] = networkAccessList.VpcStatus
							}

							networkAccessListList = append(networkAccessListList, networkAccessListMap)
						}

						readOnlyDbInstanceListMap["network_access_list"] = networkAccessListList
					}

					if readOnlyDbInstanceList.DBMajorVersion != nil {
						readOnlyDbInstanceListMap["db_major_version"] = readOnlyDbInstanceList.DBMajorVersion
					}

					if readOnlyDbInstanceList.DBNodeSet != nil {
						dbNodeSetList := []interface{}{}
						for _, dbNodeSet := range readOnlyDbInstanceList.DBNodeSet {
							dbNodeSetMap := map[string]interface{}{}

							if dbNodeSet.Role != nil {
								dbNodeSetMap["role"] = dbNodeSet.Role
							}

							if dbNodeSet.Zone != nil {
								dbNodeSetMap["zone"] = dbNodeSet.Zone
							}

							dbNodeSetList = append(dbNodeSetList, dbNodeSetMap)
						}

						readOnlyDbInstanceListMap["db_node_set"] = dbNodeSetList
					}

					if readOnlyDbInstanceList.IsSupportTDE != nil {
						readOnlyDbInstanceListMap["is_support_t_d_e"] = readOnlyDbInstanceList.IsSupportTDE
					}

					if readOnlyDbInstanceList.DBEngine != nil {
						readOnlyDbInstanceListMap["db_engine"] = readOnlyDbInstanceList.DBEngine
					}

					if readOnlyDbInstanceList.DBEngineConfig != nil {
						readOnlyDbInstanceListMap["db_engine_config"] = readOnlyDbInstanceList.DBEngineConfig
					}

					readOnlyDbInstanceListList = append(readOnlyDbInstanceListList, readOnlyDbInstanceListMap)
				}

				readOnlyGroupMap["read_only_db_instance_list"] = readOnlyDbInstanceListList
			}

			if readOnlyGroup.Rebalance != nil {
				readOnlyGroupMap["rebalance"] = readOnlyGroup.Rebalance
			}

			if readOnlyGroup.DBInstanceNetInfo != nil {
				dbInstanceNetInfoList := []interface{}{}
				for _, dbInstanceNetInfo := range readOnlyGroup.DBInstanceNetInfo {
					dbInstanceNetInfoMap := map[string]interface{}{}

					if dbInstanceNetInfo.Address != nil {
						dbInstanceNetInfoMap["address"] = dbInstanceNetInfo.Address
					}

					if dbInstanceNetInfo.Ip != nil {
						dbInstanceNetInfoMap["ip"] = dbInstanceNetInfo.Ip
					}

					if dbInstanceNetInfo.Port != nil {
						dbInstanceNetInfoMap["port"] = dbInstanceNetInfo.Port
					}

					if dbInstanceNetInfo.NetType != nil {
						dbInstanceNetInfoMap["net_type"] = dbInstanceNetInfo.NetType
					}

					if dbInstanceNetInfo.Status != nil {
						dbInstanceNetInfoMap["status"] = dbInstanceNetInfo.Status
					}

					if dbInstanceNetInfo.VpcId != nil {
						dbInstanceNetInfoMap["vpc_id"] = dbInstanceNetInfo.VpcId
					}

					if dbInstanceNetInfo.SubnetId != nil {
						dbInstanceNetInfoMap["subnet_id"] = dbInstanceNetInfo.SubnetId
					}

					if dbInstanceNetInfo.ProtocolType != nil {
						dbInstanceNetInfoMap["protocol_type"] = dbInstanceNetInfo.ProtocolType
					}

					dbInstanceNetInfoList = append(dbInstanceNetInfoList, dbInstanceNetInfoMap)
				}

				readOnlyGroupMap["db_instance_net_info"] = dbInstanceNetInfoList
			}

			if readOnlyGroup.NetworkAccessList != nil {
				networkAccessListList := []interface{}{}
				for _, networkAccessList := range readOnlyGroup.NetworkAccessList {
					networkAccessListMap := map[string]interface{}{}

					if networkAccessList.ResourceId != nil {
						networkAccessListMap["resource_id"] = networkAccessList.ResourceId
					}

					if networkAccessList.ResourceType != nil {
						networkAccessListMap["resource_type"] = networkAccessList.ResourceType
					}

					if networkAccessList.VpcId != nil {
						networkAccessListMap["vpc_id"] = networkAccessList.VpcId
					}

					if networkAccessList.Vip != nil {
						networkAccessListMap["vip"] = networkAccessList.Vip
					}

					if networkAccessList.Vip6 != nil {
						networkAccessListMap["vip6"] = networkAccessList.Vip6
					}

					if networkAccessList.Vport != nil {
						networkAccessListMap["vport"] = networkAccessList.Vport
					}

					if networkAccessList.SubnetId != nil {
						networkAccessListMap["subnet_id"] = networkAccessList.SubnetId
					}

					if networkAccessList.VpcStatus != nil {
						networkAccessListMap["vpc_status"] = networkAccessList.VpcStatus
					}

					networkAccessListList = append(networkAccessListList, networkAccessListMap)
				}

				readOnlyGroupMap["network_access_list"] = networkAccessListList
			}

			ids = append(ids, *readOnlyGroup.ReadOnlyGroupId)
			tmpList = append(tmpList, readOnlyGroupMap)
		}

		_ = d.Set("read_only_group_list", tmpList)
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
