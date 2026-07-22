package cwp

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cwpv20180228 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cwp/v20180228"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCwpMachines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCwpMachinesRead,
		Schema: map[string]*schema.Schema{
			"machine_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "类型 machine's 可用区\nCVM: Cloud Virtual Machine\nBM: BMECM: Edge Computing Machine\nLH: Lighthouse\nOther: Hybrid Cloud 可用区",
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "过滤器 criteria\n<li>Ips - String - 必填: 无 - 查询 通过 IP</li>\n<li>Names - String - 必填: 无 - 查询 通过 实例名称</li>\n<li>InstanceIds - String - 必填: 无 - 实例 ID 对于 查询 </li>\n<li>状态 - String - 必填: 无 - 客户端 online 状态 (OFFLINE: offline/shut down | ONLINE: online | UNINSTALLED: 不 installed | AGENT_OFFLINE: agent offline | AGENT_SHUTDOWN: agent shut down)</li>\n<li>版本 - String 必填: 无 - 当前 edition ( PRO_VERSION: Pro Edition | BASIC_VERSION: Basic Edition | Flagship: Ultimate Edition | ProtectedMachines: Pro + Ultimate Editions)</li>\n<li>Risk - String - 必填: 无 - risky 主机 (yes)</li>\n<li>Os - String - 必填: 无 - operating 系统 (值 的 DescribeMachineOsList)</li>\nEach 过滤器 criterion 支持 仅 一个 值\n<li>Quuid - String - 必填: 无 - CVM 实例 UUID. Maximum 值: 100.</li>\n<li>AddedOnTheFifteen - String 必填: 无 - 是否query 仅 hosts added within last 15 days (1: yes) </li>\n<li> TagId - String 必填: 无 - 查询 列表 hosts associated 使用 指定 标签 </li>。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "名称 过滤器 键",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "一个或多个过滤值",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"exact_match": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Fuzzy search。",
						},
					},
				},
			},

			"machine_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Machine 地域 For 示例，ap-guangzhou 和 ap-shanghai。",
			},

			"project_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "ID 列表 Businesses 到 其中 machines belong。",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"machines": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "列表 hosts。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"machine_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "主机名",
						},
						"machine_os": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "主机 System。",
						},
						"machine_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "主机 状态\n<li>OFFLINE: Offline</li>\n<li>ONLINE: Online</li>\n<li>SHUTDOWN: Shut down</li>\n<li>UNINSTALLED: Unprotected</li>。",
						},
						"agent_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ONLINE: Protected; OFFLINE: Offline; UNINSTALLED: Not installed。",
						},
						"instance_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "RUNNING; STOPPED; EXPIRED (awaiting recycling)。",
						},
						"uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Yunjing 客户端 UUID. 如果 客户端 是 offline 对于 long 时间， 空 字符串 是 返回。",
						},
						"quuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CVM 或 BM Machine Unique UUID。",
						},
						"vul_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 vulnerabilities。",
						},
						"machine_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "主机 IP",
						},
						"is_pro_version": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否edition 是 Pro Edition\n<li>true: yes</li>\n<li>false: 无</li>。",
						},
						"machine_wan_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "公网 IP 地址 的 主机",
						},
						"pay_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "主机 状态\n<li>POSTPAY: postpaid，indicating pay-作为-您-go 模式 </li>\n<li>PREPAY: prepaid，indicating monthly subscription 模式</li>。",
						},
						"malware_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 Trojans。",
						},
						"tag": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "标签 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rid": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Associated 标签 ID。",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签 名称",
									},
									"tag_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "标签 ID。",
									},
								},
							},
						},
						"baseline_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 baseline risks。",
						},
						"cyber_attack_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 网络 risks。",
						},
						"security_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Risk 状态\n<li>SAFE: Safe</li>\n<li>RISK: Risk</li>\n<li>UNKNOWN: Unknown</li>。",
						},
						"invasion_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 intrusion events。",
						},
						"region_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "地域 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "地域 identifiers，such 作为 ap-guangzhou，ap-shanghai，和 ap-beijing。",
									},
									"region_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Chinese 名称 地域，such 作为 South China (Guangzhou)，East China (Shanghai Finance)，和 North China (Beijing)。",
									},
									"region_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "地域 ID",
									},
									"region_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "地域 代码，such 作为 gz，sh，和 bj。",
									},
									"region_name_en": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "English 名称 地域",
									},
								},
							},
						},
						"instance_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例状态: TERMINATED_PRO_VERSION - terminated。",
						},
						"license_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Tamper-proof; authorization 状态: 1 - authorized; 0 - unauthorized。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "项目 ID",
						},
						"has_asset_scan": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether there 是 可用 asset scanning API: 0 - 无; 1 - yes。",
						},
						"machine_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Machine 可用区 类型 CVM - Cloud Virtual Machine; BM: Bare Metal; ECM: Edge Computing Machine; LH: Lightweight Application Server; Other: Hybrid Cloud 可用区",
						},
						"kernel_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Kernel 版本",
						},
						"protect_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protection 版本: BASIC_VERSION - Basic Edition; PRO_VERSION - Professional Edition; Flagship - Ultimate Edition; GENERAL_DISCOUNT - Inclusive Edition。",
						},
						"cloud_tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Cloud 标签 Information\n注意：此字段可能返回 null，表示无法获取有效值。",
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
						"is_added_on_the_fifteen": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether 主机 added within last 15 days: 0: 无; 1: yes\n注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"ip_list": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "主机 IP List\n注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network\n注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"machine_extra_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Additional 信息\n注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"wan_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "公网 IP 地址\n注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"private_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "内网 IP 地址\n注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"network_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Network 类型 1: VPC 网络; 2: Basic Network; 3: Non-Tencent Cloud Network\n注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"network_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network 名称，返回vpc_id 在 case 的 VPC 网络\n注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例 ID\n注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"host_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "主机名\n注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "备注\n注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"agent_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "主机 安全 agent 版本",
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

func dataSourceTencentCloudCwpMachinesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cwp_machines.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(nil)
		ctx           = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service       = CwpService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		machineType   string
		machineRegion string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("machine_type"); ok {
		paramMap["MachineType"] = helper.String(v.(string))
		machineType = v.(string)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*cwpv20180228.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := cwpv20180228.Filter{}
			if v, ok := filtersMap["name"].(string); ok && v != "" {
				filter.Name = helper.String(v)
			}

			if v, ok := filtersMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				for i := range valuesSet {
					values := valuesSet[i].(string)
					filter.Values = append(filter.Values, helper.String(values))
				}
			}

			if v, ok := filtersMap["exact_match"].(bool); ok {
				filter.ExactMatch = helper.Bool(v)
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	if v, ok := d.GetOk("machine_region"); ok {
		paramMap["MachineRegion"] = helper.String(v.(string))
		machineRegion = v.(string)
	}

	if v, ok := d.GetOk("project_ids"); ok {
		projectIdsList := []*uint64{}
		projectIdsSet := v.(*schema.Set).List()
		for i := range projectIdsSet {
			projectIds := projectIdsSet[i].(int)
			projectIdsList = append(projectIdsList, helper.IntUint64(projectIds))
		}

		paramMap["ProjectIds"] = projectIdsList
	}

	var respData []*cwpv20180228.Machine
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCwpMachinesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	machinesList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, machines := range respData {
			machinesMap := map[string]interface{}{}
			if machines.MachineName != nil {
				machinesMap["machine_name"] = machines.MachineName
			}

			if machines.MachineOs != nil {
				machinesMap["machine_os"] = machines.MachineOs
			}

			if machines.MachineStatus != nil {
				machinesMap["machine_status"] = machines.MachineStatus
			}

			if machines.AgentStatus != nil {
				machinesMap["agent_status"] = machines.AgentStatus
			}

			if machines.InstanceStatus != nil {
				machinesMap["instance_status"] = machines.InstanceStatus
			}

			if machines.Uuid != nil {
				machinesMap["uuid"] = machines.Uuid
			}

			if machines.Quuid != nil {
				machinesMap["quuid"] = machines.Quuid
			}

			if machines.VulNum != nil {
				machinesMap["vul_num"] = machines.VulNum
			}

			if machines.MachineIp != nil {
				machinesMap["machine_ip"] = machines.MachineIp
			}

			if machines.IsProVersion != nil {
				machinesMap["is_pro_version"] = machines.IsProVersion
			}

			if machines.MachineWanIp != nil {
				machinesMap["machine_wan_ip"] = machines.MachineWanIp
			}

			if machines.PayMode != nil {
				machinesMap["pay_mode"] = machines.PayMode
			}

			if machines.MalwareNum != nil {
				machinesMap["malware_num"] = machines.MalwareNum
			}

			tagList := make([]map[string]interface{}, 0, len(machines.Tag))
			if machines.Tag != nil {
				for _, tag := range machines.Tag {
					tagMap := map[string]interface{}{}

					if tag.Rid != nil {
						tagMap["rid"] = tag.Rid
					}

					if tag.Name != nil {
						tagMap["name"] = tag.Name
					}

					if tag.TagId != nil {
						tagMap["tag_id"] = tag.TagId
					}

					tagList = append(tagList, tagMap)
				}

				machinesMap["tag"] = tagList
			}
			if machines.BaselineNum != nil {
				machinesMap["baseline_num"] = machines.BaselineNum
			}

			if machines.CyberAttackNum != nil {
				machinesMap["cyber_attack_num"] = machines.CyberAttackNum
			}

			if machines.SecurityStatus != nil {
				machinesMap["security_status"] = machines.SecurityStatus
			}

			if machines.InvasionNum != nil {
				machinesMap["invasion_num"] = machines.InvasionNum
			}

			regionInfoMap := map[string]interface{}{}
			if machines.RegionInfo != nil {
				if machines.RegionInfo.Region != nil {
					regionInfoMap["region"] = machines.RegionInfo.Region
				}

				if machines.RegionInfo.RegionName != nil {
					regionInfoMap["region_name"] = machines.RegionInfo.RegionName
				}

				if machines.RegionInfo.RegionId != nil {
					regionInfoMap["region_id"] = machines.RegionInfo.RegionId
				}

				if machines.RegionInfo.RegionCode != nil {
					regionInfoMap["region_code"] = machines.RegionInfo.RegionCode
				}

				if machines.RegionInfo.RegionNameEn != nil {
					regionInfoMap["region_name_en"] = machines.RegionInfo.RegionNameEn
				}

				machinesMap["region_info"] = []interface{}{regionInfoMap}
			}

			if machines.InstanceState != nil {
				machinesMap["instance_state"] = machines.InstanceState
			}

			if machines.LicenseStatus != nil {
				machinesMap["license_status"] = machines.LicenseStatus
			}

			if machines.ProjectId != nil {
				machinesMap["project_id"] = machines.ProjectId
			}

			if machines.HasAssetScan != nil {
				machinesMap["has_asset_scan"] = machines.HasAssetScan
			}

			if machines.MachineType != nil {
				machinesMap["machine_type"] = machines.MachineType
			}

			if machines.KernelVersion != nil {
				machinesMap["kernel_version"] = machines.KernelVersion
			}

			if machines.ProtectType != nil {
				machinesMap["protect_type"] = machines.ProtectType
			}

			cloudTagsList := make([]map[string]interface{}, 0, len(machines.CloudTags))
			if machines.CloudTags != nil {
				for _, cloudTags := range machines.CloudTags {
					cloudTagsMap := map[string]interface{}{}
					if cloudTags.TagKey != nil {
						cloudTagsMap["tag_key"] = cloudTags.TagKey
					}

					if cloudTags.TagValue != nil {
						cloudTagsMap["tag_value"] = cloudTags.TagValue
					}

					cloudTagsList = append(cloudTagsList, cloudTagsMap)
				}

				machinesMap["cloud_tags"] = cloudTagsList
			}

			if machines.IsAddedOnTheFifteen != nil {
				machinesMap["is_added_on_the_fifteen"] = machines.IsAddedOnTheFifteen
			}

			if machines.IpList != nil {
				machinesMap["ip_list"] = machines.IpList
			}

			if machines.VpcId != nil {
				machinesMap["vpc_id"] = machines.VpcId
			}

			machineExtraInfoMap := map[string]interface{}{}
			if machines.MachineExtraInfo != nil {
				if machines.MachineExtraInfo.WanIP != nil {
					machineExtraInfoMap["wan_ip"] = machines.MachineExtraInfo.WanIP
				}

				if machines.MachineExtraInfo.PrivateIP != nil {
					machineExtraInfoMap["private_ip"] = machines.MachineExtraInfo.PrivateIP
				}

				if machines.MachineExtraInfo.NetworkType != nil {
					machineExtraInfoMap["network_type"] = machines.MachineExtraInfo.NetworkType
				}

				if machines.MachineExtraInfo.NetworkName != nil {
					machineExtraInfoMap["network_name"] = machines.MachineExtraInfo.NetworkName
				}

				if machines.MachineExtraInfo.InstanceID != nil {
					machineExtraInfoMap["instance_id"] = machines.MachineExtraInfo.InstanceID
				}

				if machines.MachineExtraInfo.HostName != nil {
					machineExtraInfoMap["host_name"] = machines.MachineExtraInfo.HostName
				}

				machinesMap["machine_extra_info"] = []interface{}{machineExtraInfoMap}
			}

			if machines.InstanceId != nil {
				machinesMap["instance_id"] = machines.InstanceId
			}

			if machines.Remark != nil {
				machinesMap["remark"] = machines.Remark
			}

			if machines.AgentVersion != nil {
				machinesMap["agent_version"] = machines.AgentVersion
			}

			machinesList = append(machinesList, machinesMap)
		}

		_ = d.Set("machines", machinesList)
	}

	d.SetId(strings.Join([]string{machineType, machineRegion}, tccommon.FILED_SP))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), machinesList); e != nil {
			return e
		}
	}

	return nil
}
