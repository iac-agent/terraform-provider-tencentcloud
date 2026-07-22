package tat

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTatCommand() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTatCommandCreate,
		Read:   resourceTencentCloudTatCommandRead,
		Update: resourceTencentCloudTatCommandUpdate,
		Delete: resourceTencentCloudTatCommandDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"command_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "命令名称 名称 可以 是 up 到 60 bytes，和 contain [-z]，[A-Z]，[0-9] 和 [_-.]。",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Command 内容 最大 长度 是 64 KB。",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "命令描述 最大 长度 是 120 字符。",
			},

			"command_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "命令类型 `SHELL`，`POWERSHELL` 和 `BAT` 是 支持. 默认值为 `SHELL`。",
			},

			"working_directory": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Command execution 路径 默认值为 /root 对于 `SHELL` commands 和 C:/Program Files/qcloudtat_agent/workdir 对于 `POWERSHELL` commands。",
			},

			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Command 超时 周期 默认值：60 秒. 取值范围：[1，86400]。",
			},

			"enable_parameter": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "是否enable 自定义 参数 功能.此 不能 是 modified once 创建.默认值：`false`。",
			},

			"default_parameters": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "默认值 值 的 自定义 参数 值 当 它 是 已启用. 字段 类型 是 JSON encoded 字符串. For 示例, {\"varA\": \"222\"}.`键` 是 名称 的 自定义 参数 和 值 是 默认值 值. Both `键` 和 `值` 是 strings.如果 无 参数 值 是 提供 在 `InvokeCommand` API, 默认值 值 是 使用.Up 到 20 自定义 参数 是 支持. 名称 的 自定义 参数 不能 exceed 64 字符 和 可以 contain [-z], [A-Z], [0-9] 和 [-_].",
			},

			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "标签 bound 到 command. At most 10 标签 是 allowed。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "标签键",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "标签值",
						},
					},
				},
			},

			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用户名 用于execute command 在 CVM 或 Lighthouse 实例. principle 的 least privilege 是 best practice 对于 权限 management. We recommend 您 execute TAT commands 作为 general 用户 By 默认值， root 用户 是 用于execute commands 在 Linux 和 System 用户 是 使用 在 Windows。",
			},

			"output_cos_bucket_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "COS 存储桶 URL 对于 uploading logs. URL 必须 start 使用 `https`，such 作为 `https://BucketName-123454321.cos.ap-beijing.myqcloud.com`。",
			},

			"output_cos_key_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "COS 存储桶 directory 其中 logs 是 saved. Check below 对于 规则 的 directory 名称1. It 必须 是 combination 的 数量，letters，和 visible 字符. Up 到 60 字符 是 allowed.2. Use slash (/) 到 create subdirectory.3. Consecutive dots (.) 和 slashes (/) 是 不 allowed. It 可以 不 start 使用 slash (/)。",
			},

			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Command 创建时间。",
			},

			"updated_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Command 更新时间。",
			},

			"formatted_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Formatted 描述 command. 此 参数 是 空 字符串 对于 用户 commands 和 包含values 对于 公有 commands。",
			},

			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Command 创建者 `TAT` 表示a 公有 command 和 `USER` 表示a personal command。",
			},
		},
	}
}

func resourceTencentCloudTatCommandCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tat_command.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		request   = tat.NewCreateCommandRequest()
		response  *tat.CreateCommandResponse
		commandId string
	)

	if v, ok := d.GetOk("command_name"); ok {
		request.CommandName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(tccommon.StringToBase64(v.(string)))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("command_type"); ok {
		request.CommandType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("working_directory"); ok {
		request.WorkingDirectory = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("timeout"); ok {
		request.Timeout = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("enable_parameter"); ok {
		request.EnableParameter = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("default_parameters"); ok {
		request.DefaultParameters = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			if dMap, ok := item.(map[string]interface{}); ok && dMap != nil {
				tag := tat.Tag{}
				if v, ok := dMap["key"]; ok {
					tag.Key = helper.String(v.(string))
				}
				if v, ok := dMap["value"]; ok {
					tag.Value = helper.String(v.(string))
				}

				request.Tags = append(request.Tags, &tag)
			}
		}
	}

	if v, ok := d.GetOk("username"); ok {
		request.Username = helper.String(v.(string))
	}

	if v, ok := d.GetOk("output_cos_bucket_url"); ok {
		request.OutputCOSBucketUrl = helper.String(v.(string))
	}

	if v, ok := d.GetOk("output_cos_key_prefix"); ok {
		request.OutputCOSKeyPrefix = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTatClient().CreateCommand(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create tat command failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tat command failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.CommandId == nil {
		return fmt.Errorf("CommandId is nil.")
	}

	commandId = *response.Response.CommandId

	d.SetId(commandId)
	return resourceTencentCloudTatCommandRead(d, meta)
}

func resourceTencentCloudTatCommandRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tat_command.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TatService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	commandId := d.Id()

	command, err := service.DescribeTatCommand(ctx, commandId)

	if err != nil {
		return err
	}

	if command == nil {
		d.SetId("")
		return fmt.Errorf("resource `command` %s does not exist", commandId)
	}

	if command.CommandName != nil {
		_ = d.Set("command_name", command.CommandName)
	}

	if command.Content != nil {
		content, err := tccommon.Base64ToString(*command.Content)
		if err != nil {
			return fmt.Errorf("`Content` [%v] base64 to string failed, err: %v.", *command.Content, err)
		}
		_ = d.Set("content", content)
	}

	if command.Description != nil {
		_ = d.Set("description", command.Description)
	}

	if command.CommandType != nil {
		_ = d.Set("command_type", command.CommandType)
	}

	if command.WorkingDirectory != nil {
		_ = d.Set("working_directory", command.WorkingDirectory)
	}

	if command.Timeout != nil {
		_ = d.Set("timeout", command.Timeout)
	}

	if command.EnableParameter != nil {
		_ = d.Set("enable_parameter", command.EnableParameter)
	}

	if command.DefaultParameters != nil {
		_ = d.Set("default_parameters", command.DefaultParameters)
	}

	if command.Tags != nil {
		tagsList := []interface{}{}
		for _, tags := range command.Tags {
			tagsMap := map[string]interface{}{}
			if tags.Key != nil {
				tagsMap["key"] = tags.Key
			}
			if tags.Value != nil {
				tagsMap["value"] = tags.Value
			}

			tagsList = append(tagsList, tagsMap)
		}
		_ = d.Set("tags", tagsList)
	}

	if command.Username != nil {
		_ = d.Set("username", command.Username)
	}

	if command.OutputCOSBucketUrl != nil {
		_ = d.Set("output_cos_bucket_url", command.OutputCOSBucketUrl)
	}

	if command.OutputCOSKeyPrefix != nil {
		_ = d.Set("output_cos_key_prefix", command.OutputCOSKeyPrefix)
	}

	if command.CreatedTime != nil {
		_ = d.Set("created_time", command.CreatedTime)
	}

	if command.UpdatedTime != nil {
		_ = d.Set("updated_time", command.UpdatedTime)
	}

	if command.FormattedDescription != nil {
		_ = d.Set("formatted_description", command.FormattedDescription)
	}

	if command.CreatedBy != nil {
		_ = d.Set("created_by", command.CreatedBy)
	}

	return nil
}

func resourceTencentCloudTatCommandUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tat_command.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := tat.NewModifyCommandRequest()

	if d.HasChange("enable_parameter") {
		return fmt.Errorf("`enable_parameter` do not support change now.")
	}

	if d.HasChange("tags") {
		return fmt.Errorf("`tags` do not support change now.")
	}

	commandId := d.Id()
	request.CommandId = &commandId

	if d.HasChange("command_name") {
		if v, ok := d.GetOk("command_name"); ok {
			request.CommandName = helper.String(v.(string))
		}
	}

	if d.HasChange("content") {
		if v, ok := d.GetOk("content"); ok {
			request.Content = helper.String(tccommon.StringToBase64(v.(string)))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("command_type") {
		if v, ok := d.GetOk("command_type"); ok {
			request.CommandType = helper.String(v.(string))
		}
	}

	if d.HasChange("working_directory") {
		if v, ok := d.GetOk("working_directory"); ok {
			request.WorkingDirectory = helper.String(v.(string))
		}
	}

	if d.HasChange("timeout") {
		if v, ok := d.GetOk("timeout"); ok {
			request.Timeout = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("default_parameters") {
		if v, ok := d.GetOk("default_parameters"); ok {
			request.DefaultParameters = helper.String(v.(string))
		}
	}

	if d.HasChange("username") {
		if v, ok := d.GetOk("username"); ok {
			request.Username = helper.String(v.(string))
		}
	}

	if d.HasChange("output_cos_bucket_url") {
		if v, ok := d.GetOk("output_cos_bucket_url"); ok {
			request.OutputCOSBucketUrl = helper.String(v.(string))
		}
	}

	if d.HasChange("output_cos_key_prefix") {
		if v, ok := d.GetOk("output_cos_key_prefix"); ok {
			request.OutputCOSKeyPrefix = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTatClient().ModifyCommand(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tat command failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTatCommandRead(d, meta)
}

func resourceTencentCloudTatCommandDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tat_command.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TatService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	commandId := d.Id()

	if err := service.DeleteTatCommandById(ctx, commandId); err != nil {
		return err
	}

	return nil
}
