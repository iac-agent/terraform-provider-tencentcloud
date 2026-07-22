package cfw

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfwv20190904 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCfwCcnInstanceRegionStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCfwCcnInstanceRegionStatusRead,
		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CCN ID。",
			},

			"instance_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "列表 实例 IDs associated 使用 CCN 对于 querying 流量 steering 网络 部署 状态",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"routing_mode": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Traffic steering routing 方法，0: multi-路由 表，1: 策略 routing。",
			},

			"region_fw_status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "列表 regional firewall 流量 steering 网络 状态",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "地域",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Traffic steering 网络 部署 状态\n1. `NotDeployed` Firewall 集群 不 deployed.\n2. `Deployed` Firewall 集群 deployed，但 流量 steering 网络 不 创建.\n3. `Auto` Firewall 集群 deployed，和 流量 steering 网络 创建 使用 automatically selected 网络 segment.\n4. `Custom` Firewall 集群 deployed，和 流量 steering 网络 创建 使用 用户-defined 网络 segment。",
						},
						"cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CIDR 的 流量 steering 网络，空 如果 流量 steering 网络 是 不 deployed。",
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

func dataSourceTencentCloudCfwCcnInstanceRegionStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cfw_ccn_instance_region_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		ccnId   string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("ccn_id"); ok {
		paramMap["CcnId"] = helper.String(v.(string))
		ccnId = v.(string)
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsList := []*string{}
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			instanceIdsList = append(instanceIdsList, helper.String(instanceIds))
		}

		paramMap["InstanceIds"] = instanceIdsList
	}

	if v, ok := d.GetOkExists("routing_mode"); ok {
		paramMap["RoutingMode"] = helper.IntUint64(v.(int))
	}

	var respData []*cfwv20190904.RegionFwStatus
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCfwCcnInstanceRegionStatusByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	regionFwStatusList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, regionFwStatus := range respData {
			regionFwStatusMap := map[string]interface{}{}
			if regionFwStatus.Region != nil {
				regionFwStatusMap["region"] = regionFwStatus.Region
			}

			if regionFwStatus.Status != nil {
				regionFwStatusMap["status"] = regionFwStatus.Status
			}

			if regionFwStatus.Cidr != nil {
				regionFwStatusMap["cidr"] = regionFwStatus.Cidr
			}

			regionFwStatusList = append(regionFwStatusList, regionFwStatusMap)
		}

		_ = d.Set("region_fw_status", regionFwStatusList)
	}

	d.SetId(ccnId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
