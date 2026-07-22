package wedata

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataProjectCreate,
		Read:   resourceTencentCloudWedataProjectRead,
		Update: resourceTencentCloudWedataProjectUpdate,
		Delete: resourceTencentCloudWedataProjectDelete,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Project basic 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Project identifier，English 名称 starting 使用 letter，可以 contain letters，numbers，和 underscores，不能 exceed 32 字符。",
						},
						"display_name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Project display 名称，可以 是 Chinese 名称 starting 使用 letter，可以 contain letters，numbers，和 underscores，不能 exceed 32 字符。",
						},
						"project_model": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Project 模式，SIMPLE (默认值): Simple 模式 STANDARD: Standard 模式",
						},
					},
				},
			},

			"dlc_info": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "DLC binding 集群 信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"compute_resources": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "DLC 资源名称 (need 到 add 角色 Uin 到 DLC，otherwise 可能 不 是 able 到 obtain resources)。",
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
							Description: "Access 账号 (仅 effective 对于 standard 模式 projects 和 必填 对于 standard 模式)，用于submit DLC tasks.\nIt 是 recommended 到 使用 指定 sub-账号 和 集合 corresponding 数据库 表 permissions 对于 sub-账号; 任务 runner 模式 可能 cause 任务 failure 当 responsible person leaves; main 账号 模式 是 不 easy 对于 权限 control 当 多个 projects have different permissions.\n\nEnum 值:\n- TASK_RUNNER (任务 Runner)\n- OWNER (Main 账号 模式)\n- SUB (Sub 账号 模式)。",
						},
						"sub_account_uin": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Sub-账号 ID (仅 effective 对于 standard 模式 projects)，当 AccessAccount 是 在 sub-账号 模式， sub-账号 ID 信息 needs 到 是 指定，other modes do 不 need 到 是 指定。",
						},
					},
				},
			},

			"resource_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "列表 bound 资源 组 IDs。",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description:  "Item 状态: 0: 已禁用，1: 已启用",
			},

			// computed
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "项目 ID",
			},
		},
	}
}

func resourceTencentCloudWedataProjectCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_project.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = wedatav20250806.NewCreateProjectRequest()
		response  = wedatav20250806.NewCreateProjectResponse()
		projectId string
	)

	if projectMap, ok := helper.InterfacesHeadMap(d, "project"); ok {
		projectRequest := wedatav20250806.ProjectRequest{}
		if v, ok := projectMap["project_name"].(string); ok && v != "" {
			projectRequest.ProjectName = helper.String(v)
		}

		if v, ok := projectMap["display_name"].(string); ok && v != "" {
			projectRequest.DisplayName = helper.String(v)
		}

		if v, ok := projectMap["project_model"].(string); ok && v != "" {
			projectRequest.ProjectModel = helper.String(v)
		}

		request.Project = &projectRequest
	}

	if dLCInfoMap, ok := helper.InterfacesHeadMap(d, "dlc_info"); ok {
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

		request.DLCInfo = &dLCClusterInfo
	}

	if v, ok := d.GetOk("resource_ids"); ok {
		resourceIdsSet := v.(*schema.Set).List()
		for i := range resourceIdsSet {
			resourceIds := resourceIdsSet[i].(string)
			request.ResourceIds = append(request.ResourceIds, helper.String(resourceIds))
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().CreateProjectWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Create wedata project failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create wedata project failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Data.ProjectId == nil {
		return fmt.Errorf("ProjectId is nil.")
	}

	projectId = *response.Response.Data.ProjectId
	d.SetId(projectId)

	// set status
	if v, ok := d.GetOkExists("status"); ok {
		if v.(int) == 0 {
			disableReq := wedatav20250806.NewDisableProjectRequest()
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().DisableProjectWithContext(ctx, disableReq)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, disableReq.GetAction(), disableReq.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil || result.Response.Data == nil {
					return resource.NonRetryableError(fmt.Errorf("Disable wedata project failed, Response is nil."))
				}

				if result.Response.Data.Status != nil && *result.Response.Data.Status {
					return nil
				}

				return resource.NonRetryableError(fmt.Errorf("Disable wedata project failed, Status is false."))
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s disable wedata project failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		}
	}

	return resourceTencentCloudWedataProjectRead(d, meta)
}

func resourceTencentCloudWedataProjectRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_project.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		projectId = d.Id()
	)

	respData, err := service.DescribeWedataProjectById(ctx, projectId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_wedata_project` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	dMapProject := make(map[string]interface{}, 0)
	if respData.ProjectName != nil {
		dMapProject["project_name"] = *respData.ProjectName
	}

	if respData.DisplayName != nil {
		dMapProject["display_name"] = *respData.DisplayName
	}

	if respData.ProjectModel != nil {
		dMapProject["project_model"] = *respData.ProjectModel
	}

	_ = d.Set("project", []interface{}{dMapProject})

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.ProjectId != nil {
		_ = d.Set("project_id", respData.ProjectId)
	}

	return nil
}

func resourceTencentCloudWedataProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_project.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		projectId = d.Id()
	)

	immutableArgs := []string{"dlc_info", "resource_ids"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("project.0.display_name") {
		request := wedatav20250806.NewUpdateProjectRequest()
		if v, ok := d.GetOk("project.0.display_name"); ok {
			request.DisplayName = helper.String(v.(string))
		}

		request.ProjectId = &projectId
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().UpdateProjectWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update wedata project failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if d.HasChange("status") {
		if v, ok := d.GetOkExists("status"); ok {
			if v.(int) == 1 {
				enableReq := wedatav20250806.NewEnableProjectRequest()
				enableReq.ProjectId = &projectId
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().EnableProjectWithContext(ctx, enableReq)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, enableReq.GetAction(), enableReq.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil || result.Response.Data == nil {
						return resource.NonRetryableError(fmt.Errorf("Enable wedata project failed, Response is nil."))
					}

					if result.Response.Data.Status != nil && *result.Response.Data.Status {
						return nil
					}

					return resource.NonRetryableError(fmt.Errorf("Enable wedata project failed, Status is false."))
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s enable wedata project failed, reason:%+v", logId, reqErr)
					return reqErr
				}
			} else {
				disableReq := wedatav20250806.NewDisableProjectRequest()
				disableReq.ProjectId = &projectId
				reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().DisableProjectWithContext(ctx, disableReq)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, disableReq.GetAction(), disableReq.ToJsonString(), result.ToJsonString())
					}

					if result == nil || result.Response == nil || result.Response.Data == nil {
						return resource.NonRetryableError(fmt.Errorf("Disable wedata project failed, Response is nil."))
					}

					if result.Response.Data.Status != nil && *result.Response.Data.Status {
						return nil
					}

					return resource.NonRetryableError(fmt.Errorf("Disable wedata project failed, Status is false."))
				})

				if reqErr != nil {
					log.Printf("[CRITAL]%s disable wedata project failed, reason:%+v", logId, reqErr)
					return reqErr
				}
			}
		}
	}

	return resourceTencentCloudWedataProjectRead(d, meta)
}

func resourceTencentCloudWedataProjectDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_project.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = wedatav20250806.NewDeleteProjectRequest()
		response  = wedatav20250806.NewDeleteProjectResponse()
		projectId = d.Id()
	)

	request.ProjectId = &projectId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().DeleteProjectWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil || result.Response.Data.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete project failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete wedata project failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if *response.Response.Data.Status {
		return nil
	}

	return fmt.Errorf("Delete project %s failed, Status is false.", projectId)
}
