/*
Provides an elastic kubernetes service container instance (offlined).

~> **NOTE:**  This resource was offline and no longer supported.

# Example Usage

```
data "tencentcloud_security_groups" "group" {
}

	data "tencentcloud_availability_zones_by_product" "zone" {
	  product = "cvm"
	}

	resource "tencentcloud_vpc" "vpc" {
	  cidr_block = "10.0.0.0/24"
	  name       = "tf-test-eksci"
	}

	resource "tencentcloud_subnet" "sub" {
	  availability_zone = data.tencentcloud_availability_zones_by_product.zone.zones[0].name
	  cidr_block        = "10.0.0.0/24"
	  name              = "sub"
	  vpc_id            = tencentcloud_vpc.vpc.id
	}

	resource "tencentcloud_cbs_storage" "cbs" {
	  availability_zone = data.tencentcloud_availability_zones_by_product.zone.zones[0].name
	  storage_name      = "cbs1"
	  storage_size      = 10
	  storage_type      = "CLOUD_PREMIUM"
	}

	resource "tencentcloud_eks_container_instance" "eci1" {
	  name = "foo"
	  vpc_id = tencentcloud_vpc.vpc.id
	  subnet_id = tencentcloud_subnet.sub.id
	  cpu = 2
	  cpu_type = "intel"
	  restart_policy = "Always"
	  memory = 4
	  security_groups = [data.tencentcloud_security_groups.group.security_groups[0].security_group_id]
	  cbs_volume {
	    name = "vol1"
	    disk_id = tencentcloud_cbs_storage.cbs.id
	  }
	  container {
	    name = "redis1"
	    image = "redis"
	    liveness_probe {
	      init_delay_seconds = 1
	      timeout_seconds = 3
	      period_seconds = 11
	      success_threshold = 1
	      failure_threshold = 3
	      http_get_path = "/"
	      http_get_port = 443
	      http_get_scheme = "HTTPS"
	    }
	    readiness_probe {
	      init_delay_seconds = 1
	      timeout_seconds = 3
	      period_seconds = 10
	      success_threshold = 1
	      failure_threshold = 3
	      tcp_socket_port = 81
	    }
	  }
	  container {
	    name = "nginx"
	    image = "nginx"
	  }
	  init_container {
	    name = "alpine"
	    image = "alpine:latest"
	  }
	}

```

# Import

```
terraform import tencentcloud_eks_container_instance.foo container-instance-id
```
*/
package tke

import (
	"context"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceEksCiProbeConfig() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"init_delay_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "数量 秒 after 容器 has started before probes 是 initiated。",
		},
		"timeout_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "数量 秒 after 其中 probe times out.\n默认为 1 second. Minimum 值 是 `1`。",
		},
		"period_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "How often (在 秒) 到 perform probe. 默认为 10 秒. Minimum 值 是 `1`。",
		},
		"success_threshold": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Minimum consecutive successes 对于 probe 到 是 considered successful after having failed. 默认值：`1`. Must 是 1 对于 liveness. Minimum 值 是 `1`。",
		},
		"failure_threshold": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Minimum consecutive failures 对于 probe 到 是 considered failed after having succeeded.默认值：`3`. Minimum 值 是 `1`。",
		},
		"http_get_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "HttpGet detection 路径",
		},
		"http_get_port": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "HttpGet detection 端口",
		},
		"http_get_scheme": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "HttpGet detection scheme. 可用值：`HTTP`，`HTTPS`。",
		},
		"exec_commands": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "列表 execution commands。",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"tcp_socket_port": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "TCP Socket detection 端口",
		},
	}
}

func getEksCiProbeConfig(dMap map[string]interface{}) *tke.LivenessOrReadinessProbe {

	probe := &tke.Probe{}
	httpGet := &tke.HttpGet{}
	exec := &tke.Exec{}
	probeConfig := &tke.LivenessOrReadinessProbe{}

	if v, ok := dMap["init_delay_seconds"]; ok {
		probe.InitialDelaySeconds = helper.IntInt64(v.(int))
	}
	if v, ok := dMap["timeout_seconds"]; ok {
		probe.TimeoutSeconds = helper.IntInt64(v.(int))
	}
	if v, ok := dMap["period_seconds"]; ok {
		probe.PeriodSeconds = helper.IntInt64(v.(int))
	}
	if v, ok := dMap["success_threshold"]; ok {
		probe.SuccessThreshold = helper.IntInt64(v.(int))
	}
	if v, ok := dMap["failure_threshold"]; ok {
		probe.FailureThreshold = helper.IntInt64(v.(int))
	}

	if v, ok := dMap["http_get_path"]; ok {
		httpGet.Path = helper.String(v.(string))
	}
	if v, ok := dMap["http_get_port"]; ok {
		httpGet.Port = helper.IntInt64(v.(int))
	}
	if v, ok := dMap["http_get_scheme"]; ok {
		httpGet.Scheme = helper.String(v.(string))
	}
	if v, ok := dMap["exec_commands"]; ok {
		exec.Commands = helper.InterfacesStringsPoint(v.([]interface{}))
	}
	if v := dMap["tcp_socket_port"].(int); v != 0 {
		probeConfig.TcpSocket = &tke.TcpSocket{
			Port: helper.IntUint64(v),
		}
	} else {
		probeConfig.HttpGet = httpGet
	}

	probeConfig.Probe = probe
	probeConfig.Exec = exec

	return probeConfig
}

func getImageRegistryCredentials(raw []interface{}) []*tke.ImageRegistryCredential {
	var credentials []*tke.ImageRegistryCredential

	for i := range raw {
		cred := &tke.ImageRegistryCredential{}
		item := raw[i].(map[string]interface{})
		if v, ok := item["server"]; ok {
			cred.Server = helper.String(v.(string))
		}
		if v, ok := item["username"]; ok {
			cred.Username = helper.String(v.(string))
		}
		if v, ok := item["password"]; ok {
			cred.Password = helper.String(v.(string))
		}
		if v, ok := item["name"]; ok {
			cred.Name = helper.String(v.(string))
		}
		credentials = append(credentials, cred)
	}

	return credentials
}

func getEksCiContainerList(d *schema.ResourceData, isInit bool) []*tke.Container {
	key := "container"
	if isInit {
		key = "init_container"
	}
	raw := d.Get(key).([]interface{})
	containers := make([]*tke.Container, 0, len(raw))

	for i := range raw {
		var (
			item       = raw[i].(map[string]interface{})
			image      = item["image"].(string)
			name       = item["name"].(string)
			cpu        = item["cpu"].(float64)
			workingDir = item["working_dir"].(string)
			memory     = item["memory"].(float64)
			container  = &tke.Container{
				Image:           helper.String(image),
				Name:            helper.String(name),
				EnvironmentVars: []*tke.EnvironmentVariable{},
			}
		)

		if cpu != 0 {
			container.Cpu = helper.Float64(cpu)
		}
		if memory != 0 {
			container.Memory = helper.Float64(memory)
		}
		if workingDir != "" {
			container.WorkingDir = helper.String(workingDir)
		}
		if v, ok := item["commands"]; ok {
			container.Commands = helper.InterfacesStringsPoint(v.([]interface{}))
		}
		if v, ok := item["args"]; ok {
			container.Args = helper.InterfacesStringsPoint(v.([]interface{}))
		}
		if v, ok := item["env_vars"]; ok {
			envVars := v.(map[string]interface{})
			for key, value := range envVars {
				variable := &tke.EnvironmentVariable{
					Name:  &key,
					Value: helper.String(value.(string)),
				}
				container.EnvironmentVars = append(container.EnvironmentVars, variable)
			}
		}

		if !isInit {
			if v := item["liveness_probe"]; v != nil {
				probConfig := v.([]interface{})
				if len(probConfig) > 0 {
					container.LivenessProbe = getEksCiProbeConfig(probConfig[0].(map[string]interface{}))
				}
			}

			if v := item["readiness_probe"]; v != nil {
				probConfig := v.([]interface{})
				if len(probConfig) > 0 {
					container.ReadinessProbe = getEksCiProbeConfig(probConfig[0].(map[string]interface{}))
				}
			}
		}

		containers = append(containers, container)
	}

	return containers
}

func getEksCiVolume(nfsVolumes, cbsVolumes []interface{}) *tke.EksCiVolume {
	eksCiVolume := &tke.EksCiVolume{}

	if len(nfsVolumes) > 0 {
		eksCiVolume.NfsVolumes = []*tke.NfsVolume{}
		for i := range nfsVolumes {
			var (
				item     = nfsVolumes[i].(map[string]interface{})
				name     = item["name"].(string)
				server   = item["server"].(string)
				path     = item["path"].(string)
				readOnly = item["read_only"].(bool)
				volume   = &tke.NfsVolume{
					Name:     &name,
					Server:   &server,
					Path:     &path,
					ReadOnly: &readOnly,
				}
			)
			eksCiVolume.NfsVolumes = append(eksCiVolume.NfsVolumes, volume)
		}
	}
	if len(cbsVolumes) > 0 {
		eksCiVolume.CbsVolumes = []*tke.CbsVolume{}
		for i := range cbsVolumes {
			var (
				item   = cbsVolumes[i].(map[string]interface{})
				name   = item["name"].(string)
				diskId = item["disk_id"].(string)
				volume = &tke.CbsVolume{
					Name:      &name,
					CbsDiskId: &diskId,
				}
			)
			eksCiVolume.CbsVolumes = append(eksCiVolume.CbsVolumes, volume)
		}
	}
	return eksCiVolume
}

func getEksCiCreateRequest(d *schema.ResourceData) *tke.CreateEKSContainerInstancesRequest {
	var (
		name           = d.Get("name").(string)
		securityGroups = d.Get("security_groups").([]interface{})
		vpcId          = d.Get("vpc_id").(string)
		subnetId       = d.Get("subnet_id").(string)
		memory         = d.Get("memory").(float64)
		cpu            = d.Get("cpu").(float64)
	)
	request := tke.NewCreateEKSContainerInstancesRequest()
	request.EksCiName = &name
	request.VpcId = &vpcId
	request.SubnetId = &subnetId
	request.Memory = &memory
	request.Cpu = &cpu

	if v, ok := d.GetOk("restart_policy"); ok {
		request.RestartPolicy = helper.String(v.(string))
	}

	nfsVolume, hasNfsVolume := d.GetOk("nfs_volume")
	cbsVolume, hasCbsVolume := d.GetOk("cbs_volume")
	if hasNfsVolume || hasCbsVolume {
		request.EksCiVolume = getEksCiVolume(nfsVolume.([]interface{}), cbsVolume.([]interface{}))
	}

	request.SecurityGroupIds = helper.InterfacesStringsPoint(securityGroups)
	request.Containers = getEksCiContainerList(d, false)
	request.InitContainers = getEksCiContainerList(d, true)

	dnsNames := d.Get("dns_names_servers").([]interface{})
	dnsSearches := d.Get("dns_searches").([]interface{})
	dnsOptions := d.Get("dns_config_options").(map[string]interface{})

	if len(dnsNames) > 0 || len(dnsSearches) > 0 || len(dnsOptions) > 0 {
		request.DnsConfig = &tke.DNSConfig{}
		if len(dnsNames) > 0 {
			request.DnsConfig.Nameservers = helper.InterfacesStringsPoint(dnsNames)
		}
		if len(dnsSearches) > 0 {
			request.DnsConfig.Searches = helper.InterfacesStringsPoint(dnsSearches)
		}
		if len(dnsOptions) > 0 {
			var opts []*tke.DNSConfigOption
			for k, v := range dnsOptions {
				opts = append(opts, &tke.DNSConfigOption{
					Name:  &k,
					Value: helper.String(v.(string)),
				})
			}
			request.DnsConfig.Options = opts
		}
	}

	eipIds := d.Get("existed_eip_ids").([]interface{})
	if len(eipIds) > 0 {
		request.ExistedEipIds = helper.InterfacesStringsPoint(eipIds)
	} else {
		deletePolicy, ok := d.GetOk("eip_delete_policy")
		serviceProvider := d.Get("eip_service_provider").(string)
		maxBandwidthOut := d.Get("eip_max_bandwidth_out").(int)
		request.AutoCreateEip = helper.Bool(d.Get("auto_create_eip").(bool))
		if ok || serviceProvider != "" || maxBandwidthOut != 0 {
			request.AutoCreateEipAttribute = &tke.EipAttribute{
				InternetServiceProvider: &serviceProvider,
				InternetMaxBandwidthOut: helper.IntUint64(maxBandwidthOut),
			}
			if deletePolicy.(bool) {
				request.AutoCreateEipAttribute.DeletePolicy = helper.String("true")
			} else {
				request.AutoCreateEipAttribute.DeletePolicy = helper.String("Never")
			}
		}
	}

	if v, ok := d.GetOk("cpu_type"); ok {
		request.CpuType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("gpu_type"); ok {
		request.GpuType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("gpu_count"); ok {
		request.GpuCount = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("cam_role_name"); ok {
		request.CamRoleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("image_registry_credential"); ok {
		request.ImageRegistryCredentials = getImageRegistryCredentials(v.([]interface{}))
	}

	return request
}

func getContainerDataInterface(resContainers []*tke.Container, isInit bool) []interface{} {
	containers := make([]interface{}, 0)
	for i := range resContainers {
		item := resContainers[i]
		container := make(map[string]interface{})
		container["name"] = item.Name
		container["image"] = item.Image
		container["commands"] = helper.StringsInterfaces(item.Commands)
		container["args"] = helper.StringsInterfaces(item.Args)
		container["cpu"] = item.Cpu
		container["memory"] = item.Memory
		container["working_dir"] = item.WorkingDir

		if len(item.EnvironmentVars) > 0 {
			envVars := make(map[string]interface{})
			for i := range item.EnvironmentVars {
				env := item.EnvironmentVars[i]
				envVars[*env.Name] = *env.Value
			}
			container["env_vars"] = envVars
		}

		if len(item.VolumeMounts) > 0 {
			volumes := make([]interface{}, 0, len(item.VolumeMounts))
			for i := range item.VolumeMounts {
				volume := item.VolumeMounts[i]
				raw := make(map[string]interface{})
				raw["name"] = volume.Name
				raw["path"] = volume.MountPath
				raw["read_only"] = volume.ReadOnly
				raw["sub_path"] = volume.SubPath
				raw["sub_path_expr"] = volume.SubPathExpr
				raw["mount_propagation"] = volume.MountPropagation
				volumes = append(volumes, raw)
			}
			container["volume_mount"] = volumes
		}

		if !isInit {
			if item.LivenessProbe != nil {
				probe := make(map[string]interface{})
				probeRaw := make([]interface{}, 0)
				if item.LivenessProbe.Probe != nil {
					probe["init_delay_seconds"] = item.LivenessProbe.Probe.InitialDelaySeconds
					probe["timeout_seconds"] = item.LivenessProbe.Probe.TimeoutSeconds
					probe["period_seconds"] = item.LivenessProbe.Probe.PeriodSeconds
					probe["success_threshold"] = item.LivenessProbe.Probe.SuccessThreshold
					probe["failure_threshold"] = item.LivenessProbe.Probe.FailureThreshold
				}
				if item.LivenessProbe.HttpGet != nil {
					probe["http_get_port"] = item.LivenessProbe.HttpGet.Port
					probe["http_get_path"] = item.LivenessProbe.HttpGet.Path
				}
				if item.LivenessProbe.TcpSocket != nil {
					probe["tcp_socket_port"] = item.LivenessProbe.TcpSocket.Port
				}

				probeRaw = append(probeRaw, probe)
				container["liveness_probe"] = probeRaw
			}

			if item.ReadinessProbe != nil {
				probe := make(map[string]interface{})
				probeRaw := make([]interface{}, 0)

				if item.ReadinessProbe.Probe != nil {
					probe["init_delay_seconds"] = item.ReadinessProbe.Probe.InitialDelaySeconds
					probe["timeout_seconds"] = item.ReadinessProbe.Probe.TimeoutSeconds
					probe["period_seconds"] = item.ReadinessProbe.Probe.PeriodSeconds
					probe["success_threshold"] = item.ReadinessProbe.Probe.SuccessThreshold
					probe["failure_threshold"] = item.ReadinessProbe.Probe.FailureThreshold
				}
				if item.ReadinessProbe.HttpGet != nil {
					probe["http_get_port"] = item.ReadinessProbe.HttpGet.Port
					probe["http_get_path"] = item.ReadinessProbe.HttpGet.Path
				}
				if item.ReadinessProbe.TcpSocket != nil {
					probe["tcp_socket_port"] = item.ReadinessProbe.TcpSocket.Port
				}
				probeRaw = append(probeRaw, probe)
				container["readiness_probe"] = probeRaw
			}
		}
		containers = append(containers, container)
	}
	return containers
}

func setDataFromEksDescribeResponse(res *tke.EksCi, d *schema.ResourceData) error {
	d.SetId(*res.EksCiId)
	_ = d.Set("name", res.EksCiName)
	_ = d.Set("status", res.Status)
	_ = d.Set("vpc_id", res.VpcId)
	_ = d.Set("subnet_id", res.SubnetId)
	_ = d.Set("security_groups", helper.StringsInterfaces(res.SecurityGroupIds))
	_ = d.Set("cpu", res.Cpu)
	_ = d.Set("memory", res.Memory)

	if res.EksCiVolume != nil {
		if len(res.EksCiVolume.NfsVolumes) > 0 {
			volumes := make([]interface{}, 0)
			for i := range res.EksCiVolume.NfsVolumes {
				raw := make(map[string]interface{})
				vol := res.EksCiVolume.NfsVolumes[i]
				raw["name"] = vol.Name
				raw["server"] = vol.Server
				raw["path"] = vol.Path
				raw["read_only"] = vol.ReadOnly
				volumes = append(volumes, raw)
			}
			_ = d.Set("nfs_volume", volumes)
		}

		if len(res.EksCiVolume.CbsVolumes) > 0 {
			volumes := make([]interface{}, 0)
			for i := range res.EksCiVolume.CbsVolumes {
				raw := make(map[string]interface{})
				vol := res.EksCiVolume.CbsVolumes[i]
				raw["name"] = vol.Name
				raw["disk_id"] = vol.CbsDiskId
				volumes = append(volumes, raw)
			}
			_ = d.Set("cbs_volume", volumes)
		}

	}

	_ = d.Set("container", getContainerDataInterface(res.Containers, false))

	if len(res.InitContainers) > 0 {
		_ = d.Set("init_container", getContainerDataInterface(res.InitContainers, true))
	}

	if res.CpuType != nil {
		_ = d.Set("cpu_type", res.CpuType)
	}

	if res.GpuType != nil {
		_ = d.Set("gpu_type", res.GpuType)
	}

	if res.GpuCount != nil {
		_ = d.Set("gpu_count", res.GpuCount)
	}

	if res.CamRoleName != nil {
		_ = d.Set("cam_role_name", res.CamRoleName)
	}

	if res.AutoCreatedEipId != nil {
		_ = d.Set("auto_create_eip_id", res.AutoCreatedEipId)
	}

	if res.EipAddress != nil {
		_ = d.Set("eip_address", res.EipAddress)
	}

	if res.PrivateIp != nil {
		_ = d.Set("private_ip", res.PrivateIp)
	}

	return nil
}

func resourceEksCiContainerSchema(isInitContainer bool) map[string]*schema.Schema {
	schemaMap := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "名称 Container。",
		},
		"image": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Image 的 Container。",
		},
		"commands": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Container launch command 列表。",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"args": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Container launch argument 列表。",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"env_vars": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Map 的 环境 variables 的 容器 OS。",
		},
		"cpu": {
			Type:        schema.TypeFloat,
			Optional:    true,
			Description: "数量 cpu core 的 容器。",
		},
		"memory": {
			Type:        schema.TypeFloat,
			Optional:    true,
			Description: "Memory 大小 的 容器。",
		},
		"working_dir": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Container working directory。",
		},
		"volume_mount": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "列表 卷 mount informations。",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Volume 名称",
					},
					"path": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Volume mount 路径",
					},
					"read_only": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "是否volume 是 read-仅。",
					},
					"sub_path": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Volume mount sub-路径",
					},
					"sub_path_expr": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Volume mount sub-路径 expression。",
					},
					"mount_propagation": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Volume mount propagation。",
					},
				},
			},
		},
	}
	if !isInitContainer {
		schemaMap["liveness_probe"] = &schema.Schema{
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Configuration block 的 LivenessProbe。",
			Elem: &schema.Resource{
				Schema: resourceEksCiProbeConfig(),
			},
		}
		schemaMap["readiness_probe"] = &schema.Schema{
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Configuration block 的 ReadinessProbe。",
			Elem: &schema.Resource{
				Schema: resourceEksCiProbeConfig(),
			},
		}

	}
	return schemaMap
}

func ResourceTencentCloudEksContainerInstance() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource was offline and no longer supported.",
		Read:               resourceTencentcloudEKSContainerInstanceRead,
		Create:             resourceTencentcloudEKSContainerInstanceCreate,
		Update:             resourceTencentcloudEKSContainerInstanceUpdate,
		Delete:             resourceTencentcloudEKSContainerInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "名称 EKS 容器 实例。",
			},
			"container": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "列表 容器。",
				Elem: &schema.Resource{
					Schema: resourceEksCiContainerSchema(false),
				},
			},
			"security_groups": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "列表 安全组 ID",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "子网 ID 容器 实例。",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "私有网络 ID",
			},
			"memory": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "Memory 大小. Check https://intl.云.tencent.com/document/product/457/34057 对于 规格 references。",
			},
			"cpu": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "数量 CPU 核数 Check https://intl.云.tencent.com/document/product/457/34057 对于 规格 references。",
			},
			// optional
			"cpu_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "类型 cpu，其中 可以 集合 到 `intel` 或 `amd`. It also support 备份 列表 like `amd,intel` 其中 表示using `intel` 当 `amd` sold out。",
			},
			"gpu_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "类型 GPU. Check https://intl.云.tencent.com/document/product/457/34057 对于 规格 references。",
			},
			"gpu_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Count 的 GPU. Check https://intl.云.tencent.com/document/product/457/34057 对于 规格 references。",
			},
			"restart_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"Always", "Never", "OnFailure"}),
				Description:  "Container 实例 restart 策略. 可用值：`Always`，`Never`，`OnFailure`。",
			},
			"image_registry_credential": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "列表 credentials 其中 pull 从 镜像 registry。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "地址 的 镜像 registry。",
						},
						"username": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "用户名",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "密码",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "名称 credential。",
						},
					},
				},
			},
			"cbs_volume": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "列表 CBS 卷。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "名称 CBS 卷。",
						},
						"disk_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID CBS。",
						},
					},
				},
			},
			"nfs_volume": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "列表 NFS 卷。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "名称 NFS 卷。",
						},
						"server": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "NFS 服务器 地址",
						},
						"path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "NFS 卷 路径",
						},
						"read_only": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "表示是否volume 是 read 仅. 默认为 `false`。",
						},
					},
				},
			},
			"init_container": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "列表 initialized 容器。",
				Elem: &schema.Resource{
					Schema: resourceEksCiContainerSchema(true),
				},
			},
			// DNSConfig
			"dns_names_servers": {
				Type:        schema.TypeList,
				ForceNew:    true,
				Optional:    true,
				Description: "IP Addresses 的 DNS Servers。",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"dns_searches": {
				Type:        schema.TypeList,
				ForceNew:    true,
				Optional:    true,
				Description: "列表 DNS Search 域名",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"dns_config_options": {
				Type:        schema.TypeMap,
				ForceNew:    true,
				Optional:    true,
				Description: "Map 的 DNS 配置 options。",
			},
			// End of DNSConfig
			"existed_eip_ids": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "Existed EIP ID List 其中 用于bind 容器 实例. Conflict 使用 `auto_create_eip` 和 auto create EIP options。",
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"auto_create_eip"},
			},
			"auto_create_eip": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				Description:   "表示是否create EIP instead 的 指定existing EIPs. Conflict 使用 `existed_eip_ids`。",
				ConflictsWith: []string{"existed_eip_ids"},
			},
			"eip_delete_policy": {
				Type: schema.TypeBool,
				// flatten field must not be required
				Optional:      true,
				Description:   "表示weather EIP release 或 不 after 实例 删除. Conflict 使用 `existed_eip_ids`。",
				ConflictsWith: []string{"existed_eip_ids"},
			},
			"eip_service_provider": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "EIP 服务 provider. 默认为 `BGP`，值 `CMCC`,`CTCC`,`CUCC` 是 可用 对于 whitelist customer. Conflict 使用 `existed_eip_ids`。",
				ConflictsWith: []string{"existed_eip_ids"},
			},
			"eip_max_bandwidth_out": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "Maximum outgoing 带宽 到 公有 网络，measured 在 Mbps (Mega bits per second). Conflict 使用 `existed_eip_ids`。",
				ConflictsWith: []string{"existed_eip_ids"},
			},
			"cam_role_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "被授权访问的 CAM 角色名称",
			},
			// computed
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Container 实例状态",
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Container 实例 创建时间。",
			},
			"auto_create_eip_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID EIP 其中 create automatically。",
			},
			"eip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "EIP 地址",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "内网 IP 地址",
			},
		},
	}
}

func resourceTencentcloudEKSContainerInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eks_cluster.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := EksService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var (
		instance *tke.EksCi
		has      bool
	)

	outErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		var err error
		instance, has, err = service.DescribeEksContainerInstanceById(ctx, d.Id())
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if !has {
			return resource.NonRetryableError(fmt.Errorf("EKS container instance %s not exists", d.Id()))
		}
		if shouldEksCiRetryReading(*instance.Status) {
			return resource.RetryableError(fmt.Errorf("EKS container instance %s is now %s, retrying", d.Id(), *instance.Status))
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}

	if !has {
		d.SetId("")
		return nil
	}

	if err := setDataFromEksDescribeResponse(instance, d); err != nil {
		return err
	}

	return nil
}

func resourceTencentcloudEKSContainerInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eks_cluster.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := EksService{client: client}

	request := getEksCiCreateRequest(d)

	id, err := service.CreateEksContainerInstances(ctx, request)

	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceTencentcloudEKSContainerInstanceRead(d, meta)
}

func resourceTencentcloudEKSContainerInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eks_cluster.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	id := d.Id()
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := EksService{client: client}

	request := tke.NewUpdateEKSContainerInstanceRequest()
	request.EksCiId = &id
	var updateAttrs []string

	if d.HasChange("restart_policy") {
		updateAttrs = append(updateAttrs, "restart_policy")
		policy := d.Get("restart_policy").(string)
		request.RestartPolicy = helper.String(policy)
	}

	if d.HasChange("name") {
		updateAttrs = append(updateAttrs, "restart_policy")
		request.Name = helper.String(d.Get("name").(string))
	}

	nfsVolume, hasNfsVolume := d.GetOk("nfs_volume")
	cbsVolume, hasCbsVolume := d.GetOk("cbs_volume")
	if d.HasChange("nfs_volume") || d.HasChange("cbs_volume") {
		if d.HasChange("nfs_volume") {
			updateAttrs = append(updateAttrs, "nfs_volume")
		}
		if d.HasChange("cbs_volume") {
			updateAttrs = append(updateAttrs, "cbs_volume")
		}

		request.EksCiVolume = getEksCiVolume(nfsVolume.([]interface{}), cbsVolume.([]interface{}))
		if !hasNfsVolume {
			request.EksCiVolume.NfsVolumes = []*tke.NfsVolume{}
		}
		if !hasCbsVolume {
			request.EksCiVolume.CbsVolumes = []*tke.CbsVolume{}
		}
	}

	if d.HasChange("container") {
		updateAttrs = append(updateAttrs, "container")
		request.Containers = getEksCiContainerList(d, false)
	}

	if d.HasChange("init_container") {
		updateAttrs = append(updateAttrs, "init_container")
		request.InitContainers = getEksCiContainerList(d, true)
	}

	if d.HasChange("image_registry_credential") {
		updateAttrs = append(updateAttrs, "image_registry_credential")
		request.ImageRegistryCredentials = getImageRegistryCredentials(d.Get("image_registry_credential").([]interface{}))
	}
	if len(updateAttrs) > 0 {
		err := service.UpdateEksContainerInstances(ctx, request)
		if err != nil {
			return err
		}
	}

	return resourceTencentcloudEKSContainerInstanceRead(d, meta)
}

func resourceTencentcloudEKSContainerInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_eks_cluster.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	id := d.Id()
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := EksService{client: client}

	request := tke.NewDeleteEKSContainerInstancesRequest()
	request.EksCiIds = []*string{&id}

	if err := service.DeleteEksContainerInstance(ctx, request); err != nil {
		return err
	}

	return nil
}

func shouldEksCiRetryReading(status string) bool {
	return status == "Pending" || status == "Updating" || status == "Creating"
}
