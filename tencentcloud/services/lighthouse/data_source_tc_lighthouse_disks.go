package lighthouse

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudLighthouseInstanceDisks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseInstanceDisksRead,
		Schema: map[string]*schema.Schema{
			"disk_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "列表 磁盘 ids。",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "过滤器 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Fields 到 是 filtered. 有效 names: `磁盘-ID`: Filters 通过 磁盘 ID; `实例-ID`: 过滤器 通过 实例 ID; `磁盘-名称`: 过滤器 通过 磁盘 名称; `可用区`: 过滤器 通过 可用区; `磁盘-usage`: 过滤器 通过 磁盘 usage(Values: `SYSTEM_DISK` 或 `DATA_DISK`); `磁盘-state`: 过滤器 通过 磁盘 state。",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "值 的 字段。",
						},
					},
				},
			},

			"disk_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Cloud 磁盘 信息 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk ID。",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability 可用区",
						},
						"disk_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk 名称",
						},
						"disk_usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk usage。",
						},
						"disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk 类型",
						},
						"disk_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk 计费类型",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk 大小。",
						},
						"renew_flag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "续费标识",
						},
						"disk_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk state. 有效 值:`PENDING`，`UNATTACHED`，`ATTACHING`，`ATTACHED`，`DETACHING`，`SHUTDOWN`，`CREATED_FAILED`，`TERMINATING`，`DELETING`，`FREEZING`。",
						},
						"attached": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Disk attach state。",
						},
						"delete_with_instance": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否release 使用 实例。",
						},
						"latest_operation": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Latest operation。",
						},
						"latest_operation_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Latest operation state。",
						},
						"latest_operation_request_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Latest operation 请求 ID。",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Created 时间. Expressed according 到 ISO8601 standard，和 使用 UTC 时间. 格式 是 `YYYY-MM-DDThh:mm:ssZ`。",
						},
						"expired_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "过期时间. Expressed according 到 ISO8601 standard，和 使用 UTC 时间. 格式 是 `YYYY-MM-DDThh:mm:ssZ`。",
						},
						"isolated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Isolated 时间. Expressed according 到 ISO8601 standard，和 使用 UTC 时间. 格式 是 `YYYY-MM-DDThh:mm:ssZ`。",
						},
						"disk_backup_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 existing 备份 points 的 云 磁盘。",
						},
						"disk_backup_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 备份 points 配额 对于 云 磁盘。",
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

func dataSourceTencentCloudLighthouseInstanceDisksRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_lighthouse_instance_disks.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	diskIds := make([]string, 0)
	for _, diskId := range d.Get("disk_ids").(*schema.Set).List() {
		diskIds = append(diskIds, diskId.(string))
	}
	filters := make([]*lighthouse.Filter, 0)
	if v, ok := d.GetOk("filters"); ok {
		filterSet := v.([]interface{})

		for _, item := range filterSet {
			filter := lighthouse.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			filters = append(filters, &filter)
		}
	}
	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	disks, err := service.DescribeLighthouseDisk(ctx, diskIds, filters)
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	diskList := make([]map[string]interface{}, 0)
	for _, disk := range disks {
		diskMap := make(map[string]interface{})
		if disk.DiskId != nil {
			diskMap["disk_id"] = disk.DiskId
			ids = append(ids, *disk.DiskId)
		}

		if disk.InstanceId != nil {
			diskMap["instance_id"] = disk.InstanceId
		}

		if disk.Zone != nil {
			diskMap["zone"] = disk.Zone
		}

		if disk.DiskName != nil {
			diskMap["disk_name"] = disk.DiskName
		}

		if disk.DiskUsage != nil {
			diskMap["disk_usage"] = disk.DiskUsage
		}

		if disk.DiskType != nil {
			diskMap["disk_type"] = disk.DiskType
		}

		if disk.DiskChargeType != nil {
			diskMap["disk_charge_type"] = disk.DiskChargeType
		}

		if disk.DiskSize != nil {
			diskMap["disk_size"] = disk.DiskSize
		}

		if disk.RenewFlag != nil {
			diskMap["renew_flag"] = disk.RenewFlag
		}

		if disk.DiskState != nil {
			diskMap["disk_state"] = disk.DiskState
		}

		if disk.Attached != nil {
			diskMap["attached"] = disk.Attached
		}

		if disk.DeleteWithInstance != nil {
			diskMap["delete_with_instance"] = disk.DeleteWithInstance
		}

		if disk.LatestOperation != nil {
			diskMap["latest_operation"] = disk.LatestOperation
		}

		if disk.LatestOperationState != nil {
			diskMap["latest_operation_state"] = disk.LatestOperationState
		}

		if disk.LatestOperationRequestId != nil {
			diskMap["latest_operation_request_id"] = disk.LatestOperationRequestId
		}

		if disk.CreatedTime != nil {
			diskMap["created_time"] = disk.CreatedTime
		}

		if disk.ExpiredTime != nil {
			diskMap["expired_time"] = disk.ExpiredTime
		}

		if disk.IsolatedTime != nil {
			diskMap["isolated_time"] = disk.IsolatedTime
		}

		if disk.DiskBackupCount != nil {
			diskMap["disk_backup_count"] = disk.DiskBackupCount
		}

		if disk.DiskBackupQuota != nil {
			diskMap["disk_backup_quota"] = disk.DiskBackupQuota
		}

		diskList = append(diskList, diskMap)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("disk_list", diskList)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), diskList); e != nil {
			return e
		}
	}
	return nil
}
