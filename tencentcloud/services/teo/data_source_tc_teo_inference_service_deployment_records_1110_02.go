package teo

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// _ datasource.DataSource is intentionally omitted: the factory function
// NewTeoInferenceServiceDeploymentRecords111002DataSource returns
// datasource.DataSource, which already performs the implicit compile-time
// check at the return statement. Only the extended-interface assertion is kept
// here to verify the Configure method signature promoted from
// fw.DataSourceWithConfigure.
var _ datasource.DataSourceWithConfigure = &TeoInferenceServiceDeploymentRecords111002DataSource{}

// NewTeoInferenceServiceDeploymentRecords111002DataSource is the factory
// referenced by tencentcloud/framework/registry.go to register this data
// source.
func NewTeoInferenceServiceDeploymentRecords111002DataSource() datasource.DataSource {
	return &TeoInferenceServiceDeploymentRecords111002DataSource{}
}

// TeoInferenceServiceDeploymentRecords111002DataSource implements
// datasource.DataSource for the
// tencentcloud_teo_inference_service_deployment_records_1110_02 data source.
type TeoInferenceServiceDeploymentRecords111002DataSource struct {
	fw.DataSourceWithConfigure
}

// TeoInferenceServiceDeploymentRecords111002DataSourceModel maps the schema
// attributes.
type TeoInferenceServiceDeploymentRecords111002DataSourceModel struct {
	ZoneId    types.String `tfsdk:"zone_id"`
	ServiceId types.String `tfsdk:"service_id"`
	SortBy    types.String `tfsdk:"sort_by"`
	SortOrder types.String `tfsdk:"sort_order"`
	RecordSet types.List   `tfsdk:"record_set"`
}

// teoInferenceDeploymentRecordModel is the element type of record_set.
type teoInferenceDeploymentRecordModel struct {
	RecordId               types.String `tfsdk:"record_id"`
	Operation              types.String `tfsdk:"operation"`
	Status                 types.String `tfsdk:"status"`
	Duration               types.Int64  `tfsdk:"duration"`
	InferenceServiceConfig types.List   `tfsdk:"inference_service_config"`
	CreateTime             types.String `tfsdk:"create_time"`
	ActiveStatus           types.String `tfsdk:"active_status"`
}

// teoInferenceServiceConfigModel is the element type of
// record_set.inference_service_config.
type teoInferenceServiceConfigModel struct {
	ListenPort     types.Int64 `tfsdk:"listen_port"`
	RequestPaths   types.Set   `tfsdk:"request_paths"`
	Containers     types.List  `tfsdk:"containers"`
	ResourceConfig types.List  `tfsdk:"resource_config"`
	AffinityConfig types.List  `tfsdk:"affinity_config"`
}

// teoInferenceContainerConfigModel is the element type of
// record_set.inference_service_config.containers.
type teoInferenceContainerConfigModel struct {
	ImageType            types.String `tfsdk:"image_type"`
	TcrRepositoryConfig  types.List   `tfsdk:"tcr_repository_config"`
	StartupCommand       types.String `tfsdk:"startup_command"`
	EnvironmentVariables types.List   `tfsdk:"environment_variables"`
}

// teoInferenceTcrRepositoryConfigModel is the element type of
// record_set.inference_service_config.containers.tcr_repository_config.
type teoInferenceTcrRepositoryConfigModel struct {
	TcrType    types.String `tfsdk:"tcr_type"`
	Image      types.String `tfsdk:"image"`
	RegistryId types.String `tfsdk:"registry_id"`
	RegionName types.String `tfsdk:"region_name"`
}

// teoInferenceEnvironmentVariableModel is the element type of
// record_set.inference_service_config.containers.environment_variables.
type teoInferenceEnvironmentVariableModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// teoInferenceResourceConfigModel is the element type of
// record_set.inference_service_config.resource_config.
type teoInferenceResourceConfigModel struct {
	ScalingMode          types.String `tfsdk:"scaling_mode"`
	HardwareSpec         types.String `tfsdk:"hardware_spec"`
	HardwareSpecId       types.String `tfsdk:"hardware_spec_id"`
	HardwareConfig       types.List   `tfsdk:"hardware_config"`
	AutoScalingConfig    types.List   `tfsdk:"auto_scaling_config"`
	ManualInstanceConfig types.List   `tfsdk:"manual_instance_config"`
	Concurrency          types.Int64  `tfsdk:"concurrency"`
}

// teoInferenceHardwareConfigModel is the element type of
// record_set.inference_service_config.resource_config.hardware_config.
type teoInferenceHardwareConfigModel struct {
	GpuNum   types.Float64 `tfsdk:"gpu_num"`
	CpuNum   types.Float64 `tfsdk:"cpu_num"`
	MemSize  types.Int64   `tfsdk:"mem_size"`
	DiskSize types.Int64   `tfsdk:"disk_size"`
}

// teoInferenceAutoScalingConfigModel is the element type of
// record_set.inference_service_config.resource_config.auto_scaling_config.
type teoInferenceAutoScalingConfigModel struct {
	MinInstanceCount types.Int64 `tfsdk:"min_instance_count"`
	ScalingPolicies  types.List  `tfsdk:"scaling_policies"`
}

// teoInferenceScalingPolicyModel is the element type of
// ...auto_scaling_config.scaling_policies.
type teoInferenceScalingPolicyModel struct {
	PolicyName             types.String `tfsdk:"policy_name"`
	PolicyType             types.String `tfsdk:"policy_type"`
	ScheduledScalingPolicy types.List   `tfsdk:"scheduled_scaling_policy"`
}

// teoInferenceScheduledScalingPolicyModel is the element type of
// ...scaling_policies.scheduled_scaling_policy.
type teoInferenceScheduledScalingPolicyModel struct {
	ScheduledActions types.List   `tfsdk:"scheduled_actions"`
	EffectiveRange   types.List   `tfsdk:"effective_range"`
	TimeZone         types.String `tfsdk:"time_zone"`
}

// teoInferenceScheduledScalingActionModel is the element type of
// ...scheduled_scaling_policy.scheduled_actions.
type teoInferenceScheduledScalingActionModel struct {
	CronExpression   types.String `tfsdk:"cron_expression"`
	MinInstanceCount types.Int64  `tfsdk:"min_instance_count"`
}

// teoInferenceScheduledScalingEffectiveRangeModel is the element type of
// ...scheduled_scaling_policy.effective_range.
type teoInferenceScheduledScalingEffectiveRangeModel struct {
	EffectiveType types.String `tfsdk:"effective_type"`
	StartDate     types.String `tfsdk:"start_date"`
	EndDate       types.String `tfsdk:"end_date"`
}

// teoInferenceManualInstanceConfigModel is the element type of
// record_set.inference_service_config.resource_config.manual_instance_config.
type teoInferenceManualInstanceConfigModel struct {
	FixedInstanceCount types.Int64 `tfsdk:"fixed_instance_count"`
}

// teoInferenceAffinityConfigModel is the element type of
// record_set.inference_service_config.affinity_config.
type teoInferenceAffinityConfigModel struct {
	Switch                  types.String `tfsdk:"switch"`
	AffinityMode            types.String `tfsdk:"affinity_mode"`
	SessionIdAffinityConfig types.List   `tfsdk:"session_id_affinity_config"`
}

// teoInferenceSessionIdAffinityConfigModel is the element type of
// ...affinity_config.session_id_affinity_config.
type teoInferenceSessionIdAffinityConfigModel struct {
	Source     types.String `tfsdk:"source"`
	HeaderName types.String `tfsdk:"header_name"`
}

// Attribute type maps used when converting nested objects. They must stay in
// sync with the schema declared in Schema().
var (
	teoInferenceScheduledActionAttrTypes = map[string]attr.Type{
		"cron_expression":    types.StringType,
		"min_instance_count": types.Int64Type,
	}

	teoInferenceScheduledEffectiveRangeAttrTypes = map[string]attr.Type{
		"effective_type": types.StringType,
		"start_date":     types.StringType,
		"end_date":       types.StringType,
	}

	teoInferenceScheduledPolicyAttrTypes = map[string]attr.Type{
		"scheduled_actions": types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceScheduledActionAttrTypes}},
		"effective_range":   types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceScheduledEffectiveRangeAttrTypes}},
		"time_zone":         types.StringType,
	}

	teoInferenceScalingPolicyAttrTypes = map[string]attr.Type{
		"policy_name":              types.StringType,
		"policy_type":              types.StringType,
		"scheduled_scaling_policy": types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceScheduledPolicyAttrTypes}},
	}

	teoInferenceAutoScalingAttrTypes = map[string]attr.Type{
		"min_instance_count": types.Int64Type,
		"scaling_policies":   types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceScalingPolicyAttrTypes}},
	}

	teoInferenceManualInstanceAttrTypes = map[string]attr.Type{
		"fixed_instance_count": types.Int64Type,
	}

	teoInferenceHardwareConfigAttrTypes = map[string]attr.Type{
		"gpu_num":   types.Float64Type,
		"cpu_num":   types.Float64Type,
		"mem_size":  types.Int64Type,
		"disk_size": types.Int64Type,
	}

	teoInferenceResourceConfigAttrTypes = map[string]attr.Type{
		"scaling_mode":           types.StringType,
		"hardware_spec":          types.StringType,
		"hardware_spec_id":       types.StringType,
		"hardware_config":        types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceHardwareConfigAttrTypes}},
		"auto_scaling_config":    types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceAutoScalingAttrTypes}},
		"manual_instance_config": types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceManualInstanceAttrTypes}},
		"concurrency":            types.Int64Type,
	}

	teoInferenceSessionIdAffinityAttrTypes = map[string]attr.Type{
		"source":      types.StringType,
		"header_name": types.StringType,
	}

	teoInferenceAffinityConfigAttrTypes = map[string]attr.Type{
		"switch":                     types.StringType,
		"affinity_mode":              types.StringType,
		"session_id_affinity_config": types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceSessionIdAffinityAttrTypes}},
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

	teoInferenceContainerConfigAttrTypes = map[string]attr.Type{
		"image_type":            types.StringType,
		"tcr_repository_config": types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceTcrRepositoryConfigAttrTypes}},
		"startup_command":       types.StringType,
		"environment_variables": types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceEnvironmentVariableAttrTypes}},
	}

	teoInferenceServiceConfigAttrTypes = map[string]attr.Type{
		"listen_port":     types.Int64Type,
		"request_paths":   types.SetType{ElemType: types.StringType},
		"containers":      types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceContainerConfigAttrTypes}},
		"resource_config": types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceResourceConfigAttrTypes}},
		"affinity_config": types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceAffinityConfigAttrTypes}},
	}

	teoInferenceDeploymentRecordAttrTypes = map[string]attr.Type{
		"record_id":                types.StringType,
		"operation":                types.StringType,
		"status":                   types.StringType,
		"duration":                 types.Int64Type,
		"inference_service_config": types.ListType{ElemType: types.ObjectType{AttrTypes: teoInferenceServiceConfigAttrTypes}},
		"create_time":              types.StringType,
		"active_status":            types.StringType,
	}
)

func (d *TeoInferenceServiceDeploymentRecords111002DataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_inference_service_deployment_records_1110_02"
}

func (d *TeoInferenceServiceDeploymentRecords111002DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the deployment history records of a TEO inference service.",
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
				Description: "Field used for sorting. Valid values: create-time. Default value: create-time.",
			},
			"sort_order": schema.StringAttribute{
				Optional:    true,
				Description: "Sorting order. Valid values: asc, desc. Default value: desc.",
			},
			// ---- Query outputs: Computed ----
			"record_set": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Deployment history record list of the inference service.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"record_id": schema.StringAttribute{
							Computed:    true,
							Description: "Deployment record ID.",
						},
						"operation": schema.StringAttribute{
							Computed:    true,
							Description: "Deployment operation type. Valid values: create, update, resume, stop.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Deployment status. Valid values: processing, succeeded, failed.",
						},
						"duration": schema.Int64Attribute{
							Computed:    true,
							Description: "Deployment duration in seconds.",
						},
						"inference_service_config": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Configuration used by this inference service deployment.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"listen_port": schema.Int64Attribute{
										Computed:    true,
										Description: "Port that the model service listens on.",
									},
									"request_paths": schema.SetAttribute{
										Computed:    true,
										ElementType: types.StringType,
										Description: "Request path list of the inference service.",
									},
									"containers": schema.ListNestedAttribute{
										Computed:    true,
										Description: "Container configurations of the inference service.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"image_type": schema.StringAttribute{
													Computed:    true,
													Description: "Image type. Valid values: TCR.",
												},
												"tcr_repository_config": schema.ListNestedAttribute{
													Computed:    true,
													Description: "TCR image repository information. Required when ImageType is TCR.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"tcr_type": schema.StringAttribute{
																Computed:    true,
																Description: "TCR service type. Valid values: Personal, Enterprise.",
															},
															"image": schema.StringAttribute{
																Computed:    true,
																Description: "Image address.",
															},
															"registry_id": schema.StringAttribute{
																Computed:    true,
																Description: "Image repository instance ID. Required when TCRType is Enterprise.",
															},
															"region_name": schema.StringAttribute{
																Computed:    true,
																Description: "Region name.",
															},
														},
													},
												},
												"startup_command": schema.StringAttribute{
													Computed:    true,
													Description: "Command executed when the container starts. The image Entrypoint/CMD is used when it is not set. Maximum 1024 characters.",
												},
												"environment_variables": schema.ListNestedAttribute{
													Computed:    true,
													Description: "Environment variables of the container at runtime. Up to 50 variables are supported.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"key": schema.StringAttribute{
																Computed:    true,
																Description: "Variable name. Only letters, digits and underscores are allowed, and it must start with a letter or an underscore. Maximum 64 characters.",
															},
															"value": schema.StringAttribute{
																Computed:    true,
																Description: "Variable value. Any visible characters are supported. Maximum 2048 characters.",
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
													Description: "Scaling mode. Valid values: Auto (adjust instance count automatically based on request volume), Manual (fixed instance count).",
												},
												"hardware_spec": schema.StringAttribute{
													Computed:    true,
													Description: "Hardware specification identifier. Deprecated, use HardwareSpecId instead.",
												},
												"hardware_spec_id": schema.StringAttribute{
													Computed:    true,
													Description: "Unique hardware specification ID. It can be obtained from the DescribeInferenceHardwareSpecifications API.",
												},
												"hardware_config": schema.ListNestedAttribute{
													Computed:    true,
													Description: "Hardware configuration of the inference service. When it is not specified, the default values of the selected HardwareSpecId are used.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"gpu_num": schema.Float64Attribute{
																Computed:    true,
																Description: "Number of GPUs allocated to a single instance. Only integer values are supported currently and the value must be one of AllowedGPUNums of the selected HardwareSpecId.",
															},
															"cpu_num": schema.Float64Attribute{
																Computed:    true,
																Description: "Number of CPU cores allocated to a single instance. Only integer values are supported currently.",
															},
															"mem_size": schema.Int64Attribute{
																Computed:    true,
																Description: "Memory size allocated to a single instance, in MB. It must be a multiple of 1024 when specified.",
															},
															"disk_size": schema.Int64Attribute{
																Computed:    true,
																Description: "Temporary disk size allocated to a single instance, in MB. It must be a multiple of 1024 when specified.",
															},
														},
													},
												},
												"auto_scaling_config": schema.ListNestedAttribute{
													Computed:    true,
													Description: "Auto scaling configuration of the inference service. Required when ScalingMode is Auto.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"min_instance_count": schema.Int64Attribute{
																Computed:    true,
																Description: "Minimum instance count. It does not take effect when a scaling policy is configured and within its effective period.",
															},
															"scaling_policies": schema.ListNestedAttribute{
																Computed:    true,
																Description: "Scaling policy list. Up to 5 policies are supported.",
																NestedObject: schema.NestedAttributeObject{
																	Attributes: map[string]schema.Attribute{
																		"policy_name": schema.StringAttribute{
																			Computed:    true,
																			Description: "Policy name. 1 to 30 characters. Policy names must be unique within the same service.",
																		},
																		"policy_type": schema.StringAttribute{
																			Computed:    true,
																			Description: "Policy type, which cannot be modified after creation. Valid values: ScheduledScaling.",
																		},
																		"scheduled_scaling_policy": schema.ListNestedAttribute{
																			Computed:    true,
																			Description: "Scheduled scaling configuration. Required when PolicyType is ScheduledScaling.",
																			NestedObject: schema.NestedAttributeObject{
																				Attributes: map[string]schema.Attribute{
																					"scheduled_actions": schema.ListNestedAttribute{
																						Computed:    true,
																						Description: "Scheduled scaling action list. At least 1 and at most 10 actions are supported.",
																						NestedObject: schema.NestedAttributeObject{
																							Attributes: map[string]schema.Attribute{
																								"cron_expression": schema.StringAttribute{
																									Computed:    true,
																									Description: "Cron expression describing when the scheduled scaling action is triggered. 5-field standard Cron format without seconds and years.",
																								},
																								"min_instance_count": schema.Int64Attribute{
																									Computed:    true,
																									Description: "Minimum instance count that the inference service is scaled to when this scheduled scaling action is triggered.",
																								},
																							},
																						},
																					},
																					"effective_range": schema.ListNestedAttribute{
																						Computed:    true,
																						Description: "Effective range of the scheduled scaling policy.",
																						NestedObject: schema.NestedAttributeObject{
																							Attributes: map[string]schema.Attribute{
																								"effective_type": schema.StringAttribute{
																									Computed:    true,
																									Description: "Effective range type. Valid values: LongTerm (long-term effective), Custom (custom start and end dates).",
																								},
																								"start_date": schema.StringAttribute{
																									Computed:    true,
																									Description: "Start date of the effective range. Required when EffectiveType is Custom.",
																								},
																								"end_date": schema.StringAttribute{
																									Computed:    true,
																									Description: "End date of the effective range. Required when EffectiveType is Custom and must not be earlier than StartDate.",
																								},
																							},
																						},
																					},
																					"time_zone": schema.StringAttribute{
																						Computed:    true,
																						Description: "Time zone in IANA time zone format used to interpret scheduled actions, for example UTC or Asia/Shanghai. UTC is used when it is not set.",
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
													Description: "Manually configured instance settings of the inference service. Required when ScalingMode is Manual.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"fixed_instance_count": schema.Int64Attribute{
																Computed:    true,
																Description: "Fixed instance count.",
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
													Computed:    true,
													Description: "Master switch of the inference service affinity. Valid values: On, Off.",
												},
												"affinity_mode": schema.StringAttribute{
													Computed:    true,
													Description: "Affinity mode of the inference service. Valid values: SessionId. Default value: SessionId.",
												},
												"session_id_affinity_config": schema.ListNestedAttribute{
													Computed:    true,
													Description: "Affinity configuration based on the session ID. Required when AffinityMode is SessionId.",
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"source": schema.StringAttribute{
																Computed:    true,
																Description: "Location where the session ID parameter is passed. Valid values: Header. Default value: Header.",
															},
															"header_name": schema.StringAttribute{
																Computed:    true,
																Description: "Request header name used to pass the session ID. 1 to 64 characters, only letters, digits and hyphens are allowed. Default value: EO-Infer-Session-Id.",
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
							Description: "Deployment start time in ISO date format.",
						},
						"active_status": schema.StringAttribute{
							Computed:    true,
							Description: "Whether this deployment configuration is the currently effective one. Valid values: active (currently effective), inactive (historical or abnormal version).",
						},
					},
				},
			},
		},
	}
}

// teoInferenceDeploymentRecordsPageSize is the page size used when walking the
// paginated deployment history API. 100 is the documented maximum of the Limit
// parameter, so the number of round trips stays minimal.
const teoInferenceDeploymentRecordsPageSize = int64(100)

// readDeploymentRecords aggregates every page of the deployment history and
// returns the flattened models. The API is paginated (Offset/Limit), so all
// pages are collected to keep the produced data source complete.
func (d *TeoInferenceServiceDeploymentRecords111002DataSource) readDeploymentRecords(ctx context.Context, config TeoInferenceServiceDeploymentRecords111002DataSourceModel) ([]teoInferenceDeploymentRecordModel, error) {
	records := make([]teoInferenceDeploymentRecordModel, 0)
	offset := int64(0)

	for {
		request := teov20220901.NewDescribeInferenceServiceDeploymentRecordsRequest()
		request.ZoneId = config.ZoneId.ValueStringPointer()
		request.ServiceId = config.ServiceId.ValueStringPointer()
		if !config.SortBy.IsNull() && !config.SortBy.IsUnknown() {
			request.SortBy = config.SortBy.ValueStringPointer()
		}
		if !config.SortOrder.IsNull() && !config.SortOrder.IsUnknown() {
			request.SortOrder = config.SortOrder.ValueStringPointer()
		}
		request.Offset = helper.Int64(offset)
		request.Limit = helper.Int64(teoInferenceDeploymentRecordsPageSize)

		response, err := d.Client().UseTeoClient().DescribeInferenceServiceDeploymentRecordsWithContext(ctx, request)
		if err != nil {
			return nil, fmt.Errorf("DescribeInferenceServiceDeploymentRecords failed at offset %d: %w", offset, err)
		}

		// A nil response means the API call itself failed; it must not be
		// presented as an empty result.
		if response == nil || response.Response == nil {
			return nil, fmt.Errorf("DescribeInferenceServiceDeploymentRecords returned an empty response at offset %d", offset)
		}

		page := response.Response.RecordSet
		recordDiags := diag.Diagnostics{}
		for _, record := range page {
			records = append(records, flattenTeoInferenceDeploymentRecord(ctx, record, &recordDiags))
		}
		if recordDiags.HasError() {
			return nil, fmt.Errorf("failed to flatten DescribeInferenceServiceDeploymentRecords response at offset %d: %v", offset, recordDiags)
		}

		// Stop on the last partial page, on a page that does not advance the
		// offset, or once the reported total has been collected.
		if int64(len(page)) < teoInferenceDeploymentRecordsPageSize {
			break
		}
		if response.Response.TotalCount != nil && int64(len(records)) >= *response.Response.TotalCount {
			break
		}
		offset += int64(len(page))
	}

	return records, nil
}

// Read queries the deployment history records of an inference service and
// writes them into resp.State.
func (d *TeoInferenceServiceDeploymentRecords111002DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TeoInferenceServiceDeploymentRecords111002DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	records, err := d.readDeploymentRecords(ctx, config)
	if err != nil {
		log.Printf("[CRITAL] DescribeInferenceServiceDeploymentRecords failed, zoneId[%s], serviceId[%s], reason[%v]",
			config.ZoneId.ValueString(), config.ServiceId.ValueString(), err)
		resp.Diagnostics.AddError(
			"Error querying TEO inference service deployment records",
			err.Error(),
		)
		return
	}

	// An empty RecordSet is a valid business state: the service simply has no
	// deployment history yet.
	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: teoInferenceDeploymentRecordAttrTypes}, records)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.RecordSet = list
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// flattenTeoInferenceDeploymentRecord converts one SDK record into the
// strongly typed model stored in state.
func flattenTeoInferenceDeploymentRecord(ctx context.Context, record *teov20220901.InferenceServiceDeploymentRecord, diags *diag.Diagnostics) teoInferenceDeploymentRecordModel {
	out := teoInferenceDeploymentRecordModel{
		RecordId:               types.StringNull(),
		Operation:              types.StringNull(),
		Status:                 types.StringNull(),
		Duration:               types.Int64Null(),
		InferenceServiceConfig: types.ListNull(types.ObjectType{AttrTypes: teoInferenceServiceConfigAttrTypes}),
		CreateTime:             types.StringNull(),
		ActiveStatus:           types.StringNull(),
	}
	if record == nil {
		return out
	}

	out.RecordId = types.StringPointerValue(record.RecordId)
	out.Operation = types.StringPointerValue(record.Operation)
	out.Status = types.StringPointerValue(record.Status)
	out.Duration = types.Int64PointerValue(record.Duration)
	out.CreateTime = types.StringPointerValue(record.CreateTime)
	out.ActiveStatus = types.StringPointerValue(record.ActiveStatus)
	out.InferenceServiceConfig = flattenTeoInferenceServiceConfig(ctx, record.InferenceServiceConfig, diags)
	return out
}

// flattenTeoInferenceServiceConfig converts the SDK inference service
// configuration into its state list value.
func flattenTeoInferenceServiceConfig(ctx context.Context, config *teov20220901.InferenceServiceConfig, diags *diag.Diagnostics) types.List {
	objectType := types.ObjectType{AttrTypes: teoInferenceServiceConfigAttrTypes}
	if config == nil {
		return types.ListNull(objectType)
	}

	model := teoInferenceServiceConfigModel{
		ListenPort:     types.Int64PointerValue(config.ListenPort),
		RequestPaths:   types.SetNull(types.StringType),
		Containers:     types.ListNull(types.ObjectType{AttrTypes: teoInferenceContainerConfigAttrTypes}),
		ResourceConfig: types.ListNull(types.ObjectType{AttrTypes: teoInferenceResourceConfigAttrTypes}),
		AffinityConfig: types.ListNull(types.ObjectType{AttrTypes: teoInferenceAffinityConfigAttrTypes}),
	}

	// RequestPaths is only set when the API returns it: a nil slice means "not
	// returned" and stays null, while an empty slice is a valid empty set.
	if config.RequestPaths != nil {
		paths := make([]string, 0, len(config.RequestPaths))
		for _, path := range config.RequestPaths {
			if path == nil {
				continue
			}
			paths = append(paths, *path)
		}
		pathsSet, d := types.SetValueFrom(ctx, types.StringType, paths)
		diags.Append(d...)
		if d.HasError() {
			return types.ListNull(objectType)
		}
		model.RequestPaths = pathsSet
	}

	if config.Containers != nil {
		containers := make([]teoInferenceContainerConfigModel, 0, len(config.Containers))
		for _, container := range config.Containers {
			containers = append(containers, flattenTeoInferenceContainerConfig(ctx, container, diags))
		}
		containersList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: teoInferenceContainerConfigAttrTypes}, containers)
		diags.Append(d...)
		if d.HasError() {
			return types.ListNull(objectType)
		}
		model.Containers = containersList
	}

	if config.ResourceConfig != nil {
		resourceConfig, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: teoInferenceResourceConfigAttrTypes},
			[]teoInferenceResourceConfigModel{flattenTeoInferenceResourceConfig(ctx, config.ResourceConfig, diags)})
		diags.Append(d...)
		if d.HasError() {
			return types.ListNull(objectType)
		}
		model.ResourceConfig = resourceConfig
	}

	if config.AffinityConfig != nil {
		affinityConfig, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: teoInferenceAffinityConfigAttrTypes},
			[]teoInferenceAffinityConfigModel{flattenTeoInferenceAffinityConfig(ctx, config.AffinityConfig)})
		diags.Append(d...)
		if d.HasError() {
			return types.ListNull(objectType)
		}
		model.AffinityConfig = affinityConfig
	}

	result, d := types.ListValueFrom(ctx, objectType, []teoInferenceServiceConfigModel{model})
	diags.Append(d...)
	if d.HasError() {
		return types.ListNull(objectType)
	}
	return result
}

// flattenTeoInferenceContainerConfig converts one SDK container configuration.
func flattenTeoInferenceContainerConfig(ctx context.Context, container *teov20220901.InferenceContainerConfig, diags *diag.Diagnostics) teoInferenceContainerConfigModel {
	out := teoInferenceContainerConfigModel{
		ImageType:            types.StringNull(),
		TcrRepositoryConfig:  types.ListNull(types.ObjectType{AttrTypes: teoInferenceTcrRepositoryConfigAttrTypes}),
		StartupCommand:       types.StringNull(),
		EnvironmentVariables: types.ListNull(types.ObjectType{AttrTypes: teoInferenceEnvironmentVariableAttrTypes}),
	}
	if container == nil {
		return out
	}

	out.ImageType = types.StringPointerValue(container.ImageType)
	out.StartupCommand = types.StringPointerValue(container.StartupCommand)

	if container.TcrRepositoryConfig != nil {
		tcr := teoInferenceTcrRepositoryConfigModel{
			TcrType:    types.StringPointerValue(container.TcrRepositoryConfig.TCRType),
			Image:      types.StringPointerValue(container.TcrRepositoryConfig.Image),
			RegistryId: types.StringPointerValue(container.TcrRepositoryConfig.RegistryId),
			RegionName: types.StringPointerValue(container.TcrRepositoryConfig.RegionName),
		}
		tcrList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: teoInferenceTcrRepositoryConfigAttrTypes}, []teoInferenceTcrRepositoryConfigModel{tcr})
		diags.Append(d...)
		if d.HasError() {
			return out
		}
		out.TcrRepositoryConfig = tcrList
	}

	if container.EnvironmentVariables != nil {
		variables := make([]teoInferenceEnvironmentVariableModel, 0, len(container.EnvironmentVariables))
		for _, variable := range container.EnvironmentVariables {
			if variable == nil {
				continue
			}
			variables = append(variables, teoInferenceEnvironmentVariableModel{
				Key:   types.StringPointerValue(variable.Key),
				Value: types.StringPointerValue(variable.Value),
			})
		}
		variablesList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: teoInferenceEnvironmentVariableAttrTypes}, variables)
		diags.Append(d...)
		if d.HasError() {
			return out
		}
		out.EnvironmentVariables = variablesList
	}

	return out
}

// flattenTeoInferenceResourceConfig converts the SDK resource configuration.
func flattenTeoInferenceResourceConfig(ctx context.Context, config *teov20220901.InferenceResourceConfig, diags *diag.Diagnostics) teoInferenceResourceConfigModel {
	out := teoInferenceResourceConfigModel{
		ScalingMode:          types.StringNull(),
		HardwareSpec:         types.StringNull(),
		HardwareSpecId:       types.StringNull(),
		HardwareConfig:       types.ListNull(types.ObjectType{AttrTypes: teoInferenceHardwareConfigAttrTypes}),
		AutoScalingConfig:    types.ListNull(types.ObjectType{AttrTypes: teoInferenceAutoScalingAttrTypes}),
		ManualInstanceConfig: types.ListNull(types.ObjectType{AttrTypes: teoInferenceManualInstanceAttrTypes}),
		Concurrency:          types.Int64Null(),
	}
	if config == nil {
		return out
	}

	out.ScalingMode = types.StringPointerValue(config.ScalingMode)
	out.HardwareSpec = types.StringPointerValue(config.HardwareSpec)
	out.HardwareSpecId = types.StringPointerValue(config.HardwareSpecId)
	out.Concurrency = types.Int64PointerValue(config.Concurrency)

	if config.HardwareConfig != nil {
		hardwareConfig := teoInferenceHardwareConfigModel{
			GpuNum:   types.Float64PointerValue(config.HardwareConfig.GPUNum),
			CpuNum:   types.Float64PointerValue(config.HardwareConfig.CPUNum),
			MemSize:  types.Int64PointerValue(config.HardwareConfig.MemSize),
			DiskSize: types.Int64PointerValue(config.HardwareConfig.DiskSize),
		}
		hardwareConfigList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: teoInferenceHardwareConfigAttrTypes}, []teoInferenceHardwareConfigModel{hardwareConfig})
		diags.Append(d...)
		if d.HasError() {
			return out
		}
		out.HardwareConfig = hardwareConfigList
	}

	if config.AutoScalingConfig != nil {
		autoScalingConfig := teoInferenceAutoScalingConfigModel{
			MinInstanceCount: types.Int64PointerValue(config.AutoScalingConfig.MinInstanceCount),
			ScalingPolicies:  types.ListNull(types.ObjectType{AttrTypes: teoInferenceScalingPolicyAttrTypes}),
		}
		if config.AutoScalingConfig.ScalingPolicies != nil {
			policies := make([]teoInferenceScalingPolicyModel, 0, len(config.AutoScalingConfig.ScalingPolicies))
			for _, policy := range config.AutoScalingConfig.ScalingPolicies {
				policies = append(policies, flattenTeoInferenceScalingPolicy(ctx, policy, diags))
			}
			policiesList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: teoInferenceScalingPolicyAttrTypes}, policies)
			diags.Append(d...)
			if d.HasError() {
				return out
			}
			autoScalingConfig.ScalingPolicies = policiesList
		}
		autoScalingConfigList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: teoInferenceAutoScalingAttrTypes}, []teoInferenceAutoScalingConfigModel{autoScalingConfig})
		diags.Append(d...)
		if d.HasError() {
			return out
		}
		out.AutoScalingConfig = autoScalingConfigList
	}

	if config.ManualInstanceConfig != nil {
		manualInstanceConfig := teoInferenceManualInstanceConfigModel{
			FixedInstanceCount: types.Int64PointerValue(config.ManualInstanceConfig.FixedInstanceCount),
		}
		manualInstanceConfigList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: teoInferenceManualInstanceAttrTypes}, []teoInferenceManualInstanceConfigModel{manualInstanceConfig})
		diags.Append(d...)
		if d.HasError() {
			return out
		}
		out.ManualInstanceConfig = manualInstanceConfigList
	}

	return out
}

// flattenTeoInferenceScalingPolicy converts one SDK scaling policy.
func flattenTeoInferenceScalingPolicy(ctx context.Context, policy *teov20220901.InferenceScalingPolicy, diags *diag.Diagnostics) teoInferenceScalingPolicyModel {
	out := teoInferenceScalingPolicyModel{
		PolicyName:             types.StringNull(),
		PolicyType:             types.StringNull(),
		ScheduledScalingPolicy: types.ListNull(types.ObjectType{AttrTypes: teoInferenceScheduledPolicyAttrTypes}),
	}
	if policy == nil {
		return out
	}

	out.PolicyName = types.StringPointerValue(policy.PolicyName)
	out.PolicyType = types.StringPointerValue(policy.PolicyType)

	if policy.ScheduledScalingPolicy != nil {
		scheduled := teoInferenceScheduledScalingPolicyModel{
			ScheduledActions: types.ListNull(types.ObjectType{AttrTypes: teoInferenceScheduledActionAttrTypes}),
			EffectiveRange:   types.ListNull(types.ObjectType{AttrTypes: teoInferenceScheduledEffectiveRangeAttrTypes}),
			TimeZone:         types.StringPointerValue(policy.ScheduledScalingPolicy.TimeZone),
		}

		if policy.ScheduledScalingPolicy.ScheduledActions != nil {
			actions := make([]teoInferenceScheduledScalingActionModel, 0, len(policy.ScheduledScalingPolicy.ScheduledActions))
			for _, action := range policy.ScheduledScalingPolicy.ScheduledActions {
				if action == nil {
					continue
				}
				actions = append(actions, teoInferenceScheduledScalingActionModel{
					CronExpression:   types.StringPointerValue(action.CronExpression),
					MinInstanceCount: types.Int64PointerValue(action.MinInstanceCount),
				})
			}
			actionsList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: teoInferenceScheduledActionAttrTypes}, actions)
			diags.Append(d...)
			if d.HasError() {
				return out
			}
			scheduled.ScheduledActions = actionsList
		}

		if policy.ScheduledScalingPolicy.EffectiveRange != nil {
			effectiveRange := teoInferenceScheduledScalingEffectiveRangeModel{
				EffectiveType: types.StringPointerValue(policy.ScheduledScalingPolicy.EffectiveRange.EffectiveType),
				StartDate:     types.StringPointerValue(policy.ScheduledScalingPolicy.EffectiveRange.StartDate),
				EndDate:       types.StringPointerValue(policy.ScheduledScalingPolicy.EffectiveRange.EndDate),
			}
			effectiveRangeList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: teoInferenceScheduledEffectiveRangeAttrTypes}, []teoInferenceScheduledScalingEffectiveRangeModel{effectiveRange})
			diags.Append(d...)
			if d.HasError() {
				return out
			}
			scheduled.EffectiveRange = effectiveRangeList
		}

		scheduledList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: teoInferenceScheduledPolicyAttrTypes}, []teoInferenceScheduledScalingPolicyModel{scheduled})
		diags.Append(d...)
		if d.HasError() {
			return out
		}
		out.ScheduledScalingPolicy = scheduledList
	}

	return out
}

// flattenTeoInferenceAffinityConfig converts the SDK affinity configuration.
func flattenTeoInferenceAffinityConfig(ctx context.Context, config *teov20220901.InferenceAffinityConfig) teoInferenceAffinityConfigModel {
	out := teoInferenceAffinityConfigModel{
		Switch:                  types.StringNull(),
		AffinityMode:            types.StringNull(),
		SessionIdAffinityConfig: types.ListNull(types.ObjectType{AttrTypes: teoInferenceSessionIdAffinityAttrTypes}),
	}
	if config == nil {
		return out
	}

	out.Switch = types.StringPointerValue(config.Switch)
	out.AffinityMode = types.StringPointerValue(config.AffinityMode)

	if config.SessionIdAffinityConfig != nil {
		sessionIdAffinityConfig := teoInferenceSessionIdAffinityConfigModel{
			Source:     types.StringPointerValue(config.SessionIdAffinityConfig.Source),
			HeaderName: types.StringPointerValue(config.SessionIdAffinityConfig.HeaderName),
		}
		sessionIdAffinityConfigList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: teoInferenceSessionIdAffinityAttrTypes}, []teoInferenceSessionIdAffinityConfigModel{sessionIdAffinityConfig})
		if d.HasError() {
			// The value is a plain object built from SDK scalars, so a failure
			// here can only mean an internal type mismatch; keep the remaining
			// fields usable instead of aborting the whole record.
			return out
		}
		out.SessionIdAffinityConfig = sessionIdAffinityConfigList
	}

	return out
}
