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

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/fw"
)

// teoFunctionReplicaV3IdSeparator separates the three identity segments of this
// resource. Replica names only allow `a-z`, `0-9` and `-`, so it never appears
// inside a segment.
const teoFunctionReplicaV3IdSeparator = "#"

// _ resource.Resource is intentionally omitted: the factory function
// NewTeoFunctionReplicaV3Resource returns resource.Resource, which already
// performs the implicit compile-time check at the return statement. Only the
// extended-interface assertions are kept here.
var (
	_ resource.ResourceWithConfigure   = &TeoFunctionReplicaV3Resource{}
	_ resource.ResourceWithImportState = &TeoFunctionReplicaV3Resource{}
)

// NewTeoFunctionReplicaV3Resource is the factory referenced by
// tencentcloud/framework/registry.go to register this resource.
func NewTeoFunctionReplicaV3Resource() resource.Resource {
	return &TeoFunctionReplicaV3Resource{}
}

// TeoFunctionReplicaV3Resource implements resource.Resource for
// tencentcloud_teo_function_replica_v3.
//
// Identity: the create API returns no resource level identifier, so the
// identity is the composite `<zone_id>#<function_id>#<replica_name>`. All three
// segments are immutable because they are required to locate the replica and
// the cloud API has no rename capability.
type TeoFunctionReplicaV3Resource struct {
	fw.ResourceWithConfigure
}

// TeoFunctionReplicaV3ResourceModel maps the schema attributes.
type TeoFunctionReplicaV3ResourceModel struct {
	Id          types.String `tfsdk:"id"`
	ZoneId      types.String `tfsdk:"zone_id"`
	FunctionId  types.String `tfsdk:"function_id"`
	ReplicaName types.String `tfsdk:"replica_name"`
	Content     types.String `tfsdk:"content"`
	Remark      types.String `tfsdk:"remark"`
	CreatedOn   types.String `tfsdk:"created_on"`
	ModifiedOn  types.String `tfsdk:"modified_on"`
}

func (r *TeoFunctionReplicaV3Resource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "tencentcloud_teo_function_replica_v3"
}

func (r *TeoFunctionReplicaV3Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource to create a TEO edge function replica. " +
			"The resource manages exactly one replica of one function; deleting it only removes that " +
			"single replica and never touches the parent function or its sibling replicas.",
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
				Description: "ID of the zone to which the function belongs.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"function_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the edge function.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"replica_name": schema.StringAttribute{
				Required: true,
				Description: "Name of the edge function replica. 1-50 characters, only `a-z`, `0-9` and `-` are " +
					"allowed; `-` cannot be used alone, repeated, or placed at the beginning or the end. " +
					"The name must be unique within the same function.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content": schema.StringAttribute{
				Required: true,
				Description: "Content of the edge function replica. Only JavaScript code is supported, and the " +
					"maximum size is 5MB. Note that the content is stored in the Terraform state file, so the " +
					"state file must be protected accordingly.",
			},
			"remark": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the edge function replica. Maximum 50 characters.",
			},
			"created_on": schema.StringAttribute{
				Computed:    true,
				Description: "Creation time of the edge function replica.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified_on": schema.StringAttribute{
				Computed:    true,
				Description: "Last modification time of the edge function replica.",
			},
		},
	}
}

// Create handles `terraform apply` for new resources.
func (r *TeoFunctionReplicaV3Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TeoFunctionReplicaV3ResourceModel
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

	service := NewTeoService(r.Client())
	if err := service.CreateTeoFunctionReplica(
		ctx,
		zoneId,
		functionId,
		replicaName,
		plan.Content.ValueString(),
		stringPtrFromValue(plan.Remark),
	); err != nil {
		resp.Diagnostics.AddError(
			"Error creating TEO function replica",
			fmt.Sprintf("CreateFunctionReplica failed for zone %s function %s replica %s: %v",
				zoneId, functionId, replicaName, err),
		)
		return
	}

	// The create API only returns a request id, so the identity is composed from
	// the inputs. It is written to state before any read back so that a failing
	// read back still leaves the replica manageable by `terraform destroy`.
	plan.Id = types.StringValue(buildTeoFunctionReplicaV3Id(zoneId, functionId, replicaName))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := service.DescribeTeoFunctionReplica(ctx, zoneId, functionId, replicaName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO function replica after create",
			fmt.Sprintf("DescribeFunctionReplicas failed for zone %s function %s replica %s: %v",
				zoneId, functionId, replicaName, err),
		)
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"TEO function replica not readable after create",
			fmt.Sprintf("replica %s of function %s in zone %s was created but cannot be read back",
				replicaName, functionId, zoneId),
		)
		return
	}

	r.refreshModel(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read handles `terraform refresh` and post-apply state refresh.
func (r *TeoFunctionReplicaV3Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TeoFunctionReplicaV3ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, functionId, replicaName, ok := splitTeoFunctionReplicaV3Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO function replica id is broken, id is %q (expected `<zone_id>#<function_id>#<replica_name>`)",
				state.Id.ValueString()),
		)
		return
	}

	if r.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	service := NewTeoService(r.Client())
	item, err := service.DescribeTeoFunctionReplica(ctx, zoneId, functionId, replicaName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO function replica",
			fmt.Sprintf("DescribeFunctionReplicas failed for zone %s function %s replica %s: %v",
				zoneId, functionId, replicaName, err),
		)
		return
	}

	// A missing replica is the only confirmed "does not exist" signal: the list
	// query succeeded and no item carries this exact name. Nil responses and API
	// errors are handled above and never remove the resource from state.
	if item == nil {
		resp.Diagnostics.AddWarning(
			"Resource not found",
			fmt.Sprintf("TEO function replica %s not found, removing from state.", state.Id.ValueString()),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	state.ZoneId = types.StringValue(zoneId)
	state.FunctionId = types.StringValue(functionId)
	state.ReplicaName = types.StringValue(replicaName)
	r.refreshModel(&state, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles in-place updates. The cloud API can only modify content and
// remark; every other attribute forces a replacement.
func (r *TeoFunctionReplicaV3Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TeoFunctionReplicaV3ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state TeoFunctionReplicaV3ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, functionId, replicaName, ok := splitTeoFunctionReplicaV3Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO function replica id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	if r.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	service := NewTeoService(r.Client())
	if err := service.ModifyTeoFunctionReplica(
		ctx,
		zoneId,
		functionId,
		replicaName,
		plan.Content.ValueString(),
		stringPtrFromValue(plan.Remark),
	); err != nil {
		resp.Diagnostics.AddError(
			"Error updating TEO function replica",
			fmt.Sprintf("ModifyFunctionReplica failed for zone %s function %s replica %s: %v",
				zoneId, functionId, replicaName, err),
		)
		return
	}

	plan.Id = state.Id
	plan.ZoneId = types.StringValue(zoneId)
	plan.FunctionId = types.StringValue(functionId)
	plan.ReplicaName = types.StringValue(replicaName)

	item, err := service.DescribeTeoFunctionReplica(ctx, zoneId, functionId, replicaName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading TEO function replica after update",
			fmt.Sprintf("DescribeFunctionReplicas failed for zone %s function %s replica %s: %v",
				zoneId, functionId, replicaName, err),
		)
		return
	}
	if item == nil {
		resp.Diagnostics.AddError(
			"TEO function replica not readable after update",
			fmt.Sprintf("replica %s of function %s in zone %s was updated but cannot be read back",
				replicaName, functionId, zoneId),
		)
		return
	}

	r.refreshModel(&plan, item)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles `terraform destroy` for a single replica.
func (r *TeoFunctionReplicaV3Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TeoFunctionReplicaV3ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneId, functionId, replicaName, ok := splitTeoFunctionReplicaV3Id(state.Id.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("TEO function replica id is broken, id is %q", state.Id.ValueString()),
		)
		return
	}

	if r.Client() == nil {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider client is not available. This is unexpected; please report it.",
		)
		return
	}

	service := NewTeoService(r.Client())
	if err := service.DeleteTeoFunctionReplica(ctx, zoneId, functionId, replicaName); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting TEO function replica",
			fmt.Sprintf("DeleteFunctionReplica failed for zone %s function %s replica %s: %v",
				zoneId, functionId, replicaName, err),
		)
		return
	}
}

// ImportState supports `terraform import` using the
// `<zone_id>#<function_id>#<replica_name>` form.
func (r *TeoFunctionReplicaV3Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	zoneId, functionId, replicaName, ok := splitTeoFunctionReplicaV3Id(req.ID)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected import ID in the form `<zone_id>#<function_id>#<replica_name>`, got %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("zone_id"), zoneId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("function_id"), functionId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("replica_name"), replicaName)...)
}

// refreshModel copies cloud-owned fields into the model.
func (r *TeoFunctionReplicaV3Resource) refreshModel(m *TeoFunctionReplicaV3ResourceModel, item *teov20220901.FunctionReplica) {
	if item == nil {
		return
	}
	if item.Content != nil {
		m.Content = types.StringValue(*item.Content)
	} else {
		m.Content = types.StringNull()
	}
	if item.Remark != nil && *item.Remark != "" {
		m.Remark = types.StringValue(*item.Remark)
	} else {
		m.Remark = types.StringNull()
	}
	if item.CreatedOn != nil {
		m.CreatedOn = types.StringValue(*item.CreatedOn)
	} else {
		m.CreatedOn = types.StringNull()
	}
	if item.ModifiedOn != nil {
		m.ModifiedOn = types.StringValue(*item.ModifiedOn)
	} else {
		m.ModifiedOn = types.StringNull()
	}
}

// buildTeoFunctionReplicaV3Id joins the three identity segments.
func buildTeoFunctionReplicaV3Id(zoneId, functionId, replicaName string) string {
	return strings.Join([]string{zoneId, functionId, replicaName}, teoFunctionReplicaV3IdSeparator)
}

// splitTeoFunctionReplicaV3Id splits `<zone_id>#<function_id>#<replica_name>`.
// Replica names may not contain the separator, so exactly three non empty
// segments are expected.
func splitTeoFunctionReplicaV3Id(id string) (string, string, string, bool) {
	parts := strings.Split(id, teoFunctionReplicaV3IdSeparator)
	if len(parts) != 3 {
		return "", "", "", false
	}
	for _, part := range parts {
		if part == "" {
			return "", "", "", false
		}
	}
	return parts[0], parts[1], parts[2], true
}

// stringPtrFromValue converts a types.String into *string, returning nil for
// null and unknown values so the cloud side keeps its default.
func stringPtrFromValue(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	s := v.ValueString()
	return &s
}
