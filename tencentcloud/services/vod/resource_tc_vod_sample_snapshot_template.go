package vod

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVodSampleSnapshotTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVodSampleSnapshotTemplateCreate,
		Read:   resourceTencentCloudVodSampleSnapshotTemplateRead,
		Update: resourceTencentCloudVodSampleSnapshotTemplateUpdate,
		Delete: resourceTencentCloudVodSampleSnapshotTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"sample_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Sampled screencapturing 类型 有效值：Percent: 通过 percent. Time: 通过 时间间隔。",
			},

			"sample_interval": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Sampling 间隔. 如果 `SampleType` 是 `Percent`，sampling 将 是 performed 在 间隔 的 指定 percentage. 如果 `SampleType` 是 `Time`，sampling 将 是 performed 在 指定 时间间隔 （秒）。",
			},

			"sub_app_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "VOD [应用](https://intl.云.tencent.com/document/product/266/14574) ID. For customers who activate VOD 服务 从 December 25，2023，如果 they want 到 访问 resources 在 VOD 应用 (whether 它's 默认值 应用 或 newly 创建 一个)，they 必须 fill 在 此 字段 使用 应用 ID。",
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "名称 sampled screencapturing template. Length 限制: 64 字符。",
			},

			"width": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Maximum 值 的 宽度 (或 long side) 的 screenshot （像素）。 取值范围：0 和 [128，4,096]. 如果 both `宽度` 和 `高度` 是 0， resolution 将 是 same 作为 该 的 来源 视频; 如果 `宽度` 是 0，但 `高度` 是 不 0，`宽度` 将 是 proportionally scaled; 如果 `宽度` 是 不 0，但 `高度` 是 0，`高度` 将 是 proportionally scaled; 如果 both `宽度` 和 `高度` 是 不 0， 自定义 resolution 将 是 使用.默认值：0。",
			},

			"height": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Maximum 值 的 高度 (或 short side) 的 screenshot （像素）。 取值范围：0 和 [128，4,096]. 如果 both `宽度` 和 `高度` 是 0， resolution 将 是 same 作为 该 的 来源 视频; 如果 `宽度` 是 0，但 `高度` 是 不 0，`宽度` 将 是 proportionally scaled; 如果 `宽度` 是 不 0，但 `高度` 是 0，`高度` 将 是 proportionally scaled; 如果 both `宽度` 和 `高度` 是 不 0， 自定义 resolution 将 是 使用.默认值：0。",
			},

			"resolution_adaptive": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Resolution adaption. 有效值：open: 已启用 In 此 case，`宽度` 表示 long side 的 视频，while `高度` short side; close: 已禁用 In 此 case，`宽度` 表示 宽度 的 视频，while `高度` 高度.默认值：open。",
			},

			"format": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Image 格式 有效值：jpg，png. 默认值：jpg。",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "模板描述 Length 限制: 256 字符。",
			},

			"fill_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Fill 类型 Fill refers 到 way 的 processing screenshot 当 its aspect ratio 是 different 从 该 的 来源 视频. following fill types 是 支持: stretch: stretch. screenshot 将 是 stretched frame 通过 frame 到 match aspect ratio 的 来源 视频，其中 可能 make screenshot shorter 或 longer; black: fill 使用 black. 此 选项 retains aspect ratio 的 来源 视频 对于 screenshot 和 fills unmatched area 使用 black color blocks. white: fill 使用 white. 此 选项 retains aspect ratio 的 来源 视频 对于 screenshot 和 fills unmatched area 使用 white color blocks. gauss: fill 使用 Gaussian blur. 此 选项 retains aspect ratio 的 来源 视频 对于 screenshot 和 fills unmatched area 使用 Gaussian blur.默认值：black。",
			},
		},
	}
}

func resourceTencentCloudVodSampleSnapshotTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_sample_snapshot_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = vod.NewCreateSampleSnapshotTemplateRequest()
		response = vod.NewCreateSampleSnapshotTemplateResponse()
		subAppId string
	)

	if v, ok := d.GetOk("sample_type"); ok {
		request.SampleType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("sample_interval"); ok {
		request.SampleInterval = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("sub_app_id"); ok {
		subAppId = helper.IntToStr(v.(int))
		request.SubAppId = helper.IntUint64(v.(int))
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
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().CreateSampleSnapshotTemplate(request)
		if e != nil {
			if sdkError, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "FailedOperation" && sdkError.Message == "invalid vod user" {
					return resource.RetryableError(e)
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), e.Error())
			return resource.NonRetryableError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vod sampleSnapshotTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition := *response.Response.Definition

	d.SetId(subAppId + tccommon.FILED_SP + helper.UInt64ToStr(definition))

	return resourceTencentCloudVodSampleSnapshotTemplateRead(d, meta)
}

func resourceTencentCloudVodSampleSnapshotTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_sample_snapshot_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

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
	}
	sampleSnapshotTemplate, err := service.DescribeVodSampleSnapshotTemplateById(ctx, uint64(subAppId), helper.StrToUInt64(definition))
	if err != nil {
		return err
	}

	if sampleSnapshotTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VodSampleSnapshotTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if subAppId != 0 {
		_ = d.Set("sub_app_id", subAppId)
	}
	if sampleSnapshotTemplate.SampleType != nil {
		_ = d.Set("sample_type", sampleSnapshotTemplate.SampleType)
	}

	if sampleSnapshotTemplate.SampleInterval != nil {
		_ = d.Set("sample_interval", sampleSnapshotTemplate.SampleInterval)
	}

	if sampleSnapshotTemplate.Name != nil {
		_ = d.Set("name", sampleSnapshotTemplate.Name)
	}

	if sampleSnapshotTemplate.Width != nil {
		_ = d.Set("width", sampleSnapshotTemplate.Width)
	}

	if sampleSnapshotTemplate.Height != nil {
		_ = d.Set("height", sampleSnapshotTemplate.Height)
	}

	if sampleSnapshotTemplate.ResolutionAdaptive != nil {
		_ = d.Set("resolution_adaptive", sampleSnapshotTemplate.ResolutionAdaptive)
	}

	if sampleSnapshotTemplate.Format != nil {
		_ = d.Set("format", sampleSnapshotTemplate.Format)
	}

	if sampleSnapshotTemplate.Comment != nil {
		_ = d.Set("comment", sampleSnapshotTemplate.Comment)
	}

	if sampleSnapshotTemplate.FillType != nil {
		_ = d.Set("fill_type", sampleSnapshotTemplate.FillType)
	}

	return nil
}

func resourceTencentCloudVodSampleSnapshotTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_sample_snapshot_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vod.NewModifySampleSnapshotTemplateRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("sample snapshot id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[0]
	definition := idSplit[1]

	immutableArgs := []string{"sub_app_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.Definition = helper.StrToUint64Point(definition)
	request.SubAppId = helper.StrToUint64Point(subAppId)

	if d.HasChange("sample_type") || d.HasChange("sample_interval") || d.HasChange("name") || d.HasChange("width") || d.HasChange("height") || d.HasChange("resolution_adaptive") || d.HasChange("format") || d.HasChange("comment") || d.HasChange("fill_type") {
		if v, ok := d.GetOk("sample_type"); ok {
			request.SampleType = helper.String(v.(string))
		}
		if v, ok := d.GetOkExists("sample_interval"); ok {
			request.SampleInterval = helper.IntUint64(v.(int))
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
		if v, ok := d.GetOk("format"); ok {
			request.Format = helper.String(v.(string))
		}
		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}
		if v, ok := d.GetOk("fill_type"); ok {
			request.FillType = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().ModifySampleSnapshotTemplate(request)
		if e != nil {
			if sdkError, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "FailedOperation" && sdkError.Message == "invalid vod user" {
					return resource.RetryableError(e)
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), e.Error())
			return resource.NonRetryableError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update vod sampleSnapshotTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudVodSampleSnapshotTemplateRead(d, meta)
}

func resourceTencentCloudVodSampleSnapshotTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_sample_snapshot_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("sample snapshot id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[0]
	definition := idSplit[1]

	if err := service.DeleteVodSampleSnapshotTemplateById(ctx, helper.StrToUInt64(subAppId), helper.StrToUInt64(definition)); err != nil {
		return err
	}

	return nil
}
