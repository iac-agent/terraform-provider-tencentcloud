package cfw

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCfwNatFwSwitches() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCfwNatFwSwitchesRead,
		Schema: map[string]*schema.Schema{
			"nat_ins_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "过滤器 NAT firewall 实例 到 其中 NAT firewall 子网 switch belongs。",
			},
			"status": {
				Optional:    true,
				Type:        schema.TypeInt,
				Deprecated:  "It has been deprecated from version 1.82.37. Please use `enable` instead.",
				Description: "Switch 状态，1 open; 0 close。",
			},
			"enable": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Switch 启用 状态，1 open; 0 close。",
			},
			"data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "NAT border firewall switch 列表 数据。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet ID。",
						},
						"subnet_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet 名称",
						},
						"subnet_cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IPv4 CIDR。",
						},
						"route_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Route ID。",
						},
						"route_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Route 名称",
						},
						"cvm_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cvm Num。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "私有网络 ID",
						},
						"vpc_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Vpc 名称",
						},
						"enable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Effective 状态",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Switch 状态",
						},
						"nat_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "NAT gatway ID。",
						},
						"nat_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "NAT gatway 名称",
						},
						"nat_ins_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "NAT firewall 实例 ID。",
						},
						"nat_ins_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "NAT firewall 实例名称",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域",
						},
						"abnormal": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "是否switch 是 abnormal，0: normal，1: abnormal。",
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

func dataSourceTencentCloudCfwNatFwSwitchesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cfw_nat_fw_switches.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		data    []*cfw.NatSwitchListData
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("nat_ins_id"); ok {
		paramMap["NatInsId"] = v.(string)
	}

	if v, ok := d.GetOkExists("status"); ok {
		paramMap["Status"] = v.(int)
	}

	if v, ok := d.GetOkExists("enable"); ok {
		paramMap["Enable"] = v.(int)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCfwNatFwSwitchesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		data = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(data))
	if data != nil {
		for _, natSwitchListData := range data {
			natSwitchListDataMap := map[string]interface{}{}
			if natSwitchListData.Id != nil {
				natSwitchListDataMap["id"] = natSwitchListData.Id
			}

			if natSwitchListData.SubnetId != nil {
				natSwitchListDataMap["subnet_id"] = natSwitchListData.SubnetId
			}

			if natSwitchListData.SubnetName != nil {
				natSwitchListDataMap["subnet_name"] = natSwitchListData.SubnetName
			}

			if natSwitchListData.SubnetCidr != nil {
				natSwitchListDataMap["subnet_cidr"] = natSwitchListData.SubnetCidr
			}

			if natSwitchListData.RouteId != nil {
				natSwitchListDataMap["route_id"] = natSwitchListData.RouteId
			}

			if natSwitchListData.RouteName != nil {
				natSwitchListDataMap["route_name"] = natSwitchListData.RouteName
			}

			if natSwitchListData.CvmNum != nil {
				natSwitchListDataMap["cvm_num"] = natSwitchListData.CvmNum
			}

			if natSwitchListData.VpcId != nil {
				natSwitchListDataMap["vpc_id"] = natSwitchListData.VpcId
			}

			if natSwitchListData.VpcName != nil {
				natSwitchListDataMap["vpc_name"] = natSwitchListData.VpcName
			}

			if natSwitchListData.Enable != nil {
				natSwitchListDataMap["enable"] = natSwitchListData.Enable
			}

			if natSwitchListData.Status != nil {
				natSwitchListDataMap["status"] = natSwitchListData.Status
			}

			if natSwitchListData.NatId != nil {
				natSwitchListDataMap["nat_id"] = natSwitchListData.NatId
			}

			if natSwitchListData.NatName != nil {
				natSwitchListDataMap["nat_name"] = natSwitchListData.NatName
			}

			if natSwitchListData.NatInsId != nil {
				natSwitchListDataMap["nat_ins_id"] = natSwitchListData.NatInsId
			}

			if natSwitchListData.NatInsName != nil {
				natSwitchListDataMap["nat_ins_name"] = natSwitchListData.NatInsName
			}

			if natSwitchListData.Region != nil {
				natSwitchListDataMap["region"] = natSwitchListData.Region
			}

			if natSwitchListData.Abnormal != nil {
				natSwitchListDataMap["abnormal"] = natSwitchListData.Abnormal
			}

			tmpList = append(tmpList, natSwitchListDataMap)
		}

		_ = d.Set("data", tmpList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
