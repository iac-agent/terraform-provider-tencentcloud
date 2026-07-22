package cvm

import (
	"context"
	"errors"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

var (
	errEIPNotFound = errors.New("eip not found")
)

func DataSourceTencentCloudEip() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source has been deprecated in Terraform TencentCloud provider version 1.20.0. Please use 'tencentcloud_eips' instead.",
		Read:               dataSourceTencentCloudEipRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "One 或 more 名称/值 pairs 到 过滤器。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "键 的 过滤器，有效 keys: `地址-ID`,`地址-名称`,`地址-ip`。",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "值 的 过滤器。",
						},
					},
				},
			},
			"include_arrears": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "是否IP 是 arrears。",
			},
			"include_blocked": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "是否IP 是 blocked。",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An EIP ID indicate uniqueness 的 certain EIP， 其中 可以 是 用于instance binding 或 网络 interface binding。",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An 公网 IP 地址 对于 EIP。",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "状态 EIP，there 是 several 状态 like `BIND`，`UNBIND`，和 `BIND_ENI`。",
			},
		},
	}
}

func dataSourceTencentCloudEipRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_eip.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	vpcService := svcvpc.NewVpcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	filter := make(map[string][]string)
	filters, ok := d.GetOk("filter")
	if ok {
		for _, v := range filters.(*schema.Set).List() {
			vv := v.(map[string]interface{})
			name := vv["name"].(string)
			filter[name] = []string{}
			for _, vvv := range vv["values"].([]interface{}) {
				filter[name] = append(filter[name], vvv.(string))
			}
		}
	}

	var eips []*vpc.Address
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		eips, errRet = vpcService.DescribeEipByFilter(ctx, filter)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	includeArrears := false
	if v, ok := d.GetOk("include_arrears"); ok {
		includeArrears = v.(bool)
	}
	includeBlocked := false
	if v, ok := d.GetOk("include_blocked"); ok {
		includeBlocked = v.(bool)
	}

	if len(eips) == 0 {
		return errEIPNotFound
	}

	var filteredEips []*vpc.Address
	for _, eip := range eips {
		if *eip.IsArrears && !includeArrears {
			continue
		}
		if *eip.IsBlocked && !includeBlocked {
			continue
		}
		filteredEips = append(filteredEips, eip)
	}

	if len(filteredEips) == 0 {
		return errEIPNotFound
	}

	eip := filteredEips[0]
	d.SetId(*eip.AddressId)
	_ = d.Set("public_ip", *eip.AddressIp)
	_ = d.Set("status", *eip.AddressStatus)

	return nil
}
