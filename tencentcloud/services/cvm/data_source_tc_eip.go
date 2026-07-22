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
				Description: "One or more 名称/值 pairs to filter。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "键 of the filter，valid keys: `地址-id`,`地址-名称`,`地址-ip`。",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "值 of the filter。",
						},
					},
				},
			},
			"include_arrears": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "是否IP is arrears。",
			},
			"include_blocked": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "是否IP is blocked。",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An EIP id indicate the uniqueness of a certain EIP， which can be 用于instance binding or network interface binding。",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An 公网 IP 地址 for the EIP。",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The 状态 EIP，there are several 状态 like `BIND`，`UNBIND`，and `BIND_ENI`。",
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
