package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoJustInTimeTranscodeTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoJustInTimeTranscodeTemplateCreate,
		Read:   resourceTencentCloudTeoJustInTimeTranscodeTemplateRead,
		Delete: resourceTencentCloudTeoJustInTimeTranscodeTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID.",
			},

			"template_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Just-in-time transcode template name.",
			},

			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Template description.",
			},

			"video_stream_switch": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Enable video stream switch, valid values: on, off.",
			},

			"audio_stream_switch": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Enable audio stream switch, valid values: on, off.",
			},

			"video_template": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Video stream configuration parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"codec": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Video stream encoding format, valid values: H.264, H.265.",
						},
						"fps": {
							Type:        schema.TypeFloat,
							Optional:    true,
							ForceNew:    true,
							Description: "Video frame rate, range: [0, 30].",
						},
						"bitrate": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "Video stream bitrate, range: 0 and [128, 10000].",
						},
						"resolution_adaptive": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Resolution adaptive, valid values: open, close.",
						},
						"width": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "Maximum width of video stream, range: 0 and [128, 1920].",
						},
						"height": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "Maximum height of video stream, range: 0 and [128, 1080].",
						},
						"fill_type": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Fill type, valid values: stretch, black, white, gauss.",
						},
					},
				},
			},

			"audio_template": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Audio stream configuration parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"codec": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Audio stream encoding format, valid value: libfdk_aac.",
						},
						"audio_channel": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "Audio channel count, valid value: 2.",
						},
					},
				},
			},

			// computed
			"template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Just-in-time transcode template unique identifier.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template creation time.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template last modified time.",
			},

			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template type, valid values: preset, custom.",
			},
		},
	}
}

func resourceTencentCloudTeoJustInTimeTranscodeTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_just_in_time_transcode_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = teov20220901.NewCreateJustInTimeTranscodeTemplateRequest()
		response   = teov20220901.NewCreateJustInTimeTranscodeTemplateResponse()
		zoneId     string
		templateId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	if v, ok := d.GetOk("video_stream_switch"); ok {
		request.VideoStreamSwitch = helper.String(v.(string))
	}

	if v, ok := d.GetOk("audio_stream_switch"); ok {
		request.AudioStreamSwitch = helper.String(v.(string))
	}

	if v, ok := d.GetOk("video_template"); ok {
		for _, item := range v.([]interface{}) {
			videoTemplateMap := item.(map[string]interface{})
			videoTemplateInfo := teov20220901.VideoTemplateInfo{}
			if v, ok := videoTemplateMap["codec"].(string); ok && v != "" {
				videoTemplateInfo.Codec = helper.String(v)
			}
			if v, ok := videoTemplateMap["fps"].(float64); ok && v != 0 {
				videoTemplateInfo.Fps = helper.Float64(v)
			}
			if v, ok := videoTemplateMap["bitrate"].(int); ok && v != 0 {
				videoTemplateInfo.Bitrate = helper.IntUint64(v)
			}
			if v, ok := videoTemplateMap["resolution_adaptive"].(string); ok && v != "" {
				videoTemplateInfo.ResolutionAdaptive = helper.String(v)
			}
			if v, ok := videoTemplateMap["width"].(int); ok && v != 0 {
				videoTemplateInfo.Width = helper.IntUint64(v)
			}
			if v, ok := videoTemplateMap["height"].(int); ok && v != 0 {
				videoTemplateInfo.Height = helper.IntUint64(v)
			}
			if v, ok := videoTemplateMap["fill_type"].(string); ok && v != "" {
				videoTemplateInfo.FillType = helper.String(v)
			}
			request.VideoTemplate = &videoTemplateInfo
		}
	}

	if v, ok := d.GetOk("audio_template"); ok {
		for _, item := range v.([]interface{}) {
			audioTemplateMap := item.(map[string]interface{})
			audioTemplateInfo := teov20220901.AudioTemplateInfo{}
			if v, ok := audioTemplateMap["codec"].(string); ok && v != "" {
				audioTemplateInfo.Codec = helper.String(v)
			}
			if v, ok := audioTemplateMap["audio_channel"].(int); ok && v != 0 {
				audioTemplateInfo.AudioChannel = helper.IntUint64(v)
			}
			request.AudioTemplate = &audioTemplateInfo
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateJustInTimeTranscodeTemplateWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo just in time transcode template failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo just in time transcode template failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.TemplateId == nil {
		return fmt.Errorf("TemplateId is nil.")
	}

	templateId = *response.Response.TemplateId
	d.SetId(strings.Join([]string{zoneId, templateId}, tccommon.FILED_SP))
	return resourceTencentCloudTeoJustInTimeTranscodeTemplateRead(d, meta)
}

func resourceTencentCloudTeoJustInTimeTranscodeTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_just_in_time_transcode_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	templateId := idSplit[1]

	respData, err := service.DescribeTeoJustInTimeTranscodeTemplateById(ctx, zoneId, templateId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_teo_just_in_time_transcode_template` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	if respData.TemplateId != nil {
		_ = d.Set("template_id", respData.TemplateId)
	}

	if respData.TemplateName != nil {
		_ = d.Set("template_name", respData.TemplateName)
	}

	if respData.Comment != nil {
		_ = d.Set("comment", respData.Comment)
	}

	if respData.Type != nil {
		_ = d.Set("type", respData.Type)
	}

	if respData.VideoStreamSwitch != nil {
		_ = d.Set("video_stream_switch", respData.VideoStreamSwitch)
	}

	if respData.AudioStreamSwitch != nil {
		_ = d.Set("audio_stream_switch", respData.AudioStreamSwitch)
	}

	if respData.VideoTemplate != nil {
		videoTemplateList := make([]map[string]interface{}, 0, 1)
		videoTemplateMap := map[string]interface{}{}
		if respData.VideoTemplate.Codec != nil {
			videoTemplateMap["codec"] = respData.VideoTemplate.Codec
		}
		if respData.VideoTemplate.Fps != nil {
			videoTemplateMap["fps"] = respData.VideoTemplate.Fps
		}
		if respData.VideoTemplate.Bitrate != nil {
			videoTemplateMap["bitrate"] = respData.VideoTemplate.Bitrate
		}
		if respData.VideoTemplate.ResolutionAdaptive != nil {
			videoTemplateMap["resolution_adaptive"] = respData.VideoTemplate.ResolutionAdaptive
		}
		if respData.VideoTemplate.Width != nil {
			videoTemplateMap["width"] = respData.VideoTemplate.Width
		}
		if respData.VideoTemplate.Height != nil {
			videoTemplateMap["height"] = respData.VideoTemplate.Height
		}
		if respData.VideoTemplate.FillType != nil {
			videoTemplateMap["fill_type"] = respData.VideoTemplate.FillType
		}
		videoTemplateList = append(videoTemplateList, videoTemplateMap)
		_ = d.Set("video_template", videoTemplateList)
	}

	if respData.AudioTemplate != nil {
		audioTemplateList := make([]map[string]interface{}, 0, 1)
		audioTemplateMap := map[string]interface{}{}
		if respData.AudioTemplate.Codec != nil {
			audioTemplateMap["codec"] = respData.AudioTemplate.Codec
		}
		if respData.AudioTemplate.AudioChannel != nil {
			audioTemplateMap["audio_channel"] = respData.AudioTemplate.AudioChannel
		}
		audioTemplateList = append(audioTemplateList, audioTemplateMap)
		_ = d.Set("audio_template", audioTemplateList)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}

	return nil
}

func resourceTencentCloudTeoJustInTimeTranscodeTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_just_in_time_transcode_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewDeleteJustInTimeTranscodeTemplatesRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	templateId := idSplit[1]

	request.ZoneId = &zoneId
	request.TemplateIds = []*string{&templateId}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteJustInTimeTranscodeTemplatesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo just in time transcode template failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
