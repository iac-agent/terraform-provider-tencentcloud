package vpc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpcNetworkInterfaceLimit() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcNetworkInterfaceLimitRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID CVM 实例 或 ENI 到 查询。",
			},

			"eni_quantity": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Quota 的 ENIs mounted 到 CVM 实例 在 standard way。",
			},

			"eni_private_ip_address_quantity": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Quota 的 IP addresses 该 可以 是 allocated 到 each standard-mounted ENI。",
			},

			"extend_eni_quantity": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Quota 的 ENIs mounted 到 CVM 实例 作为 extensionNote: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
			},

			"extend_eni_private_ip_address_quantity": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Quota 的 IP addresses 该 可以 是 allocated 到 each extension-mounted ENI.注意: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
			},

			"sub_eni_quantity": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "配额 的 relayed ENIsNote: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
			},

			"sub_eni_private_ip_address_quantity": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "配额 的 IPs 该 可以 是 assigned 到 each relayed ENI.注意: 此 字段 可能 返回 `null`，indicating 该 无 有效 值 可以 是 获取。",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
		},
	}
}

func dataSourceTencentCloudVpcNetworkInterfaceLimitRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_network_interface_limit.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var instanceId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var networkInterfaceLimit *vpc.DescribeNetworkInterfaceLimitResponseParams

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcNetworkInterfaceLimit(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		networkInterfaceLimit = result
		return nil
	})
	if err != nil {
		return err
	}

	limitMap := map[string]interface{}{}

	if networkInterfaceLimit.EniQuantity != nil {
		_ = d.Set("eni_quantity", networkInterfaceLimit.EniQuantity)
		limitMap["eni_quantity"] = networkInterfaceLimit.EniQuantity
	}

	if networkInterfaceLimit.EniPrivateIpAddressQuantity != nil {
		_ = d.Set("eni_private_ip_address_quantity", networkInterfaceLimit.EniPrivateIpAddressQuantity)
		limitMap["eni_private_ip_address_quantity"] = networkInterfaceLimit.EniPrivateIpAddressQuantity
	}

	if networkInterfaceLimit.ExtendEniQuantity != nil {
		_ = d.Set("extend_eni_quantity", networkInterfaceLimit.ExtendEniQuantity)
		limitMap["extend_eni_quantity"] = networkInterfaceLimit.ExtendEniQuantity
	}

	if networkInterfaceLimit.ExtendEniPrivateIpAddressQuantity != nil {
		_ = d.Set("extend_eni_private_ip_address_quantity", networkInterfaceLimit.ExtendEniPrivateIpAddressQuantity)
		limitMap["extend_eni_private_ip_address_quantity"] = networkInterfaceLimit.ExtendEniPrivateIpAddressQuantity
	}

	if networkInterfaceLimit.SubEniQuantity != nil {
		_ = d.Set("sub_eni_quantity", networkInterfaceLimit.SubEniQuantity)
		limitMap["sub_eni_quantity"] = networkInterfaceLimit.SubEniQuantity
	}

	if networkInterfaceLimit.SubEniPrivateIpAddressQuantity != nil {
		_ = d.Set("sub_eni_private_ip_address_quantity", networkInterfaceLimit.SubEniPrivateIpAddressQuantity)
		limitMap["sub_eni_private_ip_address_quantity"] = networkInterfaceLimit.SubEniPrivateIpAddressQuantity
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), limitMap); e != nil {
			return e
		}
	}
	return nil
}
