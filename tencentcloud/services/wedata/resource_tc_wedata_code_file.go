package wedata

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataCodeFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataCodeFileCreate,
		Read:   resourceTencentCloudWedataCodeFileRead,
		Update: resourceTencentCloudWedataCodeFileUpdate,
		Delete: resourceTencentCloudWedataCodeFileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "项目 ID",
			},

			"code_file_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "代码 文件 名称",
			},

			"parent_folder_path": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Parent 文件夹 路径，对于 示例 /aaa/bbb/ccc，路径 头部 必须 start 使用 slash，root directory pass /。",
			},

			"code_file_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "代码 文件 配置。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"params": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Advanced runtime 参数，variable substitution，map-json String,String。",
						},
						"notebook_session_info": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Notebook kernel 会话 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"notebook_session_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Session ID。",
									},
									"notebook_session_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Session 名称",
									},
								},
							},
						},
					},
				},
			},

			"code_file_content": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "代码 文件 内容",
			},

			// computed
			"code_file_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "代码 文件 ID。",
			},

			"access_scope": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Permission 范围: SHARED，PRIVATE。",
			},

			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "full 路径 的 节点，/aaa/bbb/ccc.ipynb，consists 的 names 的 each 节点。",
			},
		},
	}
}

func resourceTencentCloudWedataCodeFileCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_code_file.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = wedatav20250806.NewCreateCodeFileRequest()
		response   = wedatav20250806.NewCreateCodeFileResponse()
		projectId  string
		codeFileId string
	)

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOk("code_file_name"); ok {
		request.CodeFileName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("parent_folder_path"); ok {
		request.ParentFolderPath = helper.String(v.(string))
	}

	if codeFileConfigMap, ok := helper.InterfacesHeadMap(d, "code_file_config"); ok {
		codeFileConfig := wedatav20250806.CodeFileConfig{}
		if v, ok := codeFileConfigMap["params"].(string); ok && v != "" {
			codeFileConfig.Params = helper.String(v)
		}

		if notebookSessionInfoMap, ok := helper.ConvertInterfacesHeadToMap(codeFileConfigMap["notebook_session_info"]); ok {
			notebookSessionInfo := wedatav20250806.NotebookSessionInfo{}
			if v, ok := notebookSessionInfoMap["notebook_session_id"].(string); ok && v != "" {
				notebookSessionInfo.NotebookSessionId = helper.String(v)
			}

			if v, ok := notebookSessionInfoMap["notebook_session_name"].(string); ok && v != "" {
				notebookSessionInfo.NotebookSessionName = helper.String(v)
			}

			codeFileConfig.NotebookSessionInfo = &notebookSessionInfo
		}

		request.CodeFileConfig = &codeFileConfig
	}

	if v, ok := d.GetOk("code_file_content"); ok {
		request.CodeFileContent = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().CreateCodeFileWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("create wedata code file failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create wedata code file failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Data.CodeFileId == nil {
		return fmt.Errorf("CodeFileId is nil.")
	}

	codeFileId = *response.Response.Data.CodeFileId
	d.SetId(strings.Join([]string{projectId, codeFileId}, tccommon.FILED_SP))
	return resourceTencentCloudWedataCodeFileRead(d, meta)
}

func resourceTencentCloudWedataCodeFileRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_code_file.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	codeFileId := idSplit[1]

	respData, err := service.DescribeWedataCodeFileById(ctx, projectId, codeFileId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_wedata_code_file` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.ProjectId != nil {
		_ = d.Set("project_id", respData.ProjectId)
	}

	if respData.CodeFileName != nil {
		_ = d.Set("code_file_name", respData.CodeFileName)
	}

	if respData.ParentFolderPath != nil {
		if *respData.ParentFolderPath == "" {
			_ = d.Set("parent_folder_path", "/")
		} else {
			_ = d.Set("parent_folder_path", *respData.ParentFolderPath)
		}
	}

	if respData.CodeFileConfig != nil {
		codeFileConfigMap := map[string]interface{}{}
		if respData.CodeFileConfig.Params != nil {
			codeFileConfigMap["params"] = respData.CodeFileConfig.Params
		}

		notebookSessionInfoMap := map[string]interface{}{}
		if respData.CodeFileConfig.NotebookSessionInfo != nil {
			if respData.CodeFileConfig.NotebookSessionInfo.NotebookSessionId != nil {
				notebookSessionInfoMap["notebook_session_id"] = respData.CodeFileConfig.NotebookSessionInfo.NotebookSessionId
			}

			if respData.CodeFileConfig.NotebookSessionInfo.NotebookSessionName != nil {
				notebookSessionInfoMap["notebook_session_name"] = respData.CodeFileConfig.NotebookSessionInfo.NotebookSessionName
			}

			codeFileConfigMap["notebook_session_info"] = []interface{}{notebookSessionInfoMap}
		}

		_ = d.Set("code_file_config", []interface{}{codeFileConfigMap})
	}

	if respData.CodeFileContent != nil {
		_ = d.Set("code_file_content", respData.CodeFileContent)
	}

	if respData.CodeFileId != nil {
		_ = d.Set("code_file_id", respData.CodeFileId)
	}

	if respData.AccessScope != nil {
		_ = d.Set("access_scope", *respData.AccessScope)
	}

	if respData.Path != nil {
		_ = d.Set("path", *respData.Path)
	}

	return nil
}

func resourceTencentCloudWedataCodeFileUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_code_file.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	codeFileId := idSplit[1]

	needChange := false
	mutableArgs := []string{"code_file_config", "code_file_content"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := wedatav20250806.NewUpdateCodeFileRequest()
		if codeFileConfigMap, ok := helper.InterfacesHeadMap(d, "code_file_config"); ok {
			codeFileConfig := wedatav20250806.CodeFileConfig{}
			if v, ok := codeFileConfigMap["params"].(string); ok && v != "" {
				codeFileConfig.Params = helper.String(v)
			}

			if notebookSessionInfoMap, ok := helper.ConvertInterfacesHeadToMap(codeFileConfigMap["notebook_session_info"]); ok {
				notebookSessionInfo := wedatav20250806.NotebookSessionInfo{}
				if v, ok := notebookSessionInfoMap["notebook_session_id"].(string); ok && v != "" {
					notebookSessionInfo.NotebookSessionId = helper.String(v)
				}

				if v, ok := notebookSessionInfoMap["notebook_session_name"].(string); ok && v != "" {
					notebookSessionInfo.NotebookSessionName = helper.String(v)
				}

				codeFileConfig.NotebookSessionInfo = &notebookSessionInfo
			}

			request.CodeFileConfig = &codeFileConfig
		}

		if v, ok := d.GetOk("code_file_content"); ok {
			request.CodeFileContent = helper.String(v.(string))
		}

		request.ProjectId = &projectId
		request.CodeFileId = &codeFileId
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().UpdateCodeFileWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update wedata code file failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudWedataCodeFileRead(d, meta)
}

func resourceTencentCloudWedataCodeFileDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_code_file.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = wedatav20250806.NewDeleteCodeFileRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	projectId := idSplit[0]
	codeFileId := idSplit[1]

	request.ProjectId = &projectId
	request.CodeFileId = &codeFileId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().DeleteCodeFileWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil || result.Response.Data.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete wedata code file failed, Response is nil."))
		}

		if !*result.Response.Data.Status {
			return resource.NonRetryableError(fmt.Errorf("Delete wedata code file failed, Status is false."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete wedata code file failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
