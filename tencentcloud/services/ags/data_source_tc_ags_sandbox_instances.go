// Package ags provides the framework reference implementation of the AGS
// resources and data sources. This file implements the
// tencentcloud_ags_sandbox_instances data source backed by the
// DescribeSandboxInstances API.
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

// agsSandboxInstancesMountOptionModel is the element type of mount_options.
type agsSandboxInstancesMountOptionModel struct {
	Name      types.String `tfsdk:"name"`
	MountPath types.String `tfsdk:"mount_path"`
	SubPath   types.String `tfsdk:"sub_path"`
	ReadOnly  types.Bool   `tfsdk:"read_only"`
}

var agsSandboxInstancesMountOptionAttrTypes = map[string]attr.Type{
	"name":       types.StringType,
	"mount_path": types.StringType,
	"sub_path":   types.StringType,
	"read_only":  types.BoolType,
}

// agsSandboxInstancesEnvVarModel is the element type of custom_configuration.env.
type agsSandboxInstancesEnvVarModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

var agsSandboxInstancesEnvVarAttrTypes = map[string]attr.Type{
	"name":  types.StringType,
	"value": types.StringType,
}

// agsSandboxInstancesPortModel is the element type of custom_configuration.ports.
type agsSandboxInstancesPortModel struct {
	Name     types.String `tfsdk:"name"`
	Port     types.Int64  `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
}

var agsSandboxInstancesPortAttrTypes = map[string]attr.Type{
	"name":     types.StringType,
	"port":     types.Int64Type,
	"protocol": types.StringType,
}

// agsSandboxInstancesResourcesModel maps custom_configuration.resources.
type agsSandboxInstancesResourcesModel struct {
	CPU     types.String `tfsdk:"cpu"`
	Memory  types.String `tfsdk:"memory"`
	Storage types.String `tfsdk:"storage"`
}

var agsSandboxInstancesResourcesAttrTypes = map[string]attr.Type{
	"cpu":     types.StringType,
	"memory":  types.StringType,
	"storage": types.StringType,
}

// agsSandboxInstancesHttpGetModel is the element type of custom_configuration.probe.http_get.
type agsSandboxInstancesHttpGetModel struct {
	Path   types.String `tfsdk:"path"`
	Port   types.Int64  `tfsdk:"port"`
	Scheme types.String `tfsdk:"scheme"`
}

var agsSandboxInstancesHttpGetAttrTypes = map[string]attr.Type{
	"path":   types.StringType,
	"port":   types.Int64Type,
	"scheme": types.StringType,
}

// agsSandboxInstancesProbeModel maps custom_configuration.probe.
type agsSandboxInstancesProbeModel struct {
	HttpGet          types.Object `tfsdk:"http_get"`
	ReadyTimeoutMs   types.Int64  `tfsdk:"ready_timeout_ms"`
	ProbeTimeoutMs   types.Int64  `tfsdk:"probe_timeout_ms"`
	ProbePeriodMs    types.Int64  `tfsdk:"probe_period_ms"`
	SuccessThreshold types.Int64  `tfsdk:"success_threshold"`
	FailureThreshold types.Int64  `tfsdk:"failure_threshold"`
}

var agsSandboxInstancesProbeAttrTypes = map[string]attr.Type{
	"http_get":          types.ObjectType{AttrTypes: agsSandboxInstancesHttpGetAttrTypes},
	"ready_timeout_ms":  types.Int64Type,
	"probe_timeout_ms":  types.Int64Type,
	"probe_period_ms":   types.Int64Type,
	"success_threshold": types.Int64Type,
	"failure_threshold": types.Int64Type,
}

// agsSandboxInstancesDNSConfigModel maps custom_configuration.dns_config.
type agsSandboxInstancesDNSConfigModel struct {
	Servers  types.List `tfsdk:"servers"`
	Searches types.List `tfsdk:"searches"`
	Options  types.List `tfsdk:"options"`
}

var agsSandboxInstancesDNSConfigAttrTypes = map[string]attr.Type{
	"servers":  types.ListType{ElemType: types.StringType},
	"searches": types.ListType{ElemType: types.StringType},
	"options":  types.ListType{ElemType: types.StringType},
}

// agsSandboxInstancesCapabilitiesModel is the element type of
// custom_configuration.security_context.capabilities.
type agsSandboxInstancesCapabilitiesModel struct {
	Add  types.List `tfsdk:"add"`
	Drop types.List `tfsdk:"drop"`
}

var agsSandboxInstancesCapabilitiesAttrTypes = map[string]attr.Type{
	"add":  types.ListType{ElemType: types.StringType},
	"drop": types.ListType{ElemType: types.StringType},
}

// agsSandboxInstancesSecurityContextModel maps custom_configuration.security_context.
type agsSandboxInstancesSecurityContextModel struct {
	Privileged   types.Bool   `tfsdk:"privileged"`
	Capabilities types.Object `tfsdk:"capabilities"`
}

var agsSandboxInstancesSecurityContextAttrTypes = map[string]attr.Type{
	"privileged":   types.BoolType,
	"capabilities": types.ObjectType{AttrTypes: agsSandboxInstancesCapabilitiesAttrTypes},
}

// agsSandboxInstancesCustomConfigurationModel maps custom_configuration.
type agsSandboxInstancesCustomConfigurationModel struct {
	Image             types.String `tfsdk:"image"`
	ImageRegistryType types.String `tfsdk:"image_registry_type"`
	ImageDigest       types.String `tfsdk:"image_digest"`
	Command           types.List   `tfsdk:"command"`
	Args              types.List   `tfsdk:"args"`
	Env               types.List   `tfsdk:"env"`
	Ports             types.List   `tfsdk:"ports"`
	Resources         types.Object `tfsdk:"resources"`
	Probe             types.Object `tfsdk:"probe"`
	DNSConfig         types.Object `tfsdk:"dns_config"`
	SecurityContext   types.Object `tfsdk:"security_context"`
}

var agsSandboxInstancesCustomConfigurationAttrTypes = map[string]attr.Type{
	"image":               types.StringType,
	"image_registry_type": types.StringType,
	"image_digest":        types.StringType,
	"command":             types.ListType{ElemType: types.StringType},
	"args":                types.ListType{ElemType: types.StringType},
	"env":                 types.ListType{ElemType: types.ObjectType{AttrTypes: agsSandboxInstancesEnvVarAttrTypes}},
	"ports":               types.ListType{ElemType: types.ObjectType{AttrTypes: agsSandboxInstancesPortAttrTypes}},
	"resources":           types.ObjectType{AttrTypes: agsSandboxInstancesResourcesAttrTypes},
	"probe":               types.ObjectType{AttrTypes: agsSandboxInstancesProbeAttrTypes},
	"dns_config":          types.ObjectType{AttrTypes: agsSandboxInstancesDNSConfigAttrTypes},
	"security_context":    types.ObjectType{AttrTypes: agsSandboxInstancesSecurityContextAttrTypes},
}

// agsSandboxInstancesWAAConfigurationModel maps
// computer_configuration.waa_configuration.
type agsSandboxInstancesWAAConfigurationModel struct {
	ImageId types.String `tfsdk:"image_id"`
}

var agsSandboxInstancesWAAConfigurationAttrTypes = map[string]attr.Type{
	"image_id": types.StringType,
}

// agsSandboxInstancesOSWorldConfigurationModel maps
// computer_configuration.os_world_configuration.
type agsSandboxInstancesOSWorldConfigurationModel struct {
	Version types.String `tfsdk:"version"`
}

var agsSandboxInstancesOSWorldConfigurationAttrTypes = map[string]attr.Type{
	"version": types.StringType,
}

// agsSandboxInstancesComputerConfigurationModel maps computer_configuration.
type agsSandboxInstancesComputerConfigurationModel struct {
	WAAConfiguration     types.Object `tfsdk:"waa_configuration"`
	OSWorldConfiguration types.Object `tfsdk:"os_world_configuration"`
}

var agsSandboxInstancesComputerConfigurationAttrTypes = map[string]attr.Type{
	"waa_configuration":      types.ObjectType{AttrTypes: agsSandboxInstancesWAAConfigurationAttrTypes},
	"os_world_configuration": types.ObjectType{AttrTypes: agsSandboxInstancesOSWorldConfigurationAttrTypes},
}

// agsSandboxInstancesMetadataModel is the element type of metadata.
type agsSandboxInstancesMetadataModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

var agsSandboxInstancesMetadataAttrTypes = map[string]attr.Type{
	"name":  types.StringType,
	"value": types.StringType,
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
	"mount_options":          types.ListType{ElemType: types.ObjectType{AttrTypes: agsSandboxInstancesMountOptionAttrTypes}},
	"custom_configuration":   types.ObjectType{AttrTypes: agsSandboxInstancesCustomConfigurationAttrTypes},
	"computer_configuration": types.ObjectType{AttrTypes: agsSandboxInstancesComputerConfigurationAttrTypes},
	"network_mode":           types.StringType,
	"metadata":               types.ListType{ElemType: types.ObjectType{AttrTypes: agsSandboxInstancesMetadataAttrTypes}},
	"auth_mode":              types.StringType,
}

// agsSandboxInstancesItemAttributes returns the attribute map of one
// instance_set element, shared by the schema and the nested conversions.
func agsSandboxInstancesItemAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"instance_id": schema.StringAttribute{
			Computed:    true,
			Description: "Unique identifier of the sandbox instance.",
		},
		"tool_id": schema.StringAttribute{
			Computed:    true,
			Description: "Id of the sandbox tool the instance belongs to.",
		},
		"tool_name": schema.StringAttribute{
			Computed:    true,
			Description: "Name of the sandbox tool the instance belongs to.",
		},
		"status": schema.StringAttribute{
			Computed: true,
			Description: "Status of the instance. Valid values: `STARTING`, `RUNNING`, `STOPPING`, `STOPPED`, " +
				"`STOP_FAILED`, `FAILED`, `RESIZING`.",
		},
		"persistent": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether the instance is persistent.",
		},
		"timeout_seconds": schema.Int64Attribute{
			Computed:    true,
			Description: "Timeout in seconds. `null` means no timeout is set.",
		},
		"expires_at": schema.StringAttribute{
			Computed:    true,
			Description: "Expiration time in ISO 8601 format. `null` means no expiration time.",
		},
		"stop_reason": schema.StringAttribute{
			Computed: true,
			Description: "Reason why the instance was stopped. Valid values: `manual`, `timeout`, `error`, " +
				"`system`. It is only set when the status is `STOPPED`, `STOP_FAILED` or `FAILED`. When the " +
				"provider stops the instance and the stop fails, the status is `STOP_FAILED` and the reason " +
				"is `error`.",
		},
		"create_time": schema.StringAttribute{
			Computed:    true,
			Description: "Creation time of the instance in ISO 8601 format.",
		},
		"update_time": schema.StringAttribute{
			Computed:    true,
			Description: "Last update time of the instance in ISO 8601 format.",
		},
		"mount_options": schema.ListNestedAttribute{
			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: agsSandboxInstancesMountOptionAttributes()},
			Description:  "Storage mount options.",
		},
		"custom_configuration": schema.SingleNestedAttribute{
			Computed:    true,
			Attributes:  agsSandboxInstancesCustomConfigurationAttributes(),
			Description: "Custom configuration of the sandbox instance.",
		},
		"computer_configuration": schema.SingleNestedAttribute{
			Computed:    true,
			Attributes:  agsSandboxInstancesComputerConfigurationAttributes(),
			Description: "Configuration of desktop-like sandbox instances.",
		},
		"network_mode": schema.StringAttribute{
			Computed: true,
			Description: "Network mode. Valid values: `PUBLIC` (public network access), `SANDBOX` (no network), " +
				"`INTERNAL_SERVICE` (Tencent Cloud internal public services). It can override the network " +
				"configuration of the tool level. However, if a tool itself does not support VPC networks, " +
				"selecting the VPC mode in the instance settings has no effect.",
		},
		"metadata": schema.ListNestedAttribute{
			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: agsSandboxInstancesMetadataAttributes()},
			Description:  "Metadata of the sandbox instance.",
		},
		"auth_mode": schema.StringAttribute{
			Computed: true,
			Description: "Authentication mode of sandbox access. Valid values: `DEFAULT` (default, i.e. TOKEN " +
				"authentication), `TOKEN` (all port accesses must carry the TOKEN), `NONE` (all port accesses " +
				"do not need the TOKEN), `PUBLIC` (public mode, i.e. the ENVD management port (49983) access " +
				"must carry the TOKEN while other ports do not need it). Default value: `DEFAULT`.",
		},
	}
}

func agsSandboxInstancesMountOptionAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Computed:    true,
			Description: "Name of the storage configuration in the sandbox tool.",
		},
		"mount_path": schema.StringAttribute{
			Computed:    true,
			Description: "Local mount path inside the sandbox instance (optional). By default it inherits the storage configuration of the tool.",
		},
		"sub_path": schema.StringAttribute{
			Computed:    true,
			Description: "Sub path of the sandbox instance storage mount (optional).",
		},
		"read_only": schema.BoolAttribute{
			Computed:    true,
			Description: "Read and write permission of the sandbox instance storage mount (optional). By default it inherits the storage configuration of the tool.",
		},
	}
}

func agsSandboxInstancesCustomConfigurationAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"image": schema.StringAttribute{
			Computed:    true,
			Description: "Image address.",
		},
		"image_registry_type": schema.StringAttribute{
			Computed: true,
			Description: "Image registry type. Valid values: `enterprise`, `personal`, `custom`. Enumerated " +
				"values: `enterprise`: TCR enterprise container image service; `personal`: CCR personal " +
				"container image service.",
		},
		"image_digest": schema.StringAttribute{
			Computed:    true,
			Description: "Image digest.",
		},
		"command": schema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
			Description: "Startup command.",
		},
		"args": schema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
			Description: "Startup arguments.",
		},
		"env": schema.ListNestedAttribute{
			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: agsSandboxInstancesEnvVarAttributes()},
			Description:  "Environment variables.",
		},
		"ports": schema.ListNestedAttribute{
			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: agsSandboxInstancesPortAttributes()},
			Description:  "Port configuration.",
		},
		"resources": schema.SingleNestedAttribute{
			Computed:    true,
			Attributes:  agsSandboxInstancesResourcesAttributes(),
			Description: "Resource configuration.",
		},
		"probe": schema.SingleNestedAttribute{
			Computed:    true,
			Attributes:  agsSandboxInstancesProbeAttributes(),
			Description: "Probe configuration.",
		},
		"dns_config": schema.SingleNestedAttribute{
			Computed:    true,
			Attributes:  agsSandboxInstancesDNSConfigAttributes(),
			Description: "Sandbox DNS configuration.",
		},
		"security_context": schema.SingleNestedAttribute{
			Computed:    true,
			Attributes:  agsSandboxInstancesSecurityContextAttributes(),
			Description: "Sandbox security context.",
		},
	}
}

func agsSandboxInstancesEnvVarAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Computed:    true,
			Description: "Name of the environment variable.",
		},
		"value": schema.StringAttribute{
			Computed:    true,
			Description: "Value of the environment variable.",
		},
	}
}

func agsSandboxInstancesPortAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Computed:    true,
			Description: "Name of the port.",
		},
		"port": schema.Int64Attribute{
			Computed:    true,
			Description: "Port number.",
		},
		"protocol": schema.StringAttribute{
			Computed:    true,
			Description: "Protocol.",
		},
	}
}

func agsSandboxInstancesResourcesAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"cpu": schema.StringAttribute{
			Computed:    true,
			Description: "Amount of cpu resources.",
		},
		"memory": schema.StringAttribute{
			Computed:    true,
			Description: "Amount of memory resources.",
		},
		"storage": schema.StringAttribute{
			Computed: true,
			Description: "Custom disk size. Valid values: `1Gi`, `5Gi`, `10Gi`, `20Gi`. Note: this field may " +
				"return `null`, which means no valid value is obtained.",
		},
	}
}

func agsSandboxInstancesProbeAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"http_get": schema.SingleNestedAttribute{
			Computed:    true,
			Attributes:  agsSandboxInstancesHttpGetAttributes(),
			Description: "HTTP GET probe configuration.",
		},
		"ready_timeout_ms": schema.Int64Attribute{
			Computed:    true,
			Description: "Timeout in milliseconds for the health check to become ready.",
		},
		"probe_timeout_ms": schema.Int64Attribute{
			Computed:    true,
			Description: "Timeout in milliseconds of a single health check probe.",
		},
		"probe_period_ms": schema.Int64Attribute{
			Computed:    true,
			Description: "Interval in milliseconds between two health check probes.",
		},
		"success_threshold": schema.Int64Attribute{
			Computed:    true,
			Description: "Success threshold of the health check.",
		},
		"failure_threshold": schema.Int64Attribute{
			Computed:    true,
			Description: "Failure threshold of the health check.",
		},
	}
}

func agsSandboxInstancesHttpGetAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
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
	}
}

func agsSandboxInstancesDNSConfigAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"servers": schema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
			Description: "DNS server addresses. Valid IP addresses are required. Default value: `10.0.0.1`.",
		},
		"searches": schema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
			Description: "Search domains, corresponding to the `search` directive of resolv.conf.",
		},
		"options": schema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
			Description: "Configuration items, corresponding to the options of resolv.conf.",
		},
	}
}

func agsSandboxInstancesSecurityContextAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"privileged": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether privileged mode is enabled. Default value: `true`.",
		},
		"capabilities": schema.SingleNestedAttribute{
			Computed:    true,
			Attributes:  agsSandboxInstancesCapabilitiesAttributes(),
			Description: "POSIX capabilities that can be added to or removed from the running container.",
		},
	}
}

func agsSandboxInstancesCapabilitiesAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"add": schema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
			Description: "Capabilities to add. Reference value: `[\"SYS_ADMIN\"]`.",
		},
		"drop": schema.ListAttribute{
			Computed:    true,
			ElementType: types.StringType,
			Description: "Capabilities to remove. Reference value: `[\"ALL\"]`.",
		},
	}
}

func agsSandboxInstancesComputerConfigurationAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"waa_configuration": schema.SingleNestedAttribute{
			Computed:    true,
			Attributes:  agsSandboxInstancesWAAConfigurationAttributes(),
			Description: "Configuration of the WAA sandbox tool.",
		},
		"os_world_configuration": schema.SingleNestedAttribute{
			Computed:    true,
			Attributes:  agsSandboxInstancesOSWorldConfigurationAttributes(),
			Description: "Configuration of the built-in OSWorld.",
		},
	}
}

func agsSandboxInstancesWAAConfigurationAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"image_id": schema.StringAttribute{
			Computed:    true,
			Description: "Id of the custom WAA image.",
		},
	}
}

func agsSandboxInstancesOSWorldConfigurationAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"version": schema.StringAttribute{
			Computed: true,
			Description: "Version of the built-in OSWorld. Valid values: `osworld1` (OSWorld v1), `osworld2` " +
				"(OSWorld v2). Default value: `osworld1`.",
		},
	}
}

func agsSandboxInstancesMetadataAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Computed:    true,
			Description: "Name of the metadata.",
		},
		"value": schema.StringAttribute{
			Computed:    true,
			Description: "Value of the metadata.",
		},
	}
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
				Description: "Id of the data source, a random token.",
			},
			// ---- Query input ----
			"instance_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of sandbox instance ids to query. If it is empty, all instances are queried. At most 100 ids are supported.",
			},
			// ---- Query output: Computed ----
			// Outputs carry the resource data only — never a
			// total_count/xxx_count aggregate, never paging knobs.
			"instance_set": schema.ListNestedAttribute{
				Computed:     true,
				NestedObject: schema.NestedAttributeObject{Attributes: agsSandboxInstancesItemAttributes()},
				Description:  "List of sandbox instances.",
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

	instances, err := service.DescribeSandboxInstances(ctx, describeSandboxInstancesRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error querying AGS sandbox instances",
			fmt.Sprintf("DescribeSandboxInstances failed: %v", err),
		)

		return
	}

	// An empty result is a genuinely empty list, not a failure.
	items := make([]AgsSandboxInstancesItemModel, 0, len(instances))
	for _, item := range instances {
		if item == nil {
			continue
		}

		itemModel, diags := flattenAgsSandboxInstance(ctx, item)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		items = append(items, itemModel)
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: AgsSandboxInstancesItemAttrTypes}, items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.InstanceSet = list

	// The state id is a random token (JSON resource_id = fwHelper.BuildToken()).
	config.Id = types.StringValue(fwHelper.BuildToken())

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// flattenAgsSandboxInstance maps one SDK sandbox instance onto the item model.
func flattenAgsSandboxInstance(ctx context.Context, item *agsv20250920.SandboxInstance) (AgsSandboxInstancesItemModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	itemModel := AgsSandboxInstancesItemModel{
		InstanceId:            fwHelper.StringPointerValueOrNull(item.InstanceId),
		ToolId:                fwHelper.StringPointerValueOrNull(item.ToolId),
		ToolName:              fwHelper.StringPointerValueOrNull(item.ToolName),
		Status:                fwHelper.StringPointerValueOrNull(item.Status),
		Persistent:            fwHelper.BoolValueOrNull(item.Persistent),
		TimeoutSeconds:        fwHelper.Int64ValueFromUint(item.TimeoutSeconds),
		ExpiresAt:             fwHelper.StringPointerValueOrNull(item.ExpiresAt),
		StopReason:            fwHelper.StringPointerValueOrNull(item.StopReason),
		CreateTime:            fwHelper.StringPointerValueOrNull(item.CreateTime),
		UpdateTime:            fwHelper.StringPointerValueOrNull(item.UpdateTime),
		MountOptions:          types.ListNull(types.ObjectType{AttrTypes: agsSandboxInstancesMountOptionAttrTypes}),
		CustomConfiguration:   types.ObjectNull(agsSandboxInstancesCustomConfigurationAttrTypes),
		ComputerConfiguration: types.ObjectNull(agsSandboxInstancesComputerConfigurationAttrTypes),
		NetworkMode:           fwHelper.StringPointerValueOrNull(item.NetworkMode),
		Metadata:              types.ListNull(types.ObjectType{AttrTypes: agsSandboxInstancesMetadataAttrTypes}),
		AuthMode:              fwHelper.StringPointerValueOrNull(item.AuthMode),
	}

	mountOptions := make([]agsSandboxInstancesMountOptionModel, 0, len(item.MountOptions))
	for _, mountOption := range item.MountOptions {
		if mountOption == nil {
			continue
		}

		mountOptions = append(mountOptions, agsSandboxInstancesMountOptionModel{
			Name:      fwHelper.StringPointerValueOrNull(mountOption.Name),
			MountPath: fwHelper.StringPointerValueOrNull(mountOption.MountPath),
			SubPath:   fwHelper.StringPointerValueOrNull(mountOption.SubPath),
			ReadOnly:  fwHelper.BoolValueOrNull(mountOption.ReadOnly),
		})
	}

	if len(mountOptions) > 0 {
		mountOptionsList, listDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: agsSandboxInstancesMountOptionAttrTypes}, mountOptions)
		diags.Append(listDiags...)
		if diags.HasError() {
			return itemModel, diags
		}

		itemModel.MountOptions = mountOptionsList
	}

	if item.CustomConfiguration != nil {
		customConfiguration, customDiags := flattenAgsSandboxInstancesCustomConfiguration(ctx, item.CustomConfiguration)
		diags.Append(customDiags...)
		if diags.HasError() {
			return itemModel, diags
		}

		itemModel.CustomConfiguration = customConfiguration
	}

	if item.ComputerConfiguration != nil {
		computerConfiguration, computerDiags := flattenAgsSandboxInstancesComputerConfiguration(ctx, item.ComputerConfiguration)
		diags.Append(computerDiags...)
		if diags.HasError() {
			return itemModel, diags
		}

		itemModel.ComputerConfiguration = computerConfiguration
	}

	metadata := make([]agsSandboxInstancesMetadataModel, 0, len(item.Metadata))
	for _, metadataVar := range item.Metadata {
		if metadataVar == nil {
			continue
		}

		metadata = append(metadata, agsSandboxInstancesMetadataModel{
			Name:  fwHelper.StringPointerValueOrNull(metadataVar.Name),
			Value: fwHelper.StringPointerValueOrNull(metadataVar.Value),
		})
	}

	if len(metadata) > 0 {
		metadataList, listDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: agsSandboxInstancesMetadataAttrTypes}, metadata)
		diags.Append(listDiags...)
		if diags.HasError() {
			return itemModel, diags
		}

		itemModel.Metadata = metadataList
	}

	return itemModel, diags
}

// flattenAgsSandboxInstancesCustomConfiguration maps the SDK custom
// configuration onto the nested object value.
func flattenAgsSandboxInstancesCustomConfiguration(ctx context.Context, item *agsv20250920.CustomConfigurationDetail) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := agsSandboxInstancesCustomConfigurationModel{
		Image:             fwHelper.StringPointerValueOrNull(item.Image),
		ImageRegistryType: fwHelper.StringPointerValueOrNull(item.ImageRegistryType),
		ImageDigest:       fwHelper.StringPointerValueOrNull(item.ImageDigest),
		Command:           fwHelper.ListFromPointers[string](ctx, types.StringType, item.Command, &diags),
		Args:              fwHelper.ListFromPointers[string](ctx, types.StringType, item.Args, &diags),
		Env:               types.ListNull(types.ObjectType{AttrTypes: agsSandboxInstancesEnvVarAttrTypes}),
		Ports:             types.ListNull(types.ObjectType{AttrTypes: agsSandboxInstancesPortAttrTypes}),
		Resources:         types.ObjectNull(agsSandboxInstancesResourcesAttrTypes),
		Probe:             types.ObjectNull(agsSandboxInstancesProbeAttrTypes),
		DNSConfig:         types.ObjectNull(agsSandboxInstancesDNSConfigAttrTypes),
		SecurityContext:   types.ObjectNull(agsSandboxInstancesSecurityContextAttrTypes),
	}
	if diags.HasError() {
		return types.ObjectNull(agsSandboxInstancesCustomConfigurationAttrTypes), diags
	}

	env := make([]agsSandboxInstancesEnvVarModel, 0, len(item.Env))
	for _, envVar := range item.Env {
		if envVar == nil {
			continue
		}

		env = append(env, agsSandboxInstancesEnvVarModel{
			Name:  fwHelper.StringPointerValueOrNull(envVar.Name),
			Value: fwHelper.StringPointerValueOrNull(envVar.Value),
		})
	}

	if len(env) > 0 {
		envList, listDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: agsSandboxInstancesEnvVarAttrTypes}, env)
		diags.Append(listDiags...)
		if diags.HasError() {
			return types.ObjectNull(agsSandboxInstancesCustomConfigurationAttrTypes), diags
		}

		model.Env = envList
	}

	ports := make([]agsSandboxInstancesPortModel, 0, len(item.Ports))
	for _, port := range item.Ports {
		if port == nil {
			continue
		}

		ports = append(ports, agsSandboxInstancesPortModel{
			Name:     fwHelper.StringPointerValueOrNull(port.Name),
			Port:     fwHelper.Int64ValueOrNull(port.Port),
			Protocol: fwHelper.StringPointerValueOrNull(port.Protocol),
		})
	}

	if len(ports) > 0 {
		portsList, listDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: agsSandboxInstancesPortAttrTypes}, ports)
		diags.Append(listDiags...)
		if diags.HasError() {
			return types.ObjectNull(agsSandboxInstancesCustomConfigurationAttrTypes), diags
		}

		model.Ports = portsList
	}

	if item.Resources != nil {
		resources := agsSandboxInstancesResourcesModel{
			CPU:     fwHelper.StringPointerValueOrNull(item.Resources.CPU),
			Memory:  fwHelper.StringPointerValueOrNull(item.Resources.Memory),
			Storage: fwHelper.StringPointerValueOrNull(item.Resources.Storage),
		}

		resourcesObject, objectDiags := types.ObjectValueFrom(ctx, agsSandboxInstancesResourcesAttrTypes, resources)
		diags.Append(objectDiags...)
		if diags.HasError() {
			return types.ObjectNull(agsSandboxInstancesCustomConfigurationAttrTypes), diags
		}

		model.Resources = resourcesObject
	}

	if item.Probe != nil {
		probe, probeDiags := flattenAgsSandboxInstancesProbe(ctx, item.Probe)
		diags.Append(probeDiags...)
		if diags.HasError() {
			return types.ObjectNull(agsSandboxInstancesCustomConfigurationAttrTypes), diags
		}

		model.Probe = probe
	}

	if item.DNSConfig != nil {
		dnsConfig := agsSandboxInstancesDNSConfigModel{
			Servers:  fwHelper.ListFromPointers[string](ctx, types.StringType, item.DNSConfig.Servers, &diags),
			Searches: fwHelper.ListFromPointers[string](ctx, types.StringType, item.DNSConfig.Searches, &diags),
			Options:  fwHelper.ListFromPointers[string](ctx, types.StringType, item.DNSConfig.Options, &diags),
		}
		if diags.HasError() {
			return types.ObjectNull(agsSandboxInstancesCustomConfigurationAttrTypes), diags
		}

		dnsConfigObject, objectDiags := types.ObjectValueFrom(ctx, agsSandboxInstancesDNSConfigAttrTypes, dnsConfig)
		diags.Append(objectDiags...)
		if diags.HasError() {
			return types.ObjectNull(agsSandboxInstancesCustomConfigurationAttrTypes), diags
		}

		model.DNSConfig = dnsConfigObject
	}

	if item.SecurityContext != nil {
		securityContext := agsSandboxInstancesSecurityContextModel{
			Privileged:   fwHelper.BoolValueOrNull(item.SecurityContext.Privileged),
			Capabilities: types.ObjectNull(agsSandboxInstancesCapabilitiesAttrTypes),
		}

		if item.SecurityContext.Capabilities != nil {
			capabilities := agsSandboxInstancesCapabilitiesModel{
				Add:  fwHelper.ListFromPointers[string](ctx, types.StringType, item.SecurityContext.Capabilities.Add, &diags),
				Drop: fwHelper.ListFromPointers[string](ctx, types.StringType, item.SecurityContext.Capabilities.Drop, &diags),
			}
			if diags.HasError() {
				return types.ObjectNull(agsSandboxInstancesCustomConfigurationAttrTypes), diags
			}

			capabilitiesObject, objectDiags := types.ObjectValueFrom(ctx, agsSandboxInstancesCapabilitiesAttrTypes, capabilities)
			diags.Append(objectDiags...)
			if diags.HasError() {
				return types.ObjectNull(agsSandboxInstancesCustomConfigurationAttrTypes), diags
			}

			securityContext.Capabilities = capabilitiesObject
		}

		securityContextObject, objectDiags := types.ObjectValueFrom(ctx, agsSandboxInstancesSecurityContextAttrTypes, securityContext)
		diags.Append(objectDiags...)
		if diags.HasError() {
			return types.ObjectNull(agsSandboxInstancesCustomConfigurationAttrTypes), diags
		}

		model.SecurityContext = securityContextObject
	}

	return types.ObjectValueFrom(ctx, agsSandboxInstancesCustomConfigurationAttrTypes, model)
}

// flattenAgsSandboxInstancesProbe maps the SDK probe configuration onto the
// nested object value.
func flattenAgsSandboxInstancesProbe(ctx context.Context, item *agsv20250920.ProbeConfiguration) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := agsSandboxInstancesProbeModel{
		HttpGet:          types.ObjectNull(agsSandboxInstancesHttpGetAttrTypes),
		ReadyTimeoutMs:   fwHelper.Int64ValueOrNull(item.ReadyTimeoutMs),
		ProbeTimeoutMs:   fwHelper.Int64ValueOrNull(item.ProbeTimeoutMs),
		ProbePeriodMs:    fwHelper.Int64ValueOrNull(item.ProbePeriodMs),
		SuccessThreshold: fwHelper.Int64ValueOrNull(item.SuccessThreshold),
		FailureThreshold: fwHelper.Int64ValueOrNull(item.FailureThreshold),
	}

	if item.HttpGet != nil {
		httpGet := agsSandboxInstancesHttpGetModel{
			Path:   fwHelper.StringPointerValueOrNull(item.HttpGet.Path),
			Port:   fwHelper.Int64ValueOrNull(item.HttpGet.Port),
			Scheme: fwHelper.StringPointerValueOrNull(item.HttpGet.Scheme),
		}

		httpGetObject, objectDiags := types.ObjectValueFrom(ctx, agsSandboxInstancesHttpGetAttrTypes, httpGet)
		diags.Append(objectDiags...)
		if diags.HasError() {
			return types.ObjectNull(agsSandboxInstancesProbeAttrTypes), diags
		}

		model.HttpGet = httpGetObject
	}

	return types.ObjectValueFrom(ctx, agsSandboxInstancesProbeAttrTypes, model)
}

// flattenAgsSandboxInstancesComputerConfiguration maps the SDK computer
// configuration onto the nested object value.
func flattenAgsSandboxInstancesComputerConfiguration(ctx context.Context, item *agsv20250920.ComputerConfiguration) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := agsSandboxInstancesComputerConfigurationModel{
		WAAConfiguration:     types.ObjectNull(agsSandboxInstancesWAAConfigurationAttrTypes),
		OSWorldConfiguration: types.ObjectNull(agsSandboxInstancesOSWorldConfigurationAttrTypes),
	}

	if item.WAAConfiguration != nil {
		waaConfiguration := agsSandboxInstancesWAAConfigurationModel{
			ImageId: fwHelper.StringPointerValueOrNull(item.WAAConfiguration.ImageId),
		}

		waaConfigurationObject, objectDiags := types.ObjectValueFrom(ctx, agsSandboxInstancesWAAConfigurationAttrTypes, waaConfiguration)
		diags.Append(objectDiags...)
		if diags.HasError() {
			return types.ObjectNull(agsSandboxInstancesComputerConfigurationAttrTypes), diags
		}

		model.WAAConfiguration = waaConfigurationObject
	}

	if item.OSWorldConfiguration != nil {
		osWorldConfiguration := agsSandboxInstancesOSWorldConfigurationModel{
			Version: fwHelper.StringPointerValueOrNull(item.OSWorldConfiguration.Version),
		}

		osWorldConfigurationObject, objectDiags := types.ObjectValueFrom(ctx, agsSandboxInstancesOSWorldConfigurationAttrTypes, osWorldConfiguration)
		diags.Append(objectDiags...)
		if diags.HasError() {
			return types.ObjectNull(agsSandboxInstancesComputerConfigurationAttrTypes), diags
		}

		model.OSWorldConfiguration = osWorldConfigurationObject
	}

	return types.ObjectValueFrom(ctx, agsSandboxInstancesComputerConfigurationAttrTypes, model)
}
