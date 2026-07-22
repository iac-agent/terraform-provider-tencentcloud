package cvm

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCvmInstancesModification() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCvmInstancesModificationRead,

		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "One or more instance ID to be queried. It can be obtained from the 实例 ID in the returned 值 of API DescribeInstances. The 最大instances in batch for each request is 20。",
			},
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "The upper 限制 of Filters for each request is 10 and the upper 限制 for Filter.Values is 2。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Fields to be filtered。",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "值 of the field。",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			"instance_type_config_status_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The 列表 model configurations that can be adjusted by the instance。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State 描述",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态 描述 information。",
						},
						"instance_type_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Configuration information。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Availability 可用区",
									},
									"instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例类型",
									},
									"instance_family": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance family。",
									},
									"gpu": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The 数量 GPU kernels，in cores。",
									},
									"cpu": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The 数量 CPU kernels，in cores。",
									},
									"memory": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Memory capacity (in GB)。",
									},
									"fpga": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The 数量 FPGA kernels，in cores。",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCvmInstancesModificationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cvm_instances_modification.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = cvm.NewDescribeInstancesModificationRequest()
		response = cvm.NewDescribeInstancesModificationResponse()
	)
	if v, ok := d.GetOk("instance_ids"); ok {
		request.InstanceIds = helper.InterfacesStringsPoint(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("filters"); ok {
		filters := make([]*cvm.Filter, 0)
		for _, item := range v.(*schema.Set).List() {
			filter := item.(map[string]interface{})
			name := filter["name"].(string)
			filters = append(filters, &cvm.Filter{
				Name:   &name,
				Values: helper.StringsStringsPoint(filter["values"].([]string)),
			})
		}
		request.Filters = filters
	}

	instanceTypeConfigStatusList := make([]map[string]interface{}, 0)

	var innerErr error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, innerErr = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().DescribeInstancesModification(request)
		if innerErr != nil {
			return tccommon.RetryError(innerErr)
		}
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	for _, instanceTypeConfigStatusSetItem := range response.Response.InstanceTypeConfigStatusSet {
		instanceTypeConfigStatus := make(map[string]interface{})
		instanceTypeConfigStatus["status"] = instanceTypeConfigStatusSetItem.Status
		instanceTypeConfigStatus["message"] = instanceTypeConfigStatusSetItem.Message

		instanceTypeConfigMaps := make([]map[string]interface{}, 0)
		instanceTypeConfigMap := make(map[string]interface{})
		instanceTypeConfig := instanceTypeConfigStatusSetItem.InstanceTypeConfig
		instanceTypeConfigMap["zone"] = instanceTypeConfig.Zone
		ids = append(ids, *instanceTypeConfig.InstanceType)
		instanceTypeConfigMap["instance_type"] = instanceTypeConfig.InstanceType
		instanceTypeConfigMap["instance_family"] = instanceTypeConfig.InstanceFamily
		instanceTypeConfigMap["gpu"] = instanceTypeConfig.GPU
		instanceTypeConfigMap["cpu"] = instanceTypeConfig.CPU
		instanceTypeConfigMap["memory"] = instanceTypeConfig.Memory
		instanceTypeConfigMap["fpga"] = instanceTypeConfig.FPGA
		instanceTypeConfigMaps = append(instanceTypeConfigMaps, instanceTypeConfigMap)
		instanceTypeConfigStatus["instance_type_config"] = instanceTypeConfigMaps

		instanceTypeConfigStatusList = append(instanceTypeConfigStatusList, instanceTypeConfigStatus)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("instance_type_config_status_list", instanceTypeConfigStatusList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set instance list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), instanceTypeConfigStatusList); err != nil {
			return err
		}
	}
	return nil

}
