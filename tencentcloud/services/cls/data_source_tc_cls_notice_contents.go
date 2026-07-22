package cls

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClsNoticeContents() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClsNoticeContentsRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "过滤条件。支持的过滤器名称：名称（通知内容模板名称）、noticeContentId（通知内容模板ID）。每个请求最多支持 10 个过滤器，每个过滤器值列表最多支持 100 个值。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤字段名称。有效值：名称、noticeContentId。",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "过滤字段值。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"notice_content_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "通知内容模板列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"notice_content_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "通知内容模板 ID。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "注意内容模板名称。",
						},
						"type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "语言类型。 0：中文，1：英文。",
						},
						"flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "模板标志。 0：用户自定义，1：系统内置。",
						},
						"uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "创建者主帐户 ID。",
						},
						"sub_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "创建者/修改者子账户ID。",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "创建时间（Unix 时间戳，以秒为单位）。",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "更新时间（Unix 时间戳，以秒为单位）。",
						},
						"notice_contents": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "注意内容模板详细信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "通道类型。有效值：Email、短信、微信、电话、微信、钉钉、Lark、Http。",
									},
									"trigger_content": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "报警触发通知内容模板。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"title": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "通知内容模板标题。",
												},
												"content": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "通知内容模板主体。",
												},
												"headers": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "请求标头（仅适用于自定义回调通道）。",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"recovery_content": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "报警恢复通知内容模板。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"title": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "通知内容模板标题。",
												},
												"content": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "通知内容模板主体。",
												},
												"headers": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "请求标头（仅适用于自定义回调通道）。",
													Elem: &schema.Schema{
														Type: schema.TypeString,
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

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudClsNoticeContentsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cls_notice_contents.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	filters := buildClsNoticeContentsFilters(d)

	respData, reqErr := service.DescribeClsNoticeContentsByFilter(ctx, filters)
	if reqErr != nil {
		return reqErr
	}

	noticeContentList := flattenClsNoticeContentTemplateList(respData)
	_ = d.Set("notice_content_list", noticeContentList)

	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}

func buildClsNoticeContentsFilters(d *schema.ResourceData) []*cls.Filter {
	if v, ok := d.GetOk("filters"); ok {
		rawList := v.([]interface{})
		filters := make([]*cls.Filter, 0, len(rawList))
		for _, item := range rawList {
			filterMap := item.(map[string]interface{})
			filter := &cls.Filter{}
			if key, ok := filterMap["key"].(string); ok && key != "" {
				filter.Key = helper.String(key)
			}

			if vals, ok := filterMap["values"]; ok {
				valList := vals.([]interface{})
				for _, v := range valList {
					val := v.(string)
					filter.Values = append(filter.Values, &val)
				}
			}

			filters = append(filters, filter)
		}

		return filters
	}

	return nil
}

func flattenClsNoticeContentTemplateList(items []*cls.NoticeContentTemplate) []map[string]interface{} {
	list := make([]map[string]interface{}, 0, len(items))
	for _, tmpl := range items {
		tmplMap := map[string]interface{}{}

		if tmpl.NoticeContentId != nil {
			tmplMap["notice_content_id"] = tmpl.NoticeContentId
		}

		if tmpl.Name != nil {
			tmplMap["name"] = tmpl.Name
		}

		if tmpl.Type != nil {
			tmplMap["type"] = int(*tmpl.Type)
		}

		if tmpl.Flag != nil {
			tmplMap["flag"] = int(*tmpl.Flag)
		}

		if tmpl.Uin != nil {
			tmplMap["uin"] = int(*tmpl.Uin)
		}

		if tmpl.SubUin != nil {
			tmplMap["sub_uin"] = int(*tmpl.SubUin)
		}

		if tmpl.CreateTime != nil {
			tmplMap["create_time"] = int(*tmpl.CreateTime)
		}

		if tmpl.UpdateTime != nil {
			tmplMap["update_time"] = int(*tmpl.UpdateTime)
		}

		if tmpl.NoticeContents != nil {
			tmplMap["notice_contents"] = flattenClsNoticeContentList(tmpl.NoticeContents)
		}

		list = append(list, tmplMap)
	}

	return list
}

func flattenClsNoticeContentList(items []*cls.NoticeContent) []map[string]interface{} {
	list := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		itemMap := map[string]interface{}{}

		if item.Type != nil {
			itemMap["type"] = item.Type
		}

		if item.TriggerContent != nil {
			itemMap["trigger_content"] = []map[string]interface{}{flattenClsNoticeContentInfo(item.TriggerContent)}
		}

		if item.RecoveryContent != nil {
			itemMap["recovery_content"] = []map[string]interface{}{flattenClsNoticeContentInfo(item.RecoveryContent)}
		}

		list = append(list, itemMap)
	}

	return list
}

func flattenClsNoticeContentInfo(info *cls.NoticeContentInfo) map[string]interface{} {
	infoMap := map[string]interface{}{}

	if info.Title != nil {
		infoMap["title"] = info.Title
	}

	if info.Content != nil {
		infoMap["content"] = info.Content
	}

	if info.Headers != nil {
		headers := make([]string, 0, len(info.Headers))
		for _, h := range info.Headers {
			if h != nil {
				headers = append(headers, *h)
			}
		}

		infoMap["headers"] = headers
	}

	return infoMap
}
