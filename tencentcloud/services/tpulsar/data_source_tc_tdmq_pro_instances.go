package tpulsar

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctdmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmq"

	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTdmqProInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqProInstancesRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "查询 condition 过滤器。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "名称 过滤器 参数。",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "值",
						},
					},
				},
			},
			"instances": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "实例 信息 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 ID",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例名称",
						},
						"instance_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 版本",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例状态，0-creating，1-normal，2-isolating，3-destroyed，4-abnormal，5-delivery failure，6-allocation change，7-allocation failure。",
						},
						"config_display": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 配置 规格名称",
						},
						"max_tps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Peak TPS。",
						},
						"max_storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Storage 容量，（GB）。",
						},
						"expire_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例 过期时间，（毫秒）。",
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Automatic renewal mark，0 表示default state ( 用户 has 不 集合 它，该 是， initial state 是 manual renewal)，1 表示automatic renewal，2 表示that automatic renewal 是 不 指定 (用户 setting)。",
						},
						"pay_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "0-postpaid，1-prepaid。",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remarks注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"spec_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 Configuration ID。",
						},
						"scalable_tps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Elastic TPS outside specification注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 的 VPC注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet id注意：此字段可能返回 null，表示无法获取有效值。",
						},
						"max_band_width": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Peak 带宽. 单位：mbps。",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "标签列表",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签键",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "标签值",
									},
								},
							},
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
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

func dataSourceTencentCloudTdmqProInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tdmq_pro_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service   = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		instances []*tdmq.PulsarProInstance
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*tdmq.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := tdmq.Filter{}
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
		paramMap["filters"] = tmpSet
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTdmqProInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		instances = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instances))
	tmpList := make([]map[string]interface{}, 0, len(instances))

	if instances != nil {
		for _, pulsarProInstance := range instances {
			pulsarProInstanceMap := map[string]interface{}{}

			if pulsarProInstance.InstanceId != nil {
				pulsarProInstanceMap["instance_id"] = pulsarProInstance.InstanceId
			}

			if pulsarProInstance.InstanceName != nil {
				pulsarProInstanceMap["instance_name"] = pulsarProInstance.InstanceName
			}

			if pulsarProInstance.InstanceVersion != nil {
				pulsarProInstanceMap["instance_version"] = pulsarProInstance.InstanceVersion
			}

			if pulsarProInstance.Status != nil {
				pulsarProInstanceMap["status"] = pulsarProInstance.Status
			}

			if pulsarProInstance.ConfigDisplay != nil {
				pulsarProInstanceMap["config_display"] = pulsarProInstance.ConfigDisplay
			}

			if pulsarProInstance.MaxTps != nil {
				pulsarProInstanceMap["max_tps"] = pulsarProInstance.MaxTps
			}

			if pulsarProInstance.MaxStorage != nil {
				pulsarProInstanceMap["max_storage"] = pulsarProInstance.MaxStorage
			}

			if pulsarProInstance.ExpireTime != nil {
				pulsarProInstanceMap["expire_time"] = pulsarProInstance.ExpireTime
			}

			if pulsarProInstance.AutoRenewFlag != nil {
				pulsarProInstanceMap["auto_renew_flag"] = pulsarProInstance.AutoRenewFlag
			}

			if pulsarProInstance.PayMode != nil {
				pulsarProInstanceMap["pay_mode"] = pulsarProInstance.PayMode
			}

			if pulsarProInstance.Remark != nil {
				pulsarProInstanceMap["remark"] = pulsarProInstance.Remark
			}

			if pulsarProInstance.SpecName != nil {
				pulsarProInstanceMap["spec_name"] = pulsarProInstance.SpecName
			}

			if pulsarProInstance.ScalableTps != nil {
				pulsarProInstanceMap["scalable_tps"] = pulsarProInstance.ScalableTps
			}

			if pulsarProInstance.VpcId != nil {
				pulsarProInstanceMap["vpc_id"] = pulsarProInstance.VpcId
			}

			if pulsarProInstance.SubnetId != nil {
				pulsarProInstanceMap["subnet_id"] = pulsarProInstance.SubnetId
			}

			if pulsarProInstance.MaxBandWidth != nil {
				pulsarProInstanceMap["max_band_width"] = pulsarProInstance.MaxBandWidth
			}

			if pulsarProInstance.Tags != nil {
				tagsList := []interface{}{}
				for _, tags := range pulsarProInstance.Tags {
					tagsMap := map[string]interface{}{}

					if tags.TagKey != nil {
						tagsMap["tag_key"] = tags.TagKey
					}

					if tags.TagValue != nil {
						tagsMap["tag_value"] = tags.TagValue
					}

					tagsList = append(tagsList, tagsMap)
				}

				pulsarProInstanceMap["tags"] = tagsList
			}

			if pulsarProInstance.CreateTime != nil {
				pulsarProInstanceMap["create_time"] = pulsarProInstance.CreateTime
			}

			ids = append(ids, *pulsarProInstance.InstanceId)
			tmpList = append(tmpList, pulsarProInstanceMap)
		}

		_ = d.Set("instances", tmpList)
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
