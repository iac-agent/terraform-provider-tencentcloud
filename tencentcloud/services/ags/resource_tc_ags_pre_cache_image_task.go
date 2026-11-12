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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// agsPreCacheImageTaskPollInterval is the interval used while waiting for a
// pre-cache image task to reach a stable state.
const agsPreCacheImageTaskPollInterval = 3 * time.Second

var (
	_ resource.ResourceWithConfigure   = &AgsPreCacheImageTaskResource{}
	_ resource.ResourceWithImportState = &AgsPreCacheImageTaskResource{}
)

// NewAgsPreCacheImageTaskResource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewAgsPreCacheImageTaskResource() resource.Resource {
	return &AgsPreCacheImageTaskResource{}
}

// AgsPreCacheImageTaskResource implements resource.Resource for
// tencentcloud_ags_pre_cache_image_task.
//
// CreatePreCacheImageTask has no matching update API, so every configurable
// attribute is marked with RequiresReplace: changing any of them destroys the
// pre-cache image task and creates a new one. There is therefore no Update
// implementation.
type AgsPreCacheImageTaskResource struct {
	fw.ResourceWithConfigure
	fw.WithImportByID
	fw.WithNoUpdate
}

// AgsPreCacheImageTaskResourceModel maps the schema attributes.
type AgsPreCacheImageTaskResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	PreCacheImageId      types.String `tfsdk:"pre_cache_image_id"`
	Image                types.String `tfsdk:"image"`
	ImageRegistryType    types.String `tfsdk:"image_registry_type"`
	ImageDigest          types.String `tfsdk:"image_digest"`
	Status               types.String `tfsdk:"status"`
	Message              types.String `tfsdk:"message"`
	CreateTime           types.String `tfsdk:"create_time"`
	ReadyTime            types.String `tfsdk:"ready_time"`
	LastUsedTime         types.String `tfsdk:"last_used_time"`
	SourceType           types.String `tfsdk:"source_type"`
	CachedImageSizeBytes types.Int64  `tfsdk:"cached_image_size_bytes"`
}

func (r *AgsPreCacheImageTaskResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_ags_pre_cache_image_task"
}

func (r *AgsPreCacheImageTaskResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create an AGS pre-cache image task.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Id of the resource, which is the same as `pre_cache_image_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pre_cache_image_id": schema.StringAttribute{
				Computed:    true,
				Description: "Id of the pre-cache image task.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"image": schema.StringAttribute{
				Required: true,
				Description: "Image address. Only `repository:tag`, `repository@sha256:<64 hex digits>` and " +
					"`repository:tag@sha256:<64 hex digits>` are supported. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"image_registry_type": schema.StringAttribute{
				Required: true,
				Description: "Image registry type. Valid values: `enterprise`, `personal`, `custom`. " +
					"`enterprise` is the TCR enterprise image registry and `personal` is the CCR personal image " +
					"registry. It cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"image_digest": schema.StringAttribute{
				Computed:    true,
				Description: "Digest of the image, resolved and reported by the cloud side.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Status of the pre-cache image task.",
			},
			"message": schema.StringAttribute{
				Computed:    true,
				Description: "Description of the pre-cache image task status.",
			},
			"create_time": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the pre-cache image task, reported by the cloud side as a string.",
			},
			"ready_time": schema.StringAttribute{
				Computed:    true,
				Description: "Time when the pre-cache image task completed, reported by the cloud side as a string.",
			},
			"last_used_time": schema.StringAttribute{
				Computed: true,
				Description: "Time when the pre-cached image was last used by a sandbox instance. It changes while " +
					"the resource is managed, so it is only reported by the cloud side.",
			},
			"source_type": schema.StringAttribute{
				Computed: true,
				Description: "Source type of the pre-cache image resource. Valid values: `EXPLICIT` (created " +
					"manually), `AUTO` (created automatically), `TCR_AUTO` (pre-cached automatically by TCR).",
			},
			"cached_image_size_bytes": schema.Int64Attribute{
				Computed:    true,
				Description: "Storage size of the pre-cached image, in bytes.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *AgsPreCacheImageTaskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AgsPreCacheImageTaskResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := agsv20250920.NewCreatePreCacheImageTaskRequest()
	if !plan.Image.IsNull() && !plan.Image.IsUnknown() {
		request.Image = plan.Image.ValueStringPointer()
	}
	if !plan.ImageRegistryType.IsNull() && !plan.ImageRegistryType.IsUnknown() {
		request.ImageRegistryType = plan.ImageRegistryType.ValueStringPointer()
	}

	service := NewAgsService(r.Client())
	preCacheImageId, err := service.CreatePreCacheImageTask(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating AGS pre-cache image task",
			fmt.Sprintf("CreatePreCacheImageTask failed: %v", err),
		)
		return
	}

	plan.Id = types.StringValue(preCacheImageId)
	plan.PreCacheImageId = types.StringValue(preCacheImageId)
	// The create response may already report some fields, but the authoritative
	// values come from the describe API below. The identity is written first so
	// a failed wait still leaves the resource in state and `terraform destroy`
	// can delete it.
	normalizePreCacheImageTaskUnknowns(&plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.waitPreCacheImageTaskReady(ctx, service, preCacheImageId, tccommon.WriteRetryTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for AGS pre-cache image task to become ready",
			fmt.Sprintf("pre-cache image task %s was created but is not ready yet: %v", preCacheImageId, err),
		)
		return
	}

	refreshPreCacheImageTaskModel(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and the post-apply state refresh.
func (r *AgsPreCacheImageTaskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AgsPreCacheImageTaskResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	preCacheImageId, diags := preCacheImageTaskIdOf(&state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewAgsService(r.Client())
	item, err := service.DescribePreCacheImageTaskById(ctx, preCacheImageId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading AGS pre-cache image task",
			fmt.Sprintf("DescribePreCacheImageTask failed for %s: %v", preCacheImageId, err),
		)
		return
	}

	// A nil result means the query succeeded and the pre-cache image task no
	// longer exists. A nil API response is handled as an error inside the
	// service and never reaches this branch, so state is not dropped on a
	// failed request.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"AGS pre-cache image task not found",
			fmt.Sprintf("pre-cache image task %s does not exist any more, removing it from state.", preCacheImageId),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	refreshPreCacheImageTaskModel(&state, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete handles `terraform destroy`. It only removes the pre-cache image entry
// itself: it does not affect running sandbox instances and it does not touch
// the source image registry.
func (r *AgsPreCacheImageTaskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AgsPreCacheImageTaskResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	preCacheImageId, diags := preCacheImageTaskIdOf(&state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	service := NewAgsService(r.Client())

	// The item may already be gone, for example when it was removed outside of
	// Terraform. Existence is checked first so a repeated destroy is
	// idempotent instead of reporting an error.
	item, describeErr := service.DescribePreCacheImageTaskById(ctx, preCacheImageId)
	if describeErr == nil && item == nil {
		return
	}

	request := agsv20250920.NewDeletePreCacheImageRequest()
	request.PreCacheImageId = &preCacheImageId
	if item != nil {
		if item.Image != "" {
			request.Image = &item.Image
		}
		if item.ImageDigest != "" {
			request.ImageDigest = &item.ImageDigest
		}
		if item.ImageRegistryType != "" {
			request.ImageRegistryType = &item.ImageRegistryType
		}
	}

	if err := service.DeletePreCacheImage(ctx, request); err != nil {
		// The API documentation does not define an "already deleted" error
		// code, so existence is checked again instead of guessing codes: an
		// item that is really gone makes the delete idempotent, everything else
		// is an error.
		remaining, checkErr := service.DescribePreCacheImageTaskById(ctx, preCacheImageId)
		if checkErr != nil || remaining != nil {
			resp.Diagnostics.AddError(
				"Error deleting AGS pre-cache image task",
				fmt.Sprintf("DeletePreCacheImage failed for %s: %v", preCacheImageId, err),
			)
		}
		return
	}

	// The entry may stay visible for a while. Wait until it is gone so a later
	// apply does not hit a "still exists" error.
	if err := r.waitPreCacheImageTaskGone(ctx, service, preCacheImageId, tccommon.WriteRetryTimeout); err != nil {
		resp.Diagnostics.AddWarning(
			"AGS pre-cache image task is still being deleted",
			fmt.Sprintf("pre-cache image task %s was deleted but is still visible: %v", preCacheImageId, err),
		)
	}
}

// ImportState supports `terraform import` with the pre-cache image id.
func (r *AgsPreCacheImageTaskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError(
			"Invalid import id",
			"The import id must be the id of an existing AGS pre-cache image task.",
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("pre_cache_image_id"), req.ID)...)
}

// refreshPreCacheImageTaskModel copies the cloud state into the model. The
// model is fully rebuilt from the describe response so a single State.Set is
// enough after every lifecycle callback.
func refreshPreCacheImageTaskModel(model *AgsPreCacheImageTaskResourceModel, item *AgsPreCacheImageTask) {
	if item == nil {
		return
	}

	model.PreCacheImageId = types.StringValue(item.PreCacheImageId)
	model.Id = types.StringValue(item.PreCacheImageId)
	model.Image = stringValueOrNull(item.Image)
	model.ImageRegistryType = stringValueOrNull(item.ImageRegistryType)
	model.ImageDigest = stringValueOrNull(item.ImageDigest)
	model.Status = stringValueOrNull(item.Status)
	model.Message = stringValueOrNull(item.Message)
	model.CreateTime = stringValueOrNull(item.CreateTime)
	model.ReadyTime = stringValueOrNull(item.ReadyTime)
	model.LastUsedTime = stringValueOrNull(item.LastUsedTime)
	model.SourceType = stringValueOrNull(item.SourceType)
	model.CachedImageSizeBytes = types.Int64Value(item.CachedImageSizeBytes)
}

// normalizePreCacheImageTaskUnknowns replaces the computed attributes that are
// still unknown at create time with null, so the identity can be stored in
// state before the cloud has reported those values.
func normalizePreCacheImageTaskUnknowns(model *AgsPreCacheImageTaskResourceModel) {
	if model.ImageDigest.IsUnknown() {
		model.ImageDigest = types.StringNull()
	}
	if model.Status.IsUnknown() {
		model.Status = types.StringNull()
	}
	if model.Message.IsUnknown() {
		model.Message = types.StringNull()
	}
	if model.CreateTime.IsUnknown() {
		model.CreateTime = types.StringNull()
	}
	if model.ReadyTime.IsUnknown() {
		model.ReadyTime = types.StringNull()
	}
	if model.LastUsedTime.IsUnknown() {
		model.LastUsedTime = types.StringNull()
	}
	if model.SourceType.IsUnknown() {
		model.SourceType = types.StringNull()
	}
	if model.CachedImageSizeBytes.IsUnknown() {
		model.CachedImageSizeBytes = types.Int64Null()
	}
}

// preCacheImageTaskIdOf returns the identity stored in the model.
func preCacheImageTaskIdOf(model *AgsPreCacheImageTaskResourceModel) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !model.PreCacheImageId.IsNull() && !model.PreCacheImageId.IsUnknown() && model.PreCacheImageId.ValueString() != "" {
		return model.PreCacheImageId.ValueString(), diags
	}
	if !model.Id.IsNull() && !model.Id.IsUnknown() && model.Id.ValueString() != "" {
		return model.Id.ValueString(), diags
	}

	diags.AddError(
		"Invalid AGS pre-cache image task state",
		"The resource state does not contain a pre-cache image id.",
	)
	return "", diags
}

// waitPreCacheImageTaskReady waits until the pre-cache image task reports a
// ready time.
//
// DescribePreCacheImageTask does not enumerate its status values and does not
// document a failure status, so the terminal condition is the documented
// ReadyTime field. A message reported without a ready time is not treated as a
// failure: the cloud can report progress messages for a task that is still
// running.
func (r *AgsPreCacheImageTaskResource) waitPreCacheImageTaskReady(ctx context.Context, service AgsService, preCacheImageId string, timeout time.Duration) (*AgsPreCacheImageTask, error) {
	deadline := time.Now().Add(timeout)
	for {
		item, err := service.DescribePreCacheImageTaskById(ctx, preCacheImageId)
		if err != nil {
			return nil, err
		}
		if item == nil {
			return nil, fmt.Errorf("pre-cache image task %s disappeared while waiting for it to become ready", preCacheImageId)
		}

		if item.ReadyTime != "" {
			return item, nil
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timeout waiting for pre-cache image task %s to become ready, current status is %q, message is %q",
				preCacheImageId, item.Status, item.Message)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(agsPreCacheImageTaskPollInterval):
		}
	}
}

// waitPreCacheImageTaskGone waits until the deleted pre-cache image task is no
// longer visible.
func (r *AgsPreCacheImageTaskResource) waitPreCacheImageTaskGone(ctx context.Context, service AgsService, preCacheImageId string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		item, err := service.DescribePreCacheImageTaskById(ctx, preCacheImageId)
		if err != nil {
			return err
		}
		if item == nil {
			return nil
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("pre-cache image task %s is still present with status %q after %s",
				preCacheImageId, item.Status, timeout)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(agsPreCacheImageTaskPollInterval):
		}
	}
}

// stringValueOrNull converts a plain string into a typed value, mapping the
// empty string to null so absent cloud fields do not look like set values.
func stringValueOrNull(value string) types.String {
	if value == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}
