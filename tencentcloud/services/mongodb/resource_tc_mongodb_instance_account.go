package mongodb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

//internal version: replace import begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
//internal version: replace import end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

func ResourceTencentCloudMongodbInstanceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMongodbInstanceAccountCreate,
		Read:   resourceTencentCloudMongodbInstanceAccountRead,
		Update: resourceTencentCloudMongodbInstanceAccountUpdate,
		Delete: resourceTencentCloudMongodbInstanceAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "实例 ID， 格式 是: cmgo-9d0p6umb.Same 作为 实例 ID displayed 在 云 数据库 console 页面。",
			},

			"user_name": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "new 账号 名称 Its 格式 requirements 是 作为 follows: character 范围 [1,32]. Characters 在 范围 的 [A,Z]，[,z]，[1,9] 作为 well 作为 underscore _ 和 dash - 可以 是 input。",
			},

			"password": {
				Optional:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
				Description: "New 账号 密码 密码 complexity requirements 是 作为 follows: character 长度 范围 [8,32]. 包含at least letters，numbers 和 special 字符 (exclamation point!，在@，pound sign #，percent sign %，caret ^，asterisk *，parentheses ()，underscore _)。",
			},

			"mongo_user_password": {
				Optional:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "密码 corresponding 到 mongouser 账号 mongouser 是 系统 默认值 账号，其中 是 密码 集合 当 creating 实例。",
			},

			"user_desc": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "账号 备注",
			},

			"auth_role": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "read 和 write 权限信息 的 账号",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mask": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "权限信息 的 当前 账号 0: No 权限. 1: read-仅. 2: Write 仅. 3: Read 和 write。",
						},
						"namespace": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Refers 到 名称 数据库 使用 当前 账号 permissions.*: 表示all databases. db.名称: 表示database 的 特定 名称",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMongodbInstanceAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_account.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = mongodb.NewCreateAccountUserRequest()
		response   = mongodb.NewCreateAccountUserResponse()
		instanceId string
		userName   string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_name"); ok {
		userName = v.(string)
		request.UserName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	if v, ok := d.GetOk("mongo_user_password"); ok {
		request.MongoUserPassword = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_desc"); ok {
		request.UserDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("auth_role"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			auth := mongodb.Auth{}
			if v, ok := dMap["mask"]; ok {
				auth.Mask = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["namespace"]; ok {
				auth.NameSpace = helper.String(v.(string))
			}
			request.AuthRole = append(request.AuthRole, &auth)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().CreateAccountUser(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mongodb instanceAccount failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId + tccommon.FILED_SP + userName)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	if response != nil && response.Response != nil {
		if err = service.DescribeAsyncRequestInfo(ctx, helper.UInt64ToStr(*response.Response.FlowId), 3*tccommon.ReadRetryTimeout); err != nil {
			return err
		}
	}

	return resourceTencentCloudMongodbInstanceAccountRead(d, meta)
}

func resourceTencentCloudMongodbInstanceAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_account.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userName := idSplit[1]

	instanceAccount, err := service.DescribeMongodbInstanceAccountById(ctx, instanceId, userName)
	if err != nil {
		return err
	}

	if instanceAccount == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MongodbInstanceAccount` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if instanceAccount.UserName != nil {
		_ = d.Set("user_name", instanceAccount.UserName)
	}

	if instanceAccount.UserDesc != nil {
		_ = d.Set("user_desc", instanceAccount.UserDesc)
	}

	if instanceAccount.AuthRole != nil {
		authRoleList := []interface{}{}
		for _, authRole := range instanceAccount.AuthRole {
			authRoleMap := map[string]interface{}{}

			if authRole.Mask != nil {
				authRoleMap["mask"] = authRole.Mask
			}

			if authRole.NameSpace != nil {
				authRoleMap["namespace"] = authRole.NameSpace
			}

			authRoleList = append(authRoleList, authRoleMap)
		}

		_ = d.Set("auth_role", authRoleList)

	}

	return nil
}

func resourceTencentCloudMongodbInstanceAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_account.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := mongodb.NewSetAccountUserPrivilegeRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userName := idSplit[1]

	request.InstanceId = &instanceId
	request.UserName = &userName

	immutableArgs := []string{"user_desc"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("auth_role") {
		if v, ok := d.GetOk("auth_role"); ok {
			for _, item := range v.([]interface{}) {
				auth := mongodb.Auth{}
				dMap := item.(map[string]interface{})
				if v, ok := dMap["mask"]; ok {
					auth.Mask = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["namespace"]; ok {
					auth.NameSpace = helper.String(v.(string))
				}
				request.AuthRole = append(request.AuthRole, &auth)
			}
		}

		var response *mongodb.SetAccountUserPrivilegeResponse
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMongodbClient().SetAccountUserPrivilege(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mongodb instanceAccount failed, reason:%+v", logId, err)
			return err
		}

		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

		if response != nil && response.Response != nil {
			if err = service.DescribeAsyncRequestInfo(ctx, helper.UInt64ToStr(*response.Response.FlowId), 3*tccommon.ReadRetryTimeout); err != nil {
				return err
			}
		}
	}

	if d.HasChange("password") {
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		password := d.Get("password").(string)
		err := service.ResetInstancePassword(ctx, instanceId, userName, password)
		if err != nil {
			return err
		}

	}

	return resourceTencentCloudMongodbInstanceAccountRead(d, meta)
}

func resourceTencentCloudMongodbInstanceAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mongodb_instance_account.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userName := idSplit[1]

	var mongoUserPassword string
	if v, ok := d.GetOk("mongo_user_password"); ok {
		mongoUserPassword = v.(string)
	}

	if err := service.DeleteMongodbInstanceAccountById(ctx, instanceId, userName, mongoUserPassword); err != nil {
		return err
	}

	return nil
}
