package ags

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var (
	_ resource.ResourceWithConfigure   = &AgsGPUResourceGroupResource{}
	_ resource.ResourceWithImportState = &AgsGPUResourceGroupResource{}
)

// NewAgsGPUResourceGroupResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewAgsGPUResourceGroupResource() resource.Resource {
	return &AgsGPUResourceGroupResource{}
}

// AgsGPUResourceGroupResource implements resource.Resource for
// tencentcloud_ags_gpu_resource_group. The AGS product does not expose a
// public creation API for GPU resource groups (they are purchased through
// other channels), so this resource only supports import and name updates.
type AgsGPUResourceGroupResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
}

// AgsGPUResourceGroupResourceModel maps the schema attributes.
type AgsGPUResourceGroupResourceModel struct {
	Id              types.String `tfsdk:"id"`
	ResourceGroupId types.String `tfsdk:"resource_group_id"`
	Name            types.String `tfsdk:"name"`
	GPUType         types.String `tfsdk:"gpu_type"`
	GPUTypeAlias    types.String `tfsdk:"gpu_type_alias"`
	GPUMemoryGiB    types.Int64  `tfsdk:"gpu_memory_gib"`
	CPU             types.Int64  `tfsdk:"cpu"`
	Memory          types.Int64  `tfsdk:"memory"`
	GPUNum          types.Int64  `tfsdk:"gpu_num"`
	TotalCount      types.Int64  `tfsdk:"total_count"`
	AvailableCount  types.Int64  `tfsdk:"available_count"`
	AllocatedCount  types.Int64  `tfsdk:"allocated_count"`
}

func (r *AgsGPUResourceGroupResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_ags_gpu_resource_group"
}

func (r *AgsGPUResourceGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an AGS GPU resource group resource. GPU resource groups cannot be created or released through a public API: purchase them in the AGS console first, then import them with `terraform import tencentcloud_ags_gpu_resource_group.example <resource_group_id>`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Id of the resource, which is the same as `resource_group_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_group_id": schema.StringAttribute{
				Computed:    true,
				Description: "Id of the GPU resource group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the GPU resource group. After trimming whitespace it must not be empty and can be at most 128 bytes long.",
			},
			"gpu_type": schema.StringAttribute{
				Computed:    true,
				Description: "GPU type identifier of the fixed spec, for example `5090d`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gpu_type_alias": schema.StringAttribute{
				Computed:    true,
				Description: "Display name of the GPU type. When it is empty the cloud side displays `gpu_type`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gpu_memory_gib": schema.Int64Attribute{
				Computed:    true,
				Description: "Total memory of all GPUs on one instance, in GiB. `0` means the display metadata is not maintained. It does not need to be multiplied by `gpu_num`.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"cpu": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of CPU cores on one instance.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"memory": schema.Int64Attribute{
				Computed:    true,
				Description: "Host memory of one instance, in GiB. It is not the GPU memory.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"gpu_num": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of GPU cards on one instance.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"total_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Total number of resources in the group. Destroyed resources and failed resources whose compensation has completed are excluded.",
			},
			"available_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of resources that are Active, unoccupied and physically claimable.",
			},
			"allocated_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of resources that are still occupied or bound, including those being claimed, in use or whose release has not finished.",
			},
		},
	}
}

// Create is not supported: the AGS product does not expose a public API to
// create GPU resource groups, they are purchased in the console.
func (r *AgsGPUResourceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError(
		"AGS GPU resource group cannot be created",
		"The AGS product does not provide a public API to create GPU resource groups. Purchase the resource group in the AGS console first, then import it: terraform import tencentcloud_ags_gpu_resource_group.example <resource_group_id>",
	)
}

// Read handles `terraform refresh` and the post-apply state refresh.
func (r *AgsGPUResourceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AgsGPUResourceGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceGroupId, diags := agsGpuResourceGroupIdOf(&state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewAgsService(r.Client())
	item, err := service.DescribeGPUResourceGroupById(ctx, resourceGroupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading AGS GPU resource group",
			fmt.Sprintf("DescribeGPUResourceGroups failed for %s: %v", resourceGroupId, err),
		)
		return
	}

	// A nil result means the query succeeded and the group no longer exists. A
	// nil API response is handled as an error inside the service and never
	// reaches this branch, so state is not dropped on a failed request.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"AGS GPU resource group not found",
			fmt.Sprintf("GPU resource group %s does not exist any more, removing it from state.", resourceGroupId),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(refreshGPUResourceGroupModel(&state, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update renames the GPU resource group with ModifyGPUResourceGroup.
func (r *AgsGPUResourceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AgsGPUResourceGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state AgsGPUResourceGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceGroupId, diags := agsGpuResourceGroupIdOf(&state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := agsv20250920.NewModifyGPUResourceGroupRequest()
	request.ResourceGroupId = helper.String(resourceGroupId)
	request.Name = plan.Name.ValueStringPointer()

	service := NewAgsService(r.Client())
	if err := service.ModifyGPUResourceGroup(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating AGS GPU resource group",
			fmt.Sprintf("ModifyGPUResourceGroup failed for %s: %v", resourceGroupId, err),
		)
		return
	}

	item, err := service.DescribeGPUResourceGroupById(ctx, resourceGroupId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading AGS GPU resource group after update",
			fmt.Sprintf("DescribeGPUResourceGroups failed for %s: %v", resourceGroupId, err),
		)
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"AGS GPU resource group not readable after update",
			fmt.Sprintf("GPU resource group %s was updated but cannot be read back", resourceGroupId),
		)
		return
	}

	plan.Id = state.Id
	plan.ResourceGroupId = state.ResourceGroupId
	resp.Diagnostics.Append(refreshGPUResourceGroupModel(&plan, item)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete is not supported: the AGS product does not expose a public API to
// release GPU resource groups.
func (r *AgsGPUResourceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError(
		"AGS GPU resource group cannot be deleted",
		"The AGS product does not provide a public API to release GPU resource groups. Release the resource group in the AGS console instead.",
	)
}

// ImportState supports `terraform import` with the GPU resource group id.
func (r *AgsGPUResourceGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError(
			"Invalid import id",
			"The import id must be the id of an existing GPU resource group.",
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("resource_group_id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), "")...)
}

func uint64PointerValue(value *uint64) types.Int64 {
	if value == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*value))
}

func refreshGPUResourceGroupModel(model *AgsGPUResourceGroupResourceModel, item *agsv20250920.GPUResourceGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	if item == nil {
		return diags
	}

	if item.ResourceGroupId != nil && *item.ResourceGroupId != "" {
		model.ResourceGroupId = types.StringValue(*item.ResourceGroupId)
		model.Id = types.StringValue(*item.ResourceGroupId)
	}
	model.Name = stringPointerValue(item.Name)

	if item.Spec != nil {
		model.GPUType = stringPointerValue(item.Spec.GPUType)
		model.GPUTypeAlias = stringPointerValue(item.Spec.GPUTypeAlias)
		model.GPUMemoryGiB = uint64PointerValue(item.Spec.GPUMemoryGiB)
		model.CPU = uint64PointerValue(item.Spec.CPU)
		model.Memory = uint64PointerValue(item.Spec.Memory)
		model.GPUNum = uint64PointerValue(item.Spec.GPUNum)
	}

	if item.Statistics != nil {
		model.TotalCount = int64PointerValue(item.Statistics.TotalCount)
		model.AvailableCount = int64PointerValue(item.Statistics.AvailableCount)
		model.AllocatedCount = int64PointerValue(item.Statistics.AllocatedCount)
	}

	return diags
}

func agsGpuResourceGroupIdOf(model *AgsGPUResourceGroupResourceModel) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !model.ResourceGroupId.IsNull() && !model.ResourceGroupId.IsUnknown() && model.ResourceGroupId.ValueString() != "" {
		return model.ResourceGroupId.ValueString(), diags
	}
	if !model.Id.IsNull() && !model.Id.IsUnknown() && model.Id.ValueString() != "" {
		return model.Id.ValueString(), diags
	}

	diags.AddError(
		"Invalid AGS GPU resource group state",
		"The resource state does not contain a GPU resource group id.",
	)
	return "", diags
}
