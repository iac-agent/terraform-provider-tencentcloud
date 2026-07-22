package cvm

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCvmChcHosts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCvmChcHostsRead,
		Schema: map[string]*schema.Schema{
			"chc_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "CHC 主机 ID. Up 到 100 实例 per 请求 是 allowed. ChcIds 和 Filters 不能 是 指定 在 same 时间。",
			},

			"filters": {
				Optional: true,
				Type:     schema.TypeList,
				Description: "- `zone` 按可用区域过滤，例如 ap-guangzhou-1。有效值：请参见【地域和可用区】(https://www.tencentcloud.com/document/product/213/6091?from_cn_redirect=1)。\n" +
					"- `instance-name` Filter by the instance name.\n" +
					"- `instance-state` Filter by the instance status. For status details, see [InstanceStatus](https://www.tencentcloud.com/document/api/213/15753?from_cn_redirect=1#InstanceStatus).\n" +
					"- `device-type` Filter by the device type.\n" +
					"- `vpc-id` Filter by the unique VPC ID.\n" +
					"- `subnet-id` Filter by the unique VPC subnet ID.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤名称",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "过滤器 值。",
						},
					},
				},
			},

			"chc_host_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 返回 实例。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"chc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CHC 主机 ID。",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例名称",
						},
						"serial_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Server serial 数量。",
						},
						"instance_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CHC 主机 状态&lt;br/&gt;&lt;ul&gt;&lt;li&gt;REGISTERED: CHC 主机 是 registered，但 out-的-band 网络 和 部署 网络 是 不 已配置.&lt;/li&gt;&lt;li&gt;VPC_READY: out-的-band 网络 和 部署 网络 是 已配置.&lt;/li&gt;&lt;li&gt;PREPARED: It&#39;s ready 和 可以 是 associated 使用 CVM.&lt;/li&gt;&lt;li&gt;ONLINE: It&#39;s already associated 使用 CVM.&lt;/li&gt;&lt;/ul&gt;。",
						},
						"device_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Device type注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"placement": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Availability 可用区",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID availability 可用区 其中 实例 resides. You 可以 call [DescribeZones](https://www.tencentcloud.com/document/product/213/35071) API 和 obtain ID 在 返回 可用区 字段。",
									},
									"project_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ID 项目 到 其中 实例 belongs. 此 参数 可以 是 获取 从 projectId 返回 通过 DescribeProject. 如果 此 是 left 空， 默认值 项目 是 使用。",
									},
									"host_ids": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "ID 列表 CDHs 从 其中 实例 可以 是 创建. 如果 您 have purchased CDHs 和 指定this 参数， 实例 您 purchase 将 是 randomly deployed 在 CDHs。",
									},
									"host_ips": {
										Type:        schema.TypeSet,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Computed:    true,
										Deprecated:  "It has been deprecated from version 1.81.108.",
										Description: "IPs 的 hosts 到 create CVMs。",
									},
									"host_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID CDH 到 其中 实例 belongs，仅 使用 作为 output 参数。",
									},
								},
							},
						},
						"bmc_virtual_private_cloud": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Out-的-band network注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "私有网络 ID 在 格式 的 vpc-xxx. To obtain 有效 VPC IDs，您 可以 日志 在 到 [console](https://console.tencentcloud.com/vpc/vpc?rid=1) 或 call DescribeVpcEx API 和 look 对于 unVpcId 字段 在 response. 如果 您 指定DEFAULT 对于 both VpcId 和 SubnetId 当 creating 实例， 默认值 VPC 将 是 使用。",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC 子网 ID 在 格式 子网-xxx. To obtain 有效 子网 IDs，您 可以 日志 在 到 [console](https://console.tencentcloud.com/vpc/vpc?rid=1) 或 call DescribeSubnets 和 look 对于 unSubnetId 字段 在 response. 如果 您 指定DEFAULT 对于 both SubnetId 和 VpcId 当 creating 实例， 默认值 VPC 将 是 使用。",
									},
									"as_vpc_gateway": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否use CVM 实例 作为 公有 网关. 公有 网关 是 仅 可用 当 实例 has 公有 IP 和 resides 在 VPC. 有效 值:&lt;br&gt;&lt;li&gt;TRUE: yes;&lt;br&gt;&lt;li&gt;FALSE: 无&lt;br&gt;&lt;br&gt;默认值：FALSE。",
									},
									"private_ip_addresses": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "数组 VPC 子网 IPs. You 可以 使用 此 参数 当 creating 实例 或 modifying VPC attributes 的 实例. Currently 您 可以 指定multiple IPs 在 一个 子网 仅 当 creating 多个 实例 在 same 时间。",
									},
									"ipv6_address_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 IPv6 addresses randomly generated 对于 ENI。",
									},
								},
							},
						},
						"bmc_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Out-的-band 网络 IP注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"bmc_security_group_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Out-的-band 网络 安全 组 ID注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"deploy_virtual_private_cloud": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Deployment network注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "私有网络 ID 在 格式 的 vpc-xxx. To obtain 有效 VPC IDs，您 可以 日志 在 到 [console](https://console.tencentcloud.com/vpc/vpc?rid=1) 或 call DescribeVpcEx API 和 look 对于 unVpcId 字段 在 response. 如果 您 指定DEFAULT 对于 both VpcId 和 SubnetId 当 creating 实例， 默认值 VPC 将 是 使用。",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC 子网 ID 在 格式 子网-xxx. To obtain 有效 子网 IDs，您 可以 日志 在 到 [console](https://console.tencentcloud.com/vpc/vpc?rid=1) 或 call DescribeSubnets 和 look 对于 unSubnetId 字段 在 response. 如果 您 指定DEFAULT 对于 both SubnetId 和 VpcId 当 creating 实例， 默认值 VPC 将 是 使用。",
									},
									"as_vpc_gateway": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否use CVM 实例 作为 公有 网关. 公有 网关 是 仅 可用 当 实例 has 公有 IP 和 resides 在 VPC. 有效 值:&lt;br&gt;&lt;li&gt;TRUE: yes;&lt;br&gt;&lt;li&gt;FALSE: 无&lt;br&gt;&lt;br&gt;默认值：FALSE。",
									},
									"private_ip_addresses": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "数组 VPC 子网 IPs. You 可以 使用 此 参数 当 creating 实例 或 modifying VPC attributes 的 实例. Currently 您 可以 指定multiple IPs 在 一个 子网 仅 当 creating 多个 实例 在 same 时间。",
									},
									"ipv6_address_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 IPv6 addresses randomly generated 对于 ENI。",
									},
								},
							},
						},
						"deploy_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment 网络 IP注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"deploy_security_group_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Deployment 网络 安全 组 ID注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"cvm_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID associated CVM注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Server 创建时间。",
						},
						"hardware_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 hardware 描述，包括 CPU 核数，内存 容量 和 磁盘 容量.注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CPU 核数 的 CHC host注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory 容量 的 CHC 主机 (单位: GB)注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"disk": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk 容量 的 CHC host注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"bmc_mac": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "MAC 地址 assigned under out-的-band network注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"deploy_mac": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "MAC 地址 assigned under 部署 network注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"tenant_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Management typeHOSTING: HostingTENANT: Leasing注意：此字段可能返回 null，表示无法获取有效值。",
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

func dataSourceTencentCloudCvmChcHostsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cvm_chc_hosts.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("chc_ids"); ok {
		chcIdsSet := v.(*schema.Set).List()
		paramMap["ChcIds"] = helper.InterfacesStringsPoint(chcIdsSet)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*cvm.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := cvm.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var chcHostSet []*cvm.ChcHost

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCvmChcHostsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		chcHostSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(chcHostSet))
	tmpList := make([]map[string]interface{}, 0, len(chcHostSet))

	if chcHostSet != nil {
		for _, chcHost := range chcHostSet {
			chcHostMap := map[string]interface{}{}

			if chcHost.ChcId != nil {
				chcHostMap["chc_id"] = chcHost.ChcId
			}

			if chcHost.InstanceName != nil {
				chcHostMap["instance_name"] = chcHost.InstanceName
			}

			if chcHost.SerialNumber != nil {
				chcHostMap["serial_number"] = chcHost.SerialNumber
			}

			if chcHost.InstanceState != nil {
				chcHostMap["instance_state"] = chcHost.InstanceState
			}

			if chcHost.DeviceType != nil {
				chcHostMap["device_type"] = chcHost.DeviceType
			}

			if chcHost.Placement != nil {
				placementMap := map[string]interface{}{}

				if chcHost.Placement.Zone != nil {
					placementMap["zone"] = chcHost.Placement.Zone
				}

				if chcHost.Placement.ProjectId != nil {
					placementMap["project_id"] = chcHost.Placement.ProjectId
				}

				if chcHost.Placement.HostIds != nil {
					placementMap["host_ids"] = chcHost.Placement.HostIds
				}

				// It has been deprecated from version 1.81.108
				//if chcHost.Placement.HostIps != nil {
				//	placementMap["host_ips"] = chcHost.Placement.HostIps
				//}

				if chcHost.Placement.HostId != nil {
					placementMap["host_id"] = chcHost.Placement.HostId
				}

				chcHostMap["placement"] = []interface{}{placementMap}
			}

			if chcHost.BmcVirtualPrivateCloud != nil {
				bmcVirtualPrivateCloudMap := map[string]interface{}{}

				if chcHost.BmcVirtualPrivateCloud.VpcId != nil {
					bmcVirtualPrivateCloudMap["vpc_id"] = chcHost.BmcVirtualPrivateCloud.VpcId
				}

				if chcHost.BmcVirtualPrivateCloud.SubnetId != nil {
					bmcVirtualPrivateCloudMap["subnet_id"] = chcHost.BmcVirtualPrivateCloud.SubnetId
				}

				if chcHost.BmcVirtualPrivateCloud.AsVpcGateway != nil {
					bmcVirtualPrivateCloudMap["as_vpc_gateway"] = chcHost.BmcVirtualPrivateCloud.AsVpcGateway
				}

				if chcHost.BmcVirtualPrivateCloud.PrivateIpAddresses != nil {
					bmcVirtualPrivateCloudMap["private_ip_addresses"] = chcHost.BmcVirtualPrivateCloud.PrivateIpAddresses
				}

				if chcHost.BmcVirtualPrivateCloud.Ipv6AddressCount != nil {
					bmcVirtualPrivateCloudMap["ipv6_address_count"] = chcHost.BmcVirtualPrivateCloud.Ipv6AddressCount
				}

				chcHostMap["bmc_virtual_private_cloud"] = []interface{}{bmcVirtualPrivateCloudMap}
			}

			if chcHost.BmcIp != nil {
				chcHostMap["bmc_ip"] = chcHost.BmcIp
			}

			if chcHost.BmcSecurityGroupIds != nil {
				chcHostMap["bmc_security_group_ids"] = chcHost.BmcSecurityGroupIds
			}

			if chcHost.DeployVirtualPrivateCloud != nil {
				deployVirtualPrivateCloudMap := map[string]interface{}{}

				if chcHost.DeployVirtualPrivateCloud.VpcId != nil {
					deployVirtualPrivateCloudMap["vpc_id"] = chcHost.DeployVirtualPrivateCloud.VpcId
				}

				if chcHost.DeployVirtualPrivateCloud.SubnetId != nil {
					deployVirtualPrivateCloudMap["subnet_id"] = chcHost.DeployVirtualPrivateCloud.SubnetId
				}

				if chcHost.DeployVirtualPrivateCloud.AsVpcGateway != nil {
					deployVirtualPrivateCloudMap["as_vpc_gateway"] = chcHost.DeployVirtualPrivateCloud.AsVpcGateway
				}

				if chcHost.DeployVirtualPrivateCloud.PrivateIpAddresses != nil {
					deployVirtualPrivateCloudMap["private_ip_addresses"] = chcHost.DeployVirtualPrivateCloud.PrivateIpAddresses
				}

				if chcHost.DeployVirtualPrivateCloud.Ipv6AddressCount != nil {
					deployVirtualPrivateCloudMap["ipv6_address_count"] = chcHost.DeployVirtualPrivateCloud.Ipv6AddressCount
				}

				chcHostMap["deploy_virtual_private_cloud"] = []interface{}{deployVirtualPrivateCloudMap}
			}

			if chcHost.DeployIp != nil {
				chcHostMap["deploy_ip"] = chcHost.DeployIp
			}

			if chcHost.DeploySecurityGroupIds != nil {
				chcHostMap["deploy_security_group_ids"] = chcHost.DeploySecurityGroupIds
			}

			if chcHost.CvmInstanceId != nil {
				chcHostMap["cvm_instance_id"] = chcHost.CvmInstanceId
			}

			if chcHost.CreatedTime != nil {
				chcHostMap["created_time"] = chcHost.CreatedTime
			}

			if chcHost.HardwareDescription != nil {
				chcHostMap["hardware_description"] = chcHost.HardwareDescription
			}

			if chcHost.CPU != nil {
				chcHostMap["cpu"] = chcHost.CPU
			}

			if chcHost.Memory != nil {
				chcHostMap["memory"] = chcHost.Memory
			}

			if chcHost.Disk != nil {
				chcHostMap["disk"] = chcHost.Disk
			}

			if chcHost.BmcMAC != nil {
				chcHostMap["bmc_mac"] = chcHost.BmcMAC
			}

			if chcHost.DeployMAC != nil {
				chcHostMap["deploy_mac"] = chcHost.DeployMAC
			}

			if chcHost.TenantType != nil {
				chcHostMap["tenant_type"] = chcHost.TenantType
			}

			ids = append(ids, *chcHost.ChcId)
			tmpList = append(tmpList, chcHostMap)
		}

		_ = d.Set("chc_host_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
