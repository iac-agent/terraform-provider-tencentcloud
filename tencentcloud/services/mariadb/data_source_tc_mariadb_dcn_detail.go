package mariadb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMariadbDcnDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbDcnDetailRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "实例 ID",
			},
			"dcn_details": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "DCN synchronization details。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例名称",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 其中 实例 resides。",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability 可用区 其中 实例 resides。",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 IP 地址",
						},
						"vipv6": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 IPv6 地址",
						},
						"vport": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例端口",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例状态",
						},
						"status_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例状态 描述",
						},
						"dcn_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "DCN flag. 有效值：`1` (primary)，`2` (disaster recovery)。",
						},
						"dcn_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "DCN 状态 有效值：`0` (none)，`1` (creating)，`2` (syncing)，`3` (disconnected)。",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CPU 核数 的 实例。",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例 内存 容量 （GB）。",
						},
						"storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例 存储 容量 （GB）。",
						},
						"pay_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Billing 模式",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 实例 在 格式 的 2006-01-02 15:04:05。",
						},
						"period_end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "过期时间 的 实例 在 格式 的 2006-01-02 15:04:05。",
						},
						"instance_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例类型 有效值：`1` (dedicated primary 实例)，`2` (non-dedicated primary 实例)，`3` (non-dedicated disaster recovery 实例)，`4` (dedicated disaster recovery 实例)。",
						},
						"replica_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Configuration 信息 的 DCN 复制. 此 字段 是 null 对于 primary 实例.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ro_replication_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DCN running 状态 有效值：`START` (running)，`STOP` (pause)注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"delay_replication_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Delayed 复制 类型 有效值：`DEFAULT` (无 延迟)，`DUE_TIME` (指定 复制 时间)注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"due_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specified 时间 对于 delayed replication注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"replication_delay": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 秒 到 延迟 replication注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"replica_status": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "DCN 复制 状态 此 字段 是 null 对于 primary 实例.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DCN running 状态 有效值：`START` (running)，`STOP` (pause).注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"delay": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "当前 延迟，其中 takes 延迟 值 的 副本 实例。",
									},
								},
							},
						},
						"encrypt_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether KMS 是 已启用",
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

func dataSourceTencentCloudMariadbDcnDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mariadb_dcn_detail.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		dcnDetails []*mariadb.DcnDetailItem
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbDcnDetailByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		dcnDetails = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(dcnDetails))
	tmpList := make([]map[string]interface{}, 0, len(dcnDetails))

	if dcnDetails != nil {
		for _, dcnDetailItem := range dcnDetails {
			dcnDetailItemMap := map[string]interface{}{}

			if dcnDetailItem.InstanceId != nil {
				dcnDetailItemMap["instance_id"] = dcnDetailItem.InstanceId
			}

			if dcnDetailItem.InstanceName != nil {
				dcnDetailItemMap["instance_name"] = dcnDetailItem.InstanceName
			}

			if dcnDetailItem.Region != nil {
				dcnDetailItemMap["region"] = dcnDetailItem.Region
			}

			if dcnDetailItem.Zone != nil {
				dcnDetailItemMap["zone"] = dcnDetailItem.Zone
			}

			if dcnDetailItem.Vip != nil {
				dcnDetailItemMap["vip"] = dcnDetailItem.Vip
			}

			if dcnDetailItem.Vipv6 != nil {
				dcnDetailItemMap["vipv6"] = dcnDetailItem.Vipv6
			}

			if dcnDetailItem.Vport != nil {
				dcnDetailItemMap["vport"] = dcnDetailItem.Vport
			}

			if dcnDetailItem.Status != nil {
				dcnDetailItemMap["status"] = dcnDetailItem.Status
			}

			if dcnDetailItem.StatusDesc != nil {
				dcnDetailItemMap["status_desc"] = dcnDetailItem.StatusDesc
			}

			if dcnDetailItem.DcnFlag != nil {
				dcnDetailItemMap["dcn_flag"] = dcnDetailItem.DcnFlag
			}

			if dcnDetailItem.DcnStatus != nil {
				dcnDetailItemMap["dcn_status"] = dcnDetailItem.DcnStatus
			}

			if dcnDetailItem.Cpu != nil {
				dcnDetailItemMap["cpu"] = dcnDetailItem.Cpu
			}

			if dcnDetailItem.Memory != nil {
				dcnDetailItemMap["memory"] = dcnDetailItem.Memory
			}

			if dcnDetailItem.Storage != nil {
				dcnDetailItemMap["storage"] = dcnDetailItem.Storage
			}

			if dcnDetailItem.PayMode != nil {
				dcnDetailItemMap["pay_mode"] = dcnDetailItem.PayMode
			}

			if dcnDetailItem.CreateTime != nil {
				dcnDetailItemMap["create_time"] = dcnDetailItem.CreateTime
			}

			if dcnDetailItem.PeriodEndTime != nil {
				dcnDetailItemMap["period_end_time"] = dcnDetailItem.PeriodEndTime
			}

			if dcnDetailItem.InstanceType != nil {
				dcnDetailItemMap["instance_type"] = dcnDetailItem.InstanceType
			}

			if dcnDetailItem.ReplicaConfig != nil {
				replicaConfigMap := map[string]interface{}{}

				if dcnDetailItem.ReplicaConfig.RoReplicationMode != nil {
					replicaConfigMap["ro_replication_mode"] = dcnDetailItem.ReplicaConfig.RoReplicationMode
				}

				if dcnDetailItem.ReplicaConfig.DelayReplicationType != nil {
					replicaConfigMap["delay_replication_type"] = dcnDetailItem.ReplicaConfig.DelayReplicationType
				}

				if dcnDetailItem.ReplicaConfig.DueTime != nil {
					replicaConfigMap["due_time"] = dcnDetailItem.ReplicaConfig.DueTime
				}

				if dcnDetailItem.ReplicaConfig.ReplicationDelay != nil {
					replicaConfigMap["replication_delay"] = dcnDetailItem.ReplicaConfig.ReplicationDelay
				}

				dcnDetailItemMap["replica_config"] = []interface{}{replicaConfigMap}
			}

			if dcnDetailItem.ReplicaStatus != nil {
				replicaStatusMap := map[string]interface{}{}

				if dcnDetailItem.ReplicaStatus.Status != nil {
					replicaStatusMap["status"] = dcnDetailItem.ReplicaStatus.Status
				}

				if dcnDetailItem.ReplicaStatus.Delay != nil {
					replicaStatusMap["delay"] = dcnDetailItem.ReplicaStatus.Delay
				}

				dcnDetailItemMap["replica_status"] = []interface{}{replicaStatusMap}
			}

			if dcnDetailItem.EncryptStatus != nil {
				dcnDetailItemMap["encrypt_status"] = dcnDetailItem.EncryptStatus
			}

			ids = append(ids, *dcnDetailItem.InstanceId)
			tmpList = append(tmpList, dcnDetailItemMap)
		}

		_ = d.Set("dcn_details", tmpList)
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
