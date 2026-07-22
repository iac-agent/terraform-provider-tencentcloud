package ckafka

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCkafkaInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaInstancesRead,

		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter by instance ID。",
			},
			"search_word": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by 实例名称，support fuzzy query。",
			},
			"tag_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Matches the 标签键 值",
			},
			"status": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "(Filter Criteria) The 状态 instance. 0: Create，1: Run，2: Delete，do not fill the default return all。",
			},
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The field that needs to be filtered。",
						},
						"values": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "The filtered 值 of the field。",
						},
					},
				},
				Description: "Filter. filter.名称 supports ('Ip'，'VpcId'，'SubNetId'，'InstanceType','实例 ID')，filter.values can pass up to 10 values。",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Deprecated:  "This parameter is deprecated and will be removed in a future version. The data source now automatically retrieves all instances.",
				Description: "The page start 偏移量，默认为 `0`。",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Deprecated:  "This parameter is deprecated and will be removed in a future version. The data source now automatically retrieves all instances.",
				Description: "The 数量 pages，默认为 `10`。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 ckafka users. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance ID。",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The 实例名称",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Virtual IP。",
						},
						"vport": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Virtual PORT。",
						},
						"vip_list": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Virtual IP。",
									},
									"vport": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Virtual PORT。",
									},
								},
							},
							Description: "Virtual IP entities。",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The 状态 instance. 0: Created，1: Running，2: Delete: 5 Quarantined，-1 Creation failed。",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance bandwidth，in Mbps。",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The storage size of the instance，（GB）。",
						},
						"zone_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Availability 可用区 ID",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VpcId，如果为空，表示that it is the underlying network。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "子网 ID",
						},
						"renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否instance is renewed，the int enumeration 值: 1 表示auto-renewal，and 2 表示that it is not automatically renewed。",
						},
						"healthy": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例状态 int: 1 表示health，2 表示alarm，and 3 表示abnormal 实例状态",
						},
						"healthy_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例状态 information。",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The time when the instance was created。",
						},
						"expire_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The instance 过期时间。",
						},
						"is_internal": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否为an internal customer. A 值 of 1 表示an internal customer。",
						},
						"topic_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The 数量 topics。",
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签键",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签值",
									},
								},
							},
							Description: "标签 information。",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Kafka 版本 information. Note: This field may return null，indicating that a valid 值 could not be retrieved。",
						},
						"zone_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Description: "Across Availability Zones. Note: This field may return null，indicating that a valid 值 could not be retrieved。",
						},
						"cvm": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ckafka sale 类型 Note: This field may return null，indicating that a valid 值 could not be retrieved。",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ckafka 实例类型 Note: This field may return null，indicating that a valid 值 could not be retrieved。",
						},
						"disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk 类型 Note: This field may return null，indicating that a valid 值 could not be retrieved。",
						},
						"max_topic_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The 最大topics in the current specifications. Note: This field may return null，indicating that a valid 值 could not be retrieved.。",
						},
						"max_partition_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The 最大Partitions for the current specifications. Note: This field may return null，indicating that a valid 值 could not be retrieved。",
						},
						"rebalance_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Schedule the upgrade configuration time. Note: This field may return null，indicating that a valid 值 could not be retrieved.。",
						},
						"partition_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The current 数量 instances. Note: This field may return null，indicating that a valid 值 could not be retrieved.。",
						},
						"public_network_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 Internet bandwidth. Note: This field may return null，indicating that a valid 值 could not be retrieved.。",
						},
						"public_network": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Internet bandwidth 值 Note: This field may return null，indicating that a valid 值 could not be retrieved.。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCkafkaInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ckafka_instances.read")()

	ctx := context.Background()
	ckafkaService := CkafkaService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	// Build param map for service function
	param := make(map[string]interface{})
	if v, ok := d.GetOk("instance_ids"); ok {
		param["instance_ids"] = helper.InterfacesStringsPoint(v.([]interface{}))
	}
	if v, ok := d.GetOk("search_word"); ok {
		param["search_word"] = helper.String(v.(string))
	}
	if v, ok := d.GetOk("tag_key"); ok {
		param["tag_key"] = helper.String(v.(string))
	}
	if v, ok := d.GetOk("status"); ok {
		param["status"] = helper.InterfacesIntInt64Point(v.([]interface{}))
	}
	if v, ok := d.GetOk("filters"); ok {
		filterParams := v.([]interface{})
		filters := make([]*ckafka.Filter, 0)
		for _, filterParam := range filterParams {
			filterParamMap := filterParam.(map[string]interface{})
			filters = append(filters, &ckafka.Filter{
				Name:   helper.String(filterParamMap["name"].(string)),
				Values: helper.InterfacesStringsPoint(filterParamMap["values"].([]interface{})),
			})
		}
		param["filters"] = filters
	}

	// Note: offset and limit are deprecated and ignored - function always retrieves all results

	// Call service function with automatic pagination and retry
	kafkaInstanceDetails, err := ckafkaService.DescribeInstancesByFilter(ctx, param)
	if err != nil {
		return err
	}

	result := make([]map[string]interface{}, 0)
	ids := make([]string, 0)
	for _, kafkaInstanceDetail := range kafkaInstanceDetails {
		kafkaInstanceDetailMap := make(map[string]interface{})
		ids = append(ids, *kafkaInstanceDetail.InstanceId)
		kafkaInstanceDetailMap["instance_id"] = kafkaInstanceDetail.InstanceId
		kafkaInstanceDetailMap["instance_name"] = kafkaInstanceDetail.InstanceName
		kafkaInstanceDetailMap["vip"] = kafkaInstanceDetail.Vip
		kafkaInstanceDetailMap["vport"] = kafkaInstanceDetail.Vport
		kafkaInstanceDetailMap["status"] = kafkaInstanceDetail.Status
		kafkaInstanceDetailMap["bandwidth"] = kafkaInstanceDetail.Bandwidth
		kafkaInstanceDetailMap["disk_size"] = kafkaInstanceDetail.DiskSize
		kafkaInstanceDetailMap["zone_id"] = kafkaInstanceDetail.ZoneId
		kafkaInstanceDetailMap["vpc_id"] = kafkaInstanceDetail.VpcId
		kafkaInstanceDetailMap["subnet_id"] = kafkaInstanceDetail.SubnetId
		kafkaInstanceDetailMap["renew_flag"] = kafkaInstanceDetail.RenewFlag
		kafkaInstanceDetailMap["healthy"] = kafkaInstanceDetail.Healthy
		kafkaInstanceDetailMap["healthy_message"] = kafkaInstanceDetail.HealthyMessage
		kafkaInstanceDetailMap["create_time"] = kafkaInstanceDetail.CreateTime
		kafkaInstanceDetailMap["expire_time"] = kafkaInstanceDetail.ExpireTime
		kafkaInstanceDetailMap["is_internal"] = kafkaInstanceDetail.IsInternal
		kafkaInstanceDetailMap["topic_num"] = kafkaInstanceDetail.TopicNum
		kafkaInstanceDetailMap["version"] = kafkaInstanceDetail.Version
		kafkaInstanceDetailMap["cvm"] = kafkaInstanceDetail.Cvm
		kafkaInstanceDetailMap["instance_type"] = kafkaInstanceDetail.InstanceType
		kafkaInstanceDetailMap["max_topic_number"] = kafkaInstanceDetail.MaxTopicNumber
		kafkaInstanceDetailMap["max_partition_number"] = kafkaInstanceDetail.MaxPartitionNumber
		kafkaInstanceDetailMap["rebalance_time"] = kafkaInstanceDetail.RebalanceTime
		kafkaInstanceDetailMap["partition_number"] = kafkaInstanceDetail.PartitionNumber
		kafkaInstanceDetailMap["public_network_charge_type"] = kafkaInstanceDetail.PublicNetworkChargeType
		kafkaInstanceDetailMap["public_network"] = kafkaInstanceDetail.PublicNetwork

		vipList := make([]map[string]string, 0)
		for _, vip := range kafkaInstanceDetail.VipList {
			vipList = append(vipList, map[string]string{
				"vip":   *vip.Vip,
				"vport": *vip.Vport,
			})
		}
		kafkaInstanceDetailMap["vip_list"] = vipList

		tags := make([]map[string]string, 0)
		for _, tag := range kafkaInstanceDetail.Tags {
			tags = append(tags, map[string]string{
				"tag_key":   *tag.TagKey,
				"tag_value": *tag.TagValue,
			})
		}
		kafkaInstanceDetailMap["tags"] = tags

		zoneIds := make([]int64, 0)
		for _, zoneId := range kafkaInstanceDetail.ZoneIds {
			zoneIds = append(zoneIds, *zoneId)
		}
		kafkaInstanceDetailMap["zone_ids"] = zoneIds

		result = append(result, kafkaInstanceDetailMap)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("instance_list", result)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
