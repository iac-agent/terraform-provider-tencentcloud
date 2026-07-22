package crs

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudRedisInstanceZoneInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRedisInstanceZoneInfoRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ID 实例。",
			},

			"replica_groups": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 实例 节点 groups。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node 组 ID",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node 组 名称",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "he availability 可用区 ID 节点，such 作为 ap-guangzhou-1。",
						},
						"role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "节点 组 类型，master 是 primary 节点，和 副本 是 副本 节点。",
						},
						"redis_nodes": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Node 组 节点 列表。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"keys": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 节点 keys。",
									},
									"slot": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Node slot distribution。",
									},
									"node_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "节点 ID",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Node 状态",
									},
									"role": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Node 角色",
									},
								},
							},
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

func dataSourceTencentCloudRedisInstanceZoneInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_redis_instance_zone_info.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var instanceId string

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var replicaGroups []*redis.ReplicaGroup

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRedisInstanceZoneInfoByFilter(ctx, paramMap)
		if e != nil {
			if ee, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if ee.Code == "FailedOperation.UnSupportError" {
					return resource.NonRetryableError(e)
				}
			}
			return tccommon.RetryError(e)
		}
		replicaGroups = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(replicaGroups))

	if replicaGroups != nil {
		for _, replicaGroup := range replicaGroups {
			replicaGroupMap := map[string]interface{}{}

			if replicaGroup.GroupId != nil {
				replicaGroupMap["group_id"] = replicaGroup.GroupId
			}

			if replicaGroup.GroupName != nil {
				replicaGroupMap["group_name"] = replicaGroup.GroupName
			}

			if replicaGroup.ZoneId != nil {
				replicaGroupMap["zone_id"] = replicaGroup.ZoneId
			}

			if replicaGroup.Role != nil {
				replicaGroupMap["role"] = replicaGroup.Role
			}

			if replicaGroup.RedisNodes != nil {
				redisNodesList := []interface{}{}
				for _, redisNodes := range replicaGroup.RedisNodes {
					redisNodesMap := map[string]interface{}{}

					if redisNodes.Keys != nil {
						redisNodesMap["keys"] = redisNodes.Keys
					}

					if redisNodes.Slot != nil {
						redisNodesMap["slot"] = redisNodes.Slot
					}

					if redisNodes.NodeId != nil {
						redisNodesMap["node_id"] = redisNodes.NodeId
					}

					if redisNodes.Status != nil {
						redisNodesMap["status"] = redisNodes.Status
					}

					if redisNodes.Role != nil {
						redisNodesMap["role"] = redisNodes.Role
					}

					redisNodesList = append(redisNodesList, redisNodesMap)
				}

				replicaGroupMap["redis_nodes"] = redisNodesList
			}

			tmpList = append(tmpList, replicaGroupMap)
		}

		_ = d.Set("replica_groups", tmpList)
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
