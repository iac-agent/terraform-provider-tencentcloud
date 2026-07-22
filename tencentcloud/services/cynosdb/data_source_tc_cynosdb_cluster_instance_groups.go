package cynosdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCynosdbClusterInstanceGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbClusterInstanceGroupsRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "集群的ID。",
			},

			"instance_grp_info_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "实例组列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "应用程序 ID。",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "集群的ID。",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创造时间。",
						},
						"deleted_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "删除时间。",
						},
						"instance_grp_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例组ID。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地位。",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例组类型。哈哈组； ro-只读组。",
						},
						"updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "更新时间。",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "内网IP。",
						},
						"vport": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "内网端口。",
						},
						"wan_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "公共域名。",
						},
						"wan_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "公共IP。",
						},
						"wan_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "公共端口。",
						},
						"wan_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "公共地位。",
						},
						"instance_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "实例组包含实例信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "用户 Uin。",
									},
									"app_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "用户应用程序 ID。",
									},
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "集群的id。",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "集群的名称。",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例的id。",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例的名称。",
									},
									"project_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "项目的id。",
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "地区。",
									},
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "可用区。",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例的状态。",
									},
									"status_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例状态中文描述。",
									},
									"db_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "数据库类型。",
									},
									"db_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "数据库版本。",
									},
									"cpu": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "CPU，单位：CORE。",
									},
									"memory": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "内存，单位：GB。",
									},
									"storage": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "存储，单位：GB。",
									},
									"instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例类型。",
									},
									"instance_role": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例角色。",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "更新时间。",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "创造时间。",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC网络ID。",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "子网 ID。",
									},
									"vip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例内网IP。",
									},
									"vport": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "实例内网VPort。",
									},
									"pay_mode": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "付费模式。",
									},
									"period_end_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例过期时间。",
									},
									"destroy_deadline_text": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "毁掉最后期限。",
									},
									"isolate_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "隔离时间。",
									},
									"net_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "网型。",
									},
									"wan_domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "公共领域。",
									},
									"wan_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "公共IP。",
									},
									"wan_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "公共端口。",
									},
									"wan_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "公共地位。",
									},
									"destroy_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例销毁时间。",
									},
									"cynos_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cynos 内核版本。",
									},
									"processing_task": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "任务正在处理中。",
									},
									"renew_flag": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "更新旗帜。",
									},
									"min_cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "无服务器实例最低 CPU。",
									},
									"max_cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "无服务器实例最大 CPU。",
									},
									"serverless_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Serverless实例状态，可选值：resumepause。",
									},
									"storage_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Prepaid Storage Id。注意：该字段可能返回null，表示取不到有效值。",
									},
									"storage_pay_mode": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "存储付费类型。",
									},
									"physical_zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "物理区域。",
									},
									"business_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "业务类型。注意：该字段可能返回null，表示取不到有效值。",
									},
									"tasks": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "任务列表。注意：该字段可能返回null，表示取不到有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"task_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "任务自增ID。注意：该字段可能返回null，表示取不到有效值。",
												},
												"task_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "任务类型。注意：该字段可能返回null，表示取不到有效值。",
												},
												"task_status": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "任务状态。注意：该字段可能返回null，表示取不到有效值。",
												},
												"object_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "任务ID（集群ID|实例组ID|实例ID）。注意：该字段可能返回null，表示取不到有效值。",
												},
												"object_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "对象类型。注意：该字段可能返回null，表示取不到有效值。",
												},
											},
										},
									},
									"is_freeze": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "是否冻结。注意：该字段可能返回null，表示取不到有效值。",
									},
									"resource_tags": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "资源标签。注意：该字段可能返回null，表示取不到有效值。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"tag_key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "标签的关键。",
												},
												"tag_value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "标签的值。",
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

func dataSourceTencentCloudCynosdbClusterInstanceGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cynosdb_cluster_instance_groups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var clusterId string
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var instanceGrpInfoList []*cynosdb.CynosdbInstanceGrp

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClusterInstanceGrps(ctx, clusterId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instanceGrpInfoList = result.Response.InstanceGrpInfoList
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceGrpInfoList))
	tmpList := make([]map[string]interface{}, 0, len(instanceGrpInfoList))
	for _, instanceGrpInfo := range instanceGrpInfoList {
		ids = append(ids, *instanceGrpInfo.InstanceGrpId)
		instanceGrpInfoMap := make(map[string]interface{})
		instanceGrpInfoMap["app_id"] = instanceGrpInfo.AppId
		instanceGrpInfoMap["cluster_id"] = instanceGrpInfo.ClusterId
		instanceGrpInfoMap["created_time"] = instanceGrpInfo.CreatedTime
		instanceGrpInfoMap["deleted_time"] = instanceGrpInfo.DeletedTime
		instanceGrpInfoMap["instance_grp_id"] = instanceGrpInfo.InstanceGrpId
		instanceGrpInfoMap["status"] = instanceGrpInfo.Status
		instanceGrpInfoMap["type"] = instanceGrpInfo.Type
		instanceGrpInfoMap["updated_time"] = instanceGrpInfo.UpdatedTime
		instanceGrpInfoMap["vip"] = instanceGrpInfo.Vip
		instanceGrpInfoMap["vport"] = instanceGrpInfo.Vport
		instanceGrpInfoMap["wan_domain"] = instanceGrpInfo.WanDomain
		instanceGrpInfoMap["wan_ip"] = instanceGrpInfo.WanIP
		instanceGrpInfoMap["wan_port"] = instanceGrpInfo.WanPort
		instanceGrpInfoMap["wan_status"] = instanceGrpInfo.WanStatus
		if instanceGrpInfo.InstanceSet != nil {
			instances := make([]map[string]interface{}, 0)
			for _, instance := range instanceGrpInfo.InstanceSet {
				instanceMap := make(map[string]interface{})
				instanceMap["uin"] = instance.Uin
				instanceMap["app_id"] = instance.AppId
				instanceMap["cluster_id"] = instance.ClusterId
				instanceMap["cluster_name"] = instance.ClusterName
				instanceMap["instance_id"] = instance.InstanceId
				instanceMap["instance_name"] = instance.InstanceName
				instanceMap["project_id"] = instance.ProjectId
				instanceMap["region"] = instance.Region
				instanceMap["zone"] = instance.Zone
				instanceMap["status"] = instance.Status
				instanceMap["status_desc"] = instance.StatusDesc
				instanceMap["db_type"] = instance.DbType
				instanceMap["db_version"] = instance.DbVersion
				instanceMap["cpu"] = instance.Cpu
				instanceMap["memory"] = instance.Memory
				instanceMap["storage"] = instance.Storage
				instanceMap["instance_type"] = instance.InstanceType
				instanceMap["instance_role"] = instance.InstanceRole
				instanceMap["update_time"] = instance.UpdateTime
				instanceMap["create_time"] = instance.CreateTime
				instanceMap["vpc_id"] = instance.VpcId
				instanceMap["subnet_id"] = instance.SubnetId
				instanceMap["vip"] = instance.Vip
				instanceMap["vport"] = instance.Vport
				instanceMap["pay_mode"] = instance.PayMode
				instanceMap["period_end_time"] = instance.PeriodEndTime
				instanceMap["destroy_deadline_text"] = instance.DestroyDeadlineText
				instanceMap["isolate_time"] = instance.IsolateTime
				instanceMap["net_type"] = instance.NetType
				instanceMap["wan_domain"] = instance.WanDomain
				instanceMap["wan_ip"] = instance.WanIP
				instanceMap["wan_port"] = instance.WanPort
				instanceMap["wan_status"] = instance.WanStatus
				instanceMap["destroy_time"] = instance.DestroyTime
				instanceMap["cynos_version"] = instance.CynosVersion
				instanceMap["processing_task"] = instance.ProcessingTask
				instanceMap["renew_flag"] = instance.RenewFlag
				instanceMap["min_cpu"] = instance.MinCpu
				instanceMap["max_cpu"] = instance.MaxCpu
				instanceMap["serverless_status"] = instance.ServerlessStatus
				instanceMap["storage_id"] = instance.StorageId
				instanceMap["storage_pay_mode"] = instance.StoragePayMode
				instanceMap["physical_zone"] = instance.PhysicalZone
				instanceMap["business_type"] = instance.BusinessType
				instanceMap["is_freeze"] = instance.IsFreeze
				tasks := make([]map[string]interface{}, 0)
				if instance.Tasks != nil {
					for _, task := range instance.Tasks {
						taskMap := make(map[string]interface{})
						taskMap["task_id"] = task.TaskId
						taskMap["task_type"] = task.TaskType
						taskMap["task_status"] = task.TaskStatus
						taskMap["object_id"] = task.ObjectId
						taskMap["object_type"] = task.ObjectType

						tasks = append(tasks, taskMap)
					}
					instanceMap["tasks"] = tasks
				}
				tags := make([]map[string]interface{}, 0)
				if instance.ResourceTags != nil {
					for _, tag := range instance.ResourceTags {
						tagMap := make(map[string]interface{})
						tagMap["tag_key"] = tag.TagKey
						tagMap["tag_value"] = tag.TagValue

						tags = append(tags, tagMap)
					}
					instanceMap["resource_tags"] = tags
				}
				instances = append(instances, instanceMap)
				instanceGrpInfoMap["instance_set"] = instances
			}
		}
		tmpList = append(tmpList, instanceGrpInfoMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("instance_grp_info_list", tmpList)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
