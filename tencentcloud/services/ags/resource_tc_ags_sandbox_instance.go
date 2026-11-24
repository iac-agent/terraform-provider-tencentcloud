package ags

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

const agsSandboxInstancePollInterval = 3 * time.Second

const (
	AGS_SANDBOX_INSTANCE_STATUS_STARTING = "STARTING"
	AGS_SANDBOX_INSTANCE_STATUS_RUNNING  = "RUNNING"
	AGS_SANDBOX_INSTANCE_STATUS_FAILED   = "FAILED"
)

var (
	_ resource.ResourceWithConfigure   = &AgsSandboxInstanceResource{}
	_ resource.ResourceWithImportState = &AgsSandboxInstanceResource{}
)

// NewAgsSandboxInstanceResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewAgsSandboxInstanceResource() resource.Resource {
	return &AgsSandboxInstanceResource{}
}

// AgsSandboxInstanceResource implements resource.Resource for
// tencentcloud_ags_sandbox_instance.
type AgsSandboxInstanceResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// AgsSandboxInstanceResourceModel maps the schema attributes.
//
// Timeout and ClientToken are create-only inputs: the describe API reports the
// timeout in seconds (TimeoutSeconds) and never echoes back a client token, so
// Read keeps whatever value is already stored in state.
type AgsSandboxInstanceResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	InstanceId          types.String `tfsdk:"instance_id"`
	ToolId              types.String `tfsdk:"tool_id"`
	ToolName            types.String `tfsdk:"tool_name"`
	SourceInstanceId    types.String `tfsdk:"source_instance_id"`
	SnapshotId          types.String `tfsdk:"snapshot_id"`
	Timeout             types.String `tfsdk:"timeout"`
	ClientToken         types.String `tfsdk:"client_token"`
	MountOptions        types.List   `tfsdk:"mount_options"`
	CustomConfiguration types.Object `tfsdk:"custom_configuration"`
	AuthMode            types.String `tfsdk:"auth_mode"`
	Metadata            types.List   `tfsdk:"metadata"`
	Status              types.String `tfsdk:"status"`
	Persistent          types.Bool   `tfsdk:"persistent"`
	TimeoutSeconds      types.Int64  `tfsdk:"timeout_seconds"`
	ExpiresAt           types.String `tfsdk:"expires_at"`
	StopReason          types.String `tfsdk:"stop_reason"`
	NetworkMode         types.String `tfsdk:"network_mode"`
	CreateTime          types.String `tfsdk:"create_time"`
	UpdateTime          types.String `tfsdk:"update_time"`
}

type agsMountOptionModel struct {
	Name      types.String `tfsdk:"name"`
	MountPath types.String `tfsdk:"mount_path"`
	SubPath   types.String `tfsdk:"sub_path"`
	ReadOnly  types.Bool   `tfsdk:"read_only"`
}

type agsMetadataModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

var agsMountOptionAttributes = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Optional:    true,
		Description: "Name of the storage configuration in the sandbox tool.",
	},
	"mount_path": schema.StringAttribute{
		Optional:    true,
		Description: "Local mount path inside the sandbox instance. When it is not specified, the storage configuration of the tool is inherited.",
	},
	"sub_path": schema.StringAttribute{
		Optional:    true,
		Description: "Sub path of the storage mount inside the sandbox instance.",
	},
	"read_only": schema.BoolAttribute{
		Optional:    true,
		Description: "Whether the storage mount is read only. When it is not specified, the storage configuration of the tool is inherited.",
	},
}

var agsMetadataAttributes = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Optional:    true,
		Description: "Metadata name.",
	},
	"value": schema.StringAttribute{
		Optional:    true,
		Description: "Metadata value.",
	},
}

var (
	agsMountOptionAttrTypes = attributeTypesOf(agsMountOptionAttributes)
	agsMetadataAttrTypes    = attributeTypesOf(agsMetadataAttributes)

	agsMountOptionObjectType = types.ObjectType{AttrTypes: agsMountOptionAttrTypes}
	agsMetadataObjectType    = types.ObjectType{AttrTypes: agsMetadataAttrTypes}
)

func (r *AgsSandboxInstanceResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_ags_sandbox_instance"
}

func (r *AgsSandboxInstanceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to start an AGS sandbox instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Id of the resource, which is the same as `instance_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Computed:    true,
				Description: "Id of the sandbox instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tool_id": schema.StringAttribute{
				Optional:    true,
				Description: "Id of the sandbox tool. Either `tool_id` or `tool_name` must be specified.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tool_name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the sandbox tool. Either `tool_id` or `tool_name` must be specified.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "Id of the source sandbox instance to clone from.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"snapshot_id": schema.StringAttribute{
				Optional:    true,
				Description: "Id of the sandbox snapshot to start from.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"timeout": schema.StringAttribute{
				Optional: true,
				Description: "Timeout after which the instance is automatically reclaimed. Supported formats are `5m`, `300s` and `1h`. " +
					"Default value: `5m`. Minimum 30s, maximum 24h. It is not reported by the describe API and cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client_token": schema.StringAttribute{
				Optional: true,
				Description: "Idempotency token with a maximum length of 64 characters. When it is not specified, a uuid is generated " +
					"internally so that a retried create call cannot create the same instance twice. It is not reported by the describe API.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mount_options": schema.ListNestedAttribute{
				Optional:     true,
				NestedObject: schema.NestedAttributeObject{Attributes: agsMountOptionAttributes},
				Description:  "Storage mount configuration of the sandbox instance.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"custom_configuration": schema.SingleNestedAttribute{
				Optional:    true,
				Attributes:  agsCustomConfigurationAttributes,
				Description: "Custom configuration of the sandbox instance.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"auth_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Authentication mode of sandbox access. Valid values: `DEFAULT`, `TOKEN`, `NONE`, `PUBLIC`. " +
					"Default value: `DEFAULT`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"metadata": schema.ListNestedAttribute{
				Optional:     true,
				NestedObject: schema.NestedAttributeObject{Attributes: agsMetadataAttributes},
				Description:  "Metadata of the sandbox instance.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed: true,
				Description: "Status of the sandbox instance. Valid values: `STARTING`, `RUNNING`, `STOPPING`, `STOPPED`, " +
					"`STOP_FAILED`, `FAILED`, `RESIZING`.",
			},
			"persistent": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the instance is persistent.",
			},
			"timeout_seconds": schema.Int64Attribute{
				Computed:    true,
				Description: "Timeout of the instance in seconds. A null value means no timeout is set.",
			},
			"expires_at": schema.StringAttribute{
				Computed:    true,
				Description: "Expiration time of the instance in ISO8601 format. A null value means no expiration time.",
			},
			"stop_reason": schema.StringAttribute{
				Computed:    true,
				Description: "Reason why the instance was stopped: `manual`, `timeout`, `error` or `system`. It is only set when the status is `STOPPED`, `STOP_FAILED` or `FAILED`.",
			},
			"network_mode": schema.StringAttribute{
				Computed:    true,
				Description: "Network mode of the instance. Valid values: `PUBLIC`, `SANDBOX`, `INTERNAL_SERVICE`.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the sandbox instance in ISO8601 format.",
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "Last update time of the sandbox instance in ISO8601 format.",
			},
		},
	}
}

// Create handles `terraform apply` for new resources. StartSandboxInstance has
// no in-place update API, so every configurable attribute is marked with
// RequiresReplace and Update is a no-op.
func (r *AgsSandboxInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AgsSandboxInstanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request, diags := r.buildCreateRequest(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewAgsService(r.Client())
	instance, err := service.StartSandboxInstance(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error starting AGS sandbox instance",
			fmt.Sprintf("StartSandboxInstance failed: %v", err),
		)
		return
	}

	instanceId := *instance.InstanceId
	plan.Id = types.StringValue(instanceId)
	plan.InstanceId = types.StringValue(instanceId)
	plan.Status = types.StringNull()
	plan.Persistent = types.BoolNull()
	plan.TimeoutSeconds = types.Int64Null()
	plan.ExpiresAt = types.StringNull()
	plan.StopReason = types.StringNull()
	plan.NetworkMode = types.StringNull()
	plan.CreateTime = types.StringNull()
	plan.UpdateTime = types.StringNull()
	r.normalizeUnknownAttributes(&plan)

	// Persist the identity before waiting: if the wait below fails, the resource
	// is still in state and `terraform destroy` can stop it.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.waitSandboxInstanceRunning(ctx, service, instanceId, tccommon.WriteRetryTimeout)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"AGS sandbox instance is not running yet",
			fmt.Sprintf("sandbox instance %s was started but is not running yet: %v", instanceId, err),
		)
		return
	}

	resp.Diagnostics.Append(r.refreshModel(ctx, &plan, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and the post-apply state refresh.
func (r *AgsSandboxInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AgsSandboxInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId, diags := sandboxInstanceIdOf(&state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewAgsService(r.Client())
	item, err := service.DescribeSandboxInstanceById(ctx, instanceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading AGS sandbox instance",
			fmt.Sprintf("DescribeSandboxInstances failed for %s: %v", instanceId, err),
		)
		return
	}

	// A nil result means the query succeeded and the instance no longer exists. A
	// nil API response is handled as an error inside the service and never
	// reaches this branch, so state is not dropped on a failed request.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"AGS sandbox instance not found",
			fmt.Sprintf("sandbox instance %s does not exist any more, removing it from state.", instanceId),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(r.refreshModel(ctx, &state, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update is a no-op: StartSandboxInstance has no matching update API and every
// configurable attribute is marked with RequiresReplace in the schema.
func (r *AgsSandboxInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AgsSandboxInstanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles `terraform destroy`. It stops the sandbox instance; the cloud
// side reclaims the instance once it is stopped.
func (r *AgsSandboxInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AgsSandboxInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceId, diags := sandboxInstanceIdOf(&state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewAgsService(r.Client())
	if err := service.StopSandboxInstance(ctx, instanceId); err != nil {
		// The API documentation does not define an "already deleted" error code,
		// so existence is checked again instead of guessing codes: an instance
		// that is really gone makes the delete idempotent, everything else is an
		// error.
		item, describeErr := service.DescribeSandboxInstanceById(ctx, instanceId)
		if describeErr != nil || item != nil {
			resp.Diagnostics.AddError(
				"Error deleting AGS sandbox instance",
				fmt.Sprintf("StopSandboxInstance failed for %s: %v", instanceId, err),
			)
		}
		return
	}
}

// ImportState supports `terraform import` with the sandbox instance id.
func (r *AgsSandboxInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError(
			"Invalid import id",
			"The import id must be the id of an existing sandbox instance.",
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_id"), req.ID)...)
}

// buildCreateRequest maps the plan onto a StartSandboxInstance request.
func (r *AgsSandboxInstanceResource) buildCreateRequest(ctx context.Context, plan *AgsSandboxInstanceResourceModel) (*agsv20250920.StartSandboxInstanceRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	request := agsv20250920.NewStartSandboxInstanceRequest()
	if !plan.ToolId.IsNull() && !plan.ToolId.IsUnknown() {
		request.ToolId = plan.ToolId.ValueStringPointer()
	}
	if !plan.ToolName.IsNull() && !plan.ToolName.IsUnknown() {
		request.ToolName = plan.ToolName.ValueStringPointer()
	}
	if !plan.SourceInstanceId.IsNull() && !plan.SourceInstanceId.IsUnknown() {
		request.SourceInstanceId = plan.SourceInstanceId.ValueStringPointer()
	}
	if !plan.SnapshotId.IsNull() && !plan.SnapshotId.IsUnknown() {
		request.SnapshotId = plan.SnapshotId.ValueStringPointer()
	}
	if !plan.Timeout.IsNull() && !plan.Timeout.IsUnknown() {
		request.Timeout = plan.Timeout.ValueStringPointer()
	}
	// When the user does not configure client_token a uuid is generated as an
	// idempotency key: the value is fixed once per apply, so a retried
	// StartSandboxInstance call (see the retry loop in the service layer) cannot
	// create the same instance twice.
	if !plan.ClientToken.IsNull() && !plan.ClientToken.IsUnknown() && plan.ClientToken.ValueString() != "" {
		request.ClientToken = plan.ClientToken.ValueStringPointer()
	} else {
		request.ClientToken = helper.String(helper.BuildUUID())
	}

	mountOptions, mountDiags := expandMountOptions(ctx, plan.MountOptions)
	diags.Append(mountDiags...)
	if diags.HasError() {
		return nil, diags
	}
	request.MountOptions = mountOptions

	customConfiguration, customDiags := expandCustomConfiguration(ctx, plan.CustomConfiguration)
	diags.Append(customDiags...)
	if diags.HasError() {
		return nil, diags
	}
	request.CustomConfiguration = customConfiguration

	if !plan.AuthMode.IsNull() && !plan.AuthMode.IsUnknown() {
		request.AuthMode = plan.AuthMode.ValueStringPointer()
	}

	metadata, metadataDiags := expandMetadata(ctx, plan.Metadata)
	diags.Append(metadataDiags...)
	if diags.HasError() {
		return nil, diags
	}
	request.Metadata = metadata

	return request, diags
}

// refreshModel copies the cloud state into the model. Only fields reported by
// the cloud are overwritten, so create-only attributes that cannot be read back
// (timeout, client_token) keep the value already stored in the model.
func (r *AgsSandboxInstanceResource) refreshModel(ctx context.Context, model *AgsSandboxInstanceResourceModel, item *agsv20250920.SandboxInstance) diag.Diagnostics {
	var diags diag.Diagnostics

	if item == nil {
		return diags
	}

	if item.InstanceId != nil && *item.InstanceId != "" {
		model.InstanceId = types.StringValue(*item.InstanceId)
		model.Id = types.StringValue(*item.InstanceId)
	}
	if item.ToolId != nil {
		model.ToolId = types.StringValue(*item.ToolId)
	}
	if item.ToolName != nil {
		model.ToolName = types.StringValue(*item.ToolName)
	}
	model.Status = stringPointerValue(item.Status)
	model.Persistent = boolPointerValue(item.Persistent)
	model.ExpiresAt = stringPointerValue(item.ExpiresAt)
	model.StopReason = stringPointerValue(item.StopReason)
	model.NetworkMode = stringPointerValue(item.NetworkMode)
	model.CreateTime = stringPointerValue(item.CreateTime)
	model.UpdateTime = stringPointerValue(item.UpdateTime)
	if item.TimeoutSeconds != nil {
		model.TimeoutSeconds = types.Int64Value(int64(*item.TimeoutSeconds))
	} else {
		model.TimeoutSeconds = types.Int64Null()
	}
	if item.AuthMode != nil {
		model.AuthMode = types.StringValue(*item.AuthMode)
	}

	mountOptions, mountDiags := flattenMountOptions(ctx, item.MountOptions)
	diags.Append(mountDiags...)
	if diags.HasError() {
		return diags
	}
	model.MountOptions = mountOptions

	customConfiguration, customDiags := flattenCustomConfiguration(ctx, item.CustomConfiguration)
	diags.Append(customDiags...)
	if diags.HasError() {
		return diags
	}
	model.CustomConfiguration = customConfiguration

	metadata, metadataDiags := flattenMetadata(ctx, item.Metadata)
	diags.Append(metadataDiags...)
	if diags.HasError() {
		return diags
	}
	model.Metadata = metadata

	return diags
}

// waitSandboxInstanceRunning waits until the sandbox instance reaches RUNNING.
// STARTING is a valid intermediate state, and the instance is still usable for
// destroy when it is stuck, so a timeout is reported as a warning instead of an
// error and the state stored before waiting is kept.
func (r *AgsSandboxInstanceResource) waitSandboxInstanceRunning(ctx context.Context, service AgsService, instanceId string, timeout time.Duration) (*agsv20250920.SandboxInstance, error) {
	deadline := time.Now().Add(timeout)
	for {
		item, err := service.DescribeSandboxInstanceById(ctx, instanceId)
		if err != nil {
			return nil, err
		}
		if item == nil {
			return nil, fmt.Errorf("sandbox instance %s disappeared while waiting for it to become running", instanceId)
		}

		status := ""
		if item.Status != nil {
			status = *item.Status
		}
		switch status {
		case AGS_SANDBOX_INSTANCE_STATUS_RUNNING:
			return item, nil
		case AGS_SANDBOX_INSTANCE_STATUS_FAILED:
			return nil, fmt.Errorf("sandbox instance %s is in FAILED status", instanceId)
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timeout waiting for sandbox instance %s to become RUNNING, current status is %q", instanceId, status)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(agsSandboxInstancePollInterval):
		}
	}
}

// normalizeUnknownAttributes replaces the optional+computed attributes that are
// still unknown at create time with null, so the identity can be stored in state
// before the cloud has reported those values.
func (r *AgsSandboxInstanceResource) normalizeUnknownAttributes(model *AgsSandboxInstanceResourceModel) {
	if model.AuthMode.IsUnknown() {
		model.AuthMode = types.StringNull()
	}
	if model.MountOptions.IsUnknown() {
		model.MountOptions = types.ListNull(agsMountOptionObjectType)
	}
	if model.CustomConfiguration.IsUnknown() {
		model.CustomConfiguration = types.ObjectNull(agsCustomConfigAttrTypes)
	}
	if model.Metadata.IsUnknown() {
		model.Metadata = types.ListNull(agsMetadataObjectType)
	}
}

// sandboxInstanceIdOf returns the identity stored in the model.
func sandboxInstanceIdOf(model *AgsSandboxInstanceResourceModel) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !model.InstanceId.IsNull() && !model.InstanceId.IsUnknown() && model.InstanceId.ValueString() != "" {
		return model.InstanceId.ValueString(), diags
	}
	if !model.Id.IsNull() && !model.Id.IsUnknown() && model.Id.ValueString() != "" {
		return model.Id.ValueString(), diags
	}

	diags.AddError(
		"Invalid AGS sandbox instance state",
		"The resource state does not contain a sandbox instance id.",
	)
	return "", diags
}

func expandMountOptions(ctx context.Context, list types.List) ([]*agsv20250920.MountOption, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var models []agsMountOptionModel
	diags := list.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, diags
	}
	if len(models) == 0 {
		return nil, diags
	}

	mountOptions := make([]*agsv20250920.MountOption, 0, len(models))
	for _, model := range models {
		mountOption := &agsv20250920.MountOption{}
		if !model.Name.IsNull() && !model.Name.IsUnknown() {
			mountOption.Name = model.Name.ValueStringPointer()
		}
		if !model.MountPath.IsNull() && !model.MountPath.IsUnknown() {
			mountOption.MountPath = model.MountPath.ValueStringPointer()
		}
		if !model.SubPath.IsNull() && !model.SubPath.IsUnknown() {
			mountOption.SubPath = model.SubPath.ValueStringPointer()
		}
		if !model.ReadOnly.IsNull() && !model.ReadOnly.IsUnknown() {
			mountOption.ReadOnly = model.ReadOnly.ValueBoolPointer()
		}
		mountOptions = append(mountOptions, mountOption)
	}

	return mountOptions, diags
}

func flattenMountOptions(ctx context.Context, items []*agsv20250920.MountOption) (types.List, diag.Diagnostics) {
	if len(items) == 0 {
		return types.ListNull(agsMountOptionObjectType), nil
	}

	models := make([]agsMountOptionModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		models = append(models, agsMountOptionModel{
			Name:      stringPointerValue(item.Name),
			MountPath: stringPointerValue(item.MountPath),
			SubPath:   stringPointerValue(item.SubPath),
			ReadOnly:  boolPointerValue(item.ReadOnly),
		})
	}
	if len(models) == 0 {
		return types.ListNull(agsMountOptionObjectType), nil
	}

	return types.ListValueFrom(ctx, agsMountOptionObjectType, models)
}

func expandMetadata(ctx context.Context, list types.List) ([]*agsv20250920.MetadataVar, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var models []agsMetadataModel
	diags := list.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, diags
	}
	if len(models) == 0 {
		return nil, diags
	}

	metadata := make([]*agsv20250920.MetadataVar, 0, len(models))
	for _, model := range models {
		metadataVar := &agsv20250920.MetadataVar{}
		if !model.Name.IsNull() && !model.Name.IsUnknown() {
			metadataVar.Name = model.Name.ValueStringPointer()
		}
		if !model.Value.IsNull() && !model.Value.IsUnknown() {
			metadataVar.Value = model.Value.ValueStringPointer()
		}
		metadata = append(metadata, metadataVar)
	}

	return metadata, diags
}

func flattenMetadata(ctx context.Context, items []*agsv20250920.MetadataVar) (types.List, diag.Diagnostics) {
	if len(items) == 0 {
		return types.ListNull(agsMetadataObjectType), nil
	}

	models := make([]agsMetadataModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		models = append(models, agsMetadataModel{
			Name:  stringPointerValue(item.Name),
			Value: stringPointerValue(item.Value),
		})
	}
	if len(models) == 0 {
		return types.ListNull(agsMetadataObjectType), nil
	}

	return types.ListValueFrom(ctx, agsMetadataObjectType, models)
}
