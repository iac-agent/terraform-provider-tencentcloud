package ckafka

import (
	"context"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCkafkaAcls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaAclsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the ckafka instance。",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ACL 资源类型 Valid values are `UNKNOWN`，`ANY`，`TOPIC`，`GROUP`，`CLUSTER`，`TRANSACTIONAL_ID`. Currently，only `TOPIC` is available，and other fields will be 用于future ACLs compatible with open-来源 Kafka。",
			},
			"resource_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ACL 资源名称，which is related to `resource_type`. For example，if `resource_type` is `TOPIC`，this field 表示topic 名称; if `resource_type` is `GROUP`，this field 表示group 名称",
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "主机 substr 用于querying。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"acl_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 ckafka acls. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ACL 资源类型",
						},
						"resource_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ACL 资源名称，which is related to `resource_type`。",
						},
						"operation_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ACL operation 模式",
						},
						"permission_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ACL permission 类型，valid values are `UNKNOWN`，`ANY`，`DENY`，`ALLOW`，and `ALLOW` by default. Currently，CKafka supports `ALLOW` (equivalent to allow list)，and other fields will be 用于future ACLs compatible with open-来源 Kafka。",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP 地址 allowed to access。",
						},
						"principal": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "用户 which can access. `*` means that any 用户 can access。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCkafkaAclsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ckafka_acls.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	params := make(map[string]interface{})
	params["instance_id"] = d.Get("instance_id").(string)
	params["resource_type"] = d.Get("resource_type").(string)
	params["resource_name"] = d.Get("resource_name").(string)
	if v, ok := d.GetOk("host"); ok {
		params["host"] = v.(string)
	}

	ckafkaService := CkafkaService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	aclInfos, err := ckafkaService.DescribeAclByFilter(ctx, params)
	if err != nil {
		return err
	}
	aclList := make([]map[string]interface{}, 0, len(aclInfos))
	ids := make([]string, 0, len(aclInfos))
	for _, acl := range aclInfos {
		aclList = append(aclList, map[string]interface{}{
			"resource_type":   CKAFKA_ACL_RESOURCE_TYPE_TO_STRING[*acl.ResourceType],
			"resource_name":   *acl.ResourceName,
			"operation_type":  CKAFKA_ACL_OPERATION_TO_STRING[*acl.Operation],
			"permission_type": CKAFKA_PERMISSION_TYPE_TO_STRING[*acl.PermissionType],
			"host":            *acl.Host,
			"principal":       strings.TrimLeft(*acl.Principal, CKAFKA_ACL_PRINCIPAL_STR),
		})

		ids = append(ids, params["instance_id"].(string)+tccommon.FILED_SP+CKAFKA_PERMISSION_TYPE_TO_STRING[*acl.PermissionType]+
			tccommon.FILED_SP+strings.TrimLeft(*acl.Principal, CKAFKA_ACL_PRINCIPAL_STR)+tccommon.FILED_SP+*acl.Host+tccommon.FILED_SP+
			CKAFKA_ACL_OPERATION_TO_STRING[*acl.Operation]+tccommon.FILED_SP+CKAFKA_ACL_RESOURCE_TYPE_TO_STRING[*acl.ResourceType]+
			tccommon.FILED_SP+*acl.ResourceName)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("acl_list", aclList)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), aclList); e != nil {
			return e
		}
	}

	return nil
}
