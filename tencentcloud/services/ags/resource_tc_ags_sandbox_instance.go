package ags

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
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

// agsSandboxInstancePollInterval is the interval used while waiting for a
// sandbox instance to reach a stable state.
const agsSandboxInstancePollInterval = 5 * time.Second

// Sandbox instance status values documented by DescribeSandboxInstances.
const (
	AGS_SANDBOX_INSTANCE_STATUS_STARTING = "STARTING"
	AGS_SANDBOX_INSTANCE_STATUS_RUNNING  = "RUNNING"
	AGS_SANDBOX_INSTANCE_STATUS_FAILED   = "FAILED"
	AGS_SANDBOX_INSTANCE_STATUS_STOPPED  = "STOPPED"
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
// Timeout, SourceInstanceId, SnapshotId and ClientToken are create-only
// inputs: DescribeSandboxInstances does not report them back in the original
// request form, so Read keeps whatever value is already stored in state.
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
	NetworkMode         types.String `tfsdk:"network_mode"`
	Persistent          types.Bool   `tfsdk:"persistent"`
	Status              types.String `tfsdk:"status"`
	TimeoutSeconds      types.Int64  `tfsdk:"timeout_seconds"`
	ExpiresAt           types.String `tfsdk:"expires_at"`
	StopReason          types.String `tfsdk:"stop_reason"`
	CreateTime          types.String `tfsdk:"create_time"`
	UpdateTime          types.String `tfsdk:"update_time"`
}

type agsMountOptionModel struct {
	Name      types.String `tfsdk:"name"`
	MountPath types.String `tfsdk:"mount_path"`
	SubPath   types.String `tfsdk:"sub_path"`
	ReadOnly  types.Bool   `tfsdk:"read_only"`
}

var agsMountOptionAttributes = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Optional:    true,
		Description: "Name of the storage configuration in the sandbox tool.",
	},
	"mount_path": schema.StringAttribute{
		Optional:    true,
		Description: "Local mount path inside the sandbox instance. By default it inherits the storage configuration of the tool.",
	},
	"sub_path": schema.StringAttribute{
		Optional:    true,
		Description: "Sub path of the storage mount inside the sandbox instance.",
	},
	"read_only": schema.BoolAttribute{
		Optional:    true,
		Description: "Whether the storage mount is read only. By default it inherits the storage configuration of the tool.",
	},
}

var (
	agsMountOptionAttrTypes = attributeTypesOf(agsMountOptionAttributes)

	agsMountOptionObjectType = types.ObjectType{AttrTypes: agsMountOptionAttrTypes}
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
				Optional: true,
				Description: "Id of the sandbox tool. Either `tool_id` or `tool_name` must be specified. " +
					"It cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tool_name": schema.StringAttribute{
				Optional: true,
				Description: "Name of the sandbox tool. Either `tool_id` or `tool_name` must be specified. " +
					"It cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "Id of the source sandbox instance to clone from. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"snapshot_id": schema.StringAttribute{
				Optional:    true,
				Description: "Id of the sandbox snapshot to start from. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"timeout": schema.StringAttribute{
				Optional: true,
				Description: "Timeout after which the instance is recycled automatically, for example `5m`, `300s` " +
					"or `1h`. Default value: `5m`. The minimum is 30 seconds and the maximum is 24 hours. " +
					"It cannot be changed after creation and it is not reported back by the describe API.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client_token": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Idempotency token with a maximum length of 64 characters. When it is not specified, " +
					"a generated uuid is used so a retried create cannot start the same instance twice.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mount_options": schema.ListNestedAttribute{
				Optional:     true,
				NestedObject: schema.NestedAttributeObject{Attributes: agsMountOptionAttributes},
				Description:  "Storage mount configuration of the sandbox instance. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"custom_configuration": schema.SingleNestedAttribute{
				Optional:    true,
				Attributes:  agsCustomConfigurationAttributes,
				Description: "Custom configuration of the sandbox instance. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"auth_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Authentication mode of the sandbox. Valid values: `DEFAULT`, `TOKEN`, `NONE`, `PUBLIC`. " +
					"Default value: `DEFAULT`. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"metadata": schema.ListNestedAttribute{
				Optional:     true,
				NestedObject: schema.NestedAttributeObject{Attributes: agsEnvVarAttributes},
				Description:  "Metadata of the sandbox instance. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"network_mode": schema.StringAttribute{
				Computed:    true,
				Description: "Network mode of the sandbox instance. Valid values: `PUBLIC`, `SANDBOX`, `INTERNAL_SERVICE`.",
			},
			"persistent": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the instance is persistent.",
			},
			"status": schema.StringAttribute{
				Computed: true,
				Description: "Status of the sandbox instance. Valid values: `STARTING`, `RUNNING`, `STOPPING`, " +
					"`STOPPED`, `STOP_FAILED`, `FAILED`, `RESIZING`.",
			},
			"timeout_seconds": schema.Int64Attribute{
				Computed:    true,
				Description: "Timeout of the instance in seconds. It is null when no timeout is configured.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"expires_at": schema.StringAttribute{
				Computed:    true,
				Description: "Expiration time of the instance in ISO 8601 format. It is null when no timeout is configured.",
			},
			"stop_reason": schema.StringAttribute{
				Computed: true,
				Description: "Reason why the instance was stopped. Valid values: `manual`, `timeout`, `error`, " +
					"`system`. It only has a value when the status is `STOPPED`, `STOP_FAILED` or `FAILED`.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the sandbox instance in ISO 8601 format.",
			},
			"update_time": schema.StringAttribute{
				Computed:    true,
				Description: "Last update time of the sandbox instance in ISO 8601 format.",
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
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
	instanceId, err := service.StartSandboxInstance(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error starting AGS sandbox instance",
			fmt.Sprintf("StartSandboxInstance failed: %v", err),
		)
		return
	}

	plan.Id = types.StringValue(instanceId)
	plan.InstanceId = types.StringValue(instanceId)
	r.normalizeUnknownAttributes(&plan)

	// Persist the identity before waiting: if the wait below fails, the resource
	// is still in state and `terraform destroy` can stop it.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.waitSandboxInstanceRunning(ctx, service, instanceId, tccommon.WriteRetryTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for AGS sandbox instance to become running",
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

// Update handles in-place updates. Every configurable attribute of a sandbox
// instance is create-only and marked with RequiresReplace in the schema, so
// there is nothing to update in place.
func (r *AgsSandboxInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state AgsSandboxInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan AgsSandboxInstanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Id = state.Id
	plan.InstanceId = state.InstanceId
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles `terraform destroy`. Sandbox instances are recycled
// automatically when their timeout expires, so destroy only stops tracking the
// instance and stops it to avoid keeping it alive until then.
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
		// A vanished instance is already gone, which makes the delete
		// idempotent. Everything else is a real failure.
		item, describeErr := service.DescribeSandboxInstanceById(ctx, instanceId)
		if describeErr != nil || item != nil {
			resp.Diagnostics.AddError(
				"Error stopping AGS sandbox instance",
				fmt.Sprintf("StopSandboxInstance failed for %s: %v", instanceId, err),
			)
		}
		return
	}

	if err := r.waitSandboxInstanceGone(ctx, service, instanceId, tccommon.WriteRetryTimeout); err != nil {
		resp.Diagnostics.AddWarning(
			"AGS sandbox instance is still being stopped",
			fmt.Sprintf("sandbox instance %s was stopped but is still visible: %v", instanceId, err),
		)
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
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_id"), types.StringValue(req.ID))...)
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
	// client_token is used as an idempotency key for create: when it is not
	// configured a uuid is generated once per apply, so a retried
	// StartSandboxInstance call cannot start the same instance twice.
	if !plan.ClientToken.IsNull() && !plan.ClientToken.IsUnknown() {
		request.ClientToken = plan.ClientToken.ValueStringPointer()
	} else {
		request.ClientToken = helper.String(helper.BuildUUID())
	}
	if !plan.AuthMode.IsNull() && !plan.AuthMode.IsUnknown() {
		request.AuthMode = plan.AuthMode.ValueStringPointer()
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

	metadata, metadataDiags := expandMetadataVars(ctx, plan.Metadata)
	diags.Append(metadataDiags...)
	if diags.HasError() {
		return nil, diags
	}
	request.Metadata = metadata

	return request, diags
}

// refreshModel copies the cloud state into the model. Only fields reported by
// the cloud are overwritten, so create-only attributes that cannot be read back
// (timeout, source_instance_id, snapshot_id, client_token) keep the value
// already stored in the model.
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
	model.AuthMode = stringPointerValue(item.AuthMode)
	model.NetworkMode = stringPointerValue(item.NetworkMode)
	model.Persistent = boolPointerValue(item.Persistent)
	model.Status = stringPointerValue(item.Status)
	model.ExpiresAt = stringPointerValue(item.ExpiresAt)
	model.StopReason = stringPointerValue(item.StopReason)
	model.CreateTime = stringPointerValue(item.CreateTime)
	model.UpdateTime = stringPointerValue(item.UpdateTime)
	if item.TimeoutSeconds != nil {
		model.TimeoutSeconds = types.Int64Value(int64(*item.TimeoutSeconds))
	} else {
		model.TimeoutSeconds = types.Int64Null()
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

	metadata, metadataDiags := flattenMetadataVars(ctx, item.Metadata)
	diags.Append(metadataDiags...)
	if diags.HasError() {
		return diags
	}
	model.Metadata = metadata

	return diags
}

// waitSandboxInstanceRunning waits until the sandbox instance leaves the
// starting state. DescribeSandboxInstances documents STARTING, RUNNING,
// STOPPING, STOPPED, STOP_FAILED, FAILED and RESIZING as the possible status
// values.
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
		case AGS_SANDBOX_INSTANCE_STATUS_FAILED, AGS_SANDBOX_INSTANCE_STATUS_STOPPED:
			reason := ""
			if item.StopReason != nil {
				reason = *item.StopReason
			}
			return nil, fmt.Errorf("sandbox instance %s is in %s status: %s", instanceId, status, reason)
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

// waitSandboxInstanceGone waits until the stopped sandbox instance is no longer
// visible.
func (r *AgsSandboxInstanceResource) waitSandboxInstanceGone(ctx context.Context, service AgsService, instanceId string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		item, err := service.DescribeSandboxInstanceById(ctx, instanceId)
		if err != nil {
			return err
		}
		if item == nil {
			return nil
		}

		status := ""
		if item.Status != nil {
			status = *item.Status
		}
		if status == AGS_SANDBOX_INSTANCE_STATUS_STOPPED {
			return nil
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("sandbox instance %s is still present with status %q after %s", instanceId, status, timeout)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(agsSandboxInstancePollInterval):
		}
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

// normalizeUnknownAttributes replaces the optional+computed attributes that are
// still unknown at create time with null, so the identity can be stored in state
// before the cloud has reported those values.
func (r *AgsSandboxInstanceResource) normalizeUnknownAttributes(model *AgsSandboxInstanceResourceModel) {
	if model.AuthMode.IsUnknown() {
		model.AuthMode = types.StringNull()
	}
	if model.ClientToken.IsUnknown() {
		model.ClientToken = types.StringNull()
	}
	if model.MountOptions.IsUnknown() {
		model.MountOptions = types.ListNull(agsMountOptionObjectType)
	}
	if model.CustomConfiguration.IsUnknown() {
		model.CustomConfiguration = types.ObjectNull(agsCustomConfigAttrTypes)
	}
	if model.Metadata.IsUnknown() {
		model.Metadata = types.ListNull(agsEnvVarObjectType)
	}
	if model.Status.IsUnknown() {
		model.Status = types.StringNull()
	}
	if model.NetworkMode.IsUnknown() {
		model.NetworkMode = types.StringNull()
	}
	if model.Persistent.IsUnknown() {
		model.Persistent = types.BoolNull()
	}
	if model.TimeoutSeconds.IsUnknown() {
		model.TimeoutSeconds = types.Int64Null()
	}
	if model.ExpiresAt.IsUnknown() {
		model.ExpiresAt = types.StringNull()
	}
	if model.StopReason.IsUnknown() {
		model.StopReason = types.StringNull()
	}
	if model.CreateTime.IsUnknown() {
		model.CreateTime = types.StringNull()
	}
	if model.UpdateTime.IsUnknown() {
		model.UpdateTime = types.StringNull()
	}
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

func expandMetadataVars(ctx context.Context, list types.List) ([]*agsv20250920.MetadataVar, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var models []agsEnvVarModel
	diags := list.ElementsAs(ctx, &models, false)
	if diags.HasError() {
		return nil, diags
	}
	if len(models) == 0 {
		return nil, diags
	}

	metadataVars := make([]*agsv20250920.MetadataVar, 0, len(models))
	for _, model := range models {
		metadataVar := &agsv20250920.MetadataVar{}
		if !model.Name.IsNull() && !model.Name.IsUnknown() {
			metadataVar.Name = model.Name.ValueStringPointer()
		}
		if !model.Value.IsNull() && !model.Value.IsUnknown() {
			metadataVar.Value = model.Value.ValueStringPointer()
		}
		metadataVars = append(metadataVars, metadataVar)
	}

	return metadataVars, diags
}

func flattenMetadataVars(ctx context.Context, items []*agsv20250920.MetadataVar) (types.List, diag.Diagnostics) {
	if len(items) == 0 {
		return types.ListNull(agsEnvVarObjectType), nil
	}

	models := make([]agsEnvVarModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		models = append(models, agsEnvVarModel{
			Name:  stringPointerValue(item.Name),
			Value: stringPointerValue(item.Value),
		})
	}
	if len(models) == 0 {
		return types.ListNull(agsEnvVarObjectType), nil
	}

	return types.ListValueFrom(ctx, agsEnvVarObjectType, models)
}
