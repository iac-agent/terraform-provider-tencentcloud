package eb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudEbEventTransform() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEbEventTransformCreate,
		Read:   resourceTencentCloudEbEventTransformRead,
		Update: resourceTencentCloudEbEventTransformUpdate,
		Delete: resourceTencentCloudEbEventTransformDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"event_bus_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "事件 bus ID。",
			},

			"rule_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ruleId。",
			},

			"transformations": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "A 列表 transformation 规则，currently 仅 一个。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"extraction": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Describe how 到 extract 数据。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"extraction_input_path": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "JsonPath，如果未指定， 默认值 $。",
									},
									"format": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "值: `TEXT`，`JSON`。",
									},
									"text_params": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Only Text needs 到 是 passed。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"separator": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "`Comma`，`|`，`tab`，`space`，`newline`，`%`，`#`， 限制 长度 是 1。",
												},
												"regex": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Fill 在 regular expression: 长度 128。",
												},
											},
										},
									},
								},
							},
						},
						"etl_filter": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Describe how 到 过滤器 数据。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filter": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Grammatical Rules 是 consistent。",
									},
								},
							},
						},
						"transform": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Describe how 到 convert 数据。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"output_structs": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Describe how 数据 是 transformed。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Corresponding 到 键 在 output json。",
												},
												"value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "You 可以 fill 在 json-路径 和 also support constants 或 built-在 keyword date types。",
												},
												"value_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "数据 类型 值，可选 值: `STRING`，`NUMBER`，`BOOLEAN`，`NULL`，`SYS_VARIABLE`，`JSONPATH`。",
												},
											},
										},
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

func resourceTencentCloudEbEventTransformCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eb_event_transform.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request          = eb.NewCreateTransformationRequest()
		response         = eb.NewCreateTransformationResponse()
		eventBusId       string
		ruleId           string
		transformationId string
	)
	if v, ok := d.GetOk("event_bus_id"); ok {
		eventBusId = v.(string)
		request.EventBusId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_id"); ok {
		ruleId = v.(string)
		request.RuleId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("transformations"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			transformation := eb.Transformation{}
			if extractionMap, ok := helper.InterfaceToMap(dMap, "extraction"); ok {
				extraction := eb.Extraction{}
				if v, ok := extractionMap["extraction_input_path"]; ok {
					extraction.ExtractionInputPath = helper.String(v.(string))
				}
				if v, ok := extractionMap["format"]; ok {
					extraction.Format = helper.String(v.(string))
				}
				if textParamsMap, ok := helper.InterfaceToMap(extractionMap, "text_params"); ok {
					textParams := eb.TextParams{}
					if v, ok := textParamsMap["separator"]; ok {
						textParams.Separator = helper.String(v.(string))
					}
					if v, ok := textParamsMap["regex"]; ok {
						textParams.Regex = helper.String(v.(string))
					}
					extraction.TextParams = &textParams
				}
				transformation.Extraction = &extraction
			}
			if etlFilterMap, ok := helper.InterfaceToMap(dMap, "etl_filter"); ok {
				etlFilter := eb.EtlFilter{}
				if v, ok := etlFilterMap["filter"]; ok {
					etlFilter.Filter = helper.String(v.(string))
				}
				transformation.EtlFilter = &etlFilter
			}
			if transformMap, ok := helper.InterfaceToMap(dMap, "transform"); ok {
				transform := eb.Transform{}
				if v, ok := transformMap["output_structs"]; ok {
					for _, item := range v.([]interface{}) {
						outputStructsMap := item.(map[string]interface{})
						outputStructParam := eb.OutputStructParam{}
						if v, ok := outputStructsMap["key"]; ok {
							outputStructParam.Key = helper.String(v.(string))
						}
						if v, ok := outputStructsMap["value"]; ok {
							outputStructParam.Value = helper.String(v.(string))
						}
						if v, ok := outputStructsMap["value_type"]; ok {
							outputStructParam.ValueType = helper.String(v.(string))
						}
						transform.OutputStructs = append(transform.OutputStructs, &outputStructParam)
					}
				}
				transformation.Transform = &transform
			}
			request.Transformations = append(request.Transformations, &transformation)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEbClient().CreateTransformation(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create eb ebTransform failed, reason:%+v", logId, err)
		return err
	}

	transformationId = *response.Response.TransformationId
	d.SetId(eventBusId + tccommon.FILED_SP + ruleId + tccommon.FILED_SP + transformationId)

	return resourceTencentCloudEbEventTransformRead(d, meta)
}

func resourceTencentCloudEbEventTransformRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eb_event_transform.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := EbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	eventBusId := idSplit[0]
	ruleId := idSplit[1]
	transformationId := idSplit[2]

	ebTransform, err := service.DescribeEbEventTransformById(ctx, eventBusId, ruleId, transformationId)
	if err != nil {
		return err
	}

	if ebTransform == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `EbEventTransform` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("event_bus_id", eventBusId)
	_ = d.Set("rule_id", ruleId)

	if ebTransform != nil {
		transformationsMap := map[string]interface{}{}

		if ebTransform.Extraction != nil {
			extractionMap := map[string]interface{}{}

			if ebTransform.Extraction.ExtractionInputPath != nil {
				extractionMap["extraction_input_path"] = ebTransform.Extraction.ExtractionInputPath
			}

			if ebTransform.Extraction.Format != nil {
				extractionMap["format"] = ebTransform.Extraction.Format
			}

			if ebTransform.Extraction.TextParams != nil {
				textParamsMap := map[string]interface{}{}

				if ebTransform.Extraction.TextParams.Separator != nil {
					textParamsMap["separator"] = ebTransform.Extraction.TextParams.Separator
				}

				if ebTransform.Extraction.TextParams.Regex != nil {
					textParamsMap["regex"] = ebTransform.Extraction.TextParams.Regex
				}

				extractionMap["text_params"] = []interface{}{textParamsMap}
			}

			transformationsMap["extraction"] = []interface{}{extractionMap}
		}

		if ebTransform.EtlFilter != nil {
			etlFilterMap := map[string]interface{}{}

			if ebTransform.EtlFilter.Filter != nil {
				etlFilterMap["filter"] = ebTransform.EtlFilter.Filter
			}

			transformationsMap["etl_filter"] = []interface{}{etlFilterMap}
		}

		if ebTransform.Transform != nil {
			transformMap := map[string]interface{}{}

			if ebTransform.Transform.OutputStructs != nil {
				outputStructsList := []interface{}{}
				for _, outputStructs := range ebTransform.Transform.OutputStructs {
					outputStructsMap := map[string]interface{}{}

					if outputStructs.Key != nil {
						outputStructsMap["key"] = outputStructs.Key
					}

					if outputStructs.Value != nil {
						outputStructsMap["value"] = outputStructs.Value
					}

					if outputStructs.ValueType != nil {
						outputStructsMap["value_type"] = outputStructs.ValueType
					}

					outputStructsList = append(outputStructsList, outputStructsMap)
				}

				transformMap["output_structs"] = outputStructsList
			}

			transformationsMap["transform"] = []interface{}{transformMap}
		}

		_ = d.Set("transformations", []interface{}{transformationsMap})

	}

	return nil
}

func resourceTencentCloudEbEventTransformUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eb_event_transform.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := eb.NewUpdateTransformationRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	eventBusId := idSplit[0]
	ruleId := idSplit[1]
	transformationId := idSplit[2]

	request.EventBusId = &eventBusId
	request.RuleId = &ruleId
	request.TransformationId = &transformationId

	immutableArgs := []string{"event_bus_id", "rule_id"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("transformations") {
		if v, ok := d.GetOk("transformations"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				transformation := eb.Transformation{}
				if extractionMap, ok := helper.InterfaceToMap(dMap, "extraction"); ok {
					extraction := eb.Extraction{}
					if v, ok := extractionMap["extraction_input_path"]; ok {
						extraction.ExtractionInputPath = helper.String(v.(string))
					}
					if v, ok := extractionMap["format"]; ok {
						extraction.Format = helper.String(v.(string))
					}
					if textParamsMap, ok := helper.InterfaceToMap(extractionMap, "text_params"); ok {
						textParams := eb.TextParams{}
						if v, ok := textParamsMap["separator"]; ok {
							textParams.Separator = helper.String(v.(string))
						}
						if v, ok := textParamsMap["regex"]; ok {
							textParams.Regex = helper.String(v.(string))
						}
						extraction.TextParams = &textParams
					}
					transformation.Extraction = &extraction
				}
				if etlFilterMap, ok := helper.InterfaceToMap(dMap, "etl_filter"); ok {
					etlFilter := eb.EtlFilter{}
					if v, ok := etlFilterMap["filter"]; ok {
						etlFilter.Filter = helper.String(v.(string))
					}
					transformation.EtlFilter = &etlFilter
				}
				if transformMap, ok := helper.InterfaceToMap(dMap, "transform"); ok {
					transform := eb.Transform{}
					if v, ok := transformMap["output_structs"]; ok {
						for _, item := range v.([]interface{}) {
							outputStructsMap := item.(map[string]interface{})
							outputStructParam := eb.OutputStructParam{}
							if v, ok := outputStructsMap["key"]; ok {
								outputStructParam.Key = helper.String(v.(string))
							}
							if v, ok := outputStructsMap["value"]; ok {
								outputStructParam.Value = helper.String(v.(string))
							}
							if v, ok := outputStructsMap["value_type"]; ok {
								outputStructParam.ValueType = helper.String(v.(string))
							}
							transform.OutputStructs = append(transform.OutputStructs, &outputStructParam)
						}
					}
					transformation.Transform = &transform
				}
				request.Transformations = append(request.Transformations, &transformation)
			}
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEbClient().UpdateTransformation(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update eb ebTransform failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudEbEventTransformRead(d, meta)
}

func resourceTencentCloudEbEventTransformDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eb_event_transform.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := EbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	eventBusId := idSplit[0]
	ruleId := idSplit[1]
	transformationId := idSplit[2]

	if err := service.DeleteEbEventTransformById(ctx, eventBusId, ruleId, transformationId); err != nil {
		return err
	}

	return nil
}
