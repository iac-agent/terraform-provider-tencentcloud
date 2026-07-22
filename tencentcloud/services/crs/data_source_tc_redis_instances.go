package crs

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudRedisInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentRedisInstancesRead,
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID 可用 可用区",
			},
			"search_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "键 words 用于match results，和 键 words 可以 是: 实例 ID，实例名称 和 IP 地址",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID 项目 到 其中 redis 实例 belongs。",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "数量 limitation 的 results 对于 查询。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 的 redis 实例。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// Computed values
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 redis 实例. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"redis_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID redis 实例。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 redis 实例。",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Available 可用区 到 其中 redis 实例 belongs。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID 项目 到 其中 redis 实例 belongs。",
						},
						"type_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例类型 Refer 到 `数据.tencentcloud_redis_zone_config.列表.type_id` get 可用 值。",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Deprecated:  "It has been deprecated from version 1.33.1. Please use 'type_id' instead.",
							Description: "实例类型 可用值：`master_slave_redis`，`master_slave_ckv`，`cluster_ckv`，`cluster_redis` 和 `standalone_redis`。",
						},
						"redis_shard_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 实例 分片。",
						},
						"redis_replicas_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 实例 copies。",
						},
						"mem_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory 大小 （MB）。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current 状态 实例，maybe: `init`，`processing`，`online`，`isolate` 和 `todelete`。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID vpc 使用 其中 实例 是 associated。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID vpc 子网。",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP 地址 的 实例。",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "端口 用于access redis 实例。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "时间 当 实例 是 创建。",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "标签 的 实例。",
						},
						// payment
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "charge 类型 实例. 有效 值 是 `POSTPAID` 和 `PREPAID`。",
						},
						"node_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "列表 实例 节点 信息. Currently，信息 about 节点 类型 (master 或 副本) 和 节点 availability 可用区 可以 是 passed 在。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"master": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "表示是否node 是 master。",
									},
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ID master 或 副本 节点。",
									},
									"zone_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ID availability 可用区 的 master 或 副本 节点。",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentRedisInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_redis_instances.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := RedisService{client: client}
	region := client.Region

	var (
		zone      string
		searchKey string
		projectId int64 = -1
		limit     int64 = -1
	)

	if temp, ok := d.GetOk("zone"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			zone = tempStr
		}
	}
	if temp, ok := d.GetOk("search_key"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			searchKey = tempStr
		}
	}

	if temp, ok := d.GetOkExists("project_id"); ok {
		tempInt := temp.(int)
		if tempInt >= 0 {
			projectId = int64(tempInt)
		}
	}

	if temp, ok := d.GetOk("limit"); ok {
		tempInt := temp.(int)
		if tempInt >= 0 {
			limit = int64(tempInt)
		}
	}

	tags := helper.GetTags(d, "tags")

	instances, err := service.DescribeInstances(ctx, zone, searchKey, projectId, limit)
	if err != nil {
		return err
	}

	var instanceList = make([]interface{}, 0, len(instances))

instanceLoop:
	for _, instance := range instances {
		if len(tags) > 0 {
			// filter by tags, must match all tags
			for k, v := range tags {
				if instance.Tags[k] != v {
					continue instanceLoop
				}
			}
		}

		var instanceDes = make(map[string]interface{})
		instanceDes["redis_id"] = instance.RedisId
		instanceDes["name"] = instance.Name
		instanceDes["zone"] = instance.Zone
		instanceDes["project_id"] = instance.ProjectId
		instanceDes["type"] = instance.Type
		instanceDes["mem_size"] = instance.MemSize
		instanceDes["status"] = instance.Status
		instanceDes["vpc_id"] = instance.VpcId
		instanceDes["subnet_id"] = instance.SubnetId
		instanceDes["ip"] = instance.Ip
		instanceDes["port"] = instance.Port
		instanceDes["create_time"] = instance.CreateTime
		instanceDes["tags"] = instance.Tags
		instanceDes["redis_shard_num"] = instance.RedisShardNum
		instanceDes["redis_replicas_num"] = instance.RedisReplicasNum
		instanceDes["type_id"] = instance.TypeId
		instanceDes["charge_type"] = instance.BillingMode
		instanceDes["node_info"] = instance.NodeInfo
		instanceList = append(instanceList, instanceDes)
	}

	if err := d.Set("instance_list", instanceList); err != nil {
		log.Printf("[CRITAL]%s provider set redis instances fail, reason:%s\n", logId, err.Error())
	}
	d.SetId("redis_instances_list" + region)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), instanceList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
		}
	}
	return nil
}
