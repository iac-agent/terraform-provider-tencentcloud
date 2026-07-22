package vpc

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudNatGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudNatGatewaysRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID VPC。",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 NAT 网关。",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID NAT 网关。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// Computed values
			"nats": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information 列表 dedicated NATs。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID NAT 网关。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID VPC。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 NAT 网关。",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State 的 NAT 网关。",
						},
						"max_concurrent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "upper 限制 的 concurrent 连接 的 NAT 网关， 可用 值 include: 1000000,3000000,10000000. 默认为 1000000。",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最大 公有 网络 output 带宽 的 NAT 网关 (单位: Mbps)， 可用 值 include: 20,50,100,200,500,1000,2000,5000. 默认为 100。",
						},
						"assigned_eip_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "EIP IP 地址 集合 bound 到 网关. 值 的 在 least 1。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 NAT 网关。",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "可用 标签 within 此 NAT 网关。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudNatGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_nat_gateways.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := vpc.NewDescribeNatGatewaysRequest()

	params := make(map[string]string)
	if v, ok := d.GetOk("id"); ok {
		params["nat-gateway-id"] = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		params["nat-gateway-name"] = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		params["vpc-id"] = v.(string)
	}
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
	result := make([]*vpc.NatGateway, 0)
	limit := uint64(NAT_DESCRIBE_LIMIT)
	request.Limit = &limit
	for {
		var response *vpc.DescribeNatGatewaysResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeNatGateways(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read NAT gateway failed, reason:%s\n", logId, err.Error())
			return err
		} else {
			result = append(result, response.Response.NatGatewaySet...)
			if len(response.Response.NatGatewaySet) < NAT_DESCRIBE_LIMIT {
				break
			} else {
				offset = offset + limit
				request.Offset = &offset
			}
		}
	}
	ids := make([]string, 0, len(result))
	natList := make([]map[string]interface{}, 0, len(result))
	for _, nat := range result {
		mapping := map[string]interface{}{
			"id":               *nat.NatGatewayId,
			"vpc_id":           *nat.VpcId,
			"name":             *nat.NatGatewayName,
			"max_concurrent":   *nat.MaxConcurrentConnection,
			"bandwidth":        *nat.InternetMaxBandwidthOut,
			"state":            *nat.State,
			"assigned_eip_set": flattenAddressList((*nat).PublicIpAddressSet),
			"create_time":      *nat.CreatedTime,
		}
		if nat.TagSet != nil {
			tags := make(map[string]interface{}, len(nat.TagSet))
			for _, t := range nat.TagSet {
				tags[*t.Key] = *t.Value
			}
			mapping["tags"] = tags
		}
		natList = append(natList, mapping)
		ids = append(ids, *nat.NatGatewayId)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("nats", natList); e != nil {
		log.Printf("[CRITAL]%s provider set NAT list fail, reason:%s\n", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), natList); e != nil {
			return e
		}
	}

	return nil

}
