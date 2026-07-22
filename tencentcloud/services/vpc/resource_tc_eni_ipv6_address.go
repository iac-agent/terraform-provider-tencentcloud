package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudEniIpv6Address() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEniIpv6AddressCreate,
		Read:   resourceTencentCloudEniIpv6AddressRead,
		Delete: resourceTencentCloudEniIpv6AddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"network_interface_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "ENI 实例 `ID`，在 form 的 `eni-m6dyj72l`。",
			},

			"ipv6_addresses": {
				Optional:      true,
				Type:          schema.TypeSet,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"ipv6_address_count"},
				Description:   "指定 `IPv6` 地址 列表，up 到 10 可以 是 指定 在 时间. Combined 使用 input 参数 `Ipv6AddressCount` 到 calculate 配额. Mandatory 一个 使用 Ipv6AddressCount。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "`IPv6` 地址，在 form 的: `3402:4e00:20:100:0:8cd9:2a67:71f3`。",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "描述",
						},
						"primary": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "是否master `IP`。",
						},
						"address_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "`EIP` 实例 `ID`，such 作为:`eip-hxlqja90`。",
						},
						"is_wan_ip_blocked": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "是否public 网络 IP 是 blocked。",
						},
						"state": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "`IPv6` 地址 状态: `PENDING`: pending，`MIGRATING`: migrating，`DELETING`: deleting，`AVAILABLE`: 可用。",
						},
					},
				},
			},

			"ipv6_address_count": {
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Type:          schema.TypeInt,
				ConflictsWith: []string{"ipv6_addresses"},
				Description:   "数量 automatically assigned IPv6 addresses 和 总数 数量 私有 IP addresses 不能 exceed 配额. 此 should 是 combined 使用 input 参数 `ipv6_addresses` 对于 配额 calculation. At least 一个 的 them，either 此 或 'Ipv6Addresses'，必须 是 提供。",
			},
		},
	}
}

func resourceTencentCloudEniIpv6AddressCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eni_ipv6_address.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		ctx                = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		vpcService         = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request            = vpc.NewAssignIpv6AddressesRequest()
		response           = vpc.NewAssignIpv6AddressesResponse()
		networkInterfaceId string
	)

	if v, ok := d.GetOk("network_interface_id"); ok {
		request.NetworkInterfaceId = helper.String(v.(string))
		networkInterfaceId = v.(string)
	}

	if v, ok := d.GetOk("ipv6_addresses"); ok {
		for _, item := range v.(*schema.Set).List() {
			dMap := item.(map[string]interface{})
			ipv6Address := vpc.Ipv6Address{}
			if v, ok := dMap["address"]; ok {
				ipv6Address.Address = helper.String(v.(string))
			}

			if v, ok := dMap["primary"]; ok {
				ipv6Address.Primary = helper.Bool(v.(bool))
			}

			if v, ok := dMap["address_id"]; ok {
				ipv6Address.AddressId = helper.String(v.(string))
			}

			if v, ok := dMap["description"]; ok {
				ipv6Address.Description = helper.String(v.(string))
			}

			if v, ok := dMap["is_wan_ip_blocked"]; ok {
				ipv6Address.IsWanIpBlocked = helper.Bool(v.(bool))
			}

			if v, ok := dMap["state"]; ok {
				ipv6Address.State = helper.String(v.(string))
			}

			request.Ipv6Addresses = append(request.Ipv6Addresses, &ipv6Address)
		}
	}

	if v, ok := d.GetOkExists("ipv6_address_count"); ok {
		request.Ipv6AddressCount = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().AssignIpv6Addresses(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create vpc ipv6EniAddress failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create vpc ipv6EniAddress failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.Ipv6AddressSet == nil || len(response.Response.Ipv6AddressSet) < 1 {
		return fmt.Errorf("assign ipv6 addresses failed.")
	}

	d.SetId(networkInterfaceId)

	// wait
	if response.Response.RequestId != nil {
		err = vpcService.DescribeVpcTaskResult(ctx, response.Response.RequestId)
		if err != nil {
			return err
		}
	} else {
		time.Sleep(15 * time.Second)
	}

	return resourceTencentCloudEniIpv6AddressRead(d, meta)
}

func resourceTencentCloudEniIpv6AddressRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eni_ipv6_address.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		ctx                = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service            = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		networkInterfaceId = d.Id()
	)

	enis, err := service.DescribeEniById(ctx, []string{networkInterfaceId})
	if err != nil {
		return err
	}

	if len(enis) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcIpv6EniAddress` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	eni := enis[0]
	ipv6s := make([]map[string]interface{}, 0, len(eni.Ipv6AddressSet))
	for _, ipv6 := range eni.Ipv6AddressSet {
		ipv6s = append(ipv6s, map[string]interface{}{
			"address":           ipv6.Address,
			"primary":           ipv6.Primary,
			"address_id":        ipv6.AddressId,
			"description":       ipv6.Description,
			"is_wan_ip_blocked": ipv6.IsWanIpBlocked,
			"state":             ipv6.State,
		})
	}

	_ = d.Set("network_interface_id", networkInterfaceId)
	_ = d.Set("ipv6_addresses", ipv6s)
	_ = d.Set("ipv6_address_count", len(eni.Ipv6AddressSet))

	return nil
}

func resourceTencentCloudEniIpv6AddressDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eni_ipv6_address.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		ctx                = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service            = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		networkInterfaceId = d.Id()
	)

	enis, err := service.DescribeEniById(ctx, []string{networkInterfaceId})
	if err != nil {
		return err
	}

	if len(enis) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcIpv6EniAddress` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	eni := enis[0]
	ipv6s := make([]*string, 0, len(eni.Ipv6AddressSet))
	for _, ipv6 := range eni.Ipv6AddressSet {
		ipv6s = append(ipv6s, ipv6.Address)
	}

	if err := service.DeleteEniIpv6AddressById(ctx, networkInterfaceId, ipv6s); err != nil {
		return err
	}

	return nil
}
