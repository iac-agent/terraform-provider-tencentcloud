package cls

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsIndex() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsIndexCreate,
		Read:   resourceTencentCloudClsIndexRead,
		Update: resourceTencentCloudClsIndexUpdate,
		Delete: resourceTencentCloudClsIndexDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"topic_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "日志主题ID。",
			},
			"rule": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "索引规则。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"full_text": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "全文索引配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"case_sensitive": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "区分大小写。",
									},
									"tokenizer": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "全文索引分隔符。字符串中的每个字符代表一个分隔符。",
									},
									"contain_z_h": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "是否包含汉字。",
									},
								},
							},
						},
						"key_value": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "键值索引配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"case_sensitive": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "区分大小写。",
									},
									"key_values": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "待创建索引的键值对信息。最多可以配置 100 个键值对。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Required: true,
													Description: "当需要为字段配置键值或元字段索引时，元字段Key不需要以__TAG__为前缀。并且是一致的" +
														"with the one when logs are uploaded. __TAG__. will be prefixed automatically for display in the console..",
												},
												"value": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "字段索引描述信息。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "字段类型。有效值：long、text、double。",
															},
															"tokenizer": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "字段分隔符，仅当字段类型为文本时才有意义。输入字符串中的每个字符代表一个分隔符。",
															},
															"sql_flag": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "字段是否启用分析功能。",
															},
															"contain_z_h": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "是否包含汉字。",
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
						"tag": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "元字段索引配置。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"case_sensitive": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "区分大小写。",
									},
									"key_values": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "待创建索引的键值对信息。最多可以配置 100 个键值对。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Required: true,
													Description: "当需要为字段配置键值或元字段索引时，元字段Key不需要以__TAG__为前缀。并且是一致的" +
														"with the one when logs are uploaded. __TAG__. will be prefixed automatically for display in the console..",
												},
												"value": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "字段索引描述信息。",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "字段类型。有效值：long、text、double。",
															},
															"tokenizer": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "字段分隔符，仅当字段类型为文本时才有意义。输入字符串中的每个字符代表一个分隔符。",
															},
															"sql_flag": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "字段是否启用分析功能。",
															},
															"contain_z_h": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "是否包含汉字。",
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
						"dynamic_index": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "键值索引是自动配置的。如果为空，则表示该功能未启用。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "索引自动配置开关。",
									},
								},
							},
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "是否生效。默认值：true。",
			},
			"include_internal_fields": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "全文索引的内部字段标记。默认值：假。有效值： false：排除内部字段； true：包括内部字段。",
			},
			"metadata_flag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "元数据标志。默认值：0。有效值：0：全文索引（包括启用键值索引的元数据字段）； 1：全文索引（包括所有元数据字段）； 2：全文索引（不包括元数据字段）.。",
			},
		},
	}
}

func resourceTencentCloudClsIndexCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_index.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = cls.NewCreateIndexRequest()
		indexId string
	)

	if v, ok := d.GetOk("topic_id"); ok {
		request.TopicId = helper.String(v.(string))
		indexId = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "rule"); ok {
		ruleInfo := cls.RuleInfo{}
		if fullTextMap, ok := helper.InterfaceToMap(dMap, "full_text"); ok {
			fullTextInfo := cls.FullTextInfo{}
			if v, ok := fullTextMap["case_sensitive"]; ok {
				fullTextInfo.CaseSensitive = helper.Bool(v.(bool))
			}

			if v, ok := fullTextMap["tokenizer"]; ok {
				fullTextInfo.Tokenizer = helper.String(v.(string))
			}

			if v, ok := fullTextMap["contain_z_h"]; ok {
				fullTextInfo.ContainZH = helper.Bool(v.(bool))
			}

			ruleInfo.FullText = &fullTextInfo
		}

		if ruleKeyValueMap, ok := helper.InterfaceToMap(dMap, "key_value"); ok {
			ruleKeyValueInfo := cls.RuleKeyValueInfo{}
			if v, ok := ruleKeyValueMap["case_sensitive"]; ok {
				ruleKeyValueInfo.CaseSensitive = helper.Bool(v.(bool))
			}

			if v, ok := ruleKeyValueMap["key_values"]; ok {
				for _, keyValue := range v.([]interface{}) {
					keyValueInfo := cls.KeyValueInfo{}
					keyValueMap := keyValue.(map[string]interface{})
					if v, ok := keyValueMap["key"]; ok {
						keyValueInfo.Key = helper.String(v.(string))
					}

					if valueMap, ok := helper.InterfaceToMap(keyValueMap, "value"); ok {
						valueInfo := cls.ValueInfo{}
						if v, ok := valueMap["type"]; ok {
							valueInfo.Type = helper.String(v.(string))
						}

						if v, ok := valueMap["tokenizer"]; ok {
							valueInfo.Tokenizer = helper.String(v.(string))
						}

						if v, ok := valueMap["sql_flag"]; ok {
							valueInfo.SqlFlag = helper.Bool(v.(bool))
						}

						if v, ok := valueMap["contain_z_h"]; ok {
							valueInfo.ContainZH = helper.Bool(v.(bool))
						}

						keyValueInfo.Value = &valueInfo
					}

					ruleKeyValueInfo.KeyValues = append(ruleKeyValueInfo.KeyValues, &keyValueInfo)
				}
			}

			ruleInfo.KeyValue = &ruleKeyValueInfo
		}

		if tagMap, ok := helper.InterfaceToMap(dMap, "tag"); ok {
			ruleTagInfo := cls.RuleTagInfo{}
			if v, ok := tagMap["case_sensitive"]; ok {
				ruleTagInfo.CaseSensitive = helper.Bool(v.(bool))
			}

			if v, ok := tagMap["key_values"]; ok {
				for _, keyValue := range v.([]interface{}) {
					keyValueInfo := cls.KeyValueInfo{}
					keyValueMap := keyValue.(map[string]interface{})
					if v, ok := keyValueMap["key"]; ok {
						keyValueInfo.Key = helper.String(v.(string))
					}

					if valueMap, ok := helper.InterfaceToMap(keyValueMap, "value"); ok {
						valueInfo := cls.ValueInfo{}
						if v, ok := valueMap["type"]; ok {
							valueInfo.Type = helper.String(v.(string))
						}

						if v, ok := valueMap["tokenizer"]; ok {
							valueInfo.Tokenizer = helper.String(v.(string))
						}

						if v, ok := valueMap["sql_flag"]; ok {
							valueInfo.SqlFlag = helper.Bool(v.(bool))
						}

						if v, ok := valueMap["contain_z_h"]; ok {
							valueInfo.ContainZH = helper.Bool(v.(bool))
						}

						keyValueInfo.Value = &valueInfo
					}

					ruleTagInfo.KeyValues = append(ruleTagInfo.KeyValues, &keyValueInfo)
				}
			}

			ruleInfo.Tag = &ruleTagInfo
		}

		if dynamicIndexMap, ok := helper.InterfaceToMap(dMap, "dynamic_index"); ok {
			dynamicIndexInfo := cls.DynamicIndex{}
			if v, ok := dynamicIndexMap["status"]; ok {
				dynamicIndexInfo.Status = helper.Bool(v.(bool))
			}

			ruleInfo.DynamicIndex = &dynamicIndexInfo
		}

		request.Rule = &ruleInfo
	}

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("include_internal_fields"); ok {
		request.IncludeInternalFields = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("metadata_flag"); ok {
		request.MetadataFlag = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateIndex(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls index failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(indexId)

	return resourceTencentCloudClsIndexRead(d, meta)
}

func resourceTencentCloudClsIndexRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_index.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = cls.NewDescribeIndexRequest()
		result  *cls.DescribeIndexResponse
		id      = d.Id()
	)

	request.TopicId = &id
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().DescribeIndex(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		result = response
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read cls index failed, reason:%s\n", logId, err.Error())
		return err
	}

	res := result.Response
	if res.TopicId != nil {
		_ = d.Set("topic_id", res.TopicId)
	}

	if res.Rule != nil {
		ruleMap := map[string]interface{}{}

		if res.Rule.FullText != nil {
			FullTextMap := map[string]interface{}{}
			if res.Rule.FullText.CaseSensitive != nil {
				FullTextMap["case_sensitive"] = res.Rule.FullText.CaseSensitive
			}

			if res.Rule.FullText.Tokenizer != nil {
				FullTextMap["tokenizer"] = res.Rule.FullText.Tokenizer
			}

			if res.Rule.FullText.ContainZH != nil {
				FullTextMap["contain_z_h"] = res.Rule.FullText.ContainZH
			}

			ruleMap["full_text"] = []interface{}{FullTextMap}
		}

		if res.Rule.KeyValue != nil {
			RuleKeyValueMap := map[string]interface{}{}
			if res.Rule.KeyValue.CaseSensitive != nil {
				RuleKeyValueMap["case_sensitive"] = res.Rule.KeyValue.CaseSensitive
			}

			if res.Rule.KeyValue.KeyValues != nil {
				keyValuesList := []interface{}{}
				for _, keyValueInfo := range res.Rule.KeyValue.KeyValues {
					keyValueInfoMap := map[string]interface{}{}
					if keyValueInfo.Key != nil {
						keyValueInfoMap["key"] = keyValueInfo.Key
					}

					if keyValueInfo.Value != nil {
						valueInfoMap := map[string]interface{}{}
						if keyValueInfo.Value.Type != nil {
							valueInfoMap["type"] = keyValueInfo.Value.Type
						}

						if keyValueInfo.Value.Tokenizer != nil {
							valueInfoMap["tokenizer"] = keyValueInfo.Value.Tokenizer
						}

						if keyValueInfo.Value.SqlFlag != nil {
							valueInfoMap["sql_flag"] = keyValueInfo.Value.SqlFlag
						}

						if keyValueInfo.Value.ContainZH != nil {
							valueInfoMap["contain_z_h"] = keyValueInfo.Value.ContainZH
						}

						keyValueInfoMap["value"] = []interface{}{valueInfoMap}
					}

					keyValuesList = append(keyValuesList, keyValueInfoMap)
				}

				RuleKeyValueMap["key_values"] = keyValuesList
			}

			ruleMap["key_value"] = []interface{}{RuleKeyValueMap}
		}

		if res.Rule.Tag != nil {
			ruleTagMap := map[string]interface{}{
				"case_sensitive": res.Rule.Tag.CaseSensitive,
			}
			if res.Rule.Tag.KeyValues != nil {
				keyValuesList := []interface{}{}
				for _, keyValueInfo := range res.Rule.Tag.KeyValues {
					keyValueInfoMap := map[string]interface{}{
						"key": keyValueInfo.Key,
					}

					if keyValueInfo.Value != nil {
						valueInfoMap := map[string]interface{}{
							"type":        keyValueInfo.Value.Type,
							"tokenizer":   keyValueInfo.Value.Tokenizer,
							"sql_flag":    keyValueInfo.Value.SqlFlag,
							"contain_z_h": keyValueInfo.Value.ContainZH,
						}

						keyValueInfoMap["value"] = []interface{}{valueInfoMap}
					}

					keyValuesList = append(keyValuesList, keyValueInfoMap)
				}

				ruleTagMap["key_values"] = keyValuesList
			}

			ruleMap["tag"] = []interface{}{ruleTagMap}
		}

		if res.Rule.DynamicIndex != nil {
			dynamicIndexMap := map[string]interface{}{}
			if res.Rule.DynamicIndex.Status != nil {
				dynamicIndexMap["status"] = res.Rule.DynamicIndex.Status
			}

			ruleMap["dynamic_index"] = []interface{}{dynamicIndexMap}
		}

		_ = d.Set("rule", []interface{}{ruleMap})
	}

	if res.Status != nil {
		_ = d.Set("status", res.Status)
	}

	if res.IncludeInternalFields != nil {
		_ = d.Set("include_internal_fields", res.IncludeInternalFields)
	}

	if res.MetadataFlag != nil {
		_ = d.Set("metadata_flag", res.MetadataFlag)
	}

	return nil
}

func resourceTencentCloudClsIndexUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_index.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = cls.NewModifyIndexRequest()
		id      = d.Id()
	)

	immutableArgs := []string{"topic_id"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.TopicId = &id
	if dMap, ok := helper.InterfacesHeadMap(d, "rule"); ok {
		ruleInfo := cls.RuleInfo{}
		if fullTextMap, ok := helper.InterfaceToMap(dMap, "full_text"); ok {
			fullTextInfo := cls.FullTextInfo{}
			if v, ok := fullTextMap["case_sensitive"]; ok {
				fullTextInfo.CaseSensitive = helper.Bool(v.(bool))
			}

			if v, ok := fullTextMap["tokenizer"]; ok {
				fullTextInfo.Tokenizer = helper.String(v.(string))
			}

			if v, ok := fullTextMap["contain_z_h"]; ok {
				fullTextInfo.ContainZH = helper.Bool(v.(bool))
			}

			ruleInfo.FullText = &fullTextInfo
		}

		if ruleKeyValueMap, ok := helper.InterfaceToMap(dMap, "key_value"); ok {
			ruleKeyValueInfo := cls.RuleKeyValueInfo{}
			if v, ok := ruleKeyValueMap["case_sensitive"]; ok {
				ruleKeyValueInfo.CaseSensitive = helper.Bool(v.(bool))
			}

			if v, ok := ruleKeyValueMap["key_values"]; ok {
				for _, keyValue := range v.([]interface{}) {
					keyValueInfo := cls.KeyValueInfo{}
					keyValueMap := keyValue.(map[string]interface{})
					if v, ok := keyValueMap["key"]; ok {
						keyValueInfo.Key = helper.String(v.(string))
					}

					if v, ok := keyValueMap["value"]; ok {
						valueMap := v.([]interface{})[0].(map[string]interface{})
						valueInfo := cls.ValueInfo{}
						if v, ok := valueMap["type"]; ok {
							valueInfo.Type = helper.String(v.(string))
						}

						if v, ok := valueMap["tokenizer"]; ok {
							valueInfo.Tokenizer = helper.String(v.(string))
						}

						if v, ok := valueMap["sql_flag"]; ok {
							valueInfo.SqlFlag = helper.Bool(v.(bool))
						}

						if v, ok := valueMap["contain_z_h"]; ok {
							valueInfo.ContainZH = helper.Bool(v.(bool))
						}

						keyValueInfo.Value = &valueInfo
					}

					ruleKeyValueInfo.KeyValues = append(ruleKeyValueInfo.KeyValues, &keyValueInfo)
				}
			}

			ruleInfo.KeyValue = &ruleKeyValueInfo
		}

		if tagMap, ok := helper.InterfaceToMap(dMap, "tag"); ok {
			ruleTagInfo := cls.RuleTagInfo{}
			if v, ok := tagMap["case_sensitive"]; ok {
				ruleTagInfo.CaseSensitive = helper.Bool(v.(bool))
			}

			if v, ok := tagMap["key_values"]; ok {
				for _, keyValue := range v.([]interface{}) {
					keyValueInfo := cls.KeyValueInfo{}
					keyValueMap := keyValue.(map[string]interface{})
					if v, ok := keyValueMap["key"]; ok {
						keyValueInfo.Key = helper.String(v.(string))
					}

					if v, ok := keyValueMap["value"]; ok {
						valueMap := v.([]interface{})[0].(map[string]interface{})
						valueInfo := cls.ValueInfo{}
						if v, ok := valueMap["type"]; ok {
							valueInfo.Type = helper.String(v.(string))
						}

						if v, ok := valueMap["tokenizer"]; ok {
							valueInfo.Tokenizer = helper.String(v.(string))
						}

						if v, ok := valueMap["sql_flag"]; ok {
							valueInfo.SqlFlag = helper.Bool(v.(bool))
						}

						if v, ok := valueMap["contain_z_h"]; ok {
							valueInfo.ContainZH = helper.Bool(v.(bool))
						}

						keyValueInfo.Value = &valueInfo
					}

					ruleTagInfo.KeyValues = append(ruleTagInfo.KeyValues, &keyValueInfo)
				}
			}

			ruleInfo.Tag = &ruleTagInfo
		}

		if dynamicIndexMap, ok := helper.InterfaceToMap(dMap, "dynamic_index"); ok {
			dynamicIndexInfo := cls.DynamicIndex{}
			if v, ok := dynamicIndexMap["status"]; ok {
				dynamicIndexInfo.Status = helper.Bool(v.(bool))
			}

			ruleInfo.DynamicIndex = &dynamicIndexInfo
		}

		request.Rule = &ruleInfo
	}

	if v, ok := d.GetOk("status"); ok {
		request.Status = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("include_internal_fields"); ok {
		request.IncludeInternalFields = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("metadata_flag"); ok {
		request.MetadataFlag = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyIndex(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudClsIndexRead(d, meta)
}

func resourceTencentCloudClsIndexDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cos_shipper.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id      = d.Id()
	)

	if err := service.DeleteClsIndex(ctx, id); err != nil {
		return err
	}

	return nil
}
