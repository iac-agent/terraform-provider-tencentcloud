package cdh

import (
	"context"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCdhInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdhInstancesRead,

		Schema: map[string]*schema.Schema{
			"host_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID CDH 实例 到 是 queried。",
			},
			"host_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "名称 CDH 实例 到 是 queried。",
			},
			"host_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "State 的 CDH 实例 到 是 queried. 有效值：`PENDING`，`LAUNCH_FAILURE`，`RUNNING`，`EXPIRED`。",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "可用 可用区 该 CDH 实例 locates 在。",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "项目 CDH belongs 到。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// computed
			"cdh_instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An 信息 列表 cdh 实例. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID CDH 实例。",
						},
						"host_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 CDH 实例。",
						},
						"host_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 CDH 实例。",
						},
						"host_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "State 的 CDH 实例。",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "charge 类型 CDH 实例。",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用 可用区 该 CDH 实例 locates 在。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "项目 CDH belongs 到。",
						},
						"prepaid_renew_flag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "自动续费标识",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 CDH 实例。",
						},
						"expired_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "过期时间 的 CDH 实例。",
						},
						"cvm_instance_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "ID 的 CVM 实例 该 have been 创建 在 CDH 实例。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cage_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cage ID CDH 实例. 此 参数 是 仅 有效 对于 CDH 实例 在 cages 的 finance availability zones。",
						},
						"host_resource": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "An 信息 列表 主机 资源. Each element 包含following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cpu_total_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 总数 CPU 核数 的 实例。",
									},
									"cpu_available_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "数量 可用 CPU 核数 的 实例。",
									},
									"memory_total_size": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "实例 内存 总数 容量，单位 （GB）。",
									},
									"memory_available_size": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "实例 内存 可用 容量，单位 （GB）。",
									},
									"disk_total_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "实例 磁盘 总数 容量，单位 （GB）。",
									},
									"disk_available_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "实例 磁盘 可用 容量，单位 （GB）。",
									},
									"disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "类型 磁盘。",
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

func dataSourceTencentCloudCdhInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cdh_instances.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cdhService := CdhService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	filter := make(map[string]string)
	if v, ok := d.GetOk("host_id"); ok {
		filter["host-id"] = v.(string)
	}
	if v, ok := d.GetOk("host_name"); ok {
		filter["host-name"] = v.(string)
	}
	if v, ok := d.GetOk("host_state"); ok {
		filter["host-state"] = v.(string)
	}
	if v, ok := d.GetOk("availability_zone"); ok {
		filter["zone"] = v.(string)
	}
	if v, ok := d.GetOk("project_id"); ok {
		filter["project-id"] = strconv.FormatInt(int64(v.(int)), 10)
	}

	var instances []*cvm.HostItem
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instances, errRet = cdhService.DescribeCdhInstanceByFilter(ctx, filter)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}

	instanceList := make([]map[string]interface{}, 0, len(instances))
	ids := make([]string, 0, len(instances))
	for _, instance := range instances {
		mapping := map[string]interface{}{
			"host_id":            instance.HostId,
			"host_name":          instance.HostName,
			"host_type":          instance.HostType,
			"host_state":         instance.HostState,
			"charge_type":        instance.HostChargeType,
			"availability_zone":  instance.Placement.Zone,
			"project_id":         instance.Placement.ProjectId,
			"prepaid_renew_flag": instance.RenewFlag,
			"create_time":        instance.CreatedTime,
			"expired_time":       instance.ExpiredTime,
			"cvm_instance_ids":   helper.StringsInterfaces(instance.InstanceIds),
			"cage_id":            instance.CageId,
		}
		hostResource := map[string]interface{}{
			"cpu_total_num":         instance.HostResource.CpuTotal,
			"cpu_available_num":     instance.HostResource.CpuAvailable,
			"memory_total_size":     instance.HostResource.MemTotal,
			"memory_available_size": instance.HostResource.MemAvailable,
			"disk_total_size":       instance.HostResource.DiskTotal,
			"disk_available_size":   instance.HostResource.DiskAvailable,
			"disk_type":             instance.HostResource.DiskType,
		}
		mapping["host_resource"] = []map[string]interface{}{hostResource}
		instanceList = append(instanceList, mapping)
		ids = append(ids, *instance.HostId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("cdh_instance_list", instanceList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set cdh instance list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), instanceList); err != nil {
			return err
		}
	}
	return nil
}
