package vpc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpcSgSnapshotFileContent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcSgSnapshotFileContentRead,
		Schema: map[string]*schema.Schema{
			"snapshot_policy_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Snapshot 策略 IDs。",
			},

			"snapshot_file_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Snapshot 文件 ID。",
			},

			"security_group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "安全组 ID",
			},

			"instance_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "安全组 ID",
			},

			"backup_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Backup 时间。",
			},

			"operator": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "操作者",
			},

			"original_data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Original 数据。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_index": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "索引 数量 安全 组 规则，其中 dynamically changes 使用 规则. 此 参数 可以 是 获取 via `DescribeSecurityGroupPolicies` API 和 使用 使用 `版本` 字段 在 返回 值 的 API。",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "协议 有效值：TCP，UDP，ICMP，ICMPv6，ALL。",
						},
						"port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "端口 (`all`， 单个 端口，或 端口 范围).注意: 如果 `协议` 值 是 集合 到 `ALL`， `端口` 值 also needs 到 是 集合 到 `all`。",
						},
						"service_template": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "协议 端口 ID 或 协议 端口 组 ID ServiceTemplate 和 协议+端口 是 mutually exclusive。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "协议 端口 ID，such 作为 `ppm-f5n1f8da`。",
									},
									"service_group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "协议 端口 组 ID，such 作为 `ppmg-f5n1f8da`。",
									},
								},
							},
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Either `CidrBlock` 或 `Ipv6CidrBlock 可以 是 指定. 注意 该 如果 `0.0.0.0/n` 是 entered，它 是 mapped 到 0.0.0.0/0。",
						},
						"ipv6_cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CIDR block 或 IPv6 (mutually exclusive)。",
						},
						"security_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "安全 组 实例 ID，such 作为 `sg-ohuuioma`。",
						},
						"address_template": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "IP 地址 ID 或 IP 地址 组 ID",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID IP 地址，such 作为 `ipm-2uw6ujo6`。",
									},
									"address_group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID IP 地址 组，such 作为 `ipmg-2uw6ujo6`。",
									},
								},
							},
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ACCEPT 或 DROP。",
						},
						"policy_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security 组 策略 描述",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "last 修改时间 的 安全 组。",
						},
					},
				},
			},

			"backup_data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Backup 数据。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_index": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "索引 数量 安全 组 规则，其中 dynamically changes 使用 规则. 此 参数 可以 是 获取 via `DescribeSecurityGroupPolicies` API 和 使用 使用 `版本` 字段 在 返回 值 的 API。",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "协议 有效值：TCP，UDP，ICMP，ICMPv6，ALL。",
						},
						"port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "端口 (`all`， 单个 端口，或 端口 范围).注意: 如果 `协议` 值 是 集合 到 `ALL`， `端口` 值 also needs 到 是 集合 到 `all`。",
						},
						"service_template": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "协议 端口 ID 或 协议 端口 组 ID ServiceTemplate 和 协议+端口 是 mutually exclusive。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "协议 端口 ID，such 作为 `ppm-f5n1f8da`。",
									},
									"service_group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "协议 端口 组 ID，such 作为 `ppmg-f5n1f8da`。",
									},
								},
							},
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Either `CidrBlock` 或 `Ipv6CidrBlock 可以 是 指定. 注意 该 如果 `0.0.0.0/n` 是 entered，它 是 mapped 到 0.0.0.0/0。",
						},
						"ipv6_cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CIDR block 或 IPv6 (mutually exclusive)。",
						},
						"security_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "安全 组 实例 ID，such 作为 `sg-ohuuioma`。",
						},
						"address_template": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "IP 地址 ID 或 IP 地址 组 ID",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID IP 地址，such 作为 `ipm-2uw6ujo6`。",
									},
									"address_group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID IP 地址 组，such 作为 `ipmg-2uw6ujo6`。",
									},
								},
							},
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ACCEPT 或 DROP。",
						},
						"policy_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security 组 策略 描述",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "last 修改时间 的 安全 组。",
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

func dataSourceTencentCloudVpcSgSnapshotFileContentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_sg_snapshot_file_content.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("snapshot_policy_id"); ok {
		paramMap["SnapshotPolicyId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("snapshot_file_id"); ok {
		paramMap["SnapshotFileId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_id"); ok {
		paramMap["SecurityGroupId"] = helper.String(v.(string))
	}

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var content *vpc.DescribeSgSnapshotFileContentResponseParams

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcSgSnapshotFileContent(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		content = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	tmpList := make([]map[string]interface{}, 0)

	if content.InstanceId != nil {
		_ = d.Set("instance_id", content.InstanceId)
	}

	if content.BackupTime != nil {
		_ = d.Set("backup_time", content.BackupTime)
	}

	if content.Operator != nil {
		_ = d.Set("operator", content.Operator)
	}

	if content.OriginalData != nil {
		for _, securityGroupPolicy := range content.OriginalData {
			securityGroupPolicyMap := map[string]interface{}{}

			if securityGroupPolicy.PolicyIndex != nil {
				securityGroupPolicyMap["policy_index"] = securityGroupPolicy.PolicyIndex
			}

			if securityGroupPolicy.Protocol != nil {
				securityGroupPolicyMap["protocol"] = securityGroupPolicy.Protocol
			}

			if securityGroupPolicy.Port != nil {
				securityGroupPolicyMap["port"] = securityGroupPolicy.Port
			}

			if securityGroupPolicy.ServiceTemplate != nil {
				serviceTemplateMap := map[string]interface{}{}

				if securityGroupPolicy.ServiceTemplate.ServiceId != nil {
					serviceTemplateMap["service_id"] = securityGroupPolicy.ServiceTemplate.ServiceId
				}

				if securityGroupPolicy.ServiceTemplate.ServiceGroupId != nil {
					serviceTemplateMap["service_group_id"] = securityGroupPolicy.ServiceTemplate.ServiceGroupId
				}

				securityGroupPolicyMap["service_template"] = []interface{}{serviceTemplateMap}
			}

			if securityGroupPolicy.CidrBlock != nil {
				securityGroupPolicyMap["cidr_block"] = securityGroupPolicy.CidrBlock
			}

			if securityGroupPolicy.Ipv6CidrBlock != nil {
				securityGroupPolicyMap["ipv6_cidr_block"] = securityGroupPolicy.Ipv6CidrBlock
			}

			if securityGroupPolicy.SecurityGroupId != nil {
				securityGroupPolicyMap["security_group_id"] = securityGroupPolicy.SecurityGroupId
			}

			if securityGroupPolicy.AddressTemplate != nil {
				addressTemplateMap := map[string]interface{}{}

				if securityGroupPolicy.AddressTemplate.AddressId != nil {
					addressTemplateMap["address_id"] = securityGroupPolicy.AddressTemplate.AddressId
				}

				if securityGroupPolicy.AddressTemplate.AddressGroupId != nil {
					addressTemplateMap["address_group_id"] = securityGroupPolicy.AddressTemplate.AddressGroupId
				}

				securityGroupPolicyMap["address_template"] = []interface{}{addressTemplateMap}
			}

			if securityGroupPolicy.Action != nil {
				securityGroupPolicyMap["action"] = securityGroupPolicy.Action
			}

			if securityGroupPolicy.PolicyDescription != nil {
				securityGroupPolicyMap["policy_description"] = securityGroupPolicy.PolicyDescription
			}

			if securityGroupPolicy.ModifyTime != nil {
				securityGroupPolicyMap["modify_time"] = securityGroupPolicy.ModifyTime
			}

			ids = append(ids, *securityGroupPolicy.SecurityGroupId)
			tmpList = append(tmpList, securityGroupPolicyMap)
		}

		_ = d.Set("original_data", tmpList)
	}

	if content.BackupData != nil {
		for _, securityGroupPolicy := range content.BackupData {
			securityGroupPolicyMap := map[string]interface{}{}

			if securityGroupPolicy.PolicyIndex != nil {
				securityGroupPolicyMap["policy_index"] = securityGroupPolicy.PolicyIndex
			}

			if securityGroupPolicy.Protocol != nil {
				securityGroupPolicyMap["protocol"] = securityGroupPolicy.Protocol
			}

			if securityGroupPolicy.Port != nil {
				securityGroupPolicyMap["port"] = securityGroupPolicy.Port
			}

			if securityGroupPolicy.ServiceTemplate != nil {
				serviceTemplateMap := map[string]interface{}{}

				if securityGroupPolicy.ServiceTemplate.ServiceId != nil {
					serviceTemplateMap["service_id"] = securityGroupPolicy.ServiceTemplate.ServiceId
				}

				if securityGroupPolicy.ServiceTemplate.ServiceGroupId != nil {
					serviceTemplateMap["service_group_id"] = securityGroupPolicy.ServiceTemplate.ServiceGroupId
				}

				securityGroupPolicyMap["service_template"] = []interface{}{serviceTemplateMap}
			}

			if securityGroupPolicy.CidrBlock != nil {
				securityGroupPolicyMap["cidr_block"] = securityGroupPolicy.CidrBlock
			}

			if securityGroupPolicy.Ipv6CidrBlock != nil {
				securityGroupPolicyMap["ipv6_cidr_block"] = securityGroupPolicy.Ipv6CidrBlock
			}

			if securityGroupPolicy.SecurityGroupId != nil {
				securityGroupPolicyMap["security_group_id"] = securityGroupPolicy.SecurityGroupId
			}

			if securityGroupPolicy.AddressTemplate != nil {
				addressTemplateMap := map[string]interface{}{}

				if securityGroupPolicy.AddressTemplate.AddressId != nil {
					addressTemplateMap["address_id"] = securityGroupPolicy.AddressTemplate.AddressId
				}

				if securityGroupPolicy.AddressTemplate.AddressGroupId != nil {
					addressTemplateMap["address_group_id"] = securityGroupPolicy.AddressTemplate.AddressGroupId
				}

				securityGroupPolicyMap["address_template"] = []interface{}{addressTemplateMap}
			}

			if securityGroupPolicy.Action != nil {
				securityGroupPolicyMap["action"] = securityGroupPolicy.Action
			}

			if securityGroupPolicy.PolicyDescription != nil {
				securityGroupPolicyMap["policy_description"] = securityGroupPolicy.PolicyDescription
			}

			if securityGroupPolicy.ModifyTime != nil {
				securityGroupPolicyMap["modify_time"] = securityGroupPolicy.ModifyTime
			}

			ids = append(ids, *securityGroupPolicy.SecurityGroupId)
			tmpList = append(tmpList, securityGroupPolicyMap)
		}

		_ = d.Set("backup_data", tmpList)
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
