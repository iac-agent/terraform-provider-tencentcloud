package tse

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTseNacosReplicas() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTseNacosReplicasRead,
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
							Description: "Available area Name注意：此字段可能返回 null，表示有效值不可用。",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Available area ID注意：此字段可能返回 null，表示有效值不可用。",
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

func dataSourceTencentCloudTseNacosReplicasRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tse_nacos_replicas.read")()
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

	var replicas []*tse.NacosReplica

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTseNacosReplicasByFilter(ctx, paramMap)
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
		for _, nacosReplica := range replicas {
			nacosReplicaMap := map[string]interface{}{}

			if nacosReplica.Name != nil {
				nacosReplicaMap["name"] = nacosReplica.Name
			}

			if nacosReplica.Role != nil {
				nacosReplicaMap["role"] = nacosReplica.Role
			}

			if nacosReplica.Status != nil {
				nacosReplicaMap["status"] = nacosReplica.Status
			}

			if nacosReplica.SubnetId != nil {
				nacosReplicaMap["subnet_id"] = nacosReplica.SubnetId
			}

			if nacosReplica.Zone != nil {
				nacosReplicaMap["zone"] = nacosReplica.Zone
			}

			if nacosReplica.ZoneId != nil {
				nacosReplicaMap["zone_id"] = nacosReplica.ZoneId
			}

			if nacosReplica.VpcId != nil {
				nacosReplicaMap["vpc_id"] = nacosReplica.VpcId
			}

			tmpList = append(tmpList, nacosReplicaMap)
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
