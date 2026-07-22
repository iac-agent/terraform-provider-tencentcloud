package cam

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCamRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamRoleCreate,
		Read:   resourceTencentCloudCamRoleRead,
		Update: resourceTencentCloudCamRoleUpdate,
		Delete: resourceTencentCloudCamRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "名称 CAM 角色",
			},
			"document": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, olds, news string, d *schema.ResourceData) bool {
					var oldJson interface{}
					err := json.Unmarshal([]byte(olds), &oldJson)
					if err != nil {
						return olds == news
					}
					var newJson interface{}
					err = json.Unmarshal([]byte(news), &newJson)
					if err != nil {
						return olds == news
					}
					flag := reflect.DeepEqual(oldJson, newJson)
					return flag
				},
				Description: "Document 的 CAM 角色 syntax refers 到 [CAM POLICY](https://intl.云.tencent.com/document/product/598/10604). There 是 some notes 当 使用 此 para 在 terraform: 1. elements 在 json claimed supporting two types 作为 `字符串` 和 `数组` 仅 support 类型 `数组`; 2. Terraform does 不 support `root` syntax，当 appears，它 必须 是 replaced 使用 uin 它 stands 对于。",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "描述 CAM 角色",
			},
			"console_login": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "表示是否CAM 角色 可以 login 或 不。",
			},
			"session_duration": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "最大 validity 周期 的 temporary 键 对于 creating 角色",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "创建时间 的 CAM 角色",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "last 更新时间 的 CAM 角色",
			},
			"role_arn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "RoleArn Information 对于 Roles。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A 列表 标签 用于associate different resources。",
			},
		},
	}
}

func resourceTencentCloudCamRoleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_role.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	name := d.Get("name").(string)
	document := d.Get("document").(string)
	camService := CamService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	documentErr := camService.PolicyDocumentForceCheck(document)
	if documentErr != nil {
		return documentErr
	}

	request := cam.NewCreateRoleRequest()
	request.RoleName = &name
	request.PolicyDocument = &document
	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("console_login"); ok {
		loginBool := v.(bool)
		loginInt := uint64(1)
		if !loginBool {
			loginInt = uint64(0)
		}
		request.ConsoleLogin = &loginInt
	}
	if v, ok := d.GetOkExists("session_duration"); ok {
		request.SessionDuration = helper.IntUint64(v.(int))
	}

	var response *cam.CreateRoleResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCamClient().CreateRole(request)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				//check if read empty
				if strings.Contains(errCode, "RoleNameInUse") {
					return resource.NonRetryableError(e)
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.RoleId == nil {
			return resource.NonRetryableError(fmt.Errorf("Create CAM role failed, Response is nil."))
		}

		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create CAM role failed, reason:%s\n", logId, err.Error())
		return err
	}

	d.SetId(*response.Response.RoleId)

	//get really instance then read
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	roleId := d.Id()

	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, e := camService.DescribeRoleById(ctx, roleId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if instance == nil {
			return resource.RetryableError(fmt.Errorf("creation not done"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM role failed, reason:%s\n", logId, err.Error())
		return err
	}

	//modify tags
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		resourceName := tccommon.BuildTagResourceName("cam", "role", "", roleId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	time.Sleep(5 * time.Second)
	return resourceTencentCloudCamRoleRead(d, meta)
}

func resourceTencentCloudCamRoleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_role.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	roleId := d.Id()
	camService := CamService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var instance *cam.RoleInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := camService.DescribeRoleById(ctx, roleId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM role failed, reason:%s\n", logId, err.Error())
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", instance.RoleName)
	_ = d.Set("document", instance.PolicyDocument)
	_ = d.Set("session_duration", instance.SessionDuration)
	_ = d.Set("create_time", instance.AddTime)
	_ = d.Set("update_time", instance.UpdateTime)
	if instance.Description != nil {
		_ = d.Set("description", instance.Description)
	}

	if instance.RoleArn != nil {
		_ = d.Set("role_arn", instance.RoleArn)
	}

	if instance.ConsoleLogin != nil {
		if int(*instance.ConsoleLogin) == 1 {
			_ = d.Set("console_login", true)
		} else {
			_ = d.Set("console_login", false)
		}
	} else {
		_ = d.Set("console_login", false)
	}

	//tags
	tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	tags, err := tagService.DescribeResourceTags(ctx, "cam", "role", "", roleId)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudCamRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_role.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	d.Partial(true)

	roleId := d.Id()

	description := ""
	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			description = v.(string)
		}
		mDescRequest := cam.NewUpdateRoleDescriptionRequest()
		mDescRequest.Description = &description
		mDescRequest.RoleId = &roleId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCamClient().UpdateRoleDescription(mDescRequest)

			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, mDescRequest.GetAction(), mDescRequest.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, mDescRequest.GetAction(), mDescRequest.ToJsonString(), response.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update CAM role description failed, reason:%s\n", logId, err.Error())
			return err
		}

	}
	document := ""
	if d.HasChange("document") {

		document = d.Get("document").(string)
		camService := CamService{
			client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
		}
		documentErr := camService.PolicyDocumentForceCheck(document)
		if documentErr != nil {
			return documentErr
		}
		mDocRequest := cam.NewUpdateAssumeRolePolicyRequest()
		mDocRequest.PolicyDocument = &document
		mDocRequest.RoleId = &roleId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCamClient().UpdateAssumeRolePolicy(mDocRequest)

			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, mDocRequest.GetAction(), mDocRequest.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, mDocRequest.GetAction(), mDocRequest.ToJsonString(), response.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update CAM role document failed, reason:%s\n", logId, err.Error())
			return err
		}

	}

	d.Partial(false)

	//tag
	if d.HasChange("tags") {
		oldInterface, newInterface := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldInterface.(map[string]interface{}), newInterface.(map[string]interface{}))
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		resourceName := tccommon.BuildTagResourceName("cam", "role", "", roleId)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}

	}

	if d.HasChange("console_login") {
		consoleLoginRequest := cam.NewUpdateRoleConsoleLoginRequest()

		if v, ok := d.GetOkExists("console_login"); ok {
			loginBool := v.(bool)
			loginInt := int64(1)
			if !loginBool {
				loginInt = int64(0)
			}
			consoleLoginRequest.ConsoleLogin = &loginInt
		}

		consoleLoginRequest.RoleId = helper.StrToInt64Point(roleId)

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCamClient().UpdateRoleConsoleLogin(consoleLoginRequest)

			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, consoleLoginRequest.GetAction(), consoleLoginRequest.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, consoleLoginRequest.GetAction(), consoleLoginRequest.ToJsonString(), response.ToJsonString())
			}
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update CAM role console login failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	if d.HasChange("session_duration") {
		request := cam.NewUpdateRoleSessionDurationRequest()
		request.RoleId = helper.StrToUint64Point(roleId)
		if v, ok := d.GetOkExists("session_duration"); ok {
			request.SessionDuration = helper.IntUint64(v.(int))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCamClient().UpdateRoleSessionDuration(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update CAM role session duration failed, reason:%s\n", logId, err.Error())
			return err
		}
	}
	return resourceTencentCloudCamRoleRead(d, meta)
}

func resourceTencentCloudCamRoleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_role.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	roleId := d.Id()
	camService := CamService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := camService.DeleteRoleById(ctx, roleId)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CAM role failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
