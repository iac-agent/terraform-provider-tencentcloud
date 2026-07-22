package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataWorkflowFolders() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataWorkflowFoldersRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "项目 ID",
			},

			"parent_folder_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Parent 文件夹 absolute 路径，对于 示例 /abc/de，如果 它 是 root directory，pass /。",
			},

			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Paginated 文件夹 查询 结果",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "项目 ID",
						},
						"folder_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Folder ID。",
						},
						"folder_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "指定absolute 路径 的 文件夹。",
						},
						"create_user_uin": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "创建者 ID。",
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

func dataSourceTencentCloudWedataWorkflowFoldersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_workflow_folders.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("project_id"); ok {
		paramMap["ProjectId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("parent_folder_path"); ok {
		paramMap["ParentFolderPath"] = helper.String(v.(string))
	}

	var respData []*wedatav20250806.WorkflowFolder
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataWorkflowFoldersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(respData))
	itemsList := make([]map[string]interface{}, 0, len(respData))

	for _, items := range respData {
		itemsMap := map[string]interface{}{}

		if items.ProjectId != nil {
			itemsMap["project_id"] = items.ProjectId
		}

		if items.FolderId != nil {
			itemsMap["folder_id"] = items.FolderId
			ids = append(ids, *items.FolderId)
		}

		if items.FolderPath != nil {
			itemsMap["folder_path"] = items.FolderPath
		}

		if items.CreateUserUin != nil {
			itemsMap["create_user_uin"] = items.CreateUserUin
		}

		itemsList = append(itemsList, itemsMap)
	}

	_ = d.Set("data", itemsList)

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), itemsList); e != nil {
			return e
		}
	}

	return nil
}
