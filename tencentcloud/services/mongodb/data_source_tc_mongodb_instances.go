package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMongodbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMongodbInstancesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID Mongodb 实例 到 是 queried。",
			},
			"instance_name_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 prefix 的 Mongodb 实例。",
			},
			"cluster_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(MONGODB_CLUSTER_TYPE),
				Description:  "类型 Mongodb 集群，和 可用 值 include 副本 集合 集群(expressed 使用 `REPLSET`)，sharding 集群(expressed 使用 `SHARD`)。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 的 Mongodb 实例 到 是 queried。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于存储结果。",
			},

			// computed
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 实例. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID Mongodb 实例。",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 Mongodb 实例。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID 项目 其中 实例 belongs。",
						},
						"cluster_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 Mongodb 集群。",
						},
						"available_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用 可用区 的 Mongodb。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID VPC。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 子网。",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "状态 Mongodb，和 可用 值 include pending initialization(expressed 使用 0)， processing(expressed 使用 1)，running(expressed 使用 2) 和 expired(expressed 使用 -2)。",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP 的 Mongodb 实例。",
						},
						"vport": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "IP 端口 的 Mongodb 实例。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 Mongodb 实例。",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "版本 的 Mongodb 引擎。",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 cpu's core。",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory 大小。",
						},
						"volume": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk 大小。",
						},
						"machine_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 Mongodb 实例。",
						},
						"shard_quantity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 sharding。",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "标签 的 Mongodb 实例。",
						},
						// payment
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "charge 类型 实例。",
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Auto 续费标识",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudMongodbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mongodb_instances.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := ""
	clusterType := -1
	name := ""
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}
	if v, ok := d.GetOk("cluster_type"); ok {
		vv := v.(string)
		if vv == MONGODB_CLUSTER_TYPE_REPLSET {
			clusterType = 0
		} else if vv == MONGODB_CLUSTER_TYPE_SHARD {
			clusterType = 1
		}
	}
	if v, ok := d.GetOk("instance_name_prefix"); ok {
		name = v.(string)
	}

	tags := helper.GetTags(d, "tags")

	mongodbService := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	mongodbs, err := mongodbService.DescribeInstancesByFilter(ctx, instanceId, clusterType)
	if err != nil {
		return err
	}

	instanceList := make([]map[string]interface{}, 0, len(mongodbs))
	ids := make([]string, 0, len(mongodbs))

instancesLoop:
	for _, mongo := range mongodbs {
		if nilFields := tccommon.CheckNil(mongo, map[string]string{
			"InstanceId":        "instance id",
			"InstanceName":      "instance name",
			"ProjectId":         "project id",
			"ClusterType":       "cluster type",
			"Zone":              "available zone",
			"VpcId":             "vpc id",
			"SubnetId":          "subnet id",
			"Status":            "status",
			"Vip":               "vip",
			"Vport":             "vport",
			"CreateTime":        "create time",
			"MongoVersion":      "engine version",
			"CpuNum":            "cpu number",
			"Memory":            "memory",
			"Volume":            "volume",
			"MachineType":       "machine type",
			"ReplicationSetNum": "shard quantity",
		}); len(nilFields) > 0 {
			return fmt.Errorf("mongodb %v are nil", nilFields)
		}

		if !strings.Contains(*mongo.InstanceName, name) {
			continue
		}

		respTags := make(map[string]string, len(mongo.Tags))
		for _, t := range mongo.Tags {
			if t.TagKey == nil {
				return errors.New("mongodb tag key is nil")
			}
			if t.TagValue == nil {
				return errors.New("mongodb tag value is nil")
			}

			respTags[*t.TagKey] = *t.TagValue
		}

		for k, v := range tags {
			if value, ok := respTags[k]; !ok || v != value {
				continue instancesLoop
			}
		}

		switch *mongo.MachineType {
		case MONGODB_MACHINE_TYPE_TGIO:
			*mongo.MachineType = MONGODB_MACHINE_TYPE_HIO10G

		case MONGODB_MACHINE_TYPE_GIO:
			*mongo.MachineType = MONGODB_MACHINE_TYPE_HIO
		}

		clusterType := MONGODB_CLUSTER_TYPE_REPLSET
		if *mongo.ClusterType == 1 {
			clusterType = MONGODB_CLUSTER_TYPE_SHARD
		}

		instance := map[string]interface{}{
			"instance_id":    mongo.InstanceId,
			"instance_name":  mongo.InstanceName,
			"project_id":     mongo.ProjectId,
			"cluster_type":   clusterType,
			"available_zone": mongo.Zone,
			"vpc_id":         mongo.VpcId,
			"subnet_id":      mongo.SubnetId,
			"status":         mongo.Status,
			"vip":            mongo.Vip,
			"vport":          mongo.Vport,
			"create_time":    mongo.CreateTime,
			"engine_version": mongo.MongoVersion,
			"cpu":            mongo.CpuNum,
			"memory":         *mongo.Memory / 1024,
			"volume":         *mongo.Volume / 1024,
			"machine_type":   mongo.MachineType,
			"shard_quantity": mongo.ReplicationSetNum,
			"tags":           respTags,
			"charge_type":    MONGODB_CHARGE_TYPE[*mongo.PayMode],
		}
		if MONGODB_CHARGE_TYPE[*mongo.PayMode] == MONGODB_CHARGE_TYPE_PREPAID {
			instance["auto_renew_flag"] = *mongo.AutoRenewFlag
		}
		instanceList = append(instanceList, instance)
		ids = append(ids, *mongo.InstanceId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if err = d.Set("instance_list", instanceList); err != nil {
		log.Printf("[CRITAL]%s provider set mongodb instance list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), instanceList); err != nil {
			return err
		}
	}

	return nil
}
