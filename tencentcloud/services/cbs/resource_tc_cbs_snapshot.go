package cbs

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
)

func ResourceTencentCloudCbsSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCbsSnapshotCreate,
		Read:   resourceTencentCloudCbsSnapshotRead,
		Update: resourceTencentCloudCbsSnapshotUpdate,
		Delete: resourceTencentCloudCbsSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"storage_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID CBS 其中 此 快照 创建 从。",
			},
			"snapshot_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(2, 60),
				Description:  "名称 快照。",
			},
			"disk_usage": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"SYSTEM_DISK", "DATA_DISK"}),
				Description:  "类型 云 磁盘 associated 使用 快照: SYSTEM_DISK: 系统 磁盘; DATA_DISK: 数据 磁盘. 如果未填写 在， 快照 类型 将 是 consistent 使用 云 磁盘 类型 此 参数 是 使用 在 some scenarios 其中 users need 到 create 数据 磁盘 快照 从 系统 磁盘 对于 shared 使用。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "可用 标签 within 此 CBS Snapshot。",
			},
			// computed
			"storage_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Volume 的 存储 其中 此 快照 创建 从。",
			},
			"snapshot_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "状态 快照。",
			},
			"disk_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Deprecated:  "It has been deprecated from version 1.82.14. Please use `disk_usage` instead.",
				Description: "Types 的 CBS 其中 此 快照 创建 从。",
			},
			"percent": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Snapshot creation progress percentage. 如果 快照 has 创建 successfully， constant 值 是 100。",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "创建时间 的 快照。",
			},
		},
	}
}

func resourceTencentCloudCbsSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_snapshot.create")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		cbsService = CbsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request    = cbs.NewCreateSnapshotRequest()
		response   = cbs.NewCreateSnapshotResponse()
		snapshotId string
	)

	if v, ok := d.GetOk("storage_id"); ok {
		request.DiskId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("snapshot_name"); ok {
		request.SnapshotName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("disk_usage"); ok {
		request.DiskUsage = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCbsClient().CreateSnapshot(request)
		if e != nil {
			return tccommon.RetryError(e, tccommon.InternalError)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create cbs snapshot, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cbs snapshot failed, reason:%s\n ", logId, err.Error())
		return err
	}

	if response.Response.SnapshotId == nil {
		return fmt.Errorf("SnapshotId is nil.")
	}

	snapshotId = *response.Response.SnapshotId
	d.SetId(snapshotId)

	// wait
	err = resource.Retry(20*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		snapshot, e := cbsService.DescribeSnapshotById(ctx, snapshotId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if snapshot == nil {
			return resource.RetryableError(fmt.Errorf("cbs snapshot is nil"))
		}

		if *snapshot.SnapshotState == CBS_SNAPSHOT_STATUS_CREATING {
			return resource.RetryableError(fmt.Errorf("cbs snapshot status is still %s", *snapshot.SnapshotState))
		}

		if *snapshot.SnapshotState == CBS_SNAPSHOT_STATUS_NORMAL {
			return nil
		}

		return resource.NonRetryableError(fmt.Errorf("cbs snapshot status is %s, we won't wait for it finish.", *snapshot.SnapshotState))
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cbs snapshot failed, reason:%s\n ", logId, err.Error())
		return err
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("cvm", "snapshot", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudCbsSnapshotRead(d, meta)
}

func resourceTencentCloudCbsSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_snapshot.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		tcClient   = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		cbsService = CbsService{client: tcClient}
		snapshotId = d.Id()
		snapshot   *cbs.Snapshot
		e          error
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		snapshot, e = cbsService.DescribeSnapshotById(ctx, snapshotId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read cbs snapshot failed, reason:%s\n ", logId, err.Error())
		return err
	}

	if snapshot == nil {
		d.SetId("")
		return nil
	}

	if snapshot.DiskId != nil {
		_ = d.Set("storage_id", snapshot.DiskId)
	}

	if snapshot.SnapshotName != nil {
		_ = d.Set("snapshot_name", snapshot.SnapshotName)
	}

	if snapshot.DiskUsage != nil {
		_ = d.Set("disk_usage", snapshot.DiskUsage)
	}

	if snapshot.DiskSize != nil {
		_ = d.Set("storage_size", snapshot.DiskSize)
	}

	if snapshot.SnapshotState != nil {
		_ = d.Set("snapshot_status", snapshot.SnapshotState)
	}

	if snapshot.DiskUsage != nil {
		_ = d.Set("disk_type", snapshot.DiskUsage)
	}

	if snapshot.Percent != nil {
		_ = d.Set("percent", snapshot.Percent)
	}

	if snapshot.CreateTime != nil {
		_ = d.Set("create_time", snapshot.CreateTime)
	}

	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "cvm", "snapshot", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)
	return nil
}

func resourceTencentCloudCbsSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_snapshot.update")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		tcClient   = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		snapshotId = d.Id()
	)

	if d.HasChange("snapshot_name") {
		snapshotName := d.Get("snapshot_name").(string)
		cbsService := CbsService{client: tcClient}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			e := cbsService.ModifySnapshotName(ctx, snapshotId, snapshotName)
			if e != nil {
				return tccommon.RetryError(e)
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update cbs snapshot name failed, reason:%s\n ", logId, err.Error())
			return err
		}
	}

	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("cvm", "snapshot", tcClient.Region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceTencentCloudCbsSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_snapshot.delete")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		cbsService = CbsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		snapshotId = d.Id()
	)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := cbsService.DeleteSnapshot(ctx, snapshotId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete cbs snapshot failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}
