package lighthouse

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudLighthouseInstanceSnapshot1() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseInstanceSnapshot1Create,
		Read:   resourceTencentCloudLighthouseInstanceSnapshot1Read,
		Update: resourceTencentCloudLighthouseInstanceSnapshot1Update,
		Delete: resourceTencentCloudLighthouseInstanceSnapshot1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of the instance for which to create a snapshot.",
			},
			"snapshot_name": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Snapshot name, which can contain up to 60 characters.",
			},
			"tags": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Tags to associate with the snapshot.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Required:    true,
							Type:        schema.TypeString,
							Description: "Tag key.",
						},
						"value": {
							Required:    true,
							Type:        schema.TypeString,
							Description: "Tag value.",
						},
					},
				},
			},
			"snapshot_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Snapshot ID.",
			},
			"disk_usage": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Disk type for which the snapshot was created. Values: `SYSTEM_DISK` (system disk).",
			},
			"disk_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "ID of the disk for which the snapshot was created.",
			},
			"disk_size": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Size of the disk for which the snapshot was created, in GB.",
			},
			"snapshot_state": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Snapshot state. Values: `NORMAL`, `CREATING`, `ROLLBACKING`.",
			},
			"percent": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Snapshot creation progress percentage. After the snapshot is created successfully, the value of this field will be 100.",
			},
			"latest_operation": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Latest operation of the snapshot. Only recorded when creating or rolling back the snapshot. Values: `CreateInstanceSnapshot`, `RollbackInstanceSnapshot`.",
			},
			"latest_operation_state": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Latest operation state of the snapshot. Only recorded when creating or rolling back the snapshot. Values: `SUCCESS`, `OPERATING`, `FAILED`.",
			},
			"latest_operation_request_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The unique request ID for the latest operation of the snapshot. Only recorded when creating or rolling back the snapshot.",
			},
			"created_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Snapshot creation time.",
			},
		},
	}
}

func resourceTencentCloudLighthouseInstanceSnapshot1Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_instance_snapshot_1.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = lighthouse.NewCreateInstanceSnapshotRequest()
		response   = lighthouse.NewCreateInstanceSnapshotResponse()
		snapshotId string
	)
	request.InstanceId = helper.String(d.Get("instance_id").(string))

	if v, ok := d.GetOk("snapshot_name"); ok {
		request.SnapshotName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := v.([]interface{})
		request.Tags = make([]*lighthouse.Tag, 0, len(tags))
		for _, tag := range tags {
			tagMap := tag.(map[string]interface{})
			request.Tags = append(request.Tags, &lighthouse.Tag{
				Key:   helper.String(tagMap["key"].(string)),
				Value: helper.String(tagMap["value"].(string)),
			})
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().CreateInstanceSnapshot(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create lighthouse instance_snapshot_1 failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.SnapshotId == nil || *response.Response.SnapshotId == "" {
		log.Printf("[CRITAL]%s create lighthouse instance_snapshot_1 returned empty SnapshotId, logId=%s", logId, logId)
		return fmt.Errorf("create lighthouse instance_snapshot_1 returned empty SnapshotId")
	}

	snapshotId = *response.Response.SnapshotId
	d.SetId(snapshotId)

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"NORMAL"}, 20*tccommon.ReadRetryTimeout, time.Second, service.LighthouseSnapshotStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseInstanceSnapshot1Read(d, meta)
}

func resourceTencentCloudLighthouseInstanceSnapshot1Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_instance_snapshot_1.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	snapshotId := d.Id()

	snapshot, err := service.DescribeLighthouseSnapshotById(ctx, snapshotId)
	if err != nil {
		return err
	}

	if snapshot == nil {
		log.Printf("[WARN]%s resource `lighthouse_instance_snapshot_1` [%s] not found, please check if it has been deleted.\n", logId, snapshotId)
		d.SetId("")
		return nil
	}

	if snapshot.SnapshotId != nil {
		_ = d.Set("snapshot_id", snapshot.SnapshotId)
	}

	if snapshot.SnapshotName != nil {
		_ = d.Set("snapshot_name", snapshot.SnapshotName)
	}

	if snapshot.DiskUsage != nil {
		_ = d.Set("disk_usage", snapshot.DiskUsage)
	}

	if snapshot.DiskId != nil {
		_ = d.Set("disk_id", snapshot.DiskId)
	}

	if snapshot.DiskSize != nil {
		_ = d.Set("disk_size", snapshot.DiskSize)
	}

	if snapshot.SnapshotState != nil {
		_ = d.Set("snapshot_state", snapshot.SnapshotState)
	}

	if snapshot.Percent != nil {
		_ = d.Set("percent", snapshot.Percent)
	}

	if snapshot.LatestOperation != nil {
		_ = d.Set("latest_operation", snapshot.LatestOperation)
	}

	if snapshot.LatestOperationState != nil {
		_ = d.Set("latest_operation_state", snapshot.LatestOperationState)
	}

	if snapshot.LatestOperationRequestId != nil {
		_ = d.Set("latest_operation_request_id", snapshot.LatestOperationRequestId)
	}

	if snapshot.CreatedTime != nil {
		_ = d.Set("created_time", snapshot.CreatedTime)
	}

	if snapshot.Tags != nil {
		tags := make([]map[string]interface{}, 0, len(snapshot.Tags))
		for _, tag := range snapshot.Tags {
			tagMap := make(map[string]interface{})
			if tag.Key != nil {
				tagMap["key"] = *tag.Key
			}
			if tag.Value != nil {
				tagMap["value"] = *tag.Value
			}
			tags = append(tags, tagMap)
		}
		_ = d.Set("tags", tags)
	}

	return nil
}

func resourceTencentCloudLighthouseInstanceSnapshot1Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_instance_snapshot_1.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := lighthouse.NewModifySnapshotAttributeRequest()

	snapshotId := d.Id()

	request.SnapshotId = &snapshotId

	if d.HasChange("snapshot_name") {
		if v, ok := d.GetOk("snapshot_name"); ok {
			request.SnapshotName = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseLighthouseClient().ModifySnapshotAttribute(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update lighthouse instance_snapshot_1 failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudLighthouseInstanceSnapshot1Read(d, meta)
}

func resourceTencentCloudLighthouseInstanceSnapshot1Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_lighthouse_instance_snapshot_1.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := LightHouseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	snapshotId := d.Id()

	if err := service.DeleteLighthouseSnapshotById(ctx, snapshotId); err != nil {
		return err
	}

	return nil
}
