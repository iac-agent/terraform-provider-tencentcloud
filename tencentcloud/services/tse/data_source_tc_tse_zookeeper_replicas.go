package tse

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTseZookeeperReplicas() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTseZookeeperReplicasRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "engine instance ID。",
			},

			"replicas": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Engine instance replica information。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称",
						},
						"role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "角色",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet ID注意：此字段可能返回 null，表示有效值不可用。",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Available area ID注意：此字段可能返回 null，表示有效值不可用。",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Available area ID注意：此字段可能返回 null，表示有效值不可用。",
						},
						"alias_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "alias注意：此字段可能返回 null，表示有效值不可用。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC ID注意：此字段可能返回 null，表示有效值不可用。",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudTseZookeeperReplicasRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tse_zookeeper_replicas.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	instanceId := ""

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := TseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var replicas []*tse.ZookeeperReplica

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTseZookeeperReplicasByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		replicas = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(replicas))

	if replicas != nil {
		for _, zookeeperReplica := range replicas {
			zookeeperReplicaMap := map[string]interface{}{}

			if zookeeperReplica.Name != nil {
				zookeeperReplicaMap["name"] = zookeeperReplica.Name
			}

			if zookeeperReplica.Role != nil {
				zookeeperReplicaMap["role"] = zookeeperReplica.Role
			}

			if zookeeperReplica.Status != nil {
				zookeeperReplicaMap["status"] = zookeeperReplica.Status
			}

			if zookeeperReplica.SubnetId != nil {
				zookeeperReplicaMap["subnet_id"] = zookeeperReplica.SubnetId
			}

			if zookeeperReplica.Zone != nil {
				zookeeperReplicaMap["zone"] = zookeeperReplica.Zone
			}

			if zookeeperReplica.ZoneId != nil {
				zookeeperReplicaMap["zone_id"] = zookeeperReplica.ZoneId
			}

			if zookeeperReplica.AliasName != nil {
				zookeeperReplicaMap["alias_name"] = zookeeperReplica.AliasName
			}

			if zookeeperReplica.VpcId != nil {
				zookeeperReplicaMap["vpc_id"] = zookeeperReplica.VpcId
			}

			tmpList = append(tmpList, zookeeperReplicaMap)
		}

		_ = d.Set("replicas", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash([]string{instanceId}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
