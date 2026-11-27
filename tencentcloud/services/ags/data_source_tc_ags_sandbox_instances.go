package ags

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	fwHelper "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/helper"
)

var _ datasource.DataSourceWithConfigure = &AgsSandboxInstancesDataSource{}

// NewAgsSandboxInstancesDataSource is the factory referenced by
// tencentcloud/framework/registry.go to register this data source.
func NewAgsSandboxInstancesDataSource() datasource.DataSource {
	return &AgsSandboxInstancesDataSource{}
}

// AgsSandboxInstancesDataSource implements datasource.DataSource for
// tencentcloud_ags_sandbox_instances.
type AgsSandboxInstancesDataSource struct {
	fw.DataSourceWithConfigure
}

// AgsSandboxInstancesDataSourceModel maps the schema attributes.
type AgsSandboxInstancesDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	InstanceIds types.List   `tfsdk:"instance_ids"`
	InstanceSet types.List   `tfsdk:"instance_set"`
}

// agsSandboxInstancesMountOptionModel is the element type of the
// custom_configuration mount_options list.
type agsSandboxInstancesMountOptionModel struct {
	Name      types.String `tfsdk:"name"`
	MountPath types.String `tfsdk:"mount_path"`
	SubPath   types.String `tfsdk:"sub_path"`
	ReadOnly  types.Bool   `tfsdk:"read_only"`
}

// agsSandboxInstancesEnvVarModel is the element type of the
// custom_configuration env list.
type agsSandboxInstancesEnvVarModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

// agsSandboxInstancesPortModel is the element type of the
// custom_configuration ports list.
type agsSandboxInstancesPortModel struct {
	Name     types.String `tfsdk:"name"`
	Port     types.Int64  `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
}

// agsSandboxInstancesResourcesModel maps the custom_configuration resources
// object.
type agsSandboxInstancesResourcesModel struct {
	Cpu     types.String `tfsdk:"cpu"`
	Memory  types.String `tfsdk:"memory"`
	Storage types.String `tfsdk:"storage"`
}

// agsSandboxInstancesHttpGetModel is the element type of the
// custom_configuration probe http_get object.
type agsSandboxInstancesHttpGetModel struct {
	Path   types.String `tfsdk:"path"`
	Port   types.Int64  `tfsdk:"port"`
	Scheme types.String `tfsdk:"scheme"`
}

// agsSandboxInstancesProbeModel maps the custom_configuration probe object.
type agsSandboxInstancesProbeModel struct {
	HttpGet          types.Object `tfsdk:"http_get"`
	ReadyTimeoutMs   types.Int64  `tfsdk:"ready_timeout_ms"`
	ProbeTimeoutMs   types.Int64  `tfsdk:"probe_timeout_ms"`
	ProbePeriodMs    types.Int64  `tfsdk:"probe_period_ms"`
	SuccessThreshold types.Int64  `tfsdk:"success_threshold"`
	FailureThreshold types.Int64  `tfsdk:"failure_threshold"`
}

// agsSandboxInstancesDnsConfigModel maps the custom_configuration
// dns_config object.
type agsSandboxInstancesDnsConfigModel struct {
	Servers  types.Set `tfsdk:"servers"`
	Searches types.Set `tfsdk:"searches"`
	Options  types.Set `tfsdk:"options"`
}

// agsSandboxInstancesCapabilitiesModel is the element type of the
// custom_configuration security_context capabilities object.
type agsSandboxInstancesCapabilitiesModel struct {
	Add  types.Set `tfsdk:"add"`
	Drop types.Set `tfsdk:"drop"`
}

// agsSandboxInstancesSecurityContextModel maps the custom_configuration
// security_context object.
type agsSandboxInstancesSecurityContextModel struct {
	Privileged   types.Bool   `tfsdk:"privileged"`
	Capabilities types.Object `tfsdk:"capabilities"`
}

// agsSandboxInstancesCustomConfigurationModel maps the custom_configuration
// object.
type agsSandboxInstancesCustomConfigurationModel struct {
	Image             types.String `tfsdk:"image"`
	ImageRegistryType types.String `tfsdk:"image_registry_type"`
	ImageDigest       types.String `tfsdk:"image_digest"`
	Command           types.Set    `tfsdk:"command"`
	Args              types.Set    `tfsdk:"args"`
	Env               types.List   `tfsdk:"env"`
	Ports             types.List   `tfsdk:"ports"`
	Resources         types.Object `tfsdk:"resources"`
	Probe             types.Object `tfsdk:"probe"`
	DnsConfig         types.Object `tfsdk:"dns_config"`
	SecurityContext   types.Object `tfsdk:"security_context"`
}

// agsSandboxInstancesWaaConfigurationModel is the element type of the
// computer_configuration waa_configuration object.
type agsSandboxInstancesWaaConfigurationModel struct {
	ImageId types.String `tfsdk:"image_id"`
}

// agsSandboxInstancesOsWorldConfigurationModel is the element type of the
// computer_configuration os_world_configuration object.
type agsSandboxInstancesOsWorldConfigurationModel struct {
	Version types.String `tfsdk:"version"`
}

// agsSandboxInstancesComputerConfigurationModel maps the
// computer_configuration object.
type agsSandboxInstancesComputerConfigurationModel struct {
	WaaConfiguration     types.Object `tfsdk:"waa_configuration"`
	OsWorldConfiguration types.Object `tfsdk:"os_world_configuration"`
}

// agsSandboxInstancesMetadataVarModel is the element type of the metadata
// list.
type agsSandboxInstancesMetadataVarModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

// AgsSandboxInstancesItemModel is the element type of instance_set.
type AgsSandboxInstancesItemModel struct {
	InstanceId            types.String `tfsdk:"instance_id"`
	ToolId                types.String `tfsdk:"tool_id"`
	ToolName              types.String `tfsdk:"tool_name"`
	Status                types.String `tfsdk:"status"`
	Persistent            types.Bool   `tfsdk:"persistent"`
	TimeoutSeconds        types.Int64  `tfsdk:"timeout_seconds"`
	ExpiresAt             types.String `tfsdk:"expires_at"`
	StopReason            types.String `tfsdk:"stop_reason"`
	CreateTime            types.String `tfsdk:"create_time"`
	UpdateTime            types.String `tfsdk:"update_time"`
	MountOptions          types.List   `tfsdk:"mount_options"`
	CustomConfiguration   types.Object `tfsdk:"custom_configuration"`
	ComputerConfiguration types.Object `tfsdk:"computer_configuration"`
	NetworkMode           types.String `tfsdk:"network_mode"`
	Metadata              types.List   `tfsdk:"metadata"`
	AuthMode              types.String `tfsdk:"auth_mode"`
}

// AgsSandboxInstancesMountOptionAttrTypes is the object type of the
// custom_configuration mount_options elements.
var AgsSandboxInstancesMountOptionAttrTypes = map[string]attr.Type{
	"name":       types.StringType,
	"mount_path": types.StringType,
	"sub_path":   types.StringType,
	"read_only":  types.BoolType,
}

// AgsSandboxInstancesEnvVarAttrTypes is the object type of the
// custom_configuration env elements.
var AgsSandboxInstancesEnvVarAttrTypes = map[string]attr.Type{
	"name":  types.StringType,
	"value": types.StringType,
}

// AgsSandboxInstancesPortAttrTypes is the object type of the
// custom_configuration ports elements.
var AgsSandboxInstancesPortAttrTypes = map[string]attr.Type{
	"name":     types.StringType,
	"port":     types.Int64Type,
	"protocol": types.StringType,
}

// AgsSandboxInstancesResourcesAttrTypes is the object type of the
// custom_configuration resources object.
var AgsSandboxInstancesResourcesAttrTypes = map[string]attr.Type{
	"cpu":     types.StringType,
	"memory":  types.StringType,
	"storage": types.StringType,
}

// AgsSandboxInstancesHttpGetAttrTypes is the object type of the
// custom_configuration probe http_get object.
var AgsSandboxInstancesHttpGetAttrTypes = map[string]attr.Type{
	"path":   types.StringType,
	"port":   types.Int64Type,
	"scheme": types.StringType,
}

// AgsSandboxInstancesProbeAttrTypes is the object type of the
// custom_configuration probe object.
var AgsSandboxInstancesProbeAttrTypes = map[string]attr.Type{
	"http_get":          types.ObjectType{AttrTypes: AgsSandboxInstancesHttpGetAttrTypes},
	"ready_timeout_ms":  types.Int64Type,
	"probe_timeout_ms":  types.Int64Type,
	"probe_period_ms":   types.Int64Type,
	"success_threshold": types.Int64Type,
	"failure_threshold": types.Int64Type,
}

// AgsSandboxInstancesDnsConfigAttrTypes is the object type of the
// custom_configuration dns_config object.
var AgsSandboxInstancesDnsConfigAttrTypes = map[string]attr.Type{
	"servers":  types.SetType{ElemType: types.StringType},
	"searches": types.SetType{ElemType: types.StringType},
	"options":  types.SetType{ElemType: types.StringType},
}

// AgsSandboxInstancesCapabilitiesAttrTypes is the object type of the
// custom_configuration security_context capabilities object.
var AgsSandboxInstancesCapabilitiesAttrTypes = map[string]attr.Type{
	"add":  types.SetType{ElemType: types.StringType},
	"drop": types.SetType{ElemType: types.StringType},
}

// AgsSandboxInstancesSecurityContextAttrTypes is the object type of the
// custom_configuration security_context object.
var AgsSandboxInstancesSecurityContextAttrTypes = map[string]attr.Type{
	"privileged":   types.BoolType,
	"capabilities": types.ObjectType{AttrTypes: AgsSandboxInstancesCapabilitiesAttrTypes},
}

// AgsSandboxInstancesCustomConfigurationAttrTypes is the object type of the
// custom_configuration object.
var AgsSandboxInstancesCustomConfigurationAttrTypes = map[string]attr.Type{
	"image":               types.StringType,
	"image_registry_type": types.StringType,
	"image_digest":        types.StringType,
	"command":             types.SetType{ElemType: types.StringType},
	"args":                types.SetType{ElemType: types.StringType},
	"env":                 types.ListType{ElemType: types.ObjectType{AttrTypes: AgsSandboxInstancesEnvVarAttrTypes}},
	"ports":               types.ListType{ElemType: types.ObjectType{AttrTypes: AgsSandboxInstancesPortAttrTypes}},
	"resources":           types.ObjectType{AttrTypes: AgsSandboxInstancesResourcesAttrTypes},
	"probe":               types.ObjectType{AttrTypes: AgsSandboxInstancesProbeAttrTypes},
	"dns_config":          types.ObjectType{AttrTypes: AgsSandboxInstancesDnsConfigAttrTypes},
	"security_context":    types.ObjectType{AttrTypes: AgsSandboxInstancesSecurityContextAttrTypes},
}

// AgsSandboxInstancesWaaConfigurationAttrTypes is the object type of the
// computer_configuration waa_configuration object.
var AgsSandboxInstancesWaaConfigurationAttrTypes = map[string]attr.Type{
	"image_id": types.StringType,
}

// AgsSandboxInstancesOsWorldConfigurationAttrTypes is the object type of the
// computer_configuration os_world_configuration object.
var AgsSandboxInstancesOsWorldConfigurationAttrTypes = map[string]attr.Type{
	"version": types.StringType,
}

// AgsSandboxInstancesComputerConfigurationAttrTypes is the object type of the
// computer_configuration object.
var AgsSandboxInstancesComputerConfigurationAttrTypes = map[string]attr.Type{
	"waa_configuration":      types.ObjectType{AttrTypes: AgsSandboxInstancesWaaConfigurationAttrTypes},
	"os_world_configuration": types.ObjectType{AttrTypes: AgsSandboxInstancesOsWorldConfigurationAttrTypes},
}

// AgsSandboxInstancesMetadataVarAttrTypes is the object type of the metadata
// list elements.
var AgsSandboxInstancesMetadataVarAttrTypes = map[string]attr.Type{
	"name":  types.StringType,
	"value": types.StringType,
}

// AgsSandboxInstancesItemAttrTypes is the object type of the instance_set
// elements.
var AgsSandboxInstancesItemAttrTypes = map[string]attr.Type{
	"instance_id":            types.StringType,
	"tool_id":                types.StringType,
	"tool_name":              types.StringType,
	"status":                 types.StringType,
	"persistent":             types.BoolType,
	"timeout_seconds":        types.Int64Type,
	"expires_at":             types.StringType,
	"stop_reason":            types.StringType,
	"create_time":            types.StringType,
	"update_time":            types.StringType,
	"mount_options":          types.ListType{ElemType: types.ObjectType{AttrTypes: AgsSandboxInstancesMountOptionAttrTypes}},
	"custom_configuration":   types.ObjectType{AttrTypes: AgsSandboxInstancesCustomConfigurationAttrTypes},
	"computer_configuration": types.ObjectType{AttrTypes: AgsSandboxInstancesComputerConfigurationAttrTypes},
	"network_mode":           types.StringType,
	"metadata":               types.ListType{ElemType: types.ObjectType{AttrTypes: AgsSandboxInstancesMetadataVarAttrTypes}},
	"auth_mode":              types.StringType,
}

func (d *AgsSandboxInstancesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "tencentcloud_ags_sandbox_instances"
}

func (d *AgsSandboxInstancesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to query the list of AGS sandbox instances.",
		Attributes: map[string]schema.Attribute{
			// The id carries no cloud meaning — it is a random token
			// generated at read time (see the end of Read).
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the data source, a random token.",
			},
			// ---- Query input ----
			"instance_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of sandbox instance IDs to query. If empty, all sandbox instances are queried. Up to 100 IDs are supported.",
			},
			// ---- Query output: Computed ----
			"instance_set": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Sandbox instance list.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier of the sandbox instance.",
						},
						"tool_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the sandbox tool the instance belongs to.",
						},
						"tool_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the sandbox tool the instance belongs to.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Instance status: STARTING (starting), RUNNING (running), STOPPING (stopping), STOPPED (stopped), STOP_FAILED (stop failed), FAILED (failed), RESIZING (specification updating).",
						},
						"persistent": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the instance is persistent.",
						},
						"timeout_seconds": schema.Int64Attribute{
							Computed:    true,
							Description: "Timeout in seconds; null means no timeout is set.",
						},
						"expires_at": schema.StringAttribute{
							Computed:    true,
							Description: "Expiration time (ISO 8601 format); null means no expiration time.",
						},
						"stop_reason": schema.StringAttribute{
							Computed:    true,
							Description: "Stop reason: manual (manual), timeout (timeout), error (error), system (system). Only present when the status is STOPPED, STOP_FAILED or FAILED. When the provider fails to stop the instance, the status is STOP_FAILED and the reason is error.",
						},
						"create_time": schema.StringAttribute{
							Computed:    true,
							Description: "Creation time (ISO 8601 format).",
						},
						"update_time": schema.StringAttribute{
							Computed:    true,
							Description: "Update time (ISO 8601 format).",
						},
						"mount_options": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Storage mount options.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Computed:    true,
										Description: "Name of the storage configuration in the sandbox tool.",
									},
									"mount_path": schema.StringAttribute{
										Computed:    true,
										Description: "Local mount path of the sandbox instance (optional), inherits the storage configuration of the tool by default.",
									},
									"sub_path": schema.StringAttribute{
										Computed:    true,
										Description: "Storage mount sub-path of the sandbox instance (optional).",
									},
									"read_only": schema.BoolAttribute{
										Computed:    true,
										Description: "Read/write permission of the sandbox instance storage mount (optional), inherits the storage configuration of the tool by default.",
									},
								},
							},
						},
						"custom_configuration": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Custom configuration of the sandbox instance.",
							Attributes: map[string]schema.Attribute{
								"image": schema.StringAttribute{
									Computed:    true,
									Description: "Image address.",
								},
								"image_registry_type": schema.StringAttribute{
									Computed:    true,
									Description: "Image registry type: `enterprise`, `personal`, `custom`. Enumerated values: enterprise: TCR enterprise container image service; personal: CCR personal container image service.",
								},
								"image_digest": schema.StringAttribute{
									Computed:    true,
									Description: "Image digest.",
								},
								"command": schema.SetAttribute{
									Computed:    true,
									ElementType: types.StringType,
									Description: "Start command.",
								},
								"args": schema.SetAttribute{
									Computed:    true,
									ElementType: types.StringType,
									Description: "Start arguments.",
								},
								"env": schema.ListNestedAttribute{
									Computed:    true,
									Description: "Environment variables.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Computed:    true,
												Description: "Environment variable name.",
											},
											"value": schema.StringAttribute{
												Computed:    true,
												Description: "Environment variable value.",
											},
										},
									},
								},
								"ports": schema.ListNestedAttribute{
									Computed:    true,
									Description: "Port configuration.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Computed:    true,
												Description: "Port name.",
											},
											"port": schema.Int64Attribute{
												Computed:    true,
												Description: "Port.",
											},
											"protocol": schema.StringAttribute{
												Computed:    true,
												Description: "Protocol.",
											},
										},
									},
								},
								"resources": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "Resource configuration.",
									Attributes: map[string]schema.Attribute{
										"cpu": schema.StringAttribute{
											Computed:    true,
											Description: "CPU resource amount.",
										},
										"memory": schema.StringAttribute{
											Computed:    true,
											Description: "Memory resource amount.",
										},
										"storage": schema.StringAttribute{
											Computed:    true,
											Description: "Custom disk size. Enumerated values: 1Gi: 1Gi; 5Gi: 5Gi; 10Gi: 10Gi; 20Gi: 20Gi. Note: this field may return null, indicating that no valid value is obtained.",
										},
									},
								},
								"probe": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "Probe configuration.",
									Attributes: map[string]schema.Attribute{
										"http_get": schema.SingleNestedAttribute{
											Computed:    true,
											Description: "HTTP GET probe configuration.",
											Attributes: map[string]schema.Attribute{
												"path": schema.StringAttribute{
													Computed:    true,
													Description: "Path.",
												},
												"port": schema.Int64Attribute{
													Computed:    true,
													Description: "Port.",
												},
												"scheme": schema.StringAttribute{
													Computed:    true,
													Description: "Protocol.",
												},
											},
										},
										"ready_timeout_ms": schema.Int64Attribute{
											Computed:    true,
											Description: "Health check ready timeout.",
										},
										"probe_timeout_ms": schema.Int64Attribute{
											Computed:    true,
											Description: "Health check single probe timeout.",
										},
										"probe_period_ms": schema.Int64Attribute{
											Computed:    true,
											Description: "Health check interval.",
										},
										"success_threshold": schema.Int64Attribute{
											Computed:    true,
											Description: "Health check success threshold.",
										},
										"failure_threshold": schema.Int64Attribute{
											Computed:    true,
											Description: "Health check failure threshold.",
										},
									},
								},
								"dns_config": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "Sandbox DNS configuration.",
									Attributes: map[string]schema.Attribute{
										"servers": schema.SetAttribute{
											Computed:    true,
											ElementType: types.StringType,
											Description: "DNS server addresses. Format: valid IP addresses. Default value: 10.0.0.1.",
										},
										"searches": schema.SetAttribute{
											Computed:    true,
											ElementType: types.StringType,
											Description: "Search domains (corresponding to the search directive of resolv.conf).",
										},
										"options": schema.SetAttribute{
											Computed:    true,
											ElementType: types.StringType,
											Description: "Configuration items (corresponding to resolv.conf options).",
										},
									},
								},
								"security_context": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "Sandbox security context.",
									Attributes: map[string]schema.Attribute{
										"privileged": schema.BoolAttribute{
											Computed:    true,
											Description: "Whether to enable privilege. Default value: true.",
										},
										"capabilities": schema.SingleNestedAttribute{
											Computed:    true,
											Description: "Capabilities refer to POSIX capabilities that can be added to or removed from a running container.",
											Attributes: map[string]schema.Attribute{
												"add": schema.SetAttribute{
													Computed:    true,
													ElementType: types.StringType,
													Description: "Capabilities to add. Value reference: [\"SYS_ADMIN\"].",
												},
												"drop": schema.SetAttribute{
													Computed:    true,
													ElementType: types.StringType,
													Description: "Capabilities to remove. Value reference: [\"ALL\"].",
												},
											},
										},
									},
								},
							},
						},
						"computer_configuration": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Configuration of desktop-computer sandbox instances.",
							Attributes: map[string]schema.Attribute{
								"waa_configuration": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "WAA sandbox tool configuration.",
									Attributes: map[string]schema.Attribute{
										"image_id": schema.StringAttribute{
											Computed:    true,
											Description: "Custom WAA image ID.",
										},
									},
								},
								"os_world_configuration": schema.SingleNestedAttribute{
									Computed:    true,
									Description: "Built-in OSWorld configuration.",
									Attributes: map[string]schema.Attribute{
										"version": schema.StringAttribute{
											Computed:    true,
											Description: "Built-in OSWorld version. Enumerated values: osworld1: OSWorld v1; osworld2: OSWorld v2. Default value: osworld1.",
										},
									},
								},
							},
						},
						"network_mode": schema.StringAttribute{
							Computed:    true,
							Description: "Network mode. Enumerated values: PUBLIC: public network access; SANDBOX: no network; INTERNAL_SERVICE: Tencent Cloud internal public services. It can override the network configuration of the tool. But if a tool itself does not support VPC networks, selecting the VPC mode in the instance settings is ineffective.",
						},
						"metadata": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Metadata of the sandbox instance.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Computed:    true,
										Description: "Metadata name.",
									},
									"value": schema.StringAttribute{
										Computed:    true,
										Description: "Metadata value.",
									},
								},
							},
						},
						"auth_mode": schema.StringAttribute{
							Computed:    true,
							Description: "Sandbox access authentication mode. Enumerated values: DEFAULT: default, i.e. TOKEN authentication; TOKEN: Token authentication, i.e. all port access requires a TOKEN; NONE: no authentication, i.e. all port access does not require a TOKEN; PUBLIC: public mode, i.e. the ENVD management port (49983) access requires a TOKEN while other ports do not. Default value: DEFAULT.",
						},
					},
				},
			},
		},
	}
}

// Read queries the cloud API and writes the result into resp.State.
func (d *AgsSandboxInstancesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config AgsSandboxInstancesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	describeSandboxInstancesRequest := agsv20250920.NewDescribeSandboxInstancesRequest()
	// An optional list filter that the SDK takes as []*string —
	// fwHelper.PointersFromList[string] (null/unknown -> nil, i.e. no filter).
	describeSandboxInstancesRequest.InstanceIds = fwHelper.PointersFromList[string](ctx, config.InstanceIds, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewAgsService(d.Client())

	instanceSet, err := service.DescribeSandboxInstances(ctx, describeSandboxInstancesRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying AGS sandbox instances",
			fmt.Sprintf("DescribeSandboxInstances failed: %v", err),
		)

		return
	}

	// An empty result is a genuinely empty list, not a failure.
	items := make([]AgsSandboxInstancesItemModel, 0, len(instanceSet))
	for _, item := range instanceSet {
		if item == nil {
			continue
		}

		itemModel := AgsSandboxInstancesItemModel{
			InstanceId:     types.StringPointerValue(item.InstanceId),
			ToolId:         types.StringPointerValue(item.ToolId),
			ToolName:       types.StringPointerValue(item.ToolName),
			Status:         types.StringPointerValue(item.Status),
			Persistent:     types.BoolPointerValue(item.Persistent),
			TimeoutSeconds: fwHelper.Int64ValueFromUint(item.TimeoutSeconds),
			ExpiresAt:      types.StringPointerValue(item.ExpiresAt),
			StopReason:     types.StringPointerValue(item.StopReason),
			CreateTime:     types.StringPointerValue(item.CreateTime),
			UpdateTime:     types.StringPointerValue(item.UpdateTime),
			NetworkMode:    types.StringPointerValue(item.NetworkMode),
			AuthMode:       types.StringPointerValue(item.AuthMode),
		}

		itemModel.MountOptions = flattenAgsSandboxInstancesMountOptions(ctx, item.MountOptions, &resp.Diagnostics)
		itemModel.CustomConfiguration = flattenAgsSandboxInstancesCustomConfiguration(ctx, item.CustomConfiguration, &resp.Diagnostics)
		itemModel.ComputerConfiguration = flattenAgsSandboxInstancesComputerConfiguration(ctx, item.ComputerConfiguration, &resp.Diagnostics)
		itemModel.Metadata = flattenAgsSandboxInstancesMetadata(ctx, item.Metadata, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		items = append(items, itemModel)
	}

	instanceSetList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: AgsSandboxInstancesItemAttrTypes}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.InstanceSet = instanceSetList

	// The state id is a random token (JSON resource_id = fwHelper.BuildToken()).
	config.Id = types.StringValue(fwHelper.BuildToken())

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// flattenAgsSandboxInstancesMountOptions maps the cloud mount options list
// into the instance_set mount_options list.
func flattenAgsSandboxInstancesMountOptions(ctx context.Context, cloud []*agsv20250920.MountOption, diags *diag.Diagnostics) types.List {
	if cloud == nil {
		return types.ListNull(types.ObjectType{AttrTypes: AgsSandboxInstancesMountOptionAttrTypes})
	}

	items := make([]agsSandboxInstancesMountOptionModel, 0, len(cloud))
	for _, item := range cloud {
		if item == nil {
			continue
		}

		items = append(items, agsSandboxInstancesMountOptionModel{
			Name:      types.StringPointerValue(item.Name),
			MountPath: types.StringPointerValue(item.MountPath),
			SubPath:   types.StringPointerValue(item.SubPath),
			ReadOnly:  types.BoolPointerValue(item.ReadOnly),
		})
	}

	list, listDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: AgsSandboxInstancesMountOptionAttrTypes}, items)
	diags.Append(listDiags...)

	return list
}

// flattenAgsSandboxInstancesCustomConfiguration maps the cloud custom
// configuration into the instance_set custom_configuration object.
func flattenAgsSandboxInstancesCustomConfiguration(ctx context.Context, cloud *agsv20250920.CustomConfigurationDetail, diags *diag.Diagnostics) types.Object {
	if cloud == nil {
		return types.ObjectNull(AgsSandboxInstancesCustomConfigurationAttrTypes)
	}

	var model agsSandboxInstancesCustomConfigurationModel
	model.Image = types.StringPointerValue(cloud.Image)
	model.ImageRegistryType = types.StringPointerValue(cloud.ImageRegistryType)
	model.ImageDigest = types.StringPointerValue(cloud.ImageDigest)
	model.Command = fwHelper.SetFromPointers[string](ctx, types.StringType, cloud.Command, diags)
	model.Args = fwHelper.SetFromPointers[string](ctx, types.StringType, cloud.Args, diags)
	model.Resources = flattenAgsSandboxInstancesResources(ctx, cloud.Resources, diags)
	model.Probe = flattenAgsSandboxInstancesProbe(ctx, cloud.Probe, diags)
	model.DnsConfig = flattenAgsSandboxInstancesDnsConfig(ctx, cloud.DNSConfig, diags)
	model.SecurityContext = flattenAgsSandboxInstancesSecurityContext(ctx, cloud.SecurityContext, diags)

	if cloud.Env == nil {
		model.Env = types.ListNull(types.ObjectType{AttrTypes: AgsSandboxInstancesEnvVarAttrTypes})
	} else {
		envItems := make([]agsSandboxInstancesEnvVarModel, 0, len(cloud.Env))
		for _, env := range cloud.Env {
			if env == nil {
				continue
			}

			envItems = append(envItems, agsSandboxInstancesEnvVarModel{
				Name:  types.StringPointerValue(env.Name),
				Value: types.StringPointerValue(env.Value),
			})
		}

		envList, envDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: AgsSandboxInstancesEnvVarAttrTypes}, envItems)
		diags.Append(envDiags...)
		model.Env = envList
	}

	if cloud.Ports == nil {
		model.Ports = types.ListNull(types.ObjectType{AttrTypes: AgsSandboxInstancesPortAttrTypes})
	} else {
		portItems := make([]agsSandboxInstancesPortModel, 0, len(cloud.Ports))
		for _, port := range cloud.Ports {
			if port == nil {
				continue
			}

			portItems = append(portItems, agsSandboxInstancesPortModel{
				Name:     types.StringPointerValue(port.Name),
				Port:     types.Int64PointerValue(port.Port),
				Protocol: types.StringPointerValue(port.Protocol),
			})
		}

		portList, portDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: AgsSandboxInstancesPortAttrTypes}, portItems)
		diags.Append(portDiags...)
		model.Ports = portList
	}

	obj, objDiags := types.ObjectValueFrom(ctx, AgsSandboxInstancesCustomConfigurationAttrTypes, model)
	diags.Append(objDiags...)

	return obj
}

// flattenAgsSandboxInstancesResources maps the cloud resource configuration
// into the custom_configuration resources object.
func flattenAgsSandboxInstancesResources(ctx context.Context, cloud *agsv20250920.ResourceConfiguration, diags *diag.Diagnostics) types.Object {
	if cloud == nil {
		return types.ObjectNull(AgsSandboxInstancesResourcesAttrTypes)
	}

	model := agsSandboxInstancesResourcesModel{
		Cpu:     types.StringPointerValue(cloud.CPU),
		Memory:  types.StringPointerValue(cloud.Memory),
		Storage: types.StringPointerValue(cloud.Storage),
	}

	obj, objDiags := types.ObjectValueFrom(ctx, AgsSandboxInstancesResourcesAttrTypes, model)
	diags.Append(objDiags...)

	return obj
}

// flattenAgsSandboxInstancesProbe maps the cloud probe configuration into the
// custom_configuration probe object.
func flattenAgsSandboxInstancesProbe(ctx context.Context, cloud *agsv20250920.ProbeConfiguration, diags *diag.Diagnostics) types.Object {
	if cloud == nil {
		return types.ObjectNull(AgsSandboxInstancesProbeAttrTypes)
	}

	model := agsSandboxInstancesProbeModel{
		HttpGet:          flattenAgsSandboxInstancesHttpGet(ctx, cloud.HttpGet, diags),
		ReadyTimeoutMs:   types.Int64PointerValue(cloud.ReadyTimeoutMs),
		ProbeTimeoutMs:   types.Int64PointerValue(cloud.ProbeTimeoutMs),
		ProbePeriodMs:    types.Int64PointerValue(cloud.ProbePeriodMs),
		SuccessThreshold: types.Int64PointerValue(cloud.SuccessThreshold),
		FailureThreshold: types.Int64PointerValue(cloud.FailureThreshold),
	}

	obj, objDiags := types.ObjectValueFrom(ctx, AgsSandboxInstancesProbeAttrTypes, model)
	diags.Append(objDiags...)

	return obj
}

// flattenAgsSandboxInstancesHttpGet maps the cloud HTTP GET probe
// configuration into the probe http_get object.
func flattenAgsSandboxInstancesHttpGet(ctx context.Context, cloud *agsv20250920.HttpGetAction, diags *diag.Diagnostics) types.Object {
	if cloud == nil {
		return types.ObjectNull(AgsSandboxInstancesHttpGetAttrTypes)
	}

	model := agsSandboxInstancesHttpGetModel{
		Path:   types.StringPointerValue(cloud.Path),
		Port:   types.Int64PointerValue(cloud.Port),
		Scheme: types.StringPointerValue(cloud.Scheme),
	}

	obj, objDiags := types.ObjectValueFrom(ctx, AgsSandboxInstancesHttpGetAttrTypes, model)
	diags.Append(objDiags...)

	return obj
}

// flattenAgsSandboxInstancesDnsConfig maps the cloud DNS configuration into
// the custom_configuration dns_config object.
func flattenAgsSandboxInstancesDnsConfig(ctx context.Context, cloud *agsv20250920.DNSConfig, diags *diag.Diagnostics) types.Object {
	if cloud == nil {
		return types.ObjectNull(AgsSandboxInstancesDnsConfigAttrTypes)
	}

	model := agsSandboxInstancesDnsConfigModel{
		Servers:  fwHelper.SetFromPointers[string](ctx, types.StringType, cloud.Servers, diags),
		Searches: fwHelper.SetFromPointers[string](ctx, types.StringType, cloud.Searches, diags),
		Options:  fwHelper.SetFromPointers[string](ctx, types.StringType, cloud.Options, diags),
	}

	obj, objDiags := types.ObjectValueFrom(ctx, AgsSandboxInstancesDnsConfigAttrTypes, model)
	diags.Append(objDiags...)

	return obj
}

// flattenAgsSandboxInstancesSecurityContext maps the cloud security context
// into the custom_configuration security_context object.
func flattenAgsSandboxInstancesSecurityContext(ctx context.Context, cloud *agsv20250920.SecurityContext, diags *diag.Diagnostics) types.Object {
	if cloud == nil {
		return types.ObjectNull(AgsSandboxInstancesSecurityContextAttrTypes)
	}

	var model agsSandboxInstancesSecurityContextModel
	model.Privileged = types.BoolPointerValue(cloud.Privileged)
	if cloud.Capabilities == nil {
		model.Capabilities = types.ObjectNull(AgsSandboxInstancesCapabilitiesAttrTypes)
	} else {
		capModel := agsSandboxInstancesCapabilitiesModel{
			Add:  fwHelper.SetFromPointers[string](ctx, types.StringType, cloud.Capabilities.Add, diags),
			Drop: fwHelper.SetFromPointers[string](ctx, types.StringType, cloud.Capabilities.Drop, diags),
		}

		capObj, capDiags := types.ObjectValueFrom(ctx, AgsSandboxInstancesCapabilitiesAttrTypes, capModel)
		diags.Append(capDiags...)
		model.Capabilities = capObj
	}

	obj, objDiags := types.ObjectValueFrom(ctx, AgsSandboxInstancesSecurityContextAttrTypes, model)
	diags.Append(objDiags...)

	return obj
}

// flattenAgsSandboxInstancesComputerConfiguration maps the cloud computer
// configuration into the instance_set computer_configuration object.
func flattenAgsSandboxInstancesComputerConfiguration(ctx context.Context, cloud *agsv20250920.ComputerConfiguration, diags *diag.Diagnostics) types.Object {
	if cloud == nil {
		return types.ObjectNull(AgsSandboxInstancesComputerConfigurationAttrTypes)
	}

	var model agsSandboxInstancesComputerConfigurationModel
	if cloud.WAAConfiguration == nil {
		model.WaaConfiguration = types.ObjectNull(AgsSandboxInstancesWaaConfigurationAttrTypes)
	} else {
		waaModel := agsSandboxInstancesWaaConfigurationModel{
			ImageId: types.StringPointerValue(cloud.WAAConfiguration.ImageId),
		}

		waaObj, waaDiags := types.ObjectValueFrom(ctx, AgsSandboxInstancesWaaConfigurationAttrTypes, waaModel)
		diags.Append(waaDiags...)
		model.WaaConfiguration = waaObj
	}

	if cloud.OSWorldConfiguration == nil {
		model.OsWorldConfiguration = types.ObjectNull(AgsSandboxInstancesOsWorldConfigurationAttrTypes)
	} else {
		osWorldModel := agsSandboxInstancesOsWorldConfigurationModel{
			Version: types.StringPointerValue(cloud.OSWorldConfiguration.Version),
		}

		osWorldObj, osWorldDiags := types.ObjectValueFrom(ctx, AgsSandboxInstancesOsWorldConfigurationAttrTypes, osWorldModel)
		diags.Append(osWorldDiags...)
		model.OsWorldConfiguration = osWorldObj
	}

	obj, objDiags := types.ObjectValueFrom(ctx, AgsSandboxInstancesComputerConfigurationAttrTypes, model)
	diags.Append(objDiags...)

	return obj
}

// flattenAgsSandboxInstancesMetadata maps the cloud metadata list into the
// instance_set metadata list.
func flattenAgsSandboxInstancesMetadata(ctx context.Context, cloud []*agsv20250920.MetadataVar, diags *diag.Diagnostics) types.List {
	if cloud == nil {
		return types.ListNull(types.ObjectType{AttrTypes: AgsSandboxInstancesMetadataVarAttrTypes})
	}

	items := make([]agsSandboxInstancesMetadataVarModel, 0, len(cloud))
	for _, item := range cloud {
		if item == nil {
			continue
		}

		items = append(items, agsSandboxInstancesMetadataVarModel{
			Name:  types.StringPointerValue(item.Name),
			Value: types.StringPointerValue(item.Value),
		})
	}

	list, listDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: AgsSandboxInstancesMetadataVarAttrTypes}, items)
	diags.Append(listDiags...)

	return list
}
