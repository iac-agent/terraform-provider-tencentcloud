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

func ResourceTencentCloudMpsContentReviewTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsContentReviewTemplateCreate,
		Read:   resourceTencentCloudMpsContentReviewTemplateRead,
		Update: resourceTencentCloudMpsContentReviewTemplateUpdate,
		Delete: resourceTencentCloudMpsContentReviewTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "内容 review 模板名称，长度 限制: 64 字符。",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "内容 review 模板描述 信息，长度 限制: 256 字符。",
			},

			"porn_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Control 参数 对于 porn 镜像。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"img_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Porn 镜像 Identification Control Parameters。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Porn screen 任务 switch，可选 值:ON/OFF。",
									},
									"label_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Porn 镜像 过滤器 标签，如果 review 结果 包含selected 标签， 结果 将 是 返回. 如果 过滤器 标签 是 空，all review results 将 是 返回. 可选 值 是:porn，vulgar，intimacy，sexy。",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging suspected violations. 当 smart review reaches score above，它 是 considered suspected violations. 如果 它 是 不 filled， 默认为 90 points. 取值范围：0~100。",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging whether manual review 为必填项 对于 violations. 当 intelligent review reaches score above，它 是 considered 该 manual review 为必填项. 如果 它 是 不 filled， 默认为 0. 取值范围：0~100。",
									},
								},
							},
						},
						"asr_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Voice pornography control 参数。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Voice pornography 任务 switch，可选 值:ON/OFF。",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging suspected violations. 当 smart review reaches score above，它 是 considered suspected violations. 如果 它 是 不 filled， 默认为 100 points. 取值范围：0~100。",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging whether manual review 为必填项 对于 violations. 当 intelligent review reaches score above，它 是 considered 该 manual review 为必填项. 如果 它 是 不 filled， 默认为 75 points. 取值范围：0~100。",
									},
								},
							},
						},
						"ocr_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Ocr pornography control 参数。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Ocr pornography 任务 switch，可选 值:ON/OFF。",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging suspected violations. 当 smart review reaches score above，它 是 considered suspected violations. 如果 它 是 不 filled， 默认为 100 points. 取值范围：0~100。",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging whether manual review 为必填项 对于 violations. 当 intelligent review reaches score above，它 是 considered 该 manual review 为必填项. 如果 它 是 不 filled， 默认为 75 points. 取值范围：0~100。",
									},
								},
							},
						},
					},
				},
			},

			"terrorism_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Control 参数 对于 unsafe 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"img_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Terrorism 镜像 任务 control 参数。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Terrorism 镜像 任务 switch，可选 值:ON/OFF。",
									},
									"label_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Terrorism 镜像 过滤器 标签，如果 review 结果 包含selected 标签， 结果 将 是 返回，如果 过滤器 标签 是 空，all review results 将 是 返回， 可选 值 是:guns，crowd，bloody，police，banners，militant，explosion，terrorists，scenario。",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging suspected violations. 当 smart review reaches score above，它 是 considered suspected violations. 如果 它 是 不 filled， 默认为 90 points. 取值范围：0~100。",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging whether manual review 为必填项 对于 violations. 当 intelligent review reaches score above，它 是 considered 该 manual review 为必填项. 如果 它 是 不 filled， 默认为 80 points. 取值范围：0~100。",
									},
								},
							},
						},
						"ocr_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Ocr terrorism 任务 Control Parameters。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Ocr terrorism 镜像 任务 switch，可选 值:ON/OFF。",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging suspected violations. 当 smart review reaches score above，它 是 considered suspected violations. 如果 它 是 不 filled， 默认为 100 points. 取值范围：0~100。",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging whether manual review 为必填项 对于 violations. 当 intelligent review reaches score above，它 是 considered 该 manual review 为必填项. 如果 它 是 不 filled， 默认为 75 points. 取值范围：0~100。",
									},
								},
							},
						},
					},
				},
			},

			"political_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Political control 参数。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"img_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Political 镜像 control 参数。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Political 镜像 任务 switch，可选 值:ON/OFF。",
									},
									"label_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Political 镜像 过滤器 标签，如果 review 结果 包含selected 标签， 结果 将 是 返回，如果 过滤器 标签 是 空，all review results 将 是 返回， 可选 值 是:violation_photo，politician，entertainment，sport，entrepreneur，scholar，celebrity，military。",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging suspected violations. 当 smart review reaches score above，它 是 considered suspected violations. 如果 它 是 不 filled， 默认为 97 points. 取值范围：0~100。",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging whether manual review 为必填项 对于 violations. 当 intelligent review reaches score above，它 是 considered 该 manual review 为必填项. 如果 它 是 不 filled， 默认为 95 points. 取值范围：0~100。",
									},
								},
							},
						},
						"asr_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Political asr control 参数。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Political asr 任务 switch，可选 值:ON/OFF。",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging suspected violations. 当 smart review reaches score above，它 是 considered suspected violations. 如果 它 是 不 filled， 默认为 100 points. 取值范围：0~100。",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging whether manual review 为必填项 对于 violations. 当 intelligent review reaches score above，它 是 considered 该 manual review 为必填项. 如果 它 是 不 filled， 默认为 75 points. 取值范围：0~100。",
									},
								},
							},
						},
						"ocr_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Political ocr control 参数。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Political ocr 任务 switch，可选 值:ON/OFF。",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging suspected violations. 当 smart review reaches score above，它 是 considered suspected violations. 如果 它 是 不 filled， 默认为 100 points. 取值范围：0~100。",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging whether manual review 为必填项 对于 violations. 当 intelligent review reaches score above，它 是 considered 该 manual review 为必填项. 如果 它 是 不 filled， 默认为 75 points. 取值范围：0~100。",
									},
								},
							},
						},
					},
				},
			},

			"prohibited_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Prohibited control 参数. Prohibited 内容 includes:abuse，drug-related violations.注意: 此 参数 是 不 yet 支持。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asr_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Voice Prohibition Control Parameters.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Voice Prohibition 任务 switch，可选 值:ON/OFF。",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging suspected violations. 当 smart review reaches score above，它 是 considered suspected violations. 如果 它 是 不 filled， 默认为 100 points. 取值范围：0~100。",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging whether manual review 为必填项 对于 violations. 当 intelligent review reaches score above，它 是 considered 该 manual review 为必填项. 如果 它 是 不 filled， 默认为 75 points. 取值范围：0~100。",
									},
								},
							},
						},
						"ocr_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Ocr Prohibition Control Parameters.注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Ocr Prohibition 任务 switch，可选 值:ON/OFF。",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging suspected violations. 当 smart review reaches score above，它 是 considered suspected violations. 如果 它 是 不 filled， 默认为 100 points. 取值范围：0~100。",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging whether manual review 为必填项 对于 violations. 当 intelligent review reaches score above，它 是 considered 该 manual review 为必填项. 如果 它 是 不 filled， 默认为 75 points. 取值范围：0~100。",
									},
								},
							},
						},
					},
				},
			},

			"user_define_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "用户-Defined 内容 Moderation Control Parameters。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"face_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "用户-defined face review control 参数。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "用户-defined face review 任务 switch，可选 值:ON/OFF。",
									},
									"label_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "用户-defined face review 标签， review 结果 包含selected 标签 和 返回result，如果 过滤器 标签 是 空，all review results 是 返回. 如果 您 want 到 使用 标签 filtering 函数，当 adding face 库，您 need 到 add corresponding character 标签The 最大tags 是 10，和 长度 的 each 标签 是 up 到 16 字符。",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging suspected violations. 当 smart review reaches score above，它 是 considered suspected violations. 如果 它 是 不 filled， 默认为 100 points. 取值范围：0~100。",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging whether manual review 为必填项 对于 violations. 当 intelligent review reaches score above，它 是 considered 该 manual review 为必填项. 如果 它 是 不 filled， 默认为 75 points. 取值范围：0~100。",
									},
								},
							},
						},
						"asr_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "用户-defined asr text review control 参数。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "用户-defined asr review 任务 switch，可选 值:ON/OFF。",
									},
									"label_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "用户-defined asr 标签， review 结果 包含selected 标签 和 返回result，如果 过滤器 标签 是 空，all review results 是 返回. 如果 您 want 到 使用 标签 filtering 函数，当 adding asr 库，您 need 到 add corresponding character 标签The 最大tags 是 10，和 长度 的 each 标签 是 up 到 16 字符。",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging suspected violations. 当 smart review reaches score above，它 是 considered suspected violations. 如果 它 是 不 filled， 默认为 100 points. 取值范围：0~100。",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging whether manual review 为必填项 对于 violations. 当 intelligent review reaches score above，它 是 considered 该 manual review 为必填项. 如果 它 是 不 filled， 默认为 75 points. 取值范围：0~100。",
									},
								},
							},
						},
						"ocr_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "用户-defined ocr text review control 参数。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "用户-defined ocr text review 任务 switch，可选 值:ON/OFF。",
									},
									"label_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "用户-defined ocr 标签， review 结果 包含selected 标签 和 返回result，如果 过滤器 标签 是 空，all review results 是 返回. 如果 您 want 到 使用 标签 filtering 函数，当 adding ocr 库，您 need 到 add corresponding character 标签The 最大tags 是 10，和 长度 的 each 标签 是 up 到 16 字符。",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging suspected violations. 当 smart review reaches score above，它 是 considered suspected violations. 如果 它 是 不 filled， 默认为 100 points. 取值范围：0~100。",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "score 阈值 对于 judging whether manual review 为必填项 对于 violations. 当 intelligent review reaches score above，它 是 considered 该 manual review 为必填项. 如果 它 是 不 filled， 默认为 75 points. 取值范围：0~100。",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMpsContentReviewTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_content_review_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = mps.NewCreateContentReviewTemplateRequest()
		response   = mps.NewCreateContentReviewTemplateResponse()
		definition string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "porn_configure"); ok {
		pornConfigureInfo := mps.PornConfigureInfo{}
		if imgReviewInfoMap, ok := helper.InterfaceToMap(dMap, "img_review_info"); ok {
			pornImgReviewTemplateInfo := mps.PornImgReviewTemplateInfo{}
			if v, ok := imgReviewInfoMap["switch"]; ok {
				pornImgReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := imgReviewInfoMap["label_set"]; ok {
				labelSetSet := v.(*schema.Set).List()
				for i := range labelSetSet {
					if labelSetSet[i] != nil {
						labelSet := labelSetSet[i].(string)
						pornImgReviewTemplateInfo.LabelSet = append(pornImgReviewTemplateInfo.LabelSet, &labelSet)
					}
				}
			}
			if v, ok := imgReviewInfoMap["block_confidence"]; ok {
				pornImgReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := imgReviewInfoMap["review_confidence"]; ok {
				pornImgReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			pornConfigureInfo.ImgReviewInfo = &pornImgReviewTemplateInfo
		}
		if asrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "asr_review_info"); ok {
			pornAsrReviewTemplateInfo := mps.PornAsrReviewTemplateInfo{}
			if v, ok := asrReviewInfoMap["switch"]; ok {
				pornAsrReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := asrReviewInfoMap["block_confidence"]; ok {
				pornAsrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := asrReviewInfoMap["review_confidence"]; ok {
				pornAsrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			pornConfigureInfo.AsrReviewInfo = &pornAsrReviewTemplateInfo
		}
		if ocrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "ocr_review_info"); ok {
			pornOcrReviewTemplateInfo := mps.PornOcrReviewTemplateInfo{}
			if v, ok := ocrReviewInfoMap["switch"]; ok {
				pornOcrReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := ocrReviewInfoMap["block_confidence"]; ok {
				pornOcrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := ocrReviewInfoMap["review_confidence"]; ok {
				pornOcrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			pornConfigureInfo.OcrReviewInfo = &pornOcrReviewTemplateInfo
		}
		request.PornConfigure = &pornConfigureInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "terrorism_configure"); ok {
		terrorismConfigureInfo := mps.TerrorismConfigureInfo{}
		if imgReviewInfoMap, ok := helper.InterfaceToMap(dMap, "img_review_info"); ok {
			terrorismImgReviewTemplateInfo := mps.TerrorismImgReviewTemplateInfo{}
			if v, ok := imgReviewInfoMap["switch"]; ok {
				terrorismImgReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := imgReviewInfoMap["label_set"]; ok {
				labelSetSet := v.(*schema.Set).List()
				for i := range labelSetSet {
					if labelSetSet[i] != nil {
						labelSet := labelSetSet[i].(string)
						terrorismImgReviewTemplateInfo.LabelSet = append(terrorismImgReviewTemplateInfo.LabelSet, &labelSet)
					}
				}
			}
			if v, ok := imgReviewInfoMap["block_confidence"]; ok {
				terrorismImgReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := imgReviewInfoMap["review_confidence"]; ok {
				terrorismImgReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			terrorismConfigureInfo.ImgReviewInfo = &terrorismImgReviewTemplateInfo
		}
		if ocrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "ocr_review_info"); ok {
			terrorismOcrReviewTemplateInfo := mps.TerrorismOcrReviewTemplateInfo{}
			if v, ok := ocrReviewInfoMap["switch"]; ok {
				terrorismOcrReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := ocrReviewInfoMap["block_confidence"]; ok {
				terrorismOcrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := ocrReviewInfoMap["review_confidence"]; ok {
				terrorismOcrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			terrorismConfigureInfo.OcrReviewInfo = &terrorismOcrReviewTemplateInfo
		}
		request.TerrorismConfigure = &terrorismConfigureInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "political_configure"); ok {
		politicalConfigureInfo := mps.PoliticalConfigureInfo{}
		if imgReviewInfoMap, ok := helper.InterfaceToMap(dMap, "img_review_info"); ok {
			politicalImgReviewTemplateInfo := mps.PoliticalImgReviewTemplateInfo{}
			if v, ok := imgReviewInfoMap["switch"]; ok {
				politicalImgReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := imgReviewInfoMap["label_set"]; ok {
				labelSetSet := v.(*schema.Set).List()
				for i := range labelSetSet {
					if labelSetSet[i] != nil {
						labelSet := labelSetSet[i].(string)
						politicalImgReviewTemplateInfo.LabelSet = append(politicalImgReviewTemplateInfo.LabelSet, &labelSet)
					}
				}
			}
			if v, ok := imgReviewInfoMap["block_confidence"]; ok {
				politicalImgReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := imgReviewInfoMap["review_confidence"]; ok {
				politicalImgReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			politicalConfigureInfo.ImgReviewInfo = &politicalImgReviewTemplateInfo
		}
		if asrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "asr_review_info"); ok {
			politicalAsrReviewTemplateInfo := mps.PoliticalAsrReviewTemplateInfo{}
			if v, ok := asrReviewInfoMap["switch"]; ok {
				politicalAsrReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := asrReviewInfoMap["block_confidence"]; ok {
				politicalAsrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := asrReviewInfoMap["review_confidence"]; ok {
				politicalAsrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			politicalConfigureInfo.AsrReviewInfo = &politicalAsrReviewTemplateInfo
		}
		if ocrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "ocr_review_info"); ok {
			politicalOcrReviewTemplateInfo := mps.PoliticalOcrReviewTemplateInfo{}
			if v, ok := ocrReviewInfoMap["switch"]; ok {
				politicalOcrReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := ocrReviewInfoMap["block_confidence"]; ok {
				politicalOcrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := ocrReviewInfoMap["review_confidence"]; ok {
				politicalOcrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			politicalConfigureInfo.OcrReviewInfo = &politicalOcrReviewTemplateInfo
		}
		request.PoliticalConfigure = &politicalConfigureInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "prohibited_configure"); ok {
		prohibitedConfigureInfo := mps.ProhibitedConfigureInfo{}
		if asrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "asr_review_info"); ok {
			prohibitedAsrReviewTemplateInfo := mps.ProhibitedAsrReviewTemplateInfo{}
			if v, ok := asrReviewInfoMap["switch"]; ok {
				prohibitedAsrReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := asrReviewInfoMap["block_confidence"]; ok {
				prohibitedAsrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := asrReviewInfoMap["review_confidence"]; ok {
				prohibitedAsrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			prohibitedConfigureInfo.AsrReviewInfo = &prohibitedAsrReviewTemplateInfo
		}
		if ocrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "ocr_review_info"); ok {
			prohibitedOcrReviewTemplateInfo := mps.ProhibitedOcrReviewTemplateInfo{}
			if v, ok := ocrReviewInfoMap["switch"]; ok {
				prohibitedOcrReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := ocrReviewInfoMap["block_confidence"]; ok {
				prohibitedOcrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := ocrReviewInfoMap["review_confidence"]; ok {
				prohibitedOcrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			prohibitedConfigureInfo.OcrReviewInfo = &prohibitedOcrReviewTemplateInfo
		}
		request.ProhibitedConfigure = &prohibitedConfigureInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "user_define_configure"); ok {
		userDefineConfigureInfo := mps.UserDefineConfigureInfo{}
		if faceReviewInfoMap, ok := helper.InterfaceToMap(dMap, "face_review_info"); ok {
			userDefineFaceReviewTemplateInfo := mps.UserDefineFaceReviewTemplateInfo{}
			if v, ok := faceReviewInfoMap["switch"]; ok {
				userDefineFaceReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := faceReviewInfoMap["label_set"]; ok {
				labelSetSet := v.(*schema.Set).List()
				for i := range labelSetSet {
					if labelSetSet[i] != nil {
						labelSet := labelSetSet[i].(string)
						userDefineFaceReviewTemplateInfo.LabelSet = append(userDefineFaceReviewTemplateInfo.LabelSet, &labelSet)
					}
				}
			}
			if v, ok := faceReviewInfoMap["block_confidence"]; ok {
				userDefineFaceReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := faceReviewInfoMap["review_confidence"]; ok {
				userDefineFaceReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			userDefineConfigureInfo.FaceReviewInfo = &userDefineFaceReviewTemplateInfo
		}
		if asrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "asr_review_info"); ok {
			userDefineAsrTextReviewTemplateInfo := mps.UserDefineAsrTextReviewTemplateInfo{}
			if v, ok := asrReviewInfoMap["switch"]; ok {
				userDefineAsrTextReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := asrReviewInfoMap["label_set"]; ok {
				labelSetSet := v.(*schema.Set).List()
				for i := range labelSetSet {
					if labelSetSet[i] != nil {
						labelSet := labelSetSet[i].(string)
						userDefineAsrTextReviewTemplateInfo.LabelSet = append(userDefineAsrTextReviewTemplateInfo.LabelSet, &labelSet)
					}
				}
			}
			if v, ok := asrReviewInfoMap["block_confidence"]; ok {
				userDefineAsrTextReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := asrReviewInfoMap["review_confidence"]; ok {
				userDefineAsrTextReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			userDefineConfigureInfo.AsrReviewInfo = &userDefineAsrTextReviewTemplateInfo
		}
		if ocrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "ocr_review_info"); ok {
			userDefineOcrTextReviewTemplateInfo := mps.UserDefineOcrTextReviewTemplateInfo{}
			if v, ok := ocrReviewInfoMap["switch"]; ok {
				userDefineOcrTextReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			// if v, ok := ocrReviewInfoMap["label_set"]; ok {
			// 	userDefineOcrTextReviewTemplateInfo.LabelSet = []*string{helper.String(v.(string))}
			// }
			if v, ok := ocrReviewInfoMap["label_set"]; ok {
				labelSetSet := v.(*schema.Set).List()
				for i := range labelSetSet {
					if labelSetSet[i] != nil {
						labelSet := labelSetSet[i].(string)
						userDefineOcrTextReviewTemplateInfo.LabelSet = append(userDefineOcrTextReviewTemplateInfo.LabelSet, &labelSet)
					}
				}
			}
			if v, ok := ocrReviewInfoMap["block_confidence"]; ok {
				userDefineOcrTextReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := ocrReviewInfoMap["review_confidence"]; ok {
				userDefineOcrTextReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			userDefineConfigureInfo.OcrReviewInfo = &userDefineOcrTextReviewTemplateInfo
		}
		request.UserDefineConfigure = &userDefineConfigureInfo
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().CreateContentReviewTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps contentReviewTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition = helper.Int64ToStr(*response.Response.Definition)
	d.SetId(definition)

	return resourceTencentCloudMpsContentReviewTemplateRead(d, meta)
}

func resourceTencentCloudMpsContentReviewTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_content_review_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	definition := d.Id()

	contentReviewTemplate, err := service.DescribeMpsContentReviewTemplateById(ctx, definition)
	if err != nil {
		return err
	}

	if contentReviewTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsContentReviewTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if contentReviewTemplate.Name != nil {
		_ = d.Set("name", contentReviewTemplate.Name)
	}

	if contentReviewTemplate.Comment != nil {
		_ = d.Set("comment", contentReviewTemplate.Comment)
	}

	if contentReviewTemplate.PornConfigure != nil {
		pornConfigureMap := map[string]interface{}{}

		if contentReviewTemplate.PornConfigure.ImgReviewInfo != nil {
			imgReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.PornConfigure.ImgReviewInfo.Switch != nil {
				imgReviewInfoMap["switch"] = contentReviewTemplate.PornConfigure.ImgReviewInfo.Switch
			}

			if contentReviewTemplate.PornConfigure.ImgReviewInfo.LabelSet != nil {
				imgReviewInfoMap["label_set"] = contentReviewTemplate.PornConfigure.ImgReviewInfo.LabelSet
			}

			if contentReviewTemplate.PornConfigure.ImgReviewInfo.BlockConfidence != nil {
				imgReviewInfoMap["block_confidence"] = contentReviewTemplate.PornConfigure.ImgReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.PornConfigure.ImgReviewInfo.ReviewConfidence != nil {
				imgReviewInfoMap["review_confidence"] = contentReviewTemplate.PornConfigure.ImgReviewInfo.ReviewConfidence
			}

			pornConfigureMap["img_review_info"] = []interface{}{imgReviewInfoMap}
		}

		if contentReviewTemplate.PornConfigure.AsrReviewInfo != nil {
			asrReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.PornConfigure.AsrReviewInfo.Switch != nil {
				asrReviewInfoMap["switch"] = contentReviewTemplate.PornConfigure.AsrReviewInfo.Switch
			}

			if contentReviewTemplate.PornConfigure.AsrReviewInfo.BlockConfidence != nil {
				asrReviewInfoMap["block_confidence"] = contentReviewTemplate.PornConfigure.AsrReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.PornConfigure.AsrReviewInfo.ReviewConfidence != nil {
				asrReviewInfoMap["review_confidence"] = contentReviewTemplate.PornConfigure.AsrReviewInfo.ReviewConfidence
			}

			pornConfigureMap["asr_review_info"] = []interface{}{asrReviewInfoMap}
		}

		if contentReviewTemplate.PornConfigure.OcrReviewInfo != nil {
			ocrReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.PornConfigure.OcrReviewInfo.Switch != nil {
				ocrReviewInfoMap["switch"] = contentReviewTemplate.PornConfigure.OcrReviewInfo.Switch
			}

			if contentReviewTemplate.PornConfigure.OcrReviewInfo.BlockConfidence != nil {
				ocrReviewInfoMap["block_confidence"] = contentReviewTemplate.PornConfigure.OcrReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.PornConfigure.OcrReviewInfo.ReviewConfidence != nil {
				ocrReviewInfoMap["review_confidence"] = contentReviewTemplate.PornConfigure.OcrReviewInfo.ReviewConfidence
			}

			pornConfigureMap["ocr_review_info"] = []interface{}{ocrReviewInfoMap}
		}

		_ = d.Set("porn_configure", []interface{}{pornConfigureMap})
	}

	if contentReviewTemplate.TerrorismConfigure != nil {
		terrorismConfigureMap := map[string]interface{}{}

		if contentReviewTemplate.TerrorismConfigure.ImgReviewInfo != nil {
			imgReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.TerrorismConfigure.ImgReviewInfo.Switch != nil {
				imgReviewInfoMap["switch"] = contentReviewTemplate.TerrorismConfigure.ImgReviewInfo.Switch
			}

			if contentReviewTemplate.TerrorismConfigure.ImgReviewInfo.LabelSet != nil {
				imgReviewInfoMap["label_set"] = contentReviewTemplate.TerrorismConfigure.ImgReviewInfo.LabelSet
			}

			if contentReviewTemplate.TerrorismConfigure.ImgReviewInfo.BlockConfidence != nil {
				imgReviewInfoMap["block_confidence"] = contentReviewTemplate.TerrorismConfigure.ImgReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.TerrorismConfigure.ImgReviewInfo.ReviewConfidence != nil {
				imgReviewInfoMap["review_confidence"] = contentReviewTemplate.TerrorismConfigure.ImgReviewInfo.ReviewConfidence
			}

			terrorismConfigureMap["img_review_info"] = []interface{}{imgReviewInfoMap}
		}

		if contentReviewTemplate.TerrorismConfigure.OcrReviewInfo != nil {
			ocrReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.TerrorismConfigure.OcrReviewInfo.Switch != nil {
				ocrReviewInfoMap["switch"] = contentReviewTemplate.TerrorismConfigure.OcrReviewInfo.Switch
			}

			if contentReviewTemplate.TerrorismConfigure.OcrReviewInfo.BlockConfidence != nil {
				ocrReviewInfoMap["block_confidence"] = contentReviewTemplate.TerrorismConfigure.OcrReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.TerrorismConfigure.OcrReviewInfo.ReviewConfidence != nil {
				ocrReviewInfoMap["review_confidence"] = contentReviewTemplate.TerrorismConfigure.OcrReviewInfo.ReviewConfidence
			}

			terrorismConfigureMap["ocr_review_info"] = []interface{}{ocrReviewInfoMap}
		}

		_ = d.Set("terrorism_configure", []interface{}{terrorismConfigureMap})
	}

	if contentReviewTemplate.PoliticalConfigure != nil {
		politicalConfigureMap := map[string]interface{}{}

		if contentReviewTemplate.PoliticalConfigure.ImgReviewInfo != nil {
			imgReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.PoliticalConfigure.ImgReviewInfo.Switch != nil {
				imgReviewInfoMap["switch"] = contentReviewTemplate.PoliticalConfigure.ImgReviewInfo.Switch
			}

			if contentReviewTemplate.PoliticalConfigure.ImgReviewInfo.LabelSet != nil {
				imgReviewInfoMap["label_set"] = contentReviewTemplate.PoliticalConfigure.ImgReviewInfo.LabelSet
			}

			if contentReviewTemplate.PoliticalConfigure.ImgReviewInfo.BlockConfidence != nil {
				imgReviewInfoMap["block_confidence"] = contentReviewTemplate.PoliticalConfigure.ImgReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.PoliticalConfigure.ImgReviewInfo.ReviewConfidence != nil {
				imgReviewInfoMap["review_confidence"] = contentReviewTemplate.PoliticalConfigure.ImgReviewInfo.ReviewConfidence
			}

			politicalConfigureMap["img_review_info"] = []interface{}{imgReviewInfoMap}
		}

		if contentReviewTemplate.PoliticalConfigure.AsrReviewInfo != nil {
			asrReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.PoliticalConfigure.AsrReviewInfo.Switch != nil {
				asrReviewInfoMap["switch"] = contentReviewTemplate.PoliticalConfigure.AsrReviewInfo.Switch
			}

			if contentReviewTemplate.PoliticalConfigure.AsrReviewInfo.BlockConfidence != nil {
				asrReviewInfoMap["block_confidence"] = contentReviewTemplate.PoliticalConfigure.AsrReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.PoliticalConfigure.AsrReviewInfo.ReviewConfidence != nil {
				asrReviewInfoMap["review_confidence"] = contentReviewTemplate.PoliticalConfigure.AsrReviewInfo.ReviewConfidence
			}

			politicalConfigureMap["asr_review_info"] = []interface{}{asrReviewInfoMap}
		}

		if contentReviewTemplate.PoliticalConfigure.OcrReviewInfo != nil {
			ocrReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.PoliticalConfigure.OcrReviewInfo.Switch != nil {
				ocrReviewInfoMap["switch"] = contentReviewTemplate.PoliticalConfigure.OcrReviewInfo.Switch
			}

			if contentReviewTemplate.PoliticalConfigure.OcrReviewInfo.BlockConfidence != nil {
				ocrReviewInfoMap["block_confidence"] = contentReviewTemplate.PoliticalConfigure.OcrReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.PoliticalConfigure.OcrReviewInfo.ReviewConfidence != nil {
				ocrReviewInfoMap["review_confidence"] = contentReviewTemplate.PoliticalConfigure.OcrReviewInfo.ReviewConfidence
			}

			politicalConfigureMap["ocr_review_info"] = []interface{}{ocrReviewInfoMap}
		}

		_ = d.Set("political_configure", []interface{}{politicalConfigureMap})
	}

	if contentReviewTemplate.ProhibitedConfigure != nil {
		prohibitedConfigureMap := map[string]interface{}{}

		if contentReviewTemplate.ProhibitedConfigure.AsrReviewInfo != nil {
			asrReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.ProhibitedConfigure.AsrReviewInfo.Switch != nil {
				asrReviewInfoMap["switch"] = contentReviewTemplate.ProhibitedConfigure.AsrReviewInfo.Switch
			}

			if contentReviewTemplate.ProhibitedConfigure.AsrReviewInfo.BlockConfidence != nil {
				asrReviewInfoMap["block_confidence"] = contentReviewTemplate.ProhibitedConfigure.AsrReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.ProhibitedConfigure.AsrReviewInfo.ReviewConfidence != nil {
				asrReviewInfoMap["review_confidence"] = contentReviewTemplate.ProhibitedConfigure.AsrReviewInfo.ReviewConfidence
			}

			prohibitedConfigureMap["asr_review_info"] = []interface{}{asrReviewInfoMap}
		}

		if contentReviewTemplate.ProhibitedConfigure.OcrReviewInfo != nil {
			ocrReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.ProhibitedConfigure.OcrReviewInfo.Switch != nil {
				ocrReviewInfoMap["switch"] = contentReviewTemplate.ProhibitedConfigure.OcrReviewInfo.Switch
			}

			if contentReviewTemplate.ProhibitedConfigure.OcrReviewInfo.BlockConfidence != nil {
				ocrReviewInfoMap["block_confidence"] = contentReviewTemplate.ProhibitedConfigure.OcrReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.ProhibitedConfigure.OcrReviewInfo.ReviewConfidence != nil {
				ocrReviewInfoMap["review_confidence"] = contentReviewTemplate.ProhibitedConfigure.OcrReviewInfo.ReviewConfidence
			}

			prohibitedConfigureMap["ocr_review_info"] = []interface{}{ocrReviewInfoMap}
		}

		_ = d.Set("prohibited_configure", []interface{}{prohibitedConfigureMap})
	}

	if contentReviewTemplate.UserDefineConfigure != nil {
		userDefineConfigureMap := map[string]interface{}{}

		if contentReviewTemplate.UserDefineConfigure.FaceReviewInfo != nil {
			faceReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.UserDefineConfigure.FaceReviewInfo.Switch != nil {
				faceReviewInfoMap["switch"] = contentReviewTemplate.UserDefineConfigure.FaceReviewInfo.Switch
			}

			if contentReviewTemplate.UserDefineConfigure.FaceReviewInfo.LabelSet != nil {
				faceReviewInfoMap["label_set"] = contentReviewTemplate.UserDefineConfigure.FaceReviewInfo.LabelSet
			}

			if contentReviewTemplate.UserDefineConfigure.FaceReviewInfo.BlockConfidence != nil {
				faceReviewInfoMap["block_confidence"] = contentReviewTemplate.UserDefineConfigure.FaceReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.UserDefineConfigure.FaceReviewInfo.ReviewConfidence != nil {
				faceReviewInfoMap["review_confidence"] = contentReviewTemplate.UserDefineConfigure.FaceReviewInfo.ReviewConfidence
			}

			userDefineConfigureMap["face_review_info"] = []interface{}{faceReviewInfoMap}
		}

		if contentReviewTemplate.UserDefineConfigure.AsrReviewInfo != nil {
			asrReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.UserDefineConfigure.AsrReviewInfo.Switch != nil {
				asrReviewInfoMap["switch"] = contentReviewTemplate.UserDefineConfigure.AsrReviewInfo.Switch
			}

			if contentReviewTemplate.UserDefineConfigure.AsrReviewInfo.LabelSet != nil {
				asrReviewInfoMap["label_set"] = contentReviewTemplate.UserDefineConfigure.AsrReviewInfo.LabelSet
			}

			if contentReviewTemplate.UserDefineConfigure.AsrReviewInfo.BlockConfidence != nil {
				asrReviewInfoMap["block_confidence"] = contentReviewTemplate.UserDefineConfigure.AsrReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.UserDefineConfigure.AsrReviewInfo.ReviewConfidence != nil {
				asrReviewInfoMap["review_confidence"] = contentReviewTemplate.UserDefineConfigure.AsrReviewInfo.ReviewConfidence
			}

			userDefineConfigureMap["asr_review_info"] = []interface{}{asrReviewInfoMap}
		}

		if contentReviewTemplate.UserDefineConfigure.OcrReviewInfo != nil {
			ocrReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.UserDefineConfigure.OcrReviewInfo.Switch != nil {
				ocrReviewInfoMap["switch"] = contentReviewTemplate.UserDefineConfigure.OcrReviewInfo.Switch
			}

			if contentReviewTemplate.UserDefineConfigure.OcrReviewInfo.LabelSet != nil {
				// labelSet:=contentReviewTemplate.UserDefineConfigure.OcrReviewInfo.LabelSet
				// var ret string
				// for _, label:=range labelSet {
				// 	ret += *label
				// }
				// ocrReviewInfoMap["label_set"] = ret
				ocrReviewInfoMap["label_set"] = contentReviewTemplate.UserDefineConfigure.OcrReviewInfo.LabelSet
			}

			if contentReviewTemplate.UserDefineConfigure.OcrReviewInfo.BlockConfidence != nil {
				ocrReviewInfoMap["block_confidence"] = contentReviewTemplate.UserDefineConfigure.OcrReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.UserDefineConfigure.OcrReviewInfo.ReviewConfidence != nil {
				ocrReviewInfoMap["review_confidence"] = contentReviewTemplate.UserDefineConfigure.OcrReviewInfo.ReviewConfidence
			}

			userDefineConfigureMap["ocr_review_info"] = []interface{}{ocrReviewInfoMap}
		}

		_ = d.Set("user_define_configure", []interface{}{userDefineConfigureMap})
	}

	return nil
}

func resourceTencentCloudMpsContentReviewTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_content_review_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := mps.NewModifyContentReviewTemplateRequest()

	definition := d.Id()

	request.Definition = helper.StrToInt64Point(definition)

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}
	}

	if d.HasChange("porn_configure") {
		if dMap, ok := helper.InterfacesHeadMap(d, "porn_configure"); ok {
			pornConfigureInfo := mps.PornConfigureInfoForUpdate{}
			if imgReviewInfoMap, ok := helper.InterfaceToMap(dMap, "img_review_info"); ok {
				pornImgReviewTemplateInfo := mps.PornImgReviewTemplateInfoForUpdate{}
				if v, ok := imgReviewInfoMap["switch"]; ok {
					pornImgReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := imgReviewInfoMap["label_set"]; ok {
					labelSetSet := v.(*schema.Set).List()
					for i := range labelSetSet {
						if labelSetSet[i] != nil {
							labelSet := labelSetSet[i].(string)
							pornImgReviewTemplateInfo.LabelSet = append(pornImgReviewTemplateInfo.LabelSet, &labelSet)
						}
					}
				}
				if v, ok := imgReviewInfoMap["block_confidence"]; ok {
					pornImgReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := imgReviewInfoMap["review_confidence"]; ok {
					pornImgReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				pornConfigureInfo.ImgReviewInfo = &pornImgReviewTemplateInfo
			}
			if asrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "asr_review_info"); ok {
				pornAsrReviewTemplateInfo := mps.PornAsrReviewTemplateInfoForUpdate{}
				if v, ok := asrReviewInfoMap["switch"]; ok {
					pornAsrReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := asrReviewInfoMap["block_confidence"]; ok {
					pornAsrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := asrReviewInfoMap["review_confidence"]; ok {
					pornAsrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				pornConfigureInfo.AsrReviewInfo = &pornAsrReviewTemplateInfo
			}
			if ocrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "ocr_review_info"); ok {
				pornOcrReviewTemplateInfo := mps.PornOcrReviewTemplateInfoForUpdate{}
				if v, ok := ocrReviewInfoMap["switch"]; ok {
					pornOcrReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := ocrReviewInfoMap["block_confidence"]; ok {
					pornOcrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := ocrReviewInfoMap["review_confidence"]; ok {
					pornOcrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				pornConfigureInfo.OcrReviewInfo = &pornOcrReviewTemplateInfo
			}
			request.PornConfigure = &pornConfigureInfo
		}
	}

	if d.HasChange("terrorism_configure") {
		if dMap, ok := helper.InterfacesHeadMap(d, "terrorism_configure"); ok {
			terrorismConfigureInfo := mps.TerrorismConfigureInfoForUpdate{}
			if imgReviewInfoMap, ok := helper.InterfaceToMap(dMap, "img_review_info"); ok {
				terrorismImgReviewTemplateInfo := mps.TerrorismImgReviewTemplateInfoForUpdate{}
				if v, ok := imgReviewInfoMap["switch"]; ok {
					terrorismImgReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := imgReviewInfoMap["label_set"]; ok {
					labelSetSet := v.(*schema.Set).List()
					for i := range labelSetSet {
						if labelSetSet[i] != nil {
							labelSet := labelSetSet[i].(string)
							terrorismImgReviewTemplateInfo.LabelSet = append(terrorismImgReviewTemplateInfo.LabelSet, &labelSet)
						}
					}
				}
				if v, ok := imgReviewInfoMap["block_confidence"]; ok {
					terrorismImgReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := imgReviewInfoMap["review_confidence"]; ok {
					terrorismImgReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				terrorismConfigureInfo.ImgReviewInfo = &terrorismImgReviewTemplateInfo
			}
			if ocrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "ocr_review_info"); ok {
				terrorismOcrReviewTemplateInfo := mps.TerrorismOcrReviewTemplateInfoForUpdate{}
				if v, ok := ocrReviewInfoMap["switch"]; ok {
					terrorismOcrReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := ocrReviewInfoMap["block_confidence"]; ok {
					terrorismOcrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := ocrReviewInfoMap["review_confidence"]; ok {
					terrorismOcrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				terrorismConfigureInfo.OcrReviewInfo = &terrorismOcrReviewTemplateInfo
			}
			request.TerrorismConfigure = &terrorismConfigureInfo
		}
	}

	if d.HasChange("political_configure") {
		if dMap, ok := helper.InterfacesHeadMap(d, "political_configure"); ok {
			politicalConfigureInfo := mps.PoliticalConfigureInfoForUpdate{}
			if imgReviewInfoMap, ok := helper.InterfaceToMap(dMap, "img_review_info"); ok {
				politicalImgReviewTemplateInfo := mps.PoliticalImgReviewTemplateInfoForUpdate{}
				if v, ok := imgReviewInfoMap["switch"]; ok {
					politicalImgReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := imgReviewInfoMap["label_set"]; ok {
					labelSetSet := v.(*schema.Set).List()
					for i := range labelSetSet {
						if labelSetSet[i] != nil {
							labelSet := labelSetSet[i].(string)
							politicalImgReviewTemplateInfo.LabelSet = append(politicalImgReviewTemplateInfo.LabelSet, &labelSet)
						}
					}
				}
				if v, ok := imgReviewInfoMap["block_confidence"]; ok {
					politicalImgReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := imgReviewInfoMap["review_confidence"]; ok {
					politicalImgReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				politicalConfigureInfo.ImgReviewInfo = &politicalImgReviewTemplateInfo
			}
			if asrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "asr_review_info"); ok {
				politicalAsrReviewTemplateInfo := mps.PoliticalAsrReviewTemplateInfoForUpdate{}
				if v, ok := asrReviewInfoMap["switch"]; ok {
					politicalAsrReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := asrReviewInfoMap["block_confidence"]; ok {
					politicalAsrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := asrReviewInfoMap["review_confidence"]; ok {
					politicalAsrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				politicalConfigureInfo.AsrReviewInfo = &politicalAsrReviewTemplateInfo
			}
			if ocrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "ocr_review_info"); ok {
				politicalOcrReviewTemplateInfo := mps.PoliticalOcrReviewTemplateInfoForUpdate{}
				if v, ok := ocrReviewInfoMap["switch"]; ok {
					politicalOcrReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := ocrReviewInfoMap["block_confidence"]; ok {
					politicalOcrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := ocrReviewInfoMap["review_confidence"]; ok {
					politicalOcrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				politicalConfigureInfo.OcrReviewInfo = &politicalOcrReviewTemplateInfo
			}
			request.PoliticalConfigure = &politicalConfigureInfo
		}
	}

	if d.HasChange("prohibited_configure") {
		if dMap, ok := helper.InterfacesHeadMap(d, "prohibited_configure"); ok {
			prohibitedConfigureInfo := mps.ProhibitedConfigureInfoForUpdate{}
			if asrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "asr_review_info"); ok {
				prohibitedAsrReviewTemplateInfo := mps.ProhibitedAsrReviewTemplateInfoForUpdate{}
				if v, ok := asrReviewInfoMap["switch"]; ok {
					prohibitedAsrReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := asrReviewInfoMap["block_confidence"]; ok {
					prohibitedAsrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := asrReviewInfoMap["review_confidence"]; ok {
					prohibitedAsrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				prohibitedConfigureInfo.AsrReviewInfo = &prohibitedAsrReviewTemplateInfo
			}
			if ocrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "ocr_review_info"); ok {
				prohibitedOcrReviewTemplateInfo := mps.ProhibitedOcrReviewTemplateInfoForUpdate{}
				if v, ok := ocrReviewInfoMap["switch"]; ok {
					prohibitedOcrReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := ocrReviewInfoMap["block_confidence"]; ok {
					prohibitedOcrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := ocrReviewInfoMap["review_confidence"]; ok {
					prohibitedOcrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				prohibitedConfigureInfo.OcrReviewInfo = &prohibitedOcrReviewTemplateInfo
			}
			request.ProhibitedConfigure = &prohibitedConfigureInfo
		}
	}

	if d.HasChange("user_define_configure") {
		if dMap, ok := helper.InterfacesHeadMap(d, "user_define_configure"); ok {
			userDefineConfigureInfo := mps.UserDefineConfigureInfoForUpdate{}
			if faceReviewInfoMap, ok := helper.InterfaceToMap(dMap, "face_review_info"); ok {
				userDefineFaceReviewTemplateInfo := mps.UserDefineFaceReviewTemplateInfoForUpdate{}
				if v, ok := faceReviewInfoMap["switch"]; ok {
					userDefineFaceReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := faceReviewInfoMap["label_set"]; ok {
					labelSetSet := v.(*schema.Set).List()
					for i := range labelSetSet {
						if labelSetSet[i] != nil {
							labelSet := labelSetSet[i].(string)
							userDefineFaceReviewTemplateInfo.LabelSet = append(userDefineFaceReviewTemplateInfo.LabelSet, &labelSet)
						}
					}
				}
				if v, ok := faceReviewInfoMap["block_confidence"]; ok {
					userDefineFaceReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := faceReviewInfoMap["review_confidence"]; ok {
					userDefineFaceReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				userDefineConfigureInfo.FaceReviewInfo = &userDefineFaceReviewTemplateInfo
			}
			if asrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "asr_review_info"); ok {
				userDefineAsrTextReviewTemplateInfo := mps.UserDefineAsrTextReviewTemplateInfoForUpdate{}
				if v, ok := asrReviewInfoMap["switch"]; ok {
					userDefineAsrTextReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := asrReviewInfoMap["label_set"]; ok {
					labelSetSet := v.(*schema.Set).List()
					for i := range labelSetSet {
						if labelSetSet[i] != nil {
							labelSet := labelSetSet[i].(string)
							userDefineAsrTextReviewTemplateInfo.LabelSet = append(userDefineAsrTextReviewTemplateInfo.LabelSet, &labelSet)
						}
					}
				}
				if v, ok := asrReviewInfoMap["block_confidence"]; ok {
					userDefineAsrTextReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := asrReviewInfoMap["review_confidence"]; ok {
					userDefineAsrTextReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				userDefineConfigureInfo.AsrReviewInfo = &userDefineAsrTextReviewTemplateInfo
			}
			if ocrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "ocr_review_info"); ok {
				userDefineOcrTextReviewTemplateInfo := mps.UserDefineOcrTextReviewTemplateInfoForUpdate{}
				if v, ok := ocrReviewInfoMap["switch"]; ok {
					userDefineOcrTextReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				// do not support to modify user_define_configure.ocr_review_info.label_set
				// if v, ok := ocrReviewInfoMap["label_set"]; ok {
				// 	labelSetSet := v.(*schema.Set).List()
				// 	for i := range labelSetSet {
				// 		if labelSetSet[i] != nil {
				// 			labelSet := labelSetSet[i].(string)
				// 			userDefineOcrTextReviewTemplateInfo.LabelSet = append(userDefineOcrTextReviewTemplateInfo.LabelSet, &labelSet)
				// 		}
				// 	}
				// }
				if v, ok := ocrReviewInfoMap["block_confidence"]; ok {
					userDefineOcrTextReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := ocrReviewInfoMap["review_confidence"]; ok {
					userDefineOcrTextReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				userDefineConfigureInfo.OcrReviewInfo = &userDefineOcrTextReviewTemplateInfo
			}
			request.UserDefineConfigure = &userDefineConfigureInfo
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().ModifyContentReviewTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mps contentReviewTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMpsContentReviewTemplateRead(d, meta)
}

func resourceTencentCloudMpsContentReviewTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_content_review_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	definition := d.Id()

	if err := service.DeleteMpsContentReviewTemplateById(ctx, definition); err != nil {
		return err
	}

	return nil
}
