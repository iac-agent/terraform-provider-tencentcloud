package cam

import (
	"context"
	"fmt"
	"log"
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

func ResourceTencentCloudCamUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCamUserCreate,
		Read:   resourceTencentCloudCamUserRead,
		Update: resourceTencentCloudCamUserUpdate,
		Delete: resourceTencentCloudCamUserDelete,
		Importer: &schema.ResourceImporter{
			State: helper.ImportWithDefaultValue(map[string]interface{}{
				"remark":              "",
				"force_delete":        false,
				"use_api":             true,
				"console_login":       false,
				"need_reset_password": true,
			}),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "名称 CAM 用户",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "备注 的 CAM 用户",
			},
			"force_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate 是否force deletes CAM 用户 如果 集合 false， API secret 键 将 是 checked 和 failed 当 exists; otherwise 用户 将 是 删除 directly. 默认为 false。",
			},
			"use_api": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicate 是否generate API secret 键 或 不。",
			},
			"console_login": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate 是否CAM 用户 可以 login 到 web console 或 不。",
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Sensitive:    true,
				ValidateFunc: tccommon.ValidateAsConfigPassword,
				Description:  "密码 的 CAM 用户 密码 should 是 在 least 8 字符 和 无 more 比 32 字符，includes uppercase letters，lowercase letters，numbers 和 special 字符. Only 必填 当 `console_login` 是 true. 如果 不 集合， random 密码 将 是 automatically generated。",
			},
			"need_reset_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicate 是否CAM 用户 need 到 reset 密码 当 first logins。",
			},
			"phone_num": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Phone 数量 CAM 用户",
			},
			"country_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Country 代码 的 phone 数量，对于 示例: '86'。",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email 的 CAM 用户",
			},
			"uin": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Uin 的 CAM 用户",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Secret 键 的 CAM 用户",
			},
			"secret_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Secret ID CAM 用户",
			},
			"uid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID CAM 用户",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A 列表 标签 用于associate different resources。",
			},
		},
	}
}

func resourceTencentCloudCamUserCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_user.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := cam.NewAddUserRequest()
	request.Name = helper.String(d.Get("name").(string))
	//optional
	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("use_api"); ok {
		apiBool := v.(bool)
		apiInt := uint64(1)
		if !apiBool {
			apiInt = uint64(0)
		}
		request.UseApi = &apiInt
	}
	if v, ok := d.GetOkExists("console_login"); ok {
		loginBool := v.(bool)
		loginInt := uint64(1)
		if !loginBool {
			loginInt = uint64(0)
		}
		request.ConsoleLogin = &loginInt
	}
	if v, ok := d.GetOkExists("need_reset_password"); ok {
		resetBool := v.(bool)
		resetInt := uint64(1)
		if !resetBool {
			resetInt = uint64(0)
		}
		request.NeedResetPassword = &resetInt
	}
	if v, ok := d.GetOk("phone_num"); ok {
		request.PhoneNum = helper.String(v.(string))
	}
	if v, ok := d.GetOk("country_code"); ok {
		request.CountryCode = helper.String(v.(string))
	}
	if v, ok := d.GetOk("email"); ok {
		request.Email = helper.String(v.(string))
	}
	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	var response *cam.AddUserResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCamClient().AddUser(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				if strings.Contains(errCode, "SubUserNameInUse") {
					return resource.NonRetryableError(e)
				}
			}
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create user failed, reason:%s\n", logId, err.Error())
		return err
	}
	if response.Response.Uid == nil {
		return fmt.Errorf("[CHECK][CAM user][Create] check: CAM user id returns nil")
	}

	d.SetId(*response.Response.Name)
	_ = d.Set("secret_key", *response.Response.SecretKey)
	_ = d.Set("password", *response.Response.Password)
	_ = d.Set("secret_id", *response.Response.SecretId)

	//get really instance then read
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	camService := CamService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, e := camService.DescribeUserById(ctx, *response.Response.Name)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if instance == nil {
			return resource.RetryableError(fmt.Errorf("creation not done"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s wait for CAM user ready failed, reason:%s\n", logId, err.Error())
		return err
	}

	//modify tags
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := tccommon.BuildTagResourceName("cam", "uin", region, helper.UInt64ToStr(*response.Response.Uin))
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	time.Sleep(10 * time.Second)
	return resourceTencentCloudCamUserRead(d, meta)
}

func resourceTencentCloudCamUserRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_user.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	deleteForce := false
	if v, ok := d.GetOkExists("force_delete"); ok {
		deleteForce = v.(bool)
		_ = d.Set("force_delete", deleteForce)
	}

	userId := d.Id()
	camService := CamService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var instance *cam.GetUserResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := camService.DescribeUserById(ctx, userId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM user failed, reason:%s\n", logId, err.Error())
		return err
	}

	if instance == nil || instance.Response == nil || instance.Response.Uid == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", userId)
	_ = d.Set("uin", int(*instance.Response.Uin))
	_ = d.Set("uid", int(*instance.Response.Uid))
	_ = d.Set("remark", instance.Response.Remark)
	_ = d.Set("phone_num", instance.Response.PhoneNum)
	_ = d.Set("country_code", instance.Response.CountryCode)
	_ = d.Set("email", instance.Response.Email)
	if int(*instance.Response.ConsoleLogin) == 0 {
		_ = d.Set("console_login", false)
	} else if int(*instance.Response.ConsoleLogin) == 1 {
		_ = d.Set("console_login", true)
	}

	//tags
	tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
	tags, err := tagService.DescribeResourceTags(ctx, "cam", "uin", region, helper.UInt64ToStr(*instance.Response.Uin))
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudCamUserUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_user.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	userId := d.Id()

	var updateAttrs []string

	request := cam.NewUpdateUserRequest()
	request.Name = &userId
	if d.HasChange("remark") {
		request.Remark = helper.String(d.Get("remark").(string))
		updateAttrs = append(updateAttrs, "remark")
	}

	if d.HasChange("console_login") {
		consoleLogin := d.Get("console_login").(bool)
		consoleLogin64 := uint64(0)
		if consoleLogin {
			consoleLogin64 = uint64(1)
		}
		request.ConsoleLogin = &consoleLogin64
		updateAttrs = append(updateAttrs, "console_login")
	}

	if d.HasChange("password") {
		password := d.Get("password").(string)
		request.Password = helper.String(password)
		updateAttrs = append(updateAttrs, "password")
	}

	if d.HasChange("need_reset_password") {
		resetBool := d.Get("need_reset_password").(bool)
		resetBool64 := uint64(0)
		if resetBool {
			resetBool64 = uint64(1)
		}
		request.NeedResetPassword = &resetBool64
		updateAttrs = append(updateAttrs, "need_reset_password")
	}

	if d.HasChange("phone_num") || d.HasChange("country_code") {
		request.PhoneNum = helper.String(d.Get("phone_num").(string))
		updateAttrs = append(updateAttrs, "phone_num")
		request.CountryCode = helper.String(d.Get("country_code").(string))
		updateAttrs = append(updateAttrs, "country_code")
	}
	if d.HasChange("email") {
		request.Email = helper.String(d.Get("email").(string))
		updateAttrs = append(updateAttrs, "email")
	}

	if len(updateAttrs) > 0 {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCamClient().UpdateUser(request)

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
			log.Printf("[CRITAL]%s update CAM user description failed, reason:%s\n", logId, err.Error())
			return err
		}
	}

	//tag
	if d.HasChange("tags") {
		camService := CamService{
			client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
		}

		var instance *cam.GetUserResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := camService.DescribeUserById(ctx, userId)
			if e != nil {
				return tccommon.RetryError(e)
			}
			instance = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read CAM user failed, reason:%s\n", logId, err.Error())
			return err
		}

		if instance == nil || instance.Response == nil || instance.Response.Uid == nil {
			d.SetId("")
			return nil
		}

		oldInterface, newInterface := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldInterface.(map[string]interface{}), newInterface.(map[string]interface{}))
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := tccommon.BuildTagResourceName("cam", "uin", region, helper.UInt64ToStr(*instance.Response.Uin))
		err = tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}

	}

	return nil
}

func resourceTencentCloudCamUserDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cam_user.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	userId := d.Id()
	request := cam.NewDeleteUserRequest()
	request.Name = &userId

	//check is force delete or not
	deleteForce := false
	if v, ok := d.GetOkExists("force_delete"); ok {
		deleteForce = v.(bool)
	}

	request.Force = helper.BoolToInt64Pointer(deleteForce)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCamClient().DeleteUser(request)
		if e != nil {
			log.Printf("[CRITAL]%s reason[%s]\n", logId, e.Error())
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CAM user failed, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
