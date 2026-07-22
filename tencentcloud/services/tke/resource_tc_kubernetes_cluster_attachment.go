package tke

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svcas "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/as"
)

func ResourceTencentCloudKubernetesClusterAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKubernetesClusterAttachmentCreate,
		Read:   resourceTencentCloudKubernetesClusterAttachmentRead,
		Delete: resourceTencentCloudKubernetesClusterAttachmentDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID 集群。",
			},

			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID CVM 实例，此 cvm 将 reinstall 系统。",
			},

			"image_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "ID Node 镜像。",
			},

			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Sensitive:    true,
				Description:  "密码 到 访问，should 是 集合 如果 `key_ids` 不 集合。",
				ValidateFunc: tccommon.ValidateAsConfigPassword,
			},

			"key_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "键 pair 到 使用 对于 实例，它 looks like skey-16jig7tx，它 should 是 集合 如果 `密码` 不 集合。",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "主机 名称 attached 实例. Dot (.) 和 dash (-) 不能 是 使用 作为 first 和 last 字符 的 HostName 和 不能 是 使用 consecutively. Windows 示例: 长度 的 名称 character 是 [2，15]，letters (capitalization 是 不 restricted)，numbers 和 dashes (-) 是 allowed，dots (.) 是 不 支持，和 不 all numbers 是 allowed. Examples 的 other types (Linux，etc.): character 长度 是 [2，60]，和 多个 dots 是 allowed. There 是 segment between dots. Each segment allows letters (使用 无 limitation 在 capitalization)，numbers 和 dashes (-)。",
			},

			"worker_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Deploy machine 配置 信息 的 'WORKER'，commonly 用于attach existing 实例。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_target": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "挂载目标 默认为 不 mounting。",
						},
						"docker_graph_path": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Description: "Docker graph 路径 默认为 determined 通过 平台 (currently /var/lib/containerd 对于 containerd-based nodes)。",
						},
						"data_disk": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    11,
							Description: "数据盘配置",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ForceNew:     true,
										Default:      "CLOUD_PREMIUM",
										Description:  "Types 的 磁盘. 有效 值: `LOCAL_BASIC`，`LOCAL_SSD`，`CLOUD_BASIC`，`CLOUD_PREMIUM`，`CLOUD_SSD`，`CLOUD_HSSD`，`CLOUD_TSSD` 和 `CLOUD_BSSD`。",
										ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.SYSTEM_DISK_ALLOW_TYPE),
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Optional:    true,
										ForceNew:    true,
										Default:     0,
										Description: "Volume 的 磁盘 （GB）。 默认为 `0`。",
									},
									"file_system": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Default:     "",
										Description: "File 系统，e.g. `ext3/ext4/xfs`。",
									},
									"auto_format_and_mount": {
										Type:        schema.TypeBool,
										Optional:    true,
										ForceNew:    true,
										Default:     false,
										Description: "Indicate 是否auto 格式 和 mount 或 不. 默认为 `false`。",
									},
									"mount_target": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Default:     "",
										Description: "挂载目标",
									},
									"disk_partition": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "名称 device 或 分区 到 mount. NOTE: 此 argument doesn't support setting 在 节点 池，或 将 leads 到 mount 错误",
									},
								},
							},
						},
						"extra_args": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "Custom 参数 信息 related 到 节点. 此 是 white-列表 参数。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"user_data": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Base64-encoded 用户 Data text， 长度 限制 是 16KB。",
						},
						"pre_start_user_script": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Base64-encoded 用户 脚本，executed before initializing 节点，currently 仅 effective 对于 adding existing nodes。",
						},
						"taints": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "Node taint。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "键 的 taint。",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "值 的 taint。",
									},
									"effect": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Effect 的 taint。",
									},
								},
							},
						},
						"is_schedule": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Default:     true,
							Deprecated:  "This argument was deprecated, use `unschedulable` instead.",
							Description: "Indicate 到 调度 adding 节点 或 不. 默认为 true。",
						},
						"desired_pod_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "Indicate 到 集合 desired pod 数量 在 节点. 有效 当 集群 是 podCIDR。",
						},
						"gpu_args": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "GPU 驱动 参数。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mig_enable": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "是否enable MIG。",
									},
									"driver": {
										Type:         schema.TypeMap,
										Optional:     true,
										Description:  "GPU 驱动 版本 格式 like: `{ 版本: String，名称: String }`. `版本`: 版本 的 GPU 驱动 或 CUDA; `名称`: 名称 GPU 驱动 或 CUDA。",
										ValidateFunc: tccommon.ValidateTkeGpuDriverVersion,
									},
									"cuda": {
										Type:         schema.TypeMap,
										Optional:     true,
										Description:  "CUDA 版本 格式 like: `{ 版本: String，名称: String }`. `版本`: 版本 的 GPU 驱动 或 CUDA; `名称`: 名称 GPU 驱动 或 CUDA。",
										ValidateFunc: tccommon.ValidateTkeGpuDriverVersion,
									},
									"cudnn": {
										Type:         schema.TypeMap,
										Optional:     true,
										Description:  "cuDNN 版本 格式 like: `{ 版本: String，名称: String，doc_name: String，dev_name: String }`. `版本`: cuDNN 版本; `名称`: cuDNN 名称; `doc_name`: Doc 名称 cuDNN; `dev_name`: Dev 名称 cuDNN。",
										ValidateFunc: tccommon.ValidateTkeGpuDriverVersion,
									},
									"custom_driver": {
										Type:        schema.TypeMap,
										Optional:    true,
										Description: "Custom GPU 驱动. 格式 like: `{地址: String}`. `地址`: URL 的 自定义 GPU 驱动 地址",
									},
								},
							},
						},
					},
				},
			},

			"worker_config_overrides": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Override variable worker_config，commonly 用于attach existing 实例。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_target": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Deprecated:  "This argument was no longer supported by TencentCloud TKE.",
							Description: "挂载目标 默认为 不 mounting。",
						},
						"docker_graph_path": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Computed:    true,
							Deprecated:  "This argument was no longer supported by TencentCloud TKE.",
							Description: "Docker graph 路径 默认为 determined 通过 平台 (currently /var/lib/containerd 对于 containerd-based nodes)。",
						},
						"data_disk": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    11,
							Description: "数据盘配置",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ForceNew:     true,
										Default:      "CLOUD_PREMIUM",
										Description:  "Types 的 磁盘. 有效 值: `LOCAL_BASIC`，`LOCAL_SSD`，`CLOUD_BASIC`，`CLOUD_PREMIUM`，`CLOUD_SSD`，`CLOUD_HSSD`，`CLOUD_TSSD` 和 `CLOUD_BSSD`。",
										ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.SYSTEM_DISK_ALLOW_TYPE),
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Optional:    true,
										ForceNew:    true,
										Default:     0,
										Description: "Volume 的 磁盘 （GB）。 默认为 `0`。",
									},
									"file_system": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Default:     "",
										Description: "File 系统，e.g. `ext3/ext4/xfs`。",
									},
									"auto_format_and_mount": {
										Type:        schema.TypeBool,
										Optional:    true,
										ForceNew:    true,
										Default:     false,
										Description: "Indicate 是否auto 格式 和 mount 或 不. 默认为 `false`。",
									},
									"mount_target": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Default:     "",
										Description: "挂载目标",
									},
									"disk_partition": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "名称 device 或 分区 到 mount. NOTE: 此 argument doesn't support setting 在 节点 池，或 将 leads 到 mount 错误",
									},
								},
							},
						},
						"extra_args": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Deprecated:  "This argument was no longer supported by TencentCloud TKE.",
							Description: "Custom 参数 信息 related 到 节点. 此 是 white-列表 参数。",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"user_data": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Deprecated:  "This argument was no longer supported by TencentCloud TKE.",
							Description: "Base64-encoded 用户 Data text， 长度 限制 是 16KB。",
						},
						"pre_start_user_script": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Deprecated:  "This argument was no longer supported by TencentCloud TKE.",
							Description: "Base64-encoded 用户 脚本，executed before initializing 节点，currently 仅 effective 对于 adding existing nodes。",
						},
						"is_schedule": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Default:     true,
							Deprecated:  "This argument was deprecated, use `unschedulable` instead.",
							Description: "Indicate 到 调度 adding 节点 或 不. 默认为 true。",
						},
						"desired_pod_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "Indicate 到 集合 desired pod 数量 在 节点. 有效 当 集群 是 podCIDR。",
						},
						"gpu_args": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "GPU 驱动 参数。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mig_enable": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "是否enable MIG。",
									},
									"driver": {
										Type:         schema.TypeMap,
										Optional:     true,
										Description:  "GPU 驱动 版本 格式 like: `{ 版本: String，名称: String }`. `版本`: 版本 的 GPU 驱动 或 CUDA; `名称`: 名称 GPU 驱动 或 CUDA。",
										ValidateFunc: tccommon.ValidateTkeGpuDriverVersion,
									},
									"cuda": {
										Type:         schema.TypeMap,
										Optional:     true,
										Description:  "CUDA 版本 格式 like: `{ 版本: String，名称: String }`. `版本`: 版本 的 GPU 驱动 或 CUDA; `名称`: 名称 GPU 驱动 或 CUDA。",
										ValidateFunc: tccommon.ValidateTkeGpuDriverVersion,
									},
									"cudnn": {
										Type:         schema.TypeMap,
										Optional:     true,
										Description:  "cuDNN 版本 格式 like: `{ 版本: String，名称: String，doc_name: String，dev_name: String }`. `版本`: cuDNN 版本; `名称`: cuDNN 名称; `doc_name`: Doc 名称 cuDNN; `dev_name`: Dev 名称 cuDNN。",
										ValidateFunc: tccommon.ValidateTkeGpuDriverVersion,
									},
									"custom_driver": {
										Type:        schema.TypeMap,
										Optional:    true,
										Description: "Custom GPU 驱动. 格式 like: `{地址: String}`. `地址`: URL 的 自定义 GPU 驱动 地址",
									},
								},
							},
						},
					},
				},
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Labels 的 tke attachment exits CVM。",
			},

			"unschedulable": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     0,
				Description: "Sets 是否joining 节点 participates 在 调度. 默认为 `0`，其中 表示 它 participates 在 scheduling. Non-zero(eg: `1`) 数量 表示 它 does 不 participate 在 scheduling。",
			},

			"security_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "A 列表 安全 组 IDs after attach 到 集群。",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State 的 节点。",
			},
		},
	}
}

func resourceTencentCloudKubernetesClusterAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		instanceId string
		clusterId  string
	)
	var (
		request  = tke.NewAddExistedInstancesRequest()
		response = tke.NewAddExistedInstancesResponse()
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
	}

	request.InstanceIds = []*string{helper.String(instanceId)}

	if v, ok := d.GetOk("image_id"); ok {
		request.ImageId = helper.String(v.(string))
	}

	loginSettings := tke.LoginSettings{}
	if v, ok := d.GetOk("password"); ok {
		loginSettings.Password = helper.String(v.(string))
	}
	request.LoginSettings = &loginSettings

	if instanceAdvancedSettingsMap, ok := helper.InterfacesHeadMap(d, "worker_config"); ok {
		instanceAdvancedSettings := tke.InstanceAdvancedSettings{}
		if v, ok := instanceAdvancedSettingsMap["mount_target"]; ok {
			instanceAdvancedSettings.MountTarget = helper.String(v.(string))
		}
		if v, ok := instanceAdvancedSettingsMap["data_disk"]; ok {
			for _, item := range v.([]interface{}) {
				dataDisksMap := item.(map[string]interface{})
				dataDisk := tke.DataDisk{}
				if v, ok := dataDisksMap["disk_type"]; ok {
					dataDisk.DiskType = helper.String(v.(string))
				}
				if v, ok := dataDisksMap["file_system"]; ok {
					dataDisk.FileSystem = helper.String(v.(string))
				}
				if v, ok := dataDisksMap["auto_format_and_mount"]; ok {
					dataDisk.AutoFormatAndMount = helper.Bool(v.(bool))
				}
				if v, ok := dataDisksMap["mount_target"]; ok {
					dataDisk.MountTarget = helper.String(v.(string))
				}
				if v, ok := dataDisksMap["disk_partition"]; ok {
					dataDisk.DiskPartition = helper.String(v.(string))
				}
				if v, ok := dataDisksMap["disk_size"]; ok {
					dataDisk.DiskSize = helper.IntInt64(v.(int))
				}
				instanceAdvancedSettings.DataDisks = append(instanceAdvancedSettings.DataDisks, &dataDisk)
			}
		}
		if v, ok := instanceAdvancedSettingsMap["user_data"]; ok {
			instanceAdvancedSettings.UserScript = helper.String(v.(string))
		}
		if v, ok := instanceAdvancedSettingsMap["pre_start_user_script"]; ok {
			instanceAdvancedSettings.PreStartUserScript = helper.String(v.(string))
		}
		if v, ok := instanceAdvancedSettingsMap["taints"]; ok {
			for _, item := range v.([]interface{}) {
				taintsMap := item.(map[string]interface{})
				taint := tke.Taint{}
				if v, ok := taintsMap["key"]; ok {
					taint.Key = helper.String(v.(string))
				}
				if v, ok := taintsMap["value"]; ok {
					taint.Value = helper.String(v.(string))
				}
				if v, ok := taintsMap["effect"]; ok {
					taint.Effect = helper.String(v.(string))
				}
				instanceAdvancedSettings.Taints = append(instanceAdvancedSettings.Taints, &taint)
			}
		}
		if v, ok := instanceAdvancedSettingsMap["docker_graph_path"]; ok {
			instanceAdvancedSettings.DockerGraphPath = helper.String(v.(string))
		}
		if v, ok := instanceAdvancedSettingsMap["desired_pod_num"]; ok {
			instanceAdvancedSettings.DesiredPodNumber = helper.IntInt64(v.(int))
		}
		if gPUArgsMap, ok := helper.ConvertInterfacesHeadToMap(instanceAdvancedSettingsMap["gpu_args"]); ok {
			gPUArgs := tke.GPUArgs{}
			if v, ok := gPUArgsMap["mig_enable"]; ok {
				gPUArgs.MIGEnable = helper.Bool(v.(bool))
			}
			instanceAdvancedSettings.GPUArgs = &gPUArgs
		}
		request.InstanceAdvancedSettings = &instanceAdvancedSettings
	}

	if v, ok := d.GetOk("hostname"); ok {
		request.HostName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_groups"); ok {
		securityGroupIdsSet := v.([]interface{})
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, helper.String(securityGroupIds))
		}
	}

	if v, ok := d.GetOk("worker_config_overrides"); ok {
		for _, item := range v.([]interface{}) {
			instanceAdvancedSettingsOverridesMap := item.(map[string]interface{})
			instanceAdvancedSettings := tke.InstanceAdvancedSettings{}
			if v, ok := instanceAdvancedSettingsOverridesMap["mount_target"]; ok {
				instanceAdvancedSettings.MountTarget = helper.String(v.(string))
			}
			if v, ok := instanceAdvancedSettingsOverridesMap["data_disk"]; ok {
				for _, item := range v.([]interface{}) {
					dataDisksMap := item.(map[string]interface{})
					dataDisk := tke.DataDisk{}
					if v, ok := dataDisksMap["disk_type"]; ok {
						dataDisk.DiskType = helper.String(v.(string))
					}
					if v, ok := dataDisksMap["file_system"]; ok {
						dataDisk.FileSystem = helper.String(v.(string))
					}
					if v, ok := dataDisksMap["auto_format_and_mount"]; ok {
						dataDisk.AutoFormatAndMount = helper.Bool(v.(bool))
					}
					if v, ok := dataDisksMap["mount_target"]; ok {
						dataDisk.MountTarget = helper.String(v.(string))
					}
					if v, ok := dataDisksMap["disk_partition"]; ok {
						dataDisk.DiskPartition = helper.String(v.(string))
					}
					if v, ok := dataDisksMap["disk_size"]; ok {
						dataDisk.DiskSize = helper.IntInt64(v.(int))
					}
					instanceAdvancedSettings.DataDisks = append(instanceAdvancedSettings.DataDisks, &dataDisk)
				}
			}
			if v, ok := instanceAdvancedSettingsOverridesMap["user_data"]; ok {
				instanceAdvancedSettings.UserScript = helper.String(v.(string))
			}
			if v, ok := instanceAdvancedSettingsOverridesMap["pre_start_user_script"]; ok {
				instanceAdvancedSettings.PreStartUserScript = helper.String(v.(string))
			}
			if v, ok := instanceAdvancedSettingsOverridesMap["docker_graph_path"]; ok {
				instanceAdvancedSettings.DockerGraphPath = helper.String(v.(string))
			}
			if v, ok := instanceAdvancedSettingsOverridesMap["desired_pod_num"]; ok {
				instanceAdvancedSettings.DesiredPodNumber = helper.IntInt64(v.(int))
			}
			if gPUArgsMap, ok := helper.ConvertInterfacesHeadToMap(instanceAdvancedSettingsOverridesMap["gpu_args"]); ok {
				gPUArgs2 := tke.GPUArgs{}
				if v, ok := gPUArgsMap["mig_enable"]; ok {
					gPUArgs2.MIGEnable = helper.Bool(v.(bool))
				}
				instanceAdvancedSettings.GPUArgs = &gPUArgs2
			}
			request.InstanceAdvancedSettingsOverrides = append(request.InstanceAdvancedSettingsOverrides, &instanceAdvancedSettings)
		}
	}

	if err := resourceTencentCloudKubernetesClusterAttachmentCreatePostFillRequest0(ctx, request); err != nil {
		return err
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeClient().AddExistedInstancesWithContext(ctx, request)
		if e != nil {
			return resourceTencentCloudKubernetesClusterAttachmentCreateRequestOnError0(ctx, request, e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create kubernetes cluster attachment failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	if err := resourceTencentCloudKubernetesClusterAttachmentCreatePostHandleResponse0(ctx, response); err != nil {
		return err
	}

	d.SetId(strings.Join([]string{instanceId, clusterId}, "_"))

	return resourceTencentCloudKubernetesClusterAttachmentRead(d, meta)
}

func resourceTencentCloudKubernetesClusterAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), "_")
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	clusterId := idSplit[1]

	_ = d.Set("instance_id", instanceId)

	_ = d.Set("cluster_id", clusterId)

	respData, err := service.DescribeKubernetesClusterAttachmentById(ctx, clusterId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `kubernetes_cluster_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	respData1, err := service.DescribeKubernetesClusterAttachmentById1(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData1 == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `kubernetes_cluster_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData1.LoginSettings != nil {
		if respData1.LoginSettings.KeyIds != nil {
			_ = d.Set("key_ids", respData1.LoginSettings.KeyIds)
		}

	}

	if respData1.SecurityGroupIds != nil {
		_ = d.Set("security_groups", respData1.SecurityGroupIds)
	}

	if respData1.ImageId != nil {
		_ = d.Set("image_id", respData1.ImageId)
	}

	var respData2 *tke.Instance
	reqErr2 := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKubernetesClusterAttachmentById2(ctx, instanceId, clusterId)
		if e != nil {
			return resourceTencentCloudKubernetesClusterAttachmentReadRequestOnError2(ctx, result, e)
		}
		if err := resourceTencentCloudKubernetesClusterAttachmentReadRequestOnSuccess2(ctx, result); err != nil {
			return err
		}
		respData2 = result
		return nil
	})
	if reqErr2 != nil {
		log.Printf("[CRITAL]%s read kubernetes cluster attachment failed, reason:%+v", logId, reqErr2)
		return reqErr2
	}

	if respData2 == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `kubernetes_cluster_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData2.InstanceAdvancedSettings != nil {
		if respData2.InstanceAdvancedSettings.Unschedulable != nil {
			_ = d.Set("unschedulable", respData2.InstanceAdvancedSettings.Unschedulable)
		}

	}

	if respData2.InstanceState != nil {
		_ = d.Set("state", respData2.InstanceState)
	}

	return nil
}

func resourceTencentCloudKubernetesClusterAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), "_")
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	clusterId := idSplit[1]

	var (
		request  = tke.NewDeleteClusterInstancesRequest()
		response = tke.NewDeleteClusterInstancesResponse()
	)

	request.ClusterId = helper.String(clusterId)

	request.InstanceIds = []*string{helper.String(instanceId)}

	instanceDeleteMode := "retain"
	request.InstanceDeleteMode = &instanceDeleteMode

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeClient().DeleteClusterInstancesWithContext(ctx, request)
		if e != nil {
			return resourceTencentCloudKubernetesClusterAttachmentDeleteRequestOnError0(ctx, e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete kubernetes cluster attachment failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil
}
