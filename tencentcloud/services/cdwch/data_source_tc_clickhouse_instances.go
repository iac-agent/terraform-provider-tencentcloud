package cdwch

import (
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdwch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func DataSourceTencentCloudClickhouseInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClickhouseInstancesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Search 通过 实例 ID，support exact matching。",
			},
			"instance_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Search 通过 实例名称，support fuzzy matching。",
			},
			"tags": {
				Optional:    true,
				Type:        schema.TypeMap,
				Description: "标签 过滤器，多个 标签 必须 是 matched 在 same 时间。",
			},
			"vips": {
				Optional:    true,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "VIP 地址 列表 对于 filtering 实例。",
			},
			"is_simple": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "是否return simplified 信息。",
			},
			"instance_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 ClickHouse 实例。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID，such 作为 `cdwch-xxxx`。",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例名称",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例状态: Init，Serving，Deleted，Deleting，Modify。",
						},
						"status_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态 描述",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 版本",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域，such 作为 `ap-guangzhou`。",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability 可用区，such 作为 `ap-guangzhou-3`。",
						},
						"region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "地域 ID",
						},
						"region_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域 描述",
						},
						"zone_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用区 描述",
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
						"access_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access 地址，such 作为 `10.0.0.1:9000`。",
						},
						"eip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Elastic IP 地址",
						},
						"ch_proxy_vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CHProxy VIP 地址",
						},
						"pay_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Payment 模式: `hour` 或 `prepay`。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "过期时间。",
						},
						"renew_flag": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Auto-renewal flag。",
						},
						"ha": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "High availability: `true` 或 `false`。",
						},
						"ha_zk": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "ZooKeeper high availability。",
						},
						"is_elastic": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否为an elastic 实例。",
						},
						"kind": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例类型: `外部`，`本地`，或 `yunti`。",
						},
						"monitor": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Monitoring 信息。",
						},
						"has_cls_topic": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether CLS 主题 是 已启用",
						},
						"cls_topic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLS 主题 ID。",
						},
						"cls_log_set_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLS 日志 集合 ID。",
						},
						"enable_xml_config": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether XML 配置 是 支持。",
						},
						"cos_bucket_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "COS 存储桶名称",
						},
						"can_attach_cbs": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether CBS 可以 是 attached。",
						},
						"can_attach_cbs_lvm": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether CBS LVM 可以 是 attached。",
						},
						"can_attach_cos": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether COS 可以 是 attached。",
						},
						"upgrade_versions": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Upgradeable versions。",
						},
						"flow_msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Workflow 消息",
						},
						"master_summary": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Master 节点 summary 信息。",
							Elem:        nodesSummarySchema(),
						},
						"common_summary": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Common 节点 summary 信息。",
							Elem:        nodesSummarySchema(),
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "标签列表",
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
						"components": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "组件列表",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "组件名称",
									},
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "组件版本",
									},
								},
							},
						},
						"instance_state_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "实例 state details。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_state": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例 state。",
									},
									"flow_create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Workflow 创建时间。",
									},
									"flow_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Workflow 名称",
									},
									"flow_progress": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Workflow progress。",
									},
									"instance_state_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例 state 描述",
									},
									"flow_msg": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Workflow 消息",
									},
									"process_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Process 名称",
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

func nodesSummarySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"spec": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "规格名称",
			},
			"node_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "节点数量",
			},
			"core": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "CPU 核数",
			},
			"memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Memory 大小 （GB）。",
			},
			"disk": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Disk 大小 （GB）。",
			},
			"disk_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Disk 类型",
			},
			"disk_desc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Disk 描述",
			},
			"attach_cbs_spec": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Attached CBS 规格。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk 类型",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk 大小 （GB）。",
						},
						"disk_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 disks。",
						},
						"disk_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk 描述",
						},
					},
				},
			},
			"sub_product_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Sub-产品类型",
			},
			"spec_core": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Specification CPU 核数",
			},
			"spec_memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Specification 内存。",
			},
			"disk_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "数量 disks。",
			},
			"max_disk_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum 磁盘 大小。",
			},
			"encrypt": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Encryption 状态",
			},
		},
	}
}

func dataSourceTencentCloudClickhouseInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clickhouse_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	// Build request
	request := cdwch.NewDescribeInstancesNewRequest()

	if v, ok := d.GetOk("instance_id"); ok {
		request.SearchInstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.SearchInstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vips"); ok {
		vipList := v.([]interface{})
		for _, vip := range vipList {
			request.Vips = append(request.Vips, helper.String(vip.(string)))
		}
	}

	if v, ok := d.GetOkExists("is_simple"); ok {
		request.IsSimple = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsMap := v.(map[string]interface{})
		for key, value := range tagsMap {
			searchTag := cdwch.SearchTags{
				TagKey:   helper.String(key),
				TagValue: helper.String(value.(string)),
			}
			request.SearchTags = append(request.SearchTags, &searchTag)
		}
	}

	var instancesList []*cdwch.InstanceInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCdwchClient().DescribeInstancesNew(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil {
			e = fmt.Errorf("DescribeInstancesNew response is nil")
			return resource.NonRetryableError(e)
		}
		instancesList = response.Response.InstancesList
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read ClickHouse instances failed, reason:%+v", logId, err)
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(instancesList))
	if instancesList != nil {
		for _, instance := range instancesList {
			instanceMap := flattenInstanceInfo(instance)
			tmpList = append(tmpList, instanceMap)
		}
		_ = d.Set("instance_list", tmpList)
	}

	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}

func flattenInstanceInfo(instance *cdwch.InstanceInfo) map[string]interface{} {
	instanceMap := map[string]interface{}{}

	if instance.InstanceId != nil {
		instanceMap["instance_id"] = instance.InstanceId
	}
	if instance.InstanceName != nil {
		instanceMap["instance_name"] = instance.InstanceName
	}
	if instance.Status != nil {
		instanceMap["status"] = instance.Status
	}
	if instance.StatusDesc != nil {
		instanceMap["status_desc"] = instance.StatusDesc
	}
	if instance.Version != nil {
		instanceMap["version"] = instance.Version
	}
	if instance.Region != nil {
		instanceMap["region"] = instance.Region
	}
	if instance.Zone != nil {
		instanceMap["zone"] = instance.Zone
	}
	if instance.RegionId != nil {
		instanceMap["region_id"] = instance.RegionId
	}
	if instance.RegionDesc != nil {
		instanceMap["region_desc"] = instance.RegionDesc
	}
	if instance.ZoneDesc != nil {
		instanceMap["zone_desc"] = instance.ZoneDesc
	}
	if instance.VpcId != nil {
		instanceMap["vpc_id"] = instance.VpcId
	}
	if instance.SubnetId != nil {
		instanceMap["subnet_id"] = instance.SubnetId
	}
	if instance.AccessInfo != nil {
		instanceMap["access_info"] = instance.AccessInfo
	}
	if instance.Eip != nil {
		instanceMap["eip"] = instance.Eip
	}
	if instance.CHProxyVip != nil {
		instanceMap["ch_proxy_vip"] = instance.CHProxyVip
	}
	if instance.PayMode != nil {
		instanceMap["pay_mode"] = instance.PayMode
	}
	if instance.CreateTime != nil {
		instanceMap["create_time"] = instance.CreateTime
	}
	if instance.ExpireTime != nil {
		instanceMap["expire_time"] = instance.ExpireTime
	}
	if instance.RenewFlag != nil {
		instanceMap["renew_flag"] = instance.RenewFlag
	}
	if instance.HA != nil {
		instanceMap["ha"] = instance.HA
	}
	if instance.HAZk != nil {
		instanceMap["ha_zk"] = instance.HAZk
	}
	if instance.IsElastic != nil {
		instanceMap["is_elastic"] = instance.IsElastic
	}
	if instance.Kind != nil {
		instanceMap["kind"] = instance.Kind
	}
	if instance.Monitor != nil {
		instanceMap["monitor"] = instance.Monitor
	}
	if instance.HasClsTopic != nil {
		instanceMap["has_cls_topic"] = instance.HasClsTopic
	}
	if instance.ClsTopicId != nil {
		instanceMap["cls_topic_id"] = instance.ClsTopicId
	}
	if instance.ClsLogSetId != nil {
		instanceMap["cls_log_set_id"] = instance.ClsLogSetId
	}
	if instance.EnableXMLConfig != nil {
		instanceMap["enable_xml_config"] = instance.EnableXMLConfig
	}
	if instance.CosBucketName != nil {
		instanceMap["cos_bucket_name"] = instance.CosBucketName
	}
	if instance.CanAttachCbs != nil {
		instanceMap["can_attach_cbs"] = instance.CanAttachCbs
	}
	if instance.CanAttachCbsLvm != nil {
		instanceMap["can_attach_cbs_lvm"] = instance.CanAttachCbsLvm
	}
	if instance.CanAttachCos != nil {
		instanceMap["can_attach_cos"] = instance.CanAttachCos
	}
	if instance.UpgradeVersions != nil {
		instanceMap["upgrade_versions"] = instance.UpgradeVersions
	}
	if instance.FlowMsg != nil {
		instanceMap["flow_msg"] = instance.FlowMsg
	}

	// Master Summary
	if instance.MasterSummary != nil {
		instanceMap["master_summary"] = []interface{}{flattenNodesSummary(instance.MasterSummary)}
	}

	// Common Summary
	if instance.CommonSummary != nil {
		instanceMap["common_summary"] = []interface{}{flattenNodesSummary(instance.CommonSummary)}
	}

	// Tags
	if instance.Tags != nil {
		tagsList := make([]interface{}, 0, len(instance.Tags))
		for _, tag := range instance.Tags {
			tagMap := map[string]interface{}{}
			if tag.TagKey != nil {
				tagMap["tag_key"] = tag.TagKey
			}
			if tag.TagValue != nil {
				tagMap["tag_value"] = tag.TagValue
			}
			tagsList = append(tagsList, tagMap)
		}
		instanceMap["tags"] = tagsList
	}

	// Components
	if instance.Components != nil {
		componentsList := make([]interface{}, 0, len(instance.Components))
		for _, component := range instance.Components {
			componentMap := map[string]interface{}{}
			if component.Name != nil {
				componentMap["name"] = component.Name
			}
			if component.Version != nil {
				componentMap["version"] = component.Version
			}
			componentsList = append(componentsList, componentMap)
		}
		instanceMap["components"] = componentsList
	}

	// Instance State Info
	if instance.InstanceStateInfo != nil {
		stateInfo := instance.InstanceStateInfo
		stateInfoMap := map[string]interface{}{}
		if stateInfo.InstanceState != nil {
			stateInfoMap["instance_state"] = stateInfo.InstanceState
		}
		if stateInfo.FlowCreateTime != nil {
			stateInfoMap["flow_create_time"] = stateInfo.FlowCreateTime
		}
		if stateInfo.FlowName != nil {
			stateInfoMap["flow_name"] = stateInfo.FlowName
		}
		if stateInfo.FlowProgress != nil {
			stateInfoMap["flow_progress"] = stateInfo.FlowProgress
		}
		if stateInfo.InstanceStateDesc != nil {
			stateInfoMap["instance_state_desc"] = stateInfo.InstanceStateDesc
		}
		if stateInfo.FlowMsg != nil {
			stateInfoMap["flow_msg"] = stateInfo.FlowMsg
		}
		if stateInfo.ProcessName != nil {
			stateInfoMap["process_name"] = stateInfo.ProcessName
		}
		instanceMap["instance_state_info"] = []interface{}{stateInfoMap}
	}

	return instanceMap
}

func flattenNodesSummary(summary *cdwch.NodesSummary) map[string]interface{} {
	summaryMap := map[string]interface{}{}

	if summary.Spec != nil {
		summaryMap["spec"] = summary.Spec
	}
	if summary.NodeSize != nil {
		summaryMap["node_size"] = summary.NodeSize
	}
	if summary.Core != nil {
		summaryMap["core"] = summary.Core
	}
	if summary.Memory != nil {
		summaryMap["memory"] = summary.Memory
	}
	if summary.Disk != nil {
		summaryMap["disk"] = summary.Disk
	}
	if summary.DiskType != nil {
		summaryMap["disk_type"] = summary.DiskType
	}
	if summary.DiskDesc != nil {
		summaryMap["disk_desc"] = summary.DiskDesc
	}
	if summary.SubProductType != nil {
		summaryMap["sub_product_type"] = summary.SubProductType
	}
	if summary.SpecCore != nil {
		summaryMap["spec_core"] = summary.SpecCore
	}
	if summary.SpecMemory != nil {
		summaryMap["spec_memory"] = summary.SpecMemory
	}
	if summary.DiskCount != nil {
		summaryMap["disk_count"] = summary.DiskCount
	}
	if summary.MaxDiskSize != nil {
		summaryMap["max_disk_size"] = summary.MaxDiskSize
	}
	if summary.Encrypt != nil {
		summaryMap["encrypt"] = summary.Encrypt
	}

	// AttachCBSSpec
	if summary.AttachCBSSpec != nil {
		cbsSpec := summary.AttachCBSSpec
		cbsSpecMap := map[string]interface{}{}
		if cbsSpec.DiskType != nil {
			cbsSpecMap["disk_type"] = cbsSpec.DiskType
		}
		if cbsSpec.DiskSize != nil {
			cbsSpecMap["disk_size"] = cbsSpec.DiskSize
		}
		if cbsSpec.DiskCount != nil {
			cbsSpecMap["disk_count"] = cbsSpec.DiskCount
		}
		if cbsSpec.DiskDesc != nil {
			cbsSpecMap["disk_desc"] = cbsSpec.DiskDesc
		}
		summaryMap["attach_cbs_spec"] = []interface{}{cbsSpecMap}
	}

	return summaryMap
}
