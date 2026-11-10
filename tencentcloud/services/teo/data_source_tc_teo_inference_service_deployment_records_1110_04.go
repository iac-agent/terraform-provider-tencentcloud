package teo

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

var _ datasource.DataSourceWithConfigure = &TeoInferenceServiceDeploymentRecords111004DataSource{}

// NewTeoInferenceServiceDeploymentRecords111004DataSource is the factory
// referenced by tencentcloud/framework/registry.go to register this data source.
func NewTeoInferenceServiceDeploymentRecords111004DataSource() datasource.DataSource {
	return &TeoInferenceServiceDeploymentRecords111004DataSource{}
}

// TeoInferenceServiceDeploymentRecords111004DataSource implements
// datasource.DataSource for the
// tencentcloud_teo_inference_service_deployment_records_1110_04 data source.
type TeoInferenceServiceDeploymentRecords111004DataSource struct {
	fw.DataSourceWithConfigure
}

// TeoInferenceServiceDeploymentRecords111004DataSourceModel maps the schema attributes.
type TeoInferenceServiceDeploymentRecords111004DataSourceModel struct {
	ZoneId     types.String `tfsdk:"zone_id"`
	ServiceId  types.String `tfsdk:"service_id"`
	SortBy     types.String `tfsdk:"sort_by"`
	SortOrder  types.String `tfsdk:"sort_order"`
	Offset     types.Int64  `tfsdk:"offset"`
	Limit      types.Int64  `tfsdk:"limit"`
	TotalCount types.Int64  `tfsdk:"total_count"`
	RecordSet  types.List   `tfsdk:"record_set"`
}

func (d *TeoInferenceServiceDeploymentRecords111004DataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_inference_service_deployment_records_1110_04"
}

func (d *TeoInferenceServiceDeploymentRecords111004DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the deployment history records of a TEO inference service. " +
			"This data source is read-only and does not modify any inference service or its deployment history.",
		Attributes: map[string]schema.Attribute{
			// Query inputs.
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
				Description: "Sorting field. Valid values: `create-time`. Default value: `create-time`.",
			},
			"sort_order": schema.StringAttribute{
				Optional:    true,
				Description: "Sorting order. Valid values: `asc`, `desc`. Default value: `desc`.",
			},
			"offset": schema.Int64Attribute{
				Optional:    true,
				Description: "Pagination offset. Default value: 0.",
			},
			"limit": schema.Int64Attribute{
				Optional:    true,
				Description: "Pagination limit. Default value: 20, maximum value: 100.",
			},
			// Query outputs.
			"total_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Total number of deployment history records that match the query.",
			},
			"record_set": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Inference service deployment history records.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"record_id": schema.StringAttribute{
							Computed:    true,
							Description: "Deployment record ID.",
						},
						"operation": schema.StringAttribute{
							Computed: true,
							Description: "Deployment operation type. Valid values: `create`, `update`, `resume`, " +
								"`stop`.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Deployment status. Valid values: `processing`, `succeeded`, `failed`.",
						},
						"duration": schema.Int64Attribute{
							Computed:    true,
							Description: "Deployment duration in seconds.",
						},
						"inference_service_config": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Configuration of this inference service deployment.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"listen_port": schema.Int64Attribute{
										Computed:    true,
										Description: "Port that the model service listens on.",
									},
									"request_paths": schema.SetAttribute{
										ElementType: types.StringType,
										Computed:    true,
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
												"tcr_repository_config": schema.ListNestedAttribute{
													Computed: true,
													Description: "TCR image repository information. It is required when " +
														"`image_type` is `TCR`. This field may return null, which means no valid " +
														"value can be obtained.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"tcr_type": schema.StringAttribute{
																Computed: true,
																Description: "TCR service type. Valid values: `Personal`, " +
																	"`Enterprise`.",
															},
															"image": schema.StringAttribute{
																Computed:    true,
																Description: "Image address.",
															},
															"registry_id": schema.StringAttribute{
																Computed: true,
																Description: "Image repository instance ID. It is required " +
																	"when `tcr_type` is `Enterprise`.",
															},
															"region_name": schema.StringAttribute{
																Computed:    true,
																Description: "Region name.",
															},
														},
													},
												},
												"startup_command": schema.StringAttribute{
													Computed: true,
													Description: "Command executed when the container starts. If it is " +
														"not specified, the Entrypoint/CMD of the image is used. The maximum " +
														"length is 1024 characters.",
												},
												"environment_variables": schema.ListNestedAttribute{
													Computed: true,
													Description: "Environment variables of the container. " +
														"A maximum of 50 variables are supported.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"key": schema.StringAttribute{
																Computed: true,
																Description: "Variable name. Only letters, digits and " +
																	"underscores are allowed, and it must start with a letter or " +
																	"an underscore. The maximum length is 64 characters.",
															},
															"value": schema.StringAttribute{
																Computed: true,
																Description: "Variable value. Any visible character is " +
																	"supported. The maximum length is 2048 characters.",
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
													Computed: true,
													Description: "Scaling mode. Valid values: `Auto` (adjust the " +
														"number of instances automatically according to the request volume), " +
														"`Manual` (use a fixed number of instances).",
												},
												"hardware_spec": schema.StringAttribute{
													Computed: true,
													Description: "Hardware specification identifier. This field is " +
														"deprecated, use `hardware_spec_id` instead.",
												},
												"hardware_spec_id": schema.StringAttribute{
													Computed: true,
													Description: "Unique ID of the hardware specification. It can be " +
														"obtained through the DescribeInferenceHardwareSpecifications API.",
												},
												"hardware_config": schema.ListNestedAttribute{
													Computed: true,
													Description: "Hardware configuration of the inference service. If it " +
														"is not specified, the default values of the selected " +
														"`hardware_spec_id` are used.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"gpu_num": schema.Float64Attribute{
																Computed: true,
																Description: "Number of GPUs allocated to a single " +
																	"instance of the inference service.",
															},
															"cpu_num": schema.Float64Attribute{
																Computed: true,
																Description: "Number of CPU cores allocated to a " +
																	"single instance of the inference service.",
															},
															"mem_size": schema.Int64Attribute{
																Computed: true,
																Description: "Memory size allocated to a single " +
																	"instance of the inference service, in MB.",
															},
															"disk_size": schema.Int64Attribute{
																Computed: true,
																Description: "Temporary disk size allocated to a " +
																	"single instance of the inference service, in MB.",
															},
														},
													},
												},
												"auto_scaling_config": schema.ListNestedAttribute{
													Computed: true,
													Description: "Auto scaling configuration of the inference " +
														"service. It is required when `scaling_mode` is `Auto`. This field " +
														"may return null, which means no valid value can be obtained.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"min_instance_count": schema.Int64Attribute{
																Computed: true,
																Description: "Minimum number of instances. It does " +
																	"not take effect while a scaling policy is in its " +
																	"validity period.",
															},
															"scaling_policies": schema.ListNestedAttribute{
																Computed: true,
																Description: "Scaling policy list. A maximum of 5 " +
																	"policies are supported.",
																NestedObject: schema.NestedAttributeObject{
																	Attributes: map[string]schema.Attribute{
																		"policy_name": schema.StringAttribute{
																			Computed: true,
																			Description: "Policy name. The length " +
																				"must be 1 to 30 characters. Policy names " +
																				"must be unique within the same service.",
																		},
																		"policy_type": schema.StringAttribute{
																			Computed: true,
																			Description: "Policy type. It cannot be " +
																				"modified after creation. Valid values: " +
																				"`ScheduledScaling`.",
																		},
																		"scheduled_scaling_policy": schema.ListNestedAttribute{
																			Computed: true,
																			Description: "Scheduled scaling " +
																				"configuration. It is required when " +
																				"`policy_type` is `ScheduledScaling`.",
																			NestedObject: schema.NestedAttributeObject{
																				Attributes: map[string]schema.Attribute{
																					"scheduled_actions": schema.ListNestedAttribute{
																						Computed: true,
																						Description: "Scheduled scaling " +
																							"action list. At least 1 and at " +
																							"most 10 actions are supported.",
																						NestedObject: schema.NestedAttributeObject{
																							Attributes: map[string]schema.Attribute{
																								"cron_expression": schema.StringAttribute{
																									Computed: true,
																									Description: "Cron " +
																										"expression describing " +
																										"when the scheduled " +
																										"scaling action is " +
																										"triggered. Standard " +
																										"5-field cron format: " +
																										"minute hour day month " +
																										"weekday. Second and " +
																										"year fields are not " +
																										"supported.",
																								},
																								"min_instance_count": schema.Int64Attribute{
																									Computed: true,
																									Description: "Minimum " +
																										"number of instances " +
																										"the inference service " +
																										"is scaled to when this " +
																										"action is hit.",
																								},
																							},
																						},
																					},
																					"effective_range": schema.ListNestedAttribute{
																						Computed: true,
																						Description: "Validity range " +
																							"of the scheduled scaling policy.",
																						NestedObject: schema.NestedAttributeObject{
																							Attributes: map[string]schema.Attribute{
																								"effective_type": schema.StringAttribute{
																									Computed: true,
																									Description: "Validity type. " +
																										"Valid values: `LongTerm` " +
																										"(permanently valid), " +
																										"`Custom` (custom start " +
																										"and end dates).",
																								},
																								"start_date": schema.StringAttribute{
																									Computed: true,
																									Description: "Validity " +
																										"start date. It is " +
																										"required when " +
																										"`effective_type` is " +
																										"`Custom`.",
																								},
																								"end_date": schema.StringAttribute{
																									Computed: true,
																									Description: "Validity " +
																										"end date. It is " +
																										"required when " +
																										"`effective_type` is " +
																										"`Custom` and must " +
																										"not be earlier than " +
																										"`start_date`.",
																								},
																							},
																						},
																					},
																					"time_zone": schema.StringAttribute{
																						Computed: true,
																						Description: "IANA time zone used " +
																							"to interpret the trigger time " +
																							"in `scheduled_actions`, for " +
																							"example `UTC`, `Asia/Shanghai`. " +
																							"UTC is used by default.",
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
													Computed: true,
													Description: "Manual instance configuration of the inference " +
														"service. It is required when `scaling_mode` is `Manual`. This " +
														"field may return null, which means no valid value can be obtained.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"fixed_instance_count": schema.Int64Attribute{
																Computed:    true,
																Description: "Fixed number of instances.",
															},
														},
													},
												},
												"concurrency": schema.Int64Attribute{
													Computed:    true,
													Description: "Concurrency of a single instance. Default value: 1.",
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
													Computed: true,
													Description: "Master switch of the inference service affinity. " +
														"Valid values: `On`, `Off`.",
												},
												"affinity_mode": schema.StringAttribute{
													Computed: true,
													Description: "Affinity mode of the inference service. Valid " +
														"values: `SessionId`. Default value: `SessionId`.",
												},
												"session_id_affinity_config": schema.ListNestedAttribute{
													Computed: true,
													Description: "Affinity configuration of the inference service. " +
														"It is required when `affinity_mode` is `SessionId`.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"source": schema.StringAttribute{
																Computed: true,
																Description: "Position where the session ID " +
																	"parameter is passed. Valid values: `Header`. Default " +
																	"value: `Header`.",
															},
															"header_name": schema.StringAttribute{
																Computed: true,
																Description: "Name of the request header that " +
																	"carries the session ID. It is required when `source` " +
																	"is `Header`. Default value: `EO-Infer-Session-Id`.",
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
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "Deployment creation time in ISO date format.",
						},
						"active_status": schema.StringAttribute{
							Computed: true,
							Description: "Whether this deployment configuration is the currently effective one. " +
								"Valid values: `active`, `inactive`.",
						},
					},
				},
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *TeoInferenceServiceDeploymentRecords111004DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TeoInferenceServiceDeploymentRecords111004DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := teov20220901.NewDescribeInferenceServiceDeploymentRecordsRequest()
	request.ZoneId = config.ZoneId.ValueStringPointer()
	request.ServiceId = config.ServiceId.ValueStringPointer()
	if !config.SortBy.IsNull() && !config.SortBy.IsUnknown() {
		request.SortBy = config.SortBy.ValueStringPointer()
	}
	if !config.SortOrder.IsNull() && !config.SortOrder.IsUnknown() {
		request.SortOrder = config.SortOrder.ValueStringPointer()
	}
	if !config.Offset.IsNull() && !config.Offset.IsUnknown() {
		request.Offset = config.Offset.ValueInt64Pointer()
	}
	if !config.Limit.IsNull() && !config.Limit.IsUnknown() {
		request.Limit = config.Limit.ValueInt64Pointer()
	}

	client := d.Client()
	if client == nil {
		resp.Diagnostics.AddError(
			"Provider is not configured",
			"The TencentCloud client has not been configured for this data source.",
		)
		return
	}

	response, err := client.UseTeoClient().DescribeInferenceServiceDeploymentRecordsWithContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying TEO inference service deployment records",
			fmt.Sprintf("DescribeInferenceServiceDeploymentRecords failed: %v", err),
		)
		return
	}

	// A nil response means the API call did not return any payload, which is a
	// failure rather than an empty result set.
	if response == nil || response.Response == nil {
		resp.Diagnostics.AddError(
			"Error querying TEO inference service deployment records",
			"DescribeInferenceServiceDeploymentRecords returned an empty response.",
		)
		return
	}

	records := make([]TeoInferenceServiceDeploymentRecordState, 0, len(response.Response.RecordSet))
	for _, record := range response.Response.RecordSet {
		if record == nil {
			continue
		}

		records = append(records, TeoInferenceServiceDeploymentRecordState{
			RecordId:               types.StringPointerValue(record.RecordId),
			Operation:              types.StringPointerValue(record.Operation),
			Status:                 types.StringPointerValue(record.Status),
			Duration:               types.Int64PointerValue(record.Duration),
			InferenceServiceConfig: flattenTeoInferenceServiceConfig(ctx, record.InferenceServiceConfig, resp),
			CreateTime:             types.StringPointerValue(record.CreateTime),
			ActiveStatus:           types.StringPointerValue(record.ActiveStatus),
		})
	}

	recordSet, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: teoInferenceServiceDeploymentRecordAttrTypes(),
	}, records)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.TotalCount = types.Int64PointerValue(response.Response.TotalCount)
	config.RecordSet = recordSet
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func teoInferenceServiceDeploymentRecordAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"record_id":                types.StringType,
		"operation":                types.StringType,
		"status":                   types.StringType,
		"duration":                 types.Int64Type,
		"inference_service_config": types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceServiceConfigAttrTypes()}},
		"create_time":              types.StringType,
		"active_status":            types.StringType,
	}
}

func teoInferenceServiceConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"listen_port":     types.Int64Type,
		"request_paths":   types.SetType{ElemType: types.StringType},
		"containers":      types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceContainerConfigAttrTypes()}},
		"resource_config": types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceResourceConfigAttrTypes()}},
		"affinity_config": types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceAffinityConfigAttrTypes()}},
	}
}

func teoInferenceContainerConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"image_type":            types.StringType,
		"tcr_repository_config": types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceTcrRepositoryConfigAttrTypes()}},
		"startup_command":       types.StringType,
		"environment_variables": types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceEnvironmentVariableAttrTypes()}},
	}
}

func teoInferenceTcrRepositoryConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"tcr_type":    types.StringType,
		"image":       types.StringType,
		"registry_id": types.StringType,
		"region_name": types.StringType,
	}
}

func teoInferenceEnvironmentVariableAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	}
}

func teoInferenceResourceConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"scaling_mode":           types.StringType,
		"hardware_spec":          types.StringType,
		"hardware_spec_id":       types.StringType,
		"hardware_config":        types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceHardwareConfigAttrTypes()}},
		"auto_scaling_config":    types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceAutoScalingConfigAttrTypes()}},
		"manual_instance_config": types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceManualInstanceConfigAttrTypes()}},
		"concurrency":            types.Int64Type,
	}
}

func teoInferenceHardwareConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"gpu_num":   types.Float64Type,
		"cpu_num":   types.Float64Type,
		"mem_size":  types.Int64Type,
		"disk_size": types.Int64Type,
	}
}

func teoInferenceAutoScalingConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"min_instance_count": types.Int64Type,
		"scaling_policies":   types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceScalingPolicyAttrTypes()}},
	}
}

func teoInferenceScalingPolicyAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"policy_name":              types.StringType,
		"policy_type":              types.StringType,
		"scheduled_scaling_policy": types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceScheduledScalingPolicyAttrTypes()}},
	}
}

func teoInferenceScheduledScalingPolicyAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"scheduled_actions": types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceScheduledScalingActionAttrTypes()}},
		"effective_range":   types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceScheduledScalingEffectiveRangeAttrTypes()}},
		"time_zone":         types.StringType,
	}
}

func teoInferenceScheduledScalingActionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cron_expression":    types.StringType,
		"min_instance_count": types.Int64Type,
	}
}

func teoInferenceScheduledScalingEffectiveRangeAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"effective_type": types.StringType,
		"start_date":     types.StringType,
		"end_date":       types.StringType,
	}
}

func teoInferenceManualInstanceConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"fixed_instance_count": types.Int64Type,
	}
}

func teoInferenceAffinityConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"switch":                     types.StringType,
		"affinity_mode":              types.StringType,
		"session_id_affinity_config": types.ListType{ElemType: types.ObjectType{AttrTypes: teoSessionIdAffinityConfigAttrTypes()}},
	}
}

func teoSessionIdAffinityConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"source":      types.StringType,
		"header_name": types.StringType,
	}
}

// The types below mirror the schema hierarchy and are exported so that tests
// can decode resp.State into strongly typed values.

// TeoInferenceServiceDeploymentRecordState is the state model of one record_set element.
type TeoInferenceServiceDeploymentRecordState struct {
	RecordId               types.String `tfsdk:"record_id"`
	Operation              types.String `tfsdk:"operation"`
	Status                 types.String `tfsdk:"status"`
	Duration               types.Int64  `tfsdk:"duration"`
	InferenceServiceConfig types.List   `tfsdk:"inference_service_config"`
	CreateTime             types.String `tfsdk:"create_time"`
	ActiveStatus           types.String `tfsdk:"active_status"`
}

// TeoInferenceServiceConfigState is the state model of one inference_service_config element.
type TeoInferenceServiceConfigState struct {
	ListenPort     types.Int64 `tfsdk:"listen_port"`
	RequestPaths   types.Set   `tfsdk:"request_paths"`
	Containers     types.List  `tfsdk:"containers"`
	ResourceConfig types.List  `tfsdk:"resource_config"`
	AffinityConfig types.List  `tfsdk:"affinity_config"`
}

// TeoInferenceContainerConfigState is the state model of one containers element.
type TeoInferenceContainerConfigState struct {
	ImageType            types.String `tfsdk:"image_type"`
	TcrRepositoryConfig  types.List   `tfsdk:"tcr_repository_config"`
	StartupCommand       types.String `tfsdk:"startup_command"`
	EnvironmentVariables types.List   `tfsdk:"environment_variables"`
}

// TeoInferenceTcrRepositoryConfigState is the state model of one tcr_repository_config element.
type TeoInferenceTcrRepositoryConfigState struct {
	TcrType    types.String `tfsdk:"tcr_type"`
	Image      types.String `tfsdk:"image"`
	RegistryId types.String `tfsdk:"registry_id"`
	RegionName types.String `tfsdk:"region_name"`
}

// TeoInferenceEnvironmentVariableState is the state model of one environment_variables element.
type TeoInferenceEnvironmentVariableState struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// TeoInferenceResourceConfigState is the state model of one resource_config element.
type TeoInferenceResourceConfigState struct {
	ScalingMode          types.String `tfsdk:"scaling_mode"`
	HardwareSpec         types.String `tfsdk:"hardware_spec"`
	HardwareSpecId       types.String `tfsdk:"hardware_spec_id"`
	HardwareConfig       types.List   `tfsdk:"hardware_config"`
	AutoScalingConfig    types.List   `tfsdk:"auto_scaling_config"`
	ManualInstanceConfig types.List   `tfsdk:"manual_instance_config"`
	Concurrency          types.Int64  `tfsdk:"concurrency"`
}

// TeoInferenceHardwareConfigState is the state model of one hardware_config element.
type TeoInferenceHardwareConfigState struct {
	GpuNum   types.Float64 `tfsdk:"gpu_num"`
	CpuNum   types.Float64 `tfsdk:"cpu_num"`
	MemSize  types.Int64   `tfsdk:"mem_size"`
	DiskSize types.Int64   `tfsdk:"disk_size"`
}

// TeoInferenceAutoScalingConfigState is the state model of one auto_scaling_config element.
type TeoInferenceAutoScalingConfigState struct {
	MinInstanceCount types.Int64 `tfsdk:"min_instance_count"`
	ScalingPolicies  types.List  `tfsdk:"scaling_policies"`
}

// TeoInferenceScalingPolicyState is the state model of one scaling_policies element.
type TeoInferenceScalingPolicyState struct {
	PolicyName             types.String `tfsdk:"policy_name"`
	PolicyType             types.String `tfsdk:"policy_type"`
	ScheduledScalingPolicy types.List   `tfsdk:"scheduled_scaling_policy"`
}

// TeoInferenceScheduledScalingPolicyState is the state model of one scheduled_scaling_policy element.
type TeoInferenceScheduledScalingPolicyState struct {
	ScheduledActions types.List   `tfsdk:"scheduled_actions"`
	EffectiveRange   types.List   `tfsdk:"effective_range"`
	TimeZone         types.String `tfsdk:"time_zone"`
}

// TeoInferenceScheduledScalingActionState is the state model of one scheduled_actions element.
type TeoInferenceScheduledScalingActionState struct {
	CronExpression   types.String `tfsdk:"cron_expression"`
	MinInstanceCount types.Int64  `tfsdk:"min_instance_count"`
}

// TeoInferenceScheduledScalingEffectiveRangeState is the state model of one effective_range element.
type TeoInferenceScheduledScalingEffectiveRangeState struct {
	EffectiveType types.String `tfsdk:"effective_type"`
	StartDate     types.String `tfsdk:"start_date"`
	EndDate       types.String `tfsdk:"end_date"`
}

// TeoInferenceManualInstanceConfigState is the state model of one manual_instance_config element.
type TeoInferenceManualInstanceConfigState struct {
	FixedInstanceCount types.Int64 `tfsdk:"fixed_instance_count"`
}

// TeoInferenceAffinityConfigState is the state model of one affinity_config element.
type TeoInferenceAffinityConfigState struct {
	Switch                  types.String `tfsdk:"switch"`
	AffinityMode            types.String `tfsdk:"affinity_mode"`
	SessionIdAffinityConfig types.List   `tfsdk:"session_id_affinity_config"`
}

// TeoSessionIdAffinityConfigState is the state model of one session_id_affinity_config element.
type TeoSessionIdAffinityConfigState struct {
	Source     types.String `tfsdk:"source"`
	HeaderName types.String `tfsdk:"header_name"`
}

// flattenTeoInferenceServiceConfig converts the SDK model of a deployment
// configuration into a single-element list. A nil model is represented as a
// null list so that "not returned by the cloud" stays distinguishable from
// "returned but empty".
func flattenTeoInferenceServiceConfig(ctx context.Context, config *teov20220901.InferenceServiceConfig, resp *datasource.ReadResponse) types.List {
	elemType := types.ObjectType{AttrTypes: teoInferenceServiceConfigAttrTypes()}
	if config == nil {
		return types.ListNull(elemType)
	}

	requestPaths := make([]attr.Value, 0, len(config.RequestPaths))
	for _, path := range config.RequestPaths {
		if path == nil {
			continue
		}
		requestPaths = append(requestPaths, types.StringValue(*path))
	}
	requestPathSet, diags := types.SetValue(types.StringType, requestPaths)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return types.ListNull(elemType)
	}

	containers := make([]TeoInferenceContainerConfigState, 0, len(config.Containers))
	for _, container := range config.Containers {
		if container == nil {
			continue
		}

		tcrConfig := types.ListNull(types.ObjectType{AttrTypes: teoInferenceTcrRepositoryConfigAttrTypes()})
		if container.TcrRepositoryConfig != nil {
			value, d := types.ListValueFrom(ctx, types.ObjectType{
				AttrTypes: teoInferenceTcrRepositoryConfigAttrTypes(),
			}, []TeoInferenceTcrRepositoryConfigState{
				{
					TcrType:    types.StringPointerValue(container.TcrRepositoryConfig.TCRType),
					Image:      types.StringPointerValue(container.TcrRepositoryConfig.Image),
					RegistryId: types.StringPointerValue(container.TcrRepositoryConfig.RegistryId),
					RegionName: types.StringPointerValue(container.TcrRepositoryConfig.RegionName),
				},
			})
			resp.Diagnostics.Append(d...)
			if resp.Diagnostics.HasError() {
				return types.ListNull(elemType)
			}
			tcrConfig = value
		}

		envVariables := make([]TeoInferenceEnvironmentVariableState, 0, len(container.EnvironmentVariables))
		for _, env := range container.EnvironmentVariables {
			if env == nil {
				continue
			}
			envVariables = append(envVariables, TeoInferenceEnvironmentVariableState{
				Key:   types.StringPointerValue(env.Key),
				Value: types.StringPointerValue(env.Value),
			})
		}
		envList, d := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: teoInferenceEnvironmentVariableAttrTypes(),
		}, envVariables)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return types.ListNull(elemType)
		}

		containers = append(containers, TeoInferenceContainerConfigState{
			ImageType:            types.StringPointerValue(container.ImageType),
			TcrRepositoryConfig:  tcrConfig,
			StartupCommand:       types.StringPointerValue(container.StartupCommand),
			EnvironmentVariables: envList,
		})
	}
	containerList, d := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: teoInferenceContainerConfigAttrTypes(),
	}, containers)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return types.ListNull(elemType)
	}

	resourceConfig := types.ListNull(types.ObjectType{AttrTypes: teoInferenceResourceConfigAttrTypes()})
	if config.ResourceConfig != nil {
		value, d := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: teoInferenceResourceConfigAttrTypes(),
		}, []TeoInferenceResourceConfigState{flattenTeoInferenceResourceConfig(ctx, config.ResourceConfig, resp)})
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return types.ListNull(elemType)
		}
		resourceConfig = value
	}

	affinityConfig := types.ListNull(types.ObjectType{AttrTypes: teoInferenceAffinityConfigAttrTypes()})
	if config.AffinityConfig != nil {
		value, d := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: teoInferenceAffinityConfigAttrTypes(),
		}, []TeoInferenceAffinityConfigState{flattenTeoInferenceAffinityConfig(ctx, config.AffinityConfig, resp)})
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return types.ListNull(elemType)
		}
		affinityConfig = value
	}

	result, d := types.ListValueFrom(ctx, elemType, []TeoInferenceServiceConfigState{
		{
			ListenPort:     types.Int64PointerValue(config.ListenPort),
			RequestPaths:   requestPathSet,
			Containers:     containerList,
			ResourceConfig: resourceConfig,
			AffinityConfig: affinityConfig,
		},
	})
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return types.ListNull(elemType)
	}
	return result
}

func flattenTeoInferenceResourceConfig(ctx context.Context, config *teov20220901.InferenceResourceConfig, resp *datasource.ReadResponse) TeoInferenceResourceConfigState {
	hardwareConfig := types.ListNull(types.ObjectType{AttrTypes: teoInferenceHardwareConfigAttrTypes()})
	if config.HardwareConfig != nil {
		value, d := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: teoInferenceHardwareConfigAttrTypes(),
		}, []TeoInferenceHardwareConfigState{
			{
				GpuNum:   types.Float64PointerValue(config.HardwareConfig.GPUNum),
				CpuNum:   types.Float64PointerValue(config.HardwareConfig.CPUNum),
				MemSize:  types.Int64PointerValue(config.HardwareConfig.MemSize),
				DiskSize: types.Int64PointerValue(config.HardwareConfig.DiskSize),
			},
		})
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return TeoInferenceResourceConfigState{}
		}
		hardwareConfig = value
	}

	autoScalingConfig := types.ListNull(types.ObjectType{AttrTypes: teoInferenceAutoScalingConfigAttrTypes()})
	if config.AutoScalingConfig != nil {
		value, d := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: teoInferenceAutoScalingConfigAttrTypes(),
		}, []TeoInferenceAutoScalingConfigState{
			flattenTeoInferenceAutoScalingConfig(ctx, config.AutoScalingConfig, resp),
		})
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return TeoInferenceResourceConfigState{}
		}
		autoScalingConfig = value
	}

	manualInstanceConfig := types.ListNull(types.ObjectType{AttrTypes: teoInferenceManualInstanceConfigAttrTypes()})
	if config.ManualInstanceConfig != nil {
		value, d := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: teoInferenceManualInstanceConfigAttrTypes(),
		}, []TeoInferenceManualInstanceConfigState{
			{FixedInstanceCount: types.Int64PointerValue(config.ManualInstanceConfig.FixedInstanceCount)},
		})
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return TeoInferenceResourceConfigState{}
		}
		manualInstanceConfig = value
	}

	return TeoInferenceResourceConfigState{
		ScalingMode:          types.StringPointerValue(config.ScalingMode),
		HardwareSpec:         types.StringPointerValue(config.HardwareSpec),
		HardwareSpecId:       types.StringPointerValue(config.HardwareSpecId),
		HardwareConfig:       hardwareConfig,
		AutoScalingConfig:    autoScalingConfig,
		ManualInstanceConfig: manualInstanceConfig,
		Concurrency:          types.Int64PointerValue(config.Concurrency),
	}
}

func flattenTeoInferenceAutoScalingConfig(ctx context.Context, config *teov20220901.InferenceAutoScalingConfig, resp *datasource.ReadResponse) TeoInferenceAutoScalingConfigState {
	policies := make([]TeoInferenceScalingPolicyState, 0, len(config.ScalingPolicies))
	for _, policy := range config.ScalingPolicies {
		if policy == nil {
			continue
		}
		policies = append(policies, flattenTeoInferenceScalingPolicy(ctx, policy, resp))
	}
	policyList, d := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: teoInferenceScalingPolicyAttrTypes(),
	}, policies)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return TeoInferenceAutoScalingConfigState{}
	}

	return TeoInferenceAutoScalingConfigState{
		MinInstanceCount: types.Int64PointerValue(config.MinInstanceCount),
		ScalingPolicies:  policyList,
	}
}

func flattenTeoInferenceScalingPolicy(ctx context.Context, policy *teov20220901.InferenceScalingPolicy, resp *datasource.ReadResponse) TeoInferenceScalingPolicyState {
	scheduledPolicy := types.ListNull(types.ObjectType{AttrTypes: teoInferenceScheduledScalingPolicyAttrTypes()})
	if policy.ScheduledScalingPolicy != nil {
		value, d := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: teoInferenceScheduledScalingPolicyAttrTypes(),
		}, []TeoInferenceScheduledScalingPolicyState{
			flattenTeoInferenceScheduledScalingPolicy(ctx, policy.ScheduledScalingPolicy, resp),
		})
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return TeoInferenceScalingPolicyState{}
		}
		scheduledPolicy = value
	}

	return TeoInferenceScalingPolicyState{
		PolicyName:             types.StringPointerValue(policy.PolicyName),
		PolicyType:             types.StringPointerValue(policy.PolicyType),
		ScheduledScalingPolicy: scheduledPolicy,
	}
}

func flattenTeoInferenceScheduledScalingPolicy(ctx context.Context, policy *teov20220901.InferenceScheduledScalingPolicy, resp *datasource.ReadResponse) TeoInferenceScheduledScalingPolicyState {
	actions := make([]TeoInferenceScheduledScalingActionState, 0, len(policy.ScheduledActions))
	for _, action := range policy.ScheduledActions {
		if action == nil {
			continue
		}
		actions = append(actions, TeoInferenceScheduledScalingActionState{
			CronExpression:   types.StringPointerValue(action.CronExpression),
			MinInstanceCount: types.Int64PointerValue(action.MinInstanceCount),
		})
	}
	actionList, d := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: teoInferenceScheduledScalingActionAttrTypes(),
	}, actions)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return TeoInferenceScheduledScalingPolicyState{}
	}

	effectiveRange := types.ListNull(types.ObjectType{AttrTypes: teoInferenceScheduledScalingEffectiveRangeAttrTypes()})
	if policy.EffectiveRange != nil {
		value, d := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: teoInferenceScheduledScalingEffectiveRangeAttrTypes(),
		}, []TeoInferenceScheduledScalingEffectiveRangeState{
			{
				EffectiveType: types.StringPointerValue(policy.EffectiveRange.EffectiveType),
				StartDate:     types.StringPointerValue(policy.EffectiveRange.StartDate),
				EndDate:       types.StringPointerValue(policy.EffectiveRange.EndDate),
			},
		})
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return TeoInferenceScheduledScalingPolicyState{}
		}
		effectiveRange = value
	}

	return TeoInferenceScheduledScalingPolicyState{
		ScheduledActions: actionList,
		EffectiveRange:   effectiveRange,
		TimeZone:         types.StringPointerValue(policy.TimeZone),
	}
}

func flattenTeoInferenceAffinityConfig(ctx context.Context, config *teov20220901.InferenceAffinityConfig, resp *datasource.ReadResponse) TeoInferenceAffinityConfigState {
	sessionConfig := types.ListNull(types.ObjectType{AttrTypes: teoSessionIdAffinityConfigAttrTypes()})
	if config.SessionIdAffinityConfig != nil {
		value, d := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: teoSessionIdAffinityConfigAttrTypes(),
		}, []TeoSessionIdAffinityConfigState{
			{
				Source:     types.StringPointerValue(config.SessionIdAffinityConfig.Source),
				HeaderName: types.StringPointerValue(config.SessionIdAffinityConfig.HeaderName),
			},
		})
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return TeoInferenceAffinityConfigState{}
		}
		sessionConfig = value
	}

	return TeoInferenceAffinityConfigState{
		Switch:                  types.StringPointerValue(config.Switch),
		AffinityMode:            types.StringPointerValue(config.AffinityMode),
		SessionIdAffinityConfig: sessionConfig,
	}
}
