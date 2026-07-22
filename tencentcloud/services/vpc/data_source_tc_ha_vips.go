package vpc

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudHaVips() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudHaVipsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "名称 HA VIP 长度 的 character 是 limited 到 1-60。",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID HA VIP 到 是 queried。",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VPC ID HA VIP 到 是 queried。",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "子网 ID HA VIP 到 是 queried。",
			},
			"address_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateIp,
				Description:  "EIP 的 HA VIP 到 是 queried。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// Computed values
			"ha_vip_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information 列表 dedicated HA VIPs。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID HA VIP",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 HA VIP",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC ID。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "子网 ID",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Virtual IP 地址，它 必须 不 是 occupied 和 在 此 VPC 网络 segment. 如果 不 集合，它 将 是 assigned after 资源 创建 automatically。",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State 的 HA VIP 有效值：`AVAILABLE`，`UNBIND`。",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID 该 是 associated。",
						},
						"network_interface_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network interface ID 该 是 associated。",
						},
						"address_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "EIP 该 是 associated。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 HA VIP",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudHaVipsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ha_vips.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vpc.NewDescribeHaVipsRequest()

	params := make(map[string]string)
	if v, ok := d.GetOk("id"); ok {
		params["havip-id"] = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		params["havip-name"] = v.(string)
	}
	if v, ok := d.GetOk("address_ip"); ok {
		params["address-ip"] = v.(string)
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		params["subnet-id"] = v.(string)
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
	result := make([]*vpc.HaVip, 0)
	limit := uint64(HAVIP_DESCRIBE_LIMIT)
	request.Limit = &limit
	for {
		var response *vpc.DescribeHaVipsResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DescribeHaVips(request)
			if e != nil {
				return tccommon.RetryError(errors.WithStack(e))
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read HA VIP failed, reason:%+v", logId, err)
			return err
		} else {
			result = append(result, response.Response.HaVipSet...)
			if len(response.Response.HaVipSet) < HAVIP_DESCRIBE_LIMIT {
				break
			} else {
				offset = offset + limit
			}
		}
	}
	ids := make([]string, 0, len(result))
	haVipList := make([]map[string]interface{}, 0, len(result))
	for _, haVip := range result {
		mapping := map[string]interface{}{
			"id":          *haVip.HaVipId,
			"vip":         *haVip.Vip,
			"name":        *haVip.HaVipName,
			"state":       *haVip.State,
			"vpc_id":      *haVip.VpcId,
			"subnet_id":   *haVip.SubnetId,
			"create_time": *haVip.CreatedTime,
		}
		if haVip.NetworkInterfaceId != nil {
			mapping["network_interface_id"] = *haVip.NetworkInterfaceId
		}
		if haVip.AddressIp != nil {
			mapping["address_ip"] = *haVip.AddressIp
		}
		if haVip.InstanceId != nil {
			mapping["instance_id"] = *haVip.InstanceId
		}
		haVipList = append(haVipList, mapping)
		ids = append(ids, *haVip.HaVipId)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("ha_vip_list", haVipList); e != nil {
		log.Printf("[CRITAL]%s provider set haVip list fail, reason:%s\n", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), haVipList); e != nil {
			return e
		}
	}

	return nil

}
