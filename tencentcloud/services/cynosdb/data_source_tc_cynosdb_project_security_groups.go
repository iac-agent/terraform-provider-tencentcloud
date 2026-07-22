package cynosdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCynosdbProjectSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbProjectSecurityGroupsRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "项目编号。",
			},
			"search_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "搜索关键词。",
			},
			"groups": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "安全组详细信息。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "项目编号。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间，时间格式：yyyy mm dd hh:mm:ss。",
						},
						"inbound": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "入站规则。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "行动。",
									},
									"cidr_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CidrIp。",
									},
									"port_range": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "端口范围。",
									},
									"ip_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IP协议。",
									},
									"service_module": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "服务模块。",
									},
									"address_module": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "地址模块。",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID。",
									},
									"desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "描述。",
									},
								},
							},
						},
						"outbound": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "出站规则。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "行动。",
									},
									"cidr_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cidr Ip。",
									},
									"port_range": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "端口范围。",
									},
									"ip_protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IP协议。",
									},
									"service_module": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "服务模块。",
									},
									"address_module": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "地址模块。",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID。",
									},
									"desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "描述。",
									},
								},
							},
						},
						"security_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "安全组ID。",
						},
						"security_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "安全组名称。",
						},
						"security_group_remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "安全组注释。",
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

func dataSourceTencentCloudCynosdbProjectSecurityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cynosdb_project_security_groups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		groups  []*cynosdb.SecurityGroup
	)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("project_id"); v != nil {
		paramMap["ProjectId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("search_key"); ok {
		paramMap["SearchKey"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbProjectSecurityGroupsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		groups = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(groups))
	tmpList := make([]map[string]interface{}, 0, len(groups))

	if groups != nil {
		for _, securityGroup := range groups {
			securityGroupMap := map[string]interface{}{}

			if securityGroup.ProjectId != nil {
				securityGroupMap["project_id"] = securityGroup.ProjectId
			}

			if securityGroup.CreateTime != nil {
				securityGroupMap["create_time"] = securityGroup.CreateTime
			}

			if securityGroup.Inbound != nil {
				inboundList := []interface{}{}
				for _, inbound := range securityGroup.Inbound {
					inboundMap := map[string]interface{}{}

					if inbound.Action != nil {
						inboundMap["action"] = inbound.Action
					}

					if inbound.CidrIp != nil {
						inboundMap["cidr_ip"] = inbound.CidrIp
					}

					if inbound.PortRange != nil {
						inboundMap["port_range"] = inbound.PortRange
					}

					if inbound.IpProtocol != nil {
						inboundMap["ip_protocol"] = inbound.IpProtocol
					}

					if inbound.ServiceModule != nil {
						inboundMap["service_module"] = inbound.ServiceModule
					}

					if inbound.AddressModule != nil {
						inboundMap["address_module"] = inbound.AddressModule
					}

					if inbound.Id != nil {
						inboundMap["id"] = inbound.Id
					}

					if inbound.Desc != nil {
						inboundMap["desc"] = inbound.Desc
					}

					inboundList = append(inboundList, inboundMap)
				}

				securityGroupMap["inbound"] = inboundList
			}

			if securityGroup.Outbound != nil {
				outboundList := []interface{}{}
				for _, outbound := range securityGroup.Outbound {
					outboundMap := map[string]interface{}{}

					if outbound.Action != nil {
						outboundMap["action"] = outbound.Action
					}

					if outbound.CidrIp != nil {
						outboundMap["cidr_ip"] = outbound.CidrIp
					}

					if outbound.PortRange != nil {
						outboundMap["port_range"] = outbound.PortRange
					}

					if outbound.IpProtocol != nil {
						outboundMap["ip_protocol"] = outbound.IpProtocol
					}

					if outbound.ServiceModule != nil {
						outboundMap["service_module"] = outbound.ServiceModule
					}

					if outbound.AddressModule != nil {
						outboundMap["address_module"] = outbound.AddressModule
					}

					if outbound.Id != nil {
						outboundMap["id"] = outbound.Id
					}

					if outbound.Desc != nil {
						outboundMap["desc"] = outbound.Desc
					}

					outboundList = append(outboundList, outboundMap)
				}

				securityGroupMap["outbound"] = outboundList
			}

			if securityGroup.SecurityGroupId != nil {
				securityGroupMap["security_group_id"] = securityGroup.SecurityGroupId
			}

			if securityGroup.SecurityGroupName != nil {
				securityGroupMap["security_group_name"] = securityGroup.SecurityGroupName
			}

			if securityGroup.SecurityGroupRemark != nil {
				securityGroupMap["security_group_remark"] = securityGroup.SecurityGroupRemark
			}

			ids = append(ids, *securityGroup.SecurityGroupId)
			tmpList = append(tmpList, securityGroupMap)
		}

		_ = d.Set("groups", tmpList)
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
