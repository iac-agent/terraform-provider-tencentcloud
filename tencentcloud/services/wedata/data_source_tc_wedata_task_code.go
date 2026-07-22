package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataTaskCode() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataTaskCodeRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "项目 ID",
			},

			"task_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "任务 ID",
			},

			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Retrieves 任务 代码 结果",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code_info": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "代码 内容",
						},
						"code_file_size": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "代码 文件 大小. 单位: KB。",
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

func dataSourceTencentCloudWedataTaskCodeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_task_code.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		projectId string
		taskId    string
	)
	service := WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
		paramMap["ProjectId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_id"); ok {
		taskId = v.(string)
		paramMap["TaskId"] = helper.String(v.(string))
	}

	var respData *wedatav20250806.TaskCodeResult
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataTaskCodeByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	dataMap := map[string]interface{}{}

	if respData.CodeInfo != nil {
		dataMap["code_info"] = respData.CodeInfo
	}

	if respData.CodeFileSize != nil {
		dataMap["code_file_size"] = respData.CodeFileSize
	}

	_ = d.Set("data", []interface{}{dataMap})

	d.SetId(projectId + tccommon.FILED_SP + taskId)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataMap); e != nil {
			return e
		}
	}

	return nil
}
