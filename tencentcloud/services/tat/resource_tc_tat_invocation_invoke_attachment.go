package tat

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTatInvocationInvokeAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTatInvocationInvokeAttachmentCreate,
		Read:   resourceTencentCloudTatInvocationInvokeAttachmentRead,
		Delete: resourceTencentCloudTatInvocationInvokeAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID 实例 about 到 execute commands. Supported 实例 types: CVM LIGHTHOUSE。",
			},

			"working_directory": {
				Optional:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Default:     "/root",
				Description: "Command execution 路径 默认值为 /root 对于 SHELL commands 和 C:Program Filesqcloudtat_agentworkdir 对于 POWERSHELL commands。",
			},

			"timeout": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Command 超时 周期 默认值：60 秒. 取值范围：[1，86400]。",
			},

			"parameters": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Custom 参数 的 Command. 字段 类型 是 JSON encoded 字符串. For 示例，{varA: 222}.键 是 名称 自定义 参数 和 值 是 默认值 Both 键 和 值 是 strings.如果 无 参数 值 是 提供， DefaultParameters 是 使用.Up 到 20 自定义 参数 是 支持. 名称 自定义 参数 不能 exceed 64 字符 和 可以 contain [-z]，[A-Z]，[0-9] 和 [-_]。",
			},

			"username": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "用户名 用于execute command 在 CVM 或 Lighthouse 实例. principle 的 least privilege 是 best practice 对于 权限 management. We recommend 您 execute TAT commands 作为 general 用户 By 默认值， 用户 root 是 用于execute commands 在 Linux 和 用户 System 是 使用 在 Windows。",
			},

			"output_cos_bucket_url": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "COS 存储桶 URL 对于 uploading logs. URL 必须 start 使用 https，such 作为 https://BucketName-123454321.cos.ap-beijing.myqcloud.com。",
			},

			"output_cos_key_prefix": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "COS 存储桶 directory 其中 logs 是 saved; Check below 对于 规则 的 directory 名称: 1 It 必须 是 combination 的 数量，letters，和 visible 字符，Up 到 60 字符 是 allowed; 2 Use slash (/) 到 create subdirectory; 3 可以 不 是 使用 作为 文件夹 名称; It 不能 start 使用 slash (/)，和 不能 contain consecutive slashes。",
			},

			"command_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "命令 ID",
			},
		},
	}
}

func resourceTencentCloudTatInvocationInvokeAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tat_invocation_invoke_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request      = tat.NewInvokeCommandRequest()
		response     = tat.NewInvokeCommandResponse()
		invocationId string
		instanceId   string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceIds = []*string{helper.String(v.(string))}
	}

	if v, ok := d.GetOk("working_directory"); ok {
		request.WorkingDirectory = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("timeout"); ok {
		request.Timeout = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("parameters"); ok {
		request.Parameters = helper.String(v.(string))
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

	if v, ok := d.GetOk("command_id"); ok {
		request.CommandId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTatClient().InvokeCommand(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tat invocation failed, reason:%+v", logId, err)
		return err
	}

	invocationId = *response.Response.InvocationId
	d.SetId(invocationId + tccommon.FILED_SP + instanceId)

	return resourceTencentCloudTatInvocationInvokeAttachmentRead(d, meta)
}

func resourceTencentCloudTatInvocationInvokeAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tat_invocation_invoke_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TatService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	invocationId := idSplit[0]
	instanceId := idSplit[1]

	invocation, err := service.DescribeTatInvocationById(ctx, invocationId)
	if err != nil {
		return err
	}

	if invocation == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TatInvocation` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if invocation.WorkingDirectory != nil {
		_ = d.Set("working_directory", invocation.WorkingDirectory)
	}

	if invocation.Timeout != nil {
		_ = d.Set("timeout", invocation.Timeout)
	}

	if invocation.Parameters != nil {
		_ = d.Set("parameters", invocation.Parameters)
	}

	if invocation.Username != nil {
		_ = d.Set("username", invocation.Username)
	}

	if invocation.OutputCOSBucketUrl != nil {
		_ = d.Set("output_cos_bucket_url", invocation.OutputCOSBucketUrl)
	}

	if invocation.OutputCOSKeyPrefix != nil {
		_ = d.Set("output_cos_key_prefix", invocation.OutputCOSKeyPrefix)
	}

	if invocation.CommandId != nil {
		_ = d.Set("command_id", invocation.CommandId)
	}

	return nil
}

func resourceTencentCloudTatInvocationInvokeAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tat_invocation_invoke_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := TatService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	invocationId := idSplit[0]
	instanceId := idSplit[1]

	if err := service.DeleteTatInvocationById(ctx, invocationId, instanceId); err != nil {
		return err
	}

	return nil
}
