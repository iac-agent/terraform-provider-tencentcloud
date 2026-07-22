package ssm

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSsmProductSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSsmProductSecretCreate,
		Read:   resourceTencentCloudSsmProductSecretRead,
		Update: resourceTencentCloudSsmProductSecretUpdate,
		Delete: resourceTencentCloudSsmProductSecretDelete,
		Schema: map[string]*schema.Schema{
			"secret_name": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Credential 名称，其中 必须 是 唯一 在 same 地域 It 可以 contain 128 bytes 的 letters，digits，hyphens，和 underscores 和 必须 begin 使用 letter 或 digit。",
			},
			"user_name_prefix": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Prefix 的 用户 账号 名称，其中 是 指定 通过 您 和 可以 contain up 到 8 字符.Supported character sets include:Digits: [0，9].Lowercase letters: [，z].Uppercase letters: [A，Z].Special symbols: underscore. prefix 必须 begin 使用 letter。",
			},
			"product_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "名称 Tencent Cloud 服务 bound 到 credential，such 作为 `Mysql`，`Tdsql-mysql`，`Tdsql_C_Mysql`. 您 可以 使用 dataSource `tencentcloud_ssm_products` 到 查询 支持 products。",
			},
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Tencent Cloud 服务 实例 ID。",
			},
			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "描述，其中 是 用于describe purpose 在 detail 和 可以 contain up 到 2,048 bytes。",
			},
			"kms_key_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "指定KMS CMK 该 encrypts credential. 如果此参数为空， CMK 创建 通过 Secrets Manager 通过 默认值 将 是 用于encryption.You 可以 also 指定a 自定义 KMS CMK 创建 在 same 地域 对于 加密。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 的 secret。",
			},
			"domains": {
				Required:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "域名 名称 账号 在 form 的 IP. You 可以 enter `%`。",
			},
			"privileges_list": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "列表 permissions 该 need 到 是 granted 当 credential 是 bound 到 Tencent Cloud 服务。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"privilege_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Permission 名称 有效值：`GlobalPrivileges`，`DatabasePrivileges`，`TablePrivileges`，`ColumnPrivileges`. 当 权限 是 `DatabasePrivileges`， 数据库 名称 必须 是 指定 通过 `Database` 参数; 当 权限 是 `TablePrivileges`， 数据库 名称 和 表 名称 在 数据库 必须 是 指定 通过 `Database` 和 `TableName` 参数; 当 权限 是 `ColumnPrivileges`， 数据库 名称，表 名称 在 数据库，和 列 名称 在 表 必须 是 指定 通过 `Database`，`TableName`，和 `ColumnName` 参数。",
						},
						"privileges": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Permission 列表. For `Mysql` 服务，可选 权限 值 是: 1. 有效 值 的 `GlobalPrivileges`: SELECT,INSERT,UPDATE,DELETE,CREATE，PROCESS，DROP,REFERENCES,INDEX,ALTER,SHOW DATABASES,CREATE TEMPORARY TABLES,LOCK TABLES,EXECUTE,CREATE VIEW,SHOW VIEW,CREATE ROUTINE,ALTER ROUTINE,EVENT,TRIGGER. 注意: 如果 此 参数 是 不 passed 在，它 表示 到 clear 权限. 2. 有效 值 的 `DatabasePrivileges`: SELECT,INSERT,UPDATE,DELETE,CREATE，DROP,REFERENCES,INDEX,ALTER,CREATE TEMPORARY TABLES,LOCK TABLES,EXECUTE,CREATE VIEW,SHOW VIEW,CREATE ROUTINE,ALTER ROUTINE,EVENT,TRIGGER. 注意: 如果 此 参数 是 不 passed 在，它 表示 到 clear 权限. 3. 有效 值 的 `TablePrivileges`: SELECT,INSERT,UPDATE,DELETE,CREATE，DROP,REFERENCES,INDEX,ALTER,CREATE VIEW,SHOW VIEW，TRIGGER. 注意: 如果 此 参数 是 不 passed 在，它 表示 到 clear 权限. 4. 有效 值 的 `ColumnPrivileges`: SELECT,INSERT,UPDATE,REFERENCES.注意: 如果 此 参数 是 不 passed 在，它 表示 到 clear 权限。",
						},
						"database": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "此 值 takes effect 仅 当 `PrivilegeName` 是 `DatabasePrivileges`。",
						},
						"table_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "此 值 takes effect 仅 当 `PrivilegeName` 是 `TablePrivileges`，和 `Database` 参数 为必填项 在 此 case 到 explicitly indicate 数据库 实例。",
						},
						"column_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "此 值 takes effect 仅 当 `PrivilegeName` 是 `ColumnPrivileges`，和 following 参数 为必填项 在 此 case:Database: explicitly indicate 数据库 实例.TableName: explicitly indicate 表。",
						},
					},
				},
			},
			"rotation_begin_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Computed:    true,
				Description: "用户-Defined rotation 开始时间 在 格式 的 2006-01-02 15:04:05.当 `EnableRotation` 是 `True`，此 参数 为必填项。",
			},
			"enable_rotation": {
				Optional:    true,
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "指定是否enable rotation，当 secret 状态 是 `已禁用`，rotation 将 是 已禁用 `True` - 启用，`False` - do 不 启用. 如果 此 参数 是 不 指定，`False` 将 是 使用 通过 默认值。",
			},
			"rotation_frequency": {
				Optional:    true,
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Rotation 频率 在 days. 默认值：1 day。",
			},
			"status": {
				Optional:     true,
				Type:         schema.TypeString,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"Enabled", "Disabled"}),
				Description:  "Enable 或 Disable Secret. 有效 值 是 `已启用` 或 `已禁用`. 默认为 `已启用`。",
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Credential 创建时间 在 UNIX 时间戳 格式",
			},
			"secret_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "`0`: 用户-defined secret. `1`: Tencent Cloud services secret. `2`: SSH 键 secret. `3`: Tencent Cloud API 键 secret. 注意: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
			},
		},
	}
}

func resourceTencentCloudSsmProductSecretCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_product_secret.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request    = ssm.NewCreateProductSecretRequest()
		response   = ssm.NewCreateProductSecretResponse()
		secretInfo *SecretInfo
		secretName string
	)

	if v, ok := d.GetOk("secret_name"); ok {
		request.SecretName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_name_prefix"); ok {
		request.UserNamePrefix = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product_name"); ok {
		request.ProductName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceID = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domains"); ok {
		domainsSet := v.(*schema.Set).List()
		for i := range domainsSet {
			domains := domainsSet[i].(string)
			request.Domains = append(request.Domains, &domains)
		}
	}

	if v, ok := d.GetOk("privileges_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			productPrivilegeUnit := ssm.ProductPrivilegeUnit{}
			if v, ok := dMap["privilege_name"]; ok {
				productPrivilegeUnit.PrivilegeName = helper.String(v.(string))
			}
			if v, ok := dMap["privileges"]; ok {
				privilegesSet := v.(*schema.Set).List()
				for i := range privilegesSet {
					privileges := privilegesSet[i].(string)
					productPrivilegeUnit.Privileges = append(productPrivilegeUnit.Privileges, &privileges)
				}
			}
			if v, ok := dMap["database"]; ok {
				productPrivilegeUnit.Database = helper.String(v.(string))
			}
			if v, ok := dMap["table_name"]; ok {
				productPrivilegeUnit.TableName = helper.String(v.(string))
			}
			if v, ok := dMap["column_name"]; ok {
				productPrivilegeUnit.ColumnName = helper.String(v.(string))
			}
			request.PrivilegesList = append(request.PrivilegesList, &productPrivilegeUnit)
		}
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		request.KmsKeyId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rotation_begin_time"); ok {
		request.RotationBeginTime = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("enable_rotation"); ok {
		request.EnableRotation = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("rotation_frequency"); ok {
		request.RotationFrequency = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSsmClient().CreateProductSecret(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create ssm productSecret failed, reason:%+v", logId, err)
		return err
	}

	secretName = *response.Response.SecretName
	d.SetId(secretName)
	flowId := *response.Response.FlowID
	conf := tccommon.BuildStateChangeConf([]string{}, []string{"1"}, tccommon.ReadRetryTimeout, time.Second, service.SsmProductSecretStateRefreshFunc(flowId, []string{"0"}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	// update status if disabled
	if v, ok := d.GetOk("status"); ok {
		status := v.(string)
		if status == "Disabled" {
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			err = service.DisableSecret(ctx, secretName)
			if err != nil {
				return err
			}
		}
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			secretInfo, err = service.DescribeSecretByName(ctx, secretName)
			if err != nil {
				return tccommon.RetryError(err)
			}

			return nil
		})

		if err != nil {
			return err
		}

		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("ssm", "secret", tcClient.Region, secretInfo.resourceId)
		if err = tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudSsmProductSecretRead(d, meta)
}

func resourceTencentCloudSsmProductSecretRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_product_secret.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		secretInfo *SecretInfo
		secretName = d.Id()
	)

	productSecret, err := service.DescribeSecretById(ctx, secretName, 1)
	if err != nil {
		return err
	}

	if productSecret == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SsmProductSecret` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if productSecret.SecretName != nil {
		_ = d.Set("secret_name", productSecret.SecretName)
	}

	if productSecret.ProductName != nil {
		_ = d.Set("product_name", productSecret.ProductName)
	}

	if productSecret.ResourceID != nil {
		_ = d.Set("instance_id", productSecret.ResourceID)
	}

	if productSecret.Description != nil {
		_ = d.Set("description", productSecret.Description)
	}

	if productSecret.KmsKeyId != nil {
		_ = d.Set("kms_key_id", productSecret.KmsKeyId)
	}

	if productSecret.RotationBeginTime != nil {
		_ = d.Set("rotation_begin_time", productSecret.RotationBeginTime)
	}

	if productSecret.RotationStatus != nil {
		_ = d.Set("enable_rotation", helper.Bool(true))
		if *productSecret.RotationStatus == 0 {
			_ = d.Set("enable_rotation", helper.Bool(false))
		}
	}

	if productSecret.RotationFrequency != nil {
		_ = d.Set("rotation_frequency", productSecret.RotationFrequency)
	}

	if productSecret.CreateTime != nil {
		_ = d.Set("create_time", productSecret.CreateTime)
	}

	if productSecret.SecretType != nil {
		_ = d.Set("secret_type", productSecret.SecretType)
	}

	outErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		secretInfo, err = service.DescribeSecretByName(ctx, secretName)
		if err != nil {
			return tccommon.RetryError(err)
		}

		return nil
	})

	if outErr != nil {
		return outErr
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "ssm", "secret", tcClient.Region, secretInfo.resourceId)
	if err != nil {
		return err
	}

	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudSsmProductSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_product_secret.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		ssmService = SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		secretName = d.Id()
	)

	immutableArgs := []string{
		"user_name_prefix", "product_name", "instance_id",
		"domains", "privileges_list", "kms_key_id",
	}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("description") {
		request := ssm.NewUpdateDescriptionRequest()
		request.SecretName = &secretName

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSsmClient().UpdateDescription(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update ssm productSecret failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("status") {
		service := SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

		if v, ok := d.GetOk("status"); ok {
			status := v.(string)
			if status == "Disabled" {
				err := service.DisableSecret(ctx, secretName)
				if err != nil {
					return err
				}
			} else {
				err := service.EnableSecret(ctx, secretName)
				if err != nil {
					return err
				}
			}
		}
	}

	if d.HasChange("enable_rotation") || d.HasChange("rotation_begin_time") || d.HasChange("rotation_frequency") {
		request := ssm.NewUpdateRotationStatusRequest()
		request.SecretName = &secretName

		if v, ok := d.GetOk("rotation_begin_time"); ok {
			request.RotationBeginTime = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("enable_rotation"); ok {
			request.EnableRotation = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("rotation_frequency"); ok {
			request.Frequency = helper.IntInt64(v.(int))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSsmClient().UpdateRotationStatus(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update ssm productSecret failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)

		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))
		secretInfo, err := ssmService.DescribeSecretByName(ctx, secretName)
		if err != nil {
			return err
		}

		resourceName := tccommon.BuildTagResourceName("ssm", "secret", tcClient.Region, secretInfo.resourceId)
		if err = tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

	}

	return resourceTencentCloudSsmProductSecretRead(d, meta)
}

func resourceTencentCloudSsmProductSecretDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssm_product_secret.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		secretName = d.Id()
	)

	// disable before destroy
	err := service.DisableSecret(ctx, secretName)
	if err != nil {
		return err
	}

	if err = service.DeleteSsmProductSecretById(ctx, secretName); err != nil {
		return err
	}

	return nil
}
