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

func DataSourceTencentCloudVpnGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpnGatewaysRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "名称 VPN gateway. The length of character is limited to 1-60。",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID VPN gateway。",
			},
			"public_ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateIp,
				Description:  "Public ip 地址 of the VPN gateway。",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "可用区 of the VPN gateway。",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID VPC。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 of the VPN gateway to be queried。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// Computed values
			"gateway_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information 列表 the dedicated gateways。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID VPN gateway。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 VPN gateway。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID VPC。",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum public network output bandwidth of VPN gateway (unit: Mbps)。",
						},
						"public_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public ip of the VPN gateway。",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 gateway instance。",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State of the VPN gateway。",
						},
						"prepaid_renew_flag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Flag 表示是否renew or not。",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Charge 类型 VPN gateway。",
						},
						"expired_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "过期时间 of the VPN gateway when 计费类型 is `PREPAID`。",
						},
						"is_address_blocked": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "表示是否ip 地址 is blocked。",
						},
						"new_purchase_plan": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The plan of new purchase。",
						},
						"restrict_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Restrict state of VPN gateway。",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用区 of the VPN gateway。",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "A 列表 标签 用于associate different resources。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 of the VPN gateway。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudVpnGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpn_gateways.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region

	request := vpc.NewDescribeVpnGatewaysRequest()

	params := make(map[string]string)
	if v, ok := d.GetOk("id"); ok {
		params["vpn-gateway-id"] = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		params["vpn-gateway-name"] = v.(string)
	}
	if v, ok := d.GetOk("public_ip_address"); ok {
		params["public-ip-address"] = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		params["vpc-id"] = v.(string)
	}
	if v, ok := d.GetOk("zone"); ok {
		params["zone"] = v.(string)
	}

	tags := helper.GetTags(d, "tags")

	request.Filters = make([]*vpc.FilterObject, 0, len(params))
	for k, v := range params {
		filter := &vpc.FilterObject{
			Name:   helper.String(k),
			Values: []*string{helper.String(v)},
		}
		request.Filters = append(request.Filters, filter)
	}
	offset := uint64(0)
	request.Offset = &offset
	result := make([]*vpc.VpnGateway, 0)
	limit := uint64(svcvpc.VPN_DESCRIBE_LIMIT)
	request.Limit = &limit
	for {
		var response *vpc.DescribeVpnGatewaysResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeVpnGateways(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read VPN gateway failed, reason:%s\n ", logId, err.Error())
			return err
		} else {
			result = append(result, response.Response.VpnGatewaySet...)
			if len(response.Response.VpnGatewaySet) < svcvpc.VPN_DESCRIBE_LIMIT {
				break
			} else {
				offset = offset + limit
			}
		}
	}
	ids := make([]string, 0, len(result))
	gatewayList := make([]map[string]interface{}, 0, len(result))
	for _, gateway := range result {
		//tags
		respTags, err := tagService.DescribeResourceTags(ctx, "vpc", "vpngw", region, *gateway.VpnGatewayId)
		if err != nil {
			return err
		}
		if len(tags) > 0 {
			if !reflect.DeepEqual(respTags, tags) {
				continue
			}
		}

		mapping := map[string]interface{}{
			"id":                 *gateway.VpnGatewayId,
			"name":               *gateway.VpnGatewayName,
			"public_ip_address":  *gateway.PublicIpAddress,
			"create_time":        *gateway.CreatedTime,
			"prepaid_renew_flag": *gateway.RenewFlag,
			"state":              *gateway.State,
			"charge_type":        *gateway.InstanceChargeType,
			"expired_time":       *gateway.ExpiredTime,
			"is_address_blocked": *gateway.IsAddressBlocked,
			"bandwidth":          int(*gateway.InternetMaxBandwidthOut),
			"new_purchase_plan":  *gateway.NewPurchasePlan,
			"restrict_state":     *gateway.RestrictState,
			"zone":               *gateway.Zone,
			"type":               *gateway.Type,
			"tags":               respTags,
		}
		gatewayList = append(gatewayList, mapping)
		ids = append(ids, *gateway.VpnGatewayId)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("gateway_list", gatewayList); e != nil {
		log.Printf("[CRITAL]%s provider set gateway list fail, reason:%s\n ", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), gatewayList); e != nil {
			return e
		}
	}

	return nil

}
