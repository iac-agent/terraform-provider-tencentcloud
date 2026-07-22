package cvm

import (
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCvmSyncImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmSyncImageCreate,
		Read:   resourceTencentCloudCvmSyncImageRead,
		Delete: resourceTencentCloudCvmSyncImageDelete,

		Schema: map[string]*schema.Schema{
			"image_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Image ID. 指定 镜像 必须 meet following requirement: images 必须 是 在 `NORMAL` state。",
			},

			"destination_regions": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "列表 destination regions 对于 synchronization. Limits: It 必须 是 有效 地域 For 自定义 镜像， destination 地域 不能 是 来源 地域 For shared 镜像， destination 地域 必须 是 来源 地域，其中 表示to create copy 的 镜像 作为 自定义 镜像 在 same 地域",
			},

			"dry_run": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Checks whether 镜像 synchronization 可以 是 initiated。",
			},

			"image_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Destination 镜像 名称",
			},

			"image_set_required": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "是否return ID 镜像 创建 在 destination 地域",
			},

			"encrypt": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "是否synchronize 作为 encrypted 自定义 镜像. 默认值为 `false`. Synchronization 到 encrypted 自定义 镜像 是 仅 支持 within same 地域",
			},

			"kms_key_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "KMS 键 ID 使用 当 synchronizing 到 encrypted 自定义 镜像. 此 参数 是 有效 仅 synchronizing 到 encrypted 镜像. 如果 KmsKeyId 是 不 指定， 默认值 CBS 云 product KMS 键 是 使用。",
			},

			"image_set": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image ID。",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 的 镜像。",
						},
					},
				},
				Description: "ID 镜像 创建 在 destination 地域",
			},
		},
	}
}

func resourceTencentCloudCvmSyncImageCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_sync_image.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := cvm.NewSyncImagesRequest()
	response := cvm.NewSyncImagesResponse()
	imageId := d.Get("image_id").(string)
	request.ImageIds = []*string{&imageId}

	if v, ok := d.GetOk("destination_regions"); ok {
		destinationRegionsSet := v.(*schema.Set).List()
		for i := range destinationRegionsSet {
			destinationRegions := destinationRegionsSet[i].(string)
			request.DestinationRegions = append(request.DestinationRegions, &destinationRegions)
		}
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request.DryRun = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("image_name"); ok {
		request.ImageName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("image_set_required"); ok {
		request.ImageSetRequired = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("encrypt"); ok {
		request.Encrypt = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		request.KmsKeyId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().SyncImages(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm syncImages failed, reason:%+v", logId, err)
		return err
	}

	if response == nil || response.Response == nil || response.Response.ImageSet == nil {
		err = fmt.Errorf("Response is nil")
		return err
	}

	d.SetId(imageId)

	imageSetList := []interface{}{}
	for _, image := range response.Response.ImageSet {
		imageMap := map[string]interface{}{}

		if image.ImageId != nil {
			imageMap["image_id"] = image.ImageId
		}

		if image.Region != nil {
			imageMap["region"] = image.Region
		}

		imageSetList = append(imageSetList, imageMap)
	}

	_ = d.Set("image_set", imageSetList)

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"NORMAL"}, 20*tccommon.ReadRetryTimeout, time.Second, service.CvmSyncImagesStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCvmSyncImageRead(d, meta)
}

func resourceTencentCloudCvmSyncImageRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_sync_image.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmSyncImageDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_sync_image.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
