package emr

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudEmrNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEmrNodesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster 实例 ID, 实例 ID 是 作为 follows: emr-xxxxxxxx.",
			},
			"node_flag": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Node ID, the value is:
				- all: Means to get all type nodes, except cdb information.
				- master: Indicates that the master node information is obtained.
				- core: Indicates that the core node information is obtained.
				- task: indicates obtaining task node information.
				- common: means to get common node information.
				- router: Indicates obtaining router node information.
				- db: Indicates that the cdb information for the normal state is obtained.
				- recyle: Indicates that the node information in the Recycle Bin isolation, including the cdb information, is obtained.
				- renew: Indicates that all node information to be renewed, including cddb information, is obtained, and the auto-renewal node will not be returned.
				
				Note: Only the above values are now supported, entering other values will cause an error.`,
			},
			"hardware_resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "all",
				Description: "Resource 类型: Support all/主机/pod, 默认值 是 all.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Page 数量, 使用 默认值 值 的 0, 表示 first 页面.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "数量 返回 per 页面, 默认值 值 是 100, 和 最大 值 是 100.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used 到 save results.",
			},
			"nodes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List 的 节点 details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "User APPID.",
						},
						"serial_no": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Serial 数量.",
						},
						"order_no": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Machine 实例 ID.",
						},
						"wan_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "master 节点 是 bound 到 Internet IP 地址.",
						},
						"flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node 类型. 0: common 节点; 1: master 节点; 2: core 节点; 3: 任务 节点.",
						},
						"spec": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node specifications.",
						},
						"cpu_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number 的 节点 cores.",
						},
						"mem_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node 内存.",
						},
						"mem_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node 内存 描述.",
						},
						"region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "节点 是 located 在 地域.",
						},
						"zone_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Zone 其中 节点 是 located.",
						},
						"apply_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application 时间.",
						},
						"free_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Release 时间.",
						},
						"disk_size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Hard 磁盘 大小.",
						},
						"name_tag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node 描述.",
						},
						"services": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node 部署 服务.",
						},
						"storage_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk 类型.",
						},
						"root_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "大小 的 系统 磁盘.",
						},
						"charge_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "类型 的 payment.",
						},
						"cdb_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database IP.",
						},
						"cdb_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Database 端口.",
						},
						"hw_disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Hard 磁盘 容量.",
						},
						"hw_disk_size_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Hard 磁盘 容量 描述.",
						},
						"hw_mem_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory 容量.",
						},
						"hw_mem_size_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Memory 容量 描述.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expiration 时间.",
						},
						"emr_resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node 资源 ID.",
						},
						"is_auto_renew": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Renewal logo.",
						},
						"device_class": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Device identity.",
						},
						"mutable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Supports variations.",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Intranet IP.",
						},
						"destroyable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether 此 节点 是 destroyable, 1 可以 是 destroyed, 0 是 不 destroyable.",
						},
						"auto_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether 它 是 autoscaling 节点, 0 是 normal 节点, 和 1 是 autoscaling 节点.",
						},
						"hardware_resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource 类型, 主机/pod.",
						},
						"is_dynamic_spec": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Floating specifications, 1 yes, 0 无.",
						},
						"dynamic_pod_spec": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Floating 规格 值 json 字符串.",
						},
						"support_modify_pay_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether 到 support change billing 类型 1 Yes 和 0 No.",
						},
						"cdb_node_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Database 信息.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DB 实例.",
									},
									"ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Database IP.",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Database 端口.",
									},
									"mem_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Database 内存 specifications.",
									},
									"volume": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Database 磁盘 specifications.",
									},
									"service": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "服务 identity.",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Expiration 时间.",
									},
									"apply_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Application 时间.",
									},
									"pay_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "类型 的 payment.",
									},
									"expire_flag": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Expired ID.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Database 状态.",
									},
									"is_auto_renew": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Renewal identity.",
									},
									"serial_no": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Database 字符串.",
									},
									"zone_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Zone ID.",
									},
									"region_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Region ID.",
									},
								},
							},
						},
						"mc_multi_disks": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Multi-云 磁盘.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 的 云 disks 的 此 类型.",
									},
									"type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Disk 类型.",
									},
									"volume": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "大小 的 云 磁盘.",
									},
								},
							},
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "label 的 节点 binding.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag 键.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag 值.",
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

func dataSourceTencentCloudEmrNodesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_emr_nodes.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	instanceId := d.Get("instance_id").(string)
	nodeFlag := d.Get("node_flag").(string)
	offset := d.Get("offset").(int)
	limit := d.Get("limit").(int)
	hardwareResourceType := d.Get("hardware_resource_type").(string)

	emrServer := EMRService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var nodes []*emr.NodeHardwareInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		var errRet error
		nodes, errRet = emrServer.DescribeClusterNodes(ctx, instanceId, nodeFlag, hardwareResourceType, offset, limit)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	emrNodes := make([]map[string]interface{}, 0)
	ids := make([]string, 0)
	for _, node := range nodes {
		mcMultiDisks := node.MCMultiDisk
		cdbNodeInfo := node.CdbNodeInfo
		tags := node.Tags
		mcMultiDiskList := make([]map[string]interface{}, 0)
		cdbNodeInfoMap := make(map[string]interface{})
		tagList := make([]map[string]interface{}, 0)
		for _, mcMultiDisk := range mcMultiDisks {
			tmpMCMultiDisk := make(map[string]interface{})
			tmpMCMultiDisk["count"] = mcMultiDisk.Count
			tmpMCMultiDisk["type"] = mcMultiDisk.Type
			tmpMCMultiDisk["volume"] = mcMultiDisk.Volume
			mcMultiDiskList = append(mcMultiDiskList, tmpMCMultiDisk)
		}
		for _, tag := range tags {
			tmpTag := make(map[string]interface{})
			tmpTag["tag_key"] = tag.TagKey
			tmpTag["tag_value"] = tag.TagValue
			tagList = append(tagList, tmpTag)
		}

		if cdbNodeInfo != nil {
			cdbNodeInfoMap["instance_name"] = cdbNodeInfo.InstanceName
			cdbNodeInfoMap["ip"] = cdbNodeInfo.Ip
			cdbNodeInfoMap["port"] = cdbNodeInfo.Port
			cdbNodeInfoMap["mem_size"] = cdbNodeInfo.MemSize
			cdbNodeInfoMap["volume"] = cdbNodeInfo.Volume
			cdbNodeInfoMap["service"] = cdbNodeInfo.Service
			cdbNodeInfoMap["expire_time"] = cdbNodeInfo.ExpireTime
			cdbNodeInfoMap["apply_time"] = cdbNodeInfo.ApplyTime
			cdbNodeInfoMap["pay_type"] = cdbNodeInfo.PayType
			cdbNodeInfoMap["expire_flag"] = cdbNodeInfo.ExpireFlag
			cdbNodeInfoMap["status"] = cdbNodeInfo.Status
			cdbNodeInfoMap["is_auto_renew"] = cdbNodeInfo.IsAutoRenew
			cdbNodeInfoMap["serial_no"] = cdbNodeInfo.SerialNo
			cdbNodeInfoMap["zone_id"] = cdbNodeInfo.ZoneId
			cdbNodeInfoMap["region_id"] = cdbNodeInfo.RegionId
		}

		nodeMap := map[string]interface{}{
			"app_id":                  node.AppId,
			"serial_no":               node.SerialNo,
			"order_no":                node.OrderNo,
			"wan_ip":                  node.WanIp,
			"flag":                    node.Flag,
			"spec":                    node.Spec,
			"cpu_num":                 node.CpuNum,
			"mem_size":                node.MemSize,
			"mem_desc":                node.MemDesc,
			"region_id":               node.RegionId,
			"zone_id":                 node.ZoneId,
			"apply_time":              node.ApplyTime,
			"free_time":               node.FreeTime,
			"disk_size":               node.DiskSize,
			"name_tag":                node.NameTag,
			"services":                node.Services,
			"storage_type":            node.StorageType,
			"root_size":               node.RootSize,
			"charge_type":             node.ChargeType,
			"cdb_ip":                  node.CdbIp,
			"cdb_port":                node.CdbPort,
			"hw_disk_size":            node.HwDiskSize,
			"hw_disk_size_desc":       node.HwDiskSizeDesc,
			"hw_mem_size":             node.HwMemSize,
			"hw_mem_size_desc":        node.HwMemSizeDesc,
			"expire_time":             node.ExpireTime,
			"emr_resource_id":         node.EmrResourceId,
			"is_auto_renew":           node.IsAutoRenew,
			"device_class":            node.DeviceClass,
			"mutable":                 node.Mutable,
			"ip":                      node.Ip,
			"destroyable":             node.Destroyable,
			"auto_flag":               node.AutoFlag,
			"hardware_resource_type":  node.HardwareResourceType,
			"is_dynamic_spec":         node.IsDynamicSpec,
			"dynamic_pod_spec":        node.DynamicPodSpec,
			"support_modify_pay_mode": node.SupportModifyPayMode,
			"cdb_node_info":           []map[string]interface{}{cdbNodeInfoMap},
			"mc_multi_disks":          mcMultiDiskList,
			"tags":                    tagList,
		}
		ids = append(ids, (string)(*node.AppId))
		emrNodes = append(emrNodes, nodeMap)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("nodes", emrNodes)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), emrNodes); err != nil {
			return err
		}
	}
	return nil
}
