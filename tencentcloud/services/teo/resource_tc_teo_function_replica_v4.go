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
// NewTeoFunctionReplicaV4Resource returns resource.Resource, which already
// performs the implicit compile-time check at the return statement. Only the
// extended-interface assertions are kept here.
var (
	_ resource.ResourceWithConfigure   = &TeoFunctionReplicaV4Resource{}
	_ resource.ResourceWithImportState = &TeoFunctionReplicaV4Resource{}
)

// NewTeoFunctionReplicaV4Resource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewTeoFunctionReplicaV4Resource() resource.Resource {
	return &TeoFunctionReplicaV4Resource{}
}

// TeoFunctionReplicaV4Resource implements resource.Resource for
// tencentcloud_teo_function_replica_v4.
type TeoFunctionReplicaV4Resource struct {
	fw.ResourceWithConfigure
}

// TeoFunctionReplicaV4ResourceModel maps the schema attributes.
//
// ReplicaName keeps the user supplied value: the list API documents its
// replica-name filter as a fuzzy match, so the value written into state is the
// exact name this resource manages and never a value echoed back by the cloud.
type TeoFunctionReplicaV4ResourceModel struct {
	Id          types.String `tfsdk:"id"`
	ZoneId      types.String `tfsdk:"zone_id"`
	FunctionId  types.String `tfsdk:"function_id"`
	ReplicaName types.String `tfsdk:"replica_name"`
	Content     types.String `tfsdk:"content"`
	Remark      types.String `tfsdk:"remark"`
}

func (r *TeoFunctionReplicaV4Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_function_replica_v4"
}

func (r *TeoFunctionReplicaV4Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to manage a TEO edge function replica. " +
			"A replica carries an independent copy of the edge function code; deleting this resource only removes " +
			"the replica it manages and leaves the function itself and the other replicas under the same function untouched.",
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
				Description: "Zone ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"function_id": schema.StringAttribute{
				Required:    true,
				Description: "Function ID. The replica is created under this function and cannot be moved to another function in place.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"replica_name": schema.StringAttribute{
				Required: true,
				Description: "Edge function replica name. It must be 1-50 characters long and may contain only lowercase " +
					"letters, digits and `-`; `-` cannot be used alone or consecutively, and cannot be placed at the " +
					"beginning or end. The name must be unique under the same function. It identifies the replica, " +
					"so changing it creates a new replica and removes the old one.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content": schema.StringAttribute{
				Required: true,
				Description: "Edge function replica content. Currently only JavaScript code is supported, up to 5 MB. " +
					"Changes are applied in place.",
			},
			"remark": schema.StringAttribute{
				Optional: true,
				Description: "Edge function replica description, up to 50 characters. Removing this attribute leaves " +
					"the description on the cloud side unchanged.",
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *TeoFunctionReplicaV4Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeoFunctionReplicaV4ResourceModel
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
	if zoneId == "" || functionId == "" || replicaName == "" {
		resp.Diagnostics.AddError(
			"Missing required argument",
			"zone_id, function_id and replica_name are required and cannot be empty.",
		)
		return
	}

	request := teov20220901.NewCreateFunctionReplicaRequest()
	request.ZoneId = helper.String(zoneId)
	request.FunctionId = helper.String(functionId)
	request.ReplicaName = helper.String(replicaName)
	request.Content = helper.String(plan.Content.ValueString())
	if !plan.Remark.IsNull() && !plan.Remark.IsUnknown() {
		request.Remark = helper.String(plan.Remark.ValueString())
	}

	if _, err := r.Client().UseTeoClient().CreateFunctionReplicaWithContext(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error creating TEO function replica",
			fmt.Sprintf("CreateFunctionReplica failed: %v", err),
		)
		return
	}

	// CreateFunctionReplica only returns a request id, so the identity is built
	// from the three attributes that locate the replica. It is persisted before
	// the read back so that a failed read still leaves a destroyable resource.
	plan.Id = types.StringValue(buildTeoFunctionReplicaV4Id(zoneId, functionId, replicaName))
	plan.ReplicaName = types.StringValue(replicaName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	replica, err := r.read(ctx, zoneId, functionId, replicaName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO function replica after create",
			fmt.Sprintf("DescribeFunctionReplicas failed: %v", err),
		)
		return
	}
	if replica == nil {
		resp.Diagnostics.AddError(
			"TEO function replica not readable after create",
			fmt.Sprintf("replica %s was created but cannot be read back", plan.Id.ValueString()),
		)
		return
	}

	r.flushModel(&plan, zoneId, functionId, replicaName, replica)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *TeoFunctionReplicaV4Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TeoFunctionReplicaV4ResourceModel
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

	id := state.Id.ValueString()
	zoneId, functionId, replicaName, ok := splitTeoFunctionReplicaV4Id(id)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource id",
			fmt.Sprintf("expected the id to be `<zone_id>%s<function_id>%s<replica_name>`, got %q",
				tccommon.FILED_SP, tccommon.FILED_SP, id),
		)
		return
	}

	replica, err := r.read(ctx, zoneId, functionId, replicaName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO function replica",
			fmt.Sprintf("DescribeFunctionReplicas failed: %v", err),
		)
		return
	}

	// A successful query that returns no exact match means the replica no longer
	// exists; a nil response is an API failure and is reported as an error above.
	if replica == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("TEO function replica %s not found, removing from state.", id),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	state.Id = types.StringValue(id)
	state.ZoneId = types.StringValue(zoneId)
	state.FunctionId = types.StringValue(functionId)
	r.flushModel(&state, zoneId, functionId, replicaName, replica)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates.
func (r *TeoFunctionReplicaV4Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TeoFunctionReplicaV4ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state TeoFunctionReplicaV4ResourceModel
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

	zoneId, functionId, replicaName, ok := splitTeoFunctionReplicaV4Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource id",
			fmt.Sprintf("expected the id to be `<zone_id>%s<function_id>%s<replica_name>`, got %q",
				tccommon.FILED_SP, tccommon.FILED_SP, state.Id.ValueString()),
		)
		return
	}

	request := teov20220901.NewModifyFunctionReplicaRequest()
	request.ZoneId = helper.String(zoneId)
	request.FunctionId = helper.String(functionId)
	request.ReplicaName = helper.String(replicaName)
	if !plan.Content.IsNull() && !plan.Content.IsUnknown() {
		request.Content = helper.String(plan.Content.ValueString())
	}
	// An omitted remark is not sent: whether the API treats an omitted value as
	// "keep" or "clear" is not documented, and sending an empty string could
	// silently wipe the description. An explicit empty string is passed through
	// so that clearing remains possible.
	if !plan.Remark.IsNull() && !plan.Remark.IsUnknown() {
		request.Remark = helper.String(plan.Remark.ValueString())
	}

	if _, err := r.Client().UseTeoClient().ModifyFunctionReplicaWithContext(ctx, request); err != nil {
		resp.Diagnostics.AddError(
			"Error updating TEO function replica",
			fmt.Sprintf("ModifyFunctionReplica failed: %v", err),
		)
		return
	}

	replica, err := r.read(ctx, zoneId, functionId, replicaName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO function replica after update",
			fmt.Sprintf("DescribeFunctionReplicas failed: %v", err),
		)
		return
	}
	if replica == nil {
		resp.Diagnostics.AddError(
			"TEO function replica not readable after update",
			fmt.Sprintf("replica %s was updated but cannot be read back", state.Id.ValueString()),
		)
		return
	}

	plan.Id = state.Id
	plan.ZoneId = types.StringValue(zoneId)
	plan.FunctionId = types.StringValue(functionId)
	r.flushModel(&plan, zoneId, functionId, replicaName, replica)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles resource deletion. DeleteFunctionReplica accepts a list of
// replica names, but this resource only ever deletes the single replica it
// manages.
func (r *TeoFunctionReplicaV4Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeoFunctionReplicaV4ResourceModel
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

	zoneId, functionId, replicaName, ok := splitTeoFunctionReplicaV4Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource id",
			fmt.Sprintf("expected the id to be `<zone_id>%s<function_id>%s<replica_name>`, got %q",
				tccommon.FILED_SP, tccommon.FILED_SP, state.Id.ValueString()),
		)
		return
	}

	request := teov20220901.NewDeleteFunctionReplicaRequest()
	request.ZoneId = helper.String(zoneId)
	request.FunctionId = helper.String(functionId)
	request.ReplicaNames = []*string{helper.String(replicaName)}

	if _, err := r.Client().UseTeoClient().DeleteFunctionReplicaWithContext(ctx, request); err != nil {
		// The API does not document which error code (if any) is returned when
		// the replica is already gone, so no code is swallowed here: any failure
		// keeps the id in state and lets the next apply retry.
		resp.Diagnostics.AddError(
			"Error deleting TEO function replica",
			fmt.Sprintf("DeleteFunctionReplica failed: %v", err),
		)
		return
	}
}

// ImportState parses `<zone_id>#<function_id>#<replica_name>` and fills in the
// identity attributes so that the following Read has everything it needs.
func (r *TeoFunctionReplicaV4Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	zoneId, functionId, replicaName, ok := splitTeoFunctionReplicaV4Id(req.ID)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid import id",
			fmt.Sprintf("expected `<zone_id>%s<function_id>%s<replica_name>`, got %q",
				tccommon.FILED_SP, tccommon.FILED_SP, req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("zone_id"), zoneId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("function_id"), functionId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("replica_name"), replicaName)...)
}

// buildTeoFunctionReplicaV4Id joins the identity segments of a replica.
func buildTeoFunctionReplicaV4Id(zoneId, functionId, replicaName string) string {
	return strings.Join([]string{zoneId, functionId, replicaName}, tccommon.FILED_SP)
}

// splitTeoFunctionReplicaV4Id parses `<zone_id>#<function_id>#<replica_name>`.
func splitTeoFunctionReplicaV4Id(id string) (zoneId, functionId, replicaName string, ok bool) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 3 {
		return "", "", "", false
	}
	for _, p := range parts {
		if p == "" {
			return "", "", "", false
		}
	}
	return parts[0], parts[1], parts[2], true
}

// read fetches the single replica located by the identity segments.
//
// An API failure (nil response / nil Response) is returned as an error; a
// successful query without an exact replica-name match returns (nil, nil),
// which means "the replica does not exist".
func (r *TeoFunctionReplicaV4Resource) read(ctx context.Context, zoneId, functionId, replicaName string) (*teov20220901.FunctionReplica, error) {
	request := teov20220901.NewDescribeFunctionReplicasRequest()
	request.ZoneId = helper.String(zoneId)
	request.FunctionId = helper.String(functionId)
	// The replica-name filter is a fuzzy match, so the result set is narrowed
	// down to the exact name below instead of trusting the filter.
	request.Filters = []*teov20220901.AdvancedFilter{
		{
			Name:   helper.String("replica-name"),
			Values: []*string{helper.String(replicaName)},
		},
	}
	request.Limit = helper.Int64(200)

	response, err := r.Client().UseTeoClient().DescribeFunctionReplicasWithContext(ctx, request)
	if err != nil {
		return nil, err
	}
	if response == nil || response.Response == nil {
		return nil, fmt.Errorf("DescribeFunctionReplicas returned an empty response for %s", buildTeoFunctionReplicaV4Id(zoneId, functionId, replicaName))
	}

	for _, replica := range response.Response.FunctionReplicas {
		if replica == nil || replica.ReplicaName == nil {
			continue
		}
		if *replica.ReplicaName == replicaName {
			return replica, nil
		}
	}

	return nil, nil
}

// flushModel copies the cloud owned fields into the model. replica_name always
// keeps the configured value because it is the identity segment.
func (r *TeoFunctionReplicaV4Resource) flushModel(m *TeoFunctionReplicaV4ResourceModel, zoneId, functionId, replicaName string, replica *teov20220901.FunctionReplica) {
	if replica == nil {
		return
	}

	m.ZoneId = types.StringValue(zoneId)
	m.FunctionId = types.StringValue(functionId)
	m.ReplicaName = types.StringValue(replicaName)
	if replica.Content != nil {
		m.Content = types.StringValue(*replica.Content)
	}
	m.Remark = types.StringPointerValue(replica.Remark)
}
