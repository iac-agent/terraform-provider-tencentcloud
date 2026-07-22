package vpc

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudNats() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource has been deprecated in Terraform TencentCloud provider version 1.18.0. Please use 'tencentcloud_nat_gateways' instead.",
		Read:               dataSourceTencentCloudNatsRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID 对于 NAT Gateway。",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "私有网络 ID 对于 NAT Gateway。",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "名称 对于 NAT Gateway。",
			},
			"state": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "NAT 网关 状态 有效值：0，1，2. 0: Running，1: Unavailable，2: Be 在 arrears 和 out 的 服务。",
			},
			"max_concurrent": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "upper 限制 的 concurrent 连接 的 NAT 网关，对于 示例: `1000000`，`3000000`，`10000000`。",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "最大 公有 网络 output 带宽 的 网关 (单位: Mbps)，对于 示例: `10`，`20`，`50`，`100`，`200`，`500`，`1000`，`2000`，`5000`。",
			},
			"nats": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information 列表 dedicated tunnels。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 对于 NAT Gateway。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "私有网络 ID 对于 NAT Gateway。",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 对于 NAT Gateway。",
						},
						"state": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "NAT 网关 状态，`0`: Running，`1`: Unavailable，`2`: Be 在 arrears 和 out 的 服务。",
						},
						"max_concurrent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "upper 限制 的 concurrent 连接 的 NAT 网关，对于 示例: `1000000`，`3000000`，`10000000`。",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最大 公有 网络 output 带宽 的 网关 (单位: Mbps)，对于 示例: `10`，`20`，`50`，`100`，`200`，`500`，`1000`，`2000`，`5000`。",
						},
						"assigned_eip_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Elastic IP arrays bound 到 网关。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 NAT 网关。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudNatsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_nats.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := vpc.NewDescribeNatGatewaysRequest()
	request.Offset = helper.Uint64(0)
	request.Limit = helper.Uint64(100)

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
		log.Printf("[CRITAL]%s read nat gateway failed, reason:%s\n ", logId, err.Error())
		return err
	}

	ids := make([]string, 0, len(response.Response.NatGatewaySet))
	natList := make([]map[string]interface{}, 0, len(response.Response.NatGatewaySet))
	for _, nat := range response.Response.NatGatewaySet {
		networkState := 0
		switch *nat.NetworkState {
		case "AVAILABLE":
			networkState = 0
		case "UNAVAILABLE":
			networkState = 1
		case "INSUFFICIENT":
			networkState = 2
		}

		if state, ok := d.GetOk("state"); ok && networkState != state.(int) {
			continue
		}
		if max_concurrent, ok := d.GetOk("max_concurrent"); ok && *nat.MaxConcurrentConnection != uint64(max_concurrent.(int)) {
			continue
		}
		if bandwidth, ok := d.GetOk("bandwidth"); ok && *nat.InternetMaxBandwidthOut != uint64(bandwidth.(int)) {
			continue
		}

		mapping := map[string]interface{}{
			"id":               *nat.NatGatewayId,
			"vpc_id":           *nat.VpcId,
			"name":             *nat.NatGatewayName,
			"max_concurrent":   *nat.MaxConcurrentConnection,
			"bandwidth":        *nat.InternetMaxBandwidthOut,
			"state":            networkState,
			"assigned_eip_set": flattenAddressList((*nat).PublicIpAddressSet),
			"create_time":      *nat.CreatedTime,
		}
		natList = append(natList, mapping)
		ids = append(ids, *nat.NatGatewayId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("nats", natList); e != nil {
		log.Printf("[CRITAL]%s provider set clb list fail, reason:%s\n ", logId, e.Error())
		return e
	}

	return nil
}
