package dnspod

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDnspodRecordLineList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDnspodRecordLineListRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "域名",
			},

			"domain_grade": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "域名 级别 + Old packages: D_FREE，D_PLUS，D_EXTRA，D_EXPERT，D_ULTRA correspond 到 free 包，personal luxury，enterprise 1，enterprise 2，enterprise 3. + New packages: DP_FREE，DP_PLUS，DP_EXTRA，DP_EXPERT，DP_ULTRA correspond 到 new free，personal professional，enterprise basic，enterprise standard，enterprise flagship。",
			},

			"domain_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "域名 ID. 参数 DomainId has higher 优先级 比 参数 域名 如果 参数 DomainId 是 passed， 参数 域名 将 是 ignored. You 可以 find all Domains 和 DomainIds through DescribeDomainList interface。",
			},

			"line_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Line 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Line 名称",
						},
						"line_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Line ID。",
						},
					},
				},
			},

			"line_group_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Line 组 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"line_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Line 组 ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Line 组名称",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group 类型",
						},
						"line_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "列表 lines included 在 line 组。",
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

func dataSourceTencentCloudDnspodRecordLineListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dnspod_record_line_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		domain        string
		lineList      []*dnspod.LineInfo
		lineGroupList []*dnspod.LineGroupInfo
		e             error
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("domain"); ok {
		paramMap["Domain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain_grade"); ok {
		paramMap["DomainGrade"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("domain_id"); ok {
		paramMap["DomainId"] = helper.IntUint64(v.(int))
	}

	service := DnspodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	// var lineList []*dnspod.LineInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		lineList, lineGroupList, e = service.DescribeDnspodRecordLineListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		// lineList = result
		return nil
	})
	if err != nil {
		return err
	}

	// ids := make([]string, 0, len(lineList))
	if lineList != nil {
		tmpList := make([]map[string]interface{}, 0, len(lineList))
		for _, lineInfo := range lineList {
			lineInfoMap := map[string]interface{}{}

			if lineInfo.Name != nil {
				lineInfoMap["name"] = lineInfo.Name
			}

			if lineInfo.LineId != nil {
				lineInfoMap["line_id"] = lineInfo.LineId
			}

			// ids = append(ids, *lineInfo.Domain)
			tmpList = append(tmpList, lineInfoMap)
		}

		_ = d.Set("line_list", tmpList)
	}

	if lineGroupList != nil {
		tmpList := make([]map[string]interface{}, 0, len(lineGroupList))
		for _, lineGroupInfo := range lineGroupList {
			lineGroupInfoMap := map[string]interface{}{}

			if lineGroupInfo.LineId != nil {
				lineGroupInfoMap["line_id"] = lineGroupInfo.LineId
			}

			if lineGroupInfo.Name != nil {
				lineGroupInfoMap["name"] = lineGroupInfo.Name
			}

			if lineGroupInfo.Type != nil {
				lineGroupInfoMap["type"] = lineGroupInfo.Type
			}

			if lineGroupInfo.LineList != nil {
				lineGroupInfoMap["line_list"] = lineGroupInfo.LineList
			}

			// ids = append(ids, *lineGroupInfo.Domain)
			tmpList = append(tmpList, lineGroupInfoMap)
		}

		_ = d.Set("line_group_list", tmpList)
	}

	d.SetId(helper.DataResourceIdHash(domain))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		e = tccommon.WriteToFile(output.(string), map[string]interface{}{
			"line_list":       lineList,
			"line_group_list": lineGroupList,
		})
		if e != nil {
			return e
		}
	}
	return nil
}
