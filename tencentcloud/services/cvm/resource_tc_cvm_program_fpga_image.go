package cvm

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCvmProgramFpgaImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmProgramFpgaImageCreate,
		Read:   resourceTencentCloudCvmProgramFpgaImageRead,
		Delete: resourceTencentCloudCvmProgramFpgaImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID 信息 的 实例。",
			},

			"fpga_url": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "COS URL 地址 的 FPGA 镜像 文件。",
			},

			"dbd_fs": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "DBDF 数量 FPGA card 在 实例，如果 left blank， FPGA 镜像 将 是 burned 到 all FPGA cards owned 通过 实例 通过 默认值。",
			},

			"dry_run": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Trial run，将 不 perform actual burning 操作， 默认为 False。",
			},
		},
	}
}

func resourceTencentCloudCvmProgramFpgaImageCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_program_fpga_image.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = cvm.NewProgramFpgaImageRequest()
	)
	instanceId := d.Get("instance_id").(string)
	request.InstanceId = helper.String(instanceId)

	if v, ok := d.GetOk("fpga_url"); ok {
		request.FPGAUrl = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dbd_fs"); ok {
		dBDFsSet := v.(*schema.Set).List()
		for i := range dBDFsSet {
			dBDFs := dBDFsSet[i].(string)
			request.DBDFs = append(request.DBDFs, &dBDFs)
		}
	}

	if v, _ := d.GetOk("dry_run"); v != nil {
		request.DryRun = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ProgramFpgaImage(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm programFpgaImage failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudCvmProgramFpgaImageRead(d, meta)
}

func resourceTencentCloudCvmProgramFpgaImageRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_program_fpga_image.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmProgramFpgaImageDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_program_fpga_image.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
