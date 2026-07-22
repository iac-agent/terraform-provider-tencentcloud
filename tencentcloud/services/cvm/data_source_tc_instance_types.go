package cvm

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	svccbs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cbs"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudInstanceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudInstanceTypesRead,

		Schema: map[string]*schema.Schema{
			"cpu_core_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "数量 CPU 核数 的 实例。",
			},
			"gpu_core_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "数量 GPU cores 的 实例。",
			},
			"memory_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "实例 内存 容量，单位 （GB）。",
			},
			"availability_zone": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"filter"},
				Description:   "可用 可用区 该 CVM 实例 locates 在. 此 字段 是 conflict 使用 `过滤器`。",
			},
			"filter": {
				Type:          schema.TypeSet,
				Optional:      true,
				MaxItems:      10,
				ConflictsWith: []string{"availability_zone"},
				Description:   "One 或 more 名称/值 pairs 到 过滤器. 此 字段 是 conflict 使用 `availability_zone`。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "过滤名称 有效值：`可用区`，`实例-family`，`实例-类型`，`实例-charge-类型` 和 `sort-keys`。",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "过滤器 值。",
						},
					},
				},
			},
			"cbs_filter": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Cbs 过滤器。",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_types": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Description: "硬盘介质类型。值范围:\n" +
								"	- CLOUD_BASIC: Represents ordinary Cloud Block Storage;\n" +
								"	- CLOUD_PREMIUM: Represents high-performance Cloud Block Storage;\n" +
								"	- CLOUD_SSD: Represents SSD Cloud Block Storage;\n" +
								"	- CLOUD_HSSD: Represents enhanced SSD Cloud Block Storage.",
						},
						"disk_charge_type": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "支付模式。值范围:\n" +
								"	- PREPAID: Prepaid;\n" +
								"	- POSTPAID_BY_HOUR: Post-payment.",
						},
						"disk_usage": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "系统盘或数据盘。值范围:\n" +
								"	- SYSTEM_DISK: Represents the system disk;\n" +
								"	- DATA_DISK: Represents the data disk.",
						},
					},
				},
			},
			"exclude_sold_out": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate 到 过滤器 实例 types 该 是 sold out 或 不，默认为 false。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// Computed values.
			"instance_types": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An 信息 列表 cvm 实例. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用 可用区 该 CVM 实例 locates 在。",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 实例。",
						},
						"cpu_core_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 CPU 核数 的 实例。",
						},
						"gpu_core_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 GPU cores 的 实例。",
						},
						"memory_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例 内存 容量，单位 （GB）。",
						},
						"family": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 series 的 实例。",
						},
						"instance_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Charge 类型 实例。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sell 状态 实例。",
						},
						"network_card": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Network card 类型，对于 示例: 25 表示 25G 网络 card。",
						},
						"type_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例类型 display 名称",
						},
						"sold_out_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reason 对于 sold out 状态",
						},
						"instance_bandwidth": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Internal 网络 带宽，单位: Gbps。",
						},
						"instance_pps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Network packet forwarding 容量，单位: 10K PPS。",
						},
						"storage_block_amount": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 本地 存储 blocks。",
						},
						"cpu_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Processor model。",
						},
						"fpga": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 FPGA cores。",
						},
						"gpu_count": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Physical GPU card count mapped 到 实例. vGPU 类型 是 less 比 1，direct-attach GPU 类型 是 greater 比 或 equal 到 1。",
						},
						"frequency": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CPU 频率 信息。",
						},
						"status_category": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stock 状态 category. 有效值：EnoughStock，NormalStock，UnderStock，WithoutStock。",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 备注 信息。",
						},
						"cbs_configs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "CBS 配置 cbs_configs 是 populated 当 cbs_filter 是 added。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"available": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否configuration 是 可用。",
									},
									"disk_charge_type": {
										Type:     schema.TypeString,
										Computed: true,
										Description: "支付模式。值范围:\n" +
											"	- PREPAID: Prepaid;\n" +
											"	- POSTPAID_BY_HOUR: Post-payment.",
									},
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "availability 可用区 到 其中 Cloud Block Storage belongs。",
									},
									"instance_family": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "实例 family。",
									},
									"disk_type": {
										Type:     schema.TypeString,
										Computed: true,
										Description: "硬盘介质类型。值范围:\n" +
											"	- CLOUD_BASIC: Represents ordinary Cloud Block Storage;\n" +
											"	- CLOUD_PREMIUM: Represents high-performance Cloud Block Storage;\n" +
											"	- CLOUD_SSD: Represents SSD Cloud Block Storage;\n" +
											"	- CLOUD_HSSD: Represents enhanced SSD Cloud Block Storage.",
									},
									"step_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum step 大小 change 在 云 磁盘 大小，（GB）。",
									},
									"extra_performance_range": {
										Computed:    true,
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Description: "Extra performance 范围。",
									},
									"device_class": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Device class。",
									},
									"disk_usage": {
										Type:     schema.TypeString,
										Computed: true,
										Description: "云盘类型。值范围:\n" +
											"	- SYSTEM_DISK: Represents the system disk;\n" +
											"	- DATA_DISK: Represents the data disk.",
									},
									"min_disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "最小 configurable 云 磁盘 大小，（GB）。",
									},
									"max_disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "最大 configurable 云 磁盘 大小，（GB）。",
									},
								},
							},
						},
						"local_disk_type_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "列表 本地 磁盘 specifications. Empty 如果 实例类型 does 不 support 本地 disks。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Local 磁盘 类型",
									},
									"partition_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Local 磁盘 分区 类型",
									},
									"min_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum 大小 的 本地 磁盘，（GB）。",
									},
									"max_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum 大小 的 本地 磁盘，（GB）。",
									},
									"required": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether 本地 磁盘 为必填项 当 purchasing. 有效值：REQUIRED，OPTIONAL。",
									},
								},
							},
						},
						"price": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "实例 pricing 信息。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"unit_price": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Subsequent 单位 价格，使用 在 postpaid 模式，单位: CNY。",
									},
									"charge_unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subsequent billing 单位. 有效值：HOUR，GB。",
									},
									"original_price": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Original 价格 对于 prepaid 模式，单位: CNY。",
									},
									"discount_price": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Discount 价格 对于 prepaid 模式，单位: CNY。",
									},
									"discount": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Discount 速率. For 示例，20.0 表示 20% 关闭。",
									},
									"unit_price_discount": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Subsequent discount 单位 价格，使用 在 postpaid 模式，单位: CNY。",
									},
									"unit_price_second_step": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Subsequent 单位 价格 对于 时间 范围 (96，360) hours 在 postpaid 模式，单位: CNY。",
									},
									"unit_price_discount_second_step": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Subsequent discount 单位 价格 对于 时间 范围 (96，360) hours 在 postpaid 模式，单位: CNY。",
									},
									"unit_price_third_step": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "指定original 价格 的 subsequent 总数 costs 使用 usage 时间间隔 exceeding 360 hr 在 postpaid billing 模式 measurement 单位: usd。",
									},
									"unit_price_discount_third_step": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Discounted 价格 的 subsequent 总数 费用 对于 usage 时间间隔 exceeding 360 hr 在 postpaid billing 模式 measurement 单位: usd。",
									},
								},
							},
						},
						"externals": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Extended attributes。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"release_address": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "是否release 地址",
									},
									"unsupport_networks": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Unsupported 网络 types. 有效值：BASIC (basic 网络)，VPC1.0 (VPC 1.0)。",
									},
									"storage_block_attr": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "HDD 本地 存储 attributes。",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Storage block 类型",
												},
												"min_size": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Minimum 大小 的 存储 block，（GB）。",
												},
												"max_size": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Maximum 大小 的 存储 block，（GB）。",
												},
											},
										},
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

func dataSourceTencentCloudInstanceTypesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_instance_types.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cvmService := CvmService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	isExcludeSoldOut := d.Get("exclude_sold_out").(bool)
	cpu, cpuOk := d.GetOk("cpu_core_count")
	gpu, gpuOk := d.GetOk("gpu_core_count")
	memory, memoryOk := d.GetOk("memory_size")
	var instanceSellTypes []*cvm.InstanceTypeQuotaItem
	var errRet error
	var err error
	typeList := make([]map[string]interface{}, 0)
	ids := make([]string, 0)

	var zone string
	var zone_in = 0
	if v, ok := d.GetOk("availability_zone"); ok {
		zone = v.(string)
		zone_in = 1
	}
	filters := d.Get("filter").(*schema.Set).List()
	filterMap := make(map[string][]string, len(filters)+zone_in)
	for _, v := range filters {
		item := v.(map[string]interface{})
		name := item["name"].(string)
		values := item["values"].([]interface{})
		filterValues := make([]string, 0, len(values))
		for _, value := range values {
			filterValues = append(filterValues, value.(string))
		}
		filterMap[name] = filterValues
	}
	if zone != "" {
		filterMap["zone"] = []string{zone}
	}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instanceSellTypes, errRet = cvmService.DescribeInstancesSellTypeByFilter(ctx, filterMap)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, instanceType := range instanceSellTypes {
		flag := true
		if cpuOk && int64(cpu.(int)) != *instanceType.Cpu {
			flag = false
		}
		if gpuOk && int64(gpu.(int)) != *instanceType.Gpu {
			flag = false
		}
		if memoryOk && int64(memory.(int)) != *instanceType.Memory {
			flag = false
		}
		if isExcludeSoldOut && CVM_SOLD_OUT_STATUS == *instanceType.Status {
			flag = false
		}

		if flag {
			mapping := map[string]interface{}{
				"availability_zone":    instanceType.Zone,
				"cpu_core_count":       instanceType.Cpu,
				"gpu_core_count":       instanceType.Gpu,
				"memory_size":          instanceType.Memory,
				"family":               instanceType.InstanceFamily,
				"instance_type":        instanceType.InstanceType,
				"instance_charge_type": instanceType.InstanceChargeType,
				"status":               instanceType.Status,
				"network_card":         instanceType.NetworkCard,
				"type_name":            instanceType.TypeName,
				"sold_out_reason":      instanceType.SoldOutReason,
				"instance_bandwidth":   instanceType.InstanceBandwidth,
				"instance_pps":         instanceType.InstancePps,
				"storage_block_amount": instanceType.StorageBlockAmount,
				"cpu_type":             instanceType.CpuType,
				"fpga":                 instanceType.Fpga,
				"gpu_count":            instanceType.GpuCount,
				"frequency":            instanceType.Frequency,
				"status_category":      instanceType.StatusCategory,
				"remark":               instanceType.Remark,
			}

			// Map local_disk_type_list
			localDiskList := make([]interface{}, 0)
			for _, localDisk := range instanceType.LocalDiskTypeList {
				localDiskList = append(localDiskList, map[string]interface{}{
					"type":           localDisk.Type,
					"partition_type": localDisk.PartitionType,
					"min_size":       localDisk.MinSize,
					"max_size":       localDisk.MaxSize,
					"required":       localDisk.Required,
				})
			}
			mapping["local_disk_type_list"] = localDiskList

			// Map price
			priceList := make([]interface{}, 0)
			if instanceType.Price != nil {
				priceList = append(priceList, map[string]interface{}{
					"unit_price":                      instanceType.Price.UnitPrice,
					"charge_unit":                     instanceType.Price.ChargeUnit,
					"original_price":                  instanceType.Price.OriginalPrice,
					"discount_price":                  instanceType.Price.DiscountPrice,
					"discount":                        instanceType.Price.Discount,
					"unit_price_discount":             instanceType.Price.UnitPriceDiscount,
					"unit_price_second_step":          instanceType.Price.UnitPriceSecondStep,
					"unit_price_discount_second_step": instanceType.Price.UnitPriceDiscountSecondStep,
					"unit_price_third_step":           instanceType.Price.UnitPriceThirdStep,
					"unit_price_discount_third_step":  instanceType.Price.UnitPriceDiscountThirdStep,
				})
			}
			mapping["price"] = priceList

			// Map externals
			externalsList := make([]interface{}, 0)
			if instanceType.Externals != nil {
				externalsMap := map[string]interface{}{
					"release_address":    instanceType.Externals.ReleaseAddress,
					"unsupport_networks": instanceType.Externals.UnsupportNetworks,
				}

				// Map storage_block_attr
				storageBlockList := make([]interface{}, 0)
				if instanceType.Externals.StorageBlockAttr != nil {
					storageBlockList = append(storageBlockList, map[string]interface{}{
						"type":     instanceType.Externals.StorageBlockAttr.Type,
						"min_size": instanceType.Externals.StorageBlockAttr.MinSize,
						"max_size": instanceType.Externals.StorageBlockAttr.MaxSize,
					})
				}
				externalsMap["storage_block_attr"] = storageBlockList
				externalsList = append(externalsList, externalsMap)
			}
			mapping["externals"] = externalsList

			typeList = append(typeList, mapping)
			ids = append(ids, *instanceType.InstanceType)
		}
	}

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	cbsService := svccbs.NewCbsService(client)
	cbsFilterParams := make(map[string]interface{})
	var hasCbsFilter bool
	if dMap, ok := helper.InterfacesHeadMap(d, "cbs_filter"); ok {
		if v, ok := dMap["disk_types"].([]interface{}); ok && len(v) > 0 {
			cbsFilterParams["disk_types"] = helper.InterfacesStrings(v)
		}
		if v, ok := dMap["disk_charge_type"].(string); ok && v != "" {
			cbsFilterParams["disk_charge_type"] = v
		}
		if v, ok := dMap["disk_usage"].(string); ok && v != "" {
			cbsFilterParams["disk_usage"] = v
		}
		hasCbsFilter = true
	}
	if hasCbsFilter {
		for idx, t := range typeList {
			filterParams := make(map[string]interface{})
			for k, v := range cbsFilterParams {
				filterParams[k] = v
			}

			if v, ok := t["availability_zone"].(*string); ok && v != nil {
				filterParams["availability_zone"] = *v
			}
			if v, ok := t["cpu_core_count"].(*int64); ok && v != nil {
				filterParams["cpu_core_count"] = *v
			}
			if v, ok := t["memory_size"].(*int64); ok && v != nil {
				filterParams["memory_size"] = *v
			}
			if v, ok := t["family"].(*string); ok && v != nil {
				filterParams["family"] = *v
			}
			diskConfigSet, err := cbsService.DescribeDiskConfigQuota(ctx, filterParams)
			if err != nil {
				return err
			}
			cbsConfigList := make([]interface{}, 0)
			for _, diskConfig := range diskConfigSet {
				cbsConfigList = append(cbsConfigList, map[string]interface{}{
					"available":               diskConfig.Available,
					"disk_charge_type":        diskConfig.DiskChargeType,
					"zone":                    diskConfig.Zone,
					"instance_family":         diskConfig.InstanceFamily,
					"disk_type":               diskConfig.DiskType,
					"step_size":               diskConfig.StepSize,
					"extra_performance_range": diskConfig.ExtraPerformanceRange,
					"device_class":            diskConfig.DeviceClass,
					"disk_usage":              diskConfig.DiskUsage,
					"min_disk_size":           diskConfig.MinDiskSize,
					"max_disk_size":           diskConfig.MaxDiskSize,
				})
			}
			typeList[idx]["cbs_configs"] = cbsConfigList
		}
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("instance_types", typeList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set instance type list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), typeList); err != nil {
			return err
		}
	}
	return nil
}
