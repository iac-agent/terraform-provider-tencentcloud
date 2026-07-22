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

func ResourceTencentCloudMpsImageSpriteTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsImageSpriteTemplateCreate,
		Read:   resourceTencentCloudMpsImageSpriteTemplateRead,
		Update: resourceTencentCloudMpsImageSpriteTemplateUpdate,
		Delete: resourceTencentCloudMpsImageSpriteTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"sample_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Sampling 类型，可选 值:Percent/Time。",
			},

			"sample_interval": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Sampling 间隔.当 SampleType 是 Percent，指定percentage 的 sampling 间隔.当 SampleType 是 Time，指定sampling 间隔 时间 （秒）。",
			},

			"row_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "数量 rows 在 small 镜像 在 sprite。",
			},

			"column_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "数量 columns 在 small 镜像 在 sprite。",
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Image sprite 模板名称，长度 限制: 64 字符。",
			},

			"width": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "最大 值 的 宽度 (或 long side) 的 small 镜像 在 sprite 镜像，取值范围：0 和 [128，4096]，单位: 像素.当 宽度 和 高度 是 both 0， resolution 是 same.当 宽度 是 0 和 高度 是 不 0，宽度 是 scaled proportionally.当 宽度 是 不 0 和 高度 是 0，高度 是 scaled proportionally.当 both 宽度 和 高度 是 不 0， resolution 是 指定 通过 用户默认值：0。",
			},

			"height": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "最大 值 的 高度 (或 short side) 的 small 镜像 在 sprite 镜像，取值范围：0 和 [128，4096]，单位: 像素.当 宽度 和 高度 是 both 0， resolution 是 same.当 宽度 是 0 和 高度 是 不 0，宽度 是 scaled proportionally.当 宽度 是 不 0 和 高度 是 0，高度 是 scaled proportionally.当 both 宽度 和 高度 是 不 0， resolution 是 指定 通过 用户默认值：0。",
			},

			"resolution_adaptive": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Adaptive resolution，可选 值:open: At 此 时间，宽度 表示 long side 的 视频，高度 表示 short side 的 视频.close: At 此 point，宽度 表示 宽度 的 视频，和 高度 表示 高度 的 视频.默认值：open。",
			},

			"fill_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filling 类型，当 aspect ratio 的 视频 流 配置 是 inconsistent 使用 aspect ratio 的 original 视频， processing 方法 对于 transcoding 是 filling. 可选 filling 类型:stretch: Stretching，stretching each frame 到 fill entire screen，其中 可能 cause transcoded 视频 到 是 squashed 或 stretched.black: Leave black，keep 视频 aspect ratio unchanged，和 fill rest 的 edge 使用 black.默认值：black。",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "模板描述 信息，长度 限制: 256 字符。",
			},

			"format": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Image 格式， 值 可以 是 jpg，png，webp. 默认为 jpg。",
			},
		},
	}
}

func resourceTencentCloudMpsImageSpriteTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_image_sprite_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = mps.NewCreateImageSpriteTemplateRequest()
		response   = mps.NewCreateImageSpriteTemplateResponse()
		definition uint64
	)
	if v, ok := d.GetOk("sample_type"); ok {
		request.SampleType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("sample_interval"); ok {
		request.SampleInterval = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("row_count"); ok {
		request.RowCount = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("column_count"); ok {
		request.ColumnCount = helper.IntUint64(v.(int))
	}

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

	if v, ok := d.GetOk("fill_type"); ok {
		request.FillType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	if v, ok := d.GetOk("format"); ok {
		request.Format = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().CreateImageSpriteTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps imageSpriteTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition = *response.Response.Definition
	d.SetId(helper.UInt64ToStr(definition))

	return resourceTencentCloudMpsImageSpriteTemplateRead(d, meta)
}

func resourceTencentCloudMpsImageSpriteTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_image_sprite_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	definition := d.Id()

	imageSpriteTemplate, err := service.DescribeMpsImageSpriteTemplateById(ctx, definition)
	if err != nil {
		return err
	}

	if imageSpriteTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsImageSpriteTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if imageSpriteTemplate.SampleType != nil {
		_ = d.Set("sample_type", imageSpriteTemplate.SampleType)
	}

	if imageSpriteTemplate.SampleInterval != nil {
		_ = d.Set("sample_interval", imageSpriteTemplate.SampleInterval)
	}

	if imageSpriteTemplate.RowCount != nil {
		_ = d.Set("row_count", imageSpriteTemplate.RowCount)
	}

	if imageSpriteTemplate.ColumnCount != nil {
		_ = d.Set("column_count", imageSpriteTemplate.ColumnCount)
	}

	if imageSpriteTemplate.Name != nil {
		_ = d.Set("name", imageSpriteTemplate.Name)
	}

	if imageSpriteTemplate.Width != nil {
		_ = d.Set("width", imageSpriteTemplate.Width)
	}

	if imageSpriteTemplate.Height != nil {
		_ = d.Set("height", imageSpriteTemplate.Height)
	}

	if imageSpriteTemplate.ResolutionAdaptive != nil {
		_ = d.Set("resolution_adaptive", imageSpriteTemplate.ResolutionAdaptive)
	}

	if imageSpriteTemplate.FillType != nil {
		_ = d.Set("fill_type", imageSpriteTemplate.FillType)
	}

	if imageSpriteTemplate.Comment != nil {
		_ = d.Set("comment", imageSpriteTemplate.Comment)
	}

	if imageSpriteTemplate.Format != nil {
		_ = d.Set("format", imageSpriteTemplate.Format)
	}

	return nil
}

func resourceTencentCloudMpsImageSpriteTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_image_sprite_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := mps.NewModifyImageSpriteTemplateRequest()

	definition := d.Id()

	request.Definition = helper.StrToUint64Point(definition)

	mutableArgs := []string{"sample_type", "sample_interval", "row_count", "column_count", "name", "width", "height", "resolution_adaptive", "fill_type", "comment", "format"}

	needChange := false

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		if v, ok := d.GetOk("sample_type"); ok {
			request.SampleType = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("sample_interval"); ok {
			request.SampleInterval = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("row_count"); ok {
			request.RowCount = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("column_count"); ok {
			request.ColumnCount = helper.IntUint64(v.(int))
		}

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

		if v, ok := d.GetOk("fill_type"); ok {
			request.FillType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}

		if v, ok := d.GetOk("format"); ok {
			request.Format = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().ModifyImageSpriteTemplate(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mps imageSpriteTemplate failed, reason:%+v", logId, err)
			return err
		}

	}

	return resourceTencentCloudMpsImageSpriteTemplateRead(d, meta)
}

func resourceTencentCloudMpsImageSpriteTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_image_sprite_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	definition := d.Id()

	if err := service.DeleteMpsImageSpriteTemplateById(ctx, definition); err != nil {
		return err
	}

	return nil
}
