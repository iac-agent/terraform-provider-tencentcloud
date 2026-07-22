package mps

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMpsSnapshotByTimeoffsetTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsSnapshotByTimeoffsetTemplateCreate,
		Read:   resourceTencentCloudMpsSnapshotByTimeoffsetTemplateRead,
		Update: resourceTencentCloudMpsSnapshotByTimeoffsetTemplateUpdate,
		Delete: resourceTencentCloudMpsSnapshotByTimeoffsetTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Snapshot 通过 timeoffset 模板名称，长度 限制: 64 字符。",
			},

			"width": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "最大 值 的 快照 宽度 (或 long side)，取值范围：0 和 [128，4096]，单位: 像素.当 宽度 和 高度 是 both 0， resolution 是 same.当 宽度 是 0 和 高度 是 不 0，宽度 是 scaled proportionally.当 宽度 是 不 0 和 高度 是 0，高度 是 scaled proportionally.当 both 宽度 和 高度 是 不 0， resolution 是 指定 通过 用户默认值：0。",
			},

			"height": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "最大 值 的 快照 高度 (或 short side)，取值范围：0 和 [128，4096]，单位: 像素.当 宽度 和 高度 是 both 0， resolution 是 same.当 宽度 是 0 和 高度 是 不 0，宽度 是 scaled proportionally.当 宽度 是 不 0 和 高度 是 0，高度 是 scaled proportionally.当 both 宽度 和 高度 是 不 0， resolution 是 指定 通过 用户默认值：0。",
			},

			"resolution_adaptive": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Adaptive resolution，可选 值:open: At 此 时间，宽度 表示 long side 的 视频，高度 表示 short side 的 视频.close: At 此 point，宽度 表示 宽度 的 视频，和 高度 表示 高度 的 视频.默认值：open。",
			},

			"format": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Image 格式， 值 可以 是 jpg，png，webp. 默认为 jpg。",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "模板描述 信息，长度 限制: 256 字符。",
			},

			"fill_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filling 类型，当 aspect ratio 的 视频 流 配置 是 inconsistent 使用 aspect ratio 的 original 视频， processing 方法 对于 transcoding 是 filling. 可选 filling 类型:stretch: Stretching，stretching each frame 到 fill entire screen，其中 可能 cause transcoded 视频 到 是 squashed 或 stretched.black: Leave black，keep 视频 aspect ratio unchanged，和 fill rest 的 edge 使用 black.white: Leave blank，keep aspect ratio 的 视频，和 fill rest 的 edge 使用 white.gauss: Gaussian blur，keep aspect ratio 的 视频 unchanged，和 使用 Gaussian blur 对于 rest 的 edge.默认值：black。",
			},
		},
	}
}

func resourceTencentCloudMpsSnapshotByTimeoffsetTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_snapshot_by_timeoffset_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = mps.NewCreateSnapshotByTimeOffsetTemplateRequest()
		response   = mps.NewCreateSnapshotByTimeOffsetTemplateResponse()
		definition uint64
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("width"); ok {
		request.Width = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("height"); ok {
		request.Height = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("resolution_adaptive"); ok {
		request.ResolutionAdaptive = helper.String(v.(string))
	}

	if v, ok := d.GetOk("format"); ok {
		request.Format = helper.String(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	if v, ok := d.GetOk("fill_type"); ok {
		request.FillType = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().CreateSnapshotByTimeOffsetTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps snapshotByTimeoffsetTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition = *response.Response.Definition
	d.SetId(helper.UInt64ToStr(definition))

	return resourceTencentCloudMpsSnapshotByTimeoffsetTemplateRead(d, meta)
}

func resourceTencentCloudMpsSnapshotByTimeoffsetTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_snapshot_by_timeoffset_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	definition := d.Id()

	snapshotByTimeoffsetTemplate, err := service.DescribeMpsSnapshotByTimeoffsetTemplateById(ctx, definition)
	if err != nil {
		return err
	}

	if snapshotByTimeoffsetTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsSnapshotByTimeoffsetTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if snapshotByTimeoffsetTemplate.Name != nil {
		_ = d.Set("name", snapshotByTimeoffsetTemplate.Name)
	}

	if snapshotByTimeoffsetTemplate.Width != nil {
		_ = d.Set("width", snapshotByTimeoffsetTemplate.Width)
	}

	if snapshotByTimeoffsetTemplate.Height != nil {
		_ = d.Set("height", snapshotByTimeoffsetTemplate.Height)
	}

	if snapshotByTimeoffsetTemplate.ResolutionAdaptive != nil {
		_ = d.Set("resolution_adaptive", snapshotByTimeoffsetTemplate.ResolutionAdaptive)
	}

	if snapshotByTimeoffsetTemplate.Format != nil {
		_ = d.Set("format", snapshotByTimeoffsetTemplate.Format)
	}

	if snapshotByTimeoffsetTemplate.Comment != nil {
		_ = d.Set("comment", snapshotByTimeoffsetTemplate.Comment)
	}

	if snapshotByTimeoffsetTemplate.FillType != nil {
		_ = d.Set("fill_type", snapshotByTimeoffsetTemplate.FillType)
	}

	return nil
}

func resourceTencentCloudMpsSnapshotByTimeoffsetTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_snapshot_by_timeoffset_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	needChange := false

	request := mps.NewModifySnapshotByTimeOffsetTemplateRequest()

	definition := d.Id()

	request.Definition = helper.StrToUint64Point(definition)

	mutableArgs := []string{"name", "width", "height", "resolution_adaptive", "format", "comment", "fill_type"}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("width"); ok {
			request.Width = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("height"); ok {
			request.Height = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("resolution_adaptive"); ok {
			request.ResolutionAdaptive = helper.String(v.(string))
		}

		if v, ok := d.GetOk("format"); ok {
			request.Format = helper.String(v.(string))
		}

		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}

		if v, ok := d.GetOk("fill_type"); ok {
			request.FillType = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().ModifySnapshotByTimeOffsetTemplate(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mps snapshotByTimeoffsetTemplate failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudMpsSnapshotByTimeoffsetTemplateRead(d, meta)
}

func resourceTencentCloudMpsSnapshotByTimeoffsetTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_snapshot_by_timeoffset_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	definition := d.Id()

	if err := service.DeleteMpsSnapshotByTimeoffsetTemplateById(ctx, definition); err != nil {
		return err
	}

	return nil
}
