package vpc

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudEnis() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEnisRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
				ConflictsWith: []string{"vpc_id", "subnet_id", "instance_id", "security_group", "name", "description", "ipv4", "tags"},
				Description:   "ID ENIs 到 是 queried. Conflict 使用 `vpc_id`,`subnet_id`,`instance_id`,`security_group`,`名称`,`ipv4` 和 `标签`。",
			},
			"vpc_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "ID vpc 到 是 queried. Conflict 使用 `ids`。",
			},
			"subnet_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "ID 子网 within 此 vpc 到 是 queried. Conflict 使用 `ids`。",
			},
			"instance_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "ID 实例 其中 bind ENI. Conflict 使用 `ids`。",
			},
			"security_group": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "A 集合 的 安全 组 IDs 其中 bind ENI. Conflict 使用 `ids`。",
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "名称 ENI 到 是 queried. Conflict 使用 `ids`。",
			},
			"description": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "描述 ENI. Conflict 使用 `ids`。",
			},
			"ipv4": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "Intranet IP 的 ENI. Conflict 使用 `ids`。",
			},
			"tags": {
				Type:          schema.TypeMap,
				Optional:      true,
				ConflictsWith: []string{"ids"},
				Description:   "标签 的 ENI. Conflict 使用 `ids`。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// computed
			"enis": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An 信息 列表 ENIs. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID ENI。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 ENI。",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "描述 ENI。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID vpc。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 子网 within 此 vpc。",
						},
						"security_groups": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "A 集合 的 安全 组 IDs 其中 bind ENI。",
						},
						"primary": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "表示是否IP 是 primary。",
						},
						"mac": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "MAC 地址",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "States 的 ENI。",
						},
						"ipv4s": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A 集合 的 intranet IPv4s。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Intranet IP。",
									},
									"primary": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "表示是否IP 是 primary。",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "描述 IP。",
									},
								},
							},
						},
						"ipv6s": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A 集合 的 intranet IPv6s。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "`IPv6` 地址，such 作为 `3402:4e00:20:100:0:8cd9:2a67:71f3`。",
									},
									"primary": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否为a primary `IP`。",
									},
									"address_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "`ID` 的 `EIP` 实例，such 作为 `eip-hxlqja90`。",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "消息 描述",
									},
									"is_wan_ip_blocked": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否public IP 是 blocked。",
									},
								},
							},
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 实例 其中 bind ENI。",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "标签 的 ENI。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 ENI。",
						},
						"cdc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CDC 实例 ID",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudEnisRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_enis.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		ids      []string
		vpcId    *string
		subnetId *string
		cvmId    *string
		sgId     *string
		name     *string
		desc     *string
		ipv4     *string
	)

	if raw, ok := d.GetOk("ids"); ok {
		ids = helper.InterfacesStrings(raw.(*schema.Set).List())
	}

	if raw, ok := d.GetOk("vpc_id"); ok {
		vpcId = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("subnet_id"); ok {
		subnetId = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("instance_id"); ok {
		cvmId = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("security_group"); ok {
		sgId = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("name"); ok {
		name = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("description"); ok {
		desc = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("ipv4"); ok {
		ipv4 = helper.String(raw.(string))
	}
	tags := helper.GetTags(d, "tags")

	var (
		respEnis []*vpc.NetworkInterface
		err      error
	)

	if len(ids) > 0 {
		respEnis, err = service.DescribeEniById(ctx, ids)
	} else {
		respEnis, err = service.DescribeEniByFilters(ctx, vpcId, subnetId, cvmId, sgId, name, desc, ipv4, tags)
	}

	if err != nil {
		return err
	}

	enis := make([]map[string]interface{}, 0, len(respEnis))
	eniIds := make([]string, 0, len(respEnis))

	for _, eni := range respEnis {
		ipv4s := make([]map[string]interface{}, 0, len(eni.PrivateIpAddressSet))
		for _, ipv4 := range eni.PrivateIpAddressSet {
			ipv4s = append(ipv4s, map[string]interface{}{
				"ip":          ipv4.PrivateIpAddress,
				"primary":     ipv4.Primary,
				"description": eni.NetworkInterfaceDescription,
			})
		}

		ipv6s := make([]map[string]interface{}, 0, len(eni.Ipv6AddressSet))

		for _, ipv6 := range eni.Ipv6AddressSet {
			ipv6s = append(ipv6s, map[string]interface{}{
				"address":           ipv6.Address,
				"primary":           ipv6.Primary,
				"address_id":        ipv6.AddressId,
				"description":       ipv6.Description,
				"is_wan_ip_blocked": ipv6.IsWanIpBlocked,
			})
		}

		sgs := make([]string, 0, len(eni.GroupSet))
		for _, sg := range eni.GroupSet {
			sgs = append(sgs, *sg)
		}

		respTags := make(map[string]string, len(eni.TagSet))
		for _, tag := range eni.TagSet {
			respTags[*tag.Key] = *tag.Value
		}

		eniIds = append(eniIds, *eni.NetworkInterfaceId)

		m := map[string]interface{}{
			"id":              eni.NetworkInterfaceId,
			"name":            eni.NetworkInterfaceName,
			"description":     eni.NetworkInterfaceDescription,
			"vpc_id":          eni.VpcId,
			"subnet_id":       eni.SubnetId,
			"primary":         eni.Primary,
			"mac":             eni.MacAddress,
			"state":           eni.State,
			"create_time":     eni.CreatedTime,
			"ipv4s":           ipv4s,
			"ipv6s":           ipv6s,
			"security_groups": sgs,
			"tags":            respTags,
			"cdc_id":          eni.CdcId,
		}

		if eni.Attachment != nil {
			m["instance_id"] = eni.Attachment.InstanceId
		}

		enis = append(enis, m)
	}

	_ = d.Set("enis", enis)
	d.SetId(helper.DataResourceIdsHash(eniIds))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), enis); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%v]",
				logId, output.(string), err)
			return err
		}
	}

	return nil
}
