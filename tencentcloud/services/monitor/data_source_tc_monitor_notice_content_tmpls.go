package monitor

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitorv20230616 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20230616"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMonitorNoticeContentTmpls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMonitorNoticeContentTmplsRead,
		Schema: map[string]*schema.Schema{
			"tmpl_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "模板 ID list for query。",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tmpl_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "模板名称 for query。",
			},
			"notice_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Notice 模板 ID for query。",
			},
			"tmpl_language": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Template 语言 for query. 有效值：`en`，`zh`。",
			},
			"monitor_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "监控类型 for query. Valid 值: `MT_QCE`。",
			},

			"notice_content_tmpl_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Notification 内容模板 list。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tmpl_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "模板 ID",
						},
						"tmpl_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "模板名称",
						},
						"monitor_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "监控类型",
						},
						"tmpl_language": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template 语言",
						},
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建者 uin。",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "创建时间。",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "更新时间。",
						},
						"tmpl_contents_json": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template 内容 in JSON 格式",
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

func dataSourceTencentCloudMonitorNoticeContentTmplsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_notice_content_tmpls.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})

	if v, ok := d.GetOk("tmpl_ids"); ok {
		tmplIDsSet := v.(*schema.Set).List()
		tmplIDs := make([]*string, 0, len(tmplIDsSet))
		for _, item := range tmplIDsSet {
			tmplID := item.(string)
			tmplIDs = append(tmplIDs, helper.String(tmplID))
		}
		paramMap["TmplIDs"] = tmplIDs
	}

	if v, ok := d.GetOk("tmpl_name"); ok {
		paramMap["TmplName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("notice_id"); ok {
		paramMap["NoticeID"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tmpl_language"); ok {
		paramMap["TmplLanguage"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("monitor_type"); ok {
		paramMap["MonitorType"] = helper.String(v.(string))
	}

	var respData []*monitorv20230616.NoticeContentTmpl

	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeNoticeContentTmplsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	noticeContentTmplList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, tmpl := range respData {
			tmplMap := map[string]interface{}{}

			if tmpl.TmplID != nil {
				tmplMap["tmpl_id"] = tmpl.TmplID
			}

			if tmpl.TmplName != nil {
				tmplMap["tmpl_name"] = tmpl.TmplName
			}

			if tmpl.MonitorType != nil {
				tmplMap["monitor_type"] = tmpl.MonitorType
			}

			if tmpl.TmplLanguage != nil {
				tmplMap["tmpl_language"] = tmpl.TmplLanguage
			}

			if tmpl.Creator != nil {
				tmplMap["creator"] = tmpl.Creator
			}

			if tmpl.CreateTime != nil {
				tmplMap["create_time"] = tmpl.CreateTime
			}

			if tmpl.UpdateTime != nil {
				tmplMap["update_time"] = tmpl.UpdateTime
			}

			if tmpl.TmplContents != nil {
				contentsJSON, err := json.Marshal(tmpl.TmplContents)
				if err == nil {
					tmplMap["tmpl_contents_json"] = string(contentsJSON)
				}
			}

			noticeContentTmplList = append(noticeContentTmplList, tmplMap)
		}

		_ = d.Set("notice_content_tmpl_list", noticeContentTmplList)
	}

	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
