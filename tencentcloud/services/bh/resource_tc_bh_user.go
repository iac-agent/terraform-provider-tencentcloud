package bh

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bhv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBhUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBhUserCreate,
		Read:   resourceTencentCloudBhUserRead,
		Update: resourceTencentCloudBhUserUpdate,
		Delete: resourceTencentCloudBhUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "用户名，3-20 字符，必须 start 使用 English letter 和 不能 contain 字符 other 比 `letters`，`numbers`，`.`，`_`，`-`。",
			},

			"real_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "用户's real 名称，最大 长度 20 字符，不能 contain whitespace 字符。",
			},

			"phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Input 在 格式 的 \"country 代码|phone 数量\", e.g.: \"+86|xxxxxxxx\". At least 一个 的 phone 和 email 参数 必须 是 提供.",
			},

			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email 地址 At least 一个 的 phone 和 email 参数 必须 是 提供。",
			},

			"validate_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "User effective 时间, e.g.: \"2021-09-22T00:00:00+00:00\". 如果 effective 和 expiration times 是 不 filled, 用户 将 是 有效 permanently.",
			},

			"validate_to": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "User expiration 时间, e.g.: \"2021-09-23T00:00:00+00:00\". 如果 effective 和 expiration times 是 不 filled, 用户 将 是 有效 permanently.",
			},

			"group_id_set": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "用户 组 ID 集合 到 其中 用户 belongs。",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"auth_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Authentication 方法，0 - 本地，1 - LDAP，2 - OAuth. 默认为 0 如果 不 提供。",
			},

			"validate_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Access 时间 restriction， 字符串 composed 的 0 和 1 使用 长度 168 (7 * 24)，representing 时间 slots allowed 对于 用户 在 week. Nth character 在 字符串 表示 Nth hour 在 week，0 - 不 allowed 到 访问，1 - allowed 到 访问。",
			},

			"department_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Department ID 到 其中 用户 belongs, e.g.: \"1.2.3\".",
			},

			// computed
			"user_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "用户 ID。",
			},
		},
	}
}

func resourceTencentCloudBhUserCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_user.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = bhv20230418.NewCreateUserRequest()
		response = bhv20230418.NewCreateUserResponse()
		userId   string
	)

	if v, ok := d.GetOk("user_name"); ok {
		request.UserName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("real_name"); ok {
		request.RealName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("phone"); ok {
		request.Phone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("email"); ok {
		request.Email = helper.String(v.(string))
	}

	if v, ok := d.GetOk("validate_from"); ok {
		request.ValidateFrom = helper.String(v.(string))
	}

	if v, ok := d.GetOk("validate_to"); ok {
		request.ValidateTo = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_id_set"); ok {
		groupIdSetSet := v.(*schema.Set).List()
		for i := range groupIdSetSet {
			groupIdSet := groupIdSetSet[i].(int)
			request.GroupIdSet = append(request.GroupIdSet, helper.IntUint64(groupIdSet))
		}
	}

	if v, ok := d.GetOkExists("auth_type"); ok {
		request.AuthType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("validate_time"); ok {
		request.ValidateTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("department_id"); ok {
		request.DepartmentId = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().CreateUserWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create bh user failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create bh user failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Id == nil {
		return fmt.Errorf("Id is nil.")
	}

	userId = helper.UInt64ToStr(*response.Response.Id)
	d.SetId(userId)
	return resourceTencentCloudBhUserRead(d, meta)
}

func resourceTencentCloudBhUserRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_user.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = BhService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		userId  = d.Id()
	)

	respData, err := service.DescribeBhUserById(ctx, userId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_bh_user` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.UserName != nil {
		_ = d.Set("user_name", respData.UserName)
	}

	if respData.RealName != nil {
		_ = d.Set("real_name", respData.RealName)
	}

	if respData.Phone != nil {
		_ = d.Set("phone", respData.Phone)
	}

	if respData.Email != nil {
		_ = d.Set("email", respData.Email)
	}

	if respData.ValidateFrom != nil {
		_ = d.Set("validate_from", respData.ValidateFrom)
	}

	if respData.ValidateTo != nil {
		_ = d.Set("validate_to", respData.ValidateTo)
	}

	if respData.GroupSet != nil {
		groupIdSetList := make([]uint64, 0, len(respData.GroupSet))
		for _, item := range respData.GroupSet {
			if item.Id != nil {
				groupIdSetList = append(groupIdSetList, *item.Id)
			}
		}

		_ = d.Set("group_id_set", groupIdSetList)
	}

	if respData.AuthType != nil {
		_ = d.Set("auth_type", respData.AuthType)
	}

	if respData.ValidateTime != nil {
		_ = d.Set("validate_time", respData.ValidateTime)
	}

	if respData.DepartmentId != nil {
		dResp, err := service.DescribeBhDepartments(ctx)
		if err != nil {
			return err
		}

		if dResp == nil {
			return fmt.Errorf("Departments is nil")
		}

		if dResp.Enabled != nil && *dResp.Enabled {
			_ = d.Set("department_id", respData.DepartmentId)
		}
	}

	if respData.Id != nil {
		_ = d.Set("user_id", respData.Id)
	}

	return nil
}

func resourceTencentCloudBhUserUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_user.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId  = tccommon.GetLogId(tccommon.ContextNil)
		ctx    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		userId = d.Id()
	)

	needChange := false
	mutableArgs := []string{"real_name", "phone", "email", "validate_from", "validate_to", "group_id_set", "auth_type", "validate_time", "department_id"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := bhv20230418.NewModifyUserRequest()
		if v, ok := d.GetOk("real_name"); ok {
			request.RealName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("phone"); ok {
			request.Phone = helper.String(v.(string))
		}

		if v, ok := d.GetOk("email"); ok {
			request.Email = helper.String(v.(string))
		}

		if v, ok := d.GetOk("validate_from"); ok {
			request.ValidateFrom = helper.String(v.(string))
		}

		if v, ok := d.GetOk("validate_to"); ok {
			request.ValidateTo = helper.String(v.(string))
		}

		if v, ok := d.GetOk("group_id_set"); ok {
			groupIdSetSet := v.(*schema.Set).List()
			for i := range groupIdSetSet {
				groupIdSet := groupIdSetSet[i].(int)
				request.GroupIdSet = append(request.GroupIdSet, helper.IntUint64(groupIdSet))
			}
		}

		if v, ok := d.GetOkExists("auth_type"); ok {
			request.AuthType = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("validate_time"); ok {
			request.ValidateTime = helper.String(v.(string))
		}

		if v, ok := d.GetOk("department_id"); ok {
			request.DepartmentId = helper.String(v.(string))
		}

		request.Id = helper.StrToUint64Point(userId)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().ModifyUserWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update bh user failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudBhUserRead(d, meta)
}

func resourceTencentCloudBhUserDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bh_user.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = bhv20230418.NewDeleteUsersRequest()
		userId  = d.Id()
	)

	request.IdSet = append(request.IdSet, helper.StrToUint64Point(userId))
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBhV20230418Client().DeleteUsersWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete bh user failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
