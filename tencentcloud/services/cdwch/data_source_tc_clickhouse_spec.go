package cdwch

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdwch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClickhouseSpec() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClickhouseSpecRead,
		Schema: map[string]*schema.Schema{
			"zone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Regional 信息。",
			},

			"pay_mode": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Billing 类型，PREPAID 表示 annual 和 monthly subscription，POSTPAID_BY_HOUR 表示 pay-作为-您-go billing。",
			},

			"is_elastic": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Is 它 elastic。",
			},

			"common_spec": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Zookeeper 节点 规格 描述",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "规格名称",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 CPU 核数",
						},
						"mem": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory 大小，单位 G。",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Classification 标签，STANDARD/BIGDATA/HIGHIO respectively represent standard/big 数据/high IO。",
						},
						"system_disk": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "System 磁盘 描述 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Disk 类型",
									},
									"disk_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Disk 类型 描述",
									},
									"min_disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum 磁盘 大小，单位 G。",
									},
									"max_disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum 磁盘 大小，单位 G。",
									},
									"disk_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 disks。",
									},
								},
							},
						},
						"data_disk": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Data 磁盘 描述 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Disk 类型",
									},
									"disk_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Disk 类型 描述",
									},
									"min_disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum 磁盘 大小，单位 G。",
									},
									"max_disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum 磁盘 大小，单位 G。",
									},
									"disk_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 disks。",
									},
								},
							},
						},
						"max_node_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最大nodes 限制",
						},
						"available": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否为available，false 表示 sold out。",
						},
						"compute_spec_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specification 描述 信息。",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "规格名称",
						},
						"instance_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Inventory。",
						},
					},
				},
			},

			"data_spec": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Data 节点 规格 描述",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "规格名称",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 CPU 核数",
						},
						"mem": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory 大小，单位 G。",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Classification 标签，STANDARD/BIGDATA/HIGHIO respectively represent standard/big 数据/high IO。",
						},
						"system_disk": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "System 磁盘 描述 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Disk 类型",
									},
									"disk_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Disk 类型 描述",
									},
									"min_disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum 磁盘 大小，单位 G。",
									},
									"max_disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum 磁盘 大小，单位 G。",
									},
									"disk_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 disks。",
									},
								},
							},
						},
						"data_disk": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Data 磁盘 描述 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Disk 类型",
									},
									"disk_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Disk 类型 描述",
									},
									"min_disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum 磁盘 大小，单位 G。",
									},
									"max_disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum 磁盘 大小，单位 G。",
									},
									"disk_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 disks。",
									},
								},
							},
						},
						"max_node_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "最大nodes 限制",
						},
						"available": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否为available，false 表示 sold out。",
						},
						"compute_spec_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specification 描述 信息。",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "规格名称",
						},
						"instance_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Inventory。",
						},
					},
				},
			},

			"attach_cbs_spec": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Cloud 磁盘 列表。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk 类型",
						},
						"disk_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk 类型 描述",
						},
						"min_disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum 磁盘 大小，单位 G。",
						},
						"max_disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum 磁盘 大小，单位 G。",
						},
						"disk_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 disks。",
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

func dataSourceTencentCloudClickhouseSpecRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clickhouse_spec.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone"); ok {
		paramMap["Zone"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("pay_mode"); ok {
		paramMap["PayMode"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("is_elastic"); ok {
		paramMap["IsElastic"] = helper.Bool(v.(bool))
	}

	service := CdwchService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var spec *cdwch.DescribeSpecResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClickhouseSpecByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		spec = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0)

	if spec.CommonSpec != nil {
		commonSpec := spec.CommonSpec
		commonSpecTmpList := make([]map[string]interface{}, 0, len(commonSpec))

		for _, resourceSpec := range commonSpec {
			resourceSpecMap := map[string]interface{}{}

			if resourceSpec.Name != nil {
				resourceSpecMap["name"] = resourceSpec.Name
			}

			if resourceSpec.Cpu != nil {
				resourceSpecMap["cpu"] = resourceSpec.Cpu
			}

			if resourceSpec.Mem != nil {
				resourceSpecMap["mem"] = resourceSpec.Mem
			}

			if resourceSpec.Type != nil {
				resourceSpecMap["type"] = resourceSpec.Type
			}

			if resourceSpec.SystemDisk != nil {
				systemDiskMap := map[string]interface{}{}

				if resourceSpec.SystemDisk.DiskType != nil {
					systemDiskMap["disk_type"] = resourceSpec.SystemDisk.DiskType
				}

				if resourceSpec.SystemDisk.DiskDesc != nil {
					systemDiskMap["disk_desc"] = resourceSpec.SystemDisk.DiskDesc
				}

				if resourceSpec.SystemDisk.MinDiskSize != nil {
					systemDiskMap["min_disk_size"] = resourceSpec.SystemDisk.MinDiskSize
				}

				if resourceSpec.SystemDisk.MaxDiskSize != nil {
					systemDiskMap["max_disk_size"] = resourceSpec.SystemDisk.MaxDiskSize
				}

				if resourceSpec.SystemDisk.DiskCount != nil {
					systemDiskMap["disk_count"] = resourceSpec.SystemDisk.DiskCount
				}

				resourceSpecMap["system_disk"] = []interface{}{systemDiskMap}
			}

			if resourceSpec.DataDisk != nil {
				dataDiskMap := map[string]interface{}{}

				if resourceSpec.DataDisk.DiskType != nil {
					dataDiskMap["disk_type"] = resourceSpec.DataDisk.DiskType
				}

				if resourceSpec.DataDisk.DiskDesc != nil {
					dataDiskMap["disk_desc"] = resourceSpec.DataDisk.DiskDesc
				}

				if resourceSpec.DataDisk.MinDiskSize != nil {
					dataDiskMap["min_disk_size"] = resourceSpec.DataDisk.MinDiskSize
				}

				if resourceSpec.DataDisk.MaxDiskSize != nil {
					dataDiskMap["max_disk_size"] = resourceSpec.DataDisk.MaxDiskSize
				}

				if resourceSpec.DataDisk.DiskCount != nil {
					dataDiskMap["disk_count"] = resourceSpec.DataDisk.DiskCount
				}

				resourceSpecMap["data_disk"] = []interface{}{dataDiskMap}
			}

			if resourceSpec.MaxNodeSize != nil {
				resourceSpecMap["max_node_size"] = resourceSpec.MaxNodeSize
			}

			if resourceSpec.Available != nil {
				resourceSpecMap["available"] = resourceSpec.Available
			}

			if resourceSpec.ComputeSpecDesc != nil {
				resourceSpecMap["compute_spec_desc"] = resourceSpec.ComputeSpecDesc
			}

			if resourceSpec.DisplayName != nil {
				resourceSpecMap["display_name"] = resourceSpec.DisplayName
			}

			if resourceSpec.InstanceQuota != nil {
				resourceSpecMap["instance_quota"] = resourceSpec.InstanceQuota
			}

			commonSpecTmpList = append(commonSpecTmpList, resourceSpecMap)
		}
		tmpList = append(tmpList, commonSpecTmpList...)
		_ = d.Set("common_spec", commonSpecTmpList)
	}

	if spec.DataSpec != nil {
		dataSpec := spec.DataSpec
		dataSpecTmpList := make([]map[string]interface{}, 0, len(dataSpec))

		for _, resourceSpec := range dataSpec {
			resourceSpecMap := map[string]interface{}{}

			if resourceSpec.Name != nil {
				resourceSpecMap["name"] = resourceSpec.Name
			}

			if resourceSpec.Cpu != nil {
				resourceSpecMap["cpu"] = resourceSpec.Cpu
			}

			if resourceSpec.Mem != nil {
				resourceSpecMap["mem"] = resourceSpec.Mem
			}

			if resourceSpec.Type != nil {
				resourceSpecMap["type"] = resourceSpec.Type
			}

			if resourceSpec.SystemDisk != nil {
				systemDiskMap := map[string]interface{}{}

				if resourceSpec.SystemDisk.DiskType != nil {
					systemDiskMap["disk_type"] = resourceSpec.SystemDisk.DiskType
				}

				if resourceSpec.SystemDisk.DiskDesc != nil {
					systemDiskMap["disk_desc"] = resourceSpec.SystemDisk.DiskDesc
				}

				if resourceSpec.SystemDisk.MinDiskSize != nil {
					systemDiskMap["min_disk_size"] = resourceSpec.SystemDisk.MinDiskSize
				}

				if resourceSpec.SystemDisk.MaxDiskSize != nil {
					systemDiskMap["max_disk_size"] = resourceSpec.SystemDisk.MaxDiskSize
				}

				if resourceSpec.SystemDisk.DiskCount != nil {
					systemDiskMap["disk_count"] = resourceSpec.SystemDisk.DiskCount
				}

				resourceSpecMap["system_disk"] = []interface{}{systemDiskMap}
			}

			if resourceSpec.DataDisk != nil {
				dataDiskMap := map[string]interface{}{}

				if resourceSpec.DataDisk.DiskType != nil {
					dataDiskMap["disk_type"] = resourceSpec.DataDisk.DiskType
				}

				if resourceSpec.DataDisk.DiskDesc != nil {
					dataDiskMap["disk_desc"] = resourceSpec.DataDisk.DiskDesc
				}

				if resourceSpec.DataDisk.MinDiskSize != nil {
					dataDiskMap["min_disk_size"] = resourceSpec.DataDisk.MinDiskSize
				}

				if resourceSpec.DataDisk.MaxDiskSize != nil {
					dataDiskMap["max_disk_size"] = resourceSpec.DataDisk.MaxDiskSize
				}

				if resourceSpec.DataDisk.DiskCount != nil {
					dataDiskMap["disk_count"] = resourceSpec.DataDisk.DiskCount
				}

				resourceSpecMap["data_disk"] = []interface{}{dataDiskMap}
			}

			if resourceSpec.MaxNodeSize != nil {
				resourceSpecMap["max_node_size"] = resourceSpec.MaxNodeSize
			}

			if resourceSpec.Available != nil {
				resourceSpecMap["available"] = resourceSpec.Available
			}

			if resourceSpec.ComputeSpecDesc != nil {
				resourceSpecMap["compute_spec_desc"] = resourceSpec.ComputeSpecDesc
			}

			if resourceSpec.DisplayName != nil {
				resourceSpecMap["display_name"] = resourceSpec.DisplayName
			}

			if resourceSpec.InstanceQuota != nil {
				resourceSpecMap["instance_quota"] = resourceSpec.InstanceQuota
			}

			dataSpecTmpList = append(dataSpecTmpList, resourceSpecMap)
		}
		tmpList = append(tmpList, dataSpecTmpList...)
		_ = d.Set("data_spec", dataSpecTmpList)
	}

	if spec.AttachCBSSpec != nil {
		attachCBSSpec := spec.AttachCBSSpec
		attachCBSSpecTmpList := make([]map[string]interface{}, 0, len(attachCBSSpec))

		for _, diskSpec := range attachCBSSpec {
			diskSpecMap := map[string]interface{}{}

			if diskSpec.DiskType != nil {
				diskSpecMap["disk_type"] = diskSpec.DiskType
			}

			if diskSpec.DiskDesc != nil {
				diskSpecMap["disk_desc"] = diskSpec.DiskDesc
			}

			if diskSpec.MinDiskSize != nil {
				diskSpecMap["min_disk_size"] = diskSpec.MinDiskSize
			}

			if diskSpec.MaxDiskSize != nil {
				diskSpecMap["max_disk_size"] = diskSpec.MaxDiskSize
			}

			if diskSpec.DiskCount != nil {
				diskSpecMap["disk_count"] = diskSpec.DiskCount
			}

			attachCBSSpecTmpList = append(attachCBSSpecTmpList, diskSpecMap)
		}
		tmpList = append(tmpList, attachCBSSpecTmpList...)
		_ = d.Set("attach_cbs_spec", attachCBSSpecTmpList)
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
