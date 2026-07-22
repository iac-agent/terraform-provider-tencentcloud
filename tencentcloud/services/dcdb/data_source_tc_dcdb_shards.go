package dcdb

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDcdbShards() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbShardsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "实例 ID",
			},

			"shard_instance_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "分片 实例 ids。",
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "分片 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID",
						},
						"shard_serial_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "分片 serial ID。",
						},
						"shard_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "分片 实例 ID",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "状态",
						},
						"status_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态 描述",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "私有网络 ID",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "子网 ID",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "项目 ID",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用区",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "内存， 单位 是 GB。",
						},
						"storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "内存， 单位 是 GB。",
						},
						"period_end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "过期时间。",
						},
						"node_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "节点 count。",
						},
						"storage_usage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "存储 usage。",
						},
						"memory_usage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "内存 usage。",
						},
						"proxy_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "proxy 版本",
						},
						"paymode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "付费模式",
						},
						"shard_master_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "分片 master 可用区",
						},
						"shard_slave_zones": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "分片 slave zones。",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CPU 核数",
						},
						"range": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "范围 的 分片 键",
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

func dataSourceTencentCloudDcdbShardsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dcdb_shards.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("shard_instance_ids"); ok {
		shard_instance_idsSet := v.(*schema.Set).List()
		ids := make([]*string, 0, len(shard_instance_idsSet))
		for _, vv := range shard_instance_idsSet {
			ids = append(ids, helper.String(vv.(string)))
		}
		paramMap["shard_instance_ids"] = ids
	}

	dcdbService := DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var shards []*dcdb.DCDBShardInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := dcdbService.DescribeDcdbShardsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		shards = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Dcdb shards failed, reason:%+v", logId, err)
		return err
	}

	// shardList := []interface{}{}
	shardList := make([]map[string]interface{}, 0, len(shards))
	ids := make([]string, 0, len(shards))
	if shards != nil {
		for _, shard := range shards {
			shardMap := map[string]interface{}{}
			if shard.InstanceId != nil {
				shardMap["instance_id"] = shard.InstanceId
			}
			if shard.ShardSerialId != nil {
				shardMap["shard_serial_id"] = shard.ShardSerialId
			}
			if shard.ShardInstanceId != nil {
				shardMap["shard_instance_id"] = shard.ShardInstanceId
			}
			if shard.Status != nil {
				shardMap["status"] = shard.Status
			}
			if shard.StatusDesc != nil {
				shardMap["status_desc"] = shard.StatusDesc
			}
			if shard.CreateTime != nil {
				shardMap["create_time"] = shard.CreateTime
			}
			if shard.VpcId != nil {
				shardMap["vpc_id"] = shard.VpcId
			}
			if shard.SubnetId != nil {
				shardMap["subnet_id"] = shard.SubnetId
			}
			if shard.ProjectId != nil {
				shardMap["project_id"] = shard.ProjectId
			}
			if shard.Region != nil {
				shardMap["region"] = shard.Region
			}
			if shard.Zone != nil {
				shardMap["zone"] = shard.Zone
			}
			if shard.Memory != nil {
				shardMap["memory"] = shard.Memory
			}
			if shard.Storage != nil {
				shardMap["storage"] = shard.Storage
			}
			if shard.PeriodEndTime != nil {
				shardMap["period_end_time"] = shard.PeriodEndTime
			}
			if shard.NodeCount != nil {
				shardMap["node_count"] = shard.NodeCount
			}
			if shard.StorageUsage != nil {
				shardMap["storage_usage"] = shard.StorageUsage
			}
			if shard.MemoryUsage != nil {
				shardMap["memory_usage"] = shard.MemoryUsage
			}
			if shard.ProxyVersion != nil {
				shardMap["proxy_version"] = shard.ProxyVersion
			}
			if shard.Paymode != nil {
				shardMap["paymode"] = shard.Paymode
			}
			if shard.ShardMasterZone != nil {
				shardMap["shard_master_zone"] = shard.ShardMasterZone
			}
			if shard.ShardSlaveZones != nil {
				shardMap["shard_slave_zones"] = shard.ShardSlaveZones
			}
			if shard.Cpu != nil {
				shardMap["cpu"] = shard.Cpu
			}
			if shard.Range != nil {
				shardMap["range"] = shard.Range
			}

			ids = append(ids, *shard.ShardInstanceId)
			shardList = append(shardList, shardMap)
		}
		d.SetId(helper.DataResourceIdsHash(ids))
		err = d.Set("list", shardList)
		if err != nil {
			log.Printf("[CRITAL]%s set Dcdb shards failed, reason:%+v", logId, err)
			return err
		}
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), shardList); e != nil {
			return e
		}
	}

	return nil
}
