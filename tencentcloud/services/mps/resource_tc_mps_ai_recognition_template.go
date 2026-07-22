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

func ResourceTencentCloudMpsAiRecognitionTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsAiRecognitionTemplateCreate,
		Read:   resourceTencentCloudMpsAiRecognitionTemplateRead,
		Update: resourceTencentCloudMpsAiRecognitionTemplateUpdate,
		Delete: resourceTencentCloudMpsAiRecognitionTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Ai recognition 模板名称，长度 限制: 64 字符。",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Ai recognition 模板描述 信息，长度 限制: 256 字符。",
			},

			"face_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Face recognition control 参数。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Ai face recognition 任务 switch，可选 值:ON/OFF。",
						},
						"score": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: "Face recognition 过滤器 score，当 recognition 结果 reaches score above， recognition 结果 将 是 返回. 默认为 95 points. 取值范围：0 - 100。",
						},
						"default_library_label_set": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Default face 过滤器 标签，指定tag 的 默认值 face 该 needs 到 是 返回. 如果未填写 或 空，all 默认值 face results 将 是 返回. 标签 可选 值:entertainment，sport，politician。",
						},
						"user_define_library_label_set": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "用户-defined face 过滤器 标签，指定tag 的 用户-defined face 该 needs 到 是 返回. 如果未填写 或 空，all 自定义 face results 将 是 返回. 最大tags 是 100，和 长度 的 each 标签 是 up 到 16 字符。",
						},
						"face_library": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Face 库 selection，可选 值:Default，UserDefine，All默认值：All，使用 系统 默认值 face 库 和 用户-defined face 库。",
						},
					},
				},
			},

			"ocr_full_text_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Ocr full text control 参数。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Ocr full text recognition 任务 switch，可选 值:ON/OFF。",
						},
					},
				},
			},

			"ocr_words_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Ocr words recognition control 参数。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Ocr words recognition 任务 switch，可选 值:ON/OFF。",
						},
						"label_set": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Keyword 过滤器 标签，指定label 的 keyword 到 是 返回. 如果未填写 或 空，all results 将 是 返回. 最大tags 是 10，和 长度 的 each 标签 是 up 到 16 字符。",
						},
					},
				},
			},

			"asr_full_text_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Asr full text recognition control 参数。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Asr full text recognition 任务 switch，可选 值:ON/OFF。",
						},
						"subtitle_format": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Generated subtitle 文件 格式，如果 left blank 或 blank 字符串 表示 无 subtitle 文件 将 是 generated，可选 值:vtt: Generate WebVTT subtitle files。",
						},
					},
				},
			},

			"asr_words_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Asr word recognition control 参数。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Asr word recognition 任务 switch，可选 值:ON/OFF。",
						},
						"label_set": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Keyword 过滤器 标签，指定label 的 keyword 到 是 返回. 如果未填写 或 空，all results 将 是 返回. 最大tags 是 10，和 长度 的 each 标签 是 up 到 16 字符。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMpsAiRecognitionTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_ai_recognition_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = mps.NewCreateAIRecognitionTemplateRequest()
		response   = mps.NewCreateAIRecognitionTemplateResponse()
		definition int64
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "face_configure"); ok {
		faceConfigureInfo := mps.FaceConfigureInfo{}
		if v, ok := dMap["switch"]; ok {
			faceConfigureInfo.Switch = helper.String(v.(string))
		}
		if v, ok := dMap["score"]; ok {
			faceConfigureInfo.Score = helper.Float64(v.(float64))
		}
		if v, ok := dMap["default_library_label_set"]; ok {
			defaultLibraryLabelSetSet := v.(*schema.Set).List()
			for i := range defaultLibraryLabelSetSet {
				defaultLibraryLabelSet := defaultLibraryLabelSetSet[i].(string)
				faceConfigureInfo.DefaultLibraryLabelSet = append(faceConfigureInfo.DefaultLibraryLabelSet, &defaultLibraryLabelSet)
			}
		}
		if v, ok := dMap["user_define_library_label_set"]; ok {
			userDefineLibraryLabelSetSet := v.(*schema.Set).List()
			for i := range userDefineLibraryLabelSetSet {
				userDefineLibraryLabelSet := userDefineLibraryLabelSetSet[i].(string)
				faceConfigureInfo.UserDefineLibraryLabelSet = append(faceConfigureInfo.UserDefineLibraryLabelSet, &userDefineLibraryLabelSet)
			}
		}
		if v, ok := dMap["face_library"]; ok {
			faceConfigureInfo.FaceLibrary = helper.String(v.(string))
		}
		request.FaceConfigure = &faceConfigureInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "ocr_full_text_configure"); ok {
		ocrFullTextConfigureInfo := mps.OcrFullTextConfigureInfo{}
		if v, ok := dMap["switch"]; ok {
			ocrFullTextConfigureInfo.Switch = helper.String(v.(string))
		}
		request.OcrFullTextConfigure = &ocrFullTextConfigureInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "ocr_words_configure"); ok {
		ocrWordsConfigureInfo := mps.OcrWordsConfigureInfo{}
		if v, ok := dMap["switch"]; ok {
			ocrWordsConfigureInfo.Switch = helper.String(v.(string))
		}
		if v, ok := dMap["label_set"]; ok {
			labelSetSet := v.(*schema.Set).List()
			for i := range labelSetSet {
				labelSet := labelSetSet[i].(string)
				ocrWordsConfigureInfo.LabelSet = append(ocrWordsConfigureInfo.LabelSet, &labelSet)
			}
		}
		request.OcrWordsConfigure = &ocrWordsConfigureInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "asr_full_text_configure"); ok {
		asrFullTextConfigureInfo := mps.AsrFullTextConfigureInfo{}
		if v, ok := dMap["switch"]; ok {
			asrFullTextConfigureInfo.Switch = helper.String(v.(string))
		}
		if v := dMap["subtitle_format"]; v != "" {
			asrFullTextConfigureInfo.SubtitleFormat = helper.String(v.(string))
		}
		request.AsrFullTextConfigure = &asrFullTextConfigureInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "asr_words_configure"); ok {
		asrWordsConfigureInfo := mps.AsrWordsConfigureInfo{}
		if v, ok := dMap["switch"]; ok {
			asrWordsConfigureInfo.Switch = helper.String(v.(string))
		}
		if v, ok := dMap["label_set"]; ok {
			labelSetSet := v.(*schema.Set).List()
			for i := range labelSetSet {
				labelSet := labelSetSet[i].(string)
				asrWordsConfigureInfo.LabelSet = append(asrWordsConfigureInfo.LabelSet, &labelSet)
			}
		}
		request.AsrWordsConfigure = &asrWordsConfigureInfo
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().CreateAIRecognitionTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps aiRecognitionTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition = *response.Response.Definition
	d.SetId(helper.Int64ToStr(definition))

	return resourceTencentCloudMpsAiRecognitionTemplateRead(d, meta)
}

func resourceTencentCloudMpsAiRecognitionTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_ai_recognition_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	definition := d.Id()

	aiRecognitionTemplate, err := service.DescribeMpsAiRecognitionTemplateById(ctx, definition)
	if err != nil {
		return err
	}

	if aiRecognitionTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsAiRecognitionTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if aiRecognitionTemplate.Name != nil {
		_ = d.Set("name", aiRecognitionTemplate.Name)
	}

	if aiRecognitionTemplate.Comment != nil {
		_ = d.Set("comment", aiRecognitionTemplate.Comment)
	}

	if aiRecognitionTemplate.FaceConfigure != nil {
		faceConfigureMap := map[string]interface{}{}

		if aiRecognitionTemplate.FaceConfigure.Switch != nil {
			faceConfigureMap["switch"] = aiRecognitionTemplate.FaceConfigure.Switch
		}

		if aiRecognitionTemplate.FaceConfigure.Score != nil {
			faceConfigureMap["score"] = aiRecognitionTemplate.FaceConfigure.Score
		}

		if aiRecognitionTemplate.FaceConfigure.DefaultLibraryLabelSet != nil {
			faceConfigureMap["default_library_label_set"] = aiRecognitionTemplate.FaceConfigure.DefaultLibraryLabelSet
		}

		if aiRecognitionTemplate.FaceConfigure.UserDefineLibraryLabelSet != nil {
			faceConfigureMap["user_define_library_label_set"] = aiRecognitionTemplate.FaceConfigure.UserDefineLibraryLabelSet
		}

		if aiRecognitionTemplate.FaceConfigure.FaceLibrary != nil {
			faceConfigureMap["face_library"] = aiRecognitionTemplate.FaceConfigure.FaceLibrary
		}

		_ = d.Set("face_configure", []interface{}{faceConfigureMap})
	}

	if aiRecognitionTemplate.OcrFullTextConfigure != nil {
		ocrFullTextConfigureMap := map[string]interface{}{}

		if aiRecognitionTemplate.OcrFullTextConfigure.Switch != nil {
			ocrFullTextConfigureMap["switch"] = aiRecognitionTemplate.OcrFullTextConfigure.Switch
		}

		_ = d.Set("ocr_full_text_configure", []interface{}{ocrFullTextConfigureMap})
	}

	if aiRecognitionTemplate.OcrWordsConfigure != nil {
		ocrWordsConfigureMap := map[string]interface{}{}

		if aiRecognitionTemplate.OcrWordsConfigure.Switch != nil {
			ocrWordsConfigureMap["switch"] = aiRecognitionTemplate.OcrWordsConfigure.Switch
		}

		if aiRecognitionTemplate.OcrWordsConfigure.LabelSet != nil {
			ocrWordsConfigureMap["label_set"] = aiRecognitionTemplate.OcrWordsConfigure.LabelSet
		}

		_ = d.Set("ocr_words_configure", []interface{}{ocrWordsConfigureMap})
	}

	if aiRecognitionTemplate.AsrFullTextConfigure != nil {
		asrFullTextConfigureMap := map[string]interface{}{}

		if aiRecognitionTemplate.AsrFullTextConfigure.Switch != nil {
			asrFullTextConfigureMap["switch"] = aiRecognitionTemplate.AsrFullTextConfigure.Switch
		}

		if aiRecognitionTemplate.AsrFullTextConfigure.SubtitleFormat != nil {
			asrFullTextConfigureMap["subtitle_format"] = aiRecognitionTemplate.AsrFullTextConfigure.SubtitleFormat
		}

		_ = d.Set("asr_full_text_configure", []interface{}{asrFullTextConfigureMap})
	}

	if aiRecognitionTemplate.AsrWordsConfigure != nil {
		asrWordsConfigureMap := map[string]interface{}{}

		if aiRecognitionTemplate.AsrWordsConfigure.Switch != nil {
			asrWordsConfigureMap["switch"] = aiRecognitionTemplate.AsrWordsConfigure.Switch
		}

		if aiRecognitionTemplate.AsrWordsConfigure.LabelSet != nil {
			asrWordsConfigureMap["label_set"] = aiRecognitionTemplate.AsrWordsConfigure.LabelSet
		}

		_ = d.Set("asr_words_configure", []interface{}{asrWordsConfigureMap})
	}

	return nil
}

func resourceTencentCloudMpsAiRecognitionTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_ai_recognition_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := mps.NewModifyAIRecognitionTemplateRequest()

	definition := d.Id()

	request.Definition = helper.StrToInt64Point(definition)

	mutableArgs := []string{"name", "comment", "face_configure", "ocr_full_text_configure", "ocr_words_configure", "asr_full_text_configure", "asr_words_configure"}

	needChange := false

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

		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "face_configure"); ok {
			faceConfigureInfo := mps.FaceConfigureInfoForUpdate{}
			if v, ok := dMap["switch"]; ok {
				faceConfigureInfo.Switch = helper.String(v.(string))
			}
			if v, ok := dMap["score"]; ok {
				faceConfigureInfo.Score = helper.Float64(v.(float64))
			}
			if v, ok := dMap["default_library_label_set"]; ok {
				defaultLibraryLabelSetSet := v.(*schema.Set).List()
				for i := range defaultLibraryLabelSetSet {
					defaultLibraryLabelSet := defaultLibraryLabelSetSet[i].(string)
					faceConfigureInfo.DefaultLibraryLabelSet = append(faceConfigureInfo.DefaultLibraryLabelSet, &defaultLibraryLabelSet)
				}
			}
			if v, ok := dMap["user_define_library_label_set"]; ok {
				userDefineLibraryLabelSetSet := v.(*schema.Set).List()
				for i := range userDefineLibraryLabelSetSet {
					userDefineLibraryLabelSet := userDefineLibraryLabelSetSet[i].(string)
					faceConfigureInfo.UserDefineLibraryLabelSet = append(faceConfigureInfo.UserDefineLibraryLabelSet, &userDefineLibraryLabelSet)
				}
			}
			if v, ok := dMap["face_library"]; ok {
				faceConfigureInfo.FaceLibrary = helper.String(v.(string))
			}
			request.FaceConfigure = &faceConfigureInfo
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "ocr_full_text_configure"); ok {
			ocrFullTextConfigureInfo := mps.OcrFullTextConfigureInfoForUpdate{}
			if v, ok := dMap["switch"]; ok {
				ocrFullTextConfigureInfo.Switch = helper.String(v.(string))
			}
			request.OcrFullTextConfigure = &ocrFullTextConfigureInfo
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "ocr_words_configure"); ok {
			ocrWordsConfigureInfo := mps.OcrWordsConfigureInfoForUpdate{}
			if v, ok := dMap["switch"]; ok {
				ocrWordsConfigureInfo.Switch = helper.String(v.(string))
			}
			if v, ok := dMap["label_set"]; ok {
				labelSetSet := v.(*schema.Set).List()
				for i := range labelSetSet {
					labelSet := labelSetSet[i].(string)
					ocrWordsConfigureInfo.LabelSet = append(ocrWordsConfigureInfo.LabelSet, &labelSet)
				}
			}
			request.OcrWordsConfigure = &ocrWordsConfigureInfo
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "asr_full_text_configure"); ok {
			asrFullTextConfigureInfo := mps.AsrFullTextConfigureInfoForUpdate{}
			if v, ok := dMap["switch"]; ok {
				asrFullTextConfigureInfo.Switch = helper.String(v.(string))
			}
			if v := dMap["subtitle_format"]; v != "" {
				asrFullTextConfigureInfo.SubtitleFormat = helper.String(v.(string))
			}
			request.AsrFullTextConfigure = &asrFullTextConfigureInfo
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "asr_words_configure"); ok {
			asrWordsConfigureInfo := mps.AsrWordsConfigureInfoForUpdate{}
			if v, ok := dMap["switch"]; ok {
				asrWordsConfigureInfo.Switch = helper.String(v.(string))
			}
			if v, ok := dMap["label_set"]; ok {
				labelSetSet := v.(*schema.Set).List()
				for i := range labelSetSet {
					labelSet := labelSetSet[i].(string)
					asrWordsConfigureInfo.LabelSet = append(asrWordsConfigureInfo.LabelSet, &labelSet)
				}
			}
			request.AsrWordsConfigure = &asrWordsConfigureInfo
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().ModifyAIRecognitionTemplate(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mps aiRecognitionTemplate failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudMpsAiRecognitionTemplateRead(d, meta)
}

func resourceTencentCloudMpsAiRecognitionTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_ai_recognition_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	definition := d.Id()

	if err := service.DeleteMpsAiRecognitionTemplateById(ctx, definition); err != nil {
		return err
	}

	return nil
}
