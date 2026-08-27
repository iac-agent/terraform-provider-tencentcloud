package teo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
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
				Description: "Inference service name. The length is limited to 30 characters. Only lowercase letters, numbers, and hyphens are supported. It must start with a letter and end with a number or letter. Duplicates are not supported.",
			},

			"listen_port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Port that the model service needs to listen to. Only integers between 1-65535 are supported.",
			},

			"containers": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Container configuration of the inference service. Currently only supports setting 1 container.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Image type. Valid values: `TCR` (Tencent Cloud Container Registry image).",
						},
						"tcr_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "TCR image repository configuration. Required when ImageType is TCR.",
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
										Description: "Image repository instance ID. Required when TCRType is Enterprise.",
									},
									"region_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Region name.",
									},
								},
							},
						},
						"startup_command": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Command executed when the container starts. Defaults to the image's Entrypoint/CMD if not specified. Up to 1024 characters.",
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
										Description: "Variable name. Only uppercase and lowercase letters, numbers, and underscores are allowed. It must start with a letter or underscore. Length limit is 64 characters.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Variable value. Supports any visible characters such as letters, numbers, symbols, etc. Length limit is 2048 characters.",
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
				Description: "Resource configuration of the inference service.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scaling_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Scaling mode. Valid values: `Auto` (auto-scale based on request volume), `Manual` (manually set fixed instance count).",
						},
						"hardware_spec": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Hardware specification.",
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
										Optional:    true,
										Description: "Minimum number of instances. Will not take effect when scaling policy is configured and within its validity period.",
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
													Description: "Policy name. Length limit is 1-30 characters. Must be unique within the same service.",
												},
												"policy_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Policy type. Cannot be modified after creation. Valid values: `ScheduledScaling`.",
												},
												"scheduled_scaling_policy": {
													Type:        schema.TypeList,
													Required:    true,
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
																			Description: "Cron expression describing the trigger time. Uses 5-field standard Cron format: minute hour day month weekday.",
																		},
																		"min_instance_count": {
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Minimum number of instances to adjust to after hitting this action.",
																		},
																	},
																},
															},
															"effective_range": {
																Type:        schema.TypeList,
																Required:    true,
																MaxItems:    1,
																Description: "Effective range of the scheduled scaling policy.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"effective_type": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Effective type. Valid values: `LongTerm` (long-term), `Custom` (custom start/end dates).",
																		},
																		"start_date": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Effective start date. Required when EffectiveType is Custom.",
																		},
																		"end_date": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Effective end date. Required when EffectiveType is Custom and must not be earlier than StartDate.",
																		},
																	},
																},
															},
															"time_zone": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Timezone, using IANA timezone identifiers. Defaults to UTC.",
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
							Description: "Concurrency per instance. Default value is 1.",
						},
					},
				},
			},

			"affinity_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Inference service affinity configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Inference service affinity switch. Valid values: `On`, `Off`.",
						},
						"affinity_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Inference service affinity mode. Valid values: `SessionId`. Default value: `SessionId`.",
						},
						"source": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The location where the session ID parameter is passed. Valid values: `Header`. Default value: `Header`.",
						},
						"header_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Request header name for passing the session ID. Required when Source is Header. Length limit: 1-64 characters. Only letters, numbers, and hyphens are supported. Default value: `EO-Infer-Session-Id`.",
						},
					},
				},
			},

			"request_paths": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Request path list of the inference service. Up to 20 paths.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description. Length limit is 60 characters.",
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

			"scaling_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Scaling status. Valid values: `Normal`, `ScalingOut`, `ScalingIn`.",
			},

			"current_instance_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Current number of running instances.",
			},

			"inference_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Inference access URL, through which the underlying model can be accessed for inference.",
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
	}
}

func resourceTencentCloudTeoInferenceServiceV1Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_service_v1.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = teov20220901.NewCreateInferenceServiceRequest()
		response  = teov20220901.NewCreateInferenceServiceResponse()
		zoneId    string
		serviceId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("listen_port"); ok {
		request.ListenPort = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("containers"); ok {
		request.Containers = buildInferenceContainerConfigs(v.([]interface{}))
	}

	if v, ok := d.GetOk("resource_config"); ok {
		request.ResourceConfig = buildInferenceResourceConfig(v.([]interface{}))
	}

	if v, ok := d.GetOk("affinity_config"); ok {
		request.AffinityConfig = buildInferenceAffinityConfig(v.([]interface{}))
	}

	if v, ok := d.GetOk("request_paths"); ok {
		request.RequestPaths = helper.InterfacesStringsPoint(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
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

	if response.Response.ServiceId == nil || *response.Response.ServiceId == "" {
		return fmt.Errorf("ServiceId is nil.")
	}

	serviceId = *response.Response.ServiceId
	d.SetId(strings.Join([]string{zoneId, serviceId}, tccommon.FILED_SP))

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
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	serviceId := idSplit[1]

	respData, err := service.DescribeTeoInferenceServiceById(ctx, zoneId, serviceId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_teo_inference_service_v1` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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
		_ = d.Set("listen_port", int(*respData.ListenPort))
	}

	if respData.Containers != nil {
		_ = d.Set("containers", flattenInferenceContainerConfigs(respData.Containers))
	}

	if respData.ResourceConfig != nil {
		_ = d.Set("resource_config", flattenInferenceResourceConfig(respData.ResourceConfig))
	}

	// InferenceService struct does not contain AffinityConfig directly,
	// so we preserve the existing state value.
	if respData.RequestPaths != nil {
		_ = d.Set("request_paths", helper.StringsInterfaces(respData.RequestPaths))
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.ScalingStatus != nil {
		_ = d.Set("scaling_status", respData.ScalingStatus)
	}

	if respData.CurrentInstanceCount != nil {
		_ = d.Set("current_instance_count", int(*respData.CurrentInstanceCount))
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
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	serviceId := idSplit[1]

	needChange := false
	mutableArgs := []string{"listen_port", "containers", "resource_config", "affinity_config", "request_paths", "description"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := teov20220901.NewModifyInferenceServiceRequest()
		request.ZoneId = &zoneId
		request.ServiceId = &serviceId

		if v, ok := d.GetOk("listen_port"); ok {
			request.ListenPort = helper.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("containers"); ok {
			request.Containers = buildInferenceContainerConfigsForModify(v.([]interface{}))
		}

		if v, ok := d.GetOk("resource_config"); ok {
			request.ResourceConfig = buildInferenceResourceConfigForModify(v.([]interface{}))
		}

		if v, ok := d.GetOk("affinity_config"); ok {
			request.AffinityConfig = buildInferenceAffinityConfig(v.([]interface{}))
		}

		if v, ok := d.GetOk("request_paths"); ok {
			request.RequestPaths = helper.InterfacesStringsPoint(v.(*schema.Set).List())
		}

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
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

	return resourceTencentCloudTeoInferenceServiceV1Read(d, meta)
}

func resourceTencentCloudTeoInferenceServiceV1Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_inference_service_v1.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewOperateInferenceServiceRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	serviceId := idSplit[1]

	request.ZoneId = &zoneId
	request.ServiceId = &serviceId
	request.Operation = helper.String("Delete")

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
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

// buildInferenceContainerConfigs builds InferenceContainerConfig from schema data
func buildInferenceContainerConfigs(containerList []interface{}) []*teov20220901.InferenceContainerConfig {
	if len(containerList) == 0 {
		return nil
	}

	result := make([]*teov20220901.InferenceContainerConfig, 0, len(containerList))
	for _, item := range containerList {
		v := item.(map[string]interface{})
		config := &teov20220901.InferenceContainerConfig{}

		if imageType, ok := v["image_type"]; ok && imageType.(string) != "" {
			config.ImageType = helper.String(imageType.(string))
		}

		if tcrList, ok := v["tcr_config"]; ok {
			config.TcrRepositoryConfig = buildInferenceTCRRepositoryConfig(tcrList.([]interface{}))
		}

		if startupCmd, ok := v["startup_command"]; ok && startupCmd.(string) != "" {
			config.StartupCommand = helper.String(startupCmd.(string))
		}

		if envVars, ok := v["environment_variables"]; ok {
			config.EnvironmentVariables = buildInferenceEnvironmentVariables(envVars.([]interface{}))
		}

		result = append(result, config)
	}

	return result
}

// buildInferenceContainerConfigsForModify builds InferenceContainerConfigForModify from schema data
func buildInferenceContainerConfigsForModify(containerList []interface{}) []*teov20220901.InferenceContainerConfigForModify {
	if len(containerList) == 0 {
		return nil
	}

	result := make([]*teov20220901.InferenceContainerConfigForModify, 0, len(containerList))
	for _, item := range containerList {
		v := item.(map[string]interface{})
		config := &teov20220901.InferenceContainerConfigForModify{}

		if imageType, ok := v["image_type"]; ok && imageType.(string) != "" {
			config.ImageType = helper.String(imageType.(string))
		}

		if tcrList, ok := v["tcr_config"]; ok {
			config.TcrRepositoryConfig = buildInferenceTCRRepositoryConfig(tcrList.([]interface{}))
		}

		if startupCmd, ok := v["startup_command"]; ok && startupCmd.(string) != "" {
			config.StartupCommand = helper.String(startupCmd.(string))
		}

		if envVars, ok := v["environment_variables"]; ok {
			config.EnvironmentVariables = buildInferenceEnvironmentVariables(envVars.([]interface{}))
		}

		result = append(result, config)
	}

	return result
}

// buildInferenceTCRRepositoryConfig builds InferenceTCRRepositoryConfig from schema data
func buildInferenceTCRRepositoryConfig(tcrList []interface{}) *teov20220901.InferenceTCRRepositoryConfig {
	if len(tcrList) == 0 {
		return nil
	}

	v := tcrList[0].(map[string]interface{})
	config := &teov20220901.InferenceTCRRepositoryConfig{}

	if tcrType, ok := v["tcr_type"]; ok && tcrType.(string) != "" {
		config.TCRType = helper.String(tcrType.(string))
	}

	if image, ok := v["image"]; ok && image.(string) != "" {
		config.Image = helper.String(image.(string))
	}

	if registryId, ok := v["registry_id"]; ok && registryId.(string) != "" {
		config.RegistryId = helper.String(registryId.(string))
	}

	if regionName, ok := v["region_name"]; ok && regionName.(string) != "" {
		config.RegionName = helper.String(regionName.(string))
	}

	return config
}

// buildInferenceEnvironmentVariables builds InferenceEnvironmentVariable list from schema data
func buildInferenceEnvironmentVariables(envVarList []interface{}) []*teov20220901.InferenceEnvironmentVariable {
	if len(envVarList) == 0 {
		return nil
	}

	result := make([]*teov20220901.InferenceEnvironmentVariable, 0, len(envVarList))
	for _, item := range envVarList {
		v := item.(map[string]interface{})
		envVar := &teov20220901.InferenceEnvironmentVariable{}

		if key, ok := v["key"]; ok && key.(string) != "" {
			envVar.Key = helper.String(key.(string))
		}

		if value, ok := v["value"]; ok && value.(string) != "" {
			envVar.Value = helper.String(value.(string))
		}

		result = append(result, envVar)
	}

	return result
}

// buildInferenceResourceConfig builds InferenceResourceConfig from schema data
func buildInferenceResourceConfig(rcList []interface{}) *teov20220901.InferenceResourceConfig {
	if len(rcList) == 0 {
		return nil
	}

	v := rcList[0].(map[string]interface{})
	config := &teov20220901.InferenceResourceConfig{}

	if scalingMode, ok := v["scaling_mode"]; ok && scalingMode.(string) != "" {
		config.ScalingMode = helper.String(scalingMode.(string))
	}

	if hardwareSpec, ok := v["hardware_spec"]; ok && hardwareSpec.(string) != "" {
		config.HardwareSpec = helper.String(hardwareSpec.(string))
	}

	if autoScalingList, ok := v["auto_scaling_config"]; ok {
		config.AutoScalingConfig = buildInferenceAutoScalingConfig(autoScalingList.([]interface{}))
	}

	if manualInstanceList, ok := v["manual_instance_config"]; ok {
		config.ManualInstanceConfig = buildInferenceManualInstanceConfig(manualInstanceList.([]interface{}))
	}

	if concurrency, ok := v["concurrency"]; ok {
		config.Concurrency = helper.Int64(int64(concurrency.(int)))
	}

	return config
}

// buildInferenceResourceConfigForModify builds InferenceResourceConfigForModify from schema data
func buildInferenceResourceConfigForModify(rcList []interface{}) *teov20220901.InferenceResourceConfigForModify {
	if len(rcList) == 0 {
		return nil
	}

	v := rcList[0].(map[string]interface{})
	config := &teov20220901.InferenceResourceConfigForModify{}

	if scalingMode, ok := v["scaling_mode"]; ok && scalingMode.(string) != "" {
		config.ScalingMode = helper.String(scalingMode.(string))
	}

	if autoScalingList, ok := v["auto_scaling_config"]; ok {
		config.AutoScalingConfig = buildInferenceAutoScalingConfig(autoScalingList.([]interface{}))
	}

	if manualInstanceList, ok := v["manual_instance_config"]; ok {
		config.ManualInstanceConfig = buildInferenceManualInstanceConfig(manualInstanceList.([]interface{}))
	}

	if concurrency, ok := v["concurrency"]; ok {
		config.Concurrency = helper.Int64(int64(concurrency.(int)))
	}

	return config
}

// buildInferenceAutoScalingConfig builds InferenceAutoScalingConfig from schema data
func buildInferenceAutoScalingConfig(ascList []interface{}) *teov20220901.InferenceAutoScalingConfig {
	if len(ascList) == 0 {
		return nil
	}

	v := ascList[0].(map[string]interface{})
	config := &teov20220901.InferenceAutoScalingConfig{}

	if minCount, ok := v["min_instance_count"]; ok {
		config.MinInstanceCount = helper.Int64(int64(minCount.(int)))
	}

	if policies, ok := v["scaling_policies"]; ok {
		config.ScalingPolicies = buildInferenceScalingPolicies(policies.([]interface{}))
	}

	return config
}

// buildInferenceScalingPolicies builds InferenceScalingPolicy list from schema data
func buildInferenceScalingPolicies(policyList []interface{}) []*teov20220901.InferenceScalingPolicy {
	if len(policyList) == 0 {
		return nil
	}

	result := make([]*teov20220901.InferenceScalingPolicy, 0, len(policyList))
	for _, item := range policyList {
		v := item.(map[string]interface{})
		policy := &teov20220901.InferenceScalingPolicy{}

		if policyName, ok := v["policy_name"]; ok && policyName.(string) != "" {
			policy.PolicyName = helper.String(policyName.(string))
		}

		if policyType, ok := v["policy_type"]; ok && policyType.(string) != "" {
			policy.PolicyType = helper.String(policyType.(string))
		}

		if scheduledList, ok := v["scheduled_scaling_policy"]; ok {
			policy.ScheduledScalingPolicy = buildInferenceScheduledScalingPolicy(scheduledList.([]interface{}))
		}

		result = append(result, policy)
	}

	return result
}

// buildInferenceScheduledScalingPolicy builds InferenceScheduledScalingPolicy from schema data
func buildInferenceScheduledScalingPolicy(sspList []interface{}) *teov20220901.InferenceScheduledScalingPolicy {
	if len(sspList) == 0 {
		return nil
	}

	v := sspList[0].(map[string]interface{})
	policy := &teov20220901.InferenceScheduledScalingPolicy{}

	if actions, ok := v["scheduled_actions"]; ok {
		policy.ScheduledActions = buildInferenceScheduledScalingActions(actions.([]interface{}))
	}

	if effectiveRange, ok := v["effective_range"]; ok {
		policy.EffectiveRange = buildInferenceScheduledScalingEffectiveRange(effectiveRange.([]interface{}))
	}

	if timeZone, ok := v["time_zone"]; ok && timeZone.(string) != "" {
		policy.TimeZone = helper.String(timeZone.(string))
	}

	return policy
}

// buildInferenceScheduledScalingActions builds InferenceScheduledScalingAction list from schema data
func buildInferenceScheduledScalingActions(actionList []interface{}) []*teov20220901.InferenceScheduledScalingAction {
	if len(actionList) == 0 {
		return nil
	}

	result := make([]*teov20220901.InferenceScheduledScalingAction, 0, len(actionList))
	for _, item := range actionList {
		v := item.(map[string]interface{})
		action := &teov20220901.InferenceScheduledScalingAction{}

		if cronExpr, ok := v["cron_expression"]; ok && cronExpr.(string) != "" {
			action.CronExpression = helper.String(cronExpr.(string))
		}

		if minCount, ok := v["min_instance_count"]; ok {
			action.MinInstanceCount = helper.Int64(int64(minCount.(int)))
		}

		result = append(result, action)
	}

	return result
}

// buildInferenceScheduledScalingEffectiveRange builds InferenceScheduledScalingEffectiveRange from schema data
func buildInferenceScheduledScalingEffectiveRange(erList []interface{}) *teov20220901.InferenceScheduledScalingEffectiveRange {
	if len(erList) == 0 {
		return nil
	}

	v := erList[0].(map[string]interface{})
	r := &teov20220901.InferenceScheduledScalingEffectiveRange{}

	if effectiveType, ok := v["effective_type"]; ok && effectiveType.(string) != "" {
		r.EffectiveType = helper.String(effectiveType.(string))
	}

	if startDate, ok := v["start_date"]; ok && startDate.(string) != "" {
		r.StartDate = helper.String(startDate.(string))
	}

	if endDate, ok := v["end_date"]; ok && endDate.(string) != "" {
		r.EndDate = helper.String(endDate.(string))
	}

	return r
}

// buildInferenceManualInstanceConfig builds InferenceManualInstanceConfig from schema data
func buildInferenceManualInstanceConfig(micList []interface{}) *teov20220901.InferenceManualInstanceConfig {
	if len(micList) == 0 {
		return nil
	}

	v := micList[0].(map[string]interface{})
	config := &teov20220901.InferenceManualInstanceConfig{}

	if fixedCount, ok := v["fixed_instance_count"]; ok {
		config.FixedInstanceCount = helper.Int64(int64(fixedCount.(int)))
	}

	return config
}

// buildInferenceAffinityConfig builds InferenceAffinityConfig from schema data
func buildInferenceAffinityConfig(acList []interface{}) *teov20220901.InferenceAffinityConfig {
	if len(acList) == 0 {
		return nil
	}

	v := acList[0].(map[string]interface{})
	config := &teov20220901.InferenceAffinityConfig{}

	if sw, ok := v["switch"]; ok && sw.(string) != "" {
		config.Switch = helper.String(sw.(string))
	}

	if affinityMode, ok := v["affinity_mode"]; ok && affinityMode.(string) != "" {
		config.AffinityMode = helper.String(affinityMode.(string))
	}

	// Build SessionIdAffinityConfig when source or header_name is set
	sessionIdConfig := &teov20220901.SessionIdAffinityConfig{}
	hasSessionIdConfig := false

	if source, ok := v["source"]; ok && source.(string) != "" {
		sessionIdConfig.Source = helper.String(source.(string))
		hasSessionIdConfig = true
	}

	if headerName, ok := v["header_name"]; ok && headerName.(string) != "" {
		sessionIdConfig.HeaderName = helper.String(headerName.(string))
		hasSessionIdConfig = true
	}

	if hasSessionIdConfig {
		config.SessionIdAffinityConfig = sessionIdConfig
	}

	return config
}

// flattenInferenceContainerConfigs converts InferenceContainerConfig list to map list
func flattenInferenceContainerConfigs(containers []*teov20220901.InferenceContainerConfig) []map[string]interface{} {
	if len(containers) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(containers))
	for _, c := range containers {
		item := make(map[string]interface{})

		if c.ImageType != nil {
			item["image_type"] = *c.ImageType
		}

		if c.TcrRepositoryConfig != nil {
			item["tcr_config"] = flattenInferenceTCRRepositoryConfig(c.TcrRepositoryConfig)
		}

		if c.StartupCommand != nil {
			item["startup_command"] = *c.StartupCommand
		}

		if c.EnvironmentVariables != nil {
			item["environment_variables"] = flattenInferenceEnvironmentVariables(c.EnvironmentVariables)
		}

		result = append(result, item)
	}

	return result
}

// flattenInferenceTCRRepositoryConfig converts InferenceTCRRepositoryConfig to map list
func flattenInferenceTCRRepositoryConfig(config *teov20220901.InferenceTCRRepositoryConfig) []map[string]interface{} {
	if config == nil {
		return nil
	}

	item := make(map[string]interface{})

	if config.TCRType != nil {
		item["tcr_type"] = *config.TCRType
	}

	if config.Image != nil {
		item["image"] = *config.Image
	}

	if config.RegistryId != nil {
		item["registry_id"] = *config.RegistryId
	}

	if config.RegionName != nil {
		item["region_name"] = *config.RegionName
	}

	return []map[string]interface{}{item}
}

// flattenInferenceEnvironmentVariables converts InferenceEnvironmentVariable list to map list
func flattenInferenceEnvironmentVariables(envVars []*teov20220901.InferenceEnvironmentVariable) []map[string]interface{} {
	if len(envVars) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(envVars))
	for _, ev := range envVars {
		item := make(map[string]interface{})

		if ev.Key != nil {
			item["key"] = *ev.Key
		}

		if ev.Value != nil {
			item["value"] = *ev.Value
		}

		result = append(result, item)
	}

	return result
}

// flattenInferenceResourceConfig converts InferenceResourceConfig to map list
func flattenInferenceResourceConfig(config *teov20220901.InferenceResourceConfig) []map[string]interface{} {
	if config == nil {
		return nil
	}

	item := make(map[string]interface{})

	if config.ScalingMode != nil {
		item["scaling_mode"] = *config.ScalingMode
	}

	if config.HardwareSpec != nil {
		item["hardware_spec"] = *config.HardwareSpec
	}

	if config.AutoScalingConfig != nil {
		item["auto_scaling_config"] = flattenInferenceAutoScalingConfig(config.AutoScalingConfig)
	}

	if config.ManualInstanceConfig != nil {
		item["manual_instance_config"] = flattenInferenceManualInstanceConfig(config.ManualInstanceConfig)
	}

	if config.Concurrency != nil {
		item["concurrency"] = int(*config.Concurrency)
	}

	return []map[string]interface{}{item}
}

// flattenInferenceAutoScalingConfig converts InferenceAutoScalingConfig to map list
func flattenInferenceAutoScalingConfig(config *teov20220901.InferenceAutoScalingConfig) []map[string]interface{} {
	if config == nil {
		return nil
	}

	item := make(map[string]interface{})

	if config.MinInstanceCount != nil {
		item["min_instance_count"] = int(*config.MinInstanceCount)
	}

	if config.ScalingPolicies != nil {
		item["scaling_policies"] = flattenInferenceScalingPolicies(config.ScalingPolicies)
	}

	return []map[string]interface{}{item}
}

// flattenInferenceScalingPolicies converts InferenceScalingPolicy list to map list
func flattenInferenceScalingPolicies(policies []*teov20220901.InferenceScalingPolicy) []map[string]interface{} {
	if len(policies) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(policies))
	for _, p := range policies {
		item := make(map[string]interface{})

		if p.PolicyName != nil {
			item["policy_name"] = *p.PolicyName
		}

		if p.PolicyType != nil {
			item["policy_type"] = *p.PolicyType
		}

		if p.ScheduledScalingPolicy != nil {
			item["scheduled_scaling_policy"] = flattenInferenceScheduledScalingPolicy(p.ScheduledScalingPolicy)
		}

		result = append(result, item)
	}

	return result
}

// flattenInferenceScheduledScalingPolicy converts InferenceScheduledScalingPolicy to map list
func flattenInferenceScheduledScalingPolicy(policy *teov20220901.InferenceScheduledScalingPolicy) []map[string]interface{} {
	if policy == nil {
		return nil
	}

	item := make(map[string]interface{})

	if policy.ScheduledActions != nil {
		item["scheduled_actions"] = flattenInferenceScheduledScalingActions(policy.ScheduledActions)
	}

	if policy.EffectiveRange != nil {
		item["effective_range"] = flattenInferenceScheduledScalingEffectiveRange(policy.EffectiveRange)
	}

	if policy.TimeZone != nil {
		item["time_zone"] = *policy.TimeZone
	}

	return []map[string]interface{}{item}
}

// flattenInferenceScheduledScalingActions converts InferenceScheduledScalingAction list to map list
func flattenInferenceScheduledScalingActions(actions []*teov20220901.InferenceScheduledScalingAction) []map[string]interface{} {
	if len(actions) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(actions))
	for _, a := range actions {
		item := make(map[string]interface{})

		if a.CronExpression != nil {
			item["cron_expression"] = *a.CronExpression
		}

		if a.MinInstanceCount != nil {
			item["min_instance_count"] = int(*a.MinInstanceCount)
		}

		result = append(result, item)
	}

	return result
}

// flattenInferenceScheduledScalingEffectiveRange converts InferenceScheduledScalingEffectiveRange to map list
func flattenInferenceScheduledScalingEffectiveRange(r *teov20220901.InferenceScheduledScalingEffectiveRange) []map[string]interface{} {
	if r == nil {
		return nil
	}

	item := make(map[string]interface{})

	if r.EffectiveType != nil {
		item["effective_type"] = *r.EffectiveType
	}

	if r.StartDate != nil {
		item["start_date"] = *r.StartDate
	}

	if r.EndDate != nil {
		item["end_date"] = *r.EndDate
	}

	return []map[string]interface{}{item}
}

// flattenInferenceManualInstanceConfig converts InferenceManualInstanceConfig to map list
func flattenInferenceManualInstanceConfig(config *teov20220901.InferenceManualInstanceConfig) []map[string]interface{} {
	if config == nil {
		return nil
	}

	item := make(map[string]interface{})

	if config.FixedInstanceCount != nil {
		item["fixed_instance_count"] = int(*config.FixedInstanceCount)
	}

	return []map[string]interface{}{item}
}
