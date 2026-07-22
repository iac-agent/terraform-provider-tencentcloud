package cbs

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCbsSnapshots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCbsSnapshotsRead,

		Schema: map[string]*schema.Schema{
			"snapshot_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID 快照 到 是 queried。",
			},
			"snapshot_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 快照 到 是 queried。",
			},
			"storage_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID CBS 其中 此 快照 创建 从。",
			},
			"storage_usage": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CBS_STORAGE_USAGE),
				Description:  "Types 的 CBS 其中 此 快照 创建 从，和 可用 值 include `SYSTEM_DISK` 和 `DATA_DISK`。",
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID 项目 within 快照。",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "可用 可用区 该 CBS 实例 locates 在。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"snapshot_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 快照. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snapshot_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 快照。",
						},
						"snapshot_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 快照。",
						},
						"storage_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID CBS 其中 此 快照 创建 从。",
						},
						"storage_usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Types 的 CBS 其中 此 快照 创建 从。",
						},
						"storage_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Volume 的 存储 其中 此 快照 创建 从。",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用 可用区 该 CBS 实例 locates 在。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID 项目 within 快照。",
						},
						"percent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Snapshot creation progress percentage。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 快照。",
						},
						"encrypt": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "表示是否snapshot 是 encrypted。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCbsSnapshotsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cbs_snapshots.read")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		cbsService = CbsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	params := make(map[string]string)
	if v, ok := d.GetOk("snapshot_id"); ok {
		params["snapshot-id"] = v.(string)
	}

	if v, ok := d.GetOk("snapshot_name"); ok {
		params["snapshot-name"] = v.(string)
	}

	if v, ok := d.GetOk("storage_id"); ok {
		params["disk-id"] = v.(string)
	}

	if v, ok := d.GetOk("storage_usage"); ok {
		params["disk-usage"] = v.(string)
	}

	if v, ok := d.GetOk("project_id"); ok {
		params["project-id"] = v.(string)
	}

	if v, ok := d.GetOk("availability_zone"); ok {
		params["zone"] = v.(string)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		snapshots, e := cbsService.DescribeSnapshotsByFilter(ctx, params)
		if e != nil {
			return tccommon.RetryError(e)
		}

		ids := make([]string, 0, len(snapshots))
		snapshotList := make([]map[string]interface{}, 0, len(snapshots))
		for _, snapshot := range snapshots {
			mapping := map[string]interface{}{
				"snapshot_id":       *snapshot.SnapshotId,
				"snapshot_name":     *snapshot.SnapshotName,
				"storage_id":        *snapshot.DiskId,
				"storage_usage":     *snapshot.DiskUsage,
				"storage_size":      *snapshot.DiskSize,
				"availability_zone": *snapshot.Placement.Zone,
				"project_id":        *snapshot.Placement.ProjectId,
				"percent":           *snapshot.Percent,
				"create_time":       *snapshot.CreateTime,
				"encrypt":           *snapshot.Encrypt,
			}

			snapshotList = append(snapshotList, mapping)
			ids = append(ids, *snapshot.SnapshotId)
		}

		d.SetId(helper.DataResourceIdsHash(ids))
		if e = d.Set("snapshot_list", snapshotList); e != nil {
			log.Printf("[CRITAL]%s provider set snapshot list fail, reason:%s\n ", logId, e.Error())
			return resource.NonRetryableError(e)
		}

		output, ok := d.GetOk("result_output_file")
		if ok && output.(string) != "" {
			if e := tccommon.WriteToFile(output.(string), snapshotList); e != nil {
				return resource.NonRetryableError(e)
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read cbs snapshots failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}
