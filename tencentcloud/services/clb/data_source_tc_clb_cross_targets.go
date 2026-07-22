package clb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClbCrossTargets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbCrossTargetsRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "查询云服务器和弹性网卡的过滤条件：vpc-id - 字符串 - 必填：否 - （过滤条件）按VPC ID过滤，例如vpc-12345678。 ip - String - 必填：否 - （过滤条件）按真实服务器IP过滤，如192.168.0.1。 listener-id - String - 必填：否 - （过滤条件）按监听器 ID 过滤，如 lbl-12345678。 location-id - String - 必填：否 - （过滤条件）按七层监听的转发规则ID过滤，如loc-12345678。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤器名称。",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "过滤值。",
						},
					},
				},
			},

			"cross_target_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "交叉目标设定。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"local_vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CLB实例的VPC ID。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CVM或弹性网卡实例的VPC ID。",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CVM或弹性网卡实例的IP地址。",
						},
						"vpc_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CVM或弹性网卡实例的VPC名称。",
						},
						"eni_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CVM实例的网卡ID。",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CVM实例ID。注意：该字段可能返回null，表示未找到有效值。",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CVM 实例的名称。注意：该字段可能返回 null，表示未找到有效值。",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CVM或弹性网卡实例所属地域。",
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

func dataSourceTencentCloudClbCrossTargetsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clb_cross_targets.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*clb.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := clb.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["Filters"] = tmpSet
	}

	service := ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var crossTargetSet []*clb.CrossTargets

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClbCrossTargetsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		crossTargetSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(crossTargetSet))
	tmpList := make([]map[string]interface{}, 0, len(crossTargetSet))

	if crossTargetSet != nil {
		for _, crossTargets := range crossTargetSet {
			crossTargetsMap := map[string]interface{}{}

			if crossTargets.LocalVpcId != nil {
				crossTargetsMap["local_vpc_id"] = crossTargets.LocalVpcId
			}

			if crossTargets.VpcId != nil {
				crossTargetsMap["vpc_id"] = crossTargets.VpcId
			}

			if crossTargets.IP != nil {
				crossTargetsMap["ip"] = crossTargets.IP
			}

			if crossTargets.VpcName != nil {
				crossTargetsMap["vpc_name"] = crossTargets.VpcName
			}

			if crossTargets.EniId != nil {
				crossTargetsMap["eni_id"] = crossTargets.EniId
			}

			if crossTargets.InstanceId != nil {
				crossTargetsMap["instance_id"] = crossTargets.InstanceId
			}

			if crossTargets.InstanceName != nil {
				crossTargetsMap["instance_name"] = crossTargets.InstanceName
			}

			if crossTargets.Region != nil {
				crossTargetsMap["region"] = crossTargets.Region
			}

			ids = append(ids, *crossTargets.VpcId)
			tmpList = append(tmpList, crossTargetsMap)
		}

		_ = d.Set("cross_target_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
