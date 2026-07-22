package ci

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCiMediaSpeechRecognitionTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaSpeechRecognitionTemplateCreate,
		Read:   resourceTencentCloudCiMediaSpeechRecognitionTemplateRead,
		Update: resourceTencentCloudCiMediaSpeechRecognitionTemplateUpdate,
		Delete: resourceTencentCloudCiMediaSpeechRecognitionTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "存储桶名称",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "模板名称 仅 支持 `Chinese`，`English`，`numbers`，`_`，`-` 和 `*`。",
			},

			"speech_recognition": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "音频 配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine_model_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Engine model 类型，divided into phone scene 和 non-phone scene，phone scene: 8k_zh: phone 8k Chinese Mandarin general (可以 是 用于dual-channel 音频)，8k_zh_s: phone 8k Chinese Mandarin speaker separation (仅 对于 monophonic 音频)，8k_en: Telephone 8k English; non-telephone scene: 16k_zh: 16k Mandarin Chinese，16k_zh_video: 16k 音频 和 视频 字段，16k_en: 16k English，16k_ca: 16k Cantonese，16k_ja: 16k Japanese，16k_zh_edu: Chinese education，16k_en_edu: English education，16k_zh_medical: medical，16k_th: Thai，16k_zh_dialect: multi-dialect，支持 23 dialects。",
						},
						"channel_num": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "数量 voice channels: 1 表示 mono. EngineModelType 支持 仅 mono 对于 non-telephone scenarios，和 2 表示 dual channels (仅 8k_zh 引擎 model 支持 dual channels，其中 should correspond 到 both sides 的 call)。",
						},
						"res_text_format": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Recognition 结果 返回 form: 0 表示 recognition 结果 text (包括 segmented 时间 stamps)，1 是 detailed recognition 结果 在 word 级别 granularity，without punctuation，和 includes speech 速率 值 ( 列表 word 时间 stamps，generally 用于generate subtitle scenes)，2 Detailed recognition results 在 word-级别 granularity (包括 punctuation 和 speech 速率 值).。",
						},
						"filter_dirty": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "是否filter dirty words (currently 支持 Mandarin Chinese 引擎): 0 表示 不 到 过滤器 dirty words，1 表示 到 过滤器 dirty words，2 表示 到 replace dirty words 使用 *， 默认值为 0。",
						},
						"filter_modal": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "是否pass modal particles (currently 支持 Mandarin Chinese 引擎): 0 表示 不 到 过滤器 modal particles，1 表示 partial filtering，2 表示 strict filtering，和 默认值为 0。",
						},
						"convert_num_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "是否perform intelligent conversion 的 Arabic numerals (currently 支持 Mandarin Chinese 引擎): 0 表示 无 conversion，directly output Chinese numbers，1 表示 intelligently convert 到 Arabic numerals according 到 scene，3 表示 启用 math-related digital conversion， 默认值为 0。",
						},
						"speaker_diarization": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "是否enable speaker separation: 0 表示 不 已启用，1 表示 已启用 (仅 支持 8k_zh，16k_zh，16k_zh_video，monophonic 音频)， 默认值为 0，注意: 8K telephony scenarios suggest 使用 dual-channel 到 distinguish between two parties，集合 ChannelNum=2 是 enough，无 need 到 启用 speaker separation。",
						},
						"speaker_number": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "数量 speakers 到 是 separated (need 到 是 使用 在 conjunction 使用 enabling speaker separation)，取值范围：0-10，0 表示 automatic separation (currently 仅 支持 <= 6 people)，1-10 表示 数量 指定 speakers 到 是 separated. 默认值为 0。",
						},
						"filter_punc": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "是否filter punctuation (currently 支持 Mandarin Chinese 引擎): 0 表示 无 filtering，1 表示 filtering end-的-sentence punctuation，2 表示 filtering all punctuation， 默认值为 0。",
						},
						"output_file_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "输出文件 类型，可选 txt，srt. 默认为 txt。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCiMediaSpeechRecognitionTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_speech_recognition_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request = cos.CreateMediaSpeechRecognitionTemplateOptions{
			Tag: "SpeechRecognition",
		}
		templateId string
		bucket     string
	)

	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	} else {
		return errors.New("get bucket failed!")
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "speech_recognition"); ok {
		speechRecognition := cos.SpeechRecognition{}
		if v, ok := dMap["engine_model_type"]; ok {
			speechRecognition.EngineModelType = v.(string)
		}
		if v, ok := dMap["channel_num"]; ok {
			speechRecognition.ChannelNum = v.(string)
		}
		if v, ok := dMap["res_text_format"]; ok {
			speechRecognition.ResTextFormat = v.(string)
		}
		if v, ok := dMap["filter_dirty"]; ok {
			speechRecognition.FilterDirty = v.(string)
		}
		if v, ok := dMap["filter_modal"]; ok {
			speechRecognition.FilterModal = v.(string)
		}
		if v, ok := dMap["convert_num_mode"]; ok {
			speechRecognition.ConvertNumMode = v.(string)
		}
		if v, ok := dMap["speaker_diarization"]; ok {
			speechRecognition.SpeakerDiarization = v.(string)
		}
		if v, ok := dMap["speaker_number"]; ok {
			speechRecognition.SpeakerNumber = v.(string)
		}
		if v, ok := dMap["filter_punc"]; ok {
			speechRecognition.FilterPunc = v.(string)
		}
		if v, ok := dMap["output_file_type"]; ok {
			speechRecognition.OutputFileType = v.(string)
		}
		request.SpeechRecognition = &speechRecognition
	}

	var response *cos.CreateMediaTemplateResult
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket).CI.CreateMediaSpeechRecognitionTemplate(ctx, &request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "CreateMediaSpeechRecognitionTemplate", request, result)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaSpeechRecognitionTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = response.Template.TemplateId
	d.SetId(bucket + tccommon.FILED_SP + templateId)

	return resourceTencentCloudCiMediaSpeechRecognitionTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaSpeechRecognitionTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_speech_recognition_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	mediaSpeechRecognitionTemplate, err := service.DescribeCiMediaTemplateById(ctx, bucket, templateId)
	if err != nil {
		return err
	}

	if mediaSpeechRecognitionTemplate == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	_ = d.Set("bucket", bucket)

	if mediaSpeechRecognitionTemplate.Name != "" {
		_ = d.Set("name", mediaSpeechRecognitionTemplate.Name)
	}

	if mediaSpeechRecognitionTemplate.SpeechRecognition != nil {
		speechRecognitionMap := map[string]interface{}{}

		if mediaSpeechRecognitionTemplate.SpeechRecognition.EngineModelType != "" {
			speechRecognitionMap["engine_model_type"] = mediaSpeechRecognitionTemplate.SpeechRecognition.EngineModelType
		}

		if mediaSpeechRecognitionTemplate.SpeechRecognition.ChannelNum != "" {
			speechRecognitionMap["channel_num"] = mediaSpeechRecognitionTemplate.SpeechRecognition.ChannelNum
		}

		if mediaSpeechRecognitionTemplate.SpeechRecognition.ResTextFormat != "" {
			speechRecognitionMap["res_text_format"] = mediaSpeechRecognitionTemplate.SpeechRecognition.ResTextFormat
		}

		if mediaSpeechRecognitionTemplate.SpeechRecognition.FilterDirty != "" {
			speechRecognitionMap["filter_dirty"] = mediaSpeechRecognitionTemplate.SpeechRecognition.FilterDirty
		}

		if mediaSpeechRecognitionTemplate.SpeechRecognition.FilterModal != "" {
			speechRecognitionMap["filter_modal"] = mediaSpeechRecognitionTemplate.SpeechRecognition.FilterModal
		}

		if mediaSpeechRecognitionTemplate.SpeechRecognition.ConvertNumMode != "" {
			speechRecognitionMap["convert_num_mode"] = mediaSpeechRecognitionTemplate.SpeechRecognition.ConvertNumMode
		}

		if mediaSpeechRecognitionTemplate.SpeechRecognition.SpeakerDiarization != "" {
			speechRecognitionMap["speaker_diarization"] = mediaSpeechRecognitionTemplate.SpeechRecognition.SpeakerDiarization
		}

		if mediaSpeechRecognitionTemplate.SpeechRecognition.SpeakerNumber != "" {
			speechRecognitionMap["speaker_number"] = mediaSpeechRecognitionTemplate.SpeechRecognition.SpeakerNumber
		}

		if mediaSpeechRecognitionTemplate.SpeechRecognition.FilterPunc != "" {
			speechRecognitionMap["filter_punc"] = mediaSpeechRecognitionTemplate.SpeechRecognition.FilterPunc
		}

		if mediaSpeechRecognitionTemplate.SpeechRecognition.OutputFileType != "" {
			speechRecognitionMap["output_file_type"] = mediaSpeechRecognitionTemplate.SpeechRecognition.OutputFileType
		}

		_ = d.Set("speech_recognition", []interface{}{speechRecognitionMap})
	}

	return nil
}

func resourceTencentCloudCiMediaSpeechRecognitionTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_speech_recognition_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := cos.CreateMediaSpeechRecognitionTemplateOptions{
		Tag: "SpeechRecognition",
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if d.HasChange("speech_recognition") {
		if dMap, ok := helper.InterfacesHeadMap(d, "speech_recognition"); ok {
			speechRecognition := cos.SpeechRecognition{}
			if v, ok := dMap["engine_model_type"]; ok {
				speechRecognition.EngineModelType = v.(string)
			}
			if v, ok := dMap["channel_num"]; ok {
				speechRecognition.ChannelNum = v.(string)
			}
			if v, ok := dMap["res_text_format"]; ok {
				speechRecognition.ResTextFormat = v.(string)
			}
			if v, ok := dMap["filter_dirty"]; ok {
				speechRecognition.FilterDirty = v.(string)
			}
			if v, ok := dMap["filter_modal"]; ok {
				speechRecognition.FilterModal = v.(string)
			}
			if v, ok := dMap["convert_num_mode"]; ok {
				speechRecognition.ConvertNumMode = v.(string)
			}
			if v, ok := dMap["speaker_diarization"]; ok {
				speechRecognition.SpeakerDiarization = v.(string)
			}
			if v, ok := dMap["speaker_number"]; ok {
				speechRecognition.SpeakerNumber = v.(string)
			}
			if v, ok := dMap["filter_punc"]; ok {
				speechRecognition.FilterPunc = v.(string)
			}
			if v, ok := dMap["output_file_type"]; ok {
				speechRecognition.OutputFileType = v.(string)
			}
			request.SpeechRecognition = &speechRecognition
		}
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket).CI.UpdateMediaSpeechRecognitionTemplate(ctx, &request, templateId)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaSpeechRecognitionTemplate", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaSpeechRecognitionTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaSpeechRecognitionTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaSpeechRecognitionTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_speech_recognition_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	if err := service.DeleteCiMediaTemplateById(ctx, bucket, templateId); err != nil {
		return err
	}

	return nil
}
