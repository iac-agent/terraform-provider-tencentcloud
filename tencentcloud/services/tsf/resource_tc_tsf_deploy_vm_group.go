package tsf

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTsfDeployVmGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfDeployVmGroupCreate,
		Read:   resourceTencentCloudTsfDeployVmGroupRead,
		Delete: resourceTencentCloudTsfDeployVmGroupDelete,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "group id。",
			},

			"pkg_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "program package ID。",
			},

			"startup_parameters": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "start args of group。",
			},

			"deploy_desc": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "group 描述",
			},

			"force_start": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "是否allow forced start。",
			},

			"enable_health_check": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "是否enable health check。",
			},

			"health_check_settings": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "When enabling health check，configure the health check settings。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"liveness_probe": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Survival health check. 注意：此字段可能返回 null，表示未找到有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Health check method. HTTP: check through HTTP interface; CMD: check through executing command; TCP: check through establishing TCP connection. 注意：此字段可能返回 null，表示未找到有效值。",
									},
									"initial_delay_seconds": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "The time 延迟 for the container to start the health check. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"timeout_seconds": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "The maximum timeout 周期 for each health check response. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"period_seconds": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "The 时间间隔 for performing health checks. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"success_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "The 数量 consecutive successful health checks 必填 for the backend container to transition from failure to success. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"failure_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "The 数量 consecutive successful health checks 必填 for the backend container to transition from success to failure. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"scheme": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The 协议 用于HTTP health checks. HTTP and HTTPS are supported. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "The 端口 用于health checks，ranging from 1 to 65535. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"path": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The request 路径 for HTTP health checks. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"command": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Computed:    true,
										Description: "The command to be executed for command health checks. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "类型 readiness probe. TSF_DEFAULT represents the default readiness probe of TSF，while K8S_NATIVE represents the native readiness probe of Kubernetes. If this field is not specified，the native readiness probe of Kubernetes is used by default. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
						"readiness_probe": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Readiness health check. 注意：此字段可能返回 null，表示无法获取有效值。",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The health check method. HTTP 表示checking through an HTTP interface，CMD 表示checking through executing a command，and TCP 表示checking through establishing a TCP connection. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"initial_delay_seconds": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "The time to 延迟 the start of the container health check. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"timeout_seconds": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "The maximum timeout 周期 for each health check response. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"period_seconds": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "The 时间间隔 for performing health checks. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"success_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "The 数量 consecutive successful health checks 必填 for the backend container to transition from failure to success. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"failure_threshold": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "The 数量 consecutive successful health checks 必填 for the backend container to transition from success to failure. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"scheme": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The 协议 用于HTTP health checks. HTTP and HTTPS are supported. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "The 端口 用于health checks，ranging from 1 to 65535. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"path": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The request 路径 for HTTP health checks. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"command": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Computed:    true,
										Description: "The command to be executed for command check. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "类型 readiness probe. TSF_DEFAULT represents the default readiness probe of TSF，while K8S_NATIVE represents the native readiness probe of Kubernetes. If this field is not specified，the native readiness probe of Kubernetes is used by default. 注意：此字段可能返回 null，表示无法获取有效值。",
									},
								},
							},
						},
					},
				},
			},

			"update_type": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Update method: 0 for fast update，1 for rolling update。",
			},

			"deploy_beta_enable": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "是否enable beta batch。",
			},

			"deploy_batch": {
				Optional: true,
				Computed: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeFloat,
				},
				Description: "The ratio of instances participating in each batch during rolling release。",
			},

			"deploy_exe_mode": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The execution method of rolling release。",
			},

			"deploy_wait_time": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "The 时间间隔 for each batch during rolling release。",
			},

			"start_script": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The base64-encoded startup script。",
			},

			"stop_script": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The base64-encoded stop script。",
			},

			"incremental_deployment": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "是否perform incremental deployment. The 默认值为 false，which means full update。",
			},

			"jdk_name": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "JDK 名称: konaJDK or openJDK。",
			},

			"jdk_version": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "JDK 版本: 8 or 11(openJDK only support 8)。",
			},

			"agent_profile_list": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "javaagent info: SERVICE_AGENT/OT_AGENT。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Agent 类型",
						},
						"agent_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Agent 版本",
						},
					},
				},
			},

			"warmup_setting": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "warmup setting。",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "是否enable preheating。",
						},
						"warmup_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "warmup time。",
						},
						"curvature": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Preheating curvature，with a 值 between 1 and 5。",
						},
						"enabled_protection": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "是否enable preheating protection. If protection is 已启用 and more than 50% of nodes are in preheating state，preheating will be aborted。",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTsfDeployVmGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_deploy_vm_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request = tsf.NewDeployGroupRequest()
		groupId string
	)
	if v, ok := d.GetOk("group_id"); ok {
		groupId = v.(string)
		request.GroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("pkg_id"); ok {
		request.PkgId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("startup_parameters"); ok {
		request.StartupParameters = helper.String(v.(string))
	}

	if v, ok := d.GetOk("deploy_desc"); ok {
		request.DeployDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("force_start"); ok && v != nil {
		request.ForceStart = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("enable_health_check"); ok && v != nil {
		request.EnableHealthCheck = helper.Bool(v.(bool))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "health_check_settings"); ok {
		healthCheckSettings := tsf.HealthCheckSettings{}
		if livenessProbeMap, ok := helper.InterfaceToMap(dMap, "liveness_probe"); ok {
			healthCheckSetting := tsf.HealthCheckSetting{}
			if v, ok := livenessProbeMap["action_type"]; ok {
				healthCheckSetting.ActionType = helper.String(v.(string))
			}
			if v, ok := livenessProbeMap["initial_delay_seconds"]; ok && v != nil {
				healthCheckSetting.InitialDelaySeconds = helper.IntUint64(v.(int))
			}
			if v, ok := livenessProbeMap["timeout_seconds"]; ok && v != nil {
				healthCheckSetting.TimeoutSeconds = helper.IntUint64(v.(int))
			}
			if v, ok := livenessProbeMap["period_seconds"]; ok && v != nil {
				healthCheckSetting.PeriodSeconds = helper.IntUint64(v.(int))
			}
			if v, ok := livenessProbeMap["success_threshold"]; ok && v != nil {
				healthCheckSetting.SuccessThreshold = helper.IntUint64(v.(int))
			}
			if v, ok := livenessProbeMap["failure_threshold"]; ok && v != nil {
				healthCheckSetting.FailureThreshold = helper.IntUint64(v.(int))
			}
			if v, ok := livenessProbeMap["scheme"]; ok {
				healthCheckSetting.Scheme = helper.String(v.(string))
			}
			if v, ok := livenessProbeMap["port"]; ok && v != nil {
				healthCheckSetting.Port = helper.IntUint64(v.(int))
			}
			if v, ok := livenessProbeMap["path"]; ok {
				healthCheckSetting.Path = helper.String(v.(string))
			}
			if v, ok := livenessProbeMap["command"]; ok {
				commandSet := v.(*schema.Set).List()
				for i := range commandSet {
					command := commandSet[i].(string)
					healthCheckSetting.Command = append(healthCheckSetting.Command, &command)
				}
			}
			if v, ok := livenessProbeMap["type"]; ok {
				healthCheckSetting.Type = helper.String(v.(string))
			}
			healthCheckSettings.LivenessProbe = &healthCheckSetting
		}
		if readinessProbeMap, ok := helper.InterfaceToMap(dMap, "readiness_probe"); ok {
			healthCheckSetting := tsf.HealthCheckSetting{}
			if v, ok := readinessProbeMap["action_type"]; ok {
				healthCheckSetting.ActionType = helper.String(v.(string))
			}
			if v, ok := readinessProbeMap["initial_delay_seconds"]; ok && v != nil {
				healthCheckSetting.InitialDelaySeconds = helper.IntUint64(v.(int))
			}
			if v, ok := readinessProbeMap["timeout_seconds"]; ok && v != nil {
				healthCheckSetting.TimeoutSeconds = helper.IntUint64(v.(int))
			}
			if v, ok := readinessProbeMap["period_seconds"]; ok && v != nil {
				healthCheckSetting.PeriodSeconds = helper.IntUint64(v.(int))
			}
			if v, ok := readinessProbeMap["success_threshold"]; ok && v != nil {
				healthCheckSetting.SuccessThreshold = helper.IntUint64(v.(int))
			}
			if v, ok := readinessProbeMap["failure_threshold"]; ok && v != nil {
				healthCheckSetting.FailureThreshold = helper.IntUint64(v.(int))
			}
			if v, ok := readinessProbeMap["scheme"]; ok {
				healthCheckSetting.Scheme = helper.String(v.(string))
			}
			if v, ok := readinessProbeMap["port"]; ok && v != nil {
				healthCheckSetting.Port = helper.IntUint64(v.(int))
			}
			if v, ok := readinessProbeMap["path"]; ok {
				healthCheckSetting.Path = helper.String(v.(string))
			}
			if v, ok := readinessProbeMap["command"]; ok {
				commandSet := v.(*schema.Set).List()
				for i := range commandSet {
					command := commandSet[i].(string)
					healthCheckSetting.Command = append(healthCheckSetting.Command, &command)
				}
			}
			if v, ok := readinessProbeMap["type"]; ok && v != nil {
				healthCheckSetting.Type = helper.String(v.(string))
			}
			healthCheckSettings.ReadinessProbe = &healthCheckSetting
		}
		request.HealthCheckSettings = &healthCheckSettings
	}

	if v, ok := d.GetOkExists("update_type"); ok && v != nil {
		request.UpdateType = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("deploy_beta_enable"); ok && v != nil {
		request.DeployBetaEnable = helper.Bool(v.(bool))
	}

	if v, _ := d.GetOk("deploy_batch"); v != nil {
		deployBatchListSet := v.(*schema.Set).List()
		for i := range deployBatchListSet {
			deployBatchList := deployBatchListSet[i].(float64)
			request.DeployBatch = append(request.DeployBatch, &deployBatchList)
		}
	}

	if v, ok := d.GetOk("deploy_exe_mode"); ok {
		request.DeployExeMode = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("deploy_wait_time"); ok && v != nil {
		request.DeployWaitTime = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("start_script"); ok {
		request.StartScript = helper.String(v.(string))
	}

	if v, ok := d.GetOk("stop_script"); ok {
		request.StopScript = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("incremental_deployment"); ok && v != nil {
		request.IncrementalDeployment = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("jdk_name"); ok {
		request.JdkName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("jdk_version"); ok {
		request.JdkVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("agent_profile_list"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			agentProfile := tsf.AgentProfile{}
			if v, ok := dMap["agent_type"]; ok {
				agentProfile.AgentType = helper.String(v.(string))
			}
			if v, ok := dMap["agent_version"]; ok {
				agentProfile.AgentVersion = helper.String(v.(string))
			}
			request.AgentProfileList = append(request.AgentProfileList, &agentProfile)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "warmup_setting"); ok {
		warmupSetting := tsf.WarmupSetting{}
		if v, ok := dMap["enabled"]; ok && v != nil {
			warmupSetting.Enabled = helper.Bool(v.(bool))
		}
		if v, ok := dMap["warmup_time"]; ok && v != nil {
			warmupSetting.WarmupTime = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["curvature"]; ok && v != nil {
			warmupSetting.Curvature = helper.IntUint64(v.(int))
		}
		if v, ok := dMap["enabled_protection"]; ok && v != nil {
			warmupSetting.EnabledProtection = helper.Bool(v.(bool))
		}
		request.WarmupSetting = &warmupSetting
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTsfClient().DeployGroup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf deployVmGroup failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(groupId)

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		groupInfo, err := service.DescribeTsfStartGroupById(ctx, groupId)
		if err != nil {
			return tccommon.RetryError(err)
		}
		if groupInfo == nil {
			err = fmt.Errorf("group %s not exists", groupId)
			return resource.NonRetryableError(err)
		}
		if *groupInfo.GroupStatus == "Running" {
			return nil
		}
		if *groupInfo.GroupStatus == "Waiting" || *groupInfo.GroupStatus == "Updating" {
			return resource.RetryableError(fmt.Errorf("deploy vm group status is %s", *groupInfo.GroupStatus))
		}
		err = fmt.Errorf("deploy vm group status is %v, we won't wait for it finish", *groupInfo.GroupStatus)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s deploy vm group, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudTsfDeployVmGroupRead(d, meta)
}

func resourceTencentCloudTsfDeployVmGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_deploy_vm_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	groupId := d.Id()

	deployVmGroup, err := service.DescribeTsfDeployVmGroupById(ctx, groupId)
	if err != nil {
		return err
	}

	if deployVmGroup == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfDeployVmGroup` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if deployVmGroup.GroupId != nil {
		_ = d.Set("group_id", deployVmGroup.GroupId)
	}

	if deployVmGroup.PackageId != nil {
		_ = d.Set("pkg_id", deployVmGroup.PackageId)
	}

	if deployVmGroup.StartupParameters != nil {
		_ = d.Set("startup_parameters", deployVmGroup.StartupParameters)
	}

	if deployVmGroup.DeployDesc != nil {
		_ = d.Set("deploy_desc", deployVmGroup.DeployDesc)
	}

	if deployVmGroup.EnableHealthCheck != nil {
		_ = d.Set("enable_health_check", deployVmGroup.EnableHealthCheck)
	}

	if deployVmGroup.HealthCheckSettings != nil {
		healthCheckSettingsMap := map[string]interface{}{}

		if deployVmGroup.HealthCheckSettings.LivenessProbe != nil {
			livenessProbeMap := map[string]interface{}{}

			if deployVmGroup.HealthCheckSettings.LivenessProbe.ActionType != nil {
				livenessProbeMap["action_type"] = deployVmGroup.HealthCheckSettings.LivenessProbe.ActionType
			}

			if deployVmGroup.HealthCheckSettings.LivenessProbe.InitialDelaySeconds != nil {
				livenessProbeMap["initial_delay_seconds"] = deployVmGroup.HealthCheckSettings.LivenessProbe.InitialDelaySeconds
			}

			if deployVmGroup.HealthCheckSettings.LivenessProbe.TimeoutSeconds != nil {
				livenessProbeMap["timeout_seconds"] = deployVmGroup.HealthCheckSettings.LivenessProbe.TimeoutSeconds
			}

			if deployVmGroup.HealthCheckSettings.LivenessProbe.PeriodSeconds != nil {
				livenessProbeMap["period_seconds"] = deployVmGroup.HealthCheckSettings.LivenessProbe.PeriodSeconds
			}

			if deployVmGroup.HealthCheckSettings.LivenessProbe.SuccessThreshold != nil {
				livenessProbeMap["success_threshold"] = deployVmGroup.HealthCheckSettings.LivenessProbe.SuccessThreshold
			}

			if deployVmGroup.HealthCheckSettings.LivenessProbe.FailureThreshold != nil {
				livenessProbeMap["failure_threshold"] = deployVmGroup.HealthCheckSettings.LivenessProbe.FailureThreshold
			}

			if deployVmGroup.HealthCheckSettings.LivenessProbe.Scheme != nil {
				livenessProbeMap["scheme"] = deployVmGroup.HealthCheckSettings.LivenessProbe.Scheme
			}

			if deployVmGroup.HealthCheckSettings.LivenessProbe.Port != nil {
				livenessProbeMap["port"] = deployVmGroup.HealthCheckSettings.LivenessProbe.Port
			}

			if deployVmGroup.HealthCheckSettings.LivenessProbe.Path != nil {
				livenessProbeMap["path"] = deployVmGroup.HealthCheckSettings.LivenessProbe.Path
			}

			if deployVmGroup.HealthCheckSettings.LivenessProbe.Command != nil {
				livenessProbeMap["command"] = deployVmGroup.HealthCheckSettings.LivenessProbe.Command
			}

			if deployVmGroup.HealthCheckSettings.LivenessProbe.Type != nil {
				livenessProbeMap["type"] = deployVmGroup.HealthCheckSettings.LivenessProbe.Type
			}

			healthCheckSettingsMap["liveness_probe"] = []interface{}{livenessProbeMap}
		}

		if deployVmGroup.HealthCheckSettings.ReadinessProbe != nil {
			readinessProbeMap := map[string]interface{}{}

			if deployVmGroup.HealthCheckSettings.ReadinessProbe.ActionType != nil {
				readinessProbeMap["action_type"] = deployVmGroup.HealthCheckSettings.ReadinessProbe.ActionType
			}

			if deployVmGroup.HealthCheckSettings.ReadinessProbe.InitialDelaySeconds != nil {
				readinessProbeMap["initial_delay_seconds"] = deployVmGroup.HealthCheckSettings.ReadinessProbe.InitialDelaySeconds
			}

			if deployVmGroup.HealthCheckSettings.ReadinessProbe.TimeoutSeconds != nil {
				readinessProbeMap["timeout_seconds"] = deployVmGroup.HealthCheckSettings.ReadinessProbe.TimeoutSeconds
			}

			if deployVmGroup.HealthCheckSettings.ReadinessProbe.PeriodSeconds != nil {
				readinessProbeMap["period_seconds"] = deployVmGroup.HealthCheckSettings.ReadinessProbe.PeriodSeconds
			}

			if deployVmGroup.HealthCheckSettings.ReadinessProbe.SuccessThreshold != nil {
				readinessProbeMap["success_threshold"] = deployVmGroup.HealthCheckSettings.ReadinessProbe.SuccessThreshold
			}

			if deployVmGroup.HealthCheckSettings.ReadinessProbe.FailureThreshold != nil {
				readinessProbeMap["failure_threshold"] = deployVmGroup.HealthCheckSettings.ReadinessProbe.FailureThreshold
			}

			if deployVmGroup.HealthCheckSettings.ReadinessProbe.Scheme != nil {
				readinessProbeMap["scheme"] = deployVmGroup.HealthCheckSettings.ReadinessProbe.Scheme
			}

			if deployVmGroup.HealthCheckSettings.ReadinessProbe.Port != nil {
				readinessProbeMap["port"] = deployVmGroup.HealthCheckSettings.ReadinessProbe.Port
			}

			if deployVmGroup.HealthCheckSettings.ReadinessProbe.Path != nil {
				readinessProbeMap["path"] = deployVmGroup.HealthCheckSettings.ReadinessProbe.Path
			}

			if deployVmGroup.HealthCheckSettings.ReadinessProbe.Command != nil {
				readinessProbeMap["command"] = deployVmGroup.HealthCheckSettings.ReadinessProbe.Command
			}

			if deployVmGroup.HealthCheckSettings.ReadinessProbe.Type != nil {
				readinessProbeMap["type"] = deployVmGroup.HealthCheckSettings.ReadinessProbe.Type
			}

			healthCheckSettingsMap["readiness_probe"] = []interface{}{readinessProbeMap}
		}

		_ = d.Set("health_check_settings", []interface{}{healthCheckSettingsMap})
	}

	if deployVmGroup.UpdateType != nil {
		_ = d.Set("update_type", deployVmGroup.UpdateType)
	}

	if deployVmGroup.DeployBetaEnable != nil {
		_ = d.Set("deploy_beta_enable", deployVmGroup.DeployBetaEnable)
	}

	if deployVmGroup.DeployBatch != nil {
		_ = d.Set("deploy_batch", deployVmGroup.DeployBatch)
	}

	if deployVmGroup.DeployExeMode != nil {
		_ = d.Set("deploy_exe_mode", deployVmGroup.DeployExeMode)
	}

	if deployVmGroup.DeployWaitTime != nil {
		_ = d.Set("deploy_wait_time", deployVmGroup.DeployWaitTime)
	}

	if deployVmGroup.StartScript != nil {
		_ = d.Set("start_script", deployVmGroup.StartScript)
	}

	if deployVmGroup.StopScript != nil {
		_ = d.Set("stop_script", deployVmGroup.StopScript)
	}

	if deployVmGroup.AgentProfileList != nil {
		agentProfileListList := []interface{}{}
		for _, agentProfileList := range deployVmGroup.AgentProfileList {
			agentProfileListMap := map[string]interface{}{}

			if agentProfileList.AgentType != nil {
				agentProfileListMap["agent_type"] = agentProfileList.AgentType
			}

			if agentProfileList.AgentVersion != nil {
				agentProfileListMap["agent_version"] = agentProfileList.AgentVersion
			}

			agentProfileListList = append(agentProfileListList, agentProfileListMap)
		}

		_ = d.Set("agent_profile_list", agentProfileListList)

	}

	if deployVmGroup.WarmupSetting != nil {
		warmupSettingMap := map[string]interface{}{}

		if deployVmGroup.WarmupSetting.Enabled != nil {
			warmupSettingMap["enabled"] = deployVmGroup.WarmupSetting.Enabled
		}

		if deployVmGroup.WarmupSetting.WarmupTime != nil {
			warmupSettingMap["warmup_time"] = deployVmGroup.WarmupSetting.WarmupTime
		}

		if deployVmGroup.WarmupSetting.Curvature != nil {
			warmupSettingMap["curvature"] = deployVmGroup.WarmupSetting.Curvature
		}

		if deployVmGroup.WarmupSetting.EnabledProtection != nil {
			warmupSettingMap["enabled_protection"] = deployVmGroup.WarmupSetting.EnabledProtection
		}

		_ = d.Set("warmup_setting", []interface{}{warmupSettingMap})
	}

	return nil
}

func resourceTencentCloudTsfDeployVmGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_deploy_vm_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
