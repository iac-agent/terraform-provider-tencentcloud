package vod

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudVodImageSpriteTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVodImageSpriteTemplateCreate,
		Read:   resourceTencentCloudVodImageSpriteTemplateRead,
		Update: resourceTencentCloudVodImageSpriteTemplateUpdate,
		Delete: resourceTencentCloudVodImageSpriteTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"sample_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"Percent", "Time"}),
				Description:  "Sampling 类型 有效值：`Percent`，`Time`. `Percent`: 通过 percent. `Time`: 通过 时间间隔。",
			},
			"sample_interval": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Sampling 间隔. 如果 `sample_type` 是 `Percent`，sampling 将 是 performed 在 间隔 的 指定 percentage. 如果 `sample_type` 是 `Time`，sampling 将 是 performed 在 指定 时间间隔 （秒）。",
			},
			"row_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Subimage row count 的 镜像 sprite。",
			},
			"column_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Subimage 列 count 的 镜像 sprite。",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 64),
				Description:  "名称 时间 point screen capturing template. Length 限制: 64 字符。",
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 256),
				Description:  "模板描述 Length 限制: 256 字符。",
			},
			"fill_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "black",
				Description: "Fill refers 到 way 的 processing screenshot 当 its aspect ratio 是 different 从 该 的 来源 视频. following fill types 是 支持: `stretch`: stretch. screenshot 将 是 stretched frame 通过 frame 到 match aspect ratio 的 来源 视频，其中 可能 make screenshot shorter 或 longer; `black`: fill 使用 black. 此 选项 retains aspect ratio 的 来源 视频 对于 screenshot 和 fills unmatched area 使用 black color blocks. 默认值：`black`。",
			},
			"width": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Maximum 值 的 `宽度` (或 long side) 的 screenshot （像素）。 取值范围：0 和 [128，4,096]. 如果 both `宽度` 和 `高度` 是 `0`， resolution 将 是 same 作为 该 的 来源 视频; 如果 `宽度` 是 `0`，但 `高度` 是 不 `0`，宽度 将 是 proportionally scaled; 如果 `宽度` 是 不 `0`，但 `高度` 是 `0`，`高度` 将 是 proportionally scaled; 如果 both `宽度` 和 `高度` 是 不 `0`， 自定义 resolution 将 是 使用. 默认值：`0`。",
			},
			"height": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Maximum 值 的 `高度` (或 short side) 的 screenshot （像素）。 取值范围：0 和 [128，4,096]. 如果 both `宽度` 和 `高度` 是 `0`， resolution 将 是 same 作为 该 的 来源 视频; 如果 `宽度` 是 `0`，但 `高度` 是 不 `0`，`宽度` 将 是 proportionally scaled; 如果 `宽度` 是 不 `0`，但 `高度` 是 `0`，`高度` 将 是 proportionally scaled; 如果 both `宽度` 和 `高度` 是 不 `0`， 自定义 resolution 将 是 使用. 默认值：`0`。",
			},
			"resolution_adaptive": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Resolution adaption. 有效值：`true`,`false`. `true`: 已启用 In 此 case，`宽度` 表示 long side 的 视频，while `高度` short side; `false`: 已禁用 In 此 case，`宽度` 表示 宽度 的 视频，while `高度` 高度. 默认值：`true`。",
			},
			"sub_app_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "VOD [应用](https://intl.云.tencent.com/document/product/266/14574) ID. For customers who activate VOD 服务 从 December 25，2023，如果 they want 到 访问 resources 在 VOD 应用 (whether 它's 默认值 应用 或 newly 创建 一个)，they 必须 fill 在 此 字段 使用 应用 ID。",
			},
			"format": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Image 格式, 有效 值:\n" +
					"- jpg: jpg format;\n" +
					"- png: png format;\n" +
					"- webp: webp format;\n" +
					"Default value: jpg.",
			},
			// computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "创建时间 的 template 在 ISO date 格式",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "最后修改时间 的 template 在 ISO date 格式",
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "模板 类型, 值 范围:\n" +
					"- Preset: system preset template;\n" +
					"- Custom: user-defined templates.",
			},
		},
	}
}

func resourceTencentCloudVodImageSpriteTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_image_sprite_template.create")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = vod.NewCreateImageSpriteTemplateRequest()
	)

	request.SampleType = helper.String(d.Get("sample_type").(string))
	request.SampleInterval = helper.IntUint64(d.Get("sample_interval").(int))
	request.RowCount = helper.IntUint64(d.Get("row_count").(int))
	request.ColumnCount = helper.IntUint64(d.Get("column_count").(int))
	request.Name = helper.String(d.Get("name").(string))
	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}
	request.FillType = helper.String((d.Get("fill_type")).(string))
	request.Width = helper.IntUint64(d.Get("width").(int))
	request.Height = helper.IntUint64(d.Get("height").(int))
	request.ResolutionAdaptive = helper.String(RESOLUTION_ADAPTIVE_TO_STRING[d.Get("resolution_adaptive").(bool)])
	var resourceId string
	if v, ok := d.GetOk("sub_app_id"); ok {
		subAppId := v.(int)
		resourceId += helper.IntToStr(subAppId)
		resourceId += tccommon.FILED_SP
		request.SubAppId = helper.IntUint64(subAppId)
	}
	if v, ok := d.GetOk("format"); ok {
		request.Format = helper.String(v.(string))
	}

	var response *vod.CreateImageSpriteTemplateResponse
	var err error
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().CreateImageSpriteTemplate(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "FailedOperation" && sdkError.Message == "invalid vod user" {
					return resource.RetryableError(err)
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if response == nil || response.Response == nil {
		return fmt.Errorf("for image sprite template creation, response is nil")
	}
	resourceId += strconv.FormatUint(*response.Response.Definition, 10)
	d.SetId(resourceId)

	return resourceTencentCloudVodImageSpriteTemplateRead(d, meta)
}

func resourceTencentCloudVodImageSpriteTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_image_sprite_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		subAppId   int
		definition string
		client     = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		vodService = VodService{client: client}
	)
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) == 2 {
		subAppId = helper.StrToInt(idSplit[0])
		definition = idSplit[1]
	} else {
		definition = d.Id()
	}
	// waiting for refreshing cache
	time.Sleep(30 * time.Second)
	template, has, err := vodService.DescribeImageSpriteTemplatesById(ctx, definition, subAppId)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("sample_type", template.SampleType)
	_ = d.Set("sample_interval", template.SampleInterval)
	_ = d.Set("row_count", template.RowCount)
	_ = d.Set("column_count", template.ColumnCount)
	_ = d.Set("name", template.Name)
	_ = d.Set("comment", template.Comment)
	_ = d.Set("fill_type", template.FillType)
	_ = d.Set("width", template.Width)
	_ = d.Set("height", template.Height)
	_ = d.Set("resolution_adaptive", *template.ResolutionAdaptive == "open")
	_ = d.Set("create_time", template.CreateTime)
	_ = d.Set("update_time", template.UpdateTime)
	_ = d.Set("format", template.Format)
	_ = d.Set("type", template.Type)
	if subAppId != 0 {
		_ = d.Set("sub_app_id", subAppId)
	}

	return nil
}

func resourceTencentCloudVodImageSpriteTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_image_sprite_template.update")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = vod.NewModifyImageSpriteTemplateRequest()
		changeFlag = false
		subAppId   int
		definition string
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) == 2 {
		subAppId = helper.StrToInt(idSplit[0])
		definition = idSplit[1]
		request.SubAppId = helper.IntUint64(subAppId)
	} else {
		definition = d.Id()
		if v, ok := d.GetOk("sub_app_id"); ok {
			request.SubAppId = helper.IntUint64(v.(int))
		}
	}
	request.Definition = helper.StrToUint64Point(definition)
	immutableArgs := []string{"sub_app_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("sample_type") {
		changeFlag = true
		request.SampleType = helper.String(d.Get("sample_type").(string))
	}
	if d.HasChange("sample_interval") {
		changeFlag = true
		request.SampleInterval = helper.IntUint64(d.Get("sample_interval").(int))
	}
	if d.HasChange("row_count") {
		changeFlag = true
		request.RowCount = helper.IntUint64(d.Get("row_count").(int))
	}
	if d.HasChange("column_count") {
		changeFlag = true
		request.ColumnCount = helper.IntUint64(d.Get("column_count").(int))
	}
	if d.HasChange("name") {
		changeFlag = true
		request.Name = helper.String(d.Get("name").(string))
	}
	if d.HasChange("comment") {
		changeFlag = true
		request.Comment = helper.String(d.Get("comment").(string))
	}
	if d.HasChange("fill_type") {
		changeFlag = true
		request.FillType = helper.String(d.Get("fill_type").(string))
	}
	if d.HasChange("width") || d.HasChange("height") || d.HasChange("resolution_adaptive") {
		changeFlag = true
		request.Width = helper.IntUint64(d.Get("width").(int))
		request.Height = helper.IntUint64(d.Get("height").(int))
		request.ResolutionAdaptive = helper.String(RESOLUTION_ADAPTIVE_TO_STRING[d.Get("resolution_adaptive").(bool)])
	}
	if d.HasChange("format") {
		changeFlag = true
		request.Format = helper.String(d.Get("format").(string))
	}

	if changeFlag {
		var err error
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			_, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().ModifyImageSpriteTemplate(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
				return tccommon.RetryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}

		return resourceTencentCloudVodImageSpriteTemplateRead(d, meta)
	}

	return nil
}

func resourceTencentCloudVodImageSpriteTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_image_sprite_template.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		subAppId   int
		definition string
	)
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) == 2 {
		subAppId = helper.StrToInt(idSplit[0])
		definition = idSplit[1]
	} else {
		definition = d.Id()
		if v, ok := d.GetOk("sub_app_id"); ok {
			subAppId = v.(int)
		}
	}
	vodService := VodService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	if err := vodService.DeleteImageSpriteTemplate(ctx, definition, uint64(subAppId)); err != nil {
		return err
	}

	return nil
}
