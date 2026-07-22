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
				Description: "Launch configuration 名称",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "用于保存结果。",
			},
			"configuration_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A 列表 configuration. Each element 包含following attributes:",
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
							Description: "Launch configuration 名称",
						},
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID available image，for example `img-8toqc6s3`。",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID project to which the configuration belongs. 默认值为 0。",
						},
						"instance_types": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "实例类型 列表 the scaling configuration。",
						},
						"system_disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System disk category of the scaling configuration。",
						},
						"system_disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "System disk size of the scaling configuration （GB）。",
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
										Description: "类型 disk。",
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Volume of disk （GB）。 默认为 `0`。",
									},
									"snapshot_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data disk snapshot ID。",
									},
									"delete_with_instance": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "表示是否disk remove after instance terminated。",
									},
								},
							},
						},
						"internet_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Charge types for network traffic。",
						},
						"internet_max_bandwidth_out": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Max bandwidth of Internet access in Mbps。",
						},
						"public_ip_assigned": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "指定是否assign an 公网 IP 地址",
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
							Description: "Security groups to which the instance belongs。",
						},
						"enhanced_security_service": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否activate cloud security service。",
						},
						"enhanced_monitor_service": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "是否activate cloud monitor service。",
						},
						"user_data": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Base64-encoded 用户 Data text。",
						},
						"instance_tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "A 标签列表 associates with an instance。",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current 状态 a launch configuration。",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time when the launch configuration was created。",
						},
						"disk_type_policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy of cloud disk 类型",
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
