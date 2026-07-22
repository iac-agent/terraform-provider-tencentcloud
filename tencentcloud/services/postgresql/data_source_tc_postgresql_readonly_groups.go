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
				Description: "Filter condition. The primary ID must be specified in the 格式 of db-master-instance-id to filter results，or else null will be returned。",
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
				Description: "Sorting criterion. Valid values:ROGroupId，CreateTime，名称",
			},

			"order_by_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting 顺序 Valid values:desc，asc。",
			},

			"read_only_group_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 read-only groups。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"read_only_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "read-only group id注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"read_only_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "read-only group name注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "project id注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"master_db_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "master instance id注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"min_delay_eliminate_reserve": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum 数量 Reserved Instances注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"max_replay_latency": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "延迟 space size threshold。",
						},
						"replay_latency_eliminate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "延迟 size switch。",
						},
						"max_replay_lag": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "延迟 time size threshold。",
						},
						"replay_lag_eliminate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "延迟 time switch。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "virtual network id。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "subnet-id注意：此字段可能返回 null，表示无法获取有效值。",
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
							Description: "instance details。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The 地域 to which the instance belongs，such as: ap-guangzhou，corresponding to the 地域 field of the RegionSet。",
									},
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Availability 可用区 to which the instance belongs，such as: ap-guangzhou-3，corresponding to the 可用区 field of ZoneSet。",
									},
									"project_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "project ID。",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "private network ID。",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "subnet ID。",
									},
									"db_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "instance ID。",
									},
									"db_instance_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例名称",
									},
									"db_instance_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例状态，respectively: applying (applying)，init (to be initialized)，initing (initializing)，running (running)，limited run (limited run)，isolated (isolated)，recycling (recycling )，recycled (recycled)，job running (task execution)，offline (offline)，migrating (migration)，expanding (expanding)，waitSwitch (waiting for switching)，switching (switching)，readonly (read-only )，restarting (restarting)，network changing (network changing)，upgrading (kernel 版本 upgrade)。",
									},
									"db_instance_memory": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "the memory size allocated by the instance，unit: GB。",
									},
									"db_instance_storage": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "the size of the storage space allocated by the instance，unit: GB。",
									},
									"db_instance_cpu": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "the 数量 CPUs allocated by the instance。",
									},
									"db_instance_class": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "sales specification ID。",
									},
									"db_instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例类型，the types are: 1. primary (primary instance); 2. readonly (read-only instance); 3. guard (disaster recovery instance); 4. temp (temporary instance)。",
									},
									"db_instance_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "instance 版本，currently only supports standard (dual machine high availability 版本，one master and one slave)。",
									},
									"db_charset": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "instance DB character set。",
									},
									"db_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "PostgreSQL 版本",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "instance 创建时间。",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The time when the instance performed the last update。",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "instance 过期时间。",
									},
									"isolated_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "instance isolation time。",
									},
									"pay_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "billing 模式，1. prepaid (subscription，prepaid); 2. postpaid (billing by volume，postpaid)。",
									},
									"auto_renew": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "auto-renew，1: auto-renew，0: no auto-renew。",
									},
									"db_instance_net_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "instance network connection information。",
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
													Description: "connection 端口 地址",
												},
												"net_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "network 类型，1. inner (intranet 地址 of the basic network); 2. private (intranet 地址 of the private network); 3. public (extranet 地址 of the basic network or private network);。",
												},
												"status": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "network connection 状态，1. initing (unopened); 2. opened (opened); 3. closed (closed); 4. opening (opening); 5. closing (closed);。",
												},
												"vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "private network ID注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "subnet ID注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"protocol_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The 协议 类型 for connecting to the database，currently supported: postgresql，mssql (MSSQL compatible syntax)注意：此字段可能返回 null，表示无法获取有效值。",
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
										Description: "Uid of the instance。",
									},
									"support_ipv6": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "是否instance supports Ipv6，1: support，0: not support。",
									},
									"tag_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "标签 information bound to the instance注意：此字段可能返回 null，表示无法获取有效值。",
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
										Description: "Master instance information，only returned when the instance is read-only注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"read_only_instance_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 read-only instances注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"status_in_readonly_group": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "状态 read-only instance in the read-only group注意：此字段可能返回 null，表示无法获取有效值。",
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
										Description: "Instance network information list (this field is obsolete)注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"resource_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Network 资源 ID，实例 ID or RO group id注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"resource_type": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "资源类型，1-instance 2-RO group注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "private network ID注意：此字段可能返回 null，表示无法获取有效值。",
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
													Description: "access port注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "subnet ID注意：此字段可能返回 null，表示无法获取有效值。",
												},
												"vpc_status": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Network 状态，1-applying，2-using，3-deleting，4-deleted注意：此字段可能返回 null，表示无法获取有效值。",
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
										Description: "Instance node information注意：此字段可能返回 null，表示无法获取有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"role": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Node 类型，the 值 can be:Primary，representing the primary node;Standby，stands for standby node。",
												},
												"zone": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Availability 可用区 where the node is located，such as ap-guangzhou-1。",
												},
											},
										},
									},
									"is_support_t_d_e": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "是否instance supports TDE data encryption 0: not supported，1: supported注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"db_engine": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Database engine that supports:1. postgresql (cloud database PostgreSQL);2. mssql_compatible (MSSQL compatible - cloud database PostgreSQL);注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"db_engine_config": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration information for the database engine注意：此字段可能返回 null，表示无法获取有效值。",
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
							Description: "network information。",
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
										Description: "connection 端口 地址",
									},
									"net_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "network 类型，1. inner (intranet 地址 of the basic network); 2. private (intranet 地址 of the private network); 3. public (extranet 地址 of the basic network or private network);。",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "network connection 状态，1. initing (unopened); 2. opened (opened); 3. closed (closed); 4. opening (opening); 5. closing (closed);。",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "private network ID注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "subnet ID注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"protocol_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The 协议 类型 for connecting to the database，currently supported: postgresql，mssql (MSSQL compatible syntax)注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"network_access_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Read-only 列表 group network information (this field is obsolete)注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network 资源 ID，实例 ID or RO group id注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"resource_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "资源类型，1-instance 2-RO group注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "private network ID注意：此字段可能返回 null，表示无法获取有效值。",
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
										Description: "access port注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "subnet ID注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"vpc_status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Network 状态，1-applying，2-using，3-deleting，4-deleted注意：此字段可能返回 null，表示无法获取有效值。",
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
