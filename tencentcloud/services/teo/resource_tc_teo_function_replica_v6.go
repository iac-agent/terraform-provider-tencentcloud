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
)

// _ resource.Resource is intentionally omitted: the factory function
// NewTeoFunctionReplicaV6Resource returns resource.Resource, which already
// performs the implicit compile-time check at the return statement. Only the
// extended-interface assertions are kept here.
var (
	_ resource.ResourceWithConfigure   = &TeoFunctionReplicaV6Resource{}
	_ resource.ResourceWithImportState = &TeoFunctionReplicaV6Resource{}
)

// NewTeoFunctionReplicaV6Resource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewTeoFunctionReplicaV6Resource() resource.Resource {
	return &TeoFunctionReplicaV6Resource{}
}

// TeoFunctionReplicaV6Resource implements resource.Resource for
// tencentcloud_teo_function_replica_v6.
type TeoFunctionReplicaV6Resource struct {
	fw.ResourceWithConfigure
}

// TeoFunctionReplicaV6ResourceModel maps the schema attributes.
type TeoFunctionReplicaV6ResourceModel struct {
	Id          types.String `tfsdk:"id"`
	ZoneId      types.String `tfsdk:"zone_id"`
	FunctionId  types.String `tfsdk:"function_id"`
	ReplicaName types.String `tfsdk:"replica_name"`
	Content     types.String `tfsdk:"content"`
	Remark      types.String `tfsdk:"remark"`
}

func (r *TeoFunctionReplicaV6Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_function_replica_v6"
}

func (r *TeoFunctionReplicaV6Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a TEO edge function replica. " +
			"CreateFunctionReplica does not return a resource-level id, so the resource id is the composite key " +
			"`<zone_id>#<function_id>#<replica_name>`. Changing zone_id, function_id or replica_name forces recreation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier in the form `<zone_id>#<function_id>#<replica_name>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_id": schema.StringAttribute{
				Required:    true,
				Description: "Site ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"function_id": schema.StringAttribute{
				Required:    true,
				Description: "Function ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"replica_name": schema.StringAttribute{
				Required: true,
				Description: "Edge function replica name. Length should be 1-50 characters, allowed characters are a-z, 0-9 and -, " +
					"and - cannot be used alone or consecutively, nor placed at the beginning or the end. " +
					"Replica names must be unique under the same FunctionId.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content": schema.StringAttribute{
				Required:    true,
				Description: "Edge function replica content. Currently only JavaScript code is supported, with a maximum size of 5MB.",
			},
			"remark": schema.StringAttribute{
				Optional:    true,
				Description: "Edge function replica description. Maximum 50 characters.",
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *TeoFunctionReplicaV6Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeoFunctionReplicaV6ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	zoneId := plan.ZoneId.ValueString()
	functionId := plan.FunctionId.ValueString()
	replicaName := plan.ReplicaName.ValueString()

	var remark *string
	if !plan.Remark.IsNull() && !plan.Remark.IsUnknown() {
		remark = plan.Remark.ValueStringPointer()
	}

	service := NewTeoService(r.Client())
	if err := service.CreateTeoFunctionReplica(ctx, zoneId, functionId, replicaName, plan.Content.ValueString(), remark); err != nil {
		resp.Diagnostics.AddError(
			"Error creating TEO function replica",
			fmt.Sprintf("CreateFunctionReplica failed for function %s replica %s: %v", functionId, replicaName, err),
		)
		return
	}

	plan.Id = types.StringValue(strings.Join([]string{zoneId, functionId, replicaName}, tccommon.FILED_SP))

	// Persist the identity before the read-back so a partial failure keeps the
	// id available for reconciliation.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := service.DescribeTeoFunctionReplicaById(ctx, zoneId, functionId, replicaName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO function replica after create",
			fmt.Sprintf("DescribeFunctionReplicas failed for function %s replica %s: %v", functionId, replicaName, err),
		)
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"TEO function replica not readable after create",
			fmt.Sprintf("function %s replica %s was created but cannot be read back", functionId, replicaName),
		)
		return
	}
	r.refreshModel(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *TeoFunctionReplicaV6Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TeoFunctionReplicaV6ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	zoneId, functionId, replicaName, ok := splitTeoFunctionReplicaV6Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO function replica id is broken, id is %q (expected `<zone_id>%s<function_id>%s<replica_name>`)",
				state.Id.ValueString(), tccommon.FILED_SP, tccommon.FILED_SP),
		)
		return
	}

	service := NewTeoService(r.Client())
	item, err := service.DescribeTeoFunctionReplicaById(ctx, zoneId, functionId, replicaName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO function replica",
			fmt.Sprintf("DescribeFunctionReplicas failed for function %s replica %s: %v", functionId, replicaName, err),
		)
		return
	}

	// Queried successfully but nothing matched => the replica is confirmed gone.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("TEO function replica %s not found, removing from state.", state.Id.ValueString()),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	// zone_id / function_id / replica_name are the identity segments; they are
	// refreshed from the id. ZoneId is not echoed back by
	// DescribeFunctionReplicas, so the state keeps the configured value.
	state.ZoneId = types.StringValue(zoneId)
	state.FunctionId = types.StringValue(functionId)
	state.ReplicaName = types.StringValue(replicaName)
	r.refreshModel(&state, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place modifications.
func (r *TeoFunctionReplicaV6Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TeoFunctionReplicaV6ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state TeoFunctionReplicaV6ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	zoneId, functionId, replicaName, ok := splitTeoFunctionReplicaV6Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO function replica id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	// content is a required attribute: an update plan always carries a known
	// value, and it is sent unconditionally.
	var content *string
	if !plan.Content.IsNull() && !plan.Content.IsUnknown() {
		content = plan.Content.ValueStringPointer()
	}

	// remark is optional. Three states must stay distinguishable:
	//   Unknown -> value not decided yet, nothing to send
	//   Null    -> user removed the attribute from the configuration; whether
	//              the API treats an omitted value as "keep" or "clear" is NOT
	//              confirmed, so the ambiguous transition fails loudly instead
	//              of guessing a direction
	//   value   -> send it
	var remark *string
	switch {
	case plan.Remark.IsUnknown():
		// nothing to send
	case plan.Remark.IsNull():
		if !state.Remark.IsNull() {
			resp.Diagnostics.AddError(
				"Clear semantics of remark is not confirmed",
				"Removing the remark attribute from the configuration requires knowing whether "+
					"ModifyFunctionReplica treats an omitted value as \"keep\" or \"clear\". "+
					"This semantics is not documented yet; refusing to guess. "+
					"Keep the remark attribute set, or destroy and recreate the replica to change it.",
			)
			return
		}
		// remark was already absent on both sides: nothing to send.
	default:
		remark = plan.Remark.ValueStringPointer()
	}

	service := NewTeoService(r.Client())
	if err := service.ModifyTeoFunctionReplica(ctx, zoneId, functionId, replicaName, content, remark); err != nil {
		resp.Diagnostics.AddError(
			"Error updating TEO function replica",
			fmt.Sprintf("ModifyFunctionReplica failed for function %s replica %s: %v", functionId, replicaName, err),
		)
		return
	}

	item, err := service.DescribeTeoFunctionReplicaById(ctx, zoneId, functionId, replicaName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO function replica after update",
			fmt.Sprintf("DescribeFunctionReplicas failed for function %s replica %s: %v", functionId, replicaName, err),
		)
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"TEO function replica not readable after update",
			fmt.Sprintf("function %s replica %s was updated but cannot be read back", functionId, replicaName),
		)
		return
	}

	plan.Id = state.Id
	r.refreshModel(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles `terraform destroy`.
func (r *TeoFunctionReplicaV6Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeoFunctionReplicaV6ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	zoneId, functionId, replicaName, ok := splitTeoFunctionReplicaV6Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO function replica id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	service := NewTeoService(r.Client())

	// Pre-read: DeleteFunctionReplica's behaviour for a non-existent replica is
	// not documented, so existence is confirmed via DescribeFunctionReplicas
	// first. A confirmed "already gone" makes destroy idempotent.
	item, err := service.DescribeTeoFunctionReplicaById(ctx, zoneId, functionId, replicaName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO function replica before delete",
			fmt.Sprintf("DescribeFunctionReplicas failed for function %s replica %s: %v", functionId, replicaName, err),
		)
		return
	}
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("TEO function replica %s not found, it may have been already deleted.", state.Id.ValueString()),
		)
		return
	}

	if err := service.DeleteTeoFunctionReplica(ctx, zoneId, functionId, replicaName); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting TEO function replica",
			fmt.Sprintf("DeleteFunctionReplica failed for function %s replica %s: %v", functionId, replicaName, err),
		)
		return
	}
}

// ImportState supports `terraform import` using the
// `<zone_id>#<function_id>#<replica_name>` form.
func (r *TeoFunctionReplicaV6Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	zoneId, functionId, replicaName, ok := splitTeoFunctionReplicaV6Id(req.ID)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in the form `<zone_id>%s<function_id>%s<replica_name>`, got %q",
				tccommon.FILED_SP, tccommon.FILED_SP, req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("zone_id"), zoneId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("function_id"), functionId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("replica_name"), replicaName)...)
}

// refreshModel copies cloud-owned fields into the model.
//
// DescribeFunctionReplicas returns FunctionId, ReplicaName, Content, Remark,
// CreatedOn and ModifiedOn. Only the configurable fields are refreshed;
// zone_id is not part of the response and is preserved from the id.
func (r *TeoFunctionReplicaV6Resource) refreshModel(m *TeoFunctionReplicaV6ResourceModel, item *teov20220901.FunctionReplica) {
	if item == nil {
		return
	}
	if item.Content != nil {
		m.Content = types.StringPointerValue(item.Content)
	}
	if item.Remark != nil {
		m.Remark = types.StringPointerValue(item.Remark)
	} else {
		m.Remark = types.StringNull()
	}
}

// splitTeoFunctionReplicaV6Id parses `<zone_id>#<function_id>#<replica_name>`.
func splitTeoFunctionReplicaV6Id(id string) (string, string, string, bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}
