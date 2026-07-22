package cls

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsConfigCreate,
		Read:   resourceTencentCloudClsConfigRead,
		Delete: resourceTencentCloudClsConfigDelete,
		Update: resourceTencentCloudClsConfigUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "集合配置名称。",
			},
			"output": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "采集配置的日志主题ID（TopicId）。",
			},
			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "包含文件名的日志收集路径。需要收集文件。",
			},
			"log_type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "需要收集的日志类型。有效值： json_log：JSON格式日志； delimiter_log：以分隔格式记录； imalist_log：极简日志； multiline_log：以多行格式记录；" +
					"fullregex_log: log in full regex format. Default value: minimalist_log.",
			},
			"extract_rule": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "提取规则。如果设置了ExtractRule，则必须设置LogType。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "时间字段键名称。 time_key 和 time_format 必须成对出现。",
						},
						"time_format": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "时间字段格式。更多信息请参见C语言strftime函数时间格式说明的输出参数。",
						},
						"delimiter": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "分隔日志的分隔符，仅当log_type为delimiter_log时有效。",
						},
						"log_regex": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "全日志匹配规则，仅当log_type为fullregex_log时有效。",
						},
						"begin_regex": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "首行匹配规则，仅当log_type为multiline_log或fullregex_log时有效。",
						},
						"keys": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "每个提取字段的键名称。空键表示放弃该字段。该参数仅当log_type为delimiter_log时有效。 json_log 日志使用 JSON 本身的密钥。",
						},
						"filter_key_regex": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "记录要过滤的键和相应的正则表达式。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "要过滤的日志键。",
									},
									"regex": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "key对应的过滤规则正则表达式。",
									},
								},
							},
						},
						"un_match_up_load_switch": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "是否上传解析失败的日志。有效值：true：是；假：没有。",
						},
						"un_match_log_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "不匹配的日志密钥。当 UnMatchUpLoadSwitch 为 true 时必需。",
						},
						"backtracking": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "增量采集模式下回滚的数据大小。默认值：-1（完整集合）。",
						},
						"is_gbk": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "GBK编码。默认 0。 注意： - 目前，当值为 0 时，表示 UTF-8 编码。",
						},
						"json_standard": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "标准 json。默认 0。",
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "系统日志协议、tcp 或 udp。该值可以是 tcp 或 udp。仅当LogType为service_syslog时有效。其他类型无需填写。",
						},
						"address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "syslog系统日志采集指定采集器监听的地址和端口。该参数仅当LogType为service_syslog时有效。其他类型无需填写。",
						},
						"parse_protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "解析协议。该参数仅当LogType为service_syslog时有效。其他类型无需填写。",
						},
						"metadata_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "元数据类型。 0：不使用元数据信息； 1：使用机器组元数据； 2：使用用户定义的元数据； 3：使用采集配置路径。注意：COS 导入不支持该字段。",
						},
						"path_regex": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "元数据路径正则表达式。",
						},
						"meta_tags": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "元数据标签。注： - MetadataType 为 2 时必填。 - COS 导入不支持该字段。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "标签键。",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "标签值。",
									},
								},
							},
						},
					},
				},
			},
			"exclude_paths": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "收集路径阻止列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "类型。有效值：文件、路径。",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "具体内容对应Type。",
						},
					},
				},
			},
			"user_define_rule": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "自定义采集规则，是序列化的JSON字符串。 LogType为user_define_log时必填。",
			},
		},
	}
}

func resourceTencentCloudClsConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_config.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = cls.NewCreateConfigRequest()
		response *cls.CreateConfigResponse
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}
	if v, ok := d.GetOk("output"); ok {
		request.Output = helper.String(v.(string))
	}
	if v, ok := d.GetOk("path"); ok {
		request.Path = helper.String(v.(string))
	}
	if v, ok := d.GetOk("log_type"); ok {
		request.LogType = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "extract_rule"); ok {
		extractRule := cls.ExtractRuleInfo{}
		if v, ok := dMap["time_key"].(string); ok && v != "" {
			extractRule.TimeKey = helper.String(v)
		}
		if v, ok := dMap["time_format"].(string); ok && v != "" {
			extractRule.TimeFormat = helper.String(v)
		}
		if v, ok := dMap["delimiter"].(string); ok && v != "" {
			extractRule.Delimiter = helper.String(v)
		}
		if v, ok := dMap["log_regex"].(string); ok && v != "" {
			extractRule.LogRegex = helper.String(v)
		}
		if v, ok := dMap["begin_regex"].(string); ok && v != "" {
			extractRule.BeginRegex = helper.String(v)
		}
		if v, ok := dMap["keys"]; ok {
			keys := v.(*schema.Set).List()
			for _, key := range keys {
				extractRule.Keys = append(extractRule.Keys, helper.String(key.(string)))
			}
		}
		if v, ok := dMap["filter_key_regex"]; ok {
			keyRegexs := make([]*cls.KeyRegexInfo, 0, 10)
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				keyRegex := cls.KeyRegexInfo{}
				if v, ok := dMap["key"]; ok {
					keyRegex.Key = helper.String(v.(string))
				}
				if v, ok := dMap["regex"]; ok {
					keyRegex.Regex = helper.String(v.(string))
				}
				keyRegexs = append(keyRegexs, &keyRegex)
			}
			extractRule.FilterKeyRegex = keyRegexs
		}
		if v, ok := dMap["un_match_up_load_switch"]; ok {
			extractRule.UnMatchUpLoadSwitch = helper.Bool(v.(bool))
		}
		if v, ok := dMap["un_match_log_key"].(string); ok && v != "" {
			extractRule.UnMatchLogKey = helper.String(v)
		}
		if v, ok := dMap["backtracking"]; ok {
			extractRule.Backtracking = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["is_gbk"]; ok {
			extractRule.IsGBK = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["json_standard"]; ok {
			extractRule.JsonStandard = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["protocol"].(string); ok && v != "" {
			extractRule.Protocol = helper.String(v)
		}
		if v, ok := dMap["address"].(string); ok && v != "" {
			extractRule.Address = helper.String(v)
		}
		if v, ok := dMap["parse_protocol"].(string); ok && v != "" {
			extractRule.ParseProtocol = helper.String(v)
		}
		if v, ok := dMap["metadata_type"]; ok {
			extractRule.MetadataType = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["path_regex"].(string); ok && v != "" {
			extractRule.PathRegex = helper.String(v)
		}
		if v, ok := dMap["meta_tags"]; ok {
			for _, item := range v.([]interface{}) {
				metaTagsMap := item.(map[string]interface{})
				metaTagInfo := cls.MetaTagInfo{}
				if v, ok := metaTagsMap["key"]; ok {
					metaTagInfo.Key = helper.String(v.(string))
				}
				if v, ok := metaTagsMap["value"]; ok {
					metaTagInfo.Value = helper.String(v.(string))
				}
				extractRule.MetaTags = append(extractRule.MetaTags, &metaTagInfo)
			}
		}
		request.ExtractRule = &extractRule
	}
	if v, ok := d.GetOk("exclude_paths"); ok {
		excludePaths := make([]*cls.ExcludePathInfo, 0, 10)
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			excludePath := cls.ExcludePathInfo{}
			if v, ok := dMap["type"].(string); ok && v != "" {
				excludePath.Type = helper.String(v)
			}
			if v, ok := dMap["value"].(string); ok && v != "" {
				excludePath.Value = helper.String(v)
			}
			excludePaths = append(excludePaths, &excludePath)
		}
		request.ExcludePaths = excludePaths
	}
	if v, ok := d.GetOk("user_define_rule"); ok {
		request.UserDefineRule = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls config extra failed, reason:%+v", logId, err)
		return err
	}

	id := *response.Response.ConfigId
	d.SetId(id)
	return resourceTencentCloudClsConfigRead(d, meta)
}

func resourceTencentCloudClsConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	configId := d.Id()

	config, err := service.DescribeClsConfigById(ctx, configId)
	if err != nil {
		return err
	}

	if config == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ClsConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if config.Name != nil {
		_ = d.Set("name", config.Name)
	}

	if config.Output != nil {
		_ = d.Set("output", config.Output)
	}

	if config.Path != nil {
		_ = d.Set("path", config.Path)
	}

	if config.LogType != nil {
		_ = d.Set("log_type", config.LogType)
	}

	if config.ExtractRule != nil {
		extractRuleMap := map[string]interface{}{}

		if config.ExtractRule.TimeKey != nil {
			extractRuleMap["time_key"] = config.ExtractRule.TimeKey
		}

		if config.ExtractRule.TimeFormat != nil {
			extractRuleMap["time_format"] = config.ExtractRule.TimeFormat
		}

		if config.ExtractRule.Delimiter != nil {
			extractRuleMap["delimiter"] = config.ExtractRule.Delimiter
		}

		if config.ExtractRule.LogRegex != nil {
			extractRuleMap["log_regex"] = config.ExtractRule.LogRegex
		}

		if config.ExtractRule.BeginRegex != nil {
			extractRuleMap["begin_regex"] = config.ExtractRule.BeginRegex
		}

		if config.ExtractRule.Keys != nil {
			extractRuleMap["keys"] = config.ExtractRule.Keys
		}

		if config.ExtractRule.FilterKeyRegex != nil {
			filterKeyRegexList := []interface{}{}
			for _, filterKeyRegex := range config.ExtractRule.FilterKeyRegex {
				filterKeyRegexMap := map[string]interface{}{}

				if filterKeyRegex.Key != nil {
					filterKeyRegexMap["key"] = filterKeyRegex.Key
				}

				if filterKeyRegex.Regex != nil {
					filterKeyRegexMap["regex"] = filterKeyRegex.Regex
				}

				filterKeyRegexList = append(filterKeyRegexList, filterKeyRegexMap)
			}

			extractRuleMap["filter_key_regex"] = filterKeyRegexList
		}

		if config.ExtractRule.UnMatchUpLoadSwitch != nil {
			extractRuleMap["un_match_up_load_switch"] = config.ExtractRule.UnMatchUpLoadSwitch
		}

		if config.ExtractRule.UnMatchLogKey != nil {
			extractRuleMap["un_match_log_key"] = config.ExtractRule.UnMatchLogKey
		}

		if config.ExtractRule.Backtracking != nil {
			extractRuleMap["backtracking"] = config.ExtractRule.Backtracking
		}

		if config.ExtractRule.IsGBK != nil {
			extractRuleMap["is_gbk"] = config.ExtractRule.IsGBK
		}

		if config.ExtractRule.JsonStandard != nil {
			extractRuleMap["json_standard"] = config.ExtractRule.JsonStandard
		}

		if config.ExtractRule.Protocol != nil {
			extractRuleMap["protocol"] = config.ExtractRule.Protocol
		}

		if config.ExtractRule.Address != nil {
			extractRuleMap["address"] = config.ExtractRule.Address
		}

		if config.ExtractRule.ParseProtocol != nil {
			extractRuleMap["parse_protocol"] = config.ExtractRule.ParseProtocol
		}

		if config.ExtractRule.MetadataType != nil {
			extractRuleMap["metadata_type"] = config.ExtractRule.MetadataType
		}

		if config.ExtractRule.PathRegex != nil {
			extractRuleMap["path_regex"] = config.ExtractRule.PathRegex
		}

		if config.ExtractRule.MetaTags != nil {
			metaTagsList := []interface{}{}
			for _, metaTags := range config.ExtractRule.MetaTags {
				metaTagsMap := map[string]interface{}{}

				if metaTags.Key != nil {
					metaTagsMap["key"] = metaTags.Key
				}

				if metaTags.Value != nil {
					metaTagsMap["value"] = metaTags.Value
				}

				metaTagsList = append(metaTagsList, metaTagsMap)
			}

			extractRuleMap["meta_tags"] = metaTagsList
		}

		_ = d.Set("extract_rule", []interface{}{extractRuleMap})
	}

	if config.ExcludePaths != nil {
		excludePathsList := []interface{}{}
		for _, excludePath := range config.ExcludePaths {
			excludePathsMap := map[string]interface{}{}

			if excludePath.Type != nil {
				excludePathsMap["type"] = excludePath.Type
			}

			if excludePath.Value != nil {
				excludePathsMap["value"] = excludePath.Value
			}

			excludePathsList = append(excludePathsList, excludePathsMap)
		}

		_ = d.Set("exclude_paths", excludePathsList)
	}

	if config.UserDefineRule != nil {
		_ = d.Set("user_define_rule", config.UserDefineRule)
	}

	return nil
}

func resourceTencentCloudClsConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_config.update")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := cls.NewModifyConfigRequest()

	request.ConfigId = helper.String(d.Id())

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}
	if d.HasChange("output") {
		if v, ok := d.GetOk("output"); ok {
			request.Output = helper.String(v.(string))
		}
	}
	if d.HasChange("path") {
		if v, ok := d.GetOk("path"); ok {
			request.Path = helper.String(v.(string))
		}
	}
	if d.HasChange("log_type") || d.HasChange("extract_rule") {
		if v, ok := d.GetOk("log_type"); ok {
			request.LogType = helper.String(v.(string))
		}
		if dMap, ok := helper.InterfacesHeadMap(d, "extract_rule"); ok {
			extractRule := cls.ExtractRuleInfo{}
			if v, ok := dMap["time_key"].(string); ok && v != "" {
				extractRule.TimeKey = helper.String(v)
			}
			if v, ok := dMap["time_format"].(string); ok && v != "" {
				extractRule.TimeFormat = helper.String(v)
			}
			if v, ok := dMap["delimiter"].(string); ok && v != "" {
				extractRule.Delimiter = helper.String(v)
			}
			if v, ok := dMap["log_regex"].(string); ok && v != "" {
				extractRule.LogRegex = helper.String(v)
			}
			if v, ok := dMap["begin_regex"].(string); ok && v != "" {
				extractRule.BeginRegex = helper.String(v)
			}
			if v, ok := dMap["keys"]; ok {
				keys := v.(*schema.Set).List()
				for _, key := range keys {
					extractRule.Keys = append(extractRule.Keys, helper.String(key.(string)))
				}
			}
			if v, ok := dMap["filter_key_regex"]; ok {
				keyRegexs := make([]*cls.KeyRegexInfo, 0, 10)
				for _, item := range v.([]interface{}) {
					dMap := item.(map[string]interface{})
					keyRegex := cls.KeyRegexInfo{}
					if v, ok := dMap["key"]; ok {
						keyRegex.Key = helper.String(v.(string))
					}
					if v, ok := dMap["regex"]; ok {
						keyRegex.Regex = helper.String(v.(string))
					}
					keyRegexs = append(keyRegexs, &keyRegex)
				}
				extractRule.FilterKeyRegex = keyRegexs
			}
			if v, ok := dMap["un_match_up_load_switch"]; ok {
				extractRule.UnMatchUpLoadSwitch = helper.Bool(v.(bool))
			}
			if v, ok := dMap["un_match_log_key"].(string); ok && v != "" {
				extractRule.UnMatchLogKey = helper.String(v)
			}
			if v, ok := dMap["backtracking"]; ok {
				extractRule.Backtracking = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["is_gbk"]; ok {
				extractRule.IsGBK = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["json_standard"]; ok {
				extractRule.JsonStandard = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["protocol"].(string); ok && v != "" {
				extractRule.Protocol = helper.String(v)
			}
			if v, ok := dMap["address"].(string); ok && v != "" {
				extractRule.Address = helper.String(v)
			}
			if v, ok := dMap["parse_protocol"].(string); ok && v != "" {
				extractRule.ParseProtocol = helper.String(v)
			}
			if v, ok := dMap["metadata_type"]; ok {
				extractRule.MetadataType = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["path_regex"].(string); ok && v != "" {
				extractRule.PathRegex = helper.String(v)
			}
			if v, ok := dMap["meta_tags"]; ok {
				for _, item := range v.([]interface{}) {
					metaTagsMap := item.(map[string]interface{})
					metaTagInfo := cls.MetaTagInfo{}
					if v, ok := metaTagsMap["key"]; ok {
						metaTagInfo.Key = helper.String(v.(string))
					}
					if v, ok := metaTagsMap["value"]; ok {
						metaTagInfo.Value = helper.String(v.(string))
					}
					extractRule.MetaTags = append(extractRule.MetaTags, &metaTagInfo)
				}
			}
			request.ExtractRule = &extractRule
		}
	}
	if d.HasChange("exclude_paths") {
		if v, ok := d.GetOk("exclude_paths"); ok {
			excludePaths := make([]*cls.ExcludePathInfo, 0, 10)
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				excludePath := cls.ExcludePathInfo{}
				if v, ok := dMap["type"].(string); ok && v != "" {
					excludePath.Type = helper.String(v)
				}
				if v, ok := dMap["value"].(string); ok && v != "" {
					excludePath.Value = helper.String(v)
				}
				excludePaths = append(excludePaths, &excludePath)
			}
			request.ExcludePaths = excludePaths
		}
	}

	if d.HasChange("user_define_rule") {
		if v, ok := d.GetOk("user_define_rule"); ok {
			request.UserDefineRule = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyConfig(request)
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

	return resourceTencentCloudClsConfigRead(d, meta)
}

func resourceTencentCloudClsConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_config.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	id := d.Id()

	if err := service.DeleteClsConfig(ctx, id); err != nil {
		return err
	}

	return nil
}
