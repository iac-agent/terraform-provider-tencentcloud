package tcr

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTcrReplicationInstanceSyncStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcrReplicationInstanceSyncStatusRead,
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "master registry ID。",
			},

			"replication_registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "synchronization 实例 ID",
			},

			"replication_region_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "synchronization 实例 地域 ID",
			},

			"show_replication_log": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "是否display synchronization 日志。",
			},

			"replication_status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "sync 状态",
			},

			"replication_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "sync 完整 时间。",
			},

			"replication_log": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "sync 日志. 注意：此字段可能返回 null，表示无法获取有效值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "资源类型 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "来源 镜像. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"destination": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "destination 资源. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "sync 状态 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "开始时间. 注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "结束时间. 注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudTcrReplicationInstanceSyncStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tcr_replication_instance_sync_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var replicationRegistryId string

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("registry_id"); ok {
		paramMap["RegistryId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("replication_registry_id"); ok {
		paramMap["ReplicationRegistryId"] = helper.String(v.(string))
		replicationRegistryId = v.(string)
	}

	if v, _ := d.GetOk("replication_region_id"); v != nil {
		paramMap["ReplicationRegionId"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("show_replication_log"); v != nil {
		paramMap["ShowReplicationLog"] = helper.Bool(v.(bool))
	}

	service := TCRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var (
		result *tcr.DescribeReplicationInstanceSyncStatusResponseParams
		e      error
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e = service.DescribeTcrReplicationInstanceSyncStatusByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if result.ReplicationStatus != nil {
		_ = d.Set("replication_status", result.ReplicationStatus)
	}

	if result.ReplicationTime != nil {
		_ = d.Set("replication_time", result.ReplicationTime)
	}

	if result.ReplicationLog != nil {
		replicationLogMap := map[string]interface{}{}
		replicationLog := result.ReplicationLog
		if replicationLog.ResourceType != nil {
			replicationLogMap["resource_type"] = replicationLog.ResourceType
		}

		if replicationLog.Source != nil {
			replicationLogMap["source"] = replicationLog.Source
		}

		if replicationLog.Destination != nil {
			replicationLogMap["destination"] = replicationLog.Destination
		}

		if replicationLog.Status != nil {
			replicationLogMap["status"] = replicationLog.Status
		}

		if replicationLog.StartTime != nil {
			replicationLogMap["start_time"] = replicationLog.StartTime
		}

		if replicationLog.EndTime != nil {
			replicationLogMap["end_time"] = replicationLog.EndTime
		}
		_ = d.Set("replication_log", []map[string]interface{}{replicationLogMap})
	}

	d.SetId(helper.DataResourceIdHash(replicationRegistryId))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
