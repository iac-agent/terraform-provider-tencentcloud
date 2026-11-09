package teo

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// _ resource.Resource is intentionally omitted: the factory function
// NewTeoFunctionReplicaV8Resource returns resource.Resource, which already
// performs the implicit compile-time check at the return statement. Only the
// extended-interface assertions are kept here.
var (
	_ resource.ResourceWithConfigure   = &TeoFunctionReplicaV8Resource{}
	_ resource.ResourceWithImportState = &TeoFunctionReplicaV8Resource{}
)

// NewTeoFunctionReplicaV8Resource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewTeoFunctionReplicaV8Resource() resource.Resource {
	return &TeoFunctionReplicaV8Resource{}
}

// TeoFunctionReplicaV8Resource implements resource.Resource for
// tencentcloud_teo_function_replica_v8.
type TeoFunctionReplicaV8Resource struct {
	fw.ResourceWithConfigure
}

// TeoFunctionReplicaV8ResourceModel maps the schema attributes.
//
// The API has no resource-level id: the identity is the composite
// `<zone_id>#<function_id>#<replica_name>` (joined with tccommon.FILED_SP),
// which every CRUD call must be able to re-derive.
type TeoFunctionReplicaV8ResourceModel struct {
	Id          types.String `tfsdk:"id"`
	ZoneId      types.String `tfsdk:"zone_id"`
	FunctionId  types.String `tfsdk:"function_id"`
	ReplicaName types.String `tfsdk:"replica_name"`
	Content     types.String `tfsdk:"content"`
	Remark      types.String `tfsdk:"remark"`
}

func (r *TeoFunctionReplicaV8Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_function_replica_v8"
}

func (r *TeoFunctionReplicaV8Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a TEO edge function replica " +
			"(terraform-plugin-framework implementation).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier in the form `<zone_id>#<function_id>#<replica_name>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_id": schema.StringAttribute{
				Required: true,
				Description: "Zone ID. The TEO zone that hosts the edge function. " +
					"Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"function_id": schema.StringAttribute{
				Required:    true,
				Description: "Edge function ID. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"replica_name": schema.StringAttribute{
				Required: true,
				Description: "Edge function replica name. Limited to 1-50 characters; allowed characters are " +
					"a-z, 0-9 and -, where - cannot be used alone, consecutively, or at the beginning or end. " +
					"Replica names must be unique under the same FunctionId. " +
					"The API offers no rename capability, so changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content": schema.StringAttribute{
				Required: true,
				Description: "Edge function replica content. Currently only JavaScript code is supported, " +
					"with a maximum size of 5MB.",
			},
			"remark": schema.StringAttribute{
				Optional: true,
				Description: "Edge function replica description, up to 50 characters. " +
					"Whether the Modify API treats an omitted value as keep or clear is not documented; " +
					"removing this attribute from the configuration is rejected in Update instead of guessing.",
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *TeoFunctionReplicaV8Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeoFunctionReplicaV8ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := teov20220901.NewCreateFunctionReplicaRequest()
	request.ZoneId = plan.ZoneId.ValueStringPointer()
	request.FunctionId = plan.FunctionId.ValueStringPointer()
	request.ReplicaName = plan.ReplicaName.ValueStringPointer()
	request.Content = plan.Content.ValueStringPointer()
	if !plan.Remark.IsNull() && !plan.Remark.IsUnknown() {
		request.Remark = plan.Remark.ValueStringPointer()
	}

	// The identity is composite: the Create API response carries only a
	// RequestId, so it is assembled from the request itself.
	id := strings.Join([]string{
		plan.ZoneId.ValueString(),
		plan.FunctionId.ValueString(),
		plan.ReplicaName.ValueString(),
	}, tccommon.FILED_SP)
	plan.Id = types.StringValue(id)

	service := NewTeoService(r.Client())
	if err := service.CreateTeoFunctionReplica(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error creating TEO function replica",
			fmt.Sprintf("CreateFunctionReplica failed for replica %q under function %q: %v",
				plan.ReplicaName.ValueString(), plan.FunctionId.ValueString(), err),
		)
		return
	}

	// Persist the composite identity BEFORE the read-back so a refresh failure
	// still leaves a destroyable resource in state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := service.DescribeTeoFunctionReplicaByName(ctx,
		plan.ZoneId.ValueString(), plan.FunctionId.ValueString(), plan.ReplicaName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading TEO function replica after create", err.Error())
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"TEO function replica not readable after create",
			fmt.Sprintf("replica %q under function %q was created but cannot be read back",
				plan.ReplicaName.ValueString(), plan.FunctionId.ValueString()),
		)
		return
	}
	refreshTeoFunctionReplicaV8Model(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *TeoFunctionReplicaV8Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TeoFunctionReplicaV8ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, functionId, replicaName, ok := splitTeoFunctionReplicaV8Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO function replica id is broken, id is %q (expected `<zone_id>%s<function_id>%s<replica_name>`)",
				state.Id.ValueString(), tccommon.FILED_SP, tccommon.FILED_SP),
		)
		return
	}

	service := NewTeoService(r.Client())
	item, err := service.DescribeTeoFunctionReplicaByName(ctx, zoneId, functionId, replicaName)
	if err != nil {
		resp.Diagnostics.AddError("Error reading TEO function replica", err.Error())
		return
	}
	// A successful listing without a matching entry is the confirmed
	// not-found signal; API failures are returned as errors above and never
	// reach this branch.
	if item == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ZoneId = types.StringValue(zoneId)
	state.FunctionId = types.StringValue(functionId)
	state.ReplicaName = types.StringValue(replicaName)
	refreshTeoFunctionReplicaV8Model(&state, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications. Only content and remark can change;
// the identity attributes are guarded by RequiresReplace plan modifiers.
func (r *TeoFunctionReplicaV8Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TeoFunctionReplicaV8ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state TeoFunctionReplicaV8ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, functionId, replicaName, ok := splitTeoFunctionReplicaV8Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO function replica id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	request := teov20220901.NewModifyFunctionReplicaRequest()
	request.ZoneId = helper.String(zoneId)
	request.FunctionId = helper.String(functionId)
	request.ReplicaName = helper.String(replicaName)

	// Only send the fields whose value actually changed; the Modify API
	// semantics for an omitted Content/Remark (keep vs clear) are not
	// documented, so nothing is sent unless it changed.
	if !plan.Content.Equal(state.Content) {
		request.Content = plan.Content.ValueStringPointer()
	}

	if !plan.Remark.Equal(state.Remark) {
		switch {
		case plan.Remark.IsNull():
			// Null means the user removed the attribute from the configuration.
			// Whether the API treats an omitted Remark as keep or clear is not
			// documented (ClearSemantics unknown), so fail loudly instead of
			// guessing: clearing would silently drop user data, keeping would
			// produce a permanent diff.
			resp.Diagnostics.AddError(
				"Clear semantics of remark is not confirmed",
				"Removing `remark` from the configuration requires knowing whether "+
					"ModifyFunctionReplica treats an omitted value as \"keep\" or \"clear\". "+
					"Set an explicit value instead, or replace the resource.",
			)
			return
		default:
			request.Remark = plan.Remark.ValueStringPointer()
		}
	}

	if request.Content != nil || request.Remark != nil {
		service := NewTeoService(r.Client())
		if err := service.ModifyTeoFunctionReplica(ctx, request); err != nil {
			resp.Diagnostics.AddError(
				"Error updating TEO function replica",
				fmt.Sprintf("ModifyFunctionReplica failed for replica %q under function %q: %v",
					replicaName, functionId, err),
			)
			return
		}

		item, err := service.DescribeTeoFunctionReplicaByName(ctx, zoneId, functionId, replicaName)
		if err != nil {
			resp.Diagnostics.AddError("Error reading TEO function replica after update", err.Error())
			return
		}
		if item == nil {
			resp.Diagnostics.AddError(
				"TEO function replica not readable after update",
				fmt.Sprintf("replica %q under function %q was updated but cannot be read back",
					replicaName, functionId),
			)
			return
		}
		refreshTeoFunctionReplicaV8Model(&plan, item)
	}

	plan.Id = state.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion.
func (r *TeoFunctionReplicaV8Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeoFunctionReplicaV8ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, functionId, replicaName, ok := splitTeoFunctionReplicaV8Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO function replica id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	service := NewTeoService(r.Client())
	if err := service.DeleteTeoFunctionReplicaByName(ctx, zoneId, functionId, replicaName); err != nil {
		// No documented error code marks "replica already gone" for
		// DeleteFunctionReplica, so any error fails the delete and keeps the
		// id in state for a retry.
		resp.Diagnostics.AddError(
			"Error deleting TEO function replica",
			fmt.Sprintf("DeleteFunctionReplica failed for replica %q under function %q: %v",
				replicaName, functionId, err),
		)
		return
	}
}

// ImportState supports `terraform import` using the
// `<zone_id>#<function_id>#<replica_name>` composite id.
func (r *TeoFunctionReplicaV8Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if _, _, _, ok := splitTeoFunctionReplicaV8Id(req.ID); !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in the form `<zone_id>%s<function_id>%s<replica_name>`, got %q",
				tccommon.FILED_SP, tccommon.FILED_SP, req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// refreshTeoFunctionReplicaV8Model copies the API-returned fields into the
// model. zone_id / function_id / replica_name keep the values derived from
// the composite id (they are the exact request-side strings); content and
// remark are refreshed from the cloud side, mapping an absent pointer to null.
func refreshTeoFunctionReplicaV8Model(m *TeoFunctionReplicaV8ResourceModel, item *teov20220901.FunctionReplica) {
	if item == nil {
		return
	}
	if item.Content != nil {
		m.Content = types.StringPointerValue(item.Content)
	}
	m.Remark = types.StringPointerValue(item.Remark)
}
