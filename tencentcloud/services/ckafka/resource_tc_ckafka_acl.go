package ckafka

import (
	"context"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudCkafkaAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCkafkaAclCreate,
		Read:   resourceTencentCloudCkafkaAclRead,
		Delete: resourceTencentCloudCkafkaAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID ckafka 实例。",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "TOPIC",
				ForceNew:    true,
				Description: "ACL 资源类型 有效 值 是 `UNKNOWN`，`ANY`，`TOPIC`，`GROUP`，`CLUSTER`，`TRANSACTIONAL_ID`. 和 `TOPIC` 通过 默认值. Currently，仅 `TOPIC` 是 可用，和 other 字段 将 是 用于future ACLs compatible 使用 open-来源 Kafka。",
			},
			"resource_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ACL 资源名称，其中 是 related 到 `resource_type`. For 示例，如果 `resource_type` 是 `TOPIC`，此 字段 表示topic 名称; 如果 `resource_type` 是 `GROUP`，此 字段 表示group 名称",
			},
			"operation_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ACL operation 模式 有效值：`UNKNOWN`，`ANY`，`ALL`，`READ`，`WRITE`，`CREATE`，`DELETE`，`ALTER`，`DESCRIBE`，`CLUSTER_ACTION`，`DESCRIBE_CONFIGS` 和 `ALTER_CONFIGS`。",
			},
			"permission_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ALLOW",
				ForceNew:    true,
				Description: "ACL 权限 类型 有效值：`UNKNOWN`，`ANY`，`DENY`，`ALLOW`. 和 `ALLOW` 通过 默认值. Currently，CKafka 支持 `ALLOW` (equivalent 到 allow 列表)，和 other 字段 将 是 用于future ACLs compatible 使用 open-来源 Kafka。",
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "*",
				ForceNew:    true,
				Description: "默认为 *，其中 表示 该 any 主机 可以 访问 它. Support filling 在 IP 或 网络 segment，和 support `;`separation。",
			},
			"principal": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "*",
				ForceNew:    true,
				Description: "用户 列表. 默认值为 `*`，其中 表示 该 any 用户 可以 访问. 当前 用户 可以 仅 是 一个 included 在 用户 列表. For 示例: `root` meaning 用户 root 可以 访问。",
			},
		},
	}
}

func resourceTencentCloudCkafkaAclCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_acl.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Get("instance_id").(string)
	resourceType := d.Get("resource_type").(string)
	resourceName := d.Get("resource_name").(string)
	operation := d.Get("operation_type").(string)
	permissionType := d.Get("permission_type").(string)
	host := d.Get("host").(string)
	principal := d.Get("principal").(string)

	ckafkaService := CkafkaService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	if err := ckafkaService.CreateAcl(ctx, instanceId, resourceType, resourceName, operation, permissionType, host, principal); err != nil {
		return fmt.Errorf("[CRITAL]%s create ckafka user failed, reason:%+v", logId, err)
	}
	d.SetId(instanceId + tccommon.FILED_SP + permissionType + tccommon.FILED_SP + principal + tccommon.FILED_SP + host + tccommon.FILED_SP + operation + tccommon.FILED_SP + resourceType + tccommon.FILED_SP + resourceName)

	return resourceTencentCloudCkafkaAclRead(d, meta)
}

func resourceTencentCloudCkafkaAclRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_acl.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	ckafkaService := CkafkaService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	id := d.Id()
	info, has, err := ckafkaService.DescribeAclByAclId(ctx, id)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}
	items := strings.Split(id, tccommon.FILED_SP)
	_ = d.Set("instance_id", items[0])
	_ = d.Set("resource_type", CKAFKA_ACL_RESOURCE_TYPE_TO_STRING[*info.ResourceType])
	_ = d.Set("resource_name", info.ResourceName)
	_ = d.Set("operation_type", CKAFKA_ACL_OPERATION_TO_STRING[*info.Operation])
	_ = d.Set("permission_type", CKAFKA_PERMISSION_TYPE_TO_STRING[*info.PermissionType])
	_ = d.Set("host", info.Host)
	_ = d.Set("principal", strings.TrimPrefix(*info.Principal, CKAFKA_ACL_PRINCIPAL_STR))

	return nil
}

func resourceTencentCloudCkafkaAclDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ckafka_user.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	ckafkaService := CkafkaService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	if err := ckafkaService.DeleteAcl(ctx, d.Id()); err != nil {
		return err
	}
	return nil
}
