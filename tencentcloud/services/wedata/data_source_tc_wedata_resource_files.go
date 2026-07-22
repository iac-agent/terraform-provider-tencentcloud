package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataResourceFiles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataResourceFilesRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "项目 ID",
			},

			"resource_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resource file 名称 (fuzzy search keyword)。",
			},

			"parent_folder_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "指定path of the file's parent folder (for example /a/b/c，querying resource files under the folder c)。",
			},

			"create_user_uin": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "创建者 ID. obtain through the DescribeCurrentUserInfo API。",
			},

			"modify_time_start": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "更新时间 range. 指定start time in yyyy-MM-dd HH:MM:ss 格式",
			},

			"modify_time_end": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "更新时间 range. 指定end time in yyyy-MM-dd HH:MM:ss 格式",
			},

			"create_time_start": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "创建时间 range. 指定start time in yyyy-MM-dd HH:MM:ss 格式",
			},

			"create_time_end": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "创建时间 range. 指定termination time in yyyy-MM-dd HH:MM:ss 格式",
			},

			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Retrieve the resource file list。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Resource file ID。",
						},
						"resource_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Resource file 名称",
						},
						"file_extension_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "指定resource file 类型",
						},
						"local_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Resource 路径",
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

func dataSourceTencentCloudWedataResourceFilesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_resource_files.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("project_id"); ok {
		paramMap["ProjectId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_name"); ok {
		paramMap["ResourceName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("parent_folder_path"); ok {
		paramMap["ParentFolderPath"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("create_user_uin"); ok {
		paramMap["CreateUserUin"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("modify_time_start"); ok {
		paramMap["ModifyTimeStart"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("modify_time_end"); ok {
		paramMap["ModifyTimeEnd"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("create_time_start"); ok {
		paramMap["CreateTimeStart"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("create_time_end"); ok {
		paramMap["CreateTimeEnd"] = helper.String(v.(string))
	}

	var respData []*wedatav20250806.ResourceFileItem
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataResourceFilesByFilter(ctx, paramMap)
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

		if items.ResourceId != nil {
			ids = append(ids, *items.ResourceId)
			itemsMap["resource_id"] = items.ResourceId
		}

		if items.ResourceName != nil {
			itemsMap["resource_name"] = items.ResourceName
		}

		if items.FileExtensionType != nil {
			itemsMap["file_extension_type"] = items.FileExtensionType
		}

		if items.LocalPath != nil {
			itemsMap["local_path"] = items.LocalPath
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
