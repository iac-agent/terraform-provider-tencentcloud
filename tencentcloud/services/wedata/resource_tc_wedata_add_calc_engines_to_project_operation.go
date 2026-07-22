package wedata

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataAddCalcEnginesToProjectOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataAddCalcEnginesToProjectOperationCreate,
		Read:   resourceTencentCloudWedataAddCalcEnginesToProjectOperationRead,
		Delete: resourceTencentCloudWedataAddCalcEnginesToProjectOperationDelete,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "项目 ID 到 是 modified。",
			},

			"dlc_info": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "DLC 集群 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compute_resources": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "DLC 资源 names (need 到 add 角色 Uin 到 DLC，otherwise resources 可能 不 是 可用)。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "DLC 地域",
						},
						"default_database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "指定default 数据库 对于 DLC 集群。",
						},
						"standard_mode_env_tag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cluster 配置 标签 (仅 effective 对于 standard 模式 projects 和 必填 对于 standard 模式). Enum 值:\n- Prod (Production 环境)\n- Dev (Development 环境)。",
						},
						"access_account": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Access 账号 (仅 effective 对于 standard 模式 projects 和 必填 对于 standard 模式)，用于submit DLC tasks.\nIt 是 recommended 到 使用 指定 sub-账号 和 集合 corresponding 数据库 表 permissions 对于 sub-账号; 任务 runner 模式 可能 cause 任务 failures 当 responsible person leaves; main 账号 模式 是 不 easy 对于 权限 control 当 多个 projects have different permissions.\n\nEnum 值:\n- TASK_RUNNER (任务 Runner)\n- OWNER (Main 账号 模式)\n- SUB (Sub-账号 模式)。",
						},
						"sub_account_uin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Sub-账号 ID (仅 effective 对于 standard 模式 projects)，当 AccessAccount 是 在 sub-账号 模式， sub-账号 ID 信息 needs 到 是 指定，other modes do 不 need 到 是 指定。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudWedataAddCalcEnginesToProjectOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_add_calc_engines_to_project_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = wedatav20250806.NewAddCalcEnginesToProjectRequest()
		projectId string
	)

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOk("dlc_info"); ok {
		for _, item := range v.([]interface{}) {
			dLCInfoMap := item.(map[string]interface{})
			dLCClusterInfo := wedatav20250806.DLCClusterInfo{}
			if v, ok := dLCInfoMap["compute_resources"]; ok {
				computeResourcesSet := v.(*schema.Set).List()
				for i := range computeResourcesSet {
					computeResources := computeResourcesSet[i].(string)
					dLCClusterInfo.ComputeResources = append(dLCClusterInfo.ComputeResources, helper.String(computeResources))
				}
			}

			if v, ok := dLCInfoMap["region"].(string); ok && v != "" {
				dLCClusterInfo.Region = helper.String(v)
			}

			if v, ok := dLCInfoMap["default_database"].(string); ok && v != "" {
				dLCClusterInfo.DefaultDatabase = helper.String(v)
			}

			if v, ok := dLCInfoMap["standard_mode_env_tag"].(string); ok && v != "" {
				dLCClusterInfo.StandardModeEnvTag = helper.String(v)
			}

			if v, ok := dLCInfoMap["access_account"].(string); ok && v != "" {
				dLCClusterInfo.AccessAccount = helper.String(v)
			}

			if v, ok := dLCInfoMap["sub_account_uin"].(string); ok && v != "" {
				dLCClusterInfo.SubAccountUin = helper.String(v)
			}

			request.DLCInfo = append(request.DLCInfo, &dLCClusterInfo)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().AddCalcEnginesToProjectWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create wedata add calc engines to project operation failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(projectId)
	return resourceTencentCloudWedataAddCalcEnginesToProjectOperationRead(d, meta)
}

func resourceTencentCloudWedataAddCalcEnginesToProjectOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_add_calc_engines_to_project_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudWedataAddCalcEnginesToProjectOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_add_calc_engines_to_project_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
