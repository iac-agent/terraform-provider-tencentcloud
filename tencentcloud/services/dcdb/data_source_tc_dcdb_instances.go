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

func DataSourceTencentCloudDcdbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbInstancesRead,
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "实例 ids。",
			},

			"search_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "search 名称，support instancename，VIP，all。",
			},

			"search_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "search 键，support fuzzy 查询。",
			},

			"project_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional:    true,
				Description: "项目 ids。",
			},

			"excluster_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "集群 excluster 类型",
			},

			"is_filter_excluster": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "search according 到 集群 excluter 类型",
			},

			"is_filter_vpc": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "search according 到 vpc。",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "私有网络 ID，有效 当 IsFilterVpc 是 true。",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "子网 ID，有效 当 IsFilterVpc 是 true。",
			},

			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "实例 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例名称",
						},
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "app ID。",
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
						"vpc_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "私有网络 ID",
						},
						"subnet_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "子网 ID",
						},
						"status_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态 描述",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "状态",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VIP",
						},
						"vport": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "vport。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "auto 续费标识",
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
						"shard_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "分片数",
						},
						"period_end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "过期时间。",
						},
						"isolated_timestamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "isolated 时间。",
						},
						"uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "账号 uin。",
						},
						"shard_detail": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "分片 detail。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"shard_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "分片 实例 ID",
									},
									"shard_serial_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "分片 serial ID。",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "分片 状态",
									},
									"createtime": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "分片 创建时间。",
									},
									"memory": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "内存。",
									},
									"storage": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "存储。",
									},
									"shard_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "分片 ID。",
									},
									"node_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "节点 count。",
									},
									"cpu": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "CPU 核数",
									},
								},
							},
						},
						"node_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "节点 count。",
						},
						"is_tmp": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "tmp 实例 mark。",
						},
						"wan_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "wan 域名",
						},
						"wan_vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "wan VIP",
						},
						"wan_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "wan 端口",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "更新时间。",
						},
						"db_engine": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "db 引擎。",
						},
						"db_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "db 引擎 版本",
						},
						"paymode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "付费模式",
						},
						"wan_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "wan 状态，0:nonactivated，1:activated，2:closed，3:activating。",
						},
						"is_audit_supported": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "aduit support，0:support，1:unsupport。",
						},
						"instance_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例类型",
						},
						"resource_tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "资源 标签",
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

func dataSourceTencentCloudDcdbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dcdb_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_ids"); ok {
		instance_idsSet := v.(*schema.Set).List()
		ids := make([]*string, 0, len(instance_idsSet))
		for _, vv := range instance_idsSet {
			ids = append(ids, helper.String(vv.(string)))
		}
		paramMap["instance_ids"] = ids
	}

	if v, ok := d.GetOk("search_name"); ok {
		paramMap["search_name"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_key"); ok {
		paramMap["search_key"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_ids"); ok {
		project_idsSet := v.(*schema.Set).List()
		ids := make([]*int64, 0, len(project_idsSet))
		for _, vv := range project_idsSet {
			ids = append(ids, helper.IntInt64(vv.(int)))
		}
		paramMap["project_ids"] = ids
	}

	if v, ok := d.GetOk("excluster_type"); ok {
		paramMap["excluster_type"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("is_filter_excluster"); v != nil {
		paramMap["is_filter_excluster"] = helper.Bool(v.(bool))
	}

	if v, _ := d.GetOk("is_filter_vpc"); v != nil {
		paramMap["is_filter_vpc"] = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		paramMap["vpc_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		paramMap["subnet_id"] = helper.String(v.(string))
	}

	dcdbService := DcdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var instances []*dcdb.DCDBInstanceInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := dcdbService.DescribeDcdbInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instances = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Dcdb instances failed, reason:%+v", logId, err)
		return err
	}

	ids := make([]string, 0, len(instances))
	instanceList := make([]map[string]interface{}, 0, len(instances))
	if instances != nil {
		for _, instance := range instances {
			instanceMap := map[string]interface{}{}
			if instance.InstanceId != nil {
				instanceMap["instance_id"] = instance.InstanceId
			}
			if instance.InstanceName != nil {
				instanceMap["instance_name"] = instance.InstanceName
			}
			if instance.AppId != nil {
				instanceMap["app_id"] = instance.AppId
			}
			if instance.ProjectId != nil {
				instanceMap["project_id"] = instance.ProjectId
			}
			if instance.Region != nil {
				instanceMap["region"] = instance.Region
			}
			if instance.VpcId != nil {
				instanceMap["vpc_id"] = instance.VpcId
			}
			if instance.SubnetId != nil {
				instanceMap["subnet_id"] = instance.SubnetId
			}
			if instance.StatusDesc != nil {
				instanceMap["status_desc"] = instance.StatusDesc
			}
			if instance.Status != nil {
				instanceMap["status"] = instance.Status
			}
			if instance.Vip != nil {
				instanceMap["vip"] = instance.Vip
			}
			if instance.Vport != nil {
				instanceMap["vport"] = instance.Vport
			}
			if instance.CreateTime != nil {
				instanceMap["create_time"] = instance.CreateTime
			}
			if instance.AutoRenewFlag != nil {
				instanceMap["auto_renew_flag"] = instance.AutoRenewFlag
			}
			if instance.Memory != nil {
				instanceMap["memory"] = instance.Memory
			}
			if instance.Storage != nil {
				instanceMap["storage"] = instance.Storage
			}
			if instance.ShardCount != nil {
				instanceMap["shard_count"] = instance.ShardCount
			}
			if instance.PeriodEndTime != nil {
				instanceMap["period_end_time"] = instance.PeriodEndTime
			}
			if instance.IsolatedTimestamp != nil {
				instanceMap["isolated_timestamp"] = instance.IsolatedTimestamp
			}
			if instance.Uin != nil {
				instanceMap["uin"] = instance.Uin
			}
			if instance.ShardDetail != nil {
				shardDetailList := []interface{}{}
				for _, shardDetail := range instance.ShardDetail {
					shardDetailMap := map[string]interface{}{}
					if shardDetail.ShardInstanceId != nil {
						shardDetailMap["shard_instance_id"] = shardDetail.ShardInstanceId
					}
					if shardDetail.ShardSerialId != nil {
						shardDetailMap["shard_serial_id"] = shardDetail.ShardSerialId
					}
					if shardDetail.Status != nil {
						shardDetailMap["status"] = shardDetail.Status
					}
					if shardDetail.Createtime != nil {
						shardDetailMap["createtime"] = shardDetail.Createtime
					}
					if shardDetail.Memory != nil {
						shardDetailMap["memory"] = shardDetail.Memory
					}
					if shardDetail.Storage != nil {
						shardDetailMap["storage"] = shardDetail.Storage
					}
					if shardDetail.ShardId != nil {
						shardDetailMap["shard_id"] = shardDetail.ShardId
					}
					if shardDetail.NodeCount != nil {
						shardDetailMap["node_count"] = shardDetail.NodeCount
					}
					if shardDetail.Cpu != nil {
						shardDetailMap["cpu"] = shardDetail.Cpu
					}

					shardDetailList = append(shardDetailList, shardDetailMap)
				}
				instanceMap["shard_detail"] = shardDetailList
			}
			if instance.NodeCount != nil {
				instanceMap["node_count"] = instance.NodeCount
			}
			if instance.IsTmp != nil {
				instanceMap["is_tmp"] = instance.IsTmp
			}
			if instance.WanDomain != nil {
				instanceMap["wan_domain"] = instance.WanDomain
			}
			if instance.WanVip != nil {
				instanceMap["wan_vip"] = instance.WanVip
			}
			if instance.WanPort != nil {
				instanceMap["wan_port"] = instance.WanPort
			}
			if instance.UpdateTime != nil {
				instanceMap["update_time"] = instance.UpdateTime
			}
			if instance.DbEngine != nil {
				instanceMap["db_engine"] = instance.DbEngine
			}
			if instance.DbVersion != nil {
				instanceMap["db_version"] = instance.DbVersion
			}
			if instance.Paymode != nil {
				instanceMap["paymode"] = instance.Paymode
			}
			if instance.WanStatus != nil {
				instanceMap["wan_status"] = instance.WanStatus
			}
			if instance.IsAuditSupported != nil {
				instanceMap["is_audit_supported"] = instance.IsAuditSupported
			}
			if instance.InstanceType != nil {
				instanceMap["instance_type"] = instance.InstanceType
			}
			if instance.ResourceTags != nil {
				resourceTagsList := []interface{}{}
				for _, resourceTags := range instance.ResourceTags {
					resourceTagsMap := map[string]interface{}{}
					if resourceTags.TagKey != nil {
						resourceTagsMap["tag_key"] = resourceTags.TagKey
					}
					if resourceTags.TagValue != nil {
						resourceTagsMap["tag_value"] = resourceTags.TagValue
					}

					resourceTagsList = append(resourceTagsList, resourceTagsMap)
				}
				instanceMap["resource_tags"] = resourceTagsList
			}
			ids = append(ids, *instance.InstanceId)
			instanceList = append(instanceList, instanceMap)
		}
		d.SetId(helper.DataResourceIdsHash(ids))
		_ = d.Set("list", instanceList)
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), instanceList); e != nil {
			return e
		}
	}

	return nil
}
