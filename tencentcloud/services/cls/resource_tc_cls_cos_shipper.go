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

func ResourceTencentCloudClsCosShipper() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsCosShipperCreate,
		Read:   resourceTencentCloudClsCosShipperRead,
		Delete: resourceTencentCloudClsCosShipperDelete,
		Update: resourceTencentCloudClsCosShipperUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"topic_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "待创建的发货规则所属日志主题ID。",
			},
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "要创建的发货规则中的目标存储桶。",
			},
			"prefix": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "要创建的发货规则中的发货目录的前缀。",
			},
			"shipper_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "运输规则名称。",
			},
			"interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "运送时间间隔（以秒为单位）。默认值：300。取值范围：300~900。",
			},
			"max_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "要传送的文件的最大大小（以 MB 为单位）。默认值：256。取值范围：100~256。",
			},
			"filter_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Description: "已发送日志的过滤规则。只有符合规则的日志才能被发送。所有规则都是AND关系，最多可以添加5条规则。" +
					"If the array is empty, no filtering will be performed, and all logs will be shipped.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤规则键。",
						},
						"regex": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤规则。",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤规则值。",
						},
					},
				},
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "发送日志的分区规则，可以用strftime时间格式表示。",
			},
			"compress": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"format": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "压缩格式。有效值：gzip、lzop、none（无压缩）。",
						},
					},
				},
				Description: "发送日志的压缩配置。",
			},
			"content": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "发送日志内容的格式配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"format": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "内容格式。有效值：json、csv、parquet。",
						},
						"csv": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"print_key": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "是否在 CSV 文件的第一行打印密钥。",
									},
									"keys": {
										Type:        schema.TypeSet,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "键名。注意：该字段可能返回null，表示取不到有效值。",
									},
									"delimiter": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "字段分隔符。",
									},
									"escape_char": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "字段分隔符。",
									},
									"non_existing_field": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "用于填充不存在字段的内容。",
									},
								},
							},
							Description: "CSV格式内容描述。注意：该字段可能返回null，表示取不到有效值。",
						},
						"json": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_tag": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "启用标志。",
									},
									"meta_fields": {
										Type:        schema.TypeSet,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "元数据信息列表\n注意：该字段可能返回null，表示取不到有效值。",
									},
								},
							},
							Description: "JSON格式内容描述。注意：该字段可能返回null，表示取不到有效值。",
						},
						"parquet": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"parquet_key_info": {
										Type:        schema.TypeList,
										Required:    true,
										MinItems:    1,
										Description: "Parquet 列定义数组。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key_name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Parquet 文件中的列名称。",
												},
												"key_type": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"string", "boolean", "int32", "int64", "float", "double"}),
													Description: "列的数据类型。有效值：字符串、布尔值、int32、int64、float、double。",
												},
												"key_non_existing_field": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "当字段不存在或解析失败时要分配的值。",
												},
											},
										},
									},
								},
							},
							Description: "Parquet格式内容描述。注意：该字段可能返回null，表示取不到有效值。",
						},
					},
				},
			},
			"filename_mode": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "命名运输文件。有效值：0（随机数）； 1（按运输时间）。默认值：0。",
			},
			"start_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "数据传送的开始时间，不能早于日志主题的生命周期开始时间。如果不指定该参数，则默认为创建数据传送任务时的时间。",
			},
			"end_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "数据传输的结束时间，不能设置为将来的时间。如果不指定该参数，则表示持续发送数据。",
			},
			"storage_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "COS 桶存储类型。支持：STANDARD_IA、ARCHIVE、DEEP_ARCHIVE、STANDARD、MAZ_STANDARD、MAZ_STANDARD_IA、INTELLIGENT_TIERING。",
			},
		},
	}
}

func resourceTencentCloudClsCosShipperCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cos_shipper.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = cls.NewCreateShipperRequest()
		response *cls.CreateShipperResponse
	)

	if v, ok := d.GetOk("topic_id"); ok {
		request.TopicId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bucket"); ok {
		request.Bucket = helper.String(v.(string))
	}

	if v, ok := d.GetOk("prefix"); ok {
		request.Prefix = helper.String(v.(string))
	}

	if v, ok := d.GetOk("shipper_name"); ok {
		request.ShipperName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("interval"); ok {
		request.Interval = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("max_size"); ok {
		request.MaxSize = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("filter_rules"); ok {
		filterRules := make([]*cls.FilterRuleInfo, 0, 10)
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			filterRule := cls.FilterRuleInfo{}
			if v, ok := dMap["key"]; ok {
				filterRule.Key = helper.String(v.(string))
			}
			if v, ok := dMap["regex"]; ok {
				filterRule.Regex = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				filterRule.Value = helper.String(v.(string))
			}
			filterRules = append(filterRules, &filterRule)
		}
		request.FilterRules = filterRules
	}

	if v, ok := d.GetOk("partition"); ok {
		request.Partition = helper.String(v.(string))
	}

	if v, ok := d.GetOk("compress"); ok {
		compresses := make([]*cls.CompressInfo, 0, 10)
		if len(v.([]interface{})) != 1 {
			return fmt.Errorf("need only one compress.")
		}
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			compress := cls.CompressInfo{}
			if v, ok := dMap["format"]; ok {
				compress.Format = helper.String(v.(string))
			}
			compresses = append(compresses, &compress)
		}
		request.Compress = compresses[0]
	}

	if v, ok := d.GetOk("content"); ok {
		contents := make([]*cls.ContentInfo, 0, 10)
		if len(v.([]interface{})) != 1 {
			return fmt.Errorf("need only one content.")
		}
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			content := cls.ContentInfo{}
			if v, ok := dMap["format"]; ok {
				content.Format = helper.String(v.(string))
			}
			if v, ok := dMap["csv"]; ok {
				if len(v.([]interface{})) == 1 {
					csv := v.([]interface{})[0].(map[string]interface{})
					csvInfo := cls.CsvInfo{}
					csvInfo.PrintKey = helper.Bool(csv["print_key"].(bool))
					keys := csv["keys"].(*schema.Set).List()
					for _, key := range keys {
						csvInfo.Keys = append(csvInfo.Keys, helper.String(key.(string)))
					}
					csvInfo.Delimiter = helper.String(csv["delimiter"].(string))
					csvInfo.EscapeChar = helper.String(csv["escape_char"].(string))
					csvInfo.NonExistingField = helper.String(csv["non_existing_field"].(string))
					content.Csv = &csvInfo
				}
			}
			if v, ok := dMap["json"]; ok {
				if len(v.([]interface{})) == 1 {

					json := v.([]interface{})[0].(map[string]interface{})
					jsonInfo := cls.JsonInfo{}
					jsonInfo.EnableTag = helper.Bool(json["enable_tag"].(bool))
					metaFields := json["meta_fields"].(*schema.Set).List()
					for _, metaField := range metaFields {
						jsonInfo.MetaFields = append(jsonInfo.MetaFields, helper.String(metaField.(string)))
					}
					content.Json = &jsonInfo
				}
			}
			if v, ok := dMap["parquet"]; ok {
				if len(v.([]interface{})) == 1 {
					parquet := v.([]interface{})[0].(map[string]interface{})
					parquetInfo := cls.ParquetInfo{}

					if keyInfos, ok := parquet["parquet_key_info"]; ok {
						parquetKeyInfoList := keyInfos.([]interface{})
						parquetInfo.ParquetKeyInfo = make([]*cls.ParquetKeyInfo, 0, len(parquetKeyInfoList))

						for _, keyInfo := range parquetKeyInfoList {
							keyInfoMap := keyInfo.(map[string]interface{})
							parquetKeyInfo := &cls.ParquetKeyInfo{
								KeyName: helper.String(keyInfoMap["key_name"].(string)),
								KeyType: helper.String(keyInfoMap["key_type"].(string)),
							}
							if v, ok := keyInfoMap["key_non_existing_field"]; ok {
								parquetKeyInfo.KeyNonExistingField = helper.String(v.(string))
							}
							parquetInfo.ParquetKeyInfo = append(parquetInfo.ParquetKeyInfo, parquetKeyInfo)
						}
					}
					content.Parquet = &parquetInfo
				}
			}
			contents = append(contents, &content)
		}
		request.Content = contents[0]
	}

	if v, ok := d.GetOkExists("filename_mode"); ok {
		request.FilenameMode = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("start_time"); ok {
		request.StartTime = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOkExists("end_time"); ok {
		request.EndTime = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("storage_type"); ok {
		request.StorageType = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateShipper(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create cls cos shipper failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls cos shipper failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.ShipperId == nil {
		return fmt.Errorf("ShipperId is nil.")
	}

	id := *response.Response.ShipperId
	d.SetId(id)
	return resourceTencentCloudClsCosShipperRead(d, meta)
}

func resourceTencentCloudClsCosShipperRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cos_shipper.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	id := d.Id()

	shipper, err := service.DescribeClsCosShipperById(ctx, id)

	if err != nil {
		return err
	}

	if shipper == nil {
		d.SetId("")
		return fmt.Errorf("resource `Shipper` %s does not exist", id)
	}

	_ = d.Set("topic_id", shipper.TopicId)
	_ = d.Set("bucket", shipper.Bucket)
	_ = d.Set("prefix", shipper.Prefix)
	_ = d.Set("shipper_name", shipper.ShipperName)
	if shipper.Interval != nil {
		_ = d.Set("interval", shipper.Interval)
	}
	if shipper.MaxSize != nil {
		_ = d.Set("max_size", shipper.MaxSize)
	}

	if shipper.FilterRules != nil {
		filterRules := make([]interface{}, 0, 100)
		for _, v := range shipper.FilterRules {
			filterRule := map[string]interface{}{
				"key":   v.Key,
				"regex": v.Regex,
				"value": v.Value,
			}
			filterRules = append(filterRules, filterRule)
		}
		_ = d.Set("filter_rules", filterRules)
	}

	if shipper.Partition != nil {
		_ = d.Set("partition", shipper.Partition)
	}

	if shipper.Compress != nil {
		compress := map[string]interface{}{
			"format": shipper.Compress.Format,
		}
		_ = d.Set("compress", []interface{}{compress})
	}

	if shipper.Content != nil {
		content := map[string]interface{}{
			"format": shipper.Content.Format,
		}
		if shipper.Content.Csv != nil {
			csv := map[string]interface{}{
				"print_key":          shipper.Content.Csv.PrintKey,
				"keys":               shipper.Content.Csv.Keys,
				"delimiter":          shipper.Content.Csv.Delimiter,
				"escape_char":        shipper.Content.Csv.EscapeChar,
				"non_existing_field": shipper.Content.Csv.NonExistingField,
			}
			content["csv"] = []interface{}{csv}
		}
		if shipper.Content.Json != nil {
			json := map[string]interface{}{
				"enable_tag":  shipper.Content.Json.EnableTag,
				"meta_fields": shipper.Content.Json.MetaFields,
			}
			content["json"] = []interface{}{json}
		}
		if shipper.Content.Parquet != nil {
			parquetKeyInfoList := make([]interface{}, 0, len(shipper.Content.Parquet.ParquetKeyInfo))

			for _, keyInfo := range shipper.Content.Parquet.ParquetKeyInfo {
				parquetKeyInfoMap := map[string]interface{}{
					"key_name": keyInfo.KeyName,
					"key_type": keyInfo.KeyType,
				}
				if keyInfo.KeyNonExistingField != nil {
					parquetKeyInfoMap["key_non_existing_field"] = keyInfo.KeyNonExistingField
				}
				parquetKeyInfoList = append(parquetKeyInfoList, parquetKeyInfoMap)
			}

			parquet := map[string]interface{}{
				"parquet_key_info": parquetKeyInfoList,
			}
			content["parquet"] = []interface{}{parquet}
		}
		_ = d.Set("content", []interface{}{content})
	}

	if shipper.FilenameMode != nil {
		_ = d.Set("filename_mode", shipper.FilenameMode)
	}

	if shipper.StartTime != nil {
		_ = d.Set("start_time", shipper.StartTime)
	}

	if shipper.EndTime != nil {
		_ = d.Set("end_time", shipper.EndTime)
	}

	if shipper.StorageType != nil {
		_ = d.Set("storage_type", shipper.StorageType)
	}

	return nil
}

func resourceTencentCloudClsCosShipperUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cos_shipper.update")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := cls.NewModifyShipperRequest()

	immutableArgs := []string{"start_time", "end_time"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.ShipperId = helper.String(d.Id())

	if d.HasChange("bucket") {
		if v, ok := d.GetOk("bucket"); ok {
			request.Bucket = helper.String(v.(string))
		}
	}

	if d.HasChange("prefix") {
		if v, ok := d.GetOk("prefix"); ok {
			request.Prefix = helper.String(v.(string))
		}
	}

	if d.HasChange("shipper_name") {
		if v, ok := d.GetOk("shipper_name"); ok {
			request.ShipperName = helper.String(v.(string))
		}
	}

	if d.HasChange("interval") {
		if v, ok := d.GetOkExists("interval"); ok {
			request.Interval = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("max_size") {
		if v, ok := d.GetOkExists("max_size"); ok {
			request.MaxSize = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("filter_rules") {
		if v, ok := d.GetOk("filter_rules"); ok {
			filterRules := make([]*cls.FilterRuleInfo, 0, 10)
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				filterRule := cls.FilterRuleInfo{}
				if v, ok := dMap["key"]; ok {
					filterRule.Key = helper.String(v.(string))
				}
				if v, ok := dMap["regex"]; ok {
					filterRule.Regex = helper.String(v.(string))
				}
				if v, ok := dMap["value"]; ok {
					filterRule.Value = helper.String(v.(string))
				}
				filterRules = append(filterRules, &filterRule)
			}
			request.FilterRules = filterRules
		}
	}

	if d.HasChange("partition") {
		if v, ok := d.GetOk("partition"); ok {
			request.Partition = helper.String(v.(string))
		}
	}

	if d.HasChange("compress") {
		if v, ok := d.GetOk("compress"); ok {
			compresses := make([]*cls.CompressInfo, 0, 10)
			if len(v.([]interface{})) != 1 {
				return fmt.Errorf("need only one compress.")
			}
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				compress := cls.CompressInfo{}
				if v, ok := dMap["format"]; ok {
					compress.Format = helper.String(v.(string))
				}
				compresses = append(compresses, &compress)
			}
			request.Compress = compresses[0]
		}
	}

	if d.HasChange("content") {
		if v, ok := d.GetOk("content"); ok {
			contents := make([]*cls.ContentInfo, 0, 10)
			if len(v.([]interface{})) != 1 {
				return fmt.Errorf("need only one content.")
			}
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				content := cls.ContentInfo{}
				if v, ok := dMap["format"]; ok {
					content.Format = helper.String(v.(string))
				}
				if v, ok := dMap["csv"]; ok {
					if len(v.([]interface{})) == 1 {
						csv := v.([]interface{})[0].(map[string]interface{})
						csvInfo := cls.CsvInfo{}
						csvInfo.PrintKey = helper.Bool(csv["print_key"].(bool))
						keys := csv["keys"].(*schema.Set).List()
						for _, key := range keys {
							csvInfo.Keys = append(csvInfo.Keys, helper.String(key.(string)))
						}
						csvInfo.Delimiter = helper.String(csv["delimiter"].(string))
						csvInfo.EscapeChar = helper.String(csv["escape_char"].(string))
						csvInfo.NonExistingField = helper.String(csv["non_existing_field"].(string))
						content.Csv = &csvInfo
					}
				}
				if v, ok := dMap["json"]; ok {
					if len(v.([]interface{})) == 1 {

						json := v.([]interface{})[0].(map[string]interface{})
						jsonInfo := cls.JsonInfo{}
						jsonInfo.EnableTag = helper.Bool(json["enable_tag"].(bool))
						metaFields := json["meta_fields"].(*schema.Set).List()
						for _, metaField := range metaFields {
							jsonInfo.MetaFields = append(jsonInfo.MetaFields, helper.String(metaField.(string)))
						}
						content.Json = &jsonInfo
					}
				}
				if v, ok := dMap["parquet"]; ok {
					if len(v.([]interface{})) == 1 {
						parquet := v.([]interface{})[0].(map[string]interface{})
						parquetInfo := cls.ParquetInfo{}

						if keyInfos, ok := parquet["parquet_key_info"]; ok {
							parquetKeyInfoList := keyInfos.([]interface{})
							parquetInfo.ParquetKeyInfo = make([]*cls.ParquetKeyInfo, 0, len(parquetKeyInfoList))

							for _, keyInfo := range parquetKeyInfoList {
								keyInfoMap := keyInfo.(map[string]interface{})
								parquetKeyInfo := &cls.ParquetKeyInfo{
									KeyName: helper.String(keyInfoMap["key_name"].(string)),
									KeyType: helper.String(keyInfoMap["key_type"].(string)),
								}
								if v, ok := keyInfoMap["key_non_existing_field"]; ok {
									parquetKeyInfo.KeyNonExistingField = helper.String(v.(string))
								}
								parquetInfo.ParquetKeyInfo = append(parquetInfo.ParquetKeyInfo, parquetKeyInfo)
							}
						}
						content.Parquet = &parquetInfo
					}
				}
				contents = append(contents, &content)
			}
			request.Content = contents[0]
		}
	}

	if d.HasChange("filename_mode") {
		if v, ok := d.GetOkExists("filename_mode"); ok {
			request.FilenameMode = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("storage_type") {
		if v, ok := d.GetOk("storage_type"); ok {
			request.StorageType = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyShipper(request)
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

	return resourceTencentCloudClsCosShipperRead(d, meta)
}

func resourceTencentCloudClsCosShipperDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_cos_shipper.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	id := d.Id()

	if err := service.DeleteClsCosShipper(ctx, id); err != nil {
		return err
	}

	return nil
}
