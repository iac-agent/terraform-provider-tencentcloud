package vpn

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"context"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpnConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpnConnectionsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "名称 VPN 连接. 长度 的 character 是 limited 到 1-60。",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID VPN 连接。",
			},
			"vpn_gateway_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VPN 网关 ID VPN 连接。",
			},
			"customer_gateway_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Customer 网关 ID VPN 连接。",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID VPC。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 的 VPN 连接 到 是 queried。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// Computed values
			"connection_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information 列表 dedicated connections。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID VPN 连接。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 VPN 连接。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID VPC。",
						},
						"customer_gateway_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID customer 网关。",
						},
						"vpn_gateway_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID VPN 网关。",
						},
						"pre_share_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pre-shared 键 的 VPN 连接。",
						},
						"security_group_policy": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Security 组 策略 的 VPN 连接。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"local_cidr_block": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Local cidr block。",
									},
									"remote_cidr_block": {
										Type:        schema.TypeSet,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Remote cidr block 列表。",
									},
								},
							},
						},
						"ike_proto_encry_algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Proto encrypt algorithm 的 IKE operation 规格。",
						},
						"ike_proto_authen_algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Proto authenticate algorithm 的 IKE operation 规格。",
						},
						"ike_exchange_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Exchange 模式 的 IKE operation 规格。",
						},
						"ike_local_identity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Local identity 的 IKE operation 规格。",
						},
						"ike_remote_identity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remote identity 的 IKE operation 规格。",
						},
						"ike_local_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Local 地址 的 IKE operation 规格。",
						},
						"ike_remote_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remote 地址 的 IKE operation 规格。",
						},
						"ike_local_fqdn_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Local FQDN 名称 IKE operation 规格。",
						},
						"ike_remote_fqdn_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remote FQDN 名称 IKE operation 规格。",
						},
						"ike_dh_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DH 组名称 的 IKE operation 规格。",
						},
						"ike_sa_lifetime_seconds": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "SA lifetime 的 IKE operation 规格，单位 是 `second`。",
						},
						"ike_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "版本 的 IKE operation 规格。",
						},
						"ipsec_encrypt_algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Encrypt algorithm 的 IPSEC operation 规格。",
						},
						"ipsec_integrity_algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Integrity algorithm 的 IPSEC operation 规格。",
						},
						"ipsec_sa_lifetime_seconds": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "SA lifetime 的 IPSEC operation 规格，单位 是 `second`。",
						},
						"ipsec_pfs_dh_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "PFS DH 组名称 的 IPSEC operation 规格。",
						},
						"ipsec_sa_lifetime_traffic": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "SA lifetime 流量 的 IPSEC operation 规格，单位 是 `KB`。",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "A 列表 标签 用于associate different resources。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 VPN 连接。",
						},
						"vpn_proto": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Vpn proto 的 VPN 连接。",
						},
						"encrypt_proto": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Encrypt proto 的 VPN 连接。",
						},
						"route_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Route 类型 VPN 连接。",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State 的 VPN 连接。",
						},
						"net_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Net 状态 VPN 连接。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudVpnConnectionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpn_connections.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region

	request := vpc.NewDescribeVpnConnectionsRequest()

	params := make(map[string]string)
	if v, ok := d.GetOk("id"); ok {
		params["vpn-connection-id"] = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		params["vpn-connection-name"] = v.(string)
	}
	if v, ok := d.GetOk("vpn_gateway_id"); ok {
		params["vpn-gateway-id"] = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		params["vpc-id"] = v.(string)
	}
	if v, ok := d.GetOk("customer_gateway_id"); ok {
		params["customer-gateway-id"] = v.(string)
	}

	tags := helper.GetTags(d, "tags")

	request.Filters = make([]*vpc.Filter, 0, len(params))
	for k, v := range params {
		filter := &vpc.Filter{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, filter)
	}
	offset := uint64(0)
	request.Offset = &offset
	result := make([]*vpc.VpnConnection, 0)
	limit := uint64(svcvpc.VPN_DESCRIBE_LIMIT)
	request.Limit = &limit
	for {
		var response *vpc.DescribeVpnConnectionsResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeVpnConnections(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read VPN connection failed, reason:%s\n", logId, err.Error())
			return err
		} else {
			result = append(result, response.Response.VpnConnectionSet...)
			if len(response.Response.VpnConnectionSet) < svcvpc.VPN_DESCRIBE_LIMIT {
				break
			} else {
				offset = offset + limit
			}
		}
	}
	ids := make([]string, 0, len(result))
	connectionList := make([]map[string]interface{}, 0, len(result))
	for _, connection := range result {
		//tags
		respTags, err := tagService.DescribeResourceTags(ctx, "vpc", "vpnx", region, *connection.VpnConnectionId)
		if err != nil {
			return err
		}
		if len(tags) > 0 {
			if !reflect.DeepEqual(respTags, tags) {
				continue
			}
		}

		mapping := map[string]interface{}{
			"id":                         *connection.VpnConnectionId,
			"name":                       *connection.VpnConnectionName,
			"vpc_id":                     *connection.VpcId,
			"vpn_gateway_id":             *connection.VpnGatewayId,
			"customer_gateway_id":        *connection.CustomerGatewayId,
			"ike_proto_authen_algorithm": *connection.IKEOptionsSpecification.PropoAuthenAlgorithm,
			"ike_proto_encry_algorithm":  *connection.IKEOptionsSpecification.PropoEncryAlgorithm,
			"ike_exchange_mode":          *connection.IKEOptionsSpecification.ExchangeMode,
			"ike_dh_group_name":          *connection.IKEOptionsSpecification.DhGroupName,
			"ike_sa_lifetime_seconds":    int(*connection.IKEOptionsSpecification.IKESaLifetimeSeconds),
			"ike_version":                *connection.IKEOptionsSpecification.IKEVersion,
			"ike_local_identity":         *connection.IKEOptionsSpecification.LocalIdentity,
			"ike_local_address":          *connection.IKEOptionsSpecification.LocalAddress,
			"ike_local_fqdn_name":        *connection.IKEOptionsSpecification.LocalFqdnName,
			"ike_remote_identity":        *connection.IKEOptionsSpecification.RemoteIdentity,
			"ike_remote_address":         *connection.IKEOptionsSpecification.RemoteAddress,
			"ike_remote_fqdn_name":       *connection.IKEOptionsSpecification.RemoteFqdnName,
			"ipsec_sa_lifetime_seconds":  int(*connection.IPSECOptionsSpecification.IPSECSaLifetimeSeconds),
			"ipsec_encrypt_algorithm":    *connection.IPSECOptionsSpecification.EncryptAlgorithm,
			"ipsec_integrity_algorithm":  *connection.IPSECOptionsSpecification.IntegrityAlgorith,
			"ipsec_pfs_dh_group":         *connection.IPSECOptionsSpecification.PfsDhGroup,
			"ipsec_sa_lifetime_traffic":  int(*connection.IPSECOptionsSpecification.IPSECSaLifetimeTraffic),
			"security_group_policy":      svcvpc.FlattenVpnSPDList(connection.SecurityPolicyDatabaseSet),
			"net_status":                 *connection.NetStatus,
			"state":                      *connection.State,
			"create_time":                *connection.CreatedTime,
			"vpn_proto":                  *connection.VpnProto,
			"encrypt_proto":              *connection.EncryptProto,
			"route_type":                 *connection.RouteType,
			"tags":                       respTags,
		}
		connectionList = append(connectionList, mapping)
		ids = append(ids, *connection.VpnConnectionId)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("connection_list", connectionList); e != nil {
		log.Printf("[CRITAL]%s provider set connection list fail, reason:%s\n", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), connectionList); e != nil {
			return e
		}
	}

	return nil

}
