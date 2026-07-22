package cvm

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudInstancesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID 实例 到 是 queried。",
			},
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 128),
				Description:  "名称 实例 到 是 queried。",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "可用 可用区 该 CVM 实例 locates 在。",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "项目 CVM belongs 到。",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID vpc 到 是 queried。",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID vpc subnetwork。",
			},
			"dedicated_cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Exclusive 集群 ID",
			},
			"instance_set_ids": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      100,
				ConflictsWith: []string{"instance_id", "instance_name", "availability_zone", "project_id", "vpc_id", "subnet_id", "tags"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "实例 集合 ids，max 长度 是 100，conflict 使用 other 字段。",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 的 实例。",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},

			// computed
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An 信息 列表 cvm 实例. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 实例。",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "名称 实例。",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 实例。",
						},
						"dedicated_cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Exclusive 集群 ID",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "数量 CPU 核数 的 实例。",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "实例 内存 容量，单位 （GB）。",
						},
						"os_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "实例 os 名称",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "可用 可用区 该 CVM 实例 locates 在。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "项目 CVM belongs 到。",
						},
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 镜像。",
						},
						"instance_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "charge 类型 实例。",
						},
						"system_disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "类型 系统 磁盘。",
						},
						"system_disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size 的 系统 磁盘。",
						},
						"system_disk_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image ID 系统 磁盘。",
						},
						"rack_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "rack ID 实例 资源 池 到 其中 实例 belongs。",
						},
						"data_disks": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "An 信息 列表 数据 磁盘. Each element 包含following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data_disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "类型 数据 磁盘。",
									},
									"data_disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Size 的 数据 磁盘。",
									},
									"data_disk_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Image ID 数据 磁盘。",
									},
									"delete_with_instance": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "表示是否data 磁盘 是 destroyed 使用 实例。",
									},
								},
							},
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID vpc。",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID vpc subnetwork。",
						},
						"internet_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "charge 类型 实例。",
						},
						"internet_max_bandwidth_out": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Public 网络 最大 output 带宽 的 实例。",
						},
						"allocate_public_ip": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "表示是否public ip 是 assigned。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "状态 实例。",
						},
						"public_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public IP 的 实例。",
						},
						"private_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Private IP 的 实例。",
						},
						"security_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Security groups 的 实例。",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "标签 的 实例。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "创建时间 的 实例。",
						},
						"expired_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "过期时间 的 实例。",
						},
						"instance_charge_type_prepaid_renew_flag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "way 该 CVM 实例 将 是 renew automatically 或 不 当 它 reach end 的 prepaid tenancy。",
						},
						"cam_role_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "被授权访问的 CAM 角色名称",
						},
						"uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Globally 唯一 ID 实例。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_instances.read")()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		cvmService     = CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceSetIds []*string
	)

	filter := make(map[string]string)
	if v, ok := d.GetOk("instance_id"); ok {
		filter["instance-id"] = v.(string)
	}

	if v, ok := d.GetOk("instance_name"); ok {
		filter["instance-name"] = v.(string)
	}

	if v, ok := d.GetOk("availability_zone"); ok {
		filter["zone"] = v.(string)
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		filter["project-id"] = fmt.Sprintf("%d", v.(int))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		filter["vpc-id"] = v.(string)
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		filter["subnet-id"] = v.(string)
	}

	if v, ok := d.GetOk("dedicated_cluster_id"); ok {
		filter["dedicated-cluster-id"] = v.(string)
	}

	if v, ok := d.GetOk("instance_set_ids"); ok {
		instanceSetIds = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOk("tags"); ok {
		for key, value := range v.(map[string]interface{}) {
			filter["tag:"+key] = value.(string)
		}
	}

	var instances []*cvm.Instance
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instances, errRet = cvmService.DescribeInstanceByFilter(ctx, instanceSetIds, filter)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
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
			"instance_id":                instance.InstanceId,
			"instance_name":              instance.InstanceName,
			"instance_type":              instance.InstanceType,
			"dedicated_cluster_id":       instance.DedicatedClusterId,
			"cpu":                        instance.CPU,
			"memory":                     instance.Memory,
			"os_name":                    instance.OsName,
			"availability_zone":          instance.Placement.Zone,
			"project_id":                 instance.Placement.ProjectId,
			"rack_id":                    instance.Placement.RackId,
			"image_id":                   instance.ImageId,
			"instance_charge_type":       instance.InstanceChargeType,
			"system_disk_type":           instance.SystemDisk.DiskType,
			"system_disk_size":           instance.SystemDisk.DiskSize,
			"system_disk_id":             instance.SystemDisk.DiskId,
			"vpc_id":                     instance.VirtualPrivateCloud.VpcId,
			"subnet_id":                  instance.VirtualPrivateCloud.SubnetId,
			"internet_charge_type":       instance.InternetAccessible.InternetChargeType,
			"internet_max_bandwidth_out": instance.InternetAccessible.InternetMaxBandwidthOut,
			"allocate_public_ip":         instance.InternetAccessible.PublicIpAssigned,
			"status":                     instance.InstanceState,
			"security_groups":            helper.StringsInterfaces(instance.SecurityGroupIds),
			"tags":                       flattenCvmTagsMapping(instance.Tags),
			"create_time":                instance.CreatedTime,
			"expired_time":               instance.ExpiredTime,
			"instance_charge_type_prepaid_renew_flag": instance.RenewFlag,
			"cam_role_name": instance.CamRoleName,
			"uuid":          instance.Uuid,
		}

		if len(instance.PublicIpAddresses) > 0 {
			mapping["public_ip"] = *instance.PublicIpAddresses[0]
		}

		if len(instance.PrivateIpAddresses) > 0 {
			mapping["private_ip"] = *instance.PrivateIpAddresses[0]
		}

		dataDisks := make([]map[string]interface{}, 0, len(instance.DataDisks))
		for _, v := range instance.DataDisks {
			dataDisk := map[string]interface{}{
				"data_disk_type":       v.DiskType,
				"data_disk_size":       v.DiskSize,
				"data_disk_id":         v.DiskId,
				"delete_with_instance": v.DeleteWithInstance,
			}

			dataDisks = append(dataDisks, dataDisk)
		}

		mapping["data_disks"] = dataDisks
		instanceList = append(instanceList, mapping)
		ids = append(ids, *instance.InstanceId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("instance_list", instanceList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set instance list fail, reason:%s\n ", logId, err.Error())
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
