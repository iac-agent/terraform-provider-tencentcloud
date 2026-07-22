package vpc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudVpcCvmInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcCvmInstancesRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "过滤器 condition. `RouteTableIds` 和 `Filters` 不能 是 指定 在 same 时间. vpc-ID - String - (过滤器 condition) VPC 实例 ID，such 作为 `vpc-f49l6u0z`;实例-类型 - String - (过滤器 condition) CVM 实例 ID;实例-名称 - String - (过滤器 condition) CVM 名称",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "attribute 名称 如果 more 比 一个 过滤器 exists， logical relation between these Filters 是 `AND`。",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Attribute 值 如果 多个 值 exist 在 一个 过滤器， logical relationship between these 值 是 `OR`. For `bool` 参数， 有效 值 include `TRUE` 和 `FALSE`。",
						},
					},
				},
			},

			"instance_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "列表 CVM 实例。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC 实例 ID",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "子网实例 ID",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CVM 实例 ID。",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CVM 名称",
						},
						"instance_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CVM 状态",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CPU 核数 在 实例 (在 core)。",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例's 内存 容量. 单位：GB。",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间。",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例类型",
						},
						"eni_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例 ENI 配额 (包括 primary ENIs)。",
						},
						"eni_ip_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Private IP quoata 对于 实例 ENIs (包括 primary ENIs)。",
						},
						"instance_eni_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 ENIs (包括 primary ENIs) bound 到 实例。",
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

func dataSourceTencentCloudVpcCvmInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_vpc_cvm_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*vpc.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := vpc.Filter{}
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

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var instanceSet []*vpc.CvmInstance

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcCvmInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instanceSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceSet))
	tmpList := make([]map[string]interface{}, 0, len(instanceSet))

	if instanceSet != nil {
		for _, cvmInstance := range instanceSet {
			cvmInstanceMap := map[string]interface{}{}

			if cvmInstance.VpcId != nil {
				cvmInstanceMap["vpc_id"] = cvmInstance.VpcId
			}

			if cvmInstance.SubnetId != nil {
				cvmInstanceMap["subnet_id"] = cvmInstance.SubnetId
			}

			if cvmInstance.InstanceId != nil {
				cvmInstanceMap["instance_id"] = cvmInstance.InstanceId
			}

			if cvmInstance.InstanceName != nil {
				cvmInstanceMap["instance_name"] = cvmInstance.InstanceName
			}

			if cvmInstance.InstanceState != nil {
				cvmInstanceMap["instance_state"] = cvmInstance.InstanceState
			}

			if cvmInstance.CPU != nil {
				cvmInstanceMap["cpu"] = cvmInstance.CPU
			}

			if cvmInstance.Memory != nil {
				cvmInstanceMap["memory"] = cvmInstance.Memory
			}

			if cvmInstance.CreatedTime != nil {
				cvmInstanceMap["created_time"] = cvmInstance.CreatedTime
			}

			if cvmInstance.InstanceType != nil {
				cvmInstanceMap["instance_type"] = cvmInstance.InstanceType
			}

			if cvmInstance.EniLimit != nil {
				cvmInstanceMap["eni_limit"] = cvmInstance.EniLimit
			}

			if cvmInstance.EniIpLimit != nil {
				cvmInstanceMap["eni_ip_limit"] = cvmInstance.EniIpLimit
			}

			if cvmInstance.InstanceEniCount != nil {
				cvmInstanceMap["instance_eni_count"] = cvmInstance.InstanceEniCount
			}

			ids = append(ids, *cvmInstance.InstanceId)
			tmpList = append(tmpList, cvmInstanceMap)
		}

		_ = d.Set("instance_set", tmpList)
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
