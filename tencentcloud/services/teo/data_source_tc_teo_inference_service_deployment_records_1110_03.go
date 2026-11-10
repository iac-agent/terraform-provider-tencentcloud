package teo

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// _ datasource.DataSource is intentionally omitted: the factory function
// NewTeoInferenceServiceDeploymentRecords111003DataSource returns
// datasource.DataSource, which already performs the compile-time check at the
// return statement. Only the extended-interface assertion is kept here to
// verify the Configure method signature promoted from
// fw.DataSourceWithConfigure.
var _ datasource.DataSourceWithConfigure = &TeoInferenceServiceDeploymentRecords111003DataSource{}

// NewTeoInferenceServiceDeploymentRecords111003DataSource is the factory
// referenced by tencentcloud/framework/registry.go to register this data
// source.
func NewTeoInferenceServiceDeploymentRecords111003DataSource() datasource.DataSource {
	return &TeoInferenceServiceDeploymentRecords111003DataSource{}
}

// TeoInferenceServiceDeploymentRecords111003DataSource implements
// datasource.DataSource for
// tencentcloud_teo_inference_service_deployment_records_1110_03.
type TeoInferenceServiceDeploymentRecords111003DataSource struct {
	fw.DataSourceWithConfigure
}

// TeoInferenceServiceDeploymentRecords111003DataSourceModel maps the schema
// attributes.
type TeoInferenceServiceDeploymentRecords111003DataSourceModel struct {
	ZoneId     types.String `tfsdk:"zone_id"`
	ServiceId  types.String `tfsdk:"service_id"`
	SortBy     types.String `tfsdk:"sort_by"`
	SortOrder  types.String `tfsdk:"sort_order"`
	TotalCount types.Int64  `tfsdk:"total_count"`
	RecordSet  types.List   `tfsdk:"record_set"`
}

var (
	teoInferenceScheduledScalingActionAttrTypes = map[string]attr.Type{
		"cron_expression":    types.StringType,
		"min_instance_count": types.Int64Type,
	}

	teoInferenceScheduledScalingEffectiveRangeAttrTypes = map[string]attr.Type{
		"effective_type": types.StringType,
		"start_date":     types.StringType,
		"end_date":       types.StringType,
	}

	teoInferenceScheduledScalingPolicyAttrTypes = map[string]attr.Type{
		"scheduled_actions": types.ListType{ElemType: types.ObjectType{
			AttrTypes: teoInferenceScheduledScalingActionAttrTypes,
		}},
		"effective_range": types.ListType{ElemType: types.ObjectType{
			AttrTypes: teoInferenceScheduledScalingEffectiveRangeAttrTypes,
		}},
		"time_zone": types.StringType,
	}

	teoInferenceScalingPolicyAttrTypes = map[string]attr.Type{
		"policy_name": types.StringType,
		"policy_type": types.StringType,
		"scheduled_scaling_policy": types.ListType{ElemType: types.ObjectType{
			AttrTypes: teoInferenceScheduledScalingPolicyAttrTypes,
		}},
	}

	teoInferenceAutoScalingConfigAttrTypes = map[string]attr.Type{
		"min_instance_count": types.Int64Type,
		"scaling_policies": types.ListType{ElemType: types.ObjectType{
			AttrTypes: teoInferenceScalingPolicyAttrTypes,
		}},
	}

	teoInferenceManualInstanceConfigAttrTypes = map[string]attr.Type{
		"fixed_instance_count": types.Int64Type,
	}

	teoInferenceHardwareConfigAttrTypes = map[string]attr.Type{
		"gpu_num":   types.Float64Type,
		"cpu_num":   types.Float64Type,
		"mem_size":  types.Int64Type,
		"disk_size": types.Int64Type,
	}

	teoInferenceResourceConfigAttrTypes = map[string]attr.Type{
		"scaling_mode":     types.StringType,
		"hardware_spec":    types.StringType,
		"hardware_spec_id": types.StringType,
		"hardware_config": types.ListType{ElemType: types.ObjectType{
			AttrTypes: teoInferenceHardwareConfigAttrTypes,
		}},
		"auto_scaling_config": types.ListType{ElemType: types.ObjectType{
			AttrTypes: teoInferenceAutoScalingConfigAttrTypes,
		}},
		"manual_instance_config": types.ListType{ElemType: types.ObjectType{
			AttrTypes: teoInferenceManualInstanceConfigAttrTypes,
		}},
		"concurrency": types.Int64Type,
	}

	teoInferenceEnvironmentVariableAttrTypes = map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	}

	teoInferenceTcrRepositoryConfigAttrTypes = map[string]attr.Type{
		"tcr_type":    types.StringType,
		"image":       types.StringType,
		"registry_id": types.StringType,
		"region_name": types.StringType,
	}

	teoInferenceContainerAttrTypes = map[string]attr.Type{
		"image_type": types.StringType,
		"tcr_repository_config": types.ListType{ElemType: types.ObjectType{
			AttrTypes: teoInferenceTcrRepositoryConfigAttrTypes,
		}},
		"startup_command": types.StringType,
		"environment_variables": types.ListType{ElemType: types.ObjectType{
			AttrTypes: teoInferenceEnvironmentVariableAttrTypes,
		}},
	}

	teoInferenceSessionIdAffinityConfigAttrTypes = map[string]attr.Type{
		"source":      types.StringType,
		"header_name": types.StringType,
	}

	teoInferenceAffinityConfigAttrTypes = map[string]attr.Type{
		"switch":        types.StringType,
		"affinity_mode": types.StringType,
		"session_id_affinity_config": types.ListType{ElemType: types.ObjectType{
			AttrTypes: teoInferenceSessionIdAffinityConfigAttrTypes,
		}},
	}

	teoInferenceServiceConfigAttrTypes = map[string]attr.Type{
		"listen_port":   types.Int64Type,
		"request_paths": types.SetType{ElemType: types.StringType},
		"containers": types.ListType{ElemType: types.ObjectType{
			AttrTypes: teoInferenceContainerAttrTypes,
		}},
		"resource_config": types.ListType{ElemType: types.ObjectType{
			AttrTypes: teoInferenceResourceConfigAttrTypes,
		}},
		"affinity_config": types.ListType{ElemType: types.ObjectType{
			AttrTypes: teoInferenceAffinityConfigAttrTypes,
		}},
	}

	teoInferenceDeploymentRecordAttrTypes = map[string]attr.Type{
		"record_id":  types.StringType,
		"operation":  types.StringType,
		"status":     types.StringType,
		"duration":   types.Int64Type,
		"inference_service_config": types.ListType{ElemType: types.ObjectType{
			AttrTypes: teoInferenceServiceConfigAttrTypes,
		}},
		"create_time":   types.StringType,
		"active_status": types.StringType,
	}
)

func (d *TeoInferenceServiceDeploymentRecords111003DataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_inference_service_deployment_records_1110_03"
}

func (d *TeoInferenceServiceDeploymentRecords111003DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the deployment records of a TEO inference service.",
		Attributes: map[string]schema.Attribute{
			// ---- Query inputs: Required / Optional only ----
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Zone ID.",
			},
			"service_id": schema.StringAttribute{
				Required:    true,
				Description: "Inference service ID.",
			},
			"sort_by": schema.StringAttribute{
				Optional:    true,
				Description: "Field used to sort the deployment records. Valid value: `create-time` (deployment creation time). The cloud default is `create-time`.",
			},
			"sort_order": schema.StringAttribute{
				Optional:    true,
				Description: "Sorting order. Valid values: `asc` (ascending), `desc` (descending). The cloud default is `desc`.",
			},

			// ---- Query outputs: Computed only ----
			"total_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Total number of deployment records.",
			},
			"record_set": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Deployment records of the inference service.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"record_id": schema.StringAttribute{
							Computed:    true,
							Description: "Deployment record ID.",
						},
						"operation": schema.StringAttribute{
							Computed:    true,
							Description: "Deployment operation type. Valid values: `create`, `update`, `resume`, `stop`.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Deployment status. Valid values: `processing`, `succeeded`, `failed`.",
						},
						"duration": schema.Int64Attribute{
							Computed:    true,
							Description: "Deployment duration, in seconds.",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "Deployment start time in ISO date format.",
						},
						"active_status": schema.StringAttribute{
							Computed:    true,
							Description: "Whether this deployment configuration is the currently effective one. Valid values: `active`, `inactive`.",
						},
						"inference_service_config": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Configuration of this inference service deployment.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"listen_port": schema.Int64Attribute{
										Computed:    true,
										Description: "Port the model service listens on.",
									},
									"request_paths": schema.SetAttribute{
										Computed:    true,
										ElementType: types.StringType,
										Description: "Request paths of the inference service.",
									},
									"containers": schema.ListNestedAttribute{
										Computed:    true,
										Description: "Container configuration of the inference service.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"image_type": schema.StringAttribute{
													Computed:    true,
													Description: "Image type. Valid value: `TCR`.",
												},
												"startup_command": schema.StringAttribute{
													Computed:    true,
													Description: "Command executed when the container starts. When it is not set, the Entrypoint/CMD of the image is used.",
												},
												"tcr_repository_config": schema.ListNestedAttribute{
													Computed:    true,
													Description: "TCR image repository information.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"tcr_type": schema.StringAttribute{
																Computed:    true,
																Description: "TCR service type. Valid values: `Personal`, `Enterprise`.",
															},
															"image": schema.StringAttribute{
																Computed:    true,
																Description: "Image address.",
															},
															"registry_id": schema.StringAttribute{
																Computed:    true,
																Description: "Image registry instance ID. Required when TCRType is `Enterprise`.",
															},
															"region_name": schema.StringAttribute{
																Computed:    true,
																Description: "Region name.",
															},
														},
													},
												},
												"environment_variables": schema.ListNestedAttribute{
													Computed:    true,
													Description: "Environment variables of the container runtime.",
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
									"resource_config": schema.ListNestedAttribute{
										Computed:    true,
										Description: "Resource configuration of the inference service.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"scaling_mode": schema.StringAttribute{
													Computed:    true,
													Description: "Scaling mode. Valid values: `Auto`, `Manual`.",
												},
												"hardware_spec": schema.StringAttribute{
													Computed:    true,
													Description: "Hardware specification identifier. Deprecated, use `hardware_spec_id` instead.",
												},
												"hardware_spec_id": schema.StringAttribute{
													Computed:    true,
													Description: "Unique identifier of the hardware specification.",
												},
												"concurrency": schema.Int64Attribute{
													Computed:    true,
													Description: "Concurrency of a single instance.",
												},
												"hardware_config": schema.ListNestedAttribute{
													Computed:    true,
													Description: "Hardware configuration of the inference service.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"gpu_num": schema.Float64Attribute{
																Computed:    true,
																Description: "Number of GPU cards allocated to a single instance.",
															},
															"cpu_num": schema.Float64Attribute{
																Computed:    true,
																Description: "Number of CPU cores allocated to a single instance.",
															},
															"mem_size": schema.Int64Attribute{
																Computed:    true,
																Description: "Memory size allocated to a single instance, in MB.",
															},
															"disk_size": schema.Int64Attribute{
																Computed:    true,
																Description: "Temporary disk size allocated to a single instance, in MB.",
															},
														},
													},
												},
												"auto_scaling_config": schema.ListNestedAttribute{
													Computed:    true,
													Description: "Auto scaling configuration. Required when ScalingMode is `Auto`.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"min_instance_count": schema.Int64Attribute{
																Computed:    true,
																Description: "Minimum number of instances.",
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
																			Computed:    true,
																			Description: "Policy type. Valid value: `ScheduledScaling`.",
																		},
																		"scheduled_scaling_policy": schema.ListNestedAttribute{
																			Computed:    true,
																			Description: "Scheduled scaling configuration. Required when PolicyType is `ScheduledScaling`.",
																			NestedObject: schema.NestedAttributeObject{
																				Attributes: map[string]schema.Attribute{
																					"time_zone": schema.StringAttribute{
																						Computed:    true,
																						Description: "Time zone (IANA time zone name) used to evaluate the scheduled actions. The cloud default is `UTC`.",
																					},
																					"scheduled_actions": schema.ListNestedAttribute{
																						Computed:    true,
																						Description: "Scheduled scaling actions.",
																						NestedObject: schema.NestedAttributeObject{
																							Attributes: map[string]schema.Attribute{
																								"cron_expression": schema.StringAttribute{
																									Computed:    true,
																									Description: "Cron expression describing when the scheduled action is triggered. Standard five-field cron format, seconds and years are not supported.",
																								},
																								"min_instance_count": schema.Int64Attribute{
																									Computed:    true,
																									Description: "Minimum number of instances the inference service is scaled to when this action is triggered.",
																								},
																							},
																						},
																					},
																					"effective_range": schema.ListNestedAttribute{
																						Computed:    true,
																						Description: "Validity range of the scheduled scaling policy.",
																						NestedObject: schema.NestedAttributeObject{
																							Attributes: map[string]schema.Attribute{
																								"effective_type": schema.StringAttribute{
																									Computed:    true,
																									Description: "Validity type. Valid values: `LongTerm`, `Custom`.",
																								},
																								"start_date": schema.StringAttribute{
																									Computed:    true,
																									Description: "Start date of the validity range. Set when EffectiveType is `Custom`.",
																								},
																								"end_date": schema.StringAttribute{
																									Computed:    true,
																									Description: "End date of the validity range. Set when EffectiveType is `Custom`.",
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
														},
													},
												},
												"manual_instance_config": schema.ListNestedAttribute{
													Computed:    true,
													Description: "Manual instance configuration. Required when ScalingMode is `Manual`.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"fixed_instance_count": schema.Int64Attribute{
																Computed:    true,
																Description: "Fixed number of instances.",
															},
														},
													},
												},
											},
										},
									},
									"affinity_config": schema.ListNestedAttribute{
										Computed:    true,
										Description: "Affinity configuration of the inference service.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"switch": schema.StringAttribute{
													Computed:    true,
													Description: "Master switch of the inference service affinity. Valid values: `On`, `Off`.",
												},
												"affinity_mode": schema.StringAttribute{
													Computed:    true,
													Description: "Affinity mode. Valid value: `SessionId`.",
												},
												"session_id_affinity_config": schema.ListNestedAttribute{
													Computed:    true,
													Description: "Affinity configuration used when AffinityMode is `SessionId`.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"source": schema.StringAttribute{
																Computed:    true,
																Description: "Where the session ID parameter is passed. Valid value: `Header`.",
															},
															"header_name": schema.StringAttribute{
																Computed:    true,
																Description: "Name of the request header carrying the session ID.",
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
					},
				},
			},
		},
	}
}

// Read queries the deployment records of the inference service and writes the
// result into resp.State.
func (d *TeoInferenceServiceDeploymentRecords111003DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TeoInferenceServiceDeploymentRecords111003DataSourceModel
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

	zoneId := ""
	if !config.ZoneId.IsNull() && !config.ZoneId.IsUnknown() {
		zoneId = config.ZoneId.ValueString()
	}
	serviceId := ""
	if !config.ServiceId.IsNull() && !config.ServiceId.IsUnknown() {
		serviceId = config.ServiceId.ValueString()
	}
	if zoneId == "" || serviceId == "" {
		resp.Diagnostics.AddError(
			"Missing required query parameter",
			"Both zone_id and service_id are required and cannot be empty.",
		)
		return
	}

	param := map[string]interface{}{
		"ZoneId":    zoneId,
		"ServiceId": serviceId,
	}
	if !config.SortBy.IsNull() && !config.SortBy.IsUnknown() && config.SortBy.ValueString() != "" {
		param["SortBy"] = config.SortBy.ValueString()
	}
	if !config.SortOrder.IsNull() && !config.SortOrder.IsUnknown() && config.SortOrder.ValueString() != "" {
		param["SortOrder"] = config.SortOrder.ValueString()
	}

	service := NewTeoService(d.Client())
	records, totalCount, err := service.DescribeTeoInferenceServiceDeploymentRecordsByFilter(ctx, param)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying TEO inference service deployment records",
			fmt.Sprintf("DescribeInferenceServiceDeploymentRecords failed: %v", err),
		)
		return
	}

	// An empty record set is a valid business state; a nil response is turned
	// into an error by the service layer, so it never reaches this point.
	items, diags := flattenTeoInferenceDeploymentRecords(ctx, records)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.ZoneId = types.StringValue(zoneId)
	config.ServiceId = types.StringValue(serviceId)
	// TotalCount is optional in the API response; when the cloud does not
	// return it, fall back to the number of records actually collected.
	if totalCount <= 0 {
		totalCount = int64(len(records))
	}
	config.TotalCount = types.Int64Value(totalCount)
	config.RecordSet = items

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// flattenTeoInferenceDeploymentRecords converts the SDK record set into the
// terraform list value of record_set. A nested object that is absent in the
// response is rendered as a null list, so "not returned" stays distinguishable
// from "returned empty".
func flattenTeoInferenceDeploymentRecords(ctx context.Context, records []*teo.InferenceServiceDeploymentRecord) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	elements := make([]attr.Value, 0, len(records))
	for _, record := range records {
		if record == nil {
			continue
		}

		serviceConfig, d := flattenTeoInferenceServiceConfig(ctx, record.InferenceServiceConfig)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: teoInferenceDeploymentRecordAttrTypes}), diags
		}

		obj, d := types.ObjectValue(teoInferenceDeploymentRecordAttrTypes, map[string]attr.Value{
			"record_id":                types.StringPointerValue(record.RecordId),
			"operation":                types.StringPointerValue(record.Operation),
			"status":                   types.StringPointerValue(record.Status),
			"duration":                 types.Int64PointerValue(record.Duration),
			"inference_service_config": serviceConfig,
			"create_time":              types.StringPointerValue(record.CreateTime),
			"active_status":            types.StringPointerValue(record.ActiveStatus),
		})
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: teoInferenceDeploymentRecordAttrTypes}), diags
		}
		elements = append(elements, obj)
	}

	list, d := types.ListValue(types.ObjectType{AttrTypes: teoInferenceDeploymentRecordAttrTypes}, elements)
	diags.Append(d...)
	return list, diags
}

// flattenTeoInferenceServiceConfig converts one InferenceServiceConfig into the
// nested list used by record_set.inference_service_config.
func flattenTeoInferenceServiceConfig(ctx context.Context, config *teo.InferenceServiceConfig) (types.List, diag.Diagnostics) {
	objType := types.ObjectType{AttrTypes: teoInferenceServiceConfigAttrTypes}
	if config == nil {
		return types.ListNull(objType), nil
	}

	var diags diag.Diagnostics

	requestPathValues := make([]attr.Value, 0, len(config.RequestPaths))
	for _, path := range config.RequestPaths {
		requestPathValues = append(requestPathValues, types.StringPointerValue(path))
	}
	requestPaths, d := types.SetValue(types.StringType, requestPathValues)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	containers, d := flattenTeoInferenceContainers(ctx, config.Containers)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	resourceConfig, d := flattenTeoInferenceResourceConfig(ctx, config.ResourceConfig)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	affinityConfig, d := flattenTeoInferenceAffinityConfig(ctx, config.AffinityConfig)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	obj, d := types.ObjectValue(teoInferenceServiceConfigAttrTypes, map[string]attr.Value{
		"listen_port":     types.Int64PointerValue(config.ListenPort),
		"request_paths":   requestPaths,
		"containers":      containers,
		"resource_config": resourceConfig,
		"affinity_config": affinityConfig,
	})
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	list, d := types.ListValue(objType, []attr.Value{obj})
	diags.Append(d...)
	return list, diags
}

// flattenTeoInferenceContainers converts the container configuration list.
func flattenTeoInferenceContainers(ctx context.Context, containers []*teo.InferenceContainerConfig) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	objType := types.ObjectType{AttrTypes: teoInferenceContainerAttrTypes}
	elements := make([]attr.Value, 0, len(containers))
	for _, container := range containers {
		if container == nil {
			continue
		}

		tcrConfig, d := flattenTeoInferenceTcrRepositoryConfig(ctx, container.TcrRepositoryConfig)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}

		envVars, d := flattenTeoInferenceEnvironmentVariables(ctx, container.EnvironmentVariables)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}

		obj, d := types.ObjectValue(teoInferenceContainerAttrTypes, map[string]attr.Value{
			"image_type":            types.StringPointerValue(container.ImageType),
			"tcr_repository_config": tcrConfig,
			"startup_command":       types.StringPointerValue(container.StartupCommand),
			"environment_variables": envVars,
		})
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}
		elements = append(elements, obj)
	}

	list, d := types.ListValue(objType, elements)
	diags.Append(d...)
	return list, diags
}

// flattenTeoInferenceTcrRepositoryConfig converts the TCR repository
// configuration. The API may return null, which becomes a null list.
func flattenTeoInferenceTcrRepositoryConfig(ctx context.Context, config *teo.InferenceTCRRepositoryConfig) (types.List, diag.Diagnostics) {
	objType := types.ObjectType{AttrTypes: teoInferenceTcrRepositoryConfigAttrTypes}
	if config == nil {
		return types.ListNull(objType), nil
	}

	obj, diags := types.ObjectValue(teoInferenceTcrRepositoryConfigAttrTypes, map[string]attr.Value{
		"tcr_type":    types.StringPointerValue(config.TCRType),
		"image":       types.StringPointerValue(config.Image),
		"registry_id": types.StringPointerValue(config.RegistryId),
		"region_name": types.StringPointerValue(config.RegionName),
	})
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	list, d := types.ListValue(objType, []attr.Value{obj})
	diags.Append(d...)
	return list, diags
}

// flattenTeoInferenceEnvironmentVariables converts the environment variables.
// The API does not define an ordering, so they are sorted by key to keep the
// output deterministic across reads.
func flattenTeoInferenceEnvironmentVariables(ctx context.Context, envVars []*teo.InferenceEnvironmentVariable) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	objType := types.ObjectType{AttrTypes: teoInferenceEnvironmentVariableAttrTypes}
	elements := make([]attr.Value, 0, len(envVars))
	sortedEnvVars := make([]*teo.InferenceEnvironmentVariable, 0, len(envVars))
	for _, envVar := range envVars {
		if envVar == nil {
			continue
		}
		sortedEnvVars = append(sortedEnvVars, envVar)
	}
	sort.SliceStable(sortedEnvVars, func(i, j int) bool {
		return types.StringPointerValue(sortedEnvVars[i].Key).ValueString() <
			types.StringPointerValue(sortedEnvVars[j].Key).ValueString()
	})

	for _, envVar := range sortedEnvVars {
		obj, d := types.ObjectValue(teoInferenceEnvironmentVariableAttrTypes, map[string]attr.Value{
			"key":   types.StringPointerValue(envVar.Key),
			"value": types.StringPointerValue(envVar.Value),
		})
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}
		elements = append(elements, obj)
	}

	list, d := types.ListValue(objType, elements)
	diags.Append(d...)
	return list, diags
}

// flattenTeoInferenceResourceConfig converts the resource configuration.
func flattenTeoInferenceResourceConfig(ctx context.Context, config *teo.InferenceResourceConfig) (types.List, diag.Diagnostics) {
	objType := types.ObjectType{AttrTypes: teoInferenceResourceConfigAttrTypes}
	if config == nil {
		return types.ListNull(objType), nil
	}

	var diags diag.Diagnostics

	hardwareConfig, d := flattenTeoInferenceHardwareConfig(ctx, config.HardwareConfig)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	autoScalingConfig, d := flattenTeoInferenceAutoScalingConfig(ctx, config.AutoScalingConfig)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	manualInstanceConfig, d := flattenTeoInferenceManualInstanceConfig(ctx, config.ManualInstanceConfig)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	obj, d := types.ObjectValue(teoInferenceResourceConfigAttrTypes, map[string]attr.Value{
		"scaling_mode":           types.StringPointerValue(config.ScalingMode),
		"hardware_spec":          types.StringPointerValue(config.HardwareSpec),
		"hardware_spec_id":       types.StringPointerValue(config.HardwareSpecId),
		"hardware_config":        hardwareConfig,
		"auto_scaling_config":    autoScalingConfig,
		"manual_instance_config": manualInstanceConfig,
		"concurrency":            types.Int64PointerValue(config.Concurrency),
	})
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	list, d := types.ListValue(objType, []attr.Value{obj})
	diags.Append(d...)
	return list, diags
}

// flattenTeoInferenceHardwareConfig converts the hardware configuration.
func flattenTeoInferenceHardwareConfig(ctx context.Context, config *teo.InferenceHardwareConfig) (types.List, diag.Diagnostics) {
	objType := types.ObjectType{AttrTypes: teoInferenceHardwareConfigAttrTypes}
	if config == nil {
		return types.ListNull(objType), nil
	}

	obj, diags := types.ObjectValue(teoInferenceHardwareConfigAttrTypes, map[string]attr.Value{
		"gpu_num":   types.Float64PointerValue(config.GPUNum),
		"cpu_num":   types.Float64PointerValue(config.CPUNum),
		"mem_size":  types.Int64PointerValue(config.MemSize),
		"disk_size": types.Int64PointerValue(config.DiskSize),
	})
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	list, d := types.ListValue(objType, []attr.Value{obj})
	diags.Append(d...)
	return list, diags
}

// flattenTeoInferenceAutoScalingConfig converts the auto scaling configuration.
// The API may return null, which becomes a null list.
func flattenTeoInferenceAutoScalingConfig(ctx context.Context, config *teo.InferenceAutoScalingConfig) (types.List, diag.Diagnostics) {
	objType := types.ObjectType{AttrTypes: teoInferenceAutoScalingConfigAttrTypes}
	if config == nil {
		return types.ListNull(objType), nil
	}

	var diags diag.Diagnostics

	policies, d := flattenTeoInferenceScalingPolicies(ctx, config.ScalingPolicies)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	obj, d := types.ObjectValue(teoInferenceAutoScalingConfigAttrTypes, map[string]attr.Value{
		"min_instance_count": types.Int64PointerValue(config.MinInstanceCount),
		"scaling_policies":   policies,
	})
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	list, d := types.ListValue(objType, []attr.Value{obj})
	diags.Append(d...)
	return list, diags
}

// flattenTeoInferenceScalingPolicies converts the scaling policy list.
func flattenTeoInferenceScalingPolicies(ctx context.Context, policies []*teo.InferenceScalingPolicy) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	objType := types.ObjectType{AttrTypes: teoInferenceScalingPolicyAttrTypes}
	elements := make([]attr.Value, 0, len(policies))
	for _, policy := range policies {
		if policy == nil {
			continue
		}

		scheduled, d := flattenTeoInferenceScheduledScalingPolicy(ctx, policy.ScheduledScalingPolicy)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}

		obj, d := types.ObjectValue(teoInferenceScalingPolicyAttrTypes, map[string]attr.Value{
			"policy_name":              types.StringPointerValue(policy.PolicyName),
			"policy_type":              types.StringPointerValue(policy.PolicyType),
			"scheduled_scaling_policy": scheduled,
		})
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}
		elements = append(elements, obj)
	}

	list, d := types.ListValue(objType, elements)
	diags.Append(d...)
	return list, diags
}

// flattenTeoInferenceScheduledScalingPolicy converts the scheduled scaling
// configuration.
func flattenTeoInferenceScheduledScalingPolicy(ctx context.Context, config *teo.InferenceScheduledScalingPolicy) (types.List, diag.Diagnostics) {
	objType := types.ObjectType{AttrTypes: teoInferenceScheduledScalingPolicyAttrTypes}
	if config == nil {
		return types.ListNull(objType), nil
	}

	var diags diag.Diagnostics

	actions, d := flattenTeoInferenceScheduledScalingActions(ctx, config.ScheduledActions)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	effectiveRange, d := flattenTeoInferenceScheduledScalingEffectiveRange(ctx, config.EffectiveRange)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	obj, d := types.ObjectValue(teoInferenceScheduledScalingPolicyAttrTypes, map[string]attr.Value{
		"scheduled_actions": actions,
		"effective_range":   effectiveRange,
		"time_zone":         types.StringPointerValue(config.TimeZone),
	})
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	list, d := types.ListValue(objType, []attr.Value{obj})
	diags.Append(d...)
	return list, diags
}

// flattenTeoInferenceScheduledScalingActions converts the scheduled scaling
// actions.
func flattenTeoInferenceScheduledScalingActions(ctx context.Context, actions []*teo.InferenceScheduledScalingAction) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	objType := types.ObjectType{AttrTypes: teoInferenceScheduledScalingActionAttrTypes}
	elements := make([]attr.Value, 0, len(actions))
	for _, action := range actions {
		if action == nil {
			continue
		}
		obj, d := types.ObjectValue(teoInferenceScheduledScalingActionAttrTypes, map[string]attr.Value{
			"cron_expression":    types.StringPointerValue(action.CronExpression),
			"min_instance_count": types.Int64PointerValue(action.MinInstanceCount),
		})
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(objType), diags
		}
		elements = append(elements, obj)
	}

	list, d := types.ListValue(objType, elements)
	diags.Append(d...)
	return list, diags
}

// flattenTeoInferenceScheduledScalingEffectiveRange converts the validity range
// of a scheduled scaling policy.
func flattenTeoInferenceScheduledScalingEffectiveRange(ctx context.Context, config *teo.InferenceScheduledScalingEffectiveRange) (types.List, diag.Diagnostics) {
	objType := types.ObjectType{AttrTypes: teoInferenceScheduledScalingEffectiveRangeAttrTypes}
	if config == nil {
		return types.ListNull(objType), nil
	}

	obj, diags := types.ObjectValue(teoInferenceScheduledScalingEffectiveRangeAttrTypes, map[string]attr.Value{
		"effective_type": types.StringPointerValue(config.EffectiveType),
		"start_date":     types.StringPointerValue(config.StartDate),
		"end_date":       types.StringPointerValue(config.EndDate),
	})
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	list, d := types.ListValue(objType, []attr.Value{obj})
	diags.Append(d...)
	return list, diags
}

// flattenTeoInferenceManualInstanceConfig converts the manual instance
// configuration. The API may return null, which becomes a null list.
func flattenTeoInferenceManualInstanceConfig(ctx context.Context, config *teo.InferenceManualInstanceConfig) (types.List, diag.Diagnostics) {
	objType := types.ObjectType{AttrTypes: teoInferenceManualInstanceConfigAttrTypes}
	if config == nil {
		return types.ListNull(objType), nil
	}

	obj, diags := types.ObjectValue(teoInferenceManualInstanceConfigAttrTypes, map[string]attr.Value{
		"fixed_instance_count": types.Int64PointerValue(config.FixedInstanceCount),
	})
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	list, d := types.ListValue(objType, []attr.Value{obj})
	diags.Append(d...)
	return list, diags
}

// flattenTeoInferenceAffinityConfig converts the affinity configuration.
func flattenTeoInferenceAffinityConfig(ctx context.Context, config *teo.InferenceAffinityConfig) (types.List, diag.Diagnostics) {
	objType := types.ObjectType{AttrTypes: teoInferenceAffinityConfigAttrTypes}
	if config == nil {
		return types.ListNull(objType), nil
	}

	var diags diag.Diagnostics

	sessionIdConfig, d := flattenTeoInferenceSessionIdAffinityConfig(ctx, config.SessionIdAffinityConfig)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	obj, d := types.ObjectValue(teoInferenceAffinityConfigAttrTypes, map[string]attr.Value{
		"switch":        types.StringPointerValue(config.Switch),
		"affinity_mode": types.StringPointerValue(config.AffinityMode),
		"session_id_affinity_config": sessionIdConfig,
	})
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	list, d := types.ListValue(objType, []attr.Value{obj})
	diags.Append(d...)
	return list, diags
}

// flattenTeoInferenceSessionIdAffinityConfig converts the session ID affinity
// configuration.
func flattenTeoInferenceSessionIdAffinityConfig(ctx context.Context, config *teo.SessionIdAffinityConfig) (types.List, diag.Diagnostics) {
	objType := types.ObjectType{AttrTypes: teoInferenceSessionIdAffinityConfigAttrTypes}
	if config == nil {
		return types.ListNull(objType), nil
	}

	obj, diags := types.ObjectValue(teoInferenceSessionIdAffinityConfigAttrTypes, map[string]attr.Value{
		"source":      types.StringPointerValue(config.Source),
		"header_name": types.StringPointerValue(config.HeaderName),
	})
	if diags.HasError() {
		return types.ListNull(objType), diags
	}

	list, d := types.ListValue(objType, []attr.Value{obj})
	diags.Append(d...)
	return list, diags
}
