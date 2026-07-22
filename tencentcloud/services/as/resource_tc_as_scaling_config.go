package as

import (
	"context"
	"fmt"
	"log"

	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAsScalingConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsScalingConfigCreate,
		Read:   resourceTencentCloudAsScalingConfigRead,
		Update: resourceTencentCloudAsScalingConfigUpdate,
		Delete: resourceTencentCloudAsScalingConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"configuration_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "名称 launch 配置。",
			},
			"image_id": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"image_id", "image_family"},
				Description:  "An 可用 镜像 ID 对于 cvm 实例。",
			},
			"image_family": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"image_id", "image_family"},
				Description:  "Image Family 名称 Either Image ID 或 Image Family 名称 必须 是 提供，但 不 both。",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Specifys 到 其中 项目 配置 belongs。",
			},
			"instance_types": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    10,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specified types 的 CVM 实例。",
			},
			"system_disk_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      SYSTEM_DISK_TYPE_CLOUD_PREMIUM,
				ValidateFunc: tccommon.ValidateAllowedStringValue(SYSTEM_DISK_ALLOW_TYPE),
				Description:  "类型 CVM 磁盘. 有效值：`CLOUD_PREMIUM` 和 `CLOUD_SSD`. 默认为 `CLOUD_PREMIUM`. 有效 当 disk_type_policy 是 ORIGINAL。",
			},
			"system_disk_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     50,
				Description: "Volume 的 系统 磁盘 （GB）。 默认为 `50`。",
			},
			"data_disk": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    11,
				Description: "数据盘配置",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      SYSTEM_DISK_TYPE_CLOUD_PREMIUM,
							ValidateFunc: tccommon.ValidateAllowedStringValue(SYSTEM_DISK_ALLOW_TYPE),
							Description:  "Types 的 磁盘. 有效值：`CLOUD_PREMIUM` 和 `CLOUD_SSD`. 有效 当 disk_type_policy 是 ORIGINAL。",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Volume 的 磁盘 （GB）。 默认为 `0`。",
						},
						"snapshot_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Data 磁盘 快照 ID。",
						},
						"delete_with_instance": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "表示是否disk remove after 实例 terminated. 默认为 `false`。",
						},
					},
				},
			},
			// payment
			"instance_charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Charge 类型 实例. 有效 值 是 `PREPAID`，`POSTPAID_BY_HOUR`，`SPOTPAID`，`CDCPAID`. 默认为 `POSTPAID_BY_HOUR`. NOTE: `SPOTPAID` 实例 必须 集合 `spot_instance_type` 和 `spot_max_price` 在 same 时间。",
			},
			"instance_charge_type_prepaid_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue(svccvm.CVM_PREPAID_PERIOD),
				Description:  "tenancy (在 month) 的 prepaid 实例，NOTE: 它 仅 works 当 instance_charge_type 是 集合 到 `PREPAID`. 有效 值 是 `1`，`2`，`3`，`4`，`5`，`6`，`7`，`8`，`9`，`10`，`11`，`12`，`24`，`36`。",
			},
			"instance_charge_type_prepaid_renew_flag": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(svccvm.CVM_PREPAID_RENEW_FLAG),
				Description:  "自动续费标识 有效值：`NOTIFY_AND_AUTO_RENEW`: notify upon expiration 和 renew automatically，`NOTIFY_AND_MANUAL_RENEW`: notify upon expiration 但 do 不 renew automatically，`DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. 默认值：`NOTIFY_AND_MANUAL_RENEW`. 如果 此 参数 是 指定 作为 `NOTIFY_AND_AUTO_RENEW`， 实例 将 是 automatically renewed 在 monthly basis 如果 账号 balance 是 sufficient. NOTE: 它 仅 works 当 instance_charge_type 是 集合 到 `PREPAID`。",
			},
			"spot_instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"one-time"}),
				Description:  "类型 spot 实例，仅 support `一个-时间` now. 注意: 它 仅 works 当 instance_charge_type 是 集合 到 `SPOTPAID`。",
			},
			"spot_max_price": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringNumber,
				Description:  "Max 价格 的 spot 实例, 是 格式 的 decimal 字符串, 对于 示例 \"0.50\". 注意: 它 仅 works 当 instance_charge_type 是 集合 到 `SPOTPAID`.",
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      INTERNET_CHARGE_TYPE_TRAFFIC_POSTPAID_BY_HOUR,
				ValidateFunc: tccommon.ValidateAllowedStringValue(INTERNET_CHARGE_ALLOW_TYPE),
				Description:  "Charge types 对于 网络 流量. 有效值：`BANDWIDTH_PREPAID`，`TRAFFIC_POSTPAID_BY_HOUR` 和 `BANDWIDTH_PACKAGE`。",
			},
			"internet_max_bandwidth_out": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Max 带宽 的 Internet 访问 在 Mbps. 默认为 `0`。",
			},
			"public_ip_assigned": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "指定是否assign 公网 IP 地址",
			},
			"bandwidth_package_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Bandwidth 包 ID。",
			},
			"ipv4_address_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"WanIP", "HighQualityEIP", "AntiDDoSEIP"}),
				Description:  "AddressType. 默认值：WanIP. For beta users 的 dedicated IP. 值 可以 是: HighQualityEIP: Dedicated IP. 注意 该 dedicated IPs 是 仅 可用 在 partial regions. For beta users 的 Anti-DDoS IP， 值 可以 是: AntiDDoSEIP: Anti-DDoS EIP. 注意 该 Anti-DDoS IPs 是 仅 可用 在 partial regions。",
			},
			"anti_ddos_package_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Anti-DDoS 服务 包 ID. 此 为必填项 当 您 want 到 请求 AntiDDoS IP。",
			},
			"is_keep_eip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "是否delete bound EIP 当 实例 是 destroyed. Range 的 值: True: retain EIP; False: 不 retain EIP. 注意 该 当 IPv4AddressType 字段 指定EIP 类型， 默认值 behavior 是 不 到 retain EIP. WanIP 是 unaffected 通过 此 字段 和 将 always 是 删除 使用 实例. Changing 此 字段 配置 将 take effect immediately 对于 resources already bound 到 scaling 组。",
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ValidateFunc:  tccommon.ValidateAsConfigPassword,
				ConflictsWith: []string{"keep_image_login"},
				Description:   "密码 到 访问。",
			},
			"key_ids": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"keep_image_login"},
				Description:   "ID 列表 keys。",
			},
			"keep_image_login": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"password", "key_ids"},
				Description:   "指定是否keep original settings 的 CVM 镜像. And 它 可以't 是 使用 使用 密码 或 key_ids together。",
			},
			"security_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Security groups 到 其中 CVM 实例 belongs。",
			},
			"enhanced_security_service": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "To 指定是否enable 云 安全 服务. 默认为 `TRUE`。",
			},
			"enhanced_monitor_service": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "To 指定是否enable 云 监控 服务. 默认为 `TRUE`。",
			},
			"enhanced_automation_tools_service": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "To 指定是否enable 云 automation tools 服务。",
			},
			"user_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ase64-encoded 用户 Data text， 长度 限制 是 16KB。",
			},
			"instance_tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A 列表 标签 用于associate different resources。",
			},
			"disk_type_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      SCALING_DISK_TYPE_POLICY_ORIGINAL,
				ValidateFunc: tccommon.ValidateAllowedStringValue(SCALING_DISK_TYPE_ALLOW_POLICY),
				Description:  "Policy 的 云 磁盘 类型 有效值：`ORIGINAL` 和 `AUTOMATIC`. 默认为 `ORIGINAL`。",
			},
			"cam_role_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "被授权访问的 CAM 角色名称",
			},
			"host_name_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Related settings 的 云 服务器 hostname (HostName)。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "主机 名称 云 服务器; dots (.) 和 dashes (-) 不能 是 使用 作为 first 和 last 字符 的 HostName，和 不能 是 使用 consecutively; Windows 实例 是 不 支持; other types (Linux，etc.) 实例: character 长度 是 [2，40]，它 是 allowed 到 support 多个 dots，和 there 是 paragraph between dots，和 each paragraph 是 allowed 到 consist 的 letters (无 uppercase 和 lowercase restrictions)，numbers 和 dashes (-). Pure numbers 是 不 allowed。",
						},
						"host_name_style": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "style 的 主机 名称 云 服务器， 值 范围 includes `ORIGINAL` 和 `UNIQUE`， 默认为 `ORIGINAL`; `ORIGINAL`， AS directly passes HostName filled 在 input 参数 到 CVM，和 CVM 可能 append sequence 到 HostName 数量， HostName 的 实例 在 scaling 组 将 conflict; `UNIQUE`， HostName filled 在 作为 参数 是 equivalent 到 主机名 prefix，AS 和 CVM 将 expand 它，和 HostName 的 实例 在 scaling 组 可以 是 guaranteed 到 是 唯一。",
						},
					},
				},
			},
			"instance_name_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Settings 的 CVM 实例 names。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "CVM 实例名称",
						},
						"instance_name_style": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(INSTANCE_NAME_STYLE),
							Default:      INSTANCE_NAME_ORIGINAL,
							Description:  "类型 CVM 实例名称 有效值：`ORIGINAL` 和 `UNIQUE`. 默认为 `ORIGINAL`。",
						},
					},
				},
			},

			"disaster_recover_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Placement 组 ID Only 一个 是 allowed。",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"dedicated_cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Dedicated 集群 ID",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "标签 的 launch 配置。",
			},
			// Computed values
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current statues 的 launch 配置。",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "时间 当 launch 配置 是 创建。",
			},
		},
	}
}

func resourceTencentCloudAsScalingConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_config.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := as.NewCreateLaunchConfigurationRequest()

	v := d.Get("configuration_name")
	request.LaunchConfigurationName = helper.String(v.(string))

	if v, ok := d.GetOk("image_id"); ok {
		request.ImageId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("image_family"); ok {
		request.ImageFamily = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.IntUint64(v.(int))
	}

	v = d.Get("instance_types")
	instanceTypes := v.([]interface{})
	request.InstanceTypes = make([]*string, 0, len(instanceTypes))
	for i := range instanceTypes {
		instanceType := instanceTypes[i].(string)
		request.InstanceTypes = append(request.InstanceTypes, &instanceType)
	}

	request.SystemDisk = &as.SystemDisk{}
	if v, ok := d.GetOk("system_disk_type"); ok {
		request.SystemDisk.DiskType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("system_disk_size"); ok {
		request.SystemDisk.DiskSize = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("data_disk"); ok {
		dataDisks := v.([]interface{})
		request.DataDisks = make([]*as.DataDisk, 0, len(dataDisks))
		for _, d := range dataDisks {
			value := d.(map[string]interface{})
			diskType := value["disk_type"].(string)
			diskSize := uint64(value["disk_size"].(int))
			snapshotId := value["snapshot_id"].(string)
			deleteWithInstance := value["delete_with_instance"].(bool)
			dataDisk := as.DataDisk{
				DiskType:           &diskType,
				DiskSize:           &diskSize,
				DeleteWithInstance: &deleteWithInstance,
			}
			if snapshotId != "" {
				dataDisk.SnapshotId = &snapshotId
			}
			request.DataDisks = append(request.DataDisks, &dataDisk)
		}
	}

	request.InternetAccessible = &as.InternetAccessible{}
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request.InternetAccessible.InternetChargeType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
		request.InternetAccessible.InternetMaxBandwidthOut = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOkExists("public_ip_assigned"); ok {
		publicIpAssigned := v.(bool)
		request.InternetAccessible.PublicIpAssigned = &publicIpAssigned
	}
	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		request.InternetAccessible.BandwidthPackageId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("ipv4_address_type"); ok {
		request.InternetAccessible.IPv4AddressType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("anti_ddos_package_id"); ok {
		request.InternetAccessible.AntiDDoSPackageId = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("is_keep_eip"); ok {
		request.InternetAccessible.IsKeepEIP = helper.Bool(v.(bool))
	}

	request.LoginSettings = &as.LoginSettings{}
	if v, ok := d.GetOk("password"); ok {
		request.LoginSettings.Password = helper.String(v.(string))
	}
	if v, ok := d.GetOk("key_ids"); ok {
		keyIds := v.([]interface{})
		request.LoginSettings.KeyIds = make([]*string, 0, len(keyIds))
		for i := range keyIds {
			keyId := keyIds[i].(string)
			request.LoginSettings.KeyIds = append(request.LoginSettings.KeyIds, &keyId)
		}
	}
	if v, ok := d.GetOk("keep_image_login"); ok {
		keepImageLogin := v.(bool)
		request.LoginSettings.KeepImageLogin = &keepImageLogin
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroups := v.([]interface{})
		request.SecurityGroupIds = make([]*string, 0, len(securityGroups))
		for i := range securityGroups {
			securityGroup := securityGroups[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroup)
		}
	}

	request.EnhancedService = &as.EnhancedService{}

	if v, ok := d.GetOkExists("enhanced_security_service"); ok {
		securityService := v.(bool)
		request.EnhancedService.SecurityService = &as.RunSecurityServiceEnabled{
			Enabled: &securityService,
		}
	}
	if v, ok := d.GetOkExists("enhanced_monitor_service"); ok {
		monitorService := v.(bool)
		request.EnhancedService.MonitorService = &as.RunMonitorServiceEnabled{
			Enabled: &monitorService,
		}
	}
	if v, ok := d.GetOkExists("enhanced_automation_tools_service"); ok {
		automationToolsService := v.(bool)
		request.EnhancedService.AutomationToolsService = &as.RunAutomationServiceEnabled{
			Enabled: &automationToolsService,
		}
	}

	if v, ok := d.GetOk("user_data"); ok {
		request.UserData = helper.String(v.(string))
	}

	chargeType, ok := d.Get("instance_charge_type").(string)
	if !ok || chargeType == "" {
		chargeType = INSTANCE_CHARGE_TYPE_POSTPAID
	}

	if chargeType == INSTANCE_CHARGE_TYPE_SPOTPAID {
		spotMaxPrice := d.Get("spot_max_price").(string)
		spotInstanceType := d.Get("spot_instance_type").(string)
		request.InstanceMarketOptions = &as.InstanceMarketOptionsRequest{
			MarketType: helper.String("spot"),
			SpotOptions: &as.SpotMarketOptions{
				MaxPrice:         &spotMaxPrice,
				SpotInstanceType: &spotInstanceType,
			},
		}
	}

	if chargeType == INSTANCE_CHARGE_TYPE_PREPAID {
		period := d.Get("instance_charge_type_prepaid_period").(int)
		renewFlag := d.Get("instance_charge_type_prepaid_renew_flag").(string)
		request.InstanceChargePrepaid = &as.InstanceChargePrepaid{
			Period:    helper.IntInt64(period),
			RenewFlag: &renewFlag,
		}
	}

	request.InstanceChargeType = &chargeType

	if v, ok := d.GetOk("instance_types_check_policy"); ok {
		request.InstanceTypesCheckPolicy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_tags"); ok {
		tags := v.(map[string]interface{})
		request.InstanceTags = make([]*as.InstanceTag, 0, len(tags))
		for k, t := range tags {
			key := k
			value := t.(string)
			tag := as.InstanceTag{
				Key:   &key,
				Value: &value,
			}
			request.InstanceTags = append(request.InstanceTags, &tag)
		}
	}

	if v, ok := d.GetOk("disk_type_policy"); ok {
		request.DiskTypePolicy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cam_role_name"); ok {
		request.CamRoleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("host_name_settings"); ok {
		settings := make([]*as.HostNameSettings, 0, 10)
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			settingsInfo := as.HostNameSettings{}
			if hostName, ok := dMap["host_name"]; ok {
				settingsInfo.HostName = helper.String(hostName.(string))
			}
			if hostNameStyle, ok := dMap["host_name_style"]; ok {
				settingsInfo.HostNameStyle = helper.String(hostNameStyle.(string))
			}
			settings = append(settings, &settingsInfo)
		}
		request.HostNameSettings = settings[0]
	}

	if v, ok := d.GetOk("instance_name_settings"); ok {
		settings := make([]*as.InstanceNameSettings, 0, 10)
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			settingsInfo := as.InstanceNameSettings{}
			if instanceName, ok := dMap["instance_name"]; ok {
				settingsInfo.InstanceName = helper.String(instanceName.(string))
			}
			if instanceNameStyle, ok := dMap["instance_name_style"]; ok {
				settingsInfo.InstanceNameStyle = helper.String(instanceNameStyle.(string))
			}
			settings = append(settings, &settingsInfo)
		}
		request.InstanceNameSettings = settings[0]
	}

	if v, ok := d.GetOk("disaster_recover_group_ids"); ok {
		disasterRecoverGroupIds := v.([]interface{})
		request.DisasterRecoverGroupIds = make([]*string, 0, len(disasterRecoverGroupIds))
		for i := range disasterRecoverGroupIds {
			subnetId := disasterRecoverGroupIds[i].(string)
			request.DisasterRecoverGroupIds = append(request.DisasterRecoverGroupIds, &subnetId)
		}
	}

	if v, ok := d.GetOk("dedicated_cluster_id"); ok {
		request.DedicatedClusterId = helper.String(v.(string))
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		for tagKey, tagValue := range tags {
			tag := as.Tag{
				ResourceType: helper.String("launch-configuration"),
				Key:          helper.String(tagKey),
				Value:        helper.String(tagValue),
			}

			request.Tags = append(request.Tags, &tag)
		}
	}

	var launchConfigurationId string
	err := resource.Retry(4*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().CreateLaunchConfiguration(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		}

		if response.Response.LaunchConfigurationId == nil {
			return resource.NonRetryableError(fmt.Errorf("Launch configuration id is nil"))
		}
		launchConfigurationId = *response.Response.LaunchConfigurationId
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(launchConfigurationId)

	return resourceTencentCloudAsScalingConfigRead(d, meta)
}

func resourceTencentCloudAsScalingConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	configurationId := d.Id()
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		config, has, e := asService.DescribeLaunchConfigurationById(ctx, configurationId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if has == 0 {
			d.SetId("")
			return nil
		}
		_ = d.Set("configuration_name", *config.LaunchConfigurationName)
		_ = d.Set("status", *config.LaunchConfigurationStatus)

		if config.ImageId != nil {
			_ = d.Set("image_id", *config.ImageId)
		}
		if config.ImageFamily != nil {
			_ = d.Set("image_family", *config.ImageFamily)
		}

		_ = d.Set("project_id", *config.ProjectId)
		_ = d.Set("instance_types", helper.StringsInterfaces(config.InstanceTypes))
		_ = d.Set("system_disk_size", *config.SystemDisk.DiskSize)
		_ = d.Set("data_disk", flattenDataDiskMappings(config.DataDisks))
		_ = d.Set("internet_charge_type", *config.InternetAccessible.InternetChargeType)
		_ = d.Set("internet_max_bandwidth_out", *config.InternetAccessible.InternetMaxBandwidthOut)
		_ = d.Set("public_ip_assigned", *config.InternetAccessible.PublicIpAssigned)
		_ = d.Set("key_ids", helper.StringsInterfaces(config.LoginSettings.KeyIds))
		_ = d.Set("security_group_ids", helper.StringsInterfaces(config.SecurityGroupIds))
		_ = d.Set("enhanced_security_service", *config.EnhancedService.SecurityService.Enabled)
		_ = d.Set("enhanced_monitor_service", *config.EnhancedService.MonitorService.Enabled)
		if config.EnhancedService.AutomationToolsService.Enabled != nil {
			_ = d.Set("enhanced_automation_tools_service", *config.EnhancedService.AutomationToolsService.Enabled)
		}
		_ = d.Set("user_data", helper.PString(config.UserData))
		_ = d.Set("instance_tags", flattenInstanceTagsMapping(config.InstanceTags))
		_ = d.Set("disk_type_policy", *config.DiskTypePolicy)

		_ = d.Set("cam_role_name", *config.CamRoleName)

		if config.InternetAccessible != nil {
			if config.InternetAccessible.BandwidthPackageId != nil {
				_ = d.Set("bandwidth_package_id", config.InternetAccessible.BandwidthPackageId)
			}

			if config.InternetAccessible.IPv4AddressType != nil {
				_ = d.Set("ipv4_address_type", config.InternetAccessible.IPv4AddressType)
			}

			if config.InternetAccessible.AntiDDoSPackageId != nil {
				_ = d.Set("anti_ddos_package_id", config.InternetAccessible.AntiDDoSPackageId)
			}

			if config.InternetAccessible.IsKeepEIP != nil {
				_ = d.Set("is_keep_eip", config.InternetAccessible.IsKeepEIP)
			}
		}

		if config.HostNameSettings != nil {
			isEmptySettings := true
			settings := map[string]interface{}{}
			if config.HostNameSettings.HostName != nil {
				isEmptySettings = false
				settings["host_name"] = config.HostNameSettings.HostName
			}
			if config.HostNameSettings.HostNameStyle != nil {
				isEmptySettings = false
				settings["host_name_style"] = config.HostNameSettings.HostNameStyle
			}
			if !isEmptySettings {
				_ = d.Set("host_name_settings", []interface{}{settings})
			}
		}

		if config.InstanceNameSettings != nil {
			settings := make([]map[string]interface{}, 0)
			setting := map[string]interface{}{
				"instance_name":       config.InstanceNameSettings.InstanceName,
				"instance_name_style": config.InstanceNameSettings.InstanceNameStyle,
			}
			name, nameOk := setting["instance_name"].(string)
			style, styleOk := setting["instance_name_style"].(string)
			if nameOk && name != "" || styleOk && style != "" {
				settings = append(settings, setting)
				_ = d.Set("instance_name_settings", settings)
			}
		}

		if config.SystemDisk.DiskType != nil {
			_ = d.Set("system_disk_type", *config.SystemDisk.DiskType)
		}

		if _, ok := d.GetOk("instance_charge_type"); ok || *config.InstanceChargeType != INSTANCE_CHARGE_TYPE_POSTPAID {
			_ = d.Set("instance_charge_type", *config.InstanceChargeType)
		}

		if config.InstanceMarketOptions != nil && config.InstanceMarketOptions.SpotOptions != nil {
			_ = d.Set("spot_instance_type", config.InstanceMarketOptions.SpotOptions.SpotInstanceType)
			_ = d.Set("spot_max_price", config.InstanceMarketOptions.SpotOptions.MaxPrice)
		}

		if config.InstanceChargePrepaid != nil {
			_ = d.Set("instance_charge_type_prepaid_renew_flag", config.InstanceChargePrepaid.RenewFlag)
		}

		if len(config.DisasterRecoverGroupIds) > 0 {
			_ = d.Set("disaster_recover_group_ids", helper.StringsInterfaces(config.DisasterRecoverGroupIds))
		} else {
			_ = d.Set("disaster_recover_group_ids", []string{})
		}

		if config.DedicatedClusterId != nil {
			_ = d.Set("dedicated_cluster_id", config.DedicatedClusterId)
		}

		if config.Tags != nil && len(config.Tags) > 0 {
			_ = d.Set("tags", flattenTagsMapping(config.Tags))
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudAsScalingConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_config.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := as.NewModifyLaunchConfigurationAttributesRequest()

	configurationId := d.Id()
	request.LaunchConfigurationId = &configurationId

	if d.HasChange("configuration_name") {
		if v, ok := d.GetOk("configuration_name"); ok {
			request.LaunchConfigurationName = helper.String(v.(string))
		}
	}

	if d.HasChange("image_id") {
		if v, ok := d.GetOk("image_id"); ok {
			request.ImageId = helper.String(v.(string))
		}
	}

	if d.HasChange("project_id") {
		return fmt.Errorf("`project_id` do not support change now.")
	}

	if d.HasChange("instance_types") {
		if v, ok := d.GetOk("instance_types"); ok {
			instanceTypes := v.([]interface{})
			request.InstanceTypes = make([]*string, 0, len(instanceTypes))
			for i := range instanceTypes {
				instanceType := instanceTypes[i].(string)
				request.InstanceTypes = append(request.InstanceTypes, &instanceType)
			}
		}
	}

	if d.HasChange("system_disk_type") || d.HasChange("system_disk_size") {
		request.SystemDisk = &as.SystemDisk{}
		if v, ok := d.GetOk("system_disk_type"); ok {
			request.SystemDisk.DiskType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("system_disk_size"); ok {
			request.SystemDisk.DiskSize = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("data_disk") {
		if v, ok := d.GetOk("data_disk"); ok {
			dataDisks := v.([]interface{})
			request.DataDisks = make([]*as.DataDisk, 0, len(dataDisks))
			for _, d := range dataDisks {
				value := d.(map[string]interface{})
				diskType := value["disk_type"].(string)
				diskSize := uint64(value["disk_size"].(int))
				snapshotId := value["snapshot_id"].(string)
				deleteWithInstance := value["delete_with_instance"].(bool)
				dataDisk := as.DataDisk{
					DiskType:           &diskType,
					DiskSize:           &diskSize,
					DeleteWithInstance: &deleteWithInstance,
				}
				if snapshotId != "" {
					dataDisk.SnapshotId = &snapshotId
				}
				request.DataDisks = append(request.DataDisks, &dataDisk)
			}
		}
	}

	if d.HasChange("internet_charge_type") || d.HasChange("internet_max_bandwidth_out") || d.HasChange("public_ip_assigned") ||
		d.HasChange("bandwidth_package_id") || d.HasChange("ipv4_address_type") || d.HasChange("anti_ddos_package_id") || d.HasChange("is_keep_eip") {
		request.InternetAccessible = &as.InternetAccessible{}
		if v, ok := d.GetOk("internet_charge_type"); ok {
			request.InternetAccessible.InternetChargeType = helper.String(v.(string))
		}
		if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
			request.InternetAccessible.InternetMaxBandwidthOut = helper.IntUint64(v.(int))
		}
		if v, ok := d.GetOkExists("public_ip_assigned"); ok {
			publicIpAssigned := v.(bool)
			request.InternetAccessible.PublicIpAssigned = &publicIpAssigned
		}
		if v, ok := d.GetOk("bandwidth_package_id"); ok {
			request.InternetAccessible.BandwidthPackageId = helper.String(v.(string))
		}
		if v, ok := d.GetOk("ipv4_address_type"); ok {
			request.InternetAccessible.IPv4AddressType = helper.String(v.(string))
		}
		if v, ok := d.GetOk("anti_ddos_package_id"); ok {
			request.InternetAccessible.AntiDDoSPackageId = helper.String(v.(string))
		}
		if v, ok := d.GetOkExists("is_keep_eip"); ok {
			request.InternetAccessible.IsKeepEIP = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("security_group_ids") {
		if v, ok := d.GetOk("security_group_ids"); ok {
			securityGroups := v.([]interface{})
			request.SecurityGroupIds = make([]*string, 0, len(securityGroups))
			for i := range securityGroups {
				securityGroup := securityGroups[i].(string)
				request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroup)
			}
		}
	}

	if d.HasChange("enhanced_security_service") || d.HasChange("enhanced_monitor_service") || d.HasChange("enhanced_automation_tools_service") {
		request.EnhancedService = &as.EnhancedService{}

		if v, ok := d.GetOkExists("enhanced_security_service"); ok {
			securityService := v.(bool)
			request.EnhancedService.SecurityService = &as.RunSecurityServiceEnabled{
				Enabled: &securityService,
			}
		}
		if v, ok := d.GetOkExists("enhanced_monitor_service"); ok {
			monitorService := v.(bool)
			request.EnhancedService.MonitorService = &as.RunMonitorServiceEnabled{
				Enabled: &monitorService,
			}
		}
		if v, ok := d.GetOkExists("enhanced_automation_tools_service"); ok {
			automationToolsService := v.(bool)
			request.EnhancedService.AutomationToolsService = &as.RunAutomationServiceEnabled{
				Enabled: &automationToolsService,
			}
		}
	}

	if d.HasChange("user_data") {
		if v, ok := d.GetOk("user_data"); ok {
			request.UserData = helper.String(v.(string))
		}
	}

	if d.HasChange("instance_charge_type") {
		chargeType, ok := d.Get("instance_charge_type").(string)
		if !ok || chargeType == "" {
			chargeType = INSTANCE_CHARGE_TYPE_POSTPAID
		}

		if chargeType == INSTANCE_CHARGE_TYPE_SPOTPAID {
			spotMaxPrice := d.Get("spot_max_price").(string)
			spotInstanceType := d.Get("spot_instance_type").(string)
			request.InstanceMarketOptions = &as.InstanceMarketOptionsRequest{
				MarketType: helper.String("spot"),
				SpotOptions: &as.SpotMarketOptions{
					MaxPrice:         &spotMaxPrice,
					SpotInstanceType: &spotInstanceType,
				},
			}
		}

		if chargeType == INSTANCE_CHARGE_TYPE_PREPAID {
			period := d.Get("instance_charge_type_prepaid_period").(int)
			renewFlag := d.Get("instance_charge_type_prepaid_renew_flag").(string)
			request.InstanceChargePrepaid = &as.InstanceChargePrepaid{
				Period:    helper.IntInt64(period),
				RenewFlag: &renewFlag,
			}
		}

		request.InstanceChargeType = &chargeType
	}

	if d.HasChange("instance_types_check_policy") {
		if v, ok := d.GetOk("instance_types_check_policy"); ok {
			request.InstanceTypesCheckPolicy = helper.String(v.(string))
		}
	}

	if d.HasChange("instance_tags") {
		if v, ok := d.GetOk("instance_tags"); ok {
			tags := v.(map[string]interface{})
			request.InstanceTags = make([]*as.InstanceTag, 0, len(tags))
			for k, t := range tags {
				key := k
				value := t.(string)
				tag := as.InstanceTag{
					Key:   &key,
					Value: &value,
				}
				request.InstanceTags = append(request.InstanceTags, &tag)
			}
		}
	}

	if d.HasChange("disk_type_policy") {
		if v, ok := d.GetOk("disk_type_policy"); ok {
			request.DiskTypePolicy = helper.String(v.(string))
		}
	}

	if d.HasChange("cam_role_name") {
		if v, ok := d.GetOk("cam_role_name"); ok {
			request.CamRoleName = helper.String(v.(string))
		}
	}

	if d.HasChange("host_name_settings") {
		if v, ok := d.GetOk("host_name_settings"); ok {
			settings := make([]*as.HostNameSettings, 0, 10)
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				settingsInfo := as.HostNameSettings{}
				if hostName, ok := dMap["host_name"]; ok {
					settingsInfo.HostName = helper.String(hostName.(string))
				}
				if hostNameStyle, ok := dMap["host_name_style"]; ok {
					settingsInfo.HostNameStyle = helper.String(hostNameStyle.(string))
				}
				settings = append(settings, &settingsInfo)
			}
			request.HostNameSettings = settings[0]
		}
	}

	if d.HasChange("instance_name_settings") {
		if v, ok := d.GetOk("instance_name_settings"); ok {
			settings := make([]*as.InstanceNameSettings, 0, 10)
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				settingsInfo := as.InstanceNameSettings{}
				if instanceName, ok := dMap["instance_name"]; ok {
					settingsInfo.InstanceName = helper.String(instanceName.(string))
				}
				if instanceNameStyle, ok := dMap["instance_name_style"]; ok {
					settingsInfo.InstanceNameStyle = helper.String(instanceNameStyle.(string))
				}
				settings = append(settings, &settingsInfo)
			}
			request.InstanceNameSettings = settings[0]
		}
	}

	if d.HasChange("image_family") {
		if v, ok := d.GetOk("image_family"); ok {
			request.ImageFamily = helper.String(v.(string))
		}
	}

	if d.HasChange("password") || d.HasChange("key_ids") || d.HasChange("keep_image_login") {
		request.LoginSettings = &as.LoginSettings{}
		if v, ok := d.GetOk("password"); ok {
			request.LoginSettings.Password = helper.String(v.(string))
		}
		if v, ok := d.GetOk("key_ids"); ok {
			keyIds := v.([]interface{})
			request.LoginSettings.KeyIds = make([]*string, 0, len(keyIds))
			for i := range keyIds {
				keyId := keyIds[i].(string)
				request.LoginSettings.KeyIds = append(request.LoginSettings.KeyIds, &keyId)
			}
		}
		if v, ok := d.GetOk("keep_image_login"); ok {
			keepImageLogin := v.(bool)
			request.LoginSettings.KeepImageLogin = &keepImageLogin
		}
	}

	if d.HasChange("disaster_recover_group_ids") {
		if v, ok := d.GetOk("disaster_recover_group_ids"); ok {
			disasterRecoverGroupIds := v.([]interface{})
			request.DisasterRecoverGroupIds = make([]*string, 0, len(disasterRecoverGroupIds))
			for i := range disasterRecoverGroupIds {
				subnetId := disasterRecoverGroupIds[i].(string)
				request.DisasterRecoverGroupIds = append(request.DisasterRecoverGroupIds, &subnetId)
			}
		}
	}

	if d.HasChange("dedicated_cluster_id") {
		if v, ok := d.GetOk("dedicated_cluster_id"); ok {
			request.DedicatedClusterId = helper.String(v.(string))
		}
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(client)
		region := client.Region

		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))

		resourceName := tccommon.BuildTagResourceName("as", "launch-configuration", region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}

	err := resource.Retry(4*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().ModifyLaunchConfigurationAttributes(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), e.Error())
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		}

		return nil
	})
	if err != nil {
		return err
	}

	return resourceTencentCloudAsScalingConfigRead(d, meta)
}

func resourceTencentCloudAsScalingConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_config.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	configurationId := d.Id()
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := asService.DeleteLaunchConfiguration(ctx, configurationId)
	if err != nil {
		return err
	}

	return nil
}
