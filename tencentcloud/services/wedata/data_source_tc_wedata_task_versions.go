package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataTaskVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataTaskVersionsRead,
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

			"task_version_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SAVE 版本\nSUBMIT 版本\n默认为 SAVE。",
			},

			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Task 版本 list。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"create_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "创建时间。",
						},
						"version_num": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "版本 number。",
						},
						"create_user_uin": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "创建者 ID。",
						},
						"version_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Saved 版本 ID。",
						},
						"version_remark": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "版本 描述",
						},
						"approve_status": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Approval 状态 (only for 提交 版本)。",
						},
						"status": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Production 状态 (only for 提交 版本)。",
						},
						"approve_user_uin": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Approver (only for 提交 版本)。",
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

func dataSourceTencentCloudWedataTaskVersionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_task_versions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("project_id"); ok {
		paramMap["ProjectId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_id"); ok {
		paramMap["TaskId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_version_type"); ok {
		paramMap["TaskVersionType"] = helper.String(v.(string))
	}

	var respData []*wedatav20250806.TaskVersion
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataTaskVersionsByFilter(ctx, paramMap)
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

		if items.CreateTime != nil {
			itemsMap["create_time"] = items.CreateTime
		}

		if items.VersionNum != nil {
			itemsMap["version_num"] = items.VersionNum
		}

		if items.CreateUserUin != nil {
			itemsMap["create_user_uin"] = items.CreateUserUin
		}

		if items.VersionId != nil {
			itemsMap["version_id"] = items.VersionId
		}

		if items.VersionRemark != nil {
			itemsMap["version_remark"] = items.VersionRemark
		}

		if items.ApproveStatus != nil {
			itemsMap["approve_status"] = items.ApproveStatus
		}

		if items.Status != nil {
			itemsMap["status"] = items.Status
		}

		if items.ApproveUserUin != nil {
			itemsMap["approve_user_uin"] = items.ApproveUserUin
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
