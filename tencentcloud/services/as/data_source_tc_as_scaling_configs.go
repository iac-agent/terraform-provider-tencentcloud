package as

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudAsScalingConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAsScalingConfigRead,

		Schema: map[string]*schema.Schema{
			"configuration_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "启动配置 ID",
			},
			"configuration_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Launch 配置 名称",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"configuration_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 配置. Each element 包含following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"configuration_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "启动配置 ID",
						},
						"configuration_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Launch 配置 名称",
						},
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID 可用 镜像，对于 示例 `img-8toqc6s3`。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID 项目 到 其中 配置 belongs. 默认值为 0。",
						},
						"instance_types": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "实例类型 列表 scaling 配置。",
						},
						"system_disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System 磁盘 category 的 scaling 配置。",
						},
						"system_disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "System 磁盘 大小 的 scaling 配置 （GB）。",
						},
						"data_disk": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "数据盘配置",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "类型 磁盘。",
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Volume 的 磁盘 （GB）。 默认为 `0`。",
									},
									"snapshot_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data 磁盘 快照 ID。",
									},
									"delete_with_instance": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "表示是否disk remove after 实例 terminated。",
									},
								},
							},
						},
						"internet_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Charge types 对于 网络 流量。",
						},
						"internet_max_bandwidth_out": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Max 带宽 的 Internet 访问 在 Mbps。",
						},
						"public_ip_assigned": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "指定是否assign 公网 IP 地址",
						},
						"key_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "ID 列表 login keys。",
						},
						"security_group_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Security groups 到 其中 实例 belongs。",
						},
						"enhanced_security_service": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否activate 云 安全 服务。",
						},
						"enhanced_monitor_service": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否activate 云 监控 服务。",
						},
						"user_data": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Base64-encoded 用户 Data text。",
						},
						"instance_tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "A 标签列表 associates 使用 实例。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current 状态 launch 配置。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "时间 当 launch 配置 是 创建。",
						},
						"disk_type_policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy 的 云 磁盘 类型",
						},
						"version_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "版本 Number。",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAsScalingConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_as_scaling_configs.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	configurationId := ""
	configurationName := ""
	if v, ok := d.GetOk("configuration_id"); ok {
		configurationId = v.(string)
	}
	if v, ok := d.GetOk("configuration_name"); ok {
		configurationName = v.(string)
	}

	configs, err := asService.DescribeLaunchConfigurationByFilter(ctx, configurationId, configurationName)
	if err != nil {
		return err
	}

	configurationList := make([]map[string]interface{}, 0, len(configs))
	for _, config := range configs {
		mapping := map[string]interface{}{
			"configuration_id":           *config.LaunchConfigurationId,
			"configuration_name":         *config.LaunchConfigurationName,
			"image_id":                   *config.ImageId,
			"project_id":                 *config.ProjectId,
			"instance_types":             helper.StringsInterfaces(config.InstanceTypes),
			"system_disk_size":           *config.SystemDisk.DiskSize,
			"data_disk":                  flattenDataDiskMappings(config.DataDisks),
			"internet_charge_type":       *config.InternetAccessible.InternetChargeType,
			"internet_max_bandwidth_out": *config.InternetAccessible.InternetMaxBandwidthOut,
			"public_ip_assigned":         *config.InternetAccessible.PublicIpAssigned,
			"key_ids":                    helper.StringsInterfaces(config.LoginSettings.KeyIds),
			"security_group_ids":         helper.StringsInterfaces(config.SecurityGroupIds),
			"enhanced_security_service":  *config.EnhancedService.SecurityService.Enabled,
			"enhanced_monitor_service":   *config.EnhancedService.MonitorService.Enabled,
			"user_data":                  helper.PString(config.UserData),
			"instance_tags":              flattenInstanceTagsMapping(config.InstanceTags),
			"status":                     *config.LaunchConfigurationStatus,
			"create_time":                *config.CreatedTime,
			"disk_type_policy":           *config.DiskTypePolicy,
			"version_number":             *config.VersionNumber,
		}
		if config.SystemDisk.DiskType != nil {
			mapping["system_disk_type"] = *config.SystemDisk.DiskType
		}
		configurationList = append(configurationList, mapping)
	}

	d.SetId("ConfigurationList" + configurationId + configurationName)
	err = d.Set("configuration_list", configurationList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set configuration list fail, reason:%s\n ", logId, err.Error())
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err = tccommon.WriteToFile(output.(string), configurationList); err != nil {
			return err
		}
	}

	return nil
}
