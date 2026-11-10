package teo

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var _ datasource.DataSourceWithConfigure = &TeoInferenceServiceDeploymentRecords111001DataSource{}

// NewTeoInferenceServiceDeploymentRecords111001DataSource is the factory
// referenced by tencentcloud/framework/registry.go to register this data source.
func NewTeoInferenceServiceDeploymentRecords111001DataSource() datasource.DataSource {
	return &TeoInferenceServiceDeploymentRecords111001DataSource{}
}

// TeoInferenceServiceDeploymentRecords111001DataSource implements
// datasource.DataSource for
// tencentcloud_teo_inference_service_deployment_records_1110_01.
type TeoInferenceServiceDeploymentRecords111001DataSource struct {
	fw.DataSourceWithConfigure
}

// TeoInferenceServiceDeploymentRecords111001DataSourceModel maps the schema
// attributes.
type TeoInferenceServiceDeploymentRecords111001DataSourceModel struct {
	ZoneId    types.String `tfsdk:"zone_id"`
	ServiceId types.String `tfsdk:"service_id"`
	SortBy    types.String `tfsdk:"sort_by"`
	SortOrder types.String `tfsdk:"sort_order"`
	RecordSet types.List   `tfsdk:"record_set"`
}

// teoDeploymentRecordItemModel is the element type of record_set.
type teoDeploymentRecordItemModel struct {
	RecordId               types.String `tfsdk:"record_id"`
	Operation              types.String `tfsdk:"operation"`
	Status                 types.String `tfsdk:"status"`
	Duration               types.Int64  `tfsdk:"duration"`
	InferenceServiceConfig types.Object `tfsdk:"inference_service_config"`
	CreateTime             types.String `tfsdk:"create_time"`
	ActiveStatus           types.String `tfsdk:"active_status"`
}

// teoInferenceServiceConfigModel maps the inference_service_config attribute.
type teoInferenceServiceConfigModel struct {
	ListenPort     types.Int64  `tfsdk:"listen_port"`
	RequestPaths   types.Set    `tfsdk:"request_paths"`
	Containers     types.List   `tfsdk:"containers"`
	ResourceConfig types.Object `tfsdk:"resource_config"`
	AffinityConfig types.Object `tfsdk:"affinity_config"`
}

// teoContainerItemModel is the element type of containers.
type teoContainerItemModel struct {
	ImageType            types.String `tfsdk:"image_type"`
	TcrRepositoryConfig  types.Object `tfsdk:"tcr_repository_config"`
	StartupCommand       types.String `tfsdk:"startup_command"`
	EnvironmentVariables types.List   `tfsdk:"environment_variables"`
}

// teoTCRRepositoryConfigModel maps the tcr_repository_config attribute.
type teoTCRRepositoryConfigModel struct {
	TCRType    types.String `tfsdk:"tcr_type"`
	Image      types.String `tfsdk:"image"`
	RegistryId types.String `tfsdk:"registry_id"`
	RegionName types.String `tfsdk:"region_name"`
}

// teoEnvironmentVariableModel is the element type of environment_variables.
type teoEnvironmentVariableModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// teoResourceConfigModel maps the resource_config attribute.
type teoResourceConfigModel struct {
	ScalingMode          types.String `tfsdk:"scaling_mode"`
	HardwareSpec         types.String `tfsdk:"hardware_spec"`
	HardwareSpecId       types.String `tfsdk:"hardware_spec_id"`
	HardwareConfig       types.Object `tfsdk:"hardware_config"`
	AutoScalingConfig    types.Object `tfsdk:"auto_scaling_config"`
	ManualInstanceConfig types.Object `tfsdk:"manual_instance_config"`
	Concurrency          types.Int64  `tfsdk:"concurrency"`
}

// teoHardwareConfigModel maps the hardware_config attribute.
type teoHardwareConfigModel struct {
	GPUNum   types.Float64 `tfsdk:"gpu_num"`
	CPUNum   types.Float64 `tfsdk:"cpu_num"`
	MemSize  types.Int64   `tfsdk:"mem_size"`
	DiskSize types.Int64   `tfsdk:"disk_size"`
}

// teoAutoScalingConfigModel maps the auto_scaling_config attribute.
type teoAutoScalingConfigModel struct {
	MinInstanceCount types.Int64 `tfsdk:"min_instance_count"`
	ScalingPolicies  types.List  `tfsdk:"scaling_policies"`
}

// teoScalingPolicyModel is the element type of scaling_policies.
type teoScalingPolicyModel struct {
	PolicyName             types.String `tfsdk:"policy_name"`
	PolicyType             types.String `tfsdk:"policy_type"`
	ScheduledScalingPolicy types.Object `tfsdk:"scheduled_scaling_policy"`
}

// teoScheduledScalingPolicyModel maps the scheduled_scaling_policy attribute.
type teoScheduledScalingPolicyModel struct {
	ScheduledActions types.List   `tfsdk:"scheduled_actions"`
	EffectiveRange   types.Object `tfsdk:"effective_range"`
	TimeZone         types.String `tfsdk:"time_zone"`
}

// teoScheduledScalingActionModel is the element type of scheduled_actions.
type teoScheduledScalingActionModel struct {
	CronExpression   types.String `tfsdk:"cron_expression"`
	MinInstanceCount types.Int64  `tfsdk:"min_instance_count"`
}

// teoScheduledScalingEffectiveRangeModel maps the effective_range attribute.
type teoScheduledScalingEffectiveRangeModel struct {
	EffectiveType types.String `tfsdk:"effective_type"`
	StartDate     types.String `tfsdk:"start_date"`
	EndDate       types.String `tfsdk:"end_date"`
}

// teoManualInstanceConfigModel maps the manual_instance_config attribute.
type teoManualInstanceConfigModel struct {
	FixedInstanceCount types.Int64 `tfsdk:"fixed_instance_count"`
}

// teoAffinityConfigModel maps the affinity_config attribute.
type teoAffinityConfigModel struct {
	Switch                  types.String `tfsdk:"switch"`
	AffinityMode            types.String `tfsdk:"affinity_mode"`
	SessionIdAffinityConfig types.Object `tfsdk:"session_id_affinity_config"`
}

// teoSessionIdAffinityConfigModel maps the session_id_affinity_config
// attribute.
type teoSessionIdAffinityConfigModel struct {
	Source     types.String `tfsdk:"source"`
	HeaderName types.String `tfsdk:"header_name"`
}

func (d *TeoInferenceServiceDeploymentRecords111001DataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_inference_service_deployment_records_1110_01"
}

func (d *TeoInferenceServiceDeploymentRecords111001DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the TEO inference service deployment record list.",
		Attributes: map[string]schema.Attribute{
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Zone ID.",
			},
			"service_id": schema.StringAttribute{
				Required:    true,
				Description: "Inference service ID.",
			},
			"sort_by": schema.StringAttribute{
				Optional: true,
				Description: "Sort field. Valid values: `create-time` (deployment creation time). " +
					"Default value: `create-time`.",
			},
			"sort_order": schema.StringAttribute{
				Optional: true,
				Description: "Sort order. Valid values: `asc` (ascending), `desc` (descending). " +
					"Default value: `desc`.",
			},
			"record_set": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of inference service deployment records.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"record_id": schema.StringAttribute{
							Computed:    true,
							Description: "Deployment record ID.",
						},
						"operation": schema.StringAttribute{
							Computed: true,
							Description: "Deployment operation type. Valid values: `create`, `update`, " +
								"`resume`, `stop`.",
						},
						"status": schema.StringAttribute{
							Computed: true,
							Description: "Deployment status. Valid values: `processing`, `succeeded`, " +
								"`failed`.",
						},
						"duration": schema.Int64Attribute{
							Computed:    true,
							Description: "Deployment duration, in seconds.",
						},
						"inference_service_config": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Configuration of this inference service deployment.",
							Attributes: map[string]schema.Attribute{
								"listen_port": schema.Int64Attribute{
									Computed:    true,
									Description: "Port that the model service needs to listen on.",
								},
								"request_paths": schema.SetAttribute{
									Computed:    true,
									ElementType: types.StringType,
									Description: "Request path list of the inference service.",
								},
								"containers": schema.ListNestedAttribute{
									Computed:    true,
									Description: "Container configuration of the inference service.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"image_type": schema.StringAttribute{
												Computed:    true,
												Description: "Image type. Valid values: `TCR`.",
											},
											"tcr_repository_config": schema.SingleNestedAttribute{
												Computed:    true,
												Description: "TCR repository information.",
												Attributes: map[string]schema.Attribute{
													"tcr_type": schema.StringAttribute{
														Computed: true,
														Description: "TCR service type. Valid values: " +
															"`Personal`, `Enterprise`.",
													},
													"image": schema.StringAttribute{
														Computed:    true,
														Description: "Image address.",
													},
													"registry_id": schema.StringAttribute{
														Computed:    true,
														Description: "TCR registry instance ID.",
													},
													"region_name": schema.StringAttribute{
														Computed:    true,
														Description: "Region name.",
													},
												},
											},
											"startup_command": schema.StringAttribute{
												Computed: true,
												Description: "Command executed when the container starts. " +
													"When it is left empty, the image Entrypoint/CMD is used.",
											},
											"environment_variables": schema.ListNestedAttribute{
												Computed:    true,
												Description: "Environment variables of the container.",
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"key": schema.StringAttribute{
															Computed:    true,
															Description: "Variable name.",
														},
														"value": schema.StringAttribute{
															Computed:    true,
															Description: "Variable value.",
														},
													},
												},
											},
										},
									},
								},
								"resource_config": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "Resource configuration of the inference service.",
									Attributes: map[string]schema.Attribute{
										"scaling_mode": schema.StringAttribute{
											Computed:    true,
											Description: "Scaling mode. Valid values: `Auto`, `Manual`.",
										},
										"hardware_spec": schema.StringAttribute{
											Computed: true,
											Description: "Hardware specification identifier. Deprecated, " +
												"use `hardware_spec_id` instead.",
										},
										"hardware_spec_id": schema.StringAttribute{
											Computed: true,
											Description: "Hardware specification unique ID, which can be " +
												"obtained from DescribeInferenceHardwareSpecifications.",
										},
										"hardware_config": schema.SingleNestedAttribute{
											Computed:    true,
											Description: "Hardware configuration of the inference service.",
											Attributes: map[string]schema.Attribute{
												"gpu_num": schema.Float64Attribute{
													Computed: true,
													Description: "Number of GPU cards allocated to a single " +
														"instance of the inference service.",
												},
												"cpu_num": schema.Float64Attribute{
													Computed: true,
													Description: "Number of CPU cores allocated to a single " +
														"instance of the inference service.",
												},
												"mem_size": schema.Int64Attribute{
													Computed: true,
													Description: "Memory size allocated to a single instance " +
														"of the inference service, in MB.",
												},
												"disk_size": schema.Int64Attribute{
													Computed: true,
													Description: "Temporary disk size allocated to a single " +
														"instance of the inference service, in MB.",
												},
											},
										},
										"auto_scaling_config": schema.SingleNestedAttribute{
											Computed:    true,
											Description: "Auto scaling configuration of the inference service.",
											Attributes: map[string]schema.Attribute{
												"min_instance_count": schema.Int64Attribute{
													Computed:    true,
													Description: "Minimum instance count.",
												},
												"scaling_policies": schema.ListNestedAttribute{
													Computed:    true,
													Description: "Scaling policy list.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"policy_name": schema.StringAttribute{
																Computed:    true,
																Description: "Policy name.",
															},
															"policy_type": schema.StringAttribute{
																Computed: true,
																Description: "Policy type. Valid values: " +
																	"`ScheduledScaling`.",
															},
															"scheduled_scaling_policy": schema.SingleNestedAttribute{
																Computed:    true,
																Description: "Scheduled scaling configuration.",
																Attributes: map[string]schema.Attribute{
																	"scheduled_actions": schema.ListNestedAttribute{
																		Computed:    true,
																		Description: "Scheduled action list.",
																		NestedObject: schema.NestedAttributeObject{
																			Attributes: map[string]schema.Attribute{
																				"cron_expression": schema.StringAttribute{
																					Computed:    true,
																					Description: "Cron expression.",
																				},
																				"min_instance_count": schema.Int64Attribute{
																					Computed:    true,
																					Description: "Minimum instance count.",
																				},
																			},
																		},
																	},
																	"effective_range": schema.SingleNestedAttribute{
																		Computed:    true,
																		Description: "Effective range.",
																		Attributes: map[string]schema.Attribute{
																			"effective_type": schema.StringAttribute{
																				Computed: true,
																				Description: "Effective type. Valid values: " +
																					"`LongTerm`, `Custom`.",
																			},
																			"start_date": schema.StringAttribute{
																				Computed:    true,
																				Description: "Start date.",
																			},
																			"end_date": schema.StringAttribute{
																				Computed:    true,
																				Description: "End date.",
																			},
																		},
																	},
																	"time_zone": schema.StringAttribute{
																		Computed:    true,
																		Description: "Time zone.",
																	},
																},
															},
														},
													},
												},
											},
										},
										"manual_instance_config": schema.SingleNestedAttribute{
											Computed:    true,
											Description: "Manual instance configuration of the inference service.",
											Attributes: map[string]schema.Attribute{
												"fixed_instance_count": schema.Int64Attribute{
													Computed:    true,
													Description: "Fixed instance count.",
												},
											},
										},
										"concurrency": schema.Int64Attribute{
											Computed:    true,
											Description: "Concurrency of a single instance.",
										},
									},
								},
								"affinity_config": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "Affinity configuration of the inference service.",
									Attributes: map[string]schema.Attribute{
										"switch": schema.StringAttribute{
											Computed:    true,
											Description: "Affinity switch. Valid values: `On`, `Off`.",
										},
										"affinity_mode": schema.StringAttribute{
											Computed:    true,
											Description: "Affinity mode. Valid values: `SessionId`.",
										},
										"session_id_affinity_config": schema.SingleNestedAttribute{
											Computed:    true,
											Description: "Session ID affinity configuration.",
											Attributes: map[string]schema.Attribute{
												"source": schema.StringAttribute{
													Computed: true,
													Description: "Where the session ID is passed. Valid values: " +
														"`Header`.",
												},
												"header_name": schema.StringAttribute{
													Computed:    true,
													Description: "Header name that carries the session ID.",
												},
											},
										},
									},
								},
							},
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "Deployment initiation time in ISO date format.",
						},
						"active_status": schema.StringAttribute{
							Computed: true,
							Description: "Whether this deployment configuration is the currently active one. " +
								"Valid values: `active`, `inactive`.",
						},
					},
				},
			},
		},
	}
}

// teoEnvironmentVariableAttrTypes returns the attr.Type map of an
// environment_variables element.
func teoEnvironmentVariableAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	}
}

// teoTCRRepositoryConfigAttrTypes returns the attr.Type map of the
// tcr_repository_config attribute.
func teoTCRRepositoryConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"tcr_type":    types.StringType,
		"image":       types.StringType,
		"registry_id": types.StringType,
		"region_name": types.StringType,
	}
}

// teoContainerAttrTypes returns the attr.Type map of a containers element.
func teoContainerAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"image_type": types.StringType,
		"tcr_repository_config": types.ObjectType{
			AttrTypes: teoTCRRepositoryConfigAttrTypes(),
		},
		"startup_command": types.StringType,
		"environment_variables": types.ListType{
			ElemType: types.ObjectType{AttrTypes: teoEnvironmentVariableAttrTypes()},
		},
	}
}

// teoHardwareConfigAttrTypes returns the attr.Type map of the hardware_config
// attribute.
func teoHardwareConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"gpu_num":   types.Float64Type,
		"cpu_num":   types.Float64Type,
		"mem_size":  types.Int64Type,
		"disk_size": types.Int64Type,
	}
}

// teoScheduledScalingActionAttrTypes returns the attr.Type map of a
// scheduled_actions element.
func teoScheduledScalingActionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cron_expression":    types.StringType,
		"min_instance_count": types.Int64Type,
	}
}

// teoScheduledScalingEffectiveRangeAttrTypes returns the attr.Type map of the
// effective_range attribute.
func teoScheduledScalingEffectiveRangeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"effective_type": types.StringType,
		"start_date":     types.StringType,
		"end_date":       types.StringType,
	}
}

// teoScheduledScalingPolicyAttrTypes returns the attr.Type map of the
// scheduled_scaling_policy attribute.
func teoScheduledScalingPolicyAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"scheduled_actions": types.ListType{
			ElemType: types.ObjectType{AttrTypes: teoScheduledScalingActionAttrTypes()},
		},
		"effective_range": types.ObjectType{
			AttrTypes: teoScheduledScalingEffectiveRangeAttrTypes(),
		},
		"time_zone": types.StringType,
	}
}

// teoScalingPolicyAttrTypes returns the attr.Type map of a scaling_policies
// element.
func teoScalingPolicyAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"policy_name": types.StringType,
		"policy_type": types.StringType,
		"scheduled_scaling_policy": types.ObjectType{
			AttrTypes: teoScheduledScalingPolicyAttrTypes(),
		},
	}
}

// teoAutoScalingConfigAttrTypes returns the attr.Type map of the
// auto_scaling_config attribute.
func teoAutoScalingConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"min_instance_count": types.Int64Type,
		"scaling_policies": types.ListType{
			ElemType: types.ObjectType{AttrTypes: teoScalingPolicyAttrTypes()},
		},
	}
}

// teoManualInstanceConfigAttrTypes returns the attr.Type map of the
// manual_instance_config attribute.
func teoManualInstanceConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"fixed_instance_count": types.Int64Type,
	}
}

// teoResourceConfigAttrTypes returns the attr.Type map of the resource_config
// attribute.
func teoResourceConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"scaling_mode":        types.StringType,
		"hardware_spec":       types.StringType,
		"hardware_spec_id":    types.StringType,
		"hardware_config":     types.ObjectType{AttrTypes: teoHardwareConfigAttrTypes()},
		"auto_scaling_config": types.ObjectType{AttrTypes: teoAutoScalingConfigAttrTypes()},
		"manual_instance_config": types.ObjectType{
			AttrTypes: teoManualInstanceConfigAttrTypes(),
		},
		"concurrency": types.Int64Type,
	}
}

// teoSessionIdAffinityConfigAttrTypes returns the attr.Type map of the
// session_id_affinity_config attribute.
func teoSessionIdAffinityConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"source":      types.StringType,
		"header_name": types.StringType,
	}
}

// teoAffinityConfigAttrTypes returns the attr.Type map of the affinity_config
// attribute.
func teoAffinityConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"switch":        types.StringType,
		"affinity_mode": types.StringType,
		"session_id_affinity_config": types.ObjectType{
			AttrTypes: teoSessionIdAffinityConfigAttrTypes(),
		},
	}
}

// teoInferenceServiceConfigAttrTypes returns the attr.Type map of the
// inference_service_config attribute.
func teoInferenceServiceConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"listen_port":   types.Int64Type,
		"request_paths": types.SetType{ElemType: types.StringType},
		"containers": types.ListType{
			ElemType: types.ObjectType{AttrTypes: teoContainerAttrTypes()},
		},
		"resource_config": types.ObjectType{AttrTypes: teoResourceConfigAttrTypes()},
		"affinity_config": types.ObjectType{AttrTypes: teoAffinityConfigAttrTypes()},
	}
}

// teoDeploymentRecordAttrTypes returns the attr.Type map of a record_set
// element.
func teoDeploymentRecordAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"record_id": types.StringType,
		"operation": types.StringType,
		"status":    types.StringType,
		"duration":  types.Int64Type,
		"inference_service_config": types.ObjectType{
			AttrTypes: teoInferenceServiceConfigAttrTypes(),
		},
		"create_time":   types.StringType,
		"active_status": types.StringType,
	}
}

// teoDeploymentRecordToModel converts an SDK deployment record into the
// record_set element model. Nil pointers and nil nested objects become null
// values.
func teoDeploymentRecordToModel(ctx context.Context, record *teo.InferenceServiceDeploymentRecord) teoDeploymentRecordItemModel {
	var item teoDeploymentRecordItemModel
	if record == nil {
		return item
	}

	item.RecordId = types.StringPointerValue(record.RecordId)
	item.Operation = types.StringPointerValue(record.Operation)
	item.Status = types.StringPointerValue(record.Status)
	item.Duration = types.Int64PointerValue(record.Duration)
	item.CreateTime = types.StringPointerValue(record.CreateTime)
	item.ActiveStatus = types.StringPointerValue(record.ActiveStatus)

	if record.InferenceServiceConfig == nil {
		item.InferenceServiceConfig = types.ObjectNull(teoInferenceServiceConfigAttrTypes())
	} else {
		item.InferenceServiceConfig = teoInferenceServiceConfigToObject(
			ctx, record.InferenceServiceConfig)
	}
	return item
}

// teoInferenceServiceConfigToObject converts an SDK inference service config
// into a types.Object value.
func teoInferenceServiceConfigToObject(_ context.Context, config *teo.InferenceServiceConfig) types.Object {
	var model teoInferenceServiceConfigModel
	model.ListenPort = types.Int64PointerValue(config.ListenPort)

	if config.RequestPaths != nil {
		paths := make([]attr.Value, 0, len(config.RequestPaths))
		for _, p := range config.RequestPaths {
			paths = append(paths, types.StringPointerValue(p))
		}
		model.RequestPaths = types.SetValueMust(types.StringType, paths)
	} else {
		model.RequestPaths = types.SetNull(types.StringType)
	}

	containers := make([]teoContainerItemModel, 0, len(config.Containers))
	for _, container := range config.Containers {
		containers = append(containers, teoContainerToModel(container))
	}
	model.Containers = types.ListValueMust(
		types.ObjectType{AttrTypes: teoContainerAttrTypes()}, containersToAttrValues(containers))

	if config.ResourceConfig == nil {
		model.ResourceConfig = types.ObjectNull(teoResourceConfigAttrTypes())
	} else {
		model.ResourceConfig = teoResourceConfigToObject(config.ResourceConfig)
	}

	if config.AffinityConfig == nil {
		model.AffinityConfig = types.ObjectNull(teoAffinityConfigAttrTypes())
	} else {
		model.AffinityConfig = teoAffinityConfigToObject(config.AffinityConfig)
	}

	return types.ObjectValueMust(teoInferenceServiceConfigAttrTypes(), map[string]attr.Value{
		"listen_port":     model.ListenPort,
		"request_paths":   model.RequestPaths,
		"containers":      model.Containers,
		"resource_config": model.ResourceConfig,
		"affinity_config": model.AffinityConfig,
	})
}

// teoContainerToModel converts an SDK container config into the containers
// element model.
func teoContainerToModel(container *teo.InferenceContainerConfig) teoContainerItemModel {
	var model teoContainerItemModel
	if container == nil {
		return model
	}

	model.ImageType = types.StringPointerValue(container.ImageType)
	if container.TcrRepositoryConfig == nil {
		model.TcrRepositoryConfig = types.ObjectNull(teoTCRRepositoryConfigAttrTypes())
	} else {
		var tcrModel teoTCRRepositoryConfigModel
		tcrModel.TCRType = types.StringPointerValue(container.TcrRepositoryConfig.TCRType)
		tcrModel.Image = types.StringPointerValue(container.TcrRepositoryConfig.Image)
		tcrModel.RegistryId = types.StringPointerValue(container.TcrRepositoryConfig.RegistryId)
		tcrModel.RegionName = types.StringPointerValue(container.TcrRepositoryConfig.RegionName)
		model.TcrRepositoryConfig = types.ObjectValueMust(
			teoTCRRepositoryConfigAttrTypes(),
			map[string]attr.Value{
				"tcr_type":    tcrModel.TCRType,
				"image":       tcrModel.Image,
				"registry_id": tcrModel.RegistryId,
				"region_name": tcrModel.RegionName,
			})
	}
	model.StartupCommand = types.StringPointerValue(container.StartupCommand)

	envs := make([]teoEnvironmentVariableModel, 0, len(container.EnvironmentVariables))
	for _, env := range container.EnvironmentVariables {
		var envModel teoEnvironmentVariableModel
		if env != nil {
			envModel.Key = types.StringPointerValue(env.Key)
			envModel.Value = types.StringPointerValue(env.Value)
		}
		envs = append(envs, envModel)
	}
	model.EnvironmentVariables = types.ListValueMust(
		types.ObjectType{AttrTypes: teoEnvironmentVariableAttrTypes()},
		envVarsToAttrValues(envs))
	return model
}

// teoResourceConfigToObject converts an SDK resource config into a
// types.Object value.
func teoResourceConfigToObject(config *teo.InferenceResourceConfig) types.Object {
	var model teoResourceConfigModel
	model.ScalingMode = types.StringPointerValue(config.ScalingMode)
	model.HardwareSpec = types.StringPointerValue(config.HardwareSpec)
	model.HardwareSpecId = types.StringPointerValue(config.HardwareSpecId)
	model.Concurrency = types.Int64PointerValue(config.Concurrency)

	if config.HardwareConfig == nil {
		model.HardwareConfig = types.ObjectNull(teoHardwareConfigAttrTypes())
	} else {
		var hwModel teoHardwareConfigModel
		hwModel.GPUNum = types.Float64PointerValue(config.HardwareConfig.GPUNum)
		hwModel.CPUNum = types.Float64PointerValue(config.HardwareConfig.CPUNum)
		hwModel.MemSize = types.Int64PointerValue(config.HardwareConfig.MemSize)
		hwModel.DiskSize = types.Int64PointerValue(config.HardwareConfig.DiskSize)
		model.HardwareConfig = types.ObjectValueMust(
			teoHardwareConfigAttrTypes(),
			map[string]attr.Value{
				"gpu_num":   hwModel.GPUNum,
				"cpu_num":   hwModel.CPUNum,
				"mem_size":  hwModel.MemSize,
				"disk_size": hwModel.DiskSize,
			})
	}

	if config.AutoScalingConfig == nil {
		model.AutoScalingConfig = types.ObjectNull(teoAutoScalingConfigAttrTypes())
	} else {
		model.AutoScalingConfig = teoAutoScalingConfigToObject(config.AutoScalingConfig)
	}

	if config.ManualInstanceConfig == nil {
		model.ManualInstanceConfig = types.ObjectNull(teoManualInstanceConfigAttrTypes())
	} else {
		model.ManualInstanceConfig = types.ObjectValueMust(
			teoManualInstanceConfigAttrTypes(),
			map[string]attr.Value{
				"fixed_instance_count": types.Int64PointerValue(config.ManualInstanceConfig.FixedInstanceCount),
			})
	}

	return types.ObjectValueMust(teoResourceConfigAttrTypes(), map[string]attr.Value{
		"scaling_mode":           model.ScalingMode,
		"hardware_spec":          model.HardwareSpec,
		"hardware_spec_id":       model.HardwareSpecId,
		"hardware_config":        model.HardwareConfig,
		"auto_scaling_config":    model.AutoScalingConfig,
		"manual_instance_config": model.ManualInstanceConfig,
		"concurrency":            model.Concurrency,
	})
}

// teoAutoScalingConfigToObject converts an SDK auto scaling config into a
// types.Object value.
func teoAutoScalingConfigToObject(config *teo.InferenceAutoScalingConfig) types.Object {
	var model teoAutoScalingConfigModel
	model.MinInstanceCount = types.Int64PointerValue(config.MinInstanceCount)

	policies := make([]teoScalingPolicyModel, 0, len(config.ScalingPolicies))
	for _, policy := range config.ScalingPolicies {
		var policyModel teoScalingPolicyModel
		if policy != nil {
			policyModel.PolicyName = types.StringPointerValue(policy.PolicyName)
			policyModel.PolicyType = types.StringPointerValue(policy.PolicyType)
			if policy.ScheduledScalingPolicy == nil {
				policyModel.ScheduledScalingPolicy = types.ObjectNull(teoScheduledScalingPolicyAttrTypes())
			} else {
				policyModel.ScheduledScalingPolicy = teoScheduledScalingPolicyToObject(policy.ScheduledScalingPolicy)
			}
		}
		policies = append(policies, policyModel)
	}
	policyValues := make([]attr.Value, 0, len(policies))
	for _, p := range policies {
		policyValues = append(policyValues, types.ObjectValueMust(teoScalingPolicyAttrTypes(), map[string]attr.Value{
			"policy_name":              p.PolicyName,
			"policy_type":              p.PolicyType,
			"scheduled_scaling_policy": p.ScheduledScalingPolicy,
		}))
	}
	model.ScalingPolicies = types.ListValueMust(
		types.ObjectType{AttrTypes: teoScalingPolicyAttrTypes()}, policyValues)

	return types.ObjectValueMust(teoAutoScalingConfigAttrTypes(), map[string]attr.Value{
		"min_instance_count": model.MinInstanceCount,
		"scaling_policies":   model.ScalingPolicies,
	})
}

// teoScheduledScalingPolicyToObject converts an SDK scheduled scaling policy
// into a types.Object value.
func teoScheduledScalingPolicyToObject(policy *teo.InferenceScheduledScalingPolicy) types.Object {
	var model teoScheduledScalingPolicyModel

	actions := make([]attr.Value, 0, len(policy.ScheduledActions))
	for _, action := range policy.ScheduledActions {
		if action == nil {
			actions = append(actions, types.ObjectNull(teoScheduledScalingActionAttrTypes()))
			continue
		}
		actions = append(actions, types.ObjectValueMust(teoScheduledScalingActionAttrTypes(), map[string]attr.Value{
			"cron_expression":    types.StringPointerValue(action.CronExpression),
			"min_instance_count": types.Int64PointerValue(action.MinInstanceCount),
		}))
	}
	model.ScheduledActions = types.ListValueMust(
		types.ObjectType{AttrTypes: teoScheduledScalingActionAttrTypes()}, actions)

	if policy.EffectiveRange == nil {
		model.EffectiveRange = types.ObjectNull(teoScheduledScalingEffectiveRangeAttrTypes())
	} else {
		model.EffectiveRange = types.ObjectValueMust(
			teoScheduledScalingEffectiveRangeAttrTypes(),
			map[string]attr.Value{
				"effective_type": types.StringPointerValue(policy.EffectiveRange.EffectiveType),
				"start_date":     types.StringPointerValue(policy.EffectiveRange.StartDate),
				"end_date":       types.StringPointerValue(policy.EffectiveRange.EndDate),
			})
	}
	model.TimeZone = types.StringPointerValue(policy.TimeZone)

	return types.ObjectValueMust(teoScheduledScalingPolicyAttrTypes(), map[string]attr.Value{
		"scheduled_actions": model.ScheduledActions,
		"effective_range":   model.EffectiveRange,
		"time_zone":         model.TimeZone,
	})
}

// teoAffinityConfigToObject converts an SDK affinity config into a
// types.Object value.
func teoAffinityConfigToObject(config *teo.InferenceAffinityConfig) types.Object {
	var model teoAffinityConfigModel
	model.Switch = types.StringPointerValue(config.Switch)
	model.AffinityMode = types.StringPointerValue(config.AffinityMode)

	if config.SessionIdAffinityConfig == nil {
		model.SessionIdAffinityConfig = types.ObjectNull(teoSessionIdAffinityConfigAttrTypes())
	} else {
		model.SessionIdAffinityConfig = types.ObjectValueMust(
			teoSessionIdAffinityConfigAttrTypes(),
			map[string]attr.Value{
				"source":      types.StringPointerValue(config.SessionIdAffinityConfig.Source),
				"header_name": types.StringPointerValue(config.SessionIdAffinityConfig.HeaderName),
			})
	}

	return types.ObjectValueMust(teoAffinityConfigAttrTypes(), map[string]attr.Value{
		"switch":                     model.Switch,
		"affinity_mode":              model.AffinityMode,
		"session_id_affinity_config": model.SessionIdAffinityConfig,
	})
}

// containersToAttrValues converts container models to a slice of attr.Value.
func containersToAttrValues(containers []teoContainerItemModel) []attr.Value {
	values := make([]attr.Value, 0, len(containers))
	for _, c := range containers {
		values = append(values, types.ObjectValueMust(teoContainerAttrTypes(), map[string]attr.Value{
			"image_type":            c.ImageType,
			"tcr_repository_config": c.TcrRepositoryConfig,
			"startup_command":       c.StartupCommand,
			"environment_variables": c.EnvironmentVariables,
		}))
	}
	return values
}

// envVarsToAttrValues converts environment variable models to a slice of
// attr.Value.
func envVarsToAttrValues(envs []teoEnvironmentVariableModel) []attr.Value {
	values := make([]attr.Value, 0, len(envs))
	for _, e := range envs {
		values = append(values, types.ObjectValueMust(teoEnvironmentVariableAttrTypes(), map[string]attr.Value{
			"key":   e.Key,
			"value": e.Value,
		}))
	}
	return values
}

// Read queries the cloud API and writes the result into resp.State.
func (d *TeoInferenceServiceDeploymentRecords111001DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TeoInferenceServiceDeploymentRecords111001DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if d.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	service := NewTeoService(d.Client())
	records, _, err := service.DescribeInferenceServiceDeploymentRecords(
		ctx,
		config.ZoneId.ValueString(),
		config.ServiceId.ValueString(),
		config.SortBy.ValueString(),
		config.SortOrder.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying TEO inference service deployment records",
			fmt.Sprintf("DescribeInferenceServiceDeploymentRecords failed: %v", err),
		)
		return
	}

	items := make([]teoDeploymentRecordItemModel, 0, len(records))
	for _, record := range records {
		items = append(items, teoDeploymentRecordToModel(ctx, record))
	}

	recordValues := make([]attr.Value, 0, len(items))
	for _, item := range items {
		recordValues = append(recordValues, types.ObjectValueMust(teoDeploymentRecordAttrTypes(), map[string]attr.Value{
			"record_id":                item.RecordId,
			"operation":                item.Operation,
			"status":                   item.Status,
			"duration":                 item.Duration,
			"inference_service_config": item.InferenceServiceConfig,
			"create_time":              item.CreateTime,
			"active_status":            item.ActiveStatus,
		}))
	}

	list, diags := types.ListValue(
		types.ObjectType{AttrTypes: teoDeploymentRecordAttrTypes()}, recordValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.RecordSet = list
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
