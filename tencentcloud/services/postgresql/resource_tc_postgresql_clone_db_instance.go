package postgresql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresqlCloneDbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlCloneDbInstanceCreate,
		Read:   resourceTencentCloudPostgresqlCloneDbInstanceRead,
		Delete: resourceTencentCloudPostgresqlCloneDbInstanceDelete,
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID original 实例 到 是 cloned。",
			},

			"spec_code": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Purchasable 代码，其中 可以 是 获取 从 `SpecCode` 字段 在 返回值 的 [DescribeClasses](https://intl.云.tencent.com/document/api/409/89019?from_cn_redirect=1) API。",
			},

			"storage": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "实例 存储 容量 （GB）。",
			},

			"period": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Purchase 时长，在 months.\n- Prepaid: Supports `1`，`2`，`3`，`4`，`5`，`6`，`7`，`8`，`9`，`10`，`11`，`12`，`24`，和 `36`.\n- Pay-作为-您-go: Only 支持 `1`。",
			},

			"auto_renew_flag": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Renewal Flag:\n\n- `0`: manual renewal\n`1`: auto-renewal\n\n默认值：0。",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "私有网络 ID 在 格式 的 `vpc-xxxxxxx`，其中 可以 是 获取 在 console 或 从 `unVpcId` 字段 在 返回值 的 [DescribeVpcEx](https://intl.云.tencent.com/document/api/215/1372?from_cn_redirect=1) API。",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC 子网 ID 在 格式 的 `子网-xxxxxxxx`，其中 可以 是 获取 在 console 或 从 `unSubnetId` 字段 在 返回值 的 [DescribeSubnets](https://intl.云.tencent.com/document/api/215/15784?from_cn_redirect=1) API。",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "新购买的实例名称，最多60个字母、数字或符号（-_）。如果不指定该参数，则默认显示“未命名”。",
			},

			"instance_charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "实例 billing 类型，其中 currently 支持:\n\n- PREPAID: Prepaid，i.e.，monthly subscription\n- POSTPAID_BY_HOUR: Pay-作为-您-go，i.e.，pay 通过 consumption\n\n默认值：PREPAID。",
			},

			"security_group_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "Security 组 的 实例，其中 可以 是 获取 从 `sgld` 字段 在 返回值 的 [DescribeSecurityGroups](https://intl.云.tencent.com/document/api/215/15808?from_cn_redirect=1) API. 如果 此 参数 是 不 指定， 默认值 安全 组 将 是 bound。",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "项目 ID",
			},

			"tag_list": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "信息 的 标签 到 是 bound 使用 实例，其中 是 left 空 通过 默认值. 此 参数 可以 是 获取 从 `标签` 字段 在 返回值 的 [DescribeTags](https://intl.云.tencent.com/document/api/651/35316?from_cn_redirect=1) API。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "标签键",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "标签值",
						},
					},
				},
			},

			"db_node_set": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Deployment 信息 的 实例 节点，其中 将 display 信息 的 each AZ 当 实例 节点 是 deployed across 多个 AZs.\nThe 信息 的 AZ 可以 是 获取 从 `可用区` 字段 在 返回值 的 [DescribeZones](https://intl.云.tencent.com/document/api/409/16769?from_cn_redirect=1) API。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Node 类型 有效 值:\n`Primary`;\n`Standby`。",
						},
						"zone": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "AZ 其中 节点 resides，such 作为 ap-guangzhou-1。",
						},
						"dedicated_cluster_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Dedicated 集群 ID。",
						},
					},
				},
			},

			"activity_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Campaign ID。",
			},

			"backup_set_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Basic 备份 集合 ID。",
			},

			"recovery_target_time": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Restoration point 在 时间。",
			},

			"sync_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Primary-standby sync 模式，其中 支持:\nSemi-sync: Semi-sync\nAsync: Asynchronous\nDefault 值 对于 primary 实例: Semi-sync\nDefault 值 对于 read-仅 实例: Async。",
			},
		},
	}
}

func resourceTencentCloudPostgresqlCloneDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_clone_db_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = postgresv20170312.NewCloneDBInstanceRequest()
		response = postgresv20170312.NewCloneDBInstanceResponse()
	)

	if v, ok := d.GetOk("db_instance_id"); ok {
		request.DBInstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("spec_code"); ok {
		request.SpecCode = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("storage"); ok {
		request.Storage = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("period"); ok {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("auto_renew_flag"); ok {
		request.AutoRenewFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, helper.String(securityGroupIds))
		}
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		request.ProjectId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("tag_list"); ok {
		for _, item := range v.([]interface{}) {
			tagListMap := item.(map[string]interface{})
			tag := postgresv20170312.Tag{}
			if v, ok := tagListMap["tag_key"]; ok {
				tag.TagKey = helper.String(v.(string))
			}

			if v, ok := tagListMap["tag_value"]; ok {
				tag.TagValue = helper.String(v.(string))
			}

			request.TagList = append(request.TagList, &tag)
		}
	}

	if v, ok := d.GetOk("db_node_set"); ok {
		for _, item := range v.([]interface{}) {
			dBNodeSetMap := item.(map[string]interface{})
			dBNode := postgresv20170312.DBNode{}
			if v, ok := dBNodeSetMap["role"]; ok {
				dBNode.Role = helper.String(v.(string))
			}

			if v, ok := dBNodeSetMap["zone"]; ok {
				dBNode.Zone = helper.String(v.(string))
			}

			if v, ok := dBNodeSetMap["dedicated_cluster_id"]; ok {
				dBNode.DedicatedClusterId = helper.String(v.(string))
			}

			request.DBNodeSet = append(request.DBNodeSet, &dBNode)
		}
	}

	if v, ok := d.GetOkExists("activity_id"); ok {
		request.ActivityId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("backup_set_id"); ok {
		request.BackupSetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("recovery_target_time"); ok {
		request.RecoveryTargetTime = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sync_mode"); ok {
		request.SyncMode = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresV20170312Client().CloneDBInstanceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create postgresql clone db instance failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create postgresql clone db instance failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.DBInstanceId == nil {
		return fmt.Errorf("DBInstanceId is nil.")
	}

	// wait
	dBInstanceId := *response.Response.DBInstanceId
	if _, err := (&resource.StateChangeConf{
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
		Pending:    []string{},
		Refresh:    resourcePostgresqlCloneDbInstanceCreateStateRefreshFunc_0_0(ctx, dBInstanceId),
		Target:     []string{"running"},
		Timeout:    1800 * time.Second,
	}).WaitForStateContext(ctx); err != nil {
		return err
	}

	d.SetId(dBInstanceId)
	return resourceTencentCloudPostgresqlCloneDbInstanceRead(d, meta)
}

func resourceTencentCloudPostgresqlCloneDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_clone_db_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresqlCloneDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_clone_db_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourcePostgresqlCloneDbInstanceCreateStateRefreshFunc_0_0(ctx context.Context, dBInstanceId string) resource.StateRefreshFunc {
	var req *postgresv20170312.DescribeDBInstanceAttributeRequest
	return func() (interface{}, string, error) {
		meta := tccommon.ProviderMetaFromContext(ctx)
		if meta == nil {
			return nil, "", fmt.Errorf("resource data can not be nil")
		}

		if req == nil {
			d := tccommon.ResourceDataFromContext(ctx)
			if d == nil {
				return nil, "", fmt.Errorf("resource data can not be nil")
			}

			_ = d
			req = postgresv20170312.NewDescribeDBInstanceAttributeRequest()
			req.DBInstanceId = helper.String(dBInstanceId)
		}

		resp, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresV20170312Client().DescribeDBInstanceAttributeWithContext(ctx, req)
		if err != nil {
			return nil, "", err
		}

		if resp == nil || resp.Response == nil {
			return nil, "", nil
		}

		state := fmt.Sprintf("%v", *resp.Response.DBInstance.DBInstanceStatus)
		return resp.Response, state, nil
	}
}
