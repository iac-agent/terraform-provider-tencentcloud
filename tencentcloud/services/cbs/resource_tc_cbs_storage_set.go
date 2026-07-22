package cbs

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCbsStorageSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCbsStorageSetCreate,
		Read:   resourceTencentCloudCbsStorageSetRead,
		Update: resourceTencentCloudCbsStorageSetUpdate,
		Delete: resourceTencentCloudCbsStorageSetDelete,

		Schema: map[string]*schema.Schema{
			"storage_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "类型 CBS medium. 有效值：CLOUD_BASIC: HDD 云 磁盘，CLOUD_PREMIUM: Premium Cloud Storage，CLOUD_BSSD: General Purpose SSD，CLOUD_SSD: SSD，CLOUD_HSSD: Enhanced SSD，CLOUD_TSSD: Tremendous SSD。",
			},
			"storage_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Volume 的 CBS，和 单位 是 GB。",
			},
			"disk_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "数量 disks 到 是 purchased. Default 1。",
			},
			"charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CBS_CHARGE_TYPE_POSTPAID,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CBS_CHARGE_TYPE),
				Description:  "charge 类型 CBS 实例. Support `POSTPAID_BY_HOUR` 和 `DEDICATED_CLUSTER_PAID`. 默认为 `POSTPAID_BY_HOUR`。",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "可用 可用区 该 CBS 实例 locates 在。",
			},
			"dedicated_cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Exclusive 集群 ID",
			},
			"storage_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(2, 60),
				Description:  "名称 CBS. 最大 长度 可以 不 exceed 60 bytes。",
			},
			"snapshot_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID 快照. 如果 指定，创建 CBS 通过 此 快照。",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "ID 项目 到 其中 实例 belongs。",
			},
			"kms_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "可选 参数. 当 purchasing 加密 磁盘，customize 键 当 此 参数 是 passed 在， `encrypt` 参数 need 是 集合。",
			},
			"encrypt": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "表示是否CBS 是 encrypted。",
			},
			"throughput_performance": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Add extra performance 到 数据 磁盘. Only works 当 磁盘 类型 是 `CLOUD_TSSD` 或 `CLOUD_HSSD`。",
			},
			// computed
			"storage_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "状态 CBS. 有效值：UNATTACHED，ATTACHING，ATTACHED，DETACHING，EXPANDING，ROLLBACKING，TORECYCLE 和 DUMPING。",
			},
			"attached": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "表示是否CBS 是 mounted CVM。",
			},
			"disk_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "磁盘 ID 列表。",
			},
		},
	}
}

func resourceTencentCloudCbsStorageSetCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_storage_set.create")()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		request   = cbs.NewCreateDisksRequest()
		diskCount int
	)

	request.DiskName = helper.String(d.Get("storage_name").(string))
	request.DiskType = helper.String(d.Get("storage_type").(string))
	request.DiskSize = helper.IntUint64(d.Get("storage_size").(int))
	if v, ok := d.GetOk("disk_count"); ok {
		diskCount = v.(int)
		request.DiskCount = helper.Uint64(uint64(diskCount))
	}

	request.Placement = &cbs.Placement{
		Zone: helper.String(d.Get("availability_zone").(string)),
	}

	if v, ok := d.GetOk("dedicated_cluster_id"); ok {
		request.Placement.DedicatedClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		request.Placement.ProjectId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("snapshot_id"); ok {
		request.SnapshotId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		request.KmsKeyId = helper.String(v.(string))
	}

	if _, ok := d.GetOk("encrypt"); ok {
		request.Encrypt = helper.String("ENCRYPT")
	}

	if v, ok := d.GetOk("throughput_performance"); ok {
		request.ThroughputPerformance = helper.IntUint64(v.(int))
	}

	chargeType := d.Get("charge_type").(string)

	request.DiskChargeType = &chargeType

	storageIds := make([]*string, 0)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCbsClient().CreateDisks(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())

			ee, ok := e.(*sdkErrors.TencentCloudSDKError)
			if ok && tccommon.IsContains(CVM_RETRYABLE_ERROR, ee.Code) {
				time.Sleep(1 * time.Second) // 需要重试的话，等待1s进行重试
				return resource.RetryableError(fmt.Errorf("cbs create error: %s, retrying", e.Error()))
			}

			return resource.NonRetryableError(e)
		}

		if len(response.Response.DiskIdSet) < diskCount {
			err := fmt.Errorf("number of instances is less than %s", strconv.Itoa(diskCount))
			return resource.NonRetryableError(err)
		}

		storageIds = response.Response.DiskIdSet
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cbs failed, reason:%s\n ", logId, err.Error())
		return err
	}

	_ = d.Set("disk_ids", storageIds)
	d.SetId(helper.StrListToStr(storageIds))

	return nil
}

func resourceTencentCloudCbsStorageSetRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_storage_set.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		cbsService = CbsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		storageSet []*cbs.Disk
		errRet     error
		storageId  = d.Id()
	)

	storageSet, errRet = cbsService.DescribeDiskSetByIds(ctx, storageId)
	if errRet != nil {
		return errRet
	}

	if storageSet == nil {
		d.SetId("")
		return nil
	}

	storage := storageSet[0]

	_ = d.Set("disk_count", len(storageSet))
	_ = d.Set("storage_type", storage.DiskType)
	_ = d.Set("storage_size", storage.DiskSize)
	_ = d.Set("availability_zone", storage.Placement.Zone)
	_ = d.Set("dedicated_cluster_id", storage.Placement.DedicatedClusterId)
	_ = d.Set("storage_name", d.Get("storage_name"))
	_ = d.Set("project_id", storage.Placement.ProjectId)
	_ = d.Set("encrypt", storage.Encrypt)
	_ = d.Set("storage_status", storage.DiskState)
	_ = d.Set("attached", storage.Attached)
	_ = d.Set("charge_type", storage.DiskChargeType)
	_ = d.Set("throughput_performance", storage.ThroughputPerformance)

	if storage.KmsKeyId != nil {
		_ = d.Set("kms_key_id", storage.KmsKeyId)
	}

	return nil
}

func resourceTencentCloudCbsStorageSetUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_storage_set.update")()

	return fmt.Errorf("`tencentcloud_cbs_storage_set` do not support change now.")
}

func resourceTencentCloudCbsStorageSetDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cbs_storage_set.delete")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		storageId  = d.Id()
		cbsService = CbsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := cbsService.DeleteDiskSetByIds(ctx, storageId)
		if e != nil {
			log.Printf("[CRITAL][first delete]%s api[%s] fail, reason[%s]\n", logId, "delete", e.Error())
			ee, ok := e.(*sdkErrors.TencentCloudSDKError)
			if ok && tccommon.IsContains(CVM_RETRYABLE_ERROR, ee.Code) {
				time.Sleep(1 * time.Second) // 需要重试的话，等待1s进行重试
				return resource.RetryableError(fmt.Errorf("[first delete]cvm delete error: %s, retrying", ee.Error()))
			}

			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete cbs failed, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}
