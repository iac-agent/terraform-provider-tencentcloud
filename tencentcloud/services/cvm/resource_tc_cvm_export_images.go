package cvm

import (
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCvmExportImages() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmExportImagesCreate,
		Read:   resourceTencentCloudCvmExportImagesRead,
		Delete: resourceTencentCloudCvmExportImagesDelete,
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "COS 存储桶名称",
			},

			"image_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Image ID。",
			},

			"file_name_prefix": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Prefix 的 exported 文件。",
			},

			"export_format": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "格式 的 exported 镜像 文件. 有效值：RAW，QCOW2，VHD 和 VMDK. 默认值：RAW。",
			},

			"only_export_root_disk": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "是否export 仅 系统 磁盘。",
			},

			"dry_run": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Check 是否image 可以 是 exported。",
			},

			"role_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "角色 名称 (默认值：CVM_QcsRole). Before exporting images，make sure 角色 exists，和 它 has write 权限 到 COS。",
			},
		},
	}
}

func resourceTencentCloudCvmExportImagesCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_export_images.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request        = cvm.NewExportImagesRequest()
		imageId        string
		bucketName     string
		fileNamePrefix string
	)
	imageId = d.Get("image_id").(string)
	bucketName = d.Get("bucket_name").(string)
	fileNamePrefix = d.Get("file_name_prefix").(string)
	request.ImageIds = []*string{&imageId}
	request.BucketName = helper.String(bucketName)
	request.FileNamePrefixList = []*string{&fileNamePrefix}

	if v, ok := d.GetOk("export_format"); ok {
		request.ExportFormat = helper.String(v.(string))
	}

	if v, _ := d.GetOk("only_export_root_disk"); v != nil {
		request.OnlyExportRootDisk = helper.Bool(v.(bool))
	}

	if v, _ := d.GetOk("dry_run"); v != nil {
		request.DryRun = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("role_name"); ok {
		request.RoleName = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ExportImages(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm exportImages failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(imageId)

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"NORMAL"}, 20*tccommon.ReadRetryTimeout, time.Second, service.CvmSyncImagesStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCvmExportImagesRead(d, meta)
}

func resourceTencentCloudCvmExportImagesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_export_images.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmExportImagesDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_export_images.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
