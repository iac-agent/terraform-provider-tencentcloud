package cbs

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCbsStoragesSet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCbsStoragesSetRead,

		Schema: map[string]*schema.Schema{
			"storage_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID CBS 到 是 queried。",
			},
			"storage_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 CBS 到 是 queried。",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "可用 可用区 该 CBS 实例 locates 在。",
			},
			"dedicated_cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Exclusive 集群 ID",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID 项目 使用 其中 CBS 是 associated。",
			},
			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CBS_STORAGE_TYPE),
				Description:  "过滤器 通过 云 磁盘 media 类型 (`CLOUD_BASIC`: HDD 云 磁盘 | `CLOUD_PREMIUM`: Premium Cloud Storage | `CLOUD_SSD`: SSD 云 磁盘)。",
			},
			"storage_usage": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "过滤器 通过 云 磁盘 类型 (`SYSTEM_DISK`: 系统 磁盘 | `DATA_DISK`: 数据 磁盘)。",
			},
			"charge_type": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List 过滤器 通过 磁盘 计费类型 (`POSTPAID_BY_HOUR` | `PREPAID` | `CDCPAID` | `DEDICATED_CLUSTER_PAID`)。",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"portable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "过滤器 通过 是否disk 是 portable (Boolean `true` 或 `false`)。",
			},
			"storage_state": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List 过滤器 通过 磁盘 state (`UNATTACHED` | `ATTACHING` | `ATTACHED` | `DETACHING` | `EXPANDING` | `ROLLBACKING` | `TORECYCLE`)。",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"instance_ips": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List 过滤器 通过 attached 实例 公有 或 私有 IPs。",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"instance_name": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List 过滤器 通过 attached 实例名称",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"tag_keys": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List 过滤器 通过 标签 keys。",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"tag_values": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List 过滤器 通过 标签 值。",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"storage_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 存储. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID CBS。",
						},
						"storage_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 CBS。",
						},
						"storage_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Types 的 存储 medium。",
						},
						"storage_usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Types 的 CBS。",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用区 的 CBS。",
						},
						"dedicated_cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Exclusive 集群 ID",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID 项目。",
						},
						"storage_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Volume 的 CBS。",
						},
						"attached": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "表示是否CBS 是 mounted CVM。",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID CVM 实例 该 是 mounted 通过 此 CBS。",
						},
						"kms_key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Kms 键 ID。",
						},
						"encrypt": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "表示是否CBS 是 encrypted。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 CBS。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态 CBS。",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "可用 标签 within 此 CBS。",
						},
						"prepaid_renew_flag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "way 该 CBS 实例 将 是 renew automatically 或 不 当 它 reach end 的 prepaid tenancy。",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pay 类型 CBS 实例。",
						},
						"throughput_performance": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Add extra performance 到 数据 磁盘. Only works 当 磁盘 类型 是 `CLOUD_TSSD` 或 `CLOUD_HSSD`。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCbsStoragesSetRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cbs_storages_set.read")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		cbsService = CbsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	params := make(map[string]interface{})
	if v, ok := d.GetOk("storage_id"); ok {
		params["disk-id"] = v.(string)
	}

	if v, ok := d.GetOk("storage_name"); ok {
		params["disk-name"] = v.(string)
	}

	if v, ok := d.GetOk("availability_zone"); ok {
		params["zone"] = v.(string)
	}

	if v, ok := d.GetOk("dedicated_cluster_id"); ok {
		params["dedicated-cluster-id"] = v.(string)
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		params["project-id"] = fmt.Sprintf("%d", v.(int))
	}

	if v, ok := d.GetOk("storage_type"); ok {
		params["disk-type"] = v.(string)
	}

	if v, ok := d.GetOk("storage_usage"); ok {
		params["disk-usage"] = v.(string)
	}

	if v, ok := d.GetOk("charge_type"); ok {
		params["disk-charge-type"] = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOk("portable"); ok {
		if v.(bool) {
			params["portable"] = "TRUE"
		} else {
			params["portable"] = "FALSE"
		}
	}

	if v, ok := d.GetOk("storage_state"); ok {
		params["disk-state"] = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOk("instance_ips"); ok {
		params["instance-ip-address"] = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		params["instance-name"] = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOk("tag_keys"); ok {
		params["tag-key"] = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOk("tag_values"); ok {
		params["tag-value"] = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	storages, e := cbsService.DescribeDisksInParallelByFilter(ctx, params)
	if e != nil {
		return e
	}

	ids := make([]string, 0, len(storages))
	storageList := make([]map[string]interface{}, 0, len(storages))
	for _, storage := range storages {
		mapping := map[string]interface{}{
			"storage_id":             storage.DiskId,
			"storage_name":           storage.DiskName,
			"storage_usage":          storage.DiskUsage,
			"storage_type":           storage.DiskType,
			"availability_zone":      storage.Placement.Zone,
			"dedicated_cluster_id":   storage.Placement.DedicatedClusterId,
			"project_id":             storage.Placement.ProjectId,
			"storage_size":           storage.DiskSize,
			"attached":               storage.Attached,
			"instance_id":            storage.InstanceId,
			"kms_key_id":             storage.KmsKeyId,
			"encrypt":                storage.Encrypt,
			"create_time":            storage.CreateTime,
			"status":                 storage.DiskState,
			"prepaid_renew_flag":     storage.RenewFlag,
			"charge_type":            storage.DiskChargeType,
			"throughput_performance": storage.ThroughputPerformance,
		}

		if storage.Tags != nil {
			tags := make(map[string]interface{}, len(storage.Tags))
			for _, t := range storage.Tags {
				tags[*t.Key] = *t.Value
			}

			mapping["tags"] = tags
		}

		storageList = append(storageList, mapping)
		ids = append(ids, *storage.DiskId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e = d.Set("storage_list", storageList); e != nil {
		log.Printf("[CRITAL]%s provider set storage list fail, reason:%s\n ", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), storageList); e != nil {
			return e
		}
	}

	return nil

}
