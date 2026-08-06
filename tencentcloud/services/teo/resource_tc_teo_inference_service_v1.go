package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoInferenceServiceV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoInferenceServiceV1Create,
		Read:   resourceTencentCloudTeoInferenceServiceV1Read,
		Update: resourceTencentCloudTeoInferenceServiceV1Update,
		Delete: resourceTencentCloudTeoInferenceServiceV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Inference service name. The name is up to 30 characters, only supports lowercase letters, numbers, hyphens, starts with a letter, ends with a number or letter, and does not support duplicates.",
			},
			"listen_port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The port that the inference service listens on. Only supports integers between 1-65535.",
			},
			"containers": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Container configuration for the inference service. Currently, only 1 container is supported.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Image type. Valid values: `TCR` (Tencent Cloud Container Registry).",
						},
						"tcr_repository_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "TCR repository configuration. Required when ImageType is TCR.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tcr_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "TCR service type. Valid values: `Personal`, `Enterprise`.",
									},
									"image": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Image address.",
									},
									"registry_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Registry instance ID. Required when TCRType is Enterprise.",
									},
									"region_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Region name.",
									},
								},
							},
						},
						"startup_command": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Command executed when the container starts. Defaults to the image's Entrypoint/CMD if not specified. Max length: 1024 characters.",
						},
						"environment_variables": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Environment variables for the container runtime. Up to 10 variables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Variable name. Only uppercase and lowercase letters, numbers, and underscores, must start with a letter or underscore. Max length: 64 characters.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Variable value. Any visible characters. Max length: 2048 characters.",
									},
								},
							},
						},
					},
				},
			},
			"resource_config": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Resource configuration for the inference service.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scaling_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Scaling mode. Valid values: `Auto` (auto scaling based on request volume), `Manual` (fixed instance count).",
						},
						"hardware_spec": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Hardware specification. Note: This field can only be set during creation and cannot be modified afterwards.",
						},
						"auto_scaling_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Auto scaling configuration. Required when ScalingMode is Auto.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min_instance_count": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Minimum number of instances.",
									},
									"scaling_policies": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Scaling policy list. Up to 5 policies.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"policy_name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Policy name. Length: 1-30 characters. Must be unique within the same service.",
												},
												"policy_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Policy type. Valid values: `ScheduledScaling` (scheduled scaling).",
												},
												"scheduled_scaling_policy": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Scheduled scaling policy configuration. Required when PolicyType is ScheduledScaling.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"scheduled_actions": {
																Type:        schema.TypeList,
																Required:    true,
																Description: "Scheduled scaling action list. At least 1, up to 10.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"cron_expression": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Cron expression for triggering the scheduled scaling action. Uses 5-field standard Cron format: minute hour day month weekday.",
																		},
																		"min_instance_count": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Minimum number of instances to adjust to when this scheduled scaling action is triggered.",
																		},
																	},
																},
															},
															"effective_range": {
																Type:        schema.TypeList,
																Required:    true,
																MaxItems:    1,
																Description: "Effective range for the scheduled scaling policy.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"effective_type": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Effective type. Valid values: `LongTerm`, `Custom`.",
																		},
																		"start_date": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Start date for the effective range. Required when EffectiveType is Custom.",
																		},
																		"end_date": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "End date for the effective range. Required when EffectiveType is Custom and must not be earlier than StartDate.",
																		},
																	},
																},
															},
															"time_zone": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Timezone for the scheduled actions, e.g., UTC, Asia/Shanghai. Defaults to UTC.",
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
						"manual_instance_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Manual instance configuration. Required when ScalingMode is Manual.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fixed_instance_count": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Fixed instance count.",
									},
								},
							},
						},
						"concurrency": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Concurrency per instance. Default: 1.",
						},
					},
				},
			},
			"request_paths": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Request path list for the inference service. Up to 20 paths.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description. Max length: 60 characters.",
			},
			"operation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Operation to perform on the inference service. Valid values: `Stop`, `Resume`. This field is not persisted to state.",
			},
			// computed
			"service_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Inference service ID.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Inference service status. Valid values: `Deploying`, `Running`, `Stopping`, `Stopped`, `Exception`, `Banned`.",
			},
			"inference_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Inference access URL.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time in ISO date format.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modification time in ISO date format.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(tccommon.ReadRetryTimeout),
			Update: schema.DefaultTimeout(tccommon.ReadRetryTimeout),
			Delete: schema.DefaultTimeout(tccommon.ReadRetryTimeout),
		},
	}
}

func resourceTencentCloudTeoInferenceServiceV1Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_service_v1.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = teo.NewCreateInferenceServiceRequest()
		response = teo.NewCreateInferenceServiceResponse()
		zoneId   string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(zoneId)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("listen_port"); ok {
		request.ListenPort = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("containers"); ok {
		for _, item := range v.([]interface{}) {
			containerMap := item.(map[string]interface{})
			container := teo.InferenceContainerConfig{}

			if v, ok := containerMap["image_type"].(string); ok && v != "" {
				container.ImageType = helper.String(v)
			}

			if v, ok := containerMap["tcr_repository_config"]; ok {
				for _, tcrItem := range v.([]interface{}) {
					tcrMap := tcrItem.(map[string]interface{})
					tcrConfig := teo.InferenceTCRRepositoryConfig{}

					if v, ok := tcrMap["tcr_type"].(string); ok && v != "" {
						tcrConfig.TCRType = helper.String(v)
					}

					if v, ok := tcrMap["image"].(string); ok && v != "" {
						tcrConfig.Image = helper.String(v)
					}

					if v, ok := tcrMap["registry_id"].(string); ok && v != "" {
						tcrConfig.RegistryId = helper.String(v)
					}

					if v, ok := tcrMap["region_name"].(string); ok && v != "" {
						tcrConfig.RegionName = helper.String(v)
					}

					container.TcrRepositoryConfig = &tcrConfig
				}
			}

			if v, ok := containerMap["startup_command"].(string); ok && v != "" {
				container.StartupCommand = helper.String(v)
			}

			if v, ok := containerMap["environment_variables"]; ok {
				for _, envItem := range v.([]interface{}) {
					envMap := envItem.(map[string]interface{})
					envVar := teo.InferenceEnvironmentVariable{}

					if v, ok := envMap["key"].(string); ok && v != "" {
						envVar.Key = helper.String(v)
					}

					if v, ok := envMap["value"].(string); ok && v != "" {
						envVar.Value = helper.String(v)
					}

					container.EnvironmentVariables = append(container.EnvironmentVariables, &envVar)
				}
			}

			request.Containers = append(request.Containers, &container)
		}
	}

	if v, ok := d.GetOk("resource_config"); ok {
		for _, item := range v.([]interface{}) {
			rcMap := item.(map[string]interface{})
			resConfig := teo.InferenceResourceConfig{}

			if v, ok := rcMap["scaling_mode"].(string); ok && v != "" {
				resConfig.ScalingMode = helper.String(v)
			}

			if v, ok := rcMap["hardware_spec"].(string); ok && v != "" {
				resConfig.HardwareSpec = helper.String(v)
			}

			if v, ok := rcMap["auto_scaling_config"]; ok {
				for _, ascItem := range v.([]interface{}) {
					ascMap := ascItem.(map[string]interface{})
					autoScalingConfig := teo.InferenceAutoScalingConfig{}

					if v, ok := ascMap["min_instance_count"].(int); ok {
						autoScalingConfig.MinInstanceCount = helper.IntInt64(v)
					}

					if v, ok := ascMap["scaling_policies"]; ok {
						for _, spItem := range v.([]interface{}) {
							spMap := spItem.(map[string]interface{})
							scalingPolicy := teo.InferenceScalingPolicy{}

							if v, ok := spMap["policy_name"].(string); ok && v != "" {
								scalingPolicy.PolicyName = helper.String(v)
							}

							if v, ok := spMap["policy_type"].(string); ok && v != "" {
								scalingPolicy.PolicyType = helper.String(v)
							}

							if v, ok := spMap["scheduled_scaling_policy"]; ok {
								for _, sspItem := range v.([]interface{}) {
									sspMap := sspItem.(map[string]interface{})
									scheduledScalingPolicy := teo.InferenceScheduledScalingPolicy{}

									if v, ok := sspMap["scheduled_actions"]; ok {
										for _, saItem := range v.([]interface{}) {
											saMap := saItem.(map[string]interface{})
											scheduledAction := teo.InferenceScheduledScalingAction{}

											if v, ok := saMap["cron_expression"].(string); ok && v != "" {
												scheduledAction.CronExpression = helper.String(v)
											}

											if v, ok := saMap["min_instance_count"].(int); ok {
												scheduledAction.MinInstanceCount = helper.IntInt64(v)
											}

											scheduledScalingPolicy.ScheduledActions = append(scheduledScalingPolicy.ScheduledActions, &scheduledAction)
										}
									}

									if v, ok := sspMap["effective_range"]; ok {
										for _, erItem := range v.([]interface{}) {
											erMap := erItem.(map[string]interface{})
											effectiveRange := teo.InferenceScheduledScalingEffectiveRange{}

											if v, ok := erMap["effective_type"].(string); ok && v != "" {
												effectiveRange.EffectiveType = helper.String(v)
											}

											if v, ok := erMap["start_date"].(string); ok && v != "" {
												effectiveRange.StartDate = helper.String(v)
											}

											if v, ok := erMap["end_date"].(string); ok && v != "" {
												effectiveRange.EndDate = helper.String(v)
											}

											scheduledScalingPolicy.EffectiveRange = &effectiveRange
										}
									}

									if v, ok := sspMap["time_zone"].(string); ok && v != "" {
										scheduledScalingPolicy.TimeZone = helper.String(v)
									}

									scalingPolicy.ScheduledScalingPolicy = &scheduledScalingPolicy
								}
							}

							autoScalingConfig.ScalingPolicies = append(autoScalingConfig.ScalingPolicies, &scalingPolicy)
						}
					}

					resConfig.AutoScalingConfig = &autoScalingConfig
				}
			}

			if v, ok := rcMap["manual_instance_config"]; ok {
				for _, micItem := range v.([]interface{}) {
					micMap := micItem.(map[string]interface{})
					manualInstanceConfig := teo.InferenceManualInstanceConfig{}

					if v, ok := micMap["fixed_instance_count"].(int); ok {
						manualInstanceConfig.FixedInstanceCount = helper.IntInt64(v)
					}

					resConfig.ManualInstanceConfig = &manualInstanceConfig
				}
			}

			if v, ok := rcMap["concurrency"].(int); ok {
				resConfig.Concurrency = helper.IntInt64(v)
			}

			request.ResourceConfig = &resConfig
		}
	}

	if v, ok := d.GetOk("request_paths"); ok {
		pathsSet := v.(*schema.Set)
		for _, path := range pathsSet.List() {
			request.RequestPaths = append(request.RequestPaths, helper.String(path.(string)))
		}
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateInferenceServiceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo inference_service failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo inference_service failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s teo inference_service create response, ServiceId: %v", logId, response.Response.ServiceId)

	if response.Response.ServiceId == nil || *response.Response.ServiceId == "" {
		return fmt.Errorf("Create teo inference_service failed, ServiceId is nil or empty.")
	}

	d.SetId(zoneId + tccommon.FILED_SP + *response.Response.ServiceId)
	return resourceTencentCloudTeoInferenceServiceV1Read(d, meta)
}

func resourceTencentCloudTeoInferenceServiceV1Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_service_v1.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("resource id is broken, id is %s", d.Id())
	}
	zoneId := idSplit[0]
	serviceId := idSplit[1]

	respData, err := service.DescribeTeoInferenceServiceV1ById(ctx, zoneId, serviceId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `teo_inference_service_v1` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	if respData.ServiceId != nil {
		_ = d.Set("service_id", respData.ServiceId)
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.ListenPort != nil {
		_ = d.Set("listen_port", respData.ListenPort)
	}

	if respData.Containers != nil {
		containersList := make([]map[string]interface{}, 0, len(respData.Containers))
		for _, container := range respData.Containers {
			containerMap := map[string]interface{}{}

			if container.ImageType != nil {
				containerMap["image_type"] = container.ImageType
			}

			if container.TcrRepositoryConfig != nil {
				tcrConfigList := make([]map[string]interface{}, 0, 1)
				tcrConfigMap := map[string]interface{}{}
				if container.TcrRepositoryConfig.TCRType != nil {
					tcrConfigMap["tcr_type"] = container.TcrRepositoryConfig.TCRType
				}
				if container.TcrRepositoryConfig.Image != nil {
					tcrConfigMap["image"] = container.TcrRepositoryConfig.Image
				}
				if container.TcrRepositoryConfig.RegistryId != nil {
					tcrConfigMap["registry_id"] = container.TcrRepositoryConfig.RegistryId
				}
				if container.TcrRepositoryConfig.RegionName != nil {
					tcrConfigMap["region_name"] = container.TcrRepositoryConfig.RegionName
				}
				tcrConfigList = append(tcrConfigList, tcrConfigMap)
				containerMap["tcr_repository_config"] = tcrConfigList
			}

			if container.StartupCommand != nil {
				containerMap["startup_command"] = container.StartupCommand
			}

			if container.EnvironmentVariables != nil {
				envVarsList := make([]map[string]interface{}, 0, len(container.EnvironmentVariables))
				for _, envVar := range container.EnvironmentVariables {
					envVarMap := map[string]interface{}{}
					if envVar.Key != nil {
						envVarMap["key"] = envVar.Key
					}
					if envVar.Value != nil {
						envVarMap["value"] = envVar.Value
					}
					envVarsList = append(envVarsList, envVarMap)
				}
				containerMap["environment_variables"] = envVarsList
			}

			containersList = append(containersList, containerMap)
		}
		_ = d.Set("containers", containersList)
	}

	if respData.ResourceConfig != nil {
		rc := respData.ResourceConfig
		rcList := make([]map[string]interface{}, 0, 1)
		rcMap := map[string]interface{}{}

		if rc.ScalingMode != nil {
			rcMap["scaling_mode"] = rc.ScalingMode
		}

		if rc.HardwareSpec != nil {
			rcMap["hardware_spec"] = rc.HardwareSpec
		}

		if rc.AutoScalingConfig != nil {
			ascList := make([]map[string]interface{}, 0, 1)
			ascMap := map[string]interface{}{}

			if rc.AutoScalingConfig.MinInstanceCount != nil {
				ascMap["min_instance_count"] = rc.AutoScalingConfig.MinInstanceCount
			}

			if rc.AutoScalingConfig.ScalingPolicies != nil {
				spList := make([]map[string]interface{}, 0, len(rc.AutoScalingConfig.ScalingPolicies))
				for _, sp := range rc.AutoScalingConfig.ScalingPolicies {
					spMap := map[string]interface{}{}
					if sp.PolicyName != nil {
						spMap["policy_name"] = sp.PolicyName
					}
					if sp.PolicyType != nil {
						spMap["policy_type"] = sp.PolicyType
					}
					if sp.ScheduledScalingPolicy != nil {
						sspList := make([]map[string]interface{}, 0, 1)
						sspMap := map[string]interface{}{}

						if sp.ScheduledScalingPolicy.ScheduledActions != nil {
							saList := make([]map[string]interface{}, 0, len(sp.ScheduledScalingPolicy.ScheduledActions))
							for _, sa := range sp.ScheduledScalingPolicy.ScheduledActions {
								saMap := map[string]interface{}{}
								if sa.CronExpression != nil {
									saMap["cron_expression"] = sa.CronExpression
								}
								if sa.MinInstanceCount != nil {
									saMap["min_instance_count"] = sa.MinInstanceCount
								}
								saList = append(saList, saMap)
							}
							sspMap["scheduled_actions"] = saList
						}

						if sp.ScheduledScalingPolicy.EffectiveRange != nil {
							erList := make([]map[string]interface{}, 0, 1)
							erMap := map[string]interface{}{}
							if sp.ScheduledScalingPolicy.EffectiveRange.EffectiveType != nil {
								erMap["effective_type"] = sp.ScheduledScalingPolicy.EffectiveRange.EffectiveType
							}
							if sp.ScheduledScalingPolicy.EffectiveRange.StartDate != nil {
								erMap["start_date"] = sp.ScheduledScalingPolicy.EffectiveRange.StartDate
							}
							if sp.ScheduledScalingPolicy.EffectiveRange.EndDate != nil {
								erMap["end_date"] = sp.ScheduledScalingPolicy.EffectiveRange.EndDate
							}
							erList = append(erList, erMap)
							sspMap["effective_range"] = erList
						}

						if sp.ScheduledScalingPolicy.TimeZone != nil {
							sspMap["time_zone"] = sp.ScheduledScalingPolicy.TimeZone
						}

						sspList = append(sspList, sspMap)
						spMap["scheduled_scaling_policy"] = sspList
					}
					spList = append(spList, spMap)
				}
				ascMap["scaling_policies"] = spList
			}

			ascList = append(ascList, ascMap)
			rcMap["auto_scaling_config"] = ascList
		}

		if rc.ManualInstanceConfig != nil {
			micList := make([]map[string]interface{}, 0, 1)
			micMap := map[string]interface{}{}
			if rc.ManualInstanceConfig.FixedInstanceCount != nil {
				micMap["fixed_instance_count"] = rc.ManualInstanceConfig.FixedInstanceCount
			}
			micList = append(micList, micMap)
			rcMap["manual_instance_config"] = micList
		}

		if rc.Concurrency != nil {
			rcMap["concurrency"] = rc.Concurrency
		}

		rcList = append(rcList, rcMap)
		_ = d.Set("resource_config", rcList)
	}

	if respData.RequestPaths != nil {
		_ = d.Set("request_paths", helper.StringsInterfaces(respData.RequestPaths))
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.InferenceURL != nil {
		_ = d.Set("inference_url", respData.InferenceURL)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}

	return nil
}

func resourceTencentCloudTeoInferenceServiceV1Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_service_v1.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		serviceId string
		zoneId    string
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("resource id is broken, id is %s", d.Id())
	}
	zoneId = idSplit[0]
	serviceId = idSplit[1]

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	if zoneId == "" {
		return fmt.Errorf("zone_id is required for update teo inference_service.")
	}

	immutableArgs := []string{"zone_id", "name"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChangesExcept("operation") {
		request := teo.NewModifyInferenceServiceRequest()
		request.ZoneId = &zoneId
		request.ServiceId = &serviceId

		if v, ok := d.GetOkExists("listen_port"); ok {
			request.ListenPort = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("containers"); ok {
			for _, item := range v.([]interface{}) {
				containerMap := item.(map[string]interface{})
				container := teo.InferenceContainerConfigForModify{}

				if v, ok := containerMap["image_type"].(string); ok && v != "" {
					container.ImageType = helper.String(v)
				}

				if v, ok := containerMap["tcr_repository_config"]; ok {
					for _, tcrItem := range v.([]interface{}) {
						tcrMap := tcrItem.(map[string]interface{})
						tcrConfig := teo.InferenceTCRRepositoryConfig{}

						if v, ok := tcrMap["tcr_type"].(string); ok && v != "" {
							tcrConfig.TCRType = helper.String(v)
						}

						if v, ok := tcrMap["image"].(string); ok && v != "" {
							tcrConfig.Image = helper.String(v)
						}

						if v, ok := tcrMap["registry_id"].(string); ok && v != "" {
							tcrConfig.RegistryId = helper.String(v)
						}

						if v, ok := tcrMap["region_name"].(string); ok && v != "" {
							tcrConfig.RegionName = helper.String(v)
						}

						container.TcrRepositoryConfig = &tcrConfig
					}
				}

				if v, ok := containerMap["startup_command"].(string); ok && v != "" {
					container.StartupCommand = helper.String(v)
				}

				if v, ok := containerMap["environment_variables"]; ok {
					for _, envItem := range v.([]interface{}) {
						envMap := envItem.(map[string]interface{})
						envVar := teo.InferenceEnvironmentVariable{}

						if v, ok := envMap["key"].(string); ok && v != "" {
							envVar.Key = helper.String(v)
						}

						if v, ok := envMap["value"].(string); ok && v != "" {
							envVar.Value = helper.String(v)
						}

						container.EnvironmentVariables = append(container.EnvironmentVariables, &envVar)
					}
				}

				request.Containers = append(request.Containers, &container)
			}
		}

		if v, ok := d.GetOk("resource_config"); ok {
			for _, item := range v.([]interface{}) {
				rcMap := item.(map[string]interface{})
				resConfig := teo.InferenceResourceConfigForModify{}

				if v, ok := rcMap["scaling_mode"].(string); ok && v != "" {
					resConfig.ScalingMode = helper.String(v)
				}

				if v, ok := rcMap["auto_scaling_config"]; ok {
					for _, ascItem := range v.([]interface{}) {
						ascMap := ascItem.(map[string]interface{})
						autoScalingConfig := teo.InferenceAutoScalingConfig{}

						if v, ok := ascMap["min_instance_count"].(int); ok {
							autoScalingConfig.MinInstanceCount = helper.IntInt64(v)
						}

						if v, ok := ascMap["scaling_policies"]; ok {
							for _, spItem := range v.([]interface{}) {
								spMap := spItem.(map[string]interface{})
								scalingPolicy := teo.InferenceScalingPolicy{}

								if v, ok := spMap["policy_name"].(string); ok && v != "" {
									scalingPolicy.PolicyName = helper.String(v)
								}

								if v, ok := spMap["policy_type"].(string); ok && v != "" {
									scalingPolicy.PolicyType = helper.String(v)
								}

								if v, ok := spMap["scheduled_scaling_policy"]; ok {
									for _, sspItem := range v.([]interface{}) {
										sspMap := sspItem.(map[string]interface{})
										scheduledScalingPolicy := teo.InferenceScheduledScalingPolicy{}

										if v, ok := sspMap["scheduled_actions"]; ok {
											for _, saItem := range v.([]interface{}) {
												saMap := saItem.(map[string]interface{})
												scheduledAction := teo.InferenceScheduledScalingAction{}

												if v, ok := saMap["cron_expression"].(string); ok && v != "" {
													scheduledAction.CronExpression = helper.String(v)
												}

												if v, ok := saMap["min_instance_count"].(int); ok {
													scheduledAction.MinInstanceCount = helper.IntInt64(v)
												}

												scheduledScalingPolicy.ScheduledActions = append(scheduledScalingPolicy.ScheduledActions, &scheduledAction)
											}
										}

										if v, ok := sspMap["effective_range"]; ok {
											for _, erItem := range v.([]interface{}) {
												erMap := erItem.(map[string]interface{})
												effectiveRange := teo.InferenceScheduledScalingEffectiveRange{}

												if v, ok := erMap["effective_type"].(string); ok && v != "" {
													effectiveRange.EffectiveType = helper.String(v)
												}

												if v, ok := erMap["start_date"].(string); ok && v != "" {
													effectiveRange.StartDate = helper.String(v)
												}

												if v, ok := erMap["end_date"].(string); ok && v != "" {
													effectiveRange.EndDate = helper.String(v)
												}

												scheduledScalingPolicy.EffectiveRange = &effectiveRange
											}
										}

										if v, ok := sspMap["time_zone"].(string); ok && v != "" {
											scheduledScalingPolicy.TimeZone = helper.String(v)
										}

										scalingPolicy.ScheduledScalingPolicy = &scheduledScalingPolicy
									}
								}

								autoScalingConfig.ScalingPolicies = append(autoScalingConfig.ScalingPolicies, &scalingPolicy)
							}
						}

						resConfig.AutoScalingConfig = &autoScalingConfig
					}
				}

				if v, ok := rcMap["manual_instance_config"]; ok {
					for _, micItem := range v.([]interface{}) {
						micMap := micItem.(map[string]interface{})
						manualInstanceConfig := teo.InferenceManualInstanceConfig{}

						if v, ok := micMap["fixed_instance_count"].(int); ok {
							manualInstanceConfig.FixedInstanceCount = helper.IntInt64(v)
						}

						resConfig.ManualInstanceConfig = &manualInstanceConfig
					}
				}

				if v, ok := rcMap["concurrency"].(int); ok {
					resConfig.Concurrency = helper.IntInt64(v)
				}

				request.ResourceConfig = &resConfig
			}
		}

		if v, ok := d.GetOk("request_paths"); ok {
			pathsSet := v.(*schema.Set)
			for _, path := range pathsSet.List() {
				request.RequestPaths = append(request.RequestPaths, helper.String(path.(string)))
			}
		}

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyInferenceServiceWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update teo inference_service failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	if d.HasChange("operation") {
		_, newOperation := d.GetChange("operation")
		if newOperation != nil && newOperation.(string) != "" {
			operateRequest := teo.NewOperateInferenceServiceRequest()
			operateRequest.ZoneId = &zoneId
			operateRequest.ServiceId = &serviceId
			operateRequest.Operation = helper.String(newOperation.(string))

			operateErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().OperateInferenceServiceWithContext(ctx, operateRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, operateRequest.GetAction(), operateRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if operateErr != nil {
				log.Printf("[CRITAL]%s operate teo inference_service failed, reason:%+v", logId, operateErr)
				return operateErr
			}
		}
	}

	return resourceTencentCloudTeoInferenceServiceV1Read(d, meta)
}

func resourceTencentCloudTeoInferenceServiceV1Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_service_v1.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		serviceId string
		zoneId    string
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("resource id is broken, id is %s", d.Id())
	}
	zoneId = idSplit[0]
	serviceId = idSplit[1]

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	if zoneId == "" {
		return fmt.Errorf("zone_id is required for delete teo inference_service.")
	}

	request := teo.NewOperateInferenceServiceRequest()
	request.ZoneId = &zoneId
	request.ServiceId = &serviceId
	request.Operation = helper.String("Delete")

	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().OperateInferenceServiceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo inference_service failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
